package constants

// MCP tool parameter names to eliminate hardcoded strings throughout the codebase
const (
	// Core task parameters
	ParamID        = "id"
	ParamBaslik    = "baslik"
	ParamAciklama  = "aciklama"
	ParamOncelik   = "oncelik"
	ParamDurum     = "durum"
	ParamProjeID   = "proje_id"
	ParamSonTarih  = "son_tarih"
	ParamEtiketler = "etiketler"
	ParamParentID  = "parent_id"
	ParamGorevID   = "gorev_id"

	// Task management parameters
	ParamTumProjeler = "tum_projeler"
	ParamSirala      = "sirala"
	ParamFiltre      = "filtre"
	ParamEtiket      = "etiket"
	ParamLimit       = "limit"
	ParamOffset      = "offset"
	ParamOnay        = "onay"

	// Dependency parameters
	ParamKaynakID     = "kaynak_id"
	ParamHedefID      = "hedef_id"
	ParamBaglantiTipi = "baglanti_tipi"
	ParamYeniParentID = "yeni_parent_id"

	// Project parameters
	ParamIsim  = "isim"
	ParamTanim = "tanim"

	// Template parameters
	ParamTemplateID = "template_id"
	ParamDegerler   = "degerler"
	ParamKategori   = "kategori"

	// AI context parameters
	ParamTaskID  = "task_id"
	ParamQuery   = "query"
	ParamUpdates = "updates"
)

// MCP tool names to eliminate hardcoded strings
const (
	ToolGorevListele            = "gorev_listele"
	ToolGorevDetay              = "gorev_detay"
	ToolGorevOlustur            = "gorev_olustur"
	ToolGorevGuncelle           = "gorev_guncelle"
	ToolGorevDuzenle            = "gorev_duzenle"
	ToolGorevSil                = "gorev_sil"
	ToolGorevAltgorevOlustur    = "gorev_altgorev_olustur"
	ToolGorevUstDegistir        = "gorev_ust_degistir"
	ToolGorevHiyerarsiGoster    = "gorev_hiyerarsi_goster"
	ToolGorevBagimlilikEkle     = "gorev_bagimlilik_ekle"
	ToolProjeOlustur            = "proje_olustur"
	ToolProjeListele            = "proje_listele"
	ToolProjeGorevleri          = "proje_gorevleri"
	ToolProjeAktifYap           = "proje_aktif_yap"
	ToolAktifProjeGoster        = "aktif_proje_goster"
	ToolAktifProjeKaldir        = "aktif_proje_kaldir"
	ToolTemplateListele         = "template_listele"
	ToolTemplatedenGorevOlustur = "templateden_gorev_olustur"
	ToolOzetGoster              = "ozet_goster"
	ToolGorevSetActive          = "gorev_set_active"
	ToolGorevGetActive          = "gorev_get_active"
	ToolGorevRecent             = "gorev_recent"
	ToolGorevContextSummary     = "gorev_context_summary"
	ToolGorevBatchUpdate        = "gorev_batch_update"
	ToolGorevNLPQuery           = "gorev_nlp_query"
	ToolFileWatcherStart        = "file_watcher_start"
	ToolFileWatcherStop         = "file_watcher_stop"
	ToolFileWatcherStatus       = "file_watcher_status"
)

// JSON response field names
const (
	ResponseGorevler    = "gorevler"
	ResponseProjeler    = "projeler"
	ResponseTemplateler = "templateler"
	ResponseOzet        = "ozet"
	ResponseSayi        = "sayi"
	ResponseToplam      = "toplam"
	ResponseSayfa       = "sayfa"
	ResponseHasNextPage = "has_next_page"
	ResponseHasPrevPage = "has_prev_page"
	ResponseSuccess     = "success"
	ResponseMessage     = "message"
	ResponseError       = "error"
	ResponseData        = "data"
)

// Database field names to eliminate hardcoded strings
const (
	DBFieldID              = "id"
	DBFieldBaslik          = "baslik"
	DBFieldAciklama        = "aciklama"
	DBFieldDurum           = "durum"
	DBFieldOncelik         = "oncelik"
	DBFieldProjeID         = "proje_id"
	DBFieldParentID        = "parent_id"
	DBFieldSonTarih        = "son_tarih"
	DBFieldOlusturmaTarih  = "olusturma_tarih"
	DBFieldGuncellemeTarih = "guncelleme_tarih"
	DBFieldIsim            = "isim"
	DBFieldTanim           = "tanim"
	DBFieldAktif           = "aktif"
)

// Sort parameter values
const (
	SortSonTarihAsc  = "son_tarih_asc"
	SortSonTarihDesc = "son_tarih_desc"
	SortOncelikAsc   = "oncelik_asc"
	SortOncelikDesc  = "oncelik_desc"
)

// Filter parameter values
const (
	FilterAcil    = "acil"
	FilterGecmis  = "gecmis"
	FilterBuggun  = "buggun"
	FilterBuHafta = "bu_hafta"
)
