# 🌍 i18n Translation Refactoring Task

**Created:** 2025-10-11
**Status:** 🚧 In Progress
**Assignee:** AI Assistant
**Priority:** HIGH
**Effort:** 24-34 hours (3-4 days)

## 📊 Executive Summary

Complete migration of all hardcoded strings (Turkish and English) to the i18n translation system with DRY (Don't Repeat Yourself) principles at the core. This ensures maintainability, consistency, and international readiness.

---

## 🎯 Objectives

### Primary Goals
1. **100% i18n Coverage**: ALL user-facing strings under translation system
2. **DRY Compliance**: Eliminate redundant translation keys through smart patterns
3. **Hierarchical Structure**: Logical, maintainable key organization
4. **Bilingual Parity**: Perfect TR ↔ EN translation coverage
5. **Zero Regressions**: Maintain all functionality and tests

### Success Metrics
- ✅ 0 hardcoded user-facing strings in codebase
- ✅ 100% TR/EN translation parity
- ✅ All 70 tests passing
- ✅ Translation key reuse rate > 60%
- ✅ Average key path depth ≤ 4 levels

---

## 📈 Current State Analysis

### Existing i18n Infrastructure
```
✅ Library: nicksnyder/go-i18n/v2
✅ Manager: internal/i18n/manager.go
✅ Helper: internal/i18n/helpers.go
✅ Locales: tr.json (695 lines), en.json (684 lines)
✅ Function: i18n.T(key, templateData)
```

### Hardcoded Strings Inventory

**Total Identified:** ~1,200+ strings across 28 files

| Category | File Count | Est. Strings | Priority |
|----------|------------|--------------|----------|
| MCP Handlers | 1 | ~343 | P0 (Critical) |
| Core Business Logic | 18 | ~600 | P0 (Critical) |
| API Layer | 3 | ~80 | P1 (High) |
| Constants | 2 | ~50 | P1 (High) |
| Helpers/Formatters | 4 | ~100 | P2 (Medium) |
| Test Code | N/A | ~200 | P3 (Low - skip) |

### Most Affected Files
```
internal/mcp/handlers.go              ~343 strings  ⚠️  CRITICAL
internal/gorev/veri_yonetici.go       ~120 strings
internal/gorev/template_yonetici.go   ~80 strings
internal/gorev/is_yonetici.go         ~90 strings
internal/gorev/batch_processor.go     ~60 strings
internal/api/server.go                ~40 strings
internal/constants/messages.go        ~35 strings
```

---

## 🔑 DRY Translation Key Design

### Core Principles

1. **Reusable Components**: Common patterns become template-based keys
2. **Entity-Agnostic**: `{{.Entity}}` instead of hardcoded "task", "project"
3. **Operation-Agnostic**: `{{.Operation}}` instead of hardcoded actions
4. **Parameter Injection**: All dynamic data via `{{.Param}}`
5. **Hierarchy Depth**: Max 4 levels (e.g., `common.validation.required`)

### Proposed Key Hierarchy

```
root/
├── common/                    # DRY patterns (MOST REUSED)
│   ├── validation/            # Input validation messages
│   │   ├── required           # "{{.Param}} gerekli"
│   │   ├── invalid            # "{{.Param}} için geçersiz değer"
│   │   ├── not_found          # "{{.Entity}} bulunamadı"
│   │   ├── already_exists     # "{{.Entity}} zaten mevcut"
│   │   ├── empty              # "{{.Field}} boş olamaz"
│   │   └── format             # "Geçersiz {{.Type}} formatı"
│   ├── operations/            # CRUD operation results
│   │   ├── create_success     # "✓ {{.Entity}} oluşturuldu"
│   │   ├── create_failed      # "{{.Entity}} oluşturulamadı"
│   │   ├── update_success     # "✓ {{.Entity}} güncellendi"
│   │   ├── update_failed      # "{{.Entity}} güncellenemedi"
│   │   ├── delete_success     # "✓ {{.Entity}} silindi"
│   │   ├── delete_failed      # "{{.Entity}} silinemedi"
│   │   ├── fetch_failed       # "{{.Entity}} alınamadı"
│   │   └── save_failed        # "{{.Entity}} kaydedilemedi"
│   ├── entities/              # Domain entities
│   │   ├── task               # "görev"
│   │   ├── project            # "proje"
│   │   ├── template           # "template"
│   │   ├── tag                # "etiket"
│   │   ├── dependency         # "bağımlılık"
│   │   └── subtask            # "alt görev"
│   ├── fields/                # Common field descriptions
│   │   ├── id                 # "{{.Entity}} ID'si"
│   │   ├── title              # "{{.Entity}} başlığı"
│   │   ├── description        # "{{.Entity}} açıklaması"
│   │   ├── status             # "{{.Entity}} durumu"
│   │   └── priority           # "{{.Entity}} önceliği"
│   └── labels/                # UI labels
│       ├── created_at         # "Oluşturulma Tarihi"
│       ├── updated_at         # "Güncelleme Tarihi"
│       ├── due_date           # "Son Tarih"
│       └── count              # "{{.Entity}} sayısı"
├── mcp/                       # MCP-specific strings
│   ├── handlers/              # Tool handler messages
│   ├── tools/                 # Tool descriptions
│   └── params/                # Parameter descriptions
├── api/                       # API layer strings
│   ├── errors/                # HTTP error responses
│   ├── workspace/             # Workspace messages
│   └── health/                # Health check messages
├── business/                  # Business logic strings
│   ├── task/                  # Task operations
│   ├── project/               # Project operations
│   ├── template/              # Template operations
│   └── ai/                    # AI context messages
└── cli/                       # CLI interface strings
    ├── commands/              # Command descriptions
    ├── flags/                 # Flag descriptions
    └── output/                # CLI output messages
```

### DRY Pattern Examples

#### ❌ BEFORE (Duplicate Keys)
```json
{
  "taskNotFound": "Görev bulunamadı: {{.Error}}",
  "projectNotFound": "Proje bulunamadı: {{.Error}}",
  "templateNotFound": "Template bulunamadı: {{.Error}}",
  "tagNotFound": "Etiket bulunamadı: {{.Error}}"
}
```

#### ✅ AFTER (Single Reusable Key)
```json
{
  "common": {
    "validation": {
      "not_found": "{{.Entity}} bulunamadı: {{.Error}}"
    },
    "entities": {
      "task": "Görev",
      "project": "Proje",
      "template": "Template",
      "tag": "Etiket"
    }
  }
}
```

#### Usage in Code
```go
// Before
return mcp.NewToolResultError("Görev bulunamadı: " + err.Error())

// After - DRY!
entity := i18n.T("common.entities.task")
return mcp.NewToolResultError(i18n.T("common.validation.not_found", map[string]interface{}{
    "Entity": entity,
    "Error": err.Error(),
}))
```

---

## 📋 Phase-by-Phase Breakdown

### Phase 1: DRY Key Hierarchy Design ⏱️ 4-6 hours

**Objective:** Design and implement reusable translation key structure

#### Tasks
- [ ] 1.1 Analyze all existing 695 TR keys for duplication patterns
- [ ] 1.2 Design `common.*` hierarchy (validation, operations, entities, fields, labels)
- [ ] 1.3 Create DRY pattern library document
- [ ] 1.4 Update `tr.json` with new `common.*` keys
- [ ] 1.5 Update `en.json` with new `common.*` keys
- [ ] 1.6 Add helper functions to `internal/i18n/helpers.go` for DRY usage
- [ ] 1.7 Document key naming conventions

#### Deliverables
```
✅ docs/i18n/DRY_PATTERN_LIBRARY.md
✅ Updated tr.json with ~50 new common.* keys
✅ Updated en.json with ~50 new common.* keys
✅ New helpers in internal/i18n/helpers.go:
   - FormatEntityNotFound(entity, error)
   - FormatOperationSuccess(entity, operation, details)
   - FormatOperationFailed(entity, operation, error)
   - FormatValidationError(param, issue)
```

#### Acceptance Criteria
- All common patterns identified and keyed
- Zero redundant keys in new additions
- Helper functions reduce boilerplate by 70%
- Documentation clear and comprehensive

---

### Phase 2: MCP Handlers Migration ⏱️ 6-8 hours

**Objective:** Migrate `internal/mcp/handlers.go` (~343 strings)

#### File Analysis
```bash
# Hardcoded string patterns found:
return mcp.NewToolResultError("Turkish hardcoded")     # 156 occurrences
return mcp.NewToolResultText("Turkish hardcoded")      # 98 occurrences
fmt.Sprintf("Turkish %s", var)                         # 89 occurrences
```

#### Migration Strategy

**Error Messages (156 strings)**
```go
// BEFORE
return mcp.NewToolResultError("id parametresi gerekli")
return mcp.NewToolResultError("task_ids array formatında olmalı")
return mcp.NewToolResultError(fmt.Sprintf("geçersiz task ID index %d", i))

// AFTER
return mcp.NewToolResultError(i18n.T("common.validation.required", map[string]interface{}{
    "Param": "id",
}))
return mcp.NewToolResultError(i18n.T("common.validation.invalid_format", map[string]interface{}{
    "Param": "task_ids",
    "Type": "array",
}))
return mcp.NewToolResultError(i18n.T("mcp.validation.invalid_index", map[string]interface{}{
    "Index": i,
}))
```

**Success Messages (98 strings)**
```go
// BEFORE
return mcp.NewToolResultText(fmt.Sprintf("✓ Görev düzenlendi: %s", id))
return mcp.NewToolResultText(fmt.Sprintf("✓ Proje oluşturuldu: %s (ID: %s)", proje.Name, proje.ID))

// AFTER
entity := i18n.T("common.entities.task")
return mcp.NewToolResultText(i18n.T("common.operations.update_success", map[string]interface{}{
    "Entity": entity,
    "Title": id,
}))
```

#### Tasks
- [ ] 2.1 Backup handlers.go
- [ ] 2.2 Migrate error messages (lines 400-1500) - ~156 strings
- [ ] 2.3 Migrate success messages (all fmt.Sprintf success patterns) - ~98 strings
- [ ] 2.4 Migrate display text (user-facing output) - ~89 strings
- [ ] 2.5 Run tests: `go test internal/mcp/handlers_test.go`
- [ ] 2.6 Run integration tests: `go test test/integration_test.go`
- [ ] 2.7 Manual smoke test: Call 10 random MCP tools via CLI

#### Deliverables
```
✅ Migrated handlers.go (0 hardcoded user-facing strings)
✅ Added ~150 new translation keys (using DRY patterns)
✅ All tests passing (70/70)
✅ Migration report: Before/After comparison
```

#### Acceptance Criteria
- Zero `grep '".*[ğüşıöçĞÜŞİÖÇ].*"' handlers.go` matches (exclude comments)
- All MCP tools work identically in TR and EN
- No test failures
- Code review: DRY principles followed

---

### Phase 3: Core Business Logic Migration ⏱️ 8-10 hours

**Objective:** Migrate 18 `internal/gorev/*.go` files (~600 strings)

#### File Priority List

| Priority | File | Est. Strings | Notes |
|----------|------|--------------|-------|
| P0 | veri_yonetici.go | ~120 | Database layer errors |
| P0 | is_yonetici.go | ~90 | Business logic |
| P0 | template_yonetici.go | ~80 | Template operations |
| P1 | batch_processor.go | ~60 | Bulk operations |
| P1 | search_engine.go | ~50 | Search messages |
| P1 | export_import.go | ~50 | Export/import |
| P2 | ai_context_yonetici.go | ~40 | AI context |
| P2 | file_watcher.go | ~30 | File watching |
| P2 | suggestion_engine.go | ~30 | Suggestions |
| P3 | modeller.go | ~25 | Status/priority labels |
| P3 | auto_state_manager.go | ~25 | Auto state |

#### Common Patterns in These Files

**Database Errors**
```go
// Pattern: "{{entity}} {{operation}} failed: {{error}}"
// Before: "görev güncellenemedi: %v"
// After: i18n.T("common.operations.update_failed", data)
```

**Validation Errors**
```go
// Pattern: "{{param}} required" or "invalid {{param}}"
// Before: "proje_id belirtilmedi"
// After: i18n.T("common.validation.required", map[string]interface{}{"Param": "proje_id"})
```

**Business Rule Violations**
```go
// Before: "bu görev silinemez, önce {{count}} alt görev silinmeli"
// After: i18n.T("business.task.cannot_delete_has_subtasks", data)
```

#### Tasks
- [ ] 3.1 **veri_yonetici.go** - Database errors (~120 strings)
  - [ ] SQL error messages
  - [ ] Constraint violation messages
  - [ ] Transaction errors
- [ ] 3.2 **is_yonetici.go** - Business logic (~90 strings)
  - [ ] Task state transition errors
  - [ ] Dependency validation
  - [ ] Business rule violations
- [ ] 3.3 **template_yonetici.go** - Template operations (~80 strings)
  - [ ] Template validation
  - [ ] Placeholder errors
  - [ ] Example value messages
- [ ] 3.4 **batch_processor.go** - Bulk operations (~60 strings)
  - [ ] Batch validation
  - [ ] Progress messages
  - [ ] Summary reports
- [ ] 3.5 **Other P1/P2 files** - Remaining 11 files (~200 strings)
  - [ ] search_engine.go
  - [ ] export_import.go
  - [ ] ai_context_yonetici.go
  - [ ] file_watcher.go
  - [ ] suggestion_engine.go
  - [ ] modeller.go
  - [ ] auto_state_manager.go
- [ ] 3.6 Run full test suite after each file
- [ ] 3.7 Integration testing

#### Deliverables
```
✅ 18 migrated business logic files
✅ ~200 new translation keys (many DRY reused)
✅ All 70 tests passing
✅ Business logic functions identically in TR/EN
```

#### Acceptance Criteria
- No hardcoded Turkish strings in non-test files
- Business rules work in both languages
- Error messages clear and helpful
- Performance unchanged

---

### Phase 4: API Layer Migration ⏱️ 4-6 hours

**Objective:** Migrate `internal/api/*.go` files (3 files, ~80 strings)

#### Files
```
internal/api/server.go              ~40 strings
internal/api/workspace_manager.go   ~25 strings
internal/api/mcp_bridge.go          ~15 strings
```

#### Common Patterns

**HTTP Errors**
```go
// Before: "workspace ID required"
// After: i18n.T("api.validation.workspace_id_required")
```

**Server Status**
```go
// Before: "server starting on port 5082"
// After: i18n.T("api.server.starting", map[string]interface{}{"Port": port})
```

#### Tasks
- [ ] 4.1 Migrate server.go
  - [ ] Server lifecycle messages
  - [ ] HTTP error responses
  - [ ] Route registration messages
- [ ] 4.2 Migrate workspace_manager.go
  - [ ] Workspace validation
  - [ ] Registration messages
  - [ ] Context errors
- [ ] 4.3 Migrate mcp_bridge.go
  - [ ] Bridge errors
  - [ ] Protocol messages
- [ ] 4.4 API integration testing

#### Deliverables
```
✅ 3 migrated API files
✅ ~30 new API-specific keys
✅ HTTP responses work in both languages
✅ API tests passing
```

---

### Phase 5: Constants & Helpers ⏱️ 2-4 hours

**Objective:** Migrate `internal/constants/*.go` and helpers

#### Files
```
internal/constants/messages.go      ~35 strings (status/priority labels)
internal/constants/param_names.go   ~15 strings (descriptions)
internal/mcp/tool_helpers.go        ~25 strings (formatters)
internal/gorev/modeller.go          ~20 strings (emoji + label maps)
```

#### Migration Strategy

**Status/Priority Labels** (Currently in modeller.go)
```go
// BEFORE
var StatusLabels = map[string]string{
    "pending": "Beklemede",
    "in_progress": "Devam Ediyor",
    "completed": "Tamamlandı",
}

// AFTER
func GetStatusLabel(status string) string {
    return i18n.T(fmt.Sprintf("common.status.%s", status))
}
```

#### Tasks
- [ ] 5.1 Migrate status/priority labels to i18n
- [ ] 5.2 Update all references to use i18n.T()
- [ ] 5.3 Migrate tool helper formatters
- [ ] 5.4 Update constant descriptions

#### Deliverables
```
✅ Status/priority labels i18n-ready
✅ Helper functions updated
✅ Constants migrated
```

---

### Phase 6: Documentation Strings ⏱️ 2-3 hours

**Objective:** Migrate CLI descriptions and help texts

#### Files
```
cmd/gorev/main.go                   ~40 strings
cmd/gorev/mcp_commands.go           ~30 strings
internal/mcp/tool_registry.go       ~50 strings
```

#### Tasks
- [ ] 6.1 Migrate CLI command descriptions
- [ ] 6.2 Migrate flag descriptions
- [ ] 6.3 Migrate MCP tool descriptions
- [ ] 6.4 Update examples and help texts

---

### Phase 7: Testing & Validation ⏱️ 2-4 hours

**Objective:** Comprehensive validation and testing

#### Tasks
- [ ] 7.1 **Translation Coverage Test**
  - [ ] Write tool to detect hardcoded strings
  - [ ] Scan all .go files (exclude tests)
  - [ ] Report any remaining hardcoded strings
- [ ] 7.2 **Missing Key Detection**
  - [ ] Run app in TR mode, log all i18n.T() calls
  - [ ] Run app in EN mode, log all i18n.T() calls
  - [ ] Compare: any missing translations?
- [ ] 7.3 **Consistency Validation**
  - [ ] Verify all DRY keys exist in both TR and EN
  - [ ] Check for orphaned keys (defined but not used)
  - [ ] Validate key naming conventions
- [ ] 7.4 **Functional Testing**
  - [ ] Run all 70 integration tests in TR mode
  - [ ] Run all 70 integration tests in EN mode
  - [ ] Manual smoke testing of major workflows
- [ ] 7.5 **Performance Testing**
  - [ ] Benchmark i18n.T() call overhead
  - [ ] Ensure no significant performance regression
- [ ] 7.6 **Documentation**
  - [ ] Update CLAUDE.md with i18n guidelines
  - [ ] Create i18n development guide
  - [ ] Document DRY pattern usage

#### Deliverables
```
✅ Translation coverage: 100%
✅ Missing keys: 0
✅ Orphaned keys: Cleaned up
✅ All tests passing in both languages
✅ Performance: No regression
✅ Documentation: Complete
```

---

## 🧪 Testing Strategy

### Unit Tests
```bash
# Test i18n helpers
go test internal/i18n/helpers_test.go -v

# Test i18n manager
go test internal/i18n/manager_test.go -v
```

### Integration Tests
```bash
# Run all tests in Turkish
export GOREV_LANG=tr
go test ./... -v

# Run all tests in English
export GOREV_LANG=en
go test ./... -v
```

### Manual Testing Checklist
- [ ] Create task (TR and EN)
- [ ] Update task status (TR and EN)
- [ ] Delete task with confirmation (TR and EN)
- [ ] Batch operations (TR and EN)
- [ ] Export/import (TR and EN)
- [ ] Template usage (TR and EN)
- [ ] Error scenarios (TR and EN)

---

## 📐 DRY Pattern Library

### Pattern 1: Entity Not Found
```json
{
  "common": {
    "validation": {
      "not_found": "{{.Entity}} bulunamadı: {{.Error}}"
    }
  }
}
```

**Usage:**
```go
entity := i18n.T("common.entities.task")
return i18n.T("common.validation.not_found", map[string]interface{}{
    "Entity": entity,
    "Error": err.Error(),
})
```

### Pattern 2: Operation Failed
```json
{
  "common": {
    "operations": {
      "update_failed": "{{.Entity}} güncellenemedi: {{.Error}}"
    }
  }
}
```

**Usage:**
```go
entity := i18n.T("common.entities.project")
return i18n.T("common.operations.update_failed", map[string]interface{}{
    "Entity": entity,
    "Error": err.Error(),
})
```

### Pattern 3: Required Parameter
```json
{
  "common": {
    "validation": {
      "required": "{{.Param}} parametresi gerekli"
    }
  }
}
```

**Usage:**
```go
return i18n.T("common.validation.required", map[string]interface{}{
    "Param": "task_id",
})
```

### Pattern 4: Operation Success
```json
{
  "common": {
    "operations": {
      "create_success": "✓ {{.Entity}} oluşturuldu: {{.Title}} (ID: {{.Id}})"
    }
  }
}
```

**Usage:**
```go
entity := i18n.T("common.entities.task")
return i18n.T("common.operations.create_success", map[string]interface{}{
    "Entity": entity,
    "Title": task.Title,
    "Id": task.ID,
})
```

---

## 🚨 Risk Mitigation

### Risk 1: Breaking Existing Functionality
**Mitigation:**
- Migrate one file at a time
- Run tests after each file
- Keep backups of original files
- Use git branches for each phase

### Risk 2: Missing Translations
**Mitigation:**
- Automated key detection tool
- Fallback to key name if translation missing
- Pre-commit hook to check translation parity

### Risk 3: Performance Regression
**Mitigation:**
- Benchmark before/after
- Cache frequently used translations
- Profile i18n.T() overhead

### Risk 4: DRY Pattern Complexity
**Mitigation:**
- Clear documentation with examples
- Helper functions to simplify usage
- Code review to catch violations

---

## 📊 Progress Tracking

### Overall Progress
- [ ] Phase 1: DRY Key Hierarchy Design (0%)
- [ ] Phase 2: MCP Handlers (0%)
- [ ] Phase 3: Core Business Logic (0%)
- [ ] Phase 4: API Layer (0%)
- [ ] Phase 5: Constants & Helpers (0%)
- [ ] Phase 6: Documentation Strings (0%)
- [ ] Phase 7: Testing & Validation (0%)

### Metrics Dashboard
```
Translation Coverage: 0% → 100%
Hardcoded Strings: ~1200 → 0
Translation Keys: 695 → ~900 (with DRY reuse)
Key Reuse Rate: 0% → 60%+
Tests Passing: 70/70 → 70/70
```

---

## ✅ Definition of Done

- [ ] Zero hardcoded user-facing strings in non-test files
- [ ] 100% TR/EN translation parity
- [ ] All 70 integration tests passing in both languages
- [ ] DRY key reuse rate > 60%
- [ ] Performance: No regression > 5%
- [ ] Documentation: i18n guide complete
- [ ] Code review: Approved by maintainer
- [ ] Smoke testing: All major workflows tested in TR and EN

---

## 📚 References

- [go-i18n Documentation](https://github.com/nicksnyder/go-i18n)
- [Gorev i18n System](../gorev-mcpserver/internal/i18n/)
- [DRY Principles](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself)
- [i18n Best Practices](https://phrase.com/blog/posts/i18n-best-practices/)

---

**Last Updated:** 2025-10-11
**Next Review:** After Phase 1 completion
