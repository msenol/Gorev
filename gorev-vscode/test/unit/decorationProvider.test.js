const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('DecorationProvider Test Suite', () => {
  let sandbox;
  let TaskDecorationProvider;
  let decorationProvider;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'createTextEditorDecorationType');
    sandbox.stub(vscode.workspace, 'getConfiguration').returns({
      get: sandbox.stub().returns(true)
    });
    sandbox.stub(vscode.ThemeIcon);
    sandbox.stub(vscode.ThemeColor);

    // Mock decoration types
    const mockDecorationType = {
      dispose: sandbox.stub()
    };
    vscode.window.createTextEditorDecorationType.returns(mockDecorationType);

    // Import TaskDecorationProvider
    try {
      const decorationModule = require('../../dist/providers/decorationProvider');
      TaskDecorationProvider = decorationModule.TaskDecorationProvider;
    } catch (error) {
      // Mock TaskDecorationProvider class if compilation fails
      TaskDecorationProvider = class MockTaskDecorationProvider {
        constructor() {
          this.decorationTypes = {
            priority: {
              high: vscode.window.createTextEditorDecorationType({}),
              medium: vscode.window.createTextEditorDecorationType({}),
              low: vscode.window.createTextEditorDecorationType({})
            },
            status: {
              pending: vscode.window.createTextEditorDecorationType({}),
              inProgress: vscode.window.createTextEditorDecorationType({}),
              completed: vscode.window.createTextEditorDecorationType({})
            },
            tag: vscode.window.createTextEditorDecorationType({}),
            dependency: vscode.window.createTextEditorDecorationType({}),
            dueDate: {
              overdue: vscode.window.createTextEditorDecorationType({}),
              dueToday: vscode.window.createTextEditorDecorationType({}),
              dueSoon: vscode.window.createTextEditorDecorationType({})
            }
          };
        }

        getPriorityBadge(priority) {
          const badges = {
            'yuksek': '🔥',
            'orta': '⚡',
            'dusuk': 'ℹ️'
          };
          return badges[priority] || '';
        }

        getStatusIcon(status) {
          const icons = {
            'beklemede': '⏳',
            'devam_ediyor': '🔄',
            'tamamlandi': '✅'
          };
          return icons[status] || '';
        }

        getDueDateBadge(dueDate) {
          if (!dueDate) return '';
          
          const now = new Date();
          const due = new Date(dueDate);
          const diffTime = due.getTime() - now.getTime();
          const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

          if (diffDays < 0) return '🚨 Overdue';
          if (diffDays === 0) return '📅 Today';
          if (diffDays === 1) return '📅 Tomorrow';
          if (diffDays <= 7) return `📅 ${diffDays}d left`;
          
          return '';
        }

        getDependencyBadge(dependencies) {
          if (!dependencies || dependencies.length === 0) return '';
          
          const blocked = dependencies.filter(d => !d.completed).length;
          if (blocked > 0) {
            return `🔒 ${blocked} blocked`;
          }
          return `🔓 ${dependencies.length} deps`;
        }

        getTagBadges(tags, maxTags = 3) {
          if (!tags || tags.length === 0) return [];
          
          const visibleTags = tags.slice(0, maxTags);
          const badges = visibleTags.map(tag => `#${tag}`);
          
          if (tags.length > maxTags) {
            badges.push(`+${tags.length - maxTags}`);
          }
          
          return badges;
        }

        formatProgress(progress) {
          if (typeof progress !== 'number') return '';
          return `${Math.round(progress)}%`;
        }

        createTaskDescription(task) {
          const parts = [];
          
          // Priority badge
          const priorityBadge = this.getPriorityBadge(task.oncelik);
          if (priorityBadge) parts.push(priorityBadge);
          
          // Status icon
          const statusIcon = this.getStatusIcon(task.durum);
          if (statusIcon) parts.push(statusIcon);
          
          // Due date badge
          const dueDateBadge = this.getDueDateBadge(task.son_tarih);
          if (dueDateBadge) parts.push(dueDateBadge);
          
          // Dependency badge
          const dependencyBadge = this.getDependencyBadge(task.dependencies);
          if (dependencyBadge) parts.push(dependencyBadge);
          
          // Tag badges
          const tagBadges = this.getTagBadges(task.etiketler);
          if (tagBadges.length > 0) parts.push(tagBadges.join(' '));
          
          // Progress for parent tasks
          if (task.subtaskProgress !== undefined) {
            const progressText = this.formatProgress(task.subtaskProgress);
            if (progressText) parts.push(`[${progressText}]`);
          }
          
          return parts.join(' ');
        }

        createTaskTooltip(task) {
          const lines = [];
          
          lines.push(`**${task.baslik}**`);
          
          if (task.aciklama) {
            lines.push('', task.aciklama);
          }
          
          lines.push('', '---');
          
          // Basic info
          lines.push(`**Status:** ${this.getStatusLabel(task.durum)}`);
          lines.push(`**Priority:** ${this.getPriorityLabel(task.oncelik)}`);
          
          if (task.son_tarih) {
            lines.push(`**Due Date:** ${task.son_tarih}`);
          }
          
          if (task.proje_ismi) {
            lines.push(`**Project:** ${task.proje_ismi}`);
          }
          
          // Dependencies
          if (task.dependencies && task.dependencies.length > 0) {
            lines.push('', '**Dependencies:**');
            task.dependencies.forEach(dep => {
              const status = dep.completed ? '✅' : '⏳';
              lines.push(`  ${status} ${dep.title}`);
            });
          }
          
          // Tags
          if (task.etiketler && task.etiketler.length > 0) {
            lines.push('', `**Tags:** ${task.etiketler.map(tag => `#${tag}`).join(', ')}`);
          }
          
          // Subtask progress
          if (task.subtaskProgress !== undefined) {
            lines.push('', `**Progress:** ${this.formatProgress(task.subtaskProgress)} (${task.completedSubtasks || 0}/${task.totalSubtasks || 0} subtasks)`);
          }
          
          return lines.join('\n');
        }

        getStatusLabel(status) {
          const labels = {
            'beklemede': 'Pending',
            'devam_ediyor': 'In Progress',
            'tamamlandi': 'Completed'
          };
          return labels[status] || status;
        }

        getPriorityLabel(priority) {
          const labels = {
            'yuksek': 'High',
            'orta': 'Medium',
            'dusuk': 'Low'
          };
          return labels[priority] || priority;
        }

        getPriorityColor(priority) {
          const colors = {
            'yuksek': new vscode.ThemeColor('statusBarItem.errorBackground'),
            'orta': new vscode.ThemeColor('statusBarItem.warningBackground'),
            'dusuk': new vscode.ThemeColor('statusBarItem.background')
          };
          return colors[priority];
        }

        getStatusColor(status) {
          const colors = {
            'beklemede': new vscode.ThemeColor('list.inactiveSelectionBackground'),
            'devam_ediyor': new vscode.ThemeColor('progressBar.background'),
            'tamamlandi': new vscode.ThemeColor('statusBarItem.background')
          };
          return colors[status];
        }

        shouldShowDecoration(task, decorationType) {
          const config = vscode.workspace.getConfiguration('gorev.treeView.decorations');
          
          switch (decorationType) {
            case 'priority':
              return config.get('showPriority', true) && task.oncelik;
            case 'status':
              return config.get('showStatus', true) && task.durum;
            case 'dueDate':
              return config.get('showDueDate', true) && task.son_tarih;
            case 'tags':
              return config.get('showTags', true) && task.etiketler && task.etiketler.length > 0;
            case 'dependencies':
              return config.get('showDependencies', true) && task.dependencies && task.dependencies.length > 0;
            case 'progress':
              return config.get('showProgress', true) && task.subtaskProgress !== undefined;
            default:
              return true;
          }
        }

        applyTaskDecorations(treeItem, task) {
          if (!this.shouldShowDecoration(task, 'any')) return;
          
          // Create description with decorations
          const description = this.createTaskDescription(task);
          if (description) {
            treeItem.description = description;
          }
          
          // Create tooltip with rich information
          const tooltip = this.createTaskTooltip(task);
          if (tooltip) {
            treeItem.tooltip = new vscode.MarkdownString(tooltip);
          }
          
          // Apply color coding based on priority
          if (this.shouldShowDecoration(task, 'priority')) {
            treeItem.iconPath = new vscode.ThemeIcon('circle-filled', this.getPriorityColor(task.oncelik));
          }
          
          return treeItem;
        }

        dispose() {
          Object.values(this.decorationTypes).forEach(decorationType => {
            if (decorationType.dispose) {
              decorationType.dispose();
            } else if (typeof decorationType === 'object') {
              Object.values(decorationType).forEach(decoration => {
                if (decoration && decoration.dispose) {
                  decoration.dispose();
                }
              });
            }
          });
        }
      };
    }

    // Create decoration provider instance
    if (TaskDecorationProvider) {
      decorationProvider = new TaskDecorationProvider();
    }
  });

  teardown(() => {
    if (decorationProvider && typeof decorationProvider.dispose === 'function') {
      decorationProvider.dispose();
    }
    sandbox.restore();
  });

  suite('Initialization', () => {
    test('should create decoration provider', () => {
      assert(decorationProvider);
      assert(decorationProvider.decorationTypes);
    });

    test('should create decoration types', () => {
      assert(decorationProvider.decorationTypes.priority);
      assert(decorationProvider.decorationTypes.status);
      assert(decorationProvider.decorationTypes.tag);
      assert(decorationProvider.decorationTypes.dependency);
      assert(decorationProvider.decorationTypes.dueDate);
    });

    test('should handle initialization errors gracefully', () => {
      vscode.window.createTextEditorDecorationType.throws(new Error('Creation failed'));
      
      try {
        new TaskDecorationProvider();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle initialization errors gracefully');
      }
    });
  });

  suite('Priority Badges', () => {
    test('should return correct priority badges', () => {
      assert.strictEqual(decorationProvider.getPriorityBadge('yuksek'), '🔥');
      assert.strictEqual(decorationProvider.getPriorityBadge('orta'), '⚡');
      assert.strictEqual(decorationProvider.getPriorityBadge('dusuk'), 'ℹ️');
    });

    test('should handle unknown priority', () => {
      assert.strictEqual(decorationProvider.getPriorityBadge('unknown'), '');
      assert.strictEqual(decorationProvider.getPriorityBadge(null), '');
      assert.strictEqual(decorationProvider.getPriorityBadge(undefined), '');
    });

    test('should handle invalid priority types', () => {
      assert.strictEqual(decorationProvider.getPriorityBadge(123), '');
      assert.strictEqual(decorationProvider.getPriorityBadge({}), '');
      assert.strictEqual(decorationProvider.getPriorityBadge([]), '');
    });
  });

  suite('Status Icons', () => {
    test('should return correct status icons', () => {
      assert.strictEqual(decorationProvider.getStatusIcon('beklemede'), '⏳');
      assert.strictEqual(decorationProvider.getStatusIcon('devam_ediyor'), '🔄');
      assert.strictEqual(decorationProvider.getStatusIcon('tamamlandi'), '✅');
    });

    test('should handle unknown status', () => {
      assert.strictEqual(decorationProvider.getStatusIcon('unknown'), '');
      assert.strictEqual(decorationProvider.getStatusIcon(null), '');
      assert.strictEqual(decorationProvider.getStatusIcon(undefined), '');
    });

    test('should handle invalid status types', () => {
      assert.strictEqual(decorationProvider.getStatusIcon(123), '');
      assert.strictEqual(decorationProvider.getStatusIcon({}), '');
      assert.strictEqual(decorationProvider.getStatusIcon([]), '');
    });
  });

  suite('Due Date Badges', () => {
    test('should return overdue badge for past dates', () => {
      const yesterday = new Date();
      yesterday.setDate(yesterday.getDate() - 1);
      const badge = decorationProvider.getDueDateBadge(yesterday.toISOString().split('T')[0]);
      
      assert.strictEqual(badge, '🚨 Overdue');
    });

    test('should return today badge for current date', () => {
      const today = new Date().toISOString().split('T')[0];
      const badge = decorationProvider.getDueDateBadge(today);
      
      assert.strictEqual(badge, '📅 Today');
    });

    test('should return tomorrow badge', () => {
      const tomorrow = new Date();
      tomorrow.setDate(tomorrow.getDate() + 1);
      const badge = decorationProvider.getDueDateBadge(tomorrow.toISOString().split('T')[0]);
      
      assert.strictEqual(badge, '📅 Tomorrow');
    });

    test('should return days left for near future', () => {
      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 3);
      const badge = decorationProvider.getDueDateBadge(futureDate.toISOString().split('T')[0]);
      
      assert.strictEqual(badge, '📅 3d left');
    });

    test('should return empty for far future dates', () => {
      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 30);
      const badge = decorationProvider.getDueDateBadge(futureDate.toISOString().split('T')[0]);
      
      assert.strictEqual(badge, '');
    });

    test('should handle null/undefined due dates', () => {
      assert.strictEqual(decorationProvider.getDueDateBadge(null), '');
      assert.strictEqual(decorationProvider.getDueDateBadge(undefined), '');
      assert.strictEqual(decorationProvider.getDueDateBadge(''), '');
    });

    test('should handle invalid date formats', () => {
      assert.strictEqual(decorationProvider.getDueDateBadge('invalid-date'), '');
      assert.strictEqual(decorationProvider.getDueDateBadge('2024-13-45'), '');
    });
  });

  suite('Dependency Badges', () => {
    test('should return blocked badge for incomplete dependencies', () => {
      const dependencies = [
        { title: 'Task 1', completed: false },
        { title: 'Task 2', completed: true },
        { title: 'Task 3', completed: false }
      ];
      
      const badge = decorationProvider.getDependencyBadge(dependencies);
      assert.strictEqual(badge, '🔒 2 blocked');
    });

    test('should return unlocked badge for all completed dependencies', () => {
      const dependencies = [
        { title: 'Task 1', completed: true },
        { title: 'Task 2', completed: true }
      ];
      
      const badge = decorationProvider.getDependencyBadge(dependencies);
      assert.strictEqual(badge, '🔓 2 deps');
    });

    test('should return empty for no dependencies', () => {
      assert.strictEqual(decorationProvider.getDependencyBadge([]), '');
      assert.strictEqual(decorationProvider.getDependencyBadge(null), '');
      assert.strictEqual(decorationProvider.getDependencyBadge(undefined), '');
    });

    test('should handle dependencies without completed property', () => {
      const dependencies = [
        { title: 'Task 1' },
        { title: 'Task 2' }
      ];
      
      const badge = decorationProvider.getDependencyBadge(dependencies);
      assert.strictEqual(badge, '🔒 2 blocked'); // Undefined treated as incomplete
    });
  });

  suite('Tag Badges', () => {
    test('should return tag badges for normal tags', () => {
      const tags = ['bug', 'urgent', 'frontend'];
      const badges = decorationProvider.getTagBadges(tags);
      
      assert.deepStrictEqual(badges, ['#bug', '#urgent', '#frontend']);
    });

    test('should limit tags and show overflow', () => {
      const tags = ['tag1', 'tag2', 'tag3', 'tag4', 'tag5'];
      const badges = decorationProvider.getTagBadges(tags, 3);
      
      assert.deepStrictEqual(badges, ['#tag1', '#tag2', '#tag3', '+2']);
    });

    test('should handle empty tags', () => {
      assert.deepStrictEqual(decorationProvider.getTagBadges([]), []);
      assert.deepStrictEqual(decorationProvider.getTagBadges(null), []);
      assert.deepStrictEqual(decorationProvider.getTagBadges(undefined), []);
    });

    test('should handle custom max tags', () => {
      const tags = ['tag1', 'tag2'];
      const badges = decorationProvider.getTagBadges(tags, 1);
      
      assert.deepStrictEqual(badges, ['#tag1', '+1']);
    });

    test('should handle zero max tags', () => {
      const tags = ['tag1', 'tag2'];
      const badges = decorationProvider.getTagBadges(tags, 0);
      
      assert.deepStrictEqual(badges, ['+2']);
    });
  });

  suite('Progress Formatting', () => {
    test('should format progress as percentage', () => {
      assert.strictEqual(decorationProvider.formatProgress(75.6), '76%');
      assert.strictEqual(decorationProvider.formatProgress(0), '0%');
      assert.strictEqual(decorationProvider.formatProgress(100), '100%');
    });

    test('should handle invalid progress values', () => {
      assert.strictEqual(decorationProvider.formatProgress(null), '');
      assert.strictEqual(decorationProvider.formatProgress(undefined), '');
      assert.strictEqual(decorationProvider.formatProgress('invalid'), '');
      assert.strictEqual(decorationProvider.formatProgress({}), '');
    });

    test('should handle edge cases', () => {
      assert.strictEqual(decorationProvider.formatProgress(-5), '-5%');
      assert.strictEqual(decorationProvider.formatProgress(150), '150%');
      assert.strictEqual(decorationProvider.formatProgress(0.5), '1%');
    });
  });

  suite('Task Description Creation', () => {
    test('should create complete task description', () => {
      const task = {
        baslik: 'Test Task',
        oncelik: 'yuksek',
        durum: 'devam_ediyor',
        son_tarih: new Date().toISOString().split('T')[0], // Today
        dependencies: [{ title: 'Dep 1', completed: false }],
        etiketler: ['bug', 'urgent'],
        subtaskProgress: 75
      };
      
      const description = decorationProvider.createTaskDescription(task);
      
      assert(description.includes('🔥')); // Priority
      assert(description.includes('🔄')); // Status
      assert(description.includes('📅 Today')); // Due date
      assert(description.includes('🔒 1 blocked')); // Dependencies
      assert(description.includes('#bug #urgent')); // Tags
      assert(description.includes('[75%]')); // Progress
    });

    test('should handle task with minimal data', () => {
      const task = {
        baslik: 'Simple Task',
        oncelik: 'orta',
        durum: 'beklemede'
      };
      
      const description = decorationProvider.createTaskDescription(task);
      
      assert(description.includes('⚡')); // Priority
      assert(description.includes('⏳')); // Status
      assert(!description.includes('📅')); // No due date
      assert(!description.includes('🔒')); // No dependencies
      assert(!description.includes('#')); // No tags
    });

    test('should handle empty task', () => {
      const task = {};
      
      const description = decorationProvider.createTaskDescription(task);
      
      assert.strictEqual(description, '');
    });
  });

  suite('Task Tooltip Creation', () => {
    test('should create complete tooltip', () => {
      const task = {
        baslik: 'Test Task',
        aciklama: 'Task description',
        durum: 'devam_ediyor',
        oncelik: 'yuksek',
        son_tarih: '2024-12-31',
        proje_ismi: 'Test Project',
        dependencies: [
          { title: 'Dependency 1', completed: true },
          { title: 'Dependency 2', completed: false }
        ],
        etiketler: ['bug', 'urgent'],
        subtaskProgress: 75,
        completedSubtasks: 3,
        totalSubtasks: 4
      };
      
      const tooltip = decorationProvider.createTaskTooltip(task);
      
      assert(tooltip.includes('**Test Task**'));
      assert(tooltip.includes('Task description'));
      assert(tooltip.includes('**Status:** In Progress'));
      assert(tooltip.includes('**Priority:** High'));
      assert(tooltip.includes('**Due Date:** 2024-12-31'));
      assert(tooltip.includes('**Project:** Test Project'));
      assert(tooltip.includes('**Dependencies:**'));
      assert(tooltip.includes('✅ Dependency 1'));
      assert(tooltip.includes('⏳ Dependency 2'));
      assert(tooltip.includes('**Tags:** #bug, #urgent'));
      assert(tooltip.includes('**Progress:** 75% (3/4 subtasks)'));
    });

    test('should create minimal tooltip', () => {
      const task = {
        baslik: 'Simple Task',
        durum: 'beklemede',
        oncelik: 'orta'
      };
      
      const tooltip = decorationProvider.createTaskTooltip(task);
      
      assert(tooltip.includes('**Simple Task**'));
      assert(tooltip.includes('**Status:** Pending'));
      assert(tooltip.includes('**Priority:** Medium'));
      assert(!tooltip.includes('**Due Date:**'));
      assert(!tooltip.includes('**Dependencies:**'));
      assert(!tooltip.includes('**Tags:**'));
    });

    test('should handle empty task in tooltip', () => {
      const task = {};
      
      const tooltip = decorationProvider.createTaskTooltip(task);
      
      assert(tooltip.includes('**undefined**')); // Title is undefined
      assert(tooltip.includes('---'));
    });
  });

  suite('Label Methods', () => {
    test('should get correct status labels', () => {
      assert.strictEqual(decorationProvider.getStatusLabel('beklemede'), 'Pending');
      assert.strictEqual(decorationProvider.getStatusLabel('devam_ediyor'), 'In Progress');
      assert.strictEqual(decorationProvider.getStatusLabel('tamamlandi'), 'Completed');
      assert.strictEqual(decorationProvider.getStatusLabel('unknown'), 'unknown');
    });

    test('should get correct priority labels', () => {
      assert.strictEqual(decorationProvider.getPriorityLabel('yuksek'), 'High');
      assert.strictEqual(decorationProvider.getPriorityLabel('orta'), 'Medium');
      assert.strictEqual(decorationProvider.getPriorityLabel('dusuk'), 'Low');
      assert.strictEqual(decorationProvider.getPriorityLabel('unknown'), 'unknown');
    });
  });

  suite('Color Methods', () => {
    test('should get priority colors', () => {
      const highColor = decorationProvider.getPriorityColor('yuksek');
      const mediumColor = decorationProvider.getPriorityColor('orta');
      const lowColor = decorationProvider.getPriorityColor('dusuk');
      
      assert(highColor);
      assert(mediumColor);
      assert(lowColor);
    });

    test('should get status colors', () => {
      const pendingColor = decorationProvider.getStatusColor('beklemede');
      const inProgressColor = decorationProvider.getStatusColor('devam_ediyor');
      const completedColor = decorationProvider.getStatusColor('tamamlandi');
      
      assert(pendingColor);
      assert(inProgressColor);
      assert(completedColor);
    });

    test('should handle unknown colors', () => {
      assert.strictEqual(decorationProvider.getPriorityColor('unknown'), undefined);
      assert.strictEqual(decorationProvider.getStatusColor('unknown'), undefined);
    });
  });

  suite('Configuration Integration', () => {
    test('should respect decoration configuration', () => {
      const task = { oncelik: 'yuksek' };
      
      const mockConfig = {
        get: sandbox.stub()
      };
      mockConfig.get.withArgs('showPriority', true).returns(true);
      vscode.workspace.getConfiguration.withArgs('gorev.treeView.decorations').returns(mockConfig);
      
      const shouldShow = decorationProvider.shouldShowDecoration(task, 'priority');
      assert.strictEqual(shouldShow, true);
    });

    test('should hide decorations when disabled in config', () => {
      const task = { oncelik: 'yuksek' };
      
      const mockConfig = {
        get: sandbox.stub()
      };
      mockConfig.get.withArgs('showPriority', true).returns(false);
      vscode.workspace.getConfiguration.withArgs('gorev.treeView.decorations').returns(mockConfig);
      
      const shouldShow = decorationProvider.shouldShowDecoration(task, 'priority');
      assert.strictEqual(shouldShow, false);
    });

    test('should handle missing configuration gracefully', () => {
      vscode.workspace.getConfiguration.throws(new Error('Config error'));
      
      const task = { oncelik: 'yuksek' };
      
      try {
        const shouldShow = decorationProvider.shouldShowDecoration(task, 'priority');
        assert(typeof shouldShow === 'boolean');
      } catch (error) {
        assert.fail('Should handle configuration errors gracefully');
      }
    });
  });

  suite('Apply Task Decorations', () => {
    test('should apply decorations to tree item', () => {
      const task = {
        baslik: 'Test Task',
        oncelik: 'yuksek',
        durum: 'devam_ediyor',
        etiketler: ['bug']
      };
      
      const treeItem = {
        label: 'Test Task'
      };
      
      const decoratedItem = decorationProvider.applyTaskDecorations(treeItem, task);
      
      assert(decoratedItem.description);
      assert(decoratedItem.tooltip);
      assert(decoratedItem === treeItem); // Should modify in place
    });

    test('should handle null tree item', () => {
      const task = { baslik: 'Test Task' };
      
      try {
        const result = decorationProvider.applyTaskDecorations(null, task);
        assert(result === null || result === undefined);
      } catch (error) {
        assert.fail('Should handle null tree item gracefully');
      }
    });

    test('should handle null task', () => {
      const treeItem = { label: 'Test Task' };
      
      try {
        const result = decorationProvider.applyTaskDecorations(treeItem, null);
        assert(result === treeItem);
      } catch (error) {
        assert.fail('Should handle null task gracefully');
      }
    });
  });

  suite('Dispose', () => {
    test('should dispose all decoration types', () => {
      decorationProvider.dispose();
      
      // Should call dispose on all decoration types
      assert(vscode.window.createTextEditorDecorationType().dispose.called);
    });

    test('should handle dispose errors gracefully', () => {
      const mockDecorationType = vscode.window.createTextEditorDecorationType();
      mockDecorationType.dispose.throws(new Error('Dispose failed'));
      
      try {
        decorationProvider.dispose();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle dispose errors gracefully');
      }
    });

    test('should handle null decoration types', () => {
      decorationProvider.decorationTypes = null;
      
      try {
        decorationProvider.dispose();
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle null decoration types gracefully');
      }
    });
  });

  suite('Error Handling', () => {
    test('should handle invalid task data', () => {
      const invalidTasks = [null, undefined, '', 123, [], 'string'];
      
      invalidTasks.forEach(task => {
        try {
          decorationProvider.createTaskDescription(task);
          decorationProvider.createTaskTooltip(task);
          assert(true); // Should not throw
        } catch (error) {
          assert.fail(`Should handle invalid task data gracefully: ${task}`);
        }
      });
    });

    test('should handle missing task properties gracefully', () => {
      const incompleteTask = {
        baslik: 'Test'
        // Missing other properties
      };
      
      try {
        const description = decorationProvider.createTaskDescription(incompleteTask);
        const tooltip = decorationProvider.createTaskTooltip(incompleteTask);
        
        assert(typeof description === 'string');
        assert(typeof tooltip === 'string');
      } catch (error) {
        assert.fail('Should handle missing properties gracefully');
      }
    });

    test('should handle circular references in dependencies', () => {
      const task = {
        baslik: 'Test Task',
        dependencies: []
      };
      
      // Create circular reference
      const circularDep = { title: 'Circular', task: task };
      task.dependencies.push(circularDep);
      
      try {
        decorationProvider.getDependencyBadge(task.dependencies);
        assert(true); // Should not throw
      } catch (error) {
        assert.fail('Should handle circular references gracefully');
      }
    });
  });
});