# 🧪 Gorev VS Code Extension - Debug & Test Guide

## 🚀 Hızlı Başlangıç

1. **VS Code'da Extension'ı Başlatma:**
   ```bash
   cd gorev-vscode
   code .
   # F5 tuşuna basarak Extension Development Host'u başlatın
   ```

2. **Otomatik Test Verisi:**
   - Extension ilk açıldığında görev yoksa otomatik olarak test verisi oluşturma önerir
   - "Evet, Oluştur" seçeneğine tıklayın

## 📊 Test Verileri İçeriği

### Projeler (5 adet)
- 🚀 **Yeni Web Sitesi** - Frontend geliştirme projesi
- 📱 **Mobil Uygulama** - iOS/Android uygulama
- 🔧 **Backend API** - RESTful API geliştirme
- 📊 **Veri Analitiği** - Dashboard ve raporlama
- 🔒 **Güvenlik Güncellemeleri** - Penetrasyon testi

### Görevler (23 adet)
- **Yüksek Öncelikli**: Kritik görevler, bazıları gecikmiş
- **Orta Öncelikli**: Normal geliştirme görevleri
- **Düşük Öncelikli**: Dokümantasyon ve iyileştirmeler

### Görev Durumları
- 🔵 **Beklemede**: Çoğu görev
- 🟡 **Devam Ediyor**: 5 görev aktif
- ✅ **Tamamlandı**: 4 görev tamamlanmış

### Özel Durumlar
- ⚠️ **Gecikmiş Görevler**: 3 görev (kırmızı uyarı)
- 📅 **Bugün Bitenler**: 2 görev
- 🔗 **Bağımlılıklar**: 5 görev arası bağımlılık

### Etiketler
- `urgent`, `critical` - Acil görevler
- `frontend`, `backend`, `mobile` - Teknoloji alanları
- `feature`, `bug`, `enhancement` - Görev türleri
- `security`, `performance` - Kalite metrikleri

## 🛠️ Debug Komutları

### Command Palette (Ctrl+Shift+P)
- **Gorev Debug: Seed Test Data** - Test verileri oluştur
- **Gorev Debug: Clear Test Data** - Tüm verileri temizle

### Status Bar
- Sol altta **🧪 Debug Mode** göstergesi
- Tıklayarak test verisi oluşturabilirsiniz

## 🎯 Test Senaryoları

### 1. Gruplama Testi
- Status'e göre grupla: Beklemede, Devam Ediyor, Tamamlandı
- Önceliğe göre grupla: Yüksek, Orta, Düşük
- Projeye göre grupla: Her proje ayrı grup
- Tarihe göre grupla: Gecikmiş, Bugün, Bu Hafta

### 2. Filtreleme Testi
- 🔍 "urgent" ile arama
- Yüksek öncelikli görevleri filtrele
- Gecikmiş görevleri göster
- Belirli bir projenin görevleri

### 3. Drag & Drop Testi
- Görevi "Beklemede"den "Devam Ediyor"a sürükle
- Öncelik grupları arası taşı
- Projeler arası görev taşı
- Bağımlılık oluştur (görev üzerine bırak)

### 4. Inline Edit Testi
- F2 ile başlık düzenle
- Sağ tık > Quick Status Change
- Sağ tık > Quick Priority Change
- Sağ tık > Quick Date Change

### 5. Çoklu Seçim Testi
- Ctrl+Click ile birden fazla görev seç
- Toplu durum güncelleme
- Toplu silme işlemi

## 🔍 Sorun Giderme

### Server Bağlantısı
```bash
# Server'ı manuel başlatma
cd ../gorev-mcpserver
./gorev serve --debug
```

### Extension Yenileme
- Ctrl+R: Görevleri yenile
- F1 > Developer: Reload Window

### Log Kontrolü
- Output panel > Gorev sekmesi
- Console'da hata mesajları

## 💡 İpuçları

1. **Performans Testi**: 20+ görev ile UI tepki sürelerini test edin
2. **Edge Case'ler**: Boş projeler, uzun başlıklar, çok etiket
3. **Görsel Test**: Farklı tema ve renk ayarlarında deneyin
4. **Accessibility**: Keyboard navigation ve screen reader