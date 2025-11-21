package api

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/msenol/gorev/internal/mcp"
	ws "github.com/msenol/gorev/internal/websocket"
)

// handleMCPToolCall handles MCP tool calls forwarded from proxy
func (s *APIServer) handleMCPToolCall(c *fiber.Ctx) error {
	// Extract method from path (supports slash-separated methods like "tools/list")
	path := c.Path()
	prefix := "/api/v1/mcp/"
	if !strings.HasPrefix(path, prefix) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid MCP path",
		})
	}
	toolName := strings.TrimPrefix(path, prefix)

	// Get workspace context from middleware
	workspaceID := c.Locals("workspace_id")

	if workspaceID == nil || workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing workspace context",
		})
	}

	// Get workspace context from manager
	wsCtx, err := s.workspaceManager.GetWorkspaceContext(workspaceID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get workspace context: %v", err),
		})
	}

	// Debug: Log EventEmitter type
	fmt.Fprintf(os.Stderr, "[MCP Bridge] Tool: %s, Workspace: %s, EventEmitter type: %T\n",
		toolName, wsCtx.ID, wsCtx.EventEmitter)

	// Get business logic manager from workspace context
	isYonetici := wsCtx.IsYonetici

	// Create MCP handlers for this workspace
	handlers := mcp.YeniHandlers(isYonetici)

	// Parse request body as MCP tool parameters
	var params map[string]interface{}
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
	}

	// Call the appropriate MCP tool handler
	result, err := s.dispatchMCPTool(handlers, toolName, params, wsCtx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Tool execution failed: %v", err),
		})
	}

	return c.JSON(result)
}

