package gorev

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type IsYonetici struct {
	veriYonetici VeriYoneticiInterface
}

func YeniIsYonetici(veriYonetici VeriYoneticiInterface) *IsYonetici {
	return &IsYonetici{
		veriYonetici: veriYonetici,
	}
}

// VeriYonetici returns the data manager interface
func (iy *IsYonetici) VeriYonetici() VeriYoneticiInterface {
	return iy.veriYonetici
}

func (iy *IsYonetici) GorevOlustur(baslik, aciklama, oncelik, projeID, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	var sonTarih *time.Time
	if sonTarihStr != "" {
		t, err := time.Parse("2006-01-02", sonTarihStr)
		if err != nil {
			return nil, fmt.Errorf("geçersiz son tarih formatı (YYYY-AA-GG olmalı): %w", err)
		}
		sonTarih = &t
	}

	gorev := &Gorev{
		ID:              uuid.New().String(),
		Baslik:          baslik,
		Aciklama:        aciklama,
		Oncelik:         oncelik,
		Durum:           "beklemede",
		ProjeID:         projeID,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
		SonTarih:        sonTarih,
	}

	if err := iy.veriYonetici.GorevKaydet(gorev); err != nil {
		return nil, fmt.Errorf("görev kaydedilemedi: %w", err)
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := iy.veriYonetici.EtiketleriGetirVeyaOlustur(etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf("etiketler işlenemedi: %w", err)
		}
		if err := iy.veriYonetici.GorevEtiketleriniAyarla(gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf("görev etiketleri ayarlanamadı: %w", err)
		}
		gorev.Etiketler = etiketler
	}

	return gorev, nil
}

func (iy *IsYonetici) GorevListele(durum, sirala, filtre string) ([]*Gorev, error) {
	gorevler, err := iy.veriYonetici.GorevleriGetir(durum, sirala, filtre)
	if err != nil {
		return nil, err
	}

	// Performans optimizasyonu: Tüm görevlerin ID'lerini topla
	if len(gorevler) == 0 {
		return gorevler, nil
	}

	gorevIDs := make([]string, len(gorevler))
	for i, gorev := range gorevler {
		gorevIDs[i] = gorev.ID
	}

	// Tek sorguda tüm görevlerin bağımlılık sayılarını hesapla (N+1 sorgu problemi çözüldü)
	bagimliSayilari, err := iy.veriYonetici.BulkBagimlilikSayilariGetir(gorevIDs)
	if err != nil {
		// Hata durumunda bile devam et, sadece bağımlılık sayıları 0 olarak kalır
		bagimliSayilari = make(map[string]int)
	}

	// Tek sorguda tüm görevlerin tamamlanmamış bağımlılık sayılarını hesapla
	tamamlanmamisSayilari, err := iy.veriYonetici.BulkTamamlanmamiaBagimlilikSayilariGetir(gorevIDs)
	if err != nil {
		// Hata durumunda bile devam et, sadece tamamlanmamış bağımlılık sayıları 0 olarak kalır
		tamamlanmamisSayilari = make(map[string]int)
	}

	// TODO: Bu göreve bağımlı olanları hesaplamak için de bulk query eklenebilir
	// Şimdilik bu kısmı basitleştirip performans sorununu çözmek öncelikli

	// BuGoreveBagimliSayisi hesaplaması için tüm görevleri kontrol et
	buGoreveBagimliSayilari := make(map[string]int)
	for _, gorev := range gorevler {
		buGoreveBagimliSayilari[gorev.ID] = 0
	}

	// Her görev için bu göreve bağımlı olan görevleri say
	for _, gorev := range gorevler {
		baglantilar, err := iy.veriYonetici.BaglantilariGetir(gorev.ID)
		if err == nil {
			for _, baglanti := range baglantilar {
				// Eğer bu görev kaynak ise (başka görevler buna bağımlı)
				if baglanti.KaynakID == gorev.ID {
					buGoreveBagimliSayilari[gorev.ID]++
				}
			}
		}
	}

	// Her görev için hesaplanan değerleri ata
	for _, gorev := range gorevler {
		if count, exists := bagimliSayilari[gorev.ID]; exists {
			gorev.BagimliGorevSayisi = count
		}

		if count, exists := tamamlanmamisSayilari[gorev.ID]; exists {
			gorev.TamamlanmamisBagimlilikSayisi = count
		}

		if count, exists := buGoreveBagimliSayilari[gorev.ID]; exists {
			gorev.BuGoreveBagimliSayisi = count
		}
	}

	return gorevler, nil
}

