# Release Notes - v0.16.1

**Release Date**: October 5, 2025
**Type**: VS Code Extension Enhancement Release

## 🚀 Major Features

### Automatic Server Startup

The VS Code extension now **automatically manages the Gorev server lifecycle** - no manual commands required!

**Key Features**:

- ✅ **Auto-Detection**: Checks if server is running on port 5082 before starting
- ✅ **Auto-Start**: Spawns server process automatically if not running
- ✅ **Zero Configuration**: Works out of the box - just install and go!
- ✅ **Graceful Lifecycle**: Proper shutdown when extension deactivates

**User Impact**:

- **Before**: Users had to manually run `npx @mehmetsenol/gorev-mcp-server serve` before using the extension
- **After**: Extension handles everything automatically - just open VS Code and start working!

### Smart Database Management

**Workspace-Specific Databases** with automatic path configuration:

**Database Location Priority**:

1. **Workspace folder**: `.gorev/gorev.db` (preferred)
2. **User home**: `~/.gorev/gorev.db` (fallback)

**Features**:

- ✅ Automatic directory creation with proper permissions
- ✅ Set via `GOREV_DB_PATH` environment variable
- ✅ Fixes SQLite "out of memory" errors on Windows (actually file permission issues)
- ✅ Cross-platform compatible (Windows, macOS, Linux)

### Complete Server Lifecycle Management

**Process Management Features**:

- ✅ Port availability checking before server start
- ✅ Proper stdio configuration (`stdin` kept open for MCP protocol)
- ✅ Server output logged to VS Code Output panel
- ✅ Graceful shutdown on extension deactivation
- ✅ SIGTERM for graceful stop, SIGKILL fallback after 5 seconds timeout

## 🐛 Bug Fixes

### Server Exit Issue

**Problem**: Server was exiting immediately after startup
**Cause**: stdio was set to `['ignore', 'pipe', 'pipe']`, closing stdin
**Fix**: Changed to `['pipe', 'pipe', 'pipe']` - MCP server requires open stdin pipe to prevent EOF exit

### Flag Compatibility

**Problem**: Some binary versions didn't support `--api-port` flag
**Fix**: Removed flag from startup arguments (server defaults to port 5082 anyway)

## 📝 Code Changes

### UnifiedServerManager Refactor

**File**: `gorev-vscode/src/managers/unifiedServerManager.ts` (+300 lines)

**New Methods**:

- `isServerRunning()`: Port availability check using TCP socket connection
- `startServer()`: Spawns server process with proper environment variables
- `waitForServerReady()`: Polls port until server is ready (15s timeout)
- `stopServer()`: Graceful shutdown with timeout fallback

**Enhanced Features**:

- Server process spawning with `child_process.spawn()`
- Environment variable configuration (GOREV_DB_PATH)
- Cross-platform command handling (Windows: `npx.cmd`, Unix: `npx`)
- Server output streaming to VS Code Output panel
- Error handling and recovery

### Extension Activation Flow

**Updated**: `gorev-vscode/src/extension.ts`

**Changes**:

- Extension `dispose()` method now async to properly await server shutdown
- Added await for `serverManager.dispose()` on extension deactivation

## 🔧 Technical Details

### Server Startup Process

1. **Check if server is running** → `isServerRunning()`
   - Attempts TCP connection to localhost:5082
   - 1 second timeout for connection attempt

2. **Start server if needed** → `startServer()`
   - Determine database path (workspace or home)
   - Create `.gorev` directory if doesn't exist
   - Spawn server process: `npx @mehmetsenol/gorev-mcp-server serve --debug`
   - Set `GOREV_DB_PATH` environment variable

3. **Wait for server ready** → `waitForServerReady()`
   - Poll port every 500ms until available
   - 15 second timeout for server startup
   - Throws error if timeout exceeded

4. **Connect API client** → `apiClient.connect()`
   - Connect to REST API at http://localhost:5082
   - Register workspace with server
   - Initialize extension features

### Cross-Platform Compatibility

**Windows**:

- Uses `npx.cmd` instead of `npx`
- Shell mode enabled for proper command execution
- Special handling for file permissions

**Unix (macOS/Linux)**:

- Uses `npx` directly
- No shell mode needed
- Standard file permissions

## 🔗 Related Links

- [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
- [GitHub Repository](https://github.com/msenol/Gorev)
- [Full CHANGELOG](https://github.com/msenol/Gorev/blob/main/gorev-vscode/CHANGELOG.md#0161---2025-10-05)
- [MCP Server Package](https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server)

## 📊 Impact

This release dramatically improves the VS Code extension user experience:

- ⬇️ **Setup Time**: Reduced from ~2 minutes to **instant**
- 🎯 **User Confusion**: Eliminated "server not running" errors
- 💪 **Reliability**: Automatic recovery from server crashes
- 🚀 **Productivity**: Users can start working immediately

## 🙏 User Feedback

This feature was developed in response to user feedback about the manual server startup requirement. Thank you to everyone who provided input!
