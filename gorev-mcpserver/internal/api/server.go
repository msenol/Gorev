package api

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/msenol/gorev/internal/api/middleware"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	ws "github.com/msenol/gorev/internal/websocket"
)

// APIServer represents the HTTP API server
type APIServer struct {
	app              *fiber.App
	port             string
	isYonetici       *gorev.IsYonetici // Legacy single workspace manager (deprecated)
	workspaceManager *WorkspaceManager // Multi-workspace manager
	handlers         interface{}       // MCP Handlers for export/import operations
	wsHub            *ws.Hub           // WebSocket hub for real-time updates
}

// SetMigrationsFS sets the embedded migrations filesystem for workspace manager
func (s *APIServer) SetMigrationsFS(migrationsFS fs.FS) {
	s.workspaceManager.SetMigrationsFS(migrationsFS)
}

// NewAPIServer creates a new API server instance
func NewAPIServer(port string, isYonetici *gorev.IsYonetici) *APIServer {
	if port == "" {
		port = "5082"
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
		DisableStartupMessage: false,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))

	// CORS for localhost development only
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5000,http://localhost:5001,http://localhost:5002,http://localhost:5003", // Restrict to localhost only
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,X-Workspace-Id,X-Workspace-Path,X-Workspace-Name", // Add workspace headers
		AllowCredentials: false,
	}))

	// Create WebSocket hub
	wsHub := ws.NewHub()

	// Create workspace manager with hub reference
	workspaceManager := NewWorkspaceManager()
	workspaceManager.wsHub = wsHub

	server := &APIServer{
		app:              app,
		port:             port,
		isYonetici:       isYonetici,
		workspaceManager: workspaceManager,
		wsHub:            wsHub,
	}

	// Start WebSocket hub in background
	go wsHub.Run()

	// Workspace detection middleware (must be after CORS, before routes)
	app.Use(middleware.WorkspaceMiddleware(server.workspaceManager))

	// Setup routes
	server.setupRoutes()

	// Register MCP bridge routes (for proxy support)
	server.registerMCPBridgeRoutes()

	// Register WebSocket routes
	server.registerWebSocketRoutes()

	return server
}

// Global variable to track server start time
var serverStartTime = time.Now()

// App returns the underlying Fiber app instance
func (s *APIServer) App() *fiber.App {
	return s.app
}

// SetHandlers sets the MCP handlers for export/import operations
func (s *APIServer) SetHandlers(handlers interface{}) {
	s.handlers = handlers
}

// setupRoutes configures all API routes
func (s *APIServer) setupRoutes() {
	// Health check endpoint (for daemon detection, no /v1 prefix)
	s.app.Get("/api/health", func(c *fiber.Ctx) error {
		uptime := time.Since(serverStartTime)
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": "v0.16.1", // TODO: get from build variable
			"uptime":  uptime.Seconds(),
			"time":    time.Now().Unix(),
		})
	})

	api := s.app.Group("/api/v1")

	// Health check (legacy, v1 prefix)
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// Task routes
	api.Get("/tasks", s.getTasks)
	api.Post("/tasks", s.createTask)
	api.Get("/tasks/:id", s.getTask)
	api.Put("/tasks/:id", s.updateTask)
	api.Delete("/tasks/:id", s.deleteTask)
	api.Post("/tasks/from-template", s.createTaskFromTemplate)

	// Project routes
	api.Get("/projects", s.getProjects)
	api.Post("/projects", s.createProject)
	api.Get("/projects/:id", s.getProject)
	api.Get("/projects/:id/tasks", s.getProjectTasks)
	api.Put("/projects/:id/activate", s.activateProject)

	// Template routes
	api.Get("/templates", s.getTemplates)

	// Summary routes
	api.Get("/summary", s.getSummary)

	// Language routes
	api.Get("/language", s.getLanguage)
	api.Post("/language", s.setLanguage)

	// Subtask routes
	api.Post("/tasks/:id/subtasks", s.createSubtask)
	api.Put("/tasks/:id/parent", s.changeParent)
	api.Get("/tasks/:id/hierarchy", s.getHierarchy)

	// Dependency routes
	api.Post("/tasks/:id/dependencies", s.addDependency)
	api.Delete("/tasks/:id/dependencies/:dep_id", s.removeDependency)

	// Active project routes
	api.Get("/active-project", s.getActiveProject)

	// Workspace routes
	api.Post("/workspaces/register", s.registerWorkspaceHandler)
	api.Get("/workspaces", s.listWorkspacesHandler)
	api.Get("/workspaces/:id", s.getWorkspaceHandler)
	api.Delete("/workspaces/:id", s.unregisterWorkspaceHandler)
	api.Delete("/active-project", s.removeActiveProject)

	// Export/Import routes
	api.Post("/export", s.exportData)
	api.Post("/import", s.importData)

	// MCP Protocol routes (for AI assistant integration)
	api.Post("/mcp/*", s.handleMCPToolCall)

	// Note: Static web UI is served by ServeStaticFiles() called from main.go after server creation
}

