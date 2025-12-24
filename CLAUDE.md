# CLAUDE.md

This file provides essential guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

**Last Updated:** December 24, 2025 | **Version:** v0.17.0

[🇺🇸 English](CLAUDE.en.md) | [🇹🇷 Türkçe](CLAUDE.md)

## 🚀 Recent Update (December 24, 2025) - Smart Shutdown & Client Tracking

**⚠️ BREAKING CHANGES**

- **Smart Shutdown**: VS Code extension now checks active clients before stopping daemon
  - Prevents killing MCP proxy (kilocode) when VS Code closes
  - Client count tracking via `/api/v1/daemon/clients/count` endpoint
  - Client registration and heartbeat system with 5-minute TTL

- **Client Tracking System** (NEW):
  - `internal/daemon/client_tracker.go` - Thread-safe client tracking
  - `internal/config/shared_config.go` - Shared configuration for all components
  - MCP Proxy registers as client on startup, sends heartbeat every 60s
  - Extension smart shutdown: if other clients exist, daemon keeps running

- **Unified MCP Tools (December 23, 2025)**: 27 individual tools consolidated into 8 unified handlers
  - Tool Count: 41 tools → 26 tools (37% reduction)
  - Unified Tools:
    - `aktif_proje` (actions: set|get|clear) - replaces 3 active project tools
    - `gorev_bulk` (operations: transition|tag|update) - replaces 3 bulk operation tools
    - `gorev_hierarchy` (actions: create_subtask|change_parent|show) - replaces 3 hierarchy tools
    - `gorev_filter_profile` (actions: save|load|list|delete) - replaces 4 filter profile tools
    - `gorev_ide` (actions: detect|install|uninstall|status|update) - replaces 5 IDE tools
    - `gorev_context` (actions: set_active|get_active|recent|summary) - replaces 4 AI context tools
    - `gorev_search` (modes: nlp|advanced|history) - replaces 3 search tools
    - `gorev_ozet` (actions: show) - summary tool
  - **Action Constants**: Added `constants.ActionXxx` to eliminate magic strings

- **Enhanced Summary Dashboard**: New HTML dashboard with:
  - Progress bars for task status
  - Due date summary with overdue/today/this week counts
  - High priority and recent tasks widgets
  - Active project display

**Previous Updates:**

**v0.17.0 - English Field Names Migration (October 11, 2025)** ⚠️ **BREAKING CHANGES**

- **Database Schema**: All column names migrated from Turkish to English (automatic migration 000011)
- **Go Backend**: 55+ files updated - all struct JSON tags now use English field names
- **TypeScript Frontend**: 20 files updated - API types and components use English fields
- **VS Code Extension**: Server auto-start improvements - reliable health checks, comprehensive error handling, 60s timeout
- **Template Placeholders**: `{{baslik}}` → `{{title}}`, `{{aciklama}}` → `{{description}}`
- **Backward Compatibility**: Domain terms (`gorevler`, `projeler`) remain Turkish
- See `docs/MIGRATION_GUIDE_v0.17.md` for upgrade instructions

**Recent Updates (December 4, 2025) - v0.17.0:**

- **MCP Tool Parameter Migration Complete**: All MCP tool parameters now use English names
  - Schema definitions in `tool_registry.go` and `mcp_bridge.go` fully migrated
  - Handler parameter extraction using `constants.ParamXxx` consistently
  - Parameters: `name`, `definition`, `status`, `confirm`, `category`, `tag`, `order_by`, `filter`, `all_projects`
  - All test files updated to use English parameter names
  - Files: `internal/mcp/handlers.go`, `internal/mcp/tool_registry.go`, `internal/api/mcp_bridge.go`

**Previous (December 2, 2025):**

- **Subtasks API & Tree View Fix**: Complete subtask support in VS Code extension
  - New `GET /api/v1/tasks/:id/subtasks` endpoint for retrieving task subtasks
  - `getTask` handler now includes nested subtasks in response
  - Added `subtasks` → `alt_gorevler` mapping in VS Code API client
  - Tree view now correctly shows expandable tasks with subtask hierarchies
  - Files: `internal/api/server.go`, `gorev-vscode/src/api/client.ts`

