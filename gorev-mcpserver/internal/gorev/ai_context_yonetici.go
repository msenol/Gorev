package gorev

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/msenol/gorev/internal/constants"
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
	veriYonetici     VeriYoneticiInterface
	autoStateManager *AutoStateManager
	mu               sync.RWMutex // Protects concurrent access to context operations
}

// YeniAIContextYonetici creates a new AI context manager
func YeniAIContextYonetici(vy VeriYoneticiInterface) *AIContextYonetici {
	return &AIContextYonetici{
		veriYonetici: vy,
	}
}

// SetAutoStateManager sets the auto state manager for enhanced integration
func (acy *AIContextYonetici) SetAutoStateManager(asm *AutoStateManager) {
	acy.autoStateManager = asm
}

// SetActiveTask sets the active task for the AI session
func (acy *AIContextYonetici) SetActiveTask(ctx context.Context, taskID string) error {
	acy.mu.Lock()
	defer acy.mu.Unlock()

	// Validate task exists
	gorev, err := acy.veriYonetici.GorevGetir(ctx, taskID)
	if err != nil {
		return errors.New(i18n.TEntityNotFound(i18n.FromContext(ctx), "task", err))
	}

	// Get current context
	context, err := acy.getContextUnsafe()
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
		if len(context.RecentTasks) > constants.MaxRecentTasks {
			context.RecentTasks = context.RecentTasks[:constants.MaxRecentTasks]
		}
	}

	// Save context
	if err := acy.saveContextUnsafe(context); err != nil {
		return errors.New(i18n.T("error.contextSaveFailed", map[string]interface{}{"Error": err}))
	}

	// Record interaction
	if err := acy.recordInteraction(ctx, taskID, "set_active", nil); err != nil {
		return errors.New(i18n.T("error.interactionSaveFailed", map[string]interface{}{"Error": err}))
	}

	// Auto-transition to constants.TaskStatusInProgress if task is in constants.TaskStatusPending
	if gorev.Status == constants.TaskStatusPending {
		gorev.Status = constants.TaskStatusInProgress
		if err := acy.veriYonetici.GorevGuncelle(ctx, gorev.ID, map[string]interface{}{"status": constants.TaskStatusInProgress}); err != nil {
			return errors.New(i18n.T("error.statusUpdateFailed", map[string]interface{}{"Error": err}))
		}
	}

	return nil
}

// GetActiveTask returns the current active task
func (acy *AIContextYonetici) GetActiveTask(ctx context.Context) (*Gorev, error) {
	acy.mu.RLock()
	defer acy.mu.RUnlock()

	context, err := acy.getContextUnsafe()
	if err != nil {
		return nil, err
	}

	if context.ActiveTaskID == "" {
		return nil, nil
	}

	return acy.veriYonetici.GorevGetir(ctx, context.ActiveTaskID)
}

// GetRecentTasks returns the recent tasks interacted with
func (acy *AIContextYonetici) GetRecentTasks(ctx context.Context, limit int) ([]*Gorev, error) {
	acy.mu.RLock()
	defer acy.mu.RUnlock()

	context, err := acy.getContextUnsafe()
	if err != nil {
		return nil, err
	}

	if limit <= 0 || limit > len(context.RecentTasks) {
		limit = len(context.RecentTasks)
	}

	tasks := make([]*Gorev, 0, limit)
	for i := 0; i < limit; i++ {
		gorev, err := acy.veriYonetici.GorevGetir(ctx, context.RecentTasks[i])
		if err != nil {
			continue // Skip if task not found
		}
		tasks = append(tasks, gorev)
	}

	return tasks, nil
}

