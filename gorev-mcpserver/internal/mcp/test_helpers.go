package mcp

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
)

// DRY test helper structures and functions

// TestCase represents a generic test case structure
type TestCase struct {
	Name           string
	Input          interface{}
	Expected       interface{}
	ExpectedError  error
	ShouldFail     bool
	Setup          func() interface{}
	Cleanup        func()
	PreConditions  func() error
	PostConditions func(result interface{}) error
}

// ValidationTestCase specific for parameter validation tests
type ValidationTestCase struct {
	TestCase
	Params        map[string]interface{}
	ParamName     string
	ExpectedValue interface{}
	ValidOptions  []string
	Required      bool
}

// HandlerTestCase specific for MCP handler tests
type HandlerTestCase struct {
	TestCase
	HandlerName  string
	Params       map[string]interface{}
	ExpectedType string // "success", "error", "content"
	ContentCheck func(content interface{}) bool
}

// TableDrivenTest runs a series of test cases using DRY patterns
func TableDrivenTest(t *testing.T, testName string, cases []TestCase, testFunc func(*testing.T, TestCase)) {
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s_%s", testName, tc.Name), func(t *testing.T) {
			// Setup
			var setupData interface{}
			if tc.Setup != nil {
				setupData = tc.Setup()
			}

			// Pre-conditions
			if tc.PreConditions != nil {
				if err := tc.PreConditions(); err != nil {
					t.Fatalf("Pre-condition failed: %v", err)
				}
			}

			// Execute test
			testFunc(t, tc)

			// Post-conditions
			if tc.PostConditions != nil {
				if err := tc.PostConditions(setupData); err != nil {
					t.Errorf("Post-condition failed: %v", err)
				}
			}

			// Cleanup
			if tc.Cleanup != nil {
				tc.Cleanup()
			}
		})
	}
}

// TestEnvironment provides a reusable test environment
type TestEnvironment struct {
	Handlers    *Handlers
	DataManager *gorev.VeriYonetici
	Cleanup     func()
}

// SetupTestEnvironment creates a standard test environment
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	// Initialize i18n for consistent testing
	if err := i18n.Initialize("tr"); err != nil {
		t.Fatalf("Failed to initialize i18n: %v", err)
	}

	// Create in-memory database
	dataManager, err := gorev.YeniVeriYonetici(constants.TestDatabaseURI, constants.TestMigrationsPath)
	if err != nil {
		t.Fatalf("Failed to create test data manager: %v", err)
	}

	// Create business logic manager
	isYonetici := gorev.YeniIsYonetici(dataManager)

	// Create handlers
	handlers := YeniHandlers(isYonetici)

	cleanup := func() {
		// Cleanup resources if needed
	}

	return &TestEnvironment{
		Handlers:    handlers,
		DataManager: dataManager,
		Cleanup:     cleanup,
	}
}

// AssertEqual provides DRY assertion helper
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertError provides DRY error assertion helper
func AssertError(t *testing.T, err error, shouldError bool, message string) {
	if shouldError && err == nil {
		t.Errorf("%s: expected error but got none", message)
	}
	if !shouldError && err != nil {
		t.Errorf("%s: unexpected error: %v", message, err)
	}
}

// AssertContains checks if a string contains substring
func AssertContains(t *testing.T, haystack, needle string, message string) {
	if haystack == "" && needle != "" {
		t.Errorf("%s: haystack is empty but needle is not", message)
	}
	// Simple substring check - for production use a proper contains function
	found := false
	if len(needle) <= len(haystack) {
		for i := 0; i <= len(haystack)-len(needle); i++ {
			if haystack[i:i+len(needle)] == needle {
				found = true
				break
			}
		}
	}
	if !found {
		t.Errorf("%s: '%s' not found in '%s'", message, needle, haystack)
	}
}

