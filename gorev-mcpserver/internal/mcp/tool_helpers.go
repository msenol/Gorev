package mcp

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// ParameterValidator handles common parameter validation patterns
type ParameterValidator struct {
	lang string
}

// NewParameterValidator creates a new parameter validator
func NewParameterValidator(lang string) *ParameterValidator {
	return &ParameterValidator{lang: lang}
}

// SetLanguage updates the language for the validator
func (pv *ParameterValidator) SetLanguage(lang string) {
	pv.lang = lang
}

// ValidateRequiredString validates a required string parameter
func (pv *ParameterValidator) ValidateRequiredString(params map[string]interface{}, paramName string) (string, *mcp.CallToolResult) {
	value, ok := params[paramName].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", mcp.NewToolResultError(i18n.TRequiredParam(pv.lang, paramName))
	}
	return strings.TrimSpace(value), nil
}

// ValidateOptionalString validates an optional string parameter
func (pv *ParameterValidator) ValidateOptionalString(params map[string]interface{}, paramName string) string {
	if value, ok := params[paramName].(string); ok {
		return strings.TrimSpace(value)
	}
	return ""
}

// ValidateEnum validates an enum parameter
func (pv *ParameterValidator) ValidateEnum(params map[string]interface{}, paramName string, validValues []string, required bool) (string, *mcp.CallToolResult) {
	value, exists := params[paramName].(string)

	if !exists || value == "" {
		if required {
			return "", mcp.NewToolResultError(i18n.TValidation(pv.lang, "param_required_with_values", paramName, map[string]interface{}{
				"Values": strings.Join(validValues, ", "),
			}))
		}
		return "", nil
	}

	for _, valid := range validValues {
		if value == valid {
			return value, nil
		}
	}

	return "", mcp.NewToolResultError(i18n.TInvalidValue(pv.lang, paramName, value, validValues))
}

// ValidateNumber validates a number parameter
func (pv *ParameterValidator) ValidateNumber(params map[string]interface{}, paramName string, defaultValue int) int {
	if value, ok := params[paramName].(float64); ok && value > 0 {
		return int(value)
	}
	return defaultValue
}

// ValidateBool validates a boolean parameter
func (pv *ParameterValidator) ValidateBool(params map[string]interface{}, paramName string) bool {
	if value, ok := params[paramName].(bool); ok {
		return value
	}
	return false
}

// TaskFormatter handles common task formatting patterns
type TaskFormatter struct{}

// NewTaskFormatter creates a new task formatter
func NewTaskFormatter() *TaskFormatter {
	return &TaskFormatter{}
}

// FormatTaskBasic formats a task in basic format with ID
func (tf *TaskFormatter) FormatTaskBasic(baslik, id string) string {
	shortID := id
	if len(id) > constants.ShortIDLength {
		shortID = id[:constants.ShortIDLength]
	}
	return fmt.Sprintf("**%s** (ID: %s)", baslik, shortID)
}

// FormatTaskWithStatus formats a task with status emoji
func (tf *TaskFormatter) FormatTaskWithStatus(baslik, id, durum string) string {
	statusEmoji := tf.GetStatusEmoji(durum)
	shortID := id
	if len(id) > constants.ShortIDLength {
		shortID = id[:constants.ShortIDLength]
	}
	return fmt.Sprintf("%s %s (ID: %s)", statusEmoji, baslik, shortID)
}

// FormatSuccessMessage formats a success message consistently
func (tf *TaskFormatter) FormatSuccessMessage(action, title, id string) string {
	return fmt.Sprintf("%s%s: %s (ID: %s)", constants.PrefixSuccess, action, title, id)
}

// GetStatusEmoji returns appropriate emoji for task status
func (tf *TaskFormatter) GetStatusEmoji(durum string) string {
	switch durum {
	case constants.TaskStatusCompleted:
		return constants.EmojiStatusCompleted
	case constants.TaskStatusInProgress:
		return constants.EmojiStatusInProgress
	case constants.TaskStatusPending:
		return constants.EmojiStatusPending
	case constants.TaskStatusCancelled:
		return constants.EmojiStatusCancelled
	default:
		return constants.EmojiStatusUnknown
	}
}

// GetPriorityEmoji returns appropriate emoji for task priority
func (tf *TaskFormatter) GetPriorityEmoji(oncelik string) string {
	switch oncelik {
	case constants.PriorityHigh:
		return constants.EmojiPriorityHigh
	case constants.PriorityMedium:
		return constants.EmojiPriorityMedium
	case constants.PriorityLow:
		return constants.EmojiPriorityLow
	default:
		return constants.EmojiPriorityUnknown
	}
}

// ErrorFormatter handles consistent error formatting
type ErrorFormatter struct{}

// NewErrorFormatter creates a new error formatter
func NewErrorFormatter() *ErrorFormatter {
	return &ErrorFormatter{}
}

// FormatNotFoundError formats "not found" errors consistently
func (ef *ErrorFormatter) FormatNotFoundError(entityType, id string) *mcp.CallToolResult {
	return mcp.NewToolResultError(i18n.TEntityNotFoundByID("tr", entityType, id))
}

