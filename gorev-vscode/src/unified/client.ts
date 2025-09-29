import { EventEmitter } from 'events';
import { Logger } from '../utils/logger';
import { ClientInterface } from '../interfaces/client';
import { MCPClient } from '../mcp/client';
import { ApiClient } from '../api/client';
import { MCPToolResult } from '../mcp/types';
import * as vscode from 'vscode';

export type ClientMode = 'api' | 'mcp' | 'auto';

export interface UnifiedClientConfig {
  mode: ClientMode;
  apiBaseURL?: string;
  mcpServerPath?: string;
  retryAttempts?: number;
  retryDelay?: number;
}

export class UnifiedClient extends EventEmitter implements ClientInterface {
  private mcpClient: MCPClient;
  private apiClient: ApiClient;
  private currentMode: 'api' | 'mcp' | null = null;
  private config: UnifiedClientConfig;
  private connected = false;
  private autoDetectionInProgress = false;

  constructor(config: UnifiedClientConfig = { mode: 'auto' }) {
    super();

    this.config = {
      apiBaseURL: 'http://localhost:5082',
      retryAttempts: 3,
      retryDelay: 1000,
      ...config
    };

    this.mcpClient = new MCPClient();
    this.apiClient = new ApiClient(this.config.apiBaseURL);

    this.setupEventHandlers();
  }

  private setupEventHandlers(): void {
    // MCP Client events
    this.mcpClient.on('connected', () => {
      if (this.currentMode === 'mcp') {
        this.connected = true;
        this.emit('connected', { mode: 'mcp' });
        Logger.info('[UnifiedClient] Connected via MCP');
      }
    });

    this.mcpClient.on('disconnected', (data) => {
      if (this.currentMode === 'mcp') {
        this.connected = false;
        this.emit('disconnected', { mode: 'mcp', ...data });
        Logger.info('[UnifiedClient] Disconnected from MCP');
      }
    });

    this.mcpClient.on('error', (error) => {
      if (this.currentMode === 'mcp') {
        this.emit('error', { mode: 'mcp', error });
        Logger.error('[UnifiedClient] MCP error:', error);
      }
    });

    // API Client events
    this.apiClient.on('connected', () => {
      if (this.currentMode === 'api') {
        this.connected = true;
        this.emit('connected', { mode: 'api' });
        Logger.info('[UnifiedClient] Connected via API');
      }
    });

    this.apiClient.on('disconnected', () => {
      if (this.currentMode === 'api') {
        this.connected = false;
        this.emit('disconnected', { mode: 'api' });
        Logger.info('[UnifiedClient] Disconnected from API');
      }
    });

    this.apiClient.on('error', (error) => {
      if (this.currentMode === 'api') {
        this.emit('error', { mode: 'api', error });
        Logger.error('[UnifiedClient] API error:', error);
      }
    });
  }

  async connect(): Promise<void> {
    if (this.connected) {
      Logger.warn('[UnifiedClient] Already connected');
      return;
    }

    Logger.info(`[UnifiedClient] Connecting with mode: ${this.config.mode}`);

    try {
      switch (this.config.mode) {
        case 'api':
          await this.connectAPI();
          break;
        case 'mcp':
          await this.connectMCP();
          break;
        case 'auto':
          await this.autoDetectAndConnect();
          break;
        default:
          throw new Error(`Invalid mode: ${this.config.mode}`);
      }
    } catch (error) {
      Logger.error('[UnifiedClient] Connection failed:', error);
      throw error;
    }
  }

  private async connectAPI(): Promise<void> {
    Logger.info('[UnifiedClient] Attempting API connection...');
    await this.apiClient.connect();
    this.currentMode = 'api';
  }

  private async connectMCP(): Promise<void> {
    Logger.info('[UnifiedClient] Attempting MCP connection...');
    await this.mcpClient.connect(this.config.mcpServerPath);
    this.currentMode = 'mcp';
  }

  private async autoDetectAndConnect(): Promise<void> {
    if (this.autoDetectionInProgress) {
      Logger.warn('[UnifiedClient] Auto-detection already in progress');
      return;
    }

    this.autoDetectionInProgress = true;
    Logger.info('[UnifiedClient] Starting auto-detection...');

    try {
      // First, try to detect if API server is running by checking different ports
      const apiPorts = [5082, 5080, 5081, 5083];
      let apiConnected = false;

      for (const port of apiPorts) {
        try {
          Logger.debug(`[UnifiedClient] Trying API on port ${port}...`);
          const testApiClient = new ApiClient(`http://localhost:${port}`);
          await testApiClient.connect();
          testApiClient.disconnect();

          // Update our API client to use this port
          this.apiClient = new ApiClient(`http://localhost:${port}`);
          this.setupApiEventHandlers();

          await this.connectAPI();
          apiConnected = true;
          Logger.info(`[UnifiedClient] Auto-detected API server on port ${port}`);
          break;
        } catch (error) {
          Logger.debug(`[UnifiedClient] Port ${port} not available:`, error);
          continue;
        }
      }

      // If API failed, fall back to MCP
      if (!apiConnected) {
        Logger.info('[UnifiedClient] API not available, falling back to MCP...');
        try {
          await this.connectMCP();
          Logger.info('[UnifiedClient] Auto-detected MCP server');
        } catch (mcpError) {
          Logger.error('[UnifiedClient] Both API and MCP connections failed');
          throw new Error('Could not connect to either API or MCP server');
        }
      }
    } finally {
      this.autoDetectionInProgress = false;
    }
  }

