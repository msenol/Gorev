package gorev

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
)

// ExportFormat defines the structure of exported data
type ExportFormat struct {
	Version      string               `json:"version"`
	Metadata     ExportMetadata       `json:"metadata"`
	Projects     []*Proje             `json:"projects"`
	Tasks        []*Gorev             `json:"tasks"`
	Tags         []*Etiket            `json:"tags"`
	TaskTags     []TaskTagAssociation `json:"task_tags"`
	Templates    []*GorevTemplate     `json:"templates"`
	Dependencies []*Baglanti          `json:"dependencies"`
	AIContext    []*AIInteraction     `json:"ai_context,omitempty"`
}

// ExportMetadata contains metadata about the export
type ExportMetadata struct {
	ExportDate      time.Time `json:"export_date"`
	GorevVersion    string    `json:"gorev_version"`
	DatabaseVersion string    `json:"database_version"`
	TotalTasks      int       `json:"total_tasks"`
	TotalProjects   int       `json:"total_projects"`
	ExportedBy      string    `json:"exported_by,omitempty"`
	Description     string    `json:"description,omitempty"`
}

// TaskTagAssociation represents the many-to-many relationship between tasks and tags
type TaskTagAssociation struct {
	TaskID string `json:"task_id"`
	TagID  string `json:"tag_id"`
}

// ExportOptions contains options for data export
type ExportOptions struct {
	Format              string     `json:"format"` // json, csv
	OutputPath          string     `json:"output_path"`
	DateRange           *DateRange `json:"date_range,omitempty"`
	ProjectFilter       []string   `json:"project_filter,omitempty"`
	IncludeCompleted    bool       `json:"include_completed"`
	IncludeDependencies bool       `json:"include_dependencies"`
	IncludeMetadata     bool       `json:"include_metadata"`
	IncludeAIContext    bool       `json:"include_ai_context"`
	IncludeTemplates    bool       `json:"include_templates"`
}

// ImportOptions contains options for data import
type ImportOptions struct {
	FilePath           string            `json:"file_path"`
	ImportMode         string            `json:"import_mode"`         // merge, replace
	ConflictResolution string            `json:"conflict_resolution"` // skip, overwrite, prompt
	PreserveIDs        bool              `json:"preserve_ids"`
	ProjectMapping     map[string]string `json:"project_mapping,omitempty"`
	DryRun             bool              `json:"dry_run"`
}

// DateRange represents a date range for filtering
type DateRange struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

// ConflictResolution represents conflicts during import
type ConflictResolution struct {
	Type       string      `json:"type"` // task, project, tag, template
	Existing   interface{} `json:"existing"`
	Incoming   interface{} `json:"incoming"`
	Resolution string      `json:"resolution"` // skip, overwrite, rename
	NewValue   interface{} `json:"new_value,omitempty"`
}

// ImportResult contains the result of an import operation
type ImportResult struct {
	Success           bool                 `json:"success"`
	ImportedTasks     int                  `json:"imported_tasks"`
	ImportedProjects  int                  `json:"imported_projects"`
	ImportedTags      int                  `json:"imported_tags"`
	ImportedTemplates int                  `json:"imported_templates"`
	Conflicts         []ConflictResolution `json:"conflicts,omitempty"`
	Errors            []string             `json:"errors,omitempty"`
	Warnings          []string             `json:"warnings,omitempty"`
}

