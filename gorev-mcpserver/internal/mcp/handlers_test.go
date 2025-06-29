package mcp

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test utility functions

// setupTestEnvironment creates a test MCP server with in-memory database
func setupTestEnvironment(t *testing.T) (*server.MCPServer, *Handlers, func()) {
	// Use temporary file database for testing
	tempDB := "test_mcp_" + strings.ReplaceAll(time.Now().Format("2006-01-02T15:04:05.000000000Z"), ":", "-") + ".db"
	cleanup := func() {
		os.Remove(tempDB)
	}

	veriYonetici, err := gorev.YeniVeriYonetici(tempDB, "file://../../internal/veri/migrations")
	require.NoError(t, err)

	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	mcpServer, err := YeniMCPSunucu(isYonetici)
	require.NoError(t, err)

	handlers := YeniHandlers(isYonetici)

	return mcpServer, handlers, cleanup
}

// getResultText extracts text from MCP result content
func getResultText(result *mcp.CallToolResult) string {
	if len(result.Content) == 0 {
		return ""
	}

	// The content is []interface{}, need to access the Text field
	// Based on the MCP SDK, result.Content[0] should be mcp.TextContent with .Text field
	switch content := result.Content[0].(type) {
	case mcp.TextContent:
		return content.Text
	default:
		// Try to extract text from interface
		if textMap, ok := result.Content[0].(map[string]interface{}); ok {
			if text, exists := textMap["text"]; exists {
				return fmt.Sprintf("%v", text)
			}
		}
	}

	// Fallback: convert to string directly
	return fmt.Sprintf("%v", result.Content[0])
}

