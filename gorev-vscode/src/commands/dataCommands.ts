import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import * as path from 'path';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { Logger } from '../utils/logger';
import { ExportDialog } from '../ui/exportDialog';
import { ImportWizard } from '../ui/importWizard';

/**
 * Export/Import data commands for the Gorev VS Code extension
 */
export function registerDataCommands(
  context: vscode.ExtensionContext,
  mcpClient: MCPClient,
  providers: CommandContext
): void {
  Logger.info('Registering data export/import commands');

  // Export Data Command - Opens comprehensive export dialog
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.EXPORT_DATA, async () => {
      try {
        if (!mcpClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        const exportDialog = new ExportDialog(context, mcpClient);
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
        if (!mcpClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        const importWizard = new ImportWizard(context, mcpClient, providers);
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
    vscode.commands.registerCommand(COMMANDS.EXPORT_CURRENT_VIEW, async (element?: any) => {
      try {
        if (!mcpClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        await exportCurrentView(mcpClient, providers, element);
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
        if (!mcpClient.isConnected()) {
          vscode.window.showWarningMessage(t('connection.notConnected'));
          return;
        }

        await quickExport(mcpClient);
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
  mcpClient: MCPClient,
  providers: CommandContext,
  element?: any
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
    const exportOptions: any = {
      output_path: saveUri.fsPath,
      format: format,
      include_completed: true,
      include_dependencies: true,
      include_templates: false,
      include_ai_context: false
    };

    // If exporting from a specific task/project context, add filters
    if (element && element.contextValue) {
      if (element.contextValue === 'project' || element.contextValue === 'project-active') {
        exportOptions.project_filter = [element.id];
      }
    }

    progress.report({ increment: 30, message: t('export.exporting') });

    // Call MCP export tool
    const result = await mcpClient.callTool('gorev_export', exportOptions);

    progress.report({ increment: 80, message: t('export.completing') });

    if (result.isError) {
      throw new Error(result.content[0]?.text || 'Export failed');
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
async function quickExport(mcpClient: MCPClient): Promise<void> {
  Logger.info('Performing quick export');

  // Generate default filename with timestamp
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-').split('T')[0];
  const defaultFileName = `gorev-quick-export-${timestamp}.json`;
  
  // Use Downloads folder or workspace root
  const downloadsPath = path.join(require('os').homedir(), 'Downloads');
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

    // Call MCP export tool
    const result = await mcpClient.callTool('gorev_export', exportOptions);

    if (result.isError) {
      throw new Error(result.content[0]?.text || 'Quick export failed');
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
export function validateExportOptions(options: any): { isValid: boolean; errors: string[] } {
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
export async function estimateExportSize(mcpClient: MCPClient, options: any): Promise<string> {
  try {
    // Get summary to estimate size
    const summaryResult = await mcpClient.callTool('ozet_goster');
    if (summaryResult.isError) {
      return t('export.sizeUnknown');
    }

    // Parse summary for task/project counts
    const summaryText = summaryResult.content[0]?.text || '';
    const taskMatch = summaryText.match(/(\d+).*görev/i);
    const projectMatch = summaryText.match(/(\d+).*proje/i);

    const taskCount = taskMatch ? parseInt(taskMatch[1]) : 0;
    const projectCount = projectMatch ? parseInt(projectMatch[1]) : 0;

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
