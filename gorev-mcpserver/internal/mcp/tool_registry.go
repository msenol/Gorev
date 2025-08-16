package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolRegistry handles MCP tool registration and schema definitions
type ToolRegistry struct {
	handlers *Handlers
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry(handlers *Handlers) *ToolRegistry {
	return &ToolRegistry{
		handlers: handlers,
	}
}

// RegisterAllTools registers all MCP tools with the server
func (tr *ToolRegistry) RegisterAllTools(s *server.MCPServer) {
	tr.registerTaskManagementTools(s)
	tr.registerProjectManagementTools(s)
	tr.registerTemplateTools(s)
	tr.registerAIContextTools(s)
	tr.registerFileWatcherTools(s)
	tr.registerAdvancedTools(s)
}

// registerTaskManagementTools registers core task management tools
func (tr *ToolRegistry) registerTaskManagementTools(s *server.MCPServer) {
	// Görev oluştur (deprecated)
	s.AddTool(mcp.Tool{
		Name:        "gorev_olustur",
		Description: "⚠️ KULLANIM DIŞI: Bu tool v0.10.0'dan beri kullanımdan kaldırılmıştır. Lütfen 'templateden_gorev_olustur' tool'unu kullanın.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.GorevOlustur)

	// Görev listele
	s.AddTool(mcp.Tool{
		Name:        "gorev_listele",
		Description: "Görevleri durum, proje, son teslim tarihi gibi kriterlere göre filtreleyerek ve sıralayarak listeler.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"durum": map[string]interface{}{
					"type":        "string",
					"description": "Filtrelenecek görev durumu (beklemede, devam_ediyor, tamamlandi).",
					"enum":        []string{"beklemede", "devam_ediyor", "tamamlandi"},
				},
				"sirala": map[string]interface{}{
					"type":        "string",
					"description": "Sıralama kriteri (son_tarih_asc, son_tarih_desc).",
					"enum":        []string{"son_tarih_asc", "son_tarih_desc"},
				},
				"filtre": map[string]interface{}{
					"type":        "string",
					"description": "Özel filtre türü (acil: 7 gün içinde bitenler, gecmis: vadesi geçenler).",
					"enum":        []string{"acil", "gecmis"},
				},
				"etiket": map[string]interface{}{
					"type":        "string",
					"description": "Etiket adına göre filtrele.",
				},
				"tum_projeler": map[string]interface{}{
					"type":        "boolean",
					"description": "true ise tüm projelerden görevleri göster, false ise sadece aktif projeden.",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Döndürülecek maksimum görev sayısı (varsayılan: 50, maksimum: 200).",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Atlanacak görev sayısı, sayfalama için (varsayılan: 0).",
				},
			},
		},
	}, tr.handlers.GorevListele)

	// Görev detay
	s.AddTool(mcp.Tool{
		Name:        "gorev_detay",
		Description: "Bir görevin tüm detaylarını markdown formatında gösterir. Bağımlılıklar, etiketler, son tarih gibi bilgileri içerir.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin benzersiz ID'si.",
				},
			},
			Required: []string{"id"},
		},
	}, tr.handlers.GorevDetay)

	// Görev güncelle
	s.AddTool(mcp.Tool{
		Name:        "gorev_guncelle",
		Description: "Bir görevin durumunu günceller (beklemede, devam_ediyor, tamamlandi).",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin benzersiz ID'si.",
				},
				"durum": map[string]interface{}{
					"type":        "string",
					"description": "Yeni durum.",
					"enum":        []string{"beklemede", "devam_ediyor", "tamamlandi", "iptal"},
				},
			},
			Required: []string{"id", "durum"},
		},
	}, tr.handlers.GorevGuncelle)

	// Görev düzenle
	s.AddTool(mcp.Tool{
		Name:        "gorev_duzenle",
		Description: "Bir görevin başlık, açıklama, öncelik gibi özelliklerini düzenler.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin benzersiz ID'si.",
				},
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": "Yeni başlık.",
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": "Yeni açıklama.",
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": "Yeni öncelik seviyesi.",
					"enum":        []string{"dusuk", "orta", "yuksek"},
				},
				"son_tarih": map[string]interface{}{
					"type":        "string",
					"description": "Yeni son tarih (YYYY-MM-DD formatında).",
				},
			},
			Required: []string{"id"},
		},
	}, tr.handlers.GorevDuzenle)

	// Görev sil
	s.AddTool(mcp.Tool{
		Name:        "gorev_sil",
		Description: "Bir görevi kalıcı olarak siler. DİKKAT: Bu işlem geri alınamaz!",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Silinecek görevin benzersiz ID'si.",
				},
				"onay": map[string]interface{}{
					"type":        "boolean",
					"description": "Silme işlemini onaylamak için true olmalı.",
				},
			},
			Required: []string{"id", "onay"},
		},
	}, tr.handlers.GorevSil)
}

