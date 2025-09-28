package gorev

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewIDEDetector(t *testing.T) {
	detector := NewIDEDetector()
	if detector == nil {
		t.Fatal("NewIDEDetector returned nil")
	}

	if detector.detectedIDEs == nil {
		t.Error("detectedIDEs map should be initialized")
	}
}

func TestIDETypes(t *testing.T) {
	testCases := []struct {
		ideType  IDEType
		expected string
	}{
		{IDETypeVSCode, "vscode"},
		{IDETypeCursor, "cursor"},
		{IDETypeWindsurf, "windsurf"},
	}

	for _, tc := range testCases {
		if string(tc.ideType) != tc.expected {
			t.Errorf("IDE type %v should equal %s, got %s", tc.ideType, tc.expected, string(tc.ideType))
		}
	}
}

func TestIDEInfo(t *testing.T) {
	ide := &IDEInfo{
		Type:           IDETypeVSCode,
		Name:           "Visual Studio Code",
		ExecutablePath: "/usr/bin/code",
		ConfigPath:     "/home/user/.config/Code",
		ExtensionsPath: "/home/user/.vscode/extensions",
		Version:        "1.80.0",
		IsInstalled:    true,
	}

	if ide.Type != IDETypeVSCode {
		t.Error("IDE type mismatch")
	}

	if ide.Name != "Visual Studio Code" {
		t.Error("IDE name mismatch")
	}

	if !ide.IsInstalled {
		t.Error("IDE should be marked as installed")
	}
}

func TestFileExists(t *testing.T) {
	// Test with non-existent file
	if fileExists("/nonexistent/file/path") {
		t.Error("fileExists should return false for non-existent file")
	}

	// Test with empty path
	if fileExists("") {
		t.Error("fileExists should return false for empty path")
	}
}

func TestDirExists(t *testing.T) {
	// Test with non-existent directory
	if dirExists("/nonexistent/directory/path") {
		t.Error("dirExists should return false for non-existent directory")
	}

	// Test with empty path
	if dirExists("") {
		t.Error("dirExists should return false for empty path")
	}
}

func TestGetDetectedIDE(t *testing.T) {
	detector := NewIDEDetector()

	// Test with non-existent IDE
	ide, exists := detector.GetDetectedIDE(IDETypeVSCode)
	if exists {
		t.Error("GetDetectedIDE should return false for non-detected IDE")
	}
	if ide != nil {
		t.Error("IDE should be nil for non-detected IDE")
	}

	// Add a mock IDE
	mockIDE := &IDEInfo{
		Type:        IDETypeVSCode,
		Name:        "Mock VS Code",
		IsInstalled: true,
	}
	detector.detectedIDEs[IDETypeVSCode] = mockIDE

	// Test with existing IDE
	ide, exists = detector.GetDetectedIDE(IDETypeVSCode)
	if !exists {
		t.Error("GetDetectedIDE should return true for detected IDE")
	}
	if ide == nil {
		t.Fatal("IDE should not be nil for detected IDE")
	}
	if ide.Name != "Mock VS Code" {
		t.Error("IDE name mismatch")
	}
}

func TestGetAllDetectedIDEs(t *testing.T) {
	detector := NewIDEDetector()

	// Test with empty detector
	ides := detector.GetAllDetectedIDEs()
	if len(ides) != 0 {
		t.Error("GetAllDetectedIDEs should return empty map for new detector")
	}

	// Add mock IDEs
	detector.detectedIDEs[IDETypeVSCode] = &IDEInfo{Type: IDETypeVSCode, Name: "VS Code"}
	detector.detectedIDEs[IDETypeCursor] = &IDEInfo{Type: IDETypeCursor, Name: "Cursor"}

	ides = detector.GetAllDetectedIDEs()
	if len(ides) != 2 {
		t.Errorf("Expected 2 detected IDEs, got %d", len(ides))
	}

	if _, exists := ides[IDETypeVSCode]; !exists {
		t.Error("VS Code should be in detected IDEs")
	}

	if _, exists := ides[IDETypeCursor]; !exists {
		t.Error("Cursor should be in detected IDEs")
	}
}

// BenchmarkIDEDetection benchmarks the IDE detection process
func BenchmarkIDEDetection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		detector := NewIDEDetector()
		// This will typically fail in test environment but that's ok for benchmark
		detector.DetectAllIDEs()
	}
}

