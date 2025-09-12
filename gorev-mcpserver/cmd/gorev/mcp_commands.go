package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/msenol/gorev/internal/constants"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/msenol/gorev/internal/mcp"
	"github.com/spf13/cobra"
)

var debugMode bool

// mcpCommand creates the root mcp command
func createMCPCommand() *cobra.Command {
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: i18n.T("cli.mcpTest"),
		Long:  i18n.T("cli.mcpDescription"),
	}

	// Global debug flag
	mcpCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, i18n.T("cli.debug"))

	// Add subcommands
	mcpCmd.AddCommand(
		createMCPListCommand(),
		createMCPCallCommand(),
		createMCPListTasksCommand(),
		createMCPCreateTaskCommand(),
		createMCPTaskDetailCommand(),
		createMCPProjectsCommand(),
	)

	return mcpCmd
}

// createMCPListCommand lists all available MCP tools
func createMCPListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: i18n.T("cli.list"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get database and migrations paths
			dbPath := getDatabasePath()
			migrationsPath := getMigrationsPath()

			if debugMode {
				fmt.Fprintf(os.Stderr, "Debug: Using database: %s\n", dbPath)
				fmt.Fprintf(os.Stderr, "Debug: Using migrations: %s\n", migrationsPath)
			}

			// Initialize managers
			veriYonetici, err := gorev.YeniVeriYonetici(dbPath, migrationsPath)
			if err != nil {
				return errors.New(i18n.T("error.dataManagerCreate", map[string]interface{}{"Error": err}))
			}
			defer func() { _ = veriYonetici.Kapat() }()

			// Get registered tools
			tools := mcp.ListTools()

			fmt.Println(i18n.T("display.availableTools"))
			fmt.Println("=" + strings.Repeat("=", 50))

			for _, tool := range tools {
				fmt.Printf("\n%s\n", tool.Name)
				fmt.Printf("  %s\n", i18n.T("display.toolDescription", map[string]interface{}{"Description": tool.Description}))

				if tool.InputSchema != nil {
					// Parse schema to show parameters
					if props, ok := tool.InputSchema["properties"].(map[string]interface{}); ok {
						fmt.Println("  Parametreler:")
						for param, schema := range props {
							if schemaMap, ok := schema.(map[string]interface{}); ok {
								paramType := schemaMap["type"]
								desc := schemaMap["description"]
								fmt.Printf("    - %s (%v): %v\n", param, paramType, desc)
							}
						}
					}

					if required, ok := tool.InputSchema["required"].([]interface{}); ok && len(required) > 0 {
						fmt.Print("  Zorunlu: ")
						for i, req := range required {
							if i > 0 {
								fmt.Print(", ")
							}
							fmt.Print(req)
						}
						fmt.Println()
					}
				}
			}

			return nil
		},
	}
}

