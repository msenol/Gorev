# Changelog

All notable changes to Gorev will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.15.22] - 2025-09-21

### Fixed
- **Critical i18n System Enhancement**: Fixed remaining untranslated keys for NPX package environments
  - **Root Cause**: Missing `error.noArguments` and `error.parameterRequired` keys in locale files
  - **Solution**: Added missing error keys to both Turkish and English locale files
  - **Implementation**: Updated embedded locale system with missing keys for NPX compatibility
  - **Fallback Enhancement**: Added fallback strings in handlers for all error keys
  - **User Impact**: Complete i18n coverage ensures robust operation across all deployment methods
  - **Affected Components**: MCP server, NPX package, VS Code extension

### Changed
- **Version Synchronization**: Updated all components to v0.15.22 for consistency
  - **gorev-mcpserver**: v0.15.20 → v0.15.22
  - **gorev-npm**: v0.15.20 → v0.15.22
  - **gorev-vscode**: v0.15.21 → v0.15.22
  - **VSIX Package**: Created gorev-vscode-0.15.22.vsix (325.75 KB, 147 files)

## [v0.15.20] - 2025-09-21

### Fixed
- **Migration State Repair**: Fixed "table already exists" errors in NPX package environments
  - **Root Cause**: Existing databases had tables but missing migration state records
  - **Solution**: Added `repairMigrationStateIfNeeded()` function to auto-detect and repair migration state
  - **Implementation**: Checks for existing core tables (projeler, gorevler, baglantilar) and reconstructs migration history
  - **Safety**: Non-destructive repair that only adds missing migration records
  - **User Impact**: NPX package now works seamlessly with existing databases

### Added
- **Linux ARM64 Support**: Added native binary support for ARM64 Linux platforms
  - **New Binary**: `gorev-npm/binaries/linux-arm64/gorev` for Raspberry Pi and ARM64 servers
  - **Cross-platform Build**: Complete multi-architecture support (Windows x64, macOS x64/ARM64, Linux x64/ARM64)
  - **NPM Package**: Automatic platform detection and binary download

### Changed
- **Version Update**: NPM package version bumped to v0.15.20 for proper semver tracking

## [v0.15.19] - 2025-09-21

### Fixed
- **Critical VS Code Extension Bug**: Fixed task dependency display issue where dependencies weren't showing in VS Code extension
  - **MCP Server Fix**: Expanded dependency type support in `handlers.go` to include "blocker" and "depends_on" types alongside "onceki"
  - **Parser Enhancement**: Fixed VS Code markdown parser format detection for mixed format with emoji status icons (⏳, 🔄, ✅)
  - **Regex Update**: Updated legacy parser regex to support both emoji and text status formats
  - **Priority Support**: Added support for priority letters (Y, O, D) in addition to Turkish text priorities
  - **Result**: All 8 tasks now display with real UUIDs (no fallback IDs) and dependencies show correctly
- **Localization**: Added missing i18n keys for better consistency in Turkish and English interfaces

## [v0.15.8] - 2025-09-20

### Changed
- **Version Standardization**: Standardized versions across all components to v0.15.8
- **Build System Improvements**: Enhanced cross-platform build process
- **Release Process**: Automated release artifact generation and publishing

### Added
- **Windows Support**: Complete Windows binary support with .exe packaging
- **Multi-Architecture macOS**: Support for both Intel (amd64) and Apple Silicon (arm64) Macs
- **Automated Checksums**: SHA256 checksum generation for all release artifacts

### Fixed
- **Build Script Issues**: Fixed truncated build scripts and improved reliability
- **Version Consistency**: Ensured all components use consistent version numbering

## [v0.15.5] - 2025-09-18

### Fixed
- **Critical NPX Fix**: Resolved "error.dataManagerInit" issue in NPX package environments
  - NPX package couldn't find migration files, causing database initialization failures
  - Embedded all migration files directly into Go binary using `//go:embed`
  - NPX package now works without GOREV_ROOT environment variable

