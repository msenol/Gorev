package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	testinghelpers "github.com/msenol/gorev/internal/testing"
	"github.com/spf13/cobra"
)

var (
	seedLang    string
	seedMinimal bool
	seedForce   bool
)

// createSeedCommand creates the seed-test-data CLI command
func createSeedCommand() *cobra.Command {
	seedCmd := &cobra.Command{
		Use:   "seed-test-data",
		Short: "Seed database with sample test data",
		Long: `Seed the database with realistic sample data for testing.

This command creates:
  - 3 sample projects (Mobil Uygulama, Backend API, Web Dashboard)
  - 15 sample tasks with various statuses and priorities
  - 3-level deep subtask hierarchies
  - Task dependencies
  - Tags/labels

Use --minimal for quick tests with only 3 tasks.`,
		Example: `  # Seed with full sample data (Turkish)
  gorev seed-test-data

  # Seed with English data
  gorev seed-test-data --lang=en

  # Seed with minimal data (3 tasks)
  gorev seed-test-data --minimal

  # Force re-seed (clear existing data first)
  gorev seed-test-data --force`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSeedTestData()
		},
	}

	seedCmd.Flags().StringVar(&seedLang, "lang", "tr", "Language for sample data (tr/en)")
	seedCmd.Flags().BoolVar(&seedMinimal, "minimal", false, "Create minimal test data (3 tasks)")
	seedCmd.Flags().BoolVar(&seedForce, "force", false, "Force re-seed (warning: clears existing tasks in project)")

	return seedCmd
}

// runSeedTestData executes the seed operation
func runSeedTestData() error {
	// Initialize i18n
	if seedLang != "" {
		if err := i18n.Initialize(seedLang); err != nil {
			return fmt.Errorf("failed to initialize i18n: %w", err)
		}
	}

	// Find or create database
	dbPath := findOrCreateTestDatabase()
	if dbPath == "" {
		return fmt.Errorf("could not determine database path")
	}

	fmt.Printf("🗄️  Using database: %s\n", dbPath)

	// Ensure directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Initialize database
	migrationsFS, err := getEmbeddedMigrationsFS()
	if err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	veriYonetici, err := gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, migrationsFS)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() { _ = veriYonetici.Kapat() }()

	// Create default templates if not exist
	if err := veriYonetici.VarsayilanTemplateleriOlustur(nil); err != nil {
		log.Printf("Warning: failed to create default templates: %v", err)
	}

	// Create IsYonetici
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Configure seeder
	config := &testinghelpers.SeederConfig{
		Language:        seedLang,
		WorkspaceID:     "", // Local mode
		IncludeSubtasks: !seedMinimal,
		IncludeDeps:     !seedMinimal,
		Minimal:         seedMinimal,
	}

	// Check if data already exists
	if !seedForce {
		existingTasks, err := isYonetici.GorevListele(nil, nil)
		if err == nil && len(existingTasks) > 0 {
			return fmt.Errorf("database already contains %d tasks. Use --force to overwrite", len(existingTasks))
		}
	}

	fmt.Printf("🌱 Seeding database with sample data...\n")
	fmt.Printf("   Language: %s\n", seedLang)
	fmt.Printf("   Mode: %s\n", map[bool]string{true: "minimal", false: "full"}[seedMinimal])

	// Create seeder and seed data
	seeder := testinghelpers.NewTestDataSeeder(isYonetici, config)
	result, err := seeder.SeedAll()
	if err != nil {
		return fmt.Errorf("seeding failed: %w", err)
	}

	// Print summary
	fmt.Printf("\n✅ Seeding completed successfully!\n")
	fmt.Printf(result.Summary())

	// Print project details
	fmt.Printf("\n📁 Projects created:\n")
	for i, p := range result.Projects {
		fmt.Printf("   %d. %s (ID: %s)\n", i+1, p.Name, p.ID)
	}

	// Print task summary by status
	statusCounts := make(map[string]int)
	for _, t := range result.Tasks {
		statusCounts[t.Status]++
	}
	fmt.Printf("\n📋 Tasks by status:\n")
	for status, count := range statusCounts {
		statusEmoji := map[string]string{
			"beklemede":    "⏳",
			"devam_ediyor": "🔄",
			"tamamlandi":   "✅",
			"iptal":        "❌",
		}[status]
		fmt.Printf("   %s %s: %d\n", statusEmoji, status, count)
	}

	fmt.Printf("\n🚀 You can now start the server with: gorev serve\n")
	fmt.Printf("🌐 Web UI will be available at: http://localhost:5082\n")

	return nil
}

// findOrCreateTestDatabase finds an existing database or determines where to create one
func findOrCreateTestDatabase() string {
	// 1. Check GOREV_DB_PATH environment variable
	if dbPath := os.Getenv("GOREV_DB_PATH"); dbPath != "" {
		return dbPath
	}

	// 2. Check for workspace database in current directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	workspaceDBPath := filepath.Join(cwd, ".gorev", "gorev.db")
	if _, err := os.Stat(workspaceDBPath); err == nil {
		return workspaceDBPath
	}

	// 3. Check home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		homeDBPath := filepath.Join(homeDir, ".gorev", "gorev.db")
		if _, err := os.Stat(homeDBPath); err == nil {
			return homeDBPath
		}
	}

	// 4. Create new workspace database in current directory
	return workspaceDBPath
}
