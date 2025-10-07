# Quick Start Guide

**Version**: v0.16.0
**Est. Time**: 10 minutes
**Last Updated**: October 5, 2025

Get up and running with Gorev in minutes!

---

## What is Gorev

Gorev is a task management system designed for AI assistants (Claude, Copilot, etc.) with:

- **MCP Protocol**: 41 tools for AI-powered task management
- **Embedded Web UI**: Browser-based interface at http://localhost:5082
- **VS Code Extension**: Rich IDE integration (optional)
- **Multi-Workspace**: Isolated databases per project
- **Template System**: 6 default templates for structured tasks

---

## Installation (2 minutes)

### Option 1: NPM (Recommended)

```bash
# Use directly with NPX (no installation)
npx @mehmetsenol/gorev-mcp-server@latest
```

Or install globally:

```bash
npm install -g @mehmetsenol/gorev-mcp-server
gorev daemon --detach
```

### Option 2: Download Binary

**Linux/macOS**:

```bash
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash
```

**Windows PowerShell**:

```powershell
irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

---

## First Steps (3 minutes)

### 1. Initialize Workspace

```bash
cd /path/to/your/project
npx @mehmetsenol/gorev-mcp-server init
```

**Output**:

```
✓ Created .gorev/ directory
✓ Initialized database: gorev.db
✓ Workspace registered: my-project
```

### 2. Start Server

```bash
npx @mehmetsenol/gorev-mcp-server@latest
```

**Output**:

```
🚀 Gorev MCP Server v0.16.0
📦 Web UI: http://localhost:5082
🔌 MCP: Listening on stdio
✅ Ready for connections
```

### 3. Access Web UI

Open browser: **http://localhost:5082**

You'll see:

- Workspace switcher (top-right)
- Project sidebar (left)
- Task list (center)
- "New Task" button

---

## Create Your First Task (2 minutes)

### Via Web UI

1. Click **"New Task"** button
2. Select **"bug"** template
3. Fill in fields:
   - Title: `Fix login button`
   - Module: `auth`
   - Environment: `production`
   - Steps: `1. Click login\n2. Button doesn't respond`
   - Expected: `User logs in`
   - Actual: `Nothing happens`
   - Priority: `high`
4. Click **"Create Task"**

✅ Task created and appears in task list!

### Via CLI

```bash
npx @mehmetsenol/gorev-mcp-server task create \
  --template bug \
  --field baslik="Fix login button" \
  --field modul="auth" \
  --field ortam="production" \
  --field oncelik="yuksek"
```

---

## Configure AI Assistant (3 minutes)

### Claude Desktop

**File**: `~/.config/Claude/claude_desktop_config.json` (Linux/macOS)
**File**: `%APPDATA%\Claude\claude_desktop_config.json` (Windows)

```json
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
```

**Restart Claude Desktop**

### Test Connection

In Claude:

```
List all my tasks
```

Claude should respond with your task list!

### VS Code (Kilo Code Extension)

**File**: `.kilocode/mcp.json` in workspace root

```json
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
```

**Reload VS Code window**

---

## Common Tasks

### Create Project

**Web UI**: Sidebar → "New Project" → Enter name

**CLI**:

```bash
npx @mehmetsenol/gorev-mcp-server project create --name "My Project"
```

**AI Assistant**:

```
Create a new project called "Backend API"
```

### List Tasks

**Web UI**: View in center panel (auto-updates)

**CLI**:

```bash
npx @mehmetsenol/gorev-mcp-server task list
```

**AI Assistant**:

```
Show me all my tasks
```

### Update Task Status

**Web UI**: Click task → Change status dropdown

**CLI**:

```bash
npx @mehmetsenol/gorev-mcp-server task update --id <ID> --status completed
```

**AI Assistant**:

```
Mark the login button task as completed
```

### Switch Workspace

**Web UI**: Top-right dropdown → Select workspace

**CLI**: Change directory and run commands

```bash
cd /path/to/other/project
npx @mehmetsenol/gorev-mcp-server task list
```

---

## Next Steps

