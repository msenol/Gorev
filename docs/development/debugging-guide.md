# Debugging Guide

This guide helps you debug issues with the Gorev VS Code extension and MCP server.

## Common Issues and Solutions

### Issue: Task and Project Lists Not Showing

#### Step 1: Clean Database

```bash
cd gorev-mcpserver
rm -f gorev.db
```

#### Step 2: Start MCP Server

```bash
cd gorev-mcpserver
./gorev serve
```

Keep the server running.

#### Step 3: Debug VS Code Extension

1. Open the `gorev-vscode` folder in VS Code
2. Press F5 to launch Extension Development Host
3. The extension will automatically try to connect to the MCP server

#### Step 4: Create Test Data

In the extension:

1. Open Command Palette (Ctrl+Shift+P)
2. Run "Gorev Debug: Seed Test Data" command
3. Test data will be created

#### Step 5: Verify

1. Click on the Gorev icon in the left activity bar
2. Check Tasks and Projects panels
3. If lists are empty, click the refresh button

## Examining Debug Logs

### 1. Extension Host Output

- View > Output
- Select "Extension Host" from dropdown
- Look for Gorev-related logs

### 2. Developer Tools

- Help > Toggle Developer Tools
- Check Console tab for error messages

### 3. MCP Server Logs

- Check the MCP server terminal for error messages

## Troubleshooting

### "Not connected to MCP server" Message

- Check if `gorev.serverPath` is correct in settings.json
- Verify MCP server is running
- On Windows, there might be WSL path conversion issues

### Parser Errors

- The markdown format from MCP server might have changed
- Check parse errors in Output panel
- Update regex patterns in `src/utils/markdownParser.ts`

### Empty TreeView

- Check if MCP connection exists
- Verify `gorev_listele` and `proje_listele` commands work
- Test if parser works correctly

## Testing the Parser

```bash
cd gorev-vscode
npx ts-node -e "
import { MarkdownParser } from './src/utils/markdownParser';

// Test markdown
const testMd = \`## Test Task
- **ID:** 123
- **Status:** beklemede
- **Priority:** yuksek\`;

const tasks = MarkdownParser.parseGorevListesi(testMd);
console.log('Parsed tasks:', tasks);
"
```

## VS Code Extension Debugging Tips

### 1. Enable Verbose Logging

Add to your launch.json:

```json
{
  "type": "extensionHost",
  "request": "launch",
  "name": "Run Extension",
  "args": ["--extensionDevelopmentPath=${workspaceFolder}", "--log=verbose"]
}
```

### 2. Use Breakpoints

- Set breakpoints in TypeScript files
- Use conditional breakpoints for specific scenarios
- Use logpoints for non-intrusive debugging

### 3. Test Specific Commands

```typescript
// In debug console
await vscode.commands.executeCommand('gorev.createTask');
```

### 4. Check Extension Context

```typescript
// In debug console
vscode.extensions.getExtension('mehmetsenol.gorev-vscode')
```

## MCP Server Debugging

### 1. Enable Debug Mode

```bash
./gorev serve --debug
```

### 2. Check Database

```bash
sqlite3 gorev.db
.tables
SELECT * FROM gorevler;
```

### 3. Test MCP Tools Directly

Use the MCP protocol to test tools directly:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "gorev_listele",
    "arguments": {}
  },
  "id": 1
}
```

## Performance Debugging

### 1. Profile Extension Startup

- Use VS Code's built-in profiler
- Check for slow activation events
- Optimize heavy operations

### 2. Memory Usage

- Monitor extension host memory
- Check for memory leaks
- Use Chrome DevTools for heap snapshots

### 3. TreeView Performance

- Test with large datasets
- Check rendering performance
- Optimize data providers

## Reporting Issues

When reporting issues, please include:

1. VS Code version
2. Extension version
3. MCP server version
4. Operating system
5. Error messages from Output panel
6. Steps to reproduce

## Additional Resources

- [VS Code Extension API](https://code.visualstudio.com/api)
- [VS Code Extension Testing](https://code.visualstudio.com/api/working-with-extensions/testing-extension)
- [MCP Protocol Documentation](https://modelcontextprotocol.io/docs)
