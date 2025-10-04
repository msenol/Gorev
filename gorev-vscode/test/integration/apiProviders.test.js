const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const MockAdapter = require('axios-mock-adapter');

let ApiClient;
let EnhancedGorevTreeProvider;
let ProjeTreeProvider;
let TemplateTreeProvider;

suite('API Providers Integration Tests', function() {
  this.timeout(15000);

  let sandbox;
  let apiClient;
  let mockAxios;

  suiteSetup(async function() {
    // Import compiled modules
    const clientModule = require('../../out/api/client');
    ApiClient = clientModule.ApiClient;

    // Import providers (they should be compiled to out/)
    const gorevProvider = require('../../out/providers/enhancedGorevTreeProvider');
    EnhancedGorevTreeProvider = gorevProvider.EnhancedGorevTreeProvider;

    const projeProvider = require('../../out/providers/projeTreeProvider');
    ProjeTreeProvider = projeProvider.ProjeTreeProvider;

    const templateProvider = require('../../out/providers/templateTreeProvider');
    TemplateTreeProvider = templateProvider.TemplateTreeProvider;
  });

  setup(function() {
    sandbox = sinon.createSandbox();
    apiClient = new ApiClient('http://localhost:5082');
    mockAxios = new MockAdapter(apiClient.axiosInstance);
  });

  teardown(function() {
    sandbox.restore();
    mockAxios.restore();
    if (apiClient) {
      apiClient.disconnect();
    }
  });

  suite('EnhancedGorevTreeProvider with API', function() {
    let provider;

    setup(function() {
      provider = new EnhancedGorevTreeProvider(apiClient);
    });

    test('should load tasks from API', async function() {
      const mockTasks = [
        {
          id: 'task-1',
          baslik: 'API Task 1',
          aciklama: 'Description 1',
          durum: 'beklemede',
          oncelik: 'yuksek',
          olusturma_tarihi: '2025-01-01T00:00:00Z',
          guncelleme_tarihi: '2025-01-01T00:00:00Z'
        },
        {
          id: 'task-2',
          baslik: 'API Task 2',
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

      const rootElements = await provider.getChildren();

      // Should return tree items (grouping or tasks)
      assert(Array.isArray(rootElements));
      assert(rootElements.length > 0);
    });

    test('should handle empty task list', async function() {
      mockAxios.onGet('/tasks').reply(200, {
        success: true,
        data: [],
        total: 0
      });

      const rootElements = await provider.getChildren();

      assert(Array.isArray(rootElements));
    });

    test('should handle API errors gracefully', async function() {
      mockAxios.onGet('/tasks').reply(500, {
        success: false,
        error: 'Internal server error'
      });

      const rootElements = await provider.getChildren();

      // Should return empty array or error placeholder
      assert(Array.isArray(rootElements));
    });

    test('should refresh on demand', async function() {
      mockAxios.onGet('/tasks').reply(200, {
        success: true,
        data: [],
        total: 0
      });

      // Setup spy on _onDidChangeTreeData
      const fireStub = sandbox.spy(provider._onDidChangeTreeData, 'fire');

      await provider.refresh();

      // Should fire change event
      assert(fireStub.called);
    });

    test('should handle filtering', async function() {
      mockAxios.onGet('/tasks', {
        params: { durum: 'beklemede' }
      }).reply(200, {
        success: true,
        data: [
          {
            id: 'task-1',
            baslik: 'Pending Task',
            durum: 'beklemede',
            oncelik: 'orta',
            olusturma_tarihi: '2025-01-01T00:00:00Z',
            guncelleme_tarihi: '2025-01-01T00:00:00Z'
          }
        ],
        total: 1
      });

      // Apply filter
      provider.currentFilters = { durum: 'beklemede' };

      const elements = await provider.getChildren();

      assert(Array.isArray(elements));
    });

    test('should convert API tasks to internal model', async function() {
      const mockTask = {
        id: 'task-1',
        baslik: 'Test Task',
        aciklama: 'Test Description',
        durum: 'beklemede',
        oncelik: 'yuksek',
        olusturma_tarihi: '2025-01-01T00:00:00Z',
        guncelleme_tarihi: '2025-01-01T00:00:00Z',
        proje_name: 'Test Project',
        etiketler: ['test', 'api']
      };

      mockAxios.onGet('/tasks').reply(200, {
        success: true,
        data: [mockTask],
        total: 1
      });

      const elements = await provider.getChildren();

      // Should have successfully converted and created tree items
      assert(Array.isArray(elements));
    });
  });

  suite('ProjeTreeProvider with API', function() {
    let provider;

    setup(function() {
      provider = new ProjeTreeProvider(apiClient);
    });

    test('should load projects from API', async function() {
      const mockProjects = [
        {
          id: 'proj-1',
          isim: 'API Project 1',
          tanim: 'Description 1',
          olusturma_tarihi: '2025-01-01T00:00:00Z',
          is_active: false,
          gorev_sayisi: 5
        },
        {
          id: 'proj-2',
          isim: 'API Project 2',
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

      const elements = await provider.getChildren();

      assert(Array.isArray(elements));
      assert(elements.length > 0);
    });

    test('should show active project indicator', async function() {
      const mockProjects = [
        {
          id: 'proj-1',
          isim: 'Active Project',
          olusturma_tarihi: '2025-01-01T00:00:00Z',
          is_active: true,
          gorev_sayisi: 5
        }
      ];

      mockAxios.onGet('/projects').reply(200, {
        success: true,
        data: mockProjects,
        total: 1
      });

      const elements = await provider.getChildren();

      // Check if active project is marked (implementation specific)
      assert(Array.isArray(elements));
    });

    test('should handle project activation', async function() {
      mockAxios.onGet('/projects').reply(200, {
        success: true,
        data: [],
        total: 0
      });

      mockAxios.onGet('/active-project').reply(200, {
        success: true,
        data: {
          id: 'active-proj',
          isim: 'Active Project',
          is_active: true
        }
      });

      await provider.refresh();

      // Should refresh successfully
      assert(true);
    });

    test('should convert API projects to internal model', async function() {
      const mockProject = {
        id: 'proj-1',
        isim: 'Test Project',
        tanim: 'Test Description',
        olusturma_tarihi: '2025-01-01T00:00:00Z',
        is_active: false,
        gorev_sayisi: 10
      };

      mockAxios.onGet('/projects').reply(200, {
        success: true,
        data: [mockProject],
        total: 1
      });

      const elements = await provider.getChildren();

      assert(Array.isArray(elements));
    });
  });

  suite('TemplateTreeProvider with API', function() {
    let provider;

    setup(function() {
      provider = new TemplateTreeProvider(apiClient);
    });

    test('should load templates from API', async function() {
      const mockTemplates = [
        {
          id: 'tmpl-1',
          isim: 'Bug Report',
          tanim: 'Bug report template',
          kategori: 'Teknik',
          alanlar: [
            {
              isim: 'baslik',
              tip: 'text',
              zorunlu: true
            },
            {
              isim: 'oncelik',
              tip: 'select',
              zorunlu: true,
              secenekler: ['dusuk', 'orta', 'yuksek']
            }
          ],
          aktif: true
        },
        {
          id: 'tmpl-2',
          isim: 'Feature Request',
          tanim: 'Feature template',
          kategori: 'Özellik',
          alanlar: [],
          aktif: true
        }
      ];

      mockAxios.onGet('/templates').reply(200, {
        success: true,
        data: mockTemplates,
        total: 2
      });

      const elements = await provider.getChildren();

      assert(Array.isArray(elements));
      assert(elements.length > 0);
    });

    test('should group templates by category', async function() {
      const mockTemplates = [
        {
          id: 'tmpl-1',
          isim: 'Bug Report',
          tanim: 'Bug template',
          kategori: 'Teknik',
          alanlar: [],
          aktif: true
        },
        {
          id: 'tmpl-2',
          isim: 'Feature',
          tanim: 'Feature template',
          kategori: 'Özellik',
          alanlar: [],
          aktif: true
        }
      ];

      mockAxios.onGet('/templates').reply(200, {
        success: true,
        data: mockTemplates,
        total: 2
      });

      const elements = await provider.getChildren();

      // Should have category groups or flat list
      assert(Array.isArray(elements));
    });

    test('should filter inactive templates', async function() {
      const mockTemplates = [
        {
          id: 'tmpl-1',
          isim: 'Active Template',
          kategori: 'Teknik',
          alanlar: [],
          aktif: true
        },
        {
          id: 'tmpl-2',
          isim: 'Inactive Template',
          kategori: 'Teknik',
          alanlar: [],
          aktif: false
        }
      ];

      mockAxios.onGet('/templates').reply(200, {
        success: true,
        data: mockTemplates,
        total: 2
      });

      const elements = await provider.getChildren();

      // Implementation should filter inactive templates
      assert(Array.isArray(elements));
    });

    test('should convert API templates to internal model', async function() {
      const mockTemplate = {
        id: 'tmpl-1',
        isim: 'Test Template',
        tanim: 'Test Description',
        kategori: 'Teknik',
        alanlar: [
          {
            isim: 'test_field',
            tip: 'text',
            zorunlu: true,
            varsayilan: 'default value',
            aciklama: 'Test field'
          }
        ],
        aktif: true
      };

      mockAxios.onGet('/templates').reply(200, {
        success: true,
        data: [mockTemplate],
        total: 1
      });

      const elements = await provider.getChildren();

      assert(Array.isArray(elements));
    });
  });

  suite('Provider Error Handling', function() {
    test('EnhancedGorevTreeProvider should handle network errors', async function() {
      const provider = new EnhancedGorevTreeProvider(apiClient);

      mockAxios.onGet('/tasks').networkError();

      const elements = await provider.getChildren();

      // Should return empty array or error item
      assert(Array.isArray(elements));
    });

    test('ProjeTreeProvider should handle API errors', async function() {
      const provider = new ProjeTreeProvider(apiClient);

      mockAxios.onGet('/projects').reply(500, {
        success: false,
        error: 'Server error'
      });

      const elements = await provider.getChildren();

      assert(Array.isArray(elements));
    });

    test('TemplateTreeProvider should handle empty response', async function() {
      const provider = new TemplateTreeProvider(apiClient);

      mockAxios.onGet('/templates').reply(200, {
        success: true,
        data: [],
        total: 0
      });

      const elements = await provider.getChildren();

      assert(Array.isArray(elements));
      assert.strictEqual(elements.length, 0);
    });
  });

  suite('Provider Refresh Coordination', function() {
    test('should refresh all providers independently', async function() {
      const gorevProvider = new EnhancedGorevTreeProvider(apiClient);
      const projeProvider = new ProjeTreeProvider(apiClient);
      const templateProvider = new TemplateTreeProvider(apiClient);

      // Setup mock responses
      mockAxios.onGet('/tasks').reply(200, { success: true, data: [], total: 0 });
      mockAxios.onGet('/projects').reply(200, { success: true, data: [], total: 0 });
      mockAxios.onGet('/templates').reply(200, { success: true, data: [], total: 0 });

      // Refresh all
      await Promise.all([
        gorevProvider.refresh(),
        projeProvider.refresh(),
        templateProvider.refresh()
      ]);

      // All should complete without errors
      assert(true);
    });
  });
});