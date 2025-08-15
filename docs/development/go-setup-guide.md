# Go Development Environment Setup Guide

Bu doküman Gorev projesi için Go development environment'ın kurulumunu açıklar.

## 🎯 Önkoşullar

Gorev projesi aşağıdaki gereksinimleri karşılamalıdır:

- **Go Version**: 1.21 veya üzeri
- **Platform**: Linux, macOS, Windows
- **WSL**: Windows kullanıcıları için WSL2 desteklenir

## 📥 Go Kurulumu

### Option 1: Official Binary (Recommended)

```bash
# Latest Go version'u indirin (Linux/WSL)
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz

# Önceki kurulumu temizleyin
sudo rm -rf /usr/local/go

# Yeni version'u yükleyin
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# PATH'e ekleyin (bash/zsh)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GOBIN=$GOPATH/bin' >> ~/.bashrc

# Shell'i yeniden başlatın
source ~/.bashrc
```

### Option 2: Package Manager

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# macOS (Homebrew)
brew install go

# Windows (Chocolatey)
choco install golang
```

### Option 3: Version Manager (Advanced)

```bash
# g (Go version manager)
curl -sSL https://git.io/g-install | sh -s
g install latest
```

## ✅ Kurulum Doğrulama

```bash
# Go version kontrol
go version
# Expected: go version go1.21.5 linux/amd64

# GOPATH kontrol
go env GOPATH
# Expected: /home/[username]/go

# GOROOT kontrol  
go env GOROOT
# Expected: /usr/local/go

# Workspace test
go env
```

## 🔧 Gorev Projesi için Gerekli Tools

```bash
# Code formatting
go install golang.org/x/tools/cmd/goimports@latest

# Linting (optional but recommended)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Test coverage
go install golang.org/x/tools/cmd/cover@latest

# Database migrations (used in project)
go install -tags 'postgres sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## 🚀 Proje Build ve Test

```bash
# Ana dizinden
cd /path/to/Gorev

# Dependencies'leri indir
make deps

# Build
make build

# Test
make test

# Full workflow
make all
```

## 🛠️ Development Workflow

### Daily Development Commands

```bash
# Server geliştirme
make server-run          # Server'ı debug mode'da çalıştır
make server-test         # Unit testleri çalıştır
make server-coverage     # Coverage raporunu oluştur

# Extension geliştirme  
make extension-build     # TypeScript compile
make extension-test      # Extension testleri

# Quality checks
make fmt                 # Code formatting
make lint                # Linting
make pre-commit          # Commit öncesi tüm kontroller
```

### VS Code Integration

`.vscode/settings.json`:

```json
{
    "go.gopath": "${env:GOPATH}",
    "go.goroot": "${env:GOROOT}",
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v", "-race"],
    "go.coverOnSave": true,
    "go.coverageDisplayFormat": "highlight"
}
```

## 🐛 Troubleshooting

### Common Issues

#### 1. "go: command not found"

```bash
# PATH ayarlarını kontrol edin
echo $PATH | grep go

# Shell'i yeniden başlatın
source ~/.bashrc

# Go installation'ı doğrulayın
ls -la /usr/local/go/bin/go
```

#### 2. "permission denied"

```bash
# GOPATH permissions'ı düzeltin
mkdir -p $GOPATH/bin
chmod 755 $GOPATH
chmod 755 $GOPATH/bin
```

#### 3. Module issues

```bash
# Module cache'i temizleyin
go clean -modcache

# Dependencies'leri yeniden indirin
cd gorev-mcpserver
go mod download
go mod tidy
```

#### 4. WSL-specific issues

```bash
# Windows PATH pollution'ını önleyin
echo 'export PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"' >> ~/.bashrc

# Go binary'yi WSL native olarak kurun
sudo apt install golang-go  # Simple option
# OR official binary'yi Linux section'dan takip edin
```

## 📊 Performance Tips

```bash
# Build cache'i enable
export GOCACHE=$HOME/.cache/go-build

# Module proxy (faster downloads)
export GOPROXY=https://proxy.golang.org,direct

# Disable CGO if not needed
export CGO_ENABLED=0

# Parallel compilation
export GOMAXPROCS=$(nproc)
```

## 🔄 Environment Variables

Proje için önerilen `.bashrc` / `.zshrc` ayarları:

```bash
# Go Environment
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT/bin:$GOBIN

# Gorev Project
export GOREV_LANG=tr                    # Default language
export GOREV_DEBUG=false               # Debug mode
export GOREV_DB_PATH="./gorev.db"      # Database path

# Performance
export GOCACHE=$HOME/.cache/go-build
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
```

## 📚 Useful Resources

- [Official Go Documentation](https://golang.org/doc/)
- [Go Modules Reference](https://golang.org/ref/mod)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

Bu guide'ı takip ettikten sonra `make all` komutu başarıyla çalışmalıdır. Sorun yaşarsanız GitHub Issues'da bildirebilirsiniz.
