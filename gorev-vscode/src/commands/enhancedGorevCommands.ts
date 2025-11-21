import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ApiClient } from '../api/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { GorevDurum, GorevOncelik } from '../models/common';
import { GroupingStrategy, SortingCriteria, TaskFilter } from '../models/treeModels';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';

export function registerEnhancedGorevCommands(
  context: vscode.ExtensionContext,
  apiClient: ApiClient,
  providers: CommandContext
): void {
  const treeProvider = providers.gorevTreeProvider as EnhancedGorevTreeProvider;

  // Select Task Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SELECT_TASK, (taskId: string, event?: unknown) => {
      const evt = event as { ctrlKey?: boolean; metaKey?: boolean; shiftKey?: boolean } | undefined;
      const multiSelect = evt?.ctrlKey || evt?.metaKey || false;
      const rangeSelect = evt?.shiftKey || false;
      treeProvider.selectTask(taskId, multiSelect, rangeSelect);
    })
  );

  // Set Grouping Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SET_GROUPING, async () => {
      const items = [
        { label: t('grouping.none'), value: GroupingStrategy.None },
        { label: t('grouping.byStatus'), value: GroupingStrategy.ByStatus },
        { label: t('grouping.byPriority'), value: GroupingStrategy.ByPriority },
        { label: t('grouping.byProject'), value: GroupingStrategy.ByProject },
        { label: t('grouping.byTag'), value: GroupingStrategy.ByTag },
        { label: t('grouping.byDueDate'), value: GroupingStrategy.ByDueDate },
      ];

      const selected = await vscode.window.showQuickPick(items, {
        placeHolder: t('grouping.selectStrategy'),
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
        { label: t('sorting.byTitle'), value: SortingCriteria.Title },
        { label: t('sorting.byPriority'), value: SortingCriteria.Priority },
        { label: t('sorting.byDueDate'), value: SortingCriteria.DueDate },
        { label: t('sorting.byCreatedDate'), value: SortingCriteria.CreatedDate },
        { label: t('sorting.byStatus'), value: SortingCriteria.Status },
      ];

      const selectedCriteria = await vscode.window.showQuickPick(criteria, {
        placeHolder: t('sorting.selectCriteria'),
      });

      if (!selectedCriteria) return;

      const order = await vscode.window.showQuickPick(
        [
          { label: t('sorting.ascending'), value: true },
          { label: t('sorting.descending'), value: false },
        ],
        {
          placeHolder: t('sorting.selectOrder'),
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
          { label: t('filter.searchByText'), value: 'search' },
          { label: t('filter.byStatus'), value: 'status' },
          { label: t('filter.byPriority'), value: 'priority' },
          { label: t('filter.byTag'), value: 'tag' },
          { label: t('filter.byDueDate'), value: 'dueDate' },
          { label: t('filter.quickFilters'), value: 'quick' },
        ],
        {
          placeHolder: t('filter.selectType'),
        }
      );

      if (!action) return;

      const filter: Partial<TaskFilter> = {};

      switch (action.value) {
        case 'search': {
          const query = await vscode.window.showInputBox({
            prompt: t('filter.searchPrompt'),
            placeHolder: t('filter.searchPlaceholder'),
          });
          if (query) {
            filter.searchQuery = query;
          }
          break;
        }

        case 'status': {
          const statuses = await vscode.window.showQuickPick(
            [
              { label: t('status.pending'), value: GorevDurum.Beklemede, picked: true },
              { label: t('status.inProgress'), value: GorevDurum.DevamEdiyor, picked: true },
              { label: t('status.completed'), value: GorevDurum.Tamamlandi, picked: false },
            ],
            {
              placeHolder: t('filter.selectStatuses'),
              canPickMany: true,
            }
          );
          if (statuses && statuses.length > 0) {
            // For now, we only support single status filter
            filter.durum = statuses[0].value;
          }
          break;
        }

        case 'priority': {
          const priorities = await vscode.window.showQuickPick(
            [
              { label: t('priority.high'), value: GorevOncelik.Yuksek },
              { label: t('priority.medium'), value: GorevOncelik.Orta },
              { label: t('priority.low'), value: GorevOncelik.Dusuk },
            ],
            {
              placeHolder: t('filter.selectPriorities'),
              canPickMany: true,
            }
          );
          if (priorities && priorities.length > 0) {
            // For now, we only support single priority filter
            filter.oncelik = priorities[0].value;
          }
          break;
        }

        case 'tag': {
          const tag = await vscode.window.showInputBox({
            prompt: t('filter.tagPrompt'),
            placeHolder: t('filter.tagPlaceholder'),
          });
          if (tag) {
            filter.tags = [tag];
          }
          break;
        }

        case 'dueDate': {
          const dateFilter = await vscode.window.showQuickPick(
            [
              { label: t('filter.overdueTasks'), value: 'overdue' },
              { label: t('filter.dueToday'), value: 'today' },
              { label: t('filter.dueThisWeek'), value: 'week' },
              { label: t('filter.dueThisMonth'), value: 'month' },
              { label: t('filter.hasDueDate'), value: 'hasDue' },
              { label: t('filter.noDueDate'), value: 'noDue' },
            ],
            {
              placeHolder: t('filter.selectDueDateFilter'),
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
              case 'week': {
                const nextWeek = new Date(today);
                nextWeek.setDate(nextWeek.getDate() + 7);
                filter.dueDateRange = { start: today, end: nextWeek };
                break;
              }
              case 'month': {
                const nextMonth = new Date(today);
                nextMonth.setMonth(nextMonth.getMonth() + 1);
                filter.dueDateRange = { start: today, end: nextMonth };
                break;
              }
              case 'hasDue':
                // Will be handled by checking if son_tarih exists
                break;
              case 'noDue':
                // Will be handled by checking if son_tarih is null
                break;
            }
          }
          break;
        }

        case 'quick': {
          const quickFilter = await vscode.window.showQuickPick(
            [
              { label: t('filter.highPriority'), value: 'highPriority' },
              { label: t('filter.overdue'), value: 'overdue' },
              { label: t('filter.today'), value: 'today' },
              { label: t('filter.inProgress'), value: 'inProgress' },
              { label: t('filter.completedToday'), value: 'completedToday' },
            ],
            {
              placeHolder: t('filter.selectQuickFilter'),
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
              case 'today': {
                const today = new Date();
                const tomorrow = new Date(today);
                tomorrow.setDate(tomorrow.getDate() + 1);
                filter.dueDateRange = { start: today, end: tomorrow };
                break;
              }
              case 'inProgress':
                filter.durum = GorevDurum.DevamEdiyor;
                break;
              case 'completedToday':
                filter.durum = GorevDurum.Tamamlandi;
                // Note: Completed date filtering requires backend support for completion timestamps
                // Currently showing all completed tasks regardless of completion date
                break;
            }
          }
          break;
        }
      }

      treeProvider.updateFilter(filter);
    })
  );

  // Clear Filter Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CLEAR_FILTER, () => {
      treeProvider.updateFilter({});
      vscode.window.showInformationMessage(t('filter.cleared'));
    })
  );

  // Toggle Show Completed Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.TOGGLE_SHOW_COMPLETED, () => {
      const config = vscode.workspace.getConfiguration('gorev.treeView');
      const current = config.get('showCompleted', true);
      config.update('showCompleted', !current, vscode.ConfigurationTarget.Global);
      vscode.window.showInformationMessage(
        !current ? t('toggle.showingCompleted') : t('toggle.hidingCompleted')
      );
    })
  );

  // Select All Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SELECT_ALL, () => {
      // Note: Full select all functionality requires TreeProvider enhancement
      // For now, show informational message about current selection
      const selectedTasks = treeProvider.getSelectedTasks();
      vscode.window.showInformationMessage(
        t('selectAll.currentlySelected', { count: selectedTasks.length.toString() })
      );
    })
  );

  // Deselect All Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DESELECT_ALL, () => {
      // Clear all selections in the tree provider
      treeProvider.clearSelection();
      vscode.window.showInformationMessage(t('selectAll.allCleared'));
    })
  );

  // Bulk Update Status Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.BULK_UPDATE_STATUS, async () => {
      const selectedTasks = treeProvider.getSelectedTasks();
      
      if (selectedTasks.length === 0) {
        vscode.window.showWarningMessage(t('bulk.noTasksSelected'));
        return;
      }

      const newStatus = await vscode.window.showQuickPick(
        [
          { label: t('status.pending'), value: GorevDurum.Beklemede },
          { label: t('status.inProgress'), value: GorevDurum.DevamEdiyor },
          { label: t('status.completed'), value: GorevDurum.Tamamlandi },
        ],
        {
          placeHolder: t('bulk.updateStatusPlaceholder', selectedTasks.length.toString()),
        }
      );

      if (!newStatus) return;

      try {
        vscode.window.withProgress(
          {
            location: vscode.ProgressLocation.Notification,
            title: t('bulk.updatingTasks'),
            cancellable: false,
          },
          async (progress) => {
            let completed = 0;
            for (const task of selectedTasks) {
              await apiClient.updateTask(task.id, {
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
          t('bulk.updatedTasks', selectedTasks.length.toString(), newStatus.label)
        );
        await treeProvider.refresh();
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : String(error);
        vscode.window.showErrorMessage(t('bulk.failedToUpdate', errorMessage));
      }
    })
  );

  // Bulk Delete Command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.BULK_DELETE, async () => {
      const selectedTasks = treeProvider.getSelectedTasks();
      
      if (selectedTasks.length === 0) {
        vscode.window.showWarningMessage(t('bulk.noTasksSelected'));
        return;
      }

      const confirm = await vscode.window.showWarningMessage(
        t('bulk.deleteConfirm', selectedTasks.length.toString()),
        t('confirm.yes'),
        t('confirm.no')
      );

      if (confirm !== t('confirm.yes')) return;

      try {
        vscode.window.withProgress(
          {
            location: vscode.ProgressLocation.Notification,
            title: t('bulk.deletingTasks'),
            cancellable: false,
          },
          async (progress) => {
            let completed = 0;
            for (const task of selectedTasks) {
              await apiClient.deleteTask(task.id);
              completed++;
              progress.report({
                increment: (completed / selectedTasks.length) * 100,
                message: `${completed}/${selectedTasks.length}`,
              });
            }
          }
        );

        vscode.window.showInformationMessage(
          t('bulk.deletedTasks', selectedTasks.length.toString())
        );
        await treeProvider.refresh();
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : String(error);
        vscode.window.showErrorMessage(t('bulk.failedToDelete', errorMessage));
      }
    })
  );
}
