package gorev

import (
	"context"
	"database/sql"
	"time"
)

// VeriYoneticiInterface defines the data access interface
type VeriYoneticiInterface interface {
	GorevKaydet(ctx context.Context, gorev *Gorev) error
	GorevGetir(ctx context.Context, id string) (*Gorev, error)
	GorevleriGetir(ctx context.Context, durum, sirala, filtre string) ([]*Gorev, error)
	GorevGuncelle(ctx context.Context, taskID string, params interface{}) error
	GorevSil(ctx context.Context, id string) error
	ProjeKaydet(ctx context.Context, proje *Proje) error
	ProjeGetir(ctx context.Context, id string) (*Proje, error)
	ProjeleriGetir(ctx context.Context) ([]*Proje, error)
	ProjeGorevleriGetir(ctx context.Context, projeID string) ([]*Gorev, error)
	AktifProjeAyarla(ctx context.Context, projeID string) error
	AktifProjeGetir(ctx context.Context) (string, error)
	AktifProjeKaldir(ctx context.Context) error
	BaglantiEkle(ctx context.Context, baglanti *Baglanti) error
	BaglantiSil(ctx context.Context, kaynakID, hedefID string) error
	BaglantilariGetir(ctx context.Context, gorevID string) ([]*Baglanti, error)
	BulkBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error)
	BulkTamamlanmamiaBagimlilikSayilariGetir(gorevIDs []string) (map[string]int, error)
	BulkBuGoreveBagimliSayilariGetir(gorevIDs []string) (map[string]int, error)
	EtiketleriGetirVeyaOlustur(ctx context.Context, isimler []string) ([]*Etiket, error)
	GorevEtiketleriniAyarla(ctx context.Context, gorevID string, etiketler []*Etiket) error
	TemplateOlustur(ctx context.Context, template *GorevTemplate) error
	TemplateListele(ctx context.Context, kategori string) ([]*GorevTemplate, error)
	TemplateGetir(ctx context.Context, templateID string) (*GorevTemplate, error)
	TemplateAliasIleGetir(ctx context.Context, alias string) (*GorevTemplate, error)
	TemplateIDVeyaAliasIleGetir(ctx context.Context, idOrAlias string) (*GorevTemplate, error)
	TemplatedenGorevOlustur(ctx context.Context, templateID string, degerler map[string]string) (*Gorev, error)
	VarsayilanTemplateleriOlustur(ctx context.Context) error
	AltGorevleriGetir(ctx context.Context, parentID string) ([]*Gorev, error)
	TumAltGorevleriGetir(ctx context.Context, parentID string) ([]*Gorev, error)
	UstGorevleriGetir(ctx context.Context, gorevID string) ([]*Gorev, error)
	GorevHiyerarsiGetir(ctx context.Context, gorevID string) (*GorevHiyerarsi, error)
	ParentIDGuncelle(ctx context.Context, gorevID, yeniParentID string) error
	DaireBagimliligiKontrolEt(ctx context.Context, gorevID, hedefParentID string) (bool, error)
	AltGorevOlustur(ctx context.Context, parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string) (*Gorev, error)

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
	GorevDetay(ctx context.Context, taskID string) (*Gorev, error)
	GorevListele(ctx context.Context, filters map[string]interface{}) ([]*Gorev, error)
	GorevOlustur(ctx context.Context, params map[string]interface{}) (string, error)
	GorevBagimlilikGetir(ctx context.Context, taskID string) ([]*Gorev, error)
	GetDB() (*sql.DB, error)

	Kapat() error
}