// dispatchMCPTool dispatches MCP tool call to appropriate handler
func (s *APIServer) dispatchMCPTool(handlers *mcp.Handlers, toolName string, params map[string]interface{}, wsCtx *WorkspaceContext) (interface{}, error) {
	// Map tool name to handler method
	// Note: Handler methods use PascalCase (GorevListele, not HandleGorevListele)
	// All 41 MCP tools are supported
	var result interface{}
	var err error

	switch toolName {
	// Task management tools (6)
	case "gorev_listele":
		result, err = handlers.GorevListele(params)
	case "gorev_detay":
		result, err = handlers.GorevDetay(params)
	case "gorev_guncelle":
		result, err = handlers.GorevGuncelle(params)
		if err == nil {
			// Emit task updated event
			if taskID, ok := params["id"].(string); ok {
				wsCtx.EventEmitter.EmitTaskUpdated(wsCtx.ID, taskID, params)
			}
		}
	case "gorev_duzenle":
		result, err = handlers.GorevDuzenle(params)
		if err == nil {
			// Emit task updated event
			if taskID, ok := params["id"].(string); ok {
				wsCtx.EventEmitter.EmitTaskUpdated(wsCtx.ID, taskID, params)
			}
		}
	case "gorev_sil":
		result, err = handlers.GorevSil(params)
		if err == nil {
			// Emit task deleted event
			if taskID, ok := params["id"].(string); ok {
				wsCtx.EventEmitter.EmitTaskDeleted(wsCtx.ID, taskID)
			}
		}

	// Unified hierarchy handler
	case "gorev_hierarchy":
		result, err = handlers.GorevHierarchy(params)
		if err == nil {
			action, _ := params["action"].(string)
			if action == "create_subtask" {
				// Extract subtask ID and emit task created event
				if taskID := extractTaskIDFromResult(result); taskID != "" {
					wsCtx.EventEmitter.EmitTaskCreated(wsCtx.ID, taskID, params)
				} else {
					wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
				}
			} else {
				wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
			}
		}

	// Dependency (kept separate - frequently used)
	case "gorev_bagimlilik_ekle":
		result, err = handlers.GorevBagimlilikEkle(params)
		if err == nil {
			wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
		}

	// Template tools (2)
	case "template_listele":
		result, err = handlers.TemplateListele(params)
	case "templateden_gorev_olustur":
		result, err = handlers.TemplatedenGorevOlustur(params)
		if err == nil {
			// Extract task ID from result and emit task created event
			if taskID := extractTaskIDFromResult(result); taskID != "" {
				fmt.Fprintf(os.Stderr, "[MCP Bridge] ✅ Emitting task_created event: workspace=%s, taskID=%s\n", wsCtx.ID, taskID)
				wsCtx.EventEmitter.EmitTaskCreated(wsCtx.ID, taskID, params)
			} else {
				fmt.Fprintf(os.Stderr, "[MCP Bridge] ⚠ Could not extract task ID, emitting workspace_sync instead\n")
				// Fallback to workspace sync if ID extraction fails
				wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
			}
		}

	// Project tools (6)
	case "proje_olustur":
		result, err = handlers.ProjeOlustur(params)
		if err == nil {
			// Extract project ID and emit project created event
			if projectID := extractTaskIDFromResult(result); projectID != "" {
				wsCtx.EventEmitter.EmitProjectCreated(wsCtx.ID, projectID, params)
			} else {
				wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
			}
		}
	case "proje_listele":
		result, err = handlers.ProjeListele(params)
	case "proje_gorevleri":
		result, err = handlers.ProjeGorevleri(params)

	// Unified active project handler
	case "aktif_proje":
		result, err = handlers.AktifProje(params)
		if err == nil {
			action, _ := params["action"].(string)
			if action == "set" || action == "clear" {
				wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
			}
		}

	// Summary tool (1)
	case "ozet_goster":
		result, err = handlers.OzetGoster(params)

	// Unified bulk operations handler
	case "gorev_bulk":
		result, err = handlers.GorevBulk(params)
		if err == nil {
			wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
		}

	// Unified filter profile handler
	case "gorev_filter_profile":
		result, err = handlers.GorevFilterProfile(params)

	// Unified file watch handler
	case "gorev_file_watch":
		result, err = handlers.GorevFileWatch(params)

	// Unified IDE management handler
	case "ide_manage":
		result, err = handlers.IDEManage(params)

	// Unified AI context handler
	case "gorev_context":
		result, err = handlers.GorevContext(params)
		if err == nil {
			action, _ := params["action"].(string)
			if action == "set_active" {
				wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
			}
		}

	// Unified search handler
	case "gorev_search":
		result, err = handlers.GorevSearch(params)

	// Data export/import (kept separate - distinct operations)
	case "gorev_export":
		result, err = handlers.GorevExport(params)
	case "gorev_import":
		result, err = handlers.GorevImport(params)
		if err == nil {
			wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
		}

	// AI advanced (kept separate - complex operations)
	case "gorev_suggestions":
		result, err = handlers.GorevSuggestions(params)
	case "gorev_intelligent_create":
		result, err = handlers.GorevIntelligentCreate(params)
		if err == nil {
			// Extract task ID and emit task created event
			if taskID := extractTaskIDFromResult(result); taskID != "" {
				wsCtx.EventEmitter.EmitTaskCreated(wsCtx.ID, taskID, params)
			} else {
				wsCtx.EventEmitter.EmitWorkspaceSync(wsCtx.ID)
			}
		}

	// MCP Protocol methods
	case "initialize":
		// Return proper MCP initialize response
		result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "gorev",
				"version": "0.16.1",
			},
		}
		err = nil

	case "tools/list":
		// Return list of 24 optimized MCP tools (reduced from 45)
		tools := []map[string]interface{}{
			// === CORE TOOLS (11) ===
			// Task CRUD
			{"name": "gorev_listele", "description": "List and filter tasks", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"durum": map[string]interface{}{"type": "string"}, "limit": map[string]interface{}{"type": "number"}}}},
			{"name": "gorev_detay", "description": "Show task details", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "string"}}, "required": []string{"id"}}},
			{"name": "gorev_guncelle", "description": "Update task fields", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "string"}, "durum": map[string]interface{}{"type": "string"}, "oncelik": map[string]interface{}{"type": "string"}}, "required": []string{"id"}}},
			{"name": "gorev_duzenle", "description": "Edit task content", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "string"}, "title": map[string]interface{}{"type": "string"}, "description": map[string]interface{}{"type": "string"}}, "required": []string{"id"}}},
			{"name": "gorev_sil", "description": "Delete task", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"id": map[string]interface{}{"type": "string"}}, "required": []string{"id"}}},

			// Templates
			{"name": "template_listele", "description": "List templates", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}},
			{"name": "templateden_gorev_olustur", "description": "Create task from template", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"template_id": map[string]interface{}{"type": "string", "description": "Template ID or alias (bug, feature, research)"}, "degerler": map[string]interface{}{"type": "object"}}, "required": []string{"template_id", "degerler"}}},

			// Projects
			{"name": "proje_listele", "description": "List projects", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}},
			{"name": "proje_olustur", "description": "Create project", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"isim": map[string]interface{}{"type": "string"}, "tanim": map[string]interface{}{"type": "string"}}, "required": []string{"isim"}}},
			{"name": "proje_gorevleri", "description": "List project tasks", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"proje_id": map[string]interface{}{"type": "string"}}, "required": []string{"proje_id"}}},
			{"name": "gorev_bagimlilik_ekle", "description": "Add task dependency", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"gorev_id": map[string]interface{}{"type": "string"}, "bagli_gorev_id": map[string]interface{}{"type": "string"}}, "required": []string{"gorev_id", "bagli_gorev_id"}}},

			// === UNIFIED TOOLS (8) ===
			{"name": "aktif_proje", "description": "Manage active project (unified: set|get|clear)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"action": map[string]interface{}{"type": "string", "enum": []string{"set", "get", "clear"}}, "proje_id": map[string]interface{}{"type": "string"}}, "required": []string{"action"}}},
			{"name": "gorev_hierarchy", "description": "Manage task hierarchy (unified: create_subtask|change_parent|show)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"action": map[string]interface{}{"type": "string", "enum": []string{"create_subtask", "change_parent", "show"}}, "parent_id": map[string]interface{}{"type": "string"}, "title": map[string]interface{}{"type": "string"}}, "required": []string{"action"}}},
			{"name": "gorev_bulk", "description": "Bulk operations (unified: transition|tag|update)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"operation": map[string]interface{}{"type": "string", "enum": []string{"transition", "tag", "update"}}, "ids": map[string]interface{}{"type": "array"}, "data": map[string]interface{}{"type": "object"}}, "required": []string{"operation", "ids"}}},
			{"name": "gorev_filter_profile", "description": "Filter profiles (unified: save|load|list|delete)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"action": map[string]interface{}{"type": "string", "enum": []string{"save", "load", "list", "delete"}}, "name": map[string]interface{}{"type": "string"}}, "required": []string{"action"}}},
			{"name": "gorev_file_watch", "description": "File watching (unified: add|remove|list|stats)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"action": map[string]interface{}{"type": "string", "enum": []string{"add", "remove", "list", "stats"}}, "file_path": map[string]interface{}{"type": "string"}}, "required": []string{"action"}}},
			{"name": "ide_manage", "description": "IDE management (unified: detect|install|uninstall|status|update)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"action": map[string]interface{}{"type": "string", "enum": []string{"detect", "install", "uninstall", "status", "update"}}, "ide": map[string]interface{}{"type": "string"}}, "required": []string{"action"}}},
			{"name": "gorev_context", "description": "AI context (unified: set_active|get_active|recent|summary)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"action": map[string]interface{}{"type": "string", "enum": []string{"set_active", "get_active", "recent", "summary"}}, "gorev_id": map[string]interface{}{"type": "string"}}, "required": []string{"action"}}},
			{"name": "gorev_search", "description": "Search tasks (unified: nlp|advanced|history)", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"mode": map[string]interface{}{"type": "string", "enum": []string{"nlp", "advanced", "history"}}, "query": map[string]interface{}{"type": "string"}, "arama_metni": map[string]interface{}{"type": "string"}}, "required": []string{"mode"}}},

			// === SPECIAL TOOLS (5) ===
			{"name": "ozet_goster", "description": "Show workspace summary", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{}}},
			{"name": "gorev_export", "description": "Export tasks", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"format": map[string]interface{}{"type": "string"}}}},
			{"name": "gorev_import", "description": "Import tasks", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"data": map[string]interface{}{"type": "object"}}, "required": []string{"data"}}},
			{"name": "gorev_suggestions", "description": "Get AI task suggestions", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"context": map[string]interface{}{"type": "string"}}}},
			{"name": "gorev_intelligent_create", "description": "AI-powered task creation", "inputSchema": map[string]interface{}{"type": "object", "properties": map[string]interface{}{"title": map[string]interface{}{"type": "string", "description": "Task title"}, "description": map[string]interface{}{"type": "string", "description": "Task description"}, "auto_split": map[string]interface{}{"type": "boolean", "description": "Auto-split into subtasks"}, "estimate_time": map[string]interface{}{"type": "boolean", "description": "Estimate task duration"}, "smart_priority": map[string]interface{}{"type": "boolean", "description": "AI-suggested priority"}, "suggest_template": map[string]interface{}{"type": "boolean", "description": "Suggest matching template"}, "proje_id": map[string]interface{}{"type": "string", "description": "Project ID"}}, "required": []string{"title"}}},
		}
		result = map[string]interface{}{
			"tools": tools,
		}
		err = nil

	case "notifications/initialized":
		// Client notification that initialization is complete
		result = map[string]interface{}{}
		err = nil

	case "resources/list":
		// Return empty resources list (we don't support resources yet)
		result = map[string]interface{}{
			"resources": []interface{}{},
		}
		err = nil

	case "resources/templates/list":
		// Return empty resource templates list
		result = map[string]interface{}{
			"resourceTemplates": []interface{}{},
		}
		err = nil

	case "tools/call":
		// Execute a tool call - body should have "name" and "arguments"
		toolCallName, ok := params["name"].(string)
		if !ok {
			return nil, fmt.Errorf(i18n.T("error.toolsCallNameRequired"))
		}

		// Get arguments (optional)
		var toolArgs map[string]interface{}
		if args, ok := params["arguments"].(map[string]interface{}); ok {
			toolArgs = args
		} else {
			toolArgs = make(map[string]interface{})
		}

		// Recursively call dispatchMCPTool with the actual tool name
		fmt.Fprintf(os.Stderr, "[MCP Bridge] tools/call executing: %s with args: %v\n", toolCallName, toolArgs)
		return s.dispatchMCPTool(handlers, toolCallName, toolArgs, wsCtx)

	default:
		return nil, fmt.Errorf(i18n.T("error.unknownTool", map[string]interface{}{"Tool": toolName}))
	}

	return result, err
}