func (iy *IsYonetici) GorevDurumGuncelle(id, durum string) error {
	// Validate status values
	validStatuses := []string{"beklemede", "devam_ediyor", "tamamlandi", "iptal"}
	isValidStatus := false
	for _, validStatus := range validStatuses {
		if durum == validStatus {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return fmt.Errorf("geçersiz durum: %s. Geçerli durumlar: %v", durum, validStatuses)
	}

	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return fmt.Errorf("görev bulunamadı: %w", err)
	}

	// Eğer görev "devam_ediyor" durumuna geçiyorsa, bağımlılıkları kontrol et
	if durum == "devam_ediyor" && gorev.Durum == "beklemede" {
		bagimli, tamamlanmamislar, err := iy.GorevBagimliMi(id)
		if err != nil {
			return fmt.Errorf("bağımlılık kontrolü başarısız: %w", err)
		}

		if !bagimli {
			return fmt.Errorf("bu görev başlatılamaz, önce şu görevler tamamlanmalı: %v", tamamlanmamislar)
		}
	}

	// Eğer görev "tamamlandi" durumuna geçiyorsa, tüm alt görevlerin tamamlandığını kontrol et
	if durum == "tamamlandi" && gorev.Durum != "tamamlandi" {
		altGorevler, err := iy.veriYonetici.AltGorevleriGetir(id)
		if err != nil {
			return fmt.Errorf("alt görevler kontrol edilemedi: %w", err)
		}

		for _, altGorev := range altGorevler {
			if altGorev.Durum != "tamamlandi" {
				return fmt.Errorf("bu görev tamamlanamaz, önce tüm alt görevler tamamlanmalı")
			}
		}
	}

	gorev.Durum = durum
	gorev.GuncellemeTarih = time.Now()

	return iy.veriYonetici.GorevGuncelle(gorev)
}

func (iy *IsYonetici) ProjeOlustur(isim, tanim string) (*Proje, error) {
	proje := &Proje{
		ID:              uuid.New().String(),
		Isim:            isim,
		Tanim:           tanim,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}

	if err := iy.veriYonetici.ProjeKaydet(proje); err != nil {
		return nil, fmt.Errorf("proje kaydedilemedi: %w", err)
	}

	return proje, nil
}

func (iy *IsYonetici) GorevGetir(id string) (*Gorev, error) {
	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return nil, fmt.Errorf("görev bulunamadı: %w", err)
	}
	return gorev, nil
}

func (iy *IsYonetici) ProjeGetir(id string) (*Proje, error) {
	proje, err := iy.veriYonetici.ProjeGetir(id)
	if err != nil {
		return nil, fmt.Errorf("proje bulunamadı: %w", err)
	}
	return proje, nil
}

func (iy *IsYonetici) GorevDuzenle(id, baslik, aciklama, oncelik, projeID, sonTarihStr string, baslikVar, aciklamaVar, oncelikVar, projeVar, sonTarihVar bool) error {
	// Önce mevcut görevi al
	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return fmt.Errorf("görev bulunamadı: %w", err)
	}

	// Sadece belirtilen alanları güncelle
	if baslikVar && baslik != "" {
		gorev.Baslik = baslik
	}
	if aciklamaVar {
		gorev.Aciklama = aciklama
	}
	if oncelikVar && oncelik != "" {
		gorev.Oncelik = oncelik
	}
	if projeVar {
		// Proje değiştiriliyorsa, tüm alt görevleri de taşı
		if gorev.ProjeID != projeID {
			altGorevler, err := iy.veriYonetici.TumAltGorevleriGetir(id)
			if err != nil {
				return fmt.Errorf("alt görevler alınamadı: %w", err)
			}

			// Tüm alt görevlerin projesini güncelle
			for _, altGorev := range altGorevler {
				altGorev.ProjeID = projeID
				altGorev.GuncellemeTarih = time.Now()
				if err := iy.veriYonetici.GorevGuncelle(altGorev); err != nil {
					return fmt.Errorf("alt görev güncellenemedi: %w", err)
				}
			}
		}
		gorev.ProjeID = projeID
	}
	if sonTarihVar {
		if sonTarihStr == "" {
			gorev.SonTarih = nil
		} else {
			t, err := time.Parse("2006-01-02", sonTarihStr)
			if err != nil {
				return fmt.Errorf("geçersiz son tarih formatı (YYYY-AA-GG olmalı): %w", err)
			}
			gorev.SonTarih = &t
		}
	}

	gorev.GuncellemeTarih = time.Now()

	return iy.veriYonetici.GorevGuncelle(gorev)
}

