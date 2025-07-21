import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { Gorev, GorevDurum, GorevOncelik } from '../models/gorev';
import { Logger } from '../utils/logger';

/**
 * TreeView item'ları için inline düzenleme sağlayıcı
 */
export class InlineEditProvider {
    private editingItem: any | null = null;
    private originalLabel: string | null = null;

    constructor(private mcpClient: MCPClient) {}

    /**
     * Inline düzenlemeyi başlatır
     */
    async startEdit(item: any): Promise<void> {
        if (!item || !item.task) {
            return;
        }

        this.editingItem = item;
        this.originalLabel = item.task.baslik;

        const newTitle = await vscode.window.showInputBox({
            prompt: vscode.l10n.t('inlineEdit.editTaskTitle'),
            value: item.task.baslik,
            validateInput: (value) => {
                if (!value || value.trim().length === 0) {
                    return vscode.l10n.t('inlineEdit.taskTitleEmpty');
                }
                if (value.length > 200) {
                    return vscode.l10n.t('inlineEdit.taskTitleTooLong');
                }
                return null;
            }
        });

        if (newTitle !== undefined) {
            if (newTitle !== this.originalLabel) {
                await this.saveEdit(item.task, newTitle);
            }
        }

        this.editingItem = null;
        this.originalLabel = null;
    }

    /**
     * Düzenlemeyi kaydeder
     */
    private async saveEdit(task: Gorev, newTitle: string): Promise<void> {
        try {
            await this.mcpClient.callTool('gorev_duzenle', {
                id: task.id,
                baslik: newTitle
            });

            vscode.window.showInformationMessage(vscode.l10n.t('inlineEdit.taskTitleUpdated'));
            Logger.info(`Task ${task.id} title updated to: ${newTitle}`);
        } catch (error) {
            vscode.window.showErrorMessage(vscode.l10n.t('inlineEdit.updateFailed', error));
            Logger.error('Failed to update task title:', error);
        }
    }

