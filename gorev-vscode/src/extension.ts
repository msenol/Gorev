import * as vscode from 'vscode';
import { ApiClient } from './api/client';
import { UnifiedServerManager } from './managers/unifiedServerManager';
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

let serverManager: UnifiedServerManager;
let apiClient: ApiClient;
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

  // Create UnifiedServerManager
  const apiHost = Config.get<string>('apiHost') || 'localhost';
  const apiPort = Config.get<number>('apiPort') || 5082;
  serverManager = new UnifiedServerManager(apiHost, apiPort);

  // Initialize RefreshManager first
  refreshManager = RefreshManager.getInstance();

  // Initialize UI components
  statusBarManager = new StatusBarManager();

  // Initialize server connection and register workspace
  try {
    await serverManager.initialize();
    apiClient = serverManager.getApiClient();
    Logger.info('Successfully initialized server connection and registered workspace');

    // Update status bar with workspace context
    const workspaceContext = serverManager.getWorkspaceContext();
    if (workspaceContext) {
      statusBarManager.setWorkspaceContext(workspaceContext);
      vscode.window.showInformationMessage(
        `Gorev: Workspace "${workspaceContext.workspaceName}" registered successfully`
      );
    }
  } catch (error) {
    Logger.warn('Failed to initialize server connection:', error);
    apiClient = serverManager.getApiClient(); // Get client anyway for offline mode
    statusBarManager.setWorkspaceContext(undefined); // No workspace registered
    vscode.window.showWarningMessage(
      'Gorev: Could not connect to API server. Please make sure the server is running.'
    );
  }
  gorevTreeProvider = new EnhancedGorevTreeProvider(apiClient);
  projeTreeProvider = new ProjeTreeProvider(apiClient);
  templateTreeProvider = new TemplateTreeProvider(apiClient);

  // Note: Other providers will be integrated with RefreshManager in future iterations
  // For now, only EnhancedGorevTreeProvider implements the RefreshProvider interface

  // Initialize filter toolbar
  filterToolbar = new FilterToolbar(apiClient, (filter) => {
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

  // Listen for API client events
  apiClient.on('connected', () => {
    Logger.info('[Extension] Connected to API server');
    statusBarManager.setConnectionStatus(true, 'api');
  });

  apiClient.on('disconnected', () => {
    Logger.info('[Extension] Disconnected from API server');
    statusBarManager.setConnectionStatus(false, 'api');
  });

  apiClient.on('error', (error: any) => {
    Logger.error('[Extension] API error:', error);
  });

  // Register commands with API client
  registerCommands(context, apiClient, {
    gorevTreeProvider,
    projeTreeProvider,
    templateTreeProvider,
    statusBarManager,
    filterToolbar,
  });

  // Register debug commands if in development mode
  if (isDevelopment) {
    const { registerDebugCommands } = await import('./commands/debugCommands');
    registerDebugCommands(context, apiClient, {
      gorevTreeProvider,
      projeTreeProvider,
      templateTreeProvider,
      statusBarManager,
      filterToolbar,
    });
  }

  // Try to connect to API server
  try {
    await apiClient.checkHealth();
    Logger.info('[Extension] API server is available');

    // Development modda otomatik test verisi önerisi
    if (isDevelopment) {
      setTimeout(async () => {
        try {
          // Görev sayısını kontrol et
          const response = await apiClient.getTasks({ tum_projeler: true });
          const hasNoTasks = response.total === 0;

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
  } catch (error) {
    Logger.error('[Extension] Failed to connect to API server:', error);
    vscode.window.showWarningMessage('Gorev API server is not running. Please start the server with: ./gorev serve --api-port 5082');
  }

  // Set up refresh interval with RefreshManager
  const refreshInterval = Config.get('refreshInterval') as number;
  if (refreshInterval > 0) {
    Logger.info(`[Extension] Setting up refresh interval: ${refreshInterval} seconds`);

    const intervalId = setInterval(async () => {
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

  // Dispose UnifiedServerManager (this will also disconnect apiClient)
  if (serverManager) {
    serverManager.dispose();
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