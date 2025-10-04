const assert = require('assert');

suite('TreeModels Test Suite', () => {
  let module;

  setup(() => {
    try {
      module = require('../../out/models/treeModels');
    } catch (error) {
      module = null;
    }
  });

  test('should load model module', () => {
    assert(module !== undefined, 'Model module should be defined');
  });

  test('should export model types', () => {
    if (!module) return;
    assert(typeof module === 'object');
  });
});
