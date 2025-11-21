import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import { COMMANDS } from '../utils/constants';
import { t } from '../utils/l10n';
import { WorkspaceContext } from '../models/workspace';

export class StatusBarManager implements vscode.Disposable {
  private statusBarItem: vscode.StatusBarItem;
  private databaseModeItem: vscode.StatusBarItem;
  private workspaceItem: vscode.StatusBarItem;
  private visible = false;
  private currentDatabaseMode = 'auto';
  private currentDatabasePath = '';
  private workspaceContext: WorkspaceContext | undefined;

  constructor() {
    this.statusBarItem = vscode.window.createStatusBarItem(
      vscode.StatusBarAlignment.Left,
      100
    );
    this.statusBarItem.command = COMMANDS.SHOW_SUMMARY;

    // Database mode indicator
    this.databaseModeItem = vscode.window.createStatusBarItem(
      vscode.StatusBarAlignment.Left,
      99
    );
    this.updateDatabaseModeDisplay();

    // Workspace indicator
    this.workspaceItem = vscode.window.createStatusBarItem(
      vscode.StatusBarAlignment.Left,
      98
    );
    this.updateWorkspaceDisplay();
  }

  show(): void {
    this.visible = true;
    this.statusBarItem.show();
    this.databaseModeItem.show();
    this.workspaceItem.show();
  }

  hide(): void {
    this.visible = false;
    this.statusBarItem.hide();
    this.databaseModeItem.hide();
    this.workspaceItem.hide();
  }

  isVisible(): boolean {
    return this.visible;
  }

  setConnecting(): void {
    this.statusBarItem.text = t('statusBar.connecting');
    this.statusBarItem.tooltip = t('statusBar.connectingTooltip');
    this.statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
  }

  setConnected(): void {
    this.statusBarItem.text = t('statusBar.connected');
    this.statusBarItem.tooltip = t('statusBar.connectedTooltip');
    this.statusBarItem.backgroundColor = undefined;
  }

  setDisconnected(): void {
    this.statusBarItem.text = t('statusBar.disconnected');
    this.statusBarItem.tooltip = t('statusBar.disconnectedTooltip');
    this.statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.errorBackground');
  }

  setConnectionStatus(connected: boolean, mode?: string): void {
    if (connected) {
      this.setConnected();
      if (mode) {
        // Update tooltip to show connection mode
        this.statusBarItem.tooltip = `${t('statusBar.connectedTooltip')} (${mode.toUpperCase()})`;
      }
    } else {
      this.setDisconnected();
    }
  }

  update(taskCount?: number, activeProject?: string): void {
    if (taskCount !== undefined) {
      const projectText = activeProject ? t('statusBar.inProject', activeProject) : '';
      this.statusBarItem.text = t('statusBar.taskCount', taskCount.toString(), projectText);
      this.statusBarItem.tooltip = t('statusBar.taskCountTooltip', taskCount.toString(), projectText);
    }
  }

  setDatabaseMode(mode: string, databasePath?: string): void {
    this.currentDatabaseMode = mode;
    this.currentDatabasePath = databasePath || '';
    this.updateDatabaseModeDisplay();
  }

  private updateDatabaseModeDisplay(): void {
    const mode = this.currentDatabaseMode;
    let icon = '';
    let text = '';

    if (mode === 'workspace') {
      icon = '📁';
      text = 'Workspace';
    } else if (mode === 'global') {
      icon = '🌐';
      text = 'Global';
    } else {
      // auto mode
      const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
      if (workspaceFolder) {
        const workspaceDbDir = path.join(workspaceFolder.uri.fsPath, '.gorev');

        if (fs.existsSync(workspaceDbDir)) {
          icon = '📁';
          text = 'Auto (Workspace)';
        } else {
          icon = '🌐';
          text = 'Auto (Global)';
        }
      } else {
        icon = '🌐';
        text = 'Global';
      }
    }

    this.databaseModeItem.text = `${icon} ${text}`;

    let tooltip = `Database Mode: ${text}`;
    if (this.currentDatabasePath) {
      tooltip += `\nPath: ${this.currentDatabasePath}`;
    }
    this.databaseModeItem.tooltip = tooltip;
  }

  setWorkspaceContext(context: WorkspaceContext | undefined): void {
    this.workspaceContext = context;
    this.updateWorkspaceDisplay();
  }

  private updateWorkspaceDisplay(): void {
    if (!this.workspaceContext) {
      this.workspaceItem.text = '$(folder) No Workspace';
      this.workspaceItem.tooltip = 'No workspace registered';
      this.workspaceItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
    } else {
      this.workspaceItem.text = `$(folder) ${this.workspaceContext.workspaceName}`;
      this.workspaceItem.tooltip = `Workspace: ${this.workspaceContext.workspaceName}\nID: ${this.workspaceContext.workspaceId}\nPath: ${this.workspaceContext.workspacePath}`;
      this.workspaceItem.backgroundColor = undefined;
    }
  }

  dispose(): void {
    this.statusBarItem.dispose();
    this.databaseModeItem.dispose();
    this.workspaceItem.dispose();
  }
}
