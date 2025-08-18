package i18n

import (
	"fmt"
	"strings"

	"github.com/msenol/gorev/internal/constants"
)

// DRY i18n helper functions to eliminate repeated patterns

// TCommon gets a common pattern translation
func TCommon(key string, data map[string]interface{}) string {
	return T("common."+key, data)
}

// TParam gets parameter description using DRY patterns
func TParam(paramName string) string {
	// Check if specific param description exists
	specificKey := "tools.params.descriptions." + paramName
	if HasKey(specificKey) {
		return T(specificKey, nil)
	}

	// Fallback to common field patterns
	commonKey := "common.fields." + paramName
	if HasKey(commonKey) {
		return TCommon("fields."+paramName, nil)
	}

	// Default fallback
	return paramName + " parameter"
}

// TValidation gets validation message using DRY patterns
func TValidation(validationType string, param string, extra map[string]interface{}) string {
	data := map[string]interface{}{
		"Param": param,
	}

	// Merge extra data if provided
	if extra != nil {
		for k, v := range extra {
			data[k] = v
		}
	}

	return T("validation."+validationType, data)
}

// BuildFieldDescription builds field description using common patterns
func BuildFieldDescription(prefix, entity, field string) string {
	data := map[string]interface{}{
		"Prefix": TCommon("prefixes."+prefix, nil),
		"Entity": TCommon("entities."+entity, nil),
		"Field":  field,
	}

	return TCommon("patterns.new_field", data)
}

// BuildIDDescription builds ID description using DRY patterns
func BuildIDDescription(entity, idType string) string {
	data := map[string]interface{}{
		"Entity": entity,
		"Type":   idType,
	}

	if idType == "unique" {
		return TCommon("fields.task_id", map[string]interface{}{"Type": TCommon("fields.id_base", nil)})
	}

	return TCommon("fields.task_id", data)
}

// BuildPaginationDescription builds pagination descriptions using common patterns
func BuildPaginationDescription(paginationType, entity string, defaultVal, maxVal int) string {
	data := map[string]interface{}{
		"Entity":  entity,
		"Default": defaultVal,
		"Max":     maxVal,
	}

	return TCommon("pagination."+paginationType+"_pattern", data)
}

// BuildPrefixedDescription builds descriptions with common prefixes
func BuildPrefixedDescription(prefix, target string) string {
	prefixText := TCommon("prefixes."+prefix, nil)
	if prefixText == "" {
		prefixText = prefix
	}

	return fmt.Sprintf("%s %s", prefixText, target)
}

