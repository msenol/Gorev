const vscode = require('vscode');
const sinon = require('sinon');
const MockAdapter = require('axios-mock-adapter');

class TestHelper {
  constructor() {
    this.sandbox = sinon.createSandbox();
    this.stubs = {};
  }

  /**
   * Create a mock API client with axios mock adapter
   */
  createMockAPIClient() {
    try {
      const clientModule = require('../../out/api/client');
      const ApiClient = clientModule.ApiClient;

      const client = new ApiClient('http://localhost:5082');
      const mockAxios = new MockAdapter(client.axiosInstance);

      return { client, mockAxios };
    } catch (error) {
      // If compiled code not available, create a mock client
      const axios = require('axios');
      const axiosInstance = axios.create({ baseURL: 'http://localhost:5082' });
      const mockAxios = new MockAdapter(axiosInstance);

      const mockClient = {
        axiosInstance,
        baseURL: 'http://localhost:5082',
        isConnected: () => true,
        connect: async () => {},
        disconnect: () => {},
        on: () => {},
        emit: () => {}
      };

      return { client: mockClient, mockAxios };
    }
  }

  /**
   * Setup mock API client with default responses
   */
  setupMockAPIClient(mockAxios) {
    // Health check
    mockAxios.onGet('/health').reply(200, { status: 'ok' });

    // Task operations
    mockAxios.onGet(/\/tasks\/[^/]+$/).reply(200, {
      success: true,
      data: {
        id: 'test-id',
        baslik: 'Test Task',
        aciklama: 'Test Description',
        durum: 'beklemede',
        oncelik: 'orta',
        olusturma_tarihi: '2025-01-01T00:00:00Z',
        guncelleme_tarihi: '2025-01-01T00:00:00Z'
      }
    });

    mockAxios.onGet('/tasks').reply(200, {
      success: true,
      data: [],
      total: 0
    });

    mockAxios.onPost('/tasks').reply(200, {
      success: true,
      message: '✅ Görev başarıyla oluşturuldu!'
    });

    mockAxios.onPut(/\/tasks\/[^/]+$/).reply(200, {
      success: true,
      message: '✅ Görev güncellendi!'
    });

    mockAxios.onDelete(/\/tasks\/[^/]+$/).reply(200, {
      success: true,
      message: '✅ Görev silindi!'
    });

    mockAxios.onPost('/tasks/subtask').reply(200, {
      success: true,
      message: '✅ Alt görev oluşturuldu!'
    });

    mockAxios.onGet(/\/tasks\/[^/]+\/hierarchy$/).reply(200, {
      success: true,
      data: {
        task: {},
        parents: [],
        children: [],
        stats: { total: 0, completed: 0, in_progress: 0 }
      }
    });

    // Project operations
    mockAxios.onGet(/\/projects\/[^/]+$/).reply(200, {
      success: true,
      data: {
        id: 'test-project-id',
        isim: 'Test Project',
        tanim: 'Test Description'
      }
    });

    mockAxios.onGet('/projects').reply(200, {
      success: true,
      data: [],
      total: 0
    });

    mockAxios.onPost('/projects').reply(200, {
      success: true,
      message: '✅ Proje başarıyla oluşturuldu!'
    });

    mockAxios.onGet(/\/projects\/[^/]+\/tasks$/).reply(200, {
      success: true,
      data: [],
      total: 0
    });

    mockAxios.onGet('/projects/active').reply(200, {
      success: true,
      data: null
    });

    mockAxios.onPost('/projects/active').reply(200, {
      success: true,
      message: '✅ Aktif proje ayarlandı!'
    });

    // Template operations
    mockAxios.onGet(/\/templates\/[^/]+$/).reply(200, {
      success: true,
      data: {
        id: 'test-template-id',
        isim: 'Test Template',
        kategori: 'Test'
      }
    });

    mockAxios.onGet('/templates').reply(200, {
      success: true,
      data: [],
      total: 0
    });

    mockAxios.onPost('/tasks/from-template').reply(200, {
      success: true,
      message: '✅ Template kullanılarak görev oluşturuldu!'
    });

    // Summary
    mockAxios.onGet('/summary').reply(200, {
      success: true,
      data: {
        total_tasks: 0,
        total_projects: 0,
        status_counts: {},
        priority_counts: {}
      }
    });

    // Export/Import operations
    mockAxios.onPost('/export').reply(200, {
      success: true,
      message: '✅ Export completed successfully'
    });

    mockAxios.onPost('/import').reply(200, {
      success: true,
      message: '✅ Import completed successfully'
    });

    // Dependencies
    mockAxios.onPost('/tasks/dependency').reply(200, {
      success: true,
      message: '✅ Bağımlılık eklendi!'
    });

    return mockAxios;
  }

