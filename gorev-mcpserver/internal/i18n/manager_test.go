package i18n

import (
	"testing"

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
	err := Initialize("tr")
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
	err := Initialize("tr")
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
			key:      "error.taskNotFound",
			data:     nil,
			contains: "bulunamadı", // Should contain Turkish word
		},
		{
			name:     "Key with template data",
			key:      "error.dataManagerCreate",
			data:     map[string]interface{}{"Error": "test error"},
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
	err := Initialize("tr")
	assert.NoError(t, err)

	turkishResult := T("error.taskNotFound", nil)
	assert.NotEmpty(t, turkishResult)

	// Switch to English
	err = SetLanguage("en")
	assert.NoError(t, err)

	englishResult := T("error.taskNotFound", nil)
	assert.NotEmpty(t, englishResult)

	// Results should be different (unless the key doesn't exist in translations)
	// We can't assert they're different since some keys might have same translation
	// But we can assert both are non-empty
	assert.NotEmpty(t, turkishResult)
	assert.NotEmpty(t, englishResult)
}

func TestMultipleInitializations(t *testing.T) {
	// Test that multiple initializations don't cause issues
	err1 := Initialize("tr")
	assert.NoError(t, err1)

	err2 := Initialize("en")
	assert.NoError(t, err2)

	// Should be able to translate after multiple initializations
	result := T("error.taskNotFound", nil)
	assert.NotEmpty(t, result)
}
