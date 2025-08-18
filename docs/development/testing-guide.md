# Testing Guide

This guide covers testing procedures for both the Gorev VS Code extension and MCP server.

## VS Code Extension Testing

### Unit Tests

Located in `gorev-vscode/test/unit/`:
- `markdownParser.test.js` - Tests for parsing MCP responses (207 LOC)
- `mcpClient.test.js` - Tests for MCP client functionality (100 LOC)
- `treeProviders.test.js` - Tests for tree view providers (149 LOC)
- `enhancedGorevTreeProvider.test.js` - Enhanced TreeView functionality tests (389 LOC)
- `taskDetailPanel.test.js` - WebView panel tests (396 LOC)
- `logger.test.js` - Logging utility tests (237 LOC)
- `models.test.js` - TypeScript model validation (273 LOC)
- `utils.test.js` - Utility function tests (307 LOC)

#### Running Unit Tests
```bash
cd gorev-vscode
npm test
# Or specific tests
npm test -- --grep "markdownParser"
```

### Integration Tests

Located in `gorev-vscode/test/integration/`:
- `extension.test.js` - Tests extension activation and command registration

#### Running Integration Tests
```bash
cd gorev-vscode
npm test -- --grep "integration"
```

### End-to-End Tests

Located in `gorev-vscode/test/e2e/`:
- `workflow.test.js` - Complete user workflows including server connection

#### Running E2E Tests
```bash
cd gorev-vscode
# Ensure MCP server is built first
cd ../gorev-mcpserver && make build && cd ../gorev-vscode
npm test -- --grep "E2E"
```

### Manual Testing

#### Advanced Filtering Toolbar Test

1. **Launch Extension:**
   ```bash
   cd gorev-vscode
   code .
   # Press F5 to launch Extension Development Host
   ```

2. **Status Bar Controls:**
   Check for new buttons in the status bar:
   - 🔍 Search
   - 🔧 Filter
   - 📑 Profiles

3. **Search Test:**
   - Click "🔍 Search" button
   - Search for terms like "bug", "frontend", or "urgent"
   - Verify tasks are filtered correctly

4. **Advanced Filter Menu:**
   - Click "🔧 Filter" button
   - Multi-select quick pick should open
   - Test these filters:
     - **Status**: Pending, In Progress, Completed
     - **Priority**: High, Medium, Low
     - **Special Filters**: Overdue, Due Today, Due This Week
     - **Project**: Dynamic project list

5. **Filter Profiles:**
   - Select multiple filters
   - Click Save button (💾)
   - Name the profile (e.g., "Urgent Tasks")
   - Load saved profile from "📑 Profiles" button

6. **Command Palette Tests (Ctrl+Shift+P):**
   - `Gorev: Search Tasks`
   - `Gorev: Show Filter Menu`
   - `Gorev Filter: Show Overdue Tasks`
   - `Gorev Filter: Show High Priority Tasks`
   - `Gorev Filter: Filter by Tag`

7. **Active Filter Indicator:**
   - When filters are applied, status bar shows "🔧 X filters active"
   - Click to clear all filters

### Test Coverage

Generate coverage reports:
```bash
cd gorev-vscode
npm run test-coverage  # or npm run coverage
# Report shows file-by-file coverage analysis
# Current coverage: 100% (full coverage achieved)
```

#### Custom Coverage Tool
The project includes a custom coverage analysis tool (`test-coverage.js`) that:
- Analyzes TypeScript source files and their test coverage
- Provides detailed LOC (Lines of Code) metrics
- Shows test/source ratio (optimized with DRY patterns)
- Identifies any regression in test coverage

## MCP Server Testing

### Unit Tests

Located in `gorev-mcpserver/internal/gorev/`:
- `veri_yonetici_test.go` - Data layer tests
- `is_yonetici_test.go` - Business logic tests
- `ai_context_yonetici_test.go` - AI context thread-safety tests

Located in `gorev-mcpserver/internal/mcp/`:
- `handlers_test.go` - MCP protocol handler tests (all 25 tools)
- `server_test.go` - MCP server initialization tests
- `test_helpers.go` - DRY test infrastructure (NEW)
- `table_driven_test.go` - Table-driven test patterns (NEW)
- `concurrency_test.go` - DRY concurrency testing (NEW)
- `benchmark_test.go` - Standardized benchmark suite (NEW)
- `dry_validation_test.go` - Focused validation tests (NEW)

### DRY Test Patterns

The testing infrastructure implements comprehensive DRY patterns to eliminate duplicate test code:

#### Table-Driven Test Structure
```go
type TestCase struct {
    Name        string
    Args        map[string]interface{}
    ExpectError bool
    ErrorMsg    string
    Setup       func(*testing.T, *MCPServer)
    Cleanup     func(*testing.T, *MCPServer)
    Validate    func(*testing.T, *mcp.CallToolResult)
}
```

