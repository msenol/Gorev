package gorev

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
)

func TestExportData(t *testing.T) {
	// Test için geçici database
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer vy.Kapat()

	iy := YeniIsYonetici(vy)

	// Test verisi oluştur
	setupTestData(t, vy)

	tests := []struct {
		name     string
		options  ExportOptions
		wantErr  bool
		validate func(*testing.T, *ExportFormat)
	}{
		{
			name: "Complete export with all data",
			options: ExportOptions{
				OutputPath:          "/tmp/test.json",
				IncludeCompleted:    true,
				IncludeDependencies: true,
				IncludeTemplates:    true,
				IncludeAIContext:    true,
			},
			wantErr: false,
			validate: func(t *testing.T, data *ExportFormat) {
				if data.Version == "" {
					t.Error("Version should not be empty")
				}
				if data.Metadata.ExportDate.IsZero() {
					t.Error("Export date should be set")
				}
				if len(data.Projects) == 0 {
					t.Error("Should export projects")
				}
				if len(data.Tasks) == 0 {
					t.Error("Should export tasks")
				}
			},
		},
		{
			name: "Export without completed tasks",
			options: ExportOptions{
				OutputPath:          "/tmp/test.json",
				IncludeCompleted:    false,
				IncludeDependencies: true,
				IncludeTemplates:    false,
				IncludeAIContext:    false,
			},
			wantErr: false,
			validate: func(t *testing.T, data *ExportFormat) {
				// Tamamlanmış görevlerin olmaması kontrol et
				for _, task := range data.Tasks {
					if task.Durum == constants.TaskStatusCompleted {
						t.Errorf("Found completed task when IncludeCompleted=false: %s", task.ID)
					}
				}
			},
		},
		{
			name: "Export with project filter",
			options: ExportOptions{
				OutputPath:          "/tmp/test.json",
				IncludeCompleted:    true,
				IncludeDependencies: true,
				ProjectFilter:       []string{"test-project-1"},
			},
			wantErr: false,
			validate: func(t *testing.T, data *ExportFormat) {
				// Sadece belirli projedeki görevlerin export edilmesi
				for _, task := range data.Tasks {
					if task.ProjeID != "" && task.ProjeID != "test-project-1" {
						t.Errorf("Found task from unexpected project: %s", task.ProjeID)
					}
				}
			},
		},
		{
			name: "Export with date range",
			options: ExportOptions{
				OutputPath:       "/tmp/test.json",
				IncludeCompleted: true,
				DateRange: &DateRange{
					From: &[]time.Time{time.Now().AddDate(0, 0, -7)}[0], // Son 7 gün
					To:   &[]time.Time{time.Now()}[0],
				},
			},
			wantErr: false,
			validate: func(t *testing.T, data *ExportFormat) {
				// Tarih aralığı kontrolü
				for _, task := range data.Tasks {
					if task.OlusturmaTarih.Before(time.Now().AddDate(0, 0, -7)) {
						t.Errorf("Found task outside date range: %s", task.ID)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := iy.ExportData(tt.options)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExportData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, data)
			}
		})
	}
}

func TestSaveExportToFile(t *testing.T) {
	// Test için geçici database
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer vy.Kapat()

	iy := YeniIsYonetici(vy)

	// Test verisi oluştur
	setupTestData(t, vy)

	// Export data oluştur
	data, err := iy.ExportData(ExportOptions{
		OutputPath:          "/tmp/export_for_save_test.json",
		IncludeCompleted:    true,
		IncludeDependencies: true,
	})
	if err != nil {
		t.Fatalf("Failed to create export data: %v", err)
	}

	tests := []struct {
		name       string
		outputPath string
		format     string
		wantErr    bool
		validate   func(*testing.T, string)
	}{
		{
			name:       "Save as JSON",
			outputPath: filepath.Join(t.TempDir(), "test_export.json"),
			format:     "json",
			wantErr:    false,
			validate: func(t *testing.T, path string) {
				// JSON dosyasının geçerli olduğunu kontrol et
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read exported file: %v", err)
				}

				var exportData ExportFormat
				if err := json.Unmarshal(content, &exportData); err != nil {
					t.Fatalf("Invalid JSON format: %v", err)
				}

				if exportData.Version == "" {
					t.Error("Version should not be empty in exported file")
				}
			},
		},
		{
			name:       "Save as CSV",
			outputPath: filepath.Join(t.TempDir(), "test_export.csv"),
			format:     "csv",
			wantErr:    false,
			validate: func(t *testing.T, path string) {
				// CSV dosyasının var olduğunu kontrol et
				if _, err := os.Stat(path); os.IsNotExist(err) {
					t.Fatalf("CSV file was not created: %s", path)
				}

				// CSV içeriğinin boş olmadığını kontrol et
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read CSV file: %v", err)
				}

				if len(content) == 0 {
					t.Error("CSV file is empty")
				}
			},
		},
		{
			name:       "Invalid format",
			outputPath: filepath.Join(t.TempDir(), "test_export.xml"),
			format:     "xml",
			wantErr:    true,
		},
		{
			name:       "Invalid path",
			outputPath: "/nonexistent/directory/test.json",
			format:     "json",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := ExportOptions{
				OutputPath: tt.outputPath,
				Format:     tt.format,
			}
			err := iy.SaveExportToFile(data, options)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveExportToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, tt.outputPath)
			}
		})
	}
}

