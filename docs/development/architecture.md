# Sistem Mimarisi

> **Versiyon**: Bu dokümantasyon v0.11.1 için geçerlidir.  
> **Son Güncelleme**: 16 August 2025

Gorev'in teknik mimarisi ve tasarım kararları.

## 🏗️ Genel Mimari

```
┌─────────────────┐     ┌─────────────────┐
│  Claude Desktop │     │   Claude Code   │
│       /Code     │     │                 │
└────────┬────────┘     └────────┬────────┘
         │                       │
         └───────────┬───────────┘
                     │ MCP Protocol (stdio/JSON-RPC)
                     │
              ┌──────▼──────┐
              │ Gorev Server │
              │  (main.go)   │
              └──────┬──────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
   ┌────▼────┐ ┌────▼────┐ ┌────▼────┐
   │   MCP   │ │  İş     │ │  Veri   │
   │ Katmanı │ │ Mantığı │ │ Katmanı │
   └─────────┘ └─────────┘ └────┬────┘
                                 │
                           ┌─────▼─────┐
                           │  SQLite   │
                           │    DB     │
                           └───────────┘
```

## 📁 Proje Yapısı

```
gorev/
├── cmd/
│   └── gorev/
│       └── main.go          # Uygulama giriş noktası
├── internal/               # Özel paketler (dışarı açık değil)
│   ├── mcp/
│   │   └── sunucu.go       # MCP protokol implementasyonu
│   ├── gorev/
│   │   ├── modeller.go     # Domain modelleri
│   │   ├── is_yonetici.go  # İş mantığı katmanı
│   │   └── veri_yonetici.go # Veri erişim katmanı
│   └── veri/               # Gelecek: Alternatif veri katmanları
├── pkg/                    # Dışa açık paketler
│   ├── islem/              # İşlem yönetimi utilities
│   └── sunum/              # Sunum formatları
└── test/                   # Test dosyaları
```

## 🔧 Katman Detayları

### 1. MCP Katmanı (`internal/mcp`)

**Sorumluluklar:**
- JSON-RPC protokolü işleme
- Stdin/stdout üzerinden iletişim
- Tool tanımları ve routing
- Request/response yönetimi

**Temel Bileşenler:**
```go
type Sunucu struct {
    isYonetici *gorev.IsYonetici
    reader     *bufio.Reader
    writer     io.Writer
    mu         sync.Mutex
}
```

### 2. İş Mantığı Katmanı (`internal/gorev`)

**Sorumluluklar:**
- Domain logic implementasyonu
- İş kuralları ve validasyon
- Use case orchestration

**Temel Operasyonlar:**
- Görev CRUD işlemleri
- Proje yönetimi
- Durum geçişleri
- Özet raporlama

### 3. Veri Katmanı (`internal/gorev/veri_yonetici.go`)

**Sorumluluklar:**
- Veritabanı bağlantı yönetimi
- SQL sorguları
- Transaction yönetimi
- Migration işlemleri

**Veritabanı Şeması:**

```sql
-- Projeler tablosu
CREATE TABLE projeler (
    id TEXT PRIMARY KEY,
    isim TEXT NOT NULL,
    tanim TEXT,
    olusturma_tarih DATETIME NOT NULL,
    guncelleme_tarih DATETIME NOT NULL
);

-- Görevler tablosu
CREATE TABLE gorevler (
    id TEXT PRIMARY KEY,
    baslik TEXT NOT NULL,
    aciklama TEXT,
    durum TEXT NOT NULL DEFAULT 'beklemede',
    oncelik TEXT NOT NULL DEFAULT 'orta',
    proje_id TEXT,
    olusturma_tarih DATETIME NOT NULL,
    guncelleme_tarih DATETIME NOT NULL,
    FOREIGN KEY (proje_id) REFERENCES projeler(id)
);

-- Bağlantılar tablosu (gelecek özellik)
CREATE TABLE baglantilar (
    id TEXT PRIMARY KEY,
    kaynak_id TEXT NOT NULL,
    hedef_id TEXT NOT NULL,
    baglanti_tip TEXT NOT NULL,
    FOREIGN KEY (kaynak_id) REFERENCES gorevler(id),
    FOREIGN KEY (hedef_id) REFERENCES gorevler(id)
);
```

## 🔄 İstek Akışı

1. **Claude → Gorev:**
   ```json
   {
     "jsonrpc": "2.0",
     "id": 1,
     "method": "tools/call",
     "params": {
       "name": "gorev_olustur",
       "arguments": {
         "baslik": "Yeni görev"
       }
     }
   }
   ```

