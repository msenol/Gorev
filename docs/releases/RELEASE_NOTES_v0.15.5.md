# Release Notes v0.15.5

**Release Date:** September 18, 2025
**Version:** v0.15.5

## 🌟 Overview

Release v0.15.5 addresses critical NPX package issues and introduces embedded migrations architecture. This release ensures that the NPX package works seamlessly across all environments without requiring additional configuration, making installation and usage significantly more reliable.

## 🚀 What's New

### Critical NPX Fix

- **Resolved "error.dataManagerInit" Issue**: Fixed critical database initialization failures in NPX package environments
- **Migration File Access**: NPX package can now properly access migration files without external dependencies
- **No Environment Variables Required**: NPX package works without GOREV_ROOT environment variable
- **Enhanced Reliability**: Significantly improved stability for NPX-based installations

### Embedded Migrations Architecture

- **Complete Migration Embedding**: All migration files now embedded directly into Go binary using `//go:embed`
- **New Architecture**: Added `cmd/gorev/migrations_embed.go` with embedded filesystem support
- **Unified API**: `YeniVeriYoneticiWithEmbeddedMigrations()` function for embedded FS support
- **Automatic Extraction**: `migrateDBWithFS()` extracts embedded migrations to temporary directory
- **Fallback System**: `createVeriYonetici()` unified helper with embedded/filesystem fallback
- **Template Integration**: All template commands updated to use embedded migrations

## 📦 Components Updated

### Go Server (gorev-mcpserver)

- **Version**: v0.15.5
- **New Files**: `cmd/gorev/migrations_embed.go`
- **Enhanced**: `internal/gorev/veri_yonetici.go` with embed.FS support
- **Architecture**: Embedded migrations system

### NPM Package (gorev-npm)

- **Version**: 0.15.5
- **Package**: @mehmetsenol/gorev-mcp-server@0.15.5
- **Fixed**: Database initialization in NPX environments
- **Registry**: Published to npm registry

## 🛠️ Installation

### NPX Package (Recommended)

```bash
# Now works reliably without configuration
npx @mehmetsenol/gorev-mcp-server@latest

# Global installation
npm install -g @mehmetsenol/gorev-mcp-server
```

### MCP Client Configuration

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

### Binary Installation

```bash
# Linux
wget https://github.com/msenol/Gorev/releases/download/v0.15.5/gorev-v0.15.5-linux-amd64.tar.gz
tar -xzf gorev-v0.15.5-linux-amd64.tar.gz

# macOS
wget https://github.com/msenol/Gorev/releases/download/v0.15.5/gorev-v0.15.5-darwin-amd64.tar.gz
tar -xzf gorev-v0.15.5-darwin-amd64.tar.gz

# Windows
wget https://github.com/msenol/Gorev/releases/download/v0.15.5/gorev-v0.15.5-windows-amd64.zip
unzip gorev-v0.15.5-windows-amd64.zip
```

## 🔧 Technical Details

### Embedded Migrations System

- **Go Embed**: Uses `//go:embed` directive for file embedding
- **Filesystem Abstraction**: embed.FS interface for embedded file access
- **Temporary Extraction**: Migrations extracted to temporary directory during initialization
- **Backward Compatibility**: Maintains compatibility with existing filesystem migrations
- **Performance Impact**: Slight binary size increase, one-time extraction overhead

### NPX Package Architecture

- **Self-Contained**: No external file dependencies
- **Configuration-Free**: Works without environment variables
- **Error Handling**: Improved error messages and fallback mechanisms
- **Cross-Platform**: Consistent behavior across all supported platforms

### Database Initialization Flow

1. **Check Embedded**: Attempt to use embedded migrations first
2. **Fallback**: If embedded fails, fall back to filesystem migrations
3. **Extract**: Extract embedded migrations to temporary directory
4. **Initialize**: Run database initialization with extracted files
5. **Cleanup**: Temporary files handled automatically

## 🐛 Bug Fixes

### NPX Package Issues

- **Database Initialization**: Fixed "error.dataManagerInit" preventing NPX package startup
- **Migration Access**: Resolved inability to access migration files in NPX environments
- **Environment Dependencies**: Eliminated requirement for GOREV_ROOT environment variable
- **Error Messages**: Improved error reporting for initialization failures

### Migration System

- **File Access**: Reliable migration file access across all deployment methods
- **Path Resolution**: Enhanced path resolution for embedded and filesystem migrations
- **Error Handling**: Better error handling during migration extraction and execution

## 📋 Release Artifacts

### Binaries

- `gorev-linux-amd64` - Linux binary with embedded migrations
- `gorev-darwin-amd64` - macOS Intel binary with embedded migrations
- `gorev-darwin-arm64` - macOS Apple Silicon binary with embedded migrations
- `gorev-windows-amd64.exe` - Windows binary with embedded migrations

### Packages

- `gorev-v0.15.5-linux-amd64.tar.gz` - Linux package
- `gorev-v0.15.5-darwin-amd64.tar.gz` - macOS Intel package
- `gorev-v0.15.5-darwin-arm64.tar.gz` - macOS Apple Silicon package
- `gorev-v0.15.5-windows-amd64.zip` - Windows package

### NPM Package

- `mehmetsenol-gorev-mcp-server-0.15.5.tgz` - NPM package with embedded migrations

### Verification

- `checksums.txt` - SHA256 checksums for all artifacts

## 🔮 What's Next

### Planned Improvements

- **Performance Optimization**: Further optimization of migration extraction process
- **Enhanced Error Handling**: More detailed error reporting and recovery options
- **Documentation Updates**: Expanded NPX usage documentation

### Future Features

- **Hot Migration Updates**: Dynamic migration updates without restart
- **Migration Rollback**: Enhanced migration rollback capabilities
- **Custom Migration Paths**: Support for custom migration directories

## 📝 Upgrade Notes

### For NPX Users

- **Automatic Fix**: The critical NPX issue is automatically resolved upon update
- **No Configuration Required**: Remove any workaround environment variables
- **Immediate Benefits**: Expect significantly more reliable operation

### For Binary Users

- **Minor Size Increase**: Binary size slightly increased due to embedded migrations
- **No Configuration Changes**: Existing configurations continue to work
- **Enhanced Reliability**: Improved migration handling across all deployment methods

### For Developers

- **New API**: `YeniVeriYoneticiWithEmbeddedMigrations()` available for embedded migration support
- **Backward Compatibility**: Existing `YeniVeriYonetici()` continues to work
- **Migration Files**: Migration files now embedded in binary, reducing deployment complexity

## 🤝 Contributing

We welcome contributions! Please see our [contributing guidelines](../docs/development/contributing.md) for more information.

## 📞 Support

For issues and questions:

- **GitHub Issues**: [Create an issue](https://github.com/msenol/Gorev/issues)
- **Documentation**: [Full documentation](../docs/README.md)
- **Community**: Join our community discussions

---

**Previous Release**: [v0.15.4](./RELEASE_NOTES_v0.15.4.md)
**Next Release**: [v0.15.8](./RELEASE_NOTES_v0.15.8.md)

*Last Updated: September 18, 2025*
