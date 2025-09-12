package gorev

import (
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
)

// TestNewBatchProcessor tests the NewBatchProcessor constructor
func TestNewBatchProcessor(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	
	if bp == nil {
		t.Fatal("NewBatchProcessor returned nil")
	}
	
	if bp.veriYonetici == nil {
		t.Error("veriYonetici not properly set")
	}
	
	if bp.aiContextManager != nil {
		t.Error("aiContextManager should be nil initially")
	}
}

// TestSetAIContextManager tests the SetAIContextManager function
func TestSetAIContextManager(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	acy := YeniAIContextYonetici(vy)
	
	// Initially should be nil
	if bp.aiContextManager != nil {
		t.Error("aiContextManager should be nil initially")
	}
	
	// Set the AI context manager
	bp.SetAIContextManager(acy)
	
	// Verify it was set
	if bp.aiContextManager == nil {
		t.Error("aiContextManager not set properly")
	}
	if bp.aiContextManager != acy {
		t.Error("aiContextManager does not match the set value")
	}
}

// TestProcessBatchUpdate_Constructor tests ProcessBatchUpdate with basic scenarios
func TestProcessBatchUpdate_Basic(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	
	// Create test task
	testTask := &Gorev{
		ID:              "task-1",
		Baslik:         "Test Task",
		Aciklama:       "Test Description",
		Durum:          constants.TaskStatusPending,
		Oncelik:        constants.PriorityMedium,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask
	
	testCases := []struct {
		name                string
		requests            []BatchUpdateRequest
		expectedSuccessful  int
		expectedFailed      int
		expectedWarnings    int
		setupError          func(*MockVeriYonetici)
	}{
		{
			name:               "Empty batch",
			requests:           []BatchUpdateRequest{},
			expectedSuccessful: 0,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Single valid update - status",
			requests: []BatchUpdateRequest{
				{
					TaskID: "task-1",
					Updates: map[string]interface{}{
						"durum": constants.TaskStatusInProgress,
					},
				},
			},
			expectedSuccessful: 1,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Single valid update - priority",
			requests: []BatchUpdateRequest{
				{
					TaskID: "task-1",
					Updates: map[string]interface{}{
						"oncelik": constants.PriorityHigh,
					},
				},
			},
			expectedSuccessful: 1,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Single valid update - title",
			requests: []BatchUpdateRequest{
				{
					TaskID: "task-1",
					Updates: map[string]interface{}{
						"baslik": "Updated Title",
					},
				},
			},
			expectedSuccessful: 1,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Single valid update - description",
			requests: []BatchUpdateRequest{
				{
					TaskID: "task-1",
					Updates: map[string]interface{}{
						"aciklama": "Updated Description",
					},
				},
			},
			expectedSuccessful: 1,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Single valid update - due date",
			requests: []BatchUpdateRequest{
				{
					TaskID: "task-1",
					Updates: map[string]interface{}{
						"son_tarih": "2024-12-31",
					},
				},
			},
			expectedSuccessful: 1,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Non-existent task",
			requests: []BatchUpdateRequest{
				{
					TaskID: "non-existent",
					Updates: map[string]interface{}{
						"durum": constants.TaskStatusCompleted,
					},
				},
			},
			expectedSuccessful: 0,
			expectedFailed:     1,
			expectedWarnings:   0,
		},
		{
			name: "Mixed valid and invalid updates",
			requests: []BatchUpdateRequest{
				{
					TaskID: "task-1",
					Updates: map[string]interface{}{
						"baslik": "Valid Update",
					},
				},
				{
					TaskID: "non-existent",
					Updates: map[string]interface{}{
						"baslik": "Invalid Update",
					},
				},
			},
			expectedSuccessful: 1,
			expectedFailed:     1,
			expectedWarnings:   0,
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
			
			result, err := bp.ProcessBatchUpdate(tc.requests)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result == nil {
				t.Fatal("Result should not be nil")
			}
			
			if len(result.Successful) != tc.expectedSuccessful {
				t.Errorf("Expected %d successful updates, got %d", tc.expectedSuccessful, len(result.Successful))
			}
			
			if len(result.Failed) != tc.expectedFailed {
				t.Errorf("Expected %d failed updates, got %d", tc.expectedFailed, len(result.Failed))
			}
			
			if len(result.Warnings) != tc.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d", tc.expectedWarnings, len(result.Warnings))
			}
			
			if result.TotalProcessed != len(tc.requests) {
				t.Errorf("Expected total processed %d, got %d", len(tc.requests), result.TotalProcessed)
			}
			
			if result.ExecutionTime == 0 {
				t.Error("ExecutionTime should be greater than 0")
			}
			
			if result.Summary == "" {
				t.Error("Summary should not be empty")
			}
		})
	}
}

