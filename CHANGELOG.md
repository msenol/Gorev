# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Test Infrastructure Improvements** (30 June 2025)
  - MCP Server Test Coverage:
    - Added comprehensive integration tests for all 16 MCP tools achieving 75.1% coverage
    - Created `handlers_test.go` with 561 lines covering complete protocol compliance
    - Created `server_test.go` for MCP server initialization tests
  - VS Code Extension Test Coverage:
    - Achieved 50.9% file coverage (up from 0%) with 19 files tested
    - Added 7 new unit test files totaling 2,700 LOC
    - Created custom test coverage analysis tool (`test-coverage.js`)
    - Test files added:
      - `enhancedGorevTreeProvider.test.js` (389 LOC) - TreeView functionality
      - `taskDetailPanel.test.js` (396 LOC) - WebView panel testing
      - `logger.test.js` (237 LOC) - Logging utility tests
      - `models.test.js` (273 LOC) - TypeScript model validation
      - `utils.test.js` (307 LOC) - Utility function tests

### Fixed
- Database migration issues with `etiketler` table in test environment
- Concurrent access test failures when using in-memory SQLite database
- TypeScript compilation error in `markdownParser.ts` (`proje_ismi` vs `proje_isim`)
- Test coverage reporting in VS Code extension

### Changed
- Improved `gorevEtiketleriniGetir` to handle missing tables gracefully
- Enhanced test coverage analysis to better detect TypeScript imports
- Updated `package.json` with test coverage scripts (`npm run coverage`)

### Documentation
- Updated ROADMAP.md with Filter State Persistence Issue as Task #7

## [0.7.0-beta.1] - 2025-06-29

### Added

#### VS Code Extension - Enhanced UI Features
- **Enhanced TreeView**: Professional task management with grouping, multi-select, color coding
  - Grouping by status/priority/tag/project/due date
  - Multi-select support (Ctrl/Cmd+Click)
  - Expandable/collapsible categories
  - Priority-based color coding
  - Quick completion checkboxes
  - Badges (task counts, due date warnings)
- **Drag & Drop Support**: Intuitive task management with visual feedback
  - Move tasks between projects
  - Change status by dragging
  - Reorder priorities
  - Create dependencies by dropping tasks on each other
  - Visual indicators and animations
- **Inline Editing**: Quick task editing directly in TreeView
  - F2 or double-click to edit
  - Context menu for status/priority changes
  - Inline date picker
  - Escape to cancel, Enter to save
- **Advanced Filtering Toolbar**: Powerful search and filter system
  - Search bar with real-time filtering
  - Advanced filters (status, priority, tags, dates)
  - Saved filter profiles
  - Status bar integration
  - Quick filter shortcuts
- **Rich Task Detail Panel**: Comprehensive task view with markdown editor
  - Split-view markdown editor
  - Live preview
  - Dependency visualization graph
  - Activity timeline
  - Template field indicators
- **Template Wizard UI**: Multi-step interface for template-based task creation
  - Template search and filter
  - Dynamic form generation
  - Field validation
  - Preview before creation
  - Template categories with icons
- **Comprehensive Test Suite**: Unit, integration, and E2E tests
  - Unit tests for all major components
  - Integration tests for extension features
  - End-to-end workflow tests
  - Test fixtures and helpers
  - Coverage reporting with c8

#### MCP Server Improvements
- **Path Resolution**: Fixed database and migration paths to work from any directory
  - `getDatabasePath()` function for executable-relative paths
  - `getMigrationsPath()` function for migration discovery
  - Works correctly when running from different directories

### Changed
- **Enhanced Commands**: Additional 10+ commands for new UI features
- **Configuration**: Extended settings for grouping, sorting, drag-drop behavior
- **Tree Providers**: Complete rewrite with enhanced functionality
- **Template Parser**: Updated to handle new MCP server response format
- **VS Code Extension Version**: Bumped to 0.3.0 for bug fixes and improvements

### Fixed
- **Template Display**: Fixed markdown parser to correctly parse template list responses
- **TreeView Classes**: Exported tree item classes to allow VS Code instantiation
- **TypeScript Errors**: Fixed filter interface property names (Turkish equivalents)
- **Path Issues**: Fixed gorev command execution from different directories
- **Tag Display in VS Code UI** (29 June 2025)
  - Fixed tags not showing in TreeView when tasks created via CLI
  - Updated `GorevListele` handler to include tags and due dates in response
  - Updated `ProjeGorevleri` handler to include tags for all task statuses
