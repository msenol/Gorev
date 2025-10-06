# Template System Guide

**Version**: v0.16.0
**Last Updated**: October 5, 2025
**Feature Status**: Production Ready ✅

---

## Overview

Gorev's **template system** provides structured, standardized task creation with predefined fields and validation rules. Templates ensure consistency, reduce errors, and speed up task creation for common workflows like bug reports, feature requests, and documentation tasks.

### Key Features

- ✅ **6 Default Templates**: Bug reports, features, research, refactoring, tests, documentation
- ✅ **Template Aliases**: Human-readable shortcuts (`bug`, `feature`, `research`, etc.)
- ✅ **Field Validation**: Required vs. optional fields with type checking
- ✅ **Dynamic Forms**: Auto-generated UI based on template schema
- ✅ **Title Patterns**: Automatic title formatting with placeholders
- ✅ **Custom Templates**: Create your own templates (future enhancement)
- ✅ **Multi-Language**: Turkish and English field names

---

## Concepts

### Template

A **template** is a blueprint for creating tasks with:

- **Unique ID**: UUID identifier
- **Alias**: Short, memorable name (`bug`, `feature`)
- **Category**: Grouping (Technical, Feature, Process, Research)
- **Title Pattern**: Format string for auto-generated titles
- **Field Schema**: List of fields with types, validation, and defaults

### Template Alias

**Template aliases** are human-readable shortcuts introduced in v0.16.0:

| Alias | Full Template Name | UUID (varies per workspace) |
|-------|-------------------|----------------------------|
| `bug` | Bug Raporu | `39f28dbd-...` |
| `feature` | Özellik Geliştirme | `7a3c9f2e-...` |
| `research` | Araştırma | `5d8b4a1c-...` |
| `refactor` | Refactoring | `2f6e8c9a-...` |
| `test` | Test Yazma | `9b1d5f3c-...` |
| `doc` | Dokümantasyon | `4c7a2e8b-...` |

**Benefits**:

- No need to remember or look up UUIDs
- Consistent across all workspaces
- Easier to use in AI assistant prompts
- Portable between environments

### Field Types

Templates support various field types:

| Type | Description | Example |
|------|-------------|---------|
| `text` | Free-form text | Task title, description |
| `textarea` | Multi-line text | Detailed description, code snippets |
| `select` | Dropdown selection | Priority (low, medium, high) |
| `multiselect` | Multiple selections | Tags, categories |
| `date` | Date picker | Due date, deadline |
| `number` | Numeric input | Story points, estimate hours |
| `checkbox` | Boolean value | Is blocking, needs review |

---

## Default Templates

### 1. Bug Report (`bug`)

**Purpose**: Document software bugs with reproduction steps

**Category**: Technical

**Title Pattern**: `🐛 [{{modul}}] {{baslik}}`

**Fields**:

| Field | Type | Required | Default | Options |
|-------|------|----------|---------|---------|
| `baslik` | text | ✅ | - | - |
| `aciklama` | textarea | ✅ | - | - |
| `modul` | text | ✅ | - | - |
| `ortam` | select | ✅ | - | development, staging, production |
| `adimlar` | textarea | ✅ | - | - |
| `beklenen` | textarea | ✅ | - | - |
| `mevcut` | textarea | ✅ | - | - |
| `ekler` | text | ❌ | - | - |
| `cozum` | textarea | ❌ | - | - |
| `oncelik` | select | ✅ | orta | dusuk, orta, yuksek |
| `etiketler` | text | ❌ | bug | - |

**Example Usage**:

```bash
# CLI
gorev task create --template bug \
  --field baslik="Login butonu çalışmıyor" \
  --field modul="auth" \
  --field ortam="production" \
  --field adimlar="1. Login sayfasına git\n2. Email ve şifre gir\n3. Login butonuna tıkla" \
  --field beklenen="Kullanıcı ana sayfaya yönlendirilmeli" \
  --field mevcut="Hiçbir şey olmuyor" \
  --field oncelik="yuksek"
```

**AI Assistant Prompt**:

```
Create a bug task for login button not responding in production
```

### 2. Feature Development (`feature`)