// TestProcessBatchUpdate_DryRun tests dry run functionality
func TestProcessBatchUpdate_DryRun(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	
	// Create test task
	testTask := &Gorev{
		ID:              "task-1",
		Baslik:         "Test Task",
		Durum:          constants.TaskStatusPending,
		Oncelik:        constants.PriorityMedium,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask
	
	testCases := []struct {
		name               string
		request            BatchUpdateRequest
		expectedSuccessful int
		expectedFailed     int
		expectedWarnings   int
	}{
		{
			name: "Valid dry run",
			request: BatchUpdateRequest{
				TaskID: "task-1",
				Updates: map[string]interface{}{
					"durum": constants.TaskStatusInProgress,
				},
				DryRun: true,
			},
			expectedSuccessful: 0, // Dry run doesn't execute
			expectedFailed:     0,
			expectedWarnings:   1, // Should have dry run warning
		},
		{
			name: "Invalid dry run - non-existent task",
			request: BatchUpdateRequest{
				TaskID: "non-existent",
				Updates: map[string]interface{}{
					"durum": constants.TaskStatusInProgress,
				},
				DryRun: true,
			},
			expectedSuccessful: 0,
			expectedFailed:     1, // Should fail validation
			expectedWarnings:   0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := bp.ProcessBatchUpdate([]BatchUpdateRequest{tc.request})
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if len(result.Successful) != tc.expectedSuccessful {
				t.Errorf("Expected %d successful updates, got %d", tc.expectedSuccessful, len(result.Successful))
			}
			
			if len(result.Failed) != tc.expectedFailed {
				t.Errorf("Expected %d failed updates, got %d", tc.expectedFailed, len(result.Failed))
			}
			
			if len(result.Warnings) != tc.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d", tc.expectedWarnings, len(result.Warnings))
			}
		})
	}
}

// TestProcessBatchUpdate_WithWarnings tests scenarios that generate warnings
func TestProcessBatchUpdate_WithWarnings(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	
	// Create test task
	testTask := &Gorev{
		ID:              "task-1",
		Baslik:         "Test Task",
		Durum:          constants.TaskStatusPending,
		Oncelik:        constants.PriorityMedium,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask
	
	testCases := []struct {
		name               string
		request            BatchUpdateRequest
		expectedSuccessful int
		expectedFailed     int
		expectedWarnings   int
	}{
		{
			name: "Empty title warning",
			request: BatchUpdateRequest{
				TaskID: "task-1",
				Updates: map[string]interface{}{
					"baslik": "", // Empty title should generate warning
				},
			},
			expectedSuccessful: 0, // No update should be made
			expectedFailed:     0,
			expectedWarnings:   1,
		},
		{
			name: "Empty title with spaces warning",
			request: BatchUpdateRequest{
				TaskID: "task-1",
				Updates: map[string]interface{}{
					"baslik": "   ", // Whitespace-only title should generate warning
				},
			},
			expectedSuccessful: 0,
			expectedFailed:     0,
			expectedWarnings:   1,
		},
		{
			name: "Invalid date format warning",
			request: BatchUpdateRequest{
				TaskID: "task-1",
				Updates: map[string]interface{}{
					"son_tarih": "invalid-date", // Invalid date format should generate warning
				},
			},
			expectedSuccessful: 0,
			expectedFailed:     0,
			expectedWarnings:   1,
		},
		{
			name: "Valid date clearing",
			request: BatchUpdateRequest{
				TaskID: "task-1",
				Updates: map[string]interface{}{
					"son_tarih": "", // Empty date should clear due date
				},
			},
			expectedSuccessful: 1, // Should be successful
			expectedFailed:     0,
			expectedWarnings:   0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := bp.ProcessBatchUpdate([]BatchUpdateRequest{tc.request})
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if len(result.Successful) != tc.expectedSuccessful {
				t.Errorf("Expected %d successful updates, got %d", tc.expectedSuccessful, len(result.Successful))
			}
			
			if len(result.Failed) != tc.expectedFailed {
				t.Errorf("Expected %d failed updates, got %d", tc.expectedFailed, len(result.Failed))
			}
			
			if len(result.Warnings) != tc.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d", tc.expectedWarnings, len(result.Warnings))
			}
		})
	}
}

