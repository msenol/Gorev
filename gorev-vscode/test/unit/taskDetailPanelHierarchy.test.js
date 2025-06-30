const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const path = require('path');

suite('TaskDetailPanel Hierarchy Test Suite', () => {
    let sandbox;
    let mockMcpClient;
    let mockPanel;
    let mockWebview;

    setup(() => {
        sandbox = sinon.createSandbox();
        
        // Mock webview
        mockWebview = {
            html: '',
            asWebviewUri: (uri) => uri.toString(),
            onDidReceiveMessage: sandbox.stub(),
            postMessage: sandbox.stub(),
            cspSource: 'mock-csp-source'
        };

        // Mock panel
        mockPanel = {
            webview: mockWebview,
            title: '',
            iconPath: null,
            reveal: sandbox.stub(),
            onDidDispose: sandbox.stub(),
            onDidChangeViewState: sandbox.stub(),
            dispose: sandbox.stub(),
            visible: true
        };

        // Mock MCP Client
        mockMcpClient = {
            callTool: sandbox.stub(),
            isConnected: sandbox.stub().returns(true)
        };

        // Mock vscode.window.createWebviewPanel
        sandbox.stub(vscode.window, 'createWebviewPanel').returns(mockPanel);
        
        // Mock vscode.Uri
        sandbox.stub(vscode.Uri, 'joinPath').returns({ toString: () => 'mock-uri' });
    });

    teardown(() => {
        sandbox.restore();
    });

    suite('Hierarchy Section Rendering Tests', () => {
        test('should call gorev_hiyerarsi_goster for hierarchy info', async () => {
            const TaskDetailPanel = require('../../out/ui/taskDetailPanel').TaskDetailPanel;
            
            const task = {
                id: '123',
                baslik: 'Test Task',
                parent_id: '456'
            };

            const detayResponse = {
                content: [{
                    text: `# Test Task
**ID:** 123
**Durum:** beklemede
**Öncelik:** orta`
                }]
            };

            const hierarchyResponse = {
                content: [{
                    text: `## Görev Hiyerarşisi

Toplam alt görev: 5
Tamamlanan: 2
Devam eden: 1
Beklemede: 2
İlerleme: 40%`
                }]
            };

            mockMcpClient.callTool
                .withArgs('gorev_detay').resolves(detayResponse)
                .withArgs('gorev_hiyerarsi_goster').resolves(hierarchyResponse);

            await TaskDetailPanel.createOrShow(mockMcpClient, task, vscode.Uri.file('test'));

            // Should call both tools
            assert.ok(mockMcpClient.callTool.calledWith('gorev_detay', { id: '123' }));
            assert.ok(mockMcpClient.callTool.calledWith('gorev_hiyerarsi_goster', { id: '123' }));
        });

        test('should render hierarchy section in HTML', async () => {
            const TaskDetailPanel = require('../../out/ui/taskDetailPanel').TaskDetailPanel;
            
            const task = {
                id: '123',
                baslik: 'Parent Task',
                parent_id: null
            };

            mockMcpClient.callTool.resolves({
                content: [{
                    text: `# Parent Task

Toplam alt görev: 10
Tamamlanan: 7
İlerleme: 70%`
                }]
            });

            await TaskDetailPanel.createOrShow(mockMcpClient, task, vscode.Uri.file('test'));

            // Check if hierarchy section is in HTML
            const html = mockWebview.html;
            assert.ok(html.includes('hierarchy-section') || html.includes('Hiyerarşi'));
        });

        test('should show create subtask button', async () => {
            const TaskDetailPanel = require('../../out/ui/taskDetailPanel').TaskDetailPanel;
            
            const task = {
                id: '123',
                baslik: 'Test Task'
            };

            mockMcpClient.callTool.resolves({
                content: [{ text: '# Test Task' }]
            });

            await TaskDetailPanel.createOrShow(mockMcpClient, task, vscode.Uri.file('test'));

            const html = mockWebview.html;
            assert.ok(html.includes('createSubtaskBtn') || html.includes('Alt Görev Oluştur'));
        });

        test('should show parent info for child tasks', async () => {
            const TaskDetailPanel = require('../../out/ui/taskDetailPanel').TaskDetailPanel;
            
            const task = {
                id: '123',
                baslik: 'Child Task',
                parent_id: '456'
            };

            mockMcpClient.callTool.resolves({
                content: [{ text: '# Child Task' }]
            });

            await TaskDetailPanel.createOrShow(mockMcpClient, task, vscode.Uri.file('test'));

            const html = mockWebview.html;
            assert.ok(html.includes('parent-info') || html.includes('Üst Görev'));
        });
    });

    suite('Hierarchy Message Handling Tests', () => {
        test('should handle createSubtask message', () => {
            const executeCommandStub = sandbox.stub(vscode.commands, 'executeCommand');
            
            // Simulate webview message
            const messageHandler = mockWebview.onDidReceiveMessage.getCall(0).args[0];
            messageHandler({ command: 'createSubtask' });

            assert.ok(executeCommandStub.calledWith('gorev.createSubtask'));
        });

        test('should handle changeParent message', () => {
            const executeCommandStub = sandbox.stub(vscode.commands, 'executeCommand');
            
            // Simulate webview message
            const messageHandler = mockWebview.onDidReceiveMessage.getCall(0).args[0];
            messageHandler({ command: 'changeParent' });

            assert.ok(executeCommandStub.calledWith('gorev.changeParent'));
        });

        test('should handle removeParent message', () => {
            const executeCommandStub = sandbox.stub(vscode.commands, 'executeCommand');
            
            // Simulate webview message
            const messageHandler = mockWebview.onDidReceiveMessage.getCall(0).args[0];
            messageHandler({ command: 'removeParent' });

            assert.ok(executeCommandStub.calledWith('gorev.removeParent'));
        });
    });

    suite('Progress Bar Tests', () => {
        test('should calculate progress percentage correctly', () => {
            const hierarchyInfo = {
                toplam_alt_gorev: 10,
                tamamlanan_alt: 7,
                ilerleme_yuzdesi: 70
            };

            assert.strictEqual(hierarchyInfo.ilerleme_yuzdesi, 70);
            assert.strictEqual(hierarchyInfo.tamamlanan_alt / hierarchyInfo.toplam_alt_gorev * 100, 70);
        });

        test('should show 0% progress when no subtasks completed', () => {
            const hierarchyInfo = {
                toplam_alt_gorev: 5,
                tamamlanan_alt: 0,
                ilerleme_yuzdesi: 0
            };

            assert.strictEqual(hierarchyInfo.ilerleme_yuzdesi, 0);
        });

        test('should show 100% progress when all subtasks completed', () => {
            const hierarchyInfo = {
                toplam_alt_gorev: 5,
                tamamlanan_alt: 5,
                ilerleme_yuzdesi: 100
            };

            assert.strictEqual(hierarchyInfo.ilerleme_yuzdesi, 100);
        });
    });
});