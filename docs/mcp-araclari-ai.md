# 🤖 AI Context Management MCP Araçları

Bu dokümanda Gorev'in AI-optimized context management araçları detaylı olarak açıklanmıştır.

## Genel Bakış

Gorev v0.9.0 ile birlikte, AI asistanlarla daha verimli çalışmak için özel olarak tasarlanmış 6 yeni MCP aracı eklendi. Bu araçlar, görev durumlarını otomatik yönetir, bağlamı korur ve doğal dil sorgularını destekler.

## 🎯 gorev_set_active

Aktif görevi ayarlar ve otomatik durum geçişi sağlar.

### Parametreler
- `task_id` (string, zorunlu): Aktif olarak ayarlanacak görevin ID'si

### Özellikler
- Görev "beklemede" durumundaysa otomatik olarak "devam_ediyor" durumuna geçirir
- AI oturum bağlamında aktif görevi saklar
- Son 10 görevi recent tasks listesinde tutar

### Örnek Kullanım
```json
{
  "task_id": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
}
```

### Yanıt
```
✅ Görev f47ac10b-58cc-4372-a567-0e02b2c3d479 başarıyla aktif görev olarak ayarlandı.
```

## 📍 gorev_get_active

Mevcut aktif görevi detaylarıyla birlikte getirir.

### Parametreler
Parametre almaz.

### Yanıt Formatı
Markdown formatında görev detayları:
- Görev başlığı
- Durum, öncelik, proje bilgileri
- Açıklama (varsa)

### Örnek Yanıt
```markdown
# Aktif Görev: API Dokümantasyonu Yaz

## 📋 Genel Bilgiler
- **ID:** f47ac10b-58cc-4372-a567-0e02b2c3d479
- **Durum:** devam_ediyor
- **Öncelik:** yuksek
- **Proje:** backend-project
```

## 📚 gorev_recent

Son etkileşime geçilen görevleri listeler.

### Parametreler
- `limit` (number, opsiyonel): Döndürülecek görev sayısı (varsayılan: 5)

### Yanıt Formatı
Markdown liste formatında son görevler:
```markdown
## 📚 Son Görevler

1. **API Dokümantasyonu Yaz** (ID: f47ac10b)
   - Durum: devam_ediyor | Öncelik: yuksek

2. **Unit Test Ekle** (ID: a8b9c0d1)
   - Durum: beklemede | Öncelik: orta
```

## 📊 gorev_context_summary

AI oturum özetini ve istatistikleri sunar.

### Parametreler
Parametre almaz.

### Yanıt İçeriği
- Aktif görev bilgisi
- Oturum istatistikleri (oluşturulan, güncellenen, tamamlanan)
- Öncelikli görevler listesi
- Blokajlar (bağımlılık bekleyen görevler)

### Örnek Yanıt
```markdown
## 🤖 AI Oturum Özeti

### 🎯 Aktif Görev
**API Dokümantasyonu Yaz** (devam_ediyor)

### 📊 Oturum İstatistikleri
- Oluşturulan: 3
- Güncellenen: 5
- Tamamlanan: 2

### 🔥 Öncelikli Görevler
- **Kritik Bug Fix** (ID: abc123)
- **Performance Optimization** (ID: def456)

### 🚫 Blokajlar
- **Deploy to Production** (ID: ghi789) - 2 bağımlılık bekliyor
```

## 🔄 gorev_batch_update

Birden fazla görevi tek seferde günceller.

### Parametreler
- `updates` (array, zorunlu): Güncelleme listesi

### Güncelleme Formatı
```json
{
  "updates": [
    {
      "id": "task-id-1",
      "updates": {
        "durum": "tamamlandi"
      }
    },
    {
      "id": "task-id-2",
      "updates": {
        "durum": "devam_ediyor"
      }
    }
  ]
}
```

### Desteklenen Güncellemeler
- `durum`: beklemede, devam_ediyor, tamamlandi

### Yanıt
```markdown
## 📦 Toplu Güncelleme Sonucu

**Toplam İşlenen:** 2
**Başarılı:** 2
**Başarısız:** 0

### ✅ Başarılı Güncellemeler
- task-id-1
- task-id-2
```

## 🔍 gorev_nlp_query

Doğal dil sorgularıyla görev arama.

