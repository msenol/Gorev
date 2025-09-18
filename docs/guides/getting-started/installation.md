# Installation Guide

> **Version**: This documentation is valid for v0.15.5+
> **Last Updated**: September 18, 2025

Complete installation instructions for Gorev on all platforms.

## 📋 Prerequisites

- **Operating System**: Linux, macOS, Windows
- **MCP Compatible Editor**: Claude Desktop, VS Code (with MCP extension), Windsurf, Cursor, Zed, or other MCP-supported editors
- **Docker** (optional, for container deployment)

## 🚀 Quick Installation

### 🚀 NPX Easy Installation (Recommended)

**For MCP clients, simply add to your configuration:**

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

**Traditional Installation:**

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

### ✅ Verify Installation

```bash
gorev version
gorev help
```

## 🔧 MCP Editor Configuration

### 🤖 Claude Desktop

Add this configuration to your Claude Desktop config file:

**File Locations:**
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "gorev": {
      "command": "/usr/local/bin/gorev",
      "args": ["serve"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev"
      }
    }
  }
}
```

> **Windows Users**: Use the full path like `"C:\\Users\\YourName\\AppData\\Local\\Programs\\gorev\\gorev.exe"`

### 💻 VS Code

**Option 1: Gorev VS Code Extension (Recommended)**

Install from the [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode):

```bash
code --install-extension mehmetsenol.gorev-vscode
```

**Option 2: MCP Extension**

Install an MCP extension and add to `settings.json`:

```json
{
  "mcp.servers": {
    "gorev": {
      "command": "/usr/local/bin/gorev",
      "args": ["serve"]
    }
  }
}
```

## 🐳 Docker Installation

### Run with Docker

```bash
# Pull the image
docker pull ghcr.io/msenol/gorev:latest

# Create volume for data persistence
docker volume create gorev-data

# Run the server
docker run -d --name gorev-server \
  -v gorev-data:/data \
  -p 8080:8080 \
  ghcr.io/msenol/gorev:latest serve
```

### MCP Configuration with Docker

```json
{
  "mcpServers": {
    "gorev": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-v", "gorev-data:/data",
        "ghcr.io/msenol/gorev:latest",
        "serve"
      ]
    }
  }
}
```

## 📁 Data Directory Locations

| Platform | Default Location | Environment Variable |
|----------|------------------|---------------------|
| Windows | `%USERPROFILE%\.gorev\` | `GOREV_DATA_DIR` |
| macOS | `~/.gorev/` | `GOREV_DATA_DIR` |
| Linux | `~/.gorev/` | `GOREV_DATA_DIR` |
| Docker | `/data` (volume) | N/A |

## ⚙️ Advanced Configuration

### Custom Port

```json
{
  "mcpServers": {
    "gorev": {
      "command": "gorev",
      "args": ["serve", "--port", "8080"]
    }
  }
}
```

### Debug Mode

```json
{
  "mcpServers": {
    "gorev": {
      "command": "gorev",
      "args": ["serve", "--debug"],
      "env": {
        "GOREV_LOG_LEVEL": "debug"
      }
    }
  }
}
```

### Multiple Instances

```json
{
  "mcpServers": {
    "gorev-personal": {
      "command": "gorev",
      "args": ["serve", "--port", "8081"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev-personal"
      }
    },
    "gorev-work": {
      "command": "gorev",
      "args": ["serve", "--port", "8082"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev-work"
      }
    }
  }
}
```

## 🛠️ Platform-Specific Installation

### 🍎 macOS

**Homebrew (Recommended):**
```bash
brew tap msenol/gorev
brew install gorev
```

**Manual Binary:**
```bash
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-darwin-amd64 -o gorev
chmod +x gorev
sudo mv gorev /usr/local/bin/
```

**Security Note:** First run may require bypassing Gatekeeper:
- System Preferences → Security & Privacy → General → "Allow Anyway"

### 🪟 Windows

**PowerShell (Admin required):**
```powershell
# Download and install
$url = "https://github.com/msenol/gorev/releases/latest/download/gorev-windows-amd64.exe"
$path = "$env:LOCALAPPDATA\Programs\gorev"
New-Item -ItemType Directory -Force -Path $path
Invoke-WebRequest -Uri $url -OutFile "$path\gorev.exe"

# Add to PATH
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$path", "User")
```

**Scoop (Alternative):**
```powershell
scoop bucket add gorev https://github.com/msenol/scoop-gorev
scoop install gorev
```

### 🐧 Linux

**Package Managers:**

**Debian/Ubuntu:**
```bash
sudo add-apt-repository ppa:msenol/gorev
sudo apt update && sudo apt install gorev
```

**Arch Linux (AUR):**
```bash
yay -S gorev-bin
```

**Manual Binary:**
```bash
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev
sudo mv gorev /usr/local/bin/
```

## 🧪 Testing Your Installation

### 1. CLI Test
```bash
gorev version
gorev serve --test
```

### 2. MCP Editor Test
Restart your MCP editor and test with your AI assistant:
```
"Is Gorev working? Create a test task."
```

### 3. Check Logs
```bash
# View server logs
tail -f ~/.gorev/logs/gorev.log
```

## 🆘 Troubleshooting

### Common Issues

**"Command not found":**
- Restart your terminal/editor
- Check PATH environment variable
- Verify binary location

**MCP Connection Failed:**
1. Close your editor completely
2. Check JSON syntax in config file
3. Test `gorev serve` manually
4. Restart your editor

**Database Locked Error:**
```bash
pkill -f gorev
rm ~/.gorev/gorev.db-wal ~/.gorev/gorev.db-shm
```

**VS Code Extension Issues:**
- Ensure extension is up to date
- Check Developer Console for errors
- Refresh MCP server list

### Platform-Specific Issues

**Windows - "Windows protected your PC":**
- Click "More info" → "Run anyway"
- Or add Windows Defender exclusion

**macOS - "Cannot be opened because developer cannot be verified":**
```bash
xattr -d com.apple.quarantine /usr/local/bin/gorev
```

**Linux - Permission denied:**
```bash
chmod +x /usr/local/bin/gorev
sudo restorecon -v /usr/local/bin/gorev  # SELinux systems
```

## 🔄 Updating

### Automatic Update
```bash
gorev update
```

### Package Manager Update
```bash
brew upgrade gorev        # macOS
scoop update gorev       # Windows
sudo apt upgrade gorev   # Debian/Ubuntu
```

## ❌ Uninstalling

### Remove Gorev
```bash
# Package manager
brew uninstall gorev      # macOS
scoop uninstall gorev    # Windows
sudo apt remove gorev    # Linux

# Manual removal
sudo rm /usr/local/bin/gorev
rm -rf ~/.gorev  # WARNING: This deletes all your data
```

## 📚 Next Steps

After installation:

1. **[Quick Start Guide](../user/usage.md)** - Learn basic usage
2. **[MCP Tools Reference](../user/mcp-tools.md)** - Complete tool documentation
3. **[VS Code Extension Guide](../user/vscode-extension.md)** - Visual interface
4. **[Troubleshooting](../../debugging/)** - Common issues and solutions

## 💬 Support

Having trouble? Get help:

- [GitHub Issues](https://github.com/msenol/gorev/issues)
- [Discussions](https://github.com/msenol/gorev/discussions)
- [Documentation](https://github.com/msenol/gorev/tree/main/docs)

---

*This installation guide was created with assistance from Claude (Anthropic)*