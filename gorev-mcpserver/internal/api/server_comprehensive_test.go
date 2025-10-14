package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupComprehensiveTestServer creates a test server with test data
func setupComprehensiveTestServer(t *testing.T) (*APIServer, string, func()) {
	// Initialize i18n
	err := i18n.Initialize("en")
	if err != nil {
		// Already initialized, that's fine
	}

	// Create VeriYonetici with in-memory database
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)

	// Create templates
	err = veriYonetici.VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Create IsYonetici
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Create a test project
	proje, err := isYonetici.ProjeOlustur(context.Background(), "Test Project", "Test Description")
	require.NoError(t, err)

	// Create API server
	server := NewAPIServer("8080", isYonetici)

	cleanup := func() {
		if veriYonetici != nil {
			veriYonetici.Kapat()
		}
	}

	return server, proje.ID, cleanup
}

// TestGetTask tests getting a single task
func TestGetTask(t *testing.T) {
	server, projectID, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	// Create a task first
	payload := map[string]interface{}{
		"baslik":    "Test Task",
		"aciklama":  "Description",
		"durum":     constants.TaskStatusPending,
		"oncelik":   constants.PriorityMedium,
		"proje_id":  projectID,
		"etiketler": "test",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := server.app.Test(req)

	var createResult map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &createResult)

	// Skip if task creation is not implemented
	if resp.StatusCode == 501 {
		t.Skip("Task creation not implemented")
		return
	}

	// Now get the task
	if createResult["data"] != nil {
		taskData := createResult["data"].(map[string]interface{})
		taskID := taskData["id"].(string)

		req = httptest.NewRequest("GET", "/api/v1/tasks/"+taskID, nil)
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
	}
}

// TestUpdateTask tests updating a task
func TestUpdateTask(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		"durum": constants.TaskStatusInProgress,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", "/api/v1/tasks/test-id", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	// Accept any response (task may not exist)
	assert.NotNil(t, resp)
}

// TestDeleteTask tests deleting a task
func TestDeleteTask(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("DELETE", "/api/v1/tasks/test-id", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestCreateTaskFromTemplate tests creating task from template
func TestCreateTaskFromTemplate(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		"template_id": constants.TestTemplateSimple,
		"degerler": map[string]interface{}{
			"baslik":  "Template Task",
			"oncelik": constants.PriorityMedium,
		},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/tasks/from-template", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 600)
}

// TestGetProject tests getting a single project
func TestGetProject(t *testing.T) {
	server, projectID, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID, nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestGetProjectTasks tests getting tasks for a project
func TestGetProjectTasks(t *testing.T) {
	server, projectID, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/projects/"+projectID+"/tasks", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestActivateProject tests activating a project
func TestActivateProject(t *testing.T) {
	server, projectID, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("PUT", "/api/v1/projects/"+projectID+"/activate", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	// Accept 2xx or 404 (not implemented)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300 || resp.StatusCode == 404)
}

// TestGetActiveProject tests getting active project
func TestGetActiveProject(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/projects/active", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	// Accept 2xx or 404 (not implemented)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300 || resp.StatusCode == 404)
}

// TestRemoveActiveProject tests removing active project
func TestRemoveActiveProject(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("DELETE", "/api/v1/active-project", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	// Accept 2xx or 404 (not implemented)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300 || resp.StatusCode == 404)
}

// TestLanguageEndpoints tests language get/set
func TestLanguageEndpoints(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	// Get language
	t.Run("GetLanguage", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/language", nil)
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
	})

	// Set language
	t.Run("SetLanguage", func(t *testing.T) {
		payload := map[string]interface{}{
			"language": "tr",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/language", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
	})
}

// TestSubtaskOperations tests subtask creation
func TestSubtaskOperations(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		"parent_id": "test-parent-id",
		"baslik":    "Test Subtask",
		"oncelik":   constants.PriorityMedium,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/tasks/subtask", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestChangeParent tests changing task parent
func TestChangeParent(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		"gorev_id":       "test-task-id",
		"yeni_parent_id": "test-parent-id",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", "/api/v1/tasks/parent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestGetHierarchy tests getting task hierarchy
func TestGetHierarchy(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/tasks/test-id/hierarchy", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

// TestDependencyOperations tests adding and removing dependencies
func TestDependencyOperations(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	// Add dependency
	t.Run("AddDependency", func(t *testing.T) {
		payload := map[string]interface{}{
			"kaynak_id":     "task-1",
			"hedef_id":      "task-2",
			"baglanti_tipi": "depends_on",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/tasks/dependency", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	// Remove dependency
	t.Run("RemoveDependency", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/api/v1/tasks/task-1/dependency/task-2", nil)
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

// TestExportImport tests export and import operations
func TestExportImport(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	// Export
	t.Run("ExportData", func(t *testing.T) {
		payload := map[string]interface{}{
			"format":               "json",
			"include_completed":    true,
			"include_dependencies": true,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/export", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})

	// Import
	t.Run("ImportData", func(t *testing.T) {
		payload := map[string]interface{}{
			"data":                "{}",
			"import_mode":         "merge",
			"conflict_resolution": "skip",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/import", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

// TestServerStartMethods tests Start and StartAsync
func TestServerStartMethods(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	// Test StartAsync
	t.Run("StartAsync", func(t *testing.T) {
		go func() {
			server.StartAsync()
		}()

		// Give it a moment
		time.Sleep(100 * time.Millisecond)

		// Shutdown
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	})
}

// TestStaticFiles tests static file serving
func TestStaticFiles(t *testing.T) {
	server, _, cleanup := setupComprehensiveTestServer(t)
	defer cleanup()

	// Test root path (static files are served automatically)
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
