# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- Filter state persistence issue - users can now clear filters without restarting VS Code
  - Added `clearFilters()` method to EnhancedGorevTreeProvider
  - Updated FilterToolbar to properly reset filter state
  - Added keyboard shortcut `Ctrl+Alt+R` / `Cmd+Alt+R` to clear filters
  - Fixed filter update logic to detect and handle empty filter objects

## [0.7.0-beta.1] - 2025-06-30

### Added
- **Test Infrastructure**:
  - Comprehensive edge case testing for all 16 MCP tools
  - Template unit tests with 100% handler coverage
  - Concurrent operation tests for thread safety
  - Testing framework evaluation and decision documentation
  - MCP package coverage increased from 75.1% to 81.5%
  - Created `handlers_edge_cases_test.go` (600+ LOC)
  - Created `template_yonetici_test.go` (400+ LOC)
  - Added `docs/testing-framework-decision.md`
  - Added `docs/test-coverage-phase1-summary.md`

- **VS Code Extension**:
  - Enhanced TreeView with grouping, multi-select, and priority-based color coding
  - Drag & Drop support for moving tasks, changing status, and creating dependencies
  - Inline editing with F2/double-click, context menus, and date pickers
  - Advanced filtering toolbar with search, filters, and saved profiles
  - Rich task detail panel with markdown editor and dependency visualization
  - Template wizard UI with multi-step interface and dynamic forms
  - Comprehensive test suite achieving 50.9% file coverage (19 files tested)

### Fixed
- Database migration issues with `etiketler` table in tests
- Concurrent access test failures (switched to file-based database)
- TypeScript compilation error in `markdownParser.ts`
- Tag display in VS Code UI when tasks created via CLI
- Project task count showing as 0 in TreeView
- Task detail panel UI issues in dark theme

### Changed
- Testing framework decision: Continue with testify (152x faster than ginkgo)
- Enhanced `GorevListele` and `ProjeGorevleri` handlers to include tags and due dates

### Discovered
- Input validation gaps (whitespace-only titles accepted)
- Missing enum validation for priority and status values
- These issues are documented for future improvements

## [0.6.0] - 2025-06-29

### Added
- Task Template System with predefined templates (Bug Report, Feature Request, Technical Debt, Research)
- Task Dependencies with validation
- Due Dates with filtering for urgent/overdue tasks
- Tagging System with multiple tags per task
- Database Schema Management using golang-migrate
- New MCP tools: `gorev_bagimlilik_ekle`, `template_listele`, `templateden_gorev_olustur`
- Enhanced filtering and sorting parameters for `gorev_listele`

### Changed
- GorevOlustur now accepts 6 parameters (added sonTarih, etiketler)
- GorevListele now accepts 3 parameters (added sirala, filtre)
- VeriYonetici constructor requires migrations path

## [0.5.0] - 2025-06-28

### Added
- Initial MCP server implementation with 13 core tools
- VS Code extension with TreeView providers
- SQLite database storage
- Project management features
- Active project context
- Basic task CRUD operations

### Features
- Task management (create, list, update, delete)
- Project management
- Task status tracking (beklemede, devam_ediyor, tamamlandi)
- Priority levels (dusuk, orta, yuksek)
- Summary statistics

---

*For detailed documentation of all features, see the [docs/](docs/) directory.*