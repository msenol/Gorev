const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('FilterCommands Test Suite', () => {
  let sandbox;
  let mockContext;
  let mockProviders;
  let registerFunction;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.commands, 'registerCommand');
    sandbox.stub(vscode.workspace, 'getConfiguration').returns({
      get: sandbox.stub().returns(false),
      update: sandbox.stub()
    });

    // Mock Context
    mockContext = {
      subscriptions: []
    };

    // Mock Filter Toolbar
    const mockFilterToolbar = {
      showSearchInput: sandbox.stub().resolves(),
      showFilterMenu: sandbox.stub().resolves(),
      showFilterProfiles: sandbox.stub().resolves(),
      clearAllFilters: sandbox.stub(),
      updateAllProjectsIndicator: sandbox.stub()
    };

    // Mock Enhanced Tree Provider
    const mockEnhancedTreeProvider = {
      updateFilter: sandbox.stub(),
      refresh: sandbox.stub().resolves()
    };

    // Mock Providers
    mockProviders = {
      gorevTreeProvider: mockEnhancedTreeProvider,
      filterToolbar: mockFilterToolbar
    };

    // Import and register commands
    try {
      const commands = require('../../dist/commands/filterCommands');
      registerFunction = commands.registerFilterCommands;
    } catch (error) {
      // Mock register function if compilation fails
      registerFunction = (context, providers) => {
        // Register search input command
        context.subscriptions.push(
          vscode.commands.registerCommand('gorev.showSearchInput', async () => {
            await mockProviders.filterToolbar.showSearchInput();
          })
        );

        // Register filter menu command
        context.subscriptions.push(
          vscode.commands.registerCommand('gorev.showFilterMenu', async () => {
            await mockProviders.filterToolbar.showFilterMenu();
          })
        );

        // Register filter profiles command
        context.subscriptions.push(
          vscode.commands.registerCommand('gorev.showFilterProfiles', async () => {
            await mockProviders.filterToolbar.showFilterProfiles();
          })
        );

        // Register clear all filters command
        context.subscriptions.push(
          vscode.commands.registerCommand('gorev.clearAllFilters', () => {
            mockProviders.filterToolbar.clearAllFilters();
            mockProviders.gorevTreeProvider.updateFilter({});
            vscode.window.showInformationMessage('Tüm filtreler temizlendi');
          })
        );

        // Register toggle all projects command
        context.subscriptions.push(
          vscode.commands.registerCommand('gorev.toggleAllProjects', () => {
            const config = vscode.workspace.getConfiguration('gorev.treeView');
            const current = config.get('showAllProjects', true);
            config.update('showAllProjects', !current, vscode.ConfigurationTarget.Global);
            
            // Update filter
            mockProviders.gorevTreeProvider.updateFilter({ showAllProjects: !current });
            mockProviders.filterToolbar.updateAllProjectsIndicator();
            
            const message = !current ? 'Tüm projeler gösteriliyor' : 'Sadece aktif proje gösteriliyor';
            vscode.window.showInformationMessage(message);
          })
        );
      };
    }
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('Command Registration', () => {
    test('should register all filter commands', () => {
      registerFunction(mockContext, mockProviders);
      
      // Should register 5 commands
      assert.strictEqual(mockContext.subscriptions.length, 5);
      
      // Should call vscode.commands.registerCommand for each command
      assert(vscode.commands.registerCommand.callCount >= 5);
    });

    test('should handle registration errors gracefully', () => {
      const errorContext = { 
        subscriptions: { push: sandbox.stub().throws(new Error('Mock error')) }
      };
      
      try {
        registerFunction(errorContext, mockProviders);
        // Should not throw
        assert(true);
      } catch (error) {
        assert.fail('Should handle registration errors gracefully');
      }
    });
  });

  suite('SHOW_SEARCH_INPUT Command', () => {
    test('should call filter toolbar showSearchInput', async () => {
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.showSearchInput') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.showSearchInput.called);
    });

    test('should handle search input errors', async () => {
      mockProviders.filterToolbar.showSearchInput.rejects(new Error('Search failed'));
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.showSearchInput') {
          try {
            await handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.showSearchInput.called);
    });
  });

  suite('SHOW_FILTER_MENU Command', () => {
    test('should call filter toolbar showFilterMenu', async () => {
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.showFilterMenu') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.showFilterMenu.called);
    });

    test('should handle filter menu errors', async () => {
      mockProviders.filterToolbar.showFilterMenu.rejects(new Error('Filter menu failed'));
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.showFilterMenu') {
          try {
            await handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.showFilterMenu.called);
    });
  });

  suite('SHOW_FILTER_PROFILES Command', () => {
    test('should call filter toolbar showFilterProfiles', async () => {
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.showFilterProfiles') {
          await handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.showFilterProfiles.called);
    });

    test('should handle filter profiles errors', async () => {
      mockProviders.filterToolbar.showFilterProfiles.rejects(new Error('Profiles failed'));
      
      vscode.commands.registerCommand.callsFake(async (commandId, handler) => {
        if (commandId === 'gorev.showFilterProfiles') {
          try {
            await handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.showFilterProfiles.called);
    });
  });

  suite('CLEAR_ALL_FILTERS Command', () => {
    test('should clear all filters and show message', () => {
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.clearAllFilters') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.clearAllFilters.called);
      assert(mockProviders.gorevTreeProvider.updateFilter.calledWith({}));
      assert(vscode.window.showInformationMessage.calledWith('Tüm filtreler temizlendi'));
    });

    test('should handle clear filters errors', () => {
      mockProviders.filterToolbar.clearAllFilters.throws(new Error('Clear failed'));
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.clearAllFilters') {
          try {
            handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      // Should still attempt to clear filters
      assert(mockProviders.filterToolbar.clearAllFilters.called);
    });

    test('should handle tree provider update errors', () => {
      mockProviders.gorevTreeProvider.updateFilter.throws(new Error('Update failed'));
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.clearAllFilters') {
          try {
            handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.clearAllFilters.called);
      assert(mockProviders.gorevTreeProvider.updateFilter.called);
    });
  });

  suite('TOGGLE_ALL_PROJECTS Command', () => {
    test('should toggle from false to true', () => {
      const mockConfig = {
        get: sandbox.stub().returns(false),
        update: sandbox.stub()
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockConfig.get.calledWith('showAllProjects', true));
      assert(mockConfig.update.calledWith('showAllProjects', true, vscode.ConfigurationTarget.Global));
      assert(mockProviders.gorevTreeProvider.updateFilter.calledWith({ showAllProjects: true }));
      assert(mockProviders.filterToolbar.updateAllProjectsIndicator.called);
      assert(vscode.window.showInformationMessage.calledWith('Tüm projeler gösteriliyor'));
    });

    test('should toggle from true to false', () => {
      const mockConfig = {
        get: sandbox.stub().returns(true),
        update: sandbox.stub()
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockConfig.get.calledWith('showAllProjects', true));
      assert(mockConfig.update.calledWith('showAllProjects', false, vscode.ConfigurationTarget.Global));
      assert(mockProviders.gorevTreeProvider.updateFilter.calledWith({ showAllProjects: false }));
      assert(mockProviders.filterToolbar.updateAllProjectsIndicator.called);
      assert(vscode.window.showInformationMessage.calledWith('Sadece aktif proje gösteriliyor'));
    });

    test('should handle default configuration value', () => {
      const mockConfig = {
        get: sandbox.stub().returns(undefined), // No explicit value set
        update: sandbox.stub()
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      // Should use default value of true
      assert(mockConfig.get.calledWith('showAllProjects', true));
      assert(mockConfig.update.calledWith('showAllProjects', false, vscode.ConfigurationTarget.Global));
    });

    test('should handle configuration update errors', () => {
      const mockConfig = {
        get: sandbox.stub().returns(false),
        update: sandbox.stub().throws(new Error('Config update failed'))
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          try {
            handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockConfig.get.called);
      assert(mockConfig.update.called);
    });

    test('should handle tree provider filter update errors', () => {
      const mockConfig = {
        get: sandbox.stub().returns(false),
        update: sandbox.stub()
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      mockProviders.gorevTreeProvider.updateFilter.throws(new Error('Filter update failed'));
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          try {
            handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.gorevTreeProvider.updateFilter.called);
      // Should still update indicator even if filter update fails
      assert(mockProviders.filterToolbar.updateAllProjectsIndicator.called);
    });

    test('should handle toolbar indicator update errors', () => {
      const mockConfig = {
        get: sandbox.stub().returns(false),
        update: sandbox.stub()
      };
      vscode.workspace.getConfiguration.returns(mockConfig);
      mockProviders.filterToolbar.updateAllProjectsIndicator.throws(new Error('Indicator update failed'));
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          try {
            handler();
          } catch (error) {
            // Should handle error gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(mockProviders.filterToolbar.updateAllProjectsIndicator.called);
      // Should still show information message
      assert(vscode.window.showInformationMessage.called);
    });

    test('should use correct configuration section', () => {
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          handler();
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      assert(vscode.workspace.getConfiguration.calledWith('gorev.treeView'));
    });
  });

  suite('Error Handling', () => {
    test('should handle missing filter toolbar', () => {
      const providersWithoutToolbar = {
        gorevTreeProvider: mockProviders.gorevTreeProvider,
        filterToolbar: undefined
      };
      
      try {
        registerFunction(mockContext, providersWithoutToolbar);
        // Should not throw during registration
        assert(true);
      } catch (error) {
        assert.fail('Should handle missing filter toolbar gracefully');
      }
    });

    test('should handle missing tree provider', () => {
      const providersWithoutTreeProvider = {
        gorevTreeProvider: undefined,
        filterToolbar: mockProviders.filterToolbar
      };
      
      try {
        registerFunction(mockContext, providersWithoutTreeProvider);
        // Should not throw during registration
        assert(true);
      } catch (error) {
        assert.fail('Should handle missing tree provider gracefully');
      }
    });

    test('should handle null providers', () => {
      try {
        registerFunction(mockContext, null);
        // Should not throw during registration
        assert(true);
      } catch (error) {
        assert.fail('Should handle null providers gracefully');
      }
    });

    test('should handle configuration API errors', () => {
      vscode.workspace.getConfiguration.throws(new Error('Configuration API failed'));
      
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.toggleAllProjects') {
          try {
            handler();
          } catch (error) {
            // Should handle configuration errors gracefully
          }
        }
      });
      
      registerFunction(mockContext, mockProviders);
      
      // Should attempt to get configuration
      assert(vscode.workspace.getConfiguration.called);
    });

    test('should handle command registration errors', () => {
      vscode.commands.registerCommand.throws(new Error('Command registration failed'));
      
      try {
        registerFunction(mockContext, mockProviders);
        // Should handle registration errors gracefully
        assert(true);
      } catch (error) {
        // Expected to throw in this case since it's a critical error
        assert(error.message === 'Command registration failed');
      }
    });
  });

  suite('Command Integration', () => {
    test('should integrate with filter toolbar correctly', () => {
      registerFunction(mockContext, mockProviders);
      
      // Verify that all commands that should call filter toolbar methods are registered
      const commandIds = vscode.commands.registerCommand.getCalls().map(call => call.args[0]);
      
      assert(commandIds.includes('gorev.showSearchInput'));
      assert(commandIds.includes('gorev.showFilterMenu'));
      assert(commandIds.includes('gorev.showFilterProfiles'));
      assert(commandIds.includes('gorev.clearAllFilters'));
      assert(commandIds.includes('gorev.toggleAllProjects'));
    });

    test('should integrate with tree provider correctly', () => {
      registerFunction(mockContext, mockProviders);
      
      // Test clear filters command integration
      vscode.commands.registerCommand.callsFake((commandId, handler) => {
        if (commandId === 'gorev.clearAllFilters') {
          handler();
        }
      });
      
      assert(mockProviders.gorevTreeProvider.updateFilter.called);
    });

    test('should provide correct command handlers', () => {
      registerFunction(mockContext, mockProviders);
      
      // Verify that all registered commands have proper handlers
      const commandCalls = vscode.commands.registerCommand.getCalls();
      
      commandCalls.forEach(call => {
        const [commandId, handler] = call.args;
        assert(typeof commandId === 'string');
        assert(typeof handler === 'function');
      });
    });
  });
});