// GetContextSummary returns an AI-optimized summary of the current context
func (acy *AIContextYonetici) GetContextSummary(ctx context.Context) (*AIContextSummary, error) {
	summary := &AIContextSummary{}

	// Get active task
	activeTask, _ := acy.GetActiveTask(ctx)
	summary.ActiveTask = activeTask

	// Get recent tasks
	recentTasks, _ := acy.GetRecentTasks(ctx, constants.DefaultRecentTaskLimit)
	summary.RecentTasks = recentTasks

	// Get working project
	if activeTask != nil && activeTask.ProjeID != "" {
		proje, _ := acy.veriYonetici.ProjeGetir(ctx, activeTask.ProjeID)
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
	gorevler, _ := acy.veriYonetici.GorevleriGetir(ctx, constants.TaskStatusPending, "", "")
	for _, g := range gorevler {
		if g.Priority == constants.PriorityHigh {
			summary.NextPriorities = append(summary.NextPriorities, g)
			if len(summary.NextPriorities) >= constants.MaxSummaryItems {
				break
			}
		}
	}

	// Get blockers (tasks with unfinished dependencies)
	for _, g := range gorevler {
		if g.UncompletedDependencyCount > 0 {
			summary.Blockers = append(summary.Blockers, g)
			if len(summary.Blockers) >= constants.MaxSummaryItems {
				break
			}
		}
	}

	return summary, nil
}

// RecordTaskView records when a task is viewed and auto-transitions state
func (acy *AIContextYonetici) RecordTaskView(ctx context.Context, taskID string) error {
	// Record interaction
	if err := acy.recordInteraction(ctx, taskID, "viewed", nil); err != nil {
		return err
	}

	// Get task
	gorev, err := acy.veriYonetici.GorevGetir(ctx, taskID)
	if err != nil {
		return err
	}

	// Auto-transition to constants.TaskStatusInProgress if in constants.TaskStatusPending
	if gorev.Status == constants.TaskStatusPending {
		gorev.Status = constants.TaskStatusInProgress
		if err := acy.veriYonetici.GorevGuncelle(ctx, gorev.ID, map[string]interface{}{"status": constants.TaskStatusInProgress}); err != nil {
			return errors.New(i18n.T("error.autoStatusUpdateFailed", map[string]interface{}{"Error": err}))
		}
		// Record the state change
		if err := acy.recordInteraction(ctx, taskID, "updated", map[string]interface{}{
			"auto_state_change": true,
			"from":              constants.TaskStatusPending,
			"to":                constants.TaskStatusInProgress,
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
	acy.mu.RLock()
	defer acy.mu.RUnlock()
	return acy.getContextUnsafe()
}

func (acy *AIContextYonetici) saveContext(context *AIContext) error {
	acy.mu.Lock()
	defer acy.mu.Unlock()
	return acy.saveContextUnsafe(context)
}

// getContextUnsafe is the internal method without mutex protection
func (acy *AIContextYonetici) getContextUnsafe() (*AIContext, error) {
	return acy.veriYonetici.AIContextGetir()
}

// saveContextUnsafe is the internal method without mutex protection
func (acy *AIContextYonetici) saveContextUnsafe(context *AIContext) error {
	return acy.veriYonetici.AIContextKaydet(context)
}

func (acy *AIContextYonetici) RecordInteraction(ctx context.Context, taskID, actionType string, context interface{}) error {
	return acy.recordInteraction(ctx, taskID, actionType, context)
}

// recordInteraction is the internal method for recording interactions
func (acy *AIContextYonetici) recordInteraction(ctx context.Context, taskID, actionType string, context interface{}) error {
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
	return acy.veriYonetici.AIInteractionlariGetir(constants.MaxInteractionHistory) // Get last 50 interactions for session
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
	if len(context.RecentTasks) > constants.MaxRecentTasks {
		context.RecentTasks = context.RecentTasks[:constants.MaxRecentTasks]
	}

	// Save updated context
	return acy.saveContext(context)
}

// BatchUpdate represents a single update in a batch operation
type BatchUpdate struct {
	ID      string                 `json:"id"`
	Updates map[string]interface{} `json:"updates"`
}

// BatchUpdate performs multiple task updates in a single operation
func (acy *AIContextYonetici) BatchUpdate(ctx context.Context, updates []BatchUpdate) (*BatchUpdateResult, error) {
	result := &BatchUpdateResult{
		Successful: []string{},
		Failed:     []BatchUpdateError{},
	}

	for _, update := range updates {
		// Validate task exists
		_, err := acy.veriYonetici.GorevGetir(ctx, update.ID)
		if err != nil {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: update.ID,
				Error:  i18n.TEntityNotFound(i18n.FromContext(ctx), "task", err),
			})
			continue
		}

		// Apply updates based on fields
		updateFields := make(map[string]interface{})

		// Validate and collect all supported field updates
		if durum, ok := update.Updates["status"].(string); ok {
			// Validate status values
			validStatuses := constants.GetValidTaskStatuses()
			isValid := false
			for _, status := range validStatuses {
				if durum == status {
					isValid = true
					break
				}
			}
			if !isValid {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: update.ID,
					Error:  i18n.T("error.invalidStatusBatch", map[string]interface{}{"Status": durum}),
				})
				continue
			}
			updateFields["status"] = durum
		}

		if oncelik, ok := update.Updates["priority"].(string); ok {
			// Validate priority values
			validPriorities := constants.GetValidPriorities()
			isValid := false
			for _, priority := range validPriorities {
				if oncelik == priority {
					isValid = true
					break
				}
			}
			if !isValid {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: update.ID,
					Error:  i18n.T("error.invalidPriorityBatch", map[string]interface{}{"Priority": oncelik}),
				})
				continue
			}
			updateFields["priority"] = oncelik
		}

		if baslik, ok := update.Updates["title"].(string); ok {
			if strings.TrimSpace(baslik) == "" {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: update.ID,
					Error:  i18n.T("error.titleCannotBeEmpty"),
				})
				continue
			}
			updateFields["title"] = baslik
		}

		if aciklama, ok := update.Updates["description"].(string); ok {
			updateFields["description"] = aciklama
		}

		if sonTarih, ok := update.Updates["due_date"].(string); ok {
			if sonTarih != "" {
				// Validate date format (YYYY-MM-DD)
				if _, err := time.Parse("2006-01-02", sonTarih); err != nil {
					result.Failed = append(result.Failed, BatchUpdateError{
						TaskID: update.ID,
						Error:  i18n.T("error.invalidDateFormatBatch", map[string]interface{}{"Date": sonTarih}),
					})
					continue
				}
			}
			updateFields["due_date"] = sonTarih
		}

		// Apply all validated updates at once
		if len(updateFields) > 0 {
			if err := acy.veriYonetici.GorevGuncelle(ctx, update.ID, updateFields); err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: update.ID,
					Error:  i18n.T("error.taskUpdateError", map[string]interface{}{"Error": err}),
				})
				continue
			}
		}

		result.Successful = append(result.Successful, update.ID)

		// Record batch operation
		if err := acy.recordInteraction(ctx, update.ID, "bulk_operation", update.Updates); err != nil {
			// Log but don't fail the operation
			// fmt.Printf("interaction kaydetme hatası: %v\n", err)
		}
	}

	result.TotalProcessed = len(updates)
	return result, nil
}

