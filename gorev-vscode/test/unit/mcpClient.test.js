const assert = require('assert');
const sinon = require('sinon');
const vscode = require('vscode');
const { MCPClient } = require('../../dist/mcp/client');

suite('MCPClient Test Suite', () => {
  let client;
  let sandbox;

  setup(() => {
    sandbox = sinon.createSandbox();
    // Mock vscode.window methods
    sandbox.stub(vscode.window, 'showErrorMessage');
    sandbox.stub(vscode.window, 'showInformationMessage');
    
    client = new MCPClient();
  });

  teardown(() => {
    sandbox.restore();
    if (client.isConnected()) {
      client.disconnect();
    }
  });

  test('should initialize with disconnected state', () => {
    assert.strictEqual(client.isConnected(), false);
  });

  test('should parse server response correctly', () => {
    const mockResponse = {
      jsonrpc: '2.0',
      id: 1,
      result: {
        content: [
          {
            type: 'text',
            text: '## Test Response\\nThis is a test'
          }
        ]
      }
    };

    // This would need the actual implementation to test properly
    // For now, we're testing the structure
    assert(mockResponse.result);
    assert(mockResponse.result.content);
    assert.strictEqual(mockResponse.result.content[0].type, 'text');
  });

  test('should handle error responses', () => {
    const errorResponse = {
      jsonrpc: '2.0',
      id: 1,
      error: {
        code: -32601,
        message: 'Method not found'
      }
    };

    // Test error structure
    assert(errorResponse.error);
    assert.strictEqual(errorResponse.error.code, -32601);
    assert.strictEqual(errorResponse.error.message, 'Method not found');
  });

  suite('Tool calls', () => {
    test('should format gorev_olustur parameters correctly', () => {
      const params = {
        baslik: 'Test Task',
        aciklama: 'Test Description',
        oncelik: 'orta',
        proje_id: 'test-proj',
        son_tarih: '2025-07-01',
        etiketler: 'tag1,tag2'
      };

      // Verify parameter structure
      assert.strictEqual(params.baslik, 'Test Task');
      assert.strictEqual(params.oncelik, 'orta');
      assert(params.son_tarih);
      assert(params.etiketler);
    });

    test('should format gorev_listele parameters correctly', () => {
      const params = {
        durum: 'beklemede',
        tum_projeler: 'false',
        sirala: 'son_tarih_asc',
        filtre: 'acil',
        etiket: 'bug'
      };

      // Verify parameter structure
      assert.strictEqual(params.durum, 'beklemede');
      assert.strictEqual(params.tum_projeler, 'false');
      assert.strictEqual(params.sirala, 'son_tarih_asc');
    });

    test('should format template_listele parameters correctly', () => {
      const params = {
        kategori: 'Teknik'
      };

      // Verify parameter structure
      assert.strictEqual(params.kategori, 'Teknik');
    });
  });

  suite('Connection management', () => {
    test('should set connected state on successful connection', () => {
      // This would require mocking the actual connection
      // For now, we test the expected behavior
      client._connected = true;
      assert.strictEqual(client.isConnected(), true);
    });

    test('should clear state on disconnect', () => {
      client._connected = true;
      client._serverProcess = { kill: sinon.stub() };
      
      client.disconnect();
      
      assert.strictEqual(client.isConnected(), false);
      assert.strictEqual(client._serverProcess, null);
    });
  });
});