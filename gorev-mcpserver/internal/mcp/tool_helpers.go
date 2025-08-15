package mcp

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// ParameterValidator handles common parameter validation patterns
type ParameterValidator struct{}

// NewParameterValidator creates a new parameter validator
func NewParameterValidator() *ParameterValidator {
	return &ParameterValidator{}
}

// ValidateRequiredString validates a required string parameter
func (pv *ParameterValidator) ValidateRequiredString(params map[string]interface{}, paramName string) (string, *mcp.CallToolResult) {
	value, ok := params[paramName].(string)
	if !ok || strings.TrimSpace(value) == "" {
		return "", mcp.NewToolResultError(fmt.Sprintf("%s parametresi gerekli", paramName))
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
			return "", mcp.NewToolResultError(fmt.Sprintf("%s parametresi gerekli. Geçerli değerler: %s",
				paramName, strings.Join(validValues, ", ")))
		}
		return "", nil
	}

	for _, valid := range validValues {
		if value == valid {
			return value, nil
		}
	}

	return "", mcp.NewToolResultError(fmt.Sprintf("%s için geçersiz değer: %s. Geçerli değerler: %s",
		paramName, value, strings.Join(validValues, ", ")))
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
	return fmt.Sprintf("**%s** (ID: %s)", baslik, id[:8])
}

// FormatTaskWithStatus formats a task with status emoji
func (tf *TaskFormatter) FormatTaskWithStatus(baslik, id, durum string) string {
	statusEmoji := tf.GetStatusEmoji(durum)
	return fmt.Sprintf("%s %s (ID: %s)", statusEmoji, baslik, id[:8])
}

// FormatSuccessMessage formats a success message consistently
func (tf *TaskFormatter) FormatSuccessMessage(action, title, id string) string {
	return fmt.Sprintf("✓ %s: %s (ID: %s)", action, title, id)
}

// GetStatusEmoji returns appropriate emoji for task status
func (tf *TaskFormatter) GetStatusEmoji(durum string) string {
	switch durum {
	case "tamamlandi":
		return "✅"
	case "devam_ediyor":
		return "🔄"
	case "beklemede":
		return "⏳"
	case "iptal":
		return "❌"
	default:
		return "⚪"
	}
}

// GetPriorityEmoji returns appropriate emoji for task priority
func (tf *TaskFormatter) GetPriorityEmoji(oncelik string) string {
	switch oncelik {
	case "yuksek":
		return "🔥"
	case "orta":
		return "⚡"
	case "dusuk":
		return "ℹ️"
	default:
		return "ℹ️"
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
	return mcp.NewToolResultError(fmt.Sprintf("%s bulunamadı: %s", entityType, id))
}

// FormatOperationError formats operation errors consistently
func (ef *ErrorFormatter) FormatOperationError(operation string, err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("%s işlemi başarısız: %v", operation, err))
}

// FormatValidationError formats validation errors consistently
func (ef *ErrorFormatter) FormatValidationError(message string) *mcp.CallToolResult {
	return mcp.NewToolResultError(fmt.Sprintf("Doğrulama hatası: %s", message))
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
		result.WriteString("*Görev bulunamadı*\n")
		return result.String()
	}

	result.WriteString(fmt.Sprintf("**%d görev bulundu**\n\n", len(tasks)))

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
		paramValidator: NewParameterValidator(),
		errorFormatter: NewErrorFormatter(),
	}
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
	validStatuses := []string{"beklemede", "devam_ediyor", "tamamlandi", "iptal"}
	return cv.paramValidator.ValidateEnum(params, "durum", validStatuses, required)
}

// ValidateTaskPriority validates task priority
func (cv *CommonValidators) ValidateTaskPriority(params map[string]interface{}, required bool) (string, *mcp.CallToolResult) {
	validPriorities := []string{"dusuk", "orta", "yuksek"}
	return cv.paramValidator.ValidateEnum(params, "oncelik", validPriorities, required)
}

// ValidatePagination validates pagination parameters
func (cv *CommonValidators) ValidatePagination(params map[string]interface{}) (limit, offset int) {
	limit = cv.paramValidator.ValidateNumber(params, "limit", 50)
	offset = cv.paramValidator.ValidateNumber(params, "offset", 0)

	// Ensure reasonable limits
	if limit > 200 {
		limit = 200
	}
	if limit < 1 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
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
