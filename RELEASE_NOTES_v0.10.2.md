# 🚀 Gorev v0.10.2 Release Notes

**Release Date**: July 17, 2025  
**Version**: v0.10.2  
**Git Commit**: da074a4

## 📋 Overview

This release includes significant improvements to the MCP debug system, enhanced documentation, and comprehensive test data seeding capabilities. The focus has been on improving developer experience with better debugging tools and fixing critical pagination issues.

## 🆕 New Features

### 🛠️ Enhanced MCP Debug System
- **New MCP Debug Commands**: Added comprehensive CLI commands for debugging MCP server functionality
  - `gorev mcp list` - List all available MCP tools
  - `gorev mcp call <tool> <args>` - Direct tool invocation for testing
- **Debug Configuration**: New debug settings for VS Code extension
  - `gorev.debug.useWrapper` - Enable debug wrapper for MCP communication logging
  - `gorev.debug.logPath` - Configure debug log directory
  - `gorev.debug.serverTimeout` - Configurable server timeout
- **Debug Documentation**: Added comprehensive debugging guides
  - `docs/debugging/MCP_DEBUG_GUIDE.md` - Complete debugging guide
  - `docs/debugging/DEBUG_CONNECTION_ISSUE.md` - Connection issue analysis
  - `docs/debugging/VS_CODE_CONNECTION_FIX_SUMMARY.md` - Fix summary

### 📊 Enhanced Test Data Seeding
- **Template-Based Test Data**: Updated test data seeder to use templates instead of deprecated `gorev_olustur`
- **Realistic Test Scenarios**: Added comprehensive test scenarios using all template types
- **Enhanced Debugging**: Better logging and error handling in test data generation
- **Template Mapping**: Added comprehensive template mapping documentation

## 🔧 Improvements

### 🐛 Critical Bug Fixes
- **Fixed Pagination Logic**: Resolved critical pagination bug where subtasks appeared twice
- **Fixed Duplicate Task Display**: Tasks no longer appear both as independent items and under their parent
- **Fixed Infinite Loop**: Resolved infinite loop issue in VS Code when requesting pages beyond available data

### 💻 VS Code Extension (v0.4.6)
- **Enhanced Duplicate Detection**: Added detailed duplicate detection logging with context
- **Show All Projects Toggle**: New configuration option `gorev.treeView.showAllProjects`
- **Keyboard Shortcuts**: Added `Ctrl+Alt+P` / `Cmd+Alt+P` for toggle all projects
- **Improved Markdown Parser**: Better handling of Turkish priority names and emojis
- **Reduced Pagination Default**: Changed default page size from 100 to 10 for better performance

### 🗄️ MCP Server Improvements
- **Enhanced Database Path Resolution**: Improved logic for finding database and migrations
- **GOREV_ROOT Support**: Better support for GOREV_ROOT environment variable
- **Helper Methods**: Added `CallTool` method for direct tool invocation
- **Template Enforcement**: Continued enforcement of template-based task creation

## 📚 Documentation Updates

### 📖 New Documentation
- **Pagination Fix Documentation**: `docs/pagination-fix.md` - Detailed technical explanation
- **Template Mapping Guide**: `docs/test-data-seeder-template-mapping.md` - Complete template usage guide
- **Enhanced MCP Tools Reference**: Updated `docs/user-guide/mcp-tools.md` with latest tools

### 📝 Updated Documentation
- **README.md**: Updated version to v0.10.2 and added latest features
- **CLAUDE.md**: Updated with v0.10.2 changes and technical details
- **CHANGELOG.md**: Comprehensive changelog with all recent improvements

## 🧪 Technical Improvements

### 🔧 Build System
- **Version Management**: Updated build system to v0.10.2
- **Enhanced Makefile**: Added new build targets and improved dependency management
- **Path Resolution**: Improved executable and migration path resolution

### 🗃️ Database & Performance
- **Optimized Pagination**: Rewritten pagination logic for better performance
- **Index Optimization**: Added database indexes for improved query performance
- **Query Optimization**: Bulk operations for N+1 query prevention

## 🐛 Bug Fixes

### Critical Fixes
- **Fixed Subtask Pagination**: Subtasks now correctly appear with their parent regardless of pagination window
- **Fixed Task Count Display**: Now shows root task count instead of total task count
- **Fixed Orphan Task Logic**: Removed problematic orphan task checking that caused duplicates

### Minor Fixes
- **VS Code Extension**: Fixed various UI issues in dark theme
- **TypeScript Compilation**: Resolved compilation errors in markdown parser
- **Package Management**: Fixed package-lock.json handling

## 🔄 Breaking Changes

None in this release. The template system requirement from v0.10.0 continues to be enforced.

## 📈 Performance Improvements

- **Reduced Token Usage**: Optimized response formatting to prevent token limit errors
- **Faster Pagination**: Improved pagination performance with better query optimization
- **Enhanced Caching**: Better caching mechanisms for frequently accessed data

## 🧰 Developer Experience

### 🛠️ Debug Tools
- **MCP Communication Logging**: Debug wrapper logs all MCP communication
- **Enhanced Error Messages**: Better error reporting with context
- **Test Data Generation**: Improved test data seeding with realistic scenarios

### 📋 Testing
- **Enhanced Test Coverage**: Improved test scenarios for pagination and template usage
- **Better Mock Data**: More realistic test data generation
- **Integration Testing**: Enhanced integration test coverage

## 📦 Installation & Upgrade

### New Installation
```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash

# Windows (PowerShell)
irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

### Upgrade from Previous Version
```bash
# Backup your data
cp gorev.db gorev.db.backup

# Download new version
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev

# Test version
./gorev version
```

## 🔗 Links

- **GitHub Release**: https://github.com/msenol/gorev/releases/tag/v0.10.2
- **VS Code Extension**: https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode
- **Documentation**: https://github.com/msenol/gorev/tree/main/docs
- **Issue Tracker**: https://github.com/msenol/gorev/issues

## 🤝 Contributors

Thanks to everyone who contributed to this release through bug reports, feature requests, and testing.

## 🔜 What's Next

- Enhanced AI context management features
- Improved template system with custom templates
- Better integration with more MCP-compatible editors
- Advanced reporting and analytics features

---

**Full Changelog**: https://github.com/msenol/gorev/compare/v0.10.1...v0.10.2