# Kurulum Rehberi

Gorev'i sisteminize kurmanın detaylı adımları.

## 📋 Gereksinimler

- **İşletim Sistemi**: Linux, macOS, Windows
- **MCP Uyumlu Editör**: Claude Desktop, VS Code (MCP extension ile), Windsurf, Cursor, Zed veya diğer MCP destekli editörler
- **Docker** (opsiyonel, konteyner kurulumu için)

## 🚀 Platform Bazlı Kurulum

### 🪟 Windows Kurulumu

#### Yöntem 1: Binary İndirme

```powershell
# PowerShell'de (Administrator olarak çalıştırın)

# 1. Gorev dizini oluştur
New-Item -ItemType Directory -Force -Path "C:\Program Files\gorev"

# 2. Binary'yi indir
Invoke-WebRequest -Uri "https://github.com/msenol/gorev/releases/latest/download/gorev-windows-amd64.exe" -OutFile "$env:TEMP\gorev.exe"

# 3. Program Files'a taşı
Move-Item "$env:TEMP\gorev.exe" "C:\Program Files\gorev\gorev.exe" -Force

# 4. PATH'e ekle (kalıcı)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Program Files\gorev", [EnvironmentVariableTarget]::Machine)

# 5. PowerShell'i yeniden başlat ve test et
gorev version
```

#### Yöntem 2: Scoop ile Kurulum (Alternatif)

```powershell
# Scoop kurulu değilse önce kur
irm get.scoop.sh | iex

# Gorev'i kur
scoop bucket add gorev https://github.com/msenol/scoop-gorev
scoop install gorev
```

#### Windows Defender İstisnası

```powershell
# Administrator PowerShell'de
Add-MpPreference -ExclusionPath "C:\Program Files\gorev\gorev.exe"
```

#### MCP Editör Konfigürasyonu (Windows)

##### Claude Desktop

1. Konfigürasyon dosyasını açın:
   ```
   %APPDATA%\Claude\claude_desktop_config.json
   ```

2. Şu içeriği ekleyin:
   ```json
   {
     "mcpServers": {
       "gorev": {
         "command": "C:\\Program Files\\gorev\\gorev.exe",
         "args": ["serve"],
         "env": {
           "GOREV_DATA_DIR": "%USERPROFILE%\\.gorev"
         }
       }
     }
   }
   ```

### 🍎 macOS Kurulumu

#### Yöntem 1: Homebrew ile Kurulum (Önerilen)

```bash
# Homebrew tap ekle
brew tap msenol/gorev

# Gorev'i kur
brew install gorev

# Test et
gorev version
```

#### Yöntem 2: Binary İndirme

```bash
# 1. Binary'yi indir
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-darwin-amd64 -o gorev

# 2. Çalıştırma izni ver
chmod +x gorev

# 3. Güvenlik kontrolünü geç (ilk çalıştırmada)
xattr -d com.apple.quarantine gorev

# 4. /usr/local/bin'e taşı
sudo mv gorev /usr/local/bin/

# 5. Test et
gorev version
```

#### macOS Gatekeeper Uyarısı

İlk çalıştırmada "geliştirici doğrulanamadı" hatası alırsanız:

1. **System Preferences → Security & Privacy → General**
2. "gorev was blocked" mesajının yanındaki **"Allow Anyway"** butonuna tıklayın
3. Veya Terminal'de: `sudo spctl --master-disable` (güvenlik riskli)

#### MCP Editör Konfigürasyonu (macOS)

##### Claude Desktop

1. Konfigürasyon dosyasını açın:
   ```bash
   open ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

2. Şu içeriği ekleyin:
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

### 🐧 Linux Kurulumu

#### Yöntem 1: Sistem Paket Yöneticileri

**Debian/Ubuntu (APT)**
```bash
# PPA ekle
sudo add-apt-repository ppa:msenol/gorev
sudo apt update

# Kur
sudo apt install gorev
```

**Fedora/RHEL (DNF)**
```bash
# Repo ekle
sudo dnf config-manager --add-repo https://github.com/msenol/gorev/releases/latest/download/gorev.repo

# Kur
sudo dnf install gorev
```

**Arch Linux (AUR)**
```bash
# yay kullanarak
yay -S gorev-bin

# Veya manual
git clone https://aur.archlinux.org/gorev-bin.git
cd gorev-bin
makepkg -si
```

#### Yöntem 2: Binary İndirme

```bash
# Binary'yi indir
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-linux-amd64 -o gorev

# Çalıştırma izni ver
chmod +x gorev

# Sistem geneline kur
sudo mv gorev /usr/local/bin/

