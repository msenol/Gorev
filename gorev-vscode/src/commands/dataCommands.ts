import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import * as path from 'path';
import * as os from 'os';
import { ApiClient, ExportRequest } from '../api/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { Logger } from '../utils/logger';
import { ExportDialog } from '../ui/exportDialog';
import { ImportWizard } from '../ui/importWizard';

/** Extended export options for validation (allows optional fields) */
type ExportOptions = Partial<ExportRequest> & {
  date_range?: {
    from?: string;
    to?: string;
  };
};

/**
 * Export/Import data commands for the Gorev VS Code extension
 */
export function registerDataCommands(
  context: vscode.ExtensionContext,
  apiClient: ApiClient,
  providers: CommandContext
): void {
  Logger.info('Registering data export/import commands');

  // Export Data Command - Opens comprehensive export dialog
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.EXPORT_DATA, async () => {
      try {
        if (!apiClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        const exportDialog = new ExportDialog(context, apiClient);
        await exportDialog.show();
      } catch (error) {
        Logger.error('Export data command failed', error);
        vscode.window.showErrorMessage(
          t('error.exportFailed', { error: String(error) })
        );
      }
    })
  );

  // Import Data Command - Opens import wizard
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.IMPORT_DATA, async () => {
      try {
        if (!apiClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        const importWizard = new ImportWizard(context, apiClient, providers);
        await importWizard.show();
      } catch (error) {
        Logger.error('Import data command failed', error);
        vscode.window.showErrorMessage(
          t('error.importFailed', { error: String(error) })
        );
      }
    })
  );

  // Export Current View Command - Exports filtered tasks from current view
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.EXPORT_CURRENT_VIEW, async (element?: vscode.TreeItem) => {
      try {
        if (!apiClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        await exportCurrentView(apiClient, providers, element);
      } catch (error) {
        Logger.error('Export current view command failed', error);
        vscode.window.showErrorMessage(
          t('error.exportCurrentViewFailed', { error: String(error) })
        );
      }
    })
  );

  // Quick Export Command - One-click JSON export with default settings
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.QUICK_EXPORT, async () => {
      try {
        if (!apiClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        await quickExport(apiClient);
      } catch (error) {
        Logger.error('Quick export command failed', error);
        vscode.window.showErrorMessage(
          t('error.quickExportFailed', { error: String(error) })
        );
      }
    })
  );
}

/**
 * Export current filtered view with context-appropriate settings
 */
async function exportCurrentView(
  apiClient: ApiClient,
  providers: CommandContext,
  element?: vscode.TreeItem
): Promise<void> {
  Logger.info('Exporting current view', { element });

  // Show file save dialog
  const defaultFileName = `gorev-export-${new Date().toISOString().split('T')[0]}.json`;
  const saveUri = await vscode.window.showSaveDialog({
    defaultUri: vscode.Uri.file(path.join(vscode.workspace.rootPath || '', defaultFileName)),
    filters: {
      'JSON Files': ['json'],
      'CSV Files': ['csv'],
      'All Files': ['*']
    },
    title: t('export.selectLocation')
  });

  if (!saveUri) {
    return; // User cancelled
  }

  // Determine format from file extension
  const format = path.extname(saveUri.fsPath).toLowerCase() === '.csv' ? 'csv' : 'json';

  // Show progress
  await vscode.window.withProgress({
    location: vscode.ProgressLocation.Notification,
    title: t('export.inProgress'),
    cancellable: false
  }, async (progress) => {
    progress.report({ increment: 10, message: t('export.preparing') });

    // Build export options based on current view context
    const exportOptions: ExportOptions = {
      output_path: saveUri.fsPath,
      format: format,
      include_completed: true,
      include_dependencies: true,
      include_templates: false,
      include_ai_context: false
    };

    // If exporting from a specific task/project context, add filters
    if (element && element.contextValue && element.id) {
      if (element.contextValue === 'project' || element.contextValue === 'project-active') {
        exportOptions.project_filter = [element.id];
      }
    }

    progress.report({ increment: 30, message: t('export.exporting') });

    // Call REST API export endpoint
    const result = await apiClient.exportData(exportOptions as ExportRequest);

    progress.report({ increment: 80, message: t('export.completing') });

    if (!result.success) {
      throw new Error(result.message || 'Export failed');
    }

    progress.report({ increment: 100, message: t('export.complete') });

    // Show success message with option to open file
    const openAction = t('export.openFile');
    const action = await vscode.window.showInformationMessage(
      t('export.success', { path: saveUri.fsPath }),
      openAction
    );

    if (action === openAction) {
      await vscode.commands.executeCommand('vscode.open', saveUri);
    }
  });
}