### Added
- **Embedded Migrations Architecture**: Complete migration embedding system
  - New `cmd/gorev/migrations_embed.go` with embedded filesystem support
  - `YeniVeriYoneticiWithEmbeddedMigrations()` function for embedded FS support
  - `migrateDBWithFS()` extracts embedded migrations to temporary directory
  - `createVeriYonetici()` unified helper with embedded/filesystem fallback
  - All template commands updated to use embedded migrations

### Changed
- **Binary Size**: Slight increase due to embedded migration files
- **Performance**: Temporary directory extraction during first migration

### Technical
- Maintains backward compatibility with existing filesystem migrations
- Enhanced `internal/gorev/veri_yonetici.go` with embed.FS support
- Updated package versions: gorev-mcpserver v0.15.5, gorev-npm 0.15.5

## [v0.15.4] - 2025-09-18

### 🚀 Features
- **NPX Easy Installation System** - Complete NPM package distribution for effortless setup
  - **@mehmetsenol/gorev-mcp-server Package**: New NPM package enabling `npx @mehmetsenol/gorev-mcp-server@latest` usage
  - **Cross-Platform Binary Support**: Automatic binary download for Windows, macOS, Linux (amd64/arm64)
  - **Zero Installation Setup**: Users can run Gorev without manual binary installation steps
  - **Simple MCP Configuration**: Easy addition to `mcp.json` with `"command": "npx", "args": ["@mehmetsenol/gorev-mcp-server@latest"]`
  - **GitHub Actions Pipeline**: Automated NPM publishing with multi-platform binary builds
  - **Platform Detection**: Intelligent platform and architecture detection for correct binary selection
  - **Fallback Mechanisms**: Robust error handling and fallback to latest releases

### 🔧 Implementation
- **gorev-npm Module**: Complete NPM wrapper package structure
  - `package.json`: NPM package configuration with cross-platform support
  - `index.js`: Platform detection and binary wrapper with stdio passthrough
  - `postinstall.js`: Automatic binary download from GitHub releases
  - `bin/gorev-mcp`: Executable entry point for NPX usage
- **CI/CD Enhancement**: Multi-stage GitHub Actions workflow
  - Cross-platform binary building (Windows, macOS, Linux)
  - NPM package testing on multiple Node.js versions
  - Automated NPM publishing with artifact management
  - Release automation with GitHub releases

### 🔧 VS Code Extension NPX Integration (v0.6.11)
- **New Server Mode Configuration**: Added `gorev.serverMode` setting ("npx" | "binary")
- **NPX Mode as Default**: Zero-installation setup for users
- **MCP Client Enhancement**: Automatic NPX vs binary mode detection
- **Smart Path Validation**: Server path only required for binary mode
- **Localization Support**: Turkish/English messages for NPX configuration
- **User Experience**: Eliminates need for manual binary installation

### 📚 Documentation
- **Installation Guides Updated**: Both Turkish and English README files updated
  - Added NPX installation as the primary, recommended method
  - Comprehensive MCP configuration examples for Claude Desktop, VS Code, Cursor
  - Platform-specific installation paths and configuration locations
- **CLAUDE.md Enhancement**: Added NPM package development commands and architecture

### 🎯 User Experience
- **Windows Users**: Eliminates complex installation steps and PATH configuration
- **MCP Clients**: Universal compatibility with single configuration format
- **Developers**: Easy testing with `npx @mehmetsenol/gorev-mcp-server@latest --help`
- **CI/CD**: Simple integration without binary management complexity

## [v0.15.3] - 2025-09-18

### 🔧 Fixed
- **VS Code Extension Dependency Display** - Critical fix for dependency visualization in VS Code extension
  - **Root Cause**: MCP handlers (`gorev_listele`, `gorev_detay`) were not including dependency count information in markdown output
  - **Solution**: Enhanced `GorevListele` and `GorevDetay` handlers to include dependency information using `gorevBagimlilikBilgisi` helper
  - **Impact**: VS Code extension now displays 🔒/🔓 icons, dependency counts, and proper task blocking indicators
  - **Fields**: Now properly transmits `bagimli_gorev_sayisi`, `tamamlanmamis_bagimlilik_sayisi`, `bu_goreve_bagimli_sayisi`
