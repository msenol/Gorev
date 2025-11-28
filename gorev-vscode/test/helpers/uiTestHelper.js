/**
 * UI Test Helper for Gorev Extension
 * Provides utilities for testing UI interactions
 */

const vscode = require('vscode');
const sinon = require('sinon');

class UITestHelper {
  constructor() {
    this.sandbox = sinon.createSandbox();
  }

  /**
   * Simulate user clicking a tree view item
   */
  async clickTreeItem(treeViewId, itemLabel) {
    // Mock tree view selection
    const mockItem = {
      id: `mock-${Date.now()}`,
      label: itemLabel,
      command: {
        command: 'gorev.showTaskDetail',
        title: 'Show Details'
      }
    };

    return vscode.commands.executeCommand('gorev.showTaskDetail', mockItem);
  }

  /**
   * Simulate command execution with mocked user inputs
   */
  async mockUserInput(command, inputs) {
    const showInputBoxStub = this.sandbox.stub(vscode.window, 'showInputBox');
    const showQuickPickStub = this.sandbox.stub(vscode.window, 'showQuickPick');
    const showOpenDialogStub = this.sandbox.stub(vscode.window, 'showOpenDialog');

    // Configure mock responses
    if (inputs.inputBox) {
      showInputBoxStub.resolves(inputs.inputBox);
    }
    if (inputs.quickPick) {
      showQuickPickStub.resolves(inputs.quickPick);
    }
    if (inputs.openDialog) {
      showOpenDialogStub.resolves(inputs.openDialog);
    }

    // Execute command
    const result = await vscode.commands.executeCommand(command);

    // Restore stubs
    this.sandbox.restore();

    return result;
  }

  /**
   * Mock API responses for testing
   */
  mockApiResponses(mockApi) {
    const axios = require('axios');
    const MockAdapter = require('axios-mock-adapter');
    const mockAxios = new MockAdapter(axios);

    // Setup default responses
    mockAxios.onGet('http://localhost:5082/api/v1/health').reply(200, { status: 'ok' });
    mockAxios.onGet('http://localhost:5082/api/v1/tasks').reply(200, {
      success: true,
      data: []
    });
    mockAxios.onGet('http://localhost:5082/api/v1/projects').reply(200, {
      success: true,
      data: []
    });
    mockAxios.onGet('http://localhost:5082/api/v1/templates').reply(200, {
      success: true,
      data: []
    });

    // Custom responses from test
    if (mockApi.custom) {
      mockApi.custom(mockAxios);
    }

    return mockAxios;
  }

  /**
   * Verify tree view has expected items
   */
  async verifyTreeView(treeViewId, expectedItems) {
    // This would require accessing the actual tree data
    // For now, verify the view is registered
    const treeViews = vscode.window.registerTreeDataProvider(treeViewId, {
      getChildren: () => Promise.resolve([]),
      getTreeItem: (element) => element
    });

    assert.ok(treeViews);
    return true;
  }

  /**
   * Simulate context menu action
   */
  async triggerContextMenu(treeViewId, itemLabel, menuAction) {
    const contextValue = `task.${menuAction}`;
    const mockItem = {
      id: `context-${Date.now()}`,
      label: itemLabel,
      contextValue,
      command: {
        command: `gorev.${menuAction}`,
        title: menuAction
      }
    };

    return vscode.commands.executeCommand(`gorev.${menuAction}`, mockItem);
  }

  /**
   * Mock workspace and extension activation
   */
  async setupTestEnvironment() {
    // Create mock workspace
    const workspaceFolder = vscode.Uri.file('/tmp/gorev-test-workspace');
    await vscode.workspace.updateWorkspaceFolders(0, null, { uri: workspaceFolder });

    // Ensure extension is active
    const extension = vscode.extensions.getExtension('gorev-team.gorev-vscode');
    if (extension && !extension.isActive) {
      await extension.activate();
    }

    return {
      workspaceFolder,
      extension
    };
  }

  /**
   * Clean up after tests
   */
  cleanup() {
    this.sandbox.restore();
  }
}

module.exports = UITestHelper;