- **E2E Testing Infrastructure**: Comprehensive Playwright-based testing framework
  - Test data seeding with `./gorev seed-test-data` CLI command
  - Page Object pattern for maintainable tests
  - Interactive test runner script (`./test-runner.sh`)
  - Manual test guide with 79 scenarios
  - Files: `cmd/gorev/seed.go`, `internal/testing/seeder.go`, `test-runner.sh`

- **Web UI Test Attributes**: Added `data-testid` for E2E testing
  - TaskList, TaskCard, Header components enhanced
  - Files: `gorev-web/src/components/*.tsx`

- **Critical Bug Fixes - Workspace Isolation**: Fixed workspace_id not being set on task creation
- **API Field Mapping Fix**: Fixed English→Turkish field mapping for v0.17.0 API
- **Heavy Development Warnings**: Added prominent warnings across all platforms

**Previous (November 22, 2025):**

- **Multilingual Template Support** - Templates now support multiple languages (TR/EN)
  - Database schema extended with `language_code` and `base_template_id` fields (migration 000012)
  - Template pairs created for common templates (bug, feature) in Turkish and English
  - Language-aware template selection based on `GOREV_LANG` environment variable
  - Files: `internal/gorev/template_yonetici.go`, `internal/veri/migrations/000012_add_template_multilang.*.sql`

- **i18n Phase 3 Complete** - Context-aware language propagation system
  - Environment variable `GOREV_LANG` now properly propagates through all layers (CLI, MCP, API)
  - Per-request language selection for MCP handlers with fallback hierarchy
  - Files: `internal/i18n/manager.go`, `internal/i18n/helpers.go`, `internal/mcp/handlers.go`

- **VS Code Extension: Rule 15 Compliance Achieved**
  - Eliminated all 242 ESLint warnings → 0 warnings (100% clean)
  - Maintained 100% test pass rate (104/104 tests)
  - Files: 13 files refactored including `ui/taskDetailPanel.ts`, `providers/*.ts`, `commands/*.ts`

**Previous (v0.16.3 - October 6, 2025):**

- MCP Tool Parameter Transformation Fixes, Bulk Operations, VS Code Dependency Counters

## 📋 Project Overview

**Gorev** is a three-module task management project with MCP (Model Context Protocol) integration:

1. **gorev-mcpserver**: Core MCP server (Go) - Task management for AI assistants
   - Includes embedded Web UI (React + TypeScript) 🌐
   - REST API server (Fiber framework)
   - Automatically available at http://localhost:5082
2. **gorev-vscode**: Optional VS Code extension - Rich visual interface
3. **gorev-web**: React + TypeScript source code (for development only)

**Core Features**: 26 optimized MCP tools (7 unified handlers), unlimited subtask hierarchy, task dependencies, template system with aliases (bug, feature, etc.), data export/import, IDE extension management, file watching, bilingual support (TR/EN), AI context management, enhanced NLP processing, advanced search & filtering, fuzzy matching, filter profiles.

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
- **REST API Server**: 24 endpoints for VS Code extension (Fiber framework, port 5082)
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

### Build & Run

```bash
# Build Process (Important: Web UI is built first, then embedded in Go binary)
cd gorev-mcpserver
make build                             # Builds Web UI first, then Go binary
make build-all                         # Cross-platform builds (Linux/macOS/Windows)

# Run server
./gorev serve                          # Normal mode
./gorev serve --debug                  # Debug mode with verbose logging
./gorev serve --port 5082              # Custom port
./gorev daemon --detach                # Start daemon in background (recommended for MCP)
```

### Testing

```bash
# Run all tests
make test                              # Root: runs both server and extension tests
cd gorev-mcpserver && make test        # Server tests only (~71% coverage)
cd gorev-vscode && npm test            # Extension tests (100% coverage)

# Specific test commands
cd gorev-mcpserver
go test -v ./internal/mcp/             # Test specific package
go test -v -run TestGorevOlustur ./... # Run single test by name
go test -v -race ./...                 # Race condition detection
make test-coverage                     # Generate coverage report (coverage.html)

# VS Code Extension Testing
cd gorev-vscode
npm test                               # Run extension tests
npm run test:coverage                  # Extension test coverage

# E2E Testing (Playwright)
cd gorev-vscode
npm run test:e2e                       # Run all E2E tests
npm run test:e2e:api                   # API integration tests only
npm run test:e2e:ui                    # Web UI tests only (requires web server)

# Test Data Seeding
./gorev seed-test-data                 # Seed with full Turkish data
./gorev seed-test-data --lang=en       # Seed with English data
./gorev seed-test-data --minimal       # Quick seed with 3 tasks
./gorev seed-test-data --force         # Overwrite existing data
```

