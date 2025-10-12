# Phase 2: i18n Refactoring Guide (v0.17.0)

## Status: READY TO EXECUTE
**Phase 1 Complete:** Core infrastructure stable, all tests passing

---

## Overview

This guide provides step-by-step instructions for refactoring 227 i18n function calls to support per-request language selection.

### Completed (Phase 1)
✅ `internal/context/language.go` - Context helpers
✅ `internal/i18n/manager.go` - TWithLang() function
✅ `internal/i18n/manager_test.go` - Bilingual tests (all passing)
✅ `internal/mcp/handlers.go` - extractLanguage() method added
✅ Backward compatible - existing code still works

### Remaining (Phase 2)
🔄 227 i18n function calls in MCP handlers
🔄 20+ i18n helper functions in helpers.go

---

## Step 1: Update i18n Helper Functions

**File:** `internal/i18n/helpers.go`

**Pattern:** Add `lang string` as first parameter to all helper functions

### Example Transformation:

```go
// BEFORE:
func TFetchFailed(entity string, err error) string {
    return T("common.operations.fetch_failed", map[string]interface{}{
        "Entity": entity,
        "Error":  err.Error(),
    })
}

// AFTER:
func TFetchFailed(lang string, entity string, err error) string {
    return TWithLang(lang, "common.operations.fetch_failed", map[string]interface{}{
        "Entity": entity,
        "Error":  err.Error(),
    })
}
```

### Functions to Update (20+):

1. `TCommon(lang string, key string, data map[string]interface{}) string`
2. `TParam(lang string, paramName string) string`
3. `TValidation(lang string, validationType string, param string, extra map[string]interface{}) string`
4. `TRequiredParam(lang string, param string) string`
5. `TRequiredArray(lang string, param string) string`
6. `TRequiredObject(lang string, param string) string`
7. `TEntityNotFound(lang string, entity string, err error) string`
8. `TEntityNotFoundByID(lang string, entity, id string) string`
9. `TOperationFailed(lang string, operation, entity string, err error) string`
10. `TSuccess(lang string, operation, entity string, details map[string]interface{}) string`
11. `TInvalidValue(lang string, param, value string, validValues []string) string`
12. `TInvalidStatus(lang string, status string, validStatuses []string) string`
13. `TInvalidPriority(lang string, priority string) string`
14. `TInvalidDate(lang string, dateValue string) string`
15. `TInvalidFormat(lang string, formatType, value string) string`
16. `TCreateFailed(lang string, entity string, err error) string`
17. `TUpdateFailed(lang string, entity string, err error) string`
18. `TDeleteFailed(lang string, entity string, err error) string`
19. `TFetchFailed(lang string, entity string, err error) string`
20. `TSaveFailed(lang string, entity string, err error) string`
21. `TLoadFailed(lang string, entity string, err error) string`
22. `TSearchFailed(lang string, entity string, err error) string`
23. `TMarkdownLabel(lang string, key string, value interface{}) string`
24. `TListItem(lang string, key string, value interface{}) string`
25. `TStatus(lang string, status string) string`
26. `TPriority(lang string, priority string) string`
27. `TAddFailed(lang string, entity string, err error) string`
28. `TRemoveFailed(lang string, entity string, err error) string`

**Automated Approach:**

```bash
# Backup first
cp internal/i18n/helpers.go internal/i18n/helpers.go.backup

# Use regex replacement (sed/awk/perl) or manual editing
# Pattern 1: Add lang parameter to function signatures
# Pattern 2: Replace i18n.T( with i18n.TWithLang(lang,
```

**Verification:**
```bash
go build -o /dev/null ./internal/i18n/...
```

---

## Step 2: Update MCP Handler Functions

**Files:**
- `internal/mcp/handlers.go` (157 calls)
- `internal/mcp/tool_helpers.go` (50 calls)
- `internal/mcp/tool_registry.go` (20 calls)

### Pattern:

1. Add `lang := h.extractLanguage()` at the start of each handler
2. Update all i18n helper calls to pass `lang` as first parameter

### Example Transformation:

```go
// BEFORE:
func (h *Handlers) OzetGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
    ozet, err := h.isYonetici.OzetAl()
    if err != nil {
        return mcp.NewToolResultError(i18n.TFetchFailed("summary", err)), nil
    }

    metin := i18n.T("headers.summaryReport")
    return mcp.NewToolResultText(metin), nil
}

// AFTER:
func (h *Handlers) OzetGoster(params map[string]interface{}) (*mcp.CallToolResult, error) {
    lang := h.extractLanguage()

    ozet, err := h.isYonetici.OzetAl()
    if err != nil {
        return mcp.NewToolResultError(i18n.TFetchFailed(lang, "summary", err)), nil
    }

    metin := i18n.TWithLang(lang, "headers.summaryReport")
    return mcp.NewToolResultText(metin), nil
}
```

### Systematic Approach:

**Step 2.1:** Find all handler functions
```bash
grep -n "func (h \*Handlers)" internal/mcp/handlers.go | wc -l
# Output: ~30 handlers
```

**Step 2.2:** For each handler, add language extraction
```go
lang := h.extractLanguage()
```

**Step 2.3:** Update all i18n calls in that handler
- `i18n.T(` → `i18n.TWithLang(lang,`
- `i18n.TFoo(` → `i18n.TFoo(lang,` (for all helpers)

**Step 2.4:** Verify after each file
```bash
go build -o /dev/null ./internal/mcp/...
go test ./internal/mcp/... -run TestHandlerName
```

---

## Step 3: Update tool_helpers.go

**File:** `internal/mcp/tool_helpers.go`

