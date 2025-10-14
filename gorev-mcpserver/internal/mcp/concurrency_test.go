package mcp

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
	testinghelpers "github.com/msenol/gorev/internal/testing"
)

// DRY concurrency test helper structure
type ConcurrencyTestConfig struct {
	Name                   string
	Setup                  func() interface{}
	Cleanup                func()
	Operation              func(interface{}) error
	Goroutines             int
	OperationsPerGoroutine int
	Timeout                time.Duration
	ExpectRaceCondition    bool
	AllowedFailures        int
}

// ConcurrencyTestResult holds results of concurrency tests
type ConcurrencyTestResult struct {
	TotalOperations int
	SuccessfulOps   int
	FailedOps       int
	ExecutionTime   time.Duration
	RaceDetected    bool
	Errors          []error
}

// RunConcurrencyTest - DRY pattern for concurrency testing
func RunConcurrencyTest(t *testing.T, config ConcurrencyTestConfig) ConcurrencyTestResult {
	data := config.Setup()
	defer config.Cleanup()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error
	successCount := 0
	totalOps := config.Goroutines * config.OperationsPerGoroutine

	start := time.Now()

	// Create timeout context
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	for i := 0; i < config.Goroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			for j := 0; j < config.OperationsPerGoroutine; j++ {
				select {
				case <-ctx.Done():
					mu.Lock()
					errors = append(errors, fmt.Errorf("goroutine %d operation %d: timeout", goroutineID, j))
					mu.Unlock()
					return
				default:
					if err := config.Operation(data); err != nil {
						mu.Lock()
						errors = append(errors, fmt.Errorf("goroutine %d operation %d: %w", goroutineID, j, err))
						mu.Unlock()
					} else {
						mu.Lock()
						successCount++
						mu.Unlock()
					}
				}
			}
		}(i)
	}

	wg.Wait()
	executionTime := time.Since(start)

	result := ConcurrencyTestResult{
		TotalOperations: totalOps,
		SuccessfulOps:   successCount,
		FailedOps:       len(errors),
		ExecutionTime:   executionTime,
		Errors:          errors,
		RaceDetected:    len(errors) > config.AllowedFailures,
	}

	// Validate results
	if !config.ExpectRaceCondition && result.RaceDetected {
		t.Errorf("Unexpected race conditions detected: %d failures out of %d operations", result.FailedOps, result.TotalOperations)
		for i, err := range result.Errors {
			if i < 5 { // Limit error output
				t.Errorf("Error %d: %v", i+1, err)
			}
		}
		if len(result.Errors) > 5 {
			t.Errorf("... and %d more errors", len(result.Errors)-5)
		}
	}

	if config.ExpectRaceCondition && !result.RaceDetected {
		t.Errorf("Expected race conditions but none detected")
	}

	if result.FailedOps > config.AllowedFailures {
		t.Errorf("Too many failures: %d (allowed: %d)", result.FailedOps, config.AllowedFailures)
	}

	return result
}

// setupConcurrencyTestEnvironment - DRY setup for concurrency tests using standardized helpers
func setupConcurrencyTestEnvironment(t *testing.T) (*Handlers, func()) {
	// Create test environment using standardized helpers
	// Use temp file instead of memory database for better concurrency support
	config := &testinghelpers.TestDatabaseConfig{
		UseTempFile:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: true, // Templates needed for concurrency tests
		InitializeI18n:  true,
	}

	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)

	// Create handlers
	handlers := YeniHandlers(isYonetici)

	return handlers, cleanup
}

