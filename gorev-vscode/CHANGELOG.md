# Change Log

All notable changes to the "gorev-vscode" extension will be documented in this file.

## [0.5.1] - 2025-08-19

### Changed
- **Server Compatibility Update**: Enhanced compatibility with Gorev MCP Server v0.11.1
- **Template Alias Support**: Ready for new template alias system (bug, feature, task, etc.)
- **Performance Optimization**: Optimized for 500x faster server response times
- **Resource Management**: Enhanced MCP client connection stability with FileWatcher improvements

### Fixed
- **Connection Stability**: Improved MCP server connection handling for better reliability  
- **Template System**: Enhanced template selection compatibility with server alias system
- **Error Handling**: Better error messages for server communication issues

### Technical
- **MCP Protocol**: Full compatibility with server's enhanced MCP tool registration
- **FileWatcher Integration**: Ready for automatic file monitoring capabilities
- **Thread Safety**: Compatible with server's thread-safe AI context management
- **i18n Consistency**: Maintains bilingual support with server's i18n enhancements

## [0.5.0] - 2025-07-21

### Added
- **Complete Bilingual Support (English/Turkish)**
  - Automatic language detection based on VS Code language settings (vscode.env.language)
  - Localized all UI strings across 36 source files (500+ translations)
  - Added bundle.l10n.json for runtime localization
  - Added package.nls.json files for VS Code marketplace metadata
  - Turkish README.tr.md for Turkish users
  - All commands, notifications, error messages, and UI elements now support both languages

### Changed
- Migrated from hardcoded strings to VS Code's modern l10n API
- Updated all user-facing strings to use vscode.l10n.t() for dynamic translation
- Enhanced user experience for international users

### Technical Details
- Localized components:
  - Commands (21 commands with localized titles and descriptions)
  - TreeView providers (tasks, projects, templates)
  - UI components (filter toolbar, status bar, task detail panel)
  - Inline editing and drag-drop operations
  - Debug tools and test data seeders
- Translation files:
  - `l10n/bundle.l10n.json` - English base strings
  - `l10n/bundle.l10n.tr.json` - Turkish translations
  - `package.nls.json` - English package metadata
  - `package.nls.tr.json` - Turkish package metadata

## [0.4.0] - 2025-07-10

### ⚠️ BREAKING CHANGES

- **Template Usage is Now Mandatory**: Direct task creation via `gorev_olustur` is no longer supported. All tasks must be created using templates.
- The "Create Task" (Ctrl+Shift+G) command now opens the template wizard
- The "Quick Create Task" command now opens the quick template selection

### Changed

- Updated to support MCP server v0.10.0 with mandatory template requirement
- Modified task creation commands to redirect to template selection
- `gorev.createTask` command now executes `gorev.openTemplateWizard`
- `gorev.quickCreateTask` command now executes `gorev.quickCreateFromTemplate`

### Migration Guide

To create tasks in v0.4.0:
1. Use the template wizard (Ctrl+Shift+G or "Create Task" command)
2. Select a template from the available options
3. Fill in the required fields for the template
4. The task will be created with consistent structure

Available templates include:
- Bug Report v2 (detailed bug tracking)
- Spike Research (time-boxed investigations)
- Performance Issue (performance optimization tasks)
- Security Fix (security vulnerability fixes)
- Refactoring (code quality improvements)
- And more standard templates...

## [0.3.9] - 2025-07-10

### Fixed
- Pagination logic in MCP server v0.9.1 that was causing incomplete task lists
- VS Code extension now correctly displays all tasks when there are many subtasks

### Changed
- Updated to work with MCP server v0.9.1 pagination improvements

## [0.3.8] - 2025-07-09

### Fixed
- Task count display issue where only 38 tasks were shown instead of all 147 tasks
- MCP server pagination logic that was counting all tasks but only paginating root tasks
- Empty response for second page (offset 100+) when fetching tasks

### Changed
- Updated MCP server pagination to correctly handle all tasks (root and subtasks)
- Improved task fetching to display complete task hierarchy across all pages

## [0.3.7] - 2025-07-09

### Fixed
- Task list not showing any tasks due to parser not recognizing new status emoji (🔄)
- Subtask hierarchy not being preserved in TreeView
- Multiline task descriptions not being parsed correctly
- Parser format detection not recognizing compact format v0.8.1+

### Changed
- Enhanced MarkdownParser to support all task status emojis (⏳, 🚀, ✅, ✓, 🔄)
- Improved compact format parser with proper hierarchy support
- Better handling of multiline descriptions where ID appears on a separate line

## [0.3.6] - 2025-07-09

### Added
- Screenshot support in package.json for VS Code marketplace
- Gallery banner configuration for better marketplace appearance

### Changed
- Updated extension description to highlight key features

## [0.3.5] - 2025-07-09

### Added
- Pagination support integration with MCP server v0.9.0
- Configuration option `gorev.pagination.pageSize` (default: 100)
- Automatic token limit prevention

### Changed
- Updated to work with MCP server v0.9.0 AI features
- Improved performance with large task lists

### Fixed
- Token limit errors when displaying many tasks
- Performance issues with large projects

## [0.3.4] - 2025-07-05

### Added
- Visual progress bars for parent tasks showing subtask completion
- Priority badges (🔥⚡ℹ️) with color coding
- Smart due date formatting (Today, Tomorrow, 3d left, etc.)
- Enhanced dependency badges with lock/unlock status (🔒🔓🔗)
- Colored pill badges for tags with configurable limit
- Rich markdown tooltips with progress visualization
- TaskDecorationProvider for managing visual decorations
- 9 new configuration options for visual preferences

### Fixed
- Dependency data transmission between MCP handlers and UI
- MarkdownParser to correctly parse dependency information

## [0.3.3] - 2025-06-30

### Fixed
- Fixed circular progress chart percentage not being visible in task detail panel
- Implemented CSS overlay solution with absolute positioning for percentage text
- Progress percentage now displays correctly in the center of the circular progress indicator
- Enhanced hierarchy progress parsing with more flexible pattern matching
- Added fallback calculation when server response doesn't include percentage
- Fixed potential NaN values in progress display with proper validation

## [0.3.2] - 2025-06-30

### Added
- Enhanced TreeView with grouping, multi-select, and priority-based color coding
- Drag & Drop support for moving tasks, changing status, and creating dependencies
- Inline editing with F2/double-click, context menus, and date pickers
- Advanced filtering toolbar with search, filters, and saved profiles
- Rich task detail panel with markdown editor and dependency visualization
- Template wizard UI with multi-step interface and dynamic forms
- Hierarchical task display with tree structure
- Progress indicators showing subtask completion percentage

### Fixed
- Filter state persistence issue
- Tag display when tasks created via CLI
- Project task count showing as 0 in TreeView
- Task detail panel UI issues in dark theme
- Single-click task selection in TreeView

## [0.3.0] - 2025-06-29

### Added
- Subtask hierarchy support with unlimited depth
- Visual hierarchy indicators in TreeView
- Progress tracking for parent tasks

## [0.2.0] - 2025-06-28

### Added
- Task templates support
- Tag filtering in TreeView
- Due date indicators

## [0.1.0] - 2025-06-27

### Initial Release
- Basic TreeView for tasks and projects
- MCP client integration
- Status bar integration
- Command palette commands