const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const { DragDropController } = require('../../out/providers/dragDropController');
const { 
    DragDataType, 
    DropTargetType,
    GroupingStrategy 
} = require('../../out/models/treeModels');
const { GorevDurum, GorevOncelik } = require('../../out/models/common');

suite('DragDropController Test Suite', () => {
    let sandbox;
    let mockMcpClient;
    let controller;
    let mockDataTransfer;

    setup(() => {
        sandbox = sinon.createSandbox();
        
        // Mock MCP Client
        mockMcpClient = {
            callTool: sandbox.stub(),
            isConnected: sandbox.stub().returns(true)
        };

        // Mock DataTransfer
        mockDataTransfer = {
            get: sandbox.stub(),
            set: sandbox.stub()
        };

        // Mock vscode.workspace.getConfiguration
        sandbox.stub(vscode.workspace, 'getConfiguration').returns({
            get: (key, defaultValue) => defaultValue
        });

        controller = new DragDropController(mockMcpClient);
    });

    teardown(() => {
        sandbox.restore();
    });

    suite('handleDrag Tests', () => {
        test('should handle single task drag', async () => {
            const task = {
                id: '123',
                baslik: 'Test Task',
                parent_id: null
            };

            const source = [{
                task: task,
                parent: { groupKey: 'beklemede' }
            }];

            const mockDataTransferItem = sandbox.stub();
            mockDataTransfer.set = sandbox.spy((type, item) => {
                assert.strictEqual(type, DragDataType.Task);
                assert.ok(item);
            });

            await controller.handleDrag(source, mockDataTransfer, {});
            
            assert.ok(mockDataTransfer.set.calledOnce);
        });

        test('should handle multiple tasks drag', async () => {
            const tasks = [
                { task: { id: '1', baslik: 'Task 1' } },
                { task: { id: '2', baslik: 'Task 2' } },
                { task: { id: '3', baslik: 'Task 3' } }
            ];

            mockDataTransfer.set = sandbox.spy((type, item) => {
                assert.strictEqual(type, DragDataType.Tasks);
            });

            await controller.handleDrag(tasks, mockDataTransfer, {});
            
            assert.ok(mockDataTransfer.set.calledOnce);
        });
    });

    suite('handleDrop Tests', () => {
        test('should handle drop on task for parent change', async () => {
            const sourceTask = {
                id: '123',
                baslik: 'Source Task',
                parent_id: null
            };

            const targetTask = {
                id: '456',
                baslik: 'Target Task'
            };

            const dragData = {
                type: DragDataType.Task,
                task: sourceTask
            };

            mockDataTransfer.get.returns({
                value: dragData
            });

            // Mock quick pick for action selection
            sandbox.stub(vscode.window, 'showQuickPick').resolves({
                value: 'make_subtask'
            });

            mockMcpClient.callTool.resolves();

            const target = {
                task: targetTask
            };

            await controller.handleDrop(target, mockDataTransfer, {});

            assert.ok(mockMcpClient.callTool.calledWith('gorev_ust_degistir', {
                id: sourceTask.id,
                yeni_parent_id: targetTask.id
            }));
        });

        test('should handle drop on empty area to remove parent', async () => {
            const task = {
                id: '123',
                baslik: 'Child Task',
                parent_id: '456'
            };

            const dragData = {
                type: DragDataType.Task,
                task: task
            };

            mockDataTransfer.get.returns({
                value: dragData
            });

            mockMcpClient.callTool.resolves();

            // Drop on empty area (target is undefined)
            await controller.handleDrop(undefined, mockDataTransfer, {});

            assert.ok(mockMcpClient.callTool.calledWith('gorev_ust_degistir', {
                id: task.id,
                yeni_parent_id: ''
            }));
        });

        test('should handle circular dependency error', async () => {
            const sourceTask = {
                id: '123',
                baslik: 'Source Task'
            };

            const targetTask = {
                id: '456',
                baslik: 'Target Task'
            };

            const dragData = {
                type: DragDataType.Task,
                task: sourceTask
            };

            mockDataTransfer.get.returns({
                value: dragData
            });

            sandbox.stub(vscode.window, 'showQuickPick').resolves({
                value: 'make_subtask'
            });

            const errorStub = sandbox.stub(vscode.window, 'showErrorMessage');
            
            mockMcpClient.callTool.rejects(new Error('dairesel bağımlılık'));

            const target = {
                task: targetTask
            };

            try {
                await controller.handleDrop(target, mockDataTransfer, {});
            } catch (e) {
                // Expected error
            }

            assert.ok(errorStub.calledWith('Bu işlem dairesel bağımlılık oluşturur!'));
        });

        test('should handle same project requirement error', async () => {
            const sourceTask = {
                id: '123',
                baslik: 'Source Task',
                proje_id: 'proj1'
            };

            const targetTask = {
                id: '456',
                baslik: 'Target Task',
                proje_id: 'proj2'
            };

            const dragData = {
                type: DragDataType.Task,
                task: sourceTask
            };

            mockDataTransfer.get.returns({
                value: dragData
            });

            sandbox.stub(vscode.window, 'showQuickPick').resolves({
                value: 'make_subtask'
            });

            const errorStub = sandbox.stub(vscode.window, 'showErrorMessage');
            
            mockMcpClient.callTool.rejects(new Error('aynı projede olmalı'));

            const target = {
                task: targetTask
            };

            try {
                await controller.handleDrop(target, mockDataTransfer, {});
            } catch (e) {
                // Expected error
            }

            assert.ok(errorStub.calledWith('Alt görev ve üst görev aynı projede olmalı!'));
        });
    });

    suite('canDrop Tests', () => {
        test('should allow drop on task when parent change is enabled', () => {
            const dragData = {
                type: DragDataType.Task,
                task: { id: '123' }
            };

            mockDataTransfer.get.returns({
                value: dragData
            });

            const target = {
                task: { id: '456' }
            };

            const canDrop = controller.canDrop(target, mockDataTransfer);
            assert.ok(canDrop);
        });

        test('should allow drop on empty area for tasks with parent', () => {
            const dragData = {
                type: DragDataType.Task,
                task: { 
                    id: '123',
                    parent_id: '456'
                }
            };

            mockDataTransfer.get.withArgs(DragDataType.Task).returns({
                value: dragData
            });

            const canDrop = controller.canDrop(undefined, mockDataTransfer);
            assert.ok(canDrop);
        });

        test('should not allow drop on empty area for root tasks', () => {
            const dragData = {
                type: DragDataType.Task,
                task: { 
                    id: '123',
                    parent_id: null
                }
            };

            mockDataTransfer.get.withArgs(DragDataType.Task).returns({
                value: dragData
            });

            const canDrop = controller.canDrop(undefined, mockDataTransfer);
            assert.strictEqual(canDrop, false);
        });
    });

    suite('Multiple Tasks Drop Tests', () => {
        test('should handle multiple tasks drop on empty area', async () => {
            const tasks = [
                { id: '1', baslik: 'Task 1', parent_id: '999' },
                { id: '2', baslik: 'Task 2', parent_id: '999' },
                { id: '3', baslik: 'Task 3', parent_id: '888' }
            ];

            const dragData = {
                type: DragDataType.Tasks,
                tasks: tasks
            };

            mockDataTransfer.get.returns({
                value: dragData
            });

            mockMcpClient.callTool.resolves();

            const progressStub = sandbox.stub();
            sandbox.stub(vscode.window, 'withProgress').callsFake(async (options, task) => {
                await task({ report: progressStub });
            });

            await controller.handleDrop(undefined, mockDataTransfer, {});

            // Should call gorev_ust_degistir for each task with parent
            assert.strictEqual(mockMcpClient.callTool.callCount, 3);
            tasks.forEach((task, index) => {
                assert.ok(mockMcpClient.callTool.getCall(index).calledWith('gorev_ust_degistir', {
                    id: task.id,
                    yeni_parent_id: ''
                }));
            });
        });
    });
});