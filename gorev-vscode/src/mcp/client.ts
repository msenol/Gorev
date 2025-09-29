import { spawn, ChildProcess } from 'child_process';
import { EventEmitter } from 'events';
import { Logger } from '../utils/logger';
import { ClientInterface } from '../interfaces/client';
import {
  MCPRequest,
  MCPResponse,
  MCPNotification,
  MCPInitializeParams,
  MCPToolCallParams,
  MCPToolResult,
  MCPTool,
} from './types';
import * as vscode from 'vscode';

export class MCPClient extends EventEmitter implements ClientInterface {
  private process: ChildProcess | null = null;
  private requestId = 0;
  private pendingRequests = new Map<number, {
    resolve: (value: any) => void;
    reject: (error: any) => void;
    timeout?: NodeJS.Timeout;
  }>();
  private connected = false;
  private buffer = '';
  private tools = new Map<string, MCPTool>();

  constructor() {
    super();
  }

  async connect(serverPath?: string): Promise<void> {
    if (this.connected) {
      throw new Error('Already connected to MCP server');
    }

    // Get server mode from configuration
    const config = vscode.workspace.getConfiguration('gorev');
    const serverMode = config.get<string>('serverMode', 'npx');

    let command: string;
    let args: string[];

    if (serverMode === 'npx') {
      // Windows needs cmd wrapper, others can use npx directly
      if (process.platform === 'win32') {
        command = 'cmd';
        args = ['/c', 'npx', '-y', '@mehmetsenol/gorev-mcp-server@latest'];
        Logger.info('Connecting to MCP server via NPX (Windows): cmd /c npx -y @mehmetsenol/gorev-mcp-server@latest');
      } else {
        command = 'npx';
        args = ['-y', '@mehmetsenol/gorev-mcp-server@latest'];
        Logger.info('Connecting to MCP server via NPX: npx -y @mehmetsenol/gorev-mcp-server@latest');
      }
    } else {
      // Binary mode
      if (!serverPath) {
        throw new Error('Server path is required for binary mode');
      }
      command = serverPath;
      args = []; // Remove 'serve' argument - binary runs in MCP mode by default
      Logger.info(`Connecting to MCP server at: ${serverPath}`);
    }

    try {
      // Set working directory to server's directory
      const path = require('path');

      // Get database mode from configuration
      const databaseMode = vscode.workspace.getConfiguration('gorev').get<string>('databaseMode', 'auto');

      // Prepare environment variables
      const env: any = {
        ...process.env,
        // Set GOREV_ROOT to a data directory (fallback)
        GOREV_ROOT: serverMode === 'npx' ? path.join(require('os').homedir(), '.gorev') : path.join(path.dirname(serverPath || ''), '..', 'data')
      };

      // Determine database path based on mode and workspace
      const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
      let databasePath: string | null = null;

      if (databaseMode === 'global') {
        // Force global mode - don't set GOREV_DB_PATH, let server use default
        Logger.info('Using global database mode');
      } else if (databaseMode === 'workspace' && workspaceFolder) {
        // Force workspace mode
        databasePath = path.join(workspaceFolder.uri.fsPath, '.gorev', 'gorev.db');
        Logger.info(`Using workspace database mode: ${databasePath}`);
      } else if (databaseMode === 'auto' && workspaceFolder) {
        // Auto mode - check if .gorev directory exists
        const workspaceDbPath = path.join(workspaceFolder.uri.fsPath, '.gorev', 'gorev.db');
        const workspaceDbDir = path.join(workspaceFolder.uri.fsPath, '.gorev');

        try {
          const fs = require('fs');
          if (fs.existsSync(workspaceDbDir)) {
            databasePath = workspaceDbPath;
            Logger.info(`Auto-detected workspace database: ${workspaceDbPath}`);
          } else {
            Logger.info('Auto mode: No .gorev directory found, using global database');
          }
        } catch (error) {
          Logger.warn(`Failed to check workspace database directory: ${error}`);
          Logger.info('Auto mode: Falling back to global database');
        }
      } else {
        Logger.info('Using global database mode (no workspace or global mode selected)');
      }

      // Set GOREV_DB_PATH if we determined a specific path
      if (databasePath) {
        env.GOREV_DB_PATH = databasePath;
        Logger.info(`Setting GOREV_DB_PATH: ${databasePath}`);
      }

      // Store database mode info for status bar
      this.emit('databaseModeChanged', {
        mode: databaseMode,
        path: databasePath,
        workspaceFolder: workspaceFolder?.uri.fsPath
      });

      // Spawn the MCP server process
      const spawnOptions: any = {
        stdio: ['pipe', 'pipe', 'pipe'],
        shell: process.platform === 'win32' && serverMode === 'npx', // Use shell on Windows for cmd wrapper
        env: env,
      };
      // Set working directory for NPX or binary mode
      if (serverMode === 'binary' && serverPath) {
        const serverDir = path.dirname(serverPath);
        if (serverDir) {
          spawnOptions.cwd = serverDir;
          Logger.debug(`Setting working directory to: ${serverDir}`);
        }
      }

      Logger.debug(`Spawning process: ${command} ${args.join(' ')}`);
      Logger.debug(`Spawn options: ${JSON.stringify(spawnOptions)}`);

      this.process = spawn(command, args, spawnOptions);
      
      if (!this.process) {
        throw new Error('Failed to spawn process');
      }
      
      Logger.debug(`Process spawned with PID: ${this.process.pid}`);

      this.setupProcessHandlers();

      // Wait for process to be ready
      await new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
          reject(new Error('Server startup timeout'));
        }, 5000);

