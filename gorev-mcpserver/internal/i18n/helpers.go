package i18n

import (
	"context"
	"fmt"
	"strings"

	"github.com/msenol/gorev/internal/constants"
)

// ==================== CONTEXT HELPERS ====================

// contextKey is a private type for context keys to avoid collisions
type contextKey string

// languageKey is the context key for language storage
const languageKey contextKey = "language"

// WithLanguage stores the language preference in context
func WithLanguage(ctx context.Context, lang string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, languageKey, lang)
}

// FromContext extracts the language from context, defaulting to Turkish
func FromContext(ctx context.Context) string {
	if ctx == nil {
		return "tr"
	}

	if lang, ok := ctx.Value(languageKey).(string); ok && lang != "" {
		return lang
	}

	return "tr"
}

// ==================== DRY i18n helper functions ====================

// TCommon gets a common pattern translation
func TCommon(lang string, key string, data map[string]interface{}) string {
	return TWithLang(lang, "common."+key, data)
}

// TParam gets parameter description using DRY patterns
func TParam(lang string, paramName string) string {
	// Check if specific param description exists
	specificKey := "tools.params.descriptions." + paramName
	if HasKey(specificKey) {
		return TWithLang(lang, specificKey, nil)
	}

	// Fallback to common field patterns
	commonKey := "common.fields." + paramName
	if HasKey(commonKey) {
		return TCommon(lang, "fields."+paramName, nil)
	}

	// Default fallback
	return paramName + " parameter"
}

// TValidation gets validation message using DRY patterns
func TValidation(lang string, validationType string, param string, extra map[string]interface{}) string {
	data := map[string]interface{}{
		"Param": param,
	}

	// Merge extra data if provided
	if extra != nil {
		for k, v := range extra {
			data[k] = v
		}
	}

	return TWithLang(lang, "validation."+validationType, data)
}

// BuildFieldDescription builds field description using common patterns
func BuildFieldDescription(lang string, prefix, entity, field string) string {
	data := map[string]interface{}{
		"Prefix": TCommon(lang, "prefixes."+prefix, nil),
		"Entity": TCommon(lang, "entities."+entity, nil),
		"Field":  field,
	}

	return TCommon(lang, "patterns.new_field", data)
}

// BuildIDDescription builds ID description using DRY patterns
func BuildIDDescription(lang string, entity, idType string) string {
	data := map[string]interface{}{
		"Entity": entity,
		"Type":   idType,
	}

	if idType == "unique" {
		return TCommon(lang, "fields.task_id", map[string]interface{}{"Type": TCommon(lang, "fields.id_base", nil)})
	}

	return TCommon(lang, "fields.task_id", data)
}

// BuildPaginationDescription builds pagination descriptions using common patterns
func BuildPaginationDescription(lang string, paginationType, entity string, defaultVal, maxVal int) string {
	data := map[string]interface{}{
		"Entity":  entity,
		"Default": defaultVal,
		"Max":     maxVal,
	}

	return TCommon(lang, "pagination."+paginationType+"_pattern", data)
}

// BuildPrefixedDescription builds descriptions with common prefixes
func BuildPrefixedDescription(lang string, prefix, target string) string {
	prefixText := TCommon(lang, "prefixes."+prefix, nil)
	if prefixText == "" {
		prefixText = prefix
	}

	return fmt.Sprintf("%s %s", prefixText, target)
}

// GetCommonSuffix gets common suffix patterns
func GetCommonSuffix(lang string, suffixType string) string {
	return TCommon(lang, "suffixes."+suffixType, nil)
}

