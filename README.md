# 🚀 Gorev

<div align="center">

**Last Updated:** July 18, 2025 | **Version:** v0.11.0

[🇺🇸 English](README.en.md) | [🇹🇷 Türkçe](README.md)

> ⚠️ **BREAKING CHANGE (v0.10.0)**: `gorev_olustur` tool artık kullanılmıyor! Template kullanımı zorunlu hale getirildi. [Detaylar](#breaking-change-template-zorunluluğu)

![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)
![MCP](https://img.shields.io/badge/MCP-Compatible-4A154B?style=flat-square&logo=anthropic)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Test Coverage](https://img.shields.io/badge/Coverage-84.6%25-brightgreen?style=flat-square)
![Platform](https://img.shields.io/badge/Platform-Linux%20|%20macOS%20|%20Windows-blue?style=flat-square)

**Modern task management system with Turkish support, designed for MCP-compatible AI assistants (Claude, VS Code, Windsurf, Cursor)**

[Özellikler](#-özellikler) • [Kurulum](#-kurulum) • [Kullanım](#-kullanım) • [Dokümantasyon](#-dokümantasyon) • [Katkıda Bulunma](#-katkıda-bulunma)

</div>

## 🌍 English Summary

**Gorev** is a powerful **Model Context Protocol (MCP)** server written in Go that provides task management capabilities to AI assistants (Claude, VS Code, Windsurf, Cursor). It features unlimited subtask hierarchy, dependency management, tagging system, and templates for structured task creation. 

**Key Features**: Natural language task creation, project organization, due date tracking, AI context management, 29 MCP tools, and optional VS Code extension with rich visual interface.

**Quick Start**: [Installation Guide](README.en.md#-installation) | [VS Code Extension](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)

---

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
- **🌳 Alt görev hiyerarşisi** - Sınırsız derinlikte görev ağacı yapısı
- **📊 İlerleme takibi** - Ana görevde alt görev tamamlanma yüzdesi
- **📁 File System Watcher** - Dosya değişikliklerini izleme ve otomatik görev durum geçişleri
- **🔔 Otomatik Durum Güncelleme** - Dosya değişikliklerinde "beklemede" → "devam_ediyor" otomasyonu
- **⚙️ Konfigürasyon Yönetimi** - Ignore patterns ve izleme kuralları özelleştirmesi

### 🤖 AI Entegrasyonu
- **Doğal dil işleme** - AI asistanlarla konuşarak görev yönetimi
- **Çoklu editör desteği** - Claude, VS Code, Windsurf, Cursor, Zed
- **Bağlamsal anlama** - Akıllı komut yorumlama
- **MCP standardı** - Tüm MCP uyumlu araçlarla uyumluluk
- **🌍 Uluslararası Destek** - Türkçe ve İngilizce tam dil desteği (v0.11.0+)

### 🎨 VS Code Extension Özellikleri (Opsiyonel)
- **İki Dil Desteği** - Türkçe ve İngilizce arayüz (v0.5.0+) 🌍
- **TreeView Panelleri** - Görev, proje ve şablon listeleri
- **Görsel Arayüz** - Tıkla ve kullan deneyimi
- **Status Bar** - Anlık durum bilgisi
- **Komut Paleti** - Hızlı erişim (Ctrl+Shift+G)
- **Renk Kodlaması** - Öncelik bazlı görsel ayırt etme
- **Context Menüler** - Sağ tık işlemleri
- **Otomatik Dil Algılama** - VS Code diline göre otomatik arayüz dili
- **[Marketplace'den İndir](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** 🚀
- **[Diğer IDE'lere Kurulum](docs/vscode-extension-kurulum.md)** (Cursor, Windsurf, VSCodium vb.)

## 📦 Kurulum

### 🚀 Otomatik Kurulum (Önerilen)

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash

# Specific version
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | VERSION=v0.11.0 bash
```

> **Not**: Kurulum sonrası eğer "unable to open database file" hatası alırsanız, GOREV_ROOT environment variable'ını ayarlayın:
> ```bash
> export GOREV_ROOT=/path/to/gorev-mcpserver  # Projenin bulunduğu dizin
> echo 'export GOREV_ROOT=/path/to/gorev-mcpserver' >> ~/.bashrc  # Kalıcı yapmak için
> ```

### Manuel Kurulum

<details>
<summary><b>🪟 Windows</b></summary>

**Otomatik Kurulum (PowerShell):**
```powershell
# PowerShell'de çalıştırın (Admin yetkisi gerekmez)
irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex

# Veya belirli versiyon için:
$env:VERSION="v0.10.0"; irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

**Manuel Kurulum:**
```powershell
# PowerShell (Admin olarak çalıştırın)
New-Item -ItemType Directory -Force -Path "C:\Program Files\gorev"
Invoke-WebRequest -Uri "https://github.com/msenol/gorev/releases/latest/download/gorev-windows-amd64.exe" -OutFile "C:\Program Files\gorev\gorev.exe"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\gorev", [EnvironmentVariableTarget]::Machine)

# GOREV_ROOT ayarla
[Environment]::SetEnvironmentVariable("GOREV_ROOT", "$env:APPDATA\gorev", [EnvironmentVariableTarget]::User)

# Test
gorev version
```

</details>

<details>
<summary><b>🍎 macOS</b></summary>

```bash
# Binary indirme (Homebrew desteği yakında)
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-darwin-amd64 -o gorev
chmod +x gorev
sudo mv gorev /usr/local/bin/
```

</details>

<details>
<summary><b>🐧 Linux</b></summary>

```bash
# Binary indirme
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev
sudo mv gorev /usr/local/bin/

# Test
gorev version
```

</details>

<details>
<summary><b>🐳 Docker</b></summary>

```bash
# Docker image yakında gelecek
# docker pull ghcr.io/msenol/gorev:latest
# docker run -v ~/.gorev:/data ghcr.io/msenol/gorev serve
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
        "GOREV_DATA_DIR": "~/.gorev",
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

> **🌍 Dil Desteği**: `GOREV_LANG` environment variable ile dil seçimi:
> - `tr` (varsayılan) - Türkçe
> - `en` - English
> 
> Alternatif olarak `--lang` parametresi: `gorev serve --lang=en`

</details>

<details>
<summary><b>💻 VS Code</b></summary>

#### Seçenek 1: Gorev VS Code Extension (Tavsiye Edilen)

1. **Extension'ı Yükleyin**:
   - **[VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** üzerinden
   - Veya komut satırından:
   ```bash
   code --install-extension mehmetsenol.gorev-vscode
   ```
   - Veya VS Code içinde: Extensions → "gorev" ara → Install

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
"Proje dosyalarını izlemeye başla"
"Dosya değişikliklerinde otomatik durum geçişi yap"
"Watch listesini göster"
"Git ignore kurallarını file watcher'a ekle"
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

- 📦 [Installation Guide](docs/getting-started/installation.md) - Platform-specific installation instructions
- 📖 [Usage Guide](docs/user-guide/usage.md) - Detailed usage examples
- 🛠 [MCP Tools](docs/user-guide/mcp-tools.md) - Complete reference for 29 MCP tools
- 🤖 [AI MCP Tools](docs/mcp-araclari-ai.md) - AI context management tools (v0.9.0)
- 🏗 [System Architecture](docs/development/architecture.md) - Technical details
- 🗺️ [Roadmap](ROADMAP.md) - Development roadmap and future plans
- 💻 [Contributing Guide](docs/development/contributing.md) - How to contribute
- 🎨 [VS Code Extension](docs/user-guide/vscode-extension.md) - Extension documentation
- 🚀 **[VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** - Extension'ı indir

### AI Assistant Documentation
- 🤖 [CLAUDE.md](CLAUDE.md) - Turkish AI assistant guidance
- 🌍 [CLAUDE.en.md](CLAUDE.en.md) - English AI assistant guidance
- 📋 [MCP Tools Reference](docs/MCP_TOOLS_REFERENCE.md) - Detailed MCP tool documentation
- 📚 [Development History](docs/DEVELOPMENT_HISTORY.md) - Complete project history

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

# Test çalıştır (84.6% overall coverage)
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

- **Versiyon**: v0.11.0
- **Test Coverage**: %84.6
- **Go Version**: 1.22+
- **MCP SDK**: mark3labs/mcp-go v0.6.0
- **Database**: SQLite (embedded)

## 🤝 Topluluk

- 📦 [GitHub Releases](https://github.com/msenol/gorev/releases)
- 🐛 [Issue Tracker](https://github.com/msenol/gorev/issues)
- 💬 [Discussions](https://github.com/msenol/gorev/discussions)
- 📖 [Wiki](https://github.com/msenol/gorev/wiki)

## ⚠️ BREAKING CHANGE: Template Zorunluluğu

### v0.10.0'dan İtibaren Template Kullanımı Zorunludur!

`gorev_olustur` tool artık kullanılamaz. Tüm görevler template kullanılarak oluşturulmalıdır.

#### 🔄 Eski Kullanım (Artık Çalışmaz):
```bash
gorev_olustur baslik="Bug fix" aciklama="..." oncelik="yuksek"
```

#### ✅ Yeni Kullanım (Zorunlu):
```bash
# 1. Önce template listesini görün
template_listele

# 2. Template kullanarak görev oluşturun
templateden_gorev_olustur template_id='bug_report_v2' degerler={
  'baslik': 'Login bug',
  'aciklama': 'Kullanıcı giriş yapamıyor',
  'modul': 'auth',
  'severity': 'high',
  ...
}
```

#### 🆕 Yeni Template'ler:
- **bug_report_v2** - Gelişmiş bug raporu (severity, steps, environment)
- **spike_research** - Time-boxed araştırma görevleri
- **performance_issue** - Performans sorunları (metrics, targets)
- **security_fix** - Güvenlik düzeltmeleri (CVSS, components)
- **refactoring** - Kod iyileştirme (code smell, strategy)

#### 🎯 Neden Template Zorunlu?
- **Tutarlılık**: Her görev belirli standartlara uygun
- **Kalite**: Zorunlu alanlar eksik bilgi girişini engeller
- **Otomasyon**: Template tipine göre otomatik workflow
- **Raporlama**: Görev tipine göre detaylı metrikler

## 📄 Lisans

Bu proje [MIT Lisansı](LICENSE) altında lisanslanmıştır.

---

<div align="center">

Made with ❤️ by [msenol](https://github.com/msenol/gorev/graphs/contributors)

📚 *Documentation enhanced by Claude (Anthropic) - Your AI pair programming assistant*

**[⬆ Başa Dön](#-gorev)**

</div>