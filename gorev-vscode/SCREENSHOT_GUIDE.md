# VS Code Marketplace Screenshot Guide

## Screenshot Requirements

- **Format**: PNG or JPG
- **Recommended Size**: 1280x800 pixels
- **Maximum Size**: 2MB per image
- **Location**: `images/screenshots/` directory

## Required Screenshots

1. **Main TreeView** (`screenshot-1-treeview.png`)
   - Show the Gorev TreeView with tasks organized by status
   - Include different priority indicators (🔥⚡ℹ️)
   - Show subtask hierarchy with progress bars
   - Include tags and due dates

2. **Task Detail Panel** (`screenshot-2-task-detail.png`)
   - Show rich task detail view
   - Include markdown editor
   - Show dependencies section
   - Display progress visualization

3. **Filtering Toolbar** (`screenshot-3-filters.png`)
   - Show the filter toolbar in action
   - Display search functionality
   - Show dropdown filters

4. **Command Palette** (`screenshot-4-commands.png`)
   - Show command palette with "Gorev:" commands
   - Highlight quick task creation (Ctrl+Shift+G)

5. **Status Bar** (`screenshot-5-statusbar.png`)
   - Show status bar with active project
   - Include task counts

6. **Dark Theme** (`screenshot-6-dark-theme.png`)
   - Same as screenshot 1 but in dark theme
   - Shows theme compatibility

## How to Add Screenshots

1. Take screenshots and save them in `images/screenshots/`
2. Update `package.json`:

```json
{
  "galleryBanner": {
    "color": "#2A2A2A",
    "theme": "dark"
  },
  "screenshots": [
    {
      "path": "images/screenshots/screenshot-1-treeview.png"
    },
    {
      "path": "images/screenshots/screenshot-2-task-detail.png"
    },
    {
      "path": "images/screenshots/screenshot-3-filters.png"
    },
    {
      "path": "images/screenshots/screenshot-4-commands.png"
    },
    {
      "path": "images/screenshots/screenshot-5-statusbar.png"
    },
    {
      "path": "images/screenshots/screenshot-6-dark-theme.png"
    }
  ]
}
```

3. Update version to 0.3.6
4. Run `vsce publish`

## Tips for Good Screenshots

- ✅ Use realistic task names and descriptions
- ✅ Show variety in task statuses and priorities
- ✅ Include Turkish content to show i18n support
- ✅ Clean, uncluttered workspace
- ✅ Good contrast and readability
- ❌ Avoid showing sensitive information
- ❌ Don't show error states
- ❌ Avoid empty states

## Example Task Content for Screenshots

```
Projeler:
- E-ticaret Sitesi (Aktif)
- Mobil Uygulama v2.0
- API Dokümantasyonu

Görevler:
- 🔥 [devam_ediyor] Ödeme sistemi entegrasyonu (3 gün kaldı)
  └─ ✅ Stripe API kurulumu
  └─ 🔄 Test senaryoları yazımı
  └─ ⏳ Production deployment
  
- ⚡ [beklemede] Kullanıcı profil sayfası
  🏷️ frontend, ui/ux
  
- ℹ️ [tamamlandı] Veritabanı şeması tasarımı
```
