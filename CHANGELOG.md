# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0] - 2025-06-25

### Added

#### Active Project Management
- **New feature**: Active project context for task management
- **New MCP tools**:
  - `proje_aktif_yap` - Set a project as the active project
  - `aktif_proje_goster` - Display current active project details
  - `aktif_proje_kaldir` - Remove active project setting
- **Database changes**:
  - Added `aktif_proje` table to store persistent active project setting
  - Table uses CHECK constraint to ensure only one active project (id=1)
- **Enhanced existing tools**:
  - `gorev_olustur` - Now accepts optional `proje_id` parameter; uses active project by default if not specified
  - `gorev_listele` - Added `tum_projeler` boolean parameter; filters by active project by default

### Changed
- **Breaking change**: `GorevOlustur` function now takes 4 parameters (added `projeID`)
- Task creation feedback now includes project name when task is assigned to a project
- Task listing title shows active project name when filtering by active project

### Technical
- Added `veri_yonetici_ext.go` for active project database operations
- Added `AktifProjeAyarla`, `AktifProjeGetir`, `AktifProjeKaldir` methods to VeriYonetici
- Updated VeriYoneticiInterface with active project methods
- Updated IsYonetici to support active project operations
- Enhanced MCP handlers to utilize active project context

## [1.2.0] - 2025-06-25

### Added

#### Unit Testing Infrastructure
- **Comprehensive unit tests** for business logic layer with 88.2% code coverage
- `veri_yonetici_test.go` - Tests for data access layer (VeriYonetici)
  - CRUD operations testing
  - SQL injection protection tests
  - NULL handling tests
  - Concurrent access tests
  - Edge case validation
- `is_yonetici_test.go` - Tests for business logic layer (IsYonetici)
  - Mock implementation of VeriYoneticiInterface
  - Business logic validation
  - Error handling scenarios
  - Partial update logic tests
- `veri_yonetici_interface.go` - Interface for dependency injection and mocking

#### New MCP Tools
- `gorev_detay` - Display detailed task information in markdown format
- `gorev_duzenle` - Edit task title, description, priority, or project assignment
- `gorev_sil` - Delete tasks with confirmation safety
- `proje_listele` - List all projects with task counts
- `proje_gorevleri` - List tasks for a specific project grouped by status

#### Features
- Full markdown support in task descriptions
- Partial update capability for task editing (only specified fields are updated)
- Task count display in project listings
- Status-based grouping in project task views
- Comprehensive integration tests for all new tools

### Changed
- Task descriptions now support full markdown formatting
- Improved error messages to be more user-friendly
- Updated MCP handler signatures to match mark3labs/mcp-go v0.6.0 API
- **Refactored IsYonetici to use VeriYoneticiInterface for better testability**

### Documentation
- Updated `docs/mcp-araclari.md` with detailed documentation for all new tools
- Added examples and response formats for each tool
- Updated future features roadmap

### Technical
- Added `GorevDetayAl`, `ProjeDetayAl`, `GorevDuzenle`, `GorevSil` methods to business logic layer
- Added `ProjeGetir`, `GorevSil`, `ProjeGorevleriGetir` methods to data access layer
- Fixed all integration tests to work with new MCP API
- Added helper function for extracting text from MCP results in tests
- **Implemented dependency injection pattern for better testability**
- **Added table-driven test patterns following Go conventions**
- **Test coverage includes: edge cases, SQL injection, concurrent access, NULL handling**

## [1.0.0] - 2024-12-15

### Added
- Initial release with core MCP server functionality
- Basic task management tools: `gorev_olustur`, `gorev_listele`, `gorev_guncelle`
- Project management: `proje_olustur`
- System overview: `ozet_goster`
- SQLite database backend
- Clean architecture implementation
- Turkish domain language support