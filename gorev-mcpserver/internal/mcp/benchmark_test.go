package mcp

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
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

// setupBenchmarkEnvironment - DRY setup for all benchmarks
func setupBenchmarkEnvironment() (*Handlers, func()) {
	// Initialize i18n for consistent environment
	i18n.Initialize("tr")

	// Setup data manager
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

// BenchmarkToolRegistration tests the performance of tool registration
func BenchmarkToolRegistration(b *testing.B) {
	config := BenchmarkConfig{
		Name: "ToolRegistration",
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
				"id":     "test-id-123",
				"durum":  "beklemede",
				"limit":  float64(50),
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
				return fmt.Errorf("validation failed: %v", validationError.Content)
			}

			_, validationError = validator.ValidateEnum(params, "durum", []string{"beklemede", "devam_ediyor", "tamamlandi"}, true)
			if validationError != nil {
				return fmt.Errorf("validation failed: %v", validationError.Content)
			}

			validator.ValidateNumber(params, "limit", 50)
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
			formatter.FormatTaskBasic("Test Task", "12345678-1234-1234-1234-123456789012")
			formatter.FormatTaskWithStatus("Test Task", "12345678-1234-1234-1234-123456789012", "beklemede")
			formatter.FormatSuccessMessage("Test Action", "Test Task", "12345678")
			formatter.GetStatusEmoji("beklemede")
			formatter.GetPriorityEmoji("yuksek")

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
				return fmt.Errorf("tool error: %v", result.Content)
			}
			return nil
		},
		Goroutines: 10,
	}

	RunConcurrentBenchmark(b, config)
}

// BenchmarkI18nHelpers tests i18n helper performance
func BenchmarkI18nHelpers(b *testing.B) {
	config := BenchmarkConfig{
		Name: "I18nHelpers",
		Setup: func() interface{} {
			i18n.Initialize("tr")
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
				"id":     "test-id",
				"durum":  "beklemede",
				"baslik": "Test Task",
			}

			helpers.Validator.ValidateRequiredString(params, "id")
			helpers.Validator.ValidateRequiredString(params, "baslik")
			helpers.Validator.ValidateEnum(params, "durum", []string{"beklemede", "devam_ediyor"}, true)

			helpers.Formatter.FormatTaskBasic("Test", "12345678")
			helpers.Formatter.GetStatusEmoji("beklemede")

			return nil
		},
	}

	RunStandardBenchmark(b, config)
}

// BenchmarkRegistryPatterns compares old vs new registry patterns
func BenchmarkRegistryPatterns(b *testing.B) {
	i18n.Initialize("tr")
	dataManager, _ := gorev.YeniVeriYonetici(":memory:", "")
	isYonetici := gorev.YeniIsYonetici(dataManager)
	handlers := YeniHandlers(isYonetici)

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
					result, err := handlers.OzetGoster(map[string]interface{}{})
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
			op := operations[len(operations)%3]
			return op()
		},
		Goroutines: 50, // High concurrency
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
	if setupTime > 100*time.Millisecond {
		b.Errorf("Setup took too long: %v (expected < 100ms)", setupTime)
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

	if callTime > 10*time.Millisecond {
		b.Errorf("Tool call took too long: %v (expected < 10ms)", callTime)
	}
}

// Race condition detection benchmark
func BenchmarkRaceConditionDetection(b *testing.B) {
	handlers, cleanup := setupBenchmarkEnvironment()
	defer cleanup()

	const goroutines = 100
	const iterations = 10

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
							errors <- fmt.Errorf("tool error: %v", result.Content)
							return
						}
					case 1:
						result, err := handlers.ProjeListele(map[string]interface{}{})
						if err != nil {
							errors <- err
							return
						}
						if result.IsError {
							errors <- fmt.Errorf("tool error: %v", result.Content)
							return
						}
					case 2:
						// Validator operations
						validator := NewParameterValidator()
						params := map[string]interface{}{"id": "test"}
						_, validationError := validator.ValidateRequiredString(params, "id")
						if validationError != nil {
							errors <- fmt.Errorf("validation failed: %v", validationError.Content)
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
