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
            prompt: 'Görev başlığını düzenle',
            value: item.task.baslik,
            validateInput: (value) => {
                if (!value || value.trim().length === 0) {
                    return 'Görev başlığı boş olamaz';
                }
                if (value.length > 200) {
                    return 'Görev başlığı 200 karakterden uzun olamaz';
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

            vscode.window.showInformationMessage('Görev başlığı güncellendi');
            Logger.info(`Task ${task.id} title updated to: ${newTitle}`);
        } catch (error) {
            vscode.window.showErrorMessage(`Güncelleme başarısız: ${error}`);
            Logger.error('Failed to update task title:', error);
        }
    }

    /**
     * Hızlı durum değiştirme menüsü
     */
    async quickStatusChange(task: Gorev): Promise<void> {
        const items = [
            { 
                label: '$(circle-outline) Beklemede', 
                value: GorevDurum.Beklemede,
                description: task.durum === GorevDurum.Beklemede ? 'Mevcut durum' : ''
            },
            { 
                label: '$(sync~spin) Devam Ediyor', 
                value: GorevDurum.DevamEdiyor,
                description: task.durum === GorevDurum.DevamEdiyor ? 'Mevcut durum' : ''
            },
            { 
                label: '$(check) Tamamlandı', 
                value: GorevDurum.Tamamlandi,
                description: task.durum === GorevDurum.Tamamlandi ? 'Mevcut durum' : ''
            }
        ];

        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: 'Yeni durumu seçin',
            title: `"${task.baslik}" için durum değiştir`
        });

        if (selected && selected.value !== task.durum) {
            try {
                Logger.info(`[QuickStatusChange] Updating task ${task.id} from ${task.durum} to ${selected.value}`);
                
                const result = await this.mcpClient.callTool('gorev_guncelle', {
                    id: task.id,
                    durum: selected.value
                });
                
                Logger.info(`[QuickStatusChange] MCP response:`, JSON.stringify(result));

                vscode.window.showInformationMessage('Görev durumu güncellendi');
                Logger.info(`[QuickStatusChange] Task ${task.id} status updated to: ${selected.value}`);
                
                // Force a command execution to refresh all trees
                await vscode.commands.executeCommand('gorev.refreshTasks');
            } catch (error) {
                vscode.window.showErrorMessage(`Durum güncellemesi başarısız: ${error}`);
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
                label: '$(flame) Yüksek Öncelik', 
                value: GorevOncelik.Yuksek,
                description: task.oncelik === GorevOncelik.Yuksek ? 'Mevcut öncelik' : ''
            },
            { 
                label: '$(dash) Orta Öncelik', 
                value: GorevOncelik.Orta,
                description: task.oncelik === GorevOncelik.Orta ? 'Mevcut öncelik' : ''
            },
            { 
                label: '$(arrow-down) Düşük Öncelik', 
                value: GorevOncelik.Dusuk,
                description: task.oncelik === GorevOncelik.Dusuk ? 'Mevcut öncelik' : ''
            }
        ];

        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: 'Yeni önceliği seçin',
            title: `"${task.baslik}" için öncelik değiştir`
        });

        if (selected && selected.value !== task.oncelik) {
            try {
                await this.mcpClient.callTool('gorev_duzenle', {
                    id: task.id,
                    oncelik: selected.value
                });

                vscode.window.showInformationMessage('Görev önceliği güncellendi');
                Logger.info(`Task ${task.id} priority updated to: ${selected.value}`);
            } catch (error) {
                vscode.window.showErrorMessage(`Öncelik güncellemesi başarısız: ${error}`);
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
            prompt: 'Son tarihi girin (YYYY-MM-DD)',
            value: currentDate,
            placeHolder: '2024-12-31',
            validateInput: (value) => {
                if (!value) {
                    return null; // Boş bırakılabilir
                }
                const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
                if (!dateRegex.test(value)) {
                    return 'Geçersiz tarih formatı. YYYY-MM-DD kullanın';
                }
                const date = new Date(value);
                if (isNaN(date.getTime())) {
                    return 'Geçersiz tarih';
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
                    newDate ? 'Son tarih güncellendi' : 'Son tarih kaldırıldı'
                );
                Logger.info(`Task ${task.id} due date updated to: ${newDate || 'none'}`);
            } catch (error) {
                vscode.window.showErrorMessage(`Tarih güncellemesi başarısız: ${error}`);
                Logger.error('Failed to update task due date:', error);
            }
        }
    }

    /**
     * Detaylı düzenleme dialog'u
     */
    async showDetailedEdit(task: Gorev): Promise<void> {
        const options = [
            { label: '$(edit) Başlık Düzenle', action: 'title' },
            { label: '$(note) Açıklama Düzenle', action: 'description' },
            { label: '$(circle-outline) Durum Değiştir', action: 'status' },
            { label: '$(flame) Öncelik Değiştir', action: 'priority' },
            { label: '$(calendar) Son Tarih Değiştir', action: 'dueDate' },
            { label: '$(tag) Etiketleri Düzenle', action: 'tags' }
        ];

        const selected = await vscode.window.showQuickPick(options, {
            placeHolder: 'Ne düzenlemek istiyorsunuz?'
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
            prompt: 'Görev açıklamasını düzenle',
            value: task.aciklama || '',
            placeHolder: 'Görev hakkında detaylı açıklama...'
        });

        if (newDescription !== undefined && newDescription !== task.aciklama) {
            try {
                await this.mcpClient.callTool('gorev_duzenle', {
                    id: task.id,
                    aciklama: newDescription
                });

                vscode.window.showInformationMessage('Görev açıklaması güncellendi');
                Logger.info(`Task ${task.id} description updated`);
            } catch (error) {
                vscode.window.showErrorMessage(`Açıklama güncellemesi başarısız: ${error}`);
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
            prompt: 'Etiketleri düzenle (virgülle ayırın)',
            value: currentTags,
            placeHolder: 'bug, frontend, urgent',
            validateInput: (value) => {
                if (value) {
                    const tags = value.split(',').map(t => t.trim());
                    for (const tag of tags) {
                        if (tag.length > 50) {
                            return 'Etiketler 50 karakterden uzun olamaz';
                        }
                        if (!/^[a-zA-Z0-9_-]+$/.test(tag)) {
                            return 'Etiketler sadece harf, rakam, tire ve alt çizgi içerebilir';
                        }
                    }
                }
                return null;
            }
        });

        if (newTags !== undefined && newTags !== currentTags) {
            // Etiket güncelleme MCP tool'u henüz yok, bu yüzden şimdilik bilgi mesajı
            vscode.window.showInformationMessage(
                'Etiket güncelleme özelliği yakında eklenecek'
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