    /**
     * Hızlı durum değiştirme menüsü
     */
    async quickStatusChange(task: Gorev): Promise<void> {
        const items = [
            { 
                label: vscode.l10n.t('inlineEdit.pending'), 
                value: GorevDurum.Beklemede,
                description: task.durum === GorevDurum.Beklemede ? vscode.l10n.t('inlineEdit.currentStatus') : ''
            },
            { 
                label: vscode.l10n.t('inlineEdit.inProgress'), 
                value: GorevDurum.DevamEdiyor,
                description: task.durum === GorevDurum.DevamEdiyor ? vscode.l10n.t('inlineEdit.currentStatus') : ''
            },
            { 
                label: vscode.l10n.t('inlineEdit.completed'), 
                value: GorevDurum.Tamamlandi,
                description: task.durum === GorevDurum.Tamamlandi ? vscode.l10n.t('inlineEdit.currentStatus') : ''
            }
        ];

        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: vscode.l10n.t('inlineEdit.selectNewStatus'),
            title: vscode.l10n.t('inlineEdit.changeStatusFor', task.baslik)
        });

        if (selected && selected.value !== task.durum) {
            try {
                Logger.info(`[QuickStatusChange] Updating task ${task.id} from ${task.durum} to ${selected.value}`);
                
                const result = await this.mcpClient.callTool('gorev_guncelle', {
                    id: task.id,
                    durum: selected.value
                });
                
                Logger.info(`[QuickStatusChange] MCP response:`, JSON.stringify(result));

                vscode.window.showInformationMessage(vscode.l10n.t('inlineEdit.taskStatusUpdated'));
                Logger.info(`[QuickStatusChange] Task ${task.id} status updated to: ${selected.value}`);
                
                // Force a command execution to refresh all trees
                await vscode.commands.executeCommand('gorev.refreshTasks');
            } catch (error) {
                vscode.window.showErrorMessage(vscode.l10n.t('inlineEdit.statusUpdateFailed', error));
                Logger.error('[QuickStatusChange] Failed to update task status:', error);
            }
        }
    }

    /**
     * Hızlı öncelik değiştirme menüsü
     */
    async quickPriorityChange(task: Gorev): Promise<void> {
        const items = [
            { 
                label: vscode.l10n.t('inlineEdit.highPriority'), 
                value: GorevOncelik.Yuksek,
                description: task.oncelik === GorevOncelik.Yuksek ? vscode.l10n.t('inlineEdit.currentPriority') : ''
            },
            { 
                label: vscode.l10n.t('inlineEdit.mediumPriority'), 
                value: GorevOncelik.Orta,
                description: task.oncelik === GorevOncelik.Orta ? vscode.l10n.t('inlineEdit.currentPriority') : ''
            },
            { 
                label: vscode.l10n.t('inlineEdit.lowPriority'), 
                value: GorevOncelik.Dusuk,
                description: task.oncelik === GorevOncelik.Dusuk ? vscode.l10n.t('inlineEdit.currentPriority') : ''
            }
        ];

        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: vscode.l10n.t('inlineEdit.selectNewPriority'),
            title: vscode.l10n.t('inlineEdit.changePriorityFor', task.baslik)
        });

        if (selected && selected.value !== task.oncelik) {
            try {
                await this.mcpClient.callTool('gorev_duzenle', {
                    id: task.id,
                    oncelik: selected.value
                });

                vscode.window.showInformationMessage(vscode.l10n.t('inlineEdit.taskPriorityUpdated'));
                Logger.info(`Task ${task.id} priority updated to: ${selected.value}`);
            } catch (error) {
                vscode.window.showErrorMessage(vscode.l10n.t('inlineEdit.priorityUpdateFailed', error));
                Logger.error('Failed to update task priority:', error);
            }
        }
    }

    /**
     * Hızlı tarih değiştirme
     */
    async quickDateChange(task: Gorev): Promise<void> {
        const currentDate = task.son_tarih || '';
        
        const newDate = await vscode.window.showInputBox({
            prompt: vscode.l10n.t('inlineEdit.enterDueDate'),
            value: currentDate,
            placeHolder: '2024-12-31',
            validateInput: (value) => {
                if (!value) {
                    return null; // Boş bırakılabilir
                }
                const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
                if (!dateRegex.test(value)) {
                    return vscode.l10n.t('inlineEdit.invalidDateFormat');
                }
                const date = new Date(value);
                if (isNaN(date.getTime())) {
                    return vscode.l10n.t('inlineEdit.invalidDate');
                }
                return null;
            }
        });

        if (newDate !== undefined && newDate !== currentDate) {
            try {
                await this.mcpClient.callTool('gorev_duzenle', {
                    id: task.id,
                    son_tarih: newDate || null
                });

                vscode.window.showInformationMessage(
                    newDate ? vscode.l10n.t('inlineEdit.dueDateUpdated') : vscode.l10n.t('inlineEdit.dueDateRemoved')
                );
                Logger.info(`Task ${task.id} due date updated to: ${newDate || 'none'}`);
            } catch (error) {
                vscode.window.showErrorMessage(vscode.l10n.t('inlineEdit.dateUpdateFailed', error));
                Logger.error('Failed to update task due date:', error);
            }
        }
    }

    /**
     * Detaylı düzenleme dialog'u
     */
    async showDetailedEdit(task: Gorev): Promise<void> {
        const options = [
            { label: vscode.l10n.t('inlineEdit.editTitle'), action: 'title' },
            { label: vscode.l10n.t('inlineEdit.editDescription'), action: 'description' },
            { label: vscode.l10n.t('inlineEdit.changeStatus'), action: 'status' },
            { label: vscode.l10n.t('inlineEdit.changePriority'), action: 'priority' },
            { label: vscode.l10n.t('inlineEdit.changeDueDate'), action: 'dueDate' },
            { label: vscode.l10n.t('inlineEdit.editTags'), action: 'tags' }
        ];

        const selected = await vscode.window.showQuickPick(options, {
            placeHolder: vscode.l10n.t('inlineEdit.whatToEdit')
        });

        if (!selected) return;

        switch (selected.action) {
            case 'title':
                await this.startEdit({ task });
                break;
            case 'description':
                await this.editDescription(task);
                break;
            case 'status':
                await this.quickStatusChange(task);
                break;
            case 'priority':
                await this.quickPriorityChange(task);
                break;
            case 'dueDate':
                await this.quickDateChange(task);
                break;
            case 'tags':
                await this.editTags(task);
                break;
        }
    }

    /**
     * Açıklama düzenleme
     */
    private async editDescription(task: Gorev): Promise<void> {
        const newDescription = await vscode.window.showInputBox({
            prompt: vscode.l10n.t('inlineEdit.editTaskDescription'),
            value: task.aciklama || '',
            placeHolder: vscode.l10n.t('inlineEdit.descriptionPlaceholder')
        });

        if (newDescription !== undefined && newDescription !== task.aciklama) {
            try {
                await this.mcpClient.callTool('gorev_duzenle', {
                    id: task.id,
                    aciklama: newDescription
                });

                vscode.window.showInformationMessage(vscode.l10n.t('inlineEdit.descriptionUpdated'));
                Logger.info(`Task ${task.id} description updated`);
            } catch (error) {
                vscode.window.showErrorMessage(vscode.l10n.t('inlineEdit.descriptionUpdateFailed', error));
                Logger.error('Failed to update task description:', error);
            }
        }
    }

    /**
     * Etiket düzenleme
     */
    private async editTags(task: Gorev): Promise<void> {
        const currentTags = task.etiketler?.join(', ') || '';
        
        const newTags = await vscode.window.showInputBox({
            prompt: vscode.l10n.t('inlineEdit.editTagsPrompt'),
            value: currentTags,
            placeHolder: 'bug, frontend, urgent',
            validateInput: (value) => {
                if (value) {
                    const tags = value.split(',').map(t => t.trim());
                    for (const tag of tags) {
                        if (tag.length > 50) {
                            return vscode.l10n.t('inlineEdit.tagsTooLong');
                        }
                        if (!/^[a-zA-Z0-9_-]+$/.test(tag)) {
                            return vscode.l10n.t('inlineEdit.tagsInvalidChars');
                        }
                    }
                }
                return null;
            }
        });

        if (newTags !== undefined && newTags !== currentTags) {
            // Etiket güncelleme MCP tool'u henüz yok, bu yüzden şimdilik bilgi mesajı
            vscode.window.showInformationMessage(
                vscode.l10n.t('inlineEdit.tagsFeatureComingSoon')
            );
        }
    }

    /**
     * Düzenleme modunda mı kontrolü
     */
    isEditing(): boolean {
        return this.editingItem !== null;
    }

    /**
     * Mevcut düzenlemeyi iptal et
     */
    cancelEdit(): void {
        this.editingItem = null;
        this.originalLabel = null;
    }
}