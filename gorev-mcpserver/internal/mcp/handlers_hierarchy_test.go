package mcp

import (
	"fmt"
	"strings"
	"testing"

	"github.com/msenol/gorev/internal/gorev"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test for gorevHiyerarsiYazdirVeIsaretle
func TestGorevHiyerarsiYazdirVeIsaretle(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test tasks
	parentTask := &gorev.Gorev{
		ID:      "parent-123",
		Baslik:  "Parent Task",
		Durum:   "beklemede",
		Oncelik: "yuksek",
		ProjeID: "proj-1",
	}

	childTask1 := &gorev.Gorev{
		ID:       "child-1",
		Baslik:   "Child 1",
		Durum:    "devam_ediyor",
		Oncelik:  "orta",
		ProjeID:  "proj-1",
		ParentID: parentTask.ID,
	}

	childTask2 := &gorev.Gorev{
		ID:       "child-2",
		Baslik:   "Child 2",
		Durum:    "tamamlandi",
		Oncelik:  "dusuk",
		ProjeID:  "proj-1",
		ParentID: parentTask.ID,
	}

	grandchildTask := &gorev.Gorev{
		ID:       "grandchild-1",
		Baslik:   "Grandchild",
		Durum:    "beklemede",
		Oncelik:  "orta",
		ProjeID:  "proj-1",
		ParentID: childTask1.ID,
	}

	gorevMap := map[string]*gorev.Gorev{
		parentTask.ID:     parentTask,
		childTask1.ID:     childTask1,
		childTask2.ID:     childTask2,
		grandchildTask.ID: grandchildTask,
	}

	tests := []struct {
		name              string
		gorev             *gorev.Gorev
		gorevMap          map[string]*gorev.Gorev
		seviye            int
		projeGoster       bool
		initialShownIDs   map[string]bool
		expectShownIDs    []string
		expectNotShownIDs []string
	}{
		{
			name:            "mark parent and all children",
			gorev:           parentTask,
			gorevMap:        gorevMap,
			seviye:          0,
			projeGoster:     true,
			initialShownIDs: make(map[string]bool),
			expectShownIDs: []string{
				parentTask.ID,
				childTask1.ID,
				childTask2.ID,
				grandchildTask.ID,
			},
		},
		{
			name:        "mark only specific subtree",
			gorev:       childTask1,
			gorevMap:    gorevMap,
			seviye:      1,
			projeGoster: false,
			initialShownIDs: map[string]bool{
				parentTask.ID: true, // Parent already shown
			},
			expectShownIDs: []string{
				childTask1.ID,
				grandchildTask.ID,
			},
			expectNotShownIDs: []string{
				childTask2.ID, // Sibling not shown
			},
		},
		{
			name:     "skip already shown tasks",
			gorev:    parentTask,
			gorevMap: gorevMap,
			seviye:   0,
			initialShownIDs: map[string]bool{
				parentTask.ID: true,
				childTask1.ID: true,
			},
			expectShownIDs: []string{
				childTask2.ID,     // Only unshown child
				grandchildTask.ID, // Grandchild of shown parent
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shownIDs := make(map[string]bool)
			for k, v := range tt.initialShownIDs {
				shownIDs[k] = v
			}

			result := handlers.gorevHiyerarsiYazdirVeIsaretle(
				tt.gorev,
				tt.gorevMap,
				tt.seviye,
				tt.projeGoster,
				shownIDs,
			)

			// Check that output contains task information
			assert.Contains(t, result, tt.gorev.Baslik)

			// Verify shown IDs
			for _, id := range tt.expectShownIDs {
				assert.True(t, shownIDs[id], "Expected %s to be marked as shown", id)
			}

			// Verify not shown IDs
			for _, id := range tt.expectNotShownIDs {
				assert.False(t, shownIDs[id], "Expected %s to NOT be marked as shown", id)
			}
		})
	}
}

// Test for gorevHiyerarsiYazdirInternal
func TestGorevHiyerarsiYazdirInternal(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a complex hierarchy
	rootTask := &gorev.Gorev{
		ID:                            "root-task",
		Baslik:                        "Root Task with Dependencies",
		Durum:                         "devam_ediyor",
		Oncelik:                       "yuksek",
		ProjeID:                       "proj-main",
		BagimliGorevSayisi:            3,
		TamamlanmamisBagimlilikSayisi: 2,
		Etiketler:                     []*gorev.Etiket{{Isim: "important"}, {Isim: "milestone"}},
	}

	completedChild := &gorev.Gorev{
		ID:        "completed-child",
		Baslik:    "Completed Subtask",
		Durum:     "tamamlandi",
		Oncelik:   "orta",
		ProjeID:   "proj-main",
		ParentID:  rootTask.ID,
		Etiketler: []*gorev.Etiket{{Isim: "done"}},
	}

	inProgressChild := &gorev.Gorev{
		ID:                            "progress-child",
		Baslik:                        "In Progress Subtask",
		Durum:                         "devam_ediyor",
		Oncelik:                       "yuksek",
		ProjeID:                       "proj-main",
		ParentID:                      rootTask.ID,
		BagimliGorevSayisi:            1,
		TamamlanmamisBagimlilikSayisi: 1,
	}

	deepChild := &gorev.Gorev{
		ID:       "deep-child",
		Baslik:   "Deep Nested Task",
		Durum:    "beklemede",
		Oncelik:  "dusuk",
		ProjeID:  "proj-main",
		ParentID: inProgressChild.ID,
	}

	gorevMap := map[string]*gorev.Gorev{
		rootTask.ID:        rootTask,
		completedChild.ID:  completedChild,
		inProgressChild.ID: inProgressChild,
		deepChild.ID:       deepChild,
	}

	tests := []struct {
		name            string
		gorev           *gorev.Gorev
		seviye          int
		projeGoster     bool
		expectPrefix    bool
		expectCompleted bool
		expectTags      bool
		expectDeps      bool
		expectChildren  int
	}{
		{
			name:            "root task with all features",
			gorev:           rootTask,
			seviye:          0,
			projeGoster:     true,
			expectPrefix:    false,
			expectCompleted: false,
			expectTags:      true,
			expectDeps:      true,
			expectChildren:  2,
		},
		{
			name:            "completed child task",
			gorev:           completedChild,
			seviye:          1,
			projeGoster:     false,
			expectPrefix:    true,
			expectCompleted: true,
			expectTags:      true,
			expectDeps:      false,
			expectChildren:  0,
		},
		{
			name:            "in-progress with dependencies",
			gorev:           inProgressChild,
			seviye:          1,
			projeGoster:     false,
			expectPrefix:    true,
			expectCompleted: false,
			expectTags:      false,
			expectDeps:      true,
			expectChildren:  1,
		},
		{
			name:            "deep nested task",
			gorev:           deepChild,
			seviye:          3,
			projeGoster:     true,
			expectPrefix:    true,
			expectCompleted: false,
			expectTags:      false,
			expectDeps:      false,
			expectChildren:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shownIDs := make(map[string]bool)
			result := handlers.gorevHiyerarsiYazdirInternal(
				tt.gorev,
				gorevMap,
				tt.seviye,
				tt.projeGoster,
				shownIDs,
			)

			// Check basic structure
			lines := strings.Split(result, "\n")
			assert.NotEmpty(t, lines)

			// Check indentation
			if tt.seviye > 0 {
				assert.True(t, strings.HasPrefix(lines[0], strings.Repeat("  ", tt.seviye)))
			}

			// Check prefix
			if tt.expectPrefix {
				assert.Contains(t, result, "└─")
			}

			// Check completed status
			if tt.expectCompleted {
				assert.Contains(t, result, "~~")
			}

			// Check tags
			if tt.expectTags && len(tt.gorev.Etiketler) > 0 {
				for _, tag := range tt.gorev.Etiketler {
					assert.Contains(t, result, tag)
				}
			}

			// Check dependencies
			if tt.expectDeps {
				assert.Contains(t, result, "🔒")
			}

			// Check children count
			childCount := 0
			for _, line := range lines {
				if strings.Contains(line, "└─") && !strings.Contains(line, tt.gorev.Baslik) {
					childCount++
				}
			}
			// Note: This is approximate as children may have their own children
			if tt.expectChildren > 0 {
				assert.True(t, childCount >= tt.expectChildren)
			}
		})
	}
}

