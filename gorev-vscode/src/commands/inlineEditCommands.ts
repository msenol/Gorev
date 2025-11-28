import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ApiClient } from '../api/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { InlineEditProvider } from '../providers/inlineEditProvider';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';
import { GorevTreeItem } from '../providers/gorevTreeProvider';
import { Logger } from '../utils/logger';

export function registerInlineEditCommands(
    context: vscode.ExtensionContext,
    apiClient: ApiClient,
    providers: CommandContext
): void {
    const editProvider = new InlineEditProvider(apiClient);
    const treeProvider = providers.gorevTreeProvider as EnhancedGorevTreeProvider;

    // Edit Task Title (F2)
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.EDIT_TASK_TITLE, async (item?: GorevTreeItem) => {
            // Eğer item yoksa, seçili görevi al
            if (!item) {
                const selectedTasks = treeProvider.getSelectedTasks();
                if (selectedTasks.length === 1) {
                    item = new GorevTreeItem(selectedTasks[0]);
                } else {
                    vscode.window.showWarningMessage(t('inlineEdit.selectTaskToEdit'));
                    return;
                }
            }

            await editProvider.startEdit(item);
            await treeProvider.refresh();
        })
    );

    // Quick Status Change
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.QUICK_STATUS_CHANGE, async (item: GorevTreeItem) => {
            Logger.info(t('inlineEdit.quickStatusStartLog'));
            Logger.info(t('inlineEdit.quickStatusItemTypeLog'), item?.constructor?.name);
            Logger.info(t('inlineEdit.quickStatusHasTaskLog'), !!item?.task);
            
            if (item?.task) {
                Logger.info(t('inlineEdit.quickStatusTaskIdLog'), item.task.id);
                Logger.info(t('inlineEdit.quickStatusTaskTitleLog'), item.task.baslik);
                Logger.info(t('inlineEdit.quickStatusCurrentStatusLog'), item.task.durum);
            }
            
            if (!item || !item.task) {
                Logger.warn(t('inlineEdit.quickStatusNoTaskLog'));
                vscode.window.showWarningMessage(t('inlineEdit.selectTask'));
                return;
            }

            try {
                Logger.info(t('inlineEdit.quickStatusCallingEditLog'));
                await editProvider.quickStatusChange(item.task);
                // Add a small delay to ensure the backend has processed the update
                await new Promise(resolve => setTimeout(resolve, 100));
                Logger.info(t('inlineEdit.quickStatusRefreshingLog'));
                await treeProvider.refresh();
                
                // Also refresh the project tree if it exists
                const projeTreeProvider = (global as unknown as { projeTreeProvider?: { refresh: () => Promise<void> } }).projeTreeProvider;
                if (projeTreeProvider) {
                    await projeTreeProvider.refresh();
                }
            } catch (error) {
                Logger.error(t('inlineEdit.quickStatusFailed'), error);
            }
        })
    );

    // Quick Priority Change
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.QUICK_PRIORITY_CHANGE, async (item: GorevTreeItem) => {
            if (!item || !item.task) {
                vscode.window.showWarningMessage(t('inlineEdit.selectTask'));
                return;
            }

            await editProvider.quickPriorityChange(item.task);
            await treeProvider.refresh();
        })
    );

    // Quick Date Change
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.QUICK_DATE_CHANGE, async (item: GorevTreeItem) => {
            if (!item || !item.task) {
                vscode.window.showWarningMessage(t('inlineEdit.selectTask'));
                return;
            }

            await editProvider.quickDateChange(item.task);
            await treeProvider.refresh();
        })
    );

    // Detailed Edit
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.DETAILED_EDIT, async (item: GorevTreeItem) => {
            if (!item || !item.task) {
                vscode.window.showWarningMessage(t('inlineEdit.selectTask'));
                return;
            }

            await editProvider.showDetailedEdit(item.task);
            await treeProvider.refresh();
        })
    );

    // Double-click to edit title
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.onTreeItemDoubleClick', async (item: GorevTreeItem) => {
            if (item && item.task) {
                await editProvider.startEdit(item);
                await treeProvider.refresh();
            }
        })
    );

    // ESC to cancel edit
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.cancelEdit', () => {
            if (editProvider.isEditing()) {
                editProvider.cancelEdit();
                vscode.window.showInformationMessage(t('inlineEdit.editCancelled'));
            }
        })
    );
}
