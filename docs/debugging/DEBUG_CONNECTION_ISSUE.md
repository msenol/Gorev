# VS Code MCP Connection Issue - Debug Report & Solution

## Problem Summary

The VS Code extension cannot reliably connect to the MCP server. Based on the logs, the connection initially works (we see successful API calls) but then starts timing out with "Request timeout" errors.

## Root Cause

The MCP server is using `log.Printf()` statements that output to stderr, which interferes with the MCP protocol communication. The MCP protocol expects clean JSON-RPC messages on stdin/stdout, but these log messages are contaminating the communication channel.

### Evidence from logs

1. `/home/msenol/.vscode-server/data/logs/20250709T154654/exthost13/output_logging_20250709T164742/1-Gorev.log` shows:
   - Initial successful API calls (gorev_listele, proje_listele, template_listele)
   - Then timeout errors: "Request timeout for tools/call"
   - Server error: "sunucu başlatılamadı: context canceled"

2. When running the server manually:

   ```
   2025/07/09 21:55:03 Migration öncesi veritabanı versiyonu: 6, dirty: false
   2025/07/09 21:55:03 Migration sonrası veritabanı versiyonu: 6, dirty: false
   2025/07/09 21:55:03 Veritabanı başarıyla migrate edildi.
   Gorev MCP sunucusu başlatılıyor...
   ```

   These log messages are sent to stderr/stdout and interfere with JSON-RPC protocol.

## Solution

### Option 1: Quick Fix - Disable Logging (Recommended for immediate fix)

1. **Create a patch to disable logging in veri_yonetici.go:**

```bash
# Create a backup
cp gorev-mcpserver/internal/gorev/veri_yonetici.go gorev-mcpserver/internal/gorev/veri_yonetici.go.backup

# Comment out all log statements
sed -i 's/log\.Printf/\/\/log.Printf/g' gorev-mcpserver/internal/gorev/veri_yonetici.go
sed -i 's/log\.Println/\/\/log.Println/g' gorev-mcpserver/internal/gorev/veri_yonetici.go
```

2. **Remove the startup message in main.go:**
Edit `/mnt/f/Development/Projects/Gorev/gorev-mcpserver/cmd/gorev/main.go` and comment out line 233:

```go
// fmt.Fprintln(os.Stderr, "Gorev MCP sunucusu başlatılıyor...")
```

3. **Rebuild the server:**

```bash
cd gorev-mcpserver
make build
```

### Option 2: Proper Fix - Use File-based Logging

1. **Create a logger that writes to a file instead of stderr:**

Create `/mnt/f/Development/Projects/Gorev/gorev-mcpserver/internal/gorev/logger.go`:

```go
package gorev

import (
    "log"
    "os"
    "path/filepath"
)

var fileLogger *log.Logger

func init() {
    // Create logs directory
    logDir := filepath.Join(os.TempDir(), "gorev-logs")
    os.MkdirAll(logDir, 0755)
    
    // Open log file
    logFile, err := os.OpenFile(
        filepath.Join(logDir, "gorev-mcp.log"),
        os.O_CREATE|os.O_WRONLY|os.O_APPEND,
        0666,
    )
    if err != nil {
        // If we can't create a log file, create a no-op logger
        fileLogger = log.New(os.NewFile(0, os.DevNull), "", 0)
        return
    }
    
    fileLogger = log.New(logFile, "", log.LstdFlags)
}

func LogPrintf(format string, v ...interface{}) {
    if fileLogger != nil {
        fileLogger.Printf(format, v...)
    }
}

func LogPrintln(v ...interface{}) {
    if fileLogger != nil {
        fileLogger.Println(v...)
    }
}
```

2. **Replace all log calls in veri_yonetici.go:**

```bash
# In veri_yonetici.go, replace:
# log.Printf -> LogPrintf
# log.Println -> LogPrintln
```

### Option 3: Environment Variable Control

Add an environment variable check to disable logging when running as MCP server:

```go
// At the top of veri_yonetici.go
var debugMode = os.Getenv("GOREV_DEBUG") == "true"

// Wrap all log statements:
if debugMode {
    log.Printf("Migration öncesi versiyon alınamadı: %v", err)
}
```

## Testing the Fix

After applying the fix:

1. **Rebuild the server:**

```bash
cd gorev-mcpserver
make build
```

2. **Test manually:**

```bash
# Should produce NO output to stderr except for JSON-RPC messages
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | ./build/gorev serve
```

3. **Test in VS Code:**

- Restart VS Code
- Check the Gorev output panel
- Verify no timeout errors occur

## Additional Recommendations

1. **Update VS Code extension timeout:**
   In `/mnt/f/Development/Projects/Gorev/gorev-vscode/src/mcp/client.ts`, consider increasing the timeout from 10 seconds to 30 seconds:

   ```typescript
   }, 30000); // 30 second timeout
   ```

2. **Add connection retry logic:**
   The extension should retry connection attempts if the initial connection fails.

3. **Better error handling:**
   The MCP server should validate that it's running in a proper MCP environment before starting the stdio server.

## Verification Steps

1. Check that the server binary exists and is executable:

   ```bash
   ls -la /mnt/f/Development/Projects/Gorev/gorev-mcpserver/build/gorev
   ```

2. Verify the server starts without output:

   ```bash
   ./build/gorev serve 2>/dev/null
   # Should wait for input without printing anything
   ```

3. Monitor VS Code logs:

   ```bash
   tail -f /home/msenol/.vscode-server/data/logs/*/exthost*/output_logging_*/1-Gorev.log
   ```

## Log File Locations

- VS Code Extension Logs: `/home/msenol/.vscode-server/data/logs/*/exthost*/output_logging_*/1-Gorev.log`
- MCP Server Logs (after fix): `/tmp/gorev-logs/gorev-mcp.log`

## Current Configuration

- Server Path: `/mnt/f/Development/Projects/Gorev/gorev-mcpserver/build/gorev`
- Auto-connect: Enabled
- Refresh Interval: 30 seconds
- Request Timeout: 10 seconds (should be increased)
