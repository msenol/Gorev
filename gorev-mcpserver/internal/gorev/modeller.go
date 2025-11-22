package gorev

import "time"

// Gorev temel görev yapısı (task structure)
type Gorev struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	ProjeID     string     `json:"proje_id,omitempty"`
	ProjeName   string     `json:"proje_name,omitempty"` // Project name for Web UI/VS Code
	ParentID    string     `json:"parent_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []*Etiket  `json:"tags,omitempty"`
	Subtasks    []*Gorev   `json:"subtasks,omitempty"`
	Level       int        `json:"level,omitempty"`
	// Dependency counters - For TreeView display (omitempty removed - send 0 values too)
	DependencyCount            int `json:"dependency_count"`
	UncompletedDependencyCount int `json:"uncompleted_dependency_count"`
	DependentOnThisCount       int `json:"dependent_on_this_count"`
}

// Etiket görevleri kategorize etmek için kullanılır (tag for categorizing tasks)
type Etiket struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Proje görevleri gruplamak için kullanılır (project for grouping tasks)
type Proje struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Definition string    `json:"definition"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TaskCount  int       `json:"task_count"`
}

// Ozet sistem durumu özeti (summary of system status)
type Ozet struct {
	TotalProjects       int `json:"total_projects"`
	TotalTasks          int `json:"total_tasks"`
	PendingTasks        int `json:"pending_tasks"`
	InProgressTasks     int `json:"in_progress_tasks"`
	CompletedTasks      int `json:"completed_tasks"`
	HighPriorityTasks   int `json:"high_priority_tasks"`
	MediumPriorityTasks int `json:"medium_priority_tasks"`
	LowPriorityTasks    int `json:"low_priority_tasks"`
}

// Baglanti görevler arası bağlantı (connection between tasks)
type Baglanti struct {
	ID             string `json:"id"`
	SourceID       string `json:"source_id"`
	TargetID       string `json:"target_id"`
	ConnectionType string `json:"connection_type"`
}

// GorevTemplate görev oluşturma şablonu (task creation template)
type GorevTemplate struct {
	ID                  string            `json:"id"`
	Name                string            `json:"name"`
	Definition          string            `json:"definition"`
	Alias               string            `json:"alias"` // Short alias (e.g. bug, feature, research)
	DefaultTitle        string            `json:"default_title"`
	DescriptionTemplate string            `json:"description_template"`
	Fields              []TemplateAlan    `json:"fields"`
	SampleValues        map[string]string `json:"sample_values"`
	Category            string            `json:"category"`
	Active              bool              `json:"active"`
	// Multi-language support
	LanguageCode   string  `json:"language_code"`    // tr, en, etc.
	BaseTemplateID *string `json:"base_template_id"` // Groups templates by language
}

// TemplateAlan template'deki özelleştirilebilir alanlar (customizable template fields)
type TemplateAlan struct {
	Name     string   `json:"name"`
	Type     string   `json:"type"` // text, select, date, number
	Required bool     `json:"required"`
	Default  string   `json:"default"`
	Options  []string `json:"options,omitempty"`
}

// GorevHiyerarsi görev hiyerarşi bilgilerini tutar (task hierarchy information)
type GorevHiyerarsi struct {
	Gorev              *Gorev   `json:"gorev"`
	ParentTasks        []*Gorev `json:"parent_tasks,omitempty"`
	TotalSubtasks      int      `json:"total_subtasks"`
	CompletedSubtasks  int      `json:"completed_subtasks"`
	InProgressSubtasks int      `json:"in_progress_subtasks"`
	PendingSubtasks    int      `json:"pending_subtasks"`
	ProgressPercentage float64  `json:"progress_percentage"`
}

// Note: AIInteraction, AIContext, and FilterProfile structs are defined in their respective manager files
// to avoid circular dependencies and maintain clear ownership

// SearchHistory arama geçmişi (search history)
type SearchHistory struct {
	ID        string    `json:"id"`
	Query     string    `json:"query"`
	Results   int       `json:"results"`
	Timestamp time.Time `json:"timestamp"`
}

// FileWatch dosya izleme kaydı (file watch record)
type FileWatch struct {
	ID        string    `json:"id"`
	FilePath  string    `json:"file_path"`
	TaskID    string    `json:"task_id"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FileChange dosya değişikliği kaydı (file change record)
type FileChange struct {
	ID         string    `json:"id"`
	WatchID    string    `json:"watch_id"`
	ChangeType string    `json:"change_type"` // created, modified, deleted
	Timestamp  time.Time `json:"timestamp"`
}
