# MCP Configuration Examples

Production-ready MCP server configuration examples for various AI coding assistants.

> **Version**: v0.16.0
> **Last Updated**: 4 October 2025

---

## Prerequisites

Install Gorev MCP server via NPM:

```bash
npm install -g @mehmetsenol/gorev-mcp-server
```

Or use directly with `npx`:

```bash
npx @mehmetsenol/gorev-mcp-server --version
```

---

## Configuration by IDE

### Claude Code (Desktop App)

**Location**: `~/.claude/config.json` (macOS/Linux) or `%APPDATA%\.claude\config.json` (Windows)

```json
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

**With custom database path**:

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      "env": {
        "GOREV_LANG": "en",
        "GOREV_DB_PATH": "/path/to/custom/.gorev/gorev.db"
      }
    }
  }
}
```

---

### Kilo Code (VS Code Extension)

**Location**: `.kilocode/mcp.json` in your workspace root

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      "env": {
        "GOREV_LANG": "en"
      },
      "alwaysAllow": [
        "ozet_goster",
        "proje_listele",
        "gorev_listele",
        "template_listele",
        "gorev_search_advanced",
        "gorev_context_summary",
        "gorev_ide_status",
        "templateden_gorev_olustur",
        "gorev_detay",
        "gorev_guncelle",
        "gorev_altgorev_olustur",
        "gorev_hiyerarsi_goster",
        "gorev_bagimlilik_ekle",
        "gorev_file_watch_add",
        "gorev_file_watch_stats",
        "gorev_filter_profile_save",
        "gorev_filter_profile_list",
        "gorev_nlp_query",
        "gorev_batch_update",
        "gorev_ust_degistir",
        "gorev_export",
        "gorev_search_history",
        "gorev_set_active",
        "gorev_get_active",
        "gorev_duzenle",
        "proje_gorevleri",
        "gorev_file_watch_list",
        "gorev_filter_profile_load",
        "gorev_recent",
        "gorev_filter_profile_delete",
        "gorev_file_watch_remove",
        "gorev_ide_detect",
        "aktif_proje_ayarla",
        "aktif_proje_kaldir",
        "gorev_ide_uninstall",
        "gorev_sil"
      ]
    }
  }
}
```

---

### Cursor (AI Code Editor)

**Location**: `.cursor/mcp.json` in your workspace root

```json
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

---

### Windsurf (Codeium IDE)

**Location**: `.windsurf/mcp.json` in your workspace root

```json
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

---

## Configuration Options

### Environment Variables

| Variable | Description | Default | Values |
|----------|-------------|---------|--------|
| `GOREV_LANG` | Interface language | `tr` | `en`, `tr` |
| `GOREV_DB_PATH` | Custom database path | Auto-detected | Full path to `.gorev/gorev.db` |

### Language Options

**Turkish** (default):
```json
"env": {
  "GOREV_LANG": "tr"
}
```

**English**:
```json
"env": {
  "GOREV_LANG": "en"
}
```

### Workspace-Specific Database

**Auto-detection** (recommended):
- No `GOREV_DB_PATH` specified
- Gorev will look for `.gorev/gorev.db` in current workspace
- Falls back to global database at `~/.gorev/gorev.db`

**Custom path**:
```json
"env": {
  "GOREV_DB_PATH": "/path/to/project/.gorev/gorev.db"
}
```

---

## Initialization

### First-time Setup

1. **Initialize workspace**:
```bash
cd /path/to/your/project
npx @mehmetsenol/gorev-mcp-server init
```

2. **Create MCP config** (see examples above for your IDE)

3. **Restart your IDE/extension**

4. **Verify connection**:
   - Check MCP server status in your IDE
   - Should show "gorev" as connected

### Testing Connection

Use these commands in your AI assistant:

```
List all my projects
```

```
Show me available task templates
```

```
Create a new project called "Test Project"
```

---

## Available MCP Tools

Gorev provides **41 MCP tools** across 10 categories:

### Core Features
- **Task Management**: 6 tools
- **Subtask Management**: 3 tools
- **Templates**: 2 tools
- **Project Management**: 6 tools
- **AI Context Management**: 6 tools

### Advanced Features
- **Search & Filtering**: 6 tools (FTS5, fuzzy matching)
- **Data Export/Import**: 2 tools (JSON/CSV)
- **IDE Management**: 5 tools
- **File Watching**: 4 tools
- **Reporting**: 1 tool

**Full documentation**: See [MCP Tools Reference](../tr/mcp-araclari.md)

---

## Template Aliases

Quick shortcuts for common task templates:

| Alias | Template | Description |
|-------|----------|-------------|
| `bug` | Bug Raporu | Bug report with details |
| `feature` | Özellik Geliştirme | Feature development |
| `research` | Araştırma | Research task |
| `refactor` | Refactoring | Code refactoring |
| `test` | Test Yazma | Test writing |
| `doc` | Dokümantasyon | Documentation |

**Usage example**:
```
Create a bug task using the bug template for login issue
```

---

## Troubleshooting

### MCP Server Not Connecting

1. **Check NPM installation**:
```bash
npx @mehmetsenol/gorev-mcp-server --version
```

2. **Verify config syntax** (JSON must be valid)

3. **Check IDE logs** for error messages

4. **Restart IDE** after config changes

### Database Not Found

1. **Initialize workspace**:
```bash
npx @mehmetsenol/gorev-mcp-server init
```

2. **Check GOREV_DB_PATH** if using custom path

3. **Verify file permissions** on `.gorev/` directory

### Permission Errors

**macOS/Linux**:
```bash
chmod 755 ~/.gorev
chmod 644 ~/.gorev/gorev.db
```

**Windows**: Run terminal as Administrator

---

## Multi-Workspace Setup

### Option 1: Auto-Detection (Recommended)

Don't specify `GOREV_DB_PATH`. Gorev will automatically:
1. Look for `.gorev/gorev.db` in current workspace
2. Create new database if not found
3. Register workspace with Web UI

### Option 2: Explicit Path Per Workspace

Create workspace-specific config:

**Project A** (`.kilocode/mcp.json`):
```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      "env": {
        "GOREV_DB_PATH": "/projects/project-a/.gorev/gorev.db"
      }
    }
  }
}
```

**Project B** (`.kilocode/mcp.json`):
```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["-y", "@mehmetsenol/gorev-mcp-server", "serve"],
      "env": {
        "GOREV_DB_PATH": "/projects/project-b/.gorev/gorev.db"
      }
    }
  }
}
```

---

## Web UI Access

When MCP server is running, Web UI is available at:

```
http://localhost:5082
```

Features:
- Multi-workspace switcher
- Project and task management
- Visual task creation
- Export/import data
- Bilingual interface (TR/EN)

---

## Security Notes

1. **Local-only**: MCP server runs locally, no external connections
2. **File Permissions**: Ensure `.gorev/` directory has proper permissions
3. **Database Backup**: Regular backups recommended (use `gorev_export` tool)
4. **Sensitive Data**: Database is not encrypted, don't store secrets

---

## Additional Resources

- **MCP Tools Reference**: [docs/legacy/tr/mcp-araclari.md](../legacy/tr/mcp-araclari.md)
- **NPM Package**: [npmjs.com/package/@mehmetsenol/gorev-mcp-server](https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server)
- **GitHub Repository**: [github.com/msenol/gorev](https://github.com/msenol/gorev)
- **VS Code Extension Guide**: [docs/guides/user/vscode-data-export-import.md](../guides/user/vscode-data-export-import.md)

---

**Need help?** Open an issue at [GitHub Issues](https://github.com/msenol/gorev/issues)
