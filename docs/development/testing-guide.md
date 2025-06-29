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