# Gorev VS Code Extension

<p align="center">
  <img src="media/icon.png" alt="Gorev Logo" width="128" height="128">
</p>

> ⚠️ **BREAKING CHANGE (v0.4.0)**: Template kullanımı artık zorunludur! Doğrudan görev oluşturma kaldırıldı. Tüm görevler template kullanılarak oluşturulmalıdır. [Detaylar](#breaking-change-template-zorunluluğu)

> 🌐 **v0.5.0'da YENİ**: Tam iki dilli destek! Extension artık VS Code dil ayarınıza göre otomatik olarak Türkçe veya İngilizce görüntülenir.

Gorev için güçlü ve kullanıcı dostu VS Code extension'ı. MCP protokolü üzerinden Gorev sunucusuyla iletişim kurarak gelişmiş görev yönetimi özellikleri sunar.

## 🚀 Özellikler

### 🌍 İki Dilli Destek (v0.5.0+)
- **Otomatik Dil Algılama**: VS Code dil ayarına göre arayüz dili
- **Tam Yerelleştirme**: Tüm UI elemanları, komutlar ve mesajlar Türkçe ve İngilizce
- **Sorunsuz Deneyim**: Ayar gerektirmez - VS Code dilinizle otomatik çalışır

### Enhanced TreeView
- **Gruplama**: Durum, öncelik, proje, etiket veya son tarihe göre görevleri gruplayın
- **Çoklu Seçim**: Ctrl/Cmd+Click ile birden fazla görevi seçin
- **Sıralama**: Başlık, öncelik, son tarih ve daha fazlasına göre sıralayın
- **Renk Kodlaması**: Öncelik bazlı görsel ayırt edicilik
- **Hızlı Tamamlama**: Checkbox ile tek tıkla tamamlama
- **Badges**: Görev sayıları ve son tarih uyarıları

### Drag & Drop Desteği
- 🔄 Görevleri projeler arası taşıma
- 📊 Sürükleyerek durum değiştirme
- 🎯 Öncelik sıralaması değiştirme
- 🔗 Bağımlılık oluşturma (görev üzerine bırakarak)
- ✨ Görsel geri bildirim ve animasyonlar

### Inline Düzenleme
- ✏️ F2 veya double-click ile hızlı düzenleme
- 📝 Context menu ile durum/öncelik değiştirme
- 📅 Inline tarih seçici
- ❌ Escape ile iptal, ✅ Enter ile kaydet

### Gelişmiş Filtreleme
- 🔍 Gerçek zamanlı arama
- 🎛️ Gelişmiş filtreler (durum, öncelik, etiket, tarih)
- 💾 Kayıtlı filtre profilleri
- 📊 Status bar entegrasyonu
- ⚡ Hızlı filtre kısayolları

### Task Dependencies (v0.3.4 NEW!)
- 🏷️ **Dependency Badges**: TreeView'da görsel bağımlılık göstergeleri
  - `[🔗3]`: Bu task 3 göreve bağımlı
  - `[🔗2 ⚠️1]`: 2 bağımlılık, 1 tanesi tamamlanmamış
  - `[← 2]`: 2 task bu göreve bağımlı
- ➕ **Add Dependency**: Context menu ile kolay bağımlılık ekleme
- 📋 **Always Visible Dependencies**: Task detail'de her zaman gösterilen dependency section
- ⚠️ **Smart Warnings**: Tamamlanmamış bağımlılık uyarıları

### Unlimited Subtask Hierarchy (v0.3.4 NEW!)
- 🌳 **Infinite Nesting**: Sınırsız derinlikte subtask oluşturma
- 📊 **Progress Tracking**: Ana task'ların otomatik progress hesaplaması
  - `📎 2/5`: 5 subtask'tan 2'si tamamlandı
- 🔄 **Visual Hierarchy**: TreeView'da indentasyon ile hiyerarşi gösterimi
- 🎯 **Smart Business Rules**:
  - Ana task'lar tüm subtask'lar tamamlanmadan completion'a geçemez
  - Subtask'ı olan task'lar silinemez
  - Subtask'lar parent'ın projesini otomatik inherit eder
- 🏗️ **MCP Integration**: Server-side hierarchy support ile güçlü backend

### Zengin Görev Detayı
- 📝 Split-view markdown editörü
- 👁️ Canlı önizleme
- 🔗 Bağımlılık bilgileri (her zaman görünür)
- 📊 Gelişmiş progress indicator (fixed v0.3.3)
- 🏷️ Template alan göstergeleri
- 🎨 Enhanced theming (dark/light mode improvements)

### Template Wizard
- 🧙 Çok adımlı arayüz
- 🔍 Template arama ve filtreleme
- 📋 Dinamik form oluşturma
- ✅ Alan doğrulama
- 👁️ Oluşturma öncesi önizleme

### Görev Yönetimi
- ✅ Hızlı görev oluşturma (Ctrl+Shift+G)
- 📝 Görev detaylarını görüntüleme
- 🔄 Durum güncelleme
- 🗑️ Toplu silme işlemleri
- 🏷️ Etiket yönetimi
- 📅 Son tarih takibi
- 🔗 Bağımlılık yönetimi

### Proje Yönetimi
- 📁 Proje oluşturma ve yönetimi
- 🎯 Aktif proje seçimi
- 📊 Proje bazlı görev görüntüleme
- 📈 Görev istatistikleri

### Template Sistemi
- 📋 Hazır görev şablonları (Bug, Feature, Technical Debt, Research)
- ⚡ Template wizard ile hızlı görev oluşturma
- 🔧 Özelleştirilebilir alanlar
- 🎨 Kategori bazlı organizasyon

## 📦 Kurulum

1. VS Code'da extension'ı yükleyin
2. Gorev MCP sunucusunun yolunu ayarlayın:
   ```json
   "gorev.serverPath": "/path/to/gorev"
   ```

## ⚙️ Konfigürasyon

### Temel Ayarlar
| Ayar | Açıklama | Varsayılan |
|------|----------|------------|
| `gorev.serverPath` | Gorev sunucu yolu | - |
| `gorev.autoConnect` | Başlangıçta otomatik bağlan | `true` |
| `gorev.showStatusBar` | Status bar'ı göster | `true` |
| `gorev.refreshInterval` | Otomatik yenileme aralığı (saniye) | `30` |

### TreeView Ayarları
| Ayar | Açıklama | Varsayılan |
|------|----------|------------|
| `gorev.treeView.grouping` | Gruplama stratejisi | `status` |
| `gorev.treeView.sorting` | Sıralama kriteri | `priority` |
| `gorev.treeView.sortAscending` | Artan sıralama | `false` |
| `gorev.treeView.showCompleted` | Tamamlanan görevleri göster | `true` |
| `gorev.treeView.showEmptyGroups` | Boş grupları göster | `false` |

### Drag & Drop Ayarları
| Ayar | Açıklama | Varsayılan |
|------|----------|------------|
| `gorev.dragDrop.allowTaskMove` | Görev taşımaya izin ver | `true` |
| `gorev.dragDrop.allowStatusChange` | Durum değiştirmeye izin ver | `true` |
| `gorev.dragDrop.allowPriorityChange` | Öncelik değiştirmeye izin ver | `true` |
| `gorev.dragDrop.allowProjectMove` | Projeler arası taşıma | `true` |
| `gorev.dragDrop.allowDependencyCreate` | Bağımlılık oluşturma | `true` |
| `gorev.dragDrop.showDropIndicator` | Drop göstergelerini göster | `true` |
| `gorev.dragDrop.animateOnDrop` | Drop animasyonları | `true` |

## 🎮 Klavye Kısayolları

- `Ctrl+Shift+G`: Hızlı görev oluştur
- `Ctrl+Shift+P`: Proje oluştur
- `Ctrl+Shift+T`: Template wizard'ı aç
- `Ctrl+R`: Görevleri yenile (TreeView odaktayken)
- `F2`: Görev başlığını düzenle
- `Delete`: Seçili görevi sil
- `Enter`: Görev detaylarını göster

## 📋 Komutlar

### Görev Komutları
- `Gorev: Create Task` - Yeni görev oluştur
- `Gorev: Edit Task` - Görevi düzenle
- `Gorev: Delete Task` - Görevi sil
- `Gorev: Complete Task` - Görevi tamamla
- `Gorev: Start Task` - Göreve başla
- `Gorev: Show Task Detail` - Görev detayını göster
- `Gorev: Add Dependency` - Bağımlılık ekle

### Proje Komutları
- `Gorev: Create Project` - Yeni proje oluştur
- `Gorev: Set Active Project` - Aktif proje seç
- `Gorev: Clear Active Project` - Aktif projeyi kaldır

### Template Komutları
- `Gorev: Create Task from Template` - Template'den görev oluştur
- `Gorev: Show Template Wizard` - Template wizard'ı göster
- `Gorev: Refresh Templates` - Template'leri yenile

### Genel Komutlar
- `Gorev: Connect to Server` - Sunucuya bağlan
- `Gorev: Disconnect` - Bağlantıyı kes
- `Gorev: Refresh` - Tüm verileri yenile
- `Gorev: Show Summary` - Özet bilgileri göster
- `Gorev: Show Search Input` - Arama kutusunu göster
- `Gorev: Show Advanced Filter` - Gelişmiş filtreleri göster
- `Gorev: Toggle Grouping` - Gruplama modunu değiştir
- `Gorev: Clear Filters` - Tüm filtreleri temizle

## 🛠️ Geliştirme

```bash
# Bağımlılıkları yükle
npm install

# TypeScript'i derle
npm run compile

# Watch modunda çalıştır
npm run watch

# Extension'ı paketle
npm run package
```

### 🧪 Test

```bash
# Tüm testleri çalıştır
npm test

# Test coverage raporu
npm run test-coverage

# Watch modunda test
npm run test-watch
```

Test suite şunları içerir:
- **Unit Tests**: Markdown parser, MCP client, tree providers
- **Integration Tests**: Extension activation, command registration
- **E2E Tests**: Tam kullanıcı iş akışları

## 📝 Lisans

MIT

## 🆕 What's New in v0.3.4

### Major Features Added:
- 🎯 **Task Dependencies**: Visual dependency tracking with TreeView badges
- 🌳 **Unlimited Subtask Hierarchy**: Infinite nesting with visual tree structure
- 📊 **Smart Progress Tracking**: Automatic parent task completion based on subtasks  
- 🔗 **Add Dependency Command**: Easy dependency creation via context menu
- 📋 **Always-Visible Dependencies**: Enhanced task detail panel 
- 🐛 **Progress Display Fix**: Circular progress percentage now visible
- 🎨 **Theme Improvements**: Better dark/light mode support
- ⚡ **Performance**: Optimized TreeView rendering and parsing

### Subtask System Highlights:
- **Infinite Depth**: Create subtasks under subtasks with no limits
- **Visual Hierarchy**: TreeView shows indented structure with progress indicators
- **Business Rules**: Smart completion and deletion constraints
- **MCP Backend**: Server-side hierarchy support with recursive queries

### Bug Fixes:
- Fixed progress percentage display in task detail panel
- Enhanced hierarchy parsing with flexible pattern matching
- Improved dependency section visibility
- Fixed circular progress chart rendering
- Enhanced CSP-compliant event handling

## 🤝 Katkıda Bulunma

Pull request'ler kabul edilir. Büyük değişiklikler için lütfen önce bir issue açın.

## ⚠️ BREAKING CHANGE: Template Zorunluluğu

### v0.4.0'dan İtibaren Template Kullanımı Zorunludur!

`gorev_olustur` komutu artık kullanılamaz. Tüm görevler template kullanılarak oluşturulmalıdır.

#### 🔄 Eski Kullanım (Artık Çalışmaz):
- "Create Task" (Ctrl+Shift+G) - Eskiden dialog açardı
- "Quick Create Task" - Eskiden hızlı görev oluşturma dialog'u açardı

#### ✅ Yeni Kullanım (Zorunlu):
- **"Create Task" (Ctrl+Shift+G)** → Template Wizard'ı açar
- **"Quick Create Task"** → Hızlı template seçimi açar
- Context menu'den **"Create from Template"** seçeneğini kullanın

#### 🆕 Kullanılabilir Template'ler:
- **Bug Raporu v2** - Detaylı bug takibi (severity, steps, environment)
- **Spike Araştırma** - Time-boxed teknik araştırmalar
- **Performans Sorunu** - Performans optimizasyon görevleri
- **Güvenlik Düzeltmesi** - Güvenlik açığı düzeltmeleri
- **Refactoring** - Kod kalitesi iyileştirmeleri
- **Ve diğer standart template'ler**...

#### 🎯 Neden Template Zorunlu?
- **Tutarlılık**: Her görev belirli standartlara uygun
- **Kalite**: Zorunlu alanlar eksik bilgi girişini engeller
- **Otomasyon**: Template tipine göre otomatik workflow
- **Raporlama**: Görev tipine göre detaylı metrikler