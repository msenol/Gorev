import * as vscode from 'vscode';
import { t } from '../utils/l10n';
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
        
        // Load showAllProjects from configuration
        if (!this.config.filters) {
            this.config.filters = {};
        }
        this.config.filters.showAllProjects = config.get('showAllProjects', true);
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
        // Root level getChildren call
        
        if (!this.mcpClient.isConnected()) {
            Logger.warn('[EnhancedGorevTreeProvider] MCP client not connected');
            return [new EmptyTreeViewItem(t('enhancedTree.notConnected'))];
        }

        // Root level
        if (!element) {
            try {
                await this.loadTasks();
                const items = this.createRootItems();
                return items;
            } catch (error) {
                Logger.error('Failed to load tasks:', error);
                return [new EmptyTreeViewItem(t('enhancedTree.loadFailed'))];
            }
        }

        // Grup altındaki görevler
        if (element instanceof GroupTreeViewItem) {
            return this.createTaskItems(element);
        }

        // Task altındaki alt görevler
        if (element instanceof TaskTreeViewItem && element.task.alt_gorevler && element.task.alt_gorevler.length > 0) {
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
        // Tamamlanmış görevleri gizle
        if (!this.config.showCompleted) {
            this.filteredTasks = this.filteredTasks.filter(
                task => task.durum !== GorevDurum.Tamamlandi
            );
        }

        if (this.filteredTasks.length === 0) {
            return [new EmptyTreeViewItem(this.getEmptyMessage())];
        }

        // Gruplama yoksa direkt görevleri göster
        if (this.config.grouping === GroupingStrategy.None) {
            const sortedTasks = TreeViewUtils.sortTasks(
                this.filteredTasks, 
                this.config.sorting, 
                this.config.sortAscending
            );
            // Show all tasks (temporarily disable root filtering due to subtask hierarchy issues)
            const rootTasks = sortedTasks; // TODO: Fix subtask hierarchy display
            return rootTasks.map(task => new TaskTreeViewItem(task, this.selection));
        }

        // Görevleri grupla
        const groups = TreeViewUtils.groupTasks(this.filteredTasks, this.config.grouping);
        const groupItems: EnhancedTreeViewItem[] = [];

        // Grupları sırala
        const sortedGroupKeys = Array.from(groups.keys()).sort((a, b) =>
            GroupingStrategyProvider.compareGroups(a, b, this.config.grouping)
        );

        // Grup item'larını oluştur
        for (const groupKey of sortedGroupKeys) {
            const tasksInGroup = groups.get(groupKey)!;

            // Boş grupları gizle
            if (!this.config.showEmptyGroups && tasksInGroup.length === 0) {
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

        // Improved hierarchy handling: Show root tasks and orphaned subtasks
        const allTaskIds = new Set(sortedTasks.map(t => t.id));
        const rootTasks = sortedTasks.filter(task => !task.parent_id);
        const orphanedSubtasks = sortedTasks.filter(task =>
            task.parent_id && !allTaskIds.has(task.parent_id)
        );

        Logger.debug(`[EnhancedGorevTreeProvider] Group ${group.groupKey}: ${sortedTasks.length} total, ${rootTasks.length} root, ${orphanedSubtasks.length} orphaned`);

        // Combine root tasks and orphaned subtasks to prevent task loss
        const visibleTasks = [...rootTasks, ...orphanedSubtasks];

        if (sortedTasks.length > 0 && visibleTasks.length === 0) {
            Logger.warn(`[EnhancedGorevTreeProvider] No visible tasks in group ${group.groupKey} after hierarchy filtering!`);
            Logger.warn(`[EnhancedGorevTreeProvider] Sample tasks:`, sortedTasks.slice(0, 3).map(t => ({
                id: t.id,
                baslik: t.baslik,
                parent_id: t.parent_id,
                seviye: t.seviye
            })));

            // Show all tasks to prevent complete loss of visibility
            return sortedTasks.map(task => {
                const item = new TaskTreeViewItem(task, this.selection, group);
                this.decorationProvider.updateTaskDecoration(task, item);
                return item;
            });
        }

        return visibleTasks.map(task => {
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
            return t('enhancedTree.searchNoResults', this.config.filters.searchQuery);
        }
        if (Object.keys(this.config.filters).length > 0) {
            return t('enhancedTree.filterNoResults');
        }
        return t('enhancedTree.noTasks');
    }

    /**
     * Görevleri yükler
     */
    private async loadTasks(): Promise<void> {
        try {
            // Önce projeleri yükle
            // Loading projects
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
            
            // Get page size from configuration (default should match package.json)
            const pageSize = vscode.workspace.getConfiguration('gorev').get<number>('pagination.pageSize', 100);
            
            // Check if we should show all projects or just active project
            // Default to true to show all projects
            const showAllProjects = this.config.filters?.showAllProjects !== false;
            
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
                            // Active project loaded
                        }
                    }
                } catch (err) {
                    Logger.warn('[EnhancedGorevTreeProvider] Failed to get active project:', err);
                }
            }
            
            // Initialize tasks array
            this.tasks = [];
            let offset = 0;
            let hasMoreTasks = true;
            let totalTaskCount = 0;
            
            // Fetch all tasks with pagination
            while (hasMoreTasks) {
                // Fetching tasks
                
                const result = await this.mcpClient.callTool('gorev_listele', {
                    tum_projeler: showAllProjects,
                    limit: pageSize,
                    offset: offset
                });
                
                if (!result || !result.content || !result.content[0]) {
                    Logger.warn('[EnhancedGorevTreeProvider] No content in response');
                    break;
                }
                
                const responseText = result.content[0].text;
                
                // Check for pagination info: "Görevler (1-100 / 147)"
                const paginationMatch = responseText.match(/Görevler \((\d+)-(\d+) \/ (\d+)\)/);
                if (paginationMatch) {
                    const [_, start, end, total] = paginationMatch;
                    totalTaskCount = parseInt(total);
                    Logger.info(`[EnhancedGorevTreeProvider] Pagination: ${start}-${end} / ${total}`);
                    
                    // Parse the markdown content to extract tasks
                    const pageTasks = MarkdownParser.parseGorevListesi(responseText);
                    Logger.info('[EnhancedGorevTreeProvider] Parsed tasks from page:', pageTasks.length);
                    
                    // Add to our task list
                    this.tasks.push(...pageTasks);
                    
                    // Update offset for next page based on server response, not pageSize
                    // This handles server-side response limits correctly
                    offset = parseInt(end);
                    
                    // Check if we need to fetch more
                    if (offset >= totalTaskCount) {
                        hasMoreTasks = false;
                    }
                } else {
                    // No pagination info, parse and assume this is the last page
                    const pageTasks = MarkdownParser.parseGorevListesi(responseText);
                    Logger.info('[EnhancedGorevTreeProvider] Parsed tasks from page:', pageTasks.length);
                    this.tasks.push(...pageTasks);
                    hasMoreTasks = false;
                }
                
                // Safety check to prevent infinite loop
                if (offset > 1000 || this.tasks.length > 1000) {
                    Logger.warn(`[EnhancedGorevTreeProvider] Safety limit reached at offset=${offset}, tasks=${this.tasks.length}`);
                    Logger.warn('[EnhancedGorevTreeProvider] Consider increasing safety limits if your project has >1000 tasks');
                    break;
                }
            }
            
            Logger.info(`[EnhancedGorevTreeProvider] Task loading summary:`);
            Logger.info(`  - Total tasks fetched: ${this.tasks.length}`);
            Logger.info(`  - Expected total: ${totalTaskCount}`);
            Logger.info(`  - Page size used: ${pageSize}`);
            Logger.info(`  - Show all projects: ${showAllProjects}`);
            Logger.info(`  - Active project ID: ${activeProjectId || 'N/A'}`);

            if (this.tasks.length !== totalTaskCount && totalTaskCount > 0) {
                Logger.warn(`[EnhancedGorevTreeProvider] TASK COUNT MISMATCH: Expected ${totalTaskCount}, got ${this.tasks.length}`);
            }
            
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
                // Tasks loaded successfully
                
                const tasksWithProjectId = this.tasks.filter(t => t.proje_id);
                const tasksWithoutProjectId = this.tasks.filter(t => !t.proje_id);
                
                // Task categorization complete
                
                // Warn about tasks without project_id
                if (tasksWithoutProjectId.length > 0) {
                    Logger.warn('[EnhancedGorevTreeProvider] Found tasks without project_id:', 
                        tasksWithoutProjectId.map(t => ({ id: t.id, baslik: t.baslik })));
                    Logger.warn('[EnhancedGorevTreeProvider] These tasks may not appear when filtering by project');
                }
            } else {
                Logger.warn('[EnhancedGorevTreeProvider] No tasks parsed from response');
            }
            
            // IMPORTANT: Set filtered tasks after loading
            this.filteredTasks = [...this.tasks];
            // Tasks filtered and ready
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
     * Mevcut filtreyi döndürür
     */
    getFilter(): TaskFilter {
        return this.config.filters || {};
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
     * Tüm seçimleri temizler
     */
    clearSelection(): void {
        this.selection.selectedTasks.clear();
        this.selection.lastSelectedTask = undefined;
        this.selection.anchorTask = undefined;
        this._onDidChangeTreeData.fire();
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
            t('enhancedTree.totalTasks', this.tasks.length.toString()),
        ];

        if (completedCount > 0) {
            lines.push(t('enhancedTree.completedTasks', completedCount.toString()));
        }

        if (overdueCount > 0) {
            lines.push(t('enhancedTree.overdueTasks', overdueCount.toString()));
        }

        const highPriorityCount = this.tasks.filter(t => t.oncelik === GorevOncelik.Yuksek).length;
        if (highPriorityCount > 0) {
            lines.push(t('enhancedTree.highPriorityTasks', highPriorityCount.toString()));
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
        md.appendMarkdown(`**${t('enhancedTree.priority')}** ${priorityBadge} <span style="color: ${priorityColor}">${this.task.oncelik}</span>\n\n`);

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
                    dueDateText = `⚠️ <span style="color: red">${t('enhancedTree.daysOverdue', Math.abs(diffDays).toString())}</span>`;
                } else if (diffDays === 0) {
                    dueDateText = `📅 <span style="color: orange">${t('enhancedTree.today')}</span>`;
                } else if (diffDays === 1) {
                    dueDateText = `📅 <span style="color: orange">${t('enhancedTree.tomorrow')}</span>`;
                } else if (diffDays <= 7) {
                    dueDateText = `📅 ${t('enhancedTree.daysLeft', diffDays.toString())}`;
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
