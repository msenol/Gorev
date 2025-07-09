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
                `Gorev debug mode ${!currentValue ? 'enabled' : 'disabled'}. Please restart VS Code for changes to take effect.`
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
                    vscode.window.showInformationMessage('No debug logs found. Enable debug mode and restart VS Code.');
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
                        placeHolder: 'Select a debug log to view'
                    }
                );
                
                if (selected) {
                    const uri = vscode.Uri.file(path.join(debugPath, selected.label));
                    await vscode.commands.executeCommand('vscode.open', uri);
                }
            } catch (e) {
                vscode.window.showErrorMessage(`Failed to read debug logs: ${e}`);
            }
        })
    );
    
    // Clear debug logs
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.clearDebugLogs', async () => {
            const debugConfig = getDebugConfig();
            const debugPath = debugConfig.debugLogPath;
            
            const answer = await vscode.window.showWarningMessage(
                'Are you sure you want to clear all debug logs?',
                'Yes',
                'No'
            );
            
            if (answer === 'Yes') {
                try {
                    const files = fs.readdirSync(debugPath)
                        .filter(f => f.startsWith('mcp-session-') || f.startsWith('stdin-') || f.startsWith('stdout-') || f.endsWith('.log'));
                    
                    for (const file of files) {
                        fs.unlinkSync(path.join(debugPath, file));
                    }
                    
                    vscode.window.showInformationMessage(`Cleared ${files.length} debug log files.`);
                } catch (e) {
                    vscode.window.showErrorMessage(`Failed to clear debug logs: ${e}`);
                }
            }
        })
    );
    
    // Test MCP connection
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.testConnection', async () => {
            outputChannel.show();
            outputChannel.appendLine('=== Testing MCP Connection ===');
            outputChannel.appendLine(`Time: ${new Date().toISOString()}`);
            
            const config = vscode.workspace.getConfiguration('gorev');
            const serverPath = config.get<string>('mcp.serverPath') || config.get<string>('serverPath');
            
            if (!serverPath) {
                outputChannel.appendLine('ERROR: No server path configured');
                vscode.window.showErrorMessage('Please configure the Gorev server path in settings.');
                return;
            }
            
            outputChannel.appendLine(`Server path: ${serverPath}`);
            
            // Check if file exists
            if (!fs.existsSync(serverPath)) {
                outputChannel.appendLine('ERROR: Server file does not exist');
                vscode.window.showErrorMessage(`Server not found at: ${serverPath}`);
                return;
            }
            
            // Check if executable
            try {
                fs.accessSync(serverPath, fs.constants.X_OK);
                outputChannel.appendLine('✓ Server file is executable');
            } catch {
                outputChannel.appendLine('WARNING: Server file may not be executable');
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
                    outputChannel.appendLine('✓ Server responded to initialize request');
                    outputChannel.appendLine(`Response: ${response.substring(0, 200)}...`);
                } else {
                    outputChannel.appendLine('✗ No response from server');
                }
                
                if (error) {
                    outputChannel.appendLine(`STDERR: ${error}`);
                }
                
            } catch (e) {
                outputChannel.appendLine(`ERROR: ${e}`);
            }
            
            outputChannel.appendLine('=== Test Complete ===');
        })
    );
}

function getFileDescription(filename: string): string {
    if (filename.startsWith('mcp-session-')) {
        return 'Main debug log';
    } else if (filename.startsWith('stdin-')) {
        return 'Input messages (VS Code → Server)';
    } else if (filename.startsWith('stdout-')) {
        return 'Output messages (Server → VS Code)';
    }
    return 'Debug log';
}

function getFileSize(filepath: string): string {
    try {
        const stats = fs.statSync(filepath);
        const size = stats.size;
        if (size < 1024) {
            return `${size} bytes`;
        } else if (size < 1024 * 1024) {
            return `${(size / 1024).toFixed(1)} KB`;
        } else {
            return `${(size / (1024 * 1024)).toFixed(1)} MB`;
        }
    } catch {
        return '';
    }
}