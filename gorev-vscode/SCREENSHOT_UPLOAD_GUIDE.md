# VS Code Extension Screenshot Upload Guide

## 📸 Screenshot Ekleme Adımları

### 1. Screenshot'ları Hazırlayın

1. VS Code'u açın ve Gorev extension'ını aktif edin
2. Örnek görevler oluşturun (çeşitli durum ve önceliklerde)
3. Her bir özellik için screenshot alın:
   - **TreeView** görünümü
   - **Task Detail Panel**
   - **Filter Toolbar**
   - **Command Palette**
   - **Status Bar**
   - **Dark Theme** görünümü

### 2. Screenshot'ları Düzenleyin

- Boyut: 1280x800 px önerilen
- Format: PNG
- Dosya boyutu: Max 2MB
- İsimlendirme:
  - `screenshot-1-treeview.png`
  - `screenshot-2-task-detail.png`
  - vb.

### 3. Dosyaları Yerleştirin

```bash
# Screenshot dizinine kopyalayın
cp /path/to/your/screenshots/*.png gorev-vscode/images/screenshots/
```

### 4. Extension'ı Güncelleyin ve Yayınlayın

```bash
cd gorev-vscode

# Compile et
npm run compile

# Yeni VSIX paketi oluştur
npx vsce package

# Marketplace'e yayınla
npx vsce publish
```

## 📝 Screenshot İçeriği Önerileri

### TreeView Screenshot

```
Gorev
├── 📁 Tasks
│   ├── 🔵 Devam Ediyor (3)
│   │   ├── 🔥 Ödeme sistemi entegrasyonu ████░░ 75%
│   │   │   └─ ✅ Stripe API kurulumu
│   │   │   └─ 🔄 Test senaryoları
│   │   │   └─ ⏳ Production deploy
│   │   └── ⚡ Kullanıcı profil sayfası
│   ├── ⚪ Beklemede (5)
│   └── ✅ Tamamlandı (8)
├── 📁 Projects (3)
└── 📁 Templates (4)
```

### Task Detail Screenshot

- Başlık: "Ödeme sistemi entegrasyonu"
- Markdown editor açık
- Dependency bölümü görünür
- Progress göstergesi
- Tags: payment, critical, backend

### Filter Toolbar Screenshot

- Search box'ta "api" yazılı
- Status dropdown açık
- Priority filter seçili
- Clear filter butonu görünür

## 🎨 Renk ve Tema Önerileri

- Light theme için temiz, okunabilir arkaplan
- Dark theme için VS Code'un default dark theme'i
- Öncelik renklerinin net görünmesi
- Progress bar'ların belirgin olması

## ✅ Kontrol Listesi

- [ ] Tüm screenshot'lar 1280x800 px boyutunda
- [ ] Dosya boyutları 2MB'dan küçük
- [ ] Dosya isimleri doğru
- [ ] `images/screenshots/` dizinine kaydedildi
- [ ] package.json'da version 0.3.6
- [ ] Hassas bilgi içermiyor
- [ ] Profesyonel görünüm

## 🚀 Yayınlama

Screenshot'lar eklendikten sonra:

1. `npm run compile`
2. `npx vsce package`
3. `npx vsce publish`

Marketplace'te güncellemenin görünmesi 5-15 dakika sürebilir.

## 📌 Not

Screenshot'lar olmadan da extension yayınlanabilir. Screenshot'ları sonradan da ekleyebilirsiniz. Her güncelleme için version numarasını artırmayı unutmayın!
