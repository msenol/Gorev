package constants

import (
	"testing"
)

func TestIsValidTaskStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Valid pending status", TaskStatusPending, true},
		{"Valid in progress status", TaskStatusInProgress, true},
		{"Valid completed status", TaskStatusCompleted, true},
		{"Valid cancelled status", TaskStatusCancelled, true},
		{"Invalid status", "invalid_status", false},
		{"Empty status", "", false},
		{"Case sensitive test", "PENDING", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidTaskStatus(tt.status)
			if result != tt.expected {
				t.Errorf("IsValidTaskStatus(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestIsValidPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		expected bool
	}{
		{"Valid high priority", PriorityHigh, true},
		{"Valid medium priority", PriorityMedium, true},
		{"Valid low priority", PriorityLow, true},
		{"Invalid priority", "invalid_priority", false},
		{"Empty priority", "", false},
		{"Case sensitive test", "HIGH", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidPriority(tt.priority)
			if result != tt.expected {
				t.Errorf("IsValidPriority(%q) = %v, want %v", tt.priority, result, tt.expected)
			}
		})
	}
}

func TestIsValidDependencyType(t *testing.T) {
	tests := []struct {
		name     string
		depType  string
		expected bool
	}{
		{"Valid blocker dependency", DependencyTypeBlocker, true},
		{"Valid depends_on dependency", DependencyTypeDependsOn, true},
		{"Invalid dependency type", "invalid_type", false},
		{"Empty dependency type", "", false},
		{"Case sensitive test", "BLOCKER", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidDependencyType(tt.depType)
			if result != tt.expected {
				t.Errorf("IsValidDependencyType(%q) = %v, want %v", tt.depType, result, tt.expected)
			}
		})
	}
}

func TestGetValidTaskStatuses(t *testing.T) {
	statuses := GetValidTaskStatuses()

	// Check that it contains expected statuses
	expectedStatuses := []string{TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted, TaskStatusCancelled}
	if len(statuses) != len(expectedStatuses) {
		t.Errorf("GetValidTaskStatuses() returned %d statuses, expected %d", len(statuses), len(expectedStatuses))
	}

	// Check each expected status is present
	for _, expected := range expectedStatuses {
		found := false
		for _, status := range statuses {
			if status == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetValidTaskStatuses() missing expected status: %s", expected)
		}
	}
}

func TestGetValidPriorities(t *testing.T) {
	priorities := GetValidPriorities()

	// Check that it contains expected priorities
	expectedPriorities := []string{PriorityHigh, PriorityMedium, PriorityLow}
	if len(priorities) != len(expectedPriorities) {
		t.Errorf("GetValidPriorities() returned %d priorities, expected %d", len(priorities), len(expectedPriorities))
	}

	// Check each expected priority is present
	for _, expected := range expectedPriorities {
		found := false
		for _, priority := range priorities {
			if priority == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetValidPriorities() missing expected priority: %s", expected)
		}
	}
}

func TestGetValidDependencyTypes(t *testing.T) {
	depTypes := GetValidDependencyTypes()

	// Check that it contains expected dependency types
	expectedTypes := []string{DependencyTypeBlocker, DependencyTypeDependsOn}
	if len(depTypes) != len(expectedTypes) {
		t.Errorf("GetValidDependencyTypes() returned %d types, expected %d", len(depTypes), len(expectedTypes))
	}

	// Check each expected type is present
	for _, expected := range expectedTypes {
		found := false
		for _, depType := range depTypes {
			if depType == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetValidDependencyTypes() missing expected type: %s", expected)
		}
	}
}
