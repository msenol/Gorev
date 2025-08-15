package gorev

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// NLPProcessor handles natural language queries for task management
type NLPProcessor struct {
	TimeZone *time.Location
}

// NewNLPProcessor creates a new NLP processor with the system timezone
func NewNLPProcessor() *NLPProcessor {
	return &NLPProcessor{
		TimeZone: time.Local,
	}
}

// QueryIntent represents the parsed intent from natural language
type QueryIntent struct {
	Action     string                 `json:"action"`
	Filters    map[string]interface{} `json:"filters"`
	Parameters map[string]interface{} `json:"parameters"`
	TimeRange  *TimeRange             `json:"time_range,omitempty"`
	Confidence float64                `json:"confidence"`
	Raw        string                 `json:"raw"`
}

// TimeRange represents a parsed time range
type TimeRange struct {
	Start    *time.Time `json:"start,omitempty"`
	End      *time.Time `json:"end,omitempty"`
	Relative string     `json:"relative,omitempty"`
}

// ProcessQuery analyzes natural language and returns structured intent
func (nlp *NLPProcessor) ProcessQuery(query string) (*QueryIntent, error) {
	intent := &QueryIntent{
		Filters:    make(map[string]interface{}),
		Parameters: make(map[string]interface{}),
		Confidence: 0.0,
		Raw:        query,
	}

	// Normalize query
	normalized := strings.ToLower(strings.TrimSpace(query))

	// Parse time expressions first
	if timeRange := nlp.parseTimeExpressions(normalized); timeRange != nil {
		intent.TimeRange = timeRange
		intent.Confidence += 0.3
	}

	// Parse action intent
	action := nlp.parseAction(normalized)
	if action != "" {
		intent.Action = action
		intent.Confidence += 0.4
	}

	// Parse filters and tags
	if filters := nlp.parseFilters(normalized); len(filters) > 0 {
		for k, v := range filters {
			intent.Filters[k] = v
		}
		intent.Confidence += 0.3
	}

	// Parse task references
	if refs := nlp.parseTaskReferences(normalized); len(refs) > 0 {
		intent.Parameters["task_references"] = refs
		intent.Confidence += 0.2
	}

	log.Printf("NLP Query processed: %s -> Action: %s, Confidence: %.2f",
		query, intent.Action, intent.Confidence)

	return intent, nil
}

// parseAction determines the main action from the query
func (nlp *NLPProcessor) parseAction(query string) string {
	actionPatterns := map[string][]string{
		"list": {
			"görevleri göster", "listele", "görevler", "ne var", "neler var",
			"show tasks", "list tasks", "what tasks", "tasks",
		},
		"create": {
			"görev oluştur", "yeni görev", "ekle", "oluştur", "yap",
			"create task", "new task", "add task", "make task",
		},
		"update": {
			"güncelle", "değiştir", "düzenle", "revize et",
			"update", "modify", "edit", "change",
		},
		"complete": {
			"tamamla", "bitir", "kapat", "hallettim", "bitti",
			"complete", "finish", "close", "done",
		},
		"delete": {
			"sil", "kaldır", "iptal et",
			"delete", "remove", "cancel",
		},
		"search": {
			"ara", "bul", "araştır",
			"search", "find", "look for",
		},
		"status": {
			"durum", "durumu", "nasıl", "ne durumda",
			"status", "state", "how is",
		},
	}

	for action, patterns := range actionPatterns {
		for _, pattern := range patterns {
			if strings.Contains(query, pattern) {
				return action
			}
		}
	}

	return "list" // Default action
}

