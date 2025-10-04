import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ApiClient, ApiError, Task, Project } from '../api/client';
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
import { RefreshManager, RefreshTarget, RefreshProvider, RefreshReason, RefreshPriority } from '../managers/refreshManager';
import { performanceMonitor, measureAsync } from '../utils/performance';
import { debounceConfig } from '../utils/debounce';

/**
 * Gelişmiş görev TreeView provider'ı
 * Now with differential updates and RefreshManager integration
 */
export class EnhancedGorevTreeProvider implements vscode.TreeDataProvider<EnhancedTreeViewItem>, vscode.TreeDragAndDropController<EnhancedTreeViewItem>, RefreshProvider {
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

    // Differential update support
    private previousTasksHash: string = '';
    private previousProjectsHash: string = '';
    private treeDataCache: Map<string, EnhancedTreeViewItem> = new Map();
    private expansionStateCache: Map<string, boolean> = new Map();

    // Configuration change debouncing
    private debouncedConfigChange: ReturnType<typeof debounceConfig>;


    constructor(private apiClient: ApiClient) {
        // Initialize API client
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
        this.dragDropController = new DragDropController(apiClient);
        this.dropMimeTypes = this.dragDropController.dropMimeTypes;
        this.dragMimeTypes = this.dragDropController.dragMimeTypes;
        
        // Decoration provider
        this.decorationProvider = new TaskDecorationProvider();

        // Initialize debounced configuration change handler
        this.debouncedConfigChange = debounceConfig(this.handleConfigurationChange.bind(this));

        // Konfigürasyon değişikliklerini dinle (tek listener, RefreshManager ile koordineli)
        this.loadConfiguration();
        vscode.workspace.onDidChangeConfiguration(e => {
            if (e.affectsConfiguration('gorev.treeView')) {
                this.debouncedConfigChange();
            }
        });

        // Register with RefreshManager
        const refreshManager = RefreshManager.getInstance();
        refreshManager.registerProvider([RefreshTarget.TASKS], this);
    }