// Test pagination edge cases with hierarchy
func TestHierarchyWithPagination(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur()
	require.NoError(t, err)

	// Create a large hierarchy
	var taskIDs []string

	// Create multiple root tasks
	for i := 0; i < 10; i++ {
		rootResult, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
			"template_id": "feature_request",
			"degerler": map[string]interface{}{
				"baslik":    fmt.Sprintf("Root Task %d", i),
				"aciklama":  "Root task",
				"oncelik":   "orta",
				"modul":     "test",
				"kullanici": "user",
			},
		})
		rootID := extractTaskIDFromText(getResultText(rootResult))
		taskIDs = append(taskIDs, rootID)

		// Create subtasks for some roots
		if i%2 == 0 {
			for j := 0; j < 3; j++ {
				subResult, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
					"parent_id": rootID,
					"baslik":    fmt.Sprintf("Subtask %d-%d", i, j),
				})
				subID := extractTaskIDFromText(getResultText(subResult))

				// Create grandchildren for first subtask
				if j == 0 {
					for k := 0; k < 2; k++ {
						handlers.GorevAltGorevOlustur(map[string]interface{}{
							"parent_id": subID,
							"baslik":    fmt.Sprintf("Grandchild %d-%d-%d", i, j, k),
						})
					}
				}
			}
		}
	}

	// Test listing with pagination
	tests := []struct {
		name           string
		limit          int
		offset         int
		expectRootMin  int
		expectTotalMin int
	}{
		{
			name:           "first page",
			limit:          5,
			offset:         0,
			expectRootMin:  5,
			expectTotalMin: 5, // At least the roots
		},
		{
			name:           "second page",
			limit:          5,
			offset:         5,
			expectRootMin:  5,
			expectTotalMin: 5,
		},
		{
			name:           "page with subtasks",
			limit:          3,
			offset:         0,
			expectRootMin:  3,
			expectTotalMin: 6, // Roots 0, 2 have subtasks
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.GorevListele(map[string]interface{}{
				"limit":  float64(tt.limit),
				"offset": float64(tt.offset),
			})
			require.NoError(t, err)

			text := getResultText(result)
			lines := strings.Split(text, "\n")

			rootCount := 0
			totalCount := 0
			for _, line := range lines {
				if strings.Contains(line, "Root Task") {
					rootCount++
				}
				if strings.Contains(line, "Task") || strings.Contains(line, "child") {
					totalCount++
				}
			}

			assert.GreaterOrEqual(t, rootCount, tt.expectRootMin)
			assert.GreaterOrEqual(t, totalCount, tt.expectTotalMin)
		})
	}
}

