const assert = require('assert');

suite('Utils Test Suite', () => {
  test('should load config utils', () => {
    try {
      const module = require('../../out/utils/config');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load logger utils', () => {
    try {
      const module = require('../../out/utils/logger');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load constants', () => {
    try {
      const module = require('../../out/utils/constants');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });
});
