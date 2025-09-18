# Kullanım Kılavuzu

> **Versiyon**: Bu dokümantasyon v0.15.5 için geçerlidir.
> **Son Güncelleme**: 18 Eylül 2025

Gorev ile görev yönetiminin temelleri ve gelişmiş özellikler.

## 🎯 Temel Kavramlar

### Görev (Task)
- Yapılacak işlerin temel birimi
- Her görevin benzersiz bir ID'si vardır
- **Durum**: `beklemede`, `devam_ediyor`, `tamamlandi`
- **Öncelik**: `dusuk`, `orta`, `yuksek`
- **Alt görevler**: Sınırsız derinlikte hiyerarşik yapı
- **Bağımlılıklar**: Görevler arası ilişki kurma

### Proje (Project)
- Görevleri gruplamak için kullanılır
- Birden fazla görev içerebilir
- Aktif proje sistemi ile hızlı işlemler
- İsteğe bağlıdır

### Şablonlar (Templates)
- Standart görev yapıları
- Hızlı görev oluşturma
- Template alias sistemi (bug, feature, research vs.)

## 🤖 AI Asistan ile Kullanım

### Görev Yönetimi Komutları

```
"Yeni bir görev oluştur: API dokümantasyonu yaz"
"Bug raporu şablonundan görev oluştur: Login sorunu"
"Acil görevleri listele"
"bug etiketi olan görevleri göster"
"Mobile App v2 projesini aktif yap"
"5 numaralı görevi tamamlandı olarak işaretle"
```

### Proje Yönetimi

```
"Sprint planning için yeni proje oluştur"
"Aktif projedeki görevleri listele"
"Proje durumunu göster"
"Tüm projelerdeki görevleri listele"
```

### Gelişmiş Arama ve Filtreleme (v0.15.0+)

```
"API ile ilgili görevleri ara"
"Son 7 gündeki tamamlanan görevleri bul"
"Yüksek öncelikli bekleyen görevleri filtrele"
"Frontend etiketli devam eden görevleri göster"
```

### File Watching ve Otomatik Durum Geçişleri

```
"Proje dosyalarını izlemeye başla"
"Dosya değişikliklerinde otomatik durum geçişini etkinleştir"
"İzleme listesini göster"
"Git ignore kurallarını file watcher'a ekle"
```

## 🔧 Komut Satırı Kullanımı

### Sunucuyu Başlatma
```bash
# Normal modda başlat
gorev serve

# Debug modunda başlat
gorev serve --debug

# Belirli port ile
gorev serve --port 8080

# Türkçe dil ile
gorev serve --lang=tr
```

### Veritabanı Yönetimi
```bash
# Workspace veritabanı başlat (.gorev/gorev.db)
gorev init

# Global veritabanı başlat (~/.gorev/gorev.db)
gorev init --global

# Versiyon bilgisi
gorev version
```

### Template Alias Sistemi (v0.11.1+)
```bash
# Template alias'larını görüntüle
gorev template aliases

# Template'leri listele
gorev template list

# Belirli template'i incele
gorev template show bug
```

## 📋 Template Alias Referansı

### Mevcut Alias'lar (v0.15.5)

- **`bug`** → bug-report: Hata raporları ve düzeltmeler
- **`feature`** → feature: Yeni özellik ve geliştirmeler
- **`research`** → research: Araştırma ve inceleme görevleri
- **`spike`** → research: Teknik araştırma ve prototipler
- **`security`** → security: Güvenlik ile ilgili görevler
- **`performance`** → performance: Performans optimizasyonu
- **`refactor`** → refactor: Kod yeniden düzenleme
- **`debt`** → technical-debt: Teknik borç temizliği
- **`bug2`** → bug-report-v2: Gelişmiş bug raporu

### Kullanım Örnekleri

```bash
# AI asistan komutları:
"bug alias'ını kullanarak görev oluştur: Database connection timeout"
"feature şablonundan yeni görev: User profile page"
"research template'i ile analiz görevi oluştur"
```

## 🗂️ Görev Hiyerarşisi ve Bağımlılıklar

### Alt Görev Yönetimi

```
"Ana görev 15'e alt görev ekle: Frontend komponenti geliştir"
"Görev 23'ün alt görevlerini listele"
"Alt görev hiyerarşisini göster"
```

### Görev Bağımlılıkları

```
"Görev 10 görev 5'e bağımlı olsun"
"Bağımlılıkları olan görevleri listele"
"Görev 12'nin bağımlılıklarını göster"
```

