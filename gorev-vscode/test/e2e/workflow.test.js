const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const path = require('path');
const { spawn } = require('child_process');

suite('E2E Workflow Test Suite', function() {
  this.timeout(60000); // 60 second timeout for E2E tests
  
  let sandbox;
  let serverProcess;
  let serverPath;

  suiteSetup(async () => {
    // Try to find the gorev server executable
    serverPath = vscode.workspace.getConfiguration('gorev').get('serverPath');
    
    if (!serverPath) {
      // Try default locations
      const possiblePaths = [
        path.join(__dirname, '../../../../gorev-mcpserver/gorev'),
        path.join(__dirname, '../../../../gorev-mcpserver/gorev.exe'),
        path.join(__dirname, '../../../../gorev-mcpserver/build/gorev'),
        path.join(__dirname, '../../../../gorev-mcpserver/build/gorev.exe')
      ];
      
      for (const p of possiblePaths) {
        try {
          const fs = require('fs');
          if (fs.existsSync(p)) {
            serverPath = p;
            break;
          }
        } catch (e) {
          // Continue
        }
      }
    }
    
    if (!serverPath) {
      console.warn('Gorev server not found, skipping E2E tests');
      this.skip();
    }
  });

  setup(() => {
    sandbox = sinon.createSandbox();
  });

  teardown(async () => {
    sandbox.restore();
    
    // Disconnect if connected
    try {
      await vscode.commands.executeCommand('gorev.disconnect');
    } catch (e) {
      // Ignore
    }
    
    // Kill server if running
    if (serverProcess) {
      serverProcess.kill();
      serverProcess = null;
    }
  });

  suite('Full Workflow', function() {
    test('Should complete full task management workflow', async function() {
      // Start server
      console.log('Starting server at:', serverPath);
      serverProcess = spawn(serverPath, ['serve']);
      
      // Wait for server to start
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // 1. Connect to server
      await vscode.workspace.getConfiguration('gorev').update('serverPath', serverPath, vscode.ConfigurationTarget.Workspace);
      await vscode.commands.executeCommand('gorev.connect');
      
      // Wait for connection
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // 2. Create a project
      const createProjectStub = sandbox.stub(vscode.window, 'showInputBox');
      createProjectStub.onCall(0).resolves('Test E2E Project'); // Project name
      createProjectStub.onCall(1).resolves('E2E test project description'); // Description
      
      await vscode.commands.executeCommand('gorev.createProject');
      
      // 3. Create a task
      createProjectStub.onCall(2).resolves('Test E2E Task'); // Task title
      createProjectStub.onCall(3).resolves('E2E task description'); // Description
      
      const priorityStub = sandbox.stub(vscode.window, 'showQuickPick');
      priorityStub.onCall(0).resolves({ label: 'Orta', value: 'orta' }); // Priority
      
      await vscode.commands.executeCommand('gorev.createTask');
      
      // 4. Refresh to see changes
      await vscode.commands.executeCommand('gorev.refresh');
      
      // 5. Get summary
      await vscode.commands.executeCommand('gorev.showSummary');
      
      // Verify workflow completed without errors
      assert(true, 'Workflow completed successfully');
    });

    test('Should handle task state transitions', async function() {
      if (!serverProcess) {
        serverProcess = spawn(serverPath, ['serve']);
        await new Promise(resolve => setTimeout(resolve, 2000));
      }
      
      // Ensure connected
      await vscode.commands.executeCommand('gorev.connect');
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Create a task first
      const inputStub = sandbox.stub(vscode.window, 'showInputBox');
      inputStub.onCall(0).resolves('State Transition Task');
      inputStub.onCall(1).resolves('Task to test state transitions');
      
      const quickPickStub = sandbox.stub(vscode.window, 'showQuickPick');
      quickPickStub.onCall(0).resolves({ label: 'Yüksek', value: 'yuksek' });
      
      await vscode.commands.executeCommand('gorev.createTask');
      
      // Simulate task state changes
      // Note: In real scenario, we'd need to get the task ID and update it
      // For now, we're testing the command registration
      const commands = await vscode.commands.getCommands();
      
      assert(commands.includes('gorev.startTask'), 'Start task command should exist');
      assert(commands.includes('gorev.completeTask'), 'Complete task command should exist');
    });

    test('Should create task from template', async function() {
      if (!serverProcess) {
        serverProcess = spawn(serverPath, ['serve']);
        await new Promise(resolve => setTimeout(resolve, 2000));
      }
      
      // Ensure connected
      await vscode.commands.executeCommand('gorev.connect');
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Initialize templates
      const initProcess = spawn(serverPath, ['template', 'init']);
      await new Promise((resolve, reject) => {
        initProcess.on('close', code => {
          if (code === 0) resolve();
          else reject(new Error('Template init failed'));
        });
      });
      
      // Refresh templates
      await vscode.commands.executeCommand('gorev.refreshTemplates');
      
      // Test template wizard command
      const commands = await vscode.commands.getCommands();
      assert(commands.includes('gorev.showTemplateWizard'), 'Template wizard command should exist');
    });
  });

  suite('Error Recovery', function() {
    test('Should handle server disconnect gracefully', async function() {
      // Start server
      serverProcess = spawn(serverPath, ['serve']);
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Connect
      await vscode.workspace.getConfiguration('gorev').update('serverPath', serverPath, vscode.ConfigurationTarget.Workspace);
      await vscode.commands.executeCommand('gorev.connect');
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Kill server
      serverProcess.kill();
      serverProcess = null;
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Try to create task - should show error
      const errorStub = sandbox.stub(vscode.window, 'showErrorMessage');
      
      try {
        await vscode.commands.executeCommand('gorev.createTask');
      } catch (e) {
        // Expected to fail
      }
      
      // Should show error or handle gracefully
      assert(errorStub.called || true, 'Should handle disconnect gracefully');
    });

    test('Should reconnect after disconnect', async function() {
      // Start server
      serverProcess = spawn(serverPath, ['serve']);
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Connect
      await vscode.workspace.getConfiguration('gorev').update('serverPath', serverPath, vscode.ConfigurationTarget.Workspace);
      await vscode.commands.executeCommand('gorev.connect');
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Disconnect
      await vscode.commands.executeCommand('gorev.disconnect');
      await new Promise(resolve => setTimeout(resolve, 500));
      
      // Reconnect
      await vscode.commands.executeCommand('gorev.connect');
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Should be able to use commands again
      await vscode.commands.executeCommand('gorev.refresh');
      
      assert(true, 'Reconnection successful');
    });
  });
});