/**
 * Quick export with default settings to Downloads folder
 */
async function quickExport(apiClient: ApiClient): Promise<void> {
  Logger.info('Performing quick export');

  // Generate default filename with timestamp
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-').split('T')[0];
  const defaultFileName = `gorev-quick-export-${timestamp}.json`;
  
  // Use Downloads folder or workspace root
  const downloadsPath = path.join(os.homedir(), 'Downloads');
  const defaultPath = path.join(downloadsPath, defaultFileName);

  // Show progress
  await vscode.window.withProgress({
    location: vscode.ProgressLocation.Notification,
    title: t('export.quickExporting'),
    cancellable: false
  }, async (progress) => {
    progress.report({ increment: 20, message: t('export.preparing') });

    const exportOptions = {
      output_path: defaultPath,
      format: 'json',
      include_completed: true,
      include_dependencies: true,
      include_templates: false,
      include_ai_context: false
    };

    progress.report({ increment: 50, message: t('export.exporting') });

    // Call REST API export endpoint
    const result = await apiClient.exportData(exportOptions as ExportRequest);

    if (!result.success) {
      throw new Error(result.message || 'Quick export failed');
    }

    progress.report({ increment: 100, message: t('export.complete') });

    // Show success notification
    const openAction = t('export.openFile');
    const openFolderAction = t('export.openFolder');
    
    const action = await vscode.window.showInformationMessage(
      t('export.quickSuccess', { filename: defaultFileName }),
      openAction,
      openFolderAction
    );

    if (action === openAction) {
      await vscode.commands.executeCommand('vscode.open', vscode.Uri.file(defaultPath));
    } else if (action === openFolderAction) {
      await vscode.commands.executeCommand('revealFileInOS', vscode.Uri.file(defaultPath));
    }
  });
}

/**
 * Utility function to validate export options
 */
export function validateExportOptions(options: ExportOptions): { isValid: boolean; errors: string[] } {
  const errors: string[] = [];

  if (!options.output_path) {
    errors.push(t('validation.outputPathRequired'));
  }

  if (options.format && !['json', 'csv'].includes(options.format)) {
    errors.push(t('validation.invalidFormat'));
  }

  if (options.date_range) {
    if (options.date_range.from && options.date_range.to) {
      const fromDate = new Date(options.date_range.from);
      const toDate = new Date(options.date_range.to);
      
      if (fromDate > toDate) {
        errors.push(t('validation.invalidDateRange'));
      }
    }
  }

  return {
    isValid: errors.length === 0,
    errors
  };
}

/**
 * Utility function to estimate export size
 */
export async function estimateExportSize(apiClient: ApiClient, _options: ExportOptions): Promise<string> {
  try {
    // Get summary to estimate size using REST API
    const summaryResult = await apiClient.getSummary();
    if (!summaryResult.success || !summaryResult.data) {
      return t('export.sizeUnknown');
    }

    // Extract task and project counts from summary data
    const taskCount = summaryResult.data.tasks || 0;
    const projectCount = summaryResult.data.projects || 0;

    // Rough estimation: ~500 bytes per task, ~200 bytes per project
    const estimatedBytes = (taskCount * 500) + (projectCount * 200);
    
    if (estimatedBytes < 1024) {
      return `${estimatedBytes} bytes`;
    } else if (estimatedBytes < 1024 * 1024) {
      return `${Math.round(estimatedBytes / 1024)} KB`;
    } else {
      return `${Math.round(estimatedBytes / (1024 * 1024))} MB`;
    }
  } catch (error) {
    Logger.warn('Failed to estimate export size', error);
    return t('export.sizeUnknown');
  }
}
