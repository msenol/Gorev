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
    this.statusBarItem.text = '$(sync~spin) Gorev: Connecting...';
    this.statusBarItem.tooltip = 'Connecting to Gorev server...';
    this.statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
  }

  setConnected(): void {
    this.statusBarItem.text = '$(check) Gorev: Connected';
    this.statusBarItem.tooltip = 'Connected to Gorev server\nClick to show summary';
    this.statusBarItem.backgroundColor = undefined;
  }

  setDisconnected(): void {
    this.statusBarItem.text = '$(x) Gorev: Disconnected';
    this.statusBarItem.tooltip = 'Not connected to Gorev server';
    this.statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.errorBackground');
  }

  update(taskCount?: number, activeProject?: string): void {
    if (taskCount !== undefined) {
      const projectText = activeProject ? ` (${activeProject})` : '';
      this.statusBarItem.text = `$(checklist) Gorev: ${taskCount} tasks${projectText}`;
      this.statusBarItem.tooltip = `${taskCount} tasks in total${projectText}\nClick to show summary`;
    }
  }

  dispose(): void {
    this.statusBarItem.dispose();
  }
}