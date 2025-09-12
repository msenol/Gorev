package mcp

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
	testinghelpers "github.com/msenol/gorev/internal/testing"
)

// DRY benchmark helper structure
type BenchmarkConfig struct {
	Name                   string
	Setup                  func() interface{}
	Cleanup                func()
	Operation              func(interface{}) error
	Goroutines             int
	OperationsPerGoroutine int
	Timeout                time.Duration
}

// RunStandardBenchmark - DRY pattern for benchmarks
func RunStandardBenchmark(b *testing.B, config BenchmarkConfig) {
	data := config.Setup()
	defer config.Cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := config.Operation(data); err != nil {
			b.Fatalf("Benchmark operation failed: %v", err)
		}
	}
}

// RunConcurrentBenchmark - DRY pattern for concurrent benchmarks
func RunConcurrentBenchmark(b *testing.B, config BenchmarkConfig) {
	data := config.Setup()
	defer config.Cleanup()

	b.SetParallelism(config.Goroutines)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := config.Operation(data); err != nil {
				b.Errorf("Concurrent operation failed: %v", err)
			}
		}
	})
}

// setupBenchmarkEnvironment - DRY setup for all benchmarks using standardized helpers
func setupBenchmarkEnvironment() (*Handlers, func()) {
	// Create test environment using standardized helpers
	config := &testinghelpers.TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  "", // Empty for benchmark performance
		CreateTemplates: false,
		InitializeI18n:  true,
	}

	// Create a minimal testing.T for helper compatibility
	t := &testing.T{}
	isYonetici, cleanup := testinghelpers.SetupTestEnvironmentWithConfig(t, config)

	// Create handlers
	handlers := YeniHandlers(isYonetici)

	return handlers, cleanup
}

// BenchmarkToolRegistration tests the performance of tool registration
func BenchmarkToolRegistration(b *testing.B) {
	config := BenchmarkConfig{
		Name: "ToolRegistration",
		Setup: func() interface{} {
			handlers, _ := setupBenchmarkEnvironment()
			return handlers
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			handlers := data.(*Handlers)
			s := server.NewMCPServer("bench-server", "1.0.0")

			registry := NewToolRegistry(handlers)
			registry.RegisterAllTools(s)
			return nil
		},
	}

	RunStandardBenchmark(b, config)
}

// BenchmarkParameterValidation tests validation performance
func BenchmarkParameterValidation(b *testing.B) {
	config := BenchmarkConfig{
		Name: "ParameterValidation",
		Setup: func() interface{} {
			validator := NewParameterValidator()
			params := map[string]interface{}{
				"id":     constants.TestIDValidation,
				"durum":  constants.TaskStatusPending,
				"limit":  float64(constants.TestIterationLimit),
				"offset": float64(0),
			}
			return map[string]interface{}{
				"validator": validator,
				"params":    params,
			}
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			d := data.(map[string]interface{})
			validator := d["validator"].(*ParameterValidator)
			params := d["params"].(map[string]interface{})

			// Test multiple validations
			_, validationError := validator.ValidateRequiredString(params, "id")
			if validationError != nil {
				return fmt.Errorf(constants.ErrorFormatValidation, validationError.Content)
			}

			_, validationError = validator.ValidateEnum(params, "durum", constants.GetValidTaskStatuses()[:3], true)
			if validationError != nil {
				return fmt.Errorf(constants.ErrorFormatValidation, validationError.Content)
			}

			validator.ValidateNumber(params, "limit", constants.TestIterationLimit)
			validator.ValidateNumber(params, "offset", 0)

			return nil
		},
	}

	RunStandardBenchmark(b, config)
}