// TestProcessBatchUpdate_WithAIContext tests batch updates with AI context manager
func TestProcessBatchUpdate_WithAIContext(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	acy := YeniAIContextYonetici(vy)
	bp.SetAIContextManager(acy)
	
	// Create test task
	testTask := &Gorev{
		ID:              "task-1",
		Baslik:         "Test Task",
		Durum:          constants.TaskStatusPending,
		Oncelik:        constants.PriorityMedium,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask
	
	request := BatchUpdateRequest{
		TaskID: "task-1",
		Updates: map[string]interface{}{
			"baslik": "Updated with AI Context",
		},
	}
	
	result, err := bp.ProcessBatchUpdate([]BatchUpdateRequest{request})
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	
	if len(result.Successful) != 1 {
		t.Errorf("Expected 1 successful update, got %d", len(result.Successful))
	}
	
	// Verify AI interaction was recorded
	if len(vy.interactions) == 0 {
		t.Error("Expected AI interaction to be recorded")
	}
	
	if len(vy.interactions) > 0 {
		interaction := vy.interactions[0]
		if interaction.GorevID != "task-1" {
			t.Errorf("Expected interaction for task-1, got %s", interaction.GorevID)
		}
		if interaction.ActionType != "batch_update" {
			t.Errorf("Expected action type 'batch_update', got %s", interaction.ActionType)
		}
	}
}

// TestBatchProcessor_ValidationFunctions tests the helper validation functions
func TestBatchProcessor_ValidationFunctions(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	
	// Create test task for validateUpdateRequest
	testTask := &Gorev{
		ID:              "task-1",
		Baslik:         "Test Task",
		Durum:          constants.TaskStatusPending,
		Oncelik:        constants.PriorityMedium,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask
	
	// Test validateStatus
	t.Run("validateStatus", func(t *testing.T) {
		validStatusTests := []struct {
			status   string
			expected bool
		}{
			{constants.TaskStatusPending, true},
			{constants.TaskStatusInProgress, true},
			{constants.TaskStatusCompleted, true},
			{constants.TaskStatusCancelled, true},
			{"invalid-status", false},
			{"", false},
		}
		
		for _, test := range validStatusTests {
			result := bp.validateStatus(test.status)
			if result != test.expected {
				t.Errorf("validateStatus(%q) = %v, expected %v", test.status, result, test.expected)
			}
		}
	})
	
	// Test validatePriority
	t.Run("validatePriority", func(t *testing.T) {
		validPriorityTests := []struct {
			priority string
			expected bool
		}{
			{constants.PriorityLow, true},
			{constants.PriorityMedium, true},
			{constants.PriorityHigh, true},
			{"invalid-priority", false},
			{"", false},
		}
		
		for _, test := range validPriorityTests {
			result := bp.validatePriority(test.priority)
			if result != test.expected {
				t.Errorf("validatePriority(%q) = %v, expected %v", test.priority, result, test.expected)
			}
		}
	})
	
	// Test validateStatusTransition
	t.Run("validateStatusTransition", func(t *testing.T) {
		validTransitionTests := []struct {
			from     string
			to       string
			expected bool
			desc     string
		}{
			// Valid transitions
			{constants.TaskStatusPending, constants.TaskStatusInProgress, true, "pending -> in progress"},
			{constants.TaskStatusPending, constants.TaskStatusCancelled, true, "pending -> cancelled"},
			{constants.TaskStatusInProgress, constants.TaskStatusPending, true, "in progress -> pending"},
			{constants.TaskStatusInProgress, constants.TaskStatusCompleted, true, "in progress -> completed"},
			{constants.TaskStatusInProgress, constants.TaskStatusCancelled, true, "in progress -> cancelled"},
			{constants.TaskStatusCompleted, constants.TaskStatusInProgress, true, "completed -> in progress (reopen)"},
			{constants.TaskStatusCancelled, constants.TaskStatusPending, true, "cancelled -> pending (reactivate)"},
			
			// Invalid transitions
			{constants.TaskStatusPending, constants.TaskStatusCompleted, false, "pending -> completed (not allowed)"},
			{constants.TaskStatusCompleted, constants.TaskStatusPending, false, "completed -> pending (not allowed)"},
			{constants.TaskStatusCompleted, constants.TaskStatusCancelled, false, "completed -> cancelled (not allowed)"},
			{constants.TaskStatusCancelled, constants.TaskStatusCompleted, false, "cancelled -> completed (not allowed)"},
			{constants.TaskStatusCancelled, constants.TaskStatusInProgress, false, "cancelled -> in progress (not allowed)"},
			
			// Same status
			{constants.TaskStatusPending, constants.TaskStatusPending, false, "pending -> pending (same status)"},
			{constants.TaskStatusInProgress, constants.TaskStatusInProgress, false, "in progress -> in progress (same status)"},
			
			// Invalid statuses
			{"invalid-from", constants.TaskStatusInProgress, false, "invalid from status"},
			{constants.TaskStatusPending, "invalid-to", false, "invalid to status"},
			{"", constants.TaskStatusInProgress, false, "empty from status"},
			{constants.TaskStatusPending, "", false, "empty to status"},
		}
		
		for _, test := range validTransitionTests {
			result := bp.validateStatusTransition(test.from, test.to)
			if result != test.expected {
				t.Errorf("validateStatusTransition(%q, %q) = %v, expected %v (%s)", 
					test.from, test.to, result, test.expected, test.desc)
			}
		}
	})
	
	// Test validateUpdateRequest
	t.Run("validateUpdateRequest", func(t *testing.T) {
		validRequestTests := []struct {
			request     BatchUpdateRequest
			expectError bool
			desc        string
		}{
			// Valid requests
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"baslik": "Valid Title"},
				},
				false,
				"valid request with title update",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"durum": constants.TaskStatusInProgress},
				},
				false,
				"valid request with status update",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"oncelik": constants.PriorityHigh},
				},
				false,
				"valid request with priority update",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"aciklama": "New description"},
				},
				false,
				"valid request with description update",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{}, // Empty updates
				},
				false,
				"valid request with no updates",
			},
			
			// Invalid requests
			{
				BatchUpdateRequest{
					TaskID:  "non-existent",
					Updates: map[string]interface{}{"baslik": "Title"},
				},
				true,
				"invalid request - non-existent task",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"durum": "invalid-status"},
				},
				true,
				"invalid request - invalid status",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"oncelik": "invalid-priority"},
				},
				true,
				"invalid request - invalid priority",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"baslik": ""},
				},
				true,
				"invalid request - empty title",
			},
			{
				BatchUpdateRequest{
					TaskID:  "task-1",
					Updates: map[string]interface{}{"baslik": "   "},
				},
				true,
				"invalid request - whitespace-only title",
			},
		}
		
		for _, test := range validRequestTests {
			err := bp.validateUpdateRequest(test.request)
			if test.expectError && err == nil {
				t.Errorf("validateUpdateRequest: expected error for %s, but got none", test.desc)
			}
			if !test.expectError && err != nil {
				t.Errorf("validateUpdateRequest: expected no error for %s, but got: %v", test.desc, err)
			}
		}
	})
	
	// Test checkDependenciesCompleted (with mock dependencies)
	t.Run("checkDependenciesCompleted", func(t *testing.T) {
		// This would need more complex mocking for dependency management
		// For now, just test the basic call
		canStart, err := bp.checkDependenciesCompleted("task-1")
		if err != nil {
			t.Errorf("checkDependenciesCompleted returned error: %v", err)
		}
		// With our mock, no dependencies should exist, so should return true
		if !canStart {
			t.Error("checkDependenciesCompleted should return true for task with no dependencies")
		}
	})
}

