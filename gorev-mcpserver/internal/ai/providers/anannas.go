// Package providers implements AI provider integrations
// This file contains the Anannas provider implementation
package providers

import (
	"context"
	"fmt"
)

// AnannasProvider implements the Provider interface for Anannas
// Uses the shared OpenAICompatibleClient (DRY compliance)
type AnannasProvider struct {
	*BaseProvider
}

// NewAnannasProvider creates a new Anannas provider
func NewAnannasProvider(apiKey string) *AnannasProvider {
	return &AnannasProvider{
		BaseProvider: NewBaseProvider(ProviderAnannas, apiKey),
	}
}

// Name returns the provider name
func (p *AnannasProvider) Name() string {
	return "anannas"
}

// Type returns the provider type
func (p *AnannasProvider) Type() ProviderType {
	return ProviderAnannas
}

// ChatCompletion sends a chat completion request to Anannas
// Adds Anannas-specific handling
func (p *AnannasProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
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

// Supports checks if Anannas supports a specific feature
func (p *AnannasProvider) Supports(feature string) bool {
	// Anannas supports all features through its model routing
	supportedFeatures := map[string]bool{
		"chat":               true,
		"task_creation":      true,
		"task_decomposition": true,
		"semantic_search":    true,
		"time_estimation":    true,
		"project_analytics":  true,
		"streaming":          true,
	}

	return supportedFeatures[feature]
}

// GetRecommendedModel returns a recommended model for a given operation
func (p *AnannasProvider) GetRecommendedModel(operation string) string {
	// Map operations to recommended models
	// Anannas uses provider/model format
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

// ValidateAPIKey checks if the API key is valid for Anannas
func (p *AnannasProvider) ValidateAPIKey(ctx context.Context) error {
	// Try to list models as a validation
	_, err := p.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("invalid Anannas API key: %w", err)
	}
	return nil
}
