package gorev

import (
	"testing"
	"time"
)

// TestBulkTagOperation tests the BulkTagOperation functionality
func TestBatchProcessor_BulkTagOperation(t *testing.T) {
	// Create mock VeriYonetici
	vy := NewMockVeriYonetici()
	
	// Create BatchProcessor
	bp := NewBatchProcessor(vy)
	
	// Create AI Context Manager
	acy := YeniAIContextYonetici(vy)
	bp.SetAIContextManager(acy)
	
	// Create test tasks
	task1 := &Gorev{
		ID:              "task-1",
		Baslik:          "Task 1",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now().Add(-1 * time.Hour),
		GuncellemeTarih: time.Now().Add(-1 * time.Hour),
		Etiketler:       []*Etiket{},
	}
	
	task2 := &Gorev{
		ID:              "task-2",
		Baslik:          "Task 2",
		Durum:           "devam_ediyor",
		Oncelik:         "yuksek",
		OlusturmaTarih:  time.Now().Add(-2 * time.Hour),
		GuncellemeTarih: time.Now().Add(-2 * time.Hour),
		Etiketler:       []*Etiket{},
	}
	
	task3 := &Gorev{
		ID:              "task-3",
		Baslik:          "Task 3",
		Durum:           "tamamlandi",
		Oncelik:         "dusuk",
		OlusturmaTarih:  time.Now().Add(-3 * time.Hour),
		GuncellemeTarih: time.Now().Add(-3 * time.Hour),
		Etiketler:       []*Etiket{},
	}
	
	// Add existing tag to task3
	existingTag := &Etiket{
		ID:   "tag-existing",
		Isim: "existing",
	}
	task3.Etiketler = []*Etiket{existingTag}
	vy.tags["existing"] = existingTag
	
	// Add tasks to mock
	vy.gorevler["task-1"] = task1
	vy.gorevler["task-2"] = task2
	vy.gorevler["task-3"] = task3
	
	// Test cases for BulkTagOperation
	testCases := []struct {
		name           string
		request        BulkTagOperationRequest
		expectedResult struct {
			successful  int
			failed      int
			warnings    int
			errorResult bool
		}
		description string
	}{
		{
			name: "Add tags to tasks",
			request: BulkTagOperationRequest{
				TaskIDs:   []string{"task-1", "task-2"},
				Tags:      []string{"important", "frontend"},
				Operation: "add",
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  2,
				failed:      0,
				warnings:    0,
				errorResult: false,
			},
			description: "Should add tags to tasks that don't have them",
		},
		{
			name: "Remove tags from tasks",
			request: BulkTagOperationRequest{
				TaskIDs:   []string{"task-3"},
				Tags:      []string{"existing"},
				Operation: "remove",
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  1,
				failed:      0,
				warnings:    0,
				errorResult: false,
			},
			description: "Should remove existing tags from tasks",
		},
		{
			name: "Replace tags on tasks",
			request: BulkTagOperationRequest{
				TaskIDs:   []string{"task-1", "task-3"},
				Tags:      []string{"replaced", "new-tag"},
				Operation: "replace",
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  2,
				failed:      0,
				warnings:    0,
				errorResult: false,
			},
			description: "Should replace all existing tags with new ones",
		},
		{
			name: "No changes needed",
			request: BulkTagOperationRequest{
				TaskIDs:   []string{"task-1"},
				Tags:      []string{"replaced", "new-tag"},
				Operation: "replace",
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      0,
				warnings:    1,
				errorResult: false,
			},
			description: "Should skip task if tags are already set",
		},
		{
			name: "Invalid operation",
			request: BulkTagOperationRequest{
				TaskIDs:   []string{"task-1", "task-2"},
				Tags:      []string{"tag1", "tag2"},
				Operation: "invalid",
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      0,
				warnings:    0,
				errorResult: true,
			},
			description: "Should return error for invalid operation",
		},
		{
			name: "Dry run operation",
			request: BulkTagOperationRequest{
				TaskIDs:   []string{"task-1", "task-2"},
				Tags:      []string{"dry-run-tag"},
				Operation: "add",
				DryRun:    true,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      0,
				warnings:    2,
				errorResult: false,
			},
			description: "Should not modify tasks but return warnings in dry run mode",
		},
		{
			name: "Non-existent task",
			request: BulkTagOperationRequest{
				TaskIDs:   []string{"non-existent"},
				Tags:      []string{"tag1"},
				Operation: "add",
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      1,
				warnings:    0,
				errorResult: false,
			},
			description: "Should handle non-existent tasks gracefully",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset task tags for consistent test state
			if tc.name != "No changes needed" {
				// Reset task1 and task2 tags
				task1.Etiketler = []*Etiket{}
				task2.Etiketler = []*Etiket{}
				// Ensure task3 has the existing tag
				task3.Etiketler = []*Etiket{existingTag}
			}
			
			result, err := bp.BulkTagOperation(tc.request)
			
			if tc.expectedResult.errorResult {
				if err == nil {
					t.Errorf("Expected error but got none for %s", tc.description)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.description, err)
				return
			}
			
			// Check result counts
			if len(result.Successful) != tc.expectedResult.successful {
				t.Errorf("Expected %d successful operations, got %d for %s", 
					tc.expectedResult.successful, len(result.Successful), tc.description)
			}
			
			if len(result.Failed) != tc.expectedResult.failed {
				t.Errorf("Expected %d failed operations, got %d for %s", 
					tc.expectedResult.failed, len(result.Failed), tc.description)
			}
			
			if len(result.Warnings) != tc.expectedResult.warnings {
				t.Errorf("Expected %d warnings, got %d for %s", 
					tc.expectedResult.warnings, len(result.Warnings), tc.description)
			}
			
			// Verify task tags were updated correctly
			if !tc.request.DryRun && tc.expectedResult.errorResult == false {
				switch tc.request.Operation {
				case "add":
					// Check that tags were added for successful tasks
					for _, taskID := range result.Successful {
						task := vy.gorevler[taskID]
						for _, tagName := range tc.request.Tags {
							found := false
							for _, tag := range task.Etiketler {
								if tag.Isim == tagName {
									found = true
									break
								}
							}
							if !found {
								t.Errorf("Tag %s was not added to task %s for %s", 
									tagName, taskID, tc.description)
							}
						}
					}
					
				case "remove":
					// Check that tags were removed for successful tasks
					for _, taskID := range result.Successful {
						task := vy.gorevler[taskID]
						for _, tagName := range tc.request.Tags {
							for _, tag := range task.Etiketler {
								if tag.Isim == tagName {
									t.Errorf("Tag %s was not removed from task %s for %s", 
										tagName, taskID, tc.description)
								}
							}
						}
					}
					
				case "replace":
					// Check that tags were replaced for successful tasks
					for _, taskID := range result.Successful {
						task := vy.gorevler[taskID]
						
						// Check tag count matches
						if len(task.Etiketler) != len(tc.request.Tags) {
							t.Errorf("Expected %d tags, got %d for task %s in %s", 
								len(tc.request.Tags), len(task.Etiketler), taskID, tc.description)
						}
						
						// Check each tag exists
						for _, tagName := range tc.request.Tags {
							found := false
							for _, tag := range task.Etiketler {
								if tag.Isim == tagName {
									found = true
									break
								}
							}
							if !found {
								t.Errorf("Tag %s was not found in replaced tags for task %s in %s", 
									tagName, taskID, tc.description)
							}
						}
					}
				}
			}
		})
	}
}

