/**
 * UnifiedServerManager
 *
 * Manages both MCP server connection and REST API client for VS Code extension.
 * Provides workspace registration and context management.
 *
 * Key responsibilities:
 * - Initialize and manage API client connection
 * - Register workspace with server on activation
 * - Inject workspace headers into API requests
 * - Coordinate server lifecycle (start/stop/health checks)
 * - Auto-persist workspaceId to workspace settings for consistency
 */

import * as vscode from 'vscode';
import * as net from 'net';
import * as fs from 'fs';
import * as path from 'path';
import { spawn, ChildProcess } from 'child_process';
import { EventEmitter } from 'events';
import { ApiClient } from '../api/client';
import { Logger } from '../utils/logger';
import {
  WorkspaceContext,
  WorkspaceInfo
} from '../models/workspace';

export interface ServerStatus {
  connected: boolean;
  workspaceRegistered: boolean;
  lastHealthCheck?: Date;
  serverVersion?: string;
}

const WORKSPACE_ID_CONFIG_KEY = 'gorev.workspaceId';

export class UnifiedServerManager extends EventEmitter {
  private apiClient: ApiClient;
  private workspaceContext: WorkspaceContext | undefined;
  private serverStatus: ServerStatus;
  private healthCheckInterval?: NodeJS.Timeout;
  private serverProcess?: ChildProcess;
  private readonly HEALTH_CHECK_INTERVAL = 30000; // 30 seconds
  private readonly SERVER_START_TIMEOUT = 60000; // 60 seconds (increased for first-time setup)
  private readonly PORT_CHECK_RETRY_INTERVAL = 1000; // 1 second
  private readonly connectionMode: string;
  private readonly localServerPath: string | undefined;

  constructor(
    private readonly apiHost = 'localhost',
    private readonly apiPort = 5082
  ) {
    super();

    // Read connection mode and local server path from VS Code configuration
    const config = vscode.workspace.getConfiguration();
    this.connectionMode = config.get<string>('gorev.connectionMode', 'auto');
    this.localServerPath = config.get<string>('gorev.serverPath', '');

    Logger.info(`[UnifiedServerManager] Connection mode: ${this.connectionMode}`);
    if (this.connectionMode === 'local') {
      if (this.localServerPath) {
        Logger.info(`[UnifiedServerManager] Local server path: ${this.localServerPath}`);
      } else {
        Logger.warn('[UnifiedServerManager] Local mode enabled but gorev.serverPath not configured, falling back to npm package');
      }
    }

    // Initialize API client
    this.apiClient = new ApiClient(`http://${apiHost}:${apiPort}`);

    // Initialize server status
    this.serverStatus = {
      connected: false,
      workspaceRegistered: false
    };

    // Listen to API client events
    this.setupApiClientListeners();
  }

  /**
   * Get the API client instance
   */
  getApiClient(): ApiClient {
    return this.apiClient;
  }

  /**
   * Get current workspace context
   */
  getWorkspaceContext(): WorkspaceContext | undefined {
    return this.workspaceContext;
  }

  /**
   * Get current server status
   */
  getServerStatus(): ServerStatus {
    return { ...this.serverStatus };
  }

  /**
   * Initialize server connection and register workspace
   */
  async initialize(): Promise<void> {
    try {
      // Step 0: Check if server is already running, start if needed
      const isRunning = await this.isServerRunning();
      if (!isRunning) {
        Logger.info('[UnifiedServerManager] Server not running, starting automatically...');
        await this.startServer();
        await this.waitForServerReady();
      } else {
        Logger.info('[UnifiedServerManager] Server already running');
      }

      // Step 1: Connect to API server
      Logger.info('[UnifiedServerManager] Connecting to API server...');
      await this.apiClient.connect();
      this.serverStatus.connected = true;
      this.serverStatus.lastHealthCheck = new Date();
      Logger.info('[UnifiedServerManager] Connected to API server');

      // Step 2: Register current workspace
      const workspaceFolder = this.getCurrentWorkspaceFolder();
      if (workspaceFolder) {
        Logger.info(`[UnifiedServerManager] Registering workspace: ${workspaceFolder.uri.fsPath}`);
        await this.registerWorkspace(workspaceFolder);
      } else {
        Logger.warn('[UnifiedServerManager] No workspace folder found, skipping registration');
      }

      // Step 3: Start health check monitoring
      this.startHealthCheckMonitoring();

      // Emit initialized event
      this.emit('initialized', this.serverStatus);

    } catch (error) {
      Logger.error('[UnifiedServerManager] Initialization failed:', error);
      this.serverStatus.connected = false;
      this.serverStatus.workspaceRegistered = false;
      throw error;
    }
  }

