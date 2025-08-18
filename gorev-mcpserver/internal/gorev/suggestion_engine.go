package gorev

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/msenol/gorev/internal/constants"
)

// SuggestionEngine provides intelligent suggestions for task management
type SuggestionEngine struct {
	veriYonetici     VeriYoneticiInterface
	aiContextManager *AIContextYonetici
}

// NewSuggestionEngine creates a new suggestion engine
func NewSuggestionEngine(vy VeriYoneticiInterface) *SuggestionEngine {
	return &SuggestionEngine{
		veriYonetici: vy,
	}
}

// SetAIContextManager sets the AI context manager for better suggestions
func (se *SuggestionEngine) SetAIContextManager(acm *AIContextYonetici) {
	se.aiContextManager = acm
}

// Suggestion represents a suggested action
type Suggestion struct {
	Type        string                 `json:"type"`              // "next_action", "similar_task", "template", "deadline_risk"
	Priority    string                 `json:"priority"`          // "high", "medium", "low"
	Title       string                 `json:"title"`             // Human readable title
	Description string                 `json:"description"`       // Detailed description
	Action      string                 `json:"action"`            // Suggested action to take
	Context     map[string]interface{} `json:"context"`           // Additional context data
	Confidence  float64                `json:"confidence"`        // Confidence score (0-1)
	TaskID      string                 `json:"task_id,omitempty"` // Related task ID if applicable
}

// SuggestionRequest represents a request for suggestions
type SuggestionRequest struct {
	SessionID    string   `json:"session_id,omitempty"`
	ActiveTaskID string   `json:"active_task_id,omitempty"`
	Limit        int      `json:"limit,omitempty"`
	Types        []string `json:"types,omitempty"` // Filter by suggestion types
}

// SuggestionResponse contains suggestions and metadata
type SuggestionResponse struct {
	Suggestions   []Suggestion  `json:"suggestions"`
	TotalCount    int           `json:"total_count"`
	GeneratedAt   time.Time     `json:"generated_at"`
	ExecutionTime time.Duration `json:"execution_time"`
}

// GetSuggestions returns intelligent suggestions based on context
func (se *SuggestionEngine) GetSuggestions(request SuggestionRequest) (*SuggestionResponse, error) {
	startTime := time.Now()

	log.Printf("Generating suggestions: sessionId=%s, activeTask=%s", request.SessionID, request.ActiveTaskID)

	var allSuggestions []Suggestion

	// Generate different types of suggestions
	if len(request.Types) == 0 || contains(request.Types, "next_action") {
		nextActions, err := se.generateNextActionSuggestions(request)
		if err != nil {
			log.Printf("Failed to generate next action suggestions: error=%v", err)
		} else {
			allSuggestions = append(allSuggestions, nextActions...)
		}
	}

	if len(request.Types) == 0 || contains(request.Types, "similar_task") {
		similarTasks, err := se.generateSimilarTaskSuggestions(request)
		if err != nil {
			log.Printf("Failed to generate similar task suggestions: error=%v", err)
		} else {
			allSuggestions = append(allSuggestions, similarTasks...)
		}
	}

	if len(request.Types) == 0 || contains(request.Types, "template") {
		templateSuggestions, err := se.generateTemplateSuggestions(request)
		if err != nil {
			log.Printf("Failed to generate template suggestions: error=%v", err)
		} else {
			allSuggestions = append(allSuggestions, templateSuggestions...)
		}
	}

	if len(request.Types) == 0 || contains(request.Types, "deadline_risk") {
		deadlineRisks, err := se.generateDeadlineRiskSuggestions(request)
		if err != nil {
			log.Printf("Failed to generate deadline risk suggestions: error=%v", err)
		} else {
			allSuggestions = append(allSuggestions, deadlineRisks...)
		}
	}

	// Sort by priority and confidence
	sort.Slice(allSuggestions, func(i, j int) bool {
		// First by priority (high > medium > low)
		priorityScore := map[string]int{"high": constants.PriorityScoreHigh, "medium": constants.PriorityScoreMedium, "low": constants.PriorityScoreLow}
		if priorityScore[allSuggestions[i].Priority] != priorityScore[allSuggestions[j].Priority] {
			return priorityScore[allSuggestions[i].Priority] > priorityScore[allSuggestions[j].Priority]
		}
		// Then by confidence
		return allSuggestions[i].Confidence > allSuggestions[j].Confidence
	})

	// Apply limit
	limit := request.Limit
	if limit <= 0 {
		limit = constants.DefaultSuggestionLimit // Default limit
	}
	if len(allSuggestions) > limit {
		allSuggestions = allSuggestions[:limit]
	}

	response := &SuggestionResponse{
		Suggestions:   allSuggestions,
		TotalCount:    len(allSuggestions),
		GeneratedAt:   time.Now(),
		ExecutionTime: time.Since(startTime),
	}

	log.Printf("Suggestions generated: count=%d, duration=%v", len(allSuggestions), response.ExecutionTime)

	return response, nil
}