// registerProjectManagementTools registers project management tools
func (tr *ToolRegistry) registerProjectManagementTools(s *server.MCPServer) {
	// Proje oluştur
	s.AddTool(mcp.Tool{
		Name:        "proje_olustur",
		Description: "Yeni bir proje oluşturur.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"isim": map[string]interface{}{
					"type":        "string",
					"description": "Projenin adı.",
				},
				"tanim": map[string]interface{}{
					"type":        "string",
					"description": "Projenin açıklaması.",
				},
			},
			Required: []string{"isim", "tanim"},
		},
	}, tr.handlers.ProjeOlustur)

	// Proje listele
	s.AddTool(mcp.Tool{
		Name:        "proje_listele",
		Description: "Tüm projeleri görev sayıları ile birlikte listeler.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.ProjeListele)

	// Proje görevleri
	s.AddTool(mcp.Tool{
		Name:        "proje_gorevleri",
		Description: "Belirtilen projeye ait görevleri durum gruplarına göre organize ederek gösterir.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": "Projenin benzersiz ID'si.",
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Döndürülecek maksimum görev sayısı (varsayılan: 50).",
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": "Atlanacak görev sayısı (varsayılan: 0).",
				},
			},
			Required: []string{"proje_id"},
		},
	}, tr.handlers.ProjeGorevleri)

	// Aktif proje ayarla
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_ayarla",
		Description: "Çalışılacak aktif projeyi ayarlar. Yeni görevler bu projeye otomatik atanır.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": "Aktif yapılacak projenin ID'si.",
				},
			},
			Required: []string{"proje_id"},
		},
	}, tr.handlers.AktifProjeAyarla)

	// Aktif proje göster
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_goster",
		Description: "Şu anda aktif olan projeyi gösterir.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.AktifProjeGoster)

	// Aktif proje kaldır
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_kaldir",
		Description: "Aktif proje ayarını kaldırır. Artık hiçbir proje aktif olmaz.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.AktifProjeKaldir)
}

// registerTemplateTools registers template-related tools
func (tr *ToolRegistry) registerTemplateTools(s *server.MCPServer) {
	// Template listele
	s.AddTool(mcp.Tool{
		Name:        "template_listele",
		Description: "Mevcut görev template'lerini listeler. Template'ler tutarlı görev oluşturmak için kullanılır.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"kategori": map[string]interface{}{
					"type":        "string",
					"description": "Template kategorisine göre filtrele.",
				},
			},
		},
	}, tr.handlers.TemplateListele)

	// Template'den görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "templateden_gorev_olustur",
		Description: "Belirtilen template kullanarak yeni görev oluşturur. Bu, v0.10.0+ sürümlerinde görev oluşturmanın TEK yoludur.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"template_id": map[string]interface{}{
					"type":        "string",
					"description": "Kullanılacak template'in ID'si. 'template_listele' ile mevcut template'leri görebilirsiniz.",
				},
				"degerler": map[string]interface{}{
					"type":        "object",
					"description": "Template alanları için değerler. Her template'in farklı gerekli alanları vardır.",
				},
			},
			Required: []string{"template_id", "degerler"},
		},
	}, tr.handlers.TemplatedenGorevOlustur)
}

// registerAIContextTools registers AI context management tools
func (tr *ToolRegistry) registerAIContextTools(s *server.MCPServer) {
	// AI aktif görev ayarla
	s.AddTool(mcp.Tool{
		Name:        "gorev_set_active",
		Description: "AI session için aktif görev ayarlar ve otomatik olarak durumunu 'devam_ediyor' yapar.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": "Aktif yapılacak görevin ID'si.",
				},
			},
			Required: []string{"task_id"},
		},
	}, tr.handlers.GorevSetActive)

	// AI aktif görevi getir
	s.AddTool(mcp.Tool{
		Name:        "gorev_get_active",
		Description: "AI session için şu anda aktif olan görevi döndürür.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.GorevGetActive)

	// Son görevleri getir
	s.AddTool(mcp.Tool{
		Name:        "gorev_recent",
		Description: "Son etkileşimde bulunulan görevleri döndürür.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "number",
					"description": "Döndürülecek görev sayısı (varsayılan: 5).",
				},
			},
		},
	}, tr.handlers.GorevRecent)

	// Context özeti
	s.AddTool(mcp.Tool{
		Name:        "gorev_context_summary",
		Description: "AI session için optimize edilmiş context özeti. Aktif görev, son görevler, öncelikler ve blokajları gösterir.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.GorevContextSummary)

	// Toplu güncelleme
	s.AddTool(mcp.Tool{
		Name:        "gorev_batch_update",
		Description: "Birden fazla görevi tek seferde günceller. Verimli toplu işlemler için kullanılır.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"updates": map[string]interface{}{
					"type":        "array",
					"description": "Güncelleme listesi. Her öğe {id: string, updates: object} formatında olmalı.",
				},
			},
			Required: []string{"updates"},
		},
	}, tr.handlers.GorevBatchUpdate)

	// Doğal dil sorgusu
	s.AddTool(mcp.Tool{
		Name:        "gorev_nlp_query",
		Description: "Doğal dil ile görev arama. Örn: 'yüksek öncelikli', 'bugün üzerinde çalıştığım', 'etiket:bug'",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Doğal dil sorgusu.",
				},
			},
			Required: []string{"query"},
		},
	}, tr.handlers.GorevNLPQuery)
}

