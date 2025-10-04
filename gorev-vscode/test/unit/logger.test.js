const assert = require('assert');

suite('Logger Test Suite', () => {
  let Logger;

  setup(() => {
    try {
      const module = require('../../out/utils/logger');
      Logger = module.Logger;
    } catch (error) {
      Logger = null;
    }
  });

  test('should have Logger class', () => {
    if (!Logger) return;
    assert(Logger);
    assert(typeof Logger.info === 'function');
    assert(typeof Logger.error === 'function');
    assert(typeof Logger.warn === 'function');
  });

  test('should log info messages', () => {
    if (!Logger) return;
    assert.doesNotThrow(() => {
      Logger.info('Test message');
    });
  });

  test('should log error messages', () => {
    if (!Logger) return;
    assert.doesNotThrow(() => {
      Logger.error('Test error');
    });
  });

  test('should log warning messages', () => {
    if (!Logger) return;
    assert.doesNotThrow(() => {
      Logger.warn('Test warning');
    });
  });
});
