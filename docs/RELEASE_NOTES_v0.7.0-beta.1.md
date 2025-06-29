# Gorev v0.7.0-beta.1 Release Notes

## 📊 Değişiklik Özeti

**Kategori**: Feature/Enhancement  
**Etkilenen Paketler**:
- `gorev-mcpserver`: Path resolution improvements for database and migrations
- `gorev-vscode`: Major UI enhancements with 20+ new features
- `docs/`: Updated documentation for new features

**Breaking Changes**: Yok

## 🚀 Yeni Özellikler

### VS Code Extension - Gelişmiş UI

#### 1. Enhanced TreeView
- Gruplama desteği (durum/öncelik/etiket/proje/tarih)
- Çoklu seçim (Ctrl/Cmd+Click)
- Öncelik bazlı renk kodlaması
- Hızlı tamamlama checkbox'ları
- Görev sayısı ve son tarih badge'leri

#### 2. Drag & Drop Sistemi
- Projeler arası görev taşıma
- Sürükle-bırak ile durum değiştirme
- Öncelik sıralaması
- Bağımlılık oluşturma
- Görsel geri bildirim ve animasyonlar

#### 3. Inline Düzenleme
- F2 veya double-click ile hızlı düzenleme
- Context menu entegrasyonu
- Inline tarih seçici
- Escape/Enter kısayolları

#### 4. Gelişmiş Filtreleme
- Gerçek zamanlı arama
- Çoklu filtre kriterleri
- Kayıtlı filtre profilleri
- Status bar entegrasyonu
- Hızlı filtre butonları

#### 5. Zengin Görev Detay Paneli
- Split-view markdown editörü
- Canlı önizleme
- Bağımlılık görselleştirme
- Aktivite zaman çizelgesi
- Template alan göstergeleri

#### 6. Template Wizard
- Çok adımlı arayüz
- Dinamik form oluşturma
- Alan doğrulama
- Önizleme desteği
- Kategori bazlı organizasyon

#### 7. Test Suite
- Unit testler (markdownParser, mcpClient, treeProviders)
- Integration testler (extension features)
- E2E testler (full workflows)
- Coverage raporlama (c8)

### MCP Server İyileştirmeleri
- `getDatabasePath()`: Executable-relative database path
- `getMigrationsPath()`: Automatic migration discovery
- Farklı dizinlerden çalıştırma desteği

## 🐛 Düzeltmeler

1. **Template Display**: Markdown parser güncellendi, template listesi doğru parse ediliyor
2. **TreeView Classes**: Export edilen class'lar VS Code tarafından instantiate edilebiliyor
3. **TypeScript Errors**: Filter interface property isimleri düzeltildi (Türkçe karşılıkları)
4. **Path Issues**: gorev komutu farklı dizinlerden çalıştırılabilir

## 📝 Güncellenen Dosyalar

### ✅ CHANGELOG.md
- Version: 0.7.0-dev
- Tüm yeni özellikler ve düzeltmeler eklendi

### ✅ CLAUDE.md
- Son güncelleme tarihi: 28 June 2025
- v0.7.0-dev değişiklikleri eklendi
- Important Files bölümü güncellendi

### ✅ gorev-vscode/README.md
- Tüm yeni özellikler detaylandırıldı
- Konfigürasyon seçenekleri güncellendi
- 21 komut dokumentasyonu eklendi
- Test bölümü eklendi

### ✅ gorev-mcpserver/Makefile
- Version: 0.7.0-dev

### ✅ gorev-vscode/package.json
- Version: 0.2.0
- Test dependencies eklendi (mocha, sinon, c8)
- Yeni konfigürasyon seçenekleri

### ✅ gorev-vscode/TASKS.md
- Tüm görevler tamamlandı olarak işaretlendi

## ⚠️ Dikkat Edilecekler

- [x] go.mod version sync kontrolü (0.7.0-dev)
- [x] VS Code extension version (0.2.0)
- [ ] Docker image version tag'i güncellenmeli
- [ ] GitHub release hazırlığı yapılmalı
- [ ] VS Code Marketplace için paketleme

## 🎯 Sonraki Adımlar

1. Test coverage'ı artır (hedef: >80%)
2. VS Code Marketplace için dokümantasyon hazırla
3. Demo GIF'leri oluştur
4. CI/CD pipeline kurulumu
5. Auto-update mekanizması

## 📊 Proje İstatistikleri

- **Yeni TypeScript Dosyaları**: 20+
- **Yeni Komutlar**: 10+
- **Yeni Konfigürasyon Seçenekleri**: 15+
- **Test Dosyaları**: 8
- **Toplam Değişiklik**: 1167 ekleme, 506 silme