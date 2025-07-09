import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';

export interface DebugConfig {
    useDebugWrapper: boolean;
    debugLogPath: string;
    serverTimeout: number;
}

export function getDebugConfig(): DebugConfig {
    const config = vscode.workspace.getConfiguration('gorev');
    
    return {
        useDebugWrapper: config.get<boolean>('debug.useWrapper', false),
        debugLogPath: config.get<string>('debug.logPath', '/tmp/gorev-debug'),
        serverTimeout: config.get<number>('debug.serverTimeout', 5000)
    };
}

export function getServerPath(): string {
    const config = vscode.workspace.getConfiguration('gorev');
    const debugConfig = getDebugConfig();
    
    // If debug wrapper is enabled, use it
    if (debugConfig.useDebugWrapper) {
        const wrapperPath = path.join(
            path.dirname(config.get<string>('mcp.serverPath') || ''),
            '..',
            'debug-wrapper.sh'
        );
        
        if (fs.existsSync(wrapperPath)) {
            console.log(`[Gorev] Using debug wrapper: ${wrapperPath}`);
            return wrapperPath;
        } else {
            console.warn(`[Gorev] Debug wrapper not found at: ${wrapperPath}`);
        }
    }
    
    // Fall back to normal server path
    return config.get<string>('mcp.serverPath') || '';
}

export function showDebugInfo(outputChannel: vscode.OutputChannel): void {
    const debugConfig = getDebugConfig();
    
    if (debugConfig.useDebugWrapper) {
        outputChannel.appendLine('=== Debug Mode Enabled ===');
        outputChannel.appendLine(`Debug logs will be written to: ${debugConfig.debugLogPath}`);
        outputChannel.appendLine(`Server timeout: ${debugConfig.serverTimeout}ms`);
        outputChannel.appendLine('');
        
        // Show latest debug log if available
        try {
            const files = fs.readdirSync(debugConfig.debugLogPath)
                .filter(f => f.startsWith('mcp-session-'))
                .sort()
                .reverse();
            
            if (files.length > 0) {
                const latestLog = path.join(debugConfig.debugLogPath, files[0]);
                outputChannel.appendLine(`Latest debug log: ${latestLog}`);
                
                // Also create a command to open the log
                vscode.commands.executeCommand('vscode.open', vscode.Uri.file(latestLog));
            }
        } catch (e) {
            // Debug directory might not exist yet
        }
    }
}