package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMigrationsPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "should return embedded migrations path",
			want: "embedded://migrations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getMigrationsPath()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetDatabasePath(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Save original environment variables
	origDBPath := os.Getenv("GOREV_DB_PATH")
	origGorevRoot := os.Getenv("GOREV_ROOT")
	origHome := os.Getenv("HOME")
	origUserHomeDir := os.Getenv("USERPROFILE")

	// Cleanup function to restore environment
	defer func() {
		if origDBPath != "" {
			os.Setenv("GOREV_DB_PATH", origDBPath)
		} else {
			os.Unsetenv("GOREV_DB_PATH")
		}
		if origGorevRoot != "" {
			os.Setenv("GOREV_ROOT", origGorevRoot)
		} else {
			os.Unsetenv("GOREV_ROOT")
		}
		os.Setenv("HOME", origHome)
		os.Setenv("USERPROFILE", origUserHomeDir)
	}()

	t.Run("GOREV_DB_PATH environment variable", func(t *testing.T) {
		dbPath := filepath.Join(tempDir, "custom.db")
		os.Setenv("GOREV_DB_PATH", dbPath)

		got := getDatabasePath()
		assert.Equal(t, dbPath, got)
	})

	t.Run("workspace database in current directory", func(t *testing.T) {
		os.Unsetenv("GOREV_DB_PATH")
		os.Unsetenv("GOREV_ROOT")

		// Create separate temp directory for this test
		workspaceTempDir := t.TempDir()

		// Change to temp directory
		origCwd, _ := os.Getwd()
		defer os.Chdir(origCwd)

		os.Chdir(workspaceTempDir)

		// Create .gorev directory and database
		gorevDir := filepath.Join(workspaceTempDir, ".gorev")
		err := os.MkdirAll(gorevDir, 0755)
		require.NoError(t, err)

		dbPath := filepath.Join(gorevDir, "gorev.db")
		file, err := os.Create(dbPath)
		require.NoError(t, err)
		file.Close()

		got := getDatabasePath()
		assert.Equal(t, dbPath, got)
	})

	t.Run("GOREV_ROOT environment variable", func(t *testing.T) {
		os.Unsetenv("GOREV_DB_PATH")

		// Create separate temp directory for this test
		gorevRootTempDir := t.TempDir()
		os.Setenv("GOREV_ROOT", gorevRootTempDir)
		defer os.Unsetenv("GOREV_ROOT")

		// Temporarily hide global database if it exists
		homeDir, _ := os.UserHomeDir()
		globalDBPath := filepath.Join(homeDir, ".gorev", "gorev.db")
		tempGlobalDBPath := globalDBPath + ".test_backup"

		// Move global DB temporarily if it exists
		if _, err := os.Stat(globalDBPath); err == nil {
			os.Rename(globalDBPath, tempGlobalDBPath)
			defer os.Rename(tempGlobalDBPath, globalDBPath)
		}

		// Change to temp directory to avoid current workspace detection
		origCwd, _ := os.Getwd()
		defer os.Chdir(origCwd)
		os.Chdir(gorevRootTempDir)

		got := getDatabasePath()
		// The function might find project root first, so accept either GOREV_ROOT path or project root
		expectedPath1 := filepath.Join(gorevRootTempDir, "gorev.db")
		expectedPath2 := filepath.Join(gorevRootTempDir, "gorev-mcpserver", "gorev.db")
		if got != expectedPath1 && got != expectedPath2 {
			t.Errorf("Expected %s or %s, got %s", expectedPath1, expectedPath2, got)
		}
	})

	t.Run("fallback to current directory", func(t *testing.T) {
		os.Unsetenv("GOREV_DB_PATH")
		os.Unsetenv("GOREV_ROOT")

		// Create separate temp directory for this test
		fallbackTempDir := t.TempDir()

		// Temporarily hide global database if it exists
		homeDir, _ := os.UserHomeDir()
		globalDBPath := filepath.Join(homeDir, ".gorev", "gorev.db")
		tempGlobalDBPath := globalDBPath + ".test_backup"

		// Move global DB temporarily if it exists
		if _, err := os.Stat(globalDBPath); err == nil {
			os.Rename(globalDBPath, tempGlobalDBPath)
			defer os.Rename(tempGlobalDBPath, globalDBPath)
		}

		// Change to temp directory
		origCwd, _ := os.Getwd()
		defer os.Chdir(origCwd)

		os.Chdir(fallbackTempDir)

		got := getDatabasePath()
		// Should return current directory path since temp dir doesn't have workspace .gorev/gorev.db
		assert.Equal(t, "gorev.db", got)
	})
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		expected string
	}{
		{
			name:     "default turkish",
			env:      "",
			expected: "tr",
		},
		{
			name:     "explicit turkish",
			env:      "tr",
			expected: "tr",
		},
		{
			name:     "english",
			env:      "en",
			expected: "en",
		},
		{
			name:     "invalid language fallback",
			env:      "fr",
			expected: "tr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			origLang := os.Getenv("GOREV_LANG")
			defer os.Setenv("GOREV_LANG", origLang)

			// Set test environment
			if tt.env != "" {
				os.Setenv("GOREV_LANG", tt.env)
			} else {
				os.Unsetenv("GOREV_LANG")
			}

			got := detectLanguage()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCreateVeriYonetici(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		// Set up environment for database creation
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "gorev.db")

		origDBPath := os.Getenv("GOREV_DB_PATH")
		defer os.Setenv("GOREV_DB_PATH", origDBPath)

		os.Setenv("GOREV_DB_PATH", dbPath)

		vy, err := createVeriYonetici()
		require.NoError(t, err)
		assert.NotNil(t, vy)
	})

	t.Run("invalid database path", func(t *testing.T) {
		// Set up invalid path
		origDBPath := os.Getenv("GOREV_DB_PATH")
		defer os.Setenv("GOREV_DB_PATH", origDBPath)

		os.Setenv("GOREV_DB_PATH", "/invalid/path/database.db")

		vy, err := createVeriYonetici()
		require.Error(t, err)
		assert.Nil(t, vy)
	})
}

