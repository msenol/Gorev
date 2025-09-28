# Release Notes v0.15.4

**Release Date:** September 18, 2025
**Version:** v0.15.4

## 🌟 Overview

Release v0.15.4 introduces the revolutionary NPX Easy Installation System, making Gorev accessible to users without complex installation procedures. This release includes a complete NPM package distribution system, VS Code extension NPX integration, and comprehensive CI/CD automation for seamless user experience.

## 🚀 What's New

### NPX Easy Installation System

- **@mehmetsenol/gorev-mcp-server Package**: Brand new NPM package enabling `npx @mehmetsenol/gorev-mcp-server@latest` usage
- **Cross-Platform Binary Support**: Automatic binary download for Windows, macOS, Linux (amd64/arm64)
- **Zero Installation Setup**: Users can run Gorev without manual binary installation steps
- **Simple MCP Configuration**: Easy addition to `mcp.json` with `"command": "npx", "args": ["@mehmetsenol/gorev-mcp-server@latest"]`
- **GitHub Actions Pipeline**: Automated NPM publishing with multi-platform binary builds
- **Platform Detection**: Intelligent platform and architecture detection for correct binary selection
- **Fallback Mechanisms**: Robust error handling and fallback to latest releases

### VS Code Extension NPX Integration (v0.6.11)

- **New Server Mode Configuration**: Added `gorev.serverMode` setting ("npx" | "binary")
- **NPX Mode as Default**: Zero-installation setup for users
- **MCP Client Enhancement**: Automatic NPX vs binary mode detection
- **Smart Path Validation**: Server path only required for binary mode
- **Localization Support**: Turkish/English messages for NPX configuration
- **User Experience**: Eliminates need for manual binary installation

## 📦 Components Updated

### Go Server (gorev-mcpserver)

- **Version**: v0.15.4
- **NPM Integration**: Complete NPM wrapper functionality
- **Cross-Platform**: Enhanced multi-platform binary support

### NPM Package (gorev-npm) - NEW MODULE

- **Version**: 0.15.4
- **Package**: @mehmetsenol/gorev-mcp-server@0.15.4
- **Files**:
  - `package.json`: NPM package configuration with cross-platform support
  - `index.js`: Platform detection and binary wrapper with stdio passthrough
  - `postinstall.js`: Automatic binary download from GitHub releases
  - `bin/gorev-mcp`: Executable entry point for NPX usage

### VS Code Extension (gorev-vscode)

- **Version**: 0.6.11
- **NPX Integration**: Complete NPX mode support
- **Configuration**: New server mode settings

## 🛠️ Installation

### NPX Package (NEW - Recommended)

```bash
# Instant usage without installation
npx @mehmetsenol/gorev-mcp-server@latest

# Global installation
npm install -g @mehmetsenol/gorev-mcp-server

# Version-specific usage
npx @mehmetsenol/gorev-mcp-server@0.15.4
```

### MCP Client Configuration

#### Claude Desktop

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

#### VS Code

```json
{
  "servers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"]
    }
  }
}
```

#### Cursor

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"]
    }
  }
}
```

### Binary Installation (Traditional)

```bash
# Linux
wget https://github.com/msenol/Gorev/releases/download/v0.15.4/gorev-v0.15.4-linux-amd64.tar.gz
tar -xzf gorev-v0.15.4-linux-amd64.tar.gz

# macOS
wget https://github.com/msenol/Gorev/releases/download/v0.15.4/gorev-v0.15.4-darwin-amd64.tar.gz
tar -xzf gorev-v0.15.4-darwin-amd64.tar.gz

