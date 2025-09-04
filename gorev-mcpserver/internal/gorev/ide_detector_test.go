package gorev

import (
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
