// Package ai provides tests for AI service with graceful degradation scenarios
package ai

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/ai/providers"
)

// MockProvider is a mock AI provider for testing
type MockProvider struct {
	name       string
	providerType providers.ProviderType
	shouldFail bool
	delay      time.Duration
}

func (m *MockProvider) Name() string {
	return m.name
}

func (m *MockProvider) Type() providers.ProviderType {
	return m.providerType
}

func (m *MockProvider) BaseURL() string {
	return "https://mock.example.com"
}

func (m *MockProvider) ListModels(ctx context.Context) ([]providers.ModelInfo, error) {
	if m.shouldFail {
		return nil, errors.New("provider unavailable")
	}
	return []providers.ModelInfo{
		{ID: "mock/gpt-4o-mini", Name: "Mock GPT-4o Mini", ContextWindow: 128000},
	}, nil
}

func (m *MockProvider) ValidateModel(ctx context.Context, model string) error {
	return nil
}

func (m *MockProvider) ChatCompletion(ctx context.Context, req providers.ChatRequest) (*providers.ChatResponse, error) {
	if m.shouldFail {
		return nil, errors.New("provider unavailable")
	}

	// Simulate delay if set
	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(m.delay):
		}
	}

	// Return mock response based on operation
	return &providers.ChatResponse{
		Choices: []providers.ChatChoice{
			{
				Message: providers.ChatMessage{
					Role:    "assistant",
					Content: m.getMockResponse(req),
				},
			},
		},
		Usage: providers.ChatUsage{
			TotalTokens: 100,
		},
	}, nil
}

func (m *MockProvider) HealthCheck(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("provider unavailable")
	}
	return nil
}

func (m *MockProvider) Supports(feature string) bool {
	return true
}

func (m *MockProvider) GetRecommendedModel(operation AIOperation) string {
	return "mock/gpt-4o-mini"
}

// getMockResponse returns appropriate mock responses based on the request
func (m *MockProvider) getMockResponse(req providers.ChatRequest) string {
	// Check system prompt to determine operation
	for _, msg := range req.Messages {
		if msg.Role == "system" {
			content := msg.Content
			switch {
			case containsSubstring(content, "breaking down complex tasks") || containsSubstring(content, "karmaşık görevleri"):
				return `{
  "subtasks": [
    {
      "title": "Subtask 1",
      "description": "First subtask",
      "estimated_hours": 2.0,
      "dependencies": []
    },
    {
      "title": "Subtask 2",
      "description": "Second subtask",
      "estimated_hours": 1.5,
      "dependencies": ["subtask-1"]
    }
  ]
}`
			case containsSubstring(content, "estimating task durations") || containsSubstring(content, "görev süreleri tahmin"):
				return `{
  "estimated_hours": 4.5,
  "confidence": 0.8,
  "reasoning": "Based on similar tasks"
}`
			case containsSubstring(content, "project analytics") || containsSubstring(content, "proje analizi"):
				return `{
  "critical_path": ["task-1", "task-2"],
  "risk_level": "medium",
  "risk_factors": ["tight deadline", "complex dependencies"],
  "recommendations": ["add buffer time", "assign more resources"],
  "bottlenecks": ["frontend resources"]
}`
			}
		}
	}

	// Default chat response
	return "Mock AI response"
}

// containsSubstring checks if a string contains a substring (case-insensitive)
func containsSubstring(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return len(s) >= len(substr) && findSubstring(s, substr)
}

