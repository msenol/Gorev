// Package providers implements AI provider integrations
// This file contains the OpenRouter provider implementation
package providers

import (
	"context"
	"fmt"
)

// OpenRouterProvider implements the Provider interface for OpenRouter
// Uses the shared OpenAICompatibleClient (DRY compliance)
type OpenRouterProvider struct {
	*BaseProvider
}

// NewOpenRouterProvider creates a new OpenRouter provider
func NewOpenRouterProvider(apiKey string) *OpenRouterProvider {
	return &OpenRouterProvider{
		BaseProvider: NewBaseProvider(ProviderOpenRouter, apiKey),
	}
}

// Name returns the provider name
func (p *OpenRouterProvider) Name() string {
	return "openrouter"
}

// Type returns the provider type
func (p *OpenRouterProvider) Type() ProviderType {
	return ProviderOpenRouter
}

// ChatCompletion sends a chat completion request to OpenRouter
// Adds OpenRouter-specific headers and handling
func (p *OpenRouterProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Validate request
	if req.Model == "" {
		return nil, fmt.Errorf("model is required")
	}
	if len(req.Messages) == 0 {
		return nil, fmt.Errorf("at least one message is required")
	}

	// Call base implementation (uses shared HTTP client)
	return p.BaseProvider.ChatCompletion(ctx, req)
}

// Supports checks if OpenRouter supports a specific feature
func (p *OpenRouterProvider) Supports(feature string) bool {
	// OpenRouter supports all features through its model routing
	supportedFeatures := map[string]bool{
		"chat":               true,
		"task_creation":      true,
		"task_decomposition": true,
		"semantic_search":    true,
		"time_estimation":    true,
		"project_analytics":  true,
		"streaming":          true,
		"images":             true,
	}

	return supportedFeatures[feature]
}

// GetRecommendedModel returns a recommended model for a given operation
func (p *OpenRouterProvider) GetRecommendedModel(operation string) string {
	// Map operations to recommended models
	recommendations := map[string]string{
		"task_creation":      "openai/gpt-4o-mini",
		"task_decomposition": "openai/gpt-4o-mini",
		"semantic_search":    "openai/gpt-4o-mini",
		"time_estimation":    "openai/gpt-4o-mini",
		"project_analytics":  "openai/gpt-4o",
		"prioritization":     "openai/gpt-4o-mini",
		"nl_command":         "openai/gpt-4o-mini",
		"chat":               "openai/gpt-4o-mini",
		"suggest":            "openai/gpt-4o-mini",
	}

	if model, ok := recommendations[operation]; ok {
		return model
	}

	// Default to a cost-effective model
	return "openai/gpt-4o-mini"
}

// ValidateAPIKey checks if the API key is valid for OpenRouter
func (p *OpenRouterProvider) ValidateAPIKey(ctx context.Context) error {
	// Try to list models as a validation
	_, err := p.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("invalid OpenRouter API key: %w", err)
	}
	return nil
}
