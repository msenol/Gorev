import * as vscode from 'vscode';
import { Gorev, GorevDurum, GorevOncelik } from './gorev';

/**
 * Gruplama stratejileri
 */
export enum GroupingStrategy {
    None = 'none',
    ByStatus = 'status',
    ByPriority = 'priority',
    ByProject = 'project',
    ByTag = 'tag',
    ByDueDate = 'dueDate'
}

/**
 * Sıralama kriterleri
 */
export enum SortingCriteria {
    Title = 'title',
    Priority = 'priority',
    DueDate = 'dueDate',
    CreatedDate = 'createdDate',
    Status = 'status'
}

/**
 * TreeView item türleri
 */
export enum TreeItemType {
    Group = 'group',
    Task = 'task',
    LoadMore = 'loadMore',
    Empty = 'empty'
}

/**
 * Grup item'ı için model
 */
export interface GroupTreeItem {
    type: TreeItemType.Group;
    label: string;
    groupKey: string;
    groupType: GroupingStrategy;
    taskCount: number;
    collapsibleState: vscode.TreeItemCollapsibleState;
    children: Gorev[];
    iconPath?: vscode.ThemeIcon;
    description?: string;
}

/**
 * Görev item'ı için model
 */
export interface TaskTreeItem {
    type: TreeItemType.Task;
    task: Gorev;
    parent?: GroupTreeItem;
}

/**
 * Boş durum item'ı
 */
export interface EmptyTreeItem {
    type: TreeItemType.Empty;
    message: string;
}

/**
 * Load more item'ı (pagination için)
 */
export interface LoadMoreTreeItem {
    type: TreeItemType.LoadMore;
    parent?: GroupTreeItem;
    offset: number;
    limit: number;
}

/**
 * TreeView item türlerinin union'ı
 */
export type EnhancedTreeItem = GroupTreeItem | TaskTreeItem | EmptyTreeItem | LoadMoreTreeItem;

/**
 * Filtre modeli
 */
export interface TaskFilter {
    searchQuery?: string;
    durum?: GorevDurum;
    oncelik?: GorevOncelik;
    projeId?: string;
    tags?: string[];
    dueDateRange?: {
        start?: Date;
        end?: Date;
    };
    overdue?: boolean;
    dueToday?: boolean;
    dueThisWeek?: boolean;
    hasTag?: boolean;
    hasDependency?: boolean;
    showAllProjects?: boolean;
}

/**
 * TreeView konfigürasyonu
 */
export interface TreeViewConfig {
    grouping: GroupingStrategy;
    sorting: SortingCriteria;
    sortAscending: boolean;
    showCompleted: boolean;
    showEmptyGroups: boolean;
    expandedGroups: Set<string>;
    filters: TaskFilter;
}

/**
 * Çoklu seçim için selection model
 */
export interface TaskSelection {
    selectedTasks: Set<string>; // task id'leri
    lastSelectedTask?: string;
    anchorTask?: string; // shift+click için
}

/**
 * TreeView event'leri
 */
export interface TreeViewEvents {
    onTaskSelected?: (task: Gorev) => void;
    onTasksSelected?: (tasks: Gorev[]) => void;
    onGroupCollapsed?: (groupKey: string) => void;
    onGroupExpanded?: (groupKey: string) => void;
    onFilterChanged?: (filter: TaskFilter) => void;
    onGroupingChanged?: (grouping: GroupingStrategy) => void;
    onSortingChanged?: (sorting: SortingCriteria, ascending: boolean) => void;
}

/**
 * Grup metadata'sı
 */
export interface GroupMetadata {
    totalTasks: number;
    completedTasks: number;
    highPriorityTasks: number;
    overdueTasks: number;
    todayTasks: number;
    averageCompletionTime?: number;
}

/**
 * TreeView için utility fonksiyonlar
 */
export class TreeViewUtils {
    /**
     * Görevleri belirtilen stratejiye göre gruplar
     */
    static groupTasks(tasks: Gorev[], strategy: GroupingStrategy): Map<string, Gorev[]> {
        const groups = new Map<string, Gorev[]>();

        if (strategy === GroupingStrategy.None) {
            groups.set('all', tasks);
            return groups;
        }

        tasks.forEach(task => {
            const groupKey = this.getGroupKey(task, strategy);
            if (!groups.has(groupKey)) {
                groups.set(groupKey, []);
            }
            const group = groups.get(groupKey);
            if (group) {
                group.push(task);
            }
        });

        return groups;
    }

    /**
     * Görev için grup anahtarını döndürür
     */
    private static getGroupKey(task: Gorev, strategy: GroupingStrategy): string {
        switch (strategy) {
            case GroupingStrategy.ByStatus:
                return task.durum;
            case GroupingStrategy.ByPriority:
                return task.oncelik;
            case GroupingStrategy.ByProject:
                return task.proje_id || 'no-project';
            case GroupingStrategy.ByTag:
                return task.etiketler?.length ? task.etiketler[0].isim : 'no-tag';
            case GroupingStrategy.ByDueDate:
                return this.getDueDateGroup(task.son_tarih);
            default:
                return 'all';
        }
    }

