import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ClientInterface } from '../interfaces/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { InlineEditProvider } from '../providers/inlineEditProvider';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';
import { Logger } from '../utils/logger';

export function registerInlineEditCommands(
    context: vscode.ExtensionContext,
    mcpClient: ClientInterface,
    providers: CommandContext
): void {
    const editProvider = new InlineEditProvider(mcpClient);
    const treeProvider = providers.gorevTreeProvider as EnhancedGorevTreeProvider;

    // Edit Task Title (F2)
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.EDIT_TASK_TITLE, async (item?: any) => {
            // Eğer item yoksa, seçili görevi al
            if (!item) {
                const selectedTasks = treeProvider.getSelectedTasks();
                if (selectedTasks.length === 1) {
                    item = { task: selectedTasks[0] };
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
        vscode.commands.registerCommand(COMMANDS.QUICK_STATUS_CHANGE, async (item: any) => {
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
                const projeTreeProvider = (global as any).projeTreeProvider;
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
        vscode.commands.registerCommand(COMMANDS.QUICK_PRIORITY_CHANGE, async (item: any) => {
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
        vscode.commands.registerCommand(COMMANDS.QUICK_DATE_CHANGE, async (item: any) => {
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
        vscode.commands.registerCommand(COMMANDS.DETAILED_EDIT, async (item: any) => {
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
        vscode.commands.registerCommand('gorev.onTreeItemDoubleClick', async (item: any) => {
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
