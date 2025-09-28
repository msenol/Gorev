# Release Notes v0.15.8

**Release Date:** September 20, 2025  
**Version:** v0.15.8

## 🌟 Overview

Release v0.15.8 focuses on version standardization across all components, enhanced cross-platform support, and improved build automation. This release ensures consistent versioning and provides comprehensive platform coverage including Windows support and multi-architecture macOS binaries.

## 🚀 What's New

### Version Standardization

- **Unified Versioning**: All components now consistently use v0.15.8
- **Build System**: Enhanced build automation with improved reliability
- **Release Process**: Streamlined artifact generation and publishing workflow

### Enhanced Platform Support

- **Windows Support**: Complete Windows binary support with .exe packaging
- **Multi-Architecture macOS**: Native support for both Intel (amd64) and Apple Silicon (arm64) Macs
- **Cross-Platform Binaries**: Optimized binaries for all major platforms

### Build & Release Improvements

- **Automated Checksums**: SHA256 checksum generation for all release artifacts
- **Fixed Build Scripts**: Resolved truncated build script issues
- **Improved Reliability**: Enhanced build process with better error handling

## 📦 Components Updated

### Go Server (gorev-mcpserver)

- **Version**: v0.15.8
- **Build**: Enhanced cross-platform compilation
- **Binaries**: Linux, macOS (Intel/ARM), Windows

### VS Code Extension (gorev-vscode)

- **Version**: 0.15.8
- **Package**: gorev-vscode-0.15.8.vsix
- **Marketplace**: Published to VS Code Marketplace

### NPM Package (gorev-npm)

- **Version**: 0.15.8
- **Package**: @mehmetsenol/gorev-mcp-server@0.15.8
- **Registry**: Published to npm registry

## 🛠️ Installation

### VS Code Extension

```bash
# Install from marketplace
# Search for "Gorev" in VS Code extensions

# Or install from .vsix file
code --install-extension gorev-vscode-0.15.8.vsix
```

### NPM Package

```bash
# Install globally
npm install -g @mehmetsenol/gorev-mcp-server

# Or use with npx
npx @mehmetsenol/gorev-mcp-server
```

### Binaries

```bash
# Linux
wget https://github.com/msenol/Gorev/releases/download/v0.15.8/gorev-v0.15.8-linux-amd64.tar.gz
tar -xzf gorev-v0.15.8-linux-amd64.tar.gz

# macOS Intel
wget https://github.com/msenol/Gorev/releases/download/v0.15.8/gorev-v0.15.8-darwin-amd64.tar.gz
tar -xzf gorev-v0.15.8-darwin-amd64.tar.gz

# macOS Apple Silicon
wget https://github.com/msenol/Gorev/releases/download/v0.15.8/gorev-v0.15.8-darwin-arm64.tar.gz
tar -xzf gorev-v0.15.8-darwin-arm64.tar.gz

# Windows
wget https://github.com/msenol/Gorev/releases/download/v0.15.8/gorev-v0.15.8-windows-amd64.zip
unzip gorev-v0.15.8-windows-amd64.zip
```

## 🔧 Technical Details

### Build System Changes

- **Fixed Build Scripts**: Resolved issues with truncated build automation
- **Enhanced Cross-Platform**: Improved multi-architecture build process
- **Automated Artifacts**: Streamlined release artifact generation

### Version Management

- **Consistent Versioning**: All components now use synchronized version numbers
- **Automated Updates**: Version updates are now automated across all components
- **Release Tagging**: Proper git tagging and release management

### Platform Support Matrix

| Platform | Architecture | Binary | Status |
|----------|--------------|--------|--------|
| Linux | AMD64 | gorev-linux-amd64 | ✅ |
| macOS | Intel (AMD64) | gorev-darwin-amd64 | ✅ |
| macOS | Apple Silicon (ARM64) | gorev-darwin-arm64 | ✅ |
| Windows | AMD64 | gorev-windows-amd64.exe | ✅ |

## 📋 Release Artifacts

### Binaries

- `gorev-linux-amd64` - Linux binary
- `gorev-darwin-amd64` - macOS Intel binary
- `gorev-darwin-arm64` - macOS Apple Silicon binary
- `gorev-windows-amd64.exe` - Windows binary

### Packages

- `gorev-v0.15.8-linux-amd64.tar.gz` - Linux package
- `gorev-v0.15.8-darwin-amd64.tar.gz` - macOS Intel package
- `gorev-v0.15.8-darwin-arm64.tar.gz` - macOS Apple Silicon package
- `gorev-v0.15.8-windows-amd64.zip` - Windows package

### Extension & NPM

- `gorev-vscode-0.15.8.vsix` - VS Code extension package
- `mehmetsenol-gorev-mcp-server-0.15.8.tgz` - NPM package

### Verification

- `checksums.txt` - SHA256 checksums for all artifacts

## 🐛 Bug Fixes

### Build System

- **Fixed Truncated Scripts**: Resolved issues with build scripts getting truncated
- **Improved Reliability**: Enhanced build process with better error handling
- **Version Consistency**: Fixed version mismatches across components

### Release Process

- **Automated Checksums**: Added automatic checksum generation
- **Streamlined Publishing**: Improved release artifact publishing workflow
- **Enhanced Documentation**: Updated documentation with new version information

## 🔮 What's Next

### Planned Features

- **Enhanced Documentation**: Further documentation improvements
- **Performance Optimizations**: Ongoing performance enhancements
- **Additional Platform Support**: Potential support for additional platforms

### Maintenance

- **Continuous Integration**: Ongoing CI/CD improvements
- **Testing**: Expanded test coverage
- **Documentation**: Regular documentation updates

## 📝 Upgrade Notes

### For Users

- **No Breaking Changes**: This is a maintenance release with no breaking changes
- **Recommended Update**: All users are encouraged to update to v0.15.8
- **Enhanced Platform Support**: Better support for all major platforms

### For Developers

- **Build System**: Updated build system with improved reliability
- **Version Management**: Streamlined version management across components
- **Release Process**: Enhanced release automation and artifact generation

## 🤝 Contributing

We welcome contributions! Please see our [contributing guidelines](../docs/development/contributing.md) for more information.

## 📞 Support

For issues and questions:

- **GitHub Issues**: [Create an issue](https://github.com/msenol/Gorev/issues)
- **Documentation**: [Full documentation](../docs/README.md)
- **Community**: Join our community discussions

---

**Previous Release**: [v0.15.5](./RELEASE_NOTES_v0.15.5.md)  
**Next Release**: TBD

*Last Updated: September 20, 2025*