// generateNextActionSuggestions suggests next actions based on priorities and context
func (se *SuggestionEngine) generateNextActionSuggestions(request SuggestionRequest) ([]Suggestion, error) {
	var suggestions []Suggestion

	// Get high priority pending tasks
	gorevler, err := se.veriYonetici.GorevListele(map[string]interface{}{"durum": constants.TaskStatusPending})
	if err != nil {
		return suggestions, err
	}

	// Filter for high priority tasks
	var highPriorityTasks []*Gorev
	for _, gorev := range gorevler {
		if gorev.Oncelik == constants.PriorityHigh {
			highPriorityTasks = append(highPriorityTasks, gorev)
		}
	}

	// Suggest starting high priority tasks
	for i, gorev := range highPriorityTasks {
		if i >= constants.MaxSuggestionsToShow { // Limit to top suggestions
			break
		}

		// Check if task can be started (no blocking dependencies)
		canStart, err := se.checkCanStartTask(gorev.ID)
		if err != nil {
			continue
		}

		if canStart {
			suggestions = append(suggestions, Suggestion{
				Type:        "next_action",
				Priority:    "high",
				Title:       "Yüksek öncelikli görevi başlat",
				Description: fmt.Sprintf("'%s' görevi yüksek öncelikli ve başlamaya hazır", gorev.Baslik),
				Action:      fmt.Sprintf("gorev_guncelle id='%s' durum='devam_ediyor'", gorev.ID),
				Context: map[string]interface{}{
					"task_title":    gorev.Baslik,
					"task_priority": gorev.Oncelik,
					"task_id":       gorev.ID,
				},
				Confidence: constants.ConfidenceVeryHigh,
				TaskID:     gorev.ID,
			})
		}
	}

	// Suggest completing tasks in progress
	devamEdenGorevler, err := se.veriYonetici.GorevListele(map[string]interface{}{"durum": constants.TaskStatusInProgress})
	if err == nil {
		for i, gorev := range devamEdenGorevler {
			if i >= 2 { // Limit to top 2
				break
			}

			suggestions = append(suggestions, Suggestion{
				Type:        "next_action",
				Priority:    "medium",
				Title:       "Devam eden görevi tamamla",
				Description: fmt.Sprintf("'%s' görevi devam ediyor, tamamlamayı düşünün", gorev.Baslik),
				Action:      fmt.Sprintf("gorev_guncelle id='%s' durum='tamamlandi'", gorev.ID),
				Context: map[string]interface{}{
					"task_title": gorev.Baslik,
					"task_id":    gorev.ID,
				},
				Confidence: constants.ConfidenceMedium,
				TaskID:     gorev.ID,
			})
		}
	}

	return suggestions, nil
}

