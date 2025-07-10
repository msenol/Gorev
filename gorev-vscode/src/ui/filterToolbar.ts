import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { GorevDurum, GorevOncelik } from '../models/common';
import { TaskFilter, SortingCriteria } from '../models/treeModels';
import { Logger } from '../utils/logger';

/**
 * Gelişmiş filtreleme toolbar'ı
 */
export class FilterToolbar {
    private statusBarItems: vscode.StatusBarItem[] = [];
    private quickPick: vscode.QuickPick<FilterQuickPickItem> | undefined;
    private activeFilters: TaskFilter = {};
    private savedProfiles: Map<string, TaskFilter> = new Map();

    constructor(
        private mcpClient: MCPClient,
        private onFilterChange: (filter: TaskFilter) => void
    ) {
        this.loadSavedProfiles();
        this.createStatusBarItems();
    }

    /**
     * Status bar öğelerini oluştur
     */
    private createStatusBarItems(): void {
        // Arama butonu
        const searchItem = vscode.window.createStatusBarItem(
            vscode.StatusBarAlignment.Left,
            100
        );
        searchItem.text = '$(search) Ara';
        searchItem.tooltip = 'Görevlerde ara';
        searchItem.command = 'gorev.showSearchInput';
        this.statusBarItems.push(searchItem);

        // Filtre butonu
        const filterItem = vscode.window.createStatusBarItem(
            vscode.StatusBarAlignment.Left,
            99
        );
        filterItem.text = '$(filter) Filtrele';
        filterItem.tooltip = 'Gelişmiş filtreleme';
        filterItem.command = 'gorev.showFilterMenu';
        this.statusBarItems.push(filterItem);

        // Aktif filtre göstergesi
        const activeFilterItem = vscode.window.createStatusBarItem(
            vscode.StatusBarAlignment.Left,
            98
        );
        activeFilterItem.text = '';
        activeFilterItem.tooltip = 'Aktif filtreler';
        activeFilterItem.command = 'gorev.clearAllFilters';
        this.statusBarItems.push(activeFilterItem);

        // Kayıtlı profiller
        const profileItem = vscode.window.createStatusBarItem(
            vscode.StatusBarAlignment.Left,
            97
        );
        profileItem.text = '$(bookmark) Profiller';
        profileItem.tooltip = 'Filtre profilleri';
        profileItem.command = 'gorev.showFilterProfiles';
        this.statusBarItems.push(profileItem);

        // Tüm projeler toggle
        const allProjectsItem = vscode.window.createStatusBarItem(
            vscode.StatusBarAlignment.Left,
            96
        );
        allProjectsItem.text = '$(globe) Tüm Projeler';
        allProjectsItem.tooltip = 'Tüm projeler/Aktif proje arasında geçiş yap';
        allProjectsItem.command = 'gorev.toggleAllProjects';
        this.statusBarItems.push(allProjectsItem);
    }

    /**
     * Toolbar'ı göster
     */
    show(): void {
        this.statusBarItems.forEach(item => item.show());
        this.updateAllProjectsIndicator();
    }

    /**
     * Toolbar'ı gizle
     */
    hide(): void {
        this.statusBarItems.forEach(item => item.hide());
    }

    /**
     * Arama input'unu göster
     */
    async showSearchInput(): Promise<void> {
        const searchQuery = await vscode.window.showInputBox({
            prompt: 'Görev başlığı veya açıklamasında ara',
            placeHolder: 'örn: bug fix, frontend, urgent',
            value: this.activeFilters.searchQuery || ''
        });

        if (searchQuery !== undefined) {
            this.updateFilter({ searchQuery: searchQuery || undefined });
        }
    }

