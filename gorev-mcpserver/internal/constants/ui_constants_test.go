package constants

import (
	"testing"
)

func TestGetStatusEmoji(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Pending status", TaskStatusPending, EmojiStatusPending},
		{"In progress status", TaskStatusInProgress, EmojiStatusInProgress},
		{"Completed status", TaskStatusCompleted, EmojiStatusCompleted},
		{"Cancelled status", TaskStatusCancelled, EmojiStatusCancelled},
		{"Unknown status", "unknown", EmojiStatusUnknown},
		{"Empty status", "", EmojiStatusUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStatusEmoji(tt.status)
			if result != tt.expected {
				t.Errorf("GetStatusEmoji(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetPriorityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		expected string
	}{
		{"High priority", PriorityHigh, EmojiPriorityHigh},
		{"Medium priority", PriorityMedium, EmojiPriorityMedium},
		{"Low priority", PriorityLow, EmojiPriorityLow},
		{"Unknown priority", "unknown", EmojiPriorityUnknown},
		{"Empty priority", "", EmojiPriorityUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPriorityEmoji(tt.priority)
			if result != tt.expected {
				t.Errorf("GetPriorityEmoji(%q) = %v, want %v", tt.priority, result, tt.expected)
			}
		})
	}
}

func TestGetStatusDisplay(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{"Valid pending status", TaskStatusPending, StatusDisplayMap[TaskStatusPending]},
		{"Valid in progress status", TaskStatusInProgress, StatusDisplayMap[TaskStatusInProgress]},
		{"Valid completed status", TaskStatusCompleted, StatusDisplayMap[TaskStatusCompleted]},
		{"Valid cancelled status", TaskStatusCancelled, StatusDisplayMap[TaskStatusCancelled]},
		{"Unknown status", "unknown", EmojiStatusUnknown + " Bilinmeyen"},
		{"Empty status", "", EmojiStatusUnknown + " Bilinmeyen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStatusDisplay(tt.status)
			if result != tt.expected {
				t.Errorf("GetStatusDisplay(%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetPriorityDisplay(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		expected string
	}{
		{"Valid high priority", PriorityHigh, PriorityDisplayMap[PriorityHigh]},
		{"Valid medium priority", PriorityMedium, PriorityDisplayMap[PriorityMedium]},
		{"Valid low priority", PriorityLow, PriorityDisplayMap[PriorityLow]},
		{"Unknown priority", "unknown", EmojiPriorityUnknown + " Bilinmeyen"},
		{"Empty priority", "", EmojiPriorityUnknown + " Bilinmeyen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPriorityDisplay(tt.priority)
			if result != tt.expected {
				t.Errorf("GetPriorityDisplay(%q) = %v, want %v", tt.priority, result, tt.expected)
			}
		})
	}
}

func TestGetSuggestionPriorityEmoji(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		expected string
	}{
		{"High priority suggestion", "high", EmojiSuggestionHigh},
		{"Medium priority suggestion", "medium", EmojiSuggestionMedium},
		{"Low priority suggestion", "low", EmojiSuggestionLow},
		{"Unknown priority suggestion", "unknown", EmojiPriorityUnknown},
		{"Empty priority suggestion", "", EmojiPriorityUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSuggestionPriorityEmoji(tt.priority)
			if result != tt.expected {
				t.Errorf("GetSuggestionPriorityEmoji(%q) = %v, want %v", tt.priority, result, tt.expected)
			}
		})
	}
}

func TestStatusDisplayMapConsistency(t *testing.T) {
	// Test that all valid task statuses have display mappings
	validStatuses := GetValidTaskStatuses()
	for _, status := range validStatuses {
		if _, exists := StatusDisplayMap[status]; !exists {
			t.Errorf("StatusDisplayMap missing entry for valid status: %s", status)
		}
	}
}

func TestPriorityDisplayMapConsistency(t *testing.T) {
	// Test that all valid priorities have display mappings
	validPriorities := GetValidPriorities()
	for _, priority := range validPriorities {
		if _, exists := PriorityDisplayMap[priority]; !exists {
			t.Errorf("PriorityDisplayMap missing entry for valid priority: %s", priority)
		}
	}
}
