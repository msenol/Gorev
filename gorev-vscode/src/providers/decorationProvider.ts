import * as vscode from 'vscode';
import { Gorev, GorevDurum, GorevOncelik } from '../models/gorev';
import { TaskTreeViewItem } from './enhancedGorevTreeProvider';
import { COLORS } from '../utils/constants';
import { Logger } from '../utils/logger';

/**
 * Task decoration provider for visual enhancements
 */
export class TaskDecorationProvider implements vscode.FileDecorationProvider {
    private _onDidChangeFileDecorations = new vscode.EventEmitter<vscode.Uri | vscode.Uri[]>();
    readonly onDidChangeFileDecorations = this._onDidChangeFileDecorations.event;

    private decorations = new Map<string, vscode.FileDecoration>();

    provideFileDecoration(uri: vscode.Uri): vscode.FileDecoration | undefined {
        return this.decorations.get(uri.toString());
    }

    /**
     * Updates decoration for a task
     */
    updateTaskDecoration(task: Gorev, treeItem: TaskTreeViewItem): void {
        const decoration = this.createTaskDecoration(task);
        const uri = this.getTaskUri(task);
        
        this.decorations.set(uri.toString(), decoration);
        this._onDidChangeFileDecorations.fire(uri);
    }

    /**
     * Creates decoration for a task based on its properties
     */
    private createTaskDecoration(task: Gorev): vscode.FileDecoration {
        const badges: string[] = [];
        let color: string | undefined;
        let tooltip = '';

        // Status badge
        const statusBadge = this.getStatusBadge(task);
        if (statusBadge) {
            badges.push(statusBadge.badge);
            if (statusBadge.tooltip) {
                tooltip += statusBadge.tooltip + '\n';
            }
        }

        // Priority badge
        const priorityBadge = this.getPriorityBadge(task);
        if (priorityBadge) {
            badges.push(priorityBadge.badge);
            color = priorityBadge.color;
            if (priorityBadge.tooltip) {
                tooltip += priorityBadge.tooltip + '\n';
            }
        }

        // Due date badge
        const dueDateBadge = this.getDueDateBadge(task);
        if (dueDateBadge) {
            badges.push(dueDateBadge.badge);
            if (!color && dueDateBadge.color) {
                color = dueDateBadge.color;
            }
            if (dueDateBadge.tooltip) {
                tooltip += dueDateBadge.tooltip + '\n';
            }
        }

        // Dependency badge
        const depBadge = this.getDependencyBadge(task);
        if (depBadge) {
            badges.push(depBadge.badge);
            if (depBadge.tooltip) {
                tooltip += depBadge.tooltip + '\n';
            }
        }

        // Progress badge for parent tasks
        const progressBadge = this.getProgressBadge(task);
        if (progressBadge) {
            badges.push(progressBadge.badge);
            if (progressBadge.tooltip) {
                tooltip += progressBadge.tooltip + '\n';
            }
        }

        // Tag badges
        const tagBadges = this.getTagBadges(task);
        badges.push(...tagBadges);

        return {
            badge: badges.join(' '),
            color: color ? new vscode.ThemeColor(color) : undefined,
            tooltip: tooltip.trim(),
            propagate: false
        };
    }

    /**
     * Gets status badge based on task status
     */
    private getStatusBadge(task: Gorev): { badge: string; tooltip?: string } | undefined {
        if (task.durum === GorevDurum.Tamamlandi) {
            return { badge: '✓', tooltip: 'Tamamlandı' };
        }
        if (task.durum === GorevDurum.DevamEdiyor) {
            return { badge: '▶', tooltip: 'Devam ediyor' };
        }
        return undefined;
    }

    /**
     * Gets priority badge with appropriate icon and color
     */
    private getPriorityBadge(task: Gorev): { badge: string; color?: string; tooltip?: string } | undefined {
        if (task.durum === GorevDurum.Tamamlandi) {
            return undefined; // No priority badge for completed tasks
        }

        switch (task.oncelik) {
            case GorevOncelik.Yuksek:
                return { 
                    badge: '🔥', 
                    color: COLORS.HIGH_PRIORITY,
                    tooltip: 'Yüksek öncelik' 
                };
            case GorevOncelik.Orta:
                return { 
                    badge: '⚡', 
                    color: COLORS.MEDIUM_PRIORITY,
                    tooltip: 'Orta öncelik' 
                };
            case GorevOncelik.Dusuk:
                return { 
                    badge: 'ℹ', 
                    color: COLORS.LOW_PRIORITY,
                    tooltip: 'Düşük öncelik' 
                };
        }
    }

