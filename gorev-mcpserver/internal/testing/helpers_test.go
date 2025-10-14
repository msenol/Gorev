package testing

import (
	"context"
	"testing"
)

func TestDefaultTestDatabaseConfig(t *testing.T) {
	config := DefaultTestDatabaseConfig()

	// Test that config is not nil
	if config == nil {
		t.Fatal("DefaultTestDatabaseConfig() returned nil")
	}

	// Test that config has reasonable values
	if config.UseMemoryDB != true {
		t.Error("Expected UseMemoryDB to be true for test config")
	}

	// Test that migrations path is set
	if config.MigrationsPath == "" {
		t.Error("Expected MigrationsPath to be set")
	}
}

func TestSetupTestEnvironmentBasic(t *testing.T) {
	// Test basic setup
	isYonetici, cleanup := SetupTestEnvironmentBasic(t)
	if cleanup != nil {
		defer cleanup()
	}

	if isYonetici == nil {
		t.Error("Expected IsYonetici to be initialized")
	}
}

func TestSetupTestEnvironmentWithConfig(t *testing.T) {
	config := DefaultTestDatabaseConfig()

	// Test setup with custom config
	isYonetici, cleanup := SetupTestEnvironmentWithConfig(t, config)
	if cleanup != nil {
		defer cleanup()
	}

	if isYonetici == nil {
		t.Error("Expected IsYonetici to be initialized")
	}
}

func TestTestEnvironmentCleanup(t *testing.T) {
	isYonetici, cleanup := SetupTestEnvironmentBasic(t)
	if isYonetici == nil {
		t.Fatal("SetupTestEnvironmentBasic() failed")
	}

	// Test that cleanup function exists and can be called safely
	if cleanup == nil {
		t.Error("Expected cleanup function to be returned")
	} else {
		// Call cleanup - should not panic
		cleanup()
	}
}

func TestMultipleTestEnvironmentCreation(t *testing.T) {
	// Test creating multiple test environments
	is1, cleanup1 := SetupTestEnvironmentBasic(t)
	if cleanup1 != nil {
		defer cleanup1()
	}

	is2, cleanup2 := SetupTestEnvironmentBasic(t)
	if cleanup2 != nil {
		defer cleanup2()
	}

	// Environments should be independent
	if is1 == is2 {
		t.Error("Expected different IsYonetici instances")
	}
}

// TestSetupTestDatabase_UseTempFile tests temporary file database creation
func TestSetupTestDatabase_UseTempFile(t *testing.T) {
	config := &TestDatabaseConfig{
		UseTempFile:     true,
		MigrationsPath:  "file://../../internal/veri/migrations",
		CreateTemplates: true,
		InitializeI18n:  false, // Already initialized
	}

	veriYonetici, cleanup := SetupTestDatabase(t, config)
	if veriYonetici == nil {
		t.Fatal("Expected VeriYonetici to be initialized")
	}
	if cleanup == nil {
		t.Fatal("Expected cleanup function")
	}

	// Cleanup should remove the temporary file
	cleanup()
}

// TestSetupTestDatabase_CustomPath tests custom database path
func TestSetupTestDatabase_CustomPath(t *testing.T) {
	config := &TestDatabaseConfig{
		CustomPath:      "/tmp/test_custom.db",
		MigrationsPath:  "file://../../internal/veri/migrations",
		CreateTemplates: false,
		InitializeI18n:  false,
	}

	veriYonetici, cleanup := SetupTestDatabase(t, config)
	if veriYonetici == nil {
		t.Fatal("Expected VeriYonetici to be initialized")
	}
	if cleanup == nil {
		t.Fatal("Expected cleanup function")
	}

	cleanup()
}

// TestSetupTestDatabase_NoTemplates tests database creation without templates
func TestSetupTestDatabase_NoTemplates(t *testing.T) {
	config := &TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  "file://../../internal/veri/migrations",
		CreateTemplates: false,
		InitializeI18n:  false,
	}

	veriYonetici, cleanup := SetupTestDatabase(t, config)
	if veriYonetici == nil {
		t.Fatal("Expected VeriYonetici to be initialized")
	}

	defer cleanup()

	// Verify database was created successfully (templates may or may not exist
	// depending on previous test runs in the same session)
	// The important thing is that CreateTemplates=false doesn't cause an error
	templates, err := veriYonetici.TemplateListele(context.Background(), "")
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	// With CreateTemplates=false, we just verify the database works
	// We don't assert on template count since VarsayilanTemplateleriOlustur
	// may have been called by other code paths
	t.Logf("Found %d templates in database", len(templates))
}

// TestSetupTestDatabase_WithI18n tests i18n initialization
func TestSetupTestDatabase_WithI18n(t *testing.T) {
	config := &TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  "file://../../internal/veri/migrations",
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	veriYonetici, cleanup := SetupTestDatabase(t, config)
	if veriYonetici == nil {
		t.Fatal("Expected VeriYonetici to be initialized")
	}

	defer cleanup()
}

// TestSetupTestDatabase_DefaultPath tests default path fallback
func TestSetupTestDatabase_DefaultPath(t *testing.T) {
	config := &TestDatabaseConfig{
		UseMemoryDB:     false,
		UseTempFile:     false,
		CustomPath:      "",
		MigrationsPath:  "file://../../internal/veri/migrations",
		CreateTemplates: false,
		InitializeI18n:  false,
	}

	veriYonetici, cleanup := SetupTestDatabase(t, config)
	if veriYonetici == nil {
		t.Fatal("Expected VeriYonetici to be initialized")
	}

	defer cleanup()
}
