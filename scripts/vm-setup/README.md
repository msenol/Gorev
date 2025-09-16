# Gorev VirtualBox Linux VM Setup Guide

Bu klasör, Gorev projesini VirtualBox Linux VM'de test etmek için kapsamlı kurulum script'lerini içerir.

## 📁 Script'ler

### 1. `01-install-prerequisites.sh`
**Temel gereksinimler kurulumu**
- Go 1.23.5 kurulumu
- Node.js 20 LTS kurulumu
- VS Code kurulumu
- Build tools ve development dependencies
- Git konfigürasyonu
- Workspace dizinleri oluşturma

```bash
chmod +x 01-install-prerequisites.sh
./01-install-prerequisites.sh
```

### 2. `02-build-gorev.sh`
**Proje build ve konfigürasyon**
- GitHub'dan proje klonlama
- Go dependencies yükleme
- MCP server build etme
- Database'leri başlatma (global ve workspace)
- VS Code extension dependencies
- Development aliases oluşturma

```bash
chmod +x 02-build-gorev.sh
./02-build-gorev.sh
```

### 3. `03-setup-vscode.sh`
**VS Code extension kurulumu**
- Extension dependencies yükleme
- TypeScript compilation
- VSIX paketi oluşturma
- Extension'ı VS Code'a yükleme
- Workspace settings konfigürasyonu
- Debug ve launch configurations

```bash
chmod +x 03-setup-vscode.sh
./03-setup-vscode.sh
```

### 4. `04-run-tests.sh`
**Kapsamlı test suite**
- Prerequisites verification
- Project structure validation
- Database tests
- Unit tests (Go)
- Integration tests
- Extension tests
- Performance tests
- Test coverage raporları

```bash
chmod +x 04-run-tests.sh
./04-run-tests.sh
```

### 5. `05-debug-helper.sh`
**Debug ve troubleshooting araçları**
- Interaktif debug menüsü
- System health check
- Database inspection tools
- Server debug modes
- Extension debugging
- Log analysis
- Performance profiling
- Clean & reset options

```bash
chmod +x 05-debug-helper.sh
./05-debug-helper.sh
```

## 🚀 Hızlı Kurulum

### Tam Otomatik Kurulum
```bash
# Tüm script'leri sırayla çalıştır
cd ~/Projects/Gorev/scripts/vm-setup
./01-install-prerequisites.sh
source ~/.bashrc
./02-build-gorev.sh
./03-setup-vscode.sh
./04-run-tests.sh
```

### Manuel Adım Adım
```bash
# 1. Prerequisites
./01-install-prerequisites.sh
source ~/.bashrc

# 2. Build project
./02-build-gorev.sh

# 3. VS Code setup
./03-setup-vscode.sh

# 4. Test everything
./04-run-tests.sh

# 5. Debug if needed
./05-debug-helper.sh
```

## 📋 Sistem Gereksinimleri

### Minimum Donanım
- **RAM**: 4GB (8GB önerilen)
- **Disk**: 10GB boş alan
- **CPU**: 2 core (4 core önerilen)
- **İnternet**: GitHub'a erişim

### İşletim Sistemi
- Ubuntu 20.04+ LTS
- Debian 11+
- Linux Mint 20+
- Fedora 34+
- CentOS 8+

## 🔧 Post-Installation

### Aliases (Kurulum sonrası kullanılabilir)
```bash
# Navigation
gorev-cd           # cd ~/workspace/Gorev
gorev-server       # cd ~/workspace/Gorev/gorev-mcpserver
gorev-ext          # cd ~/workspace/Gorev/gorev-vscode

# Server commands
gorev-serve        # Start server with debug
gorev-build        # Build server
gorev-test         # Run server tests
gorev-clean        # Clean build artifacts

# Extension commands
gorev-ext-compile  # Compile extension
gorev-ext-test     # Run extension tests
gorev-ext-package  # Package VSIX

# Database commands
gorev-db-global    # Open global database
gorev-db-workspace # Open workspace database

# Development commands
gorev-logs         # Find log files
gorev-status       # Git status
gorev-pull         # Git pull
```

### Development Workflow
```bash
# 1. Start server
gorev-serve

# 2. Open VS Code for extension development
gorev-ext
code .
# Press F5 to start debugging

# 3. Test changes
gorev-test
gorev-ext-test

# 4. Package extension
gorev-ext-package
```

