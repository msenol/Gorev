# MCP Tools Reference

Complete reference for all 41 MCP tools implemented by the Gorev server, moved from CLAUDE.md for better organization.

**Note**: `gorev_olustur` was removed in v0.11.1 - use `templateden_gorev_olustur` with template aliases instead.

## 🆕 NEW: IDE Management Tools (v0.12.0)

The Gorev MCP server now includes comprehensive IDE extension management capabilities. These tools enable automatic detection, installation, updating, and management of the Gorev VS Code extension across multiple IDEs.

## Task Management

### 1. gorev_listele

- **Purpose**: List tasks
- **Parameters**: durum?, tum_projeler?, sirala?, filtre?, etiket?, limit?, offset?
- **Options**:
  - tum_projeler: if false/omitted, shows only active project tasks
  - sirala: son_tarih_asc, son_tarih_desc
  - filtre: acil (due in 7 days), gecmis (overdue)
  - etiket: filter by tag name
  - limit: maximum number of tasks to return (default: 50)
  - offset: number of tasks to skip for pagination (default: 0)

### 2. gorev_detay

- **Purpose**: Show detailed task info in markdown
- **Parameters**: id
- **Features**: Shows due dates, tags, and dependencies with status indicators

### 3. gorev_guncelle

- **Purpose**: Update task status
- **Parameters**: id, durum
- **Validation**: Dependencies must be completed before allowing "devam_ediyor" status

### 4. gorev_duzenle

- **Purpose**: Edit task properties
- **Parameters**: id, baslik?, aciklama?, oncelik?, proje_id?, son_tarih?

### 5. gorev_sil

- **Purpose**: Delete task
- **Parameters**: id, onay
- **Safety**: Prevents deletion if task has subtasks

### 6. gorev_bagimlilik_ekle

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
- **Notes**:
  - template_id: Can be either UUID or alias (e.g., `bug`, `feature`, `research`)
  - degerler: Object with field values for the template
  - **Template Aliases Available**: `bug`, `bug2`, `feature`, `research`, `spike`, `security`, `performance`, `refactor`, `debt`
- **Discovery**: Use `gorev template aliases` CLI command to see all shortcuts
- **Examples**:
  - Using alias: `template_id='bug'`
  - Using UUID: `template_id='550e8400-e29b-41d4-a716-446655440000'`

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

## IDE Management

### 27. gorev_ide_detect

- **Purpose**: Detect installed IDEs on the system (VS Code, Cursor, Windsurf)
- **Parameters**: (none)
- **Features**:
  - Automatically discovers IDE installations across all platforms
  - Returns detailed information including executable paths and versions
  - Detects configuration and extension directories
  - Cross-platform support (Windows, macOS, Linux)

### 28. gorev_ide_install

- **Purpose**: Install Gorev extension to specified IDE(s)
- **Parameters**: ide_type
- **Options**:
  - ide_type: "vscode", "cursor", "windsurf", or "all"
  - "all" installs to all detected IDEs
- **Features**:
  - Downloads latest extension version from GitHub releases
  - Automatic VSIX file management with cleanup
  - Version checking to avoid duplicate installations
  - Progress reporting and error handling

### 29. gorev_ide_uninstall

- **Purpose**: Remove Gorev extension from specified IDE
- **Parameters**: ide_type, extension_id?
- **Options**:
  - ide_type: "vscode", "cursor", or "windsurf"
  - extension_id: Extension identifier (default: mehmetsenol.gorev-vscode)
- **Features**: Safe extension removal with confirmation

### 30. gorev_ide_status

- **Purpose**: Check installation status of Gorev extension in all detected IDEs
- **Parameters**: (none)
- **Features**:
  - Shows installation status for each detected IDE
  - Displays installed versions vs latest available version
  - Identifies outdated installations requiring updates
  - Version comparison with GitHub releases

### 31. gorev_ide_update

- **Purpose**: Update Gorev extension to latest version
- **Parameters**: ide_type
- **Options**:
  - ide_type: "vscode", "cursor", "windsurf", or "all"
  - "all" updates all detected IDEs with outdated extensions
- **Features**:
  - Automatic latest version detection from GitHub
  - Selective or bulk update operations
  - Version verification and rollback support

## Data Export/Import

### 32. gorev_export

- **Purpose**: Export tasks, projects and related data to file in JSON or CSV format
- **Parameters**: output_path, format?, include_completed?, include_dependencies?, include_templates?, include_ai_context?, project_filter?, date_range?
- **Options**:
  - output_path: Path where the exported file will be saved (required)
  - format: Export format - "json" or "csv" (default: "json")
  - include_completed: Include completed tasks (default: true)
  - include_dependencies: Include task dependencies (default: true)
  - include_templates: Include templates (default: false)
  - include_ai_context: Include AI context data (default: false)
  - project_filter: Array of project IDs to export (optional)
  - date_range: Object with "from" and "to" date fields in ISO 8601 format (optional)
- **Use Cases**: Backup, data sharing, migration, reporting
- **VS Code Integration**: Available through Extension commands - Export Data, Export Current View, Quick Export
- **Output**: Creates export file with comprehensive task management data
- **Example**:

  ```json
  {
    "output_path": "/path/to/export.json",
    "format": "json",
    "include_completed": true,
    "include_dependencies": true,
    "project_filter": ["project-1", "project-2"],
    "date_range": {
      "from": "2025-01-01T00:00:00Z",
      "to": "2025-12-31T23:59:59Z"
    }
  }
  ```

### 33. gorev_import

- **Purpose**: Import previously exported data back into the system
- **Parameters**: file_path, import_mode?, conflict_resolution?, preserve_ids?, dry_run?, project_mapping?
- **Options**:
  - file_path: Path to the file to import (required)
  - import_mode: Import strategy - "merge" (add to existing) or "replace" (replace all) (default: "merge")
  - conflict_resolution: How to handle conflicts - "skip", "overwrite", or "prompt" (default: "skip")
  - preserve_ids: Preserve original IDs from export (default: false)
  - dry_run: Only analyze, don't make changes (default: false)
  - project_mapping: Object mapping old project IDs to new ones (optional)
- **Features**:
  - Conflict detection and resolution
  - Data validation and integrity checks
  - Dry run mode for safe preview
  - Project remapping capabilities
  - Detailed import statistics and error reporting
- **Use Cases**: Data restoration, migration between instances, bulk data updates
- **VS Code Integration**: Available through Extension Import Data command with multi-step wizard UI
- **Example**:

  ```json
  {
    "file_path": "/path/to/export.json",
    "import_mode": "merge",
    "conflict_resolution": "skip",
    "preserve_ids": false,
    "dry_run": true,
    "project_mapping": {
      "old-project-id": "new-project-id"
    }
  }
  ```

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
