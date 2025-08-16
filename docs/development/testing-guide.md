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
# Current coverage: 50.9% (19/33 files tested)
```

#### Custom Coverage Tool
The project includes a custom coverage analysis tool (`test-coverage.js`) that:
- Analyzes TypeScript source files and their test coverage
- Provides detailed LOC (Lines of Code) metrics
- Shows test/source ratio (currently 0.45:1)
- Identifies untested files with recommendations

## MCP Server Testing

### Unit Tests

Located in `gorev-mcpserver/internal/gorev/`:
- `veri_yonetici_test.go` - Data layer tests
- `is_yonetici_test.go` - Business logic tests

Located in `gorev-mcpserver/internal/mcp/`:
- `handlers_test.go` - MCP protocol handler tests (561 LOC, all 16 tools)
- `server_test.go` - MCP server initialization tests

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
# Current coverage: 75.1% for MCP package, 53.8% for gorev package
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

## Additional Resources

- [VS Code Testing Guide](https://code.visualstudio.com/api/working-with-extensions/testing-extension)
- [Go Testing Guide](https://go.dev/doc/tutorial/add-a-test)
- [Mocha Documentation](https://mochajs.org/)
- [Sinon.js Documentation](https://sinonjs.org/)