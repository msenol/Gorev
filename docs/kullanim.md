# Kullanım Kılavuzu

Gorev ile görev yönetiminin temelleri.

## 🎯 Temel Kavramlar

### Görev (Task)
- Yapılacak işlerin temel birimi
- Her görevin benzersiz bir ID'si vardır
- Durum: `beklemede`, `devam_ediyor`, `tamamlandi`
- Öncelik: `dusuk`, `orta`, `yuksek`

### Proje (Project)
- Görevleri gruplamak için kullanılır
- Birden fazla görev içerebilir
- İsteğe bağlıdır

## 🔧 Komut Satırı Kullanımı

### Sunucuyu Başlatma
```bash
# Normal modda başlat
gorev serve

# Debug modunda başlat
gorev serve --debug

# Özel veri dizini ile
gorev serve --data-dir /path/to/data
```

### Versiyon Kontrolü
```bash
gorev version
```

## 💬 Claude ile Kullanım

### Görev Oluşturma

**Basit görev:**
```
Yeni bir görev oluştur: "Dokümantasyonu güncelle"
```

**Detaylı görev:**
```
Yüksek öncelikli bir görev oluştur:
Başlık: API endpoint'lerini test et
Açıklama: Tüm REST API endpoint'lerinin unit testlerini yaz
```

### Görevleri Listeleme

**Tüm görevler:**
```
Görevleri listele
```

**Duruma göre filtreleme:**
```
Beklemedeki görevleri göster
Devam eden görevleri listele
Tamamlanmış görevleri göster
```

### Görev Durumu Güncelleme

```
[görev-id] görevini "devam ediyor" olarak güncelle
[görev-id] görevini tamamlandı olarak işaretle
```

### Proje Yönetimi

**Proje oluşturma:**
```
"Web Uygulaması" adında yeni bir proje oluştur
```

**Proje özeti:**
```
Proje özetini göster
Genel durumu özetle
```

## 📊 Gelişmiş Kullanım

### Görev Organizasyonu

1. **Proje Bazlı Çalışma:**
   ```
   1. "Backend API" projesi oluştur
   2. Bu proje için görevler ekle:
      - Kullanıcı authentication endpoint'i
      - Veritabanı migration'ları
      - API dokümantasyonu
   ```

2. **Öncelik Yönetimi:**
   ```
   Yüksek öncelikli görevleri listele
   En acil 3 görevi göster
   ```

3. **Durum Takibi:**
   ```
   Bugün tamamlanan görevleri göster
   Devam eden görevlerin özetini ver
   ```

### Workflow Örnekleri

#### Sprint Planlama
```
1. "Sprint 1" projesi oluştur
2. Sprint görevlerini ekle (her biri için tahmini süre)
3. Görevleri öncelik sırasına göre listele
4. İlk görevi "devam ediyor" olarak işaretle
```

#### Bug Takibi
```
1. "Buglar" projesi oluştur
2. Yüksek öncelikli bug görevi ekle:
   - Başlık: "Login sayfası 404 hatası"
   - Açıklama: "Production'da login sayfası açılmıyor"
3. Görevi "devam ediyor" olarak güncelle
4. Çözüldüğünde "tamamlandı" olarak işaretle
```

#### Günlük Planlama
```
1. Bugünkü görevleri listele
2. En yüksek öncelikli görevi seç
3. "Devam ediyor" olarak işaretle
4. Tamamlandığında güncelle
5. Günlük özet raporu al
```

## 🎨 İpuçları ve Püf Noktaları

### 1. Etkili Görev Başlıkları
- ❌ "Bug fix"
- ✅ "Kullanıcı giriş formunda email validasyonu düzelt"

### 2. Açıklama Kullanımı
- Bağlam bilgisi ekleyin
- Kabul kriterleri belirtin
- İlgili kaynakları not edin

### 3. Öncelik Stratejisi
- **Yüksek**: Acil ve önemli (production buglar, kritik özellikler)
- **Orta**: Önemli ama acil değil (yeni özellikler, iyileştirmeler)
- **Düşük**: Ne acil ne önemli (nice-to-have özellikler)

### 4. Durum Yönetimi
- Aynı anda sadece 1-3 görev "devam ediyor" durumunda olmalı
- Görevleri küçük, yönetilebilir parçalara bölün
- Tamamlanan görevleri düzenli olarak gözden geçirin

## 🔍 Sık Kullanılan Komutlar

### Hızlı Başlangıç
```
"Todo uygulaması" projesi oluştur ve şu görevleri ekle:
- Frontend tasarımı (orta öncelik)
- Backend API geliştirme (yüksek öncelik)  
- Veritabanı şeması (yüksek öncelik)
- Test yazma (orta öncelik)
- Deployment setup (düşük öncelik)
```

### Durum Raporu
```
Şu bilgileri ver:
- Toplam görev sayısı
- Duruma göre dağılım
- Önceliğe göre dağılım
- Aktif proje sayısı
```

### Temizlik
```
Tamamlanmış görevleri listele ve hangilerinin arşivlenebileceğini belirt
```

## ❓ Sorun Giderme

### Görev ID'si Bulma
```
"API test" içeren görevi bul ve ID'sini göster
```

### Toplu Güncelleme
```
"Backend" projesindeki tüm bekleyen görevleri listele
```

### Veri Yedekleme
Gorev otomatik olarak SQLite veritabanını kullanır. Yedekleme için:
```bash
cp ~/.gorev/data/gorev.db ~/.gorev/data/gorev.db.backup
```

## 🚀 Sonraki Adımlar

- [MCP Araçları Referansı](mcp-araclari.md) - Tüm komutların detaylı açıklaması
- [Örnekler](ornekler.md) - Gerçek kullanım senaryoları
- [Mimari](mimari.md) - Sistem nasıl çalışır?