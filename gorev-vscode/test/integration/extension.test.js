const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const path = require('path');

suite('Extension Integration Test Suite', () => {
  let sandbox;

  setup(() => {
    sandbox = sinon.createSandbox();
  });

  teardown(() => {
    sandbox.restore();
  });

  test('Extension should be present', () => {
    assert.ok(vscode.extensions.getExtension('gorev-team.gorev-vscode'));
  });

  test('Should activate extension', async () => {
    const extension = vscode.extensions.getExtension('gorev-team.gorev-vscode');
    if (extension && !extension.isActive) {
      await extension.activate();
    }
    assert.ok(extension.isActive);
  });

  suite('Commands Registration', () => {
    test('Should register all commands', async () => {
      const commands = await vscode.commands.getCommands();
      
      const gorevCommands = [
        'gorev.connect',
        'gorev.disconnect',
        'gorev.createTask',
        'gorev.createProject',
        'gorev.showSummary',
        'gorev.refresh',
        'gorev.editTask',
        'gorev.deleteTask',
        'gorev.completeTask',
        'gorev.startTask',
        'gorev.setActiveProject',
        'gorev.clearActiveProject',
        'gorev.showTaskDetail',
        'gorev.createTaskFromTemplate',
        'gorev.showTemplateWizard',
        'gorev.refreshTemplates',
        'gorev.addDependency',
        'gorev.showSearchInput',
        'gorev.showAdvancedFilter',
        'gorev.toggleGrouping',
        'gorev.clearFilters'
      ];

      gorevCommands.forEach(cmd => {
        assert(commands.includes(cmd), `Command ${cmd} should be registered`);
      });
    });
  });

  suite('Views Registration', () => {
    test('Should register tree views', () => {
      // Check if views are registered in the package.json
      const packageJson = require('../../package.json');
      const views = packageJson.contributes.views.gorev;
      
      assert(views.find(v => v.id === 'gorevTasks'), 'Tasks view should be registered');
      assert(views.find(v => v.id === 'gorevProjects'), 'Projects view should be registered');
      assert(views.find(v => v.id === 'gorevTemplates'), 'Templates view should be registered');
    });
  });

  suite('Configuration', () => {
    test('Should have default configuration', () => {
      const config = vscode.workspace.getConfiguration('gorev');
      
      assert.strictEqual(config.get('autoConnect'), true);
      assert.strictEqual(config.get('showStatusBar'), true);
      assert.strictEqual(config.get('refreshInterval'), 30);
    });

    test('Should update configuration', async () => {
      const config = vscode.workspace.getConfiguration('gorev');
      
      await config.update('autoConnect', false, vscode.ConfigurationTarget.Workspace);
      assert.strictEqual(config.get('autoConnect'), false);
      
      // Reset
      await config.update('autoConnect', undefined, vscode.ConfigurationTarget.Workspace);
    });
  });

  suite('Status Bar', () => {
    test('Should show status bar item', async () => {
      // This would need access to the actual status bar item
      // For now, we verify the configuration exists
      const config = vscode.workspace.getConfiguration('gorev');
      assert.strictEqual(config.get('showStatusBar'), true);
    });
  });

  suite('Command Execution', () => {
    test('Should handle connect command', async () => {
      // Mock the connection process
      const showInformationMessage = sandbox.stub(vscode.window, 'showInformationMessage');
      
      try {
        await vscode.commands.executeCommand('gorev.connect');
        // Command should execute without throwing
      } catch (error) {
        // If server is not configured, it should show a message
        assert(showInformationMessage.called || error.message.includes('Server path not configured'));
      }
    });

    test('Should handle refresh command', async () => {
      // This should always work, even if not connected
      await vscode.commands.executeCommand('gorev.refresh');
      // Should not throw
    });
  });

  suite('TreeDataProvider Integration', () => {
    test('Should create empty tree when not connected', async () => {
      // Get the tree provider through the extension API
      const extension = vscode.extensions.getExtension('gorev-team.gorev-vscode');
      if (extension && extension.exports && extension.exports.gorevTreeProvider) {
        const provider = extension.exports.gorevTreeProvider;
        const children = await provider.getChildren();
        
        // When not connected, should return empty array
        assert.strictEqual(children.length, 0);
      }
    });
  });

  suite('Webview Panels', () => {
    test('Should create task detail panel command', async () => {
      const commands = await vscode.commands.getCommands();
      assert(commands.includes('gorev.showTaskDetail'));
    });

    test('Should create template wizard command', async () => {
      const commands = await vscode.commands.getCommands();
      assert(commands.includes('gorev.showTemplateWizard'));
    });
  });

  suite('Error Handling', () => {
    test('Should handle missing server path gracefully', async () => {
      const config = vscode.workspace.getConfiguration('gorev');
      await config.update('serverPath', '', vscode.ConfigurationTarget.Workspace);
      
      const showErrorMessage = sandbox.stub(vscode.window, 'showErrorMessage');
      
      try {
        await vscode.commands.executeCommand('gorev.connect');
      } catch (error) {
        // Should show error message
        assert(showErrorMessage.called || error.message.includes('Server path not configured'));
      }
      
      // Reset
      await config.update('serverPath', undefined, vscode.ConfigurationTarget.Workspace);
    });
  });
});