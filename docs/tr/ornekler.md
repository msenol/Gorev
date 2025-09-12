# Gorev Kullanım Örnekleri

> **Versiyon**: Bu dokümantasyon v0.7.0-beta.1 için geçerlidir.  
> **Son Güncelleme**: 29 June 2025

Bu dokümanda Gorev'in çeşitli kullanım senaryolarını ve örneklerini bulabilirsiniz.

## İçindekiler

- [Temel Kullanım](#temel-kullanım)
- [İleri Seviye Senaryolar](#ileri-seviye-senaryolar)
- [Şablon Kullanımı](#şablon-kullanımı)
- [Bağımlılık Yönetimi](#bağımlılık-yönetimi)
- [Etiketleme ve Filtreleme](#etiketleme-ve-filtreleme)
- [Proje Yönetimi İş Akışları](#proje-yönetimi-iş-akışları)

## Temel Kullanım

### Basit Görev Oluşturma

```bash
# Claude Desktop/Code üzerinden
"Yeni bir görev oluştur: Veritabanı şemasını güncelle"

# MCP komutu olarak
gorev_olustur(
  baslik="Veritabanı şemasını güncelle",
  aciklama="User tablosuna last_login alanı ekle",
  oncelik="orta"
)
```

### Görev Durumu Güncelleme

```bash
# Görevi başlat
gorev_guncelle(id=1, durum="devam_ediyor")

# Görevi tamamla
gorev_guncelle(id=1, durum="tamamlandı")
```

## İleri Seviye Senaryolar

### Sprint Planlama

```bash
# Sprint projesi oluştur
proje_olustur(isim="Sprint 2025-W04", tanim="4. hafta sprint görevleri")

# Aktif proje yap
proje_aktif_yap(proje_id=5)

# Sprint görevlerini ekle
gorev_olustur(
  baslik="API rate limiting implementasyonu",
  aciklama="Redis tabanlı rate limiting ekle",
  oncelik="yuksek",
  son_tarih="2025-01-31"
)

gorev_olustur(
  baslik="Kullanıcı profil sayfası",
  aciklama="Profil düzenleme ve görüntüleme",
  oncelik="orta",
  son_tarih="2025-01-30",
  etiketler="frontend,ux"
)
```

### Bug Tracking İş Akışı

```bash
# Bug şablonundan görev oluştur
templateden_gorev_olustur(
  template_id="bug-report",
  degerler={
    "baslik": "Login butonu çalışmıyor",
    "adimlar": "1. Login sayfasına git\n2. Bilgileri gir\n3. Butona tıkla",
    "beklenen": "Kullanıcı giriş yapmalı",
    "gerceklesen": "Hiçbir şey olmuyor",
    "oncelik": "yuksek"
  }
)

# Bug'ı etiketle
gorev_duzenle(id=15, etiketler="bug,production,urgent")

# Acil bug'ları listele
gorev_listele(filtre="acil", etiket="bug")
```

## Şablon Kullanımı

### Mevcut Şablonları Görüntüleme

```bash
# Tüm şablonları listele
template_listele()

# Kategori bazlı listeleme
template_listele(kategori="gelistirme")
```

### Feature Request Şablonu

```bash
templateden_gorev_olustur(
  template_id="feature-request",
  degerler={
    "baslik": "Dark mode desteği",
    "amac": "Kullanıcıların göz yorgunluğunu azaltmak",
    "fayda": "Gece kullanımında konfor sağlar",
    "kabul_kriterleri": "- Tema değiştirme toggle'ı\n- Ayarlar kaydedilmeli\n- Tüm sayfalar desteklemeli"
  }
)
```

### Teknik Borç Şablonu

```bash
templateden_gorev_olustur(
  template_id="technical-debt",
  degerler={
    "baslik": "Legacy API client refactor",
    "mevcut_durum": "3 yıllık kod, test yok",
    "onerilen_cozum": "Modern HTTP client kullan",
    "etki": "Bakım maliyeti %40 azalır",
    "oncelik": "orta"
  }
)
```

## Bağımlılık Yönetimi

### Sıralı Görev Zinciri

```bash
# Ana görevler
gorev_olustur(baslik="Database migration planlama", oncelik="yuksek")  # id: 20
gorev_olustur(baslik="Migration script yazma", oncelik="yuksek")       # id: 21
gorev_olustur(baslik="Test ortamında deneme", oncelik="yuksek")       # id: 22
gorev_olustur(baslik="Production deployment", oncelik="yuksek")        # id: 23

# Bağımlılıkları tanımla
gorev_bagimlilik_ekle(kaynak_id=21, hedef_id=20, baglanti_tipi="tamamlanmali")
gorev_bagimlilik_ekle(kaynak_id=22, hedef_id=21, baglanti_tipi="tamamlanmali")
gorev_bagimlilik_ekle(kaynak_id=23, hedef_id=22, baglanti_tipi="tamamlanmali")
```

### Paralel Görev Grupları

```bash
# Frontend ve backend paralel geliştirme
gorev_olustur(baslik="API endpoint tasarımı", oncelik="yuksek")           # id: 30
gorev_olustur(baslik="Frontend mockup hazırlama", oncelik="orta")        # id: 31
gorev_olustur(baslik="Backend implementation", oncelik="yuksek")          # id: 32
gorev_olustur(baslik="Frontend implementation", oncelik="yuksek")         # id: 33
gorev_olustur(baslik="Integration testing", oncelik="yuksek")            # id: 34

# Bağımlılıklar
gorev_bagimlilik_ekle(kaynak_id=32, hedef_id=30, baglanti_tipi="tamamlanmali")
gorev_bagimlilik_ekle(kaynak_id=33, hedef_id=31, baglanti_tipi="tamamlanmali")
gorev_bagimlilik_ekle(kaynak_id=34, hedef_id=32, baglanti_tipi="tamamlanmali")
gorev_bagimlilik_ekle(kaynak_id=34, hedef_id=33, baglanti_tipi="tamamlanmali")
```

## Etiketleme ve Filtreleme

### Etiket Bazlı İş Akışı

```bash
# Görevleri kategorize et
gorev_duzenle(id=10, etiketler="backend,api,v2")
gorev_duzenle(id=11, etiketler="frontend,react,v2")
gorev_duzenle(id=12, etiketler="devops,deployment,v2")

# Etiketlere göre filtrele
gorev_listele(etiket="v2")                    # Tüm v2 görevleri
gorev_listele(etiket="backend", durum="beklemede")  # Bekleyen backend görevleri
```

### Acil Görev Yönetimi

```bash
# Önümüzdeki 7 gün içinde deadline'ı olan görevler
gorev_listele(filtre="acil")

# Gecikmiş görevler
gorev_listele(filtre="gecmis")

# Yüksek öncelikli ve acil görevler
gorev_listele(filtre="acil", durum="beklemede")
```

## Proje Yönetimi İş Akışları

### Çoklu Proje Takibi

```bash
# Farklı projeler oluştur
proje_olustur(isim="Website Redesign", tanim="Kurumsal site yenileme")
proje_olustur(isim="Mobile App v2", tanim="Mobil uygulama major güncelleme")
proje_olustur(isim="Backend Refactor", tanim="Mikroservis mimarisine geçiş")

# Proje bazlı görev oluşturma
gorev_olustur(
  baslik="Homepage tasarımı",
  oncelik="yuksek",
  proje_id=1,
  etiketler="design,ui"
)

gorev_olustur(
  baslik="Push notification sistemi",
  oncelik="orta",
  proje_id=2,
  etiketler="mobile,feature"
)

# Proje durumunu görüntüle
proje_listele()  # Tüm projeleri görev sayılarıyla listele
proje_gorevleri(proje_id=1)  # Specific proje görevleri
```

### Haftalık Rapor Hazırlama

```bash
# Genel özet
ozet_goster()

# Tamamlanan görevleri listele
gorev_listele(durum="tamamlandı", tum_projeler=true)

# Her proje için durum
for proje_id in [1, 2, 3]:
    proje_gorevleri(proje_id=proje_id)
```

### Milestone Tracking

```bash
# Milestone projesi oluştur
proje_olustur(isim="Q1 2025 Milestones", tanim="İlk çeyrek hedefleri")

# Major milestone'ları görev olarak ekle
gorev_olustur(
  baslik="v2.0 Release",
  aciklama="Major version release with new features",
  oncelik="yuksek",
  son_tarih="2025-03-31",
  etiketler="milestone,release"
)

# Alt görevleri milestone'a bağla
gorev_olustur(baslik="Feature freeze", oncelik="yuksek", son_tarih="2025-03-15")
gorev_olustur(baslik="Beta testing", oncelik="yuksek", son_tarih="2025-03-22")
gorev_olustur(baslik="Documentation update", oncelik="orta", son_tarih="2025-03-25")

# Bağımlılıkları kur
gorev_bagimlilik_ekle(kaynak_id=45, hedef_id=44, baglanti_tipi="tamamlanmali")
gorev_bagimlilik_ekle(kaynak_id=46, hedef_id=44, baglanti_tipi="tamamlanmali")
gorev_bagimlilik_ekle(kaynak_id=47, hedef_id=45, baglanti_tipi="tamamlanmali")
```

## İpuçları ve En İyi Uygulamalar

1. **Tutarlı Etiketleme**: Takım içinde standart etiket seti belirleyin
2. **Düzenli Güncelleme**: Görev durumlarını güncel tutun
3. **Bağımlılık Planlama**: Karmaşık projelerde bağımlılıkları önceden tanımlayın
4. **Şablon Kullanımı**: Tekrarlayan görev türleri için şablon kullanın
5. **Tarih Takibi**: Önemli görevlere mutlaka son_tarih ekleyin

## İlgili Dokümantasyon

- [MCP Araçları Referansı](../guides/user/mcp-tools.md)
- [Kullanım Kılavuzu](../guides/user/usage.md)
- [API Referansı](../api/reference.md)

---

<div align="center">

*✨ Bu kapsamlı örnek senaryolar Claude (Anthropic) ile birlikte hazırlanmıştır - Gerçek dünya kullanım örnekleriyle zenginleştirilmiş dokümantasyon*

</div>