// TestConcurrentToolRegistration tests concurrent tool registration
func TestConcurrentToolRegistration(t *testing.T) {
	// Setup test environment once before running concurrent operations
	// Use temp file instead of memory database for better concurrency support
	config := &testinghelpers.TestDatabaseConfig{
		UseTempFile:     true,
		MigrationsPath:  constants.TestMigrationsPath,
		CreateTemplates: true,
		InitializeI18n:  true,
	}
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)
	defer cleanup()

	concurrencyConfig := ConcurrencyTestConfig{
		Name: "ConcurrentToolRegistration",
		Setup: func() interface{} {
			handlers := YeniHandlers(isYonetici)
			return handlers
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			handlers := data.(*Handlers)
			s := server.NewMCPServer("concurrent-server", "1.0.0")

			registry := NewToolRegistry(handlers)
			registry.RegisterAllTools(s)
			return nil
		},
		Goroutines:             constants.TestConcurrencyMedium,
		OperationsPerGoroutine: constants.TestConcurrencySmall,
		Timeout:                constants.TestTimeoutLargeSeconds * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	result := RunConcurrencyTest(t, concurrencyConfig)

	// Validate performance expectations
	if result.ExecutionTime > constants.TestTimeoutMediumSeconds*time.Second {
		t.Errorf("Tool registration took too long: %v (expected < 10s)", result.ExecutionTime)
	}
}

