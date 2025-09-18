import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
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
  mcpClient: MCPClient,
  providers: CommandContext
): void {
  // Register all command groups
  registerGorevCommands(context, mcpClient, providers);
  registerProjeCommands(context, mcpClient, providers);
  registerTemplateCommands(context, mcpClient, providers);
  registerEnhancedGorevCommands(context, mcpClient, providers);
  registerInlineEditCommands(context, mcpClient, providers);
  registerDataCommands(context, mcpClient, providers);
  registerDatabaseCommands(context, mcpClient, providers);
  
  if (providers.filterToolbar) {
    registerFilterCommands(context, mcpClient, providers);
  }

  // Register general commands
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SHOW_SUMMARY, async () => {
      try {
        if (!mcpClient.isConnected()) {
          vscode.window.showWarningMessage('Not connected to Gorev server');
          return;
        }
        const result = await mcpClient.callTool('ozet_goster');
        const summaryPanel = vscode.window.createWebviewPanel(
          'gorevSummary',
          'Gorev Summary',
          vscode.ViewColumn.One,
          {}
        );
        
        summaryPanel.webview.html = getSummaryHtml(result.content[0].text);
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to show summary: ${error}`);
      }
    })
  );

  // Connect command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CONNECT, async () => {
      try {
        await connectToServer(mcpClient, providers);
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to connect: ${error}`);
      }
    })
  );

  // Disconnect command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DISCONNECT, () => {
      mcpClient.disconnect();
      providers.statusBarManager.setDisconnected();
      vscode.window.showInformationMessage('Disconnected from Gorev server');
    })
  );
}

async function connectToServer(mcpClient: MCPClient, providers: CommandContext): Promise<void> {
  const config = vscode.workspace.getConfiguration('gorev');
  const serverMode = config.get<string>('serverMode', 'npx');
  let serverPath = config.get<string>('serverPath');

  // Validate configuration based on server mode
  if (serverMode === 'binary') {
    if (!serverPath) {
      throw new Error('Gorev server path not configured for binary mode. Please set gorev.serverPath in settings.');
    }

    // Convert WSL path to Windows path if needed
    if (process.platform === 'win32' && serverPath.startsWith('/mnt/')) {
      // Convert /mnt/f/... to F:\...
      const drive = serverPath.charAt(5).toUpperCase();
      serverPath = drive + ':\\' + serverPath.substring(7).replace(/\//g, '\\');
      Logger.debug(`Converted WSL path to Windows path: ${serverPath}`);
    }
  } else {
    // NPX mode - no server path validation needed
    Logger.info('Using NPX mode - no server path required');
  }

  providers.statusBarManager.setConnecting();
  
  try {
    await mcpClient.connect(serverPath);
    providers.statusBarManager.setConnected();
    
    // Refresh all views after connection - sequentially to avoid overwhelming the MCP server
    await providers.gorevTreeProvider.refresh();
    await providers.projeTreeProvider.refresh();
    await providers.templateTreeProvider.refresh();
    
    vscode.window.showInformationMessage('Connected to Gorev server');
  } catch (error) {
    providers.statusBarManager.setDisconnected();
    throw error;
  }
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
