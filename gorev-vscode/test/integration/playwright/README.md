# Playwright UI Testing Suite

This directory contains comprehensive UI testing for the Gorev VS Code Extension using Playwright. These tests verify actual UI functionality, not just extension loading.

## Overview

The test suite includes:

1. **API Integration Tests** - Verify REST API endpoints work correctly
2. **Task Workflow Tests** - Test complete user workflows (create, edit, delete tasks)
3. **UI Component Tests** - Test tree view interactions, context menus, filters
4. **Mock Server** - Lightweight API server for testing without requiring the actual Gorev daemon

## Prerequisites

```bash
# Install dependencies
npm install

# Install Playwright browsers (one-time setup)
npx playwright install chromium
```

## Test Structure

```
test/integration/playwright/
├── README.md                 # This file
├── mock-server.ts            # Mock API server for testing
├── api-integration.spec.ts   # API endpoint tests
├── task-workflow.spec.ts     # UI workflow tests
└── playwright.config.ts      # Playwright configuration
```

## Running Tests

### All Tests
```bash
npm run test:ui
```

### Specific Test Suites
```bash
# API integration tests only
npm run test:api

# UI workflow tests only
npm run test:workflow
```

### Debug Modes
```bash
# Run with browser visible (headed mode)
npm run test:ui:headed

# Debug mode with Playwright Inspector
npm run test:ui:debug
```

### Individual Test Files
```bash
npx playwright test api-integration.spec.ts
npx playwright test task-workflow.spec.ts
```

## Test Categories

### 1. API Integration Tests

These tests verify the REST API endpoints work correctly:
- ✅ Health check endpoint
- ✅ Task CRUD operations (Create, Read, Update, Delete)
- ✅ Subtask operations
- ✅ Project operations
- ✅ Template operations
- ✅ Summary statistics
- ✅ Pagination
- ✅ Error handling (404, 400)
- ✅ CORS headers

**Run**: `npm run test:api`

### 2. Task Workflow Tests

These tests simulate real user interactions:
- ✅ Load tasks from API and display in tree view
- ✅ Filter tasks by project
- ✅ Create new task from template
- ✅ Edit task via context menu
- ✅ Update task status
- ✅ Display subtasks hierarchy
- ✅ Delete task
- ✅ Search tasks by title
- ✅ Refresh task list
- ✅ Task priority display
- ✅ Task statistics in sidebar

**Run**: `npm run test:workflow`

## Mock Server

The `mock-server.ts` provides a lightweight Express.js server that mimics the Gorev API:

- **Port**: 5083 (configurable)
- **Endpoints**: All major API endpoints
- **Test Data**: Pre-populated with sample tasks and projects
- **In-Memory Storage**: Changes are lost on server restart
- **CORS Enabled**: For cross-origin requests

### Mock Data

The server includes:
- 1 test project
- 3 tasks (1 completed, 1 in progress, 1 pending)
- 1 subtask
- 2 templates (Bug Report, Feature Request)

## Writing New Tests

### Example: Basic Test

```typescript
import { test, expect } from '@playwright/test';

test('should display task title', async ({ page }) => {
  await page.goto('http://localhost:5001');
  await expect(page.locator('text=Task Title')).toBeVisible();
});
```

### Example: API Test

```typescript
import { test, expect } from '@playwright/test';

test('should create new task', async ({ page }) => {
  const response = await page.request.post('http://localhost:5083/api/v1/tasks/from-template', {
    data: {
      template_id: 'bug-report',
      values: {
        title: 'Test Bug',
        description: 'Test description'
      }
    }
  });

  expect(response.status()).toBe(200);
  const data = await response.json();
  expect(data.success).toBe(true);
});
```

### Example: User Interaction Test

```typescript
import { test, expect } from '@playwright/test';

test('should edit task', async ({ page }) => {
  await page.goto('http://localhost:5001');

  // Right-click on task
  const taskItem = page.locator('[data-testid="task-item"]').first();
  await taskItem.click({ button: 'right' });

  // Click edit
  await page.locator('[data-testid="menu-item-edit"]').click();

  // Update title
  await page.locator('[data-testid="input-title"]').clear();
  await page.locator('[data-testid="input-title"]').fill('Updated Title');

  // Save
  await page.locator('[data-testid="save-button"]').click();

  // Verify
  await expect(page.locator('text=Updated Title')).toBeVisible();
});
```

## Test Data Attributes

To make tests more resilient, use `data-testid` attributes in the Web UI:

```html
<!-- Good for testing -->
<button data-testid="create-task-button">Create Task</button>
<div data-testid="task-item">
  <span data-testid="task-title">Task Title</span>
  <span data-testid="status-badge">pending</span>
</div>

<!-- Avoid -->
<button>Create Task</button>
<div>
  <span>Task Title</span>
</div>
```

## Configuration

