import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { GorevDurum, GorevOncelik } from '../models/common';
import { TaskDetailPanel } from '../ui/taskDetailPanel';
// import { GorevTreeItem } from '../providers/gorevTreeProvider';

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
    vscode.commands.registerCommand(COMMANDS.SHOW_TASK_DETAIL, async (item: any) => {
      try {
        // Use the new TaskDetailPanel
        await TaskDetailPanel.createOrShow(
          mcpClient,
          item.task,
          context.extensionUri
        );
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to show task details: ${error}`);
      }
    })
  );

  // Update Task Status
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.UPDATE_TASK_STATUS, async (item: any) => {
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
    vscode.commands.registerCommand(COMMANDS.DELETE_TASK, async (item: any) => {
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

  // Create Subtask
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_SUBTASK, async (item: any) => {
      try {
        const baslik = await vscode.window.showInputBox({
          prompt: 'Subtask title',
          placeHolder: `Create subtask for "${item.task.baslik}"`,
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return 'Subtask title is required';
            }
            return null;
          },
        });

        if (!baslik) return;

        const aciklama = await vscode.window.showInputBox({
          prompt: 'Subtask description (optional)',
          placeHolder: 'Enter subtask description',
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

        await mcpClient.callTool('gorev_altgorev_olustur', {
          parent_id: item.task.id,
          baslik,
          aciklama: aciklama || '',
          oncelik: oncelik.value,
        });

        vscode.window.showInformationMessage('Subtask created successfully');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to create subtask: ${error}`);
      }
    })
  );

  // Change Parent
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CHANGE_PARENT, async (item: any) => {
      try {
        // Get all tasks except the current one and its subtasks
        const result = await mcpClient.callTool('gorev_listele', {
          tum_projeler: true,
        });
        
        if (!result || !result.content || !result.content[0]) {
          throw new Error('Failed to fetch tasks');
        }

        const tasks = result.content[0].text
          .split('\n')
          .filter((line: string) => line.includes('ID:'))
          .map((line: string) => {
            const idMatch = line.match(/ID:\s*([a-f0-9-]+)/);
            const titleMatch = line.match(/\[.+\]\s+(.+)\s+\(/);
            return idMatch && titleMatch ? { id: idMatch[1], baslik: titleMatch[1] } : null;
          })
          .filter((task: any) => task && task.id !== item.task.id);

        const parentChoice = await vscode.window.showQuickPick(
          [
            { label: 'No parent (make root task)', value: null },
            ...tasks.map((t: any) => ({ label: t.baslik, value: t.id }))
          ],
          {
            placeHolder: 'Select new parent task',
          }
        );

        if (!parentChoice) return;

        await mcpClient.callTool('gorev_ust_degistir', {
          gorev_id: item.task.id,
          yeni_parent_id: parentChoice.value || '',
        });

        vscode.window.showInformationMessage('Parent changed successfully');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to change parent: ${error}`);
      }
    })
  );

  // Remove Parent (make root task)
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.REMOVE_PARENT, async (item: any) => {
      try {
        await mcpClient.callTool('gorev_ust_degistir', {
          gorev_id: item.task.id,
          yeni_parent_id: '',
        });

        vscode.window.showInformationMessage('Task is now a root task');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to remove parent: ${error}`);
      }
    })
  );

  // Add Dependency
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.ADD_DEPENDENCY, async (item: any) => {
      try {
        // Get all tasks except the current one
        const result = await mcpClient.callTool('gorev_listele', {
          tum_projeler: true,
        });
        
        if (!result || !result.content || !result.content[0]) {
          throw new Error('Failed to fetch tasks');
        }

        // Parse tasks from the result
        const tasks = result.content[0].text
          .split('\n')
          .filter((line: string) => line.includes('ID:'))
          .map((line: string) => {
            const idMatch = line.match(/ID:\s*([a-f0-9-]+)/);
            const titleMatch = line.match(/\[.+\]\s+(.+?)\s+\(/);
            const statusMatch = line.match(/\[(beklemede|devam_ediyor|tamamlandi)\]/);
            return idMatch && titleMatch ? { 
              id: idMatch[1], 
              baslik: titleMatch[1].trim(),
              durum: statusMatch ? statusMatch[1] : 'beklemede'
            } : null;
          })
          .filter((task: any) => task && task.id !== item.task.id);

        if (tasks.length === 0) {
          vscode.window.showInformationMessage('No other tasks available to add as dependency');
          return;
        }

        // Create quick pick items with status icons
        const quickPickItems = tasks.map((t: any) => ({
          label: `${t.durum === 'tamamlandi' ? '✓' : t.durum === 'devam_ediyor' ? '▶' : '○'} ${t.baslik}`,
          description: t.durum === 'tamamlandi' ? 'Completed' : t.durum === 'devam_ediyor' ? 'In Progress' : 'Pending',
          value: t.id
        }));

        const selectedTask = await vscode.window.showQuickPick(
          quickPickItems,
          {
            placeHolder: `Select a task that "${item.task.baslik}" depends on`,
          }
        );

        if (!selectedTask) return;

        // Add the dependency
        await mcpClient.callTool('gorev_bagimlilik_ekle', {
          kaynak_id: item.task.id,
          hedef_id: selectedTask.value,
          baglanti_tipi: 'bagimli', // depends on
        });

        vscode.window.showInformationMessage('Dependency added successfully');
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to add dependency: ${error}`);
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