// ExportData exports data from the database according to the given options
func (iy *IsYonetici) ExportData(ctx context.Context, options ExportOptions) (*ExportFormat, error) {
	if iy.veriYonetici == nil {
		return nil, fmt.Errorf(i18n.T("error.dataManagerNotInitialized", nil))
	}

	// Validate export data options (no file path required)
	if err := iy.validateExportDataOptions(options); err != nil {
		return nil, fmt.Errorf(i18n.T("error.invalidExportOptions", map[string]interface{}{"Error": err}))
	}

	exportData := &ExportFormat{
		Version: "v0.11.1",
		Metadata: ExportMetadata{
			ExportDate:      time.Now(),
			GorevVersion:    "v0.11.1",
			DatabaseVersion: "1.9", // Current migration version
			ExportedBy:      "gorev_export",
			Description:     i18n.T("export.generatedDescription", nil),
		},
	}

	// Export projects using IsYonetici method
	projects, err := iy.ProjeListele(ctx)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.failedToExportProjects", map[string]interface{}{"Error": err}))
	}

	// Filter projects if specified
	if len(options.ProjectFilter) > 0 {
		filteredProjects := []*Proje{}
		projectFilterMap := make(map[string]bool)
		for _, pid := range options.ProjectFilter {
			projectFilterMap[pid] = true
		}
		for _, project := range projects {
			if projectFilterMap[project.ID] {
				filteredProjects = append(filteredProjects, project)
			}
		}
		projects = filteredProjects
	}

	exportData.Projects = projects
	exportData.Metadata.TotalProjects = len(projects)

	// Create project ID map for task filtering
	projectIDMap := make(map[string]bool)
	for _, project := range projects {
		projectIDMap[project.ID] = true
	}

	// Export tasks with filters
	filters := make(map[string]interface{})

	// Apply project filter
	if len(options.ProjectFilter) > 0 {
		// Will be filtered below since GorevListele doesn't support multiple project filter
	}

	// Apply date range filter
	if options.DateRange != nil {
		// Note: Current GorevListele doesn't support date range, will filter after retrieval
	}

	// Get all tasks first
	allTasks, err := iy.veriYonetici.GorevListele(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.failedToExportTasks", map[string]interface{}{"Error": err}))
	}

	// Filter tasks based on options
	filteredTasks := []*Gorev{}
	for _, task := range allTasks {
		// Filter by project if specified
		if len(options.ProjectFilter) > 0 {
			if task.ProjeID == "" || !projectIDMap[task.ProjeID] {
				continue
			}
		}

		// Filter by completed status
		if !options.IncludeCompleted && task.Status == constants.TaskStatusCompleted {
			continue
		}

		// Filter by date range
		if options.DateRange != nil {
			if options.DateRange.From != nil && task.CreatedAt.Before(*options.DateRange.From) {
				continue
			}
			if options.DateRange.To != nil && task.CreatedAt.After(*options.DateRange.To) {
				continue
			}
		}

		filteredTasks = append(filteredTasks, task)
	}

	exportData.Tasks = filteredTasks
	exportData.Metadata.TotalTasks = len(filteredTasks)

	// Export tags - collect unique tags from all filtered tasks
	tagMap := make(map[string]*Etiket)
	for _, task := range filteredTasks {
		if task.Tags != nil {
			for _, tag := range task.Tags {
				tagMap[tag.ID] = tag
			}
		}
	}

	// Convert map to slice
	tags := []*Etiket{}
	for _, tag := range tagMap {
		tags = append(tags, tag)
	}
	exportData.Tags = tags

	// Export task-tag associations
	taskTags, err := iy.exportTaskTagAssociations(ctx, filteredTasks)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.failedToExportTaskTags", map[string]interface{}{"Error": err}))
	}
	exportData.TaskTags = taskTags

	// Export dependencies if requested
	if options.IncludeDependencies {
		dependencies, err := iy.exportTaskDependencies(ctx, filteredTasks)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.failedToExportDependencies", map[string]interface{}{"Error": err}))
		}
		exportData.Dependencies = dependencies
	}

	// Export templates if requested
	if options.IncludeTemplates {
		templates, err := iy.veriYonetici.TemplateListele(ctx, "")
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.failedToExportTemplates", map[string]interface{}{"Error": err}))
		}
		exportData.Templates = templates
	}

	// Export AI context if requested
	if options.IncludeAIContext {
		aiContext, err := iy.exportAIContext(ctx, filteredTasks)
		if err != nil {
			// Log warning but continue with export - AI context is optional
			fmt.Printf("Warning: AI context export failed: %v\n", err)
			exportData.AIContext = []*AIInteraction{}
		} else {
			exportData.AIContext = aiContext
		}
	}

	return exportData, nil
}