- **Compilation Fix** - Resolved missing `log` import in `export_import.go`
  - Fixed build failure due to undefined log package in import/export logging statements

### 🧪 Tests
- **Enhanced Test Coverage** - Added comprehensive dependency parsing tests
  - New MarkdownParser tests for task list dependency parsing
  - New MarkdownParser tests for task detail dependency parsing
  - Validates proper extraction of all dependency count fields

### 🛡️ Quality Assurance
- **Rule 15 Compliance** - Complete root cause analysis without workarounds
- **Architecture Reuse** - Leveraged existing dependency calculation infrastructure
- **Regression Prevention** - Comprehensive test coverage added

## [v0.15.2] - 2025-09-18

### 🔧 Fixed
- **Import/Export System Logging** - Fixed logging inconsistencies in data import/export operations
  - Replaced `fmt.Printf` with `log.Printf` for proper log formatting
  - Added detailed import conflict logging with task IDs and resolution strategies
  - Enhanced error logging for failed task creation during import
  - Fixed task-tag association and dependency creation log messages
- **AI Interaction Error Handling** - Improved file watcher AI interaction error reporting
  - Fixed AI interaction save error message to use proper i18n key `error.interactionSaveFailed`
  - Enhanced error context for file change interaction recording
- **VS Code Extension Duplicate Task Handling** - Resolved duplicate task display issues
  - Added `removeDuplicateTasks()` method to filter duplicate tasks by ID
  - Implemented intelligent duplicate resolution keeping most recently updated tasks
  - Enhanced task loading with automatic duplicate filtering
  - Added debug logging for duplicate detection and removal

### 🛡️ Security & Quality
- **Rule 15 Compliance** - All fixes maintain zero technical debt approach
- **Logging Standardization** - Consistent logging patterns across all modules
- **Error Message Localization** - Proper i18n integration for user-facing error messages

## [v0.15.0] - 2025-09-17

### ✨ Added
- **Advanced Search and Filtering System** - Major new feature
  - FTS5 full-text search with SQLite virtual tables for high-performance searching
  - 6 new MCP tools: `gorev_search_advanced`, `gorev_search_suggestions`, `gorev_search_history`, `gorev_filter_profile_save`, `gorev_filter_profile_load`, `gorev_filter_profile_list`
  - Fuzzy string matching using Levenshtein distance algorithm
  - NLP integration for intelligent query parsing
  - Thread-safe concurrent access with comprehensive error handling
  - Filter profile management system for saved search combinations
  - Search history tracking with analytics
- **VirtualBox Linux VM Setup Scripts** - Complete development environment automation
  - 7 comprehensive setup scripts for Ubuntu/Debian/Fedora/CentOS
  - Modular design with error handling and comprehensive logging
  - Automated Go, Node.js, VS Code installation
  - Project building, testing, and extension compilation
  - Debug tools and troubleshooting helpers
- **Database Migration 000010** - FTS5 search infrastructure
  - `gorevler_fts` virtual table for full-text search
  - `filter_profiles` table for saved filter combinations
  - `search_history` table for search analytics

### 🔧 Fixed
- **MCP Schema Validation Error** - Fixed missing `items` property in `gorev_batch_update` tool array schema
- **Localization Compliance** - Added 30+ new i18n keys for search functionality
- **Thread Safety** - Enhanced SearchEngine and FilterProfileManager with proper mutex protection

### 📚 Changed
- **Tool Count Update** - Updated from 31 to 42 MCP tools in documentation
- **Version Bump** - Updated to v0.15.0 across all project files
- **Documentation** - Enhanced README.md with advanced search features
- **ROADMAP** - Marked Advanced Search as completed

