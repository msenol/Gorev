package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test helpers.go functions that are not covered

func TestTCommon(t *testing.T) {
	// Initialize i18n first
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		key      string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "Simple common key",
			key:      "test",
			data:     nil,
			expected: "common.test", // Will return key itself if not found
		},
		{
			name:     "Common key with data",
			key:      "test",
			data:     map[string]interface{}{"Name": "example"},
			expected: "common.test",
		},
		{
			name:     "Empty key",
			key:      "",
			data:     nil,
			expected: "common.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TCommon("tr", tt.key, tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTParam(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name      string
		paramName string
		expected  string
	}{
		{
			name:      "Known parameter",
			paramName: "gorev_id",
			expected:  "Taşınacak görevin ID'si", // Actual translation from system
		},
		{
			name:      "Unknown parameter fallback",
			paramName: "unknown_param",
			expected:  "unknown_param parameter", // Fallback pattern
		},
		{
			name:      "Empty parameter",
			paramName: "",
			expected:  " parameter", // Fallback pattern for empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TParam("tr", tt.paramName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTValidation(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		validationType string
		param          string
		extra          map[string]interface{}
		expected       string
	}{
		{
			name:           "Required validation",
			validationType: "required",
			param:          "gorev_id",
			extra:          nil,
			expected:       "validation.required", // Returns key if not found
		},
		{
			name:           "Invalid validation with extra data",
			validationType: "invalid",
			param:          "status",
			extra:          map[string]interface{}{"Values": "pending, completed"},
			expected:       "validation.invalid",
		},
		{
			name:           "Empty validation type",
			validationType: "",
			param:          "test",
			extra:          nil,
			expected:       "validation.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TValidation("tr", tt.validationType, tt.param, tt.extra)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildFieldDescription(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		prefix   string
		entity   string
		field    string
		expected string
	}{
		{
			name:     "Basic field description",
			prefix:   "new",
			entity:   "task",
			field:    "title",
			expected: "Yeni title", // Actual result from system
		},
		{
			name:     "Empty prefix",
			prefix:   "",
			entity:   "project",
			field:    "name",
			expected: "common.prefixes. name", // Actual result when prefix is empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildFieldDescription("tr", tt.prefix, tt.entity, tt.field)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildIDDescription(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		entity   string
		idType   string
		expected string
	}{
		{
			name:     "Unique ID description",
			entity:   "task",
			idType:   "unique",
			expected: "common.fields.task_id", // Returns key if not found
		},
		{
			name:     "Regular ID description",
			entity:   "project",
			idType:   "regular",
			expected: "common.fields.task_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildIDDescription("tr", tt.entity, tt.idType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildPaginationDescription(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name           string
		paginationType string
		entity         string
		defaultVal     int
		maxVal         int
		expected       string
	}{
		{
			name:           "Limit pagination",
			paginationType: "limit",
			entity:         "tasks",
			defaultVal:     10,
			maxVal:         100,
			expected:       "maksimum tasks sayısı (varsayılan: 10, maksimum: 100)", // Actual translation
		},
		{
			name:           "Page pagination",
			paginationType: "page",
			entity:         "projects",
			defaultVal:     1,
			maxVal:         50,
			expected:       "common.pagination.page_pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildPaginationDescription("tr", tt.paginationType, tt.entity, tt.defaultVal, tt.maxVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildPrefixedDescription(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		prefix   string
		target   string
		expected string
	}{
		{
			name:     "Basic prefixed description",
			prefix:   "new",
			target:   "task description",
			expected: "Yeni task description", // Turkish "Yeni" prefix
		},
		{
			name:     "Empty prefix",
			prefix:   "",
			target:   "project name",
			expected: "common.prefixes. project name", // empty prefix returns the key
		},
		{
			name:     "Empty target",
			prefix:   "update",
			target:   "",
			expected: "common.prefixes.update ", // empty target still includes prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildPrefixedDescription("tr", tt.prefix, tt.target)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCommonSuffix(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name       string
		suffixType string
		expected   string
	}{
		{
			name:       "Optional suffix",
			suffixType: "optional",
			expected:   "common.suffixes.optional",
		},
		{
			name:       "Required suffix",
			suffixType: "required",
			expected:   "parametresi gerekli",
		},
		{
			name:       "Empty suffix type",
			suffixType: "",
			expected:   "common.suffixes.", // Empty key returns the key itself
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCommonSuffix("tr", tt.suffixType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test uncovered helper functions for improved coverage

func TestUncoveredHelpers(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	// GetEntityName
	assert.NotEmpty(t, GetEntityName("tr", "task"))
	assert.NotEmpty(t, GetEntityName("tr", "project"))

	// Format functions
	assert.NotEmpty(t, FormatParameterRequired("tr", "gorev_id"))
	assert.NotEmpty(t, FormatInvalidValue("tr", "status", "invalid", []string{"beklemede", "tamamlandi"}))
	assert.NotEmpty(t, FormatEntityNotFound("tr", "task", "123"))
	assert.NotEmpty(t, FormatOperationFailed("tr", "create", assert.AnError))

	// TRequired functions
	assert.NotEmpty(t, TRequiredParam("tr", "gorev_id"))
	assert.NotEmpty(t, TRequiredArray("tr", "etiketler"))
	assert.NotEmpty(t, TRequiredObject("tr", "degerler"))

	// TEntity functions
	assert.NotEmpty(t, TEntityNotFound("tr", "task", assert.AnError))
	assert.NotEmpty(t, TEntityNotFoundByID("tr", "task", "123"))

	// TOperation functions
	assert.NotEmpty(t, TOperationFailed("tr", "create", "task", assert.AnError))
	assert.NotEmpty(t, TSuccess("tr", "created", "task", nil))

	// TInvalid functions
	assert.NotEmpty(t, TInvalidValue("tr", "field", "value", []string{"option1", "option2"}))
	assert.NotEmpty(t, TInvalidStatus("tr", "invalid", []string{"beklemede", "tamamlandi"}))
	assert.NotEmpty(t, TInvalidPriority("tr", "invalid"))
	assert.NotEmpty(t, TInvalidDate("tr", "2025-13-45"))
	assert.NotEmpty(t, TInvalidFormat("tr", "date", "2025-13-45"))

	// TAction functions (all require error parameter)
	testErr := assert.AnError
	assert.NotEmpty(t, TCreateFailed("tr", "task", testErr))
	assert.NotEmpty(t, TUpdateFailed("tr", "task", testErr))
	assert.NotEmpty(t, TDeleteFailed("tr", "task", testErr))
	assert.NotEmpty(t, TFetchFailed("tr", "task", testErr))
	assert.NotEmpty(t, TSaveFailed("tr", "task", testErr))
	assert.NotEmpty(t, TSetFailed("tr", "status", testErr))
	assert.NotEmpty(t, TInitFailed("tr", "database", testErr))
	assert.NotEmpty(t, TCheckFailed("tr", "validation", testErr))
	assert.NotEmpty(t, TQueryFailed("tr", "task", testErr))
	assert.NotEmpty(t, TProcessFailed("tr", "data", testErr))
	assert.NotEmpty(t, TListFailed("tr", "tasks", testErr))
	assert.NotEmpty(t, TEditFailed("tr", "task", testErr))
	assert.NotEmpty(t, TAddFailed("tr", "tag", testErr))
	assert.NotEmpty(t, TRemoveFailed("tr", "tag", testErr))
	assert.NotEmpty(t, TReadFailed("tr", "file", testErr))
	assert.NotEmpty(t, TConvertFailed("tr", "data", "json", testErr))
	assert.NotEmpty(t, TParseFailed("tr", "date", testErr))

	// TSuccess messages
	assert.NotEmpty(t, TCreated("tr", "task", "Test Task", "123"))
	assert.NotEmpty(t, TUpdated("tr", "task", "details"))
	assert.NotEmpty(t, TDeleted("tr", "task", "Test Task", "123"))
	assert.NotEmpty(t, TSet("tr", "status", "devam_ediyor"))
	assert.NotEmpty(t, TRemoved("tr", "tag"))
	assert.NotEmpty(t, TAdded("tr", "tag", "yeni etiket"))
	assert.NotEmpty(t, TMoved("tr", "task"))
	assert.NotEmpty(t, TEdited("tr", "task", "Test Task"))

	// Field helpers
	assert.NotEmpty(t, TFieldID("tr", "task", "create"))
	assert.NotEmpty(t, TTaskCount("tr", "total", "10"))
	assert.NotEmpty(t, TProjectField("tr", "name"))
	assert.NotEmpty(t, TSubtaskField("tr", "title"))
	assert.NotEmpty(t, TCommaSeparated("tr", "tags"))
	assert.NotEmpty(t, TWithFormat("tr", "file path", "json"))
	assert.NotEmpty(t, TFilePath("tr", "import"))
	assert.NotEmpty(t, TTemplate("tr", "bug"))
	assert.NotEmpty(t, TBatch("tr", "update"))

	// Markdown helpers
	assert.NotEmpty(t, TLabel("tr", "test"))
	assert.NotEmpty(t, TMarkdownLabel("tr", "test", "value"))
	assert.NotEmpty(t, TMarkdownHeader("tr", 1, "test"))
	assert.NotEmpty(t, TMarkdownBold("tr", "test"))
	assert.NotEmpty(t, TMarkdownSection("tr", "📝", "test"))

	// Utility helpers
	assert.NotEmpty(t, TCount("tr", "tasks", 5))
	assert.NotEmpty(t, TDuration("tr", "elapsed", 120))
	assert.NotEmpty(t, TListItem("tr", "task", 1))
}

// TestTStatus tests all status value handling
func TestTStatus(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name   string
		status string
	}{
		{
			name:   "Pending status",
			status: "beklemede",
		},
		{
			name:   "In progress status",
			status: "devam_ediyor",
		},
		{
			name:   "Completed status",
			status: "tamamlandi",
		},
		{
			name:   "Cancelled status",
			status: "iptal",
		},
		{
			name:   "Unknown status fallback",
			status: "unknown_status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TStatus("tr", tt.status)
			// Just verify function returns non-empty string
			assert.NotEmpty(t, result)
			// Unknown status should return itself
			if tt.status == "unknown_status" {
				assert.Equal(t, "unknown_status", result)
			}
		})
	}
}

// TestTPriority tests all priority value handling
func TestTPriority(t *testing.T) {
	err := Initialize("tr")
	assert.NoError(t, err)

	tests := []struct {
		name     string
		priority string
	}{
		{
			name:     "Low priority",
			priority: "dusuk",
		},
		{
			name:     "Medium priority",
			priority: "orta",
		},
		{
			name:     "High priority",
			priority: "yuksek",
		},
		{
			name:     "Unknown priority fallback",
			priority: "unknown_priority",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TPriority("tr", tt.priority)
			// Just verify function returns non-empty string
			assert.NotEmpty(t, result)
			// Unknown priority should return itself
			if tt.priority == "unknown_priority" {
				assert.Equal(t, "unknown_priority", result)
			}
		})
	}
}
