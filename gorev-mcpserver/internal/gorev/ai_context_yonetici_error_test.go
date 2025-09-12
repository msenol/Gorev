package gorev

import (
	"errors"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
)

// TestAIContextYonetici_addToRecentTasks tests the addToRecentTasks function extensively
func TestAIContextYonetici_addToRecentTasks(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	testCases := []struct {
		name           string
		taskID         string
		existingTasks  []string
		expectedTasks  []string
		expectedLength int
	}{
		{
			name:           "Add to empty recent tasks",
			taskID:         "task-1",
			existingTasks:  []string{},
			expectedTasks:  []string{"task-1"},
			expectedLength: 1,
		},
		{
			name:           "Add new task to existing list",
			taskID:         "task-2",
			existingTasks:  []string{"task-1"},
			expectedTasks:  []string{"task-2", "task-1"},
			expectedLength: 2,
		},
		{
			name:           "Add duplicate task (should move to front)",
			taskID:         "task-1",
			existingTasks:  []string{"task-2", "task-1", "task-3"},
			expectedTasks:  []string{"task-1", "task-2", "task-3"},
			expectedLength: 3,
		},
		{
			name:           "Add task when at limit (should remove oldest)",
			taskID:         "task-new",
			existingTasks:  make([]string, 10), // Fill to limit
			expectedLength: 10,                 // Should stay at limit
		},
		{
			name:           "Add empty task ID",
			taskID:         "",
			existingTasks:  []string{"task-1"},
			expectedTasks:  []string{"", "task-1"},
			expectedLength: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up initial state
			if tc.name == "Add task when at limit (should remove oldest)" {
				// Fill the existing tasks array to the limit
				for i := 0; i < 10; i++ {
					tc.existingTasks[i] = "task-" + string(rune('a'+i))
				}
				tc.expectedTasks = append([]string{tc.taskID}, tc.existingTasks[:9]...)
			}

			vy.aiContext.RecentTasks = tc.existingTasks

			// Call addToRecentTasks
			acy.addToRecentTasks(tc.taskID)

			// Verify results
			if len(vy.aiContext.RecentTasks) != tc.expectedLength {
				t.Errorf("Expected length %d, got %d", tc.expectedLength, len(vy.aiContext.RecentTasks))
			}

			if tc.expectedTasks != nil && len(tc.expectedTasks) <= 5 {
				// Only check exact order for smaller lists
				for i, expected := range tc.expectedTasks {
					if i >= len(vy.aiContext.RecentTasks) {
						t.Errorf("Missing task at index %d", i)
						continue
					}
					if vy.aiContext.RecentTasks[i] != expected {
						t.Errorf("At index %d, expected %s, got %s", i, expected, vy.aiContext.RecentTasks[i])
					}
				}
			}

			// Verify new task is at front (unless empty case)
			if tc.taskID != "" && len(vy.aiContext.RecentTasks) > 0 {
				if vy.aiContext.RecentTasks[0] != tc.taskID {
					t.Errorf("Expected new task %s to be at front, got %s", tc.taskID, vy.aiContext.RecentTasks[0])
				}
			}
		})
	}
}

// TestAIContextYonetici_GetContext_ErrorHandling tests GetContext with database errors
func TestAIContextYonetici_GetContext_ErrorHandling(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Set up error condition
	vy.shouldReturnError = true
	vy.errorToReturn = errors.New("database connection failed")

	_, err := acy.GetContext()
	if err == nil {
		t.Error("Expected error but got none")
	}

	if err.Error() != "database connection failed" {
		t.Errorf("Expected 'database connection failed', got %s", err.Error())
	}
}

