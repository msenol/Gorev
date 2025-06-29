# Gorev VS Code Extension - Sorun Giderme

## TreeView Refresh Sorunu

### Sorun
TreeView'da görevler görünmüyor veya refresh edilmiyor.

### Debug Adımları

1. **Output Panel'i Kontrol Edin**
   - VS Code'da `View > Output` menüsünü açın
   - Dropdown'dan "Gorev" kanalını seçin
   - Log mesajlarını kontrol edin

2. **MCP Server Bağlantısını Kontrol Edin**
   - Status bar'da MCP bağlantı durumunu kontrol edin
   - Bağlı değilse, Command Palette'den `Gorev: Connect to MCP Server` komutunu çalıştırın

3. **Detaylı Log'ları Etkinleştirin**
   - Extension zaten debug modunda çalışıyor
   - Output panel'de şu log'ları aramalısınız:
     - `[EnhancedGorevTreeProvider] Calling gorev_listele...`
     - `[EnhancedGorevTreeProvider] Raw MCP response:`
     - `[EnhancedGorevTreeProvider] Parsed tasks count:`
     - `[MarkdownParser] Raw markdown length:`
     - `[MarkdownParser] Task match found:`

4. **Response Formatını Kontrol Edin**
   MCP server'dan gelen response şu formatta olmalı:
   ```markdown
   ## Görev Listesi

   - [beklemede] Görev başlığı (orta öncelik)
     Görev açıklaması
     ID: uuid-string
   ```

5. **Yaygın Sorunlar**

   a) **Görevler parse edilemiyor**
   - Log'larda `[MarkdownParser] Task line:` satırlarını kontrol edin
   - Response formatının beklenen formatta olduğundan emin olun

   b) **MCP server bağlantısı kopuyor**
   - gorev-mcpserver binary'sinin yolunun doğru ayarlandığından emin olun
   - Settings'de `gorev.serverPath` değerini kontrol edin

   c) **TreeView boş görünüyor**
   - Aktif proje ayarlanmış olabilir, `Gorev: Remove Active Project` komutunu deneyin
   - Veya `tum_projeler: true` parametresiyle tüm görevleri listelemek için TreeView'ı refresh edin

### Developer Console'u Kullanma

Extension geliştirme modunda çalışıyorsanız:
1. `Help > Toggle Developer Tools` menüsünü açın
2. Console sekmesine gidin
3. `[MarkdownParser]` prefix'li log'ları kontrol edin

### Hata Raporlama

Sorun devam ediyorsa, lütfen şu bilgilerle bir issue açın:
- VS Code versiyonu
- Gorev extension versiyonu
- Output panel'deki tüm log'lar
- Kullandığınız işletim sistemi