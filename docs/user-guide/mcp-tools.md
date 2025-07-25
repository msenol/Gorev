# MCP Araçları Referansı

Gorev'in sağladığı 25 MCP tool'unun detaylı açıklaması.

> ⚠️ **BREAKING CHANGE (v0.10.0)**: `gorev_olustur` artık kullanılamaz! Template kullanımı zorunludur. Detaylar için [templateden_gorev_olustur](#templateden_gorev_olustur) bölümüne bakın.

> **Not**: AI Context Management araçları için [AI MCP Araçları Dokümantasyonu](../mcp-araclari-ai.md)'na bakın.

## 📋 Araç Listesi

### Görev Yönetimi
1. [gorev_olustur](#gorev_olustur) - ⚠️ **DEPRECATED (v0.10.0)** - Template kullanımı zorunlu
2. [gorev_listele](#gorev_listele) - Görevleri listeleme
3. [gorev_detay](#gorev_detay) - Görev detaylarını görüntüleme (markdown)
4. [gorev_guncelle](#gorev_guncelle) - Görev durumu güncelleme
5. [gorev_duzenle](#gorev_duzenle) - Görev bilgilerini düzenleme
6. [gorev_sil](#gorev_sil) - Görev silme
7. [gorev_bagimlilik_ekle](#gorev_bagimlilik_ekle) - Görevler arası bağımlılık oluşturma

### Subtask Yönetimi (v0.8.0+)
8. [gorev_altgorev_olustur](#gorev_altgorev_olustur) - Alt görev oluşturma
9. [gorev_ust_degistir](#gorev_ust_degistir) - Görevin üst görevini değiştirme
10. [gorev_hiyerarsi_goster](#gorev_hiyerarsi_goster) - Görev hiyerarşisini gösterme

### Görev Şablonları
11. [template_listele](#template_listele) - Görev şablonlarını listeleme
12. [templateden_gorev_olustur](#templateden_gorev_olustur) - Şablondan görev oluşturma

### Proje Yönetimi
13. [proje_olustur](#proje_olustur) - Yeni proje oluşturma
14. [proje_listele](#proje_listele) - Tüm projeleri listeleme
15. [proje_gorevleri](#proje_gorevleri) - Bir projenin görevlerini listeleme
16. [proje_aktif_yap](#proje_aktif_yap) - Projeyi aktif olarak ayarlama
17. [aktif_proje_goster](#aktif_proje_goster) - Aktif projeyi görüntüleme
18. [aktif_proje_kaldir](#aktif_proje_kaldir) - Aktif proje ayarını kaldırma

### Raporlama
19. [ozet_goster](#ozet_goster) - Sistem özeti görüntüleme

### AI Context Management (v0.9.0+)
20. [gorev_set_active](#ai-araçları) - Aktif görevi ayarlama ve otomatik durum yönetimi
21. [gorev_get_active](#ai-araçları) - Aktif görevi görüntüleme
22. [gorev_recent](#ai-araçları) - Son görüntülenen görevleri listeleme
23. [gorev_context_summary](#ai-araçları) - AI oturum özeti
24. [gorev_batch_update](#ai-araçları) - Toplu görev güncelleme
25. [gorev_nlp_query](#ai-araçları) - Doğal dil ile görev arama

> **Detaylı bilgi için**: [AI MCP Araçları Dokümantasyonu](../mcp-araclari-ai.md)

---

## gorev_olustur

⚠️ **DEPRECATED (v0.10.0)**: Bu araç artık kullanılamaz! Template kullanımı zorunludur.

### Migration
`gorev_olustur` yerine artık [templateden_gorev_olustur](#templateden_gorev_olustur) kullanılmalıdır.

**Eski Kullanım:**
```json
{
  "name": "gorev_olustur",
  "arguments": {
    "baslik": "Bug fix",
    "aciklama": "Login sorunu",
    "oncelik": "yuksek"
  }
}
```

**Yeni Kullanım:**
```json
{
  "name": "templateden_gorev_olustur", 
  "arguments": {
    "template_id": "bug_raporu_id",
    "degerler": {
      "baslik": "Bug fix",
      "aciklama": "Login sorunu", 
      "modul": "Authentication",
      "ortam": "production",
      "oncelik": "yuksek"
    }
  }
}
```

### Error Message
Bu araç çağrıldığında aşağıdaki hata mesajı döner:

> ❌ gorev_olustur artık kullanılmıyor!
> 
> Template kullanımı zorunludur. Önce kullanılabilir template'leri görmek için:
> template_listele
> 
> Sonra template kullanarak görev oluşturmak için:
> templateden_gorev_olustur template_id='...' degerler={...}

### ✅ Çözüm
Artık [templateden_gorev_olustur](#templateden_gorev_olustur) kullanın. Bu daha iyi çünkü:
- **Tutarlılık**: Her görev belirli standartlara uygun
- **Kalite**: Zorunlu alanlar eksik bilgi girişini engeller  
- **Otomasyon**: Template tipine göre otomatik workflow
- **Raporlama**: Görev tipine göre detaylı metrikler

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Görev oluşturuldu: API authentication implementasyonu (ID: 550e8400-e29b-41d4-a716-446655440000)\n  Proje: E-ticaret Sitesi"
  }]
}
```

---

## gorev_listele

Görevleri filtreleyerek listeler.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama | Varsayılan |
|-----------|-----|---------|----------|------------|
| `durum` | string | ❌ | Filtrelenecek durum: `beklemede`, `devam_ediyor`, `tamamlandi` | Tümü |
| `tum_projeler` | boolean | ❌ | Tüm projelerdeki görevleri göster | `false` |
| `sirala` | string | ❌ | Sıralama: `son_tarih_asc`, `son_tarih_desc` | - |
| `filtre` | string | ❌ | Zaman filtresi: `acil` (7 gün içinde), `gecmis` (gecikmiş) | - |
| `etiket` | string | ❌ | Etiket adına göre filtreleme | - |
| `limit` | number | ❌ | Maksimum görev sayısı (pagination) | 50 |
| `offset` | number | ❌ | Kaç görev atlanacak (pagination) | 0 |

### Örnek Kullanım

**Tüm görevler:**
```json
{
  "name": "gorev_listele",
  "arguments": {}
}
```

**Duruma göre filtreleme:**
```json
{
  "name": "gorev_listele",
  "arguments": {
    "durum": "devam_ediyor"
  }
}
```

**Tüm projelerdeki görevler:**
```json
{
  "name": "gorev_listele",
  "arguments": {
    "tum_projeler": true
  }
}
```

**Acil görevler (7 gün içinde son tarih):**
```json
{
  "name": "gorev_listele",
  "arguments": {
    "filtre": "acil",
    "sirala": "son_tarih_asc"
  }
}
```

**Etiketle filtreleme:**
```json
{
  "name": "gorev_listele",
  "arguments": {
    "etiket": "güvenlik"
  }
}
```

**Pagination kullanımı:**
```json
{
  "name": "gorev_listele",
  "arguments": {
    "limit": 100,
    "offset": 100
  }
}
```

**Not:** 
- `tum_projeler` parametresi `false` veya verilmezse ve aktif proje varsa, sadece aktif projenin görevleri listelenir.
- Pagination özelliği v0.8.1+ sürümünde eklenmiştir. Büyük görev listeleri için token limit hatalarını önler.

### Yanıt

```markdown
## Görev Listesi

- [devam_ediyor] API authentication implementasyonu (yuksek öncelik)
  JWT tabanlı authentication sistemi kur. Refresh token desteği olmalı.
- [beklemede] README dosyasını güncelle (orta öncelik)
- [tamamlandi] Veritabanı şemasını oluştur (yuksek öncelik)
  User ve Task tabloları oluşturuldu.
```

---

## gorev_detay

Bir görevin detaylı bilgilerini markdown formatında görüntüler. Bağımlılık bilgileri her zaman gösterilir (boş olsa bile).

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `id` | string | ✅ | Görev ID'si |

### Örnek Kullanım

```json
{
  "name": "gorev_detay",
  "arguments": {
    "id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Yanıt

```markdown
# API authentication implementasyonu

## 📋 Genel Bilgiler
- **ID:** 550e8400-e29b-41d4-a716-446655440000
- **Durum:** devam_ediyor
- **Öncelik:** yuksek
- **Oluşturma Tarihi:** 2024-01-15 14:30:00
- **Son Güncelleme:** 2024-01-16 10:45:00
- **Proje:** E-ticaret Sitesi

## 📝 Açıklama
JWT tabanlı authentication sistemi kur. Refresh token desteği olmalı.

### Yapılacaklar:
- [ ] JWT library entegrasyonu
- [ ] User authentication endpoint
- [ ] Token refresh mekanizması
- [ ] Rate limiting

---

*Son güncelleme: 28 June 2025*
```

---

## gorev_guncelle

Bir görevin durumunu günceller.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `id` | string | ✅ | Görev ID'si |
| `durum` | string | ✅ | Yeni durum: `beklemede`, `devam_ediyor`, `tamamlandi` |

### Örnek Kullanım

```json
{
  "name": "gorev_guncelle",
  "arguments": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "durum": "tamamlandi"
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Görev güncellendi: 550e8400-e29b-41d4-a716-446655440000 → tamamlandi"
  }]
}
```

---

## gorev_duzenle

Bir görevin başlık, açıklama, öncelik veya proje bilgilerini düzenler.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama | Değerler |
|-----------|-----|---------|----------|----------|
| `id` | string | ✅ | Görev ID'si | - |
| `baslik` | string | ❌ | Yeni başlık | - |
| `aciklama` | string | ❌ | Yeni açıklama (markdown destekler) | - |
| `oncelik` | string | ❌ | Yeni öncelik seviyesi | `dusuk`, `orta`, `yuksek` |
| `proje_id` | string | ❌ | Yeni proje ID'si | - |

**Not:** En az bir düzenleme alanı belirtilmelidir.

### Örnek Kullanım

**Başlık ve açıklama güncelleme:**
```json
{
  "name": "gorev_duzenle",
  "arguments": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "baslik": "JWT Authentication Sistemı v2",
    "aciklama": "## JWT Authentication\n\n- Refresh token desteği\n- Role-based access control\n- Session management"
  }
}
```

**Öncelik değiştirme:**
```json
{
  "name": "gorev_duzenle",
  "arguments": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "oncelik": "dusuk"
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Görev düzenlendi: 550e8400-e29b-41d4-a716-446655440000"
  }]
}
```

---

## gorev_sil

Bir görevi kalıcı olarak siler.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `id` | string | ✅ | Görev ID'si |
| `onay` | boolean | ✅ | Silme işlemini onaylamak için `true` olmalı |

### Örnek Kullanım

```json
{
  "name": "gorev_sil",
  "arguments": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "onay": true
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Görev silindi: API authentication implementasyonu (ID: 550e8400-e29b-41d4-a716-446655440000)"
  }]
}
```

**Uyarı:** Bu işlem geri alınamaz!

---

## proje_olustur

Yeni bir proje oluşturur. Projeler görevleri gruplamak için kullanılır.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama | Varsayılan |
|-----------|-----|---------|----------|------------|
| `isim` | string | ✅ | Proje ismi | - |
| `tanim` | string | ❌ | Proje tanımı/açıklaması | "" |

### Örnek Kullanım

**Basit proje:**
```json
{
  "name": "proje_olustur",
  "arguments": {
    "isim": "E-ticaret Sitesi"
  }
}
```

**Detaylı proje:**
```json
{
  "name": "proje_olustur",
  "arguments": {
    "isim": "Mobil Uygulama v2.0",
    "tanim": "React Native ile cross-platform mobil uygulama. iOS ve Android desteği, offline çalışma özelliği."
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Proje oluşturuldu: Mobil Uygulama v2.0 (ID: 6ba7b810-9dad-11d1-80b4-00c04fd430c8)"
  }]
}
```

---

## proje_listele

Sistemdeki tüm projeleri görev sayılarıyla birlikte listeler.

### Parametreler

Bu araç parametre almaz.

### Örnek Kullanım

```json
{
  "name": "proje_listele",
  "arguments": {}
}
```

### Yanıt

```markdown
## Proje Listesi

### E-ticaret Sitesi
- **ID:** 6ba7b810-9dad-11d1-80b4-00c04fd430c8
- **Tanım:** Online satış platformu geliştirme projesi
- **Oluşturma:** 15 Jan 2024, 10:00
- **Görev Sayısı:** 12

### Mobil Uygulama v2.0
- **ID:** 6ba7b814-9dad-11d1-80b4-00c04fd430c8
- **Tanım:** React Native ile cross-platform mobil uygulama
- **Oluşturma:** 20 Jan 2024, 14:30
- **Görev Sayısı:** 8
```

---

## proje_gorevleri

Belirtilen projeye ait görevleri durum gruplarına göre listeler.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama | Varsayılan |
|-----------|-----|---------|----------|------------|
| `proje_id` | string | ✅ | Proje ID'si | - |
| `limit` | number | ❌ | Maksimum görev sayısı (pagination) | 50 |
| `offset` | number | ❌ | Kaç görev atlanacak (pagination) | 0 |

### Örnek Kullanım

```json
{
  "name": "proje_gorevleri",
  "arguments": {
    "proje_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
  }
}
```

### Yanıt

```markdown
## E-ticaret Sitesi - Görevler

### 🔵 Devam Ediyor
- **API authentication implementasyonu** (yuksek öncelik)
  JWT tabanlı authentication sistemi kur
  `ID: 550e8400-e29b-41d4-a716-446655440000`

### ⚪ Beklemede  
- **Ödeme sistemi entegrasyonu** (yuksek öncelik)
  Stripe ve PayPal entegrasyonu
  `ID: 550e8400-e29b-41d4-a716-446655440001`
- **Ürün arama özelliği** (orta öncelik)
  Elasticsearch ile gelişmiş arama
  `ID: 550e8400-e29b-41d4-a716-446655440002`

### ✅ Tamamlandı
- ~~Veritabanı şeması tasarımı~~ (yuksek öncelik)
  `ID: 550e8400-e29b-41d4-a716-446655440003`
```

---

## ozet_goster

Sistem genelinde görev ve proje özetini gösterir.

### Parametreler

Bu araç parametre almaz.

### Örnek Kullanım

```json
{
  "name": "ozet_goster",
  "arguments": {}
}
```

### Yanıt

```markdown
## Özet Rapor

**Toplam Proje:** 3
**Toplam Görev:** 15

### Durum Dağılımı
- Beklemede: 8
- Devam Ediyor: 3
- Tamamlandı: 4

### Öncelik Dağılımı
- Yüksek: 5
- Orta: 7
- Düşük: 3
```

---

## proje_aktif_yap

Bir projeyi aktif proje olarak ayarlar. Aktif proje ayarlandığında, `templateden_gorev_olustur` ve `gorev_listele` komutları varsayılan olarak bu projeyi kullanır.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `proje_id` | string | ✅ | Aktif yapılacak proje ID'si |

### Örnek Kullanım

```json
{
  "name": "proje_aktif_yap",
  "arguments": {
    "proje_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Aktif proje ayarlandı: E-ticaret Sitesi"
  }]
}
```

---

## aktif_proje_goster

Mevcut aktif projeyi ve detaylarını gösterir.

### Parametreler

Bu araç parametre almaz.

### Örnek Kullanım

```json
{
  "name": "aktif_proje_goster",
  "arguments": {}
}
```

### Yanıt

```markdown
## Aktif Proje

**Proje:** E-ticaret Sitesi
**ID:** 6ba7b810-9dad-11d1-80b4-00c04fd430c8
**Açıklama:** Online satış platformu geliştirme projesi
**Görev Sayısı:** 12
```

Aktif proje yoksa:
```
Henüz aktif proje ayarlanmamış.
```

---

## aktif_proje_kaldir

Aktif proje ayarını kaldırır. Bu işlemden sonra görev oluşturma ve listeleme işlemleri tüm projeler üzerinde çalışır.

### Parametreler

Bu araç parametre almaz.

### Örnek Kullanım

```json
{
  "name": "aktif_proje_kaldir",
  "arguments": {}
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Aktif proje ayarı kaldırıldı."
  }]
}
```

---

## template_listele

Kullanılabilir görev şablonlarını listeler. Şablonlar görev oluşturmayı hızlandırır ve standartlaştırır.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `kategori` | string | ❌ | Filtrelenecek kategori (Teknik, Özellik, Araştırma) |

### Örnek Kullanım

```json
{
  "name": "template_listele",
  "arguments": {
    "kategori": "Teknik"
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "## 📋 Görev Template'leri\n\n### Teknik\n\n#### Bug Raporu\n- **ID:** `39f28dbd-10f3-454c-8b35-52ae6b7ea391`\n- **Açıklama:** Yazılım hatası bildirimi için detaylı template\n- **Başlık Şablonu:** `🐛 [{{modul}}] {{baslik}}`\n- **Alanlar:**\n  - `baslik` (text) *(zorunlu)*\n  - `aciklama` (text) *(zorunlu)*\n  - `modul` (text) *(zorunlu)*\n  - `ortam` (select) *(zorunlu)* - seçenekler: development, staging, production\n  - `adimlar` (text) *(zorunlu)*\n  - `beklenen` (text) *(zorunlu)*\n  - `mevcut` (text) *(zorunlu)*\n  - `ekler` (text)\n  - `cozum` (text)\n  - `oncelik` (select) *(zorunlu)* - varsayılan: orta - seçenekler: dusuk, orta, yuksek\n  - `etiketler` (text) - varsayılan: bug\n\n💡 **Kullanım:** `templateden_gorev_olustur` komutunu template ID'si ve alan değerleriyle kullanın."
  }]
}
```

---

## templateden_gorev_olustur

Seçilen şablonu kullanarak özelleştirilmiş bir görev oluşturur.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `template_id` | string | ✅ | Kullanılacak template'in ID'si |
| `degerler` | object | ✅ | Template alanları için değerler (key-value çiftleri) |

### Örnek Kullanım

```json
{
  "name": "templateden_gorev_olustur",
  "arguments": {
    "template_id": "39f28dbd-10f3-454c-8b35-52ae6b7ea391",
    "degerler": {
      "baslik": "Login butonu çalışmıyor",
      "aciklama": "Kullanıcı giriş sayfasında login butonu tıklamaya yanıt vermiyor",
      "modul": "auth",
      "ortam": "production",
      "adimlar": "1. Login sayfasına git\n2. Email ve şifre gir\n3. Login butonuna tıkla",
      "beklenen": "Kullanıcı ana sayfaya yönlendirilmeli",
      "mevcut": "Hiçbir şey olmuyor, buton tepki vermiyor",
      "oncelik": "yuksek",
      "etiketler": "bug,acil,auth"
    }
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Template kullanılarak görev oluşturuldu: 🐛 [auth] Login butonu çalışmıyor (ID: d7f4e8b9-2a1c-4f5e-9d3b-8c1a2e3f4d5b)"
  }]
}
```

---

## gorev_bagimlilik_ekle

Görevler arası bağımlılık oluşturur. Bir görevin başka bir göreve bağımlı olmasını sağlar.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `kaynak_id` | string | ✅ | Önce tamamlanması gereken görevin ID'si |
| `hedef_id` | string | ✅ | Bağımlı görevin ID'si |
| `baglanti_tipi` | string | ✅ | Bağlantı tipi (genellikle "onceki") |

### Örnek Kullanım

```json
{
  "name": "gorev_bagimlilik_ekle",
  "arguments": {
    "kaynak_id": "550e8400-e29b-41d4-a716-446655440000",
    "hedef_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "baglanti_tipi": "onceki"
  }
}
```

### Yanıt

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Bağımlılık eklendi: 550e8400-e29b-41d4-a716-446655440000 -> 7c9e6679-7425-40de-944b-e07fc1f90ae7 (onceki)"
  }]
}
```

**Not:** Bağımlılık eklendikten sonra, hedef görev "devam_ediyor" durumuna geçmek için kaynak görevin "tamamlandi" durumunda olması gerekir.

---

## 🔄 Gelecek Sürümlerde Eklenecek Araçlar

### Planlanan Araçlar

1. <s>**gorev_bagla** - Görevler arası bağımlılık oluşturma</s> ✅ Eklendi (gorev_bagimlilik_ekle)
2. **gorev_ara** - Görevlerde arama yapma
3. **gorev_filtrele** - Çoklu kriterlere göre filtreleme
4. <s>**gorev_etiketle** - Görevlere etiket ekleme</s> ✅ Eklendi (gorev_olustur ile)
5. **gorev_not_ekle** - Göreve not/yorum ekleme
6. **proje_sil** - Proje silme (görevleriyle birlikte)
7. **rapor_olustur** - Detaylı proje raporları
8. **proje_ihrac** - Projeyi JSON/Markdown formatında dışa aktarma
9. **proje_ice_aktar** - JSON formatında proje içe aktarma
10. **gorev_istatistik** - Görev tamamlanma süreleri ve istatistikler

### Özellik Önerileri

Yeni araç önerileri için [GitHub Issues](https://github.com/msenol/gorev/issues) üzerinden talepte bulunabilirsiniz.

---

## 💡 Kullanım İpuçları

### 1. ID Yönetimi
- Görev ID'leri UUID formatındadır
- Claude genellikle son oluşturulan görevin ID'sini hatırlar
- ID yerine görev başlığı ile referans verebilirsiniz

### 2. Durum Geçişleri
Önerilen durum geçiş sırası:
```
beklemede → devam_ediyor → tamamlandi
```

### 3. Öncelik Seviyeleri
- **yuksek**: Acil ve kritik işler
- **orta**: Normal iş akışı
- **dusuk**: İleride yapılabilecek işler

### 4. Hata Durumları

| Hata Kodu | Açıklama | Çözüm |
|-----------|----------|-------|
| -32602 | Geçersiz parametreler | Parametre tiplerini kontrol edin |
| -32000 | İşlem hatası | Görev ID'sinin doğru olduğundan emin olun |

---

## gorev_altgorev_olustur

Ana görevin altında yeni bir alt görev oluşturur.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama | Varsayılan |
|-----------|-----|---------|----------|------------|
| `parent_id` | string | ✅ | Ana görev ID'si | - |
| `baslik` | string | ✅ | Alt görev başlığı | - |
| `aciklama` | string | ❌ | Alt görev açıklaması | "" |
| `oncelik` | string | ❌ | Öncelik seviyesi | `orta` |
| `son_tarih` | string | ❌ | Son teslim tarihi (YYYY-AA-GG) | - |
| `etiketler` | string | ❌ | Virgülle ayrılmış etiketler | - |

### Örnek Kullanım

```json
{
  "name": "gorev_altgorev_olustur",
  "arguments": {
    "parent_id": "550e8400-e29b-41d4-a716-446655440000",
    "baslik": "API endpoint'lerini test et",
    "aciklama": "Tüm REST API endpoint'lerinin unit test'lerini yaz",
    "oncelik": "yuksek"
  }
}
```

---

## gorev_ust_degistir

Bir görevin üst görevini değiştirir veya kök seviyeye taşır.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `gorev_id` | string | ✅ | Taşınacak görev ID'si |
| `yeni_parent_id` | string | ❌ | Yeni ana görev ID'si (boş ise kök seviyeye taşır) |

### Örnek Kullanım

```json
{
  "name": "gorev_ust_degistir",
  "arguments": {
    "gorev_id": "550e8400-e29b-41d4-a716-446655440001",
    "yeni_parent_id": "550e8400-e29b-41d4-a716-446655440002"
  }
}
```

---

## gorev_hiyerarsi_goster

Bir görevin tüm hiyerarşisini (üst görevler ve alt görevler) gösterir.

### Parametreler

| Parametre | Tip | Zorunlu | Açıklama |
|-----------|-----|---------|----------|
| `gorev_id` | string | ✅ | Hiyerarşisi gösterilecek görev ID'si |

### Örnek Kullanım

```json
{
  "name": "gorev_hiyerarsi_goster",
  "arguments": {
    "gorev_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

### Çıktı Formatı

```markdown
# Görev Hiyerarşisi: Ana Proje

## 📊 Hiyerarşi İstatistikleri
- **Toplam alt görev**: 3
- **Tamamlanan**: 1 (33%)
- **Devam eden**: 2 (67%)

## 🌳 Üst Görevler
*Bu görev kök seviyededir*

## 📋 Alt Görevler
└─ [🔄] Backend API (yüksek öncelik)
  └─ [✓] Veritabanı tasarımı (orta öncelik)
  └─ [⏳] API endpoint'leri (yüksek öncelik)
```

---

## 📚 İlgili Dokümantasyon

- [Kullanım Kılavuzu](kullanim.md) - Pratik kullanım örnekleri
- [Örnekler](ornekler.md) - Gerçek dünya senaryoları
- [API Referansı](api-referans.md) - Programatik erişim