package constants

// Test database constants to eliminate hardcoded values in test files
const (
	// TestDatabaseURI for in-memory test database
	TestDatabaseURI = ":memory:"

	// TestMigrationsPath for test database migrations
	TestMigrationsPath = "file://../../internal/veri/migrations"

	// TestMigrationsPathShort for tests in internal/veri
	TestMigrationsPathShort = "file://migrations"

	// TestMigrationsPathLong for deeply nested test files
	TestMigrationsPathLong = "file://../../../internal/veri/migrations"

	// TestMigrationsPathIntegration for integration test files
	TestMigrationsPathIntegration = "file://../internal/veri/migrations"
)

// Test data constants for consistent test values
const (
	// TestTaskID common test task UUID
	TestTaskID = "12345678-1234-1234-1234-123456789012"

	// TestProjectID common test project UUID
	TestProjectID = "87654321-4321-4321-4321-210987654321"

	// TestUserID common test user identifier
	TestUserID = "test-user-001"

	// TestTemplateID common test template UUID
	TestTemplateID = "template-1234-5678-9012-345678901234"

	// TestTaskShortID common test task short ID (8 chars)
	TestTaskShortID = "12345678"

	// Test ID variants for different test scenarios
	TestIDBasic      = "test-id"
	TestIDConcurrent = "test-id-concurrent"
	TestIDStress     = "test-stress"
	TestIDDryPattern = "test-dry-pattern"
	TestIDWorkflow   = "test-workflow"
	TestIDValidation = "test-id-123"

	// Test project IDs for different scenarios
	TestProjectIDEdge       = "test-edge-project"
	TestProjectIDTemplate   = "test-template-edge-project"
	TestProjectIDDep        = "test-dep-project"
	TestProjectIDPerf       = "test-perf-project"
	TestProjectIDBug        = "test-bug-project"
	TestProjectIDTechDebt   = "test-tech-debt-project"
	TestProjectIDConcurrent = "test-concurrent-project"

	// Test helper IDs for generic test helpers
	TestHelperProjectID = "test-project-id"
	TestHelperTaskID    = "test-task-id"
	TestParentTaskID    = "parent-123"
	TestChildTaskID     = "child-123"
)

// Test language and localization constants
const (
	// DefaultTestLanguage for i18n testing
	DefaultTestLanguage = "tr"

	// AlternateTestLanguage for i18n testing
	AlternateTestLanguage = "en"
)

// Test iteration and limit constants
const (
	// TestIterationSmall for small test loops
	TestIterationSmall = 10

	// TestIterationMedium for medium test loops
	TestIterationMedium = 50

	// TestIterationLarge for large test loops
	TestIterationLarge = 100

	// TestIterationStress for stress testing
	TestIterationStress = 1000

	// TestTimeoutShort for quick operations (milliseconds)
	TestTimeoutShort = 100

	// TestTimeoutMedium for medium operations (milliseconds)
	TestTimeoutMedium = 500

	// TestTimeoutLong for long operations (milliseconds)
	TestTimeoutLong = 1000

	// TestConcurrencySmall for small concurrent tests
	TestConcurrencySmall = 5

	// TestConcurrencyMedium for medium concurrent tests
	TestConcurrencyMedium = 10

	// TestConcurrencyLarge for large concurrent tests
	TestConcurrencyLarge = 50
)

// Test time duration constants (multiplier values for time.Second, time.Millisecond)
const (
	// TestTimeoutShortSeconds for quick timeout tests (seconds)
	TestTimeoutShortSeconds = 5

	// TestTimeoutMediumSeconds for medium timeout tests (seconds)
	TestTimeoutMediumSeconds = 10

	// TestTimeoutLongSeconds for long timeout tests (seconds)
	TestTimeoutLongSeconds = 15

	// TestTimeoutStressSeconds for stress test timeouts (seconds)
	TestTimeoutStressSeconds = 25

	// TestTimeoutLargeSeconds for large operation timeouts (seconds)
	TestTimeoutLargeSeconds = 30

	// TestSetupTimeoutMs maximum setup time in milliseconds
	TestSetupTimeoutMs = 100

	// TestCallTimeoutMs maximum call time in milliseconds
	TestCallTimeoutMs = 10

	// TestRaceMicroseconds small delay for race condition tests (microseconds)
	TestRaceMicroseconds = 1
)