// findSubstring finds a substring in a string
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// toLower converts a string to lowercase
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		   (s == substr ||
		    len(s) > len(substr) && (
		       s[:len(substr)] == substr ||
		       s[len(s)-len(substr):] == substr ||
		       containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestAIServiceChat tests the Chat method with various scenarios
func TestAIServiceChat(t *testing.T) {
	registry := providers.GetRegistry()
	pm := NewPromptManager()

	tests := []struct {
		name        string
		projectID   string
		shouldFail  bool
		expectError bool
	}{
		{
			name:        "successful chat",
			projectID:   "test-project-1",
			shouldFail:  false,
			expectError: false,
		},
		{
			name:        "AI not configured",
			projectID:   "non-existent-project",
			shouldFail:  false,
			expectError: true,
		},
		{
			name:        "AI provider failure",
			projectID:   "test-project-fail",
			shouldFail:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test service with logger
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := &AIService{
				registry:      registry,
				promptManager: pm,
				logger:        logger,
			}

			// Register mock provider if needed
			if !tt.shouldFail && tt.projectID != "non-existent-project" {
				provider := &MockProvider{name: "mock", providerType: providers.ProviderOpenRouter, shouldFail: false}
				registry.Register(tt.projectID, provider)
			} else if tt.shouldFail {
				provider := &MockProvider{name: "mock-fail", providerType: providers.ProviderOpenRouter, shouldFail: true}
				registry.Register(tt.projectID, provider)
			}

			// Test chat
			ctx := context.Background()
			response, err := service.Chat(ctx, tt.projectID, "Hello AI")

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if response == "" {
				t.Errorf("Expected non-empty response")
			}
		})
	}
}

// TestAIServiceDecomposeTask tests task decomposition with various scenarios
func TestAIServiceDecomposeTask(t *testing.T) {
	registry := providers.GetRegistry()
	pm := NewPromptManager()

	tests := []struct {
		name        string
		projectID   string
		title       string
		description string
		shouldFail  bool
		expectError bool
		expectCount int
	}{
		{
			name:        "successful decomposition",
			projectID:   "test-project-1",
			title:       "Build web application",
			description: "Create a full-stack web application with React frontend",
			shouldFail:  false,
			expectError: false,
			expectCount: 2,
		},
		{
			name:        "AI not configured",
			projectID:   "non-existent-project",
			title:       "Some task",
			description: "Some description",
			shouldFail:  false,
			expectError: true,
			expectCount: 0,
		},
		{
			name:        "AI provider failure",
			projectID:   "test-project-fail",
			title:       "Some task",
			description: "Some description",
			shouldFail:  true,
			expectError: true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test service with logger
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := &AIService{
				registry:      registry,
				promptManager: pm,
				logger:        logger,
			}

			// Register mock provider if needed
			if !tt.shouldFail && tt.projectID != "non-existent-project" {
				provider := &MockProvider{name: "mock", providerType: providers.ProviderOpenRouter, shouldFail: false}
				registry.Register(tt.projectID, provider)
			} else if tt.shouldFail {
				provider := &MockProvider{name: "mock-fail", providerType: providers.ProviderOpenRouter, shouldFail: true}
				registry.Register(tt.projectID, provider)
			}

			// Test decomposition
			ctx := context.Background()
			subtasks, err := service.DecomposeTask(ctx, tt.projectID, "task-1", tt.title, tt.description, 3)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(subtasks) != tt.expectCount {
				t.Errorf("Expected %d subtasks, got %d", tt.expectCount, len(subtasks))
			}
		})
	}
}

// TestAIServiceEstimateTime tests time estimation with various scenarios
func TestAIServiceEstimateTime(t *testing.T) {
	registry := providers.GetRegistry()
	pm := NewPromptManager()

	tests := []struct {
		name        string
		projectID   string
		title       string
		description string
		tags        []string
		shouldFail  bool
		expectError bool
	}{
		{
			name:        "successful estimation",
			projectID:   "test-project-1",
			title:       "API Development",
			description: "Build REST API endpoints",
			tags:        []string{"backend", "api"},
			shouldFail:  false,
			expectError: false,
		},
		{
			name:        "AI not configured",
			projectID:   "non-existent-project",
			title:       "Some task",
			description: "Some description",
			tags:        []string{},
			shouldFail:  false,
			expectError: true,
		},
		{
			name:        "AI provider failure",
			projectID:   "test-project-fail",
			title:       "Some task",
			description: "Some description",
			tags:        []string{},
			shouldFail:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test service with logger
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := &AIService{
				registry:      registry,
				promptManager: pm,
				logger:        logger,
			}

			// Register mock provider if needed
			if !tt.shouldFail && tt.projectID != "non-existent-project" {
				provider := &MockProvider{name: "mock", providerType: providers.ProviderOpenRouter, shouldFail: false}
				registry.Register(tt.projectID, provider)
			} else if tt.shouldFail {
				provider := &MockProvider{name: "mock-fail", providerType: providers.ProviderOpenRouter, shouldFail: true}
				registry.Register(tt.projectID, provider)
			}

			// Test estimation
			ctx := context.Background()
			result, err := service.EstimateTime(ctx, tt.projectID, "task-1", tt.title, tt.description, tt.tags)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.EstimatedHours <= 0 {
				t.Errorf("Expected positive estimated hours, got %f", result.EstimatedHours)
			}

			if result.Method != "ai" {
				t.Errorf("Expected method 'ai', got '%s'", result.Method)
			}
		})
	}
}

