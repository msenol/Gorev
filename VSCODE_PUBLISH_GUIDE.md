# 📦 VS Code Extension Marketplace Yayınlama Rehberi

## 1. Publisher Hesabı Oluşturma

1. **Azure DevOps'a Giriş**:
   - https://dev.azure.com/ adresine gidin
   - Microsoft hesabıyla giriş yapın

2. **Personal Access Token (PAT) Oluşturma**:
   - User Settings → Personal Access Tokens
   - "New Token" → Name: `vsce-publish`
   - Organization: All accessible organizations
   - Expiration: 90 days (veya daha uzun)
   - Scopes: Custom defined → "Marketplace" → "Manage" ✓
   - Token'ı kopyalayın ve güvenli saklayın!

3. **Publisher Oluşturma**:
   - https://marketplace.visualstudio.com/manage
   - "Create Publisher" tıklayın
   - Publisher ID: `msenol`
   - Display Name: `Mehmet Şenol`
   - Description: `Gorev Task Management System Developer`

## 2. VSCE ile Giriş

```bash
# Token ile giriş
vsce login msenol
# PAT token'ınızı yapıştırın
```

## 3. Extension'ı Yayınlama

```bash
cd gorev-vscode

# Önce paketleyin (test için)
vsce package

# Sonra yayınlayın
vsce publish

# Veya specific version ile
vsce publish 0.2.0

# Veya minor version bump ile
vsce publish minor
```

## 4. Marketplace URL'si

Yayınlandıktan sonra extension'ınız şu adreste olacak:
```
https://marketplace.visualstudio.com/items?itemName=msenol.gorev-vscode
```

## 5. VS Code'da Kurulum

```bash
# Command line'dan
code --install-extension msenol.gorev-vscode

# Veya VS Code içinden
# Extensions → "gorev" ara → Install
```

## 6. Güncelleme Yayınlama

```bash
# Version bump et
npm version patch  # 0.2.0 → 0.2.1
# veya
npm version minor  # 0.2.0 → 0.3.0

# Tekrar yayınla
vsce publish
```

## 7. İstatistikleri Görüntüleme

https://marketplace.visualstudio.com/manage/publishers/msenol

- İndirme sayıları
- Ratings
- Reviews
- Version history

## ⚠️ Dikkat Edilecekler

1. **Icon**: 128x128 veya 256x256 PNG olmalı (SVG kabul edilmiyor)
2. **README**: Marketplace'de görünecek, kaliteli olmalı
3. **Version**: Her publish'te version artmalı
4. **Categories**: Doğru kategoriler seçili olmalı
5. **Keywords**: Aranabilirlik için önemli

## 🎯 Checklist

- [ ] Publisher hesabı oluşturuldu
- [ ] PAT token alındı
- [ ] package.json metadata'ları eksiksiz
- [ ] Icon PNG formatında
- [ ] README profesyonel
- [ ] .vscodeignore dosyası var
- [ ] Test edildi
- [ ] Version güncellendi