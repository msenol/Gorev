# Changelog

All notable changes to the Gorev project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.11.0] - 2025-07-21 - Full Internationalization Release

### Added
- **🌍 Complete MCP Server Internationalization** 
  - Added `go-i18n/v2` library for professional internationalization support
  - **270+ strings converted** from hardcoded Turkish to internationalization system
  - Full Turkish (default) and English language support
  - Dynamic language switching without server restart
  - Language detection hierarchy: CLI flag → GOREV_LANG env → LANG env → Turkish default

### New Translation Infrastructure
- **`internal/i18n/manager.go`** - Complete translation management system
- **`locales/tr.json`** - 270+ Turkish translation keys with organized structure
- **`locales/en.json`** - Complete English translations matching Turkish functionality
- Template data support for dynamic values using {{.Variable}} syntax

### Files Internationalized
- **`internal/mcp/handlers.go`** - All 25 MCP tools with error messages and descriptions
- **`internal/gorev/is_yonetici.go`** - Business logic error messages (25+ strings)
- **`internal/gorev/template_yonetici.go`** - Template system messages (15+ strings)
- **`internal/gorev/veri_yonetici.go`** - Database operation errors (8+ strings)
- **`internal/gorev/ai_context_yonetici.go`** - AI context management messages (5+ strings)
- **`cmd/gorev/main.go`** - CLI command descriptions and help text
- **`cmd/gorev/mcp_commands.go`** - Debug command interfaces (10+ strings)

### Language Features
- **`--lang` CLI flag** for explicit language selection (tr, en)
- **`GOREV_LANG` environment variable** support for configuration
- **Automatic fallback** to system `LANG` environment variable
- All error messages, success messages, and UI text translated
- **Backward compatibility** - No breaking changes to existing Turkish interfaces

### VS Code Extension Updates (from previous release)
- **Bilingual Support (English/Turkish)** for VS Code Extension
  - Automatic language detection based on VS Code language setting
  - Complete localization of all UI strings (500+ translations)
  - Localized command palette, menus, and notifications
  - English and Turkish README files
  - Package.nls files for marketplace localization

### Changed
- **MCP Server** now fully internationalized with professional i18n system
- **Language preference** can be set via environment variables or CLI flags
- **User experience** enhanced for international users
- VS Code extension now supports both English and Turkish interfaces
- All hardcoded strings replaced with localized versions using vscode.l10n.t() API

### Technical Details
- Added dependencies: `go-i18n/v2 v2.6.0`, `golang.org/x/text v0.23.0`
- Maintains 100% backward compatibility
- No breaking changes to MCP tool APIs
- Ready for easy addition of more languages in the future

## [0.10.3] - 2025-07-18 - Test Infrastructure & Stability Release

### 🚀 Major Features
- **Circular Dependency Detection**: Added comprehensive task hierarchy circular dependency prevention
- **Enhanced Task Detail Display**: Parent task information now visible in task details (`GorevDetay`)
- **Active Project Auto-Setup**: Template-based task creation now automatically handles project context
- **Enterprise-Level Test Infrastructure**: Massive test coverage improvements across both modules

### 📊 Test Infrastructure Enhancement
- **MCP Server**: Coverage improved 66.0% → 84.6% (+18.6 percentage points)
  - 3 new comprehensive test files (2,334 lines)
  - `handlers_coverage_test.go` (1,525 lines) - complete MCP handler coverage
  - `handlers_hierarchy_test.go` (523 lines) - task hierarchy and pagination testing
  - `server_coverage_test.go` (286 lines) - MCP server infrastructure testing
- **VS Code Extension**: Coverage improved 55.6% → 100.0% (+44.4 percentage points)
  - 15 new unit test files (~3,000 lines)
  - Complete test coverage for all 36 source files
  - Advanced mock implementations for VS Code APIs

### 🔧 Bug Fixes & Enhancements
- Fixed circular dependency validation in task hierarchy operations
- Enhanced `GorevUstDegistir` with proper dependency checking
- Added "Üst Görev" field display in task detail view
- Improved test isolation to prevent cross-test interference
- Fixed template validation error messages
- Enhanced error handling with contextual Turkish messages

### 🛠️ Technical Improvements
- All tests migrated from deprecated `gorev_olustur` to template-based creation
- Table-driven tests following Go best practices
- Jest/Mocha patterns for TypeScript testing
- Comprehensive edge case testing: SQL injection, Unicode, concurrent access
- Performance benchmarking for bulk operations
- Production-ready test infrastructure

## MCP Server

### [0.10.2] - 2025-07-17

### Added
- Enhanced MCP Debug System with CLI commands
- `gorev mcp list` - List all available MCP tools
- `gorev mcp call <tool> <args>` - Direct tool invocation for testing