    /**
     * Gelişmiş filtre menüsünü göster
     */
    async showFilterMenu(): Promise<void> {
        this.quickPick = vscode.window.createQuickPick<FilterQuickPickItem>();
        this.quickPick.title = 'Gelişmiş Filtreleme';
        this.quickPick.placeholder = 'Filtre seçin veya arama yapın';
        this.quickPick.canSelectMany = true;
        
        // Filtre seçeneklerini oluştur
        const items: FilterQuickPickItem[] = [
            // Durum filtreleri
            ...Object.values(GorevDurum).map(durum => ({
                label: `$(circle-outline) ${this.getDurumLabel(durum)}`,
                description: 'Durum',
                value: { durum },
                filterType: 'durum' as const,
                picked: this.activeFilters.durum === durum
            })),
            
            // Öncelik filtreleri
            ...Object.values(GorevOncelik).map(oncelik => ({
                label: `$(arrow-up) ${this.getOncelikLabel(oncelik)}`,
                description: 'Öncelik',
                value: { oncelik },
                filterType: 'oncelik' as const,
                picked: this.activeFilters.oncelik === oncelik
            })),
            
            // Özel filtreler
            {
                label: '$(globe) Tüm Projeler',
                description: 'Tüm projelerdeki görevleri göster',
                value: { showAllProjects: true },
                filterType: 'special' as const,
                picked: this.activeFilters.showAllProjects !== false
            },
            {
                label: '$(warning) Gecikmiş Görevler',
                description: 'Son tarihi geçmiş görevler',
                value: { overdue: true },
                filterType: 'special' as const,
                picked: this.activeFilters.overdue === true
            },
            {
                label: '$(calendar) Bugün Biten',
                description: 'Son tarihi bugün olan görevler',
                value: { dueToday: true },
                filterType: 'special' as const,
                picked: this.activeFilters.dueToday === true
            },
            {
                label: '$(calendar) Bu Hafta Biten',
                description: 'Son tarihi bu hafta olan görevler',
                value: { dueThisWeek: true },
                filterType: 'special' as const,
                picked: this.activeFilters.dueThisWeek === true
            },
            {
                label: '$(tag) Etiketli Görevler',
                description: 'En az bir etiketi olan görevler',
                value: { hasTag: true },
                filterType: 'special' as const,
                picked: this.activeFilters.hasTag === true
            },
            {
                label: '$(link) Bağımlılığı Olan',
                description: 'Bağımlılığı olan görevler',
                value: { hasDependency: true },
                filterType: 'special' as const,
                picked: this.activeFilters.hasDependency === true
            }
        ];

        // Proje listesini al ve filtre olarak ekle
        try {
            const projectResult = await this.mcpClient.callTool('proje_listele', {});
            const projects = this.parseProjects(projectResult.content[0].text);
            
            items.push(...projects.map(project => ({
                label: `$(folder) ${project.isim}`,
                description: 'Proje',
                value: { projeId: project.id },
                filterType: 'proje' as const,
                picked: this.activeFilters.projeId === project.id
            })));
        } catch (error) {
            Logger.error('Failed to load projects for filter:', error);
        }

        this.quickPick.items = items;
        this.quickPick.selectedItems = items.filter(item => item.picked);

        // Seçim değişikliklerini dinle
        this.quickPick.onDidChangeSelection(selection => {
            const newFilter: TaskFilter = {};
            
            selection.forEach(item => {
                switch (item.filterType) {
                    case 'durum':
                        newFilter.durum = item.value.durum;
                        break;
                    case 'oncelik':
                        newFilter.oncelik = item.value.oncelik;
                        break;
                    case 'proje':
                        newFilter.projeId = item.value.projeId;
                        break;
                    case 'special':
                        Object.assign(newFilter, item.value);
                        break;
                }
            });

            // Mevcut arama sorgusu varsa koru
            if (this.activeFilters.searchQuery) {
                newFilter.searchQuery = this.activeFilters.searchQuery;
            }

            this.updateFilter(newFilter);
        });

        // Butonlar ekle
        this.quickPick.buttons = [
            {
                iconPath: new vscode.ThemeIcon('save'),
                tooltip: 'Filtreyi Kaydet'
            },
            {
                iconPath: new vscode.ThemeIcon('clear-all'),
                tooltip: 'Tüm Filtreleri Temizle'
            }
        ];

        this.quickPick.onDidTriggerButton(button => {
            if (button.tooltip === 'Filtreyi Kaydet') {
                this.saveCurrentFilter();
            } else if (button.tooltip === 'Tüm Filtreleri Temizle') {
                this.clearAllFilters();
                this.quickPick?.hide();
            }
        });

        this.quickPick.onDidHide(() => {
            this.quickPick?.dispose();
            this.quickPick = undefined;
        });

        this.quickPick.show();
    }

    /**
     * Filtre profillerini göster
     */
    async showFilterProfiles(): Promise<void> {
        const profiles = Array.from(this.savedProfiles.entries()).map(([name, filter]) => ({
            label: name,
            description: this.getFilterDescription(filter),
            filter
        }));

        if (profiles.length === 0) {
            vscode.window.showInformationMessage('Kayıtlı filtre profili bulunmuyor.');
            return;
        }

        const selected = await vscode.window.showQuickPick(profiles, {
            placeHolder: 'Bir filtre profili seçin'
        });

        if (selected) {
            this.updateFilter(selected.filter);
            vscode.window.showInformationMessage(`"${selected.label}" filtre profili uygulandı.`);
        }
    }

    /**
     * Mevcut filtreyi kaydet
     */
    private async saveCurrentFilter(): Promise<void> {
        if (Object.keys(this.activeFilters).length === 0) {
            vscode.window.showWarningMessage('Kaydedilecek aktif filtre yok.');
            return;
        }

        const name = await vscode.window.showInputBox({
            prompt: 'Filtre profili adı',
            placeHolder: 'örn: Yüksek Öncelikli Buglar'
        });

        if (name) {
            this.savedProfiles.set(name, { ...this.activeFilters });
            this.saveSavedProfiles();
            vscode.window.showInformationMessage(`"${name}" filtre profili kaydedildi.`);
        }
    }

