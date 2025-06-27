# Gemini Proje Entegrasyon Kılavuzu: Gorev

Bu doküman, Gemini'nin `gorev` projesiyle etkileşim kurması için gerekli bağlamı ve komutları sağlar.

## Proje Özeti

`gorev`, Go dilinde yazılmış, Model Context Protocol (MCP) tabanlı bir görev yönetimi sunucusudur. Kullanıcıların doğal dil komutları aracılığıyla görevleri ve projeleri yönetmesine olanak tanır. Temel veri depolama birimi olarak SQLite kullanır.

**Anahtar Özellikler:**
- Hiyerarşik görev yönetimi (görevler projelere atanabilir).
- Görevler için başlık, açıklama, durum (`beklemede`, `devam_ediyor`, `tamamlandi`), ve öncelik (`dusuk`, `orta`, `yuksek`) gibi zengin metadata.
- Proje bazlı görev gruplaması.
- Aktif bir proje belirleyerek varsayılan olarak o projeye görev ekleme.
- MCP üzerinden AI asistanlarla (Gemini gibi) entegrasyon.

## Mimari ve Teknoloji

- **Dil:** Go (Golang)
- **Protokol:** Model Context Protocol (MCP)
- **Veritabanı:** SQLite
- **CLI Çerçevesi:** Cobra
- **Bağımlılıklar:** `mcp-go`, `go-sqlite3`, `cobra`, `uuid`.
- **Test:** Go'nun standart test kütüphanesi ve `testify`.

Proje, `internal` ve `pkg` dizinleri altında modüler bir yapıya sahiptir.
- `internal/gorev`: Ana iş mantığı ve veri yönetimi.
- `internal/mcp`: MCP sunucu implementasyonu.
- `cmd/gorev`: Komut satırı arayüzü.

## Temel Komutlar ve Kullanım Senaryoları

Gemini, `gorev` sunucusuyla etkileşim kurarken aşağıdaki komutları kullanabilir.

### 1. Görev Yönetimi

- **Yeni Görev Oluşturma:**
  - `yeni bir görev oluştur: "README dosyasını güncelle"`
  - `yüksek öncelikli bir görev oluştur: Başlık: "API bug'ını düzelt", Açıklama: "Kullanıcı girişi 500 hatası veriyor."`
  - `aktif projeye "testleri yaz" görevi ekle`

- **Görevleri Listeleme:**
  - `görevleri listele`
  - `beklemedeki görevleri göster`
  - `tamamlanmış görevleri listele`
  - `"API" projesindeki görevleri göster`

- **Görev Detaylarını Görüntüleme:**
  - `[görev-id] ID'li görevin detaylarını göster`
  - `başı "API" ile başlayan görevi bul`

- **Görevi Güncelleme:**
  - `[görev-id] görevinin başlığını "Yeni Başlık" olarak değiştir`
  - `[görev-id] görevinin durumunu "devam_ediyor" yap`
  - `[görev-id] görevinin önceliğini "yuksek" olarak ayarla`

- **Görevi Silme:**
  - `[görev-id] ID'li görevi sil`

### 2. Proje Yönetimi

- **Yeni Proje Oluşturma:**
  - `"Web Sitesi Yenileme" adında yeni bir proje oluştur`

- **Projeleri Listeleme:**
  - `tüm projeleri listele`

- **Aktif Proje Yönetimi:**
  - `"API Geliştirme" projesini aktif proje yap`
  - `aktif proje hangisi?`
  - `aktif proje ayarını kaldır`

### 3. Raporlama ve Özet

- `proje özetini göster`
- `genel durumu özetle` (toplam görev, duruma göre dağılım vb.)
- `yüksek öncelikli ve beklemede olan görevleri listele`

## Dosya Sistemi ve Yapılandırma

- **Veritabanı Dosyası:** Varsayılan olarak `~/.gorev/data/gorev.db` konumunda saklanır. Bu konum `--data-dir` parametresi ile değiştirilebilir.
- **Uygulama Binary'si:** `gorev` (veya Windows'ta `gorev.exe`).
- **Başlatma Komutu:** `gorev serve`

## Geliştirme ve Test

- **Testleri Çalıştırma:** `make test`
- **Kod Kapsamı Raporu:** `make test-coverage`
- **Linting:** `make lint`
- **Derleme:** `make build`
