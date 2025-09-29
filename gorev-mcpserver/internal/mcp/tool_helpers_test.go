package mcp

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParameterValidator tests the ParameterValidator helper functions
func TestParameterValidator(t *testing.T) {
	// Setup i18n for testing
	err := i18n.Initialize("tr")
	require.NoError(t, err)

	validator := NewParameterValidator()

	t.Run("ValidateRequiredString - success", func(t *testing.T) {
		params := map[string]interface{}{
			"id": "test-task-id",
		}

		result, toolResult := validator.ValidateRequiredString(params, "id")
		require.NoError(t, err)
		assert.Nil(t, toolResult)
		assert.Equal(t, "test-task-id", result)
	})

	t.Run("ValidateRequiredString - missing", func(t *testing.T) {
		params := map[string]interface{}{
			"title": "Test Task",
		}

		result, toolResult := validator.ValidateRequiredString(params, "id")
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateRequiredString - empty", func(t *testing.T) {
		params := map[string]interface{}{
			"id": "",
		}

		result, toolResult := validator.ValidateRequiredString(params, "id")
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateRequiredString - whitespace only", func(t *testing.T) {
		params := map[string]interface{}{
			"id": "   ",
		}

		result, toolResult := validator.ValidateRequiredString(params, "id")
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateOptionalString - present", func(t *testing.T) {
		params := map[string]interface{}{
			"description": "Task description",
		}

		result := validator.ValidateOptionalString(params, "description")
		assert.Equal(t, "Task description", result)
	})

	t.Run("ValidateOptionalString - missing", func(t *testing.T) {
		params := map[string]interface{}{
			"title": "Test Task",
		}

		result := validator.ValidateOptionalString(params, "description")
		assert.Equal(t, "", result)
	})

	t.Run("ValidateOptionalString - empty", func(t *testing.T) {
		params := map[string]interface{}{
			"description": "",
		}

		result := validator.ValidateOptionalString(params, "description")
		assert.Equal(t, "", result)
	})

	t.Run("ValidateEnum - required - valid", func(t *testing.T) {
		params := map[string]interface{}{
			"durum": "devam_ediyor",
		}

		result, toolResult := validator.ValidateEnum(params, "durum", constants.GetValidTaskStatuses(), true)
		require.NoError(t, err)
		assert.Nil(t, toolResult)
		assert.Equal(t, "devam_ediyor", result)
	})

	t.Run("ValidateEnum - required - invalid", func(t *testing.T) {
		params := map[string]interface{}{
			"durum": "invalid_status",
		}

		result, toolResult := validator.ValidateEnum(params, "durum", constants.GetValidTaskStatuses(), true)
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateEnum - required - missing", func(t *testing.T) {
		params := map[string]interface{}{
			"title": "Test Task",
		}

		result, toolResult := validator.ValidateEnum(params, "durum", constants.GetValidTaskStatuses(), true)
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateEnum - optional - missing", func(t *testing.T) {
		params := map[string]interface{}{
			"title": "Test Task",
		}

		result, toolResult := validator.ValidateEnum(params, "durum", constants.GetValidTaskStatuses(), false)
		require.NoError(t, err)
		assert.Nil(t, toolResult)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateNumber - valid", func(t *testing.T) {
		params := map[string]interface{}{
			"limit": 25.0,
		}

		result := validator.ValidateNumber(params, "limit", 50)
		assert.Equal(t, 25, result)
	})

	t.Run("ValidateNumber - missing", func(t *testing.T) {
		params := map[string]interface{}{
			"title": "Test Task",
		}

		result := validator.ValidateNumber(params, "limit", 50)
		assert.Equal(t, 50, result)
	})

	t.Run("ValidateNumber - zero", func(t *testing.T) {
		params := map[string]interface{}{
			"limit": 0.0,
		}

		result := validator.ValidateNumber(params, "limit", 50)
		assert.Equal(t, 50, result)
	})

	t.Run("ValidateNumber - negative", func(t *testing.T) {
		params := map[string]interface{}{
			"limit": -5.0,
		}

		result := validator.ValidateNumber(params, "limit", 50)
		assert.Equal(t, 50, result)
	})

	t.Run("ValidateBool - true", func(t *testing.T) {
		params := map[string]interface{}{
			"tum_projeler": true,
		}

		result := validator.ValidateBool(params, "tum_projeler")
		assert.True(t, result)
	})

	t.Run("ValidateBool - false", func(t *testing.T) {
		params := map[string]interface{}{
			"tum_projeler": false,
		}

		result := validator.ValidateBool(params, "tum_projeler")
		assert.False(t, result)
	})

	t.Run("ValidateBool - missing", func(t *testing.T) {
		params := map[string]interface{}{
			"title": "Test Task",
		}

		result := validator.ValidateBool(params, "tum_projeler")
		assert.False(t, result)
	})
}