- **Project Task Count Display** (29 June 2025)
  - Fixed "0 tasks" showing for all projects in TreeView
  - Updated MarkdownParser to correctly parse "Görev Sayısı" field
- **Task Detail Panel UI Issues** (29 June 2025)
  - Fixed action buttons (Status, Edit, Delete) not visible in dark theme
  - Fixed markdown editor toolbar completely missing
  - Replaced inline event handlers to comply with CSP
  - Fixed edit button to properly pass task data to edit dialog
  - Fixed delete button to use VS Code's native confirmation dialog
  - Added fallback text/emoji for markdown editor buttons
  - Improved dark theme compatibility with explicit color values
- **TreeView Task Selection** (29 June 2025)
  - Tasks now open detail panel on single click
  - Updated TaskTreeViewItem command from selectTask to showTaskDetail
- **Dependency Graph Feature Removed** (29 June 2025)
  - Removed non-functional dependency graph visualization
  - Cleaned up related files and command registrations

### Technical
- **New Files**: 20+ new TypeScript files for enhanced UI
- **Test Infrastructure**: Complete test setup with mocha, sinon, and VS Code test APIs
- **Markdown Parser**: Comprehensive parser for all MCP response types
- **Debug Support**: Test data seeder and debug commands
- **Icons**: Custom SVG icons for tasks, priorities, and templates

### Documentation
- **TASKS.md**: Added comprehensive roadmap with 11 active development tasks
  - DevOps Pipeline and CI/CD automation
  - Test coverage improvements (target 95%)
  - UI/UX and accessibility enhancements
  - Multi-user system and authorization
  - Performance optimizations
  - External service integrations (GitHub, Jira, Slack)
  - Analytics dashboard
  - Advanced search and filtering
  - Dependency visualization in TreeView
  - Subtask system with unlimited hierarchy
  - AI-Powered task enrichment system (NEW)

## [0.6.0] - 2025-06-27

### Added

#### VS Code Extension (gorev-vscode)
- **New module**: Optional VS Code extension for visual task management
- **TreeView panels**: Tasks, Projects, and Templates with visual hierarchy
- **Command palette**: 11 commands including quick task creation (Ctrl+Shift+G)
- **Status bar**: Real-time connection status and task statistics
- **Context menus**: Right-click operations for tasks and projects
- **Theme support**: Priority-based color coding
- **MCP client**: TypeScript implementation connecting to Gorev server
- **Configuration**: Extension settings for server path, auto-connect, and refresh interval

### Changed
- **Project structure**: Reorganized into two modules: `gorev-mcpserver` and `gorev-vscode`
- **Package structure**: Updated module path from `github.com/msenol/gorev` to `github.com/msenol/gorev`
- **Go version**: Updated from 1.21 to 1.22 in go.mod
- **Version**: Aligned version numbers (0.6.0) across VERSION file and Makefile
- **Configuration**: Removed docker-based gorev server configuration from `.mcp.json`

### Fixed
- **Import paths**: Fixed all internal import paths to use the correct module name
- **Test coverage**: Updated badge to reflect actual coverage (53.8%)

### Documentation
- **Enhanced**: Complete documentation overhaul with Claude (Anthropic) assistance
- **Added**: VS Code extension documentation (vscode-extension.md)
- **Added**: Extension README (gorev-vscode/README.md)
- **Added**: Platform-specific installation guides for Windows, macOS, and Linux
- **Added**: MCP editor integration guides for VS Code, Windsurf, Cursor
- **Added**: Comprehensive examples documentation (ornekler.md)
- **Added**: API reference documentation (api-referans.md)
- **Added**: Developer guide (gelistirme.md) with extension development section
- **Updated**: Main README to reflect two-module architecture
- **Updated**: CLAUDE.md to document both modules
- **Updated**: Installation guide with VS Code extension setup
- **Fixed**: All placeholder URLs replaced with actual GitHub repository URLs

## [0.5.0] - 2025-06-27