// HasKey checks if a translation key exists (mock implementation for now)
func HasKey(key string) bool {
	// This would check if the key exists in the translation bundle
	// For now, we'll use simple heuristics based on our known structure

	commonParams := map[string]bool{
		"tools.params.descriptions.id_field":       true,
		"tools.params.descriptions.task_id":        true,
		"tools.params.descriptions.parent_id":      true,
		"tools.params.descriptions.gorev_id":       true,
		"tools.params.descriptions.proje_id":       true,
		"tools.params.descriptions.template_id":    true,
		"tools.params.descriptions.durum":          true,
		"tools.params.descriptions.baslik":         true,
		"tools.params.descriptions.aciklama":       true,
		"tools.params.descriptions.oncelik":        true,
		"tools.params.descriptions.son_tarih":      true,
		"tools.params.descriptions.onay":           true,
		"tools.params.descriptions.durum_filter":   true,
		"tools.params.descriptions.sirala":         true,
		"tools.params.descriptions.filtre":         true,
		"tools.params.descriptions.etiket":         true,
		"tools.params.descriptions.tum_projeler":   true,
		"tools.params.descriptions.limit":          true,
		"tools.params.descriptions.offset":         true,
		"tools.params.descriptions.isim":           true,
		"tools.params.descriptions.tanim":          true,
		"tools.params.descriptions.kategori":       true,
		"tools.params.descriptions.degerler":       true,
		"tools.params.descriptions.file_path":      true,
		"tools.params.descriptions.baglanti_tipi":  true,
		"tools.params.descriptions.query":          true,
		"tools.params.descriptions.updates":        true,
		"tools.params.descriptions.etiketler":      true,
		"tools.params.descriptions.kaynak_id":      true,
		"tools.params.descriptions.hedef_id":       true,
		"tools.params.descriptions.yeni_parent_id": true,
	}

	return commonParams[key]
}

// GetEntityName gets entity name in current language
func GetEntityName(lang string, entity string) string {
	return TCommon(lang, "entities."+entity, nil)
}

// FormatParameterRequired formats the "parameter required" message
func FormatParameterRequired(lang string, paramName string) string {
	return fmt.Sprintf("%s %s", paramName, GetCommonSuffix(lang, "required"))
}

// FormatInvalidValue formats the "invalid value" message
func FormatInvalidValue(lang string, paramName, value string, validValues []string) string {
	return fmt.Sprintf("%s %s: %s. Geçerli değerler: %s",
		paramName,
		GetCommonSuffix(lang, "invalid"),
		value,
		strings.Join(validValues, ", "))
}

// FormatEntityNotFound formats the "entity not found" message
func FormatEntityNotFound(lang string, entityType, id string) string {
	return fmt.Sprintf("%s bulunamadı: %s", entityType, id)
}

// FormatOperationFailed formats the "operation failed" message
func FormatOperationFailed(lang string, operation string, err error) string {
	return fmt.Sprintf("%s işlemi başarısız: %v", operation, err)
}

// ==================== NEW DRY HELPER FUNCTIONS ====================

// TRequiredParam formats the "parameter required" message using DRY pattern
func TRequiredParam(lang string, param string) string {
	return TCommon(lang, "validation.required", map[string]interface{}{
		"Param": param,
	})
}

// TRequiredArray formats the "parameter required and must be array" message
func TRequiredArray(lang string, param string) string {
	return TCommon(lang, "validation.required_array", map[string]interface{}{
		"Param": param,
	})
}

// TRequiredObject formats the "parameter required and must be object" message
func TRequiredObject(lang string, param string) string {
	return TCommon(lang, "validation.required_object", map[string]interface{}{
		"Param": param,
	})
}

// TEntityNotFound formats the "entity not found" message using DRY pattern
func TEntityNotFound(lang string, entity string, err error) string {
	return TCommon(lang, "validation.not_found", map[string]interface{}{
		"Entity": TCommon(lang, "entities."+entity, nil),
		"Error":  err.Error(),
	})
}

// TEntityNotFoundByID formats the "entity not found by ID" message
func TEntityNotFoundByID(lang string, entity, id string) string {
	return TCommon(lang, "validation.not_found_id", map[string]interface{}{
		"Entity": TCommon(lang, "entities."+entity, nil),
		"Id":     id,
	})
}

// TOperationFailed formats operation failure messages using DRY pattern
func TOperationFailed(lang string, operation, entity string, err error) string {
	return TCommon(lang, "operations."+operation+"_failed", map[string]interface{}{
		"Entity": TCommon(lang, "entities."+entity, nil),
		"Error":  err.Error(),
	})
}

// TSuccess formats success messages using DRY pattern
func TSuccess(lang string, operation, entity string, details map[string]interface{}) string {
	data := map[string]interface{}{
		"Entity": TCommon(lang, "entities."+entity, nil),
	}

	// Merge details if provided
	if details != nil {
		for k, v := range details {
			data[k] = v
		}
	}

	return TCommon(lang, "success."+operation, data)
}

