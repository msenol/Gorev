package gorev

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// BatchProcessor handles batch operations on multiple tasks
type BatchProcessor struct {
	veriYonetici     VeriYoneticiInterface
	aiContextManager *AIContextYonetici
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(vy VeriYoneticiInterface) *BatchProcessor {
	return &BatchProcessor{
		veriYonetici: vy,
	}
}

// SetAIContextManager sets the AI context manager for interaction tracking
func (bp *BatchProcessor) SetAIContextManager(acm *AIContextYonetici) {
	bp.aiContextManager = acm
}

// BatchUpdateRequest represents a single update in a batch operation
type BatchUpdateRequest struct {
	TaskID  string                 `json:"task_id"`
	Updates map[string]interface{} `json:"updates"`
	DryRun  bool                   `json:"dry_run,omitempty"`
}

// BatchUpdateResult represents the result of a batch update operation
type BatchUpdateResult struct {
	Successful     []string             `json:"successful"`
	Failed         []BatchUpdateError   `json:"failed"`
	Warnings       []BatchUpdateWarning `json:"warnings"`
	TotalProcessed int                  `json:"total_processed"`
	ExecutionTime  time.Duration        `json:"execution_time"`
	Summary        string               `json:"summary"`
}

// BatchUpdateError represents a failed update in a batch operation
type BatchUpdateError struct {
	TaskID string `json:"task_id"`
	Error  string `json:"error"`
	Field  string `json:"field,omitempty"`
}

// BatchUpdateWarning represents a warning during batch operation
type BatchUpdateWarning struct {
	TaskID  string `json:"task_id"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// BulkStatusTransitionRequest represents a bulk status change request
type BulkStatusTransitionRequest struct {
	TaskIDs           []string `json:"task_ids"`
	NewStatus         string   `json:"new_status"`
	Force             bool     `json:"force,omitempty"`
	CheckDependencies bool     `json:"check_dependencies,omitempty"`
	DryRun            bool     `json:"dry_run,omitempty"`
}

// BulkTagOperationRequest represents a bulk tag operation request
type BulkTagOperationRequest struct {
	TaskIDs   []string `json:"task_ids"`
	Tags      []string `json:"tags"`
	Operation string   `json:"operation"` // "add", "remove", "replace"
	DryRun    bool     `json:"dry_run,omitempty"`
}

// BulkDeleteRequest represents a bulk delete request
type BulkDeleteRequest struct {
	TaskIDs        []string `json:"task_ids"`
	Confirmation   string   `json:"confirmation"`
	Force          bool     `json:"force,omitempty"`
	DeleteSubtasks bool     `json:"delete_subtasks,omitempty"`
	DryRun         bool     `json:"dry_run,omitempty"`
}

// ProcessBatchUpdate performs multiple task updates in a single transaction
func (bp *BatchProcessor) ProcessBatchUpdate(requests []BatchUpdateRequest) (*BatchUpdateResult, error) {
	startTime := time.Now()

	result := &BatchUpdateResult{
		Successful: []string{},
		Failed:     []BatchUpdateError{},
		Warnings:   []BatchUpdateWarning{},
	}

	log.Printf("Starting batch update operation: count=%d", len(requests))

	for _, request := range requests {
		if request.DryRun {
			// Validate without executing
			if err := bp.validateUpdateRequest(request); err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: request.TaskID,
					Error:  fmt.Sprintf("validation failed: %v", err),
				})
			} else {
				result.Warnings = append(result.Warnings, BatchUpdateWarning{
					TaskID:  request.TaskID,
					Message: "dry run - would be updated",
				})
			}
			continue
		}

		// Validate task exists
		task, err := bp.veriYonetici.GorevDetay(request.TaskID)
		if err != nil {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: request.TaskID,
				Error:  i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}),
			})
			continue
		}

		// Apply updates field by field
		updated := false
		warnings := []string{}

		// Update status
		if newStatus, ok := request.Updates["durum"].(string); ok {
			if bp.validateStatusTransition(task.Durum, newStatus) {
				task.Durum = newStatus
				updated = true
			} else {
				warnings = append(warnings, fmt.Sprintf("invalid status transition: %s -> %s", task.Durum, newStatus))
			}
		}

		// Update priority
		if newPriority, ok := request.Updates["oncelik"].(string); ok {
			if bp.validatePriority(newPriority) {
				task.Oncelik = newPriority
				updated = true
			} else {
				warnings = append(warnings, fmt.Sprintf("invalid priority: %s", newPriority))
			}
		}

		// Update title
		if newTitle, ok := request.Updates["baslik"].(string); ok {
			if strings.TrimSpace(newTitle) != "" {
				task.Baslik = strings.TrimSpace(newTitle)
				updated = true
			} else {
				warnings = append(warnings, "empty title not allowed")
			}
		}

		// Update description
		if newDesc, ok := request.Updates["aciklama"].(string); ok {
			task.Aciklama = newDesc
			updated = true
		}

		// Update due date
		if dueDateStr, ok := request.Updates["son_tarih"].(string); ok {
			if dueDateStr == "" {
				task.SonTarih = nil
				updated = true
			} else if dueDate, err := time.Parse("2006-01-02", dueDateStr); err == nil {
				task.SonTarih = &dueDate
				updated = true
			} else {
				warnings = append(warnings, fmt.Sprintf("invalid date format: %s", dueDateStr))
			}
		}

		// Update tags
		if tagsRaw, ok := request.Updates["etiketler"]; ok {
			if tagNames, ok := tagsRaw.([]string); ok {
				tags, err := bp.veriYonetici.EtiketleriGetirVeyaOlustur(tagNames)
				if err != nil {
					warnings = append(warnings, fmt.Sprintf("tag processing failed: %v", err))
				} else {
					if err := bp.veriYonetici.GorevEtiketleriniAyarla(request.TaskID, tags); err != nil {
						warnings = append(warnings, fmt.Sprintf("tag assignment failed: %v", err))
					} else {
						updated = true
					}
				}
			}
		}

		// Save the task if any updates were made
		if updated {
			updateParams := make(map[string]interface{})
			if newStatus := request.Updates["status"]; newStatus != nil {
				updateParams["durum"] = newStatus
			}
			if newPriority := request.Updates["priority"]; newPriority != nil {
				updateParams["oncelik"] = newPriority
			}

			if err := bp.veriYonetici.GorevGuncelle(request.TaskID, updateParams); err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: request.TaskID,
					Error:  fmt.Sprintf("update failed: %v", err),
				})
				continue
			}

			result.Successful = append(result.Successful, request.TaskID)

			// Record interaction
			if bp.aiContextManager != nil {
				bp.aiContextManager.RecordInteraction(request.TaskID, "batch_update", request.Updates)
			}
		}

		// Add warnings if any
		for _, warning := range warnings {
			result.Warnings = append(result.Warnings, BatchUpdateWarning{
				TaskID:  request.TaskID,
				Message: warning,
			})
		}
	}

	result.TotalProcessed = len(requests)
	result.ExecutionTime = time.Since(startTime)
	result.Summary = fmt.Sprintf("Processed %d tasks: %d successful, %d failed, %d warnings",
		result.TotalProcessed, len(result.Successful), len(result.Failed), len(result.Warnings))

	log.Printf("Batch update completed: total=%d, successful=%d, failed=%d, warnings=%d, duration=%v", result.TotalProcessed, len(result.Successful), len(result.Failed), len(result.Warnings), result.ExecutionTime)

	return result, nil
}

// BulkStatusTransition changes status for multiple tasks
func (bp *BatchProcessor) BulkStatusTransition(request BulkStatusTransitionRequest) (*BatchUpdateResult, error) {
	startTime := time.Now()

	result := &BatchUpdateResult{
		Successful: []string{},
		Failed:     []BatchUpdateError{},
		Warnings:   []BatchUpdateWarning{},
	}

	log.Printf("Starting bulk status transition: count=%d, newStatus=%s", len(request.TaskIDs), request.NewStatus)

	// Validate new status
	if !bp.validateStatus(request.NewStatus) {
		return nil, fmt.Errorf("invalid status: %s", request.NewStatus)
	}

	for _, taskID := range request.TaskIDs {
		if request.DryRun {
			task, err := bp.veriYonetici.GorevDetay(taskID)
			if err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: taskID,
					Error:  i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}),
				})
			} else if !bp.validateStatusTransition(task.Durum, request.NewStatus) {
				result.Warnings = append(result.Warnings, BatchUpdateWarning{
					TaskID:  taskID,
					Message: fmt.Sprintf("invalid transition: %s -> %s", task.Durum, request.NewStatus),
				})
			} else {
				result.Warnings = append(result.Warnings, BatchUpdateWarning{
					TaskID:  taskID,
					Message: "dry run - would be updated",
				})
			}
			continue
		}

		// Get task
		task, err := bp.veriYonetici.GorevDetay(taskID)
		if err != nil {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: taskID,
				Error:  i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}),
			})
			continue
		}

		// Check current status
		if task.Durum == request.NewStatus {
			result.Warnings = append(result.Warnings, BatchUpdateWarning{
				TaskID:  taskID,
				Message: fmt.Sprintf("already in status %s", request.NewStatus),
			})
			continue
		}

		// Validate transition
		if !request.Force && !bp.validateStatusTransition(task.Durum, request.NewStatus) {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: taskID,
				Error:  fmt.Sprintf("invalid status transition: %s -> %s", task.Durum, request.NewStatus),
				Field:  "durum",
			})
			continue
		}

		// Check dependencies if required
		if request.CheckDependencies && request.NewStatus == constants.TaskStatusInProgress {
			canStart, err := bp.checkDependenciesCompleted(taskID)
			if err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: taskID,
					Error:  fmt.Sprintf("dependency check failed: %v", err),
				})
				continue
			}
			if !canStart {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: taskID,
					Error:  "task has incomplete dependencies",
					Field:  "dependencies",
				})
				continue
			}
		}

		// Update status
		if err := bp.veriYonetici.GorevGuncelle(taskID, map[string]interface{}{"durum": request.NewStatus}); err != nil {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: taskID,
				Error:  fmt.Sprintf("update failed: %v", err),
			})
			continue
		}

		result.Successful = append(result.Successful, taskID)

		// Record interaction
		if bp.aiContextManager != nil {
			bp.aiContextManager.RecordInteraction(taskID, "bulk_status_change", map[string]interface{}{
				"old_status": task.Durum,
				"new_status": request.NewStatus,
			})
		}
	}

	result.TotalProcessed = len(request.TaskIDs)
	result.ExecutionTime = time.Since(startTime)
	result.Summary = fmt.Sprintf("Status transition to '%s': %d successful, %d failed",
		request.NewStatus, len(result.Successful), len(result.Failed))

	log.Printf("Bulk status transition completed: newStatus=%s, total=%d, successful=%d, failed=%d, duration=%v", request.NewStatus, result.TotalProcessed, len(result.Successful), len(result.Failed), result.ExecutionTime)

	return result, nil
}

// BulkTagOperation adds, removes, or replaces tags for multiple tasks
func (bp *BatchProcessor) BulkTagOperation(request BulkTagOperationRequest) (*BatchUpdateResult, error) {
	startTime := time.Now()

	result := &BatchUpdateResult{
		Successful: []string{},
		Failed:     []BatchUpdateError{},
		Warnings:   []BatchUpdateWarning{},
	}

	log.Printf("Starting bulk tag operation: count=%d, operation=%s, tags=%v", len(request.TaskIDs), request.Operation, request.Tags)

	// Validate operation
	if request.Operation != "add" && request.Operation != "remove" && request.Operation != "replace" {
		return nil, fmt.Errorf("invalid operation: %s. Must be 'add', 'remove', or 'replace'", request.Operation)
	}

	// Get or create tags for add/replace operations
	var tags []*Etiket
	if request.Operation == "add" || request.Operation == "replace" {
		var err error
		tags, err = bp.veriYonetici.EtiketleriGetirVeyaOlustur(request.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to get/create tags: %v", err)
		}
	}

	for _, taskID := range request.TaskIDs {
		if request.DryRun {
			if _, err := bp.veriYonetici.GorevDetay(taskID); err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: taskID,
					Error:  i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}),
				})
			} else {
				result.Warnings = append(result.Warnings, BatchUpdateWarning{
					TaskID:  taskID,
					Message: fmt.Sprintf("dry run - would %s tags: %v", request.Operation, request.Tags),
				})
			}
			continue
		}

		// Get current task
		task, err := bp.veriYonetici.GorevDetay(taskID)
		if err != nil {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: taskID,
				Error:  i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}),
			})
			continue
		}

		var newTags []*Etiket
		var updated bool

		switch request.Operation {
		case "add":
			// Add new tags to existing ones
			existingTagMap := make(map[string]*Etiket)
			for _, tag := range task.Etiketler {
				existingTagMap[tag.Isim] = tag
			}

			newTags = task.Etiketler
			for _, tag := range tags {
				if _, exists := existingTagMap[tag.Isim]; !exists {
					newTags = append(newTags, tag)
					updated = true
				}
			}

		case "remove":
			// Remove specified tags
			removeMap := make(map[string]bool)
			for _, tagName := range request.Tags {
				removeMap[tagName] = true
			}

			for _, tag := range task.Etiketler {
				if !removeMap[tag.Isim] {
					newTags = append(newTags, tag)
				} else {
					updated = true
				}
			}

		case "replace":
			// Replace all tags
			newTags = tags
			updated = len(task.Etiketler) != len(tags)
			if !updated {
				// Check if tags are actually different
				existingNames := make(map[string]bool)
				for _, tag := range task.Etiketler {
					existingNames[tag.Isim] = true
				}
				for _, tag := range tags {
					if !existingNames[tag.Isim] {
						updated = true
						break
					}
				}
			}
		}

		if updated {
			if err := bp.veriYonetici.GorevEtiketleriniAyarla(taskID, newTags); err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: taskID,
					Error:  fmt.Sprintf("tag operation failed: %v", err),
				})
				continue
			}

			result.Successful = append(result.Successful, taskID)

			// Record interaction
			if bp.aiContextManager != nil {
				bp.aiContextManager.RecordInteraction(taskID, "bulk_tag_operation", map[string]interface{}{
					"operation": request.Operation,
					"tags":      request.Tags,
				})
			}
		} else {
			result.Warnings = append(result.Warnings, BatchUpdateWarning{
				TaskID:  taskID,
				Message: "no changes needed",
			})
		}
	}

	result.TotalProcessed = len(request.TaskIDs)
	result.ExecutionTime = time.Since(startTime)
	result.Summary = fmt.Sprintf("Tag %s operation: %d successful, %d failed",
		request.Operation, len(result.Successful), len(result.Failed))

	log.Printf("Bulk tag operation completed: operation=%s, total=%d, successful=%d, failed=%d, duration=%v", request.Operation, result.TotalProcessed, len(result.Successful), len(result.Failed), result.ExecutionTime)

	return result, nil
}

// BulkDelete deletes multiple tasks with safety checks
func (bp *BatchProcessor) BulkDelete(request BulkDeleteRequest) (*BatchUpdateResult, error) {
	startTime := time.Now()

	result := &BatchUpdateResult{
		Successful: []string{},
		Failed:     []BatchUpdateError{},
		Warnings:   []BatchUpdateWarning{},
	}

	log.Printf("Starting bulk delete operation: count=%d", len(request.TaskIDs))

	// Safety check: require confirmation
	expectedConfirmation := fmt.Sprintf("DELETE %d TASKS", len(request.TaskIDs))
	if !request.Force && request.Confirmation != expectedConfirmation {
		return nil, fmt.Errorf("confirmation required: '%s'", expectedConfirmation)
	}

	for _, taskID := range request.TaskIDs {
		if request.DryRun {
			if _, err := bp.veriYonetici.GorevDetay(taskID); err != nil {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: taskID,
					Error:  i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}),
				})
			} else {
				result.Warnings = append(result.Warnings, BatchUpdateWarning{
					TaskID:  taskID,
					Message: "dry run - would be deleted",
				})
			}
			continue
		}

		// Check if task exists
		task, err := bp.veriYonetici.GorevDetay(taskID)
		if err != nil {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: taskID,
				Error:  i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}),
			})
			continue
		}

		// Check for subtasks
		if !request.DeleteSubtasks {
			subtasks, err := bp.veriYonetici.AltGorevleriGetir(taskID)
			if err == nil && len(subtasks) > 0 {
				result.Failed = append(result.Failed, BatchUpdateError{
					TaskID: taskID,
					Error:  fmt.Sprintf("task has %d subtasks, use delete_subtasks=true to force", len(subtasks)),
				})
				continue
			}
		}

		// Check for dependencies (tasks that depend on this one)
		// This would require a reverse dependency lookup - for now, we'll skip this check

		// Delete the task
		if err := bp.veriYonetici.GorevSil(taskID); err != nil {
			result.Failed = append(result.Failed, BatchUpdateError{
				TaskID: taskID,
				Error:  fmt.Sprintf("deletion failed: %v", err),
			})
			continue
		}

		result.Successful = append(result.Successful, taskID)

		// Record interaction
		if bp.aiContextManager != nil {
			bp.aiContextManager.RecordInteraction(taskID, "bulk_delete", map[string]interface{}{
				"task_title": task.Baslik,
				"force":      request.Force,
			})
		}
	}

	result.TotalProcessed = len(request.TaskIDs)
	result.ExecutionTime = time.Since(startTime)
	result.Summary = fmt.Sprintf("Bulk delete: %d successful, %d failed",
		len(result.Successful), len(result.Failed))

	log.Printf("Bulk delete completed: total=%d, successful=%d, failed=%d, duration=%v", result.TotalProcessed, len(result.Successful), len(result.Failed), result.ExecutionTime)

	return result, nil
}

// Helper validation methods

func (bp *BatchProcessor) validateUpdateRequest(request BatchUpdateRequest) error {
	// Check if task exists
	if _, err := bp.veriYonetici.GorevDetay(request.TaskID); err != nil {
		return fmt.Errorf(i18n.T("error.taskNotFound", map[string]interface{}{"Error": err}))
	}

	// Validate individual fields
	if status, ok := request.Updates["durum"].(string); ok {
		if !bp.validateStatus(status) {
			return fmt.Errorf("invalid status: %s", status)
		}
	}

	if priority, ok := request.Updates["oncelik"].(string); ok {
		if !bp.validatePriority(priority) {
			return fmt.Errorf("invalid priority: %s", priority)
		}
	}

	if title, ok := request.Updates["baslik"].(string); ok {
		if strings.TrimSpace(title) == "" {
			return fmt.Errorf("title cannot be empty")
		}
	}

	return nil
}

func (bp *BatchProcessor) validateStatus(status string) bool {
	validStatuses := constants.GetValidTaskStatuses()
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func (bp *BatchProcessor) validatePriority(priority string) bool {
	validPriorities := constants.GetValidPriorities()
	for _, valid := range validPriorities {
		if priority == valid {
			return true
		}
	}
	return false
}

func (bp *BatchProcessor) validateStatusTransition(from, to string) bool {
	// Define valid transitions
	validTransitions := map[string][]string{
		constants.TaskStatusPending:    {constants.TaskStatusInProgress, constants.TaskStatusCancelled},
		constants.TaskStatusInProgress: {constants.TaskStatusPending, constants.TaskStatusCompleted, constants.TaskStatusCancelled},
		constants.TaskStatusCompleted:  {constants.TaskStatusInProgress}, // Allow reopening
		constants.TaskStatusCancelled:  {constants.TaskStatusPending},    // Allow reactivation
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowedTo := range allowed {
		if to == allowedTo {
			return true
		}
	}
	return false
}

func (bp *BatchProcessor) checkDependenciesCompleted(taskID string) (bool, error) {
	dependencies, err := bp.veriYonetici.GorevBagimlilikGetir(taskID)
	if err != nil {
		return false, err
	}

	for _, dep := range dependencies {
		if dep.Durum != constants.TaskStatusCompleted {
			return false, nil
		}
	}

	return true, nil
}
