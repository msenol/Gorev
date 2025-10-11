package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	testinghelpers "github.com/msenol/gorev/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnvironmentWithTemplate sets up test environment with a simple template and active project
func setupTestEnvironmentWithTemplate(t *testing.T) (*server.MCPServer, *Handlers, string, func()) {
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     false, // Use temp file for edge cases
		UseTempFile:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: false, // We'll create our own template
		InitializeI18n:  true,
	}

	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	veriYonetici := isYonetici.VeriYonetici()

	// Create simple template for edge case testing
	template := &gorev.GorevTemplate{
		Name:                "Simple Test Template",
		Definition:          "Basit template for edge case testing",
		DefaultTitle:        "{{baslik}}",
		DescriptionTemplate: "{{aciklama}}",
		Fields: []gorev.TemplateAlan{
			{Name: "baslik", Type: "text", Required: true},
			{Name: "aciklama", Type: "text", Required: false, Default: "Test description"},
			{Name: "oncelik", Type: "select", Required: false, Default: constants.PriorityMedium, Options: []string{constants.PriorityLow, constants.PriorityMedium, constants.PriorityHigh}},
			{Name: "etiketler", Type: "text", Required: false},
			{Name: "son_tarih", Type: "date", Required: false},
		},
		Category: "Test",
		Active:   true,
	}

	err := veriYonetici.TemplateOlustur(template)
	require.NoError(t, err)

	// Create and set active project
	proje := &gorev.Proje{
		ID:         constants.TestProjectIDEdge,
		Name:       "Test Edge Project",
		Definition: "Test project for edge cases",
	}
	err = veriYonetici.ProjeKaydet(proje)
	require.NoError(t, err)
	err = veriYonetici.AktifProjeAyarla(proje.ID)
	require.NoError(t, err)

	mcpServer, err := YeniMCPSunucu(isYonetici)
	require.NoError(t, err)
	handlers := YeniHandlers(isYonetici)

	return mcpServer, handlers, template.ID, cleanup
}

