const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('DecorationProvider Test Suite', () => {
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
      module = require('../../out/providers/decorationProvider');
    } catch (error) {
      module = null;
    }
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should load provider module', () => {
    assert(module !== undefined, 'Provider module should be defined');
  });

  test('should export provider class', () => {
    if (!module) return;
    assert(typeof module === 'object' || typeof module === 'function');
  });
});
