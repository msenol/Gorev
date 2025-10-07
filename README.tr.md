# 🚀 Gorev

<div align="center">

**Last Updated:** October 6, 2025 | **Version:** v0.16.3

[🇺🇸 English](README.en.md) | [🇹🇷 Türkçe](README.md)

> 🎉 **YENİ v0.16.3**: MCP araç parametre dönüşümü düzeltmeleri + %100 test başarısı! [Yeniliklere Bak](#-v0163-yenilikleri)

> ⚠️ **BREAKING CHANGE (v0.10.0)**: `gorev_olustur` tool artık kullanılmıyor! Template kullanımı zorunlu hale getirildi. [Detaylar](#breaking-change-template-zorunluluğu)

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)
![MCP](https://img.shields.io/badge/MCP-Compatible-4A154B?style=flat-square&logo=anthropic)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Test Coverage](https://img.shields.io/badge/Coverage-71%25-yellow?style=flat-square)
![Platform](https://img.shields.io/badge/Platform-Linux%20|%20macOS%20|%20Windows-blue?style=flat-square)

**Modern task management system with Turkish support, designed for MCP-compatible AI assistants (Claude, VS Code, Windsurf, Cursor)**

[Özellikler](#-özellikler) • [Kurulum](#-kurulum) • [Kullanım](#-kullanım) • [Dokümantasyon](#-dokümantasyon) • [Katkıda Bulunma](#-katkıda-bulunma)

</div>

## 🌍 English Summary

**Gorev** is a powerful **Model Context Protocol (MCP)** server written in Go that provides task management capabilities to AI assistants (Claude, VS Code, Windsurf, Cursor). It features unlimited subtask hierarchy, dependency management, tagging system, and templates for structured task creation.

**Key Features**: Natural language task creation, project organization, due date tracking, AI context management, enhanced NLP processing, advanced search & filtering with FTS5, 24 optimized MCP tools (unified from 45), and optional VS Code extension with rich visual interface.

**Quick Start**: [Installation Guide](README.en.md#-installation) | [VS Code Extension](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)

---

## 🎯 Gorev Nedir

Gorev, **Model Context Protocol (MCP)** standardını kullanarak MCP uyumlu tüm AI editörler (Claude Desktop, VS Code, Windsurf, Cursor, Zed vb.) ile doğal dilde iletişim kurabilen, Go dilinde yazılmış güçlü bir görev yönetim sunucusudur. Proje yönetimi, görev takibi ve organizasyon ihtiyaçlarınızı AI asistanlarının yetenekleriyle birleştirerek verimliliğinizi artırır.

### 🏗️ Üç Modüllü Yapı

1. **gorev-mcpserver** - Go dilinde yazılmış MCP server (ana bileşen)
   - Embedded Web UI 🌐 - React arayüzü binary'ye gömülü (YENİ! v0.16.0)
   - REST API server (Fiber framework)
   - MCP protokol desteği
2. **gorev-vscode** - VS Code extension'ı (opsiyonel görsel arayüz)
3. **gorev-web** - React + TypeScript kaynak kodu (development)

MCP protokolü sayesinde server'a herhangi bir MCP uyumlu editörden bağlanabilirsiniz. Web arayüzü `npx @mehmetsenol/gorev-mcp-server serve` komutuyla otomatik olarak http://localhost:5082 adresinde hazır olur. VS Code extension'ı ise IDE içinde zengin görsel deneyim sunar.

### 🔌 Daemon Mimarisi (v0.16.0+)

Gorev, **arka plan daemon process** olarak çalışır ve şu avantajları sağlar:

**Temel Özellikler:**

- **Tek Instance Yönetimi**: Lock dosyası (`~/.gorev-daemon/.lock`) port çakışmalarını önler
- **Çoklu İstemci Desteği**: Birden fazla MCP istemcisi (Claude, VS Code, Windsurf, Cursor) aynı anda bağlanabilir
- **Otomatik Başlatma**: VS Code extension daemon'u otomatik tespit eder ve başlatır (v0.16.2+)
- **Sağlık İzleme**: `/api/health` endpoint'i ile gerçek zamanlı durum kontrolü
- **WebSocket Desteği**: Gerçek zamanlı görev güncelleme olayları (deneysel)

**Hızlı Başlangıç:**

```bash
# Daemon'u arka planda başlat
gorev daemon --detach

# Daemon durumunu kontrol et
curl http://localhost:5082/api/health

# Web arayüzü otomatik olarak hazır
open http://localhost:5082
```

**Mimari Bileşenler:**

- **Lock Dosyası**: `~/.gorev-daemon/.lock` PID, port, versiyon ve daemon URL içerir
- **REST API Server**: VS Code extension için 23 endpoint (Fiber framework)
- **MCP Proxy**: stdio MCP protokol isteklerini internal handler'lara yönlendirir
- **WebSocket Server**: Görev güncellemeleri için gerçek zamanlı olay yayını
- **Workspace Manager**: SHA256 tabanlı ID'lerle çoklu workspace desteği

**VS Code Entegrasyonu:**
Extension daemon yaşam döngüsünü otomatik yönetir:

1. Aktivasyonda daemon'un çalışıp çalışmadığını kontrol eder (lock dosyasını okur)
2. Çalışmıyorsa daemon'u başlatır
3. Tüm işlemler için REST API'ye bağlanır
4. Deaktivasyonda daemon'u kapatır (eğer extension başlattıysa)

Detaylı teknik özellikler için [Daemon Mimari Dokümantasyonu](docs/architecture/daemon-architecture.md)'na bakın.

## 🎉 v0.16.3 Yenilikleri

### 🔧 MCP Araç Parametre Dönüşüm Düzeltmeleri (6 Ekim 2025)

**gorev_bulk** - Tüm 3 operasyon artık tamamen çalışıyor:

- **`update` operasyonu**: `{ids: [], data: {}}` → `{updates: [{id, ...alanlar}]}` dönüşümü düzgün çalışıyor
- **`transition` operasyonu**: Hem `durum` hem `yeni_durum` parametrelerini kabul ediyor
- **`tag` operasyonu**: Hem `operation` hem `tag_operation` parametrelerini kabul ediyor
- **Test sonucu**: %100 başarı oranı (5/5 operasyon production'da test edildi)

**gorev_guncelle** - Çoklu alan güncelleme desteği eklendi:

- `durum` (durum), `oncelik` (öncelik) veya her ikisini birden güncelleyebilir
- En az bir parametre gerekli (validasyon)
- Mevcut kodla geriye dönük uyumlu

**gorev_search (gelişmiş mod)** - Akıllı sorgu ayrıştırma eklendi:

- **Örnek**: `"durum:devam_ediyor oncelik:yuksek tags:frontend"`
- Doğal dil sorgularından filtreleri otomatik olarak çıkarır
- Boşlukla ayrılmış key:value çiftleri ile çoklu filtre desteği
- Mevcut filtre parametreleriyle sorunsuz çalışır

**VS Code Tree View** - Bağımlılık göstergeleri artık görünür:

- 🔒 (bloke), 🔓 (bloke değil), 🔗 (bağımlı) ikonları düzgün gösteriliyor
- JSON serileştirme sorunu düzeltildi (bağımlılık sayaçlarından `omitempty` kaldırıldı)
- Tüm bağımlılık ilişkileri artık tree yapısında görünür

**Doğrulama**: Kilocode AI kapsamlı test raporu ile %100 başarı oranı onaylandı

---

### 🐛 Önceki Güncellemeler (v0.16.2 - 5 Ekim 2025)

- **NPM Binary Güncelleme Hatası**: NPM paket yükseltmelerinde eski binary'lerin korunması hatası düzeltildi
  - Paket boyutu 78.4 MB'tan 6.9 KB'ye düşürüldü (binary'ler artık GitHub'dan indiriliyor)
- **VS Code Otomatik Başlatma**: Extension artık server'ı otomatik olarak başlatıyor

### 🌐 Embedded Web UI (v0.16.0)

- **Sıfır Yapılandırma**: Modern React arayüzü Go binary'sine gömülü
- **Anında Erişim**: http://localhost:5082 adresinde otomatik olarak hazır
- **Tam Özellikler**: Görevler, projeler, şablonlar, alt görevler ve bağımlılıklar
- **Dil Senkronizasyonu**: Türkçe/İngilizce değiştirici MCP server ile senkronize
- **Ayrı Kurulum Yok**: Sadece `npx @mehmetsenol/gorev-mcp-server serve` komutuyla hazır!

### 🗂️ Çoklu Workspace Desteği (v0.16.0)

- **İzole Workspace'ler**: Her proje klasörü kendi görev veritabanına sahip
- **Workspace Değiştirici**: Web UI'da workspace'ler arası sorunsuz geçiş
- **Otomatik Tespit**: Mevcut klasördeki `.gorev/` dizinini otomatik algılar

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

- **🔍 Gelişmiş Arama & Filtreleme** - FTS5 full-text search ve akıllı filtreleme (v0.15.0)
- **🧠 Gelişmiş NLP İşleme** - Akıllı doğal dil anlama ve yorumlama (v0.14.0)
- **Doğal dil işleme** - AI asistanlarla konuşarak görev yönetimi
- **Çoklu editör desteği** - Claude, VS Code, Windsurf, Cursor, Zed
- **Bağlamsal anlama** - Akıllı komut yorumlama ve parametre çıkarımı
- **MCP standardı** - Tüm MCP uyumlu araçlarla uyumluluk
- **🌐 MCP Registry** - Otomatik keşfedilebilirlik ve kolay kurulum (v0.15.24+)
- **🌍 Uluslararası Destek** - Türkçe ve İngilizce tam dil desteği (v0.11.0+)
- **⚡ Thread Safety** - 100% race condition free operation (v0.14.0)

### 🌐 Web UI Özellikleri (YENİ! v0.16.0)

- **Modern React Arayüzü** - TypeScript + Vite ile hızlı ve responsive
- **Proje Bazlı Görünüm** - Proje kartları ve gerçek zamanlı istatistikler
- **Görev Yönetimi** - CRUD işlemleri template sistemi ile
- **Alt Görev Görünümü** - Hiyerarşik görev listesi (collapse/expand)
- **Bağımlılık Göstergesi** - Visual dependency indicators (🔗 count + ⚠️ incomplete)
- **Durum Yönetimi** - Inline dropdown'larla hızlı güncelleme
- **Gelişmiş Filtreleme** - Durum, öncelik, proje bazlı filtreleme
- **🌍 Dil Değiştirici** - Türkçe/İngilizce arasında geçiş, MCP sunucusu ile senkronize
- **Responsive Tasarım** - Tailwind CSS ile mobil uyumlu
- **Gerçek Zamanlı Sync** - React Query ile otomatik veri güncelleme
- **🚀 Kurulum Gerektirmez**: `npx @mehmetsenol/gorev-mcp-server serve` komutuyla otomatik aktif!
- **Embedded UI**: Go binary'sine gömülü, ayrı kurulum yok
- **Erişim**: http://localhost:5082 (varsayılan port)

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
- **[Diğer IDE'lere Kurulum](docs/legacy/tr/vscode-extension-kurulum.md)** (Cursor, Windsurf, VSCodium vb.)

## 📦 Kurulum

### 🚀 NPM ile Kolay Kurulum (Önerilen!)

> ⚠️ **Windows Kullanıcıları**: NPM kullanımı için Node.js kurulumu gereklidir. [Node.js indirin](https://nodejs.org/) ve yükledikten sonra sistemi yeniden başlatın.

#### Global Kurulum

```bash
npm install -g @mehmetsenol/gorev-mcp-server
```

Veya doğrudan NPX ile kullanın (kurulum gerektirmez):

```bash
npx @mehmetsenol/gorev-mcp-server serve
```

#### MCP Client Konfigürasyonu

**Claude Desktop için:**

```json
// ~/.config/Claude/claude_desktop_config.json (Linux)
// ~/Library/Application Support/Claude/claude_desktop_config.json (macOS)
// %APPDATA%/Claude/claude_desktop_config.json (Windows)
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest",
        "mcp-proxy"
      ],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

**Kilo Code (VS Code Extension) için:**

```json
// .kilocode/mcp.json (workspace root)
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest",
        "mcp-proxy"
      ],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

**Cursor için:**

```json
// .cursor/mcp.json (workspace root)
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": [
        "-y",
        "@mehmetsenol/gorev-mcp-server@latest",
        "mcp-proxy"
      ],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

> 📚 **Daha fazla örnek**: [MCP Configuration Examples](docs/guides/mcp-config-examples.md)

#### 🔧 Windows NPX Kurulum Sorunu Giderme

Eğer `ENOENT: spawn npx` hatası alıyorsanız:

1. **Node.js kurulu mu kontrol edin:**

   ```cmd
   node --version
   npm --version
   npx --version
   ```

2. **Node.js kurulumunu yapın:**
   - [Node.js websitesinden](https://nodejs.org/) LTS sürümü indirin
   - Installer'ı çalıştırırken "Add to PATH" seçeneğini işaretleyin
   - Kurulum sonrası bilgisayarı yeniden başlatın

3. **NPX ayrı kurulumu (gerekirse):**

   ```cmd
   npm install -g npx
   ```

4. **PATH kontrolü:**

   ```cmd
   echo %PATH%
   ```

   Node.js paths (`C:\Program Files\nodejs\`) görünmeli.

### 🔧 Geleneksel Kurulum (Otomatik)

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash

# Specific version
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | VERSION=v0.11.0 bash
```

> **Not**: Kurulum sonrası eğer "unable to open database file" hatası alırsanız, GOREV_ROOT environment variable'ını ayarlayın:
>
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
$env:VERSION="v0.15.4"; irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
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
>
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

### AI Asistan ile Örnek Komutlar (v0.14.0 Enhanced NLP)

```
"Yeni bir görev oluştur: API dokümantasyonu yazılacak yarın deadline ile"
"Bu hafta için yüksek öncelikli görevleri göster"
"Bug etiketli açık görevleri listele"
"Mobil App v2 projesini aktif yap"
"Sprint planlaması için yeni proje oluştur"
"Görev #5'i tamamlandı olarak işaretle"
"Son oluşturduğum görev nasıl gidiyor?"
"Frontend kategorisindeki görevleri göster"
"Bugün deadline olan acil görevler var mı?"
"Feature request şablonundan yeni görev oluştur"
"Proje dosyalarını izlemeye başla"
"Dosya değişikliklerinde otomatik durum geçişi yap"
"Watch listesini göster"
```

> 💡 **İpucu**: Bu komutlar Claude, VS Code Copilot, Windsurf AI, Cursor ve diğer MCP uyumlu AI asistanlarla kullanılabilir.

### CLI Komutları

```bash
# Daemon yönetimi (önerilen)
gorev daemon --detach        # Daemon'u arka planda başlat
gorev daemon-status          # Daemon durumunu kontrol et
gorev daemon-stop            # Çalışan daemon'u durdur
gorev mcp-proxy              # MCP proxy (AI asistanlar için)

# Geliştirme/test (foreground mod)
gorev serve                  # Normal mod
gorev serve --debug          # Debug modunda
gorev serve --port 8080      # Farklı port

# Görev işlemleri
gorev task list              # Görevleri listele
gorev task create            # Yeni görev oluştur
gorev task show <id>         # Görev detayı

# Proje işlemleri
gorev project list           # Projeleri listele
gorev project create         # Yeni proje oluştur

# Diğer
gorev version                # Versiyon bilgisi
gorev help                   # Yardım
```

## 📚 Dokümantasyon

Detaylı dokümantasyon için [docs/](docs/) klasörüne bakın:

### Başlangıç

- 🚀 [Hızlı Başlangıç](docs/guides/getting-started/quick-start.md) - 10 dakikada kurulum ve kullanım
- 📦 [Kurulum Kılavuzu](docs/guides/getting-started/installation.md) - Platform-specific installation
- 🆘 [Sorun Giderme](docs/guides/getting-started/troubleshooting.md) - Yaygın sorunlar ve çözümleri
- 🔄 [Göç Kılavuzu (v0.15→v0.16)](docs/migration/v0.15-to-v0.16.md) - v0.15'ten yükseltme

### Özellikler

- 🌐 [Web UI Kılavuzu](docs/guides/features/web-ui.md) - Gömülü React arayüz dokümantasyonu
- 🗂️ [Multi-Workspace Desteği](docs/guides/features/multi-workspace.md) - Çoklu proje yönetimi
- 📋 [Template Sistemi](docs/guides/features/template-system.md) - Yapılandırılmış görev oluşturma
- 🤖 [AI Context Yönetimi](docs/guides/features/ai-context-management.md) - AI asistan entegrasyonu

### Referans

- 🛠️ [MCP Araçları](gorev-mcpserver/docs/mcp-araclari.md) - 24 optimize MCP aracının komple referansı (45'ten birleştirildi)
- 🔧 [MCP Konfigürasyon Örnekleri](docs/guides/mcp-config-examples.md) - IDE kurulum kılavuzları
- 📖 [Kullanım Kılavuzu](docs/guides/user/usage.md) - Detaylı kullanım örnekleri
- 🎨 [VS Code Extension](docs/guides/user/vscode-extension.md) - Extension dokümantasyonu

### Geliştirme

- 🏗️ [Sistem Mimarisi](docs/architecture/architecture-v2.md) - Teknik detaylar
- 💻 [Katkı Kılavuzu](CONTRIBUTING.md) - Nasıl katkıda bulunulur
- 🗺️ [Yol Haritası](ROADMAP.md) - Geliştirme planları
- 📚 [Geliştirme Geçmişi](docs/development/TASKS.md) - Proje geçmişi
- 🚀 **[VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** - Extension'ı indir

### AI Assistant Documentation

- 🤖 [CLAUDE.md](CLAUDE.md) - Turkish AI assistant guidance
- 🌍 [CLAUDE.en.md](CLAUDE.en.md) - English AI assistant guidance
- 📋 [MCP Tools Reference](docs/api/MCP_TOOLS_REFERENCE.md) - Detailed MCP tool documentation
- 📚 [Development History](docs/development/TASKS.md) - Complete project history

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

- Go 1.23+
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

Detaylı bilgi için [Geliştirici Rehberi](docs/development/contributing.md)'ne bakın.

## 📊 Proje Durumu

- **Versiyon**: v0.16.2 🚀
- **Test Coverage**: %75+ (Comprehensive test coverage with ongoing improvements)
- **Go Version**: 1.23+
- **MCP SDK**: mark3labs/mcp-go v0.6.0
- **Database**: SQLite (embedded)
- **Security**: Production-ready audit compliant
- **Thread Safety**: 100% race condition free

## 🤝 Topluluk

- 📦 [GitHub Releases](https://github.com/msenol/gorev/releases)
- 🐛 [Issue Tracker](https://github.com/msenol/gorev/issues)
- 💬 [Discussions](https://github.com/msenol/gorev/discussions)
- 📖 [Wiki](https://github.com/msenol/gorev/wiki)

## ⚠️ BREAKING CHANGE: Template Zorunluluğu

### v0.10.0'dan İtibaren Template Kullanımı Zorunludur

`gorev_olustur` tool artık kullanılamaz. Tüm görevler template kullanılarak oluşturulmalıdır.

#### 🔄 Eski Kullanım (Artık Çalışmaz)

```bash
gorev_olustur baslik="Bug fix" aciklama="..." oncelik="yuksek"
```

#### ✅ Yeni Kullanım (Zorunlu)

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

#### 🆕 Yeni Template'ler

- **bug_report_v2** - Gelişmiş bug raporu (severity, steps, environment)
- **spike_research** - Time-boxed araştırma görevleri
- **performance_issue** - Performans sorunları (metrics, targets)
- **security_fix** - Güvenlik düzeltmeleri (CVSS, components)
- **refactoring** - Kod iyileştirme (code smell, strategy)

#### 🎯 Neden Template Zorunlu

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