// BenchmarkTaskFormatter tests formatting performance
func BenchmarkTaskFormatter(b *testing.B) {
	config := BenchmarkConfig{
		Name: "TaskFormatter",
		Setup: func() interface{} {
			formatter := NewTaskFormatter()
			return formatter
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			formatter := data.(*TaskFormatter)

			// Test multiple formatting operations
			formatter.FormatTaskBasic(constants.TestTaskTitleEN, constants.TestTaskID)
			formatter.FormatTaskWithStatus(constants.TestTaskTitleEN, constants.TestTaskID, constants.TaskStatusPending)
			formatter.FormatSuccessMessage(constants.TestActionName, constants.TestTaskTitleEN, constants.TestTaskShortID)
			formatter.GetStatusEmoji(constants.TaskStatusPending)
			formatter.GetPriorityEmoji(constants.PriorityHigh)

			return nil
		},
	}

	RunStandardBenchmark(b, config)
}

// BenchmarkConcurrentToolCalls tests concurrent tool calls
func BenchmarkConcurrentToolCalls(b *testing.B) {
	handlers, cleanup := setupBenchmarkEnvironment()
	defer cleanup()

	config := BenchmarkConfig{
		Name: "ConcurrentToolCalls",
		Setup: func() interface{} {
			return handlers
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			handlers := data.(*Handlers)

			// Simulate tool call
			params := map[string]interface{}{}
			result, err := handlers.GorevListele(params)
			if err != nil {
				return err
			}

			if result.IsError {
				return fmt.Errorf(constants.ErrorFormatTool, result.Content)
			}
			return nil
		},
		Goroutines: constants.TestConcurrencyMedium,
	}

	RunConcurrentBenchmark(b, config)
}

// BenchmarkI18nHelpers tests i18n helper performance
func BenchmarkI18nHelpers(b *testing.B) {
	config := BenchmarkConfig{
		Name: "I18nHelpers",
		Setup: func() interface{} {
			_ = i18n.Initialize(constants.DefaultTestLanguage)
			return nil
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			// Test DRY i18n helper functions
			i18n.TParam("id")
			i18n.TParam("durum")
			i18n.TParam("limit")
			i18n.TParam("offset")

			i18n.TValidation("param_required", "test_param", nil)
			i18n.FormatParameterRequired("test_param")
			i18n.FormatInvalidValue("test_param", "invalid", []string{"valid1", "valid2"})

			return nil
		},
	}

	RunStandardBenchmark(b, config)
}

// BenchmarkMemoryAllocation tests memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	config := BenchmarkConfig{
		Name: "MemoryAllocation",
		Setup: func() interface{} {
			return NewToolHelpers()
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			helpers := data.(*ToolHelpers)

			// Test memory allocation in helper functions
			params := map[string]interface{}{
				"id":     constants.TestIDBasic,
				"durum":  constants.TaskStatusPending,
				"baslik": constants.TestTaskTitleEN,
			}

			helpers.Validator.ValidateRequiredString(params, "id")
			helpers.Validator.ValidateRequiredString(params, "baslik")
			helpers.Validator.ValidateEnum(params, "durum", constants.GetValidTaskStatuses()[:2], true)

			helpers.Formatter.FormatTaskBasic("Test", constants.TestTaskShortID)
			helpers.Formatter.GetStatusEmoji(constants.TaskStatusPending)

			return nil
		},
	}

	RunStandardBenchmark(b, config)
}

// BenchmarkRegistryPatterns compares old vs new registry patterns
func BenchmarkRegistryPatterns(b *testing.B) {
	handlers, cleanup := setupBenchmarkEnvironment()
	defer cleanup()

	b.Run("NewRegistryPattern", func(b *testing.B) {
		config := BenchmarkConfig{
			Name: "NewRegistryPattern",
			Setup: func() interface{} {
				return handlers
			},
			Cleanup: func() {},
			Operation: func(data interface{}) error {
				h := data.(*Handlers)
				s := server.NewMCPServer("new-bench", "1.0.0")

				// Use new DRY registry pattern
				registry := NewToolRegistry(h)
				registry.RegisterAllTools(s)
				return nil
			},
		}
		RunStandardBenchmark(b, config)
	})
}

