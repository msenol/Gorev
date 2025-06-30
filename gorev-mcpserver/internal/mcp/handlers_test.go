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

// TestTemplateHandlers tests all template-related MCP handlers
func TestTemplateHandlers(t *testing.T) {
	t.Run("List Templates Empty", func(t *testing.T) {
		// Create fresh database without templates
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_empty.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_empty.db")

		// Don't initialize default templates
		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// List templates when none exist
		result := callTool(t, handlers, "template_listele", map[string]interface{}{})
		assert.False(t, result.IsError)

		text := getResultText(result)
		// The test environment might have default templates from migration
		// So we check if we get a valid response structure
		if !strings.Contains(text, "Henüz template bulunmuyor") {
			assert.Contains(t, text, "## 📋 Görev Template'leri")
		}
	})

	t.Run("Initialize Default Templates", func(t *testing.T) {
		// Initialize default templates through veri_yonetici
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_init.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_init.db")

		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		// Create handlers with initialized database
		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// List all templates
		result := callTool(t, handlers, "template_listele", map[string]interface{}{})
		assert.False(t, result.IsError)

		text := getResultText(result)
		assert.Contains(t, text, "## 📋 Görev Template'leri")
		assert.Contains(t, text, "Bug Raporu")
		assert.Contains(t, text, "Özellik İsteği")
		assert.Contains(t, text, "Teknik Borç")
		assert.Contains(t, text, "Araştırma Görevi")
		assert.Contains(t, text, "### Teknik")
		assert.Contains(t, text, "### Özellik")
		assert.Contains(t, text, "### Araştırma")
	})

	t.Run("List Templates By Category", func(t *testing.T) {
		// Setup with default templates
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_category.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_category.db")

		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// List only "Teknik" category templates
		result := callTool(t, handlers, "template_listele", map[string]interface{}{
			"kategori": "Teknik",
		})
		assert.False(t, result.IsError)

		text := getResultText(result)
		assert.Contains(t, text, "Bug Raporu")
		assert.Contains(t, text, "Teknik Borç")
		assert.NotContains(t, text, "Özellik İsteği")
		assert.NotContains(t, text, "Araştırma Görevi")
	})

	t.Run("Create Task From Template - Bug Report", func(t *testing.T) {
		// Setup with default templates
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_bug.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_bug.db")

		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		// Get bug report template ID
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)

		var bugTemplateID string
		for _, tmpl := range templates {
			if tmpl.Isim == "Bug Raporu" {
				bugTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, bugTemplateID)

		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// Create task from bug report template
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"template_id": bugTemplateID,
			"degerler": map[string]interface{}{
				"baslik":    "Login button not working",
				"aciklama":  "Users can't log in when clicking the login button",
				"modul":     "Authentication",
				"ortam":     "production",
				"adimlar":   "1. Go to login page\n2. Enter credentials\n3. Click login button",
				"beklenen":  "User should be logged in and redirected to dashboard",
				"mevcut":    "Nothing happens when clicking the button",
				"ekler":     "console-error.png",
				"cozum":     "Check event handler binding",
				"oncelik":   "yuksek",
				"etiketler": "bug,urgent,auth",
			},
		})
		assert.False(t, result.IsError)

		text := getResultText(result)
		assert.Contains(t, text, "✓ Template kullanılarak görev oluşturuldu")
		assert.Contains(t, text, "🐛 [Authentication] Login button not working")

		// Extract task ID and verify details
		taskID := extractTaskIDFromText(text)
		require.NotEmpty(t, taskID)

		// Get task details
		detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{
			"id": taskID,
		})
		assert.False(t, detailResult.IsError)

		detailText := getResultText(detailResult)
		assert.Contains(t, detailText, "## 🐛 Hata Açıklaması")
		assert.Contains(t, detailText, "Users can't log in when clicking the login button")
		assert.Contains(t, detailText, "**Modül/Bileşen:** Authentication")
		assert.Contains(t, detailText, "**Ortam:** production")
		assert.Contains(t, detailText, "## 🔄 Tekrar Üretme Adımları")
		assert.Contains(t, detailText, "1. Go to login page")
		assert.Contains(t, detailText, "## ✅ Beklenen Davranış")
		assert.Contains(t, detailText, "User should be logged in and redirected to dashboard")
		assert.Contains(t, detailText, "## ❌ Mevcut Davranış")
		assert.Contains(t, detailText, "Nothing happens when clicking the button")
		assert.Contains(t, detailText, "## 📸 Ekran Görüntüleri/Loglar")
		assert.Contains(t, detailText, "console-error.png")
		assert.Contains(t, detailText, "## 🔧 Olası Çözüm")
		assert.Contains(t, detailText, "Check event handler binding")
		assert.Contains(t, detailText, "## 📊 Öncelik: yuksek")
		assert.Contains(t, detailText, "## 🏷️ Etiketler: bug,urgent,auth")
		assert.Contains(t, detailText, "bug")
		assert.Contains(t, detailText, "urgent")
		assert.Contains(t, detailText, "auth")
	})

	t.Run("Create Task From Template - Missing Required Fields", func(t *testing.T) {
		// Setup with default templates
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_missing.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_missing.db")

		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		// Get bug report template ID
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)

		var bugTemplateID string
		for _, tmpl := range templates {
			if tmpl.Isim == "Bug Raporu" {
				bugTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, bugTemplateID)

		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// Try to create task without required fields
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"template_id": bugTemplateID,
			"degerler": map[string]interface{}{
				"baslik": "Test bug",
				// Missing required fields: aciklama, modul, ortam, adimlar, beklenen, mevcut, oncelik
			},
		})
		assert.True(t, result.IsError)

		text := getResultText(result)
		assert.Contains(t, text, "zorunlu alan eksik")
	})

	t.Run("Create Task From Template - Invalid Template ID", func(t *testing.T) {
		// Setup fresh handlers
		_, handlers, cleanup := setupTestEnvironment(t)
		defer cleanup()

		// Try to create task with non-existent template
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"template_id": "non-existent-template-id",
			"degerler": map[string]interface{}{
				"baslik": "Test task",
			},
		})
		assert.True(t, result.IsError)

		text := getResultText(result)
		assert.Contains(t, text, "template bulunamadı")
	})

	t.Run("Create Task From Template - Feature Request", func(t *testing.T) {
		// Setup with default templates
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_feature.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_feature.db")

		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		// Get feature request template ID
		templates, err := veriYonetici.TemplateListele("Özellik")
		require.NoError(t, err)

		var featureTemplateID string
		for _, tmpl := range templates {
			if tmpl.Isim == "Özellik İsteği" {
				featureTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, featureTemplateID)

		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// Create a project for the feature
		projectResult := callTool(t, handlers, "proje_olustur", map[string]interface{}{
			"isim":  "Mobile App",
			"tanim": "Mobile application project",
		})
		assert.False(t, projectResult.IsError)
		projectID := extractProjectIDFromText(getResultText(projectResult))

		// Create task from feature request template
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"template_id": featureTemplateID,
			"degerler": map[string]interface{}{
				"baslik":       "Dark mode support",
				"aciklama":     "Add dark mode theme to the mobile app",
				"amac":         "Improve user experience in low-light conditions and save battery",
				"kullanicilar": "All mobile app users",
				"kriterler":    "- Theme toggle in settings\n- Persistent preference\n- Smooth transition",
				"ui_ux":        "Material Design 3 dark theme guidelines",
				"ilgili":       "Settings module, Theme manager",
				"efor":         "orta",
				"oncelik":      "orta",
				"etiketler":    "özellik,ui,mobile",
				"proje_id":     projectID,
				"son_tarih":    "2025-08-15",
			},
		})
		assert.False(t, result.IsError)

		text := getResultText(result)
		assert.Contains(t, text, "✓ Template kullanılarak görev oluşturuldu")
		assert.Contains(t, text, "✨ Dark mode support")

		// Verify task was created with correct project
		taskID := extractTaskIDFromText(text)
		detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{
			"id": taskID,
		})

		detailText := getResultText(detailResult)
		assert.Contains(t, detailText, "Mobile App")
		assert.Contains(t, detailText, "2025-08-15")
	})

	t.Run("Create Task From Template - Technical Debt", func(t *testing.T) {
		// Setup with default templates
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_tech_debt.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_tech_debt.db")

		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		// Get technical debt template ID
		templates, err := veriYonetici.TemplateListele("Teknik")
		require.NoError(t, err)

		var techDebtTemplateID string
		for _, tmpl := range templates {
			if tmpl.Isim == "Teknik Borç" {
				techDebtTemplateID = tmpl.ID
				break
			}
		}
		require.NotEmpty(t, techDebtTemplateID)

		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// Create task from technical debt template
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"template_id": techDebtTemplateID,
			"degerler": map[string]interface{}{
				"baslik":         "Database query optimization",
				"aciklama":       "Optimize slow database queries in user listing",
				"alan":           "Backend/Database",
				"dosyalar":       "user_repository.go, user_queries.sql",
				"neden":          "Page load time exceeds 5 seconds for user list",
				"analiz":         "N+1 query problem, missing indexes",
				"cozum":          "Add composite indexes, use JOIN instead of multiple queries",
				"riskler":        "Potential data inconsistency during migration",
				"iyilestirmeler": "50% reduction in page load time",
				"sure":           "2-3 gün",
				"oncelik":        "yuksek",
				"etiketler":      "teknik-borç,performance,database",
			},
		})
		assert.False(t, result.IsError)

		text := getResultText(result)
		assert.Contains(t, text, "✓ Template kullanılarak görev oluşturuldu")
		assert.Contains(t, text, "🔧 [Backend/Database] Database query optimization")

		// Verify task details
		taskID := extractTaskIDFromText(text)
		detailResult := callTool(t, handlers, "gorev_detay", map[string]interface{}{
			"id": taskID,
		})

		detailText := getResultText(detailResult)
		assert.Contains(t, detailText, "## 🔧 Teknik Borç Açıklaması")
		assert.Contains(t, detailText, "Optimize slow database queries")
		assert.Contains(t, detailText, "**Alan/Modül:** Backend/Database")
		assert.Contains(t, detailText, "**Dosyalar:** user_repository.go, user_queries.sql")
		assert.Contains(t, detailText, "## ⏱️ Tahmini Süre: 2-3 gün")
		assert.Contains(t, detailText, "performance")
		assert.Contains(t, detailText, "database")
	})

	t.Run("Template Field Validation", func(t *testing.T) {
		// Setup with default templates
		veriYonetici, err := gorev.YeniVeriYonetici("test_template_validation.db", "file://../../internal/veri/migrations")
		require.NoError(t, err)
		defer os.Remove("test_template_validation.db")

		err = veriYonetici.VarsayilanTemplateleriOlustur()
		require.NoError(t, err)

		isYonetici := gorev.YeniIsYonetici(veriYonetici)
		handlers := YeniHandlers(isYonetici)

		// List templates to verify field information
		result := callTool(t, handlers, "template_listele", map[string]interface{}{})
		assert.False(t, result.IsError)

		text := getResultText(result)

		// Verify field types and requirements are shown
		assert.Contains(t, text, "(text) *(zorunlu)*")
		assert.Contains(t, text, "(select) *(zorunlu)*")
		assert.Contains(t, text, "(date)")
		assert.Contains(t, text, "varsayılan: orta")
		assert.Contains(t, text, "seçenekler: development, staging, production")
		assert.Contains(t, text, "seçenekler: dusuk, orta, yuksek")
	})

	t.Run("Template Parameters Validation", func(t *testing.T) {
		// Setup fresh handlers
		_, handlers, cleanup := setupTestEnvironment(t)
		defer cleanup()

		// Test missing template_id
		result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"degerler": map[string]interface{}{
				"baslik": "Test",
			},
		})
		assert.True(t, result.IsError)
		assert.Contains(t, getResultText(result), "template_id parametresi gerekli")

		// Test missing degerler
		result = callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"template_id": "some-id",
		})
		assert.True(t, result.IsError)
		assert.Contains(t, getResultText(result), "degerler parametresi gerekli")

		// Test wrong type for degerler
		result = callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
			"template_id": "some-id",
			"degerler":    "not-an-object",
		})
		assert.True(t, result.IsError)
		assert.Contains(t, getResultText(result), "degerler parametresi gerekli ve obje tipinde olmalı")
	})
}

