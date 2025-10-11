# Changelog

All notable changes to Gorev MCP Server will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.17.0] - 2025-10-11

### BREAKING CHANGES ⚠️

This is a **major refactoring release** that renames all database column names and API field names from Turkish to English. Domain terminology (`gorevler`, `projeler`) remains Turkish as per project convention.

**⚠️ IMPORTANT**: This release includes **automatic database migration (000011)** that will run when you first start Gorev v0.17.0. Your data will be preserved, but the schema will be updated.

### Changed

#### Database Schema (Migration 000011)
- **gorevler table**: `baslik` → `title`, `aciklama` → `description`, `durum` → `status`, `oncelik` → `priority`, `proje_id` → `project_id`, `olusturma_tarih` → `created_at`, `guncelleme_tarih` → `updated_at`, `son_tarih` → `due_date`
- **projeler table**: `isim` → `name`, `tanim` → `definition`, `olusturma_tarih` → `created_at`, `guncelleme_tarih` → `updated_at`
- **etiketler table**: `isim` → `name`
- **gorev_etiketleri table**: `gorev_id` → `task_id`, `etiket_id` → `tag_id`
- **baglantilar table**: `kaynak_id` → `source_id`, `hedef_id` → `target_id`, `baglanti_tip` → `connection_type`
- **gorev_templateleri table**: `isim` → `name`, `tanim` → `definition`, `varsayilan_baslik` → `default_title`, `aciklama_template` → `description_template`, `ornek_degerler` → `sample_values`, `kategori` → `category`, `aktif` → `active`, `alanlar` → `fields`
- **ai_interactions table**: `gorev_id` → `task_id`
- **aktif_proje table**: `proje_id` → `project_id`
- **gorevler_fts table**: Rebuilt with new column names (`task_id`, `title`, `description`, `tags`)
- **gorev_hiyerarsi view**: Recreated with all new English column names

#### Go Backend (55+ files)
- **Core Models** (`modeller.go`): All struct field JSON tags updated to English
  - `Gorev`: `Title`, `Description`, `Status`, `Priority`, `ProjectID`, `ParentID`, `CreatedAt`, `UpdatedAt`, `DueDate`
  - `Proje`: `Name`, `Definition`, `CreatedAt`, `UpdatedAt`
  - `Etiket`: `Name`
  - `Baglanti`: `SourceID`, `TargetID`, `ConnectionType`
  - `GorevTemplate`: `Name`, `Definition`, `DefaultTitle`, `DescriptionTemplate`, `SampleValues`, `Category`, `Active`, `Fields`
  - `AIInteraction`: `TaskID`, `ActionType`, `Context`, `Timestamp`
- **Database Layer** (18 files): All SQL queries updated with new column names
- **MCP & API Layer** (8 files): All parameter handling updated
- **Template Placeholders**: All built-in templates now use English placeholders
  - `{{baslik}}` → `{{title}}`
  - `{{aciklama}}` → `{{description}}`
  - `{{oncelik}}` → `{{priority}}`

#### TypeScript Frontend (20 files)
- **API Types** (`types/api.ts`, `interfaces/gorev.ts`): All interfaces updated
- **API Client** (`api.ts`): All request/response structures updated
- **UI Components**: All data access updated to use English field names
- **React Query Hooks**: Query key and data structures updated

#### Documentation (25+ files)
- All API documentation updated with new field names
- MCP Tools Reference updated
- README files updated across all modules

### Fixed

- **VS Code Extension Server Auto-Start** (Rule 15 compliant fix):
  - `isServerRunning()`: Now checks both port and health endpoint (prevents false positives)
  - `startServer()`: Comprehensive error handling with actionable user messages
  - Timeout increased: 15s → 60s for reliable first-time installation
  - User feedback: "Starting server..." notification with "Show Logs" button
  - Error detection: Package not found, spawn failures, immediate exits
- FTS5 (Full-Text Search) virtual table configuration and triggers
- Template field substitution to work with English field names
- All integration tests (70/70 passing)
- Foreign key constraint handling for `ProjeID` field (empty string → NULL conversion)

### Migration Notes

**Users upgrading from v0.16.x should**:
1. **Backup your database** before upgrading (`.gorev/gorev.db`)
2. Read `docs/MIGRATION_GUIDE_v0.17.md` for detailed migration instructions
3. Update any custom scripts/tools that reference old field names
4. VS Code extension users: Update to v0.17.0 simultaneously

**Migration is automatic but irreversible** - the database schema will be permanently changed.

## [0.16.3] - 2025-10-06

### Added

- `parseQueryFilters()` helper function for advanced search query parsing

### Fixed

- **gorev_bulk**: Complete rewrite with proper parameter transformation for all 3 operations
  - `update` operation: Transforms `{ids: [], data: {}}` → `{updates: [{id, ...fields}]}`
  - `transition` operation: Accepts both `durum` and `yeni_durum` parameters
  - `tag` operation: Accepts both `operation` and `tag_operation` parameters
- **gorev_guncelle**: Extended to support both `durum` and `oncelik` updates simultaneously
- **gorev_search**: Advanced mode now parses query strings like `"durum:X oncelik:Y"` into filters
- **VS Code Tree View**: Dependency counters now display correctly (removed `omitempty` from JSON tags)
  - `bagimli_gorev_sayisi`, `tamamlanmamis_bagimlilik_sayisi`, `bu_goreve_bagimli_sayisi`

### Changed

- Unified MCP handlers now accept multiple parameter format variations for backward compatibility

## [0.16.2] - 2025-10-05

### Fixed

- NPM package binary update mechanism (78.4 MB → 6.9 KB package size)
- VS Code extension auto-start functionality

### Added

- Embedded Web UI (React + TypeScript) at http://localhost:5082
- Multi-workspace support with SHA256-based workspace IDs
- Template aliases: `bug`, `feature`, `research`, `refactor`, `test`, `doc`

## [0.16.1] - 2025-10-04

### Changed

- Version bump and documentation updates

## [0.16.0] - 2025-10-03

### Added

- 24 unified MCP tools (reduced from 45)
- 8 unified handlers: Active Project, Hierarchy, Bulk Ops, Filter Profiles, File Watch, IDE Management, AI Context, Search

### Changed

- Major refactoring of MCP tool architecture for better maintainability

---

**Note:** This changelog was created on 2025-10-06. Previous releases may not have complete entries.
