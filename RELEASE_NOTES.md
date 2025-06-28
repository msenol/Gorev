# 🎉 Gorev v0.7.0-dev - İlk Public Release

## 🚀 Öne Çıkanlar

**Gorev**, MCP (Model Context Protocol) uyumlu AI editörlerle (Claude Desktop, VS Code, Windsurf, Cursor) entegre çalışan, Türkçe destekli modern bir görev yönetim sistemidir.

### 🌟 Ana Özellikler

- **16 MCP Tool** ile kapsamlı görev yönetimi
- **Go ile yazılmış** hızlı ve güvenilir MCP server
- **VS Code Extension** ile zengin görsel arayüz (opsiyonel)
- **%88.2 test coverage** ile yüksek kod kalitesi
- **Türkçe ve İngilizce** dokümantasyon

## 📦 İndirme

### Binary'ler

| Platform | Dosya | Boyut |
|----------|-------|-------|
| 🐧 Linux | [gorev-linux-amd64](https://github.com/msenol/Gorev/releases/download/v0.7.0-dev/gorev-linux-amd64) | 10.5 MB |
| 🍎 macOS | [gorev-darwin-amd64](https://github.com/msenol/Gorev/releases/download/v0.7.0-dev/gorev-darwin-amd64) | 7.0 MB |
| 🪟 Windows | [gorev-windows-amd64.exe](https://github.com/msenol/Gorev/releases/download/v0.7.0-dev/gorev-windows-amd64.exe) | 7.0 MB |

### Kurulum

```bash
# Linux/macOS
chmod +x gorev-*
sudo mv gorev-* /usr/local/bin/gorev

# Windows
# gorev-windows-amd64.exe dosyasını PATH'e ekleyin
```

## 🛠 MCP Araçları (16 Tool)

### Görev Yönetimi
- `gorev_olustur` - Yeni görev oluştur
- `gorev_listele` - Görevleri listele
- `gorev_detay` - Görev detaylarını göster
- `gorev_guncelle` - Görev durumunu güncelle
- `gorev_duzenle` - Görev özelliklerini düzenle
- `gorev_sil` - Görev sil
- `gorev_bagimlilik_ekle` - Görev bağımlılığı ekle

### Şablon Sistemi
- `template_listele` - Hazır şablonları listele
- `templateden_gorev_olustur` - Şablondan görev oluştur

### Proje Yönetimi
- `proje_olustur` - Yeni proje oluştur
- `proje_listele` - Projeleri listele
- `proje_gorevleri` - Proje görevlerini listele
- `proje_aktif_yap` - Aktif proje belirle
- `aktif_proje_goster` - Aktif projeyi göster
- `aktif_proje_kaldir` - Aktif proje ayarını kaldır

### Raporlama
- `ozet_goster` - Genel özet istatistikleri

## 🔧 Teknik Detaylar

- **Go Version**: 1.22+
- **MCP SDK**: mark3labs/mcp-go v0.6.0
- **Database**: SQLite (embedded)
- **Test Coverage**: %88.2

## 📚 Dokümantasyon

- [Kurulum Rehberi](docs/kurulum.md)
- [Kullanım Kılavuzu](docs/kullanim.md)
- [MCP Araçları Referansı](docs/mcp-araclari.md)
- [VS Code Extension](docs/vscode-extension.md)

## 🤝 Katkıda Bulunma

Projeye katkıda bulunmak için [Geliştirici Rehberi](docs/gelistirme.md)'ni inceleyin.

## 📄 Lisans

MIT Lisansı - Detaylar için [LICENSE](LICENSE) dosyasına bakın.

---

**Not**: Bu bir geliştirme sürümüdür (v0.7.0-dev). Üretim ortamında kullanmadan önce test edin.