const vscode = require('vscode');
const sinon = require('sinon');
const { mockMCPResponses } = require('../fixtures/mockData');
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
    const clientModule = require('../../out/api/client');
    const ApiClient = clientModule.ApiClient;

    const client = new ApiClient('http://localhost:5082');
    const mockAxios = new MockAdapter(client.axiosInstance);

    return { client, mockAxios };
  }

  /**
   * Setup mock API client with default responses
   */
  setupMockAPIClient(mockAxios) {
    // Setup default successful responses
    mockAxios.onGet('/health').reply(200, { status: 'ok' });

    mockAxios.onGet('/tasks').reply(200, {
      success: true,
      data: [],
      total: 0
    });

    mockAxios.onGet('/projects').reply(200, {
      success: true,
      data: [],
      total: 0
    });

    mockAxios.onGet('/templates').reply(200, {
      success: true,
      data: [],
      total: 0
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
   * Create a mock MCP client
   */
  createMockMCPClient() {
    return {
      isConnected: this.sandbox.stub().returns(true),
      connect: this.sandbox.stub().resolves(),
      disconnect: this.sandbox.stub(),
      callTool: this.sandbox.stub(),
      _onDidConnect: new vscode.EventEmitter(),
      _onDidDisconnect: new vscode.EventEmitter(),
      onDidConnect: function() { return this._onDidConnect.event; },
      onDidDisconnect: function() { return this._onDidDisconnect.event; }
    };
  }

  /**
   * Setup mock MCP client with responses
   */
  setupMockMCPClient(client, responses = mockMCPResponses) {
    // Setup default responses
    client.callTool.withArgs('gorev_listele').resolves(responses.gorev_listele);
    client.callTool.withArgs('proje_listele').resolves(responses.proje_listele);
    client.callTool.withArgs('template_listele').resolves(responses.template_listele);
    client.callTool.withArgs('ozet_goster').resolves(responses.ozet_goster);
    
    // Setup success responses for create operations
    client.callTool.withArgs('gorev_olustur').resolves({
      content: [{ type: 'text', text: '✅ Görev başarıyla oluşturuldu!' }]
    });
    
    client.callTool.withArgs('proje_olustur').resolves({
      content: [{ type: 'text', text: '✅ Proje başarıyla oluşturuldu!' }]
    });
    
    return client;
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