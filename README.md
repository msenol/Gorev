# Gorev

Güçlü ve esnek görev yönetimi için Model Context Protocol (MCP) sunucusu.

## Özellikler

- 🎯 Hiyerarşik görev organizasyonu
- 📝 Markdown formatında görev açıklamaları
- ✏️ Esnek görev düzenleme (başlık, açıklama, öncelik, proje)
- 🗑️ Güvenli görev silme işlemleri
- 📁 Proje bazlı görev gruplandırma
- 🔄 Gerçek zamanlı senkronizasyon
- 📊 Zengin metadata desteği
- 🚀 Yüksek performanslı Go implementasyonu
- 🛠️ MCP protokolü ile AI entegrasyonu
- 🎪 **Aktif Proje Yönetimi** - Varsayılan proje seçimi ile hızlı görev oluşturma

## Kurulum

### Binary ile Kurulum

```bash
# En son sürümü indir
curl -L https://github.com/yourusername/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev
./gorev
```

### Kaynak Koddan Derleme

```bash
git clone https://github.com/yourusername/gorev.git
cd gorev
go build -o gorev cmd/gorev/main.go
./gorev
```

## Claude Desktop Konfigürasyonu

`claude_desktop_config.json` dosyanıza ekleyin:

```json
{
  "mcpServers": {
    "gorev": {
      "command": "/path/to/gorev",
      "args": ["serve"]
    }
  }
}
```

## Kullanım

### Temel Komutlar

```bash
# Sunucuyu başlat
gorev serve

# Versiyon bilgisi
gorev version

# Yardım
gorev help
```

### MCP İle Kullanım

Claude'a şu komutları verebilirsiniz:

#### Görev Yönetimi
- "Yeni bir görev oluştur"
- "Görevleri listele"
- "Görev detaylarını göster"
- "Görev bilgilerini düzenle"
- "Görevi sil"
- "Görev durumunu güncelle"

#### Proje Yönetimi
- "Yeni proje oluştur"
- "Projeleri listele"
- "Projenin görevlerini göster"
- "Proje özeti göster"
- "Projeyi aktif yap" - Varsayılan proje olarak ayarla
- "Aktif projeyi göster"
- "Aktif proje ayarını kaldır"

## Mimari

```
gorev/
├── cmd/gorev/        # Ana uygulama giriş noktası
├── internal/         # İç paketler
│   ├── mcp/         # MCP protokol implementasyonu
│   ├── gorev/       # İş mantığı
│   └── veri/        # Veri katmanı
├── pkg/             # Dışa açık paketler
│   ├── islem/       # İşlem yönetimi
│   └── sunum/       # Sunum katmanı
└── test/            # Test dosyaları
```

## Geliştirme

### Test Altyapısı

- **88.2% kod kapsama** oranı ile kapsamlı birim testleri
- Dependency injection pattern ile test edilebilir mimari
- SQL injection koruması testleri
- Concurrent erişim testleri
- Edge case validasyonları

```bash
# Testleri çalıştır
make test

# Test kapsama raporu oluştur
make test-coverage

# Race condition kontrolü
go test -race ./...

# Lint kontrolü
golangci-lint run

# Tüm platformlar için derle
make build-all
```

## Lisans

MIT License