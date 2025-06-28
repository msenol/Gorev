# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **Package structure**: Updated module path from `github.com/yourusername/gorev` to `github.com/msenol/gorev`
- **Configuration**: Removed docker-based gorev server configuration from `.mcp.json`

### Fixed
- **Import paths**: Fixed all internal import paths to use the correct module name

### Documentation
- **Enhanced**: Complete documentation overhaul with Claude (Anthropic) assistance
- **Added**: Platform-specific installation guides for Windows, macOS, and Linux
- **Added**: MCP editor integration guides for VS Code, Windsurf, Cursor
- **Added**: Comprehensive examples documentation (ornekler.md)
- **Added**: API reference documentation (api-referans.md)
- **Added**: Developer guide (gelistirme.md)
- **Updated**: Clarified MCP compatibility with multiple AI-enabled editors
- **Fixed**: All placeholder URLs replaced with actual GitHub repository URLs

## [0.5.0] - 2025-06-27

### Added

#### Task Template System (Görev Şablonları)
- **New feature**: Predefined templates for consistent task creation
- **Database changes**: 
  - Added `gorev_templateleri` table via migration 000004
  - Stores template definitions with dynamic fields
- **Four default templates**:
  - **Bug Raporu** - Detailed bug reporting template with environment, steps, expected/actual behavior
  - **Özellik İsteği** - Feature request template with user stories and acceptance criteria
  - **Teknik Borç** - Technical debt template for refactoring tasks
  - **Araştırma Görevi** - Research task template with objectives and evaluation criteria
- **Template features**:
  - Dynamic field types: text, select, date, number
  - Field validation (required/optional)
  - Default values and select options
  - Template placeholders using `{{field_name}}` syntax
  - Automatic tag and priority assignment
- **New MCP tools**:
  - `template_listele` - List available templates with optional category filter
  - `templateden_gorev_olustur` - Create tasks from templates with custom field values
- **CLI commands**:
  - `gorev template list [kategori]` - List templates by category
  - `gorev template show <template-id>` - Show detailed template information
  - `gorev template init` - Initialize default templates
- **Integration**: Templates are automatically created on database initialization

### Technical
- Added `template_yonetici.go` for template management logic
- Added `GorevTemplate` and `TemplateAlan` structs to domain model
- Extended `VeriYoneticiInterface` with 5 new template methods
- Modified `veri_yonetici.go` to auto-initialize templates after migration
- Added template handlers to MCP server with full schema definitions

## [0.4.0] - 2025-06-27

### Added

#### Task Due Dates (Son Tarih)
- **New feature**: Tasks can now have due dates for deadline tracking
- **Database changes**: Added `son_tarih` column to `gorevler` table via migration
- **Enhanced tools**:
  - `gorev_olustur` - Accepts `son_tarih` parameter in YYYY-MM-DD format
  - `gorev_duzenle` - Can edit due dates
  - `gorev_listele` - New sorting options: `son_tarih_asc`, `son_tarih_desc`
  - `gorev_listele` - New filters: `acil` (due in 7 days), `gecmis` (overdue)
  - `gorev_detay` - Displays due dates in task details

#### Task Tagging System (Etiketler)
- **New feature**: Tasks can be categorized with multiple tags
- **Database changes**:
  - Added `etiketler` table for storing tags
  - Added `gorev_etiketleri` table for many-to-many relationships
- **Enhanced tools**:
  - `gorev_olustur` - Accepts `etiketler` parameter (comma-separated tags)
  - `gorev_listele` - Can filter by `etiket` parameter
  - `gorev_detay` - Displays tags in task details
- Tags are automatically created if they don't exist

#### Task Dependencies (Görev Bağımlılıkları)
- **New feature**: Tasks can have dependencies that must be completed before starting
- **New MCP tool**: `gorev_bagimlilik_ekle` - Create dependencies between tasks
- **Business logic**:
  - Tasks cannot be moved to "devam_ediyor" status if dependencies are incomplete
  - New `GorevBagimliMi` function checks dependency satisfaction
- **Enhanced tools**:
  - `gorev_guncelle` - Validates dependencies before status changes
  - `gorev_detay` - Shows dependencies with completion status indicators:
    - ✅ for completed dependencies
    - ⏳ for pending dependencies
    - Warning message if task cannot be started

#### Database Schema Management
- **Implemented** database schema versioning with `golang-migrate/migrate`
- **Schema files** organized in `internal/veri/migrations/`
- **Schema versions**:
  1. `000001_initial_schema.up.sql` - Base tables
  2. `000002_add_due_date_to_gorevler.up.sql` - Due date support
  3. `000003_add_tags.up.sql` - Tagging system

### Changed

#### Breaking Changes
- **`GorevOlustur`** function signature: Now takes 6 parameters (added `sonTarihStr`, `etiketIsimleri`)
- **`GorevListele`** function signature: Now takes 3 parameters (added `sirala`, `filtre`)
- **`VeriYonetici`** constructor: Now requires migrations path parameter
- **Method renames**:
  - `GorevDetayAl` → `GorevGetir`
  - `ProjeDetayAl` → `ProjeGetir`

#### Improvements
- Enhanced task detail view with better dependency visualization
- Improved error messages for dependency violations
- Better test coverage with dependency validation tests

### Technical
- Added `GorevBagimliMi` method to check task dependencies
- Added `BaglantiEkle` and `BaglantilariGetir` methods for dependency management
- Added `EtiketleriGetirVeyaOlustur` and `GorevEtiketleriniAyarla` for tag management
- Updated mock implementations in tests to support new features
- Integration tests updated for new functionality

## [0.3.0] - 2025-06-25

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

## [0.2.0] - 2025-06-25

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

## [0.1.0] - 2024-12-15

### Added
- Initial release with core MCP server functionality
- Basic task management tools: `gorev_olustur`, `gorev_listele`, `gorev_guncelle`
- Project management: `proje_olustur`
- System overview: `ozet_goster`
- SQLite database backend
- Clean architecture implementation
- Turkish domain language support