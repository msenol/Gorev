# CLAUDE.md

This file provides essential guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

**Last Updated:** October 6, 2025 | **Version:** v0.16.3

[🇺🇸 English](CLAUDE.en.md) | [🇹🇷 Türkçe](CLAUDE.md)

## 🚀 Recent Major Update

**v0.16.3 - MCP Tool Parameter Transformation Fixes (October 6, 2025)**

- **gorev_bulk**: All 3 operations (update/transition/tag) now working with proper parameter transformation
- **gorev_guncelle**: Extended to support both `durum` and `oncelik` updates simultaneously
- **gorev_search**: Advanced mode now supports query parsing (`durum:X oncelik:Y` → filters)
- **VS Code Tree View**: Dependency counters (🔒/🔓/🔗) now display correctly (JSON `omitempty` fix)
- **Backward Compatibility**: Unified handlers accept multiple parameter formats for flexibility

**Previous (v0.16.2 - October 5, 2025):**

- NPM Binary Update Fix, VS Code Auto-Start, Embedded Web UI, Multi-Workspace, Template Aliases

## 📋 Project Overview

**Gorev** is a three-module task management project with MCP (Model Context Protocol) integration:

1. **gorev-mcpserver**: Core MCP server (Go) - Task management for AI assistants
   - Includes embedded Web UI (React + TypeScript) 🌐
   - REST API server (Fiber framework)
   - Automatically available at http://localhost:5082
2. **gorev-vscode**: Optional VS Code extension - Rich visual interface
3. **gorev-web**: React + TypeScript source code (for development only)

**Core Features**: 24 optimized MCP tools (unified from 45), unlimited subtask hierarchy, task dependencies, template system, data export/import, IDE extension management, file watching, bilingual support (TR/EN), AI context management, enhanced NLP processing, advanced search & filtering, fuzzy matching, filter profiles.

## 🏗️ Architecture

**Core Layers**:

- `cmd/gorev/` - Entry point, CLI commands, daemon management
- `internal/mcp/` - MCP protocol layer (handlers, tools, helpers)
- `internal/api/` - REST API server (Fiber), embedded Web UI
- `internal/gorev/` - Business logic, data access (SQLite)
- `internal/daemon/` - Lock file, health checks, process management
- `internal/websocket/` - Real-time update broadcasts
- `internal/i18n/` - Internationalization (TR/EN)
- `gorev-npm/` - NPM package with auto-download
- `gorev-web/` - React + TypeScript UI source (embedded in binary)

### 🔌 Daemon Architecture (v0.16.0+)

Gorev runs as a **background daemon process** with multi-client support:

- **Lock File Mechanism**: `~/.gorev-daemon/.lock` ensures single instance, provides service discovery
- **Multi-Client MCP Proxy**: Multiple AI assistants (Claude, Windsurf, Cursor) can connect simultaneously
- **REST API Server**: 23 endpoints for VS Code extension (Fiber framework, port 5082)
- **WebSocket Server**: Real-time task update broadcasts (experimental)
- **VS Code Auto-Start**: Extension automatically detects and starts daemon
- **Workspace Isolation**: SHA256-based workspace IDs for multi-project support

**Key Files:**

- `cmd/gorev/daemon.go` - Daemon lifecycle management
- `cmd/gorev/mcp_proxy.go` - Multi-client MCP routing
- `internal/daemon/lockfile.go` - Single instance enforcement
- `internal/api/mcp_bridge.go` - MCP-to-REST API bridge

See [Daemon Architecture Documentation](docs/architecture/daemon-architecture.md) for detailed technical specifications.

## 🔧 Development Commands

```bash
# Build & Run
make build                # Build for current platform
make test                 # Run all tests with coverage
./gorev serve --debug     # Run with debug logging

# Database
./gorev init              # Initialize workspace DB (.gorev/gorev.db)

# VS Code Extension
cd gorev-vscode
npm install && npm run compile

# Web UI (embedded in binary, auto-available at :5082)
./gorev serve --api-port 5082
```

## 🛠️ MCP Tools Summary

**24 Optimized MCP Tools** (reduced from 45 via unification):

- **Core Tools (11)**: Task CRUD (5), Templates (2), Projects (3), Dependencies (1)
- **Unified Tools (8)**: Active Project, Hierarchy, Bulk Ops, Filter Profiles, File Watch, IDE Management, AI Context, Search
- **Special Tools (5)**: Summary, Export, Import, AI Suggestions, Intelligent Create

