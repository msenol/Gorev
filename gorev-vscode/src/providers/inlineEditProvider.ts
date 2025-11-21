import * as vscode from 'vscode';
import { ApiClient } from '../api/client';
import { Gorev, GorevDurum, GorevOncelik } from '../models/gorev';
import { Logger } from '../utils/logger';
import { t } from '../utils/l10n';

/**
 * TreeView item'ları için inline düzenleme sağlayıcı
 */
export class InlineEditProvider {
    private editingItem: { task?: Gorev } | null = null;
    private originalLabel: string | null = null;

    constructor(private apiClient: ApiClient) {}

    /**
     * Inline düzenlemeyi başlatır
     */
    async startEdit(item: { task?: Gorev }): Promise<void> {
        if (!item || !item.task) {
            return;
        }

        this.editingItem = item;
        this.originalLabel = item.task.baslik;

        const newTitle = await vscode.window.showInputBox({
            prompt: t('inlineEdit.editTaskTitle'),
            value: item.task.baslik,
            validateInput: (value) => {
                if (!value || value.trim().length === 0) {
                    return t('inlineEdit.taskTitleEmpty');
                }
                if (value.length > 200) {
                    return t('inlineEdit.taskTitleTooLong');
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
            await this.apiClient.updateTask(task.id, {
                baslik: newTitle
            });

            vscode.window.showInformationMessage(t('inlineEdit.taskTitleUpdated'));
            Logger.info(`Task ${task.id} title updated to: ${newTitle}`);
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : String(error);
            vscode.window.showErrorMessage(t('inlineEdit.updateFailed', errorMessage));
            Logger.error('Failed to update task title:', error);
        }
    }

    /**
     * Hızlı durum değiştirme menüsü
     */
    async quickStatusChange(task: Gorev): Promise<void> {
        const items = [
            { 
                label: t('inlineEdit.pending'), 
                value: GorevDurum.Beklemede,
                description: task.durum === GorevDurum.Beklemede ? t('inlineEdit.currentStatus') : ''
            },
            { 
                label: t('inlineEdit.inProgress'), 
                value: GorevDurum.DevamEdiyor,
                description: task.durum === GorevDurum.DevamEdiyor ? t('inlineEdit.currentStatus') : ''
            },
            { 
                label: t('inlineEdit.completed'), 
                value: GorevDurum.Tamamlandi,
                description: task.durum === GorevDurum.Tamamlandi ? t('inlineEdit.currentStatus') : ''
            }
        ];

        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: t('inlineEdit.selectNewStatus'),
            title: t('inlineEdit.changeStatusFor', task.baslik)
        });

        if (selected && selected.value !== task.durum) {
            try {
                Logger.info(`[QuickStatusChange] Updating task ${task.id} from ${task.durum} to ${selected.value}`);

                const result = await this.apiClient.updateTask(task.id, {
                    durum: selected.value
                });

                Logger.info(`[QuickStatusChange] API response:`, JSON.stringify(result));

                vscode.window.showInformationMessage(t('inlineEdit.taskStatusUpdated'));
                Logger.info(`[QuickStatusChange] Task ${task.id} status updated to: ${selected.value}`);
                
                // Force a command execution to refresh all trees
                await vscode.commands.executeCommand('gorev.refreshTasks');
            } catch (error) {
                const errorMessage = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(t('inlineEdit.statusUpdateFailed', errorMessage));
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
                label: t('inlineEdit.highPriority'), 
                value: GorevOncelik.Yuksek,
                description: task.oncelik === GorevOncelik.Yuksek ? t('inlineEdit.currentPriority') : ''
            },
            { 
                label: t('inlineEdit.mediumPriority'), 
                value: GorevOncelik.Orta,
                description: task.oncelik === GorevOncelik.Orta ? t('inlineEdit.currentPriority') : ''
            },
            { 
                label: t('inlineEdit.lowPriority'), 
                value: GorevOncelik.Dusuk,
                description: task.oncelik === GorevOncelik.Dusuk ? t('inlineEdit.currentPriority') : ''
            }
        ];

        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: t('inlineEdit.selectNewPriority'),
            title: t('inlineEdit.changePriorityFor', task.baslik)
        });

