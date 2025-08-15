package gorev

import (
	"strings"
	"testing"
)

func TestNLPProcessor_ProcessQuery(t *testing.T) {
	nlp := NewNLPProcessor()

	tests := []struct {
		name           string
		query          string
		expectedAction string
		minConfidence  float64
	}{
		{
			name:           "Turkish task list",
			query:          "görevleri göster",
			expectedAction: "list",
			minConfidence:  0.7,
		},
		{
			name:           "English task list",
			query:          "show tasks",
			expectedAction: "list",
			minConfidence:  0.7,
		},
		{
			name:           "Turkish task creation",
			query:          "yeni görev oluştur: Frontend API entegrasyonu",
			expectedAction: "create",
			minConfidence:  0.7,
		},
		{
			name:           "English task creation",
			query:          "create task: Update user authentication",
			expectedAction: "create",
			minConfidence:  0.7,
		},
		{
			name:           "Turkish completion",
			query:          "görev #123 tamamla",
			expectedAction: "complete",
			minConfidence:  0.6,
		},
		{
			name:           "English completion",
			query:          "complete task #456",
			expectedAction: "complete",
			minConfidence:  0.6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent, err := nlp.ProcessQuery(tt.query)
			if err != nil {
				t.Fatalf("ProcessQuery failed: %v", err)
			}

			if intent.Action != tt.expectedAction {
				t.Errorf("Expected action %s, got %s", tt.expectedAction, intent.Action)
			}

			if intent.Confidence < tt.minConfidence {
				t.Errorf("Confidence too low: got %.2f, expected >= %.2f", 
					intent.Confidence, tt.minConfidence)
			}

			if intent.Raw != tt.query {
				t.Errorf("Raw query mismatch: got %s, expected %s", intent.Raw, tt.query)
			}
		})
	}
}

func TestNLPProcessor_ParseTimeExpressions(t *testing.T) {
	nlp := NewNLPProcessor()

	tests := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "Turkish today",
			query:    "bugün yapmam gereken görevler",
			expected: true,
		},
		{
			name:     "English today",
			query:    "tasks for today",
			expected: true,
		},
		{
			name:     "Turkish tomorrow",
			query:    "yarın deadline olan görevler",
			expected: true,
		},
		{
			name:     "English tomorrow",
			query:    "tasks due tomorrow",
			expected: true,
		},
		{
			name:     "Turkish this week",
			query:    "bu hafta tamamlanması gereken görevler",
			expected: true,
		},
		{
			name:     "English this week",
			query:    "tasks for this week",
			expected: true,
		},
		{
			name:     "Specific date",
			query:    "2025-12-25 tarihli görevler",
			expected: true,
		},
		{
			name:     "No time expression",
			query:    "tüm görevleri göster",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeRange := nlp.parseTimeExpressions(tt.query)
			hasTime := timeRange != nil

			if hasTime != tt.expected {
				t.Errorf("Expected time expression: %v, got: %v", tt.expected, hasTime)
			}

			if hasTime && timeRange.Start == nil {
				t.Error("Time range found but start time is nil")
			}
		})
	}
}

func TestNLPProcessor_ParseFilters(t *testing.T) {
	nlp := NewNLPProcessor()

	tests := []struct {
		name           string
		query          string
		expectedFilter string
		expectedValue  interface{}
	}{
		{
			name:           "Turkish high priority",
			query:          "yüksek öncelik görevleri",
			expectedFilter: "priority",
			expectedValue:  "high",
		},
		{
			name:           "English urgent",
			query:          "urgent tasks",
			expectedFilter: "priority",
			expectedValue:  "urgent",
		},
		{
			name:           "Turkish open status",
			query:          "açık görevleri göster",
			expectedFilter: "status",
			expectedValue:  "open",
		},
		{
			name:           "English completed status",
			query:          "show completed tasks",
			expectedFilter: "status",
			expectedValue:  "completed",
		},
		{
			name:           "Frontend category",
			query:          "frontend görevleri",
			expectedFilter: "category",
			expectedValue:  "frontend",
		},
		{
			name:           "Backend category",
			query:          "backend tasks",
			expectedFilter: "category",
			expectedValue:  "backend",
		},
		{
			name:           "Tag filter",
			query:          "etiket:bug görevleri",
			expectedFilter: "tags",
			expectedValue:  []string{"bug"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filters := nlp.parseFilters(tt.query)

			value, exists := filters[tt.expectedFilter]
			if !exists {
				t.Errorf("Expected filter %s not found", tt.expectedFilter)
				return
			}

			// Handle different value types
			switch expected := tt.expectedValue.(type) {
			case string:
				if value != expected {
					t.Errorf("Expected %s, got %v", expected, value)
				}
			case []string:
				actual, ok := value.([]string)
				if !ok {
					t.Errorf("Expected []string, got %T", value)
					return
				}
				if len(actual) != len(expected) {
					t.Errorf("Expected %v, got %v", expected, actual)
					return
				}
				for i, v := range expected {
					if actual[i] != v {
						t.Errorf("Expected %v, got %v", expected, actual)
						break
					}
				}
			}
		})
	}
}

