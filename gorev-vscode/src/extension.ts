import * as vscode from 'vscode';
import { MCPClient } from './mcp/client';
import { UnifiedClient } from './unified/client';
import { EnhancedGorevTreeProvider } from './providers/enhancedGorevTreeProvider';
import { ProjeTreeProvider } from './providers/projeTreeProvider';
import { TemplateTreeProvider } from './providers/templateTreeProvider';
import { registerCommands } from './commands';
import { StatusBarManager } from './ui/statusBar';
import { FilterToolbar } from './ui/filterToolbar';
import { Logger, LogLevel } from './utils/logger';
import { Config } from './utils/config';
import { COMMANDS } from './utils/constants';
import { initializeL10n } from './utils/l10n';
import { RefreshManager, RefreshTarget, RefreshReason, RefreshPriority } from './managers/refreshManager';
import { measureAsync } from './utils/performance';
import { debounceConfig } from './utils/debounce';

let mcpClient: MCPClient;
let unifiedClient: UnifiedClient;
let statusBarManager: StatusBarManager;
let filterToolbar: FilterToolbar;
let gorevTreeProvider: EnhancedGorevTreeProvider;
let projeTreeProvider: ProjeTreeProvider;
let templateTreeProvider: TemplateTreeProvider;
let refreshManager: RefreshManager;

let context: vscode.ExtensionContext;
let debouncedConfigHandler: ReturnType<typeof debounceConfig>;

export async function activate(extensionContext: vscode.ExtensionContext) {
  context = extensionContext;

  // Initialize L10n system first
  await initializeL10n(context);

  Logger.info('Extension starting...');

  // Check if we're in development mode
  const isDevelopment = extensionContext.extensionMode === vscode.ExtensionMode.Development;

  // Initialize configuration
  Config.initialize(context);

  // Create clients
  mcpClient = new MCPClient();

  // Get client mode from configuration (default to 'auto' for automatic detection)
  const clientMode = Config.get<string>('clientMode') || 'auto';
  unifiedClient = new UnifiedClient({ mode: clientMode as 'api' | 'mcp' | 'auto' });

  // Initialize RefreshManager first
  refreshManager = RefreshManager.getInstance();

  // Initialize UI components - use unified client for new components
  statusBarManager = new StatusBarManager();
  gorevTreeProvider = new EnhancedGorevTreeProvider(unifiedClient);
  projeTreeProvider = new ProjeTreeProvider(unifiedClient);
  templateTreeProvider = new TemplateTreeProvider(unifiedClient);

  // Note: Other providers will be integrated with RefreshManager in future iterations
  // For now, only EnhancedGorevTreeProvider implements the RefreshProvider interface

  // Initialize filter toolbar
  filterToolbar = new FilterToolbar(unifiedClient, (filter) => {
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

  // Listen for unified client events
  unifiedClient.on('connected', (data: any) => {
    Logger.info(`[Extension] Connected via ${data.mode} mode`);
    statusBarManager.setConnectionStatus(true, data.mode);
  });

  unifiedClient.on('disconnected', (data: any) => {
    Logger.info(`[Extension] Disconnected from ${data.mode} mode`);
    statusBarManager.setConnectionStatus(false, data.mode);
  });

  unifiedClient.on('error', (data: any) => {
    Logger.error(`[Extension] ${data.mode} error:`, data.error);
  });

  // Listen for database mode changes (MCP only)
  mcpClient.on('databaseModeChanged', (data: any) => {
    statusBarManager.setDatabaseMode(data.mode, data.path);
  });

  // Register commands with unified client
  registerCommands(context, unifiedClient, {
    gorevTreeProvider,
    projeTreeProvider,
    templateTreeProvider,
    statusBarManager,
    filterToolbar,
  });

  // Register debug commands if in development mode
  if (isDevelopment) {
    const { registerDebugCommands } = await import('./commands/debugCommands');
    registerDebugCommands(context, unifiedClient, {
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
          const result = await unifiedClient.callTool('gorev_listele', { tum_projeler: true });
          const hasNoTasks = result.content[0].text.includes('Henüz görev bulunmuyor');

          if (hasNoTasks) {
            const { t } = await import('./utils/l10n');
            const answer = await vscode.window.showInformationMessage(
              t('debug.noTasksFound'),
              t('debug.yesCreate'),
              t('debug.no')
            );

            if (answer === t('debug.yesCreate')) {
              await vscode.commands.executeCommand('gorev.debug.seedTestData');
            }
          }
        } catch (error) {
          Logger.debug('[Extension] Failed to check test data status:', error);
        }
      }, 2000); // 2 saniye bekle
    }
  }

  // Set up refresh interval with RefreshManager
  const refreshInterval = Config.get('refreshInterval') as number;
  if (refreshInterval > 0) {
    Logger.info(`[Extension] Setting up refresh interval: ${refreshInterval} seconds`);

    const intervalId = setInterval(async () => {
      if (unifiedClient.isConnected()) {
        try {
          // Use RefreshManager for coordinated refresh
          await refreshManager.requestRefresh(
            RefreshReason.INTERVAL,
            [RefreshTarget.ALL],
            RefreshPriority.LOW
          );
        } catch (error) {
          const { t } = await import('./utils/l10n');
          Logger.error(t('log.failedRefreshViews'), error);
        }
      }
    }, refreshInterval * 1000);

    context.subscriptions.push({
      dispose: () => clearInterval(intervalId),
    });
  } else {
    Logger.info('[Extension] Auto-refresh disabled (refreshInterval = 0)');
  }

  // Status bar setup
  if (Config.get('showStatusBar')) {
    statusBarManager.show();
    context.subscriptions.push(statusBarManager);
  }
  
  // Show filter toolbar
  filterToolbar.show();
  context.subscriptions.push(filterToolbar);

  // Initialize debounced configuration handler
  debouncedConfigHandler = debounceConfig(handleConfigurationChange);

  // Listen for configuration changes (consolidated single handler)
  context.subscriptions.push(
    vscode.workspace.onDidChangeConfiguration((e) => {
      if (e.affectsConfiguration('gorev')) {
        debouncedConfigHandler();
      }
    })
  );

  const { t } = await import('./utils/l10n');
  Logger.info(t('extension.activated'));
}

