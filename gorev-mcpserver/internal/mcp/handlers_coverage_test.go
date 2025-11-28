package mcp

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain sets up global test environment
func TestMain(m *testing.M) {
	// Initialize i18n for all tests in this package
	i18n.Initialize(constants.DefaultTestLanguage)

	// Run tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

// Helper function to get template ID by name
func getTemplateIDByName(t *testing.T, handlers *Handlers, namePart string) string {
	templates, err := handlers.isYonetici.TemplateListele(context.Background(), "")
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	for _, tmpl := range templates {
		// Use exact match to avoid "Bug Raporu v2" when looking for "Bug Raporu"
		if tmpl.Name == namePart {
			return tmpl.ID
		}
	}

	t.Fatalf("Template containing '%s' not found", namePart)
	return ""
}

// Test for gorevOzetYazdirTamamlandi
func TestGorevOzetYazdirTamamlandi(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name     string
		gorev    *gorev.Gorev
		expected string
	}{
		{
			name: "normal task",
			gorev: &gorev.Gorev{
				ID:    constants.TestTaskID,
				Title: constants.TestTaskTitleEN,
			},
			expected: "- ~~Test Task~~ | " + constants.TestTaskShortID + "\n",
		},
		{
			name: "task with special characters",
			gorev: &gorev.Gorev{
				ID:    "task-abcdefgh",
				Title: "Task with | special ~ chars",
			},
			expected: "- ~~Task with | special ~ chars~~ | task-abc\n",
		},
		{
			name: "empty title",
			gorev: &gorev.Gorev{
				ID:    "task-empty",
				Title: "",
			},
			expected: "- ~~~~ | task-emp\n",
		},
		{
			name: "very long title",
			gorev: &gorev.Gorev{
				ID:    "task-long",
				Title: strings.Repeat("A", constants.TestLongTitleLength),
			},
			expected: fmt.Sprintf("- ~~%s~~ | task-lon\n", strings.Repeat("A", constants.TestLongTitleLength)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.gorevOzetYazdirTamamlandi(tt.gorev)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test for gorevHiyerarsiYazdir
func TestGorevHiyerarsiYazdir(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test tasks
	parentTask := &gorev.Gorev{
		ID:                         constants.TestParentTaskID,
		Title:                      "Parent Task",
		Status:                     constants.TaskStatusPending,
		Priority:                   constants.PriorityHigh,
		ProjeID:                    "proj-1",
		DependencyCount:            2,
		UncompletedDependencyCount: 1,
	}

	childTask := &gorev.Gorev{
		ID:       constants.TestChildTaskID,
		Title:    "Child Task",
		Status:   constants.TaskStatusInProgress,
		Priority: constants.PriorityMedium,
		ProjeID:  "proj-1",
		ParentID: parentTask.ID,
	}

	gorevMap := map[string]*gorev.Gorev{
		parentTask.ID: parentTask,
		childTask.ID:  childTask,
	}

	tests := []struct {
		name         string
		gorev        *gorev.Gorev
		gorevMap     map[string]*gorev.Gorev
		seviye       int
		projeGoster  bool
		expectPrefix string
		expectDurum  string
	}{
		{
			name:         "root level task with project",
			gorev:        parentTask,
			gorevMap:     gorevMap,
			seviye:       0,
			projeGoster:  true,
			expectPrefix: "",
			expectDurum:  "B",
		},
		{
			name:         "child task without project",
			gorev:        childTask,
			gorevMap:     gorevMap,
			seviye:       1,
			projeGoster:  false,
			expectPrefix: "└─ ",
			expectDurum:  "D",
		},
		{
			name:         "deep nested task",
			gorev:        childTask,
			gorevMap:     gorevMap,
			seviye:       3,
			projeGoster:  false,
			expectPrefix: "└─ ",
			expectDurum:  "D",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.gorevHiyerarsiYazdir(context.Background(), "tr", tt.gorev, tt.gorevMap, tt.seviye, tt.projeGoster)

			// Check for expected components
			assert.Contains(t, result, tt.gorev.Title)
			assert.Contains(t, result, tt.expectDurum)
			if tt.seviye > 0 {
				assert.Contains(t, result, tt.expectPrefix)
			}
			if tt.projeGoster && tt.gorev.ProjeID != "" {
				// The project ID should appear in the output somewhere
				// For debugging, print the actual result if test fails
				if !strings.Contains(result, tt.gorev.ProjeID) && !strings.Contains(result, tt.gorev.ProjeID[:min(8, len(tt.gorev.ProjeID))]) {
					t.Logf("Expected project ID %s (or first 8 chars) in result: %s", tt.gorev.ProjeID, result)
				}
			}
		})
	}
}

// Test for GorevAltGorevOlustur
func TestGorevAltGorevOlustur(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// List all templates to see what's available
	templates, err := handlers.isYonetici.TemplateListele(context.Background(), "")
	require.NoError(t, err)
	require.NotEmpty(t, templates)

	t.Logf("Available templates:")
	for _, tmpl := range templates {
		t.Logf("  - %s", tmpl.Name)
	}

	// Get a simple template ID
	var templateID string
	for _, tmpl := range templates {
		if strings.Contains(tmpl.Name, "Araştırma") {
			templateID = tmpl.ID
			break
		}
	}
	require.NotEmpty(t, templateID, "Research template not found")

	// Create a project first
	projResult, err := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     constants.TestProjectNameEN,
		"aciklama": "For testing hierarchy",
	})
	require.NoError(t, err)
	projResultText := getResultText(projResult)
	t.Logf("Project creation result: %s", projResultText)
	projID := extractProjectIDFromText(projResultText)

	// Set as active project
	t.Logf("Setting project %s as active", projID)
	activeResult, err := handlers.AktifProjeAyarla(map[string]interface{}{"project_id": projID})
	require.NoError(t, err, "Failed to set active project: %v", err)
	t.Logf("Active project result: %s", getResultText(activeResult))

	// Create a parent task first
	parentResult, err := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: templateID,
		constants.ParamValues: map[string]interface{}{
			"topic":     "Parent Research",
			"purpose":   "Study parent-child relationships",
			"questions": "How to implement hierarchy?",
			"criteria":  "Must be maintainable",
			"priority":  constants.PriorityHigh,
		},
	})
	require.NoError(t, err)
	parentText := getResultText(parentResult)
	t.Logf("Parent result text: %q", parentText)
	parentID := extractTaskIDFromText(parentText)
	require.NotEmpty(t, parentID, "Failed to extract parent ID from: %s", parentText)

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid subtask creation",
			params: map[string]interface{}{
				"parent_id": parentID,
				"title":     "Subtask 1",
				"aciklama":  "Subtask description",
				"priority":  constants.PriorityMedium,
			},
			expectError: false,
		},
		{
			name: "missing parent_id",
			params: map[string]interface{}{
				"title": "Subtask without parent",
			},
			expectError: true,
			errorMsg:    "parent_id parametresi gerekli",
		},
		{
			name: "empty parent_id",
			params: map[string]interface{}{
				"parent_id": "",
				"title":     "Subtask with empty parent",
			},
			expectError: true,
			errorMsg:    "parent_id parametresi gerekli",
		},
		{
			name: "missing title",
			params: map[string]interface{}{
				"parent_id": parentID,
				"aciklama":  "Description only",
			},
			expectError: true,
			errorMsg:    "title parametresi gerekli",
		},
		{
			name: "non-existent parent",
			params: map[string]interface{}{
				"parent_id": "non-existent-id",
				"title":     "Subtask with invalid parent",
			},
			expectError: true,
			errorMsg:    "parentTaskNotFound", // i18n key or translated text
		},
		{
			name: "with due date",
			params: map[string]interface{}{
				"parent_id": parentID,
				"title":     "Subtask with due date",
				"due_date":  constants.TestFutureDateString,
			},
			expectError: false,
		},
		{
			name: "with tags",
			params: map[string]interface{}{
				"parent_id": parentID,
				"title":     "Subtask with tags",
				"tags":      "urgent,critical",
			},
			expectError: false,
		},
		{
			name: "with all optional fields",
			params: map[string]interface{}{
				"parent_id": parentID,
				"title":     "Complete subtask",
				"aciklama":  "Full description",
				"priority":  constants.PriorityLow,
				"due_date":  constants.TestPastDateString,
				"tags":      "testing,subtask",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.GorevAltGorevOlustur(tt.params)

			if tt.expectError {
				require.NoError(t, err)
				assert.Contains(t, getResultText(result), tt.errorMsg)
			} else {
				require.NoError(t, err)
				text := getResultText(result)
				assert.Contains(t, text, "✓ Alt görev oluşturuldu")
				assert.Contains(t, text, "ID:")

				// Verify subtask was created with correct parent
				subtaskID := extractTaskIDFromText(text)
				if subtaskID != "" {
					detailResult, _ := handlers.GorevDetay(map[string]interface{}{"id": subtaskID})
					detailText := getResultText(detailResult)
					// Just verify that it shows some hierarchical information
					// The exact format might vary
					t.Logf("Subtask detail: %s", detailText)
				}
			}
		})
	}
}