  /**
   * Register workspace with server
   *
   * Flow:
   * 1. Check if workspaceId already exists in workspace settings
   * 2. If yes, use existing ID (maintains consistency across sessions)
   * 3. If no, register with server, get returned ID, save to workspace settings
   */
  async registerWorkspace(workspaceFolder: vscode.WorkspaceFolder): Promise<WorkspaceInfo> {
    try {
      const workspacePath = workspaceFolder.uri.fsPath;
      const workspaceName = workspaceFolder.name;

      // Step 1: Check for existing workspaceId in workspace settings
      let workspaceId = this.getSavedWorkspaceId();
      const isNewWorkspace = !workspaceId;

      if (workspaceId) {
        Logger.info(`[UnifiedServerManager] Using saved workspaceId: ${workspaceId}`);
      } else {
        Logger.info(`[UnifiedServerManager] No saved workspaceId, will register new workspace: ${workspaceName}`);
      }

      // Step 2: Call workspace registration endpoint
      const response = await this.apiClient.registerWorkspace({
        path: workspacePath,
        name: workspaceName,
        workspace_id: workspaceId // Use saved ID or let server generate one
      });

      if (!response.success) {
        throw new Error('Workspace registration failed: Unknown error');
      }

      const workspaceInfo = response.workspace;
      const finalWorkspaceId = workspaceId || response.workspace_id;

      // Step 3: Save workspaceId to workspace settings if this is a new registration
      if (isNewWorkspace && finalWorkspaceId) {
        await this.saveWorkspaceId(finalWorkspaceId);
        Logger.info(`[UnifiedServerManager] Saved new workspaceId to workspace settings: ${finalWorkspaceId}`);
      }

      // Set workspace context
      this.workspaceContext = {
        workspaceId: finalWorkspaceId,
        workspacePath: workspaceInfo.path || workspacePath,
        workspaceName: workspaceInfo.name || workspaceName
      };

      // Inject workspace headers into API client
      this.apiClient.setWorkspaceHeaders(this.workspaceContext);

      this.serverStatus.workspaceRegistered = true;
      Logger.info(`[UnifiedServerManager] Workspace registered successfully: ID=${this.workspaceContext.workspaceId}`);

      // Emit workspace registered event
      this.emit('workspaceRegistered', workspaceInfo);

      return workspaceInfo;

    } catch (error) {
      Logger.error('[UnifiedServerManager] Workspace registration failed:', error);
      this.serverStatus.workspaceRegistered = false;
      throw error;
    }
  }

  /**
   * Get saved workspaceId from workspace settings
   */
  private getSavedWorkspaceId(): string | undefined {
    const config = vscode.workspace.getConfiguration();
    const workspaceId = config.get<string>(WORKSPACE_ID_CONFIG_KEY);
    return workspaceId && workspaceId.trim() !== '' ? workspaceId : undefined;
  }

  /**
   * Save workspaceId to workspace settings (.vscode/settings.json)
   */
  private async saveWorkspaceId(workspaceId: string): Promise<void> {
    try {
      const config = vscode.workspace.getConfiguration();
      // Save to workspace scope (creates/updates .vscode/settings.json)
      await config.update(
        WORKSPACE_ID_CONFIG_KEY,
        workspaceId,
        vscode.ConfigurationTarget.Workspace
      );
      Logger.info(`[UnifiedServerManager] WorkspaceId saved to .vscode/settings.json`);
    } catch (error) {
      // Log but don't fail - workspaceId will work for this session
      Logger.warn('[UnifiedServerManager] Failed to save workspaceId to settings:', error);
    }
  }

  /**
   * Dispose resources and cleanup
   * Rule 15: Smart shutdown - check active clients before stopping
   */
  async dispose(): Promise<void> {
    // Stop health check monitoring
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
      this.healthCheckInterval = undefined;
    }

    // Disconnect API client
    if (this.apiClient) {
      this.apiClient.disconnect();
    }

    // SMART SHUTDOWN: Check active clients before stopping server
    if (this.serverProcess) {
      const shouldStop = await this.shouldStopServer();
      if (shouldStop) {
        Logger.info('[UnifiedServerManager] No other clients detected, stopping server...');
        await this.stopServer();
      } else {
        Logger.info('[UnifiedServerManager] Other clients active, leaving server running...');
        this.serverProcess = undefined; // Just clear the reference, don't stop
      }
    }

