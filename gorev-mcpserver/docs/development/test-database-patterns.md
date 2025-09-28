# Test Database Patterns - Standardization Guide

This document provides standardized patterns for test database initialization to reduce code duplication and improve maintainability.

## Problem Statement

Currently, the codebase has 40+ instances of duplicate database initialization patterns across test files. Each test manually creates database connections with similar but slightly different configurations.

## Solution: Standardized Test Helpers

**New Helper Location**: `internal/testing/helpers.go`

### Available Helpers

#### 1. `DefaultTestDatabaseConfig()`

Returns default configuration for most test scenarios:

- Memory database (`:memory:`)
- Default migrations path
- Creates default templates
- Initializes i18n

#### 2. `SetupTestDatabase(t *testing.T, config *TestDatabaseConfig)`

Creates configured test database with proper cleanup.

**Returns**: `*gorev.VeriYonetici, cleanup func()`

#### 3. `SetupTestEnvironmentBasic(t *testing.T)`

Creates complete test environment with business logic layer.

**Returns**: `*gorev.IsYonetici, cleanup func()`

### Migration Path

**Current Pattern** (40+ occurrences):

```go
veriYonetici, err := gorev.YeniVeriYonetici(constants.TestDatabaseURI, constants.TestMigrationsPath)
require.NoError(t, err)
defer veriYonetici.Kapat()

err = veriYonetici.VarsayilanTemplateleriOlustur()
require.NoError(t, err)
```

**New Standardized Pattern**:

```go
import "github.com/msenol/gorev/internal/testing"

func TestExample(t *testing.T) {
    isYonetici, cleanup := testing.SetupTestEnvironmentBasic(t)
    defer cleanup()
    
    // Test logic here
}
```

### Custom Configuration Example

```go
func TestCustomDatabase(t *testing.T) {
    config := &testing.TestDatabaseConfig{
        UseTempFile:     true,  // Use file database
        MigrationsPath:  customPath,
        CreateTemplates: false, // Skip template creation
        InitializeI18n:  true,
    }
    
    veriYonetici, cleanup := testing.SetupTestDatabase(t, config)
    defer cleanup()
    
    // Test logic here
}
```

## Benefits

1. **DRY Compliance**: Eliminates 40+ duplicate patterns
2. **Consistency**: Standardized database setup across all tests
3. **Maintainability**: Central location for test configuration changes
4. **Flexibility**: Configurable for different test scenarios
5. **Rule 15 Compliance**: No shortcuts or workarounds

## Implementation Status

- ✅ Helper functions created (`internal/testing/helpers.go`)
- ✅ Documentation completed
- ⏳ **Migration**: 40+ test files can be updated to use helpers
- ⏳ **Adoption**: Teams can migrate existing tests gradually

## Recommended Next Steps

1. **Gradual Migration**: Update test files one at a time to use helpers
2. **New Tests**: All new tests should use standardized helpers
3. **Documentation**: Update test writing guidelines to reference these helpers
4. **Code Review**: Ensure new tests use standardized patterns

## Impact Assessment

**Files Affected**: 40+ test files across the codebase
**Code Reduction**: ~300 lines of duplicate code can be eliminated
**Maintenance**: Single point of change for test database configuration
**Risk**: Low - helpers are additive, existing tests continue to work

This standardization significantly improves codebase maintainability while adhering to Rule 15 principles of comprehensive solutions without shortcuts.