// TestBulkStatusTransition_Basic tests basic bulk status transition functionality
func TestBulkStatusTransition_Basic(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	
	// Create test tasks with different statuses
	task1 := &Gorev{
		ID:              "task-1",
		Baslik:         "Task 1",
		Durum:          constants.TaskStatusPending,
		OlusturmaTarih: time.Now(),
	}
	task2 := &Gorev{
		ID:              "task-2", 
		Baslik:         "Task 2",
		Durum:          constants.TaskStatusPending,
		OlusturmaTarih: time.Now(),
	}
	task3 := &Gorev{
		ID:              "task-3",
		Baslik:         "Task 3", 
		Durum:          constants.TaskStatusCompleted, // Already completed
		OlusturmaTarih: time.Now(),
	}
	
	vy.gorevler["task-1"] = task1
	vy.gorevler["task-2"] = task2
	vy.gorevler["task-3"] = task3
	
	testCases := []struct {
		name               string
		request            BulkStatusTransitionRequest
		expectedSuccessful int
		expectedFailed     int
		expectedWarnings   int
		expectError        bool
	}{
		{
			name: "Valid status transition - pending to in progress",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"task-1", "task-2"},
				NewStatus: constants.TaskStatusInProgress,
			},
			expectedSuccessful: 2,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Invalid status transition - pending to completed", 
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"task-1"},
				NewStatus: constants.TaskStatusCompleted,
			},
			expectedSuccessful: 0, // This transition is not allowed by validateStatusTransition
			expectedFailed:     1, // Should fail validation
			expectedWarnings:   0,
		},
		{
			name: "Empty task list",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{},
				NewStatus: constants.TaskStatusInProgress,
			},
			expectedSuccessful: 0,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
		{
			name: "Non-existent task",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"non-existent"},
				NewStatus: constants.TaskStatusInProgress,
			},
			expectedSuccessful: 0,
			expectedFailed:     1,
			expectedWarnings:   0,
		},
		{
			name: "Already in target status (warning)",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"task-3"}, // Already completed
				NewStatus: constants.TaskStatusCompleted,
			},
			expectedSuccessful: 0,
			expectedFailed:     0,
			expectedWarnings:   1, // Should warn already in status
		},
		{
			name: "Mixed scenarios",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"task-1", "non-existent", "task-3"}, // valid, invalid, completed->in progress (valid)
				NewStatus: constants.TaskStatusInProgress,
			},
			expectedSuccessful: 2, // task-1 and task-3 (completed -> in progress is allowed)
			expectedFailed:     1, // non-existent
			expectedWarnings:   0, // task-3 transition is valid according to validateStatusTransition
		},
		{
			name: "Invalid status",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"task-1"},
				NewStatus: "invalid-status",
			},
			expectError: true, // Should return error for invalid status
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset task states for each test
			task1.Durum = constants.TaskStatusPending
			task2.Durum = constants.TaskStatusPending
			task3.Durum = constants.TaskStatusCompleted
			
			result, err := bp.BulkStatusTransition(tc.request)
			
			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result == nil {
				t.Fatal("Result should not be nil")
			}
			
			if len(result.Successful) != tc.expectedSuccessful {
				t.Errorf("Expected %d successful transitions, got %d", tc.expectedSuccessful, len(result.Successful))
			}
			
			if len(result.Failed) != tc.expectedFailed {
				t.Errorf("Expected %d failed transitions, got %d", tc.expectedFailed, len(result.Failed))
			}
			
			if len(result.Warnings) != tc.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d", tc.expectedWarnings, len(result.Warnings))
			}
			
			if result.TotalProcessed != len(tc.request.TaskIDs) {
				t.Errorf("Expected total processed %d, got %d", len(tc.request.TaskIDs), result.TotalProcessed)
			}
			
			if result.ExecutionTime == 0 {
				t.Error("ExecutionTime should be greater than 0")
			}
			
			if result.Summary == "" {
				t.Error("Summary should not be empty")
			}
		})
	}
}

