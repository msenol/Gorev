# MCP Araçları Referansı - v0.17.0

**24 optimize MCP aracı** için eksiksiz referans (v0.16.0'da 45 araçtan birleştirildi).

**Son Güncelleme:** 22 Kasım 2025 | **Sürüm:** v0.17.0 | **Test Durumu:** ✅ %100 Başarı

---

## 📋 Araç Kategorileri

### TEMEL ARAÇLAR (11)

1. `gorev_listele` - Görevleri listele ve filtrele
2. `gorev_detay` - Detaylı görev bilgisi göster
3. `gorev_guncelle` - Görev durumu/önceliği güncelle ⭐ **v0.16.3 GENİŞLETİLDİ**
4. `gorev_duzenle` - Görev içeriği ve özelliklerini düzenle
5. `gorev_sil` - Güvenlik kontrolleriyle görev sil
6. `gorev_bagimlilik_ekle` - Görev bağımlılıkları oluştur
7. `template_listele` - Kullanılabilir şablonları listele
8. `templateden_gorev_olustur` - Şablondan görev oluştur
9. `proje_olustur` - Yeni proje oluştur
10. `proje_listele` - Projeleri görev sayılarıyla listele
11. `proje_gorevleri` - Bir projenin görevlerini göster

### BİRLEŞİK ARAÇLAR (8)

12. `aktif_proje` - Aktif proje yönetimi (set|get|clear)
13. `gorev_hierarchy` - Görev hiyerarşi işlemleri
14. `gorev_bulk` - Toplu işlemler ⭐ **v0.16.3 DÜZELTİLDİ**
15. `gorev_filter_profile` - Filtre profili yönetimi (create|list|get|update|delete)
16. `gorev_file_watch` - Dosya izleme işlemleri (add|list|remove|get)
17. `gorev_ide_manage` - IDE uzantı yönetimi (detect|status|restart)
18. `gorev_ai_context` - AI context yönetimi (store|retrieve|analyze|clear)
19. `gorev_search` - Gelişmiş arama ve NLP sorguları ⭐ **v0.16.3 DÜZELTİLDİ**

### ÖZEL ARAÇLAR (5)

20. `ozet_goster` - Çalışma alanı özetini göster
21. `gorev_export` - Görevleri dışa aktar (json|csv|markdown)
22. `gorev_import` - Görevleri içe aktar (json|csv)
23. `gorev_intelligent_create` - AI önerileriyle görev oluştur
24. `gorev_nlp_query` - Doğal dil sorguları

---

## 🔍 Hızlı Referans Tablosu

| Araç | Kategori | Tür | v0.16.3 | Kullanım Sıklığı |
|------|----------|-----|---------|------------------|
| `gorev_listele` | Temel | Okuma | - | ⭐⭐⭐⭐⭐ |
| `gorev_detay` | Temel | Okuma | - | ⭐⭐⭐⭐⭐ |
| `gorev_guncelle` | Temel | Yazma | ⭐ Genişletildi | ⭐⭐⭐⭐⭐ |
| `gorev_duzenle` | Temel | Yazma | - | ⭐⭐⭐⭐ |
| `gorev_sil` | Temel | Yazma | - | ⭐⭐⭐ |
| `gorev_bagimlilik_ekle` | Temel | Yazma | - | ⭐⭐⭐ |
| `template_listele` | Temel | Okuma | - | ⭐⭐⭐⭐ |
| `templateden_gorev_olustur` | Temel | Yazma | - | ⭐⭐⭐⭐ |
| `proje_olustur` | Temel | Yazma | - | ⭐⭐⭐ |
| `proje_listele` | Temel | Okuma | - | ⭐⭐⭐⭐ |
| `proje_gorevleri` | Temel | Okuma | - | ⭐⭐⭐⭐ |
| `aktif_proje` | Birleşik | Karma | - | ⭐⭐⭐⭐ |
| `gorev_hierarchy` | Birleşik | Karma | - | ⭐⭐⭐ |
| `gorev_bulk` | Birleşik | Yazma | ⭐ Düzeltildi | ⭐⭐⭐⭐ |
| `gorev_filter_profile` | Birleşik | Karma | - | ⭐⭐⭐ |
| `gorev_file_watch` | Birleşik | Karma | - | ⭐⭐ |
| `gorev_ide_manage` | Birleşik | Karma | - | ⭐⭐ |
| `gorev_ai_context` | Birleşik | Karma | - | ⭐⭐⭐ |
| `gorev_search` | Birleşik | Okuma | ⭐ Düzeltildi | ⭐⭐⭐⭐⭐ |
| `ozet_goster` | Özel | Okuma | - | ⭐⭐⭐⭐ |
| `gorev_export` | Özel | Okuma | - | ⭐⭐⭐ |
| `gorev_import` | Özel | Yazma | - | ⭐⭐ |
| `gorev_intelligent_create` | Özel | Yazma | - | ⭐⭐⭐⭐ |
| `gorev_nlp_query` | Özel | Okuma | - | ⭐⭐⭐⭐ |

---

## 🔧 Detaylı Araç Spesifikasyonları

### 1. gorev_listele

Görevleri filtrele ve listele. Hiyerarşik ağaç yapısında çıktı verir.

**Parametreler:**

- `status` (string, opsiyonel): Durum filtresi
  - Değerler: `beklemede`, `devam_ediyor`, `tamamlandi`
- `tum_projeler` (boolean, opsiyonel): `true` ise tüm projeler, `false`/boş ise sadece aktif proje
- `sirala` (string, opsiyonel): Sıralama düzeni
  - `son_tarih_asc`: Son tarihe göre artan
  - `son_tarih_desc`: Son tarihe göre azalan
- `filtre` (string, opsiyonel): Özel filtreler
  - `acil`: 7 gün içinde teslim edilecek görevler
  - `gecmis`: Süresi geçmiş görevler
- `etiket` (string, opsiyonel): Etiket adına göre filtrele
- `limit` (number, opsiyonel): Gösterilecek maksimum görev sayısı (varsayılan: 50)
- `offset` (number, opsiyonel): Pagination için atlama sayısı (varsayılan: 0)

**Çıktı Formatı:** Hiyerarşik ağaç yapısı (üst görevler → alt görevler), renk kodlu durum göstergeleri

**Örnek:**

```json
{
  "status": "devam_ediyor",
  "tum_projeler": false,
  "sirala": "son_tarih_asc",
  "limit": 20
}
```

---

### 2. gorev_detay

Bir görevin tüm detaylarını markdown formatında gösterir.

**Parametreler:**

- `id` (string, zorunlu): Görev ID'si (UUID formatı)

**Çıktı:** Son tarihler, etiketler, bağımlılıklar, alt görevler, AI notları, dosya izlemeleri dahil kapsamlı detaylar

**Örnek:**

```json
{
  "id": "abc12345-67de-89fg-hijk-lmnopqrstuvw"
}
```

---

### 3. gorev_guncelle ⭐ **v0.16.3 GÜNCELLENDİ**

Görev durumunu ve/veya önceliğini günceller.

**v0.16.3'te Yenilikler:**

- Artık hem `status` hem de `priority` parametrelerini destekliyor (önceden sadece durum)
- Durumu, önceliği veya her ikisini birden güncelleyebilir
- En az bir parametre gereklidir

**Parametreler:**

- `id` (string, zorunlu): Görev ID'si
- `status` (string, opsiyonel): Yeni durum
  - `beklemede`: Görev henüz başlamadı
  - `devam_ediyor`: Aktif çalışılıyor
  - `tamamlandi`: Tamamlandı
- `priority` (string, opsiyonel): Yeni öncelik ⭐ **YENİ**
  - `dusuk`: Düşük öncelik
  - `orta`: Normal öncelik
  - `yuksek`: Yüksek öncelik

**Durum Geçiş Kuralları:**

- `devam_ediyor`: Tüm bağımlılıklar tamamlanmış olmalı
- `tamamlandi`: Tüm alt görevler tamamlanmış olmalı

**Örnekler:**

```json
// Sadece durum güncelle
{
  "id": "abc12345",
  "status": "devam_ediyor"
}

// Sadece öncelik güncelle (YENİ v0.16.3'te)
{
  "id": "abc12345",
  "priority": "yuksek"
}

// Her ikisini birden güncelle (YENİ v0.16.3'te)
{
  "id": "abc12345",
  "status": "devam_ediyor",
  "priority": "yuksek"
}
```

---

### 4. gorev_duzenle

Görev içeriğini ve özelliklerini düzenler.

**Parametreler:**

- `id` (string, zorunlu): Görev ID'si
- `title` (string, opsiyonel): Yeni başlık
- `description` (string, opsiyonel): Yeni açıklama (markdown destekli)
- `priority` (string, opsiyonel): Yeni öncelik (dusuk, orta, yuksek)
- `proje_id` (string, opsiyonel): Yeni proje ID'si (alt görevler de taşınır)
- `son_tarih` (string, opsiyonel): Yeni son tarih (YYYY-AA-GG formatı)

**Not:** Herhangi bir parametre değiştirilebilir, diğerleri aynı kalır.

**Örnek:**

```json
{
  "id": "abc12345",
  "title": "API dokümantasyonu güncelle",
  "priority": "yuksek",
  "son_tarih": "2025-10-15"
}
```

---

### 5. gorev_sil

Güvenlik kontrolleriyle görev siler.

**Parametreler:**

- `id` (string, zorunlu): Görev ID'si
- `onay` (boolean, zorunlu): Silme onayı (güvenlik için)

**Güvenlik Kontrolleri:**

- Alt görevleri olan görevler silinemez (önce alt görevleri silin)
- Onay parametresi `true` olmalı

**Örnek:**

```json
{
  "id": "abc12345",
  "onay": true
}
```

---

### 6. gorev_bagimlilik_ekle

İki görev arasında bağımlılık oluşturur.

**Parametreler:**

- `kaynak_id` (string, zorunlu): Kaynak görev ID'si
- `hedef_id` (string, zorunlu): Hedef görev ID'si
- `baglanti_tipi` (string, zorunlu): Bağlantı tipi
  - `engelliyor`: Kaynak görev tamamlanmadan hedef görev başlayamaz
  - `iliskili`: İlişkili görevler (bilgi amaçlı)

**Dairesel Bağımlılık Kontrolü:** Sistem otomatik olarak dairesel bağımlılıkları engeller

**Örnek:**

```json
{
  "kaynak_id": "abc12345",
  "hedef_id": "def67890",
  "baglanti_tipi": "engelliyor"
}
```

---

### 7. template_listele

Kullanılabilir görev şablonlarını listeler.

**Parametreler:**

- `kategori` (string, opsiyonel): Kategori filtresi
  - Örnekler: `Teknik`, `Özellik`, `Araştırma`, `Refaktör`, `Test`, `Dokümantasyon`

**Çıktı:** Şablon ID'leri, isimleri, kategorileri ve kullanılabilir alanlar

**Şablon Takma Adları (v0.16.2+):** `bug`, `feature`, `research`, `refactor`, `test`, `doc`

**Örnek:**

```json
{
  "kategori": "Teknik"
}
```

---

### 8. templateden_gorev_olustur

Seçilen şablondan özelleştirilmiş görev oluşturur.

**Parametreler:**

- `template_id` (string, zorunlu): Şablon ID'si veya takma adı
- `values` (object, zorunlu): Şablon alanları için key-value çiftleri

**Şablon Takma Adları Kullanımı:**

```json
// UUID yerine takma ad kullanın
{
  "template_id": "bug",
  "values": {
    "hatanin_adi": "Kullanıcı girişi başarısız",
    "sertlik": "kritik",
    "adimlar": "1. Giriş sayfasını aç\n2. Yanlış şifre gir\n3. 500 hatası al"
  }
}
```

---

### 9. proje_olustur

Yeni proje oluşturur.

**Parametreler:**

- `name` (string, zorunlu): Proje adı
- `definition` (string, zorunlu): Proje açıklaması

**Örnek:**

```json
{
  "name": "Web Sitesi Yenileme",
  "definition": "Kurumsal web sitesini modern teknolojilerle yeniden inşa etme projesi"
}
```

---

### 10. proje_listele

Tüm projeleri görev sayılarıyla listeler.

**Parametreler:** Yok

**Çıktı:** Proje ID, isim, açıklama ve görev sayıları (toplam, beklemede, devam ediyor, tamamlandı)

---

### 11. proje_gorevleri

Bir projenin görevlerini duruma göre gruplandırarak listeler.

**Parametreler:**

- `proje_id` (string, zorunlu): Proje ID'si
- `limit` (number, opsiyonel): Gösterilecek maksimum görev sayısı (varsayılan: 50)
- `offset` (number, opsiyonel): Pagination için atlama sayısı (varsayılan: 0)

**Örnek:**

```json
{
  "proje_id": "xyz78901",
  "limit": 30
}
```

---

### 12. aktif_proje

Aktif proje yönetimi (action tabanlı birleşik araç).

**Parametreler:**

- `action` (string, zorunlu): İşlem tipi
  - `set`: Aktif projeyi ayarla
  - `get`: Aktif projeyi göster
  - `clear`: Aktif proje ayarını kaldır
- `proje_id` (string, action=set için zorunlu): Proje ID'si

**Örnekler:**

```json
// Aktif projeyi ayarla
{
  "action": "set",
  "proje_id": "xyz78901"
}

// Aktif projeyi göster
{
  "action": "get"
}

// Aktif proje ayarını kaldır
{
  "action": "clear"
}
```

---

### 13. gorev_hierarchy

Görev hiyerarşi işlemleri (action tabanlı birleşik araç).

**Parametreler:**

- `action` (string, zorunlu): İşlem tipi
  - `show`: Hiyerarşiyi göster
  - `move`: Görevi taşı
  - `create_subtask`: Alt görev oluştur
- `gorev_id` (string, zorunlu): Hedef görev ID'si
- `yeni_parent_id` (string, action=move için): Yeni üst görev ID'si (boş string = kök)
- `title`, `description`, `priority`, `son_tarih`, `tags` (action=create_subtask için)

**Örnekler:**

```json
// Hiyerarşiyi göster
{
  "action": "show",
  "gorev_id": "abc12345"
}

// Görevi taşı
{
  "action": "move",
  "gorev_id": "abc12345",
  "yeni_parent_id": "def67890"
}

// Alt görev oluştur
{
  "action": "create_subtask",
  "gorev_id": "abc12345",
  "title": "Alt görev başlığı",
  "priority": "yuksek"
}
```

---

### 14. gorev_bulk ⭐ **v0.16.3 DÜZELTİLDİ**

Toplu işlemler aracı (update, transition, tag operations).

**v0.16.3'te Yenilikler:**

- Tüm 3 işlem için parametre dönüşümüyle kapsamlı yeniden yazım
- %100 başarı oranı onaylandı
- Geriye dönük uyumlu parametre formatları

**Parametreler:**

- `operation` (string, zorunlu): İşlem tipi
  - `update`: Toplu güncelleme
  - `transition`: Toplu durum geçişi
  - `tag`: Toplu etiket işlemleri

---

#### İşlem: update ⭐ **v0.16.3'te DÜZELTİLDİ**

Birden fazla görevi aynı değişikliklerle günceller.

**Parametreler:**

- `operation`: `"update"`
- `ids` (array, zorunlu): Görev ID'leri dizisi
- `data` (object, zorunlu): Uygulanacak değişiklikler
  - `status` (string, opsiyonel): Yeni durum
  - `priority` (string, opsiyonel): Yeni öncelik
  - `proje_id` (string, opsiyonel): Yeni proje ID'si
  - `son_tarih` (string, opsiyonel): Yeni son tarih

**İç Dönüşüm:** `{ids: [], data: {}}` → `{updates: [{id, ...fields}]}`

**Örnek:**

```json
{
  "operation": "update",
  "ids": ["abc123", "def456", "ghi789"],
  "data": {
    "priority": "yuksek",
    "status": "devam_ediyor"
  }
}
```

**Çıktı:** İşlenen görev sayısı + başarı/başarısızlık detayları

---

#### İşlem: transition ⭐ **v0.16.3'te DÜZELTİLDİ**

Birden fazla görevi yeni duruma geçir.

**Parametreler:**

- `operation`: `"transition"`
- `ids` (array, zorunlu): Görev ID'leri dizisi
- `status` veya `yeni_durum` (string, zorunlu): Hedef durum ⭐ **Her iki parametre adı da destekleniyor**
- `force` (boolean, opsiyonel): Bağımlılık kontrollerini atla (varsayılan: false)
- `check_dependencies` (boolean, opsiyonel): Bağımlılıkları kontrol et (varsayılan: true)
- `dry_run` (boolean, opsiyonel): Deneme modu (değişiklik yapmadan simüle et)

**Örnek:**

```json
{
  "operation": "transition",
  "ids": ["abc123", "def456"],
  "status": "tamamlandi",
  "check_dependencies": true
}
```

---

#### İşlem: tag ⭐ **v0.16.3'te DÜZELTİLDİ**

Birden fazla görevde toplu etiket işlemleri.

**Parametreler:**

- `operation`: `"tag"`
- `ids` (array, zorunlu): Görev ID'leri dizisi
- `operation` veya `tag_operation` (string, zorunlu): Etiket işlemi ⭐ **Her iki parametre adı da destekleniyor**
  - `add`: Etiketleri ekle
  - `remove`: Etiketleri kaldır
  - `replace`: Etiketleri değiştir
- `tags` (array, zorunlu): Etiket isimleri dizisi

**Örnekler:**

```json
// Etiket ekle
{
  "operation": "tag",
  "ids": ["abc123", "def456"],
  "tag_operation": "add",
  "tags": ["acil", "backend"]
}

// Etiketleri değiştir
{
  "operation": "tag",
  "ids": ["abc123"],
  "operation": "replace",
  "tags": ["tamamlandi", "dağıtıldı"]
}
```

**Performans:** 11-33ms işlem süresi (v0.16.3 test sonuçları)

---

### 15. gorev_filter_profile

Filtre profili yönetimi (action tabanlı birleşik araç).

**Parametreler:**

- `action` (string, zorunlu): İşlem tipi
  - `create`: Yeni profil oluştur
  - `list`: Tüm profilleri listele
  - `get`: Profil detaylarını getir
  - `update`: Profili güncelle
  - `delete`: Profili sil
- `id` (string, action=get|update|delete için zorunlu): Profil ID'si
- `name` (string, action=create|update için): Profil adı
- `description` (string, opsiyonel): Profil açıklaması
- `filters` (object, action=create|update için): Filtre konfigürasyonu

**Örnek:**

```json
{
  "action": "create",
  "name": "Acil Backend Görevleri",
  "description": "Yüksek öncelikli backend görevleri",
  "filters": {
    "status": ["beklemede", "devam_ediyor"],
    "priority": ["yuksek"],
    "tags": ["backend", "acil"]
  }
}
```

---

### 16. gorev_file_watch

Dosya izleme işlemleri (action tabanlı birleşik araç).

**Parametreler:**

- `action` (string, zorunlu): İşlem tipi
  - `add`: Dosya/dizin izleme ekle
  - `list`: Tüm izlemeleri listele
  - `remove`: İzlemeyi kaldır
  - `get`: İzleme detaylarını getir
- `gorev_id` (string, action=add için zorunlu): Görev ID'si
- `watch_id` (string, action=remove|get için zorunlu): İzleme ID'si
- `file_path` (string, action=add için zorunlu): İzlenecek dosya/dizin yolu
- `watch_type` (string, action=add için): İzleme tipi (file, directory)

**Örnek:**

```json
{
  "action": "add",
  "gorev_id": "abc12345",
  "file_path": "/path/to/project/src",
  "watch_type": "directory"
}
```

---

### 17. gorev_ide_manage

IDE uzantı yönetimi (action tabanlı birleşik araç).

**Parametreler:**

- `action` (string, zorunlu): İşlem tipi
  - `detect`: IDE'leri tespit et
  - `status`: IDE uzantı durumunu göster
  - `restart`: Daemon/uzantıyı yeniden başlat

**Örnekler:**

```json
// IDE'leri tespit et
{
  "action": "detect"
}

// Daemon'u yeniden başlat
{
  "action": "restart"
}
```

**Çıktı:** VS Code, JetBrains, Cursor, Windsurf gibi yüklü IDE'lerin listesi ve durumları

---

### 18. gorev_ai_context

AI context yönetimi (action tabanlı birleşik araç).

**Parametreler:**

- `action` (string, zorunlu): İşlem tipi
  - `store`: Context bilgisini sakla
  - `retrieve`: Context bilgisini getir
  - `analyze`: Context'i analiz et
  - `clear`: Context'i temizle
- `gorev_id` (string, zorunlu): Görev ID'si
- `context_data` (object, action=store için): Saklanacak context verisi
- `context_type` (string, opsiyonel): Context tipi (conversation, code_review, analysis)

**Örnek:**

```json
{
  "action": "store",
  "gorev_id": "abc12345",
  "context_data": {
    "type": "code_review",
    "files": ["src/api/handler.go"],
    "summary": "REST API handler için code review tamamlandı"
  }
}
```

---

### 19. gorev_search ⭐ **v0.16.3 DÜZELTİLDİ**

Gelişmiş arama ve NLP sorguları (action tabanlı birleşik araç).

**v0.16.3'te Yenilikler:**

- Gelişmiş mod artık sorgu dizelerini `"status:X priority:Y"` formatında filtrelere ayrıştırıyor
- Boşlukla ayrılmış key:value çiftleri otomatik olarak çıkarılıyor
- Mevcut filtre parametresiyle sorunsuz çalışıyor

**Parametreler:**

- `action` (string, opsiyonel): İşlem tipi (varsayılan: "advanced")
  - `advanced`: Gelişmiş arama (FTS5 + filtreler) ⭐ **v0.16.3'te geliştirildi**
  - `nlp`: Doğal dil sorgusu
  - `suggestions`: Akıllı öneriler
  - `history`: Arama geçmişi
- `query` (string, gelişmiş/nlp/öneriler için zorunlu): Arama terimi veya sorgu
- `filters` (object, gelişmiş arama için opsiyonel): Filtre konfigürasyonu
- `limit` (number, opsiyonel): Sonuç sayısı (varsayılan: 50)

---

#### Action: advanced ⭐ **v0.16.3'te GELİŞTİRİLDİ**

FTS5 tam metin arama + gelişmiş filtreler.

**Sorgu Dizesi Ayrıştırma (YENİ v0.16.3):**

```json
// Sorgu dizesindeki key:value çiftleri otomatik olarak filtrelere dönüştürülür
{
  "action": "advanced",
  "query": "status:devam_ediyor priority:yuksek API dokümantasyonu"
}
// İç ayrıştırma: status ve priority filtrelere çıkarılır, "API dokümantasyonu" FTS5 aramasında kullanılır
```

**Manuel Filtreler:**

```json
{
  "action": "advanced",
  "query": "veritabanı",
  "filters": {
    "status": ["beklemede", "devam_ediyor"],
    "priority": ["yuksek"],
    "tags": ["backend"],
    "olusturulma_baslangic": "2025-10-01",
    "son_tarih_bitis": "2025-10-31"
  },
  "limit": 20
}
```

**Bulanık Arama (Fuzzy Matching):**

```json
{
  "action": "advanced",
  "query": "databas",  // "database" bul
  "enable_fuzzy": true,
  "fuzzy_threshold": 2
}
```

**Performans:** 6-67ms (FTS5 tam metin arama + relevance scoring, v0.16.3 test sonuçları)

---

#### Action: nlp

Doğal dil sorgularını yapılandırılmış filtrelere dönüştürür.

**Örnek:**

```json
{
  "action": "nlp",
  "query": "bu hafta teslim edilecek yüksek öncelikli görevleri göster"
}
```

**Çıktı:** NLP motoru sorguyu ayrıştırır ve uygun filtreleri uygular

---

#### Action: suggestions

Arama terimi için akıllı öneriler üretir.

**Örnek:**

```json
{
  "action": "suggestions",
  "query": "API"
}
```

**Çıktı:** NLP tabanlı ve geçmiş aramalara dayalı öneri listesi

---

#### Action: history

Arama geçmişini getirir.

**Örnek:**

```json
{
  "action": "history",
  "limit": 10
}
```

---

### 20. ozet_goster

Çalışma alanı özetini gösterir.

**Parametreler:** Yok

**Çıktı:**

- Aktif proje bilgisi
- Toplam proje sayısı
- Toplam görev sayısı
- Durum bazlı görev dağılımı
- Öncelik bazlı görev dağılımı
- Yaklaşan son tarihler
- Son aktivite özeti

---

### 21. gorev_export

Görevleri çeşitli formatlarda dışa aktarır.

**Parametreler:**

- `format` (string, zorunlu): Çıktı formatı
  - `json`: JSON formatı
  - `csv`: CSV formatı
  - `markdown`: Markdown formatı
- `proje_id` (string, opsiyonel): Proje ID'si (belirtilmezse tüm projeler)
- `status` (string, opsiyonel): Durum filtresi
- `include_subtasks` (boolean, opsiyonel): Alt görevleri dahil et (varsayılan: true)

**Örnek:**

```json
{
  "format": "markdown",
  "proje_id": "xyz78901",
  "status": "tamamlandi",
  "include_subtasks": true
}
```

---

### 22. gorev_import

Görevleri içe aktarır.

**Parametreler:**

- `format` (string, zorunlu): Girdi formatı
  - `json`: JSON formatı
  - `csv`: CSV formatı
- `data` (string, zorunlu): İçe aktarılacak veri (JSON dizesi veya CSV içeriği)
- `proje_id` (string, opsiyonel): Hedef proje ID'si

**Örnek:**

```json
{
  "format": "json",
  "data": "[{\"title\":\"New Task\",\"priority\":\"high\"}]",
  "proje_id": "xyz78901"
}
```

---

### 23. gorev_intelligent_create

AI önerileriyle görev oluşturur.

**Parametreler:**

- `title` (string, zorunlu): Görev başlığı
- `description` (string, opsiyonel): Görev açıklaması
- `enable_ai_suggestions` (boolean, opsiyonel): AI önerilerini etkinleştir (varsayılan: true)

**AI Önerileri:**

- Otomatik öncelik tahmini
- Önerilen etiketler
- Benzer görev tespiti
- Önerilen şablon
- Tahmini tamamlanma süresi

**Örnek:**

```json
{
  "title": "Kullanıcı kimlik doğrulama API'si",
  "description": "JWT tabanlı kimlik doğrulama endpoint'leri oluştur",
  "enable_ai_suggestions": true
}
```

**Çıktı:**

```json
{
  "gorev_id": "abc12345",
  "ai_suggestions": {
    "priority": "yuksek",
    "tags": ["backend", "güvenlik", "API"],
    "sablonlar": ["feature"],
    "benzer_gorevler": ["xyz78901"],
    "tahmini_sure": "3-5 gün"
  }
}
```

---

### 24. gorev_nlp_query

Doğal dil sorgularını işler ve çalıştırır.

**Parametreler:**

- `query` (string, zorunlu): Doğal dil sorgusu

**Desteklenen Sorgu Türleri:**

- Görev listeleme: "bu hafta teslim edilecek görevler"
- Görev oluşturma: "API dokümantasyonu için yüksek öncelikli görev oluştur"
- Durum güncelleme: "görev abc12345'i tamamlandı olarak işaretle"
- Arama: "veritabanı ile ilgili görevleri bul"
- İstatistik: "kaç tane yüksek öncelikli görev var"

**Örnek:**

```json
{
  "query": "gelecek hafta teslim edilecek yüksek öncelikli görevleri göster"
}
```

**Çıktı:** Sorgu ayrıştırılır, uygun MCP aracı çağrılır ve sonuçlar döndürülür

---

## 📊 Sürüm Geçmişi

### v0.16.3 (6 Ekim 2025) - Kritik Düzeltmeler

**Düzeltilen Araçlar:**

1. **gorev_bulk** ⭐ - Kapsamlı yeniden yazım
   - `update` işlemi: `{ids: [], data: {}}` → `{updates: [{id, ...fields}]}` dönüşümü
   - `transition` işlemi: Hem `status` hem de `yeni_durum` parametrelerini kabul eder
   - `tag` işlemi: Hem `operation` hem de `tag_operation` parametrelerini kabul eder
   - Test sonucu: %100 başarı (2/2 update, 1/1 transition, 2/2 tag)

2. **gorev_guncelle** ⭐ - Öncelik desteği eklendi
   - Artık hem `status` hem de `priority` güncellemelerini destekliyor
   - Her ikisi de birlikte veya ayrı ayrı güncellenebilir
   - En az bir parametre gerekli
   - Test sonucu: %100 başarı (durum ✓, öncelik ✓, her ikisi ✓)

3. **gorev_search** ⭐ - Sorgu ayrıştırma geliştirildi
   - Gelişmiş mod artık `"status:X priority:Y"` gibi sorguları ayrıştırıyor
   - `parseQueryFilters()` yardımcı fonksiyonu eklendi
   - Boşlukla ayrılmış key:value çiftleri otomatik olarak çıkarılıyor
   - Test sonucu: %100 başarı (tek filtre 8 sonuç, çoklu filtre 4 sonuç, etiket 21 sonuç)

4. **VS Code Tree View** ⭐ - Bağımlılık sayaçları düzeltildi
   - `omitempty` JSON tag'leri kaldırıldı (bagimli_gorev_sayisi, tamamlanmamis_bagimlilik_sayisi, bu_goreve_bagimli_sayisi)
   - Bağımlılık göstergeleri (🔒/🔓/🔗) artık doğru görüntüleniyor

**Performans Metrikleri:**

- Toplu işlemler: 11-33ms işlem süresi
- Gelişmiş arama: 6-67ms (FTS5 tam metin arama + relevance scoring)
- Başarı oranı: %100 (Kilocode AI kapsamlı test raporu tarafından doğrulandı)

### v0.16.2 (5 Ekim 2025)

- NPM binary güncelleme mekanizması (78.4 MB → 6.9 KB paket boyutu)
- Gömülü Web UI (React + TypeScript) - http://localhost:5082
- Çoklu çalışma alanı desteği (SHA256 tabanlı workspace ID'leri)
- Şablon takma adları: `bug`, `feature`, `research`, `refactor`, `test`, `doc`

### v0.16.0 (3 Ekim 2025) - Araç Birleştirme

- **45 araçtan 24 optimize araca** indirgeme
- 8 birleşik handler: Aktif Proje, Hiyerarşi, Toplu İşlemler, Filtre Profilleri, Dosya İzleme, IDE Yönetimi, AI Context, Arama
- Action tabanlı routing ile gelişmiş bakım kolaylığı

---

## 🚀 Performans ve En İyi Uygulamalar

### Performans Metrikleri (v0.16.3)

| Araç | Ortalama Süre | Kullanım Senaryosu |
|------|---------------|-------------------|
| `gorev_listele` | 5-15ms | Küçük/orta çalışma alanları (<1000 görev) |
| `gorev_detay` | 3-8ms | Tek görev sorgusu |
| `gorev_search` (advanced) | 6-67ms | FTS5 tam metin arama + filtreler |
| `gorev_bulk` (update) | 11-33ms | 2-10 görev toplu güncelleme |
| `gorev_bulk` (transition) | 15-40ms | Bağımlılık kontrollerini içerir |
| `gorev_nlp_query` | 50-150ms | NLP ayrıştırma + araç çağrısı |

### En İyi Uygulamalar

1. **Toplu İşlemler için `gorev_bulk` Kullanın**
   - Tek tek güncelleme yerine toplu güncelleme
   - Her görev için ~3-5ms tasarruf
   - Örnek: 10 görev → 30-50ms tasarruf

2. **FTS5 Tam Metin Arama**
   - Tam eşleşme yerine FTS5 kullanın
   - `enable_fuzzy: true` ile typo toleransı
   - `fuzzy_threshold: 2` dengeli accuracy için

3. **Filtre Profillerini Kaydedin**
   - Sık kullanılan filtreleri profillere kaydedin
   - Manuel filtre parametrelerini tekrar yazmayın
   - Paylaşılabilir ve yeniden kullanılabilir

4. **NLP Sorgularını Akıllıca Kullanın**
   - Basit sorgular için gelişmiş aramayı tercih edin
   - Karmaşık doğal dil istekleri için NLP kullanın
   - NLP ~50-100ms ek yük getirir

5. **Pagination Uygulayın**
   - Büyük sonuç kümeleri için `limit` ve `offset` kullanın
   - Varsayılan limit: 50 görev
   - Önerilen: Sayfa başına 20-50 görev

6. **Hiyerarşik Yapı**
   - Karmaşık projelerde alt görevleri kullanın
   - Derinlik: 3-5 seviye önerilen maksimum
   - Performans: Hiyerarşi derinliği başına ~2-3ms

7. **Bağımlılık Yönetimi**
   - Dairesel bağımlılıklardan kaçının (sistem kontrol eder)
   - `force: true` dikkatli kullanın
   - Bağımlılık kontrolü: Her geçişte ~5-10ms

---

## 🔧 Sorun Giderme

### Yaygın Hatalar

#### 1. "Görev bulunamadı"

**Neden:** Yanlış görev ID'si veya görev silinmiş
**Çözüm:** `gorev_listele` ile ID'yi doğrulayın

#### 2. "Bağımlılıklar tamamlanmamış"

**Neden:** `devam_ediyor` durumuna geçerken bağımlı görevler henüz tamamlanmamış
**Çözüm:**

- `gorev_detay` ile bağımlılıkları kontrol edin
- Önce bağımlı görevleri tamamlayın
- Veya `force: true` kullanın (dikkatli!)

#### 3. "Alt görevler tamamlanmamış"

**Neden:** `tamamlandi` durumuna geçerken alt görevler henüz tamamlanmamış
**Çözüm:**

- `gorev_hierarchy` ile alt görevleri kontrol edin
- Önce tüm alt görevleri tamamlayın

#### 4. "Geçersiz durum geçişi"

**Neden:** İzin verilmeyen durum değişikliği
**Çözüm:** Geçerli durumlar: `beklemede`, `devam_ediyor`, `tamamlandi`

#### 5. "Dairesel bağımlılık tespit edildi"

**Neden:** Görev A → B → C → A gibi dairesel bir zincir oluşturuluyor
**Çözüm:** Bağımlılık yapısını yeniden tasarlayın

#### 6. "Alt görevleri olan görev silinemez"

**Neden:** Güvenlik önlemi - veri kaybını önler
**Çözüm:** Önce tüm alt görevleri silin

### Debug İpuçları

1. **Verbose Logging:**

   ```bash
   gorev serve --debug
   ```

2. **API Health Check:**

   ```bash
   curl http://localhost:5082/api/health
   ```

3. **Görev ID Doğrulama:**

   ```json
   {"action": "advanced", "query": "görev_id"}
   ```

4. **Bağımlılık Grafiği:**

   ```json
   {"action": "show", "gorev_id": "abc12345"}
   ```

---

## 🔌 API Entegrasyonu

### Daemon Port Yapılandırması

```bash
# Varsayılan port: 5082
gorev serve

# Özel port
gorev serve --api-port 8080

# Daemon modu (arka planda)
gorev daemon --detach
```

### Lock File Konumu

- **Linux/macOS:** `~/.gorev-daemon/.lock`
- **Windows:** `%USERPROFILE%\.gorev-daemon\.lock`

**Lock File İçeriği:**

```json
{
  "pid": 12345,
  "port": "5082",
  "start_time": "2025-10-06T10:30:00Z",
  "daemon_url": "http://localhost:5082",
  "version": "0.16.3"
}
```

### WebSocket Endpoint (Deneysel)

```javascript
const ws = new WebSocket('ws://localhost:5082/ws');
ws.onmessage = (event) => {
  const update = JSON.parse(event.data);
  console.log('Görev güncellendi:', update);
};
```

### REST API Endpoints (23 Endpoint)

- `GET /api/health` - Daemon sağlık kontrolü
- `GET /api/gorevler` - Görev listesi
- `POST /api/gorev` - Yeni görev
- `GET /api/gorev/:id` - Görev detayı
- `PUT /api/gorev/:id` - Görev güncelleme
- `DELETE /api/gorev/:id` - Görev silme
- [21 endpoint daha...]

**Not:** REST API öncelikle VS Code uzantısı için tasarlanmıştır. MCP istemcileri için MCP araçlarını kullanın.

---

## 📝 Notlar

- Tüm araçlar Turkish domain language kullanır (gorev, proje, durum, vb.)
- Görev açıklamaları tam markdown formatını destekler
- Tarih formatı: `YYYY-AA-GG` (örn: `2025-10-30`)
- ID'ler UUID v4 formatındadır
- FTS5 tam metin arama SQLite extension gerektirir (genellikle dahildir)
- Bulanık arama Levenshtein distance algoritması kullanır
- Arama geçmişi otomatik olarak kaydedilir
- Filtre profilleri JSON formatında saklanır
- NLP özellikleri AI Context Management sistemi ile entegredir

---

## 🔗 İlgili Kaynaklar

- [İngilizce MCP Araçları Referansı](../../docs/api/MCP_TOOLS_REFERENCE.md)
- [Daemon Mimarisi Dokümantasyonu](../../docs/architecture/daemon-architecture.md)
- [v0.16.3 Release Notes](../../docs/releases/RELEASE_NOTES_v0.16.3.md)
- [CHANGELOG](../CHANGELOG.md)
- [NPM Package](https://www.npmjs.com/package/@mehmetsenol/gorev-mcp-server)
- [GitHub Repository](https://github.com/msenol/Gorev)

---

**Son Güncelleme:** 6 Ekim 2025 | **Doğrulayan:** Kilocode AI Test Raporu | **Durum:** Üretim Hazır ✅
