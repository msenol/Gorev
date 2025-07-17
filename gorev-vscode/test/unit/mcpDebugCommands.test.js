const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const fs = require('fs');
const path = require('path');

// Mock child_process before requiring the module
const mockChildProcess = {
  spawn: sinon.stub()
};
require.cache[require.resolve('child_process')] = {
  exports: mockChildProcess
};

suite('MCPDebugCommands Test Suite', () => {
  let sandbox;
  let mockContext;
  let mockOutputChannel;
  let mockConfig;
  let registerFunction;
  let mockDebugConfig;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showWarningMessage');
    sandbox.stub(vscode.window, 'showQuickPick');
    sandbox.stub(vscode.commands, 'registerCommand');
    sandbox.stub(vscode.commands, 'executeCommand');
    sandbox.stub(vscode.Uri, 'file');

    // Mock workspace configuration
    mockConfig = {
      get: sandbox.stub(),
      update: sandbox.stub().resolves()
    };
    sandbox.stub(vscode.workspace, 'getConfiguration').returns(mockConfig);

    // Mock output channel
    mockOutputChannel = {
      show: sandbox.stub(),
      appendLine: sandbox.stub()
    };

    // Mock context
    mockContext = {
      subscriptions: []
    };

    // Mock fs operations
    sandbox.stub(fs, 'readdirSync');
    sandbox.stub(fs, 'existsSync');
    sandbox.stub(fs, 'accessSync');
    sandbox.stub(fs, 'unlinkSync');
    sandbox.stub(fs, 'statSync');

    // Mock debug config
    mockDebugConfig = {
      debugLogPath: '/tmp/gorev-debug',
      useDebugWrapper: false,
      serverTimeout: 5000
    };

    // Mock getDebugConfig and showDebugInfo
    const debugConfig = require('../../src/debug/debugConfig');
    sandbox.stub(debugConfig, 'getDebugConfig').returns(mockDebugConfig);
    sandbox.stub(debugConfig, 'showDebugInfo');

    // Reset child_process mock
    mockChildProcess.spawn.reset();

    // Import function under test
    registerFunction = require('../../src/commands/mcpDebugCommands').registerMCPDebugCommands;
  });

  teardown(() => {
    sandbox.restore();
  });

  test('registerMCPDebugCommands should register all commands', () => {
    registerFunction(mockContext, mockOutputChannel);

    assert.strictEqual(vscode.commands.registerCommand.callCount, 4);
    assert(vscode.commands.registerCommand.calledWith('gorev.toggleDebugMode'));
    assert(vscode.commands.registerCommand.calledWith('gorev.showDebugLogs'));
    assert(vscode.commands.registerCommand.calledWith('gorev.clearDebugLogs'));
    assert(vscode.commands.registerCommand.calledWith('gorev.testConnection'));
  });

  suite('Toggle Debug Mode Command', () => {
    let toggleCallback;

    setup(() => {
      registerFunction(mockContext, mockOutputChannel);
      toggleCallback = vscode.commands.registerCommand.getCall(0).args[1];
    });

    test('should toggle debug mode from false to true', async () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(false);
      
      await toggleCallback();

      assert(mockConfig.update.calledWith('debug.useWrapper', true, vscode.ConfigurationTarget.Workspace));
      assert(vscode.window.showInformationMessage.calledWith(
        'Gorev debug mode enabled. Please restart VS Code for changes to take effect.'
      ));
    });

    test('should toggle debug mode from true to false', async () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      
      await toggleCallback();

      assert(mockConfig.update.calledWith('debug.useWrapper', false, vscode.ConfigurationTarget.Workspace));
      assert(vscode.window.showInformationMessage.calledWith(
        'Gorev debug mode disabled. Please restart VS Code for changes to take effect.'
      ));
    });

    test('should show debug info when enabling', async () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(false);
      const debugConfig = require('../../src/debug/debugConfig');
      
      await toggleCallback();

      assert(debugConfig.showDebugInfo.calledWith(mockOutputChannel));
    });

    test('should not show debug info when disabling', async () => {
      mockConfig.get.withArgs('debug.useWrapper', false).returns(true);
      const debugConfig = require('../../src/debug/debugConfig');
      
      await toggleCallback();

      assert(!debugConfig.showDebugInfo.called);
    });

    test('should handle config update error', async () => {
      const error = new Error('Config update failed');
      mockConfig.update.rejects(error);
      
      await toggleCallback();
      
      // Should not throw, but handle gracefully
      assert(mockConfig.update.called);
    });
  });

  suite('Show Debug Logs Command', () => {
    let showLogsCallback;

    setup(() => {
      registerFunction(mockContext, mockOutputChannel);
      showLogsCallback = vscode.commands.registerCommand.getCall(1).args[1];
    });

    test('should show debug logs when files exist', async () => {
      const mockFiles = ['mcp-session-123.log', 'stdin-456.log', 'stdout-789.log'];
      fs.readdirSync.returns(mockFiles);
      fs.statSync.returns({ size: 1024 });
      
      const selectedFile = { label: 'mcp-session-123.log' };
      vscode.window.showQuickPick.resolves(selectedFile);
      
      await showLogsCallback();

      assert(fs.readdirSync.calledWith('/tmp/gorev-debug'));
      assert(vscode.window.showQuickPick.called);
      assert(vscode.commands.executeCommand.calledWith('vscode.open'));
    });

    test('should show message when no debug logs found', async () => {
      fs.readdirSync.returns([]);
      
      await showLogsCallback();

      assert(vscode.window.showInformationMessage.calledWith(
        'No debug logs found. Enable debug mode and restart VS Code.'
      ));
    });

    test('should filter log files correctly', async () => {
      const allFiles = [
        'mcp-session-123.log',
        'other-file.txt',
        'stdin-456.log',
        'stdout-789.log',
        'another-file.json'
      ];
      fs.readdirSync.returns(allFiles);
      fs.statSync.returns({ size: 1024 });
      
      await showLogsCallback();

      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      assert.strictEqual(items.length, 3);
      assert(items.some(item => item.label === 'mcp-session-123.log'));
      assert(items.some(item => item.label === 'stdin-456.log'));
      assert(items.some(item => item.label === 'stdout-789.log'));
    });

    test('should sort files in reverse order (newest first)', async () => {
      const mockFiles = ['file-1.log', 'file-2.log', 'file-3.log'];
      fs.readdirSync.returns(mockFiles);
      fs.statSync.returns({ size: 1024 });
      
      await showLogsCallback();

      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      // Files should be in reverse order
      assert.strictEqual(items[0].label, 'file-3.log');
      assert.strictEqual(items[1].label, 'file-2.log');
      assert.strictEqual(items[2].label, 'file-1.log');
    });

    test('should handle directory read error', async () => {
      const error = new Error('Permission denied');
      fs.readdirSync.throws(error);
      
      await showLogsCallback();

      assert(vscode.window.showErrorMessage.calledWith('Failed to read debug logs: Error: Permission denied'));
    });

    test('should handle user cancellation', async () => {
      fs.readdirSync.returns(['mcp-session-123.log']);
      fs.statSync.returns({ size: 1024 });
      vscode.window.showQuickPick.resolves(undefined);
      
      await showLogsCallback();

      assert(!vscode.commands.executeCommand.calledWith('vscode.open'));
    });

    test('should include file descriptions', async () => {
      fs.readdirSync.returns(['mcp-session-123.log', 'stdin-456.log', 'stdout-789.log']);
      fs.statSync.returns({ size: 1024 });
      
      await showLogsCallback();

      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      const mcpItem = items.find(item => item.label === 'mcp-session-123.log');
      const stdinItem = items.find(item => item.label === 'stdin-456.log');
      const stdoutItem = items.find(item => item.label === 'stdout-789.log');
      
      assert.strictEqual(mcpItem.description, 'Main debug log');
      assert.strictEqual(stdinItem.description, 'Input messages (VS Code → Server)');
      assert.strictEqual(stdoutItem.description, 'Output messages (Server → VS Code)');
    });

    test('should include file sizes', async () => {
      fs.readdirSync.returns(['test.log']);
      fs.statSync.returns({ size: 1536 }); // 1.5 KB
      
      await showLogsCallback();

      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      assert.strictEqual(items[0].detail, '1.5 KB');
    });
  });

  suite('Clear Debug Logs Command', () => {
    let clearLogsCallback;

    setup(() => {
      registerFunction(mockContext, mockOutputChannel);
      clearLogsCallback = vscode.commands.registerCommand.getCall(2).args[1];
    });

    test('should clear debug logs when user confirms', async () => {
      const mockFiles = ['mcp-session-123.log', 'stdin-456.log', 'debug.log'];
      fs.readdirSync.returns(mockFiles);
      vscode.window.showWarningMessage.resolves('Yes');
      
      await clearLogsCallback();

      assert(vscode.window.showWarningMessage.calledWith(
        'Are you sure you want to clear all debug logs?',
        'Yes',
        'No'
      ));
      assert.strictEqual(fs.unlinkSync.callCount, 3);
      assert(vscode.window.showInformationMessage.calledWith('Cleared 3 debug log files.'));
    });

    test('should not clear logs when user cancels', async () => {
      vscode.window.showWarningMessage.resolves('No');
      
      await clearLogsCallback();

      assert(!fs.readdirSync.called);
      assert(!fs.unlinkSync.called);
    });

    test('should handle no logs to clear', async () => {
      fs.readdirSync.returns([]);
      vscode.window.showWarningMessage.resolves('Yes');
      
      await clearLogsCallback();

      assert(!fs.unlinkSync.called);
      assert(vscode.window.showInformationMessage.calledWith('Cleared 0 debug log files.'));
    });

    test('should handle clear error', async () => {
      const error = new Error('Permission denied');
      fs.readdirSync.throws(error);
      vscode.window.showWarningMessage.resolves('Yes');
      
      await clearLogsCallback();

      assert(vscode.window.showErrorMessage.calledWith('Failed to clear debug logs: Error: Permission denied'));
    });

    test('should filter log files for clearing', async () => {
      const allFiles = [
        'mcp-session-123.log',
        'other-file.txt',
        'stdin-456.log',
        'stdout-789.log',
        'debug.log',
        'another-file.json'
      ];
      fs.readdirSync.returns(allFiles);
      vscode.window.showWarningMessage.resolves('Yes');
      
      await clearLogsCallback();

      assert.strictEqual(fs.unlinkSync.callCount, 4); // Only log files
    });
  });

  suite('Test Connection Command', () => {
    let testConnectionCallback;
    let mockProcess;

    setup(() => {
      registerFunction(mockContext, mockOutputChannel);
      testConnectionCallback = vscode.commands.registerCommand.getCall(3).args[1];
      
      mockProcess = {
        stdout: { on: sandbox.stub() },
        stderr: { on: sandbox.stub() },
        stdin: { write: sandbox.stub() },
        kill: sandbox.stub()
      };
      mockChildProcess.spawn.returns(mockProcess);
    });

    test('should test connection with valid server path', async () => {
      const serverPath = '/usr/local/bin/gorev';
      mockConfig.get.withArgs('mcp.serverPath').returns(serverPath);
      fs.existsSync.withArgs(serverPath).returns(true);
      
      const clock = sandbox.useFakeTimers();
      const promise = testConnectionCallback();
      clock.tick(3000);
      await promise;
      
      assert(mockOutputChannel.show.called);
      assert(mockOutputChannel.appendLine.calledWith('=== Testing MCP Connection ==='));
      assert(mockOutputChannel.appendLine.calledWith(`Server path: ${serverPath}`));
      assert(mockOutputChannel.appendLine.calledWith('✓ Server file is executable'));
      
      clock.restore();
    });

    test('should handle missing server path configuration', async () => {
      mockConfig.get.returns(undefined);
      
      await testConnectionCallback();

      assert(mockOutputChannel.appendLine.calledWith('ERROR: No server path configured'));
      assert(vscode.window.showErrorMessage.calledWith('Please configure the Gorev server path in settings.'));
    });

    test('should handle non-existent server file', async () => {
      const serverPath = '/nonexistent/gorev';
      mockConfig.get.withArgs('mcp.serverPath').returns(serverPath);
      fs.existsSync.withArgs(serverPath).returns(false);
      
      await testConnectionCallback();

      assert(mockOutputChannel.appendLine.calledWith('ERROR: Server file does not exist'));
      assert(vscode.window.showErrorMessage.calledWith(`Server not found at: ${serverPath}`));
    });

    test('should handle non-executable server file', async () => {
      const serverPath = '/usr/local/bin/gorev';
      mockConfig.get.withArgs('mcp.serverPath').returns(serverPath);
      fs.existsSync.withArgs(serverPath).returns(true);
      fs.accessSync.throws(new Error('Permission denied'));
      
      await testConnectionCallback();

      assert(mockOutputChannel.appendLine.calledWith('WARNING: Server file may not be executable'));
    });

    test('should test MCP initialization', async () => {
      const serverPath = '/usr/local/bin/gorev';
      mockConfig.get.withArgs('mcp.serverPath').returns(serverPath);
      fs.existsSync.withArgs(serverPath).returns(true);
      
      const clock = sandbox.useFakeTimers();
      const promise = testConnectionCallback();
      
      // Simulate server response
      const onData = mockProcess.stdout.on.getCall(0).args[1];
      onData(Buffer.from('{"jsonrpc":"2.0","id":1,"result":{}}'));
      
      clock.tick(3000);
      await promise;
      
      assert(mockChildProcess.spawn.calledWith(serverPath, ['serve']));
      assert(mockProcess.stdin.write.called);
      assert(mockOutputChannel.appendLine.calledWith('✓ Server responded to initialize request'));
      
      clock.restore();
    });

    test('should handle server stderr output', async () => {
      const serverPath = '/usr/local/bin/gorev';
      mockConfig.get.withArgs('mcp.serverPath').returns(serverPath);
      fs.existsSync.withArgs(serverPath).returns(true);
      
      const clock = sandbox.useFakeTimers();
      const promise = testConnectionCallback();
      
      // Simulate server error
      const onError = mockProcess.stderr.on.getCall(0).args[1];
      onError(Buffer.from('Server error occurred'));
      
      clock.tick(3000);
      await promise;
      
      assert(mockOutputChannel.appendLine.calledWith('STDERR: Server error occurred'));
      
      clock.restore();
    });

    test('should handle no server response', async () => {
      const serverPath = '/usr/local/bin/gorev';
      mockConfig.get.withArgs('mcp.serverPath').returns(serverPath);
      fs.existsSync.withArgs(serverPath).returns(true);
      
      const clock = sandbox.useFakeTimers();
      const promise = testConnectionCallback();
      clock.tick(3000);
      await promise;
      
      assert(mockOutputChannel.appendLine.calledWith('✗ No response from server'));
      
      clock.restore();
    });

    test('should handle spawn error', async () => {
      const serverPath = '/usr/local/bin/gorev';
      mockConfig.get.withArgs('mcp.serverPath').returns(serverPath);
      fs.existsSync.withArgs(serverPath).returns(true);
      mockChildProcess.spawn.throws(new Error('Spawn failed'));
      
      await testConnectionCallback();

      assert(mockOutputChannel.appendLine.calledWith('ERROR: Error: Spawn failed'));
    });

    test('should use fallback server path', async () => {
      mockConfig.get.withArgs('mcp.serverPath').returns(undefined);
      mockConfig.get.withArgs('serverPath').returns('/fallback/gorev');
      fs.existsSync.returns(true);
      
      await testConnectionCallback();

      assert(mockOutputChannel.appendLine.calledWith('Server path: /fallback/gorev'));
    });
  });

  suite('Helper Functions', () => {
    const { registerMCPDebugCommands } = require('../../src/commands/mcpDebugCommands');
    
    test('getFileDescription should return correct descriptions', () => {
      // We need to test this indirectly through the showDebugLogs command
      fs.readdirSync.returns(['mcp-session-123.log', 'stdin-456.log', 'stdout-789.log', 'other.log']);
      fs.statSync.returns({ size: 1024 });
      
      registerFunction(mockContext, mockOutputChannel);
      const showLogsCallback = vscode.commands.registerCommand.getCall(1).args[1];
      
      showLogsCallback();
      
      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      assert(items.some(item => item.description === 'Main debug log'));
      assert(items.some(item => item.description === 'Input messages (VS Code → Server)'));
      assert(items.some(item => item.description === 'Output messages (Server → VS Code)'));
      assert(items.some(item => item.description === 'Debug log'));
    });

    test('getFileSize should format sizes correctly', () => {
      // Test different file sizes through the showDebugLogs command
      fs.readdirSync.returns(['small.log', 'medium.log', 'large.log']);
      
      // Mock different file sizes
      fs.statSync.onCall(0).returns({ size: 500 }); // bytes
      fs.statSync.onCall(1).returns({ size: 1536 }); // KB
      fs.statSync.onCall(2).returns({ size: 2097152 }); // MB
      
      registerFunction(mockContext, mockOutputChannel);
      const showLogsCallback = vscode.commands.registerCommand.getCall(1).args[1];
      
      showLogsCallback();
      
      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      assert.strictEqual(items[0].detail, '500 bytes');
      assert.strictEqual(items[1].detail, '1.5 KB');
      assert.strictEqual(items[2].detail, '2.0 MB');
    });
  });

  suite('Edge Cases', () => {
    test('should handle undefined output channel', () => {
      assert.doesNotThrow(() => {
        registerFunction(mockContext, undefined);
      });
    });

    test('should handle context without subscriptions', () => {
      const badContext = {};
      
      assert.doesNotThrow(() => {
        registerFunction(badContext, mockOutputChannel);
      });
    });

    test('should handle file stat error', async () => {
      fs.readdirSync.returns(['test.log']);
      fs.statSync.throws(new Error('Stat failed'));
      
      registerFunction(mockContext, mockOutputChannel);
      const showLogsCallback = vscode.commands.registerCommand.getCall(1).args[1];
      
      await showLogsCallback();
      
      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      assert.strictEqual(items[0].detail, '');
    });
  });
});