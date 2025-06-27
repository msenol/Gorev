# CLAUDE.md

This file provides guidance to AI assistants using MCP (Model Context Protocol) when working with code in this repository. Compatible with Claude Code, VS Code with MCP extension, Windsurf, Cursor, and other MCP-enabled editors.

## Last Updated: 27 June 2025

> 🤖 **Documentation Note**: This comprehensive technical guide was enhanced and structured with the assistance of Claude (Anthropic), demonstrating the power of AI-assisted documentation in modern software development.

### Recent Changes
- **Added Task Template System** - Predefined templates for bug reports, feature requests, technical debt, and research tasks
- **Added Task Dependencies** - Tasks can now have dependencies that must be completed before starting
- **Added Due Dates** - Tasks can have deadlines with filtering for urgent/overdue tasks
- **Added Tagging System** - Tasks can be categorized with multiple tags
- **Database Schema Management** - Using golang-migrate for version control
- **Enhanced gorev_listele** - Added sorting (sirala) and filtering (filtre, etiket) parameters
- **Enhanced gorev_olustur** - Now accepts son_tarih (due date) and etiketler (tags) parameters
- **Enhanced gorev_detay** - Shows dependencies with completion status indicators
- **New MCP tools**: 
  - `gorev_bagimlilik_ekle` for creating task dependencies
  - `template_listele` for listing available task templates
  - `templateden_gorev_olustur` for creating tasks from templates
- **New CLI commands**:
  - `gorev template list [kategori]` - List templates by category
  - `gorev template show <template-id>` - Show template details
  - `gorev template init` - Initialize default templates
- **Breaking changes**: 
  - GorevOlustur now takes 6 parameters (added sonTarihStr, etiketIsimleri)
  - GorevListele now takes 3 parameters (added sirala, filtre)
  - VeriYonetici constructor requires migrations path

## Project Overview

Gorev is an MCP (Model Context Protocol) server written in Go that provides task management capabilities to AI assistants across all MCP-compatible editors (Claude Desktop, VS Code, Windsurf, Cursor, Zed, etc.). The project uses the community MCP SDK (`mark3labs/mcp-go`).

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
cmd/gorev/main.go                  → Entry point, CLI commands (cobra)
internal/mcp/                      → MCP protocol layer
  ├── handlers.go                 → MCP tool implementations
  └── server.go                  → MCP server setup
internal/gorev/                   → Business logic layer
  ├── modeller.go                → Domain models (Gorev, Proje, Ozet)
  ├── is_yonetici.go             → Business logic orchestration
  ├── is_yonetici_test.go        → Business logic unit tests
  ├── veri_yonetici.go           → Data access layer (SQLite)
  ├── veri_yonetici_test.go      → Data access layer unit tests
  └── veri_yonetici_interface.go → Interface for dependency injection
```

### Key Design Decisions

1. **Turkish Domain Language**: Core domain concepts use Turkish terms (gorev=task, proje=project, durum=status, oncelik=priority)
2. **MCP SDK Integration**: Uses `mark3labs/mcp-go` v0.6.0 for MCP protocol implementation
3. **SQLite Storage**: Single-file database for simplicity and portability
4. **No External State**: Each MCP request is stateless, no session management

## Development Commands

```bash
# Build
make build                 # Build for current platform
make build-all            # Build for all platforms (linux, darwin, windows)

# Test
make test                 # Run all tests with coverage
make test-coverage        # Generate HTML coverage report
go test -race ./...       # Run with race detector

# Code Quality
make fmt                  # Format code with gofmt
make lint                 # Run golangci-lint (must be installed)
go vet ./...             # Run go vet

# Dependencies
make deps                 # Download and tidy dependencies

# Docker
make docker-build         # Build Docker image
make docker-run          # Run Docker container

