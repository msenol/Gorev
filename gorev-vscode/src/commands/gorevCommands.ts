import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ApiClient, ApiError } from '../api/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { GorevDurum, GorevOncelik } from '../models/common';
import { Gorev } from '../models/gorev';
import { TaskDetailPanel } from '../ui/taskDetailPanel';
import { Logger } from '../utils/logger';
import { GorevTreeItem } from '../providers/gorevTreeProvider';

export function registerGorevCommands(
  context: vscode.ExtensionContext,
  apiClient: ApiClient,
  providers: CommandContext
): void {
  // Initialize API client
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
        vscode.window.showErrorMessage(t('error.refreshTasks') + `: ${error}`);
      }
    })
  );

  // Show Task Detail
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SHOW_TASK_DETAIL, async (item: GorevTreeItem) => {
      try {
        // Use the new TaskDetailPanel
        await TaskDetailPanel.createOrShow(
          apiClient,
          item.task,
          context.extensionUri
        );
      } catch (error) {
        vscode.window.showErrorMessage(t('error.showTaskDetail') + `: ${error}`);
      }
    })
  );

  // Update Task Status
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.UPDATE_TASK_STATUS, async (item: GorevTreeItem) => {
      try {
        const newStatus = await vscode.window.showQuickPick(
          [
            { label: t('status.pending'), value: GorevDurum.Beklemede },
            { label: t('status.inProgress'), value: GorevDurum.DevamEdiyor },
            { label: t('status.completed'), value: GorevDurum.Tamamlandi },
          ],
          {
            placeHolder: t('input.selectStatus'),
          }
        );

        if (!newStatus) return;

        // Use REST API to update task
        await apiClient.updateTask(item.task.id, {
          durum: newStatus.value,
        });

        vscode.window.showInformationMessage(t('success.taskStatusUpdated'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[UpdateTaskStatus] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('error.updateTaskStatus') + `: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(t('error.updateTaskStatus') + `: ${error}`);
        }
      }
    })
  );

  // Delete Task
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DELETE_TASK, async (item: GorevTreeItem) => {
      try {
        const confirm = await vscode.window.showWarningMessage(
          t('confirm.deleteTask', item.task.baslik),
          t('confirm.yes'),
          t('confirm.no')
        );

        if (confirm !== t('confirm.yes')) return;

        // Use REST API to delete task
        await apiClient.deleteTask(item.task.id);

        vscode.window.showInformationMessage(t('success.taskDeleted'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[DeleteTask] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('error.deleteTask') + `: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(t('error.deleteTask') + `: ${error}`);
        }
      }
    })
  );

  // Create Subtask
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CREATE_SUBTASK, async (item: GorevTreeItem) => {
      try {
        const baslik = await vscode.window.showInputBox({
          prompt: t('input.subtaskTitle'),
          placeHolder: t('placeholder.createSubtask', item.task.baslik),
          validateInput: (value) => {
            if (!value || value.trim().length === 0) {
              return t('validation.subtaskTitleRequired');
            }
            return null;
          },
        });

        if (!baslik) return;

        const aciklama = await vscode.window.showInputBox({
          prompt: t('input.subtaskDescription'),
          placeHolder: t('placeholder.subtaskDescription'),
        });

        const oncelik = await vscode.window.showQuickPick(
          [
            { label: t('priority.high'), value: GorevOncelik.Yuksek },
            { label: t('priority.medium'), value: GorevOncelik.Orta },
            { label: t('priority.low'), value: GorevOncelik.Dusuk },
          ],
          {
            placeHolder: t('input.selectPriority'),
          }
        );

        if (!oncelik) return;

        // Use REST API to create subtask
        await apiClient.createSubtask(item.task.id, {
          baslik,
          aciklama: aciklama || '',
          oncelik: oncelik.value,
        });

        vscode.window.showInformationMessage(t('success.subtaskCreated'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[CreateSubtask] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('error.createSubtask') + `: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(t('error.createSubtask') + `: ${error}`);
        }
      }
    })
  );

  // Change Parent
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CHANGE_PARENT, async (item: GorevTreeItem | { task: Gorev }) => {
      try {
        // Support both GorevTreeItem (from tree view) and { task: Gorev } (from task detail panel)
        const task = 'task' in item ? item.task : item;
        
        // Use REST API to get all tasks
        const response = await apiClient.getTasks({
          tum_projeler: true,
        });

        if (!response.success || !response.data) {
          throw new Error(t('error.fetchTasks'));
        }

        // Filter out current task and show only tasks from the same project
        // (parent and child must be in the same project)
        const availableTasks = response.data.filter(taskItem => 
          taskItem.id !== task.id && 
          taskItem.proje_id === task.proje_id
        );

        const parentChoice = await vscode.window.showQuickPick(
          [
            { label: t('parent.noParent'), value: null },
            ...availableTasks.map(taskItem => ({ label: taskItem.baslik, value: taskItem.id }))
          ],
          {
            placeHolder: t('input.selectParentTask'),
          }
        );

        if (!parentChoice) return;

        // Use REST API to change parent
        await apiClient.changeParent(task.id, parentChoice.value || '');

        vscode.window.showInformationMessage(t('success.parentChanged'));
        await providers.gorevTreeProvider.refresh();
        
        // Refresh task detail panel if open
        await TaskDetailPanel.refreshIfOpen(task.id);
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[ChangeParent] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('error.changeParent') + `: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(t('error.changeParent') + `: ${error}`);
        }
      }
    })
  );

  // Remove Parent (make root task)
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.REMOVE_PARENT, async (item: GorevTreeItem | { task: Gorev }) => {
      try {
        // Support both GorevTreeItem (from tree view) and { task: Gorev } (from task detail panel)
        const task = 'task' in item ? item.task : item;
        
        // Use REST API to remove parent (empty string makes it root)
        await apiClient.changeParent(task.id, '');

        vscode.window.showInformationMessage(t('success.taskIsRootNow'));
        await providers.gorevTreeProvider.refresh();
        
        // Refresh task detail panel if open
        await TaskDetailPanel.refreshIfOpen(task.id);
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[RemoveParent] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('error.removeParent') + `: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(t('error.removeParent') + `: ${error}`);
        }
      }
    })
  );

  // Add Dependency
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.ADD_DEPENDENCY, async (item: GorevTreeItem | { task: Gorev }) => {
      try {
        // Support both GorevTreeItem (from tree view) and { task: Gorev } (from task detail panel)
        const task = 'task' in item ? item.task : item;
        
        // Use REST API to get all tasks
        const response = await apiClient.getTasks({
          tum_projeler: true,
        });

        if (!response.success || !response.data) {
          throw new Error(t('error.fetchTasks'));
        }

        // Filter out current task
        const availableTasks = response.data.filter(t => t.id !== task.id);

        if (availableTasks.length === 0) {
          vscode.window.showInformationMessage(t('info.noTasksForDependency'));
          return;
        }

        // Create quick pick items with status icons
        const quickPickItems = availableTasks.map(taskItem => ({
          label: `${taskItem.durum === 'tamamlandi' ? '✓' : taskItem.durum === 'devam_ediyor' ? '▶' : '○'} ${taskItem.baslik}`,
          description: taskItem.durum === 'tamamlandi' ? t('status.completed') : taskItem.durum === 'devam_ediyor' ? t('status.inProgress') : t('status.pending'),
          value: taskItem.id
        }));

        const selectedTask = await vscode.window.showQuickPick(
          quickPickItems,
          {
            placeHolder: t('input.selectDependency', task.baslik),
          }
        );

        if (!selectedTask) return;

        // Use REST API to add dependency
        await apiClient.addDependency(task.id, {
          kaynak_id: selectedTask.value,
          baglanti_tipi: 'bagimli', // depends on
        });

        vscode.window.showInformationMessage(t('success.dependencyAdded'));
        await providers.gorevTreeProvider.refresh();
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[AddDependency] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(t('error.addDependency') + `: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(t('error.addDependency') + `: ${error}`);
        }
      }
    })
  );

}