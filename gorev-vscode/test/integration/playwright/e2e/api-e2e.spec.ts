/**
 * API End-to-End Tests
 *
 * Comprehensive E2E tests for the Gorev REST API.
 * Uses the real Go server for integration testing.
 *
 * Run with: SKIP_WEB_SERVER=true npx playwright test e2e/api-e2e.spec.ts
 */

import { test, expect } from '@playwright/test';

// Test configuration
const API_URL = process.env.API_URL || 'http://localhost:5082';

test.describe('API E2E - Health and Status', () => {
  test('should return healthy status', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/health`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.status).toBe('ok');
  });

  test('should return summary statistics', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/summary`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data).toHaveProperty('total_tasks');
    expect(data.data).toHaveProperty('pending_tasks');
    expect(data.data).toHaveProperty('in_progress_tasks');
    expect(data.data).toHaveProperty('completed_tasks');
    expect(data.data).toHaveProperty('total_projects');
  });
});

test.describe('API E2E - Projects', () => {
  let createdProjectId: string | null = null;

  test('should list all projects', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/projects`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
  });

  test('should create a new project', async ({ request }) => {
    const projectData = {
      name: `E2E Test Project ${Date.now()}`,
      description: 'Project created during E2E testing',
    };

    const response = await request.post(`${API_URL}/api/v1/projects`, {
      data: projectData,
    });

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data).toHaveProperty('id');
    expect(data.data.name).toBe(projectData.name);

    createdProjectId = data.data.id;
  });

  test('should get project by ID', async ({ request }) => {
    // First, create a project
    const createResponse = await request.post(`${API_URL}/api/v1/projects`, {
      data: {
        name: `Test Project ${Date.now()}`,
        description: 'Test',
      },
    });
    const created = await createResponse.json();
    const projectId = created.data.id;

    // Then fetch it
    const response = await request.get(`${API_URL}/api/v1/projects/${projectId}`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data.id).toBe(projectId);
  });

  test('should set active project', async ({ request }) => {
    // Get projects
    const listResponse = await request.get(`${API_URL}/api/v1/projects`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const projectId = listData.data[0].id;
    const response = await request.post(`${API_URL}/api/v1/projects/${projectId}/set-active`);
    expect(response.ok()).toBeTruthy();
  });

  test('should get tasks for a project', async ({ request }) => {
    // Get projects
    const listResponse = await request.get(`${API_URL}/api/v1/projects`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const projectId = listData.data[0].id;
    const response = await request.get(`${API_URL}/api/v1/projects/${projectId}/tasks`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
  });
});

test.describe('API E2E - Tasks', () => {
  test('should list all tasks', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/tasks`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
  });

  test('should list tasks with pagination', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/tasks?limit=5&offset=0`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data.length).toBeLessThanOrEqual(5);
  });

  test('should filter tasks by status', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/tasks?status=beklemede`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);

    // All returned tasks should have the filtered status
    for (const task of data.data) {
      expect(task.status).toBe('beklemede');
    }
  });

  test('should filter tasks by priority', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/tasks?priority=yuksek`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);

    // All returned tasks should have the filtered priority
    for (const task of data.data) {
      expect(task.priority).toBe('yuksek');
    }
  });

  test('should create task from template', async ({ request }) => {
    // First get templates
    const templatesResponse = await request.get(`${API_URL}/api/v1/templates`);
    const templatesData = await templatesResponse.json();

    if (!templatesData.data || templatesData.data.length === 0) {
      test.skip();
      return;
    }

    const template = templatesData.data[0];
    const taskData = {
      template_alias: template.alias || 'bug',
      values: {
        title: `E2E Test Task ${Date.now()}`,
        description: 'Task created during E2E testing',
      },
      status: 'beklemede',
      priority: 'orta',
    };

    const response = await request.post(`${API_URL}/api/v1/tasks/from-template`, {
      data: taskData,
    });

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data).toHaveProperty('id');
  });

  test('should get task by ID', async ({ request }) => {
    // Get tasks
    const listResponse = await request.get(`${API_URL}/api/v1/tasks?limit=1`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const taskId = listData.data[0].id;
    const response = await request.get(`${API_URL}/api/v1/tasks/${taskId}`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data.id).toBe(taskId);
  });

  test('should update task status', async ({ request }) => {
    // Get a task
    const listResponse = await request.get(`${API_URL}/api/v1/tasks?limit=1`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const task = listData.data[0];
    const newStatus = task.status === 'beklemede' ? 'devam_ediyor' : 'beklemede';

    const response = await request.put(`${API_URL}/api/v1/tasks/${task.id}`, {
      data: { status: newStatus },
    });

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data.status).toBe(newStatus);
  });

  test('should update task priority', async ({ request }) => {
    // Get a task
    const listResponse = await request.get(`${API_URL}/api/v1/tasks?limit=1`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const task = listData.data[0];
    const newPriority = task.priority === 'yuksek' ? 'orta' : 'yuksek';

    const response = await request.put(`${API_URL}/api/v1/tasks/${task.id}`, {
      data: { priority: newPriority },
    });

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data.priority).toBe(newPriority);
  });

  test('should return 404 for non-existent task', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/tasks/non-existent-id-12345`);
    expect(response.status()).toBe(404);
  });
});

