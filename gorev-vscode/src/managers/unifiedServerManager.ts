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
 */

import * as vscode from 'vscode';
import * as crypto from 'crypto';
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

export class UnifiedServerManager extends EventEmitter {
  private apiClient: ApiClient;
  private workspaceContext: WorkspaceContext | undefined;
  private serverStatus: ServerStatus;
  private healthCheckInterval?: NodeJS.Timeout;
  private serverProcess?: ChildProcess;
  private readonly HEALTH_CHECK_INTERVAL = 30000; // 30 seconds
  private readonly SERVER_START_TIMEOUT = 60000; // 60 seconds (increased for first-time setup)
  private readonly PORT_CHECK_RETRY_INTERVAL = 1000; // 1 second

  constructor(
    private readonly apiHost = 'localhost',
    private readonly apiPort = 5082
  ) {
    super();

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
   */
  async registerWorkspace(workspaceFolder: vscode.WorkspaceFolder): Promise<WorkspaceInfo> {
    try {
      const workspacePath = workspaceFolder.uri.fsPath;
      const workspaceName = workspaceFolder.name;

      Logger.info(`[UnifiedServerManager] Registering workspace: ${workspaceName} at ${workspacePath}`);

      // Call workspace registration endpoint
      const response = await this.apiClient.registerWorkspace({
        path: workspacePath,
        name: workspaceName
      });

      if (!response.success) {
        throw new Error('Workspace registration failed: Unknown error');
      }

      const workspaceInfo = response.workspace;

      // Set workspace context
      this.workspaceContext = {
        workspaceId: response.workspace_id,
        workspacePath: workspaceInfo.path,
        workspaceName: workspaceInfo.name
      };

      // Inject workspace headers into API client
      this.apiClient.setWorkspaceHeaders(this.workspaceContext);

      this.serverStatus.workspaceRegistered = true;
      Logger.info(`[UnifiedServerManager] Workspace registered successfully: ID=${response.workspace_id}`);

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
   * Dispose resources and cleanup
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

    // Stop server if we started it
    if (this.serverProcess) {
      await this.stopServer();
    }

    // Clear workspace context
    this.workspaceContext = undefined;
    this.serverStatus.connected = false;
    this.serverStatus.workspaceRegistered = false;

    Logger.info('[UnifiedServerManager] Disposed');
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
   * Generate workspace ID from path (matches server-side logic)
   */
  private generateWorkspaceId(path: string): string {
    const hash = crypto.createHash('sha256').update(path).digest('hex');
    return hash.substring(0, 16); // First 8 bytes (16 hex chars)
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
   * Start the Gorev server process
   */
  private async startServer(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        Logger.info('[UnifiedServerManager] Starting Gorev server...');

        // Show user notification
        vscode.window.showInformationMessage('Gorev: Starting server...', 'Show Logs').then((selection) => {
          if (selection === 'Show Logs') {
            vscode.commands.executeCommand('workbench.action.output.show');
          }
        });

        // Use npx to run the server
        const command = process.platform === 'win32' ? 'npx.cmd' : 'npx';
        // Note: Not passing --api-port as it defaults to 5082 (matches apiPort setting)
        // This ensures compatibility with all versions of the binary
        const args = ['@mehmetsenol/gorev-mcp-server', 'serve', '--debug'];

        // Determine database path
        // Priority: Workspace folder > User home directory
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

        Logger.info(`[UnifiedServerManager] Running command: ${command} ${args.join(' ')} (port: ${this.apiPort})`);
        Logger.info(`[UnifiedServerManager] Database path: ${dbPath}`);

        let errorOutput = '';

        this.serverProcess = spawn(command, args, {
          // Note: stdin must be 'pipe' (not 'ignore') to keep MCP server alive
          // We don't send any MCP commands, but the server needs stdin open
          stdio: ['pipe', 'pipe', 'pipe'],
          shell: process.platform === 'win32',
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

            // Check for common errors
            if (output.includes('command not found') || output.includes('not found')) {
              Logger.error(`[UnifiedServerManager] npx command failed: ${output}`);
            } else if (output.includes('ENOENT')) {
              Logger.error(`[UnifiedServerManager] Package not found: ${output}`);
            } else {
              Logger.warn(`[Gorev Server] ${output}`);
            }
          });
        }

        this.serverProcess.on('error', (error) => {
          Logger.error('[UnifiedServerManager] Server process spawn error:', error);
          vscode.window.showErrorMessage(
            `Gorev: Failed to start server. ${error.message}. ` +
            'Please ensure @mehmetsenol/gorev-mcp-server is installed.',
            'Install Package',
            'Show Logs'
          ).then((selection) => {
            if (selection === 'Install Package') {
              vscode.window.showInformationMessage(
                'Run in terminal: npm install -g @mehmetsenol/gorev-mcp-server'
              );
            } else if (selection === 'Show Logs') {
              vscode.commands.executeCommand('workbench.action.output.show');
            }
          });
          reject(new Error(`Failed to start server: ${error.message}`));
        });

        this.serverProcess.on('exit', (code, signal) => {
          Logger.info(`[UnifiedServerManager] Server process exited with code ${code}, signal ${signal}`);

          // If process exited immediately with error, show helpful message
          if (code !== 0 && code !== null) {
            const errorMsg = errorOutput || 'Unknown error';
            Logger.error(`[UnifiedServerManager] Server startup failed: ${errorMsg}`);

            if (errorMsg.includes('not found') || errorMsg.includes('ENOENT')) {
              vscode.window.showErrorMessage(
                'Gorev: Server package not found. Please install it first.',
                'Install Instructions'
              ).then((selection) => {
                if (selection === 'Install Instructions') {
                  vscode.window.showInformationMessage(
                    'Run: npm install -g @mehmetsenol/gorev-mcp-server'
                  );
                }
              });
            }
          }

          this.serverProcess = undefined;
        });

        // Give server a moment to start and check for immediate failures
        setTimeout(() => {
          if (this.serverProcess && !this.serverProcess.killed) {
            Logger.info('[UnifiedServerManager] Server process spawned successfully');
            resolve();
          } else {
            reject(new Error('Server process terminated immediately after spawn'));
          }
        }, 2000); // Increased from 1000ms to 2000ms for better reliability

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