func TestNLPProcessor_ParseTaskReferences(t *testing.T) {
	nlp := NewNLPProcessor()

	tests := []struct {
		name      string
		query     string
		expected  []string
	}{
		{
			name:     "Turkish task ID",
			query:    "görev #123 tamamla",
			expected: []string{"id:123"},
		},
		{
			name:     "English task ID",
			query:    "complete task #456",
			expected: []string{"id:456"},
		},
		{
			name:     "Turkish recent task",
			query:    "son oluşturduğum görev",
			expected: []string{"recent:1"},
		},
		{
			name:     "English recent task",
			query:    "latest task",
			expected: []string{"recent:1"},
		},
		{
			name:     "Title reference",
			query:    "\"API Integration\" görevini güncelle",
			expected: []string{"title:API Integration"},
		},
		{
			name:     "Multiple IDs",
			query:    "görev #123 ve #456 tamamla",
			expected: []string{"id:123", "id:456"},
		},
		{
			name:     "No references",
			query:    "tüm görevleri listele",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refs := nlp.parseTaskReferences(tt.query)

			if len(refs) != len(tt.expected) {
				t.Errorf("Expected %d references, got %d", len(tt.expected), len(refs))
				return
			}

			for i, expected := range tt.expected {
				if refs[i] != expected {
					t.Errorf("Expected %s, got %s", expected, refs[i])
				}
			}
		})
	}
}

func TestNLPProcessor_ExtractTaskContent(t *testing.T) {
	nlp := NewNLPProcessor()

	tests := []struct {
		name           string
		query          string
		expectedTitle  string
		expectedDesc   string
		expectDueDate  bool
	}{
		{
			name:          "Turkish simple task",
			query:         "yeni görev oluştur: API entegrasyonu",
			expectedTitle: "API entegrasyonu",
			expectedDesc:  "",
		},
		{
			name:          "English simple task",
			query:         "create task: Update documentation",
			expectedTitle: "Update documentation",
			expectedDesc:  "",
		},
		{
			name:          "Turkish task with description",
			query:         "yeni görev: Frontend - Kullanıcı arayüzü geliştirme",
			expectedTitle: "Frontend",
			expectedDesc:  "Kullanıcı arayüzü geliştirme",
		},
		{
			name:          "Task with colon separator",
			query:         "görev oluştur: Backend API: Database connection optimizasyonu",
			expectedTitle: "Backend API",
			expectedDesc:  "Database connection optimizasyonu",
		},
		{
			name:          "Task with due date",
			query:         "yarın deadline: Test senaryoları yazma",
			expectedTitle: "Test senaryoları yazma",
			expectDueDate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := nlp.ExtractTaskContent(tt.query)

			if title, ok := content["title"].(string); ok {
				if title != tt.expectedTitle {
					t.Errorf("Expected title %s, got %s", tt.expectedTitle, title)
				}
			} else if tt.expectedTitle != "" {
				t.Error("Expected title but got none")
			}

			if tt.expectedDesc != "" {
				if desc, ok := content["description"].(string); ok {
					if desc != tt.expectedDesc {
						t.Errorf("Expected description %s, got %s", tt.expectedDesc, desc)
					}
				} else {
					t.Error("Expected description but got none")
				}
			}

			if tt.expectDueDate {
				if _, ok := content["due_date"]; !ok {
					t.Error("Expected due date but got none")
				}
			}
		})
	}
}

func TestNLPProcessor_FormatResponse(t *testing.T) {
	nlp := NewNLPProcessor()

	tests := []struct {
		name     string
		action   string
		results  interface{}
		lang     string
		contains string
	}{
		{
			name:     "Turkish task creation",
			action:   "create",
			results:  "API Integration",
			lang:     "tr",
			contains: "başarıyla oluşturuldu",
		},
		{
			name:     "English task creation",
			action:   "create",
			results:  "API Integration",
			lang:     "en",
			contains: "created successfully",
		},
		{
			name:     "Turkish empty list",
			action:   "list",
			results:  nil,
			lang:     "tr",
			contains: "bulunamadı",
		},
		{
			name:     "English empty list",
			action:   "list",
			results:  nil,
			lang:     "en",
			contains: "No tasks found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := nlp.FormatResponse(tt.action, tt.results, tt.lang)

			if !strings.Contains(response, tt.contains) {
				t.Errorf("Expected response to contain %s, got: %s", tt.contains, response)
			}
		})
	}
}