## 📊 Raporlama ve İstatistikler

### Proje İstatistikleri

```
"Proje progress raporunu göster"
"Bu ayki tamamlanan görev sayısını ver"
"Öncelik dağılımını analiz et"
```

### Zaman Yönetimi

```
"Son tarih yaklaşan görevleri göster"
"Geciken görevleri listele"
"Bu hafta bitirilmesi gereken görevleri bul"
```

## 🔍 Gelişmiş Arama Özellikleri (v0.15.0)

### Fuzzy Search
- Yazım hatalarını tolere eden arama
- Benzer kelimeleri bulma
- Akıllı öneri sistemi

### Filter Profiles
- Kayıtlı arama filtrelerinin yönetimi
- Sık kullanılan filtreleri saklama
- Hızlı filtre uygulama

### Arama Geçmişi
- Önceki aramaları takip etme
- Arama istatistikleri
- Popüler arama terimlerini analiz

## 🔄 Data Export/Import (v0.12.0+)

### Veri Dışa Aktarma

```
"Tüm görevleri JSON formatında dışa aktar"
"Aktif projedeki görevleri CSV olarak çıkart"
"Sadece tamamlanan görevleri dışa aktar"
```

### Veri İçe Aktarma

```
"Backup dosyasından görevleri geri yükle"
"Dry run modunda import işlemini test et"
"Çakışan görevler için çözüm stratejisi belirle"
```

## 🛠️ IDE Entegrasyonu

### VS Code Extension

```
"VS Code uzantısını otomatik kur"
"Uzantı durumunu kontrol et"
"Extension'u güncelle"
```

### Multi-IDE Desteği
- VS Code
- Cursor
- Windsurf
- Claude Desktop
- Tüm MCP uyumlu editörler

## 🌐 Dil ve Yerelleştirme

### Dil Ayarları

```bash
# Çevre değişkeni ile
export GOREV_LANG=tr
gorev serve

# Komut satırı ile
gorev serve --lang=tr
```

### Desteklenen Diller
- **Türkçe (tr)**: Tam dil desteği
- **İngilizce (en)**: Tam dil desteği

## ⚡ Performans ve Thread Safety (v0.14.0+)

### Concurrent Access
- %100 thread-safe operasyonlar
- Race condition koruması
- Yüksek performanslı eşzamanlı erişim

### Memory Optimization
- %15-20 bellek tasarrufu
- Optimize edilmiş veritabanı sorguları
- %30 daha hızlı başlangıç

## 🔒 Güvenlik ve En İyi Pratikler

### Veritabanı Güvenliği
- SQLite encryption desteği
- Backup ve recovery prosedürleri
- Veri bütünlüğü kontrolü

### API Güvenliği
- MCP protokol standardları
- Güvenli parametre validasyonu
- Error handling best practices

## 🚨 Sorun Giderme

### Yaygın Sorunlar

**1. MCP Bağlantı Sorunu**
```bash
# Server durumunu kontrol et
gorev serve --debug

# Port kullanımını kontrol et
netstat -tlnp | grep 8080
```

**2. Veritabanı Kilit Sorunu**
```bash
# Server'ı yeniden başlat
pkill gorev
gorev serve
```

**3. VS Code Extension Çalışmıyor**
```
- VS Code'u yeniden başlat
- Developer: Reload Window komutunu çalıştır
- MCP server'ın çalıştığından emin ol
```

### Debug Yöntemleri

```bash
# Detaylı loglama
gorev serve --debug --log-level trace

# Veritabanı durumu
ls -la ~/.gorev/

# System information
gorev version --verbose
```

## 📚 İleri Seviye Kullanım

### Batch Operations
- Toplu görev işlemleri
- Mass update operasyonları
- Bulk import/export

### API Customization
- Özel MCP tool'lar
- Custom template'ler
- Workflow automation

### Integration Patterns
- CI/CD entegrasyonu
- Project management tools
- Time tracking systems

## 🔗 Faydalı Kaynaklar

- **[Kurulum Kılavuzu](kurulum.md)** - Detaylı kurulum talimatları
- **[MCP Araçları](mcp-araclari.md)** - 48 MCP tool referansı
- **[GitHub Repository](https://github.com/msenol/gorev)** - Kaynak kod ve issue'lar
- **[VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** - Extension
- **[Community Discussions](https://github.com/msenol/gorev/discussions)** - Topluluk desteği

---

*Bu dokümantasyon Claude (Anthropic) ile birlikte hazırlanmıştır*