// TestIDEDetectionError tests error handling in IDE detection
func TestIDEDetectionError(t *testing.T) {
	detector := NewIDEDetector()

	// Test detection with no IDEs installed (expected in test environment)
	detectedIDEs, err := detector.DetectAllIDEs()

	// This should not error, just return empty results
	if err != nil {
		t.Errorf("DetectAllIDEs should not error in test environment: %v", err)
	}

	// In test environment, no IDEs will be detected
	if len(detectedIDEs) > 0 {
		t.Logf("Note: %d IDEs detected in test environment", len(detectedIDEs))
		// This is actually good - means the detection is working!
		for ideType, ide := range detectedIDEs {
			t.Logf("Detected: %s - %s at %s", ideType, ide.Name, ide.ExecutablePath)
		}
	}
}

// TestDetectIDEUnsupportedType tests error handling for unsupported IDE types
func TestDetectIDEUnsupportedType(t *testing.T) {
	detector := NewIDEDetector()

	// Test with invalid IDE type
	_, err := detector.detectIDE("unsupported")
	if err == nil {
		t.Error("Expected error for unsupported IDE type")
	}
}

// TestGetDetectedIDECoverage improves coverage for GetDetectedIDE
func TestGetDetectedIDECoverage(t *testing.T) {
	detector := NewIDEDetector()

	// Test with empty detector
	ide, exists := detector.GetDetectedIDE(IDETypeVSCode)
	if exists {
		t.Error("GetDetectedIDE should return false for empty detector")
	}
	if ide != nil {
		t.Error("GetDetectedIDE should return nil for empty detector")
	}

	// Test with multiple IDE types
	for _, ideType := range []IDEType{IDETypeVSCode, IDETypeCursor, IDETypeWindsurf} {
		ide, exists = detector.GetDetectedIDE(ideType)
		if exists {
			t.Errorf("GetDetectedIDE should return false for %s in empty detector", ideType)
		}
		if ide != nil {
			t.Errorf("GetDetectedIDE should return nil for %s in empty detector", ideType)
		}
	}
}

// TestGetAllDetectedIDEsCoverage improves coverage for GetAllDetectedIDEs
func TestGetAllDetectedIDEsCoverage(t *testing.T) {
	detector := NewIDEDetector()

	// Test with empty detector (already covered in existing test)
	ides := detector.GetAllDetectedIDEs()
	if len(ides) != 0 {
		t.Error("GetAllDetectedIDEs should return empty map for new detector")
	}

	// Test with single IDE
	detector.detectedIDEs[IDETypeVSCode] = &IDEInfo{Type: IDETypeVSCode, Name: "VS Code"}
	ides = detector.GetAllDetectedIDEs()
	if len(ides) != 1 {
		t.Errorf("Expected 1 detected IDE, got %d", len(ides))
	}

	// Test that returned map is a copy (modifying it shouldn't affect detector)
	ides[IDETypeCursor] = &IDEInfo{Type: IDETypeCursor, Name: "Cursor"}
	if _, exists := detector.detectedIDEs[IDETypeCursor]; exists {
		t.Error("Modifying returned map should not affect detector")
	}
}

// TestIsExtensionInstalled tests extension installation detection
func TestIsExtensionInstalled(t *testing.T) {
	detector := NewIDEDetector()

	// Test with no IDEs detected
	installed, err := detector.IsExtensionInstalled(IDETypeVSCode, "mehmetsenol.gorev-vscode")
	if err == nil {
		t.Error("Expected error when IDE is not detected")
	}
	if installed {
		t.Error("Should return false when IDE is not detected")
	}

	// Test with IDE detected but no extensions path
	detector.detectedIDEs[IDETypeVSCode] = &IDEInfo{
		Type:        IDETypeVSCode,
		Name:        "VS Code",
		IsInstalled: true,
	}
	installed, err = detector.IsExtensionInstalled(IDETypeVSCode, "mehmetsenol.gorev-vscode")
	if err == nil {
		t.Error("Expected error when extensions path is empty")
	}
	if installed {
		t.Error("Should return false when extensions path is empty")
	}

	// Test with valid IDE info but non-existent extensions path
	detector.detectedIDEs[IDETypeVSCode].ExtensionsPath = "/nonexistent/extensions/path"
	installed, err = detector.IsExtensionInstalled(IDETypeVSCode, "mehmetsenol.gorev-vscode")
	if err != nil {
		t.Logf("Expected error for non-existent extensions path: %v", err)
	}
	if installed {
		t.Error("Should return false for non-existent extensions path")
	}
}

