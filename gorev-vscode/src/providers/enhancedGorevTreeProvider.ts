import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { Gorev, GorevDurum, GorevOncelik } from '../models/gorev';
import { 
    EnhancedTreeItem, 
    GroupTreeItem, 
    TaskTreeItem, 
    EmptyTreeItem,
    TreeItemType,
    GroupingStrategy,
    SortingCriteria,
    TreeViewConfig,
    TaskFilter,
    TaskSelection,
    TreeViewEvents,
    TreeViewUtils
} from '../models/treeModels';
import { GroupingStrategyProvider } from './groupingStrategy';
import { DragDropController } from './dragDropController';
import { TaskDecorationProvider } from './decorationProvider';
import { ICONS, COLORS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';
import { MarkdownParser } from '../utils/markdownParser';

/**
 * Gelişmiş görev TreeView provider'ı
 */
export class EnhancedGorevTreeProvider implements vscode.TreeDataProvider<EnhancedTreeViewItem>, vscode.TreeDragAndDropController<EnhancedTreeViewItem> {
    private _onDidChangeTreeData = new vscode.EventEmitter<EnhancedTreeViewItem | undefined | null | void>();
    readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

    private tasks: Gorev[] = [];
    private filteredTasks: Gorev[] = [];
    private projects: Map<string, any> = new Map(); // Proje ID -> Proje bilgisi
    private config: TreeViewConfig;
    private selection: TaskSelection;
    private events: TreeViewEvents = {};
    
    // Drag & Drop
    public readonly dragDropController: DragDropController;
    readonly dropMimeTypes: readonly string[];
    readonly dragMimeTypes: readonly string[];
    
    // Decoration provider
    private decorationProvider: TaskDecorationProvider;

    constructor(private mcpClient: MCPClient) {
        // Varsayılan konfigürasyon
        this.config = {
            grouping: GroupingStrategy.ByStatus,
            sorting: SortingCriteria.Priority,
            sortAscending: false,
            showCompleted: true,
            showEmptyGroups: false,
            expandedGroups: GroupingStrategyProvider.getDefaultExpandedGroups(GroupingStrategy.ByStatus),
            filters: {}
        };

        // Boş selection
        this.selection = {
            selectedTasks: new Set<string>()
        };

        // Drag & Drop controller
        this.dragDropController = new DragDropController(mcpClient);
        this.dropMimeTypes = this.dragDropController.dropMimeTypes;
        this.dragMimeTypes = this.dragDropController.dragMimeTypes;
        
        // Decoration provider
        this.decorationProvider = new TaskDecorationProvider();

        // Konfigürasyon değişikliklerini dinle
        this.loadConfiguration();
        vscode.workspace.onDidChangeConfiguration(e => {
            if (e.affectsConfiguration('gorev.treeView')) {
                this.loadConfiguration();
                this.refresh();
            }
        });
    }

    /**
     * TreeView konfigürasyonunu yükler
     */
    private loadConfiguration(): void {
        const config = vscode.workspace.getConfiguration('gorev.treeView');
        
        this.config.grouping = config.get('grouping', GroupingStrategy.ByStatus) as GroupingStrategy;
        this.config.sorting = config.get('sorting', SortingCriteria.Priority) as SortingCriteria;
        this.config.sortAscending = config.get('sortAscending', false);
        this.config.showCompleted = config.get('showCompleted', true);
        this.config.showEmptyGroups = config.get('showEmptyGroups', false);
    }

    /**
     * TreeView item'ını döndürür
     */
    getTreeItem(element: EnhancedTreeViewItem): vscode.TreeItem {
        return element;
    }

    /**
     * Alt elemanları döndürür
     */
    async getChildren(element?: EnhancedTreeViewItem): Promise<EnhancedTreeViewItem[]> {
        Logger.debug('[EnhancedGorevTreeProvider] getChildren called with element:', element);
        Logger.debug('[EnhancedGorevTreeProvider] Current tasks count:', this.tasks.length);
        Logger.debug('[EnhancedGorevTreeProvider] Current filtered tasks count:', this.filteredTasks.length);
        
        if (!this.mcpClient.isConnected()) {
            Logger.warn('[EnhancedGorevTreeProvider] MCP client not connected');
            return [new EmptyTreeViewItem('MCP sunucusuna bağlı değil')];
        }

        // Root level
        if (!element) {
            Logger.debug('[EnhancedGorevTreeProvider] Loading root items...');
            try {
                await this.loadTasks();
                const items = this.createRootItems();
                Logger.debug('[EnhancedGorevTreeProvider] Returning', items.length, 'root items');
                return items;
            } catch (error) {
                Logger.error('Failed to load tasks:', error);
                return [new EmptyTreeViewItem('Görevler yüklenemedi')];
            }
        }

        // Grup altındaki görevler
        if (element instanceof GroupTreeViewItem) {
            Logger.debug('[EnhancedGorevTreeProvider] Loading children for group:', element.groupKey);
            return this.createTaskItems(element);
        }

        // Task altındaki alt görevler
        if (element instanceof TaskTreeViewItem && element.task.alt_gorevler && element.task.alt_gorevler.length > 0) {
            Logger.debug('[EnhancedGorevTreeProvider] Loading subtasks for task:', element.task.id);
            return this.createSubtaskItems(element.task);
        }

        return [];
    }

    /**
     * Root level item'larını oluşturur
     */
    private createRootItems(): EnhancedTreeViewItem[] {
        // Filtreleme uygula
        this.filteredTasks = TreeViewUtils.filterTasks(this.tasks, this.config.filters);
        Logger.debug('[EnhancedGorevTreeProvider] After filtering:', this.filteredTasks.length, 'tasks');

        // Tamamlanmış görevleri gizle
        if (!this.config.showCompleted) {
            this.filteredTasks = this.filteredTasks.filter(
                task => task.durum !== GorevDurum.Tamamlandi
            );
            Logger.debug('[EnhancedGorevTreeProvider] After hiding completed:', this.filteredTasks.length, 'tasks');
        }

        if (this.filteredTasks.length === 0) {
            Logger.debug('[EnhancedGorevTreeProvider] No tasks to show, returning empty message');
            return [new EmptyTreeViewItem(this.getEmptyMessage())];
        }

        // Gruplama yoksa direkt görevleri göster
        Logger.debug('[EnhancedGorevTreeProvider] Grouping strategy:', this.config.grouping);
        if (this.config.grouping === GroupingStrategy.None) {
            const sortedTasks = TreeViewUtils.sortTasks(
                this.filteredTasks, 
                this.config.sorting, 
                this.config.sortAscending
            );
            // Sadece root level görevleri göster (parent_id olmayan)
            const rootTasks = sortedTasks.filter(task => !task.parent_id);
            return rootTasks.map(task => new TaskTreeViewItem(task, this.selection));
        }

        // Görevleri grupla
        const groups = TreeViewUtils.groupTasks(this.filteredTasks, this.config.grouping);
        Logger.debug('[EnhancedGorevTreeProvider] Created groups:', groups.size, 'groups');
        Logger.debug('[EnhancedGorevTreeProvider] Group keys:', Array.from(groups.keys()));
        
        const groupItems: EnhancedTreeViewItem[] = [];

        // Grupları sırala
        const sortedGroupKeys = Array.from(groups.keys()).sort((a, b) => 
            GroupingStrategyProvider.compareGroups(a, b, this.config.grouping)
        );
        Logger.debug('[EnhancedGorevTreeProvider] Sorted group keys:', sortedGroupKeys);

        // Grup item'larını oluştur
        for (const groupKey of sortedGroupKeys) {
            const tasksInGroup = groups.get(groupKey)!;
            Logger.debug(`[EnhancedGorevTreeProvider] Group ${groupKey} has ${tasksInGroup.length} tasks`);
            
            // Boş grupları gizle
            if (!this.config.showEmptyGroups && tasksInGroup.length === 0) {
                Logger.debug(`[EnhancedGorevTreeProvider] Skipping empty group: ${groupKey}`);
                continue;
            }

            const groupItem = new GroupTreeViewItem(
                groupKey,
                this.config.grouping,
                tasksInGroup,
                this.config.expandedGroups.has(groupKey),
                this.projects
            );

            groupItems.push(groupItem);
        }

        Logger.debug('[EnhancedGorevTreeProvider] Total group items created:', groupItems.length);
        return groupItems;
    }

    /**
     * Grup altındaki görev item'larını oluşturur
     */
    private createTaskItems(group: GroupTreeViewItem): TaskTreeViewItem[] {
        // Grup içindeki görevleri sırala
        const sortedTasks = TreeViewUtils.sortTasks(
            group.tasks,
            this.config.sorting,
            this.config.sortAscending
        );

        // Sadece root level görevleri göster (parent_id olmayan)
        const rootTasks = sortedTasks.filter(task => !task.parent_id);
        return rootTasks.map(task => {
            const item = new TaskTreeViewItem(task, this.selection, group);
            this.decorationProvider.updateTaskDecoration(task, item);
            return item;
        });
    }

    /**
     * Alt görev item'larını oluşturur
     */
    private createSubtaskItems(parentTask: Gorev): TaskTreeViewItem[] {
        if (!parentTask.alt_gorevler || parentTask.alt_gorevler.length === 0) {
            return [];
        }

        // Alt görevleri sırala
        const sortedSubtasks = TreeViewUtils.sortTasks(
            parentTask.alt_gorevler,
            this.config.sorting,
            this.config.sortAscending
        );

        return sortedSubtasks.map(task => {
            const item = new TaskTreeViewItem(task, this.selection);
            this.decorationProvider.updateTaskDecoration(task, item);
            return item;
        });
    }

    /**
     * Boş durum mesajını döndürür
     */
    private getEmptyMessage(): string {
        if (this.config.filters.searchQuery) {
            return `"${this.config.filters.searchQuery}" için sonuç bulunamadı`;
        }
        if (Object.keys(this.config.filters).length > 0) {
            return 'Filtrelere uygun görev bulunamadı';
        }
        return 'Henüz görev yok';
    }

    /**
     * Görevleri yükler
     */
    private async loadTasks(): Promise<void> {
        try {
            // Önce projeleri yükle
            Logger.debug('[EnhancedGorevTreeProvider] Loading projects...');
            const projectsResult = await this.mcpClient.callTool('proje_listele', {});
            if (projectsResult && projectsResult.content && projectsResult.content[0]) {
                const projectsText = projectsResult.content[0].text;
                const projeler = MarkdownParser.parseProjeListesi(projectsText);
                this.projects.clear();
                for (const proje of projeler) {
                    if (proje.id) {
                        this.projects.set(proje.id, proje);
                    }
                }
                Logger.debug(`[EnhancedGorevTreeProvider] Loaded ${this.projects.size} projects`);
            }
            
            // Get page size from configuration
            const pageSize = vscode.workspace.getConfiguration('gorev').get<number>('pagination.pageSize', 100);
            
            // Check if we should show all projects or just active project
            const showAllProjects = this.config.filters?.showAllProjects || false;
            
            // First get the active project to assign project_id to tasks
            let activeProjectId = '';
            if (!showAllProjects) {
                try {
                    const activeProjectResult = await this.mcpClient.callTool('aktif_proje_goster', {});
                    if (activeProjectResult && activeProjectResult.content && activeProjectResult.content[0]) {
                        const activeProjectText = activeProjectResult.content[0].text;
                        const idMatch = activeProjectText.match(/ID:\s*([a-f0-9-]+)/);
                        if (idMatch) {
                            activeProjectId = idMatch[1];
                            Logger.debug('[EnhancedGorevTreeProvider] Active project ID:', activeProjectId);
                        }
                    }
                } catch (err) {
                    Logger.warn('[EnhancedGorevTreeProvider] Failed to get active project:', err);
                }
            }
            
            Logger.debug('[EnhancedGorevTreeProvider] Calling gorev_listele with page size:', pageSize, 'showAllProjects:', showAllProjects);
            const result = await this.mcpClient.callTool('gorev_listele', {
                tum_projeler: showAllProjects,
                limit: pageSize,
                offset: 0
            });
            
            // Debug: Log raw response
            Logger.debug('[EnhancedGorevTreeProvider] Raw MCP response:', JSON.stringify(result, null, 2));
            
            if (result && result.content && result.content[0]) {
                const responseText = result.content[0].text;
                Logger.debug('[EnhancedGorevTreeProvider] Content text length:', responseText.length);
                Logger.debug('[EnhancedGorevTreeProvider] Content text (first 500 chars):', responseText.substring(0, 500));
                
                // Parse the markdown content to extract tasks
                this.tasks = MarkdownParser.parseGorevListesi(responseText);
                
                // Debug: Log parsed tasks
                Logger.info('[EnhancedGorevTreeProvider] Parsed tasks count:', this.tasks.length);
                
                // If tasks don't have project_id and we're showing active project tasks only, assign it
                if (!showAllProjects && activeProjectId && this.tasks.length > 0) {
                    for (const task of this.tasks) {
                        if (!task.proje_id || task.proje_id === '') {
                            task.proje_id = activeProjectId;
                            Logger.debug(`[EnhancedGorevTreeProvider] Assigned project_id ${activeProjectId} to task: ${task.baslik}`);
                        }
                    }
                }
                
                if (this.tasks.length > 0) {
                    Logger.debug('[EnhancedGorevTreeProvider] First task:', JSON.stringify(this.tasks[0], null, 2));
                    Logger.debug('[EnhancedGorevTreeProvider] Tasks with project_id:', this.tasks.filter(t => t.proje_id).length);
                    Logger.debug('[EnhancedGorevTreeProvider] Tasks without project_id:', this.tasks.filter(t => !t.proje_id).length);
                } else {
                    Logger.warn('[EnhancedGorevTreeProvider] No tasks parsed from response');
                    // Log a few lines to debug
                    const lines = responseText.split('\n').slice(0, 10);
                    Logger.debug('[EnhancedGorevTreeProvider] First 10 lines of response:', lines);
                }
                
                // IMPORTANT: Set filtered tasks after loading
                this.filteredTasks = [...this.tasks];
                Logger.debug('[EnhancedGorevTreeProvider] After filtering:', this.filteredTasks.length, 'tasks');
            } else {
                Logger.error('[EnhancedGorevTreeProvider] Invalid MCP response structure:', JSON.stringify(result, null, 2));
                this.tasks = [];
                this.filteredTasks = [];
            }
        } catch (error) {
            Logger.error('Failed to load tasks:', error);
            throw error;
        }
    }

    /**
     * TreeView'ı yeniler
     */
    async refresh(): Promise<void> {
        Logger.info('[EnhancedGorevTreeProvider] Refreshing tree view...');
        try {
            // Clear cached data to force full reload
            this.tasks = [];
            this.filteredTasks = [];
            
            await this.loadTasks();
            
            // Fire change event with undefined to refresh entire tree
            this._onDidChangeTreeData.fire(undefined);
            
            Logger.info('[EnhancedGorevTreeProvider] Tree view refreshed successfully');
        } catch (error) {
            Logger.error('[EnhancedGorevTreeProvider] Failed to refresh tree view:', error);
            throw error;
        }
    }

    /**
     * Filtreleri günceller
     */
    updateFilter(filter: Partial<TaskFilter>): void {
        this.config.filters = { ...this.config.filters, ...filter };
        this.events.onFilterChanged?.(this.config.filters);
        this._onDidChangeTreeData.fire();
    }

    /**
     * Tüm filtreleri temizler
     */
    clearFilters(): void {
        this.config.filters = {};
        this.events.onFilterChanged?.(this.config.filters);
        this._onDidChangeTreeData.fire();
    }

    /**
     * Gruplama stratejisini değiştirir
     */
    setGrouping(grouping: GroupingStrategy): void {
        this.config.grouping = grouping;
        this.config.expandedGroups = GroupingStrategyProvider.getDefaultExpandedGroups(grouping);
        this.dragDropController.setGroupingStrategy(grouping);
        this.events.onGroupingChanged?.(grouping);
        this._onDidChangeTreeData.fire();
    }

    /**
     * Sıralama kriterini değiştirir
     */
    setSorting(criteria: SortingCriteria, ascending?: boolean): void {
        this.config.sorting = criteria;
        if (ascending !== undefined) {
            this.config.sortAscending = ascending;
        }
        this.events.onSortingChanged?.(criteria, this.config.sortAscending);
        this._onDidChangeTreeData.fire();
    }

    /**
     * Grup genişletme/daraltma durumunu günceller
     */
    toggleGroupExpansion(groupKey: string): void {
        if (this.config.expandedGroups.has(groupKey)) {
            this.config.expandedGroups.delete(groupKey);
            this.events.onGroupCollapsed?.(groupKey);
        } else {
            this.config.expandedGroups.add(groupKey);
            this.events.onGroupExpanded?.(groupKey);
        }
    }

    /**
     * Görev seçimini günceller
     */
    selectTask(taskId: string, multiSelect: boolean = false, rangeSelect: boolean = false): void {
        if (!multiSelect && !rangeSelect) {
            // Tek seçim
            this.selection.selectedTasks.clear();
            this.selection.selectedTasks.add(taskId);
            this.selection.lastSelectedTask = taskId;
            this.selection.anchorTask = taskId;
        } else if (multiSelect) {
            // Ctrl/Cmd + Click
            if (this.selection.selectedTasks.has(taskId)) {
                this.selection.selectedTasks.delete(taskId);
            } else {
                this.selection.selectedTasks.add(taskId);
            }
            this.selection.lastSelectedTask = taskId;
        } else if (rangeSelect && this.selection.anchorTask) {
            // Shift + Click
            this.selectRange(this.selection.anchorTask, taskId);
        }

        // Event'leri tetikle
        const selectedTasks = this.getSelectedTasks();
        if (selectedTasks.length === 1) {
            this.events.onTaskSelected?.(selectedTasks[0]);
        } else {
            this.events.onTasksSelected?.(selectedTasks);
        }

        this._onDidChangeTreeData.fire();
    }

    /**
     * Aralık seçimi yapar
     */
    private selectRange(fromId: string, toId: string): void {
        const fromIndex = this.filteredTasks.findIndex(t => t.id === fromId);
        const toIndex = this.filteredTasks.findIndex(t => t.id === toId);
        
        if (fromIndex === -1 || toIndex === -1) return;

        const start = Math.min(fromIndex, toIndex);
        const end = Math.max(fromIndex, toIndex);

        this.selection.selectedTasks.clear();
        for (let i = start; i <= end; i++) {
            this.selection.selectedTasks.add(this.filteredTasks[i].id!);
        }
        
        this.selection.lastSelectedTask = toId;
    }

    /**
     * Seçili görevleri döndürür
     */
    getSelectedTasks(): Gorev[] {
        return this.tasks.filter(task => 
            task.id && this.selection.selectedTasks.has(task.id)
        );
    }

    /**
     * Event handler'ları ayarlar
     */
    setEventHandlers(events: TreeViewEvents): void {
        this.events = events;
    }


    /**
     * Drag & Drop: Drag başladığında
     */
    async handleDrag(
        source: readonly EnhancedTreeViewItem[],
        dataTransfer: vscode.DataTransfer,
        token: vscode.CancellationToken
    ): Promise<void> {
        return this.dragDropController.handleDrag(source, dataTransfer, token);
    }

    /**
     * Drag & Drop: Drop yapıldığında
     */
    async handleDrop(
        target: EnhancedTreeViewItem | undefined,
        dataTransfer: vscode.DataTransfer,
        token: vscode.CancellationToken
    ): Promise<void> {
        await this.dragDropController.handleDrop(target, dataTransfer, token);
        // Drop sonrası TreeView'ı yenile
        await this.refresh();
    }
}

/**
 * Grup TreeView item'ı
 */
export class GroupTreeViewItem extends vscode.TreeItem {
    type = TreeItemType.Group;

    constructor(
        public groupKey: string,
        public groupType: GroupingStrategy,
        public tasks: Gorev[],
        expanded: boolean = true,
        private projects?: Map<string, any>
    ) {
        let label = GroupingStrategyProvider.getGroupLabel(groupKey, groupType);
        
        // Proje gruplandırması için proje ismini kullan
        if (groupType === GroupingStrategy.ByProject && projects && groupKey !== 'no-project') {
            const project = projects.get(groupKey);
            if (project && project.isim) {
                label = project.isim;
            }
        }
        
        const collapsibleState = expanded 
            ? vscode.TreeItemCollapsibleState.Expanded 
            : vscode.TreeItemCollapsibleState.Collapsed;

        super(label, collapsibleState);

        // Grup metadata'sı
        const completedCount = tasks.filter(t => t.durum === GorevDurum.Tamamlandi).length;
        const overdueCount = tasks.filter(t => 
            t.son_tarih && new Date(t.son_tarih) < new Date() && t.durum !== GorevDurum.Tamamlandi
        ).length;

        // Icon ve açıklama
        this.iconPath = GroupingStrategyProvider.getGroupIcon(groupKey, groupType);
        this.description = GroupingStrategyProvider.getGroupDescription(groupKey, tasks.length, {
            completedCount,
            overdueCount
        });

        // Badge ekle
        if (tasks.length > 0) {
            const badge = GroupingStrategyProvider.createGroupBadge(tasks.length, completedCount);
            this.description = `${badge} ${this.description || ''}`.trim();
        }

        // Context value
        this.contextValue = GroupingStrategyProvider.getGroupContextValue(groupKey, groupType);

        // Tooltip
        this.tooltip = this.createTooltip(completedCount, overdueCount);
    }

    private createTooltip(completedCount: number, overdueCount: number): string {
        const lines = [
            `${this.label}`,
            `Toplam: ${this.tasks.length} görev`,
        ];

        if (completedCount > 0) {
            lines.push(`Tamamlanan: ${completedCount}`);
        }

        if (overdueCount > 0) {
            lines.push(`Gecikmiş: ${overdueCount}`);
        }

        const highPriorityCount = this.tasks.filter(t => t.oncelik === GorevOncelik.Yuksek).length;
        if (highPriorityCount > 0) {
            lines.push(`Yüksek öncelik: ${highPriorityCount}`);
        }

        return lines.join('\n');
    }
}

/**
 * Görev TreeView item'ı
 */
export class TaskTreeViewItem extends vscode.TreeItem {
    type = TreeItemType.Task;

    constructor(
        public task: Gorev,
        private selection: TaskSelection,
        public parent?: GroupTreeViewItem
    ) {
        // Alt görevleri varsa collapsible yap
        const collapsibleState = (task.alt_gorevler && task.alt_gorevler.length > 0) 
            ? vscode.TreeItemCollapsibleState.Collapsed 
            : vscode.TreeItemCollapsibleState.None;
        
        super(task.baslik, collapsibleState);

        // Seçim durumu
        const isSelected = !!(task.id && selection.selectedTasks.has(task.id));

        // Icon
        this.iconPath = this.getTaskIcon(isSelected) as any;

        // Açıklama - configuration'a göre ayarla
        const config = vscode.workspace.getConfiguration('gorev.treeView.visuals');
        if (config.get('showPriorityBadges', true) || 
            config.get('showDueDateIndicators', true) || 
            config.get('showDependencyBadges', true) || 
            config.get('showProgressBars', true) ||
            config.get('showTagPills', true)) {
            this.description = this.getTaskDescription();
        }

        // Context value - seçim durumuna göre
        if (isSelected && selection.selectedTasks.size > 1) {
            // Çoklu seçim varsa
            this.contextValue = 'task:selected';
        } else {
            // Tek görev veya seçili değil
            this.contextValue = 'task';
        }

        // Tooltip
        this.tooltip = this.getTaskTooltip();

        // Command (tıklama işlemi) - task detayını açar
        this.command = {
            command: 'gorev.showTaskDetail',
            title: 'Show Task Detail',
            arguments: [this]
        };
    }

    private getTaskIcon(isSelected: boolean): vscode.ThemeIcon {
        let iconName: string;
        let color: string | undefined;

        // Durum bazlı icon
        if (this.task.durum === GorevDurum.Tamamlandi) {
            iconName = isSelected ? 'pass-filled' : 'pass';
            color = 'testing.iconPassed';
        } else if (this.task.durum === GorevDurum.DevamEdiyor) {
            iconName = isSelected ? 'debug-pause' : 'debug-start';
            color = 'debugIcon.pauseForeground';
        } else {
            iconName = isSelected ? 'circle-filled' : 'circle-outline';
        }

        // Öncelik rengi (tamamlanmamış görevler için)
        if (this.task.durum !== GorevDurum.Tamamlandi) {
            switch (this.task.oncelik) {
                case GorevOncelik.Yuksek:
                    color = COLORS.HIGH_PRIORITY;
                    break;
                case GorevOncelik.Orta:
                    color = COLORS.MEDIUM_PRIORITY;
                    break;
                case GorevOncelik.Dusuk:
                    color = COLORS.LOW_PRIORITY;
                    break;
            }
        }

        return color 
            ? new vscode.ThemeIcon(iconName, new vscode.ThemeColor(color))
            : new vscode.ThemeIcon(iconName);
    }

    private getTaskDescription(): string {
        const parts = [];

        // Progress indicator for parent tasks
        if (this.task.alt_gorevler && this.task.alt_gorevler.length > 0) {
            const completedCount = this.task.alt_gorevler.filter(t => t.durum === GorevDurum.Tamamlandi).length;
            const total = this.task.alt_gorevler.length;
            const percentage = Math.round((completedCount / total) * 100);
            
            // Visual progress bar
            const filledBlocks = Math.round(percentage / 10);
            const emptyBlocks = 10 - filledBlocks;
            const progressBar = '█'.repeat(filledBlocks) + '░'.repeat(emptyBlocks);
            
            parts.push(`[${progressBar}] ${percentage}%`);
            Logger.debug(`[TaskTreeViewItem] Task "${this.task.baslik}" has ${total} subtasks, ${completedCount} completed (${percentage}%)`);
        }

        // Priority indicator with colored badges (for non-completed tasks)
        if (this.task.durum !== GorevDurum.Tamamlandi) {
            switch (this.task.oncelik) {
                case GorevOncelik.Yuksek:
                    parts.push('🔥 Yüksek');
                    break;
                case GorevOncelik.Orta:
                    parts.push('⚡ Orta');
                    break;
                case GorevOncelik.Dusuk:
                    parts.push('ℹ️ Düşük');
                    break;
            }
        }

        // Due date with smart formatting
        if (this.task.son_tarih) {
            const dueDate = new Date(this.task.son_tarih);
            const today = new Date();
            today.setHours(0, 0, 0, 0);
            const tomorrow = new Date(today);
            tomorrow.setDate(tomorrow.getDate() + 1);
            const diffTime = dueDate.getTime() - today.getTime();
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
            
            if (this.task.durum !== GorevDurum.Tamamlandi) {
                if (dueDate < today) {
                    parts.push(`📅 ${Math.abs(diffDays)}g gecikmiş!`);
                } else if (dueDate.toDateString() === today.toDateString()) {
                    parts.push('📅 Bugün!');
                } else if (dueDate.toDateString() === tomorrow.toDateString()) {
                    parts.push('📅 Yarın');
                } else if (diffDays <= 7) {
                    parts.push(`📅 ${diffDays}g kaldı`);
                } else {
                    // Format as date for distant dates
                    const dateStr = dueDate.toLocaleDateString('tr-TR', { day: 'numeric', month: 'short' });
                    parts.push(`📅 ${dateStr}`);
                }
            }
        }

        // Enhanced dependency indicators
        const depParts = [];
        
        // Dependencies with visual lock/unlock icons
        if (this.task.bagimli_gorev_sayisi && this.task.bagimli_gorev_sayisi > 0) {
            if (this.task.tamamlanmamis_bagimlilik_sayisi && this.task.tamamlanmamis_bagimlilik_sayisi > 0) {
                // Blocked by incomplete dependencies
                depParts.push(`🔒 ${this.task.tamamlanmamis_bagimlilik_sayisi}/${this.task.bagimli_gorev_sayisi}`);
            } else {
                // All dependencies completed
                depParts.push(`🔓 ${this.task.bagimli_gorev_sayisi}`);
            }
        }
        
        // Tasks that depend on this
        if (this.task.bu_goreve_bagimli_sayisi && this.task.bu_goreve_bagimli_sayisi > 0) {
            depParts.push(`🔗 ${this.task.bu_goreve_bagimli_sayisi}`);
        }
        
        if (depParts.length > 0) {
            parts.push(depParts.join(' '));
        }

        // Tags as colored pills (limit to 3 for space)
        if (this.task.etiketler && this.task.etiketler.length > 0) {
            const tagPills = this.task.etiketler.slice(0, 3).map(tag => `⬤ ${tag}`);
            if (this.task.etiketler.length > 3) {
                tagPills.push(`+${this.task.etiketler.length - 3}`);
            }
            parts.push(tagPills.join(' '));
        }

        return parts.join(' │ ');
    }

    private getTaskTooltip(): string | vscode.MarkdownString {
        const md = new vscode.MarkdownString();
        md.supportHtml = true;
        
        // Title with status icon
        let statusIcon = '';
        switch (this.task.durum) {
            case GorevDurum.Tamamlandi:
                statusIcon = '✅';
                break;
            case GorevDurum.DevamEdiyor:
                statusIcon = '🔄';
                break;
            case GorevDurum.Beklemede:
                statusIcon = '⏸️';
                break;
        }
        md.appendMarkdown(`## ${statusIcon} ${this.task.baslik}\n\n`);

        // Priority with visual indicator
        let priorityBadge = '';
        let priorityColor = '';
        switch (this.task.oncelik) {
            case GorevOncelik.Yuksek:
                priorityBadge = '🔥';
                priorityColor = 'red';
                break;
            case GorevOncelik.Orta:
                priorityBadge = '⚡';
                priorityColor = 'orange';
                break;
            case GorevOncelik.Dusuk:
                priorityBadge = 'ℹ️';
                priorityColor = 'blue';
                break;
        }
        md.appendMarkdown(`**Öncelik:** ${priorityBadge} <span style="color: ${priorityColor}">${this.task.oncelik}</span>\n\n`);

        // Due date with smart formatting
        if (this.task.son_tarih) {
            const dueDate = new Date(this.task.son_tarih);
            const today = new Date();
            today.setHours(0, 0, 0, 0);
            const diffTime = dueDate.getTime() - today.getTime();
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
            
            let dueDateText = dueDate.toLocaleDateString('tr-TR', { 
                weekday: 'long', 
                year: 'numeric', 
                month: 'long', 
                day: 'numeric' 
            });
            
            if (this.task.durum !== GorevDurum.Tamamlandi) {
                if (diffDays < 0) {
                    dueDateText = `⚠️ <span style="color: red">${Math.abs(diffDays)} gün gecikmiş!</span>`;
                } else if (diffDays === 0) {
                    dueDateText = `📅 <span style="color: orange">Bugün!</span>`;
                } else if (diffDays === 1) {
                    dueDateText = `📅 <span style="color: orange">Yarın</span>`;
                } else if (diffDays <= 7) {
                    dueDateText = `📅 ${diffDays} gün kaldı`;
                }
            }
            
            md.appendMarkdown(`**Son Tarih:** ${dueDateText}\n\n`);
        }

        // Progress visualization for parent tasks
        if (this.task.alt_gorevler && this.task.alt_gorevler.length > 0) {
            const total = this.task.alt_gorevler.length;
            const completed = this.task.alt_gorevler.filter(t => t.durum === GorevDurum.Tamamlandi).length;
            const percentage = Math.round((completed / total) * 100);
            
            md.appendMarkdown(`### 📊 Alt Görev İlerlemesi\n\n`);
            
            // Visual progress bar
            const filledBlocks = Math.round(percentage / 5);
            const emptyBlocks = 20 - filledBlocks;
            const progressBar = '█'.repeat(filledBlocks) + '░'.repeat(emptyBlocks);
            
            md.appendMarkdown(`\`${progressBar}\` **${percentage}%**\n\n`);
            md.appendMarkdown(`✅ Tamamlanan: ${completed}/${total}\n\n`);
        }

        // Description
        if (this.task.aciklama) {
            md.appendMarkdown(`### 📝 Açıklama\n\n${this.task.aciklama}\n\n`);
        }

        // Dependencies visualization
        if ((this.task.bagimli_gorev_sayisi && this.task.bagimli_gorev_sayisi > 0) || 
            (this.task.bu_goreve_bagimli_sayisi && this.task.bu_goreve_bagimli_sayisi > 0)) {
            
            md.appendMarkdown(`### 🔗 Bağımlılıklar\n\n`);
            
            if (this.task.bagimli_gorev_sayisi && this.task.bagimli_gorev_sayisi > 0) {
                const completed = this.task.bagimli_gorev_sayisi - (this.task.tamamlanmamis_bagimlilik_sayisi || 0);
                const incomplete = this.task.tamamlanmamis_bagimlilik_sayisi || 0;
                
                md.appendMarkdown(`**Bu görev için beklenenler:**\n`);
                if (incomplete > 0) {
                    md.appendMarkdown(`- 🔒 ${incomplete} tamamlanmamış\n`);
                }
                if (completed > 0) {
                    md.appendMarkdown(`- 🔓 ${completed} tamamlanmış\n`);
                }
                md.appendMarkdown('\n');
            }
            
            if (this.task.bu_goreve_bagimli_sayisi && this.task.bu_goreve_bagimli_sayisi > 0) {
                md.appendMarkdown(`**Bu görevi bekleyenler:** 🔗 ${this.task.bu_goreve_bagimli_sayisi} görev\n\n`);
            }
        }

        // Tags as colored badges
        if (this.task.etiketler && this.task.etiketler.length > 0) {
            md.appendMarkdown(`### 🏷️ Etiketler\n\n`);
            this.task.etiketler.forEach(tag => {
                md.appendMarkdown(`\`${tag}\` `);
            });
            md.appendMarkdown('\n\n');
        }

        // Creation date
        if (this.task.olusturma_tarih) {
            const createdDate = new Date(this.task.olusturma_tarih);
            md.appendMarkdown(`---\n\n*Oluşturulma: ${createdDate.toLocaleDateString('tr-TR')}*`);
        }

        return md;
    }
}

/**
 * Boş durum TreeView item'ı
 */
export class EmptyTreeViewItem extends vscode.TreeItem {
    type = TreeItemType.Empty;

    constructor(message: string) {
        super(message, vscode.TreeItemCollapsibleState.None);
        this.contextValue = 'empty';
        this.iconPath = new vscode.ThemeIcon('info');
    }
}

/**
 * Enhanced TreeView item türleri
 */
export type EnhancedTreeViewItem = GroupTreeViewItem | TaskTreeViewItem | EmptyTreeViewItem;