// TestAIServiceAnalyzeProject tests project analytics with various scenarios
func TestAIServiceAnalyzeProject(t *testing.T) {
	registry := providers.GetRegistry()
	pm := NewPromptManager()

	tests := []struct {
		name          string
		projectID     string
		projectName   string
		taskCount     int
		completed     int
		pending       int
		shouldFail    bool
		expectError   bool
		expectContent bool
	}{
		{
			name:          "successful analysis",
			projectID:     "test-project-1",
			projectName:   "Web App",
			taskCount:     10,
			completed:     5,
			pending:       5,
			shouldFail:    false,
			expectError:   false,
			expectContent: true,
		},
		{
			name:          "AI not configured",
			projectID:     "non-existent-project",
			projectName:   "Some Project",
			taskCount:     5,
			completed:     2,
			pending:       3,
			shouldFail:    false,
			expectError:   true,
			expectContent: false,
		},
		{
			name:          "AI provider failure",
			projectID:     "test-project-fail",
			projectName:   "Some Project",
			taskCount:     5,
			completed:     2,
			pending:       3,
			shouldFail:    true,
			expectError:   true,
			expectContent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test service with logger
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			service := &AIService{
				registry:      registry,
				promptManager: pm,
				logger:        logger,
			}

			// Register mock provider if needed
			if !tt.shouldFail && tt.projectID != "non-existent-project" {
				provider := &MockProvider{name: "mock", providerType: providers.ProviderOpenRouter, shouldFail: false}
				registry.Register(tt.projectID, provider)
			} else if tt.shouldFail {
				provider := &MockProvider{name: "mock-fail", providerType: providers.ProviderOpenRouter, shouldFail: true}
				registry.Register(tt.projectID, provider)
			}

			// Test analysis
			ctx := context.Background()
			result, err := service.AnalyzeProject(ctx, tt.projectID, tt.projectName, tt.taskCount, tt.completed, tt.pending)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.expectContent {
				if len(result.CriticalPath) == 0 {
					t.Errorf("Expected non-empty critical path")
				}

				if result.RiskLevel == "" {
					t.Errorf("Expected risk level to be set")
				}

				if len(result.Recommendations) == 0 {
					t.Errorf("Expected non-empty recommendations")
				}
			}
		})
	}
}

// TestAIServiceTimeout tests AI service timeout handling
func TestAIServiceTimeout(t *testing.T) {
	registry := providers.GetRegistry()
	pm := NewPromptManager()

	// Create test service with slow provider
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := &AIService{
		registry:      registry,
		promptManager: pm,
		logger:        logger,
	}

	// Register a slow provider that will timeout
	provider := &MockProvider{
		name:         "slow-mock",
		providerType: providers.ProviderOpenRouter,
		delay:        5 * time.Second, // Longer than test timeout
	}
	registry.Register("slow-project", provider)

	// Test with timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := service.Chat(ctx, "slow-project", "Hello")
	if err == nil {
		t.Errorf("Expected timeout error but got none")
	}
}

// TestAIServiceConcurrentAccess tests concurrent AI service access
func TestAIServiceConcurrentAccess(t *testing.T) {
	registry := providers.GetRegistry()
	pm := NewPromptManager()

	// Create test service
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	service := &AIService{
		registry:      registry,
		promptManager: pm,
		logger:        logger,
	}

	// Register provider
	provider := &MockProvider{name: "mock", providerType: providers.ProviderOpenRouter, shouldFail: false}
	registry.Register("concurrent-project", provider)

	// Launch concurrent requests
	ctx := context.Background()
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			_, err := service.Chat(ctx, "concurrent-project", "Concurrent test")
			if err != nil {
				t.Errorf("Concurrent request failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Errorf("Timeout waiting for concurrent requests")
		}
	}
}