// Test circular dependency prevention in hierarchy
func TestCircularDependencyPrevention(t *testing.T) {
	_, handlers, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Initialize templates
	err := handlers.isYonetici.VeriYonetici().VarsayilanTemplateleriOlustur()
	require.NoError(t, err)

	// Create a chain of tasks
	task1Result, _ := handlers.TemplatedenGorevOlustur(map[string]interface{}{
		"template_id": "feature_request",
		"degerler": map[string]interface{}{
			"baslik":    "Task 1",
			"aciklama":  "First",
			"oncelik":   "orta",
			"modul":     "test",
			"kullanici": "user",
		},
	})
	task1ID := extractTaskIDFromText(getResultText(task1Result))

	task2Result, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": task1ID,
		"baslik":    "Task 2",
	})
	task2ID := extractTaskIDFromText(getResultText(task2Result))

	task3Result, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": task2ID,
		"baslik":    "Task 3",
	})
	task3ID := extractTaskIDFromText(getResultText(task3Result))

	task4Result, _ := handlers.GorevAltGorevOlustur(map[string]interface{}{
		"parent_id": task3ID,
		"baslik":    "Task 4",
	})
	task4ID := extractTaskIDFromText(getResultText(task4Result))

	// Test various circular dependency scenarios
	tests := []struct {
		name        string
		gorevID     string
		newParentID string
		shouldFail  bool
		errorMsg    string
	}{
		{
			name:        "direct circular - parent to child",
			gorevID:     task1ID,
			newParentID: task2ID,
			shouldFail:  true,
			errorMsg:    "dairesel bağımlılık",
		},
		{
			name:        "indirect circular - grandparent to grandchild",
			gorevID:     task1ID,
			newParentID: task3ID,
			shouldFail:  true,
			errorMsg:    "dairesel bağımlılık",
		},
		{
			name:        "deep circular - root to deep descendant",
			gorevID:     task1ID,
			newParentID: task4ID,
			shouldFail:  true,
			errorMsg:    "dairesel bağımlılık",
		},
		{
			name:        "valid move - sibling reparenting",
			gorevID:     task3ID,
			newParentID: task1ID,
			shouldFail:  false,
		},
		{
			name:        "valid move - to unrelated task",
			gorevID:     task4ID,
			newParentID: "",
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handlers.GorevUstDegistir(map[string]interface{}{
				"gorev_id":       tt.gorevID,
				"yeni_parent_id": tt.newParentID,
			})
			require.NoError(t, err)

			text := getResultText(result)
			if tt.shouldFail {
				assert.Contains(t, text, tt.errorMsg)
			} else {
				assert.Contains(t, text, "✓")
			}
		})
	}
}
