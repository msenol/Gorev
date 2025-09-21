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
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
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

	// Load Turkish translations (default)
	trPath := getLocaleFilePath("tr.json")
	if _, err := os.Stat(trPath); err == nil {
		_, err = bundle.LoadMessageFile(trPath)
		if err == nil {
			loadSuccess = true
		}
	}

	// Load English translations
	enPath := getLocaleFilePath("en.json")
	if _, err := os.Stat(enPath); err == nil {
		_, err = bundle.LoadMessageFile(enPath)
		if err == nil {
			loadSuccess = true
		}
	}

	if !loadSuccess {
		return fmt.Errorf("failed to load any translations from filesystem and embedded fallback failed: %w", err)
	}

	// Create localizer for the specified language with Turkish fallback
	var localizer *i18n.Localizer
	if lang == "en" {
		localizer = i18n.NewLocalizer(bundle, "en", "tr")
	} else {
		// Default to Turkish for any other language or empty lang
		localizer = i18n.NewLocalizer(bundle, "tr")
	}

	globalManager = &Manager{
		bundle:    bundle,
		localizer: localizer,
	}

	return nil
}

// T translates a message key with optional template data
func T(messageID string, templateData ...map[string]interface{}) string {
	if globalManager == nil {
		// Fallback to messageID if i18n is not initialized
		return messageID
	}

	var data map[string]interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}

	msg, err := globalManager.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	if err != nil {
		// Return the message ID if translation fails
		return messageID
	}

	return msg
}

// SetLanguage changes the current language
func SetLanguage(lang string) error {
	if globalManager == nil {
		return fmt.Errorf("i18n manager not initialized")
	}

	var localizer *i18n.Localizer
	if lang == "en" {
		localizer = i18n.NewLocalizer(globalManager.bundle, "en", "tr")
	} else {
		localizer = i18n.NewLocalizer(globalManager.bundle, "tr")
	}

	globalManager.localizer = localizer
	return nil
}

// GetCurrentLanguage returns the current language code
func GetCurrentLanguage() string {
	if globalManager == nil {
		return "tr"
	}

	// Check if we have English as primary language
	msg, err := globalManager.localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "lang.code",
	})

	if err != nil || msg == "lang.code" {
		return "tr" // Default fallback
	}

	return msg
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
