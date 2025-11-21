import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import * as path from 'path';
import { ApiClient, ImportRequest } from '../api/client';
import { CommandContext } from '../commands/index';
import { Logger } from '../utils/logger';

/**
 * Multi-step import wizard using WebView panel
 */
export class ImportWizard {
  private panel: vscode.WebviewPanel | undefined;
  private context: vscode.ExtensionContext;
  private apiClient: ApiClient;
  private providers: CommandContext;

  constructor(context: vscode.ExtensionContext, apiClient: ApiClient, providers: CommandContext) {
    this.context = context;
    this.apiClient = apiClient;
    this.providers = providers;
  }

  async show(): Promise<void> {
    Logger.info('Opening import wizard');

    // Create and show panel
    this.panel = vscode.window.createWebviewPanel(
      'gorevImportWizard',
      t('import.wizardTitle'),
      vscode.ViewColumn.One,
      {
        enableScripts: true,
        localResourceRoots: [
          vscode.Uri.joinPath(this.context.extensionUri, 'media'),
          vscode.Uri.joinPath(this.context.extensionUri, 'dist')
        ]
      }
    );

    // Set panel icon
    this.panel.iconPath = {
      light: vscode.Uri.joinPath(this.context.extensionUri, 'media', 'import-light.svg'),
      dark: vscode.Uri.joinPath(this.context.extensionUri, 'media', 'import-dark.svg')
    };

    // Handle messages from webview
    this.panel.webview.onDidReceiveMessage(
      async (message) => {
        await this.handleMessage(message);
      },
      undefined,
      this.context.subscriptions
    );

    // Handle panel disposal
    this.panel.onDidDispose(
      () => {
        this.panel = undefined;
      },
      null,
      this.context.subscriptions
    );

    // Set initial HTML content
    this.panel.webview.html = this.getWebviewContent();

    // Load initial data
    await this.loadInitialData();
  }

  private async handleMessage(message: { command: string; [key: string]: unknown }): Promise<void> {
    Logger.debug('Import wizard received message', message);

    switch (message.command) {
      case 'selectFile':
        await this.selectImportFile();
        break;

      case 'analyzeFile':
        await this.analyzeImportFile(message.filePath as string);
        break;

      case 'loadProjects':
        await this.loadProjects();
        break;

      case 'performDryRun':
        await this.performDryRun(message.options as Record<string, unknown>);
        break;

      case 'startImport':
        await this.startImport(message.options as Record<string, unknown>);
        break;

      case 'cancel':
        this.panel?.dispose();
        break;

      default:
        Logger.warn('Unknown message command in import wizard', message.command);
    }
  }

  private async loadInitialData(): Promise<void> {
    if (!this.panel) return;

    try {
      // Load projects for project mapping
      await this.loadProjects();
    } catch (error) {
      Logger.error('Failed to load initial import wizard data', error);
      this.panel.webview.postMessage({
        command: 'showError',
        message: t('error.loadFailed', { error: String(error) })
      });
    }
  }

  private async selectImportFile(): Promise<void> {
    if (!this.panel) return;

    const fileUri = await vscode.window.showOpenDialog({
      canSelectFiles: true,
      canSelectFolders: false,
      canSelectMany: false,
      filters: {
        'JSON Files': ['json'],
        'CSV Files': ['csv'],
        'All Files': ['*']
      },
      title: t('import.selectFile')
    });

    if (fileUri && fileUri[0]) {
      const filePath = fileUri[0].fsPath;
      this.panel.webview.postMessage({
        command: 'setSelectedFile',
        path: filePath
      });

      // Automatically analyze the file
      await this.analyzeImportFile(filePath);
    }
  }

  private async analyzeImportFile(filePath: string): Promise<void> {
    if (!this.panel) return;

    try {
      // Show analysis in progress
      this.panel.webview.postMessage({
        command: 'analysisStarted'
      });

      // Determine format from file extension
      const format = path.extname(filePath).toLowerCase() === '.csv' ? 'csv' : 'json';

      // Call dry run to analyze the file
      const dryRunOptions = {
        file_path: filePath,
        format: format,
        dry_run: true,
        conflict_resolution: 'skip',
        project_mapping: {}
      };

      // Call REST API import endpoint
      const result = await this.apiClient.importData(dryRunOptions);

      if (!result.success) {
        throw new Error(result.message || 'Analysis failed');
      }

      // Parse analysis results
      const analysisResult = this.parseAnalysisResult(result.message || '');
      
      this.panel.webview.postMessage({
        command: 'analysisCompleted',
        result: analysisResult,
        format: format
      });

    } catch (error) {
      Logger.error('File analysis failed', error);
      this.panel.webview.postMessage({
        command: 'analysisFailed',
        error: String(error)
      });
    }
  }

