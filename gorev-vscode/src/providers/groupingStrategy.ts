import * as vscode from 'vscode';
import { GroupingStrategy, GroupTreeItem, TreeItemType } from '../models/treeModels';
import { GorevDurum, GorevOncelik } from '../models/common';
import { ICONS, COLORS } from '../utils/constants';

/**
 * Gruplama stratejileri için label ve icon sağlayıcı
 */
export class GroupingStrategyProvider {
    /**
     * Grup için label döndürür
     */
    static getGroupLabel(groupKey: string, strategy: GroupingStrategy): string {
        switch (strategy) {
            case GroupingStrategy.ByStatus:
                return this.getStatusLabel(groupKey);
            case GroupingStrategy.ByPriority:
                return this.getPriorityLabel(groupKey);
            case GroupingStrategy.ByDueDate:
                return this.getDueDateLabel(groupKey);
            case GroupingStrategy.ByTag:
                return groupKey === 'no-tag' ? 'Etiketsiz' : `#${groupKey}`;
            case GroupingStrategy.ByProject:
                return groupKey === 'no-project' ? 'Projesiz' : groupKey;
            default:
                return 'Tüm Görevler';
        }
    }

    /**
     * Grup için icon döndürür
     */
    static getGroupIcon(groupKey: string, strategy: GroupingStrategy): vscode.ThemeIcon {
        switch (strategy) {
            case GroupingStrategy.ByStatus:
                return this.getStatusIcon(groupKey);
            case GroupingStrategy.ByPriority:
                return this.getPriorityIcon(groupKey);
            case GroupingStrategy.ByDueDate:
                return this.getDueDateIcon(groupKey);
            case GroupingStrategy.ByTag:
                return new vscode.ThemeIcon('tag');
            case GroupingStrategy.ByProject:
                return new vscode.ThemeIcon('folder');
            default:
                return new vscode.ThemeIcon('list-unordered');
        }
    }

    /**
     * Grup için açıklama döndürür
     */
    static getGroupDescription(groupKey: string, taskCount: number, metadata?: any): string {
        const countText = `${taskCount} görev`;
        
        switch (groupKey) {
            case 'overdue':
                return `${countText} • Gecikmiş!`;
            case 'today':
                return `${countText} • Bugün`;
            case 'this-week':
                return `${countText} • Bu hafta`;
            default:
                return countText;
        }
    }

    private static getStatusLabel(status: string): string {
        const labels: { [key: string]: string } = {
            [GorevDurum.Beklemede]: 'Beklemede',
            [GorevDurum.DevamEdiyor]: 'Devam Ediyor',
            [GorevDurum.Tamamlandi]: 'Tamamlandı'
        };
        return labels[status] || status;
    }

    private static getStatusIcon(status: string): vscode.ThemeIcon {
        switch (status) {
            case GorevDurum.Beklemede:
                return new vscode.ThemeIcon('circle-outline');
            case GorevDurum.DevamEdiyor:
                return new vscode.ThemeIcon('sync~spin');
            case GorevDurum.Tamamlandi:
                return new vscode.ThemeIcon('check');
            default:
                return new vscode.ThemeIcon('question');
        }
    }

    private static getPriorityLabel(priority: string): string {
        const labels: { [key: string]: string } = {
            [GorevOncelik.Yuksek]: 'Yüksek Öncelik',
            [GorevOncelik.Orta]: 'Orta Öncelik',
            [GorevOncelik.Dusuk]: 'Düşük Öncelik'
        };
        return labels[priority] || priority;
    }

    private static getPriorityIcon(priority: string): vscode.ThemeIcon {
        let iconName: string;
        let color: string | undefined;

        switch (priority) {
            case GorevOncelik.Yuksek:
                iconName = 'flame';
                color = COLORS.HIGH_PRIORITY;
                break;
            case GorevOncelik.Orta:
                iconName = 'warning';
                color = COLORS.MEDIUM_PRIORITY;
                break;
            case GorevOncelik.Dusuk:
                iconName = 'info';
                color = COLORS.LOW_PRIORITY;
                break;
            default:
                iconName = 'circle-outline';
        }

        return color 
            ? new vscode.ThemeIcon(iconName, new vscode.ThemeColor(color))
            : new vscode.ThemeIcon(iconName);
    }

