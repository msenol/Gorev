package mcp

import (
	"log/slog"

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
	tr.registerSearchTools(s)
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
					"type":     "string",
					"aciklama": i18n.TParam("tr", "durum_filter"),
					"enum":     []string{constants.TaskStatusPending, constants.TaskStatusInProgress, constants.TaskStatusCompleted},
				},
				"sirala": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "sirala"),
					"enum":     []string{"son_tarih_asc", "son_tarih_desc"},
				},
				"filtre": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "filtre"),
					"enum":     []string{"acil", "gecmis"},
				},
				"etiket": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "etiket"),
				},
				"tum_projeler": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.TParam("tr", "tum_projeler"),
				},
				"limit": map[string]interface{}{
					"type":     "number",
					"aciklama": i18n.TParam("tr", "limit"),
				},
				"offset": map[string]interface{}{
					"type":     "number",
					"aciklama": i18n.TParam("tr", "offset"),
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
					"type":     "string",
					"aciklama": i18n.TParam("tr", "id_field"),
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
					"type":     "string",
					"aciklama": i18n.TParam("tr", "id_field"),
				},
				"durum": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "durum"),
					"enum":     constants.GetValidTaskStatuses(),
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
					"type":     "string",
					"aciklama": i18n.TParam("tr", "id_field"),
				},
				"baslik": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "baslik"),
				},
				"aciklama": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "aciklama"),
				},
				"oncelik": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "oncelik"),
					"enum":     constants.GetValidPriorities(),
				},
				"son_tarih": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TParam("tr", "son_tarih"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "delete"),
				},
				"onay": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.TParam("tr", "onay"),
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
					"type":     "string",
					"aciklama": i18n.TProjectField("tr", "name"),
				},
				"tanim": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TProjectField("tr", "description"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "project", "unique"),
				},
				"limit": map[string]interface{}{
					"type":     "number",
					"aciklama": i18n.TTaskCount("return_max", "fifty"),
				},
				"offset": map[string]interface{}{
					"type":     "number",
					"aciklama": i18n.TTaskCount("skip", "zero"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "project", "active"),
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
					"type":     "string",
					"aciklama": i18n.TTemplate("tr", "filter"),
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
					"type":     "string",
					"aciklama": i18n.TTemplate("tr", "template_id"),
				},
				constants.ParamValues: map[string]interface{}{
					"type":     "object",
					"aciklama": i18n.TTemplate("tr", "fields"),
				},
			},
			Required: []string{constants.ParamTemplateID, constants.ParamValues},
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "active"),
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
					"type":     "number",
					"aciklama": i18n.TTaskCount("return", "five"),
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
					"type":     "array",
					"aciklama": i18n.TBatch("tr", "updates"),
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"id": map[string]interface{}{
								"type":     "string",
								"aciklama": "Task ID to update",
							},
							"baslik": map[string]interface{}{
								"type":     "string",
								"aciklama": "New task title",
							},
							"aciklama": map[string]interface{}{
								"type":     "string",
								"aciklama": "New task description",
							},
							"durum": map[string]interface{}{
								"type":     "string",
								"aciklama": "New task status",
								"enum":     []string{"beklemede", "devam_ediyor", "tamamlandi", "iptal"},
							},
							"oncelik": map[string]interface{}{
								"type":     "string",
								"aciklama": "New task priority",
								"enum":     []string{"dusuk", "orta", "yuksek"},
							},
							"son_tarih": map[string]interface{}{
								"type":     "string",
								"aciklama": "New due date (YYYY-MM-DD format)",
							},
							"proje_id": map[string]interface{}{
								"type":     "string",
								"aciklama": "New project ID",
							},
						},
						"required": []string{"id"},
					},
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
					"type":     "string",
					"aciklama": i18n.TBatch("tr", "query"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "simple"),
				},
				"file_path": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TFilePath("tr", "watch"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "simple"),
				},
				"file_path": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TFilePath("tr", "remove"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "simple"),
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