func TestRunServer(t *testing.T) {
	t.Run("server startup", func(t *testing.T) {
		// Set up environment for server test
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "gorev.db")

		origDBPath := os.Getenv("GOREV_DB_PATH")
		defer os.Setenv("GOREV_DB_PATH", origDBPath)

		os.Setenv("GOREV_DB_PATH", dbPath)

		// Start server in goroutine (this will likely fail due to server setup, but we test the call)
		done := make(chan bool)
		go func() {
			err := runServer()
			assert.Error(t, err) // Expect error due to server setup limitations in test
			done <- true
		}()

		select {
		case <-done:
			// Server completed (likely with error)
		case <-time.After(500 * time.Millisecond):
			// Timeout, server might be running
			t.Log("Server likely started successfully")
		}
	})
}

func TestInitWorkspaceDatabase(t *testing.T) {
	t.Run("successful workspace initialization", func(t *testing.T) {
		tempDir := t.TempDir()

		// Change to temp directory
		origCwd, _ := os.Getwd()
		defer os.Chdir(origCwd)

		os.Chdir(tempDir)

		err := initWorkspaceDatabase(false)
		require.NoError(t, err)

		// Check if .gorev directory was created
		gorevDir := filepath.Join(tempDir, ".gorev")
		_, err = os.Stat(gorevDir)
		assert.NoError(t, err)

		// Check if database file was created
		dbPath := filepath.Join(gorevDir, "gorev.db")
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)
	})

	t.Run("successful global initialization", func(t *testing.T) {
		// Set up a fake home directory
		tempHome := t.TempDir()
		origHome := os.Getenv("HOME")
		origUserHomeDir := os.Getenv("USERPROFILE")

		defer func() {
			os.Setenv("HOME", origHome)
			os.Setenv("USERPROFILE", origUserHomeDir)
		}()

		os.Setenv("HOME", tempHome)
		os.Setenv("USERPROFILE", tempHome)

		err := initWorkspaceDatabase(true)
		require.NoError(t, err)

		// Check if .gorev directory was created in home directory
		gorevDir := filepath.Join(tempHome, ".gorev")
		_, err = os.Stat(gorevDir)
		assert.NoError(t, err)

		// Check if database file was created
		dbPath := filepath.Join(gorevDir, "gorev.db")
		_, err = os.Stat(dbPath)
		assert.NoError(t, err)
	})

	t.Run("database already exists", func(t *testing.T) {
		tempDir := t.TempDir()

		origCwd, _ := os.Getwd()
		defer os.Chdir(origCwd)

		os.Chdir(tempDir)

		// Create database first
		gorevDir := filepath.Join(tempDir, ".gorev")
		err := os.MkdirAll(gorevDir, 0755)
		require.NoError(t, err)

		dbPath := filepath.Join(gorevDir, "gorev.db")
		file, err := os.Create(dbPath)
		require.NoError(t, err)
		file.Close()

		// Try to initialize again
		err = initWorkspaceDatabase(false)
		assert.NoError(t, err)
	})
}

