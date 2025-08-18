# CLAUDE.md

This file provides essential guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

**Last Updated:** 18 August 2025 | **Version:** v0.11.1

[🇺🇸 English](CLAUDE.en.md) | [🇹🇷 Türkçe](CLAUDE.md)

> 🤖 **Documentation Note**: This technical guide is optimized for token efficiency while maintaining essential information for AI assistants.

## 🚀 Recent Major Update

**v0.11.1 - Ultra-Detailed DRY Compliance & Thread-Safety (18 Aug 2025)**
- **Phase 7 Ultra-DRY Implementation**: Industry-leading DRY compliance achieved
  - **700+ total violations eliminated** across 7 comprehensive phases
  - **Template & Parameter Constants**: All hardcoded `"template_id"`, `"degerler"` strings replaced with constants
  - **Magic Number Elimination**: Replaced all hardcoded test numbers with context-specific constants
  - **Emoji Constants Enforcement**: All `"✅"` hardcoded emojis replaced with `constants.EmojiStatusCompleted`
  - **Zero DRY violations remaining**: Complete string duplication elimination
- **Thread-Safety Enhancement**: AI Context Manager race condition fix
  - Added `sync.RWMutex` protection to `AIContextYonetici` struct
  - Protected all context operations with read-write locks
  - 50-goroutine concurrent access test with race detector
- **Major Refactoring**: handlers.go architectural improvement
  - **File size reduction**: 3,060 lines → 2,362 lines (-698 lines, 23% reduction)
  - **New Architecture**: Extracted `tool_registry.go` (570 lines) and `tool_helpers.go` (286 lines)
  - **Eliminated Code Smells**: Replaced 703-line RegisterTools method with 4-line delegation
- **Rule 15 Compliance**: Zero technical debt, comprehensive DRY implementation

## 📋 Project Overview

**Gorev** is a two-module MCP (Model Context Protocol) server written in Go:

1. **gorev-mcpserver**: Core MCP server (Go) - Task management for AI assistants
2. **gorev-vscode**: Optional VS Code extension - Rich visual interface

**Core Features**: 25 MCP tools, unlimited subtask hierarchy, task dependencies, template system, bilingual support (TR/EN), AI context management.

## 🏗️ Architecture

```
cmd/gorev/main.go                  → Entry point, CLI commands (cobra)
internal/mcp/handlers.go           → MCP handlers (refactored, 2,362 lines)
internal/mcp/tool_registry.go      → MCP tool registration (570 lines)
internal/mcp/tool_helpers.go       → Validation & formatting utilities (286 lines)
internal/gorev/is_yonetici.go      → Business logic orchestration
internal/gorev/veri_yonetici.go    → Data access layer (SQLite)
internal/i18n/manager.go           → Internationalization system
internal/i18n/helpers.go           → DRY i18n patterns (NEW)
locales/[tr|en].json              → Translation files
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

# Development
make run                  # Build and run server
./gorev serve --debug     # Run with debug logging
./gorev serve --lang=en   # Run with English language
```

### VS Code Extension (gorev-vscode)
```bash
cd gorev-vscode
npm install               # Install dependencies  
npm run compile           # Compile TypeScript
npm test                  # Run tests
# Press F5 in VS Code to launch extension development host
```

## 🛠️ MCP Tools Summary

**25 MCP Tools** organized in 6 categories:
- **Task Management**: 7 tools (gorev_listele, gorev_detay, etc.)
- **Subtask Management**: 3 tools (gorev_altgorev_olustur, etc.) 
- **Templates**: 2 tools (template_listele, templateden_gorev_olustur)
- **Project Management**: 6 tools (proje_olustur, proje_listele, etc.)
- **AI Context**: 6 tools (gorev_set_active, gorev_nlp_query, etc.)
- **Reporting**: 1 tool (ozet_goster)

> **💡 Important**: Use `templateden_gorev_olustur` for task creation (gorev_olustur deprecated in v0.10.0)

## 🗄️ Database Schema

**SQLite database** with 9 tables + 1 view:
- **gorevler**: Tasks (with parent_id for hierarchy)
- **projeler**: Projects  
- **baglantilar**: Task dependencies
- **etiketler**, **gorev_etiketleri**: Tagging system
- **gorev_templateleri**: Task templates
- **ai_interactions**, **ai_context**: AI session management
- **aktif_proje**: Active project setting
- **gorev_hiyerarsi** (VIEW): Recursive hierarchy queries

Migrations: `gorev-mcpserver/internal/veri/migrations/` (handled by golang-migrate)

## 📝 Code Style

- **Domain Language**: Turkish terms for domain concepts (gorev=task, proje=project)
- **Technical Terms**: English for technical concepts and comments
- **Error Handling**: Always return explicit errors, use `mcp.NewToolResultError()`
- **Go Idioms**: Follow Go conventions, prefer composition over inheritance
- **i18n**: Use `i18n.T("key", templateData)` for user-facing strings

## 🧪 Testing Strategy

- **Unit Tests**: Business logic (`internal/gorev/`) - 81.3% coverage
- **DRY Test Patterns**: 12 comprehensive test files with reusable infrastructure
- **Table-Driven Tests**: TestCase structs, BenchmarkConfig, ConcurrencyTestConfig
- **Integration Tests**: MCP handlers (`test/integration_test.go`)  
- **VS Code Extension**: 100% test coverage with comprehensive mocks
- **Test Database**: Use `:memory:` SQLite for tests

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

- **MCP Tools Reference**: @docs/MCP_TOOLS_REFERENCE.md
- **Development History**: @docs/DEVELOPMENT_HISTORY.md  
- **Architecture Details**: Project structure above + clean architecture pattern
- **Database Migrations**: @internal/veri/migrations/
- **Testing Guide**: DRY patterns, table-driven tests, 81.3% server coverage, 100% extension coverage
- **Version Management**: Build-time injection via Makefile LDFLAGS

## 🚨 Rule 15: Comprehensive Problem-Solving & Zero Technical Debt

**ZERO TOLERANCE for shortcuts, workarounds, or temporary fixes**

### Core Principles:
1. **NO Workarounds**: Every problem requires root cause analysis and proper solution
2. **NO Code Duplication**: DRY (Don't Repeat Yourself) principle is absolute
3. **NO Technical Debt**: "Şimdilik böyle kalsın" approach is FORBIDDEN
4. **NO Quick Fixes**: Every solution must be production-ready
5. **NO Disabled Tests/Lints**: Fix test and lint errors instead of bypassing them

### Implementation Rules:
- ❌ `"temporary"`, `"workaround"`, `"quick fix"`, `"hotfix"`, `"band-aid"`
- ❌ `@ts-ignore`, `@ts-expect-error`, `eslint-disable`, `//nolint`
- ❌ Hardcoded values "for speed"
- ❌ Copy-paste solutions
- ✅ Root cause analysis
- ✅ Proper abstraction and reusability
- ✅ Comprehensive testing
- ✅ Clean, maintainable code

### For Gorev Project Specifically:
- **Template Enforcement**: All tasks MUST use templates (v0.10.0+)
- **Domain Terms**: Turkish domain terminology must be preserved
- **i18n Compliance**: Use `i18n.T()` for all user-facing strings
- **Test Coverage**: Maintain %81.3+ (server), %100 (extension)
- **Atomic Tasks**: Even atomic tasks must be comprehensive

### Quality Checklist:
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