# 🚀 Gorev

<div align="center">

**Last Updated:** September 4, 2025 | **Version:** v0.13.1

[🇺🇸 English](README.en.md) | [🇹🇷 Türkçe](README.md)

> ⚠️ **BREAKING CHANGE (v0.10.0)**: The `gorev_olustur` tool is no longer available! Template usage is now mandatory. [Details](#breaking-change-template-requirement)

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)
![MCP](https://img.shields.io/badge/MCP-Compatible-4A154B?style=flat-square&logo=anthropic)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Test Coverage](https://img.shields.io/badge/Coverage-84.6%25-brightgreen?style=flat-square)
![Platform](https://img.shields.io/badge/Platform-Linux%20|%20macOS%20|%20Windows-blue?style=flat-square)

**Modern task management system with Turkish support, designed for MCP-compatible AI assistants (Claude, VS Code, Windsurf, Cursor)**

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Documentation](#-documentation) • [Contributing](#-contributing)

</div>

## 🎯 What is Gorev?

Gorev is a powerful **Model Context Protocol (MCP)** server written in Go that provides task management capabilities to all MCP-compatible AI editors (Claude Desktop, VS Code, Windsurf, Cursor, Zed, etc.). It combines project management, task tracking, and organization needs with the power of AI assistants to boost your productivity.

### 🏗️ Two-Module Architecture

1. **gorev-mcpserver** - MCP server written in Go (core component)
2. **gorev-vscode** - VS Code extension (optional visual interface)

Thanks to the MCP protocol, you can connect to the server from any MCP-compatible editor. The VS Code extension provides a rich visual experience.

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

## 📦 Installation

### 🚀 Automatic Installation (Recommended)

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

### MCP Editor Integration

#### 🤖 Claude Desktop

Add to your configuration file:
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "gorev": {
      "command": "/path/to/gorev",
      "args": ["serve"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev"
      }
    }
  }
}
```

#### 💻 VS Code

**Option 1: Gorev VS Code Extension (Recommended)**

Install from [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)

```bash
code --install-extension mehmetsenol.gorev-vscode
```

**Option 2: MCP Extension**

Install MCP extension and add to `settings.json`:

```json
{
  "mcp.servers": {
    "gorev": {
      "command": "/path/to/gorev",
      "args": ["serve"]
    }
  }
}
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

- 📦 [Installation Guide](docs/getting-started/installation.md) - Platform-specific installation instructions
- 📖 [Usage Guide](docs/user-guide/usage.md) - Detailed usage examples
- 🛠 [MCP Tools](docs/user-guide/mcp-tools.md) - Complete reference for 29 MCP tools
- 🤖 [AI MCP Tools](docs/mcp-araclari-ai.md) - AI context management tools (v0.9.0)
- 🏗 [System Architecture](docs/development/architecture.md) - Technical details
- 🗺️ [Roadmap](ROADMAP.md) - Development roadmap and future plans
- 💻 [Contributing Guide](docs/development/contributing.md) - How to contribute
- 🎨 [VS Code Extension](docs/user-guide/vscode-extension.md) - Extension documentation

### AI Assistant Documentation
- 🌍 [CLAUDE.en.md](CLAUDE.en.md) - English AI assistant guidance
- 🤖 [CLAUDE.md](CLAUDE.md) - Turkish AI assistant guidance
- 📋 [MCP Tools Reference](docs/MCP_TOOLS_REFERENCE.md) - Detailed MCP tool documentation
- 📚 [Development History](docs/DEVELOPMENT_HISTORY.md) - Complete project history

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
- Go 1.22+
- Make (optional)
- golangci-lint (for code quality)

### Commands

```bash
# Download dependencies
make deps

# Run tests (84.6% overall coverage)
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

- **Version**: v0.11.0
- **Test Coverage**: 84.6%
- **Go Version**: 1.22+
- **MCP SDK**: mark3labs/mcp-go v0.6.0
- **Database**: SQLite (embedded)

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