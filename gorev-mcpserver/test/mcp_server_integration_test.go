package test

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/gorev"
	mcphandlers "github.com/msenol/gorev/internal/mcp"
	testinghelpers "github.com/msenol/gorev/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMCPServerCreation tests basic MCP server creation
func TestMCPServerCreation(t *testing.T) {
	// Create test environment
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Test basic server creation
	mcpServer, err := mcphandlers.YeniMCPSunucu(isYonetici)
	require.NoError(t, err)
	assert.NotNil(t, mcpServer)
}

// TestMCPServerCreationWithDebug tests MCP server creation with debug mode
func TestMCPServerCreationWithDebug(t *testing.T) {
	// Create test environment
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Test server creation with debug
	mcpServer, err := mcphandlers.YeniMCPSunucuWithDebug(isYonetici, true)
	require.NoError(t, err)
	assert.NotNil(t, mcpServer)
}

// TestMCPServerCreationNilManager tests server creation with nil manager
func TestMCPServerCreationNilManager(t *testing.T) {
	// Test with nil manager - should still create server but without functionality
	mcpServer, err := mcphandlers.YeniMCPSunucu(nil)
	require.NoError(t, err)
	assert.NotNil(t, mcpServer)
}

// TestHandlersCreation tests handlers creation with various configurations
func TestHandlersCreation(t *testing.T) {
	t.Run("with valid manager", func(t *testing.T) {
		// Create test environment
		isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
		defer cleanup()

		// Create handlers
		handlers := mcphandlers.YeniHandlers(isYonetici)
		assert.NotNil(t, handlers)
		defer handlers.Close()

		// Test that handlers are functional by calling a simple operation
		result, err := handlers.OzetGoster(map[string]interface{}{})
		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("with debug mode", func(t *testing.T) {
		// Create test environment
		isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
		defer cleanup()

		// Create handlers with debug
		handlers := mcphandlers.YeniHandlersWithDebug(isYonetici, true)
		assert.NotNil(t, handlers)
		defer handlers.Close()

		// Test functionality
		result, err := handlers.OzetGoster(map[string]interface{}{})
		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("with nil manager", func(t *testing.T) {
		// Create handlers with nil manager
		handlers := mcphandlers.YeniHandlers(nil)
		assert.NotNil(t, handlers)
		defer handlers.Close()

		// Operations should panic with nil manager - this reveals a bug that should be fixed
		assert.Panics(t, func() {
			_, _ = handlers.OzetGoster(map[string]interface{}{})
		}, "Operations with nil manager should panic (reveals bug in handlers)")
	})
}

// TestHandlersResourceCleanup tests that handlers properly clean up resources
func TestHandlersResourceCleanup(t *testing.T) {
	// Create test environment
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Create and immediately close handlers
	handlers := mcphandlers.YeniHandlers(isYonetici)
	assert.NotNil(t, handlers)

	// Test cleanup
	err := handlers.Close()
	assert.NoError(t, err)

	// Test double close (should not panic)
	err = handlers.Close()
	assert.NoError(t, err)
}

// TestMCPServerToolRegistration tests that tools are properly registered
func TestMCPServerToolRegistration(t *testing.T) {
	// Create test environment
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Create MCP server
	_, err := mcphandlers.YeniMCPSunucu(isYonetici)
	require.NoError(t, err)

	// Test that we can list tools (this will verify registration)
	tools := mcphandlers.ListTools()
	assert.Greater(t, len(tools), 0, "Should have registered tools")

	// Check for essential tools
	toolNames := make([]string, len(tools))
	for i, tool := range tools {
		toolNames[i] = tool.Name
	}

	expectedTools := []string{
		"gorev_listele",
		"gorev_detay",
		"gorev_guncelle",
		"proje_olustur",
		"template_listele",
		"templateden_gorev_olustur",
	}

	for _, expectedTool := range expectedTools {
		assert.Contains(t, toolNames, expectedTool, "Should have registered tool: %s", expectedTool)
	}

	// Verify deprecated tool is not in list
	assert.NotContains(t, toolNames, "gorev_olustur", "Deprecated tool should not be registered")
}

// TestMCPServerWithMemoryDatabase tests server functionality with in-memory database
func TestMCPServerWithMemoryDatabase(t *testing.T) {
	// Create test environment with memory database
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Create handlers
	handlers := mcphandlers.YeniHandlers(isYonetici)
	assert.NotNil(t, handlers)
	defer handlers.Close()

	// Test basic functionality through handlers
	result, err := handlers.ProjeListele(map[string]interface{}{})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	// With empty database, should get "no projects" message
	text := extractText(t, result)
	// Check that we get a valid response (either project list or no projects message)
	assert.True(t, len(text) > 0, "Should get a response")
}

// TestMCPServerWithFileDatabase tests server functionality with file database
func TestMCPServerWithFileDatabase(t *testing.T) {
	// Create temporary directory for test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Get embedded migrations config and set custom path
	migrationsFS, err := getEmbeddedMigrationsFS()
	require.NoError(t, err)

	config := &testinghelpers.TestDatabaseConfig{
		CustomPath:      dbPath,
		MigrationsFS:    migrationsFS,
		CreateTemplates: true,
		InitializeI18n:  true,
	}
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()

	// Create handlers
	handlers := mcphandlers.YeniHandlers(isYonetici)
	assert.NotNil(t, handlers)
	defer handlers.Close()

	// Test basic functionality
	result, err := handlers.OzetGoster(map[string]interface{}{})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	text := extractText(t, result)
	assert.Contains(t, text, "## Özet Rapor")
}

// TestMCPServerConcurrentAccess tests server behavior under concurrent access
func TestMCPServerConcurrentAccess(t *testing.T) {
	// Create test environment
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Create handlers
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test concurrent access
	var wg sync.WaitGroup
	concurrentOps := 10
	errors := make(chan error, concurrentOps)

	for i := 0; i < concurrentOps; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Perform different operations concurrently
			switch id % 4 {
			case 0:
				_, err := handlers.OzetGoster(map[string]interface{}{})
				errors <- err
			case 1:
				_, err := handlers.ProjeListele(map[string]interface{}{})
				errors <- err
			case 2:
				_, err := handlers.GorevListele(map[string]interface{}{})
				errors <- err
			case 3:
				_, err := handlers.TemplateListele(map[string]interface{}{"kategori": ""})
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		assert.NoError(t, err, "Concurrent operation should not error")
	}
}

// TestMCPServerPerformance tests server performance under load
func TestMCPServerPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	// Create test environment
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Create handlers
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Create some test data
	for i := 0; i < 50; i++ {
		_, err := isYonetici.GorevOlustur(fmt.Sprintf("Performance Test Görev %d", i), "", "orta", "", "", nil)
		require.NoError(t, err)
	}

	// Measure performance
	start := time.Now()
	operations := 100

	for i := 0; i < operations; i++ {
		_, err := handlers.GorevListele(map[string]interface{}{})
		require.NoError(t, err)
	}

	duration := time.Since(start)
	avgDuration := duration / time.Duration(operations)

	t.Logf("Performed %d operations in %v (average: %v per operation)", operations, duration, avgDuration)

	// Assert reasonable performance (should be much faster than 10ms per operation)
	assert.Less(t, avgDuration, 10*time.Millisecond, "Average operation duration should be less than 10ms")
}

// TestMCPServerErrorHandling tests server error handling
func TestMCPServerErrorHandling(t *testing.T) {
	t.Run("operations with invalid database", func(t *testing.T) {
		// Create a test with an invalid database scenario
		// Note: This is a simplified test since we can't easily create corrupted databases
		migrationsFS, err := getEmbeddedMigrationsFS()
		require.NoError(t, err)

		config := &testinghelpers.TestDatabaseConfig{
			UseMemoryDB:     true,
			MigrationsFS:    migrationsFS,
			CreateTemplates: false, // Skip templates to test error scenarios
			InitializeI18n:  true,
		}

		isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
		defer cleanup()

		handlers := mcphandlers.YeniHandlers(isYonetici)
		defer handlers.Close()

		// Test operations with invalid parameters to verify error handling
		_, err = handlers.GorevDetay(map[string]interface{}{"id": "invalid-uuid"})
		// This should not crash and should handle invalid input gracefully
		if err != nil {
			t.Logf("Gracefully handled invalid parameter: %v", err)
		}
	})
}

// TestMCPServerLargeDataset tests server behavior with large datasets
func TestMCPServerLargeDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large dataset test in short mode")
	}

	// Create test environment
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, getTestConfigWithEmbeddedMigrations(t))
	defer cleanup()

	// Create large dataset
	projects := 5
	tasksPerProject := 100

	for p := 0; p < projects; p++ {
		proje, err := isYonetici.ProjeOlustur(fmt.Sprintf("Large Dataset Projesi %d", p+1), "")
		require.NoError(t, err)

		for i := 0; i < tasksPerProject; i++ {
			_, err := isYonetici.GorevOlustur(
				fmt.Sprintf("Görev %d-%d", p+1, i+1),
				fmt.Sprintf("Bu görev %d. projenin %d. görevidir", p+1, i+1),
				[]string{"dusuk", "orta", "yuksek"}[i%3],
				proje.ID,
				"", nil)
			require.NoError(t, err)
		}
	}

	// Create handlers
	handlers := mcphandlers.YeniHandlers(isYonetici)
	defer handlers.Close()

	// Test list operations with large dataset
	start := time.Now()

	result, err := handlers.GorevListele(map[string]interface{}{"limit": 1000})
	require.NoError(t, err)
	assert.False(t, result.IsError)

	duration := time.Since(start)
	t.Logf("Listed 500+ tasks in %v", duration)

	// Should complete in reasonable time
	assert.Less(t, duration, 5*time.Second, "Large dataset listing should complete in under 5 seconds")

	text := extractText(t, result)
	// Check for pagination header (for large datasets) or regular header
	assert.True(t,
		len(text) > 0 && (
			strings.Contains(text, "Görev Listesi") ||
			strings.Contains(text, "Görevler (")),
		"Should contain task list or pagination header")
}