### Parametreler
- `query` (string, zorunlu): Türkçe doğal dil sorgusu

### Desteklenen Sorgu Türleri

#### Zaman Bazlı
- "bugün üzerinde çalıştığım görevler"
- "son oluşturduğum görev"
- "son oluşturulan 5 görev"

#### Durum Bazlı
- "yüksek öncelikli görevler"
- "tamamlanmamış görevler"
- "devam eden görevler"
- "tamamlanan görevler"

#### Özel Filtreler
- "blokaj olan görevler"
- "acil görevler"
- "gecikmiş görevler"

#### Etiket Bazlı
- "etiket:bug"
- "tag:feature"

#### Genel Arama
- Başlık ve açıklamada kelime araması
- Birden fazla kelime AND mantığıyla aranır

### Örnek Kullanım
```json
{
  "query": "yüksek öncelikli bug etiketli görevler"
}
```

### Yanıt
```markdown
## 🔍 Arama Sonuçları: "yüksek öncelikli bug etiketli görevler"

2 görev bulundu:

- **Kritik Login Bug** (Y) | Login sisteminde hata | 🏷️ bug | 48d92f10
- **Database Connection Error** (Y) | Veritabanı bağlantı sorunu | 🏷️ bug, critical | 6ea83c29
```

## Otomatik Durum Yönetimi

### gorev_detay ile Entegrasyon
`gorev_detay` aracı artık görüntülenen görevi otomatik olarak "devam_ediyor" durumuna geçirir:
- Sadece "beklemede" durumundaki görevler etkilenir
- AI interaction kaydı tutulur
- Kullanıcıya bildirim yapılmaz (sessiz güncelleme)

### Durum Geçiş Kuralları
1. **beklemede → devam_ediyor**: Görev görüntülendiğinde veya aktif ayarlandığında
2. **devam_ediyor → tamamlandi**: Manuel güncelleme gerekir
3. Alt görevli görevler tüm alt görevler tamamlanmadan "tamamlandi" yapılamaz

## Kullanım Senaryoları

### Senaryo 1: Günlük Çalışma Akışı
```
AI: "Bugün üzerinde çalışacağım görevleri göster"
→ gorev_nlp_query("bugün")

AI: "İlk görevi aktif yap"
→ gorev_set_active(task_id)

AI: "Aktif görev detaylarını göster"
→ gorev_get_active()
```

### Senaryo 2: Toplu Durum Güncelleme
```
AI: "Test edilen 3 görevi tamamlandı olarak işaretle"
→ gorev_batch_update({
    updates: [
      {id: "1", updates: {durum: "tamamlandi"}},
      {id: "2", updates: {durum: "tamamlandi"}},
      {id: "3", updates: {durum: "tamamlandi"}}
    ]
  })
```

### Senaryo 3: Akıllı Görev Bulma
```
AI: "Acil bug'ları listele"
→ gorev_nlp_query("acil etiket:bug")

AI: "Blokajda olan yüksek öncelikli görevleri bul"
→ gorev_nlp_query("blokaj yüksek öncelik")
```

## Performans ve Limitler

- Recent tasks maksimum 10 görev saklar (FIFO)
- NLP query sonuçları pagination desteklemez (tüm eşleşenler döner)
- Batch update maksimum 100 görev işleyebilir
- Context summary maksimum 5 öncelikli görev ve 5 blokaj gösterir

## Hata Durumları

### gorev_set_active
- "görev bulunamadı": Geçersiz task_id
- "task_id parametresi gerekli": Parametre eksik

### gorev_batch_update
- "updates parametresi gerekli ve dizi olmalı": Yanlış format
- Her başarısız güncelleme için detaylı hata mesajı

### gorev_nlp_query
- "query parametresi gerekli": Boş sorgu
- Eşleşme bulunamazsa: "Eşleşen görev bulunamadı"

## Gelecek Geliştirmeler

1. **Tahmin Sistemi**: estimated_hours ve actual_hours kullanımı
2. **Akıllı Önceliklendirme**: AI'nın görev önceliklerini öğrenmesi
3. **Otomatik Kategorizasyon**: NLP ile otomatik etiketleme
4. **Bağlam Tabanlı Öneriler**: Çalışma alışkanlıklarına göre görev önerileri