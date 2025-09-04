package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/gorev"
)

// TestIDEIntegration_CompleteFlow tests the complete IDE management workflow
func TestIDEIntegration_CompleteFlow(t *testing.T) {
	// Skip this test in CI environments where IDEs won't be installed
	if os.Getenv("CI") == "true" || os.Getenv("GITHUB_ACTIONS") == "true" {
		t.Skip("Skipping IDE integration tests in CI environment")
	}

	t.Run("IDE Detection Flow", func(t *testing.T) {
		// Create IDE detector
		detector := gorev.NewIDEDetector()
		if detector == nil {
			t.Fatal("Failed to create IDE detector")
		}

		// Detect all IDEs
		detectedIDEs, err := detector.DetectAllIDEs()
		if err != nil {
			t.Errorf("IDE detection failed: %v", err)
		}

		// Log detected IDEs for debugging
		t.Logf("Detected %d IDEs", len(detectedIDEs))
		for ideType, ide := range detectedIDEs {
			t.Logf("Found %s: %s at %s (version: %s)", ideType, ide.Name, ide.ExecutablePath, ide.Version)
		}

		// Test individual IDE retrieval
		for ideType := range detectedIDEs {
			ide, exists := detector.GetDetectedIDE(ideType)
			if !exists {
				t.Errorf("GetDetectedIDE should return true for detected IDE %s", ideType)
			}
			if ide == nil {
				t.Errorf("GetDetectedIDE should return IDE info for %s", ideType)
			}
		}

		// Test GetAllDetectedIDEs
		allIDEs := detector.GetAllDetectedIDEs()
		if len(allIDEs) != len(detectedIDEs) {
			t.Errorf("GetAllDetectedIDEs returned %d IDEs, expected %d", len(allIDEs), len(detectedIDEs))
		}
	})

	t.Run("Configuration Management Flow", func(t *testing.T) {
		// Create config manager
		configManager := gorev.NewIDEConfigManager()
		if configManager == nil {
			t.Fatal("Failed to create IDE config manager")
		}

		// Test default configuration
		if !configManager.ShouldCheckForUpdates() {
			t.Log("Note: First run should check for updates")
		}

		// Test configuration setters
		originalAutoInstall := configManager.GetConfig().AutoInstall
		err := configManager.SetAutoInstall(true)
		if err != nil {
			t.Errorf("SetAutoInstall failed: %v", err)
		}

		if !configManager.GetConfig().AutoInstall {
			t.Error("AutoInstall should be true after SetAutoInstall(true)")
		}

		// Restore original setting
		configManager.SetAutoInstall(originalAutoInstall)

		// Test update check timing
		err = configManager.UpdateLastCheckTime()
		if err != nil {
			t.Errorf("UpdateLastCheckTime failed: %v", err)
		}

		// After updating, should not need immediate check
		if configManager.ShouldCheckForUpdates() {
			t.Error("Should not need immediate update check after UpdateLastCheckTime")
		}
	})

	t.Run("Extension Installer Flow", func(t *testing.T) {
		detector := gorev.NewIDEDetector()
		installer := gorev.NewExtensionInstaller(detector)

		if installer == nil {
			t.Fatal("Failed to create extension installer")
		}

		// Test cleanup functionality
		err := installer.Cleanup()
		if err != nil {
			t.Errorf("Cleanup failed: %v", err)
		}

		// Create new installer after cleanup
		installer = gorev.NewExtensionInstaller(detector)

		// Test with mock extension info
		extensionInfo := &gorev.ExtensionInfo{
			ID:          "test.extension",
			Name:        "Test Extension",
			Version:     "1.0.0",
			DownloadURL: "https://example.com/nonexistent.vsix",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// This should fail due to invalid URL, but should handle gracefully
		result, err := installer.InstallExtension(ctx, gorev.IDETypeVSCode, extensionInfo)
		if result == nil {
			t.Error("InstallExtension should return result even on failure")
		}

		if result != nil && result.Success {
			t.Error("Installation should fail with invalid URL")
		}

		// Test installation to all IDEs
		results, err := installer.InstallToAllIDEs(ctx, extensionInfo)
		if err != nil {
			// In test environment, this may fail due to no IDEs detected
			t.Logf("InstallToAllIDEs failed as expected in test environment: %v", err)
		}

		if len(results) == 0 {
			t.Log("Note: No IDEs detected, which is expected in test environment")
		}

		// Test uninstall
		uninstallResult, err := installer.UninstallExtension(gorev.IDETypeVSCode, "test.extension")
		if uninstallResult == nil {
			t.Error("UninstallExtension should return result even on failure")
		}
	})
}

// TestIDEIntegration_MockedExternalCalls tests IDE management with mocked external dependencies
func TestIDEIntegration_MockedExternalCalls(t *testing.T) {
	t.Run("Mocked GitHub API Integration", func(t *testing.T) {
		// Create mock GitHub API server
		mockRelease := map[string]interface{}{
			"tag_name": "v1.2.3",
			"name":     "Test Release v1.2.3",
			"assets": []map[string]interface{}{
				{
					"name":                 "gorev-vscode-1.2.3.vsix",
					"browser_download_url": "http://example.com/test.vsix",
				},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/releases/latest") {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(mockRelease)
			} else {
				w.WriteHeader(404)
			}
		}))
		defer server.Close()

		detector := gorev.NewIDEDetector()
		installer := gorev.NewExtensionInstaller(detector)

		ctx := context.Background()

		// Note: This will still fail because we can't easily mock the actual GitHub API call
		// but it tests the error handling pathway
		info, err := installer.GetLatestExtensionInfo(ctx, "test", "repo")
		if err == nil && info != nil {
			t.Log("Note: GitHub API call succeeded (unexpected in test environment)")
		} else {
			t.Logf("Expected error from GitHub API: %v", err)
		}
	})

	t.Run("Mocked VSIX Download", func(t *testing.T) {
		// Create mock VSIX content
		mockVSIXContent := "PK\x03\x04test-vsix-content"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/zip")
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(mockVSIXContent)))
			w.Write([]byte(mockVSIXContent))
		}))
		defer server.Close()

		detector := gorev.NewIDEDetector()
		installer := gorev.NewExtensionInstaller(detector)

		extensionInfo := &gorev.ExtensionInfo{
			ID:          "test.extension",
			Name:        "Test Extension",
			Version:     "1.0.0",
			DownloadURL: server.URL + "/test.vsix",
		}

		ctx := context.Background()

		// This should work with our mock server
		result, _ := installer.InstallExtension(ctx, gorev.IDETypeVSCode, extensionInfo)

		// Installation will likely fail because VS Code is not detected,
		// but download should succeed
		if result == nil {
			t.Error("Should return result")
		}

		if result != nil && result.Success {
			t.Log("Installation succeeded (VS Code must be installed)")
		} else {
			t.Log("Installation failed as expected (VS Code not detected or installation failed)")
		}
	})

	t.Run("File System Operations", func(t *testing.T) {
		detector := gorev.NewIDEDetector()
		installer := gorev.NewExtensionInstaller(detector)

		// Test download directory creation
		downloadPath := installer.GetDownloadPath()
		if downloadPath == "" {
			t.Error("Download path should not be empty")
		}

		// Verify directory exists
		if _, err := os.Stat(downloadPath); os.IsNotExist(err) {
			t.Error("Download directory should exist")
		}

		// Test cleanup
		testFile := filepath.Join(downloadPath, "test-cleanup.txt")
		err := os.WriteFile(testFile, []byte("test"), 0644)
		if err != nil {
			t.Errorf("Failed to create test file: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Error("Test file should exist")
		}

		// Run cleanup
		err = installer.Cleanup()
		if err != nil {
			t.Errorf("Cleanup failed: %v", err)
		}

		// Verify directory is cleaned up
		if _, err := os.Stat(downloadPath); !os.IsNotExist(err) {
			t.Error("Download directory should be removed after cleanup")
		}
	})
}