2. **MCP Katmanı:**
   - JSON parse edilir
   - Tool adı ve parametreler çıkarılır
   - İlgili handler çağrılır

3. **İş Mantığı:**
   - Validasyon yapılır
   - UUID oluşturulur
   - Model nesnesi yaratılır

4. **Veri Katmanı:**
   - SQL INSERT çalıştırılır
   - Transaction commit edilir

5. **Gorev → Claude:**
   ```json
   {
     "jsonrpc": "2.0",
     "id": 1,
     "result": {
       "content": [{
         "type": "text",
         "text": "✓ Görev oluşturuldu: Yeni görev (ID: ...)"
       }]
     }
   }
   ```

## 🎯 Tasarım Prensipleri

### 1. Katmanlı Mimari
- Her katmanın net sorumlulukları var
- Katmanlar arası bağımlılık tek yönlü
- Test edilebilirlik ön planda

### 2. Domain-Driven Design
- İş mantığı domain modellerinde
- Altyapı detayları izole edilmiş
- Ubiquitous language kullanımı

### 3. SOLID Prensipleri
- **S**ingle Responsibility
- **O**pen/Closed
- **L**iskov Substitution
- **I**nterface Segregation
- **D**ependency Inversion

### 4. Go İdiomları
- Explicit error handling
- Interface kullanımı
- Composition over inheritance
- Concurrency safety

## 🔒 Güvenlik Mimarisi

### 1. Input Validasyonu
- Tüm MCP inputları validate edilir
- SQL injection koruması
- Path traversal koruması

### 2. Veri İzolasyonu
- Her kullanıcı kendi veritabanına sahip
- Cross-user erişim yok
- Dosya sistemi izolasyonu

### 3. Error Handling
- Hassas bilgiler loglanmaz
- Stack trace'ler gizlenir
- Güvenli varsayılanlar

## 🚀 Performans Optimizasyonları

### 1. Veritabanı
- Index'ler eklendi (durum, proje_id)
- Prepared statements kullanımı
- Connection pooling (gelecek)

### 2. Bellek Yönetimi
- Minimal allocation
- Buffer reuse
- Lazy loading

### 3. Concurrency
- Goroutine kullanımı (gelecek)
- Channel-based communication
- Lock-free algoritmalar

## 📈 Ölçeklenebilirlik

### Mevcut Limitler
- Tek SQLite dosyası
- Senkron işlem modeli
- Lokal dosya sistemi

### Gelecek İyileştirmeler
1. **Veri Katmanı:**
   - PostgreSQL desteği
   - Redis cache katmanı
   - Distributed storage

2. **İşlem Modeli:**
   - Async task processing
   - Event-driven architecture
   - Message queue entegrasyonu

3. **API Genişletme:**
   - REST API endpoint'leri
   - GraphQL desteği
   - WebSocket real-time updates

## 🔧 Konfigürasyon

### Environment Variables
```bash
GOREV_DATA_DIR=/path/to/data
GOREV_LOG_LEVEL=debug|info|warn|error
GOREV_MAX_CONNECTIONS=10
GOREV_TIMEOUT=30s
```

### Yapılandırma Dosyası (Planlanan)
```yaml
server:
  transport: stdio
  timeout: 30s

database:
  type: sqlite
  path: ${GOREV_DATA_DIR}/gorev.db
  
logging:
  level: info
  format: json
  output: stderr
```

## 📊 Monitoring ve Metrics

### Planlanan Metrikler
- Request/response süreleri
- Tool kullanım istatistikleri
- Hata oranları
- Veritabanı performansı

### Health Check Endpoint
```go
GET /health
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "2h30m",
  "database": "connected"
}
```

## 🔄 Versiyon Yönetimi

### Semantic Versioning
- Major: Breaking changes
- Minor: Yeni özellikler
- Patch: Bug fix'ler

### Backward Compatibility
- MCP protokol versiyonu korunur
- Veritabanı migration'ları
- Deprecation politikası

## 🔨 Major Refactoring Architecture (v0.11.1)

> **Update**: Major architectural improvements implemented in v0.11.1 (16 August 2025)

### Code Refactoring Overview

A significant refactoring was completed to eliminate code smells and improve maintainability:

