# Release Notes v0.15.3

**Release Date:** September 18, 2025
**Version:** v0.15.3

## 🌟 Overview

Release v0.15.3 addresses critical VS Code extension issues related to dependency visualization and provides essential compilation fixes. This release ensures that task dependencies are properly displayed in the VS Code extension and resolves build issues for improved development experience.

## 🚀 What's New

### VS Code Extension Dependency Display Fix

- **Critical Fix**: Resolved dependency visualization issues in VS Code extension
- **Enhanced MCP Handlers**: Updated `gorev_listele` and `gorev_detay` handlers to include comprehensive dependency information
- **Visual Indicators**: VS Code extension now properly displays 🔒/🔓 icons and dependency counts
- **Task Blocking**: Proper task blocking indicators for dependent tasks
- **Dependency Fields**: Complete transmission of dependency information:
  - `bagimli_gorev_sayisi` - Number of tasks this task depends on
  - `tamamlanmamis_bagimlilik_sayisi` - Number of incomplete dependencies
  - `bu_goreve_bagimli_sayisi` - Number of tasks depending on this task

### Compilation Improvements

- **Build Fix**: Resolved missing `log` import in `export_import.go`
- **Clean Build**: Eliminated build failures due to undefined packages
- **Logging Consistency**: Proper log formatting throughout the codebase

## 📦 Components Updated

### Go Server (gorev-mcpserver)

- **Version**: v0.15.3
- **Fixed**: MCP handlers for proper dependency information transmission
- **Enhanced**: `gorevBagimlilikBilgisi` helper function integration
- **Resolved**: Compilation issues in export/import functionality

### VS Code Extension

- **Dependency Display**: Now properly shows dependency information
- **Visual Indicators**: Enhanced task status visualization
- **Bug Resolution**: Dependency parsing and display issues resolved

## 🛠️ Installation

### Binary Installation

```bash
# Linux
wget https://github.com/msenol/Gorev/releases/download/v0.15.3/gorev-v0.15.3-linux-amd64.tar.gz
tar -xzf gorev-v0.15.3-linux-amd64.tar.gz

# macOS
wget https://github.com/msenol/Gorev/releases/download/v0.15.3/gorev-v0.15.3-darwin-amd64.tar.gz
tar -xzf gorev-v0.15.3-darwin-amd64.tar.gz

# Windows
wget https://github.com/msenol/Gorev/releases/download/v0.15.3/gorev-v0.15.3-windows-amd64.zip
unzip gorev-v0.15.3-windows-amd64.zip
```

### MCP Client Configuration

```json
{
  "mcpServers": {
    "gorev": {
      "command": "/path/to/gorev",
      "args": ["serve"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

## 🔧 Technical Details

### Dependency Information Enhancement

- **Root Cause**: MCP handlers were not transmitting dependency count information
- **Solution**: Enhanced handlers to use `gorevBagimlilikBilgisi` helper function
- **Architecture**: Leveraged existing dependency calculation infrastructure
- **Impact**: Complete dependency information now available in VS Code extension

### MCP Handler Updates

- **GorevListele Handler**: Now includes dependency information in task listings
- **GorevDetay Handler**: Enhanced task detail view with comprehensive dependency data
- **Data Transmission**: Proper formatting of dependency fields for VS Code extension consumption

### Build System Improvements

- **Import Resolution**: Fixed missing log package imports
- **Compilation**: Clean compilation across all modules
- **Error Handling**: Improved error logging in import/export operations

## 🧪 Testing Enhancements

### Comprehensive Test Coverage

- **MarkdownParser Tests**: New tests for task list dependency parsing
- **Dependency Parsing**: Validation of proper extraction of dependency count fields
- **Task Detail Tests**: New tests for task detail dependency parsing
- **Regression Prevention**: Comprehensive test coverage to prevent future issues

### Quality Assurance

- **Rule 15 Compliance**: Complete root cause analysis without workarounds
- **Architecture Reuse**: Leveraged existing dependency calculation infrastructure
- **Test-Driven Fixes**: Tests added to validate fixes and prevent regressions

## 🐛 Bug Fixes

### VS Code Extension Issues

- **Dependency Display**: Fixed missing dependency information in task views
- **Visual Indicators**: Resolved issues with 🔒/🔓 dependency icons
- **Task Blocking**: Fixed task blocking indicators for dependent tasks
- **Data Parsing**: Enhanced markdown parsing for dependency information

### Compilation Issues

- **Missing Imports**: Fixed undefined `log` package in export_import.go
- **Build Failures**: Resolved compilation errors preventing successful builds
- **Logging**: Consistent log formatting throughout the application

### Data Transmission

- **MCP Handlers**: Fixed incomplete dependency data transmission
- **Field Mapping**: Proper mapping of dependency fields to VS Code extension
- **Information Loss**: Prevented loss of dependency information during MCP communication

## 📋 Release Artifacts

### Binaries

- `gorev-linux-amd64` - Linux binary with dependency fixes
- `gorev-darwin-amd64` - macOS Intel binary with dependency fixes
- `gorev-darwin-arm64` - macOS Apple Silicon binary with dependency fixes
- `gorev-windows-amd64.exe` - Windows binary with dependency fixes

### Packages

- `gorev-v0.15.3-linux-amd64.tar.gz` - Linux package
- `gorev-v0.15.3-darwin-amd64.tar.gz` - macOS Intel package
- `gorev-v0.15.3-darwin-arm64.tar.gz` - macOS Apple Silicon package
- `gorev-v0.15.3-windows-amd64.zip` - Windows package

### Verification

- `checksums.txt` - SHA256 checksums for all artifacts

## 🔮 What's Next

### Planned Improvements

- **Enhanced Dependency Features**: More advanced dependency management
- **VS Code Extension**: Additional visual improvements
- **Performance**: Optimization of dependency calculations

### Future Features

- **Dependency Graphs**: Visual dependency graph representations
- **Smart Dependencies**: Automatic dependency detection
- **Dependency Templates**: Template-based dependency creation

## 📝 Upgrade Notes

### For VS Code Extension Users

- **Immediate Benefits**: Dependency information now properly displayed
- **Visual Improvements**: Enhanced task status indicators
- **Dependency Management**: Better understanding of task relationships

### For Developers

- **Build Issues Resolved**: Clean compilation without import errors
- **Testing**: Enhanced test coverage for dependency features
- **Architecture**: Improved dependency information handling

### For Users

- **No Configuration Changes**: Existing configurations continue to work
- **Enhanced Features**: Better dependency visibility in VS Code
- **Stability**: Improved build stability and reliability

## 🤝 Contributing

We welcome contributions! Please see our [contributing guidelines](../docs/development/contributing.md) for more information.

## 📞 Support

For issues and questions:

- **GitHub Issues**: [Create an issue](https://github.com/msenol/Gorev/issues)
- **Documentation**: [Full documentation](../docs/README.md)
- **Community**: Join our community discussions

---

**Previous Release**: [v0.15.2](./RELEASE_NOTES_v0.15.2.md)
**Next Release**: [v0.15.4](./RELEASE_NOTES_v0.15.4.md)

*Last Updated: September 18, 2025*