  private async loadProjects(): Promise<void> {
    if (!this.panel) return;

    try {
      const result = await this.apiClient.getProjects();
      if (!result.success || !result.data) {
        throw new Error('Failed to load projects');
      }

      this.panel.webview.postMessage({
        command: 'setProjects',
        projects: result.data
      });

    } catch (error) {
      Logger.error('Failed to load projects for import wizard', error);
    }
  }

  private async performDryRun(options: Record<string, unknown>): Promise<void> {
    if (!this.panel) return;

    try {
      this.panel.webview.postMessage({
        command: 'dryRunStarted'
      });

      // Perform dry run
      const dryRunOptions = {
        ...options,
        dry_run: true
      } as ImportRequest;

      // Call REST API import endpoint
      const result = await this.apiClient.importData(dryRunOptions);

      if (!result.success) {
        throw new Error(result.message || 'Dry run failed');
      }

      // Parse dry run results
      const dryRunResult = this.parseDryRunResult(result.message || '');
      
      this.panel.webview.postMessage({
        command: 'dryRunCompleted',
        result: dryRunResult
      });

    } catch (error) {
      Logger.error('Dry run failed', error);
      this.panel.webview.postMessage({
        command: 'dryRunFailed',
        error: String(error)
      });
    }
  }

  private async startImport(options: Record<string, unknown>): Promise<void> {
    if (!this.panel) return;

    try {
      // Show progress in webview
      this.panel.webview.postMessage({
        command: 'importStarted'
      });

      // Perform import with progress updates
      await vscode.window.withProgress({
        location: vscode.ProgressLocation.Notification,
        title: t('import.inProgress'),
        cancellable: false
      }, async (progress) => {
        progress.report({ increment: 10, message: t('import.preparing') });

        // Remove dry_run flag for actual import
        const importOptions = {
          ...options,
          dry_run: false
        } as ImportRequest;

        // Call REST API import endpoint
        const result = await this.apiClient.importData(importOptions);

        progress.report({ increment: 80, message: t('import.processing') });

        if (!result.success) {
          throw new Error(result.message || 'Import failed');
        }

        progress.report({ increment: 100, message: t('import.complete') });

        // Parse import results
        const importResult = this.parseImportResult(result.message || '');

        // Notify webview of success
        this.panel?.webview.postMessage({
          command: 'importCompleted',
          result: importResult
        });

        // Refresh all views after successful import
        await this.refreshViews();

        // Show success message
        vscode.window.showInformationMessage(
          t('import.success', { 
            tasks: importResult.tasksImported,
            projects: importResult.projectsImported
          })
        );

        // Close wizard after successful import
        setTimeout(() => {
          this.panel?.dispose();
        }, 3000);
      });

    } catch (error) {
      Logger.error('Import failed', error);
      
      this.panel.webview.postMessage({
        command: 'importFailed',
        error: String(error)
      });

      vscode.window.showErrorMessage(
        t('error.importFailed', { error: String(error) })
      );
    }
  }

  private async refreshViews(): Promise<void> {
    try {
      // Refresh all providers sequentially to avoid overwhelming the MCP server
      await this.providers.gorevTreeProvider.refresh();
      await this.providers.projeTreeProvider.refresh();
      await this.providers.templateTreeProvider.refresh();
    } catch (error) {
      Logger.warn('Failed to refresh some views after import', error);
    }
  }

  private parseAnalysisResult(text: string): {
    totalTasks: number;
    totalProjects: number;
    warnings: string[];
    errors: string[];
    preview: string[];
  } {
    // Parse the analysis result from MCP tool response
    const result = {
      totalTasks: 0,
      totalProjects: 0,
      warnings: [] as string[],
      errors: [] as string[],
      preview: [] as string[]
    };

    try {
      // Try to parse as JSON if it looks like structured data
      if (text.trim().startsWith('{')) {
        const parsed = JSON.parse(text);
        return { ...result, ...parsed };
      }

      // Parse text format
      const lines = text.split('\n');
      for (const line of lines) {
        if (line.includes('toplam') && line.includes('görev')) {
          const match = line.match(/(\d+)/);
          if (match) result.totalTasks = parseInt(match[1]);
        }
        if (line.includes('toplam') && line.includes('proje')) {
          const match = line.match(/(\d+)/);
          if (match) result.totalProjects = parseInt(match[1]);
        }
        if (line.includes('uyarı') || line.includes('warning')) {
          result.warnings.push(line);
        }
        if (line.includes('hata') || line.includes('error')) {
          result.errors.push(line);
        }
      }
    } catch (error) {
      Logger.warn('Failed to parse analysis result', error);
      result.errors.push('Failed to parse analysis result');
    }

    return result;
  }