**Purpose**: Plan and track new feature development

**Category**: Feature

**Title Pattern**: `✨ {{baslik}}`

**Fields**:

| Field | Type | Required | Default | Options |
|-------|------|----------|---------|---------|
| `baslik` | text | ✅ | - | - |
| `aciklama` | textarea | ✅ | - | - |
| `hedef_kullanici` | text | ✅ | - | - |
| `kullanici_hikayesi` | textarea | ✅ | - | - |
| `kabul_kriterleri` | textarea | ✅ | - | - |
| `teknik_detaylar` | textarea | ❌ | - | - |
| `tasarim_linki` | text | ❌ | - | - |
| `oncelik` | select | ✅ | orta | dusuk, orta, yuksek |
| `tahmini_sure` | text | ❌ | - | - |
| `etiketler` | text | ❌ | feature | - |

**Example**:

```
Title: Dark mode theme support
Target User: All users
User Story: As a user, I want to switch between light and dark themes so that I can reduce eye strain in low-light environments
Acceptance Criteria:
  - Toggle switch in settings
  - Persists preference to localStorage
  - Affects all UI components
  - Smooth transition animation
```

### 3. Research Task (`research`)

**Purpose**: Investigation, analysis, and proof-of-concept work

**Category**: Research

**Title Pattern**: `🔍 {{baslik}}`

**Fields**:

| Field | Type | Required | Default | Options |
|-------|------|----------|---------|---------|
| `baslik` | text | ✅ | - | - |
| `arastirma_sorusu` | textarea | ✅ | - | - |
| `hipotez` | textarea | ❌ | - | - |
| `yontem` | textarea | ✅ | - | - |
| `basari_kriterleri` | textarea | ✅ | - | - |
| `kaynaklar` | textarea | ❌ | - | - |
| `bulgular` | textarea | ❌ | - | - |
| `oneri` | textarea | ❌ | - | - |
| `oncelik` | select | ✅ | orta | dusuk, orta, yuksek |
| `etiketler` | text | ❌ | research | - |

**Example**:

```
Research Question: Should we migrate from REST to GraphQL?
Hypothesis: GraphQL will reduce API calls by 50% and improve mobile performance
Method:
  1. Benchmark current REST API performance
  2. Build proof-of-concept GraphQL server
  3. Compare query complexity and response times
  4. Analyze client-side caching benefits
Success Criteria:
  - Performance metrics documented
  - POC demonstrates query flexibility
  - Migration effort estimated
```

### 4. Code Refactoring (`refactor`)

**Purpose**: Code quality improvements without changing functionality

**Category**: Technical

**Title Pattern**: `♻️ {{baslik}}`

**Fields**:

| Field | Type | Required | Default | Options |
|-------|------|----------|---------|---------|
| `baslik` | text | ✅ | - | - |
| `mevcut_durum` | textarea | ✅ | - | - |
| `hedef_durum` | textarea | ✅ | - | - |
| `neden` | textarea | ✅ | - | - |
| `dosyalar` | textarea | ✅ | - | - |
| `test_plani` | textarea | ✅ | - | - |
| `oncelik` | select | ✅ | orta | dusuk, orta, yuksek |
| `etiketler` | text | ❌ | refactor | - |

**Example**:

```
Title: Extract authentication logic into middleware
Current State: Auth checks scattered across 15 route handlers
Target State: Centralized auth middleware with role-based access control
Why: Reduce duplication, improve security, easier to maintain
Files: src/routes/*.ts, src/middleware/auth.ts (new)
Test Plan: Unit tests for middleware, integration tests for protected routes
```

### 5. Test Writing (`test`)

**Purpose**: Plan and track test creation

**Category**: Technical

**Title Pattern**: `🧪 {{baslik}}`

**Fields**:

| Field | Type | Required | Default | Options |
|-------|------|----------|---------|---------|
| `baslik` | text | ✅ | - | - |
| `test_turu` | select | ✅ | - | unit, integration, e2e, performance |
| `kapsam` | textarea | ✅ | - | - |
| `test_senaryolari` | textarea | ✅ | - | - |
| `beklenen_kapsama` | text | ❌ | - | - |
| `oncelik` | select | ✅ | orta | dusuk, orta, yuksek |
| `etiketler` | text | ❌ | test | - |

