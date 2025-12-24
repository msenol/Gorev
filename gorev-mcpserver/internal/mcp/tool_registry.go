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
	tr.registerUnifiedTools(s) // 8 unified tools replacing 27 individual tools
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
				"status": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "durum_filter"),
					"enum":        []string{constants.TaskStatusPending, constants.TaskStatusInProgress, constants.TaskStatusCompleted},
				},
				"sort": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "sirala"),
					"enum":        []string{"due_date_asc", "due_date_desc"},
				},
				"filter": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "filtre"),
					"enum":        []string{"urgent", "overdue"},
				},
				"tag": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "etiket"),
				},
				"all_projects": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.TParam("tr", "tum_projeler"),
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": i18n.TParam("tr", "limit"),
				},
				"offset": map[string]interface{}{
					"type":        "number",
					"description": i18n.TParam("tr", "offset"),
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
					"description": i18n.TParam("tr", "id_field"),
				},
			},
			Required: []string{"id"},
		},
	}, tr.handlers.GorevDetay)

	// Görev güncelle (status update)
	s.AddTool(mcp.Tool{
		Name:        "gorev_guncelle",
		Description: i18n.T("tools.descriptions.gorev_guncelle", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "id_field"),
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "durum"),
					"enum":        constants.GetValidTaskStatuses(),
				},
				"priority": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "oncelik"),
					"enum":        constants.GetValidPriorities(),
				},
			},
			Required: []string{"id"},
		},
	}, tr.handlers.GorevGuncelle)

	// Görev düzenle (content edit)
	s.AddTool(mcp.Tool{
		Name:        "gorev_duzenle",
		Description: i18n.T("tools.descriptions.gorev_duzenle", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "id_field"),
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "baslik"),
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "aciklama"),
				},
				"priority": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "oncelik"),
					"enum":        constants.GetValidPriorities(),
				},
				"due_date": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "son_tarih"),
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
					"description": i18n.TFieldID("tr", "task", "delete"),
				},
				"confirm": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.TParam("tr", "onay"),
				},
			},
			Required: []string{"id", "confirm"},
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
				"name": map[string]interface{}{
					"type":        "string",
					"description": i18n.TProjectField("tr", "name"),
				},
				"definition": map[string]interface{}{
					"type":        "string",
					"description": i18n.TProjectField("tr", "description"),
				},
			},
			Required: []string{"name"},
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
				"project_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("tr", "project", "unique"),
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
			Required: []string{"project_id"},
		},
	}, tr.handlers.ProjeGorevleri)

	// Active project tools replaced by unified "aktif_proje" tool with actions: set|get|clear
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
				"category": map[string]interface{}{
					"type":        "string",
					"description": i18n.TTemplate("tr", "filter"),
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
					"description": i18n.TTemplate("tr", "template_id"),
				},
				constants.ParamValues: map[string]interface{}{
					"type":        "object",
					"description": i18n.TTemplate("tr", "fields"),
				},
			},
			Required: []string{constants.ParamTemplateID, constants.ParamValues},
		},
	}, tr.handlers.TemplatedenGorevOlustur)
}

// registerAIContextTools was removed - AI context tools now use unified handlers
// All AI context tools replaced by unified handlers:
// - gorev_set_active, gorev_get_active, gorev_recent, gorev_context_summary → gorev_context
// - gorev_batch_update → gorev_bulk (operation: update)
// - gorev_nlp_query → gorev_search (mode: nlp)

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
					"description": i18n.TFieldID("tr", "task", "simple"),
				},
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFilePath("tr", "watch"),
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
					"description": i18n.TFieldID("tr", "task", "simple"),
				},
				"file_path": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFilePath("tr", "remove"),
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
					"description": i18n.TFieldID("tr", "task", "simple"),
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

// registerSearchTools was removed - search tools now use unified handlers
// All search tools replaced by unified handlers:
// - gorev_search_advanced, gorev_nlp_query, gorev_search_history → gorev_search (modes: advanced|nlp|history)
// - gorev_filter_profile_save, gorev_filter_profile_load, gorev_filter_profile_list, gorev_filter_profile_delete → gorev_filter_profile

