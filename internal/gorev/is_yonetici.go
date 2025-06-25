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

func (iy *IsYonetici) GorevOlustur(baslik, aciklama, oncelik, projeID string) (*Gorev, error) {
	gorev := &Gorev{
		ID:              uuid.New().String(),
		Baslik:          baslik,
		Aciklama:        aciklama,
		Oncelik:         oncelik,
		Durum:           "beklemede",
		ProjeID:         projeID,
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}

	if err := iy.veriYonetici.GorevKaydet(gorev); err != nil {
		return nil, fmt.Errorf("görev kaydedilemedi: %w", err)
	}

	return gorev, nil
}

func (iy *IsYonetici) GorevListele(durum string) ([]*Gorev, error) {
	return iy.veriYonetici.GorevleriGetir(durum)
}

func (iy *IsYonetici) GorevDurumGuncelle(id, durum string) error {
	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return fmt.Errorf("görev bulunamadı: %w", err)
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

func (iy *IsYonetici) GorevDetayAl(id string) (*Gorev, error) {
	gorev, err := iy.veriYonetici.GorevGetir(id)
	if err != nil {
		return nil, fmt.Errorf("görev bulunamadı: %w", err)
	}
	return gorev, nil
}

func (iy *IsYonetici) ProjeDetayAl(id string) (*Proje, error) {
	proje, err := iy.veriYonetici.ProjeGetir(id)
	if err != nil {
		return nil, fmt.Errorf("proje bulunamadı: %w", err)
	}
	return proje, nil
}

func (iy *IsYonetici) GorevDuzenle(id, baslik, aciklama, oncelik, projeID string, baslikVar, aciklamaVar, oncelikVar, projeVar bool) error {
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
	gorevler, err := iy.veriYonetici.GorevleriGetir("")
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
