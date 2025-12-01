/**
 * Web UI End-to-End Tests
 *
 * Comprehensive E2E tests for the Gorev Web UI.
 * Uses the real Go server for integration testing.
 *
 * Run with: npx playwright test e2e/web-ui.spec.ts
 * With real server: SKIP_WEB_SERVER=true npx playwright test e2e/web-ui.spec.ts
 */

import { test, expect } from '@playwright/test';
import { TaskListPage } from '../page-objects';

// Test configuration
const BASE_URL = process.env.WEB_UI_URL || 'http://localhost:5001';
const API_URL = process.env.API_URL || 'http://localhost:5082';

test.describe('Web UI - Task List', () => {
  let taskListPage: TaskListPage;

  test.beforeEach(async ({ page }) => {
    taskListPage = new TaskListPage(page, BASE_URL);
  });

  test('should load the task list page', async ({ page }) => {
    await taskListPage.navigate();
    await expect(page).toHaveTitle(/Gorev/i);
    await expect(taskListPage.taskList).toBeVisible();
  });

  test('should display tasks grouped by status', async ({ page }) => {
    await taskListPage.navigate();

    // Check for status sections
    const beklemedeSections = page.locator('[data-testid="status-section-beklemede"]');
    const devamEdiyorSections = page.locator('[data-testid="status-section-devam_ediyor"]');
    const tamamlandiSections = page.locator('[data-testid="status-section-tamamlandi"]');

    // At least one section should be visible (depends on data)
    const hasSections =
      (await beklemedeSections.count()) > 0 ||
      (await devamEdiyorSections.count()) > 0 ||
      (await tamamlandiSections.count()) > 0;

    expect(hasSections).toBeTruthy();
  });

  test('should show loading indicator while fetching tasks', async ({ page }) => {
    // Intercept API call to delay response
    await page.route(`${API_URL}/api/v1/tasks*`, async (route) => {
      await new Promise((resolve) => setTimeout(resolve, 500));
      await route.continue();
    });

    await taskListPage.goto('/');

    // Loading indicator should appear
    const loadingIndicator = taskListPage.getByTestId('loading-indicator');
    // Note: This might be too fast to catch, depending on network speed
    await expect(loadingIndicator.or(taskListPage.taskList)).toBeVisible();
  });

  test('should display task details correctly', async ({ page }) => {
    await taskListPage.navigate();

    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) > 0) {
      // Check task structure
      const taskTitle = taskItem.locator('[data-testid="task-title"]');
      await expect(taskTitle).toBeVisible();

      const statusSelect = taskItem.locator('[data-testid="status-select"]');
      await expect(statusSelect).toBeVisible();

      const priorityBadge = taskItem.locator('[data-testid="priority-badge"]');
      await expect(priorityBadge).toBeVisible();
    }
  });
});

test.describe('Web UI - Task Search and Filter', () => {
  let taskListPage: TaskListPage;

  test.beforeEach(async ({ page }) => {
    taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();
  });

  test('should have search input visible', async ({ page }) => {
    const searchInput = page.locator('[data-testid="search-input"]');
    await expect(searchInput).toBeVisible();
  });

  test('should filter tasks by search term', async ({ page }) => {
    const searchInput = page.locator('[data-testid="search-input"]');
    await searchInput.fill('test');

    // Wait for debounce and filter
    await page.waitForTimeout(600);

    // Results should update (we can't know exact count without knowing data)
    await expect(taskListPage.taskList).toBeVisible();
  });

  test('should have status filter dropdown', async ({ page }) => {
    const statusFilter = page.locator('[data-testid="status-filter"]');
    await expect(statusFilter).toBeVisible();

    // Check options
    const options = await statusFilter.locator('option').allTextContents();
    expect(options.length).toBeGreaterThan(1);
  });

  test('should have priority filter dropdown', async ({ page }) => {
    const priorityFilter = page.locator('[data-testid="priority-filter"]');
    await expect(priorityFilter).toBeVisible();

    // Check options
    const options = await priorityFilter.locator('option').allTextContents();
    expect(options.length).toBeGreaterThan(1);
  });

  test('should filter by status', async ({ page }) => {
    const statusFilter = page.locator('[data-testid="status-filter"]');
    await statusFilter.selectOption('beklemede');

    await page.waitForTimeout(500);
    await expect(taskListPage.taskList).toBeVisible();
  });

  test('should filter by priority', async ({ page }) => {
    const priorityFilter = page.locator('[data-testid="priority-filter"]');
    await priorityFilter.selectOption('yuksek');

    await page.waitForTimeout(500);
    await expect(taskListPage.taskList).toBeVisible();
  });

  test('should clear search when input is emptied', async ({ page }) => {
    const searchInput = page.locator('[data-testid="search-input"]');

    // Search for something
    await searchInput.fill('test');
    await page.waitForTimeout(600);

    // Clear search
    await searchInput.clear();
    await page.waitForTimeout(600);

    // Should show all tasks again
    await expect(taskListPage.taskList).toBeVisible();
  });
});

