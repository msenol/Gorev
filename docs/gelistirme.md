# Gorev Geliştirici Rehberi

Bu dokümanda Gorev projesine katkıda bulunmak isteyenler için geliştirme ortamı kurulumu, kod standartları ve katkı süreçleri açıklanmaktadır.

## İçindekiler

- [Geliştirme Ortamı Kurulumu](#geliştirme-ortamı-kurulumu)
- [Proje Yapısı](#proje-yapısı)
- [Kod Standartları](#kod-standartları)
- [Test Yazma](#test-yazma)
- [Yeni Özellik Ekleme](#yeni-özellik-ekleme)
- [MCP Tool Ekleme](#mcp-tool-ekleme)
- [Debugging](#debugging)
- [Katkıda Bulunma](#katkıda-bulunma)

## Geliştirme Ortamı Kurulumu

### Gereksinimler

- Go 1.22 veya üzeri
- Git
- Make (opsiyonel, Makefile kullanımı için)
- golangci-lint (kod kalitesi için)
- Docker (opsiyonel, konteyner testleri için)

### Kurulum Adımları

```bash
# Projeyi klonla
git clone https://github.com/msenol/gorev.git
cd gorev/gorev-mcpserver

# Bağımlılıkları indir
make deps
# veya
go mod download

# Projeyi derle
make build
# veya
go build -o gorev cmd/gorev/main.go

# Testleri çalıştır
make test
# veya
go test ./...
```

### IDE Ayarları

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

## Proje Yapısı

```
gorev/
├── gorev-mcpserver/             # MCP server projesi
│   ├── cmd/
│   │   └── gorev/
│   │       └── main.go          # Ana uygulama giriş noktası
│   ├── internal/
│   │   ├── gorev/               # Domain logic
│   │   │   ├── modeller.go      # Veri modelleri
│   │   │   ├── is_yonetici.go   # Business logic
│   │   │   ├── veri_yonetici.go # Data access layer
│   │   │   ├── template_yonetici.go # Template yönetimi
│   │   │   └── *_test.go        # Unit testler
│   │   └── mcp/                 # MCP protokol katmanı
│   │       ├── server.go        # MCP server
│   │       └── handlers.go      # Tool handler'ları
│   ├── migrations/              # Veritabanı migration'ları
│   └── test/                    # Integration testler
├── gorev-vscode/                # VS Code extension
├── docs/                        # Dokümantasyon
└── scripts/                     # Yardımcı scriptler
```

### Paket Açıklamaları

- **cmd/gorev**: CLI komutları ve server başlatma
- **internal/gorev**: Core business logic ve domain modelleri
- **internal/mcp**: MCP protokol implementasyonu
- **migrations**: SQL migration dosyaları (golang-migrate formatı)

## Kod Standartları

### Genel Kurallar

1. **Go idiomlarını takip et**: Effective Go ve Go Code Review Comments'i oku
2. **Türkçe domain terimleri**: Görev, proje, durum gibi domain terimlerini Türkçe kullan
3. **İngilizce teknik terimler**: Kod yorumları ve teknik terimler İngilizce
4. **Error handling**: Explicit error döndür, panic kullanma

### Naming Conventions

```go
// Domain modelleri - Türkçe
type Gorev struct { ... }
type Proje struct { ... }

// Interface'ler - Türkçe + -ci/-ici eki
type VeriYonetici interface { ... }
type IsYonetici interface { ... }

// Method isimleri - Türkçe fiil + İngilizce nesne (gerekirse)
func (v *veriYonetici) GorevOlustur(...) { ... }
func (v *veriYonetici) ProjeListele(...) { ... }

// Sabitler - UPPER_SNAKE_CASE
const VERITABANI_VERSIYON = "1.2.0"

// Private değişkenler - camelCase
var aktifProjeID int
```

### Code Style

```go
// İyi: Kısa ve açık fonksiyonlar
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

// Kötü: Uzun ve karmaşık fonksiyonlar
func (v *veriYonetici) HepsiniYap(id int) error {
    // 100+ satır kod...
}
```

### Error Messages

```go
// Türkçe kullanıcı mesajları
return fmt.Errorf("görev bulunamadı: %d", id)
return fmt.Errorf("geçersiz durum değeri: %s", durum)

// Context ile wrap etme
if err != nil {
    return fmt.Errorf("veritabanı bağlantısı kurulamadı: %w", err)
}
```

## Test Yazma

### Unit Test Yapısı

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
            name:    "başarılı oluşturma",
            baslik:  "Test görevi",
            oncelik: "orta",
            wantErr: false,
        },
        {
            name:    "boş başlık",
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
                t.Error("hata beklendi ama nil döndü")
            }
            if !tc.wantErr && err != nil {
                t.Errorf("beklenmeyen hata: %v", err)
            }
        })
    }
}
```

### Test Utilities

```go
// test/test_helpers.go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("test db açılamadı: %v", err)
    }
    
    // Migration'ları çalıştır
    if err := runMigrations(db); err != nil {
        t.Fatalf("migration başarısız: %v", err)
    }
    
    return db
}
```

### Integration Test

```go
// test/integration_test.go
func TestMCPToolIntegration(t *testing.T) {
    // MCP server başlat
    server := setupTestServer(t)
    defer server.Close()
    
    // Tool çağrısı yap
    response, err := server.CallTool("gorev_olustur", map[string]interface{}{
        "baslik":  "Integration test",
        "oncelik": "yuksek",
    })
    
    // Sonucu kontrol et
    assert.NoError(t, err)
    assert.Contains(t, response.Text, "başarıyla oluşturuldu")
}
```

## Yeni Özellik Ekleme

### 1. Domain Model Güncelleme

```go
// internal/gorev/modeller.go
type Gorev struct {
    // Mevcut alanlar...
    YeniAlan string `json:"yeni_alan,omitempty"`
}
```

### 2. Migration Ekleme

```sql
-- migrations/004_yeni_alan.up.sql
ALTER TABLE gorevler ADD COLUMN yeni_alan TEXT DEFAULT '';

-- migrations/004_yeni_alan.down.sql
ALTER TABLE gorevler DROP COLUMN yeni_alan;
```

### 3. Data Layer Güncelleme

```go
// internal/gorev/veri_yonetici.go
func (v *veriYonetici) gorevleriTara(rows *sql.Rows) ([]Gorev, error) {
    // Scan'e yeni alan ekle
    err := rows.Scan(
        &gorev.ID,
        // diğer alanlar...
        &gorev.YeniAlan,
    )
}
```

### 4. Business Logic Güncelleme

```go
// internal/gorev/is_yonetici.go
func (i *isYonetici) GorevOlustur(..., yeniAlan string) (*Gorev, error) {
    // Validation ekle
    if err := validateYeniAlan(yeniAlan); err != nil {
        return nil, err
    }
    
    // Veri katmanını çağır
    return i.veriYonetici.GorevOlustur(..., yeniAlan)
}
```

## MCP Tool Ekleme

### 1. Handler Fonksiyonu Yaz

```go
// internal/mcp/handlers.go
func (h *Handler) handleYeniTool(args map[string]interface{}) (*ToolResult, error) {
    // Parametreleri parse et
    param1, ok := args["param1"].(string)
    if !ok {
        return nil, fmt.Errorf("param1 gerekli")
    }
    
    // Business logic çağır
    result, err := h.isYonetici.YeniIslem(param1)
    if err != nil {
        return nil, mcp.NewToolResultError(err.Error())
    }
    
    // Sonucu döndür
    return &ToolResult{
        Content: []Content{{
            Type: "text",
            Text: fmt.Sprintf("✅ İşlem başarılı: %v", result),
        }},
    }, nil
}
```

### 2. Tool'u Kaydet

```go
// internal/mcp/handlers.go - RegisterTools()
tools = append(tools, Tool{
    Name:        "yeni_tool",
    Description: "Yeni işlem yapar",
    InputSchema: InputSchema{
        Type: "object",
        Properties: map[string]Property{
            "param1": {
                Type:        "string",
                Description: "Parametre açıklaması",
            },
        },
        Required: []string{"param1"},
    },
})
```

### 3. Dokümantasyon Ekle

`docs/mcp-araclari.md` dosyasına yeni tool'u ekle.

### 4. Test Yaz

```go
// test/integration_test.go
func TestYeniTool(t *testing.T) {
    // Test senaryoları...
}
```

## Debugging

### Debug Mode

```bash
# Debug log'ları aktif
./gorev serve --debug

# Veya environment variable
DEBUG=true ./gorev serve
```

### Logging

```go
import "log/slog"

// Debug log
slog.Debug("işlem başladı", "id", gorevID, "durum", durum)

// Error log
slog.Error("veritabanı hatası", "error", err)
```

### Profiling

```go
import _ "net/http/pprof"

// main.go'da
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

## VS Code Extension Geliştirme

### Extension Kurulumu

```bash
cd gorev-vscode

# Bağımlılıkları yükle
npm install

# TypeScript derle
npm run compile

# Watch mode (geliştirme için)
npm run watch
```

### Extension Test Etme

1. VS Code'da `gorev-vscode` klasörünü aç
2. F5 tuşuna bas (veya Run > Start Debugging)
3. Yeni VS Code penceresi açılacak (Extension Development Host)
4. Extension'ı test et

### Extension Yapısı

```
gorev-vscode/
├── src/
│   ├── extension.ts          # Ana giriş noktası
│   ├── mcp/
│   │   ├── client.ts        # MCP client implementasyonu
│   │   └── types.ts         # TypeScript tipleri
│   ├── commands/            # Komut handler'ları
│   ├── providers/           # TreeView provider'ları
│   └── models/              # Data modelleri
├── package.json             # Extension manifest
└── tsconfig.json           # TypeScript config
```

### Yeni Komut Ekleme

1. **package.json'a komut ekle**:
```json
{
  "contributes": {
    "commands": [
      {
        "command": "gorev.newCommand",
        "title": "Gorev: New Command"
      }
    ]
  }
}
```

2. **Command handler ekle**:
```typescript
// src/commands/newCommand.ts
export async function newCommand() {
    // Komut implementasyonu
}
```

3. **Extension.ts'de kaydet**:
```typescript
context.subscriptions.push(
    vscode.commands.registerCommand('gorev.newCommand', newCommand)
);
```

### Extension Debugging

1. **Output Channel kullan**:
```typescript
const outputChannel = vscode.window.createOutputChannel('Gorev');
outputChannel.appendLine('Debug mesajı');
```

2. **Breakpoint koy**: VS Code'da TypeScript dosyalarına breakpoint ekle

3. **Debug Console**: Extension Development Host'ta Debug Console'u kontrol et

## Katkıda Bulunma

### Pull Request Süreci

1. **Issue Aç**: Önce bir issue açarak ne yapmak istediğini açıkla
2. **Fork & Branch**: Projeyi fork'la ve feature branch oluştur
   ```bash
   git checkout -b feature/yeni-ozellik
   ```
3. **Kod Yaz**: Kod standartlarına uygun şekilde geliştir
4. **Test Yaz**: %80+ coverage hedefle
5. **Commit**: Anlamlı commit mesajları kullan
   ```bash
   git commit -m "feat: yeni özellik ekle"
   git commit -m "fix: hata düzelt"
   git commit -m "docs: dokümantasyon güncelle"
   ```
6. **Push & PR**: Branch'i push'la ve PR aç

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types:
- `feat`: Yeni özellik
- `fix`: Hata düzeltme
- `docs`: Dokümantasyon
- `style`: Formatting, missing semi-colons, etc.
- `refactor`: Kod düzenleme
- `test`: Test ekleme/düzeltme
- `chore`: Maintenance

### Code Review Checklist

- [ ] Testler yazıldı ve geçiyor
- [ ] Dokümantasyon güncellendi
- [ ] Kod standartlarına uygun
- [ ] Breaking change yok (varsa dokümante edildi)
- [ ] Performance etkileri düşünüldü

## Sık Karşılaşılan Sorunlar

### SQLite Locked Error

```go
// Çözüm: WAL mode kullan
db.Exec("PRAGMA journal_mode=WAL")
```

### Import Cycle

```go
// Çözüm: Interface kullan
type VeriYoneticiInterface interface {
    GorevOlustur(...) (*Gorev, error)
}
```

### Test Isolation

```go
// Her test için yeni DB
func TestXXX(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    // ...
}
```

## Faydalı Kaynaklar

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [MCP Specification](https://modelcontextprotocol.io/docs)
- [SQLite Best Practices](https://www.sqlite.org/bestpractice.html)

## İlgili Dokümantasyon

- [Sistem Mimarisi](mimari.md)
- [API Referansı](api-referans.md)
- [MCP Araçları](mcp-araclari.md)
- [VS Code Extension](vscode-extension.md)

---

<div align="center">

*💻 Bu geliştirici rehberi Claude (Anthropic) ile işbirliği içinde oluşturulmuştur - AI & İnsan: Mükemmel dokümantasyon takımı!*

</div>