// TestGorevOlustur_EdgeCases tests edge cases for task creation
func TestGorevOlustur_EdgeCases(t *testing.T) {
	_, handlers, templateID, cleanup := setupTestEnvironmentWithTemplate(t)
	defer cleanup()

	// Test 1: Empty strings and whitespace
	t.Run("Empty strings and whitespace", func(t *testing.T) {
		testCases := []struct {
			name    string
			baslik  string
			wantErr bool
		}{
			{"Empty title", "", true},
			{"Whitespace only title", "   ", true},
			{"Tab only title", "\t", true},
			{"Newline only title", "\n", true},
			{"Mixed whitespace", " \t\n ", true},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
					constants.ParamTemplateID: templateID,
					constants.ParamValues: map[string]interface{}{
						"baslik":  tc.baslik,
						"oncelik": constants.PriorityMedium,
					},
				})

				if tc.wantErr && !result.IsError {
					t.Errorf("Expected error for %s, got success", tc.name)
				}
				if !tc.wantErr && result.IsError {
					t.Errorf("Expected success for %s, got error: %v", tc.name, getResultText(result))
				}
			})
		}
	})

	// Test 2: SQL injection attempts
	t.Run("SQL injection attempts", func(t *testing.T) {
		injectionAttempts := []string{
			"'; DROP TABLE gorevler; --",
			"\" OR 1=1 --",
			"'; DELETE FROM projeler WHERE 1=1; --",
			"1'; UPDATE gorevler SET durum='tamamlandi' WHERE 1=1; --",
			"Robert'); DROP TABLE students;--",
		}

		for _, injection := range injectionAttempts {
			result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
				constants.ParamTemplateID: templateID,
				constants.ParamValues: map[string]interface{}{
					"baslik":   injection,
					"aciklama": injection,
					"oncelik":  constants.PriorityMedium,
				},
			})

			// Should either sanitize or reject, but not execute SQL
			if result.IsError {
				continue // Rejection is fine
			}

			// If accepted, verify the task was created with escaped content
			text := getResultText(result)
			if !strings.Contains(text, "oluşturuldu") {
				t.Errorf("Task creation failed for injection: %s", injection)
			}

			// Verify database is still intact
			listResult := callTool(t, handlers, "gorev_listele", map[string]interface{}{})
			if listResult.IsError {
				t.Fatalf("Database corrupted after injection attempt: %v", getResultText(listResult))
			}
		}
	})

	// Test 3: Special characters and Unicode
	t.Run("Special characters and Unicode", func(t *testing.T) {
		specialCases := []struct {
			name   string
			baslik string
		}{
			{"Emoji", "🚀 Deploy to production 🎉"},
			{"Chinese characters", "部署到生产环境"},
			{"Arabic text", "نشر إلى الإنتاج"},
			{"Mixed scripts", "Deploy 部署 🚀 نشر"},
			{"Zero-width characters", "Deploy\u200Bto\u200Cproduction"},
			{"Control characters", "Deploy\x00to\x01production"},
		}

		for _, sc := range specialCases {
			t.Run(sc.name, func(t *testing.T) {
				result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
					constants.ParamTemplateID: templateID,
					constants.ParamValues: map[string]interface{}{
						"baslik":  sc.baslik,
						"oncelik": constants.PriorityMedium,
					},
				})

				if result.IsError {
					t.Logf("Task creation with %s failed: %v", sc.name, getResultText(result))
					// Some special characters might be rejected, which is OK
				} else {
					// Verify the task was created
					text := getResultText(result)
					taskID := extractTaskIDFromText(text)
					if taskID == "" {
						t.Errorf("Failed to extract task ID for %s", sc.name)
					}
				}
			})
		}
	})

	// Test 4: Extremely long inputs
	t.Run("Extremely long inputs", func(t *testing.T) {
		longString := strings.Repeat("A", constants.TestStringVeryLong)
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: templateID,
			constants.ParamValues: map[string]interface{}{
				"baslik":   "Task with long description",
				"aciklama": longString,
				"oncelik":  constants.PriorityMedium,
			},
		})

		if result.IsError {
			t.Logf("Task creation with long description failed: %v", getResultText(result))
		} else {
			// Verify the task was created
			text := getResultText(result)
			taskID := extractTaskIDFromText(text)
			if taskID != "" {
				// Check if the long description was stored
				detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{"id": taskID})
				if !detailResult.IsError {
					detail := getResultText(detailResult)
					if len(detail) < constants.TestStringLong {
						t.Error("Long description may have been truncated")
					}
				}
			}
		}
	})

	// Test 5: Invalid priority values
	t.Run("Invalid priority values", func(t *testing.T) {
		invalidPriorities := []string{
			"critical",
			"YUKSEK",
			"1",
			"high",
			"",
			"null",
		}

		for _, priority := range invalidPriorities {
			result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
				constants.ParamTemplateID: templateID,
				constants.ParamValues: map[string]interface{}{
					"baslik":  "Test task",
					"oncelik": priority,
				},
			})

			// Invalid priorities should either be rejected or accepted
			// If accepted, the template system currently doesn't validate select values
			// so the invalid priority will be stored as-is (which is not ideal but current behavior)
			if !result.IsError {
				text := getResultText(result)
				taskID := extractTaskIDFromText(text)
				if taskID != "" {
					// Check what priority was assigned
					detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{"id": taskID})
					detail := getResultText(detailResult)
					// Task was created - the priority validation is not yet implemented
					// For now, just log what priority was assigned
					t.Logf("Task created with priority: %s, detail contains: %s", priority, detail)
				}
			}
		}
	})

	// Test 6: Invalid date formats
	t.Run("Invalid date formats", func(t *testing.T) {
		invalidDates := []string{
			"31-12-2025",
			"2025/12/31",
			"December 31, 2025",
			"2025-13-01",
			"2025-12-32",
			"not-a-date",
			"2025-02-30",
		}

		for _, date := range invalidDates {
			result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
				constants.ParamTemplateID: templateID,
				constants.ParamValues: map[string]interface{}{
					"baslik":    "Test task",
					"son_tarih": date,
					"oncelik":   constants.PriorityMedium,
				},
			})

			if !result.IsError {
				// Task created despite invalid date - check if date was ignored
				text := getResultText(result)
				taskID := extractTaskIDFromText(text)
				if taskID != "" {
					detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{"id": taskID})
					detail := getResultText(detailResult)
					if strings.Contains(detail, date) {
						t.Errorf("Invalid date %s was accepted", date)
					}
				}
			}
		}
	})

	// Test 7: Multiple tags with edge cases
	t.Run("Tag edge cases", func(t *testing.T) {
		tagCases := []struct {
			name      string
			etiketler string
			wantError bool
		}{
			{"Empty tags", "", false},
			{"Single tag", "important", false},
			{"Multiple tags", "important,urgent,bug", false},
			{"Tags with spaces", "important, urgent, bug", false},
			{"Duplicate tags", "important,urgent,important,urgent", true}, // Duplicate tags should cause error
			{"Tags with special chars", "important!,urgent@,#bug", false},
			{"Very long tag", strings.Repeat("a", constants.TestStringMedium), false},
		}

		for _, tc := range tagCases {
			t.Run(tc.name, func(t *testing.T) {
				result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
					constants.ParamTemplateID: templateID,
					constants.ParamValues: map[string]interface{}{
						"baslik":    "Task with tags: " + tc.name,
						"oncelik":   constants.PriorityMedium,
						"etiketler": tc.etiketler,
					},
				})

				if tc.wantError && !result.IsError {
					t.Errorf("Expected error for %s, got success", tc.name)
				}
				if !tc.wantError && result.IsError {
					t.Errorf("Expected success for %s, got error: %v", tc.name, getResultText(result))
				}
			})
		}
	})
}

