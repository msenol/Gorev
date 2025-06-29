# Getting Started with Gorev VS Code Extension

Bu rehber, Gorev VS Code Extension'ı kullanmaya başlamanız için gereken tüm adımları içerir.

## İçindekiler

- [Kurulum](#kurulum)
- [İlk Yapılandırma](#ilk-yapılandırma)
- [Temel Kullanım](#temel-kullanım)
- [Gelişmiş Özellikler](#gelişmiş-özellikler)
- [Sorun Giderme](#sorun-giderme)

## Kurulum

### Ön Gereksinimler

1. **VS Code**: Version 1.95.0 veya üzeri
2. **Gorev MCP Server**: [Kurulum rehberi](../../../docs/kurulum.md)
3. **Node.js**: Extension development için (opsiyonel)

### Extension Kurulum Seçenekleri

#### 1. Marketplace'den (Yakında)
```
1. VS Code'u açın
2. Extensions paneline gidin (Ctrl+Shift+X)
3. "Gorev Task Orchestrator" arayın
4. Install butonuna tıklayın
```

#### 2. VSIX Dosyasından
```bash
# VSIX dosyasını indirin
curl -L https://github.com/msenol/gorev/releases/latest/download/gorev-vscode-0.1.0.vsix -o gorev.vsix

# VS Code'da yükleyin
code --install-extension gorev.vsix
```

#### 3. Kaynak Koddan
```bash
# Repository'yi klonlayın
git clone https://github.com/msenol/gorev.git
cd gorev/gorev-vscode

# Bağımlılıkları yükleyin
npm install

# Extension'ı derleyin
npm run compile

# VS Code'da test edin (F5)
```

## İlk Yapılandırma

### 1. Server Path Ayarlama

Extension'ın MCP server'a bağlanabilmesi için server binary yolunu belirtmeniz gerekir.

**Ayarları açın**: `Ctrl+,` veya File → Preferences → Settings

**Gorev ayarlarını bulun**: Arama kutusuna "gorev" yazın

**Server Path'i ayarlayın**:
- Windows: `C:\Program Files\gorev\gorev.exe`
- macOS/Linux: `/usr/local/bin/gorev`

Alternatif olarak `settings.json` dosyasını düzenleyin:

```json
{
  "gorev.serverPath": "/usr/local/bin/gorev",
  "gorev.autoConnect": true,
  "gorev.showStatusBar": true
}
```

### 2. İlk Bağlantı

1. VS Code'u yeniden başlatın
2. Activity Bar'da Gorev ikonunu göreceksiniz
3. İkona tıklayarak Gorev panelini açın
4. Otomatik bağlanma kapalıysa, Command Palette'den `Gorev: Connect` komutunu çalıştırın

### 3. Bağlantıyı Doğrulama

- Status bar'da 🟢 simgesi: Bağlantı başarılı
- Status bar'da 🔴 simgesi: Bağlantı yok
- Output panelinde (View → Output → Gorev) detaylı logları kontrol edin

## Temel Kullanım

### Görev Oluşturma

#### Yöntem 1: TreeView Üzerinden
1. Gorev panelinde "+" ikonuna tıklayın
2. Form alanlarını doldurun:
   - Başlık (zorunlu)
   - Açıklama (opsiyonel, markdown destekli)
   - Öncelik (Düşük/Orta/Yüksek)
   - Proje (aktif proje otomatik seçili)
3. "Create" butonuna tıklayın

#### Yöntem 2: Komut Paleti
1. `Ctrl+Shift+P` ile Command Palette'i açın
2. "Gorev: Create Task" yazın ve Enter
3. Form alanlarını doldurun

#### Yöntem 3: Hızlı Oluşturma
1. `Ctrl+Shift+G` kısayolunu kullanın
2. Sadece başlık girin
3. Görev varsayılan değerlerle oluşturulur

### Görev Yönetimi

#### Durum Güncelleme
1. Görev üzerine sağ tıklayın
2. "Update Status" seçin
3. Yeni durumu seçin:
   - Beklemede
   - Devam Ediyor
   - Tamamlandı

#### Görev Detayları
- Görev üzerine tıklayarak detayları görüntüleyin
- Detay görünümünde:
  - Tam açıklama
  - Son tarih
  - Etiketler
  - Bağımlılıklar

#### Görev Silme
1. Görev üzerine sağ tıklayın
2. "Delete Task" seçin
3. Onay dialogunda "Yes" tıklayın

### Proje Yönetimi

#### Yeni Proje Oluşturma
1. Projects panelinde "+" ikonuna tıklayın
2. Proje adı ve açıklama girin
3. "Create" tıklayın

#### Aktif Proje Belirleme
1. Proje üzerine sağ tıklayın
2. "Set as Active" seçin
3. Yeni görevler otomatik olarak bu projeye atanır

### Şablon Kullanımı

1. Templates panelini açın
2. Kullanmak istediğiniz şablonu bulun
3. Şablon üzerine tıklayın
4. Gerekli alanları doldurun
5. "Create from Template" tıklayın

## Gelişmiş Özellikler

### Filtreleme ve Sıralama

#### Durum Bazlı Filtreleme
TreeView otomatik olarak görevleri duruma göre gruplar:
- 📋 Beklemede
- 🔄 Devam Ediyor
- ✅ Tamamlandı

#### Öncelik Renklendirmesi
- 🔴 Yüksek öncelik (kırmızı)
- 🟡 Orta öncelik (sarı)
- 🟢 Düşük öncelik (yeşil)

### Kısayollar

| Kısayol | Açıklama |
|---------|----------|
| `Ctrl+Shift+G` | Hızlı görev oluştur |
| `F5` | Listeleri yenile |
| `Delete` | Seçili görevi sil |

### Status Bar Özellikleri

Status bar'a tıklayarak:
- Toplam görev sayısı
- Tamamlanan görev sayısı
- Duruma göre dağılım
- Proje istatistikleri

### Otomatik Yenileme

`gorev.refreshInterval` ayarı ile otomatik yenileme süresini belirleyin:

```json
{
  "gorev.refreshInterval": 30  // 30 saniyede bir yenile
}
```

## Sorun Giderme

### Bağlantı Sorunları

**Problem**: Extension server'a bağlanamıyor

**Çözümler**:
1. Server'ın çalıştığını kontrol edin: `gorev serve`
2. Server path'inin doğru olduğunu kontrol edin
3. Windows'ta tam path kullanın: `C:\\Program Files\\gorev\\gorev.exe`
4. Firewall/antivirus ayarlarını kontrol edin

### TreeView Boş Görünüyor

**Problem**: Görevler listesi boş

**Çözümler**:
1. Refresh butonuna tıklayın
2. Output panelinde hataları kontrol edin
3. Server'ın doğru veritabanına bağlandığını kontrol edin

### Performance Sorunları

**Problem**: Extension yavaş çalışıyor

**Çözümler**:
1. `gorev.refreshInterval` değerini artırın
2. Debug mode'u kapatın: `"gorev.debug": false`
3. Çok sayıda görev varsa sayfalama özelliğini bekleyin

### Debug Logları

Detaylı hata ayıklama için:

1. Settings'de debug'ı açın:
```json
{
  "gorev.debug": true
}
```

2. Output panelini kontrol edin:
   - View → Output
   - Dropdown'dan "Gorev" seçin

3. Developer Tools'u açın:
   - Help → Toggle Developer Tools
   - Console sekmesini kontrol edin

## Sonraki Adımlar

- [Komut Referansı](../api/commands.md) - Tüm komutların detaylı açıklaması
- [API Dokümantasyonu](../api/README.md) - Extension API'si
- [Ana Dokümantasyon](../../../README.md) - Proje genel bilgileri

## Yardım ve Destek

- **GitHub Issues**: https://github.com/msenol/gorev/issues
- **Discussions**: https://github.com/msenol/gorev/discussions
- **Discord**: (Yakında)

---

<div align="center">
💡 İpucu: Extension ile ilgili önerilerinizi GitHub Issues'da paylaşın!
</div>