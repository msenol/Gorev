# Gorev API Reference

> **Version**: This documentation is valid for v0.15.22+
> **Last Updated**: September 21, 2025

Complete API reference for Gorev MCP server, data models, and programmatic interfaces.

## Table of Contents

- [Data Models](#data-models)
- [MCP Protocol](#mcp-protocol)
- [CLI Commands](#cli-commands)
- [Database Schema](#database-schema)
- [Error Handling](#error-handling)
- [Configuration](#configuration)

## Data Models

### Task (Gorev)

Core task model with full feature support.

```go
type Gorev struct {
    ID              int               `json:"id"`
    Baslik          string            `json:"baslik"`          // Title
    Aciklama        string            `json:"aciklama"`        // Description (Markdown)
    Durum           string            `json:"durum"`           // Status
    Oncelik         string            `json:"oncelik"`         // Priority
    ProjeID         *int              `json:"proje_id,omitempty"`
    UstGorevID      *int              `json:"ust_gorev_id,omitempty"` // Parent task
    OlusturmaTarih  time.Time         `json:"olusturma_tarih"`
    GuncellemeTarih time.Time         `json:"guncelleme_tarih"`
    SonTarih        *time.Time        `json:"son_tarih,omitempty"`
    Etiketler       []string          `json:"etiketler,omitempty"`
    AltGorevler     []Gorev           `json:"alt_gorevler,omitempty"`
    TamamlanmaYuzdesi int             `json:"tamamlanma_yuzdesi,omitempty"`
}
```

**Field Descriptions:**

| Field | Type | Description | Values |
|-------|------|-------------|--------|
| `ID` | `int` | Auto-increment primary key | |
| `Baslik` | `string` | Task title (required, max 200 chars) | |
| `Aciklama` | `string` | Detailed description (Markdown supported) | |
| `Durum` | `string` | Task status | `beklemede`, `devam_ediyor`, `tamamlandi` |
| `Oncelik` | `string` | Task priority | `dusuk`, `orta`, `yuksek` |
| `ProjeID` | `*int` | Associated project ID (optional) | |
| `UstGorevID` | `*int` | Parent task ID for subtasks | |
| `SonTarih` | `*time.Time` | Due date (optional) | |
| `Etiketler` | `[]string` | Task tags/labels | |
| `AltGorevler` | `[]Gorev` | Child tasks (subtasks) | |
| `TamamlanmaYuzdesi` | `int` | Completion percentage (calculated) | 0-100 |

### Project (Proje)

Project organization model.

```go
type Proje struct {
    ID              int       `json:"id"`
    Isim            string    `json:"isim"`            // Name
    Tanim           string    `json:"tanim"`           // Description
    OlusturmaTarih  time.Time `json:"olusturma_tarih"`
    GuncellemeTarih time.Time `json:"guncelleme_tarih"`
    GorevSayisi     int       `json:"gorev_sayisi,omitempty"`
    AktifMi         bool      `json:"aktif_mi"`
}
```

### Task Template (GorevTemplate)

Template system for structured task creation.

```go
type GorevTemplate struct {
    ID                int                     `json:"id"`
    Isim              string                  `json:"isim"`              // Template name
    Tanim             string                  `json:"tanim"`             // Description
    VarsayilanBaslik  string                  `json:"varsayilan_baslik"` // Default title
    AciklamaTemplate  string                  `json:"aciklama_template"` // Description template
    Alanlar          []TemplateAlan          `json:"alanlar"`           // Custom fields
    OrnekDegerler    map[string]interface{}  `json:"ornek_degerler"`    // Example values
    Kategori         string                  `json:"kategori"`          // Category
    Aktif            bool                    `json:"aktif"`             // Active status
}
```

**Available Templates:**
- `bug-report` - Bug reports and fixes
- `feature` - New features and enhancements
- `task` - General tasks and activities
- `meeting` - Meeting planning and notes
- `research` - Research and investigation tasks

### Template Field (TemplateAlan)

Custom fields for templates.

```go
type TemplateAlan struct {
    Isim       string   `json:"isim"`                        // Field name
    Tip        string   `json:"tip"`                         // Field type
    Zorunlu    bool     `json:"zorunlu"`                     // Required field
    Varsayilan string   `json:"varsayilan,omitempty"`        // Default value
    Secenekler []string `json:"secenekler,omitempty"`        // Options for select
}
```

**Field Types:**
- `text` - Text input
- `number` - Numeric input
- `select` - Dropdown selection
- `date` - Date picker
- `textarea` - Multi-line text

## MCP Protocol

### Available Tools

Gorev provides 25+ MCP tools for comprehensive task management:

#### Core Task Management
- `gorev_template_olustur` - Create task from template
- `gorev_listele` - List tasks with filtering
- `gorev_detay` - Get task details
- `gorev_guncelle` - Update task properties
- `gorev_sil` - Delete task
- `gorev_durum_degistir` - Change task status

#### Project Management
- `proje_olustur` - Create new project
- `proje_listele` - List projects
- `proje_aktif_yap` - Set active project
- `proje_detay` - Get project details

#### Advanced Features
- `gorev_etiket_ekle` - Add tags to task
- `gorev_etiket_kaldir` - Remove tags from task
- `gorev_son_tarih` - Set due date
- `gorev_bagimliligi_ekle` - Add task dependency
- `gorev_alt_olustur` - Create subtask
- `gorev_ust_degistir` - Change parent task

#### Search and Filtering
- `gorev_ara` - Search tasks
- `etiket_listele` - List all tags
- `gorev_oncelik_filtrele` - Filter by priority
- `gorev_durum_filtrele` - Filter by status

#### Templates and AI
- `template_listele` - List available templates
- `template_detay` - Get template details
- `ai_context_yonetici` - AI context management (v0.9.0+)

### Tool Schema Example

```json
{
  "name": "gorev_template_olustur",
  "description": "Create a new task using a template",
  "inputSchema": {
    "type": "object",
    "properties": {
      "template": {
        "type": "string",
        "description": "Template name (bug-report, feature, task, meeting, research)"
      },
      "title": {
        "type": "string",
        "description": "task title"
      },
      "priority": {
        "type": "string",
        "description": "Priority level (low, medium, high)",
        "enum": ["low", "medium", "high"]
      },
      "due_date": {
        "type": "string",
        "description": "Due date in YYYY-MM-DD format"
      },
      "tags": {
        "type": "string",
        "description": "Comma-separated tags"
      }
    },
    "required": ["template", "title"]
  }
}
```

### Response Format

All MCP tools return responses in this format:

```json
{
  "content": [
    {
      "type": "text",
      "text": "✅ Task created successfully: Fix login bug (#42)"
    }
  ]
}
```

## CLI Commands

### Server Management

```bash
# Start MCP server
gorev serve [--port PORT] [--debug] [--config PATH]

# Server with custom configuration
gorev serve --port 8080 --debug --config ./custom-config.json
```

### Task Operations

```bash
# List tasks
gorev task list [--status STATUS] [--priority PRIORITY] [--project PROJECT]

# Create task
gorev task create --title "Task title" [--description "Description"]

# Show task details
gorev task show <ID>

# Update task
gorev task update <ID> --status completed
```

### Project Operations

```bash
# List projects
gorev project list

# Create project
gorev project create --name "Project name" [--description "Description"]

# Set active project
gorev project active <ID>
```

### Utility Commands

```bash
# Show version
gorev version

# Show help
gorev help

# Health check
gorev serve --test
```

## Database Schema

### Tables

**gorevler** (tasks)
```sql
CREATE TABLE gorevler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    baslik TEXT NOT NULL,
    aciklama TEXT DEFAULT '',
    durum TEXT DEFAULT 'beklemede',
    oncelik TEXT DEFAULT 'orta',
    proje_id INTEGER REFERENCES projeler(id),
    ust_gorev_id INTEGER REFERENCES gorevler(id),
    olusturma_tarih DATETIME DEFAULT CURRENT_TIMESTAMP,
    guncelleme_tarih DATETIME DEFAULT CURRENT_TIMESTAMP,
    son_tarih DATETIME,
    tamamlanma_yuzdesi INTEGER DEFAULT 0
);
```

**projeler** (projects)
```sql
CREATE TABLE projeler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    isim TEXT NOT NULL UNIQUE,
    tanim TEXT DEFAULT '',
    olusturma_tarih DATETIME DEFAULT CURRENT_TIMESTAMP,
    guncelleme_tarih DATETIME DEFAULT CURRENT_TIMESTAMP,
    aktif_mi BOOLEAN DEFAULT 0
);
```

**gorev_etiketler** (task tags)
```sql
CREATE TABLE gorev_etiketler (
    gorev_id INTEGER REFERENCES gorevler(id),
    etiket TEXT NOT NULL,
    PRIMARY KEY (gorev_id, etiket)
);
```

**gorev_templates** (task templates)
```sql
CREATE TABLE gorev_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    isim TEXT NOT NULL UNIQUE,
    tanim TEXT DEFAULT '',
    varsayilan_baslik TEXT DEFAULT '',
    aciklama_template TEXT DEFAULT '',
    kategori TEXT DEFAULT 'general',
    aktif BOOLEAN DEFAULT 1,
    olusturma_tarih DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Indexes

```sql
-- Task indexes for performance
CREATE INDEX idx_gorevler_durum ON gorevler(durum);
CREATE INDEX idx_gorevler_oncelik ON gorevler(oncelik);
CREATE INDEX idx_gorevler_proje_id ON gorevler(proje_id);
CREATE INDEX idx_gorevler_son_tarih ON gorevler(son_tarih);

-- Tag indexes
CREATE INDEX idx_gorev_etiketler_etiket ON gorev_etiketler(etiket);
```

## Error Handling

### Error Types

```go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### Common Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid parameters |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Duplicate resource |
| 422 | Validation Error - Data validation failed |
| 500 | Internal Server Error - Unexpected error |

### Example Error Response

```json
{
  "content": [
    {
      "type": "text", 
      "text": "❌ Error: Task not found (ID: 999)"
    }
  ]
}
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GOREV_DATA_DIR` | Data directory path | `~/.gorev` |
| `GOREV_PORT` | Server port | `8080` |
| `GOREV_LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |
| `GOREV_LANG` | Language (tr, en) | `tr` |
| `GOREV_DB_PATH` | Database file path | `${GOREV_DATA_DIR}/gorev.db` |

### Configuration File

Example `gorev.json`:

```json
{
  "server": {
    "port": 8080,
    "host": "localhost",
    "debug": false
  },
  "database": {
    "path": "./gorev.db",
    "migrations_path": "./migrations"
  },
  "i18n": {
    "default_language": "tr",
    "supported_languages": ["tr", "en"]
  },
  "logging": {
    "level": "info",
    "format": "json",
    "output": "stdout"
  }
}
```

## Internationalization

### Supported Languages

- **Turkish (tr)** - Primary language, full support
- **English (en)** - Full translation support (v0.11.0+)

### Language Detection Priority

1. `--lang` CLI flag
2. `GOREV_LANG` environment variable  
3. `LANG` environment variable
4. Turkish (default)

### Usage

```bash
# Use English interface
GOREV_LANG=en gorev serve

# Use Turkish interface (default)
GOREV_LANG=tr gorev serve
```

## Webhooks and Events

### Event Types

```go
type Event struct {
    Type      string      `json:"type"`
    Timestamp time.Time   `json:"timestamp"`
    Data      interface{} `json:"data"`
}
```

**Available Events:**
- `task.created`
- `task.updated` 
- `task.deleted`
- `task.status_changed`
- `project.created`
- `project.updated`

### Webhook Configuration

```json
{
  "webhooks": {
    "enabled": true,
    "endpoints": [
      {
        "url": "https://example.com/webhook",
        "events": ["task.created", "task.completed"],
        "secret": "webhook-secret"
      }
    ]
  }
}
```

## Performance Considerations

### Database Optimization

- SQLite with WAL mode for better concurrency
- Proper indexing on frequently queried fields
- Connection pooling for multiple requests
- Regular VACUUM operations for maintenance

### Memory Usage

- Lazy loading of subtasks and relations
- Pagination for large result sets
- Configurable cache sizes
- Efficient JSON serialization

### Scaling

- Single binary deployment
- Horizontal scaling via load balancer
- Database replication support
- Stateless server design

---

*This API reference was created with assistance from Claude (Anthropic)*