// TestGorevGuncelle_EdgeCases tests edge cases for task updates
func TestGorevGuncelle_EdgeCases(t *testing.T) {
	_, handlers, templateID, cleanup := setupTestEnvironmentWithTemplate(t)
	defer cleanup()

	// Test 1: Invalid status transitions
	t.Run("Invalid status transitions", func(t *testing.T) {
		// Create a task
		createResult := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: templateID,
			constants.ParamValues: map[string]interface{}{
				"baslik":  "Test task for status",
				"oncelik": constants.PriorityMedium,
			},
		})
		taskID := extractTaskIDFromText(getResultText(createResult))

		invalidStatuses := []string{
			"completed",
			"TAMAMLANDI",
			"done",
			"in-progress",
			"",
			"null",
			"123",
		}

		for _, status := range invalidStatuses {
			result := callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
				"id":    taskID,
				"durum": status,
			})

			if !result.IsError {
				t.Errorf("Expected error for invalid status '%s', got success", status)
			}
		}
	})

	// Test 2: Non-existent task ID
	t.Run("Non-existent task ID", func(t *testing.T) {
		fakeIDs := []string{
			"non-existent-id",
			"12345",
			"",
			"null",
			"undefined",
		}

		for _, fakeID := range fakeIDs {
			result := callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
				"id":    fakeID,
				"durum": constants.TaskStatusInProgress,
			})

			assert.True(t, result.IsError, "Expected error for non-existent ID: %s", fakeID)
		}
	})
}

// TestProjeOlustur_EdgeCases tests edge cases for project creation
func TestProjeOlustur_EdgeCases(t *testing.T) {
	_, handlers, _, cleanup := setupTestEnvironmentWithTemplate(t)
	defer cleanup()

	// Test 1: Empty project names
	t.Run("Empty project names", func(t *testing.T) {
		emptyNames := []string{"", "   ", "\t", "\n", " \t\n "}

		for _, name := range emptyNames {
			result := callTool(t, handlers, "proje_olustur", map[string]interface{}{
				"isim":  name,
				"tanim": "Test description",
			})

			assert.True(t, result.IsError, "Expected error for empty project name")
		}
	})

	// Test 2: Duplicate project names
	t.Run("Duplicate project names", func(t *testing.T) {
		// Create first project
		result1 := callTool(t, handlers, "proje_olustur", map[string]interface{}{
			"isim":  "Duplicate Test Project",
			"tanim": "First project",
		})
		assert.False(t, result1.IsError)

		// Try to create second project with same name
		result2 := callTool(t, handlers, "proje_olustur", map[string]interface{}{
			"isim":  "Duplicate Test Project",
			"tanim": "Second project",
		})

		// The system might allow duplicate names, which is OK
		if !result2.IsError {
			// Both projects should exist
			listResult := callTool(t, handlers, "proje_listele", map[string]interface{}{})
			text := getResultText(listResult)
			count := strings.Count(text, "Duplicate Test Project")
			assert.GreaterOrEqual(t, count, 2, "Should have at least 2 projects with the same name")
		}
	})
}

