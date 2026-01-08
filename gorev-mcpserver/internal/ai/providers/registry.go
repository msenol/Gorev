// Package providers implements AI provider integrations
// This file contains the provider registry for managing provider instances
package providers

import (
	"context"
	"fmt"
	"sync"
)

// IsValidProvider checks if a provider type is valid
func IsValidProvider(provider string) bool {
	return provider == string(ProviderOpenRouter) || provider == string(ProviderAnannas)
}

// Registry manages provider instances
// Thread-safe singleton pattern (DRY compliance)
type Registry struct {
	mu        sync.RWMutex
	providers map[string]Provider // project_id -> Provider
}

// global registry instance
var globalRegistry *Registry
var once sync.Once

// GetRegistry returns the global provider registry (singleton)
func GetRegistry() *Registry {
	once.Do(func() {
		globalRegistry = &Registry{
			providers: make(map[string]Provider),
		}
	})
	return globalRegistry
}

// Register registers a provider for a project
func (r *Registry) Register(projectID string, provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[projectID] = provider
}

// Unregister removes a provider for a project
func (r *Registry) Unregister(projectID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.providers, projectID)
}

// Get retrieves a provider for a project
func (r *Registry) Get(projectID string) (Provider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, ok := r.providers[projectID]
	return provider, ok
}

// CreateProvider creates a provider instance based on type and API key
// Factory function for provider creation (DRY compliance)
func CreateProvider(providerType ProviderType, apiKey string) (Provider, error) {
	if !IsValidProvider(string(providerType)) {
		return nil, fmt.Errorf("invalid provider type: %s", providerType)
	}

	switch providerType {
	case ProviderOpenRouter:
		return NewOpenRouterProvider(apiKey), nil
	case ProviderAnannas:
		return NewAnannasProvider(apiKey), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", providerType)
	}
}

// RegisterFromConfig registers a provider from configuration
func (r *Registry) RegisterFromConfig(projectID string, providerType ProviderType, apiKey string) error {
	provider, err := CreateProvider(providerType, apiKey)
	if err != nil {
		return fmt.Errorf("create provider: %w", err)
	}

	// Validate API key
	ctx := context.Background()
	if err := provider.HealthCheck(ctx); err != nil {
		return fmt.Errorf("provider health check failed: %w", err)
	}

	r.Register(projectID, provider)
	return nil
}

// ListProjects returns all project IDs with registered providers
func (r *Registry) ListProjects() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projects := make([]string, 0, len(r.providers))
	for projectID := range r.providers {
		projects = append(projects, projectID)
	}
	return projects
}

// Clear removes all registered providers
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers = make(map[string]Provider)
}

// ProviderInfo contains information about a registered provider
type ProviderInfo struct {
	ProjectID string
	Name      string
	Type      ProviderType
	BaseURL   string
}

// GetInfo returns information about a registered provider
func (r *Registry) GetInfo(projectID string) (*ProviderInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, ok := r.providers[projectID]
	if !ok {
		return nil, false
	}

	return &ProviderInfo{
		ProjectID: projectID,
		Name:      provider.Name(),
		Type:      provider.Type(),
		BaseURL:   provider.BaseURL(),
	}, true
}

// ListInfo returns information about all registered providers
func (r *Registry) ListInfo() []ProviderInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info := make([]ProviderInfo, 0, len(r.providers))
	for projectID, provider := range r.providers {
		info = append(info, ProviderInfo{
			ProjectID: projectID,
			Name:      provider.Name(),
			Type:      provider.Type(),
			BaseURL:   provider.BaseURL(),
		})
	}
	return info
}
