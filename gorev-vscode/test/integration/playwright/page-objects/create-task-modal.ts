/**
 * Create Task Modal Page Object
 *
 * Page object for the task creation modal/form.
 */

import { expect, Locator } from '@playwright/test';
import BasePage from './base-page';

export interface TaskFormData {
  title: string;
  description?: string;
  template?: string;
  status?: string;
  priority?: string;
  dueDate?: string;
  tags?: string[];
}

export class CreateTaskModal extends BasePage {
  // Locators
  get modal(): Locator {
    return this.getByTestId('create-task-modal');
  }

  get templateSelector(): Locator {
    return this.getByTestId('template-selector');
  }

  get templateOptions(): Locator {
    return this.getAllByTestId('template-option');
  }

  get titleInput(): Locator {
    return this.getByTestId('input-title');
  }

  get descriptionInput(): Locator {
    return this.getByTestId('input-description');
  }

  get statusSelect(): Locator {
    return this.getByTestId('select-status');
  }

  get prioritySelect(): Locator {
    return this.getByTestId('select-priority');
  }

  get dueDateInput(): Locator {
    return this.getByTestId('input-due-date');
  }

  get tagsInput(): Locator {
    return this.getByTestId('input-tags');
  }

  get submitButton(): Locator {
    return this.getByTestId('submit-task-button');
  }

  get cancelButton(): Locator {
    return this.getByTestId('cancel-button');
  }

  get formError(): Locator {
    return this.getByTestId('form-error');
  }

  // Actions

  /**
   * Wait for modal to be visible
   */
  async waitForModal(timeout = 10000): Promise<void> {
    await expect(this.modal).toBeVisible({ timeout });
  }

  /**
   * Wait for modal to be hidden
   */
  async waitForModalClosed(timeout = 10000): Promise<void> {
    await expect(this.modal).toBeHidden({ timeout });
  }

  /**
   * Select a template
   */
  async selectTemplate(templateName: string): Promise<void> {
    await this.templateSelector.click();
    await this.page.locator(`[data-testid="template-option"]:has-text("${templateName}")`).click();
  }

  /**
   * Fill the title field
   */
  async fillTitle(title: string): Promise<void> {
    await this.titleInput.fill(title);
  }

  /**
   * Fill the description field
   */
  async fillDescription(description: string): Promise<void> {
    await this.descriptionInput.fill(description);
  }

  /**
   * Select status
   */
  async selectStatus(status: string): Promise<void> {
    await this.statusSelect.selectOption(status);
  }

  /**
   * Select priority
   */
  async selectPriority(priority: string): Promise<void> {
    await this.prioritySelect.selectOption(priority);
  }

  /**
   * Set due date
   */
  async setDueDate(date: string): Promise<void> {
    await this.dueDateInput.fill(date);
  }

  /**
   * Add tags
   */
  async addTags(tags: string[]): Promise<void> {
    for (const tag of tags) {
      await this.tagsInput.fill(tag);
      await this.page.keyboard.press('Enter');
    }
  }

  /**
   * Fill the entire form
   */
  async fillForm(data: TaskFormData): Promise<void> {
    if (data.template) {
      await this.selectTemplate(data.template);
    }

    await this.fillTitle(data.title);

    if (data.description) {
      await this.fillDescription(data.description);
    }

    if (data.status) {
      await this.selectStatus(data.status);
    }

    if (data.priority) {
      await this.selectPriority(data.priority);
    }

    if (data.dueDate) {
      await this.setDueDate(data.dueDate);
    }

    if (data.tags && data.tags.length > 0) {
      await this.addTags(data.tags);
    }
  }

  /**
   * Submit the form
   */
  async submit(): Promise<void> {
    await this.submitButton.click();
  }

  /**
   * Cancel and close the modal
   */
  async cancel(): Promise<void> {
    await this.cancelButton.click();
  }

  /**
   * Create a task with the given data
   */
  async createTask(data: TaskFormData): Promise<void> {
    await this.waitForModal();
    await this.fillForm(data);
    await this.submit();
    await this.waitForModalClosed();
  }

  /**
   * Check if form has validation errors
   */
  async hasValidationError(): Promise<boolean> {
    return await this.formError.isVisible();
  }

  /**
   * Get validation error message
   */
  async getValidationError(): Promise<string | null> {
    if (await this.hasValidationError()) {
      return await this.formError.textContent();
    }
    return null;
  }

  /**
   * Check if submit button is enabled
   */
  async isSubmitEnabled(): Promise<boolean> {
    return await this.submitButton.isEnabled();
  }

  /**
   * Get available templates
   */
  async getAvailableTemplates(): Promise<string[]> {
    await this.templateSelector.click();
    const templates: string[] = [];
    const count = await this.templateOptions.count();
    for (let i = 0; i < count; i++) {
      const text = await this.templateOptions.nth(i).textContent();
      if (text) {
        templates.push(text);
      }
    }
    await this.page.keyboard.press('Escape'); // Close dropdown
    return templates;
  }
}

export default CreateTaskModal;
