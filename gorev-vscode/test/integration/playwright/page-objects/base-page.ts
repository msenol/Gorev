/**
 * Base Page Object
 *
 * Common functionality shared across all page objects.
 */

import { Page, Locator, expect } from '@playwright/test';

export class BasePage {
  protected page: Page;
  protected baseUrl: string;

  constructor(page: Page, baseUrl: string = 'http://localhost:5001') {
    this.page = page;
    this.baseUrl = baseUrl;
  }

  /**
   * Navigate to the page
   */
  async goto(path: string = ''): Promise<void> {
    await this.page.goto(`${this.baseUrl}${path}`);
  }

  /**
   * Wait for page to load
   */
  async waitForLoad(): Promise<void> {
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Get element by test ID
   */
  getByTestId(testId: string): Locator {
    return this.page.locator(`[data-testid="${testId}"]`);
  }

  /**
   * Wait for element to be visible
   */
  async waitForTestId(testId: string, timeout = 10000): Promise<Locator> {
    const element = this.getByTestId(testId);
    await expect(element).toBeVisible({ timeout });
    return element;
  }

  /**
   * Check if element exists
   */
  async hasTestId(testId: string): Promise<boolean> {
    const count = await this.getByTestId(testId).count();
    return count > 0;
  }

  /**
   * Click element by test ID
   */
  async clickTestId(testId: string): Promise<void> {
    await this.getByTestId(testId).click();
  }

  /**
   * Fill input by test ID
   */
  async fillTestId(testId: string, value: string): Promise<void> {
    await this.getByTestId(testId).fill(value);
  }

  /**
   * Get text content of element
   */
  async getTextContent(testId: string): Promise<string | null> {
    return await this.getByTestId(testId).textContent();
  }

  /**
   * Take a screenshot
   */
  async screenshot(name: string): Promise<Buffer> {
    return await this.page.screenshot({ path: `test-results/screenshots/${name}.png` });
  }

  /**
   * Wait for network to be idle
   */
  async waitForNetwork(): Promise<void> {
    await this.page.waitForLoadState('networkidle');
  }

  /**
   * Get all elements matching test ID
   */
  getAllByTestId(testId: string): Locator {
    return this.page.locator(`[data-testid="${testId}"]`);
  }

  /**
   * Count elements matching test ID
   */
  async countTestId(testId: string): Promise<number> {
    return await this.getByTestId(testId).count();
  }

  /**
   * Wait for loading to finish
   */
  async waitForLoadingToFinish(): Promise<void> {
    const loadingIndicator = this.getByTestId('loading-indicator');
    const count = await loadingIndicator.count();
    if (count > 0) {
      await expect(loadingIndicator).toBeHidden({ timeout: 30000 });
    }
  }

  /**
   * Check for error message
   */
  async hasError(): Promise<boolean> {
    return await this.hasTestId('error-message');
  }

  /**
   * Get error message text
   */
  async getErrorMessage(): Promise<string | null> {
    if (await this.hasError()) {
      return await this.getTextContent('error-message');
    }
    return null;
  }
}

export default BasePage;