// TestBulkStatusTransition_DryRun tests dry run functionality for status transitions
func TestBulkStatusTransition_DryRun(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	
	// Create test task
	testTask := &Gorev{
		ID:              "task-1",
		Baslik:         "Test Task",
		Durum:          constants.TaskStatusPending,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask
	
	testCases := []struct {
		name               string
		request            BulkStatusTransitionRequest
		expectedSuccessful int
		expectedFailed     int
		expectedWarnings   int
	}{
		{
			name: "Valid dry run",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"task-1"},
				NewStatus: constants.TaskStatusInProgress,
				DryRun:    true,
			},
			expectedSuccessful: 0, // Dry run doesn't execute
			expectedFailed:     0,
			expectedWarnings:   1, // Should have dry run warning
		},
		{
			name: "Invalid dry run - non-existent task",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"non-existent"},
				NewStatus: constants.TaskStatusInProgress,
				DryRun:    true,
			},
			expectedSuccessful: 0,
			expectedFailed:     1, // Should fail validation
			expectedWarnings:   0,
		},
		{
			name: "Invalid transition dry run",
			request: BulkStatusTransitionRequest{
				TaskIDs:   []string{"task-1"},
				NewStatus: "invalid-status", // Will be caught by validation first
				DryRun:    true,
			},
			expectedSuccessful: 0,
			expectedFailed:     0,
			expectedWarnings:   0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test dry run with invalid status should return error before processing
			if tc.request.NewStatus == "invalid-status" {
				_, err := bp.BulkStatusTransition(tc.request)
				if err == nil {
					t.Error("Expected error for invalid status but got none")
				}
				return
			}
			
			result, err := bp.BulkStatusTransition(tc.request)
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if len(result.Successful) != tc.expectedSuccessful {
				t.Errorf("Expected %d successful transitions, got %d", tc.expectedSuccessful, len(result.Successful))
			}
			
			if len(result.Failed) != tc.expectedFailed {
				t.Errorf("Expected %d failed transitions, got %d", tc.expectedFailed, len(result.Failed))
			}
			
			if len(result.Warnings) != tc.expectedWarnings {
				t.Errorf("Expected %d warnings, got %d", tc.expectedWarnings, len(result.Warnings))
			}
		})
	}
}

