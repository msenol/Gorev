# MCP Araçları Dokümantasyonu

Bu dokümanda Gorev MCP Server'ın sunduğu tüm araçların detaylı açıklamaları bulunmaktadır.

## 📋 Araç Listesi (Toplam: 25)

### Görev Yönetimi

#### 1. `gorev_olustur`
Yeni bir görev oluşturur.

**Parametreler:**
- `baslik` (zorunlu): Görev başlığı
- `aciklama` (opsiyonel): Görev açıklaması (markdown destekli)
- `oncelik` (opsiyonel): Öncelik seviyesi (yuksek, orta, dusuk) - varsayılan: orta
- `proje_id` (opsiyonel): Proje ID'si - belirtilmezse aktif proje kullanılır
- `son_tarih` (opsiyonel): Son tarih (YYYY-AA-GG formatında)
- `etiketler` (opsiyonel): Virgülle ayrılmış etiket listesi

**Örnek:**
```json
{
  "baslik": "API dokümantasyonu yaz",
  "aciklama": "REST API endpoint'lerini dokümante et",
  "oncelik": "yuksek",
  "son_tarih": "2025-07-15",
  "etiketler": "dokümantasyon,api"
}
```

#### 2. `gorev_listele`
Görevleri listeler ve filtreler.

**Parametreler:**
- `durum` (opsiyonel): Filtrelenecek durum (beklemede, devam_ediyor, tamamlandi)
- `tum_projeler` (opsiyonel): true ise tüm projelerdeki görevler, false/belirtilmezse sadece aktif proje
- `sirala` (opsiyonel): son_tarih_asc veya son_tarih_desc
- `filtre` (opsiyonel): acil (7 gün içinde) veya gecmis (süresi geçmiş)
- `etiket` (opsiyonel): Filtrelenecek etiket ismi
- `limit` (opsiyonel): Gösterilecek maksimum görev sayısı (varsayılan: 50)
- `offset` (opsiyonel): Atlanacak görev sayısı - pagination için (varsayılan: 0)

**Çıktı:** Hiyerarşik ağaç yapısında görev listesi

#### 3. `gorev_detay`
Bir görevin detaylı bilgilerini markdown formatında gösterir.

**Parametreler:**
- `id` (zorunlu): Görev ID'si

**Çıktı:** Son tarihler, etiketler, bağımlılıklar ve alt görevler dahil tüm detaylar

#### 4. `gorev_guncelle`
Görev durumunu günceller.

**Parametreler:**
- `id` (zorunlu): Görev ID'si
- `durum` (zorunlu): Yeni durum (beklemede, devam_ediyor, tamamlandi)

**Notlar:**
- "devam_ediyor" durumuna geçmek için bağımlılıklar tamamlanmış olmalı
- "tamamlandi" durumuna geçmek için tüm alt görevler tamamlanmış olmalı

#### 5. `gorev_duzenle`
Görev özelliklerini düzenler.

**Parametreler:**
- `id` (zorunlu): Görev ID'si
- `baslik` (opsiyonel): Yeni başlık
- `aciklama` (opsiyonel): Yeni açıklama
- `oncelik` (opsiyonel): Yeni öncelik
- `proje_id` (opsiyonel): Yeni proje ID'si - alt görevler de taşınır
- `son_tarih` (opsiyonel): Yeni son tarih

#### 6. `gorev_sil`
Görevi siler.

**Parametreler:**
- `id` (zorunlu): Görev ID'si
- `onay` (zorunlu): Silme onayı (true/false)

**Not:** Alt görevi olan görevler silinemez.

#### 7. `gorev_bagimlilik_ekle`
İki görev arasında bağımlılık oluşturur.

**Parametreler:**
- `kaynak_id` (zorunlu): Kaynak görev ID'si
- `hedef_id` (zorunlu): Hedef görev ID'si
- `baglanti_tipi` (zorunlu): Bağlantı tipi (örn: 'engelliyor', 'ilişkili')

### Alt Görev Yönetimi

#### 8. `gorev_altgorev_olustur`
Mevcut bir görevin altına yeni görev oluşturur.

