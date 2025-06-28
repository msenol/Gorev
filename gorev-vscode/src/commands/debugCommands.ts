import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { CommandContext } from './index';
import { TestDataSeeder } from '../debug/testDataSeeder';
import { Logger } from '../utils/logger';

export function registerDebugCommands(
    context: vscode.ExtensionContext,
    mcpClient: MCPClient,
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
                vscode.window.showErrorMessage(`Test data seeding failed: ${error}`);
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
                vscode.window.showErrorMessage(`Test data clearing failed: ${error}`);
            }
        })
    );

    // Debug modda olduğumuzu belirt
    vscode.commands.executeCommand('setContext', 'debugMode', true);
    
    // Status bar'a debug göstergesi ekle
    const debugStatusBar = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 1000);
    debugStatusBar.text = '$(beaker) Debug Mode';
    debugStatusBar.tooltip = 'Gorev Debug Mode Active\nClick to seed test data';
    debugStatusBar.command = 'gorev.debug.seedTestData';
    debugStatusBar.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
    debugStatusBar.show();
    
    context.subscriptions.push(debugStatusBar);

    Logger.info('Debug commands registered - Test data seeding available');
}