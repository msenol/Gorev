import * as vscode from 'vscode';
import { ApiClient, ApiError } from '../api/client';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';
import { ProjeTreeProvider } from '../providers/projeTreeProvider';
import { TemplateTreeProvider } from '../providers/templateTreeProvider';
import { StatusBarManager } from '../ui/statusBar';
import { FilterToolbar } from '../ui/filterToolbar';
import { COMMANDS } from '../utils/constants';
import { Logger } from '../utils/logger';
import { registerGorevCommands } from './gorevCommands';
import { registerProjeCommands } from './projeCommands';
import { registerTemplateCommands } from './templateCommands';
import { registerEnhancedGorevCommands } from './enhancedGorevCommands';
import { registerInlineEditCommands } from './inlineEditCommands';
import { registerFilterCommands } from './filterCommands';
import { registerDataCommands } from './dataCommands';
import { registerDatabaseCommands } from './databaseCommands';

export interface CommandContext {
  gorevTreeProvider: EnhancedGorevTreeProvider;
  projeTreeProvider: ProjeTreeProvider;
  templateTreeProvider: TemplateTreeProvider;
  statusBarManager: StatusBarManager;
  filterToolbar?: FilterToolbar;
}

export function registerCommands(
  context: vscode.ExtensionContext,
  apiClient: ApiClient,
  providers: CommandContext
): void {
  // Register all command groups
  registerGorevCommands(context, apiClient, providers);
  registerProjeCommands(context, apiClient, providers);
  registerTemplateCommands(context, apiClient, providers);
  registerEnhancedGorevCommands(context, apiClient, providers);
  registerInlineEditCommands(context, apiClient, providers);
  registerDataCommands(context, apiClient, providers);
  registerDatabaseCommands(context, apiClient, providers);
  
  if (providers.filterToolbar) {
    registerFilterCommands(context, apiClient, providers);
  }

  // Initialize API client for general commands

  // Register general commands
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SHOW_SUMMARY, async () => {
      try {
        if (!apiClient.isConnected()) {
          vscode.window.showWarningMessage('Not connected to Gorev server');
          return;
        }

        // Use REST API to get summary
        const response = await apiClient.getSummary();

        if (!response.success || !response.data) {
          vscode.window.showErrorMessage('Failed to get summary from server');
          return;
        }

        const summaryPanel = vscode.window.createWebviewPanel(
          'gorevSummary',
          'Gorev Summary',
          vscode.ViewColumn.One,
          {}
        );

        // Format summary data as HTML
        summaryPanel.webview.html = getSummaryHtml(formatSummaryData(response.data));
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[ShowSummary] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(`Failed to show summary: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(`Failed to show summary: ${error}`);
        }
      }
    })
  );

  // Connect command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CONNECT, async () => {
      try {
        await connectToServer(apiClient, providers);
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to connect: ${error}`);
      }
    })
  );

  // Disconnect command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DISCONNECT, () => {
      apiClient.disconnect();
      providers.statusBarManager.setDisconnected();
      vscode.window.showInformationMessage('Disconnected from Gorev server');
    })
  );
}

async function connectToServer(client: ApiClient, providers: CommandContext): Promise<void> {
  providers.statusBarManager.setConnecting();

  try {
    // Use the unified client's connect method which handles auto-detection
    await client.connect();
    providers.statusBarManager.setConnected();

    // Refresh all views after connection - sequentially to avoid overwhelming the server
    await providers.gorevTreeProvider.refresh();
    await providers.projeTreeProvider.refresh();
    await providers.templateTreeProvider.refresh();

    vscode.window.showInformationMessage('Connected to Gorev server');
  } catch (error) {
    providers.statusBarManager.setDisconnected();
    throw error;
  }
}

/**
 * Format summary data from API response to markdown-like text
 */
function formatSummaryData(data: any): string {
  let text = `# Özet Rapor\n\n`;

  if (data.toplam_proje !== undefined) {
    text += `**Toplam Proje:** ${data.toplam_proje}\n`;
  }
  if (data.toplam_gorev !== undefined) {
    text += `**Toplam Görev:** ${data.toplam_gorev}\n\n`;
  }

  if (data.durum_dagilimi) {
    text += `### Durum Dağılımı\n`;
    text += `- Beklemede: ${data.durum_dagilimi.beklemede || 0}\n`;
    text += `- Devam Ediyor: ${data.durum_dagilimi.devam_ediyor || 0}\n`;
    text += `- Tamamlandı: ${data.durum_dagilimi.tamamlandi || 0}\n\n`;
  }

  if (data.oncelik_dagilimi) {
    text += `### Öncelik Dağılımı\n`;
    text += `- Yüksek: ${data.oncelik_dagilimi.yuksek || 0}\n`;
    text += `- Orta: ${data.oncelik_dagilimi.orta || 0}\n`;
    text += `- Düşük: ${data.oncelik_dagilimi.dusuk || 0}\n`;
  }

  return text;
}

function getSummaryHtml(content: string): string {
  return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: var(--vscode-font-family);
            color: var(--vscode-foreground);
            background-color: var(--vscode-editor-background);
            padding: 20px;
            line-height: 1.6;
        }
        h1, h2, h3 {
            color: var(--vscode-foreground);
        }
        pre {
            background-color: var(--vscode-textBlockQuote-background);
            padding: 10px;
            border-radius: 4px;
            overflow-x: auto;
        }
    </style>
</head>
<body>
    ${content.replace(/\n/g, '<br>')}
</body>
</html>`;
}
