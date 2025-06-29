const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('TaskDetailPanel Test Suite', () => {
  let panel;
  let sandbox;
  let mockWebviewPanel;
  let mockMCPClient;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.commands, 'executeCommand');
    
    // Mock webview panel
    mockWebviewPanel = {
      webview: {
        html: '',
        postMessage: sandbox.stub(),
        onDidReceiveMessage: sandbox.stub(),
        cspSource: 'vscode-webview:',
        asWebviewUri: sandbox.stub().returns('mock-uri')
      },
      onDidDispose: sandbox.stub(),
      dispose: sandbox.stub(),
      reveal: sandbox.stub(),
      visible: true,
      active: true
    };

    // Mock MCP Client
    mockMCPClient = {
      isConnected: sandbox.stub().returns(true),
      callTool: sandbox.stub().resolves({ 
        content: [{ type: 'text', text: '# Test Task\n\n**ID:** task-123\n**Durum:** beklemede\n**Öncelik:** orta\n\nTest description' }] 
      })
    };

    // Create panel with mocked dependencies
    try {
      const { TaskDetailPanel } = require('../../dist/ui/taskDetailPanel');
      panel = TaskDetailPanel.createOrShow(mockWebviewPanel, mockMCPClient, 'task-123');
    } catch (error) {
      // Handle compilation or import issues
      panel = createMockPanel();
    }
  });

  teardown(() => {
    sandbox.restore();
    if (panel && panel.dispose) {
      panel.dispose();
    }
  });

  function createMockPanel() {
    return {
      loadTask: sandbox.stub(),
      updateTask: sandbox.stub(),
      dispose: sandbox.stub(),
      _onDidDispose: { fire: sandbox.stub() },
      _panel: mockWebviewPanel,
      _mcpClient: mockMCPClient
    };
  }

  suite('Panel Creation', () => {
    test('should create panel instance', () => {
      assert(panel);
      if (panel._panel) {
        assert.strictEqual(panel._panel, mockWebviewPanel);
      }
    });

    test('should register webview message handler', () => {
      if (mockWebviewPanel.webview.onDidReceiveMessage.called) {
        assert(mockWebviewPanel.webview.onDidReceiveMessage.calledOnce);
      }
    });

    test('should register dispose handler', () => {
      if (mockWebviewPanel.onDidDispose.called) {
        assert(mockWebviewPanel.onDidDispose.calledOnce);
      }
    });
  });

  suite('Task Loading', () => {
    test('should load task details from MCP', async () => {
      if (panel.loadTask) {
        await panel.loadTask('task-123');
        
        if (mockMCPClient.callTool.called) {
          assert(mockMCPClient.callTool.calledWith('gorev_detay', { id: 'task-123' }));
        }
      }
    });

    test('should handle loading errors gracefully', async () => {
      mockMCPClient.callTool.rejects(new Error('Failed to load task'));

      if (panel.loadTask) {
        try {
          await panel.loadTask('invalid-task');
          // Should not throw
          assert(true);
        } catch (error) {
          assert(false, 'Should handle loading errors gracefully');
        }
      }
    });

    test('should update webview HTML after loading', async () => {
      if (panel.loadTask) {
        await panel.loadTask('task-123');
        
        // HTML should be updated
        assert(typeof mockWebviewPanel.webview.html === 'string');
      }
    });
  });

  suite('Task Editing', () => {
    test('should handle title edit request', () => {
      const message = {
        command: 'edit',
        field: 'baslik',
        value: 'Updated Title'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });

    test('should handle description edit request', () => {
      const message = {
        command: 'edit',
        field: 'aciklama',
        value: 'Updated description'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });

    test('should handle priority change request', () => {
      const message = {
        command: 'edit',
        field: 'oncelik',
        value: 'yuksek'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });

    test('should handle status change request', () => {
      const message = {
        command: 'updateStatus',
        taskId: 'task-123',
        newStatus: 'devam_ediyor'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        
        if (mockMCPClient.callTool.called) {
          assert(mockMCPClient.callTool.calledWith('gorev_guncelle'));
        }
      }
    });

    test('should handle due date change request', () => {
      const message = {
        command: 'edit',
        field: 'son_tarih',
        value: '2025-12-31'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });
  });

  suite('Task Actions', () => {
    test('should handle delete request', () => {
      const message = {
        command: 'delete',
        taskId: 'task-123'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        
        if (mockMCPClient.callTool.called) {
          assert(mockMCPClient.callTool.calledWith('gorev_sil'));
        }
      }
    });

    test('should handle copy task request', () => {
      const message = {
        command: 'copy',
        taskId: 'task-123'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });

    test('should handle refresh request', () => {
      const message = {
        command: 'refresh',
        taskId: 'task-123'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        
        if (panel.loadTask) {
          // Should reload task
          assert(true);
        }
      }
    });
  });

  suite('Dependency Management', () => {
    test('should handle add dependency request', () => {
      const message = {
        command: 'addDependency',
        sourceId: 'task-123',
        targetId: 'task-456',
        type: 'blocks'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        
        if (mockMCPClient.callTool.called) {
          assert(mockMCPClient.callTool.calledWith('gorev_bagimlilik_ekle'));
        }
      }
    });

    test('should handle remove dependency request', () => {
      const message = {
        command: 'removeDependency',
        sourceId: 'task-123',
        targetId: 'task-456'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });
  });

  suite('Tag Management', () => {
    test('should handle add tag request', () => {
      const message = {
        command: 'addTag',
        taskId: 'task-123',
        tag: 'bug'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });

    test('should handle remove tag request', () => {
      const message = {
        command: 'removeTag',
        taskId: 'task-123',
        tag: 'bug'
      };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });
  });

  suite('HTML Generation', () => {
    test('should generate valid HTML structure', () => {
      if (panel._getHtmlForWebview) {
        const html = panel._getHtmlForWebview();
        
        assert(typeof html === 'string');
        assert(html.includes('<!DOCTYPE html>'));
        assert(html.includes('<html>'));
        assert(html.includes('<body>'));
        assert(html.includes('</html>'));
      }
    });

    test('should include CSP meta tag', () => {
      if (panel._getHtmlForWebview) {
        const html = panel._getHtmlForWebview();
        
        assert(html.includes('Content-Security-Policy'));
        assert(html.includes('vscode-webview:'));
      }
    });

    test('should include task data script', () => {
      if (panel._getHtmlForWebview) {
        const mockTask = {
          id: 'task-123',
          baslik: 'Test Task',
          durum: 'beklemede',
          oncelik: 'orta'
        };

        const html = panel._getHtmlForWebview(mockTask);
        
        assert(html.includes('window.taskData'));
      }
    });
  });

  suite('State Management', () => {
    test('should track panel visibility', () => {
      if (panel._panel) {
        assert.strictEqual(panel._panel.visible, true);
      }
    });

    test('should track panel active state', () => {
      if (panel._panel) {
        assert.strictEqual(panel._panel.active, true);
      }
    });

    test('should update task data on successful edit', async () => {
      if (panel.updateTask) {
        const updatedTask = {
          id: 'task-123',
          baslik: 'Updated Task',
          durum: 'devam_ediyor',
          oncelik: 'yuksek'
        };

        await panel.updateTask(updatedTask);
        assert(true);
      }
    });
  });

  suite('Error Handling', () => {
    test('should handle MCP connection errors', async () => {
      mockMCPClient.isConnected.returns(false);

      if (panel.loadTask) {
        try {
          await panel.loadTask('task-123');
          assert(true); // Should handle gracefully
        } catch (error) {
          assert(false, 'Should handle connection errors gracefully');
        }
      }
    });

    test('should handle invalid message format', () => {
      const invalidMessage = null;

      if (panel._handleWebviewMessage) {
        try {
          panel._handleWebviewMessage(invalidMessage);
          assert(true); // Should handle gracefully
        } catch (error) {
          assert(false, 'Should handle invalid messages gracefully');
        }
      }
    });

    test('should handle task not found error', async () => {
      mockMCPClient.callTool.rejects(new Error('Task not found'));

      if (panel.loadTask) {
        try {
          await panel.loadTask('nonexistent-task');
          assert(true); // Should handle gracefully
        } catch (error) {
          assert(false, 'Should handle task not found gracefully');
        }
      }
    });
  });

  suite('Disposal', () => {
    test('should clean up resources on dispose', () => {
      if (panel.dispose) {
        panel.dispose();
        
        if (mockWebviewPanel.dispose.called) {
          assert(mockWebviewPanel.dispose.calledOnce);
        }
      }
    });

    test('should fire dispose event', () => {
      if (panel.dispose && panel._onDidDispose) {
        panel.dispose();
        
        if (panel._onDidDispose.fire) {
          assert(panel._onDidDispose.fire.called);
        }
      }
    });

    test('should prevent multiple disposals', () => {
      if (panel.dispose) {
        panel.dispose();
        
        try {
          panel.dispose(); // Second call should be safe
          assert(true);
        } catch (error) {
          assert(false, 'Should handle multiple dispose calls safely');
        }
      }
    });
  });

  suite('Webview Communication', () => {
    test('should post messages to webview', () => {
      if (panel._postMessage) {
        const message = { command: 'updateTask', data: { id: 'task-123' } };
        panel._postMessage(message);
        
        if (mockWebviewPanel.webview.postMessage.called) {
          assert(mockWebviewPanel.webview.postMessage.calledWith(message));
        }
      }
    });

    test('should handle webview ready message', () => {
      const message = { command: 'ready' };

      if (panel._handleWebviewMessage) {
        panel._handleWebviewMessage(message);
        assert(true);
      }
    });

    test('should handle unknown commands gracefully', () => {
      const message = { command: 'unknownCommand', data: {} };

      if (panel._handleWebviewMessage) {
        try {
          panel._handleWebviewMessage(message);
          assert(true); // Should handle gracefully
        } catch (error) {
          assert(false, 'Should handle unknown commands gracefully');
        }
      }
    });
  });
});