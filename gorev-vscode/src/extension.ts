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
  Logger.info(vscode.l10n.t('extension.starting'));
  
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
    // If the filter object is empty, clear all filters
    if (Object.keys(filter).length === 0) {
      gorevTreeProvider.clearFilters();
    } else {
      gorevTreeProvider.updateFilter(filter);
    }
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
              vscode.l10n.t('debug.noTasksFound'),
              vscode.l10n.t('debug.yesCreate'),
              vscode.l10n.t('debug.no')
            );
            
            if (answer === vscode.l10n.t('debug.yesCreate')) {
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
          Logger.error(vscode.l10n.t('log.failedRefreshViews'), error);
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

  Logger.info(vscode.l10n.t('extension.activated'));
}

export function deactivate() {
  Logger.info(vscode.l10n.t('extension.deactivated'));
  
  if (mcpClient) {
    mcpClient.disconnect();
  }
  
  if (statusBarManager) {
    statusBarManager.dispose();
  }

  Logger.info(vscode.l10n.t('extension.deactivated'));
}


async function refreshAllViews(): Promise<void> {
  if (!mcpClient.isConnected()) {
    Logger.warn(vscode.l10n.t('log.cannotRefreshViews'));
    return;
  }
  
  try {
    // Refresh sequentially to avoid overwhelming the MCP server
    await gorevTreeProvider.refresh();
    await projeTreeProvider.refresh();
    await templateTreeProvider.refresh();
    
    statusBarManager.update();
  } catch (error) {
    Logger.error(vscode.l10n.t('log.errorRefreshingViews'), error);
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