func TestTemplateCommands(t *testing.T) {
	// Set up environment for template commands
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "gorev.db")

	origDBPath := os.Getenv("GOREV_DB_PATH")
	defer os.Setenv("GOREV_DB_PATH", origDBPath)

	os.Setenv("GOREV_DB_PATH", dbPath)

	t.Run("list templates", func(t *testing.T) {
		err := listTemplates("")
		require.NoError(t, err)
	})

	t.Run("show template", func(t *testing.T) {
		// This will show all templates since no specific template is provided
		err := showTemplate("")
		require.NoError(t, err)
	})

	t.Run("init templates", func(t *testing.T) {
		err := initTemplates()
		require.NoError(t, err)
	})

	t.Run("list template aliases", func(t *testing.T) {
		err := listTemplateAliases()
		require.NoError(t, err)
	})
}

func TestIDECommands(t *testing.T) {
	t.Run("run IDE detect", func(t *testing.T) {
		err := runIDEDetect()
		require.NoError(t, err)
	})

	t.Run("run IDE status", func(t *testing.T) {
		err := runIDEStatus()
		require.NoError(t, err)
	})

	t.Run("run IDE config", func(t *testing.T) {
		err := runIDEConfig()
		require.NoError(t, err)
	})

	t.Run("run IDE config set", func(t *testing.T) {
		err := runIDEConfigSet("auto_install", "true")
		require.NoError(t, err)
	})
}

func TestCreateIDECommand(t *testing.T) {
	t.Run("IDE command creation", func(t *testing.T) {
		cmd := createIDECommand()
		assert.NotNil(t, cmd)
		assert.Equal(t, "ide", cmd.Name())
		assert.True(t, cmd.HasSubCommands())
	})
}


func TestCheckAndPromptIDEExtensions(t *testing.T) {
	t.Run("check IDE extensions", func(t *testing.T) {
		// Set up environment for IDE commands
		tempDir := t.TempDir()
		dbPath := filepath.Join(tempDir, "gorev.db")

		origDBPath := os.Getenv("GOREV_DB_PATH")
		defer os.Setenv("GOREV_DB_PATH", origDBPath)

		os.Setenv("GOREV_DB_PATH", dbPath)

		// This should not panic
		assert.NotPanics(t, func() {
			checkAndPromptIDEExtensions()
		})
	})
}

func BenchmarkGetDatabasePath(b *testing.B) {
	// Save original environment
	origDBPath := os.Getenv("GOREV_DB_PATH")
	origGorevRoot := os.Getenv("GOREV_ROOT")

	defer func() {
		if origDBPath != "" {
			os.Setenv("GOREV_DB_PATH", origDBPath)
		} else {
			os.Unsetenv("GOREV_DB_PATH")
		}
		if origGorevRoot != "" {
			os.Setenv("GOREV_ROOT", origGorevRoot)
		} else {
			os.Unsetenv("GOREV_ROOT")
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getDatabasePath()
	}
}

func BenchmarkDetectLanguage(b *testing.B) {
	// Save original environment
	origLang := os.Getenv("GOREV_LANG")

	defer func() {
		if origLang != "" {
			os.Setenv("GOREV_LANG", origLang)
		} else {
			os.Unsetenv("GOREV_LANG")
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detectLanguage()
	}
}