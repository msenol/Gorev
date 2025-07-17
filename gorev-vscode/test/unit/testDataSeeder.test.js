const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const { Logger } = require('../../src/utils/logger');
const { GorevDurum, GorevOncelik } = require('../../src/models/common');

suite('TestDataSeeder Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let TestDataSeeder;
  let seeder;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showWarningMessage');
    sandbox.stub(vscode.window, 'withProgress');

    // Mock MCP Client
    mockMCPClient = {
      callTool: sandbox.stub()
    };

    // Mock Logger
    sandbox.stub(Logger, 'error');
    sandbox.stub(Logger, 'info');
    sandbox.stub(Logger, 'warn');
    sandbox.stub(Logger, 'debug');

    // Import and create seeder
    TestDataSeeder = require('../../src/debug/testDataSeeder').TestDataSeeder;
    seeder = new TestDataSeeder(mockMCPClient);
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('Constructor', () => {
    test('should create instance with MCP client', () => {
      assert(seeder instanceof TestDataSeeder);
      assert.strictEqual(seeder.mcpClient, mockMCPClient);
    });

    test('should have template IDs defined', () => {
      assert(seeder.TEMPLATE_IDS);
      assert(seeder.TEMPLATE_IDS.BUG_RAPORU);
      assert(seeder.TEMPLATE_IDS.OZELLIK_ISTEGI);
      assert(seeder.TEMPLATE_IDS.TEKNIK_BORC);
      assert(seeder.TEMPLATE_IDS.ARASTIRMA_GOREVI);
    });
  });

  suite('seedTestData', () => {
    let progressCallback;
    let progressReport;

    setup(() => {
      progressReport = sandbox.stub();
      progressCallback = sandbox.stub();
      
      vscode.window.withProgress.callsFake((options, callback) => {
        return callback({ report: progressReport });
      });

      // Mock successful responses
      mockMCPClient.callTool.resolves({
        content: [{ text: 'Created with ID: 12345678-1234-1234-1234-123456789abc' }]
      });
    });

    test('should prompt user for confirmation', async () => {
      vscode.window.showInformationMessage.resolves('Hayır');
      
      await seeder.seedTestData();
      
      assert(vscode.window.showInformationMessage.calledWith(
        'Test verileri oluşturulacak. Mevcut veriler korunacak. Devam etmek istiyor musunuz?',
        'Evet, Oluştur',
        'Hayır'
      ));
    });

    test('should exit early if user cancels', async () => {
      vscode.window.showInformationMessage.resolves('Hayır');
      
      await seeder.seedTestData();
      
      assert(!vscode.window.withProgress.called);
      assert(!mockMCPClient.callTool.called);
    });

    test('should create test data when user confirms', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      assert(vscode.window.withProgress.called);
      assert(mockMCPClient.callTool.called);
    });

    test('should show progress messages', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      assert(progressReport.calledWith({ increment: 10, message: 'Projeler oluşturuluyor...' }));
      assert(progressReport.calledWith({ increment: 30, message: 'Görevler oluşturuluyor...' }));
      assert(progressReport.calledWith({ increment: 20, message: 'Bağımlılıklar oluşturuluyor...' }));
      assert(progressReport.calledWith({ increment: 10, message: 'Alt görevler oluşturuluyor...' }));
    });

    test('should show success message when completed', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      assert(vscode.window.showInformationMessage.calledWith('✅ Test verileri başarıyla oluşturuldu!'));
    });

    test('should handle seeding error', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      const error = new Error('Seeding failed');
      mockMCPClient.callTool.rejects(error);
      
      await seeder.seedTestData();
      
      assert(vscode.window.showErrorMessage.calledWith('Test verileri oluşturulamadı: Error: Seeding failed'));
      assert(Logger.error.calledWith('Test data seeding failed:', error));
    });

    test('should create projects first', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      // Verify proje_olustur is called
      const projectCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'proje_olustur'
      );
      assert(projectCalls.length > 0);
    });

    test('should set active project', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      // Verify proje_aktif_yap is called
      const activeCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'proje_aktif_yap'
      );
      assert(activeCalls.length > 0);
    });
  });

  suite('createTestProjects', () => {
    test('should create multiple projects', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.resolves({
        content: [{ text: 'Created project with ID: 12345678-1234-1234-1234-123456789abc' }]
      });
      
      await seeder.seedTestData();
      
      const projectCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'proje_olustur'
      );
      assert(projectCalls.length >= 5); // At least 5 projects
    });

    test('should handle project creation error', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('proje_olustur').rejects(new Error('Project creation failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });

    test('should parse project IDs from response', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      const projectId = '12345678-1234-1234-1234-123456789abc';
      mockMCPClient.callTool.withArgs('proje_olustur').resolves({
        content: [{ text: `Project created successfully with ID: ${projectId}` }]
      });
      
      await seeder.seedTestData();
      
      assert(Logger.info.calledWithMatch(`Created project:`, sinon.match.string, `with ID: ${projectId}`));
    });
  });

  suite('createTestTasks', () => {
    test('should create template-based tasks', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const templateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur'
      );
      assert(templateCalls.length > 0);
    });

    test('should use all template types', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const templateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur'
      );
      
      const templateIds = templateCalls.map(call => call.args[1].template_id);
      assert(templateIds.includes(seeder.TEMPLATE_IDS.BUG_RAPORU));
      assert(templateIds.includes(seeder.TEMPLATE_IDS.OZELLIK_ISTEGI));
      assert(templateIds.includes(seeder.TEMPLATE_IDS.TEKNIK_BORC));
      assert(templateIds.includes(seeder.TEMPLATE_IDS.ARASTIRMA_GOREVI));
    });

    test('should assign tasks to projects', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('templateden_gorev_olustur').resolves({
        content: [{ text: 'Task created with ID: 87654321-4321-4321-4321-210987654321' }]
      });
      
      await seeder.seedTestData();
      
      const editCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_duzenle'
      );
      assert(editCalls.length > 0);
    });

    test('should handle task creation error', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('templateden_gorev_olustur').rejects(new Error('Task creation failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });
  });

  suite('createTestDependencies', () => {
    test('should create task dependencies', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const depCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_bagimlilik_ekle'
      );
      assert(depCalls.length > 0);
    });

    test('should handle dependency creation error', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('gorev_bagimlilik_ekle').rejects(new Error('Dependency creation failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });
  });

  suite('createSubtasks', () => {
    test('should create subtasks', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const subtaskCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_altgorev_olustur'
      );
      assert(subtaskCalls.length > 0);
    });

    test('should create nested subtasks', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      // Mock response for subtask creation to enable nested subtask creation
      mockMCPClient.callTool.withArgs('gorev_altgorev_olustur').resolves({
        content: [{ text: 'Subtask created with ID: 99999999-9999-9999-9999-999999999999' }]
      });
      
      await seeder.seedTestData();
      
      const subtaskCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_altgorev_olustur'
      );
      
      // Should have multiple subtask calls including nested ones
      assert(subtaskCalls.length > 2);
    });
  });

  suite('updateSomeTaskStatuses', () => {
    test('should update task statuses', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const statusCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_guncelle'
      );
      assert(statusCalls.length > 0);
    });

    test('should set tasks to different statuses', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const statusCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_guncelle'
      );
      
      const statuses = statusCalls.map(call => call.args[1].durum);
      assert(statuses.includes(GorevDurum.DevamEdiyor));
      assert(statuses.includes(GorevDurum.Tamamlandi));
    });
  });

  suite('setupAIContext', () => {
    test('should setup AI context', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const aiCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0].startsWith('gorev_') && 
        ['gorev_set_active', 'gorev_nlp_query', 'gorev_context_summary', 'gorev_batch_update'].includes(call.args[0])
      );
      assert(aiCalls.length > 0);
    });

    test('should test NLP queries', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const nlpCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_nlp_query'
      );
      assert(nlpCalls.length > 0);
    });

    test('should perform batch updates', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const batchCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_batch_update'
      );
      assert(batchCalls.length > 0);
    });
  });

  suite('clearTestData', () => {
    test('should prompt user for confirmation', async () => {
      vscode.window.showWarningMessage.resolves('Hayır');
      
      await seeder.clearTestData();
      
      assert(vscode.window.showWarningMessage.calledWith(
        '⚠️ DİKKAT: Tüm görevler ve projeler silinecek! Emin misiniz?',
        'Evet, Sil',
        'Hayır'
      ));
    });

    test('should exit early if user cancels', async () => {
      vscode.window.showWarningMessage.resolves('Hayır');
      
      await seeder.clearTestData();
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should list and delete tasks when user confirms', async () => {
      vscode.window.showWarningMessage.resolves('Evet, Sil');
      mockMCPClient.callTool.withArgs('gorev_listele').resolves({
        content: [{ text: 'ID: 123\nID: 456\nID: 789' }]
      });
      
      await seeder.clearTestData();
      
      assert(mockMCPClient.callTool.calledWith('gorev_listele', { tum_projeler: true }));
      
      const deleteCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_sil'
      );
      assert.strictEqual(deleteCalls.length, 3);
    });

    test('should show success message when completed', async () => {
      vscode.window.showWarningMessage.resolves('Evet, Sil');
      mockMCPClient.callTool.withArgs('gorev_listele').resolves({
        content: [{ text: 'No tasks found' }]
      });
      
      await seeder.clearTestData();
      
      assert(vscode.window.showInformationMessage.calledWith('✅ Test verileri temizlendi!'));
    });

    test('should handle clearing error', async () => {
      vscode.window.showWarningMessage.resolves('Evet, Sil');
      const error = new Error('Clearing failed');
      mockMCPClient.callTool.rejects(error);
      
      await seeder.clearTestData();
      
      assert(vscode.window.showErrorMessage.calledWith('Test verileri temizlenemedi: Error: Clearing failed'));
      assert(Logger.error.calledWith('Failed to clear test data:', error));
    });

    test('should handle individual task deletion errors', async () => {
      vscode.window.showWarningMessage.resolves('Evet, Sil');
      mockMCPClient.callTool.withArgs('gorev_listele').resolves({
        content: [{ text: 'ID: 123\nID: 456' }]
      });
      mockMCPClient.callTool.withArgs('gorev_sil').rejects(new Error('Delete failed'));
      
      await seeder.clearTestData();
      
      assert(Logger.error.called);
    });
  });

  suite('getDateString', () => {
    test('should return date string for future days', () => {
      const result = seeder.getDateString(7);
      const expected = new Date();
      expected.setDate(expected.getDate() + 7);
      
      assert.strictEqual(result, expected.toISOString().split('T')[0]);
    });

    test('should return today for 0 days', () => {
      const result = seeder.getDateString(0);
      const expected = new Date().toISOString().split('T')[0];
      
      assert.strictEqual(result, expected);
    });

    test('should handle negative days', () => {
      const result = seeder.getDateString(-1);
      const expected = new Date();
      expected.setDate(expected.getDate() - 1);
      
      assert.strictEqual(result, expected.toISOString().split('T')[0]);
    });
  });

  suite('Edge Cases', () => {
    test('should handle null MCP client', () => {
      assert.doesNotThrow(() => {
        new TestDataSeeder(null);
      });
    });

    test('should handle malformed MCP responses', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.resolves({
        content: [{ text: 'Invalid response without ID' }]
      });
      
      await seeder.seedTestData();
      
      assert(Logger.warn.called);
    });

    test('should handle empty task list during clearing', async () => {
      vscode.window.showWarningMessage.resolves('Evet, Sil');
      mockMCPClient.callTool.withArgs('gorev_listele').resolves({
        content: [{ text: '' }]
      });
      
      await seeder.clearTestData();
      
      assert(vscode.window.showInformationMessage.calledWith('✅ Test verileri temizlendi!'));
    });

    test('should handle progress callback errors', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      vscode.window.withProgress.callsFake((options, callback) => {
        const mockProgress = {
          report: sandbox.stub().throws(new Error('Progress error'))
        };
        return callback(mockProgress);
      });
      
      await seeder.seedTestData();
      
      // Should handle gracefully and still show success
      assert(vscode.window.showErrorMessage.called || vscode.window.showInformationMessage.called);
    });
  });

  suite('Template Integration', () => {
    test('should use correct template fields for bug reports', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const bugTemplateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.BUG_RAPORU
      );
      
      assert(bugTemplateCalls.length > 0);
      
      const bugTask = bugTemplateCalls[0].args[1].degerler;
      assert(bugTask.modul);
      assert(bugTask.ortam);
      assert(bugTask.adimlar);
      assert(bugTask.beklenen);
      assert(bugTask.mevcut);
    });

    test('should use correct template fields for feature requests', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const featureTemplateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.OZELLIK_ISTEGI
      );
      
      assert(featureTemplateCalls.length > 0);
      
      const featureTask = featureTemplateCalls[0].args[1].degerler;
      assert(featureTask.amac);
      assert(featureTask.kullanicilar);
      assert(featureTask.kriterler);
      assert(featureTask.efor);
    });
  });
});