**Example**:

```
Title: API endpoint integration tests
Test Type: integration
Scope: All REST API endpoints in /api/tasks/*
Test Scenarios:
  - GET /api/tasks returns task list with correct pagination
  - POST /api/tasks creates task with valid data
  - POST /api/tasks rejects invalid data with 400 error
  - PUT /api/tasks/:id updates existing task
  - DELETE /api/tasks/:id removes task
Expected Coverage: 90%+ for all API routes
```

### 6. Documentation (`doc`)

**Purpose**: Write or update documentation

**Category**: Process

**Title Pattern**: `📝 {{baslik}}`

**Fields**:

| Field | Type | Required | Default | Options |
|-------|------|----------|---------|---------|
| `baslik` | text | ✅ | - | - |
| `dokuman_turu` | select | ✅ | - | README, API, Guide, Tutorial, Changelog |
| `hedef_okuyucu` | text | ✅ | - | - |
| `kapsam` | textarea | ✅ | - | - |
| `yapisal_tasarim` | textarea | ❌ | - | - |
| `oncelik` | select | ✅ | orta | dusuk, orta, yuksek |
| `etiketler` | text | ❌ | documentation | - |

**Example**:

```
Title: Web UI user guide
Document Type: Guide
Target Audience: End users (developers using Gorev)
Scope:
  - Getting started with Web UI
  - Multi-workspace management
  - Task creation workflows
  - Template usage
  - Troubleshooting common issues
```

---

## Usage

### CLI

#### List Templates

```bash
gorev template list
```

**Output**:

```
Available Templates:

Technical:
  - bug (Bug Raporu)
    ID: 39f28dbd-10f3-454c-8b35-52ae6b7ea391
    Description: Yazılım hatası bildirimi için detaylı template

Feature:
  - feature (Özellik Geliştirme)
    ID: 7a3c9f2e-5d1b-4c8a-9f6e-2b3d4a5c6e7f
    Description: Yeni özellik geliştirme planlaması
...
```

#### Create Task from Template (Using Alias)

```bash
gorev task create --template bug \
  --field baslik="API timeout hatası" \
  --field modul="backend" \
  --field ortam="production" \
  --field adimlar="1. /api/tasks endpoint'ini çağır" \
  --field beklenen="2 saniyede yanıt" \
  --field mevcut="30 saniye bekleyip timeout" \
  --field oncelik="yuksek"
```

#### Create Task from Template (Using UUID)

```bash
gorev task create --template 39f28dbd-10f3-454c-8b35-52ae6b7ea391 \
  --field baslik="Memory leak" \
  --field modul="server" \
  --field ortam="production"
```

### MCP Protocol

#### List Templates

```json
{
  "name": "template_listele",
  "arguments": {}
}
```

**Response**:

```markdown
## 📋 Görev Template'leri

### Teknik

#### Bug Raporu
- **ID:** 39f28dbd-10f3-454c-8b35-52ae6b7ea391
- **Alias:** bug
- **Açıklama:** Yazılım hatası bildirimi için detaylı template
- **Başlık Şablonu:** 🐛 [{{modul}}] {{baslik}}
...
```

#### Create Task from Template

```json
{
  "name": "templateden_gorev_olustur",
  "arguments": {
    "template_alias": "bug",
    "degerler": {
      "baslik": "Login butonu çalışmıyor",
      "aciklama": "Kullanıcı giriş sayfasında login butonu tıklamaya yanıt vermiyor",
      "modul": "auth",
      "ortam": "production",
      "adimlar": "1. Login sayfasına git\n2. Email ve şifre gir\n3. Login butonuna tıkla",
      "beklenen": "Kullanıcı ana sayfaya yönlendirilmeli",
      "mevcut": "Hiçbir şey olmuyor",
      "oncelik": "yuksek"
    }
  }
}
```

**Response**:

