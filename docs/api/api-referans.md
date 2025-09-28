# API Referans Dokümantasyonu

**Version:** v0.15.24
**Last Updated:** 28 September 2025
**Status:** Production Ready

---

## 📋 İçindekiler

- [Genel Bakış](#genel-bakış)
- [Authentication & Setup](#authentication--setup)
- [MCP Tools Reference](#mcp-tools-reference)
- [Error Handling](#error-handling)
- [Response Formats](#response-formats)
- [Rate Limiting](#rate-limiting)
- [Best Practices](#best-practices)

## 🔍 Genel Bakış

Gorev MCP Server, Model Context Protocol (MCP) standardını kullanarak AI asistanlarına kapsamlı görev yönetimi yetenekleri sağlar. Server 41+ aktif MCP tool ile zengin bir API sunar.

### Desteklenen Protokoller

- **MCP Version:** 2024-11-05
- **Transport:** stdio, HTTP
- **Authentication:** Environment-based
- **Data Format:** JSON

### Server Endpoint

```
npx @mehmetsenol/gorev-mcp-server@latest
```

## 🔐 Authentication & Setup

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GOREV_LANG` | ❌ | `tr` | Interface language (tr/en) |
| `GOREV_DB_PATH` | ❌ | `.gorev/gorev.db` | Database file path |
| `GOREV_ROOT` | ❌ | current dir | Project root directory |
| `GOREV_DEBUG` | ❌ | `false` | Enable debug logging |

### MCP Client Configuration

#### Claude Desktop

```json
{
  "mcpServers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "tr"
      }
    }
  }
}
```

#### VS Code MCP Extension

```json
{
  "servers": {
    "gorev": {
      "command": "npx",
      "args": ["@mehmetsenol/gorev-mcp-server@latest"],
      "env": {
        "GOREV_LANG": "en"
      }
    }
  }
}
```

## 🛠️ MCP Tools Reference

### 1. Görev Yönetimi

#### `gorev_listele`

Lists tasks with filtering and pagination support.

**Parameters:**

```typescript
interface GorevListeleParams {
  durum?: "beklemede" | "devam_ediyor" | "tamamlandi";
  tum_projeler?: boolean;
  sirala?: "son_tarih_asc" | "son_tarih_desc";
  filtre?: "acil" | "gecmis";
  etiket?: string;
  limit?: number;
  offset?: number;
}
```

**Response:**

```json
{
  "content": [{
    "type": "text",
    "text": "## Görev Listesi\n\n- [devam_ediyor] API authentication implementasyonu..."
  }]
}
```

**Example:**

```json
{
  "name": "gorev_listele",
  "arguments": {
    "durum": "devam_ediyor",
    "limit": 10,
    "sirala": "son_tarih_asc"
  }
}
```

#### `gorev_detay`

Get detailed task information in markdown format.

**Parameters:**

```typescript
interface GorevDetayParams {
  id: string;  // UUID format
}
```

**Response:**

```json
{
  "content": [{
    "type": "text",
    "text": "# Task Title\n\n## 📋 Genel Bilgiler\n- **ID:** uuid\n- **Durum:** status..."
  }]
}
```

**Error Cases:**

- `task_not_found`: Task with given ID doesn't exist
- `invalid_uuid`: Malformed task ID

#### `templateden_gorev_olustur`

Create tasks from predefined templates.

**Parameters:**

```typescript
interface TemplatedenGorevOlusturParams {
  template_id: string;
  degerler: Record<string, any>;
}
```

**Template Values Schema:**

```typescript
interface BugReportTemplate {
  baslik: string;        // Required
  aciklama: string;      // Required
  modul: string;         // Required
  ortam: "development" | "staging" | "production"; // Required
  adimlar: string;       // Required
  beklenen: string;      // Required
  mevcut: string;        // Required
  ekler?: string;        // Optional
  cozum?: string;        // Optional
  oncelik: "dusuk" | "orta" | "yuksek"; // Default: orta
  etiketler?: string;    // Default: bug
}
```

**Example:**

```json
{
  "name": "templateden_gorev_olustur",
  "arguments": {
    "template_id": "bug-report",
    "degerler": {
      "baslik": "Login butonu çalışmıyor",
      "aciklama": "Kullanıcı giriş sayfasında login butonu tepki vermiyor",
      "modul": "authentication",
      "ortam": "production",
      "adimlar": "1. Login sayfası\n2. Email/şifre gir\n3. Login buton tıkla",
      "beklenen": "Ana sayfaya yönlendirme",
      "mevcut": "Hiçbir şey olmuyor",
      "oncelik": "yuksek"
    }
  }
}
```

### 2. Proje Yönetimi

#### `proje_olustur`

Create and manage projects.

**Parameters:**

```typescript
interface ProjeOlusturParams {
  isim: string;      // Required, 1-255 chars
  tanim?: string;    // Optional description
}
```

#### `aktif_proje_ayarla`

Set active project for context.

**Parameters:**

```typescript
interface AktifProjeAyarlaParams {
  proje_id: string;  // UUID of existing project
}
```

### 3. Gelişmiş Arama & Filtreleme

#### `gorev_search_advanced`

Advanced search with FTS5 and fuzzy matching.

**Parameters:**

```typescript
interface GorevSearchAdvancedParams {
  query?: string;
  filters?: {
    durum?: string[];
    oncelik?: string[];
    proje_id?: string;
    etiket?: string[];
    tarih_baslangic?: string;  // YYYY-MM-DD
    tarih_bitis?: string;      // YYYY-MM-DD
  };
  use_fuzzy_search?: boolean;  // Default: true
  fuzzy_threshold?: number;    // 0.0-1.0, Default: 0.6
  max_results?: number;        // Default: 50, Max: 100
}
```

**Response:**

```json
{
  "content": [{
    "type": "text",
    "text": "## Arama Sonuçları\n\n**Sorgu:** authentication\n**Bulunan:** 5 görev\n\n..."
  }]
}
```

### 4. Veri İşlemleri

#### `gorev_export`

Export tasks and projects to JSON/CSV.

**Parameters:**

```typescript
interface GorevExportParams {
  output_path: string;
  format?: "json" | "csv";          // Default: json
  include_completed?: boolean;       // Default: true
  include_dependencies?: boolean;    // Default: true
  include_templates?: boolean;       // Default: false
  proje_id?: string;                // Export specific project
}
```

#### `gorev_import`

Import tasks with conflict resolution.

**Parameters:**

```typescript
interface GorevImportParams {
  file_path: string;
  import_mode?: "merge" | "replace";           // Default: merge
  conflict_resolution?: "skip" | "overwrite"; // Default: skip
  dry_run?: boolean;                           // Default: false
}
```

### 5. AI Context Management

#### `gorev_set_active`

Set active task for AI context.

**Parameters:**

```typescript
interface GorevSetActiveParams {
  task_id: string;
}
```

#### `gorev_nlp_query`

Natural language task search.

**Parameters:**

```typescript
interface GorevNlpQueryParams {
  query: string;  // Natural language query
}
```

**Example:**

```json
{
  "name": "gorev_nlp_query",
  "arguments": {
    "query": "bu hafta tamamlanması gereken acil görevler"
  }
}
```

## ❌ Error Handling

### Standard Error Format

```json
{
  "error": {
    "code": "error_code",
    "message": "Human readable error message",
    "details": {
      "field": "Additional context"
    }
  }
}
```

### Common Error Codes

| Code | Description | Resolution |
|------|-------------|------------|
| `invalid_parameters` | Missing or invalid parameters | Check parameter types and requirements |
| `task_not_found` | Task ID doesn't exist | Verify task ID with `gorev_listele` |
| `project_not_found` | Project ID doesn't exist | Check project with `proje_listele` |
| `template_not_found` | Template ID doesn't exist | List templates with `template_listele` |
| `database_error` | SQLite database error | Check database permissions and disk space |
| `file_not_found` | Import/export file not found | Verify file path exists |
| `permission_denied` | File system permission error | Check file/directory permissions |
| `validation_error` | Data validation failed | Check required fields and formats |

### Error Handling Best Practices

1. **Always handle errors gracefully**

```typescript
try {
  const result = await mcp.call("gorev_listele", params);
} catch (error) {
  if (error.code === "task_not_found") {
    // Handle missing task
  } else {
    // Handle other errors
  }
}
```

2. **Check parameter validity before calling**

```typescript
if (!isValidUUID(taskId)) {
  throw new Error("Invalid task ID format");
}
```

3. **Use dry_run for destructive operations**

```typescript
// Test import first
await mcp.call("gorev_import", { ...params, dry_run: true });
// Then do actual import
await mcp.call("gorev_import", params);
```

## 📊 Response Formats

### Text Response

Most tools return markdown-formatted text content:

```json
{
  "content": [{
    "type": "text",
    "text": "## Markdown Content\n\nFormatted response..."
  }]
}
```

### Structured Data Response

Some tools return structured data for programmatic use:

```json
{
  "content": [{
    "type": "text",
    "text": "✅ Operation successful",
    "metadata": {
      "task_id": "uuid",
      "created_at": "2025-09-28T10:30:00Z"
    }
  }]
}
```

## ⚡ Rate Limiting

### Default Limits

- **Requests per minute:** 60
- **Concurrent requests:** 5
- **Database operations:** No specific limit

### Best Practices

1. **Batch operations when possible**
2. **Use pagination for large datasets**
3. **Cache frequently accessed data**
4. **Implement exponential backoff for retries**

## 🎯 Best Practices

### 1. Task Management

- Always use templates for task creation
- Include meaningful descriptions
- Set appropriate priorities and due dates
- Use consistent tagging conventions

### 2. Project Organization

- Create projects for logical groupings
- Use descriptive project names
- Set active project for context
- Regularly archive completed projects

### 3. Search & Filtering

- Use fuzzy search for better matching
- Apply multiple filters for precision
- Save filter profiles for repeated use
- Use natural language queries when appropriate

### 4. Data Management

- Regular exports for backup
- Test imports with dry_run first
- Maintain consistent data formats
- Monitor database size and performance

### 5. Error Handling

- Always check for errors in responses
- Implement proper retry logic
- Log errors for debugging
- Provide user-friendly error messages

## 📚 Related Documentation

- **[MCP Tools Turkish Reference](../tr/mcp-araclari.md)** - Detailed Turkish documentation
- **[Usage Guide](../guides/user/usage.md)** - Basic usage patterns
- **[Architecture Guide](../architecture/technical-specification-v2.md)** - System architecture
- **[Development Guide](../development/contributing.md)** - Contributing guidelines

---

> 💡 **Note**: This API reference is for v0.15.24+. For older versions, check the relevant release documentation.
