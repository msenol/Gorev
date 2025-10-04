const assert = require('assert');
const TestHelper = require('../utils/testHelper');

suite('ApiClient Test Suite', () => {
  let helper;
  let mockApiClient;
  let mockAxios;

  setup(() => {
    helper = new TestHelper();
    const result = helper.createMockAPIClient();
    mockApiClient = result.client;
    mockAxios = result.mockAxios;
    helper.setupMockAPIClient(mockAxios);
  });

  teardown(() => {
    helper.cleanup();
  });

  suite('Connection', () => {
    test('should create API client instance', () => {
      assert(mockApiClient);
      assert.strictEqual(mockApiClient.baseURL, 'http://localhost:5082');
    });

    test('should have axios instance', () => {
      assert(mockApiClient.axiosInstance);
    });

    test('should respond to health check', async () => {
      const response = await mockApiClient.axiosInstance.get('/health');
      assert.strictEqual(response.status, 200);
      assert.strictEqual(response.data.status, 'ok');
    });
  });

  suite('Task Operations', () => {
    test('should fetch tasks', async () => {
      const response = await mockApiClient.axiosInstance.get('/tasks');
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
      assert(Array.isArray(response.data.data));
    });

    test('should fetch single task', async () => {
      const response = await mockApiClient.axiosInstance.get('/tasks/test-id');
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
      assert(response.data.data);
    });

    test('should create task', async () => {
      const response = await mockApiClient.axiosInstance.post('/tasks', {
        baslik: 'Test Task'
      });
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });

    test('should update task', async () => {
      const response = await mockApiClient.axiosInstance.put('/tasks/test-id', {
        durum: 'tamamlandi'
      });
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });

    test('should delete task', async () => {
      const response = await mockApiClient.axiosInstance.delete('/tasks/test-id');
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });
  });

  suite('Project Operations', () => {
    test('should fetch projects', async () => {
      const response = await mockApiClient.axiosInstance.get('/projects');
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });

    test('should fetch active project', async () => {
      const response = await mockApiClient.axiosInstance.get('/projects/active');
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });

    test('should create project', async () => {
      const response = await mockApiClient.axiosInstance.post('/projects', {
        isim: 'Test Project'
      });
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });
  });

  suite('Template Operations', () => {
    test('should fetch templates', async () => {
      const response = await mockApiClient.axiosInstance.get('/templates');
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });

    test('should create task from template', async () => {
      const response = await mockApiClient.axiosInstance.post('/tasks/from-template', {
        template_id: 'test',
        degerler: {}
      });
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });
  });

  suite('Export/Import Operations', () => {
    test('should export data', async () => {
      const response = await mockApiClient.axiosInstance.post('/export', {
        format: 'json'
      });
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });

    test('should import data', async () => {
      const response = await mockApiClient.axiosInstance.post('/import', {
        data: {}
      });
      assert.strictEqual(response.status, 200);
      assert(response.data.success);
    });
  });
});