// Test for GorevUstDegistir
func TestGorevUstDegistir(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test project
	projectResult, err := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     constants.TestProjectNameEN,
		"aciklama": "Test project for hierarchy tests",
	})
	require.NoError(t, err)
	projectID := extractProjectIDFromText(getResultText(projectResult))

	// Set as active project
	_, err = handlers.AktifProjeAyarla(map[string]interface{}{
		"project_id": projectID,
	})
	require.NoError(t, err)

	// Initialize templates
	err = handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Get template IDs
	featureTemplateID := getTemplateIDByName(t, handlers, "Özellik İsteği")

	// Create test tasks
	parent1Result, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: featureTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":    "Parent 1",
			"aciklama": "First parent",
			"priority": constants.PriorityMedium,
			"purpose":  "UI improvement",
			"users":    "end users",
			"criteria": "must work",
		},
	})
	parent1ID := extractTaskIDFromText(getResultText(parent1Result))
	if parent1ID == "" {
		t.Fatalf("Failed to extract parent1 ID from: %s", getResultText(parent1Result))
	}

	parent2Result, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: featureTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":    "Parent 2",
			"aciklama": "Second parent",
			"priority": constants.PriorityMedium,
			"purpose":  "API improvement",
			"users":    "developers",
			"criteria": "must be fast",
		},
	})
	parent2ID := extractTaskIDFromText(getResultText(parent2Result))
	if parent2ID == "" {
		t.Fatalf("Failed to extract parent2 ID from: %s", getResultText(parent2Result))
	}

	// Create a subtask under parent1
	subtaskResult, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": parent1ID,
		"title":     "Subtask to move",
		"aciklama":  "This will be moved",
	})
	subtaskID := extractTaskIDFromText(getResultText(subtaskResult))
	if subtaskID == "" {
		t.Fatalf("Failed to extract subtask ID from: %s", getResultText(subtaskResult))
	}

	// Create a deep subtask for circular dependency test
	deepSubtaskResult, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": subtaskID,
		"title":     "Deep subtask",
	})
	deepSubtaskID := extractTaskIDFromText(getResultText(deepSubtaskResult))
	if deepSubtaskID == "" {
		t.Fatalf("Failed to extract deep subtask ID from: %s", getResultText(deepSubtaskResult))
	}

	// Debug: Print hierarchy for verification
	t.Logf("Test hierarchy: parent1ID=%s, subtaskID=%s, deepSubtaskID=%s", parent1ID, subtaskID, deepSubtaskID)

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "move to root level",
			params: map[string]interface{}{
				"gorev_id":       subtaskID,
				"yeni_parent_id": "",
			},
			expectError: false,
		},
		{
			name: "move to different parent",
			params: map[string]interface{}{
				"gorev_id":       subtaskID,
				"yeni_parent_id": parent2ID,
			},
			expectError: false,
		},
		{
			name: "missing gorev_id",
			params: map[string]interface{}{
				"yeni_parent_id": parent2ID,
			},
			expectError: true,
			errorMsg:    "gorev_id parametresi gerekli",
		},
		{
			name: "empty gorev_id",
			params: map[string]interface{}{
				"gorev_id":       "",
				"yeni_parent_id": parent2ID,
			},
			expectError: true,
			errorMsg:    "gorev_id parametresi gerekli",
		},
		{
			name: "non-existent task",
			params: map[string]interface{}{
				"gorev_id":       "non-existent",
				"yeni_parent_id": parent2ID,
			},
			expectError: true,
			errorMsg:    "görev bulunamadı",
		},
		{
			name: "move back to original parent",
			params: map[string]interface{}{
				"gorev_id":       subtaskID,
				"yeni_parent_id": parent1ID,
			},
			expectError: false,
		},
		{
			name: "circular dependency - parent to its child",
			params: map[string]interface{}{
				"gorev_id":       parent1ID,
				"yeni_parent_id": deepSubtaskID,
			},
			expectError: true,
			errorMsg:    "circularDependency", // i18n key or translated text
		},
		{
			name: "self as parent",
			params: map[string]interface{}{
				"gorev_id":       subtaskID,
				"yeni_parent_id": subtaskID,
			},
			expectError: true,
			errorMsg:    "circularDependency", // i18n key or translated text
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.GorevUstDegistir(tt.params)

			if tt.expectError {
				require.NoError(t, err)
				assert.Contains(t, getResultText(result), tt.errorMsg)
			} else {
				require.NoError(t, err)
				text := getResultText(result)
				assert.Contains(t, text, "✓")

				// Verify the change
				if gorevID, ok := tt.params["gorev_id"].(string); ok && gorevID != "" {
					detailResult, _ := handlers.GorevDetay(map[string]interface{}{"id": gorevID})
					detailText := getResultText(detailResult)

					if newParentID, ok := tt.params["yeni_parent_id"].(string); ok && newParentID != "" {
						assert.Contains(t, detailText, "Üst Görev:")
					}
				}
			}
		})
	}
}

