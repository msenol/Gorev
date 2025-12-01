// Package testing provides test utilities and data seeding for the Gorev project.
package testing

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/msenol/gorev/internal/testing/fixtures"
)

// SeederConfig configures the test data seeder
type SeederConfig struct {
	Language        string // "tr" or "en" (default: "tr")
	WorkspaceID     string // Workspace ID for centralized mode
	IncludeSubtasks bool   // Create subtask hierarchies (default: true)
	IncludeDeps     bool   // Create dependencies (default: true)
	Minimal         bool   // Use minimal data set (3 tasks instead of 15)
}

// DefaultSeederConfig returns the default seeder configuration
func DefaultSeederConfig() *SeederConfig {
	return &SeederConfig{
		Language:        "tr",
		WorkspaceID:     "",
		IncludeSubtasks: true,
		IncludeDeps:     true,
		Minimal:         false,
	}
}

// SeedResult contains the results of seeding operations
type SeedResult struct {
	Projects     []*gorev.Proje
	Tasks        []*gorev.Gorev
	Subtasks     []*gorev.Gorev
	Dependencies []DependencyResult
	Tags         []*gorev.Etiket
}

// DependencyResult represents a created dependency
type DependencyResult struct {
	SourceID string
	TargetID string
	Type     string
}

// TestDataSeeder handles seeding test data into the database
type TestDataSeeder struct {
	isYonetici *gorev.IsYonetici
	config     *SeederConfig
	ctx        context.Context
}

// NewTestDataSeeder creates a new test data seeder
func NewTestDataSeeder(isYonetici *gorev.IsYonetici, config *SeederConfig) *TestDataSeeder {
	if config == nil {
		config = DefaultSeederConfig()
	}

	// Set language in context
	ctx := i18n.WithLanguage(context.Background(), config.Language)

	return &TestDataSeeder{
		isYonetici: isYonetici,
		config:     config,
		ctx:        ctx,
	}
}

// SeedAll seeds all test data: projects, tasks, subtasks, dependencies
func (s *TestDataSeeder) SeedAll() (*SeedResult, error) {
	result := &SeedResult{}

	// 1. Seed projects
	projects, err := s.SeedProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to seed projects: %w", err)
	}
	result.Projects = projects

	// 2. Seed tasks
	tasks, err := s.SeedTasks(projects)
	if err != nil {
		return nil, fmt.Errorf("failed to seed tasks: %w", err)
	}
	result.Tasks = tasks

	// 3. Seed subtasks (if enabled)
	if s.config.IncludeSubtasks && !s.config.Minimal {
		subtasks, err := s.SeedSubtasks(tasks)
		if err != nil {
			return nil, fmt.Errorf("failed to seed subtasks: %w", err)
		}
		result.Subtasks = subtasks
	}

	// 4. Seed dependencies (if enabled)
	if s.config.IncludeDeps && !s.config.Minimal {
		deps, err := s.SeedDependencies(tasks)
		if err != nil {
			return nil, fmt.Errorf("failed to seed dependencies: %w", err)
		}
		result.Dependencies = deps
	}

	// 5. Set first project as active
	if len(projects) > 0 {
		if err := s.isYonetici.AktifProjeAyarla(s.ctx, projects[0].ID); err != nil {
			// Non-fatal error, just log it
			fmt.Printf("Warning: failed to set active project: %v\n", err)
		}
	}

	return result, nil
}

// SeedProjects creates sample projects
func (s *TestDataSeeder) SeedProjects() ([]*gorev.Proje, error) {
	var sampleProjects []fixtures.SampleProject
	if s.config.Minimal {
		sampleProjects = fixtures.MinimalSampleProjects
	} else {
		sampleProjects = fixtures.SampleProjects
	}

	projects := make([]*gorev.Proje, 0, len(sampleProjects))

	for _, sp := range sampleProjects {
		name := sp.NameTR
		definition := sp.DefinitionTR
		if s.config.Language == "en" {
			name = sp.NameEN
			definition = sp.DefinitionEN
		}

		proje, err := s.isYonetici.ProjeOlustur(s.ctx, name, definition)
		if err != nil {
			return nil, fmt.Errorf("failed to create project '%s': %w", name, err)
		}
		projects = append(projects, proje)
	}

	return projects, nil
}

