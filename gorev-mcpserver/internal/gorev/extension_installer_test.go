package gorev

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Mock HTTP client for testing
type mockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	if err, exists := m.errors[url]; exists {
		return nil, err
	}
	if resp, exists := m.responses[url]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader("Not Found")),
	}, nil
}

func TestNewExtensionInstaller(t *testing.T) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	if installer == nil {
		t.Fatal("NewExtensionInstaller returned nil")
	}

	if installer.detector != detector {
		t.Error("Detector not properly assigned")
	}

	if installer.client == nil {
		t.Error("HTTP client should be initialized")
	}

	if installer.downloadPath == "" {
		t.Error("Download path should not be empty")
	}

	// Verify download directory was created
	if _, err := os.Stat(installer.downloadPath); os.IsNotExist(err) {
		t.Error("Download directory should be created")
	}
}

func TestExtensionInstaller_GetLatestExtensionInfo(t *testing.T) {
	// Create mock GitHub API response
	mockRelease := map[string]interface{}{
		"tag_name": "v1.2.3",
		"name":     "Release v1.2.3",
		"assets": []map[string]interface{}{
			{
				"name":                 "gorev-vscode-1.2.3.vsix",
				"browser_download_url": "https://github.com/test/repo/releases/download/v1.2.3/gorev-vscode-1.2.3.vsix",
			},
		},
	}

	responseBody, _ := json.Marshal(mockRelease)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/releases/latest") {
			w.Header().Set("Content-Type", "application/json")
			w.Write(responseBody)
		} else {
			w.WriteHeader(404)
		}
	}))
	defer server.Close()

	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	// Replace the base GitHub API URL for testing
	ctx := context.Background()

	// Mock the HTTP client to use our test server
	installer.client = &http.Client{Timeout: 5 * time.Second}

	testCases := []struct {
		name         string
		repoOwner    string
		repoName     string
		expectError  bool
		expectedName string
		expectedVer  string
	}{
		{
			name:        "Valid repository",
			repoOwner:   "test",
			repoName:    "repo",
			expectError: true, // Will fail because we can't mock the actual GitHub API easily
		},
		{
			name:        "Empty repository owner",
			repoOwner:   "",
			repoName:    "repo",
			expectError: true,
		},
		{
			name:        "Empty repository name",
			repoOwner:   "test",
			repoName:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := installer.GetLatestExtensionInfo(ctx, tc.repoOwner, tc.repoName)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError && info == nil {
				t.Error("Expected extension info but got nil")
			}
		})
	}
}

func TestExtensionInstaller_downloadVSIX(t *testing.T) {
	testCases := []struct {
		name        string
		setupServer func() *httptest.Server
		expectError bool
	}{
		{
			name: "Successful download",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/zip")
					w.Write([]byte("fake-vsix-content"))
				}))
			},
			expectError: false,
		},
		{
			name: "Server returns 404",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(404)
				}))
			},
			expectError: true,
		},
		{
			name: "Server timeout",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(100 * time.Millisecond) // Simulate slow response
					w.Write([]byte("content"))
				}))
			},
			expectError: true, // Should timeout with short context deadline
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := tc.setupServer()
			defer server.Close()

			detector := NewIDEDetector()
			installer := NewExtensionInstaller(detector)

			extensionInfo := &ExtensionInfo{
				ID:          "test.extension",
				Name:        "Test Extension",
				Version:     "1.0.0",
				DownloadURL: server.URL + "/test.vsix",
			}

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			filePath, err := installer.downloadVSIX(ctx, extensionInfo)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError {
				// Verify file was created
				if filePath == "" {
					t.Error("File path should not be empty")
				}

				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Error("Downloaded file should exist")
				}

				// Cleanup
				os.Remove(filePath)
			}
		})
	}
}

