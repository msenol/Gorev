package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// E2E Test Scenarios for Multi-Workspace Support

// TestE2E_MultipleWorkspaceRegistration tests registering multiple workspaces
func TestE2E_MultipleWorkspaceRegistration(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "gorev-e2e-multi-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewWorkspaceManager()
	app := setupTestServer(manager)

	// Create three test workspaces
	workspaces := []struct {
		name string
		path string
	}{
		{"Project Alpha", filepath.Join(tempDir, "alpha")},
		{"Project Beta", filepath.Join(tempDir, "beta")},
		{"Project Gamma", filepath.Join(tempDir, "gamma")},
	}

	registeredIDs := make([]string, 0, len(workspaces))

	// Register all workspaces
	for _, ws := range workspaces {
		err := os.MkdirAll(ws.path, 0755)
		require.NoError(t, err)

		payload := map[string]string{
			"name": ws.name,
			"path": ws.path,
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/workspaces/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, true, result["success"])
		assert.NotEmpty(t, result["workspace_id"])

		registeredIDs = append(registeredIDs, result["workspace_id"].(string))
	}

	// Verify all workspaces are listed
	req := httptest.NewRequest("GET", "/api/workspaces", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&listResp)
	assert.Equal(t, true, listResp["success"])
	assert.Equal(t, float64(len(workspaces)), listResp["total"])

	respWorkspaces := listResp["workspaces"].([]interface{})
	assert.Len(t, respWorkspaces, len(workspaces))

	// Verify each workspace has correct details
	for i, wsInterface := range respWorkspaces {
		ws := wsInterface.(map[string]interface{})
		assert.Equal(t, workspaces[i].name, ws["name"])
		assert.Equal(t, workspaces[i].path, ws["path"])
		assert.Contains(t, registeredIDs, ws["id"])
	}
}

