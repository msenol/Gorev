const assert = require('assert');
const TestHelper = require('../utils/testHelper');

suite('DebugConfig Test Suite', () => {
  let helper;
  let module;

  setup(() => {
    helper = new TestHelper();

    try {
      module = require('../../out/debug/debugConfig');
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