### Added

#### Task Template System (Görev Şablonları)
- **New feature**: Predefined templates for consistent task creation
- **Database changes**: 
  - Added `gorev_templateleri` table via migration 000004
  - Stores template definitions with dynamic fields
- **Four default templates**:
  - **Bug Raporu** - Detailed bug reporting template with environment, steps, expected/actual behavior
  - **Özellik İsteği** - Feature request template with user stories and acceptance criteria
  - **Teknik Borç** - Technical debt template for refactoring tasks
  - **Araştırma Görevi** - Research task template with objectives and evaluation criteria
- **Template features**:
  - Dynamic field types: text, select, date, number
  - Field validation (required/optional)
  - Default values and select options
  - Template placeholders using `{{field_name}}` syntax
  - Automatic tag and priority assignment
- **New MCP tools**:
  - `template_listele` - List available templates with optional category filter
  - `templateden_gorev_olustur` - Create tasks from templates with custom field values
- **CLI commands**:
  - `gorev template list [kategori]` - List templates by category
  - `gorev template show <template-id>` - Show detailed template information
  - `gorev template init` - Initialize default templates
- **Integration**: Templates are automatically created on database initialization

### Technical
- Added `template_yonetici.go` for template management logic
- Added `GorevTemplate` and `TemplateAlan` structs to domain model
- Extended `VeriYoneticiInterface` with 5 new template methods
- Modified `veri_yonetici.go` to auto-initialize templates after migration
- Added template handlers to MCP server with full schema definitions

## [0.4.0] - 2025-06-27

### Added

#### Task Due Dates (Son Tarih)
- **New feature**: Tasks can now have due dates for deadline tracking
- **Database changes**: Added `son_tarih` column to `gorevler` table via migration
- **Enhanced tools**:
  - `gorev_olustur` - Accepts `son_tarih` parameter in YYYY-MM-DD format
  - `gorev_duzenle` - Can edit due dates
  - `gorev_listele` - New sorting options: `son_tarih_asc`, `son_tarih_desc`
  - `gorev_listele` - New filters: `acil` (due in 7 days), `gecmis` (overdue)
  - `gorev_detay` - Displays due dates in task details

#### Task Tagging System (Etiketler)
- **New feature**: Tasks can be categorized with multiple tags
- **Database changes**:
  - Added `etiketler` table for storing tags
  - Added `gorev_etiketleri` table for many-to-many relationships
- **Enhanced tools**:
  - `gorev_olustur` - Accepts `etiketler` parameter (comma-separated tags)
  - `gorev_listele` - Can filter by `etiket` parameter
  - `gorev_detay` - Displays tags in task details
- Tags are automatically created if they don't exist

#### Task Dependencies (Görev Bağımlılıkları)
- **New feature**: Tasks can have dependencies that must be completed before starting
- **New MCP tool**: `gorev_bagimlilik_ekle` - Create dependencies between tasks
- **Business logic**:
  - Tasks cannot be moved to "devam_ediyor" status if dependencies are incomplete
  - New `GorevBagimliMi` function checks dependency satisfaction
- **Enhanced tools**:
  - `gorev_guncelle` - Validates dependencies before status changes
  - `gorev_detay` - Shows dependencies with completion status indicators:
    - ✅ for completed dependencies
    - ⏳ for pending dependencies
    - Warning message if task cannot be started

#### Database Schema Management
- **Implemented** database schema versioning with `golang-migrate/migrate`
- **Schema files** organized in `internal/veri/migrations/`
- **Schema versions**:
  1. `000001_initial_schema.up.sql` - Base tables
  2. `000002_add_due_date_to_gorevler.up.sql` - Due date support
  3. `000003_add_tags.up.sql` - Tagging system

### Changed

#### Breaking Changes
- **`GorevOlustur`** function signature: Now takes 6 parameters (added `sonTarihStr`, `etiketIsimleri`)
- **`GorevListele`** function signature: Now takes 3 parameters (added `sirala`, `filtre`)
- **`VeriYonetici`** constructor: Now requires migrations path parameter
- **Method renames**:
  - `GorevDetayAl` → `GorevGetir`
  - `ProjeDetayAl` → `ProjeGetir`