#### Benchmark Configuration
```go
type BenchmarkConfig struct {
    Name        string
    ToolName    string
    Args        map[string]interface{}
    Setup       func(*testing.B, *MCPServer)
    Cleanup     func(*testing.B, *MCPServer)
    Iterations  int
    Parallel    bool
}
```

#### Concurrency Test Patterns
```go
type ConcurrencyTestConfig struct {
    Name         string
    ToolName     string
    Args         map[string]interface{}
    Goroutines   int
    Operations   int
    Setup        func(*testing.T, *MCPServer)
    Validate     func(*testing.T, []interface{})
    ExpectRaces  bool
}
```

#### Running Tests
```bash
cd gorev-mcpserver
# All tests with coverage
make test

# Specific package
go test ./internal/gorev/...
go test ./internal/mcp/...

# With race detection
go test -race ./...

# Generate coverage report
make test-coverage
# Current coverage: 81.3% (enhanced with DRY patterns)
```

### Integration Tests

Located in `gorev-mcpserver/test/`:
- `integration_test.go` - MCP handler tests

#### Running Integration Tests
```bash
cd gorev-mcpserver
go test ./test/...
```

### Manual MCP Testing

Test MCP tools directly:
```bash
# Start server in debug mode
./gorev serve --debug

# In another terminal, test with MCP protocol
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"gorev_listele","arguments":{}},"id":1}' | ./gorev serve
```

## DRY Testing Patterns (v0.11.1+)

### Constants Usage in Tests

All test files must use centralized constants from `internal/constants/test_constants.go`:

#### Template Constants
```go
// ✅ Correct - Use constants
result, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
    constants.ParamTemplateID: constants.TestTemplateFeatureRequest,
    constants.ParamDegerler: map[string]interface{}{
        "baslik": constants.TestTaskTitleEN,
    },
})

// ❌ Wrong - Hardcoded strings
result, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
    "template_id": "feature_request",
    "degerler": map[string]interface{}{
        "baslik": "Test Task",
    },
})
```

#### Test Iteration Constants
```go
// ✅ Correct - Use constants for pagination and loops
for i := 0; i < constants.TestIterationSmall; i++ {
    // test logic
}

params := map[string]interface{}{
    "limit": float64(constants.TestPaginationLimit),
}

// ❌ Wrong - Magic numbers
for i := 0; i < 10; i++ {
    // test logic
}

params := map[string]interface{}{
    "limit": float64(10),
}
```

#### Concurrency Test Constants
```go
// ✅ Correct - Use predefined concurrency levels
config := ConcurrencyTestConfig{
    Goroutines: constants.TestConcurrencyLarge, // 50
    Operations: constants.TestIterationMedium,  // 50
    Timeout:    constants.TestTimeoutLargeSeconds * time.Second,
}

// ❌ Wrong - Hardcoded concurrency values
config := ConcurrencyTestConfig{
    Goroutines: 50,
    Operations: 50,
    Timeout:    30 * time.Second,
}
```

### DRY Test Infrastructure

The testing infrastructure includes comprehensive reusable patterns:

#### Test Helper Functions
```go
// Use centralized test environment setup
env := SetupTestEnvironment(t)
defer env.Cleanup()

// Use helper functions for common operations
projectID := CreateTestProject(t, env, constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
taskID := CreateTestTask(t, env, constants.TestTemplateFeatureRequest, taskValues)
```

#### Table-Driven Test Patterns
```go
// Use standardized TestCase struct
testCases := []TestCase{
    {
        Name:       "ValidInput",
        Input:      constants.TestIDBasic,
        Expected:   expectedResult,
        ShouldFail: false,
    },
}

TableDrivenTest(t, "ValidationTest", testCases, testFunc)
```

### Rule 15 Compliance in Tests

All test code must follow Rule 15 principles:

- ✅ **Use constants** for all repeated values
- ✅ **No hardcoded strings** in test parameters
- ✅ **No magic numbers** in test configurations
- ✅ **DRY helper functions** for common test operations
- ❌ **No copy-paste** test patterns
- ❌ **No temporary** test values

## Test Data Management

### Creating Test Data

VS Code Extension:
```javascript
// Use test data seeder
vscode.commands.executeCommand('gorev.debug.seedTestData');
```

MCP Server:
```bash
# Use SQL directly
sqlite3 gorev.db
INSERT INTO projeler (id, isim, tanim) VALUES ('test-1', 'Test Project', 'Test Description');
INSERT INTO gorevler (id, baslik, proje_id, durum, oncelik) VALUES ('task-1', 'Test Task', 'test-1', 'beklemede', 'yuksek');
```