```json
{
  "content": [{
    "type": "text",
    "text": "✓ Template kullanılarak görev oluşturuldu: 🐛 [auth] Login butonu çalışmıyor (ID: d7f4e8b9-2a1c-4f5e-9d3b-8c1a2e3f4d5b)"
  }]
}
```

### AI Assistant

#### Natural Language Prompts

```
Create a bug task for the login button issue in production
```

AI assistant automatically:

1. Selects `bug` template
2. Prompts for required fields
3. Fills in defaults
4. Creates task

**More Examples**:

```
Create a feature task for dark mode support
```

```
Create a research task to investigate GraphQL migration
```

```
Create a refactoring task to clean up the authentication code
```

```
Create test cases for the new API endpoints
```

```
Create documentation for the multi-workspace feature
```

### VS Code Extension

#### Template Wizard

1. **Open Command Palette**: `Ctrl+Shift+P` (Windows/Linux) or `Cmd+Shift+P` (macOS)
2. **Run**: `Gorev: Create Task from Template`
3. **Select Template**: Choose from list or type alias
4. **Fill Fields**: Dynamic form with validation
5. **Preview**: Review task before creation
6. **Create**: Task created and added to TreeView

#### Quick Create

Right-click in Tasks TreeView → "Create Task from Template" → Select template

### Web UI

#### Template-Based Creation

1. **Click**: "New Task" button
2. **Select Template**: Grid view of all templates with descriptions
3. **Fill Form**: Auto-generated form based on template schema
4. **Validate**: Real-time field validation
5. **Create**: Task appears in task list immediately

**Features**:

- Template search and filtering
- Required field highlighting
- Default value pre-population
- Select field dropdowns
- Date picker for due dates
- Tag auto-complete

---

## Template Schema

### JSON Structure

```json
{
  "id": "39f28dbd-10f3-454c-8b35-52ae6b7ea391",
  "alias": "bug",
  "isim": "Bug Raporu",
  "kategori": "Teknik",
  "aciklama": "Yazılım hatası bildirimi için detaylı template",
  "baslik_sablonu": "🐛 [{{modul}}] {{baslik}}",
  "alanlar": [
    {
      "isim": "baslik",
      "tip": "text",
      "zorunlu": true,
      "varsayilan": null,
      "aciklama": "Hatanın kısa özeti"
    },
    {
      "isim": "oncelik",
      "tip": "select",
      "zorunlu": true,
      "varsayilan": "orta",
      "secenekler": ["dusuk", "orta", "yuksek"],
      "aciklama": "Hatanın öncelik seviyesi"
    }
  ],
  "olusturulma_tarihi": "2025-01-15T10:00:00Z"
}
```

### Database Schema

```sql
CREATE TABLE gorev_templateleri (
  id TEXT PRIMARY KEY,              -- UUID
  alias TEXT UNIQUE,                -- Short name (bug, feature)
  isim TEXT NOT NULL,               -- Display name
  kategori TEXT,                    -- Category (Teknik, Özellik, etc.)
  aciklama TEXT,                    -- Description
  baslik_sablonu TEXT,              -- Title pattern with {{placeholders}}
  alanlar TEXT NOT NULL,            -- JSON array of field schemas
  olusturulma_tarihi TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  guncelleme_tarihi TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Advanced Features

### Title Pattern Substitution

Templates support **placeholder substitution** in title patterns:

**Pattern**: `🐛 [{{modul}}] {{baslik}}`

**Field Values**:

```json
{
  "baslik": "Login hatası",
  "modul": "auth"
}
```

**Generated Title**: `🐛 [auth] Login hatası`

**Supported Placeholders**:

- `{{fieldname}}`: Direct field value
- `{{FIELDNAME}}`: Uppercase transformation
- `{{fieldname|default:value}}`: Default if empty (future)

### Field Validation

**Required Fields**:

```json
{
  "isim": "baslik",
  "zorunlu": true
}
```

Validation error if missing:

```
Error: Required field 'baslik' is missing
```

**Select Field Options**:

```json
{
  "isim": "oncelik",
  "tip": "select",
  "secenekler": ["dusuk", "orta", "yuksek"]
}
```

Validation error if invalid:

```
Error: Invalid value 'urgent' for field 'oncelik'. Must be one of: dusuk, orta, yuksek
```

**Default Values**:

```json
{
  "isim": "oncelik",
  "varsayilan": "orta"
}
```

Auto-fills if not provided by user.

### Multi-Language Support

Templates have language-specific field names:

**Turkish**:

```json
{
  "isim": "baslik",
  "aciklama": "Görev başlığı"
}
```

**English Translation** (via i18n system):

```json
{
  "name": "title",
  "description": "Task title"
}
```

Language selected via `GOREV_LANG` environment variable.

---

## Custom Templates

### Creating Custom Templates (Future Feature)

**CLI Command** (planned for v0.17.0):

```bash
gorev template create \
  --name "Code Review" \
  --alias "review" \
  --category "Process" \
  --title-pattern "👁️ {{pr_title}}" \
  --field "pr_title:text:required" \
  --field "pr_url:text:required" \
  --field "reviewers:multiselect:required:@team" \
  --field "priority:select:required:low,medium,high:medium"
