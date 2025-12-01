package testing

import (
	"context"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/require"
)

// TestDatabaseConfig configures test database setup
type TestDatabaseConfig struct {
	UseMemoryDB     bool          // Use :memory: database
	UseTempFile     bool          // Use temporary file database
	CustomPath      string        // Use custom database path
	MigrationsPath  string        // Path to migration files (used if MigrationsFS is nil)
	MigrationsFS    fs.FS         // Embedded migrations filesystem (preferred over MigrationsPath)
	CreateTemplates bool          // Create default templates
	InitializeI18n  bool          // Initialize i18n system
	SeedTestData    bool          // Seed test data after database setup
	SeederConfig    *SeederConfig // Configuration for test data seeder (used if SeedTestData is true)
}

// DefaultTestDatabaseConfig returns default configuration for test databases
func DefaultTestDatabaseConfig() *TestDatabaseConfig {
	return &TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPathMCP, // internal/testing is same level as internal/gorev
		CreateTemplates: true,
		InitializeI18n:  true,
	}
}

// SetupTestDatabase creates a test database with the given configuration
func SetupTestDatabase(t *testing.T, config *TestDatabaseConfig) (*gorev.VeriYonetici, func()) {
	// Initialize i18n if requested
	if config.InitializeI18n {
		if !i18n.IsInitialized() {
			err := i18n.Initialize(constants.DefaultTestLanguage)
			if err != nil {
				t.Logf("Warning: i18n initialization failed: %v", err)
			}
		}
	}

	// Determine database path
	var dbPath string
	var cleanup func()

	switch {
	case config.UseMemoryDB:
		dbPath = constants.TestDatabaseURI
		cleanup = func() {} // No cleanup needed for in-memory database
	case config.UseTempFile:
		dbPath = "test_" + strings.ReplaceAll(time.Now().Format("2006-01-02T15:04:05.000000000Z"), ":", "-") + ".db"
		cleanup = func() {
			os.Remove(dbPath)
		}
	case config.CustomPath != "":
		dbPath = config.CustomPath
		cleanup = func() {} // Custom path, no automatic cleanup
	default:
		dbPath = constants.TestDatabaseURI
		cleanup = func() {}
	}

	// Create database manager
	var veriYonetici *gorev.VeriYonetici
	var err error

	// Use embedded migrations if available, otherwise use file path
	if config.MigrationsFS != nil {
		veriYonetici, err = gorev.YeniVeriYoneticiWithEmbeddedMigrations(dbPath, config.MigrationsFS)
	} else {
		veriYonetici, err = gorev.YeniVeriYonetici(dbPath, config.MigrationsPath)
	}
	require.NoError(t, err)

	// Create default templates if requested
	if config.CreateTemplates {
		err = veriYonetici.VarsayilanTemplateleriOlustur(context.Background())
		require.NoError(t, err)
	}

	// Wrap cleanup to include database close
	originalCleanup := cleanup
	cleanup = func() {
		veriYonetici.Kapat()
		originalCleanup()
	}

	return veriYonetici, cleanup
}

// SetupTestEnvironmentBasic creates a basic test environment with database and business logic
func SetupTestEnvironmentBasic(t *testing.T) (*gorev.IsYonetici, func()) {
	veriYonetici, cleanup := SetupTestDatabase(t, DefaultTestDatabaseConfig())
	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	return isYonetici, cleanup
}

// SetupTestEnvironmentWithConfig creates a test environment with custom database configuration
func SetupTestEnvironmentWithConfig(t *testing.T, config *TestDatabaseConfig) (*gorev.IsYonetici, func()) {
	veriYonetici, cleanup := SetupTestDatabase(t, config)
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Seed test data if requested
	if config.SeedTestData {
		seederConfig := config.SeederConfig
		if seederConfig == nil {
			seederConfig = DefaultSeederConfig()
		}
		seeder := NewTestDataSeeder(isYonetici, seederConfig)
		_, err := seeder.SeedAll()
		require.NoError(t, err, "failed to seed test data")
	}

	return isYonetici, cleanup
}
