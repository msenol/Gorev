import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ApiClient } from '../api/client';
import { CommandContext } from './index';
import { Logger } from '../utils/logger';

export function registerFilterCommands(
    context: vscode.ExtensionContext,
    apiClient: ApiClient,
    providers: CommandContext
): void {
    const { filterToolbar } = providers;
    
    if (!filterToolbar) {
        Logger.warn(t('filter.toolbarNotInitialized'));
        return;
    }

    // Arama komutu
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.showSearchInput', async () => {
            await filterToolbar.showSearchInput();
        })
    );

    // Filtre menüsü komutu
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.showFilterMenu', async () => {
            await filterToolbar.showFilterMenu();
        })
    );

    // Filtre profilleri komutu
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.showFilterProfiles', async () => {
            await filterToolbar.showFilterProfiles();
        })
    );

    // Tüm filtreleri temizle komutu
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.clearAllFilters', () => {
            filterToolbar.clearAllFilters();
        })
    );

    // Tüm projeler / Aktif proje toggle komutu
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.toggleAllProjects', () => {
            const currentFilter = providers.gorevTreeProvider.getFilter();
            const showingAllProjects = currentFilter?.showAllProjects !== false;
            
            providers.gorevTreeProvider.updateFilter({ 
                ...currentFilter,
                showAllProjects: !showingAllProjects 
            });
            
            const message = showingAllProjects ? 
                t('filter.showingActiveProject') : 
                t('filter.showingAllProjects');
            vscode.window.showInformationMessage(message);
        })
    );

    // Önceden tanımlanmış filtreler için kısayol komutları
    
    // Gecikmiş görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterOverdue', () => {
            providers.gorevTreeProvider.updateFilter({ overdue: true });
            vscode.window.showInformationMessage(t('filter.showingOverdue'));
        })
    );

    // Bugün biten görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterDueToday', () => {
            providers.gorevTreeProvider.updateFilter({ dueToday: true });
            vscode.window.showInformationMessage(t('filter.showingDueToday'));
        })
    );

    // Bu hafta biten görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterDueThisWeek', () => {
            providers.gorevTreeProvider.updateFilter({ dueThisWeek: true });
            vscode.window.showInformationMessage(t('filter.showingDueThisWeek'));
        })
    );

    // Yüksek öncelikli görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterHighPriority', () => {
            providers.gorevTreeProvider.updateFilter({ oncelik: 'yuksek' as any });
            vscode.window.showInformationMessage(t('filter.showingHighPriority'));
        })
    );

    // Aktif proje görevleri
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterActiveProject', async () => {
            try {
                const result = await apiClient.getActiveProject();

                if (result.success && result.data) {
                    providers.gorevTreeProvider.updateFilter({ projeId: result.data.id });
                    vscode.window.showInformationMessage(t('filter.showingActiveProject'));
                } else {
                    vscode.window.showWarningMessage(t('filter.activeProjectNotFound'));
                }
            } catch (error) {
                vscode.window.showErrorMessage(t('filter.activeProjectFetchFailed'));
            }
        })
    );

    // Etiket bazlı filtreleme
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterByTag', async () => {
            const tag = await vscode.window.showInputBox({
                prompt: t('filter.tagPromptMessage'),
                placeHolder: t('filter.tagPlaceholderExample')
            });

            if (tag) {
                providers.gorevTreeProvider.updateFilter({ tags: [tag] });
                vscode.window.showInformationMessage(t('filter.showingTagged', tag));
            }
        })
    );
}
