import * as vscode from 'vscode';
import { 
    DragDataType, 
    TaskDragData, 
    TasksDragData,
    DropTargetType,
    DropTarget,
    DragDropResult,
    DragDropConfig,
    DropZoneVisual
} from '../utils/dragDropTypes';
import { Gorev, GorevDurum, GorevOncelik } from '../models/gorev';
import { GroupingStrategy } from '../models/treeModels';
import { MCPClient } from '../mcp/client';
import { Logger } from '../utils/logger';

/**
 * VS Code TreeView için Drag & Drop Controller
 */
export class DragDropController implements vscode.TreeDragAndDropController<any> {
    dropMimeTypes = [DragDataType.Task, DragDataType.Tasks];
    dragMimeTypes = [DragDataType.Task, DragDataType.Tasks];

    private config: DragDropConfig;
    private mcpClient: MCPClient;
    private currentGrouping: GroupingStrategy;

    constructor(mcpClient: MCPClient) {
        this.mcpClient = mcpClient;
        this.config = this.loadConfig();
        this.currentGrouping = GroupingStrategy.ByStatus;
    }

    /**
     * Drag başladığında çağrılır
     */
    public async handleDrag(
        source: readonly any[],
        dataTransfer: vscode.DataTransfer,
        token: vscode.CancellationToken
    ): Promise<void> {
        // Tek görev mi yoksa çoklu görev mi?
        if (source.length === 1 && source[0].task) {
            // Tek görev drag
            const dragData: TaskDragData = {
                type: DragDataType.Task,
                task: source[0].task,
                sourceGroupKey: source[0].parent?.groupKey
            };
            
            dataTransfer.set(
                DragDataType.Task,
                new vscode.DataTransferItem(dragData)
            );
            
            Logger.debug(`Dragging task: ${dragData.task.baslik}`);
        } else if (source.length > 1) {
            // Çoklu görev drag
            const tasks = source
                .filter(item => item.task)
                .map(item => item.task);
            
            if (tasks.length > 0) {
                const dragData: TasksDragData = {
                    type: DragDataType.Tasks,
                    tasks,
                    sourceGroupKey: source[0].parent?.groupKey
                };
                
                dataTransfer.set(
                    DragDataType.Tasks,
                    new vscode.DataTransferItem(dragData)
                );
                
                Logger.debug(`Dragging ${tasks.length} tasks`);
            }
        }
    }

    /**
     * Drop yapıldığında çağrılır
     */
    public async handleDrop(
        target: any | undefined,
        dataTransfer: vscode.DataTransfer,
        token: vscode.CancellationToken
    ): Promise<void> {
        // Drop hedefini belirle
        const dropTarget = this.identifyDropTarget(target);
        if (!dropTarget) {
            Logger.warn('Invalid drop target');
            return;
        }

        // Drag edilen veriyi al
        const taskData = dataTransfer.get(DragDataType.Task);
        const tasksData = dataTransfer.get(DragDataType.Tasks);

        if (taskData) {
            // Tek görev drop
            const dragData = taskData.value as TaskDragData;
            await this.handleTaskDrop(dragData.task, dropTarget, dragData.sourceGroupKey);
        } else if (tasksData) {
            // Çoklu görev drop
            const dragData = tasksData.value as TasksDragData;
            await this.handleMultipleTasksDrop(dragData.tasks, dropTarget, dragData.sourceGroupKey);
        }
    }

    /**
     * Drop hedefini belirler
     */
    private identifyDropTarget(target: any): DropTarget | null {
        if (!target) {
            return {
                type: DropTargetType.EmptyArea
            };
        }

        // Grup üzerine drop
        if (target.groupKey) {
            switch (this.currentGrouping) {
                case GroupingStrategy.ByStatus:
                    return {
                        type: DropTargetType.StatusGroup,
                        groupKey: target.groupKey,
                        newStatus: this.mapGroupKeyToStatus(target.groupKey)
                    };
                case GroupingStrategy.ByPriority:
                    return {
                        type: DropTargetType.PriorityGroup,
                        groupKey: target.groupKey,
                        newPriority: this.mapGroupKeyToPriority(target.groupKey)
                    };
                case GroupingStrategy.ByProject:
                    return {
                        type: DropTargetType.ProjectGroup,
                        groupKey: target.groupKey,
                        newProjectId: target.groupKey !== 'no-project' ? target.groupKey : undefined
                    };
                default:
                    return null;
            }
        }

        // Görev üzerine drop (bağımlılık oluşturma için)
        if (target.task) {
            return {
                type: DropTargetType.Task,
                targetTask: target.task
            };
        }

        return null;
    }