### Web UI Development

```bash
# Develop Web UI independently (gorev-web)
cd gorev-web
npm install                            # Install dependencies
npm run dev                            # Start Vite dev server (port 5173)
npm run build                          # Build production bundle
npm run preview                        # Preview production build

# Note: Web UI is automatically embedded in Go binary via build-web target
# The production build goes to gorev-mcpserver/binaries/web-ui/
```

### VS Code Extension Development

```bash
cd gorev-vscode
npm install                            # Install dependencies
npm run compile                        # Compile TypeScript
npm run watch                          # Watch mode (for development)

# Testing in VS Code
# 1. Open gorev-vscode folder in VS Code
# 2. Press F5 (Run > Start Debugging)
# 3. Test in Extension Development Host window
```

#### Connection Mode Configuration

The VS Code extension supports **4 connection modes** via `gorev.connectionMode` setting:

| Mode | Behavior | Configuration |
|------|----------|---------------|
| **auto** (default) | Automatically detects and starts daemon using npm package | No additional config needed |
| **local** | Uses local gorev binary (for development) | Set `gorev.serverPath` to binary path |
| **docker** | Uses docker-compose to run container | Set `gorev.dockerComposeFile` path |
| **remote** | Connects to existing remote server | Server must be started manually |

**Local Development Setup** (recommended for contributors):
```json
{
  "gorev.connectionMode": "local",
  "gorev.serverPath": "/home/msenol/Projects/Gorev/gorev-mcpserver/gorev"
}
```

Build the local binary first:
```bash
cd gorev-mcpserver
make build
```

Now VS Code extension will use your locally built binary instead of the npm package.

### Database

```bash
./gorev init                           # Initialize workspace database (.gorev/gorev.db)
# Migrations run automatically on first init
# Migration files: gorev-mcpserver/internal/veri/migrations/*.sql
```

### Debugging & Development

```bash
# Daemon management
./gorev daemon --detach                # Start daemon in background
./gorev daemon-status                  # Check daemon status
./gorev daemon-stop                    # Stop running daemon
curl http://localhost:5082/api/health  # Health check endpoint

# Clean build artifacts
make clean                             # Root level: cleans both modules
cd gorev-mcpserver && make clean       # Server only
cd gorev-vscode && rm -rf out/         # Extension only
```

## 🛠️ MCP Tools Summary

**26 Optimized MCP Tools** (reduced from 41 via unification - 37% reduction):

- **Core Tools (10)**: Task CRUD (5), Templates (2), Projects (2), Dependencies (1)
- **Unified Tools (8)**:
  - `aktif_proje` (actions: set|get|clear) - Active project management
  - `gorev_bulk` (operations: transition|tag|update) - Bulk task operations
  - `gorev_hierarchy` (actions: create_subtask|change_parent|show) - Task hierarchy
  - `gorev_filter_profile` (actions: save|load|list|delete) - Filter profiles
  - `gorev_ide` (actions: detect|install|uninstall|status|update) - IDE extension
  - `gorev_context` (actions: set_active|get_active|recent|summary) - AI context
  - `gorev_search` (modes: nlp|advanced|history) - Search
  - `gorev_ozet` (actions: show) - Summary dashboard
- **File Watcher Tools (4)**: Add, Remove, List, Stats
- **Special Tools (4)**: Export, Import, AI Suggestions

> **Template Aliases**: `bug`, `feature`, `research`, `refactor`, `test`, `doc`

## 🗄️ Database Schema

**12 tables + 1 view**: gorevler (tasks), projeler, baglantilar (dependencies), etiketler, gorev_templateleri, ai_interactions, ai_context, aktif_proje, gorevler_fts (full-text search), filter_profiles, search_history, gorev_hiyerarsi (VIEW)

## 📝 Code Style

