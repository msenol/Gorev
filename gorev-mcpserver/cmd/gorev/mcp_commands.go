package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

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
				return fmt.Errorf(i18n.T("error.dataManagerCreate", map[string]interface{}{"Error": err}))
			}
			defer veriYonetici.Kapat()

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
					// Convert boolean strings
					switch parts[1] {
					case "true":
						params[parts[0]] = true
					case "false":
						params[parts[0]] = false
					default:
						params[parts[0]] = parts[1]
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
				return fmt.Errorf(i18n.T("error.dataManagerCreate", map[string]interface{}{"Error": err}))
			}
			defer veriYonetici.Kapat()

			isYonetici := gorev.YeniIsYonetici(veriYonetici)
			handlers := mcp.YeniHandlers(isYonetici)

			// Call the tool
			result, err := handlers.CallTool(toolName, params)
			if err != nil {
				return fmt.Errorf(i18n.T("error.toolCallFailed", map[string]interface{}{"Error": err}))
			}

			// Display result
			if jsonOutput {
				jsonData, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf(i18n.T("error.jsonOutputFailed", map[string]interface{}{"Error": err}))
				}
				fmt.Println(string(jsonData))
			} else {
				// Pretty print for text content
				if result.Content != nil && len(result.Content) > 0 {
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
	var title string
	var description string
	var priority string
	var projectID string

	cmd := &cobra.Command{
		Use:   "create-task",
		Short: i18n.T("cli.createTask"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf(i18n.T("error.titleRequired"))
			}

			params := map[string]interface{}{
				"baslik":   title,
				"aciklama": description,
				"oncelik":  priority,
			}
			if projectID != "" {
				params["proje_id"] = projectID
			}

			return callMCPTool("gorev_olustur", params)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", i18n.T("flags.title"))
	cmd.Flags().StringVar(&description, "description", "", i18n.T("flags.taskDescription"))
	cmd.Flags().StringVar(&priority, "priority", "orta", i18n.T("flags.priority"))
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
	defer veriYonetici.Kapat()

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
	if result.Content != nil && len(result.Content) > 0 {
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
