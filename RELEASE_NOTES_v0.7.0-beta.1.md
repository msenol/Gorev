# 🚀 Gorev v0.7.0-beta.1 Release Notes

**Release Date:** June 29, 2025  
**Type:** Beta Release

## 🎉 Highlights

This beta release introduces major enhancements to the VS Code extension UI, numerous bug fixes, and comprehensive documentation updates. The release also includes a detailed roadmap for future development.

## ✨ What's New

### VS Code Extension - Enhanced UI Features

#### 🎯 Enhanced TreeView
- Professional task management with advanced grouping options
- Multi-select support (Ctrl/Cmd+Click)
- Priority-based color coding
- Quick completion checkboxes
- Expandable/collapsible categories
- Task count badges and due date warnings

#### 🖱️ Drag & Drop Support
- Intuitive task management with visual feedback
- Move tasks between projects
- Change status by dragging
- Create dependencies by dropping tasks on each other
- Ghost images and drop zone indicators

#### ✏️ Inline Editing
- F2 or double-click to edit task titles
- Context menu for quick status/priority changes
- Inline date picker
- Escape to cancel, Enter to save

#### 🔍 Advanced Filtering Toolbar
- Real-time search with debouncing
- Multi-criteria filtering (status, priority, tags, dates)
- Saved filter profiles
- Quick filter buttons (Today, This Week, Overdue, Critical)
- Filter result count display

#### 📋 Rich Task Detail Panel
- Split-view markdown editor with live preview
- Dependency visualization
- Activity timeline
- Template field indicators

#### 🧙 Template Wizard UI
- Multi-step interface for template-based task creation
- Template search and filtering
- Dynamic form generation with validation
- Preview before creation

#### 🧪 Comprehensive Test Suite
- Unit tests for all major components
- Integration tests for extension features
- End-to-end workflow tests
- Test fixtures and helpers
- Coverage reporting with c8

### MCP Server Improvements

- **Path Resolution**: Fixed database and migration paths to work from any directory
- **Enhanced Handlers**: `GorevListele` and `ProjeGorevleri` now include tags and due dates

## 🐛 Bug Fixes

### VS Code Extension
- ✅ Fixed tag display in TreeView when tasks created via CLI
- ✅ Fixed project task count showing as 0 in TreeView
- ✅ Fixed task detail panel UI issues in dark theme:
  - Action buttons now visible with proper styling
  - Markdown editor toolbar displays correctly
  - CSP-compliant event handlers
  - Edit/Delete functionality restored
- ✅ Fixed single-click task selection in TreeView
- ✅ Removed non-functional dependency graph feature

### MCP Server
- ✅ Fixed gorev command execution from different directories
- ✅ Fixed TypeScript errors with Turkish filter property names

## 📚 Documentation Updates

### New Documentation
- **TASKS.md**: Added comprehensive roadmap with 11 active development tasks:
  1. DevOps Pipeline and CI/CD automation
  2. Test coverage improvements (target 95%)
  3. UI/UX and accessibility enhancements
  4. Multi-user system and authorization
  5. Performance optimizations
  6. External service integrations (GitHub, Jira, Slack)
  7. Analytics dashboard
  8. Advanced search and filtering
  9. Dependency visualization in TreeView
  10. Subtask system with unlimited hierarchy
  11. **AI-Powered task enrichment system** (NEW)

### Updated Documentation
- CHANGELOG.md: Updated to v0.7.0-beta.1
- CLAUDE.md: Updated version and added recent changes
- README.md: Fixed version info and test coverage discrepancy

## 📊 Technical Details

- **New Files**: 20+ new TypeScript files for enhanced UI
- **Test Infrastructure**: Complete test setup with mocha, sinon, and VS Code test APIs
- **Markdown Parser**: Comprehensive parser for all MCP response types
- **Debug Support**: Test data seeder and debug commands
- **Icons**: Custom SVG icons for tasks, priorities, and templates

## 🔧 Breaking Changes

None in this release.

## 📦 Installation

### MCP Server
```bash
# Download the latest release
wget https://github.com/msenol/gorev/releases/download/v0.7.0-beta.1/gorev-linux-amd64
chmod +x gorev-linux-amd64
sudo mv gorev-linux-amd64 /usr/local/bin/gorev

# Or build from source
git clone https://github.com/msenol/gorev.git
cd gorev/gorev-mcpserver
make build
```

### VS Code Extension
1. Download `gorev-vscode-0.3.0.vsix` from the release page
2. In VS Code: `Extensions` → `...` → `Install from VSIX...`
3. Select the downloaded file

## 🚀 What's Next

We're actively working on the features outlined in TASKS.md. Key priorities include:
- AI-powered task enrichment with API key support
- Advanced search with full-text capabilities
- Performance optimizations for large datasets
- Multi-user support and authorization

## 🙏 Acknowledgments

Special thanks to all contributors and testers who helped make this release possible. This documentation was enhanced with the assistance of Claude (Anthropic).

## 📝 Feedback

We'd love to hear your feedback! Please report issues or feature requests at:
https://github.com/msenol/gorev/issues

---

**Note:** This is a beta release. While it's stable for daily use, some features are still under development. Please report any issues you encounter.