// TInvalidValue formats invalid value messages using DRY pattern
func TInvalidValue(lang string, param, value string, validValues []string) string {
	return TCommon(lang, "validation.invalid_value", map[string]interface{}{
		"Param": param,
		"Value": value,
	}) + " Geçerli değerler: " + strings.Join(validValues, ", ")
}

// TInvalidStatus formats invalid status messages
func TInvalidStatus(lang string, status string, validStatuses []string) string {
	return TCommon(lang, "validation.invalid_status", map[string]interface{}{
		"Status":        status,
		"ValidStatuses": strings.Join(validStatuses, ", "),
	})
}

// TInvalidPriority formats invalid priority messages
func TInvalidPriority(lang string, priority string) string {
	return TCommon(lang, "validation.invalid_priority", map[string]interface{}{
		"Priority": priority,
	})
}

// TInvalidDate formats invalid date messages
func TInvalidDate(lang string, dateValue string) string {
	return TCommon(lang, "validation.invalid_date", map[string]interface{}{
		"Date": dateValue,
	})
}

// TInvalidFormat formats invalid format messages
func TInvalidFormat(lang string, formatType, value string) string {
	return TCommon(lang, "validation.invalid_format", map[string]interface{}{
		"Type":  formatType,
		"Value": value,
	})
}

// Specific operation helpers for common patterns
func TCreateFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "create", entity, err)
}

func TUpdateFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "update", entity, err)
}

func TDeleteFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "delete", entity, err)
}

func TFetchFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "fetch", entity, err)
}

func TSaveFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "save", entity, err)
}

func TSetFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "set", entity, err)
}

func TInitFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "init", entity, err)
}

func TCheckFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "check", entity, err)
}

func TQueryFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "query", entity, err)
}

func TProcessFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "process", entity, err)
}

func TListFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "list", entity, err)
}

func TEditFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "edit", entity, err)
}

func TAddFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "add", entity, err)
}

func TRemoveFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "remove", entity, err)
}

func TReadFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "read", entity, err)
}

func TConvertFailed(lang string, entity, format string, err error) string {
	return TCommon(lang, "operations.convert_failed", map[string]interface{}{
		"Entity": TCommon(lang, "entities."+entity, nil),
		"Format": format,
		"Error":  err.Error(),
	})
}

func TParseFailed(lang string, entity string, err error) string {
	return TOperationFailed(lang, "parse", entity, err)
}

// Success message helpers
func TCreated(lang string, entity, title, id string) string {
	return TSuccess(lang, "created", entity, map[string]interface{}{
		"Title": title,
		"Id":    id,
	})
}

func TUpdated(lang string, entity, details string) string {
	return TSuccess(lang, "updated", entity, map[string]interface{}{
		"Details": details,
	})
}

func TDeleted(lang string, entity, title, id string) string {
	return TSuccess(lang, "deleted", entity, map[string]interface{}{
		"Title": title,
		"Id":    id,
	})
}

func TSet(lang string, entity, details string) string {
	return TSuccess(lang, "set", entity, map[string]interface{}{
		"Details": details,
	})
}

func TRemoved(lang string, entity string) string {
	return TSuccess(lang, "removed", entity, nil)
}

func TAdded(lang string, entity, details string) string {
	return TSuccess(lang, "added", entity, map[string]interface{}{
		"Details": details,
	})
}

func TMoved(lang string, entity string) string {
	return TSuccess(lang, "moved", entity, nil)
}

func TEdited(lang string, entity, title string) string {
	return TSuccess(lang, "edited", entity, map[string]interface{}{
		"Title": title,
	})
}

// ==================== FIELD DESCRIPTION HELPERS ====================

// TFieldID returns ID field descriptions using DRY patterns
func TFieldID(lang string, entityType, action string) string {
	key := fmt.Sprintf("common.fields.id_descriptions.%s_%s", entityType, action)
	return TWithLang(lang, key, nil)
}

// TTaskCount returns task count descriptions with defaults
func TTaskCount(lang string, countType string, defaultVal ...string) string {
	desc := TWithLang(lang, fmt.Sprintf("common.fields.task_count.%s", countType), nil)
	if len(defaultVal) > 0 {
		desc += " " + TWithLang(lang, fmt.Sprintf("common.fields.defaults.default_%s", defaultVal[0]), nil)
	}
	return desc
}