// TestTaskFormatter tests the TaskFormatter helper functions
func TestTaskFormatter(t *testing.T) {
	formatter := NewTaskFormatter()

	t.Run("FormatTaskBasic", func(t *testing.T) {
		result := formatter.FormatTaskBasic("Test Task", "550e8400-e29b-41d4-a716-446655440000")
		assert.Equal(t, "**Test Task** (ID: 550e8400)", result)
	})

	t.Run("FormatTaskBasic - short ID", func(t *testing.T) {
		result := formatter.FormatTaskBasic("Test Task", "short-id")
		assert.Equal(t, "**Test Task** (ID: short-id)", result)
	})

	t.Run("FormatTaskWithStatus", func(t *testing.T) {
		result := formatter.FormatTaskWithStatus("Test Task", "550e8400-e29b-41d4-a716-446655440000", "devam_ediyor")
		assert.Equal(t, "🔄 Test Task (ID: 550e8400)", result)
	})

	t.Run("FormatTaskWithStatus - completed", func(t *testing.T) {
		result := formatter.FormatTaskWithStatus("Test Task", "550e8400-e29b-41d4-a716-446655440000", "tamamlandi")
		assert.Equal(t, "✅ Test Task (ID: 550e8400)", result)
	})

	t.Run("FormatSuccessMessage", func(t *testing.T) {
		result := formatter.FormatSuccessMessage("Oluşturuldu", "Test Task", "550e8400-e29b-41d4-a716-446655440000")
		assert.Contains(t, result, "✓ Oluşturuldu: Test Task")
		assert.Contains(t, result, "ID: 550e8400")
	})

	t.Run("GetStatusEmoji - valid statuses", func(t *testing.T) {
		assert.Equal(t, constants.EmojiStatusPending, formatter.GetStatusEmoji(constants.TaskStatusPending))
		assert.Equal(t, constants.EmojiStatusInProgress, formatter.GetStatusEmoji(constants.TaskStatusInProgress))
		assert.Equal(t, constants.EmojiStatusCompleted, formatter.GetStatusEmoji(constants.TaskStatusCompleted))
		assert.Equal(t, constants.EmojiStatusCancelled, formatter.GetStatusEmoji(constants.TaskStatusCancelled))
	})

	t.Run("GetStatusEmoji - unknown", func(t *testing.T) {
		assert.Equal(t, constants.EmojiStatusUnknown, formatter.GetStatusEmoji("unknown_status"))
	})

	t.Run("GetPriorityEmoji - valid priorities", func(t *testing.T) {
		assert.Equal(t, constants.EmojiPriorityHigh, formatter.GetPriorityEmoji(constants.PriorityHigh))
		assert.Equal(t, constants.EmojiPriorityMedium, formatter.GetPriorityEmoji(constants.PriorityMedium))
		assert.Equal(t, constants.EmojiPriorityLow, formatter.GetPriorityEmoji(constants.PriorityLow))
	})

	t.Run("GetPriorityEmoji - unknown", func(t *testing.T) {
		assert.Equal(t, constants.EmojiPriorityUnknown, formatter.GetPriorityEmoji("unknown_priority"))
	})
}

// TestErrorFormatter tests the ErrorFormatter helper functions
func TestErrorFormatter(t *testing.T) {
	// Setup i18n for testing
	err := i18n.Initialize("tr")
	require.NoError(t, err)

	formatter := NewErrorFormatter()

	t.Run("FormatNotFoundError", func(t *testing.T) {
		result := formatter.FormatNotFoundError("task", "550e8400-e29b-41d4-a716-446655440000")
		assert.NotNil(t, result)
		assert.IsType(t, &mcp.CallToolResult{}, result)
	})

	t.Run("FormatOperationError", func(t *testing.T) {
		testErr := assert.AnError
		result := formatter.FormatOperationError("gorev_olusturma", testErr)
		assert.NotNil(t, result)
		assert.IsType(t, &mcp.CallToolResult{}, result)
	})

	t.Run("FormatValidationError", func(t *testing.T) {
		result := formatter.FormatValidationError("Invalid task parameters")
		assert.NotNil(t, result)
		assert.IsType(t, &mcp.CallToolResult{}, result)
	})
}