// SaveExportToFile saves export data to a file
func (iy *IsYonetici) SaveExportToFile(ctx context.Context, exportData *ExportFormat, options ExportOptions) error {
	// Validate export file options
	if err := iy.validateExportFileOptions(options); err != nil {
		return fmt.Errorf(i18n.T("error.invalidExportOptions", map[string]interface{}{"Error": err}))
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(options.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf(i18n.T("error.failedToCreateDirectory", map[string]interface{}{"Path": outputDir, "Error": err}))
	}

	switch options.Format {
	case "json", "":
		return iy.saveAsJSON(exportData, options.OutputPath)
	case "csv":
		return iy.saveAsCSV(exportData, options.OutputPath)
	default:
		return fmt.Errorf(i18n.T("error.unsupportedExportFormat", map[string]interface{}{"Format": options.Format}))
	}
}

// saveAsJSON saves export data as JSON
func (iy *IsYonetici) saveAsJSON(exportData *ExportFormat, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf(i18n.T("error.failedToCreateFile", map[string]interface{}{"Path": outputPath, "Error": err}))
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("Warning: failed to close file %s: %v\n", outputPath, cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf(i18n.T("error.failedToEncodeJSON", map[string]interface{}{"Error": err}))
	}

	return nil
}

// saveAsCSV saves export data as CSV (simplified format)
func (iy *IsYonetici) saveAsCSV(exportData *ExportFormat, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf(i18n.T("error.failedToCreateFile", map[string]interface{}{"Path": outputPath, "Error": err}))
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("Warning: failed to close file %s: %v\n", outputPath, cerr)
		}
	}()

	// Write CSV header
	header := "ID,Title,Description,Status,Priority,Project,Created,Updated,Tags\n"
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf(i18n.T("error.failedToWriteFile", map[string]interface{}{"Error": err}))
	}

	// Create project name map
	projectNames := make(map[string]string)
	for _, project := range exportData.Projects {
		projectNames[project.ID] = project.Name
	}

	// Create task tag map
	taskTagMap := make(map[string][]string)
	for _, taskTag := range exportData.TaskTags {
		taskTagMap[taskTag.TaskID] = append(taskTagMap[taskTag.TaskID], taskTag.TagID)
	}

	// Create tag name map
	tagNames := make(map[string]string)
	for _, tag := range exportData.Tags {
		tagNames[tag.ID] = tag.Name
	}

	// Write task data
	for _, task := range exportData.Tasks {
		projectName := ""
		if task.ProjeID != "" {
			projectName = projectNames[task.ProjeID]
		}

		// Build tags string
		taskTagNames := []string{}
		for _, tagID := range taskTagMap[task.ID] {
			if tagName, exists := tagNames[tagID]; exists {
				taskTagNames = append(taskTagNames, tagName)
			}
		}
		tagsStr := ""
		if len(taskTagNames) > 0 {
			tagsStr = "\"" + fmt.Sprintf("%v", taskTagNames) + "\""
		}

		// Escape CSV values
		title := escapeCSVValue(task.Title)
		description := escapeCSVValue(task.Description)

		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			task.ID,
			title,
			description,
			task.Status,
			task.Priority,
			escapeCSVValue(projectName),
			task.CreatedAt.Format(time.RFC3339),
			task.UpdatedAt.Format(time.RFC3339),
			tagsStr,
		)

		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf(i18n.T("error.failedToWriteFile", map[string]interface{}{"Error": err}))
		}
	}

	return nil
}

// escapeCSVValue escapes a string value for CSV
func escapeCSVValue(value string) string {
	if value == "" {
		return ""
	}
	// If value contains comma, quote, or newline, wrap in quotes and escape internal quotes
	needsQuoting := false
	for _, char := range value {
		if char == ',' || char == '"' || char == '\n' || char == '\r' {
			needsQuoting = true
			break
		}
	}

	if needsQuoting {
		// Escape internal quotes by doubling them
		escaped := ""
		for _, char := range value {
			if char == '"' {
				escaped += "\"\""
			} else {
				escaped += string(char)
			}
		}
		return "\"" + escaped + "\""
	}

	return value
}

