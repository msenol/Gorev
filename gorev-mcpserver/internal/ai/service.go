// Package ai provides AI provider integrations for task management
// This file contains the AI service orchestrator
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/msenol/gorev/internal/ai/providers"
)

// AIService manages AI operations across providers with graceful fallback
// Implements Rule 15 compliance - never degrades existing functionality
type AIService struct {
	registry       *providers.Registry
	promptManager  *PromptManager
	logger         *slog.Logger
}

// NewAIService creates a new AI service
func NewAIService(logger *slog.Logger) *AIService {
	if logger == nil {
		logger = slog.Default()
	}
	return &AIService{
		registry:      providers.GetRegistry(),
		promptManager: NewPromptManager(),
		logger:        logger,
	}
}

// ProjectAIConfig represents AI configuration for a project
type ProjectAIConfig struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Provider    string    `json:"provider"`
	APIKey      string    `json:"api_key"`
	Model       string    `json:"model"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SaveConfig saves AI configuration for a project and registers the provider
func (s *AIService) SaveConfig(ctx context.Context, config *ProjectAIConfig) error {
	if config.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}
	if config.Provider == "" {
		return fmt.Errorf("provider is required")
	}
	if config.APIKey == "" {
		return fmt.Errorf("api_key is required")
	}

	// Set defaults
	if config.Model == "" {
		config.Model = "openai/gpt-4o-mini"
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 2000
	}

	// Create and register provider
	providerType := providers.ProviderType(config.Provider)
	if err := s.registry.RegisterFromConfig(config.ProjectID, providerType, config.APIKey); err != nil {
		return fmt.Errorf("register provider: %w", err)
	}

	s.logger.Info("AI configuration saved",
		"project_id", config.ProjectID,
		"provider", config.Provider,
		"model", config.Model)

	return nil
}

// IsConfiguredForProject checks if AI is configured for a project
func (s *AIService) IsConfiguredForProject(projectID string) bool {
	_, ok := s.registry.Get(projectID)
	return ok
}

// GetConfig retrieves AI configuration for a project
func (s *AIService) GetConfig(ctx context.Context, projectID string) (*ProjectAIConfig, error) {
	_, ok := s.registry.Get(projectID)
	if !ok {
		return nil, fmt.Errorf("AI not configured for project: %s", projectID)
	}

	info, ok := s.registry.GetInfo(projectID)
	if !ok {
		return nil, fmt.Errorf("provider info not found: %s", projectID)
	}

	return &ProjectAIConfig{
		ProjectID: projectID,
		Provider:  info.Name,
		// API key is not returned for security
		Model:     "openai/gpt-4o-mini", // Would be retrieved from DB
	}, nil
}

// Chat sends a chat message and returns the response
func (s *AIService) Chat(ctx context.Context, projectID, message string) (string, error) {
	if !s.IsConfiguredForProject(projectID) {
		return "", ErrAINotConfigured
	}

	provider, ok := s.registry.Get(projectID)
	if !ok {
		return "", fmt.Errorf("provider not found for project: %s", projectID)
	}

	// Get recommended model
	var providerImpl interface {
		GetRecommendedModel(AIOperation) string
	}
	if p, ok := provider.(interface{ GetRecommendedModel(AIOperation) string }); ok {
		providerImpl = p
	}

	model := "openai/gpt-4o-mini"
	if providerImpl != nil {
		model = providerImpl.GetRecommendedModel(AIOperationChat)
	}

	req := providers.ChatRequest{
		Model: model,
		Messages: []providers.ChatMessage{
			{Role: "user", Content: message},
		},
		MaxTokens:   2000,
		Temperature: 0.7,
	}

	resp, err := provider.ChatCompletion(ctx, req)
	if err != nil {
		s.logger.Warn("AI chat failed", "error", err, "project_id", projectID)
		return "", fmt.Errorf("AI unavailable: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GetSuggestions generates AI-powered task suggestions
func (s *AIService) GetSuggestions(ctx context.Context, projectID string, maxResults int, contextStr string) ([]Suggestion, error) {
	if !s.IsConfiguredForProject(projectID) {
		return nil, ErrAINotConfigured
	}

	// TODO: Implement suggestion generation
	// For now, return empty slice
	return []Suggestion{}, nil
}

// DecomposeTask breaks down a complex task into subtasks
func (s *AIService) DecomposeTask(ctx context.Context, projectID, taskID, title, description string, maxDepth int) ([]SubtaskSuggestion, error) {
	if !s.IsConfiguredForProject(projectID) {
		return nil, ErrAINotConfigured
	}

	provider, ok := s.registry.Get(projectID)
	if !ok {
		return nil, fmt.Errorf("provider not found for project: %s", projectID)
	}

	// Get recommended model
	var providerImpl interface {
		GetRecommendedModel(AIOperation) string
	}
	if p, ok := provider.(interface{ GetRecommendedModel(AIOperation) string }); ok {
		providerImpl = p
	}

	model := "openai/gpt-4o-mini"
	if providerImpl != nil {
		model = providerImpl.GetRecommendedModel(AIOperationTaskDecomposition)
	}

	// Get prompts from prompt manager
	systemPrompt, _, err := s.promptManager.GetPrompt(ctx, AIOperationTaskDecomposition, "tr")
	if err != nil {
		return nil, fmt.Errorf("get prompt: %w", err)
	}

	userPrompt := s.promptManager.BuildTaskDecompositionPrompt("tr", title, description)

	req := providers.ChatRequest{
		Model: model,
		Messages: []providers.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   2000,
		Temperature: 0.7,
	}

	resp, err := provider.ChatCompletion(ctx, req)
	if err != nil {
		s.logger.Warn("AI decomposition failed", "error", err, "project_id", projectID, "task_id", taskID)
		return nil, fmt.Errorf("AI unavailable: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Parse response
	var result struct {
		Subtasks []SubtaskSuggestion `json:"subtasks"`
	}

	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("parse AI response: %w", err)
	}

	return result.Subtasks, nil
}

// SemanticSearch performs semantic search across tasks
func (s *AIService) SemanticSearch(ctx context.Context, projectID, query string, limit int) (*SearchResult, error) {
	if !s.IsConfiguredForProject(projectID) {
		return nil, ErrAINotConfigured
	}

	// TODO: Implement semantic search with embeddings
	// For now, return empty result
	return &SearchResult{
		Total:   0,
		Results: []SearchMatch{},
		Strategy: "semantic",
	}, nil
}

// EstimateTime estimates task duration using AI
func (s *AIService) EstimateTime(ctx context.Context, projectID, taskID, title, description string, tags []string) (*EstimationResult, error) {
	if !s.IsConfiguredForProject(projectID) {
		return nil, ErrAINotConfigured
	}

	provider, ok := s.registry.Get(projectID)
	if !ok {
		return nil, fmt.Errorf("provider not found for project: %s", projectID)
	}

	// Get recommended model
	var providerImpl interface {
		GetRecommendedModel(AIOperation) string
	}
	if p, ok := provider.(interface{ GetRecommendedModel(AIOperation) string }); ok {
		providerImpl = p
	}

	model := "openai/gpt-4o-mini"
	if providerImpl != nil {
		model = providerImpl.GetRecommendedModel(AIOperationTimeEstimation)
	}

	// Get prompts from prompt manager
	systemPrompt, _, err := s.promptManager.GetPrompt(ctx, AIOperationTimeEstimation, "tr")
	if err != nil {
		return nil, fmt.Errorf("get prompt: %w", err)
	}

	userPrompt := s.promptManager.BuildTimeEstimationPrompt("tr", title, description, tags)

	req := providers.ChatRequest{
		Model: model,
		Messages: []providers.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   500,
		Temperature: 0.3, // Lower temperature for more deterministic estimates
	}

	resp, err := provider.ChatCompletion(ctx, req)
	if err != nil {
		s.logger.Warn("AI estimation failed", "error", err, "project_id", projectID, "task_id", taskID)
		return nil, fmt.Errorf("AI unavailable: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Parse response
	var result EstimationResult
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("parse AI response: %w", err)
	}

	result.Method = "ai"
	return &result, nil
}

// AnalyzeProject performs project analytics using AI
func (s *AIService) AnalyzeProject(ctx context.Context, projectID, projectName string, taskCount, completedCount, pendingCount int) (*ProjectAnalysisResult, error) {
	if !s.IsConfiguredForProject(projectID) {
		return nil, ErrAINotConfigured
	}

	provider, ok := s.registry.Get(projectID)
	if !ok {
		return nil, fmt.Errorf("provider not found for project: %s", projectID)
	}

	// Get recommended model
	var providerImpl interface {
		GetRecommendedModel(AIOperation) string
	}
	if p, ok := provider.(interface{ GetRecommendedModel(AIOperation) string }); ok {
		providerImpl = p
	}

	model := "openai/gpt-4o"
	if providerImpl != nil {
		model = providerImpl.GetRecommendedModel(AIOperationProjectAnalytics)
	}

	// Get prompts from prompt manager
	systemPrompt, _, err := s.promptManager.GetPrompt(ctx, AIOperationProjectAnalytics, "tr")
	if err != nil {
		return nil, fmt.Errorf("get prompt: %w", err)
	}

	userPrompt := s.promptManager.BuildProjectAnalyticsPrompt("tr", projectName, taskCount, completedCount, pendingCount)

	req := providers.ChatRequest{
		Model: model,
		Messages: []providers.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   1500,
		Temperature: 0.5,
	}

	resp, err := provider.ChatCompletion(ctx, req)
	if err != nil {
		s.logger.Warn("AI analysis failed", "error", err, "project_id", projectID)
		return nil, fmt.Errorf("AI unavailable: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Parse response
	var result ProjectAnalysisResult
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("parse AI response: %w", err)
	}

	return &result, nil
}

// TestConnection tests the AI provider connection
func (s *AIService) TestConnection(ctx context.Context, projectID string) error {
	provider, ok := s.registry.Get(projectID)
	if !ok {
		return fmt.Errorf("provider not found for project: %s", projectID)
	}

	return provider.HealthCheck(ctx)
}

// SubtaskSuggestion represents a suggested subtask from decomposition
type SubtaskSuggestion struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	EstimatedHours float64  `json:"estimated_hours"`
	Dependencies   []string `json:"dependencies"`
}

// ProjectAnalysisResult represents AI-powered project analytics
type ProjectAnalysisResult struct {
	CriticalPath   []string `json:"critical_path"`
	RiskLevel      string   `json:"risk_level"`
	RiskFactors    []string `json:"risk_factors"`
	Recommendations []string `json:"recommendations"`
	Bottlenecks    []string `json:"bottlenecks"`
}

// Errors
var (
	ErrAINotConfigured = fmt.Errorf("AI not configured for project")
	ErrAINotAvailable  = fmt.Errorf("AI service unavailable")
)
