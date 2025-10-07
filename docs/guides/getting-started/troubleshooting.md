# Troubleshooting Guide

**Version**: v0.16.0
**Last Updated**: October 5, 2025

Common issues and solutions for Gorev MCP Server, Web UI, and VS Code Extension.

---

## Quick Diagnostics

Run these commands first to gather information:

```bash
# Version check
npx @mehmetsenol/gorev-mcp-server --version

# Database verification
ls -la .gorev/gorev.db

# Server health
curl http://localhost:5082/api/health

# Check running processes
ps aux | grep gorev
```

---

## Installation Issues

### NPM: "ENOENT: spawn npx"

**Platform**: Windows

**Symptoms**:

```
Error: spawn npx ENOENT
```

**Cause**: Node.js not installed or not in PATH

**Solution**:

```cmd
REM 1. Install Node.js
REM Download from: https://nodejs.org/

REM 2. Verify installation
node --version
npm --version
npx --version

REM 3. Restart terminal/system

REM 4. Test again
npx @mehmetsenol/gorev-mcp-server --version
```

### NPM: Package Not Found

**Symptoms**:

```
npm ERR! 404 '@mehmetsenol/gorev-mcp-server' is not in the npm registry
```

**Solutions**:

```bash
# 1. Check package name (common typo)
npm view @mehmetsenol/gorev-mcp-server

# 2. Clear NPM cache
npm cache clean --force

# 3. Use full package name with scope
npx -y @mehmetsenol/gorev-mcp-server serve

# 4. Check npm registry
npm config get registry  # Should be https://registry.npmjs.org/
```

### Binary Installation Fails

**Linux/macOS**:

```bash
# Permission denied
sudo chmod +x gorev-linux-amd64
./gorev-linux-amd64 --version

# Move to PATH
sudo mv gorev-linux-amd64 /usr/local/bin/gorev
```

**Windows**:

```powershell
# Run as Administrator
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Download manually
Invoke-WebRequest -Uri "https://github.com/msenol/gorev/releases/latest/download/gorev-windows-amd64.exe" -OutFile "gorev.exe"
```

---

## Server Issues

### Server Won't Start

**Symptoms**:

```
Error: listen EADDRINUSE: address already in use :::5082
```

**Cause**: Port 5082 already in use

**Solutions**:

```bash
# 1. Find process using port
lsof -i :5082                    # Linux/macOS
netstat -ano | findstr :5082     # Windows

# 2. Kill process
kill -9 <PID>                    # Linux/macOS
taskkill /PID <PID> /F           # Windows

# 3. Or use different port
npx @mehmetsenol/gorev-mcp-server serve --api-port 8080
```

### Server Crashes on Startup

**Symptoms**:

```
panic: runtime error: invalid memory address
```

**Solutions**:

```bash
# 1. Check database integrity
sqlite3 .gorev/gorev.db "PRAGMA integrity_check;"

# 2. Backup and reinitialize
mv .gorev/gorev.db .gorev/gorev.db.backup
npx @mehmetsenol/gorev-mcp-server init

# 3. Restore data
npx @mehmetsenol/gorev-mcp-server import --input backup.json

# 4. Check debug logs
npx @mehmetsenol/gorev-mcp-server serve --debug 2>&1 | tee gorev.log
```

### MCP Server Not Responding

**Symptoms**:

- Claude/Cursor shows "Server not connected"
- MCP tools timeout

**Solutions**:

```bash
# 1. Verify server is running
ps aux | grep gorev

# 2. Check stdio connection
# Server must be started by MCP client, not manually

# 3. Restart MCP client (Claude/Cursor)

# 4. Check MCP config syntax
cat ~/.config/Claude/claude_desktop_config.json | jq .

# 5. Test MCP server manually
echo '{"method":"tools/list"}' | npx @mehmetsenol/gorev-mcp-server serve
```

---

## Database Issues

### Database Locked

**Symptoms**:

```
Error: database is locked (code 5)
```

**Cause**: Multiple processes accessing database, or WAL mode on network drive

**Solutions**:

```bash
# 1. Close all Gorev instances
pkill gorev

# 2. Remove lock files
rm -f .gorev/gorev.db-shm .gorev/gorev.db-wal

# 3. Disable WAL mode (if on network drive)
sqlite3 .gorev/gorev.db "PRAGMA journal_mode=DELETE;"

# 4. Move workspace to local directory
# VirtualBox shared folders don't support WAL mode
mv /media/shared/workspace /tmp/workspace
cd /tmp/workspace
```

### Corrupted Database

**Symptoms**:

```
Error: database disk image is malformed
```

**Solutions**:

```bash
# 1. Backup immediately
cp .gorev/gorev.db .gorev/gorev.db.corrupt

# 2. Try to recover
sqlite3 .gorev/gorev.db ".recover" | sqlite3 .gorev/gorev-recovered.db

# 3. Verify recovered database
sqlite3 .gorev/gorev-recovered.db "SELECT COUNT(*) FROM gorevler;"

# 4. Replace if successful
mv .gorev/gorev.db .gorev/gorev.db.backup
mv .gorev/gorev-recovered.db .gorev/gorev.db

# 5. Last resort: export/import
npx @mehmetsenol/gorev-mcp-server export --output backup.json
npx @mehmetsenol/gorev-mcp-server init --force
npx @mehmetsenol/gorev-mcp-server import --input backup.json
```

### Database Migration Failed

**Symptoms**:

```
Error: migration failed at version 003
```

**Solutions**:

```bash
# 1. Check current schema version
sqlite3 .gorev/gorev.db "SELECT version FROM schema_version;"

# 2. Backup before retry
cp .gorev/gorev.db .gorev/gorev.db.pre-migration

# 3. Manual migration (if needed)
sqlite3 .gorev/gorev.db < gorev-mcpserver/internal/veri/migrations/003_add_workspace_support.sql

# 4. Update version
sqlite3 .gorev/gorev.db "UPDATE schema_version SET version = 3;"
```

---

## Web UI Issues

### Web UI Not Loading

**Symptoms**:

- Browser shows 404 or blank page
- "Cannot GET /" error

**Solutions**:

```bash
# 1. Verify server is running
curl http://localhost:5082/api/health

# 2. Check embedded files
npx @mehmetsenol/gorev-mcp-server serve --debug
# Look for: "📦 Web UI embedded in binary"

# 3. Clear browser cache
# Chrome: Ctrl+Shift+Delete → Clear cached images and files

# 4. Try different browser

# 5. Check binary version
npx @mehmetsenol/gorev-mcp-server --version
# Must be v0.16.0 or later for embedded UI
```

### API Connection Errors

**Symptoms**:

```
Network Error: Failed to fetch
ERR_CONNECTION_REFUSED
```

**Solutions**:

```bash
# 1. Verify API server running
curl http://localhost:5082/api/workspaces

# 2. Check CORS (if using dev server)
# Browser Console → Network tab → Check response headers

# 3. Verify workspace headers
# Browser DevTools → Network → Request Headers:
# X-Workspace-Id: ...
# X-Workspace-Path: ...

# 4. Restart server
pkill gorev
npx @mehmetsenol/gorev-mcp-server serve
```

### Workspace Not Loading

**Symptoms**:

- Empty workspace list
- "No workspace selected" error

**Solutions**:

```bash
# 1. Initialize workspace
cd /path/to/project
npx @mehmetsenol/gorev-mcp-server init

# 2. Verify database exists
ls -la .gorev/gorev.db

# 3. Check workspace registry
sqlite3 ~/.gorev/workspace_registry.db "SELECT * FROM workspaces;"

# 4. Manual registration
npx @mehmetsenol/gorev-mcp-server workspace register \
  --path $(pwd) \
  --name "My Project"

# 5. Clear localStorage
# Browser DevTools → Application → Local Storage → Clear All
```

### Language Not Switching

**Symptoms**:

- UI stays in Turkish despite setting English

**Solutions**:

```bash
# 1. Restart server with language flag
GOREV_LANG=en npx @mehmetsenol/gorev-mcp-server serve

# 2. Clear localStorage
# Browser DevTools → Application → Local Storage → Delete "gorev_language"

# 3. Set environment variable permanently
echo 'export GOREV_LANG=en' >> ~/.bashrc
source ~/.bashrc

# 4. Check MCP config
cat ~/.config/Claude/claude_desktop_config.json
# Should have: "env": { "GOREV_LANG": "en" }
```

---

## VS Code Extension Issues

### Extension Not Detecting Workspace

**Symptoms**:

- Task TreeView empty
- "No workspace found" in status bar

**Solutions**:

```bash
# 1. Initialize workspace
cd /path/to/vscode/workspace
npx @mehmetsenol/gorev-mcp-server init

# 2. Reload VS Code window
# Ctrl+Shift+P → "Reload Window"

# 3. Check workspace folders
# VS Code → File → Add Folder to Workspace

# 4. Verify .gorev/ directory in workspace root
ls -la .gorev/

# 5. Check extension logs
# VS Code → View → Output → Select "Gorev Extension"
```

### MCP Server Connection Failed

**Symptoms**:

- Extension shows "Server: Disconnected"
- MCP tools don't work

**Solutions**:

```bash
# 1. Check MCP config exists
ls -la .kilocode/mcp.json
# or
ls -la .cursor/mcp.json

# 2. Verify config syntax
cat .kilocode/mcp.json | jq .

# 3. Test MCP server manually
npx @mehmetsenol/gorev-mcp-server serve

# 4. Check VS Code extension compatibility
# Kilo Code extension required for MCP support

# 5. Restart VS Code
```

### Extension Commands Not Working

**Symptoms**:

- "Command 'gorev.createTask' not found"

**Solutions**:

```bash
# 1. Verify extension installed
code --list-extensions | grep gorev

# 2. Reinstall extension
code --uninstall-extension mehmetsenol.gorev-vscode
code --install-extension mehmetsenol.gorev-vscode

# 3. Check extension enabled
# VS Code → Extensions → Search "gorev" → Ensure enabled

# 4. Reload window
# Ctrl+Shift+P → "Reload Window"
```

---

## MCP Configuration Issues

### Claude Desktop: Server Not Connecting

**Symptoms**:

- Claude shows "gorev (not connected)"

**Solutions**:

```bash
# 1. Verify config file location
# Linux:   ~/.config/Claude/claude_desktop_config.json
# macOS:   ~/Library/Application Support/Claude/claude_desktop_config.json
# Windows: %APPDATA%\Claude\claude_desktop_config.json

# 2. Check config syntax
cat ~/.config/Claude/claude_desktop_config.json | jq .

# 3. Verify NPM package name
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest"
      ]
    }
  }
}

# 4. Restart Claude Desktop
pkill Claude && open -a Claude

# 5. Check Claude logs
# macOS: ~/Library/Logs/Claude/
# Look for MCP connection errors
```

### Cursor: MCP Tools Not Available

**Symptoms**:

- Cursor doesn't show Gorev tools in autocomplete

**Solutions**:

```bash
# 1. Create .cursor/mcp.json in workspace root
mkdir -p .cursor
cat > .cursor/mcp.json << 'EOF'
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest"
      ],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
EOF

# 2. Reload Cursor window

# 3. Test connection
# Cursor → Chat → Type: "List my tasks"

# 4. Check Cursor logs
# Cursor → View → Output → Select "MCP"
```

---

## Template Issues

### Template Not Found

**Symptoms**:

```
Error: template not found: bug
```

**Solutions**:

```bash
# 1. List available templates
npx @mehmetsenol/gorev-mcp-server template list

# 2. Verify alias exists
npx @mehmetsenol/gorev-mcp-server template list | grep "bug"

# 3. Use UUID if alias doesn't work
npx @mehmetsenol/gorev-mcp-server task create \
  --template 39f28dbd-10f3-454c-8b35-52ae6b7ea391

# 4. Reinitialize templates (if corrupted)
sqlite3 .gorev/gorev.db < gorev-mcpserver/internal/veri/migrations/002_default_templates.sql
```

### Required Field Missing

**Symptoms**:

```
Error: Required field 'baslik' is missing
```

**Solutions**:

```bash
# 1. Check template schema
npx @mehmetsenol/gorev-mcp-server template show --alias bug

# 2. Provide all required fields
npx @mehmetsenol/gorev-mcp-server task create --template bug \
  --field baslik="Title" \
  --field modul="Module" \
  --field ortam="production" \
  --field adimlar="Steps" \
  --field beklenen="Expected" \
  --field mevcut="Actual"

# 3. Use Web UI for guided form
# http://localhost:5082 → New Task → Select template
```

---

## Performance Issues

### Slow Task Listing

**Symptoms**:

- Task list takes > 5 seconds to load

**Solutions**:

```bash
# 1. Check database size
du -h .gorev/gorev.db

# 2. Vacuum database
sqlite3 .gorev/gorev.db "VACUUM;"

# 3. Rebuild FTS index
sqlite3 .gorev/gorev.db "DELETE FROM gorevler_fts;"
sqlite3 .gorev/gorev.db "INSERT INTO gorevler_fts SELECT id, baslik, aciklama FROM gorevler;"

# 4. Archive completed tasks
npx @mehmetsenol/gorev-mcp-server export \
  --filter status=completed \
  --output completed-archive.json
# Then delete from database

# 5. Optimize indexes
sqlite3 .gorev/gorev.db "ANALYZE;"
```

### High Memory Usage

**Symptoms**:

- Gorev process uses > 500MB RAM

**Solutions**:

```bash
# 1. Check for memory leaks
# Monitor over time
watch -n 1 'ps aux | grep gorev'

# 2. Restart server periodically
# Add to crontab for daily restart

# 3. Reduce cache size
# Edit code (future config option)

# 4. Report issue
# https://github.com/msenol/gorev/issues
```

---

## Export/Import Issues