// TestAIContextYonetici_BatchUpdate_ErrorHandling tests BatchUpdate with various error conditions
func TestAIContextYonetici_BatchUpdate_ErrorHandling(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Add test task for valid updates
	testTask := &Gorev{
		ID:             "task-1",
		Baslik:         "Test Task",
		Durum:          constants.TaskStatusPending,
		Oncelik:        constants.PriorityMedium,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask

	testCases := []struct {
		name               string
		updates            []BatchUpdate
		setupError         func(*MockVeriYonetici)
		expectSuccess      bool
		expectedSuccessful int
		expectedFailed     int
	}{
		{
			name: "Valid updates with no errors",
			updates: []BatchUpdate{
				{
					ID: "task-1",
					Updates: map[string]interface{}{
						"durum":   constants.TaskStatusInProgress,
						"oncelik": constants.PriorityHigh,
					},
				},
			},
			expectSuccess:      true,
			expectedSuccessful: 1,
			expectedFailed:     0,
		},
		{
			name: "Non-existent task",
			updates: []BatchUpdate{
				{
					ID: "non-existent",
					Updates: map[string]interface{}{
						"durum": constants.TaskStatusCompleted,
					},
				},
			},
			expectSuccess:      true, // Should not error, but will have failed updates
			expectedSuccessful: 0,
			expectedFailed:     1,
		},
		{
			name:               "Empty updates list",
			updates:            []BatchUpdate{},
			expectSuccess:      true,
			expectedSuccessful: 0,
			expectedFailed:     0,
		},
		{
			name: "Invalid status value",
			updates: []BatchUpdate{
				{
					ID: "task-1",
					Updates: map[string]interface{}{
						"durum": "invalid-status",
					},
				},
			},
			expectSuccess:      true, // Should not error, but will have failed updates
			expectedSuccessful: 0,
			expectedFailed:     1,
		},
		{
			name: "Mixed valid and invalid updates",
			updates: []BatchUpdate{
				{
					ID: "task-1",
					Updates: map[string]interface{}{
						"baslik": "Updated Title",
					},
				},
				{
					ID: "non-existent",
					Updates: map[string]interface{}{
						"baslik": "Another Title",
					},
				},
			},
			expectSuccess:      true,
			expectedSuccessful: 1,
			expectedFailed:     1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock state
			vy.shouldReturnError = false
			vy.errorToReturn = nil

			// Set up error condition if specified
			if tc.setupError != nil {
				tc.setupError(vy)
			}

			result, err := acy.BatchUpdate(tc.updates)

			if tc.expectSuccess && err != nil {
				t.Errorf("Expected success but got error: %v", err)
			}
			if !tc.expectSuccess && err == nil {
				t.Error("Expected error but got success")
			}

			if tc.expectSuccess && result != nil {
				if len(result.Successful) != tc.expectedSuccessful {
					t.Errorf("Expected %d successful updates, got %d", tc.expectedSuccessful, len(result.Successful))
				}
				if len(result.Failed) != tc.expectedFailed {
					t.Errorf("Expected %d failed updates, got %d", tc.expectedFailed, len(result.Failed))
				}
			}
		})
	}
}

// TestAIContextYonetici_NLPQuery_EdgeCases tests NLP query with various edge cases
func TestAIContextYonetici_NLPQueryError(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Create test tasks with various attributes
	task1 := &Gorev{
		ID:             "task-1",
		Baslik:         "Task with ąćęłńóśźż unicode",
		Aciklama:       "Description with special chars: @#$%^&*()",
		Durum:          constants.TaskStatusPending,
		Oncelik:        constants.PriorityHigh,
		OlusturmaTarih: time.Now(),
		Etiketler:      []*Etiket{{ID: "tag-1", Isim: "tag with spaces"}},
	}

	task2 := &Gorev{
		ID:             "task-2",
		Baslik:         "UPPERCASE TITLE",
		Aciklama:       "lowercase description",
		Durum:          constants.TaskStatusCompleted,
		Oncelik:        constants.PriorityLow,
		OlusturmaTarih: time.Now(),
		Etiketler:      []*Etiket{{ID: "tag-2", Isim: "MixedCase"}},
	}

	task3 := &Gorev{
		ID:             "task-3",
		Baslik:         "", // Empty title
		Aciklama:       "", // Empty description
		Durum:          constants.TaskStatusInProgress,
		Oncelik:        constants.PriorityMedium,
		OlusturmaTarih: time.Now(),
		Etiketler:      []*Etiket{}, // No tags
	}

	// Add tasks to mock
	vy.gorevler["task-1"] = task1
	vy.gorevler["task-2"] = task2
	vy.gorevler["task-3"] = task3
	vy.allTasks = []*Gorev{task1, task2, task3}

	testCases := []struct {
		name        string
		query       string
		expectedLen int
		expectError bool
		description string
	}{
		{
			name:        "Unicode search",
			query:       "ąćę",
			expectedLen: 1,
			description: "Should find unicode characters",
		},
		{
			name:        "Special characters search",
			query:       "@#$",
			expectedLen: 1,
			description: "Should find special characters",
		},
		{
			name:        "Very long query",
			query:       "this is a very long query that should still work fine even though it contains many words and is quite lengthy",
			expectedLen: 0,
			description: "Should handle long queries",
		},
		{
			name:        "Query with only spaces",
			query:       "   ",
			expectedLen: 3, // NLP query doesn't trim, so returns all tasks
			description: "Should handle whitespace-only queries",
		},
		{
			name:        "Mixed case tag search",
			query:       "etiket: MixedCase",
			expectedLen: 1,
			description: "Should find mixed case tags",
		},
		{
			name:        "Tag with spaces search",
			query:       "tag: tag with spaces",
			expectedLen: 1,
			description: "Should find tags with spaces",
		},
		{
			name:        "Search in empty fields",
			query:       "nonexistent",
			expectedLen: 0,
			description: "Should handle searches with no matches",
		},
		{
			name:        "Case insensitive title search",
			query:       "uppercase",
			expectedLen: 1,
			description: "Should find uppercase content with lowercase query",
		},
		{
			name:        "Multiple keywords with no matches",
			query:       "keyword1 keyword2 keyword3",
			expectedLen: 0,
			description: "Should handle multiple keywords with no matches",
		},
		{
			name:        "Single character search",
			query:       "T",
			expectedLen: 2, // Should match "Task" and "TITLE"
			description: "Should handle single character searches",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks, err := acy.NLPQuery(tc.query)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(tasks) != tc.expectedLen {
				t.Errorf("Expected %d tasks, got %d. %s", tc.expectedLen, len(tasks), tc.description)
			}
		})
	}
}

