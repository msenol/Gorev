# MCP Tools Reference - v0.17.0

Complete reference for **26 optimized MCP tools** (unified from 41 in v0.16.0).

**Last Updated**: December 24, 2025

## 📋 Tool Categories

### CORE TOOLS (11)

Essential operations for task management, templates, projects, and dependencies.

#### Task Management (5)

1. `gorev_listele` - List and filter tasks
2. `gorev_detay` - Show detailed task information
3. `gorev_guncelle` - Update task status/priority
4. `gorev_duzenle` - Edit task content and properties
5. `gorev_sil` - Delete task with safety checks

#### Templates (2)

6. `template_listele` - List available templates
7. `templateden_gorev_olustur` - Create task from template

#### Projects (3)

8. `proje_listele` - List all projects
9. `proje_olustur` - Create new project
10. `proje_gorevleri` - List tasks in project

#### Dependencies (1)

11. `gorev_bagimlilik_ekle` - Add task dependency relationship

### UNIFIED TOOLS (7)

Consolidated handlers that replace multiple specialized tools through action-based routing.

12. `aktif_proje` - Active project management (set|get|clear)
13. `gorev_bulk` - Bulk operations (transition|tag|update)
14. `gorev_hierarchy` - Task hierarchy operations (create_subtask|change_parent|show)
15. `gorev_filter_profile` - Filter profile management (save|load|list|delete)
16. `gorev_ide` - IDE extension management (detect|install|uninstall|status|update)
17. `gorev_context` - AI context management (set_active|get_active|recent|summary)
18. `gorev_search` - Task search (nlp|advanced|history)

### FILE WATCHER TOOLS (4)

Individual tools for file watching operations.

19. `gorev_file_watch_add` - Add file path to watch
20. `gorev_file_watch_remove` - Remove file path from watch
21. `gorev_file_watch_list` - List active file watches
22. `gorev_file_watch_stats` - Show file watch statistics

### SPECIAL TOOLS (4)

Advanced features for summaries, data management, and AI-powered operations.

23. `ozet_goster` - Show workspace summary with HTML dashboard
24. `gorev_suggestions` - Get AI-powered task suggestions
25. `gorev_export` - Export tasks to various formats
26. `gorev_import` - Import tasks from external sources

> **Template Aliases**: `bug`, `feature`, `research`, `refactor`, `test`, `doc`

---

## 🔧 Detailed Tool Specifications

### CORE TOOLS

#### 1. gorev_listele

**Purpose**: List and filter tasks with hierarchical display

**Parameters**:

- `durum` (optional): Filter by status (beklemede|devam_ediyor|tamamlandi|iptal)
- `tum_projeler` (optional): boolean - if true, shows all projects; if false/omitted, shows only active project
- `sirala` (optional): Sort order (son_tarih_asc|son_tarih_desc)
- `filtre` (optional): Quick filters (acil - due in 7 days, gecmis - overdue)
- `etiket` (optional): Filter by tag name
- `limit` (optional): Maximum tasks to return (default: 50)
- `offset` (optional): Number of tasks to skip for pagination (default: 0)

**Output**: Hierarchical tree structure showing tasks with status indicators, priorities, and dependencies

**Example**:

```json
{
  "durum": "devam_ediyor",
  "tum_projeler": false,
  "limit": 20
}
```

---

#### 2. gorev_detay

**Purpose**: Show detailed information about a specific task

**Parameters**:

- `id` (required): Task ID (8-character short ID or full UUID)

**Output**: Comprehensive task details including:

- Title, description, status, priority
- Due date with countdown/overdue indicators
- Tags and project association
- Subtasks and hierarchy information
- Dependencies (blocked by, blocking)
- Creation and update timestamps

**Example**:

```json
{
  "id": "abc12345"
}
```

---

#### 3. gorev_guncelle ⭐ **v0.16.3 UPDATED**

**Purpose**: Update task status and/or priority

**What's New in v0.16.3**:

- Now supports both `durum` and `oncelik` parameters (previously only status)
- Can update status, priority, or both simultaneously
- At least one parameter required

**Parameters**:

- `id` (required): Task ID
- `durum` (optional): New status (beklemede|devam_ediyor|tamamlandi|iptal)
- `oncelik` (optional): New priority (dusuk|orta|yuksek)

**Note**: At least one of `durum` or `oncelik` must be provided

**Validation**:

- Checks task dependencies before allowing status transitions
- Prevents moving to "devam_ediyor" if dependencies are incomplete

