package gorev

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// AutoStateManager handles automatic task state transitions
type AutoStateManager struct {
	veriYonetici     VeriYoneticiInterface
	inactivityTimer  time.Duration
	activeTimers     map[string]*time.Timer
	aiContextManager *AIContextYonetici
	nlpProcessor     *NLPProcessor
}

// InactivityConfig represents configuration for inactivity detection
type InactivityConfig struct {
	Duration time.Duration
	Enabled  bool
}

// YeniAutoStateManager creates a new auto state manager
func YeniAutoStateManager(vy VeriYoneticiInterface) *AutoStateManager {
	return &AutoStateManager{
		veriYonetici:    vy,
		inactivityTimer: 30 * time.Minute, // Default 30 minutes
		activeTimers:    make(map[string]*time.Timer),
		nlpProcessor:    NewNLPProcessor(),
	}
}

// SetAIContextManager sets the AI context manager for integration
func (asm *AutoStateManager) SetAIContextManager(acm *AIContextYonetici) {
	asm.aiContextManager = acm
}

// AutoTransitionToInProgress automatically transitions a task to in-progress when accessed
func (asm *AutoStateManager) AutoTransitionToInProgress(taskID string) error {
	log.Printf("Debug: Auto-transitioning task to in-progress, taskID: %s", taskID)

	// Get current task
	task, err := asm.veriYonetici.GorevDetay(taskID)
	if err != nil {
		return err
	}

	// Only transition if currently pending
	if task.Status != constants.TaskStatusPending {
		log.Printf("Task not in pending status, skipping auto-transition: taskID=%s, currentStatus=%s", taskID, task.Status)
		return nil
	}

	// Check dependencies before transitioning
	canStart, err := asm.checkDependenciesCompleted(taskID)
	if err != nil {
		return err
	}

	if !canStart {
		log.Printf("Task has incomplete dependencies, cannot auto-start: taskID=%s", taskID)
		return nil
	}

	// Transition to in-progress
	err = asm.veriYonetici.GorevGuncelle(taskID, map[string]interface{}{"status": constants.TaskStatusInProgress})
	if err != nil {
		return err
	}

	// Record the interaction
	if asm.aiContextManager != nil {
		if recErr := asm.aiContextManager.recordInteraction(taskID, "auto_transition_start", map[string]interface{}{
			"from_status": constants.TaskStatusPending,
			"to_status":   constants.TaskStatusInProgress,
			"reason":      "task_accessed",
			"timestamp":   time.Now(),
		}); recErr != nil {
			log.Printf("Failed to record AI interaction (auto_transition_start): taskID=%s, error=%v", taskID, recErr)
		}
	}

	// Start inactivity timer
	asm.startInactivityTimer(taskID)

	log.Printf("Auto-transitioned task to in-progress: taskID=%s, reason=task_accessed", taskID)

	return nil
}

// AutoTransitionToPending automatically transitions a task back to pending after inactivity
func (asm *AutoStateManager) AutoTransitionToPending(taskID string) error {
	log.Printf("Auto-transitioning task to pending due to inactivity: taskID=%s", taskID)

	// Get current task
	task, err := asm.veriYonetici.GorevDetay(taskID)
	if err != nil {
		return err
	}

	// Only transition if currently in progress
	if task.Status != constants.TaskStatusInProgress {
		log.Printf("Task not in progress, skipping auto-transition: taskID=%s, currentStatus=%s", taskID, task.Status)
		return nil
	}

	// Transition back to pending
	err = asm.veriYonetici.GorevGuncelle(taskID, map[string]interface{}{"status": constants.TaskStatusPending})
	if err != nil {
		return err
	}

	// Record the interaction
	if asm.aiContextManager != nil {
		if recErr := asm.aiContextManager.recordInteraction(taskID, "auto_transition_pause", map[string]interface{}{
			"from_status": constants.TaskStatusInProgress,
			"to_status":   constants.TaskStatusPending,
			"reason":      "inactivity_timeout",
			"timeout":     asm.inactivityTimer.String(),
			"timestamp":   time.Now(),
		}); recErr != nil {
			log.Printf("Failed to record AI interaction (auto_transition_pause): taskID=%s, error=%v", taskID, recErr)
		}
	}

	// Clear inactivity timer
	asm.clearInactivityTimer(taskID)

	log.Printf("Auto-transitioned task to pending due to inactivity: taskID=%s, timeout=%v", taskID, asm.inactivityTimer)

	return nil
}

