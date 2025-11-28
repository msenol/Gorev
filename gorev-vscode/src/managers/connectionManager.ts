import * as vscode from 'vscode';
import { Logger } from '../utils/logger';
import { t } from '../utils/l10n';

export enum ConnectionMode {
  AUTO = 'auto',
  LOCAL = 'local',
  DOCKER = 'docker',
  REMOTE = 'remote'
}

export interface ConnectionConfig {
  mode: ConnectionMode;
  apiHost: string;
  apiPort: number;
  serverPath?: string;
  dockerComposeFile?: string;
}

export class ConnectionManager {
  private static instance: ConnectionManager;
  private terminal?: vscode.Terminal;

  private constructor() {}

  static getInstance(): ConnectionManager {
    if (!ConnectionManager.instance) {
      ConnectionManager.instance = new ConnectionManager();
    }
    return ConnectionManager.instance;
  }

  async ensureConnection(config: ConnectionConfig): Promise<boolean> {
    const apiURL = `http://${config.apiHost}:${config.apiPort}`;
    
    Logger.info(`🔍 Checking daemon at ${apiURL}`);

    // 1. Try direct connection first
    try {
      const response = await fetch(`${apiURL}/api/health`, { method: 'GET' });
      if (response.ok) {
        Logger.info(`✅ Daemon already running at ${apiURL}`);
        return true;
      }
    } catch (e) {
      // Daemon not accessible, continue to start
    }

    // 2. Connection wizard
    return await this.handleConnectionMode(config);
  }

  private async handleConnectionMode(config: ConnectionConfig): Promise<boolean> {
    switch (config.mode) {
      case ConnectionMode.REMOTE:
        return await this.connectRemote(config);
      case ConnectionMode.DOCKER:
        return await this.startDocker(config);
      case ConnectionMode.LOCAL:
        return await this.startLocal(config);
      case ConnectionMode.AUTO:
      default:
        return await this.handleAuto(config);
    }
  }

  private async connectRemote(config: ConnectionConfig): Promise<boolean> {
    const apiURL = `http://${config.apiHost}:${config.apiPort}`;
    
    await vscode.window.showErrorMessage(
      t('connection.remoteNotConnected', { host: config.apiHost, port: config.apiPort }),
      t('wizard.learnMore')
    ).then(selection => {
      if (selection === t('wizard.learnMore')) {
        vscode.env.openExternal(vscode.Uri.parse('https://github.com/msenol/Gorev/blob/main/docs/remote-setup.md'));
      }
    });

    return false;
  }

  private async startDocker(config: ConnectionConfig): Promise<boolean> {
    const composeFile = config.dockerComposeFile || './docker-compose.yml';
    
    if (!config.dockerComposeFile || config.dockerComposeFile === './docker-compose.yml') {
      // Create default docker-compose.yml if it doesn't exist
      if (!await this.fileExists(composeFile)) {
        await this.createDefaultDockerCompose(composeFile);
      }
    }

    if (!await this.commandExists('docker-compose')) {
      await vscode.window.showErrorMessage(
        t('connection.dockerNotFound'),
        t('wizard.installDocker')
      );
      return false;
    }

    const terminal = this.getOrCreateTerminal('Gorev Docker');
    terminal.sendText(`docker-compose -f "${composeFile}" up -d`);
    terminal.show();

    await vscode.window.showInformationMessage(
      t('connection.dockerStarting'),
      t('wizard.wait')
    );

    // Wait for daemon
    return await this.waitForDaemon(config, 30);
  }

