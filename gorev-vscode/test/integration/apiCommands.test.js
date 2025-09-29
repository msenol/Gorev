const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const MockAdapter = require('axios-mock-adapter');

let ApiClient;
let gorevCommands;
let projeCommands;

suite('API Commands Integration Tests', function() {
  this.timeout(15000);

  let sandbox;
  let apiClient;
  let mockAxios;
  let stubs;

  suiteSetup(async function() {
    // Import compiled modules
    const clientModule = require('../../out/api/client');
    ApiClient = clientModule.ApiClient;

    // Import commands
    gorevCommands = require('../../out/commands/gorevCommands');
    projeCommands = require('../../out/commands/projeCommands');
  });

  setup(function() {
    sandbox = sinon.createSandbox();
    apiClient = new ApiClient('http://localhost:5082');
    mockAxios = new MockAdapter(apiClient.axiosInstance);

    // Setup common stubs
    stubs = {
      showInformationMessage: sandbox.stub(vscode.window, 'showInformationMessage'),
      showErrorMessage: sandbox.stub(vscode.window, 'showErrorMessage'),
      showWarningMessage: sandbox.stub(vscode.window, 'showWarningMessage'),
      showInputBox: sandbox.stub(vscode.window, 'showInputBox'),
      showQuickPick: sandbox.stub(vscode.window, 'showQuickPick'),
      withProgress: sandbox.stub(vscode.window, 'withProgress')
    };

    // Make withProgress pass through
    stubs.withProgress.callsFake(async (options, task) => {
      return await task({ report: () => {} });
    });
  });

  teardown(function() {
    sandbox.restore();
    mockAxios.restore();
    if (apiClient) {
      apiClient.disconnect();
    }
  });

  suite('Task Commands', function() {
    test('UPDATE_TASK_STATUS should update task status', async function() {
      // Mock tree item
      const treeItem = {
        task: {
          id: 'task-1',
          baslik: 'Test Task',
          durum: 'beklemede',
          oncelik: 'yuksek'
        }
      };

      // Mock user selection
      stubs.showQuickPick.resolves({ value: 'devam_ediyor' });

      // Mock API response
      mockAxios.onPut('/tasks/task-1').reply(200, {
        success: true,
        data: {
          id: 'task-1',
          baslik: 'Test Task',
          durum: 'devam_ediyor'
        }
      });

      // Execute command
      await gorevCommands.updateTaskStatusCommand(apiClient, treeItem);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });

    test('DELETE_TASK should delete task', async function() {
      const treeItem = {
        task: {
          id: 'task-1',
          baslik: 'Task to Delete'
        }
      };

      // Mock confirmation
      stubs.showWarningMessage.resolves('Sil');

      // Mock API response
      mockAxios.onDelete('/tasks/task-1').reply(200, {
        success: true,
        message: 'Task deleted'
      });

      // Execute command
      await gorevCommands.deleteTaskCommand(apiClient, treeItem);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });

    test('CREATE_SUBTASK should create subtask', async function() {
      const parentItem = {
        task: {
          id: 'parent-1',
          baslik: 'Parent Task'
        }
      };

      // Mock user inputs
      stubs.showInputBox
        .onCall(0).resolves('Subtask Title')  // baslik
        .onCall(1).resolves('Subtask Description');  // aciklama

      stubs.showQuickPick
        .onCall(0).resolves({ value: 'yuksek' });  // oncelik

      // Mock API response
      mockAxios.onPost('/tasks/parent-1/subtasks').reply(201, {
        success: true,
        data: {
          id: 'subtask-1',
          baslik: 'Subtask Title',
          parent_id: 'parent-1'
        }
      });

      // Execute command
      await gorevCommands.createSubtaskCommand(apiClient, parentItem);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });

    test('CHANGE_PARENT should change task parent', async function() {
      const taskItem = {
        task: {
          id: 'task-1',
          baslik: 'Task to Move'
        }
      };

      // Mock available tasks
      mockAxios.onGet('/tasks', {
        params: { tum_projeler: true }
      }).reply(200, {
        success: true,
        data: [
          {
            id: 'task-2',
            baslik: 'Potential Parent',
            durum: 'devam_ediyor'
          }
        ],
        total: 1
      });

      // Mock user selection
      stubs.showQuickPick.resolves({
        id: 'task-2',
        baslik: 'Potential Parent'
      });

      // Mock API response
      mockAxios.onPut('/tasks/task-1/parent').reply(200, {
        success: true,
        data: {
          id: 'task-1',
          parent_id: 'task-2'
        }
      });

      // Execute command
      await gorevCommands.changeParentCommand(apiClient, taskItem);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });

    test('REMOVE_PARENT should remove task parent', async function() {
      const taskItem = {
        task: {
          id: 'task-1',
          baslik: 'Task with Parent',
          parent_id: 'parent-1'
        }
      };

      // Mock API response
      mockAxios.onPut('/tasks/task-1/parent').reply(200, {
        success: true,
        data: {
          id: 'task-1',
          parent_id: null
        }
      });

      // Execute command
      await gorevCommands.removeParentCommand(apiClient, taskItem);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });

    test('ADD_DEPENDENCY should add dependency', async function() {
      const taskItem = {
        task: {
          id: 'task-1',
          baslik: 'Dependent Task'
        }
      };

      // Mock available tasks
      mockAxios.onGet('/tasks', {
        params: { tum_projeler: true }
      }).reply(200, {
        success: true,
        data: [
          {
            id: 'task-2',
            baslik: 'Dependency Task',
            durum: 'tamamlandi'
          }
        ],
        total: 1
      });

      // Mock user selection
      stubs.showQuickPick.resolves({
        id: 'task-2',
        baslik: 'Dependency Task'
      });

      // Mock API response
      mockAxios.onPost('/tasks/task-1/dependencies').reply(201, {
        success: true,
        message: 'Dependency added'
      });

      // Execute command
      await gorevCommands.addDependencyCommand(apiClient, taskItem);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });
  });

  suite('Project Commands', function() {
    test('CREATE_PROJECT should create project', async function() {
      // Mock user inputs
      stubs.showInputBox
        .onCall(0).resolves('New Project')  // isim
        .onCall(1).resolves('Project Description');  // tanim

      // Mock API response
      mockAxios.onPost('/projects').reply(201, {
        success: true,
        data: {
          id: 'new-proj',
          isim: 'New Project',
          tanim: 'Project Description'
        }
      });

      // Execute command
      await projeCommands.createProjectCommand(apiClient);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });

    test('SET_ACTIVE_PROJECT should activate project', async function() {
      const projectItem = {
        project: {
          id: 'proj-1',
          isim: 'Project to Activate'
        }
      };

      // Mock API response
      mockAxios.onPost('/active-project').reply(200, {
        success: true,
        message: 'Project activated'
      });

      // Execute command
      await projeCommands.setActiveProjectCommand(apiClient, projectItem);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });

    test('SET_ACTIVE_PROJECT should remove active project when none', async function() {
      // Mock projects list
      mockAxios.onGet('/projects').reply(200, {
        success: true,
        data: [
          {
            id: 'proj-1',
            isim: 'Project 1',
            is_active: false
          }
        ],
        total: 1
      });

      // User cancels selection
      stubs.showQuickPick.resolves(undefined);

      // Mock API response
      mockAxios.onDelete('/active-project').reply(200, {
        success: true,
        message: 'Active project removed'
      });

      // Execute command
      await projeCommands.setActiveProjectCommand(apiClient);

      // Verify success message
      assert(stubs.showInformationMessage.called);
    });
  });

  suite('Error Handling', function() {
    test('should handle API errors in task update', async function() {
      const treeItem = {
        task: {
          id: 'task-1',
          baslik: 'Test Task'
        }
      };

      stubs.showQuickPick.resolves({ value: 'devam_ediyor' });

      // Mock API error
      mockAxios.onPut('/tasks/task-1').reply(500, {
        success: false,
        error: 'Internal server error'
      });

      // Execute command
      await gorevCommands.updateTaskStatusCommand(apiClient, treeItem);

      // Verify error message shown
      assert(stubs.showErrorMessage.called);
    });

    test('should handle network errors in project creation', async function() {
      stubs.showInputBox
        .onCall(0).resolves('New Project')
        .onCall(1).resolves('Description');

      // Mock network error
      mockAxios.onPost('/projects').networkError();

      // Execute command
      await projeCommands.createProjectCommand(apiClient);

      // Verify error message shown
      assert(stubs.showErrorMessage.called);
    });

    test('should handle user cancellation gracefully', async function() {
      const treeItem = {
        task: {
          id: 'task-1',
          baslik: 'Test Task'
        }
      };

      // User cancels
      stubs.showQuickPick.resolves(undefined);

      // Execute command
      await gorevCommands.updateTaskStatusCommand(apiClient, treeItem);

      // Should not show error
      assert(!stubs.showErrorMessage.called);
    });

    test('should handle not found errors', async function() {
      const treeItem = {
        task: {
          id: 'nonexistent',
          baslik: 'Nonexistent Task'
        }
      };

      stubs.showQuickPick.resolves({ value: 'tamamlandi' });

      mockAxios.onPut('/tasks/nonexistent').reply(404, {
        success: false,
        error: 'Task not found'
      });

      await gorevCommands.updateTaskStatusCommand(apiClient, treeItem);

      assert(stubs.showErrorMessage.called);
      const errorMsg = stubs.showErrorMessage.firstCall.args[0];
      assert(errorMsg.includes('404') || errorMsg.includes('not found'));
    });
  });

  suite('Command Integration Flow', function() {
    test('should create subtask and update parent hierarchy', async function() {
      const parentItem = {
        task: {
          id: 'parent-1',
          baslik: 'Parent Task'
        }
      };

      stubs.showInputBox
        .onCall(0).resolves('Subtask')
        .onCall(1).resolves('Description');

      stubs.showQuickPick.onCall(0).resolves({ value: 'orta' });

      // Mock subtask creation
      mockAxios.onPost('/tasks/parent-1/subtasks').reply(201, {
        success: true,
        data: {
          id: 'subtask-1',
          baslik: 'Subtask',
          parent_id: 'parent-1'
        }
      });

      // Mock hierarchy refresh
      mockAxios.onGet('/tasks/parent-1/hierarchy').reply(200, {
        success: true,
        data: {
          gorev: { id: 'parent-1', baslik: 'Parent Task' },
          alt_gorevler: [{ id: 'subtask-1', baslik: 'Subtask' }],
          toplam_alt_gorev: 1,
          tamamlanan_alt_gorev: 0
        }
      });

      await gorevCommands.createSubtaskCommand(apiClient, parentItem);

      assert(stubs.showInformationMessage.called);
    });

    test('should handle chained dependency additions', async function() {
      const taskItem = {
        task: {
          id: 'task-3',
          baslik: 'Task 3'
        }
      };

      mockAxios.onGet('/tasks', {
        params: { tum_projeler: true }
      }).reply(200, {
        success: true,
        data: [
          { id: 'task-1', baslik: 'Task 1', durum: 'tamamlandi' },
          { id: 'task-2', baslik: 'Task 2', durum: 'tamamlandi' }
        ],
        total: 2
      });

      stubs.showQuickPick.resolves({ id: 'task-1', baslik: 'Task 1' });

      mockAxios.onPost('/tasks/task-3/dependencies').reply(201, {
        success: true,
        message: 'Dependency added'
      });

      await gorevCommands.addDependencyCommand(apiClient, taskItem);

      assert(stubs.showInformationMessage.called);
    });
  });
});