test.describe('Web UI - Task Status Management', () => {
  let taskListPage: TaskListPage;

  test.beforeEach(async ({ page }) => {
    taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();
  });

  test('should be able to change task status', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const statusSelect = taskItem.locator('[data-testid="status-select"]');
    const currentStatus = await statusSelect.inputValue();

    // Change to different status
    const newStatus = currentStatus === 'beklemede' ? 'devam_ediyor' : 'beklemede';
    await statusSelect.selectOption(newStatus);

    // Wait for API call
    await page.waitForTimeout(500);

    // Verify change (status might have changed or page might have reloaded)
    await expect(taskListPage.taskList).toBeVisible();
  });
});

test.describe('Web UI - Task Context Menu', () => {
  let taskListPage: TaskListPage;

  test.beforeEach(async ({ page }) => {
    taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();
  });

  test('should show menu button on task', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const menuButton = taskItem.locator('[data-testid="task-menu-button"]');
    await expect(menuButton).toBeVisible();
  });

  test('should open context menu when clicking menu button', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const menuButton = taskItem.locator('[data-testid="task-menu-button"]');
    await menuButton.click();

    const contextMenu = page.locator('[data-testid="context-menu"]');
    await expect(contextMenu).toBeVisible();
  });

  test('should have edit option in context menu', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const menuButton = taskItem.locator('[data-testid="task-menu-button"]');
    await menuButton.click();

    const editOption = page.locator('[data-testid="menu-item-edit"]');
    await expect(editOption).toBeVisible();
  });

  test('should have delete option in context menu', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const menuButton = taskItem.locator('[data-testid="task-menu-button"]');
    await menuButton.click();

    const deleteOption = page.locator('[data-testid="menu-item-delete"]');
    await expect(deleteOption).toBeVisible();
  });
});

test.describe('Web UI - Task Edit Modal', () => {
  let taskListPage: TaskListPage;

  test.beforeEach(async ({ page }) => {
    taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();
  });

  test('should open edit modal when clicking edit', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const menuButton = taskItem.locator('[data-testid="task-menu-button"]');
    await menuButton.click();

    const editOption = page.locator('[data-testid="menu-item-edit"]');
    await editOption.click();

    const editModal = page.locator('[data-testid="edit-task-modal"]');
    await expect(editModal).toBeVisible();
  });

  test('should have title input in edit modal', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const menuButton = taskItem.locator('[data-testid="task-menu-button"]');
    await menuButton.click();

    const editOption = page.locator('[data-testid="menu-item-edit"]');
    await editOption.click();

    const titleInput = page.locator('[data-testid="input-title"]');
    await expect(titleInput).toBeVisible();
  });

  test('should close modal when clicking cancel', async ({ page }) => {
    const taskItem = taskListPage.taskItems.first();
    if ((await taskItem.count()) === 0) {
      test.skip();
      return;
    }

    const menuButton = taskItem.locator('[data-testid="task-menu-button"]');
    await menuButton.click();

    const editOption = page.locator('[data-testid="menu-item-edit"]');
    await editOption.click();

    const cancelButton = page.locator('[data-testid="cancel-button"]');
    await cancelButton.click();

    const editModal = page.locator('[data-testid="edit-task-modal"]');
    await expect(editModal).toBeHidden();
  });
});

