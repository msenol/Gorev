# CLAUDE.md

This file provides guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

## Last Updated: 21 July 2025

> 🤖 **Documentation Note**: This comprehensive technical guide was enhanced and structured with the assistance of Claude (Anthropic), demonstrating the power of AI-assisted documentation in modern software development.

### Recent Changes

#### VS Code Extension (v0.5.0) - Complete Bilingual Support (21 July 2025)
- **English/Turkish Localization**: Complete bilingual support for international users
  - Automatic language detection using `vscode.env.language`
  - 500+ UI strings localized across all 36 source files
  - Modern VS Code l10n API implementation with `vscode.l10n.t()`
  - Bundle-based localization structure (`l10n/bundle.l10n.json`)
  - Marketplace metadata localization (`package.nls.json`)
- **Localized Components**:
  - All 21 VS Code commands with titles and descriptions
  - TreeView providers: task tree, project tree, template tree
  - UI components: filter toolbar, status bar, task detail panel, template wizard
  - Drag-drop controller with operation feedback messages
  - Inline edit provider with validation messages
  - Debug tools and test data seeders
- **Files Added**:
  - `l10n/bundle.l10n.json` - English runtime strings
  - `l10n/bundle.l10n.tr.json` - Turkish translations
  - `package.nls.json` - English VS Code marketplace metadata
  - `package.nls.tr.json` - Turkish VS Code marketplace metadata
  - `README.tr.md` - Turkish README for Turkish users
- **Technical Implementation**:
  - Replaced all hardcoded strings with l10n.t() calls
  - Maintained icon codes and formatting in translations
  - Used placeholder syntax {0}, {1} for dynamic values
  - Consistent key naming pattern (component.key)
  - Preserved all special characters and escape sequences

#### Comprehensive Test Infrastructure Enhancement (17 July 2025)
- **Massive Test Coverage Improvement**: Systematic test development across both modules
  - **MCP Server**: Coverage improved from 66.0% to 81.3% (+15.3 percentage points)
  - **VS Code Extension**: Coverage improved from 55.6% to 100.0% (+44.4 percentage points)
  - **Total Test Files Added**: 18 new comprehensive test files (~5,300 lines of test code)
- **MCP Server Test Enhancements**: 
  - Added 3 new test files: `handlers_coverage_test.go`, `handlers_hierarchy_test.go`, `server_coverage_test.go`
  - Comprehensive edge case testing: SQL injection, Unicode, concurrent access, performance benchmarks
  - Template system migration: All tests converted from deprecated `gorev_olustur` to template-based creation
  - Status validation enhancement: Added proper validation for task status transitions
  - Error handling: Complete error scenario testing with Turkish localization
- **VS Code Extension Test Infrastructure**: 
  - Created 15 new unit test files covering all previously untested components
  - Command testing: Enhanced task commands, template commands, project commands, filter commands
  - UI component testing: FilterToolbar, StatusBar, DecorationProvider, GroupingStrategy
  - Debug utility testing: Debug commands, MCP debug commands, test data seeders
  - Mock integration: Comprehensive VS Code API and MCP client mocking
- **Quality Assurance**: 
  - Table-driven tests following Go best practices
  - Jest/Mocha patterns for TypeScript testing
  - Real-world scenario testing with actual user workflows
  - Comprehensive error handling and edge case coverage
- **Technical Excellence**: 
  - Production-ready test infrastructure for both modules
  - Maintainable and well-documented test patterns
  - Integration testing with proper mock strategies
  - Performance testing with timing benchmarks

#### MCP Server (v0.10.2) - MCP Debug System & TypeScript Fix (17 July 2025)
- **Enhanced MCP Debug System**: Added comprehensive CLI commands for debugging MCP server functionality
  - `gorev mcp list` - List all available MCP tools
  - `gorev mcp call <tool> <args>` - Direct tool invocation for testing
  - Enhanced debugging capabilities with detailed MCP communication logging
