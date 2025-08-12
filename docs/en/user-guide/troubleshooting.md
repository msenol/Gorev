# Troubleshooting Guide

Comprehensive guide for resolving common issues with Gorev installation, configuration, and usage across different platforms and AI assistants.

## 🚨 Quick Solutions

### Most Common Issues

#### ❌ "Command not found: gorev"
**Problem**: Gorev binary not in PATH or not installed
**Quick Fix**:
```bash
# Check if gorev is installed
which gorev

# If not found, reinstall
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash

# Windows PowerShell
iwr -useb https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

#### ❌ "MCP server not running" (VS Code Extension)
**Problem**: Extension can't connect to Gorev server
**Quick Fix**:
1. Start server manually: `gorev serve`
2. Check VS Code settings: `gorev.serverPath`
3. Restart VS Code

#### ❌ "Template usage is now mandatory"
**Problem**: Trying to use deprecated `gorev_olustur` command
**Quick Fix**: Use templates instead:
```
# Old (doesn't work)
Create task with title "Fix bug"

# New (works)
Create bug report task:
- Title: Login form validation error
- Module: Authentication
- Priority: high
```

## 📋 Installation Issues

### Windows Installation Problems

#### PowerShell Execution Policy Error
**Symptoms**: "Execution of scripts is disabled on this system"
**Solution**:
```powershell
# Check current policy
Get-ExecutionPolicy

# If Restricted, temporarily allow scripts
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Run installer
iwr -useb https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex

# Restore policy (optional)
Set-ExecutionPolicy -ExecutionPolicy Restricted -Scope CurrentUser
```

#### Windows Defender False Positive
**Symptoms**: Installer blocked or gorev.exe quarantined
**Solution**:
1. Open Windows Security → Virus & threat protection
2. Go to "Manage settings" under "Virus & threat protection settings"
3. Add exclusion for the installation directory
4. Re-run installer

#### PATH Not Updated
**Symptoms**: `gorev` command not found after installation
**Solution**:
```cmd
# Check if directory is in PATH
echo %PATH%

# If not, add manually (PowerShell as Admin)
$env:PATH += ";$env:USERPROFILE\AppData\Local\Programs\gorev"

# Or restart terminal/VS Code to pick up changes
```

### macOS Installation Problems

#### Gatekeeper Blocking Execution
**Symptoms**: "gorev cannot be opened because it is from an unidentified developer"
**Solution**:
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/gorev

# Or allow in System Preferences
# System Preferences → Security & Privacy → General
# Click "Allow Anyway" next to the blocked app message
```

#### Homebrew Permission Issues
**Symptoms**: Permission denied when installing via Homebrew
**Solution**:
```bash
# Fix Homebrew permissions
sudo chown -R $(whoami) $(brew --prefix)/*

# Then retry installation
brew install msenol/tap/gorev
```

#### Wrong Architecture
**Symptoms**: "Bad CPU type in executable" on Apple Silicon
**Solution**:
```bash
# Check if you're running x86 binary on M1/M2
file /usr/local/bin/gorev

# Download correct ARM64 version
curl -sSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash
```

### Linux Installation Problems

#### Permission Denied
**Symptoms**: `permission denied: ./gorev`
**Solution**:
```bash
# Make binary executable
chmod +x /usr/local/bin/gorev

# Or if local installation
chmod +x ~/bin/gorev
```

#### Missing Dependencies
**Symptoms**: `error while loading shared libraries`
**Solution**:
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install libc6

# CentOS/RHEL/Fedora
sudo yum install glibc

# Alpine Linux
apk add libc6-compat
```

#### AppImage Issues
**Symptoms**: AppImage doesn't start or crashes
**Solution**:
```bash
# Install FUSE if missing
sudo apt install fuse

# Make AppImage executable
chmod +x gorev-*.AppImage