// registerSearchTools registers advanced search and filter tools
func (tr *ToolRegistry) registerSearchTools(s *server.MCPServer) {
	if tr.handlers.debug {
		slog.Debug("Registering search tools")
	}

	// Advanced search with FTS5
	s.AddTool(mcp.Tool{
		Name:        "gorev_search_advanced",
		Description: i18n.T("tools.descriptions.gorev_search_advanced", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"query": map[string]interface{}{
					"type":     "string",
					"aciklama": "Search query for FTS5 full-text search",
				},
				"filters": map[string]interface{}{
					"type":     "object",
					"aciklama": "Filter conditions (durum, oncelik, proje_id, etc.)",
				},
				"use_fuzzy_search": map[string]interface{}{
					"type":     "boolean",
					"aciklama": "Enable fuzzy search for partial matches",
					"default":  true,
				},
				"fuzzy_threshold": map[string]interface{}{
					"type":     "number",
					"aciklama": "Fuzzy search similarity threshold (0.0-1.0)",
					"default":  0.6,
				},
				"max_results": map[string]interface{}{
					"type":     "integer",
					"aciklama": "Maximum number of results to return",
					"default":  50,
				},
				"sort_by": map[string]interface{}{
					"type":     "string",
					"aciklama": "Sort field (relevance, created, updated, due_date, priority)",
					"default":  "relevance",
				},
				"sort_direction": map[string]interface{}{
					"type":     "string",
					"aciklama": "Sort direction (asc, desc)",
					"default":  "desc",
				},
				"include_completed": map[string]interface{}{
					"type":     "boolean",
					"aciklama": "Include completed tasks in results",
					"default":  false,
				},
				"search_fields": map[string]interface{}{
					"type":     "array",
					"aciklama": "Fields to search (baslik, aciklama, etiketler, proje_adi)",
					"items": map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}, tr.handlers.GorevSearchAdvanced)

	// Save filter profile
	s.AddTool(mcp.Tool{
		Name:        "gorev_filter_profile_save",
		Description: i18n.T("tools.descriptions.gorev_filter_profile_save", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"isim": map[string]interface{}{
					"type":     "string",
					"aciklama": "Profile name",
				},
				"aciklama": map[string]interface{}{
					"type":     "string",
					"aciklama": "Profile description",
				},
				"filters": map[string]interface{}{
					"type":     "object",
					"aciklama": "Filter configuration to save",
				},
				"search_query": map[string]interface{}{
					"type":     "string",
					"aciklama": "Search query to save with profile",
				},
				"is_default": map[string]interface{}{
					"type":     "boolean",
					"aciklama": "Mark as default profile",
					"default":  false,
				},
			},
			Required: []string{"name"},
		},
	}, tr.handlers.GorevFilterProfileSave)

	// Load filter profile
	s.AddTool(mcp.Tool{
		Name:        "gorev_filter_profile_load",
		Description: i18n.T("tools.descriptions.gorev_filter_profile_load", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"profile_id": map[string]interface{}{
					"type":     "integer",
					"aciklama": "Profile ID to load",
				},
				"profile_name": map[string]interface{}{
					"type":     "string",
					"aciklama": "Profile name to load",
				},
			},
		},
	}, tr.handlers.GorevFilterProfileLoad)

	// List filter profiles
	s.AddTool(mcp.Tool{
		Name:        "gorev_filter_profile_list",
		Description: i18n.T("tools.descriptions.gorev_filter_profile_list", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"defaults_only": map[string]interface{}{
					"type":     "boolean",
					"aciklama": "Only return default profiles",
					"default":  false,
				},
			},
		},
	}, tr.handlers.GorevFilterProfileList)

	// Delete filter profile
	s.AddTool(mcp.Tool{
		Name:        "gorev_filter_profile_delete",
		Description: i18n.T("tools.descriptions.gorev_filter_profile_delete", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"profile_id": map[string]interface{}{
					"type":     "integer",
					"aciklama": "Profile ID to delete",
				},
			},
			Required: []string{"profile_id"},
		},
	}, tr.handlers.GorevFilterProfileDelete)

	// Search history
	s.AddTool(mcp.Tool{
		Name:        "gorev_search_history",
		Description: i18n.T("tools.descriptions.gorev_search_history", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"limit": map[string]interface{}{
					"type":     "integer",
					"aciklama": "Maximum number of history entries to return",
					"default":  20,
				},
			},
		},
	}, tr.handlers.GorevSearchHistory)
}