## 🧪 Test Senaryoları

### 1. Server Testing
```bash
cd ~/workspace/Gorev/gorev-mcpserver

# Basic commands
./gorev version
./gorev help
./gorev template aliases

# Start debug server
./gorev serve --debug
```

### 2. Extension Testing
```bash
# VS Code'da extension debug
cd ~/workspace/Gorev/gorev-vscode
code .
# F5 tuşuna bas

# Command Palette'te test:
> Gorev Debug: Seed Test Data
> Gorev Debug: Test MCP Connection
```

### 3. Database Testing
```bash
# Global database
sqlite3 ~/.gorev/gorev.db
.tables
SELECT * FROM gorevler LIMIT 5;

# Workspace database
sqlite3 ~/workspace/Gorev/.gorev/gorev.db
.tables
```

## 🐛 Troubleshooting

### Yaygın Problemler

**1. Go not found**
```bash
source ~/.bashrc
go version
```

**2. Server build fails**
```bash
cd ~/workspace/Gorev/gorev-mcpserver
make clean
make deps
make build
```

**3. Extension tests fail**
```bash
cd ~/workspace/Gorev/gorev-vscode
rm -rf node_modules
npm install
npm run compile
```

**4. Database missing**
```bash
cd ~/workspace/Gorev/gorev-mcpserver
./gorev init --global
./gorev init
```

### Debug Helper Menu
```bash
./05-debug-helper.sh
# Menüden seçin:
# 1. System Information
# 2. Project Health Check
# 3. Database Inspection
# 10. Generate Debug Report
```

## 📊 Test Coverage

Script'ler şu alanları test eder:
- ✅ Prerequisites installation
- ✅ Go module compilation
- ✅ Database initialization
- ✅ VS Code extension compilation
- ✅ Unit tests (90%+ coverage)
- ✅ Integration tests
- ✅ Static analysis
- ✅ Performance tests
- ✅ Memory leak detection

## 📁 Kurulum Sonrası Dizin Yapısı

```
~/workspace/
└── Gorev/                          # Ana proje
    ├── gorev-mcpserver/            # Go server
    │   ├── gorev                   # Server binary
    │   ├── go.mod                  # Go modules
    │   └── coverage.html           # Test coverage
    ├── gorev-vscode/               # VS Code extension
    │   ├── out/                    # Compiled JS
    │   ├── *.vsix                  # Package files
    │   └── node_modules/           # Dependencies
    ├── .gorev/                     # Workspace database
    │   └── gorev.db
    ├── .vscode/                    # VS Code config
    │   ├── settings.json
    │   ├── launch.json
    │   └── tasks.json
    ├── debug-logs/                 # Debug output
    ├── test-results/               # Test reports
    └── scripts/vm-setup/           # Setup scripts

~/.gorev/                           # Global config
├── gorev.db                        # Global database
├── migrations/                     # Database migrations
└── locales/                        # Language files
```

## 🔗 Yararlı Komutlar

### Database Queries
```sql
-- Global database
sqlite3 ~/.gorev/gorev.db

-- Show all tables
.tables

-- Show projects
SELECT id, ad, aciklama FROM projeler;

-- Show tasks
SELECT id, baslik, durum, oncelik FROM gorevler LIMIT 10;

-- Show templates
SELECT id, ad, aciklama FROM gorev_templateleri;
```

### Server Commands
```bash
# Version check
./gorev version

# Help
./gorev help

# Template list
./gorev template aliases

# Start server
./gorev serve --debug --lang=tr

# Start with English
./gorev serve --debug --lang=en
```

### Extension Commands
```bash
# Compile
npm run compile

# Watch mode
npm run watch

# Test
npm test

# Package
npm run package

# Install development version
code --install-extension *.vsix
```

## 🎯 Sonraki Adımlar

1. **Manual Testing**: Server'ı başlatın ve komutları test edin
2. **VS Code Integration**: Extension'ı VS Code'da test edin
3. **MCP Protocol**: Claude Desktop ile bağlantı test edin
4. **Performance**: Büyük veri setleri ile test edin
5. **Development**: Kendi özelliklerinizi geliştirin

---

Bu kurulum script'leri ile VirtualBox Linux VM'nizde Gorev'i tam kapsamlı test edebilirsiniz. Her script modüler olarak tasarlanmış ve bağımsız çalışabilir.