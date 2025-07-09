# Changelog

All notable changes to the Gorev project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## MCP Server

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