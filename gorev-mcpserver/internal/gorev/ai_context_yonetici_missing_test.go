package gorev

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
)

// TestAIContextYonetici_SetAutoStateManager tests the SetAutoStateManager function
func TestAIContextYonetici_SetAutoStateManager(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Create a mock AutoStateManager
	asm := &AutoStateManager{}

	// Test setting auto state manager
	acy.SetAutoStateManager(asm)

	// Verify it was set
	if acy.autoStateManager != asm {
		t.Error("AutoStateManager was not set correctly")
	}
}

// TestAIContextYonetici_RecordInteraction tests the public RecordInteraction function
func TestAIContextYonetici_RecordInteraction(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	testCases := []struct {
		name        string
		taskID      string
		actionType  string
		context     interface{}
		shouldError bool
	}{
		{
			name:        "Valid interaction with string context",
			taskID:      "task-123",
			actionType:  "viewed",
			context:     "user clicked task",
			shouldError: false,
		},
		{
			name:        "Valid interaction with map context",
			taskID:      "task-456",
			actionType:  "updated",
			context:     map[string]interface{}{"field": "status", "value": "completed"},
			shouldError: false,
		},
		{
			name:        "Valid interaction with nil context",
			taskID:      "task-789",
			actionType:  "created",
			context:     nil,
			shouldError: false,
		},
		{
			name:        "Empty task ID",
			taskID:      "",
			actionType:  "viewed",
			context:     nil,
			shouldError: false, // We don't validate task ID in recordInteraction
		},
		{
			name:        "Empty action type",
			taskID:      "task-123",
			actionType:  "",
			context:     nil,
			shouldError: false, // We don't validate action type in recordInteraction
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := acy.RecordInteraction(context.Background(), tc.taskID, tc.actionType, tc.context)

			if tc.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// TestAIContextYonetici_getTodayInteractions tests the getTodayInteractions function
func TestAIContextYonetici_getTodayInteractions(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Create some test interactions for today
	now := time.Now()
	todayInteraction := &AIInteraction{
		ID:         "int-1",
		GorevID:    "task-1",
		ActionType: "viewed",
		Context:    "",
		Timestamp:  now,
	}

	// Mock the return value
	vy.todayInteractions = []*AIInteraction{todayInteraction}

	// Test getting today's interactions
	interactions, err := acy.getTodayInteractions()
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if len(interactions) != 1 {
		t.Errorf("Expected 1 interaction, got %d", len(interactions))
	}

	if interactions[0].GorevID != "task-1" {
		t.Errorf("Expected task-1, got %s", interactions[0].GorevID)
	}
}

// TestAIContextYonetici_getTasksFromInteractions tests the getTasksFromInteractions function
func TestAIContextYonetici_getTasksFromInteractions(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Create test tasks
	task1 := &Gorev{
		ID:        "task-1",
		Title:     "Test Task 1",
		Status:    constants.TaskStatusPending,
		Priority:  constants.PriorityMedium,
		CreatedAt: time.Now(),
	}

	task2 := &Gorev{
		ID:        "task-2",
		Title:     "Test Task 2",
		Status:    constants.TaskStatusInProgress,
		Priority:  constants.PriorityHigh,
		CreatedAt: time.Now(),
	}

	// Add tasks to mock
	vy.gorevler["task-1"] = task1
	vy.gorevler["task-2"] = task2

	// Create test interactions
	interactions := []*AIInteraction{
		{
			ID:         "int-1",
			GorevID:    "task-1",
			ActionType: "viewed",
			Timestamp:  time.Now(),
		},
		{
			ID:         "int-2",
			GorevID:    "task-2",
			ActionType: "updated",
			Timestamp:  time.Now(),
		},
		{
			ID:         "int-3",
			GorevID:    "task-1", // Duplicate - should only appear once
			ActionType: "viewed",
			Timestamp:  time.Now(),
		},
		{
			ID:         "int-4",
			GorevID:    "task-nonexistent", // Non-existent task - should be skipped
			ActionType: "viewed",
			Timestamp:  time.Now(),
		},
	}

	// Test getting tasks from interactions
	tasks, err := acy.getTasksFromInteractions(context.Background(), interactions)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 unique tasks, got %d", len(tasks))
	}

	// Verify tasks are correct (order may vary due to map iteration)
	taskIDs := make(map[string]bool)
	for _, task := range tasks {
		taskIDs[task.ID] = true
	}

	if !taskIDs["task-1"] || !taskIDs["task-2"] {
		t.Error("Expected task-1 and task-2 to be returned")
	}
}

// TestAIContextYonetici_getLastCreatedTasks tests the getLastCreatedTasks function
func TestAIContextYonetici_getLastCreatedTasks(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Create test tasks with different creation times
	baseTime := time.Now()

	task1 := &Gorev{
		ID:        "task-1",
		Title:     "Oldest Task",
		Status:    constants.TaskStatusPending,
		Priority:  constants.PriorityMedium,
		CreatedAt: baseTime.Add(-2 * time.Hour),
	}

	task2 := &Gorev{
		ID:        "task-2",
		Title:     "Middle Task",
		Status:    constants.TaskStatusPending,
		Priority:  constants.PriorityMedium,
		CreatedAt: baseTime.Add(-1 * time.Hour),
	}

	task3 := &Gorev{
		ID:        "task-3",
		Title:     "Newest Task",
		Status:    constants.TaskStatusPending,
		Priority:  constants.PriorityMedium,
		CreatedAt: baseTime,
	}

	// Add tasks to mock
	vy.gorevler["task-1"] = task1
	vy.gorevler["task-2"] = task2
	vy.gorevler["task-3"] = task3

	// Set up mock to return all tasks
	vy.allTasks = []*Gorev{task1, task2, task3}

	testCases := []struct {
		name          string
		limit         int
		expectedLen   int
		expectedFirst string // ID of first task (newest)
	}{
		{
			name:          "Get last 2 tasks",
			limit:         2,
			expectedLen:   2,
			expectedFirst: "task-3", // Newest task should be first
		},
		{
			name:          "Get last 5 tasks (more than available)",
			limit:         5,
			expectedLen:   3,
			expectedFirst: "task-3",
		},
		{
			name:          "Get last 1 task",
			limit:         1,
			expectedLen:   1,
			expectedFirst: "task-3",
		},
		{
			name:          "Get last 0 tasks",
			limit:         0,
			expectedLen:   0,
			expectedFirst: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks, err := acy.getLastCreatedTasks(context.Background(), tc.limit)
			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(tasks) != tc.expectedLen {
				t.Errorf("Expected %d tasks, got %d", tc.expectedLen, len(tasks))
			}

			if tc.expectedLen > 0 && tasks[0].ID != tc.expectedFirst {
				t.Errorf("Expected first task to be %s, got %s", tc.expectedFirst, tasks[0].ID)
			}

			// Verify sorting (newest first)
			if len(tasks) > 1 {
				for i := 1; i < len(tasks); i++ {
					if tasks[i-1].CreatedAt.Before(tasks[i].CreatedAt) {
						t.Error("Tasks are not sorted by creation date (newest first)")
					}
				}
			}
		})
	}
}

// TestAIContextYonetici_NLPQuery tests the NLP query functionality with missing branches
func TestAIContextYonetici_NLPQueryBasic(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Create test tasks
	task1 := &Gorev{
		ID:          "task-1",
		Title:       "Test Task 1",
		Description: "This is a test task",
		Status:      constants.TaskStatusPending,
		Priority:    constants.PriorityHigh,
		CreatedAt:   time.Now(),
		Tags:        []*Etiket{{ID: "tag-1", Name: "urgent"}},
	}

	task2 := &Gorev{
		ID:          "task-2",
		Title:       "Another Task",
		Description: "Different task description",
		Status:      constants.TaskStatusCompleted,
		Priority:    constants.PriorityLow,
		CreatedAt:   time.Now(),
		Tags:        []*Etiket{{ID: "tag-2", Name: "normal"}},
	}

	// Add tasks to mock
	vy.gorevler["task-1"] = task1
	vy.gorevler["task-2"] = task2
	vy.allTasks = []*Gorev{task1, task2}

	testCases := []struct {
		name        string
		query       string
		expectedLen int
		expectError bool
	}{
		{
			name:        "Search by tag",
			query:       "etiket: urgent",
			expectedLen: 1,
			expectError: false,
		},
		{
			name:        "Search by tag with 'tag:' syntax",
			query:       "tag: normal",
			expectedLen: 1,
			expectError: false,
		},
		{
			name:        "Search by project (empty result)",
			query:       "proje: test-project",
			expectedLen: 0,
			expectError: false,
		},
		{
			name:        "Search by project with 'project:' syntax",
			query:       "project: some-project",
			expectedLen: 0,
			expectError: false,
		},
		{
			name:        "Text search - single term",
			query:       "test",
			expectedLen: 1, // Only task1 has "test" in title/description
			expectError: false,
		},
		{
			name:        "Text search - multiple terms",
			query:       "test task",
			expectedLen: 1, // Only task1 has both "test" and "task"
			expectError: false,
		},
		{
			name:        "Text search - no matches",
			query:       "nonexistent keyword",
			expectedLen: 0,
			expectError: false,
		},
		{
			name:        "Empty tag query",
			query:       "etiket:",
			expectedLen: 0,
			expectError: false,
		},
		{
			name:        "Case insensitive search",
			query:       "TEST TASK",
			expectedLen: 1,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tasks, err := acy.NLPQuery(context.Background(), tc.query)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if len(tasks) != tc.expectedLen {
				t.Errorf("Expected %d tasks, got %d", tc.expectedLen, len(tasks))
			}
		})
	}
}

// TestAIContextYonetici_NLPQuery_DatabaseError tests NLP query error handling
func TestAIContextYonetici_NLPQuery_DatabaseError(t *testing.T) {
	vy := NewMockVeriYonetici()
	acy := YeniAIContextYonetici(vy)

	// Set up mock to return error
	vy.shouldReturnError = true
	vy.errorToReturn = errors.New("database connection failed")

	_, err := acy.NLPQuery(context.Background(), "some query")
	if err == nil {
		t.Error("Expected error but got none")
	}

	if err.Error() != "database connection failed" {
		t.Errorf("Expected 'database connection failed', got %s", err.Error())
	}
}
