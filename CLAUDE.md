# CLAUDE.md

This file provides guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

## Last Updated: 30 June 2025

> 🤖 **Documentation Note**: This comprehensive technical guide was enhanced and structured with the assistance of Claude (Anthropic), demonstrating the power of AI-assisted documentation in modern software development.

### Recent Changes (v0.8.0)

#### Major Features (30 June 2025)
- **Implemented Subtask System with Unlimited Hierarchy**:
  - Added `parent_id` column to tasks table with foreign key constraint
  - Created recursive CTE views for efficient hierarchy queries
  - Implemented circular dependency prevention
  - Added parent task progress tracking based on subtask completion
  - **New MCP Tools**:
    - `gorev_altgorev_olustur` - Create subtask under a parent task
    - `gorev_ust_degistir` - Move task to different parent or root
    - `gorev_hiyerarsi_goster` - Show complete task hierarchy with statistics
  - **Business Rules**:
    - Tasks cannot be deleted if they have subtasks
    - Tasks cannot be completed unless all subtasks are completed
    - Moving a task to a different project moves all its subtasks
    - Subtasks inherit parent's project
  - **Data Layer Methods**:
    - `AltGorevleriGetir` - Get direct subtasks
    - `TumAltGorevleriGetir` - Get entire subtask tree recursively
    - `UstGorevleriGetir` - Get parent hierarchy
    - `GorevHiyerarsiGetir` - Get hierarchy statistics
    - `DaireBagimliligiKontrolEt` - Prevent circular dependencies
  - **UI Enhancements**:
    - Hierarchical task display in `gorev_listele` with tree structure
    - Progress indicators showing subtask completion percentage
    - Visual hierarchy with indentation and tree connectors

### Previous Changes (v0.7.1)

#### Bug Fixes (30 June 2025)
- **Fixed Filter State Persistence Issue** in VS Code extension:
  - Added `clearFilters()` method to `EnhancedGorevTreeProvider`
  - Fixed `clearAllFilters()` in FilterToolbar to properly reset state
  - Added keyboard shortcut `Ctrl+Alt+R` / `Cmd+Alt+R` for quick filter clearing
  - Users can now clear filters without restarting VS Code

### Previous Changes (v0.7.0-beta.1)

#### Test Infrastructure Improvements (30 June 2025)
- **MCP Server Test Coverage**:
  - Improved overall MCP package coverage from 75.1% to 81.5% (+6.4%)
  - Created `handlers_edge_cases_test.go` (600+ LOC) with comprehensive edge case testing
  - Created `template_yonetici_test.go` (400+ LOC) for template unit tests
  - Enhanced `handlers_test.go` with complete template handler coverage
  - Fixed database migration issues with `etiketler` table in tests
  - Fixed concurrent access test using file-based database instead of in-memory
  - Discovered and documented validation gaps for future improvements
- **VS Code Extension Test Coverage**:
  - Achieved 50.9% file coverage (up from 0%) with 19 files tested
  - Added 7 new unit test files totaling 2,700 LOC
  - Created custom test coverage analysis tool (`test-coverage.js`)
  - Key test files added:
    - `enhancedGorevTreeProvider.test.js` (389 LOC) - TreeView functionality
    - `taskDetailPanel.test.js` (396 LOC) - WebView panel testing
    - `logger.test.js` (237 LOC) - Logging utility tests
    - `models.test.js` (273 LOC) - TypeScript model validation
    - `utils.test.js` (307 LOC) - Utility function tests
- **Bug Fixes in Tests**:
  - Fixed TypeScript compilation error in `markdownParser.ts`
  - Added table existence check in `gorevEtiketleriniGetir` to handle missing tables
  - Added `npm run coverage` script for test coverage reporting
- **Testing Framework Decision**:
  - Evaluated testify vs ginkgo for Go testing
  - Decided to continue with testify (152x faster, already integrated)
  - Created `docs/testing-framework-decision.md` documenting the rationale

#### VS Code Extension - Enhanced UI Features:
  - Enhanced TreeView with grouping, multi-select, and priority-based color coding
  - Drag & Drop support for moving tasks, changing status, and creating dependencies
  - Inline editing with F2/double-click, context menus, and date pickers
  - Advanced filtering toolbar with search, filters, and saved profiles
  - Rich task detail panel with markdown editor and dependency visualization
  - Template wizard UI with multi-step interface and dynamic forms
  - Comprehensive test suite (unit, integration, E2E) with coverage reporting