### Cleaning Test Data

```bash
# Remove database
rm -f gorev.db

# Or clean specific tables
sqlite3 gorev.db "DELETE FROM gorevler; DELETE FROM projeler;"
```

## Continuous Integration

### GitHub Actions Workflow

```yaml
name: Tests
on: [push, pull_request]

jobs:
  test-mcp-server:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - name: Run tests
        run: |
          cd gorev-mcpserver
          make test
          
  test-vscode-extension:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install and test
        run: |
          cd gorev-vscode
          npm install
          npm test
```

## Race Condition Testing (v0.11.1+)

### Go Race Detector

The Go race detector is essential for detecting race conditions in concurrent code. All tests must pass with the race detector enabled.

#### Running Tests with Race Detector
```bash
# Run all tests with race detector
cd gorev-mcpserver
go test -race ./...

# Run specific concurrent tests
go test -race -v -run TestAIContextRaceCondition ./internal/gorev/

# Build and run server with race detector
go build -race ./cmd/gorev
./gorev serve --debug
```

#### Continuous Integration
The CI pipeline includes race detection:
```bash
# In CI scripts
make test-race  # Equivalent to go test -race ./...
```

### Concurrent Testing Patterns

#### Standard Race Condition Test Pattern
```go
func TestComponentName_ConcurrentAccess(t *testing.T) {
    // Setup component under test
    component := setupTestComponent()
    
    // Error collection
    errors := make(chan error, 100)
    const numGoroutines = 50
    const operationsPerGoroutine = 10
    
    // Launch concurrent operations
    for i := 0; i < numGoroutines; i++ {
        go func(id int) {
            defer func() {
                if r := recover(); r != nil {
                    errors <- fmt.Errorf("goroutine %d panicked: %v", id, r)
                }
            }()
            
            for j := 0; j < operationsPerGoroutine; j++ {
                if err := component.ConcurrentOperation(); err != nil {
                    errors <- fmt.Errorf("goroutine %d operation failed: %w", id, err)
                    return
                }
            }
        }(i)
    }
    
    // Wait and collect results
    time.Sleep(100 * time.Millisecond)
    close(errors)
    
    var collectedErrors []error
    for err := range errors {
        collectedErrors = append(collectedErrors, err)
    }
    
    assert.Empty(t, collectedErrors, "No race conditions should occur")
}
```

#### AI Context Manager Example
The `TestAIContextRaceCondition` in `ai_context_yonetici_test.go` demonstrates comprehensive concurrent testing:

```go
func TestAIContextRaceCondition(t *testing.T) {
    // Test demonstrates 50 goroutines performing 500 total operations
    // Mix of read and write operations:
    // - SetActiveTask (write)
    // - GetActiveTask (read)
    // - GetContext (read)
    // - GetRecentTasks (read)
    
    // Verifies no race conditions occur under high concurrency
}
```

### Load Testing for Concurrency

#### MCP Tool Concurrent Access
```bash
#!/bin/bash
# test_concurrent_mcp.sh

echo "Testing concurrent MCP tool access..."

call_mcp_tool() {
    local task_id="task-$1"
    ./gorev mcp call gorev_set_active "{\"task_id\": \"$task_id\"}" 2>&1
}

# Launch 20 concurrent MCP calls
for i in {1..20}; do
    call_mcp_tool $i &
done

wait
echo "Concurrent MCP test completed"
```

### Benchmarking Concurrent Performance

#### Concurrent Benchmark Pattern
```go
func BenchmarkConcurrentAccess(b *testing.B) {
    component := setupBenchmarkComponent()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Mix of operations to benchmark
            switch rand.Intn(3) {
            case 0:
                component.WriteOperation()
            case 1:
                component.ReadOperation()
            case 2:
                component.QueryOperation()
            }
        }
    })
}
```

### Thread-Safety Verification Checklist

When testing concurrent components:

- [ ] **Race Detector**: All tests pass with `go test -race`
- [ ] **Stress Testing**: 50+ goroutines performing mixed operations
- [ ] **Error Collection**: Proper error handling from concurrent operations
- [ ] **Final State Verification**: Consistent state after concurrent access
- [ ] **Performance**: No significant degradation under concurrency
- [ ] **Deadlock Detection**: Tests complete without hanging

### Integration with CI/CD

#### GitHub Actions Race Detection
```yaml
name: Race Condition Tests
on: [push, pull_request]

jobs:
  race-detection:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - name: Run race detector
        run: |
          cd gorev-mcpserver
          go test -race -timeout=10m ./...
```

## Best Practices

1. **Test Isolation**
   - Each test should be independent
   - Use setup/teardown properly
   - Don't rely on test execution order

