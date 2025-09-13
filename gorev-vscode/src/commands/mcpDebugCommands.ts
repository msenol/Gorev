import * as vscode from 'vscode';
import { t } from '../utils/l10n';
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
                t('mcpDebug.toggleMessage', !currentValue ? t('mcpDebug.enabled') : t('mcpDebug.disabled'))
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
                    vscode.window.showInformationMessage(t('mcpDebug.noLogsFound'));
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
                        placeHolder: t('mcpDebug.selectLogFile')
                    }
                );
                
                if (selected) {
                    const uri = vscode.Uri.file(path.join(debugPath, selected.label));
                    await vscode.commands.executeCommand('vscode.open', uri);
                }
            } catch (e) {
                const errorMessage = e instanceof Error ? e.message : String(e);
                vscode.window.showErrorMessage(t('mcpDebug.readLogsFailed', errorMessage));
            }
        })
    );
    
    // Clear debug logs
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.clearDebugLogs', async () => {
            const debugConfig = getDebugConfig();
            const debugPath = debugConfig.debugLogPath;
            
            const answer = await vscode.window.showWarningMessage(
                t('mcpDebug.clearLogsConfirm'),
                t('mcpDebug.yes'),
                t('mcpDebug.no')
            );
            
            if (answer === t('mcpDebug.yes')) {
                try {
                    const files = fs.readdirSync(debugPath)
                        .filter(f => f.startsWith('mcp-session-') || f.startsWith('stdin-') || f.startsWith('stdout-') || f.endsWith('.log'));
                    
                    for (const file of files) {
                        fs.unlinkSync(path.join(debugPath, file));
                    }
                    
                    vscode.window.showInformationMessage(t('mcpDebug.clearedLogs', files.length.toString()));
                } catch (e) {
                    const errorMessage = e instanceof Error ? e.message : String(e);
                    vscode.window.showErrorMessage(t('mcpDebug.clearLogsFailed', errorMessage));
                }
            }
        })
    );
    
    // Test MCP connection
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.testConnection', async () => {
            outputChannel.show();
            outputChannel.appendLine(t('mcpDebug.testingConnection'));
            outputChannel.appendLine(t('mcpDebug.testTime', new Date().toISOString()));
            
            const config = vscode.workspace.getConfiguration('gorev');
            const serverPath = config.get<string>('mcp.serverPath') || config.get<string>('serverPath');
            
            if (!serverPath) {
                outputChannel.appendLine(t('mcpDebug.noServerPath'));
                vscode.window.showErrorMessage(t('mcpDebug.configureServerPath'));
                return;
            }
            
            outputChannel.appendLine(t('mcpDebug.serverPath', serverPath));
            
            // Check if file exists
            if (!fs.existsSync(serverPath)) {
                outputChannel.appendLine(t('mcpDebug.serverNotFound'));
                vscode.window.showErrorMessage(t('mcpDebug.serverNotFoundAt', serverPath));
                return;
            }
            
            // Check if executable
            try {
                fs.accessSync(serverPath, fs.constants.X_OK);
                outputChannel.appendLine(t('mcpDebug.serverExecutable'));
            } catch {
                outputChannel.appendLine(t('mcpDebug.serverNotExecutable'));
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
                    outputChannel.appendLine(t('mcpDebug.serverResponded'));
                    outputChannel.appendLine(t('mcpDebug.responsePreview', response.substring(0, 200)));
                } else {
                    outputChannel.appendLine(t('mcpDebug.noResponse'));
                }
                
                if (error) {
                    outputChannel.appendLine(t('mcpDebug.stderr', error));
                }
                
            } catch (e) {
                const errorMessage = e instanceof Error ? e.message : String(e);
                outputChannel.appendLine(t('mcpDebug.error', errorMessage));
            }
            
            outputChannel.appendLine(t('mcpDebug.testComplete'));
        })
    );
}

function getFileDescription(filename: string): string {
    if (filename.startsWith('mcp-session-')) {
        return t('mcpDebug.mainDebugLog');
    } else if (filename.startsWith('stdin-')) {
        return t('mcpDebug.inputMessages');
    } else if (filename.startsWith('stdout-')) {
        return t('mcpDebug.outputMessages');
    }
    return t('mcpDebug.debugLog');
}

function getFileSize(filepath: string): string {
    try {
        const stats = fs.statSync(filepath);
        const size = stats.size;
        if (size < 1024) {
            return t('mcpDebug.bytes', size.toString());
        } else if (size < 1024 * 1024) {
            return t('mcpDebug.kilobytes', (size / 1024).toFixed(1));
        } else {
            return t('mcpDebug.megabytes', (size / (1024 * 1024)).toFixed(1));
        }
    } catch {
        return '';
    }
}