- **MCP Server Improvements**:
  - Fixed path resolution for database and migrations to work from any directory
  - Added `getDatabasePath()` and `getMigrationsPath()` functions
  - Enhanced `GorevListele` and `ProjeGorevleri` handlers to include tags and due dates
- **Bug Fixes** (29 June 2025):
  - Fixed tag display in VS Code UI when tasks created via CLI
  - Fixed project task count showing as 0 in TreeView
  - Fixed task detail panel UI issues in dark theme:
    - Action buttons now visible with proper styling
    - Markdown editor toolbar displays correctly
    - CSP-compliant event handlers
    - Edit/Delete functionality restored
  - Fixed single-click task selection in TreeView
  - Removed non-functional dependency graph feature
  - Added Filter State Persistence Issue to ROADMAP.md as Task #7
- **Documentation Updates**:
  - Added comprehensive TASKS.md with 11 active development tasks
  - Documented AI-Powered task enrichment system plans
  - Updated version to 0.7.0-beta.1 for beta release

### Previous Changes (v0.5.0 - v0.6.0)
- **Added Task Template System** - Predefined templates for bug reports, feature requests, technical debt, and research tasks
- **Added Task Dependencies** - Tasks can now have dependencies that must be completed before starting
- **Added Due Dates** - Tasks can have deadlines with filtering for urgent/overdue tasks
- **Added Tagging System** - Tasks can be categorized with multiple tags
- **Database Schema Management** - Using golang-migrate for version control
- **Enhanced gorev_listele** - Added sorting (sirala) and filtering (filtre, etiket) parameters
- **Enhanced gorev_olustur** - Now accepts son_tarih (due date) and etiketler (tags) parameters
- **Enhanced gorev_detay** - Shows dependencies with completion status indicators
- **New MCP tools**: 
  - `gorev_bagimlilik_ekle` for creating task dependencies
  - `template_listele` for listing available task templates
  - `templateden_gorev_olustur` for creating tasks from templates
- **New CLI commands**:
  - `gorev template list [kategori]` - List templates by category
  - `gorev template show <template-id>` - Show template details
  - `gorev template init` - Initialize default templates
- **Breaking changes**: 
  - GorevOlustur now takes 6 parameters (added sonTarihStr, etiketIsimleri)
  - GorevListele now takes 3 parameters (added sirala, filtre)
  - VeriYonetici constructor requires migrations path

## Project Overview

Gorev is a two-module project that provides task management capabilities to AI assistants:

1. **gorev-mcpserver**: An MCP (Model Context Protocol) server written in Go that provides task management capabilities to AI assistants across all MCP-compatible editors (Claude Desktop, VS Code, Windsurf, Cursor, Zed, etc.). Uses the community MCP SDK (`mark3labs/mcp-go`).

2. **gorev-vscode**: A VS Code extension (optional) that provides a rich visual interface for task management. It connects to the MCP server and offers TreeView panels, status bar integration, and command palette commands.

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

### gorev-mcpserver (Go)
```
cmd/gorev/main.go                  → Entry point, CLI commands (cobra)
internal/mcp/                      → MCP protocol layer
  ├── handlers.go                 → MCP tool implementations
  └── server.go                  → MCP server setup
internal/gorev/                   → Business logic layer
  ├── modeller.go                → Domain models (Gorev, Proje, Ozet)
  ├── is_yonetici.go             → Business logic orchestration
  ├── is_yonetici_test.go        → Business logic unit tests
  ├── veri_yonetici.go           → Data access layer (SQLite)
  ├── veri_yonetici_test.go      → Data access layer unit tests
  └── veri_yonetici_interface.go → Interface for dependency injection
```

### gorev-vscode (TypeScript)
```
src/extension.ts                   → VS Code extension entry point
src/mcp/                          → MCP client implementation
  ├── client.ts                  → MCP protocol client
  └── types.ts                   → TypeScript type definitions
src/commands/                     → VS Code commands
  ├── gorevCommands.ts           → Task-related commands
  ├── projeCommands.ts           → Project-related commands
  └── templateCommands.ts        → Template-related commands
src/providers/                    → TreeView providers
  ├── gorevTreeProvider.ts       → Task tree view
  ├── projeTreeProvider.ts       → Project tree view
  └── templateTreeProvider.ts    → Template tree view
```

