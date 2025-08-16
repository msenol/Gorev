package i18n

import (
	"fmt"
	"strings"
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