- **VS Code Extension v0.4.6**: Published to marketplace with TypeScript dependency fixes
  - Fixed TypeScript version compatibility issue (5.7.0 → 5.8.3)
  - Fixed npm compilation paths and dependencies
  - Resolved WSL symlink permission errors with `--no-bin-links` flag
  - Enhanced duplicate detection logging with context
  - Added toggle for "Show All Projects" feature
- **Release Management**: Complete v0.10.2 release with all platform binaries
  - Updated installation scripts to default to v0.10.2
  - Created comprehensive release notes with all changes
  - Published GitHub release with checksums and binaries
  - Updated README.md and version references across all files

#### MCP Server (v0.10.1) - Critical Pagination Fix
- **Fixed Duplicate Task Display Issue** (11 July 2025):
  - Fixed critical bug where subtasks appeared twice: once as independent tasks and again under their parent
  - Root cause: Pagination was applied to all tasks (root + subtasks) instead of just root tasks
  - **Solution**: Changed pagination logic to only paginate root-level tasks
  - Subtasks now always appear with their parent, regardless of pagination window
  - Fixed infinite loop issue in VS Code when requesting pages beyond available data
  - **Technical Details**:
    - Modified `GorevListele` handler to use `kokGorevler` (root tasks) for pagination
    - Removed orphan task checking logic that caused duplicates
    - Fixed task count display to show root task count instead of total
  - **Files Updated**:
    - `internal/mcp/handlers.go` - Pagination logic rewrite

#### MCP Server (v0.10.0) - Template Usage Now Mandatory
- **BREAKING CHANGE: `gorev_olustur` Deprecated** (10 July 2025):
  - Direct task creation without templates is no longer allowed
  - All tasks must be created using `templateden_gorev_olustur`
  - Added comprehensive error message guiding users to template usage
  - **New Templates Added**:
    - `bug_report_v2` - Enhanced bug reporting with severity and environment
    - `spike_research` - Time-boxed technical research tasks
    - `performance_issue` - Performance problems with metrics
    - `security_fix` - Security vulnerabilities with CVSS scoring
    - `refactoring` - Code quality improvements with risk assessment
  - **Enhanced Validation**:
    - Strict enforcement of required fields
    - Select field value validation
    - Detailed error messages with examples
  - **Helper Functions**:
    - `templateZorunluAlanlariListele` - Lists required fields
    - `templateOrnekDegerler` - Generates example values

#### MCP Server (v0.9.2) - Enhanced Pagination & Debug Tools
- **Fixed Pagination to Show ALL Descendants** (10 July 2025):
  - Fixed critical issue where subtasks weren't shown if parent wasn't in paginated set
  - Added `gorevHiyerarsiYazdirVeIsaretle` and `gorevHiyerarsiYazdirInternal` functions
  - Tasks without visible parents now shown as root-level tasks
  - Ensures complete task hierarchy is always visible
  - **New Features**:
    - Added MCP debug CLI commands (`gorev mcp list`, `gorev mcp call`)
    - Added `CallTool` method for direct tool invocation
    - Enhanced database path resolution (checks ~/.gorev first)
  - **Files Updated**:
    - `internal/mcp/handlers.go` - New hierarchy display functions
    - `cmd/gorev/mcp_commands.go` - New debug commands
    - `internal/mcp/server.go` - Added helper methods

#### MCP Server (v0.9.1) - Pagination Fix
- **Fixed Critical Pagination Bug** (9 July 2025):
  - Fixed issue where pagination was counting all tasks (147) but only paginating through root tasks (38)
  - Second page (offset 100+) was returning empty responses
  - Now correctly paginates through all tasks including subtasks
  - VS Code extension now displays all tasks correctly
  - **Updated Functions**:
    - `GorevListele` handler - pagination logic rewritten to handle all tasks
    - Added `paginatedGorevMap` for proper hierarchy display
  - **Performance**:
    - Added `BulkBagimlilikSayilariGetir` method for N+1 query optimization
    - Added `BulkTamamlanmamiaBagimlilikSayilariGetir` method
    - Added database indexes migration (000007)