export async function deactivate() {
  const { t } = await import('./utils/l10n');
  Logger.info(t('extension.deactivated'));

  // Clean up debounced handlers
  if (debouncedConfigHandler) {
    debouncedConfigHandler.cancel();
  }

  // Dispose providers
  if (gorevTreeProvider) {
    gorevTreeProvider.dispose();
  }

  // Dispose RefreshManager
  if (refreshManager) {
    refreshManager.dispose();
  }

  if (unifiedClient) {
    unifiedClient.disconnect();
  }

  if (mcpClient) {
    mcpClient.disconnect();
  }

  if (statusBarManager) {
    statusBarManager.dispose();
  }
}


/**
 * DEPRECATED: Use RefreshManager.requestRefresh() instead
 * Legacy function kept for backward compatibility during transition
 */
async function refreshAllViews(): Promise<void> {
  Logger.warn('[Extension] refreshAllViews() is deprecated, use RefreshManager instead');

  if (!refreshManager) {
    Logger.error('[Extension] RefreshManager not initialized');
    return;
  }

  await refreshManager.requestRefresh(
    RefreshReason.MANUAL,
    [RefreshTarget.ALL],
    RefreshPriority.HIGH
  );
}

/**
 * Handle configuration changes with RefreshManager integration
 */
async function handleConfigurationChange(): Promise<void> {
  Logger.debug('[Extension] Configuration changed');

  await measureAsync(
    'config-change',
    async () => {
      // Handle status bar visibility
      const showStatusBar = Config.get('showStatusBar') as boolean;
      if (showStatusBar && !statusBarManager.isVisible()) {
        statusBarManager.show();
      } else if (!showStatusBar && statusBarManager.isVisible()) {
        statusBarManager.hide();
      }

      // Request refresh through RefreshManager
      if (refreshManager) {
        await refreshManager.requestRefresh(
          RefreshReason.CONFIG_CHANGE,
          [RefreshTarget.ALL],
          RefreshPriority.NORMAL
        );
      }
    },
    'configuration-change'
  );
}