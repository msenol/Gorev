# Gorev - Gelişmiş Görev Yönetimi ve AI Entegrasyonu VS Code için

<p align="center">
  <img src="media/icon.png" alt="Gorev Logo" width="128" height="128">
</p>

<div align="center">

[🇺🇸 English](README.md) | [🇹🇷 Türkçe](README.tr.md)

[![Version](https://img.shields.io/badge/Version-0.6.12-blue?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![Downloads](https://img.shields.io/visual-studio-marketplace/d/mehmetsenol.gorev-vscode?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![Rating](https://img.shields.io/visual-studio-marketplace/r/mehmetsenol.gorev-vscode?style=for-the-badge)](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**48 MCP aracı, sınırsız hiyerarşi ve sorunsuz AI asistan entegrasyonu ile nihai görev yönetimi güç merkezi**

> 🚀 **v0.6.12'de YENİ**: Sıfır kurulum NPX desteği! Extension artık yayınlanan @mehmetsenol/gorev-mcp-server paketini kullanarak anında kurulum sağlıyor. VS Code, Cursor, Windsurf ve diğer MCP destekli editörler için mükemmel çoklu IDE desteği.

</div>

## 🌟 Neden Gorev'i Seçmelisiniz?

Gorev, VS Code'u benzersiz yetenekleri olan **profesyonel bir görev yönetimi güç merkezine** dönüştürür:

- **🚀 Sıfır Kurulum NPX Desteği** - Binary indirme olmadan saniyeler içinde başlayın
- **🤖 48 MCP Aracı** - AI asistanlar için en kapsamlı görev yönetimi API'sı
- **🌳 Sınırsız Hiyerarşi** - Görsel ilerleme takibi ile sonsuz alt görev yuvalanması
- **🔍 Gelişmiş Arama** - Bulanık eşleştirme ve NLP ile FTS5 tam metin arama
- **🎯 Akıllı Bağımlılıklar** - Otomatik çözümleme ile görsel bağımlılık yönetimi
- **🌍 İki Dilli Destek** - Otomatik dil algılama ile 668 i18n anahtarı
- **📊 Veri Dışa/İçe Aktarma** - Çakışma çözümlemesi ile çok adımlı sihirbazlar
- **⚡ Ultra Performans** - %90 işlem azaltması ile RefreshManager

## 🚀 Sıfır Kurulum Ayarları

### 🎯 NPX Modu (Önerilen - İndirme Yok!)

Başlamanın en kolay yolu - binary kurulum gerektirmez:

1. **Extension'ı Kur**: VS Code marketplace'te "Gorev" ara
2. **Otomatik Yapılandırma**: Extension varsayılan olarak NPX modunu kullanır
3. **Çalışmaya Başla**: Hemen proje ve görevler oluştur!

Extension arka planda otomatik olarak `npx @mehmetsenol/gorev-mcp-server@latest` çalıştırır.

```json
// Varsayılan yapılandırma - kurulum gerekmez!
{
  "gorev.serverMode": "npx",     // Otomatik NPX çalıştırma
  "gorev.autoConnect": true      // Başlangıçta bağlan
}
```

### 🔧 Binary Modu (İleri Düzey Kullanıcılar)

Yerel binary kurulumu tercih eden kullanıcılar için:

```json
{
  "gorev.serverMode": "binary",
  "gorev.serverPath": "/path/to/gorev"
}
```

Binary kurulum için [kurulum kılavuzunu](https://github.com/msenol/Gorev/blob/main/README.md#-kurulum) takip edin.

## 🎯 Ana Özellikler Matrisi

| Kategori | Özellik | Açıklama | Durum |
|----------|---------|----------|-------|
| **🚀 Kurulum** | NPX Sıfır Kurulum | İndirme yok, anında kurulum | ✅ |
| **🤖 AI Entegrasyonu** | 48 MCP Aracı | AI asistanlar için tam API | ✅ |
| **🌳 Görev Yönetimi** | Sınırsız Hiyerarşi | Sonsuz alt görev yuvalanması | ✅ |
| **🔗 Bağımlılıklar** | Akıllı Çözümleme | Görsel bağımlılık yönetimi | ✅ |
| **🔍 Arama** | FTS5 Tam Metin | SQLite sanal tablolar, bulanık eşleştirme | ✅ |
| **📊 Veri Yönetimi** | Dışa/İçe Aktarma Sihirbazları | Çakışma çözümlemesi ile JSON/CSV | ✅ |
| **🎨 Görsel Arayüz** | Zengin TreeView | İlerleme çubukları, rozetler, renk kodlama | ✅ |
| **⚡ Performans** | RefreshManager | %90 işlem azaltması, geciktirme | ✅ |
| **🌍 Yerelleştirme** | İki Dilli Destek | 668 i18n anahtarı, otomatik algılama | ✅ |
| **💾 Veritabanı** | Workspace Modu | Proje özel veya global veritabanları | ✅ |
| **🎛️ Özelleştirme** | 50+ Ayar | Tam görsel ve davranışsal kontrol | ✅ |
| **🔄 Gerçek Zamanlı** | Dosya İzleme | Dosya değişikliklerinde otomatik güncellemeler | ✅ |

## 🤖 AI Asistan Entegrasyonu

### MCP Protokol Uyumluluğu

Tüm MCP uyumlu AI asistanlarıyla sorunsuz çalışır:

- **✅ Claude Desktop** - Tam konuşma entegrasyonu
- **✅ VS Code with MCP** - Yerel extension desteği
- **✅ Cursor IDE** - AI kodlama asistanı entegrasyonu
- **✅ Windsurf** - Geliştirme ortamı entegrasyonu
- **✅ Herhangi MCP İstemci** - Evrensel uyumluluk

### Doğal Dil Görev Yönetimi

AI asistanınızla doğal şekilde konuşun:

```
🗨️ "Dark mode implementasyonu için yüksek öncelikli yeni bir görev oluştur"
🗨️ "Bağımlılıkları olan tüm geciken görevleri göster"
🗨️ "#42 numaralı görevi tamamlandı olarak işaretle ve bağımlılıkları güncelle"
🗨️ "Login sorunu için bug raporu şablonu oluştur"
🗨️ "Geçen aydan tüm tamamlanan görevleri CSV'ye aktar"
```

### 48 MCP Aracı Kategorileri

| Kategori | Araçlar | Açıklama |
|----------|---------|----------|
| **Görev Yönetimi** | 6 araç | Oluştur, güncelle, listele, detay işlemleri |
| **Alt Görev İşlemleri** | 3 araç | Hiyerarşi yönetimi ve yuvalama |
| **Proje Yönetimi** | 6 araç | Proje oluşturma, etkinleştirme, istatistikler |
| **Şablon Sistemi** | 2 araç | Şablon tabanlı görev oluşturma |
| **Gelişmiş Arama** | 6 araç | FTS5 arama, öneriler, geçmiş |
| **Veri Dışa/İçe Aktarma** | 2 araç | Çok formatlı veri işlemleri |
| **Dosya İzleme** | 4 araç | Dosya sistemi izleme |
| **AI Bağlamı** | 6 araç | Bağlam yönetimi ve NLP |
| **IDE Entegrasyonu** | 5 araç | Extension yönetimi otomasyonu |
| **Gelişmiş İşlemler** | 8 araç | Toplu işlem, analitik |

## 🌳 Sınırsız Görev Hiyerarşisi

### Görsel Hiyerarşi Yönetimi

- **🔄 Sonsuz Yuvalama** - Limit olmadan görevler içinde görev oluştur
- **📊 İlerleme Takibi** - Ana görevler tamamlanma yüzdesini gösterir
- **🎯 Görsel Göstergeler** - Genişlet/daralt ile ağaç yapısı
- **⚡ Hızlı İşlemler** - Sürükle & bırak, satır içi düzenleme

### Hiyerarşi Örnekleri

```
📁 Proje: E-ticaret Platformu
├── 🚀 Kullanıcı Kimlik Doğrulama Sistemi (%75 tamamlandı)
│   ├── ✅ JWT Middleware Kurulumu
│   ├── ✅ Login Form Bileşeni
│   ├── 🔄 Şifre Doğrulama
│   │   ├── ⏳ Regex Pattern Implementasyonu
│   │   └── ⏳ Hata Mesajı Yerelleştirmesi
│   └── ⏳ Oturum Yönetimi
└── 📱 Mobil Responsive Tasarım (%25 tamamlandı)
    ├── ✅ Breakpoint Analizi
    └── ⏳ Bileşen Adaptasyonu
        ├── ⏳ Header Responsiveness
        └── ⏳ Navigasyon Menüsü
```

## 🔍 Gelişmiş Arama ve Filtreleme

### FTS5 Tam Metin Arama

Şimşek hızında arama için SQLite sanal tabloları:

- **🔍 İçerik Arama** - Başlık, açıklama, etiketlerde arama
- **🎯 Bulanık Eşleştirme** - Yazım hatalarıyla bile görev bul
- **🧠 NLP Entegrasyonu** - Doğal dil sorgu ayrıştırma
- **📊 Arama Analitiği** - Arama kalıplarını ve geçmişini takip et
- **💾 Kayıtlı Profiller** - Karmaşık filtre kombinasyonlarını sakla

### Filtreleme Yetenekleri

```typescript
// Gelişmiş filtreleme seçenekleri
{
  status: ["pending", "in_progress"],
  priority: ["high", "medium"],
  tags: ["bug", "urgent"],
  dateRange: {
    start: "2025-01-01",
    end: "2025-12-31"
  },
  project: "WebApp",
  hasDepencies: true,
  isOverdue: false
}
```

## 🔗 Akıllı Bağımlılık Yönetimi

### Görsel Bağımlılık Sistemi

- **🔒 Engellenmiş Görevler** - Engellenmiş görevler için net görsel göstergeler
- **🔓 Hazır Görevler** - Bağımlılıklar tamamlandığında otomatik çözümleme
- **🔗 Bağlantılı Görevler** - İki yönlü bağımlılık görselleştirmesi
- **⚡ Toplu İşlemler** - Aynı anda birden fazla bağımlılığı yönet

### Bağımlılık Türleri

| İkon | Durum | Açıklama |
|------|-------|----------|
| 🔒 | Engellenmiş | Tamamlanmamış bağımlılıkları var |
| 🔓 | Hazır | Tüm bağımlılıklar tamamlandı |
| 🔗 | Bağlantılı | İki yönlü bağlantıları var |
| ⚡ | Otomatik | Otomatik çözümleme etkin |

## 📊 Veri Dışa ve İçe Aktarma Sihirbazları

### Çok Adımlı Dışa Aktarma Sihirbazı

Rehberli kurulum ile gelişmiş dışa aktarma yetenekleri:

1. **📋 Format Seç** - JSON (yapısal) veya CSV (tablo)
2. **🎯 Kapsam Seç** - Mevcut görünüm, proje veya özel filtre
3. **📅 Tarih Aralığı** - Esnek tarih filtreleme seçenekleri
4. **🔧 Yapılandırma** - Bağımlılıklar, etiketler, metadata dahil et
5. **📤 Dışa Aktar** - VS Code bildirimleri ile ilerleme takibi

### Çakışma Çözümlemesi ile İçe Aktarma

Birden fazla çözümleme stratejisi ile akıllı içe aktarma sistemi:

- **🔄 Çakışmaları Atla** - Mevcut veriyi değiştirmeden bırak
- **📝 Üzerine Yaz** - İçe aktarılan veri ile değiştir
- **🔀 Birleştir** - Mevcut ve yeni verinin akıllı kombinasyonu
- **👀 Önizleme** - Uygulamadan önce değişiklikleri gör

### Dışa Aktarma Formatları

```json
// JSON Dışa Aktarma (yapısal)
{
  "tasks": [...],
  "projects": [...],
  "dependencies": [...],
  "metadata": {
    "exportDate": "2025-09-18",
    "version": "v0.6.12"
  }
}
```

```csv
// CSV Dışa Aktarma (tablo)
ID,Title,Status,Priority,Project,Tags,DueDate,Progress
1,"Auth Kurulumu",pending,high,"WebApp","güvenlik,auth","2025-10-01",0
```

## ⚡ Performans Optimizasyonları

### RefreshManager Mimarisi

%90 performans iyileştirmesi ile devrimsel yenileme sistemi:

- **🎯 Akıllı Gruplama** - Verimlilik için işlemleri grupla
- **⏱️ Öncelik Geciktirmesi** - Yüksek: 100ms, Normal: 500ms, Düşük: 2s
- **🔍 Diferansiyel Güncellemeler** - Hash tabanlı değişiklik algılama
- **📊 Performans İzleme** - Gerçek zamanlı işlem takibi
- **🚫 Sıfır Engelleme** - Engelleyici olmayan async işlemler

### Performans Metrikleri

| Metrik | Önce | Sonra | İyileştirme |
|--------|------|-------|-------------|
| Yenileme İşlemleri | 1000/dk | 100/dk | %90 azalma |
| UI Thread Engelleme | 50ms | 0ms | Sıfır engelleme |
| Bellek Kullanımı | 50MB | 35MB | %30 azalma |
| Başlangıç Süresi | 2s | 1.4s | %30 daha hızlı |

## 🎨 Zengin Görsel Arayüz

### Gelişmiş TreeView

Gelişmiş özelliklerle profesyonel seviye ağaç arayüzü:

- **📊 İlerleme Çubukları** - Ana görevler için görsel tamamlanma takibi
- **🎯 Öncelik Rozetleri** - Renk kodlu öncelik göstergeleri (🔥⚡ℹ️)
- **📅 Akıllı Tarihler** - Göreceli formatlama (Bugün, Yarın, 3g kaldı)
- **🔗 Bağımlılık İkonları** - Görsel bağımlılık durumu (🔒🔓🔗)
- **🏷️ Etiket Hapları** - Hover detayları ile renkli etiket rozetleri
- **📈 Zengin Tooltips** - İlerleme görselleştirmeli markdown tooltips

### Sürükle & Bırak İşlemleri

Tüm işlemler için sezgisel sürükle & bırak:

- **🔄 Görev Taşı** - Görsel geri bildirim ile projeler arası
- **📊 Durum Değiştir** - Durum gruplarına bırakarak güncelle
- **🎯 Öncelik Sırala** - Öncelik seviyelerini değiştirmek için sürükle
- **🔗 Bağımlılık Oluştur** - Bağımlılık oluşturmak için görevi diğerine bırak
- **✨ Görsel Geri Bildirim** - Düzgün animasyonlar ve bırakma göstergeleri

### Satır İçi Düzenleme

Ağaç görünümünden ayrılmadan hızlı düzenleme:

- **✏️ F2 veya Çift Tık** - Anında başlık düzenleme
- **📝 Bağlam Menüleri** - Durum/öncelik değişiklikleri için sağ tık
- **📅 Tarih Seçici** - Satır içi son tarih seçimi
- **⌨️ Klavye Kısayolları** - İptal için Escape, kaydetmek için Enter

## 🎛️ Kapsamlı Yapılandırma

### 50+ Yapılandırma Seçeneği

Extension'ın her yönü üzerinde tam kontrol:

#### Sunucu Ayarları (5 seçenek)
```json
{
  "gorev.serverMode": "npx|binary",
  "gorev.serverPath": "/path/to/gorev",
  "gorev.autoConnect": true,
  "gorev.showStatusBar": true,
  "gorev.refreshInterval": 300
}
```

#### TreeView Görseller (15 seçenek)
```json
{
  "gorev.treeView.grouping": "status|priority|project|tag|dueDate",
  "gorev.treeView.sorting": "title|priority|dueDate|created",
  "gorev.treeView.sortAscending": false,
  "gorev.treeView.showCompleted": true,
  "gorev.treeView.showEmptyGroups": false,
  "gorev.treeView.visuals.showProgressBars": true,
  "gorev.treeView.visuals.showPriorityBadges": true,
  "gorev.treeView.visuals.showDueDateIndicators": true,
  "gorev.treeView.visuals.showDependencyBadges": true,
  "gorev.treeView.visuals.showTagPills": true,
  "gorev.treeView.visuals.progressBarStyle": "blocks|percentage|both",
  "gorev.treeView.visuals.dueDateFormat": "relative|absolute|smart",
  "gorev.treeView.visuals.priorityBadgeStyle": "emoji|text|color",
  "gorev.treeView.visuals.tagPillLimit": 3,
  "gorev.treeView.visuals.showSubtaskProgress": true
}
```

#### Sürükle & Bırak (8 seçenek)
```json
{
  "gorev.dragDrop.allowTaskMove": true,
  "gorev.dragDrop.allowStatusChange": true,
  "gorev.dragDrop.allowPriorityChange": true,
  "gorev.dragDrop.allowProjectMove": true,
  "gorev.dragDrop.allowDependencyCreate": true,
  "gorev.dragDrop.allowParentChange": true,
  "gorev.dragDrop.showDropIndicator": true,
  "gorev.dragDrop.enableAnimation": true
}
```

#### Performans (8 seçenek)
```json
{
  "gorev.refreshManager.enableBatching": true,
  "gorev.refreshManager.batchSize": 10,
  "gorev.refreshManager.highPriorityDelay": 100,
  "gorev.refreshManager.normalPriorityDelay": 500,
  "gorev.refreshManager.lowPriorityDelay": 2000,
  "gorev.performance.enableMonitoring": true,
  "gorev.performance.slowOperationThreshold": 1000,
  "gorev.performance.maxMetrics": 1000
}
```

#### Veritabanı Modları (3 seçenek)
```json
{
  "gorev.databaseMode": "auto|workspace|global",
  "gorev.workspaceDatabase.autoDetect": true,
  "gorev.workspaceDatabase.showModeInStatusBar": true
}
```

## 🌍 Uluslararasılaşma

### Tam İki Dilli Destek

- **668 i18n Anahtarı** - Her UI elementi çevrildi
- **Otomatik Algılama** - VS Code dil ayarını takip eder
- **Diller**: İngilizce (en) ve Türkçe (tr)
- **Bağlam Duyarlı** - Kullanıma göre akıllı çeviriler

### Çeviri Örnekleri

| İngilizce | Türkçe | Bağlam |
|-----------|--------|--------|
| "Create Task" | "Görev Oluştur" | Komut |
| "High Priority" | "Yüksek Öncelik" | Öncelik rozeti |
| "Dependencies blocked" | "Bağımlılıklar engelledi" | Durum |
| "Export completed" | "Dışa aktarma tamamlandı" | Bildirim |

## 💾 Veritabanı Yönetimi

### Esnek Veritabanı Modları

#### Workspace Modu (Varsayılan)
- **📁 Proje Özel** - Her projenin kendi `.gorev/gorev.db`'si
- **🔍 Otomatik Algılama** - Workspace veritabanlarını otomatik bulur
- **📊 Durum Göstergesi** - Durum çubuğunda mevcut veritabanını gösterir

#### Global Mod
- **🌐 Paylaşılan Veritabanı** - Tüm projeler için tek veritabanı
- **🏠 Kullanıcı Dizini** - `~/.gorev/gorev.db`'de saklanır
- **🔄 Kolay Geçiş** - Komut ile modlar arası geçiş

#### Otomatik Mod
- **🤖 Akıllı Seçim** - Otomatik olarak en iyi veritabanını seçer
- **⬆️ Fallback Zinciri** - Workspace → Parent → Global
- **⚡ Sıfır Yapılandırma** - Kutudan çıktığı gibi çalışır

## 📋 50+ Mevcut Komut

### Görev İşlemleri (15 komut)
- `gorev.createTask` - Yeni görev oluştur
- `gorev.updateTaskStatus` - Görev durumunu güncelle
- `gorev.showTaskDetail` - Detaylı görev görünümü
- `gorev.deleteTask` - Görev sil
- `gorev.markAsCompleted` - Hızlı tamamlama
- `gorev.setTaskPriority` - Öncelik değiştir
- `gorev.addTaskTag` - Etiket ekle
- `gorev.setTaskDueDate` - Son tarih belirle
- `gorev.createSubtask` - Alt görev ekle
- `gorev.moveTo` - Projeye taşı
- `gorev.duplicateTask` - Görev klonla
- `gorev.addTaskNote` - Not ekle
- `gorev.linkTasks` - Bağımlılık oluştur
- `gorev.unlinkTasks` - Bağımlılık kaldır
- `gorev.showTaskHistory` - Geçmişi görüntüle

### Proje Yönetimi (8 komut)
- `gorev.createProject` - Yeni proje oluştur
- `gorev.setActiveProject` - Aktif proje belirle
- `gorev.showProjectStats` - İstatistikleri görüntüle
- `gorev.deleteProject` - Proje sil
- `gorev.renameProject` - Proje adını değiştir
- `gorev.archiveProject` - Projeyi arşivle
- `gorev.exportProject` - Proje verisini dışa aktar
- `gorev.duplicateProject` - Proje klonla

### Şablon Sistemi (7 komut)
- `gorev.openTemplateWizard` - Şablon sihirbazı
- `gorev.createFromTemplate` - Şablondan oluştur
- `gorev.quickCreateFromTemplate` - Hızlı şablon seçimi
- `gorev.refreshTemplates` - Şablonları yeniden yükle
- `gorev.initDefaultTemplates` - Varsayılanları başlat
- `gorev.showTemplateDetails` - Şablon detayları
- `gorev.exportTemplate` - Şablon dışa aktar

### Veri İşlemleri (4 komut)
- `gorev.exportData` - Veri dışa aktarma sihirbazı
- `gorev.importData` - Veri içe aktarma sihirbazı
- `gorev.exportCurrentView` - Mevcut görünümü dışa aktar
- `gorev.quickExport` - Hızlı dışa aktarma

### Filtre İşlemleri (10 komut)
- `gorev.showSearchInput` - Görev ara
- `gorev.showFilterMenu` - Filtre menüsü
- `gorev.showFilterProfiles` - Kayıtlı profiller
- `gorev.clearAllFilters` - Tüm filtreleri temizle
- `gorev.filterOverdue` - Gecikenleri göster
- `gorev.filterDueToday` - Bugün bitenleri göster
- `gorev.filterDueThisWeek` - Bu hafta bitenleri göster
- `gorev.filterHighPriority` - Yüksek öncelikli göster
- `gorev.filterActiveProject` - Aktif proje göster
- `gorev.filterByTag` - Etikete göre filtrele

### Debug Araçları (6 komut)
- `gorev.showDebugInfo` - Debug bilgisi
- `gorev.clearDebugLogs` - Logları temizle
- `gorev.testConnection` - MCP bağlantısını test et
- `gorev.refreshAllViews` - Zorla yenile
- `gorev.resetExtension` - Durumu sıfırla
- `gorev.generateTestData` - Test verisi oluştur

## 🔄 Dosya Sistemi Entegrasyonu

### Dosya İzleyici Yetenekleri

- **📁 Proje İzleme** - Proje dosyalarındaki değişiklikleri izle
- **🔄 Otomatik Güncellemeler** - Otomatik görev durumu geçişleri
- **⚡ Gerçek Zamanlı Senkronizasyon** - Dosya değişikliklerinde anında UI güncellemeleri
- **🎯 Seçmeli İzleme** - Hangi dosyaların izleneceğini yapılandır

### Entegrasyon Kalıpları

```javascript
// Dosya değişikliklerine göre otomatik durum güncellemeleri
.gitignore değişikliği → "Git Kurulumu" görevini güncelle
package.json değişikliği → "Bağımlılık Yapılandırması" görevini güncelle
README.md değişikliği → "Dokümantasyon" görevini güncelle
```

## 🏆 Gelişmiş Yetenekler

### Takma Adları Olan Şablon Sistemi

Yaygın görev türleri için önceden oluşturulmuş şablonlar:

- **🐛 Bug Raporu** (`bug`) - Yapısal bug dokümantasyonu
- **✨ Özellik İsteği** (`feature`) - Yeni özellik spesifikasyonları
- **🔬 Araştırma** (`research`) - İnceleme ve öğrenme görevleri
- **⚡ Spike** (`spike`) - Zaman sınırlı keşif
- **🔒 Güvenlik** (`security`) - Güvenlik ile ilgili görevler
- **🚀 Performans** (`performance`) - Optimizasyon görevleri
- **🔧 Refactoring** (`refactor`) - Kod iyileştirme görevleri
- **💳 Teknik Borç** (`debt`) - Kod borcu takibi

### Toplu İşlemler

Verimlilik için etkili toplu işlemler:

- **✅ Çoklu Seçim** - Çoklu seçim için Ctrl/Cmd+Tık
- **📊 Toplu Durum Güncelleme** - Birden fazla görev için durum değiştir
- **🗑️ Toplu Silme** - Aynı anda birden fazla görev sil
- **🏷️ Etiket Yönetimi** - Toplu etiket ekleme/kaldırma
- **📁 Proje Göçü** - Birden fazla görevi projeler arası taşı

### Analitik ve Raporlama

Proje içgörüleri için yerleşik analitik:

- **📊 İlerleme Takibi** - Görsel ilerleme grafikleri
- **⏱️ Zaman Analizi** - Görev tamamlama kalıpları
- **🎯 Öncelik Dağılımı** - Öncelik seviyesi analizi
- **📅 Son Tarih İçgörüleri** - Deadline uyumluluk takibi
- **👥 Bağımlılık Analizi** - Bağımlılık karmaşıklık metrikleri

## 🎮 Kullanım Örnekleri

### Başlangıç İş Akışı

```
1. 📦 Extension Kur → VS Code marketplace'te "Gorev" ara
2. 🚀 Otomatik Kurulum → Extension otomatik NPX modunu yapılandırır
3. 📁 Proje Oluştur → "Web Uygulaması Geliştirme"
4. 🎯 Görev Ekle → Yapısal görevler için şablon sihirbazını kullan
5. 🌳 Hiyerarşi İnşa Et → Sınırsız yuvalama ile alt görevler oluştur
6. 🔗 Bağımlılık Belirle → İş akışı için ilgili görevleri bağla
7. 📊 İlerleme Takip Et → Görsel ilerleme göstergelerini izle
8. 🔍 Arama Kullan → FTS5 arama ile görevleri hızlıca bul
9. 📤 Veri Dışa Aktar → CSV/JSON ile takımla ilerleme paylaş
```

### AI Asistan İş Akışı

```
🤖 "Merhaba Claude, proje görevlerimi organize etmekte yardıma ihtiyacım var"
🗨️ "Login form doğrulama sorunu için bug raporu görevi oluştur"
   → Önem derecesi, adımlar, ortam ile yapısal bug raporu oluşturur
🗨️ "Süresi geçmiş tüm yüksek öncelikli görevleri göster"
   → Dikkat gerektiren acil görevleri filtreler ve görüntüler
🗨️ "JWT middleware görevini tamamlandı olarak işaretle"
   → Durumu günceller ve bağımlı görevleri otomatik çözümler
🗨️ "Bu sprint'ten tüm tamamlanan görevleri CSV'ye aktar"
   → Sprint değerlendirme toplantısı için rapor oluşturur
```

### Gelişmiş Arama Örnekleri

```
🔍 "kimlik doğrulama bug yüksek"     → Başlık/açıklamalarda bulanık arama
🔍 "status:pending priority:high"   → Yapısal filtre sorgusu
🔍 "project:WebApp overdue"         → Proje özel geciken görevler
🔍 "tags:güvenlik,acil"             → Çoklu etiket kesişim araması
🔍 "created:geçen-hafta"            → Tarih göreceli arama
```

## 🛠️ Kurulum Yöntemleri

### Yöntem 1: VS Code Marketplace (Önerilen)
```
1. VS Code'u aç
2. Extensions'a git (Ctrl+Shift+X)
3. "Gorev" ara
4. Install'a tıkla
5. Hemen kullanmaya başla!
```

### Yöntem 2: Komut Satırı
```bash
code --install-extension mehmetsenol.gorev-vscode
```

### Yöntem 3: VSIX Dosyası
[GitHub Releases](https://github.com/msenol/Gorev/releases)'ten indir ve manuel kurulum yap.

## 🔧 Sorun Giderme

### Yaygın Sorunlar

**NPX Modu Çalışmıyor?**
```bash
# Node.js versiyonunu kontrol et (14+ gerekli)
node --version

# NPX'i doğrudan test et
npx @mehmetsenol/gorev-mcp-server@latest --version
```

**Binary Mod Bağlantı Sorunları?**
```bash
# Binary kurulumunu doğrula
gorev version

# Ayarlarda binary yolunu kontrol et
"gorev.serverPath": "/usr/local/bin/gorev"
```

**Extension Yüklenmiyor?**
1. VS Code Output → Gorev kanalını kontrol et
2. VS Code'u yeniden başlat
3. Çakışan extension'ları kontrol et
4. Extension ayarlarını sıfırla

### Debug Modu

Sorun giderme için debug logging'i etkinleştir:

```json
{
  "gorev.debug.enabled": true,
  "gorev.debug.logLevel": "debug",
  "gorev.debug.showInOutput": true
}
```

## 📈 Performans ve İstatistikler

### Extension Metrikleri

- **📊 Test Kapsama**: %100 (VS Code extension)
- **🎯 MCP Araçları**: 48 araç mevcut
- **🌍 Diller**: İngilizce + Türkçe desteği
- **💻 Platformlar**: Windows, macOS, Linux
- **⚡ Performans**: %90 işlem azaltması
- **🔧 Yapılandırma**: 50+ özelleştirilebilir ayar
- **📱 Komutlar**: 50+ mevcut komut
- **🎨 Görsel**: 15+ özelleştirme seçeneği

### Mimari Öne Çıkanlar

- **🏗️ TypeScript**: Tam tip güvenliği ile strict mode
- **🔒 Thread Safety**: Race-condition'suz işlemler
- **⚡ Async İşlemler**: Engelleyici olmayan UI etkileşimleri
- **📊 Bellek Verimli** - Akıllı önbellekleme ve temizlik
- **🔄 Reaktif Güncellemeler**: Olay güdümlü mimari
- **🎯 Modüler Tasarım**: Temiz endişe ayrımı

## 🤝 Entegrasyon Noktaları

### MCP İstemci Uyumluluğu

| İstemci | Durum | Özellikler |
|---------|-------|------------|
| **Claude Desktop** | ✅ Tam | Tüm 48 MCP aracı, konuşma entegrasyonu |
| **VS Code MCP** | ✅ Tam | Yerel extension, doğrudan entegrasyon |
| **Cursor IDE** | ✅ Tam | AI kodlama asistanı, bağlam farkındalığı |
| **Windsurf** | ✅ Tam | Geliştirme ortamı entegrasyonu |
| **Zed Editor** | 🔄 Planlandı | Gelecek MCP destek entegrasyonu |

### Geliştirme Araçları

- **Git Entegrasyonu** - Dosya değişikliklerini ve görev güncellemelerini takip et
- **Proje Şablonları** - Görev şablonları ile yeni projeler iskele
- **CI/CD Hooks** - Build ve deployment pipeline'ları ile entegre ol
- **Dokümantasyon** - Görev yapısından otomatik dokümantasyon oluştur

## 📚 Kaynaklar ve Destek

### Dokümantasyon
- 📖 [Ana Repository](https://github.com/msenol/Gorev) - Tam kaynak kod ve dokümanlar
- 🔧 [MCP Araçları Referansı](https://github.com/msenol/Gorev/blob/main/docs/mcp-araclari.md) - Tüm 48 araç dokümanlandı
- 📋 [Kurulum Kılavuzu](https://github.com/msenol/Gorev/blob/main/README.md#-kurulum) - Binary kurulum talimatları
- 🎯 [VS Code Extension Kılavuzu](https://github.com/msenol/Gorev/blob/main/docs/user-guide/vscode-extension.md) - Gelişmiş kullanım

### Topluluk ve Destek
- 🐛 [Issue Tracker](https://github.com/msenol/Gorev/issues) - Bug raporları ve özellik istekleri
- 💬 [Discussions](https://github.com/msenol/Gorev/discussions) - Topluluk tartışmaları
- ❓ [SSS](https://github.com/msenol/Gorev/wiki/FAQ) - Sık sorulan sorular
- 📧 [İletişim](mailto:me@mehmetsenol.dev) - Doğrudan geliştirici iletişimi

### Katkıda Bulunma
1. 🍴 Repository'yi fork'la
2. 🌿 Feature branch oluştur
3. ✨ Değişikliklerini yap
4. 🧪 Uygulanabilirse test ekle
5. 📝 Pull request gönder

## 📄 Lisans

Bu proje **MIT Lisansı** altında lisanslanmıştır - detaylar için [LICENSE](LICENSE) dosyasına bakın.

## 🙏 Teşekkürler

- **MCP Protokol** - Sorunsuz AI entegrasyonu sağladığı için
- **SQLite FTS5** - Güçlü tam metin arama yetenekleri için
- **VS Code API** - Genişletilebilir editör entegrasyonu için
- **Topluluk** - Geri bildirim, bug raporları ve özellik istekleri için

---

<div align="center">

**❤️ ile üretken geliştiriciler için yapıldı**

[⬆ Başa Dön](#gorev---gelişmiş-görev-yönetimi-ve-ai-entegrasyonu-vs-code-için)

**Şimdi dene:** [VS Code Marketplace'ten Kur](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)

</div>