func (iy *IsYonetici) GorevSil(id string) error {
	// Önce görevin var olduğunu kontrol et
	_, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return fmt.Errorf("görev bulunamadı: %w", err)
	}

	// Alt görevleri kontrol et
	altGorevler, err := iy.veriYonetici.AltGorevleriGetir(id)
	if err != nil {
		return fmt.Errorf("alt görevler kontrol edilemedi: %w", err)
	}

	if len(altGorevler) > 0 {
		return fmt.Errorf("bu görev silinemez, önce %d alt görev silinmeli", len(altGorevler))
	}

	return iy.veriYonetici.GorevSil(id)
}

func (iy *IsYonetici) ProjeListele() ([]*Proje, error) {
	return iy.veriYonetici.ProjeleriGetir()
}

func (iy *IsYonetici) ProjeGorevleri(projeID string) ([]*Gorev, error) {
	gorevler, err := iy.veriYonetici.ProjeGorevleriGetir(projeID)
	if err != nil {
		return nil, err
	}

	// Her görev için bağımlılık sayılarını hesapla
	for _, gorev := range gorevler {
		// Bu görevin bağımlılıklarını al (bu görev başka görevlere bağımlı)
		baglantilar, err := iy.veriYonetici.BaglantilariGetir(gorev.ID)
		if err == nil && len(baglantilar) > 0 {
			gorev.BagimliGorevSayisi = len(baglantilar)

			// Tamamlanmamış bağımlılıkları say
			tamamlanmamisSayisi := 0
			for _, baglanti := range baglantilar {
				hedefGorev, err := iy.veriYonetici.GorevGetir(baglanti.HedefID)
				if err == nil && hedefGorev.Durum != "tamamlandi" {
					tamamlanmamisSayisi++
				}
			}
			gorev.TamamlanmamisBagimlilikSayisi = tamamlanmamisSayisi
		}

		// Bu göreve bağımlı olan görevleri bul
		buGoreveBagimliSayisi := 0
		tumGorevler, _ := iy.veriYonetici.GorevleriGetir("", "", "")
		for _, digerGorev := range tumGorevler {
			digerBaglantilar, err := iy.veriYonetici.BaglantilariGetir(digerGorev.ID)
			if err == nil {
				for _, baglanti := range digerBaglantilar {
					if baglanti.HedefID == gorev.ID {
						buGoreveBagimliSayisi++
						break
					}
				}
			}
		}
		gorev.BuGoreveBagimliSayisi = buGoreveBagimliSayisi
	}

	return gorevler, nil
}

func (iy *IsYonetici) ProjeGorevSayisi(projeID string) (int, error) {
	gorevler, err := iy.veriYonetici.ProjeGorevleriGetir(projeID)
	if err != nil {
		return 0, err
	}
	return len(gorevler), nil
}

func (iy *IsYonetici) AktifProjeAyarla(projeID string) error {
	return iy.veriYonetici.AktifProjeAyarla(projeID)
}

func (iy *IsYonetici) AktifProjeGetir() (*Proje, error) {
	projeID, err := iy.veriYonetici.AktifProjeGetir()
	if err != nil {
		return nil, err
	}
	if projeID == "" {
		return nil, nil // Aktif proje yok
	}
	return iy.veriYonetici.ProjeGetir(projeID)
}

