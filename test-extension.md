# Gorev VS Code Extension Test Guide

## Advanced Filtering Toolbar Test

1. **Extension'ı Başlat:**
   ```bash
   cd gorev-vscode
   code .
   # F5 tuşuna basarak Extension Development Host'u başlatın
   ```

2. **Status Bar Kontrolleri:**
   - Alt status bar'da yeni butonlar görünmeli:
     - 🔍 Ara
     - 🔧 Filtrele  
     - 📑 Profiller

3. **Arama Testi:**
   - "🔍 Ara" butonuna tıklayın
   - "bug", "frontend", veya "urgent" gibi terimler arayın
   - Görevlerin filtrelendiğini kontrol edin

4. **Gelişmiş Filtre Menüsü:**
   - "🔧 Filtrele" butonuna tıklayın
   - Multi-select quick pick açılmalı
   - Aşağıdaki filtreleri test edin:
     - **Durum**: Beklemede, Devam Ediyor, Tamamlandı
     - **Öncelik**: Yüksek, Orta, Düşük
     - **Özel Filtreler**: Gecikmiş, Bugün Biten, Bu Hafta Biten
     - **Proje**: Dinamik proje listesi

5. **Filtre Profilleri:**
   - Birkaç filtre seçin
   - Save butonuna (💾) tıklayın
   - Profil adı verin (örn: "Acil İşler")
   - "📑 Profiller" butonundan kayıtlı profili yükleyin

6. **Command Palette Testleri (Ctrl+Shift+P):**
   - `Gorev: Search Tasks`
   - `Gorev: Show Filter Menu`
   - `Gorev Filter: Show Overdue Tasks`
   - `Gorev Filter: Show High Priority Tasks`
   - `Gorev Filter: Filter by Tag`

7. **Aktif Filtre Göstergesi:**
   - Filtre uygulandığında status bar'da "🔧 X filtre aktif" görünmeli
   - Tıklayarak tüm filtreleri temizleyebilirsiniz

## Sorun Giderme

Eğer komutlar bulunamazsa:
1. Extension'ı yeniden yükleyin (Ctrl+R)
2. Output panelinde Gorev loglarını kontrol edin
3. Developer Tools'da (Help > Toggle Developer Tools) hata mesajlarını kontrol edin

## Beklenen Davranış

- Filtreler anında uygulanmalı
- TreeView otomatik yenilenmeli
- Birden fazla filtre aynı anda çalışabilmeli
- Filtre profilleri workspace'e kaydedilmeli