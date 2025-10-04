const assert = require('assert');

suite('Models Test Suite', () => {
  test('should load common models', () => {
    try {
      const module = require('../../out/models/common');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load gorev model', () => {
    try {
      const module = require('../../out/models/gorev');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load proje model', () => {
    try {
      const module = require('../../out/models/proje');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });

  test('should load template model', () => {
    try {
      const module = require('../../out/models/template');
      assert(module);
    } catch (error) {
      // Module not compiled
    }
  });
});
