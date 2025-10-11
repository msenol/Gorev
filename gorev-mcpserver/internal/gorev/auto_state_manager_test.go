package gorev

import (
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
)

// TestYeniAutoStateManager tests the AutoStateManager constructor
func TestYeniAutoStateManager(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)

	if asm == nil {
		t.Fatal("YeniAutoStateManager returned nil")
	}

	if asm.veriYonetici == nil {
		t.Error("veriYonetici not properly set")
	}

	if asm.inactivityTimer != 30*time.Minute {
		t.Errorf("Expected default inactivity timer to be 30 minutes, got %v", asm.inactivityTimer)
	}

	if asm.activeTimers == nil {
		t.Error("activeTimers map not initialized")
	}

	if len(asm.activeTimers) != 0 {
		t.Error("activeTimers should be empty initially")
	}

	if asm.aiContextManager != nil {
		t.Error("aiContextManager should be nil initially")
	}

	if asm.nlpProcessor == nil {
		t.Error("nlpProcessor should be initialized")
	}
}

// TestSetAIContextManager tests the SetAIContextManager function
func TestAutoStateManager_SetAIContextManager(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)
	acy := YeniAIContextYonetici(vy)

	// Initially should be nil
	if asm.aiContextManager != nil {
		t.Error("aiContextManager should be nil initially")
	}

	// Set the AI context manager
	asm.SetAIContextManager(acy)

	// Verify it was set
	if asm.aiContextManager == nil {
		t.Error("aiContextManager not set properly")
	}
	if asm.aiContextManager != acy {
		t.Error("aiContextManager does not match the set value")
	}
}

// TestAutoTransitionToInProgress tests automatic transition to in-progress
func TestAutoStateManager_AutoTransitionToInProgress(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)
	acy := YeniAIContextYonetici(vy)
	asm.SetAIContextManager(acy)

	// Create test tasks with different statuses
	pendingTask := &Gorev{
		ID:        "pending-task",
		Title:     "Pending Task",
		Status:    constants.TaskStatusPending,
		CreatedAt: time.Now(),
	}

	inProgressTask := &Gorev{
		ID:        "inprogress-task",
		Title:     "In Progress Task",
		Status:    constants.TaskStatusInProgress,
		CreatedAt: time.Now(),
	}

	completedTask := &Gorev{
		ID:        "completed-task",
		Title:     "Completed Task",
		Status:    constants.TaskStatusCompleted,
		CreatedAt: time.Now(),
	}

	vy.gorevler["pending-task"] = pendingTask
	vy.gorevler["inprogress-task"] = inProgressTask
	vy.gorevler["completed-task"] = completedTask

	testCases := []struct {
		name           string
		taskID         string
		expectedStatus string
		expectError    bool
		expectTimer    bool
		description    string
	}{
		{
			name:           "Transition pending to in-progress",
			taskID:         "pending-task",
			expectedStatus: constants.TaskStatusInProgress,
			expectError:    false,
			expectTimer:    true,
			description:    "Should transition pending task to in-progress and start timer",
		},
		{
			name:           "Skip transition - already in progress",
			taskID:         "inprogress-task",
			expectedStatus: constants.TaskStatusInProgress,
			expectError:    false,
			expectTimer:    false,
			description:    "Should skip transition for task already in progress",
		},
		{
			name:           "Skip transition - already completed",
			taskID:         "completed-task",
			expectedStatus: constants.TaskStatusCompleted,
			expectError:    false,
			expectTimer:    false,
			description:    "Should skip transition for completed task",
		},
		{
			name:        "Non-existent task",
			taskID:      "non-existent",
			expectError: true,
			description: "Should return error for non-existent task",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear any existing timers
			asm.Cleanup()

			// Reset task states for each test
			if tc.taskID == "pending-task" {
				pendingTask.Status = constants.TaskStatusPending
			}

			err := asm.AutoTransitionToInProgress(tc.taskID)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for %s", tc.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.description, err)
				return
			}

			// Check if task status was updated correctly
			if tc.expectedStatus != "" {
				task, exists := vy.gorevler[tc.taskID]
				if !exists {
					t.Errorf("Task not found: %s", tc.taskID)
					return
				}

				if task.Status != tc.expectedStatus {
					t.Errorf("Expected status %s, got %s for %s", tc.expectedStatus, task.Status, tc.description)
				}
			}

			// Check if timer was started
			if tc.expectTimer {
				if _, exists := asm.activeTimers[tc.taskID]; !exists {
					t.Errorf("Expected timer to be started for %s", tc.description)
				}
			} else {
				if _, exists := asm.activeTimers[tc.taskID]; exists {
					t.Errorf("Expected no timer to be started for %s", tc.description)
				}
			}
		})
	}
}