// validateExportDataOptions validates export data options (no file path required)
func (iy *IsYonetici) validateExportDataOptions(options ExportOptions) error {
	if options.DateRange != nil {
		if options.DateRange.From != nil && options.DateRange.To != nil {
			if options.DateRange.From.After(*options.DateRange.To) {
				return fmt.Errorf(i18n.T("error.invalidDateRange", nil))
			}
		}
	}

	return nil
}

// validateExportFileOptions validates export file options (requires output path)
func (iy *IsYonetici) validateExportFileOptions(options ExportOptions) error {
	if options.OutputPath == "" {
		return fmt.Errorf(i18n.T("error.outputPathRequired", nil))
	}

	if options.Format != "" && options.Format != "json" && options.Format != "csv" {
		return fmt.Errorf(i18n.T("error.invalidFormat", map[string]interface{}{"Format": options.Format}))
	}

	return iy.validateExportDataOptions(options)
}

// exportTaskTagAssociations exports task-tag associations
func (iy *IsYonetici) exportTaskTagAssociations(ctx context.Context, tasks []*Gorev) ([]TaskTagAssociation, error) {
	taskTags := []TaskTagAssociation{}

	for _, task := range tasks {
		if task.Tags != nil {
			for _, tag := range task.Tags {
				taskTags = append(taskTags, TaskTagAssociation{
					TaskID: task.ID,
					TagID:  tag.ID,
				})
			}
		}
	}

	return taskTags, nil
}

// exportTaskDependencies exports task dependencies
func (iy *IsYonetici) exportTaskDependencies(ctx context.Context, tasks []*Gorev) ([]*Baglanti, error) {
	// Create a map of task IDs for filtering
	taskIDMap := make(map[string]bool)
	for _, task := range tasks {
		taskIDMap[task.ID] = true
	}

	// Get dependencies for each task and combine them
	dependencies := []*Baglanti{}
	for _, task := range tasks {
		taskDeps, err := iy.veriYonetici.BaglantilariGetir(ctx, task.ID)
		if err != nil {
			continue // Skip on error
		}
		dependencies = append(dependencies, taskDeps...)
	}

	// Filter dependencies to only include those between exported tasks
	filteredDependencies := []*Baglanti{}
	for _, dep := range dependencies {
		if taskIDMap[dep.SourceID] && taskIDMap[dep.TargetID] {
			filteredDependencies = append(filteredDependencies, dep)
		}
	}

	return filteredDependencies, nil
}

// exportAIContext exports AI context data
func (iy *IsYonetici) exportAIContext(ctx context.Context, tasks []*Gorev) ([]*AIInteraction, error) {
	// Create a map of task IDs for filtering
	taskIDMap := make(map[string]bool)
	for _, task := range tasks {
		taskIDMap[task.ID] = true
	}

	// This would require implementing GetAIInteractions in AIContextYonetici
	// For now, return empty slice
	return []*AIInteraction{}, nil
}