  private setupApiEventHandlers(): void {
    // Re-setup event handlers for the new API client instance
    this.apiClient.on('connected', () => {
      if (this.currentMode === 'api') {
        this.connected = true;
        this.emit('connected', { mode: 'api' });
        Logger.info('[UnifiedClient] Connected via API');
      }
    });

    this.apiClient.on('disconnected', () => {
      if (this.currentMode === 'api') {
        this.connected = false;
        this.emit('disconnected', { mode: 'api' });
        Logger.info('[UnifiedClient] Disconnected from API');
      }
    });

    this.apiClient.on('error', (error) => {
      if (this.currentMode === 'api') {
        this.emit('error', { mode: 'api', error });
        Logger.error('[UnifiedClient] API error:', error);
      }
    });
  }

  disconnect(): void {
    Logger.info(`[UnifiedClient] Disconnecting from ${this.currentMode} mode`);

    if (this.currentMode === 'api') {
      this.apiClient.disconnect();
    } else if (this.currentMode === 'mcp') {
      this.mcpClient.disconnect();
    }

    this.connected = false;
    this.currentMode = null;
    this.emit('disconnected', { mode: 'unified' });
  }

  isConnected(): boolean {
    return this.connected;
  }

  getCurrentMode(): 'api' | 'mcp' | null {
    return this.currentMode;
  }

  async callTool(name: string, params?: any): Promise<MCPToolResult | any> {
    if (!this.connected) {
      throw new Error('Not connected to any server');
    }

    Logger.info(`[UnifiedClient] Calling tool ${name} via ${this.currentMode} mode`);

    try {
      if (this.currentMode === 'api') {
        return await this.apiClient.callTool(name, params);
      } else if (this.currentMode === 'mcp') {
        return await this.mcpClient.callTool(name, params);
      } else {
        throw new Error('No active connection mode');
      }
    } catch (error) {
      Logger.error(`[UnifiedClient] Tool call failed for ${name}:`, error);

      // If API fails and we're in auto mode, try to fall back to MCP
      if (this.currentMode === 'api' && this.config.mode === 'auto') {
        Logger.warn('[UnifiedClient] API call failed, attempting MCP fallback...');
        try {
          // Disconnect from API and try MCP
          this.apiClient.disconnect();
          await this.connectMCP();
          return await this.mcpClient.callTool(name, params);
        } catch (fallbackError) {
          Logger.error('[UnifiedClient] MCP fallback also failed:', fallbackError);
          throw error; // Throw original error
        }
      }

      throw error;
    }
  }

  getTools(): any[] {
    if (this.currentMode === 'mcp') {
      return this.mcpClient.getTools();
    } else if (this.currentMode === 'api') {
      // API doesn't have explicit tool discovery, return known tools
      return [
        { name: 'gorev_listele', description: 'List tasks' },
        { name: 'proje_listele', description: 'List projects' },
        { name: 'template_listele', description: 'List templates' },
        { name: 'templateden_gorev_olustur', description: 'Create task from template' },
        { name: 'gorev_guncelle', description: 'Update task' },
        { name: 'gorev_sil', description: 'Delete task' },
        { name: 'proje_olustur', description: 'Create project' },
        { name: 'aktif_proje_ayarla', description: 'Activate project' },
        { name: 'ozet_goster', description: 'Show summary' },
      ];
    }
    return [];
  }

  // Configuration methods
  setMode(mode: ClientMode): void {
    this.config.mode = mode;
  }

  getMode(): ClientMode {
    return this.config.mode;
  }

  setApiBaseURL(url: string): void {
    this.config.apiBaseURL = url;
    this.apiClient = new ApiClient(url);
    this.setupApiEventHandlers();
  }

  getApiBaseURL(): string {
    return this.config.apiBaseURL || 'http://localhost:5082';
  }

  // Health check for both modes
  async healthCheck(): Promise<{ mode: string; healthy: boolean; details?: any }> {
    if (!this.connected) {
      return { mode: 'none', healthy: false, details: 'Not connected' };
    }

    try {
      if (this.currentMode === 'api') {
        const health = await this.apiClient.checkHealth();
        return { mode: 'api', healthy: true, details: health };
      } else if (this.currentMode === 'mcp') {
        // MCP doesn't have explicit health check, but we can check if it's connected
        return { mode: 'mcp', healthy: this.mcpClient.isConnected() };
      }
    } catch (error) {
      return { mode: this.currentMode || 'unknown', healthy: false, details: error };
    }

    return { mode: 'unknown', healthy: false };
  }

  // Get status information
  getStatus(): {
    connected: boolean;
    mode: string | null;
    config: UnifiedClientConfig;
    autoDetectionInProgress: boolean;
  } {
    return {
      connected: this.connected,
      mode: this.currentMode,
      config: this.config,
      autoDetectionInProgress: this.autoDetectionInProgress,
    };
  }
}