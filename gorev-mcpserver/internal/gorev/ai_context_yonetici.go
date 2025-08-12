package gorev

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/msenol/gorev/internal/i18n"
)

// AIInteraction represents a single AI interaction with a task
type AIInteraction struct {
	ID         string    `json:"id"`
	GorevID    string    `json:"gorev_id"`
	ActionType string    `json:"action_type"`
	Context    string    `json:"context,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// AIContext represents the current AI session context
type AIContext struct {
	ActiveTaskID string                 `json:"active_task_id,omitempty"`
	RecentTasks  []string               `json:"recent_tasks"`
	SessionData  map[string]interface{} `json:"session_data"`
	LastUpdated  time.Time              `json:"last_updated"`
}

// AIContextSummary provides a summary optimized for AI usage
type AIContextSummary struct {
	ActiveTask     *Gorev   `json:"active_task,omitempty"`
	RecentTasks    []*Gorev `json:"recent_tasks"`
	WorkingProject *Proje   `json:"working_project,omitempty"`
	SessionSummary struct {
		Created   int `json:"created"`
		Updated   int `json:"updated"`
		Completed int `json:"completed"`
	} `json:"session_summary"`
	NextPriorities []*Gorev `json:"next_priorities"`
	Blockers       []*Gorev `json:"blockers"`
}

// AIContextYonetici manages AI context and interactions
type AIContextYonetici struct {
	veriYonetici VeriYoneticiInterface
}

// YeniAIContextYonetici creates a new AI context manager
func YeniAIContextYonetici(vy VeriYoneticiInterface) *AIContextYonetici {
	return &AIContextYonetici{
		veriYonetici: vy,
	}
}

// SetActiveTask sets the active task for the AI session
func (acy *AIContextYonetici) SetActiveTask(taskID string) error {
	// Validate task exists
	gorev, err := acy.veriYonetici.GorevGetir(taskID)
	if err != nil {
		return fmt.Errorf(i18n.T("error.taskNotFoundAi", map[string]interface{}{"Error": err}))
	}

	// Get current context
	context, err := acy.GetContext()
	if err != nil {
		// Initialize new context if not exists
		context = &AIContext{
			RecentTasks: []string{},
			SessionData: make(map[string]interface{}),
		}
	}

	// Update context
	context.ActiveTaskID = taskID
	context.LastUpdated = time.Now()

	// Add to recent tasks if not already there
	if !contains(context.RecentTasks, taskID) {
		context.RecentTasks = append([]string{taskID}, context.RecentTasks...)
		if len(context.RecentTasks) > 10 {
			context.RecentTasks = context.RecentTasks[:10]
		}
	}

	// Save context
	if err := acy.saveContext(context); err != nil {
		return fmt.Errorf(i18n.T("error.contextSaveFailed", map[string]interface{}{"Error": err}))
	}

	// Record interaction
	if err := acy.recordInteraction(taskID, "set_active", nil); err != nil {
		return fmt.Errorf(i18n.T("error.interactionSaveFailed", map[string]interface{}{"Error": err}))
	}

	// Auto-transition to "devam_ediyor" if task is in "beklemede"
	if gorev.Durum == "beklemede" {
		gorev.Durum = "devam_ediyor"
		if err := acy.veriYonetici.GorevGuncelle(gorev); err != nil {
			return fmt.Errorf(i18n.T("error.statusUpdateFailed", map[string]interface{}{"Error": err}))
		}
	}

	return nil
}

// GetActiveTask returns the current active task
func (acy *AIContextYonetici) GetActiveTask() (*Gorev, error) {
	context, err := acy.GetContext()
	if err != nil {
		return nil, err
	}

	if context.ActiveTaskID == "" {
		return nil, nil
	}

	return acy.veriYonetici.GorevGetir(context.ActiveTaskID)
}

// GetRecentTasks returns the recent tasks interacted with
func (acy *AIContextYonetici) GetRecentTasks(limit int) ([]*Gorev, error) {
	context, err := acy.GetContext()
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > len(context.RecentTasks) {
		limit = len(context.RecentTasks)
	}

	tasks := make([]*Gorev, 0, limit)
	for i := 0; i < limit; i++ {
		gorev, err := acy.veriYonetici.GorevGetir(context.RecentTasks[i])
		if err != nil {
			continue // Skip if task not found
		}
		tasks = append(tasks, gorev)
	}

	return tasks, nil
}

// GetContextSummary returns an AI-optimized summary of the current context
func (acy *AIContextYonetici) GetContextSummary() (*AIContextSummary, error) {
	summary := &AIContextSummary{}

	// Get active task
	activeTask, _ := acy.GetActiveTask()
	summary.ActiveTask = activeTask

	// Get recent tasks
	recentTasks, _ := acy.GetRecentTasks(5)
	summary.RecentTasks = recentTasks

	// Get working project
	if activeTask != nil && activeTask.ProjeID != "" {
		proje, _ := acy.veriYonetici.ProjeGetir(activeTask.ProjeID)
		summary.WorkingProject = proje
	}

	// Get session summary from interactions
	interactions, err := acy.getSessionInteractions()
	if err == nil {
		for _, interaction := range interactions {
			switch interaction.ActionType {
			case "created":
				summary.SessionSummary.Created++
			case "updated":
				summary.SessionSummary.Updated++
			case "completed":
				summary.SessionSummary.Completed++
			}
		}
	}

	// Get next priorities (high priority, not completed)
	gorevler, _ := acy.veriYonetici.GorevleriGetir("beklemede", "", "")
	for _, g := range gorevler {
		if g.Oncelik == "yuksek" {
			summary.NextPriorities = append(summary.NextPriorities, g)
			if len(summary.NextPriorities) >= 5 {
				break
			}
		}
	}

	// Get blockers (tasks with unfinished dependencies)
	for _, g := range gorevler {
		if g.TamamlanmamisBagimlilikSayisi > 0 {
			summary.Blockers = append(summary.Blockers, g)
			if len(summary.Blockers) >= 5 {
				break
			}
		}
	}

	return summary, nil
}

// RecordTaskView records when a task is viewed and auto-transitions state
func (acy *AIContextYonetici) RecordTaskView(taskID string) error {
	// Record interaction
	if err := acy.recordInteraction(taskID, "viewed", nil); err != nil {
		return err
	}

	// Get task
	gorev, err := acy.veriYonetici.GorevGetir(taskID)
	if err != nil {
		return err
	}

	// Auto-transition to "devam_ediyor" if in "beklemede"
	if gorev.Durum == "beklemede" {
		gorev.Durum = "devam_ediyor"
		if err := acy.veriYonetici.GorevGuncelle(gorev); err != nil {
			return fmt.Errorf(i18n.T("error.autoStatusUpdateFailed", map[string]interface{}{"Error": err}))
		}
		// Record the state change
		if err := acy.recordInteraction(taskID, "updated", map[string]interface{}{
			"auto_state_change": true,
			"from":              "beklemede",
			"to":                "devam_ediyor",
		}); err != nil {
			return err
		}
	}

	// Update last AI interaction time
	if err := acy.updateLastInteraction(taskID); err != nil {
		return err
	}
	
	// Add task to recent tasks
	return acy.addToRecentTasks(taskID)
}

// Helper functions

func (acy *AIContextYonetici) GetContext() (*AIContext, error) {
	return acy.veriYonetici.AIContextGetir()
}

func (acy *AIContextYonetici) saveContext(context *AIContext) error {
	return acy.veriYonetici.AIContextKaydet(context)
}

func (acy *AIContextYonetici) recordInteraction(taskID, actionType string, context interface{}) error {
	// Convert context to JSON string if provided
	contextJSON := ""
	if context != nil {
		if ctxBytes, err := json.Marshal(context); err == nil {
			contextJSON = string(ctxBytes)
		}
	}

	interaction := &AIInteraction{
		ID:         "", // Will be auto-generated by database
		GorevID:    taskID,
		ActionType: actionType,
		Context:    contextJSON,
		Timestamp:  time.Now(),
	}

	return acy.veriYonetici.AIInteractionKaydet(interaction)
}

func (acy *AIContextYonetici) getSessionInteractions() ([]*AIInteraction, error) {
	return acy.veriYonetici.AIInteractionlariGetir(50) // Get last 50 interactions for session
}

func (acy *AIContextYonetici) updateLastInteraction(taskID string) error {
	return acy.veriYonetici.AILastInteractionGuncelle(taskID, time.Now())
}

func (acy *AIContextYonetici) addToRecentTasks(taskID string) error {
	// Get current AI context
	context, err := acy.GetContext()
	if err != nil {
		// If no context exists, create new one
		context = &AIContext{
			RecentTasks: []string{},
		}
	}
	
	// Remove taskID if it already exists (to move it to front)
	newRecentTasks := []string{}
	for _, id := range context.RecentTasks {
		if id != taskID {
			newRecentTasks = append(newRecentTasks, id)
		}
	}
	
	// Add taskID to front
	context.RecentTasks = append([]string{taskID}, newRecentTasks...)
	
	// Keep only last 10 tasks
	if len(context.RecentTasks) > 10 {
		context.RecentTasks = context.RecentTasks[:10]
	}
	
	// Save updated context
	return acy.saveContext(context)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// BatchUpdate represents a single update in a batch operation
type BatchUpdate struct {
	ID      string                 `json:"id"`
	Updates map[string]interface{} `json:"updates"`
}

// BatchUpdateResult represents the result of a batch update operation
type BatchUpdateResult struct {
	Successful []string `json:"successful"`
	Failed     []struct {
		ID    string `json:"id"`
		Error string `json:"error"`
	} `json:"failed"`
	TotalProcessed int `json:"total_processed"`
}

// BatchUpdate performs multiple task updates in a single operation
func (acy *AIContextYonetici) BatchUpdate(updates []BatchUpdate) (*BatchUpdateResult, error) {
	result := &BatchUpdateResult{
		Successful: []string{},
		Failed: []struct {
			ID    string `json:"id"`
			Error string `json:"error"`
		}{},
	}

	for _, update := range updates {
		// Validate task exists
		_, err := acy.veriYonetici.GorevGetir(update.ID)
		if err != nil {
			result.Failed = append(result.Failed, struct {
				ID    string `json:"id"`
				Error string `json:"error"`
			}{
				ID:    update.ID,
				Error: fmt.Sprintf("görev bulunamadı: %v", err),
			})
			continue
		}

		// Apply updates based on fields
		if durum, ok := update.Updates["durum"].(string); ok {
			gorev, _ := acy.veriYonetici.GorevGetir(update.ID)
			gorev.Durum = durum
			if err := acy.veriYonetici.GorevGuncelle(gorev); err != nil {
				result.Failed = append(result.Failed, struct {
					ID    string `json:"id"`
					Error string `json:"error"`
				}{
					ID:    update.ID,
					Error: fmt.Sprintf("durum güncelleme hatası: %v", err),
				})
				continue
			}
		}

		// TODO: Add support for other fields (priority, description, etc.)

		result.Successful = append(result.Successful, update.ID)

		// Record batch operation
		if err := acy.recordInteraction(update.ID, "bulk_operation", update.Updates); err != nil {
			// Log but don't fail the operation
			// fmt.Printf("interaction kaydetme hatası: %v\n", err)
		}
	}

	result.TotalProcessed = len(updates)
	return result, nil
}

// NLPQuery performs natural language query on tasks
func (acy *AIContextYonetici) NLPQuery(query string) ([]*Gorev, error) {
	// Normalize the query to lowercase for easier matching
	normalizedQuery := strings.ToLower(query)

	// Define query patterns and their corresponding actions
	patterns := map[string]func() ([]*Gorev, error){
		"bugün": func() ([]*Gorev, error) {
			// Tasks interacted with today
			interactions, err := acy.getTodayInteractions()
			if err != nil {
				return nil, err
			}
			return acy.getTasksFromInteractions(interactions)
		},
		"son oluşturduğum": func() ([]*Gorev, error) {
			// Last created task
			return acy.getLastCreatedTasks(1)
		},
		"son oluşturulan": func() ([]*Gorev, error) {
			// Recently created tasks
			return acy.getLastCreatedTasks(5)
		},
		"yüksek öncelik": func() ([]*Gorev, error) {
			// High priority tasks
			return acy.veriYonetici.GorevleriGetir("beklemede", "", "")
		},
		"tamamlanmamış": func() ([]*Gorev, error) {
			// Incomplete tasks
			return acy.veriYonetici.GorevleriGetir("beklemede", "", "")
		},
		"devam eden": func() ([]*Gorev, error) {
			// In progress tasks
			return acy.veriYonetici.GorevleriGetir("devam_ediyor", "", "")
		},
		"tamamlanan": func() ([]*Gorev, error) {
			// Completed tasks
			return acy.veriYonetici.GorevleriGetir("tamamlandi", "", "")
		},
		"blokaj": func() ([]*Gorev, error) {
			// Blocked tasks
			gorevler, _ := acy.veriYonetici.GorevleriGetir("", "", "")
			var blocked []*Gorev
			for _, g := range gorevler {
				if g.TamamlanmamisBagimlilikSayisi > 0 {
					blocked = append(blocked, g)
				}
			}
			return blocked, nil
		},
		"acil": func() ([]*Gorev, error) {
			// Urgent tasks (due soon)
			return acy.veriYonetici.GorevleriGetir("", "", "acil")
		},
		"gecikmiş": func() ([]*Gorev, error) {
			// Overdue tasks
			return acy.veriYonetici.GorevleriGetir("", "", "gecmis")
		},
	}

	// Check for keyword matches
	for pattern, handler := range patterns {
		if strings.Contains(normalizedQuery, pattern) {
			return handler()
		}
	}

	// Check for tag queries
	if strings.Contains(normalizedQuery, "etiket:") || strings.Contains(normalizedQuery, "tag:") {
		// Extract tag name
		parts := strings.Split(normalizedQuery, ":")
		if len(parts) > 1 {
			tagName := strings.TrimSpace(parts[1])
			// Filter by tag - we need to filter the results manually
			allTasks, err := acy.veriYonetici.GorevleriGetir("", "", "")
			if err != nil {
				return nil, err
			}
			var taggedTasks []*Gorev
			for _, task := range allTasks {
				for _, tag := range task.Etiketler {
					if strings.EqualFold(tag.Isim, tagName) {
						taggedTasks = append(taggedTasks, task)
						break
					}
				}
			}
			return taggedTasks, nil
		}
	}

	// Check for project-specific queries
	if strings.Contains(normalizedQuery, "proje:") || strings.Contains(normalizedQuery, "project:") {
		// This would need project name search functionality
		// For now, return empty
		return []*Gorev{}, nil
	}

	// Default: search in task titles and descriptions
	allTasks, err := acy.veriYonetici.GorevleriGetir("", "", "")
	if err != nil {
		return nil, err
	}

	var matchedTasks []*Gorev
	searchTerms := strings.Fields(normalizedQuery)

	for _, task := range allTasks {
		taskText := strings.ToLower(task.Baslik + " " + task.Aciklama)
		matched := true

		// Check if all search terms are present
		for _, term := range searchTerms {
			if !strings.Contains(taskText, term) {
				matched = false
				break
			}
		}

		if matched {
			matchedTasks = append(matchedTasks, task)
		}
	}

	return matchedTasks, nil
}

// Helper functions for NLP queries

func (acy *AIContextYonetici) getTodayInteractions() ([]*AIInteraction, error) {
	return acy.veriYonetici.AITodayInteractionlariGetir()
}

func (acy *AIContextYonetici) getTasksFromInteractions(interactions []*AIInteraction) ([]*Gorev, error) {
	seen := make(map[string]bool)
	var tasks []*Gorev

	for _, interaction := range interactions {
		if !seen[interaction.GorevID] {
			task, err := acy.veriYonetici.GorevGetir(interaction.GorevID)
			if err == nil {
				tasks = append(tasks, task)
				seen[interaction.GorevID] = true
			}
		}
	}

	return tasks, nil
}

func (acy *AIContextYonetici) getLastCreatedTasks(limit int) ([]*Gorev, error) {
	// Get all tasks sorted by creation date
	allTasks, err := acy.veriYonetici.GorevleriGetir("", "", "")
	if err != nil {
		return nil, err
	}

	// Sort by creation date (newest first)
	sort.Slice(allTasks, func(i, j int) bool {
		return allTasks[i].OlusturmaTarih.After(allTasks[j].OlusturmaTarih)
	})

	// Return requested number of tasks
	if limit > len(allTasks) {
		limit = len(allTasks)
	}

	return allTasks[:limit], nil
}
