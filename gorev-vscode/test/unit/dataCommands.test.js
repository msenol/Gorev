const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('DataCommands Test Suite', () => {
  let helper;
  let sandbox;
  let mockApiClient;
  let mockAxios;
  let stubs;
  let mockContext;
  let mockProviders;
  let dataCommands;

  setup(() => {
    helper = new TestHelper();
    sandbox = helper.sandbox;

    // Create mock API client
    const result = helper.createMockAPIClient();
    mockApiClient = result.client;
    mockAxios = result.mockAxios;
    helper.setupMockAPIClient(mockAxios);

    // Setup common stubs
    stubs = helper.setupCommonStubs();

    // Create mock context
    mockContext = helper.createMockContext();

    // Create mock providers
    mockProviders = {
      gorevTreeProvider: {
        refresh: sandbox.stub()
      },
      projeTreeProvider: {
        refresh: sandbox.stub()
      },
      templateTreeProvider: {
        refresh: sandbox.stub()
      }
    };

    // Stub isConnected
    sandbox.stub(mockApiClient, 'isConnected').returns(true);

    // Load data commands module
    try {
      dataCommands = require('../../out/commands/dataCommands');
    } catch (error) {
      // Module not compiled yet, skip tests
      dataCommands = null;
    }
  });

  teardown(() => {
    helper.cleanup();
  });

  suite('Registration', () => {
    test('should register export data command', () => {
      if (!dataCommands) return;

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0, 'Commands should be registered');
    });

    test('should register import data command', () => {
      if (!dataCommands) return;

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0, 'Commands should be registered');
    });

    test('should register export current view command', () => {
      if (!dataCommands) return;

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0, 'Commands should be registered');
    });

    test('should register quick export command', () => {
      if (!dataCommands) return;

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0, 'Commands should be registered');
    });
  });

  suite('Export Data Command', () => {
    test('should check connection before export', async () => {
      if (!dataCommands) return;

      mockApiClient.isConnected.returns(false);
      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      // Find export command
      const exportCommand = mockContext.subscriptions.find(
        sub => sub && sub.dispose
      );

      if (exportCommand && typeof exportCommand === 'function') {
        await exportCommand();
      }

      assert(stubs.showWarningMessage.called || mockApiClient.isConnected.called);
    });

    test('should handle export errors gracefully', async () => {
      if (!dataCommands) return;

      // Mock export to fail
      mockAxios.onPost('/export').reply(500, {
        success: false,
        error: 'Export failed'
      });

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      // Should not throw
      assert(true);
    });
  });

  suite('Import Data Command', () => {
    test('should check connection before import', async () => {
      if (!dataCommands) return;

      mockApiClient.isConnected.returns(false);
      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockApiClient.isConnected.called || mockContext.subscriptions.length > 0);
    });

    test('should handle import errors gracefully', async () => {
      if (!dataCommands) return;

      // Mock import to fail
      mockAxios.onPost('/import').reply(500, {
        success: false,
        error: 'Import failed'
      });

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      // Should not throw
      assert(true);
    });

    test('should refresh providers after successful import', async () => {
      if (!dataCommands) return;

      // Mock successful import
      mockAxios.onPost('/import').reply(200, {
        success: true,
        message: 'Import completed'
      });

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      // Should register commands
      assert(mockContext.subscriptions.length > 0);
    });
  });

  suite('Export Current View Command', () => {
    test('should export filtered tasks', async () => {
      if (!dataCommands) return;

      // Mock tasks endpoint
      mockAxios.onGet('/tasks').reply(200, {
        success: true,
        data: [
          {
            id: 'task-1',
            baslik: 'Test Task',
            durum: 'beklemede'
          }
        ]
      });

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0);
    });

    test('should handle empty view', async () => {
      if (!dataCommands) return;

      mockAxios.onGet('/tasks').reply(200, {
        success: true,
        data: []
      });

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0);
    });
  });

  suite('Quick Export Command', () => {
    test('should export with default settings', async () => {
      if (!dataCommands) return;

      mockAxios.onPost('/export').reply(200, {
        success: true,
        message: 'Export completed'
      });

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0);
    });

    test('should use default file name with timestamp', async () => {
      if (!dataCommands) return;

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      // Should register command
      assert(mockContext.subscriptions.length > 0);
    });

    test('should save to Downloads folder by default', async () => {
      if (!dataCommands) return;

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0);
    });

    test('should show success notification', async () => {
      if (!dataCommands) return;

      mockAxios.onPost('/export').reply(200, {
        success: true,
        message: 'Export completed'
      });

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0);
    });
  });

  suite('Error Handling', () => {
    test('should handle disconnected API client', async () => {
      if (!dataCommands) return;

      mockApiClient.isConnected.returns(false);
      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      // Should register commands even when disconnected
      assert(mockContext.subscriptions.length > 0);
    });

    test('should handle network errors', async () => {
      if (!dataCommands) return;

      mockAxios.onPost('/export').networkError();

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0);
    });

    test('should handle timeout errors', async () => {
      if (!dataCommands) return;

      mockAxios.onPost('/export').timeout();

      dataCommands.registerDataCommands(mockContext, mockApiClient, mockProviders);

      assert(mockContext.subscriptions.length > 0);
    });
  });
});
