package mcp

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
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

// setupConcurrencyTestEnvironment - DRY setup for concurrency tests
func setupConcurrencyTestEnvironment() (*Handlers, func()) {
	// Initialize i18n for consistent environment
	i18n.Initialize("tr")

	// Setup data manager with thread-safe in-memory DB
	dataManager, err := gorev.YeniVeriYonetici(":memory:", "")
	if err != nil {
		panic("Failed to create data manager: " + err.Error())
	}

	// Create business logic manager
	isYonetici := gorev.YeniIsYonetici(dataManager)

	// Create handlers
	handlers := YeniHandlers(isYonetici)

	cleanup := func() {
		// Cleanup resources
	}

	return handlers, cleanup
}

// TestConcurrentToolRegistration tests concurrent tool registration
func TestConcurrentToolRegistration(t *testing.T) {
	config := ConcurrencyTestConfig{
		Name: "ConcurrentToolRegistration",
		Setup: func() interface{} {
			i18n.Initialize("tr")
			dataManager, _ := gorev.YeniVeriYonetici(":memory:", "")
			isYonetici := gorev.YeniIsYonetici(dataManager)
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
		Goroutines:             10,
		OperationsPerGoroutine: 5,
		Timeout:                30 * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	result := RunConcurrencyTest(t, config)

	// Validate performance expectations
	if result.ExecutionTime > 10*time.Second {
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
				"id":     "test-id-concurrent",
				"durum":  "beklemede",
				"limit":  float64(50),
				"offset": float64(0),
			}

			// Test multiple validations concurrently
			_, err := validator.ValidateRequiredString(params, "id")
			if err != nil {
				return fmt.Errorf("validation error: %v", err)
			}

			_, err = validator.ValidateEnum(params, "durum", []string{"beklemede", "devam_ediyor", "tamamlandi"}, true)
			if err != nil {
				return fmt.Errorf("validation error: %v", err)
			}

			validator.ValidateNumber(params, "limit", 50)
			validator.ValidateNumber(params, "offset", 0)

			return nil
		},
		Goroutines:             20,
		OperationsPerGoroutine: 10,
		Timeout:                15 * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        0,
	}

	RunConcurrencyTest(t, config)
}

// TestConcurrentToolCalls tests concurrent MCP tool calls
func TestConcurrentToolCalls(t *testing.T) {
	handlers, cleanup := setupConcurrencyTestEnvironment()
	defer cleanup()

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
		Goroutines:             15,
		OperationsPerGoroutine: 8,
		Timeout:                20 * time.Second,
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
			i18n.Initialize("tr")
			return nil
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			// Test concurrent i18n access
			keys := []string{
				"tools.descriptions.gorev_listele",
				"tools.params.descriptions.id",
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
			i18n.TParam("id")
			i18n.TValidation("param_required", "test_param", nil)
			i18n.FormatParameterRequired("test_param")

			return nil
		},
		Goroutines:             25,
		OperationsPerGoroutine: 20,
		Timeout:                10 * time.Second,
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
			formatter.FormatTaskBasic("Test Task", "12345678-1234-1234-1234-123456789012")
			formatter.FormatTaskWithStatus("Test Task", "12345678-1234-1234-1234-123456789012", "beklemede")
			formatter.FormatSuccessMessage("Test Action", "Test Task", "12345678")
			formatter.GetStatusEmoji("beklemede")
			formatter.GetPriorityEmoji("yuksek")
			formatter.GetStatusEmoji("devam_ediyor")
			formatter.GetPriorityEmoji("orta")

			return nil
		},
		Goroutines:             30,
		OperationsPerGoroutine: 15,
		Timeout:                10 * time.Second,
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

	handlers, cleanup := setupConcurrencyTestEnvironment()
	defer cleanup()

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
				params := map[string]interface{}{"id": "test-stress"}
				_, err := validator.ValidateRequiredString(params, "id")
				if err != nil {
					return fmt.Errorf("validation error: %v", err)
				}

			case 2: // Formatting
				formatter.FormatTaskBasic("Stress Test", "stress-id")
				formatter.GetStatusEmoji("beklemede")

			case 3: // i18n
				i18n.T("tools.descriptions.gorev_listele", nil)
				i18n.TParam("id")
			}

			return nil
		},
		Goroutines:             50,
		OperationsPerGoroutine: 20,
		Timeout:                30 * time.Second,
		ExpectRaceCondition:    false,
		AllowedFailures:        5, // Allow some failures under extreme stress
	}

	result := RunConcurrencyTest(t, config)

	// Performance expectations under stress
	if result.ExecutionTime > 25*time.Second {
		t.Errorf("Stress test took too long: %v (expected < 25s)", result.ExecutionTime)
	}

	// Success rate should be high even under stress
	successRate := float64(result.SuccessfulOps) / float64(result.TotalOperations)
	if successRate < 0.95 {
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

			// Deliberately create race condition by not using mutex
			// This should be detected by race detector if enabled
			current := *c
			time.Sleep(1 * time.Microsecond) // Small delay to increase race chance
			*c = current + 1

			return nil
		},
		Goroutines:             10,
		OperationsPerGoroutine: 10,
		Timeout:                5 * time.Second,
		ExpectRaceCondition:    true, // We expect this to have race conditions
		AllowedFailures:        100,  // Allow many failures for this test
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
			i18n.Initialize("tr")
			return NewToolHelpers()
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			helpers := data.(*ToolHelpers)

			// Test DRY patterns concurrently
			params := map[string]interface{}{
				"id":     "test-dry-pattern",
				"durum":  "beklemede",
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
			i18n.TParam("id")

			return nil
		},
		Goroutines:             20,
		OperationsPerGoroutine: 15,
		Timeout:                15 * time.Second,
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