    /**
     * Filtreyi güncelle
     */
    private updateFilter(filter: TaskFilter): void {
        this.activeFilters = filter;
        this.updateActiveFilterDisplay();
        this.updateAllProjectsIndicator();
        this.onFilterChange(filter);
    }

    /**
     * Tüm filtreleri temizle
     */
    clearAllFilters(): void {
        this.activeFilters = {};
        this.onFilterChange({}); // Empty filter object to clear all filters
        this.updateActiveFilterDisplay();
        vscode.window.showInformationMessage('Tüm filtreler temizlendi.');
    }

    /**
     * Aktif filtre göstergesini güncelle
     */
    private updateActiveFilterDisplay(): void {
        const activeFilterItem = this.statusBarItems[2]; // Aktif filtre göstergesi
        const filterCount = Object.keys(this.activeFilters).length;

        if (filterCount > 0) {
            activeFilterItem.text = `$(filter-filled) ${filterCount} filtre aktif`;
            activeFilterItem.tooltip = `Aktif filtreler:\n${this.getFilterDescription(this.activeFilters)}\n\nTemizlemek için tıklayın`;
            activeFilterItem.show();
        } else {
            activeFilterItem.hide();
        }
    }

    /**
     * Filtre açıklamasını oluştur
     */
    private getFilterDescription(filter: TaskFilter): string {
        const parts: string[] = [];

        if (filter.searchQuery) {
            parts.push(`Arama: "${filter.searchQuery}"`);
        }
        if (filter.durum) {
            parts.push(`Durum: ${this.getDurumLabel(filter.durum)}`);
        }
        if (filter.oncelik) {
            parts.push(`Öncelik: ${this.getOncelikLabel(filter.oncelik)}`);
        }
        if (filter.projeId) {
            parts.push('Proje filtresi');
        }
        if (filter.overdue) {
            parts.push('Gecikmiş');
        }
        if (filter.dueToday) {
            parts.push('Bugün biten');
        }
        if (filter.dueThisWeek) {
            parts.push('Bu hafta biten');
        }
        if (filter.hasTag) {
            parts.push('Etiketli');
        }
        if (filter.hasDependency) {
            parts.push('Bağımlılığı olan');
        }

        return parts.join(' • ');
    }

    /**
     * Durum label'ı
     */
    private getDurumLabel(durum: GorevDurum): string {
        const labels: Record<GorevDurum, string> = {
            [GorevDurum.Beklemede]: 'Beklemede',
            [GorevDurum.DevamEdiyor]: 'Devam Ediyor',
            [GorevDurum.Tamamlandi]: 'Tamamlandı'
        };
        return labels[durum];
    }

    /**
     * Öncelik label'ı
     */
    private getOncelikLabel(oncelik: GorevOncelik): string {
        const labels: Record<GorevOncelik, string> = {
            [GorevOncelik.Dusuk]: 'Düşük',
            [GorevOncelik.Orta]: 'Orta',
            [GorevOncelik.Yuksek]: 'Yüksek'
        };
        return labels[oncelik];
    }

    /**
     * Tüm projeler göstergesini güncelle
     */
    private updateAllProjectsIndicator(): void {
        const allProjectsItem = this.statusBarItems[4]; // 5th item (index 4)
        if (allProjectsItem) {
            const showingAllProjects = this.activeFilters.showAllProjects !== false;
            allProjectsItem.text = showingAllProjects ? '$(globe) Tüm Projeler' : '$(folder) Aktif Proje';
            allProjectsItem.backgroundColor = showingAllProjects ? undefined : new vscode.ThemeColor('statusBarItem.warningBackground');
        }
    }

    /**
     * Projeleri parse et
     */
    private parseProjects(content: string): Array<{ id: string; isim: string }> {
        const projects: Array<{ id: string; isim: string }> = [];
        const lines = content.split('\n');
        
        for (const line of lines) {
            const match = line.match(/^## (.+) \(ID: ([^)]+)\)/);
            if (match) {
                projects.push({
                    isim: match[1],
                    id: match[2]
                });
            }
        }
        
        return projects;
    }

    /**
     * Kayıtlı profilleri yükle
     */
    private loadSavedProfiles(): void {
        const saved = vscode.workspace.getConfiguration('gorev').get<Record<string, TaskFilter>>('filterProfiles');
        if (saved) {
            this.savedProfiles = new Map(Object.entries(saved));
        }
    }

    /**
     * Kayıtlı profilleri kaydet
     */
    private saveSavedProfiles(): void {
        const profiles = Object.fromEntries(this.savedProfiles);
        vscode.workspace.getConfiguration('gorev').update('filterProfiles', profiles, true);
    }

    /**
     * Dispose
     */
    dispose(): void {
        this.statusBarItems.forEach(item => item.dispose());
        this.quickPick?.dispose();
    }
}

/**
 * Filter quick pick item tipi
 */
interface FilterQuickPickItem extends vscode.QuickPickItem {
    value: Partial<TaskFilter>;
    filterType: 'durum' | 'oncelik' | 'proje' | 'special';
}