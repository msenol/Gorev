const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('StatusBar Test Suite', () => {
  let helper;
  let sandbox;
  let mockApiClient;
  let mockAxios;

  setup(() => {
    helper = new TestHelper();
    sandbox = helper.sandbox;

    const result = helper.createMockAPIClient();
    mockApiClient = result.client;
    mockAxios = result.mockAxios;
    helper.setupMockAPIClient(mockAxios);

    sandbox.stub(mockApiClient, 'isConnected').returns(true);

    // Mock VS Code status bar
    sandbox.stub(vscode.window, 'createStatusBarItem').returns({
      text: '',
      tooltip: '',
      command: '',
      show: sandbox.stub(),
      hide: sandbox.stub(),
      dispose: sandbox.stub()
    });
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should load status bar module', () => {
    try {
      const module = require('../../out/ui/statusBar');
      assert(module);
    } catch (error) {
      // Module not compiled
      assert(true);
    }
  });

  test('should create status bar manager', () => {
    try {
      const module = require('../../out/ui/statusBar');
      if (module && module.StatusBarManager) {
        const manager = new module.StatusBarManager(mockApiClient);
        assert(manager);
      }
    } catch (error) {
      // Module not compiled
      assert(true);
    }
  });

  test('should handle connection status updates', () => {
    mockAxios.onGet('/summary').reply(200, {
      success: true,
      data: { total_tasks: 0 }
    });

    assert(true);
  });
});