# Test et
gorev version
```

#### MCP Editör Konfigürasyonu (Linux)

##### Claude Desktop

1. Konfigürasyon dosyasını düzenleyin:
   ```bash
   mkdir -p ~/.config/Claude
   nano ~/.config/Claude/claude_desktop_config.json
   ```

2. Şu içeriği ekleyin:
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

## 🐳 Docker ile Kurulum

### Tüm Platformlarda Çalışır

```bash
# Docker image'ı çek
docker pull ghcr.io/msenol/gorev:latest

# Volume oluştur (veri kalıcılığı için)
docker volume create gorev-data

# Test çalıştırması
docker run --rm -v gorev-data:/data ghcr.io/msenol/gorev:latest version
```

### MCP Editör Docker Konfigürasyonu

#### Claude Desktop

```json
{
  "mcpServers": {
    "gorev": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-v", "gorev-data:/data",
        "ghcr.io/msenol/gorev:latest",
        "serve"
      ]
    }
  }
}
```

## 📁 Veri Dizini Konumları

| Platform | Varsayılan Konum | Environment Variable |
|----------|------------------|---------------------|
| Windows | `%USERPROFILE%\.gorev\` | `GOREV_DATA_DIR` |
| macOS | `~/.gorev/` | `GOREV_DATA_DIR` |
| Linux | `~/.gorev/` | `GOREV_DATA_DIR` |
| Docker | `/data` (volume) | N/A |

## 🔧 Gelişmiş Konfigürasyon

### Port Değiştirme

```json
{
  "mcpServers": {
    "gorev": {
      "command": "gorev",
      "args": ["serve", "--port", "8080"]
    }
  }
}
```

### Debug Mode

```json
{
  "mcpServers": {
    "gorev": {
      "command": "gorev",
      "args": ["serve", "--debug"],
      "env": {
        "GOREV_LOG_LEVEL": "debug"
      }
    }
  }
}
```

### Çoklu Instance

```json
{
  "mcpServers": {
    "gorev-personal": {
      "command": "gorev",
      "args": ["serve", "--port", "8081"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev-personal"
      }
    },
    "gorev-work": {
      "command": "gorev",
      "args": ["serve", "--port", "8082"],
      "env": {
        "GOREV_DATA_DIR": "~/.gorev-work"
      }
    }
  }
}
```

## 🎨 VS Code Extension Kurulumu

### Seçenek 1: Gorev VS Code Extension (Opsiyonel)

Gorev VS Code Extension, MCP server'a zengin görsel arayüz sağlar. TreeView panelleri, komut paleti ve status bar desteği sunar.

#### Marketplace'den Kurulum (Yakında)
```
1. VS Code Extensions panelini aç (Ctrl+Shift+X)
2. "Gorev Task Orchestrator" ara
3. Install butonuna tıkla
```

#### Local Development Kurulumu
```bash
# Repository'yi klonla
git clone https://github.com/yourusername/gorev.git
cd gorev/gorev-vscode

# Bağımlılıkları yükle
npm install

# Extension'ı derle
npm run compile

# VS Code'da test et
# F5 tuşuna bas veya Run > Start Debugging
```

#### Extension Konfigürasyonu

VS Code ayarlarında (`settings.json`):

```json
{
  // MCP server binary yolu
  "gorev.serverPath": "/usr/local/bin/gorev",
  
  // Windows için örnek:
  // "gorev.serverPath": "C:\\Program Files\\gorev\\gorev.exe",
  
  // Otomatik bağlanma
  "gorev.autoConnect": true,
  
  // Status bar gösterimi
  "gorev.showStatusBar": true
}
```

#### Extension Kullanımı

1. **Activity Bar**: Gorev ikonuna tıklayarak paneli aç
2. **TreeView Panelleri**: 
   - Görevler (durum bazında gruplandırılmış)
   - Projeler (aktif proje vurgulanmış)
   - Şablonlar (kategori bazında listelenmiş)
3. **Komut Paleti**: `Ctrl+Shift+P` > "Gorev" yaz
4. **Hızlı Görev**: `Ctrl+Shift+G` kısayolu
5. **Status Bar**: Bağlantı durumu ve özet bilgiler

### Seçenek 2: MCP Extension ile Kullanım

Eğer Gorev Extension kullanmak istemiyorsanız, standart MCP extension ile de kullanabilirsiniz:

```json
{
  "mcp.servers": {
    "gorev": {
      "command": "/usr/local/bin/gorev",
      "args": ["serve"]
    }
  }
}
```

> **Not**: Her iki extension'ı aynı anda kullanmayın. Ya Gorev Extension ya da MCP Extension tercih edin.

## ✅ Kurulum Doğrulama

### 1. CLI Test

```bash
# Version kontrolü
gorev version

