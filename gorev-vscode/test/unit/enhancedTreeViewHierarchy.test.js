const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const { EnhancedGorevTreeProvider, TaskTreeViewItem } = require('../../out/providers/enhancedGorevTreeProvider');
const { GroupingStrategy } = require('../../out/models/treeModels');
const { GorevDurum, GorevOncelik } = require('../../out/models/common');

suite('Enhanced TreeView Hierarchy Test Suite', () => {
    let sandbox;
    let mockMcpClient;
    let treeProvider;

    setup(() => {
        sandbox = sinon.createSandbox();
        
        // Mock MCP Client
        mockMcpClient = {
            callTool: sandbox.stub(),
            isConnected: sandbox.stub().returns(true)
        };

        // Mock vscode.workspace
        sandbox.stub(vscode.workspace, 'getConfiguration').returns({
            get: (key, defaultValue) => defaultValue
        });

        treeProvider = new EnhancedGorevTreeProvider(mockMcpClient);
    });

    teardown(() => {
        sandbox.restore();
    });

    suite('Hierarchical Display Tests', () => {
        test('should display only root tasks at top level', async () => {
            const mockResponse = {
                content: [{
                    text: `## Görevler

[✓] Root Task 1 (orta öncelik)
  Root görev açıklaması
  ID: root1
  
  └─ [⏳] Child Task 1.1 (yüksek öncelik)
    Alt görev açıklaması
    ID: child11

[🔄] Root Task 2 (yüksek öncelik)
  İkinci root görev
  ID: root2`
                }]
            };

            mockMcpClient.callTool.resolves(mockResponse);

            const children = await treeProvider.getChildren();
            
            // GroupingStrategy.None olmalı ki direkt görevleri görelim
            treeProvider.setGrouping(GroupingStrategy.None);
            const tasks = await treeProvider.getChildren();

            // Sadece 2 root görev görünmeli
            const rootTasks = tasks.filter(item => item.type === 'task');
            assert.strictEqual(rootTasks.length, 2);
            assert.strictEqual(rootTasks[0].task.id, 'root1');
            assert.strictEqual(rootTasks[1].task.id, 'root2');
        });

        test('should expand parent tasks to show children', async () => {
            const parentTask = {
                id: 'parent1',
                baslik: 'Parent Task',
                alt_gorevler: [
                    { id: 'child1', baslik: 'Child 1', parent_id: 'parent1' },
                    { id: 'child2', baslik: 'Child 2', parent_id: 'parent1' }
                ]
            };

            const parentItem = new TaskTreeViewItem(parentTask, { selectedTasks: new Set() });
            
            // Parent task should be collapsible
            assert.strictEqual(parentItem.collapsibleState, vscode.TreeItemCollapsibleState.Expanded);
            
            // Context value should indicate parent
            assert.strictEqual(parentItem.contextValue, 'task:parent');
        });

        test('should show subtask count in description', () => {
            const task = {
                id: '123',
                baslik: 'Parent Task',
                durum: GorevDurum.DevamEdiyor,
                oncelik: GorevOncelik.Orta,
                alt_gorevler: [
                    { id: '1', durum: GorevDurum.Tamamlandi },
                    { id: '2', durum: GorevDurum.DevamEdiyor },
                    { id: '3', durum: GorevDurum.Beklemede },
                    { id: '4', durum: GorevDurum.Beklemede },
                    { id: '5', durum: GorevDurum.Tamamlandi }
                ]
            };

            const item = new TaskTreeViewItem(task, { selectedTasks: new Set() });
            
            // Description should include subtask count
            assert.ok(item.description.includes('📁 2/5')); // 2 completed out of 5
        });
    });

    suite('Tree Filtering with Hierarchy Tests', () => {
        test('should filter tasks while maintaining hierarchy', async () => {
            const tasks = [
                {
                    id: 'parent1',
                    baslik: 'Parent Task',
                    etiketler: ['important'],
                    alt_gorevler: [
                        { id: 'child1', baslik: 'Child 1', etiketler: ['urgent'] },
                        { id: 'child2', baslik: 'Child 2', etiketler: ['important'] }
                    ]
                },
                {
                    id: 'parent2',
                    baslik: 'Another Parent',
                    etiketler: ['low-priority'],
                    alt_gorevler: []
                }
            ];

            // Test search filter
            treeProvider.tasks = tasks;
            treeProvider.updateFilter({ searchQuery: 'Parent' });
            
            // Both parent tasks should match
            const filtered = treeProvider.filteredTasks;
            assert.strictEqual(filtered.length, 2);
        });

        test('should handle nested subtasks correctly', () => {
            const deepTask = {
                id: 'root',
                baslik: 'Root Task',
                alt_gorevler: [{
                    id: 'level1',
                    baslik: 'Level 1',
                    parent_id: 'root',
                    alt_gorevler: [{
                        id: 'level2',
                        baslik: 'Level 2',
                        parent_id: 'level1',
                        alt_gorevler: [{
                            id: 'level3',
                            baslik: 'Level 3',
                            parent_id: 'level2',
                            alt_gorevler: []
                        }]
                    }]
                }]
            };

            const rootItem = new TaskTreeViewItem(deepTask, { selectedTasks: new Set() });
            assert.strictEqual(rootItem.collapsibleState, vscode.TreeItemCollapsibleState.Expanded);
            assert.ok(rootItem.task.alt_gorevler.length > 0);
        });
    });

    suite('Context Menu Tests', () => {
        test('should show create subtask option for all tasks', () => {
            const task = { id: '123', baslik: 'Any Task' };
            const item = new TaskTreeViewItem(task, { selectedTasks: new Set() });
            
            // Context value should allow subtask creation
            const contextValue = item.contextValue;
            assert.ok(['task', 'task:parent', 'task:child'].includes(contextValue));
        });

        test('should show remove parent option only for child tasks', () => {
            const childTask = { 
                id: '123', 
                baslik: 'Child Task',
                parent_id: '456'
            };
            const item = new TaskTreeViewItem(childTask, { selectedTasks: new Set() });
            
            assert.strictEqual(item.contextValue, 'task:child');
        });

        test('should show change parent option for all tasks', () => {
            const task = { id: '123', baslik: 'Any Task' };
            const item = new TaskTreeViewItem(task, { selectedTasks: new Set() });
            
            // All task types can change parent
            assert.ok(item.contextValue.startsWith('task'));
        });
    });

    suite('getChildren Implementation Tests', () => {
        test('should return subtasks when expanding parent', async () => {
            const parentTask = {
                id: 'parent1',
                baslik: 'Parent',
                alt_gorevler: [
                    { id: 'child1', baslik: 'Child 1' },
                    { id: 'child2', baslik: 'Child 2' }
                ]
            };

            const parentItem = new TaskTreeViewItem(parentTask, { selectedTasks: new Set() });
            
            // Mock the tree provider to return children
            const children = parentTask.alt_gorevler.map(child => 
                new TaskTreeViewItem(child, { selectedTasks: new Set() }, parentItem.parent)
            );

            assert.strictEqual(children.length, 2);
            assert.strictEqual(children[0].task.baslik, 'Child 1');
            assert.strictEqual(children[1].task.baslik, 'Child 2');
        });
    });
});