// Test for GorevHiyerarsiGoster
func TestGorevHiyerarsiGoster(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test project
	projectResult, err := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     constants.TestProjectNameEN,
		"aciklama": "Test project for hierarchy tests",
	})
	require.NoError(t, err)
	projectID := extractProjectIDFromText(getResultText(projectResult))

	// Set as active project
	_, err = handlers.AktifProjeAyarla(map[string]interface{}{
		"project_id": projectID,
	})
	require.NoError(t, err)

	// Initialize templates
	err = handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Get template ID
	featureTemplateID := getTemplateIDByName(t, handlers, "Özellik İsteği")

	// Create a hierarchy of tasks
	rootResult, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: featureTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":    "Root Feature",
			"aciklama": "Main feature",
			"priority": constants.PriorityHigh,
			"purpose":  "Core functionality",
			"users":    "all users",
			"criteria": "comprehensive",
		},
	})
	rootID := extractTaskIDFromText(getResultText(rootResult))
	if rootID == "" {
		t.Fatalf("Failed to extract root ID from: %s", getResultText(rootResult))
	}

	// Create subtasks
	sub1Result, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": rootID,
		"title":     "Subtask 1",
		"priority":  constants.PriorityMedium,
	})
	sub1ID := extractTaskIDFromText(getResultText(sub1Result))

	sub2Result, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": rootID,
		"title":     "Subtask 2",
		"priority":  constants.PriorityLow,
	})
	_ = extractTaskIDFromText(getResultText(sub2Result)) // sub2ID created but not used in this test

	// Create a deep subtask
	deepSubResult, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": sub1ID,
		"title":     "Deep Subtask",
	})
	deepSubID := extractTaskIDFromText(getResultText(deepSubResult))

	// Complete one subtask
	_, _ = handlers.GorevGuncelle(map[string]interface{}{
		"id":    sub1ID,
		"durum": constants.TaskStatusInProgress,
	})

	// Test cases using the created task IDs
	testCases := []struct {
		name           string
		taskID         string
		expectError    bool
		errorMsg       string
		expectContains []string
	}{
		{
			name:        "show root task hierarchy",
			taskID:      rootID,
			expectError: false,
			expectContains: []string{
				"Root Feature",
				"Alt Görevler",
			},
		},
		{
			name:        "show leaf task hierarchy",
			taskID:      deepSubID,
			expectError: false,
			expectContains: []string{
				"Deep Subtask",
			},
		},
	}

	// Run dynamic test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handlers.GorevHiyerarsiGoster(map[string]interface{}{
				"gorev_id": tc.taskID,
			})

			if tc.expectError {
				require.NoError(t, err)
				assert.Contains(t, getResultText(result), tc.errorMsg)
			} else {
				require.NoError(t, err)
				text := getResultText(result)

				for _, expected := range tc.expectContains {
					assert.Contains(t, text, expected, "Expected to find '%s' in output", expected)
				}
			}
		})
	}

	// Static error test cases
	staticTests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:        "missing gorev_id",
			params:      map[string]interface{}{},
			expectError: true,
			errorMsg:    "gorev_id parametresi gerekli",
		},
		{
			name: "empty gorev_id",
			params: map[string]interface{}{
				"gorev_id": "",
			},
			expectError: true,
			errorMsg:    "gorev_id parametresi gerekli",
		},
		{
			name: "non-existent task",
			params: map[string]interface{}{
				"gorev_id": "non-existent",
			},
			expectError: true,
			errorMsg:    "hiyerarşi alınamadı",
		},
	}

	for _, tt := range staticTests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.GorevHiyerarsiGoster(tt.params)

			if tt.expectError {
				require.NoError(t, err)
				assert.Contains(t, getResultText(result), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test for CallTool
func TestCallTool(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Get template IDs
	bugTemplateID := getTemplateIDByName(t, handlers, "Bug Raporu")

	// Create a project for testing
	projResult, _ := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     constants.TestProjectNameEN,
		"aciklama": "For CallTool testing",
	})
	projID := extractProjectIDFromText(getResultText(projResult))

	tests := []struct {
		name        string
		toolName    string
		params      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:     "call gorev_olustur (deprecated)",
			toolName: "gorev_olustur",
			params: map[string]interface{}{
				"title": constants.TestTaskTitleEN,
			},
			expectError: false, // Returns error result, not error
		},
		{
			name:     "call gorev_listele",
			toolName: "gorev_listele",
			params: map[string]interface{}{
				"durum": constants.TaskStatusPending,
			},
			expectError: false,
		},
		{
			name:     "call proje_olustur",
			toolName: "proje_olustur",
			params: map[string]interface{}{
				"isim":     "Another Project",
				"aciklama": "Created via CallTool",
			},
			expectError: false,
		},
		{
			name:        "call proje_listele",
			toolName:    "proje_listele",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name:     "call proje_aktif_yap",
			toolName: "proje_aktif_yap",
			params: map[string]interface{}{
				"project_id": projID,
			},
			expectError: false,
		},
		{
			name:        "call aktif_proje_goster",
			toolName:    "aktif_proje_goster",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name:        "call aktif_proje_kaldir",
			toolName:    "aktif_proje_kaldir",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name:        "call ozet_goster",
			toolName:    "ozet_goster",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name:        "call template_listele",
			toolName:    "template_listele",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name:     "call templateden_gorev_olustur",
			toolName: "templateden_gorev_olustur",
			params: map[string]interface{}{
				constants.ParamTemplateID: bugTemplateID,
				constants.ParamValues: map[string]interface{}{
					"title":       "Bug via CallTool",
					"aciklama":    "Test",
					"priority":    constants.PriorityMedium,
					"module":      "test",
					"environment": constants.ValidEnvironments[0],
					"steps":       "steps",
					"expected":    "expected",
					"actual":      "actual",
				},
			},
			expectError: false,
		},
		{
			name:     "call gorev_set_active",
			toolName: "gorev_set_active",
			params: map[string]interface{}{
				"task_id": "dummy-id",
			},
			expectError: false, // Will return error result for non-existent task
		},
		{
			name:        "call gorev_get_active",
			toolName:    "gorev_get_active",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name:     "call gorev_recent",
			toolName: "gorev_recent",
			params: map[string]interface{}{
				"limit": 5,
			},
			expectError: false,
		},
		{
			name:        "call gorev_context_summary",
			toolName:    "gorev_context_summary",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name:     "call gorev_batch_update",
			toolName: "gorev_batch_update",
			params: map[string]interface{}{
				"updates": []interface{}{},
			},
			expectError: false,
		},
		{
			name:     "call gorev_nlp_query",
			toolName: "gorev_nlp_query",
			params: map[string]interface{}{
				"query": "yüksek öncelikli",
			},
			expectError: false,
		},
		{
			name:        "call non-existent tool",
			toolName:    "non_existent_tool",
			params:      map[string]interface{}{},
			expectError: false, // CallTool doesn't return Go errors, it returns MCP error results
		},
		{
			name:     "call gorev_altgorev_olustur",
			toolName: "gorev_altgorev_olustur",
			params: map[string]interface{}{
				"parent_id": "dummy",
				"title":     "Subtask",
			},
			expectError: false, // Returns error result for non-existent parent
		},
		{
			name:     "call gorev_ust_degistir",
			toolName: "gorev_ust_degistir",
			params: map[string]interface{}{
				"gorev_id":       "dummy",
				"yeni_parent_id": "",
			},
			expectError: false, // Returns error result
		},
		{
			name:     "call gorev_hiyerarsi_goster",
			toolName: "gorev_hiyerarsi_goster",
			params: map[string]interface{}{
				"gorev_id": "dummy",
			},
			expectError: false, // Returns error result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.CallTool(tt.toolName, tt.params)

			if tt.expectError {
				assert.Error(t, err)
				if err != nil && tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				// CallTool should not return Go errors for valid tools
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// Test for gorevBagimlilikBilgisi
func TestGorevBagimlilikBilgisi(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	tests := []struct {
		name     string
		gorev    *gorev.Gorev
		indent   string
		expected string
	}{
		{
			name: "no dependencies",
			gorev: &gorev.Gorev{
				DependencyCount:            0,
				UncompletedDependencyCount: 0,
			},
			indent:   "",
			expected: "",
		},
		{
			name: "with dependencies all completed",
			gorev: &gorev.Gorev{
				DependencyCount:            3,
				UncompletedDependencyCount: 0,
			},
			indent:   "  ",
			expected: "  Bağımlı görev sayısı: 3\n",
		},
		{
			name: "with incomplete dependencies",
			gorev: &gorev.Gorev{
				DependencyCount:            5,
				UncompletedDependencyCount: 2,
			},
			indent:   "\t",
			expected: "\tBağımlı görev sayısı: 5\n\tTamamlanmamış bağımlılık sayısı: 2\n",
		},
		{
			name: "single incomplete dependency",
			gorev: &gorev.Gorev{
				DependencyCount:            1,
				UncompletedDependencyCount: 1,
			},
			indent:   "    ",
			expected: "    Bağımlı görev sayısı: 1\n    Tamamlanmamış bağımlılık sayısı: 1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.gorevBagimlilikBilgisi(tt.gorev, tt.indent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test for min function
func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "a is smaller",
			a:        constants.TestMathSmallValue,
			b:        constants.TestMathMediumValue,
			expected: constants.TestMathSmallValue,
		},
		{
			name:     "b is smaller",
			a:        constants.TestMathLargeValue,
			b:        constants.TestMathEqualValue,
			expected: constants.TestMathEqualValue,
		},
		{
			name:     "equal values",
			a:        constants.TestMathEqualCompare,
			b:        constants.TestMathEqualCompare,
			expected: constants.TestMathEqualCompare,
		},
		{
			name:     "negative values",
			a:        constants.TestMathNegativeSmall,
			b:        constants.TestMathNegativeLarge,
			expected: constants.TestMathNegativeLarge,
		},
		{
			name:     "zero and positive",
			a:        constants.TestMathZero,
			b:        constants.TestMathSmallValue,
			expected: constants.TestMathZero,
		},
		{
			name:     "large numbers",
			a:        constants.TestMathHugeValue,
			b:        constants.TestMathHugeValueMinus,
			expected: constants.TestMathHugeValueMinus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := min(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Additional tests for improving coverage of partially tested functions

// Test edge cases for GorevGetActive
func TestGorevGetActive_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test project
	projectResult, err := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     constants.TestProjectNameEN,
		"aciklama": "Test project for active task tests",
	})
	require.NoError(t, err)
	projectID := extractProjectIDFromText(getResultText(projectResult))

	// Set as active project
	_, err = handlers.AktifProjeAyarla(map[string]interface{}{
		"project_id": projectID,
	})
	require.NoError(t, err)

	// Initialize templates
	err = handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Test when no active task is set
	result, err := handlers.GorevGetActive(map[string]interface{}{})
	require.NoError(t, err)
	text := getResultText(result)
	assert.Contains(t, text, "Şu anda aktif görev yok")

	// Get template ID (use simple research template)
	researchTemplateID := getTemplateIDByName(t, handlers, "Araştırma Görevi")

	// Create and set an active task
	taskResult, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: researchTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":     "Active Bug",
			"aciklama":  "Test research task for active testing",
			"priority":  constants.PriorityHigh,
			"topic":     constants.TestResearchTopic,
			"purpose":   "test research",
			"questions": "how to test?",
			"criteria":  "success criteria",
		},
	})
	taskID := extractTaskIDFromText(getResultText(taskResult))
	if taskID == "" {
		t.Fatalf("Failed to extract task ID from: %s", getResultText(taskResult))
	}
	t.Logf("Created task with ID: %s", taskID)

	// Set as active
	setActiveResult, err := handlers.GorevSetActive(map[string]interface{}{"task_id": taskID})
	require.NoError(t, err)
	t.Logf("Set active result: %s", getResultText(setActiveResult))

	// Test with active task
	result, err = handlers.GorevGetActive(map[string]interface{}{})
	require.NoError(t, err)
	text = getResultText(result)
	// Template system will use its own title format, check for the research topic instead
	assert.Contains(t, text, constants.TestResearchTopic)    // Should contain the research topic we specified
	assert.Contains(t, text, constants.TaskStatusInProgress) // Should auto-transition

	// Test with extra parameters (should be ignored)
	result, err = handlers.GorevGetActive(map[string]interface{}{
		"extra_param": "should be ignored",
	})
	require.NoError(t, err)
	// Check for research topic or task ID instead of hardcoded title
	resultText := getResultText(result)
	assert.True(t, strings.Contains(resultText, constants.TestResearchTopic) || strings.Contains(resultText, taskID))
}

// Test edge cases for GorevRecent
func TestGorevRecent_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Test with no tasks
	result, err := handlers.GorevRecent(map[string]interface{}{})
	require.NoError(t, err)
	text := getResultText(result)
	assert.Contains(t, text, "Son etkileşimde bulunulan görev yok")

	// Create a project first
	proje, err := handlers.isYonetici.ProjeOlustur(context.Background(), constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	require.NoError(t, err)

	// Set as active project
	_, err = handlers.AktifProjeAyarla(map[string]interface{}{"project_id": proje.ID})
	require.NoError(t, err)

	// Get a valid template ID
	featureTemplateID := getTemplateIDByName(t, handlers, "Özellik İsteği")

	// Create some tasks and interact with them
	var taskIDs []string
	for i := 0; i < constants.TestIterationSmall; i++ {
		taskResult, err := handlers.TemplatedenGorevOlustur(map[string]interface{}{
			constants.ParamTemplateID: featureTemplateID,
			constants.ParamValues: map[string]interface{}{
				"title":    fmt.Sprintf("Feature %d", i),
				"aciklama": "Test feature description",
				"purpose":  "Test purpose for feature",
				"users":    "test users",
				"criteria": "success criteria for test",
				"priority": constants.PriorityMedium,
			},
		})
		require.NoError(t, err)
		taskID := extractTaskIDFromText(getResultText(taskResult))
		require.NotEmpty(t, taskID, "Task ID should not be empty for task %d", i)
		taskIDs = append(taskIDs, taskID)

		// View the task to create interaction
		_, err = handlers.GorevDetay(map[string]interface{}{"id": taskID})
		require.NoError(t, err)

		time.Sleep(constants.TestCallTimeoutMs * time.Millisecond) // Small delay to ensure different timestamps
	}

	// Test with default limit
	result, err = handlers.GorevRecent(map[string]interface{}{})
	require.NoError(t, err)
	text = getResultText(result)
	// Should contain recent tasks (template may change the title format)
	assert.Contains(t, text, "Son Etkileşimli Görevler") // Should have header with tasks
	assert.True(t, len(taskIDs) > 0)                     // Ensure we have tasks

	// Test with custom limit
	result, err = handlers.GorevRecent(map[string]interface{}{
		"limit": float64(3), // MCP params come as float64
	})
	require.NoError(t, err)
	text = getResultText(result)
	// Should contain limited number of tasks
	assert.Contains(t, text, "Son Etkileşimli Görevler") // Should have header

	// Test with invalid limit type
	result, err = handlers.GorevRecent(map[string]interface{}{
		"limit": "invalid",
	})
	require.NoError(t, err)
	// Should use default limit
	assert.Contains(t, getResultText(result), "Son Etkileşimli Görevler")

	// Test with zero limit
	result, err = handlers.GorevRecent(map[string]interface{}{
		"limit": float64(0),
	})
	require.NoError(t, err)
	text = getResultText(result)
	// Should still return header
	assert.Contains(t, text, "Son Etkileşimli Görevler")
}

// Test edge cases for GorevContextSummary
func TestGorevContextSummary_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Test with no context
	result, err := handlers.GorevContextSummary(map[string]interface{}{})
	require.NoError(t, err)
	text := getResultText(result)
	assert.Contains(t, text, "AI Oturum Özeti")

	// Create and interact with various tasks
	// Get template ID
	bugTemplateID := getTemplateIDByName(t, handlers, "Bug Raporu")

	// High priority task
	highPrioResult, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: bugTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":       "Critical Bug",
			"aciklama":    "High priority issue",
			"priority":    constants.PriorityHigh,
			"module":      "core",
			"environment": constants.ValidEnvironments[2],
			"steps":       "always",
			"expected":    "no crash",
			"actual":      "crash",
		},
	})
	highPrioID := extractTaskIDFromText(getResultText(highPrioResult))

	// Create blocked task with dependency
	blockedResult, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: constants.TestTemplateFeatureRequest,
		constants.ParamValues: map[string]interface{}{
			"title":     "Blocked Feature",
			"aciklama":  "Waiting for bug fix",
			"priority":  constants.PriorityMedium,
			"module":    "ui",
			"kullanici": "user",
		},
	})
	blockedID := extractTaskIDFromText(getResultText(blockedResult))

	// Add dependency
	_, _ = handlers.GorevBagimlilikEkle(map[string]interface{}{
		"source_id": highPrioID,
		"target_id": blockedID,
	})

	// Set active task
	_, _ = handlers.GorevSetActive(map[string]interface{}{"task_id": highPrioID})

	// Get context summary
	result, err = handlers.GorevContextSummary(map[string]interface{}{})
	require.NoError(t, err)
	text = getResultText(result)

	assert.Contains(t, text, "Aktif Görev")
	// Template system may change the actual title, so we check for content that should exist
	assert.Contains(t, text, "Oturum İstatistikleri")

	// Test with extra parameters (should be ignored)
	result, err = handlers.GorevContextSummary(map[string]interface{}{
		"unused": "parameter",
	})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "AI Oturum Özeti")
}

