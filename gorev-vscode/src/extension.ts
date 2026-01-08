import * as vscode from 'vscode';
import { ApiClient } from './api/client';
import { UnifiedServerManager } from './managers/unifiedServerManager';
import { EnhancedGorevTreeProvider } from './providers/enhancedGorevTreeProvider';
import { ProjeTreeProvider } from './providers/projeTreeProvider';
import { TemplateTreeProvider } from './providers/templateTreeProvider';
import { registerCommands } from './commands';
import { StatusBarManager } from './ui/statusBar';
import { FilterToolbar } from './ui/filterToolbar';
import { Logger } from './utils/logger';
import { Config } from './utils/config';
import { initializeL10n } from './utils/l10n';
import { RefreshManager, RefreshTarget, RefreshReason, RefreshPriority } from './managers/refreshManager';
import { measureAsync } from './utils/performance';
import { debounceConfig } from './utils/debounce';
import { DatabaseFileWatcher } from './managers/databaseFileWatcher';
import { WebSocketClient } from './managers/websocketClient';
import { AIPanelProvider } from './panels/aiPanel';
import { AIStatusBar } from './statusbar/aiStatusBar';

let serverManager: UnifiedServerManager;
let apiClient: ApiClient;
let statusBarManager: StatusBarManager;
let filterToolbar: FilterToolbar;
let gorevTreeProvider: EnhancedGorevTreeProvider;
let projeTreeProvider: ProjeTreeProvider;
let templateTreeProvider: TemplateTreeProvider;
let refreshManager: RefreshManager;
let databaseWatcher: DatabaseFileWatcher;
let webSocketClient: WebSocketClient | null = null;
let aiPanelProvider: AIPanelProvider;
let aiStatusBar: AIStatusBar;

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

  // Create UnifiedServerManager with serverUrl or apiHost/apiPort
  const serverUrl = Config.get<string>('serverUrl');
  let apiHost: string;
  let apiPort: number;

  if (serverUrl) {
    // Parse serverUrl to extract host and port
    try {
      const url = new URL(serverUrl);
      apiHost = url.hostname;
      apiPort = parseInt(url.port, 10) || 5082;
      Logger.info(`[Extension] Using serverUrl: ${serverUrl}`);
    } catch {
      Logger.warn(`[Extension] Invalid serverUrl: ${serverUrl}, falling back to apiHost/apiPort`);
      apiHost = Config.get<string>('apiHost') || 'localhost';
      apiPort = Config.get<number>('apiPort') || 5082;
    }
  } else {
    apiHost = Config.get<string>('apiHost') || 'localhost';
    apiPort = Config.get<number>('apiPort') || 5082;
  }

  // Create UnifiedServerManager - workspaceId is auto-managed via workspace settings
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

  // Register all providers with RefreshManager for coordinated refresh
  refreshManager.registerProvider([RefreshTarget.PROJECTS], projeTreeProvider);
  refreshManager.registerProvider([RefreshTarget.TEMPLATES], templateTreeProvider);
  Logger.info('[Extension] All tree providers registered with RefreshManager');

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

  apiClient.on('error', (error: unknown) => {
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

  // Initialize AI components (Rule 15 compliance - gracefully degrades if AI not configured)
  aiPanelProvider = AIPanelProvider.getInstance(context.extensionUri, context);
  aiStatusBar = AIStatusBar.getInstance(apiClient);

  // Register AI commands
  context.subscriptions.push(
    vscode.commands.registerCommand('gorev.ai.showPanel', async () => {
      try {
        // Check if AI is configured
        const activeProjectResponse = await apiClient.getActiveProject();
        const projectId = activeProjectResponse.data?.id;

        if (!projectId) {
          vscode.window.showWarningMessage('AI features require an active project. Please activate a project first.');
          return;
        }

        // Show AI panel
        await aiPanelProvider.show(true);
      } catch (error) {
        Logger.error('[Extension] Failed to show AI panel:', error);
        vscode.window.showErrorMessage(`Failed to open AI panel: ${error}`);
      }
    }),

    vscode.commands.registerCommand('gorev.ai.configure', async () => {
      try {
        await vscode.commands.executeCommand('gorev.ai.showPanel');
        // Panel will open on Configure tab
      } catch (error) {
        Logger.error('[Extension] Failed to open AI configuration:', error);
        vscode.window.showErrorMessage(`Failed to open AI configuration: ${error}`);
      }
    }),

    vscode.commands.registerCommand('gorev.ai.chat', async (message?: string) => {
      try {
        if (!message) {
          message = await vscode.window.showInputBox({
            prompt: 'Enter your message for AI',
            placeHolder: 'Ask AI anything about your tasks...'
          });

          if (!message) {
            return; // User cancelled
          }
        }

        // Show AI panel on Chat tab
        await vscode.commands.executeCommand('gorev.ai.showPanel');
      } catch (error) {
        Logger.error('[Extension] Failed to send AI chat:', error);
        vscode.window.showErrorMessage(`Failed to send chat message: ${error}`);
      }
    }),

    vscode.commands.registerCommand('gorev.ai.analyzeProject', async () => {
      try {
        const activeProjectResponse = await apiClient.getActiveProject();
        const project = activeProjectResponse.data;

        if (!project) {
          vscode.window.showWarningMessage('AI features require an active project. Please activate a project first.');
          return;
        }

        // Get project statistics
        const tasksResponse = await apiClient.getProjectTasks(project.id);
        const tasks = tasksResponse.data || [];
        const completedCount = tasks.filter(t => t.durum === 'tamamlandi').length;
        const pendingCount = tasks.length - completedCount;

        // Call AI analyze via API
        const result = await apiClient.aiAnalyze({
          project_id: project.id,
          project_name: project.name,
          task_count: tasks.length,
          completed_count: completedCount,
          pending_count: pendingCount
        });

        if (result.success && result.data) {
          // Show results in AI panel
          await aiPanelProvider.show(true);

          // Send analyze results to panel
          vscode.commands.executeCommand('gorev.ai.showPanel', {
            action: 'analyze',
            data: result.data
          });
        } else {
          vscode.window.showWarningMessage('AI analysis is not configured for this project.');
        }
      } catch (error) {
        Logger.error('[Extension] Failed to analyze project with AI:', error);
        vscode.window.showErrorMessage(`Failed to analyze project: ${error}`);
      }
    }),

    vscode.commands.registerCommand('gorev.ai.estimateTask', async (taskId?: string) => {
      try {
        if (!taskId) {
          // Ask user to select a task
          const tasksResponse = await apiClient.getTasks({ tum_projeler: true, limit: 50 });
          const tasks = tasksResponse.data || [];

          if (tasks.length === 0) {
            vscode.window.showInformationMessage('No tasks found to estimate.');
            return;
          }

          const quickPickItems = tasks.map(task => ({
            label: task.baslik,
            description: task.proje_name || '',
            taskId: task.id
          }));

          const selected = await vscode.window.showQuickPick(quickPickItems, {
            placeHolder: 'Select a task to estimate time'
          });

          if (!selected) {
            return; // User cancelled
          }

          taskId = selected.taskId;
        }

        // Get task details
        const taskResponse = await apiClient.getTask(taskId);
        const task = taskResponse.data;

        if (!task) {
          vscode.window.showErrorMessage('Task not found.');
          return;
        }

        // Call AI estimate via API
        const result = await apiClient.aiEstimate({
          task_id: task.id,
          title: task.baslik,
          description: task.aciklama || '',
          tags: task.etiketler?.map(t => t.isim)
        });

        if (result.success && result.data) {
          const estimation = result.data;
          vscode.window.showInformationMessage(
            `Time Estimate for "${task.baslik}": ${estimation.estimated_hours}h ` +
            `(confidence: ${Math.round(estimation.confidence * 100)}%)\n\n` +
            `Reasoning: ${estimation.reasoning}`
          );
        } else {
          vscode.window.showWarningMessage('AI estimation is not configured for this project.');
        }
      } catch (error) {
        Logger.error('[Extension] Failed to estimate task time with AI:', error);
        vscode.window.showErrorMessage(`Failed to estimate task: ${error}`);
      }
    })
  );

  // Show AI status bar
  aiStatusBar.show();
  context.subscriptions.push(aiStatusBar);

  // Update AI status when active project changes
  apiClient.on('connected', async () => {
    try {
      const activeProjectResponse = await apiClient.getActiveProject();
      if (activeProjectResponse.data?.id) {
        await aiStatusBar.update(activeProjectResponse.data.id);
      }
    } catch (error) {
      Logger.warn('[Extension] Failed to update AI status:', error);
    }
  });

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

  // Initialize and start database file watcher for real-time sync
  databaseWatcher = new DatabaseFileWatcher(refreshManager);
  databaseWatcher.start();
  context.subscriptions.push({
    dispose: () => databaseWatcher.dispose(),
  });

  // Initialize WebSocket client for real-time updates (daemon mode)
  const workspaceContext = serverManager.getWorkspaceContext();
  if (workspaceContext && workspaceContext.workspaceId) {
    const daemonUrl = `http://${apiHost}:${apiPort}`;
    const outputChannel = vscode.window.createOutputChannel('Gorev WebSocket');
    webSocketClient = new WebSocketClient(daemonUrl, workspaceContext.workspaceId, outputChannel);
    webSocketClient.connect();

    context.subscriptions.push({
      dispose: () => {
        if (webSocketClient) {
          webSocketClient.disconnect();
        }
        outputChannel.dispose();
      }
    });

    Logger.info(`[Extension] WebSocket client initialized for workspace: ${workspaceContext.workspaceId}`);
  } else {
    Logger.warn('[Extension] WebSocket client not initialized - no workspace context');
  }

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

  // Dispose UnifiedServerManager (this will also disconnect apiClient and stop server)
  if (serverManager) {
    await serverManager.dispose();
  }

  if (statusBarManager) {
    statusBarManager.dispose();
  }

  // Dispose AI components
  if (aiStatusBar) {
    aiStatusBar.dispose();
  }
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