  private parseDryRunResult(text: string): {
    tasksToImport: number;
    projectsToImport: number;
    conflicts: string[];
    warnings: string[];
    preview: string[];
  } {
    const result = {
      tasksToImport: 0,
      projectsToImport: 0,
      conflicts: [] as string[],
      warnings: [] as string[],
      preview: [] as string[]
    };

    try {
      // Parse dry run results from text
      const lines = text.split('\n');
      for (const line of lines) {
        if (line.includes('import edilecek') && line.includes('görev')) {
          const match = line.match(/(\d+)/);
          if (match) result.tasksToImport = parseInt(match[1]);
        }
        if (line.includes('import edilecek') && line.includes('proje')) {
          const match = line.match(/(\d+)/);
          if (match) result.projectsToImport = parseInt(match[1]);
        }
        if (line.includes('çakışma') || line.includes('conflict')) {
          result.conflicts.push(line);
        }
        if (line.includes('uyarı') || line.includes('warning')) {
          result.warnings.push(line);
        }
      }
    } catch (error) {
      Logger.warn('Failed to parse dry run result', error);
    }

    return result;
  }

  private parseImportResult(text: string): {
    tasksImported: number;
    projectsImported: number;
    skipped: number;
    errors: string[];
  } {
    const result = {
      tasksImported: 0,
      projectsImported: 0,
      skipped: 0,
      errors: [] as string[]
    };

    try {
      const lines = text.split('\n');
      for (const line of lines) {
        if (line.includes('import edildi') && line.includes('görev')) {
          const match = line.match(/(\d+)/);
          if (match) result.tasksImported = parseInt(match[1]);
        }
        if (line.includes('import edildi') && line.includes('proje')) {
          const match = line.match(/(\d+)/);
          if (match) result.projectsImported = parseInt(match[1]);
        }
        if (line.includes('atlandı') || line.includes('skipped')) {
          const match = line.match(/(\d+)/);
          if (match) result.skipped = parseInt(match[1]);
        }
        if (line.includes('hata') || line.includes('error')) {
          result.errors.push(line);
        }
      }
    } catch (error) {
      Logger.warn('Failed to parse import result', error);
    }

    return result;
  }

  private parseProjectList(text: string): {id: string, name: string, isActive: boolean}[] {
    const projects: {id: string, name: string, isActive: boolean}[] = [];
    const lines = text.split('\n');
    
    for (const line of lines) {
      const match = line.match(/^\s*-\s*(.+?)\s*\(ID:\s*([^)]+)\)\s*(✅)?/);
      if (match) {
        projects.push({
          id: match[2].trim(),
          name: match[1].trim(),
          isActive: !!match[3]
        });
      }
    }
    
