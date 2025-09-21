package i18n

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Embedded locale data as strings to avoid path issues
var embeddedTRData = `{
  "lang": {
    "code": "tr"
  },
  "messages": {
    "noProjectTasks": "*Bu projeye ait görev bulunmuyor.*",
    "noProjects": "Henüz proje bulunmuyor.",
    "noTasks": "Görev bulunmuyor",
    "noActiveProject": "Henüz aktif proje ayarlanmamış",
    "noTemplates": "Henüz template bulunmuyor",
    "noTasksInProject": "{{.Project}} projesinde görev bulunmuyor",
    "taskListCount": "Toplam {{.Count}} görev",
    "projectHeader": "=== {{.Name}} ===",
    "sizeWarning": "⚠️ {{.Count}} görev daha var (sayfalama ile sınırlandı)"
  },
  "error": {
    "activeProjectRetrieve": "aktif proje getirilemedi: {{.Error}}",
    "noArguments": "Parametre belirtilmedi",
    "parameterRequired": "{{.Param}} parametresi zorunludur"
  }
}`

var embeddedENData = `{
  "lang": {
    "code": "en"
  },
  "messages": {
    "noProjectTasks": "*No tasks found for this project.*",
    "noProjects": "No projects found yet.",
    "noTasks": "No tasks found",
    "noActiveProject": "No active project set yet",
    "noTemplates": "No templates found yet",
    "noTasksInProject": "No tasks in {{.Project}} project",
    "taskListCount": "Total {{.Count}} tasks",
    "projectHeader": "=== {{.Name}} ===",
    "sizeWarning": "⚠️ {{.Count}} more tasks available (limited by pagination)"
  },
  "error": {
    "activeProjectRetrieve": "Failed to retrieve active project: {{.Error}}",
    "noArguments": "No arguments provided",
    "parameterRequired": "Parameter {{.Param}} is required"
  }
}`

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
		// Load Turkish translations from embedded data
		if _, err := bundle.ParseMessageFileBytes([]byte(embeddedTRData), "tr.json"); err == nil {
			loadSuccess = true
		}

		// Load English translations from embedded data
		if _, err := bundle.ParseMessageFileBytes([]byte(embeddedENData), "en.json"); err == nil {
			loadSuccess = true
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
