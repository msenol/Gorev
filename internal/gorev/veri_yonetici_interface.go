package gorev

// VeriYoneticiInterface defines the data access interface
type VeriYoneticiInterface interface {
	GorevKaydet(gorev *Gorev) error
	GorevGetir(id string) (*Gorev, error)
	GorevleriGetir(durum, sirala, filtre string) ([]*Gorev, error)
	GorevGuncelle(gorev *Gorev) error
	GorevSil(id string) error
	ProjeKaydet(proje *Proje) error
	ProjeGetir(id string) (*Proje, error)
	ProjeleriGetir() ([]*Proje, error)
	ProjeGorevleriGetir(projeID string) ([]*Gorev, error)
	AktifProjeAyarla(projeID string) error
	AktifProjeGetir() (string, error)
	AktifProjeKaldir() error
	BaglantiEkle(baglanti *Baglanti) error
	BaglantilariGetir(gorevID string) ([]*Baglanti, error)
	EtiketleriGetirVeyaOlustur(isimler []string) ([]*Etiket, error)
	GorevEtiketleriniAyarla(gorevID string, etiketler []*Etiket) error
	Kapat() error
}
