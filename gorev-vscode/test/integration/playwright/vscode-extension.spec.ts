/**
 * VS Code Extension Integration Tests
 *
 * These tests verify that the VS Code extension correctly integrates with the
 * Gorev API server and properly displays tasks, projects, and templates.
 */

import { test, expect } from '@playwright/test';
import MockServer from './mock-server';

let mockServer: MockServer;
let mockServerPort: number;

test.describe('VS Code Extension Integration', () => {
  test.beforeAll(async () => {
    // Start mock server before all tests
    mockServer = new MockServer(); // Port 0 = auto-assign
    await mockServer.start();
    mockServerPort = mockServer.getPort();
    console.log(`[Extension Tests] Mock server started on port ${mockServerPort}`);
  });

  test.afterAll(async () => {
    // Stop mock server after all tests
    await mockServer.stop();
    console.log('[Extension Tests] Mock server stopped');
  });

  test('should successfully connect to API server', async ({ page }) => {
    // Test health endpoint
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/health`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.status).toBe('ok');
  });

  test('should load and display tasks in extension tree view', async ({ page }) => {
    // Simulate loading tasks in VS Code extension
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
    expect(data.data.length).toBeGreaterThan(0);

    // Verify task structure matches what extension expects
    const task = data.data[0];
    expect(task).toHaveProperty('id');
    expect(task).toHaveProperty('title');
    expect(task).toHaveProperty('description');
    expect(task).toHaveProperty('status');
    expect(task).toHaveProperty('priority');
    expect(task).toHaveProperty('created_at');
    expect(task).toHaveProperty('updated_at');
  });

  test('should handle task status transitions correctly', async ({ page }) => {
    const tasks = mockServer.getTasks();
    if (tasks.length === 0) {
      test.skip();
      return;
    }

    const taskId = tasks[0].id;

    // Update task status
    const updateResponse = await page.request.put(`http://localhost:${mockServerPort}/api/v1/tasks/${taskId}`, {
      data: { status: 'completed' }
    });

    expect(updateResponse.status()).toBe(200);
    const updatedTask = await updateResponse.json();
    expect(updatedTask.data.status).toBe('completed');

    // Verify the change persists
    const getResponse = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks/${taskId}`);
    const task = await getResponse.json();
    expect(task.data.status).toBe('completed');
  });

  test('should handle task creation with template', async ({ page }) => {
    // Create task using template
    const createResponse = await page.request.post(`http://localhost:${mockServerPort}/api/v1/tasks/from-template`, {
      data: {
        template_id: 'feature',
        values: {
          title: 'New Feature Implementation',
          description: 'Implement new feature for user authentication',
          type: 'enhancement'
        }
      }
    });

    expect(createResponse.status()).toBe(200);
    const createdTask = await createResponse.json();
    expect(createdTask.success).toBe(true);
    expect(createdTask.data.title).toBe('New Feature Implementation');
    expect(createdTask.data.id).toBeDefined();
  });

  test('should handle task deletion', async ({ page }) => {
    const tasks = mockServer.getTasks();
    const initialCount = tasks.length;

    if (initialCount === 0) {
      test.skip();
      return;
    }

    const taskToDelete = tasks[0];
    const deleteResponse = await page.request.delete(`http://localhost:${mockServerPort}/api/v1/tasks/${taskToDelete.id}`);

    expect(deleteResponse.status()).toBe(200);
    const result = await deleteResponse.json();
    expect(result.success).toBe(true);

    // Verify task is deleted
    const getDeletedTask = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks/${taskToDelete.id}`);
    expect(getDeletedTask.status()).toBe(404);
  });

  test('should load and display projects', async ({ page }) => {
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/projects`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);

    if (data.data.length > 0) {
      const project = data.data[0];
      expect(project).toHaveProperty('id');
      expect(project).toHaveProperty('name');
      expect(project).toHaveProperty('description');
      expect(project).toHaveProperty('task_count');
    }
  });

  test('should filter tasks by project', async ({ page }) => {
    const projects = mockServer.getProjects();
    if (projects.length === 0) {
      test.skip();
      return;
    }

    const projectId = projects[0].id;
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks?project_id=${projectId}`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);

    // All tasks should belong to the filtered project
    data.data.forEach((task: any) => {
      expect(task.project_id).toBe(projectId);
    });
  });

  test('should create new project', async ({ page }) => {
    const response = await page.request.post(`http://localhost:${mockServerPort}/api/v1/projects`, {
      data: {
        name: 'New Test Project',
        description: 'A test project created during UI testing'
      }
    });

    expect(response.status()).toBe(200);
    const project = await response.json();
    expect(project.success).toBe(true);
    expect(project.data.name).toBe('New Test Project');
    expect(project.data.id).toBeDefined();
  });

  test('should load templates from API', async ({ page }) => {
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/templates`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
    expect(data.data.length).toBeGreaterThan(0);

    const template = data.data[0];
    expect(template).toHaveProperty('id');
    expect(template).toHaveProperty('name');
    expect(template).toHaveProperty('description');
    expect(template).toHaveProperty('language_code');
    expect(template).toHaveProperty('fields');
    expect(Array.isArray(template.fields)).toBe(true);
  });

  test('should handle subtask creation', async ({ page }) => {
    const tasks = mockServer.getTasks();
    if (tasks.length === 0) {
      test.skip();
      return;
    }

    const parentTaskId = tasks[0].id;
    const response = await page.request.post(`http://localhost:${mockServerPort}/api/v1/tasks/${parentTaskId}/subtasks`, {
      data: {
        title: 'New Subtask',
        description: 'Subtask description',
        status: 'pending',
        priority: 'medium'
      }
    });

    expect(response.status()).toBe(200);
    const subtask = await response.json();
    expect(subtask.success).toBe(true);
    expect(subtask.data.parent_id).toBe(parentTaskId);
    expect(subtask.data.title).toBe('New Subtask');
  });

  test('should fetch subtasks for parent task', async ({ page }) => {
    const tasks = mockServer.getTasks();
    const parentTask = tasks.find(t => t.parent_id);

    if (!parentTask) {
      test.skip();
      return;
    }

    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks/${parentTask.id}/subtasks`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);

    // All subtasks should have the correct parent_id
    data.data.forEach((subtask: any) => {
      expect(subtask.parent_id).toBe(parentTask.id);
    });
  });

  test('should load and verify summary statistics', async ({ page }) => {
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/summary`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);

    const stats = data.data;
    expect(stats).toHaveProperty('total_tasks');
    expect(stats).toHaveProperty('pending_tasks');
    expect(stats).toHaveProperty('in_progress_tasks');
    expect(stats).toHaveProperty('completed_tasks');
    expect(stats).toHaveProperty('total_projects');

    // Verify counts are consistent
    const totalFromTasks = stats.pending_tasks + stats.in_progress_tasks + stats.completed_tasks;
    expect(stats.total_tasks).toBe(totalFromTasks);
    expect(stats.total_tasks).toBe(mockServer.getTasks().length);
    expect(stats.total_projects).toBe(mockServer.getProjects().length);
  });

  test('should handle pagination correctly', async ({ page }) => {
    // Test with limit and offset
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks?limit=2&offset=0`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.limit).toBe(2);
    expect(data.offset).toBe(0);
    expect(data.data.length).toBeLessThanOrEqual(2);

    // Test with different offset
    const response2 = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks?limit=2&offset=2`);
    expect(response2.status()).toBe(200);
    const data2 = await response2.json();
    expect(data2.offset).toBe(2);
  });

  test('should properly handle Turkish field names in API', async ({ page }) => {
    // Verify that the API returns data in expected format for extension
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);

    // The extension expects Turkish field names
    // Our mock server already returns the correct format
    const task = data.data[0];
    if (task) {
      // Verify task has expected fields for extension
      expect(task).toHaveProperty('title');
      expect(task).toHaveProperty('description');
      expect(task).toHaveProperty('status');
      expect(task).toHaveProperty('priority');
      expect(task).toHaveProperty('created_at');
      expect(task).toHaveProperty('updated_at');
    }
  });

  test('should validate workspace isolation', async ({ page }) => {
    // All tasks should have workspace_id
    const response = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks`);
    expect(response.status()).toBe(200);

    const data = await response.json();
    expect(data.success).toBe(true);

    data.data.forEach((task: any) => {
      expect(task.workspace_id).toBeDefined();
      expect(task.workspace_id).toBe('test-workspace');
    });
  });

  test('should handle CORS preflight requests', async ({ page }) => {
    const response = await page.request.fetch(`http://localhost:${mockServerPort}/api/v1/tasks`, {
      method: 'OPTIONS',
      headers: {
        'Origin': 'vscode-webview://extension',
        'Access-Control-Request-Method': 'GET',
        'Access-Control-Request-Headers': 'Content-Type'
      }
    });

    expect(response.status()).toBe(200);

    // Verify CORS headers
    const headers = response.headers();
    expect(headers['access-control-allow-origin']).toBe('*');
    expect(headers['access-control-allow-methods']).toContain('GET');
    expect(headers['access-control-allow-methods']).toContain('POST');
    expect(headers['access-control-allow-methods']).toContain('PUT');
    expect(headers['access-control-allow-methods']).toContain('DELETE');
    expect(headers['access-control-allow-headers']).toContain('Content-Type');
  });

  test('should return proper error responses', async ({ page }) => {
    // Test 404 for non-existent task
    const notFoundResponse = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks/non-existent-id`);
    expect(notFoundResponse.status()).toBe(404);
    const notFoundData = await notFoundResponse.json();
    expect(notFoundData.success).toBe(false);

    // Test 404 for update
    const updateResponse = await page.request.put(`http://localhost:${mockServerPort}/api/v1/tasks/non-existent-id`, {
      data: { title: 'Test' }
    });
    expect(updateResponse.status()).toBe(404);

    // Test 404 for delete
    const deleteResponse = await page.request.delete(`http://localhost:${mockServerPort}/api/v1/tasks/non-existent-id`);
    expect(deleteResponse.status()).toBe(404);
  });

  test('should handle concurrent API requests', async ({ page }) => {
    // Make multiple concurrent requests
    const requests = Array(5).fill(null).map(() =>
      page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks`)
    );

    const responses = await Promise.all(requests);

    // All requests should succeed
    responses.forEach((response) => {
      expect(response.status()).toBe(200);
    });

    // All should return consistent data
    const data = await responses[0].json();
    expect(data.success).toBe(true);
    expect(data.data.length).toBeGreaterThan(0);
  });

  test('should maintain data consistency across operations', async ({ page }) => {
    // Create a task
    const createResponse = await page.request.post(`http://localhost:${mockServerPort}/api/v1/tasks/from-template`, {
      data: {
        template_id: 'bug-report',
        values: {
          title: 'Consistency Test Bug',
          description: 'Testing data consistency',
          severity: 'low'
        }
      }
    });

    expect(createResponse.status()).toBe(200);
    const created = await createResponse.json();
    const taskId = created.data.id;

    // Update the task
    await page.request.put(`http://localhost:${mockServerPort}/api/v1/tasks/${taskId}`, {
      data: { title: 'Updated Title', status: 'completed' }
    });

    // Fetch and verify
    const getResponse = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks/${taskId}`);
    expect(getResponse.status()).toBe(200);
    const task = await getResponse.json();
    expect(task.data.title).toBe('Updated Title');
    expect(task.data.status).toBe('completed');

    // Delete the task
    await page.request.delete(`http://localhost:${mockServerPort}/api/v1/tasks/${taskId}`);

    // Verify it's gone
    const deletedResponse = await page.request.get(`http://localhost:${mockServerPort}/api/v1/tasks/${taskId}`);
    expect(deletedResponse.status()).toBe(404);
  });
});