// registerAdvancedTools registers advanced and hierarchy tools
func (tr *ToolRegistry) registerAdvancedTools(s *server.MCPServer) {
	if tr.handlers.debug {
		slog.Debug("Registering advanced tools")
	}

	// Hierarchy tools replaced by unified "gorev_hierarchy" tool with actions: create_subtask|change_parent|show

	// Bağımlılık ekle
	s.AddTool(mcp.Tool{
		Name:        "gorev_bagimlilik_ekle",
		Description: i18n.T("tools.descriptions.gorev_bagimlilik_ekle", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"source_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TFieldID("tr", "task", "dependent"),
				},
				"target_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.T("common.fields.target_id", nil),
				},
				"connection_type": map[string]interface{}{
					"type":        "string",
					"description": i18n.TBatch("tr", "dependency_type"),
					"enum":        constants.GetValidDependencyTypes(),
				},
			},
			Required: []string{"source_id", "target_id", "connection_type"},
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

	// Gorev Suggestions - AI-powered task suggestions
	s.AddTool(mcp.Tool{
		Name:        "gorev_suggestions",
		Description: i18n.T("tools.descriptions.gorev_suggestions", nil),
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]interface{}{},
		},
	}, tr.handlers.GorevSuggestions)

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
					"type":        "string",
					"description": i18n.T("tools.params.export.output_path", nil),
				},
				"format": map[string]interface{}{
					"type":        "string",
					"description": i18n.T("tools.params.export.format", nil),
					"enum":        []string{"json", "csv"},
					"default":     "json",
				},
				"include_completed": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.T("tools.params.export.include_completed", nil),
					"default":     true,
				},
				"include_dependencies": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.T("tools.params.export.include_dependencies", nil),
					"default":     true,
				},
				"include_templates": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.T("tools.params.export.include_templates", nil),
					"default":     false,
				},
				"include_ai_context": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.T("tools.params.export.include_ai_context", nil),
					"default":     false,
				},
				"project_filter": map[string]interface{}{
					"type":        "array",
					"description": i18n.T("tools.params.export.project_filter", nil),
					"items": map[string]interface{}{
						"type": "string",
					},
				},
				"date_range": map[string]interface{}{
					"type":        "object",
					"description": i18n.T("tools.params.export.date_range", nil),
					"properties": map[string]interface{}{
						"from": map[string]interface{}{
							"type":        "string",
							"description": i18n.T("tools.params.export.date_from", nil),
							"format":      "date-time",
						},
						"to": map[string]interface{}{
							"type":        "string",
							"description": i18n.T("tools.params.export.date_to", nil),
							"format":      "date-time",
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
					"type":        "string",
					"description": i18n.T("tools.params.import.file_path", nil),
				},
				"import_mode": map[string]interface{}{
					"type":        "string",
					"description": i18n.T("tools.params.import.import_mode", nil),
					"enum":        []string{"merge", "replace"},
					"default":     "merge",
				},
				"conflict_resolution": map[string]interface{}{
					"type":        "string",
					"description": i18n.T("tools.params.import.conflict_resolution", nil),
					"enum":        []string{"skip", "overwrite", "prompt"},
					"default":     "skip",
				},
				"preserve_ids": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.T("tools.params.import.preserve_ids", nil),
					"default":     false,
				},
				"dry_run": map[string]interface{}{
					"type":        "boolean",
					"description": i18n.T("tools.params.import.dry_run", nil),
					"default":     false,
				},
				"project_mapping": map[string]interface{}{
					"type":        "object",
					"description": i18n.T("tools.params.import.project_mapping", nil),
					"additionalProperties": map[string]interface{}{
						"type": "string",
					},
				},
			},
			Required: []string{"file_path"},
		},
	}, tr.handlers.GorevImport)

	// IDE Management tools replaced by unified "gorev_ide" tool with actions: detect|install|uninstall|status|update
}

