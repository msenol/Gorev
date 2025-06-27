# 📋 Gorev Dokümantasyonu

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![MCP](https://img.shields.io/badge/MCP-Uyumlu-4A154B?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Linux%20|%20macOS%20|%20Windows-blue?style=for-the-badge)

**🚀 MCP uyumlu AI editörlerle (Claude, VS Code, Windsurf, Cursor) entegre çalışan, Türkçe destekli güçlü görev yönetim sistemi**

[Kurulum](#-kurulum) • [Özellikler](#-özellikler) • [Dokümantasyon](#-dokümantasyon) • [Örnekler](#-örnekler) • [Katkıda Bulunma](#-katkıda-bulunma)

</div>

---

## 🎯 Gorev Nedir?

Gorev, **Model Context Protocol (MCP)** standardını kullanarak tüm MCP uyumlu AI editörler (Claude Desktop, VS Code, Windsurf, Cursor, Zed vb.) ile sorunsuz entegre olan, Go dilinde yazılmış modern bir görev yönetim sunucusudur. Proje yönetimi, görev takibi ve organizasyon ihtiyaçlarınızı AI asistanlarının doğal dil yetenekleriyle birleştirerek güçlü bir üretkenlik aracı sunar.

### 🌟 Neden Gorev?

- **🤖 AI-Native**: MCP uyumlu tüm AI editörlerle doğal dilde görev yönetimi
- **🇹🇷 Türkçe Destek**: Arayüz ve komutlarda tam Türkçe desteği
- **⚡ Yüksek Performans**: Go ile yazılmış, minimal kaynak tüketimi
- **🔧 Kolay Kurulum**: Tek binary, sıfır bağımlılık
- **📊 Zengin Özellikler**: Şablonlar, bağımlılıklar, etiketleme ve daha fazlası

## 📚 Dokümantasyon

<table>
<tr>
<td width="33%" valign="top">

### 🚀 Başlangıç
- **[📦 Kurulum Rehberi](kurulum.md)**  
  Adım adım kurulum talimatları
  
- **[📖 Kullanım Kılavuzu](kullanim.md)**  
  Temel kullanım ve iş akışları
  
- **[💡 Örnekler](ornekler.md)**  
  Gerçek dünya senaryoları

</td>
<td width="33%" valign="top">

### 🔍 Referans
- **[🛠 MCP Araçları](mcp-araclari.md)**  
  16 MCP tool'unun detaylı referansı
  
- **[📡 API Dokümantasyonu](api-referans.md)**  
  Go API ve veri modelleri
  
- **[🏗 Sistem Mimarisi](mimari.md)**  
  Teknik tasarım ve yapı

</td>
<td width="33%" valign="top">

### 👩‍💻 Geliştirme
- **[💻 Geliştirici Rehberi](gelistirme.md)**  
  Katkıda bulunma kılavuzu
  
- **[📝 API Değişiklikleri](api-changes.md)**  
  Versiyon geçiş notları
  
- **[🐛 Sorun Giderme](kurulum.md#sorun-giderme)**  
  Yaygın sorunlar ve çözümleri

</td>
</tr>
</table>

## ⚡ Hızlı Başlangıç

### 1️⃣ Kurulum (30 saniye)

<details>
<summary><b>🐧 Linux / macOS</b></summary>

```bash
# Binary'yi indir
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev

# Sisteme kur (opsiyonel)
sudo mv gorev /usr/local/bin/
```

</details>

<details>
<summary><b>🪟 Windows</b></summary>

```powershell
# PowerShell ile indir
Invoke-WebRequest -Uri "https://github.com/msenol/gorev/releases/latest/download/gorev-windows-amd64.exe" -OutFile "gorev.exe"

# PATH'e ekle (opsiyonel)
$env:Path += ";$pwd"
```

</details>

<details>
<summary><b>🐳 Docker</b></summary>

```bash
docker run -v ~/.gorev:/data ghcr.io/msenol/gorev serve
```

</details>

### 2️⃣ MCP Editör Entegrasyonu

**Claude Desktop** için `claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "gorev": {
      "command": "/usr/local/bin/gorev",
      "args": ["serve"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev"
      }
    }
  }
}
```

**VS Code** için MCP extension ile:
```json
{
  "mcp.servers": {
    "gorev": {
      "command": "gorev",
      "args": ["serve"]
    }
  }
}
```

**Windsurf/Cursor** için ilgili MCP ayarlarını yapılandırın.

### 3️⃣ İlk Kullanım

AI asistanınıza (Claude, Copilot, Windsurf AI vb.) şunları söyleyebilirsiniz:

```
"Yeni bir proje oluştur: Mobil Uygulama v2"
"Bug raporu şablonundan yeni görev oluştur"
"Acil görevleri listele"
"Sprint planlaması yap"
```

## 🎨 Özellikler

### 📝 Görev Yönetimi
- ✅ Görev oluşturma, düzenleme, silme
- 🏷️ Çoklu etiketleme sistemi
- 📅 Son tarih ve aciliyet takibi
- 🔄 Durum yönetimi (beklemede → devam ediyor → tamamlandı)
- 🎯 Öncelik seviyeleri (düşük, orta, yüksek)

### 📁 Proje Organizasyonu
- 📊 Proje bazlı görev gruplama
- 🎯 Aktif proje yönetimi
- 📈 Proje bazlı istatistikler
- 🔀 Projeler arası görev taşıma

### 🔗 İleri Özellikler
- 🔄 Görev bağımlılıkları
- 📋 Özelleştirilebilir şablonlar
- 🔍 Gelişmiş filtreleme ve sıralama
- 📊 Detaylı raporlama

### 🤖 MCP Entegrasyonu
- 🗣️ Doğal dil komutları
- 🔧 16 özel MCP tool
- 📡 Gerçek zamanlı senkronizasyon
- 🔐 Güvenli veri yönetimi

## 📊 Desteklenen MCP Araçları

<details>
<summary><b>Tüm araçları göster</b></summary>

| Araç | Açıklama | Kategori |
|------|----------|----------|
| `gorev_olustur` | Yeni görev oluşturur | Görev |
| `gorev_listele` | Görevleri filtreler ve listeler | Görev |
| `gorev_detay` | Görev detaylarını gösterir | Görev |
| `gorev_guncelle` | Görev durumunu günceller | Görev |
| `gorev_duzenle` | Görev özelliklerini düzenler | Görev |
| `gorev_sil` | Görevi siler | Görev |
| `gorev_bagimlilik_ekle` | Görev bağımlılığı tanımlar | Görev |
| `template_listele` | Mevcut şablonları listeler | Şablon |
| `templateden_gorev_olustur` | Şablondan görev oluşturur | Şablon |
| `proje_olustur` | Yeni proje oluşturur | Proje |
| `proje_listele` | Tüm projeleri listeler | Proje |
| `proje_gorevleri` | Proje görevlerini gösterir | Proje |
| `proje_aktif_yap` | Aktif projeyi değiştirir | Proje |
| `aktif_proje_goster` | Aktif projeyi gösterir | Proje |
| `aktif_proje_kaldir` | Aktif proje ayarını kaldırır | Proje |
| `ozet_goster` | Genel istatistikleri gösterir | Rapor |

</details>

## 🛠️ Teknik Özellikler

- **Dil**: Go 1.22+
- **Veritabanı**: SQLite3 (embedded)
- **Protokol**: MCP (Model Context Protocol)
- **SDK**: mark3labs/mcp-go v0.6.0
- **Platform**: Linux, macOS, Windows
- **Mimari**: Clean Architecture, Domain-Driven Design

## 🤝 Katkıda Bulunma

Gorev'e katkıda bulunmak ister misiniz? Harika! 

1. 🍴 Projeyi fork'layın
2. 🌿 Feature branch oluşturun (`git checkout -b ozellik/harika-ozellik`)
3. 💾 Değişikliklerinizi commit'leyin (`git commit -m 'feat: harika özellik ekle'`)
4. 📤 Branch'inizi push'layın (`git push origin ozellik/harika-ozellik`)
5. 🔄 Pull Request açın

Detaylı bilgi için [Geliştirici Rehberi](gelistirme.md)'ne bakın.

## 📈 Proje Durumu

- **Versiyon**: v0.5.0
- **Durum**: Aktif Geliştirme
- **Son Güncelleme**: Haziran 2025
- **Test Coverage**: %88.2

## 🔗 Bağlantılar

- 📦 [GitHub Repository](https://github.com/msenol/gorev)
- 🐛 [Sorun Bildirme](https://github.com/msenol/gorev/issues)
- 💬 [Tartışmalar](https://github.com/msenol/gorev/discussions)
- 📖 [MCP Protokolü](https://modelcontextprotocol.io)

## 📄 Lisans

Bu proje [MIT Lisansı](../LICENSE) altında lisanslanmıştır.

---

<div align="center">

**[⬆ Başa Dön](#-gorev-dokümantasyonu)**

Made with ❤️ by [Gorev Contributors](https://github.com/msenol/gorev/graphs/contributors)

📚 *Documentation crafted with assistance from Claude (Anthropic) - Your AI documentation partner*

</div>