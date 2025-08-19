package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
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
	// Note: gorev_olustur was deprecated in v0.10.0 and removed in v0.11.1
	// Use templateden_gorev_olustur instead

	// Görev listele
	s.AddTool(mcp.Tool{
		Name:        "gorev_listele",
		Description: i18n.T("tools.descriptions.gorev_listele", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"durum": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("durum_filter"),
					"enum":        []string{constants.TaskStatusPending, constants.TaskStatusInProgress, constants.TaskStatusCompleted},
				},
				"sirala": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("sirala"),
					"enum":        []string{"son_tarih_asc", "son_tarih_desc"},
				},
				"filtre": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("filtre"),
					"enum":        []string{"acil", "gecmis"},
				},
				"etiket": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("etiket"),
				},
				"tum_projeler": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.TParam("tum_projeler"),
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": i18n.TParam("limit"),
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": i18n.TParam("offset"),
				},
			},
		},
	}, tr.handlers.GorevListele)

	// Görev detay
	s.AddTool(mcp.Tool{
		Name:        "gorev_detay",
		Description: i18n.T("tools.descriptions.gorev_detay", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("id_field"),
				},
			},
			Required: []string{"id"},
		},
	}, tr.handlers.GorevDetay)

	// Görev güncelle
	s.AddTool(mcp.Tool{
		Name:        "gorev_guncelle",
		Description: i18n.T("tools.descriptions.gorev_guncelle", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("id_field"),
				},
				"durum": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("durum"),
					"enum":        constants.GetValidTaskStatuses(),
				},
			},
			Required: []string{"id", "durum"},
		},
	}, tr.handlers.GorevGuncelle)

	// Görev düzenle
	s.AddTool(mcp.Tool{
		Name:        "gorev_duzenle",
		Description: i18n.T("tools.descriptions.gorev_duzenle", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("id_field"),
				},
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("baslik"),
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("aciklama"),
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("oncelik"),
					"enum":        constants.GetValidPriorities(),
				},
				"son_tarih": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("son_tarih"),
				},
			},
			Required: []string{"id"},
		},
	}, tr.handlers.GorevDuzenle)

	// Görev sil
	s.AddTool(mcp.Tool{
		Name:        "gorev_sil",
		Description: i18n.T("tools.descriptions.gorev_sil", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "delete"),
				},
				"onay": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.TParam("onay"),
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
		Description: i18n.T("tools.descriptions.proje_olustur", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"isim": map[string]interface{}{
					"type":        "string",
					"description": i18n.TProjectField("name"),
				},
				"tanim": map[string]interface{}{
					"type":        "string",
					"description": i18n.TProjectField("description"),
				},
			},
			Required: []string{"isim", "tanim"},
		},
	}, tr.handlers.ProjeOlustur)

	// Proje listele
	s.AddTool(mcp.Tool{
		Name:        "proje_listele",
		Description: i18n.T("tools.descriptions.proje_listele", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.ProjeListele)

	// Proje görevleri
	s.AddTool(mcp.Tool{
		Name:        "proje_gorevleri",
		Description: i18n.T("tools.descriptions.proje_gorevleri", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("project", "unique"),
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": i18n.TTaskCount("return_max", "fifty"),
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": i18n.TTaskCount("skip", "zero"),
				},
			},
			Required: []string{"proje_id"},
		},
	}, tr.handlers.ProjeGorevleri)

	// Aktif proje ayarla
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_ayarla",
		Description: i18n.T("tools.descriptions.aktif_proje_ayarla", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"proje_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("project", "active"),
				},
			},
			Required: []string{"proje_id"},
		},
	}, tr.handlers.AktifProjeAyarla)

	// Aktif proje göster
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_goster",
		Description: i18n.T("tools.descriptions.aktif_proje_goster", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.AktifProjeGoster)

	// Aktif proje kaldır
	s.AddTool(mcp.Tool{
		Name:        "aktif_proje_kaldir",
		Description: i18n.T("tools.descriptions.aktif_proje_kaldir", nil),
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
		Description: i18n.T("tools.descriptions.template_listele", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"kategori": map[string]interface{}{
					"type":        "string",
					"description": i18n.TTemplate("filter"),
				},
			},
		},
	}, tr.handlers.TemplateListele)

	// Template'den görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "templateden_gorev_olustur",
		Description: i18n.T("tools.descriptions.templateden_gorev_olustur", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				constants.ParamTemplateID: map[string]interface{}{
					"type":        "string",
					"description": i18n.TTemplate("template_id"),
				},
				constants.ParamDegerler: map[string]interface{}{
					"type":        "object",
					"description": i18n.TTemplate("fields"),
				},
			},
			Required: []string{constants.ParamTemplateID, constants.ParamDegerler},
		},
	}, tr.handlers.TemplatedenGorevOlustur)
}