// TestResponseBuilder tests the ResponseBuilder helper functions
func TestResponseBuilder(t *testing.T) {
	// Setup i18n for testing
	err := i18n.Initialize("tr")
	require.NoError(t, err)

	builder := NewResponseBuilder()

	t.Run("NewResponseBuilder", func(t *testing.T) {
		assert.NotNil(t, builder)
		assert.NotNil(t, builder.formatter)
	})

	t.Run("BuildMarkdownTaskDetail", func(t *testing.T) {
		result := builder.BuildMarkdownTaskDetail(nil)
		assert.Equal(t, "# Task Detail\n\nTask details would be formatted here.", result)
	})

	t.Run("BuildTaskList - empty list", func(t *testing.T) {
		result := builder.BuildTaskList([]interface{}{}, "Görev Listesi")
		assert.Contains(t, result, "## Görev Listesi")
		assert.Contains(t, result, "no_tasks_found")
	})

	t.Run("BuildTaskList - with tasks", func(t *testing.T) {
		tasks := []interface{}{"task1", "task2"}
		result := builder.BuildTaskList(tasks, "Görev Listesi")
		assert.Contains(t, result, "## Görev Listesi")
		assert.Contains(t, result, "tasks_found_count")
		// The i18n key doesn't exist, so it returns the key name instead of substituted value
		assert.Contains(t, result, "messages.tasks_found_count")
	})
}

// TestCommonValidators tests the CommonValidators helper functions
func TestCommonValidators(t *testing.T) {
	// Setup i18n for testing
	err := i18n.Initialize("tr")
	require.NoError(t, err)

	validators := NewCommonValidators()

	t.Run("NewCommonValidators", func(t *testing.T) {
		assert.NotNil(t, validators)
		assert.NotNil(t, validators.paramValidator)
		assert.NotNil(t, validators.errorFormatter)
	})

	t.Run("ValidateTaskID - valid", func(t *testing.T) {
		params := map[string]interface{}{
			"id": "test-task-id",
		}

		result, toolResult := validators.ValidateTaskID(params)
		assert.Nil(t, toolResult)
		assert.Equal(t, "test-task-id", result)
	})

	t.Run("ValidateTaskID - missing", func(t *testing.T) {
		params := map[string]interface{}{
			"title": "Test Task",
		}

		result, toolResult := validators.ValidateTaskID(params)
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateTaskIDField - valid", func(t *testing.T) {
		params := map[string]interface{}{
			"task_id": "test-task-id",
		}

		result, toolResult := validators.ValidateTaskIDField(params, "task_id")
		assert.Nil(t, toolResult)
		assert.Equal(t, "test-task-id", result)
	})

	t.Run("ValidateTaskStatus - valid", func(t *testing.T) {
		params := map[string]interface{}{
			"durum": "devam_ediyor",
		}

		result, toolResult := validators.ValidateTaskStatus(params, true)
		assert.Nil(t, toolResult)
		assert.Equal(t, "devam_ediyor", result)
	})

	t.Run("ValidateTaskStatus - invalid", func(t *testing.T) {
		params := map[string]interface{}{
			"durum": "invalid_status",
		}

		result, toolResult := validators.ValidateTaskStatus(params, true)
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidateTaskPriority - valid", func(t *testing.T) {
		params := map[string]interface{}{
			"oncelik": "yuksek",
		}

		result, toolResult := validators.ValidateTaskPriority(params, true)
		assert.Nil(t, toolResult)
		assert.Equal(t, "yuksek", result)
	})

	t.Run("ValidateTaskPriority - invalid", func(t *testing.T) {
		params := map[string]interface{}{
			"oncelik": "invalid_priority",
		}

		result, toolResult := validators.ValidateTaskPriority(params, true)
		assert.NotNil(t, toolResult)
		assert.True(t, toolResult.IsError)
		assert.Equal(t, "", result)
	})

	t.Run("ValidatePagination - default values", func(t *testing.T) {
		params := map[string]interface{}{}

		limit, offset := validators.ValidatePagination(params)
		assert.Equal(t, constants.DefaultTaskLimit, limit)
		assert.Equal(t, constants.DefaultPaginationOffset, offset)
	})

	t.Run("ValidatePagination - custom values", func(t *testing.T) {
		params := map[string]interface{}{
			"limit":  25.0,
			"offset": 10.0,
		}

		limit, offset := validators.ValidatePagination(params)
		assert.Equal(t, 25, limit)
		assert.Equal(t, 10, offset)
	})

	t.Run("ValidatePagination - limit enforcement", func(t *testing.T) {
		params := map[string]interface{}{
			"limit":  500.0, // Should be capped at MaxTaskLimit
			"offset": -5.0,  // Should be reset to 0
		}

		limit, offset := validators.ValidatePagination(params)
		assert.Equal(t, constants.MaxTaskLimit, limit)
		assert.Equal(t, constants.DefaultPaginationOffset, offset)
	})

	t.Run("ValidatePagination - negative limit", func(t *testing.T) {
		params := map[string]interface{}{
			"limit": -10.0,
		}

		limit, offset := validators.ValidatePagination(params)
		assert.Equal(t, constants.DefaultTaskLimit, limit)
		assert.Equal(t, constants.DefaultPaginationOffset, offset)
	})
}