**Examples**:

```json
// Update status only
{
  "id": "abc12345",
  "durum": "devam_ediyor"
}

// Update priority only (NEW in v0.16.3)
{
  "id": "abc12345",
  "oncelik": "yuksek"
}

// Update both simultaneously (NEW in v0.16.3)
{
  "id": "abc12345",
  "durum": "devam_ediyor",
  "oncelik": "yuksek"
}
```

---

#### 4. gorev_duzenle

**Purpose**: Edit task content and properties

**Parameters**:

- `id` (required): Task ID
- `baslik` (optional): New title
- `aciklama` (optional): New description (markdown supported)
- `oncelik` (optional): New priority (dusuk|orta|yuksek)
- `proje_id` (optional): Move to different project
- `son_tarih` (optional): New due date (YYYY-MM-DD format)

**Example**:

```json
{
  "id": "abc12345",
  "baslik": "Updated task title",
  "oncelik": "yuksek",
  "son_tarih": "2025-12-31"
}
```

---

#### 5. gorev_sil

**Purpose**: Delete a task with safety checks

**Parameters**:

- `id` (required): Task ID
- `onay` (required): Confirmation (must be "evet" or "yes")

**Safety Features**:

- Prevents deletion if task has subtasks (must delete subtasks first)
- Removes all dependencies automatically
- Confirmation required to prevent accidental deletion

**Example**:

```json
{
  "id": "abc12345",
  "onay": "evet"
}
```

---

#### 6. template_listele

**Purpose**: List all available task templates

**Parameters**: None required

**Output**: List of templates with:

- Template name and alias (e.g., "bug", "feature", "research")
- Description
- Default priority
- Variable placeholders

**Example**:

```json
{}
```

---

#### 7. templateden_gorev_olustur

**Purpose**: Create a task from a template

**Parameters**:

- `template_id` (required): Template ID or alias (bug|feature|research|refactor|test|doc)
- `degerler` (required): Object with variable values
- `proje_id` (optional): Project ID (uses active project if omitted)

**Template Aliases**: Use human-readable aliases instead of UUIDs:

- `bug` - Bug fix template
- `feature` - Feature implementation
- `research` - Research task
- `refactor` - Code refactoring
- `test` - Testing task
- `doc` - Documentation

**Example**:

```json
{
  "template_id": "bug",
  "degerler": {
    "baslik": "Fix login button not responding",
    "aciklama": "Users report clicking login button has no effect",
    "oncelik": "yuksek"
  }
}
```

---

#### 8. proje_listele

**Purpose**: List all projects

**Parameters**: None required

**Output**: List of projects with:

- Project name and description
- Task counts (total, completed, in progress)
- Active project indicator

---

#### 9. proje_olustur

**Purpose**: Create a new project

**Parameters**:

- `isim` (required): Project name
- `tanim` (optional): Project description

**Example**:

```json
{
  "isim": "Mobile App Redesign",
  "tanim": "Complete redesign of mobile application UI/UX"
}
```

---

#### 10. proje_gorevleri

**Purpose**: List all tasks in a specific project

**Parameters**:

- `proje_id` (required): Project ID

**Output**: Hierarchical task list for the specified project

---

#### 11. gorev_bagimlilik_ekle

**Purpose**: Create a dependency relationship between tasks

**Parameters**:

- `gorev_id` (required): Task that depends on another
- `bagli_gorev_id` (required): Task that must be completed first
- `baglanti_tipi` (optional): Dependency type (blocker|depends_on) - default: depends_on

**Validation**:

- Prevents circular dependencies
- Checks both tasks exist

**Example**:

```json
{
  "gorev_id": "abc12345",
  "bagli_gorev_id": "def67890",
  "baglanti_tipi": "blocker"
}
```

---

### UNIFIED TOOLS

#### 12. aktif_proje

**Purpose**: Manage active project (unified handler for 3 operations)

**Parameters**:

- `action` (required): "set" | "get" | "clear"
- `proje_id` (required for set): Project ID to set as active

**Actions**:

- **set**: Set active project
- **get**: Show current active project
- **clear**: Clear active project

**Examples**:

```json
// Set active project
{
  "action": "set",
  "proje_id": "proj123"
}

// Get active project
{
  "action": "get"
}

// Clear active project
{
  "action": "clear"
}
```

---

#### 13. gorev_hierarchy

**Purpose**: Manage task hierarchy (unified handler for 3 operations)

**Parameters**:

- `action` (required): "create_subtask" | "change_parent" | "show"
- `parent_id` (required for create_subtask/show): Parent task ID
- `baslik` (required for create_subtask): Subtask title
- `aciklama` (optional): Subtask description
- `gorev_id` (required for change_parent): Task ID to move
- `yeni_parent_id` (optional for change_parent): New parent (empty = root level)

**Actions**:

- **create_subtask**: Create a new subtask under parent
- **change_parent**: Move task to different parent
- **show**: Show task hierarchy tree

**Examples**:

```json
// Create subtask
{
  "action": "create_subtask",
  "parent_id": "abc12345",
  "baslik": "Implement login validation",
  "aciklama": "Add client-side validation for login form"
}

// Change parent
{
  "action": "change_parent",
  "gorev_id": "def67890",
  "yeni_parent_id": "abc12345"
}

// Show hierarchy
{
  "action": "show",
  "parent_id": "abc12345"
}
```

---

#### 14. gorev_bulk ⭐ **v0.16.3 FIXED**

**Purpose**: Perform bulk operations on multiple tasks

**What's New in v0.16.3**:

- Complete rewrite with proper parameter transformation
- All 3 operations now fully functional (update, transition, tag)
- Flexible parameter naming for backward compatibility
- 100% success rate confirmed by testing

**Parameters**:

- `operation` (required): "transition" | "tag" | "update"
- `ids` (required): Array of task IDs
- `data` (required): Operation-specific data object

**Operation: update** ⭐ **FIXED in v0.16.3**

Transforms `{ids, data}` → `{updates: [{id, ...fields}]}` automatically

```json
{
  "operation": "update",
  "ids": ["abc123", "def456"],
  "data": {
    "oncelik": "yuksek",
    "durum": "devam_ediyor"
  }
}
```

Internal transformation: Creates array of update objects, each with ID + data fields

**Operation: transition**

Accepts both `durum` and `yeni_durum` parameter names (flexible)

```json
{
  "operation": "transition",
  "ids": ["abc123", "def456"],
  "data": {
    "durum": "devam_ediyor",  // or "yeni_durum"
    "force": false,
    "check_dependencies": true,
    "dry_run": false
  }
}
```

**Operation: tag**

Accepts both `operation` and `tag_operation` parameter names (flexible)

```json
{
  "operation": "tag",
  "ids": ["abc123", "def456"],
  "data": {
    "tags": ["frontend", "urgent"],
    "operation": "add",  // or "tag_operation", values: add|remove|replace
    "dry_run": false
  }
}
```

**Performance**: 11-33ms processing time for bulk operations

---

#### 15. gorev_filter_profile

**Purpose**: Manage filter profiles (save/load search filters)

**Parameters**:

- `action` (required): "save" | "load" | "list" | "delete"
- `name` (required for save/load/delete): Profile name
- `filters` (required for save): Filter object to save

**Example**:

```json
// Save filter profile
{
  "action": "save",
  "name": "my-urgent-tasks",
  "filters": {
    "durum": "devam_ediyor",
    "oncelik": "yuksek",
    "filtre": "acil"
  }
}

// Load filter profile
{
  "action": "load",
  "name": "my-urgent-tasks"
}
```

---

#### 17. gorev_file_watch_add

**Purpose**: Add a file path to watch for automatic task creation

**Parameters**:

- `file_path` (required): Path to file to watch
- `gorev_id` (optional): Task ID to link to this file

**Example**:

```json
{
  "file_path": "/path/to/project/tasks.md",
  "gorev_id": "abc12345"
}
```

---

#### 18. gorev_file_watch_remove

**Purpose**: Remove a file path from watching

**Parameters**:

- `file_path` (required): Path to file to remove from watch

**Example**:

```json
{
  "file_path": "/path/to/project/tasks.md"
}
```

---

#### 19. gorev_file_watch_list

**Purpose**: List all active file watches

**Parameters**: None required

**Example**:

```json
{}
```

---

#### 20. gorev_file_watch_stats

**Purpose**: Show file watch statistics and activity

**Parameters**: None required

**Example**:

```json
{}
```

---

#### 17. gorev_ide

**Purpose**: Manage IDE extensions (VS Code extension management)

**Parameters**:

- `action` (required): "detect" | "install" | "uninstall" | "status" | "update"
- `ide` (optional): IDE name (vscode|cursor|windsurf)

**Example**:

```json
{
  "action": "status",
  "ide": "vscode"
}
```

---

#### 18. gorev_context

**Purpose**: AI context management for active task tracking