// Start starts the API server
func (s *APIServer) Start() error {
	log.Printf("🚀 API Server starting on port %s", s.port)
	log.Printf("📱 Web UI: http://localhost:%s", s.port)
	log.Printf("🔧 API: http://localhost:%s/api/v1", s.port)

	return s.app.Listen(":" + s.port)
}

// StartAsync starts the API server in a goroutine
func (s *APIServer) StartAsync() {
	go func() {
		if err := s.Start(); err != nil {
			log.Printf("❌ API Server error: %v", err)
		}
	}()
}

// Shutdown gracefully stops the API server
func (s *APIServer) Shutdown(ctx context.Context) error {
	log.Println("🔽 Shutting down API server...")
	return s.app.ShutdownWithContext(ctx)
}

// Handler methods

// getTasks retrieves all tasks with optional filtering
func (s *APIServer) getTasks(c *fiber.Ctx) error {
	// Create filters map based on query parameters
	filters := make(map[string]interface{})

	// Handle standard query parameters
	if durum := c.Query("durum"); durum != "" {
		filters["durum"] = durum
	}
	if tumProjeler := c.QueryBool("tum_projeler"); tumProjeler {
		filters["tum_projeler"] = true
	}
	if sirala := c.Query("sirala"); sirala != "" {
		filters["sirala"] = sirala
	}
	if filtre := c.Query("filtre"); filtre != "" {
		filters["filtre"] = filtre
	}
	if etiket := c.Query("etiket"); etiket != "" {
		filters["etiket"] = etiket
	}
	if limit := c.QueryInt("limit", 50); limit > 0 {
		filters["limit"] = limit
	}
	if offset := c.QueryInt("offset", 0); offset >= 0 {
		filters["offset"] = offset
	}

	// Call business logic with workspace context
	iy := s.getIsYoneticiFromContext(c)
	gorevler, err := iy.GorevListele(filters)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to list tasks with filters %v: %v", filters, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    gorevler,
		"total":   len(gorevler),
	})
}

// getTask retrieves a specific task by ID
func (s *APIServer) getTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Task ID is required")
	}

	// Get task details using MCP handler logic with workspace context
	iy := s.getIsYoneticiFromContext(c)
	gorev, err := iy.VeriYonetici().GorevGetir(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to get task with ID %s: %v", id, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    gorev,
	})
}

// createTask creates a new task using template (since direct creation is deprecated)
func (s *APIServer) createTask(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented, "Direct task creation is deprecated. Use /tasks/from-template endpoint with a template.")
}

// updateTask updates an existing task
func (s *APIServer) updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Task ID is required")
	}

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Handle status update - use GorevGuncelle with params map and workspace context
	iy := s.getIsYoneticiFromContext(c)
	if durum, ok := req["durum"].(string); ok {
		params := map[string]interface{}{"durum": durum}
		if err := iy.VeriYonetici().GorevGuncelle(id, params); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update task %s with status %s: %v", id, durum, err))
		}
	}

	// Get updated task
	gorev, err := iy.VeriYonetici().GorevGetir(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to get updated task with ID %s: %v", id, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    gorev,
		"message": "Task updated successfully",
	})
}

// deleteTask deletes a task
func (s *APIServer) deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Task ID is required")
	}

	iy := s.getIsYoneticiFromContext(c)
	if err := iy.VeriYonetici().GorevSil(id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete task with ID %s: %v", id, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Task deleted successfully",
	})
}