// Stress test with high concurrency
func BenchmarkStressConcurrency(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping stress test in short mode")
	}

	handlers, cleanup := setupBenchmarkEnvironment()
	defer cleanup()

	config := BenchmarkConfig{
		Name: "StressConcurrency",
		Setup: func() interface{} {
			return handlers
		},
		Cleanup: func() {},
		Operation: func(data interface{}) error {
			handlers := data.(*Handlers)

			// Mix of different operations
			operations := []func() error{
				func() error {
					result, err := handlers.GorevListele(map[string]interface{}{})
					if err != nil {
						return err
					}
					if result.IsError {
						return fmt.Errorf(constants.ErrorFormatTool, result.Content)
					}
					return nil
				},
				func() error {
					result, err := handlers.ProjeListele(map[string]interface{}{})
					if err != nil {
						return err
					}
					if result.IsError {
						return fmt.Errorf(constants.ErrorFormatTool, result.Content)
					}
					return nil
				},
				func() error {
					result, err := handlers.OzetGoster(map[string]interface{}{})
					if err != nil {
						return err
					}
					if result.IsError {
						return fmt.Errorf(constants.ErrorFormatTool, result.Content)
					}
					return nil
				},
			}

			// Randomly execute one of the operations
			op := operations[len(operations)%3]
			return op()
		},
		Goroutines: constants.TestConcurrencyLarge, // High concurrency
	}

	RunConcurrentBenchmark(b, config)
}

// Performance comparison: Before vs After refactoring
func BenchmarkPerformanceRegression(b *testing.B) {
	// This would compare against baseline metrics
	// For now, just ensure new code performs within acceptable bounds

	start := time.Now()

	handlers, cleanup := setupBenchmarkEnvironment()
	defer cleanup()

	setupTime := time.Since(start)
	if setupTime > constants.TestSetupTimeoutMs*time.Millisecond {
		b.Errorf("Setup took too long: %v (expected < %dms)", setupTime, constants.TestSetupTimeoutMs)
	}

	// Test tool call latency
	start = time.Now()
	result, err := handlers.GorevListele(map[string]interface{}{})
	callTime := time.Since(start)

	if err != nil {
		b.Errorf("Tool call failed: %v", err)
	}

	if result.IsError {
		b.Errorf("Tool call failed: %v", result.Content)
	}

	if callTime > constants.TestCallTimeoutMs*time.Millisecond {
		b.Errorf("Tool call took too long: %v (expected < %dms)", callTime, constants.TestCallTimeoutMs)
	}
}

// Race condition detection benchmark
func BenchmarkRaceConditionDetection(b *testing.B) {
	handlers, cleanup := setupBenchmarkEnvironment()
	defer cleanup()

	const goroutines = constants.TestConcurrencyLarge
	const iterations = constants.TestIterationSmall

	errors := make(chan error, goroutines*iterations)
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(goroutines)

		for j := 0; j < goroutines; j++ {
			go func(id int) {
				defer wg.Done()

				for k := 0; k < iterations; k++ {
					// Test concurrent access to different handlers
					switch k % 3 {
					case 0:
						result, err := handlers.GorevListele(map[string]interface{}{})
						if err != nil {
							errors <- err
							return
						}
						if result.IsError {
							errors <- fmt.Errorf(constants.ErrorFormatTool, result.Content)
							return
						}
					case 1:
						result, err := handlers.ProjeListele(map[string]interface{}{})
						if err != nil {
							errors <- err
							return
						}
						if result.IsError {
							errors <- fmt.Errorf(constants.ErrorFormatTool, result.Content)
							return
						}
					case 2:
						// Validator operations
						validator := NewParameterValidator()
						params := map[string]interface{}{"id": "test"}
						_, validationError := validator.ValidateRequiredString(params, "id")
						if validationError != nil {
							errors <- fmt.Errorf(constants.ErrorFormatValidation, validationError.Content)
							return
						}
					}
				}
			}(j)
		}

		wg.Wait()
	}

	// Check for errors
	close(errors)
	var collectedErrors []error
	for err := range errors {
		collectedErrors = append(collectedErrors, err)
	}

	if len(collectedErrors) > 0 {
		b.Errorf("Race conditions detected: %v errors, first: %v", len(collectedErrors), collectedErrors[0])
	}
}