// ImportData imports data from a file according to the given options
func (iy *IsYonetici) ImportData(ctx context.Context, options ImportOptions) (*ImportResult, error) {
	if iy.veriYonetici == nil {
		return nil, fmt.Errorf(i18n.T("error.dataManagerNotInitialized", nil))
	}

	// Validate import options
	if err := iy.validateImportOptions(options); err != nil {
		return nil, fmt.Errorf(i18n.T("error.invalidImportOptions", map[string]interface{}{"Error": err}))
	}

	// Load import data
	importData, err := iy.loadImportData(options.FilePath)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.failedToLoadImportData", map[string]interface{}{"Error": err}))
	}

	// Validate import data structure
	if err := iy.validateImportData(importData); err != nil {
		return nil, fmt.Errorf(i18n.T("error.invalidImportData", map[string]interface{}{"Error": err}))
	}

	result := &ImportResult{
		Success:   true,
		Conflicts: []ConflictResolution{},
		Errors:    []string{},
		Warnings:  []string{},
	}

	// If dry run, analyze conflicts without making changes
	if options.DryRun {
		conflicts, err := iy.analyzeImportConflicts(ctx, importData, options)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("error.failedToAnalyzeConflicts", map[string]interface{}{"Error": err}))
		}
		result.Conflicts = conflicts
		return result, nil
	}

	// Note: Currently no transaction support in veriYonetici, will implement without transaction
	// This is acceptable for import as most operations are already atomic at DB level

	// Import in order: Projects -> Tags -> Templates -> Tasks -> Dependencies -> AI Context

	// Import projects
	if len(importData.Projects) > 0 {
		imported, conflicts, err := iy.importProjects(ctx, importData.Projects, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Projects: %v", err))
			result.Success = false
		} else {
			result.ImportedProjects = imported
			result.Conflicts = append(result.Conflicts, conflicts...)
		}
	}

	// Import tags
	if len(importData.Tags) > 0 {
		imported, conflicts, err := iy.importTags(ctx, importData.Tags, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Tags: %v", err))
			result.Success = false
		} else {
			result.ImportedTags = imported
			result.Conflicts = append(result.Conflicts, conflicts...)
		}
	}

	// Import templates
	if len(importData.Templates) > 0 {
		imported, conflicts, err := iy.importTemplates(ctx, importData.Templates, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Templates: %v", err))
			result.Success = false
		} else {
			result.ImportedTemplates = imported
			result.Conflicts = append(result.Conflicts, conflicts...)
		}
	}

	// Import tasks
	if len(importData.Tasks) > 0 {
		imported, conflicts, err := iy.importTasks(ctx, importData.Tasks, options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Tasks: %v", err))
			result.Success = false
		} else {
			result.ImportedTasks = imported
			result.Conflicts = append(result.Conflicts, conflicts...)
		}
	}

	// Import task-tag associations
	if len(importData.TaskTags) > 0 {
		if err := iy.importTaskTagAssociations(ctx, importData.TaskTags, options); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Task-Tag associations: %v", err))
		}
	}

	// Import dependencies
	if len(importData.Dependencies) > 0 {
		if err := iy.importDependencies(ctx, importData.Dependencies, options); err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Dependencies: %v", err))
		}
	}

	return result, nil
}

// loadImportData loads and parses import data from file
func (iy *IsYonetici) loadImportData(filePath string) (*ExportFormat, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.failedToOpenFile", map[string]interface{}{"Path": filePath, "Error": err}))
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Printf("Warning: failed to close file %s: %v\n", filePath, cerr)
		}
	}()

	var importData ExportFormat
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&importData); err != nil {
		return nil, fmt.Errorf(i18n.T("error.failedToDecodeJSON", map[string]interface{}{"Error": err}))
	}

	return &importData, nil
}

// validateImportOptions validates import options
func (iy *IsYonetici) validateImportOptions(options ImportOptions) error {
	if options.FilePath == "" {
		return fmt.Errorf(i18n.T("error.filePathRequired", nil))
	}

	if _, err := os.Stat(options.FilePath); os.IsNotExist(err) {
		return fmt.Errorf(i18n.T("error.fileNotFound", map[string]interface{}{"Path": options.FilePath}))
	}

	if options.ImportMode != "" && options.ImportMode != "merge" && options.ImportMode != "replace" {
		return fmt.Errorf(i18n.T("error.invalidImportMode", map[string]interface{}{"Mode": options.ImportMode}))
	}

	if options.ConflictResolution != "" &&
		options.ConflictResolution != "skip" &&
		options.ConflictResolution != "overwrite" &&
		options.ConflictResolution != "prompt" {
		return fmt.Errorf(i18n.T("error.invalidConflictResolution", map[string]interface{}{"Resolution": options.ConflictResolution}))
	}

	return nil
}