    return projects;
  }

  private getWebviewContent(): string {
    const nonce = this.getNonce();
    
    return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; script-src 'nonce-${nonce}';">
    <title>${t('import.wizardTitle')}</title>
    <style>
        body {
            font-family: var(--vscode-font-family);
            color: var(--vscode-foreground);
            background-color: var(--vscode-editor-background);
            padding: 20px;
            margin: 0;
        }
        .step {
            display: none;
            animation: fadeIn 0.3s ease-in;
        }
        .step.active {
            display: block;
        }
        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
        .step-header {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 1px solid var(--vscode-panel-border);
        }
        .step-indicator {
            background: var(--vscode-button-background);
            color: var(--vscode-button-foreground);
            border-radius: 50%;
            width: 30px;
            height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 15px;
            font-weight: bold;
        }
        .step-title {
            font-size: 1.2em;
            font-weight: bold;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        .form-group input, .form-group select, .form-group textarea {
            width: 100%;
            padding: 8px;
            border: 1px solid var(--vscode-input-border);
            background: var(--vscode-input-background);
            color: var(--vscode-input-foreground);
            border-radius: 4px;
            box-sizing: border-box;
        }
        .form-group input[type="radio"] {
            width: auto;
            margin-right: 8px;
        }
        .radio-group {
            display: flex;
            flex-direction: column;
            gap: 10px;
            margin-top: 10px;
        }
        .radio-item {
            display: flex;
            align-items: center;
        }
        .button-group {
            display: flex;
            justify-content: space-between;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid var(--vscode-panel-border);
        }
        .btn {
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        .btn-primary {
            background: var(--vscode-button-background);
            color: var(--vscode-button-foreground);
        }
        .btn-primary:hover {
            background: var(--vscode-button-hoverBackground);
        }
        .btn-secondary {
            background: var(--vscode-button-secondaryBackground);
            color: var(--vscode-button-secondaryForeground);
        }
        .btn-secondary:hover {
            background: var(--vscode-button-secondaryHoverBackground);
        }
        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
        }
        .error {
            color: var(--vscode-errorForeground);
            background: var(--vscode-inputValidation-errorBackground);
            border: 1px solid var(--vscode-inputValidation-errorBorder);
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .warning {
            color: var(--vscode-warningForeground);
            background: var(--vscode-inputValidation-warningBackground);
            border: 1px solid var(--vscode-inputValidation-warningBorder);
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .info {
            color: var(--vscode-infoForeground);
            background: var(--vscode-inputValidation-infoBackground);
            border: 1px solid var(--vscode-inputValidation-infoBorder);
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .success {
            color: var(--vscode-testing-iconPassed);
            background: var(--vscode-inputValidation-infoBackground);
            border: 1px solid var(--vscode-testing-iconPassed);
            padding: 10px;
            border-radius: 4px;
            margin: 10px 0;
        }
        .file-input-group {
            display: flex;
            gap: 10px;
            align-items: end;
        }
        .file-input-group input {
            flex: 1;
        }
        .analysis-result {
            background: var(--vscode-editor-background);
            border: 1px solid var(--vscode-panel-border);
            border-radius: 4px;
            padding: 15px;
            margin: 20px 0;
        }
        .stat-item {
            display: flex;
            justify-content: space-between;
            margin: 5px 0;
        }
        .stat-label {
            font-weight: bold;
        }
        .stat-value {
            color: var(--vscode-button-background);
        }
        .project-mapping {
            background: var(--vscode-editor-background);
            border: 1px solid var(--vscode-panel-border);
            border-radius: 4px;
            padding: 15px;
            margin: 20px 0;
        }
        .mapping-item {
            display: grid;
            grid-template-columns: 1fr auto 1fr;
            gap: 10px;
            align-items: center;
            margin: 10px 0;
            padding: 10px;
            background: var(--vscode-input-background);
            border-radius: 4px;
        }
        .arrow {
            text-align: center;
            color: var(--vscode-descriptionForeground);
        }
        .progress-bar {
            width: 100%;
            height: 4px;
            background: var(--vscode-progressBar-background);
            border-radius: 2px;
            overflow: hidden;
            margin: 20px 0;
        }
        .progress-fill {
            height: 100%;
            background: var(--vscode-progressBar-foreground);
            width: 0%;
            transition: width 0.3s ease;
            animation: progressAnimation 1.5s ease-in-out infinite;
        }
        @keyframes progressAnimation {
            0% { transform: translateX(-100%); }
            100% { transform: translateX(100%); }
        }
        .conflict-list {
            max-height: 200px;
            overflow-y: auto;
            border: 1px solid var(--vscode-panel-border);
            border-radius: 4px;
            padding: 10px;
            margin: 10px 0;
            background: var(--vscode-input-background);
        }
        .conflict-item {
            padding: 5px 0;
            border-bottom: 1px solid var(--vscode-panel-border);
        }
        .conflict-item:last-child {
            border-bottom: none;
        }
    </style>
</head>
<body>
    <!-- Step 1: File Selection -->
    <div class="step active" id="step1">
        <div class="step-header">
            <div class="step-indicator">1</div>
            <div class="step-title">${t('import.step1.title')}</div>
        </div>
        
        <div class="form-group">
            <label for="filePath">${t('import.selectFile')}</label>
            <div class="file-input-group">
                <input type="text" id="filePath" readonly placeholder="${t('import.noFileSelected')}">
                <button class="btn btn-secondary" onclick="selectFile()">${t('import.browse')}</button>
            </div>
        </div>

        <div id="analysisStatus" style="display: none;">
            <div class="progress-bar">
                <div class="progress-fill"></div>
            </div>
            <div id="analysisMessage">${t('import.analyzing')}</div>
        </div>

        <div id="analysisResult" class="analysis-result" style="display: none;">
            <!-- Analysis results filled dynamically -->
        </div>

        <div id="analysisError" class="error" style="display: none;"></div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="cancel()">${t('common.cancel')}</button>
            <button class="btn btn-primary" onclick="nextStep(2)" id="nextBtn1" disabled>${t('common.next')}</button>
        </div>
    </div>

    <!-- Step 2: Project Mapping -->
    <div class="step" id="step2">
        <div class="step-header">
            <div class="step-indicator">2</div>
            <div class="step-title">${t('import.step2.title')}</div>
        </div>

        <div class="form-group">
            <label>${t('import.projectMapping')}</label>
            <div id="projectMappings" class="project-mapping">
                <!-- Project mappings filled dynamically -->
            </div>
        </div>

        <div class="form-group">
            <label>${t('import.conflictResolution')}</label>
            <div class="radio-group">
                <div class="radio-item">
                    <input type="radio" id="conflictSkip" name="conflictResolution" value="skip" checked>
                    <label for="conflictSkip">${t('import.conflict.skip')}</label>
                </div>
                <div class="radio-item">
                    <input type="radio" id="conflictOverwrite" name="conflictResolution" value="overwrite">
                    <label for="conflictOverwrite">${t('import.conflict.overwrite')}</label>
                </div>
                <div class="radio-item">
                    <input type="radio" id="conflictMerge" name="conflictResolution" value="merge">
                    <label for="conflictMerge">${t('import.conflict.merge')}</label>
                </div>
            </div>
        </div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="previousStep(1)">${t('common.back')}</button>
            <button class="btn btn-primary" onclick="nextStep(3)">${t('common.next')}</button>
        </div>
    </div>

    <!-- Step 3: Dry Run Preview -->
    <div class="step" id="step3">
        <div class="step-header">
            <div class="step-indicator">3</div>
            <div class="step-title">${t('import.step3.title')}</div>
        </div>

        <div class="form-group">
            <button class="btn btn-primary" onclick="performDryRun()" id="dryRunBtn">${t('import.runPreview')}</button>
        </div>

        <div id="dryRunStatus" style="display: none;">
            <div class="progress-bar">
                <div class="progress-fill"></div>
            </div>
            <div id="dryRunMessage">${t('import.previewRunning')}</div>
        </div>

        <div id="dryRunResult" class="analysis-result" style="display: none;">
            <!-- Dry run results filled dynamically -->
        </div>

        <div id="conflictsList" class="conflict-list" style="display: none;">
            <!-- Conflicts filled dynamically -->
        </div>

        <div id="dryRunError" class="error" style="display: none;"></div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="previousStep(2)">${t('common.back')}</button>
            <button class="btn btn-primary" onclick="nextStep(4)" id="nextBtn3" disabled>${t('common.next')}</button>
        </div>
    </div>

    <!-- Step 4: Import Execution -->
    <div class="step" id="step4">
        <div class="step-header">
            <div class="step-indicator">4</div>
            <div class="step-title">${t('import.step4.title')}</div>
        </div>

        <div class="info">
            <p>${t('import.finalWarning')}</p>
        </div>

        <div id="importProgress" style="display: none;">
            <div class="progress-bar">
                <div class="progress-fill"></div>
            </div>
            <div id="importMessage">${t('import.preparing')}</div>
        </div>

        <div id="importResult" style="display: none;"></div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="previousStep(3)" id="backBtn4">${t('common.back')}</button>
            <button class="btn btn-primary" onclick="startImport()" id="importBtn">${t('import.start')}</button>
        </div>
    </div>

    <script nonce="${nonce}">
        const vscode = acquireVsCodeApi();
        let currentStep = 1;
        let projects = [];
        let analysisData = null;
        let dryRunData = null;
        let selectedFormat = 'json';

        // Message handling
        window.addEventListener('message', event => {
            const message = event.data;
            
            switch (message.command) {
                case 'setSelectedFile':
                    document.getElementById('filePath').value = message.path;
                    showAnalysisProgress(true);
                    break;
                    
                case 'setProjects':
                    projects = message.projects;
                    break;
                    
                case 'analysisStarted':
                    showAnalysisProgress(true);
                    break;
                    
                case 'analysisCompleted':
                    analysisData = message.result;
                    selectedFormat = message.format;
                    showAnalysisResult(message.result);
                    showAnalysisProgress(false);
                    enableNextButton(1, true);
                    break;
                    
                case 'analysisFailed':
                    showAnalysisError(message.error);
                    showAnalysisProgress(false);
                    break;
                    
                case 'dryRunStarted':
                    showDryRunProgress(true);
                    break;
                    
                case 'dryRunCompleted':
                    dryRunData = message.result;
                    showDryRunResult(message.result);
                    showDryRunProgress(false);
                    enableNextButton(3, true);
                    break;
                    
                case 'dryRunFailed':
                    showDryRunError(message.error);
                    showDryRunProgress(false);
                    break;
                    
                case 'importStarted':
                    showImportProgress(true);
                    break;
                    
                case 'importCompleted':
                    showImportResult(true, message.result);
                    showImportProgress(false);
                    break;
                    
                case 'importFailed':
                    showImportResult(false, null, message.error);
                    showImportProgress(false);
                    break;
                    
                case 'showError':
                    showError(message.message);
                    break;
            }
        });

        function selectFile() {
            vscode.postMessage({ command: 'selectFile' });
        }

        function nextStep(step) {
            if (step === 2) {
                updateProjectMappings();
            }
            
            document.getElementById('step' + currentStep).classList.remove('active');
            currentStep = step;
            document.getElementById('step' + currentStep).classList.add('active');
        }

        function previousStep(step) {
            document.getElementById('step' + currentStep).classList.remove('active');
            currentStep = step;
            document.getElementById('step' + currentStep).classList.add('active');
        }

        function updateProjectMappings() {
            const mappingsDiv = document.getElementById('projectMappings');
            
            if (!analysisData || analysisData.totalProjects === 0) {
                mappingsDiv.innerHTML = '<p>${t('import.noProjectMapping')}</p>';
                return;
            }

            let html = '<p>${t('import.projectMapping.description')}</p>';
            
            // Create mapping interfaces for each project in import file
            // This would need to be enhanced based on actual analysis data structure
            for (let i = 0; i < analysisData.totalProjects; i++) {
                html += '<div class="mapping-item">';
                html += '<div><label>Import Project ' + (i + 1) + '</label><input type="text" readonly value="Project from file"></div>';
                html += '<div class="arrow">→</div>';
                html += '<div><select id="projectMapping' + i + '">';
                html += '<option value="">${t('import.createNew')}</option>';
                
                projects.forEach(project => {
                    html += '<option value="' + project.id + '">' + project.name + '</option>';
                });
                
                html += '</select></div>';
                html += '</div>';
            }
            
            mappingsDiv.innerHTML = html;
        }

        function performDryRun() {
            const options = gatherImportOptions();
            options.dry_run = true;
            
            vscode.postMessage({
                command: 'performDryRun',
                options: options
            });
        }

        function startImport() {
            const options = gatherImportOptions();
            
            vscode.postMessage({
                command: 'startImport',
                options: options
            });
        }

        function gatherImportOptions() {
            const filePath = document.getElementById('filePath').value;
            const conflictResolution = document.querySelector('input[name="conflictResolution"]:checked').value;
            
            const projectMapping = {};
            // Gather project mappings
            for (let i = 0; i < (analysisData?.totalProjects || 0); i++) {
                const selectElement = document.getElementById('projectMapping' + i);
                if (selectElement && selectElement.value) {
                    projectMapping['project' + i] = selectElement.value;
                }
            }

            return {
                file_path: filePath,
                format: selectedFormat,
                conflict_resolution: conflictResolution,
                project_mapping: projectMapping
            };
        }

        function showAnalysisProgress(show) {
            const statusDiv = document.getElementById('analysisStatus');
            statusDiv.style.display = show ? 'block' : 'none';
        }

        function showAnalysisResult(result) {
            const resultDiv = document.getElementById('analysisResult');
            
            let html = '<h3>${t('import.analysisResult')}</h3>';
            html += '<div class="stat-item"><span class="stat-label">${t('import.totalTasks')}:</span><span class="stat-value">' + result.totalTasks + '</span></div>';
            html += '<div class="stat-item"><span class="stat-label">${t('import.totalProjects')}:</span><span class="stat-value">' + result.totalProjects + '</span></div>';
            
            if (result.warnings && result.warnings.length > 0) {
                html += '<div class="warning"><strong>${t('import.warnings')}:</strong><ul>';
                result.warnings.forEach(warning => {
                    html += '<li>' + warning + '</li>';
                });
                html += '</ul></div>';
            }
            
            if (result.errors && result.errors.length > 0) {
                html += '<div class="error"><strong>${t('import.errors')}:</strong><ul>';
                result.errors.forEach(error => {
                    html += '<li>' + error + '</li>';
                });
                html += '</ul></div>';
            }
            
            resultDiv.innerHTML = html;
            resultDiv.style.display = 'block';
        }

        function showAnalysisError(error) {
            const errorDiv = document.getElementById('analysisError');
            errorDiv.innerHTML = '<strong>${t('import.analysisFailed')}:</strong> ' + error;
            errorDiv.style.display = 'block';
        }

        function showDryRunProgress(show) {
            const statusDiv = document.getElementById('dryRunStatus');
            const btn = document.getElementById('dryRunBtn');
            
            statusDiv.style.display = show ? 'block' : 'none';
            btn.disabled = show;
            btn.textContent = show ? '${t('import.previewRunning')}' : '${t('import.runPreview')}';
        }

        function showDryRunResult(result) {
            const resultDiv = document.getElementById('dryRunResult');
            
            let html = '<h3>${t('import.previewResult')}</h3>';
            html += '<div class="stat-item"><span class="stat-label">${t('import.tasksToImport')}:</span><span class="stat-value">' + result.tasksToImport + '</span></div>';
            html += '<div class="stat-item"><span class="stat-label">${t('import.projectsToImport')}:</span><span class="stat-value">' + result.projectsToImport + '</span></div>';
            
            resultDiv.innerHTML = html;
            resultDiv.style.display = 'block';
            
            if (result.conflicts && result.conflicts.length > 0) {
                const conflictsDiv = document.getElementById('conflictsList');
                let conflictsHtml = '<h4>${t('import.conflicts')}</h4>';
                result.conflicts.forEach(conflict => {
                    conflictsHtml += '<div class="conflict-item">' + conflict + '</div>';
                });
                conflictsDiv.innerHTML = conflictsHtml;
                conflictsDiv.style.display = 'block';
            }
        }

        function showDryRunError(error) {
            const errorDiv = document.getElementById('dryRunError');
            errorDiv.innerHTML = '<strong>${t('import.previewFailed')}:</strong> ' + error;
            errorDiv.style.display = 'block';
        }

        function showImportProgress(show) {
            const progressDiv = document.getElementById('importProgress');
            const backBtn = document.getElementById('backBtn4');
            const importBtn = document.getElementById('importBtn');
            
            progressDiv.style.display = show ? 'block' : 'none';
            backBtn.disabled = show;
            importBtn.disabled = show;
            importBtn.textContent = show ? '${t('import.importing')}' : '${t('import.start')}';
        }

        function showImportResult(success, result, error) {
            const resultDiv = document.getElementById('importResult');
            const importBtn = document.getElementById('importBtn');
            
            if (success) {
                resultDiv.innerHTML = '<div class="success"><strong>${t('import.completed')}!</strong><br>' +
                    '${t('import.importedSummary')}: ' + result.tasksImported + ' ${t('import.tasks')}, ' + 
                    result.projectsImported + ' ${t('import.projects')}</div>';
                importBtn.textContent = '${t('common.close')}';
                importBtn.onclick = cancel;
            } else {
                resultDiv.innerHTML = '<div class="error"><strong>${t('import.failed')}:</strong> ' + error + '</div>';
                importBtn.textContent = '${t('import.retry')}';
            }
            
            resultDiv.style.display = 'block';
        }

        function enableNextButton(step, enabled) {
            const btn = document.getElementById('nextBtn' + step);
            if (btn) {
                btn.disabled = !enabled;
            }
        }

        function showError(message) {
            const resultDiv = document.getElementById('importResult');
            resultDiv.innerHTML = '<div class="error">' + message + '</div>';
            resultDiv.style.display = 'block';
        }

        function cancel() {
            vscode.postMessage({ command: 'cancel' });
        }

        // Initial setup
        vscode.postMessage({ command: 'loadProjects' });
    </script>
</body>
</html>`;
  }

  private getNonce(): string {
    let text = '';
    const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    for (let i = 0; i < 32; i++) {
      text += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    return text;
  }
}