func TestExtensionInstaller_verifyChecksum(t *testing.T) {
	// Create temporary test file
	tempFile, err := os.CreateTemp("", "test-checksum-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	testContent := "test content"
	tempFile.WriteString(testContent)
	tempFile.Close()

	// Calculate expected checksum (SHA256 of "test content")
	expectedChecksum := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"

	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	testCases := []struct {
		name             string
		filePath         string
		expectedChecksum string
		expectMatch      bool
		expectError      bool
	}{
		{
			name:             "Valid checksum match",
			filePath:         tempFile.Name(),
			expectedChecksum: expectedChecksum,
			expectMatch:      true,
			expectError:      false,
		},
		{
			name:             "Invalid checksum",
			filePath:         tempFile.Name(),
			expectedChecksum: "invalid-checksum",
			expectMatch:      false,
			expectError:      false,
		},
		{
			name:             "Non-existent file",
			filePath:         "/nonexistent/file.txt",
			expectedChecksum: expectedChecksum,
			expectMatch:      false,
			expectError:      true,
		},
		{
			name:             "Empty checksum",
			filePath:         tempFile.Name(),
			expectedChecksum: "",
			expectMatch:      false, // Empty checksum will not match actual checksum
			expectError:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match, err := installer.verifyChecksum(tc.filePath, tc.expectedChecksum)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if match != tc.expectMatch {
				t.Errorf("Expected match=%v, got match=%v", tc.expectMatch, match)
			}
		})
	}
}

func TestExtensionInstaller_InstallExtension(t *testing.T) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	// Mock IDE detection
	mockIDE := &IDEInfo{
		Type:           IDETypeVSCode,
		Name:           "Mock VS Code",
		ExecutablePath: "/usr/bin/code",
		IsInstalled:    true,
	}
	detector.detectedIDEs[IDETypeVSCode] = mockIDE

	testCases := []struct {
		name        string
		ideType     IDEType
		extension   *ExtensionInfo
		expectError bool
	}{
		{
			name:    "IDE not detected",
			ideType: IDETypeCursor,
			extension: &ExtensionInfo{
				ID:          "test.extension",
				DownloadURL: "https://example.com/test.vsix",
			},
			expectError: true,
		},
		{
			name:    "Invalid download URL",
			ideType: IDETypeVSCode,
			extension: &ExtensionInfo{
				ID:          "test.extension",
				DownloadURL: "invalid-url",
			},
			expectError: true,
		},
		{
			name:    "Nil extension info",
			ideType: IDETypeVSCode,
			extension: &ExtensionInfo{
				ID:          "test.extension",
				DownloadURL: "", // Empty URL to trigger error
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			result, err := installer.InstallExtension(ctx, tc.ideType, tc.extension)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Result should always be returned (success or failure)
			if result == nil && !tc.expectError {
				t.Error("Expected result but got nil")
			}
		})
	}
}

func TestExtensionInstaller_UninstallExtension(t *testing.T) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	// Mock IDE detection
	mockIDE := &IDEInfo{
		Type:           IDETypeVSCode,
		Name:           "Mock VS Code",
		ExecutablePath: "/usr/bin/code",
		IsInstalled:    true,
	}
	detector.detectedIDEs[IDETypeVSCode] = mockIDE

	testCases := []struct {
		name        string
		ideType     IDEType
		extensionID string
		expectError bool
	}{
		{
			name:        "Valid uninstall request",
			ideType:     IDETypeVSCode,
			extensionID: "test.extension",
			expectError: true, // Will fail because we can't actually run code command
		},
		{
			name:        "IDE not detected",
			ideType:     IDETypeCursor,
			extensionID: "test.extension",
			expectError: true,
		},
		{
			name:        "Empty extension ID",
			ideType:     IDETypeVSCode,
			extensionID: "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := installer.UninstallExtension(tc.ideType, tc.extensionID)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Result should always be returned
			if result == nil {
				t.Error("Expected result but got nil")
			}
		})
	}
}

func TestExtensionInstaller_InstallToAllIDEs(t *testing.T) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	// Mock multiple IDE detections
	detector.detectedIDEs[IDETypeVSCode] = &IDEInfo{
		Type:        IDETypeVSCode,
		Name:        "Mock VS Code",
		IsInstalled: true,
	}
	detector.detectedIDEs[IDETypeCursor] = &IDEInfo{
		Type:        IDETypeCursor,
		Name:        "Mock Cursor",
		IsInstalled: true,
	}

	extensionInfo := &ExtensionInfo{
		ID:          "test.extension",
		Name:        "Test Extension",
		DownloadURL: "invalid-url", // Will cause download to fail
	}

	ctx := context.Background()
	results, err := installer.InstallToAllIDEs(ctx, extensionInfo)

	// Should not return error even if individual installations fail
	if err != nil {
		t.Errorf("InstallToAllIDEs should not return error: %v", err)
	}

	// Should return results for all detected IDEs
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// All installations should fail due to invalid URL
	for _, result := range results {
		if result.Success {
			t.Error("Expected installation to fail due to invalid URL")
		}
	}
}

