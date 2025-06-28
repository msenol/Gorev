# Gorev VS Code Extension Debug Talimatları

## Sorun: Görev ve Proje Listeleri Görünmüyor

### Adım 1: Veritabanını Temizle
```bash
cd gorev-mcpserver
rm -f gorev.db
```

### Adım 2: MCP Server'ı Başlat
```bash
cd gorev-mcpserver
./gorev serve
```
Server'ı açık bırakın.

### Adım 3: VS Code Extension'ı Debug Et
1. VS Code'da `gorev-vscode` klasörünü açın
2. F5 tuşuna basarak Extension Development Host'u başlatın
3. Extension başladığında otomatik olarak MCP server'a bağlanmaya çalışacak

### Adım 4: Test Verileri Oluştur
Extension'da:
1. Command Palette (Ctrl+Shift+P) açın
2. "Gorev Debug: Seed Test Data" komutunu çalıştırın
3. Test verileri oluşturulacak

### Adım 5: Kontrol Et
1. Sol activity bar'da Gorev ikonuna tıklayın
2. Tasks ve Projects panellerini kontrol edin
3. Eğer listeler boşsa, refresh butonuna tıklayın

## Debug Loglarını İnceleme

1. **Extension Host Output**:
   - View > Output
   - Dropdown'dan "Extension Host" seçin
   - Gorev ile ilgili logları arayın

2. **Developer Tools**:
   - Help > Toggle Developer Tools
   - Console sekmesinde hata mesajlarını kontrol edin

3. **MCP Server Logları**:
   - MCP server terminalinde hata mesajlarını kontrol edin

## Olası Sorunlar ve Çözümleri

### 1. "MCP sunucusuna bağlı değil" Mesajı
- Settings.json'da `gorev.serverPath` doğru mu kontrol edin
- MCP server çalışıyor mu kontrol edin
- Windows'ta WSL path dönüşümü sorunu olabilir

### 2. Parser Hataları
- MCP server'dan dönen markdown formatı değişmiş olabilir
- Output panelinde parse hatalarını kontrol edin
- `src/utils/markdownParser.ts` dosyasında regex pattern'leri güncelleyin

### 3. TreeView Boş Görünüyor
- MCP bağlantısı var mı kontrol edin
- `gorev_listele` ve `proje_listele` komutlarının çalıştığını doğrulayın
- Parser'ın doğru çalıştığını test edin

## Test Parser'ı Çalıştırma

```bash
cd gorev-vscode
npx ts-node -e "
import { MarkdownParser } from './src/utils/markdownParser';

// Test markdown
const testMd = \`## Test Görevi
- **ID:** 123
- **Durum:** beklemede
- **Öncelik:** yuksek\`;

const tasks = MarkdownParser.parseGorevListesi(testMd);
console.log('Parsed tasks:', tasks);
"
```