/**
 * Task List Page Object
 *
 * Page object for the main task list view in the Web UI.
 */

import { expect, Locator } from '@playwright/test';
import BasePage from './base-page';

export interface TaskInfo {
  title: string;
  description?: string;
  status: string;
  priority: string;
}

export class TaskListPage extends BasePage {
  // Locators
  get taskList(): Locator {
    return this.getByTestId('task-list');
  }

  get taskItems(): Locator {
    return this.getAllByTestId('task-item');
  }

  get createTaskButton(): Locator {
    return this.getByTestId('create-task-button');
  }

  get searchInput(): Locator {
    return this.getByTestId('search-input');
  }

  get refreshButton(): Locator {
    return this.getByTestId('refresh-button');
  }

  get filterButton(): Locator {
    return this.getByTestId('filter-button');
  }

  get projectSelector(): Locator {
    return this.getByTestId('project-selector');
  }

  get loadingIndicator(): Locator {
    return this.getByTestId('loading-indicator');
  }

  get emptyState(): Locator {
    return this.getByTestId('empty-state');
  }

  // Actions

  /**
   * Navigate to task list
   */
  async navigate(): Promise<void> {
    await this.goto('/');
    await this.waitForTaskList();
  }

  /**
   * Wait for task list to be visible
   */
  async waitForTaskList(timeout = 10000): Promise<void> {
    await expect(this.taskList).toBeVisible({ timeout });
    await this.waitForLoadingToFinish();
  }

  /**
   * Get all visible tasks
   */
  async getTaskCount(): Promise<number> {
    return await this.taskItems.count();
  }

  /**
   * Get task by index
   */
  getTaskByIndex(index: number): Locator {
    return this.taskItems.nth(index);
  }

  /**
   * Get task by title
   */
  getTaskByTitle(title: string): Locator {
    return this.taskItems.filter({ hasText: title });
  }

  /**
   * Click on a task to select it
   */
  async selectTask(title: string): Promise<void> {
    await this.getTaskByTitle(title).click();
  }

  /**
   * Double-click on a task to open details
   */
  async openTaskDetails(title: string): Promise<void> {
    await this.getTaskByTitle(title).dblclick();
  }

  /**
   * Right-click on a task to open context menu
   */
  async openTaskContextMenu(title: string): Promise<void> {
    await this.getTaskByTitle(title).click({ button: 'right' });
  }

  /**
   * Search for tasks
   */
  async searchTasks(query: string): Promise<void> {
    await this.searchInput.fill(query);
    await this.page.waitForTimeout(500); // Wait for debounce
    await this.waitForLoadingToFinish();
  }

  /**
   * Clear search
   */
  async clearSearch(): Promise<void> {
    await this.searchInput.clear();
    await this.waitForLoadingToFinish();
  }

  /**
   * Refresh task list
   */
  async refresh(): Promise<void> {
    await this.refreshButton.click();
    await this.waitForLoadingToFinish();
  }

  /**
   * Click create task button
   */
  async clickCreateTask(): Promise<void> {
    await this.createTaskButton.click();
  }

  /**
   * Select a project from the project selector
   */
  async selectProject(projectName: string): Promise<void> {
    await this.projectSelector.click();
    await this.page.locator(`[data-testid="project-option"]:has-text("${projectName}")`).click();
    await this.waitForLoadingToFinish();
  }

  /**
   * Get task status badge text
   */
  async getTaskStatus(title: string): Promise<string | null> {
    const task = this.getTaskByTitle(title);
    const statusBadge = task.locator('[data-testid="status-badge"]');
    return await statusBadge.textContent();
  }

  /**
   * Get task priority
   */
  async getTaskPriority(title: string): Promise<string | null> {
    const task = this.getTaskByTitle(title);
    const priorityBadge = task.locator('[data-testid="priority-badge"]');
    return await priorityBadge.textContent();
  }

  /**
   * Check if task is visible
   */
  async isTaskVisible(title: string): Promise<boolean> {
    const count = await this.getTaskByTitle(title).count();
    return count > 0;
  }

  /**
   * Get all task titles
   */
  async getAllTaskTitles(): Promise<string[]> {
    const titles: string[] = [];
    const count = await this.taskItems.count();
    for (let i = 0; i < count; i++) {
      const titleElement = this.taskItems.nth(i).locator('[data-testid="task-title"]');
      const title = await titleElement.textContent();
      if (title) {
        titles.push(title);
      }
    }
    return titles;
  }

  /**
   * Get task count by status
   */
  async getTaskCountByStatus(status: string): Promise<number> {
    const statusSection = this.page.locator(`[data-testid="status-section-${status}"]`);
    if (await statusSection.count() === 0) {
      return 0;
    }
    const tasks = statusSection.locator('[data-testid="task-item"]');
    return await tasks.count();
  }

  /**
   * Expand task to show subtasks
   */
  async expandTask(title: string): Promise<void> {
    const task = this.getTaskByTitle(title);
    const expandButton = task.locator('[data-testid="expand-button"]');
    if (await expandButton.count() > 0) {
      await expandButton.click();
    }
  }

  /**
   * Get subtask count for a task
   */
  async getSubtaskCount(title: string): Promise<number> {
    const task = this.getTaskByTitle(title);
    const subtasks = task.locator('[data-testid="subtask-item"]');
    return await subtasks.count();
  }

  /**
   * Change task status using dropdown
   */
  async changeTaskStatus(title: string, newStatus: string): Promise<void> {
    const task = this.getTaskByTitle(title);
    const statusSelect = task.locator('select[data-testid="status-select"]');
    await statusSelect.selectOption(newStatus);
    await this.waitForLoadingToFinish();
  }

  /**
   * Delete task from context menu
   */
  async deleteTask(title: string): Promise<void> {
    await this.openTaskContextMenu(title);
    await this.page.locator('[data-testid="menu-item-delete"]').click();
    await this.page.locator('[data-testid="confirm-delete-button"]').click();
    await this.waitForLoadingToFinish();
  }

  /**
   * Check if empty state is shown
   */
  async isEmptyStateVisible(): Promise<boolean> {
    return await this.emptyState.isVisible();
  }

  /**
   * Verify task list contains expected tasks
   */
  async verifyTasksExist(expectedTitles: string[]): Promise<void> {
    for (const title of expectedTitles) {
      await expect(this.getTaskByTitle(title)).toBeVisible();
    }
  }

  /**
   * Get statistics from sidebar
   */
  async getStatistics(): Promise<{ total: number; pending: number; inProgress: number; completed: number }> {
    const total = await this.getTextContent('stat-total');
    const pending = await this.getTextContent('stat-pending');
    const inProgress = await this.getTextContent('stat-in-progress');
    const completed = await this.getTextContent('stat-completed');

    const extractNumber = (text: string | null): number => {
      if (!text) return 0;
      const match = text.match(/\d+/);
      return match ? parseInt(match[0], 10) : 0;
    };

    return {
      total: extractNumber(total),
      pending: extractNumber(pending),
      inProgress: extractNumber(inProgress),
      completed: extractNumber(completed),
    };
  }
}

export default TaskListPage;
