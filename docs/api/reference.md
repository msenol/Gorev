# Gorev API Referansı

> **Versiyon**: Bu dokümantasyon v0.7.0-beta.1 için geçerlidir.  
> **Son Güncelleme**: 29 June 2025

Bu dokümanda Gorev'in programatik API'si, veri modelleri ve MCP protokol detayları açıklanmaktadır.

## İçindekiler

- [Veri Modelleri](#veri-modelleri)
- [MCP Protokol Detayları](#mcp-protokol-detayları)
- [Go API Kullanımı](#go-api-kullanımı)
- [Veritabanı Şeması](#veritabanı-şeması)
- [Hata Kodları](#hata-kodları)

## Veri Modelleri

### Gorev (Görev)

```go
type Gorev struct {
    ID              int               `json:"id"`
    Baslik          string            `json:"baslik"`
    Aciklama        string            `json:"aciklama"`
    Durum           string            `json:"durum"`
    Oncelik         string            `json:"oncelik"`
    ProjeID         *int              `json:"proje_id,omitempty"`
    OlusturmaTarih  time.Time         `json:"olusturma_tarih"`
    GuncellemeTarih time.Time         `json:"guncelleme_tarih"`
    SonTarih        *time.Time        `json:"son_tarih,omitempty"`
    Etiketler       []string          `json:"etiketler,omitempty"`
    Bagimliliklar   []GorevBagimlilik `json:"bagimliliklar,omitempty"`
}
```

**Alan Açıklamaları:**
- `ID`: Otomatik artan birincil anahtar
- `Baslik`: Görev başlığı (zorunlu, max 200 karakter)
- `Aciklama`: Detaylı açıklama (markdown destekli)
- `Durum`: `beklemede`, `devam_ediyor`, `tamamlandı` değerlerinden biri
- `Oncelik`: `dusuk`, `orta`, `yuksek` değerlerinden biri
- `ProjeID`: İlişkili proje ID (opsiyonel)
- `SonTarih`: Görev deadline'ı (opsiyonel)
- `Etiketler`: Görev etiketleri listesi
- `Bagimliliklar`: Bu göreve bağımlı olan görevler

### Proje

```go
type Proje struct {
    ID              int       `json:"id"`
    Isim            string    `json:"isim"`
    Tanim           string    `json:"tanim"`
    OlusturmaTarih  time.Time `json:"olusturma_tarih"`
    GuncellemeTarih time.Time `json:"guncelleme_tarih"`
    GorevSayisi     int       `json:"gorev_sayisi,omitempty"`
}
```

### GorevBagimlilik

```go
type GorevBagimlilik struct {
    ID          int    `json:"id"`
    KaynakID    int    `json:"kaynak_id"`
    HedefID     int    `json:"hedef_id"`
    BaglantiTip string `json:"baglanti_tip"`
    HedefGorev  *Gorev `json:"hedef_gorev,omitempty"`
}
```

**Bağlantı Tipleri:**
- `tamamlanmali`: Hedef görev tamamlanmadan kaynak görev başlayamaz
- `baslangic`: Hedef görev başlamadan kaynak görev başlayamaz

### GorevTemplate

```go
type GorevTemplate struct {
    ID                int                     `json:"id"`
    Isim              string                  `json:"isim"`
    Tanim             string                  `json:"tanim"`
    VarsayilanBaslik  string                  `json:"varsayilan_baslik"`
    AciklamaTemplate  string                  `json:"aciklama_template"`
    Alanlar          []TemplateAlan          `json:"alanlar"`
    OrnekDegerler    map[string]interface{}  `json:"ornek_degerler"`
    Kategori         string                  `json:"kategori"`
    Aktif            bool                    `json:"aktif"`
}
```

### TemplateAlan

```go
type TemplateAlan struct {
    Isim      string `json:"isim"`
    Tip       string `json:"tip"`       // text, number, select, date
    Zorunlu   bool   `json:"zorunlu"`
    Varsayilan string `json:"varsayilan,omitempty"`
    Secenekler []string `json:"secenekler,omitempty"` // select tipi için
}
```

## MCP Protokol Detayları

### Tool Schema Formatı

Her MCP tool'u aşağıdaki JSON Schema formatında tanımlanır:

```json
{
  "name": "tool_name",
  "description": "Tool açıklaması",
  "inputSchema": {
    "type": "object",
    "properties": {
      "param1": {
        "type": "string",
        "description": "Parametre açıklaması"
      }
    },
    "required": ["param1"]
  }
}
```

### Request/Response Formatı

**Request:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "gorev_olustur",
    "arguments": {
      "baslik": "Yeni görev",
      "aciklama": "Görev açıklaması",
      "oncelik": "orta"
    }
  }
}
```

**Response:**
```json
{
  "content": [
    {
      "type": "text",
      "text": "✅ Görev başarıyla oluşturuldu (ID: 123)"
    }
  ]
}
```

## Go API Kullanımı

### IsYonetici Interface

```go
type IsYonetici interface {
    // Görev işlemleri
    GorevOlustur(baslik, aciklama, oncelik string, projeID *int, sonTarihStr string, etiketIsimleri string) (*Gorev, error)
    GorevListele(durum string, tumProjeler bool, sirala, filtre, etiket string) ([]Gorev, error)
    GorevDetay(id int) (*Gorev, error)
    GorevGuncelle(id int, durum string) error
    GorevDuzenle(id int, baslik, aciklama, oncelik *string, projeID *int, sonTarihStr *string) error
    GorevSil(id int) error
    GorevBagimlilikEkle(kaynakID, hedefID int, baglantiTip string) error
    
    // Proje işlemleri
    ProjeOlustur(isim, tanim string) (*Proje, error)
    ProjeListele() ([]Proje, error)
    ProjeGorevleri(projeID int) (map[string][]Gorev, error)
    AktifProjeAyarla(projeID int) error
    AktifProjeGetir() (*Proje, error)
    AktifProjeKaldir() error
    
    // Template işlemleri
    TemplateListele(kategori string) ([]GorevTemplate, error)
    TemplatedenGorevOlustur(templateID int, degerler map[string]interface{}) (*Gorev, error)
    
    // Özet
    OzetGetir() (*Ozet, error)
}
```

### Örnek Kullanım

```go
package main

import (
    "github.com/msenol/gorev/internal/gorev"
    "log"
)

func main() {
    // Veri yöneticisi oluştur
    veriYonetici, err := gorev.YeniVeriYonetici("gorev.db", "migrations")
    if err != nil {
        log.Fatal(err)
    }
    defer veriYonetici.Kapat()
    
    // İş yöneticisi oluştur
    isYonetici := gorev.YeniIsYonetici(veriYonetici)
    
    // Görev oluştur
    gorev, err := isYonetici.GorevOlustur(
        "API entegrasyonu",
        "REST API implementasyonu",
        "yuksek",
        nil,
        "2025-02-01",
        "backend,api",
    )
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Görev oluşturuldu: %d - %s\n", gorev.ID, gorev.Baslik)
}
```

## Veritabanı Şeması

### Tablolar

**gorevler**
```sql
CREATE TABLE gorevler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    baslik TEXT NOT NULL,
    aciklama TEXT DEFAULT '',
    durum TEXT NOT NULL DEFAULT 'beklemede',
    oncelik TEXT NOT NULL DEFAULT 'orta',
    proje_id INTEGER,
    olusturma_tarih TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    guncelleme_tarih TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    son_tarih TIMESTAMP,
    FOREIGN KEY (proje_id) REFERENCES projeler(id) ON DELETE SET NULL,
    CHECK (durum IN ('beklemede', 'devam_ediyor', 'tamamlandı')),
    CHECK (oncelik IN ('dusuk', 'orta', 'yuksek'))
);
```

**projeler**
```sql
CREATE TABLE projeler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    isim TEXT NOT NULL UNIQUE,
    tanim TEXT DEFAULT '',
    olusturma_tarih TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    guncelleme_tarih TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**baglantilar**
```sql
CREATE TABLE baglantilar (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    kaynak_id INTEGER NOT NULL,
    hedef_id INTEGER NOT NULL,
    baglanti_tip TEXT NOT NULL DEFAULT 'tamamlanmali',
    FOREIGN KEY (kaynak_id) REFERENCES gorevler(id) ON DELETE CASCADE,
    FOREIGN KEY (hedef_id) REFERENCES gorevler(id) ON DELETE CASCADE,
    CHECK (baglanti_tip IN ('tamamlanmali', 'baslangic')),
    UNIQUE(kaynak_id, hedef_id)
);
```

**etiketler**
```sql
CREATE TABLE etiketler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    isim TEXT NOT NULL UNIQUE
);
```

**gorev_etiketleri**
```sql
CREATE TABLE gorev_etiketleri (
    gorev_id INTEGER NOT NULL,
    etiket_id INTEGER NOT NULL,
    PRIMARY KEY (gorev_id, etiket_id),
    FOREIGN KEY (gorev_id) REFERENCES gorevler(id) ON DELETE CASCADE,
    FOREIGN KEY (etiket_id) REFERENCES etiketler(id) ON DELETE CASCADE
);
```

### İndeksler

```sql
CREATE INDEX idx_gorevler_durum ON gorevler(durum);
CREATE INDEX idx_gorevler_proje ON gorevler(proje_id);
CREATE INDEX idx_gorevler_son_tarih ON gorevler(son_tarih);
CREATE INDEX idx_baglantilar_kaynak ON baglantilar(kaynak_id);
CREATE INDEX idx_baglantilar_hedef ON baglantilar(hedef_id);
```

## Hata Kodları

### MCP Hata Formatı

```go
mcp.NewToolResultError(fmt.Sprintf("Hata mesajı: %v", err))
```

### Yaygın Hatalar

| Hata | Açıklama | Çözüm |
|------|----------|-------|
| `gorev bulunamadı` | Belirtilen ID'ye sahip görev yok | Geçerli bir görev ID'si kullanın |
| `proje bulunamadı` | Belirtilen ID'ye sahip proje yok | Geçerli bir proje ID'si kullanın |
| `geçersiz durum` | Durum değeri geçersiz | beklemede, devam_ediyor, tamamlandı değerlerinden birini kullanın |
| `geçersiz öncelik` | Öncelik değeri geçersiz | dusuk, orta, yuksek değerlerinden birini kullanın |
| `bağımlılık döngüsü` | Döngüsel bağımlılık tespit edildi | Bağımlılık zincirini kontrol edin |
| `bağımlı görev tamamlanmamış` | Bağımlı görev henüz tamamlanmadı | Önce bağımlı görevi tamamlayın |

### Validation Hataları

- Başlık boş olamaz
- Başlık maksimum 200 karakter olabilir
- Tarih formatı YYYY-MM-DD olmalıdır
- Aynı görevler arasında birden fazla bağımlılık tanımlanamaz

## Özel Notlar

### Concurrency

- SQLite WAL mode kullanılır
- Okuma işlemleri paralel yapılabilir
- Yazma işlemleri serialize edilir

### Performans

- Görev listeleme için indeksler optimize edilmiştir
- Büyük projeler için sayfalama önerilir (henüz implement edilmemiş)

### Güvenlik

- SQL injection koruması: Prepared statements kullanılır
- Input validation: Tüm girişler validate edilir
- Rate limiting: MCP server seviyesinde handle edilmelidir

## API Değişiklikleri

Detaylı API değişiklikleri için [api-changes.md](api-changes.md) dosyasına bakın.

## İlgili Dokümantasyon

- [MCP Araçları](mcp-araclari.md)
- [Sistem Mimarisi](mimari.md)
- [Geliştirici Rehberi](gelistirme.md)

---

<div align="center">

*🔧 Bu API referans dokümantasyonu Claude (Anthropic) tarafından titizlikle yapılandırılmıştır - Teknik dokümantasyonda AI desteği*

</div>