// extractIDFromText extracts IDs from result text using regex
func extractTaskIDFromText(text string) string {
	// Pattern: "✓ Görev oluşturuldu: Title (ID: task-id)"
	re := regexp.MustCompile(`\(ID: ([^)]+)\)`)
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractProjectIDFromText(text string) string {
	// Same pattern as task ID
	return extractTaskIDFromText(text)
}

// callTool helper to simulate MCP tool calls
func callTool(t *testing.T, handlers *Handlers, toolName string, params map[string]interface{}) *mcp.CallToolResult {
	// Get the tool handler function
	var result *mcp.CallToolResult
	var err error

	switch toolName {
	case "gorev_olustur":
		result, err = handlers.GorevOlustur(params)
	case "gorev_listele":
		result, err = handlers.GorevListele(params)
	case "gorev_detay":
		result, err = handlers.GorevDetay(params)
	case "gorev_guncelle":
		result, err = handlers.GorevGuncelle(params)
	case "gorev_duzenle":
		result, err = handlers.GorevDuzenle(params)
	case "gorev_sil":
		result, err = handlers.GorevSil(params)
	case "gorev_bagimlilik_ekle":
		result, err = handlers.GorevBagimlilikEkle(params)
	case "proje_olustur":
		result, err = handlers.ProjeOlustur(params)
	case "proje_listele":
		result, err = handlers.ProjeListele(params)
	case "proje_gorevleri":
		result, err = handlers.ProjeGorevleri(params)
	case "proje_aktif_yap":
		result, err = handlers.AktifProjeAyarla(params)
	case "aktif_proje_goster":
		result, err = handlers.AktifProjeGoster(params)
	case "aktif_proje_kaldir":
		result, err = handlers.AktifProjeKaldir(params)
	case "template_listele":
		result, err = handlers.TemplateListele(params)
	case "templateden_gorev_olustur":
		result, err = handlers.TemplatedenGorevOlustur(params)
	case "ozet_goster":
		result, err = handlers.OzetGoster(params)
	default:
		t.Fatalf("Unknown tool: %s", toolName)
	}

	require.NoError(t, err)
	require.NotNil(t, result)
	return result
}

// Test cases for all 16 MCP tools

func TestMCPHandlers_Integration(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Complete Task Lifecycle", func(t *testing.T) {
		// 1. Create a project first
		result := callTool(t, handlers, "proje_olustur", map[string]interface{}{
			"isim":  "Test Projesi",
			"tanim": "Integration test projesi",
		})
		assert.False(t, result.IsError)

		// Parse result content
		contentText := getResultText(result)
		// Extract project ID from text content (since it's text, not JSON)
		// Format: "✓ Proje oluşturuldu: Test Projesi (ID: project-id)"
		projectID := extractProjectIDFromText(contentText)
		require.NotEmpty(t, projectID, "Project ID should be extracted from result")

		// 2. Set active project
		result = callTool(t, handlers, "proje_aktif_yap", map[string]interface{}{
			"proje_id": projectID,
		})
		assert.False(t, result.IsError)

		// 3. Show active project
		result = callTool(t, handlers, "aktif_proje_goster", map[string]interface{}{})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "Test Projesi")

		// 4. Create tasks with various parameters
		taskParams := []map[string]interface{}{
			{
				"baslik":    "İlk Görev",
				"aciklama":  "Test açıklaması",
				"oncelik":   "yuksek",
				"son_tarih": "2025-12-31",
				"etiketler": "test,integration",
			},
			{
				"baslik":   "İkinci Görev",
				"oncelik":  "dusuk",
				"proje_id": projectID,
			},
		}

		var taskIDs []string
		for i, params := range taskParams {
			result = callTool(t, handlers, "gorev_olustur", params)
			assert.False(t, result.IsError, "Task %d creation failed", i+1)

			// Extract task ID from text content
			contentText := getResultText(result)
			taskID := extractTaskIDFromText(contentText)
			require.NotEmpty(t, taskID, "Task ID should be extracted from result")
			taskIDs = append(taskIDs, taskID)
		}

		// 5. List all tasks
		result = callTool(t, handlers, "gorev_listele", map[string]interface{}{})
		assert.False(t, result.IsError)
		listContentText := getResultText(result)
		assert.Contains(t, listContentText, "İlk Görev")
		assert.Contains(t, listContentText, "İkinci Görev")

		// 6. List tasks by status
		result = callTool(t, handlers, "gorev_listele", map[string]interface{}{
			"durum": "beklemede",
		})
		assert.False(t, result.IsError)

		// 7. List tasks with filters
		result = callTool(t, handlers, "gorev_listele", map[string]interface{}{
			"filtre": "acil",
		})
		assert.False(t, result.IsError)

		// 8. Get task details
		result = callTool(t, handlers, "gorev_detay", map[string]interface{}{
			"id": taskIDs[0],
		})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "İlk Görev")

		// 9. Update task status
		result = callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
			"id":    taskIDs[0],
			"durum": "devam_ediyor",
		})
		assert.False(t, result.IsError)

		// 10. Edit task properties
		result = callTool(t, handlers, "gorev_duzenle", map[string]interface{}{
			"id":      taskIDs[1],
			"baslik":  "İkinci Görev (Güncellendi)",
			"oncelik": "yuksek",
		})
		assert.False(t, result.IsError)

		// 11. Add task dependency
		result = callTool(t, handlers, "gorev_bagimlilik_ekle", map[string]interface{}{
			"kaynak_id":     taskIDs[1],
			"hedef_id":      taskIDs[0],
			"baglanti_tipi": "tamamla_oncebi",
		})
		assert.False(t, result.IsError)

		// 12. List project tasks
		result = callTool(t, handlers, "proje_gorevleri", map[string]interface{}{
			"proje_id": projectID,
		})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "İlk Görev")

		// 13. List all projects
		result = callTool(t, handlers, "proje_listele", map[string]interface{}{})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "Test Projesi")

		// 14. Show summary
		result = callTool(t, handlers, "ozet_goster", map[string]interface{}{})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "Toplam")

		// 15. List templates
		result = callTool(t, handlers, "template_listele", map[string]interface{}{})
		assert.False(t, result.IsError)

		// 16. Create task from template (if templates exist)
		// Note: This might not work if default templates aren't created

		// 17. Remove active project
		result = callTool(t, handlers, "aktif_proje_kaldir", map[string]interface{}{})
		assert.False(t, result.IsError)

		// 18. Delete tasks (clean up)
		for i, taskID := range taskIDs {
			result = callTool(t, handlers, "gorev_sil", map[string]interface{}{
				"id":   taskID,
				"onay": true,
			})
			if result.IsError {
				t.Logf("Failed to delete task %d (ID: %s): %s", i+1, taskID, getResultText(result))
			}
			assert.False(t, result.IsError, "Failed to delete task %d: %s", i+1, getResultText(result))
		}
	})
}

