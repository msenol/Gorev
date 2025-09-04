package gorev

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultIDEConfig(t *testing.T) {
	config := DefaultIDEConfig()

	if config == nil {
		t.Fatal("DefaultIDEConfig returned nil")
	}

	if config.AutoInstall {
		t.Error("Default auto install should be false")
	}

	if config.AutoUpdate {
		t.Error("Default auto update should be false")
	}

	if config.CheckInterval != 24*time.Hour {
		t.Error("Default check interval should be 24 hours")
	}

	if config.ExtensionID != "mehmetsenol.gorev-vscode" {
		t.Error("Default extension ID mismatch")
	}

	if config.DisablePrompts {
		t.Error("Default disable prompts should be false")
	}

	expectedIDEs := []string{"vscode", "cursor", "windsurf"}
	if len(config.SupportedIDEs) != len(expectedIDEs) {
		t.Errorf("Expected %d supported IDEs, got %d", len(expectedIDEs), len(config.SupportedIDEs))
	}
}

func TestIDEConfigManagerCreation(t *testing.T) {
	manager := NewIDEConfigManager()

	if manager == nil {
		t.Fatal("NewIDEConfigManager returned nil")
	}

	if manager.config == nil {
		t.Error("Config should be initialized")
	}

	if manager.configPath == "" {
		t.Error("Config path should not be empty")
	}
}

func TestIDEConfigSaveLoad(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()

	// Create manager with custom config path
	manager := &IDEConfigManager{
		configPath: filepath.Join(tempDir, "test-config.json"),
		config:     DefaultIDEConfig(),
	}

	// Modify config
	manager.config.AutoInstall = true
	manager.config.AutoUpdate = true
	manager.config.CheckInterval = 12 * time.Hour

	// Save config
	err := manager.SaveConfig()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(manager.configPath); os.IsNotExist(err) {
		t.Error("Config file should exist after save")
	}

	// Create new manager and load
	newManager := &IDEConfigManager{
		configPath: manager.configPath,
		config:     DefaultIDEConfig(),
	}

	err = newManager.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded config
	if !newManager.config.AutoInstall {
		t.Error("AutoInstall should be true after load")
	}

	if !newManager.config.AutoUpdate {
		t.Error("AutoUpdate should be true after load")
	}

	if newManager.config.CheckInterval != 12*time.Hour {
		t.Error("CheckInterval should be 12 hours after load")
	}
}

func TestIDEConfigLoadNonExistent(t *testing.T) {
	tempDir := t.TempDir()

	manager := &IDEConfigManager{
		configPath: filepath.Join(tempDir, "nonexistent-config.json"),
		config:     DefaultIDEConfig(),
	}

	// Loading non-existent config should create default config
	err := manager.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should create default config: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(manager.configPath); os.IsNotExist(err) {
		t.Error("Config file should be created when loading non-existent config")
	}
}

func TestIDEConfigSetters(t *testing.T) {
	tempDir := t.TempDir()

	manager := &IDEConfigManager{
		configPath: filepath.Join(tempDir, "test-config.json"),
		config:     DefaultIDEConfig(),
	}

	// Test SetAutoInstall
	err := manager.SetAutoInstall(true)
	if err != nil {
		t.Fatalf("SetAutoInstall failed: %v", err)
	}
	if !manager.config.AutoInstall {
		t.Error("AutoInstall should be true")
	}

	// Test SetAutoUpdate
	err = manager.SetAutoUpdate(true)
	if err != nil {
		t.Fatalf("SetAutoUpdate failed: %v", err)
	}
	if !manager.config.AutoUpdate {
		t.Error("AutoUpdate should be true")
	}

	// Test SetDisablePrompts
	err = manager.SetDisablePrompts(true)
	if err != nil {
		t.Fatalf("SetDisablePrompts failed: %v", err)
	}
	if !manager.config.DisablePrompts {
		t.Error("DisablePrompts should be true")
	}

	// Test SetCheckInterval
	err = manager.SetCheckInterval(6 * time.Hour)
	if err != nil {
		t.Fatalf("SetCheckInterval failed: %v", err)
	}
	if manager.config.CheckInterval != 6*time.Hour {
		t.Error("CheckInterval should be 6 hours")
	}
}

func TestShouldCheckForUpdates(t *testing.T) {
	manager := NewIDEConfigManager()

	// Never checked before - should check
	if !manager.ShouldCheckForUpdates() {
		t.Error("Should check for updates when never checked before")
	}

	// Set last check to now
	manager.config.LastUpdateCheck = time.Now()

	// Should not check (interval not passed)
	if manager.ShouldCheckForUpdates() {
		t.Error("Should not check for updates when interval hasn't passed")
	}

	// Set last check to old time
	manager.config.LastUpdateCheck = time.Now().Add(-25 * time.Hour)

	// Should check (interval passed)
	if !manager.ShouldCheckForUpdates() {
		t.Error("Should check for updates when interval has passed")
	}
}

func TestUpdateLastCheckTime(t *testing.T) {
	tempDir := t.TempDir()

	manager := &IDEConfigManager{
		configPath: filepath.Join(tempDir, "test-config.json"),
		config:     DefaultIDEConfig(),
	}

	// Initial check time should be zero
	if !manager.config.LastUpdateCheck.IsZero() {
		t.Error("Initial last update check should be zero")
	}

	// Update check time
	beforeUpdate := time.Now()
	err := manager.UpdateLastCheckTime()
	if err != nil {
		t.Fatalf("UpdateLastCheckTime failed: %v", err)
	}
	afterUpdate := time.Now()

	// Check time should be between before and after
	checkTime := manager.config.LastUpdateCheck
	if checkTime.Before(beforeUpdate) || checkTime.After(afterUpdate) {
		t.Error("Last update check time should be within expected range")
	}
}

func TestIDEConfigJSONMarshaling(t *testing.T) {
	config := DefaultIDEConfig()
	config.AutoInstall = true
	config.CheckInterval = 2 * time.Hour

	// Marshal to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config to JSON: %v", err)
	}

	// Unmarshal from JSON
	var newConfig IDEConfig
	err = json.Unmarshal(data, &newConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal config from JSON: %v", err)
	}

	// Verify values
	if newConfig.AutoInstall != config.AutoInstall {
		t.Error("AutoInstall value mismatch after JSON round-trip")
	}

	if newConfig.CheckInterval != config.CheckInterval {
		t.Error("CheckInterval value mismatch after JSON round-trip")
	}

	if newConfig.ExtensionID != config.ExtensionID {
		t.Error("ExtensionID value mismatch after JSON round-trip")
	}
}

func BenchmarkConfigSaveLoad(b *testing.B) {
	tempDir := b.TempDir()

	manager := &IDEConfigManager{
		configPath: filepath.Join(tempDir, "benchmark-config.json"),
		config:     DefaultIDEConfig(),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Save and load config
		err := manager.SaveConfig()
		if err != nil {
			b.Fatalf("SaveConfig failed: %v", err)
		}

		err = manager.LoadConfig()
		if err != nil {
			b.Fatalf("LoadConfig failed: %v", err)
		}
	}
}
