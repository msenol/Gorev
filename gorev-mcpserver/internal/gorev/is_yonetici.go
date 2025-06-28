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
	return iy.veriYonetici.GorevleriGetir(durum, sirala, filtre)
}

func (iy *IsYonetici) GorevDurumGuncelle(id, durum string) error {
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

	return iy.veriYonetici.GorevSil(id)
}

func (iy *IsYonetici) ProjeListele() ([]*Proje, error) {
	return iy.veriYonetici.ProjeleriGetir()
}

func (iy *IsYonetici) ProjeGorevleri(projeID string) ([]*Gorev, error) {
	return iy.veriYonetici.ProjeGorevleriGetir(projeID)
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