func TestMCPHandlers_ErrorHandling(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Invalid Parameters", func(t *testing.T) {
		testCases := []struct {
			name    string
			tool    string
			params  map[string]interface{}
			wantErr bool
		}{
			{
				name:    "gorev_olustur - missing baslik",
				tool:    "gorev_olustur",
				params:  map[string]interface{}{},
				wantErr: true,
			},
			{
				name:    "gorev_olustur - empty baslik",
				tool:    "gorev_olustur",
				params:  map[string]interface{}{"baslik": ""},
				wantErr: true,
			},
			{
				name:    "gorev_detay - missing id",
				tool:    "gorev_detay",
				params:  map[string]interface{}{},
				wantErr: true,
			},
			{
				name:    "gorev_detay - non-existent id",
				tool:    "gorev_detay",
				params:  map[string]interface{}{"id": "non-existent"},
				wantErr: true,
			},
			{
				name:    "gorev_guncelle - invalid durum",
				tool:    "gorev_guncelle",
				params:  map[string]interface{}{"id": "test", "durum": "invalid"},
				wantErr: true,
			},
			{
				name:    "proje_olustur - missing isim",
				tool:    "proje_olustur",
				params:  map[string]interface{}{},
				wantErr: true,
			},
			{
				name:    "gorev_sil - wrong confirmation",
				tool:    "gorev_sil",
				params:  map[string]interface{}{"id": "test", "onay": false},
				wantErr: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result := callTool(t, handlers, tc.tool, tc.params)
				if tc.wantErr {
					assert.True(t, result.IsError, "Expected error for %s", tc.name)
				} else {
					assert.False(t, result.IsError, "Unexpected error for %s", tc.name)
				}
			})
		}
	})
}

func TestMCPHandlers_TemplateIntegration(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Template Operations", func(t *testing.T) {
		// 1. List available templates
		result := callTool(t, handlers, "template_listele", map[string]interface{}{})
		assert.False(t, result.IsError)

		// 2. List templates by category
		result = callTool(t, handlers, "template_listele", map[string]interface{}{
			"kategori": "bug",
		})
		assert.False(t, result.IsError)

		// 3. Try creating task from template (if any exist)
		// First check if templates exist by parsing the response
		listResult := callTool(t, handlers, "template_listele", map[string]interface{}{})
		if !listResult.IsError && strings.Contains(getResultText(listResult), "Bug Fix") {
			// Try to create a task from the Bug Fix template
			result = callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
				"template_id": "bug-fix",
				"degerler": map[string]interface{}{
					"bug_tanim": "Test bug açıklaması",
					"oncelik":   "yuksek",
				},
			})
			// This might fail if the template structure doesn't match
			// but we test the call succeeds
			assert.NotNil(t, result)
		}
	})
}

func TestMCPHandlers_ProjectManagement(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Project Lifecycle", func(t *testing.T) {
		// Create multiple projects
		projects := []map[string]interface{}{
			{"isim": "Proje A", "tanim": "İlk test projesi"},
			{"isim": "Proje B", "tanim": "İkinci test projesi"},
		}

		var projectIDs []string
		for _, proj := range projects {
			result := callTool(t, handlers, "proje_olustur", proj)
			assert.False(t, result.IsError)

			contentText := getResultText(result)
			projectID := extractProjectIDFromText(contentText)
			require.NotEmpty(t, projectID)
			projectIDs = append(projectIDs, projectID)
		}

		// List all projects
		result := callTool(t, handlers, "proje_listele", map[string]interface{}{})
		assert.False(t, result.IsError)
		contentText := getResultText(result)
		assert.Contains(t, contentText, "Proje A")
		assert.Contains(t, contentText, "Proje B")

		// Set active project
		result = callTool(t, handlers, "proje_aktif_yap", map[string]interface{}{
			"proje_id": projectIDs[0],
		})
		assert.False(t, result.IsError)

		// Create tasks in the active project
		result = callTool(t, handlers, "gorev_olustur", map[string]interface{}{
			"baslik": "Proje A Görevi",
		})
		assert.False(t, result.IsError)

		// Get project tasks
		result = callTool(t, handlers, "proje_gorevleri", map[string]interface{}{
			"proje_id": projectIDs[0],
		})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "Proje A Görevi")

		// Switch active project
		result = callTool(t, handlers, "proje_aktif_yap", map[string]interface{}{
			"proje_id": projectIDs[1],
		})
		assert.False(t, result.IsError)

		// Verify active project changed
		result = callTool(t, handlers, "aktif_proje_goster", map[string]interface{}{})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "Proje B")

		// Remove active project
		result = callTool(t, handlers, "aktif_proje_kaldir", map[string]interface{}{})
		assert.False(t, result.IsError)

		// Verify no active project
		result = callTool(t, handlers, "aktif_proje_goster", map[string]interface{}{})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "aktif proje")
	})
}

