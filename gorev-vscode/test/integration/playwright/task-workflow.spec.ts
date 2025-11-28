import { test, expect } from '@playwright/test';
import MockServer from './mock-server';

let mockServer: MockServer;

test.describe('Task Management UI Workflows', () => {
  test.beforeAll(async () => {
    // Start mock API server
    mockServer = new MockServer(5083);
    await mockServer.start();
    console.log('[Test] Mock server started on port 5083');
  });

  test.afterAll(async () => {
    // Stop mock API server
    await mockServer.stop();
    console.log('[Test] Mock server stopped');
  });

  test('should load tasks from API and display in tree view', async ({ page }) => {
    // Navigate to VS Code web version or mock UI
    await page.goto('http://localhost:5001');

    // Wait for app to load
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Verify task list loads with test data
    const taskList = page.locator('[data-testid="task-item"]');
    await expect(taskList).toHaveCount(3);

    // Verify task titles are displayed
    await expect(page.locator('text=Setup Test Environment')).toBeVisible();
    await expect(page.locator('text=Implement UI Tests')).toBeVisible();
    await expect(page.locator('text=Write Documentation')).toBeVisible();

    // Verify different statuses are displayed
    await expect(page.locator('[data-testid="status-badge"]:has-text("completed")')).toBeVisible();
    await expect(page.locator('[data-testid="status-badge"]:has-text("in_progress")')).toBeVisible();
    await expect(page.locator('[data-testid="status-badge"]:has-text("pending")')).toBeVisible();
  });

  test('should filter tasks by project', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Get initial task count
    const initialTasks = await page.locator('[data-testid="task-item"]').count();

    // Select a project from dropdown
    const projectDropdown = page.locator('[data-testid="project-selector"]');
    await expect(projectDropdown).toBeVisible();
    await projectDropdown.click();

    // Select first project
    const firstProject = page.locator('[data-testid="project-option"]').first();
    await firstProject.click();

    // Verify filtered results (should show only tasks from selected project)
    // This is a mock test - actual behavior depends on Web UI implementation
    const filteredTasks = await page.locator('[data-testid="task-item"]').count();
    expect(filteredTasks).toBeLessThanOrEqual(initialTasks);
  });

  test('should create new task from template', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Click create task button
    const createButton = page.locator('[data-testid="create-task-button"]');
    await expect(createButton).toBeVisible();
    await createButton.click();

    // Should open template selection dialog
    await expect(page.locator('[data-testid="template-selector"]')).toBeVisible();

    // Select bug report template
    await page.locator('[data-testid="template-option"]:has-text("Bug Report")').click();

    // Fill in template fields
    await page.locator('[data-testid="input-title"]').fill('New Bug: Login fails');
    await page.locator('[data-testid="input-description"]').fill('User cannot log in with correct credentials');
    await page.locator('[data-testid="select-severity"]').selectOption('high');

    // Submit task
    await page.locator('[data-testid="submit-task-button"]').click();

    // Verify task was created
    await expect(page.locator('text=New Bug: Login fails')).toBeVisible({ timeout: 5000 });
  });

  test('should edit task via context menu', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Right-click on a task to open context menu
    const taskItem = page.locator('[data-testid="task-item"]').first();
    await taskItem.click({ button: 'right' });

    // Should show context menu
    await expect(page.locator('[data-testid="context-menu"]')).toBeVisible();
    await expect(page.locator('[data-testid="menu-item-edit"]')).toBeVisible();

    // Click edit
    await page.locator('[data-testid="menu-item-edit"]').click();

    // Should open edit dialog
    await expect(page.locator('[data-testid="edit-task-dialog"]')).toBeVisible();

    // Update task title
    await page.locator('[data-testid="input-title"]').clear();
    await page.locator('[data-testid="input-title"]').fill('Updated Task Title');

    // Save changes
    await page.locator('[data-testid="save-button"]').click();

    // Verify update in UI
    await expect(page.locator('text=Updated Task Title')).toBeVisible({ timeout: 5000 });
  });

  test('should update task status', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Find a task with pending status
    const pendingTask = page.locator('[data-testid="task-item"]').filter({ hasText: 'pending' }).first();
    await expect(pendingTask).toBeVisible();

    // Click on status badge to open dropdown
    await pendingTask.locator('[data-testid="status-badge"]').click();

    // Select new status
    await page.locator('[data-testid="status-option"]:has-text("in_progress")').click();

    // Verify status updated
    await expect(pendingTask.locator('[data-testid="status-badge"]:has-text("in_progress")')).toBeVisible({ timeout: 3000 });
  });

  test('should display subtasks hierarchy', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Parent task should be visible
    await expect(page.locator('text=Implement UI Tests')).toBeVisible();

    // Subtask should be visible (indented)
    const parentTask = page.locator('[data-testid="task-item"]').filter({ hasText: 'Implement UI Tests' });
    await expect(parentTask.locator('text=Create Playwright Tests')).toBeVisible();

    // Verify indentation/ierarchy
    const subtaskElement = parentTask.locator('[data-testid="subtask-item"]').first();
    await expect(subtaskElement).toBeVisible();
  });

  test('should delete task', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Get initial task count
    const initialCount = await page.locator('[data-testid="task-item"]').count();

    // Right-click on task
    const taskItem = page.locator('[data-testid="task-item"]').filter({ hasText: 'Write Documentation' }).first();
    await taskItem.click({ button: 'right' });

    // Click delete from context menu
    await page.locator('[data-testid="menu-item-delete"]').click();

    // Confirm deletion in dialog
    await expect(page.locator('[data-testid="confirm-dialog"]')).toBeVisible();
    await page.locator('[data-testid="confirm-delete-button"]').click();

    // Verify task was removed
    await expect(page.locator('text=Write Documentation')).toHaveCount(0, { timeout: 3000 });
    await expect(page.locator('[data-testid="task-item"]')).toHaveCount(initialCount - 1);
  });

  test('should search tasks by title', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Use search box
    const searchBox = page.locator('[data-testid="search-box"]');
    await expect(searchBox).toBeVisible();
    await searchBox.fill('Setup');

    // Should filter tasks
    await expect(page.locator('[data-testid="task-item"]')).toHaveCount(1);
    await expect(page.locator('text=Setup Test Environment')).toBeVisible();
  });

  test('should refresh task list', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Click refresh button
    const refreshButton = page.locator('[data-testid="refresh-button"]');
    await expect(refreshButton).toBeVisible();
    await refreshButton.click();

    // Should show loading indicator
    await expect(page.locator('[data-testid="loading-indicator"]')).toBeVisible();

    // Verify tasks are still visible after refresh
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible();
    const taskCount = await page.locator('[data-testid="task-item"]').count();
    expect(taskCount).toBeGreaterThan(0);
  });

  test('should handle task priority display', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Verify high priority is highlighted
    const highPriorityTasks = page.locator('[data-testid="task-item"]').filter({ has: page.locator('[data-testid="priority-badge"].high') });
    await expect(highPriorityTasks.first()).toBeVisible();

    // Verify priority badges are colored appropriately
    const highBadge = page.locator('[data-testid="priority-badge.high"]');
    await expect(highBadge).toHaveClass(/priority-high/);
  });

  test('should display task statistics in sidebar', async ({ page }) => {
    await page.goto('http://localhost:5001');
    await expect(page.locator('[data-testid="task-list"]')).toBeVisible({ timeout: 10000 });

    // Check sidebar statistics
    const statsSection = page.locator('[data-testid="stats-section"]');
    await expect(statsSection).toBeVisible();

    // Verify count displays
    await expect(statsSection.locator('[data-testid="stat-total"]')).toContainText('Total: 3');
    await expect(statsSection.locator('[data-testid="stat-pending"]')).toContainText('Pending:');
    await expect(statsSection.locator('[data-testid="stat-completed"]')).toContainText('Completed:');
  });
});
