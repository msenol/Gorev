# Git Commit İşlemleri - Güvenli ve Organize

## Görev
Commit edilmeyi bekleyen dosyaları analiz et ve uygun commit stratejisi ile işle.

## İşlem Adımları

### 1. Dosya Analizi ve Filtreleme
```bash
git status --porcelain
```
- Tüm değişiklikleri kategorize et (Added, Modified, Deleted, Renamed)
- Aşağıdaki dosyaları ASLA commit etme:
  - `*.db`, `*.db-journal`, `*.sqlite` veritabanı dosyaları
  - `.env`, `.env.local` gibi environment dosyaları
  - `vendor/`, `dist/`, `build/` klasörleri
  - `*.log`, `*.tmp`, `*.cache` uzantılı dosyalar
  - IDE config dosyaları (`.idea/`, `.vscode/settings.json`)
  - Go binary dosyaları (`gorev`, `gorev.exe`)

### 2. Değişiklik Gruplandırması
Değişiklikleri mantıksal gruplara ayır:
- **Feature**: Yeni özellik eklemeleri (yeni MCP tool, veri modeli)
- **Fix**: Bug düzeltmeleri
- **Refactor**: Kod iyileştirmeleri (davranış değişikliği olmadan)
- **Docs**: Dokümantasyon güncellemeleri
- **Style**: Formatting, gofmt düzeltmeleri
- **Test**: Test eklemeleri veya düzeltmeleri
- **Chore**: Build process, Makefile, Docker güncellemeleri
- **Perf**: Performance iyileştirmeleri

### 3. Commit Stratejisi Belirleme
- Tek bir mantıksal değişiklik = Tek commit
- Farklı paketlere ait değişiklikler = Ayrı commitler
- Büyük feature = Ana commit + destekleyici commitler

### 4. Emin Olamadığın Durumlar İçin Sorgulama
Aşağıdaki durumlarda kullanıcıya mutlaka sor:
- go.mod ve go.sum dosyalarındaki değişiklikler
- Migration dosyaları (veritabanı şema değişiklikleri)
- Dockerfile veya docker-compose değişiklikleri
- Binary dosyalar
- Büyük dosyalar (> 1MB)

### 5. Commit Mesajı Formatı
Conventional Commits formatını kullan:
```
<type>(<scope>): <subject>

<body>

<footer>
```

Örnekler:
```
feat(mcp): add task dependency management tools

- Implement gorev_bagla tool for creating dependencies
- Add cycle detection algorithm
- Update veri_yonetici with new baglantilar table

Closes #12
```

## Çıktı Formatı

### 1. Değişiklik Özeti
```
📊 Değişiklik Analizi:
- Toplam dosya: X
- Eklenen: X dosya
- Değiştirilen: X dosya  
- Silinen: X dosya

📁 Etkilenen Modüller:
- internal/mcp/
- internal/gorev/
- docs/
```

### 2. Commit Planı
```
🎯 Önerilen Commit Stratejisi:

Commit 1: feat(mcp): implement MCP Go SDK integration
- internal/mcp/handlers.go
- internal/mcp/server.go
- go.mod

Commit 2: docs: add comprehensive Turkish documentation
- docs/README.md
- docs/kurulum.md
- docs/kullanim.md
- docs/mcp-araclari.md

Commit 3: test: add integration tests for MCP handlers
- test/integration_test.go
```

### 3. Dikkat Edilmesi Gerekenler
```
⚠️ Uyarılar:
- [ ] go.sum dosyası güncellenmeli (go mod tidy)
- [ ] Binary build edilip test edilmeli
- [ ] Docker image yeniden build edilmeli
```

### 4. Kullanıcı Onayı İsteme
```
❓ Emin olamadığım dosyalar:

1. arsiv-kotlin/
   - Eski Kotlin kodları arşivlenmiş
   - Git'e eklemeli miyim? (E/H)

2. gorev.db
   - SQLite veritabanı dosyası
   - Genelde commit edilmez, eklemeli miyim? (E/H)
```

## Güvenlik Kontrol Listesi

Commit öncesi mutlaka kontrol et:
- [ ] Hassas bilgi içermiyor (API keys, passwords, tokens)
- [ ] Gereksiz fmt.Println() kalmamış
- [ ] TODO/FIXME commentleri uygun
- [ ] Test dosyaları sadece test implementasyonu içeriyor
- [ ] Binary dosyalar commit edilmemiş

## Özel Kurallar (Gorev Projesi İçin)

1. **Go Standartları**: 
   - CLAUDE.md'deki "Development Commands" kullan
   - `make fmt` ve `go vet ./...` kontrollerini çalıştır
   - go.mod değişikliklerinde `make deps` çalıştır
   
2. **MCP Tool Değişiklikleri**: 
   - CLAUDE.md'deki "Adding New MCP Tools" adımlarını takip et
   - Yeni tool eklendiğinde docs/mcp-araclari.md güncellenmeli
   - Handler testleri test/integration_test.go'ya eklenmeli
   
3. **Veritabanı Değişiklikleri**:
   - Schema değişikliklerini tablolariOlustur() fonksiyonuna ekle
   - CLAUDE.md'deki "Database Schema" bölümünü güncelle
   
4. **Docker/Build**:
   - Dockerfile değişikliklerinde multi-stage build korunmalı
   - Version bilgisi Makefile'da LDFLAGS ile yönetiliyor

## Komut Örnekleri

```bash
# CLAUDE.md'deki Development Commands kullan:
make fmt              # Format kontrolü ve düzeltme
go vet ./...         # Vet kontrolü
make deps            # Modül temizliği (go mod tidy dahil)
make test            # Testleri çalıştır
make lint            # Lint kontrolü

# Staged dosyaları göster
git diff --staged --name-only

# Belirli dosyaları stage'e al
git add internal/mcp/*.go

# Interactive staging (parçalı commit için)
git add -p

# Commit with message
git commit -m "feat(mcp): implement MCP Go SDK integration"

# Amend last commit (sadece küçük düzeltmeler için)
git commit --amend
```