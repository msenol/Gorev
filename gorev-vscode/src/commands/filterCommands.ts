import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { Logger } from '../utils/logger';

export function registerFilterCommands(
    context: vscode.ExtensionContext,
    mcpClient: MCPClient,
    providers: CommandContext
): void {
    const { filterToolbar } = providers;
    
    if (!filterToolbar) {
        Logger.warn(vscode.l10n.t('filter.toolbarNotInitialized'));
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
                vscode.l10n.t('filter.showingActiveProject') : 
                vscode.l10n.t('filter.showingAllProjects');
            vscode.window.showInformationMessage(message);
        })
    );

    // Önceden tanımlanmış filtreler için kısayol komutları
    
    // Gecikmiş görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterOverdue', () => {
            providers.gorevTreeProvider.updateFilter({ overdue: true });
            vscode.window.showInformationMessage(vscode.l10n.t('filter.showingOverdue'));
        })
    );

    // Bugün biten görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterDueToday', () => {
            providers.gorevTreeProvider.updateFilter({ dueToday: true });
            vscode.window.showInformationMessage(vscode.l10n.t('filter.showingDueToday'));
        })
    );

    // Bu hafta biten görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterDueThisWeek', () => {
            providers.gorevTreeProvider.updateFilter({ dueThisWeek: true });
            vscode.window.showInformationMessage(vscode.l10n.t('filter.showingDueThisWeek'));
        })
    );

    // Yüksek öncelikli görevler
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterHighPriority', () => {
            providers.gorevTreeProvider.updateFilter({ oncelik: 'yuksek' as any });
            vscode.window.showInformationMessage(vscode.l10n.t('filter.showingHighPriority'));
        })
    );

    // Aktif proje görevleri
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterActiveProject', async () => {
            try {
                const result = await mcpClient.callTool('aktif_proje_goster', {});
                const content = result.content[0].text;
                
                // Parse active project ID
                const match = content.match(/ID: ([a-f0-9-]+)/);
                if (match) {
                    providers.gorevTreeProvider.updateFilter({ projeId: match[1] });
                    vscode.window.showInformationMessage(vscode.l10n.t('filter.showingActiveProject'));
                } else {
                    vscode.window.showWarningMessage(vscode.l10n.t('filter.activeProjectNotFound'));
                }
            } catch (error) {
                vscode.window.showErrorMessage(vscode.l10n.t('filter.activeProjectFetchFailed'));
            }
        })
    );

    // Etiket bazlı filtreleme
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.filterByTag', async () => {
            const tag = await vscode.window.showInputBox({
                prompt: vscode.l10n.t('filter.tagPromptMessage'),
                placeHolder: vscode.l10n.t('filter.tagPlaceholderExample')
            });

            if (tag) {
                providers.gorevTreeProvider.updateFilter({ tags: [tag] });
                vscode.window.showInformationMessage(vscode.l10n.t('filter.showingTagged', tag));
            }
        })
    );
}