# Development
make run                 # Build and run server
./gorev serve --debug    # Run with debug logging
```

## MCP Tools

The server implements 16 MCP tools:

### Task Management
1. **gorev_olustur**: Create new task (params: baslik, aciklama, oncelik, proje_id?, son_tarih?, etiketler?)
   - proje_id is optional; if not provided, uses active project
   - son_tarih: optional due date in YYYY-MM-DD format
   - etiketler: optional comma-separated tags
2. **gorev_listele**: List tasks (params: durum?, tum_projeler?, sirala?, filtre?, etiket?)
   - tum_projeler: if false/omitted, shows only active project tasks
   - sirala: son_tarih_asc, son_tarih_desc
   - filtre: acil (due in 7 days), gecmis (overdue)
   - etiket: filter by tag name
3. **gorev_detay**: Show detailed task info in markdown (params: id)
   - Shows due dates, tags, and dependencies with status indicators
4. **gorev_guncelle**: Update task status (params: id, durum)
   - Validates dependencies before allowing "devam_ediyor" status
5. **gorev_duzenle**: Edit task properties (params: id, baslik?, aciklama?, oncelik?, proje_id?, son_tarih?)
6. **gorev_sil**: Delete task (params: id, onay)
7. **gorev_bagimlilik_ekle**: Create task dependency (params: kaynak_id, hedef_id, baglanti_tipi)

### Task Templates
8. **template_listele**: List available templates (params: kategori?)
   - Shows predefined templates for consistent task creation
9. **templateden_gorev_olustur**: Create task from template (params: template_id, degerler)
   - degerler is an object with field values for the template

### Project Management
10. **proje_olustur**: Create project (params: isim, tanim)
11. **proje_listele**: List all projects with task counts (no params)
12. **proje_gorevleri**: List project tasks grouped by status (params: proje_id)
13. **proje_aktif_yap**: Set active project (params: proje_id)
14. **aktif_proje_goster**: Show current active project (no params)
15. **aktif_proje_kaldir**: Remove active project setting (no params)

### Reporting
16. **ozet_goster**: Show summary statistics (no params)

All tools follow the pattern in `internal/mcp/handlers.go` and are registered in `RegisterTools()`. Task descriptions support full markdown formatting.

## Testing Strategy

- **Unit Tests**: Business logic in `internal/gorev/` (88.2% coverage)
  - `veri_yonetici_test.go`: Data layer tests with SQL injection and concurrent access tests
  - `is_yonetici_test.go`: Business logic tests with mocked dependencies
- **Integration Tests**: MCP handlers in `test/integration_test.go`
- **Table-Driven Tests**: Go convention for test cases
- **Test Database**: Use `:memory:` SQLite for tests
- **Coverage Goal**: Maintain >80% code coverage

Example test pattern:
```go
func TestGorevOlustur(t *testing.T) {
    testCases := []struct {
        name    string
        input   map[string]interface{}
        wantErr bool
    }{
        // test cases...
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Adding New MCP Tools

1. Add handler method to `internal/mcp/handlers.go`
2. Register tool in `RegisterTools()` with proper schema
3. Add integration tests in `test/integration_test.go`
4. Update `docs/mcp-araclari.md` with tool documentation

## Database Schema

SQLite database with seven tables:

- **projeler**: id, isim, tanim, olusturma_tarih, guncelleme_tarih
- **gorevler**: id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih, son_tarih
- **baglantilar**: id, kaynak_id, hedef_id, baglanti_tip (for task dependencies)
- **aktif_proje**: id (CHECK id=1), proje_id (stores single active project)
- **etiketler**: id, isim (tags)
- **gorev_etiketleri**: gorev_id, etiket_id (many-to-many relationship)
- **gorev_templateleri**: id, isim, tanim, varsayilan_baslik, aciklama_template, alanlar, ornek_degerler, kategori, aktif (task templates)

Migrations are handled by golang-migrate in `internal/veri/migrations/`.

## Error Handling

- Always return explicit errors, never panic
- MCP errors use `mcp.NewToolResultError()`
- Turkish error messages for user-facing errors
- Wrap errors with context: `fmt.Errorf("context: %w", err)`

## Code Style

- Keep Turkish terms for domain concepts
- Use English for technical terms and comments
- Follow Go idioms and conventions
- Prefer composition over inheritance
- Keep functions small and focused

## Important Files

- `internal/gorev/modeller.go`: Domain model definitions (includes GorevTemplate, TemplateAlan)
- `internal/mcp/handlers.go`: MCP tool implementations (includes template handlers)
- `internal/gorev/veri_yonetici.go`: Database operations
- `internal/gorev/template_yonetici.go`: Template management operations
- `cmd/gorev/main.go`: CLI and server initialization (includes template commands)

## Version Management

Version info is injected at build time via ldflags:
- `main.version`
- `main.buildTime`
- `main.gitCommit`

The Makefile handles this automatically.