### Learn Features

- 📚 [Web UI Guide](../features/web-ui.md) - Browser interface details
- 🗂️ [Multi-Workspace](../features/multi-workspace.md) - Manage multiple projects
- 📋 [Template System](../features/template-system.md) - Structured task creation
- 🤖 [AI Context Management](../features/ai-context-management.md) - AI assistant integration

### Advanced Setup

- 🔧 [MCP Configuration Examples](../../guides/mcp-config-examples.md) - All IDE configs
- 💻 [VS Code Extension](../user/vscode-extension.md) - IDE integration
- 🔍 [MCP Tools Reference](../../legacy/tr/mcp-araclari.md) - All 41 tools

### Troubleshooting

- 🆘 [Troubleshooting Guide](troubleshooting.md) - Common issues
- 📦 [Migration Guide](../../migration/v0.15-to-v0.16.md) - Upgrade from v0.15

---

## Quick Reference

### Template Aliases

| Alias | Use For |
|-------|---------|
| `bug` | Software bugs |
| `feature` | New features |
| `research` | Investigation |
| `refactor` | Code improvement |
| `test` | Test writing |
| `doc` | Documentation |

### CLI Commands

```bash
# Server
npx @mehmetsenol/gorev-mcp-server@latest          # Start server
npx @mehmetsenol/gorev-mcp-server --version      # Version info

# Workspace
npx @mehmetsenol/gorev-mcp-server init           # Initialize
npx @mehmetsenol/gorev-mcp-server workspace list # List workspaces

# Tasks
npx @mehmetsenol/gorev-mcp-server task list      # List tasks
npx @mehmetsenol/gorev-mcp-server task create    # Create task
npx @mehmetsenol/gorev-mcp-server task show <ID> # View details

# Projects
npx @mehmetsenol/gorev-mcp-server project list   # List projects
npx @mehmetsenol/gorev-mcp-server project create # Create project
```

### Default Ports

- **Web UI**: http://localhost:5082
- **MCP Protocol**: stdio (no network port)

### Environment Variables

```bash
GOREV_LANG=en              # Interface language (en/tr)
GOREV_API_PORT=5082        # Web UI port
GOREV_DB_PATH=/custom/path # Custom database location
```

---

## Tips & Tricks

### 1. Use Template Aliases

Instead of remembering UUIDs:

```bash
# ✅ Good
npx @mehmetsenol/gorev-mcp-server task create --template bug

# ❌ Tedious
npx @mehmetsenol/gorev-mcp-server task create --template 39f28dbd-10f3-454c-8b35-52ae6b7ea391
```

### 2. Workspace Auto-Detection

Don't specify `GOREV_DB_PATH` - let Gorev find `.gorev/` automatically:

```bash
cd /projects/project-a
npx @mehmetsenol/gorev-mcp-server task list  # Uses project-a database

cd /projects/project-b
npx @mehmetsenol/gorev-mcp-server task list  # Uses project-b database
```

### 3. AI Assistant Prompts

Use natural language:

```
"Create a bug task for login issue"
"Show me high priority tasks"
"Mark authentication task as in progress"
"What should I work on next?"
```

### 4. Web UI Keyboard Shortcuts

- `N`: New task
- `Ctrl+K`: Search tasks
- `Escape`: Close modals
- `Ctrl+S`: Save task (in edit mode)

---

## Getting Help

### Documentation

- **Main README**: [README.md](../../../README.md)
- **Full Guides**: [docs/guides/](../../guides/)
- **API Reference**: [MCP Tools](../../legacy/tr/mcp-araclari.md)

### Community

- **GitHub Issues**: https://github.com/msenol/gorev/issues
- **Discussions**: https://github.com/msenol/gorev/discussions

### Report Bugs

Open issue with:

- Gorev version (`npx @mehmetsenol/gorev-mcp-server --version`)
- Error messages
- Steps to reproduce
- Expected vs. actual behavior

---

**Congratulations! 🎉** You're ready to use Gorev for AI-powered task management!

For detailed documentation, see the [guides](../../guides/) directory.
