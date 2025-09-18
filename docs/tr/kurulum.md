# Kurulum Kılavuzu

> **Versiyon**: Bu dokümantasyon v0.15.4+ için geçerlidir
> **Son Güncelleme**: 18 Eylül 2025

Gorev'in tüm platformlarda kurulumu için detaylı talimatlar.

## 📋 Ön Gereksinimler

- **İşletim Sistemi**: Linux, macOS, Windows
- **MCP Uyumlu Editör**: Claude Desktop, VS Code (MCP uzantısı ile), Windsurf, Cursor, Zed, veya diğer MCP destekli editörler
- **Docker** (isteğe bağlı, konteyner kullanımı için)

## 🚀 Hızlı Kurulum

### 🔥 NPX ile Kolay Kurulum (Önerilen!)

MCP istemcileri için `mcp.json` konfigürasyonunuza basitçe ekleyin:

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

**Claude Desktop için:**
```json
// Windows: %APPDATA%/Claude/claude_desktop_config.json
// macOS: ~/Library/Application Support/Claude/claude_desktop_config.json
// Linux: ~/.config/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

**VS Code için:**
```json
// .vscode/mcp.json
{
  "servers": {
    "gorev": {
      "command": "npx",
      "args": ["@gorev/mcp-server@latest"]
    }
  }
}
```

### 🖥️ Otomatik Kurulum (Geleneksel)

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | bash
```

**Belirli versiyon için:**
```bash
curl -fsSL https://raw.githubusercontent.com/msenol/Gorev/main/install.sh | VERSION=v0.15.4 bash
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

**Belirli versiyon için:**
```powershell
$env:VERSION="v0.15.4"; irm https://raw.githubusercontent.com/msenol/Gorev/main/install.ps1 | iex
```

### ✅ Kurulumu Doğrulama

```bash
gorev version
gorev help
```

## 🔧 MCP Editör Konfigürasyonu

### 🤖 Claude Desktop

Claude Desktop konfigürasyon dosyanıza şu ayarları ekleyin:

**Dosya Konumları:**
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

**NPX ile (Önerilen):**
```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

**Yerel kurulum ile:**
```json
{
  "mcpServers": {
    "gorev": {
      "command": "/usr/local/bin/gorev",
      "args": ["serve"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

### 💻 VS Code

**Seçenek 1: Gorev VS Code Uzantısı (Önerilen)**

[VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)'ten kurulum:

```bash
code --install-extension mehmetsenol.gorev-vscode
```

**Seçenek 2: Genel MCP Uzantısı**

1. [MCP için VS Code uzantısını](https://marketplace.visualstudio.com/items?itemName=mark3labs.mcp-vscode) kurun
2. `.vscode/mcp.json` dosyası oluşturun:

```json
{
  "servers": {
    "gorev": {
      "command": "npx",
      "args": ["@gorev/mcp-server@latest"]
    }
  }
}
```

### 🌊 Windsurf

Windsurf ayarlarınıza MCP sunucu konfigürasyonu ekleyin:

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

### ⚡ Cursor

Cursor'ın MCP konfigürasyonuna ekleyin:

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@gorev/mcp-server@latest"]
    }
  }
}
```

## 📱 Manuel Kurulum

### 📥 Binary İndirme

GitHub Releases sayfasından platform-specific binary'leri indirin:

```bash
# Linux AMD64
wget https://github.com/msenol/gorev/releases/download/v0.15.4/gorev-linux-amd64.tar.gz

# macOS
wget https://github.com/msenol/gorev/releases/download/v0.15.4/gorev-darwin-amd64.tar.gz

# Windows
curl -L -o gorev-windows-amd64.zip https://github.com/msenol/gorev/releases/download/v0.15.4/gorev-windows-amd64.zip
```

### 🔧 Kaynak Koddan Derleme

```bash
# Repository'yi klonlayın
git clone https://github.com/msenol/gorev.git
cd gorev/gorev-mcpserver

# Bağımlılıkları indirin
go mod download

# Derleyin
go build -o gorev cmd/gorev/main.go

# Kurulum dizinine taşıyın
sudo mv gorev /usr/local/bin/
```

## 🐳 Docker Kurulumu

### 🔨 Docker Image Build

```bash
git clone https://github.com/msenol/gorev.git
cd gorev
docker build -t gorev:latest .
```

### 🚀 Docker Compose

`docker-compose.yml` dosyası oluşturun:

```yaml
version: '3.8'
services:
  gorev:
    image: gorev:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data:/data
    environment:
      - GOREV_LANG=tr
      - GOREV_DB_PATH=/data/gorev.db
    command: ["gorev", "serve", "--port", "8080"]
```

Çalıştırma:
```bash
docker-compose up -d
```

## 🔧 Konfigürasyon

### 🌍 Dil Ayarları

```bash
# Türkçe
export GOREV_LANG=tr
gorev serve --lang=tr

# İngilizce
export GOREV_LANG=en
gorev serve --lang=en
```

### 📁 Veritabanı Konumu

```bash
# Özel veritabanı konumu
export GOREV_DB_PATH=/path/to/your/gorev.db
gorev serve

# Workspace veritabanı (.gorev/gorev.db)
gorev init

# Global veritabanı (~/.gorev/gorev.db)
gorev init --global
```

### 🔧 Sunucu Ayarları

```bash
# Farklı port
gorev serve --port 9090

# Debug modu
gorev serve --debug

# Belirli host
gorev serve --host 0.0.0.0
```

## 🔍 Sorun Giderme

### ❌ Yaygın Sorunlar

**1. Permission Denied (Linux/macOS)**
```bash
sudo chmod +x /usr/local/bin/gorev
```

**2. Command Not Found**
```bash
# PATH'e ekleyin
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

**3. MCP Bağlantı Sorunu**
```bash
# Sunucunun çalışıp çalışmadığını kontrol edin
gorev serve --debug

# Port kullanımda kontrolü
netstat -tlnp | grep 8080
```

**4. VS Code Uzantısı Çalışmıyor**
- VS Code'u yeniden başlatın
- MCP sunucusunun çalıştığından emin olun
- Extension Host'u yeniden yükleyin (Ctrl+Shift+P → "Developer: Reload Window")

### 📋 Sistem Gereksinimleri

- **RAM**: Minimum 512MB, önerilen 1GB+
- **Disk**: 100MB boş alan
- **CPU**: x86_64 veya ARM64
- **Network**: GitHub'a erişim (kurulum için)

### 🛠️ Debug Bilgileri

```bash
# Sistem bilgileri
gorev version
gorev serve --debug

# Veritabanı durumu
ls -la ~/.gorev/

# Loglar
tail -f ~/.gorev/gorev.log
```

## 🔗 Faydalı Linkler

- **[Ana Proje](https://github.com/msenol/gorev)** - Kaynak kod ve issue'lar
- **[VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)** - Uzantı indirme
- **[Kullanım Kılavuzu](kullanim.md)** - Detaylı kullanım örnekleri
- **[MCP Araçları](mcp-araclari.md)** - 48 MCP tool referansı
- **[GitHub Discussions](https://github.com/msenol/gorev/discussions)** - Topluluk desteği

---

*Bu dokümantasyon Claude (Anthropic) ile birlikte hazırlanmıştır*