func (iy *IsYonetici) AktifProjeKaldir() error {
	return iy.veriYonetici.AktifProjeKaldir()
}

func (iy *IsYonetici) OzetAl() (*Ozet, error) {
	gorevler, err := iy.veriYonetici.GorevleriGetir("", "", "")
	if err != nil {
		return nil, fmt.Errorf("görevler alınamadı: %w", err)
	}

	projeler, err := iy.veriYonetici.ProjeleriGetir()
	if err != nil {
		return nil, fmt.Errorf("projeler alınamadı: %w", err)
	}

	ozet := &Ozet{
		ToplamProje: len(projeler),
		ToplamGorev: len(gorevler),
	}

	for _, gorev := range gorevler {
		switch gorev.Durum {
		case "beklemede":
			ozet.BeklemedeGorev++
		case "devam_ediyor":
			ozet.DevamEdenGorev++
		case "tamamlandi":
			ozet.TamamlananGorev++
		}

		switch gorev.Oncelik {
		case "yuksek":
			ozet.YuksekOncelik++
		case "orta":
			ozet.OrtaOncelik++
		case "dusuk":
			ozet.DusukOncelik++
		}
	}

	return ozet, nil
}

func (iy *IsYonetici) GorevBagimlilikEkle(kaynakID, hedefID, baglantiTipi string) (*Baglanti, error) {
	// Görevlerin var olup olmadığını kontrol et
	_, err := iy.veriYonetici.GorevGetir(kaynakID)
	if err != nil {
		return nil, fmt.Errorf("kaynak görev bulunamadı: %w", err)
	}
	_, err = iy.veriYonetici.GorevGetir(hedefID)
	if err != nil {
		return nil, fmt.Errorf("hedef görev bulunamadı: %w", err)
	}

	baglanti := &Baglanti{
		ID:          uuid.New().String(),
		KaynakID:    kaynakID,
		HedefID:     hedefID,
		BaglantiTip: baglantiTipi,
	}

	if err := iy.veriYonetici.BaglantiEkle(baglanti); err != nil {
		return nil, fmt.Errorf("bağlantı eklenemedi: %w", err)
	}

	return baglanti, nil
}

func (iy *IsYonetici) GorevBaglantilariGetir(gorevID string) ([]*Baglanti, error) {
	return iy.veriYonetici.BaglantilariGetir(gorevID)
}

// GorevBagimliMi görevi başlatmak için tüm bağımlılıkların tamamlanıp tamamlanmadığını kontrol eder
func (iy *IsYonetici) GorevBagimliMi(gorevID string) (bool, []string, error) {
	baglantilar, err := iy.veriYonetici.BaglantilariGetir(gorevID)
	if err != nil {
		return false, nil, fmt.Errorf("bağlantılar alınamadı: %w", err)
	}

	var tamamlanmamisBagimliliklar []string

	for _, baglanti := range baglantilar {
		// Bu görev hedef konumundaysa ve bağlantı tipi "onceki" ise
		if baglanti.HedefID == gorevID && baglanti.BaglantiTip == "onceki" {
			// Kaynak görevin durumunu kontrol et
			kaynakGorev, err := iy.veriYonetici.GorevGetir(baglanti.KaynakID)
			if err != nil {
				return false, nil, fmt.Errorf("bağımlı görev bulunamadı: %w", err)
			}

			// Eğer bağımlı görev tamamlanmamışsa
			if kaynakGorev.Durum != "tamamlandi" {
				tamamlanmamisBagimliliklar = append(tamamlanmamisBagimliliklar, kaynakGorev.Baslik)
			}
		}
	}

	return len(tamamlanmamisBagimliliklar) == 0, tamamlanmamisBagimliliklar, nil
}

// TemplateListele kullanılabilir template'leri listeler
func (iy *IsYonetici) TemplateListele(kategori string) ([]*GorevTemplate, error) {
	return iy.veriYonetici.TemplateListele(kategori)
}

