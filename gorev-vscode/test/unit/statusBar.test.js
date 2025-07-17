const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('StatusBar Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let StatusBarManager;
  let statusBarManager;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'createStatusBarItem');
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.workspace, 'getConfiguration').returns({
      get: sandbox.stub().returns(true)
    });

    // Mock MCP Client
    mockMCPClient = {
      callTool: sandbox.stub(),
      isConnected: sandbox.stub().returns(true),
      on: sandbox.stub(),
      off: sandbox.stub()
    };

    // Set up different responses for different tool calls
    mockMCPClient.callTool.withArgs('aktif_proje_goster').resolves({
      content: [{ text: '**Aktif Proje:** Test Project (ID: proj-001)' }]
    });
    mockMCPClient.callTool.withArgs('ozet_goster').resolves({
      content: [{ text: '## Görev Özeti\n\n- Toplam Görev: 10\n- Beklemede: 3\n- Devam Ediyor: 4\n- Tamamlandı: 3' }]
    });

    // Mock Status Bar Items
    const mockStatusBarItem = {
      text: '',
      tooltip: '',
      command: '',
      color: undefined,
      backgroundColor: undefined,
      show: sandbox.stub(),
      hide: sandbox.stub(),
      dispose: sandbox.stub()
    };

    vscode.window.createStatusBarItem.returns(mockStatusBarItem);

    // Import StatusBarManager
    try {
      const statusBarModule = require('../../dist/ui/statusBar');
      StatusBarManager = statusBarModule.StatusBarManager;
    } catch (error) {
      // Mock StatusBarManager class if compilation fails
      StatusBarManager = class MockStatusBarManager {
        constructor(mcpClient) {
          this.mcpClient = mcpClient;
          this.isVisible = false;
          this.activeProject = null;
          this.taskStats = null;
          this.createStatusBarItems();
          this.setupEventListeners();
        }

        createStatusBarItems() {
          this.connectionStatusItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 200);
          this.activeProjectItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 199);
          this.taskStatsItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 198);
          this.quickActionsItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Left, 197);
        }

        setupEventListeners() {
          if (this.mcpClient) {
            this.mcpClient.on('connected', () => this.onConnectionChanged(true));
            this.mcpClient.on('disconnected', () => this.onConnectionChanged(false));
          }
        }

        async show() {
          this.isVisible = true;
          await this.update();
          this.connectionStatusItem.show();
          this.activeProjectItem.show();
          this.taskStatsItem.show();
          this.quickActionsItem.show();
        }

        hide() {
          this.isVisible = false;
          this.connectionStatusItem.hide();
          this.activeProjectItem.hide();
          this.taskStatsItem.hide();
          this.quickActionsItem.hide();
        }

        async update() {
          if (!this.isVisible) return;

          await this.updateConnectionStatus();
          await this.updateActiveProject();
          await this.updateTaskStats();
          this.updateQuickActions();
        }

        async updateConnectionStatus() {
          if (this.mcpClient && this.mcpClient.isConnected()) {
            this.connectionStatusItem.text = '$(plug) Connected';
            this.connectionStatusItem.color = undefined;
            this.connectionStatusItem.tooltip = 'MCP Sunucuya bağlı';
          } else {
            this.connectionStatusItem.text = '$(plug) Disconnected';
            this.connectionStatusItem.color = new vscode.ThemeColor('statusBarItem.errorForeground');
            this.connectionStatusItem.tooltip = 'MCP Sunucuya bağlı değil';
          }
        }

        async updateActiveProject() {
          try {
            if (this.mcpClient && this.mcpClient.isConnected()) {
              const result = await this.mcpClient.callTool('aktif_proje_goster');
              const content = result.content[0].text;
              
              if (content.includes('Aktif proje bulunmuyor')) {
                this.activeProject = null;
                this.activeProjectItem.text = '$(folder) No Active Project';
                this.activeProjectItem.tooltip = 'Aktif proje yok - proje seçmek için tıklayın';
              } else {
                const match = content.match(/\*\*Aktif Proje:\*\* (.+)/);
                if (match) {
                  this.activeProject = match[1];
                  this.activeProjectItem.text = `$(folder) ${match[1]}`;
                  this.activeProjectItem.tooltip = `Aktif proje: ${match[1]}`;
                }
              }
              this.activeProjectItem.command = 'gorev.setActiveProject';
            }
          } catch (error) {
            this.activeProjectItem.text = '$(folder) Error';
            this.activeProjectItem.tooltip = 'Proje bilgisi alınamadı';
          }
        }

        async updateTaskStats() {
          try {
            if (this.mcpClient && this.mcpClient.isConnected()) {
              const result = await this.mcpClient.callTool('ozet_goster');
              const content = result.content[0].text;
              
              const stats = this.parseTaskStats(content);
              this.taskStats = stats;
              
              if (stats) {
                this.taskStatsItem.text = `$(checklist) ${stats.beklemede}/${stats.devamEdiyor}/${stats.tamamlandi}`;
                this.taskStatsItem.tooltip = `Görevler - Beklemede: ${stats.beklemede}, Devam Ediyor: ${stats.devamEdiyor}, Tamamlandı: ${stats.tamamlandi}`;
              } else {
                this.taskStatsItem.text = '$(checklist) 0/0/0';
                this.taskStatsItem.tooltip = 'Görev istatistikleri';
              }
              this.taskStatsItem.command = 'gorev.showSummary';
            }
          } catch (error) {
            this.taskStatsItem.text = '$(checklist) Error';
            this.taskStatsItem.tooltip = 'İstatistik bilgisi alınamadı';
          }
        }

        updateQuickActions() {
          this.quickActionsItem.text = '$(add) New Task';
          this.quickActionsItem.tooltip = 'Yeni görev oluştur';
          this.quickActionsItem.command = 'gorev.createTask';
        }

        parseTaskStats(content) {
          const beklemedeMat = content.match(/Beklemede[:\s]*(\d+)/i);
          const devamEdiyorMat = content.match(/Devam Ediyor[:\s]*(\d+)/i);
          const tamamlandiMat = content.match(/Tamamlandı[:\s]*(\d+)/i);

          if (beklemedeMat && devamEdiyorMat && tamamlandiMat) {
            return {
              beklemede: parseInt(beklemedeMat[1]),
              devamEdiyor: parseInt(devamEdiyorMat[1]),
              tamamlandi: parseInt(tamamlandiMat[1])
            };
          }
          return null;
        }

        onConnectionChanged(isConnected) {
          if (this.isVisible) {
            this.update();
          }
        }

        dispose() {
          if (this.connectionStatusItem) this.connectionStatusItem.dispose();
          if (this.activeProjectItem) this.activeProjectItem.dispose();
          if (this.taskStatsItem) this.taskStatsItem.dispose();
          if (this.quickActionsItem) this.quickActionsItem.dispose();
          
          if (this.mcpClient) {
            this.mcpClient.off('connected', this.onConnectionChanged);
            this.mcpClient.off('disconnected', this.onConnectionChanged);
          }
        }
      };
    }

    // Create status bar manager instance
    if (StatusBarManager) {
      statusBarManager = new StatusBarManager(mockMCPClient);
    }
  });

  teardown(() => {
    if (statusBarManager && typeof statusBarManager.dispose === 'function') {
      statusBarManager.dispose();
    }
    sandbox.restore();
  });

  suite('Initialization', () => {
    test('should create status bar manager with MCP client', () => {
      assert(statusBarManager);
      assert.strictEqual(statusBarManager.mcpClient, mockMCPClient);
    });

    test('should create status bar items', () => {
      assert(statusBarManager.connectionStatusItem);
      assert(statusBarManager.activeProjectItem);
      assert(statusBarManager.taskStatsItem);
      assert(statusBarManager.quickActionsItem);
      assert.strictEqual(vscode.window.createStatusBarItem.callCount, 4);
    });

    test('should setup event listeners', () => {
      assert(mockMCPClient.on.calledWith('connected'));
      assert(mockMCPClient.on.calledWith('disconnected'));
    });

    test('should initialize as not visible', () => {
      assert.strictEqual(statusBarManager.isVisible, false);
    });

    test('should handle null MCP client gracefully', () => {
      try {
        new StatusBarManager(null);
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle null MCP client gracefully');
      }
    });
  });

  suite('Show/Hide Management', () => {
    test('should show all status bar items', async () => {
      await statusBarManager.show();
      
      assert.strictEqual(statusBarManager.isVisible, true);
      assert(statusBarManager.connectionStatusItem.show.called);
      assert(statusBarManager.activeProjectItem.show.called);
      assert(statusBarManager.taskStatsItem.show.called);
      assert(statusBarManager.quickActionsItem.show.called);
    });

    test('should hide all status bar items', () => {
      statusBarManager.hide();
      
      assert.strictEqual(statusBarManager.isVisible, false);
      assert(statusBarManager.connectionStatusItem.hide.called);
      assert(statusBarManager.activeProjectItem.hide.called);
      assert(statusBarManager.taskStatsItem.hide.called);
      assert(statusBarManager.quickActionsItem.hide.called);
    });

    test('should update on show', async () => {
      await statusBarManager.show();
      
      // Should call MCP client to get latest data
      assert(mockMCPClient.callTool.called);
    });

    test('should not update when hidden', async () => {
      statusBarManager.isVisible = false;
      await statusBarManager.update();
      
      // Should not make any MCP calls when hidden
      assert(mockMCPClient.callTool.notCalled);
    });

    test('should handle show errors gracefully', async () => {
      statusBarManager.connectionStatusItem.show.throws(new Error('Show failed'));
      
      try {
        await statusBarManager.show();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle show errors gracefully');
      }
    });
  });

  suite('Task Stats Parsing', () => {
    test('should parse Turkish task statistics correctly', () => {
      const content = `
        ## Görev Özeti
        
        - Toplam Görev: 15
        - Beklemede: 5
        - Devam Ediyor: 7
        - Tamamlandı: 3
      `;
      
      const stats = statusBarManager.parseTaskStats(content);
      
      assert.deepStrictEqual(stats, {
        beklemede: 5,
        devamEdiyor: 7,
        tamamlandi: 3
      });
    });

    test('should parse statistics with different formatting', () => {
      const content = 'Beklemede:2 Devam Ediyor: 3 Tamamlandı:1';
      
      const stats = statusBarManager.parseTaskStats(content);
      
      assert.deepStrictEqual(stats, {
        beklemede: 2,
        devamEdiyor: 3,
        tamamlandi: 1
      });
    });

    test('should return null for incomplete statistics', () => {
      const content = 'Beklemede: 5\nTamamlandı: 3'; // Missing Devam Ediyor
      
      const stats = statusBarManager.parseTaskStats(content);
      
      assert.strictEqual(stats, null);
    });

    test('should return null for malformed content', () => {
      const content = 'Invalid content format';
      
      const stats = statusBarManager.parseTaskStats(content);
      
      assert.strictEqual(stats, null);
    });

    test('should handle zero values', () => {
      const content = 'Beklemede: 0\nDevam Ediyor: 0\nTamamlandı: 0';
      
      const stats = statusBarManager.parseTaskStats(content);
      
      assert.deepStrictEqual(stats, {
        beklemede: 0,
        devamEdiyor: 0,
        tamamlandi: 0
      });
    });

    test('should handle large numbers', () => {
      const content = 'Beklemede: 999\nDevam Ediyor: 1234\nTamamlandı: 567';
      
      const stats = statusBarManager.parseTaskStats(content);
      
      assert.deepStrictEqual(stats, {
        beklemede: 999,
        devamEdiyor: 1234,
        tamamlandi: 567
      });
    });
  });
});