func TestImportData(t *testing.T) {
	// Test için geçici database
	vy, err := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer vy.Kapat()

	iy := YeniIsYonetici(vy)

	// Export data oluştur (import testi için)
	setupTestData(t, vy)
	exportData, err := iy.ExportData(ExportOptions{
		OutputPath:          "/tmp/export_for_import.json",
		IncludeCompleted:    true,
		IncludeDependencies: true,
	})
	if err != nil {
		t.Fatalf("Failed to create export data: %v", err)
	}

	// Test için yeni database (temiz import testi)
	vy2, err2 := YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	if err2 != nil {
		t.Fatalf("Failed to create target database: %v", err2)
	}
	defer vy2.Kapat()

	iy2 := YeniIsYonetici(vy2)

	tests := []struct {
		name     string
		data     *ExportFormat
		options  ImportOptions
		wantErr  bool
		validate func(*testing.T, *ImportResult)
	}{
		{
			name: "Import into empty database",
			data: exportData,
			options: ImportOptions{
				FilePath:           filepath.Join(t.TempDir(), "test_import.json"),
				ImportMode:         "merge",
				ConflictResolution: "skip",
				PreserveIDs:        true,
				DryRun:             false,
			},
			wantErr: false,
			validate: func(t *testing.T, result *ImportResult) {
				if result.ImportedTasks == 0 {
					t.Error("Should import tasks")
				}
				if result.ImportedProjects == 0 {
					t.Error("Should import projects")
				}
				if len(result.Errors) > 0 {
					t.Errorf("Should not have errors: %v", result.Errors)
				}
			},
		},
		{
			name: "Dry run import",
			data: exportData,
			options: ImportOptions{
				FilePath:           filepath.Join(t.TempDir(), "test_import_dry.json"),
				ImportMode:         "merge",
				ConflictResolution: "skip",
				PreserveIDs:        true,
				DryRun:             true,
			},
			wantErr: false,
			validate: func(t *testing.T, result *ImportResult) {
				// Dry run should not actually import data
				if result.ImportedTasks > 0 {
					t.Error("Dry run should not import tasks")
				}
				if result.ImportedProjects > 0 {
					t.Error("Dry run should not import projects")
				}
			},
		},
		{
			name: "Import with conflict resolution overwrite",
			data: exportData,
			options: ImportOptions{
				FilePath:           filepath.Join(t.TempDir(), "test_import_overwrite.json"),
				ImportMode:         "merge",
				ConflictResolution: "overwrite",
				PreserveIDs:        true,
				DryRun:             false,
			},
			wantErr: false,
			validate: func(t *testing.T, result *ImportResult) {
				// Conflicts should be resolved by overwriting
				if len(result.Conflicts) > 0 && result.ImportedTasks == 0 {
					t.Error("Should resolve conflicts and import data")
				}
			},
		},
		{
			name: "Invalid import data",
			data: &ExportFormat{
				Version: "invalid",
				Metadata: ExportMetadata{
					ExportDate: time.Now(),
				},
				Projects: []*Proje{},
				Tasks:    []*Gorev{},
			},
			options: ImportOptions{
				FilePath:           filepath.Join(t.TempDir(), "test_import_invalid.json"),
				ImportMode:         "merge",
				ConflictResolution: "skip",
				DryRun:             false,
			},
			wantErr: false, // Should handle gracefully
			validate: func(t *testing.T, result *ImportResult) {
				// Should complete without errors but import nothing
				if result.ImportedTasks != 0 || result.ImportedProjects != 0 {
					t.Error("Should not import invalid data")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Export file'ı oluştur
			err := iy.SaveExportToFile(tt.data, ExportOptions{
				OutputPath: tt.options.FilePath,
				Format:     "json",
			})
			if err != nil {
				t.Fatalf("Failed to save export file: %v", err)
			}

			result, err := iy2.ImportData(tt.options)

			if (err != nil) != tt.wantErr {
				t.Errorf("ImportData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// CSV test removed - convertToCSV is not a public method

func TestValidateExportOptions(t *testing.T) {
	iy := &IsYonetici{}

	tests := []struct {
		name    string
		options ExportOptions
		wantErr bool
	}{
		{
			name: "Valid options",
			options: ExportOptions{
				OutputPath:          "/tmp/test.json",
				IncludeCompleted:    true,
				IncludeDependencies: true,
				IncludeTemplates:    false,
				IncludeAIContext:    false,
			},
			wantErr: false,
		},
		{
			name: "Valid date range",
			options: ExportOptions{
				OutputPath:       "/tmp/test.json",
				IncludeCompleted: true,
				DateRange: &DateRange{
					From: &[]time.Time{time.Now().AddDate(0, 0, -7)}[0],
					To:   &[]time.Time{time.Now()}[0],
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid date range - From after To",
			options: ExportOptions{
				OutputPath:       "/tmp/test.json",
				IncludeCompleted: true,
				DateRange: &DateRange{
					From: &[]time.Time{time.Now()}[0],
					To:   &[]time.Time{time.Now().AddDate(0, 0, -7)}[0],
				},
			},
			wantErr: true,
		},
		{
			name: "Empty project filter",
			options: ExportOptions{
				OutputPath:       "/tmp/test.json",
				IncludeCompleted: true,
				ProjectFilter:    []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := iy.validateExportFileOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateExportFileOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImportOptions(t *testing.T) {
	iy := &IsYonetici{}

	// Test için geçici dosya oluştur
	tempFile := filepath.Join(t.TempDir(), "test.json")
	err := os.WriteFile(tempFile, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	tests := []struct {
		name    string
		options ImportOptions
		wantErr bool
	}{
		{
			name: "Valid options",
			options: ImportOptions{
				FilePath:           tempFile,
				ImportMode:         "merge",
				ConflictResolution: "skip",
				PreserveIDs:        true,
				DryRun:             false,
			},
			wantErr: false,
		},
		{
			name: "Invalid import mode",
			options: ImportOptions{
				ImportMode:         "invalid",
				ConflictResolution: "skip",
			},
			wantErr: true,
		},
		{
			name: "Invalid conflict resolution",
			options: ImportOptions{
				ImportMode:         "merge",
				ConflictResolution: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Valid project mapping",
			options: ImportOptions{
				FilePath:           tempFile,
				ImportMode:         "merge",
				ConflictResolution: "skip",
				ProjectMapping: map[string]string{
					"old-id": "new-id",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := iy.validateImportOptions(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateImportOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test helper functions

func setupTestData(t *testing.T, vy *VeriYonetici) {
	// Test projesi oluştur
	proje := &Proje{
		ID:              "test-project-1",
		Isim:            "Test Project",
		Tanim:           "Test project description",
		OlusturmaTarih:  time.Now(),
		GuncellemeTarih: time.Now(),
	}
	err := vy.ProjeKaydet(proje)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	// Test görevleri oluştur
	tasks := []*Gorev{
		{
			ID:              "test-task-1",
			Baslik:          "Test Task 1",
			Aciklama:        "Test task description",
			Durum:           constants.TaskStatusPending,
			Oncelik:         constants.PriorityHigh,
			ProjeID:         proje.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
		{
			ID:              "test-task-2",
			Baslik:          "Test Task 2",
			Aciklama:        "Another test task",
			Durum:           constants.TaskStatusCompleted,
			Oncelik:         constants.PriorityMedium,
			ProjeID:         proje.ID,
			OlusturmaTarih:  time.Now(),
			GuncellemeTarih: time.Now(),
		},
	}

	for _, task := range tasks {
		err := vy.GorevKaydet(task)
		if err != nil {
			t.Fatalf("Failed to create test task: %v", err)
		}
	}

	// Etiket test için şimdilik atlıyoruz - EtiketKaydet metodu yok
}