// validateImportData validates the structure of import data
func (iy *IsYonetici) validateImportData(importData *ExportFormat) error {
	if importData.Version == "" {
		return fmt.Errorf(i18n.T("error.missingVersion", nil))
	}

	// Validate that all task project references exist in projects list
	projectIDMap := make(map[string]bool)
	for _, project := range importData.Projects {
		projectIDMap[project.ID] = true
	}

	for _, task := range importData.Tasks {
		if task.ProjeID != "" && !projectIDMap[task.ProjeID] {
			return fmt.Errorf(i18n.T("error.invalidTaskProjectReference",
				map[string]interface{}{"TaskID": task.ID, "ProjectID": task.ProjeID}))
		}
	}

	// Validate that all task-tag references exist
	taskIDMap := make(map[string]bool)
	for _, task := range importData.Tasks {
		taskIDMap[task.ID] = true
	}

	tagIDMap := make(map[string]bool)
	for _, tag := range importData.Tags {
		tagIDMap[tag.ID] = true
	}

	for _, taskTag := range importData.TaskTags {
		if !taskIDMap[taskTag.TaskID] {
			return fmt.Errorf(i18n.T("error.invalidTaskTagTaskReference",
				map[string]interface{}{"TaskID": taskTag.TaskID}))
		}
		if !tagIDMap[taskTag.TagID] {
			return fmt.Errorf(i18n.T("error.invalidTaskTagTagReference",
				map[string]interface{}{"TagID": taskTag.TagID}))
		}
	}

	// Validate dependencies
	for _, dep := range importData.Dependencies {
		if !taskIDMap[dep.SourceID] {
			return fmt.Errorf(i18n.T("error.invalidDependencySourceReference",
				map[string]interface{}{"SourceID": dep.SourceID}))
		}
		if !taskIDMap[dep.TargetID] {
			return fmt.Errorf(i18n.T("error.invalidDependencyTargetReference",
				map[string]interface{}{"TargetID": dep.TargetID}))
		}
	}

	return nil
}

// analyzeImportConflicts analyzes potential conflicts without importing
func (iy *IsYonetici) analyzeImportConflicts(ctx context.Context, importData *ExportFormat, options ImportOptions) ([]ConflictResolution, error) {
	conflicts := []ConflictResolution{}

	// Check project conflicts
	for _, project := range importData.Projects {
		existing, err := iy.veriYonetici.ProjeGetir(ctx, project.ID)
		if err == nil && existing != nil {
			conflicts = append(conflicts, ConflictResolution{
				Type:       "project",
				Existing:   existing,
				Incoming:   project,
				Resolution: options.ConflictResolution,
			})
		}
	}

	// Check task conflicts
	for _, task := range importData.Tasks {
		existing, err := iy.veriYonetici.GorevGetir(ctx, task.ID)
		if err == nil && existing != nil {
			conflicts = append(conflicts, ConflictResolution{
				Type:       "task",
				Existing:   existing,
				Incoming:   task,
				Resolution: options.ConflictResolution,
			})
		}
	}

	// Check tag conflicts by name (using existing tags from the system)
	// Note: Since there's no EtiketGetirIsimile method, we'll need to get all tags and compare
	// This will be implemented when we create missing methods

	return conflicts, nil
}

