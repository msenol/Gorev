package i18n

import (
	"testing"

	"github.com/msenol/gorev/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name     string
		language string
		wantErr  bool
	}{
		{
			name:     "Turkish initialization",
			language: "tr",
			wantErr:  false,
		},
		{
			name:     "English initialization",
			language: "en",
			wantErr:  false,
		},
		{
			name:     "Invalid language fallback to Turkish",
			language: "invalid",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.language)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSetLanguage(t *testing.T) {
	// Initialize first
	err := Initialize(constants.DefaultTestLanguage)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		language string
		wantErr  bool
	}{
		{
			name:     "Switch to English",
			language: "en",
			wantErr:  false,
		},
		{
			name:     "Switch to Turkish",
			language: "tr",
			wantErr:  false,
		},
		{
			name:     "Invalid language",
			language: "invalid",
			wantErr:  false, // Should not error, just fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetLanguage(tt.language)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestT(t *testing.T) {
	// Initialize with Turkish
	err := Initialize(constants.DefaultTestLanguage)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		key      string
		data     map[string]interface{}
		expected string
		contains string
	}{
		{
			name:     "Simple key translation",
			key:      "common.validation.not_found",
			data:     nil,
			contains: "bulunamadı", // Should contain Turkish word
		},
		{
			name:     "Key with template data",
			key:      "common.operations.create_failed",
			data:     map[string]interface{}{"Entity": "data_manager", "Error": "test error"},
			contains: "test error",
		},
		{
			name:     "Non-existent key returns key itself",
			key:      "non.existent.key",
			data:     nil,
			expected: "non.existent.key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := T(tt.key, tt.data)

			if tt.expected != "" {
				assert.Equal(t, tt.expected, result)
			} else if tt.contains != "" {
				assert.Contains(t, result, tt.contains)
			}

			// Result should never be empty
			assert.NotEmpty(t, result)
		})
	}
}

func TestTWithLanguageSwitching(t *testing.T) {
	// Start with Turkish
	err := Initialize(constants.DefaultTestLanguage)
	assert.NoError(t, err)

	turkishResult := T("common.validation.not_found", map[string]interface{}{"Entity": "görev", "Error": "test"})
	assert.NotEmpty(t, turkishResult)

	// Switch to English
	err = SetLanguage("en")
	assert.NoError(t, err)

	englishResult := T("common.validation.not_found", map[string]interface{}{"Entity": "görev", "Error": "test"})
	assert.NotEmpty(t, englishResult)

	// Results should be different (unless the key doesn't exist in translations)
	// We can't assert they're different since some keys might have same translation
	// But we can assert both are non-empty
	assert.NotEmpty(t, turkishResult)
	assert.NotEmpty(t, englishResult)
}

func TestMultipleInitializations(t *testing.T) {
	// Test that multiple initializations don't cause issues
	err1 := Initialize(constants.DefaultTestLanguage)
	assert.NoError(t, err1)

	err2 := Initialize("en")
	assert.NoError(t, err2)

	// Should be able to translate after multiple initializations
	result := T("common.validation.not_found", map[string]interface{}{"Entity": "görev", "Error": "test"})
	assert.NotEmpty(t, result)
}

// Test additional manager.go functions not covered by existing tests

func TestGetCurrentLanguage(t *testing.T) {
	tests := []struct {
		name     string
		initLang string
		setLang  string
		expected string
	}{
		{
			name:     "Turkish language detection",
			initLang: "tr",
			setLang:  "",
			expected: "tr",
		},
		{
			name:     "English language detection",
			initLang: "en",
			setLang:  "",
			expected: "en", // English locale data has "lang.code": "en"
		},
		{
			name:     "Language switch detection",
			initLang: "tr",
			setLang:  "en",
			expected: "en", // After switching to English
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.initLang)
			assert.NoError(t, err)

			if tt.setLang != "" {
				err = SetLanguage(tt.setLang)
				assert.NoError(t, err)
			}

			result := GetCurrentLanguage()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsInitialized(t *testing.T) {
	// Reset global manager first (simulate uninitialized state)
	originalManager := globalManager
	globalManager = nil

	// Test uninitialized state
	assert.False(t, IsInitialized())

	// Test initialized state
	err := Initialize("tr")
	assert.NoError(t, err)
	assert.True(t, IsInitialized())

	// Restore original state
	globalManager = originalManager
}

func TestTWithUninitializedManager(t *testing.T) {
	// Reset global manager to simulate uninitialized state
	originalManager := globalManager
	globalManager = nil
	defer func() { globalManager = originalManager }()

	// Should return the key itself when manager is not initialized
	result := T("test.key", map[string]interface{}{"data": "value"})
	assert.Equal(t, "test.key", result)
}

func TestSetLanguageWithUninitializedManager(t *testing.T) {
	// Reset global manager to simulate uninitialized state
	originalManager := globalManager
	globalManager = nil
	defer func() { globalManager = originalManager }()

	// Should return error when manager is not initialized
	err := SetLanguage("en")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "i18n manager not initialized")
}

func TestGetCurrentLanguageWithUninitializedManager(t *testing.T) {
	// Reset global manager to simulate uninitialized state
	originalManager := globalManager
	globalManager = nil
	defer func() { globalManager = originalManager }()

	// Should return default "tr" when manager is not initialized
	result := GetCurrentLanguage()
	assert.Equal(t, "tr", result)
}

func TestInitializeWithEmbeddedFallback(t *testing.T) {
	// Test initialization with embedded locales
	// This tests the initializeWithEmbedded function indirectly
	tests := []struct {
		name string
		lang string
	}{
		{
			name: "Initialize with Turkish embedded",
			lang: "tr",
		},
		{
			name: "Initialize with English embedded",
			lang: "en",
		},
		{
			name: "Initialize with invalid lang (fallback to Turkish)",
			lang: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.lang)
			assert.NoError(t, err)
			assert.True(t, IsInitialized())

			// Test that we can translate basic keys
			result := T("error.noArguments", nil)
			assert.NotEmpty(t, result)
		})
	}
}

func TestEmbeddedLocaleData(t *testing.T) {
	// Test that embedded locale data is valid JSON and contains expected keys
	err := Initialize("tr")
	assert.NoError(t, err)

	// Test Turkish embedded keys - using keys that exist in mock
	trResult := T("tools.params.descriptions.gorev_id", nil)
	assert.NotEmpty(t, trResult)
	assert.NotEqual(t, "tools.params.descriptions.gorev_id", trResult) // Should be translated
	assert.Contains(t, trResult, "ID")                                 // Translation contains ID

	// Switch to English and test
	err = SetLanguage("en")
	assert.NoError(t, err)

	enResult := T("tools.params.descriptions.gorev_id", nil)
	assert.NotEmpty(t, enResult)
	assert.NotEqual(t, "tools.params.descriptions.gorev_id", enResult) // Should be translated
	assert.Contains(t, enResult, "ID")                                 // English translation contains ID
}

func TestHasKey(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	// Test the HasKey function by checking actual embedded keys
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{
			name:     "Known parameter key",
			key:      "tools.params.descriptions.gorev_id",
			expected: true, // This is in the HasKey mock
		},
		{
			name:     "Another known key",
			key:      "tools.params.descriptions.template_id",
			expected: true, // This should also be in the mock
		},
		{
			name:     "Non-existing key",
			key:      "non.existent.key",
			expected: false,
		},
		{
			name:     "Empty key",
			key:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasKey(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}