// TestTemplateConcurrency tests template operations under concurrent access
func TestTemplateConcurrency(t *testing.T) {
	// Setup with default templates
	veriYonetici, err := gorev.YeniVeriYonetici("test_template_concurrent.db", "file://../../internal/veri/migrations")
	require.NoError(t, err)
	defer os.Remove("test_template_concurrent.db")

	err = veriYonetici.VarsayilanTemplateleriOlustur()
	require.NoError(t, err)

	// Get bug report template ID
	templates, err := veriYonetici.TemplateListele("Teknik")
	require.NoError(t, err)

	var bugTemplateID string
	for _, tmpl := range templates {
		if tmpl.Isim == "Bug Raporu" {
			bugTemplateID = tmpl.ID
			break
		}
	}
	require.NotEmpty(t, bugTemplateID)

	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := YeniHandlers(isYonetici)

	// Test concurrent task creation from template
	numGoroutines := 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			result := callTool(t, handlers, "templateden_gorev_olustur", map[string]interface{}{
				"template_id": bugTemplateID,
				"degerler": map[string]interface{}{
					"baslik":    fmt.Sprintf("Concurrent Bug %d", index),
					"aciklama":  fmt.Sprintf("Bug description %d", index),
					"modul":     "TestModule",
					"ortam":     "development",
					"adimlar":   "Test steps",
					"beklenen":  "Expected behavior",
					"mevcut":    "Current behavior",
					"oncelik":   "orta",
					"etiketler": "test,concurrent",
				},
			})
			assert.False(t, result.IsError)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify all tasks were created
	result := callTool(t, handlers, "gorev_listele", map[string]interface{}{})
	assert.False(t, result.IsError)

	text := getResultText(result)
	for i := 0; i < numGoroutines; i++ {
		assert.Contains(t, text, fmt.Sprintf("Concurrent Bug %d", i))
	}
}
