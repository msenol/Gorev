# Release Notes - v0.16.3

**Release Date**: October 6, 2025
**Type**: Critical MCP Tool Fix + Documentation Update

## 🔧 Critical Fixes

### 1. gorev_bulk - Complete Parameter Transformation Rewrite

**Problem**: All 3 operations (update/transition/tag) were failing with parameter errors
**Root Cause**: Unified schema → operation-specific format transformation was missing

**Solution**:

- **`update` operation**: Properly transforms `{ids: [], data: {}}` → `{updates: [{id, ...fields}]}`
  - Each ID in the array gets the same data fields applied
  - Array of update objects created automatically
- **`transition` operation**: Accepts both `durum` and `yeni_durum` parameter names for flexibility
  - Backward compatible with existing code
  - Optional parameters: `force`, `check_dependencies`, `dry_run`
- **`tag` operation**: Accepts both `operation` and `tag_operation` parameter names
  - Supports `add`, `remove`, `replace` operations
  - Tags array properly forwarded to handler

**Files Changed**: `internal/mcp/handlers.go:3311-3417` (+106 lines)
**Test Result**: 100% success rate (2/2 update, 1/1 transition, 2/2 tag operations passed)

### 2. gorev_guncelle - Priority Update Support

**Problem**: Only `durum` (status) updates were supported, `oncelik` (priority) parameter was ignored
**Root Cause**: Handler only validated and processed status field

**Solution**:

- Extended validation to accept both `durum` and `oncelik` as optional parameters
- At least one parameter required for operation
- Calls appropriate update methods based on provided parameters
- Success message shows all updated fields

**Files Changed**: `internal/mcp/handlers.go:639-684` (+40 lines)
**Test Result**: 100% success (durum only ✓, oncelik only ✓, both simultaneously ✓)

### 3. gorev_search (advanced mode) - Query Filter Parsing

**Problem**: Filter queries like `"durum:devam_ediyor oncelik:yuksek"` were not being parsed
**Root Cause**: No query string parsing logic for key:value patterns

**Solution**:

- Added `parseQueryFilters()` helper function
- Parses space-separated key:value pairs from query string
- Automatically extracts filters into proper filter map
- Works seamlessly with existing filter parameter

**Files Changed**: `internal/mcp/handlers.go:3571-3593` (+23 lines)
**Test Result**: 100% success (8 results for single filter, 4 results for multi-filter, 21 results for tag filter)

### 4. VS Code Tree View - Dependency Counter Display

**Problem**: Dependency counters (🔒/🔓/🔗 icons) were not visible in VS Code tree view
**Root Cause**: `omitempty` JSON tags excluded 0 values from JSON response, causing frontend checks to fail

**Solution**:

- Removed `omitempty` from 3 dependency counter fields:
  - `bagimli_gorev_sayisi` (number of dependencies)
  - `tamamlanmamis_bagimlilik_sayisi` (incomplete dependencies)
  - `bu_goreve_bagimli_sayisi` (tasks depending on this)
- Fields now always present in JSON, even with 0 values
- VS Code tree provider now correctly displays dependency indicators

**Files Changed**: `internal/gorev/modeller.go:21-23` (3 fields)
**Test Result**: Dependency indicators now display correctly in all scenarios

## 📊 Performance Metrics

- **Bulk operations**: 11-33ms processing time across all operation types
- **Advanced search**: 6-67ms with FTS5 full-text search and relevance scoring
- **Success rate**: 100% (validated by Kilocode AI comprehensive test report)

## 📚 Documentation Updates

### Major Documentation Additions

**Daemon Architecture Documentation** (v0.16.0 feature, previously undocumented):

- Comprehensive daemon mode documentation added to README.md
- Background daemon process model explained
- Lock file mechanism (`~/.gorev-daemon/.lock`) documented
- Multi-client MCP proxy architecture detailed
- WebSocket server real-time update infrastructure
- VS Code extension auto-start integration flow

**Updated Files**:

- `README.md` + `README.tr.md`: Added daemon architecture section + v0.16.3 features
- `CLAUDE.md`: Added daemon architecture notes and v0.16.3 updates
- `gorev-mcpserver/CHANGELOG.md`: Added v0.16.3 release section
- Version bumped across all package.json files and Makefile

### Documentation Pending (Faz 4 & 5)

- `docs/api/MCP_TOOLS_REFERENCE.md`: Complete rewrite planned (41 → 24 tools)
- `gorev-mcpserver/docs/mcp-araclari.md`: Turkish version rewrite planned
- `docs/architecture/daemon-architecture.md`: New comprehensive technical guide planned

## 🎯 Testing & Validation

**Kilocode AI Test Report Results**:

- gorev_bulk: 5/5 operations passed (update, transition, tag)
- gorev_guncelle: 3/3 test scenarios passed (durum, oncelik, both)
- gorev_search: 3/3 query types passed (single filter, multi-filter, tag filter)
- VS Code tree view: Visual confirmation of dependency indicators

**Production Validation**:

- Ubuntu test environment: All fixes confirmed working
- Performance: Sub-100ms response times for all operations
- Backward compatibility: No breaking changes introduced

## 🔗 Links

- [GitHub Release v0.16.3](https://github.com/msenol/Gorev/releases/tag/v0.16.3)
- [CHANGELOG](../../gorev-mcpserver/CHANGELOG.md#0163---2025-10-06)
- [NPM Package](https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server)
- [Migration Guide](../migration/v0.15-to-v0.16.md)

## 🔄 Upgrade Instructions

**For NPM Package Users**:

```bash
# Update globally installed package
npm update -g @mehmetsenol/gorev-mcp-server

# Or use npx (always gets latest)
npx @mehmetsenol/gorev-mcp-server serve

# Verify version
gorev-mcp version  # Should show v0.16.3
```

**For VS Code Extension Users**:

- Extension will automatically use updated MCP server
- Restart VS Code to ensure daemon restarts with new version
- Check daemon version: `curl http://localhost:5082/api/health`

**For Development**:

```bash
# Pull latest changes
git pull origin main

# Rebuild
make build

# Verify version
./gorev version  # Should show v0.16.3
```

## 🙏 Acknowledgments

Special thanks to:

- **Kilocode AI**: Comprehensive testing that identified all 4 critical issues
- **Community testers**: Ubuntu environment validation
- **Contributors**: Bug reports and feature requests

This release brings Gorev MCP server to 100% operational status for all core MCP tools!

---

**Note**: This is a critical bug fix release. All users are strongly encouraged to upgrade to benefit from the MCP tool fixes and improved reliability.