2. **Mock External Dependencies**
   - Mock VS Code API in unit tests
   - Mock MCP server responses
   - Use test doubles for database

3. **Test Naming**
   - Use descriptive test names
   - Follow pattern: `should_expectedBehavior_when_condition`

4. **Test Coverage**
   - Aim for >80% coverage
   - Focus on critical paths

5. **Concurrent Testing (v0.11.1+)**
   - Always test concurrent components with race detector
   - Use high goroutine counts (50+) for stress testing
   - Verify final state consistency after concurrent operations
   - Include concurrent scenarios in integration tests
   - Don't just chase numbers

5. **Performance Testing**
   - Test with large datasets (1000+ tasks)
   - Monitor memory usage
   - Check response times

## Troubleshooting

### Common Issues

1. **"Extension not found" in tests**
   ```bash
   npm run compile
   ```

2. **"Server not found" in E2E tests**
   ```bash
   cd ../gorev-mcpserver && make build
   ```

3. **Timeout errors**
   ```javascript
   this.timeout(10000); // Increase timeout
   ```

4. **Flaky tests**
   - Add proper waits
   - Check for race conditions
   - Use stable selectors

### Debug Tips

1. **VS Code Extension**
   - Set breakpoints in test files
   - Use `console.log` for debugging
   - Check test output panel

2. **MCP Server**
   - Use `t.Logf()` for debug output
   - Run with `-v` flag: `go test -v`
   - Use debugger: `dlv test`

## Testing Refactored Code (v0.11.1+)

### Architecture Verification Testing

After the major refactoring in v0.11.1, ensure that architectural improvements maintain functionality:

#### 1. Tool Registration Testing
Verify that the new `ToolRegistry` pattern correctly registers all 25 MCP tools:

```bash
# Test tool registration
./gorev mcp list | wc -l  # Should show 25 tools

# Test specific tool categories
./gorev mcp list | grep -E "(gorev_|proje_|template_|aktif_|ozet_)"
```

#### 2. Helper Class Unit Testing
Test the extracted helper classes individually:

```go
func TestParameterValidator(t *testing.T) {
    validator := NewParameterValidator()
    
    // Test required string validation
    result, err := validator.ValidateRequiredString(params, "test_param")
    assert.NoError(t, err)
    assert.Equal(t, "expected_value", result)
}

func TestTaskFormatter(t *testing.T) {
    formatter := NewTaskFormatter()
    
    // Test status emoji formatting
    emoji := formatter.GetStatusEmoji("tamamlandi")
    assert.Equal(t, "✅", emoji)
}
```

#### 3. Integration Testing Pattern
Verify that refactored components work together correctly:

```bash
# Test complete workflow after refactoring
make test  # Run all tests
go test -v ./internal/mcp/  # Specific MCP package tests
```

#### 4. Regression Testing Checklist
After any refactoring, verify:

- [ ] **All MCP tools register correctly** (`./gorev mcp list`)
- [ ] **Build succeeds without errors** (`make build`)
- [ ] **All tests pass** (`make test`)
- [ ] **No performance degradation** (benchmark if needed)
- [ ] **API compatibility maintained** (existing clients work)

#### 5. Code Quality Verification

```bash
# Verify code formatting
make fmt

# Run static analysis
go vet ./...

# Check for code smells (manual review)
wc -l internal/mcp/*.go  # Verify reasonable file sizes
```

### Refactoring Testing Best Practices

1. **Test Before Refactoring**: Ensure comprehensive test coverage exists
2. **Incremental Testing**: Test after each refactoring step
3. **Functional Verification**: Verify behavior unchanged after refactoring
4. **Performance Testing**: Ensure refactoring doesn't degrade performance
5. **Integration Testing**: Test interactions between refactored components

### Testing the Tool Registry Pattern

The new tool registry pattern should be tested to ensure:

- All tool categories are registered
- Tools have correct schemas and handlers
- Registration order doesn't affect functionality
- Error handling works properly

```go
func TestToolRegistryCategories(t *testing.T) {
    handler := setupTestHandler()
    registry := NewToolRegistry(handler)
    server := setupTestMCPServer()
    
    registry.RegisterAllTools(server)
    
    // Verify all categories registered
    tools := server.GetRegisteredTools()
    assert.Len(t, tools, 25)  // Expected tool count
}
```

## Additional Resources

- [VS Code Testing Guide](https://code.visualstudio.com/api/working-with-extensions/testing-extension)
- [Go Testing Guide](https://go.dev/doc/tutorial/add-a-test)
- [Mocha Documentation](https://mochajs.org/)
- [Sinon.js Documentation](https://sinonjs.org/)
- [Architecture Guide](architecture.md) - Refactoring patterns and improvements
- [Concurrency Guide](concurrency-guide.md) - Thread-safety testing patterns