  /**
   * Setup common stubs for tests
   */
  setupCommonStubs() {
    // Stub VS Code window methods
    this.stubs.showInformationMessage = this.sandbox.stub(vscode.window, 'showInformationMessage');
    this.stubs.showErrorMessage = this.sandbox.stub(vscode.window, 'showErrorMessage');
    this.stubs.showWarningMessage = this.sandbox.stub(vscode.window, 'showWarningMessage');
    this.stubs.showInputBox = this.sandbox.stub(vscode.window, 'showInputBox');
    this.stubs.showQuickPick = this.sandbox.stub(vscode.window, 'showQuickPick');
    
    return this.stubs;
  }


  /**
   * Create a test workspace
   */
  async createTestWorkspace() {
    const workspaceFolder = vscode.workspace.workspaceFolders?.[0];
    if (!workspaceFolder) {
      throw new Error('No workspace folder found');
    }
    
    return {
      folder: workspaceFolder,
      config: vscode.workspace.getConfiguration('gorev')
    };
  }

  /**
   * Wait for a condition to be true
   */
  async waitFor(condition, timeout = 5000, interval = 100) {
    const startTime = Date.now();
    
    while (Date.now() - startTime < timeout) {
      if (await condition()) {
        return true;
      }
      await new Promise(resolve => setTimeout(resolve, interval));
    }
    
    throw new Error('Timeout waiting for condition');
  }

  /**
   * Simulate user input for task creation
   */
  setupTaskCreationInputs(stubs, taskData = {}) {
    const {
      title = 'Test Task',
      description = 'Test Description',
      priority = 'orta',
      dueDate = '2025-07-01',
      tags = 'test,mock'
    } = taskData;
    
    stubs.showInputBox
      .onCall(0).resolves(title)
      .onCall(1).resolves(description)
      .onCall(2).resolves(dueDate)
      .onCall(3).resolves(tags);
    
    stubs.showQuickPick
      .onCall(0).resolves({ label: 'Orta', value: priority });
  }

  /**
   * Simulate user input for project creation
   */
  setupProjectCreationInputs(stubs, projectData = {}) {
    const {
      name = 'Test Project',
      description = 'Test Project Description'
    } = projectData;
    
    stubs.showInputBox
      .onCall(0).resolves(name)
      .onCall(1).resolves(description);
  }

  /**
   * Get tree view by ID
   */
  async getTreeView(viewId) {
    // This is a simplified version - in real tests you'd need to access the actual tree view
    const extension = vscode.extensions.getExtension('gorev-team.gorev-vscode');
    if (extension && extension.isActive && extension.exports) {
      switch (viewId) {
        case 'gorevTasks':
          return extension.exports.gorevTreeProvider;
        case 'gorevProjects':
          return extension.exports.projeTreeProvider;
        case 'gorevTemplates':
          return extension.exports.templateTreeProvider;
      }
    }
    return null;
  }

  /**
   * Clean up after tests
   */
  cleanup() {
    this.sandbox.restore();
  }

  /**
   * Create a mock context for commands
   */
  createMockContext(data = {}) {
    return {
      subscriptions: [],
      workspaceState: {
        get: this.sandbox.stub(),
        update: this.sandbox.stub()
      },
      globalState: {
        get: this.sandbox.stub(),
        update: this.sandbox.stub()
      },
      extensionPath: __dirname,
      ...data
    };
  }

  /**
   * Assert that a notification was shown
   */
  assertNotification(type, messagePattern) {
    const stub = this.stubs[`show${type}Message`];
    assert(stub.called, `No ${type} message was shown`);
    
    const message = stub.firstCall.args[0];
    if (messagePattern instanceof RegExp) {
      assert(messagePattern.test(message), `Message "${message}" does not match pattern ${messagePattern}`);
    } else {
      assert(message.includes(messagePattern), `Message "${message}" does not include "${messagePattern}"`);
    }
  }
}

module.exports = TestHelper;