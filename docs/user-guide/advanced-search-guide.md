# 🔍 Gelişmiş Arama ve Filtreleme Kullanım Kılavuzu

**Gorev v0.15.0** ile birlikte gelen yeni gelişmiş arama ve filtreleme sistemi kullanım kılavuzu.

## 📋 İçindekiler

1. [Hızlı Başlangıç](#hızlı-başlangıç)
2. [Gelişmiş Arama](#gelişmiş-arama)
3. [Filtre Profilleri](#filtre-profilleri)
4. [Arama Geçmişi](#arama-geçmişi)
5. [Akıllı Öneriler](#akıllı-öneriler)
6. [Pratik Örnekler](#pratik-örnekler)
7. [İpuçları ve Püf Noktaları](#ipuçları-ve-püf-noktaları)

## 🚀 Hızlı Başlangıç

### Temel Arama
En basit haliyle metin arama yapmak için:

```bash
# Tüm görevlerde "database" kelimesini ara
gorev mcp gorev_search_advanced query="database"
```

### Çoklu Filtre ile Arama
Daha spesifik sonuçlar için birden fazla filtre kombinasyonu:

```bash
# Yüksek öncelikli, beklemede olan görevlerde "bug" ara
gorev mcp gorev_search_advanced query="bug" priority='["yuksek"]' status='["beklemede"]'
```

## 🔍 Gelişmiş Arama

### 1. Full-Text Search (FTS5)
SQLite FTS5 teknolojisi ile ultra hızlı metin arama:

- **Başlık ve açıklamada arama**: Tüm görev içeriği indekslenir
- **Kelime parçası eşleştirme**: "data" kelimesi "database" içinde bulunur
- **Performans**: Binlerce görev içinde milisaniye yanıt

### 2. Bulanık Arama (Fuzzy Search)
Yazım hatalarına toleranslı arama:

```bash
# "databse" yazım hatası ile "database" bulur
gorev mcp gorev_search_advanced query="databse" enable_fuzzy=true
```

**Eşik Ayarları:**
- `fuzzy_threshold=1`: Çok hassas (1 karakter fark)
- `fuzzy_threshold=2`: Dengelenmiş (varsayılan)
- `fuzzy_threshold=3`: Toleranslı (3 karakter fark)

### 3. Çoklu Filtre Kombinasyonları

#### Durum Filtreleri
```bash
# Beklemede ve devam eden görevler
status='["beklemede", "devam_ediyor"]'
```

#### Öncelik Filtreleri
```bash
# Yüksek ve orta öncelikli görevler
priority='["yuksek", "orta"]'
```

#### Tarih Filtreleri
```bash
# Bu aydan sonra oluşturulan görevler
created_after="2024-09-01"

# Gelecek hafta teslimi olan görevler
due_after="2024-09-20" due_before="2024-09-27"
```

#### Proje ve Etiket Filtreleri
```bash
# Belirli projelerde ara
project_ids='["proje-uuid-1", "proje-uuid-2"]'

# Belirli etiketlerde ara
tags='["bug", "critical"]'
```

## 📂 Filtre Profilleri

Sık kullanılan filtre kombinasyonlarını kaydetmek ve yeniden kullanmak.

### Profil Oluşturma
```bash
gorev mcp gorev_filter_profile_create \
  name="Acil Buglar" \
  description="Yüksek öncelikli bug görevleri" \
  filters='{
    "status": ["beklemede", "devam_ediyor"],
    "priority": ["yuksek"],
    "tags": ["bug"],
    "enable_fuzzy": false
  }'
```

### Profil Kullanımı
```bash
# Profilleri listele
gorev mcp gorev_filter_profile_list

# Belirli profili getir
gorev mcp gorev_filter_profile_get id="profil-uuid"

# Profili güncelle
gorev mcp gorev_filter_profile_update id="profil-uuid" name="Yeni İsim"

# Profili sil
gorev mcp gorev_filter_profile_delete id="profil-uuid"
```

### Örnek Kullanışlı Profiller

#### 1. Acil Görevler
```json
{
  "name": "Acil Görevler",
  "filters": {
    "priority": ["yuksek"],
    "status": ["beklemede", "devam_ediyor"],
    "due_before": "2024-09-30"
  }
}
```

#### 2. Bu Hafta Tamamlanacaklar
```json
{
  "name": "Bu Hafta Teslim",
  "filters": {
    "status": ["beklemede", "devam_ediyor"],
    "due_after": "2024-09-16",
    "due_before": "2024-09-22"
  }
}
```

#### 3. Kod Review Görevleri
```json
{
  "name": "Code Review",
  "filters": {
    "tags": ["review", "code"],
    "status": ["beklemede"]
  }
}
```

## 📜 Arama Geçmişi

Önceki aramalarınızı takip edin ve tekrarlayın.

```bash
# Son 10 aramayı görüntüle
gorev mcp gorev_search_history limit=10

# Son 50 aramayı görüntüle
gorev mcp gorev_search_history limit=50
```

**Otomatik Kayıt:**
- Her `gorev_search_advanced` çağrısı otomatik kaydedilir
- Tarih ve saat bilgisi ile saklanır
- En son aramalar en üstte görüntülenir

## 🤖 Akıllı Öneriler

NLP tabanlı akıllı arama önerileri.

```bash
# "veritaban" için öneriler al
gorev mcp gorev_search_suggestions query="veritaban"
```

**Öneriler şunları içerir:**
- **NLP Önerileri**: "veritaban" → "database", "veri tabanı", "db"
- **Geçmiş Aramalar**: Daha önce yapılan benzer aramalar
- **Yaygın Kalıplar**: Sık kullanılan arama kombinasyonları
- **Zaman Tabanlı**: "bugün", "bu hafta", "bu ay" gibi öneriler

## 💡 Pratik Örnekler

### 1. Günlük Görev Kontrolü
```bash
# Bugün yapılacak yüksek öncelikli görevler
gorev mcp gorev_search_advanced \
  priority='["yuksek"]' \
  status='["beklemede", "devam_ediyor"]' \
  due_before="2024-09-17"
```

### 2. Proje Temizliği
```bash
# Belirli bir projede tamamlanmış görevler
gorev mcp gorev_search_advanced \
  project_ids='["proje-uuid"]' \
  status='["tamamlandi"]' \
  created_before="2024-08-01"
```

### 3. Bug Avı
```bash
# Tüm bug raporları (bulanık arama ile)
gorev mcp gorev_search_advanced \
  query="bug" \
  tags='["bug", "hata", "sorun"]' \
  enable_fuzzy=true \
  fuzzy_threshold=2
```

### 4. Sprint Planlama
```bash
# Gelecek sprint için orta öncelikli görevler
gorev mcp gorev_search_advanced \
  priority='["orta"]' \
  status='["beklemede"]' \
  created_after="2024-09-01"
```

## 🎯 İpuçları ve Püf Noktaları

### 1. Performans Optimizasyonu
- **FTS5 kullanın**: Metin araması için en hızlı yöntem
- **Filtre kombinasyonları**: Önce dar filtreler, sonra geniş aramalar
- **Limit kullanın**: Büyük sonuç setleri için sayfa sayfa görüntüleme

### 2. Etkili Arama Stratejileri
- **Anahtar kelimeler**: Spesifik terimler kullanın
- **Etiket sistemi**: Görevleri kategorize etmek için etiketleri kullanın
- **Tarih aralıkları**: Zaman bazlı filtreleme ile sonuçları daraltın

### 3. Filtre Profili İpuçları
- **Anlamlı isimler**: Profillere açıklayıcı isimler verin
- **Dokümantasyon**: Description alanını kullanarak açıklama ekleyin
- **Periyodik güncelleme**: Kullanım alışkanlıklarınıza göre profilleri güncelleyin

### 4. Bulanık Arama İpuçları
- **Kısa kelimeler**: 3-4 harfli kelimeler için eşiği düşürün
- **Uzun kelimeler**: 10+ harfli kelimeler için eşiği artırın
- **Test edin**: Farklı eşik değerlerini deneyerek optimal sonuçları bulun

### 5. Hata Ayıklama
- **Sonuç bulunamadı**: Filtreleri gevşetin veya bulanık aramayı aktifleştirin
- **Çok fazla sonuç**: Daha spesifik filtreler ekleyin
- **Yavaş yanıt**: Arama terimini kısaltın veya filtre sayısını azaltın

## 🔧 Teknik Detaylar

### FTS5 Konfigürasyonu
- **İndekslenmiş alanlar**: başlık, açıklama, etiketler, proje adı
- **Tokenizer**: unicode61 (Türkçe karakter desteği)
- **Trigger sistemi**: Otomatik FTS indeks güncellemesi

### Bulanık Arama Algoritması
- **Levenshtein Distance**: Karakter düzeyinde benzerlik hesaplaması
- **Case insensitive**: Büyük/küçük harf duyarsız
- **Unicode desteği**: Türkçe karakterler desteklenir

### Performans Metrikleri
- **FTS5 arama**: ~1-5ms (10K görev)
- **Bulanık arama**: ~10-50ms (eşiğe bağlı)
- **Kombineli filtreler**: ~5-20ms
- **Profil yükleme**: ~1-2ms

## 🚀 Gelecek Özellikler

- **VS Code Extension**: Görsel arama arayüzü
- **Regex desteği**: Gelişmiş pattern matching
- **Saved search shortcuts**: Hızlı arama kısayolları
- **Export search results**: Arama sonuçlarını dışa aktarma
- **Search analytics**: Arama istatistikleri ve analizler

---

> 💡 **Not**: Bu kılavuz v0.15.0 sürümü için hazırlanmıştır. Güncellemeler için CHANGELOG.md dosyasını takip edin.