// GetCommonSuffix gets common suffix patterns
func GetCommonSuffix(suffixType string) string {
	return TCommon("suffixes."+suffixType, nil)
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
func GetEntityName(entity string) string {
	return TCommon("entities."+entity, nil)
}

// FormatParameterRequired formats the "parameter required" message
func FormatParameterRequired(paramName string) string {
	return fmt.Sprintf("%s %s", paramName, GetCommonSuffix("required"))
}

// FormatInvalidValue formats the "invalid value" message
func FormatInvalidValue(paramName, value string, validValues []string) string {
	return fmt.Sprintf("%s %s: %s. Geçerli değerler: %s",
		paramName,
		GetCommonSuffix("invalid"),
		value,
		strings.Join(validValues, ", "))
}

// FormatEntityNotFound formats the "entity not found" message
func FormatEntityNotFound(entityType, id string) string {
	return fmt.Sprintf("%s bulunamadı: %s", entityType, id)
}

// FormatOperationFailed formats the "operation failed" message
func FormatOperationFailed(operation string, err error) string {
	return fmt.Sprintf("%s işlemi başarısız: %v", operation, err)
}

// ==================== NEW DRY HELPER FUNCTIONS ====================

// TRequiredParam formats the "parameter required" message using DRY pattern
func TRequiredParam(param string) string {
	return TCommon("validation.required", map[string]interface{}{
		"Param": param,
	})
}

// TRequiredArray formats the "parameter required and must be array" message
func TRequiredArray(param string) string {
	return TCommon("validation.required_array", map[string]interface{}{
		"Param": param,
	})
}

// TRequiredObject formats the "parameter required and must be object" message
func TRequiredObject(param string) string {
	return TCommon("validation.required_object", map[string]interface{}{
		"Param": param,
	})
}

// TEntityNotFound formats the "entity not found" message using DRY pattern
func TEntityNotFound(entity string, err error) string {
	return TCommon("validation.not_found", map[string]interface{}{
		"Entity": TCommon("entities."+entity, nil),
		"Error":  err.Error(),
	})
}

// TEntityNotFoundByID formats the "entity not found by ID" message
func TEntityNotFoundByID(entity, id string) string {
	return TCommon("validation.not_found_id", map[string]interface{}{
		"Entity": TCommon("entities."+entity, nil),
		"Id":     id,
	})
}

// TOperationFailed formats operation failure messages using DRY pattern
func TOperationFailed(operation, entity string, err error) string {
	return TCommon("operations."+operation+"_failed", map[string]interface{}{
		"Entity": TCommon("entities."+entity, nil),
		"Error":  err.Error(),
	})
}

// TSuccess formats success messages using DRY pattern
func TSuccess(operation, entity string, details map[string]interface{}) string {
	data := map[string]interface{}{
		"Entity": TCommon("entities."+entity, nil),
	}

	// Merge details if provided
	if details != nil {
		for k, v := range details {
			data[k] = v
		}
	}

	return TCommon("success."+operation, data)
}

// TInvalidValue formats invalid value messages using DRY pattern
func TInvalidValue(param, value string, validValues []string) string {
	return TCommon("validation.invalid_value", map[string]interface{}{
		"Param": param,
		"Value": value,
	}) + " Geçerli değerler: " + strings.Join(validValues, ", ")
}

// TInvalidStatus formats invalid status messages
func TInvalidStatus(status string, validStatuses []string) string {
	return TCommon("validation.invalid_status", map[string]interface{}{
		"Status":        status,
		"ValidStatuses": strings.Join(validStatuses, ", "),
	})
}

// TInvalidPriority formats invalid priority messages
func TInvalidPriority(priority string) string {
	return TCommon("validation.invalid_priority", map[string]interface{}{
		"Priority": priority,
	})
}

// TInvalidDate formats invalid date messages
func TInvalidDate(dateValue string) string {
	return TCommon("validation.invalid_date", map[string]interface{}{
		"Date": dateValue,
	})
}

// TInvalidFormat formats invalid format messages
func TInvalidFormat(formatType, value string) string {
	return TCommon("validation.invalid_format", map[string]interface{}{
		"Type":  formatType,
		"Value": value,
	})
}

// Specific operation helpers for common patterns
func TCreateFailed(entity string, err error) string {
	return TOperationFailed("create", entity, err)
}

func TUpdateFailed(entity string, err error) string {
	return TOperationFailed("update", entity, err)
}

func TDeleteFailed(entity string, err error) string {
	return TOperationFailed("delete", entity, err)
}

func TFetchFailed(entity string, err error) string {
	return TOperationFailed("fetch", entity, err)
}

func TSaveFailed(entity string, err error) string {
	return TOperationFailed("save", entity, err)
}

func TSetFailed(entity string, err error) string {
	return TOperationFailed("set", entity, err)
}

func TInitFailed(entity string, err error) string {
	return TOperationFailed("init", entity, err)
}

func TCheckFailed(entity string, err error) string {
	return TOperationFailed("check", entity, err)
}

func TQueryFailed(entity string, err error) string {
	return TOperationFailed("query", entity, err)
}

func TProcessFailed(entity string, err error) string {
	return TOperationFailed("process", entity, err)
}

func TListFailed(entity string, err error) string {
	return TOperationFailed("list", entity, err)
}

func TEditFailed(entity string, err error) string {
	return TOperationFailed("edit", entity, err)
}

func TAddFailed(entity string, err error) string {
	return TOperationFailed("add", entity, err)
}

func TRemoveFailed(entity string, err error) string {
	return TOperationFailed("remove", entity, err)
}

func TReadFailed(entity string, err error) string {
	return TOperationFailed("read", entity, err)
}

func TConvertFailed(entity, format string, err error) string {
	return TCommon("operations.convert_failed", map[string]interface{}{
		"Entity": TCommon("entities."+entity, nil),
		"Format": format,
		"Error":  err.Error(),
	})
}

func TParseFailed(entity string, err error) string {
	return TOperationFailed("parse", entity, err)
}

// Success message helpers
func TCreated(entity, title, id string) string {
	return TSuccess("created", entity, map[string]interface{}{
		"Title": title,
		"Id":    id,
	})
}

func TUpdated(entity, details string) string {
	return TSuccess("updated", entity, map[string]interface{}{
		"Details": details,
	})
}

func TDeleted(entity, title, id string) string {
	return TSuccess("deleted", entity, map[string]interface{}{
		"Title": title,
		"Id":    id,
	})
}

func TSet(entity, details string) string {
	return TSuccess("set", entity, map[string]interface{}{
		"Details": details,
	})
}

func TRemoved(entity string) string {
	return TSuccess("removed", entity, nil)
}

func TAdded(entity, details string) string {
	return TSuccess("added", entity, map[string]interface{}{
		"Details": details,
	})
}

func TMoved(entity string) string {
	return TSuccess("moved", entity, nil)
}

func TEdited(entity, title string) string {
	return TSuccess("edited", entity, map[string]interface{}{
		"Title": title,
	})
}

// ==================== FIELD DESCRIPTION HELPERS ====================

// TFieldID returns ID field descriptions using DRY patterns
func TFieldID(entityType, action string) string {
	key := fmt.Sprintf("common.fields.id_descriptions.%s_%s", entityType, action)
	return T(key, nil)
}

// TTaskCount returns task count descriptions with defaults
func TTaskCount(countType string, defaultVal ...string) string {
	desc := T(fmt.Sprintf("common.fields.task_count.%s", countType), nil)
	if len(defaultVal) > 0 {
		desc += " " + T(fmt.Sprintf("common.fields.defaults.default_%s", defaultVal[0]), nil)
	}
	return desc
}

// TProjectField returns project field descriptions
func TProjectField(field string) string {
	return T(fmt.Sprintf("common.fields.project.project_%s", field), nil)
}

// TSubtaskField returns subtask field descriptions using DRY patterns
func TSubtaskField(field string) string {
	return T(fmt.Sprintf("common.fields.subtask.%s", field), nil)
}

// TCommaSeparated returns comma-separated list descriptions
func TCommaSeparated(entity string) string {
	return T("common.fields.patterns.comma_separated", map[string]interface{}{
		"Entity": entity,
	})
}

// TWithFormat returns description with format pattern
func TWithFormat(description, format string) string {
	return T("common.fields.patterns.with_format", map[string]interface{}{
		"Description": description,
		"Format":      format,
	})
}

// TFilePath returns file path descriptions
func TFilePath(action string) string {
	return T(fmt.Sprintf("common.fields.patterns.file_path.%s", action), nil)
}

// TTemplate returns template-related descriptions
func TTemplate(templateType string) string {
	return T(fmt.Sprintf("common.fields.patterns.template.%s", templateType), nil)
}

// TBatch returns batch operation descriptions
func TBatch(batchType string) string {
	return T(fmt.Sprintf("common.fields.patterns.batch.%s", batchType), nil)
}

// ==================== MARKDOWN FORMATTING HELPERS ====================

// TLabel returns a translated label for markdown formatting
func TLabel(labelKey string) string {
	return T(fmt.Sprintf("common.labels.%s", labelKey), nil)
}

// TMarkdownLabel formats a markdown label with value
func TMarkdownLabel(labelKey string, value interface{}) string {
	label := TLabel(labelKey)
	return fmt.Sprintf("**%s:** %v", label, value)
}

// TMarkdownHeader formats a markdown header with translated text
func TMarkdownHeader(level int, labelKey string) string {
	label := TLabel(labelKey)
	prefix := strings.Repeat("#", level)
	return fmt.Sprintf("%s %s", prefix, label)
}

// TMarkdownBold formats text as bold markdown
func TMarkdownBold(labelKey string) string {
	label := TLabel(labelKey)
	return fmt.Sprintf("**%s**", label)
}

// TMarkdownSection formats a section header with emoji and text
func TMarkdownSection(emoji, labelKey string) string {
	label := TLabel(labelKey)
	return fmt.Sprintf("### %s %s", emoji, label)
}

// TCount formats count messages with label
func TCount(labelKey string, count int) string {
	label := TLabel(labelKey)
	return fmt.Sprintf("**%s:** %d", label, count)
}

// TDuration formats duration with label
func TDuration(labelKey string, duration interface{}) string {
	label := TLabel(labelKey)
	return fmt.Sprintf("**%s:** %v", label, duration)
}

// TListItem formats a list item with label and value
func TListItem(labelKey string, value interface{}) string {
	label := TLabel(labelKey)
	return fmt.Sprintf("- **%s:** %v", label, value)
}

// ==================== STATUS/PRIORITY TRANSLATIONS ====================

// TStatus returns translated status text
func TStatus(status string) string {
	switch status {
	case constants.TaskStatusPending:
		return T("status.pending", nil)
	case constants.TaskStatusInProgress:
		return T("status.in_progress", nil)
	case constants.TaskStatusCompleted:
		return T("status.completed", nil)
	case constants.TaskStatusCancelled:
		return T("status.cancelled", nil)
	default:
		return status
	}
}

// TPriority returns translated priority text
func TPriority(priority string) string {
	switch priority {
	case constants.PriorityLow:
		return T("priority.low", nil)
	case constants.PriorityMedium:
		return T("priority.medium", nil)
	case constants.PriorityHigh:
		return T("priority.high", nil)
	default:
		return priority
	}
}