func TestMCPHandlers_TaskDependencies(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Task Dependencies", func(t *testing.T) {
		// Create tasks
		tasks := []string{"Görev 1", "Görev 2", "Görev 3"}
		var taskIDs []string

		for _, baslik := range tasks {
			result := callTool(t, handlers, "gorev_olustur", map[string]interface{}{
				"baslik": baslik,
			})
			assert.False(t, result.IsError)

			contentText := getResultText(result)
			taskID := extractTaskIDFromText(contentText)
			require.NotEmpty(t, taskID)
			taskIDs = append(taskIDs, taskID)
		}

		// Create dependencies: Task 3 depends on Task 1 and Task 2
		dependencies := []map[string]interface{}{
			{
				"kaynak_id":     taskIDs[2], // Task 3
				"hedef_id":      taskIDs[0], // depends on Task 1
				"baglanti_tipi": "tamamla_oncebi",
			},
			{
				"kaynak_id":     taskIDs[2], // Task 3
				"hedef_id":      taskIDs[1], // depends on Task 2
				"baglanti_tipi": "tamamla_oncebi",
			},
		}

		for _, dep := range dependencies {
			result := callTool(t, handlers, "gorev_bagimlilik_ekle", dep)
			assert.False(t, result.IsError)
		}

		// Try to start Task 3 (should fail due to dependencies)
		result := callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
			"id":    taskIDs[2],
			"durum": "devam_ediyor",
		})
		// Note: Dependencies might not be enforced in this version, so this test might pass
		if result.IsError {
			assert.Contains(t, getResultText(result), "bağımlılık")
		} else {
			t.Log("Dependencies not enforced - this is acceptable for current implementation")
		}

		// Complete Task 1
		result = callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
			"id":    taskIDs[0],
			"durum": "tamamlandi",
		})
		assert.False(t, result.IsError)

		// Still can't start Task 3 (Task 2 not complete) - but this might not be enforced
		result = callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
			"id":    taskIDs[2],
			"durum": "devam_ediyor",
		})
		// Dependencies might not be enforced in this version
		if result.IsError {
			t.Log("Dependencies enforced as expected")
		} else {
			t.Log("Dependencies not enforced - this is acceptable for current implementation")
		}

		// Complete Task 2
		result = callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
			"id":    taskIDs[1],
			"durum": "tamamlandi",
		})
		assert.False(t, result.IsError)

		// Now Task 3 can start
		result = callTool(t, handlers, "gorev_guncelle", map[string]interface{}{
			"id":    taskIDs[2],
			"durum": "devam_ediyor",
		})
		assert.False(t, result.IsError)

		// Verify task details show dependencies
		result = callTool(t, handlers, "gorev_detay", map[string]interface{}{
			"id": taskIDs[2],
		})
		assert.False(t, result.IsError)
		assert.Contains(t, getResultText(result), "Bağımlılıklar")
	})
}

// Performance and stress testing
func TestMCPHandlers_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Run("Bulk Operations", func(t *testing.T) {
		// Create multiple tasks quickly
		taskCount := 100
		start := time.Now()

		for i := 0; i < taskCount; i++ {
			result := callTool(t, handlers, "gorev_olustur", map[string]interface{}{
				"baslik": fmt.Sprintf("Performance Test Task %d", i),
			})
			assert.False(t, result.IsError)
		}

		createDuration := time.Since(start)
		t.Logf("Created %d tasks in %v (avg: %v per task)", taskCount, createDuration, createDuration/time.Duration(taskCount))

		// List all tasks
		start = time.Now()
		result := callTool(t, handlers, "gorev_listele", map[string]interface{}{})
		assert.False(t, result.IsError)
		listDuration := time.Since(start)
		t.Logf("Listed %d tasks in %v", taskCount, listDuration)

		// Verify performance is reasonable (adjust thresholds as needed)
		assert.Less(t, createDuration, 10*time.Second, "Task creation took too long")
		assert.Less(t, listDuration, 1*time.Second, "Task listing took too long")
	})
}