### Key Design Decisions

1. **Turkish Domain Language**: Core domain concepts use Turkish terms (gorev=task, proje=project, durum=status, oncelik=priority)
2. **MCP SDK Integration**: Uses `mark3labs/mcp-go` v0.6.0 for MCP protocol implementation
3. **SQLite Storage**: Single-file database for simplicity and portability
4. **No External State**: Each MCP request is stateless, no session management

## Development Commands

### MCP Server (gorev-mcpserver)
```bash
cd gorev-mcpserver

# Build
make build                 # Build for current platform
make build-all            # Build for all platforms (linux, darwin, windows)

# Test
make test                 # Run all tests with coverage
make test-coverage        # Generate HTML coverage report
go test -race ./...       # Run with race detector

# Code Quality
make fmt                  # Format code with gofmt
make lint                 # Run golangci-lint (must be installed)
go vet ./...             # Run go vet

# Dependencies
make deps                 # Download and tidy dependencies

# Docker
make docker-build         # Build Docker image
make docker-run          # Run Docker container

# Development
make run                 # Build and run server
./gorev serve --debug    # Run with debug logging
```

### VS Code Extension (gorev-vscode)
```bash
cd gorev-vscode

# Install dependencies
npm install

# Build
npm run compile          # Compile TypeScript
npm run watch           # Watch mode for development

# Test
npm test                # Run tests

# Package
npm run package         # Create .vsix package

# Development
# Press F5 in VS Code to launch extension development host
```

## MCP Tools

The server implements 19 MCP tools:

### Task Management
1. **gorev_olustur**: Create new task (params: baslik, aciklama, oncelik, proje_id?, son_tarih?, etiketler?)
   - proje_id is optional; if not provided, uses active project
   - son_tarih: optional due date in YYYY-MM-DD format
   - etiketler: optional comma-separated tags
2. **gorev_listele**: List tasks (params: durum?, tum_projeler?, sirala?, filtre?, etiket?)
   - tum_projeler: if false/omitted, shows only active project tasks
   - sirala: son_tarih_asc, son_tarih_desc
   - filtre: acil (due in 7 days), gecmis (overdue)
   - etiket: filter by tag name
3. **gorev_detay**: Show detailed task info in markdown (params: id)
   - Shows due dates, tags, and dependencies with status indicators
4. **gorev_guncelle**: Update task status (params: id, durum)
   - Validates dependencies before allowing "devam_ediyor" status
5. **gorev_duzenle**: Edit task properties (params: id, baslik?, aciklama?, oncelik?, proje_id?, son_tarih?)
6. **gorev_sil**: Delete task (params: id, onay)
   - Prevents deletion if task has subtasks
7. **gorev_bagimlilik_ekle**: Create task dependency (params: kaynak_id, hedef_id, baglanti_tipi)

### Subtask Management
8. **gorev_altgorev_olustur**: Create subtask under a parent (params: parent_id, baslik, aciklama?, oncelik?, son_tarih?, etiketler?)
   - Subtask inherits parent's project
   - parent_id: ID of the parent task
9. **gorev_ust_degistir**: Change task's parent (params: gorev_id, yeni_parent_id?)
   - yeni_parent_id: empty string moves task to root level
   - Validates circular dependencies
10. **gorev_hiyerarsi_goster**: Show task hierarchy (params: gorev_id)
   - Shows parent hierarchy, subtask statistics, and progress

### Task Templates
11. **template_listele**: List available templates (params: kategori?)
   - Shows predefined templates for consistent task creation
12. **templateden_gorev_olustur**: Create task from template (params: template_id, degerler)
   - degerler is an object with field values for the template

### Project Management
13. **proje_olustur**: Create project (params: isim, tanim)
14. **proje_listele**: List all projects with task counts (no params)
15. **proje_gorevleri**: List project tasks grouped by status (params: proje_id)
16. **proje_aktif_yap**: Set active project (params: proje_id)
17. **aktif_proje_goster**: Show current active project (no params)
18. **aktif_proje_kaldir**: Remove active project setting (no params)

