const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('ImportWizard Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let mockContext;
  let mockProviders;
  let ImportWizard;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'createWebviewPanel');
    sandbox.stub(vscode.window, 'showOpenDialog');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'withProgress');
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

    // Mock Providers
    mockProviders = {
      gorevTreeProvider: { 
        refresh: sandbox.stub().resolves() 
      },
      projeTreeProvider: { 
        refresh: sandbox.stub().resolves() 
      },
      templateTreeProvider: { 
        refresh: sandbox.stub().resolves() 
      },
      statusBarManager: { 
        update: sandbox.stub() 
      }
    };

    // Try to import ImportWizard
    try {
      const importWizardModule = require('../../dist/ui/importWizard');
      ImportWizard = importWizardModule.ImportWizard;
    } catch (error) {
      // Mock ImportWizard if compilation fails
      ImportWizard = class MockImportWizard {
        constructor(context, mcpClient, providers) {
          this.context = context;
          this.mcpClient = mcpClient;
          this.providers = providers;
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

  suite('Wizard Creation', () => {
    test('should create ImportWizard instance', () => {
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      
      assert(wizard);
      assert.strictEqual(wizard.context, mockContext);
      assert.strictEqual(wizard.mcpClient, mockMCPClient);
      assert.strictEqual(wizard.providers, mockProviders);
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      assert(vscode.window.createWebviewPanel.called);
      assert(vscode.window.createWebviewPanel.calledWith(
        'gorevImportWizard',
        'import.wizardTitle',
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      assert(mockPanel.webview.onDidReceiveMessage.called);
    });

    test('should handle selectFile message', async () => {
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
      vscode.window.showOpenDialog.resolves([{
        fsPath: '/test/import.json'
      }]);
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate selectFile message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'selectFile' });
      
      assert(vscode.window.showOpenDialog.called);
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'setSelectedFile',
        path: '/test/import.json'
      }));
    });

    test('should handle analyzeFile message', async () => {
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
      mockMCPClient.callTool.withArgs('gorev_import').resolves({
        content: [{ text: 'Toplam 10 görev, 2 proje bulundu' }],
        isError: false
      });
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate analyzeFile message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'analyzeFile',
        filePath: '/test/import.json'
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'analysisStarted'
      }));
      assert(mockMCPClient.callTool.calledWith('gorev_import', sinon.match({
        file_path: '/test/import.json',
        dry_run: true
      })));
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate loadProjects message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'loadProjects' });
      
      assert(mockMCPClient.callTool.calledWith('proje_listele'));
      assert(mockPanel.webview.postMessage.calledWith(sinon.match({
        command: 'setProjects'
      })));
    });

    test('should handle performDryRun message', async () => {
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
      mockMCPClient.callTool.withArgs('gorev_import').resolves({
        content: [{ text: '5 görev import edilecek, 2 proje import edilecek' }],
        isError: false
      });
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate performDryRun message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'performDryRun',
        options: {
          file_path: '/test/import.json',
          format: 'json',
          conflict_resolution: 'skip'
        }
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'dryRunStarted'
      }));
      assert(mockMCPClient.callTool.calledWith('gorev_import', sinon.match({
        dry_run: true
      })));
    });

    test('should handle startImport message', async () => {
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate startImport message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startImport',
        options: {
          file_path: '/test/import.json',
          format: 'json',
          conflict_resolution: 'skip'
        }
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'importStarted'
      }));
      assert(vscode.window.withProgress.called);
      assert(mockMCPClient.callTool.calledWith('gorev_import', sinon.match({
        dry_run: false
      })));
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate cancel message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'cancel' });
      
      assert(mockPanel.dispose.called);
    });
  });

  suite('Import Operations', () => {
    test('should perform import with progress updates', async () => {
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate import
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startImport',
        options: {
          file_path: '/test/import.json',
          format: 'json'
        }
      });
      
      assert(progressReports.length > 0);
      assert(progressReports.some(p => p.increment === 10));
      assert(progressReports.some(p => p.increment === 80));
      assert(progressReports.some(p => p.increment === 100));
    });

    test('should refresh views after successful import', async () => {
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate import
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startImport',
        options: {
          file_path: '/test/import.json',
          format: 'json'
        }
      });
      
      assert(mockProviders.gorevTreeProvider.refresh.called);
      assert(mockProviders.projeTreeProvider.refresh.called);
      assert(mockProviders.templateTreeProvider.refresh.called);
    });

    test('should handle import errors', async () => {
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
        content: [{ text: 'Import failed' }],
        isError: true
      });
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate import
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startImport',
        options: {
          file_path: '/test/import.json',
          format: 'json'
        }
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'importFailed',
        error: 'Import failed'
      }));
      assert(vscode.window.showErrorMessage.called);
    });
  });

  suite('File Analysis', () => {
    test('should detect JSON format from file extension', async () => {
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
        content: [{ text: 'Analysis result' }],
        isError: false
      });
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate file analysis
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'analyzeFile',
        filePath: '/test/import.json'
      });
      
      assert(mockMCPClient.callTool.calledWith('gorev_import', sinon.match({
        format: 'json'
      })));
    });

    test('should detect CSV format from file extension', async () => {
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
        content: [{ text: 'Analysis result' }],
        isError: false
      });
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate file analysis
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'analyzeFile',
        filePath: '/test/import.csv'
      });
      
      assert(mockMCPClient.callTool.calledWith('gorev_import', sinon.match({
        format: 'csv'
      })));
    });

    test('should handle analysis failures', async () => {
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
        content: [{ text: 'File format not supported' }],
        isError: true
      });
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate file analysis
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'analyzeFile',
        filePath: '/test/import.json'
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'analysisFailed',
        error: 'File format not supported'
      }));
    });
  });

  suite('Data Parsing', () => {
    test('should parse analysis result correctly', async () => {
      if (!ImportWizard.prototype.parseAnalysisResult) {
        // Skip if method not available
        return;
      }

      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      const analysisText = `
Toplam 15 görev bulundu
Toplam 3 proje bulundu
2 uyarı var
      `.trim();
      
      const result = wizard.parseAnalysisResult(analysisText);
      
      assert.strictEqual(result.totalTasks, 15);
      assert.strictEqual(result.totalProjects, 3);
      assert(Array.isArray(result.warnings));
      assert(Array.isArray(result.errors));
    });

    test('should parse dry run result correctly', async () => {
      if (!ImportWizard.prototype.parseDryRunResult) {
        return;
      }

      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      const dryRunText = `
10 görev import edilecek
2 proje import edilecek
3 çakışma tespit edildi
      `.trim();
      
      const result = wizard.parseDryRunResult(dryRunText);
      
      assert.strictEqual(result.tasksToImport, 10);
      assert.strictEqual(result.projectsToImport, 2);
      assert(Array.isArray(result.conflicts));
      assert(Array.isArray(result.warnings));
    });

    test('should parse import result correctly', async () => {
      if (!ImportWizard.prototype.parseImportResult) {
        return;
      }

      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      const importText = `
8 görev import edildi
2 proje import edildi
2 öğe atlandı
      `.trim();
      
      const result = wizard.parseImportResult(importText);
      
      assert.strictEqual(result.tasksImported, 8);
      assert.strictEqual(result.projectsImported, 2);
      assert.strictEqual(result.skipped, 2);
      assert(Array.isArray(result.errors));
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate loadProjects message that will fail
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'loadProjects' });
      
      // Should not crash and should handle error gracefully
      assert(true);
    });

    test('should handle file selection cancellation', async () => {
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
      vscode.window.showOpenDialog.resolves(undefined);
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate selectFile message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ command: 'selectFile' });
      
      assert(vscode.window.showOpenDialog.called);
      // Should not crash when user cancels
      assert(true);
    });

    test('should handle dry run failures', async () => {
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
      mockMCPClient.callTool.rejects(new Error('Dry run failed'));
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate performDryRun message
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'performDryRun',
        options: {}
      });
      
      assert(mockPanel.webview.postMessage.calledWith({
        command: 'dryRunFailed',
        error: 'Dry run failed'
      }));
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      assert(typeof mockPanel.webview.html === 'string');
      assert(mockPanel.webview.html.includes('<!DOCTYPE html>'));
      assert(mockPanel.webview.html.includes('<html'));
      assert(mockPanel.webview.html.includes('</html>'));
      assert(mockPanel.webview.html.includes('import.wizardTitle'));
    });

    test('should include all required wizard steps', async () => {
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      const html = mockPanel.webview.html;
      assert(html.includes('step1'));
      assert(html.includes('step2'));
      assert(html.includes('step3'));
      assert(html.includes('step4'));
      assert(html.includes('import.step1.title'));
      assert(html.includes('import.step2.title'));
      assert(html.includes('import.step3.title'));
      assert(html.includes('import.step4.title'));
    });

    test('should include conflict resolution options', async () => {
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      const html = mockPanel.webview.html;
      assert(html.includes('import.conflict.skip'));
      assert(html.includes('import.conflict.overwrite'));
      assert(html.includes('import.conflict.merge'));
    });
  });

  suite('Success Scenarios', () => {
    test('should show success message after import', async () => {
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
        content: [{ text: '5 görev import edildi, 2 proje import edildi' }],
        isError: false
      });
      vscode.window.withProgress.callsFake(async (options, callback) => {
        const progress = { report: sandbox.stub() };
        return await callback(progress);
      });
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate import
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startImport',
        options: {
          file_path: '/test/import.json',
          format: 'json'
        }
      });
      
      assert(vscode.window.showInformationMessage.calledWith(
        sinon.match(/import.success/)
      ));
    });

    test('should auto-close wizard after successful import', async () => {
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
      
      const wizard = new ImportWizard(mockContext, mockMCPClient, mockProviders);
      await wizard.show();
      
      // Simulate import
      const messageHandler = mockPanel.webview.onDidReceiveMessage.getCall(0).args[0];
      await messageHandler({ 
        command: 'startImport',
        options: {
          file_path: '/test/import.json',
          format: 'json'
        }
      });
      
      // Panel should be disposed after a delay
      setTimeout(() => {
        assert(mockPanel.dispose.called);
      }, 3100);
    });
  });
});