// CheckParentCompletion checks if a parent task can be completed based on subtask completion
func (asm *AutoStateManager) CheckParentCompletion(taskID string) error {
	log.Printf("Checking parent completion eligibility: taskID=%s", taskID)

	// Get task details
	task, err := asm.veriYonetici.GorevDetay(taskID)
	if err != nil {
		return err
	}

	// If task has no parent, nothing to check
	if task.ParentID == "" {
		return nil
	}

	parentID := task.ParentID

	// Get all subtasks of the parent
	subtasks, err := asm.getSubtasks(parentID)
	if err != nil {
		return err
	}

	// Check if all subtasks are completed
	allCompleted := true
	for _, subtask := range subtasks {
		if subtask.Status != constants.TaskStatusCompleted {
			allCompleted = false
			break
		}
	}

	if allCompleted {
		// Get parent task
		parentTask, err := asm.veriYonetici.GorevDetay(parentID)
		if err != nil {
			return err
		}

		// Auto-complete parent if not already completed
		if parentTask.Status != constants.TaskStatusCompleted {
			err = asm.veriYonetici.GorevGuncelle(parentID, map[string]interface{}{"status": constants.TaskStatusCompleted})
			if err != nil {
				return err
			}

			// Record the interaction
			if asm.aiContextManager != nil {
				if recErr := asm.aiContextManager.recordInteraction(parentID, "auto_complete_parent", map[string]interface{}{
					"reason":        "all_subtasks_completed",
					"subtask_count": len(subtasks),
					"timestamp":     time.Now(),
				}); recErr != nil {
					log.Printf("Failed to record AI interaction (auto_complete_parent): parentID=%s, error=%v", parentID, recErr)
				}
			}

			log.Printf("Auto-completed parent task: parentID=%s, subtaskCount=%d, reason=all_subtasks_completed", parentID, len(subtasks))

			// Recursively check grandparent
			return asm.CheckParentCompletion(parentID)
		}
	}

	return nil
}

// ScheduleInactivityCheck starts a timer for inactivity detection
func (asm *AutoStateManager) ScheduleInactivityCheck(taskID string) {
	asm.startInactivityTimer(taskID)
}

// ResetInactivityTimer resets the inactivity timer for a task
func (asm *AutoStateManager) ResetInactivityTimer(taskID string) {
	asm.clearInactivityTimer(taskID)
	asm.startInactivityTimer(taskID)
}

// OnTaskAccessed should be called whenever a task is accessed by AI
func (asm *AutoStateManager) OnTaskAccessed(taskID string) error {
	// Auto-transition to in-progress if needed
	err := asm.AutoTransitionToInProgress(taskID)
	if err != nil {
		return err
	}

	// Reset inactivity timer
	asm.ResetInactivityTimer(taskID)

	return nil
}

// OnTaskCompleted should be called when a task is completed
func (asm *AutoStateManager) OnTaskCompleted(taskID string) error {
	// Clear inactivity timer
	asm.clearInactivityTimer(taskID)

	// Check if parent can be completed
	return asm.CheckParentCompletion(taskID)
}

// SetInactivityDuration sets the inactivity timeout duration
func (asm *AutoStateManager) SetInactivityDuration(duration time.Duration) {
	asm.inactivityTimer = duration
}

// GetInactivityDuration returns the current inactivity timeout duration
func (asm *AutoStateManager) GetInactivityDuration() time.Duration {
	return asm.inactivityTimer
}

