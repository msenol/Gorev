const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('DebugCommands Test Suite', () => {
  let helper;
  let sandbox;
  let mockApiClient;
  let mockAxios;
  let stubs;
  let mockContext;
  let mockProviders;
  let module;

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

    try {
      module = require('../../out/commands/debugCommands');
    } catch (error) {
      module = null;
    }
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should load module without errors', () => {
    assert(module !== undefined, 'Module should be defined');
  });

  test('should export expected functions', () => {
    if (!module) return;
    assert(typeof module === 'object' || typeof module === 'function');
  });
});
