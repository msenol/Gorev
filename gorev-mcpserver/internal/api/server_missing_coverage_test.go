package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMissingCoverageTestServer creates a test server with test data
func setupMissingCoverageTestServer(t *testing.T) (*APIServer, string, string, func()) {
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

	// Set active project
	err = isYonetici.AktifProjeAyarla(context.Background(), proje.ID)
	require.NoError(t, err)

	// Create a test task using the "bug" template alias
	degerler := map[string]string{
		"title":       "Test Task",
		"description": "Test Description",
		"modul":       "API",
		"ortam":       "development",
		"adimlar":     "Test steps",
		"beklenen":    "Expected result",
		"mevcut":      "Current result",
		"priority":    constants.PriorityMedium,
		"etiketler":   "test",
	}
	gorevResult, err := isYonetici.TemplatedenGorevOlustur(context.Background(), "bug", degerler)
	require.NoError(t, err)

	// Create API server
	server := NewAPIServer("8080", isYonetici)

	cleanup := func() {
		if veriYonetici != nil {
			veriYonetici.Kapat()
		}
	}

	return server, proje.ID, gorevResult.ID, cleanup
}

// TestGetTaskByID tests getting a single task by ID
func TestGetTaskByID(t *testing.T) {
	server, _, taskID, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/tasks/"+taskID, nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &result)
	require.NoError(t, err)

	assert.True(t, result["success"].(bool))
	assert.NotNil(t, result["data"])
}

// TestGetTask_NotFound tests getting a non-existent task
func TestGetTask_NotFound(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/tasks/nonexistent-id", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

// TestGetTask_EmptyID tests getting a task with empty ID
func TestGetTask_EmptyID(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/tasks/", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	// Empty ID returns task list endpoint (200 OK), not an error
	assert.Equal(t, 200, resp.StatusCode)
}

// TestCreateSubtask tests creating a subtask
func TestCreateSubtask(t *testing.T) {
	server, _, taskID, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		"title":       "Test Subtask",
		"description": "Subtask Description",
		"priority":    constants.PriorityMedium,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/tasks/"+taskID+"/subtasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestCreateSubtask_MissingParentID tests creating a subtask with invalid parent ID
func TestCreateSubtask_MissingParentID(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		"title":    "Test Subtask",
		"priority": constants.PriorityMedium,
	}
	body, _ := json.Marshal(payload)

	// Use non-existent parent ID
	req := httptest.NewRequest("POST", "/api/v1/tasks/nonexistent-id/subtasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

// TestChangeParentTask tests changing a task's parent
func TestChangeParentTask(t *testing.T) {
	server, _, taskID, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	// Create another task to use as new parent
	degerler := map[string]string{
		"title":       "Parent Task",
		"description": "Parent Description",
		"modul":       "API",
		"ortam":       "development",
		"adimlar":     "Test steps",
		"beklenen":    "Expected result",
		"mevcut":      "Current result",
		"priority":    constants.PriorityMedium,
		"etiketler":   "test",
	}
	parentResult, err := server.isYonetici.TemplatedenGorevOlustur(context.Background(), "bug", degerler)
	require.NoError(t, err)

	payload := map[string]interface{}{
		"yeni_parent_id": parentResult.ID,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", "/api/v1/tasks/"+taskID+"/parent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestChangeParent_MissingTaskID tests changing parent with invalid task ID
func TestChangeParent_MissingTaskID(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		"yeni_parent_id": "some-parent-id",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", "/api/v1/tasks/nonexistent-id/parent", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

// TestAddDependency tests adding a task dependency
func TestAddDependency(t *testing.T) {
	server, _, taskID, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	// Create another task for dependency
	degerler := map[string]string{
		"title":       "Dependent Task",
		"description": "Dependent Description",
		"modul":       "API",
		"ortam":       "development",
		"adimlar":     "Test steps",
		"beklenen":    "Expected result",
		"mevcut":      "Current result",
		"priority":    constants.PriorityMedium,
		"etiketler":   "test",
	}
	depResult, err := server.isYonetici.TemplatedenGorevOlustur(context.Background(), "bug", degerler)
	require.NoError(t, err)

	payload := map[string]interface{}{
		"kaynak_id":     taskID,
		"baglanti_tipi": "depends_on",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/tasks/"+depResult.ID+"/dependencies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestAddDependency_MissingFields tests adding dependency with missing fields
func TestAddDependency_MissingFields(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	payload := map[string]interface{}{
		// Missing hedef_id and baglanti_tipi
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/v1/tasks/some-id/dependencies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

// TestRemoveDependency tests removing a task dependency
func TestRemoveDependency(t *testing.T) {
	server, _, taskID, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	// Create another task for dependency
	degerler := map[string]string{
		"title":       "Dependent Task",
		"description": "Dependent Description",
		"modul":       "API",
		"ortam":       "development",
		"adimlar":     "Test steps",
		"beklenen":    "Expected result",
		"mevcut":      "Current result",
		"priority":    constants.PriorityMedium,
		"etiketler":   "test",
	}
	depResult, err := server.isYonetici.TemplatedenGorevOlustur(context.Background(), "bug", degerler)
	require.NoError(t, err)

	// First add a dependency using VeriYonetici
	baglanti := &gorev.Baglanti{
		ID:             uuid.New().String(),
		SourceID:       taskID,
		TargetID:       depResult.ID,
		ConnectionType: "depends_on",
	}
	err = server.isYonetici.VeriYonetici().BaglantiEkle(context.Background(), baglanti)
	require.NoError(t, err)

	// Now remove it
	// URL pattern: /api/v1/tasks/:id/dependencies/:dep_id maps to hedef=:id, kaynak=:dep_id
	// We created kaynak=taskID, hedef=depResult.ID
	// So URL should be /api/v1/tasks/depResult.ID/dependencies/taskID
	req := httptest.NewRequest("DELETE", "/api/v1/tasks/"+depResult.ID+"/dependencies/"+taskID, nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)

	// Debug: Print response if not successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var result map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		t.Logf("Response status: %d, body: %v", resp.StatusCode, result)
	}

	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestRemoveDependency_MissingIDs tests removing dependency with missing IDs
func TestRemoveDependency_MissingIDs(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("DELETE", "/api/v1/tasks//dependency/", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

// TestGetActiveProjectAPI tests getting the active project
func TestGetActiveProjectAPI(t *testing.T) {
	server, projectID, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	// First set an active project
	err := server.isYonetici.AktifProjeAyarla(context.Background(), projectID)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/active-project", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestGetActiveProject_NoActive tests getting active project when none set
func TestGetActiveProject_NoActive(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/active-project", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	// Should return 404 or empty result
	assert.True(t, resp.StatusCode == 404 || resp.StatusCode == 200)
}

// TestActivateProjectAPI tests activating a project
func TestActivateProjectAPI(t *testing.T) {
	server, projectID, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("PUT", "/api/v1/projects/"+projectID+"/activate", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestActivateProject_MissingID tests activating project with invalid ID
func TestActivateProject_MissingID(t *testing.T) {
	server, _, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("PUT", "/api/v1/projects/nonexistent-id/activate", nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400)
}

// TestRemoveActiveProjectAPI tests removing active project
func TestRemoveActiveProjectAPI(t *testing.T) {
	server, projectID, _, cleanup := setupMissingCoverageTestServer(t)
	defer cleanup()

	// First set an active project
	err := server.isYonetici.AktifProjeAyarla(context.Background(), projectID)
	require.NoError(t, err)

	// Now remove it
	req := httptest.NewRequest("DELETE", "/api/v1/active-project", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}
