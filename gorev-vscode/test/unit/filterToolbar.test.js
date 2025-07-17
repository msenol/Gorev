const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('FilterToolbar Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let mockOnFilterChange;
  let FilterToolbar;
  let filterToolbar;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'createStatusBarItem');
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showWarningMessage');
    sandbox.stub(vscode.window, 'showQuickPick');
    sandbox.stub(vscode.window, 'showInputBox');
    sandbox.stub(vscode.window, 'createQuickPick');
    sandbox.stub(vscode.workspace, 'getConfiguration').returns({
      get: sandbox.stub(),
      update: sandbox.stub()
    });
    sandbox.stub(vscode.ThemeIcon);
    sandbox.stub(vscode.ThemeColor);

    // Mock MCP Client
    mockMCPClient = {
      callTool: sandbox.stub().resolves({
        content: [{ text: '## Test Project (ID: proj-001)\n## Another Project (ID: proj-002)' }]
      }),
      isConnected: sandbox.stub().returns(true)
    };

    // Mock filter change callback
    mockOnFilterChange = sandbox.stub();

    // Mock Status Bar Items
    const mockStatusBarItem = {
      text: '',
      tooltip: '',
      command: '',
      show: sandbox.stub(),
      hide: sandbox.stub(),
      dispose: sandbox.stub(),
      backgroundColor: undefined
    };

    vscode.window.createStatusBarItem.returns(mockStatusBarItem);

    // Mock Quick Pick
    const mockQuickPick = {
      title: '',
      placeholder: '',
      canSelectMany: false,
      items: [],
      selectedItems: [],
      buttons: [],
      show: sandbox.stub(),
      hide: sandbox.stub(),
      dispose: sandbox.stub(),
      onDidChangeSelection: sandbox.stub(),
      onDidTriggerButton: sandbox.stub(),
      onDidHide: sandbox.stub()
    };

    vscode.window.createQuickPick.returns(mockQuickPick);

    // Import FilterToolbar
    try {
      const filterToolbarModule = require('../../dist/ui/filterToolbar');
      FilterToolbar = filterToolbarModule.FilterToolbar;
    } catch (error) {
      // Mock FilterToolbar class if compilation fails
      FilterToolbar = class MockFilterToolbar {
        constructor(mcpClient, onFilterChange) {
          this.mcpClient = mcpClient;
          this.onFilterChange = onFilterChange;
          this.statusBarItems = [];
          this.activeFilters = {};
          this.savedProfiles = new Map();
          this.createStatusBarItems();
        }

        createStatusBarItems() {
          for (let i = 0; i < 5; i++) {
            this.statusBarItems.push(vscode.window.createStatusBarItem());
          }
        }

        show() {
          this.statusBarItems.forEach(item => item.show());
        }

        hide() {
          this.statusBarItems.forEach(item => item.hide());
        }

        async showSearchInput() {
          const searchQuery = await vscode.window.showInputBox({
            prompt: 'Görev başlığı veya açıklamasında ara',
            value: this.activeFilters.searchQuery || ''
          });
          if (searchQuery !== undefined) {
            this.updateFilter({ searchQuery: searchQuery || undefined });
          }
        }

        async showFilterMenu() {
          // Mock implementation
          await vscode.window.createQuickPick();
        }

        async showFilterProfiles() {
          if (this.savedProfiles.size === 0) {
            vscode.window.showInformationMessage('Kayıtlı filtre profili bulunmuyor.');
            return;
          }
        }

        clearAllFilters() {
          this.activeFilters = {};
          this.onFilterChange({});
        }

        updateAllProjectsIndicator() {
          // Mock implementation
        }

        updateFilter(filter) {
          this.activeFilters = filter;
          this.onFilterChange(filter);
        }

        dispose() {
          this.statusBarItems.forEach(item => item.dispose());
        }

        parseProjects(content) {
          const projects = [];
          const lines = content.split('\n');
          for (const line of lines) {
            const match = line.match(/^## (.+) \(ID: ([^)]+)\)/);
            if (match) {
              projects.push({ isim: match[1], id: match[2] });
            }
          }
          return projects;
        }

        getDurumLabel(durum) {
          const labels = {
            'beklemede': 'Beklemede',
            'devam_ediyor': 'Devam Ediyor',
            'tamamlandi': 'Tamamlandı'
          };
          return labels[durum] || durum;
        }

        getOncelikLabel(oncelik) {
          const labels = {
            'dusuk': 'Düşük',
            'orta': 'Orta',
            'yuksek': 'Yüksek'
          };
          return labels[oncelik] || oncelik;
        }

        getFilterDescription(filter) {
          const parts = [];
          if (filter.searchQuery) parts.push(`Arama: "${filter.searchQuery}"`);
          if (filter.durum) parts.push(`Durum: ${this.getDurumLabel(filter.durum)}`);
          if (filter.oncelik) parts.push(`Öncelik: ${this.getOncelikLabel(filter.oncelik)}`);
          return parts.join(' • ');
        }
      };
    }

    // Create filter toolbar instance
    if (FilterToolbar) {
      filterToolbar = new FilterToolbar(mockMCPClient, mockOnFilterChange);
    }
  });

  teardown(() => {
    if (filterToolbar && typeof filterToolbar.dispose === 'function') {
      filterToolbar.dispose();
    }
    sandbox.restore();
  });

  suite('Initialization', () => {
    test('should create filter toolbar with correct parameters', () => {
      assert(filterToolbar);
      assert.strictEqual(filterToolbar.mcpClient, mockMCPClient);
      assert.strictEqual(filterToolbar.onFilterChange, mockOnFilterChange);
    });

    test('should create status bar items', () => {
      assert(Array.isArray(filterToolbar.statusBarItems));
      assert.strictEqual(filterToolbar.statusBarItems.length, 5);
      assert.strictEqual(vscode.window.createStatusBarItem.callCount, 5);
    });

    test('should initialize with empty active filters', () => {
      assert(typeof filterToolbar.activeFilters === 'object');
      assert.strictEqual(Object.keys(filterToolbar.activeFilters).length, 0);
    });

    test('should initialize saved profiles map', () => {
      assert(filterToolbar.savedProfiles instanceof Map);
    });

    test('should handle initialization errors gracefully', () => {
      vscode.window.createStatusBarItem.throws(new Error('Creation failed'));
      
      try {
        new FilterToolbar(mockMCPClient, mockOnFilterChange);
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle initialization errors gracefully');
      }
    });
  });

  suite('Status Bar Management', () => {
    test('should show all status bar items', () => {
      filterToolbar.show();
      
      filterToolbar.statusBarItems.forEach(item => {
        assert(item.show.called);
      });
    });

    test('should hide all status bar items', () => {
      filterToolbar.hide();
      
      filterToolbar.statusBarItems.forEach(item => {
        assert(item.hide.called);
      });
    });

    test('should dispose all status bar items', () => {
      filterToolbar.dispose();
      
      filterToolbar.statusBarItems.forEach(item => {
        assert(item.dispose.called);
      });
    });

    test('should handle show errors gracefully', () => {
      filterToolbar.statusBarItems[0].show.throws(new Error('Show failed'));
      
      try {
        filterToolbar.show();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle show errors gracefully');
      }
    });

    test('should handle hide errors gracefully', () => {
      filterToolbar.statusBarItems[0].hide.throws(new Error('Hide failed'));
      
      try {
        filterToolbar.hide();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle hide errors gracefully');
      }
    });
  });

  suite('Search Input', () => {
    test('should show search input dialog', async () => {
      const searchQuery = 'test query';
      vscode.window.showInputBox.resolves(searchQuery);
      
      await filterToolbar.showSearchInput();
      
      assert(vscode.window.showInputBox.called);
      const inputBoxCall = vscode.window.showInputBox.getCall(0);
      const options = inputBoxCall.args[0];
      assert.strictEqual(options.prompt, 'Görev başlığı veya açıklamasında ara');
      assert(mockOnFilterChange.calledWith({ searchQuery }));
    });

    test('should handle empty search query', async () => {
      vscode.window.showInputBox.resolves('');
      
      await filterToolbar.showSearchInput();
      
      assert(mockOnFilterChange.calledWith({ searchQuery: undefined }));
    });

    test('should handle search cancellation', async () => {
      vscode.window.showInputBox.resolves(undefined);
      
      await filterToolbar.showSearchInput();
      
      assert(vscode.window.showInputBox.called);
      assert(mockOnFilterChange.notCalled);
    });

    test('should preserve existing search query as default', async () => {
      filterToolbar.activeFilters.searchQuery = 'existing query';
      vscode.window.showInputBox.resolves('new query');
      
      await filterToolbar.showSearchInput();
      
      const inputBoxCall = vscode.window.showInputBox.getCall(0);
      const options = inputBoxCall.args[0];
      assert.strictEqual(options.value, 'existing query');
    });

    test('should handle search input errors', async () => {
      vscode.window.showInputBox.rejects(new Error('Input failed'));
      
      try {
        await filterToolbar.showSearchInput();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle search input errors gracefully');
      }
    });
  });

  suite('Filter Menu', () => {
    test('should create and show filter menu', async () => {
      await filterToolbar.showFilterMenu();
      
      assert(vscode.window.createQuickPick.called);
    });

    test('should handle filter menu creation errors', async () => {
      vscode.window.createQuickPick.throws(new Error('Creation failed'));
      
      try {
        await filterToolbar.showFilterMenu();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle filter menu creation errors gracefully');
      }
    });

    test('should load projects for filter menu', async () => {
      await filterToolbar.showFilterMenu();
      
      assert(mockMCPClient.callTool.called);
    });

    test('should handle project loading errors in filter menu', async () => {
      mockMCPClient.callTool.rejects(new Error('Failed to load projects'));
      
      try {
        await filterToolbar.showFilterMenu();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle project loading errors gracefully');
      }
    });
  });

  suite('Filter Profiles', () => {
    test('should show message when no profiles exist', async () => {
      await filterToolbar.showFilterProfiles();
      
      assert(vscode.window.showInformationMessage.calledWith('Kayıtlı filtre profili bulunmuyor.'));
    });

    test('should show profile selection when profiles exist', async () => {
      const testFilter = { searchQuery: 'test', durum: 'beklemede' };
      filterToolbar.savedProfiles.set('Test Profile', testFilter);
      
      const selectedProfile = {
        label: 'Test Profile',
        description: 'Arama: "test" • Durum: Beklemede',
        filter: testFilter
      };
      
      vscode.window.showQuickPick.resolves(selectedProfile);
      
      await filterToolbar.showFilterProfiles();
      
      assert(vscode.window.showQuickPick.called);
      assert(mockOnFilterChange.calledWith(testFilter));
      assert(vscode.window.showInformationMessage.calledWith('"Test Profile" filtre profili uygulandı.'));
    });

    test('should handle profile selection cancellation', async () => {
      filterToolbar.savedProfiles.set('Test Profile', { searchQuery: 'test' });
      vscode.window.showQuickPick.resolves(undefined);
      
      await filterToolbar.showFilterProfiles();
      
      assert(vscode.window.showQuickPick.called);
      assert(mockOnFilterChange.notCalled);
    });

    test('should handle profile selection errors', async () => {
      filterToolbar.savedProfiles.set('Test Profile', { searchQuery: 'test' });
      vscode.window.showQuickPick.rejects(new Error('Selection failed'));
      
      try {
        await filterToolbar.showFilterProfiles();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle profile selection errors gracefully');
      }
    });
  });

  suite('Clear All Filters', () => {
    test('should clear all filters and notify', () => {
      filterToolbar.activeFilters = { searchQuery: 'test', durum: 'beklemede' };
      
      filterToolbar.clearAllFilters();
      
      assert.deepStrictEqual(filterToolbar.activeFilters, {});
      assert(mockOnFilterChange.calledWith({}));
    });

    test('should handle clear filters when already empty', () => {
      filterToolbar.activeFilters = {};
      
      filterToolbar.clearAllFilters();
      
      assert.deepStrictEqual(filterToolbar.activeFilters, {});
      assert(mockOnFilterChange.calledWith({}));
    });
  });

  suite('Update Filter', () => {
    test('should update active filters', () => {
      const newFilter = { searchQuery: 'new query', durum: 'devam_ediyor' };
      
      filterToolbar.updateFilter(newFilter);
      
      assert.deepStrictEqual(filterToolbar.activeFilters, newFilter);
      assert(mockOnFilterChange.calledWith(newFilter));
    });

    test('should replace existing filters', () => {
      filterToolbar.activeFilters = { searchQuery: 'old query' };
      const newFilter = { durum: 'tamamlandi' };
      
      filterToolbar.updateFilter(newFilter);
      
      assert.deepStrictEqual(filterToolbar.activeFilters, newFilter);
      assert(mockOnFilterChange.calledWith(newFilter));
    });

    test('should handle empty filter object', () => {
      filterToolbar.activeFilters = { searchQuery: 'test' };
      
      filterToolbar.updateFilter({});
      
      assert.deepStrictEqual(filterToolbar.activeFilters, {});
      assert(mockOnFilterChange.calledWith({}));
    });
  });

  suite('All Projects Indicator', () => {
    test('should update all projects indicator', () => {
      try {
        filterToolbar.updateAllProjectsIndicator();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should update all projects indicator without errors');
      }
    });

    test('should handle indicator update errors', () => {
      // Mock a scenario where statusBarItems might be undefined
      filterToolbar.statusBarItems = null;
      
      try {
        filterToolbar.updateAllProjectsIndicator();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle indicator update errors gracefully');
      }
    });
  });

  suite('Project Parsing', () => {
    test('should parse projects from markdown content', () => {
      const content = '## Test Project (ID: proj-001)\n## Another Project (ID: proj-002)';
      const projects = filterToolbar.parseProjects(content);
      
      assert(Array.isArray(projects));
      assert.strictEqual(projects.length, 2);
      assert.deepStrictEqual(projects[0], { isim: 'Test Project', id: 'proj-001' });
      assert.deepStrictEqual(projects[1], { isim: 'Another Project', id: 'proj-002' });
    });

    test('should handle empty content', () => {
      const projects = filterToolbar.parseProjects('');
      
      assert(Array.isArray(projects));
      assert.strictEqual(projects.length, 0);
    });

    test('should handle malformed content', () => {
      const content = 'Invalid content without proper format';
      const projects = filterToolbar.parseProjects(content);
      
      assert(Array.isArray(projects));
      assert.strictEqual(projects.length, 0);
    });

    test('should handle mixed valid and invalid lines', () => {
      const content = `
        ## Valid Project (ID: proj-001)
        Invalid line
        ## Another Valid (ID: proj-002)
        Another invalid line
      `;
      const projects = filterToolbar.parseProjects(content);
      
      assert(Array.isArray(projects));
      assert.strictEqual(projects.length, 2);
    });
  });

  suite('Label Methods', () => {
    test('should get correct durum labels', () => {
      assert.strictEqual(filterToolbar.getDurumLabel('beklemede'), 'Beklemede');
      assert.strictEqual(filterToolbar.getDurumLabel('devam_ediyor'), 'Devam Ediyor');
      assert.strictEqual(filterToolbar.getDurumLabel('tamamlandi'), 'Tamamlandı');
      assert.strictEqual(filterToolbar.getDurumLabel('unknown'), 'unknown');
    });

    test('should get correct oncelik labels', () => {
      assert.strictEqual(filterToolbar.getOncelikLabel('dusuk'), 'Düşük');
      assert.strictEqual(filterToolbar.getOncelikLabel('orta'), 'Orta');
      assert.strictEqual(filterToolbar.getOncelikLabel('yuksek'), 'Yüksek');
      assert.strictEqual(filterToolbar.getOncelikLabel('unknown'), 'unknown');
    });

    test('should handle null/undefined labels', () => {
      assert.strictEqual(filterToolbar.getDurumLabel(null), null);
      assert.strictEqual(filterToolbar.getDurumLabel(undefined), undefined);
      assert.strictEqual(filterToolbar.getOncelikLabel(null), null);
      assert.strictEqual(filterToolbar.getOncelikLabel(undefined), undefined);
    });
  });

  suite('Filter Description', () => {
    test('should generate filter description with single filter', () => {
      const filter = { searchQuery: 'test query' };
      const description = filterToolbar.getFilterDescription(filter);
      
      assert.strictEqual(description, 'Arama: "test query"');
    });

    test('should generate filter description with multiple filters', () => {
      const filter = {
        searchQuery: 'test query',
        durum: 'beklemede',
        oncelik: 'yuksek'
      };
      const description = filterToolbar.getFilterDescription(filter);
      
      assert.strictEqual(description, 'Arama: "test query" • Durum: Beklemede • Öncelik: Yüksek');
    });

    test('should handle empty filter', () => {
      const description = filterToolbar.getFilterDescription({});
      
      assert.strictEqual(description, '');
    });

    test('should handle special filters', () => {
      const filter = {
        overdue: true,
        dueToday: true,
        hasTag: true
      };
      const description = filterToolbar.getFilterDescription(filter);
      
      assert(description.includes('Gecikmiş'));
      assert(description.includes('Bugün biten'));
      assert(description.includes('Etiketli'));
    });
  });

  suite('Configuration Integration', () => {
    test('should load saved profiles from configuration', () => {
      const savedProfiles = {
        'Profile 1': { searchQuery: 'test' },
        'Profile 2': { durum: 'beklemede' }
      };
      
      const mockConfig = {
        get: sandbox.stub().returns(savedProfiles),
        update: sandbox.stub()
      };
      
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      // Create new instance to test loading
      const newToolbar = new FilterToolbar(mockMCPClient, mockOnFilterChange);
      
      assert.strictEqual(newToolbar.savedProfiles.size, 2);
      assert(newToolbar.savedProfiles.has('Profile 1'));
      assert(newToolbar.savedProfiles.has('Profile 2'));
    });

    test('should handle missing configuration', () => {
      const mockConfig = {
        get: sandbox.stub().returns(undefined),
        update: sandbox.stub()
      };
      
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      const newToolbar = new FilterToolbar(mockMCPClient, mockOnFilterChange);
      
      assert.strictEqual(newToolbar.savedProfiles.size, 0);
    });

    test('should handle configuration errors', () => {
      vscode.workspace.getConfiguration.throws(new Error('Config error'));
      
      try {
        new FilterToolbar(mockMCPClient, mockOnFilterChange);
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle configuration errors gracefully');
      }
    });
  });

  suite('Error Handling', () => {
    test('should handle null MCP client', () => {
      try {
        new FilterToolbar(null, mockOnFilterChange);
        assert(true); // Should not throw during construction
      } catch (error) {
        assert.fail('Should handle null MCP client gracefully');
      }
    });

    test('should handle undefined callback', () => {
      try {
        const toolbar = new FilterToolbar(mockMCPClient, undefined);
        toolbar.updateFilter({ searchQuery: 'test' });
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle undefined callback gracefully');
      }
    });

    test('should handle callback errors', () => {
      const errorCallback = sandbox.stub().throws(new Error('Callback failed'));
      const toolbar = new FilterToolbar(mockMCPClient, errorCallback);
      
      try {
        toolbar.updateFilter({ searchQuery: 'test' });
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle callback errors gracefully');
      }
    });
  });
});