import * as vscode from 'vscode';
import * as fs from 'fs';
import * as path from 'path';
import { getDebugConfig, showDebugInfo } from '../debug/debugConfig';

export function registerMCPDebugCommands(context: vscode.ExtensionContext, outputChannel: vscode.OutputChannel) {
    // Toggle debug mode
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.toggleDebugMode', async () => {
            const config = vscode.workspace.getConfiguration('gorev');
            const currentValue = config.get<boolean>('debug.useWrapper', false);
            
            await config.update('debug.useWrapper', !currentValue, vscode.ConfigurationTarget.Workspace);
            
            vscode.window.showInformationMessage(
                vscode.l10n.t('mcpDebug.toggleMessage', !currentValue ? vscode.l10n.t('mcpDebug.enabled') : vscode.l10n.t('mcpDebug.disabled'))
            );
            
            if (!currentValue) {
                showDebugInfo(outputChannel);
            }
        })
    );
    
    // Show debug logs
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.showDebugLogs', async () => {
            const debugConfig = getDebugConfig();
            const debugPath = debugConfig.debugLogPath;
            
            try {
                const files = fs.readdirSync(debugPath)
                    .filter(f => f.startsWith('mcp-session-') || f.startsWith('stdin-') || f.startsWith('stdout-'))
                    .sort()
                    .reverse();
                
                if (files.length === 0) {
                    vscode.window.showInformationMessage(vscode.l10n.t('mcpDebug.noLogsFound'));
                    return;
                }
                
                // Show quick pick
                const selected = await vscode.window.showQuickPick(
                    files.map(f => ({
                        label: f,
                        description: getFileDescription(f),
                        detail: getFileSize(path.join(debugPath, f))
                    })),
                    {
                        placeHolder: vscode.l10n.t('mcpDebug.selectLogFile')
                    }
                );
                
                if (selected) {
                    const uri = vscode.Uri.file(path.join(debugPath, selected.label));
                    await vscode.commands.executeCommand('vscode.open', uri);
                }
            } catch (e) {
                vscode.window.showErrorMessage(vscode.l10n.t('mcpDebug.readLogsFailed', e.toString()));
            }
        })
    );
    
    // Clear debug logs
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.clearDebugLogs', async () => {
            const debugConfig = getDebugConfig();
            const debugPath = debugConfig.debugLogPath;
            
            const answer = await vscode.window.showWarningMessage(
                vscode.l10n.t('mcpDebug.clearLogsConfirm'),
                vscode.l10n.t('mcpDebug.yes'),
                vscode.l10n.t('mcpDebug.no')
            );
            
            if (answer === vscode.l10n.t('mcpDebug.yes')) {
                try {
                    const files = fs.readdirSync(debugPath)
                        .filter(f => f.startsWith('mcp-session-') || f.startsWith('stdin-') || f.startsWith('stdout-') || f.endsWith('.log'));
                    
                    for (const file of files) {
                        fs.unlinkSync(path.join(debugPath, file));
                    }
                    
                    vscode.window.showInformationMessage(vscode.l10n.t('mcpDebug.clearedLogs', files.length.toString()));
                } catch (e) {
                    vscode.window.showErrorMessage(vscode.l10n.t('mcpDebug.clearLogsFailed', e.toString()));
                }
            }
        })
    );
    
    // Test MCP connection
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.testConnection', async () => {
            outputChannel.show();
            outputChannel.appendLine(vscode.l10n.t('mcpDebug.testingConnection'));
            outputChannel.appendLine(vscode.l10n.t('mcpDebug.testTime', new Date().toISOString()));
            
            const config = vscode.workspace.getConfiguration('gorev');
            const serverPath = config.get<string>('mcp.serverPath') || config.get<string>('serverPath');
            
            if (!serverPath) {
                outputChannel.appendLine(vscode.l10n.t('mcpDebug.noServerPath'));
                vscode.window.showErrorMessage(vscode.l10n.t('mcpDebug.configureServerPath'));
                return;
            }
            
            outputChannel.appendLine(vscode.l10n.t('mcpDebug.serverPath', serverPath));
            
            // Check if file exists
            if (!fs.existsSync(serverPath)) {
                outputChannel.appendLine(vscode.l10n.t('mcpDebug.serverNotFound'));
                vscode.window.showErrorMessage(vscode.l10n.t('mcpDebug.serverNotFoundAt', serverPath));
                return;
            }
            
            // Check if executable
            try {
                fs.accessSync(serverPath, fs.constants.X_OK);
                outputChannel.appendLine(vscode.l10n.t('mcpDebug.serverExecutable'));
            } catch {
                outputChannel.appendLine(vscode.l10n.t('mcpDebug.serverNotExecutable'));
            }
            
            // Try a simple MCP call
            try {
                const { spawn } = require('child_process');
                const testProcess = spawn(serverPath, ['serve'], {
                    stdio: ['pipe', 'pipe', 'pipe']
                });
                
                let response = '';
                let error = '';
                
                testProcess.stdout.on('data', (data: Buffer) => {
                    response += data.toString();
                });
                
                testProcess.stderr.on('data', (data: Buffer) => {
                    error += data.toString();
                });
                
                // Send initialize request
                const initRequest = JSON.stringify({
                    jsonrpc: '2.0',
                    id: 1,
                    method: 'initialize',
                    params: { capabilities: {} }
                }) + '\n';
                
                testProcess.stdin.write(initRequest);
                
                // Wait for response
                await new Promise((resolve) => {
                    setTimeout(() => {
                        testProcess.kill();
                        resolve(null);
                    }, 3000);
                });
                
                if (response) {
                    outputChannel.appendLine(vscode.l10n.t('mcpDebug.serverResponded'));
                    outputChannel.appendLine(vscode.l10n.t('mcpDebug.responsePreview', response.substring(0, 200)));
                } else {
                    outputChannel.appendLine(vscode.l10n.t('mcpDebug.noResponse'));
                }
                
                if (error) {
                    outputChannel.appendLine(vscode.l10n.t('mcpDebug.stderr', error));
                }
                
            } catch (e) {
                outputChannel.appendLine(vscode.l10n.t('mcpDebug.error', e.toString()));
            }
            
            outputChannel.appendLine(vscode.l10n.t('mcpDebug.testComplete'));
        })
    );
}

function getFileDescription(filename: string): string {
    if (filename.startsWith('mcp-session-')) {
        return vscode.l10n.t('mcpDebug.mainDebugLog');
    } else if (filename.startsWith('stdin-')) {
        return vscode.l10n.t('mcpDebug.inputMessages');
    } else if (filename.startsWith('stdout-')) {
        return vscode.l10n.t('mcpDebug.outputMessages');
    }
    return vscode.l10n.t('mcpDebug.debugLog');
}

function getFileSize(filepath: string): string {
    try {
        const stats = fs.statSync(filepath);
        const size = stats.size;
        if (size < 1024) {
            return vscode.l10n.t('mcpDebug.bytes', size.toString());
        } else if (size < 1024 * 1024) {
            return vscode.l10n.t('mcpDebug.kilobytes', (size / 1024).toFixed(1));
        } else {
            return vscode.l10n.t('mcpDebug.megabytes', (size / (1024 * 1024)).toFixed(1));
        }
    } catch {
        return '';
    }
}