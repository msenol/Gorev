# Change Log

All notable changes to the "gorev-vscode" extension will be documented in this file.

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