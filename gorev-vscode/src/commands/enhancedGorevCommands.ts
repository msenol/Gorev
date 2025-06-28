import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { GorevDurum, GorevOncelik } from '../models/common';
import { GroupingStrategy, SortingCriteria, TaskFilter } from '../models/treeModels';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';

export function registerEnhancedGorevCommands(
  context: vscode.ExtensionContext,
  mcpClient: MCPClient,
  providers: CommandContext
): void {
  const treeProvider = providers.gorevTreeProvider as EnhancedGorevTreeProvider;

  // Select Task Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SELECT_TASK, (taskId: string, event?: any) => {
      const multiSelect = event?.ctrlKey || event?.metaKey || false;
      const rangeSelect = event?.shiftKey || false;
      treeProvider.selectTask(taskId, multiSelect, rangeSelect);
    })
  );

  // Set Grouping Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SET_GROUPING, async () => {
      const items = [
        { label: 'No Grouping', value: GroupingStrategy.None },
        { label: 'By Status', value: GroupingStrategy.ByStatus },
        { label: 'By Priority', value: GroupingStrategy.ByPriority },
        { label: 'By Project', value: GroupingStrategy.ByProject },
        { label: 'By Tag', value: GroupingStrategy.ByTag },
        { label: 'By Due Date', value: GroupingStrategy.ByDueDate },
      ];

      const selected = await vscode.window.showQuickPick(items, {
        placeHolder: 'Select grouping strategy',
      });

      if (selected) {
        treeProvider.setGrouping(selected.value);
      }
    })
  );

  // Set Sorting Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SET_SORTING, async () => {
      const criteria = [
        { label: 'By Title', value: SortingCriteria.Title },
        { label: 'By Priority', value: SortingCriteria.Priority },
        { label: 'By Due Date', value: SortingCriteria.DueDate },
        { label: 'By Created Date', value: SortingCriteria.CreatedDate },
        { label: 'By Status', value: SortingCriteria.Status },
      ];

      const selectedCriteria = await vscode.window.showQuickPick(criteria, {
        placeHolder: 'Select sorting criteria',
      });

      if (!selectedCriteria) return;

      const order = await vscode.window.showQuickPick(
        [
          { label: 'Ascending', value: true },
          { label: 'Descending', value: false },
        ],
        {
          placeHolder: 'Select sort order',
        }
      );

      if (order) {
        treeProvider.setSorting(selectedCriteria.value, order.value);
      }
    })
  );

  // Filter Tasks Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.FILTER_TASKS, async () => {
      const action = await vscode.window.showQuickPick(
        [
          { label: '🔍 Search by text', value: 'search' },
          { label: '📊 Filter by status', value: 'status' },
          { label: '🎯 Filter by priority', value: 'priority' },
          { label: '🏷️ Filter by tag', value: 'tag' },
          { label: '📅 Filter by due date', value: 'dueDate' },
          { label: '⚡ Quick filters', value: 'quick' },
        ],
        {
          placeHolder: 'Select filter type',
        }
      );

      if (!action) return;

      const filter: Partial<TaskFilter> = {};

      switch (action.value) {
        case 'search':
          const query = await vscode.window.showInputBox({
            prompt: 'Search tasks',
            placeHolder: 'Enter search text',
          });
          if (query) {
            filter.searchQuery = query;
          }
          break;

        case 'status':
          const statuses = await vscode.window.showQuickPick(
            [
              { label: 'Pending', value: GorevDurum.Beklemede, picked: true },
              { label: 'In Progress', value: GorevDurum.DevamEdiyor, picked: true },
              { label: 'Completed', value: GorevDurum.Tamamlandi, picked: false },
            ],
            {
              placeHolder: 'Select statuses to show',
              canPickMany: true,
            }
          );
          if (statuses && statuses.length > 0) {
            // For now, we only support single status filter
            filter.durum = statuses[0].value;
          }
          break;

        case 'priority':
          const priorities = await vscode.window.showQuickPick(
            [
              { label: 'High', value: GorevOncelik.Yuksek },
              { label: 'Medium', value: GorevOncelik.Orta },
              { label: 'Low', value: GorevOncelik.Dusuk },
            ],
            {
              placeHolder: 'Select priorities to show',
              canPickMany: true,
            }
          );
          if (priorities && priorities.length > 0) {
            // For now, we only support single priority filter
            filter.oncelik = priorities[0].value;
          }
          break;

        case 'tag':
          const tag = await vscode.window.showInputBox({
            prompt: 'Filter by tag',
            placeHolder: 'Enter tag name',
          });
          if (tag) {
            filter.tags = [tag];
          }
          break;

        case 'dueDate':
          const dateFilter = await vscode.window.showQuickPick(
            [
              { label: 'Overdue tasks', value: 'overdue' },
              { label: 'Due today', value: 'today' },
              { label: 'Due this week', value: 'week' },
              { label: 'Due this month', value: 'month' },
              { label: 'Has due date', value: 'hasDue' },
              { label: 'No due date', value: 'noDue' },
            ],
            {
              placeHolder: 'Select due date filter',
            }
          );

          if (dateFilter) {
            const now = new Date();
            const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
            const tomorrow = new Date(today);
            tomorrow.setDate(tomorrow.getDate() + 1);

            switch (dateFilter.value) {
              case 'overdue':
                filter.dueDateRange = { end: today };
                break;
              case 'today':
                filter.dueDateRange = { start: today, end: tomorrow };
                break;
              case 'week':
                const nextWeek = new Date(today);
                nextWeek.setDate(nextWeek.getDate() + 7);
                filter.dueDateRange = { start: today, end: nextWeek };
                break;
              case 'month':
                const nextMonth = new Date(today);
                nextMonth.setMonth(nextMonth.getMonth() + 1);
                filter.dueDateRange = { start: today, end: nextMonth };
                break;
              case 'hasDue':
                // Will be handled by checking if son_tarih exists
                break;
              case 'noDue':
                // Will be handled by checking if son_tarih is null
                break;
            }
          }
          break;

        case 'quick':
          const quickFilter = await vscode.window.showQuickPick(
            [
              { label: '🔥 High priority tasks', value: 'highPriority' },
              { label: '⚠️ Overdue tasks', value: 'overdue' },
              { label: '📅 Due today', value: 'today' },
              { label: '🏃 In progress', value: 'inProgress' },
              { label: '✅ Completed today', value: 'completedToday' },
            ],
            {
              placeHolder: 'Select quick filter',
            }
          );

          if (quickFilter) {
            switch (quickFilter.value) {
              case 'highPriority':
                filter.oncelik = GorevOncelik.Yuksek;
                break;
              case 'overdue':
                filter.overdue = true;
                break;
              case 'today':
                const today = new Date();
                const tomorrow = new Date(today);
                tomorrow.setDate(tomorrow.getDate() + 1);
                filter.dueDateRange = { start: today, end: tomorrow };
                break;
              case 'inProgress':
                filter.durum = GorevDurum.DevamEdiyor;
                break;
              case 'completedToday':
                filter.durum = GorevDurum.Tamamlandi;
                // TODO: Add completed date filter when available
                break;
            }
          }
          break;
      }

      treeProvider.updateFilter(filter);
    })
  );

  // Clear Filter Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CLEAR_FILTER, () => {
      treeProvider.updateFilter({});
      vscode.window.showInformationMessage('Filters cleared');
    })
  );

  // Toggle Show Completed Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.TOGGLE_SHOW_COMPLETED, () => {
      const config = vscode.workspace.getConfiguration('gorev.treeView');
      const current = config.get('showCompleted', true);
      config.update('showCompleted', !current, vscode.ConfigurationTarget.Global);
      vscode.window.showInformationMessage(
        `${!current ? 'Showing' : 'Hiding'} completed tasks`
      );
    })
  );

  // Select All Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SELECT_ALL, () => {
      const tasks = treeProvider.getSelectedTasks();
      // TODO: Implement select all
      vscode.window.showInformationMessage('Select all not yet implemented');
    })
  );

  // Deselect All Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DESELECT_ALL, () => {
      // TODO: Implement deselect all
      vscode.window.showInformationMessage('Deselect all not yet implemented');
    })
  );

  // Bulk Update Status Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.BULK_UPDATE_STATUS, async () => {
      const selectedTasks = treeProvider.getSelectedTasks();
      
      if (selectedTasks.length === 0) {
        vscode.window.showWarningMessage('No tasks selected');
        return;
      }

      const newStatus = await vscode.window.showQuickPick(
        [
          { label: 'Pending', value: GorevDurum.Beklemede },
          { label: 'In Progress', value: GorevDurum.DevamEdiyor },
          { label: 'Completed', value: GorevDurum.Tamamlandi },
        ],
        {
          placeHolder: `Update status for ${selectedTasks.length} tasks`,
        }
      );

      if (!newStatus) return;

      try {
        vscode.window.withProgress(
          {
            location: vscode.ProgressLocation.Notification,
            title: 'Updating tasks...',
            cancellable: false,
          },
          async (progress) => {
            let completed = 0;
            for (const task of selectedTasks) {
              await mcpClient.callTool('gorev_guncelle', {
                id: task.id,
                durum: newStatus.value,
              });
              completed++;
              progress.report({
                increment: (completed / selectedTasks.length) * 100,
                message: `${completed}/${selectedTasks.length}`,
              });
            }
          }
        );

        vscode.window.showInformationMessage(
          `Updated ${selectedTasks.length} tasks to ${newStatus.label}`
        );
        await treeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to update tasks: ${error}`);
      }
    })
  );

  // Bulk Delete Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.BULK_DELETE, async () => {
      const selectedTasks = treeProvider.getSelectedTasks();
      
      if (selectedTasks.length === 0) {
        vscode.window.showWarningMessage('No tasks selected');
        return;
      }

      const confirm = await vscode.window.showWarningMessage(
        `Are you sure you want to delete ${selectedTasks.length} tasks?`,
        'Yes',
        'No'
      );

      if (confirm !== 'Yes') return;

      try {
        vscode.window.withProgress(
          {
            location: vscode.ProgressLocation.Notification,
            title: 'Deleting tasks...',
            cancellable: false,
          },
          async (progress) => {
            let completed = 0;
            for (const task of selectedTasks) {
              await mcpClient.callTool('gorev_sil', {
                id: task.id,
                onay: true,
              });
              completed++;
              progress.report({
                increment: (completed / selectedTasks.length) * 100,
                message: `${completed}/${selectedTasks.length}`,
              });
            }
          }
        );

        vscode.window.showInformationMessage(
          `Deleted ${selectedTasks.length} tasks`
        );
        await treeProvider.refresh();
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to delete tasks: ${error}`);
      }
    })
  );
}