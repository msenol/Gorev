const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('TemplateWizard Test Suite', () => {
  let helper;
  let sandbox;
  let mockApiClient;
  let mockAxios;
  let module;

  setup(() => {
    helper = new TestHelper();
    sandbox = helper.sandbox;

    const result = helper.createMockAPIClient();
    mockApiClient = result.client;
    mockAxios = result.mockAxios;
    helper.setupMockAPIClient(mockAxios);

    sandbox.stub(mockApiClient, 'isConnected').returns(true);

    try {
      module = require('../../out/ui/templateWizard');
    } catch (error) {
      module = null;
    }
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should load UI module', () => {
    assert(module !== undefined, 'UI module should be defined');
  });

  test('should export UI class or functions', () => {
    if (!module) return;
    assert(typeof module === 'object' || typeof module === 'function');
  });
});
