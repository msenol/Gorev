import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { GorevDurum, GorevOncelik } from '../models/common';
import { GorevTreeItem } from '../providers/gorevTreeProvider';

export function registerGorevCommands(
  context: vscode.ExtensionContext,
  mcpClient: MCPClient,
  providers: CommandContext
): void {
  // Create Task
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_TASK, async () => {
      try {
        const baslik = await vscode.window.showInputBox({
          prompt: 'Task title',
          placeHolder: 'Enter task title',
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return 'Task title is required';
            }
            return null;
          },
        });

        if (!baslik) return;

        const aciklama = await vscode.window.showInputBox({
          prompt: 'Task description (optional)',
          placeHolder: 'Enter task description',
        });

        const oncelik = await vscode.window.showQuickPick(
          [
            { label: 'High', value: GorevOncelik.Yuksek },
            { label: 'Medium', value: GorevOncelik.Orta },
            { label: 'Low', value: GorevOncelik.Dusuk },
          ],
          {
            placeHolder: 'Select priority',
          }
        );

        if (!oncelik) return;

        await mcpClient.callTool('gorev_olustur', {
          baslik,
          aciklama: aciklama || '',
          oncelik: oncelik.value,
        });

        vscode.window.showInformationMessage('Task created successfully');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to create task: ${error}`);
      }
    })
  );

  // Quick Create Task
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.QUICK_CREATE_TASK, async () => {
      try {
        const input = await vscode.window.showInputBox({
          prompt: 'Quick create task',
          placeHolder: 'e.g., Fix login bug (high priority)',
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return 'Task description is required';
            }
            return null;
          },
        });

        if (!input) return;

        // Parse priority from input
        let oncelik = GorevOncelik.Orta;
        let baslik = input;

        if (input.toLowerCase().includes('high') || input.toLowerCase().includes('urgent')) {
          oncelik = GorevOncelik.Yuksek;
        } else if (input.toLowerCase().includes('low')) {
          oncelik = GorevOncelik.Dusuk;
        }

        // Remove priority indicators from title
        baslik = baslik
          .replace(/\s*\(?(high|medium|low|urgent|priority)\)?/gi, '')
          .trim();

        await mcpClient.callTool('gorev_olustur', {
          baslik,
          oncelik,
        });

        vscode.window.showInformationMessage('Task created successfully');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to create task: ${error}`);
      }
    })
  );

  // Refresh Tasks
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.REFRESH_TASKS, async () => {
      try {
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to refresh tasks: ${error}`);
      }
    })
  );

  // Show Task Detail
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SHOW_TASK_DETAIL, async (item: GorevTreeItem) => {
      try {
        const result = await mcpClient.callTool('gorev_detay', {
          id: item.task.id,
        });

        const panel = vscode.window.createWebviewPanel(
          'gorevDetail',
          `Task: ${item.task.baslik}`,
          vscode.ViewColumn.Two,
          {
            enableScripts: true,
          }
        );

        panel.webview.html = getTaskDetailHtml(result.content[0].text);
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to show task details: ${error}`);
      }
    })
  );

  // Update Task Status
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.UPDATE_TASK_STATUS, async (item: GorevTreeItem) => {
      try {
        const newStatus = await vscode.window.showQuickPick(
          [
            { label: 'Pending', value: GorevDurum.Beklemede },
            { label: 'In Progress', value: GorevDurum.DevamEdiyor },
            { label: 'Completed', value: GorevDurum.Tamamlandi },
          ],
          {
            placeHolder: 'Select new status',
          }
        );

        if (!newStatus) return;

        await mcpClient.callTool('gorev_guncelle', {
          id: item.task.id,
          durum: newStatus.value,
        });

        vscode.window.showInformationMessage('Task status updated');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to update task status: ${error}`);
      }
    })
  );

  // Delete Task
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DELETE_TASK, async (item: GorevTreeItem) => {
      try {
        const confirm = await vscode.window.showWarningMessage(
          `Are you sure you want to delete "${item.task.baslik}"?`,
          'Yes',
          'No'
        );

        if (confirm !== 'Yes') return;

        await mcpClient.callTool('gorev_sil', {
          id: item.task.id,
          onay: true,
        });

        vscode.window.showInformationMessage('Task deleted');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to delete task: ${error}`);
      }
    })
  );
}

function getTaskDetailHtml(content: string): string {
  return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: var(--vscode-font-family);
            color: var(--vscode-foreground);
            background-color: var(--vscode-editor-background);
            padding: 20px;
            line-height: 1.6;
        }
        h1, h2, h3 {
            color: var(--vscode-foreground);
            margin-top: 20px;
        }
        .status {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: bold;
        }
        .status.pending { background-color: var(--vscode-badge-background); }
        .status.in-progress { background-color: var(--vscode-progressBar-background); }
        .status.completed { background-color: var(--vscode-testing-iconPassed); }
        .priority {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            margin-left: 8px;
        }
        .priority.high { color: var(--vscode-errorForeground); }
        .priority.medium { color: var(--vscode-editorWarning-foreground); }
        .priority.low { color: var(--vscode-editorInfo-foreground); }
        pre {
            background-color: var(--vscode-textBlockQuote-background);
            padding: 10px;
            border-radius: 4px;
            overflow-x: auto;
        }
        .dependency {
            padding: 8px;
            margin: 4px 0;
            border-left: 3px solid var(--vscode-badge-background);
            background-color: var(--vscode-editor-lineHighlightBackground);
        }
    </style>
</head>
<body>
    ${content.replace(/\n/g, '<br>')}
</body>
</html>`;
}