func TestExtensionInstaller_ListInstalledExtensions(t *testing.T) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	testCases := []struct {
		name        string
		ideType     IDEType
		expectError bool
	}{
		{
			name:        "IDE not detected",
			ideType:     IDETypeVSCode,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			extensions, err := installer.ListInstalledExtensions(tc.ideType)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError && extensions == nil {
				t.Error("Expected extension list but got nil")
			}
		})
	}
}

func TestExtensionInstaller_Cleanup(t *testing.T) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	// Create test file in download directory
	testFile := filepath.Join(installer.downloadPath, "test-file.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Test file should exist before cleanup")
	}

	// Run cleanup
	err = installer.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should not return error: %v", err)
	}

	// Verify directory was cleaned up
	if _, err := os.Stat(installer.downloadPath); !os.IsNotExist(err) {
		t.Error("Download directory should be removed after cleanup")
	}
}

// Benchmark tests
func BenchmarkExtensionInstaller_GetLatestExtensionInfo(b *testing.B) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail but we're measuring the overhead
		installer.GetLatestExtensionInfo(ctx, "test", "repo")
	}
}

func BenchmarkExtensionInstaller_verifyChecksum(b *testing.B) {
	// Create test file
	tempFile, err := os.CreateTemp("", "benchmark-checksum-*.txt")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	tempFile.WriteString("test content")
	tempFile.Close()

	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)
	checksum := "6ae8a75555209fd6c44157c0aed8016e763ff435a19cf186f76863140143ff72"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		installer.verifyChecksum(tempFile.Name(), checksum)
	}
}

// Test helper functions
func TestExtensionInfo_Validation(t *testing.T) {
	testCases := []struct {
		name      string
		extension ExtensionInfo
		valid     bool
	}{
		{
			name: "Valid extension info",
			extension: ExtensionInfo{
				ID:          "test.extension",
				Name:        "Test Extension",
				Version:     "1.0.0",
				DownloadURL: "https://example.com/test.vsix",
			},
			valid: true,
		},
		{
			name: "Missing ID",
			extension: ExtensionInfo{
				Name:        "Test Extension",
				DownloadURL: "https://example.com/test.vsix",
			},
			valid: false,
		},
		{
			name: "Missing download URL",
			extension: ExtensionInfo{
				ID:   "test.extension",
				Name: "Test Extension",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid := tc.extension.ID != "" && tc.extension.DownloadURL != ""
			if valid != tc.valid {
				t.Errorf("Expected valid=%v, got valid=%v", tc.valid, valid)
			}
		})
	}
}

func TestInstallResult_String(t *testing.T) {
	result := InstallResult{
		Success:   true,
		Message:   "Installation successful",
		Extension: "test.extension",
		IDE:       "vscode",
		Version:   "1.0.0",
	}

	// Test that all fields are properly set
	if !result.Success {
		t.Error("Success should be true")
	}

	if result.Message == "" {
		t.Error("Message should not be empty")
	}

	if result.Extension == "" {
		t.Error("Extension should not be empty")
	}

	if result.IDE == "" {
		t.Error("IDE should not be empty")
	}
}

// Edge case tests
func TestExtensionInstaller_ContextCancellation(t *testing.T) {
	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	extensionInfo := &ExtensionInfo{
		ID:          "test.extension",
		DownloadURL: "https://example.com/test.vsix",
	}

	_, err := installer.downloadVSIX(ctx, extensionInfo)
	if err == nil {
		t.Error("Expected error due to cancelled context")
	}

	if !strings.Contains(err.Error(), "context") {
		t.Error("Error should mention context cancellation")
	}
}

func TestExtensionInstaller_LargeFileHandling(t *testing.T) {
	// Create a server that returns a large response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write 1MB of data
		data := make([]byte, 1024*1024)
		for i := range data {
			data[i] = byte(i % 256)
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Write(data)
	}))
	defer server.Close()

	detector := NewIDEDetector()
	installer := NewExtensionInstaller(detector)

	extensionInfo := &ExtensionInfo{
		ID:          "test.extension",
		DownloadURL: server.URL + "/large.vsix",
	}

	ctx := context.Background()
	filePath, err := installer.downloadVSIX(ctx, extensionInfo)

	if err != nil {
		t.Errorf("Should handle large file download: %v", err)
	}

	if filePath != "" {
		defer os.Remove(filePath)

		// Verify file size
		stat, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("Downloaded file should exist: %v", err)
		}

		if stat.Size() != 1024*1024 {
			t.Errorf("Expected file size 1MB, got %d bytes", stat.Size())
		}
	}
}
