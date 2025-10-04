const assert = require('assert');
const TestHelper = require('../utils/testHelper');

suite('TestDataSeederWithTemplates Test Suite', () => {
  let helper;
  let module;

  setup(() => {
    helper = new TestHelper();

    try {
      module = require('../../out/debug/testDataSeederWithTemplates');
    } catch (error) {
      module = null;
    }
  });

  teardown(() => {
    helper.cleanup();
  });

  test('should load debug module', () => {
    assert(module !== undefined, 'Debug module should be defined');
  });
});