// TProjectField returns project field descriptions
func TProjectField(lang string, field string) string {
	return TWithLang(lang, fmt.Sprintf("common.fields.project.project_%s", field), nil)
}

// TSubtaskField returns subtask field descriptions using DRY patterns
func TSubtaskField(lang string, field string) string {
	return TWithLang(lang, fmt.Sprintf("common.fields.subtask.%s", field), nil)
}

// TCommaSeparated returns comma-separated list descriptions
func TCommaSeparated(lang string, entity string) string {
	return TWithLang(lang, "common.fields.patterns.comma_separated", map[string]interface{}{
		"Entity": entity,
	})
}

// TWithFormat returns description with format pattern
func TWithFormat(lang string, description, format string) string {
	return TWithLang(lang, "common.fields.patterns.with_format", map[string]interface{}{
		"Description": description,
		"Format":      format,
	})
}

// TFilePath returns file path descriptions
func TFilePath(lang string, action string) string {
	return TWithLang(lang, fmt.Sprintf("common.fields.patterns.file_path.%s", action), nil)
}

// TTemplate returns template-related descriptions
func TTemplate(lang string, templateType string) string {
	return TWithLang(lang, fmt.Sprintf("common.fields.patterns.template.%s", templateType), nil)
}

// TBatch returns batch operation descriptions
func TBatch(lang string, batchType string) string {
	return TWithLang(lang, fmt.Sprintf("common.fields.patterns.batch.%s", batchType), nil)
}

// ==================== MARKDOWN FORMATTING HELPERS ====================

// TLabel returns a translated label for markdown formatting
func TLabel(lang string, labelKey string) string {
	return TWithLang(lang, fmt.Sprintf("common.labels.%s", labelKey), nil)
}

// TMarkdownLabel formats a markdown label with value
func TMarkdownLabel(lang string, labelKey string, value interface{}) string {
	label := TLabel(lang, labelKey)
	return fmt.Sprintf("**%s:** %v", label, value)
}

// TMarkdownHeader formats a markdown header with translated text
func TMarkdownHeader(lang string, level int, labelKey string) string {
	label := TLabel(lang, labelKey)
	prefix := strings.Repeat("#", level)
	return fmt.Sprintf("%s %s", prefix, label)
}

// TMarkdownBold formats text as bold markdown
func TMarkdownBold(lang string, labelKey string) string {
	label := TLabel(lang, labelKey)
	return fmt.Sprintf("**%s**", label)
}

// TMarkdownSection formats a section header with emoji and text
func TMarkdownSection(lang string, emoji, labelKey string) string {
	label := TLabel(lang, labelKey)
	return fmt.Sprintf("### %s %s", emoji, label)
}

// TCount formats count messages with label
func TCount(lang string, labelKey string, count int) string {
	label := TLabel(lang, labelKey)
	return fmt.Sprintf("**%s:** %d", label, count)
}

// TDuration formats duration with label
func TDuration(lang string, labelKey string, duration interface{}) string {
	label := TLabel(lang, labelKey)
	return fmt.Sprintf("**%s:** %v", label, duration)
}

// TListItem formats a list item with label and value
func TListItem(lang string, labelKey string, value interface{}) string {
	label := TLabel(lang, labelKey)
	return fmt.Sprintf("- **%s:** %v", label, value)
}

// ==================== STATUS/PRIORITY TRANSLATIONS ====================

// TStatus returns translated status text
func TStatus(lang string, status string) string {
	switch status {
	case constants.TaskStatusPending:
		return TWithLang(lang, "status.pending", nil)
	case constants.TaskStatusInProgress:
		return TWithLang(lang, "status.in_progress", nil)
	case constants.TaskStatusCompleted:
		return TWithLang(lang, "status.completed", nil)
	case constants.TaskStatusCancelled:
		return TWithLang(lang, "status.cancelled", nil)
	default:
		return status
	}
}

// TPriority returns translated priority text
func TPriority(lang string, priority string) string {
	switch priority {
	case constants.PriorityLow:
		return TWithLang(lang, "priority.low", nil)
	case constants.PriorityMedium:
		return TWithLang(lang, "priority.medium", nil)
	case constants.PriorityHigh:
		return TWithLang(lang, "priority.high", nil)
	default:
		return priority
	}
}
