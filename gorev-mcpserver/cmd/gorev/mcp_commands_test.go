package main

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
)

// Test helper to capture stdout/stderr
func captureOutput(f func()) (string, string) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = wOut
	os.Stderr = wErr

	outChan := make(chan string)
	errChan := make(chan string)

	go func() {
		var buf strings.Builder
		_, _ = io.Copy(&buf, rOut)
		outChan <- buf.String()
	}()

	go func() {
		var buf strings.Builder
		_, _ = io.Copy(&buf, rErr)
		errChan <- buf.String()
	}()

	f()

	_ = wOut.Close()
	_ = wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	stdout := <-outChan
	stderr := <-errChan

	return stdout, stderr
}

// Setup test environment
func setupTestEnvironment(t *testing.T) {
	// Initialize i18n for tests
	i18n.Initialize(constants.DefaultTestLanguage)
}

// Test the main MCP command creation
func TestCreateMCPCommand(t *testing.T) {
	setupTestEnvironment(t)

	cmd := createMCPCommand()

	assert.Equal(t, "mcp", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Check that subcommands are added (check exact Use strings)
	expectedSubcommands := []string{"list", "call <tool> [param=value...]", "list-tasks", "create-task --template=<template-id-or-alias>", "task-detail <task-id>", "projects"}
	actualSubcommands := make([]string, len(cmd.Commands()))
	for i, subcmd := range cmd.Commands() {
		actualSubcommands[i] = subcmd.Use
	}

	for _, expected := range expectedSubcommands {
		assert.Contains(t, actualSubcommands, expected, "Missing subcommand: %s", expected)
	}
}

// Test parameter parsing logic directly (unit test approach)
func TestParameterParsingLogic(t *testing.T) {
	setupTestEnvironment(t)

	// Test the parameter parsing logic that happens inside createMCPCallCommand
	// This simulates the parsing logic from lines 118-151 in mcp_commands.go

	tests := []struct {
		name        string
		paramString string
		expectedKey string
		expectedVal interface{}
		expectError bool
	}{
		{
			name:        "Simple string parameter",
			paramString: "param1=value1",
			expectedKey: "param1",
			expectedVal: "value1",
			expectError: false,
		},
		{
			name:        "Boolean true parameter",
			paramString: "bool_param=true",
			expectedKey: "bool_param",
			expectedVal: true,
			expectError: false,
		},
		{
			name:        "Boolean false parameter",
			paramString: "bool_param=false",
			expectedKey: "bool_param",
			expectedVal: false,
			expectError: false,
		},
		{
			name:        "Valid JSON values parameter",
			paramString: `values={"title":"Test","description":"Description"}`,
			expectedKey: constants.ParamValues,
			expectedVal: map[string]interface{}{"title": "Test", "description": "Description"},
			expectError: false,
		},
		{
			name:        "Invalid JSON values parameter",
			paramString: `values={"title":"Test","description":invalid}`,
			expectedKey: constants.ParamValues,
			expectedVal: nil,
			expectError: true,
		},
		{
			name:        "Nested JSON parameter",
			paramString: `complex={"nested":{"value":"test"}}`,
			expectedKey: "complex",
			expectedVal: map[string]interface{}{"nested": map[string]interface{}{"value": "test"}},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parameter parsing logic
			parts := strings.SplitN(tt.paramString, "=", 2)
			if len(parts) != 2 {
				if !tt.expectError {
					t.Errorf("Expected successful parsing but got invalid format")
				}
				return
			}

			key := parts[0]
			value := parts[1]
			params := make(map[string]interface{})

			// Special handling for degerler parameter (JSON object)
			if key == constants.ParamValues {
				var degerlerMap map[string]interface{}
				if err := json.Unmarshal([]byte(value), &degerlerMap); err != nil {
					if !tt.expectError {
						t.Errorf("Expected valid JSON but got error: %v", err)
					}
					return
				} else {
					params[key] = degerlerMap
				}
			} else {
				// Convert boolean strings
				switch value {
				case "true":
					params[key] = true
				case "false":
					params[key] = false
				default:
					// Try to parse as JSON for nested objects
					var jsonValue interface{}
					if err := json.Unmarshal([]byte(value), &jsonValue); err == nil {
						params[key] = jsonValue
					} else {
						params[key] = value
					}
				}
			}

			// Verify results
			assert.Equal(t, tt.expectedKey, key)
			if !tt.expectError {
				assert.Equal(t, tt.expectedVal, params[key])
			}
		})
	}
}

// Test parameter parsing edge cases with direct logic testing
func TestParameterParsingEdgeCases(t *testing.T) {
	setupTestEnvironment(t)

	tests := []struct {
		name        string
		paramString string
		shouldParse bool
		description string
	}{
		{
			name:        "Parameter without equals sign",
			paramString: "invalid_param",
			shouldParse: false,
			description: "Should skip parameters without equals sign",
		},
		{
			name:        "Empty parameter value",
			paramString: "param=",
			shouldParse: true,
			description: "Should handle empty parameter values",
		},
		{
			name:        "Parameter with multiple equals signs",
			paramString: "param=value=with=equals",
			shouldParse: true,
			description: "Should split on first equals only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the parsing logic directly
			parts := strings.SplitN(tt.paramString, "=", 2)

			if tt.shouldParse {
				assert.Len(t, parts, 2, "Should successfully split parameter")
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]
					assert.NotEmpty(t, key, "Key should not be empty")
					// Value can be empty, that's fine

					// Test that we can create a parameter map
					params := make(map[string]interface{})
					params[key] = value
					assert.Contains(t, params, key)
				}
			} else {
				assert.Len(t, parts, 1, "Should not split parameter without equals")
			}
		})
	}
}

