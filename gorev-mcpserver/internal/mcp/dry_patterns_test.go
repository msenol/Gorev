package mcp

import (
	"strings"
	"testing"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// TestDRYPatternsBasic tests our basic DRY patterns work correctly
func TestDRYPatternsBasic(t *testing.T) {
	// Initialize i18n
	if err := i18n.Initialize(constants.DefaultTestLanguage); err != nil {
		t.Fatalf("Failed to initialize i18n: %v", err)
	}

	// Test i18n DRY helpers
	t.Run("I18nHelpers", func(t *testing.T) {
		// Test TParam function
		result := i18n.TParam("id_field")
		if result == "" {
			t.Error("TParam should return non-empty string for 'id'")
		}

		// Test FormatParameterRequired
		required := i18n.FormatParameterRequired("test_param")
		if required == "" {
			t.Error("FormatParameterRequired should return non-empty string")
		}

		// Test FormatInvalidValue
		invalid := i18n.FormatInvalidValue("durum", "invalid", constants.GetValidTaskStatuses()[:2])
		if invalid == "" {
			t.Error("FormatInvalidValue should return non-empty string")
		}
	})

	// Test validation DRY helpers
	t.Run("ValidationHelpers", func(t *testing.T) {
		validator := NewParameterValidator()
		if validator == nil {
			t.Fatal("NewParameterValidator should not return nil")
		}

		// Test valid string validation
		params := map[string]interface{}{"id": constants.TestIDBasic}
		result, err := validator.ValidateRequiredString(params, "id")
		if err != nil {
			t.Errorf("ValidateRequiredString should not fail for valid input: %v", err)
		}
		if result != constants.TestIDBasic {
			t.Errorf("Expected 'test-id', got '%s'", result)
		}

		// Test missing required parameter
		emptyParams := map[string]interface{}{}
		_, err = validator.ValidateRequiredString(emptyParams, "id")
		if err == nil {
			t.Error("ValidateRequiredString should fail for missing required parameter")
		}
	})

	// Test formatting DRY helpers
	t.Run("FormattingHelpers", func(t *testing.T) {
		formatter := NewTaskFormatter()
		if formatter == nil {
			t.Fatal("NewTaskFormatter should not return nil")
		}

		// Test basic task formatting
		result := formatter.FormatTaskBasic(constants.TestTaskTitleEN, constants.TestTaskID)
		if result == "" {
			t.Error("FormatTaskBasic should return non-empty string")
		}

		// Test status emoji
		emoji := formatter.GetStatusEmoji("beklemede")
		if emoji == "" {
			t.Error("GetStatusEmoji should return non-empty string")
		}

		// Test priority emoji
		priorityEmoji := formatter.GetPriorityEmoji(constants.PriorityHigh)
		if priorityEmoji == "" {
			t.Error("GetPriorityEmoji should return non-empty string")
		}
	})

	// Test combined DRY helpers
	t.Run("CombinedHelpers", func(t *testing.T) {
		helpers := NewToolHelpers()
		if helpers == nil {
			t.Fatal("NewToolHelpers should not return nil")
		}

		if helpers.Validator == nil {
			t.Error("ToolHelpers.Validator should not be nil")
		}

		if helpers.Formatter == nil {
			t.Error("ToolHelpers.Formatter should not be nil")
		}

		// Test combined usage
		params := map[string]interface{}{"id": constants.TestIDBasic}
		result, err := helpers.Validator.ValidateRequiredString(params, "id")
		if err != nil {
			t.Errorf("Combined validation should not fail: %v", err)
		}

		formatted := helpers.Formatter.FormatTaskBasic("Combined Test", result)
		if formatted == "" {
			t.Error("Combined formatting should return non-empty string")
		}
	})
}

// TestTableDrivenPatterns tests our table-driven test patterns
func TestTableDrivenPatterns(t *testing.T) {
	// Simple table-driven test example
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"EmptyString", "", false},
		{"NonEmptyString", "test", true},
		{"SpaceOnly", " ", false},
		{"ValidID", constants.TestIDValidation, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Check for non-empty and non-whitespace-only content
			trimmed := strings.TrimSpace(tc.input)
			result := trimmed != ""
			if result != tc.expected {
				t.Errorf("Expected %v, got %v for input '%s'", tc.expected, result, tc.input)
			}
		})
	}
}

// TestDRYAssertionHelpers tests our assertion helper functions
func TestDRYAssertionHelpers(t *testing.T) {
	t.Run("AssertEqual", func(t *testing.T) {
		// This would normally fail the test, but we'll capture it
		// AssertEqual(t, "expected", "actual", "test message")

		// Test with matching values
		AssertEqual(t, "same", "same", "should not fail")
	})

	t.Run("AssertError", func(t *testing.T) {
		// Test expecting no error
		AssertError(t, nil, false, "should not fail for nil error")

		// Note: We can't easily test the failure case without a custom testing.T
	})

	t.Run("AssertContains", func(t *testing.T) {
		AssertContains(t, "hello world", "world", "should find substring")
		AssertContains(t, "test string", "test", "should find prefix")
	})
}

// TestConcurrencyPatterns tests our concurrency test patterns work
func TestConcurrencyPatterns(t *testing.T) {
	config := ConcurrencyTestConfig{
		Name: "BasicConcurrency",
		Setup: func() interface{} {
			return map[string]interface{}{"counter": 0}
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			// Simple operation that should always succeed
			_ = data.(map[string]interface{})
			return nil
		},
		Goroutines:             3,
		OperationsPerGoroutine: 5,
		Timeout:                5 * 1000 * 1000 * 1000, // 5 seconds in nanoseconds
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	result := RunConcurrencyTest(t, config)

	if result.TotalOperations != 15 {
		t.Errorf("Expected 15 total operations, got %d", result.TotalOperations)
	}

	if result.FailedOps > 0 {
		t.Errorf("Expected 0 failed operations, got %d", result.FailedOps)
	}
}

// BenchmarkDRYPatternsSample benchmarks our DRY patterns
func BenchmarkDRYPatternsSample(b *testing.B) {
	i18n.Initialize(constants.DefaultTestLanguage)

	cases := []BenchmarkTestCase{
		{
			Name: "I18nTParam",
			Setup: func() interface{} {
				return nil
			},
			Cleanup: func() {},
			Operation: func(data interface{}) error {
				i18n.TParam("id_field")
				return nil
			},
		},
		{
			Name: "FormatterBasic",
			Setup: func() interface{} {
				return NewTaskFormatter()
			},
			Cleanup: func() {},
			Operation: func(data interface{}) error {
				formatter := data.(*TaskFormatter)
				formatter.FormatTaskBasic("Test", constants.TestIDBasic)
				return nil
			},
		},
	}

	BenchmarkRunner(b, cases)
}
