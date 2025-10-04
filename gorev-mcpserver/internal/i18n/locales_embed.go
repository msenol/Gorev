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

	// Try to load from filesystem first (for development)
	trPath := getLocaleFilePath("tr.json")
	if _, err := os.Stat(trPath); err == nil {
		if _, err := bundle.LoadMessageFile(trPath); err == nil {
			loadSuccess = true
		}
	}

	enPath := getLocaleFilePath("en.json")
	if _, err := os.Stat(enPath); err == nil {
		if _, err := bundle.LoadMessageFile(enPath); err == nil {
			loadSuccess = true
		}
	}

	// If filesystem loading failed, use embedded data
	if !loadSuccess {
		// Load Turkish translations from embedded FS
		trData, err := embeddedLocales.ReadFile("locales/tr.json")
		if err == nil {
			if _, err := bundle.ParseMessageFileBytes(trData, "tr.json"); err == nil {
				loadSuccess = true
			}
		}

		// Load English translations from embedded FS
		enData, err := embeddedLocales.ReadFile("locales/en.json")
		if err == nil {
			if _, err := bundle.ParseMessageFileBytes(enData, "en.json"); err == nil {
				loadSuccess = true
			}
		}
	}

	if !loadSuccess {
		return fmt.Errorf("failed to load any locale files (neither filesystem nor embedded)")
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