// TestIDEIntegration_ErrorHandling tests error scenarios in IDE management
func TestIDEIntegration_ErrorHandling(t *testing.T) {
	t.Run("Network Errors", func(t *testing.T) {
		detector := gorev.NewIDEDetector()
		installer := gorev.NewExtensionInstaller(detector)

		// Test with invalid URL
		extensionInfo := &gorev.ExtensionInfo{
			ID:          "test.extension",
			DownloadURL: "invalid-url-format",
		}

		ctx := context.Background()
		result, err := installer.InstallExtension(ctx, gorev.IDETypeVSCode, extensionInfo)

		if result == nil {
			t.Error("Should return result even on network error")
		}

		if result != nil && result.Success {
			t.Error("Should not succeed with invalid URL")
		}

		if err == nil {
			t.Log("Note: No error returned, handled gracefully")
		} else {
			t.Logf("Expected network error: %v", err)
		}
	})

	t.Run("Timeout Handling", func(t *testing.T) {
		// Create slow server
		slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)
			w.Write([]byte("slow response"))
		}))
		defer slowServer.Close()

		detector := gorev.NewIDEDetector()
		installer := gorev.NewExtensionInstaller(detector)

		extensionInfo := &gorev.ExtensionInfo{
			ID:          "test.extension",
			DownloadURL: slowServer.URL + "/slow.vsix",
		}

		// Create context with very short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		result, err := installer.InstallExtension(ctx, gorev.IDETypeVSCode, extensionInfo)

		if result == nil {
			t.Error("Should return result even on timeout")
		}

		// Should either timeout or complete quickly
		if err != nil && strings.Contains(err.Error(), "context") {
			t.Log("Context timeout handled correctly")
		} else if result != nil && !result.Success {
			t.Log("Installation failed as expected")
		}
	})

	t.Run("Invalid Extension Data", func(t *testing.T) {
		detector := gorev.NewIDEDetector()
		installer := gorev.NewExtensionInstaller(detector)

		testCases := []struct {
			name      string
			extension *gorev.ExtensionInfo
		}{
			{
				name:      "Nil extension info",
				extension: nil,
			},
			{
				name: "Empty extension ID",
				extension: &gorev.ExtensionInfo{
					ID:          "",
					DownloadURL: "https://example.com/test.vsix",
				},
			},
			{
				name: "Empty download URL",
				extension: &gorev.ExtensionInfo{
					ID:          "test.extension",
					DownloadURL: "",
				},
			},
		}

		ctx := context.Background()
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, _ := installer.InstallExtension(ctx, gorev.IDETypeVSCode, tc.extension)

				if result == nil && tc.extension == nil {
					// Nil extension should be handled gracefully
					return
				}

				if result == nil {
					t.Error("Should return result even with invalid data")
				}

				if result != nil && result.Success {
					t.Error("Should not succeed with invalid extension data")
				}
			})
		}
	})
}

