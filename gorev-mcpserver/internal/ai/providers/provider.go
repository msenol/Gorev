// Package providers implements AI provider integrations
// This file contains the base provider interface
package providers

import (
	"context"
	"fmt"
)

// ProviderType represents the AI provider type
type ProviderType string

const (
	ProviderOpenRouter ProviderType = "openrouter"
	ProviderAnannas    ProviderType = "anannas"
)

// ProviderBaseURL returns the base URL for a provider
func ProviderBaseURL(provider ProviderType) string {
	switch provider {
	case ProviderOpenRouter:
		return "https://openrouter.ai/api/v1"
	case ProviderAnannas:
		return "https://api.anannas.ai/v1"
	default:
		return ""
	}
}

// Provider defines the interface for AI providers (OpenRouter, Anannas, etc.)
// All providers implement this interface for consistency (DRY principle)
type Provider interface {
	// Name returns the provider name
	Name() string

	// Type returns the provider type
	Type() ProviderType

	// BaseURL returns the base URL of the provider
	BaseURL() string

	// ListModels retrieves available models from the provider
	ListModels(ctx context.Context) ([]ModelInfo, error)

	// ValidateModel checks if a model ID is valid for this provider
	ValidateModel(ctx context.Context, model string) error

	// ChatCompletion sends a chat completion request
	ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// HealthCheck verifies the provider is accessible
	HealthCheck(ctx context.Context) error

	// Supports checks if the provider supports a specific feature
	Supports(feature string) bool
}

// BaseProvider provides common functionality for all providers
// Reduces code duplication through composition (DRY principle)
type BaseProvider struct {
	providerType ProviderType
	baseURL     string
	apiKey      string
	client      *OpenAICompatibleClient
}

// NewBaseProvider creates a new base provider
func NewBaseProvider(providerType ProviderType, apiKey string) *BaseProvider {
	return &BaseProvider{
		providerType: providerType,
		baseURL:     ProviderBaseURL(providerType),
		apiKey:      apiKey,
		client:      NewOpenAICompatibleClient(ProviderBaseURL(providerType), apiKey),
	}
}

// Name returns the provider name
func (p *BaseProvider) Name() string {
	return string(p.providerType)
}

// Type returns the provider type
func (p *BaseProvider) Type() ProviderType {
	return p.providerType
}

// BaseURL returns the base URL of the provider
func (p *BaseProvider) BaseURL() string {
	return p.baseURL
}

// ListModels retrieves available models from the provider
func (p *BaseProvider) ListModels(ctx context.Context) ([]ModelInfo, error) {
	return p.client.ListModels(ctx)
}

// ValidateModel checks if a model ID is valid for this provider
func (p *BaseProvider) ValidateModel(ctx context.Context, model string) error {
	if model == "" {
		return fmt.Errorf("model ID cannot be empty")
	}

	// List all models and check if the requested model exists
	models, err := p.ListModels(ctx)
	if err != nil {
		return fmt.Errorf("list models: %w", err)
	}

	for _, m := range models {
		if m.ID == model {
			return nil
		}
	}

	return fmt.Errorf("model %s not found for provider %s", model, p.Name())
}

// ChatCompletion sends a chat completion request
func (p *BaseProvider) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	return p.client.ChatCompletion(ctx, req)
}

// HealthCheck verifies the provider is accessible
func (p *BaseProvider) HealthCheck(ctx context.Context) error {
	// Try to list models as a health check
	_, err := p.ListModels(ctx)
	return err
}

// Supports checks if the provider supports a specific feature
// Base implementation - can be overridden by specific providers
func (p *BaseProvider) Supports(feature string) bool {
	// By default, all providers support these features
	supportedFeatures := map[string]bool{
		"chat":              true,
		"task_creation":     true,
		"task_decomposition": true,
		"semantic_search":   true,
		"time_estimation":   true,
		"project_analytics": true,
	}

	return supportedFeatures[feature]
}

// SetTimeout sets the HTTP client timeout for the provider
func (p *BaseProvider) SetTimeout(timeout interface{}) {
	// Type assertion for time.Duration
	// This is a placeholder - actual implementation would use proper timeout
}

// ModelCost calculates the cost of a chat completion based on token usage
func (p *BaseProvider) ModelCost(modelID string, usage ChatUsage) (float64, error) {
	// Get model info to find pricing
	ctx := context.Background()
	models, err := p.ListModels(ctx)
	if err != nil {
		return 0, fmt.Errorf("get model info: %w", err)
	}

	var pricing Pricing
	for _, m := range models {
		if m.ID == modelID {
			pricing = m.Pricing
			break
		}
	}

	// Calculate cost
	inputCost := float64(usage.PromptTokens) / 1000000 * pricing.Input
	outputCost := float64(usage.CompletionTokens) / 1000000 * pricing.Output
	totalCost := inputCost + outputCost

	return totalCost, nil
}