// TestPriorityFormatter tests the PriorityFormatter helper functions
func TestPriorityFormatter(t *testing.T) {
	formatter := NewPriorityFormatter()

	t.Run("GetPriorityShort", func(t *testing.T) {
		assert.Equal(t, "Y", formatter.GetPriorityShort(constants.PriorityHigh))
		assert.Equal(t, "O", formatter.GetPriorityShort(constants.PriorityMedium))
		assert.Equal(t, "D", formatter.GetPriorityShort(constants.PriorityLow))
		assert.Equal(t, "?", formatter.GetPriorityShort("unknown"))
	})

	t.Run("GetPriorityEmoji", func(t *testing.T) {
		assert.Equal(t, constants.EmojiPriorityHigh, formatter.GetPriorityEmoji(constants.PriorityHigh))
		assert.Equal(t, constants.EmojiPriorityMedium, formatter.GetPriorityEmoji(constants.PriorityMedium))
		assert.Equal(t, constants.EmojiPriorityLow, formatter.GetPriorityEmoji(constants.PriorityLow))
		assert.Equal(t, "⚪", formatter.GetPriorityEmoji("unknown"))
	})
}

// TestStatusFormatter tests the StatusFormatter helper functions
func TestStatusFormatter(t *testing.T) {
	formatter := NewStatusFormatter()

	t.Run("GetStatusEmoji", func(t *testing.T) {
		assert.Equal(t, constants.EmojiStatusPending, formatter.GetStatusEmoji(constants.TaskStatusPending))
		assert.Equal(t, constants.EmojiStatusInProgress, formatter.GetStatusEmoji(constants.TaskStatusInProgress))
		assert.Equal(t, constants.EmojiStatusCompleted, formatter.GetStatusEmoji(constants.TaskStatusCompleted))
		assert.Equal(t, constants.EmojiStatusCancelled, formatter.GetStatusEmoji(constants.TaskStatusCancelled))
		assert.Equal(t, constants.EmojiStatusUnknown, formatter.GetStatusEmoji("unknown"))
	})

	t.Run("GetStatusSymbol", func(t *testing.T) {
		assert.Equal(t, constants.EmojiStatusPending, formatter.GetStatusSymbol(constants.TaskStatusPending))
		assert.Equal(t, constants.EmojiStatusInProgress, formatter.GetStatusSymbol(constants.TaskStatusInProgress))
		assert.Equal(t, constants.BulletCheck, formatter.GetStatusSymbol(constants.TaskStatusCompleted))
		assert.Equal(t, constants.EmojiStatusCancelled, formatter.GetStatusSymbol(constants.TaskStatusCancelled))
		assert.Equal(t, "?", formatter.GetStatusSymbol("unknown"))
	})

	t.Run("GetStatusShort", func(t *testing.T) {
		assert.Equal(t, "P", formatter.GetStatusShort(constants.TaskStatusPending))
		assert.Equal(t, "I", formatter.GetStatusShort(constants.TaskStatusInProgress))
		assert.Equal(t, "C", formatter.GetStatusShort(constants.TaskStatusCompleted))
		assert.Equal(t, "X", formatter.GetStatusShort(constants.TaskStatusCancelled))
		assert.Equal(t, "?", formatter.GetStatusShort("unknown"))
	})
}

// TestTaskIDFormatter tests the TaskIDFormatter helper functions
func TestTaskIDFormatter(t *testing.T) {
	formatter := NewTaskIDFormatter()

	t.Run("FormatShortID - long ID", func(t *testing.T) {
		longID := "550e8400-e29b-41d4-a716-446655440000"
		result := formatter.FormatShortID(longID)
		assert.Equal(t, "550e8400", result)
	})

	t.Run("FormatShortID - short ID", func(t *testing.T) {
		shortID := "short123"
		result := formatter.FormatShortID(shortID)
		assert.Equal(t, "short123", result)
	})

	t.Run("FormatShortID - exact length", func(t *testing.T) {
		exactID := "550e8400"
		result := formatter.FormatShortID(exactID)
		assert.Equal(t, "550e8400", result)
	})

	t.Run("FormatTaskReference", func(t *testing.T) {
		taskID := "550e8400-e29b-41d4-a716-446655440000"
		title := "Test Task"
		result := formatter.FormatTaskReference(taskID, title)
		assert.Equal(t, "Test Task (550e8400)", result)
	})
}