        if (selected && selected.value !== task.oncelik) {
            try {
                await this.apiClient.updateTask(task.id, {
                    oncelik: selected.value
                });

                vscode.window.showInformationMessage(t('inlineEdit.taskPriorityUpdated'));
                Logger.info(`Task ${task.id} priority updated to: ${selected.value}`);
            } catch (error) {
                const errorMessage = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(t('inlineEdit.priorityUpdateFailed', errorMessage));
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
            prompt: t('inlineEdit.enterDueDate'),
            value: currentDate,
            placeHolder: '2024-12-31',
            validateInput: (value) => {
                if (!value) {
                    return null; // Boş bırakılabilir
                }
                const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
                if (!dateRegex.test(value)) {
                    return t('inlineEdit.invalidDateFormat');
                }
                const date = new Date(value);
                if (isNaN(date.getTime())) {
                    return t('inlineEdit.invalidDate');
                }
                return null;
            }
        });

        if (newDate !== undefined && newDate !== currentDate) {
            try {
                await this.apiClient.updateTask(task.id, {
                    son_tarih: newDate || undefined
                });

                vscode.window.showInformationMessage(
                    newDate ? t('inlineEdit.dueDateUpdated') : t('inlineEdit.dueDateRemoved')
                );
                Logger.info(`Task ${task.id} due date updated to: ${newDate || 'none'}`);
            } catch (error) {
                const errorMessage = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(t('inlineEdit.dateUpdateFailed', errorMessage));
                Logger.error('Failed to update task due date:', error);
            }
        }
    }

    /**
     * Detaylı düzenleme dialog'u
     */
    async showDetailedEdit(task: Gorev): Promise<void> {
        const options = [
            { label: t('inlineEdit.editTitle'), action: 'title' },
            { label: t('inlineEdit.editDescription'), action: 'description' },
            { label: t('inlineEdit.changeStatus'), action: 'status' },
            { label: t('inlineEdit.changePriority'), action: 'priority' },
            { label: t('inlineEdit.changeDueDate'), action: 'dueDate' },
            { label: t('inlineEdit.editTags'), action: 'tags' }
        ];

        const selected = await vscode.window.showQuickPick(options, {
            placeHolder: t('inlineEdit.whatToEdit')
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
            prompt: t('inlineEdit.editTaskDescription'),
            value: task.aciklama || '',
            placeHolder: t('inlineEdit.descriptionPlaceholder')
        });

        if (newDescription !== undefined && newDescription !== task.aciklama) {
            try {
                await this.apiClient.updateTask(task.id, {
                    aciklama: newDescription
                });

                vscode.window.showInformationMessage(t('inlineEdit.descriptionUpdated'));
                Logger.info(`Task ${task.id} description updated`);
            } catch (error) {
                const errorMessage = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(t('inlineEdit.descriptionUpdateFailed', errorMessage));
                Logger.error('Failed to update task description:', error);
            }
        }
    }

    /**
     * Etiket düzenleme
     */
    private async editTags(task: Gorev): Promise<void> {
        const currentTags = task.etiketler?.map(t => t.isim).join(', ') || '';
        
        const newTags = await vscode.window.showInputBox({
            prompt: t('inlineEdit.editTagsPrompt'),
            value: currentTags,
            placeHolder: 'bug, frontend, urgent',
            validateInput: (value) => {
                if (value) {
                    const tags = value.split(',').map(t => t.trim());
                    for (const tag of tags) {
                        if (tag.length > 50) {
                            return t('inlineEdit.tagsTooLong');
                        }
                        if (!/^[a-zA-Z0-9_-]+$/.test(tag)) {
                            return t('inlineEdit.tagsInvalidChars');
                        }
                    }
                }
                return null;
            }
        });

        if (newTags !== undefined && newTags !== currentTags) {
            // Etiket güncelleme MCP tool'u henüz yok, bu yüzden şimdilik bilgi mesajı
            vscode.window.showInformationMessage(
                t('inlineEdit.tagsFeatureComingSoon')
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