// parseTimeExpressions extracts time-related information
func (nlp *NLPProcessor) parseTimeExpressions(query string) *TimeRange {
	now := time.Now().In(nlp.TimeZone)

	// Turkish time expressions
	turkishPatterns := map[string]func() *TimeRange{
		"bugün": func() *TimeRange {
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, nlp.TimeZone)
			end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
			return &TimeRange{Start: &start, End: &end, Relative: "today"}
		},
		"yarın": func() *TimeRange {
			tomorrow := now.Add(24 * time.Hour)
			start := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, nlp.TimeZone)
			end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
			return &TimeRange{Start: &start, End: &end, Relative: "tomorrow"}
		},
		"dün": func() *TimeRange {
			yesterday := now.Add(-24 * time.Hour)
			start := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, nlp.TimeZone)
			end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
			return &TimeRange{Start: &start, End: &end, Relative: "yesterday"}
		},
		"bu hafta": func() *TimeRange {
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			} // Sunday = 7
			start := now.Add(-time.Duration(weekday-1) * 24 * time.Hour)
			start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, nlp.TimeZone)
			end := start.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)
			return &TimeRange{Start: &start, End: &end, Relative: "this_week"}
		},
		"gelecek hafta": func() *TimeRange {
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			nextWeekStart := now.Add(time.Duration(7-weekday+1) * 24 * time.Hour)
			start := time.Date(nextWeekStart.Year(), nextWeekStart.Month(), nextWeekStart.Day(), 0, 0, 0, 0, nlp.TimeZone)
			end := start.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)
			return &TimeRange{Start: &start, End: &end, Relative: "next_week"}
		},
	}

	// English time expressions
	englishPatterns := map[string]func() *TimeRange{
		"today":     turkishPatterns["bugün"],
		"tomorrow":  turkishPatterns["yarın"],
		"yesterday": turkishPatterns["dün"],
		"this week": turkishPatterns["bu hafta"],
		"next week": turkishPatterns["gelecek hafta"],
	}

	// Merge patterns
	allPatterns := make(map[string]func() *TimeRange)
	for k, v := range turkishPatterns {
		allPatterns[k] = v
	}
	for k, v := range englishPatterns {
		allPatterns[k] = v
	}

	for pattern, timeFunc := range allPatterns {
		if strings.Contains(query, pattern) {
			return timeFunc()
		}
	}

	// Parse specific dates (YYYY-MM-DD format)
	dateRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	if matches := dateRegex.FindStringSubmatch(query); len(matches) > 1 {
		if date, err := time.Parse("2006-01-02", matches[1]); err == nil {
			date = date.In(nlp.TimeZone)
			start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, nlp.TimeZone)
			end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
			return &TimeRange{Start: &start, End: &end, Relative: "specific_date"}
		}
	}

	return nil
}

// parseFilters extracts filtering criteria
func (nlp *NLPProcessor) parseFilters(query string) map[string]interface{} {
	filters := make(map[string]interface{})

	// Tag filters
	tagRegex := regexp.MustCompile(`(?:etiket:|tag:)(\w+)`)
	if matches := tagRegex.FindAllStringSubmatch(query, -1); len(matches) > 0 {
		var tags []string
		for _, match := range matches {
			if len(match) > 1 {
				tags = append(tags, match[1])
			}
		}
		if len(tags) > 0 {
			filters["tags"] = tags
		}
	}

	// Priority filters
	priorityPatterns := map[string]string{
		"yüksek öncelik": "high",
		"high priority":  "high",
		"acil":           "urgent",
		"urgent":         "urgent",
		"düşük öncelik":  "low",
		"low priority":   "low",
	}

	for pattern, priority := range priorityPatterns {
		if strings.Contains(query, pattern) {
			filters["priority"] = priority
			break
		}
	}

	// Status filters
	statusPatterns := map[string]string{
		"açık":        "open",
		"open":        "open",
		"tamamlanan":  "completed",
		"completed":   "completed",
		"devam eden":  "in_progress",
		"in progress": "in_progress",
		"bekleyen":    "pending",
		"pending":     "pending",
	}

	for pattern, status := range statusPatterns {
		if strings.Contains(query, pattern) {
			filters["status"] = status
			break
		}
	}

	// Category filters
	if strings.Contains(query, "frontend") || strings.Contains(query, "ön yüz") {
		filters["category"] = "frontend"
	} else if strings.Contains(query, "backend") || strings.Contains(query, "arka plan") {
		filters["category"] = "backend"
	} else if strings.Contains(query, "bug") || strings.Contains(query, "hata") {
		filters["category"] = "bug"
	} else if strings.Contains(query, "feature") || strings.Contains(query, "özellik") {
		filters["category"] = "feature"
	}

	return filters
}

// parseTaskReferences extracts references to specific tasks
func (nlp *NLPProcessor) parseTaskReferences(query string) []string {
	var references []string

	// Task ID references
	idRegex := regexp.MustCompile(`(?:görev |task )?\#?(\d+)`)
	if matches := idRegex.FindAllStringSubmatch(query, -1); len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				references = append(references, "id:"+match[1])
			}
		}
	}

	// Recent task references
	recentPatterns := []string{
		"son oluşturduğum", "son görev", "en son", "last task", "latest task",
		"recent task", "son eklediğim",
	}

	for _, pattern := range recentPatterns {
		if strings.Contains(query, pattern) {
			references = append(references, "recent:1")
			break
		}
	}

	// Title-based references
	titleRegex := regexp.MustCompile(`"([^"]+)"`)
	if matches := titleRegex.FindAllStringSubmatch(query, -1); len(matches) > 0 {
		for _, match := range matches {
			if len(match) > 1 {
				references = append(references, "title:"+match[1])
			}
		}
	}

	return references
}

// BuildQuery converts NLP intent back to structured query
func (nlp *NLPProcessor) BuildQuery(intent *QueryIntent) map[string]interface{} {
	query := make(map[string]interface{})

	if intent.Action != "" {
		query["action"] = intent.Action
	}

	if len(intent.Filters) > 0 {
		query["filters"] = intent.Filters
	}

	if intent.TimeRange != nil {
		query["time_range"] = intent.TimeRange
	}

	if len(intent.Parameters) > 0 {
		query["parameters"] = intent.Parameters
	}

	query["confidence"] = intent.Confidence

	return query
}