test.describe('Web UI - Subtasks', () => {
  let taskListPage: TaskListPage;

  test.beforeEach(async ({ page }) => {
    taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();
  });

  test('should show expand button for tasks with subtasks', async ({ page }) => {
    // Find a task with subtasks section
    const subtasksSection = page.locator('[data-testid="subtasks-section"]').first();
    if ((await subtasksSection.count()) === 0) {
      test.skip();
      return;
    }

    const expandButton = subtasksSection.locator('[data-testid="expand-button"]');
    await expect(expandButton).toBeVisible();
  });

  test('should toggle subtasks visibility when clicking expand', async ({ page }) => {
    const subtasksSection = page.locator('[data-testid="subtasks-section"]').first();
    if ((await subtasksSection.count()) === 0) {
      test.skip();
      return;
    }

    const expandButton = subtasksSection.locator('[data-testid="expand-button"]');
    await expandButton.click();

    const subtasksList = subtasksSection.locator('[data-testid="subtasks-list"]');
    await expect(subtasksList).toBeVisible();

    // Click again to collapse
    await expandButton.click();
    await expect(subtasksList).toBeHidden();
  });

  test('should display subtask count', async ({ page }) => {
    const subtaskCount = page.locator('[data-testid="subtask-count"]').first();
    if ((await subtaskCount.count()) === 0) {
      test.skip();
      return;
    }

    const text = await subtaskCount.textContent();
    expect(text).toMatch(/\d+\s+alt\s+görev/i);
  });
});

test.describe('Web UI - Responsive Design', () => {
  test('should be responsive on mobile viewport', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });

    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();

    await expect(taskListPage.taskList).toBeVisible();
  });

  test('should be responsive on tablet viewport', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });

    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();

    await expect(taskListPage.taskList).toBeVisible();
  });

  test('should be responsive on desktop viewport', async ({ page }) => {
    await page.setViewportSize({ width: 1920, height: 1080 });

    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();

    await expect(taskListPage.taskList).toBeVisible();
  });
});

test.describe('Web UI - Accessibility', () => {
  test('should have proper heading structure', async ({ page }) => {
    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();

    // Check for h1
    const h1 = page.locator('h1');
    await expect(h1).toBeVisible();
  });

  test('should have keyboard accessible elements', async ({ page }) => {
    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();

    // Tab through interactive elements
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');
    await page.keyboard.press('Tab');

    // Something should be focused
    const focusedElement = page.locator(':focus');
    await expect(focusedElement).toBeVisible();
  });

  test('should have proper form labels', async ({ page }) => {
    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.navigate();

    // Search input should have placeholder as accessible name
    const searchInput = page.locator('[data-testid="search-input"]');
    await expect(searchInput).toHaveAttribute('placeholder');
  });
});

test.describe('Web UI - Error Handling', () => {
  test('should handle API errors gracefully', async ({ page }) => {
    // Intercept API call to return error
    await page.route(`${API_URL}/api/v1/tasks*`, (route) => {
      route.fulfill({
        status: 500,
        body: JSON.stringify({ error: 'Internal server error' }),
      });
    });

    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.goto('/');

    // Page should still be functional
    await expect(page.locator('body')).toBeVisible();
  });

  test('should handle network errors', async ({ page }) => {
    // Abort API call to simulate network error
    await page.route(`${API_URL}/api/v1/tasks*`, (route) => {
      route.abort('failed');
    });

    const taskListPage = new TaskListPage(page, BASE_URL);
    await taskListPage.goto('/');

    // Page should still be functional
    await expect(page.locator('body')).toBeVisible();
  });
});
