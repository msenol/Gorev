import * as vscode from 'vscode';
import { MCPClient } from './mcp/client';
import { EnhancedGorevTreeProvider } from './providers/enhancedGorevTreeProvider';
import { ProjeTreeProvider } from './providers/projeTreeProvider';
import { TemplateTreeProvider } from './providers/templateTreeProvider';
import { registerCommands } from './commands';
import { StatusBarManager } from './ui/statusBar';
import { FilterToolbar } from './ui/filterToolbar';
import { Logger, LogLevel } from './utils/logger';
import { Config } from './utils/config';
import { COMMANDS } from './utils/constants';

let mcpClient: MCPClient;
let statusBarManager: StatusBarManager;
let filterToolbar: FilterToolbar;
let gorevTreeProvider: EnhancedGorevTreeProvider;
let projeTreeProvider: ProjeTreeProvider;
let templateTreeProvider: TemplateTreeProvider;

let context: vscode.ExtensionContext;

export async function activate(extensionContext: vscode.ExtensionContext) {
  context = extensionContext;
  Logger.info('Gorev extension is starting...');
  
  // Set debug logging
  Logger.setLogLevel(LogLevel.Debug);
  
  // Check if we're in development mode
  const isDevelopment = extensionContext.extensionMode === vscode.ExtensionMode.Development;

  // Initialize configuration
  Config.initialize(context);

  // Create MCP client
  mcpClient = new MCPClient();
  
  // Initialize UI components
  statusBarManager = new StatusBarManager();
  gorevTreeProvider = new EnhancedGorevTreeProvider(mcpClient);
  projeTreeProvider = new ProjeTreeProvider(mcpClient);
  templateTreeProvider = new TemplateTreeProvider(mcpClient);
  
  // Initialize filter toolbar
  filterToolbar = new FilterToolbar(mcpClient, (filter) => {
    gorevTreeProvider.updateFilter(filter);
  });

  // Register tree data providers
  const tasksView = vscode.window.createTreeView('gorevTasks', {
    treeDataProvider: gorevTreeProvider,
    showCollapseAll: true,
    canSelectMany: true,
    dragAndDropController: gorevTreeProvider
  });

  const projectsView = vscode.window.createTreeView('gorevProjects', {
    treeDataProvider: projeTreeProvider,
    showCollapseAll: false,
  });

  const templatesView = vscode.window.createTreeView('gorevTemplates', {
    treeDataProvider: templateTreeProvider,
    showCollapseAll: true,
  });

  context.subscriptions.push(tasksView, projectsView, templatesView);

  // Register commands
  registerCommands(context, mcpClient, {
    gorevTreeProvider,
    projeTreeProvider,
    templateTreeProvider,
    statusBarManager,
    filterToolbar,
  });

  // Register debug commands if in development mode
  if (isDevelopment) {
    const { registerDebugCommands } = await import('./commands/debugCommands');
    registerDebugCommands(context, mcpClient, {
      gorevTreeProvider,
      projeTreeProvider,
      templateTreeProvider,
      statusBarManager,
      filterToolbar,
    });
  }

  // Auto-connect if configured
  if (Config.get('autoConnect')) {
    await vscode.commands.executeCommand(COMMANDS.CONNECT);
    
    // Development modda otomatik test verisi önerisi
    if (isDevelopment) {
      setTimeout(async () => {
        try {
          // Görev sayısını kontrol et
          const result = await mcpClient.callTool('gorev_listele', { tum_projeler: true });
          const hasNoTasks = result.content[0].text.includes('Henüz görev bulunmuyor');
          
          if (hasNoTasks) {
            const answer = await vscode.window.showInformationMessage(
              '🧪 Debug Mode: Görev bulunamadı. Test verileri oluşturmak ister misiniz?',
              'Evet, Oluştur',
              'Hayır'
            );
            
            if (answer === 'Evet, Oluştur') {
              await vscode.commands.executeCommand('gorev.debug.seedTestData');
            }
          }
        } catch (error) {
          // Sessizce devam et
        }
      }, 2000); // 2 saniye bekle
    }
  }

  // Set up refresh interval
  const refreshInterval = Config.get('refreshInterval') as number;
  if (refreshInterval > 0) {
    const intervalId = setInterval(async () => {
      if (mcpClient.isConnected()) {
        try {
          await refreshAllViews();
        } catch (error) {
          Logger.error('Failed to refresh views:', error);
        }
      }
    }, refreshInterval * 1000);
    
    context.subscriptions.push({
      dispose: () => clearInterval(intervalId),
    });
  }

  // Status bar setup
  if (Config.get('showStatusBar')) {
    statusBarManager.show();
    context.subscriptions.push(statusBarManager);
  }
  
  // Show filter toolbar
  filterToolbar.show();
  context.subscriptions.push(filterToolbar);

  // Listen for configuration changes
  context.subscriptions.push(
    vscode.workspace.onDidChangeConfiguration((e) => {
      if (e.affectsConfiguration('gorev')) {
        handleConfigurationChange();
      }
    })
  );

  Logger.info('Gorev extension activated successfully');
}

export function deactivate() {
  Logger.info('Gorev extension is deactivating...');
  
  if (mcpClient) {
    mcpClient.disconnect();
  }
  
  if (statusBarManager) {
    statusBarManager.dispose();
  }

  Logger.info('Gorev extension deactivated');
}


async function refreshAllViews(): Promise<void> {
  if (!mcpClient.isConnected()) {
    Logger.warn('Cannot refresh views: not connected to MCP server');
    return;
  }
  
  try {
    await Promise.all([
      gorevTreeProvider.refresh(),
      projeTreeProvider.refresh(),
      templateTreeProvider.refresh(),
    ]);
    
    statusBarManager.update();
  } catch (error) {
    Logger.error('Error refreshing views:', error);
    throw error;
  }
}

function handleConfigurationChange(): void {
  // Handle configuration changes
  const showStatusBar = Config.get('showStatusBar') as boolean;
  if (showStatusBar && !statusBarManager.isVisible()) {
    statusBarManager.show();
  } else if (!showStatusBar && statusBarManager.isVisible()) {
    statusBarManager.hide();
  }
}