### Export Fails with Timeout

**Symptoms**:

```
Error: export timed out after 30s
```

**Solutions**:

```bash
# 1. Export in smaller chunks
npx @mehmetsenol/gorev-mcp-server export \
  --filter "created_after=2025-01-01" \
  --output recent-tasks.json

# 2. Exclude completed tasks
npx @mehmetsenol/gorev-mcp-server export \
  --exclude-completed \
  --output active-tasks.json

# 3. Increase timeout (future option)

# 4. Direct database dump
sqlite3 .gorev/gorev.db ".dump" > database-dump.sql
```

### Import Conflicts

**Symptoms**:

```
Warning: 15 tasks skipped due to conflicts
```

**Solutions**:

```bash
# 1. Use overwrite mode
npx @mehmetsenol/gorev-mcp-server import \
  --input backup.json \
  --conflict-resolution overwrite

# 2. Preview conflicts with dry run
npx @mehmetsenol/gorev-mcp-server import \
  --input backup.json \
  --dry-run

# 3. Merge mode
npx @mehmetsenol/gorev-mcp-server import \
  --input backup.json \
  --import-mode merge

# 4. Manual conflict resolution
# Edit backup.json, change task IDs, import again
```

---

## Diagnostic Tools

### Enable Debug Mode

```bash
# Server debug logs
npx @mehmetsenol/gorev-mcp-server serve --debug 2>&1 | tee gorev.log

# Environment variables
export GOREV_DEBUG=1
export GOREV_LOG_LEVEL=trace
```

### Database Inspection

```bash
# Schema version
sqlite3 .gorev/gorev.db "SELECT * FROM schema_version;"

# Table counts
sqlite3 .gorev/gorev.db "
SELECT 'gorevler' as table_name, COUNT(*) FROM gorevler
UNION ALL
SELECT 'projeler', COUNT(*) FROM projeler
UNION ALL
SELECT 'baglantilar', COUNT(*) FROM baglantilar;
"

# Recent tasks
sqlite3 .gorev/gorev.db "
SELECT baslik, durum, olusturulma_tarihi
FROM gorevler
ORDER BY olusturulma_tarihi DESC
LIMIT 5;
"
```

### Health Check Script

```bash
#!/bin/bash
# health-check.sh

echo "=== Gorev Health Check ==="

# Version
echo "Version:"
npx @mehmetsenol/gorev-mcp-server --version

# Database
echo "Database:"
ls -lh .gorev/gorev.db

# Integrity
echo "Database Integrity:"
sqlite3 .gorev/gorev.db "PRAGMA integrity_check;"

# Server health
echo "Server Health:"
curl -s http://localhost:5082/api/health | jq .

# Workspace count
echo "Workspaces:"
sqlite3 ~/.gorev/workspace_registry.db "SELECT COUNT(*) FROM workspaces;"

echo "=== Health Check Complete ==="
```

---

## Getting Help

### Before Opening an Issue

Collect this information:

```bash
# System info
uname -a
node --version
npm --version

# Gorev version
npx @mehmetsenol/gorev-mcp-server --version

# Database info
ls -lh .gorev/gorev.db
sqlite3 .gorev/gorev.db "SELECT version FROM schema_version;"

# Error logs
npx @mehmetsenol/gorev-mcp-server serve --debug 2>&1 | head -50
```

### Useful Log Locations

- **Server Logs**: stdout/stderr (use `tee` to save)
- **VS Code Extension**: Output panel → "Gorev Extension"
- **Claude Desktop**: `~/Library/Logs/Claude/` (macOS)
- **Web UI**: Browser DevTools → Console
- **Database**: `.gorev/gorev.db` (use sqlite3 to inspect)

### Reporting Bugs

**GitHub Issues**: https://github.com/msenol/gorev/issues

Include:

- Gorev version
- Operating system
- Error messages (full stack trace)
- Steps to reproduce
- Expected vs. actual behavior
- Relevant logs
- Database schema version

---

## Common Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| `EADDRINUSE` | Port in use | Kill process or use different port |
| `ENOENT` | File not found | Check paths and file existence |
| `database is locked` | Concurrent access | Close other instances, disable WAL |
| `template not found` | Invalid alias/UUID | List templates, verify spelling |
| `required field missing` | Incomplete data | Provide all required template fields |
| `workspace not found` | Not initialized | Run `gorev init` |
| `migration failed` | Schema mismatch | Backup DB, retry migration |
| `connection refused` | Server not running | Start server |

---

**Still Having Issues?**

1. Search existing issues: https://github.com/msenol/gorev/issues
2. Ask in discussions: https://github.com/msenol/gorev/discussions
3. Open new issue with details above

We're here to help! 🚀
