package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	mcphandlers "github.com/msenol/gorev/internal/mcp"
	testinghelpers "github.com/msenol/gorev/internal/testing"
)

// setupTestDB creates a test database and returns VeriYonetici (legacy helper - use testinghelpers instead)
func setupTestDB(t *testing.T) *gorev.VeriYonetici {
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPathIntegration,
		CreateTemplates: true,
		InitializeI18n:  true,
	}
	veriYonetici, _ := testinghelpers.SetupTestDatabase(t, config)
	return veriYonetici
}

func TestGorevExportMCPTool(t *testing.T) {
	// Test database setup
	vy := setupTestDB(t)
	defer vy.Kapat()

	iy := gorev.YeniIsYonetici(vy)
	handlers := mcphandlers.YeniHandlers(iy)

	// Test verisi oluştur
	setupExportTestData(t, vy)

	// Create a shared temp directory for all tests in this function
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(*testing.T, *mcp.CallToolResult, map[string]interface{})
	}{
		{
			name: "Basic JSON export",
			params: map[string]interface{}{
				"output_path": filepath.Join(tempDir, "basic_export.json"),
				"format":      "json",
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}

				// Dosyanın oluşturulduğunu kontrol et
				outputPath := params["output_path"].(string)
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Error("Export file was not created")
				}
			},
		},
		{
			name: "CSV export with filters",
			params: map[string]interface{}{
				"output_path":       filepath.Join(tempDir, "filtered_export.csv"),
				"format":            "csv",
				"include_completed": false,
				"project_filter":    []string{"test-project-1"},
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}
			},
		},
		{
			name: "Export with date range",
			params: map[string]interface{}{
				"output_path": filepath.Join(tempDir, "date_range_export.json"),
				"format":      "json",
				"date_range": map[string]interface{}{
					"from": time.Now().AddDate(0, 0, -7).Format(time.RFC3339),
					"to":   time.Now().Format(time.RFC3339),
				},
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}
			},
		},
		{
			name: "Export with all options enabled",
			params: map[string]interface{}{
				"output_path":          filepath.Join(tempDir, "complete_export.json"),
				"format":               "json",
				"include_completed":    true,
				"include_dependencies": true,
				"include_templates":    true,
				"include_ai_context":   true,
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}
			},
		},
		{
			name: "Missing output_path parameter",
			params: map[string]interface{}{
				"format": "json",
			},
			wantErr: false, // MCP handlers return CallToolResult, not Go errors
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if !result.IsError {
					t.Error("Expected error for missing output_path")
				}
			},
		},
		{
			name: "Invalid format",
			params: map[string]interface{}{
				"output_path": filepath.Join(tempDir, "invalid_format.xml"),
				"format":      "xml",
			},
			wantErr: false, // MCP handlers return CallToolResult, not Go errors
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if !result.IsError {
					t.Error("Expected error for invalid format")
				}
			},
		},
		{
			name: "Invalid directory path",
			params: map[string]interface{}{
				"output_path": "/nonexistent/directory/export.json",
				"format":      "json",
			},
			wantErr: false, // MCP handlers return CallToolResult, not Go errors
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if !result.IsError {
					t.Error("Expected error for invalid directory")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.GorevExport(tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("GorevExport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result, tt.params)
			}
		})
	}
}

