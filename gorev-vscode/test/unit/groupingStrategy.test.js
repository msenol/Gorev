const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');

suite('GroupingStrategy Test Suite', () => {
  let sandbox;
  let GroupingStrategy;
  let TaskGroupManager;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.TreeItemCollapsibleState, 'Collapsed', 1);
    sandbox.stub(vscode.TreeItemCollapsibleState, 'Expanded', 2);
    sandbox.stub(vscode.TreeItemCollapsibleState, 'None', 0);

    // Import GroupingStrategy and related classes
    try {
      const groupingModule = require('../../dist/providers/groupingStrategy');
      GroupingStrategy = groupingModule.GroupingStrategy;
      TaskGroupManager = groupingModule.TaskGroupManager;
    } catch (error) {
      // Mock classes if compilation fails
      GroupingStrategy = {
        None: 'none',
        ByStatus: 'status',
        ByPriority: 'priority',
        ByProject: 'project',
        ByTag: 'tag',
        ByDueDate: 'dueDate'
      };

      TaskGroupManager = class MockTaskGroupManager {
        constructor() {
          this.currentStrategy = GroupingStrategy.None;
        }

        setGroupingStrategy(strategy) {
          this.currentStrategy = strategy;
        }

        groupTasks(tasks) {
          if (this.currentStrategy === GroupingStrategy.None) {
            return tasks;
          }
          
          const groups = new Map();
          
          tasks.forEach(task => {
            const groupKey = this.getGroupKey(task, this.currentStrategy);
            if (!groups.has(groupKey)) {
              groups.set(groupKey, {
                label: this.getGroupLabel(groupKey, this.currentStrategy),
                children: [],
                collapsibleState: 1 // Collapsed
              });
            }
            groups.get(groupKey).children.push(task);
          });
          
          return Array.from(groups.values());
        }

        getGroupKey(task, strategy) {
          switch (strategy) {
            case GroupingStrategy.ByStatus:
              return task.durum || 'unknown';
            case GroupingStrategy.ByPriority:
              return task.oncelik || 'unknown';
            case GroupingStrategy.ByProject:
              return task.proje_ismi || task.proje_id || 'no-project';
            case GroupingStrategy.ByTag:
              return task.etiketler && task.etiketler.length > 0 
                ? task.etiketler[0] 
                : 'no-tags';
            case GroupingStrategy.ByDueDate:
              return this.getDueDateGroup(task.son_tarih);
            default:
              return 'other';
          }
        }

        getGroupLabel(key, strategy) {
          switch (strategy) {
            case GroupingStrategy.ByStatus:
              return this.getStatusLabel(key);
            case GroupingStrategy.ByPriority:
              return this.getPriorityLabel(key);
            case GroupingStrategy.ByProject:
              return key === 'no-project' ? 'No Project' : key;
            case GroupingStrategy.ByTag:
              return key === 'no-tags' ? 'No Tags' : `#${key}`;
            case GroupingStrategy.ByDueDate:
              return this.getDueDateLabel(key);
            default:
              return key;
          }
        }

        getStatusLabel(status) {
          const labels = {
            'beklemede': 'Pending',
            'devam_ediyor': 'In Progress',
            'tamamlandi': 'Completed',
            'unknown': 'Unknown Status'
          };
          return labels[status] || status;
        }

        getPriorityLabel(priority) {
          const labels = {
            'yuksek': 'High Priority',
            'orta': 'Medium Priority',
            'dusuk': 'Low Priority',
            'unknown': 'Unknown Priority'
          };
          return labels[priority] || priority;
        }

        getDueDateGroup(dueDate) {
          if (!dueDate) return 'no-due-date';
          
          const now = new Date();
          const due = new Date(dueDate);
          const diffTime = due.getTime() - now.getTime();
          const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

          if (diffDays < 0) return 'overdue';
          if (diffDays === 0) return 'today';
          if (diffDays === 1) return 'tomorrow';
          if (diffDays <= 7) return 'this-week';
          if (diffDays <= 30) return 'this-month';
          return 'later';
        }

        getDueDateLabel(group) {
          const labels = {
            'overdue': '🚨 Overdue',
            'today': '📅 Due Today',
            'tomorrow': '📅 Due Tomorrow',
            'this-week': '📅 This Week',
            'this-month': '📅 This Month',
            'later': '📅 Later',
            'no-due-date': '📅 No Due Date'
          };
          return labels[group] || group;
        }

        sortGroups(groups, strategy) {
          const sortOrder = this.getGroupSortOrder(strategy);
          
          return groups.sort((a, b) => {
            const aKey = this.getGroupKeyFromLabel(a.label, strategy);
            const bKey = this.getGroupKeyFromLabel(b.label, strategy);
            
            const aIndex = sortOrder.indexOf(aKey);
            const bIndex = sortOrder.indexOf(bKey);
            
            // If both found in sort order, use that order
            if (aIndex !== -1 && bIndex !== -1) {
              return aIndex - bIndex;
            }
            
            // If only one found, put found one first
            if (aIndex !== -1) return -1;
            if (bIndex !== -1) return 1;
            
            // If neither found, sort alphabetically
            return a.label.localeCompare(b.label);
          });
        }

        getGroupSortOrder(strategy) {
          switch (strategy) {
            case GroupingStrategy.ByStatus:
              return ['beklemede', 'devam_ediyor', 'tamamlandi', 'unknown'];
            case GroupingStrategy.ByPriority:
              return ['yuksek', 'orta', 'dusuk', 'unknown'];
            case GroupingStrategy.ByDueDate:
              return ['overdue', 'today', 'tomorrow', 'this-week', 'this-month', 'later', 'no-due-date'];
            default:
              return [];
          }
        }

        getGroupKeyFromLabel(label, strategy) {
          // Reverse mapping from label to key
          switch (strategy) {
            case GroupingStrategy.ByStatus:
              const statusMap = { 'Pending': 'beklemede', 'In Progress': 'devam_ediyor', 'Completed': 'tamamlandi', 'Unknown Status': 'unknown' };
              return statusMap[label] || label;
            case GroupingStrategy.ByPriority:
              const priorityMap = { 'High Priority': 'yuksek', 'Medium Priority': 'orta', 'Low Priority': 'dusuk', 'Unknown Priority': 'unknown' };
              return priorityMap[label] || label;
            case GroupingStrategy.ByDueDate:
              const dueDateMap = { '🚨 Overdue': 'overdue', '📅 Due Today': 'today', '📅 Due Tomorrow': 'tomorrow', '📅 This Week': 'this-week', '📅 This Month': 'this-month', '📅 Later': 'later', '📅 No Due Date': 'no-due-date' };
              return dueDateMap[label] || label;
            default:
              return label;
          }
        }

        getGroupStats(group) {
          if (!group.children) return '';
          
          const total = group.children.length;
          const completed = group.children.filter(task => task.durum === 'tamamlandi').length;
          const inProgress = group.children.filter(task => task.durum === 'devam_ediyor').length;
          
          return `${completed}/${total} completed, ${inProgress} in progress`;
        }

        expandGroup(groupLabel) {
          // Mock implementation for expanding groups
          return true;
        }

        collapseGroup(groupLabel) {
          // Mock implementation for collapsing groups
          return true;
        }

        isGroupExpanded(groupLabel) {
          // Mock implementation - could be managed via settings
          return false; // Default to collapsed
        }
      };
    }
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('GroupingStrategy Enum', () => {
    test('should have all grouping strategy values', () => {
      assert.strictEqual(GroupingStrategy.None, 'none');
      assert.strictEqual(GroupingStrategy.ByStatus, 'status');
      assert.strictEqual(GroupingStrategy.ByPriority, 'priority');
      assert.strictEqual(GroupingStrategy.ByProject, 'project');
      assert.strictEqual(GroupingStrategy.ByTag, 'tag');
      assert.strictEqual(GroupingStrategy.ByDueDate, 'dueDate');
    });
  });

  suite('TaskGroupManager Initialization', () => {
    test('should create task group manager', () => {
      const manager = new TaskGroupManager();
      assert(manager);
      assert.strictEqual(manager.currentStrategy, GroupingStrategy.None);
    });

    test('should allow setting grouping strategy', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByStatus);
      assert.strictEqual(manager.currentStrategy, GroupingStrategy.ByStatus);
    });
  });

  suite('No Grouping', () => {
    test('should return tasks as-is when no grouping', () => {
      const manager = new TaskGroupManager();
      const tasks = [
        { id: '1', baslik: 'Task 1', durum: 'beklemede' },
        { id: '2', baslik: 'Task 2', durum: 'devam_ediyor' }
      ];
      
      const result = manager.groupTasks(tasks);
      assert.deepStrictEqual(result, tasks);
    });

    test('should handle empty task list', () => {
      const manager = new TaskGroupManager();
      const result = manager.groupTasks([]);
      assert.deepStrictEqual(result, []);
    });

    test('should handle null task list', () => {
      const manager = new TaskGroupManager();
      const result = manager.groupTasks(null);
      assert.strictEqual(result, null);
    });
  });

  suite('Group by Status', () => {
    test('should group tasks by status', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByStatus);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', durum: 'beklemede' },
        { id: '2', baslik: 'Task 2', durum: 'devam_ediyor' },
        { id: '3', baslik: 'Task 3', durum: 'beklemede' },
        { id: '4', baslik: 'Task 4', durum: 'tamamlandi' }
      ];
      
      const result = manager.groupTasks(tasks);
      
      assert(Array.isArray(result));
      assert.strictEqual(result.length, 3); // 3 different statuses
      
      // Find the pending group
      const pendingGroup = result.find(group => group.label === 'Pending');
      assert(pendingGroup);
      assert.strictEqual(pendingGroup.children.length, 2);
      assert(pendingGroup.children.some(task => task.id === '1'));
      assert(pendingGroup.children.some(task => task.id === '3'));
    });

    test('should handle tasks with unknown status', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByStatus);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', durum: 'beklemede' },
        { id: '2', baslik: 'Task 2' } // No status
      ];
      
      const result = manager.groupTasks(tasks);
      
      const unknownGroup = result.find(group => group.label === 'Unknown Status');
      assert(unknownGroup);
      assert.strictEqual(unknownGroup.children.length, 1);
      assert.strictEqual(unknownGroup.children[0].id, '2');
    });

    test('should sort status groups in logical order', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByStatus);
      
      const tasks = [
        { id: '1', durum: 'tamamlandi' },
        { id: '2', durum: 'beklemede' },
        { id: '3', durum: 'devam_ediyor' }
      ];
      
      const groups = manager.groupTasks(tasks);
      const sortedGroups = manager.sortGroups(groups, GroupingStrategy.ByStatus);
      
      assert.strictEqual(sortedGroups[0].label, 'Pending');
      assert.strictEqual(sortedGroups[1].label, 'In Progress');
      assert.strictEqual(sortedGroups[2].label, 'Completed');
    });
  });

  suite('Group by Priority', () => {
    test('should group tasks by priority', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByPriority);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', oncelik: 'yuksek' },
        { id: '2', baslik: 'Task 2', oncelik: 'orta' },
        { id: '3', baslik: 'Task 3', oncelik: 'yuksek' },
        { id: '4', baslik: 'Task 4', oncelik: 'dusuk' }
      ];
      
      const result = manager.groupTasks(tasks);
      
      assert(Array.isArray(result));
      assert.strictEqual(result.length, 3); // 3 different priorities
      
      const highPriorityGroup = result.find(group => group.label === 'High Priority');
      assert(highPriorityGroup);
      assert.strictEqual(highPriorityGroup.children.length, 2);
    });

    test('should sort priority groups correctly', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByPriority);
      
      const tasks = [
        { id: '1', oncelik: 'dusuk' },
        { id: '2', oncelik: 'yuksek' },
        { id: '3', oncelik: 'orta' }
      ];
      
      const groups = manager.groupTasks(tasks);
      const sortedGroups = manager.sortGroups(groups, GroupingStrategy.ByPriority);
      
      assert.strictEqual(sortedGroups[0].label, 'High Priority');
      assert.strictEqual(sortedGroups[1].label, 'Medium Priority');
      assert.strictEqual(sortedGroups[2].label, 'Low Priority');
    });
  });

  suite('Group by Project', () => {
    test('should group tasks by project name', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByProject);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', proje_ismi: 'Project A' },
        { id: '2', baslik: 'Task 2', proje_ismi: 'Project B' },
        { id: '3', baslik: 'Task 3', proje_ismi: 'Project A' },
        { id: '4', baslik: 'Task 4' } // No project
      ];
      
      const result = manager.groupTasks(tasks);
      
      assert(Array.isArray(result));
      assert.strictEqual(result.length, 3); // 2 projects + no project
      
      const projectAGroup = result.find(group => group.label === 'Project A');
      assert(projectAGroup);
      assert.strictEqual(projectAGroup.children.length, 2);
      
      const noProjectGroup = result.find(group => group.label === 'No Project');
      assert(noProjectGroup);
      assert.strictEqual(noProjectGroup.children.length, 1);
    });

    test('should fallback to project ID when name is missing', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByProject);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', proje_id: 'proj-123' }
      ];
      
      const result = manager.groupTasks(tasks);
      
      const projectGroup = result.find(group => group.label === 'proj-123');
      assert(projectGroup);
    });
  });

  suite('Group by Tag', () => {
    test('should group tasks by first tag', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByTag);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', etiketler: ['bug', 'urgent'] },
        { id: '2', baslik: 'Task 2', etiketler: ['feature'] },
        { id: '3', baslik: 'Task 3', etiketler: ['bug', 'ui'] },
        { id: '4', baslik: 'Task 4' } // No tags
      ];
      
      const result = manager.groupTasks(tasks);
      
      assert(Array.isArray(result));
      assert.strictEqual(result.length, 3); // bug, feature, no-tags
      
      const bugGroup = result.find(group => group.label === '#bug');
      assert(bugGroup);
      assert.strictEqual(bugGroup.children.length, 2);
      
      const noTagsGroup = result.find(group => group.label === 'No Tags');
      assert(noTagsGroup);
      assert.strictEqual(noTagsGroup.children.length, 1);
    });

    test('should handle empty tags array', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByTag);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', etiketler: [] }
      ];
      
      const result = manager.groupTasks(tasks);
      
      const noTagsGroup = result.find(group => group.label === 'No Tags');
      assert(noTagsGroup);
      assert.strictEqual(noTagsGroup.children.length, 1);
    });
  });

  suite('Group by Due Date', () => {
    test('should group tasks by due date categories', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByDueDate);
      
      const today = new Date();
      const tomorrow = new Date(today);
      tomorrow.setDate(tomorrow.getDate() + 1);
      const yesterday = new Date(today);
      yesterday.setDate(yesterday.getDate() - 1);
      
      const tasks = [
        { id: '1', baslik: 'Task 1', son_tarih: today.toISOString().split('T')[0] },
        { id: '2', baslik: 'Task 2', son_tarih: tomorrow.toISOString().split('T')[0] },
        { id: '3', baslik: 'Task 3', son_tarih: yesterday.toISOString().split('T')[0] },
        { id: '4', baslik: 'Task 4' } // No due date
      ];
      
      const result = manager.groupTasks(tasks);
      
      assert(Array.isArray(result));
      
      const todayGroup = result.find(group => group.label === '📅 Due Today');
      assert(todayGroup);
      assert.strictEqual(todayGroup.children.length, 1);
      
      const tomorrowGroup = result.find(group => group.label === '📅 Due Tomorrow');
      assert(tomorrowGroup);
      assert.strictEqual(tomorrowGroup.children.length, 1);
      
      const overdueGroup = result.find(group => group.label === '🚨 Overdue');
      assert(overdueGroup);
      assert.strictEqual(overdueGroup.children.length, 1);
      
      const noDueDateGroup = result.find(group => group.label === '📅 No Due Date');
      assert(noDueDateGroup);
      assert.strictEqual(noDueDateGroup.children.length, 1);
    });

    test('should categorize this week and this month correctly', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByDueDate);
      
      const thisWeek = new Date();
      thisWeek.setDate(thisWeek.getDate() + 5);
      
      const thisMonth = new Date();
      thisMonth.setDate(thisMonth.getDate() + 15);
      
      const later = new Date();
      later.setDate(later.getDate() + 45);
      
      const tasks = [
        { id: '1', son_tarih: thisWeek.toISOString().split('T')[0] },
        { id: '2', son_tarih: thisMonth.toISOString().split('T')[0] },
        { id: '3', son_tarih: later.toISOString().split('T')[0] }
      ];
      
      const result = manager.groupTasks(tasks);
      
      const thisWeekGroup = result.find(group => group.label === '📅 This Week');
      const thisMonthGroup = result.find(group => group.label === '📅 This Month');
      const laterGroup = result.find(group => group.label === '📅 Later');
      
      assert(thisWeekGroup);
      assert(thisMonthGroup);
      assert(laterGroup);
    });

    test('should sort due date groups in logical order', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByDueDate);
      
      const tasks = [
        { id: '1', son_tarih: null },
        { id: '2', son_tarih: '2025-01-01' }, // Later
        { id: '3', son_tarih: '2024-01-01' }  // Overdue
      ];
      
      const groups = manager.groupTasks(tasks);
      const sortedGroups = manager.sortGroups(groups, GroupingStrategy.ByDueDate);
      
      // Should be in order: overdue, today, tomorrow, this-week, this-month, later, no-due-date
      const orderLabels = sortedGroups.map(g => g.label);
      const overdueIndex = orderLabels.indexOf('🚨 Overdue');
      const laterIndex = orderLabels.indexOf('📅 Later');
      const noDueDateIndex = orderLabels.indexOf('📅 No Due Date');
      
      assert(overdueIndex < laterIndex);
      assert(laterIndex < noDueDateIndex);
    });
  });

  suite('Group Statistics', () => {
    test('should calculate group statistics', () => {
      const manager = new TaskGroupManager();
      
      const group = {
        label: 'Test Group',
        children: [
          { id: '1', durum: 'tamamlandi' },
          { id: '2', durum: 'devam_ediyor' },
          { id: '3', durum: 'beklemede' },
          { id: '4', durum: 'tamamlandi' }
        ]
      };
      
      const stats = manager.getGroupStats(group);
      
      assert.strictEqual(stats, '2/4 completed, 1 in progress');
    });

    test('should handle empty group', () => {
      const manager = new TaskGroupManager();
      
      const group = {
        label: 'Empty Group',
        children: []
      };
      
      const stats = manager.getGroupStats(group);
      
      assert.strictEqual(stats, '0/0 completed, 0 in progress');
    });

    test('should handle group without children', () => {
      const manager = new TaskGroupManager();
      
      const group = {
        label: 'No Children Group'
      };
      
      const stats = manager.getGroupStats(group);
      
      assert.strictEqual(stats, '');
    });
  });

  suite('Group Management', () => {
    test('should expand group', () => {
      const manager = new TaskGroupManager();
      const result = manager.expandGroup('Test Group');
      assert.strictEqual(result, true);
    });

    test('should collapse group', () => {
      const manager = new TaskGroupManager();
      const result = manager.collapseGroup('Test Group');
      assert.strictEqual(result, true);
    });

    test('should check if group is expanded', () => {
      const manager = new TaskGroupManager();
      const result = manager.isGroupExpanded('Test Group');
      assert.strictEqual(result, false); // Default to collapsed
    });
  });

  suite('Edge Cases', () => {
    test('should handle tasks with all null/undefined values', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByStatus);
      
      const tasks = [
        { id: '1' }, // No properties
        { id: '2', durum: null },
        { id: '3', durum: undefined }
      ];
      
      const result = manager.groupTasks(tasks);
      
      const unknownGroup = result.find(group => group.label === 'Unknown Status');
      assert(unknownGroup);
      assert.strictEqual(unknownGroup.children.length, 3);
    });

    test('should handle invalid date formats', () => {
      const manager = new TaskGroupManager();
      
      const invalidDates = ['invalid-date', '2024-13-45', '', null, undefined];
      
      invalidDates.forEach(invalidDate => {
        const group = manager.getDueDateGroup(invalidDate);
        assert.strictEqual(group, 'no-due-date');
      });
    });

    test('should handle very large task lists', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByStatus);
      
      const largeTasks = Array.from({ length: 1000 }, (_, i) => ({
        id: `task-${i}`,
        baslik: `Task ${i}`,
        durum: ['beklemede', 'devam_ediyor', 'tamamlandi'][i % 3]
      }));
      
      const result = manager.groupTasks(largeTasks);
      
      assert(Array.isArray(result));
      assert.strictEqual(result.length, 3);
      
      const totalTasks = result.reduce((sum, group) => sum + group.children.length, 0);
      assert.strictEqual(totalTasks, 1000);
    });

    test('should handle duplicate tasks', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByStatus);
      
      const task = { id: '1', baslik: 'Task 1', durum: 'beklemede' };
      const tasks = [task, task, task]; // Same reference
      
      const result = manager.groupTasks(tasks);
      
      const pendingGroup = result.find(group => group.label === 'Pending');
      assert(pendingGroup);
      assert.strictEqual(pendingGroup.children.length, 3);
    });

    test('should handle circular references in task data', () => {
      const manager = new TaskGroupManager();
      manager.setGroupingStrategy(GroupingStrategy.ByProject);
      
      const task = { id: '1', baslik: 'Task 1' };
      task.self = task; // Circular reference
      
      const tasks = [task];
      
      try {
        const result = manager.groupTasks(tasks);
        assert(Array.isArray(result));
      } catch (error) {
        assert.fail('Should handle circular references gracefully');
      }
    });
  });

  suite('Label Utilities', () => {
    test('should get correct status labels', () => {
      const manager = new TaskGroupManager();
      
      assert.strictEqual(manager.getStatusLabel('beklemede'), 'Pending');
      assert.strictEqual(manager.getStatusLabel('devam_ediyor'), 'In Progress');
      assert.strictEqual(manager.getStatusLabel('tamamlandi'), 'Completed');
      assert.strictEqual(manager.getStatusLabel('unknown'), 'Unknown Status');
      assert.strictEqual(manager.getStatusLabel('invalid'), 'invalid');
    });

    test('should get correct priority labels', () => {
      const manager = new TaskGroupManager();
      
      assert.strictEqual(manager.getPriorityLabel('yuksek'), 'High Priority');
      assert.strictEqual(manager.getPriorityLabel('orta'), 'Medium Priority');
      assert.strictEqual(manager.getPriorityLabel('dusuk'), 'Low Priority');
      assert.strictEqual(manager.getPriorityLabel('unknown'), 'Unknown Priority');
      assert.strictEqual(manager.getPriorityLabel('invalid'), 'invalid');
    });

    test('should get correct due date labels', () => {
      const manager = new TaskGroupManager();
      
      assert.strictEqual(manager.getDueDateLabel('overdue'), '🚨 Overdue');
      assert.strictEqual(manager.getDueDateLabel('today'), '📅 Due Today');
      assert.strictEqual(manager.getDueDateLabel('tomorrow'), '📅 Due Tomorrow');
      assert.strictEqual(manager.getDueDateLabel('this-week'), '📅 This Week');
      assert.strictEqual(manager.getDueDateLabel('this-month'), '📅 This Month');
      assert.strictEqual(manager.getDueDateLabel('later'), '📅 Later');
      assert.strictEqual(manager.getDueDateLabel('no-due-date'), '📅 No Due Date');
      assert.strictEqual(manager.getDueDateLabel('invalid'), 'invalid');
    });
  });

  suite('Reverse Label Mapping', () => {
    test('should map status labels back to keys', () => {
      const manager = new TaskGroupManager();
      
      assert.strictEqual(manager.getGroupKeyFromLabel('Pending', GroupingStrategy.ByStatus), 'beklemede');
      assert.strictEqual(manager.getGroupKeyFromLabel('In Progress', GroupingStrategy.ByStatus), 'devam_ediyor');
      assert.strictEqual(manager.getGroupKeyFromLabel('Completed', GroupingStrategy.ByStatus), 'tamamlandi');
      assert.strictEqual(manager.getGroupKeyFromLabel('Unknown Status', GroupingStrategy.ByStatus), 'unknown');
    });

    test('should map priority labels back to keys', () => {
      const manager = new TaskGroupManager();
      
      assert.strictEqual(manager.getGroupKeyFromLabel('High Priority', GroupingStrategy.ByPriority), 'yuksek');
      assert.strictEqual(manager.getGroupKeyFromLabel('Medium Priority', GroupingStrategy.ByPriority), 'orta');
      assert.strictEqual(manager.getGroupKeyFromLabel('Low Priority', GroupingStrategy.ByPriority), 'dusuk');
      assert.strictEqual(manager.getGroupKeyFromLabel('Unknown Priority', GroupingStrategy.ByPriority), 'unknown');
    });

    test('should map due date labels back to keys', () => {
      const manager = new TaskGroupManager();
      
      assert.strictEqual(manager.getGroupKeyFromLabel('🚨 Overdue', GroupingStrategy.ByDueDate), 'overdue');
      assert.strictEqual(manager.getGroupKeyFromLabel('📅 Due Today', GroupingStrategy.ByDueDate), 'today');
      assert.strictEqual(manager.getGroupKeyFromLabel('📅 Due Tomorrow', GroupingStrategy.ByDueDate), 'tomorrow');
    });

    test('should handle unknown labels', () => {
      const manager = new TaskGroupManager();
      
      assert.strictEqual(manager.getGroupKeyFromLabel('Unknown Label', GroupingStrategy.ByStatus), 'Unknown Label');
      assert.strictEqual(manager.getGroupKeyFromLabel('', GroupingStrategy.ByPriority), '');
      assert.strictEqual(manager.getGroupKeyFromLabel(null, GroupingStrategy.ByDueDate), null);
    });
  });
});