# Komutları listele
gorev --help

# Server'ı test et
gorev serve --test
```

### 2. MCP Editör Test

MCP uyumlu editörünüzü yeniden başlatın ve AI asistanınıza test edin:
```
"Gorev çalışıyor mu? Test için yeni bir görev oluştur."
```

> **Not**: VS Code için ya Gorev Extension ya da MCP extension'ının yüklü olduğundan emin olun.

### 3. Log Kontrolü

```bash
# Windows
type %USERPROFILE%\.gorev\logs\gorev.log

# macOS/Linux
tail -f ~/.gorev/logs/gorev.log
```

## 🆘 Sorun Giderme

### Windows Sorunları

**"Windows korumalı bilgisayarınızı korudu" hatası:**
- "Daha fazla bilgi" → "Yine de çalıştır"

**DLL hatası:**
```powershell
# Visual C++ Redistributable kur
winget install Microsoft.VCRedist.2022.x64
```

**PowerShell execution policy:**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### macOS Sorunları

**"cannot be opened because the developer cannot be verified":**
```bash
# Quarantine flag'ini kaldır
xattr -cr /usr/local/bin/gorev

# Veya System Preferences'tan izin ver
```

**Port permission hatası:**
```bash
# 1024'ten büyük port kullan
gorev serve --port 8080
```

### Linux Sorunları

**Permission denied:**
```bash
# Çalıştırma izni ver
chmod +x /usr/local/bin/gorev

# SELinux context düzelt (RHEL/Fedora)
sudo restorecon -v /usr/local/bin/gorev
```

**Shared library hatası:**
```bash
# Gerekli kütüphaneleri kur
sudo apt-get install libc6  # Debian/Ubuntu
sudo dnf install glibc      # Fedora
```

### Genel Sorunlar

**MCP bağlantı hatası:**
1. Editörünüzü tamamen kapatın
2. Config dosyasını kontrol edin (JSON syntax)
3. `gorev serve` komutunu manuel test edin
4. Editörünüzü yeniden başlatın

**VS Code MCP extension sorunları:**
- Extension'ın güncel olduğundan emin olun
- Developer Console'dan hata loglarını kontrol edin
- MCP server listesini yenileyin

**Database locked hatası:**
```bash
# Eski process'leri temizle
pkill -f gorev
rm ~/.gorev/gorev.db-wal
rm ~/.gorev/gorev.db-shm
```

## 🔄 Güncelleme

### Otomatik Güncelleme

```bash
# Self-update komutu
gorev update

# Veya package manager ile
brew upgrade gorev        # macOS
scoop update gorev       # Windows
sudo apt upgrade gorev   # Debian/Ubuntu
```

### Manuel Güncelleme

```bash
# Mevcut versiyonu kontrol et
gorev version

# Yedek al
cp $(which gorev) ~/gorev.backup

# Yeni versiyonu indir ve değiştir
# (Platform spesifik komutları kullan)
```

## ❌ Kaldırma

### Windows
```powershell
# Program Files'tan sil
Remove-Item -Recurse -Force "C:\Program Files\gorev"

# Veri dizinini sil (DİKKAT: Veriler silinir!)
Remove-Item -Recurse -Force "$env:USERPROFILE\.gorev"

# PATH'den kaldır
[Environment]::SetEnvironmentVariable("Path", ($env:Path -split ';' | Where-Object { $_ -ne "C:\Program Files\gorev" }) -join ';', [EnvironmentVariableTarget]::Machine)
```

### macOS
```bash
# Homebrew ile kurulduysa
brew uninstall gorev

# Manuel kurulum
sudo rm /usr/local/bin/gorev
rm -rf ~/.gorev
```

### Linux
```bash
# Package manager ile
sudo apt remove gorev      # Debian/Ubuntu
sudo dnf remove gorev      # Fedora
yay -R gorev-bin          # Arch

# Manuel kurulum
sudo rm /usr/local/bin/gorev
rm -rf ~/.gorev
```

## 📚 İlgili Dokümantasyon

- [Kullanım Kılavuzu](kullanim.md)
- [MCP Araçları Referansı](mcp-araclari.md)
- [Örnek Kullanımlar](ornekler.md)

## 💬 Destek

Kurulum sorunları için:
- [GitHub Issues](https://github.com/msenol/gorev/issues)
- [Discussions](https://github.com/msenol/gorev/discussions)

---

<div align="center">

*🤖 Bu detaylı kurulum rehberi Claude (Anthropic) yardımıyla hazırlanmıştır - AI destekli dokümantasyonun gücü!*

</div>