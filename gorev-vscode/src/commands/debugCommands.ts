import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ClientInterface } from '../interfaces/client';
import { CommandContext } from './index';
import { TestDataSeeder } from '../debug/testDataSeeder';
import { Logger } from '../utils/logger';

export function registerDebugCommands(
    context: vscode.ExtensionContext,
    mcpClient: ClientInterface,
    providers: CommandContext
): void {
    const seeder = new TestDataSeeder(mcpClient);

    // Seed Test Data Command
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.debug.seedTestData', async () => {
            try {
                await seeder.seedTestData();
                // Tüm view'ları yenile
                await Promise.all([
                    providers.gorevTreeProvider.refresh(),
                    providers.projeTreeProvider.refresh(),
                    providers.templateTreeProvider?.refresh()
                ].filter(Boolean));
            } catch (error) {
                Logger.error('Failed to seed test data:', error);
                const errorMessage = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(t('debug.seedingFailed', errorMessage));
            }
        })
    );

    // Clear Test Data Command
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.debug.clearTestData', async () => {
            try {
                await seeder.clearTestData();
                // Tüm view'ları yenile
                await Promise.all([
                    providers.gorevTreeProvider.refresh(),
                    providers.projeTreeProvider.refresh(),
                    providers.templateTreeProvider?.refresh()
                ].filter(Boolean));
            } catch (error) {
                Logger.error('Failed to clear test data:', error);
                const errorMessage = error instanceof Error ? error.message : String(error);
                vscode.window.showErrorMessage(t('debug.clearingFailed', errorMessage));
            }
        })
    );

    // Debug modda olduğumuzu belirt
    vscode.commands.executeCommand('setContext', 'debugMode', true);
    
    // Status bar'a debug göstergesi ekle
    const debugStatusBar = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 1000);
    debugStatusBar.text = t('debug.modeLabel');
    debugStatusBar.tooltip = t('debug.modeTooltip');
    debugStatusBar.command = 'gorev.debug.seedTestData';
    debugStatusBar.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
    debugStatusBar.show();
    
    context.subscriptions.push(debugStatusBar);

    Logger.info(t('debug.commandsRegistered'));
}