    /**
     * Gets due date badge with color coding
     */
    private getDueDateBadge(task: Gorev): { badge: string; color?: string; tooltip?: string } | undefined {
        if (!task.son_tarih || task.durum === GorevDurum.Tamamlandi) {
            return undefined;
        }

        const dueDate = new Date(task.son_tarih);
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const tomorrow = new Date(today);
        tomorrow.setDate(tomorrow.getDate() + 1);
        const dayAfterTomorrow = new Date(tomorrow);
        dayAfterTomorrow.setDate(dayAfterTomorrow.getDate() + 1);
        const nextWeek = new Date(today);
        nextWeek.setDate(nextWeek.getDate() + 7);

        // Calculate days difference
        const diffTime = dueDate.getTime() - today.getTime();
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (dueDate < today) {
            // Overdue
            return { 
                badge: `📅 ${Math.abs(diffDays)}g gecikmiş`, 
                color: 'errorForeground',
                tooltip: `${Math.abs(diffDays)} gün gecikmiş!` 
            };
        } else if (dueDate >= today && dueDate < tomorrow) {
            // Due today
            return { 
                badge: '📅 Bugün', 
                color: 'warningForeground',
                tooltip: 'Bugün teslim!' 
            };
        } else if (dueDate >= tomorrow && dueDate < dayAfterTomorrow) {
            // Due tomorrow
            return { 
                badge: '📅 Yarın', 
                color: 'warningForeground',
                tooltip: 'Yarın teslim' 
            };
        } else if (dueDate < nextWeek) {
            // Due this week
            return { 
                badge: `📅 ${diffDays}g`, 
                color: 'editorWarning.foreground',
                tooltip: `${diffDays} gün içinde teslim` 
            };
        } else {
            // Due later
            const weeks = Math.floor(diffDays / 7);
            if (weeks > 0) {
                return { 
                    badge: `📅 ${weeks}h`, 
                    tooltip: `${weeks} hafta içinde teslim` 
                };
            } else {
                return { 
                    badge: `📅 ${diffDays}g`, 
                    tooltip: `${diffDays} gün içinde teslim` 
                };
            }
        }
    }

    /**
     * Gets dependency badge showing blocked/blocking status
     */
    private getDependencyBadge(task: Gorev): { badge: string; tooltip?: string } | undefined {
        const badges: string[] = [];
        const tooltips: string[] = [];

        // Task has dependencies (blocked by other tasks)
        if (task.bagimli_gorev_sayisi && task.bagimli_gorev_sayisi > 0) {
            if (task.tamamlanmamis_bagimlilik_sayisi && task.tamamlanmamis_bagimlilik_sayisi > 0) {
                // Has incomplete dependencies - blocked
                badges.push(`🔒${task.tamamlanmamis_bagimlilik_sayisi}`);
                tooltips.push(`${task.tamamlanmamis_bagimlilik_sayisi} tamamlanmamış bağımlılık`);
            } else {
                // All dependencies completed
                badges.push(`🔓${task.bagimli_gorev_sayisi}`);
                tooltips.push(`Tüm ${task.bagimli_gorev_sayisi} bağımlılık tamamlandı`);
            }
        }

        // Other tasks depend on this (blocking others)
        if (task.bu_goreve_bagimli_sayisi && task.bu_goreve_bagimli_sayisi > 0) {
            badges.push(`🔗${task.bu_goreve_bagimli_sayisi}`);
            tooltips.push(`${task.bu_goreve_bagimli_sayisi} görev bunu bekliyor`);
        }

        if (badges.length > 0) {
            return {
                badge: badges.join(' '),
                tooltip: tooltips.join('\n')
            };
        }

        return undefined;
    }

    /**
     * Gets progress badge for parent tasks
     */
    private getProgressBadge(task: Gorev): { badge: string; tooltip?: string } | undefined {
        if (!task.alt_gorevler || task.alt_gorevler.length === 0) {
            return undefined;
        }

        const total = task.alt_gorevler.length;
        const completed = task.alt_gorevler.filter(t => t.durum === GorevDurum.Tamamlandi).length;
        const percentage = Math.round((completed / total) * 100);

        // Use different icons based on progress
        let icon: string;
        if (percentage === 0) {
            icon = '○';
        } else if (percentage < 25) {
            icon = '◔';
        } else if (percentage < 50) {
            icon = '◑';
        } else if (percentage < 75) {
            icon = '◕';
        } else if (percentage < 100) {
            icon = '◉';
        } else {
            icon = '●';
        }

        return {
            badge: `${icon} ${percentage}%`,
            tooltip: `Alt görevler: ${completed}/${total} tamamlandı (${percentage}%)`
        };
    }

    /**
     * Gets tag badges as colored pills
     */
    private getTagBadges(task: Gorev): string[] {
        if (!task.etiketler || task.etiketler.length === 0) {
            return [];
        }

        // Return first 3 tags as badges
        return task.etiketler.slice(0, 3).map(tag => `#${tag.isim}`);
    }

    /**
     * Creates a unique URI for a task
     */
    private getTaskUri(task: Gorev): vscode.Uri {
        return vscode.Uri.parse(`gorev:task/${task.id}`);
    }

    /**
     * Refreshes all decorations
     */
    refresh(): void {
        this._onDidChangeFileDecorations.fire([...this.decorations.keys()].map(uri => vscode.Uri.parse(uri)));
    }

    /**
     * Clears all decorations
     */
    clear(): void {
        this.decorations.clear();
        this.refresh();
    }
}