#### VS Code Extension (v0.4.6) - Enhanced Duplicate Detection
- **Enhanced Markdown Parser Logging** (11 July 2025):
  - Added detailed duplicate detection logging with context
  - Logs now include parent ID, project ID, and duplicate count
  - Improved debugging capabilities for server-side issues
  - **Files Updated**:
    - `src/utils/markdownParser.ts` - Enhanced duplicate detection logging

#### VS Code Extension (v0.3.9) - All Projects Toggle & Parser Improvements
- **Added "Show All Projects" Toggle** (10 July 2025):
  - New configuration: `gorev.treeView.showAllProjects` (default: true)
  - Keyboard shortcut: `Ctrl+Alt+P` / `Cmd+Alt+P` to toggle
  - Visual indicator in status bar (globe/folder icon)
  - Fixed issue where only active project tasks were showing
- **Enhanced Markdown Parser**:
  - Fixed parsing of Turkish priority names (düşük/orta/yüksek)
  - Better handling of tasks with emojis in titles
  - Improved subtask parsing with └─ prefix
  - Added Logger.debug instead of console.log
- **Filter Toolbar Improvements**:
  - Added "Tüm Projeler" toggle button
  - Fixed clearAllFilters to reset showAllProjects state
  - Visual feedback for active filter state

#### VS Code Extension (v0.3.8)
- **Fixed Task Display Issues** (9 July 2025):
  - Fixed pagination to fetch all tasks from MCP server
  - Updated `enhancedGorevTreeProvider` to handle multiple pages
  - Updated `markdownParser` to include all tasks (root and subtasks)
  - Added debug mode features for troubleshooting connection issues

#### MCP Server (v0.9.0) - AI Context Management
- **Implemented AI Context Management & Automation** (9 July 2025):
  - Added 6 new MCP tools for AI context management and NLP queries
  - **New Tables**: `ai_interactions`, `ai_context` for tracking AI sessions
  - **Auto-State Management**: Tasks automatically transition to "devam_ediyor" when viewed
  - **Context Persistence**: Active task maintained across AI interactions
  - **Batch Operations**: New `gorev_batch_update` for efficient bulk updates
  - **Natural Language Queries**: New `gorev_nlp_query` for intuitive task search
  - **AI-Optimized Summary**: `gorev_context_summary` provides session overview
  - **New Tools**:
    - `gorev_set_active` - Set and track active task
    - `gorev_get_active` - Get current active task
    - `gorev_recent` - List recent task interactions
    - `gorev_context_summary` - AI-optimized session summary
    - `gorev_batch_update` - Bulk task updates
    - `gorev_nlp_query` - Natural language task search
  - **Database Updates**:
    - Added `last_ai_interaction`, `estimated_hours`, `actual_hours` to tasks
    - Created AI interaction tracking tables
    - Migration 000006 adds AI context support

