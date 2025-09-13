import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { COMMANDS } from '../utils/constants';
import { t } from '../utils/l10n';

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

  update(taskCount?: number, activeProject?: string): void {
    if (taskCount !== undefined) {
      const projectText = activeProject ? t('statusBar.inProject', activeProject) : '';
      this.statusBarItem.text = t('statusBar.taskCount', taskCount.toString(), projectText);
      this.statusBarItem.tooltip = t('statusBar.taskCountTooltip', taskCount.toString(), projectText);
    }
  }

  dispose(): void {
    this.statusBarItem.dispose();
  }
}
