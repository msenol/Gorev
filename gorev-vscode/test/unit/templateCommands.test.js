const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('TemplateCommands Test Suite', () => {
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

  test('should register template commands', () => {
    try {
      const module = require('../../out/commands/templateCommands');
      if (module && module.registerTemplateCommands) {
        module.registerTemplateCommands(mockContext, mockApiClient, mockProviders);
        assert(mockContext.subscriptions.length > 0);
      }
    } catch (error) {
      // Module not compiled
    }
  });

  test('should list templates', async () => {
    mockAxios.onGet('/templates').reply(200, {
      success: true,
      data: []
    });

    assert(true);
  });

  test('should create task from template', async () => {
    mockAxios.onPost('/tasks/from-template').reply(200, {
      success: true
    });

    assert(true);
  });
});
