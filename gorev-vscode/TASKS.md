# Gorev VS Code Extension - Kalan Görevler

## 🎯 Kritik Özellikler

### 1. **Markdown Parser Implementation** (Yüksek Öncelik)
TreeView'ların çalışması için MCP response'larını parse etmek gerekiyor.

**Dosyalar:**
- `src/providers/gorevTreeProvider.ts` - `parseTasksFromContent()` metodu
- `src/providers/projeTreeProvider.ts` - `parseProjectsFromContent()` metodu
- `src/providers/templateTreeProvider.ts` - `parseTemplatesFromContent()` metodu

**Yapılacaklar:**
- MCP server'dan gelen markdown formatındaki response'ları parse et
- Görev, proje ve template objelerine dönüştür
- ID, başlık, durum, öncelik vb. alanları çıkar

### 2. **Icon Tasarımı** (Orta Öncelik)
Extension için SVG icon'lar oluştur.

**Gerekli Icon'lar:**
- `media/icon.svg` - Ana extension icon'u (Activity Bar için)
- `media/icons/task-pending.svg` - Bekleyen görev
- `media/icons/task-progress.svg` - Devam eden görev
- `media/icons/task-done.svg` - Tamamlanan görev
- `media/icons/priority-high.svg` - Yüksek öncelik
- `media/icons/priority-medium.svg` - Orta öncelik
- `media/icons/priority-low.svg` - Düşük öncelik

## 🚀 Gelişmiş Özellikler

### 3. **WebView Task Editor** (Orta Öncelik)
Detaylı görev düzenleme için zengin UI.

**Özellikler:**
- Markdown editor (açıklama için)
- Date picker (son tarih için)
- Tag input (auto-complete ile)
- Dependency graph visualization
- Real-time preview

**Dosya:** `src/providers/webviewProvider.ts` (yeni oluşturulacak)

### 4. **Context Menu Implementation** (Orta Öncelik)
TreeView item'ları için sağ tık menüleri.

**Menü Öğeleri:**
- Görev: Düzenle, Sil, Kopyala, Durumu Değiştir
- Proje: Aktif Yap, Düzenle, Sil
- Template: Kullan, Düzenle

### 5. **Dependency Visualization** (Orta Öncelik)
Görev bağımlılıklarını görselleştirme.

**Özellikler:**
- Mermaid.js veya D3.js ile graph gösterimi
- Interactive düzenleme
- Circular dependency uyarıları

### 6. **Template System UI** (Orta Öncelik)
Template'den görev oluşturma wizard'ı.

**Özellikler:**
- Template seçim dialog'u
- Dynamic form generation
- Field validation
- Preview before creation

### 7. **Advanced Filtering** (Orta Öncelik)
Görev filtreleme ve arama UI'ı.

**Özellikler:**
- Multi-criteria filtering
- Saved filter presets
- Quick filter buttons (Urgent, Overdue, etc.)
- Search highlighting

### 8. **Due Date Features** (Orta Öncelik)
Son tarih yönetimi geliştirmeleri.

**Özellikler:**
- Calendar widget
- Overdue highlighting
- Reminder notifications
- Bulk date operations

### 9. **Tag Management** (Orta Öncelik)
Etiket sistemi UI'ı.

**Özellikler:**
- Tag auto-complete
- Popular tags suggestion
- Tag filtering
- Tag colors

### 10. **Turkish Localization** (Orta Öncelik)
Türkçe dil desteği.

**Dosyalar:**
- `localization/package.nls.tr.json` - Extension manifest çevirileri
- `localization/bundle/tr.json` - UI string çevirileri

## 🔧 Teknik İyileştirmeler

### 11. **Performance Optimization** (Orta Öncelik)
Büyük veri setleri için optimizasyon.

**İyileştirmeler:**
- Virtual scrolling for TreeViews
- Lazy loading
- Caching strategy
- Debounced refresh

### 12. **Testing Suite** (Orta Öncelik)
Kapsamlı test coverage.

**Test Türleri:**
- Unit tests (MCP client, parsers)
- Integration tests (commands, providers)
- E2E tests (user workflows)

### 13. **Error Recovery** (Düşük Öncelik)
Gelişmiş hata yönetimi.

**Özellikler:**
- Offline mode
- Retry mechanisms
- User-friendly error messages
- Recovery suggestions

### 14. **Drag & Drop** (Düşük Öncelik)
Görev ve proje yönetimi için drag & drop.

**Özellikler:**
- Task reordering
- Move tasks between projects
- Priority adjustment by dragging

### 15. **Notifications** (Düşük Öncelik)
Akıllı bildirim sistemi.

**Özellikler:**
- Due date reminders
- Task completion celebrations
- Dependency unblock notifications

## 📦 Deployment

### 16. **Extension Packaging** (Düşük Öncelik)
VS Code Marketplace için hazırlık.

**Görevler:**
- README.md yazımı
- CHANGELOG.md oluşturma
- Icon ve screenshot'lar
- VSIX packaging
- Publishing workflow

### 17. **Documentation** (Düşük Öncelik)
Kullanıcı dokümantasyonu.

**İçerik:**
- Getting started guide
- Feature documentation
- Troubleshooting
- Video tutorials

## 📊 Özet

**Toplam Kalan Görev:** 17

**Öncelik Dağılımı:**
- Yüksek: 1 (Markdown Parser)
- Orta: 11
- Düşük: 5

**Tahmini Süre:**
- Kritik özellikler: 1-2 hafta
- Tüm özellikler: 4-6 hafta

**İlk Odaklanılacaklar:**
1. Markdown Parser (TreeView'ların çalışması için kritik)
2. Icon tasarımları (profesyonel görünüm için)
3. WebView editor (zengin kullanıcı deneyimi için)