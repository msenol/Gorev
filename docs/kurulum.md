# Kurulum Rehberi

Gorev'i sisteminize kurmanın farklı yolları.

## 📋 Gereksinimler

- **İşletim Sistemi**: Linux, macOS, Windows
- **Claude Desktop** veya **Claude Code**
- **Docker** (opsiyonel, konteyner kurulumu için)

## 🚀 Hızlı Kurulum

### Opsiyon 1: Binary İndirme (Önerilen)

#### Linux/macOS
```bash
# En son sürümü indir
curl -L https://github.com/yourusername/gorev/releases/latest/download/gorev-linux-amd64 -o gorev

# Çalıştırma izni ver
chmod +x gorev

# İsteğe bağlı: PATH'e ekle
sudo mv gorev /usr/local/bin/
```

#### Windows
```powershell
# PowerShell ile indir
Invoke-WebRequest -Uri "https://github.com/yourusername/gorev/releases/latest/download/gorev-windows-amd64.exe" -OutFile "gorev.exe"

# Veya manuel olarak GitHub releases sayfasından indirin
```

### Opsiyon 2: Docker ile Kurulum

```bash
# Docker image'ı çek
docker pull ghcr.io/yourusername/gorev:latest

# Test et
docker run --rm ghcr.io/yourusername/gorev:latest version
```

### Opsiyon 3: Kaynak Koddan Derleme

#### Gereksinimler
- Go 1.22 veya üzeri
- Git
- Make (opsiyonel)

#### Derleme Adımları
```bash
# Kodu klonla
git clone https://github.com/yourusername/gorev.git
cd gorev

# Dependencies'leri indir
go mod download

# Derle
go build -o gorev cmd/gorev/main.go

# Veya make ile
make build

# Tüm platformlar için derle
make build-all
```

## 🔧 Claude Entegrasyonu

### Claude Desktop Konfigürasyonu

1. Claude Desktop konfigürasyon dosyasını bulun:
   - **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
   - **Linux**: `~/.config/Claude/claude_desktop_config.json`

2. Dosyayı düzenleyin ve şunu ekleyin:

#### Binary Kurulum için:
```json
{
  "mcpServers": {
    "gorev": {
      "command": "/usr/local/bin/gorev",
      "args": ["serve"]
    }
  }
}
```

#### Docker Kurulum için:
```json
{
  "mcpServers": {
    "gorev": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-v", "gorev-data:/app/data",
        "ghcr.io/yourusername/gorev:latest",
        "serve"
      ]
    }
  }
}
```

### Claude Code Konfigürasyonu

Claude Code terminal'inde:

#### Binary Kurulum:
```bash
claude mcp add gorev /usr/local/bin/gorev serve
```

#### Docker Kurulum:
```bash
claude mcp add-json gorev '{
  "type": "stdio",
  "command": "docker",
  "args": [
    "run", "--rm", "-i",
    "-v", "gorev-data:/app/data",
    "ghcr.io/yourusername/gorev:latest",
    "serve"
  ]
}'
```

## ✅ Kurulumu Test Etme

### 1. Binary Test
```bash
# Versiyon kontrolü
gorev version

# Sunucuyu test modunda çalıştır
gorev serve --test
```

### 2. Claude'da Test

Claude'u yeniden başlatın ve şu komutu deneyin:
```
Gorev MCP sunucusu çalışıyor mu? Test için basit bir görev oluştur.
```

### 3. Bağlantı Sorunlarını Giderme

Debug modunu etkinleştirin:
```bash
# Environment variable ile
MCP_DEBUG=true gorev serve

# Veya flag ile
gorev serve --debug
```

## 📁 Veri Dizini

Gorev varsayılan olarak verileri şu konumlarda saklar:

- **Linux/macOS**: `~/.gorev/data/`
- **Windows**: `%USERPROFILE%\.gorev\data\`
- **Docker**: `/app/data` (volume mount edilmiş)

Özel konum belirtmek için:
```bash
gorev serve --data-dir /path/to/data
```

## 🔄 Güncelleme

### Binary Güncelleme
```bash
# Eski versiyonu yedekle
mv /usr/local/bin/gorev /usr/local/bin/gorev.backup

# Yeni versiyonu indir ve kur
curl -L https://github.com/yourusername/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev
sudo mv gorev /usr/local/bin/
```

### Docker Güncelleme
```bash
# Yeni image'ı çek
docker pull ghcr.io/yourusername/gorev:latest

# Claude Desktop'ı yeniden başlat
```

## ❌ Kaldırma

### Binary Kaldırma
```bash
# Binary'yi sil
sudo rm /usr/local/bin/gorev

# Veri dizinini sil (DİKKAT: Tüm veriler silinir!)
rm -rf ~/.gorev
```

### Docker Kaldırma
```bash
# Image'ı sil
docker rmi ghcr.io/yourusername/gorev:latest

# Volume'u sil (DİKKAT: Tüm veriler silinir!)
docker volume rm gorev-data
```

### Claude Konfigürasyonunu Temizle
Claude Desktop config dosyasından `gorev` girişini silin.

## 🆘 Yardım

Kurulum sorunları için:
- [GitHub Issues](https://github.com/yourusername/gorev/issues)
- [Troubleshooting Guide](https://github.com/yourusername/gorev/wiki/Troubleshooting)