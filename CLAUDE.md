# CLAUDE.md

This file provides essential guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

**Last Updated:** Auto-condensed on $(date '+%d %B %Y') | **Version:** v0.15.24

---

**Note:** This file has been auto-condensed to maintain optimal size. Critical rules and essential information are preserved.

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


## 📚 Essential References

- **MCP Tools Reference**: @docs/tr/mcp-araclari.md (Turkish documentation)
- **Architecture Details**: Project structure above + clean architecture pattern
- **Database Migrations**: @internal/veri/migrations/ (including FTS5 migration 000010)
- **Testing Guide**: DRY patterns, table-driven tests, 90%+ server coverage, 100% extension coverage
- **Version Management**: Build-time injection via Makefile LDFLAGS
- **Advanced Search Documentation**: SearchEngine, FilterProfileManager, NLP integration

---

> 💡 **Token Optimization**: Detailed information moved to `docs/` folder. This file contains only essential guidance for daily development work.