### 🏗️ Technical
- New packages: `internal/gorev/search_engine.go`, `internal/gorev/filter_profile_manager.go`
- Comprehensive test coverage: `search_engine_test.go`, `filter_profile_manager_test.go`, `search_integration_test.go`
- Enhanced error handling with proper i18n message formatting
- SQLite FTS5 integration with content synchronization triggers

## [v0.14.2] - 2025-09-14

### Added
- **Complete Workspace Database Support**: Project-specific database functionality without VS Code dependency
  - Automatic workspace detection: Server detects `.gorev/gorev.db` in current directory and parent directories
  - New `gorev init` CLI command: Initialize workspace (`gorev init`) or global (`gorev init --global`) databases
  - Database path logging: Server logs which database file is being used on startup
  - MCP client agnostic: Works with Claude Desktop, Cursor, Windsurf, and any MCP-compatible client
- **VS Code Extension Enhancements**: Smart database mode selection with visual indicators
  - New configuration: `gorev.databaseMode` replaces `gorev.databasePath` (auto/workspace/global modes)
  - Status bar indicator: Shows current database mode (📁 Workspace / 🌐 Global) with path tooltip
  - New commands: "Initialize Workspace Database" and "Switch Database Mode"
  - Auto-detection: Extension automatically detects workspace databases and switches modes

### Changed
- Database priority logic enhanced with comprehensive fallback system:
  1. `GOREV_DB_PATH` environment variable (any MCP client can set this)
  2. Current directory `.gorev/gorev.db`
  3. Parent directories `.gorev/gorev.db` (monorepo support)
  4. User home `~/.gorev/gorev.db` (global)
  5. Standard fallback locations
- `getDatabasePath()` comment updated to reflect MCP client agnostic nature
- Server version updated to v0.14.2

### Technical Details
- **Files Added**:
  - `gorev-vscode/src/commands/databaseCommands.ts`: New database management commands
- **Files Modified**:
  - `gorev-mcpserver/cmd/gorev/main.go`: Enhanced workspace detection, init command, logging
  - `gorev-vscode/package.json`: New database mode configuration and commands
  - `gorev-vscode/src/mcp/client.ts`: Workspace database detection and mode switching
  - `gorev-vscode/src/ui/statusBar.ts`: Database mode visual indicator
  - `gorev-vscode/src/extension.ts`: Database mode event handling
  - `gorev-vscode/src/utils/constants.ts`: New command constants
  - All localization files: Database mode translations (TR/EN)
- **Documentation**: CLAUDE.md updated with comprehensive v0.14.2 feature documentation

### Usage Scenarios
- **Single Project**: Auto-detects `.gorev/gorev.db`
- **Monorepo**: Each package can have its own `.gorev/` directory
- **Global Mode**: Shared database across all projects
- **VS Code-Free**: Full functionality without VS Code extension

## [v0.14.1] - 2025-09-14

### Added
- **VS Code Extension Database Path Configuration**: Users can now specify a custom database file path in VS Code extension settings
  - New configuration setting: `gorev.databasePath` in VS Code extension settings
  - Server enhancement: Added `GOREV_DB_PATH` environment variable support with priority over automatic detection
  - Automatic directory creation for custom database paths
  - Full localization support (Turkish and English) for new configuration option

### Changed
- `getDatabasePath()` function now prioritizes `GOREV_DB_PATH` environment variable as first choice for database location
- VS Code extension MCP client enhanced to read database path configuration and set appropriate environment variables

### Technical Details
- **Files Modified**:
  - `gorev-mcpserver/cmd/gorev/main.go`: Enhanced database path detection logic with environment variable support
  - `gorev-vscode/package.json`: Added `gorev.databasePath` configuration property
  - `gorev-vscode/src/mcp/client.ts`: Added configuration reading and `GOREV_DB_PATH` environment variable setting
  - `gorev-vscode/l10n/bundle.l10n*.json`: Added `config.databasePath` localization keys
- **Backward Compatibility**: Maintains existing database location logic as fallback when custom path is not configured
- **Rule 15 Compliance**: Complete implementation without workarounds or technical shortcuts

