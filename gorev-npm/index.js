#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');
const os = require('os');

/**
 * Gorev MCP Server NPM Wrapper
 *
 * This wrapper automatically detects the platform and launches the appropriate
 * Gorev MCP server binary for use with MCP clients like Claude Desktop, VS Code, etc.
 */

function getPlatform() {
    const platform = os.platform();
    const arch = os.arch();

    switch (platform) {
        case 'win32':
            return arch === 'arm64' ? 'windows-arm64' : 'windows-amd64';
        case 'darwin':
            return arch === 'arm64' ? 'darwin-arm64' : 'darwin-amd64';
        case 'linux':
            return arch === 'arm64' ? 'linux-arm64' : 'linux-amd64';
        default:
            throw new Error(`Unsupported platform: ${platform}-${arch}`);
    }
}

function getBinaryName() {
    const platform = os.platform();
    return platform === 'win32' ? 'gorev.exe' : 'gorev';
}

function getBinaryPath() {
    const platformDir = getPlatform();
    const binaryName = getBinaryName();
    const binaryPath = path.join(__dirname, 'binaries', platformDir, binaryName);

    if (!fs.existsSync(binaryPath)) {
        throw new Error(`Binary not found for platform ${platformDir}: ${binaryPath}`);
    }

    return binaryPath;
}

function runServer() {
    try {
        const binaryPath = getBinaryPath();

        // Pass through all command line arguments except the first two (node and script path)
        const args = process.argv.slice(2);

        // Default to 'serve' command if no arguments provided (for MCP usage)
        if (args.length === 0) {
            args.push('serve');
        }

        // Spawn the Gorev binary
        const child = spawn(binaryPath, args, {
            stdio: 'inherit', // Pass through stdin, stdout, stderr
            env: {
                ...process.env,
                // Set default language to Turkish if not specified
                GOREV_LANG: process.env.GOREV_LANG || 'tr'
            }
        });

        // Handle process termination
        process.on('SIGINT', () => {
            child.kill('SIGINT');
        });

        process.on('SIGTERM', () => {
            child.kill('SIGTERM');
        });

        // Exit with the same code as the child process
        child.on('exit', (code, signal) => {
            if (signal) {
                process.kill(process.pid, signal);
            } else {
                process.exit(code);
            }
        });

        // Handle errors
        child.on('error', (err) => {
            console.error('Failed to start Gorev MCP server:', err.message);
            process.exit(1);
        });

    } catch (err) {
        console.error('Error:', err.message);
        console.error('');
        console.error('This usually means the Gorev binary is not properly installed.');
        console.error('Try reinstalling the package:');
        console.error('  npm uninstall @gorev/mcp-server');
        console.error('  npm install @gorev/mcp-server@latest');
        console.error('');
        console.error('If the problem persists, please report it at:');
        console.error('  https://github.com/msenol/Gorev/issues');
        process.exit(1);
    }
}

// Show help information
function showHelp() {
    console.log('Gorev MCP Server - Task Management for AI Assistants');
    console.log('');
    console.log('Usage:');
    console.log('  npx @gorev/mcp-server [command] [options]');
    console.log('');
    console.log('Commands:');
    console.log('  serve          Start MCP server (default)');
    console.log('  init           Initialize database');
    console.log('  template init  Initialize default templates');
    console.log('  --help, -h     Show this help message');
    console.log('  --version, -v  Show version');
    console.log('');
    console.log('Environment Variables:');
    console.log('  GOREV_LANG     Language (tr/en, default: tr)');
    console.log('  GOREV_DB_PATH  Custom database path');
    console.log('');
    console.log('MCP Configuration Example:');
    console.log('Add to your mcp.json:');
    console.log('{');
    console.log('  "mcpServers": {');
    console.log('    "gorev": {');
    console.log('      "command": "npx",');
    console.log('      "args": ["@gorev/mcp-server@latest"],');
    console.log('      "env": {');
    console.log('        "GOREV_LANG": "tr"');
    console.log('      }');
    console.log('    }');
    console.log('  }');
    console.log('}');
}

// Show version
function showVersion() {
    const packageJson = require('./package.json');
    console.log(packageJson.version);
}

// Main entry point
function main() {
    const args = process.argv.slice(2);

    if (args.includes('--help') || args.includes('-h')) {
        showHelp();
        return;
    }

    if (args.includes('--version') || args.includes('-v')) {
        showVersion();
        return;
    }

    runServer();
}

// Run if this script is executed directly
if (require.main === module) {
    main();
}

module.exports = {
    getPlatform,
    getBinaryPath,
    runServer
};