// TestIDEIntegration_Performance tests performance aspects of IDE management
func TestIDEIntegration_Performance(t *testing.T) {
	t.Run("Detection Performance", func(t *testing.T) {
		start := time.Now()

		detector := gorev.NewIDEDetector()
		detectedIDEs, err := detector.DetectAllIDEs()

		elapsed := time.Since(start)

		if err != nil {
			t.Errorf("Detection failed: %v", err)
		}

		t.Logf("Detected %d IDEs in %v", len(detectedIDEs), elapsed)

		if elapsed > 5*time.Second {
			t.Errorf("IDE detection took too long: %v", elapsed)
		}
	})

	t.Run("Config Management Performance", func(t *testing.T) {
		configManager := gorev.NewIDEConfigManager()

		start := time.Now()

		// Perform multiple config operations
		for i := 0; i < 10; i++ {
			configManager.SetAutoInstall(i%2 == 0)
			configManager.SetAutoUpdate(i%2 == 1)
			configManager.UpdateLastCheckTime()
		}

		elapsed := time.Since(start)

		t.Logf("10 config operations completed in %v", elapsed)

		if elapsed > time.Second {
			t.Errorf("Config operations took too long: %v", elapsed)
		}
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		detector := gorev.NewIDEDetector()

		// Test concurrent detection calls
		numGoroutines := 5
		results := make(chan map[gorev.IDEType]*gorev.IDEInfo, numGoroutines)
		errors := make(chan error, numGoroutines)

		start := time.Now()

		for i := 0; i < numGoroutines; i++ {
			go func() {
				ides, err := detector.DetectAllIDEs()
				if err != nil {
					errors <- err
					return
				}
				results <- ides
			}()
		}

		// Collect results
		var allResults []map[gorev.IDEType]*gorev.IDEInfo
		for i := 0; i < numGoroutines; i++ {
			select {
			case result := <-results:
				allResults = append(allResults, result)
			case err := <-errors:
				t.Errorf("Concurrent detection failed: %v", err)
			case <-time.After(10 * time.Second):
				t.Fatal("Concurrent detection timed out")
			}
		}

		elapsed := time.Since(start)

		t.Logf("Completed %d concurrent detections in %v", numGoroutines, elapsed)

		// Verify all results are consistent
		if len(allResults) > 1 {
			firstResult := allResults[0]
			for i, result := range allResults[1:] {
				if len(result) != len(firstResult) {
					t.Errorf("Result %d has different number of IDEs: %d vs %d", i+1, len(result), len(firstResult))
				}
			}
		}
	})
}

// Benchmark tests
func BenchmarkIDEDetection(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector := gorev.NewIDEDetector()
		detector.DetectAllIDEs()
	}
}

func BenchmarkConfigOperations(b *testing.B) {
	configManager := gorev.NewIDEConfigManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		configManager.SetAutoInstall(i%2 == 0)
		configManager.UpdateLastCheckTime()
	}
}

// Helper functions for testing
func createMockVSIXContent() []byte {
	// Create minimal ZIP-like content that resembles a VSIX file
	return []byte("PK\x03\x04\x14\x00\x00\x00mock-vsix-content-for-testing")
}

func setupTempDownloadDir(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "gorev-test-downloads-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}