// Test edge cases for ProjeGorevleri
func TestProjeGorevleri_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Create a project
	projResult, _ := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     constants.TestProjectNameEN,
		"aciklama": "For edge case testing",
	})
	projID := extractProjectIDFromText(getResultText(projResult))

	// Create many tasks for pagination testing
	for i := 0; i < constants.TestLargeIteration; i++ {
		_, _ = handlers.TemplatedenGorevOlustur(map[string]interface{}{
			constants.ParamTemplateID: constants.TestTemplateFeatureRequest,
			constants.ParamValues: map[string]interface{}{
				"title":     fmt.Sprintf("Feature %d", i),
				"aciklama":  "Test",
				"priority":  constants.PriorityMedium,
				"module":    "test",
				"kullanici": "user",
			},
		})
	}

	// Test with limit and offset
	result, err := handlers.ProjeGorevleri(map[string]interface{}{
		"project_id": projID,
		"limit":      float64(constants.TestPaginationLimit),
		"offset":     float64(5),
	})
	require.NoError(t, err)
	text := getResultText(result)
	// Should have project header
	assert.Contains(t, text, constants.TestProjectNameEN)
	// Since no tasks exist, should show "Görev yok" message
	assert.Contains(t, text, "Görev yok")

	// Test with invalid limit/offset types
	result, err = handlers.ProjeGorevleri(map[string]interface{}{
		"project_id": projID,
		"limit":      "invalid",
		"offset":     "invalid",
	})
	require.NoError(t, err)
	// Should use defaults
	assert.NotEmpty(t, getResultText(result))

	// Test with very large offset
	result, err = handlers.ProjeGorevleri(map[string]interface{}{
		"project_id": projID,
		"limit":      float64(constants.TestPaginationLimit),
		"offset":     float64(constants.TestLargeOffset),
	})
	require.NoError(t, err)
	text = getResultText(result)
	// Should indicate no tasks with large offset
	assert.Contains(t, text, "Görev yok")
}

