const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const { GorevDurum, GorevOncelik } = require('../../src/models/gorev');
const { Logger } = require('../../src/utils/logger');

suite('InlineEditProvider Test Suite', () => {
  let sandbox;
  let mockMCPClient;
  let InlineEditProvider;
  let provider;

  setup(() => {
    sandbox = sinon.createSandbox();
    
    // Mock VS Code API
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    sandbox.stub(vscode.window, 'showWarningMessage');
    sandbox.stub(vscode.window, 'showInputBox');
    sandbox.stub(vscode.window, 'showQuickPick');
    sandbox.stub(vscode.commands, 'executeCommand');

    // Mock MCP Client
    mockMCPClient = {
      callTool: sandbox.stub().resolves({ content: [{ text: 'Success' }] })
    };

    // Mock Logger
    sandbox.stub(Logger, 'error');
    sandbox.stub(Logger, 'info');

    // Import and create provider
    InlineEditProvider = require('../../src/providers/inlineEditProvider').InlineEditProvider;
    provider = new InlineEditProvider(mockMCPClient);
  });

  teardown(() => {
    sandbox.restore();
  });

  suite('Constructor', () => {
    test('should create instance with MCP client', () => {
      assert(provider instanceof InlineEditProvider);
      assert.strictEqual(provider.mcpClient, mockMCPClient);
    });

    test('should initialize editing state', () => {
      assert.strictEqual(provider.editingItem, null);
      assert.strictEqual(provider.originalLabel, null);
    });
  });

  suite('startEdit', () => {
    const mockTask = {
      id: '123',
      baslik: 'Test Task',
      aciklama: 'Test description'
    };

    test('should start edit with valid item', async () => {
      const item = { task: mockTask };
      vscode.window.showInputBox.resolves('Updated Task Title');
      
      await provider.startEdit(item);
      
      assert(vscode.window.showInputBox.calledWith({
        prompt: 'Görev başlığını düzenle',
        value: mockTask.baslik,
        validateInput: sinon.match.func
      }));
      
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', {
        id: mockTask.id,
        baslik: 'Updated Task Title'
      }));
    });

    test('should not edit if no item provided', async () => {
      await provider.startEdit(null);
      
      assert(!vscode.window.showInputBox.called);
      assert(!mockMCPClient.callTool.called);
    });

    test('should not edit if item has no task', async () => {
      await provider.startEdit({});
      
      assert(!vscode.window.showInputBox.called);
      assert(!mockMCPClient.callTool.called);
    });

    test('should not save if user cancels', async () => {
      const item = { task: mockTask };
      vscode.window.showInputBox.resolves(undefined);
      
      await provider.startEdit(item);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should not save if title unchanged', async () => {
      const item = { task: mockTask };
      vscode.window.showInputBox.resolves(mockTask.baslik);
      
      await provider.startEdit(item);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should validate input - empty title', async () => {
      const item = { task: mockTask };
      
      await provider.startEdit(item);
      
      const validator = vscode.window.showInputBox.getCall(0).args[0].validateInput;
      const result = validator('');
      
      assert.strictEqual(result, 'Görev başlığı boş olamaz');
    });

    test('should validate input - whitespace only', async () => {
      const item = { task: mockTask };
      
      await provider.startEdit(item);
      
      const validator = vscode.window.showInputBox.getCall(0).args[0].validateInput;
      const result = validator('   ');
      
      assert.strictEqual(result, 'Görev başlığı boş olamaz');
    });

    test('should validate input - too long', async () => {
      const item = { task: mockTask };
      
      await provider.startEdit(item);
      
      const validator = vscode.window.showInputBox.getCall(0).args[0].validateInput;
      const longTitle = 'a'.repeat(201);
      const result = validator(longTitle);
      
      assert.strictEqual(result, 'Görev başlığı 200 karakterden uzun olamaz');
    });

    test('should validate input - valid title', async () => {
      const item = { task: mockTask };
      
      await provider.startEdit(item);
      
      const validator = vscode.window.showInputBox.getCall(0).args[0].validateInput;
      const result = validator('Valid Title');
      
      assert.strictEqual(result, null);
    });

    test('should show success message after save', async () => {
      const item = { task: mockTask };
      vscode.window.showInputBox.resolves('Updated Title');
      
      await provider.startEdit(item);
      
      assert(vscode.window.showInformationMessage.calledWith('Görev başlığı güncellendi'));
      assert(Logger.info.calledWith(`Task ${mockTask.id} title updated to: Updated Title`));
    });

    test('should handle save error', async () => {
      const item = { task: mockTask };
      vscode.window.showInputBox.resolves('Updated Title');
      const error = new Error('Save failed');
      mockMCPClient.callTool.rejects(error);
      
      await provider.startEdit(item);
      
      assert(vscode.window.showErrorMessage.calledWith('Güncelleme başarısız: Error: Save failed'));
      assert(Logger.error.calledWith('Failed to update task title:', error));
    });

    test('should track editing state during edit', async () => {
      const item = { task: mockTask };
      let editingState;
      
      vscode.window.showInputBox.callsFake(() => {
        editingState = provider.isEditing();
        return Promise.resolve('Updated Title');
      });
      
      await provider.startEdit(item);
      
      assert.strictEqual(editingState, true);
      assert.strictEqual(provider.isEditing(), false);
    });
  });

  suite('quickStatusChange', () => {
    const mockTask = {
      id: '123',
      baslik: 'Test Task',
      durum: GorevDurum.Beklemede
    };

    test('should show status options', async () => {
      await provider.quickStatusChange(mockTask);
      
      assert(vscode.window.showQuickPick.called);
      
      const items = vscode.window.showQuickPick.getCall(0).args[0];
      assert.strictEqual(items.length, 3);
      assert(items.some(item => item.value === GorevDurum.Beklemede));
      assert(items.some(item => item.value === GorevDurum.DevamEdiyor));
      assert(items.some(item => item.value === GorevDurum.Tamamlandi));
    });

    test('should mark current status in options', async () => {
      await provider.quickStatusChange(mockTask);
      
      const items = vscode.window.showQuickPick.getCall(0).args[0];
      const currentItem = items.find(item => item.value === mockTask.durum);
      
      assert.strictEqual(currentItem.description, 'Mevcut durum');
    });

    test('should update status when different status selected', async () => {
      const selectedItem = { value: GorevDurum.DevamEdiyor };
      vscode.window.showQuickPick.resolves(selectedItem);
      
      await provider.quickStatusChange(mockTask);
      
      assert(mockMCPClient.callTool.calledWith('gorev_guncelle', {
        id: mockTask.id,
        durum: GorevDurum.DevamEdiyor
      }));
      
      assert(vscode.window.showInformationMessage.calledWith('Görev durumu güncellendi'));
      assert(vscode.commands.executeCommand.calledWith('gorev.refreshTasks'));
    });

    test('should not update when same status selected', async () => {
      const selectedItem = { value: mockTask.durum };
      vscode.window.showQuickPick.resolves(selectedItem);
      
      await provider.quickStatusChange(mockTask);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should not update when user cancels', async () => {
      vscode.window.showQuickPick.resolves(undefined);
      
      await provider.quickStatusChange(mockTask);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should handle update error', async () => {
      const selectedItem = { value: GorevDurum.DevamEdiyor };
      vscode.window.showQuickPick.resolves(selectedItem);
      const error = new Error('Update failed');
      mockMCPClient.callTool.rejects(error);
      
      await provider.quickStatusChange(mockTask);
      
      assert(vscode.window.showErrorMessage.calledWith('Durum güncellemesi başarısız: Error: Update failed'));
      assert(Logger.error.calledWith('[QuickStatusChange] Failed to update task status:', error));
    });

    test('should log status change details', async () => {
      const selectedItem = { value: GorevDurum.DevamEdiyor };
      vscode.window.showQuickPick.resolves(selectedItem);
      
      await provider.quickStatusChange(mockTask);
      
      assert(Logger.info.calledWithMatch(`[QuickStatusChange] Updating task ${mockTask.id} from ${mockTask.durum} to ${selectedItem.value}`));
      assert(Logger.info.calledWithMatch(`[QuickStatusChange] Task ${mockTask.id} status updated to: ${selectedItem.value}`));
    });
  });

  suite('quickPriorityChange', () => {
    const mockTask = {
      id: '123',
      baslik: 'Test Task',
      oncelik: GorevOncelik.Orta
    };

    test('should show priority options', async () => {
      await provider.quickPriorityChange(mockTask);
      
      assert(vscode.window.showQuickPick.called);
      
      const items = vscode.window.showQuickPick.getCall(0).args[0];
      assert.strictEqual(items.length, 3);
      assert(items.some(item => item.value === GorevOncelik.Yuksek));
      assert(items.some(item => item.value === GorevOncelik.Orta));
      assert(items.some(item => item.value === GorevOncelik.Dusuk));
    });

    test('should mark current priority in options', async () => {
      await provider.quickPriorityChange(mockTask);
      
      const items = vscode.window.showQuickPick.getCall(0).args[0];
      const currentItem = items.find(item => item.value === mockTask.oncelik);
      
      assert.strictEqual(currentItem.description, 'Mevcut öncelik');
    });

    test('should update priority when different priority selected', async () => {
      const selectedItem = { value: GorevOncelik.Yuksek };
      vscode.window.showQuickPick.resolves(selectedItem);
      
      await provider.quickPriorityChange(mockTask);
      
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', {
        id: mockTask.id,
        oncelik: GorevOncelik.Yuksek
      }));
      
      assert(vscode.window.showInformationMessage.calledWith('Görev önceliği güncellendi'));
    });

    test('should not update when same priority selected', async () => {
      const selectedItem = { value: mockTask.oncelik };
      vscode.window.showQuickPick.resolves(selectedItem);
      
      await provider.quickPriorityChange(mockTask);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should handle update error', async () => {
      const selectedItem = { value: GorevOncelik.Yuksek };
      vscode.window.showQuickPick.resolves(selectedItem);
      const error = new Error('Priority update failed');
      mockMCPClient.callTool.rejects(error);
      
      await provider.quickPriorityChange(mockTask);
      
      assert(vscode.window.showErrorMessage.calledWith('Öncelik güncellemesi başarısız: Error: Priority update failed'));
      assert(Logger.error.calledWith('Failed to update task priority:', error));
    });
  });

  suite('quickDateChange', () => {
    const mockTask = {
      id: '123',
      baslik: 'Test Task',
      son_tarih: '2024-12-31'
    };

    test('should show date input with current value', async () => {
      await provider.quickDateChange(mockTask);
      
      assert(vscode.window.showInputBox.calledWith({
        prompt: 'Son tarihi girin (YYYY-MM-DD)',
        value: mockTask.son_tarih,
        placeHolder: '2024-12-31',
        validateInput: sinon.match.func
      }));
    });

    test('should validate date format', async () => {
      await provider.quickDateChange(mockTask);
      
      const validator = vscode.window.showInputBox.getCall(0).args[0].validateInput;
      
      assert.strictEqual(validator('2024-12-31'), null);
      assert.strictEqual(validator('invalid-date'), 'Geçersiz tarih formatı. YYYY-MM-DD kullanın');
      assert.strictEqual(validator('2024-13-01'), 'Geçersiz tarih');
      assert.strictEqual(validator(''), null); // Empty is allowed
    });

    test('should update date when valid date provided', async () => {
      vscode.window.showInputBox.resolves('2025-01-15');
      
      await provider.quickDateChange(mockTask);
      
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', {
        id: mockTask.id,
        son_tarih: '2025-01-15'
      }));
      
      assert(vscode.window.showInformationMessage.calledWith('Son tarih güncellendi'));
    });

    test('should remove date when empty string provided', async () => {
      vscode.window.showInputBox.resolves('');
      
      await provider.quickDateChange(mockTask);
      
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', {
        id: mockTask.id,
        son_tarih: null
      }));
      
      assert(vscode.window.showInformationMessage.calledWith('Son tarih kaldırıldı'));
    });

    test('should not update when date unchanged', async () => {
      vscode.window.showInputBox.resolves(mockTask.son_tarih);
      
      await provider.quickDateChange(mockTask);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should not update when user cancels', async () => {
      vscode.window.showInputBox.resolves(undefined);
      
      await provider.quickDateChange(mockTask);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should handle task with no current date', async () => {
      const taskWithoutDate = { ...mockTask, son_tarih: null };
      
      await provider.quickDateChange(taskWithoutDate);
      
      assert(vscode.window.showInputBox.calledWith(sinon.match({
        value: ''
      })));
    });

    test('should handle update error', async () => {
      vscode.window.showInputBox.resolves('2025-01-15');
      const error = new Error('Date update failed');
      mockMCPClient.callTool.rejects(error);
      
      await provider.quickDateChange(mockTask);
      
      assert(vscode.window.showErrorMessage.calledWith('Tarih güncellemesi başarısız: Error: Date update failed'));
      assert(Logger.error.calledWith('Failed to update task due date:', error));
    });
  });

  suite('showDetailedEdit', () => {
    const mockTask = {
      id: '123',
      baslik: 'Test Task',
      aciklama: 'Test description'
    };

    test('should show edit options', async () => {
      await provider.showDetailedEdit(mockTask);
      
      assert(vscode.window.showQuickPick.called);
      
      const items = vscode.window.showQuickPick.getCall(0).args[0];
      assert.strictEqual(items.length, 6);
      assert(items.some(item => item.action === 'title'));
      assert(items.some(item => item.action === 'description'));
      assert(items.some(item => item.action === 'status'));
      assert(items.some(item => item.action === 'priority'));
      assert(items.some(item => item.action === 'dueDate'));
      assert(items.some(item => item.action === 'tags'));
    });

    test('should call appropriate method for title edit', async () => {
      const titleOption = { action: 'title' };
      vscode.window.showQuickPick.resolves(titleOption);
      
      // Spy on startEdit method
      const startEditSpy = sandbox.spy(provider, 'startEdit');
      
      await provider.showDetailedEdit(mockTask);
      
      assert(startEditSpy.calledWith({ task: mockTask }));
    });

    test('should call appropriate method for status change', async () => {
      const statusOption = { action: 'status' };
      vscode.window.showQuickPick.resolves(statusOption);
      
      // Spy on quickStatusChange method
      const statusChangeSpy = sandbox.spy(provider, 'quickStatusChange');
      
      await provider.showDetailedEdit(mockTask);
      
      assert(statusChangeSpy.calledWith(mockTask));
    });

    test('should call appropriate method for priority change', async () => {
      const priorityOption = { action: 'priority' };
      vscode.window.showQuickPick.resolves(priorityOption);
      
      // Spy on quickPriorityChange method
      const priorityChangeSpy = sandbox.spy(provider, 'quickPriorityChange');
      
      await provider.showDetailedEdit(mockTask);
      
      assert(priorityChangeSpy.calledWith(mockTask));
    });

    test('should call appropriate method for date change', async () => {
      const dateOption = { action: 'dueDate' };
      vscode.window.showQuickPick.resolves(dateOption);
      
      // Spy on quickDateChange method
      const dateChangeSpy = sandbox.spy(provider, 'quickDateChange');
      
      await provider.showDetailedEdit(mockTask);
      
      assert(dateChangeSpy.calledWith(mockTask));
    });

    test('should handle description edit', async () => {
      const descOption = { action: 'description' };
      vscode.window.showQuickPick.resolves(descOption);
      vscode.window.showInputBox.resolves('Updated description');
      
      await provider.showDetailedEdit(mockTask);
      
      assert(vscode.window.showInputBox.calledWith({
        prompt: 'Görev açıklamasını düzenle',
        value: mockTask.aciklama,
        placeHolder: 'Görev hakkında detaylı açıklama...'
      }));
      
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', {
        id: mockTask.id,
        aciklama: 'Updated description'
      }));
    });

    test('should handle tags edit', async () => {
      const tagsOption = { action: 'tags' };
      vscode.window.showQuickPick.resolves(tagsOption);
      
      const mockTaskWithTags = { ...mockTask, etiketler: ['bug', 'frontend'] };
      
      await provider.showDetailedEdit(mockTaskWithTags);
      
      assert(vscode.window.showInputBox.calledWith(sinon.match({
        prompt: 'Etiketleri düzenle (virgülle ayırın)',
        value: 'bug, frontend'
      })));
    });

    test('should validate tags input', async () => {
      const tagsOption = { action: 'tags' };
      vscode.window.showQuickPick.resolves(tagsOption);
      
      await provider.showDetailedEdit(mockTask);
      
      const validator = vscode.window.showInputBox.getCall(0).args[0].validateInput;
      
      assert.strictEqual(validator('valid-tag'), null);
      assert.strictEqual(validator('invalid tag with spaces'), 'Etiketler sadece harf, rakam, tire ve alt çizgi içerebilir');
      assert.strictEqual(validator('a'.repeat(51)), 'Etiketler 50 karakterden uzun olamaz');
    });

    test('should not execute action when user cancels', async () => {
      vscode.window.showQuickPick.resolves(undefined);
      
      await provider.showDetailedEdit(mockTask);
      
      assert(!vscode.window.showInputBox.called);
      assert(!mockMCPClient.callTool.called);
    });
  });

  suite('isEditing and cancelEdit', () => {
    test('should return false when not editing', () => {
      assert.strictEqual(provider.isEditing(), false);
    });

    test('should return true when editing', async () => {
      const mockTask = { id: '123', baslik: 'Test Task' };
      const item = { task: mockTask };
      
      // Mock input box to check editing state during call
      let editingState;
      vscode.window.showInputBox.callsFake(() => {
        editingState = provider.isEditing();
        return Promise.resolve('Updated Title');
      });
      
      await provider.startEdit(item);
      
      assert.strictEqual(editingState, true);
    });

    test('should cancel edit and reset state', () => {
      // Manually set editing state
      provider.editingItem = { task: { id: '123' } };
      provider.originalLabel = 'Original Title';
      
      provider.cancelEdit();
      
      assert.strictEqual(provider.editingItem, null);
      assert.strictEqual(provider.originalLabel, null);
      assert.strictEqual(provider.isEditing(), false);
    });
  });

  suite('editDescription', () => {
    const mockTask = {
      id: '123',
      baslik: 'Test Task',
      aciklama: 'Original description'
    };

    test('should update description when changed', async () => {
      vscode.window.showInputBox.resolves('Updated description');
      
      await provider.editDescription(mockTask);
      
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', {
        id: mockTask.id,
        aciklama: 'Updated description'
      }));
      
      assert(vscode.window.showInformationMessage.calledWith('Görev açıklaması güncellendi'));
    });

    test('should not update when description unchanged', async () => {
      vscode.window.showInputBox.resolves(mockTask.aciklama);
      
      await provider.editDescription(mockTask);
      
      assert(!mockMCPClient.callTool.called);
    });

    test('should handle task with no description', async () => {
      const taskWithoutDesc = { ...mockTask, aciklama: null };
      vscode.window.showInputBox.resolves('New description');
      
      await provider.editDescription(taskWithoutDesc);
      
      assert(vscode.window.showInputBox.calledWith(sinon.match({
        value: ''
      })));
    });

    test('should handle description update error', async () => {
      vscode.window.showInputBox.resolves('Updated description');
      const error = new Error('Description update failed');
      mockMCPClient.callTool.rejects(error);
      
      await provider.editDescription(mockTask);
      
      assert(vscode.window.showErrorMessage.calledWith('Açıklama güncellemesi başarısız: Error: Description update failed'));
      assert(Logger.error.calledWith('Failed to update task description:', error));
    });
  });

  suite('editTags', () => {
    test('should show tags edit placeholder message', async () => {
      const mockTask = { id: '123', etiketler: ['bug', 'frontend'] };
      vscode.window.showInputBox.resolves('bug, frontend, urgent');
      
      await provider.editTags(mockTask);
      
      assert(vscode.window.showInformationMessage.calledWith(
        'Etiket güncelleme özelliği yakında eklenecek'
      ));
    });

    test('should handle task with no tags', async () => {
      const mockTask = { id: '123', etiketler: null };
      
      await provider.editTags(mockTask);
      
      assert(vscode.window.showInputBox.calledWith(sinon.match({
        value: ''
      })));
    });
  });

  suite('Edge Cases', () => {
    test('should handle null MCP client', () => {
      const providerWithNullClient = new InlineEditProvider(null);
      
      assert.strictEqual(providerWithNullClient.mcpClient, null);
    });

    test('should handle malformed task objects', async () => {
      const malformedTask = { id: null, baslik: null };
      
      assert.doesNotThrow(async () => {
        await provider.quickStatusChange(malformedTask);
      });
    });

    test('should handle MCP client call failures gracefully', async () => {
      const mockTask = { id: '123', baslik: 'Test Task', durum: GorevDurum.Beklemede };
      const selectedItem = { value: GorevDurum.DevamEdiyor };
      vscode.window.showQuickPick.resolves(selectedItem);
      
      mockMCPClient.callTool.rejects(new Error('Network error'));
      
      await provider.quickStatusChange(mockTask);
      
      assert(Logger.error.called);
      assert(vscode.window.showErrorMessage.called);
    });

    test('should handle missing task properties', async () => {
      const incompleteTask = { id: '123' };
      
      assert.doesNotThrow(async () => {
        await provider.quickDateChange(incompleteTask);
      });
    });
  });

  suite('Integration Tests', () => {
    test('should complete full edit workflow', async () => {
      const mockTask = {
        id: '123',
        baslik: 'Original Task',
        aciklama: 'Original description',
        durum: GorevDurum.Beklemede,
        oncelik: GorevOncelik.Orta,
        son_tarih: '2024-12-31'
      };

      // Test title edit
      const item = { task: mockTask };
      vscode.window.showInputBox.resolves('Updated Task');
      await provider.startEdit(item);
      
      // Test status change
      vscode.window.showQuickPick.resolves({ value: GorevDurum.DevamEdiyor });
      await provider.quickStatusChange(mockTask);
      
      // Test priority change
      vscode.window.showQuickPick.resolves({ value: GorevOncelik.Yuksek });
      await provider.quickPriorityChange(mockTask);
      
      // Test date change
      vscode.window.showInputBox.resolves('2025-01-15');
      await provider.quickDateChange(mockTask);
      
      // Verify all updates were called
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', sinon.match({ baslik: 'Updated Task' })));
      assert(mockMCPClient.callTool.calledWith('gorev_guncelle', sinon.match({ durum: GorevDurum.DevamEdiyor })));
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', sinon.match({ oncelik: GorevOncelik.Yuksek })));
      assert(mockMCPClient.callTool.calledWith('gorev_duzenle', sinon.match({ son_tarih: '2025-01-15' })));
    });
  });
});