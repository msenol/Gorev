package gorev

// VeriYoneticiInterface defines the data access interface
type VeriYoneticiInterface interface {
	GorevKaydet(gorev *Gorev) error
	GorevGetir(id string) (*Gorev, error)
	GorevleriGetir(durum string) ([]*Gorev, error)
	GorevGuncelle(gorev *Gorev) error
	GorevSil(id string) error
	ProjeKaydet(proje *Proje) error
	ProjeGetir(id string) (*Proje, error)
	ProjeleriGetir() ([]*Proje, error)
	ProjeGorevleriGetir(projeID string) ([]*Gorev, error)
	AktifProjeAyarla(projeID string) error
	AktifProjeGetir() (string, error)
	AktifProjeKaldir() error
	Kapat() error
}