// TestPromptManager tests the PromptManager with i18n support
func TestPromptManager(t *testing.T) {
	pm := NewPromptManager()

	tests := []struct {
		name        string
		operation   AIOperation
		lang        string
		expectError bool
	}{
		{
			name:        "Turkish task decomposition",
			operation:   AIOperationTaskDecomposition,
			lang:        "tr",
			expectError: false,
		},
		{
			name:        "English task decomposition",
			operation:   AIOperationTaskDecomposition,
			lang:        "en",
			expectError: false,
		},
		{
			name:        "Turkish time estimation",
			operation:   AIOperationTimeEstimation,
			lang:        "tr",
			expectError: false,
		},
		{
			name:        "English time estimation",
			operation:   AIOperationTimeEstimation,
			lang:        "en",
			expectError: false,
		},
		{
			name:        "Turkish project analytics",
			operation:   AIOperationProjectAnalytics,
			lang:        "tr",
			expectError: false,
		},
		{
			name:        "English project analytics",
			operation:   AIOperationProjectAnalytics,
			lang:        "en",
			expectError: false,
		},
		{
			name:        "Invalid operation",
			operation:   "invalid",
			lang:        "en",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			system, _, err := pm.GetPrompt(ctx, tt.operation, tt.lang)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if system == "" {
				t.Errorf("Expected non-empty system prompt")
			}
		})
	}
}

// TestPromptManagerBuildPrompts tests the prompt building methods
func TestPromptManagerBuildPrompts(t *testing.T) {
	pm := NewPromptManager()

	t.Run("BuildTaskDecompositionPrompt Turkish", func(t *testing.T) {
		prompt := pm.BuildTaskDecompositionPrompt("tr", "Test Task", "Test Description")

		if !contains(prompt, "Test Task") {
			t.Errorf("Expected prompt to contain task title")
		}
		if !contains(prompt, "Test Description") {
			t.Errorf("Expected prompt to contain task description")
		}
	})

	t.Run("BuildTaskDecompositionPrompt English", func(t *testing.T) {
		prompt := pm.BuildTaskDecompositionPrompt("en", "Test Task", "Test Description")

		if !contains(prompt, "Test Task") {
			t.Errorf("Expected prompt to contain task title")
		}
		if !contains(prompt, "subtasks") {
			t.Errorf("Expected prompt to mention subtasks")
		}
	})

	t.Run("BuildTimeEstimationPrompt", func(t *testing.T) {
		tags := []string{"backend", "api"}
		prompt := pm.BuildTimeEstimationPrompt("tr", "API Task", "Build API", tags)

		if !contains(prompt, "API Task") {
			t.Errorf("Expected prompt to contain task title")
		}
		if !contains(prompt, "estimated_hours") {
			t.Errorf("Expected prompt to request estimated_hours")
		}
	})

	t.Run("BuildProjectAnalyticsPrompt", func(t *testing.T) {
		prompt := pm.BuildProjectAnalyticsPrompt("en", "Web Project", 10, 5, 5)

		if !contains(prompt, "Web Project") {
			t.Errorf("Expected prompt to contain project name")
		}
		if !contains(prompt, "10") {
			t.Errorf("Expected prompt to contain total task count")
		}
		if !contains(prompt, "critical_path") {
			t.Errorf("Expected prompt to request critical_path")
		}
	})
}

// TestEstimationResultUnmarshal tests JSON unmarshaling for EstimationResult
func TestEstimationResultUnmarshal(t *testing.T) {
	jsonData := `{
		"estimated_hours": 4.5,
		"confidence": 0.8,
		"reasoning": "Based on similar tasks",
		"method": "ai"
	}`

	var result EstimationResult
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal EstimationResult: %v", err)
	}

	if result.EstimatedHours != 4.5 {
		t.Errorf("Expected EstimatedHours 4.5, got %f", result.EstimatedHours)
	}

	if result.Confidence != 0.8 {
		t.Errorf("Expected Confidence 0.8, got %f", result.Confidence)
	}

	if result.Method != "ai" {
		t.Errorf("Expected method 'ai', got '%s'", result.Method)
	}
}

// TestProjectAnalysisResultUnmarshal tests JSON unmarshaling for ProjectAnalysisResult
func TestProjectAnalysisResultUnmarshal(t *testing.T) {
	jsonData := `{
		"critical_path": ["task-1", "task-2"],
		"risk_level": "medium",
		"risk_factors": ["deadline", "complexity"],
		"recommendations": ["add buffer"],
		"bottlenecks": ["resources"]
	}`

	var result ProjectAnalysisResult
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal ProjectAnalysisResult: %v", err)
	}

	if len(result.CriticalPath) != 2 {
		t.Errorf("Expected 2 critical path items, got %d", len(result.CriticalPath))
	}

	if result.RiskLevel != "medium" {
		t.Errorf("Expected risk level 'medium', got '%s'", result.RiskLevel)
	}

	if len(result.RiskFactors) != 2 {
		t.Errorf("Expected 2 risk factors, got %d", len(result.RiskFactors))
	}
}