// FormatOperationError formats operation errors consistently
// DEPRECATED: This function receives hardcoded Turkish strings which violates DRY.
// Use specific i18n helpers instead: TCreateFailed, TUpdateFailed, TDeleteFailed, etc.
// This function is kept temporarily for backwards compatibility.
func (ef *ErrorFormatter) FormatOperationError(operation string, err error) *mcp.CallToolResult {
	// Legacy implementation - just formats the hardcoded string with error
	return mcp.NewToolResultError(fmt.Sprintf("%s: %v", operation, err))
}

// FormatValidationError formats validation errors consistently
func (ef *ErrorFormatter) FormatValidationError(message string) *mcp.CallToolResult {
	return mcp.NewToolResultError(i18n.TValidation("tr", "validation_error", "", map[string]interface{}{
		"Message": message,
	}))
}

// ResponseBuilder helps build consistent response formats
type ResponseBuilder struct {
	formatter *TaskFormatter
}

// NewResponseBuilder creates a new response builder
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		formatter: NewTaskFormatter(),
	}
}

// BuildMarkdownTaskDetail builds detailed task information in markdown
func (rb *ResponseBuilder) BuildMarkdownTaskDetail(task interface{}) string {
	// This would be implemented based on the Gorev struct
	// For now, returning a placeholder
	return "# Task Detail\n\nTask details would be formatted here."
}

// BuildTaskList builds a formatted task list
func (rb *ResponseBuilder) BuildTaskList(tasks []interface{}, title string) string {
	var result strings.Builder

	if title != "" {
		result.WriteString(fmt.Sprintf("## %s\n\n", title))
	}

	if len(tasks) == 0 {
		result.WriteString(i18n.T("messages.no_tasks_found", nil) + "\n")
		return result.String()
	}

	result.WriteString(i18n.T("messages.tasks_found_count", map[string]interface{}{
		"Count": len(tasks),
	}) + "\n\n")

	// Task formatting would be implemented here based on task structure

	return result.String()
}

// CommonValidators contains common validation logic
type CommonValidators struct {
	paramValidator *ParameterValidator
	errorFormatter *ErrorFormatter
}

// NewCommonValidators creates a new common validators instance
func NewCommonValidators() *CommonValidators {
	return &CommonValidators{
		paramValidator: NewParameterValidator("tr"), // Default to Turkish
		errorFormatter: NewErrorFormatter(),
	}
}

// SetLanguage updates the language for all validators
func (cv *CommonValidators) SetLanguage(lang string) {
	cv.paramValidator.SetLanguage(lang)
}

// ValidateTaskID validates task ID parameter
func (cv *CommonValidators) ValidateTaskID(params map[string]interface{}) (string, *mcp.CallToolResult) {
	return cv.paramValidator.ValidateRequiredString(params, "id")
}

// ValidateTaskIDField validates task ID with custom field name
func (cv *CommonValidators) ValidateTaskIDField(params map[string]interface{}, fieldName string) (string, *mcp.CallToolResult) {
	return cv.paramValidator.ValidateRequiredString(params, fieldName)
}

// ValidateTaskStatus validates task status
func (cv *CommonValidators) ValidateTaskStatus(params map[string]interface{}, required bool) (string, *mcp.CallToolResult) {
	validStatuses := constants.GetValidTaskStatuses()
	return cv.paramValidator.ValidateEnum(params, "status", validStatuses, required)
}

// ValidateTaskPriority validates task priority
func (cv *CommonValidators) ValidateTaskPriority(params map[string]interface{}, required bool) (string, *mcp.CallToolResult) {
	validPriorities := constants.GetValidPriorities()
	return cv.paramValidator.ValidateEnum(params, "priority", validPriorities, required)
}

// ValidatePagination validates pagination parameters
func (cv *CommonValidators) ValidatePagination(params map[string]interface{}) (limit, offset int) {
	limit = cv.paramValidator.ValidateNumber(params, "limit", constants.DefaultTaskLimit)
	offset = cv.paramValidator.ValidateNumber(params, "offset", constants.DefaultPaginationOffset)

	// Ensure reasonable limits
	if limit > constants.MaxTaskLimit {
		limit = constants.MaxTaskLimit
	}
	if limit < 1 {
		limit = constants.DefaultTaskLimit
	}
	if offset < 0 {
		offset = constants.DefaultPaginationOffset
	}

	return limit, offset
}

// ValidateRequiredString validates a required string parameter
func (cv *CommonValidators) ValidateRequiredString(params map[string]interface{}, paramName string) (string, *mcp.CallToolResult) {
	return cv.paramValidator.ValidateRequiredString(params, paramName)
}

// ValidateEnum validates an enum parameter
func (cv *CommonValidators) ValidateEnum(params map[string]interface{}, paramName string, validValues []string, required bool) (string, *mcp.CallToolResult) {
	return cv.paramValidator.ValidateEnum(params, paramName, validValues, required)
}

