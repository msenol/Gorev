package gorev

import (
	"database/sql"
	"time"
)

// VeriYoneticiInterface defines the data access interface
type VeriYoneticiInterface interface {
	GorevKaydet(gorev *Gorev) error
	GorevGetir(id string) (*Gorev, error)
	GorevleriGetir(durum, sirala, filtre string) ([]*Gorev, error)
	GorevGuncelle(taskID string, params interface{}) error
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
	BulkBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error)
	BulkTamamlanmamiaBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error)
	BulkBuGoreveBagimliSayilariGetir(gorevIDs []string) (map[string]int, error)
	EtiketleriGetirVeyaOlustur(isimler []string) ([]*Etiket, error)
	GorevEtiketleriniAyarla(gorevID string, etiketler []*Etiket) error
	TemplateOlustur(template *GorevTemplate) error
	TemplateListele(kategori string) ([]*GorevTemplate, error)
	TemplateGetir(templateID string) (*GorevTemplate, error)
	TemplateAliasIleGetir(alias string) (*GorevTemplate, error)
	TemplateIDVeyaAliasIleGetir(idOrAlias string) (*GorevTemplate, error)
	TemplatedenGorevOlustur(templateID string, degerler map[string]string) (*Gorev, error)
	VarsayilanTemplateleriOlustur() error
	AltGorevleriGetir(parentID string) ([]*Gorev, error)
	TumAltGorevleriGetir(parentID string) ([]*Gorev, error)
	UstGorevleriGetir(gorevID string) ([]*Gorev, error)
	GorevHiyerarsiGetir(gorevID string) (*GorevHiyerarsi, error)
	ParentIDGuncelle(gorevID, yeniParentID string) error
	DaireBagimliligiKontrolEt(gorevID, hedefParentID string) (bool, error)
	AltGorevOlustur(parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error)

	// AI Context Management methods
	AIContextGetir() (*AIContext, error)
	AIContextKaydet(context *AIContext) error
	AIInteractionKaydet(interaction *AIInteraction) error
	AIInteractionlariGetir(limit int) ([]*AIInteraction, error)
	AITodayInteractionlariGetir() ([]*AIInteraction, error)
	AILastInteractionGuncelle(taskID string, timestamp time.Time) error

	// File Watcher Integration methods (using string IDs to match existing interface)
	GorevDosyaYoluEkle(taskID string, path string) error
	GorevDosyaYoluSil(taskID string, path string) error
	GorevDosyaYollariGetir(taskID string) ([]string, error)
	DosyaYoluGorevleriGetir(path string) ([]string, error)
	AIEtkilemasimKaydet(taskID string, interactionType, data, sessionID string) error
	GorevSonAIEtkilesiminiGuncelle(taskID string, timestamp time.Time) error

	// Additional methods for NLP and auto state management
	GorevDetay(taskID string) (*Gorev, error)
	GorevListele(filters map[string]interface{}) ([]*Gorev, error)
	GorevOlustur(params map[string]interface{}) (string, error)
	GorevBagimlilikGetir(taskID string) ([]*Gorev, error)
	GetDB() (*sql.DB, error)

	Kapat() error
}
