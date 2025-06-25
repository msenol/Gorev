package gorev

import "time"

// Gorev temel görev yapısı
type Gorev struct {
	ID              string    `json:"id"`
	Baslik          string    `json:"baslik"`
	Aciklama        string    `json:"aciklama"`
	Durum           string    `json:"durum"`
	Oncelik         string    `json:"oncelik"`
	ProjeID         string    `json:"proje_id,omitempty"`
	OlusturmaTarih  time.Time `json:"olusturma_tarih"`
	GuncellemeTarih time.Time `json:"guncelleme_tarih"`
}

// Proje görevleri gruplamak için kullanılır
type Proje struct {
	ID              string    `json:"id"`
	Isim            string    `json:"isim"`
	Tanim           string    `json:"tanim"`
	OlusturmaTarih  time.Time `json:"olusturma_tarih"`
	GuncellemeTarih time.Time `json:"guncelleme_tarih"`
}

// Ozet sistem durumu özeti
type Ozet struct {
	ToplamProje     int `json:"toplam_proje"`
	ToplamGorev     int `json:"toplam_gorev"`
	BeklemedeGorev  int `json:"beklemede_gorev"`
	DevamEdenGorev  int `json:"devam_eden_gorev"`
	TamamlananGorev int `json:"tamamlanan_gorev"`
	YuksekOncelik   int `json:"yuksek_oncelik"`
	OrtaOncelik     int `json:"orta_oncelik"`
	DusukOncelik    int `json:"dusuk_oncelik"`
}

// Baglanti görevler arası bağlantı
type Baglanti struct {
	ID         string `json:"id"`
	KaynakID   string `json:"kaynak_id"`
	HedefID    string `json:"hedef_id"`
	BaglantiTip string `json:"baglanti_tip"`
}