// NLPQuery performs natural language query on tasks
func (acy *AIContextYonetici) NLPQuery(ctx context.Context, query string) ([]*Gorev, error) {
	// Use enhanced NLP processing if auto state manager is available
	if acy.autoStateManager != nil {
		result, err := acy.autoStateManager.ProcessNaturalLanguageQuery(ctx, query, "")
		if err != nil {
			// Fallback to basic NLP processing
			return acy.basicNLPQuery(ctx, query)
		}

		// Extract tasks from the structured result
		if resultMap, ok := result.(map[string]interface{}); ok {
			if tasksResult, ok := resultMap["result"]; ok {
				if tasks, ok := tasksResult.([]*Gorev); ok {
					return tasks, nil
				}
			}
		}
	}

	// Fallback to basic NLP processing
	return acy.basicNLPQuery(ctx, query)
}

// basicNLPQuery performs basic natural language query processing
func (acy *AIContextYonetici) basicNLPQuery(ctx context.Context, query string) ([]*Gorev, error) {
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
			return acy.getTasksFromInteractions(ctx, interactions)
		},
		"son oluşturduğum": func() ([]*Gorev, error) {
			// Last created task
			return acy.getLastCreatedTasks(ctx, constants.LastCreatedCount)
		},
		"son oluşturulan": func() ([]*Gorev, error) {
			// Recently created tasks
			return acy.getLastCreatedTasks(ctx, constants.RecentlyCreatedCount)
		},
		"yüksek öncelik": func() ([]*Gorev, error) {
			// High priority tasks
			return acy.veriYonetici.GorevleriGetir(ctx, constants.TaskStatusPending, "", "")
		},
		"tamamlanmamış": func() ([]*Gorev, error) {
			// Incomplete tasks
			return acy.veriYonetici.GorevleriGetir(ctx, constants.TaskStatusPending, "", "")
		},
		"devam eden": func() ([]*Gorev, error) {
			// In progress tasks
			return acy.veriYonetici.GorevleriGetir(ctx, constants.TaskStatusInProgress, "", "")
		},
		"tamamlanan": func() ([]*Gorev, error) {
			// Completed tasks
			return acy.veriYonetici.GorevleriGetir(ctx, constants.TaskStatusCompleted, "", "")
		},
		"blokaj": func() ([]*Gorev, error) {
			// Blocked tasks
			gorevler, _ := acy.veriYonetici.GorevleriGetir(ctx, "", "", "")
			var blocked []*Gorev
			for _, g := range gorevler {
				if g.UncompletedDependencyCount > 0 {
					blocked = append(blocked, g)
				}
			}
			return blocked, nil
		},
		"acil": func() ([]*Gorev, error) {
			// Urgent tasks (due soon)
			return acy.veriYonetici.GorevleriGetir(ctx, "", "", "acil")
		},
		"gecikmiş": func() ([]*Gorev, error) {
			// Overdue tasks
			return acy.veriYonetici.GorevleriGetir(ctx, "", "", "gecmis")
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
			allTasks, err := acy.veriYonetici.GorevleriGetir(ctx, "", "", "")
			if err != nil {
				return nil, err
			}
			var taggedTasks []*Gorev
			for _, task := range allTasks {
				for _, tag := range task.Tags {
					if strings.EqualFold(tag.Name, tagName) {
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
	allTasks, err := acy.veriYonetici.GorevleriGetir(ctx, "", "", "")
	if err != nil {
		return nil, err
	}

	var matchedTasks []*Gorev
	searchTerms := strings.Fields(normalizedQuery)

	for _, task := range allTasks {
		taskText := strings.ToLower(task.Title + " " + task.Description)
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

func (acy *AIContextYonetici) getTasksFromInteractions(ctx context.Context, interactions []*AIInteraction) ([]*Gorev, error) {
	seen := make(map[string]bool)
	var tasks []*Gorev

	for _, interaction := range interactions {
		if !seen[interaction.GorevID] {
			task, err := acy.veriYonetici.GorevGetir(ctx, interaction.GorevID)
			if err == nil {
				tasks = append(tasks, task)
				seen[interaction.GorevID] = true
			}
		}
	}

	return tasks, nil
}

func (acy *AIContextYonetici) getLastCreatedTasks(ctx context.Context, limit int) ([]*Gorev, error) {
	// Get all tasks sorted by creation date
	allTasks, err := acy.veriYonetici.GorevleriGetir(ctx, "", "", "")
	if err != nil {
		return nil, err
	}

	// Sort by creation date (newest first)
	sort.Slice(allTasks, func(i, j int) bool {
		return allTasks[i].CreatedAt.After(allTasks[j].CreatedAt)
	})

	// Return requested number of tasks
	if limit > len(allTasks) {
		limit = len(allTasks)
	}

	return allTasks[:limit], nil
}