// TestAutoTransitionToPending tests automatic transition back to pending
func TestAutoStateManager_AutoTransitionToPending(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)
	acy := YeniAIContextYonetici(vy)
	asm.SetAIContextManager(acy)

	// Create test tasks
	inProgressTask := &Gorev{
		ID:        "inprogress-task",
		Title:     "In Progress Task",
		Status:    constants.TaskStatusInProgress,
		CreatedAt: time.Now(),
	}

	pendingTask := &Gorev{
		ID:        "pending-task",
		Title:     "Pending Task",
		Status:    constants.TaskStatusPending,
		CreatedAt: time.Now(),
	}

	completedTask := &Gorev{
		ID:        "completed-task",
		Title:     "Completed Task",
		Status:    constants.TaskStatusCompleted,
		CreatedAt: time.Now(),
	}

	vy.gorevler["inprogress-task"] = inProgressTask
	vy.gorevler["pending-task"] = pendingTask
	vy.gorevler["completed-task"] = completedTask

	// Add timers for some tasks
	asm.startInactivityTimer("inprogress-task")
	asm.startInactivityTimer("pending-task")

	testCases := []struct {
		name           string
		taskID         string
		expectedStatus string
		expectError    bool
		expectNoTimer  bool
		description    string
	}{
		{
			name:           "Transition in-progress to pending",
			taskID:         "inprogress-task",
			expectedStatus: constants.TaskStatusPending,
			expectError:    false,
			expectNoTimer:  true,
			description:    "Should transition in-progress task back to pending and clear timer",
		},
		{
			name:           "Skip transition - already pending",
			taskID:         "pending-task",
			expectedStatus: constants.TaskStatusPending,
			expectError:    false,
			expectNoTimer:  false, // Timer may or may not be cleared in this case
			description:    "Should skip transition for task already pending",
		},
		{
			name:           "Skip transition - completed task",
			taskID:         "completed-task",
			expectedStatus: constants.TaskStatusCompleted,
			expectError:    false,
			expectNoTimer:  false,
			description:    "Should skip transition for completed task",
		},
		{
			name:        "Non-existent task",
			taskID:      "non-existent",
			expectError: true,
			description: "Should return error for non-existent task",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset task states for each test
			if tc.taskID == "inprogress-task" {
				inProgressTask.Status = constants.TaskStatusInProgress
			}

			err := asm.AutoTransitionToPending(tc.taskID)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for %s", tc.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.description, err)
				return
			}

			// Check if task status was updated correctly
			if tc.expectedStatus != "" {
				task, exists := vy.gorevler[tc.taskID]
				if !exists {
					t.Errorf("Task not found: %s", tc.taskID)
					return
				}

				if task.Status != tc.expectedStatus {
					t.Errorf("Expected status %s, got %s for %s", tc.expectedStatus, task.Status, tc.description)
				}
			}

			// Check if timer was cleared
			if tc.expectNoTimer {
				if _, exists := asm.activeTimers[tc.taskID]; exists {
					t.Errorf("Expected timer to be cleared for %s", tc.description)
				}
			}
		})
	}
}

// TestOnTaskAccessed tests the OnTaskAccessed event handler
func TestAutoStateManager_OnTaskAccessed(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)

	// Create test task
	testTask := &Gorev{
		ID:        "test-task",
		Title:     "Test Task",
		Status:    constants.TaskStatusPending,
		CreatedAt: time.Now(),
	}
	vy.gorevler["test-task"] = testTask

	// Call OnTaskAccessed
	err := asm.OnTaskAccessed("test-task")
	if err != nil {
		t.Errorf("OnTaskAccessed returned error: %v", err)
	}

	// Verify task was transitioned to in-progress
	if testTask.Status != constants.TaskStatusInProgress {
		t.Errorf("Expected task to be in progress, got %s", testTask.Status)
	}

	// Verify timer was started
	if _, exists := asm.activeTimers["test-task"]; !exists {
		t.Error("Expected timer to be started")
	}
}