// ValidationTestRunner runs validation tests using DRY patterns
func ValidationTestRunner(t *testing.T, validator *ParameterValidator, cases []ValidationTestCase) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setup
			if tc.Setup != nil {
				tc.Setup()
			}

			var result interface{}
			var err error

			// Run different validation types based on test case
			if tc.Required {
				var validationError *mcp.CallToolResult
				result, validationError = validator.ValidateRequiredString(tc.Params, tc.ParamName)
				if validationError != nil {
					err = fmt.Errorf("validation failed: %v", validationError.Content)
				}
			} else if len(tc.ValidOptions) > 0 {
				var validationError *mcp.CallToolResult
				result, validationError = validator.ValidateEnum(tc.Params, tc.ParamName, tc.ValidOptions, tc.Required)
				if validationError != nil {
					err = fmt.Errorf("validation failed: %v", validationError.Content)
				}
			} else {
				if expectedInt, ok := tc.ExpectedValue.(int); ok {
					result = validator.ValidateNumber(tc.Params, tc.ParamName, expectedInt)
				}
			}

			// Assertions
			AssertError(t, err, tc.ShouldFail, fmt.Sprintf("Validation test %s", tc.Name))

			if !tc.ShouldFail && tc.ExpectedValue != nil {
				AssertEqual(t, tc.ExpectedValue, result, fmt.Sprintf("Validation result %s", tc.Name))
			}

			// Cleanup
			if tc.Cleanup != nil {
				tc.Cleanup()
			}
		})
	}
}

// HandlerTestRunner runs handler tests using DRY patterns
func HandlerTestRunner(t *testing.T, env *TestEnvironment, cases []HandlerTestCase) {
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setup
			if tc.Setup != nil {
				tc.Setup()
			}

			// Execute handler based on name
			var result interface{}
			var err error

			switch tc.HandlerName {
			case "gorev_listele":
				result, err = env.Handlers.GorevListele(tc.Params)
			case "proje_listele":
				result, err = env.Handlers.ProjeListele(tc.Params)
			case "template_listele":
				result, err = env.Handlers.TemplateListele(tc.Params)
			case "ozet_goster":
				result, err = env.Handlers.OzetGoster(tc.Params)
			default:
				t.Fatalf("Unknown handler: %s", tc.HandlerName)
			}

			// Assertions
			AssertError(t, err, tc.ShouldFail, fmt.Sprintf("Handler test %s", tc.Name))

			if !tc.ShouldFail && result != nil {
				// Type-specific checks
				if tc.ContentCheck != nil {
					if !tc.ContentCheck(result) {
						t.Errorf("Content check failed for test %s", tc.Name)
					}
				}
			}

			// Cleanup
			if tc.Cleanup != nil {
				tc.Cleanup()
			}
		})
	}
}

// CreateTestProject creates a test project for use in tests
func CreateTestProject(t *testing.T, env *TestEnvironment, name, description string) string {
	params := map[string]interface{}{
		"isim":  name,
		"tanim": description,
	}

	result, err := env.Handlers.ProjeOlustur(params)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}

	if result.IsError {
		t.Fatalf("Failed to create test project: %v", result.Content)
	}

	// Extract project ID from result - simplified extraction
	// result.Content is []interface{} with text content
	if len(result.Content) > 0 {
		return constants.TestHelperProjectID
	}

	return ""
}

// CreateTestTask creates a test task using template for use in tests
func CreateTestTask(t *testing.T, env *TestEnvironment, templateID string, values map[string]interface{}) string {
	params := map[string]interface{}{
		constants.ParamTemplateID: templateID,
		constants.ParamDegerler:   values,
	}

	result, err := env.Handlers.TemplatedenGorevOlustur(params)
	if err != nil {
		t.Fatalf("Failed to create test task: %v", err)
	}

	if result.IsError {
		t.Fatalf("Failed to create test task: %v", result.Content)
	}

	// Extract task ID from result - simplified extraction
	// result.Content is []interface{} with text content
	if len(result.Content) > 0 {
		return constants.TestHelperTaskID
	}

	return ""
}

