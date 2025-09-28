# Gorev VS Code Extension - Alternatif IDE Kurulum Rehberi

Bu rehber, Gorev VS Code extension'ını VS Code dışındaki VS Code tabanlı editörlere nasıl kuracağınızı açıklar.

## 📦 VSIX Dosyası İndirme

İlk olarak VSIX dosyasını indirin:

### Yöntem 1: GitHub Release'den İndirme

1. [Gorev Releases](https://github.com/msenol/gorev/releases) sayfasına gidin
2. En son release'i bulun
3. Assets bölümünden `gorev-vscode-x.x.x.vsix` dosyasını indirin

### Yöntem 2: VS Code Marketplace'ten İndirme

1. [VS Code Marketplace - Gorev](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode) sayfasına gidin
2. "Download Extension" linkine **sağ tıklayın**
3. "Save Link As..." / "Bağlantıyı Farklı Kaydet" seçin
4. `.vsix` uzantılı olarak kaydedin

### Yöntem 3: Direct Link

```bash
# Terminal/PowerShell'de
curl -L -o gorev-vscode.vsix https://marketplace.visualstudio.com/_apis/public/gallery/publishers/mehmetsenol/vsextensions/gorev-vscode/latest/vspackage
```

## 🚀 IDE'ye Göre Kurulum

### Cursor

1. Command Palette açın: `Cmd/Ctrl + Shift + P`
2. "Extensions: Install from VSIX..." yazın
3. İndirdiğiniz `.vsix` dosyasını seçin
4. Cursor'ı yeniden başlatın

### Windsurf

1. Command Palette: `Cmd/Ctrl + Shift + P`
2. "Install from VSIX" komutunu arayın
3. VSIX dosyasını seçin
4. IDE'yi yeniden başlatın

### VSCodium

```bash
# Terminal'de
codium --install-extension gorev-vscode.vsix

# Veya GUI üzerinden
# Extensions panel → ⋯ (üç nokta) → Install from VSIX...
```

### Theia

```bash
# Terminal'de
theia extension:install gorev-vscode.vsix
```

### Code-Server (Browser tabanlı VS Code)

```bash
# Server'da
code-server --install-extension gorev-vscode.vsix

# Veya upload ederek
# Settings → Extensions → Install from VSIX → Upload
```

## 🔧 Manuel Kurulum (Tüm IDE'ler için)

Eğer yukarıdaki yöntemler çalışmazsa:

1. Extension klasörünü bulun:
   - **Windows**: `%USERPROFILE%\.vscode-oss\extensions\` veya `%APPDATA%\[IDE-ADI]\extensions\`
   - **macOS**: `~/.vscode-oss/extensions/` veya `~/Library/Application Support/[IDE-ADI]/extensions/`
   - **Linux**: `~/.vscode-oss/extensions/` veya `~/.config/[IDE-ADI]/extensions/`

2. VSIX dosyasını zip olarak açın:

   ```bash
   # VSIX aslında bir zip dosyasıdır
   unzip gorev-vscode.vsix -d gorev-temp
   ```

3. Extension klasörüne kopyalayın:

   ```bash
   # Linux/macOS
   cp -r gorev-temp/extension ~/.vscode-oss/extensions/mehmetsenol.gorev-vscode-0.2.0

   # Windows PowerShell
   Copy-Item -Recurse gorev-temp\extension "$env:USERPROFILE\.vscode-oss\extensions\mehmetsenol.gorev-vscode-0.2.0"
   ```

4. IDE'yi yeniden başlatın

## ⚙️ Extension Ayarları

Kurulumdan sonra, MCP server yolunu ayarlayın:

1. Settings/Preferences açın: `Cmd/Ctrl + ,`
2. "gorev" arayın
3. **Gorev: Server Path** ayarını yapın:
   - Linux/macOS: `gorev` (PATH'te ise) veya `/usr/local/bin/gorev`
   - Windows: `gorev` (PATH'te ise) veya `C:\Users\[USERNAME]\AppData\Local\Programs\gorev\gorev.bat`

## 🐛 Sorun Giderme

### Extension görünmüyor

- IDE'yi tamamen kapatıp açın
- Extension klasörünü kontrol edin
- IDE'nin VS Code API uyumluluğunu kontrol edin

### MCP server bağlanmıyor

- Terminal'de `gorev version` komutunu test edin
- Server path'ini tam yol olarak girin
- Output panelinde "Gorev" kanalını kontrol edin

### Uyumluluk sorunu

- Extension `package.json` dosyasında minimum VS Code versiyonunu kontrol edin
- IDE'nizin VS Code API versiyonunu kontrol edin

## 📝 Desteklenen IDE'ler

✅ **Tam Uyumlu:**

- Cursor
- Windsurf
- VSCodium
- Code-Server

⚠️ **Kısmi Uyumlu:**

- Theia (bazı özellikler çalışmayabilir)
- Eclipse Che (web tabanlı)

❌ **Uyumsuz:**

- Sublime Text (VS Code API yok)
- Atom (farklı extension sistemi)
- IntelliJ IDEA (farklı platform)

## 🔗 Faydalı Linkler

- [Gorev GitHub](https://github.com/msenol/gorev)
- [VS Code Marketplace](https://marketplace.visualstudio.com/items?itemName=mehmetsenol.gorev-vscode)
- [MCP Protokolü](https://github.com/modelcontextprotocol)