        const checkReady = () => {
          if (this.process && this.process.stdin && this.process.stdout) {
            clearTimeout(timeout);
            resolve(true);
          } else {
            setTimeout(checkReady, 100);
          }
        };
        checkReady();
      });

      // Initialize the connection
      await this.initialize();
      
      // Discover available tools
      await this.discoverTools();

      this.connected = true;
      this.emit('connected');
      
      Logger.info('Successfully connected to MCP server');
    } catch (error) {
      Logger.error('Connection failed:', error);
      this.cleanup();
      throw error;
    }
  }

  disconnect(): void {
    Logger.info('Disconnecting from MCP server');
    this.cleanup();
    this.emit('disconnected');
  }

  isConnected(): boolean {
    return this.connected;
  }

  async callTool(name: string, params?: any): Promise<MCPToolResult> {
    if (!this.connected) {
      throw new Error('Not connected to MCP server');
    }

    // DEBUG: Log tool call for debugging
    Logger.info(`[MCPClient] Calling tool: ${name} with params:`, params);

    const result = await this.sendRequest('tools/call', {
      name,
      arguments: params || {},
    } as MCPToolCallParams);

    // DEBUG: Log response for gorev_listele specifically
    if (name === 'gorev_listele') {
      Logger.info(`[MCPClient] gorev_listele RESPONSE:`);
      Logger.info(`[MCPClient] - Response type: ${typeof result}`);
      Logger.info(`[MCPClient] - Response keys:`, Object.keys(result || {}));
      if (result && result.content && Array.isArray(result.content)) {
        Logger.info(`[MCPClient] - Content array length: ${result.content.length}`);
        result.content.forEach((content: any, idx: number) => {
          if (content.type === 'text') {
            Logger.info(`[MCPClient] - Content ${idx} text length: ${content.text?.length}`);
            Logger.info(`[MCPClient] - Content ${idx} first 500 chars:`, content.text?.substring(0, 500));
          }
        });
      }
    }

    return result as MCPToolResult;
  }

  getTools(): MCPTool[] {
    return Array.from(this.tools.values());
  }

  private setupProcessHandlers(): void {
    if (!this.process) return;

    this.process.stdout?.on('data', (data: Buffer) => {
      const str = data.toString();
      Logger.debug('Received data:', str);
      this.handleData(str);
    });

    this.process.stderr?.on('data', (data: Buffer) => {
      const str = data.toString();
      Logger.error('MCP Server Error:', str);
      // Don't treat stderr as fatal - some servers use it for logging
    });

    this.process.on('error', (error) => {
      Logger.error('MCP Process Error:', error);
      Logger.error('Error details:', JSON.stringify(error));
      this.handleProcessError(error);
    });

    this.process.on('exit', (code, signal) => {
      Logger.info(`MCP Process exited with code ${code}, signal ${signal}`);
      this.handleProcessExit(code, signal);
    });

    this.process.on('close', (code, signal) => {
      Logger.info(`MCP Process closed with code ${code}, signal ${signal}`);
      this.handleProcessExit(code, signal);
    });
  }

  private handleData(data: string): void {
    this.buffer += data;
    Logger.debug(`Buffer size: ${this.buffer.length}, Data chunk size: ${data.length}`);
    
    // Process complete messages from buffer
    const lines = this.buffer.split('\n');
    this.buffer = lines.pop() || '';

    for (const line of lines) {
      if (line.trim()) {
        try {
          const message = JSON.parse(line);
          Logger.debug(`Parsed message (${line.length} chars):`, message);
          this.handleMessage(message);
        } catch (error) {
          Logger.error('Failed to parse MCP message:', line, error);
        }
      }
    }
  }

  private handleMessage(message: MCPResponse | MCPNotification): void {
    if ('id' in message) {
      // Response to a request
      const pending = this.pendingRequests.get(message.id);
      if (pending) {
        if (pending.timeout) {
          clearTimeout(pending.timeout);
        }
        this.pendingRequests.delete(message.id);

        if (message.error) {
          pending.reject(new Error(message.error.message));
        } else {
          pending.resolve(message.result);
        }
      }
    } else {
      // Notification
      this.handleNotification(message as MCPNotification);
    }
  }

  private handleNotification(notification: MCPNotification): void {
    Logger.debug('Received notification:', notification.method);
    this.emit('notification', notification);
  }

  private async initialize(): Promise<void> {
    const params: MCPInitializeParams = {
      protocolVersion: '2024-11-05',
      capabilities: {
        tools: {
          listChanged: true,
        },
        logging: {},
      },
      clientInfo: {
        name: 'gorev-vscode',
        version: '0.1.0',
      },
    };

    Logger.debug('Sending initialize request:', params);
    const result = await this.sendRequest('initialize', params);
    Logger.debug('Initialize response:', result);
    
    // Send initialized notification
    await this.sendNotification('initialized', {});
  }

  private async discoverTools(): Promise<void> {
    Logger.debug('Discovering tools...');
    const result = await this.sendRequest('tools/list', {});
    const tools = (result as any).tools || [];
    
    this.tools.clear();
    for (const tool of tools) {
      this.tools.set(tool.name, tool);
    }
    
    Logger.info(`Discovered ${this.tools.size} MCP tools`);
  }

  private sendRequest(method: string, params?: any): Promise<any> {
    return new Promise((resolve, reject) => {
      const id = ++this.requestId;
      const request: MCPRequest = {
        jsonrpc: '2.0',
        id,
        method,
        params,
      };

      // Get timeout from configuration
      const config = vscode.workspace.getConfiguration('gorev');
      const timeoutMs = config.get<number>('debug.serverTimeout', 5000);
      
      const timeout = setTimeout(() => {
        this.pendingRequests.delete(id);
        Logger.error(`Request timeout for ${method} (id: ${id}) after ${timeoutMs}ms`);
        reject(new Error(`Request timeout: ${method}`));
      }, timeoutMs);

      this.pendingRequests.set(id, { resolve, reject, timeout });

      const message = JSON.stringify(request) + '\n';
      Logger.debug(`Sending request: ${method} (id: ${id})`, request);
      
      try {
        if (!this.process) {
          Logger.error('Process is null when trying to send request');
          throw new Error('Process not available');
        }
        if (!this.process.stdin) {
          Logger.error('Process stdin is null');
          throw new Error('Process stdin not available');
        }
        
        this.process.stdin.write(message, (error) => {
          if (error) {
            Logger.error('Failed to write to stdin:', error);
            this.pendingRequests.delete(id);
            clearTimeout(timeout);
            reject(error);
          }
        });
      } catch (error) {
        Logger.error('Error sending request:', error);
        this.pendingRequests.delete(id);
        clearTimeout(timeout);
        reject(error);
      }
    });
  }

  private async sendNotification(method: string, params?: any): Promise<void> {
    const notification: MCPNotification = {
      jsonrpc: '2.0',
      method,
      params,
    };

    const message = JSON.stringify(notification) + '\n';
    Logger.debug(`Sending notification: ${method}`, notification);
    
    return new Promise((resolve, reject) => {
      this.process?.stdin?.write(message, (error) => {
        if (error) {
          Logger.error('Failed to send notification:', error);
          reject(error);
        } else {
          resolve();
        }
      });
    });
  }

  private handleProcessError(error: Error): void {
    Logger.error('MCP process error:', error);
    this.cleanup();
    this.emit('error', error);
  }

  private handleProcessExit(code: number | null, signal: string | null): void {
    if (this.connected) {
      Logger.warn(`MCP process exited unexpectedly (code: ${code}, signal: ${signal})`);
      this.cleanup();
      this.emit('disconnected', { unexpected: true, code, signal });
    }
  }

  private cleanup(): void {
    // Clear pending requests
    for (const [id, pending] of this.pendingRequests) {
      if (pending.timeout) {
        clearTimeout(pending.timeout);
      }
      pending.reject(new Error('Connection closed'));
    }
    this.pendingRequests.clear();

    // Kill process if still running
    if (this.process && !this.process.killed) {
      this.process.kill();
    }

    this.process = null;
    this.connected = false;
    this.buffer = '';
    this.tools.clear();
  }
}