    /**
     * Clear all internal caches
     */
    clearCache(): void {
        Logger.info('[EnhancedGorevTreeProvider] Clearing all caches');
        this.treeDataCache.clear();
        this.expansionStateCache.clear();
        this.previousTasksHash = '';
        this.previousProjectsHash = '';
        Logger.debug('[EnhancedGorevTreeProvider] All caches cleared');
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
        
        if (!this.apiClient.isConnected()) {
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
            // Apply proper hierarchy handling for flat view too
            const allTaskIds = new Set(sortedTasks.map(t => t.id));
            const rootTasks = sortedTasks.filter(task => !task.parent_id);
            const orphanedSubtasks = sortedTasks.filter(task =>
                task.parent_id && !allTaskIds.has(task.parent_id)
            );

            Logger.debug(`[EnhancedGorevTreeProvider] Flat view: ${sortedTasks.length} total, ${rootTasks.length} root, ${orphanedSubtasks.length} orphaned`);

            // Combine root tasks and orphaned subtasks to prevent task loss
            const visibleTasks = [...rootTasks, ...orphanedSubtasks];

            if (sortedTasks.length > 0 && visibleTasks.length === 0) {
                Logger.warn(`[EnhancedGorevTreeProvider] No visible tasks in flat view after hierarchy filtering!`);
                Logger.warn(`[EnhancedGorevTreeProvider] Sample tasks:`, sortedTasks.slice(0, 3).map(t => ({
                    id: t.id,
                    baslik: t.baslik,
                    parent_id: t.parent_id
                })));
                // Show all tasks to prevent complete loss of visibility
                return sortedTasks.map(task => new TaskTreeViewItem(task, this.selection));
            }

            return visibleTasks.map(task => new TaskTreeViewItem(task, this.selection));
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
     * Görevleri yükler - Now using REST API
     */
    private async loadTasks(): Promise<void> {
        try {
            // Load projects first using API
            const projectsResponse = await this.apiClient.getProjects();
            if (projectsResponse.success && projectsResponse.data) {
                this.projects.clear();
                for (const project of projectsResponse.data) {
                    this.projects.set(project.id, project);
                }
                Logger.debug(`[EnhancedGorevTreeProvider] Loaded ${this.projects.size} projects via API`);
            }

            // Get page size from configuration
            const pageSize = vscode.workspace.getConfiguration('gorev').get<number>('pagination.pageSize', 100);

            // Check if we should show all projects
            const showAllProjects = this.config.filters?.showAllProjects !== false;

            // Get active project if needed
            let activeProjectId = '';
            if (!showAllProjects) {
                try {
                    const activeResponse = await this.apiClient.getActiveProject();
                    if (activeResponse.success && activeResponse.data) {
                        activeProjectId = activeResponse.data.id;
                        Logger.debug(`[EnhancedGorevTreeProvider] Active project: ${activeProjectId}`);
                    }
                } catch (err) {
                    if (err instanceof ApiError && !err.isNotFound()) {
                        Logger.warn('[EnhancedGorevTreeProvider] Failed to get active project:', err);
                    }
                }
            }

            // Fetch all tasks with pagination using API
            this.tasks = [];
            let offset = 0;
            let hasMoreTasks = true;
            const maxIterations = 50;
            let iterationCount = 0;

            while (hasMoreTasks && iterationCount < maxIterations) {
                iterationCount++;

                const tasksResponse = await this.apiClient.getTasks({
                    tum_projeler: showAllProjects,
                    limit: pageSize,
                    offset: offset
                });

                if (!tasksResponse.success || !tasksResponse.data) {
                    Logger.warn('[EnhancedGorevTreeProvider] No data in API response');
                    break;
                }

                const apiTasks = tasksResponse.data;
                Logger.info(`[EnhancedGorevTreeProvider] Fetched ${apiTasks.length} tasks (offset: ${offset}, total: ${tasksResponse.total || 'unknown'})`);

                if (apiTasks.length === 0) {
                    hasMoreTasks = false;
                    break;
                }

                // Convert API Task to Gorev model
                const gorevTasks = apiTasks.map(task => this.convertTaskToGorev(task));
                this.tasks.push(...gorevTasks);

                // Check pagination
                const total = tasksResponse.total || apiTasks.length;
                offset += apiTasks.length;

                if (offset >= total || apiTasks.length < pageSize) {
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
            Logger.info(`  - Page size used: ${pageSize}`);
            Logger.info(`  - Show all projects: ${showAllProjects}`);
            Logger.info(`  - Active project ID: ${activeProjectId || 'N/A'}`);
            
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
            
            // IMPORTANT: Set filtered tasks after loading with duplicate removal
            this.filteredTasks = this.removeDuplicateTasks([...this.tasks]);

            // DEBUG: Log detailed task information
            Logger.info(`[EnhancedGorevTreeProvider] FINAL TASK ANALYSIS:`);
            Logger.info(`[EnhancedGorevTreeProvider] - Raw tasks loaded: ${this.tasks.length}`);
            Logger.info(`[EnhancedGorevTreeProvider] - Filtered tasks after dedup: ${this.filteredTasks.length}`);
            Logger.info(`[EnhancedGorevTreeProvider] - Task titles loaded:`);
            this.tasks.forEach((task, idx) => {
                Logger.info(`[EnhancedGorevTreeProvider]   ${idx + 1}. "${task.baslik}" (ID: ${task.id}, parent: ${task.parent_id || 'none'})`);
            });

            // Check if we have the missing tasks
            const expectedMissingTasks = ['T001: Project Structure Setup', 'T002: Initialize Rust Workspace', 'T003: Initialize TypeScript Project', 'T009: Contract Test - search_code Tool'];
            expectedMissingTasks.forEach(title => {
                const found = this.tasks.find(t => t.baslik?.includes(title) || t.baslik?.includes(title.split(':')[0]));
                if (found) {
                    Logger.info(`[EnhancedGorevTreeProvider] ✅ Found expected task: "${found.baslik}" (ID: ${found.id})`);
                } else {
                    Logger.warn(`[EnhancedGorevTreeProvider] ❌ Missing expected task: "${title}"`);
                }
            });
            // Logger.debug(`[EnhancedGorevTreeProvider] Filtered ${this.tasks.length} tasks to ${this.filteredTasks.length} after duplicate removal`); // Reduced logging
            // Tasks filtered and ready
        } catch (error) {
            Logger.error('Failed to load tasks:', error);
            throw error;
        }
    }

    /**
     * TreeView'ı yeniler - RefreshProvider interface implementation
     */
    async refresh(): Promise<void> {
        await measureAsync(
            'enhanced-tree-refresh',
            async () => {
                // Logger.debug('[EnhancedGorevTreeProvider] Starting refresh...'); // Reduced logging

                try {
                    // Load fresh data
                    await this.loadTasks();

                    // Check if data actually changed using hashes
                    const currentTasksHash = this.calculateTasksHash();
                    const currentProjectsHash = this.calculateProjectsHash();

                    const tasksChanged = currentTasksHash !== this.previousTasksHash;
                    const projectsChanged = currentProjectsHash !== this.previousProjectsHash;

                    if (tasksChanged || projectsChanged) {
                        // Logger.debug('[EnhancedGorevTreeProvider] Data changed, updating tree'); // Reduced logging

                        // Update hashes
                        this.previousTasksHash = currentTasksHash;
                        this.previousProjectsHash = currentProjectsHash;

                        // Clear cache for changed data
                        if (tasksChanged) {
                            this.treeDataCache.clear();
                        }

                        // Fire selective change event
                        this._onDidChangeTreeData.fire(undefined);

                        // Logger.debug('[EnhancedGorevTreeProvider] Tree refreshed with new data'); // Reduced logging
                    } else {
                        // Logger.debug('[EnhancedGorevTreeProvider] No data changes detected, skipping tree update'); // Reduced logging
                    }

                } catch (error) {
                    Logger.error('[EnhancedGorevTreeProvider] Failed to refresh tree view:', error);
                    throw error;
                }
            },
            'tree-refresh',
            { provider: 'enhanced-tree' }
        );
    }

    /**
     * Force full refresh bypassing differential checks
     */
    async forceRefresh(): Promise<void> {
        Logger.info('[EnhancedGorevTreeProvider] Force refreshing tree view...');

        // Clear all caches
        this.tasks = [];
        this.filteredTasks = [];
        this.treeDataCache.clear();
        this.previousTasksHash = '';
        this.previousProjectsHash = '';

        await this.refresh();
    }

    /**
     * RefreshProvider interface implementation
     */
    getName(): string {
        return 'EnhancedGorevTreeProvider';
    }

    supportsTarget(target: RefreshTarget): boolean {
        return target === RefreshTarget.TASKS || target === RefreshTarget.ALL;
    }

    /**
     * Handle configuration changes with debouncing
     */
    private async handleConfigurationChange(): Promise<void> {
        Logger.debug('[EnhancedGorevTreeProvider] Configuration changed');
        this.loadConfiguration();

        // Request refresh through RefreshManager instead of direct call
        const refreshManager = RefreshManager.getInstance();
        await refreshManager.requestRefresh(
            RefreshReason.CONFIG_CHANGE,
            [RefreshTarget.TASKS],
            RefreshPriority.NORMAL
        );
    }

    /**
     * Calculate hash of current tasks for change detection
     */
    private calculateTasksHash(): string {
        const taskData = this.tasks.map(task => ({
            id: task.id,
            baslik: task.baslik,
            durum: task.durum,
            oncelik: task.oncelik,
            olusturma_tarihi: task.olusturma_tarihi,
            guncelleme_tarihi: task.guncelleme_tarihi,
            parent_id: task.parent_id
        }));

        return this.generateHash(JSON.stringify(taskData));
    }

    /**
     * Calculate hash of current projects for change detection
     */
    private calculateProjectsHash(): string {
        const projectData = Array.from(this.projects.entries());
        return this.generateHash(JSON.stringify(projectData));
    }

    /**
     * Simple hash function for change detection
     */
    private generateHash(input: string): string {
        let hash = 0;
        for (let i = 0; i < input.length; i++) {
            const char = input.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash = hash & hash; // Convert to 32bit integer
        }
        return hash.toString();
    }

    /**
     * Remove duplicate tasks based on ID
     * Keeps the most recently updated version when duplicates exist
     */
    private removeDuplicateTasks(tasks: Gorev[]): Gorev[] {
        const taskMap = new Map<string, Gorev>();

        for (const task of tasks) {
            if (!task.id) {
                // Skip tasks without valid IDs
                Logger.warn('[EnhancedGorevTreeProvider] Skipping task without ID:', task.baslik);
                continue;
            }

            const existingTask = taskMap.get(task.id);
            if (!existingTask) {
                // First occurrence of this task ID
                taskMap.set(task.id, task);
            } else {
                // Duplicate found - keep the one with the latest update date
                const currentDate = task.guncelleme_tarihi ? new Date(task.guncelleme_tarihi) : new Date(0);
                const existingDate = existingTask.guncelleme_tarihi ? new Date(existingTask.guncelleme_tarihi) : new Date(0);

                if (currentDate > existingDate) {
                    Logger.debug(`[EnhancedGorevTreeProvider] Replacing duplicate task ${task.id} with newer version`);
                    taskMap.set(task.id, task);
                } else {
                    Logger.debug(`[EnhancedGorevTreeProvider] Keeping existing version of duplicate task ${task.id}`);
                }
            }
        }

        return Array.from(taskMap.values());
    }

    /**
     * Save expansion state before refresh
     */
    private saveExpansionState(element?: EnhancedTreeViewItem): void {
        if (element && element.id) {
            // Note: Expansion state saving will be implemented when tree item has expanded property
            // this.expansionStateCache.set(element.id, element.expanded || false);
        }
    }

    /**
     * Restore expansion state after refresh
     */
    private restoreExpansionState(element: EnhancedTreeViewItem): void {
        if (element.id && this.expansionStateCache.has(element.id)) {
            // Note: Expansion state restoration will be implemented when tree item has expanded property
            // element.expanded = this.expansionStateCache.get(element.id);
        }
    }


    /**
     * Get cache statistics for debugging
     */
    getCacheStats(): { treeDataSize: number; expansionStateSize: number } {
        return {
            treeDataSize: this.treeDataCache.size,
            expansionStateSize: this.expansionStateCache.size
        };
    }

    /**
     * Dispose and cleanup
     */
    dispose(): void {
        this.debouncedConfigChange.cancel();
        this.clearCache();

        // Unregister from RefreshManager
        const refreshManager = RefreshManager.getInstance();
        refreshManager.unregisterProvider(this);

        Logger.debug('[EnhancedGorevTreeProvider] Disposed');
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

    /**
     * Convert API Task to internal Gorev model
     */
    private convertTaskToGorev(task: Task): Gorev {
        return {
            id: task.id,
            baslik: task.baslik,
            aciklama: task.aciklama,
            durum: task.durum as GorevDurum,
            oncelik: task.oncelik as GorevOncelik,
            proje_id: task.proje_id || '',
            son_tarih: task.son_tarih,
            etiketler: task.etiketler,
            olusturma_tarihi: task.olusturma_tarihi,
            guncelleme_tarihi: task.guncelleme_tarihi,
        };
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
            const tagPills = this.task.etiketler.slice(0, 3).map(tag => `⬤ ${tag.isim}`);
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
                md.appendMarkdown(`\`${tag.isim}\` `);
            });
            md.appendMarkdown('\n\n');
        }

        // Creation date
        if (this.task.olusturma_tarihi) {
            const createdDate = new Date(this.task.olusturma_tarihi);
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