// extractTaskIDFromResult extracts task ID from MCP tool result
// Returns empty string if ID cannot be extracted
func extractTaskIDFromResult(result interface{}) string {
	// MCP result is *mcpgo.CallToolResult with Content slice
	if toolResult, ok := result.(*mcpgo.CallToolResult); ok {
		for _, content := range toolResult.Content {
			// Type assert to TextContent
			if textContent, ok := content.(mcpgo.TextContent); ok {
				// Extract ID using regex pattern "ID: <uuid>"
				// Pattern matches: ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
				re := regexp.MustCompile(`ID:\s*([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)
				if matches := re.FindStringSubmatch(textContent.Text); len(matches) > 1 {
					return matches[1]
				}
			}
		}
	}
	return ""
}

// registerMCPBridgeRoutes registers MCP bridge endpoints
func (s *APIServer) registerMCPBridgeRoutes() {
	// MCP bridge endpoint - handles all MCP tool calls from proxy
	// Use wildcard to support slash-separated methods like "tools/list"
	s.app.Post("/api/v1/mcp/*", s.handleMCPToolCall)
}

// registerWebSocketRoutes registers WebSocket endpoints for real-time updates
func (s *APIServer) registerWebSocketRoutes() {
	// WebSocket endpoint - handles real-time updates
	// Example: ws://localhost:5082/ws?workspace_id=abc123
	s.app.Get("/ws", ws.WebSocketUpgradeMiddleware(), ws.HandleWebSocket(s.wsHub))

	// WebSocket stats endpoint (for monitoring)
	s.app.Get("/api/v1/ws/stats", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"total_clients": s.wsHub.GetClientCount(),
			"timestamp":     serverStartTime.Unix(),
		})
	})
}