func TestNLPProcessor_ValidateIntent(t *testing.T) {
	nlp := NewNLPProcessor()

	tests := []struct {
		name        string
		intent      *QueryIntent
		expectError bool
	}{
		{
			name: "Valid high confidence intent",
			intent: &QueryIntent{
				Action:     "list",
				Confidence: 0.8,
				Raw:        "görevleri göster",
			},
			expectError: false,
		},
		{
			name: "Low confidence intent",
			intent: &QueryIntent{
				Action:     "list",
				Confidence: 0.2,
				Raw:        "unclear query",
			},
			expectError: true,
		},
		{
			name: "No action identified",
			intent: &QueryIntent{
				Action:     "",
				Confidence: 0.8,
				Raw:        "random text",
			},
			expectError: true,
		},
		{
			name: "Update without task reference",
			intent: &QueryIntent{
				Action:     "update",
				Confidence: 0.8,
				Raw:        "update something",
				Parameters: map[string]interface{}{},
			},
			expectError: true,
		},
		{
			name: "Update with task reference",
			intent: &QueryIntent{
				Action:     "update",
				Confidence: 0.8,
				Raw:        "update task #123",
				Parameters: map[string]interface{}{
					"task_references": []string{"id:123"},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := nlp.ValidateIntent(tt.intent)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestNLPProcessor_BuildQuery(t *testing.T) {
	nlp := NewNLPProcessor()

	intent := &QueryIntent{
		Action: "list",
		Filters: map[string]interface{}{
			"status":   "open",
			"priority": "high",
		},
		TimeRange: &TimeRange{
			Relative: "today",
		},
		Confidence: 0.8,
	}

	query := nlp.BuildQuery(intent)

	// Verify all fields are included
	if query["action"] != "list" {
		t.Errorf("Expected action 'list', got %v", query["action"])
	}

	if filters, ok := query["filters"].(map[string]interface{}); ok {
		if filters["status"] != "open" {
			t.Errorf("Expected status 'open', got %v", filters["status"])
		}
		if filters["priority"] != "high" {
			t.Errorf("Expected priority 'high', got %v", filters["priority"])
		}
	} else {
		t.Error("Expected filters to be present")
	}

	if timeRange, ok := query["time_range"].(*TimeRange); ok {
		if timeRange.Relative != "today" {
			t.Errorf("Expected relative 'today', got %v", timeRange.Relative)
		}
	} else {
		t.Error("Expected time_range to be present")
	}

	if query["confidence"] != 0.8 {
		t.Errorf("Expected confidence 0.8, got %v", query["confidence"])
	}
}

// Benchmark tests for performance
func BenchmarkNLPProcessor_ProcessQuery(b *testing.B) {
	nlp := NewNLPProcessor()
	query := "bugün tamamlanması gereken yüksek öncelik görevleri göster"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := nlp.ProcessQuery(query)
		if err != nil {
			b.Fatalf("ProcessQuery failed: %v", err)
		}
	}
}

func BenchmarkNLPProcessor_ParseTimeExpressions(b *testing.B) {
	nlp := NewNLPProcessor()
	query := "yarın deadline olan frontend görevleri listele"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nlp.parseTimeExpressions(query)
	}
}

// Helper function to create test intent
func createTestIntent(action string, confidence float64) *QueryIntent {
	return &QueryIntent{
		Action:     action,
		Filters:    make(map[string]interface{}),
		Parameters: make(map[string]interface{}),
		Confidence: confidence,
		Raw:        "test query",
	}
}

// Test for thread safety
func TestNLPProcessor_ConcurrentAccess(t *testing.T) {
	nlp := NewNLPProcessor()
	
	// Test concurrent query processing
	queries := []string{
		"görevleri göster",
		"yeni görev oluştur",
		"görev #123 tamamla",
		"bugün yapmam gereken görevler",
		"high priority tasks",
	}

	results := make(chan error, len(queries))
	
	for _, query := range queries {
		go func(q string) {
			_, err := nlp.ProcessQuery(q)
			results <- err
		}(query)
	}

	// Collect results
	for i := 0; i < len(queries); i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent query failed: %v", err)
		}
	}
}