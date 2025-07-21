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
  // Create Task - Now redirects to template wizard due to mandatory template requirement
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_TASK, async () => {
      // Redirect to template wizard since direct task creation is no longer allowed
      await vscode.commands.executeCommand(COMMANDS.OPEN_TEMPLATE_WIZARD);
    })
  );

  // Quick Create Task - Now uses quick template selection
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.QUICK_CREATE_TASK, async () => {
      // Use the quick template selection command
      await vscode.commands.executeCommand(COMMANDS.QUICK_CREATE_FROM_TEMPLATE);
    })
  );

  // Refresh Tasks
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.REFRESH_TASKS, async () => {
      try {
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(vscode.l10n.t('error.refreshTasks') + `: ${error}`);
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
        vscode.window.showErrorMessage(vscode.l10n.t('error.showTaskDetail') + `: ${error}`);
      }
    })
  );

  // Update Task Status
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.UPDATE_TASK_STATUS, async (item: any) => {
      try {
        const newStatus = await vscode.window.showQuickPick(
          [
            { label: vscode.l10n.t('status.pending'), value: GorevDurum.Beklemede },
            { label: vscode.l10n.t('status.inProgress'), value: GorevDurum.DevamEdiyor },
            { label: vscode.l10n.t('status.completed'), value: GorevDurum.Tamamlandi },
          ],
          {
            placeHolder: vscode.l10n.t('input.selectStatus'),
          }
        );

        if (!newStatus) return;

        await mcpClient.callTool('gorev_guncelle', {
          id: item.task.id,
          durum: newStatus.value,
        });

        vscode.window.showInformationMessage(vscode.l10n.t('success.taskStatusUpdated'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(vscode.l10n.t('error.updateTaskStatus') + `: ${error}`);
      }
    })
  );

  // Delete Task
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DELETE_TASK, async (item: any) => {
      try {
        const confirm = await vscode.window.showWarningMessage(
          vscode.l10n.t('confirm.deleteTask', item.task.baslik),
          vscode.l10n.t('confirm.yes'),
          vscode.l10n.t('confirm.no')
        );

        if (confirm !== vscode.l10n.t('confirm.yes')) return;

        await mcpClient.callTool('gorev_sil', {
          id: item.task.id,
          onay: true,
        });

        vscode.window.showInformationMessage(vscode.l10n.t('success.taskDeleted'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(vscode.l10n.t('error.deleteTask') + `: ${error}`);
      }
    })
  );

  // Create Subtask
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_SUBTASK, async (item: any) => {
      try {
        const baslik = await vscode.window.showInputBox({
          prompt: vscode.l10n.t('input.subtaskTitle'),
          placeHolder: vscode.l10n.t('placeholder.createSubtask', item.task.baslik),
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return vscode.l10n.t('validation.subtaskTitleRequired');
            }
            return null;
          },
        });

        if (!baslik) return;

        const aciklama = await vscode.window.showInputBox({
          prompt: vscode.l10n.t('input.subtaskDescription'),
          placeHolder: vscode.l10n.t('placeholder.subtaskDescription'),
        });

        const oncelik = await vscode.window.showQuickPick(
          [
            { label: vscode.l10n.t('priority.high'), value: GorevOncelik.Yuksek },
            { label: vscode.l10n.t('priority.medium'), value: GorevOncelik.Orta },
            { label: vscode.l10n.t('priority.low'), value: GorevOncelik.Dusuk },
          ],
          {
            placeHolder: vscode.l10n.t('input.selectPriority'),
          }
        );

        if (!oncelik) return;

        await mcpClient.callTool('gorev_altgorev_olustur', {
          parent_id: item.task.id,
          baslik,
          aciklama: aciklama || '',
          oncelik: oncelik.value,
        });

        vscode.window.showInformationMessage(vscode.l10n.t('success.subtaskCreated'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(vscode.l10n.t('error.createSubtask') + `: ${error}`);
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
          throw new Error(vscode.l10n.t('error.fetchTasks'));
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
            { label: vscode.l10n.t('parent.noParent'), value: null },
            ...tasks.map((t: any) => ({ label: t.baslik, value: t.id }))
          ],
          {
            placeHolder: vscode.l10n.t('input.selectParentTask'),
          }
        );

        if (!parentChoice) return;

        await mcpClient.callTool('gorev_ust_degistir', {
          gorev_id: item.task.id,
          yeni_parent_id: parentChoice.value || '',
        });

        vscode.window.showInformationMessage(vscode.l10n.t('success.parentChanged'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(vscode.l10n.t('error.changeParent') + `: ${error}`);
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

        vscode.window.showInformationMessage(vscode.l10n.t('success.taskIsRootNow'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(vscode.l10n.t('error.removeParent') + `: ${error}`);
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
          throw new Error(vscode.l10n.t('error.fetchTasks'));
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
          vscode.window.showInformationMessage(vscode.l10n.t('info.noTasksForDependency'));
          return;
        }

        // Create quick pick items with status icons
        const quickPickItems = tasks.map((t: any) => ({
          label: `${t.durum === 'tamamlandi' ? '✓' : t.durum === 'devam_ediyor' ? '▶' : '○'} ${t.baslik}`,
          description: t.durum === 'tamamlandi' ? vscode.l10n.t('status.completed') : t.durum === 'devam_ediyor' ? vscode.l10n.t('status.inProgress') : vscode.l10n.t('status.pending'),
          value: t.id
        }));

        const selectedTask = await vscode.window.showQuickPick(
          quickPickItems,
          {
            placeHolder: vscode.l10n.t('input.selectDependency', item.task.baslik),
          }
        );

        if (!selectedTask) return;

        // Add the dependency
        await mcpClient.callTool('gorev_bagimlilik_ekle', {
          kaynak_id: item.task.id,
          hedef_id: selectedTask.value,
          baglanti_tipi: 'bagimli', // depends on
        });

        vscode.window.showInformationMessage(vscode.l10n.t('success.dependencyAdded'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(vscode.l10n.t('error.addDependency') + `: ${error}`);
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