    // Clear workspace context
    this.workspaceContext = undefined;
    this.serverStatus.connected = false;
    this.serverStatus.workspaceRegistered = false;

    Logger.info('[UnifiedServerManager] Disposed');
  }

  /**
   * Check if server should be stopped based on active client count
   * Rule 15: Proper root cause analysis - don't kill other clients' connections
   */
  private async shouldStopServer(): Promise<boolean> {
    try {
      const clientCount = await this.getActiveClientCount();
      // VS Code itself counts as 1 client, so >1 means other clients exist
      if (clientCount > 1) {
        Logger.info(`[UnifiedServerManager] ${clientCount - 1} other client(s) active, not stopping server`);
        return false;
      }
      // No other clients, safe to stop
      Logger.info('[UnifiedServerManager] No other clients detected');
      return true;
    } catch (error) {
      Logger.warn('[UnifiedServerManager] Failed to check client count, stopping server:', error);
      // Conservative approach: if we can't check, stop the server
      // This ensures no orphaned daemon processes
      return true;
    }
  }

  /**
   * Get active client count from daemon
   */
  private async getActiveClientCount(): Promise<number> {
    try {
      return await this.apiClient.getActiveClientCount();
    } catch (error) {
      Logger.warn('[UnifiedServerManager] Failed to get client count:', error);
      return -1; // Indicate check failed
    }
  }

  /**
   * Private helper methods
   */

  private getCurrentWorkspaceFolder(): vscode.WorkspaceFolder | undefined {
    if (vscode.workspace.workspaceFolders && vscode.workspace.workspaceFolders.length > 0) {
      // For now, use the first workspace folder
      // In multi-root workspaces, we could let user choose
      return vscode.workspace.workspaceFolders[0];
    }
    return undefined;
  }

  private setupApiClientListeners(): void {
    this.apiClient.on('connected', () => {
      this.serverStatus.connected = true;
      this.serverStatus.lastHealthCheck = new Date();
      this.emit('serverConnected');
    });

    this.apiClient.on('disconnected', () => {
      this.serverStatus.connected = false;
      this.serverStatus.workspaceRegistered = false;
      this.emit('serverDisconnected');
    });

    this.apiClient.on('error', (error: Error) => {
      Logger.error('[UnifiedServerManager] API client error:', error);
      this.emit('serverError', error);
    });
  }

  private startHealthCheckMonitoring(): void {
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
    }

    this.healthCheckInterval = setInterval(async () => {
      try {
        // Simple health check - try to connect
        if (!this.apiClient.isConnected()) {
          await this.apiClient.connect();
        }
        this.serverStatus.lastHealthCheck = new Date();
      } catch (error) {
        Logger.warn('[UnifiedServerManager] Health check failed:', error);
        this.serverStatus.connected = false;
        this.serverStatus.workspaceRegistered = false;
        this.emit('healthCheckFailed', error);
      }
    }, this.HEALTH_CHECK_INTERVAL);
  }

  /**
   * Check if server is running by testing both port and health endpoint
   * Port check alone is unreliable (TIME_WAIT sockets, zombie processes)
   */
  private async isServerRunning(): Promise<boolean> {
    // First check if port is listening
    const portOpen = await this.checkPort();
    if (!portOpen) {
      return false;
    }

    // Then verify health endpoint responds correctly
    try {
      const response = await this.apiClient.checkHealth();
      return response !== undefined;
    } catch (error) {
      Logger.debug('[UnifiedServerManager] Health check failed, server not running:', error);
      return false;
    }
  }

  /**
   * Check if port is listening (TCP connection test)
   */
  private async checkPort(): Promise<boolean> {
    return new Promise((resolve) => {
      const socket = new net.Socket();

      socket.setTimeout(1000); // 1 second timeout

      socket.on('connect', () => {
        socket.destroy();
        resolve(true);
      });

      socket.on('timeout', () => {
        socket.destroy();
        resolve(false);
      });

      socket.on('error', () => {
        resolve(false);
      });

      socket.connect(this.apiPort, this.apiHost);
    });
  }

  /**
   * Start the Gorev server process based on connection mode
   */
  private async startServer(): Promise<void> {
    // Handle remote mode - don't start server, expect it to be running
    if (this.connectionMode === 'remote') {
      Logger.warn('[UnifiedServerManager] Remote mode enabled - server should be started manually');
      vscode.window.showWarningMessage(
        'Gorev: Remote mode enabled. Please ensure the Gorev server is running on the remote host.',
        'OK'
      );
      return;
    }

    return new Promise((resolve, reject) => {
      try {
        Logger.info('[UnifiedServerManager] Starting Gorev server...');
        vscode.window.showInformationMessage('Gorev: Starting server...', 'Show Logs').then((selection) => {
          if (selection === 'Show Logs') {
            vscode.commands.executeCommand('workbench.action.output.show');
          }
        });

        // Determine command and args based on connection mode
        let command: string;
        let args: string[];
        let useShell = false;

        if (this.connectionMode === 'local') {
          // Local mode: use gorev.serverPath or try to find gorev in PATH
          if (this.localServerPath && fs.existsSync(this.localServerPath)) {
            command = this.localServerPath;
            args = ['serve', '--debug'];
            Logger.info(`[UnifiedServerManager] Using local binary: ${command}`);
          } else {
            // Try to use gorev from PATH
            command = 'gorev';
            args = ['serve', '--debug'];
            Logger.info(`[UnifiedServerManager] Using gorev from PATH`);
          }
          useShell = false;
        } else if (this.connectionMode === 'docker') {
          // Docker mode: use docker-compose
          command = 'docker-compose';
          const composeFile = this.localServerPath || './docker-compose.yml';
          args = ['-f', composeFile, 'up', '-d'];
          useShell = true;
          Logger.info(`[UnifiedServerManager] Using Docker: ${command} ${args.join(' ')}`);
        } else {
          // Auto mode (default): use npm package
          command = process.platform === 'win32' ? 'npx.cmd' : 'npx';
          args = ['@mehmetsenol/gorev-mcp-server', 'serve', '--debug'];
          useShell = process.platform === 'win32';
          Logger.info(`[UnifiedServerManager] Using npm package: ${command} ${args.join(' ')}`);
        }

        // Determine database path
        const workspaceFolder = this.getCurrentWorkspaceFolder();
        const dbPath = workspaceFolder
          ? path.join(workspaceFolder.uri.fsPath, '.gorev', 'gorev.db')
          : path.join(process.env.HOME || process.env.USERPROFILE || '', '.gorev', 'gorev.db');

        // Ensure .gorev directory exists
        const dbDir = path.dirname(dbPath);
        if (!fs.existsSync(dbDir)) {
          Logger.info(`[UnifiedServerManager] Creating database directory: ${dbDir}`);
          fs.mkdirSync(dbDir, { recursive: true });
        }

        // Set environment variables for server process
        const env = {
          ...process.env,
          GOREV_DB_PATH: dbPath
        };

        Logger.info(`[UnifiedServerManager] Database path: ${dbPath}`);
        Logger.info(`[UnifiedServerManager] Running command: ${command} ${args.join(' ')} (port: ${this.apiPort})`);

        let errorOutput = '';

        this.serverProcess = spawn(command, args, {
          stdio: ['pipe', 'pipe', 'pipe'],
          shell: useShell,
          env: env
        });

        // Log server output and detect errors
        if (this.serverProcess.stdout) {
          this.serverProcess.stdout.on('data', (data) => {
            const output = data.toString().trim();
            Logger.info(`[Gorev Server] ${output}`);
          });
        }

        if (this.serverProcess.stderr) {
          this.serverProcess.stderr.on('data', (data) => {
            const output = data.toString().trim();
            errorOutput += output + '\n';

            if (output.includes('command not found') || output.includes('not found')) {
              Logger.error(`[UnifiedServerManager] Command failed: ${output}`);
            } else if (output.includes('ENOENT')) {
              Logger.error(`[UnifiedServerManager] File/executable not found: ${output}`);
            } else {
              Logger.warn(`[Gorev Server] ${output}`);
            }
          });
        }

        this.serverProcess.on('error', (error) => {
          Logger.error('[UnifiedServerManager] Server process spawn error:', error);
          this.showStartError(error);
          reject(new Error(`Failed to start server: ${error.message}`));
        });

        this.serverProcess.on('exit', (code, signal) => {
          Logger.info(`[UnifiedServerManager] Server process exited with code ${code}, signal ${signal}`);
          if (code !== 0 && code !== null) {
            const errorMsg = errorOutput || 'Unknown error';
            Logger.error(`[UnifiedServerManager] Server startup failed: ${errorMsg}`);
            this.handleStartupError(errorMsg);
          }
          this.serverProcess = undefined;
        });

        // Give server a moment to start
        setTimeout(() => {
          if (this.serverProcess && !this.serverProcess.killed) {
            Logger.info('[UnifiedServerManager] Server process spawned successfully');
            resolve();
          } else {
            reject(new Error('Server process terminated immediately after spawn'));
          }
        }, 2000);

      } catch (error) {
        Logger.error('[UnifiedServerManager] Failed to start server:', error);
        vscode.window.showErrorMessage(
          `Gorev: Failed to start server. ${error instanceof Error ? error.message : 'Unknown error'}`,
          'Show Logs'
        ).then((selection) => {
          if (selection === 'Show Logs') {
            vscode.commands.executeCommand('workbench.action.output.show');
          }
        });
        reject(error);
      }
    });
  }

  /**
   * Show error message when server fails to start
   */
  private showStartError(error: Error): void {
    let message = `Gorev: Failed to start server. ${error.message}.`;
    let actions: string[] = ['Show Logs'];

    if (this.connectionMode === 'local') {
      message += ' Please configure gorev.serverPath or ensure gorev is in PATH.';
      actions.push('Configure');
    } else if (this.connectionMode === 'docker') {
      message += ' Please ensure Docker and docker-compose are installed.';
      actions.push('Docker Info');
    } else {
      message += ' Please ensure @mehmetsenol/gorev-mcp-server is installed.';
      actions.push('Install Package');
    }

    vscode.window.showErrorMessage(message, ...actions).then((selection) => {
      if (selection === 'Show Logs') {
        vscode.commands.executeCommand('workbench.action.output.show');
      } else if (selection === 'Install Package') {
        vscode.window.showInformationMessage('Run: npm install -g @mehmetsenol/gorev-mcp-server');
      } else if (selection === 'Configure') {
        vscode.commands.executeCommand('workbench.action.openSettings', 'gorev.serverPath');
      } else if (selection === 'Docker Info') {
        vscode.window.showInformationMessage('Install Docker from: https://docs.docker.com/get-docker/');
      }
    });
  }

  /**
   * Handle server startup errors with helpful messages
   */
  private handleStartupError(errorMsg: string): void {
    if (errorMsg.includes('not found') || errorMsg.includes('ENOENT')) {
      if (this.connectionMode === 'local') {
        vscode.window.showErrorMessage(
          'Gorev: Local binary not found. Please configure gorev.serverPath in settings.',
          'Open Settings'
        ).then((selection) => {
          if (selection === 'Open Settings') {
            vscode.commands.executeCommand('workbench.action.openSettings', 'gorev.serverPath');
          }
        });
      } else if (this.connectionMode === 'docker') {
        vscode.window.showErrorMessage(
          'Gorev: Docker or docker-compose not found. Please install Docker.',
          'Docker Install'
        ).then((selection) => {
          if (selection === 'Docker Install') {
            vscode.env.openExternal(vscode.Uri.parse('https://docs.docker.com/get-docker/'));
          }
        });
      } else {
        vscode.window.showErrorMessage(
          'Gorev: Server package not found. Please install it first.',
          'Install Instructions'
        ).then((selection) => {
          if (selection === 'Install Instructions') {
            vscode.window.showInformationMessage('Run: npm install -g @mehmetsenol/gorev-mcp-server');
          }
        });
      }
    }
  }

  /**
   * Wait for server to be ready by polling the port
   */
  private async waitForServerReady(): Promise<void> {
    const startTime = Date.now();

    while (Date.now() - startTime < this.SERVER_START_TIMEOUT) {
      const isRunning = await this.isServerRunning();
      if (isRunning) {
        Logger.info('[UnifiedServerManager] Server is ready');
        return;
      }

      // Wait before next check
      await new Promise(resolve => setTimeout(resolve, this.PORT_CHECK_RETRY_INTERVAL));
    }

    throw new Error(`Server failed to start within ${this.SERVER_START_TIMEOUT}ms timeout`);
  }

  /**
   * Stop the server process gracefully
   */
  private async stopServer(): Promise<void> {
    if (!this.serverProcess) {
      return;
    }

    return new Promise((resolve) => {
      Logger.info('[UnifiedServerManager] Stopping Gorev server...');

      const timeout = setTimeout(() => {
        if (this.serverProcess && !this.serverProcess.killed) {
          Logger.warn('[UnifiedServerManager] Server did not stop gracefully, forcing kill');
          this.serverProcess.kill('SIGKILL');
        }
        resolve();
      }, 5000); // 5 second timeout for graceful shutdown

      if (this.serverProcess) {
        this.serverProcess.once('exit', () => {
          clearTimeout(timeout);
          Logger.info('[UnifiedServerManager] Server stopped successfully');
          this.serverProcess = undefined;
          resolve();
        });

        // Send SIGTERM for graceful shutdown
        this.serverProcess.kill('SIGTERM');
      } else {
        clearTimeout(timeout);
        resolve();
      }
    });
  }
}
