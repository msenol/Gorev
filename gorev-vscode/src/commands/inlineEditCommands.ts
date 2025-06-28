import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { COMMANDS } from '../utils/constants';
import { InlineEditProvider } from '../providers/inlineEditProvider';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';

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
                    vscode.window.showWarningMessage('Lütfen düzenlemek için bir görev seçin');
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
            if (!item || !item.task) {
                vscode.window.showWarningMessage('Lütfen bir görev seçin');
                return;
            }

            await editProvider.quickStatusChange(item.task);
            await treeProvider.refresh();
        })
    );

    // Quick Priority Change
    context.subscriptions.push(
        vscode.commands.registerCommand(COMMANDS.QUICK_PRIORITY_CHANGE, async (item: any) => {
            if (!item || !item.task) {
                vscode.window.showWarningMessage('Lütfen bir görev seçin');
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
                vscode.window.showWarningMessage('Lütfen bir görev seçin');
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
                vscode.window.showWarningMessage('Lütfen bir görev seçin');
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
                vscode.window.showInformationMessage('Düzenleme iptal edildi');
            }
        })
    );
}