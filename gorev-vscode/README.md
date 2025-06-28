# 🎨 Gorev VS Code Extension

Gorev MCP server için zengin görsel arayüz sunan VS Code extension'ı. TreeView panelleri, komut paleti entegrasyonu ve status bar desteği ile görev yönetimini kolaylaştırır.

![VS Code Marketplace Version](https://img.shields.io/visual-studio-marketplace/v/gorev.gorev-vscode?style=flat-square)
![VS Code Marketplace Downloads](https://img.shields.io/visual-studio-marketplace/d/gorev.gorev-vscode?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

## ✨ Özellikler

### 📊 TreeView Panelleri
- **Görevler**: Durum bazlı gruplandırma, öncelik renklendirmesi
- **Projeler**: Aktif proje vurgulama, görev sayıları
- **Şablonlar**: Kategori bazlı listeleme, hızlı görev oluşturma

### ⌨️ Komut Paleti
- `Gorev: Create Task` - Yeni görev oluştur
- `Gorev: Quick Create Task` (`Ctrl+Shift+G`) - Hızlı görev oluşturma
- `Gorev: Create Project` - Yeni proje oluştur
- `Gorev: Show Summary` - Özet istatistikleri göster
- `Gorev: Connect` - MCP server'a bağlan
- `Gorev: Disconnect` - Bağlantıyı kes

### 🎯 Context Menu İşlemleri
- Görev durumu güncelleme
- Görev silme (onaylı)
- Proje aktif yapma
- Detay görüntüleme

### 📈 Status Bar
- Bağlantı durumu göstergesi
- Toplam/tamamlanan görev sayısı
- Aktif proje bilgisi

### 🎨 Tema Desteği
- Öncelik bazlı renklendirme
- Dark/Light tema uyumu
- Özelleştirilebilir renkler

## 📦 Kurulum

### Marketplace'den (Yakında)
```
VS Code Extensions → "Gorev Task Orchestrator" ara → Install
```

### Local Kurulum
```bash
# Repository'yi klonla
git clone https://github.com/yourusername/gorev.git
cd gorev/gorev-vscode

# Bağımlılıkları yükle
npm install

# Extension'ı derle
npm run compile

# VS Code'da test et
# F5 tuşuna bas veya Run → Start Debugging
```

## ⚙️ Konfigürasyon

VS Code ayarlarında (`settings.json`):

```json
{
  // MCP server binary yolu
  "gorev.serverPath": "/path/to/gorev-mcpserver/gorev",
  
  // Otomatik bağlanma (varsayılan: true)
  "gorev.autoConnect": true,
  
  // Status bar gösterimi (varsayılan: true)
  "gorev.showStatusBar": true,
  
  // Otomatik yenileme süresi (saniye, 0 = devre dışı)
  "gorev.refreshInterval": 30,
  
  // Debug loglama (varsayılan: false)
  "gorev.debug": false
}
```

## 🚀 Kullanım

### İlk Kurulum
1. Gorev MCP server'ı yükleyin ([kurulum rehberi](../docs/kurulum.md))
2. Extension'ı yükleyin
3. `gorev.serverPath` ayarını yapın
4. VS Code'u yeniden başlatın

### Temel Kullanım
1. Activity Bar'da Gorev ikonuna tıklayın
2. TreeView'lardan görev/proje yönetin
3. `Ctrl+Shift+G` ile hızlı görev oluşturun
4. Sağ tık menüleri ile işlem yapın

### İpuçları
- 🔄 TreeView'ları yenilemek için başlıktaki refresh ikonuna tıklayın
- 📌 Aktif projeyi belirleyerek yeni görevlerin otomatik atanmasını sağlayın
- 🏷️ Şablonları kullanarak tutarlı görevler oluşturun
- 📊 Status bar'a tıklayarak özet istatistikleri görün

## 🛠️ Geliştirme

### Gereksinimler
- Node.js 16+
- npm veya yarn
- VS Code 1.95.0+

### Proje Yapısı
```
gorev-vscode/
├── src/
│   ├── extension.ts          # Ana giriş noktası
│   ├── mcp/                  # MCP client
│   ├── commands/             # Komut implementasyonları
│   ├── providers/            # TreeView provider'ları
│   ├── models/               # Data modelleri
│   └── ui/                   # UI bileşenleri
├── media/                    # İkonlar ve görseller
├── package.json              # Extension manifest
└── tsconfig.json            # TypeScript konfigürasyonu
```

### Komutlar
```bash
# Geliştirme
npm run watch            # Watch mode
npm run compile         # TypeScript derleme

# Test
npm test                # Unit testler
npm run test:e2e       # E2E testler

# Paketleme
npm run package        # VSIX paketi oluştur
npm run publish       # Marketplace'e yayınla
```

### Debug
1. VS Code'da projeyi aç
2. `F5` tuşuna bas veya Debug panelinden "Run Extension" seç
3. Yeni VS Code penceresi açılacak (Extension Development Host)
4. Output panelinde "Gorev" kanalını kontrol et

## 🐛 Bilinen Sorunlar

### Markdown Parser
TreeView'ların düzgün çalışması için MCP response'larının parse edilmesi gerekiyor. Geçici çözüm:
- Server response'ları düz metin olarak işleniyor
- Markdown formatı tam desteklenmiyor

### Icon Eksiklikleri
Extension ve TreeView ikonları henüz eklenmedi. Varsayılan VS Code ikonları kullanılıyor.

## 🤝 Katkıda Bulunma

1. Fork yapın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'feat: add amazing feature'`)
4. Branch'i push edin (`git push origin feature/amazing-feature`)
5. Pull Request açın

### Kod Standartları
- TypeScript strict mode
- ESLint kurallarına uyum
- Prettier formatlaması
- Conventional commits

## 📝 Lisans

MIT License - detaylar için [LICENSE](../LICENSE) dosyasına bakın.

## 🔗 Linkler

- [Ana Proje](https://github.com/yourusername/gorev)
- [MCP Server Dokümantasyonu](../gorev-mcpserver/README.md)
- [API Referansı](docs/api/README.md)
- [Sorun Bildirme](https://github.com/yourusername/gorev/issues)

---

<div align="center">
💡 Bu extension, Gorev MCP server'ın görsel arayüzüdür. MCP protokolü sayesinde server'a diğer editörlerden de bağlanabilirsiniz.
</div>