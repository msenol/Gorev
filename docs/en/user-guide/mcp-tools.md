# MCP Tools Reference

Complete reference for all 25+ MCP tools provided by Gorev.

> ⚠️ **BREAKING CHANGE (v0.10.0)**: The `gorev_olustur` tool is no longer available! Template usage is now mandatory. See [gorev_template_olustur](#gorev_template_olustur) for details.

> **Note**: For AI Context Management tools, see [AI MCP Tools Documentation](../mcp-araclari-ai.md).

## 📋 Tool Categories

### Core Task Management
- [gorev_template_olustur](#gorev_template_olustur) - Create task from template
- [gorev_listele](#gorev_listele) - List tasks with filtering
- [gorev_detay](#gorev_detay) - Get task details
- [gorev_guncelle](#gorev_guncelle) - Update task properties
- [gorev_sil](#gorev_sil) - Delete task
- [gorev_durum_degistir](#gorev_durum_degistir) - Change task status

### Subtask & Hierarchy
- [gorev_alt_olustur](#gorev_alt_olustur) - Create subtask
- [gorev_ust_degistir](#gorev_ust_degistir) - Change parent task
- [gorev_hiyerarsi_goster](#gorev_hiyerarsi_goster) - Show task hierarchy

### Project Management
- [proje_olustur](#proje_olustur) - Create project
- [proje_listele](#proje_listele) - List projects
- [proje_aktif_yap](#proje_aktif_yap) - Set active project
- [aktif_proje_goster](#aktif_proje_goster) - Show active project

### Templates & Search
- [template_listele](#template_listele) - List templates
- [gorev_ara](#gorev_ara) - Search tasks
- [etiket_listele](#etiket_listele) - List all tags

### Advanced Features
- [gorev_etiket_ekle](#gorev_etiket_ekle) - Add tags
- [gorev_son_tarih](#gorev_son_tarih) - Set due date
- [gorev_bagimliligi_ekle](#gorev_bagimliligi_ekle) - Add dependencies
- [ozet_goster](#ozet_goster) - Show summary

---

## gorev_template_olustur

Create a new task using a template. This is the **primary method** for task creation since v0.10.0.

### Parameters

| Parameter | Type | Required | Description | Examples |
|-----------|------|----------|-------------|----------|
| `template` | string | ✅ | Template name | `bug-report`, `feature`, `task`, `meeting`, `research` |
| `title` | string | ✅ | Task title | "Fix login bug", "Add search feature" |
| `priority` | string | ❌ | Priority level | `low`, `medium`, `high` |
| `due_date` | string | ❌ | Due date (YYYY-MM-DD) | "2025-08-20" |
| `tags` | string | ❌ | Comma-separated tags | "bug,urgent,frontend" |
| `description` | string | ❌ | Additional description | Markdown supported |
| `project_id` | number | ❌ | Project ID | 1, 2, 3 |

### Available Templates

| Template | Description | Use Case |
|----------|-------------|----------|
| `bug-report` | Bug fixes and issues | Software bugs, UI issues |
| `feature` | New features | Feature requests, enhancements |
| `task` | General tasks | General work items |
| `meeting` | Meetings and planning | Meeting prep, agenda |
| `research` | Research tasks | Investigation, analysis |

### Examples

**Basic task:**
```json
{
  "name": "gorev_template_olustur",
  "arguments": {
    "template": "task",
    "title": "Update documentation"
  }
}
```

**Bug report with details:**
```json
{
  "name": "gorev_template_olustur", 
  "arguments": {
    "template": "bug-report",
    "title": "Login page not responsive",
    "priority": "high",
    "due_date": "2025-08-20",
    "tags": "bug,frontend,urgent",
    "description": "Login form breaks on mobile devices"
  }
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "✅ Task created successfully: Login page not responsive (#42)\n📊 Project: Web Development\n🏷️ Tags: bug, frontend, urgent\n📅 Due: 2025-08-20"
  }]
}
```

---

## gorev_listele

List tasks with flexible filtering and sorting options.

### Parameters

| Parameter | Type | Required | Description | Values |
|-----------|------|----------|-------------|--------|
| `durum` | string | ❌ | Filter by status | `beklemede`, `devam_ediyor`, `tamamlandi` |
| `oncelik` | string | ❌ | Filter by priority | `dusuk`, `orta`, `yuksek` |
| `tum_projeler` | boolean | ❌ | Show all projects | `true`, `false` |
| `sirala` | string | ❌ | Sort order | `son_tarih_asc`, `son_tarih_desc` |
| `filtre` | string | ❌ | Time filter | `acil` (due in 7 days), `gecmis` (overdue) |
| `etiket` | string | ❌ | Filter by tag | Any tag name |
| `limit` | number | ❌ | Max results | Default: 50 |
| `offset` | number | ❌ | Skip results | Default: 0 |

### Examples

**All tasks:**
```json
{
  "name": "gorev_listele",
  "arguments": {}
}
```

**In-progress tasks:**
```json
{
  "name": "gorev_listele",
  "arguments": {
    "durum": "devam_ediyor"
  }
}
```

**Urgent tasks (due in 7 days):**
```json
{
  "name": "gorev_listele",
  "arguments": {
    "filtre": "acil"
  }
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "📋 **Tasks List** (5 tasks)\n\n🔄 **In Progress**\n• #42 Login page fix (High) 📅 Due: 2025-08-20\n• #43 Add search feature (Medium)\n\n⏳ **Pending**\n• #44 Update docs (Low)\n• #45 Code review (Medium)\n\n✅ **Completed**\n• #41 Bug fix deployment (High)"
  }]
}
```

---

## gorev_detay

Get detailed information about a specific task.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | number | ✅ | Task ID |

### Example

```json
{
  "name": "gorev_detay",
  "arguments": {
    "id": 42
  }
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "# 🔧 Task Details\n\n**ID:** 42\n**Title:** Login page fix\n**Status:** In Progress 🔄\n**Priority:** High 🔴\n**Project:** Web Development\n**Tags:** bug, frontend, urgent\n**Due Date:** 2025-08-20\n**Created:** 2025-08-13 10:30\n**Updated:** 2025-08-13 14:45\n\n## Description\nLogin form breaks on mobile devices. Need to fix responsive CSS.\n\n## Subtasks (2/3 completed)\n✅ Identify CSS issues\n✅ Fix mobile viewport\n⏳ Test on different devices"
  }]
}
```

---

## gorev_guncelle

Update task properties like status, priority, or description.

### Parameters

| Parameter | Type | Required | Description | Values |
|-----------|------|----------|-------------|--------|
| `id` | number | ✅ | Task ID | |
| `durum` | string | ❌ | New status | `beklemede`, `devam_ediyor`, `tamamlandi` |
| `oncelik` | string | ❌ | New priority | `dusuk`, `orta`, `yuksek` |
| `baslik` | string | ❌ | New title | |
| `aciklama` | string | ❌ | New description | Markdown supported |

### Example

```json
{
  "name": "gorev_guncelle",
  "arguments": {
    "id": 42,
    "durum": "tamamlandi",
    "oncelik": "orta"
  }
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "✅ Task updated successfully: Login page fix (#42)\n📊 Status: Completed\n⚡ Priority: Medium"
  }]
}
```

---

## proje_olustur

Create a new project for organizing tasks.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `isim` | string | ✅ | Project name |
| `tanim` | string | ❌ | Project description |
| `aktif_yap` | boolean | ❌ | Set as active project |

### Example

```json
{
  "name": "proje_olustur",
  "arguments": {
    "isim": "Mobile App v2.0",
    "tanim": "Next generation mobile application",
    "aktif_yap": true
  }
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "✅ Project created: Mobile App v2.0 (#3)\n📁 Description: Next generation mobile application\n⭐ Set as active project"
  }]
}
```

---

## template_listele

List all available task templates.

### Parameters

None required.

### Example

```json
{
  "name": "template_listele",
  "arguments": {}
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "📋 **Available Templates**\n\n🐛 **bug-report** - Bug fixes and issue tracking\n• Fields: module, environment, steps, expected\n• Example: \"Fix login authentication error\"\n\n✨ **feature** - New features and enhancements\n• Fields: requirements, acceptance_criteria, impact\n• Example: \"Add dark mode support\"\n\n📝 **task** - General tasks and activities\n• Fields: category, effort, dependencies\n• Example: \"Update project documentation\"\n\n📅 **meeting** - Meeting planning and notes\n• Fields: attendees, agenda, location, duration\n• Example: \"Sprint planning meeting\"\n\n🔍 **research** - Research and investigation\n• Fields: scope, methodology, deliverables\n• Example: \"Market analysis for new feature\""
  }]
}
```

---

## gorev_ara

Search tasks using natural language or keywords.

### Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `sorgu` | string | ✅ | Search query |
| `limit` | number | ❌ | Max results |

### Example

```json
{
  "name": "gorev_ara",
  "arguments": {
    "sorgu": "login bug mobile"
  }
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "🔍 **Search Results** for \"login bug mobile\" (3 matches)\n\n• #42 **Login page fix** (High) - bug, frontend, urgent\n  📱 Login form breaks on mobile devices\n\n• #38 **Mobile authentication** (Medium) - mobile, auth\n  🔐 Implement biometric login for mobile app\n\n• #35 **Bug tracker setup** (Low) - tools, tracking\n  🐛 Set up bug tracking system for mobile team"
  }]
}
```

---

## etiket_listele

List all unique tags used in tasks.

### Parameters

None required.

### Example

```json
{
  "name": "etiket_listele",
  "arguments": {}
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "🏷️ **All Tags** (15 tags)\n\n**Frequent:**\n• bug (8 tasks)\n• frontend (6 tasks)\n• urgent (4 tasks)\n\n**By Category:**\n• **Tech:** backend, frontend, mobile, api, database\n• **Priority:** urgent, important, low-priority\n• **Type:** bug, feature, enhancement, refactor\n• **Status:** blocked, review, testing"
  }]
}
```

---

## ozet_goster

Show system summary with task statistics and project overview.

### Parameters

None required.

### Example

```json
{
  "name": "ozet_goster",
  "arguments": {}
}
```

### Response

```json
{
  "content": [{
    "type": "text",
    "text": "📊 **Gorev System Summary**\n\n**📋 Tasks Overview**\n• Total: 45 tasks\n• Pending: 12 tasks\n• In Progress: 8 tasks\n• Completed: 25 tasks\n• Overdue: 3 tasks ⚠️\n\n**📁 Projects**\n• Total: 3 projects\n• Active: Web Development (15 tasks)\n• Mobile App v2.0 (12 tasks)\n• Documentation (18 tasks)\n\n**⚡ Priority Distribution**\n• High: 8 tasks (18%)\n• Medium: 22 tasks (49%)\n• Low: 15 tasks (33%)\n\n**🏷️ Most Used Tags**\n• bug (8), frontend (6), urgent (4)\n\n**📅 Due This Week**\n• 5 tasks due in next 7 days\n• 2 tasks overdue"
  }]
}
```

---

## Advanced Tools

### gorev_alt_olustur

Create a subtask under an existing task.

```json
{
  "name": "gorev_alt_olustur",
  "arguments": {
    "ust_gorev_id": 42,
    "template": "task",
    "title": "Test mobile responsiveness"
  }
}
```

### gorev_etiket_ekle

Add tags to a task.

```json
{
  "name": "gorev_etiket_ekle",
  "arguments": {
    "id": 42,
    "etiketler": "mobile,testing,qa"
  }
}
```

### gorev_son_tarih

Set or update task due date.

```json
{
  "name": "gorev_son_tarih",
  "arguments": {
    "id": 42,
    "son_tarih": "2025-08-25"
  }
}
```

### gorev_bagimliligi_ekle

Create dependency between tasks.

```json
{
  "name": "gorev_bagimliligi_ekle",
  "arguments": {
    "kaynak_id": 42,
    "hedef_id": 43,
    "tip": "tamamlanmali"
  }
}
```

## Error Handling

All tools return error messages in a consistent format:

```json
{
  "content": [{
    "type": "text",
    "text": "❌ Error: Task not found (ID: 999)"
  }]
}
```

Common error types:
- **Not Found**: Resource doesn't exist
- **Validation Error**: Invalid parameters
- **Constraint Error**: Business rule violation
- **Permission Error**: Access denied

## Migration from v0.9.x

If you're upgrading from v0.9.x, replace `gorev_olustur` calls with `gorev_template_olustur`:

**Old way:**
```json
{"name": "gorev_olustur", "arguments": {"baslik": "Fix bug", "oncelik": "yuksek"}}
```

**New way:**
```json
{"name": "gorev_template_olustur", "arguments": {"template": "bug-report", "title": "Fix bug", "priority": "high"}}
```

## Tips for AI Assistants

1. **Always use templates** for task creation
2. **Check task details** before making updates
3. **Use search** to find existing similar tasks
4. **Set due dates** for time-sensitive tasks
5. **Add relevant tags** for better organization
6. **Create subtasks** for complex work breakdown
7. **Set active project** for better context

---

*This MCP tools reference was created with assistance from Claude (Anthropic)*