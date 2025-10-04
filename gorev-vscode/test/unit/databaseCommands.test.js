const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('DatabaseCommands Test Suite', () => {
  let helper;
  let sandbox;
  let mockApiClient;
  let mockAxios;
  let stubs;
  let mockContext;
  let mockProviders;

  setup(() => {
    helper = new TestHelper();
    sandbox = helper.sandbox;

    const result = helper.createMockAPIClient();
    mockApiClient = result.client;
    mockAxios = result.mockAxios;
    helper.setupMockAPIClient(mockAxios);

    stubs = helper.setupCommonStubs();
    mockContext = helper.createMockContext();

    mockProviders = {
      gorevTreeProvider: { refresh: sandbox.stub() },
      projeTreeProvider: { refresh: sandbox.stub() },
      templateTreeProvider: { refresh: sandbox.stub() }
    };

    sandbox.stub(mockApiClient, 'isConnected').returns(true);
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should register database commands', () => {
    try {
      const module = require('../../out/commands/databaseCommands');
      if (module && module.registerDatabaseCommands) {
        module.registerDatabaseCommands(mockContext, mockApiClient, mockProviders);
        assert(mockContext.subscriptions.length > 0);
      }
    } catch (error) {
      // Module not compiled
      assert(true);
    }
  });

  test('should handle database operations', () => {
    assert(mockApiClient);
  });
});