    /**
     * Tek görev drop işlemi
     */
    private async handleTaskDrop(
        task: Gorev,
        dropTarget: DropTarget,
        sourceGroupKey?: string
    ): Promise<void> {
        try {
            switch (dropTarget.type) {
                case DropTargetType.StatusGroup:
                    if (this.config.allowStatusChange && dropTarget.newStatus) {
                        await this.updateTaskStatus(task, dropTarget.newStatus);
                    }
                    break;

                case DropTargetType.PriorityGroup:
                    if (this.config.allowPriorityChange && dropTarget.newPriority) {
                        await this.updateTaskPriority(task, dropTarget.newPriority);
                    }
                    break;

                case DropTargetType.ProjectGroup:
                    if (this.config.allowProjectMove && dropTarget.newProjectId !== undefined) {
                        await this.moveTaskToProject(task, dropTarget.newProjectId);
                    }
                    break;

                case DropTargetType.Task:
                    if (this.config.allowDependencyCreate && dropTarget.targetTask) {
                        await this.createDependency(task, dropTarget.targetTask);
                    }
                    break;
            }

            // UI'da başarı göstergesi
            if (this.config.animateOnDrop) {
                vscode.window.showInformationMessage(
                    `✅ ${task.baslik} başarıyla taşındı`
                );
            }
        } catch (error) {
            Logger.error('Drop operation failed:', error);
            vscode.window.showErrorMessage(`Drop işlemi başarısız: ${error}`);
        }
    }

    /**
     * Çoklu görev drop işlemi
     */
    private async handleMultipleTasksDrop(
        tasks: Gorev[],
        dropTarget: DropTarget,
        sourceGroupKey?: string
    ): Promise<void> {
        const operations: Promise<void>[] = [];

        for (const task of tasks) {
            switch (dropTarget.type) {
                case DropTargetType.StatusGroup:
                    if (this.config.allowStatusChange && dropTarget.newStatus) {
                        operations.push(this.updateTaskStatus(task, dropTarget.newStatus));
                    }
                    break;

                case DropTargetType.PriorityGroup:
                    if (this.config.allowPriorityChange && dropTarget.newPriority) {
                        operations.push(this.updateTaskPriority(task, dropTarget.newPriority));
                    }
                    break;

                case DropTargetType.ProjectGroup:
                    if (this.config.allowProjectMove && dropTarget.newProjectId !== undefined) {
                        operations.push(this.moveTaskToProject(task, dropTarget.newProjectId));
                    }
                    break;
            }
        }

        if (operations.length > 0) {
            try {
                await vscode.window.withProgress(
                    {
                        location: vscode.ProgressLocation.Notification,
                        title: `${tasks.length} görev taşınıyor...`,
                        cancellable: false
                    },
                    async (progress) => {
                        let completed = 0;
                        for (const op of operations) {
                            await op;
                            completed++;
                            progress.report({
                                increment: (completed / operations.length) * 100
                            });
                        }
                    }
                );

                vscode.window.showInformationMessage(
                    `✅ ${tasks.length} görev başarıyla taşındı`
                );
            } catch (error) {
                Logger.error('Multiple drop operation failed:', error);
                vscode.window.showErrorMessage(`İşlem başarısız: ${error}`);
            }
        }
    }

    /**
     * Görev durumunu günceller
     */
    private async updateTaskStatus(task: Gorev, newStatus: GorevDurum): Promise<void> {
        if (task.durum === newStatus) {
            return; // Zaten aynı durumda
        }

        await this.mcpClient.callTool('gorev_guncelle', {
            id: task.id,
            durum: newStatus
        });

        Logger.info(`Task ${task.id} status updated to ${newStatus}`);
    }

    /**
     * Görev önceliğini günceller
     */
    private async updateTaskPriority(task: Gorev, newPriority: GorevOncelik): Promise<void> {
        if (task.oncelik === newPriority) {
            return; // Zaten aynı öncelikte
        }

        await this.mcpClient.callTool('gorev_duzenle', {
            id: task.id,
            oncelik: newPriority
        });

        Logger.info(`Task ${task.id} priority updated to ${newPriority}`);
    }