// TestOnTaskCompleted tests the OnTaskCompleted event handler
func TestAutoStateManager_OnTaskCompleted(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)

	// Create test task
	testTask := &Gorev{
		ID:        "test-task",
		Title:     "Test Task",
		Status:    constants.TaskStatusCompleted,
		CreatedAt: time.Now(),
	}
	vy.gorevler["test-task"] = testTask

	// Start a timer first
	asm.startInactivityTimer("test-task")

	// Verify timer exists
	if _, exists := asm.activeTimers["test-task"]; !exists {
		t.Error("Expected timer to exist before OnTaskCompleted")
	}

	// Call OnTaskCompleted
	err := asm.OnTaskCompleted("test-task")
	if err != nil {
		t.Errorf("OnTaskCompleted returned error: %v", err)
	}

	// Verify timer was cleared
	if _, exists := asm.activeTimers["test-task"]; exists {
		t.Error("Expected timer to be cleared after OnTaskCompleted")
	}
}

// TestInactivityConfiguration tests inactivity duration configuration
func TestAutoStateManager_InactivityConfiguration(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)

	// Test default duration
	defaultDuration := asm.GetInactivityDuration()
	if defaultDuration != 30*time.Minute {
		t.Errorf("Expected default duration to be 30 minutes, got %v", defaultDuration)
	}

	// Test setting custom duration
	customDuration := 15 * time.Minute
	asm.SetInactivityDuration(customDuration)

	newDuration := asm.GetInactivityDuration()
	if newDuration != customDuration {
		t.Errorf("Expected custom duration %v, got %v", customDuration, newDuration)
	}

	// Test with zero duration
	asm.SetInactivityDuration(0)
	zeroDuration := asm.GetInactivityDuration()
	if zeroDuration != 0 {
		t.Errorf("Expected zero duration, got %v", zeroDuration)
	}
}

// TestCleanup tests the cleanup functionality
func TestAutoStateManager_Cleanup(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)

	// Start multiple timers
	asm.startInactivityTimer("task-1")
	asm.startInactivityTimer("task-2")
	asm.startInactivityTimer("task-3")

	// Verify timers exist
	if len(asm.activeTimers) != 3 {
		t.Errorf("Expected 3 timers, got %d", len(asm.activeTimers))
	}

	// Call cleanup
	asm.Cleanup()

	// Verify all timers were cleared
	if len(asm.activeTimers) != 0 {
		t.Errorf("Expected 0 timers after cleanup, got %d", len(asm.activeTimers))
	}
}

// TestTimerManagement tests timer start/stop/reset functionality
func TestAutoStateManager_TimerManagement(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)

	taskID := "test-task"

	// Test starting timer
	asm.startInactivityTimer(taskID)
	if _, exists := asm.activeTimers[taskID]; !exists {
		t.Error("Expected timer to be started")
	}

	// Test clearing timer
	asm.clearInactivityTimer(taskID)
	if _, exists := asm.activeTimers[taskID]; exists {
		t.Error("Expected timer to be cleared")
	}

	// Test ScheduleInactivityCheck
	asm.ScheduleInactivityCheck(taskID)
	if _, exists := asm.activeTimers[taskID]; !exists {
		t.Error("Expected timer to be started by ScheduleInactivityCheck")
	}

	// Test ResetInactivityTimer
	asm.ResetInactivityTimer(taskID)
	if _, exists := asm.activeTimers[taskID]; !exists {
		t.Error("Expected timer to exist after reset")
	}

	// Clear for cleanup
	asm.clearInactivityTimer(taskID)
}

// TestCheckParentCompletion tests the parent completion checking functionality
func TestAutoStateManager_CheckParentCompletion(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)
	acy := YeniAIContextYonetici(vy)
	asm.SetAIContextManager(acy)

	// Create parent task
	parentTask := &Gorev{
		ID:        "parent-task",
		Title:     "Parent Task",
		Status:    constants.TaskStatusInProgress,
		CreatedAt: time.Now(),
	}

	// Create child task
	childTask := &Gorev{
		ID:        "child-task",
		Title:     "Child Task",
		Status:    constants.TaskStatusCompleted,
		ParentID:  "parent-task",
		CreatedAt: time.Now(),
	}

	// Create task without parent
	orphanTask := &Gorev{
		ID:        "orphan-task",
		Title:     "Orphan Task",
		Status:    constants.TaskStatusCompleted,
		CreatedAt: time.Now(),
	}

	vy.gorevler["parent-task"] = parentTask
	vy.gorevler["child-task"] = childTask
	vy.gorevler["orphan-task"] = orphanTask

	testCases := []struct {
		name                 string
		taskID               string
		expectedParentStatus string
		expectError          bool
		description          string
	}{
		{
			name:                 "Complete parent when child completed",
			taskID:               "child-task",
			expectedParentStatus: constants.TaskStatusCompleted,
			expectError:          false,
			description:          "Should complete parent task when all children are done",
		},
		{
			name:        "Skip for task without parent",
			taskID:      "orphan-task",
			expectError: false,
			description: "Should skip completion check for task without parent",
		},
		{
			name:        "Error for non-existent task",
			taskID:      "non-existent",
			expectError: true,
			description: "Should return error for non-existent task",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset parent status for each test
			if tc.taskID == "child-task" {
				parentTask.Status = constants.TaskStatusInProgress
			}

			err := asm.CheckParentCompletion(tc.taskID)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for %s", tc.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.description, err)
				return
			}

			// Check parent status if expected
			if tc.expectedParentStatus != "" {
				parent := vy.gorevler["parent-task"]
				if parent.Status != tc.expectedParentStatus {
					t.Errorf("Expected parent status %s, got %s for %s", tc.expectedParentStatus, parent.Status, tc.description)
				}
			}
		})
	}
}

