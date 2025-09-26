package testing

import (
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