> **Template Aliases**: `bug`, `feature`, `research`, `refactor`, `test`, `doc`

## 🗄️ Database Schema

**12 tables + 1 view**: gorevler (tasks), projeler, baglantilar (dependencies), etiketler, gorev_templateleri, ai_interactions, ai_context, aktif_proje, gorevler_fts (full-text search), filter_profiles, search_history, gorev_hiyerarsi (VIEW)

## 📝 Code Style

- **Domain Language**: Turkish terms for domain concepts (gorev=task, proje=project)
- **Technical Terms**: English for technical concepts and comments
- **Error Handling**: Always return explicit errors, use `mcp.NewToolResultError()`
- **Go Idioms**: Follow Go conventions, prefer composition over inheritance
- **i18n**: Use `i18n.T("key", templateData)` for user-facing strings

## 🧪 Testing Strategy

Centralized test infrastructure with ~71% server coverage (goal: 80%+), 100% extension coverage. Unit tests for business logic, integration tests for MCP handlers, performance testing for concurrent access.

## 🔄 Adding New MCP Tools

1. Add handler method to `internal/mcp/handlers.go`
2. Register tool in `RegisterTools()` with proper schema
3. Add integration tests in `test/integration_test.go`
4. Update `docs/MCP_TOOLS_REFERENCE.md` with tool documentation

## 🌍 Language Support

**Environment Setup:**

```bash
export GOREV_LANG=en     # English
export GOREV_LANG=tr     # Turkish (default)
```

**CLI Usage:**

```bash
gorev serve --lang=en    # English interface
gorev serve --lang=tr    # Turkish interface
```


## 🚨 Rule 15: Comprehensive Problem-Solving & Zero Technical Debt

**ZERO TOLERANCE for shortcuts, workarounds, or temporary fixes**

### Core Principles

1. **NO Workarounds**: Every problem requires root cause analysis and proper solution
2. **NO Code Duplication**: DRY (Don't Repeat Yourself) principle is absolute
3. **NO Technical Debt**: "Şimdilik böyle kalsın" approach is FORBIDDEN
4. **NO Quick Fixes**: Every solution must be production-ready
5. **NO Disabled Tests/Lints**: Fix test and lint errors instead of bypassing them

### Implementation Rules

- ❌ `"temporary"`, `"workaround"`, `"quick fix"`, `"hotfix"`, `"band-aid"`
- ❌ `@ts-ignore`, `@ts-expect-error`, `eslint-disable`, `//nolint`
- ❌ Hardcoded values "for speed"
- ❌ Copy-paste solutions
- ✅ Root cause analysis
- ✅ Proper abstraction and reusability
- ✅ Comprehensive testing
- ✅ Clean, maintainable code

### For Gorev Project Specifically

- **Template Enforcement**: All tasks MUST use templates (v0.10.0+)
- **Domain Terms**: Turkish domain terminology must be preserved
- **i18n Compliance**: Use `i18n.T()` for all user-facing strings
- **Test Coverage**: Maintain ~71% server (goal: 80%+), 100% extension
- **Atomic Tasks**: Even atomic tasks must be comprehensive

### Quality Checklist

- [ ] Root cause identified and addressed
- [ ] No temporary workarounds introduced
- [ ] DRY principle followed (no duplication)
- [ ] All tests passing without disabling
- [ ] Proper error handling with context
- [ ] i18n keys used for user messages

## 🚨 Important Development Rules

1. **NEVER commit**: `*.db`, `*.log`, binary files (`gorev`, `gorev.exe`)
2. **Always run before commit**: `make fmt`, `go vet ./...`, `make deps`, `make test`
3. **Template Usage**: Mandatory since v0.10.0, use `templateden_gorev_olustur`
4. **Turkish Domain**: Keep domain concepts in Turkish, technical terms in English
5. **Error Context**: Wrap errors with context: `fmt.Errorf("context: %w", err)`
6. **i18n Strings**: Use `i18n.T()` for all user-facing messages
7. **Rule 15 Compliance**: NO workarounds, NO technical debt, NO quick fixes

---

> 💡 **Token Optimization**: Detailed information moved to `docs/` folder. This file contains only essential guidance for daily development work.