// importProjects imports project data
func (iy *IsYonetici) importProjects(ctx context.Context, projects []*Proje, options ImportOptions) (int, []ConflictResolution, error) {
	imported := 0
	conflicts := []ConflictResolution{}

	for _, project := range projects {
		projectID := project.ID

		// Generate new ID if not preserving IDs
		if !options.PreserveIDs {
			projectID = uuid.New().String()
		}

		// Check for conflicts
		existing, err := iy.veriYonetici.ProjeGetir(ctx, projectID)
		if err == nil && existing != nil {
			// Handle conflict
			switch options.ConflictResolution {
			case "skip", "":
				conflicts = append(conflicts, ConflictResolution{
					Type:       "project",
					Existing:   existing,
					Incoming:   project,
					Resolution: "skip",
				})
				continue
			case "overwrite":
				// Update existing project using ProjeKaydet (since ProjeGuncelle doesn't exist)
				project.ID = projectID
				if err := iy.veriYonetici.ProjeKaydet(ctx, project); err != nil {
					return imported, conflicts, err
				}
				conflicts = append(conflicts, ConflictResolution{
					Type:       "project",
					Existing:   existing,
					Incoming:   project,
					Resolution: "overwrite",
					NewValue:   project,
				})
			case "prompt":
				// Return conflict for user resolution
				conflicts = append(conflicts, ConflictResolution{
					Type:       "project",
					Existing:   existing,
					Incoming:   project,
					Resolution: "prompt",
				})
				continue
			}
		} else {
			// No conflict, create new project
			project.ID = projectID
			if err := iy.veriYonetici.ProjeKaydet(ctx, project); err != nil {
				return imported, conflicts, err
			}
		}

		imported++
	}

	return imported, conflicts, nil
}

// importTags imports tag data
func (iy *IsYonetici) importTags(ctx context.Context, tags []*Etiket, options ImportOptions) (int, []ConflictResolution, error) {
	imported := 0
	conflicts := []ConflictResolution{}

	for _, tag := range tags {
		// Create or get existing tags using EtiketleriGetirVeyaOlustur
		// This method already handles duplicates by name
		tagNames := []string{tag.Name}
		createdTags, err := iy.veriYonetici.EtiketleriGetirVeyaOlustur(ctx, tagNames)
		if err != nil {
			return imported, conflicts, err
		}

		// If tag was created or found, count as imported
		if len(createdTags) > 0 {
			imported++
		}
	}

	return imported, conflicts, nil
}

// importTemplates imports template data
func (iy *IsYonetici) importTemplates(ctx context.Context, templates []*GorevTemplate, options ImportOptions) (int, []ConflictResolution, error) {
	imported := 0
	conflicts := []ConflictResolution{}

	for _, template := range templates {
		templateID := template.ID

		// Generate new ID if not preserving IDs
		if !options.PreserveIDs {
			templateID = uuid.New().String()
		}

		// Check for conflicts by ID (no TemplateGetirIsimile available)
		existing, err := iy.veriYonetici.TemplateGetir(ctx, templateID)

		if err == nil && existing != nil {
			// Handle conflict
			switch options.ConflictResolution {
			case "skip", "":
				conflicts = append(conflicts, ConflictResolution{
					Type:       "template",
					Existing:   existing,
					Incoming:   template,
					Resolution: "skip",
				})
				continue
			case "overwrite":
				template.ID = existing.ID // Use existing ID
				if err := iy.veriYonetici.TemplateOlustur(ctx, template); err != nil {
					return imported, conflicts, err
				}
				conflicts = append(conflicts, ConflictResolution{
					Type:       "template",
					Existing:   existing,
					Incoming:   template,
					Resolution: "overwrite",
					NewValue:   template,
				})
			case "prompt":
				conflicts = append(conflicts, ConflictResolution{
					Type:       "template",
					Existing:   existing,
					Incoming:   template,
					Resolution: "prompt",
				})
				continue
			}
		} else {
			// No conflict, create new template
			template.ID = templateID
			if err := iy.veriYonetici.TemplateOlustur(ctx, template); err != nil {
				return imported, conflicts, err
			}
		}

		imported++
	}

	return imported, conflicts, nil
}

