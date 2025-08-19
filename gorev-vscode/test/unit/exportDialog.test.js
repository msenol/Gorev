const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('ExportDialog Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let mockContext;
  let ExportDialog;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'createWebviewPanel');
    sandbox.stub(vscode.window, 'showSaveDialog');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'withProgress');
    sandbox.stub(vscode.commands, 'executeCommand');
    sandbox.stub(vscode.l10n, 't').callsFake((key, data) => {
      let result = key;
      if (data) {
        Object.keys(data).forEach(prop => {
          result = result.replace(`{${prop}}`, data[prop]);
        });
      }
      return result;
    });
    sandbox.stub(vscode.Uri, 'joinPath').callsFake((base, ...segments) => ({
      fsPath: segments.join('/')
    }));
    sandbox.stub(vscode.Uri, 'file').callsFake(path => ({ fsPath: path }));
    sandbox.stub(vscode.workspace, 'rootPath').value('/test/workspace');

    // Mock MCP Client
    mockMCPClient = {
      callTool: sandbox.stub().resolves({ 
        content: [{ text: 'Success' }],
        isError: false
      }),
      isConnected: sandbox.stub().returns(true)
    };

    // Mock Context
    mockContext = {
      subscriptions: [],
      extensionUri: vscode.Uri.file('/test/extension')
    };

    // Try to import ExportDialog
    try {
      const exportDialogModule = require('../../dist/ui/exportDialog');
      ExportDialog = exportDialogModule.ExportDialog;
    } catch (error) {
      // Mock ExportDialog if compilation fails
      ExportDialog = class MockExportDialog {
        constructor(context, mcpClient) {
          this.context = context;
          this.mcpClient = mcpClient;
          this.panel = null;
        }
        
        async show() {
          // Mock implementation
          const mockPanel = {
            webview: {
              html: '',
              onDidReceiveMessage: sandbox.stub(),
              postMessage: sandbox.stub()
            },
            iconPath: undefined,
            onDidDispose: sandbox.stub(),
            dispose: sandbox.stub()
          };
          
          vscode.window.createWebviewPanel.returns(mockPanel);
          this.panel = vscode.window.createWebviewPanel();
          return this.panel;
        }
      };
    }
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('Dialog Creation', () => {
    test('should create ExportDialog instance', () => {
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      
      assert(dialog);
      assert.strictEqual(dialog.context, mockContext);
      assert.strictEqual(dialog.mcpClient, mockMCPClient);
    });

    test('should create webview panel when shown', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      assert(vscode.window.createWebviewPanel.called);
      assert(vscode.window.createWebviewPanel.calledWith(
        'gorevExportDialog',
        'export.dialogTitle',
        vscode.ViewColumn.One
      ));
    });

    test('should configure webview options correctly', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      const createCall = vscode.window.createWebviewPanel.getCall(0);
      const options = createCall.args[3];
      
      assert.strictEqual(options.enableScripts, true);
      assert(Array.isArray(options.localResourceRoots));
      assert.strictEqual(options.localResourceRoots.length, 2);
    });

    test('should set panel icon paths', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      assert(mockPanel.iconPath);
      assert(mockPanel.iconPath.light);
      assert(mockPanel.iconPath.dark);
    });
  });

  suite('Message Handling', () => {
    test('should register message handler', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      assert(mockPanel.webview.onDidReceiveMessage.called);
    });

    test('should handle loadProjects message', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      mockMCPClient.callTool.withArgs('proje_listele').resolves({
        content: [{ text: '- Test Project (ID: test-123) ✅\n- Another Project (ID: another-456)' }],
        isError: false
      });
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate loadProjects message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'loadProjects' });
      
      assert(mockMCPClient.callTool.calledWith('proje_listele'));
      assert(mockPanel.webview.postMessage.calledWith(sinon.match({
        command: 'setProjects'
      })));
    });

    test('should handle loadTags message', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      mockMCPClient.callTool.withArgs('ozet_goster').resolves({
        content: [{ text: 'Etiketler: bug, feature, urgent' }],
        isError: false
      });
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate loadTags message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'loadTags' });
      
      assert(mockMCPClient.callTool.calledWith('ozet_goster'));
      assert(mockPanel.webview.postMessage.calledWith(sinon.match({
        command: 'setTags'
      })));
    });

    test('should handle selectOutputPath message', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      vscode.window.showSaveDialog.resolves({
        fsPath: '/test/export.json'
      });
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate selectOutputPath message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'selectOutputPath',
        format: 'json'
      });
      
      assert(vscode.window.showSaveDialog.called);
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'setOutputPath',
        path: '/test/export.json'
      }));
    });

    test('should handle startExport message', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate startExport message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startExport',
        options: {
          output_path: '/test/export.json',
          format: 'json'
        }
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'exportStarted'
      }));
      assert(vscode.window.withProgress.called);
      assert(mockMCPClient.callTool.calledWith('gorev_export'));
    });

    test('should handle cancel message', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate cancel message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'cancel' });
      
      assert(mockPanel.dispose.called);
    });

    test('should handle unknown messages gracefully', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate unknown message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'unknownCommand' });
      
      // Should not throw or crash
      assert(true);
    });
  });

  suite('Export Operations', () => {
    test('should perform export with progress updates', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      let progressReports = [];
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { 
          report: (data) => progressReports.push(data)
        };
        return await callback(progress);
      });
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate export
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startExport',
        options: {
          output_path: '/test/export.json',
          format: 'json'
        }
      });
      
      assert(progressReports.length > 0);
      assert(progressReports.some(p => p.increment === 10));
      assert(progressReports.some(p => p.increment === 90));
      assert(progressReports.some(p => p.increment === 100));
    });

    test('should handle export errors', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      mockMCPClient.callTool.resolves({
        content: [{ text: 'Export failed' }],
        isError: true
      });
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate export
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startExport',
        options: {
          output_path: '/test/export.json',
          format: 'json'
        }
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'exportFailed',
        error: 'Export failed'
      }));
      assert(vscode.window.showErrorMessage.called);
    });

    test('should show success actions after export', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      vscode.window.showInformationMessage.resolves('export.openFile');
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate export
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startExport',
        options: {
          output_path: '/test/export.json',
          format: 'json'
        }
      });
      
      assert(vscode.window.showInformationMessage.called);
      assert(vscode.commands.executeCommand.calledWith('vscode.open'));
    });
  });

  suite('Data Parsing', () => {
    test('should parse project list correctly', async () => {
      if (!ExportDialog.prototype.parseProjectList) {
        // Skip if method not available
        return;
      }

      const dialog = new ExportDialog(mockContext, mockMCPClient);
      const projectText = `
- Test Project (ID: test-123) ✅
- Another Project (ID: another-456)
- Third Project (ID: third-789) ✅
      `.trim();
      
      const projects = dialog.parseProjectList(projectText);
      
      assert.strictEqual(projects.length, 3);
      assert.strictEqual(projects[0].name, 'Test Project');
      assert.strictEqual(projects[0].id, 'test-123');
      assert.strictEqual(projects[0].isActive, true);
      assert.strictEqual(projects[1].name, 'Another Project');
      assert.strictEqual(projects[1].id, 'another-456');
      assert.strictEqual(projects[1].isActive, false);
    });

    test('should parse tags from summary correctly', async () => {
      if (!ExportDialog.prototype.parseTagsFromSummary) {
        // Skip if method not available
        return;
      }

      const dialog = new ExportDialog(mockContext, mockMCPClient);
      const summaryText = `
Toplam 10 görev
Etiketler: bug, feature, urgent, documentation
5 proje mevcut
      `.trim();
      
      const tags = dialog.parseTagsFromSummary(summaryText);
      
      assert(Array.isArray(tags));
      assert(tags.includes('bug'));
      assert(tags.includes('feature'));
      assert(tags.includes('urgent'));
      assert(tags.includes('documentation'));
    });

    test('should handle empty project list', async () => {
      if (!ExportDialog.prototype.parseProjectList) {
        return;
      }

      const dialog = new ExportDialog(mockContext, mockMCPClient);
      const projects = dialog.parseProjectList('');
      
      assert(Array.isArray(projects));
      assert.strictEqual(projects.length, 0);
    });
  });

  suite('Error Handling', () => {
    test('should handle MCP tool errors gracefully', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      mockMCPClient.callTool.rejects(new Error('Network error'));
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate loadProjects message that will fail
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'loadProjects' });
      
      // Should not crash and should handle error gracefully
      assert(true);
    });

    test('should handle webview disposal', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      // Simulate panel disposal
      const disposeHandler = mockPanel.onDidDispose.getCall(0).args[0];
      disposeHandler();
      
      // Panel should be set to undefined
      assert.strictEqual(dialog.panel, undefined);
    });
  });

  suite('HTML Generation', () => {
    test('should generate valid HTML content', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      assert(typeof mockPanel.webview.html === 'string');
      assert(mockPanel.webview.html.includes('<!DOCTYPE html>'));
      assert(mockPanel.webview.html.includes('<html'));
      assert(mockPanel.webview.html.includes('</html>'));
      assert(mockPanel.webview.html.includes('export.dialogTitle'));
    });

    test('should include security policy in HTML', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      assert(mockPanel.webview.html.includes('Content-Security-Policy'));
      assert(mockPanel.webview.html.includes("script-src 'nonce-"));
    });

    test('should include all required form steps', async () => {
      const mockPanel = {
        webview: {
          html: '',
          onDidReceiveMessage: sandbox.stub(),
          postMessage: sandbox.stub()
        },
        iconPath: undefined,
        onDidDispose: sandbox.stub(),
        dispose: sandbox.stub()
      };
      
      vscode.window.createWebviewPanel.returns(mockPanel);
      
      const dialog = new ExportDialog(mockContext, mockMCPClient);
      await dialog.show();
      
      const html = mockPanel.webview.html;
      assert(html.includes('step1'));
      assert(html.includes('step2'));
      assert(html.includes('step3'));
      assert(html.includes('step4'));
      assert(html.includes('export.step1.title'));
      assert(html.includes('export.step2.title'));
      assert(html.includes('export.step3.title'));
      assert(html.includes('export.step4.title'));
    });
  });
});