```

**JSON Import**:

```bash
gorev template import --file my-templates.json
```

**Web UI Template Builder**:

- Drag-and-drop field creation
- Visual field type selection
- Real-time preview
- Export as JSON

### Template Sharing (Future Feature)

**Export**:

```bash
gorev template export --id 39f28dbd... --output bug-template.json
```

**Import**:

```bash
gorev template import --input bug-template.json
```

**Template Marketplace** (planned):

- Community-contributed templates
- Category browsing
- One-click install
- Rating and reviews

---

## Best Practices

### 1. Choose the Right Template

| Scenario | Recommended Template |
|----------|---------------------|
| Software bug found | `bug` |
| New feature request | `feature` |
| Need to investigate something | `research` |
| Code needs improvement | `refactor` |
| Missing tests | `test` |
| Missing documentation | `doc` |

### 2. Fill Required Fields Completely

**Bad**:

```json
{
  "baslik": "Fix bug"
}
```

**Good**:

```json
{
  "baslik": "Login button not responding",
  "modul": "authentication",
  "ortam": "production",
  "adimlar": "1. Navigate to login\n2. Enter credentials\n3. Click submit",
  "beklenen": "User redirected to dashboard",
  "mevcut": "Button click has no effect"
}
```

### 3. Use Consistent Naming

**Module Names**:

- ✅ `auth`, `api`, `frontend`, `backend`
- ❌ `Authentication Module`, `The API`, `Front-End`

**Tags**:

- ✅ `bug`, `security`, `performance`
- ❌ `Bug Report`, `Sec Issue`, `perf`

### 4. Leverage Defaults

Define sensible defaults in templates:

```json
{
  "isim": "oncelik",
  "varsayilan": "orta"
}
```

Reduces user input burden for common cases.

### 5. Template Selection Guidelines

**Use `bug` when**:

- Something is broken
- Unexpected behavior occurs
- Error messages appear

**Use `feature` when**:

- Adding new functionality
- Enhancing existing features
- User stories need tracking

**Use `research` when**:

- Evaluating new technologies
- Performance analysis needed
- Proof-of-concept required

**Use `refactor` when**:

- Code quality issues exist
- Technical debt accumulates
- Maintainability suffers

**Use `test` when**:

- Test coverage gaps identified
- New code needs testing
- Test framework changes

**Use `doc` when**:

- Documentation missing
- README outdated
- API docs needed

---

## Troubleshooting

### Issue: Template Not Found

**Symptoms**:

- "Template not found: bug"
- "Template ID invalid"

**Solutions**:

```bash
# List available templates
gorev template list

# Verify alias exists
gorev template list | grep "bug"

# Use UUID if alias doesn't work
gorev task create --template 39f28dbd-...
```

### Issue: Required Field Missing

**Symptoms**:

- "Required field 'baslik' is missing"

**Solutions**:

```bash
# Check template schema
gorev template show --alias bug

# Provide all required fields
gorev task create --template bug \
  --field baslik="Title" \
  --field modul="Module" \
  --field ortam="production" \
  --field adimlar="Steps" \
  --field beklenen="Expected" \
  --field mevcut="Actual"
