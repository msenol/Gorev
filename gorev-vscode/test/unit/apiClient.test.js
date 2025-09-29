const assert = require('assert');
const MockAdapter = require('axios-mock-adapter');

// We'll use dynamic import for ES modules in test
let ApiClient, ApiError;

suite('ApiClient Test Suite', function() {
  this.timeout(10000);

  let client;
  let mockAxios;

  suiteSetup(async function() {
    // Dynamically import the compiled module
    const clientModule = require('../../out/api/client');
    ApiClient = clientModule.ApiClient;
    ApiError = clientModule.ApiError;
  });

  setup(function() {
    // Create fresh client for each test
    client = new ApiClient('http://localhost:5082');

    // Create axios mock adapter
    mockAxios = new MockAdapter(client.axiosInstance);
  });

  teardown(function() {
    mockAxios.restore();
    if (client) {
      client.disconnect();
    }
  });

  suite('Connection Management', function() {
    test('should connect successfully', async function() {
      mockAxios.onGet('/health').reply(200, { status: 'ok' });

      await client.connect();
      assert.strictEqual(client.isConnected(), true);
    });

    test('should handle connection failure', async function() {
      mockAxios.onGet('/health').networkError();

      try {
        await client.connect();
        assert.fail('Should have thrown an error');
      } catch (error) {
        assert(error instanceof Error);
      }
    });

    test('should disconnect gracefully', async function() {
      mockAxios.onGet('/health').reply(200, { status: 'ok' });
      await client.connect();

      client.disconnect();
      assert.strictEqual(client.isConnected(), false);
    });
  });

  suite('Task Operations', function() {
    test('getTasks - should fetch tasks successfully', async function() {
      const mockTasks = [
        {
          id: 'task-1',
          baslik: 'Test Task 1',
          aciklama: 'Description 1',
          durum: 'beklemede',
          oncelik: 'yuksek',
          olusturma_tarihi: '2025-01-01T00:00:00Z',
          guncelleme_tarihi: '2025-01-01T00:00:00Z'
        },
        {
          id: 'task-2',
          baslik: 'Test Task 2',
          aciklama: 'Description 2',
          durum: 'devam_ediyor',
          oncelik: 'orta',
          olusturma_tarihi: '2025-01-02T00:00:00Z',
          guncelleme_tarihi: '2025-01-02T00:00:00Z'
        }
      ];

      mockAxios.onGet('/tasks').reply(200, {
        success: true,
        data: mockTasks,
        total: 2
      });

      const response = await client.getTasks();

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.length, 2);
      assert.strictEqual(response.data[0].id, 'task-1');
      assert.strictEqual(response.data[1].id, 'task-2');
    });

    test('getTasks - should handle filters', async function() {
      mockAxios.onGet('/tasks', {
        params: { durum: 'beklemede', oncelik: 'yuksek' }
      }).reply(200, {
        success: true,
        data: [],
        total: 0
      });

      const response = await client.getTasks({
        durum: 'beklemede',
        oncelik: 'yuksek'
      });

      assert.strictEqual(response.success, true);
    });

    test('getTask - should fetch single task', async function() {
      const mockTask = {
        id: 'task-1',
        baslik: 'Test Task',
        aciklama: 'Description',
        durum: 'beklemede',
        oncelik: 'yuksek',
        olusturma_tarihi: '2025-01-01T00:00:00Z',
        guncelleme_tarihi: '2025-01-01T00:00:00Z'
      };

      mockAxios.onGet('/tasks/task-1').reply(200, {
        success: true,
        data: mockTask
      });

      const response = await client.getTask('task-1');

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.id, 'task-1');
      assert.strictEqual(response.data.baslik, 'Test Task');
    });

    test('getTask - should handle not found', async function() {
      mockAxios.onGet('/tasks/nonexistent').reply(404, {
        success: false,
        error: 'Task not found'
      });

      try {
        await client.getTask('nonexistent');
        assert.fail('Should have thrown ApiError');
      } catch (error) {
        assert(error instanceof ApiError);
        assert.strictEqual(error.statusCode, 404);
        assert(error.isNotFound());
      }
    });

    test('createTask - should create task successfully', async function() {
      const taskData = {
        baslik: 'New Task',
        aciklama: 'New Description',
        oncelik: 'yuksek'
      };

      mockAxios.onPost('/tasks').reply(201, {
        success: true,
        data: {
          id: 'new-task-id',
          ...taskData,
          durum: 'beklemede',
          olusturma_tarihi: '2025-01-01T00:00:00Z',
          guncelleme_tarihi: '2025-01-01T00:00:00Z'
        }
      });

      const response = await client.createTask(taskData);

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.baslik, 'New Task');
    });

    test('updateTask - should update task successfully', async function() {
      mockAxios.onPut('/tasks/task-1').reply(200, {
        success: true,
        data: {
          id: 'task-1',
          baslik: 'Updated Task',
          durum: 'tamamlandi'
        }
      });

      const response = await client.updateTask('task-1', { durum: 'tamamlandi' });

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.durum, 'tamamlandi');
    });

    test('deleteTask - should delete task successfully', async function() {
      mockAxios.onDelete('/tasks/task-1').reply(200, {
        success: true,
        message: 'Task deleted'
      });

      const response = await client.deleteTask('task-1');

      assert.strictEqual(response.success, true);
    });
  });

  suite('Subtask Operations', function() {
    test('createSubtask - should create subtask successfully', async function() {
      const subtaskData = {
        baslik: 'Subtask 1',
        aciklama: 'Subtask description'
      };

      mockAxios.onPost('/tasks/parent-1/subtasks').reply(201, {
        success: true,
        data: {
          id: 'subtask-1',
          ...subtaskData,
          parent_id: 'parent-1'
        }
      });

      const response = await client.createSubtask('parent-1', subtaskData);

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.baslik, 'Subtask 1');
    });

    test('getTaskHierarchy - should fetch hierarchy', async function() {
      mockAxios.onGet('/tasks/task-1/hierarchy').reply(200, {
        success: true,
        data: {
          gorev: { id: 'task-1', baslik: 'Parent' },
          alt_gorevler: [
            { id: 'child-1', baslik: 'Child 1' },
            { id: 'child-2', baslik: 'Child 2' }
          ],
          toplam_alt_gorev: 2,
          tamamlanan_alt_gorev: 1
        }
      });

      const response = await client.getTaskHierarchy('task-1');

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.alt_gorevler.length, 2);
    });

    test('changeParent - should change parent successfully', async function() {
      mockAxios.onPut('/tasks/task-1/parent').reply(200, {
        success: true,
        data: { id: 'task-1', parent_id: 'new-parent' }
      });

      const response = await client.changeParent('task-1', 'new-parent');

      assert.strictEqual(response.success, true);
    });
  });

  suite('Project Operations', function() {
    test('getProjects - should fetch projects', async function() {
      const mockProjects = [
        {
          id: 'proj-1',
          isim: 'Project 1',
          tanim: 'Description 1',
          olusturma_tarihi: '2025-01-01T00:00:00Z',
          is_active: false,
          gorev_sayisi: 5
        },
        {
          id: 'proj-2',
          isim: 'Project 2',
          tanim: 'Description 2',
          olusturma_tarihi: '2025-01-02T00:00:00Z',
          is_active: true,
          gorev_sayisi: 3
        }
      ];

      mockAxios.onGet('/projects').reply(200, {
        success: true,
        data: mockProjects,
        total: 2
      });

      const response = await client.getProjects();

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.length, 2);
    });

    test('createProject - should create project', async function() {
      mockAxios.onPost('/projects').reply(201, {
        success: true,
        data: {
          id: 'new-proj',
          isim: 'New Project',
          tanim: 'New Description'
        }
      });

      const response = await client.createProject({
        isim: 'New Project',
        tanim: 'New Description'
      });

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.isim, 'New Project');
    });

    test('getActiveProject - should fetch active project', async function() {
      mockAxios.onGet('/active-project').reply(200, {
        success: true,
        data: {
          id: 'active-proj',
          isim: 'Active Project',
          is_active: true
        }
      });

      const response = await client.getActiveProject();

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.is_active, true);
    });

    test('activateProject - should activate project', async function() {
      mockAxios.onPost('/active-project').reply(200, {
        success: true,
        message: 'Project activated'
      });

      const response = await client.activateProject('proj-1');

      assert.strictEqual(response.success, true);
    });

    test('removeActiveProject - should remove active project', async function() {
      mockAxios.onDelete('/active-project').reply(200, {
        success: true,
        message: 'Active project removed'
      });

      const response = await client.removeActiveProject();

      assert.strictEqual(response.success, true);
    });
  });

  suite('Template Operations', function() {
    test('getTemplates - should fetch templates', async function() {
      const mockTemplates = [
        {
          id: 'tmpl-1',
          isim: 'Bug Report',
          tanim: 'Bug template',
          kategori: 'Teknik',
          alanlar: [],
          aktif: true
        }
      ];

      mockAxios.onGet('/templates').reply(200, {
        success: true,
        data: mockTemplates,
        total: 1
      });

      const response = await client.getTemplates();

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.length, 1);
    });

    test('createTaskFromTemplate - should create task from template', async function() {
      mockAxios.onPost('/tasks/from-template').reply(201, {
        success: true,
        data: {
          id: 'new-task',
          baslik: 'Bug: Login issue'
        }
      });

      const response = await client.createTaskFromTemplate({
        template_id: 'tmpl-1',
        degerler: { baslik: 'Login issue' }
      });

      assert.strictEqual(response.success, true);
    });
  });

  suite('Dependency Operations', function() {
    test('addDependency - should add dependency', async function() {
      mockAxios.onPost('/tasks/task-1/dependencies').reply(201, {
        success: true,
        message: 'Dependency added'
      });

      const response = await client.addDependency('task-1', {
        kaynak_id: 'task-2',
        baglanti_tipi: 'onceki'
      });

      assert.strictEqual(response.success, true);
    });
  });

  suite('Summary Operations', function() {
    test('getSummary - should fetch summary', async function() {
      mockAxios.onGet('/summary').reply(200, {
        success: true,
        data: {
          toplam_proje: 5,
          toplam_gorev: 20,
          durum_dagilimi: {
            beklemede: 10,
            devam_ediyor: 7,
            tamamlandi: 3
          },
          oncelik_dagilimi: {
            yuksek: 5,
            orta: 10,
            dusuk: 5
          }
        }
      });

      const response = await client.getSummary();

      assert.strictEqual(response.success, true);
      assert.strictEqual(response.data.toplam_proje, 5);
      assert.strictEqual(response.data.toplam_gorev, 20);
    });
  });

  suite('Error Handling', function() {
    test('ApiError - should have helper methods', function() {
      const notFoundError = new ApiError(404, 'Not found', '/tasks/123');
      assert(notFoundError.isNotFound());
      assert(!notFoundError.isBadRequest());
      assert(!notFoundError.isServerError());

      const badRequestError = new ApiError(400, 'Bad request', '/tasks');
      assert(!badRequestError.isNotFound());
      assert(badRequestError.isBadRequest());
      assert(!badRequestError.isServerError());

      const serverError = new ApiError(500, 'Server error', '/tasks');
      assert(!serverError.isNotFound());
      assert(!serverError.isBadRequest());
      assert(serverError.isServerError());
    });

    test('should handle network errors gracefully', async function() {
      mockAxios.onGet('/tasks').networkError();

      try {
        await client.getTasks();
        assert.fail('Should have thrown an error');
      } catch (error) {
        assert(error instanceof Error);
      }
    });

    test('should handle timeout errors', async function() {
      mockAxios.onGet('/tasks').timeout();

      try {
        await client.getTasks();
        assert.fail('Should have thrown an error');
      } catch (error) {
        assert(error instanceof Error);
      }
    });

    test('should handle API error responses', async function() {
      mockAxios.onGet('/tasks/invalid').reply(400, {
        success: false,
        error: 'Invalid task ID'
      });

      try {
        await client.getTask('invalid');
        assert.fail('Should have thrown ApiError');
      } catch (error) {
        assert(error instanceof ApiError);
        assert.strictEqual(error.statusCode, 400);
        assert(error.isBadRequest());
      }
    });
  });
});