# Windows
wget https://github.com/msenol/Gorev/releases/download/v0.15.4/gorev-v0.15.4-windows-amd64.zip
unzip gorev-v0.15.4-windows-amd64.zip
```

## 🔧 Technical Details

### NPM Package Architecture

- **Platform Detection**: Automatically detects operating system and architecture
- **Binary Download**: Downloads appropriate binary from GitHub releases during installation
- **Stdio Passthrough**: Transparent communication between MCP client and Gorev server
- **Error Handling**: Comprehensive error handling with helpful messages
- **Caching**: Downloaded binaries cached for performance

### CI/CD Enhancement

- **Multi-Stage Pipeline**: Enhanced GitHub Actions workflow
  - Cross-platform binary building (Windows, macOS, Linux)
  - NPM package testing on multiple Node.js versions (16, 18, 20)
  - Automated NPM publishing with artifact management
  - Release automation with GitHub releases
- **Quality Assurance**: Automated testing across platforms and Node.js versions
- **Release Automation**: Streamlined release process with automatic NPM publishing

### VS Code Extension Architecture

- **Server Mode Setting**: `gorev.serverMode` configuration option
- **Dynamic Configuration**: Automatic detection and configuration of server mode
- **Path Validation**: Smart validation based on selected server mode
- **Localization**: Complete Turkish/English localization for NPX features

### Platform Support Matrix

| Platform | Architecture | Binary | NPM Support | Status |
|----------|--------------|--------|-------------|--------|
| Linux | AMD64 | gorev-linux-amd64 | ✅ | ✅ |
| Linux | ARM64 | gorev-linux-arm64 | ✅ | ✅ |
| macOS | Intel (AMD64) | gorev-darwin-amd64 | ✅ | ✅ |
| macOS | Apple Silicon (ARM64) | gorev-darwin-arm64 | ✅ | ✅ |
| Windows | AMD64 | gorev-windows-amd64.exe | ✅ | ✅ |

## 🎯 User Experience Improvements

### For Windows Users

- **Eliminates Complex Installation**: No more PATH configuration or manual binary management
- **Single Command Setup**: `npx @mehmetsenol/gorev-mcp-server@latest` works immediately
- **Automatic Updates**: Always uses latest version when using `@latest`

### For MCP Clients

- **Universal Configuration**: Same configuration format works across all MCP clients
- **Zero Dependencies**: No need to install Go, Git, or other development tools
- **Instant Availability**: Ready to use in seconds

### For Developers

- **Easy Testing**: `npx @mehmetsenol/gorev-mcp-server@latest --help` for immediate testing
- **CI/CD Integration**: Simple integration without binary management complexity
- **Version Pinning**: Support for specific version usage

## 📋 Release Artifacts

### Binaries

- `gorev-linux-amd64` - Linux binary
- `gorev-linux-arm64` - Linux ARM64 binary
- `gorev-darwin-amd64` - macOS Intel binary
- `gorev-darwin-arm64` - macOS Apple Silicon binary
- `gorev-windows-amd64.exe` - Windows binary

### Packages

- `gorev-v0.15.4-linux-amd64.tar.gz` - Linux package
- `gorev-v0.15.4-linux-arm64.tar.gz` - Linux ARM64 package
- `gorev-v0.15.4-darwin-amd64.tar.gz` - macOS Intel package
- `gorev-v0.15.4-darwin-arm64.tar.gz` - macOS Apple Silicon package
- `gorev-v0.15.4-windows-amd64.zip` - Windows package

### NPM Package

- `mehmetsenol-gorev-mcp-server-0.15.4.tgz` - NPM package with cross-platform support

### VS Code Extension

- `gorev-vscode-0.6.11.vsix` - VS Code extension with NPX integration

### Verification

- `checksums.txt` - SHA256 checksums for all artifacts

## 🐛 Bug Fixes

### Installation Issues

- **Windows PATH Problems**: Eliminated need for manual PATH configuration
- **Permission Issues**: NPX handles binary permissions automatically
- **Platform Detection**: Robust platform and architecture detection
- **Download Failures**: Enhanced error handling with retry mechanisms

### VS Code Extension

- **Server Configuration**: Streamlined server configuration with NPX mode
- **Path Validation**: Improved path validation logic
- **Error Messages**: Better error messages for configuration issues

## 🔮 What's Next

### Planned Improvements

- **Binary Caching**: Enhanced caching for faster subsequent runs
- **Offline Support**: Offline mode support for cached binaries
- **Update Notifications**: Automatic update notifications in VS Code extension

### Future Features

- **Docker Support**: Docker image distribution
- **Package Managers**: Support for additional package managers (Homebrew, Chocolatey)
- **Auto-Update**: Automatic binary updates

## 📝 Upgrade Notes

### For New Users

- **Recommended**: Use NPX installation for easiest setup
- **No Installation Required**: `npx @mehmetsenol/gorev-mcp-server@latest` works immediately
- **Configuration**: Update MCP client configuration to use NPX

### For Existing Users

- **Backward Compatibility**: Binary installation continues to work
- **Migration**: Optionally migrate to NPX for easier maintenance
- **VS Code Extension**: Update to v0.6.11 for NPX support

### For Developers

- **NPM Development**: New `gorev-npm` module for NPM package development
- **Testing**: Use NPX for testing across different environments
- **CI/CD**: Enhanced GitHub Actions pipeline for automated releases

## 🤝 Contributing

We welcome contributions! Please see our [contributing guidelines](../docs/development/contributing.md) for more information.

## 📞 Support

For issues and questions:

- **GitHub Issues**: [Create an issue](https://github.com/msenol/Gorev/issues)
- **Documentation**: [Full documentation](../docs/README.md)
- **Community**: Join our community discussions

---

**Previous Release**: [v0.15.3](./RELEASE_NOTES_v0.15.3.md)
**Next Release**: [v0.15.5](./RELEASE_NOTES_v0.15.5.md)

*Last Updated: September 18, 2025*