  private async startLocal(config: ConnectionConfig): Promise<boolean> {
    let serverPath = config.serverPath;
    
    if (!serverPath) {
      // Try to find in PATH
      if (await this.commandExists('gorev')) {
        serverPath = 'gorev';
      } else {
        const install = await vscode.window.showErrorMessage(
          t('connection.localNotFound'),
          t('wizard.installGorev'),
          t('wizard.cancel')
        );
        
        if (install === t('wizard.installGorev')) {
          await this.showInstallationInstructions();
        }
        return false;
      }
    }

    if (!await this.fileExists(serverPath) && serverPath !== 'gorev') {
      await vscode.window.showErrorMessage(t('connection.serverPathNotFound', { path: serverPath }));
      return false;
    }

    const terminal = this.getOrCreateTerminal('Gorev Local');
    terminal.sendText(`"${serverPath}" daemon --detach --port ${config.apiPort}`);
    terminal.show();

    await vscode.window.showInformationMessage(
      t('connection.localStarting'),
      t('wizard.wait')
    );

    return await this.waitForDaemon(config, 15);
  }

  private async handleAuto(config: ConnectionConfig): Promise<boolean> {
    // 1. Try remote if specified
    if (config.apiHost !== 'localhost') {
      return await this.connectRemote(config);
    }

    // 2. Try Docker if docker-compose.yml exists
    if (await this.fileExists('./docker-compose.yml') && await this.commandExists('docker-compose')) {
      Logger.info("📄 Found docker-compose.yml, trying Docker mode...");
      if (await this.startDocker(config)) {
        return true;
      }
    }

    // 3. Try local
    return await this.startLocal(config);
  }

  private async waitForDaemon(config: ConnectionConfig, timeoutSeconds: number): Promise<boolean> {
    const apiURL = `http://${config.apiHost}:${config.apiPort}`;
    const startTime = Date.now();
    const timeout = timeoutSeconds * 1000;

    while (Date.now() - startTime < timeout) {
      try {
        const response = await fetch(`${apiURL}/api/health`);
        if (response.ok) {
          await vscode.window.showInformationMessage(t('connection.success'));
          return true;
        }
      } catch (e) {
        // Continue waiting
      }
      await new Promise(resolve => setTimeout(resolve, 1000));
    }

    await vscode.window.showErrorMessage(t('connection.timeout'));
    return false;
  }

  // Helper methods
  private getOrCreateTerminal(name: string): vscode.Terminal {
    if (!this.terminal || this.terminal.exitStatus) {
      this.terminal = vscode.window.createTerminal(name);
    }
    return this.terminal;
  }

  private async fileExists(path: string): Promise<boolean> {
    try {
      await vscode.workspace.fs.stat(vscode.Uri.file(path));
      return true;
    } catch {
      return false;
    }
  }

  private async commandExists(cmd: string): Promise<boolean> {
    try {
      const { exec } = await import('child_process');
      await new Promise((resolve, reject) => {
        exec(`which ${cmd}`, (error) => {
          if (error) reject(error);
          else resolve(true);
        });
      });
      return true;
    } catch {
      return false;
    }
  }

  private async showInstallationInstructions(): Promise<void> {
    const action = await vscode.window.showInformationMessage(
      t('wizard.installInstructions'),
      t('wizard.openNpm'),
      t('wizard.copyCommand')
    );

    if (action === t('wizard.openNpm')) {
      vscode.env.openExternal(vscode.Uri.parse('https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server'));
    } else if (action === t('wizard.copyCommand')) {
      await vscode.env.clipboard.writeText('npm install -g @mehmetsenol/gorev-mcp-server');
      vscode.window.showInformationMessage(t('wizard.commandCopied'));
    }
  }

  dispose(): void {
    if (this.terminal) {
      this.terminal.dispose();
    }
  }

  private async createDefaultDockerCompose(path: string): Promise<void> {
    const content = `version: '3.8'
services:
  gorev:
    image: msenol/gorev:latest
    ports:
      - "5082:5082"
    volumes:
      - ./workspace:/workspace
    environment:
      - GOREV_LANG=tr
    command: daemon --detach
`;
    await vscode.workspace.fs.writeFile(vscode.Uri.file(path), Buffer.from(content));
  }
}