## [v0.14.0] - 2025-09-13
### Added
- Thread safety enhancement with comprehensive mutex protection
- 8 new test suites for 90%+ coverage expansion
- Performance optimizations: 30% faster startup
- Enhanced NLP processor with improved AI interactions
- Auto state manager with better file system integration

### Changed  
- Modernized string handling across all modules
- Standardized error handling patterns with proper context
- Improved resource management and cleanup patterns
- Updated testing infrastructure to eliminate duplicate patterns

### Fixed
- Race conditions in AI Context Manager
- File system integration issues in auto state manager
- NLP processor error handling edge cases
- Auto state manager state transition reliability
- TypeScript compilation errors in VS Code extension

### Performance
- 15-20% memory footprint reduction
- 30% faster application startup
- Optimized database queries and connections
- Enhanced concurrent access patterns

### Security
- 100% production-ready security audit compliance
- Enhanced defensive programming patterns
- Improved error handling without information leakage


# Changelog

All notable changes to the Gorev project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.15.4] - 2025-09-18

### 🚀 Added - NPX Integration
- **NPX Easy Installation**: Zero-config setup via `npx @gorev/mcp-server@latest`
- **NPM Package Wrapper**: @gorev/mcp-server with automatic binary downloads
- **Cross-platform Support**: Windows/macOS/Linux, AMD64/ARM64 architecture support
- **GitHub Actions CI/CD**: Automated NPM publishing pipeline

### 🎨 Added - VS Code Extension NPX Mode
- **serverMode Configuration**: NPX as default mode (npx/binary options)
- **Smart MCP Client**: Automatic command/args selection based on server mode
- **Backward Compatibility**: Full support for existing binary installations
- **Version 0.6.11**: Enhanced NPX integration with localization updates

### 📚 Added - Comprehensive Documentation
- **Turkish Documentation Suite**: Complete Turkish docs (kurulum.md, kullanim.md, mcp-araclari.md)
- **Documentation Standardization**: All docs updated to v0.15.4
- **NPX Installation Guides**: Updated installation instructions for all MCP clients
- **Go Version Standardization**: Updated to Go 1.23+ across all documentation

### 🔧 Changed
- **Installation Methods**: NPX as primary installation method
- **VS Code Extension Default**: NPX mode as default server configuration
- **Documentation Structure**: Enhanced organization with Turkish localization
- **Test Coverage Reporting**: Updated from 84.6% to 90%+ in documentation

### 🐛 Fixed
- **PowerShell Syntax**: Fixed install.ps1 syntax error with VERSION variable
- **Documentation Versions**: Standardized version references across all files
- **MCP Tools Count**: Consistent 48 tools reference across documentation

## [v0.13.1] - 2025-09-04 - Test Infrastructure Standardization

### Added
- **New Testing Package**: `internal/testing/helpers.go` with standardized test infrastructure
  - `TestDatabaseConfig` struct for unified database configuration
  - `SetupTestEnvironmentWithConfig()` centralized helper function
  - `DefaultTestDatabaseConfig()` with sensible defaults
  - Support for memory DB, temp files, custom paths, migrations, templates, i18n

### Changed
- **Test Infrastructure Modernization**: Migrated 50+ duplicate patterns to standardized helpers
  - `handlers_test.go`: 30+ pattern migrations across template and performance tests
  - `integration_test.go`: 11 comprehensive test function migrations  
  - `server_coverage_test.go`: 7 patterns with various database configurations
  - `concurrency_test.go`: Concurrent testing patterns standardized
  - `handlers_edge_cases_test.go`: Complex scenarios including read-only database tests
- **Code Quality**: Net -17 lines despite adding comprehensive infrastructure
- **Import Cleanup**: Removed unused imports from test files

### Fixed
- **i18n Key Corrections**: Fixed incorrect error.* prefixes
  - `parentTaskNotFound` instead of `error.parentTaskNotFound`
  - `circularDependency` instead of `error.circularDependency` 