#### Improvements
- Enhanced task detail view with better dependency visualization
- Improved error messages for dependency violations
- Better test coverage with dependency validation tests

### Technical
- Added `GorevBagimliMi` method to check task dependencies
- Added `BaglantiEkle` and `BaglantilariGetir` methods for dependency management
- Added `EtiketleriGetirVeyaOlustur` and `GorevEtiketleriniAyarla` for tag management
- Updated mock implementations in tests to support new features
- Integration tests updated for new functionality

## [0.3.0] - 2025-06-25

### Added

#### Active Project Management
- **New feature**: Active project context for task management
- **New MCP tools**:
  - `proje_aktif_yap` - Set a project as the active project
  - `aktif_proje_goster` - Display current active project details
  - `aktif_proje_kaldir` - Remove active project setting
- **Database changes**:
  - Added `aktif_proje` table to store persistent active project setting
  - Table uses CHECK constraint to ensure only one active project (id=1)
- **Enhanced existing tools**:
  - `gorev_olustur` - Now accepts optional `proje_id` parameter; uses active project by default if not specified
  - `gorev_listele` - Added `tum_projeler` boolean parameter; filters by active project by default

### Changed
- **Breaking change**: `GorevOlustur` function now takes 4 parameters (added `projeID`)
- Task creation feedback now includes project name when task is assigned to a project
- Task listing title shows active project name when filtering by active project

### Technical
- Added `veri_yonetici_ext.go` for active project database operations
- Added `AktifProjeAyarla`, `AktifProjeGetir`, `AktifProjeKaldir` methods to VeriYonetici
- Updated VeriYoneticiInterface with active project methods
- Updated IsYonetici to support active project operations
- Enhanced MCP handlers to utilize active project context

## [0.2.0] - 2025-06-25

### Added

#### Unit Testing Infrastructure
- **Comprehensive unit tests** for business logic layer with 88.2% code coverage
- `veri_yonetici_test.go` - Tests for data access layer (VeriYonetici)
  - CRUD operations testing
  - SQL injection protection tests
  - NULL handling tests
  - Concurrent access tests
  - Edge case validation
- `is_yonetici_test.go` - Tests for business logic layer (IsYonetici)
  - Mock implementation of VeriYoneticiInterface
  - Business logic validation
  - Error handling scenarios
  - Partial update logic tests
- `veri_yonetici_interface.go` - Interface for dependency injection and mocking

#### New MCP Tools
- `gorev_detay` - Display detailed task information in markdown format
- `gorev_duzenle` - Edit task title, description, priority, or project assignment
- `gorev_sil` - Delete tasks with confirmation safety
- `proje_listele` - List all projects with task counts
- `proje_gorevleri` - List tasks for a specific project grouped by status

#### Features
- Full markdown support in task descriptions
- Partial update capability for task editing (only specified fields are updated)
- Task count display in project listings
- Status-based grouping in project task views
- Comprehensive integration tests for all new tools

### Changed
- Task descriptions now support full markdown formatting
- Improved error messages to be more user-friendly
- Updated MCP handler signatures to match mark3labs/mcp-go v0.6.0 API
- **Refactored IsYonetici to use VeriYoneticiInterface for better testability**

### Documentation
- Updated `docs/mcp-araclari.md` with detailed documentation for all new tools
- Added examples and response formats for each tool
- Updated future features roadmap

### Technical
- Added `GorevDetayAl`, `ProjeDetayAl`, `GorevDuzenle`, `GorevSil` methods to business logic layer
- Added `ProjeGetir`, `GorevSil`, `ProjeGorevleriGetir` methods to data access layer
- Fixed all integration tests to work with new MCP API
- Added helper function for extracting text from MCP results in tests
- **Implemented dependency injection pattern for better testability**
- **Added table-driven test patterns following Go conventions**
- **Test coverage includes: edge cases, SQL injection, concurrent access, NULL handling**

## [0.1.0] - 2024-12-15

### Added
- Initial release with core MCP server functionality
- Basic task management tools: `gorev_olustur`, `gorev_listele`, `gorev_guncelle`
- Project management: `proje_olustur`
- System overview: `ozet_goster`
- SQLite database backend
- Clean architecture implementation
- Turkish domain language support