// TestBulkDelete tests the BulkDelete functionality
func TestBatchProcessor_BulkDelete(t *testing.T) {
	// Create mock VeriYonetici
	vy := NewMockVeriYonetici()
	
	// Create BatchProcessor
	bp := NewBatchProcessor(vy)
	
	// Create AI Context Manager
	acy := YeniAIContextYonetici(vy)
	bp.SetAIContextManager(acy)
	
	// Create test tasks
	task1 := &Gorev{
		ID:              "task-1",
		Baslik:          "Task 1",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now().Add(-1 * time.Hour),
		GuncellemeTarih: time.Now().Add(-1 * time.Hour),
	}
	
	task2 := &Gorev{
		ID:              "task-2",
		Baslik:          "Task 2",
		Durum:           "devam_ediyor",
		Oncelik:         "yuksek",
		OlusturmaTarih:  time.Now().Add(-2 * time.Hour),
		GuncellemeTarih: time.Now().Add(-2 * time.Hour),
	}
	
	// Parent task with subtasks
	parentTask := &Gorev{
		ID:              "parent-task",
		Baslik:          "Parent Task",
		Durum:           "beklemede",
		Oncelik:         "orta",
		OlusturmaTarih:  time.Now().Add(-3 * time.Hour),
		GuncellemeTarih: time.Now().Add(-3 * time.Hour),
	}
	
	// Child task
	childTask := &Gorev{
		ID:              "child-task",
		Baslik:          "Child Task",
		Durum:           "beklemede",
		Oncelik:         "orta",
		ParentID:        "parent-task",
		OlusturmaTarih:  time.Now().Add(-3 * time.Hour),
		GuncellemeTarih: time.Now().Add(-3 * time.Hour),
	}
	
	// Add tasks to mock
	vy.gorevler["task-1"] = task1
	vy.gorevler["task-2"] = task2
	vy.gorevler["parent-task"] = parentTask
	vy.gorevler["child-task"] = childTask
	
	// Test cases for BulkDelete
	testCases := []struct {
		name           string
		request        BulkDeleteRequest
		expectedResult struct {
			successful  int
			failed      int
			warnings    int
			errorResult bool
		}
		description string
	}{
		{
			name: "Delete tasks without confirmation",
			request: BulkDeleteRequest{
				TaskIDs:      []string{"task-1", "task-2"},
				Confirmation: "WRONG CONFIRMATION",
				Force:        false,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      0,
				warnings:    0,
				errorResult: true,
			},
			description: "Should fail without proper confirmation",
		},
		{
			name: "Delete tasks with confirmation",
			request: BulkDeleteRequest{
				TaskIDs:      []string{"task-1", "task-2"},
				Confirmation: "DELETE 2 TASKS",
				Force:        false,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  2,
				failed:      0,
				warnings:    0,
				errorResult: false,
			},
			description: "Should successfully delete tasks with proper confirmation",
		},
		{
			name: "Force delete without confirmation",
			request: BulkDeleteRequest{
				TaskIDs:      []string{"task-1"},
				Confirmation: "WRONG",
				Force:        true,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  1,
				failed:      0,
				warnings:    0,
				errorResult: false,
			},
			description: "Should delete tasks without confirmation when force is true",
		},
		{
			name: "Delete parent with child without delete_subtasks flag",
			request: BulkDeleteRequest{
				TaskIDs:        []string{"parent-task"},
				Confirmation:   "DELETE 1 TASKS",
				DeleteSubtasks: false,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      1,
				warnings:    0,
				errorResult: false,
			},
			description: "Should fail to delete parent task when it has subtasks and delete_subtasks is false",
		},
		{
			name: "Delete parent with delete_subtasks flag",
			request: BulkDeleteRequest{
				TaskIDs:        []string{"parent-task"},
				Confirmation:   "DELETE 1 TASKS",
				DeleteSubtasks: true,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  1,
				failed:      0,
				warnings:    0,
				errorResult: false,
			},
			description: "Should delete parent task when delete_subtasks is true",
		},
		{
			name: "Dry run deletion",
			request: BulkDeleteRequest{
				TaskIDs:      []string{"task-1", "task-2"},
				Confirmation: "DELETE 2 TASKS",
				DryRun:       true,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      0,
				warnings:    2,
				errorResult: false,
			},
			description: "Should not delete tasks in dry run mode",
		},
		{
			name: "Delete non-existent task",
			request: BulkDeleteRequest{
				TaskIDs:      []string{"non-existent"},
				Confirmation: "DELETE 1 TASKS",
				Force:        true,
			},
			expectedResult: struct {
				successful  int
				failed      int
				warnings    int
				errorResult bool
			}{
				successful:  0,
				failed:      1,
				warnings:    0,
				errorResult: false,
			},
			description: "Should handle non-existent tasks gracefully",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset tasks for each test
			vy.gorevler["task-1"] = task1
			vy.gorevler["task-2"] = task2
			vy.gorevler["parent-task"] = parentTask
			vy.gorevler["child-task"] = childTask
			
			result, err := bp.BulkDelete(tc.request)
			
			if tc.expectedResult.errorResult {
				if err == nil {
					t.Errorf("Expected error but got none for %s", tc.description)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tc.description, err)
				return
			}
			
			// Check result counts
			if len(result.Successful) != tc.expectedResult.successful {
				t.Errorf("Expected %d successful operations, got %d for %s", 
					tc.expectedResult.successful, len(result.Successful), tc.description)
			}
			
			if len(result.Failed) != tc.expectedResult.failed {
				t.Errorf("Expected %d failed operations, got %d for %s", 
					tc.expectedResult.failed, len(result.Failed), tc.description)
			}
			
			if len(result.Warnings) != tc.expectedResult.warnings {
				t.Errorf("Expected %d warnings, got %d for %s", 
					tc.expectedResult.warnings, len(result.Warnings), tc.description)
			}
			
			// Verify tasks were actually deleted
			if !tc.request.DryRun && tc.expectedResult.errorResult == false {
				for _, taskID := range tc.request.TaskIDs {
					_, exists := vy.gorevler[taskID]
					shouldExist := true
					
					// Check if this task should have been deleted
					for _, successfulID := range result.Successful {
						if successfulID == taskID {
							shouldExist = false
							break
						}
					}
					
					// Skip verification for non-existent tasks since they were never added to the mock
					if taskID == "non-existent" {
						continue
					}
					
					if shouldExist && !exists {
						t.Errorf("Task %s was unexpectedly deleted for %s", taskID, tc.description)
					} else if !shouldExist && exists {
						t.Errorf("Task %s was not deleted as expected for %s", taskID, tc.description)
					}
				}
			}
		})
	}
}