    /**
     * Son tarih grubunu belirler
     */
    private static getDueDateGroup(dueDate?: string): string {
        if (!dueDate) return 'no-due-date';

        const due = new Date(dueDate);
        const now = new Date();
        const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const tomorrow = new Date(today);
        tomorrow.setDate(tomorrow.getDate() + 1);
        const nextWeek = new Date(today);
        nextWeek.setDate(nextWeek.getDate() + 7);

        if (due < today) return 'overdue';
        if (due >= today && due < tomorrow) return 'today';
        if (due >= tomorrow && due < nextWeek) return 'this-week';
        return 'later';
    }

    /**
     * Görevleri sıralar
     */
    static sortTasks(tasks: Gorev[], criteria: SortingCriteria, ascending: boolean): Gorev[] {
        const sorted = [...tasks].sort((a, b) => {
            let comparison = 0;

            switch (criteria) {
                case SortingCriteria.Title:
                    comparison = a.baslik.localeCompare(b.baslik);
                    break;
                case SortingCriteria.Priority:
                    comparison = this.comparePriority(a.oncelik as GorevOncelik, b.oncelik as GorevOncelik);
                    break;
                case SortingCriteria.DueDate:
                    comparison = this.compareDates(a.son_tarih, b.son_tarih);
                    break;
                case SortingCriteria.Status:
                    comparison = this.compareStatus(a.durum as GorevDurum, b.durum as GorevDurum);
                    break;
                case SortingCriteria.CreatedDate:
                    comparison = this.compareDates(a.olusturma_tarihi, b.olusturma_tarihi);
                    break;
            }

            return ascending ? comparison : -comparison;
        });

        return sorted;
    }

    private static comparePriority(a: GorevOncelik, b: GorevOncelik): number {
        const priorityOrder = {
            [GorevOncelik.Yuksek]: 3,
            [GorevOncelik.Orta]: 2,
            [GorevOncelik.Dusuk]: 1
        };
        return priorityOrder[a] - priorityOrder[b];
    }

    private static compareStatus(a: GorevDurum, b: GorevDurum): number {
        const statusOrder = {
            [GorevDurum.DevamEdiyor]: 3,
            [GorevDurum.Beklemede]: 2,
            [GorevDurum.Tamamlandi]: 1
        };
        return statusOrder[a] - statusOrder[b];
    }

    private static compareDates(a?: string, b?: string): number {
        if (!a && !b) return 0;
        if (!a) return 1;
        if (!b) return -1;
        return new Date(a).getTime() - new Date(b).getTime();
    }

    /**
     * Görevleri filtreler
     */
    static filterTasks(tasks: Gorev[], filter: TaskFilter): Gorev[] {
        return tasks.filter(task => {
            // Arama sorgusu
            if (filter.searchQuery) {
                const query = filter.searchQuery.toLowerCase();
                const matchesSearch =
                    task.baslik.toLowerCase().includes(query) ||
                    (task.aciklama && task.aciklama.toLowerCase().includes(query)) ||
                    (task.etiketler && task.etiketler.some(tag => tag.isim.toLowerCase().includes(query)));
                
                if (!matchesSearch) return false;
            }

            // Durum filtresi
            if (filter.durum) {
                if (task.durum !== filter.durum) return false;
            }

            // Öncelik filtresi
            if (filter.oncelik) {
                if (task.oncelik !== filter.oncelik) return false;
            }

            // Proje filtresi
            if (filter.projeId) {
                if (task.proje_id !== filter.projeId) return false;
            }

            // Etiket filtresi
            if (filter.tags && filter.tags.length > 0) {
                const filterTags = filter.tags;
                if (!task.etiketler || !task.etiketler.some(tag => filterTags.includes(tag.isim))) return false;
            }

            // Son tarih aralığı
            if (filter.dueDateRange) {
                if (!task.son_tarih) return false;
                const dueDate = new Date(task.son_tarih);
                if (filter.dueDateRange.start && dueDate < filter.dueDateRange.start) return false;
                if (filter.dueDateRange.end && dueDate > filter.dueDateRange.end) return false;
            }

            // Gecikmiş görevler
            if (filter.overdue) {
                const isOverdue = task.son_tarih && new Date(task.son_tarih) < new Date() && task.durum !== GorevDurum.Tamamlandi;
                if (!isOverdue) return false;
            }

            // Bugün biten görevler
            if (filter.dueToday) {
                if (!task.son_tarih) return false;
                const today = new Date();
                today.setHours(0, 0, 0, 0);
                const tomorrow = new Date(today);
                tomorrow.setDate(tomorrow.getDate() + 1);
                const dueDate = new Date(task.son_tarih);
                if (dueDate < today || dueDate >= tomorrow) return false;
            }

            // Bu hafta biten görevler
            if (filter.dueThisWeek) {
                if (!task.son_tarih) return false;
                const today = new Date();
                today.setHours(0, 0, 0, 0);
                const nextWeek = new Date(today);
                nextWeek.setDate(nextWeek.getDate() + 7);
                const dueDate = new Date(task.son_tarih);
                if (dueDate < today || dueDate >= nextWeek) return false;
            }

            // Etiketli görevler
            if (filter.hasTag) {
                if (!task.etiketler || task.etiketler.length === 0) return false;
            }

            // Bağımlılığı olan görevler
            if (filter.hasDependency) {
                if (!task.bagimliliklar || task.bagimliliklar.length === 0) return false;
            }

            return true;
        });
    }
}