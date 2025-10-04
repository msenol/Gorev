# Contributing to Gorev

> **Version**: This documentation is valid for v0.16.0+
> **Last Updated**: October 4, 2025

This document explains the development environment setup, code standards, and contribution processes for those
who want to contribute to the Gorev project.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Project Structure](#project-structure)
- [Code Standards](#code-standards)
- [Writing Tests](#writing-tests)
- [Adding New Features](#adding-new-features)
- [Adding MCP Tools](#adding-mcp-tools)
- [Debugging](#debugging)
- [VS Code Extension Development](#vs-code-extension-development)
- [Contributing Process](#contributing-process)

## Development Environment Setup

### Requirements

- Go 1.23+ or higher
- Git
- Make (optional, for Makefile usage)
- golangci-lint (for code quality)
- Docker (optional, for container tests)

### Installation Steps

```bash
# Clone the project
git clone https://github.com/msenol/gorev.git
cd gorev/gorev-mcpserver

# Download dependencies
make deps
# or
go mod download

# Build the project
make build
# or
go build -o gorev cmd/gorev/main.go

# Run tests
make test
# or
go test ./...
```

### IDE Settings

#### VS Code

`.vscode/settings.json`:

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--fast"
  ],
  "go.testFlags": ["-v"],
  "go.testTimeout": "30s"
}
```

#### GoLand/IntelliJ

- Go Modules support: Enable
- GOROOT: System Go installation
- Run gofmt on save: Enable

## Project Structure

```text
gorev/
├── gorev-mcpserver/             # MCP server project
│   ├── cmd/
│   │   └── gorev/
│   │       └── main.go          # Main application entry point
│   ├── internal/
│   │   ├── gorev/               # Domain logic
│   │   │   ├── modeller.go      # Data models
│   │   │   ├── is_yonetici.go   # Business logic
│   │   │   ├── veri_yonetici.go # Data access layer
│   │   │   ├── template_yonetici.go # Template management
│   │   │   └── *_test.go        # Unit tests
│   │   ├── mcp/                 # MCP protocol layer
│   │   │   ├── server.go        # MCP server
│   │   │   └── handlers.go      # Tool handlers
│   │   └── i18n/                # Internationalization
│   ├── migrations/              # Database migrations
│   └── test/                    # Integration tests
├── gorev-vscode/                # VS Code extension
├── docs/                        # Documentation
└── scripts/                     # Helper scripts
```

### Package Descriptions

- **cmd/gorev**: CLI commands and server startup
- **internal/gorev**: Core business logic and domain models
- **internal/mcp**: MCP protocol implementation
- **internal/i18n**: Internationalization support (Turkish/English)
- **migrations**: SQL migration files (golang-migrate format)

## Code Standards

### General Rules

1. **Follow Go idioms**: Read Effective Go and Go Code Review Comments
2. **Turkish domain terms**: Use Turkish for domain terms like Görev, Proje, Durum
3. **English technical terms**: Use English for code comments and technical terms
4. **Error handling**: Return explicit errors, don't use panic

### Naming Conventions

```go
// Domain models - Turkish
type Gorev struct { ... }
type Proje struct { ... }

// Interfaces - Turkish + -ci/-ici suffix
type VeriYonetici interface { ... }
type IsYonetici interface { ... }

// Method names - Turkish verb + English object (if needed)
func (v *veriYonetici) GorevOlustur(...) { ... }
func (v *veriYonetici) ProjeListele(...) { ... }

// Constants - UPPER_SNAKE_CASE
const VERITABANI_VERSIYON = "1.2.0"

// Private variables - camelCase
var aktifProjeID int
```

### Code Style

```go
// Good: Short and clear functions
func (v *veriYonetici) GorevSil(id int) error {
    result, err := v.db.Exec("DELETE FROM gorevler WHERE id = ?", id)
    if err != nil {
        return fmt.Errorf("görev silinirken hata: %w", err)
    }
    
    rows, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("etkilenen satır sayısı alınamadı: %w", err)
    }
    
    if rows == 0 {
        return fmt.Errorf("görev bulunamadı: %d", id)
    }
    
    return nil
}
```

### Error Messages

```go
// Turkish user messages (translated via i18n system)
return fmt.Errorf("görev bulunamadı: %d", id)
return fmt.Errorf("geçersiz durum değeri: %s", durum)

// Context wrapping
if err != nil {
    return fmt.Errorf("veritabanı bağlantısı kurulamadı: %w", err)
}
```

## Writing Tests

### Unit Test Structure

```go
func TestGorevOlustur(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    defer db.Close()
    
    veriYonetici := &veriYonetici{db: db}
    
    testCases := []struct {
        name    string
        baslik  string
        oncelik string
        wantErr bool
    }{
        {
            name:    "successful creation",
            baslik:  "Test task",
            oncelik: "orta",
            wantErr: false,
        },
        {
            name:    "empty title",
            baslik:  "",
            oncelik: "orta",
            wantErr: true,
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Act
            _, err := veriYonetici.GorevOlustur(tc.baslik, "", tc.oncelik, nil, "", "")
            
            // Assert
            if tc.wantErr && err == nil {
                t.Error("expected error but got nil")
            }
            if !tc.wantErr && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
```

## Adding New Features

### 1. Update Domain Model

```go
// internal/gorev/modeller.go
type Gorev struct {
    // Existing fields...
    NewField string `json:"new_field,omitempty"`
}
```

### 2. Add Migration

```sql
-- migrations/004_new_field.up.sql
ALTER TABLE gorevler ADD COLUMN new_field TEXT DEFAULT '';

-- migrations/004_new_field.down.sql
ALTER TABLE gorevler DROP COLUMN new_field;
```

### 3. Update Data Layer

```go
// internal/gorev/veri_yonetici.go
func (v *veriYonetici) gorevleriTara(rows *sql.Rows) ([]Gorev, error) {
    // Add new field to scan
    err := rows.Scan(
        &gorev.ID,
        // other fields...
        &gorev.NewField,
    )
}
```

## Adding MCP Tools

### 1. Write Handler Function

```go
// internal/mcp/handlers.go
func (h *Handler) handleNewTool(args map[string]interface{}) (*ToolResult, error) {
    // Parse parameters
    param1, ok := args["param1"].(string)
    if !ok {
        return nil, fmt.Errorf("param1 is required")
    }
    
    // Call business logic
    result, err := h.isYonetici.NewOperation(param1)
    if err != nil {
        return nil, mcp.NewToolResultError(err.Error())
    }
    
    // Return result
    return &ToolResult{
        Content: []Content{{
            Type: "text",
            Text: fmt.Sprintf("✅ Operation successful: %v", result),
        }},
    }, nil
}
```

### 2. Register Tool

```go
// internal/mcp/handlers.go - RegisterTools()
tools = append(tools, Tool{
    Name:        "new_tool",
    Description: "Performs new operation",
    InputSchema: InputSchema{
        Type: "object",
        Properties: map[string]Property{
            "param1": {
                Type:        "string",
                Description: "Parameter description",
            },
        },
        Required: []string{"param1"},
    },
})
```

## Debugging

### Debug Mode

```bash
# Enable debug logs
./gorev serve --debug

# Or environment variable
DEBUG=true ./gorev serve
```

### Logging

```go
import "log/slog"

// Debug log
slog.Debug("operation started", "id", gorevID, "status", durum)

// Error log
slog.Error("database error", "error", err)
```

## VS Code Extension Development

### Extension Setup

```bash
cd gorev-vscode

# Install dependencies
npm install

# Compile TypeScript
npm run compile

# Watch mode (for development)
npm run watch
```

### Testing Extension

1. Open `gorev-vscode` folder in VS Code
2. Press F5 (or Run > Start Debugging)
3. New VS Code window will open (Extension Development Host)
4. Test the extension

## Contributing Process

### Pull Request Workflow

1. **Open Issue**: First open an issue explaining what you want to do
2. **Fork & Branch**: Fork the project and create a feature branch

   ```bash
   git checkout -b feature/new-feature
   ```

3. **Write Code**: Develop according to code standards
4. **Write Tests**: Target 80%+ coverage
5. **Commit**: Use meaningful commit messages

   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve bug"
   git commit -m "docs: update documentation"
   ```

6. **Push & PR**: Push branch and open PR

### Commit Message Format

```text
<type>(<scope>): <subject>

<body>

<footer>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting, missing semi-colons, etc.
- `refactor`: Code refactoring
- `test`: Adding/fixing tests
- `chore`: Maintenance

### Code Review Checklist

- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] Code follows standards
- [ ] No breaking changes (or documented)
- [ ] Performance implications considered

## Common Issues

### SQLite Locked Error

```go
// Solution: Use WAL mode
db.Exec("PRAGMA journal_mode=WAL")
```

### Import Cycle

```go
// Solution: Use interfaces
type VeriYoneticiInterface interface {
    GorevOlustur(...) (*Gorev, error)
}
```

### Test Isolation

```go
// New DB for each test
func TestXXX(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    // ...
}
```

## Common Development Scenarios

### Scenario 1: Adding a New MCP Tool

**Step-by-step process:**

1. **Define the tool schema** in `internal/mcp/tool_registry.go`:

```go
{
    Name: "gorev_my_new_tool",
    Description: "My new tool description",
    InputSchema: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "param1": map[string]interface{}{
                "type":        "string",
                "description": "Parameter description",
            },
        },
        "required": []interface{}{"param1"},
    },
}
```

2. **Implement the handler** in `internal/mcp/handlers.go`:

```go
func (h *MCPHandlers) MyNewTool(args map[string]interface{}) (*mcp.CallToolResult, error) {
    // 1. Extract and validate parameters
    param1, ok := args["param1"].(string)
    if !ok || param1 == "" {
        return mcp.NewToolResultError("param1 gerekli"), nil
    }

    // 2. Call business logic
    result, err := h.isYonetici.DoSomething(param1)
    if err != nil {
        return mcp.NewToolResultError(fmt.Sprintf("Hata: %v", err)), nil
    }

    // 3. Format response
    output := fmt.Sprintf("✓ İşlem başarılı: %s", result)
    return mcp.NewToolResultText(output), nil
}
```

3. **Add tests** in `internal/mcp/handlers_test.go`:

```go
func TestMyNewTool(t *testing.T) {
    handlers := setupTestHandlers(t)

    args := map[string]interface{}{
        "param1": "test-value",
    }

    result, err := handlers.MyNewTool(args)
    require.NoError(t, err)
    require.Contains(t, result.Content[0].Text, "başarılı")
}
```

4. **Update documentation** in `docs/tr/mcp-araclari.md` and `docs/api/MCP_TOOLS_REFERENCE.md`

5. **Run tests and verify**:

```bash
make test
make build
./gorev serve --test
```

### Scenario 2: Fixing a Bug

**Example: Task status not updating properly**

1. **Reproduce the bug**:

```bash
# Create a failing test first
func TestGorevDurumGuncelle(t *testing.T) {
    // Setup
    handlers := setupTestHandlers(t)
    gorev := createTestGorev(t, handlers)

    // Execute
    args := map[string]interface{}{
        "id": gorev.ID,
        "durum": "tamamlandi",
    }
    result, err := handlers.GorevGuncelle(args)

    // Verify
    require.NoError(t, err)
    updated := getGorev(t, handlers, gorev.ID)
    assert.Equal(t, "tamamlandi", updated.Durum)
}
```

2. **Identify root cause**:

```bash
# Run with debug logging
./gorev serve --debug

# Or add debug prints
log.Printf("DEBUG: Updating task %d to status %s", id, durum)
```

3. **Implement fix**:

```go
func (v *veriYonetici) GorevGuncelle(id int, durum string) error {
    // BEFORE (wrong):
    // _, err := v.db.Exec("UPDATE gorevler SET durum = ?", durum)

    // AFTER (correct):
    _, err := v.db.Exec(
        "UPDATE gorevler SET durum = ?, guncelleme_tarih = CURRENT_TIMESTAMP WHERE id = ?",
        durum, id,
    )
    return err
}
```

4. **Verify fix**:

```bash
make test
# All tests should pass now
```

5. **Document the fix** in CHANGELOG.md

### Scenario 3: Adding Database Schema Changes

**Example: Adding a new field to tasks**

1. **Create migration file**:

```bash
# File: migrations/000011_add_task_priority_level.up.sql
ALTER TABLE gorevler ADD COLUMN priority_level INTEGER DEFAULT 0;
CREATE INDEX idx_gorevler_priority_level ON gorevler(priority_level);
```

```bash
# File: migrations/000011_add_task_priority_level.down.sql
DROP INDEX IF EXISTS idx_gorevler_priority_level;
ALTER TABLE gorevler DROP COLUMN priority_level;
```

2. **Update data model** in `internal/gorev/modeller.go`:

```go
type Gorev struct {
    // ... existing fields
    PriorityLevel int `json:"priority_level,omitempty"`
}
```

3. **Update data access layer** in `internal/gorev/veri_yonetici.go`:

```go
func (v *veriYonetici) GorevOlustur(gorev *Gorev) (*Gorev, error) {
    result, err := v.db.Exec(`
        INSERT INTO gorevler (baslik, aciklama, priority_level)
        VALUES (?, ?, ?)
    `, gorev.Baslik, gorev.Aciklama, gorev.PriorityLevel)
    // ...
}
```

4. **Test migration**:

```bash
# Test up migration
make build
./gorev serve --test

# Test down migration
./gorev migrate down
./gorev migrate up
```

5. **Update documentation and examples**

## Debugging Tips

### Using VS Code Debugger

**launch.json** configuration:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Gorev Server",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/gorev",
            "args": ["serve", "--debug"],
            "env": {
                "GOREV_LANG": "tr"
            }
        },
        {
            "name": "Debug Tests",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${fileDirname}",
            "args": ["-v", "-run", "${selectedText}"]
        }
    ]
}
```

### Common Debugging Commands

```bash
# Run with debug logging
./gorev serve --debug

# Run specific test with verbose output
go test -v -run TestGorevOlustur ./internal/gorev

# Check database state
sqlite3 ~/.gorev/gorev.db "SELECT * FROM gorevler LIMIT 5;"

# Profile CPU usage
go test -cpuprofile cpu.prof -bench .
go tool pprof cpu.prof

# Check memory usage
go test -memprofile mem.prof -bench .
go tool pprof mem.prof

# Race condition detection
go test -race ./...
```

### Debugging MCP Communication

```bash
# Test MCP tool directly
echo '{"method":"tools/call","params":{"name":"gorev_listele","arguments":{}}}' | \
  ./gorev serve --stdio

# Monitor MCP traffic
tail -f ~/.gorev/mcp-debug.log

# Validate JSON schemas
cat tool-schema.json | jq '.'
```

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with delve
dlv debug cmd/gorev/main.go -- serve --debug

# Set breakpoints
(dlv) break internal/mcp/handlers.go:123
(dlv) continue

# Inspect variables
(dlv) print gorev
(dlv) locals
(dlv) args
```

## Development Workflows

### Daily Development Workflow

```bash
# 1. Pull latest changes
git pull origin main

# 2. Create feature branch
git checkout -b feature/my-feature

# 3. Make changes and test frequently
make test
./gorev serve --test

# 4. Format and lint
make fmt
golangci-lint run

# 5. Commit incrementally
git add .
git commit -m "feat: add my feature"

# 6. Push and create PR
git push origin feature/my-feature
```

### Testing Workflow

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/gorev -v

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests
go test ./test -tags=integration

# Benchmark tests
go test -bench=. ./internal/gorev
```

### Release Workflow

```bash
# 1. Update version in Makefile
VERSION ?= 0.16.1

# 2. Update CHANGELOG.md
# Add new version entry with changes

# 3. Build all platforms
make build-all

# 4. Run full test suite
make test
make test-integration

# 5. Tag release
git tag -a v0.16.1 -m "Release v0.16.1"
git push origin v0.16.1

# 6. Create GitHub release
gh release create v0.16.1 --generate-notes
```

## Performance Optimization Tips

### Database Queries

```go
// ❌ Bad: N+1 queries
for _, gorev := range gorevler {
    etiketler := v.GetEtiketler(gorev.ID)  // Separate query for each task
}

// ✅ Good: Single query with JOIN
gorevler := v.GetGorevlerWithEtiketler()  // One query with JOIN
```

### Memory Management

```go
// ❌ Bad: Loading all data
gorevler := v.GetAllGorevler()  // May load thousands of tasks

// ✅ Good: Pagination
gorevler := v.GetGorevlerPaginated(limit, offset)
```

### Caching Strategy

```go
// For frequently accessed, rarely changed data
type Cache struct {
    templates map[string]*GorevTemplate
    mu        sync.RWMutex
}

func (c *Cache) GetTemplate(id string) *GorevTemplate {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.templates[id]
}
```

## Useful Resources

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [MCP Specification](https://modelcontextprotocol.io/docs)
- [SQLite Best Practices](https://www.sqlite.org/bestpractice.html)
- [Go Testing Best Practices](https://github.com/golang/go/wiki/TestComments)

## Related Documentation

- [System Architecture](../architecture/architecture-v2.md)
- [API Reference](../api/reference.md)
- [MCP Tools](../guides/user/mcp-tools.md)
- [VS Code Extension](../guides/user/vscode-extension.md)

---

*💻 This developer guide was created in collaboration with Claude (Anthropic) - AI & Human: The perfect documentation team!*