// registerAIContextTools registers AI context management tools
func (tr *ToolRegistry) registerAIContextTools(s *server.MCPServer) {
	// AI aktif görev ayarla
	s.AddTool(mcp.Tool{
		Name:        "gorev_set_active",
		Description: i18n.T("tools.descriptions.gorev_set_active", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "active"),
				},
			},
			Required: []string{"task_id"},
		},
	}, tr.handlers.GorevSetActive)

	// AI aktif görevi getir
	s.AddTool(mcp.Tool{
		Name:        "gorev_get_active",
		Description: i18n.T("tools.descriptions.gorev_get_active", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.GorevGetActive)

	// Son görevleri getir
	s.AddTool(mcp.Tool{
		Name:        "gorev_recent",
		Description: i18n.T("tools.descriptions.gorev_recent", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":        "number",
					"description": i18n.TTaskCount("return", "five"),
				},
			},
		},
	}, tr.handlers.GorevRecent)

	// Context özeti
	s.AddTool(mcp.Tool{
		Name:        "gorev_context_summary",
		Description: i18n.T("tools.descriptions.gorev_context_summary", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.GorevContextSummary)

	// Toplu güncelleme
	s.AddTool(mcp.Tool{
		Name:        "gorev_batch_update",
		Description: i18n.T("tools.descriptions.gorev_batch_update", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"updates": map[string]interface{}{
					"type":        "array",
					"description": i18n.TBatch("updates"),
				},
			},
			Required: []string{"updates"},
		},
	}, tr.handlers.GorevBatchUpdate)

	// Doğal dil sorgusu
	s.AddTool(mcp.Tool{
		Name:        "gorev_nlp_query",
		Description: i18n.T("tools.descriptions.gorev_nlp_query", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": i18n.TBatch("query"),
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
		Description: i18n.T("tools.descriptions.gorev_file_watch_add", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "simple"),
				},
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFilePath("watch"),
				},
			},
			Required: []string{"task_id", "file_path"},
		},
	}, tr.handlers.GorevFileWatchAdd)

	// Dosya izleme kaldır
	s.AddTool(mcp.Tool{
		Name:        "gorev_file_watch_remove",
		Description: i18n.T("tools.descriptions.gorev_file_watch_remove", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "simple"),
				},
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFilePath("remove"),
				},
			},
			Required: []string{"task_id", "file_path"},
		},
	}, tr.handlers.GorevFileWatchRemove)

	// Dosya izleme listesi
	s.AddTool(mcp.Tool{
		Name:        "gorev_file_watch_list",
		Description: i18n.T("tools.descriptions.gorev_file_watch_list", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "simple"),
				},
			},
			Required: []string{"task_id"},
		},
	}, tr.handlers.GorevFileWatchList)

	// Dosya izleme istatistikleri
	s.AddTool(mcp.Tool{
		Name:        "gorev_file_watch_stats",
		Description: i18n.T("tools.descriptions.gorev_file_watch_stats", nil),
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
		Description: i18n.T("tools.descriptions.gorev_altgorev_olustur", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"parent_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "parent"),
				},
				"baslik": map[string]interface{}{
					"type":        "string",
					"description": i18n.TSubtaskField("title"),
				},
				"aciklama": map[string]interface{}{
					"type":        "string",
					"description": i18n.TSubtaskField("subtask_description"),
				},
				"oncelik": map[string]interface{}{
					"type":        "string",
					"description": i18n.TSubtaskField("priority"),
					"enum":        constants.GetValidPriorities(),
				},
				"son_tarih": map[string]interface{}{
					"type":        "string",
					"description": i18n.TWithFormat(i18n.TSubtaskField("due_date"), "YYYY-MM-DD"),
				},
				"etiketler": map[string]interface{}{
					"type":        "string",
					"description": i18n.TCommaSeparated("etiket"),
				},
			},
			Required: []string{"parent_id", "baslik"},
		},
	}, tr.handlers.GorevAltGorevOlustur)

	// Üst görev değiştir
	s.AddTool(mcp.Tool{
		Name:        "gorev_ust_degistir",
		Description: i18n.T("tools.descriptions.gorev_ust_degistir", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"gorev_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "move"),
				},
				"yeni_parent_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "new_parent"),
				},
			},
			Required: []string{"gorev_id"},
		},
	}, tr.handlers.GorevUstDegistir)

	// Hiyerarşi göster
	s.AddTool(mcp.Tool{
		Name:        "gorev_hiyerarsi_goster",
		Description: i18n.T("tools.descriptions.gorev_hiyerarsi_goster", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"gorev_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "hierarchy"),
				},
			},
			Required: []string{"gorev_id"},
		},
	}, tr.handlers.GorevHiyerarsiGoster)

	// Bağımlılık ekle
	s.AddTool(mcp.Tool{
		Name:        "gorev_bagimlilik_ekle",
		Description: i18n.T("tools.descriptions.gorev_bagimlilik_ekle", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"kaynak_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("task", "dependent"),
				},
				"hedef_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.T("common.fields.target_id", nil),
				},
				"baglanti_tipi": map[string]interface{}{
					"type":        "string",
					"description": i18n.TBatch("dependency_type"),
					"enum":        constants.GetValidDependencyTypes(),
				},
			},
			Required: []string{"kaynak_id", "hedef_id", "baglanti_tipi"},
		},
	}, tr.handlers.GorevBagimlilikEkle)

	// Özet göster
	s.AddTool(mcp.Tool{
		Name:        "ozet_goster",
		Description: i18n.T("tools.descriptions.ozet_goster", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.OzetGoster)
}