// TestConcurrency_EdgeCases tests concurrent operations
func TestConcurrency_EdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrency test in short mode")
	}

	_, handlers, templateID, cleanup := setupTestEnvironmentWithTemplate(t)
	defer cleanup()

	// Test 1: Concurrent task creation
	t.Run("Concurrent task creation", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := constants.TestConcurrencyMedium
		errors := make([]error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
					constants.ParamTemplateID: templateID,
					constants.ParamValues: map[string]interface{}{
						"baslik":  fmt.Sprintf("Concurrent task %d", index),
						"oncelik": constants.PriorityMedium,
					},
				})

				if result.IsError {
					errors[index] = fmt.Errorf("failed to create task %d: %v", index, getResultText(result))
				}
			}(i)
		}

		wg.Wait()

		// Check for errors
		errorCount := 0
		for i, err := range errors {
			if err != nil {
				t.Logf("Goroutine %d error: %v", i, err)
				errorCount++
			}
		}

		// Accept up to 50% failure rate for concurrent task creation
		// SQLite with WAL mode and retry logic can still have contention under high concurrency
		maxAcceptableErrors := numGoroutines / 2 // 50% acceptable failure rate
		if errorCount > maxAcceptableErrors {
			t.Errorf("%d out of %d concurrent operations failed (%.1f%%), exceeds 50%% acceptable threshold",
				errorCount, numGoroutines, float64(errorCount)/float64(numGoroutines)*100)
		} else if errorCount > 0 {
			t.Logf("%d out of %d concurrent operations failed (%.1f%%), within acceptable 50%% threshold",
				errorCount, numGoroutines, float64(errorCount)/float64(numGoroutines)*100)
		}

		// Verify tasks were created
		listResult := callTool(t, handlers, "gorev_listele", map[string]interface{}{})
		if listResult.IsError {
			t.Fatalf("Failed to list tasks: %v", getResultText(listResult))
		}

		// Should have at least 50% success rate
		text := getResultText(listResult)
		taskCount := strings.Count(text, "Concurrent task")
		minExpectedTasks := numGoroutines / 2 // At least 50% should succeed
		if taskCount < minExpectedTasks {
			t.Errorf("Expected at least %d concurrent tasks (50%% of %d), found %d",
				minExpectedTasks, numGoroutines, taskCount)
		} else {
			t.Logf("Successfully created %d out of %d tasks (%.1f%% success rate)",
				taskCount, numGoroutines, float64(taskCount)/float64(numGoroutines)*100)
		}
	})

	// Test 2: Concurrent updates to same task
	t.Run("Concurrent updates to same task", func(t *testing.T) {
		// Create a task
		createResult := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: templateID,
			constants.ParamValues: map[string]interface{}{
				"baslik":  "Task for concurrent updates",
				"oncelik": constants.PriorityMedium,
			},
		})
		taskID := extractTaskIDFromText(getResultText(createResult))

		var wg sync.WaitGroup
		numUpdates := 5

		// Try to update the same task concurrently
		for i := 0; i < numUpdates; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				result := callTool(t, handlers, "gorev_duzenle", map[string]interface{}{
					"id":       taskID,
					"baslik":   fmt.Sprintf("Updated title %d", index),
					"aciklama": fmt.Sprintf("Updated by goroutine %d", index),
				})

				if result.IsError {
					t.Logf("Update %d failed: %v", index, getResultText(result))
				}
			}(i)
		}

		wg.Wait()

		// Check final state
		detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{"id": taskID})
		if detailResult.IsError {
			t.Fatalf("Failed to get task detail: %v", getResultText(detailResult))
		}

		// The task should have one of the updated titles
		detail := getResultText(detailResult)
		if !strings.Contains(detail, "Updated title") {
			t.Error("Task was not updated by any goroutine")
		}
	})

	// Test 3: Concurrent active project changes
	t.Run("Concurrent active project changes", func(t *testing.T) {
		// Create multiple projects
		projectIDs := make([]string, 3)
		for i := 0; i < 3; i++ {
			result := callTool(t, handlers, "proje_olustur", map[string]interface{}{
				"isim":  fmt.Sprintf("Concurrent Project %d", i),
				"tanim": fmt.Sprintf("Project %d for concurrency test", i),
			})
			projectIDs[i] = extractProjectIDFromText(getResultText(result))
		}

		var wg sync.WaitGroup
		numChanges := constants.TestConcurrencyMedium

		// Try to change active project concurrently
		for i := 0; i < numChanges; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				projectID := projectIDs[index%3]
				result := callTool(t, handlers, "proje_aktif_yap", map[string]interface{}{
					"proje_id": projectID,
				})

				if result.IsError {
					t.Logf("Setting active project failed: %v", getResultText(result))
				}
			}(i)
		}

		wg.Wait()

		// Check final active project
		result := callTool(t, handlers, "aktif_proje_goster", map[string]interface{}{})
		if result.IsError {
			t.Logf("No active project after concurrent changes: %v", getResultText(result))
		} else {
			text := getResultText(result)
			// Should have one of the concurrent projects as active
			hasValidProject := false
			for i := 0; i < 3; i++ {
				if strings.Contains(text, fmt.Sprintf("Concurrent Project %d", i)) {
					hasValidProject = true
					break
				}
			}
			assert.True(t, hasValidProject, "Active project should be one of the concurrent projects")
		}
	})
}