// createTaskFromTemplate creates a task from a template
func (s *APIServer) createTaskFromTemplate(c *fiber.Ctx) error {
	var req struct {
		TemplateID string            `json:"template_id"`
		Degerler   map[string]string `json:"degerler"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.TemplateID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Template ID is required")
	}

	// Create task from template using business logic with workspace context
	iy := s.getIsYoneticiFromContext(c)
	gorev, err := iy.TemplatedenGorevOlustur(req.TemplateID, req.Degerler)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create task from template %s: %v", req.TemplateID, err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    gorev,
		"message": "Task created from template successfully",
	})
}

// getIsYoneticiFromContext extracts workspace-specific IsYonetici from Fiber context
// Falls back to global isYonetici if workspace context is not available (backward compatibility)
func (s *APIServer) getIsYoneticiFromContext(c *fiber.Ctx) *gorev.IsYonetici {
	isYonetici := middleware.GetIsYonetici(c)
	if isYonetici == nil {
		log.Printf("[getIsYoneticiFromContext] No workspace context, using global isYonetici for %s %s",
			c.Method(), c.Path())
		return s.isYonetici
	}

	iy, ok := isYonetici.(*gorev.IsYonetici)
	if !ok {
		log.Printf("[getIsYoneticiFromContext] Type assertion failed, using global isYonetici for %s %s",
			c.Method(), c.Path())
		return s.isYonetici
	}

	wsID := middleware.GetWorkspaceID(c)
	log.Printf("[getIsYoneticiFromContext] Using workspace-specific isYonetici for workspace %s (%s %s)",
		wsID, c.Method(), c.Path())
	return iy
}

// getProjects retrieves all projects
func (s *APIServer) getProjects(c *fiber.Ctx) error {
	iy := s.getIsYoneticiFromContext(c)

	projeler, err := iy.ProjeListele()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to list projects: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    projeler,
		"total":   len(projeler),
	})
}

// getProject retrieves a specific project by ID
func (s *APIServer) getProject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Project ID is required")
	}

	iy := s.getIsYoneticiFromContext(c)
	proje, err := iy.VeriYonetici().ProjeGetir(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to get project with ID %s: %v", id, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    proje,
	})
}

// createProject creates a new project
func (s *APIServer) createProject(c *fiber.Ctx) error {
	var req struct {
		Name       string `json:"isim"`
		Definition string `json:"tanim"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Isim is required")
	}

	iy := s.getIsYoneticiFromContext(c)

	// Create project using business logic
	proje, err := iy.ProjeOlustur(req.Name, req.Definition)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create project '%s': %v", req.Name, err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    proje,
		"message": "Project created successfully",
	})
}

// getProjectTasks retrieves all tasks for a specific project
func (s *APIServer) getProjectTasks(c *fiber.Ctx) error {
	projeID := c.Params("id")
	if projeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Project ID is required")
	}

	filters := map[string]interface{}{
		"proje_id": projeID,
		"limit":    c.QueryInt("limit", 50),
		"offset":   c.QueryInt("offset", 0),
	}

	iy := s.getIsYoneticiFromContext(c)
	gorevler, err := iy.GorevListele(filters)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to list tasks for project %s: %v", projeID, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    gorevler,
		"total":   len(gorevler),
	})
}

// activateProject activates a project
func (s *APIServer) activateProject(c *fiber.Ctx) error {
	projeID := c.Params("id")
	if projeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Project ID is required")
	}

	iy := s.getIsYoneticiFromContext(c)
	if err := iy.VeriYonetici().AktifProjeAyarla(projeID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to activate project with ID %s: %v", projeID, err))
	}

	// Return updated project
	proje, err := iy.VeriYonetici().ProjeGetir(projeID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get activated project with ID %s: %v", projeID, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    proje,
		"message": "Project activated successfully",
	})
}

// getTemplates retrieves all templates with optional category filtering
func (s *APIServer) getTemplates(c *fiber.Ctx) error {
	kategori := c.Query("kategori")

	iy := s.getIsYoneticiFromContext(c)
	templateler, err := iy.VeriYonetici().TemplateListele(kategori)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to list templates with category '%s': %v", kategori, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    templateler,
		"total":   len(templateler),
	})
}

// getSummary retrieves system-wide summary statistics
func (s *APIServer) getSummary(c *fiber.Ctx) error {
	// This would need a summary method in business logic
	// For now, return basic info
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Summary endpoint - to be implemented",
		},
	})
}

// getLanguage retrieves the current language setting
func (s *APIServer) getLanguage(c *fiber.Ctx) error {
	currentLang := i18n.GetCurrentLanguage()
	return c.JSON(fiber.Map{
		"success":  true,
		"language": currentLang,
	})
}

// setLanguage changes the current language setting
func (s *APIServer) setLanguage(c *fiber.Ctx) error {
	type LanguageRequest struct {
		Language string `json:"language"`
	}

	var req LanguageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate language code
	if req.Language != "tr" && req.Language != "en" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Language must be 'tr' or 'en'",
		})
	}

	// Set the language
	if err := i18n.SetLanguage(req.Language); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("Failed to set language: %v", err),
		})
	}

	log.Printf("🌍 Language changed to: %s", req.Language)

	return c.JSON(fiber.Map{
		"success":  true,
		"language": req.Language,
		"message":  fmt.Sprintf("Language changed to %s", req.Language),
	})
}