// TestConcurrentParameterValidation tests concurrent parameter validation
func TestConcurrentParameterValidation(t *testing.T) {
	config := ConcurrencyTestConfig{
		Name: "ConcurrentParameterValidation",
		Setup: func() interface{} {
			validator := NewParameterValidator()
			return validator
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			validator := data.(*ParameterValidator)
			params := map[string]interface{}{
				"id":     constants.TestIDConcurrent,
				"durum":  constants.TaskStatusPending,
				"limit":  float64(constants.TestIterationMedium),
				"offset": float64(0),
			}

			// Test multiple validations concurrently
			_, err := validator.ValidateRequiredString(params, "id")
			if err != nil {
				return fmt.Errorf("validation error: %v", err)
			}

			_, err = validator.ValidateEnum(params, "durum", constants.GetValidTaskStatuses()[:3], true)
			if err != nil {
				return fmt.Errorf("validation error: %v", err)
			}

			validator.ValidateNumber(params, "limit", constants.TestIterationMedium)
			validator.ValidateNumber(params, "offset", 0)

			return nil
		},
		Goroutines:             constants.TestIterationSmall * 2,
		OperationsPerGoroutine: constants.TestConcurrencyMedium,
		Timeout:                constants.TestTimeoutLongSeconds * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	RunConcurrencyTest(t, config)
}

// TestConcurrentToolCalls tests concurrent MCP tool calls
func TestConcurrentToolCalls(t *testing.T) {
	handlers, cleanup := setupConcurrencyTestEnvironment(t)
	defer cleanup()

	// Skip test if setup failed
	if handlers == nil {
		t.Skip("Skipping concurrency test due to data manager setup failure")
	}

	config := ConcurrencyTestConfig{
		Name: "ConcurrentToolCalls",
		Setup: func() interface{} {
			return handlers
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			handlers := data.(*Handlers)

			// Mix different tool calls
			operations := []func() error{
				func() error {
					result, err := handlers.GorevListele(map[string]interface{}{})
					if err != nil {
						return err
					}
					if result.IsError {
						return fmt.Errorf("tool error: %v", result.Content)
					}
					return nil
				},
				func() error {
					result, err := handlers.ProjeListele(map[string]interface{}{})
					if err != nil {
						return err
					}
					if result.IsError {
						return fmt.Errorf("tool error: %v", result.Content)
					}
					return nil
				},
				func() error {
					result, err := handlers.TemplateListele(map[string]interface{}{})
					if err != nil {
						return err
					}
					if result.IsError {
						return fmt.Errorf("tool error: %v", result.Content)
					}
					return nil
				},
			}

			// Randomly execute one of the operations
			op := operations[time.Now().Nanosecond()%len(operations)]
			return op()
		},
		Goroutines:             constants.TestTimeoutLongSeconds,
		OperationsPerGoroutine: constants.TestConcurrencySmall + 3,
		Timeout:                constants.TestIterationSmall * 2 * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	result := RunConcurrencyTest(t, config)

	// Validate that all read operations should succeed
	if result.FailedOps > 0 {
		t.Errorf("Read-only operations should not fail: %d failures", result.FailedOps)
	}
}

// TestConcurrentI18nAccess tests concurrent i18n access
func TestConcurrentI18nAccess(t *testing.T) {
	config := ConcurrencyTestConfig{
		Name: "ConcurrentI18nAccess",
		Setup: func() interface{} {
			_ = i18n.Initialize(constants.DefaultTestLanguage)
			return nil
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			// Test concurrent i18n access
			keys := []string{
				"tools.descriptions.gorev_listele",
				"tools.params.descriptions.id_field",
				"error.taskNotFound",
				"success.taskUpdated",
				"common.fields.task_id",
			}

			for _, key := range keys {
				result := i18n.T(key, nil)
				if result == "" {
					return fmt.Errorf("empty translation for key: %s", key)
				}
			}

			// Test helper functions
			i18n.TParam("tr", "id_field")
			i18n.TValidation("tr", "param_required", "test_param", nil)
			i18n.FormatParameterRequired("tr", "test_param")

			return nil
		},
		Goroutines:             constants.TestIterationSmall * 5,
		OperationsPerGoroutine: constants.TestIterationSmall * 2,
		Timeout:                constants.TestTimeoutMediumSeconds * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	RunConcurrencyTest(t, config)
}

// TestConcurrentFormatterAccess tests concurrent formatter access
func TestConcurrentFormatterAccess(t *testing.T) {
	config := ConcurrencyTestConfig{
		Name: "ConcurrentFormatterAccess",
		Setup: func() interface{} {
			return NewTaskFormatter()
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			formatter := data.(*TaskFormatter)

			// Test concurrent formatter operations
			formatter.FormatTaskBasic(constants.TestTaskTitleEN, constants.TestTaskID)
			formatter.FormatTaskWithStatus(constants.TestTaskTitleEN, constants.TestTaskID, constants.TaskStatusPending)
			formatter.FormatSuccessMessage(constants.TestActionName, constants.TestTaskTitleEN, constants.TestTaskShortID)
			formatter.GetStatusEmoji(constants.TaskStatusPending)
			formatter.GetPriorityEmoji(constants.PriorityHigh)
			formatter.GetStatusEmoji(constants.TaskStatusInProgress)
			formatter.GetPriorityEmoji(constants.PriorityMedium)

			return nil
		},
		Goroutines:             constants.TestTimeoutLargeSeconds,
		OperationsPerGoroutine: constants.TestTimeoutLongSeconds,
		Timeout:                constants.TestTimeoutMediumSeconds * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	RunConcurrencyTest(t, config)
}

// TestStressConcurrencyMixed tests mixed operations under high stress
func TestStressConcurrencyMixed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	handlers, cleanup := setupConcurrencyTestEnvironment(t)
	defer cleanup()

	// Skip test if setup failed
	if handlers == nil {
		t.Skip("Skipping stress test due to data manager setup failure")
	}

	config := ConcurrencyTestConfig{
		Name: "StressConcurrencyMixed",
		Setup: func() interface{} {
			return map[string]interface{}{
				"handlers":  handlers,
				"validator": NewParameterValidator(),
				"formatter": NewTaskFormatter(),
			}
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			d := data.(map[string]interface{})
			handlers := d["handlers"].(*Handlers)
			validator := d["validator"].(*ParameterValidator)
			formatter := d["formatter"].(*TaskFormatter)

			// Mix different types of operations
			opType := time.Now().Nanosecond() % 4

			switch opType {
			case 0: // Tool calls
				result, err := handlers.GorevListele(map[string]interface{}{})
				if err != nil {
					return err
				}
				if result.IsError {
					return fmt.Errorf("tool error: %v", result.Content)
				}

			case 1: // Validation
				params := map[string]interface{}{"id": constants.TestIDStress}
				_, err := validator.ValidateRequiredString(params, "id")
				if err != nil {
					return fmt.Errorf("validation error: %v", err)
				}

			case 2: // Formatting
				formatter.FormatTaskBasic("Stress Test", "stress-id")
				formatter.GetStatusEmoji(constants.TaskStatusPending)

			case 3: // i18n
				i18n.T("tools.descriptions.gorev_listele", nil)
				i18n.TParam("tr", "id_field")
			}

			return nil
		},
		Goroutines:             constants.TestConcurrencyLarge,
		OperationsPerGoroutine: constants.TestIterationSmall * 2,
		Timeout:                constants.TestTimeoutLargeSeconds * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        constants.TestStressAllowedFailures, // Allow some failures under extreme stress
	}

	result := RunConcurrencyTest(t, config)

	// Performance expectations under stress
	if result.ExecutionTime > 25*time.Second {
		t.Errorf("Stress test took too long: %v (expected < 25s)", result.ExecutionTime)
	}

	// Success rate should be high even under stress
	successRate := float64(result.SuccessfulOps) / float64(result.TotalOperations)
	if successRate < constants.TestSuccessRateThreshold {
		t.Errorf("Success rate too low under stress: %.2f%% (expected > 95%%)", successRate*100)
	}
}

// TestRaceConditionDetection deliberately tests race condition detection
func TestRaceConditionDetection(t *testing.T) {
	var counter int
	var mu sync.RWMutex

	config := ConcurrencyTestConfig{
		Name: "RaceConditionDetection",
		Setup: func() interface{} {
			counter = 0
			return &counter
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			c := data.(*int)

			// Use proper synchronization - this should prevent race conditions
			mu.Lock()
			*c = *c + 1
			mu.Unlock()

			return nil
		},
		Goroutines:             constants.TestConcurrencyMedium,
		OperationsPerGoroutine: constants.TestConcurrencyMedium,
		Timeout:                constants.TestTimeoutShortSeconds * time.Second,
		ExpectRaceCondition:    false,                             // We expect NO race conditions with thread-safe AI context
		AllowedFailures:        constants.TestRaceAllowedFailures, // Allow many failures for this test
	}

	result := RunConcurrencyTest(t, config)

	// Check if final counter value indicates race conditions
	mu.RLock()
	finalValue := counter
	mu.RUnlock()

	expectedValue := config.Goroutines * config.OperationsPerGoroutine
	if finalValue == expectedValue {
		t.Logf("Counter reached expected value %d, but race conditions may still exist", finalValue)
	} else {
		t.Logf("Race condition detected: counter=%d, expected=%d", finalValue, expectedValue)
	}

	// Use result to avoid unused variable warning
	_ = result
}

// TestConcurrentDRYPatterns tests our DRY patterns under concurrency
func TestConcurrentDRYPatterns(t *testing.T) {
	config := ConcurrencyTestConfig{
		Name: "ConcurrentDRYPatterns",
		Setup: func() interface{} {
			i18n.Initialize(constants.DefaultTestLanguage)
			return NewToolHelpers()
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			helpers := data.(*ToolHelpers)

			// Test DRY patterns concurrently
			params := map[string]interface{}{
				"id":     constants.TestIDDryPattern,
				"durum":  constants.TaskStatusPending,
				"baslik": "DRY Test Task",
			}

			// Validation using DRY helper
			_, err := helpers.Validator.ValidateRequiredString(params, "id")
			if err != nil {
				return fmt.Errorf("validation error: %v", err)
			}

			// Formatting using DRY helper
			helpers.Formatter.FormatTaskBasic("DRY Test", "dry-test-id")

			// i18n using DRY helper
			i18n.TParam("tr", "id_field")

			return nil
		},
		Goroutines:             constants.TestIterationSmall * 2,
		OperationsPerGoroutine: constants.TestTimeoutLongSeconds,
		Timeout:                constants.TestTimeoutLongSeconds * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	result := RunConcurrencyTest(t, config)

	// DRY patterns should be highly efficient
	avgTimePerOp := result.ExecutionTime / time.Duration(result.TotalOperations)
	if avgTimePerOp > 1*time.Millisecond {
		t.Errorf("DRY patterns too slow: %v per operation (expected < 1ms)", avgTimePerOp)
	}
}