// startInactivityTimer starts the inactivity timer for a task
func (asm *AutoStateManager) startInactivityTimer(taskID string) {
	// Clear existing timer if any
	asm.clearInactivityTimer(taskID)

	// Create new timer
	timer := time.AfterFunc(asm.inactivityTimer, func() {
		err := asm.AutoTransitionToPending(taskID)
		if err != nil {
			log.Printf("Failed to auto-transition task to pending: taskID=%s, error=%v", taskID, err)
		}
	})

	asm.activeTimers[taskID] = timer
	log.Printf("Started inactivity timer: taskID=%s, duration=%v", taskID, asm.inactivityTimer)
}

// clearInactivityTimer clears the inactivity timer for a task
func (asm *AutoStateManager) clearInactivityTimer(taskID string) {
	if timer, exists := asm.activeTimers[taskID]; exists {
		timer.Stop()
		delete(asm.activeTimers, taskID)
		log.Printf("Cleared inactivity timer: taskID=%s", taskID)
	}
}

// checkDependenciesCompleted checks if all dependencies for a task are completed
func (asm *AutoStateManager) checkDependenciesCompleted(taskID string) (bool, error) {
	dependencies, err := asm.veriYonetici.GorevBagimlilikGetir(taskID)
	if err != nil {
		return false, err
	}

	for _, dep := range dependencies {
		if dep.Status != constants.TaskStatusCompleted {
			return false, nil
		}
	}

	return true, nil
}

// getSubtasks returns all subtasks for a given parent task
func (asm *AutoStateManager) getSubtasks(parentID string) ([]*Gorev, error) {
	// Query all tasks and filter by ParentID
	allTasks, err := asm.veriYonetici.GorevleriGetir("", "", "")
	if err != nil {
		return nil, err
	}

	var subtasks []*Gorev
	for _, task := range allTasks {
		if task.ParentID == parentID {
			subtasks = append(subtasks, task)
		}
	}

	return subtasks, nil
}

// Cleanup stops all active timers
func (asm *AutoStateManager) Cleanup() {
	for taskID, timer := range asm.activeTimers {
		timer.Stop()
		log.Printf("Stopped timer during cleanup: taskID=%s", taskID)
	}
	asm.activeTimers = make(map[string]*time.Timer)
}

// ProcessNaturalLanguageQuery processes natural language queries and executes corresponding actions
func (asm *AutoStateManager) ProcessNaturalLanguageQuery(query string, lang string) (interface{}, error) {
	log.Printf("Processing natural language query: query=%s, lang=%s", query, lang)

	// Parse the query using NLP processor
	intent, err := asm.nlpProcessor.ProcessQuery(query)
	if err != nil {
		return nil, err
	}

	// Validate the intent
	if err := asm.nlpProcessor.ValidateIntent(intent); err != nil {
		return nil, err
	}

	// Record the query in AI context
	if asm.aiContextManager != nil {
		if recErr := asm.aiContextManager.recordInteraction("system", "nlp_query", map[string]interface{}{
			"query":      query,
			"intent":     intent,
			"confidence": intent.Confidence,
			"timestamp":  time.Now(),
		}); recErr != nil {
			log.Printf("Failed to record AI interaction (nlp_query): error=%v", recErr)
		}
	}

	// Execute the action based on intent
	result, err := asm.executeAction(intent)
	if err != nil {
		return nil, err
	}

	// Format natural language response
	response := asm.nlpProcessor.FormatResponse(intent.Action, result, lang)

	log.Printf("Natural language query processed: query=%s, action=%s, confidence=%f", query, intent.Action, intent.Confidence)

	return map[string]interface{}{
		"response": response,
		"intent":   intent,
		"result":   result,
	}, nil
}

// executeAction executes the parsed action from NLP intent
func (asm *AutoStateManager) executeAction(intent *QueryIntent) (interface{}, error) {
	switch intent.Action {
	case "list":
		return asm.executeListAction(intent)
	case "create":
		return asm.executeCreateAction(intent)
	case "update":
		return asm.executeUpdateAction(intent)
	case "complete":
		return asm.executeCompleteAction(intent)
	case "delete":
		return asm.executeDeleteAction(intent)
	case "search":
		return asm.executeSearchAction(intent)
	case "status":
		return asm.executeStatusAction(intent)
	default:
		return nil, fmt.Errorf(i18n.T("error.unsupportedAction", map[string]interface{}{"Action": intent.Action}))
	}
}

