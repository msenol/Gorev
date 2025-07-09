package gorev

import "time"

// Gorev temel görev yapısı
type Gorev struct {
	ID              string     `json:"id"`
	Baslik          string     `json:"baslik"`
	Aciklama        string     `json:"aciklama"`
	Durum           string     `json:"durum"`
	Oncelik         string     `json:"oncelik"`
	ProjeID         string     `json:"proje_id,omitempty"`
	ParentID        string     `json:"parent_id,omitempty"`
	OlusturmaTarih  time.Time  `json:"olusturma_tarih"`
	GuncellemeTarih time.Time  `json:"guncelleme_tarih"`
	SonTarih        *time.Time `json:"son_tarih,omitempty"`
	Etiketler       []*Etiket  `json:"etiketler,omitempty"`
	AltGorevler     []*Gorev   `json:"alt_gorevler,omitempty"`
	Seviye          int        `json:"seviye,omitempty"`
	// Bağımlılık sayaçları - TreeView gösterimi için
	BagimliGorevSayisi            int `json:"bagimli_gorev_sayisi,omitempty"`
	TamamlanmamisBagimlilikSayisi int `json:"tamamlanmamis_bagimlilik_sayisi,omitempty"`
	BuGoreveBagimliSayisi         int `json:"bu_goreve_bagimli_sayisi,omitempty"`
}

// Etiket görevleri kategorize etmek için kullanılır
type Etiket struct {
	ID   string `json:"id"`
	Isim string `json:"isim"`
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
	ID          string `json:"id"`
	KaynakID    string `json:"kaynak_id"`
	HedefID     string `json:"hedef_id"`
	BaglantiTip string `json:"baglanti_tip"`
}

// GorevTemplate görev oluşturma şablonu
type GorevTemplate struct {
	ID               string            `json:"id"`
	Isim             string            `json:"isim"`
	Tanim            string            `json:"tanim"`
	VarsayilanBaslik string            `json:"varsayilan_baslik"`
	AciklamaTemplate string            `json:"aciklama_template"`
	Alanlar          []TemplateAlan    `json:"alanlar"`
	OrnekDegerler    map[string]string `json:"ornek_degerler"`
	Kategori         string            `json:"kategori"`
	Aktif            bool              `json:"aktif"`
}

// TemplateAlan template'deki özelleştirilebilir alanlar
type TemplateAlan struct {
	Isim       string   `json:"isim"`
	Tip        string   `json:"tip"` // text, select, date, number
	Zorunlu    bool     `json:"zorunlu"`
	Varsayilan string   `json:"varsayilan"`
	Secenekler []string `json:"secenekler,omitempty"`
}

// GorevHiyerarsi görev hiyerarşi bilgilerini tutar
type GorevHiyerarsi struct {
	Gorev           *Gorev   `json:"gorev"`
	UstGorevler     []*Gorev `json:"ust_gorevler,omitempty"`
	ToplamAltGorev  int      `json:"toplam_alt_gorev"`
	TamamlananAlt   int      `json:"tamamlanan_alt"`
	DevamEdenAlt    int      `json:"devam_eden_alt"`
	BeklemedeAlt    int      `json:"beklemede_alt"`
	IlerlemeYuzdesi float64  `json:"ilerleme_yuzdesi"`
}
