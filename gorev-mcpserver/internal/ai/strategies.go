// Package ai provides AI provider integrations for task management
// This file contains the strategy pattern for AI/fallback (DRY compliance)
package ai

import (
	"context"
	"fmt"
)

// SuggestionStrategy defines the interface for task suggestions
// Single abstraction that delegates to AI or rule-based (DRY compliance)
type SuggestionStrategy interface {
	// GetSuggestions returns suggested tasks for the user
	GetSuggestions(ctx context.Context, input SuggestionInput) ([]Suggestion, error)
}

// SuggestionInput represents input for suggestion generation
type SuggestionInput struct {
	ProjectID  string
	MaxResults int
	Context    string // Additional context (file path, active task, etc.)
}

// Suggestion represents a suggested task or action
type Suggestion struct {
	Type        SuggestionType `json:"type"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    string         `json:"priority,omitempty"`
	TaskID      string         `json:"task_id,omitempty"`
	Confidence  float64        `json:"confidence"`
}

// SuggestionType represents the type of suggestion
type SuggestionType string

const (
	SuggestionTypeStartTask     SuggestionType = "start_task"
	SuggestionTypeCompleteTask  SuggestionType = "complete_task"
	SuggestionTypeCreateTask    SuggestionType = "create_task"
	SuggestionTypeReviewTask    SuggestionType = "review_task"
)

// AISuggestionStrategy uses AI to generate suggestions
type AISuggestionStrategy struct {
	aiService *AIService
}

// NewAISuggestionStrategy creates an AI-powered suggestion strategy
func NewAISuggestionStrategy(aiService *AIService) *AISuggestionStrategy {
	return &AISuggestionStrategy{aiService: aiService}
}

// GetSuggestions generates suggestions using AI
func (s *AISuggestionStrategy) GetSuggestions(ctx context.Context, input SuggestionInput) ([]Suggestion, error) {
	// Check if AI is configured
	if !s.aiService.IsConfiguredForProject(input.ProjectID) {
		return nil, fmt.Errorf("AI not configured for project")
	}

	// Use AI to generate suggestions
	// This will call the AI service with appropriate prompts
	return s.aiService.GetSuggestions(ctx, input.ProjectID, input.MaxResults, input.Context)
}

// RuleBasedSuggestionStrategy uses rule-based logic for suggestions
// This is the existing behavior (fallback when AI unavailable)
type RuleBasedSuggestionStrategy struct {
	// Existing suggestion engine would be injected here
	// For now, this is a placeholder that will be integrated with existing code
}

// NewRuleBasedSuggestionStrategy creates a rule-based suggestion strategy
func NewRuleBasedSuggestionStrategy() *RuleBasedSuggestionStrategy {
	return &RuleBasedSuggestionStrategy{}
}

// GetSuggestions generates suggestions using rule-based logic
func (s *RuleBasedSuggestionStrategy) GetSuggestions(ctx context.Context, input SuggestionInput) ([]Suggestion, error) {
	// This will integrate with the existing suggestion_engine.go
	// For now, return empty slice (no-op)
	return []Suggestion{}, nil
}

// SuggestionStrategyFactory creates the appropriate strategy based on AI configuration
// Single factory function ensures no code duplication (DRY)
func SuggestionStrategyFactory(aiService *AIService, projectID string) SuggestionStrategy {
	if aiService != nil && aiService.IsConfiguredForProject(projectID) {
		return NewAISuggestionStrategy(aiService)
	}
	return NewRuleBasedSuggestionStrategy()
}

// SearchStrategy defines the interface for task search
// Single abstraction that delegates to AI (semantic) or FTS (DRY compliance)
type SearchStrategy interface {
	Search(ctx context.Context, input SearchInput) (*SearchResult, error)
}

// SearchInput represents input for search
type SearchInput struct {
	Query     string
	ProjectID string
	Limit     int
	Offset    int
}

// SearchResult represents search results
type SearchResult struct {
	Total    int              `json:"total"`
	Results  []SearchMatch    `json:"results"`
	Strategy string           `json:"strategy"` // "semantic" or "fts"
}

// SearchMatch represents a single search result
type SearchMatch struct {
	TaskID        string  `json:"task_id"`
	Title         string  `json:"title"`
	RelevanceScore float64 `json:"relevance_score"`
	Highlight     string  `json:"highlight,omitempty"`
}

// SemanticSearchStrategy uses AI for semantic search
type SemanticSearchStrategy struct {
	aiService *AIService
}

// NewSemanticSearchStrategy creates an AI-powered search strategy
func NewSemanticSearchStrategy(aiService *AIService) *SemanticSearchStrategy {
	return &SemanticSearchStrategy{aiService: aiService}
}

// Search performs semantic search using AI embeddings
func (s *SemanticSearchStrategy) Search(ctx context.Context, input SearchInput) (*SearchResult, error) {
	if !s.aiService.IsConfiguredForProject(input.ProjectID) {
		return nil, fmt.Errorf("AI not configured for project")
	}

	return s.aiService.SemanticSearch(ctx, input.ProjectID, input.Query, input.Limit)
}

// FTSSearchStrategy uses full-text search (existing behavior)
type FTSSearchStrategy struct {
	// Existing FTS search would be injected here
}

// NewFTSSearchStrategy creates a full-text search strategy
func NewFTSSearchStrategy() *FTSSearchStrategy {
	return &FTSSearchStrategy{}
}

// Search performs full-text search
func (s *FTSSearchStrategy) Search(ctx context.Context, input SearchInput) (*SearchResult, error) {
	// This will integrate with the existing FTS search in gorevler_search table
	// For now, return empty result (no-op)
	return &SearchResult{
		Total:    0,
		Results:  []SearchMatch{},
		Strategy: "fts",
	}, nil
}

// SearchStrategyFactory creates the appropriate strategy based on AI configuration
// Single factory function ensures no code duplication (DRY)
func SearchStrategyFactory(aiService *AIService, projectID string) SearchStrategy {
	if aiService != nil && aiService.IsConfiguredForProject(projectID) && IsFeatureEnabled("semantic_search") {
		return NewSemanticSearchStrategy(aiService)
	}
	return NewFTSSearchStrategy()
}

// EstimationStrategy defines the interface for time estimation
// Single abstraction that delegates to AI or historical data (DRY compliance)
type EstimationStrategy interface {
	Estimate(ctx context.Context, input EstimationInput) (*EstimationResult, error)
}

// EstimationInput represents input for time estimation
type EstimationInput struct {
	TaskID     string
	Title      string
	Desc       string
	Tags       []string
	ProjectID  string
}

// EstimationResult represents time estimation result
type EstimationResult struct {
	EstimatedHours float64 `json:"estimated_hours"`
	Confidence     float64 `json:"confidence"`
	Reasoning      string  `json:"reasoning,omitempty"`
	Method         string  `json:"method"` // "ai" or "historical"
}

// AIEstimationStrategy uses AI for time estimation
type AIEstimationStrategy struct {
	aiService *AIService
}

// NewAIEstimationStrategy creates an AI-powered estimation strategy
func NewAIEstimationStrategy(aiService *AIService) *AIEstimationStrategy {
	return &AIEstimationStrategy{aiService: aiService}
}

// Estimate performs AI-based time estimation
func (s *AIEstimationStrategy) Estimate(ctx context.Context, input EstimationInput) (*EstimationResult, error) {
	if !s.aiService.IsConfiguredForProject(input.ProjectID) {
		return nil, fmt.Errorf("AI not configured for project")
	}

	return s.aiService.EstimateTime(ctx, input.ProjectID, input.TaskID, input.Title, input.Desc, input.Tags)
}

// HistoricalEstimationStrategy uses historical data (existing behavior)
type HistoricalEstimationStrategy struct {
	// Existing estimation logic would be injected here
}

// NewHistoricalEstimationStrategy creates a historical data estimation strategy
func NewHistoricalEstimationStrategy() *HistoricalEstimationStrategy {
	return &HistoricalEstimationStrategy{}
}

// Estimate performs estimation based on historical data
func (s *HistoricalEstimationStrategy) Estimate(ctx context.Context, input EstimationInput) (*EstimationResult, error) {
	// This will integrate with the existing time estimation logic
	// For now, return default estimation (no-op)
	return &EstimationResult{
		EstimatedHours: 2.0,
		Confidence:     0.5,
		Method:         "historical",
	}, nil
}

// EstimationStrategyFactory creates the appropriate strategy based on AI configuration
// Single factory function ensures no code duplication (DRY)
func EstimationStrategyFactory(aiService *AIService, projectID string) EstimationStrategy {
	if aiService != nil && aiService.IsConfiguredForProject(projectID) && IsFeatureEnabled("time_estimation") {
		return NewAIEstimationStrategy(aiService)
	}
	return NewHistoricalEstimationStrategy()
}
