#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');

// Common VS Code log locations
const vscodeLogPaths = [
    // Windows
    path.join(os.homedir(), 'AppData', 'Roaming', 'Code', 'logs'),
    // macOS
    path.join(os.homedir(), 'Library', 'Application Support', 'Code', 'logs'),
    // Linux
    path.join(os.homedir(), '.config', 'Code', 'logs'),
    // WSL
    '/mnt/c/Users/*/AppData/Roaming/Code/logs'
];

console.log('=== VS Code Extension Debug Tool ===\n');

// Function to find the most recent extension host log
function findExtensionLogs() {
    for (const logPath of vscodeLogPaths) {
        if (logPath.includes('*')) {
            // Handle wildcard paths
            const baseDir = path.dirname(logPath);
            const pattern = path.basename(logPath);
            
            try {
                const dirs = fs.readdirSync(baseDir);
                for (const dir of dirs) {
                    const fullPath = path.join(baseDir, dir, 'logs');
                    if (fs.existsSync(fullPath)) {
                        return scanLogsDirectory(fullPath);
                    }
                }
            } catch (err) {
                // Continue to next path
            }
        } else if (fs.existsSync(logPath)) {
            return scanLogsDirectory(logPath);
        }
    }
    return null;
}

function scanLogsDirectory(logsDir) {
    try {
        const entries = fs.readdirSync(logsDir, { withFileTypes: true });
        const extHostLogs = [];
        
        for (const entry of entries) {
            if (entry.isDirectory() && entry.name.includes('exthost')) {
                const logFile = path.join(logsDir, entry.name, 'exthost.log');
                if (fs.existsSync(logFile)) {
                    const stats = fs.statSync(logFile);
                    extHostLogs.push({ path: logFile, mtime: stats.mtime });
                }
            }
        }
        
        // Sort by most recent
        extHostLogs.sort((a, b) => b.mtime - a.mtime);
        return extHostLogs[0]?.path;
    } catch (err) {
        return null;
    }
}

// Find and read VS Code logs
const logFile = findExtensionLogs();
if (logFile) {
    console.log(`Found VS Code extension log: ${logFile}\n`);
    
    try {
        const content = fs.readFileSync(logFile, 'utf8');
        const lines = content.split('\n');
        
        // Filter for Gorev extension logs
        const gorevLogs = lines.filter(line => 
            line.includes('gorev') || 
            line.includes('Gorev') || 
            line.includes('MCP') ||
            line.includes('task') ||
            line.includes('Task')
        );
        
        console.log(`Found ${gorevLogs.length} Gorev-related log entries:\n`);
        
        // Show last 50 entries
        const recentLogs = gorevLogs.slice(-50);
        recentLogs.forEach(log => console.log(log));
        
    } catch (err) {
        console.error(`Error reading log file: ${err.message}`);
    }
} else {
    console.log('Could not find VS Code extension logs.');
    console.log('Make sure VS Code is running and the extension is loaded.');
}

// Check if output channel logs exist
const outputChannelPath = path.join(process.cwd(), '.vscode', 'gorev-output.log');
if (fs.existsSync(outputChannelPath)) {
    console.log('\n\n=== Gorev Output Channel Logs ===\n');
    try {
        const content = fs.readFileSync(outputChannelPath, 'utf8');
        console.log(content);
    } catch (err) {
        console.error(`Error reading output channel log: ${err.message}`);
    }
}

console.log('\n\n=== Next Steps ===');
console.log('1. Add detailed logging to enhancedGorevTreeProvider.ts');
console.log('2. Reload VS Code window (Ctrl+R or Cmd+R)');
console.log('3. Run this script again to see the new logs');