func TestGorevImportMCPTool(t *testing.T) {
	// Test database setup
	vy := setupTestDB(t)
	defer vy.Kapat()
	iy := gorev.YeniIsYonetici(vy)
	handlers := mcphandlers.YeniHandlers(iy)

	// Test export dosyası oluştur
	exportFile := createTestExportFile(t, vy)

	tests := []struct {
		name     string
		params   map[string]interface{}
		wantErr  bool
		validate func(*testing.T, *mcp.CallToolResult, map[string]interface{})
	}{
		{
			name: "Basic import",
			params: map[string]interface{}{
				"file_path":           exportFile,
				"import_mode":         "merge",
				"conflict_resolution": "skip",
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}
			},
		},
		{
			name: "Dry run import",
			params: map[string]interface{}{
				"file_path":           exportFile,
				"import_mode":         "merge",
				"conflict_resolution": "skip",
				"dry_run":             true,
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}
				// Dry run sonuçlarını kontrol et
				if len(result.Content) > 0 {
					content := result.Content[0]
					if textContent, ok := content.(map[string]interface{}); ok {
						if text, exists := textContent["text"].(string); exists {
							if len(text) == 0 {
								t.Error("Dry run should return analysis results")
							}
						}
					}
				}
			},
		},
		{
			name: "Import with overwrite conflict resolution",
			params: map[string]interface{}{
				"file_path":           exportFile,
				"import_mode":         "merge",
				"conflict_resolution": "overwrite",
				"preserve_ids":        true,
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}
			},
		},
		{
			name: "Import with project mapping",
			params: map[string]interface{}{
				"file_path":           exportFile,
				"import_mode":         "merge",
				"conflict_resolution": "skip",
				"project_mapping": map[string]interface{}{
					"old-project-id": "new-project-id",
				},
			},
			wantErr: false,
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if result.IsError {
					t.Errorf("Expected success, got error: %v", result.Content)
				}
			},
		},
		{
			name: "Missing file_path parameter",
			params: map[string]interface{}{
				"import_mode": "merge",
			},
			wantErr: false, // MCP handlers return CallToolResult, not Go errors
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if !result.IsError {
					t.Error("Expected error for missing file_path")
				}
			},
		},
		{
			name: "Non-existent file",
			params: map[string]interface{}{
				"file_path":           "/nonexistent/file.json",
				"import_mode":         "merge",
				"conflict_resolution": "skip",
			},
			wantErr: false, // MCP handlers return CallToolResult, not Go errors
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if !result.IsError {
					t.Error("Expected error for non-existent file")
				}
			},
		},
		{
			name: "Invalid import mode",
			params: map[string]interface{}{
				"file_path":           exportFile,
				"import_mode":         "invalid",
				"conflict_resolution": "skip",
			},
			wantErr: false, // MCP handlers return CallToolResult, not Go errors
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if !result.IsError {
					t.Error("Expected error for invalid import mode")
				}
			},
		},
		{
			name: "Invalid conflict resolution",
			params: map[string]interface{}{
				"file_path":           exportFile,
				"import_mode":         "merge",
				"conflict_resolution": "invalid",
			},
			wantErr: false, // MCP handlers return CallToolResult, not Go errors
			validate: func(t *testing.T, result *mcp.CallToolResult, params map[string]interface{}) {
				if !result.IsError {
					t.Error("Expected error for invalid conflict resolution")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.GorevImport(tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("GorevImport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.validate != nil {
				tt.validate(t, result, tt.params)
			}
		})
	}
}

func TestExportImportRoundTrip(t *testing.T) {
	// Kaynak database
	sourceVY := setupTestDB(t)
	defer func() { _ = sourceVY.Kapat() }()

	sourceIY := gorev.YeniIsYonetici(sourceVY)
	sourceHandlers := mcphandlers.YeniHandlers(sourceIY)

	// Test verisi oluştur
	setupExportTestData(t, sourceVY)

	// Export
	exportPath := filepath.Join(t.TempDir(), "roundtrip_export.json")
	exportParams := map[string]interface{}{
		"output_path":          exportPath,
		"format":               "json",
		"include_completed":    true,
		"include_dependencies": true,
		"include_templates":    true,
	}

	exportResult, err := sourceHandlers.GorevExport(exportParams)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if exportResult.IsError {
		t.Fatalf("Export returned error: %v", exportResult.Content)
	}

	// Hedef database
	targetVY := setupTestDB(t)
	defer func() { _ = targetVY.Kapat() }()
	targetIY := gorev.YeniIsYonetici(targetVY)
	targetHandlers := mcphandlers.YeniHandlers(targetIY)

	// Import
	importParams := map[string]interface{}{
		"file_path":           exportPath,
		"import_mode":         "merge",
		"conflict_resolution": "overwrite",
		"preserve_ids":        true,
	}

	importResult, err := targetHandlers.GorevImport(importParams)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}
	if importResult.IsError {
		t.Fatalf("Import returned error: %v", importResult.Content)
	}

	// Veri karşılaştırması
	sourceGorevler, err := sourceVY.GorevListele(map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to list source tasks: %v", err)
	}

	targetGorevler, err := targetVY.GorevListele(map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to list target tasks: %v", err)
	}

	if len(sourceGorevler) != len(targetGorevler) {
		t.Errorf("Task count mismatch: source=%d, target=%d",
			len(sourceGorevler), len(targetGorevler))
	}

	// Proje sayısını kontrol et
	sourceProjeler, err := sourceIY.ProjeListele()
	if err != nil {
		t.Fatalf("Failed to list source projects: %v", err)
	}

	targetProjeler, err := targetIY.ProjeListele()
	if err != nil {
		t.Fatalf("Failed to list target projects: %v", err)
	}

	if len(sourceProjeler) != len(targetProjeler) {
		t.Errorf("Project count mismatch: source=%d, target=%d",
			len(sourceProjeler), len(targetProjeler))
	}
}