// importTasks imports task data
func (iy *IsYonetici) importTasks(ctx context.Context, tasks []*Gorev, options ImportOptions) (int, []ConflictResolution, error) {
	imported := 0
	conflicts := []ConflictResolution{}

	for _, task := range tasks {
		taskID := task.ID

		// Generate new ID if not preserving IDs
		if !options.PreserveIDs {
			taskID = uuid.New().String()
		}

		// Map project ID if project mapping is provided
		if task.ProjeID != "" && options.ProjectMapping != nil {
			if newProjectID, exists := options.ProjectMapping[task.ProjeID]; exists {
				task.ProjeID = newProjectID
			}
		}

		// Check for conflicts (both by ID and potentially by title for better duplicate detection)
		existing, err := iy.veriYonetici.GorevGetir(ctx, taskID)
		if err == nil && existing != nil {
			log.Printf("Import: Found existing task with ID %s, applying conflict resolution: %s", taskID, options.ConflictResolution)
			// Handle conflict
			switch options.ConflictResolution {
			case "skip", "":
				conflicts = append(conflicts, ConflictResolution{
					Type:       "task",
					Existing:   existing,
					Incoming:   task,
					Resolution: "skip",
				})
				continue
			case "overwrite":
				task.ID = taskID
				if err := iy.veriYonetici.GorevGuncelle(ctx, taskID, map[string]interface{}{
					"title":       task.Title,
					"description": task.Description,
					"status":      task.Status,
					"priority":    task.Priority,
					"project_id":  task.ProjeID,
					"due_date":    task.DueDate,
				}); err != nil {
					return imported, conflicts, err
				}
				conflicts = append(conflicts, ConflictResolution{
					Type:       "task",
					Existing:   existing,
					Incoming:   task,
					Resolution: "overwrite",
					NewValue:   task,
				})
			case "prompt":
				conflicts = append(conflicts, ConflictResolution{
					Type:       "task",
					Existing:   existing,
					Incoming:   task,
					Resolution: "prompt",
				})
				continue
			}
		} else {
			// No conflict, create new task using GorevOlustur with params map
			task.ID = taskID
			params := map[string]interface{}{
				"id":          taskID,
				"title":       task.Title,
				"description": task.Description,
				"status":      task.Status,
				"priority":    task.Priority,
				"project_id":  task.ProjeID,
				"parent_id":   task.ParentID,
				"due_date":    task.DueDate,
			}
			if _, err := iy.veriYonetici.GorevOlustur(ctx, params); err != nil {
				log.Printf("Import: Failed to create new task %s: %v", taskID, err)
				return imported, conflicts, err
			}
			log.Printf("Import: Successfully created new task %s", taskID)
		}

		imported++
	}

	return imported, conflicts, nil
}

// importTaskTagAssociations imports task-tag associations
func (iy *IsYonetici) importTaskTagAssociations(ctx context.Context, taskTags []TaskTagAssociation, options ImportOptions) error {
	for _, taskTag := range taskTags {
		// Verify task exists
		task, err := iy.veriYonetici.GorevGetir(ctx, taskTag.TaskID)
		if err != nil || task == nil {
			continue // Skip if task doesn't exist
		}

		// Get all tags for the task and add the association using GorevEtiketleriniAyarla
		// This is a simplified approach - we'd need the actual tag by ID
		// For now, skip tag associations as we don't have a direct method to add individual tags
		log.Printf("Note: Task-tag association skipped for task %s - requires tag lookup by ID", taskTag.TaskID)
	}

	return nil
}

// importDependencies imports task dependencies
func (iy *IsYonetici) importDependencies(ctx context.Context, dependencies []*Baglanti, options ImportOptions) error {
	for _, dep := range dependencies {
		// Verify both tasks exist
		source, err := iy.veriYonetici.GorevGetir(ctx, dep.SourceID)
		if err != nil || source == nil {
			continue // Skip if source task doesn't exist
		}

		target, err := iy.veriYonetici.GorevGetir(ctx, dep.TargetID)
		if err != nil || target == nil {
			continue // Skip if target task doesn't exist
		}

		// Create dependency using BaglantiEkle
		newDep := &Baglanti{
			ID:             uuid.New().String(),
			SourceID:       dep.SourceID,
			TargetID:       dep.TargetID,
			ConnectionType: dep.ConnectionType,
		}
		if err := iy.veriYonetici.BaglantiEkle(ctx, newDep); err != nil {
			// Log warning but continue
			log.Printf("Warning: Failed to create dependency from %s to %s: %v", dep.SourceID, dep.TargetID, err)
		}
	}

	return nil
}