// executeListAction lists tasks based on filters
func (asm *AutoStateManager) executeListAction(intent *QueryIntent) (interface{}, error) {
	// Build query parameters from intent
	filters := make(map[string]interface{})

	// Apply filters from intent
	for key, value := range intent.Filters {
		filters[key] = value
	}

	// Apply time range if specified
	if intent.TimeRange != nil {
		if intent.TimeRange.Start != nil {
			filters["created_after"] = intent.TimeRange.Start
		}
		if intent.TimeRange.End != nil {
			filters["created_before"] = intent.TimeRange.End
		}
	}

	// Get tasks from data manager
	tasks, err := asm.veriYonetici.GorevListele(filters)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// executeCreateAction creates a new task from natural language
func (asm *AutoStateManager) executeCreateAction(intent *QueryIntent) (interface{}, error) {
	// Extract task content from the query
	content := asm.nlpProcessor.ExtractTaskContent(intent.Raw)

	// Validate required fields
	title, ok := content["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf(i18n.T("error.taskTitleRequired"))
	}

	// Create task parameters
	taskParams := map[string]interface{}{
		"title":      title,
		"status":     constants.TaskStatusPending,
		"created_at": time.Now(),
	}

	// Add optional fields
	if desc, ok := content["description"].(string); ok && desc != "" {
		taskParams["description"] = desc
	}

	if dueDate, ok := content["due_date"].(string); ok && dueDate != "" {
		taskParams["due_date"] = dueDate
	}

	// Apply filters as task properties
	for key, value := range intent.Filters {
		switch key {
		case "priority":
			taskParams["priority"] = value
		case "category":
			taskParams["category"] = value
		case "tags":
			taskParams["tags"] = value
		}
	}

	// Create the task
	taskID, err := asm.veriYonetici.GorevOlustur(taskParams)
	if err != nil {
		return nil, err
	}

	// Record the creation
	if asm.aiContextManager != nil {
		if recErr := asm.aiContextManager.recordInteraction(taskID, "nlp_create", map[string]interface{}{
			"original_query":    intent.Raw,
			"extracted_content": content,
			"timestamp":         time.Now(),
		}); recErr != nil {
			log.Printf("Failed to record AI interaction (nlp_create): taskID=%s, error=%v", taskID, recErr)
		}
	}

	return title, nil
}

// executeUpdateAction updates an existing task
func (asm *AutoStateManager) executeUpdateAction(intent *QueryIntent) (interface{}, error) {
	refs, ok := intent.Parameters["task_references"].([]string)
	if !ok || len(refs) == 0 {
		return nil, fmt.Errorf(i18n.T("error.taskReferenceRequiredUpdate"))
	}

	// For now, handle the first reference
	taskID, err := asm.resolveTaskReference(refs[0])
	if err != nil {
		return nil, err
	}

	// Get current task
	task, err := asm.veriYonetici.GorevDetay(taskID)
	if err != nil {
		return nil, err
	}

	// Extract update content
	content := asm.nlpProcessor.ExtractTaskContent(intent.Raw)
	updateParams := make(map[string]interface{})

	// Apply updates
	if title, ok := content["title"].(string); ok && title != "" {
		updateParams["title"] = title
	}

	if desc, ok := content["description"].(string); ok && desc != "" {
		updateParams["description"] = desc
	}

	// Apply filter changes as updates
	for key, value := range intent.Filters {
		switch key {
		case "priority":
			updateParams["priority"] = value
		case "status":
			updateParams["status"] = value
		case "category":
			updateParams["category"] = value
		}
	}

	// Update the task
	err = asm.veriYonetici.GorevGuncelle(taskID, updateParams)
	if err != nil {
		return nil, err
	}

	return task.Title, nil
}

// executeCompleteAction marks a task as completed
func (asm *AutoStateManager) executeCompleteAction(intent *QueryIntent) (interface{}, error) {
	refs, ok := intent.Parameters["task_references"].([]string)
	if !ok || len(refs) == 0 {
		return nil, fmt.Errorf(i18n.T("error.taskReferenceRequiredCompletion"))
	}

	taskID, err := asm.resolveTaskReference(refs[0])
	if err != nil {
		return nil, err
	}

	// Get task details
	task, err := asm.veriYonetici.GorevDetay(taskID)
	if err != nil {
		return nil, err
	}

	// Complete the task
	err = asm.veriYonetici.GorevGuncelle(taskID, map[string]interface{}{"status": constants.TaskStatusCompleted})
	if err != nil {
		return nil, err
	}

	// Trigger auto-completion check for parent
	err = asm.OnTaskCompleted(taskID)
	if err != nil {
		log.Printf("Failed to check parent completion: taskID=%s, error=%v", taskID, err)
	}

	return task.Title, nil
}

// executeDeleteAction deletes a task
func (asm *AutoStateManager) executeDeleteAction(intent *QueryIntent) (interface{}, error) {
	refs, ok := intent.Parameters["task_references"].([]string)
	if !ok || len(refs) == 0 {
		return nil, fmt.Errorf(i18n.T("error.taskReferenceRequiredDeletion"))
	}

	taskID, err := asm.resolveTaskReference(refs[0])
	if err != nil {
		return nil, err
	}

	// Get task details before deletion
	task, err := asm.veriYonetici.GorevDetay(taskID)
	if err != nil {
		return nil, err
	}

	// Delete the task
	err = asm.veriYonetici.GorevSil(taskID)
	if err != nil {
		return nil, err
	}

	// Clear any active timers
	asm.clearInactivityTimer(taskID)

	return task.Title, nil
}

// executeSearchAction searches for tasks
func (asm *AutoStateManager) executeSearchAction(intent *QueryIntent) (interface{}, error) {
	// Use the list action with search-specific filters
	return asm.executeListAction(intent)
}

// executeStatusAction shows status of specific tasks
func (asm *AutoStateManager) executeStatusAction(intent *QueryIntent) (interface{}, error) {
	if refs, ok := intent.Parameters["task_references"].([]string); ok && len(refs) > 0 {
		taskID, err := asm.resolveTaskReference(refs[0])
		if err != nil {
			return nil, err
		}

		task, err := asm.veriYonetici.GorevDetay(taskID)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"task":   task,
			"status": task.Status,
		}, nil
	}

	// Return general status
	return asm.executeListAction(intent)
}

