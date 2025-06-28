const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const { GorevTreeProvider } = require('../../dist/providers/gorevTreeProvider');
const { ProjeTreeProvider } = require('../../dist/providers/projeTreeProvider');
const { TemplateTreeProvider } = require('../../dist/providers/templateTreeProvider');

suite('TreeProviders Test Suite', () => {
  let sandbox;
  let mockClient;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Create mock MCP client
    mockClient = {
      isConnected: sinon.stub().returns(true),
      callTool: sinon.stub()
    };
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('GorevTreeProvider', () => {
    let provider;

    setup(() => {
      provider = new GorevTreeProvider(mockClient);
    });

    test('should return empty array when not connected', async () => {
      mockClient.isConnected.returns(false);
      
      const children = await provider.getChildren();
      
      assert.strictEqual(children.length, 0);
    });

    test('should parse and return tasks', async () => {
      const mockResponse = {
        content: [{
          text: `- [beklemede] Test Task (orta öncelik)
  ID: task-123
  Proje: Test Project
  Test description`
        }]
      };
      
      mockClient.callTool.resolves(mockResponse);
      
      const children = await provider.getChildren();
      
      assert(children.length > 0);
      assert(mockClient.callTool.calledOnce);
      assert(mockClient.callTool.calledWith('gorev_listele'));
    });

    test('should handle errors gracefully', async () => {
      mockClient.callTool.rejects(new Error('Test error'));
      
      const children = await provider.getChildren();
      
      assert.strictEqual(children.length, 0);
    });

    test('should refresh tree data', () => {
      const eventFired = sinon.stub();
      provider.onDidChangeTreeData(eventFired);
      
      provider.refresh();
      
      assert(eventFired.called);
    });
  });

  suite('ProjeTreeProvider', () => {
    let provider;

    setup(() => {
      provider = new ProjeTreeProvider(mockClient);
    });

    test('should return projects with task counts', async () => {
      const mockResponse = {
        content: [{
          text: `### 🔒 Test Project
**ID:** proj-123
**Tanım:** Test description
**Görev Sayısı:** Toplam: 5, Tamamlanan: 2, Devam Eden: 1, Bekleyen: 2`
        }]
      };
      
      mockClient.callTool.resolves(mockResponse);
      
      const children = await provider.getChildren();
      
      assert(children.length > 0);
      assert(mockClient.callTool.calledWith('proje_listele'));
    });

    test('should create tree item with correct properties', () => {
      const mockProject = {
        id: 'proj-123',
        isim: 'Test Project',
        tanim: 'Test description'
      };
      
      const item = provider.getTreeItem({ project: mockProject });
      
      assert.strictEqual(item.label, 'Test Project');
      assert.strictEqual(item.description, 'Test description');
      assert.strictEqual(item.contextValue, 'proje');
    });
  });

  suite('TemplateTreeProvider', () => {
    let provider;

    setup(() => {
      provider = new TemplateTreeProvider(mockClient);
    });

    test('should return categories at root level', async () => {
      const mockResponse = {
        content: [{
          text: `### Teknik

#### Bug Report
- **ID:** \`template-123\`
- **Açıklama:** Bug report template

### Özellik

#### Feature Request
- **ID:** \`template-456\`
- **Açıklama:** Feature request template`
        }]
      };
      
      mockClient.callTool.resolves(mockResponse);
      
      const children = await provider.getChildren();
      
      assert(children.length > 0);
      assert(children.every(child => child.collapsibleState === vscode.TreeItemCollapsibleState.Expanded));
    });

    test('should return templates for category', async () => {
      // First load templates
      const mockResponse = {
        content: [{
          text: `### Teknik

#### Bug Report
- **ID:** \`template-123\`
- **Açıklama:** Bug report template`
        }]
      };
      
      mockClient.callTool.resolves(mockResponse);
      await provider.getChildren(); // Load templates
      
      // Now get children of category
      const categoryItem = { category: 'Teknik' };
      const children = await provider.getChildren(categoryItem);
      
      assert(children.length > 0);
      assert(children[0].template);
    });
  });

  suite('Tree item creation', () => {
    test('should create task tree item with priority icon', () => {
      const { TaskTreeViewItem } = require('../../dist/providers/enhancedGorevTreeProvider');
      
      const task = {
        id: 'task-123',
        baslik: 'Test Task',
        durum: 'beklemede',
        oncelik: 'yuksek',
        aciklama: 'Test description'
      };
      
      const item = new TaskTreeViewItem(task);
      
      assert.strictEqual(item.label, 'Test Task');
      assert.strictEqual(item.description, 'beklemede');
      assert.strictEqual(item.tooltip.includes('Yüksek öncelik'), true);
      assert.strictEqual(item.contextValue, 'gorev');
    });

    test('should create group tree item', () => {
      const { GroupTreeViewItem } = require('../../dist/providers/enhancedGorevTreeProvider');
      
      const item = new GroupTreeViewItem('Test Group', 5);
      
      assert.strictEqual(item.label, 'Test Group');
      assert.strictEqual(item.description, '5');
      assert.strictEqual(item.collapsibleState, vscode.TreeItemCollapsibleState.Expanded);
    });
  });
});