// ValidateNumber validates a number parameter
func (cv *CommonValidators) ValidateNumber(params map[string]interface{}, paramName string, defaultValue int) int {
	return cv.paramValidator.ValidateNumber(params, paramName, defaultValue)
}

// ValidateBool validates a boolean parameter
func (cv *CommonValidators) ValidateBool(params map[string]interface{}, paramName string) bool {
	return cv.paramValidator.ValidateBool(params, paramName)
}

// ValidateOptionalString validates an optional string parameter
func (cv *CommonValidators) ValidateOptionalString(params map[string]interface{}, paramName string) string {
	return cv.paramValidator.ValidateOptionalString(params, paramName)
}

// ToolHelpers provides common functionality for handlers
type ToolHelpers struct {
	Validator       *CommonValidators
	Formatter       *TaskFormatter
	ErrorFormatter  *ErrorFormatter
	ResponseBuilder *ResponseBuilder
}

// NewToolHelpers creates a new tool helpers instance
func NewToolHelpers() *ToolHelpers {
	return &ToolHelpers{
		Validator:       NewCommonValidators(),
		Formatter:       NewTaskFormatter(),
		ErrorFormatter:  NewErrorFormatter(),
		ResponseBuilder: NewResponseBuilder(),
	}
}

// SetLanguage updates the language for all tool helpers
func (th *ToolHelpers) SetLanguage(lang string) {
	th.Validator.SetLanguage(lang)
}

// ==================== CENTRALIZED FORMATTERS ====================

// PriorityFormatter handles priority formatting consistently
type PriorityFormatter struct{}

// NewPriorityFormatter creates a new priority formatter
func NewPriorityFormatter() *PriorityFormatter {
	return &PriorityFormatter{}
}

// GetPriorityShort returns single letter priority code
func (pf *PriorityFormatter) GetPriorityShort(priority string) string {
	switch priority {
	case constants.PriorityHigh:
		return "Y"
	case constants.PriorityMedium:
		return "O"
	case constants.PriorityLow:
		return "D"
	default:
		return "?"
	}
}

// GetPriorityEmoji returns emoji for priority
func (pf *PriorityFormatter) GetPriorityEmoji(priority string) string {
	switch priority {
	case constants.PriorityHigh:
		return constants.EmojiPriorityHigh
	case constants.PriorityMedium:
		return constants.EmojiPriorityMedium
	case constants.PriorityLow:
		return constants.EmojiPriorityLow
	default:
		return "⚪"
	}
}

// StatusFormatter handles status formatting consistently
type StatusFormatter struct{}

// NewStatusFormatter creates a new status formatter
func NewStatusFormatter() *StatusFormatter {
	return &StatusFormatter{}
}

// GetStatusEmoji returns emoji for status
func (sf *StatusFormatter) GetStatusEmoji(status string) string {
	switch status {
	case constants.TaskStatusPending:
		return constants.EmojiStatusPending
	case constants.TaskStatusInProgress:
		return constants.EmojiStatusInProgress
	case constants.TaskStatusCompleted:
		return constants.EmojiStatusCompleted
	case constants.TaskStatusCancelled:
		return constants.EmojiStatusCancelled
	default:
		return constants.EmojiStatusUnknown
	}
}

// GetStatusSymbol returns simple symbol for status (for compact display)
func (sf *StatusFormatter) GetStatusSymbol(status string) string {
	switch status {
	case constants.TaskStatusPending:
		return constants.EmojiStatusPending
	case constants.TaskStatusInProgress:
		return constants.EmojiStatusInProgress
	case constants.TaskStatusCompleted:
		return constants.BulletCheck
	case constants.TaskStatusCancelled:
		return constants.EmojiStatusCancelled
	default:
		return "?"
	}
}

// GetStatusShort returns short status code
func (sf *StatusFormatter) GetStatusShort(status string) string {
	switch status {
	case constants.TaskStatusPending:
		return "P"
	case constants.TaskStatusInProgress:
		return "I"
	case constants.TaskStatusCompleted:
		return "C"
	case constants.TaskStatusCancelled:
		return "X"
	default:
		return "?"
	}
}

// TaskIDFormatter handles task ID formatting
type TaskIDFormatter struct{}

// NewTaskIDFormatter creates a new task ID formatter
func NewTaskIDFormatter() *TaskIDFormatter {
	return &TaskIDFormatter{}
}

// FormatShortID returns truncated task ID (first 8 characters)
func (tf *TaskIDFormatter) FormatShortID(taskID string) string {
	if len(taskID) > constants.ShortIDLength {
		return taskID[:constants.ShortIDLength]
	}
	return taskID
}

// FormatTaskReference formats task with short ID and title
func (tf *TaskIDFormatter) FormatTaskReference(taskID, title string) string {
	shortID := tf.FormatShortID(taskID)
	return fmt.Sprintf("%s (%s)", title, shortID)
}

// Global formatter instances for easy access
var (
	PriorityFormat = NewPriorityFormatter()
	StatusFormat   = NewStatusFormatter()
	TaskIDFormat   = NewTaskIDFormatter()
)