// TestE2E_WorkspaceDatabaseIsolation tests that tasks are isolated per workspace
func TestE2E_WorkspaceDatabaseIsolation(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "gorev-e2e-isolation-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewWorkspaceManager()
	app := setupTestServer(manager)

	// Register two workspaces
	workspace1Path := filepath.Join(tempDir, "workspace1")
	workspace2Path := filepath.Join(tempDir, "workspace2")

	err = os.MkdirAll(workspace1Path, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(workspace2Path, 0755)
	require.NoError(t, err)

	// Register workspace 1
	ws1ID := registerWorkspace(t, app, "Workspace 1", workspace1Path)

	// Register workspace 2
	ws2ID := registerWorkspace(t, app, "Workspace 2", workspace2Path)

	// Create multiple projects in workspace 1
	project1ID_1 := createProject(t, app, ws1ID, workspace1Path, "Workspace 1", "Project Alpha 1", "Alpha project 1")
	project1ID_2 := createProject(t, app, ws1ID, workspace1Path, "Workspace 1", "Project Alpha 2", "Alpha project 2")

	// Create multiple projects in workspace 2
	project2ID_1 := createProject(t, app, ws2ID, workspace2Path, "Workspace 2", "Project Beta 1", "Beta project 1")
	project2ID_2 := createProject(t, app, ws2ID, workspace2Path, "Workspace 2", "Project Beta 2", "Beta project 2")
	project2ID_3 := createProject(t, app, ws2ID, workspace2Path, "Workspace 2", "Project Beta 3", "Beta project 3")

	// Verify workspace isolation by checking project counts
	// Get workspace 1 context and verify it has 2 projects
	ws1Context, err := manager.GetWorkspaceContext(ws1ID)
	require.NoError(t, err)
	ws1Projects, err := ws1Context.IsYonetici.ProjeListele()
	require.NoError(t, err)
	assert.Len(t, ws1Projects, 2, "Workspace 1 should have 2 projects")

	// Get workspace 2 context and verify it has 3 projects
	ws2Context, err := manager.GetWorkspaceContext(ws2ID)
	require.NoError(t, err)
	ws2Projects, err := ws2Context.IsYonetici.ProjeListele()
	require.NoError(t, err)
	assert.Len(t, ws2Projects, 3, "Workspace 2 should have 3 projects")

	// Verify workspace 1 projects don't appear in workspace 2
	ws2ProjectIDs := make([]string, len(ws2Projects))
	for i, p := range ws2Projects {
		ws2ProjectIDs[i] = p.ID
	}
	assert.NotContains(t, ws2ProjectIDs, project1ID_1)
	assert.NotContains(t, ws2ProjectIDs, project1ID_2)

	// Verify workspace 2 projects don't appear in workspace 1
	ws1ProjectIDs := make([]string, len(ws1Projects))
	for i, p := range ws1Projects {
		ws1ProjectIDs[i] = p.ID
	}
	assert.NotContains(t, ws1ProjectIDs, project2ID_1)
	assert.NotContains(t, ws1ProjectIDs, project2ID_2)
	assert.NotContains(t, ws1ProjectIDs, project2ID_3)
}

// TestE2E_ConcurrentWorkspaceAccess tests concurrent access from multiple workspaces
func TestE2E_ConcurrentWorkspaceAccess(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "gorev-e2e-concurrent-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewWorkspaceManager()
	app := setupTestServer(manager)

	// Register 5 workspaces
	numWorkspaces := 5
	workspaceIDs := make([]string, numWorkspaces)
	workspacePaths := make([]string, numWorkspaces)
	workspaceNames := make([]string, numWorkspaces)

	for i := 0; i < numWorkspaces; i++ {
		path := filepath.Join(tempDir, fmt.Sprintf("workspace%d", i))
		name := fmt.Sprintf("Workspace %d", i)
		err = os.MkdirAll(path, 0755)
		require.NoError(t, err)

		workspaceIDs[i] = registerWorkspace(t, app, name, path)
		workspacePaths[i] = path
		workspaceNames[i] = name
	}

	// Concurrently create projects in each workspace (simplified test)
	var wg sync.WaitGroup
	projectIDs := make([]string, numWorkspaces)
	errorsChan := make(chan error, numWorkspaces)

	for i := 0; i < numWorkspaces; i++ {
		wg.Add(1)
		go func(wsIndex int) {
			defer wg.Done()

			wsID := workspaceIDs[wsIndex]
			wsPath := workspacePaths[wsIndex]
			wsName := workspaceNames[wsIndex]

			// Create project (with error handling)
			projectID, err := createProjectWithError(app, wsID, wsPath, wsName, fmt.Sprintf("Project %d", wsIndex), "Test project")
			if err != nil {
				errorsChan <- fmt.Errorf("workspace %d project creation: %w", wsIndex, err)
				return
			}
			projectIDs[wsIndex] = projectID
		}(i)
	}

	wg.Wait()
	close(errorsChan)

	// Check for errors
	var errors []error
	for err := range errorsChan {
		errors = append(errors, err)
	}
	assert.Empty(t, errors, "Should have no errors during concurrent workspace access")

	// Verify all projects were created
	for i := 0; i < numWorkspaces; i++ {
		assert.NotEmpty(t, projectIDs[i], "Workspace %d should have a project", i)
	}
}

// TestE2E_WorkspaceHeaderInjection tests that workspace headers are properly handled
func TestE2E_WorkspaceHeaderInjection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gorev-e2e-headers-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewWorkspaceManager()
	app := setupTestServer(manager)

	// Register workspace
	wsPath := filepath.Join(tempDir, "workspace")
	err = os.MkdirAll(wsPath, 0755)
	require.NoError(t, err)

	wsID := registerWorkspace(t, app, "Test Workspace", wsPath)

	// Test missing headers
	req := httptest.NewRequest("GET", "/api/tasks", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	// Should still work but use default behavior (no workspace context)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test with headers
	req = httptest.NewRequest("GET", "/api/tasks", nil)
	req.Header.Set("X-Workspace-Id", wsID)
	req.Header.Set("X-Workspace-Path", wsPath)
	req.Header.Set("X-Workspace-Name", "Test Workspace")

	resp, err = app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test with incorrect workspace ID
	req = httptest.NewRequest("GET", "/api/tasks", nil)
	req.Header.Set("X-Workspace-Id", "non-existent-workspace-id")
	req.Header.Set("X-Workspace-Path", wsPath)
	req.Header.Set("X-Workspace-Name", "Test Workspace")

	resp, err = app.Test(req, -1)
	require.NoError(t, err)
	// Should handle gracefully
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestE2E_WorkspaceUnregistration tests unregistering a workspace
func TestE2E_WorkspaceUnregistration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gorev-e2e-unreg-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	manager := NewWorkspaceManager()
	app := setupTestServer(manager)

	// Register workspace
	wsPath := filepath.Join(tempDir, "workspace")
	err = os.MkdirAll(wsPath, 0755)
	require.NoError(t, err)

	wsID := registerWorkspace(t, app, "Test Workspace", wsPath)

	// Verify workspace exists
	req := httptest.NewRequest("GET", "/api/workspaces", nil)
	resp, err := app.Test(req, -1)
	require.NoError(t, err)

	var listResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&listResp)
	assert.Equal(t, float64(1), listResp["total"])

	// Unregister workspace
	req = httptest.NewRequest("DELETE", "/api/workspaces/"+wsID, nil)
	resp, err = app.Test(req, -1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify workspace is gone
	req = httptest.NewRequest("GET", "/api/workspaces", nil)
	resp, err = app.Test(req, -1)
	require.NoError(t, err)

	json.NewDecoder(resp.Body).Decode(&listResp)
	assert.Equal(t, float64(0), listResp["total"])
}

// Helper functions

func setupTestServer(manager *WorkspaceManager) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Setup routes (simplified for testing)
	apiGroup := app.Group("/api")

	apiGroup.Post("/workspaces/register", func(c *fiber.Ctx) error {
		var req WorkspaceRegistration
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"success": false, "error": err.Error()})
		}

		workspace, err := manager.RegisterWorkspace(req.Path, req.Name)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"success":      true,
			"workspace_id": workspace.ID,
		})
	})

	apiGroup.Get("/workspaces", func(c *fiber.Ctx) error {
		workspaces := manager.ListWorkspaces()

		// Convert to WorkspaceInfo for API response
		infos := make([]*WorkspaceInfo, 0, len(workspaces))
		for _, ws := range workspaces {
			infos = append(infos, ws.ToWorkspaceInfo())
		}

		return c.JSON(fiber.Map{
			"success":    true,
			"workspaces": infos,
			"total":      len(infos),
		})
	})

	apiGroup.Delete("/workspaces/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := manager.UnregisterWorkspace(id); err != nil {
			return c.Status(500).JSON(fiber.Map{"success": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"success": true})
	})

	// Simplified task endpoints for testing
	apiGroup.Get("/tasks", func(c *fiber.Ctx) error {
		wsID := c.Get("X-Workspace-Id")
		if wsID == "" {
			return c.JSON(fiber.Map{"tasks": []interface{}{}})
		}

		workspace, err := manager.GetWorkspace(wsID)
		if err != nil {
			return c.JSON(fiber.Map{"tasks": []interface{}{}})
		}

		wsContext := workspace.(*WorkspaceContext)
		filters := make(map[string]interface{})
		tasks, err := wsContext.IsYonetici.GorevListele(filters)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"tasks": tasks})
	})

	apiGroup.Post("/projects", func(c *fiber.Ctx) error {
		var req struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		wsID := c.Get("X-Workspace-Id")
		workspace, err := manager.GetWorkspace(wsID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Workspace not found"})
		}

		wsContext := workspace.(*WorkspaceContext)
		project, err := wsContext.IsYonetici.ProjeOlustur(req.Name, req.Description)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"id": project.ID})
	})

	apiGroup.Post("/tasks", func(c *fiber.Ctx) error {
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			ProjectID   string `json:"project_id"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		wsID := c.Get("X-Workspace-Id")
		workspace, err := manager.GetWorkspace(wsID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Workspace not found"})
		}

		wsContext := workspace.(*WorkspaceContext)

		// Use Feature template for testing (it has minimal required fields)
		templates, _ := wsContext.IsYonetici.TemplateListele("")
		if len(templates) == 0 {
			return c.Status(500).JSON(fiber.Map{"error": "No templates available"})
		}

		// Find a feature template or use any available template
		var templateID string
		for _, tmpl := range templates {
			if tmpl.Alias == "feature" || tmpl.Category == "Özellik" {
				templateID = tmpl.ID
				break
			}
		}
		if templateID == "" {
			templateID = templates[0].ID
		}

		values := map[string]string{
			"baslik":    req.Title,
			"aciklama":  req.Description,
			"oncelik":   "orta",
			"proje_id":  req.ProjectID,
			"etiketler": "test",
		}

		task, err := wsContext.IsYonetici.TemplatedenGorevOlustur(templateID, values)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"id": task.ID})
	})

	return app
}

func registerWorkspace(t *testing.T, app *fiber.App, name, path string) string {
	payload := map[string]string{"name": name, "path": path}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/workspaces/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	require.Equal(t, true, result["success"])
	require.NotEmpty(t, result["workspace_id"])

	return result["workspace_id"].(string)
}

func createProject(t *testing.T, app *fiber.App, wsID, wsPath, wsName, projectName, description string) string {
	projectID, err := createProjectWithError(app, wsID, wsPath, wsName, projectName, description)
	require.NoError(t, err)
	return projectID
}

func createProjectWithError(app *fiber.App, wsID, wsPath, wsName, projectName, description string) (string, error) {
	payload := map[string]string{
		"name":        projectName,
		"description": description,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Workspace-Id", wsID)
	req.Header.Set("X-Workspace-Path", wsPath)
	req.Header.Set("X-Workspace-Name", wsName)

	resp, err := app.Test(req, -1)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	id, ok := result["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("project ID is empty or invalid")
	}

	return id, nil
}

func createTask(t *testing.T, app *fiber.App, wsID, wsPath, wsName, projectID, title string) string {
	taskID, err := createTaskWithError(app, wsID, wsPath, wsName, projectID, title)
	require.NoError(t, err)
	return taskID
}

func createTaskWithError(app *fiber.App, wsID, wsPath, wsName, projectID, title string) (string, error) {
	payload := map[string]string{
		"title":       title,
		"description": "Test task",
		"project_id":  projectID,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Workspace-Id", wsID)
	req.Header.Set("X-Workspace-Path", wsPath)
	req.Header.Set("X-Workspace-Name", wsName)

	resp, err := app.Test(req, -1)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result["id"] == "" {
		return "", fmt.Errorf("task ID is empty")
	}

	// Add small delay to avoid overwhelming SQLite
	time.Sleep(10 * time.Millisecond)

	return result["id"], nil
}

func listTasks(t *testing.T, app *fiber.App, wsID, wsPath, wsName string) []map[string]interface{} {
	req := httptest.NewRequest("GET", "/api/tasks", nil)
	req.Header.Set("X-Workspace-Id", wsID)
	req.Header.Set("X-Workspace-Path", wsPath)
	req.Header.Set("X-Workspace-Name", wsName)

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	tasks, ok := result["tasks"].([]interface{})
	require.True(t, ok)

	typedTasks := make([]map[string]interface{}, len(tasks))
	for i, task := range tasks {
		typedTasks[i] = task.(map[string]interface{})
	}

	return typedTasks
}

func getTaskIDs(tasks []map[string]interface{}) []string {
	ids := make([]string, len(tasks))
	for i, task := range tasks {
		ids[i] = task["id"].(string)
	}
	return ids
}