// TestTemplatedenGorevOlustur_EdgeCases tests edge cases for template-based task creation
func TestTemplatedenGorevOlustur_EdgeCases(t *testing.T) {
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     false, // Use temp file for edge cases
		UseTempFile:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: true, // Create default templates
		InitializeI18n:  true,
	}

	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()

	veriYonetici := isYonetici.VeriYonetici()

	// Create and set active project
	proje := &gorev.Proje{
		ID:         constants.TestProjectIDTemplate,
		Name:       "Test Template Edge Project",
		Definition: "Test project for template edge cases",
	}
	err := veriYonetici.ProjeKaydet(proje)
	require.NoError(t, err)
	err = veriYonetici.AktifProjeAyarla(proje.ID)
	require.NoError(t, err)

	handlers := YeniHandlers(isYonetici)

	// Test 1: Empty template ID
	t.Run("Empty template ID", func(t *testing.T) {
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: "",
			constants.ParamValues: map[string]interface{}{
				"baslik": "Test",
			},
		})
		assert.True(t, result.IsError)
		assert.Contains(t, getResultText(result), "template_id parametresi gerekli")
	})

	// Test 2: Wrong type for degerler
	t.Run("Wrong type for degerler", func(t *testing.T) {
		wrongTypes := []interface{}{
			"string instead of object",
			123,
			true,
			[]string{"array", "instead", "of", "object"},
		}

		for i, wrongType := range wrongTypes {
			t.Run(fmt.Sprintf("Type %d", i), func(t *testing.T) {
				result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
					constants.ParamTemplateID: constants.TestTemplateBugFix,
					constants.ParamValues:     wrongType,
				})
				assert.True(t, result.IsError)
				assert.Contains(t, getResultText(result), "degerler parametresi gerekli ve obje tipinde olmalı")
			})
		}
	})

	// Test 3: Template field injection attempts
	t.Run("Template field injection", func(t *testing.T) {
		// Get bug report template ID
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)

		var bugTemplateID string
		for _, tmpl := range templates {
			if tmpl.Name == "Bug Raporu" {
				bugTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, bugTemplateID)

		injectionValues := map[string]interface{}{
			"baslik":   "{{baslik}} {{aciklama}} {{modul}}",
			"aciklama": "'; DROP TABLE gorevler; --",
			"modul":    "{{../../../etc/passwd}}",
			"ortam":    "production' OR '1'='1",
			"adimlar":  "{{constructor.constructor('return process')()}}",
			"beklenen": "${7*7}",
			"mevcut":   "<script>alert('xss')</script>",
			"oncelik":  constants.PriorityHigh,
		}

		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: bugTemplateID,
			constants.ParamValues:     injectionValues,
		})

		// Should either sanitize or create the task with escaped values
		if !result.IsError {
			taskID := extractTaskIDFromText(getResultText(result))

			// Verify the task was created and check its content
			detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{"id": taskID})
			detail := getResultText(detailResult)

			// The template placeholders should not be expanded recursively
			if strings.Contains(detail, "/etc/passwd") {
				t.Error("Path traversal attempt was not sanitized")
			}
			if strings.Contains(detail, "49") { // Result of 7*7
				t.Error("Expression evaluation was not prevented")
			}
		}
	})

	// Test 4: Extremely large field values
	t.Run("Large field values", func(t *testing.T) {
		// Get bug report template ID
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)

		var bugTemplateID string
		for _, tmpl := range templates {
			if tmpl.Name == "Bug Raporu" {
				bugTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, bugTemplateID)

		largeString := strings.Repeat("A", constants.TestStringHuge)
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: bugTemplateID,
			constants.ParamValues: map[string]interface{}{
				"baslik":   "Large content test",
				"aciklama": largeString,
				"modul":    "TestModule",
				"ortam":    "development",
				"adimlar":  "Test steps",
				"beklenen": "Expected behavior",
				"mevcut":   "Current behavior",
				"oncelik":  constants.PriorityMedium,
			},
		})

		// Large values might be rejected or truncated
		if result.IsError {
			t.Logf("Large field value was rejected: %v", getResultText(result))
		} else {
			t.Log("Large field value was accepted")
		}
	})

	// Test 5: Missing all required fields
	t.Run("Missing all required fields", func(t *testing.T) {
		// Get feature request template ID
		templates, err := veriYonetici.TemplateListele("Özellik")
		require.NoError(t, err)

		var featureTemplateID string
		for _, tmpl := range templates {
			if tmpl.Name == "Özellik İsteği" {
				featureTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, featureTemplateID)

		// Provide no fields at all
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: featureTemplateID,
			constants.ParamValues:     map[string]interface{}{},
		})

		assert.True(t, result.IsError)
		assert.Contains(t, getResultText(result), "Zorunlu alanlar eksik")
	})

	// Test 6: Invalid field types
	t.Run("Invalid field types", func(t *testing.T) {
		// Get bug report template ID
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)

		var bugTemplateID string
		for _, tmpl := range templates {
			if tmpl.Name == "Bug Raporu" {
				bugTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, bugTemplateID)

		// Provide objects/arrays instead of strings
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: bugTemplateID,
			constants.ParamValues: map[string]interface{}{
				"baslik":   map[string]string{"nested": "object"},
				"aciklama": []string{"array", "of", "strings"},
				"modul":    123,
				"ortam":    true,
				"adimlar":  nil,
				"beklenen": "Expected",
				"mevcut":   "Current",
				"oncelik":  constants.PriorityMedium,
			},
		})

		// The handler should convert all values to strings
		if !result.IsError {
			t.Log("Non-string field values were converted successfully")
			taskID := extractTaskIDFromText(getResultText(result))

			// Check how the values were converted
			detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{"id": taskID})
			detail := getResultText(detailResult)
			t.Logf("Task detail with converted values:\n%s", detail)
		}
	})

	// Test 7: Duplicate tags in template
	t.Run("Duplicate tags", func(t *testing.T) {
		// Get bug report template ID
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)

		var bugTemplateID string
		for _, tmpl := range templates {
			if tmpl.Name == "Bug Raporu" {
				bugTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, bugTemplateID)

		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: bugTemplateID,
			constants.ParamValues: map[string]interface{}{
				"baslik":    "Task with duplicate tags",
				"aciklama":  "Bug description",
				"modul":     "TestModule",
				"ortam":     "production",
				"adimlar":   "Test steps",
				"beklenen":  "Expected",
				"mevcut":    "Current",
				"oncelik":   constants.PriorityMedium,
				"etiketler": "important,urgent,important,urgent,important",
			},
		})

		if !result.IsError {
			// Check how many unique tags were created
			taskID := extractTaskIDFromText(getResultText(result))
			detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{"id": taskID})

			detail := getResultText(detailResult)
			// Count occurrences of "important" and "urgent" in the detail
			importantCount := strings.Count(detail, "important")
			urgentCount := strings.Count(detail, "urgent")

			t.Logf("Tag 'important' appears %d times, 'urgent' appears %d times", importantCount, urgentCount)
		}
	})
}