// FormatResponse generates a natural language response
func (nlp *NLPProcessor) FormatResponse(action string, results interface{}, lang string) string {
	if lang == "" {
		lang = "tr" // Default to Turkish
	}

	templates := map[string]map[string]string{
		"tr": {
			"list_empty":       "Belirtilen kriterlere uygun görev bulunamadı.",
			"list_found":       "%d görev bulundu.",
			"create_success":   "Görev başarıyla oluşturuldu: %s",
			"update_success":   "Görev güncellendi: %s",
			"complete_success": "Görev tamamlandı: %s",
			"delete_success":   "Görev silindi: %s",
			"error":            "İşlem sırasında hata oluştu: %s",
		},
		"en": {
			"list_empty":       "No tasks found matching the specified criteria.",
			"list_found":       "Found %d tasks.",
			"create_success":   "Task created successfully: %s",
			"update_success":   "Task updated: %s",
			"complete_success": "Task completed: %s",
			"delete_success":   "Task deleted: %s",
			"error":            "Error occurred during operation: %s",
		},
	}

	template := templates[lang]
	if template == nil {
		template = templates["tr"]
	}

	switch action {
	case "list":
		if results == nil {
			return template["list_empty"]
		}
		// Assume results is a slice
		if tasksJson, ok := results.([]byte); ok {
			var tasks []interface{}
			if err := json.Unmarshal(tasksJson, &tasks); err == nil {
				if len(tasks) == 0 {
					return template["list_empty"]
				}
				return fmt.Sprintf(template["list_found"], len(tasks))
			}
		}
		return template["list_empty"]

	case "create":
		if title, ok := results.(string); ok {
			return fmt.Sprintf(template["create_success"], title)
		}
		return template["create_success"]

	case "update":
		if title, ok := results.(string); ok {
			return fmt.Sprintf(template["update_success"], title)
		}
		return template["update_success"]

	case "complete":
		if title, ok := results.(string); ok {
			return fmt.Sprintf(template["complete_success"], title)
		}
		return template["complete_success"]

	case "delete":
		if title, ok := results.(string); ok {
			return fmt.Sprintf(template["delete_success"], title)
		}
		return template["delete_success"]

	default:
		if err, ok := results.(error); ok {
			return fmt.Sprintf(template["error"], err.Error())
		}
		return template["error"]
	}
}

// ValidateIntent checks if the parsed intent is actionable
func (nlp *NLPProcessor) ValidateIntent(intent *QueryIntent) error {
	if intent.Confidence < 0.3 {
		return fmt.Errorf("query confidence too low (%.2f): %s", intent.Confidence, intent.Raw)
	}

	if intent.Action == "" {
		return fmt.Errorf("no clear action identified in query: %s", intent.Raw)
	}

	// Validate action-specific requirements
	switch intent.Action {
	case "create":
		// Create actions should have some content indication
		if !strings.Contains(strings.ToLower(intent.Raw), "görev") &&
			!strings.Contains(strings.ToLower(intent.Raw), "task") {
			return fmt.Errorf("create action requires task content specification")
		}
	case "update", "complete", "delete":
		// These actions need task reference
		if refs, ok := intent.Parameters["task_references"]; !ok || len(refs.([]string)) == 0 {
			return fmt.Errorf("%s action requires task reference", intent.Action)
		}
	}

	return nil
}

// ExtractTaskContent extracts task content from natural language
func (nlp *NLPProcessor) ExtractTaskContent(query string) map[string]interface{} {
	content := make(map[string]interface{})

	// Extract title (usually the main part after action words)
	normalized := strings.ToLower(query)

	// Remove action words to get the content
	actionWords := []string{
		"görev oluştur", "yeni görev", "ekle", "oluştur", "yap",
		"create task", "new task", "add task", "make task",
	}

	title := query
	for _, action := range actionWords {
		if strings.Contains(normalized, action) {
			parts := strings.Split(query, action)
			if len(parts) > 1 {
				title = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	// Clean title
	title = strings.Trim(title, `"'`)
	if title != "" {
		content["title"] = title
	}

	// Extract description if separated by colon or dash
	if strings.Contains(title, ":") {
		parts := strings.SplitN(title, ":", 2)
		content["title"] = strings.TrimSpace(parts[0])
		content["description"] = strings.TrimSpace(parts[1])
	} else if strings.Contains(title, " - ") {
		parts := strings.SplitN(title, " - ", 2)
		content["title"] = strings.TrimSpace(parts[0])
		content["description"] = strings.TrimSpace(parts[1])
	}

	// Extract due date from time expressions
	if timeRange := nlp.parseTimeExpressions(normalized); timeRange != nil && timeRange.Start != nil {
		content["due_date"] = timeRange.Start.Format("2006-01-02T15:04:05Z07:00")
	}

	return content
}
