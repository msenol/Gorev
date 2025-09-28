# Gorev VS Code Extension Test Suite

This directory contains comprehensive tests for the Gorev VS Code extension, including unit tests, integration tests, and end-to-end tests.

## Test Structure

```
test/
├── unit/                  # Unit tests for individual components
│   ├── markdownParser.test.js
│   ├── mcpClient.test.js
│   └── treeProviders.test.js
├── integration/           # Integration tests for extension features
│   └── extension.test.js
├── e2e/                  # End-to-end workflow tests
│   └── workflow.test.js
├── fixtures/             # Test data and mocks
│   └── mockData.js
├── utils/                # Test utilities
│   └── testHelper.js
├── suite/                # Test suite configuration
│   └── index.js
└── runTest.js           # Test runner entry point
```

## Running Tests

### Prerequisites

1. Install dependencies:

```bash
npm install
```

2. Compile TypeScript:

```bash
npm run compile
```

3. Ensure the Gorev MCP server is built:

```bash
cd ../gorev-mcpserver
make build
```

### Run All Tests

```bash
npm test
```

### Run Specific Test Suites

```bash
# Run only unit tests
npm test -- --grep "unit"

# Run only integration tests
npm test -- --grep "integration"

# Run only E2E tests
npm test -- --grep "E2E"
```

### Run Tests in Watch Mode

For development, you can run tests in watch mode:

```bash
npm run test-watch
```

## Test Categories

### Unit Tests

Unit tests focus on individual components in isolation:

- **markdownParser.test.js**: Tests for parsing MCP markdown responses
- **mcpClient.test.js**: Tests for MCP client functionality
- **treeProviders.test.js**: Tests for tree view providers

### Integration Tests

Integration tests verify that different components work together:

- **extension.test.js**: Tests extension activation, command registration, and basic workflows

### End-to-End Tests

E2E tests simulate real user workflows:

- **workflow.test.js**: Complete task management workflows including server connection

## Writing Tests

### Using Test Helpers

The `TestHelper` class provides utilities for common test scenarios:

```javascript
const TestHelper = require('../utils/testHelper');

suite('My Test Suite', () => {
  let helper;

  setup(() => {
    helper = new TestHelper();
    helper.setupCommonStubs();
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should do something', async () => {
    const mockClient = helper.createMockMCPClient();
    helper.setupMockMCPClient(mockClient);
    
    // Your test logic here
  });
});
```

### Mock Data

Use the mock data from `fixtures/mockData.js`:

```javascript
const { mockTasks, mockProjects, mockMCPResponses } = require('../fixtures/mockData');
```

### Testing Commands

```javascript
test('should execute command', async () => {
  // Setup inputs
  helper.setupTaskCreationInputs(helper.stubs, {
    title: 'Test Task',
    priority: 'yuksek'
  });

  // Execute command
  await vscode.commands.executeCommand('gorev.createTask');

  // Assert results
  helper.assertNotification('Information', 'başarıyla oluşturuldu');
});
```

## Debugging Tests

1. Set breakpoints in test files
2. Use VS Code's JavaScript Debug Terminal
3. Run tests with `--inspect` flag:

```bash
node --inspect ./node_modules/.bin/mocha
```

## Coverage

To generate coverage reports:

```bash
npm run test-coverage
```

Coverage reports will be generated in the `coverage/` directory.

## Troubleshooting

### Tests fail with "Extension not found"

Make sure the extension is compiled and the manifest is correct:

```bash
npm run compile
```

### E2E tests fail with "Server not found"

1. Build the Gorev server:

```bash
cd ../gorev-mcpserver
make build
```

2. Update the server path in tests or VS Code settings

### Timeout errors

Increase timeout for slow operations:

```javascript
test('slow test', async function() {
  this.timeout(10000); // 10 seconds
  // test logic
});
```

## Best Practices

1. **Isolate tests**: Each test should be independent
2. **Use mocks**: Mock external dependencies (VS Code API, MCP server)
3. **Clean up**: Always restore stubs and clean up resources
4. **Descriptive names**: Use clear test descriptions
5. **Test edge cases**: Include error scenarios and edge cases

## CI/CD Integration

The test suite is designed to run in CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run tests
  run: |
    npm install
    npm run compile
    npm test
```