// registerUnifiedTools registers optimized unified handlers
// This replaces 27 individual tools with 7 unified handlers (37% reduction)
func (tr *ToolRegistry) registerUnifiedTools(s *server.MCPServer) {
	if tr.handlers.debug {
		slog.Debug("Registering unified tools - 7 tools replacing 27 individual tools")
	}

	// ========================================
	// Active Project Management (1 tool replaces 3)
	// ========================================

	s.AddTool(mcp.Tool{
		Name:        "aktif_proje",
		Description: i18n.T("tools.descriptions.aktif_proje", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Operation to perform",
					"enum":        constants.ValidActiveProjectActions,
				},
				"project_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "project_id"),
				},
			},
			Required: []string{"action"},
		},
	}, tr.handlers.AktifProje)

	// ========================================
	// Bulk Operations (1 tool replaces 3)
	// ========================================

	s.AddTool(mcp.Tool{
		Name:        "gorev_bulk",
		Description: i18n.T("tools.descriptions.gorev_bulk", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"operation": map[string]interface{}{
					"type":        "string",
					"description": "Bulk operation type",
					"enum":        constants.ValidBulkOperationActions,
				},
				"ids": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "string"},
					"description": i18n.TParam("tr", "task_ids"),
				},
				"data": map[string]interface{}{
					"type":        "object",
					"description": "Operation data (status, tags, etc.)",
				},
			},
			Required: []string{"operation", "ids"},
		},
	}, tr.handlers.GorevBulk)

	// ========================================
	// Hierarchy Management (1 tool replaces 3)
	// ========================================

	s.AddTool(mcp.Tool{
		Name:        "gorev_hierarchy",
		Description: i18n.T("tools.descriptions.gorev_hierarchy", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Hierarchy operation",
					"enum":        constants.ValidHierarchyActions,
				},
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "task_id"),
				},
				"parent_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "parent_id"),
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "title"),
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "description"),
				},
			},
			Required: []string{"action"},
		},
	}, tr.handlers.GorevHierarchy)

	// ========================================
	// Filter Profile Management (1 tool replaces 4)
	// ========================================

	s.AddTool(mcp.Tool{
		Name:        "gorev_filter_profile",
		Description: i18n.T("tools.descriptions.gorev_filter_profile", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "Filter profile action",
					"enum":        constants.ValidFilterProfileActions,
				},
				"profile_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "profile_id"),
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "profile_name"),
				},
				"filters": map[string]interface{}{
					"type":        "object",
					"description": i18n.TParam("tr", "filters"),
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "description"),
				},
			},
			Required: []string{"action"},
		},
	}, tr.handlers.GorevFilterProfile)

	// ========================================
	// IDE Management (1 tool replaces 5)
	// ========================================

	s.AddTool(mcp.Tool{
		Name:        "gorev_ide",
		Description: i18n.T("tools.descriptions.gorev_ide", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "IDE management action",
					"enum":        constants.ValidIDEActions,
				},
				"ide_type": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "ide_type"),
					"enum":        []string{"vscode", "cursor", "windsurf", "all"},
				},
			},
			Required: []string{"action"},
		},
	}, tr.handlers.IDEManage)

	// ========================================
	// AI Context Management (1 tool replaces 4)
	// ========================================

	s.AddTool(mcp.Tool{
		Name:        "gorev_context",
		Description: i18n.T("tools.descriptions.gorev_context", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"action": map[string]interface{}{
					"type":        "string",
					"description": "AI context action",
					"enum":        constants.ValidContextActions,
				},
				"task_id": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "task_id"),
				},
				"limit": map[string]interface{}{
					"type":        "number",
					"description": i18n.TParam("tr", "limit"),
				},
			},
			Required: []string{"action"},
		},
	}, tr.handlers.GorevContext)

	// ========================================
	// Search Operations (1 tool replaces 3)
	// ========================================

	s.AddTool(mcp.Tool{
		Name:        "gorev_search",
		Description: i18n.T("tools.descriptions.gorev_search", nil),
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"mode": map[string]interface{}{
					"type":        "string",
					"description": "Search mode",
					"enum":        constants.ValidSearchModes,
				},
				"query": map[string]interface{}{
					"type":        "string",
					"description": i18n.TParam("tr", "query"),
				},
				"filters": map[string]interface{}{
					"type":        "object",
					"description": i18n.TParam("tr", "filters"),
				},
			},
			Required: []string{"mode"},
		},
	}, tr.handlers.GorevSearch)
}