    private static getDueDateLabel(groupKey: string): string {
        const labels: { [key: string]: string } = {
            'overdue': 'Gecikmiş',
            'today': 'Bugün',
            'this-week': 'Bu Hafta',
            'later': 'Daha Sonra',
            'no-due-date': 'Tarihsiz'
        };
        return labels[groupKey] || groupKey;
    }

    private static getDueDateIcon(groupKey: string): vscode.ThemeIcon {
        switch (groupKey) {
            case 'overdue':
                return new vscode.ThemeIcon('alert', new vscode.ThemeColor('errorForeground'));
            case 'today':
                return new vscode.ThemeIcon('calendar', new vscode.ThemeColor('warningForeground'));
            case 'this-week':
                return new vscode.ThemeIcon('calendar');
            case 'later':
                return new vscode.ThemeIcon('calendar-clock');
            case 'no-due-date':
                return new vscode.ThemeIcon('calendar-remove');
            default:
                return new vscode.ThemeIcon('calendar');
        }
    }

    /**
     * Grup sıralama karşılaştırıcısı
     */
    static compareGroups(a: string, b: string, strategy: GroupingStrategy): number {
        switch (strategy) {
            case GroupingStrategy.ByStatus:
                return this.compareStatusGroups(a, b);
            case GroupingStrategy.ByPriority:
                return this.comparePriorityGroups(a, b);
            case GroupingStrategy.ByDueDate:
                return this.compareDueDateGroups(a, b);
            default:
                return a.localeCompare(b);
        }
    }

    private static compareStatusGroups(a: string, b: string): number {
        const order: { [key: string]: number } = {
            [GorevDurum.DevamEdiyor]: 1,
            [GorevDurum.Beklemede]: 2,
            [GorevDurum.Tamamlandi]: 3
        };
        return (order[a] || 99) - (order[b] || 99);
    }

    private static comparePriorityGroups(a: string, b: string): number {
        const order: { [key: string]: number } = {
            [GorevOncelik.Yuksek]: 1,
            [GorevOncelik.Orta]: 2,
            [GorevOncelik.Dusuk]: 3
        };
        return (order[a] || 99) - (order[b] || 99);
    }

    private static compareDueDateGroups(a: string, b: string): number {
        const order: { [key: string]: number } = {
            'overdue': 1,
            'today': 2,
            'this-week': 3,
            'later': 4,
            'no-due-date': 5
        };
        return (order[a] || 99) - (order[b] || 99);
    }

    /**
     * Varsayılan olarak açık olması gereken grupları döndürür
     */
    static getDefaultExpandedGroups(strategy: GroupingStrategy): Set<string> {
        const expanded = new Set<string>();

        switch (strategy) {
            case GroupingStrategy.ByStatus:
                expanded.add(GorevDurum.DevamEdiyor);
                expanded.add(GorevDurum.Beklemede);
                break;
            case GroupingStrategy.ByPriority:
                expanded.add(GorevOncelik.Yuksek);
                expanded.add(GorevOncelik.Orta);
                break;
            case GroupingStrategy.ByDueDate:
                expanded.add('overdue');
                expanded.add('today');
                expanded.add('this-week');
                break;
        }

        return expanded;
    }

    /**
     * Grup için context value döndürür (sağ tık menüsü için)
     */
    static getGroupContextValue(groupKey: string, strategy: GroupingStrategy): string {
        return `group:${strategy}:${groupKey}`;
    }

    /**
     * Boş grup için mesaj döndürür
     */
    static getEmptyGroupMessage(groupKey: string, strategy: GroupingStrategy): string {
        switch (strategy) {
            case GroupingStrategy.ByStatus:
                return `${this.getStatusLabel(groupKey)} görev yok`;
            case GroupingStrategy.ByPriority:
                return `${this.getPriorityLabel(groupKey)} görev yok`;
            case GroupingStrategy.ByDueDate:
                return `${this.getDueDateLabel(groupKey)} için görev yok`;
            default:
                return 'Bu grupta görev yok';
        }
    }

    /**
     * Grup için badge (görev sayısı göstergesi) oluşturur
     */
    static createGroupBadge(taskCount: number, completedCount?: number): string {
        if (completedCount !== undefined && completedCount > 0) {
            return `${completedCount}/${taskCount}`;
        }
        return taskCount.toString();
    }
}