// Test shortcut commands structure
func TestMCPListTasksCommand(t *testing.T) {
	setupTestEnvironment(t)

	cmd := createMCPListTasksCommand()

	assert.Equal(t, "list-tasks", cmd.Use)
	assert.NotEmpty(t, cmd.Short)

	// Check flags
	assert.True(t, cmd.Flags().Lookup("all-projects") != nil)
	assert.True(t, cmd.Flags().Lookup("status") != nil)
	assert.True(t, cmd.Flags().Lookup("limit") != nil)
	assert.True(t, cmd.Flags().Lookup("offset") != nil)
}

func TestMCPCreateTaskCommand(t *testing.T) {
	setupTestEnvironment(t)

	cmd := createMCPCreateTaskCommand()

	assert.Equal(t, "create-task --template=<template-id-or-alias>", cmd.Use)
	assert.NotEmpty(t, cmd.Short)

	// Check flags
	assert.True(t, cmd.Flags().Lookup("title") != nil)
	assert.True(t, cmd.Flags().Lookup("description") != nil)
	assert.True(t, cmd.Flags().Lookup("priority") != nil)
	assert.True(t, cmd.Flags().Lookup("project") != nil)
}

func TestMCPTaskDetailCommand(t *testing.T) {
	setupTestEnvironment(t)

	cmd := createMCPTaskDetailCommand()

	assert.Equal(t, "task-detail <task-id>", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotNil(t, cmd.Args, "Command should have Args validation")
}

func TestMCPProjectsCommand(t *testing.T) {
	setupTestEnvironment(t)

	cmd := createMCPProjectsCommand()

	assert.Equal(t, "projects", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
}

// Test JSON output flag functionality
func TestMCPCallCommandJSONOutputFlag(t *testing.T) {
	setupTestEnvironment(t)

	cmd := createMCPCallCommand()

	// Test that the JSON flag exists and can be set
	assert.True(t, cmd.Flags().Lookup("json") != nil)

	err := cmd.Flags().Set("json", "true")
	assert.NoError(t, err)

	flagValue, err := cmd.Flags().GetBool("json")
	assert.NoError(t, err)
	assert.True(t, flagValue)
}

// Test debug mode global flag
func TestDebugModeFlag(t *testing.T) {
	setupTestEnvironment(t)

	cmd := createMCPCommand()

	// Test that debug flag exists
	assert.True(t, cmd.PersistentFlags().Lookup("debug") != nil)

	// Test setting debug flag
	err := cmd.PersistentFlags().Set("debug", "true")
	assert.NoError(t, err)

	flagValue, err := cmd.PersistentFlags().GetBool("debug")
	assert.NoError(t, err)
	assert.True(t, flagValue)
}

// Test specific degerler parameter JSON parsing logic (isolated)
func TestDegerlerParameterParsing(t *testing.T) {
	setupTestEnvironment(t)

	tests := []struct {
		name        string
		jsonValue   string
		expectError bool
		description string
	}{
		{
			name:        "Valid simple JSON",
			jsonValue:   `{"baslik":"Test Task","oncelik":"yuksek"}`,
			expectError: false,
			description: "Should parse valid simple JSON",
		},
		{
			name:        "Valid complex JSON",
			jsonValue:   `{"baslik":"Complex Task","aciklama":"Long description","oncelik":"yuksek","tags":["tag1","tag2"]}`,
			expectError: false,
			description: "Should parse valid complex JSON",
		},
		{
			name:        "Invalid JSON - missing quotes",
			jsonValue:   `{baslik:"Test Task"}`,
			expectError: true,
			description: "Should fail on invalid JSON syntax",
		},
		{
			name:        "Invalid JSON - trailing comma",
			jsonValue:   `{"baslik":"Test Task","oncelik":"yuksek",}`,
			expectError: true,
			description: "Should fail on trailing comma",
		},
		{
			name:        "Empty JSON object",
			jsonValue:   `{}`,
			expectError: false,
			description: "Should accept empty JSON object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var degerlerMap map[string]interface{}
			err := json.Unmarshal([]byte(tt.jsonValue), &degerlerMap)

			if tt.expectError {
				assert.Error(t, err, "Expected JSON parsing to fail for: %s", tt.jsonValue)
			} else {
				assert.NoError(t, err, "Expected JSON parsing to succeed for: %s", tt.jsonValue)
				assert.NotNil(t, degerlerMap)
			}
		})
	}
}

// Test the template creation command integration (structure only, no actual execution)
func TestTemplateBasedTaskCreationCLIIntegration(t *testing.T) {
	setupTestEnvironment(t)

	// Test that we can create the command and it has the expected structure
	cmd := createMCPCallCommand()

	// Simulate the args that would be used for template-based task creation
	testArgs := []string{
		"templateden_gorev_olustur",
		"template_id=basic_task",
		`degerler={"baslik":"CLI Test Task","aciklama":"Created via CLI","oncelik":"yuksek"}`,
	}

	// We're not actually executing this since it would require a database,
	// but we can test that the command structure and arguments are valid
	assert.NotNil(t, cmd.RunE, "Command should have a RunE function")
	assert.Equal(t, "call <tool> [param=value...]", cmd.Use)
	assert.True(t, len(testArgs) >= 1, "Should accept the required number of arguments")

	// Test JSON parsing for the degerler parameter specifically
	degerlerJSON := `{"baslik":"CLI Test Task","aciklama":"Created via CLI","oncelik":"yuksek"}`
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(degerlerJSON), &parsed)
	assert.NoError(t, err, "degerler JSON should be valid")
	assert.Equal(t, "CLI Test Task", parsed["baslik"])
	assert.Equal(t, "Created via CLI", parsed["aciklama"])
	assert.Equal(t, "yuksek", parsed["oncelik"])
}