// generateSimilarTaskSuggestions detects similar tasks and suggests templates
func (se *SuggestionEngine) generateSimilarTaskSuggestions(request SuggestionRequest) ([]Suggestion, error) {
	var suggestions []Suggestion

	if request.ActiveTaskID == "" {
		return suggestions, nil
	}

	// Get active task
	activeTask, err := se.veriYonetici.GorevDetay(request.ActiveTaskID)
	if err != nil {
		return suggestions, err
	}

	// Find similar tasks based on title and description
	allTasks, err := se.veriYonetici.GorevListele(map[string]interface{}{})
	if err != nil {
		return suggestions, err
	}

	var similarTasks []*Gorev
	activeWords := extractKeywords(activeTask.Baslik + " " + activeTask.Aciklama)

	for _, task := range allTasks {
		if task.ID == activeTask.ID {
			continue
		}

		taskWords := extractKeywords(task.Baslik + " " + task.Aciklama)
		similarity := calculateSimilarity(activeWords, taskWords)

		if similarity > constants.SimilarityThreshold { // Similarity threshold
			similarTasks = append(similarTasks, task)
		}
	}

	// Suggest reviewing similar completed tasks
	for i, task := range similarTasks {
		if i >= constants.MaxSuggestionsToShow { // Limit to top suggestions
			break
		}

		if task.Durum == constants.TaskStatusCompleted {
			suggestions = append(suggestions, Suggestion{
				Type:        "similar_task",
				Priority:    "medium",
				Title:       "Benzer tamamlanmış görev bulundu",
				Description: fmt.Sprintf("'%s' görevi mevcut görevinize benziyor ve tamamlanmış", task.Baslik),
				Action:      fmt.Sprintf("gorev_detay id='%s'", task.ID),
				Context: map[string]interface{}{
					"similar_task_title": task.Baslik,
					"similar_task_id":    task.ID,
					"active_task_id":     request.ActiveTaskID,
				},
				Confidence: constants.ConfidenceLow,
				TaskID:     task.ID,
			})
		}
	}

	return suggestions, nil
}

// generateTemplateSuggestions recommends templates based on task content
func (se *SuggestionEngine) generateTemplateSuggestions(request SuggestionRequest) ([]Suggestion, error) {
	var suggestions []Suggestion

	// Get available templates
	templates, err := se.veriYonetici.TemplateListele("")
	if err != nil {
		return suggestions, err
	}

	// Analyze recent tasks to suggest relevant templates
	recentTasks, err := se.veriYonetici.GorevListele(map[string]interface{}{})
	if err != nil {
		return suggestions, err
	}

	// Simple keyword-based template matching
	keywordCounts := make(map[string]int)
	for _, task := range recentTasks {
		words := extractKeywords(task.Baslik + " " + task.Aciklama)
		for _, word := range words {
			keywordCounts[word]++
		}
	}

	// Suggest templates based on frequent keywords
	for _, template := range templates {
		templateWords := extractKeywords(template.Isim + " " + template.Tanim)
		score := 0
		for _, word := range templateWords {
			score += keywordCounts[word]
		}

		if score >= 2 { // Minimum relevance threshold
			suggestions = append(suggestions, Suggestion{
				Type:        "template",
				Priority:    "low",
				Title:       "İlgili template önerisi",
				Description: fmt.Sprintf("'%s' template'i son görevlerinize uygun görünüyor", template.Isim),
				Action:      fmt.Sprintf("template_listele kategori='%s'", template.Kategori),
				Context: map[string]interface{}{
					constants.ParamTemplateID: template.ID,
					"template_name":           template.Isim,
					"template_category":       template.Kategori,
					"relevance_score":         score,
				},
				Confidence: float64(score) / constants.ConfidenceNormalizer, // Normalize score
				TaskID:     "",
			})
		}
	}

	return suggestions, nil
}