// BenchmarkTestCase for benchmark test cases
type BenchmarkTestCase struct {
	Name      string
	Setup     func() interface{}
	Cleanup   func()
	Operation func(interface{}) error
	N         int // Number of iterations for this specific case
}

// BenchmarkRunner runs benchmark tests using DRY patterns
func BenchmarkRunner(b *testing.B, cases []BenchmarkTestCase) {
	for _, bc := range cases {
		b.Run(bc.Name, func(b *testing.B) {
			data := bc.Setup()
			defer bc.Cleanup()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := bc.Operation(data); err != nil {
					b.Fatalf("Benchmark operation failed: %v", err)
				}
			}
		})
	}
}

// ParameterTestCases provides common parameter validation test cases
func ParameterTestCases() []ValidationTestCase {
	return []ValidationTestCase{
		{
			TestCase: TestCase{
				Name:       "ValidID",
				ShouldFail: false,
			},
			Params:        map[string]interface{}{"id": "valid-id-123"},
			ParamName:     "id",
			ExpectedValue: "valid-id-123",
			Required:      true,
		},
		{
			TestCase: TestCase{
				Name:       "MissingRequiredID",
				ShouldFail: true,
			},
			Params:    map[string]interface{}{},
			ParamName: "id",
			Required:  true,
		},
		{
			TestCase: TestCase{
				Name:       "ValidStatus",
				ShouldFail: false,
			},
			Params:        map[string]interface{}{"durum": constants.TaskStatusPending},
			ParamName:     "durum",
			ExpectedValue: constants.TaskStatusPending,
			ValidOptions:  constants.GetValidTaskStatuses()[:3], // Exclude "iptal" for this test
			Required:      false,
		},
		{
			TestCase: TestCase{
				Name:       "InvalidStatus",
				ShouldFail: true,
			},
			Params:       map[string]interface{}{"durum": "invalid-status"},
			ParamName:    "durum",
			ValidOptions: constants.GetValidTaskStatuses()[:3], // Exclude "iptal" for this test
			Required:     false,
		},
	}
}

// HandlerTestCases provides common handler test cases
func HandlerTestCases() []HandlerTestCase {
	return []HandlerTestCase{
		{
			TestCase: TestCase{
				Name:       "ListTasksEmpty",
				ShouldFail: false,
			},
			HandlerName:  "gorev_listele",
			Params:       map[string]interface{}{},
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
				Name:       "ListProjects",
				ShouldFail: false,
			},
			HandlerName:  "proje_listele",
			Params:       map[string]interface{}{},
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
				Name:       "ListTemplates",
				ShouldFail: false,
			},
			HandlerName:  "template_listele",
			Params:       map[string]interface{}{},
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
				Name:       "ShowSummary",
				ShouldFail: false,
			},
			HandlerName:  "ozet_goster",
			Params:       map[string]interface{}{},
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
}

// I18nTestHelper helps test i18n functionality
type I18nTestHelper struct {
	OriginalLang string
}

// SetupI18nTest sets up i18n for testing and returns cleanup function
func SetupI18nTest(lang string) (*I18nTestHelper, func()) {
	helper := &I18nTestHelper{
		OriginalLang: i18n.GetCurrentLanguage(),
	}

	i18n.Initialize(lang)

	cleanup := func() {
		// Restore original language if needed
		if helper.OriginalLang != "" {
			_ = i18n.SetLanguage(helper.OriginalLang)
		}
	}

	return helper, cleanup
}

// AssertTranslation checks if a translation key returns expected value
func (h *I18nTestHelper) AssertTranslation(t *testing.T, key, expected string, data map[string]interface{}) {
	result := i18n.T(key, data)
	if result != expected {
		t.Errorf("Translation mismatch for key '%s': expected '%s', got '%s'", key, expected, result)
	}
}

// AssertTranslationExists checks if a translation key exists (not returning the key itself)
func (h *I18nTestHelper) AssertTranslationExists(t *testing.T, key string) {
	result := i18n.T(key, nil)
	if result == key {
		t.Errorf("Translation missing for key '%s'", key)
	}
}
