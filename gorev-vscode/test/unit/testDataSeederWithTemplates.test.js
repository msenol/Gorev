const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const { Logger } = require('../../src/utils/logger');
const { GorevDurum, GorevOncelik } = require('../../src/models/common');

suite('TestDataSeederWithTemplates Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let TestDataSeederWithTemplates;
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
    TestDataSeederWithTemplates = require('../../src/debug/testDataSeederWithTemplates').TestDataSeederWithTemplates;
    seeder = new TestDataSeederWithTemplates(mockMCPClient);
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('Constructor', () => {
    test('should create instance with MCP client', () => {
      assert(seeder instanceof TestDataSeederWithTemplates);
      assert.strictEqual(seeder.mcpClient, mockMCPClient);
    });

    test('should have template IDs defined', () => {
      assert(seeder.TEMPLATE_IDS);
      assert(seeder.TEMPLATE_IDS.BUG_RAPORU);
      assert(seeder.TEMPLATE_IDS.OZELLIK_ISTEGI);
      assert(seeder.TEMPLATE_IDS.TEKNIK_BORC);
      assert(seeder.TEMPLATE_IDS.ARASTIRMA_GOREVI);
    });

    test('should use same template IDs as original seeder', () => {
      const expectedIds = {
        BUG_RAPORU: '4dd56a2a-caf4-472c-8c0f-276bc8a1f880',
        OZELLIK_ISTEGI: '6b083358-9c4d-4f4e-b041-9288c05a1bb7',
        TEKNIK_BORC: '69e2b237-7c2e-4459-9d46-ea6c05aba39a',
        ARASTIRMA_GOREVI: '13f04fe2-b5b6-4fd6-8684-5eca5dc2770d'
      };
      
      assert.deepStrictEqual(seeder.TEMPLATE_IDS, expectedIds);
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

    test('should prompt user for confirmation with template-specific message', async () => {
      vscode.window.showInformationMessage.resolves('Hayır');
      
      await seeder.seedTestData();
      
      assert(vscode.window.showInformationMessage.calledWith(
        'Template-based test verileri oluşturulacak. Mevcut veriler korunacak. Devam etmek istiyor musunuz?',
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

    test('should show template-specific progress messages', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      assert(progressReport.calledWith({ increment: 10, message: 'Projeler oluşturuluyor...' }));
      assert(progressReport.calledWith({ increment: 30, message: 'Template görevleri oluşturuluyor...' }));
      assert(progressReport.calledWith({ increment: 20, message: 'Bağımlılıklar oluşturuluyor...' }));
      assert(progressReport.calledWith({ increment: 10, message: 'Alt görevler oluşturuluyor...' }));
      assert(progressReport.calledWith({ increment: 20, message: 'Görev durumları güncelleniyor...' }));
      assert(progressReport.calledWith({ increment: 10, message: 'AI context oluşturuluyor...' }));
    });

    test('should show template-specific success message', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      assert(vscode.window.showInformationMessage.calledWith('✅ Template-based test verileri başarıyla oluşturuldu!'));
    });

    test('should handle seeding error', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      const error = new Error('Seeding failed');
      mockMCPClient.callTool.rejects(error);
      
      await seeder.seedTestData();
      
      assert(vscode.window.showErrorMessage.calledWith('Test verileri oluşturulamadı: Error: Seeding failed'));
      assert(Logger.error.calledWith('Test data seeding failed:', error));
    });
  });

  suite('createTestProjects', () => {
    test('should create identical projects to original seeder', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const projectCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'proje_olustur'
      );
      
      assert.strictEqual(projectCalls.length, 5);
      
      const projectNames = projectCalls.map(call => call.args[1].isim);
      assert(projectNames.includes('🚀 Yeni Web Sitesi'));
      assert(projectNames.includes('📱 Mobil Uygulama'));
      assert(projectNames.includes('🔧 Backend API'));
      assert(projectNames.includes('📊 Veri Analitiği'));
      assert(projectNames.includes('🔒 Güvenlik Güncellemeleri'));
    });

    test('should set first project as active', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const activeCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'proje_aktif_yap'
      );
      assert(activeCalls.length > 0);
    });

    test('should handle project creation error gracefully', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('proje_olustur').rejects(new Error('Project creation failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });

    test('should parse project IDs correctly', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      const projectId = '12345678-1234-1234-1234-123456789abc';
      mockMCPClient.callTool.withArgs('proje_olustur').resolves({
        content: [{ text: `Project created successfully with ID: ${projectId}` }]
      });
      
      await seeder.seedTestData();
      
      assert(Logger.info.calledWithMatch(`Created project:`, sinon.match.string, `with ID: ${projectId}`));
    });
  });

  suite('createTemplateBasedTasks', () => {
    test('should create bug report tasks with enhanced details', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const bugTemplateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.BUG_RAPORU
      );
      
      assert(bugTemplateCalls.length >= 3);
      
      // Check first bug task has enhanced details
      const firstBugTask = bugTemplateCalls[0].args[1].degerler;
      assert.strictEqual(firstBugTask.baslik, 'Login sayfası 404 hatası veriyor');
      assert(firstBugTask.modul);
      assert(firstBugTask.ortam);
      assert(firstBugTask.adimlar);
      assert(firstBugTask.beklenen);
      assert(firstBugTask.mevcut);
      assert(firstBugTask.cozum);
    });

    test('should create feature request tasks with comprehensive fields', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const featureTemplateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.OZELLIK_ISTEGI
      );
      
      assert(featureTemplateCalls.length >= 4);
      
      // Check feature task has comprehensive fields
      const featureTask = featureTemplateCalls[0].args[1].degerler;
      assert(featureTask.amac);
      assert(featureTask.kullanicilar);
      assert(featureTask.kriterler);
      assert(featureTask.ui_ux);
      assert(featureTask.ilgili);
      assert(featureTask.efor);
    });

    test('should create technical debt tasks with detailed analysis', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const techDebtCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.TEKNIK_BORC
      );
      
      assert(techDebtCalls.length >= 3);
      
      // Check tech debt task has detailed analysis
      const techDebtTask = techDebtCalls[0].args[1].degerler;
      assert(techDebtTask.alan);
      assert(techDebtTask.dosyalar);
      assert(techDebtTask.neden);
      assert(techDebtTask.analiz);
      assert(techDebtTask.cozum);
      assert(techDebtTask.riskler);
      assert(techDebtTask.iyilestirmeler);
      assert(techDebtTask.sure);
    });

    test('should create research tasks with comprehensive investigation fields', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const researchCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.ARASTIRMA_GOREVI
      );
      
      assert(researchCalls.length >= 3);
      
      // Check research task has comprehensive fields
      const researchTask = researchCalls[0].args[1].degerler;
      assert(researchTask.konu);
      assert(researchTask.amac);
      assert(researchTask.sorular);
      assert(researchTask.kaynaklar);
      assert(researchTask.alternatifler);
      assert(researchTask.kriterler);
    });

    test('should create direct tasks without templates', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const directTaskCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_olustur'
      );
      
      assert(directTaskCalls.length >= 3);
      
      // Check direct task structure
      const directTask = directTaskCalls[0].args[1];
      assert(directTask.baslik);
      assert(directTask.aciklama);
      assert(directTask.oncelik);
      assert(directTask.etiketler);
    });

    test('should assign tasks to appropriate projects', async () => {
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

    test('should handle task creation errors', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('templateden_gorev_olustur').rejects(new Error('Task creation failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });

    test('should include due dates for time-sensitive tasks', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const researchCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.ARASTIRMA_GOREVI
      );
      
      // Some research tasks should have due dates
      const taskWithDueDate = researchCalls.find(
        call => call.args[1].degerler.son_tarih
      );
      assert(taskWithDueDate);
    });
  });

  suite('createTestDependencies', () => {
    test('should create logical task dependencies', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const depCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_bagimlilik_ekle'
      );
      
      assert(depCalls.length >= 4);
      
      // Check dependency types
      const depTypes = depCalls.map(call => call.args[1].baglanti_tipi);
      assert(depTypes.includes('blocks'));
      assert(depTypes.includes('depends_on'));
    });

    test('should handle dependency creation errors', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('gorev_bagimlilik_ekle').rejects(new Error('Dependency failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });
  });

  suite('createSubtasks', () => {
    test('should create relevant subtasks for main feature tasks', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const subtaskCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_altgorev_olustur'
      );
      
      assert(subtaskCalls.length >= 2);
      
      // Check subtask structure
      const subtask = subtaskCalls[0].args[1];
      assert(subtask.parent_id);
      assert(subtask.baslik);
      assert(subtask.aciklama);
      assert(subtask.oncelik);
      assert(subtask.etiketler);
    });

    test('should create subtasks with appropriate priorities', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const subtaskCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_altgorev_olustur'
      );
      
      const priorities = subtaskCalls.map(call => call.args[1].oncelik);
      assert(priorities.includes(GorevOncelik.Yuksek));
      assert(priorities.includes(GorevOncelik.Orta));
    });
  });

  suite('updateSomeTaskStatuses', () => {
    test('should update tasks to in-progress status', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const inProgressCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_guncelle' &&
        call.args[1].durum === GorevDurum.DevamEdiyor
      );
      
      assert(inProgressCalls.length >= 3);
    });

    test('should complete some tasks', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const completedCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_guncelle' &&
        call.args[1].durum === GorevDurum.Tamamlandi
      );
      
      assert(completedCalls.length >= 2);
    });

    test('should handle status update errors', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('gorev_guncelle').rejects(new Error('Status update failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });
  });

  suite('setupAIContext', () => {
    test('should set active task for AI context', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const activeCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_set_active'
      );
      
      assert(activeCalls.length >= 1);
    });

    test('should generate context summary', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const summaryCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_context_summary'
      );
      
      assert(summaryCalls.length >= 1);
    });

    test('should handle AI context setup errors', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('gorev_set_active').rejects(new Error('AI context failed'));
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });
  });

  suite('getDateString', () => {
    test('should generate correct date strings', () => {
      const result = seeder.getDateString(7);
      const expected = new Date();
      expected.setDate(expected.getDate() + 7);
      
      assert.strictEqual(result, expected.toISOString().split('T')[0]);
    });

    test('should handle today (0 days)', () => {
      const result = seeder.getDateString(0);
      const expected = new Date().toISOString().split('T')[0];
      
      assert.strictEqual(result, expected);
    });

    test('should handle future dates', () => {
      const result = seeder.getDateString(14);
      const expected = new Date();
      expected.setDate(expected.getDate() + 14);
      
      assert.strictEqual(result, expected.toISOString().split('T')[0]);
    });
  });

  suite('Enhanced Template Features', () => {
    test('should create security-focused bug reports', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const bugTemplateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.BUG_RAPORU
      );
      
      // Find SSL certificate bug
      const sslBug = bugTemplateCalls.find(
        call => call.args[1].degerler.baslik.includes('SSL sertifikası')
      );
      
      assert(sslBug);
      assert(sslBug.args[1].degerler.etiketler.includes('security'));
    });

    test('should create performance-focused technical debt', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const techDebtCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.TEKNIK_BORC
      );
      
      // Find Redis cache task
      const redisTask = techDebtCalls.find(
        call => call.args[1].degerler.baslik.includes('Redis cache')
      );
      
      assert(redisTask);
      assert(redisTask.args[1].degerler.iyilestirmeler.includes('performance'));
    });

    test('should create modern framework research tasks', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const researchCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.ARASTIRMA_GOREVI
      );
      
      // Find Next.js research task
      const nextjsTask = researchCalls.find(
        call => call.args[1].degerler.konu.includes('Next.js')
      );
      
      assert(nextjsTask);
      assert(nextjsTask.args[1].degerler.sorular.includes('Performance improvements'));
    });

    test('should include risk assessment in technical debt', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const techDebtCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur' &&
        call.args[1].template_id === seeder.TEMPLATE_IDS.TEKNIK_BORC
      );
      
      // All tech debt tasks should have risk assessments
      techDebtCalls.forEach(call => {
        assert(call.args[1].degerler.riskler);
      });
    });
  });

  suite('Edge Cases', () => {
    test('should handle null MCP client', () => {
      assert.doesNotThrow(() => {
        new TestDataSeederWithTemplates(null);
      });
    });

    test('should handle malformed template responses', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('templateden_gorev_olustur').resolves({
        content: [{ text: 'Invalid response without ID' }]
      });
      
      await seeder.seedTestData();
      
      assert(Logger.error.called);
    });

    test('should handle empty project list', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      mockMCPClient.callTool.withArgs('proje_olustur').resolves({
        content: [{ text: 'Failed to parse ID' }]
      });
      
      await seeder.seedTestData();
      
      // Should continue with empty project list
      assert(mockMCPClient.callTool.called);
    });

    test('should handle progress reporting errors', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      vscode.window.withProgress.callsFake((options, callback) => {
        const mockProgress = {
          report: sandbox.stub().throws(new Error('Progress error'))
        };
        return callback(mockProgress);
      });
      
      await seeder.seedTestData();
      
      // Should complete despite progress errors
      assert(vscode.window.showErrorMessage.called || mockMCPClient.callTool.called);
    });
  });

  suite('Data Quality', () => {
    test('should create realistic task data', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const templateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur'
      );
      
      // All template tasks should have meaningful titles and descriptions
      templateCalls.forEach(call => {
        const task = call.args[1].degerler;
        assert(task.baslik);
        assert(task.baslik.length > 10);
        assert(task.aciklama);
        assert(task.aciklama.length > 20);
      });
    });

    test('should distribute tasks across different projects', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const editCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'gorev_duzenle'
      );
      
      // Should have tasks assigned to different projects
      const projectIds = editCalls.map(call => call.args[1].proje_id);
      const uniqueProjects = [...new Set(projectIds)];
      assert(uniqueProjects.length > 1);
    });

    test('should create balanced workload across priorities', async () => {
      vscode.window.showInformationMessage.resolves('Evet, Oluştur');
      
      await seeder.seedTestData();
      
      const templateCalls = mockMCPClient.callTool.getCalls().filter(
        call => call.args[0] === 'templateden_gorev_olustur'
      );
      
      const priorities = templateCalls.map(call => call.args[1].degerler.oncelik);
      
      // Should have mix of priorities
      assert(priorities.includes('yuksek'));
      assert(priorities.includes('orta'));
      assert(priorities.includes('dusuk'));
    });
  });
});