#### MCP Server (v0.8.1)
- **Fixed Token Limit Errors with Pagination and Compact Formatting** (9 July 2025):
  - Added pagination support to `gorev_listele` and `proje_gorevleri` MCP tools
  - Added `limit` (default: 50) and `offset` (default: 0) parameters
  - Implemented response size estimation to prevent token limit errors
  - Maximum response size set to 20K characters to stay under 25K token limit
  - **Optimized Response Formatting**:
    - Reduced verbosity: priority shown as Y/O/D instead of yuksek/orta/dusuk
    - Condensed task details to single line with pipe separators
    - Truncated descriptions to 100 chars, IDs to 8 chars
    - Removed empty fields and unnecessary newlines
    - Simplified section headers (removed ### markdown)
  - **VS Code Extension Update**:
    - Added `gorev.pagination.pageSize` configuration (default: 100)
    - Extension now passes pagination parameters to MCP tools
  - **New/Updated Functions**:
    - `gorevResponseSizeEstimate` - Estimates response size for a task
    - `gorevOzetYazdir` - Compact task formatting (reduced by ~60%)
    - `gorevOzetYazdirTamamlandi` - Ultra-compact completed task format
    - `gorevHiyerarsiYazdir` - Optimized hierarchical task display

#### VS Code Extension (v0.3.4)
- **Enhanced TreeView Visual Indicators** (5 July 2025):
  - Added visual progress bars for parent tasks showing subtask completion
  - Implemented priority badges (🔥⚡ℹ️) with color coding
  - Added smart due date formatting (Today, Tomorrow, 3d left, etc.)
  - Enhanced dependency badges with lock/unlock status (🔒🔓🔗)
  - Converted tags to colored pill badges with configurable limit
  - Created rich markdown tooltips with progress visualization
  - Added TaskDecorationProvider for managing visual decorations
  - Added 9 new configuration options for visual preferences
  - Improved task description formatting with separators
  - **Fixed dependency data transmission** - MCP handlers now include dependency counts in markdown output
  - **Updated MarkdownParser** to parse dependency information from server responses

#### VS Code Extension (v0.3.3)
- **Fixed Progress Percentage Display Issue** (30 June 2025):
  - Fixed circular progress chart percentage not being visible in task detail panel
  - Implemented CSS overlay solution with absolute positioning for percentage text
  - Added `.percentage-overlay` class with proper centering and theme-aware styling
  - Progress percentage now displays correctly in the center of the circular progress indicator
- **Enhanced Progress Percentage Parsing** (30 June 2025):
  - Improved `parseHierarchyInfo` method with more flexible pattern matching
  - Added fallback calculation when server doesn't provide percentage
  - Added validation to ensure percentage is always a valid number (0-100)
  - Added debug logging for troubleshooting hierarchy parsing issues
  - Fixed potential NaN values in progress display

#### MCP Server (v0.8.0)

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

The server implements 25 MCP tools:

### Task Management
1. **gorev_olustur**: Create new task (params: baslik, aciklama, oncelik, proje_id?, son_tarih?, etiketler?)
   - proje_id is optional; if not provided, uses active project
   - son_tarih: optional due date in YYYY-MM-DD format
   - etiketler: optional comma-separated tags
2. **gorev_listele**: List tasks (params: durum?, tum_projeler?, sirala?, filtre?, etiket?, limit?, offset?)
   - tum_projeler: if false/omitted, shows only active project tasks
   - sirala: son_tarih_asc, son_tarih_desc
   - filtre: acil (due in 7 days), gecmis (overdue)
   - etiket: filter by tag name
   - limit: maximum number of tasks to return (default: 50)
   - offset: number of tasks to skip for pagination (default: 0)
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
15. **proje_gorevleri**: List project tasks grouped by status (params: proje_id, limit?, offset?)
   - limit: maximum number of tasks to return (default: 50)
   - offset: number of tasks to skip for pagination (default: 0)
16. **proje_aktif_yap**: Set active project (params: proje_id)
17. **aktif_proje_goster**: Show current active project (no params)
18. **aktif_proje_kaldir**: Remove active project setting (no params)

### Reporting
19. **ozet_goster**: Show summary statistics (no params)

### AI Context Management (NEW)
20. **gorev_set_active**: Set active task for AI session (params: task_id)
   - Automatically transitions task to "devam_ediyor" status
   - Maintains context across AI interactions
21. **gorev_get_active**: Get current active task (no params)
   - Returns detailed information about the active task
22. **gorev_recent**: Get recent tasks interacted with (params: limit?)
   - limit: number of recent tasks to return (default: 5)
23. **gorev_context_summary**: Get AI-optimized session summary (no params)
   - Shows active task, recent tasks, priorities, and blockers
24. **gorev_batch_update**: Batch update multiple tasks (params: updates)
   - updates: array of {id: string, updates: {durum?: string, ...}}
   - Efficient bulk operations for AI workflows
25. **gorev_nlp_query**: Natural language task search (params: query)
   - Supports queries like: "bugün üzerinde çalıştığım görevler", "yüksek öncelikli", "database ile ilgili"
   - Tag search: "etiket:bug" or "tag:frontend"
   - Status queries: "tamamlanmamış", "devam eden", "tamamlanan"
   - Time-based: "acil", "gecikmiş", "son oluşturulan"

All tools follow the pattern in `internal/mcp/handlers.go` and are registered in `RegisterTools()`. Task descriptions support full markdown formatting.

**Auto-State Management**: When `gorev_detay` is called, tasks automatically transition from "beklemede" to "devam_ediyor" status to reflect AI interaction.

**Template System**: As of v0.10.0, the `gorev_olustur` tool is deprecated. All task creation must use templates via `templateden_gorev_olustur`.

## Testing Strategy

- **Unit Tests**: Business logic in `internal/gorev/` (65.4% coverage), MCP handlers in `internal/mcp/` (81.3% coverage)
  - `veri_yonetici_test.go`: Data layer tests with SQL injection and concurrent access tests
  - `is_yonetici_test.go`: Business logic tests with mocked dependencies
  - `handlers_coverage_test.go`: Comprehensive MCP handler testing with edge cases (1,525 lines)
  - `handlers_hierarchy_test.go`: Subtask hierarchy and pagination testing (523 lines)
  - `server_coverage_test.go`: MCP server infrastructure testing (286 lines)
- **Integration Tests**: MCP handlers in `test/integration_test.go`
- **VS Code Extension Tests**: Complete test coverage (100.0%) with 27 unit test files
  - Command testing: All VS Code commands with proper mocking
  - UI component testing: FilterToolbar, StatusBar, decorations
  - Provider testing: Tree providers, grouping strategies, drag-drop
  - Debug utility testing: Debug commands, test data seeders
- **Table-Driven Tests**: Go convention for test cases
- **Test Database**: Use `:memory:` SQLite for tests  
- **Coverage Achievement**: MCP server 81.3%, VS Code extension 100.0%

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

SQLite database with nine tables and one view:

- **projeler**: id, isim, tanim, olusturma_tarih, guncelleme_tarih
- **gorevler**: id, baslik, aciklama, durum, oncelik, proje_id, parent_id, olusturma_tarih, guncelleme_tarih, son_tarih, last_ai_interaction, estimated_hours, actual_hours
- **baglantilar**: id, kaynak_id, hedef_id, baglanti_tip (for task dependencies)
- **aktif_proje**: id (CHECK id=1), proje_id (stores single active project)
- **etiketler**: id, isim (tags)
- **gorev_etiketleri**: gorev_id, etiket_id (many-to-many relationship)
- **gorev_templateleri**: id, isim, tanim, varsayilan_baslik, aciklama_template, alanlar, ornek_degerler, kategori, aktif (task templates)
- **ai_interactions**: id, gorev_id, action_type, context, timestamp (tracks AI-task interactions)
- **ai_context**: id (CHECK id=1), active_task_id, recent_tasks, session_data, last_updated (AI session context)
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
- `internal/mcp/handlers.go`: MCP tool implementations (includes template, subtask, and AI context handlers)
- `internal/gorev/veri_yonetici.go`: Database operations (includes hierarchy queries with recursive CTEs)
- `internal/gorev/is_yonetici.go`: Business logic (includes subtask validation and circular dependency checks)
- `internal/gorev/ai_context_yonetici.go`: AI context management and auto-state transitions
- `internal/gorev/template_yonetici.go`: Template management operations
- `cmd/gorev/main.go`: CLI and server initialization (includes template commands, path resolution)
- `internal/veri/migrations/000005_add_parent_id_to_gorevler.up.sql`: Subtask hierarchy migration
- `internal/veri/migrations/000006_add_ai_context_tables.up.sql`: AI context tables migration

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
- `test/`: Comprehensive test suite (unit, integration, E2E) - **100.0% coverage achieved**
  - `test/unit/`: 27 unit test files covering all source files
  - `test/integration/`: Integration tests for VS Code API interaction
  - `test/e2e/`: End-to-end workflow testing

## Version Management

Version info is injected at build time via ldflags:
- `main.version`
- `main.buildTime`
- `main.gitCommit`

The Makefile handles this automatically.