**Parameters**:

- `action` (required): "set_active" | "get_active" | "recent" | "summary"
- `gorev_id` (required for set_active): Task ID to set as active

**Example**:

```json
{
  "action": "set_active",
  "gorev_id": "abc12345"
}
```

---

#### 19. gorev_search

**Purpose**: Search tasks using NLP, advanced filters, or history

**Parameters**:

- `mode` (required): "nlp" | "advanced" | "history"
- `query` or `arama_metni` (required): Search query

**Mode: nlp** (Natural Language Processing)

Understands natural language queries:

```json
{
  "mode": "nlp",
  "query": "kullanıcı kayıt formu"
}
```

**Mode: advanced**

Supports query string parsing with key:value patterns:

```json
{
  "mode": "advanced",
  "query": "durum:devam_ediyor oncelik:yuksek tags:frontend"
}
```

Automatically parsed to:

```json
{
  "filters": {
    "durum": "devam_ediyor",
    "oncelik": "yuksek",
    "tags": "frontend"
  }
}
```

Additional parameters for advanced mode:

- `filters` (optional): Filter object (if not using query parsing)
- `use_fuzzy_search` (optional): boolean - enable fuzzy matching (default: true)
- `fuzzy_threshold` (optional): number - 0.0 to 1.0 (default: 0.6)
- `max_results` (optional): number - max results (default: 50)
- `sort_by` (optional): "relevance" | "due_date" | "priority"
- `include_completed` (optional): boolean (default: false)

**Mode: history**

Show recent search queries:

```json
{
  "mode": "history"
}
```

**Performance**: 6-67ms with FTS5 full-text search and relevance scoring

---

### SPECIAL TOOLS

#### 20. ozet_goster

**Purpose**: Show workspace summary with HTML dashboard

**Parameters**: None required

**Output** (Enhanced in v0.17.0):

- Progress bars for task status
- Due date summary (overdue, due today, due this week)
- High priority tasks widget
- Recent tasks
- Active project display

---

#### 21. gorev_suggestions

**Purpose**: Get AI-powered task suggestions

**Parameters**:

- `context` (optional): Context string for suggestions

**Output**: AI-generated task suggestions based on:

- Current project context
- Incomplete tasks
- Dependencies
- Project patterns

---

#### 22. gorev_export

**Purpose**: Export tasks to various formats

**Parameters**:

- `format` (required): "json" | "markdown" | "csv"
- `proje_id` (optional): Export specific project only

**Example**:

```json
{
  "format": "markdown"
}
```

---

#### 23. gorev_import

**Purpose**: Import tasks from external sources

**Parameters**:

- `data` (required): Import data object
- `format` (optional): Source format hint

**Example**:

```json
{
  "data": {
    "tasks": [
      {"baslik": "Task 1", "oncelik": "yuksek"},
      {"baslik": "Task 2", "oncelik": "orta"}
    ]
  }
}
```

---

## 📊 Version History

### v0.17.0 (December 24, 2025) - Smart Shutdown & Client Tracking

- **Smart Shutdown**: VS Code extension now checks active clients before stopping daemon
- **Client Tracking System**: New `/api/v1/daemon/clients/*` endpoints for multi-client support
- **MCP Proxy Registration**: MCP clients register with daemon, send heartbeat every 60s
- **gorev_ide**: Renamed from `ide_manage` to `gorev_ide` for consistency
- **File Watcher Tools**: Split from unified `gorev_file_watch` into 4 individual tools
- **Enhanced Summary Dashboard**: New HTML dashboard with progress bars and due date summary

### v0.16.3 (October 6, 2025) - Critical Fixes

- ⭐ **gorev_bulk**: All 3 operations fixed with proper parameter transformation
- ⭐ **gorev_guncelle**: Extended to support both status and priority updates
- ⭐ **gorev_search**: Advanced mode enhanced with query string parsing

### v0.16.0 (October 3, 2025) - Tool Unification

- Reduced from 45 tools to 26 optimized unified tools
- 8 unified handlers with action-based routing
- Improved maintainability and consistency

---

## 🔗 Related Documentation

- [Daemon Architecture](../architecture/daemon-architecture.md)
- [MCP Configuration Examples](../guides/mcp-config-examples.md)
- [VS Code Extension Guide](../guides/user/vscode-extension.md)
- [CHANGELOG](../../gorev-mcpserver/CHANGELOG.md)

---

**Note**: This documentation reflects the current v0.17.0 implementation. All tools are production-ready with 100% test success rate.