// Test success rate and threshold constants
const (
	// TestSuccessRateThreshold minimum success rate for stress tests (0.95 = 95%)
	TestSuccessRateThreshold = 0.95

	// TestRaceAllowedFailures allowed failures for race condition tests
	TestRaceAllowedFailures = 100

	// TestStressAllowedFailures allowed failures for stress tests
	TestStressAllowedFailures = 5
)

// Test string constants for common test values
const (
	// TestTaskTitle common test task title (Turkish)
	TestTaskTitle = "Test Görev"

	// TestTaskTitleEN common test task title (English)
	TestTaskTitleEN = "Test Task"

	// TestTaskDescription common test task description (Turkish)
	TestTaskDescription = "Test görev açıklaması"

	// TestTaskDescriptionEN common test task description (English)
	TestTaskDescriptionEN = "Test task description"

	// TestProjectName common test project name (Turkish)
	TestProjectName = "Test Projesi"

	// TestProjectNameEN common test project name (English)
	TestProjectNameEN = "Test Project"

	// TestProjectDescription common test project description (Turkish)
	TestProjectDescription = "Test proje açıklaması"

	// TestProjectDescriptionEN common test project description (English)
	TestProjectDescriptionEN = "Test project description"

	// TestActionName common test action name
	TestActionName = "Test Action"

	// Test task name variants for different scenarios
	TestTaskActive       = "Active Task"
	TestTaskHighPriority = "High Priority"
	TestTaskNormal       = "Normal Task"
	TestTaskOne          = "Task 1"
	TestTaskTwo          = "Task 2"

	// Test research topic
	TestResearchTopic = "testing"

	// TestTagName common test tag name
	TestTagName = "test-tag"

	// TestInvalidID invalid ID for negative testing
	TestInvalidID = "invalid-id"

	// TestEmptyString empty string for testing
	TestEmptyString = ""

	// TestPerformanceImprovement common test performance improvement text
	TestPerformanceImprovement = "50% reduction in page load time"
)

// Test date constants
const (
	// TestDateString standardized test date
	TestDateString = "2025-01-15"

	// TestDateTimeString standardized test datetime
	TestDateTimeString = "2025-01-15 10:30:00"

	// TestFutureDateString standardized future test date
	TestFutureDateString = "2025-12-31"

	// TestPastDateString standardized past test date
	TestPastDateString = "2025-06-30"
)

// Test template constants to eliminate template_id duplications
const (
	// TestTemplateFeatureRequest for feature request tests
	TestTemplateFeatureRequest = "feature_request"

	// TestTemplateBugFix for bug fix tests
	TestTemplateBugFix = "bug-fix"

	// TestTemplateSimple for simple template tests
	TestTemplateSimple = "simple-template"

	// TestTemplateNonExistent for error testing
	TestTemplateNonExistent = "non-existent-template-id"

	// TestTemplateSomeID for generic template tests
	TestTemplateSomeID = "some-id"
)

// Test context constants for frequently used test patterns
const (
	// TestIterationLimit for standard test loops
	TestIterationLimit = 50

	// TestPaginationLimit for pagination tests
	TestPaginationLimit = 10

	// TestLargeOffset for testing pagination with large offset
	TestLargeOffset = 1000

	// TestStressIterations for stress testing
	TestStressIterations = 1000

	// TestLongTitleLength for testing very long titles
	TestLongTitleLength = 200

	// Mathematical test constants for comparison functions
	TestMathSmallValue     = 5
	TestMathMediumValue    = 10
	TestMathLargeValue     = 20
	TestMathEqualValue     = 15
	TestMathEqualCompare   = 7
	TestMathNegativeSmall  = -5
	TestMathNegativeLarge  = -10
	TestMathZero           = 0
	TestMathHugeValue      = 1000000
	TestMathHugeValueMinus = 999999

	// Edge case test constants
	TestStringVeryLong    = 10000
	TestStringLong        = 1000
	TestStringMedium      = 100
	TestStringHuge        = 100000
	TestEdgeCaseLimit     = 100
	TestLargeIteration    = 60
	TestDescriptionRepeat = 50
)