### playwright.config.ts

Key settings:
- **baseURL**: http://localhost:5001 (Web UI dev server)
- **retries**: 2 (on CI)
- **workers**: Parallel test execution
- **reporters**: HTML, JSON, JUnit
- **screenshots**: On failure
- **Videos**: On failure

### Environment Variables

```bash
# CI mode (disable parallel tests, more retries)
CI=true npm run test:ui

# Custom port for mock server
MOCK_PORT=5083 npm run test:api
```

## Continuous Integration

Tests are designed to run in CI:
- Headless browser mode
- Automatic retries on failure
- Screenshot/video on failure
- JUnit XML output for CI integration
- HTML report for detailed analysis

Example GitHub Actions step:

```yaml
- name: Run Playwright tests
  run: npm run test:ui

- name: Upload Playwright Report
  uses: actions/upload-artifact@v3
  if: always()
  with:
    name: playwright-report
    path: playwright-report/
    retention-days: 30
```

## Debugging

### Using Playwright Inspector
```bash
npm run test:ui:debug
```

### Using Browser Developer Tools
```bash
npm run test:ui:headed
```
Then right-click and "Inspect" in the Chromium window.

### Check Test Output
```bash
# View HTML report
open playwright-report/index.html

# View test results JSON
cat test-results/results.json
```

### Screenshot on Failure
Screenshots are automatically saved to `test-results/` when tests fail.

### Video on Failure
Videos are automatically recorded and saved to `test-results/` when tests fail.

## Best Practices

1. **Use data-testid attributes** for reliable element selection
2. **Test user workflows**, not just individual components
3. **Mock external dependencies** (API, database)
4. **Clean up test data** in `afterEach` or `afterAll` hooks
5. **Use realistic test data** that mimics production
6. **Test both happy path and error scenarios**
7. **Make tests independent** - they should not depend on each other
8. **Use proper assertions** - `expect()` with meaningful messages

## Common Patterns

### Waiting for Elements
```typescript
// Wait for element to be visible
await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

// Wait for element to have text
await expect(page.locator('[data-testid="task-count"]')).toContainText('3');
```

### Handling Dialogs
```typescript
// Handle browser dialog (confirm, alert, prompt)
page.on('dialog', dialog => dialog.accept());

// Or handle with specific text
page.on('dialog', dialog => {
  expect(dialog.message()).toContainText('Delete task?');
  dialog.accept();
});
```

### Mocking API Responses
```typescript
// Intercept and mock API calls
await page.route('**/api/v1/tasks', route => {
  route.fulfill({
    json: [{ id: '1', title: 'Mock Task', status: 'pending' }]
  });
});
```

### Testing with MockServer
```typescript
import MockServer from './mock-server';

let mockServer: MockServer;

test.beforeAll(async () => {
  mockServer = new MockServer(5083);
  await mockServer.start();
});

test.afterAll(async () => {
  await mockServer.stop();
});
```

## Limitations

1. **VS Code Extension**: These tests focus on the Web UI. For VS Code extension UI testing, use `@vscode/test-electron` (already in the project).
2. **Browser Support**: Tests run on Chromium, Firefox, and WebKit. Some features may need browser-specific handling.
3. **Real-time Updates**: Mock server doesn't support WebSocket connections. Test real-time features with actual Gorev server.

## Troubleshooting

### Port Already in Use
```bash
# Kill process on port 5001
lsof -ti:5001 | xargs kill

# Or use different port
PORT=5002 npm run test:serve
```

### Playwright Browser Not Installed
```bash
npx playwright install chromium
```

### Tests Timing Out
```typescript
// Increase timeout for slow operations
test('slow test', async ({ page }) => {
  test.setTimeout(60000);
  // ...
});
```

### Element Not Found
```typescript
// Wait for element before interacting
await page.waitForSelector('[data-testid="task-item"]');
await page.click('[data-testid="task-item"]');
```

## References

- [Playwright Documentation](https://playwright.dev/)
- [Playwright API Reference](https://playwright.dev/docs/api/class-playwright)
- [VS Code Testing API](https://code.visualstudio.com/api/working-with-extensions/testing-extension)
- [Gorev Project README](../../../../README.md)

## Contributing

When adding new tests:
1. Follow the naming convention: `{feature}.spec.ts`
2. Use descriptive test names that explain what is being tested
3. Group related tests with `test.describe()`
4. Include both positive and negative test cases
5. Document complex test scenarios in comments
6. Update this README if adding new test categories

## Test Coverage Goals

- **API Endpoints**: 100% coverage
- **Task Workflows**: 90%+ coverage
- **User Interactions**: 85%+ coverage
- **Error Scenarios**: 80%+ coverage

Current status:
- ✅ API Integration: 20 tests
- ✅ Task Workflows: 12 tests
- 🔄 Total Coverage: 32 tests (growing)
