/**
 * Server Manager for E2E Tests
 *
 * Manages the real Gorev server for integration testing.
 * Can start, stop, and health check the server.
 */

import { spawn, ChildProcess, exec } from 'child_process';
import { promisify } from 'util';
import * as path from 'path';
import * as fs from 'fs';
import * as http from 'http';

const execAsync = promisify(exec);

export interface ServerConfig {
  port: number;
  host: string;
  workspacePath: string;
  language: string;
  binaryPath?: string;
  timeout: number;
}

export const defaultServerConfig: ServerConfig = {
  port: 5082,
  host: 'localhost',
  workspacePath: path.join(process.cwd(), 'test-workspace'),
  language: 'tr',
  timeout: 30000,
};

export class ServerManager {
  private config: ServerConfig;
  private serverProcess: ChildProcess | null = null;
  private isRunning = false;

  constructor(config: Partial<ServerConfig> = {}) {
    this.config = { ...defaultServerConfig, ...config };
  }

  /**
   * Find the gorev binary
   */
  private async findBinary(): Promise<string> {
    // Check if binary path is explicitly set
    if (this.config.binaryPath && fs.existsSync(this.config.binaryPath)) {
      return this.config.binaryPath;
    }

    // Check common locations
    const possiblePaths = [
      path.join(process.cwd(), '..', 'gorev-mcpserver', 'gorev'),
      path.join(process.cwd(), '..', 'gorev-mcpserver', 'binaries', 'gorev-linux'),
      path.join(process.cwd(), '..', 'gorev-mcpserver', 'binaries', 'gorev-darwin'),
      '/usr/local/bin/gorev',
    ];

    for (const binPath of possiblePaths) {
      if (fs.existsSync(binPath)) {
        return binPath;
      }
    }

    // Try to build if not found
    console.log('[ServerManager] Binary not found, attempting to build...');
    try {
      const mcpServerPath = path.join(process.cwd(), '..', 'gorev-mcpserver');
      await execAsync('make build', { cwd: mcpServerPath });
      const builtPath = path.join(mcpServerPath, 'gorev');
      if (fs.existsSync(builtPath)) {
        return builtPath;
      }
    } catch (error) {
      console.error('[ServerManager] Build failed:', error);
    }

    throw new Error('Gorev binary not found. Please build the server first with "make build" in gorev-mcpserver directory.');
  }

  /**
   * Ensure test workspace exists
   */
  private ensureWorkspace(): void {
    if (!fs.existsSync(this.config.workspacePath)) {
      fs.mkdirSync(this.config.workspacePath, { recursive: true });
    }
  }

  /**
   * Check if server is healthy
   */
  async healthCheck(): Promise<boolean> {
    return new Promise((resolve) => {
      const req = http.request(
        {
          hostname: this.config.host,
          port: this.config.port,
          path: '/api/v1/health',
          method: 'GET',
          timeout: 5000,
        },
        (res) => {
          resolve(res.statusCode === 200);
        }
      );

      req.on('error', () => resolve(false));
      req.on('timeout', () => {
        req.destroy();
        resolve(false);
      });

      req.end();
    });
  }

  /**
   * Wait for server to be ready
   */
  private async waitForReady(maxAttempts = 30, delayMs = 1000): Promise<void> {
    for (let i = 0; i < maxAttempts; i++) {
      if (await this.healthCheck()) {
        console.log(`[ServerManager] Server ready after ${i + 1} attempts`);
        return;
      }
      await new Promise((resolve) => setTimeout(resolve, delayMs));
    }
    throw new Error(`Server failed to start after ${maxAttempts} attempts`);
  }

  /**
   * Start the server
   */
  async start(): Promise<void> {
    if (this.isRunning) {
      console.log('[ServerManager] Server already running');
      return;
    }

    // Check if server is already running externally
    if (await this.healthCheck()) {
      console.log('[ServerManager] Using existing server');
      this.isRunning = true;
      return;
    }

    this.ensureWorkspace();
    const binaryPath = await this.findBinary();

    console.log(`[ServerManager] Starting server with binary: ${binaryPath}`);
    console.log(`[ServerManager] Workspace: ${this.config.workspacePath}`);

    this.serverProcess = spawn(binaryPath, ['serve', '--port', String(this.config.port), '--lang', this.config.language], {
      cwd: this.config.workspacePath,
      env: {
        ...process.env,
        GOREV_LANG: this.config.language,
      },
      stdio: ['pipe', 'pipe', 'pipe'],
    });

    this.serverProcess.stdout?.on('data', (data) => {
      console.log(`[Server] ${data.toString().trim()}`);
    });

    this.serverProcess.stderr?.on('data', (data) => {
      console.error(`[Server Error] ${data.toString().trim()}`);
    });

    this.serverProcess.on('error', (error) => {
      console.error('[ServerManager] Process error:', error);
      this.isRunning = false;
    });

    this.serverProcess.on('exit', (code) => {
      console.log(`[ServerManager] Process exited with code ${code}`);
      this.isRunning = false;
    });

    await this.waitForReady();
    this.isRunning = true;
    console.log('[ServerManager] Server started successfully');
  }

  /**
   * Stop the server
   */
  async stop(): Promise<void> {
    if (!this.serverProcess) {
      console.log('[ServerManager] No server process to stop');
      return;
    }

    console.log('[ServerManager] Stopping server...');

    return new Promise((resolve) => {
      this.serverProcess?.on('exit', () => {
        this.serverProcess = null;
        this.isRunning = false;
        console.log('[ServerManager] Server stopped');
        resolve();
      });

      this.serverProcess?.kill('SIGTERM');

      // Force kill after timeout
      setTimeout(() => {
        if (this.serverProcess) {
          console.log('[ServerManager] Force killing server...');
          this.serverProcess.kill('SIGKILL');
        }
      }, 5000);
    });
  }

  /**
   * Get the API base URL
   */
  getBaseUrl(): string {
    return `http://${this.config.host}:${this.config.port}`;
  }

  /**
   * Get server configuration
   */
  getConfig(): ServerConfig {
    return { ...this.config };
  }

  /**
   * Check if server is running
   */
  isServerRunning(): boolean {
    return this.isRunning;
  }
}

// Singleton instance for global server management
let globalServerManager: ServerManager | null = null;

export function getGlobalServerManager(config?: Partial<ServerConfig>): ServerManager {
  if (!globalServerManager) {
    globalServerManager = new ServerManager(config);
  }
  return globalServerManager;
}

export default ServerManager;