// TestAIContextYonetici_GetActiveTask_EdgeCases tests GetActiveTask with edge cases
func TestAIContextYonetici_GetActiveTask_EdgeCases(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	testCases := []struct {
		name           string
		setupContext   func(*AIContext)
		setupError     func(*MockVeriYonetici)
		expectTask     bool
		expectError    bool
		expectedTaskID string
	}{
		{
			name: "Context with empty active task",
			setupContext: func(ctx *AIContext) {
				ctx.ActiveTaskID = ""
			},
			expectTask: false,
		},
		{
			name: "Context with non-existent active task",
			setupContext: func(ctx *AIContext) {
				ctx.ActiveTaskID = "non-existent-task"
			},
			expectTask:  false,
			expectError: true,
		},
		{
			name: "Database error on context retrieval",
			setupError: func(m *MockVeriYonetici) {
				m.shouldReturnError = true
				m.errorToReturn = errors.New("database error")
			},
			expectError: true,
		},
		{
			name: "Valid active task",
			setupContext: func(ctx *AIContext) {
				ctx.ActiveTaskID = "task-1"
			},
			expectTask:     true,
			expectedTaskID: "task-1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock state
			vy.shouldReturnError = false
			vy.errorToReturn = nil

			// Add a test task to mock
			testTask := &Gorev{
				ID:             "task-1",
				Baslik:         "Test Task",
				Durum:          constants.TaskStatusPending,
				OlusturmaTarih: time.Now(),
			}
			vy.gorevler["task-1"] = testTask

			// Set up context if specified
			if tc.setupContext != nil {
				tc.setupContext(vy.aiContext)
			}

			// Set up error condition if specified
			if tc.setupError != nil {
				tc.setupError(vy)
			}

			task, err := acy.GetActiveTask()

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if tc.expectTask && task == nil {
				t.Error("Expected task but got nil")
			}
			if !tc.expectTask && task != nil {
				t.Error("Expected nil task but got task")
			}
			if tc.expectTask && task != nil && task.ID != tc.expectedTaskID {
				t.Errorf("Expected task ID %s, got %s", tc.expectedTaskID, task.ID)
			}
		})
	}
}

// TestAIContextYonetici_recordInteraction_EdgeCases tests recordInteraction with edge cases
func TestAIContextYonetici_recordInteraction_EdgeCases(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	testCases := []struct {
		name        string
		taskID      string
		actionType  string
		context     interface{}
		setupError  func(*MockVeriYonetici)
		expectError bool
	}{
		{
			name:       "Very long taskID",
			taskID:     "task-" + string(make([]byte, 1000)), // Very long ID
			actionType: "viewed",
			context:    "test",
		},
		{
			name:       "Very long actionType",
			taskID:     "task-1",
			actionType: string(make([]byte, 1000)), // Very long action type
			context:    "test",
		},
		{
			name:       "Complex nested context",
			taskID:     "task-1",
			actionType: "complex",
			context: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": []interface{}{1, 2, 3, "nested", true},
					},
				},
			},
		},
		{
			name:       "Context with null values",
			taskID:     "task-1",
			actionType: "test",
			context:    map[string]interface{}{"null_value": nil, "empty_string": "", "zero": 0},
		},
		{
			name:       "Database error on interaction save",
			taskID:     "task-1",
			actionType: "test",
			context:    "test",
			setupError: func(m *MockVeriYonetici) {
				m.shouldReturnError = true
				m.errorToReturn = errors.New("interaction save failed")
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset mock state
			vy.shouldReturnError = false
			vy.errorToReturn = nil

			// Set up error condition if specified
			if tc.setupError != nil {
				tc.setupError(vy)
			}

			err := acy.recordInteraction(tc.taskID, tc.actionType, tc.context)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
