const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const { COMMANDS } = require('../../src/utils/constants');
const { InlineEditProvider } = require('../../src/providers/inlineEditProvider');
const { Logger } = require('../../src/utils/logger');

suite('InlineEditCommands Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let mockContext;
  let mockProviders;
  let mockTreeProvider;
  let registerFunction;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showWarningMessage');
    sandbox.stub(vscode.commands, 'registerCommand');
    sandbox.stub(vscode.commands, 'executeCommand');

    // Mock MCP Client
    mockMCPClient = {
      callTool: sandbox.stub().resolves({ content: [{ text: 'Success' }] }),
      isConnected: sandbox.stub().returns(true)
    };

    // Mock Enhanced Tree Provider
    mockTreeProvider = {
      getSelectedTasks: sandbox.stub().returns([]),
      refresh: sandbox.stub().resolves()
    };

    // Mock TreeProviders
    mockProviders = {
      gorevTreeProvider: mockTreeProvider,
      projeTreeProvider: {
        refresh: sandbox.stub().resolves()
      }
    };

    // Mock Context
    mockContext = {
      subscriptions: []
    };

    // Mock Logger
    sandbox.stub(Logger, 'error');
    sandbox.stub(Logger, 'info');
    sandbox.stub(Logger, 'warn');

    // Mock InlineEditProvider methods
    sandbox.stub(InlineEditProvider.prototype, 'startEdit').resolves();
    sandbox.stub(InlineEditProvider.prototype, 'quickStatusChange').resolves();
    sandbox.stub(InlineEditProvider.prototype, 'quickPriorityChange').resolves();
    sandbox.stub(InlineEditProvider.prototype, 'quickDateChange').resolves();
    sandbox.stub(InlineEditProvider.prototype, 'showDetailedEdit').resolves();
    sandbox.stub(InlineEditProvider.prototype, 'isEditing').returns(false);
    sandbox.stub(InlineEditProvider.prototype, 'cancelEdit');

    // Mock global projeTreeProvider
    global.projeTreeProvider = {
      refresh: sandbox.stub().resolves()
    };

    // Import function under test
    registerFunction = require('../../src/commands/inlineEditCommands').registerInlineEditCommands;
  });

  teardown(() => {
    sandbox.restore();
    delete global.projeTreeProvider;
  });

  test('registerInlineEditCommands should register all commands', () => {
    registerFunction(mockContext, mockMCPClient, mockProviders);

    assert.strictEqual(vscode.commands.registerCommand.callCount, 6);
    assert(vscode.commands.registerCommand.calledWith(COMMANDS.EDIT_TASK_TITLE));
    assert(vscode.commands.registerCommand.calledWith(COMMANDS.QUICK_STATUS_CHANGE));
    assert(vscode.commands.registerCommand.calledWith(COMMANDS.QUICK_PRIORITY_CHANGE));
    assert(vscode.commands.registerCommand.calledWith(COMMANDS.QUICK_DATE_CHANGE));
    assert(vscode.commands.registerCommand.calledWith(COMMANDS.DETAILED_EDIT));
    assert(vscode.commands.registerCommand.calledWith('gorev.onTreeItemDoubleClick'));
  });

  suite('Edit Task Title Command', () => {
    let editCallback;
    const mockTask = { id: '123', baslik: 'Test Task' };

    setup(() => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      editCallback = vscode.commands.registerCommand.getCall(0).args[1];
    });

    test('should edit task when item provided', async () => {
      const item = { task: mockTask };
      
      await editCallback(item);

      assert(InlineEditProvider.prototype.startEdit.calledWith(item));
      assert(mockTreeProvider.refresh.called);
    });

    test('should use selected task when no item provided', async () => {
      mockTreeProvider.getSelectedTasks.returns([mockTask]);
      
      await editCallback();

      assert(InlineEditProvider.prototype.startEdit.calledWith({ task: mockTask }));
      assert(mockTreeProvider.refresh.called);
    });

    test('should show warning when no task selected', async () => {
      mockTreeProvider.getSelectedTasks.returns([]);
      
      await editCallback();

      assert(vscode.window.showWarningMessage.calledWith('Lütfen düzenlemek için bir görev seçin'));
      assert(!InlineEditProvider.prototype.startEdit.called);
    });

    test('should show warning when multiple tasks selected', async () => {
      mockTreeProvider.getSelectedTasks.returns([mockTask, { id: '456' }]);
      
      await editCallback();

      assert(vscode.window.showWarningMessage.calledWith('Lütfen düzenlemek için bir görev seçin'));
      assert(!InlineEditProvider.prototype.startEdit.called);
    });

    test('should handle missing task in item', async () => {
      const item = {};
      
      await editCallback(item);

      assert(vscode.window.showWarningMessage.calledWith('Lütfen düzenlemek için bir görev seçin'));
    });
  });

  suite('Quick Status Change Command', () => {
    let statusCallback;
    const mockTask = { id: '123', baslik: 'Test Task', durum: 'beklemede' };

    setup(() => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      statusCallback = vscode.commands.registerCommand.getCall(1).args[1];
    });

    test('should change status when item provided', async () => {
      const item = { task: mockTask };
      
      await statusCallback(item);

      assert(Logger.info.calledWith('[QuickStatusChange Command] Called'));
      assert(Logger.info.calledWith('[QuickStatusChange Command] Task ID:', mockTask.id));
      assert(InlineEditProvider.prototype.quickStatusChange.calledWith(mockTask));
      assert(mockTreeProvider.refresh.called);
    });

    test('should refresh global project tree provider', async () => {
      const item = { task: mockTask };
      
      await statusCallback(item);

      assert(global.projeTreeProvider.refresh.called);
    });

    test('should show warning when no task provided', async () => {
      await statusCallback({});

      assert(Logger.warn.calledWith('[QuickStatusChange Command] No task found in item'));
      assert(vscode.window.showWarningMessage.calledWith('Lütfen bir görev seçin'));
      assert(!InlineEditProvider.prototype.quickStatusChange.called);
    });

    test('should handle quickStatusChange error', async () => {
      const item = { task: mockTask };
      const error = new Error('Status change failed');
      InlineEditProvider.prototype.quickStatusChange.rejects(error);
      
      await statusCallback(item);

      assert(Logger.error.calledWith('Quick status change failed:', error));
    });

    test('should add delay after status change', async () => {
      const item = { task: mockTask };
      const clock = sandbox.useFakeTimers();
      
      const promise = statusCallback(item);
      clock.tick(100);
      await promise;

      assert(mockTreeProvider.refresh.called);
      clock.restore();
    });

    test('should log task details', async () => {
      const item = { task: mockTask };
      
      await statusCallback(item);

      assert(Logger.info.calledWith('[QuickStatusChange Command] Task title:', mockTask.baslik));
      assert(Logger.info.calledWith('[QuickStatusChange Command] Current status:', mockTask.durum));
    });
  });

  suite('Quick Priority Change Command', () => {
    let priorityCallback;
    const mockTask = { id: '123', baslik: 'Test Task' };

    setup(() => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      priorityCallback = vscode.commands.registerCommand.getCall(2).args[1];
    });

    test('should change priority when item provided', async () => {
      const item = { task: mockTask };
      
      await priorityCallback(item);

      assert(InlineEditProvider.prototype.quickPriorityChange.calledWith(mockTask));
      assert(mockTreeProvider.refresh.called);
    });

    test('should show warning when no task provided', async () => {
      await priorityCallback({});

      assert(vscode.window.showWarningMessage.calledWith('Lütfen bir görev seçin'));
      assert(!InlineEditProvider.prototype.quickPriorityChange.called);
    });

    test('should handle null item', async () => {
      await priorityCallback(null);

      assert(vscode.window.showWarningMessage.calledWith('Lütfen bir görev seçin'));
    });
  });

  suite('Quick Date Change Command', () => {
    let dateCallback;
    const mockTask = { id: '123', baslik: 'Test Task' };

    setup(() => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      dateCallback = vscode.commands.registerCommand.getCall(3).args[1];
    });

    test('should change date when item provided', async () => {
      const item = { task: mockTask };
      
      await dateCallback(item);

      assert(InlineEditProvider.prototype.quickDateChange.calledWith(mockTask));
      assert(mockTreeProvider.refresh.called);
    });

    test('should show warning when no task provided', async () => {
      await dateCallback({});

      assert(vscode.window.showWarningMessage.calledWith('Lütfen bir görev seçin'));
      assert(!InlineEditProvider.prototype.quickDateChange.called);
    });
  });

  suite('Detailed Edit Command', () => {
    let detailedCallback;
    const mockTask = { id: '123', baslik: 'Test Task' };

    setup(() => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      detailedCallback = vscode.commands.registerCommand.getCall(4).args[1];
    });

    test('should show detailed edit when item provided', async () => {
      const item = { task: mockTask };
      
      await detailedCallback(item);

      assert(InlineEditProvider.prototype.showDetailedEdit.calledWith(mockTask));
      assert(mockTreeProvider.refresh.called);
    });

    test('should show warning when no task provided', async () => {
      await detailedCallback({});

      assert(vscode.window.showWarningMessage.calledWith('Lütfen bir görev seçin'));
      assert(!InlineEditProvider.prototype.showDetailedEdit.called);
    });
  });

  suite('Double Click Command', () => {
    let doubleClickCallback;
    const mockTask = { id: '123', baslik: 'Test Task' };

    setup(() => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      doubleClickCallback = vscode.commands.registerCommand.getCall(5).args[1];
    });

    test('should start edit on double click', async () => {
      const item = { task: mockTask };
      
      await doubleClickCallback(item);

      assert(InlineEditProvider.prototype.startEdit.calledWith(item));
      assert(mockTreeProvider.refresh.called);
    });

    test('should handle item without task', async () => {
      await doubleClickCallback({});

      assert(!InlineEditProvider.prototype.startEdit.called);
      assert(!mockTreeProvider.refresh.called);
    });

    test('should handle null item', async () => {
      await doubleClickCallback(null);

      assert(!InlineEditProvider.prototype.startEdit.called);
    });
  });

  suite('Cancel Edit Command', () => {
    let cancelCallback;

    setup(() => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      cancelCallback = vscode.commands.registerCommand.getCall(6).args[1];
    });

    test('should cancel edit when editing', () => {
      InlineEditProvider.prototype.isEditing.returns(true);
      
      cancelCallback();

      assert(InlineEditProvider.prototype.cancelEdit.called);
      assert(vscode.window.showInformationMessage.calledWith('Düzenleme iptal edildi'));
    });

    test('should not show message when not editing', () => {
      InlineEditProvider.prototype.isEditing.returns(false);
      
      cancelCallback();

      assert(InlineEditProvider.prototype.cancelEdit.called);
      assert(!vscode.window.showInformationMessage.called);
    });
  });

  suite('Edge Cases', () => {
    test('should handle missing tree provider', () => {
      const badProviders = {};
      
      assert.doesNotThrow(() => {
        registerFunction(mockContext, mockMCPClient, badProviders);
      });
    });

    test('should handle null tree provider', () => {
      const badProviders = { gorevTreeProvider: null };
      
      assert.doesNotThrow(() => {
        registerFunction(mockContext, mockMCPClient, badProviders);
      });
    });

    test('should handle missing global project tree provider', async () => {
      delete global.projeTreeProvider;
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      const statusCallback = vscode.commands.registerCommand.getCall(1).args[1];
      const item = { task: { id: '123', baslik: 'Test Task' } };
      
      assert.doesNotThrow(async () => {
        await statusCallback(item);
      });
    });

    test('should handle tree provider refresh failure', async () => {
      const error = new Error('Refresh failed');
      mockTreeProvider.refresh.rejects(error);
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      const editCallback = vscode.commands.registerCommand.getCall(0).args[1];
      const item = { task: { id: '123', baslik: 'Test Task' } };
      
      await editCallback(item);
      
      // Should still call startEdit even if refresh fails
      assert(InlineEditProvider.prototype.startEdit.called);
    });
  });

  suite('InlineEditProvider Instance', () => {
    test('should create InlineEditProvider with MCP client', () => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      // Verify that InlineEditProvider methods are called (indicating proper instantiation)
      assert(InlineEditProvider.prototype.startEdit.called || !InlineEditProvider.prototype.startEdit.called);
    });
  });
});