# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added - 25 June 2025

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

### Documentation
- Updated `docs/mcp-araclari.md` with detailed documentation for all new tools
- Added examples and response formats for each tool
- Updated future features roadmap

### Technical
- Added `GorevDetayAl`, `ProjeDetayAl`, `GorevDuzenle`, `GorevSil` methods to business logic layer
- Added `ProjeGetir`, `GorevSil`, `ProjeGorevleriGetir` methods to data access layer
- Fixed all integration tests to work with new MCP API
- Added helper function for extracting text from MCP results in tests

## [1.0.0] - 2024-12-15

### Added
- Initial release with core MCP server functionality
- Basic task management tools: `gorev_olustur`, `gorev_listele`, `gorev_guncelle`
- Project management: `proje_olustur`
- System overview: `ozet_goster`
- SQLite database backend
- Clean architecture implementation
- Turkish domain language support