# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Last Updated: 25 June 2025

### Recent Changes
- Added markdown support for task descriptions
- Implemented new MCP tools: `gorev_detay`, `gorev_duzenle`, `gorev_sil`, `proje_listele`, `proje_gorevleri`
- Enhanced task editing with partial update capability
- Added comprehensive integration tests for all new tools
- Updated documentation with examples and usage patterns

## Project Overview

Gorev is an MCP (Model Context Protocol) server written in Go that provides task management capabilities to AI assistants. The project was recently converted from Kotlin to Go and uses the community MCP SDK (`mark3labs/mcp-go`).

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

```
cmd/gorev/main.go          → Entry point, CLI commands (cobra)
internal/mcp/              → MCP protocol layer
  ├── handlers.go         → MCP tool implementations
  └── server.go          → MCP server setup
internal/gorev/           → Business logic layer
  ├── modeller.go        → Domain models (Gorev, Proje, Ozet)
  ├── is_yonetici.go     → Business logic orchestration
  └── veri_yonetici.go   → Data access layer (SQLite)
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

The server implements 10 MCP tools:

### Task Management
1. **gorev_olustur**: Create new task (params: baslik, aciklama, oncelik)
2. **gorev_listele**: List tasks (params: durum filter)
3. **gorev_detay**: Show detailed task info in markdown (params: id)
4. **gorev_guncelle**: Update task status (params: id, durum)
5. **gorev_duzenle**: Edit task properties (params: id, baslik?, aciklama?, oncelik?, proje_id?)
6. **gorev_sil**: Delete task (params: id, onay)

### Project Management
7. **proje_olustur**: Create project (params: isim, tanim)
8. **proje_listele**: List all projects with task counts (no params)
9. **proje_gorevleri**: List project tasks grouped by status (params: proje_id)

### Reporting
10. **ozet_goster**: Show summary statistics (no params)

All tools follow the pattern in `internal/mcp/handlers.go` and are registered in `RegisterTools()`. Task descriptions support full markdown formatting.

## Testing Strategy

- **Unit Tests**: Business logic in `internal/gorev/`
- **Integration Tests**: MCP handlers in `test/integration_test.go`
- **Table-Driven Tests**: Go convention for test cases
- **Test Database**: Use `:memory:` SQLite for tests

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

SQLite database with three tables:

- **projeler**: id, isim, tanim, olusturma_tarih, guncelleme_tarih
- **gorevler**: id, baslik, aciklama, durum, oncelik, proje_id, olusturma_tarih, guncelleme_tarih
- **baglantilar**: id, kaynak_id, hedef_id, baglanti_tip (future feature)

Migrations are handled manually in `tablolariOlustur()`.

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

- `internal/gorev/modeller.go`: Domain model definitions
- `internal/mcp/handlers.go`: MCP tool implementations
- `internal/gorev/veri_yonetici.go`: Database operations
- `cmd/gorev/main.go`: CLI and server initialization

## Version Management

Version info is injected at build time via ldflags:
- `main.version`
- `main.buildTime`
- `main.gitCommit`

The Makefile handles this automatically.