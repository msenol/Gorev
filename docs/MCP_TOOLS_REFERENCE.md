# MCP Tools Reference

Complete reference for all 25 MCP tools implemented by the Gorev server, moved from CLAUDE.md for better organization.

## Task Management

### 1. gorev_olustur
**DEPRECATED as of v0.10.0** - Use `templateden_gorev_olustur` instead
- **Purpose**: Create new task 
- **Parameters**: baslik, aciklama, oncelik, proje_id?, son_tarih?, etiketler?
- **Notes**: 
  - proje_id is optional; if not provided, uses active project
  - son_tarih: optional due date in YYYY-MM-DD format
  - etiketler: optional comma-separated tags

### 2. gorev_listele
- **Purpose**: List tasks
- **Parameters**: durum?, tum_projeler?, sirala?, filtre?, etiket?, limit?, offset?
- **Options**:
  - tum_projeler: if false/omitted, shows only active project tasks
  - sirala: son_tarih_asc, son_tarih_desc
  - filtre: acil (due in 7 days), gecmis (overdue)
  - etiket: filter by tag name
  - limit: maximum number of tasks to return (default: 50)
  - offset: number of tasks to skip for pagination (default: 0)

### 3. gorev_detay
- **Purpose**: Show detailed task info in markdown
- **Parameters**: id
- **Features**: Shows due dates, tags, and dependencies with status indicators

### 4. gorev_guncelle
- **Purpose**: Update task status
- **Parameters**: id, durum
- **Validation**: Dependencies must be completed before allowing "devam_ediyor" status

### 5. gorev_duzenle
- **Purpose**: Edit task properties
- **Parameters**: id, baslik?, aciklama?, oncelik?, proje_id?, son_tarih?

### 6. gorev_sil
- **Purpose**: Delete task
- **Parameters**: id, onay
- **Safety**: Prevents deletion if task has subtasks

### 7. gorev_bagimlilik_ekle
- **Purpose**: Create task dependency
- **Parameters**: kaynak_id, hedef_id, baglanti_tipi

## Subtask Management

### 8. gorev_altgorev_olustur
- **Purpose**: Create subtask under a parent
- **Parameters**: parent_id, baslik, aciklama?, oncelik?, son_tarih?, etiketler?
- **Rules**: 
  - Subtask inherits parent's project
  - parent_id: ID of the parent task

### 9. gorev_ust_degistir
- **Purpose**: Change task's parent
- **Parameters**: gorev_id, yeni_parent_id?
- **Options**:
  - yeni_parent_id: empty string moves task to root level
  - Validates circular dependencies

### 10. gorev_hiyerarsi_goster
- **Purpose**: Show task hierarchy
- **Parameters**: gorev_id
- **Features**: Shows parent hierarchy, subtask statistics, and progress

## Task Templates

### 11. template_listele
- **Purpose**: List available templates
- **Parameters**: kategori?
- **Features**: Shows predefined templates for consistent task creation

### 12. templateden_gorev_olustur
- **Purpose**: Create task from template (**PREFERRED METHOD**)
- **Parameters**: template_id, degerler
- **Notes**: degerler is an object with field values for the template

## Project Management

### 13. proje_olustur
- **Purpose**: Create project
- **Parameters**: isim, tanim

### 14. proje_listele
- **Purpose**: List all projects with task counts
- **Parameters**: (no params)

### 15. proje_gorevleri
- **Purpose**: List project tasks grouped by status
- **Parameters**: proje_id, limit?, offset?
- **Options**:
  - limit: maximum number of tasks to return (default: 50)
  - offset: number of tasks to skip for pagination (default: 0)

### 16. proje_aktif_yap
- **Purpose**: Set active project
- **Parameters**: proje_id

### 17. aktif_proje_goster
- **Purpose**: Show current active project
- **Parameters**: (no params)

### 18. aktif_proje_kaldir
- **Purpose**: Remove active project setting
- **Parameters**: (no params)

## Reporting

### 19. ozet_goster
- **Purpose**: Show summary statistics
- **Parameters**: (no params)

## AI Context Management

### 20. gorev_set_active
- **Purpose**: Set active task for AI session
- **Parameters**: task_id
- **Features**:
  - Automatically transitions task to "devam_ediyor" status
  - Maintains context across AI interactions

### 21. gorev_get_active
- **Purpose**: Get current active task
- **Parameters**: (no params)
- **Features**: Returns detailed information about the active task

### 22. gorev_recent
- **Purpose**: Get recent tasks interacted with
- **Parameters**: limit?
- **Default**: limit = 5 recent tasks

### 23. gorev_context_summary
- **Purpose**: Get AI-optimized session summary
- **Parameters**: (no params)
- **Features**: Shows active task, recent tasks, priorities, and blockers

### 24. gorev_batch_update
- **Purpose**: Batch update multiple tasks
- **Parameters**: updates
- **Format**: updates = array of {id: string, updates: {durum?: string, ...}}
- **Use Case**: Efficient bulk operations for AI workflows

### 25. gorev_nlp_query
- **Purpose**: Natural language task search
- **Parameters**: query
- **Supported Queries**:
  - "bugün üzerinde çalıştığım görevler"
  - "yüksek öncelikli"
  - "database ile ilgili"
  - Tag search: "etiket:bug" or "tag:frontend"
  - Status queries: "tamamlanmamış", "devam eden", "tamamlanan"
  - Time-based: "acil", "gecikmiş", "son oluşturulan"

## Implementation Notes

- All tools follow the pattern in `internal/mcp/handlers.go`
- Tools are registered in `RegisterTools()` with proper schema
- Task descriptions support full markdown formatting
- **Auto-State Management**: `gorev_detay` automatically transitions tasks from "beklemede" to "devam_ediyor" status
- **Template System**: As of v0.10.0, direct task creation via `gorev_olustur` is deprecated

## Adding New MCP Tools

1. Add handler method to `internal/mcp/handlers.go`
2. Register tool in `RegisterTools()` with proper schema
3. Add integration tests in `test/integration_test.go`
4. Update this documentation with tool details