// createSubtask creates a subtask under a parent task
func (s *APIServer) createSubtask(c *fiber.Ctx) error {
	parentID := c.Params("id")
	if parentID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Parent task ID is required")
	}

	var req struct {
		Title       string `json:"baslik"`
		Description string `json:"aciklama"`
		Priority    string `json:"oncelik"`
		DueDate     string `json:"son_tarih"`
		Tags        string `json:"etiketler"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.Title == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Title is required")
	}

	// Parse etiketler string to []string
	var etiketler []string
	if req.Tags != "" {
		// Simple split by comma
		for _, e := range strings.Split(req.Tags, ",") {
			etiketler = append(etiketler, strings.TrimSpace(e))
		}
	}

	// Create subtask using business logic with workspace context
	iy := s.getIsYoneticiFromContext(c)
	gorev, err := iy.AltGorevOlustur(parentID, req.Title, req.Description, req.Priority, req.DueDate, etiketler)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create subtask under parent %s: %v", parentID, err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    gorev,
		"message": "Subtask created successfully",
	})
}

// changeParent changes the parent of a task
func (s *APIServer) changeParent(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Task ID is required")
	}

	var req struct {
		NewParentID string `json:"new_parent_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Change parent using business logic with workspace context
	iy := s.getIsYoneticiFromContext(c)
	if err := iy.GorevUstDegistir(taskID, req.NewParentID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to change parent for task %s: %v", taskID, err))
	}

	// Get updated task
	gorev, err := iy.VeriYonetici().GorevGetir(taskID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to get updated task with ID %s: %v", taskID, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    gorev,
		"message": "Parent changed successfully",
	})
}

// getHierarchy retrieves the full hierarchy of a task
func (s *APIServer) getHierarchy(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Task ID is required")
	}

	// Get hierarchy using business logic with workspace context
	iy := s.getIsYoneticiFromContext(c)
	hierarchy, err := iy.GorevHiyerarsiGetir(taskID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get hierarchy for task %s: %v", taskID, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    hierarchy,
	})
}

// addDependency adds a dependency between tasks
func (s *APIServer) addDependency(c *fiber.Ctx) error {
	hedefID := c.Params("id")
	if hedefID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Target task ID is required")
	}

	var req struct {
		KaynakID     string `json:"kaynak_id"`
		BaglantiTipi string `json:"baglanti_tipi"`
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.KaynakID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "kaynak_id is required")
	}
	if req.BaglantiTipi == "" {
		req.BaglantiTipi = "onceki" // default
	}

	// Add dependency using business logic with workspace context
	iy := s.getIsYoneticiFromContext(c)
	if _, err := iy.GorevBagimlilikEkle(req.KaynakID, hedefID, req.BaglantiTipi); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to add dependency from %s to %s: %v", req.KaynakID, hedefID, err))
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Dependency added successfully",
	})
}

// removeDependency removes a dependency between tasks
func (s *APIServer) removeDependency(c *fiber.Ctx) error {
	hedefID := c.Params("id")
	kaynakID := c.Params("dep_id")

	if hedefID == "" || kaynakID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Both task IDs are required")
	}

	// Remove dependency using VeriYonetici with workspace context
	iy := s.getIsYoneticiFromContext(c)
	err := iy.VeriYonetici().BaglantiSil(kaynakID, hedefID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to remove dependency: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Dependency removed successfully",
	})
}

// getActiveProject retrieves the currently active project
func (s *APIServer) getActiveProject(c *fiber.Ctx) error {
	iy := s.getIsYoneticiFromContext(c)
	aktifProjeID, err := iy.VeriYonetici().AktifProjeGetir()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to get active project: %v", err))
	}

	if aktifProjeID == "" {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    nil,
			"message": "No active project set",
		})
	}

	// Get full project details
	proje, err := iy.VeriYonetici().ProjeGetir(aktifProjeID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to get active project details: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    proje,
	})
}

// removeActiveProject removes the active project setting
func (s *APIServer) removeActiveProject(c *fiber.Ctx) error {
	iy := s.getIsYoneticiFromContext(c)
	if err := iy.VeriYonetici().AktifProjeKaldir(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to remove active project: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Active project removed successfully",
	})
}
