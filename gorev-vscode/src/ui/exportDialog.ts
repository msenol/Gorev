import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import * as path from 'path';
import { ApiClient } from '../api/client';
import { Logger } from '../utils/logger';
import { validateExportOptions, estimateExportSize } from '../commands/dataCommands';

/**
 * Multi-step export dialog using WebView panel
 */
export class ExportDialog {
  private panel: vscode.WebviewPanel | undefined;
  private context: vscode.ExtensionContext;
  private apiClient: ApiClient;

  constructor(context: vscode.ExtensionContext, apiClient: ApiClient) {
    this.context = context;
    this.apiClient = apiClient;
  }

  async show(): Promise<void> {
    Logger.info('Opening export dialog');

    // Create and show panel
    this.panel = vscode.window.createWebviewPanel(
      'gorevExportDialog',
      t('export.dialogTitle'),
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
      light: vscode.Uri.joinPath(this.context.extensionUri, 'media', 'export-light.svg'),
      dark: vscode.Uri.joinPath(this.context.extensionUri, 'media', 'export-dark.svg')
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

  private async handleMessage(message: any): Promise<void> {
    Logger.debug('Export dialog received message', message);

    switch (message.command) {
      case 'loadProjects':
        await this.loadProjects();
        break;

      case 'loadTags':
        await this.loadTags();
        break;

      case 'validateOptions':
        await this.validateExportOptions(message.options);
        break;

      case 'estimateSize':
        await this.estimateExportSize(message.options);
        break;

      case 'selectOutputPath':
        await this.selectOutputPath(message.format);
        break;

      case 'startExport':
        await this.startExport(message.options);
        break;

      case 'cancel':
        this.panel?.dispose();
        break;

      default:
        Logger.warn('Unknown message command in export dialog', message.command);
    }
  }

  private async loadInitialData(): Promise<void> {
    if (!this.panel) return;

    try {
      // Load projects and tags for filter options
      await this.loadProjects();
      await this.loadTags();

      // Set default output path
      const defaultPath = this.getDefaultOutputPath();
      this.panel.webview.postMessage({
        command: 'setDefaultPath',
        path: defaultPath
      });

    } catch (error) {
      Logger.error('Failed to load initial export dialog data', error);
      this.panel.webview.postMessage({
        command: 'showError',
        message: t('error.loadFailed', { error: String(error) })
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
      Logger.error('Failed to load projects for export dialog', error);
    }
  }

  private async loadTags(): Promise<void> {
    if (!this.panel) return;

    try {
      // TODO: Replace with REST API getTags() when available
      // For now, using empty array - tags will be available when REST endpoint is added
      const tags: string[] = [];
      this.panel.webview.postMessage({
        command: 'setTags',
        tags: tags
      });

    } catch (error) {
      Logger.error('Failed to load tags for export dialog', error);
    }
  }

  private async validateExportOptions(options: any): Promise<void> {
    if (!this.panel) return;

    const validation = validateExportOptions(options);
    this.panel.webview.postMessage({
      command: 'validationResult',
      isValid: validation.isValid,
      errors: validation.errors
    });
  }

  private async estimateExportSize(options: any): Promise<void> {
    if (!this.panel) return;

    try {
      const sizeEstimate = await estimateExportSize(this.apiClient, options);
      this.panel.webview.postMessage({
        command: 'sizeEstimate',
        size: sizeEstimate
      });
    } catch (error) {
      Logger.error('Failed to estimate export size', error);
    }
  }

  private async selectOutputPath(format: string): Promise<void> {
    if (!this.panel) return;

    const defaultFileName = `gorev-export-${new Date().toISOString().split('T')[0]}.${format}`;
    const saveUri = await vscode.window.showSaveDialog({
      defaultUri: vscode.Uri.file(path.join(vscode.workspace.rootPath || '', defaultFileName)),
      filters: format === 'csv' ? {
        'CSV Files': ['csv'],
        'All Files': ['*']
      } : {
        'JSON Files': ['json'],
        'All Files': ['*']
      },
      title: t('export.selectLocation')
    });

    if (saveUri) {
      this.panel.webview.postMessage({
        command: 'setOutputPath',
        path: saveUri.fsPath
      });
    }
  }

  private async startExport(options: any): Promise<void> {
    if (!this.panel) return;

    try {
      // Show progress in webview
      this.panel.webview.postMessage({
        command: 'exportStarted'
      });

      // Perform export with progress updates
      await vscode.window.withProgress({
        location: vscode.ProgressLocation.Notification,
        title: t('export.inProgress'),
        cancellable: false
      }, async (progress) => {
        progress.report({ increment: 10, message: t('export.preparing') });

        // Call REST API export endpoint
        const result = await this.apiClient.exportData(options);

        progress.report({ increment: 90, message: t('export.completing') });

        if (!result.success) {
          throw new Error(result.message || 'Export failed');
        }

        progress.report({ increment: 100, message: t('export.complete') });

        // Notify webview of success
        this.panel?.webview.postMessage({
          command: 'exportCompleted',
          path: options.output_path
        });

        // Show success message with action buttons
        const openAction = t('export.openFile');
        const openFolderAction = t('export.openFolder');
        
        const action = await vscode.window.showInformationMessage(
          t('export.success', { path: options.output_path }),
          openAction,
          openFolderAction
        );

        if (action === openAction) {
          await vscode.commands.executeCommand('vscode.open', vscode.Uri.file(options.output_path));
        } else if (action === openFolderAction) {
          await vscode.commands.executeCommand('revealFileInOS', vscode.Uri.file(options.output_path));
        }

        // Close dialog after successful export
        setTimeout(() => {
          this.panel?.dispose();
        }, 2000);
      });

    } catch (error) {
      Logger.error('Export failed', error);
      
      this.panel.webview.postMessage({
        command: 'exportFailed',
        error: String(error)
      });

      vscode.window.showErrorMessage(
        t('error.exportFailed', { error: String(error) })
      );
    }
  }

  private getDefaultOutputPath(): string {
    const timestamp = new Date().toISOString().split('T')[0];
    const defaultFileName = `gorev-export-${timestamp}.json`;
    
    if (vscode.workspace.rootPath) {
      return path.join(vscode.workspace.rootPath, defaultFileName);
    } else {
      const downloadsPath = path.join(require('os').homedir(), 'Downloads');
      return path.join(downloadsPath, defaultFileName);
    }
  }

  private parseProjectList(text: string): Array<{id: string, name: string, isActive: boolean}> {
    const projects: Array<{id: string, name: string, isActive: boolean}> = [];
    const lines = text.split('\n');
    
    for (const line of lines) {
      // Parse project lines with format: "- Project Name (ID: xxx) ✅"
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

  private parseTagsFromSummary(text: string): string[] {
    const tags: string[] = [];
    
    // Look for tag patterns in summary text
    const tagMatches = text.match(/etiket[^:]*:\s*([^.\n]+)/gi);
    if (tagMatches) {
      for (const match of tagMatches) {
        const tagList = match.split(':')[1];
        if (tagList) {
          const individualTags = tagList.split(',').map(t => t.trim()).filter(t => t);
          tags.push(...individualTags);
        }
      }
    }
    
    // Remove duplicates and return
    return [...new Set(tags)];
  }

  private getWebviewContent(): string {
    const nonce = this.getNonce();
    
    return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'unsafe-inline'; script-src 'nonce-${nonce}';">
    <title>${t('export.dialogTitle')}</title>
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
        .form-group input[type="checkbox"] {
            width: auto;
            margin-right: 8px;
        }
        .checkbox-group {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 10px;
            margin-top: 10px;
        }
        .checkbox-item {
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
        }
        .size-estimate {
            font-size: 0.9em;
            color: var(--vscode-descriptionForeground);
            margin-top: 5px;
        }
        .file-input-group {
            display: flex;
            gap: 10px;
            align-items: end;
        }
        .file-input-group input {
            flex: 1;
        }
        .export-summary {
            background: var(--vscode-editor-background);
            border: 1px solid var(--vscode-panel-border);
            border-radius: 4px;
            padding: 15px;
            margin: 20px 0;
        }
        .summary-item {
            display: flex;
            justify-content: space-between;
            margin: 5px 0;
        }
        .summary-label {
            font-weight: bold;
        }
    </style>
</head>
<body>
    <!-- Step 1: Format and Basic Options -->
    <div class="step active" id="step1">
        <div class="step-header">
            <div class="step-indicator">1</div>
            <div class="step-title">${t('export.step1.title')}</div>
        </div>
        
        <div class="form-group">
            <label for="format">${t('export.format')}</label>
            <select id="format">
                <option value="json">JSON</option>
                <option value="csv">CSV</option>
            </select>
        </div>

        <div class="form-group">
            <label>${t('export.includeOptions')}</label>
            <div class="checkbox-group">
                <div class="checkbox-item">
                    <input type="checkbox" id="includeCompleted" checked>
                    <label for="includeCompleted">${t('export.includeCompleted')}</label>
                </div>
                <div class="checkbox-item">
                    <input type="checkbox" id="includeDependencies" checked>
                    <label for="includeDependencies">${t('export.includeDependencies')}</label>
                </div>
                <div class="checkbox-item">
                    <input type="checkbox" id="includeTemplates">
                    <label for="includeTemplates">${t('export.includeTemplates')}</label>
                </div>
                <div class="checkbox-item">
                    <input type="checkbox" id="includeAiContext">
                    <label for="includeAiContext">${t('export.includeAiContext')}</label>
                </div>
            </div>
        </div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="cancel()">${t('common.cancel')}</button>
            <button class="btn btn-primary" onclick="nextStep(2)">${t('common.next')}</button>
        </div>
    </div>

    <!-- Step 2: Filters -->
    <div class="step" id="step2">
        <div class="step-header">
            <div class="step-indicator">2</div>
            <div class="step-title">${t('export.step2.title')}</div>
        </div>

        <div class="form-group">
            <label for="projectFilter">${t('export.projectFilter')}</label>
            <select id="projectFilter" multiple style="height: 120px;">
                <!-- Projects loaded dynamically -->
            </select>
            <small>${t('export.projectFilter.help')}</small>
        </div>

        <div class="form-group">
            <label for="statusFilter">${t('export.statusFilter')}</label>
            <div class="checkbox-group">
                <div class="checkbox-item">
                    <input type="checkbox" id="statusPending" checked>
                    <label for="statusPending">${t('status.pending')}</label>
                </div>
                <div class="checkbox-item">
                    <input type="checkbox" id="statusInProgress" checked>
                    <label for="statusInProgress">${t('status.inProgress')}</label>
                </div>
                <div class="checkbox-item">
                    <input type="checkbox" id="statusCompleted" checked>
                    <label for="statusCompleted">${t('status.completed')}</label>
                </div>
            </div>
        </div>

        <div class="form-group">
            <label for="tagFilter">${t('export.tagFilter')}</label>
            <select id="tagFilter" multiple style="height: 100px;">
                <!-- Tags loaded dynamically -->
            </select>
            <small>${t('export.tagFilter.help')}</small>
        </div>

        <div class="form-group">
            <label>${t('export.dateRange')}</label>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 10px;">
                <div>
                    <label for="dateFrom">${t('export.dateFrom')}</label>
                    <input type="date" id="dateFrom">
                </div>
                <div>
                    <label for="dateTo">${t('export.dateTo')}</label>
                    <input type="date" id="dateTo">
                </div>
            </div>
        </div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="previousStep(1)">${t('common.back')}</button>
            <button class="btn btn-primary" onclick="nextStep(3)">${t('common.next')}</button>
        </div>
    </div>

    <!-- Step 3: Output Location -->
    <div class="step" id="step3">
        <div class="step-header">
            <div class="step-indicator">3</div>
            <div class="step-title">${t('export.step3.title')}</div>
        </div>

        <div class="form-group">
            <label for="outputPath">${t('export.outputPath')}</label>
            <div class="file-input-group">
                <input type="text" id="outputPath" readonly>
                <button class="btn btn-secondary" onclick="selectOutputPath()">${t('export.browse')}</button>
            </div>
        </div>

        <div id="validationErrors" class="error" style="display: none;"></div>
        <div id="sizeEstimate" class="info" style="display: none;"></div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="previousStep(2)">${t('common.back')}</button>
            <button class="btn btn-primary" onclick="nextStep(4)">${t('common.next')}</button>
        </div>
    </div>

    <!-- Step 4: Review and Export -->
    <div class="step" id="step4">
        <div class="step-header">
            <div class="step-indicator">4</div>
            <div class="step-title">${t('export.step4.title')}</div>
        </div>

        <div class="export-summary" id="exportSummary">
            <!-- Summary filled dynamically -->
        </div>

        <div id="exportProgress" style="display: none;">
            <div class="progress-bar">
                <div class="progress-fill" id="progressFill"></div>
            </div>
            <div id="progressMessage">${t('export.preparing')}</div>
        </div>

        <div id="exportResult" style="display: none;"></div>

        <div class="button-group">
            <button class="btn btn-secondary" onclick="previousStep(3)" id="backBtn">${t('common.back')}</button>
            <button class="btn btn-primary" onclick="startExport()" id="exportBtn">${t('export.start')}</button>
        </div>
    </div>

    <script nonce="${nonce}">
        const vscode = acquireVsCodeApi();
        let currentStep = 1;
        let projects = [];
        let tags = [];

        // Message handling
        window.addEventListener('message', event => {
            const message = event.data;
            
            switch (message.command) {
                case 'setProjects':
                    projects = message.projects;
                    populateProjectSelect();
                    break;
                    
                case 'setTags':
                    tags = message.tags;
                    populateTagSelect();
                    break;
                    
                case 'setDefaultPath':
                    document.getElementById('outputPath').value = message.path;
                    break;
                    
                case 'setOutputPath':
                    document.getElementById('outputPath').value = message.path;
                    validateCurrentOptions();
                    break;
                    
                case 'validationResult':
                    showValidationResult(message.isValid, message.errors);
                    break;
                    
                case 'sizeEstimate':
                    showSizeEstimate(message.size);
                    break;
                    
                case 'exportStarted':
                    showExportProgress(true);
                    break;
                    
                case 'exportCompleted':
                    showExportResult(true, message.path);
                    break;
                    
                case 'exportFailed':
                    showExportResult(false, null, message.error);
                    break;
                    
                case 'showError':
                    showError(message.message);
                    break;
            }
        });

        function populateProjectSelect() {
            const select = document.getElementById('projectFilter');
            select.innerHTML = '';
            
            projects.forEach(project => {
                const option = document.createElement('option');
                option.value = project.id;
                option.textContent = project.name + (project.isActive ? ' ✅' : '');
                select.appendChild(option);
            });
        }

        function populateTagSelect() {
            const select = document.getElementById('tagFilter');
            select.innerHTML = '';
            
            tags.forEach(tag => {
                const option = document.createElement('option');
                option.value = tag;
                option.textContent = tag;
                select.appendChild(option);
            });
        }

        function nextStep(step) {
            if (step === 3) {
                validateCurrentOptions();
            } else if (step === 4) {
                updateExportSummary();
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

        function validateCurrentOptions() {
            const options = gatherExportOptions();
            vscode.postMessage({
                command: 'validateOptions',
                options: options
            });
            
            vscode.postMessage({
                command: 'estimateSize',
                options: options
            });
        }

        function selectOutputPath() {
            const format = document.getElementById('format').value;
            vscode.postMessage({
                command: 'selectOutputPath',
                format: format
            });
        }

        function gatherExportOptions() {
            const selectedProjects = Array.from(document.getElementById('projectFilter').selectedOptions)
                .map(option => option.value);
            const selectedTags = Array.from(document.getElementById('tagFilter').selectedOptions)
                .map(option => option.value);
            
            const statusFilter = [];
            if (document.getElementById('statusPending').checked) statusFilter.push('beklemede');
            if (document.getElementById('statusInProgress').checked) statusFilter.push('devam_ediyor');
            if (document.getElementById('statusCompleted').checked) statusFilter.push('tamamlandi');

            const options = {
                output_path: document.getElementById('outputPath').value,
                format: document.getElementById('format').value,
                include_completed: document.getElementById('includeCompleted').checked,
                include_dependencies: document.getElementById('includeDependencies').checked,
                include_templates: document.getElementById('includeTemplates').checked,
                include_ai_context: document.getElementById('includeAiContext').checked
            };

            if (selectedProjects.length > 0) {
                options.project_filter = selectedProjects;
            }

            if (selectedTags.length > 0) {
                options.tag_filter = selectedTags;
            }

            if (statusFilter.length > 0 && statusFilter.length < 3) {
                options.status_filter = statusFilter;
            }

            const dateFrom = document.getElementById('dateFrom').value;
            const dateTo = document.getElementById('dateTo').value;
            if (dateFrom || dateTo) {
                options.date_range = {};
                if (dateFrom) options.date_range.from = dateFrom;
                if (dateTo) options.date_range.to = dateTo;
            }

            return options;
        }

        function updateExportSummary() {
            const options = gatherExportOptions();
            const summary = document.getElementById('exportSummary');
            
            let html = '<h3>${t('export.summary')}</h3>';
            
            html += '<div class="summary-item"><span class="summary-label">${t('export.format')}:</span><span>' + options.format.toUpperCase() + '</span></div>';
            html += '<div class="summary-item"><span class="summary-label">${t('export.outputPath')}:</span><span>' + options.output_path + '</span></div>';
            
            const includeOptions = [];
            if (options.include_completed) includeOptions.push('${t('export.includeCompleted')}');
            if (options.include_dependencies) includeOptions.push('${t('export.includeDependencies')}');
            if (options.include_templates) includeOptions.push('${t('export.includeTemplates')}');
            if (options.include_ai_context) includeOptions.push('${t('export.includeAiContext')}');
            
            html += '<div class="summary-item"><span class="summary-label">${t('export.includeOptions')}:</span><span>' + includeOptions.join(', ') + '</span></div>';
            
            if (options.project_filter && options.project_filter.length > 0) {
                const projectNames = options.project_filter.map(id => {
                    const project = projects.find(p => p.id === id);
                    return project ? project.name : id;
                });
                html += '<div class="summary-item"><span class="summary-label">${t('export.projectFilter')}:</span><span>' + projectNames.join(', ') + '</span></div>';
            }
            
            if (options.tag_filter && options.tag_filter.length > 0) {
                html += '<div class="summary-item"><span class="summary-label">${t('export.tagFilter')}:</span><span>' + options.tag_filter.join(', ') + '</span></div>';
            }
            
            if (options.date_range) {
                const dateRange = (options.date_range.from || '') + ' - ' + (options.date_range.to || '');
                html += '<div class="summary-item"><span class="summary-label">${t('export.dateRange')}:</span><span>' + dateRange + '</span></div>';
            }
            
            summary.innerHTML = html;
        }

        function showValidationResult(isValid, errors) {
            const errorDiv = document.getElementById('validationErrors');
            
            if (isValid) {
                errorDiv.style.display = 'none';
            } else {
                errorDiv.innerHTML = '<strong>${t('validation.errors')}:</strong><ul>' + 
                    errors.map(error => '<li>' + error + '</li>').join('') + '</ul>';
                errorDiv.style.display = 'block';
            }
        }

        function showSizeEstimate(size) {
            const estimateDiv = document.getElementById('sizeEstimate');
            estimateDiv.innerHTML = '<strong>${t('export.estimatedSize')}:</strong> ' + size;
            estimateDiv.style.display = 'block';
        }

        function showExportProgress(show) {
            const progressDiv = document.getElementById('exportProgress');
            const backBtn = document.getElementById('backBtn');
            const exportBtn = document.getElementById('exportBtn');
            
            if (show) {
                progressDiv.style.display = 'block';
                backBtn.disabled = true;
                exportBtn.disabled = true;
                exportBtn.textContent = '${t('export.exporting')}';
            } else {
                progressDiv.style.display = 'none';
                backBtn.disabled = false;
                exportBtn.disabled = false;
                exportBtn.textContent = '${t('export.start')}';
            }
        }

        function showExportResult(success, path, error) {
            const resultDiv = document.getElementById('exportResult');
            const exportBtn = document.getElementById('exportBtn');
            
            showExportProgress(false);
            
            if (success) {
                resultDiv.innerHTML = '<div class="info"><strong>${t('export.completed')}!</strong><br>' + 
                    '${t('export.savedTo')}: ' + path + '</div>';
                exportBtn.textContent = '${t('common.close')}';
                exportBtn.onclick = cancel;
            } else {
                resultDiv.innerHTML = '<div class="error"><strong>${t('export.failed')}:</strong> ' + error + '</div>';
                exportBtn.textContent = '${t('export.retry')}';
            }
            
            resultDiv.style.display = 'block';
        }

        function showError(message) {
            const resultDiv = document.getElementById('exportResult');
            resultDiv.innerHTML = '<div class="error">' + message + '</div>';
            resultDiv.style.display = 'block';
        }

        function startExport() {
            const options = gatherExportOptions();
            vscode.postMessage({
                command: 'startExport',
                options: options
            });
        }

        function cancel() {
            vscode.postMessage({ command: 'cancel' });
        }

        // Initial setup
        vscode.postMessage({ command: 'loadProjects' });
        vscode.postMessage({ command: 'loadTags' });
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
