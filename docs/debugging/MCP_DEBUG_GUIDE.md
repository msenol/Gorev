# Gorev MCP Connection Debugging Guide

## Overview

This guide helps diagnose and fix connection issues between the VS Code extension and the MCP server.

## Solution

Based on the diagnosis in `DEBUG_CONNECTION_ISSUE.md`, the problem was that the MCP server was outputting log messages to stderr which interfered with the JSON-RPC protocol. This has been resolved by removing the problematic stderr logging. To apply the fix, rebuild the server:

```bash
cd gorev-mcpserver
make build
```

## Debug Tools Created

### 1. Debug Wrapper Script (`debug-wrapper.sh`)

A shell script that intercepts all communication between VS Code and the MCP server, logging to:

- `/tmp/gorev-debug/mcp-session-*.log` - Main debug log
- `/tmp/gorev-debug/stdin-*.log` - Input messages (VS Code → Server)
- `/tmp/gorev-debug/stdout-*.log` - Output messages (Server → VS Code)

### 2. Test Scripts

#### `test-mcp-server.sh`

Basic shell script to test MCP server functionality:

```bash
cd gorev-mcpserver
./test-mcp-server.sh
```

#### `test-mcp-protocol.py`

Python script for more thorough MCP protocol testing:

```bash
cd gorev-mcpserver
python3 test-mcp-protocol.py
```

### 3. Debug Server (`cmd/gorev-debug/`)

A special debug version of the server that logs all activity to a file instead of stderr:

```bash
cd gorev-mcpserver
./build-debug.sh
./build/gorev-debug serve
# Check logs in /tmp/gorev-debug/
```

### 4. VS Code Debug Commands

New commands added to the extension:

- **Toggle Debug Mode** - Enable/disable debug wrapper
- **Show Debug Logs** - View debug log files
- **Clear Debug Logs** - Remove all debug logs
- **Test MCP Connection** - Test server connectivity

Access these via Command Palette (Ctrl+Shift+P) → "Gorev Debug:"

## Using Debug Mode in VS Code

1. **Enable Debug Mode:**
   - Open VS Code settings (Ctrl+,)
   - Search for "gorev.debug"
   - Enable "Use Debug Wrapper"
   - Restart VS Code

2. **View Debug Logs:**
   - Command Palette → "Gorev Debug: Show Debug Logs"
   - Select a log file to view

3. **Test Connection:**
   - Command Palette → "Gorev Debug: Test MCP Connection"
   - Check the Output panel (View → Output → Gorev)

## Configuration Options

New settings in VS Code:

- `gorev.debug.useWrapper` - Enable debug wrapper (default: false)
- `gorev.debug.logPath` - Debug log directory (default: /tmp/gorev-debug)
- `gorev.debug.serverTimeout` - Server timeout in ms (default: 30000)

## Troubleshooting Steps

1. **Check Server Path:**

   ```bash
   ls -la /mnt/f/Development/Projects/Gorev/gorev-mcpserver/build/gorev
   ```

2. **Test Server Manually:**

   ```bash
   cd gorev-mcpserver
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | ./build/gorev serve
   ```

   Should return a JSON response, not log messages.

3. **Check VS Code Logs:**
   - View → Output → Gorev
   - Help → Toggle Developer Tools → Console

4. **Enable Debug Mode:**
   - Follow steps above to enable debug wrapper
   - Run a task operation
   - Check debug logs in `/tmp/gorev-debug/`

## Common Issues

### "Request timeout" errors

- Increase timeout: `gorev.debug.serverTimeout` to 60000
- Check if server is outputting non-JSON to stdout/stderr

### "Server not found" errors

- Verify server path in settings
- Ensure server is executable: `chmod +x build/gorev`
- Check file exists and has correct permissions

### Connection works then fails

- This is the issue addressed in this fix - log output interfering with protocol
- Rebuild server with `make build` to use commented-out logs

## Development Tips

1. **Always use file-based logging** in the MCP server, never log to stderr/stdout
2. **Test changes** with the Python test script before VS Code integration
3. **Use debug wrapper** during development to capture protocol issues
4. **Monitor both sides** - VS Code output and server debug logs

## Log File Locations

- **VS Code Extension Logs:** Check Output panel (View → Output → Gorev)
- **Debug Wrapper Logs:** `/tmp/gorev-debug/mcp-session-*.log`
- **Debug Server Logs:** `/tmp/gorev-debug/gorev-debug-*.log`
- **VS Code System Logs:** Help → Toggle Developer Tools → Console

## Next Steps

1. Rebuild the server: `cd gorev-mcpserver && make build`
2. Test with: `./test-mcp-protocol.py`
3. Restart VS Code and check if connection works
4. If issues persist, enable debug mode and check logs
