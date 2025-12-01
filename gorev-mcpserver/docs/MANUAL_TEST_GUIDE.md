# Gorev - Manuel Test Rehberi

Bu rehber, Gorev uygulamasını gerçek bir kullanıcı gibi test etmek için adım adım senaryolar içerir.

## İçindekiler

1. [Kurulum ve Hazırlık](#kurulum-ve-hazırlık)
2. [Web UI Testleri](#web-ui-testleri)
3. [VS Code Extension Testleri](#vs-code-extension-testleri)
4. [CLI ve MCP Testleri](#cli-ve-mcp-testleri)
5. [Çapraz Arayüz Testleri](#çapraz-arayüz-testleri)
6. [Edge Case ve Hata Senaryoları](#edge-case-ve-hata-senaryoları)

---

## Kurulum ve Hazırlık

### Ön Gereksinimler

- Go 1.23+
- Node.js 18+
- VS Code (extension testleri için)
- Tarayıcı (Chrome, Firefox veya Safari)

### Test Ortamını Hazırlama

```bash
# 1. Projeyi derle
cd gorev-mcpserver
make build

# 2. Test veritabanını oluştur
./gorev seed-test-data

# 3. Sunucuyu başlat
./gorev serve --debug

# 4. Web UI'ı aç
open http://localhost:5082
```

### Test Verisi Seçenekleri

```bash
# Tam veri (15 görev, 3 proje, alt görevler, bağımlılıklar)
./gorev seed-test-data

# Minimal veri (3 görev, 1 proje)
./gorev seed-test-data --minimal

# İngilizce veri
./gorev seed-test-data --lang=en

# Mevcut veriyi silip yeniden oluştur
./gorev seed-test-data --force
```

---

## Web UI Testleri

### Test Ortamı
- URL: http://localhost:5082
- Sunucu: `./gorev serve` çalışıyor olmalı

### 1. Proje Yönetimi

#### WEB-PROJ-001: Proje Listesini Görüntüle
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Web UI'ı aç | Dashboard yüklenir |
| 2 | Sol kenar çubuğunu kontrol et | Proje listesi görünür |
| 3 | Her projenin yanındaki görev sayısını kontrol et | Görev sayıları doğru |

#### WEB-PROJ-002: Proje Oluştur
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | "+" butonuna tıkla (proje bölümünde) | Proje oluşturma formu açılır |
| 2 | Proje adı gir: "Test Projesi" | Alan dolar |
| 3 | Açıklama gir: "Test amaçlı proje" | Alan dolar |
| 4 | "Kaydet" butonuna tıkla | Proje listesinde görünür |

#### WEB-PROJ-003: Aktif Proje Değiştir
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Farklı bir projeye tıkla | Proje seçilir |
| 2 | Görev listesini kontrol et | Sadece seçilen projenin görevleri görünür |
| 3 | "Tüm Projeler"e tıkla | Tüm görevler görünür |

### 2. Görev CRUD İşlemleri

#### WEB-TASK-001: Görev Listesini Görüntüle
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Bir proje seç | Görev listesi yüklenir |
| 2 | Görev kartlarını kontrol et | Başlık, durum, öncelik görünür |
| 3 | Durum gruplarını kontrol et | Görevler duruma göre gruplandırılmış |

#### WEB-TASK-002: Şablondan Görev Oluştur
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | "Yeni Görev" butonuna tıkla | Şablon seçim modalı açılır |
| 2 | "Bug Raporu" şablonunu seç | Form alanları görünür |
| 3 | Başlık gir: "Test Bug" | Alan dolar |
| 4 | Açıklama gir | Alan dolar |
| 5 | Modül ve ortam bilgilerini gir | Alanlar dolar |
| 6 | "Oluştur" butonuna tıkla | Görev listesinde görünür |

#### WEB-TASK-003: Görev Durumunu Güncelle
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | "Beklemede" durumundaki bir görevi bul | Görev görünür |
| 2 | Durum dropdown'ını aç | Seçenekler: beklemede, devam_ediyor, tamamlandi |
| 3 | "devam_ediyor" seç | Durum güncellenir, renk değişir |
| 4 | Sayfayı yenile | Değişiklik kalıcı |

#### WEB-TASK-004: Görev Önceliğini Güncelle
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Herhangi bir görevi bul | Görev görünür |
| 2 | Öncelik badge'ına tıkla veya dropdown'ı aç | Seçenekler: düşük, orta, yüksek |
| 3 | "yuksek" seç | Badge kırmızıya döner |

#### WEB-TASK-005: Görev Sil
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Bir görevin "..." menüsüne tıkla | Menü açılır |
| 2 | "Sil" seçeneğine tıkla | Onay dialogu görünür |
| 3 | "Evet, Sil" tıkla | Görev listeden kaldırılır |

### 3. Alt Görevler

#### WEB-SUB-001: Alt Görevleri Görüntüle
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Alt görevi olan bir görev bul | Genişletme ikonu görünür |
| 2 | Genişletme ikonuna tıkla | Alt görevler görünür |
| 3 | Alt görev durumlarını kontrol et | Renkli noktalar: yeşil=tamamlandı, mavi=devam, gri=beklemede |

#### WEB-SUB-002: Alt Görev Oluştur
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Bir görevin menüsünü aç | Menü görünür |
| 2 | "Alt Görev Ekle" seç | Form açılır |
| 3 | Alt görev bilgilerini gir | Alanlar dolar |
| 4 | Kaydet | Alt görev ana görevin altında görünür |

### 4. Filtreleme ve Arama

#### WEB-FILTER-001: Duruma Göre Filtrele
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Filtre dropdown'ını aç | Durum seçenekleri görünür |
| 2 | "devam_ediyor" seç | Sadece devam eden görevler görünür |
| 3 | Filtreyi temizle | Tüm görevler görünür |

#### WEB-FILTER-002: Önceliğe Göre Filtrele
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Öncelik filtresini aç | Seçenekler: düşük, orta, yüksek |
| 2 | "yuksek" seç | Sadece yüksek öncelikli görevler görünür |

#### WEB-FILTER-003: Metin Ara
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Arama kutusuna "login" yaz | Sonuçlar filtrelenir |
| 2 | Eşleşen görevleri kontrol et | Başlık veya açıklamada "login" içerenler |

### 5. Dil Değiştirme

#### WEB-LANG-001: İngilizce'ye Geç
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Dil seçiciyi bul (genellikle sağ üst) | TR/EN seçenekleri |
| 2 | "EN" seç | UI etiketleri İngilizce olur |
| 3 | Sayfayı yenile | Ayar kalıcı |

---

## VS Code Extension Testleri

### Test Ortamı
- VS Code kurulu ve açık
- Gorev extension yüklü
- `./gorev daemon --detach` çalışıyor

### 1. Bağlantı Testleri

#### VSC-CONN-001: Otomatik Bağlantı
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Gorev daemon'u başlat | Daemon çalışır |
| 2 | VS Code'u aç | Extension yüklenir |
| 3 | Durum çubuğunu kontrol et | Yeşil bağlantı göstergesi |

#### VSC-CONN-002: Manuel Bağlantı
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Ctrl+Shift+P ile komut paletini aç | Palet açılır |
| 2 | "Gorev: Connect" yaz ve çalıştır | Bağlantı kurulur |
| 3 | Durum çubuğunu kontrol et | Yeşil gösterge |

### 2. TreeView Testleri

#### VSC-TREE-001: Görev TreeView'ı
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Gorev panelini aç | TreeView görünür |
| 2 | Görevleri kontrol et | Duruma göre gruplandırılmış |
| 3 | Öncelik renklerini kontrol et | Kırmızı=yüksek, sarı=orta, yeşil=düşük |

#### VSC-TREE-002: Proje TreeView'ı
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Projeler panelini aç | Proje listesi görünür |
| 2 | Her projenin yanındaki sayıyı kontrol et | Görev sayıları doğru |

#### VSC-TREE-003: Bağımlılık Badge'leri
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Bağımlılığı olan görevi bul | Badge görünür: "[link]3" |
| 2 | Tamamlanmamış bağımlılık varsa | Uyarı: "[link]2 warning1" |

### 3. Görev İşlemleri

#### VSC-TASK-001: Şablondan Görev Oluştur
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Ctrl+Shift+G veya komut paleti | Şablon seçici açılır |
| 2 | "Bug Raporu" şablonunu seç | Form wizard başlar |
| 3 | Alanları doldur | Her adımda sonraki alan |
| 4 | Tamamla | Görev TreeView'da görünür |

#### VSC-TASK-002: Durum Güncelle (Sağ Tık)
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Bir göreve sağ tıkla | Context menu açılır |
| 2 | "Durum Güncelle" seç | Durum seçenekleri |
| 3 | Yeni durum seç | Görev güncellenir |

#### VSC-TASK-003: Görev Detayını Görüntüle
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Bir göreve tıkla | Detay paneli açılır |
| 2 | Tüm bilgileri kontrol et | Başlık, açıklama, durum, öncelik, tarihler |
| 3 | Bağımlılıklar bölümünü kontrol et | İlgili görevler listelenir |

### 4. Klavye Kısayolları

#### VSC-KB-001: Hızlı Görev Oluştur
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Ctrl+Shift+G tuşlarına bas | Şablon seçici açılır |
| 2 | Şablon seç ve devam et | Görev oluşturma başlar |

#### VSC-KB-002: Filtreleri Temizle
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Bazı filtreler uygula | Görevler filtrelenir |
| 2 | Ctrl+Alt+R tuşlarına bas | Tüm filtreler temizlenir |

### 5. Data Export/Import

#### VSC-DATA-001: Veri Dışa Aktar
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Ctrl+Shift+P ile komut paletini aç | Palet açılır |
| 2 | "Gorev: Export Data" çalıştır | Format seçici |
| 3 | JSON seç | Dosya kaydedilir |
| 4 | Dosyayı kontrol et | Tüm görevler JSON formatında |

---

## CLI ve MCP Testleri

### 1. Temel CLI Komutları

#### CLI-001: Version
```bash
./gorev version
# Beklenen: Gorev v0.17.0, Build Time, Git Commit
```

#### CLI-002: Init
```bash
mkdir /tmp/test-workspace && cd /tmp/test-workspace
../gorev init
# Beklenen: .gorev/gorev.db oluşturulur
```

#### CLI-003: Serve
```bash
./gorev serve --debug
# Beklenen: Server başlar, Web UI http://localhost:5082'de açılır
```

#### CLI-004: Daemon Mode
```bash
./gorev daemon --detach
./gorev daemon-status
./gorev daemon-stop
# Beklenen: Her komut başarılı
```

### 2. Template Komutları

#### CLI-TMPL-001: Template Listele
```bash
./gorev template list
# Beklenen: Tüm şablonlar kategoriye göre listelenir
```

#### CLI-TMPL-002: Template Göster
```bash
./gorev template show bug
# Beklenen: Bug şablonunun detayları görünür
```

#### CLI-TMPL-003: Template Aliases
```bash
./gorev template aliases
# Beklenen: Kısa alias'lar listelenir (bug, feature, debt, etc.)
```

### 3. API Health Check

```bash
# Sunucu çalışırken
curl http://localhost:5082/api/health
# Beklenen: {"status":"healthy"} veya {"status":"ok"}

curl http://localhost:5082/api/v1/tasks
# Beklenen: JSON görev listesi

curl http://localhost:5082/api/v1/projects
# Beklenen: JSON proje listesi
```

---

## Çapraz Arayüz Testleri

### INT-SYNC-001: Web'den VS Code'a Senkronizasyon
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Web UI'da yeni görev oluştur | Görev kaydedilir |
| 2 | VS Code'da TreeView'ı yenile | Yeni görev görünür |

### INT-SYNC-002: VS Code'dan Web'e Senkronizasyon
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | VS Code'da görev durumunu değiştir | Durum güncellenir |
| 2 | Web UI'ı yenile | Değişiklik yansır |

### INT-SYNC-003: CLI'dan Arayüzlere
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | API ile görev oluştur (curl) | Görev kaydedilir |
| 2 | Web UI ve VS Code'u kontrol et | Her ikisinde de görünür |

---

## Edge Case ve Hata Senaryoları

### 1. Boş Durumlar

#### EDGE-EMPTY-001: Boş Veritabanı
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Yeni veritabanı oluştur | Boş DB |
| 2 | Web UI'ı aç | "Henüz proje yok" mesajı |
| 3 | VS Code'u aç | Boş TreeView |

### 2. Uzun Metin

#### EDGE-TEXT-001: Uzun Başlık
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | 500 karakterlik başlık gir | Kaydedilir |
| 2 | Görev kartını kontrol et | Başlık kesilir/wrap edilir |

#### EDGE-TEXT-002: Markdown İçerik
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Açıklamaya markdown gir | Kaydedilir |
| 2 | Detay panelinde kontrol et | Markdown render edilir |

### 3. Özel Karakterler

#### EDGE-CHAR-001: Unicode/Emoji
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Başlığa emoji ekle: "Test " | Kaydedilir |
| 2 | Görev listesinde kontrol et | Emoji görünür |

#### EDGE-CHAR-002: SQL Injection Denemesi
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Başlığa gir: `'; DROP TABLE gorevler; --` | Input sanitize edilir |
| 2 | Veritabanını kontrol et | Tablo sağlam |

### 4. Ağ Hataları

#### EDGE-NET-001: Sunucu Kapalı
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Sunucuyu durdur | Sunucu kapanır |
| 2 | Web UI'da işlem yap | Hata mesajı görünür |
| 3 | VS Code'da işlem yap | Bağlantı hatası |

### 5. Validasyon Hataları

#### EDGE-VAL-001: Boş Zorunlu Alan
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | Görev oluştur, başlığı boş bırak | Submit et |
| 2 | Sonucu kontrol et | Validasyon hatası |

#### EDGE-VAL-002: Döngüsel Bağımlılık
| Adım | Eylem | Beklenen Sonuç |
|------|-------|----------------|
| 1 | A görevi B'ye bağımlı yap | Kaydedilir |
| 2 | B görevi A'ya bağımlı yapmayı dene | "Döngüsel bağımlılık" hatası |

---

## Test Sonuç Raporlama

### Test Durumu Özeti

| Kategori | Toplam Test | Geçen | Başarısız | Atlandı |
|----------|-------------|-------|-----------|---------|
| Web UI | 21 | | | |
| VS Code Extension | 24 | | | |
| CLI/MCP | 21 | | | |
| Cross-Interface | 3 | | | |
| Edge Cases | 10 | | | |
| **TOPLAM** | **79** | | | |

### Hata Raporlama Formatı

Bir hata bulduğunuzda:

```markdown
## Bug Report

**Test ID:** WEB-TASK-002
**Tarih:** YYYY-MM-DD
**Ortam:** Chrome 120, macOS 14.1
**Adımlar:**
1. ...
2. ...

**Beklenen Sonuç:** X
**Gerçek Sonuç:** Y
**Ekran Görüntüsü:** [varsa ekle]
```

---

## Notlar

- Testler sırayla veya bağımsız olarak çalıştırılabilir
- Her test sonrası veritabanını sıfırlamak için: `./gorev seed-test-data --force`
- Debug logları için: `./gorev serve --debug`
- i18n testleri için dil değiştir: `GOREV_LANG=en ./gorev serve`