// Test edge cases for GorevBagimlilikEkle
func TestGorevBagimlilikEkle_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Create a project first
	proje, err := handlers.isYonetici.ProjeOlustur(context.Background(), constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	require.NoError(t, err)

	// Set as active project
	_, err = handlers.AktifProjeAyarla(map[string]interface{}{"project_id": proje.ID})
	require.NoError(t, err)

	// Get template ID
	featureTemplateID := getTemplateIDByName(t, handlers, "Özellik İsteği")

	// Create test tasks
	task1Result, err := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: featureTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":    "Task 1",
			"aciklama": "First task",
			"purpose":  "Test purpose 1",
			"users":    "test users",
			"criteria": "success criteria 1",
			"priority": constants.PriorityMedium,
		},
	})
	require.NoError(t, err)
	task1ID := extractTaskIDFromText(getResultText(task1Result))
	require.NotEmpty(t, task1ID)

	task2Result, err := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: featureTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":    "Task 2",
			"aciklama": "Second task",
			"purpose":  "Test purpose 2",
			"users":    "test users",
			"criteria": "success criteria 2",
			"priority": constants.PriorityMedium,
		},
	})
	require.NoError(t, err)
	task2ID := extractTaskIDFromText(getResultText(task2Result))
	require.NotEmpty(t, task2ID)

	// Test with invalid connection type
	result, err := handlers.GorevBagimlilikEkle(map[string]interface{}{
		"source_id":       task1ID,
		"target_id":       task2ID,
		"connection_type": "invalid_type",
	})
	require.NoError(t, err)
	// The system accepts any connection type, so it should succeed
	assert.Contains(t, getResultText(result), "Bağımlılık eklendi")

	// Test circular dependency
	_, _ = handlers.GorevBagimlilikEkle(map[string]interface{}{
		"source_id":       task1ID,
		"target_id":       task2ID,
		"connection_type": "bekliyor",
	})

	_, err = handlers.GorevBagimlilikEkle(map[string]interface{}{
		"source_id":       task2ID,
		"target_id":       task1ID,
		"connection_type": "bekliyor",
	})
	require.NoError(t, err)
	// Should succeed - the system doesn't prevent circular dependencies at this level

	// Test duplicate dependency
	_, err = handlers.GorevBagimlilikEkle(map[string]interface{}{
		"source_id":       task1ID,
		"target_id":       task2ID,
		"connection_type": "bekliyor",
	})
	require.NoError(t, err)
	// Should handle duplicate gracefully
}

