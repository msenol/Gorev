import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { COMMANDS } from '../utils/constants';
import { Logger } from '../utils/logger';
import { CommandContext } from './index';

export function registerDatabaseCommands(
  context: vscode.ExtensionContext,
  mcpClient: MCPClient,
  providers: CommandContext
): void {

  // Initialize Workspace Database
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.INIT_WORKSPACE_DATABASE, async () => {
      try {
        const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
        if (!workspaceFolder) {
          vscode.window.showWarningMessage('No workspace folder found. Please open a folder first.');
          return;
        }

        const path = require('path');
        const fs = require('fs');

        const workspaceDbDir = path.join(workspaceFolder.uri.fsPath, '.gorev');
        const workspaceDbPath = path.join(workspaceDbDir, 'gorev.db');

        // Check if .gorev directory already exists
        if (fs.existsSync(workspaceDbDir)) {
          const answer = await vscode.window.showInformationMessage(
            'Workspace database directory already exists. Do you want to reinitialize it?',
            'Yes',
            'No'
          );
          if (answer !== 'Yes') {
            return;
          }
        }

        // Create .gorev directory
        fs.mkdirSync(workspaceDbDir, { recursive: true });

        vscode.window.showInformationMessage(
          `Workspace database initialized at: ${workspaceDbPath}\n\nPlease reconnect to server to use the new database.`,
          'Reconnect Now'
        ).then(selection => {
          if (selection === 'Reconnect Now') {
            vscode.commands.executeCommand(COMMANDS.DISCONNECT).then(() => {
              setTimeout(() => {
                vscode.commands.executeCommand(COMMANDS.CONNECT);
              }, 1000);
            });
          }
        });

        Logger.info(`Workspace database initialized at: ${workspaceDbPath}`);

      } catch (error) {
        Logger.error('Failed to initialize workspace database:', error);
        vscode.window.showErrorMessage(`Failed to initialize workspace database: ${error}`);
      }
    })
  );

  // Switch Database Mode
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SWITCH_DATABASE_MODE, async () => {
      try {
        const currentMode = vscode.workspace.getConfiguration('gorev').get<string>('databaseMode', 'auto');

        const items: vscode.QuickPickItem[] = [
          {
            label: '🔍 Auto',
            description: 'Auto-detect workspace database',
            detail: 'Use workspace database if .gorev/ exists, otherwise global',
            picked: currentMode === 'auto'
          },
          {
            label: '📁 Workspace',
            description: 'Always use project-local database',
            detail: 'Store database in current workspace (.gorev/gorev.db)',
            picked: currentMode === 'workspace'
          },
          {
            label: '🌐 Global',
            description: 'Always use shared database',
            detail: 'Store database in user home directory (~/.gorev/gorev.db)',
            picked: currentMode === 'global'
          }
        ];

        const selection = await vscode.window.showQuickPick(items, {
          title: 'Select Database Mode',
          placeHolder: `Current mode: ${currentMode}`,
          ignoreFocusOut: true
        });

        if (!selection) {
          return;
        }

        let newMode: string;
        if (selection.label.includes('Auto')) {
          newMode = 'auto';
        } else if (selection.label.includes('Workspace')) {
          newMode = 'workspace';
        } else {
          newMode = 'global';
        }

        if (newMode === currentMode) {
          vscode.window.showInformationMessage(`Database mode is already set to: ${newMode}`);
          return;
        }

        // Update configuration
        await vscode.workspace.getConfiguration('gorev').update('databaseMode', newMode, vscode.ConfigurationTarget.Workspace);

        vscode.window.showInformationMessage(
          `Database mode changed to: ${newMode}\n\nPlease reconnect to server to apply changes.`,
          'Reconnect Now'
        ).then(selection => {
          if (selection === 'Reconnect Now') {
            vscode.commands.executeCommand(COMMANDS.DISCONNECT).then(() => {
              setTimeout(() => {
                vscode.commands.executeCommand(COMMANDS.CONNECT);
              }, 1000);
            });
          }
        });

        Logger.info(`Database mode changed from ${currentMode} to ${newMode}`);

      } catch (error) {
        Logger.error('Failed to switch database mode:', error);
        vscode.window.showErrorMessage(`Failed to switch database mode: ${error}`);
      }
    })
  );
}