```

### Issue: Invalid Select Value

**Symptoms**:

- "Invalid value 'urgent' for field 'oncelik'"

**Solutions**:

```bash
# Check valid options
gorev template show --alias bug | grep oncelik

# Use valid option
--field oncelik="yuksek"  # Not "urgent"
```

---

## Template Aliases Reference

### Quick Reference Card

```
bug       → Bug Raporu           → 🐛 Technical
feature   → Özellik Geliştirme   → ✨ Feature
research  → Araştırma            → 🔍 Research
refactor  → Refactoring          → ♻️ Technical
test      → Test Yazma           → 🧪 Technical
doc       → Dokümantasyon        → 📝 Process
```

### Alias Usage in Different Contexts

**CLI**:

```bash
gorev task create --template bug ...
```

**MCP**:

```json
{
  "name": "templateden_gorev_olustur",
  "arguments": {
    "template_alias": "bug"
  }
}
```

**AI Assistant**:

```
Use the bug template to create a task
```

**Web UI**:

- Template selector shows aliases prominently
- Click "bug" card to select template

---

## Performance

### Benchmarks

| Operation | Time | Notes |
|-----------|------|-------|
| List templates | 2ms | Cached in memory |
| Get template by alias | < 1ms | Hash table lookup |
| Get template by UUID | 3ms | Database query |
| Create task from template | 15ms | Validation + insertion |
| Template field validation | 5ms | Per field |

### Optimization Tips

1. **Use Aliases**: Faster than UUID lookups
2. **Cache Templates**: Extension/Web UI cache templates on startup
3. **Batch Creation**: Create multiple tasks in single transaction
4. **Lazy Loading**: Load template details only when needed

---

## Migration from v0.9.x

### Breaking Change (v0.10.0)

**Removed**: `gorev_olustur` MCP tool (direct task creation)

**Required**: All task creation must use templates

**Before (v0.9.x)**:

```json
{
  "name": "gorev_olustur",
  "arguments": {
    "baslik": "Fix bug",
    "aciklama": "Some bug",
    "oncelik": "yuksek"
  }
}
```

**After (v0.10.0+)**:

```json
{
  "name": "templateden_gorev_olustur",
  "arguments": {
    "template_alias": "bug",
    "degerler": {
      "baslik": "Fix bug",
      "aciklama": "Some bug",
      "modul": "backend",
      "ortam": "production",
      "adimlar": "...",
      "beklenen": "...",
      "mevcut": "...",
      "oncelik": "yuksek"
    }
  }
}
```

**Migration Script**:

```bash
# Export existing tasks
gorev export --output tasks-backup.json

# For each task, determine appropriate template
# Re-create using template system
# (Manual process - no automatic migration)
```

---

## API Reference

### Template Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/templates` | GET | List all templates |
| `/api/templates/:id` | GET | Get template details |
| `/api/templates/by-alias/:alias` | GET | Get template by alias |
| `/api/tasks/from-template` | POST | Create task from template |

**Example**:

```bash
curl -X GET http://localhost:5082/api/templates
```

**Response**:

```json
{
  "templates": [
    {
      "id": "39f28dbd-10f3-454c-8b35-52ae6b7ea391",
      "alias": "bug",
      "name": "Bug Raporu",
      "category": "Teknik",
      "description": "Yazılım hatası bildirimi için detaylı template",
      "titlePattern": "🐛 [{{modul}}] {{baslik}}",
      "fields": [
        {
          "name": "baslik",
          "type": "text",
          "required": true,
          "default": null
        }
      ]
    }
  ]
}
```

---

## Additional Resources

- **MCP Tools Reference**: [MCP Tools Guide](../../legacy/tr/mcp-araclari.md)
- **Web UI Guide**: [Web UI Documentation](web-ui.md)
- **AI Context Management**: [AI Context Guide](ai-context-management.md)
- **GitHub Issues**: https://github.com/msenol/gorev/issues

---

**Need Help?** Open an issue at [GitHub Issues](https://github.com/msenol/gorev/issues) with the `template-system` label.
