const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const fs = require('fs');
const path = require('path');

suite('DebugConfig Test Suite', () => {
  let sandbox;
  let mockConfig;
  let mockOutputChannel;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    mockConfig = {
      get: sandbox.stub()
    };
    sandbox.stub(vscode.workspace, 'getConfiguration').returns(mockConfig);
    sandbox.stub(vscode.commands, 'executeCommand');
    sandbox.stub(vscode.Uri, 'file');

    // Mock output channel
    mockOutputChannel = {
      appendLine: sandbox.stub()
    };

    // Mock fs operations
    sandbox.stub(fs, 'existsSync');
    sandbox.stub(fs, 'readdirSync');

    // Mock console
    sandbox.stub(console, 'log');
    sandbox.stub(console, 'warn');
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('getDebugConfig', () => {
    let getDebugConfig;

    setup(() => {
      // Clear module cache and re-require
      delete require.cache[require.resolve('../../src/debug/debugConfig')];
      const debugConfig = require('../../src/debug/debugConfig');
      getDebugConfig = debugConfig.getDebugConfig;
    });

    test('should return default debug config', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(false);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/tmp/gorev-debug');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(5000);
      
      const config = getDebugConfig();
      
      assert.strictEqual(config.useDebugWrapper, false);
      assert.strictEqual(config.debugLogPath, '/tmp/gorev-debug');
      assert.strictEqual(config.serverTimeout, 5000);
    });

    test('should return custom debug config', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/custom/debug/path');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(10000);
      
      const config = getDebugConfig();
      
      assert.strictEqual(config.useDebugWrapper, true);
      assert.strictEqual(config.debugLogPath, '/custom/debug/path');
      assert.strictEqual(config.serverTimeout, 10000);
    });

    test('should call vscode.workspace.getConfiguration with gorev', () => {
      getDebugConfig();
      
      assert(vscode.workspace.getConfiguration.calledWith('gorev'));
    });

    test('should handle missing config values', () => {
      mockConfig.get.returns(undefined);
      
      const config = getDebugConfig();
      
      assert.strictEqual(config.useDebugWrapper, false);
      assert.strictEqual(config.debugLogPath, '/tmp/gorev-debug');
      assert.strictEqual(config.serverTimeout, 5000);
    });
  });

  suite('getServerPath', () => {
    let getServerPath, getDebugConfig;

    setup(() => {
      delete require.cache[require.resolve('../../src/debug/debugConfig')];
      const debugConfig = require('../../src/debug/debugConfig');
      getServerPath = debugConfig.getServerPath;
      getDebugConfig = debugConfig.getDebugConfig;
    });

    test('should return normal server path when debug wrapper disabled', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(false);
      mockConfig.get.withArgs('mcp.serverPath').returns('/usr/local/bin/gorev');
      
      const serverPath = getServerPath();
      
      assert.strictEqual(serverPath, '/usr/local/bin/gorev');
    });

    test('should return debug wrapper path when enabled and exists', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('mcp.serverPath').returns('/usr/local/bin/gorev');
      
      const expectedWrapperPath = path.join('/usr/local/bin', '..', 'debug-wrapper.sh');
      fs.existsSync.withArgs(expectedWrapperPath).returns(true);
      
      const serverPath = getServerPath();
      
      assert.strictEqual(serverPath, expectedWrapperPath);
      assert(console.log.calledWith(`[Gorev] Using debug wrapper: ${expectedWrapperPath}`));
    });

    test('should fallback to normal server when wrapper not found', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('mcp.serverPath').returns('/usr/local/bin/gorev');
      
      const expectedWrapperPath = path.join('/usr/local/bin', '..', 'debug-wrapper.sh');
      fs.existsSync.withArgs(expectedWrapperPath).returns(false);
      
      const serverPath = getServerPath();
      
      assert.strictEqual(serverPath, '/usr/local/bin/gorev');
      assert(console.warn.calledWith(`[Gorev] Debug wrapper not found at: ${expectedWrapperPath}`));
    });

    test('should handle missing server path config', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(false);
      mockConfig.get.withArgs('mcp.serverPath').returns(undefined);
      
      const serverPath = getServerPath();
      
      assert.strictEqual(serverPath, '');
    });

    test('should handle empty server path config', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(false);
      mockConfig.get.withArgs('mcp.serverPath').returns('');
      
      const serverPath = getServerPath();
      
      assert.strictEqual(serverPath, '');
    });

    test('should construct wrapper path correctly', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('mcp.serverPath').returns('/home/user/gorev/bin/gorev');
      
      const expectedWrapperPath = path.join('/home/user/gorev/bin', '..', 'debug-wrapper.sh');
      fs.existsSync.withArgs(expectedWrapperPath).returns(true);
      
      getServerPath();
      
      assert(fs.existsSync.calledWith(expectedWrapperPath));
    });
  });

  suite('showDebugInfo', () => {
    let showDebugInfo;

    setup(() => {
      delete require.cache[require.resolve('../../src/debug/debugConfig')];
      const debugConfig = require('../../src/debug/debugConfig');
      showDebugInfo = debugConfig.showDebugInfo;
    });

    test('should show debug info when wrapper enabled', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/custom/debug');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(8000);
      
      showDebugInfo(mockOutputChannel);
      
      assert(mockOutputChannel.appendLine.calledWith('=== Debug Mode Enabled ==='));
      assert(mockOutputChannel.appendLine.calledWith('Debug logs will be written to: /custom/debug'));
      assert(mockOutputChannel.appendLine.calledWith('Server timeout: 8000ms'));
    });

    test('should not show debug info when wrapper disabled', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(false);
      
      showDebugInfo(mockOutputChannel);
      
      assert(!mockOutputChannel.appendLine.called);
    });

    test('should show latest debug log when available', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/debug/path');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(5000);
      
      const mockFiles = ['mcp-session-20240101.log', 'mcp-session-20240102.log'];
      fs.readdirSync.returns(mockFiles);
      
      const latestLogPath = path.join('/debug/path', 'mcp-session-20240102.log');
      const expectedUri = { file: latestLogPath };
      vscode.Uri.file.returns(expectedUri);
      
      showDebugInfo(mockOutputChannel);
      
      assert(mockOutputChannel.appendLine.calledWith(`Latest debug log: ${latestLogPath}`));
      assert(vscode.commands.executeCommand.calledWith('vscode.open', expectedUri));
    });

    test('should handle no debug logs available', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/debug/path');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(5000);
      
      fs.readdirSync.returns([]);
      
      showDebugInfo(mockOutputChannel);
      
      assert(mockOutputChannel.appendLine.calledWith('=== Debug Mode Enabled ==='));
      assert(!vscode.commands.executeCommand.called);
    });

    test('should handle directory read error', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/debug/path');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(5000);
      
      fs.readdirSync.throws(new Error('Directory not found'));
      
      showDebugInfo(mockOutputChannel);
      
      // Should continue without error
      assert(mockOutputChannel.appendLine.calledWith('=== Debug Mode Enabled ==='));
    });

    test('should filter debug log files correctly', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/debug/path');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(5000);
      
      const allFiles = [
        'mcp-session-old.log',
        'other-file.txt',
        'mcp-session-new.log',
        'stdin-123.log'
      ];
      fs.readdirSync.returns(allFiles);
      
      showDebugInfo(mockOutputChannel);
      
      // Should pick the last mcp-session file after sorting
      const latestLogPath = path.join('/debug/path', 'mcp-session-old.log');
      assert(mockOutputChannel.appendLine.calledWith(`Latest debug log: ${latestLogPath}`));
    });

    test('should sort log files and pick latest', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/debug/path');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(5000);
      
      const files = [
        'mcp-session-20240101.log',
        'mcp-session-20240103.log', 
        'mcp-session-20240102.log'
      ];
      fs.readdirSync.returns(files);
      
      showDebugInfo(mockOutputChannel);
      
      // Should pick the latest after reverse sorting
      const latestLogPath = path.join('/debug/path', 'mcp-session-20240103.log');
      assert(mockOutputChannel.appendLine.calledWith(`Latest debug log: ${latestLogPath}`));
    });
  });

  suite('Edge Cases and Error Handling', () => {
    test('should handle null config', () => {
      vscode.workspace.getConfiguration.returns(null);
      
      const { getDebugConfig } = require('../../src/debug/debugConfig');
      
      assert.throws(() => getDebugConfig());
    });

    test('should handle undefined output channel in showDebugInfo', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      
      const { showDebugInfo } = require('../../src/debug/debugConfig');
      
      assert.doesNotThrow(() => {
        showDebugInfo(undefined);
      });
    });

    test('should handle fs.readdirSync permission errors gracefully', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/restricted/path');
      
      fs.readdirSync.throws(new Error('EACCES: permission denied'));
      
      const { showDebugInfo } = require('../../src/debug/debugConfig');
      
      assert.doesNotThrow(() => {
        showDebugInfo(mockOutputChannel);
      });
    });

    test('should handle malformed server path in getServerPath', () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('mcp.serverPath').returns(null);
      
      const { getServerPath } = require('../../src/debug/debugConfig');
      
      assert.doesNotThrow(() => {
        const result = getServerPath();
        assert.strictEqual(result, '');
      });
    });
  });

  suite('Integration Tests', () => {
    test('should work with realistic debug workflow', () => {
      // Setup realistic config
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      mockConfig.get.withArgs('debug.logPath', '/tmp/gorev-debug').returns('/home/user/.gorev/debug');
      mockConfig.get.withArgs('debug.serverTimeout', 5000).returns(10000);
      mockConfig.get.withArgs('mcp.serverPath').returns('/usr/local/bin/gorev');
      
      // Setup file system
      const wrapperPath = path.join('/usr/local/bin', '..', 'debug-wrapper.sh');
      fs.existsSync.withArgs(wrapperPath).returns(true);
      fs.readdirSync.returns(['mcp-session-latest.log']);
      
      const { getDebugConfig, getServerPath, showDebugInfo } = require('../../src/debug/debugConfig');
      
      // Test full workflow
      const config = getDebugConfig();
      const serverPath = getServerPath();
      showDebugInfo(mockOutputChannel);
      
      // Verify results
      assert.strictEqual(config.useDebugWrapper, true);
      assert.strictEqual(config.debugLogPath, '/home/user/.gorev/debug');
      assert.strictEqual(config.serverTimeout, 10000);
      assert.strictEqual(serverPath, wrapperPath);
      
      assert(mockOutputChannel.appendLine.calledWith('=== Debug Mode Enabled ==='));
      assert(console.log.calledWith(`[Gorev] Using debug wrapper: ${wrapperPath}`));
    });
  });
});