// registerAdvancedTools registers advanced and hierarchy tools
func (tr *ToolRegistry) registerAdvancedTools(s *server.MCPServer) {
	if tr.handlers.debug {
		slog.Debug("Registering advanced tools")
	}
	// Alt görev oluştur
	s.AddTool(mcp.Tool{
		Name:        "gorev_altgorev_olustur",
		Description: i18n.T("tools.descriptions.gorev_altgorev_olustur", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"parent_id": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "parent"),
				},
				"baslik": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TSubtaskField("tr", "title"),
				},
				"aciklama": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TSubtaskField("tr", "subtask_description"),
				},
				"oncelik": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TSubtaskField("tr", "priority"),
					"enum":     constants.GetValidPriorities(),
				},
				"son_tarih": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TWithFormat("tr", i18n.TSubtaskField("tr", "due_date"), "YYYY-MM-DD"),
				},
				"tags": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TCommaSeparated("tr", "etiket"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "move"),
				},
				"yeni_parent_id": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "new_parent"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "hierarchy"),
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
					"type":     "string",
					"aciklama": i18n.TFieldID("tr", "task", "dependent"),
				},
				"hedef_id": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("common.fields.target_id", nil),
				},
				"baglanti_tipi": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.TBatch("tr", "dependency_type"),
					"enum":     constants.GetValidDependencyTypes(),
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

	// Gorev Export - Data export tool
	if tr.handlers.debug {
		slog.Debug("Registering gorev_export tool")
	}
	s.AddTool(mcp.Tool{
		Name:        "gorev_export",
		Description: i18n.T("tools.descriptions.gorev_export", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"output_path": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.export.output_path", nil),
				},
				"format": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.export.format", nil),
					"enum":     []string{"json", "csv"},
					"default":  "json",
				},
				"include_completed": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.T("tools.params.export.include_completed", nil),
					"default":  true,
				},
				"include_dependencies": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.T("tools.params.export.include_dependencies", nil),
					"default":  true,
				},
				"include_templates": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.T("tools.params.export.include_templates", nil),
					"default":  false,
				},
				"include_ai_context": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.T("tools.params.export.include_ai_context", nil),
					"default":  false,
				},
				"project_filter": map[string]interface{}{
					"type":     "array",
					"aciklama": i18n.T("tools.params.export.project_filter", nil),
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"date_range": map[string]interface{}{
					"type":     "object",
					"aciklama": i18n.T("tools.params.export.date_range", nil),
					"properties": map[string]interface{}{
						"from": map[string]interface{}{
							"type":     "string",
							"aciklama": i18n.T("tools.params.export.date_from", nil),
							"format":   "date-time",
						},
						"to": map[string]interface{}{
							"type":     "string",
							"aciklama": i18n.T("tools.params.export.date_to", nil),
							"format":   "date-time",
						},
					},
				},
			},
			Required: []string{"output_path"},
		},
	}, tr.handlers.GorevExport)

	// Gorev Import - Data import tool
	if tr.handlers.debug {
		slog.Debug("Registering gorev_import tool")
	}
	s.AddTool(mcp.Tool{
		Name:        "gorev_import",
		Description: i18n.T("tools.descriptions.gorev_import", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"file_path": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.import.file_path", nil),
				},
				"import_mode": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.import.import_mode", nil),
					"enum":     []string{"merge", "replace"},
					"default":  "merge",
				},
				"conflict_resolution": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.import.conflict_resolution", nil),
					"enum":     []string{"skip", "overwrite", "prompt"},
					"default":  "skip",
				},
				"preserve_ids": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.T("tools.params.import.preserve_ids", nil),
					"default":  false,
				},
				"dry_run": map[string]interface{}{
					"type":     "boolean",
					"aciklama": i18n.T("tools.params.import.dry_run", nil),
					"default":  false,
				},
				"project_mapping": map[string]interface{}{
					"type":     "object",
					"aciklama": i18n.T("tools.params.import.project_mapping", nil),
					"additionalProperties": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"file_path"},
		},
	}, tr.handlers.GorevImport)

	// ========================================
	// IDE Management Tools
	// ========================================

	// IDE Detect - Detect installed IDEs
	if tr.handlers.debug {
		slog.Debug("Registering gorev_ide_detect tool")
	}
	s.AddTool(mcp.Tool{
		Name:        "gorev_ide_detect",
		Description: i18n.T("tools.descriptions.ide_detect", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.IDEDetect)

	// IDE Install Extension - Install Gorev extension to IDEs
	if tr.handlers.debug {
		slog.Debug("Registering gorev_ide_install tool")
	}
	s.AddTool(mcp.Tool{
		Name:        "gorev_ide_install",
		Description: i18n.T("tools.descriptions.ide_install", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"ide_type": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.ide.ide_type", nil),
					"enum":     []string{"vscode", "cursor", "windsurf", "all"},
				},
			},
			Required: []string{"ide_type"},
		},
	}, tr.handlers.IDEInstallExtension)

	// IDE Uninstall Extension - Remove Gorev extension from IDEs
	if tr.handlers.debug {
		slog.Debug("Registering gorev_ide_uninstall tool")
	}
	s.AddTool(mcp.Tool{
		Name:        "gorev_ide_uninstall",
		Description: i18n.T("tools.descriptions.ide_uninstall", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"ide_type": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.ide.ide_type", nil),
					"enum":     []string{"vscode", "cursor", "windsurf"},
				},
				"extension_id": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.ide.extension_id", nil),
					"default":  "mehmetsenol.gorev-vscode",
				},
			},
			Required: []string{"ide_type"},
		},
	}, tr.handlers.IDEUninstallExtension)

	// IDE Extension Status - Check extension installation status
	if tr.handlers.debug {
		slog.Debug("Registering gorev_ide_status tool")
	}
	s.AddTool(mcp.Tool{
		Name:        "gorev_ide_status",
		Description: i18n.T("tools.descriptions.ide_status", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.IDEExtensionStatus)

	// IDE Update Extension - Update Gorev extension to latest version
	if tr.handlers.debug {
		slog.Debug("Registering gorev_ide_update tool")
	}
	s.AddTool(mcp.Tool{
		Name:        "gorev_ide_update",
		Description: i18n.T("tools.descriptions.ide_update", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"ide_type": map[string]interface{}{
					"type":     "string",
					"aciklama": i18n.T("tools.params.ide.ide_type", nil),
					"enum":     []string{"vscode", "cursor", "windsurf", "all"},
				},
			},
			Required: []string{"ide_type"},
		},
	}, tr.handlers.IDEUpdateExtension)
}