// TemplatedenGorevOlustur template kullanarak görev oluşturur
func (iy *IsYonetici) TemplatedenGorevOlustur(templateID string, degerler map[string]string) (*Gorev, error) {
	return iy.veriYonetici.TemplatedenGorevOlustur(templateID, degerler)
}

// AltGorevOlustur mevcut bir görevin altına yeni görev oluşturur
func (iy *IsYonetici) AltGorevOlustur(parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error) {
	// Parent görevi kontrol et
	parent, err := iy.veriYonetici.GorevGetir(parentID)
	if err != nil {
		return nil, fmt.Errorf("üst görev bulunamadı: %w", err)
	}

	var sonTarih *time.Time
	if sonTarihStr != "" {
		t, err := time.Parse("2006-01-02", sonTarihStr)
		if err != nil {
			return nil, fmt.Errorf("geçersiz son tarih formatı (YYYY-AA-GG olmalı): %w", err)
		}
		sonTarih = &t
	}

	gorev := &Gorev{
		ID:              uuid.New().String(),
		Baslik:          baslik,
		Aciklama:        aciklama,
		Oncelik:         oncelik,
		Durum:           "beklemede",
		ProjeID:         parent.ProjeID, // Alt görev aynı projede olmalı
		ParentID:        parentID,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
		SonTarih:        sonTarih,
	}

	if err := iy.veriYonetici.GorevKaydet(gorev); err != nil {
		return nil, fmt.Errorf("alt görev kaydedilemedi: %w", err)
	}

	if len(etiketIsimleri) > 0 {
		etiketler, err := iy.veriYonetici.EtiketleriGetirVeyaOlustur(etiketIsimleri)
		if err != nil {
			return nil, fmt.Errorf("etiketler işlenemedi: %w", err)
		}
		if err := iy.veriYonetici.GorevEtiketleriniAyarla(gorev.ID, etiketler); err != nil {
			return nil, fmt.Errorf("etiketler ayarlanamadı: %w", err)
		}
		gorev.Etiketler = etiketler
	}

	return gorev, nil
}

// GorevUstDegistir bir görevin üst görevini değiştirir
func (iy *IsYonetici) GorevUstDegistir(gorevID, yeniParentID string) error {
	// Görevi kontrol et
	gorev, err := iy.veriYonetici.GorevGetir(gorevID)
	if err != nil {
		return fmt.Errorf("görev bulunamadı: %w", err)
	}

	// Yeni parent varsa kontrol et
	if yeniParentID != "" {
		parent, err := iy.veriYonetici.GorevGetir(yeniParentID)
		if err != nil {
			return fmt.Errorf("yeni üst görev bulunamadı: %w", err)
		}

		// Circular dependency kontrolü
		circular, err := iy.veriYonetici.DaireBagimliligiKontrolEt(gorevID, yeniParentID)
		if err != nil {
			return fmt.Errorf("dairesel bağımlılık kontrolü başarısız: %w", err)
		}
		if circular {
			return fmt.Errorf("dairesel bağımlılık tespit edildi")
		}

		// Alt görev ve üst görev aynı projede olmalı
		if gorev.ProjeID != parent.ProjeID {
			return fmt.Errorf("alt görev ve üst görev aynı projede olmalı")
		}
	}

	return iy.veriYonetici.ParentIDGuncelle(gorevID, yeniParentID)
}

// GorevHiyerarsiGetir bir görevin tam hiyerarşi bilgilerini getirir
func (iy *IsYonetici) GorevHiyerarsiGetir(gorevID string) (*GorevHiyerarsi, error) {
	return iy.veriYonetici.GorevHiyerarsiGetir(gorevID)
}

// AltGorevleriGetir bir görevin doğrudan alt görevlerini getirir
func (iy *IsYonetici) AltGorevleriGetir(parentID string) ([]*Gorev, error) {
	return iy.veriYonetici.AltGorevleriGetir(parentID)
}

// TumAltGorevleriGetir bir görevin tüm alt görev ağacını getirir
func (iy *IsYonetici) TumAltGorevleriGetir(parentID string) ([]*Gorev, error) {
	return iy.veriYonetici.TumAltGorevleriGetir(parentID)
}