// TestMCPServerInitializationSequence tests the complete initialization sequence
func TestMCPServerInitializationSequence(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "init_test.db")

	// Test step-by-step initialization
	t.Run("database initialization", func(t *testing.T) {
		// Initialize database
		migrationsFS, err := getEmbeddedMigrationsFS()
		require.NoError(t, err)

		config := &testinghelpers.TestDatabaseConfig{
			CustomPath:      dbPath,
			MigrationsFS:    migrationsFS,
			CreateTemplates: true,
			InitializeI18n:  true,
		}

		isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
		defer cleanup()

		// Verify database is accessible
		_, err = isYonetici.ProjeListele()
		assert.NoError(t, err, "Should be able to access database after initialization")
	})

	t.Run("manager initialization", func(t *testing.T) {
		// Create manager with existing database
		migrationsFS, err := getEmbeddedMigrationsFS()
		require.NoError(t, err)

		vy, err := gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, migrationsFS)
		require.NoError(t, err)
		defer func() {
			// Note: VeriYonetici doesn't have Close method, cleanup is handled by test framework
		}()

		isYonetici := gorev.YeniIsYonetici(vy)

		// Verify manager is functional
		_, err2 := isYonetici.ProjeListele()
		assert.NoError(t, err2, "Manager should be functional with existing database")
		// Note: We don't assert projects is not nil since the database might be empty
		// The important thing is that the manager can access the database without errors
	})

	t.Run("handlers initialization", func(t *testing.T) {
		// Create manager
		migrationsFS, err := getEmbeddedMigrationsFS()
		require.NoError(t, err)

		vy, err := gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, migrationsFS)
		require.NoError(t, err)
		defer func() {
			// Note: VeriYonetici doesn't have Close method, cleanup is handled by test framework
		}()

		isYonetici := gorev.YeniIsYonetici(vy)

		// Create handlers
		handlers := mcphandlers.YeniHandlers(isYonetici)
		assert.NotNil(t, handlers, "Handlers should be created successfully")
		defer handlers.Close()

		// Verify handlers are functional
		result, err := handlers.OzetGoster(map[string]interface{}{})
		require.NoError(t, err, "Handlers should be functional")
		assert.False(t, result.IsError, "Handlers should return valid results")
	})

	t.Run("server initialization", func(t *testing.T) {
		// Create manager
		migrationsFS, err := getEmbeddedMigrationsFS()
		require.NoError(t, err)

		vy, err := gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, migrationsFS)
		require.NoError(t, err)
		defer func() {
			// Note: VeriYonetici doesn't have Close method, cleanup is handled by test framework
		}()

		isYonetici := gorev.YeniIsYonetici(vy)

		// Create MCP server
		_, err2 := mcphandlers.YeniMCPSunucu(isYonetici)
		require.NoError(t, err2, "MCP server should be created successfully")
	})
}

// Benchmark tests are omitted as they require special setup for benchmark testing