- **Domain Language**: Turkish terms for domain concepts (gorev=task, proje=project)
- **Technical Terms**: English for technical concepts and comments
- **Error Handling**: Always return explicit errors, use `mcp.NewToolResultError()`
- **Go Idioms**: Follow Go conventions, prefer composition over inheritance
- **i18n**: Use `i18n.T("key", templateData)` for user-facing strings

## 🗂️ Quick File Reference

**Need to modify...**

- **MCP Tools**: `internal/mcp/handlers.go` + register in `tool_registry.go`
- **Client Tracking**: `internal/daemon/client_tracker.go`, `internal/config/shared_config.go`
- **Business Logic**: `internal/gorev/is_yonetici.go`
- **Database Access**: `internal/veri/veri_yonetici.go`
- **Database Schema**: `internal/veri/migrations/*.sql` (add new migration)
- **REST API Endpoints**: `internal/api/server.go` (28 endpoints)
- **i18n Strings**: `locales/en.toml`, `locales/tr.toml`
- **CLI Commands**: `cmd/gorev/*.go` (daemon.go, serve.go, etc.)
- **Daemon Logic**: `cmd/gorev/daemon.go`, `internal/daemon/lockfile.go`
- **MCP Proxy**: `internal/mcp/proxy.go` (client registration, heartbeat)
- **VS Code Extension**: `gorev-vscode/src/extension.ts`
- **Extension Manager**: `gorev-vscode/src/managers/unifiedServerManager.ts` (smart shutdown)
- **Extension TreeView**: `gorev-vscode/src/providers/*.ts`
- **Web UI Components**: `gorev-web/src/components/*.tsx`
- **Web UI API Client**: `gorev-web/src/api/client.ts`
- **Test Data Seeding**: `cmd/gorev/seed.go`, `internal/testing/seeder.go`
- **E2E Page Objects**: `gorev-vscode/test/integration/playwright/page-objects/*.ts`
- **E2E Test Specs**: `gorev-vscode/test/integration/playwright/e2e/*.spec.ts`

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
# Daemon mode (recommended for MCP usage)
gorev daemon --detach --lang=en    # English interface
gorev daemon --detach --lang=tr    # Turkish interface

# Direct server mode (for development/debugging)
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

## 🚨 Pre-Commit Checklist

**Before committing, ALWAYS run these commands in order:**

```bash
# 1. Format code
make fmt                               # Formats both Go and TypeScript

# 2. Update dependencies (if go.mod or package.json changed)
make deps

# 3. Run linters
make lint                              # Both Go (golangci-lint) and TS linters
go vet ./...                           # Additional Go static analysis

# 4. Run all tests (MUST pass 100%)
make test                              # Both server and extension tests

# 5. If tests pass, commit
git add .
git commit -m "feat: your message"
```

**NEVER commit:**

- Database files: `*.db`, `*.db-shm`, `*.db-wal`
- Log files: `*.log`
- Binary files: `gorev`, `gorev.exe`, `gorev-linux`, `gorev-darwin`, `gorev-windows`
- Build artifacts: `node_modules/`, `out/`, `dist/`, `coverage.out`, `coverage.html`
- Lock files: `~/.gorev-daemon/.lock`
- Temporary files: `.DS_Store`, `Thumbs.db`

**Development Standards:**

1. **Template Usage**: Mandatory since v0.10.0, use `templateden_gorev_olustur` tool
2. **Turkish Domain**: Keep domain concepts in Turkish (gorev, proje, durum), technical terms in English
3. **Error Context**: Always wrap errors with context: `fmt.Errorf("context: %w", err)`
4. **i18n Strings**: Use `i18n.T("key", templateData)` for all user-facing messages
5. **Rule 15 Compliance**: NO workarounds, NO technical debt, NO quick fixes
6. **Test Coverage**: All new code must have tests (maintain ~71% server, 100% extension)

**Commit Message Format:**

```
<type>(<scope>): <subject>

Examples:
feat(mcp): add new search filter tool
fix(api): resolve race condition in workspace manager
docs(readme): update installation instructions
test(handlers): add edge case tests for bulk operations
```

---

> 💡 **Token Optimization**: Detailed information moved to `docs/` folder. This file contains only essential guidance for daily development work.
