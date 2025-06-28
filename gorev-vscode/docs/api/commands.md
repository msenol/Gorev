# Commands API Reference

Gorev VS Code Extension tarafından sağlanan komutların detaylı API dokümantasyonu.

## Komut Listesi

### gorev.createTask

Yeni görev oluşturma dialogu açar.

**Command ID**: `gorev.createTask`

**Parameters**:
- `project?: Proje` - Opsiyonel proje. Belirtilmezse aktif proje kullanılır.

**Usage**:
```typescript
// Basit kullanım
vscode.commands.executeCommand('gorev.createTask');

// Proje ile kullanım
const project = { id: 'prj-123', isim: 'My Project' };
vscode.commands.executeCommand('gorev.createTask', project);
```

**Dialog Fields**:
- Başlık (zorunlu)
- Açıklama (opsiyonel, markdown)
- Öncelik (dropdown: Düşük/Orta/Yüksek)
- Son Tarih (date picker)
- Etiketler (comma-separated)

---

### gorev.quickCreateTask

Hızlı görev oluşturma - sadece başlık ister.

**Command ID**: `gorev.quickCreateTask`  
**Keyboard Shortcut**: `Ctrl+Shift+G` (Windows/Linux), `Cmd+Shift+G` (macOS)

**Parameters**: None

**Usage**:
```typescript
vscode.commands.executeCommand('gorev.quickCreateTask');
```

**Behavior**:
- Tek input box ile başlık alır
- Varsayılan değerler: Orta öncelik, aktif proje
- Başarılı olunca bildirim gösterir

---

### gorev.refreshTasks

Tüm TreeView'ları yeniler.

**Command ID**: `gorev.refreshTasks`

**Parameters**: None

**Usage**:
```typescript
vscode.commands.executeCommand('gorev.refreshTasks');
```

**Behavior**:
- Görev listesini yeniler
- Proje listesini yeniler
- Şablon listesini yeniler
- Status bar'ı günceller

---

### gorev.createProject

Yeni proje oluşturma dialogu açar.

**Command ID**: `gorev.createProject`

**Parameters**: None

**Usage**:
```typescript
vscode.commands.executeCommand('gorev.createProject');
```

**Dialog Fields**:
- İsim (zorunlu)
- Tanım (opsiyonel)

---

### gorev.setActiveProject

Projeyi aktif yapar.

**Command ID**: `gorev.setActiveProject`

**Parameters**:
- `project: Proje` - Aktif yapılacak proje

**Usage**:
```typescript
const project = { id: 'prj-123', isim: 'My Project' };
vscode.commands.executeCommand('gorev.setActiveProject', project);
```

**Context Menu**: Projects TreeView'da sağ tık menüsünde

---

### gorev.showTaskDetail

Görev detaylarını markdown formatında gösterir.

**Command ID**: `gorev.showTaskDetail`

**Parameters**:
- `task: Gorev` - Detayı gösterilecek görev

**Usage**:
```typescript
const task = { id: 'tsk-123', baslik: 'My Task' };
vscode.commands.executeCommand('gorev.showTaskDetail', task);
```

**Display Format**:
```markdown
# Görev Başlığı

**ID**: tsk-123  
**Durum**: Beklemede  
**Öncelik**: Yüksek  
**Proje**: My Project  
**Son Tarih**: 2025-07-15  
**Etiketler**: frontend, bug

## Açıklama
Görev açıklaması...

## Bağımlılıklar
- ✅ Bağımlı görev 1
- ⏳ Bağımlı görev 2
```

---

### gorev.updateTaskStatus

Görev durumunu güncelleme menüsü açar.

**Command ID**: `gorev.updateTaskStatus`

**Parameters**:
- `task: Gorev` - Durumu güncellenecek görev

**Usage**:
```typescript
const task = { id: 'tsk-123', baslik: 'My Task', durum: 'beklemede' };
vscode.commands.executeCommand('gorev.updateTaskStatus', task);
```

**Status Options**:
- Beklemede
- Devam Ediyor
- Tamamlandı