func TestExportImportWithLargeDataset(t *testing.T) {
	// Büyük veri seti testi
	vy := setupTestDB(t)
	defer vy.Kapat()
	iy := gorev.YeniIsYonetici(vy)
	handlers := mcphandlers.YeniHandlers(iy)

	// Büyük test verisi oluştur (100 proje, 1000 görev)
	setupLargeTestDataset(t, vy, 100, 1000)

	// Export
	exportPath := filepath.Join(t.TempDir(), "large_export.json")
	exportParams := map[string]interface{}{
		"output_path": exportPath,
		"format":      "json",
	}

	start := time.Now()
	exportResult, err := handlers.GorevExport(exportParams)
	exportDuration := time.Since(start)

	if err != nil {
		t.Fatalf("Large export failed: %v", err)
	}
	if exportResult.IsError {
		t.Fatalf("Large export returned error: %v", exportResult.Content)
	}

	t.Logf("Large export took %v", exportDuration)

	// Dosya boyutunu kontrol et
	fileInfo, err := os.Stat(exportPath)
	if err != nil {
		t.Fatalf("Failed to stat export file: %v", err)
	}
	t.Logf("Export file size: %d bytes", fileInfo.Size())

	// Import test için yeni database
	targetVY := setupTestDB(t)
	defer func() { _ = targetVY.Kapat() }()

	targetIY := gorev.YeniIsYonetici(targetVY)
	targetHandlers := mcphandlers.YeniHandlers(targetIY)

	// Import
	importParams := map[string]interface{}{
		"file_path":           exportPath,
		"import_mode":         "merge",
		"conflict_resolution": "skip",
	}

	start = time.Now()
	importResult, err := targetHandlers.GorevImport(importParams)
	importDuration := time.Since(start)

	if err != nil {
		t.Fatalf("Large import failed: %v", err)
	}
	if importResult.IsError {
		t.Fatalf("Large import returned error: %v", importResult.Content)
	}

	t.Logf("Large import took %v", importDuration)

	// Performans uyarıları
	if exportDuration > 30*time.Second {
		t.Logf("WARNING: Export took longer than 30 seconds: %v", exportDuration)
	}
	if importDuration > 60*time.Second {
		t.Logf("WARNING: Import took longer than 60 seconds: %v", importDuration)
	}
}

// Test helper functions

func setupExportTestData(t *testing.T, vy *gorev.VeriYonetici) {
	// Test projesi oluştur
	proje := &gorev.Proje{
		ID:         "test-project-1",
		Name:       "Test Export Project",
		Definition: "Project for export testing",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := vy.ProjeKaydet(proje)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Test görevleri oluştur
	tasks := []*gorev.Gorev{
		{
			ID:          "export-task-1",
			Title:       "Export Test Task 1",
			Description: "First test task for export",
			Status:      "beklemede",
			Priority:    "yuksek",
			ProjeID:     proje.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "export-task-2",
			Title:       "Export Test Task 2",
			Description: "Second test task for export",
			Status:      "tamamlandi",
			Priority:    "orta",
			ProjeID:     proje.ID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, task := range tasks {
		err := vy.GorevKaydet(task)
		if err != nil {
			t.Fatalf("Failed to create test task: %v", err)
		}
	}
}

func createTestExportFile(t *testing.T, vy *gorev.VeriYonetici) string {
	// Test verisi oluştur
	setupExportTestData(t, vy)

	// Export data oluştur
	iy := gorev.YeniIsYonetici(vy)
	exportData, err := iy.ExportData(gorev.ExportOptions{
		IncludeCompleted:    true,
		IncludeDependencies: true,
	})
	if err != nil {
		t.Fatalf("Failed to create export data: %v", err)
	}

	// Dosyaya kaydet
	exportPath := filepath.Join(t.TempDir(), "test_export.json")
	err = iy.SaveExportToFile(exportData, gorev.ExportOptions{
		OutputPath: exportPath,
		Format:     "json",
	})
	if err != nil {
		t.Fatalf("Failed to save export file: %v", err)
	}

	return exportPath
}

func setupLargeTestDataset(t *testing.T, vy *gorev.VeriYonetici, projectCount, taskCount int) {
	// Projeler oluştur
	for i := 0; i < projectCount; i++ {
		proje := &gorev.Proje{
			ID:         fmt.Sprintf("large-project-%d", i),
			Name:       fmt.Sprintf("Large Test Project %d", i),
			Definition: fmt.Sprintf("Project %d for large dataset testing", i),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		err := vy.ProjeKaydet(proje)
		if err != nil {
			t.Fatalf("Failed to create large test project %d: %v", i, err)
		}
	}

	// Görevler oluştur
	for i := 0; i < taskCount; i++ {
		projectID := fmt.Sprintf("large-project-%d", i%projectCount)
		task := &gorev.Gorev{
			ID:          fmt.Sprintf("large-task-%d", i),
			Title:       fmt.Sprintf("Large Test Task %d", i),
			Description: fmt.Sprintf("Task %d for large dataset testing", i),
			Status:      []string{"beklemede", "devam_ediyor", "tamamlandi"}[i%3],
			Priority:    []string{"dusuk", "orta", "yuksek"}[i%3],
			ProjeID:     projectID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err := vy.GorevKaydet(task)
		if err != nil {
			t.Fatalf("Failed to create large test task %d: %v", i, err)
		}
	}
}