// registerFileWatcherTools registers file watcher tools
func (tr *ToolRegistry) registerFileWatcherTools(s *server.MCPServer) {
	// Dosya izleme ekle
	s.AddTool(mcp.Tool{
		Name:        "gorev_file_watch_add",
		Description: "Bir göreve dosya yolu ekler ve otomatik durum güncelleme için izlemeye başlar.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin ID'si.",
				},
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": "İzlenecek dosya yolu.",
				},
			},
			Required: []string{"task_id", "file_path"},
		},
	}, tr.handlers.GorevFileWatchAdd)

	// Dosya izleme kaldır
	s.AddTool(mcp.Tool{
		Name:        "gorev_file_watch_remove",
		Description: "Bir görevden dosya yolu kaldırır ve izlemeyi durdurur.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin ID'si.",
				},
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": "Kaldırılacak dosya yolu.",
				},
			},
			Required: []string{"task_id", "file_path"},
		},
	}, tr.handlers.GorevFileWatchRemove)

	// Dosya izleme listesi
	s.AddTool(mcp.Tool{
		Name:        "gorev_file_watch_list",
		Description: "Bir görevin izlenen dosya yollarını listeler.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": "Görevin ID'si.",
				},
			},
			Required: []string{"task_id"},
		},
	}, tr.handlers.GorevFileWatchList)

	// Dosya izleme istatistikleri
	s.AddTool(mcp.Tool{
		Name:        "gorev_file_watch_stats",
		Description: "Dosya izleme sisteminin istatistiklerini gösterir.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.GorevFileWatchStats)
}

// registerAdvancedTools registers advanced and hierarchy tools
func (tr *ToolRegistry) registerAdvancedTools(s *server.MCPServer) {
	// Alt görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "gorev_altgorev_olustur",
		Description: "Mevcut bir görev altında alt görev oluşturur. Alt görev üst görevin projesini devralır.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"parent_id": map[string]interface{}{
					"type":        "string",
					"description": "Üst görevin ID'si.",
				},
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": "Alt görevin başlığı.",
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": "Alt görevin açıklaması.",
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": "Alt görevin öncelik seviyesi.",
					"enum":        []string{"dusuk", "orta", "yuksek"},
				},
				"son_tarih": map[string]interface{}{
					"type":        "string",
					"description": "Alt görevin son tarihi (YYYY-MM-DD).",
				},
				"etiketler": map[string]interface{}{
					"type":        "string",
					"description": "Virgülle ayrılmış etiket listesi.",
				},
			},
			Required: []string{"parent_id", "baslik"},
		},
	}, tr.handlers.GorevAltGorevOlustur)

	// Üst görev değiştir
	s.AddTool(mcp.Tool{
		Name:        "gorev_ust_degistir",
		Description: "Bir görevin üst görevini değiştirir veya kök seviyeye taşır.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"gorev_id": map[string]interface{}{
					"type":        "string",
					"description": "Taşınacak görevin ID'si.",
				},
				"yeni_parent_id": map[string]interface{}{
					"type":        "string",
					"description": "Yeni üst görevin ID'si. Boş bırakılırsa kök seviyeye taşınır.",
				},
			},
			Required: []string{"gorev_id"},
		},
	}, tr.handlers.GorevUstDegistir)

	// Hiyerarşi göster
	s.AddTool(mcp.Tool{
		Name:        "gorev_hiyerarsi_goster",
		Description: "Bir görevin tüm hiyerarşisini (üst görevler ve alt görevler) gösterir.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"gorev_id": map[string]interface{}{
					"type":        "string",
					"description": "Hiyerarşisi gösterilecek görevin ID'si.",
				},
			},
			Required: []string{"gorev_id"},
		},
	}, tr.handlers.GorevHiyerarsiGoster)

	// Bağımlılık ekle
	s.AddTool(mcp.Tool{
		Name:        "gorev_bagimlilik_ekle",
		Description: "İki görev arasında bağımlılık ilişkisi kurar.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"kaynak_id": map[string]interface{}{
					"type":        "string",
					"description": "Bağımlı olan görevin ID'si (bekleyecek olan).",
				},
				"hedef_id": map[string]interface{}{
					"type":        "string",
					"description": "Bağımlılık hedefinin ID'si (önce tamamlanması gereken).",
				},
				"baglanti_tipi": map[string]interface{}{
					"type":        "string",
					"description": "Bağımlılık türü.",
					"enum":        []string{"blocker", "depends_on"},
				},
			},
			Required: []string{"kaynak_id", "hedef_id", "baglanti_tipi"},
		},
	}, tr.handlers.GorevBagimlilikEkle)

	// Özet göster
	s.AddTool(mcp.Tool{
		Name:        "ozet_goster",
		Description: "Sistemin genel durumu, proje ve görev istatistiklerini gösterir.",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.OzetGoster)
}
