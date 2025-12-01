package testing

import (
	"context"
	"testing"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTestDataSeeder(t *testing.T) {
	isYonetici, cleanup := SetupTestEnvironmentBasic(t)
	defer cleanup()

	t.Run("with nil config uses defaults", func(t *testing.T) {
		seeder := NewTestDataSeeder(isYonetici, nil)
		require.NotNil(t, seeder)
		assert.Equal(t, "tr", seeder.config.Language)
		assert.True(t, seeder.config.IncludeSubtasks)
		assert.True(t, seeder.config.IncludeDeps)
		assert.False(t, seeder.config.Minimal)
	})

	t.Run("with custom config", func(t *testing.T) {
		config := &SeederConfig{
			Language:        "en",
			IncludeSubtasks: false,
			IncludeDeps:     false,
			Minimal:         true,
		}
		seeder := NewTestDataSeeder(isYonetici, config)
		require.NotNil(t, seeder)
		assert.Equal(t, "en", seeder.config.Language)
		assert.False(t, seeder.config.IncludeSubtasks)
		assert.False(t, seeder.config.IncludeDeps)
		assert.True(t, seeder.config.Minimal)
	})
}

func TestSeederSeedProjects(t *testing.T) {
	t.Run("creates full projects", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language: "tr",
			Minimal:  false,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		projects, err := seeder.SeedProjects()
		require.NoError(t, err)
		assert.Len(t, projects, 3)

		// Verify project names
		assert.Equal(t, "Mobil Uygulama", projects[0].Name)
		assert.Equal(t, "Backend API", projects[1].Name)
		assert.Equal(t, "Web Dashboard", projects[2].Name)
	})

	t.Run("creates minimal project", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language: "tr",
			Minimal:  true,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		projects, err := seeder.SeedProjects()
		require.NoError(t, err)
		assert.Len(t, projects, 1)
		assert.Equal(t, "Test Projesi", projects[0].Name)
	})

	t.Run("creates English projects", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		// Initialize English i18n
		err := i18n.Initialize("en")
		require.NoError(t, err)

		config := &SeederConfig{
			Language: "en",
			Minimal:  false,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		projects, err := seeder.SeedProjects()
		require.NoError(t, err)
		assert.Len(t, projects, 3)

		// Verify English project names
		assert.Equal(t, "Mobile App", projects[0].Name)
		assert.Equal(t, "Backend API", projects[1].Name)
		assert.Equal(t, "Web Dashboard", projects[2].Name)
	})
}

func TestSeederSeedTasks(t *testing.T) {
	t.Run("creates full tasks", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language: "tr",
			Minimal:  false,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		// First create projects
		projects, err := seeder.SeedProjects()
		require.NoError(t, err)

		// Then create tasks
		tasks, err := seeder.SeedTasks(projects)
		require.NoError(t, err)
		assert.Len(t, tasks, 15)

		// Verify status distribution
		statusCount := make(map[string]int)
		for _, task := range tasks {
			statusCount[task.Status]++
		}
		assert.Equal(t, 5, statusCount[constants.TaskStatusPending])
		assert.Equal(t, 4, statusCount[constants.TaskStatusInProgress])
		assert.Equal(t, 5, statusCount[constants.TaskStatusCompleted])
		assert.Equal(t, 1, statusCount[constants.TaskStatusCancelled])
	})

	t.Run("creates minimal tasks", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language: "tr",
			Minimal:  true,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		// First create projects
		projects, err := seeder.SeedProjects()
		require.NoError(t, err)

		// Then create tasks
		tasks, err := seeder.SeedTasks(projects)
		require.NoError(t, err)
		assert.Len(t, tasks, 3)

		// Verify status distribution
		statusCount := make(map[string]int)
		for _, task := range tasks {
			statusCount[task.Status]++
		}
		assert.Equal(t, 1, statusCount[constants.TaskStatusPending])
		assert.Equal(t, 1, statusCount[constants.TaskStatusInProgress])
		assert.Equal(t, 1, statusCount[constants.TaskStatusCompleted])
	})
}

func TestSeederSeedSubtasks(t *testing.T) {
	t.Run("creates subtask hierarchy", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language:        "tr",
			IncludeSubtasks: true,
			Minimal:         false,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		// First create projects and tasks
		projects, err := seeder.SeedProjects()
		require.NoError(t, err)

		tasks, err := seeder.SeedTasks(projects)
		require.NoError(t, err)

		// Then create subtasks
		subtasks, err := seeder.SeedSubtasks(tasks)
		require.NoError(t, err)
		assert.Greater(t, len(subtasks), 0)

		// Verify subtasks have parent IDs
		for _, subtask := range subtasks {
			assert.NotEmpty(t, subtask.ParentID, "subtask should have parent ID")
		}
	})
}

func TestSeederSeedDependencies(t *testing.T) {
	t.Run("creates dependencies", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language:    "tr",
			IncludeDeps: true,
			Minimal:     false,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		// First create projects and tasks
		projects, err := seeder.SeedProjects()
		require.NoError(t, err)

		tasks, err := seeder.SeedTasks(projects)
		require.NoError(t, err)

		// Then create dependencies
		deps, err := seeder.SeedDependencies(tasks)
		require.NoError(t, err)
		assert.Len(t, deps, 3)

		// Verify dependencies have correct structure
		for _, dep := range deps {
			assert.NotEmpty(t, dep.SourceID)
			assert.NotEmpty(t, dep.TargetID)
			assert.NotEmpty(t, dep.Type)
		}
	})
}

func TestSeederSeedAll(t *testing.T) {
	t.Run("seeds all data with full config", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language:        "tr",
			IncludeSubtasks: true,
			IncludeDeps:     true,
			Minimal:         false,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		result, err := seeder.SeedAll()
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify counts
		assert.Len(t, result.Projects, 3)
		assert.Len(t, result.Tasks, 15)
		assert.Greater(t, len(result.Subtasks), 0)
		assert.Len(t, result.Dependencies, 3)
	})

	t.Run("seeds minimal data", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language:        "tr",
			IncludeSubtasks: true,
			IncludeDeps:     true,
			Minimal:         true,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		result, err := seeder.SeedAll()
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify minimal counts
		assert.Len(t, result.Projects, 1)
		assert.Len(t, result.Tasks, 3)
		assert.Len(t, result.Subtasks, 0) // Minimal mode skips subtasks
		assert.Len(t, result.Dependencies, 0) // Minimal mode skips dependencies
	})

	t.Run("sets first project as active", func(t *testing.T) {
		isYonetici, cleanup := SetupTestEnvironmentBasic(t)
		defer cleanup()

		config := &SeederConfig{
			Language: "tr",
			Minimal:  true,
		}
		seeder := NewTestDataSeeder(isYonetici, config)

		result, err := seeder.SeedAll()
		require.NoError(t, err)
		require.NotNil(t, result)

		// Check active project is set
		ctx := context.Background()
		activeProje, err := isYonetici.AktifProjeGetir(ctx)
		require.NoError(t, err)
		require.NotNil(t, activeProje)
		assert.Equal(t, result.Projects[0].ID, activeProje.ID)
	})
}

func TestSeedResultSummary(t *testing.T) {
	// Test with empty result
	emptyResult := &SeedResult{}
	summary := emptyResult.Summary()
	assert.Contains(t, summary, "Seed Result Summary")
	assert.Contains(t, summary, "Projects:")
	assert.Contains(t, summary, "Tasks:")
	assert.Contains(t, summary, "Subtasks:")
	assert.Contains(t, summary, "Dependencies:")

	// Verify format includes counts
	assert.Contains(t, summary, "0")
}
