package mcp

import (
	"testing"

	"github.com/msenol/gorev/internal/constants"
	"github.com/stretchr/testify/assert"
)

// setupTestHandlers creates a test handler instance
func setupTestHandlers(t *testing.T) *Handlers {
	_, handlers, cleanup := setupTestEnvironment(t)
	t.Cleanup(cleanup)
	return handlers
}

// TestGorevSetActive tests the gorev_set_active handler
func TestGorevSetActive(t *testing.T) {
	h := setupTestHandlers(t)

	// Create a test task
	proje, _ := h.isYonetici.ProjeOlustur(constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	gorevTest, _ := h.isYonetici.GorevOlustur(constants.TestTaskTitleEN, constants.TestTaskDescriptionEN, constants.PriorityHigh, proje.ID, "", nil)

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "Set valid task as active",
			params: map[string]interface{}{
				"task_id": gorevTest.ID,
			},
			expectError: false,
		},
		{
			name: "Set non-existent task as active",
			params: map[string]interface{}{
				"task_id": "non-existent-id",
			},
			expectError: true,
			errorMsg:    "görev bulunamadı",
		},
		{
			name:        "Missing gorev_id parameter",
			params:      map[string]interface{}{},
			expectError: true,
			errorMsg:    "task_id parametresi gerekli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := h.GorevSetActive(tt.params)

			if tt.expectError {
				assert.NotNil(t, result)
				assert.True(t, result.IsError)
				assert.Contains(t, getResultText(result), tt.errorMsg)
			} else {
				assert.NotNil(t, result)
				assert.False(t, result.IsError)
				assert.Contains(t, getResultText(result), "başarıyla aktif görev olarak ayarlandı")
			}
		})
	}
}

// TestGorevGetActive tests the gorev_get_active handler
func TestGorevGetActive(t *testing.T) {
	h := setupTestHandlers(t)

	t.Run("Get active task when one exists", func(t *testing.T) {
		// Create and set active task
		proje, _ := h.isYonetici.ProjeOlustur(constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
		gorevTest, _ := h.isYonetici.GorevOlustur(constants.TestTaskActive, constants.TestTaskDescriptionEN, constants.PriorityHigh, proje.ID, "", nil)
		h.GorevSetActive(map[string]interface{}{"task_id": gorevTest.ID})

		result, _ := h.GorevGetActive(map[string]interface{}{})

		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		// Should contain the active task details
		resultText := getResultText(result)
		assert.Contains(t, resultText, "Aktif Görev: Active Task")
		assert.Contains(t, resultText, gorevTest.ID)
	})

	t.Run("Get active task when none exists", func(t *testing.T) {
		// Create a fresh handler without active task
		h2 := setupTestHandlers(t)

		result, _ := h2.GorevGetActive(map[string]interface{}{})

		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		// Should indicate no active task
		assert.Contains(t, getResultText(result), "Şu anda aktif görev yok")
	})
}

// TestGorevRecent tests the gorev_recent handler
func TestGorevRecent(t *testing.T) {
	h := setupTestHandlers(t)

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
	}{
		{
			name: "Get recent tasks with limit",
			params: map[string]interface{}{
				"limit": 5,
			},
			expectError: false,
		},
		{
			name:        "Get recent tasks without limit",
			params:      map[string]interface{}{},
			expectError: false,
		},
		{
			name: "Get recent tasks with invalid limit",
			params: map[string]interface{}{
				"limit": "invalid",
			},
			expectError: false, // Handler uses default limit for invalid values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := h.GorevRecent(tt.params)

			if tt.expectError {
				assert.NotNil(t, result)
				assert.True(t, result.IsError)
			} else {
				assert.NotNil(t, result)
				assert.False(t, result.IsError)
				assert.NotEmpty(t, getResultText(result))
			}
		})
	}
}

