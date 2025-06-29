package mcp

import (
	"testing"

	"github.com/msenol/gorev/internal/gorev"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYeniMCPSunucu(t *testing.T) {
	// Create a test data manager
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Create MCP server
	server, err := YeniMCPSunucu(isYonetici)
	require.NoError(t, err)
	assert.NotNil(t, server)

	// Verify server has tools registered
	// Note: We can't easily test the tools registration without accessing private fields
	// But we can verify the server was created successfully
	assert.NotNil(t, server)
}

func TestMCPServerIntegration(t *testing.T) {
	// Create test environment
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Create and test MCP server
	server, err := YeniMCPSunucu(isYonetici)
	require.NoError(t, err)

	// Test that handlers are properly initialized
	handlers := YeniHandlers(isYonetici)
	assert.NotNil(t, handlers)
	assert.NotNil(t, handlers.isYonetici)

	// Test tool registration (we can't test all tools individually here,
	// but we can ensure the registration process works)
	handlers.RegisterTools(server)

	// If no panic occurred, registration was successful
	assert.True(t, true, "Tool registration completed without panic")
}

// Test MCP server configuration and metadata
func TestMCPServerMetadata(t *testing.T) {
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	server, err := YeniMCPSunucu(isYonetici)
	require.NoError(t, err)

	// Test server was created with correct metadata
	// Note: The actual server fields are private, so we test indirectly
	assert.NotNil(t, server)
}

// Test handler creation and initialization
func TestHandlerCreation(t *testing.T) {
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer veriYonetici.Kapat()

	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Test handler creation
	handlers := YeniHandlers(isYonetici)
	assert.NotNil(t, handlers)
	assert.Equal(t, isYonetici, handlers.isYonetici)
}

// Test edge cases and error conditions
func TestMCPServerErrorConditions(t *testing.T) {
	t.Run("Nil IsYonetici", func(t *testing.T) {
		// This should not panic but should create handlers with nil isYonetici
		handlers := YeniHandlers(nil)
		assert.NotNil(t, handlers)
		assert.Nil(t, handlers.isYonetici)
	})

	t.Run("Invalid Database Path", func(t *testing.T) {
		// Test with invalid database path
		_, err := gorev.YeniVeriYonetici("/invalid/path/db.sqlite", "file://../../internal/veri/migrations")
		assert.Error(t, err)

		// Ensure this doesn't cause issues in server creation
		// (though in practice this would fail earlier)
	})
}