# Run directly
./gorev-*.AppImage serve
```

## 🔧 Configuration Issues

### VS Code Extension Problems

#### Extension Not Loading
**Symptoms**: Gorev icon doesn't appear in Activity Bar
**Solutions**:
1. **Check VS Code version**: Requires VS Code 1.74.0+
   ```bash
   code --version
   ```

2. **Verify installation**:
   - Open Extensions (`Ctrl+Shift+X`)
   - Search for "Gorev"
   - Ensure it's installed and enabled

3. **Check for conflicts**:
   - Disable other task management extensions temporarily
   - Restart VS Code in safe mode: `code --disable-extensions`

4. **Clear extension cache**:
   ```bash
   # Windows
   rmdir /s "%USERPROFILE%\.vscode\extensions\.obsolete"
   
   # macOS/Linux
   rm -rf ~/.vscode/extensions/.obsolete
   ```

#### Connection Timeout
**Symptoms**: "Failed to connect to MCP server" after 30 seconds
**Solutions**:
1. **Check server path in settings**:
   ```json
   {
     "gorev.serverPath": "/usr/local/bin/gorev",  // Use full path
     "gorev.timeout": 60000                       // Increase timeout
   }
   ```

2. **Start server manually first**:
   ```bash
   gorev serve --debug
   # Keep running, then restart VS Code
   ```

3. **Check firewall/antivirus**: May be blocking local connections

#### TreeView Empty or Not Refreshing
**Symptoms**: Extension loads but shows no tasks/projects
**Solutions**:
1. **Force refresh**: Click refresh button or `Ctrl+R`

2. **Check data directory**:
   ```bash
   ls -la ~/.gorev/data/
   # Should contain gorev.db
   ```

3. **Reset database** (if corrupted):
   ```bash
   mv ~/.gorev/data/gorev.db ~/.gorev/data/gorev.db.backup
   gorev serve  # Creates new database
   ```

4. **Check extension output**:
   - View → Output
   - Select "Gorev" from dropdown
   - Look for error messages

### MCP Server Configuration

#### Database Issues

**Database Locked Error**:
```bash
# Stop all gorev processes
pkill gorev

# Remove lock file
rm ~/.gorev/data/gorev.db-wal
rm ~/.gorev/data/gorev.db-shm

# Restart server
gorev serve
```

**Database Corruption**:
```bash
# Check database integrity
sqlite3 ~/.gorev/data/gorev.db "PRAGMA integrity_check;"

# If corrupted, restore from backup
cp ~/.gorev/data/gorev.db.backup ~/.gorev/data/gorev.db

# Or recreate (loses data)
rm ~/.gorev/data/gorev.db
gorev serve
```

**Migration Failures**:
```bash
# Check migration status
gorev serve --debug
# Look for migration error messages

# Force schema recreation (loses data)
rm ~/.gorev/data/gorev.db
rm ~/.gorev/data/migrations.lock
gorev serve
```

#### Port Conflicts
**Symptoms**: "Address already in use" error
**Solutions**:
```bash
# Check what's using port 3000
lsof -i :3000  # macOS/Linux
netstat -ano | findstr :3000  # Windows

# Use different port
gorev serve --port 8080

# Update VS Code settings
{
  "gorev.serverPort": 8080
}
```

## 🤖 AI Assistant Integration Issues

### Claude Desktop

#### MCP Configuration Not Loading
**Symptoms**: Gorev tools not available in Claude
**Solutions**:
1. **Check configuration file location**:
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

2. **Verify configuration syntax**:
   ```json
   {
     "mcpServers": {
       "gorev": {
         "command": "gorev",
         "args": ["serve", "--mcp"],
         "env": {
           "GOREV_LOG_LEVEL": "info"
         }
       }
     }
   }
   ```

3. **Test server manually**:
   ```bash
   gorev serve --mcp
   # Should start in MCP mode
   ```

4. **Restart Claude Desktop** completely

#### Tools Not Responding
**Symptoms**: Gorev tools appear but don't work
**Solutions**:
1. **Check server logs**:
   ```bash
   gorev serve --mcp --debug
   # Monitor for error messages
   ```

2. **Verify database access**:
   ```bash
   # Check if database is readable
   sqlite3 ~/.gorev/data/gorev.db ".tables"
   ```

3. **Update to latest version**:
   ```bash
   # Check version
   gorev version
   
   # Update if needed
   curl -sSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash
   ```

### VS Code Copilot/Windsurf/Cursor

#### MCP Extension Configuration
**Symptoms**: AI assistant doesn't see Gorev tools
**Solutions**:
1. **Install MCP extension** in your editor

2. **Configure settings.json**:
   ```json
   {
     "mcp.servers": {
       "gorev": {
         "command": "gorev",
         "args": ["serve", "--mcp"]
       }
     }
   }
   ```

3. **Restart editor** after configuration

#### Tool Discovery Issues
**Symptoms**: "No tools available" or similar messages
**Solutions**:
1. **Check MCP protocol version**: Ensure compatibility
2. **Verify server startup**: Look for "MCP server started" message
3. **Test direct connection**:
   ```bash
   echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | gorev serve --mcp --stdin
   ```

## 📱 Platform-Specific Issues

### Windows Subsystem for Linux (WSL)

#### Path Translation Problems
**Symptoms**: VS Code can't find gorev binary in WSL
**Solutions**:
1. **Install in both environments**:
   ```bash
   # In WSL
   curl -sSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash
   
   # In Windows PowerShell
   iwr -useb https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
   ```

2. **Use WSL path in VS Code settings**:
   ```json
   {
     "gorev.serverPath": "wsl gorev"
   }
   ```

#### Database Access Across Environments
**Symptoms**: Different data when accessing from Windows vs WSL
**Solution**: Use single data directory:
```bash
# In VS Code settings.json
{
  "gorev.dataDir": "\\\\wsl$\\Ubuntu\\home\\username\\.gorev\\data"
}
```

### Docker/Container Environments

#### Volume Mounting Issues
**Symptoms**: Data doesn't persist between container restarts
**Solution**:
```bash
# Mount data directory
docker run -v ~/.gorev:/root/.gorev gorev-image

