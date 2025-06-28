import * as vscode from 'vscode';

export class Config {
  private static context: vscode.ExtensionContext;
  private static readonly CONFIG_SECTION = 'gorev';

  static initialize(context: vscode.ExtensionContext): void {
    this.context = context;
  }

  static get<T>(key: string): T | undefined {
    return vscode.workspace.getConfiguration(this.CONFIG_SECTION).get<T>(key);
  }

  static async update(key: string, value: any, global = true): Promise<void> {
    await vscode.workspace
      .getConfiguration(this.CONFIG_SECTION)
      .update(key, value, global);
  }

  static getGlobalState<T>(key: string, defaultValue: T): T {
    return this.context.globalState.get(key, defaultValue);
  }

  static async updateGlobalState(key: string, value: any): Promise<void> {
    await this.context.globalState.update(key, value);
  }

  static getWorkspaceState<T>(key: string, defaultValue: T): T {
    return this.context.workspaceState.get(key, defaultValue);
  }

  static async updateWorkspaceState(key: string, value: any): Promise<void> {
    await this.context.workspaceState.update(key, value);
  }
}