const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');

suite('Subtask Commands Integration Test Suite', () => {
    let sandbox;
    let mockMcpClient;
    let mockTreeProvider;

    setup(() => {
        sandbox = sinon.createSandbox();
        
        // Mock MCP Client
        mockMcpClient = {
            callTool: sandbox.stub(),
            isConnected: sandbox.stub().returns(true)
        };

        // Mock Tree Provider
        mockTreeProvider = {
            refresh: sandbox.stub()
        };

        // Mock vscode.window methods
        sandbox.stub(vscode.window, 'showInputBox');
        sandbox.stub(vscode.window, 'showQuickPick');
        sandbox.stub(vscode.window, 'showInformationMessage');
        sandbox.stub(vscode.window, 'showErrorMessage');
    });

    teardown(() => {
        sandbox.restore();
    });

    suite('Create Subtask Command Tests', () => {
        test('should create subtask with all required fields', async () => {
            const parentTask = {
                id: 'parent123',
                baslik: 'Parent Task'
            };

            // Mock user inputs
            vscode.window.showInputBox
                .onFirstCall().resolves('New Subtask Title')
                .onSecondCall().resolves('Subtask description');
            
            vscode.window.showQuickPick.resolves({ value: 'orta' });

            mockMcpClient.callTool.resolves();

            // Execute command
            await vscode.commands.executeCommand('gorev.createSubtask', { task: parentTask });

            // Verify MCP call
            assert.ok(mockMcpClient.callTool.calledWith('gorev_alt_gorev_olustur', {
                parent_id: 'parent123',
                baslik: 'New Subtask Title',
                aciklama: 'Subtask description',
                oncelik: 'orta'
            }));

            // Verify success message
            assert.ok(vscode.window.showInformationMessage.calledWith('Subtask created successfully'));
        });

        test('should handle subtask creation cancellation', async () => {
            const parentTask = {
                id: 'parent123',
                baslik: 'Parent Task'
            };

            // User cancels at title input
            vscode.window.showInputBox.resolves(undefined);

            await vscode.commands.executeCommand('gorev.createSubtask', { task: parentTask });

            // Should not call MCP
            assert.ok(mockMcpClient.callTool.notCalled);
        });

        test('should handle subtask creation error', async () => {
            const parentTask = {
                id: 'parent123',
                baslik: 'Parent Task'
            };

            vscode.window.showInputBox
                .onFirstCall().resolves('New Subtask')
                .onSecondCall().resolves('');
            
            vscode.window.showQuickPick.resolves({ value: 'yuksek' });

            mockMcpClient.callTool.rejects(new Error('Network error'));

            await vscode.commands.executeCommand('gorev.createSubtask', { task: parentTask });

            assert.ok(vscode.window.showErrorMessage.calledWith('Failed to create subtask: Network error'));
        });
    });

    suite('Change Parent Command Tests', () => {
        test('should change parent successfully', async () => {
            const task = {
                id: 'task123',
                baslik: 'Task to Move'
            };

            // Mock task list response
            mockMcpClient.callTool
                .withArgs('gorev_listele').resolves({
                    content: [{
                        text: `[✓] Task 1 (orta öncelik)
  ID: task1
[⏳] Task 2 (yüksek öncelik)
  ID: task2
[🔄] Current Task (orta öncelik)
  ID: task123`
                    }]
                })
                .withArgs('gorev_ust_degistir').resolves();

            // User selects Task 2 as new parent
            vscode.window.showQuickPick.resolves({ value: 'task2' });

            await vscode.commands.executeCommand('gorev.changeParent', { task });

            assert.ok(mockMcpClient.callTool.calledWith('gorev_ust_degistir', {
                id: 'task123',
                yeni_parent_id: 'task2'
            }));

            assert.ok(vscode.window.showInformationMessage.calledWith('Parent changed successfully'));
        });

        test('should handle making task root (no parent)', async () => {
            const task = {
                id: 'task123',
                baslik: 'Child Task',
                parent_id: 'parent456'
            };

            mockMcpClient.callTool
                .withArgs('gorev_listele').resolves({
                    content: [{ text: 'Some tasks...' }]
                });

            // User selects "No parent"
            vscode.window.showQuickPick.resolves({ value: null });

            await vscode.commands.executeCommand('gorev.changeParent', { task });

            assert.ok(mockMcpClient.callTool.calledWith('gorev_ust_degistir', {
                id: 'task123',
                yeni_parent_id: ''
            }));
        });
    });

    suite('Remove Parent Command Tests', () => {
        test('should remove parent and make task root', async () => {
            const task = {
                id: 'child123',
                baslik: 'Child Task',
                parent_id: 'parent456'
            };

            mockMcpClient.callTool.resolves();

            await vscode.commands.executeCommand('gorev.removeParent', { task });

            assert.ok(mockMcpClient.callTool.calledWith('gorev_ust_degistir', {
                id: 'child123',
                yeni_parent_id: ''
            }));

            assert.ok(vscode.window.showInformationMessage.calledWith('Task is now a root task'));
        });

        test('should refresh tree after parent removal', async () => {
            const task = {
                id: 'child123',
                parent_id: 'parent456'
            };

            mockMcpClient.callTool.resolves();

            const providers = {
                gorevTreeProvider: mockTreeProvider
            };

            // Simulate command with providers context
            await vscode.commands.executeCommand('gorev.removeParent', { task });

            // Tree should be refreshed
            assert.ok(mockTreeProvider.refresh.called);
        });
    });

    suite('Error Handling Tests', () => {
        test('should show error for circular dependency', async () => {
            const task = { id: '123', baslik: 'Task A' };

            mockMcpClient.callTool
                .withArgs('gorev_listele').resolves({
                    content: [{ text: '[✓] Task B\n  ID: 456' }]
                })
                .withArgs('gorev_ust_degistir').rejects(new Error('dairesel bağımlılık tespit edildi'));

            vscode.window.showQuickPick.resolves({ value: '456' });

            await vscode.commands.executeCommand('gorev.changeParent', { task });

            assert.ok(vscode.window.showErrorMessage.calledWith(sinon.match(/dairesel bağımlılık/)));
        });

        test('should show error for different project constraint', async () => {
            const task = { 
                id: '123', 
                baslik: 'Task A',
                proje_id: 'proj1'
            };

            mockMcpClient.callTool
                .withArgs('gorev_listele').resolves({
                    content: [{ text: '[✓] Task B\n  ID: 456' }]
                })
                .withArgs('gorev_ust_degistir').rejects(new Error('alt görev ve üst görev aynı projede olmalı'));

            vscode.window.showQuickPick.resolves({ value: '456' });

            await vscode.commands.executeCommand('gorev.changeParent', { task });

            assert.ok(vscode.window.showErrorMessage.calledWith(sinon.match(/aynı projede/)));
        });
    });
});