### Reporting
19. **ozet_goster**: Show summary statistics (no params)

All tools follow the pattern in `internal/mcp/handlers.go` and are registered in `RegisterTools()`. Task descriptions support full markdown formatting.

## Testing Strategy

- **Unit Tests**: Business logic in `internal/gorev/` (88.2% coverage)
  - `veri_yonetici_test.go`: Data layer tests with SQL injection and concurrent access tests
  - `is_yonetici_test.go`: Business logic tests with mocked dependencies
- **Integration Tests**: MCP handlers in `test/integration_test.go`
- **Table-Driven Tests**: Go convention for test cases
- **Test Database**: Use `:memory:` SQLite for tests
- **Coverage Goal**: Maintain >80% code coverage

Example test pattern:
```go
func TestGorevOlustur(t *testing.T) {
    testCases := []struct {
        name    string
        input   map[string]interface{}
        wantErr bool
    }{
        // test cases...
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Adding New MCP Tools

1. Add handler method to `internal/mcp/handlers.go`
2. Register tool in `RegisterTools()` with proper schema
3. Add integration tests in `test/integration_test.go`
4. Update `docs/mcp-araclari.md` with tool documentation

## Database Schema

SQLite database with seven tables and one view:

- **projeler**: id, isim, tanim, olusturma_tarih, guncelleme_tarih
- **gorevler**: id, baslik, aciklama, durum, oncelik, proje_id, parent_id, olusturma_tarih, guncelleme_tarih, son_tarih
- **baglantilar**: id, kaynak_id, hedef_id, baglanti_tip (for task dependencies)
- **aktif_proje**: id (CHECK id=1), proje_id (stores single active project)
- **etiketler**: id, isim (tags)
- **gorev_etiketleri**: gorev_id, etiket_id (many-to-many relationship)
- **gorev_templateleri**: id, isim, tanim, varsayilan_baslik, aciklama_template, alanlar, ornek_degerler, kategori, aktif (task templates)
- **gorev_hiyerarsi** (VIEW): Recursive CTE view for efficient hierarchy queries with path and level information

Migrations are handled by golang-migrate in `internal/veri/migrations/`.

## Error Handling

- Always return explicit errors, never panic
- MCP errors use `mcp.NewToolResultError()`
- Turkish error messages for user-facing errors
- Wrap errors with context: `fmt.Errorf("context: %w", err)`

## Code Style

- Keep Turkish terms for domain concepts
- Use English for technical terms and comments
- Follow Go idioms and conventions
- Prefer composition over inheritance
- Keep functions small and focused

## Important Files

### gorev-mcpserver
- `internal/gorev/modeller.go`: Domain model definitions (includes GorevTemplate, TemplateAlan, GorevHiyerarsi)
- `internal/mcp/handlers.go`: MCP tool implementations (includes template and subtask handlers)
- `internal/gorev/veri_yonetici.go`: Database operations (includes hierarchy queries with recursive CTEs)
- `internal/gorev/is_yonetici.go`: Business logic (includes subtask validation and circular dependency checks)
- `internal/gorev/template_yonetici.go`: Template management operations
- `cmd/gorev/main.go`: CLI and server initialization (includes template commands, path resolution)
- `internal/veri/migrations/000005_add_parent_id_to_gorevler.up.sql`: Subtask hierarchy migration

### gorev-vscode
- `src/extension.ts`: Extension entry point and activation
- `src/mcp/client.ts`: MCP client for server communication
- `src/providers/enhancedGorevTreeProvider.ts`: Advanced TreeView with grouping and multi-select
- `src/providers/dragDropController.ts`: Drag & drop functionality
- `src/ui/filterToolbar.ts`: Advanced filtering and search
- `src/ui/taskDetailPanel.ts`: Rich task detail view with markdown editor
- `src/ui/templateWizard.ts`: Multi-step template wizard
- `src/utils/markdownParser.ts`: Comprehensive MCP response parser
- `src/commands/*.ts`: Command implementations (21 commands total)
- `package.json`: Extension manifest with commands, views, and configuration
- `test/`: Comprehensive test suite (unit, integration, E2E)

## Version Management

Version info is injected at build time via ldflags:
- `main.version`
- `main.buildTime`
- `main.gitCommit`

The Makefile handles this automatically.