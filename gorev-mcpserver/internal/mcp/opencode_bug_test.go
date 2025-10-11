package mcp

import (
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	testingHelpers "github.com/msenol/gorev/internal/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOpenCodeAIIDSerializationBug tests the specific issue reported by OpenCode.ai
// where "Expected 'id' to be a string" error occurs in gorev_listele with noTasksInProject
func TestOpenCodeAIIDSerializationBug(t *testing.T) {
	// Setup test environment
	isYonetici, cleanup := testingHelpers.SetupTestEnvironmentWithConfig(t, &testingHelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  "../../internal/veri/migrations", // Correct path to migrations
		CreateTemplates: false,                            // Skip templates for this test
		InitializeI18n:  true,
	})
	defer cleanup()

	// Create handlers with debug enabled
	handlers := &Handlers{
		isYonetici:  isYonetici,
		debug:       true,
		toolHelpers: NewToolHelpers(),
	}

	// Test 1: Create a project without any tasks (this triggers the bug)
	proje := &gorev.Proje{
		ID:         "test-project-id",
		Name:       "Test Project",
		Definition: "Test project for OpenCode.ai bug reproduction",
	}
	err := isYonetici.VeriYonetici().ProjeKaydet(proje)
	require.NoError(t, err)

	// Set the project as active
	err = isYonetici.VeriYonetici().AktifProjeAyarla(proje.ID)
	require.NoError(t, err)

	// Test 2: Call gorev_listele with empty project (reproduces OpenCode.ai bug)
	params := map[string]interface{}{
		"tum_projeler": false, // This will use active project
	}

	result, err := handlers.GorevListele(params)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Test 3: Verify the result can be JSON serialized without type errors
	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err, "MCP result should be JSON serializable")

	// Test 4: Verify the JSON doesn't contain numeric IDs in unexpected places
	var resultJSON map[string]interface{}
	err = json.Unmarshal(jsonBytes, &resultJSON)
	require.NoError(t, err)

	// Print the actual JSON for debugging
	t.Logf("MCP Result JSON: %s", string(jsonBytes))

	// Test 5: Verify specific content expectations
	content := resultJSON["content"].([]interface{})
	firstContent := content[0].(map[string]interface{})
	text := firstContent["text"].(string)

	assert.Contains(t, text, "Test Project", "Result should contain project name")
	assert.Contains(t, text, i18n.T("messages.noTasksInProject", map[string]interface{}{"Project": "Test Project"}),
		"Result should contain the no tasks message")

	// Test 6: Validate that all ID fields in the response are strings
	validateIDFields(t, resultJSON, "")
}

// validateIDFields recursively checks that all 'id' fields in the JSON structure are strings
func validateIDFields(t *testing.T, data interface{}, path string) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			currentPath := path
			if currentPath != "" {
				currentPath += "."
			}
			currentPath += key

			if key == "id" {
				// This is the critical check: all 'id' fields must be strings
				assert.IsType(t, "", value, "ID field at %s should be string, got %T with value %v", currentPath, value, value)
			}
			validateIDFields(t, value, currentPath)
		}
	case []interface{}:
		for i, item := range v {
			currentPath := path + "[" + string(rune(i+'0')) + "]"
			validateIDFields(t, item, currentPath)
		}
	}
}

// TestMCPResultTextSerialization tests the basic MCP result serialization
func TestMCPResultTextSerialization(t *testing.T) {
	// Test simple text result
	result := mcp.NewToolResultText("Simple test message")

	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err)

	var resultJSON map[string]interface{}
	err = json.Unmarshal(jsonBytes, &resultJSON)
	require.NoError(t, err)

	// Validate structure
	validateIDFields(t, resultJSON, "")
	t.Logf("Simple MCP result JSON: %s", string(jsonBytes))
}

// TestProjectObjectSerialization tests serialization of Proje objects
func TestProjectObjectSerialization(t *testing.T) {
	proje := &gorev.Proje{
		ID:         "test-project-123",
		Name:       "Test Project",
		Definition: "Test description",
	}

	// Test direct JSON serialization of Proje object
	jsonBytes, err := json.Marshal(proje)
	require.NoError(t, err)

	var projeJSON map[string]interface{}
	err = json.Unmarshal(jsonBytes, &projeJSON)
	require.NoError(t, err)

	// Validate that project ID is string
	assert.IsType(t, "", projeJSON["id"], "Project ID should be string")
	assert.Equal(t, "test-project-123", projeJSON["id"])

	t.Logf("Project JSON: %s", string(jsonBytes))
}

// TestGorevListeleWithTasks tests gorev_listele when there are actual tasks
func TestGorevListeleWithTasks(t *testing.T) {
	// Setup test environment
	isYonetici, cleanup := testingHelpers.SetupTestEnvironmentWithConfig(t, &testingHelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  "../../internal/veri/migrations",
		CreateTemplates: false,
		InitializeI18n:  true,
	})
	defer cleanup()

	// Create handlers with debug enabled
	handlers := &Handlers{
		isYonetici:  isYonetici,
		debug:       true,
		toolHelpers: NewToolHelpers(),
	}

	// Create a project
	proje := &gorev.Proje{
		ID:         "test-project-with-tasks",
		Name:       "Project With Tasks",
		Definition: "Test project with tasks",
	}
	err := isYonetici.VeriYonetici().ProjeKaydet(proje)
	require.NoError(t, err)

	// Set as active project
	err = isYonetici.VeriYonetici().AktifProjeAyarla(proje.ID)
	require.NoError(t, err)

	// Create some tasks
	task1 := &gorev.Gorev{
		ID:          "task-1",
		Title:       "Task 1",
		Description: "First task",
		Status:      "beklemede",
		Priority:    "yuksek",
		ProjeID:     proje.ID,
	}
	err = isYonetici.VeriYonetici().GorevKaydet(task1)
	require.NoError(t, err)

	task2 := &gorev.Gorev{
		ID:          "task-2",
		Title:       "Task 2",
		Description: "Second task",
		Status:      "devam_ediyor",
		Priority:    "orta",
		ProjeID:     proje.ID,
	}
	err = isYonetici.VeriYonetici().GorevKaydet(task2)
	require.NoError(t, err)

	// Call gorev_listele
	params := map[string]interface{}{
		"tum_projeler": false,
	}

	result, err := handlers.GorevListele(params)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify JSON serialization
	jsonBytes, err := json.Marshal(result)
	require.NoError(t, err, "MCP result should be JSON serializable")

	var resultJSON map[string]interface{}
	err = json.Unmarshal(jsonBytes, &resultJSON)
	require.NoError(t, err)

	// Print the actual JSON for debugging
	t.Logf("MCP Result with tasks JSON: %s", string(jsonBytes))

	// Validate that all ID fields are strings
	validateIDFields(t, resultJSON, "")

	// Verify content contains task information
	content := resultJSON["content"].([]interface{})
	firstContent := content[0].(map[string]interface{})
	text := firstContent["text"].(string)

	assert.Contains(t, text, "Task 1", "Result should contain first task")
	assert.Contains(t, text, "Task 2", "Result should contain second task")
	assert.Contains(t, text, "task-1", "Result should contain task ID")
	assert.Contains(t, text, "task-2", "Result should contain task ID")
}