// TestToolHelpers tests the ToolHelpers helper functions
func TestToolHelpers(t *testing.T) {
	helpers := NewToolHelpers()

	t.Run("NewToolHelpers", func(t *testing.T) {
		assert.NotNil(t, helpers)
		assert.NotNil(t, helpers.Validator)
		assert.NotNil(t, helpers.Formatter)
		assert.NotNil(t, helpers.ErrorFormatter)
		assert.NotNil(t, helpers.ResponseBuilder)
	})

	t.Run("ToolHelpers integration", func(t *testing.T) {
		// Test that all helper components work together
		params := map[string]interface{}{
			"id":           "test-task-id",
			"durum":        "devam_ediyor",
			"oncelik":      "yuksek",
			"limit":        25.0,
			"tum_projeler": true,
		}

		// Test validation through helpers
		taskID, toolResult := helpers.Validator.ValidateTaskID(params)
		assert.Nil(t, toolResult)
		assert.Equal(t, "test-task-id", taskID)

		status, toolResult := helpers.Validator.ValidateTaskStatus(params, false)
		assert.Nil(t, toolResult)
		assert.Equal(t, "devam_ediyor", status)

		priority, toolResult := helpers.Validator.ValidateTaskPriority(params, false)
		assert.Nil(t, toolResult)
		assert.Equal(t, "yuksek", priority)

		// Test formatting through helpers
		formattedTask := helpers.Formatter.FormatTaskWithStatus("Test Task", "550e8400-e29b-41d4-a716-446655440000", "devam_ediyor")
		assert.Contains(t, formattedTask, "🔄")
		assert.Contains(t, formattedTask, "Test Task")

		// Test pagination validation
		limit, offset := helpers.Validator.ValidatePagination(params)
		assert.Equal(t, 25, limit)
		assert.Equal(t, 0, offset)
	})
}

// TestGlobalFormatterInstances tests the global formatter instances
func TestGlobalFormatterInstances(t *testing.T) {
	t.Run("PriorityFormat global instance", func(t *testing.T) {
		assert.NotNil(t, PriorityFormat)
		assert.Equal(t, "Y", PriorityFormat.GetPriorityShort(constants.PriorityHigh))
	})

	t.Run("StatusFormat global instance", func(t *testing.T) {
		assert.NotNil(t, StatusFormat)
		assert.Equal(t, "🔄", StatusFormat.GetStatusEmoji(constants.TaskStatusInProgress))
	})

	t.Run("TaskIDFormat global instance", func(t *testing.T) {
		assert.NotNil(t, TaskIDFormat)
		longID := "550e8400-e29b-41d4-a716-446655440000"
		assert.Equal(t, "550e8400", TaskIDFormat.FormatShortID(longID))
	})
}

// TestFormatterConsistency tests that all formatters produce consistent output
func TestFormatterConsistency(t *testing.T) {
	t.Run("Status emoji consistency", func(t *testing.T) {
		taskFormatter := NewTaskFormatter()
		statusFormatter := NewStatusFormatter()

		assert.Equal(t, taskFormatter.GetStatusEmoji(constants.TaskStatusInProgress), statusFormatter.GetStatusEmoji(constants.TaskStatusInProgress))
		assert.Equal(t, taskFormatter.GetStatusEmoji(constants.TaskStatusCompleted), statusFormatter.GetStatusEmoji(constants.TaskStatusCompleted))
	})

	t.Run("Priority emoji consistency", func(t *testing.T) {
		taskFormatter := NewTaskFormatter()
		priorityFormatter := NewPriorityFormatter()

		assert.Equal(t, taskFormatter.GetPriorityEmoji(constants.PriorityHigh), priorityFormatter.GetPriorityEmoji(constants.PriorityHigh))
		assert.Equal(t, taskFormatter.GetPriorityEmoji(constants.PriorityMedium), priorityFormatter.GetPriorityEmoji(constants.PriorityMedium))
	})

	t.Run("ID formatting consistency", func(t *testing.T) {
		taskFormatter := NewTaskFormatter()
		idFormatter := NewTaskIDFormatter()

		longID := "550e8400-e29b-41d4-a716-446655440000"
		shortID := idFormatter.FormatShortID(longID)
		assert.Equal(t, shortID, "550e8400")

		// Task formatter should use short IDs
		formattedTask := taskFormatter.FormatTaskBasic("Test Task", longID)
		assert.Contains(t, formattedTask, shortID)
	})
}