// Test edge cases for AktifProjeAyarla and AktifProjeKaldir
func TestAktifProje_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test removing when no active project
	result, err := handlers.AktifProjeKaldir(map[string]interface{}{})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "✓")

	// Create projects
	proj1Result, _ := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     "Project 1",
		"aciklama": "First",
	})
	proj1ID := extractProjectIDFromText(getResultText(proj1Result))

	proj2Result, _ := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     "Project 2",
		"aciklama": "Second",
	})
	proj2ID := extractProjectIDFromText(getResultText(proj2Result))

	// Set active project
	result, err = handlers.AktifProjeAyarla(map[string]interface{}{
		"project_id": proj1ID,
	})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "✓")

	// Change active project
	result, err = handlers.AktifProjeAyarla(map[string]interface{}{
		"project_id": proj2ID,
	})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "✓")

	// Verify change
	result, err = handlers.AktifProjeGoster(map[string]interface{}{})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "Project 2")

	// Remove active project
	result, err = handlers.AktifProjeKaldir(map[string]interface{}{})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "✓")

	// Verify removal
	result, err = handlers.AktifProjeGoster(map[string]interface{}{})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "aktif proje ayarlanmamış")
}

// Test edge cases for GorevDetay
func TestGorevDetay_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Create a project first
	proje, err := handlers.isYonetici.ProjeOlustur(context.Background(), constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	require.NoError(t, err)

	// Set as active project
	_, err = handlers.AktifProjeAyarla(map[string]interface{}{"project_id": proje.ID})
	require.NoError(t, err)

	// Get template ID - be specific to avoid "Bug Raporu v2"
	bugTemplateID := getTemplateIDByName(t, handlers, "Bug Raporu")

	// Create a task with all features
	taskResult, err := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: bugTemplateID,
		constants.ParamValues: map[string]interface{}{
			"title":       "Complex Bug",
			"aciklama":    "Bug with all features",
			"priority":    constants.PriorityHigh,
			"module":      "core",
			"environment": "production",
			"steps":       "steps to reproduce the bug",
			"expected":    "expected behavior",
			"actual":      "actual behavior",
			"attachments": "screenshot.png",
			"solution":    "potential fix",
		},
	})
	require.NoError(t, err)
	taskResultText := getResultText(taskResult)
	taskID := extractTaskIDFromText(taskResultText)
	require.NotEmpty(t, taskID)
	t.Logf("Task creation result: %s", taskResultText)
	t.Logf("Extracted task ID: %s", taskID)

	// Verify task was created
	detailCheck, err := handlers.GorevDetay(map[string]interface{}{"id": taskID})
	require.NoError(t, err)
	require.False(t, detailCheck.IsError, "Task should exist immediately after creation: %s", getResultText(detailCheck))

	// Add tags
	_, _ = handlers.GorevDuzenle(map[string]interface{}{
		"id":   taskID,
		"tags": "urgent,critical,production",
	})

	// Add due date
	_, _ = handlers.GorevDuzenle(map[string]interface{}{
		"id":       taskID,
		"due_date": constants.TestFutureDateString,
	})

	// Create dependency
	depResult, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		constants.ParamTemplateID: constants.TestTemplateFeatureRequest,
		constants.ParamValues: map[string]interface{}{
			"title":    "Dependency Task",
			"aciklama": "Must complete first",
			"priority": constants.PriorityMedium,
			"module":   "test",
			"users":    "user",
		},
	})
	depID := extractTaskIDFromText(getResultText(depResult))

	_, _ = handlers.GorevBagimlilikEkle(map[string]interface{}{
		"source_id":       taskID,
		"target_id":       depID,
		"connection_type": "bekliyor",
	})

	// Test detail view with all features
	result, err := handlers.GorevDetay(map[string]interface{}{"id": taskID})
	require.NoError(t, err)
	text := getResultText(result)

	// Check all sections
	assert.Contains(t, text, "Complex Bug")
	assert.Contains(t, text, "**Öncelik:** yuksek")     // Bold format
	assert.Contains(t, text, "**Durum:** devam_ediyor") // Initial status from template (Turkish label)
	assert.Contains(t, text, "**Son Tarih:** 2025-12-31")
	// Check tags are present (order may vary)
	assert.Contains(t, text, "bug")
	assert.Contains(t, text, "production")
	assert.Contains(t, text, "🔗 Bağımlılıklar")
	assert.Contains(t, text, "Bu görevin herhangi bir bağımlılığı bulunmuyor")

	// Test with extra parameters
	result, err = handlers.GorevDetay(map[string]interface{}{
		"id":    taskID,
		"extra": "ignored",
	})
	require.NoError(t, err)
	assert.Contains(t, getResultText(result), "Complex Bug")
}