// SeedTasks creates sample tasks for the given projects
func (s *TestDataSeeder) SeedTasks(projects []*gorev.Proje) ([]*gorev.Gorev, error) {
	var sampleTasks []fixtures.SampleTask
	if s.config.Minimal {
		sampleTasks = fixtures.MinimalSampleTasks
	} else {
		sampleTasks = fixtures.SampleTasks
	}

	tasks := make([]*gorev.Gorev, 0, len(sampleTasks))

	for _, st := range sampleTasks {
		// Validate project index
		if st.ProjectIndex >= len(projects) {
			return nil, fmt.Errorf("invalid project index %d for task '%s'", st.ProjectIndex, st.Values["title"])
		}

		projeID := projects[st.ProjectIndex].ID

		// Get template
		template, err := s.isYonetici.VeriYonetici().TemplateAliasIleGetir(s.ctx, st.TemplateAlias)
		if err != nil {
			return nil, fmt.Errorf("template '%s' not found: %w", st.TemplateAlias, err)
		}

		// Calculate due date
		var sonTarih *time.Time
		if st.DueDaysOffset != 0 {
			t := time.Now().AddDate(0, 0, st.DueDaysOffset)
			sonTarih = &t
		}

		// Create task using template
		task, err := s.createTaskFromTemplate(template, st.Values, projeID, st.Priority, sonTarih, st.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to create task '%s': %w", st.Values["title"], err)
		}

		// Update status if not default
		if st.Status != constants.TaskStatusPending {
			if err := s.isYonetici.VeriYonetici().GorevGuncelle(s.ctx, task.ID, map[string]interface{}{
				"status":     st.Status,
				"updated_at": time.Now(),
			}); err != nil {
				return nil, fmt.Errorf("failed to update task status: %w", err)
			}
			task.Status = st.Status
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// SeedSubtasks creates subtask hierarchies for parent tasks
func (s *TestDataSeeder) SeedSubtasks(parentTasks []*gorev.Gorev) ([]*gorev.Gorev, error) {
	subtasks := make([]*gorev.Gorev, 0)

	for _, ss := range fixtures.SampleSubtasks {
		// Validate parent index
		if ss.ParentIndex >= len(parentTasks) {
			return nil, fmt.Errorf("invalid parent index %d", ss.ParentIndex)
		}

		parentTask := parentTasks[ss.ParentIndex]

		// Create subtask and its children recursively
		createdSubtasks, err := s.createSubtaskHierarchy(parentTask.ID, parentTask.ProjeID, ss)
		if err != nil {
			return nil, fmt.Errorf("failed to create subtask hierarchy: %w", err)
		}
		subtasks = append(subtasks, createdSubtasks...)
	}

	return subtasks, nil
}

// createSubtaskHierarchy recursively creates subtasks
func (s *TestDataSeeder) createSubtaskHierarchy(parentID, projeID string, ss fixtures.SampleSubtask) ([]*gorev.Gorev, error) {
	result := make([]*gorev.Gorev, 0)

	// Get template
	template, err := s.isYonetici.VeriYonetici().TemplateAliasIleGetir(s.ctx, ss.TemplateAlias)
	if err != nil {
		return nil, fmt.Errorf("template '%s' not found: %w", ss.TemplateAlias, err)
	}

	// Create subtask
	subtask, err := s.createSubtaskFromTemplate(template, ss.Values, projeID, parentID, ss.Priority)
	if err != nil {
		return nil, fmt.Errorf("failed to create subtask: %w", err)
	}

	// Update status if not default
	if ss.Status != constants.TaskStatusPending && ss.Status != "" {
		if err := s.isYonetici.VeriYonetici().GorevGuncelle(s.ctx, subtask.ID, map[string]interface{}{
			"status":     ss.Status,
			"updated_at": time.Now(),
		}); err != nil {
			return nil, fmt.Errorf("failed to update subtask status: %w", err)
		}
		subtask.Status = ss.Status
	}

	result = append(result, subtask)

	// Recursively create children
	for _, child := range ss.Children {
		childSubtasks, err := s.createSubtaskHierarchy(subtask.ID, projeID, child)
		if err != nil {
			return nil, err
		}
		result = append(result, childSubtasks...)
	}

	return result, nil
}

// SeedDependencies creates task dependencies
func (s *TestDataSeeder) SeedDependencies(tasks []*gorev.Gorev) ([]DependencyResult, error) {
	dependencies := make([]DependencyResult, 0)

	for _, sd := range fixtures.SampleDependencies {
		// Validate indices
		if sd.SourceIndex >= len(tasks) || sd.TargetIndex >= len(tasks) {
			return nil, fmt.Errorf("invalid dependency indices: source=%d, target=%d", sd.SourceIndex, sd.TargetIndex)
		}

		sourceTask := tasks[sd.SourceIndex]
		targetTask := tasks[sd.TargetIndex]

		// Create dependency using VeriYonetici
		baglanti := &gorev.Baglanti{
			ID:             uuid.New().String(),
			SourceID:       sourceTask.ID,
			TargetID:       targetTask.ID,
			ConnectionType: sd.Type,
		}
		err := s.isYonetici.VeriYonetici().BaglantiEkle(s.ctx, baglanti)
		if err != nil {
			return nil, fmt.Errorf("failed to create dependency %s -> %s: %w", sourceTask.ID, targetTask.ID, err)
		}

		dependencies = append(dependencies, DependencyResult{
			SourceID: sourceTask.ID,
			TargetID: targetTask.ID,
			Type:     sd.Type,
		})
	}

	return dependencies, nil
}

// createTaskFromTemplate creates a task using a template
func (s *TestDataSeeder) createTaskFromTemplate(
	template *gorev.GorevTemplate,
	values map[string]string,
	projeID string,
	priority string,
	sonTarih *time.Time,
	tags []string,
) (*gorev.Gorev, error) {
	// Build title from template
	title := values["title"]
	if title == "" && template.DefaultTitle != "" {
		title = s.processTemplate(template.DefaultTitle, values)
	}

	// Build description from template
	description := values["description"]
	if template.DescriptionTemplate != "" {
		description = s.processTemplate(template.DescriptionTemplate, values)
	}

	// Create task
	gorevObj := &gorev.Gorev{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      constants.TaskStatusPending,
		ProjeID:     projeID,
		WorkspaceID: s.config.WorkspaceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     sonTarih,
	}

	// Save task
	if err := s.isYonetici.VeriYonetici().GorevKaydet(s.ctx, gorevObj); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	// Add tags
	if len(tags) > 0 {
		etiketler, err := s.isYonetici.VeriYonetici().EtiketleriGetirVeyaOlustur(s.ctx, tags)
		if err != nil {
			return nil, fmt.Errorf("failed to create tags: %w", err)
		}
		if err := s.isYonetici.VeriYonetici().GorevEtiketleriniAyarla(s.ctx, gorevObj.ID, etiketler); err != nil {
			return nil, fmt.Errorf("failed to set task tags: %w", err)
		}
		gorevObj.Tags = etiketler
	}

	return gorevObj, nil
}

// createSubtaskFromTemplate creates a subtask using a template
func (s *TestDataSeeder) createSubtaskFromTemplate(
	template *gorev.GorevTemplate,
	values map[string]string,
	projeID string,
	parentID string,
	priority string,
) (*gorev.Gorev, error) {
	// Build title from template
	title := values["title"]
	if title == "" && template.DefaultTitle != "" {
		title = s.processTemplate(template.DefaultTitle, values)
	}

	// Build description from template
	description := values["description"]
	if template.DescriptionTemplate != "" {
		description = s.processTemplate(template.DescriptionTemplate, values)
	}

	// Create subtask using IsYonetici
	// AltGorevOlustur(ctx, parentID, baslik, aciklama, oncelik, sonTarihStr string, etiketIsimleri []string)
	subtask, err := s.isYonetici.AltGorevOlustur(s.ctx, parentID, title, description, priority, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create subtask: %w", err)
	}

	return subtask, nil
}

// processTemplate replaces placeholders in template strings
func (s *TestDataSeeder) processTemplate(templateStr string, values map[string]string) string {
	result := templateStr
	for key, value := range values {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}

// Summary returns a summary of seeded data
func (r *SeedResult) Summary() string {
	var sb strings.Builder
	sb.WriteString("=== Seed Result Summary ===\n")
	sb.WriteString(fmt.Sprintf("Projects:     %d\n", len(r.Projects)))
	sb.WriteString(fmt.Sprintf("Tasks:        %d\n", len(r.Tasks)))
	sb.WriteString(fmt.Sprintf("Subtasks:     %d\n", len(r.Subtasks)))
	sb.WriteString(fmt.Sprintf("Dependencies: %d\n", len(r.Dependencies)))
	return sb.String()
}
