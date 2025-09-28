# Geliştirici Rehberi

**Version:** v0.15.24
**Last Updated:** 28 September 2025
**Status:** Production Ready

---

## 📋 İçindekiler

- [Hızlı Başlangıç](#hızlı-başlangıç)
- [Geliştirme Ortamı Kurulumu](#geliştirme-ortamı-kurulumu)
- [Proje Mimarisi](#proje-mimarisi)
- [Kod Standartları](#kod-standartları)
- [Test Stratejisi](#test-stratejisi)
- [Contribution Süreci](#contribution-süreci)
- [Debugging & Troubleshooting](#debugging--troubleshooting)
- [Release Süreci](#release-süreci)

## 🚀 Hızlı Başlangıç

### Ön Gereksinimler

- **Go:** 1.23.2+ (required)
- **Node.js:** 18+ (NPM package ve VS Code extension için)
- **Git:** Latest version
- **VS Code:** Önerilen IDE (isteğe bağlı)

### İlk Kurulum (5 Dakika)

```bash
# 1. Repository clone
git clone https://github.com/msenol/Gorev.git
cd Gorev

# 2. Go dependencies
cd gorev-mcpserver
go mod download
go mod tidy

# 3. İlk build
make build

# 4. Test çalıştır
make test

# 5. Development server başlat
./gorev init
./gorev serve --debug
```

### Hızlı Test

```bash
# Terminal 2'de test
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./gorev serve
```

## 🛠️ Geliştirme Ortamı Kurulumu

### Go Environment Setup

```bash
# Go version kontrolü
go version  # Should be 1.23.2+

# Environment variables
export GOREV_ROOT=$(pwd)
export GOREV_LANG=en  # Development için İngilizce öneriliyor
export GOREV_DEBUG=true

# bashrc/zshrc'ye ekle
echo 'export GOREV_ROOT=/path/to/gorev' >> ~/.bashrc
echo 'export GOREV_DEBUG=true' >> ~/.bashrc
```

### Development Tools

```bash
# Essential tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest  # Hot reload

# Testing tools
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install gotest.tools/gotestsum@latest
```

### VS Code Setup

Önerilen extensions:

```json
{
  "recommendations": [
    "golang.go",
    "ms-vscode.vscode-json",
    "redhat.vscode-yaml",
    "ms-vscode.test-adapter-converter",
    "davidanson.vscode-markdownlint"
  ]
}
```

VS Code settings (`.vscode/settings.json`):

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v", "-race"],
  "go.coverOnSave": true,
  "go.coverageDecorator": {
    "type": "gutter"
  }
}
```

## 🏗️ Proje Mimarisi

### Klasör Yapısı

```
gorev/
├── gorev-mcpserver/           # Ana MCP server (Go)
│   ├── cmd/gorev/             # CLI entry point
│   ├── internal/              # Private packages
│   │   ├── mcp/               # MCP protocol handlers
│   │   ├── gorev/             # Business logic
│   │   ├── veri/              # Data layer
│   │   └── i18n/              # Internationalization
│   ├── test/                  # Integration tests
│   └── docs/                  # Server-specific docs
├── gorev-npm/                 # NPM package wrapper
├── gorev-vscode/              # VS Code extension
├── docs/                      # Project documentation
└── scripts/                   # Build & deployment scripts
```

### Kod Organizasyonu

#### Internal Packages

- **`cmd/gorev/`**: CLI commands (cobra)
- **`internal/mcp/`**: MCP protocol implementation
- **`internal/gorev/`**: Core business logic
- **`internal/veri/`**: Database abstraction layer
- **`internal/i18n/`**: Multi-language support

#### Key Design Patterns

1. **Clean Architecture**: Domain logic isolated from external concerns
2. **Dependency Injection**: Interfaces for testability
3. **Repository Pattern**: Data access abstraction
4. **Command Pattern**: CLI commands implementation

### Database Schema

```sql
-- Core tables
CREATE TABLE projeler (id, isim, tanim, olusturma_tarihi);
CREATE TABLE gorevler (id, baslik, aciklama, durum, oncelik, parent_id, proje_id);
CREATE TABLE baglantilar (id, kaynak_id, hedef_id, baglanti_tipi);
CREATE TABLE etiketler (id, isim, renk);
CREATE TABLE gorev_etiketleri (gorev_id, etiket_id);

-- Advanced features
CREATE TABLE gorev_templateleri (id, isim, kategori, sema);
CREATE TABLE ai_interactions (id, session_id, query, response);
CREATE TABLE filter_profiles (id, isim, filters, kullanici_id);
```

## 📏 Kod Standartları

### Go Code Style

**Temel Kurallar:**

- Go standard formatting (`gofmt`)
- Effective Go guidelines
- Package-level documentation
- Error wrapping with context

**Naming Conventions:**

```go
// Domain terms: Turkish
type Gorev struct {
    ID       string
    Baslik   string
    Durum    string  // beklemede, devam_ediyor, tamamlandi
    Oncelik  string  // dusuk, orta, yuksek
}

// Technical terms: English
type DatabaseManager interface {
    CreateConnection() error
    ExecuteQuery(query string) (Result, error)
}
```

**Error Handling:**

```go
// Always wrap errors with context
func (g *GorevManager) GetGorev(id string) (*Gorev, error) {
    gorev, err := g.db.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("gorev bulunamadı (id: %s): %w", id, err)
    }
    return gorev, nil
}

// Use specific error types
var (
    ErrGorevNotFound = errors.New("görev bulunamadı")
    ErrInvalidStatus = errors.New("geçersiz durum")
)
```

**Testing Conventions:**

```go
func TestGorevManager_CreateGorev(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateGorevInput
        want    *Gorev
        wantErr bool
    }{
        {
            name: "valid_gorev_creation",
            input: CreateGorevInput{
                Baslik: "Test görevi",
                Durum:  "beklemede",
            },
            want: &Gorev{
                Baslik: "Test görevi",
                Durum:  "beklemede",
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### I18n Guidelines

```go
// Good: Use i18n for user-facing messages
return i18n.T("gorev.created", map[string]interface{}{
    "title": gorev.Baslik,
    "id":    gorev.ID,
})

// Bad: Hardcoded Turkish
return fmt.Sprintf("Görev oluşturuldu: %s", gorev.Baslik)
```

## 🧪 Test Stratejisi

### Test Kategorileri

#### Unit Tests

```bash
# Run unit tests
make test

# With coverage
make test-coverage

# Specific package
go test ./internal/gorev -v
```

#### Integration Tests

```bash
# MCP server integration
go test ./test -v

# Database integration
go test ./internal/veri -v -tags integration
```

#### Performance Tests

```bash
# Benchmark tests
go test ./internal/gorev -bench=. -benchmem
```

### Test Structure

```go
func TestGorevManager_Integration(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)

    manager := gorev.NewGorevManager(db)

    // Test scenarios
    t.Run("create_and_retrieve", func(t *testing.T) {
        // Implementation
    })

    t.Run("update_status", func(t *testing.T) {
        // Implementation
    })
}
```

### Test Data Management

```go
// Test helpers
func CreateTestGorev(t *testing.T, db *sql.DB) *Gorev {
    gorev := &Gorev{
        ID:     generateTestID(),
        Baslik: "Test Görev",
        Durum:  "beklemede",
    }

    err := db.CreateGorev(gorev)
    require.NoError(t, err)

    return gorev
}

// Cleanup
func CleanupTestData(t *testing.T, db *sql.DB) {
    _, err := db.Exec("DELETE FROM gorevler WHERE baslik LIKE 'Test%'")
    require.NoError(t, err)
}
```

### Coverage Targets

- **Unit Tests:** >90%
- **Integration Tests:** >80%
- **Overall Coverage:** >85%

## 🤝 Contribution Süreci

### Git Workflow

#### Branch Strategy

```bash
# Feature development
git checkout -b feature/new-mcp-tool
git checkout -b fix/database-connection-issue
git checkout -b docs/api-documentation

# Naming convention
feature/<description>
fix/<bug-description>
docs/<documentation-area>
refactor/<area>
test/<test-area>
```

#### Commit Messages

```bash
# Format: type(scope): description
git commit -m "feat(mcp): add advanced search tool"
git commit -m "fix(db): resolve connection pool leak"
git commit -m "docs(api): add comprehensive API reference"
git commit -m "test(gorev): increase coverage to 90%"

# Types: feat, fix, docs, test, refactor, style, chore
```

### Code Review Process

#### Pull Request Template

```markdown
## 📋 Description
Brief description of changes

## 🔄 Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Test improvement

## 🧪 Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## 📚 Documentation
- [ ] Code comments updated
- [ ] API documentation updated
- [ ] User documentation updated

## ✅ Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] No breaking changes
```

#### Review Criteria

1. **Code Quality**
   - Follows Go best practices
   - Proper error handling
   - Adequate test coverage
   - Clear documentation

2. **Functionality**
   - Solves the intended problem
   - No regressions introduced
   - Edge cases handled

3. **Performance**
   - No unnecessary allocations
   - Database queries optimized
   - Memory leaks prevented

### Pre-commit Checks

```bash
# Install pre-commit hooks
make install-hooks

# Manual check
make lint
make test
make fmt
go vet ./...
```

## 🐛 Debugging & Troubleshooting

### Debug Mode

```bash
# Enable debug logging
export GOREV_DEBUG=true
./gorev serve --debug --lang=en

# Log levels
export GOREV_LOG_LEVEL=debug  # debug, info, warn, error
```

### Common Issues

#### Database Issues

```bash
# Reset database
rm .gorev/gorev.db
./gorev init

# Check database integrity
sqlite3 .gorev/gorev.db "PRAGMA integrity_check;"

# Enable SQL logging
export GOREV_SQL_DEBUG=true
```

#### MCP Connection Issues

```bash
# Test MCP connection
echo '{"jsonrpc": "2.0", "id": 1, "method": "tools/list"}' | ./gorev serve

# Enable MCP debug
export GOREV_MCP_DEBUG=true
```

#### Performance Issues

```bash
# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Profile memory usage
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

### Debugging Tools

```bash
# Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/gorev

# Race condition detection
go test -race ./...

# Memory leak detection
go test -race -memprofile=mem.prof ./...
```

## 🚀 Release Süreci

### Version Management

```bash
# Update version
make version VERSION=0.15.25

# This updates:
# - gorev-mcpserver/Makefile
# - gorev-npm/package.json
# - gorev-vscode/package.json
# - server.json
```

### Build Process

```bash
# Local build
make build-all

# Cross-platform build
make build-cross

# Docker build
make docker-build
```

### Testing Before Release

```bash
# Full test suite
make test-all

# Integration test
make test-integration

# Performance test
make test-performance

# Security scan
make security-scan
```

### Release Automation

```bash
# Use automated release script
./.claude/commands/release.md

# This handles:
# - Version updates
# - Binary builds
# - NPM publishing
# - GitHub release
# - MCP Registry publishing
# - Documentation updates
```

### Manual Release Steps

```bash
# 1. Update version
make version VERSION=0.15.25

# 2. Run tests
make test-all

# 3. Build binaries
make build-cross

# 4. Create GitHub release
gh release create v0.15.25 --title "v0.15.25 - Feature Release"

# 5. Publish NPM package
cd gorev-npm
npm publish

# 6. Update documentation
make docs-update
```

## 📚 Yararlı Kaynaklar

### Internal Documentation

- **[Architecture Guide](../architecture/technical-specification-v2.md)** - System design
- **[API Reference](../api/api-referans.md)** - Complete API documentation
- **[Testing Guide](testing-guide.md)** - Detailed testing strategies
- **[Debugging Guide](debugging-guide.md)** - Troubleshooting help

### External Resources

- **[Go Documentation](https://golang.org/doc/)** - Official Go docs
- **[MCP Specification](https://modelcontextprotocol.io/)** - MCP protocol details
- **[SQLite Documentation](https://sqlite.org/docs.html)** - Database reference
- **[Cobra Documentation](https://cobra.dev/)** - CLI framework

### Community

- **[GitHub Issues](https://github.com/msenol/Gorev/issues)** - Bug reports and feature requests
- **[GitHub Discussions](https://github.com/msenol/Gorev/discussions)** - Community discussions
- **[VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** - Extension reviews

## 🎯 Contributor Guidelines

### Yeni Özellik Ekleme

1. **Issue oluştur** veya mevcut issue'yu claim et
2. **Design document** hazırla (büyük özellikler için)
3. **Feature branch** oluştur
4. **Implementation** yap
5. **Tests** ekle
6. **Documentation** güncelle
7. **Pull request** oluştur

### Bug Fix Süreci

1. **Reproduce** et
2. **Test case** yaz (failing)
3. **Fix** implement et
4. **Test case** geçtiğini doğrula
5. **Regression test** ekle

### Documentation Contribution

1. **Typo/error** düzeltmeleri: Direkt PR
2. **Content update**: Issue oluştur -> discuss -> PR
3. **New documentation**: Design doc -> approval -> implementation

---

> 🎉 **Welcome to Gorev development!** Sorularınız için GitHub Discussions'ı kullanabilir veya issue açabilirsiniz.
