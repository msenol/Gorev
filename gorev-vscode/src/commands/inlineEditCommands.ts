import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { InlineEditProvider } from '../providers/inlineEditProvider';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';
import { Logger } from '../utils/logger';

export function registerInlineEditCommands(
    context: vscode.ExtensionContext,
    mcpClient: MCPClient,
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
                    vscode.window.showWarningMessage(vscode.l10n.t('inlineEdit.selectTaskToEdit'));
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
            Logger.info(vscode.l10n.t('inlineEdit.quickStatusStartLog'));
            Logger.info(vscode.l10n.t('inlineEdit.quickStatusItemTypeLog'), item?.constructor?.name);
            Logger.info(vscode.l10n.t('inlineEdit.quickStatusHasTaskLog'), !!item?.task);
            
            if (item?.task) {
                Logger.info(vscode.l10n.t('inlineEdit.quickStatusTaskIdLog'), item.task.id);
                Logger.info(vscode.l10n.t('inlineEdit.quickStatusTaskTitleLog'), item.task.baslik);
                Logger.info(vscode.l10n.t('inlineEdit.quickStatusCurrentStatusLog'), item.task.durum);
            }
            
            if (!item || !item.task) {
                Logger.warn(vscode.l10n.t('inlineEdit.quickStatusNoTaskLog'));
                vscode.window.showWarningMessage(vscode.l10n.t('inlineEdit.selectTask'));
                return;
            }

            try {
                Logger.info(vscode.l10n.t('inlineEdit.quickStatusCallingEditLog'));
                await editProvider.quickStatusChange(item.task);
                // Add a small delay to ensure the backend has processed the update
                await new Promise(resolve => setTimeout(resolve, 100));
                Logger.info(vscode.l10n.t('inlineEdit.quickStatusRefreshingLog'));
                await treeProvider.refresh();
                
                // Also refresh the project tree if it exists
                const projeTreeProvider = (global as any).projeTreeProvider;
                if (projeTreeProvider) {
                    await projeTreeProvider.refresh();
                }
            } catch (error) {
                Logger.error(vscode.l10n.t('inlineEdit.quickStatusFailed'), error);
            }
        })
    );

    // Quick Priority Change
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.QUICK_PRIORITY_CHANGE, async (item: any) => {
            if (!item || !item.task) {
                vscode.window.showWarningMessage(vscode.l10n.t('inlineEdit.selectTask'));
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
                vscode.window.showWarningMessage(vscode.l10n.t('inlineEdit.selectTask'));
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
                vscode.window.showWarningMessage(vscode.l10n.t('inlineEdit.selectTask'));
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
                vscode.window.showInformationMessage(vscode.l10n.t('inlineEdit.editCancelled'));
            }
        })
    );
}