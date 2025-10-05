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
  WorkspaceInfo,
  WorkspaceRegistration,
  WorkspaceRegistrationResponse
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
  private readonly SERVER_START_TIMEOUT = 15000; // 15 seconds
  private readonly PORT_CHECK_RETRY_INTERVAL = 500; // 0.5 seconds

  constructor(
    private readonly apiHost: string = 'localhost',
    private readonly apiPort: number = 5082
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
   * Check if server is running by testing if port is listening
   */
  private async isServerRunning(): Promise<boolean> {
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

        this.serverProcess = spawn(command, args, {
          // Note: stdin must be 'pipe' (not 'ignore') to keep MCP server alive
          // We don't send any MCP commands, but the server needs stdin open
          stdio: ['pipe', 'pipe', 'pipe'],
          shell: process.platform === 'win32',
          env: env
        });

        // Log server output
        if (this.serverProcess.stdout) {
          this.serverProcess.stdout.on('data', (data) => {
            Logger.info(`[Gorev Server] ${data.toString().trim()}`);
          });
        }

        if (this.serverProcess.stderr) {
          this.serverProcess.stderr.on('data', (data) => {
            Logger.warn(`[Gorev Server] ${data.toString().trim()}`);
          });
        }

        this.serverProcess.on('error', (error) => {
          Logger.error('[UnifiedServerManager] Server process error:', error);
          reject(new Error(`Failed to start server: ${error.message}`));
        });

        this.serverProcess.on('exit', (code, signal) => {
          Logger.info(`[UnifiedServerManager] Server process exited with code ${code}, signal ${signal}`);
          this.serverProcess = undefined;
        });

        // Give server a moment to start
        setTimeout(() => {
          Logger.info('[UnifiedServerManager] Server process spawned successfully');
          resolve();
        }, 1000);

      } catch (error) {
        Logger.error('[UnifiedServerManager] Failed to start server:', error);
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