### Changed
- VS Code Extension v0.4.6 published to marketplace
- Fixed TypeScript version compatibility issue (5.7.0 → 5.8.3)

### [0.10.1] - 2025-07-11

#### Fixed

- **Critical Pagination Bug**: Fixed duplicate task display issue where subtasks appeared both independently and under their parent
  - Changed pagination to only apply to root-level tasks
  - Subtasks now always appear with their parent task
  - Fixed infinite loop in VS Code when offset exceeded available tasks
  - Fixed task count display to show root tasks count instead of total

#### Changed

- **GorevListele Handler**: Completely rewritten pagination logic
  - Removed `paginatedGorevler` processing for all tasks
  - Now uses `kokGorevler` (root tasks only) for pagination
  - Removed orphan task checking that caused duplicates

### [0.10.0] - 2025-07-11

#### ⚠️ BREAKING CHANGES

- **Template Usage Now Mandatory**: Direct task creation via `gorev_olustur` is no longer supported
  - `gorev_olustur` handler now returns error message with migration instructions
  - All tasks must be created using `templateden_gorev_olustur` with appropriate templates
  - This ensures consistency, quality, and better structured task management

#### Added

- **5 New Enhanced Templates** for better task categorization:
  - `Bug Raporu v2` - Advanced bug reports with severity, steps, environment details
  - `Spike Araştırma` - Time-boxed research tasks with clear deliverables  
  - `Performans Sorunu` - Performance optimization tasks with metrics and targets
  - `Güvenlik Düzeltmesi` - Security vulnerability fixes with CVSS scoring
  - `Refactoring` - Code quality improvements with complexity metrics
- **Enhanced Template Validation**: Detailed error messages show missing required fields and examples
- **Comprehensive Test Coverage**: Added tests for deprecated handler and template workflow

#### Changed

- **GorevOlustur Handler**: Now deprecated with helpful error message guiding users to templates
- **Template System**: Enhanced with better validation and user-friendly error messages

#### Migration Guide

**Old Usage (No Longer Works):**
```bash
gorev_olustur baslik="Bug fix" aciklama="..." oncelik="yuksek"
```

**New Usage (Required):**
```bash
# 1. List available templates
template_listele

# 2. Create task from template  
templateden_gorev_olustur template_id='bug_report_v2' degerler={
  'baslik': 'Login bug',
  'aciklama': 'User cannot login',
  'modul': 'auth',
  'severity': 'high',
  ...
}
```

### [0.9.2] - 2025-07-10

#### Fixed
- **Enhanced Pagination to Show ALL Descendants**: Fixed critical issue where subtasks weren't displayed if their parent task wasn't in the current paginated set
  - Added `gorevHiyerarsiYazdirVeIsaretle` and `gorevHiyerarsiYazdirInternal` functions to track displayed tasks
  - Tasks without visible parents are now shown as root-level tasks
  - Ensures complete task hierarchy is always visible regardless of pagination boundaries

#### Added
- **MCP Debug CLI Commands**: New debugging tools for testing MCP functionality
  - `gorev mcp list` - Lists all available MCP tools
  - `gorev mcp call <tool> [params]` - Calls any MCP tool directly with parameters
  - `gorev mcp list-tasks`, `create-task`, `task-detail`, `projects` - Convenient shortcuts
  - Added `CallTool` method to Handlers for programmatic tool invocation
- **Enhanced Database Path Resolution**: Now checks `~/.gorev/gorev.db` first before relative paths

### [0.9.1] - 2025-07-09

#### Fixed
- **Critical Pagination Bug**: Fixed issue where pagination was counting all tasks (147) but only paginating through root tasks (38)
  - Second page (offset 100+) was returning empty responses
  - Now correctly paginates through all tasks including subtasks
  - VS Code extension now displays all tasks correctly

#### Changed
- Improved pagination logic in `GorevListele` handler to handle both root and subtasks
- Better task organization in paginated responses

### [0.9.0] - 2025-07-09

#### Added
- **AI Context Management & Automation System**
  - 6 new MCP tools for AI-optimized task management
  - Automatic state transitions when AI views tasks (beklemede → devam_ediyor)
  - Persistent AI context across sessions
  - Natural language query support for Turkish
  - Batch operations for efficient bulk updates
  - AI-optimized context summary dashboard

- **New MCP Tools**
  - `gorev_set_active` - Set and track active task with auto-state management
  - `gorev_get_active` - Get current active task with full details
  - `gorev_recent` - List recently interacted tasks (limit parameter)
  - `gorev_context_summary` - AI-optimized session overview with statistics
  - `gorev_batch_update` - Bulk update multiple tasks in single operation
  - `gorev_nlp_query` - Natural language task search (Turkish support)

