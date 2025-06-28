# 🚀 Gorev

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)
![MCP](https://img.shields.io/badge/MCP-Compatible-4A154B?style=flat-square&logo=anthropic)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Test Coverage](https://img.shields.io/badge/Coverage-53.8%25-yellow?style=flat-square)
![Platform](https://img.shields.io/badge/Platform-Linux%20|%20macOS%20|%20Windows-blue?style=flat-square)

**MCP uyumlu AI editörlerle (Claude, VS Code, Windsurf, Cursor) entegre çalışan, Türkçe destekli modern görev yönetim sistemi**

[Özellikler](#-özellikler) • [Kurulum](#-kurulum) • [Kullanım](#-kullanım) • [Dokümantasyon](#-dokümantasyon) • [Katkıda Bulunma](#-katkıda-bulunma)

</div>

## 🎯 Gorev Nedir?

Gorev, **Model Context Protocol (MCP)** standardını kullanarak MCP uyumlu tüm AI editörler (Claude Desktop, VS Code, Windsurf, Cursor, Zed vb.) ile doğal dilde iletişim kurabilen, Go dilinde yazılmış güçlü bir görev yönetim sunucusudur. Proje yönetimi, görev takibi ve organizasyon ihtiyaçlarınızı AI asistanlarının yetenekleriyle birleştirerek verimliliğinizi artırır.

### 🏗️ İki Modüllü Yapı

1. **gorev-mcpserver** - Go dilinde yazılmış MCP server (ana bileşen)
2. **gorev-vscode** - VS Code extension'ı (opsiyonel görsel arayüz)

MCP protokolü sayesinde server'a herhangi bir MCP uyumlu editörden bağlanabilirsiniz. VS Code extension'ı ise zengin görsel deneyim sunar.

## ✨ Özellikler

### 📝 Görev Yönetimi
- **Akıllı görev oluşturma** - Doğal dil komutlarıyla
- **Markdown desteği** - Zengin açıklama formatı
- **Durum yönetimi** - Beklemede → Devam ediyor → Tamamlandı
- **Öncelik seviyeleri** - Düşük, Orta, Yüksek
- **Esnek düzenleme** - Tüm görev özelliklerini güncelleme

### 📁 Proje Organizasyonu
- **Hiyerarşik yapı** - Projeler altında görev gruplandırma
- **Aktif proje sistemi** - Varsayılan proje ile hızlı işlem
- **Proje bazlı raporlama** - Detaylı istatistikler
- **Çoklu proje desteği** - Sınırsız proje oluşturma

### 🔗 Gelişmiş Özellikler
- **📅 Son tarih takibi** - Deadline yönetimi ve acil görev filtreleme
- **🏷️ Etiketleme sistemi** - Çoklu etiket ile kategorilendirme
- **🔄 Görev bağımlılıkları** - İlişkili görevler arası otomasyon
- **📋 Hazır şablonlar** - Bug raporu, feature request ve daha fazlası
- **🔍 Gelişmiş filtreleme** - Durum, etiket, tarih bazlı sorgulama

### 🤖 AI Entegrasyonu
- **Doğal dil işleme** - AI asistanlarla konuşarak görev yönetimi
- **Çoklu editör desteği** - Claude, VS Code, Windsurf, Cursor, Zed
- **Bağlamsal anlama** - Akıllı komut yorumlama
- **MCP standardı** - Tüm MCP uyumlu araçlarla uyumluluk

### 🎨 VS Code Extension Özellikleri (Opsiyonel)
- **TreeView Panelleri** - Görev, proje ve şablon listeleri
- **Görsel Arayüz** - Tıkla ve kullan deneyimi
- **Status Bar** - Anlık durum bilgisi
- **Komut Paleti** - Hızlı erişim (Ctrl+Shift+G)
- **Renk Kodlaması** - Öncelik bazlı görsel ayırt etme
- **Context Menüler** - Sağ tık işlemleri

## 📦 Kurulum

### Hızlı Kurulum (30 saniye)

<details>
<summary><b>🪟 Windows</b></summary>

```powershell
# PowerShell (Admin olarak çalıştırın)
New-Item -ItemType Directory -Force -Path "C:\Program Files\gorev"
Invoke-WebRequest -Uri "https://github.com/yourusername/gorev/releases/latest/download/gorev-windows-amd64.exe" -OutFile "C:\Program Files\gorev\gorev.exe"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\gorev", [EnvironmentVariableTarget]::Machine)

# Test
gorev version
```

</details>

<details>
<summary><b>🍎 macOS</b></summary>

```bash
# Homebrew ile (önerilen)
brew tap yourusername/gorev
brew install gorev

# Veya binary indirme
curl -L https://github.com/yourusername/gorev/releases/latest/download/gorev-darwin-amd64 -o gorev
chmod +x gorev
sudo mv gorev /usr/local/bin/
```

</details>

<details>
<summary><b>🐧 Linux</b></summary>

```bash
# Binary indirme
curl -L https://github.com/yourusername/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev
sudo mv gorev /usr/local/bin/

# Test
gorev version
```

</details>

<details>
<summary><b>🐳 Docker</b></summary>

```bash
docker pull ghcr.io/yourusername/gorev:latest
docker run -v ~/.gorev:/data ghcr.io/msenol/gorev serve
```

</details>

### MCP Editör Entegrasyonu

<details>
<summary><b>🤖 Claude Desktop</b></summary>

Konfigürasyon dosyası konumları:
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "gorev": {
      "command": "/path/to/gorev-mcpserver/gorev",
      "args": ["serve"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev"
      }
    }
  }
}
```

</details>

<details>
<summary><b>💻 VS Code</b></summary>

#### Seçenek 1: Gorev VS Code Extension (Tavsiye Edilen)

1. **Extension'ı Yükleyin**:
   ```bash
   # Şu an local kurulum (marketplace yayını bekliyor)
   cd gorev-vscode
   npm install
   npm run compile
   # VS Code'da F5 ile çalıştırın
   ```

2. **Extension Ayarları** (`settings.json`):
   ```json
   {
     "gorev.serverPath": "/path/to/gorev-mcpserver/gorev",
     "gorev.autoConnect": true,
     "gorev.showStatusBar": true
   }
   ```

3. **Kullanım**:
   - Activity Bar'da Gorev ikonuna tıklayın
   - `Ctrl+Shift+G` ile hızlı görev oluşturun
   - TreeView'lardan görev/proje yönetin

#### Seçenek 2: MCP Extension ile

MCP extension kurduktan sonra `settings.json`:

```json
{
  "mcp.servers": {
    "gorev": {
      "command": "/path/to/gorev-mcpserver/gorev",
      "args": ["serve"]
    }
  }
}
```

</details>

<details>
<summary><b>🌊 Windsurf</b></summary>

Windsurf ayarlarında MCP sunucusu ekleyin:

```json
{
  "mcp.servers": [
    {
      "name": "gorev",
      "command": "gorev",
      "args": ["serve"]
    }
  ]
}
```

</details>

<details>
<summary><b>⚡ Cursor</b></summary>

Cursor'da MCP entegrasyonu için:

```json
{
  "mcp.gorev": {
    "command": "gorev serve",
    "env": {
      "GOREV_DATA_DIR": "~/.gorev"
    }
  }
}
```

</details>

## 🎮 Kullanım

### AI Asistan ile Örnek Komutlar

```
"Yeni bir görev oluştur: API dokümantasyonu yazılacak"
"Acil görevleri listele"
"Bug etiketli görevleri göster"
"Mobil App v2 projesini aktif yap"
"Sprint planlaması için yeni proje oluştur"
"Görev #5'i tamamlandı olarak işaretle"
"Feature request şablonundan yeni görev oluştur"
```

> 💡 **İpucu**: Bu komutlar Claude, VS Code Copilot, Windsurf AI, Cursor ve diğer MCP uyumlu AI asistanlarla kullanılabilir.

### CLI Komutları

```bash
# Server başlatma
gorev serve                  # Normal mod
gorev serve --debug         # Debug modunda
gorev serve --port 8080     # Farklı port

# Görev işlemleri
gorev task list             # Görevleri listele
gorev task create           # Yeni görev oluştur
gorev task show <id>        # Görev detayı

# Proje işlemleri
gorev project list          # Projeleri listele
gorev project create        # Yeni proje oluştur

# Diğer
gorev version              # Versiyon bilgisi
gorev help                 # Yardım
```

## 📚 Dokümantasyon

Detaylı dokümantasyon için [docs/](docs/) klasörüne bakın:

- 📦 [Kurulum Rehberi](docs/kurulum.md) - Platform spesifik kurulum talimatları
- 📖 [Kullanım Kılavuzu](docs/kullanim.md) - Detaylı kullanım örnekleri
- 🛠 [MCP Araçları](docs/mcp-araclari.md) - 16 MCP tool referansı
- 🏗 [Sistem Mimarisi](docs/mimari.md) - Teknik detaylar
- 💻 [Geliştirici Rehberi](docs/gelistirme.md) - Katkıda bulunma kılavuzu
- 🎨 [VS Code Extension](docs/vscode-extension.md) - Extension dokümantasyonu

## 🏗 Mimari

### Proje Yapısı

```
gorev/
├── gorev-mcpserver/        # MCP Server (Go)
│   ├── cmd/gorev/         # CLI ve server entry point
│   ├── internal/
│   │   ├── mcp/           # MCP protokol katmanı
│   │   └── gorev/        # Business logic
│   └── test/              # Integration testler
├── gorev-vscode/           # VS Code Extension (TypeScript)
│   ├── src/
│   │   ├── commands/      # VS Code komutları
│   │   ├── providers/     # TreeView sağlayıcıları
│   │   └── mcp/           # MCP client
│   └── package.json       # Extension manifest
└── docs/                   # Proje dokümantasyonu
```

### Bileşen Etkileşimi

```
┌───────────────┐     ┌───────────────┐     ┌────────────────┐
│ Claude/Cursor │     │   VS Code     │     │ VS Code + Gorev│
│               │     │ + MCP Plugin  │     │   Extension    │
└──────┬───────┘     └──────┬───────┘     └───────┬────────┘
       │                      │                      │
       └──────────────────────┴──────────────────────┘
                              │ MCP Protocol
                        ┌─────┴─────┐
                        │ Gorev MCP  │
                        │   Server   │
                        └─────┬─────┘
                              │
                        ┌─────┴─────┐
                        │   SQLite   │
                        └───────────┘
```

## 🧪 Geliştirme

### Gereksinimler
- Go 1.22+
- Make (opsiyonel)
- golangci-lint (kod kalitesi için)

### Komutlar

```bash
# Bağımlılıkları indir
make deps

# Test çalıştır (53.8% coverage)
make test

# Coverage raporu
make test-coverage

# Lint kontrolü
make lint

# Build (tüm platformlar)
make build-all

# Docker image
make docker-build
```

### Katkıda Bulunma

1. Projeyi fork'layın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit'leyin (`git commit -m 'feat: add amazing feature'`)
4. Branch'inizi push'layın (`git push origin feature/amazing-feature`)
5. Pull Request açın

Detaylı bilgi için [Geliştirici Rehberi](docs/gelistirme.md)'ne bakın.

## 📊 Proje Durumu

- **Versiyon**: v0.5.0
- **Test Coverage**: %88.2
- **Go Version**: 1.22+
- **MCP SDK**: mark3labs/mcp-go v0.6.0
- **Database**: SQLite (embedded)

## 🤝 Topluluk

- 📦 [GitHub Releases](https://github.com/msenol/gorev/releases)
- 🐛 [Issue Tracker](https://github.com/msenol/gorev/issues)
- 💬 [Discussions](https://github.com/msenol/gorev/discussions)
- 📖 [Wiki](https://github.com/msenol/gorev/wiki)

## 📄 Lisans

Bu proje [MIT Lisansı](LICENSE) altında lisanslanmıştır.

---

<div align="center">

Made with ❤️ by [Gorev Contributors](https://github.com/msenol/gorev/graphs/contributors)

📚 *Documentation enhanced by Claude (Anthropic) - Your AI pair programming assistant*

**[⬆ Başa Dön](#-gorev)**

</div>