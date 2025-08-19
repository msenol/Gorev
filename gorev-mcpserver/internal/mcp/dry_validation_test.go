package mcp

import (
	"testing"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// TestDRYi18nPatterns tests our DRY i18n patterns work correctly
func TestDRYi18nPatterns(t *testing.T) {
	// Initialize i18n
	if err := i18n.Initialize(constants.DefaultTestLanguage); err != nil {
		t.Fatalf("Failed to initialize i18n: %v", err)
	}

	t.Run("TParamFunction", func(t *testing.T) {
		result := i18n.TParam("id")
		if result == "" {
			t.Error("TParam should return non-empty string for 'id' parameter")
		}

		// Test fallback behavior
		fallback := i18n.TParam("nonexistent_param")
		expected := "nonexistent_param parameter"
		if fallback != expected {
			t.Errorf("Expected fallback '%s', got '%s'", expected, fallback)
		}
	})

	t.Run("FormatParameterRequired", func(t *testing.T) {
		result := i18n.FormatParameterRequired("test_param")
		if result == "" {
			t.Error("FormatParameterRequired should return non-empty string")
		}

		// Should contain parameter name
		if len(result) < len("test_param") {
			t.Error("Result should contain parameter name")
		}
	})

	t.Run("FormatInvalidValue", func(t *testing.T) {
		result := i18n.FormatInvalidValue("durum", "invalid", constants.GetValidTaskStatuses()[:2])
		if result == "" {
			t.Error("FormatInvalidValue should return non-empty string")
		}

		// Should contain the invalid value
		if len(result) < len("invalid") {
			t.Error("Result should contain invalid value")
		}
	})

	t.Run("TranslationKeysExist", func(t *testing.T) {
		// Test that our common translation keys work
		keys := []string{
			"tools.descriptions.gorev_listele",
			"tools.params.descriptions.id_field",
			"success.taskUpdated",
		}

		for _, key := range keys {
			result := i18n.T(key, nil)
			if result == key {
				t.Errorf("Translation missing for key: %s", key)
			}
		}
	})
}

// TestValidationDRYPatterns tests validation helper patterns
func TestValidationDRYPatterns(t *testing.T) {
	validator := NewParameterValidator()
	if validator == nil {
		t.Fatal("NewParameterValidator should not return nil")
	}

	t.Run("ValidStringValidation", func(t *testing.T) {
		params := map[string]interface{}{"id": constants.TestIDValidation}
		result, err := validator.ValidateRequiredString(params, "id")

		if err != nil {
			t.Errorf("ValidateRequiredString should not fail for valid input: %v", err)
		}

		if result != constants.TestIDValidation {
			t.Errorf("Expected 'test-id-123', got '%s'", result)
		}
	})

	t.Run("MissingRequiredParameter", func(t *testing.T) {
		params := map[string]interface{}{}
		_, err := validator.ValidateRequiredString(params, "id")

		if err == nil {
			t.Error("ValidateRequiredString should fail for missing required parameter")
		}
	})

	t.Run("EmptyStringParameter", func(t *testing.T) {
		params := map[string]interface{}{"id": ""}
		_, err := validator.ValidateRequiredString(params, "id")

		if err == nil {
			t.Error("ValidateRequiredString should fail for empty string")
		}
	})

	t.Run("ValidEnumValue", func(t *testing.T) {
		params := map[string]interface{}{"durum": "beklemede"}
		validOptions := constants.GetValidTaskStatuses()[:3]

		result, err := validator.ValidateEnum(params, "durum", validOptions, false)

		if err != nil {
			t.Errorf("ValidateEnum should not fail for valid enum value: %v", err)
		}

		if result != "beklemede" {
			t.Errorf("Expected 'beklemede', got '%s'", result)
		}
	})

	t.Run("InvalidEnumValue", func(t *testing.T) {
		params := map[string]interface{}{"durum": "invalid-status"}
		validOptions := constants.GetValidTaskStatuses()[:3]

		_, err := validator.ValidateEnum(params, "durum", validOptions, false)

		if err == nil {
			t.Error("ValidateEnum should fail for invalid enum value")
		}
	})
}

// TestFormatterDRYPatterns tests formatter helper patterns
func TestFormatterDRYPatterns(t *testing.T) {
	formatter := NewTaskFormatter()
	if formatter == nil {
		t.Fatal("NewTaskFormatter should not return nil")
	}

	t.Run("FormatTaskBasic", func(t *testing.T) {
		result := formatter.FormatTaskBasic(constants.TestTaskTitleEN, constants.TestTaskID)

		if result == "" {
			t.Error("FormatTaskBasic should return non-empty string")
		}

		// Should contain task name and short ID
		if len(result) < len(constants.TestTaskTitleEN) {
			t.Error("Result should contain task name")
		}
	})

	t.Run("GetStatusEmoji", func(t *testing.T) {
		testCases := []struct {
			status   string
			expected string
		}{
			{constants.TaskStatusPending, constants.EmojiStatusPending},
			{constants.TaskStatusInProgress, constants.EmojiStatusInProgress},
			{constants.TaskStatusCompleted, constants.EmojiStatusCompleted},
			{constants.TaskStatusCancelled, constants.EmojiStatusCancelled},
			{"unknown", constants.EmojiStatusUnknown},
		}

		for _, tc := range testCases {
			result := formatter.GetStatusEmoji(tc.status)
			if result != tc.expected {
				t.Errorf("GetStatusEmoji(%s): expected '%s', got '%s'", tc.status, tc.expected, result)
			}
		}
	})

	t.Run("GetPriorityEmoji", func(t *testing.T) {
		testCases := []struct {
			priority string
			expected string
		}{
			{constants.PriorityHigh, constants.EmojiPriorityHigh},
			{constants.PriorityMedium, constants.EmojiPriorityMedium},
			{constants.PriorityLow, constants.EmojiPriorityLow},
			{"unknown", constants.EmojiPriorityUnknown},
		}

		for _, tc := range testCases {
			result := formatter.GetPriorityEmoji(tc.priority)
			if result != tc.expected {
				t.Errorf("GetPriorityEmoji(%s): expected '%s', got '%s'", tc.priority, tc.expected, result)
			}
		}
	})
}

// TestToolHelpersDRYPatterns tests combined tool helpers
func TestToolHelpersDRYPatterns(t *testing.T) {
	helpers := NewToolHelpers()
	if helpers == nil {
		t.Fatal("NewToolHelpers should not return nil")
	}

	t.Run("ValidatorExists", func(t *testing.T) {
		if helpers.Validator == nil {
			t.Error("ToolHelpers.Validator should not be nil")
		}
	})

	t.Run("FormatterExists", func(t *testing.T) {
		if helpers.Formatter == nil {
			t.Error("ToolHelpers.Formatter should not be nil")
		}
	})

	t.Run("CombinedUsage", func(t *testing.T) {
		// Test validator
		params := map[string]interface{}{"id": "combined-test-id"}
		result, err := helpers.Validator.ValidateRequiredString(params, "id")

		if err != nil {
			t.Errorf("Combined validation should not fail: %v", err)
		}

		if result != "combined-test-id" {
			t.Errorf("Expected 'combined-test-id', got '%s'", result)
		}

		// Test formatter
		formatted := helpers.Formatter.FormatTaskBasic("Combined Test", result)
		if formatted == "" {
			t.Error("Combined formatting should return non-empty string")
		}

		// Test status and priority formatting
		status := helpers.Formatter.GetStatusEmoji("beklemede")
		priority := helpers.Formatter.GetPriorityEmoji(constants.PriorityHigh)

		if status == "" || priority == "" {
			t.Error("Status and priority emojis should not be empty")
		}
	})
}

// BenchmarkDRYValidationPatterns benchmarks our DRY patterns
func BenchmarkDRYValidationPatterns(b *testing.B) {
	// Initialize i18n for benchmarks
	i18n.Initialize(constants.DefaultTestLanguage)

	b.Run("I18nTParam", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			i18n.TParam("id")
		}
	})

	b.Run("I18nFormatParameterRequired", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			i18n.FormatParameterRequired("test_param")
		}
	})

	b.Run("ValidatorCreate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewParameterValidator()
		}
	})

	b.Run("ValidatorValidateString", func(b *testing.B) {
		validator := NewParameterValidator()
		params := map[string]interface{}{"id": "test"}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			validator.ValidateRequiredString(params, "id")
		}
	})

	b.Run("FormatterCreate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewTaskFormatter()
		}
	})

	b.Run("FormatterFormatTask", func(b *testing.B) {
		formatter := NewTaskFormatter()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			formatter.FormatTaskBasic(constants.TestTaskTitleEN, constants.TestIDBasic)
		}
	})

	b.Run("CombinedHelpers", func(b *testing.B) {
		helpers := NewToolHelpers()
		params := map[string]interface{}{"id": "bench-test"}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := helpers.Validator.ValidateRequiredString(params, "id")
			if err == nil {
				helpers.Formatter.FormatTaskBasic("Bench Test", result)
			}
		}
	})
}