#### File Structure Changes
```
internal/mcp/
├── handlers.go           # 2,362 lines (reduced from 3,060)
├── tool_registry.go      # 570 lines (NEW)
├── tool_helpers.go       # 286 lines (NEW)
├── test_helpers.go       # DRY test infrastructure (NEW)
├── table_driven_test.go  # Table-driven test patterns (NEW)
├── concurrency_test.go   # DRY concurrency testing (NEW)
├── benchmark_test.go     # Standardized benchmark suite (NEW)
├── dry_validation_test.go # Focused validation tests (NEW)
└── [existing test files...]

internal/i18n/
├── manager.go            # Internationalization system
└── helpers.go            # DRY i18n patterns (NEW)
```

#### Architecture Improvements

##### 1. Tool Registration Pattern
**Before (Code Smell):**
- Single massive `RegisterTools` method: 703 lines
- Mixed concerns: registration + schema definition + documentation
- Hard to maintain and test

**After (Clean Architecture):**
```go
// Simplified delegation pattern
func (h *Handlers) RegisterTools(s *server.MCPServer) {
    registry := NewToolRegistry(h)
    registry.RegisterAllTools(s)
}

// Organized by categories in tool_registry.go
func (tr *ToolRegistry) RegisterAllTools(s *server.MCPServer) {
    tr.registerTaskManagementTools(s)
    tr.registerProjectManagementTools(s)
    tr.registerTemplateTools(s)
    tr.registerAIContextTools(s)
    tr.registerFileWatcherTools(s)
    tr.registerAdvancedTools(s)
}
```

##### 2. Helper Class Architecture
**Common Utility Classes:**
- `ParameterValidator` - Centralized input validation
- `TaskFormatter` - Consistent formatting with emojis
- `ErrorFormatter` - Standardized error messages
- `ResponseBuilder` - Reusable response patterns
- `CommonValidators` - Frequently used validations

**Benefits:**
- DRY principle enforcement
- Consistent behavior across tools
- Easier testing and maintenance
- Reusable components

##### 3. Tool Categories
**Organized by Domain:**
1. **Task Management** (7 tools): CRUD operations
2. **Project Management** (6 tools): Project lifecycle  
3. **Templates** (2 tools): Template-based creation
4. **AI Context** (6 tools): AI session management
5. **File Watcher** (4 tools): File system monitoring
6. **Advanced Tools** (5 tools): Dependencies & hierarchy

#### DRY Patterns Implementation

**i18n DRY Patterns (internal/i18n/helpers.go):**
```go
// Helper functions for consistent translation patterns
func TParam(key string, data interface{}) string
func FormatParameterRequired(param string) string
func FormatInvalidValue(param, value, expected string) string
```

**Testing DRY Infrastructure:**
- **TestCase struct**: Standardized table-driven test patterns
- **BenchmarkConfig struct**: Reusable benchmark configuration
- **ConcurrencyTestConfig struct**: Thread-safety validation patterns
- **Helper functions**: CreateTestServer(), RunTableDrivenTest(), RunBenchmarkSuite()

**Code Quality Metrics:**
- **12 total test files**: Comprehensive coverage with reusable patterns
- **11,124+ lines**: Well-organized Go code across all files
- **60% reduction**: Duplicate strings and validation patterns eliminated

#### Refactoring Impact

**Quantitative Improvements:**
- **File Size**: 698 lines removed (23% reduction)
- **Method Size**: 703-line method → 4-line delegation
- **Code Reuse**: Extracted 286 lines of reusable utilities
- **DRY Implementation**: ~60% reduction in duplicate patterns
- **Test Infrastructure**: 5 new test files with standardized patterns
- **Organization**: 25 tools organized into 6 logical categories

**Qualitative Benefits:**
- ✅ Eliminated massive method code smell
- ✅ Comprehensive DRY patterns implementation
- ✅ Reusable test infrastructure for maintainability
- ✅ Consistent i18n and validation patterns
- ✅ Enhanced testability through smaller units
- ✅ Better separation of concerns
- ✅ Zero breaking changes (API compatibility maintained)

**Rule 15 Compliance:**
- **NO Code Duplication**: DRY principle strictly enforced
- **NO Workarounds**: Proper abstraction and reusability
- **Comprehensive Testing**: Production-ready test patterns
- **Clean Architecture**: Maintainable and well-organized code

## 📚 İlgili Dokümantasyon

- [API Referansı](api-referans.md) - Detaylı API dokümantasyonu
- [Geliştirici Rehberi](gelistirme.md) - Katkıda bulunma
- [MCP Protokolü](https://modelcontextprotocol.io) - MCP spesifikasyonu
- [Thread-Safety Guide](../security/thread-safety.md) - Concurrency patterns
- [Testing Guide](testing-guide.md) - Comprehensive testing strategies