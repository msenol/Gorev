# 🚀 Gorev

<div align="center">

**Last Updated:** October 5, 2025 | **Version:** v0.16.2

[🇺🇸 English](README.md) | [🇹🇷 Türkçe](README.tr.md)

> 🎉 **NEW in v0.16.2**: Critical NPM binary update fix + VS Code auto-start! [See What's New](#-whats-new-in-v0162)

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)
![MCP](https://img.shields.io/badge/MCP-Compatible-4A154B?style=flat-square&logo=anthropic)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Test Coverage](https://img.shields.io/badge/Coverage-75%25-yellow?style=flat-square)
![Platform](https://img.shields.io/badge/Platform-Linux%20|%20macOS%20|%20Windows-blue?style=flat-square)

**Modern task management system with Turkish support, designed for MCP-compatible AI assistants (Claude, VS Code, Windsurf, Cursor)**

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Documentation](#-documentation) • [Contributing](#-contributing)

</div>

## 🎯 What is Gorev

Gorev is a powerful **Model Context Protocol (MCP)** server written in Go that provides task management capabilities to all MCP-compatible AI editors (Claude Desktop, VS Code, Windsurf, Cursor, Zed, etc.). It combines project management, task tracking, and organization needs with the power of AI assistants to boost your productivity.

### 🏗️ Three-Module Architecture

1. **gorev-mcpserver** - MCP server written in Go (core component)
   - Embedded Web UI 🌐 - React interface embedded in binary (NEW! v0.16.0)
   - REST API server (Fiber framework)
   - MCP protocol support
2. **gorev-vscode** - VS Code extension (optional visual interface)
3. **gorev-web** - React + TypeScript source code (development)

Thanks to the MCP protocol, you can connect to the server from any MCP-compatible editor. The Web UI is automatically available at http://localhost:5082 when you run `npx @mehmetsenol/gorev-mcp-server serve`. The VS Code extension provides a rich IDE-integrated experience.

## 🎉 What's New in v0.16.2

### 🐛 Critical Bug Fixes (v0.16.2)
- **NPM Binary Update Fix**: Fixed critical bug where NPM package upgrades preserved old binaries
  - Users upgrading from v0.16.1 or earlier were stuck on v0.15.24 (September 2025)
  - Package size reduced from 78.4 MB to 6.9 KB (binaries now always downloaded from GitHub)
  - All users now get latest features (REST API, Web UI, VS Code auto-start)
- **VS Code Auto-Start**: Extension now automatically starts server on activation
  - No manual `npx gorev serve` required
  - Checks if server is running, starts if needed
  - Proper database path configuration (workspace/.gorev/gorev.db)
  - Graceful server shutdown on extension deactivation

### 🌐 Embedded Web UI (v0.16.0)
- **Zero-Configuration**: Modern React interface built into Go binary
- **Instant Access**: Automatically available at http://localhost:5082
- **Full Features**: Tasks, projects, templates, subtasks, and dependencies
- **Language Sync**: Turkish/English switcher synchronized with MCP server
- **No Separate Installation**: Just run `npx @mehmetsenol/gorev-mcp-server serve` and you're ready!

### 🗂️ Multi-Workspace Support (v0.16.0)
- **Isolated Workspaces**: Each project folder gets its own task database
- **Workspace Switcher**: Seamlessly switch between workspaces in Web UI
- **Auto-Detection**: Automatically detects `.gorev/` directory in current folder
- **SHA256 IDs**: Secure workspace identification
- **VS Code Integration**: Extension auto-registers workspace on activation

### 🔌 REST API Migration
- **23 Endpoints**: Complete Fiber-based REST API
- **VS Code Extension**: Migrated from MCP to REST API for better performance
- **Type-Safe**: 100% TypeScript with zero parsing errors
- **Faster**: Direct HTTP calls vs. stdio + markdown parsing
- **Backward Compatible**: MCP protocol still fully supported

### 🏷️ Template Aliases
- **Quick Commands**: Use `bug`, `feature`, `research` instead of template IDs
- **Consistency**: Same aliases across all workspaces
- **No More UUID Hunting**: Human-readable template identifiers
- **Documentation**: Full guide at [MCP Config Examples](docs/guides/mcp-config-examples.md)

### 📦 NPM Package
- **Package Name**: `@mehmetsenol/gorev-mcp-server`
- **Global Install**: `npm install -g @mehmetsenol/gorev-mcp-server`
- **NPX Ready**: `npx @mehmetsenol/gorev-mcp-server serve` for instant use
- **Cross-Platform**: Works on Windows, macOS, and Linux

## ✨ Features

### 📝 Task Management

- **Smart task creation** - Using natural language commands
- **Markdown support** - Rich description formatting
- **Status management** - Pending → In Progress → Completed
- **Priority levels** - Low, Medium, High
- **Flexible editing** - Update all task properties

### 📁 Project Organization

- **Hierarchical structure** - Task grouping under projects
- **Active project system** - Quick operations with default project
- **Project-based reporting** - Detailed statistics
- **Multi-project support** - Unlimited project creation

### 🔗 Advanced Features

- **📅 Due date tracking** - Deadline management and urgent task filtering
- **🏷️ Tagging system** - Multi-tag categorization
- **🔄 Task dependencies** - Inter-task automation
- **📋 Ready-made templates** - Bug reports, feature requests, and more
- **🔍 Advanced filtering** - Status, tag, date-based queries
- **🌳 Subtask hierarchy** - Unlimited depth task tree structure
- **📊 Progress tracking** - Subtask completion percentage in parent tasks
- **📁 File System Watcher** - Monitor file changes and automatic task status transitions
- **🔔 Automatic Status Updates** - "pending" → "in_progress" automation on file changes
- **⚙️ Configuration Management** - Customizable ignore patterns and watch rules

### 🤖 AI Integration

- **Natural language processing** - Task management by talking to AI assistants
- **Multi-editor support** - Claude, VS Code, Windsurf, Cursor, Zed
- **Contextual understanding** - Smart command interpretation
- **MCP standard** - Compatible with all MCP-compatible tools

### 🎨 VS Code Extension Features (Optional)

- **Bilingual Support** - Turkish and English interface (v0.5.0+) 🌍
- **TreeView Panels** - Task, project, and template lists
- **Visual Interface** - Click-and-use experience
- **Status Bar** - Real-time status information
- **Command Palette** - Quick access (Ctrl+Shift+G)
- **Color Coding** - Priority-based visual distinction
- **Context Menus** - Right-click operations
- **Automatic Language Detection** - UI language based on VS Code language setting
- **[Download from Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** 🚀

### 🌐 Web UI Features (NEW! v0.16.0)

- **Modern Browser Interface** - No IDE required, works in any browser
- **Task Cards** - Rich task visualization with metadata
- **Subtask Hierarchy** - Expandable/collapsible nested tasks
- **Dependency Badges** - Visual indicators for task dependencies
- **Project Organization** - Sidebar navigation with task counts
- **Template-Based Creation** - Wizard for creating structured tasks
- **Real-time Updates** - React Query for automatic synchronization
- **Responsive Design** - Works on desktop and mobile devices
- **🌍 Language Switcher** - Toggle between Turkish/English, synchronized with MCP server
- **Quick Actions** - Edit, delete, and status updates
- **REST API Backend** - Fiber-based high-performance API server
- **🚀 No Installation Required**: Automatically active with `npx @mehmetsenol/gorev-mcp-server serve`!
- **Embedded UI**: Bundled in Go binary, no separate setup needed
- **Access**: http://localhost:5082 (default port)

## 📦 Installation

### 🚀 NPM Quick Setup (Recommended!)

> ⚠️ **Windows Users**: NPM requires Node.js installation. [Download Node.js](https://nodejs.org/) and restart your system after installation.

#### Global Installation

```bash
npm install -g @mehmetsenol/gorev-mcp-server
```

Or use directly with NPX (no installation required):

```bash
npx @mehmetsenol/gorev-mcp-server serve
```

#### MCP Client Configuration

**For Claude Desktop:**

```json
// Windows: %APPDATA%/Claude/claude_desktop_config.json
// macOS: ~/Library/Application Support/Claude/claude_desktop_config.json
// Linux: ~/.config/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

**For Kilo Code (VS Code Extension):**

```json
// .kilocode/mcp.json (workspace root)
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

**For Cursor:**

```json
// .cursor/mcp.json (workspace root)
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

> 📚 **More Examples**: [MCP Configuration Guide](docs/guides/mcp-config-examples.md)

#### 🔧 Windows NPX Troubleshooting

If you get `ENOENT: spawn npx` error:

1. **Check if Node.js is installed:**

   ```cmd
   node --version
   npm --version
   npx --version
   ```

2. **Install Node.js:**
   - Download LTS version from [Node.js website](https://nodejs.org/)
   - Check "Add to PATH" option during installation
   - Restart your computer after installation

3. **Install NPX separately (if needed):**

   ```cmd
   npm install -g npx
   ```

4. **Check PATH:**

   ```cmd
   echo %PATH%
   ```

   Should include Node.js paths (`C:\Program Files\nodejs\`).

**For Cursor:**

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

### 🔧 Traditional Installation (Automatic)

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash

# Specific version
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | VERSION=v0.11.0 bash
```

### 🪟 Windows

```powershell
# PowerShell (no admin rights required)
irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex

# Or for specific version:
$env:VERSION="v0.11.0"; irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

### 💻 VS Code Extension (Optional)

**Option 1: Gorev VS Code Extension (Recommended)**

Install from [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)

```bash
code --install-extension mehmetsenol.gorev-vscode
```

## 🎮 Usage

### AI Assistant Example Commands

```
"Create a new task: Write API documentation"
"List urgent tasks"
"Show tasks tagged with bug"
"Set Mobile App v2 project as active"
"Create new project for sprint planning"
"Mark task #5 as completed"
"Create new task from feature request template"
"Start watching project files"
"Enable automatic status transitions on file changes"
"Show watch list"
"Add git ignore rules to file watcher"
```

> 💡 **Tip**: These commands work with Claude, VS Code Copilot, Windsurf AI, Cursor, and other MCP-compatible AI assistants.

### CLI Commands

```bash
# Start server
gorev serve                  # Normal mode
gorev serve --debug         # Debug mode
gorev serve --port 8080     # Different port

# Task operations
gorev task list             # List tasks
gorev task create           # Create new task
gorev task show <id>        # Task details

# Project operations
gorev project list          # List projects
gorev project create        # Create new project

# Other
gorev version              # Version info
gorev help                 # Help
```

## 📚 Documentation

For detailed documentation, see the [docs/](docs/) folder:

### Getting Started

- 🚀 [Quick Start Guide](docs/guides/getting-started/quick-start.md) - Get up and running in 10 minutes
- 📦 [Installation Guide](docs/guides/getting-started/installation.md) - Platform-specific installation instructions
- 🆘 [Troubleshooting Guide](docs/guides/getting-started/troubleshooting.md) - Common issues and solutions
- 🔄 [Migration Guide (v0.15→v0.16)](docs/migration/v0.15-to-v0.16.md) - Upgrade from v0.15

### Features

- 🌐 [Web UI Guide](docs/guides/features/web-ui.md) - Embedded React interface documentation
- 🗂️ [Multi-Workspace Support](docs/guides/features/multi-workspace.md) - Managing multiple projects
- 📋 [Template System](docs/guides/features/template-system.md) - Structured task creation
- 🤖 [AI Context Management](docs/guides/features/ai-context-management.md) - AI assistant integration

### Reference

- 🛠️ [MCP Tools Reference](docs/legacy/tr/mcp-araclari.md) - Complete reference for 41 MCP tools
- 🔧 [MCP Configuration Examples](docs/guides/mcp-config-examples.md) - IDE setup guides
- 📖 [Usage Guide](docs/guides/user/usage.md) - Detailed usage examples
- 🎨 [VS Code Extension](docs/guides/user/vscode-extension.md) - Extension documentation

### Development

- 🏗️ [System Architecture](docs/architecture/architecture-v2.md) - Technical details
- 💻 [Contributing Guide](docs/development/contributing.md) - How to contribute
- 🗺️ [Roadmap](ROADMAP.md) - Development roadmap and future plans
- 📚 [Development History](docs/development/TASKS.md) - Complete project history

### AI Assistant Documentation

- 🌍 [CLAUDE.en.md](CLAUDE.en.md) - English AI assistant guidance
- 🤖 [CLAUDE.md](CLAUDE.md) - Turkish AI assistant guidance
- 📋 [MCP Tools Reference](docs/api/MCP_TOOLS_REFERENCE.md) - Detailed MCP tool documentation
- 📚 [Development History](docs/development/TASKS.md) - Complete project history

## 🏗 Architecture

### Project Structure

```
gorev/
├── gorev-mcpserver/        # MCP Server (Go)
│   ├── cmd/gorev/         # CLI and server entry point
│   ├── internal/
│   │   ├── mcp/           # MCP protocol layer
│   │   └── gorev/        # Business logic
│   └── test/              # Integration tests
├── gorev-vscode/           # VS Code Extension (TypeScript)
│   ├── src/
│   │   ├── commands/      # VS Code commands
│   │   ├── providers/     # TreeView providers
│   │   └── mcp/           # MCP client
│   └── package.json       # Extension manifest
└── docs/                   # Project documentation
```

## 🧪 Development

### Requirements

- Go 1.23+
- Make (optional)
- golangci-lint (for code quality)

### Commands

```bash
# Download dependencies
make deps

# Run tests (90%+ overall coverage)
make test

# Coverage report
make test-coverage

# Lint check
make lint

# Build (all platforms)
make build-all

# Docker image
make docker-build
```

## 📊 Project Status

- **Version**: v0.16.2 🚀
- **Test Coverage**: 75%+ (Comprehensive test coverage with ongoing improvements)
- **Go Version**: 1.23+
- **MCP SDK**: mark3labs/mcp-go v0.6.0
- **Database**: SQLite (embedded)
- **Security**: Production-ready audit compliant
- **Thread Safety**: 100% race condition free

## 🤝 Community

- 📦 [GitHub Releases](https://github.com/msenol/gorev/releases)
- 🐛 [Issue Tracker](https://github.com/msenol/gorev/issues)
- 💬 [Discussions](https://github.com/msenol/gorev/discussions)
- 📖 [Wiki](https://github.com/msenol/gorev/wiki)

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🚨 Breaking Change: Template Requirement

**Starting from v0.10.0**, the `gorev_olustur` tool has been removed. All task creation must now use the template system for better structure and consistency.

### Migration Guide

**Before (v0.9.x and earlier):**

```
Create a new task: Fix login bug
```

**After (v0.10.0+):**

```
Use bug-report template to create: Fix login bug
```

Available templates:

- `bug-report` - Bug reports and fixes
- `feature` - New features and enhancements  
- `task` - General tasks and activities
- `meeting` - Meeting planning and notes
- `research` - Research and investigation tasks

For more details, see [MCP Tools Documentation](docs/user-guide/mcp-tools.md#gorev_template_olustur).

---

<div align="center">

Made with ❤️ by [msenol](https://github.com/msenol/gorev/graphs/contributors)

📚 *Documentation enhanced by Claude (Anthropic) - Your AI pair programming assistant*

**[⬆ Back to Top](#-gorev)**

</div>
