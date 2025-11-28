package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Manager struct {
	bundle *i18n.Bundle
	// Note: localizer removed - we create per-request localizers for multi-client support
}

var globalManager *Manager

// Initialize sets up the i18n system with the specified language
func Initialize(lang string) error {
	// Try embedded locales first for better reliability
	err := initializeWithEmbedded(lang)
	if err == nil {
		return nil
	}

	// Fallback to original initialization for backward compatibility
	bundle := i18n.NewBundle(language.Turkish) // Default language is Turkish
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	loadSuccess := false

	// Load Turkish translations (default) - using flattened format for go-i18n compatibility
	trPath := getLocaleFilePath("tr_flat.json")
	if _, err := os.Stat(trPath); err == nil {
		_, err = bundle.LoadMessageFile(trPath)
		if err == nil {
			loadSuccess = true
		}
	}

	// Load English translations - using flattened format for go-i18n compatibility
	enPath := getLocaleFilePath("en_flat.json")
	if _, err := os.Stat(enPath); err == nil {
		_, err = bundle.LoadMessageFile(enPath)
		if err == nil {
			loadSuccess = true
		}
	}

	if !loadSuccess {
		return fmt.Errorf("failed to load any translations from filesystem and embedded fallback failed: %w", err)
	}

	// Store bundle only (localizers created per-request)
	globalManager = &Manager{
		bundle: bundle,
	}

	return nil
}

// TWithLang translates a message key with specified language and optional template data
// This is the primary translation function for multi-client support
func TWithLang(lang string, messageID string, templateData ...map[string]interface{}) string {
	if globalManager == nil {
		// Fallback to messageID if i18n is not initialized
		return messageID
	}

	// Create localizer for this specific language
	localizer := createLocalizerForLanguage(globalManager.bundle, lang)

	var data map[string]interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	if err != nil {
		// Return the message ID if translation fails
		return messageID
	}

	return msg
}

// T translates a message key with optional template data
// Backward compatible function - uses default Turkish language
// New code should prefer TWithLang() for explicit language control
func T(messageID string, templateData ...map[string]interface{}) string {
	return TWithLang("tr", messageID, templateData...)
}

// SetLanguage changes the current language
// Deprecated: This function is kept for backward compatibility only
// New code should use TWithLang() for per-request language control
func SetLanguage(lang string) error {
	if globalManager == nil {
		return fmt.Errorf("i18n manager not initialized")
	}
	// Note: No-op in new architecture (per-request localizers)
	// Kept for backward compatibility
	return nil
}

// GetCurrentLanguage returns the default language code
// Deprecated: In multi-client architecture, language is per-request
// This returns the system default only
func GetCurrentLanguage() string {
	return "tr" // System default
}

// createLocalizerForLanguage creates a localizer for the specified language
// This is a helper function used by TWithLang for per-request localization
func createLocalizerForLanguage(bundle *i18n.Bundle, lang string) *i18n.Localizer {
	// Validate and normalize language
	if lang != "en" && lang != "tr" {
		lang = "tr" // Default to Turkish for unsupported languages
	}

	if lang == "en" {
		return i18n.NewLocalizer(bundle, "en", "tr") // English with Turkish fallback
	}
	return i18n.NewLocalizer(bundle, "tr") // Turkish
}

// getLocaleFilePath returns the absolute path to a locale file
func getLocaleFilePath(filename string) string {
	// Try to get the executable directory first
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		localesPath := filepath.Join(execDir, "locales", filename)
		if _, err := os.Stat(localesPath); err == nil {
			return localesPath
		}
	}

	// Try multiple relative paths for different execution contexts
	possiblePaths := []string{
		filepath.Join("locales", filename),             // Direct execution from root
		filepath.Join("..", "..", "locales", filename), // Test execution from internal/mcp/
		filepath.Join("..", "locales", filename),       // Test execution from internal/
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Final fallback to relative path from working directory
	return filepath.Join("locales", filename)
}

// IsInitialized returns true if the i18n system is initialized
func IsInitialized() bool {
	return globalManager != nil
}