- **Test Expectations**: Updated format expectations to match current markdown output
- **Tool List**: Added missing `gorev_export` and `gorev_import` tools to server tests
- **Concurrency Tests**: Fixed database setup for concurrent testing scenarios

### Technical Details
- **98% Pattern Elimination**: YeniVeriYonetici duplicates reduced from 50+ to 1
- **40%+ Code Reduction**: Database setup boilerplate eliminated
- **Rule 15 Compliance**: Zero shortcuts, comprehensive solution
- **Zero Breaking Changes**: Internal refactoring with preserved functionality

## [0.12.0] - 2025-08-20 - VS Code Data Export/Import Integration

### Added
- **🎨 VS Code Extension v0.6.0 - Complete Data Export/Import Integration**
  - **4 New Commands**: Export Data, Import Data, Export Current View, Quick Export
  - **Multi-Step Export Dialog**: WebView-based 4-step export configuration wizard
  - **Multi-Step Import Wizard**: WebView-based 4-step import process with conflict resolution
  - **Progress Tracking**: Real-time progress reporting with VS Code progress API
  - **File Format Support**: JSON (structured) and CSV (tabular) export/import formats
  - **Advanced Filtering**: Project-specific filtering, date range support, data type selection
  - **Conflict Resolution**: Skip, overwrite, merge strategies with dry run preview
  - **70+ Localization Strings**: Complete Turkish/English localization for all export/import UI

- **🔧 MCP Server Enhancements**
  - **Template Alias System**: 9 memorable shortcuts (`bug`, `feature`, `research`, etc.)
  - **CLI Command**: `gorev template aliases` for easy discovery
  - **DRY Compliance**: Eliminated 700+ code duplication violations across 7 phases
  - **Enhanced Error Handling**: Improved i18n error messages and validation

### Changed
- **⚡ Performance & Quality Improvements**
  - **Thread-Safety**: Added comprehensive mutex protection to AI Context Manager
  - **Code Architecture**: Refactored handlers.go from 3,060 to 2,362 lines (-23% reduction)
  - **Test Coverage**: Enhanced VS Code extension to 100% test coverage
  - **Documentation Structure**: Optimized CLAUDE.md and enhanced user guides

### Fixed
- **🐛 Critical Fixes**
  - **Race Conditions**: Resolved AI context manager concurrency issues
  - **TypeScript Compilation**: Fixed ESLint and compilation errors in WSL environment
  - **Template Usage**: Complete migration from deprecated `gorev_olustur` to template-based creation

### Breaking Changes
- **🚨 Template System Mandatory**: All task creation must use `templateden_gorev_olustur` with templates
- **Deprecated Tool Removal**: `gorev_olustur` tool completely removed from MCP registry

### Technical Achievements
- **Complete Feature Integration**: VS Code extension now provides full visual interface for all MCP server export/import capabilities
- **Production-Ready UI**: WebView security, proper error handling, user experience optimization
- **Rule 15 Compliance**: Zero technical debt, comprehensive solution without shortcuts
- **Bilingual Excellence**: Consistent Turkish/English support across all components

## [0.11.1] - 2025-07-21 - Documentation Optimization Release

### Added
- **📚 Enhanced Documentation Structure**
  - Added `CLAUDE.en.md` - English version for international developers
  - Added `docs/DEVELOPMENT_HISTORY.md` - Complete project history archive
  - Added `docs/MCP_TOOLS_REFERENCE.md` - Detailed MCP tools documentation
  - Added "AI Assistant Documentation" section to README files

### Changed
- **🎯 CLAUDE.md Token Optimization**
  - Reduced from ~40KB to 8KB for Claude Code token efficiency
  - Restructured to focus on essential development guidance only
  - Moved detailed history and tool documentation to separate files
- **📖 Improved Cross-References**
  - Enhanced README.md and README.en.md with proper documentation links
  - Better navigation between Turkish and English documentation
  - Modular documentation structure for improved maintainability

### Technical Details
- Documentation restructuring improves Claude Code performance by reducing token usage
- Maintains complete information in separate, focused files
- No functional changes to codebase - documentation-only release

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