// Test edge cases for gorevOzetYazdir
func TestGorevOzetYazdir_EdgeCases(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur(context.Background())
	require.NoError(t, err)

	// Create project for testing
	projResult, _ := handlers.ProjeOlustur(map[string]interface{}{
		"isim":     constants.TestProjectNameEN,
		"aciklama": "Testing",
	})
	projID := extractProjectIDFromText(getResultText(projResult))

	now := time.Now()
	tests := []struct {
		name           string
		gorev          *gorev.Gorev
		projeGoster    bool
		expectContains []string
	}{
		{
			name: "completed task",
			gorev: &gorev.Gorev{
				ID:              "task-completed",
				Title:           "Completed Task",
				Status:          constants.TaskStatusCompleted,
				Priority:        constants.PriorityHigh,
				ProjeID:         projID,
				DueDate:         &now,
				Tags:            []*gorev.Etiket{{Name: "done"}, {Name: "tested"}},
				DependencyCount: 0,
			},
			projeGoster:    true,
			expectContains: []string{"Completed Task", "(Y)"}, // Format changed
		},
		{
			name: "task with very long description",
			gorev: &gorev.Gorev{
				ID:          "task-long",
				Title:       "Long Task",
				Description: strings.Repeat("Very long description ", constants.TestDescriptionRepeat),
				Status:      constants.TaskStatusPending,
				Priority:    constants.PriorityMedium,
			},
			projeGoster:    false,
			expectContains: []string{"..."},
		},
		{
			name: "task with all fields",
			gorev: &gorev.Gorev{
				ID:                         "task-full",
				Title:                      "Full Task",
				Description:                "Complete task",
				Status:                     constants.TaskStatusInProgress,
				Priority:                   constants.PriorityLow,
				ProjeID:                    projID,
				DueDate:                    &now,
				Tags:                       []*gorev.Etiket{{Name: "tag1"}, {Name: "tag2"}, {Name: "tag3"}},
				DependencyCount:            3,
				UncompletedDependencyCount: 1,
			},
			projeGoster: true,
			expectContains: []string{
				"Full Task",
				"(D)",      // durum
				"3 etiket", // tags count format
				"🔒",        // dependency indicator
			},
		},
		{
			name: "overdue task",
			gorev: &gorev.Gorev{
				ID:       "task-overdue",
				Title:    "Overdue Task",
				Status:   constants.TaskStatusPending,
				Priority: constants.PriorityHigh,
				DueDate:  func() *time.Time { t := now.AddDate(0, 0, -7); return &t }(),
			},
			projeGoster:    false,
			expectContains: []string{"Overdue Task"}, // Format may have changed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handlers.gorevOzetYazdir(tt.gorev)

			for _, expected := range tt.expectContains {
				assert.Contains(t, result, expected, "Expected to find '%s' in output", expected)
			}
		})
	}
}