Similar pattern to handlers, but methods on `ToolHelpers` struct:

```go
// BEFORE:
func (th *ToolHelpers) FormatError(err error) string {
    return i18n.T("error.generic", map[string]interface{}{"Error": err.Error()})
}

// AFTER:
func (th *ToolHelpers) FormatError(lang string, err error) string {
    return i18n.TWithLang(lang, "error.generic", map[string]interface{}{"Error": err.Error()})
}
```

**Note:** ToolHelpers methods will need `lang` parameter added to signature.

---

## Step 4: Update tool_registry.go

**File:** `internal/mcp/tool_registry.go`

Tool descriptions and parameter descriptions use i18n:

```go
// BEFORE:
Description: i18n.T("tools.descriptions.ozet_goster"),

// AFTER:
// Tool registry uses default language (registration happens once at startup)
// Use i18n.TWithLang("tr", ...) for consistency
Description: i18n.TWithLang("tr", "tools.descriptions.ozet_goster"),
```

**Rationale:** Tool registry is static, registered once at startup. Using default "tr" is acceptable here, as tool metadata doesn't need per-request translation.

---

## Step 5: Update Tests

**Files:** All `*_test.go` files in `internal/mcp/`

### Test Pattern:

```go
// BEFORE:
func TestOzetGoster(t *testing.T) {
    result, err := handlers.OzetGoster(params)
    assert.Contains(t, result.Text, "özet")
}

// AFTER:
func TestOzetGoster(t *testing.T) {
    // Test with Turkish
    os.Setenv("GOREV_LANG", "tr")
    defer os.Unsetenv("GOREV_LANG")

    result, err := handlers.OzetGoster(params)
    assert.Contains(t, result.Text, "özet")
}

func TestOzetGosterEnglish(t *testing.T) {
    // Test with English
    os.Setenv("GOREV_LANG", "en")
    defer os.Unsetenv("GOREV_LANG")

    result, err := handlers.OzetGoster(params)
    assert.Contains(t, result.Text, "summary")
}
```

**Add bilingual test cases for key handlers:**
- OzetGoster (Turkish + English)
- GorevListele (Turkish + English)
- TemplateListele (Turkish + English)

---

## Step 6: Compilation & Testing Strategy

### Incremental Verification:

**After each file modification:**

1. **Compile Check:**
```bash
go build -o /dev/null ./internal/i18n/...
go build -o /dev/null ./internal/mcp/...
```

2. **Run Affected Tests:**
```bash
go test -v ./internal/i18n/... -run TestTWithLang
go test -v ./internal/mcp/... -run TestHandlerName
```

3. **Full Test Suite:**
```bash
go test ./internal/... -cover
```

4. **Integration Test:**
```bash
# Start server with GOREV_LANG=en
GOREV_LANG=en ./gorev serve --debug

# In another terminal, test MCP calls
echo '{"method": "tools/call", "params": {"name": "ozet_goster"}}' | npx @modelcontextprotocol/inspector
```

---

## Step 7: Rollback Plan

If issues arise during refactoring:

**Git Safety:**
```bash
# Create branch before starting
git checkout -b feature/phase2-i18n-refactoring

# Commit Phase 1 first (checkpoint)
git add internal/context/ internal/i18n/
git commit -m "feat(i18n): Phase 1 - per-request language infrastructure"

# Create checkpoints during Phase 2
git commit -m "wip: helpers.go refactored (20 functions)"
git commit -m "wip: handlers.go partial (50/157 calls)"
```

**Rollback if needed:**
```bash
git reset --hard HEAD~1  # Undo last commit
git checkout main        # Abandon branch
```

---

## Estimated Effort

| Task | Files | Changes | Time |
|------|-------|---------|------|
| Step 1: helpers.go | 1 file | 20 functions | 1 hour |
| Step 2: handlers.go | 1 file | 157 calls | 3-4 hours |
| Step 3: tool_helpers.go | 1 file | 50 calls | 1 hour |
| Step 4: tool_registry.go | 1 file | 20 calls | 30 min |
| Step 5: Tests | 10+ files | 50+ tests | 2 hours |
| Step 6: Integration Testing | - | - | 1 hour |
| **TOTAL** | **15+ files** | **~300 changes** | **8-10 hours** |

---

## Success Criteria

**Phase 2 Complete When:**

✅ All 227 i18n calls use `TWithLang(lang, ...)` or helpers with `lang` parameter
✅ All code compiles without errors
✅ All existing tests pass
✅ New bilingual tests added and passing
✅ Manual testing: `GOREV_LANG=en` produces English responses
✅ Manual testing: `GOREV_LANG=tr` produces Turkish responses
✅ Backward compatibility maintained (no breaking changes)

---

## Next Steps After Phase 2

**Phase 3:** Business Logic (internal/gorev/)
**Phase 4:** Multilingual Templates
**Phase 5:** VS Code Extension + Workspace Language
**Phase 6:** Comprehensive Testing

---

## Quick Reference Commands

```bash
# Count remaining i18n.T() calls
grep -r "i18n\.T(" internal/mcp/*.go | grep -v "_test.go" | wc -l

# Find handlers to update
grep "func (h \*Handlers)" internal/mcp/handlers.go | cut -d'(' -f1

# Test single handler
go test -v ./internal/mcp/... -run TestOzetGoster

# Full build check
make build

# Run with English
GOREV_LANG=en ./gorev serve --debug
```

---

**Document Version:** 1.0
**Created:** 2025-10-12
**Status:** Phase 1 Complete, Phase 2 Ready to Execute
**Estimated Completion:** Phase 2 in 8-10 hours of focused work
