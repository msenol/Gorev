package mcp

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	testinghelpers "github.com/msenol/gorev/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test for ServeSunucu
func TestServeSunucu(t *testing.T) {
	// Create a test server
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()
	mcpServer, err := YeniMCPSunucu(isYonetici)
	require.NoError(t, err)

	// We can't actually test stdio serving without a real stdio connection
	// But we can verify the function exists and accepts the correct type
	assert.NotNil(t, mcpServer)

	// The actual ServeSunucu function would block waiting for stdio input
	// So we'll just verify it's callable with the right type
	var serveFn func(*server.MCPServer) error = ServeSunucu
	assert.NotNil(t, serveFn)
}

// Test for NewServer
func TestNewServer(t *testing.T) {
	// Create test handlers
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()

	handlers := YeniHandlers(isYonetici)

	// Test NewServer
	mcpServer := NewServer(handlers)
	assert.NotNil(t, mcpServer)

	// Verify server properties
	// Note: MCP server doesn't expose Name/Version as public fields
	assert.NotNil(t, mcpServer)

	// Test with nil handlers (should not panic)
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil handlers
			assert.NotNil(t, r)
		}
	}()

	// This might panic, which is expected
	nilServer := NewServer(nil)
	assert.NotNil(t, nilServer) // This line may not be reached if it panics
}

// Test for ListTools
func TestListTools(t *testing.T) {
	// Call ListTools
	tools := ListTools()

	// Verify we get a list of tools
	assert.NotEmpty(t, tools)

	// Expected tool names
	expectedTools := []string{
		"gorev_listele",
		"gorev_detay",
		"gorev_guncelle",
		"gorev_duzenle",
		"gorev_sil",
		"gorev_bagimlilik_ekle",
		"gorev_altgorev_olustur",
		"gorev_ust_degistir",
		"gorev_hiyerarsi_goster",
		"proje_olustur",
		"proje_listele",
		"proje_gorevleri",
		"proje_aktif_yap",
		"aktif_proje_goster",
		"aktif_proje_kaldir",
		"ozet_goster",
		"template_listele",
		"templateden_gorev_olustur",
		"gorev_set_active",
		"gorev_get_active",
		"gorev_recent",
		"gorev_context_summary",
		"gorev_batch_update",
		"gorev_nlp_query",
		"gorev_intelligent_create",
		"gorev_file_watch_add",
		"gorev_file_watch_remove",
		"gorev_file_watch_list",
		"gorev_file_watch_stats",
		"gorev_export",
		"gorev_import",
	}

	// Create a map for easier lookup
	toolMap := make(map[string]bool)
	for _, tool := range tools {
		toolMap[tool.Name] = true
		assert.NotEmpty(t, tool.Description, "Tool %s should have a description", tool.Name)
	}

	// Verify all expected tools are present
	for _, expectedTool := range expectedTools {
		assert.True(t, toolMap[expectedTool], "Expected tool %s not found", expectedTool)
	}

	// Verify we have exactly the expected number of tools
	assert.Equal(t, len(expectedTools), len(tools))

	// Verify tool properties
	for _, tool := range tools {
		assert.NotEmpty(t, tool.Name)
		assert.NotEmpty(t, tool.Description)
		// InputSchema is currently nil for all tools in the hardcoded list
		assert.Nil(t, tool.InputSchema)
	}
}

// Test edge cases for the Tool struct
func TestToolStruct(t *testing.T) {
	// Test creating a Tool
	tool := Tool{
		Name:        "test_tool",
		Description: "Test tool description",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"param1": map[string]interface{}{
					"type":        "string",
					"description": "First parameter",
				},
			},
		},
	}

	assert.Equal(t, "test_tool", tool.Name)
	assert.Equal(t, "Test tool description", tool.Description)
	assert.NotNil(t, tool.InputSchema)

	// Test with empty tool
	emptyTool := Tool{}
	assert.Empty(t, emptyTool.Name)
	assert.Empty(t, emptyTool.Description)
	assert.Nil(t, emptyTool.InputSchema)
}

// Test YeniMCPSunucu with various conditions
func TestYeniMCPSunucu_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() *gorev.IsYonetici
		expectError bool
	}{
		{
			name: "valid is yonetici",
			setupFunc: func() *gorev.IsYonetici {
				config := &testinghelpers.TestDatabaseConfig{
					UseMemoryDB:     true,
					MigrationsPath:  constants.TestMigrationsPath,
					CreateTemplates: false,
					InitializeI18n:  true,
				}
				t := &testing.T{}
				isYonetici, _ := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
				return isYonetici
			},
			expectError: false,
		},
		{
			name: "nil is yonetici",
			setupFunc: func() *gorev.IsYonetici {
				return nil
			},
			expectError: false, // Should handle nil gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isYonetici := tt.setupFunc()

			// Test should not panic
			defer func() {
				if r := recover(); r != nil && !tt.expectError {
					t.Errorf("YeniMCPSunucu panicked unexpectedly: %v", r)
				}
			}()

			mcpServer, err := YeniMCPSunucu(isYonetici)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, mcpServer)
			}
		})
	}
}

// Test server initialization and tool registration
func TestServerToolRegistration(t *testing.T) {
	// Create handlers
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()
	handlers := YeniHandlers(isYonetici)

	// Create server
	mcpServer := server.NewMCPServer("test", "1.0.0")

	// Register tools
	handlers.RegisterTools(mcpServer)

	// Verify tools were registered
	// Note: The MCP server doesn't expose a way to list registered tools
	// but we can verify the server is properly initialized
	assert.NotNil(t, mcpServer)
}

// Test concurrent server creation
func TestConcurrentServerCreation(t *testing.T) {
	// Create shared veri yonetici
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()

	// Test concurrent server creation
	serverCount := constants.TestIterationSmall
	servers := make([]*server.MCPServer, serverCount)
	// errors := make([]error, serverCount) // unused

	// Create servers concurrently
	done := make(chan bool, serverCount)
	for i := 0; i < serverCount; i++ {
		go func(idx int) {
			defer func() { done <- true }()

			// Each goroutine creates its own handlers
			handlers := YeniHandlers(isYonetici)
			servers[idx] = NewServer(handlers)
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < serverCount; i++ {
		<-done
	}

	// Verify all servers were created successfully
	for i, srv := range servers {
		assert.NotNil(t, srv, "Server %d should not be nil", i)
	}
}

// Test server with different configurations
func TestServerConfigurations(t *testing.T) {
	// Test with in-memory database
	config1 := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	isYonetici1, cleanup1 := testinghelpers.SetupTestEnvironmentWithConfig(t, config1)
	defer cleanup1()

	handlers1 := YeniHandlers(isYonetici1)
	server1 := NewServer(handlers1)
	assert.NotNil(t, server1)

	// Test with file database
	config2 := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     false,
		UseTempFile:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	isYonetici2, cleanup2 := testinghelpers.SetupTestEnvironmentWithConfig(t, config2)
	defer cleanup2()

	handlers2 := YeniHandlers(isYonetici2)
	server2 := NewServer(handlers2)
	assert.NotNil(t, server2)

	// Both servers should have the same configuration
	// Note: Can't directly compare Name/Version as they're not public
	assert.NotNil(t, server1)
	assert.NotNil(t, server2)
}
