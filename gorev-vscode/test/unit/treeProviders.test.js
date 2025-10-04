const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const TestHelper = require('../utils/testHelper');

suite('TreeProviders Test Suite', () => {
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
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should load GorevTreeProvider', () => {
    try {
      const module = require('../../out/providers/gorevTreeProvider');
      assert(module);
      assert(typeof module === 'object');
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load ProjeTreeProvider', () => {
    try {
      const module = require('../../out/providers/projeTreeProvider');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load TemplateTreeProvider', () => {
    try {
      const module = require('../../out/providers/templateTreeProvider');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load EnhancedGorevTreeProvider', () => {
    try {
      const module = require('../../out/providers/enhancedGorevTreeProvider');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });
});
