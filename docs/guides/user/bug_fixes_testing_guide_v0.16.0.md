# Test Planı - Bug Fixes v0.16.0

**Build Date**: 2025-10-04T02:13:28Z  
**Git Commit**: 07b43a8  
**Binary Location**: `/mnt/c/tmp/gorev-test/binaries/gorev`

## 🐛 Düzeltilen Hatalar

### 1. Batch Update Handler (CRITICAL) ✅

**Sorun**: Handler nested `{id: "x", updates: {durum: "y"}}` bekliyordu  
**Düzeltme**: Artık flat format kullanıyor `{id: "x", durum: "y"}`

**Test Adımları**:

```bash
# MCP server başlat
cd /mnt/c/tmp/gorev-test/workspace
/mnt/c/tmp/gorev-test/binaries/gorev serve

# Başka terminalde test et:
# 1. Önce birkaç görev oluştur
# 2. Batch update ile birden fazla görevi güncelle

# Beklenen: Hatasız çalışmalı, tüm görevler güncellenmeli
```

**MCP Tool Test**:

```json
{
  "name": "gorev_batch_update",
  "arguments": {
    "updates": [
      {
        "id": "task-1",
        "durum": "tamamlandi"
      },
      {
        "id": "task-2",
        "durum": "devam_ediyor",
        "oncelik": "yuksek"
      }
    ]
  }
}
```

### 2. File Watching Persistence (HIGH) ✅

**Sorun**: Dosya izlemeleri sadece RAM'de tutuluyordu  
**Düzeltme**: Artık database'e kaydediliyor

**Test Adımları**:

```bash
# 1. Server başlat
/mnt/c/tmp/gorev-test/binaries/gorev serve

# 2. Bir göreve dosya izleme ekle
# 3. Server'ı kapat
# 4. Server'ı tekrar başlat
# 5. Dosya izlemelerini listele

# Beklenen: İzlemeler kaybolmamalı
```

**MCP Tool Test**:

```json
// 1. Dosya izleme ekle
{
  "name": "gorev_file_watch_add",
  "arguments": {
    "task_id": "your-task-id",
    "file_path": "/mnt/c/tmp/gorev-test/workspace/test.go"
  }
}

// 2. Server restart sonrası listele
{
  "name": "gorev_file_watch_list",
  "arguments": {
    "task_id": "your-task-id"
  }
}
```

### 3. Filter Profile Display (MEDIUM) ✅

**Sorun**: Profil listesi sadece ID ve isim gösteriyordu  
**Düzeltme**: Artık tüm detaylar gösteriliyor

**Test Adımları**:

```bash
# 1. Server başlat
# 2. Filtre profili kaydet
# 3. Profilleri listele

# Beklenen: ID, name, description, filters, use_count vs. görmeli
```

**MCP Tool Test**:

```json
// 1. Profil kaydet
{
  "name": "gorev_filter_profile_save",
  "arguments": {
    "name": "Test Profili",
    "description": "MCP server test filtresi",
    "filters": {
      "durum": ["beklemede"],
      "oncelik": ["yuksek"]
    }
  }
}

// 2. Profilleri listele
{
  "name": "gorev_filter_profile_list",
  "arguments": {}
}
```

## 📊 Başarı Kriterleri

- [ ] Batch update flat format ile çalışıyor
- [ ] Batch update hata durumlarında detaylı mesaj veriyor
- [ ] File watch'lar server restart sonrası korunuyor
- [ ] File watch database'e kaydediliyor
- [ ] Filter profile list tüm detayları gösteriyor
- [ ] Filter profile içindeki filter kriterleri görünüyor
- [ ] Tüm MCP araçları hatasız çalışıyor

## 🔍 Debug İpuçları

**Debug mode ile çalıştır**:

```bash
/mnt/c/tmp/gorev-test/binaries/gorev serve --debug
```

**Log dosyasını izle**:

```bash
tail -f /mnt/c/tmp/gorev-test/logs/gorev-debug.log
```

**Database'i kontrol et**:

```bash
sqlite3 /mnt/c/tmp/gorev-test/workspace/.gorev/gorev.db

# File watches tablosunu kontrol
SELECT * FROM task_file_paths;

# Filter profiles tablosunu kontrol
SELECT * FROM filter_profiles;
```

## 📝 Değişen Dosyalar

1. **internal/mcp/handlers.go**
   - GorevBatchUpdate: Nested updates yerine flat format
   - GorevFilterProfileList: Detaylı markdown output

2. **internal/gorev/veri_yonetici.go**
   - GorevDosyaYoluEkle: Yeni method
   - GorevDosyaYoluSil: Yeni method
   - GorevDosyaYollariGetir: Yeni method
   - DosyaYoluGorevleriGetir: Yeni method

3. **internal/gorev/file_watcher.go**
   - loadFromDatabase: Yeni method (startup'ta DB'den yükleme)
   - AddTaskPath: DB'ye kaydetme eklendi
   - RemoveTaskPath: DB'den silme eklendi

4. **Documentation**
   - docs/tr/mcp-araclari-ai.md: Batch update örnekleri düzeltildi
   - docs/api/MCP_TOOLS_REFERENCE.md: Format açıklaması güncellendi

## ✅ Sonuç

Tüm testler başarılı olduğunda, bu 3 critical/high bug tamamen çözülmüş olacak!
