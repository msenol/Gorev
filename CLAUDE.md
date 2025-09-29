# CLAUDE.md

This file provides essential guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

**Last Updated:** 29 September 2025 | **Version:** v0.16.0

[🇺🇸 English](CLAUDE.en.md) | [🇹🇷 Türkçe](CLAUDE.md)

## 🚀 Recent Major Update

**v0.16.0 - Embedded Web UI (29 Sep 2025)**

- **Embedded Web UI**: Modern React + TypeScript interface built into Go binary
- **Zero Configuration**: Automatically available at http://localhost:5082 with `npx gorev serve`
- **Enhanced Data Models**: Subtask and dependency visualization support
- **REST API Server**: Comprehensive Fiber-based API for web UI integration
- **Language Synchronization**: Web UI language switcher syncs with MCP server (TR/EN)
- **Backward Compatible**: All existing MCP and VS Code features maintained
- **Production Ready**: Vite-built React app served via Go embed.FS

## 📋 Project Overview

**Gorev** is a three-module task management project with MCP (Model Context Protocol) integration:

1. **gorev-mcpserver**: Core MCP server (Go) - Task management for AI assistants
   - Includes embedded Web UI (React + TypeScript) 🌐
   - REST API server (Fiber framework)
   - Automatically available at http://localhost:5082
2. **gorev-vscode**: Optional VS Code extension - Rich visual interface
3. **gorev-web**: React + TypeScript source code (for development only)

**Core Features**: 41 MCP tools, unlimited subtask hierarchy, task dependencies, template system, data export/import, IDE extension management, file watching, bilingual support (TR/EN), AI context management, enhanced NLP processing, advanced search & filtering, fuzzy matching, filter profiles.

## 🏗️ Architecture

```text
cmd/gorev/main.go                  → Entry point, CLI commands (cobra)
internal/mcp/handlers.go           → MCP handlers (refactored, 2,362 lines)
internal/mcp/tool_registry.go      → MCP tool registration (570 lines)
internal/mcp/tool_helpers.go       → Validation & formatting utilities (286 lines)
internal/api/                      → REST API server (Fiber framework)
  ├── server.go                   → HTTP server, request handlers & routes
  ├── static.go                   → Embedded web UI file serving
  └── middleware/                 → CORS, logging middleware
internal/testing/helpers.go        → Standardized test infrastructure
internal/gorev/is_yonetici.go      → Business logic orchestration
internal/gorev/veri_yonetici.go    → Data access layer (SQLite)
internal/gorev/modeller.go         → Enhanced data models (subtasks, dependencies)
internal/i18n/manager.go           → Internationalization system
internal/i18n/helpers.go           → DRY i18n patterns
locales/[tr|en].json              → Translation files
gorev-npm/                        → NPM package wrapper
  ├── package.json                → NPM package configuration
  ├── index.js                    → Platform detection & binary wrapper
  ├── postinstall.js              → Auto-download from GitHub releases
  └── bin/gorev-mcp               → Executable entry point
gorev-web/                        → React + TypeScript web UI
  ├── src/
  │   ├── components/             → React components (TaskCard, Sidebar, etc.)
  │   ├── api/client.ts           → API client (React Query)
  │   ├── types/index.ts          → TypeScript type definitions
  │   └── App.tsx                 → Main application component
  ├── public/                     → Static assets
  └── package.json                → Web UI dependencies
```

## 🔧 Development Commands

### MCP Server (gorev-mcpserver)

```bash
# Essential commands
make build                 # Build for current platform
make test                  # Run all tests with coverage
make fmt                   # Format code (run before commit)
go vet ./...              # Static analysis
make deps                  # Download and tidy dependencies

# Database initialization
./gorev init               # Initialize workspace database (.gorev/gorev.db)
./gorev init --global      # Initialize global database (~/.gorev/gorev.db)

# Development
make run                  # Build and run server
./gorev serve --debug     # Run with debug logging
./gorev serve --lang=en   # Run with English language
```

### NPM Package (gorev-npm)

```bash
cd gorev-npm
npm install               # Install dependencies (node-fetch)
npm test                  # Run tests (basic validation)
npm pack                  # Create package tarball for testing
npm publish --access public  # Publish to NPM registry
```

### VS Code Extension (gorev-vscode)