// createMCPCallCommand calls a specific MCP tool
func createMCPCallCommand() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "call <tool> [param=value...]",
		Short: i18n.T("cli.call"),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolName := args[0]
			params := make(map[string]interface{})

			// Parse parameters
			for i := 1; i < len(args); i++ {
				parts := strings.SplitN(args[i], "=", 2)
				if len(parts) == 2 {
					key := parts[0]
					value := parts[1]

					// Special handling for degerler parameter (JSON object)
					if key == constants.ParamDegerler {
						var degerlerMap map[string]interface{}
						if err := json.Unmarshal([]byte(value), &degerlerMap); err == nil {
							params[key] = degerlerMap
						} else {
							return fmt.Errorf("degerler parametresi geçerli JSON objesi olmalı: %v", err)
						}
						continue
					}

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
			}

			// Get database and migrations paths
			dbPath := getDatabasePath()
			migrationsPath := getMigrationsPath()

			if debugMode {
				fmt.Fprintf(os.Stderr, "Debug: Using database: %s\n", dbPath)
				fmt.Fprintf(os.Stderr, "Debug: Using migrations: %s\n", migrationsPath)
			}

			// Initialize managers
			veriYonetici, err := gorev.YeniVeriYonetici(dbPath, migrationsPath)
			if err != nil {
				return errors.New(i18n.T("error.dataManagerCreate", map[string]interface{}{"Error": err}))
			}
			defer func() { _ = veriYonetici.Kapat() }()

			isYonetici := gorev.YeniIsYonetici(veriYonetici)
			handlers := mcp.YeniHandlers(isYonetici)

			// Call the tool
			result, err := handlers.CallTool(toolName, params)
			if err != nil {
				return errors.New(i18n.T("error.toolCallFailed", map[string]interface{}{"Error": err}))
			}

			// Display result
			if jsonOutput {
				jsonData, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
				return errors.New(i18n.T("error.jsonOutputFailed", map[string]interface{}{"Error": err}))
				}
				fmt.Println(string(jsonData))
			} else {
				// Pretty print for text content
				if len(result.Content) > 0 {
					for _, content := range result.Content {
						// Just convert to JSON and extract text
						jsonData, _ := json.Marshal(content)
						var contentMap map[string]interface{}
						if err := json.Unmarshal(jsonData, &contentMap); err == nil {
							if contentMap["type"] == "text" {
								if text, ok := contentMap["text"].(string); ok {
									fmt.Println(text)
								}
							}
						}
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, i18n.T("flags.jsonOutput"))
	return cmd
}

// Shortcut commands for common operations

func createMCPListTasksCommand() *cobra.Command {
	var allProjects bool
	var status string
	var limit int
	var offset int

	cmd := &cobra.Command{
		Use:   "list-tasks",
		Short: i18n.T("cli.listTasks"),
		RunE: func(cmd *cobra.Command, args []string) error {
			params := make(map[string]interface{})
			params["tum_projeler"] = allProjects
			if status != "" {
				params["durum"] = status
			}
			params["limit"] = fmt.Sprintf("%d", limit)
			params["offset"] = fmt.Sprintf("%d", offset)

			return callMCPTool("gorev_listele", params)
		},
	}

	cmd.Flags().BoolVar(&allProjects, "all-projects", true, i18n.T("flags.allProjects"))
	cmd.Flags().StringVar(&status, "status", "", "Durum filtresi (beklemede, devam_ediyor, tamamlandi)")
	cmd.Flags().IntVar(&limit, "limit", 50, i18n.T("flags.maxTasks"))
	cmd.Flags().IntVar(&offset, "offset", 0, i18n.T("flags.offset"))

	return cmd
}

func createMCPCreateTaskCommand() *cobra.Command {
	var templateID string
	var title string
	var description string
	var priority string
	var projectID string

	cmd := &cobra.Command{
		Use:   "create-task --template=<template-id-or-alias>",
		Short: i18n.T("cli.createTask"),
		Long:  "Create task from template. Use 'gorev template aliases' to see available shortcuts.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if templateID == "" {
				return fmt.Errorf("template ID veya alias gerekli. Kullanım: --template=bug veya --template=feature\n" +
					"Mevcut alias'ler için: gorev template aliases")
			}

			// Build template values
			degerler := make(map[string]interface{})

			if title != "" {
				degerler["baslik"] = title
			}
			if description != "" {
				degerler["aciklama"] = description
			}
			if priority != "" {
				degerler["oncelik"] = priority
			}
			if projectID != "" {
				degerler["proje_id"] = projectID
			}

			params := map[string]interface{}{
				constants.ParamTemplateID: templateID,
				constants.ParamDegerler:   degerler,
			}

			return callMCPTool("templateden_gorev_olustur", params)
		},
	}

	cmd.Flags().StringVar(&templateID, "template", "", "Template ID or alias (required). Use 'gorev template aliases' to see shortcuts")
	cmd.Flags().StringVar(&title, "title", "", i18n.T("flags.title"))
	cmd.Flags().StringVar(&description, "description", "", i18n.T("flags.taskDescription"))
	cmd.Flags().StringVar(&priority, "priority", constants.PriorityMedium, i18n.T("flags.priority"))
	cmd.Flags().StringVar(&projectID, "project", "", i18n.T("flags.projectId"))

	return cmd
}

func createMCPTaskDetailCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "task-detail <task-id>",
		Short: i18n.T("cli.showTask"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]interface{}{
				"id": args[0],
			}
			return callMCPTool("gorev_detay", params)
		},
	}
}

func createMCPProjectsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "projects",
		Short: i18n.T("cli.listProjects"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return callMCPTool("proje_listele", map[string]interface{}{})
		},
	}
}

// Helper function to call MCP tool
func callMCPTool(toolName string, params map[string]interface{}) error {
	if debugMode {
		fmt.Fprintf(os.Stderr, "Debug: Calling tool %s with params: %v\n", toolName, params)
	}
	// Get database and migrations paths
	dbPath := getDatabasePath()
	migrationsPath := getMigrationsPath()

	// Initialize managers
	veriYonetici, err := gorev.YeniVeriYonetici(dbPath, migrationsPath)
	if err != nil {
		return fmt.Errorf("veri yönetici oluşturulamadı: %w", err)
	}
			defer func() { _ = veriYonetici.Kapat() }()

	isYonetici := gorev.YeniIsYonetici(veriYonetici)
	handlers := mcp.YeniHandlers(isYonetici)

	// Call the tool
	result, err := handlers.CallTool(toolName, params)
	if err != nil {
		return fmt.Errorf("araç çağrısı başarısız: %w", err)
	}

	if debugMode {
		fmt.Fprintf(os.Stderr, "Debug: Tool returned result with %d content items\n", len(result.Content))
	}

	// Display result
	if len(result.Content) > 0 {
		for _, content := range result.Content {
			if debugMode {
				fmt.Fprintf(os.Stderr, "Debug: Content type: %T\n", content)
			}

			// Use reflection to get the Text field
			// This works with any struct that has a Text field
			contentValue := reflect.ValueOf(content)
			if contentValue.Kind() == reflect.Ptr {
				contentValue = contentValue.Elem()
			}

			if contentValue.Kind() == reflect.Struct {
				textField := contentValue.FieldByName("Text")
				if textField.IsValid() && textField.Kind() == reflect.String {
					fmt.Println(textField.String())
				}
			}
		}
	}

	return nil
}
