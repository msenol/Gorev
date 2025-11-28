import { test, expect } from '@playwright/test';
import MockServer from './mock-server';

let mockServer: MockServer;
let mockServerPort: number;

test.describe('API Integration Tests', () => {
  test.beforeAll(async () => {
    mockServer = new MockServer(); // Port 0 = auto-assign
    await mockServer.start();
    mockServerPort = mockServer.getPort();
    console.log(`[Test] Mock server started on port ${mockServerPort}`);
  });

  test.afterAll(async () => {
    await mockServer.stop();
    console.log('[Test] Mock server stopped');
  });

  test('should connect to health endpoint', async ({ page }) => {
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/health`);
    expect(response.status()).toBe(200);
    const data = await response.json();
    expect(data.status).toBe('ok');
  });

  test('should fetch all tasks', async ({ page }) => {
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
    expect(data.data.length).toBeGreaterThan(0);

    const task = data.data[0];
    expect(task).toHaveProperty('id');
    expect(task).toHaveProperty('title');
    expect(task).toHaveProperty('status');
    expect(task).toHaveProperty('priority');
  });

  test('should create task from template', async ({ page }) => {
    const taskData = {
      template_id: 'bug-report',
      values: {
        title: 'Test Bug',
        description: 'Test description',
        severity: 'high'
      }
    };

    const response = await page.request.post(`http://localhost:${mockServerPort}/api/v1/tasks/from-template`, {
      data: taskData
    });

    expect(response.status()).toBe(200);
    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data).toHaveProperty('id');
    expect(data.data.title).toBe('Test Bug');
  });

  test('should return 404 for non-existent task', async ({ page }) => {
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks/non-existent-id`);
    expect(response.status()).toBe(404);

    const data = await response.json();
    expect(data.success).toBe(false);
  });

  test('should handle CORS headers', async ({ page }) => {
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/health`);
    expect(response.status()).toBe(200);

    expect(response.headers()['access-control-allow-origin']).toBe('*');
    expect(response.headers()['access-control-allow-methods']).toContain('GET');
    expect(response.headers()['access-control-allow-headers']).toContain('Content-Type');
  });
});
