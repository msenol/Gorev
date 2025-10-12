package context

import (
	"context"
	"os"
)

type contextKey string

const languageKey contextKey = "language"

// WithLanguage adds language to context
func WithLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, languageKey, lang)
}

// GetLanguage extracts language from context with fallback priority:
// 1. Context value (from explicit WithLanguage)
// 2. Environment variable GOREV_LANG
// 3. Default "tr"
func GetLanguage(ctx context.Context) string {
	// Priority 1: Context value
	if lang, ok := ctx.Value(languageKey).(string); ok && lang != "" {
		return lang
	}

	// Priority 2: Environment variable
	if lang := os.Getenv("GOREV_LANG"); lang != "" {
		return lang
	}

	// Priority 3: Default
	return "tr"
}

// ValidateLanguage ensures language is supported
func ValidateLanguage(lang string) string {
	if lang == "en" || lang == "tr" {
		return lang
	}
	return "tr" // Default fallback
}
