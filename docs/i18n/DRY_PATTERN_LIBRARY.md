# DRY Pattern Library for i18n Migration

**Version**: 1.0.0
**Last Updated**: October 11, 2025
**Purpose**: Complete reference guide for migrating hardcoded strings to use existing i18n helpers

---

## Executive Summary

This document provides practical migration patterns using the **existing** i18n infrastructure in Gorev. All helper functions and translation keys referenced here **already exist** in:

- `internal/i18n/helpers.go` (531 lines, 40+ functions)
- `internal/i18n/locales/tr.json` (lines 314-509, complete `common.*` hierarchy)
- `internal/i18n/locales/en.json` (needs sync - Phase 1.6)

**Goal**: Migrate ~1,200 hardcoded strings to use these helpers with **zero duplication** (DRY principle).

---

## Table of Contents

1. [Quick Reference](#quick-reference)
2. [Core Helper Functions](#core-helper-functions)
3. [Migration Patterns](#migration-patterns)
4. [Entity Type Mapping](#entity-type-mapping)
5. [Real-World Examples](#real-world-examples)
6. [Common Mistakes to Avoid](#common-mistakes-to-avoid)

---

## Quick Reference

### Most Common Migrations

| Hardcoded Pattern | Helper Function | Usage |
|-------------------|----------------|-------|
| `"{{param}} gerekli"` | `TRequiredParam()` | `i18n.TRequiredParam("gorev_id")` |
| `"{{entity}} bulunamadı"` | `TEntityNotFound()` | `i18n.TEntityNotFound("task", err)` |
| `"{{entity}} oluşturulamadı"` | `TCreateFailed()` | `i18n.TCreateFailed("task", err)` |
| `"{{entity}} güncellenemedi"` | `TUpdateFailed()` | `i18n.TUpdateFailed("task", err)` |
| `"✓ {{entity}} oluşturuldu"` | `TCreated()` | `i18n.TCreated("task", title, id)` |
| `"geçersiz durum: {{status}}"` | `TInvalidStatus()` | `i18n.TInvalidStatus(status, validStatuses)` |

### Import Statement

```go
import "github.com/msenol/gorev/internal/i18n"
```

---

## Core Helper Functions

### 1. Validation Helpers (lines 176-272)

#### `TRequiredParam(param string) string`
**Translation Key**: `common.validation.required`
**Turkish**: `"{{.Param}} parametresi gerekli"`
**English**: `"{{.Param}} parameter required"`

```go
// BEFORE
return mcp.NewToolResultError("gorev_id parametresi gerekli")

// AFTER
return mcp.NewToolResultError(i18n.TRequiredParam("gorev_id"))
```

#### `TRequiredArray(param string) string`
**Translation Key**: `common.validation.required_array`
**Turkish**: `"{{.Param}} gerekli ve dizi olmalı"`

```go
// BEFORE
return mcp.NewToolResultError("etiketler gerekli ve dizi olmalı")

// AFTER
return mcp.NewToolResultError(i18n.TRequiredArray("etiketler"))
```

#### `TRequiredObject(param string) string`
**Translation Key**: `common.validation.required_object`
**Turkish**: `"{{.Param}} gerekli ve obje olmalı"`

```go
// BEFORE
return mcp.NewToolResultError("updates gerekli ve obje olmalı")

// AFTER
return mcp.NewToolResultError(i18n.TRequiredObject("updates"))
```

#### `TEntityNotFound(entity string, err error) string`
**Translation Key**: `common.validation.not_found`
**Turkish**: `"{{.Entity}} bulunamadı: {{.Error}}"`

```go
// BEFORE
return mcp.NewToolResultError("görev bulunamadı: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TEntityNotFound("task", err))
```

#### `TEntityNotFoundByID(entity, id string) string`
**Translation Key**: `common.validation.not_found_id`
**Turkish**: `"{{.Entity}} bulunamadı (ID: {{.Id}})"`

```go
// BEFORE
return mcp.NewToolResultError(fmt.Sprintf("proje bulunamadı (ID: %s)", id))

// AFTER
return mcp.NewToolResultError(i18n.TEntityNotFoundByID("project", id))
```

#### `TInvalidValue(param, value string, validValues []string) string`
**Translation Key**: `common.validation.invalid_value`
**Turkish**: `"{{.Param}} için geçersiz değer: '{{.Value}}' Geçerli değerler: ..."`

```go
// BEFORE
return mcp.NewToolResultError(fmt.Sprintf("durum için geçersiz değer: '%s'. Geçerli değerler: beklemede, devam-ediyor, tamamlandı", status))

// AFTER
return mcp.NewToolResultError(i18n.TInvalidValue("durum", status, validStatuses))
```

#### `TInvalidStatus(status string, validStatuses []string) string`
**Translation Key**: `common.validation.invalid_status`

```go
// BEFORE
return mcp.NewToolResultError("geçersiz durum: " + status)

// AFTER
return mcp.NewToolResultError(i18n.TInvalidStatus(status, constants.AllTaskStatuses))
```

#### `TInvalidPriority(priority string) string`
**Translation Key**: `common.validation.invalid_priority`

```go
// BEFORE
return mcp.NewToolResultError("geçersiz öncelik: " + priority)

// AFTER
return mcp.NewToolResultError(i18n.TInvalidPriority(priority))
```

#### `TInvalidDate(dateValue string) string`
**Translation Key**: `common.validation.invalid_date`

```go
// BEFORE
return mcp.NewToolResultError("geçersiz tarih formatı: " + date)

// AFTER
return mcp.NewToolResultError(i18n.TInvalidDate(date))
```

### 2. Operation Failure Helpers (lines 275-345)

#### `TOperationFailed(operation, entity string, err error) string`
**Generic helper** - rarely used directly, use specific helpers below instead.

#### `TCreateFailed(entity string, err error) string`
**Translation Key**: `common.operations.create_failed`

```go
// BEFORE
return mcp.NewToolResultError("görev oluşturulamadı: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TCreateFailed("task", err))
```

#### `TUpdateFailed(entity string, err error) string`
**Translation Key**: `common.operations.update_failed`

```go
// BEFORE
return mcp.NewToolResultError("proje güncellenemedi: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TUpdateFailed("project", err))
```

#### `TDeleteFailed(entity string, err error) string`
**Translation Key**: `common.operations.delete_failed`

```go
// BEFORE
return mcp.NewToolResultError("görev silinemedi: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TDeleteFailed("task", err))
```

#### `TFetchFailed(entity string, err error) string`
**Translation Key**: `common.operations.fetch_failed`

```go
// BEFORE
return mcp.NewToolResultError("görevler getirilemedi: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TFetchFailed("task", err))
```

#### `TSaveFailed(entity string, err error) string`
**Translation Key**: `common.operations.save_failed`

```go
// BEFORE
return mcp.NewToolResultError("aktif proje kaydedilemedi: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TSaveFailed("active_project", err))
```

#### `TAddFailed(entity string, err error) string`
**Translation Key**: `common.operations.add_failed`

```go
// BEFORE
return mcp.NewToolResultError("bağımlılık eklenemedi: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TAddFailed("dependency", err))
```

#### `TRemoveFailed(entity string, err error) string`
**Translation Key**: `common.operations.remove_failed`

```go
// BEFORE
return mcp.NewToolResultError("aktif proje kaldırılamadı: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TRemoveFailed("active_project", err))
```

#### `TListFailed(entity string, err error) string`
**Translation Key**: `common.operations.list_failed`

```go
// BEFORE
return mcp.NewToolResultError("projeler listelenirken hata: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TListFailed("project", err))
```

#### `TEditFailed(entity string, err error) string`
**Translation Key**: `common.operations.edit_failed`

```go
// BEFORE
return mcp.NewToolResultError("görev düzenlenemedi: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TEditFailed("task", err))
```

#### `TQueryFailed(entity string, err error) string`
**Translation Key**: `common.operations.query_failed`

```go
// BEFORE
return mcp.NewToolResultError("NLP sorgusu başarısız: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TQueryFailed("nlp_query", err))
```

#### `TParseFailed(entity string, err error) string`
**Translation Key**: `common.operations.parse_failed`

```go
// BEFORE
return mcp.NewToolResultError("JSON parse hatası: " + err.Error())

// AFTER
return mcp.NewToolResultError(i18n.TParseFailed("json", err))
```

### 3. Success Message Helpers (lines 348-392)

#### `TCreated(entity, title, id string) string`
**Translation Key**: `common.success.created`
**Turkish**: `"✓ {{.Entity}} oluşturuldu: {{.Title}} (ID: {{.Id}})"`

```go
// BEFORE
return mcp.NewToolResultText(fmt.Sprintf("✓ Görev oluşturuldu: %s (ID: %s)", title, id))

// AFTER
return mcp.NewToolResultText(i18n.TCreated("task", title, id))
```

#### `TUpdated(entity, details string) string`
**Translation Key**: `common.success.updated`

```go
// BEFORE
return mcp.NewToolResultText("✓ Görev güncellendi: " + details)

// AFTER
return mcp.NewToolResultText(i18n.TUpdated("task", details))
```

#### `TDeleted(entity, title, id string) string`
**Translation Key**: `common.success.deleted`

```go
// BEFORE
return mcp.NewToolResultText(fmt.Sprintf("✓ Görev silindi: %s (ID: %s)", title, id))

// AFTER
return mcp.NewToolResultText(i18n.TDeleted("task", title, id))
```

#### `TSet(entity, details string) string`
**Translation Key**: `common.success.set`

```go
// BEFORE
return mcp.NewToolResultText("✓ Aktif proje ayarlandı: " + name)

// AFTER
return mcp.NewToolResultText(i18n.TSet("active_project", name))
```

#### `TRemoved(entity string) string`
**Translation Key**: `common.success.removed`

```go
// BEFORE
return mcp.NewToolResultText("✓ Aktif proje kaldırıldı")

// AFTER
return mcp.NewToolResultText(i18n.TRemoved("active_project"))
```

#### `TAdded(entity, details string) string`
**Translation Key**: `common.success.added`

```go
// BEFORE
return mcp.NewToolResultText("✓ Bağımlılık eklendi: " + targetTitle)

// AFTER
return mcp.NewToolResultText(i18n.TAdded("dependency", targetTitle))
```

#### `TEdited(entity, title string) string`
**Translation Key**: `common.success.edited`

```go
// BEFORE
return mcp.NewToolResultText("✓ Görev düzenlendi: " + title)

// AFTER
return mcp.NewToolResultText(i18n.TEdited("task", title))
```

### 4. Status/Priority Translation Helpers (lines 504-531)

#### `TStatus(status string) string`
**Translation Keys**: `status.pending`, `status.in_progress`, `status.completed`, `status.cancelled`

```go
// BEFORE
var statusText string
switch status {
case "beklemede": statusText = "Beklemede"
case "devam-ediyor": statusText = "Devam Ediyor"
case "tamamlandı": statusText = "Tamamlandı"
}

// AFTER
statusText := i18n.TStatus(status)
```

#### `TPriority(priority string) string`
**Translation Keys**: `priority.low`, `priority.medium`, `priority.high`

```go
// BEFORE
var priorityText string
switch priority {
case "düşük": priorityText = "Düşük"
case "orta": priorityText = "Orta"
case "yüksek": priorityText = "Yüksek"
}

// AFTER
priorityText := i18n.TPriority(priority)
```

### 5. Markdown Formatting Helpers (lines 454-499)

#### `TLabel(labelKey string) string`
**Translation Key**: `common.labels.{{labelKey}}`

```go
// BEFORE
output += "**Başlık:** " + task.Title

// AFTER
output += i18n.TMarkdownLabel("title", task.Title)
```

#### `TMarkdownLabel(labelKey string, value interface{}) string`
Returns: `"**Label:** value"`

```go
// BEFORE
fmt.Sprintf("**Durum:** %s", status)

// AFTER
i18n.TMarkdownLabel("status", status)
```

#### `TMarkdownHeader(level int, labelKey string) string`
Returns: `"## Label"` (level controls number of #)

```go
// BEFORE
"### Alt Görevler"

// AFTER
i18n.TMarkdownHeader(3, "subtasks")
```

#### `TCount(labelKey string, count int) string`
Returns: `"**Label:** count"`

```go
// BEFORE
fmt.Sprintf("**Toplam:** %d", count)

// AFTER
i18n.TCount("total", count)
```

---

## Migration Patterns

### Pattern 1: Simple Parameter Validation

**Before**:
```go
if gorevID == "" {
    return mcp.NewToolResultError("gorev_id parametresi gerekli")
}
```

**After**:
```go
if gorevID == "" {
    return mcp.NewToolResultError(i18n.TRequiredParam("gorev_id"))
}
```

### Pattern 2: Entity Not Found with Error Context

**Before**:
```go
gorev, err := gorevStore.GetByID(gorevID)
if err != nil {
    return mcp.NewToolResultError("görev bulunamadı: " + err.Error())
}
```

**After**:
```go
gorev, err := gorevStore.GetByID(gorevID)
if err != nil {
    return mcp.NewToolResultError(i18n.TEntityNotFound("task", err))
}
```

### Pattern 3: CRUD Operation Failures

**Before**:
```go
if err := gorevStore.Create(yeniGorev); err != nil {
    return mcp.NewToolResultError("görev oluşturulamadı: " + err.Error())
}
```

**After**:
```go
if err := gorevStore.Create(yeniGorev); err != nil {
    return mcp.NewToolResultError(i18n.TCreateFailed("task", err))
}
```

### Pattern 4: Success Messages

**Before**:
```go
return mcp.NewToolResultText(fmt.Sprintf("✓ Görev oluşturuldu: %s (ID: %s)",
    yeniGorev.Title, yeniGorev.ID))
```

**After**:
```go
return mcp.NewToolResultText(i18n.TCreated("task", yeniGorev.Title, yeniGorev.ID))
```

### Pattern 5: Invalid Value with Valid Options

**Before**:
```go
validStatuses := []string{"beklemede", "devam-ediyor", "tamamlandı"}
if !contains(validStatuses, status) {
    return mcp.NewToolResultError(fmt.Sprintf(
        "durum için geçersiz değer: '%s'. Geçerli değerler: %s",
        status, strings.Join(validStatuses, ", ")))
}
```

**After**:
```go
validStatuses := []string{"beklemede", "devam-ediyor", "tamamlandı"}
if !contains(validStatuses, status) {
    return mcp.NewToolResultError(i18n.TInvalidValue("durum", status, validStatuses))
}
```

### Pattern 6: Complex Validation (Multiple Required Params)

**Before**:
```go
if durum == "" && oncelik == "" {
    return mcp.NewToolResultError("durum veya oncelik parametresi gerekli")
}
```

**After** (Option 1 - Custom message still needed):
```go
if durum == "" && oncelik == "" {
    // This specific pattern needs a dedicated key in tr.json:
    // "validation.one_of_required": "{{.Params}} parametrelerinden en az biri gerekli"
    return mcp.NewToolResultError(i18n.T("common.validation.one_of_required",
        map[string]interface{}{"Params": "durum, oncelik"}))
}
```

**After** (Option 2 - Separate checks, better UX):
```go
// Better: Guide user to provide specific param
params := []string{}
if durum == "" {
    params = append(params, "durum")
}
if oncelik == "" {
    params = append(params, "oncelik")
}
if len(params) > 0 {
    return mcp.NewToolResultError(fmt.Sprintf("%s gerekli",
        strings.Join(params, " veya ")))
}
```

### Pattern 7: Markdown Output with Labels

**Before**:
```go
output := "### Görev Detayları\n\n"
output += fmt.Sprintf("**Başlık:** %s\n", gorev.Title)
output += fmt.Sprintf("**Durum:** %s\n", gorev.Status)
output += fmt.Sprintf("**Öncelik:** %s\n", gorev.Priority)
```

**After**:
```go
output := i18n.TMarkdownHeader(3, "task_details") + "\n\n"
output += i18n.TMarkdownLabel("title", gorev.Title) + "\n"
output += i18n.TMarkdownLabel("status", i18n.TStatus(gorev.Status)) + "\n"
output += i18n.TMarkdownLabel("priority", i18n.TPriority(gorev.Priority)) + "\n"
```

---

## Entity Type Mapping

When using helper functions, pass these entity identifiers (they map to `common.entities.*` keys):

| English Entity | Helper Parameter | TR Translation | EN Translation |
|----------------|------------------|----------------|----------------|
| Task | `"task"` | `görev` | `task` |
| Project | `"project"` | `proje` | `project` |
| Template | `"template"` | `template` | `template` |
| Dependency | `"dependency"` | `bağımlılık` | `dependency` |
| Tag | `"tag"` | `etiket` | `tag` |
| Active Project | `"active_project"` | `aktif proje` | `active project` |
| Subtask | `"subtask"` | `alt görev` | `subtask` |
| AI Context | `"ai_context"` | `AI bağlamı` | `AI context` |
| Filter Profile | `"filter_profile"` | `filtre profili` | `filter profile` |
| NLP Query | `"nlp_query"` | `NLP sorgusu` | `NLP query` |

**Usage Example**:
```go
// For task operations
i18n.TCreateFailed("task", err)           // "görev oluşturulamadı: ..."
i18n.TEntityNotFound("project", err)      // "proje bulunamadı: ..."
i18n.TUpdateFailed("template", err)       // "template güncellenemedi: ..."
```

---

## Real-World Examples

### Example 1: handlers.go:632 - Remove Active Project

**BEFORE**:
```go
func (h *Handler) handleAktifProjeKaldir(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... workspace loading ...

    if err := ws.ClearActiveProject(); err != nil {
        return mcp.NewToolResultError("aktif proje kaldırılamadı: " + err.Error()), nil
    }

    return mcp.NewToolResultText("✓ Aktif proje kaldırıldı"), nil
}
```

**AFTER**:
```go
func (h *Handler) handleAktifProjeKaldir(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... workspace loading ...

    if err := ws.ClearActiveProject(); err != nil {
        return mcp.NewToolResultError(i18n.TRemoveFailed("active_project", err)), nil
    }

    return mcp.NewToolResultText(i18n.TRemoved("active_project")), nil
}
```

**Changes**: 2 strings migrated using `TRemoveFailed()` and `TRemoved()`

---

### Example 2: handlers.go:652 - Bulk Update Validation

**BEFORE**:
```go
func (h *Handler) handleGorevGuncellemeTopluca(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... argument parsing ...

    durum := getString(updates, "durum")
    oncelik := getString(updates, "oncelik")

    if durum == "" && oncelik == "" {
        return mcp.NewToolResultError("durum veya oncelik parametresi gerekli"), nil
    }

    validStatuses := []string{"beklemede", "devam-ediyor", "tamamlandı", "iptal"}
    if durum != "" && !contains(validStatuses, durum) {
        return mcp.NewToolResultError(fmt.Sprintf(
            "geçersiz durum: '%s'. Geçerli değerler: %s",
            durum, strings.Join(validStatuses, ", "))), nil
    }

    // ... rest of logic ...
}
```

**AFTER**:
```go
func (h *Handler) handleGorevGuncellemeTopluca(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... argument parsing ...

    durum := getString(updates, "durum")
    oncelik := getString(updates, "oncelik")

    if durum == "" && oncelik == "" {
        // Custom validation pattern - add to tr.json if not exists
        return mcp.NewToolResultError(i18n.T("common.validation.one_of_required",
            map[string]interface{}{"Params": "durum, oncelik"})), nil
    }

    validStatuses := []string{"beklemede", "devam-ediyor", "tamamlandı", "iptal"}
    if durum != "" && !contains(validStatuses, durum) {
        return mcp.NewToolResultError(i18n.TInvalidValue("durum", durum, validStatuses)), nil
    }

    // ... rest of logic ...
}
```

**Changes**: 2 strings migrated, 1 using custom pattern, 1 using `TInvalidValue()`

---

### Example 3: handlers.go:873 - Edit Task Complex Validation

**BEFORE**:
```go
func (h *Handler) handleGorevDuzenle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... argument parsing ...

    if baslik == "" && aciklama == "" && durum == "" && oncelik == "" &&
       sonTarih == "" && len(etiketlerArray) == 0 && yeniParentID == "" {
        return mcp.NewToolResultError(
            "en az bir düzenleme alanı belirtilmeli (baslik, aciklama, durum, oncelik, son_tarih, etiketler, yeni_parent_id)"), nil
    }

    // ... fetch task ...
    gorev, err := gorevStore.GetByID(gorevID)
    if err != nil {
        return mcp.NewToolResultError("görev bulunamadı: " + err.Error()), nil
    }

    // ... validation logic ...
    if durum != "" {
        validStatuses := []string{"beklemede", "devam-ediyor", "tamamlandı", "iptal"}
        if !contains(validStatuses, durum) {
            return mcp.NewToolResultError(fmt.Sprintf(
                "geçersiz durum: '%s'. Geçerli değerler: %s",
                durum, strings.Join(validStatuses, ", "))), nil
        }
    }

    // ... update logic ...
    if err := gorevStore.Update(gorev); err != nil {
        return mcp.NewToolResultError("görev güncellenemedi: " + err.Error()), nil
    }

    return mcp.NewToolResultText(fmt.Sprintf("✓ Görev düzenlendi: %s", gorev.Title)), nil
}
```

**AFTER**:
```go
func (h *Handler) handleGorevDuzenle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // ... argument parsing ...

    if baslik == "" && aciklama == "" && durum == "" && oncelik == "" &&
       sonTarih == "" && len(etiketlerArray) == 0 && yeniParentID == "" {
        // Complex validation - needs dedicated key:
        // "validation.at_least_one_field": "En az bir alan belirtilmeli: {{.Fields}}"
        return mcp.NewToolResultError(i18n.T("common.validation.at_least_one_field",
            map[string]interface{}{
                "Fields": "baslik, aciklama, durum, oncelik, son_tarih, etiketler, yeni_parent_id",
            })), nil
    }

    // ... fetch task ...
    gorev, err := gorevStore.GetByID(gorevID)
    if err != nil {
        return mcp.NewToolResultError(i18n.TEntityNotFound("task", err)), nil
    }

    // ... validation logic ...
    if durum != "" {
        validStatuses := []string{"beklemede", "devam-ediyor", "tamamlandı", "iptal"}
        if !contains(validStatuses, durum) {
            return mcp.NewToolResultError(i18n.TInvalidValue("durum", durum, validStatuses)), nil
        }
    }

    // ... update logic ...
    if err := gorevStore.Update(gorev); err != nil {
        return mcp.NewToolResultError(i18n.TUpdateFailed("task", err)), nil
    }

    return mcp.NewToolResultText(i18n.TEdited("task", gorev.Title)), nil
}
```

**Changes**: 5 strings migrated, 1 custom key needed, 4 using helpers

---

### Example 4: Complex Markdown Output

**BEFORE** (typical summary output):
```go
func formatTaskSummary(gorev *modeller.Gorev) string {
    output := "### Görev Detayları\n\n"
    output += fmt.Sprintf("**ID:** %s\n", gorev.ID)
    output += fmt.Sprintf("**Başlık:** %s\n", gorev.Title)
    output += fmt.Sprintf("**Durum:** %s\n", gorev.Status)
    output += fmt.Sprintf("**Öncelik:** %s\n", gorev.Priority)

    if gorev.Description != "" {
        output += fmt.Sprintf("**Açıklama:** %s\n", gorev.Description)
    }

    if gorev.DueDate != nil {
        output += fmt.Sprintf("**Son Tarih:** %s\n", gorev.DueDate.Format("2006-01-02"))
    }

    if len(gorev.Tags) > 0 {
        output += fmt.Sprintf("**Etiketler:** %s\n", strings.Join(gorev.Tags, ", "))
    }

    return output
}
```

**AFTER**:
```go
func formatTaskSummary(gorev *modeller.Gorev) string {
    output := i18n.TMarkdownHeader(3, "task_details") + "\n\n"
    output += i18n.TMarkdownLabel("id", gorev.ID) + "\n"
    output += i18n.TMarkdownLabel("title", gorev.Title) + "\n"
    output += i18n.TMarkdownLabel("status", i18n.TStatus(gorev.Status)) + "\n"
    output += i18n.TMarkdownLabel("priority", i18n.TPriority(gorev.Priority)) + "\n"

    if gorev.Description != "" {
        output += i18n.TMarkdownLabel("description", gorev.Description) + "\n"
    }

    if gorev.DueDate != nil {
        output += i18n.TMarkdownLabel("due_date", gorev.DueDate.Format("2006-01-02")) + "\n"
    }

    if len(gorev.Tags) > 0 {
        output += i18n.TMarkdownLabel("tags", strings.Join(gorev.Tags, ", ")) + "\n"
    }

    return output
}
```

**Changes**: 8 strings migrated, all using markdown helpers + `TStatus()`/`TPriority()`

---

## Common Mistakes to Avoid

### ❌ Mistake 1: Creating New Translation Keys Instead of Using Helpers

**DON'T DO THIS**:
```go
// Creating new specific keys for every variation
return mcp.NewToolResultError(i18n.T("errors.task_not_found", nil))
return mcp.NewToolResultError(i18n.T("errors.project_not_found", nil))
return mcp.NewToolResultError(i18n.T("errors.template_not_found", nil))
```

**DO THIS INSTEAD**:
```go
// Use the generic helper with entity parameter
return mcp.NewToolResultError(i18n.TEntityNotFound("task", err))
return mcp.NewToolResultError(i18n.TEntityNotFound("project", err))
return mcp.NewToolResultError(i18n.TEntityNotFound("template", err))
```

**Why**: DRY principle - one translation key serves all entities via parameter injection.

---

### ❌ Mistake 2: Ignoring Error Context

**DON'T DO THIS**:
```go
if err != nil {
    return mcp.NewToolResultError(i18n.T("errors.generic_failure", nil))
}
```

**DO THIS INSTEAD**:
```go
if err != nil {
    return mcp.NewToolResultError(i18n.TCreateFailed("task", err))
}
```

**Why**: Error context is essential for debugging - helpers include `{{.Error}}` in translations.

---

### ❌ Mistake 3: Hardcoding Entity Names in English

**DON'T DO THIS**:
```go
return mcp.NewToolResultError(i18n.TEntityNotFound("Görev", err)) // Hardcoded Turkish
```

**DO THIS INSTEAD**:
```go
return mcp.NewToolResultError(i18n.TEntityNotFound("task", err)) // Key identifier
```

**Why**: Helper resolves `"task"` → `i18n.T("common.entities.task")` → `"görev"` (TR) or `"task"` (EN).

---

### ❌ Mistake 4: Mixing Direct T() Calls with Helpers

**DON'T DO THIS** (inconsistent):
```go
if gorevID == "" {
    return mcp.NewToolResultError(i18n.T("common.validation.required",
        map[string]interface{}{"Param": "gorev_id"}))
}
```

**DO THIS INSTEAD** (consistent):
```go
if gorevID == "" {
    return mcp.NewToolResultError(i18n.TRequiredParam("gorev_id"))
}
```

**Why**: Helpers provide consistent usage patterns and reduce boilerplate.

---

### ❌ Mistake 5: Not Using Status/Priority Helpers in Output

**DON'T DO THIS**:
```go
output += fmt.Sprintf("**Durum:** %s\n", gorev.Status) // "beklemede" (raw value)
```

**DO THIS INSTEAD**:
```go
output += i18n.TMarkdownLabel("status", i18n.TStatus(gorev.Status)) // "Beklemede" (localized)
```

**Why**: Status/priority values are stored as constants but should display localized.

---

### ❌ Mistake 6: Duplicating Label Strings

**DON'T DO THIS**:
```go
output += "**Başlık:** " + gorev.Title + "\n"
output += "**Açıklama:** " + gorev.Description + "\n"
output += "**Durum:** " + gorev.Status + "\n"
```

**DO THIS INSTEAD**:
```go
output += i18n.TMarkdownLabel("title", gorev.Title) + "\n"
output += i18n.TMarkdownLabel("description", gorev.Description) + "\n"
output += i18n.TMarkdownLabel("status", i18n.TStatus(gorev.Status)) + "\n"
```

**Why**: Labels are defined once in `common.labels.*` and reused everywhere.

---

### ❌ Mistake 7: Creating Single-Use Translation Keys

**DON'T ADD THIS to tr.json**:
```json
{
  "tools": {
    "gorev_olustur": {
      "task_created_successfully": "✓ Görev başarıyla oluşturuldu: {{.Title}}"
    },
    "proje_olustur": {
      "project_created_successfully": "✓ Proje başarıyla oluşturuldu: {{.Name}}"
    }
  }
}
```

**USE EXISTING KEY INSTEAD**:
```json
{
  "common": {
    "success": {
      "created": "✓ {{.Entity}} oluşturuldu: {{.Title}} (ID: {{.Id}})"
    }
  }
}
```

**Why**: One `common.success.created` key serves all entity types - pure DRY.

---

## Validation Checklist

Before committing i18n migration changes, verify:

- [ ] **No new translation keys added** unless absolutely necessary (99% should use existing `common.*`)
- [ ] **All helpers imported**: `import "github.com/msenol/gorev/internal/i18n"`
- [ ] **Entity identifiers correct**: `"task"`, `"project"`, not `"Görev"`, `"Proje"`
- [ ] **Error context preserved**: All `TOperationFailed()` calls include `err` parameter
- [ ] **Status/Priority localized**: Use `TStatus()`/`TPriority()` in output, not raw values
- [ ] **Markdown labels consistent**: Use `TMarkdownLabel()`, not hardcoded `"**Label:**"`
- [ ] **No duplicated patterns**: Check if existing helper covers the case
- [ ] **Tests updated**: Integration tests should verify translated output

---

## Additional Resources

- **Translation Files**:
  - `internal/i18n/locales/tr.json` (lines 314-509: `common.*` structure)
  - `internal/i18n/locales/en.json` (to be synced in Phase 1.6)

- **Helper Functions**:
  - `internal/i18n/helpers.go` (lines 173-531: all DRY helpers)
  - `internal/i18n/i18n.go` (base `T()` and `TCommon()` functions)

- **Entity Constants**:
  - `internal/constants/constants.go` (status/priority constants)

- **Task Document**:
  - `docs/tasks/I18N_TRANSLATION_REFACTORING.md` (complete migration plan)

---

## Quick Migration Workflow

1. **Identify hardcoded string** in code
2. **Check this document** for matching pattern
3. **Use existing helper** if available (99% of cases)
4. **Only if truly unique**: Add to `common.*` structure with DRY pattern
5. **Test bilingual output**: `GOREV_LANG=tr` and `GOREV_LANG=en`
6. **Verify no duplication**: Search tr.json for similar keys before adding new ones

---

## Summary Statistics

- **Total Helper Functions**: 40+
- **Common Translation Keys**: 195+ lines (tr.json:314-509)
- **Entity Types Supported**: 10+
- **DRY Principle**: One key serves all entities via parameters
- **Estimated Migration Impact**: ~1,200 hardcoded strings → 40 helper functions

---

**Last Updated**: October 11, 2025
**Document Version**: 1.0.0
**Status**: Complete - Ready for Phase 2 Migration
