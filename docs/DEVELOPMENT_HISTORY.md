# Gorev Development History

This file contains the detailed development history and release notes for the Gorev project, moved from CLAUDE.md to optimize token usage.

## MCP Server (v0.11.1) - Thread-Safety Enhancement (16 August 2025)
- **AI Context Manager Race Condition Fix**: Comprehensive thread-safety implementation
  - Added `sync.RWMutex` protection to `AIContextYonetici` struct in `internal/gorev/ai_context_yonetici.go`
  - Protected all context operations: SetActiveTask, GetActiveTask, GetContext, saveContext
  - Implemented read-write lock optimization for concurrent access patterns
  - Created internal unsafe methods (getContextUnsafe, saveContextUnsafe) for use within locked sections
  - Zero breaking changes, full backward compatibility maintained
- **Enhanced Testing Infrastructure**:
  - Added comprehensive concurrent access test with 50 goroutines and 500 operations total
  - Implemented race condition detection test using Go race detector (`go test -race`)
  - Verified data integrity under high-concurrency MCP tool usage scenarios
  - Added stress testing patterns for concurrent context operations
- **Thread-Safety Implementation Details**:
  - `SetActiveTask()` and `saveContext()` use exclusive write locks (sync.RWMutex.Lock)
  - `GetActiveTask()`, `GetContext()`, `GetRecentTasks()` use shared read locks (sync.RWMutex.RLock)
  - Internal unsafe methods prevent recursive locking and deadlock scenarios
  - Proper defer patterns ensure locks are always released, even in error conditions
- **Production Readiness**:
  - Resolves data corruption issues in high-concurrency environments with multiple MCP clients
  - Maintains performance with read-write lock optimization for read-heavy workloads
  - Prevents race conditions during simultaneous AI context operations
  - Ready for production deployment with concurrent MCP tool access
- **Documentation Enhancements**:
  - Created `docs/security/thread-safety.md` - Comprehensive thread-safety guidelines
  - Created `docs/development/concurrency-guide.md` - Developer concurrency patterns
  - Updated testing documentation with race condition prevention strategies
- **Rule 15 Compliance**: Comprehensive solution addressing root cause, no workarounds or technical debt
  - Complete thread-safety implementation, not a temporary fix
  - Comprehensive testing covering all concurrent access scenarios
  - Clean abstraction separating safe public methods from internal unsafe operations
  - Production-ready implementation following Go concurrency best practices

## MCP Server (v0.11.0) - Complete Internationalization Support (21 July 2025)
- **Full Bilingual MCP Server**: Implemented complete i18n system for Gorev MCP server
  - Added `go-i18n/v2` library for professional internationalization support
  - Created comprehensive translation system with Turkish (default) and English support
  - **270+ strings converted** from hardcoded Turkish to i18n.T() function calls
  - Language detection hierarchy: CLI flag → GOREV_LANG env → LANG env → Turkish default
  - Dynamic language switching without server restart
- **New Translation Infrastructure**:
  - `internal/i18n/manager.go` - Complete translation management system
  - `locales/tr.json` - 270+ Turkish translation keys with organized structure
  - `locales/en.json` - Complete English translations matching Turkish functionality
  - Template data support for dynamic values using {{.Variable}} syntax
- **Files Internationalized**:
  - `internal/mcp/handlers.go` - All 25 MCP tools with error messages and descriptions
  - `internal/gorev/is_yonetici.go` - Business logic error messages
  - `internal/gorev/template_yonetici.go` - Template system messages
  - `internal/gorev/veri_yonetici.go` - Database operation errors
  - `internal/gorev/ai_context_yonetici.go` - AI context management messages
  - `cmd/gorev/main.go` - CLI command descriptions and help text
  - `cmd/gorev/mcp_commands.go` - Debug command interfaces
- **Language Features**:
  - Automatic language detection from environment
  - `--lang` CLI flag for explicit language selection (tr, en)
  - `GOREV_LANG` environment variable support
  - Fallback to system `LANG` environment variable
  - All error messages, success messages, and UI text translated
- **Compatibility**:
  - Maintains 100% backward compatibility with existing Turkish interfaces
  - No breaking changes to MCP tool APIs
  - Seamless integration with existing VS Code extension
  - Ready for international user adoption

## VS Code Extension (v0.5.0) - Complete Bilingual Support (21 July 2025)
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

## Comprehensive Test Infrastructure Enhancement (17 July 2025)
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

## MCP Server (v0.10.2) - MCP Debug System & TypeScript Fix (17 July 2025)
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

## MCP Server (v0.10.1) - Critical Pagination Fix
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

## MCP Server (v0.10.0) - Template Usage Now Mandatory
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

## Complete Historical Archive

For detailed release history of versions v0.9.0 and earlier, including:
- AI Context Management System (v0.9.0)
- Pagination and Performance fixes (v0.8.1, v0.9.1, v0.9.2)
- Subtask Hierarchy System (v0.8.0)
- Test Infrastructure Development (v0.7.0-beta.1)
- Template System Implementation (v0.5.0-v0.6.0)

Please refer to the CHANGELOG.md file for complete version history and detailed technical changes.