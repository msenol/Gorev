# VS Code Extension Dokümantasyonu

Gorev VS Code Extension, MCP server'a görsel arayüz sağlayan TypeScript tabanlı bir VS Code eklentisidir.

## İçindekiler

- [Genel Bakış](#genel-bakış)
- [Mimari](#mimari)
- [Özellikler](#özellikler)
- [Kurulum](#kurulum)
- [Konfigürasyon](#konfigürasyon)
- [API Referansı](#api-referansı)
- [Geliştirme](#geliştirme)
- [Sorun Giderme](#sorun-giderme)

## Genel Bakış

Gorev VS Code Extension, kullanıcıların VS Code içinden doğrudan görev yönetimi yapmasını sağlar. MCP protokolü üzerinden Gorev server'a bağlanarak zengin bir kullanıcı deneyimi sunar.

### Temel Özellikler

- **TreeView Panelleri**: Görev, proje ve şablon yönetimi
- **Komut Paleti**: Hızlı erişim komutları
- **Status Bar**: Anlık durum bilgisi
- **Context Menüler**: Sağ tık işlemleri
- **Tema Desteği**: Öncelik bazlı renklendirme

## Mimari

### Genel Yapı

```
┌─────────────────────────────────────────┐
│           VS Code Extension             │
├─────────────────────────────────────────┤
│  Commands │ TreeViews │ StatusBar │ UI │
├─────────────────────────────────────────┤
│            MCP Client Layer             │
└──────────────────┬──────────────────────┘
                   │ MCP Protocol (stdio)
                   │
┌──────────────────▼──────────────────────┐
│          Gorev MCP Server               │
└─────────────────────────────────────────┘
```

### Katmanlar

#### 1. UI Katmanı
- **TreeView Providers**: Veri görselleştirme
- **Command Handlers**: Kullanıcı etkileşimleri
- **Status Bar**: Durum gösterimi
- **WebView** (Planlanan): Gelişmiş görev editörü

#### 2. MCP Client Katmanı
- **Client**: Server bağlantı yönetimi
- **Types**: TypeScript tip tanımları
- **Protocol**: JSON-RPC mesajlaşma

#### 3. Model Katmanı
- **Gorev**: Görev veri modeli
- **Proje**: Proje veri modeli
- **Template**: Şablon veri modeli
- **Common**: Paylaşılan tipler (Durum, Öncelik)

## Özellikler

### TreeView Panelleri

#### Görevler TreeView
```typescript
interface GorevTreeItem {
  id: string;
  baslik: string;
  durum: GorevDurum;
  oncelik: GorevOncelik;
  projeId?: string;
  sonTarih?: Date;
  etiketler?: string[];
}
```

**Özellikler:**
- Durum bazlı gruplandırma (Beklemede, Devam Ediyor, Tamamlandı)
- Öncelik renklendirmesi (Yüksek: kırmızı, Orta: sarı, Düşük: yeşil)
- Son tarih gösterimi
- Etiket badges

#### Projeler TreeView
```typescript
interface ProjeTreeItem {
  id: string;
  isim: string;
  gorevSayisi: number;
  aktif: boolean;
}
```

**Özellikler:**
- Aktif proje vurgulama
- Görev sayısı gösterimi
- Hızlı aktif yapma

#### Şablonlar TreeView
```typescript
interface TemplateTreeItem {
  id: string;
  isim: string;
  kategori: string;
  alanlar: TemplateAlan[];
}
```

**Özellikler:**
- Kategori bazlı gruplandırma
- Alan bilgisi önizleme
- Hızlı görev oluşturma

### Komutlar

| Komut | ID | Kısayol | Açıklama |
|-------|-----|---------|----------|
| Create Task | `gorev.createTask` | - | Yeni görev oluşturma formu |
| Quick Create Task | `gorev.quickCreateTask` | `Ctrl+Shift+G` | Hızlı görev oluşturma |
| Create Project | `gorev.createProject` | - | Yeni proje oluştur |
| Set Active Project | `gorev.setActiveProject` | - | Projeyi aktif yap |
| Show Task Detail | `gorev.showTaskDetail` | - | Görev detaylarını göster |
| Update Task Status | `gorev.updateTaskStatus` | - | Görev durumunu güncelle |
| Delete Task | `gorev.deleteTask` | - | Görevi sil |
| Show Summary | `gorev.showSummary` | - | İstatistikleri göster |
| Connect | `gorev.connect` | - | Server'a bağlan |
| Disconnect | `gorev.disconnect` | - | Bağlantıyı kes |
| Refresh | `gorev.refreshTasks` | - | Listeleri yenile |

### Status Bar

Status bar şu bilgileri gösterir:
- Bağlantı durumu (🟢 Bağlı / 🔴 Bağlı Değil)
- Toplam görev sayısı
- Tamamlanan görev sayısı
- Aktif proje adı

Tıklandığında özet istatistikleri gösterir.

## Kurulum

### Marketplace'den Kurulum (Planlanan)
```
1. VS Code Extensions panelini aç (Ctrl+Shift+X)
2. "Gorev Task Orchestrator" ara
3. Install butonuna tıkla
```

### Local Development Kurulumu
```bash
# Repository'yi klonla
git clone https://github.com/yourusername/gorev.git
cd gorev/gorev-vscode

# Bağımlılıkları yükle
npm install

# Development build
npm run compile

# VS Code'da test et
code .
# F5 tuşuna bas
```

### VSIX Dosyasından Kurulum
```bash
# VSIX paketi oluştur
npm run package

# VS Code'da yükle
code --install-extension gorev-vscode-0.1.0.vsix
```

## Konfigürasyon

### Extension Ayarları

```json
{
  // MCP server binary'sinin tam yolu
  "gorev.serverPath": "/usr/local/bin/gorev",
  
  // Başlangıçta otomatik bağlan
  "gorev.autoConnect": true,
  
  // Status bar'ı göster
  "gorev.showStatusBar": true,
  
  // Otomatik yenileme süresi (saniye)
  // 0 = devre dışı
  "gorev.refreshInterval": 30,
  
  // Debug loglarını göster
  "gorev.debug": false,
  
  // Server connection timeout (ms)
  "gorev.connectionTimeout": 5000,
  
  // TreeView'da gösterilecek max görev sayısı
  "gorev.maxTasksPerGroup": 100
}
```

### Renk Teması Özelleştirme

`settings.json` dosyanızda:

```json
{
  "workbench.colorCustomizations": {
    "gorev.highPriorityForeground": "#ff6b6b",
    "gorev.mediumPriorityForeground": "#ffd93d",
    "gorev.lowPriorityForeground": "#6bcf7f"
  }
}
```

## API Referansı

### MCP Client API

```typescript
interface MCPClient {
  connect(): Promise<void>;
  disconnect(): Promise<void>;
  call(method: string, params: any): Promise<any>;
  isConnected(): boolean;
  onConnectionChange(callback: (connected: boolean) => void): void;
}
```

### TreeView Provider API

```typescript
interface GorevTreeDataProvider {
  refresh(): void;
  getTreeItem(element: GorevTreeItem): vscode.TreeItem;
  getChildren(element?: GorevTreeItem): Promise<GorevTreeItem[]>;
  onDidChangeTreeData: vscode.Event<GorevTreeItem | undefined>;
}
```

### Command API

```typescript
interface GorevCommands {
  createTask(proje?: Proje): Promise<void>;
  updateTaskStatus(gorev: Gorev): Promise<void>;
  deleteTask(gorev: Gorev): Promise<void>;
  showTaskDetail(gorev: Gorev): Promise<void>;
}
```

## Geliştirme

### Gereksinimler
- Node.js 16+
- npm 7+
- VS Code 1.95.0+
- TypeScript 5.0+

### Development Workflow

1. **Setup**
   ```bash
   npm install
   npm run compile
   ```

2. **Watch Mode**
   ```bash
   npm run watch
   ```

3. **Testing**
   ```bash
   # Unit tests
   npm test
   
   # E2E tests
   npm run test:e2e
   ```

4. **Linting**
   ```bash
   npm run lint
   npm run format
   ```

### Debug Yapılandırması

`.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Run Extension",
      "type": "extensionHost",
      "request": "launch",
      "args": [
        "--extensionDevelopmentPath=${workspaceFolder}"
      ],
      "outFiles": [
        "${workspaceFolder}/dist/**/*.js"
      ],
      "preLaunchTask": "${defaultBuildTask}"
    }
  ]
}
```

### Test Yazma

```typescript
// src/test/suite/extension.test.ts
import * as assert from 'assert';
import * as vscode from 'vscode';

suite('Extension Test Suite', () => {
  test('Extension should be present', () => {
    assert.ok(vscode.extensions.getExtension('gorev.gorev-vscode'));
  });

  test('Should register all commands', async () => {
    const commands = await vscode.commands.getCommands();
    assert.ok(commands.includes('gorev.createTask'));
    assert.ok(commands.includes('gorev.refreshTasks'));
  });
});
```

## Sorun Giderme

### Bağlantı Sorunları

**Sorun**: Extension server'a bağlanamıyor

**Çözümler**:
1. Server path'inin doğru olduğunu kontrol et
2. Server'ın çalıştığını doğrula: `gorev serve`
3. Windows'ta tam path kullan: `C:\\Program Files\\gorev\\gorev.exe`
4. Output panelinde "Gorev" kanalını kontrol et

### TreeView Sorunları

**Sorun**: Görevler TreeView'da görünmüyor

**Çözümler**:
1. Refresh butonuna tıkla
2. Aktif proje seçili mi kontrol et
3. Server response'larını Output'ta kontrol et
4. Debug mode'u aç: `"gorev.debug": true`

### Performance Sorunları

**Sorun**: Extension yavaş çalışıyor

**Çözümler**:
1. `gorev.refreshInterval` değerini artır
2. `gorev.maxTasksPerGroup` değerini azalt
3. Gereksiz extension'ları devre dışı bırak

### Debug Logları

Output panelinde "Gorev" kanalını seçerek detaylı logları görüntüleyebilirsiniz:

```
[Gorev] Connecting to MCP server...
[Gorev] Server path: /usr/local/bin/gorev
[Gorev] Connection established
[Gorev] Calling method: gorev_listele
[Gorev] Response received: 15 tasks
```

## Katkıda Bulunma

### Yeni Özellik Ekleme

1. Feature branch oluştur
2. Kodu implement et
3. Test yaz
4. Dokümantasyon güncelle
5. PR aç

### Code Style

- TypeScript strict mode kullan
- ESLint kurallarına uy
- Prettier ile formatla
- Meaningful variable names kullan
- JSDoc yorumları ekle

### PR Checklist

- [ ] Testler yazıldı ve geçiyor
- [ ] Dokümantasyon güncellendi
- [ ] CHANGELOG.md'ye eklendi
- [ ] Lint hataları yok
- [ ] TypeScript hataları yok

## Gelecek Özellikler

### v0.2.0 (Planlanan)
- [ ] WebView görev editörü
- [ ] Drag & drop desteği
- [ ] Bulk operations
- [ ] Görev filtreleme UI

### v0.3.0 (Planlanan)
- [ ] Gantt chart görünümü
- [ ] Notification sistemi
- [ ] Keyboard shortcuts
- [ ] Export/Import özelliği

---

<div align="center">

📚 Daha fazla bilgi için [ana dokümantasyona](../README.md) bakın.

</div>