    /**
     * Görevi başka projeye taşır
     */
    private async moveTaskToProject(task: Gorev, newProjectId: string | undefined): Promise<void> {
        const currentProjectId = task.proje_id || '';
        if (currentProjectId === newProjectId) {
            return; // Zaten aynı projede
        }

        await this.mcpClient.callTool('gorev_duzenle', {
            id: task.id,
            proje_id: newProjectId || ''
        });

        Logger.info(`Task ${task.id} moved to project ${newProjectId || 'none'}`);
    }

    /**
     * İki görev arasında bağımlılık oluşturur
     */
    private async createDependency(sourceTask: Gorev, targetTask: Gorev): Promise<void> {
        if (sourceTask.id === targetTask.id) {
            vscode.window.showWarningMessage('Bir görev kendisine bağımlı olamaz');
            return;
        }

        // Circular dependency kontrolü yapılabilir
        
        const result = await vscode.window.showQuickPick(
            [
                { label: 'Blocks', value: 'blocks', description: `"${sourceTask.baslik}" blocks "${targetTask.baslik}"` },
                { label: 'Depends on', value: 'depends_on', description: `"${sourceTask.baslik}" depends on "${targetTask.baslik}"` },
                { label: 'Related to', value: 'related', description: `"${sourceTask.baslik}" is related to "${targetTask.baslik}"` }
            ],
            {
                placeHolder: 'Select dependency type'
            }
        );

        if (result) {
            await this.mcpClient.callTool('gorev_bagimlilik_ekle', {
                kaynak_id: sourceTask.id,
                hedef_id: targetTask.id,
                baglanti_tipi: result.value
            });

            vscode.window.showInformationMessage(
                `✅ Bağımlılık oluşturuldu: ${sourceTask.baslik} ${result.value} ${targetTask.baslik}`
            );
        }
    }

    /**
     * Grup anahtarını duruma çevirir
     */
    private mapGroupKeyToStatus(groupKey: string): GorevDurum | undefined {
        const map: { [key: string]: GorevDurum } = {
            'beklemede': GorevDurum.Beklemede,
            'devam_ediyor': GorevDurum.DevamEdiyor,
            'tamamlandi': GorevDurum.Tamamlandi
        };
        return map[groupKey];
    }

    /**
     * Grup anahtarını önceliğe çevirir
     */
    private mapGroupKeyToPriority(groupKey: string): GorevOncelik | undefined {
        const map: { [key: string]: GorevOncelik } = {
            'yuksek': GorevOncelik.Yuksek,
            'orta': GorevOncelik.Orta,
            'dusuk': GorevOncelik.Dusuk
        };
        return map[groupKey];
    }

    /**
     * Konfigürasyonu yükler
     */
    private loadConfig(): DragDropConfig {
        const config = vscode.workspace.getConfiguration('gorev.dragDrop');
        return {
            allowTaskMove: config.get('allowTaskMove', true),
            allowStatusChange: config.get('allowStatusChange', true),
            allowPriorityChange: config.get('allowPriorityChange', true),
            allowProjectMove: config.get('allowProjectMove', true),
            allowDependencyCreate: config.get('allowDependencyCreate', true),
            showDropIndicator: config.get('showDropIndicator', true),
            animateOnDrop: config.get('animateOnDrop', true)
        };
    }

    /**
     * Mevcut gruplama stratejisini günceller
     */
    public setGroupingStrategy(strategy: GroupingStrategy): void {
        this.currentGrouping = strategy;
    }

    /**
     * Drop yapılabilir mi kontrolü
     */
    public canDrop(target: any, dataTransfer: vscode.DataTransfer): boolean {
        // Basit kontroller
        const hasTaskData = dataTransfer.get(DragDataType.Task) !== undefined;
        const hasTasksData = dataTransfer.get(DragDataType.Tasks) !== undefined;
        
        if (!hasTaskData && !hasTasksData) {
            return false;
        }

        const dropTarget = this.identifyDropTarget(target);
        if (!dropTarget) {
            return false;
        }

        // Konfigürasyona göre kontrol
        switch (dropTarget.type) {
            case DropTargetType.StatusGroup:
                return this.config.allowStatusChange;
            case DropTargetType.PriorityGroup:
                return this.config.allowPriorityChange;
            case DropTargetType.ProjectGroup:
                return this.config.allowProjectMove;
            case DropTargetType.Task:
                return this.config.allowDependencyCreate;
            default:
                return false;
        }
    }
}