// TestGetExtensionVersion tests extension version retrieval
func TestGetExtensionVersion(t *testing.T) {
	detector := NewIDEDetector()

	// Test with no IDEs detected
	version, err := detector.GetExtensionVersion(IDETypeVSCode, "mehmetsenol.gorev-vscode")
	if err == nil {
		t.Error("Expected error when IDE is not detected")
	}
	if version != "" {
		t.Error("Should return empty version when IDE is not detected")
	}

	// Test with IDE detected but no extensions path
	detector.detectedIDEs[IDETypeVSCode] = &IDEInfo{
		Type:        IDETypeVSCode,
		Name:        "VS Code",
		IsInstalled: true,
	}
	version, err = detector.GetExtensionVersion(IDETypeVSCode, "mehmetsenol.gorev-vscode")
	if err == nil {
		t.Error("Expected error when extensions path is empty")
	}
	if version != "" {
		t.Error("Should return empty version when extensions path is empty")
	}

	// Test with valid IDE info but non-existent extensions path
	detector.detectedIDEs[IDETypeVSCode].ExtensionsPath = "/nonexistent/extensions/path"
	version, err = detector.GetExtensionVersion(IDETypeVSCode, "mehmetsenol.gorev-vscode")
	if err == nil {
		t.Log("Note: No error for non-existent extensions path (this is expected)")
	}
	if version != "" {
		t.Error("Should return empty version for non-existent extensions path")
	}
}

// TestGetIDEVersion tests the getIDEVersion method
func TestGetIDEVersion(t *testing.T) {
	detector := NewIDEDetector()

	// Test with non-existent executable
	version := detector.getIDEVersion("/nonexistent/path", []string{"--version"})
	if version != "unknown" {
		t.Errorf("Expected 'unknown' version for non-existent executable, got '%s'", version)
	}

	// Test with empty path
	version = detector.getIDEVersion("", []string{"--version"})
	if version != "unknown" {
		t.Errorf("Expected 'unknown' version for empty path, got '%s'", version)
	}

	// Test with different args
	version = detector.getIDEVersion("/nonexistent/path", []string{"version"})
	if version != "unknown" {
		t.Errorf("Expected 'unknown' version with different args, got '%s'", version)
	}
}

// TestFileExistsWithVariousInputs tests fileExists with different inputs
func TestFileExistsWithVariousInputs(t *testing.T) {
	// Test with directory path (should return false)
	tempDir := t.TempDir()
	if fileExists(tempDir) {
		t.Error("fileExists should return false for directory path")
	}

	// Test with non-existent file
	if fileExists(filepath.Join(tempDir, "nonexistent.txt")) {
		t.Error("fileExists should return false for non-existent file")
	}

	// Test with existing file
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if !fileExists(testFile) {
		t.Error("fileExists should return true for existing file")
	}
}

// TestDirExistsWithVariousInputs tests dirExists with different inputs
func TestDirExistsWithVariousInputs(t *testing.T) {
	// Test with file path (should return false)
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)
	if dirExists(testFile) {
		t.Error("dirExists should return false for file path")
	}

	// Test with non-existent directory
	if dirExists(filepath.Join(tempDir, "nonexistent")) {
		t.Error("dirExists should return false for non-existent directory")
	}

	// Test with existing directory
	if !dirExists(tempDir) {
		t.Error("dirExists should return true for existing directory")
	}
}

// TestDetectAllIDEsThreadSafety tests thread safety of DetectAllIDEs
func TestDetectAllIDEsThreadSafety(t *testing.T) {
	detector := NewIDEDetector()

	// Call DetectAllIDEs concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := detector.DetectAllIDEs()
			if err != nil {
				t.Logf("Concurrent detection error: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestDetectAllIDEsReturnsCopy tests that DetectAllIDEs returns a copy of the map
func TestDetectAllIDEsReturnsCopy(t *testing.T) {
	detector := NewIDEDetector()

	// Get initial results
	ides1, err := detector.DetectAllIDEs()
	if err != nil {
		t.Fatalf("DetectAllIDEs failed: %v", err)
	}

	// Get results again
	ides2, err := detector.DetectAllIDEs()
	if err != nil {
		t.Fatalf("DetectAllIDEs failed: %v", err)
	}

	// Modifying the first map should not affect the second
	ides1["test"] = &IDEInfo{Type: "test", Name: "Test"}
	if _, exists := ides2["test"]; exists {
		t.Error("Modifying returned map should not affect detector")
	}
}
