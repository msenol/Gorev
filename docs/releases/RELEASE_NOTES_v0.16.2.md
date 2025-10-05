# Release Notes - v0.16.2

**Release Date**: October 5, 2025
**Type**: Critical Bug Fix Release

## 🐛 Critical Bug Fixes

### NPM Binary Update Issue (CRITICAL)

**Problem**: Users upgrading the NPM package from v0.16.1 or earlier were stuck on the v0.15.24 binary from September 2025, unable to access new features introduced in v0.16.0.

**Root Cause**: The `postinstall.js` script had logic to "skip if binary already exists", which prevented binary updates during package upgrades.

**Impact**:
- All users who upgraded from v0.16.1 or earlier couldn't access v0.16.x features:
  - ❌ REST API integration
  - ❌ Embedded Web UI at http://localhost:5082
  - ❌ VS Code auto-start feature
  - ❌ Multi-workspace support
  - ❌ SHA256-based workspace IDs

**Solution**:
- Modified `postinstall.js` to **ALWAYS** remove old binary before downloading new one
- Removed bundled binaries from NPM package distribution
- Package size dramatically reduced: **78.4 MB → 6.9 KB** 📦
- Binaries now always downloaded from GitHub releases
- Ensures all users get the correct binary version matching package version

**Files Changed**:
- `gorev-npm/postinstall.js` (lines 171-175): Added `safeUnlink(binaryPath)` before download
- `gorev-npm/package.json`: Version bump to 0.16.2
- `gorev-npm/binaries/`: Removed directory (no longer bundling binaries)

## 📚 Documentation Updates

- Updated README.md with v0.16.2 notes and critical bug fix details
- Updated CLAUDE.md with v0.16.2 release information
- Added comprehensive v0.16.2 and v0.16.1 entries to CHANGELOG.md
- Added upgrade notice to gorev-npm/README.md

## 🔧 Upgrade Instructions

**For users upgrading from v0.16.1 or earlier:**

```bash
# 1. Uninstall old version
npm uninstall -g @mehmetsenol/gorev-mcp-server

# 2. Clear NPM cache (ensures fresh download)
npm cache clean --force

# 3. Install latest version
npm install -g @mehmetsenol/gorev-mcp-server@latest

# 4. Verify version (should show v0.16.2)
gorev-mcp version
```

**For fresh installations:**

```bash
# Use npx (recommended - no installation)
npx @mehmetsenol/gorev-mcp-server serve

# Or install globally
npm install -g @mehmetsenol/gorev-mcp-server@latest
```

## 📊 Package Improvements

- **Package Size**: Reduced from 78.4 MB to 6.9 KB (99.99% reduction!)
- **Installation Speed**: Faster package download due to smaller size
- **Binary Management**: More reliable - always gets correct binary from GitHub
- **Disk Usage**: Significant reduction in node_modules size

## 🔗 Related Links

- [GitHub Release v0.16.2](https://github.com/msenol/Gorev/releases/tag/v0.16.2)
- [NPM Package](https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server)
- [Full CHANGELOG](https://github.com/msenol/Gorev/blob/main/CHANGELOG.md#0162---2025-10-05)
- [Migration Guide](../migration/v0.15-to-v0.16.md)

## 🙏 Acknowledgments

Special thanks to all users who reported the binary update issue. Your feedback is invaluable in making Gorev better for everyone!
