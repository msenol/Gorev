package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var embeddedLocales embed.FS

// initializeWithEmbedded initializes i18n with embedded locales as fallback
func initializeWithEmbedded(lang string) error {
	bundle := i18n.NewBundle(language.Turkish) // Default language is Turkish
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	loadSuccess := false
	var lastErr error

	// Try to load from filesystem first (for development) - using flattened format
	trPath := getLocaleFilePath("tr_flat.json")
	if _, err := os.Stat(trPath); err == nil {
		if _, err := bundle.LoadMessageFile(trPath); err == nil {
			loadSuccess = true
		} else {
			lastErr = err
		}
	}

	enPath := getLocaleFilePath("en_flat.json")
	if _, err := os.Stat(enPath); err == nil {
		if _, err := bundle.LoadMessageFile(enPath); err == nil {
			loadSuccess = true
		} else {
			lastErr = err
		}
	}

	// If filesystem loading failed, use embedded data
	if !loadSuccess {
		// Load Turkish translations from embedded FS - using flattened format
		trData, err := embeddedLocales.ReadFile("locales/tr_flat.json")
		if err == nil {
			if _, err := bundle.ParseMessageFileBytes(trData, "tr_flat.json"); err == nil {
				loadSuccess = true
			} else {
				lastErr = err
			}
		} else {
			lastErr = err
		}

		// Load English translations from embedded FS - using flattened format
		enData, err := embeddedLocales.ReadFile("locales/en_flat.json")
		if err == nil {
			if _, err := bundle.ParseMessageFileBytes(enData, "en_flat.json"); err == nil {
				loadSuccess = true
			} else {
				lastErr = err
			}
		} else {
			lastErr = err
		}
	}

	if !loadSuccess {
		if lastErr != nil {
			return fmt.Errorf("failed to load any locale files (neither filesystem nor embedded), last error: %w", lastErr)
		}
		return fmt.Errorf("failed to load any locale files (neither filesystem nor embedded)")
	}

	// Store bundle only (localizers created per-request)
	globalManager = &Manager{
		bundle: bundle,
	}

	return nil
}
