import { spawn, ChildProcess } from 'child_process';
import { EventEmitter } from 'events';
import { Logger } from '../utils/logger';
import {
  MCPRequest,
  MCPResponse,
  MCPNotification,
  MCPInitializeParams,
  MCPToolCallParams,
  MCPToolResult,
  MCPTool,
} from './types';

export class MCPClient extends EventEmitter {
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

  async connect(serverPath: string): Promise<void> {
    if (this.connected) {
      throw new Error('Already connected to MCP server');
    }

    Logger.info(`Connecting to MCP server at: ${serverPath}`);

    try {
      // Set working directory to server's directory
      const path = require('path');
      
      // Spawn the MCP server process
      const spawnOptions: any = {
        stdio: ['pipe', 'pipe', 'pipe'],
        shell: false,
        env: { 
          ...process.env,
          // Set GOREV_ROOT to a data directory
          GOREV_ROOT: path.join(path.dirname(serverPath), '..', 'data')
        },
      };
      const serverDir = path.dirname(serverPath);
      if (serverDir) {
        spawnOptions.cwd = serverDir;
        Logger.debug(`Setting working directory to: ${serverDir}`);
      }

      this.process = spawn(serverPath, ['serve'], spawnOptions);

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

    const result = await this.sendRequest('tools/call', {
      name,
      arguments: params || {},
    } as MCPToolCallParams);

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
    
    // Process complete messages from buffer
    const lines = this.buffer.split('\n');
    this.buffer = lines.pop() || '';

    for (const line of lines) {
      if (line.trim()) {
        try {
          const message = JSON.parse(line);
          Logger.debug('Parsed message:', message);
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

      const timeout = setTimeout(() => {
        this.pendingRequests.delete(id);
        Logger.error(`Request timeout for ${method} (id: ${id})`);
        reject(new Error(`Request timeout: ${method}`));
      }, 10000); // 10 second timeout

      this.pendingRequests.set(id, { resolve, reject, timeout });

      const message = JSON.stringify(request) + '\n';
      Logger.debug(`Sending request: ${method} (id: ${id})`, request);
      
      try {
        this.process?.stdin?.write(message, (error) => {
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