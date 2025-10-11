package constants

// MCP tool parameter names to eliminate hardcoded strings throughout the codebase
const (
	// Core task parameters
	ParamID          = "id"
	ParamTitle       = "title"
	ParamDescription = "description"
	ParamPriority    = "priority"
	ParamStatus      = "status"
	ParamProjeID     = "proje_id"
	ParamDueDate     = "due_date"
	ParamTags        = "tags"
	ParamParentID    = "parent_id"
	ParamGorevID     = "gorev_id"

	// Task management parameters
	ParamAllProjects = "all_projects"
	ParamOrderBy     = "order_by"
	ParamSort        = "sort"
	ParamFilter      = "filter"
	ParamTag         = "tag"
	ParamLimit       = "limit"
	ParamOffset      = "offset"
	ParamConfirm     = "confirm"

	// Dependency parameters
	ParamSourceID       = "source_id"
	ParamTargetID       = "target_id"
	ParamConnectionType = "connection_type"
	ParamNewParentID    = "new_parent_id"

	// Project parameters
	ParamName       = "name"
	ParamDefinition = "definition"

	// Template parameters
	ParamTemplateID = "template_id"
	ParamValues     = "values"
	ParamCategory   = "category"

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
	ResponseTemplates   = "templates"
	ResponseSummary     = "summary"
	ResponseCount       = "count"
	ResponseTotal       = "total"
	ResponsePage        = "page"
	ResponseHasNextPage = "has_next_page"
	ResponseHasPrevPage = "has_prev_page"
	ResponseSuccess     = "success"
	ResponseMessage     = "message"
	ResponseError       = "error"
	ResponseData        = "data"
)

// Database field names to eliminate hardcoded strings
const (
	DBFieldID          = "id"
	DBFieldTitle       = "title"
	DBFieldDescription = "description"
	DBFieldStatus      = "status"
	DBFieldPriority    = "priority"
	DBFieldProjeID     = "proje_id"
	DBFieldParentID    = "parent_id"
	DBFieldDueDate     = "due_date"
	DBFieldCreatedAt   = "created_at"
	DBFieldUpdatedAt   = "updated_at"
	DBFieldName        = "name"
	DBFieldDefinition  = "definition"
	DBFieldActive      = "active"
)

// Sort parameter values
const (
	SortDueDateAsc    = "due_date_asc"
	SortDueDateDesc   = "due_date_desc"
	SortPriorityAsc   = "priority_asc"
	SortPriorityDesc  = "priority_desc"
	SortCreatedAtAsc  = "created_at_asc"
	SortCreatedAtDesc = "created_at_desc"
)

// Filter parameter values
const (
	FilterUrgent   = "urgent"
	FilterOverdue  = "overdue"
	FilterToday    = "today"
	FilterThisWeek = "this_week"
)