test.describe('API E2E - Subtasks', () => {
  test('should get subtasks for a task', async ({ request }) => {
    // Get tasks with subtasks
    const listResponse = await request.get(`${API_URL}/api/v1/tasks`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const taskId = listData.data[0].id;
    const response = await request.get(`${API_URL}/api/v1/tasks/${taskId}/subtasks`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
  });

  test('should create a subtask', async ({ request }) => {
    // Get a parent task
    const listResponse = await request.get(`${API_URL}/api/v1/tasks?limit=1`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const parentId = listData.data[0].id;
    const subtaskData = {
      title: `E2E Subtask ${Date.now()}`,
      description: 'Subtask created during E2E testing',
      status: 'beklemede',
      priority: 'orta',
    };

    const response = await request.post(`${API_URL}/api/v1/tasks/${parentId}/subtasks`, {
      data: subtaskData,
    });

    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(data.data).toHaveProperty('id');
    expect(data.data.parent_id).toBe(parentId);
  });
});

test.describe('API E2E - Templates', () => {
  test('should list all templates', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/templates`);
    expect(response.ok()).toBeTruthy();

    const data = await response.json();
    expect(data.success).toBe(true);
    expect(Array.isArray(data.data)).toBe(true);
  });

  test('should have expected template structure', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/templates`);
    const data = await response.json();

    if (data.data.length === 0) {
      test.skip();
      return;
    }

    const template = data.data[0];
    expect(template).toHaveProperty('id');
    expect(template).toHaveProperty('name');
  });
});

test.describe('API E2E - Dependencies', () => {
  test('should get task dependencies', async ({ request }) => {
    // Get a task
    const listResponse = await request.get(`${API_URL}/api/v1/tasks?limit=1`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const taskId = listData.data[0].id;
    const response = await request.get(`${API_URL}/api/v1/tasks/${taskId}/dependencies`);

    // May return 200 or 404 depending on implementation
    expect([200, 404]).toContain(response.status());
  });
});

test.describe('API E2E - CORS', () => {
  test('should include CORS headers', async ({ request }) => {
    const response = await request.get(`${API_URL}/api/v1/health`);

    const headers = response.headers();
    expect(headers['access-control-allow-origin']).toBeDefined();
  });
});

test.describe('API E2E - Concurrent Requests', () => {
  test('should handle concurrent task list requests', async ({ request }) => {
    const promises = Array(5)
      .fill(null)
      .map(() => request.get(`${API_URL}/api/v1/tasks`));

    const responses = await Promise.all(promises);

    for (const response of responses) {
      expect(response.ok()).toBeTruthy();
    }
  });

  test('should handle concurrent different endpoint requests', async ({ request }) => {
    const promises = [
      request.get(`${API_URL}/api/v1/health`),
      request.get(`${API_URL}/api/v1/tasks`),
      request.get(`${API_URL}/api/v1/projects`),
      request.get(`${API_URL}/api/v1/templates`),
      request.get(`${API_URL}/api/v1/summary`),
    ];

    const responses = await Promise.all(promises);

    for (const response of responses) {
      expect(response.ok()).toBeTruthy();
    }
  });
});

test.describe('API E2E - Data Consistency', () => {
  test('should maintain data consistency after create and read', async ({ request }) => {
    // Create a project
    const createResponse = await request.post(`${API_URL}/api/v1/projects`, {
      data: {
        name: `Consistency Test ${Date.now()}`,
        description: 'Testing data consistency',
      },
    });

    const created = await createResponse.json();
    const projectId = created.data.id;

    // Read it back
    const readResponse = await request.get(`${API_URL}/api/v1/projects/${projectId}`);
    const read = await readResponse.json();

    expect(read.data.id).toBe(projectId);
    expect(read.data.name).toBe(created.data.name);
  });

  test('should maintain data consistency after update and read', async ({ request }) => {
    // Get a task
    const listResponse = await request.get(`${API_URL}/api/v1/tasks?limit=1`);
    const listData = await listResponse.json();

    if (listData.data.length === 0) {
      test.skip();
      return;
    }

    const taskId = listData.data[0].id;
    const newTitle = `Updated Title ${Date.now()}`;

    // Update it
    await request.put(`${API_URL}/api/v1/tasks/${taskId}`, {
      data: { title: newTitle },
    });

    // Read it back
    const readResponse = await request.get(`${API_URL}/api/v1/tasks/${taskId}`);
    const read = await readResponse.json();

    expect(read.data.title).toBe(newTitle);
  });
});