# Or use named volume
docker volume create gorev-data
docker run -v gorev-data:/root/.gorev gorev-image
```

#### Network Connectivity
**Symptoms**: VS Code can't connect to containerized server
**Solution**:
```bash
# Expose port
docker run -p 3000:3000 gorev-image serve

# Update VS Code settings
{
  "gorev.serverHost": "localhost",
  "gorev.serverPort": 3000
}
```

## 🔍 Diagnostic Tools

### Health Check Script

Create a diagnostic script to check your installation:

```bash
#!/bin/bash
# gorev-health-check.sh

echo "=== Gorev Health Check ==="
echo

# Check binary
echo "1. Checking gorev binary..."
if command -v gorev &> /dev/null; then
    echo "✅ gorev command found: $(which gorev)"
    gorev version
else
    echo "❌ gorev command not found"
    echo "   Run: curl -sSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash"
fi
echo

# Check data directory
echo "2. Checking data directory..."
if [ -d ~/.gorev/data ]; then
    echo "✅ Data directory exists: ~/.gorev/data"
    ls -la ~/.gorev/data/
else
    echo "⚠️  Data directory not found (will be created on first run)"
fi
echo

# Check database
echo "3. Checking database..."
if [ -f ~/.gorev/data/gorev.db ]; then
    echo "✅ Database file exists"
    echo "   Tables: $(sqlite3 ~/.gorev/data/gorev.db '.tables' 2>/dev/null || echo 'Cannot read database')"
else
    echo "⚠️  Database not found (will be created on first run)"
fi
echo

# Test server startup
echo "4. Testing server startup..."
timeout 5s gorev serve --test-connection &>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ Server starts successfully"
else
    echo "❌ Server startup failed"
    echo "   Try: gorev serve --debug"
fi
echo

# Check VS Code extension
echo "5. Checking VS Code integration..."
if command -v code &> /dev/null; then
    EXTENSION_INSTALLED=$(code --list-extensions | grep mehmetsenol.gorev-vscode)
    if [ -n "$EXTENSION_INSTALLED" ]; then
        echo "✅ VS Code extension installed: $EXTENSION_INSTALLED"
    else
        echo "⚠️  VS Code extension not installed"
        echo "   Run: code --install-extension mehmetsenol.gorev-vscode"
    fi
else
    echo "⚠️  VS Code not found in PATH"
fi

echo
echo "=== Health Check Complete ==="
```

### Debug Log Analysis

When issues occur, enable debug logging:

```bash
# Start server with debug logging
gorev serve --debug --log-file debug.log

# In another terminal, reproduce the issue, then analyze logs
grep -i error debug.log
grep -i "connection" debug.log
grep -i "template" debug.log
```

### Performance Monitoring

Monitor resource usage if experiencing performance issues:

```bash
# Monitor memory usage
ps aux | grep gorev

# Monitor database size
du -h ~/.gorev/data/gorev.db

# Monitor active connections
netstat -an | grep :3000

# Check disk space
df -h ~/.gorev
```

## 🆘 Getting Help

### Before Reporting Issues

1. **Update to latest version**:
   ```bash
   gorev version  # Check current
   # Reinstall if not latest
   ```

2. **Try safe mode**: Test with minimal configuration

3. **Reproduce with debug logging**: `gorev serve --debug`

4. **Check known issues**: Visit [GitHub Issues](https://github.com/msenol/Gorev/issues)

### Information to Include in Bug Reports

- **Environment**:
  - Operating system and version
  - Gorev version (`gorev version`)
  - VS Code version (if applicable)
  - AI assistant being used

- **Configuration**:
  - Relevant settings.json entries
  - MCP configuration (sanitized)
  - Installation method used

- **Steps to reproduce**:
  - Exact commands or actions
  - Expected vs actual behavior
  - Error messages (full text)

- **Logs**:
  - Debug output from `gorev serve --debug`
  - VS Code extension output panel
  - Any relevant system logs

### Community Support

- **GitHub Discussions**: [General questions and community help](https://github.com/msenol/Gorev/discussions)
- **GitHub Issues**: [Bug reports and feature requests](https://github.com/msenol/Gorev/issues)
- **Documentation**: [Complete documentation](https://github.com/msenol/Gorev/tree/main/docs)

### Professional Support

For enterprise users or complex deployments:
- Email: [Contact information in repository]
- Custom configuration assistance
- Performance optimization
- Integration consulting

---

*✨ This comprehensive troubleshooting guide covers the most common issues and their solutions. If you encounter a problem not covered here, please check the GitHub issues or start a discussion for community help.*