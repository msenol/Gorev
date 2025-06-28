# API Değişiklikleri

Bu dokümanda Gorev API'sindeki önemli değişiklikler ve sürüm geçişleri açıklanmaktadır.

## [0.6.0-dev] - Geliştirme Aşamasında

### Planlanan Özellikler
- GitHub Actions entegrasyonu
- Binary release otomasyonu
- Docker registry kurulumu

## [0.5.0] - 2025-06-27

### Eklenen Özellikler

#### Görev Şablon Sistemi
- **Yeni tablo**: `gorev_templateleri` - Görev şablonlarını saklar
- **Yeni MCP araçları**:
  - `template_listele` - Mevcut şablonları listele
  - `templateden_gorev_olustur` - Şablondan görev oluştur
- **4 varsayılan şablon**:
  - Bug Raporu
  - Özellik İsteği
  - Teknik Borç
  - Araştırma Görevi
- **Dinamik alan desteği**: text, select, date, number
- **Alan doğrulama**: Zorunlu/opsiyonel alanlar

### Teknik Değişiklikler
- `GorevTemplate` ve `TemplateAlan` domain modelleri eklendi
- `template_yonetici.go` dosyası eklendi
- `VeriYoneticiInterface`'e 5 yeni template metodu eklendi

## [0.4.0] - 2025-06-27

### Eklenen Özellikler

#### Son Tarih Desteği
- **Veritabanı değişikliği**: `gorevler` tablosuna `son_tarih` kolonu eklendi
- **Güncellenmiş araçlar**:
  - `gorev_olustur` - `son_tarih` parametresi (YYYY-MM-DD formatında)
  - `gorev_duzenle` - Son tarih düzenleme desteği
  - `gorev_listele` - Yeni sıralama: `son_tarih_asc`, `son_tarih_desc`
  - `gorev_listele` - Yeni filtreler: `acil` (7 gün içinde), `gecmis` (geçmiş)

#### Etiketleme Sistemi
- **Yeni tablolar**:
  - `etiketler` - Etiket tanımları
  - `gorev_etiketleri` - Many-to-many ilişki tablosu
- **Güncellenmiş araçlar**:
  - `gorev_olustur` - `etiketler` parametresi (virgülle ayrılmış)
  - `gorev_listele` - `etiket` ile filtreleme
  - `gorev_detay` - Etiketleri gösterir

#### Görev Bağımlılıkları
- **Yeni MCP aracı**: `gorev_bagimlilik_ekle` - Görevler arası bağımlılık oluştur
- **İş mantığı**:
  - Bağımlı görevler tamamlanmadan "devam_ediyor" durumuna geçilemez
  - `GorevBagimliMi` fonksiyonu bağımlılık kontrolü yapar
- **Güncellenmiş araçlar**:
  - `gorev_guncelle` - Durum değişiminde bağımlılık kontrolü
  - `gorev_detay` - Bağımlılıkları durum göstergeleriyle gösterir (✅/⏳)

### Breaking Changes
- `GorevOlustur` fonksiyonu artık 6 parametre alıyor (son_tarih, etiketler eklendi)
- `GorevListele` fonksiyonu artık 3 parametre alıyor (sirala, filtre eklendi)
- `VeriYonetici` constructor artık migrations path gerektiriyor

## [0.3.0] - 2025-06-25

### Değişiklikler
- MCP SDK entegrasyonu tamamlandı (`mark3labs/mcp-go` v0.6.0)
- Go 1.22 minimum gereksinimi belirlendi
- Modül yolu `github.com/yourusername/gorev` olarak güncellendi

## [0.2.0] - 2025-06-24

### Eklenen Özellikler
- Aktif proje sistemi
- Proje bazlı görev filtreleme
- Özet istatistikleri

## [0.1.0] - 2025-06-23

### İlk Sürüm
- Temel görev yönetimi (CRUD)
- Proje yönetimi
- SQLite veritabanı desteği
- 10 temel MCP aracı

---

> Not: Semantic Versioning kullanılmaktadır. 1.0.0 sürümü public release için planlanmaktadır.