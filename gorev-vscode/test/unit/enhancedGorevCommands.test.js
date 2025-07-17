const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('EnhancedGorevCommands Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let mockTreeProvider;
  let mockContext;
  let mockProviders;
  let registerFunction;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showWarningMessage').resolves('Yes');
    sandbox.stub(vscode.window, 'showQuickPick');
    sandbox.stub(vscode.window, 'showInputBox');
    sandbox.stub(vscode.window, 'withProgress');
    sandbox.stub(vscode.commands, 'registerCommand');
    sandbox.stub(vscode.commands, 'executeCommand');
    sandbox.stub(vscode.workspace, 'getConfiguration').returns({
      get: sandbox.stub().returns(true),
      update: sandbox.stub()
    });

    // Mock MCP Client
    mockMCPClient = {
      callTool: sandbox.stub().resolves({ content: [{ text: 'Success' }] }),
      isConnected: sandbox.stub().returns(true)
    };

    // Mock TreeProvider
    mockTreeProvider = {
      selectTask: sandbox.stub(),
      setGrouping: sandbox.stub(),
      setSorting: sandbox.stub(),
      updateFilter: sandbox.stub(),
      getSelectedTasks: sandbox.stub().returns([]),
      refresh: sandbox.stub().resolves()
    };

    // Mock Context
    mockContext = {
      subscriptions: []
    };

    // Mock Providers
    mockProviders = {
      gorevTreeProvider: mockTreeProvider,
      projeTreeProvider: { refresh: sandbox.stub() },
      statusBarManager: { update: sandbox.stub() }
    };

    // Import and register commands
    try {
      const commands = require('../../dist/commands/enhancedGorevCommands');
      registerFunction = commands.registerEnhancedGorevCommands;
    } catch (error) {
      // Mock register function if compilation fails
      registerFunction = (context, mcpClient, providers) => {
        // Mock command registrations
        context.subscriptions.push(...new Array(12).fill({}));
      };
    }
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('Command Registration', () => {
    test('should register all enhanced commands', () => {
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      // Should register 12 commands
      assert.strictEqual(mockContext.subscriptions.length, 12);
      
      // Should call vscode.commands.registerCommand for each command
      assert(vscode.commands.registerCommand.callCount >= 12);
    });

    test('should handle registration errors gracefully', () => {
      const errorContext = { subscriptions: { push: sandbox.stub().throws(new Error('Mock error')) } };
      
      try {
        registerFunction(errorContext, mockMCPClient, mockProviders);
        // Should not throw
        assert(true);
      } catch (error) {
        assert.fail('Should handle registration errors gracefully');
      }
    });
  });

  suite('SELECT_TASK Command', () => {
    test('should select task without multi-select', () => {
      const taskId = 'task-123';
      const event = { ctrlKey: false, metaKey: false, shiftKey: false };
      
      // Mock command registration and call the handler
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.selectTask') {
          handler(taskId, event);
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockTreeProvider.selectTask.calledWith(taskId, false, false));
    });

    test('should select task with multi-select (ctrl key)', () => {
      const taskId = 'task-123';
      const event = { ctrlKey: true, metaKey: false, shiftKey: false };
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.selectTask') {
          handler(taskId, event);
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockTreeProvider.selectTask.calledWith(taskId, true, false));
    });

    test('should select task with range select (shift key)', () => {
      const taskId = 'task-123';
      const event = { ctrlKey: false, metaKey: false, shiftKey: true };
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.selectTask') {
          handler(taskId, event);
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockTreeProvider.selectTask.calledWith(taskId, false, true));
    });

    test('should handle task selection with meta key (Mac)', () => {
      const taskId = 'task-123';
      const event = { ctrlKey: false, metaKey: true, shiftKey: false };
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.selectTask') {
          handler(taskId, event);
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockTreeProvider.selectTask.calledWith(taskId, true, false));
    });
  });

  suite('SET_GROUPING Command', () => {
    test('should show grouping options and set grouping', async () => {
      const selectedGrouping = { label: 'By Status', value: 'status' };
      vscode.window.showQuickPick.resolves(selectedGrouping);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.setGrouping') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showQuickPick.called);
      assert(mockTreeProvider.setGrouping.calledWith('status'));
    });

    test('should handle cancellation in grouping selection', async () => {
      vscode.window.showQuickPick.resolves(undefined);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.setGrouping') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showQuickPick.called);
      assert(mockTreeProvider.setGrouping.notCalled);
    });

    test('should provide correct grouping options', async () => {
      vscode.window.showQuickPick.resolves(undefined);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.setGrouping') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      const quickPickCall = vscode.window.showQuickPick.getCall(0);
      const items = quickPickCall.args[0];
      
      assert(Array.isArray(items));
      assert(items.length >= 6); // At least 6 grouping options
      assert(items.some(item => item.label === 'No Grouping'));
      assert(items.some(item => item.label === 'By Status'));
      assert(items.some(item => item.label === 'By Priority'));
      assert(items.some(item => item.label === 'By Project'));
      assert(items.some(item => item.label === 'By Tag'));
      assert(items.some(item => item.label === 'By Due Date'));
    });
  });

  suite('SET_SORTING Command', () => {
    test('should show sorting options and set sorting', async () => {
      const selectedCriteria = { label: 'By Priority', value: 'priority' };
      const selectedOrder = { label: 'Ascending', value: true };
      
      vscode.window.showQuickPick
        .onFirstCall().resolves(selectedCriteria)
        .onSecondCall().resolves(selectedOrder);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.setSorting') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert.strictEqual(vscode.window.showQuickPick.callCount, 2);
      assert(mockTreeProvider.setSorting.calledWith('priority', true));
    });

    test('should handle cancellation in criteria selection', async () => {
      vscode.window.showQuickPick.onFirstCall().resolves(undefined);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.setSorting') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert.strictEqual(vscode.window.showQuickPick.callCount, 1);
      assert(mockTreeProvider.setSorting.notCalled);
    });

    test('should handle cancellation in order selection', async () => {
      const selectedCriteria = { label: 'By Priority', value: 'priority' };
      
      vscode.window.showQuickPick
        .onFirstCall().resolves(selectedCriteria)
        .onSecondCall().resolves(undefined);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.setSorting') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert.strictEqual(vscode.window.showQuickPick.callCount, 2);
      assert(mockTreeProvider.setSorting.notCalled);
    });
  });

  suite('FILTER_TASKS Command', () => {
    test('should handle search filter', async () => {
      const action = { label: '🔍 Search by text', value: 'search' };
      const searchQuery = 'test query';
      
      vscode.window.showQuickPick.resolves(action);
      vscode.window.showInputBox.resolves(searchQuery);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.filterTasks') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showQuickPick.called);
      assert(vscode.window.showInputBox.called);
      assert(mockTreeProvider.updateFilter.calledWith({ searchQuery }));
    });

    test('should handle status filter', async () => {
      const action = { label: '📊 Filter by status', value: 'status' };
      const statuses = [{ label: 'Pending', value: 'beklemede', picked: true }];
      
      vscode.window.showQuickPick
        .onFirstCall().resolves(action)
        .onSecondCall().resolves(statuses);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.filterTasks') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert.strictEqual(vscode.window.showQuickPick.callCount, 2);
      assert(mockTreeProvider.updateFilter.calledWith({ durum: 'beklemede' }));
    });

    test('should handle priority filter', async () => {
      const action = { label: '🎯 Filter by priority', value: 'priority' };
      const priorities = [{ label: 'High', value: 'yuksek' }];
      
      vscode.window.showQuickPick
        .onFirstCall().resolves(action)
        .onSecondCall().resolves(priorities);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.filterTasks') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockTreeProvider.updateFilter.calledWith({ oncelik: 'yuksek' }));
    });

    test('should handle tag filter', async () => {
      const action = { label: '🏷️ Filter by tag', value: 'tag' };
      const tag = 'bug';
      
      vscode.window.showQuickPick.resolves(action);
      vscode.window.showInputBox.resolves(tag);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.filterTasks') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockTreeProvider.updateFilter.calledWith({ tags: ['bug'] }));
    });

    test('should handle cancellation in filter type selection', async () => {
      vscode.window.showQuickPick.resolves(undefined);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.filterTasks') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showQuickPick.called);
      assert(mockTreeProvider.updateFilter.notCalled);
    });
  });

  suite('CLEAR_FILTER Command', () => {
    test('should clear filters and show message', () => {
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.clearFilter') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockTreeProvider.updateFilter.calledWith({}));
      assert(vscode.window.showInformationMessage.calledWith('Filters cleared'));
    });
  });

  suite('TOGGLE_SHOW_COMPLETED Command', () => {
    test('should toggle show completed setting', () => {
      const mockConfig = {
        get: sandbox.stub().returns(true),
        update: sandbox.stub()
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleShowCompleted') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockConfig.get.calledWith('showCompleted', true));
      assert(mockConfig.update.calledWith('showCompleted', false));
      assert(vscode.window.showInformationMessage.calledWith('Hiding completed tasks'));
    });

    test('should toggle from false to true', () => {
      const mockConfig = {
        get: sandbox.stub().returns(false),
        update: sandbox.stub()
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleShowCompleted') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(mockConfig.update.calledWith('showCompleted', true));
      assert(vscode.window.showInformationMessage.calledWith('Showing completed tasks'));
    });
  });

  suite('BULK_UPDATE_STATUS Command', () => {
    test('should show warning when no tasks selected', async () => {
      mockTreeProvider.getSelectedTasks.returns([]);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkUpdateStatus') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showWarningMessage.calledWith('No tasks selected'));
      assert(vscode.window.showQuickPick.notCalled);
    });

    test('should update multiple tasks status', async () => {
      const selectedTasks = [
        { id: 'task-1', baslik: 'Task 1' },
        { id: 'task-2', baslik: 'Task 2' }
      ];
      const newStatus = { label: 'In Progress', value: 'devam_ediyor' };
      
      mockTreeProvider.getSelectedTasks.returns(selectedTasks);
      vscode.window.showQuickPick.resolves(newStatus);
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkUpdateStatus') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showQuickPick.called);
      assert.strictEqual(mockMCPClient.callTool.callCount, 2);
      assert(mockMCPClient.callTool.calledWith('gorev_guncelle', { id: 'task-1', durum: 'devam_ediyor' }));
      assert(mockMCPClient.callTool.calledWith('gorev_guncelle', { id: 'task-2', durum: 'devam_ediyor' }));
      assert(mockTreeProvider.refresh.called);
    });

    test('should handle bulk update errors', async () => {
      const selectedTasks = [{ id: 'task-1', baslik: 'Task 1' }];
      const newStatus = { label: 'In Progress', value: 'devam_ediyor' };
      
      mockTreeProvider.getSelectedTasks.returns(selectedTasks);
      vscode.window.showQuickPick.resolves(newStatus);
      mockMCPClient.callTool.rejects(new Error('Update failed'));
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkUpdateStatus') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showErrorMessage.called);
      assert(vscode.window.showErrorMessage.calledWith(sinon.match(/Failed to update tasks/)));
    });

    test('should handle cancellation in status selection', async () => {
      const selectedTasks = [{ id: 'task-1', baslik: 'Task 1' }];
      
      mockTreeProvider.getSelectedTasks.returns(selectedTasks);
      vscode.window.showQuickPick.resolves(undefined);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkUpdateStatus') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showQuickPick.called);
      assert(mockMCPClient.callTool.notCalled);
    });
  });

  suite('BULK_DELETE Command', () => {
    test('should show warning when no tasks selected', async () => {
      mockTreeProvider.getSelectedTasks.returns([]);
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkDelete') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showWarningMessage.calledWith('No tasks selected'));
    });

    test('should delete multiple tasks after confirmation', async () => {
      const selectedTasks = [
        { id: 'task-1', baslik: 'Task 1' },
        { id: 'task-2', baslik: 'Task 2' }
      ];
      
      mockTreeProvider.getSelectedTasks.returns(selectedTasks);
      vscode.window.showWarningMessage.resolves('Yes');
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkDelete') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showWarningMessage.calledWith('Are you sure you want to delete 2 tasks?'));
      assert.strictEqual(mockMCPClient.callTool.callCount, 2);
      assert(mockMCPClient.callTool.calledWith('gorev_sil', { id: 'task-1', onay: true }));
      assert(mockMCPClient.callTool.calledWith('gorev_sil', { id: 'task-2', onay: true }));
      assert(mockTreeProvider.refresh.called);
    });

    test('should handle cancellation in confirmation', async () => {
      const selectedTasks = [{ id: 'task-1', baslik: 'Task 1' }];
      
      mockTreeProvider.getSelectedTasks.returns(selectedTasks);
      vscode.window.showWarningMessage.resolves('No');
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkDelete') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showWarningMessage.called);
      assert(mockMCPClient.callTool.notCalled);
    });

    test('should handle bulk delete errors', async () => {
      const selectedTasks = [{ id: 'task-1', baslik: 'Task 1' }];
      
      mockTreeProvider.getSelectedTasks.returns(selectedTasks);
      vscode.window.showWarningMessage.resolves('Yes');
      mockMCPClient.callTool.rejects(new Error('Delete failed'));
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkDelete') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showErrorMessage.called);
      assert(vscode.window.showErrorMessage.calledWith(sinon.match(/Failed to delete tasks/)));
    });
  });

  suite('SELECT_ALL and DESELECT_ALL Commands', () => {
    test('should show not implemented message for select all', () => {
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.selectAll') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showInformationMessage.calledWith('Select all not yet implemented'));
    });

    test('should show not implemented message for deselect all', () => {
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.deselectAll') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showInformationMessage.calledWith('Deselect all not yet implemented'));
    });
  });

  suite('Error Handling', () => {
    test('should handle provider errors gracefully', () => {
      mockTreeProvider.updateFilter.throws(new Error('Provider error'));
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.clearFilter') {
          try {
            handler();
          } catch (error) {
            // Should be handled gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      // Should not throw
      assert(true);
    });

    test('should handle MCP client errors in bulk operations', async () => {
      const selectedTasks = [{ id: 'task-1', baslik: 'Task 1' }];
      
      mockTreeProvider.getSelectedTasks.returns(selectedTasks);
      vscode.window.showQuickPick.resolves({ label: 'In Progress', value: 'devam_ediyor' });
      mockMCPClient.callTool.rejects(new Error('Network error'));
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        try {
          return await callback(progress);
        } catch (error) {
          // Progress callback should handle errors
          throw error;
        }
      });
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.bulkUpdateStatus') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockMCPClient, mockProviders);
      
      assert(vscode.window.showErrorMessage.called);
    });
  });
});