**Parametreler:**
- `parent_id` (zorunlu): Üst görevin ID'si
- `baslik` (zorunlu): Alt görevin başlığı
- `aciklama` (opsiyonel): Alt görevin açıklaması
- `oncelik` (opsiyonel): Öncelik seviyesi (varsayılan: orta)
- `son_tarih` (opsiyonel): Son tarih (YYYY-AA-GG formatında)
- `etiketler` (opsiyonel): Virgülle ayrılmış etiket listesi

**Not:** Alt görev, üst görevin projesini otomatik olarak devralır.

#### 9. `gorev_ust_degistir`
Bir görevin üst görevini değiştirir veya kök göreve taşır.

**Parametreler:**
- `gorev_id` (zorunlu): Taşınacak görevin ID'si
- `yeni_parent_id` (opsiyonel): Yeni üst görevin ID'si (boş string = kök göreve taşı)

**Not:** Dairesel bağımlılık kontrolü yapılır.

#### 10. `gorev_hiyerarsi_goster`
Bir görevin tam hiyerarşisini ve alt görev istatistiklerini gösterir.

**Parametreler:**
- `gorev_id` (zorunlu): Görevin ID'si

**Çıktı:**
- Üst görev hiyerarşisi
- Alt görev istatistikleri (toplam, tamamlanan, devam eden, beklemede)
- İlerleme yüzdesi
- Doğrudan alt görevler listesi

### Görev Şablonları

#### 11. `template_listele`
Kullanılabilir görev şablonlarını listeler.

**Parametreler:**
- `kategori` (opsiyonel): Filtrelenecek kategori (Teknik, Özellik, Araştırma vb.)

#### 12. `templateden_gorev_olustur`
Seçilen şablonu kullanarak özelleştirilmiş bir görev oluşturur.

**Parametreler:**
- `template_id` (zorunlu): Şablon ID'si
- `degerler` (zorunlu): Şablon alanları için değerler (key-value çiftleri)

### Proje Yönetimi

#### 13. `proje_olustur`
Yeni proje oluşturur.

**Parametreler:**
- `isim` (zorunlu): Proje ismi
- `tanim` (zorunlu): Proje açıklaması

#### 14. `proje_listele`
Tüm projeleri görev sayılarıyla birlikte listeler.

**Parametreler:** Yok

#### 15. `proje_gorevleri`
Bir projenin görevlerini duruma göre gruplandırarak listeler.

**Parametreler:**
- `proje_id` (zorunlu): Proje ID'si
- `limit` (opsiyonel): Gösterilecek maksimum görev sayısı (varsayılan: 50)
- `offset` (opsiyonel): Atlanacak görev sayısı - pagination için (varsayılan: 0)

#### 16. `proje_aktif_yap`
Belirtilen projeyi aktif proje olarak ayarlar.

**Parametreler:**
- `proje_id` (zorunlu): Proje ID'si

#### 17. `aktif_proje_goster`
Mevcut aktif projeyi gösterir.

**Parametreler:** Yok

#### 18. `aktif_proje_kaldir`
Aktif proje ayarını kaldırır.

**Parametreler:** Yok

### Raporlama

#### 19. `ozet_goster`
Sistem genelinde özet istatistikler gösterir.

**Parametreler:** Yok

**Çıktı:**
- Toplam proje sayısı
- Toplam görev sayısı
- Durum bazlı görev dağılımı
- Öncelik bazlı görev dağılımı

## 🔧 Kullanım İpuçları

1. **Hiyerarşik Yapı**: Alt görevler kullanarak karmaşık projeleri organize edin
2. **Bağımlılık Yönetimi**: Görevler arası ilişkileri tanımlayarak iş akışı oluşturun
3. **Şablon Kullanımı**: Sık kullanılan görev tiplerini şablonlarla hızlıca oluşturun
4. **Etiketleme**: Görevleri kategorize etmek için etiketleri aktif kullanın
5. **Son Tarih Takibi**: Acil ve gecikmiş görevleri filtreleyerek önceliklendirin

## 📝 Notlar

- Tüm araçlar Turkish domain language kullanır (gorev, proje, durum, vb.)
- Görev açıklamaları full markdown formatını destekler
- Tarih formatı: YYYY-AA-GG (örn: 2025-07-30)
- ID'ler UUID formatındadır