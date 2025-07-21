import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { COMMANDS } from '../utils/constants';

export class StatusBarManager implements vscode.Disposable {
  private statusBarItem: vscode.StatusBarItem;
  private visible = false;

  constructor() {
    this.statusBarItem = vscode.window.createStatusBarItem(
      vscode.StatusBarAlignment.Left,
      100
    );
    this.statusBarItem.command = COMMANDS.SHOW_SUMMARY;
  }

  show(): void {
    this.visible = true;
    this.statusBarItem.show();
  }

  hide(): void {
    this.visible = false;
    this.statusBarItem.hide();
  }

  isVisible(): boolean {
    return this.visible;
  }

  setConnecting(): void {
    this.statusBarItem.text = vscode.l10n.t('statusBar.connecting');
    this.statusBarItem.tooltip = vscode.l10n.t('statusBar.connectingTooltip');
    this.statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
  }

  setConnected(): void {
    this.statusBarItem.text = vscode.l10n.t('statusBar.connected');
    this.statusBarItem.tooltip = vscode.l10n.t('statusBar.connectedTooltip');
    this.statusBarItem.backgroundColor = undefined;
  }

  setDisconnected(): void {
    this.statusBarItem.text = vscode.l10n.t('statusBar.disconnected');
    this.statusBarItem.tooltip = vscode.l10n.t('statusBar.disconnectedTooltip');
    this.statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.errorBackground');
  }

  update(taskCount?: number, activeProject?: string): void {
    if (taskCount !== undefined) {
      const projectText = activeProject ? vscode.l10n.t('statusBar.inProject', activeProject) : '';
      this.statusBarItem.text = vscode.l10n.t('statusBar.taskCount', taskCount.toString(), projectText);
      this.statusBarItem.tooltip = vscode.l10n.t('statusBar.taskCountTooltip', taskCount.toString(), projectText);
    }
  }

  dispose(): void {
    this.statusBarItem.dispose();
  }
}