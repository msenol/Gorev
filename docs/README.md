# Gorev Dokümantasyonu

Gorev MCP sunucusu için kapsamlı dokümantasyon.

## 📚 İçindekiler

### Başlangıç
- **[Kurulum](kurulum.md)** - Gorev'i sisteminize kurma rehberi
- **[Kullanım](kullanim.md)** - Temel kullanım ve komutlar
- **[Örnekler](ornekler.md)** - Pratik kullanım örnekleri

### Referans
- **[MCP Araçları](mcp-araclari.md)** - Tüm MCP tool'larının detaylı açıklaması
- **[API Referansı](api-referans.md)** - Go API dokümantasyonu
- **[Mimari](mimari.md)** - Sistem mimarisi ve tasarım kararları

### Geliştirme
- **[Geliştirici Rehberi](gelistirme.md)** - Katkıda bulunma ve geliştirme

## 🚀 Hızlı Başlangıç

```bash
# Binary indirme
curl -L https://github.com/yourusername/gorev/releases/latest/download/gorev-linux-amd64 -o gorev
chmod +x gorev

# Claude Desktop'a ekleme
# claude_desktop_config.json dosyanıza:
{
  "mcpServers": {
    "gorev": {
      "command": "/path/to/gorev",
      "args": ["serve"]
    }
  }
}

# Çalıştırma
./gorev serve
```

## 📋 Özellikler

- **Görev Yönetimi**: Görev oluşturma, listeleme, durum güncelleme
- **Proje Organizasyonu**: Görevleri projeler altında gruplama
- **MCP Entegrasyonu**: Claude ile sorunsuz çalışma
- **Hafif ve Hızlı**: Go ile yazılmış, minimal kaynak kullanımı

## 🔧 Gereksinimler

- Go 1.22+ (kaynak koddan derleme için)
- SQLite3 (runtime bağımlılığı)
- Claude Desktop veya Claude Code

## 📖 Dokümantasyon Yapısı

Bu dokümantasyon şu bölümlerden oluşur:

1. **Kullanıcı Rehberleri**: Gorev'i kullanmaya başlamak için
2. **Referans Dokümantasyonu**: Detaylı API ve tool açıklamaları
3. **Geliştirici Dokümantasyonu**: Katkıda bulunmak isteyenler için

## 💬 Destek

- [GitHub Issues](https://github.com/yourusername/gorev/issues)
- [Discussions](https://github.com/yourusername/gorev/discussions)

## 📄 Lisans

MIT License - Detaylar için [LICENSE](../LICENSE) dosyasına bakın.