// generateDeadlineRiskSuggestions analyzes deadline risks
func (se *SuggestionEngine) generateDeadlineRiskSuggestions(request SuggestionRequest) ([]Suggestion, error) {
	var suggestions []Suggestion

	// Get tasks with deadlines
	allTasks, err := se.veriYonetici.GorevListele(map[string]interface{}{})
	if err != nil {
		return suggestions, err
	}

	now := time.Now()
	var riskyTasks []*Gorev

	for _, task := range allTasks {
		if task.SonTarih == nil || task.Durum == constants.TaskStatusCompleted {
			continue
		}

		daysUntilDeadline := task.SonTarih.Sub(now).Hours() / 24

		// Tasks due within 3 days are high risk
		if daysUntilDeadline <= 3 && daysUntilDeadline > 0 {
			riskyTasks = append(riskyTasks, task)
		}

		// Overdue tasks are critical
		if daysUntilDeadline < 0 {
			suggestions = append(suggestions, Suggestion{
				Type:        "deadline_risk",
				Priority:    "high",
				Title:       "Gecikmiş görev!",
				Description: fmt.Sprintf("'%s' görevi %d gün gecikmiş", task.Baslik, int(-daysUntilDeadline)),
				Action:      fmt.Sprintf("gorev_detay id='%s'", task.ID),
				Context: map[string]interface{}{
					"task_title":   task.Baslik,
					"task_id":      task.ID,
					"days_overdue": int(-daysUntilDeadline),
					"deadline":     task.SonTarih.Format(constants.DateFormatISO),
				},
				Confidence: constants.ConfidenceVeryHigh,
				TaskID:     task.ID,
			})
		}
	}

	// Sort risky tasks by deadline proximity
	sort.Slice(riskyTasks, func(i, j int) bool {
		return riskyTasks[i].SonTarih.Before(*riskyTasks[j].SonTarih)
	})

	// Add deadline risk suggestions
	for i, task := range riskyTasks {
		if i >= constants.MaxSuggestionsToShow { // Limit to top suggestions
			break
		}

		daysUntil := int(task.SonTarih.Sub(now).Hours() / 24)
		suggestions = append(suggestions, Suggestion{
			Type:        "deadline_risk",
			Priority:    "high",
			Title:       "Yaklaşan son tarih",
			Description: fmt.Sprintf("'%s' görevi %d gün içinde bitiyor", task.Baslik, daysUntil),
			Action:      fmt.Sprintf("gorev_guncelle id='%s' durum='devam_ediyor'", task.ID),
			Context: map[string]interface{}{
				"task_title": task.Baslik,
				"task_id":    task.ID,
				"days_until": daysUntil,
				"deadline":   task.SonTarih.Format(constants.DateFormatISO),
			},
			Confidence: constants.ConfidenceHigh,
			TaskID:     task.ID,
		})
	}

	return suggestions, nil
}

// Helper functions

func (se *SuggestionEngine) checkCanStartTask(taskID string) (bool, error) {
	// Simple check - if no blocking dependencies, task can start
	// This could be enhanced with more sophisticated dependency analysis
	dependencies, err := se.veriYonetici.GorevBagimlilikGetir(taskID)
	if err != nil {
		return true, nil // Assume can start if we can't check dependencies
	}

	for _, dep := range dependencies {
		if dep.Durum != constants.TaskStatusCompleted {
			return false, nil
		}
	}

	return true, nil
}

func extractKeywords(text string) []string {
	// Simple keyword extraction - split on spaces and common separators
	words := strings.FieldsFunc(strings.ToLower(text), func(c rune) bool {
		return c == ' ' || c == ',' || c == '.' || c == ';' || c == ':' || c == '!' || c == '?'
	})

	// Filter out common stop words and short words
	stopWords := map[string]bool{
		"ve": true, "ile": true, "için": true, "bir": true, "bu": true, "şu": true,
		"the": true, "and": true, "or": true, "but": true, "in": true, "on": true, "at": true,
	}

	var keywords []string
	for _, word := range words {
		if len(word) > constants.MinWordLength && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

func calculateSimilarity(words1, words2 []string) float64 {
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Simple Jaccard similarity
	set1 := make(map[string]bool)
	for _, word := range words1 {
		set1[word] = true
	}

	set2 := make(map[string]bool)
	for _, word := range words2 {
		set2[word] = true
	}

	intersection := 0
	for word := range set1 {
		if set2[word] {
			intersection++
		}
	}

	union := len(set1) + len(set2) - intersection
	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
