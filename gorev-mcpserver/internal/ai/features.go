// Package ai provides AI provider integrations for task management
// This file contains feature flags and fallback modes (Rule 15 compliance)
package ai

// FeatureFlag defines whether an AI feature is enabled and how it should fallback
type FeatureFlag struct {
	// Enabled indicates if the feature is available
	Enabled bool
	// Required indicates if the feature must work (fail if unavailable)
	// If false, gracefully degrades to fallback mode
	Required bool
	// FallbackMode specifies behavior when AI is unavailable:
	// - "rule_based": Use existing rule-based logic
	// - "cached": Use cached results
	// - "disabled": Feature is unavailable (returns error/not implemented)
	// - "fts": Fall back to full-text search
	// - "historical": Fall back to historical data
	// - "nlp_processor": Fall back to existing NLP processor
	FallbackMode string
}

// AIFeatures defines all AI-powered features and their fallback behavior
var AIFeatures = map[string]FeatureFlag{
	// Task creation enhancement - AI suggests title, description, priority, tags
	// Falls back to template-only creation (existing behavior)
	"smart_task_creation": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "rule_based",
	},

	// Task decomposition - AI breaks complex tasks into subtasks
	// Feature unavailable if AI not configured (new capability)
	"task_decomposition": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "disabled",
	},

	// Semantic search - AI understands meaning and context
	// Falls back to FTS (existing behavior)
	"semantic_search": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "fts",
	},

	// Time estimation - AI predicts task duration
	// Falls back to historical data (existing behavior)
	"time_estimation": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "historical",
	},

	// File association - AI suggests task-file relationships
	// Falls back to manual file association (existing behavior)
	"file_association": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "rule_based",
	},

	// Project analytics - AI provides insights and risk assessment
	// Feature unavailable if AI not configured (new capability)
	"project_analytics": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "disabled",
	},

	// Prioritization - AI suggests task ordering
	// Falls back to rule-based sorting (existing behavior)
	"prioritization": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "rule_based",
	},

	// Natural language commands - AI understands complex queries
	// Falls back to existing NLP processor
	"nl_commands": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "nlp_processor",
	},

	// Chat - Interactive AI assistant
	// Feature unavailable if AI not configured (new capability)
	"chat": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "disabled",
	},

	// Suggestions - AI-powered task suggestions
	// Falls back to rule-based suggestions (existing behavior)
	"suggest": {
		Enabled:      true,
		Required:     false,
		FallbackMode: "rule_based",
	},
}

// IsFeatureEnabled checks if a feature is enabled
func IsFeatureEnabled(feature string) bool {
	flag, ok := AIFeatures[feature]
	return ok && flag.Enabled
}

// GetFeatureFallback returns the fallback mode for a feature
func GetFeatureFallback(feature string) string {
	flag, ok := AIFeatures[feature]
	if !ok {
		return "disabled"
	}
	return flag.FallbackMode
}

// IsFeatureRequired checks if a feature is required (cannot fallback)
func IsFeatureRequired(feature string) bool {
	flag, ok := AIFeatures[feature]
	return ok && flag.Required
}

// FeatureAvailability represents the availability state of a feature
type FeatureAvailability struct {
	Feature      string
	Available    bool   // True if AI is configured
	FallbackMode string // What happens when AI unavailable
	CanUse       bool   // True if either available or has viable fallback
}

// CheckFeature returns the availability state for a feature
func CheckFeature(feature string, aiConfigured bool) FeatureAvailability {
	flag, ok := AIFeatures[feature]
	if !ok {
		return FeatureAvailability{
			Feature:      feature,
			Available:    false,
			FallbackMode: "disabled",
			CanUse:       false,
		}
	}

	return FeatureAvailability{
		Feature:      feature,
		Available:    aiConfigured && flag.Enabled,
		FallbackMode: flag.FallbackMode,
		CanUse:       (aiConfigured && flag.Enabled) || flag.FallbackMode != "disabled",
	}
}

// AIOperation represents the type of AI operation being performed
type AIOperation string

const (
	AIOperationTaskCreation     AIOperation = "task_creation"
	AIOperationTaskDecomposition AIOperation = "task_decomposition"
	AIOperationSemanticSearch    AIOperation = "semantic_search"
	AIOperationTimeEstimation    AIOperation = "time_estimation"
	AIOperationFileAssociation   AIOperation = "file_association"
	AIOperationProjectAnalytics  AIOperation = "project_analytics"
	AIOperationPrioritization    AIOperation = "prioritization"
	AIOperationNLCommand         AIOperation = "nl_command"
	AIOperationChat              AIOperation = "chat"
	AIOperationSuggest           AIOperation = "suggest"
)

// String returns the string representation of the operation
func (o AIOperation) String() string {
	return string(o)
}

// FeatureForOperation returns the feature flag for an operation
func FeatureForOperation(op AIOperation) string {
	return string(op)
}

// ProviderType represents the AI provider type
type ProviderType string

const (
	ProviderOpenRouter ProviderType = "openrouter"
	ProviderAnannas    ProviderType = "anannas"
)

// IsValidProvider checks if a provider type is valid
func IsValidProvider(provider string) bool {
	return provider == string(ProviderOpenRouter) || provider == string(ProviderAnannas)
}

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