// TestCheckDependenciesCompleted tests the dependency checking functionality
func TestAutoStateManager_CheckDependenciesCompleted(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)

	// Create dependency tasks
	dep1 := &Gorev{
		ID:        "dep-1",
		Title:     "Dependency 1",
		Status:    constants.TaskStatusCompleted,
		CreatedAt: time.Now(),
	}

	dep2 := &Gorev{
		ID:        "dep-2",
		Title:     "Dependency 2",
		Status:    constants.TaskStatusInProgress,
		CreatedAt: time.Now(),
	}

	mainTask := &Gorev{
		ID:        "main-task",
		Title:     "Main Task",
		Status:    constants.TaskStatusPending,
		CreatedAt: time.Now(),
	}

	vy.gorevler["dep-1"] = dep1
	vy.gorevler["dep-2"] = dep2
	vy.gorevler["main-task"] = mainTask

	// Test with no dependencies (empty slice returned by mock)
	canStart, err := asm.checkDependenciesCompleted("main-task")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !canStart {
		t.Error("Expected task to be able to start with no dependencies")
	}

	// Note: The mock implementation always returns empty dependencies,
	// so we can't easily test the case with actual dependencies without
	// a more sophisticated mock or modifying the mock to return test data
}

// TestNaturalLanguageProcessing tests the NLP query processing
func TestAutoStateManager_ProcessNaturalLanguageQuery(t *testing.T) {
	vy := NewMockVeriYonetici()
	asm := YeniAutoStateManager(vy)
	acy := YeniAIContextYonetici(vy)
	asm.SetAIContextManager(acy)

	// Add some test tasks for queries to work with
	testTask1 := &Gorev{
		ID:        "task-1",
		Title:     "Complete project setup",
		Status:    constants.TaskStatusPending,
		Priority:  constants.PriorityHigh,
		CreatedAt: time.Now(),
	}

	testTask2 := &Gorev{
		ID:        "task-2",
		Title:     "Write documentation",
		Status:    constants.TaskStatusInProgress,
		Priority:  constants.PriorityMedium,
		CreatedAt: time.Now(),
	}

	vy.gorevler["task-1"] = testTask1
	vy.gorevler["task-2"] = testTask2

	testCases := []struct {
		name        string
		query       string
		lang        string
		expectError bool
		description string
	}{
		{
			name:        "Simple list query",
			query:       "show me all tasks",
			lang:        "en",
			expectError: false,
			description: "Should process basic list query",
		},
		{
			name:        "Priority query",
			query:       "show high priority tasks",
			lang:        "en",
			expectError: false,
			description: "Should process priority-based queries",
		},
		{
			name:        "Status query",
			query:       "what tasks are in progress",
			lang:        "en",
			expectError: false,
			description: "Should process status-based queries",
		},
		{
			name:        "Turkish query",
			query:       "tüm görevleri göster",
			lang:        "tr",
			expectError: false,
			description: "Should process Turkish queries",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := asm.ProcessNaturalLanguageQuery(tc.query, tc.lang)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for %s", tc.description)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.description, err)
				return
			}

			// Verify response structure
			if result == nil {
				t.Errorf("Expected non-nil result for %s", tc.description)
				return
			}

			// Check if result has expected structure
			responseMap, ok := result.(map[string]interface{})
			if !ok {
				t.Errorf("Expected result to be map[string]interface{} for %s", tc.description)
				return
			}

			if _, exists := responseMap["response"]; !exists {
				t.Errorf("Expected 'response' field in result for %s", tc.description)
			}

			if _, exists := responseMap["intent"]; !exists {
				t.Errorf("Expected 'intent' field in result for %s", tc.description)
			}
		})
	}
}