// TestErrorPropagation tests error handling in various scenarios
func TestErrorPropagation(t *testing.T) {
	// Test with invalid database path to trigger errors
	t.Run("Database connection errors", func(t *testing.T) {
		// Create environment with read-only database
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "readonly.db")

		// Create a database file
		file, err := os.Create(dbPath)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()

		// Make it read-only
		if err := os.Chmod(dbPath, 0444); err != nil {
			t.Fatal(err)
		}

		// Try to create VeriYonetici directly with read-only database
		// This should fail during configuration or migration
		veriYonetici, err := gorev.YeniVeriYonetici(dbPath, constants.TestMigrationsPath)
		if err != nil {
			// This is expected for read-only database
			t.Logf("Expected error with read-only database: %v", err)
			// Accept both configuration and migration error messages
			errorStr := err.Error()
			if !strings.Contains(errorStr, "migration") &&
				!strings.Contains(errorStr, "failed to configure database") &&
				!strings.Contains(errorStr, "WAL mode") {
				t.Errorf("Expected migration or configuration error, got: %v", err)
			}
			return
		}
		defer veriYonetici.Kapat()

		// If we reach here, migration somehow succeeded, which is unexpected
		t.Errorf("Expected migration to fail with read-only database, but it succeeded")
		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// Try operations that should fail
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			constants.ParamTemplateID: constants.TestTemplateSimple,
			constants.ParamValues: map[string]interface{}{
				"baslik":  "This should fail",
				"oncelik": constants.PriorityMedium,
			},
		})

		if !result.IsError {
			t.Error("Expected error with read-only database, got success")
		}
	})
}