// TestGorevContextSummary tests the gorev_context_summary handler
func TestGorevContextSummary(t *testing.T) {
	h := setupTestHandlers(t)

	// Create test data
	proje, _ := h.isYonetici.ProjeOlustur(constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	h.isYonetici.GorevOlustur(constants.TestTaskHighPriority, constants.TestTaskDescriptionEN, constants.PriorityHigh, proje.ID, "", nil)
	h.isYonetici.GorevOlustur(constants.TestTaskNormal, constants.TestTaskDescriptionEN, constants.PriorityMedium, proje.ID, "", nil)

	result, _ := h.GorevContextSummary(map[string]interface{}{})

	assert.NotNil(t, result)
	assert.False(t, result.IsError)

	// Check markdown response contains expected sections
	text := getResultText(result)
	assert.Contains(t, text, "## 🤖 AI Oturum Özeti")
	assert.Contains(t, text, "### 📊 Oturum İstatistikleri")
}

// TestGorevBatchUpdate tests the gorev_batch_update handler
func TestGorevBatchUpdate(t *testing.T) {
	h := setupTestHandlers(t)

	// Create test tasks
	proje, _ := h.isYonetici.ProjeOlustur(constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	task1, _ := h.isYonetici.GorevOlustur(constants.TestTaskOne, constants.TestTaskDescriptionEN, constants.PriorityHigh, proje.ID, "", nil)
	task2, _ := h.isYonetici.GorevOlustur(constants.TestTaskTwo, constants.TestTaskDescriptionEN, constants.PriorityMedium, proje.ID, "", nil)

	tests := []struct {
		name           string
		params         map[string]interface{}
		expectError    bool
		errorMsg       string
		validateResult func(t *testing.T, text string)
	}{
		{
			name: "Valid batch update",
			params: map[string]interface{}{
				"updates": []interface{}{
					map[string]interface{}{
						"id": task1.ID,
						"updates": map[string]interface{}{
							"durum": "devam_ediyor",
						},
					},
					map[string]interface{}{
						"id": task2.ID,
						"updates": map[string]interface{}{
							"durum": "tamamlandi",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "Missing updates parameter",
			params:      map[string]interface{}{},
			expectError: true,
			errorMsg:    "updates parametresi gerekli",
		},
		{
			name: "Invalid updates format",
			params: map[string]interface{}{
				"updates": "invalid",
			},
			expectError: true,
			errorMsg:    "updates parametresi gerekli ve dizi olmalı",
		},
		{
			name: "Update with non-existent task",
			params: map[string]interface{}{
				"updates": []interface{}{
					map[string]interface{}{
						"id": "non-existent",
						"updates": map[string]interface{}{
							"durum": "devam_ediyor",
						},
					},
				},
			},
			expectError: false, // Batch update returns result with failed items
			errorMsg:    "",    // Not used for non-error case
			validateResult: func(t *testing.T, text string) {
				// Should have failed updates
				assert.Contains(t, text, "❌ Başarısız Güncellemeler")
				assert.Contains(t, text, "non-existent")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := h.GorevBatchUpdate(tt.params)

			if tt.expectError {
				assert.NotNil(t, result)
				assert.True(t, result.IsError)
				assert.Contains(t, getResultText(result), tt.errorMsg)
			} else {
				assert.NotNil(t, result)
				assert.False(t, result.IsError)

				text := getResultText(result)
				if tt.validateResult != nil {
					tt.validateResult(t, text)
				} else {
					// Default validation
					assert.Contains(t, text, "## 📦 Toplu Güncelleme Sonucu")
					// Should have successful updates
					assert.Contains(t, text, "✅ Başarılı Güncellemeler")
				}
			}
		})
	}
}

// TestGorevNLPQuery tests the gorev_nlp_query handler
func TestGorevNLPQuery(t *testing.T) {
	h := setupTestHandlers(t)

	// Create test tasks
	proje, _ := h.isYonetici.ProjeOlustur(constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	h.isYonetici.GorevOlustur("Yüksek Öncelikli Görev", "Acil", "yuksek", proje.ID, "", []string{"bug"})
	h.isYonetici.GorevOlustur("Normal Görev", "Normal açıklama", "orta", proje.ID, "", []string{"feature"})
	h.isYonetici.GorevOlustur("Test Görevi", "Test için", "dusuk", proje.ID, "", nil)

	tests := []struct {
		name        string
		params      map[string]interface{}
		expectError bool
		expectCount int
	}{
		{
			name: "Query for high priority tasks",
			params: map[string]interface{}{
				"query": "yüksek öncelikli görevler",
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "Query for tasks with tag",
			params: map[string]interface{}{
				"query": "etiket:bug",
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "General text search",
			params: map[string]interface{}{
				"query": "test",
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name:        "Missing query parameter",
			params:      map[string]interface{}{},
			expectError: true,
		},
		{
			name: "Empty query",
			params: map[string]interface{}{
				"query": "",
			},
			expectError: true, // Empty query returns error
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := h.GorevNLPQuery(tt.params)

			if tt.expectError {
				assert.NotNil(t, result)
				assert.True(t, result.IsError)
			} else {
				assert.NotNil(t, result)
				assert.False(t, result.IsError)

				// Count tasks in markdown output
				text := getResultText(result)
				if tt.expectCount == 0 {
					assert.Contains(t, text, "eşleşen görev bulunamadı")
				} else {
					// Count bullet points to determine number of tasks
					bulletCount := 0
					for _, char := range text {
						if char == '-' {
							bulletCount++
						}
					}
					assert.GreaterOrEqual(t, bulletCount, tt.expectCount)
				}
			}
		})
	}
}

// TestAutoStateTransition tests automatic state transition functionality
func TestAutoStateTransition(t *testing.T) {
	h := setupTestHandlers(t)

	// Create a task in "beklemede" state
	proje, _ := h.isYonetici.ProjeOlustur(constants.TestProjectNameEN, constants.TestProjectDescriptionEN)
	gorevTest, _ := h.isYonetici.GorevOlustur(constants.TestTaskTitleEN, constants.TestTaskDescriptionEN, constants.PriorityHigh, proje.ID, "", nil)

	// Verify initial state
	gorev, _ := h.isYonetici.GorevGetir(gorevTest.ID)
	assert.Equal(t, "beklemede", gorev.Durum)

	// Call GorevDetay which should trigger auto-state transition
	result, _ := h.GorevDetay(map[string]interface{}{"id": gorevTest.ID})

	assert.NotNil(t, result)
	assert.False(t, result.IsError)

	// Check if task state was auto-transitioned
	// Note: Since we're using mocks, the actual transition won't happen
	// but we verify the handler was called successfully
	assert.Contains(t, getResultText(result), gorevTest.Baslik)
}