// TestBulkStatusTransition_WithAIContext tests bulk status transitions with AI context
func TestBulkStatusTransition_WithAIContext(t *testing.T) {
	vy := NewMockVeriYonetici()
	bp := NewBatchProcessor(vy)
	acy := YeniAIContextYonetici(vy)
	bp.SetAIContextManager(acy)
	
	// Create test task
	testTask := &Gorev{
		ID:              "task-1",
		Baslik:         "Test Task",
		Durum:          constants.TaskStatusPending,
		OlusturmaTarih: time.Now(),
	}
	vy.gorevler["task-1"] = testTask
	
	request := BulkStatusTransitionRequest{
		TaskIDs:   []string{"task-1"},
		NewStatus: constants.TaskStatusInProgress,
	}
	
	result, err := bp.BulkStatusTransition(request)
	
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	
	if len(result.Successful) != 1 {
		t.Errorf("Expected 1 successful transition, got %d", len(result.Successful))
	}
	
	// Verify AI interaction was recorded
	if len(vy.interactions) == 0 {
		t.Error("Expected AI interaction to be recorded")
	}
	
	if len(vy.interactions) > 0 {
		interaction := vy.interactions[0]
		if interaction.GorevID != "task-1" {
			t.Errorf("Expected interaction for task-1, got %s", interaction.GorevID)
		}
		if interaction.ActionType != "bulk_status_change" {
			t.Errorf("Expected action type 'bulk_status_change', got %s", interaction.ActionType)
		}
	}
}