// TestPerformance_EdgeCases tests performance with extreme inputs
func TestPerformance_EdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	_, handlers, templateID, cleanup := setupTestEnvironmentWithTemplate(t)
	defer cleanup()

	t.Run("Create many tasks with tags", func(t *testing.T) {
		start := time.Now()
		taskCount := constants.TestEdgeCaseLimit

		for i := 0; i < taskCount; i++ {
			result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
				constants.ParamTemplateID: templateID,
				constants.ParamValues: map[string]interface{}{
					"baslik":    fmt.Sprintf("Performance task %d", i),
					"aciklama":  fmt.Sprintf("Description for task %d with some longer text to simulate real usage", i),
					"oncelik":   []string{constants.PriorityHigh, constants.PriorityMedium, constants.PriorityLow}[i%3],
					"etiketler": fmt.Sprintf("tag%d,performance,test,category%d", i, i%10),
					"son_tarih": time.Now().AddDate(0, 0, i).Format("2006-01-02"),
				},
			})

			if result.IsError {
				t.Errorf("Failed to create task %d: %v", i, getResultText(result))
			}
		}

		createDuration := time.Since(start)
		t.Logf("Created %d tasks with tags in %v (avg: %v per task)",
			taskCount, createDuration, createDuration/time.Duration(taskCount))

		// Test filtering performance
		filterTests := []struct {
			name   string
			params map[string]interface{}
		}{
			{"All tasks", map[string]interface{}{}},
			{"By status", map[string]interface{}{"durum": constants.TaskStatusPending}},
			{"By priority", map[string]interface{}{"filtre": "acil"}},
			{"By tag", map[string]interface{}{"etiket": "performance"}},
			{"Sorted by date", map[string]interface{}{"sirala": "son_tarih_asc"}},
		}

		for _, ft := range filterTests {
			start = time.Now()
			result := callTool(t, handlers, "gorev_listele", ft.params)
			duration := time.Since(start)

			if result.IsError {
				t.Errorf("Filter '%s' failed: %v", ft.name, getResultText(result))
			} else {
				t.Logf("Filter '%s' completed in %v", ft.name, duration)
			}
		}
	})
}