// resolveTaskReference resolves a task reference to a task ID
func (asm *AutoStateManager) resolveTaskReference(ref string) (string, error) {
	if strings.HasPrefix(ref, "id:") {
		return strings.TrimPrefix(ref, "id:"), nil
	}

	if strings.HasPrefix(ref, "title:") {
		title := strings.TrimPrefix(ref, "title:")
		// Search for task by title - this would need implementation in VeriYonetici
		tasks, err := asm.veriYonetici.GorevListele(map[string]interface{}{
			"title_search": title,
			"limit":        1,
		})
		if err != nil {
			return "", err
		}
		if len(tasks) == 0 {
			return "", fmt.Errorf(i18n.T("error.noTaskFoundWithTitle", map[string]interface{}{"Title": title}))
		}
		return tasks[0].ID, nil
	}

	if strings.HasPrefix(ref, "recent:") {
		countStr := strings.TrimPrefix(ref, "recent:")
		count := 1
		if c, err := strconv.Atoi(countStr); err == nil {
			count = c
		}

		// Get recent tasks
		tasks, err := asm.veriYonetici.GorevListele(map[string]interface{}{
			"order_by": "created_desc",
			"limit":    count,
		})
		if err != nil {
			return "", err
		}
		if len(tasks) == 0 {
			return "", fmt.Errorf(i18n.T("error.noRecentTasksFound"))
		}
		return tasks[0].ID, nil
	}

	// Default: treat as direct task ID
	return ref, nil
}
