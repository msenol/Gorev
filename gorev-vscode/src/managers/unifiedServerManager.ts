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
  private readonly HEALTH_CHECK_INTERVAL = 30000; // 30 seconds

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
  dispose(): void {
    // Stop health check monitoring
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
      this.healthCheckInterval = undefined;
    }

    // Disconnect API client
    if (this.apiClient) {
      this.apiClient.disconnect();
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
}
