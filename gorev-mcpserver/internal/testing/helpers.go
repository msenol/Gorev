package testing

import (
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
	UseMemoryDB     bool   // Use :memory: database
	UseTempFile     bool   // Use temporary file database
	CustomPath      string // Use custom database path
	MigrationsPath  string // Path to migration files
	CreateTemplates bool   // Create default templates
	InitializeI18n  bool   // Initialize i18n system
}

// DefaultTestDatabaseConfig returns default configuration for test databases
func DefaultTestDatabaseConfig() *TestDatabaseConfig {
	return &TestDatabaseConfig{
		UseMemoryDB:     true,
		MigrationsPath:  constants.TestMigrationsPathIntegration,
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
	veriYonetici, err := gorev.YeniVeriYonetici(dbPath, config.MigrationsPath)
	require.NoError(t, err)

	// Create default templates if requested
	if config.CreateTemplates {
		err = veriYonetici.VarsayilanTemplateleriOlustur()
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
	return isYonetici, cleanup
}