**Validation**:
- Bağımlı görevler tamamlanmadan "Devam Ediyor" seçilemez

---

### gorev.deleteTask

Görevi onay dialogu ile siler.

**Command ID**: `gorev.deleteTask`

**Parameters**:
- `task: Gorev` - Silinecek görev

**Usage**:
```typescript
const task = { id: 'tsk-123', baslik: 'My Task' };
vscode.commands.executeCommand('gorev.deleteTask', task);
```

**Confirmation Dialog**:
- Title: "Görevi Sil"
- Message: "'{task.baslik}' görevini silmek istediğinizden emin misiniz?"
- Buttons: Yes/No

---

### gorev.showSummary

Özet istatistikleri modal dialog'da gösterir.

**Command ID**: `gorev.showSummary`

**Parameters**: None

**Usage**:
```typescript
vscode.commands.executeCommand('gorev.showSummary');
```

**Display Format**:
```
📊 Görev Özeti

Toplam Görev: 25
✅ Tamamlanan: 10
🔄 Devam Eden: 5
📋 Bekleyen: 10

Proje Sayısı: 3
Aktif Proje: My Project
```

---

### gorev.connect

MCP server'a bağlanır.

**Command ID**: `gorev.connect`

**Parameters**: None

**Usage**:
```typescript
vscode.commands.executeCommand('gorev.connect');
```

**Behavior**:
- Server path'i kontrol eder
- Bağlantı kurar
- Status bar'ı günceller
- TreeView'ları yeniler

---

### gorev.disconnect

MCP server bağlantısını keser.

**Command ID**: `gorev.disconnect`

**Parameters**: None

**Usage**:
```typescript
vscode.commands.executeCommand('gorev.disconnect');
```

**Behavior**:
- Aktif bağlantıyı kapatır
- Status bar'ı günceller
- TreeView'ları temizler

---

### gorev.createFromTemplate

Şablondan görev oluşturur.

**Command ID**: `gorev.createFromTemplate`

**Parameters**:
- `template: Template` - Kullanılacak şablon

**Usage**:
```typescript
const template = { 
  id: 'bug-report', 
  isim: 'Bug Raporu',
  alanlar: [...]
};
vscode.commands.executeCommand('gorev.createFromTemplate', template);
```

**Dynamic Form**:
- Şablon alanlarına göre dinamik form oluşturur
- Alan tiplerine göre input kontrolü (text, select, date, number)
- Zorunlu alan validasyonu

---

## Context Menu Commands

TreeView'larda sağ tık menüsünde görünen komutlar:

### Tasks TreeView
- `gorev.showTaskDetail` - Detayları Göster
- `gorev.updateTaskStatus` - Durumu Güncelle
- `gorev.deleteTask` - Görevi Sil

### Projects TreeView
- `gorev.setActiveProject` - Aktif Yap
- `gorev.createTask` - Bu Projede Görev Oluştur

### Templates TreeView
- `gorev.createFromTemplate` - Bu Şablondan Oluştur

## Command Registration

Extension aktivasyonunda komutlar şu şekilde kaydedilir:

```typescript
export function activate(context: vscode.ExtensionContext) {
    // Komut kayıtları
    context.subscriptions.push(
        vscode.commands.registerCommand('gorev.createTask', createTask),
        vscode.commands.registerCommand('gorev.quickCreateTask', quickCreateTask),
        // ... diğer komutlar
    );
}
```

## Error Handling

Tüm komutlar hata yönetimi içerir:

```typescript
try {
    // Komut işlemleri
} catch (error) {
    vscode.window.showErrorMessage(`Hata: ${error.message}`);
    outputChannel.appendLine(`Error in command: ${error}`);
}
```

## Extending Commands

Yeni komut eklemek için:

1. `package.json` dosyasına komut tanımı ekleyin
2. `src/commands/` klasöründe handler fonksiyonu oluşturun
3. `extension.ts` dosyasında komutu kaydedin
4. İsteğe bağlı olarak context menu'ye ekleyin

---

<div align="center">

📚 Daha fazla bilgi için [Extension Dokümantasyonu](../../README.md)

</div>