```bash
cd gorev-vscode
npm install               # Install dependencies
npm run compile           # Compile TypeScript
npm test                  # Run tests
vsce package              # Package extension as .vsix file
```

### Web UI (gorev-web)

> **📝 Not**: Web UI artık Go binary'sine embedded olarak gömülüdür. Ayrı kurulum ve çalıştırmaya gerek yoktur!

**Production Kullanım** (Otomatik):
```bash
cd gorev-mcpserver
./gorev serve --api-port 5082    # Web UI otomatik olarak http://localhost:5082 adresinde hazır
```

**Development** (Sadece web UI geliştirme için):
```bash
cd gorev-web
npm install               # Install dependencies
npm run dev               # Start development server (port 5001)
npm run build             # Build for production (output: ../gorev-mcpserver/cmd/gorev/web/dist)

# API server must be running (default port 5082)
cd gorev-mcpserver
./gorev serve --api-port 5082 --debug
```

## 🛠️ MCP Tools Summary

**41 MCP Tools** organized in 10 categories:

- **Task Management**: 6 tools (gorev_listele, gorev_detay, etc.)
- **Subtask Management**: 3 tools (gorev_altgorev_olustur, etc.)
- **Templates**: 2 tools (template_listele, templateden_gorev_olustur)
- **Project Management**: 6 tools (proje_olustur, proje_listele, etc.)
- **AI Context**: 6 tools (gorev_set_active, gorev_nlp_query, etc.)
- **Advanced Search & Filtering**: 6 tools (gorev_search_advanced, gorev_search_suggestions, gorev_search_history, gorev_filter_profile_*)
- **Data Export/Import**: 2 tools (gorev_export, gorev_import)
- **IDE Extension Management**: 5 tools (ide_detect, ide_install_extension, ide_uninstall_extension, ide_extension_status, ide_update_extension)
- **File Watching**: 4 tools (gorev_file_watch_add, etc.)
- **Advanced**: 2 tools (gorev_intelligent_create, ozet_goster)

> **💡 Template Aliases**: Use shortcuts like `bug`, `feature`, `research` with `templateden_gorev_olustur`

## 🗄️ Database Schema

**SQLite database** with 12 tables + 1 view:

- **gorevler**: Tasks (with parent_id for hierarchy)
- **projeler**: Projects
- **baglantilar**: Task dependencies
- **etiketler**, **gorev_etiketleri**: Tagging system
- **gorev_templateleri**: Task templates
- **ai_interactions**, **ai_context**: AI session management
- **aktif_proje**: Active project setting
- **gorevler_fts**: FTS5 virtual table for full-text search
- **filter_profiles**: Saved search filter combinations
- **search_history**: Search query history tracking
- **gorev_hiyerarsi** (VIEW): Recursive hierarchy queries

Migrations: `gorev-mcpserver/internal/veri/migrations/`

## 📝 Code Style

- **Domain Language**: Turkish terms for domain concepts (gorev=task, proje=project)
- **Technical Terms**: English for technical concepts and comments
- **Error Handling**: Always return explicit errors, use `mcp.NewToolResultError()`
- **Go Idioms**: Follow Go conventions, prefer composition over inheritance
- **i18n**: Use `i18n.T("key", templateData)` for user-facing strings

## 🧪 Testing Strategy

- **Standardized Test Infrastructure**: Centralized `internal/testing/helpers.go` package
- **Test Coverage**: 90%+ server coverage, 100% extension coverage
- **Unit Tests**: Business logic (`internal/gorev/`)
- **Integration Tests**: MCP handlers (`test/integration_test.go`)
- **Performance Testing**: Concurrent access, memory allocation, stress testing

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

## 📚 Essential References

- **MCP Tools Reference**: @docs/tr/mcp-araclari.md (Turkish documentation)
- **Development History**: @docs/development/TASKS.md
- **VS Code Data Export/Import**: @docs/guides/user/vscode-data-export-import.md
- **Architecture Details**: Project structure above + clean architecture pattern
- **Database Migrations**: @internal/veri/migrations/
- **Testing Guide**: DRY patterns, table-driven tests
- **Version Management**: Build-time injection via Makefile LDFLAGS

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
- **Test Coverage**: Maintain %90+ (server), %100 (extension)
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
