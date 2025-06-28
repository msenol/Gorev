# Gorev VS Code Extension Test Talimatları

## Test Adımları:

### 1. Gorev Server'ı Başlat

Yeni bir terminal aç ve:
```bash
cd /mnt/f/Development/Projects/task-orchestrator
./gorev serve
```

"Gorev MCP sunucusu başlatılıyor..." mesajını görmelisin.

### 2. VS Code'da Extension'ı Debug Et

1. VS Code'da `gorev-vscode` klasörünü aç
2. `F5` tuşuna bas (Run Extension)
3. Yeni VS Code penceresi açılacak

### 3. Extension'ı Yapılandır

Yeni pencerede:
1. `Ctrl+Shift+P` → "Preferences: Open Settings (JSON)"
2. Şu ayarları ekle:

```json
{
  "gorev.serverPath": "gorev",
  "gorev.autoConnect": false
}
```

**Windows'ta** isen tam path ver:
```json
{
  "gorev.serverPath": "F:\\Development\\Projects\\task-orchestrator\\gorev.exe"
}
```

### 4. Server'a Bağlan

1. `Ctrl+Shift+P` → "Gorev: Connect to Server"
2. Output panel'i kontrol et: View → Output → "Gorev"
3. Status bar'da "Gorev: Connected" görünmeli

### 5. Test Et

- Sol tarafta Gorev ikonu
- Create Project komutu
- Create Task komutu
- Show Summary komutu

## Sorun Giderme

### "Connection closed" hatası:
1. Server'ın çalıştığından emin ol
2. Path'in doğru olduğunu kontrol et
3. Output panel'deki debug loglarını incele

### Watch mode hatası:
Package.json'da npm scripts'leri düzeltildi, sorun çözülmüş olmalı.

## Debug Logları

Output panel'de şunları görebilirsin:
- Connection attempts
- MCP protocol messages
- Error details