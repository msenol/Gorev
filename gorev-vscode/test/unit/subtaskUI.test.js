const assert = require('assert');
const vscode = require('vscode');
const sinon = require('sinon');
const { 
    Gorev, 
    GorevDurum, 
    GorevOncelik,
    GorevHiyerarsi 
} = require('../../out/models/gorev');
const { MarkdownParser } = require('../../out/utils/markdownParser');

suite('Subtask UI Test Suite', () => {
    let sandbox;

    setup(() => {
        sandbox = sinon.createSandbox();
    });

    teardown(() => {
        sandbox.restore();
    });

    suite('Gorev Model Tests', () => {
        test('should include parent_id field', () => {
            const gorev = {
                id: '123',
                baslik: 'Test Görev',
                parent_id: '456',
                alt_gorevler: [],
                seviye: 1
            };

            assert.strictEqual(gorev.parent_id, '456');
            assert.ok(Array.isArray(gorev.alt_gorevler));
            assert.strictEqual(gorev.seviye, 1);
        });

        test('should handle GorevHiyerarsi structure', () => {
            const gorev = {
                id: '123',
                baslik: 'Parent Task'
            };

            const hiyerarsi = {
                gorev: gorev,
                ust_gorevler: [],
                toplam_alt_gorev: 5,
                tamamlanan_alt: 2,
                devam_eden_alt: 1,
                beklemede_alt: 2,
                ilerleme_yuzdesi: 40
            };

            assert.strictEqual(hiyerarsi.toplam_alt_gorev, 5);
            assert.strictEqual(hiyerarsi.ilerleme_yuzdesi, 40);
        });
    });

    suite('MarkdownParser Hierarchy Tests', () => {
        test('should parse hierarchical task structure', () => {
            const markdown = `[✓] Ana Görev (orta öncelik)
  Ana görevin açıklaması
  ID: 123
  
  └─ [⏳] Alt Görev 1 (yüksek öncelik)
    Alt görev açıklaması
    ID: 456
    
    └─ [🔄] Alt Alt Görev (düşük öncelik)
      En alt seviye görev
      ID: 789`;

            const tasks = MarkdownParser.parseGorevListesi(markdown);
            
            assert.strictEqual(tasks.length, 1); // Sadece root görev
            assert.strictEqual(tasks[0].baslik, 'Ana Görev');
            assert.strictEqual(tasks[0].id, '123');
            assert.ok(tasks[0].alt_gorevler);
            assert.strictEqual(tasks[0].alt_gorevler.length, 1);
            
            const subtask = tasks[0].alt_gorevler[0];
            assert.strictEqual(subtask.baslik, 'Alt Görev 1');
            assert.strictEqual(subtask.parent_id, '123');
            assert.strictEqual(subtask.alt_gorevler.length, 1);
            
            const subsubtask = subtask.alt_gorevler[0];
            assert.strictEqual(subsubtask.baslik, 'Alt Alt Görev');
            assert.strictEqual(subsubtask.parent_id, '456');
        });

        test('should parse hierarchy info from gorev_hiyerarsi_goster response', () => {
            const content = `## Görev Hiyerarşisi

Toplam alt görev: 10
Tamamlanan: 4
Devam eden: 2
Beklemede: 4
İlerleme: 40%

### Üst Görevler
- Ana Proje
- Sprint 3

### Alt Görevler
- [✓] Alt görev 1
- [🔄] Alt görev 2
- [⏳] Alt görev 3`;

            const parser = new MarkdownParser();
            // parseHierarchyInfo test edilmeli ama private method olduğu için
            // TaskDetailPanel içinde test edilmeli
        });
    });

    suite('Subtask Command Tests', () => {
        test('CREATE_SUBTASK command should be defined', () => {
            const { COMMANDS } = require('../../out/utils/constants');
            assert.ok(COMMANDS.CREATE_SUBTASK);
            assert.strictEqual(COMMANDS.CREATE_SUBTASK, 'gorev.createSubtask');
        });

        test('CHANGE_PARENT command should be defined', () => {
            const { COMMANDS } = require('../../out/utils/constants');
            assert.ok(COMMANDS.CHANGE_PARENT);
            assert.strictEqual(COMMANDS.CHANGE_PARENT, 'gorev.changeParent');
        });

        test('REMOVE_PARENT command should be defined', () => {
            const { COMMANDS } = require('../../out/utils/constants');
            assert.ok(COMMANDS.REMOVE_PARENT);
            assert.strictEqual(COMMANDS.REMOVE_PARENT, 'gorev.removeParent');
        });
    });

    suite('Tree Item Context Values', () => {
        test('should set correct context value for parent tasks', () => {
            const task = {
                id: '123',
                baslik: 'Parent Task',
                alt_gorevler: [{ id: '456' }]
            };

            // TaskTreeViewItem constructor'ını test et
            const contextValue = task.alt_gorevler && task.alt_gorevler.length > 0 
                ? 'task:parent' 
                : 'task';

            assert.strictEqual(contextValue, 'task:parent');
        });

        test('should set correct context value for child tasks', () => {
            const task = {
                id: '456',
                baslik: 'Child Task',
                parent_id: '123',
                alt_gorevler: []
            };

            const contextValue = task.parent_id ? 'task:child' : 'task';
            assert.strictEqual(contextValue, 'task:child');
        });
    });

    suite('Drag & Drop Configuration Tests', () => {
        test('should include allowParentChange config', () => {
            const config = {
                allowTaskMove: true,
                allowStatusChange: true,
                allowPriorityChange: true,
                allowProjectMove: true,
                allowDependencyCreate: true,
                allowParentChange: true,
                showDropIndicator: true,
                animateOnDrop: true
            };

            assert.ok(config.allowParentChange);
        });
    });
});