- **Database Changes**
  - Added migration 000006 for AI context tables
  - New tables: `ai_interactions`, `ai_context`
  - New columns in `gorevler`: `last_ai_interaction`, `estimated_hours`, `actual_hours`
  - Tracking for AI interactions: viewed, created, updated, completed, set_active, bulk_operation

- **Code Architecture**
  - Added `internal/gorev/ai_context_yonetici.go` - AI context manager
  - Added `AIContextYonetici` field to MCP Handlers
  - Integrated auto-state management into `GorevDetay` handler
  - Added comprehensive unit tests (`ai_context_yonetici_test.go`)
  - Added integration tests (`handlers_ai_test.go`)

#### Changed
- `GorevDetay` now automatically transitions tasks from "beklemede" to "devam_ediyor" when viewed
- Updated total MCP tools count from 19 to 25
- Enhanced `is_yonetici.go` with `VeriYonetici()` method for AI context access

#### Documentation
- Updated CLAUDE.md with v0.9.0 features
- Updated ROADMAP.md with AI Context Management as priority 5
- Added comprehensive AI feature documentation

### [0.8.1] - 2025-07-09

#### Fixed
- Fixed token limit errors in MCP tools (`gorev_listele` and `proje_gorevleri` exceeding 25k tokens)
- Implemented pagination support with `limit` and `offset` parameters
- Optimized response formatting to reduce token usage by ~60%
  - Priority shown as Y/O/D instead of full words (yuksek/orta/dusuk)
  - Task details condensed to single line with pipe separators
  - Descriptions truncated to 100 chars, IDs to 8 chars
  - Removed empty fields and unnecessary newlines
  - Simplified section headers (removed ### markdown)
- Added response size estimation to prevent token limit errors
- Maximum response size set to 20K characters to stay under 25K token limit

#### Changed
- Updated MCP tool schemas for `gorev_listele` and `proje_gorevleri` to include pagination parameters

## VS Code Extension

### [0.4.0] - 2025-07-11

#### ⚠️ BREAKING CHANGES

- **Template Usage Now Mandatory**: Task creation commands now redirect to template wizard
  - `gorev.createTask` (Ctrl+Shift+G) now opens Template Wizard instead of direct creation dialog
  - `gorev.quickCreateTask` now opens Quick Template Selection instead of simple input
  - Direct task creation via deprecated `gorev_olustur` is no longer supported

#### Changed

- **Updated Task Creation Flow**: All task creation now uses template system for better consistency
- **Enhanced Error Handling**: Better user guidance when attempting deprecated operations
- **TestDataSeeder Modernization**: Development test data now uses template-based creation

#### Added

- **Breaking Change Documentation**: Comprehensive migration guide in README and CHANGELOG
- **Template System Integration**: Seamless integration with MCP server v0.10.0 template requirement

#### Migration Guide

**Old Workflow:**
1. Press `Ctrl+Shift+G` → Simple task creation dialog
2. Fill basic fields (title, description, priority)
3. Task created directly

**New Workflow:**
1. Press `Ctrl+Shift+G` → Template Wizard opens
2. Select appropriate template (Bug Report, Feature Request, etc.)
3. Fill template-specific required fields
4. Task created with consistent structure and quality

**Benefits:**
- **Better Quality**: Required fields prevent incomplete tasks
- **Consistency**: All tasks follow established patterns
- **Automation**: Template-based workflow automation
- **Reporting**: Better metrics and categorization

### [0.3.9] - 2025-07-10

#### Added
- **"Show All Projects" Toggle Feature**: Easy switching between viewing all tasks vs. active project only
  - New configuration: `gorev.treeView.showAllProjects` (default: true)
  - Keyboard shortcut: `Ctrl+Alt+P` / `Cmd+Alt+P` to toggle views
  - Visual indicator in status bar with globe/folder icon
  - New command: `gorev.toggleAllProjects`

#### Fixed
- **Enhanced Markdown Parser**: Better parsing of MCP server responses
  - Fixed parsing of Turkish priority names (düşük/orta/yüksek) alongside short forms (D/O/Y)
  - Better handling of tasks with emojis in titles
  - Improved subtask parsing with └─ prefix
  - Fixed "X adet" tag count format parsing
- **Filter Toolbar Improvements**:
  - Added "Tüm Projeler" toggle button to filter options
  - Fixed `clearAllFilters` to properly reset showAllProjects state
  - Visual feedback for active filter state

#### Changed
- Replaced console.log with Logger.debug for better debugging
- Template wizard now passes values as object instead of JSON string

### [0.3.4] - 2025-07-09

#### Added
- Added `gorev.pagination.pageSize` configuration option (default: 100)
- Extension now passes pagination parameters to MCP tools

#### Changed
- Updated `enhancedGorevTreeProvider` to use pagination when calling MCP tools

## VS Code Extension

### [0.3.3] - 2025-06-30

#### Fixed
- Fixed progress percentage display issue in task detail panel
- Circular progress chart percentage was not visible
- Implemented CSS overlay solution with absolute positioning
- Progress percentage now displays correctly in the center of the progress indicator
- Enhanced hierarchy progress parsing with more flexible pattern matching
- Added fallback calculation when server response doesn't include percentage
- Fixed potential NaN values in progress display with proper validation
- Fixed dependency section not showing in task detail panel when task has no dependencies
- Dependencies section now always visible with proper fallback messages
- Enhanced dependency information display with clear status indicators

### [0.3.2] - 2025-06-30

#### Added
- Enhanced TreeView with grouping, multi-select, and priority-based color coding
- Drag & Drop support for moving tasks, changing status, and creating dependencies
- Inline editing with F2/double-click, context menus, and date pickers
- Advanced filtering toolbar with search, filters, and saved profiles
- Rich task detail panel with markdown editor and dependency visualization
- Template wizard UI with multi-step interface and dynamic forms
- Hierarchical task display with tree structure
- Progress indicators showing subtask completion percentage

#### Fixed
- Fixed filter state persistence issue
- Fixed tag display when tasks created via CLI
- Fixed project task count showing as 0 in TreeView
- Fixed task detail panel UI issues in dark theme
- Fixed single-click task selection in TreeView

## MCP Server

### [0.8.0] - 2025-06-30

#### Added
- Implemented subtask system with unlimited hierarchy
  - Added `parent_id` column to tasks table with foreign key constraint
  - Created recursive CTE views for efficient hierarchy queries
  - New MCP tools:
    - `gorev_altgorev_olustur` - Create subtask under a parent task
    - `gorev_ust_degistir` - Move task to different parent or root
    - `gorev_hiyerarsi_goster` - Show complete task hierarchy with statistics
  - Circular dependency prevention
  - Parent task progress tracking based on subtask completion

- **VS Code Extension**
  - Hierarchical task display with tree structure
  - Progress indicators showing subtask completion percentage
  - Visual hierarchy with indentation and tree connectors

### Changed
- **Business Rules**
  - Tasks cannot be deleted if they have subtasks
  - Tasks cannot be completed unless all subtasks are completed
  - Moving a task to a different project moves all its subtasks
  - Subtasks inherit parent's project

## [0.7.1] - 2025-06-30

### Fixed
- **VS Code Extension**
  - Fixed filter state persistence issue
  - Added `clearFilters()` method to `EnhancedGorevTreeProvider`
  - Fixed `clearAllFilters()` in FilterToolbar to properly reset state
  - Added keyboard shortcut `Ctrl+Alt+R` / `Cmd+Alt+R` for quick filter clearing

## [0.7.0-beta.1] - 2025-06-30

### Added
- **Test Infrastructure**
  - MCP Server test coverage improved from 75.1% to 81.5%
  - VS Code Extension achieved 50.9% file coverage
  - Created comprehensive test suites for both components
  - Added test coverage analysis tools

- **VS Code Extension Features**
  - Enhanced TreeView with grouping, multi-select, and priority-based color coding
  - Drag & Drop support for moving tasks, changing status, and creating dependencies
  - Inline editing with F2/double-click, context menus, and date pickers
  - Advanced filtering toolbar with search, filters, and saved profiles
  - Rich task detail panel with markdown editor and dependency visualization
  - Template wizard UI with multi-step interface and dynamic forms

### Fixed
- **MCP Server**
  - Fixed path resolution for database and migrations
  - Enhanced `GorevListele` and `ProjeGorevleri` handlers to include tags and due dates

- **VS Code Extension**
  - Fixed tag display when tasks created via CLI
  - Fixed project task count showing as 0 in TreeView
  - Fixed task detail panel UI issues in dark theme
  - Fixed single-click task selection in TreeView

### Removed
- Non-functional dependency graph feature

## [0.6.0] - 2025-06-29

### Added
- Task dependency system with `gorev_bagimlilik_ekle` tool
- Due dates functionality for tasks
- Enhanced filtering with date-based filters (urgent/overdue)

## [0.5.0] - 2025-06-29

### Added
- Task template system with predefined templates
- Tagging system for task categorization
- Database schema version control with golang-migrate
- New MCP tools:
  - `template_listele` - List available templates
  - `templateden_gorev_olustur` - Create tasks from templates