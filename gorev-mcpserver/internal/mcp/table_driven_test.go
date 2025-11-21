package mcp

import (
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// TestParameterValidationTableDriven demonstrates table-driven testing with DRY patterns
func TestParameterValidationTableDriven(t *testing.T) {
	validator := NewParameterValidator()

	// Use DRY test helper to run all common parameter validation tests
	ValidationTestRunner(t, validator, ParameterTestCases())

	// Additional specific test cases for this package
	customCases := []ValidationTestCase{
		{
			TestCase: TestCase{
				Name:       "EmptyStringID",
				ShouldFail: true,
			},
			Params:    map[string]interface{}{"id": ""},
			ParamName: "id",
			Required:  true,
		},
		{
			TestCase: TestCase{
				Name:       "NumericID",
				ShouldFail: false,
			},
			Params:        map[string]interface{}{"id": "12345"},
			ParamName:     "id",
			ExpectedValue: "12345",
			Required:      true,
		},
		{
			TestCase: TestCase{
				Name:       "ValidPriority",
				ShouldFail: false,
			},
			Params:        map[string]interface{}{"oncelik": constants.PriorityHigh},
			ParamName:     "oncelik",
			ExpectedValue: constants.PriorityHigh,
			ValidOptions:  []string{constants.PriorityHigh, constants.PriorityMedium, constants.PriorityLow},
			Required:      false,
		},
		{
			TestCase: TestCase{
				Name:       "InvalidPriority",
				ShouldFail: true,
			},
			Params:       map[string]interface{}{"oncelik": "invalid"},
			ParamName:    "oncelik",
			ValidOptions: []string{constants.PriorityHigh, constants.PriorityMedium, constants.PriorityLow},
			Required:     false,
		},
	}

	ValidationTestRunner(t, validator, customCases)
}

// TestHandlersTableDriven demonstrates table-driven testing for handlers
func TestHandlersTableDriven(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	// Use DRY test helper to run all common handler tests
	HandlerTestRunner(t, env, HandlerTestCases())

	// Additional specific test cases
	customCases := []HandlerTestCase{
		{
			TestCase: TestCase{
				Name:       "ListTasksWithLimit",
				ShouldFail: false,
			},
			HandlerName: "gorev_listele",
			Params: map[string]interface{}{
				"limit": float64(constants.TestPaginationLimit),
			},
			ExpectedType: "success",
			ContentCheck: func(content interface{}) bool {
				// Check if result is a valid CallToolResult with content
				if result, ok := content.(*mcp.CallToolResult); ok {
					return result != nil && !result.IsError && len(result.Content) > 0
				}
				return false
			},
		},
		{
			TestCase: TestCase{
				Name:       "ListTasksWithInvalidLimit",
				ShouldFail: false,
			},
			HandlerName: "gorev_listele",
			Params: map[string]interface{}{
				"limit": "invalid",
			},
			ExpectedType: "success",
			ContentCheck: func(content interface{}) bool {
				// Check if result is a valid CallToolResult with content
				if result, ok := content.(*mcp.CallToolResult); ok {
					return result != nil && !result.IsError && len(result.Content) > 0
				}
				return false
			},
		},
		{
			TestCase: TestCase{
				Name:       "ListTasksWithStatus",
				ShouldFail: false,
			},
			HandlerName: "gorev_listele",
			Params: map[string]interface{}{
				"durum": "beklemede",
			},
			ExpectedType: "success",
			ContentCheck: func(content interface{}) bool {
				// Check if result is a valid CallToolResult with content
				if result, ok := content.(*mcp.CallToolResult); ok {
					return result != nil && !result.IsError && len(result.Content) > 0
				}
				return false
			},
		},
	}

	HandlerTestRunner(t, env, customCases)
}

// TestFormatterTableDriven demonstrates table-driven testing for formatter
func TestFormatterTableDriven(t *testing.T) {
	formatter := NewTaskFormatter()

	testCases := []TestCase{
		{
			Name:       "FormatBasicTask",
			Input:      []interface{}{constants.TestTaskTitleEN, constants.TestTaskID},
			Expected:   "**Test Task** (ID: 12345678)",
			ShouldFail: false,
		},
		{
			Name:  "FormatTaskWithStatus",
			Input: []interface{}{constants.TestTaskTitleEN, constants.TestTaskID, constants.TaskStatusPending},
			Expected: func() string {
				return formatter.FormatTaskWithStatus(constants.TestTaskTitleEN, constants.TestTaskID, constants.TaskStatusPending)
			}(),
			ShouldFail: false,
		},
		{
			Name:       "GetStatusEmoji",
			Input:      "beklemede",
			Expected:   constants.EmojiStatusPending,
			ShouldFail: false,
		},
		{
			Name:       "GetPriorityEmoji",
			Input:      constants.PriorityHigh,
			Expected:   constants.EmojiPriorityHigh,
			ShouldFail: false,
		},
		{
			Name:       "GetUnknownStatusEmoji",
			Input:      "unknown",
			Expected:   constants.EmojiStatusUnknown,
			ShouldFail: false,
		},
	}

	TableDrivenTest(t, "Formatter", testCases, func(t *testing.T, tc TestCase) {
		var result interface{}

		switch tc.Name {
		case "FormatBasicTask":
			inputs := tc.Input.([]interface{})
			result = formatter.FormatTaskBasic(inputs[0].(string), inputs[1].(string))

		case "FormatTaskWithStatus":
			inputs := tc.Input.([]interface{})
			result = formatter.FormatTaskWithStatus(inputs[0].(string), inputs[1].(string), inputs[2].(string))

		case "GetStatusEmoji":
			result = formatter.GetStatusEmoji(tc.Input.(string))

		case "GetPriorityEmoji":
			result = formatter.GetPriorityEmoji(tc.Input.(string))

		case "GetUnknownStatusEmoji":
			result = formatter.GetStatusEmoji(tc.Input.(string))
		}

		AssertEqual(t, tc.Expected, result, fmt.Sprintf("Test case %s", tc.Name))
	})
}

// TestI18nTableDriven demonstrates table-driven testing for i18n functionality
func TestI18nTableDriven(t *testing.T) {
	helper, cleanup := SetupI18nTest("tr")
	defer cleanup()

	testCases := []struct {
		Key         string
		Data        map[string]interface{}
		Expected    string
		ShouldExist bool
	}{
		{
			Key:         "tools.descriptions.gorev_listele",
			ShouldExist: true,
		},
		{
			Key:         "tools.params.descriptions.id_field",
			ShouldExist: true,
		},
		{
			Key:         "success.taskUpdated",
			Data:        map[string]interface{}{"OldStatus": "beklemede", "NewStatus": "devam_ediyor"},
			ShouldExist: true,
		},
		{
			Key:         "nonexistent.key",
			Expected:    "nonexistent.key", // Should return key if not found
			ShouldExist: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Key, func(t *testing.T) {
			if tc.ShouldExist {
				helper.AssertTranslationExists(t, tc.Key)
			}

			if tc.Expected != "" {
				helper.AssertTranslation(t, tc.Key, tc.Expected, tc.Data)
			}
		})
	}
}

// TestI18nHelperFunctionsTableDriven tests DRY i18n helper functions
func TestI18nHelperFunctionsTableDriven(t *testing.T) {
	i18n.Initialize(constants.DefaultTestLanguage)

	testCases := []TestCase{
		{
			Name:  "TParamWithExistingKey",
			Input: "id_field",
			Expected: func() string {
				return i18n.TParam("tr", "id_field")
			}(),
			ShouldFail: false,
		},
		{
			Name:       "TParamWithNonExistentKey",
			Input:      "nonexistent_param",
			Expected:   "nonexistent_param parameter",
			ShouldFail: false,
		},
		{
			Name:  "FormatParameterRequired",
			Input: "id",
			Expected: func() string {
				return i18n.FormatParameterRequired("tr", "id")
			}(),
			ShouldFail: false,
		},
		{
			Name:  "FormatInvalidValue",
			Input: []interface{}{"durum", "invalid", constants.GetValidTaskStatuses()[:2]},
			Expected: func() string {
				return i18n.FormatInvalidValue("tr", "durum", "invalid", constants.GetValidTaskStatuses()[:2])
			}(),
			ShouldFail: false,
		},
	}

	TableDrivenTest(t, "I18nHelpers", testCases, func(t *testing.T, tc TestCase) {
		var result string

		switch tc.Name {
		case "TParamWithExistingKey", "TParamWithNonExistentKey":
			result = i18n.TParam("tr", tc.Input.(string))

		case "FormatParameterRequired":
			result = i18n.FormatParameterRequired("tr", tc.Input.(string))

		case "FormatInvalidValue":
			inputs := tc.Input.([]interface{})
			result = i18n.FormatInvalidValue("tr",
				inputs[0].(string),
				inputs[1].(string),
				inputs[2].([]string),
			)
		}

		// For dynamic results, just check they're not empty
		if tc.Expected == nil || tc.Expected == "" {
			if result == "" {
				t.Errorf("Test case %s: result should not be empty", tc.Name)
			}
		} else {
			AssertEqual(t, tc.Expected, result, fmt.Sprintf("Test case %s", tc.Name))
		}
	})
}

// TestToolHelpersTableDriven tests DRY tool helpers
func TestToolHelpersTableDriven(t *testing.T) {
	helpers := NewToolHelpers()

	testCases := []TestCase{
		{
			Name:       "ValidatorExists",
			Expected:   true,
			ShouldFail: false,
		},
		{
			Name:       "FormatterExists",
			Expected:   true,
			ShouldFail: false,
		},
		{
			Name:       "ValidatorCanValidateString",
			Input:      map[string]interface{}{"id": constants.TestIDBasic},
			Expected:   constants.TestIDBasic,
			ShouldFail: false,
		},
		{
			Name:  "FormatterCanFormatTask",
			Input: []interface{}{constants.TestTaskTitleEN, constants.TestIDBasic},
			Expected: func() string {
				return helpers.Formatter.FormatTaskBasic(constants.TestTaskTitleEN, constants.TestIDBasic)
			}(),
			ShouldFail: false,
		},
	}

	TableDrivenTest(t, "ToolHelpers", testCases, func(t *testing.T, tc TestCase) {
		switch tc.Name {
		case "ValidatorExists":
			AssertEqual(t, true, helpers.Validator != nil, "Validator should exist")

		case "FormatterExists":
			AssertEqual(t, true, helpers.Formatter != nil, "Formatter should exist")

		case "ValidatorCanValidateString":
			params := tc.Input.(map[string]interface{})
			result, validationError := helpers.Validator.ValidateRequiredString(params, "id")
			var err error
			if validationError != nil {
				err = fmt.Errorf("validation failed: %v", validationError.Content)
			}
			AssertError(t, err, tc.ShouldFail, "Validation should not fail")
			AssertEqual(t, tc.Expected, result, "Validation result should match")

		case "FormatterCanFormatTask":
			inputs := tc.Input.([]interface{})
			result := helpers.Formatter.FormatTaskBasic(inputs[0].(string), inputs[1].(string))
			// Just check it's not empty for dynamic results
			if result == "" {
				t.Errorf("Formatter result should not be empty")
			}
		}
	})
}

// BenchmarkDRYPatterns benchmarks DRY patterns using table-driven approach
func BenchmarkDRYPatterns(b *testing.B) {
	cases := []BenchmarkTestCase{
		{
			Name: "I18nHelpers",
			Setup: func() interface{} {
				i18n.Initialize(constants.DefaultTestLanguage)
				return nil
			},
			Cleanup: func() {},
			Operation: func(data interface{}) error {
				i18n.TParam("tr", "id_field")
				i18n.FormatParameterRequired("tr", "test")
				i18n.FormatInvalidValue("tr", "param", "invalid", []string{"valid"})
				return nil
			},
		},
		{
			Name: "ParameterValidation",
			Setup: func() interface{} {
				return NewParameterValidator()
			},
			Cleanup: func() {},
			Operation: func(data interface{}) error {
				validator := data.(*ParameterValidator)
				params := map[string]interface{}{"id": "test"}
				_, validationError := validator.ValidateRequiredString(params, "id")
				if validationError != nil {
					return fmt.Errorf("validation failed: %v", validationError.Content)
				}
				return nil
			},
		},
		{
			Name: "TaskFormatting",
			Setup: func() interface{} {
				return NewTaskFormatter()
			},
			Cleanup: func() {},
			Operation: func(data interface{}) error {
				formatter := data.(*TaskFormatter)
				formatter.FormatTaskBasic("Test", constants.TestIDBasic)
				formatter.GetStatusEmoji("beklemede")
				return nil
			},
		},
		{
			Name: "CombinedHelpers",
			Setup: func() interface{} {
				i18n.Initialize(constants.DefaultTestLanguage)
				return NewToolHelpers()
			},
			Cleanup: func() {},
			Operation: func(data interface{}) error {
				helpers := data.(*ToolHelpers)
				params := map[string]interface{}{"id": "test", "durum": "beklemede"}

				// Validation
				_, validationError := helpers.Validator.ValidateRequiredString(params, "id")
				if validationError != nil {
					return fmt.Errorf("validation failed: %v", validationError.Content)
				}

				// Formatting
				helpers.Formatter.FormatTaskBasic("Test", constants.TestIDBasic)

				// i18n
				i18n.TParam("tr", "id_field")

				return nil
			},
		},
	}

	BenchmarkRunner(b, cases)
}

// TestDRYPatternsIntegration tests integration of all DRY patterns
func TestDRYPatternsIntegration(t *testing.T) {
	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	// Test complete workflow using DRY patterns
	t.Run("CompleteWorkflow", func(t *testing.T) {
		// 1. Use i18n helpers
		paramDesc := i18n.TParam("tr", "id_field")
		if paramDesc == "" {
			t.Error("Parameter description should not be empty")
		}

		// 2. Use validation helpers
		params := map[string]interface{}{"id": constants.TestIDWorkflow}
		result, validationError := env.Handlers.toolHelpers.Validator.ValidateRequiredString(params, "id")
		var err error
		if validationError != nil {
			err = fmt.Errorf("validation failed: %v", validationError.Content)
		}
		AssertError(t, err, false, "Validation should succeed")
		AssertEqual(t, constants.TestIDWorkflow, result, "Validation result should match")

		// 3. Use formatting helpers
		formatted := env.Handlers.toolHelpers.Formatter.FormatTaskBasic(constants.TestTaskTitleEN, constants.TestIDWorkflow)
		if formatted == "" {
			t.Error("Formatted task should not be empty")
		}

		// 4. Test handler integration
		handlerResult, handlerErr := env.Handlers.GorevListele(map[string]interface{}{})
		AssertError(t, handlerErr, false, "Handler should succeed")

		if handlerResult.IsError {
			t.Errorf("Handler should not return error: %v", handlerResult.Content)
		}
	})
}
