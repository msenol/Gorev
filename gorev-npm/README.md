# @mehmetsenol/gorev-mcp-server

[![npm version](https://badge.fury.io/js/%40mehmetsenol%2Fgorev-mcp-server.svg)](https://badge.fury.io/js/%40mehmetsenol%2Fgorev-mcp-server)
[![Downloads](https://img.shields.io/npm/dm/@mehmetsenol/gorev-mcp-server.svg)](https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server)
[![License](https://img.shields.io/npm/l/@mehmetsenol/gorev-mcp-server.svg)](https://github.com/msenol/Gorev/blob/main/LICENSE)

**Gorev MCP Server** - Modern task management system for AI assistants via Model Context Protocol (MCP)

🇹🇷 **Turkish**: Task management with natural language support | 🇺🇸 **English**: Full bilingual support

## ⚠️ Important: v0.16.2 Upgrade Notice

**If you're upgrading from v0.16.1 or earlier**, you may need to reinstall to get the latest binary (v0.16.2):

```bash
# Uninstall old version
npm uninstall -g @mehmetsenol/gorev-mcp-server

# Clear NPM cache
npm cache clean --force

# Install latest version
npm install -g @mehmetsenol/gorev-mcp-server@latest

# Verify version (should show v0.16.2)
gorev-mcp version
```

**What was fixed in v0.16.2:**
- 🐛 Critical bug where package upgrades preserved old binaries (users were stuck on v0.15.24)
- 📦 Package size reduced from 78.4 MB to 6.9 KB
- 🚀 Binaries now always downloaded from GitHub releases
- ✨ All users now get latest features (REST API, Web UI, VS Code auto-start)

## 🚀 Quick Start

### Using npx (Recommended)

No installation required! Use directly with MCP clients:

```bash
npx @mehmetsenol/gorev-mcp-server@latest
```

### MCP Configuration

Add to your `mcp.json` configuration file:

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

### Supported MCP Clients

- ✅ **Claude Desktop** - AI assistant with MCP support
- ✅ **VS Code** - With MCP extension
- ✅ **Cursor** - AI-powered code editor
- ✅ **Windsurf** - AI development environment
- ✅ **Any MCP-compatible client**

## 🔧 Configuration Examples

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%/Claude/claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

### VS Code with MCP Extension

Add to `.vscode/mcp.json` in your workspace:

```json
{
  "servers": {
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

### Cursor IDE

Add to your Cursor MCP configuration:

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"]
    }
  }
}
```

## 🌍 Environment Variables

| Variable | Description | Default | Values |
|----------|-------------|---------|--------|
| `GOREV_LANG` | Interface language | `tr` | `tr`, `en` |
| `GOREV_DB_PATH` | Custom database path | Auto-detected | Any valid path |

## 📦 Installation Methods

### Method 1: npx (Recommended)

```bash
# No installation needed - runs directly
npx @mehmetsenol/gorev-mcp-server@latest --help
```

### Method 2: Global Installation

```bash
# Install globally
npm install -g @mehmetsenol/gorev-mcp-server

# Run directly
gorev-mcp --help
```

### Method 3: Local Installation

```bash
# Install in project
npm install @mehmetsenol/gorev-mcp-server

# Run with npx
npx @mehmetsenol/gorev-mcp-server serve
```

## 🛠️ Commands

```bash
# Start MCP server (default command)
npx @mehmetsenol/gorev-mcp-server

# Initialize database
npx @mehmetsenol/gorev-mcp-server init

# Initialize with global database
npx @mehmetsenol/gorev-mcp-server init --global

# Initialize default templates
npx @mehmetsenol/gorev-mcp-server template init

# Show help
npx @mehmetsenol/gorev-mcp-server --help

# Show version
npx @mehmetsenol/gorev-mcp-server --version
```

## ✨ Features

### 📝 Task Management

- **Natural language task creation** - Create tasks by talking to AI
- **Unlimited subtask hierarchy** - Organize tasks in unlimited depth
- **Smart status tracking** - Pending → In Progress → Completed
- **Priority levels** - Low, Medium, High priority management
- **Due date tracking** - Deadline management with overdue detection

### 📁 Project Organization

- **Hierarchical projects** - Group tasks under projects
- **Active project system** - Quick operations with default project
- **Multi-project support** - Unlimited project creation
- **Project-based reporting** - Detailed statistics per project

### 🔗 Advanced Features

- **Task dependencies** - Link related tasks with automation
- **Tagging system** - Multi-tag categorization
- **Template system** - Pre-built task templates (Bug, Feature, etc.)
- **File system watching** - Auto-update tasks based on file changes
- **Advanced search** - Full-text search with filtering
- **Data export/import** - JSON/CSV export with conflict resolution

### 🤖 AI Integration

- **48 MCP tools** - Complete API for AI assistants
- **Natural language processing** - Smart command interpretation
- **Bilingual support** - Turkish and English interfaces
- **Cross-platform** - Windows, macOS, Linux support

## 📊 MCP Tools Available

The server provides 48 MCP tools organized in categories:

- **Task Management** (6 tools): `gorev_listele`, `gorev_detay`, `gorev_guncelle`, etc.
- **Subtask Management** (3 tools): `gorev_altgorev_olustur`, etc.
- **Project Management** (6 tools): `proje_olustur`, `proje_listele`, etc.
- **Templates** (2 tools): `template_listele`, `templateden_gorev_olustur`
- **Advanced Search** (6 tools): `gorev_search_advanced`, `gorev_filter_profile_*`, etc.
- **Data Export/Import** (2 tools): `gorev_export`, `gorev_import`
- **File Watching** (4 tools): `gorev_file_watch_add`, etc.
- **AI Context** (6 tools): `gorev_set_active`, `gorev_nlp_query`, etc.
- **IDE Extension Management** (5 tools): `ide_detect`, `ide_install_extension`, etc.
- **Advanced** (8 tools): Various utility and management tools

## 🔧 Troubleshooting

### Binary Download Issues

```bash
# Clear npm cache and reinstall
npm cache clean --force
npm uninstall -g @mehmetsenol/gorev-mcp-server
npm install -g @mehmetsenol/gorev-mcp-server@latest
```

### Platform Not Supported

```bash
# Check supported platforms
npx @mehmetsenol/gorev-mcp-server --help

# Report issue with platform info
echo "Platform: $(uname -a)"
```

### MCP Connection Issues

```bash
# Test server startup
npx @mehmetsenol/gorev-mcp-server serve

# Check database initialization
npx @mehmetsenol/gorev-mcp-server init
```

## 📱 Supported Platforms

| Platform | Architecture | Status |
|----------|-------------|--------|
| Windows | x64 (amd64) | ✅ |
| macOS | x64 (amd64) | ✅ |
| macOS | ARM64 (Apple Silicon) | ✅ |
| Linux | x64 (amd64) | ✅ |
| Linux | ARM64 | ✅ |

## 🆚 Related Packages

- **[VS Code Extension](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** - Rich visual interface
- **[Main Repository](https://github.com/msenol/Gorev)** - Source code and documentation

## 📚 Documentation

- [Installation Guide](https://github.com/msenol/Gorev#-kurulum)
- [MCP Tools Reference](https://github.com/msenol/Gorev/blob/main/docs/mcp-araclari.md)
- [Development Guide](https://github.com/msenol/Gorev/blob/main/CLAUDE.md)
- [API Documentation](https://github.com/msenol/Gorev/tree/main/docs)

## 🐛 Issues & Support

- **Report Issues**: [GitHub Issues](https://github.com/msenol/Gorev/issues)
- **Feature Requests**: [GitHub Discussions](https://github.com/msenol/Gorev/discussions)
- **Documentation**: [Project Wiki](https://github.com/msenol/Gorev/wiki)

## 📄 License

MIT License - see [LICENSE](https://github.com/msenol/Gorev/blob/main/LICENSE) file for details.

## 🔄 Updates

This package automatically downloads the latest Gorev binaries from GitHub releases. To update:

```bash
# Update to latest version
npm update -g @mehmetsenol/gorev-mcp-server

# Or use npx for always latest
npx @mehmetsenol/gorev-mcp-server@latest
```

---

**Built with ❤️ by [Mehmet Senol](https://github.com/msenol)** | **Powered by Go & Model Context Protocol**
