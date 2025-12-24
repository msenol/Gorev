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
	"github.com/msenol/gorev/internal/daemon"
	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	ws "github.com/msenol/gorev/internal/websocket"
)

// APIServer represents the HTTP API server
type APIServer struct {
	app              *fiber.App
	port             string
	isYonetici       *gorev.IsYonetici     // Legacy single workspace manager (deprecated)
	workspaceManager *WorkspaceManager     // Multi-workspace manager
	handlers         interface{}           // MCP Handlers for export/import operations
	wsHub            *ws.Hub               // WebSocket hub for real-time updates
	clientTracker    *daemon.ClientTracker // Active client tracking for smart shutdown
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
		clientTracker:    daemon.NewClientTracker(),
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
	api.Post("/template/init", s.initializeTemplates)

	// Summary routes
	api.Get("/summary", s.getSummary)

	// Language routes
	api.Get("/language", s.getLanguage)
	api.Post("/language", s.setLanguage)

	// Subtask routes
	api.Get("/tasks/:id/subtasks", s.getSubtasks)
	api.Post("/tasks/:id/subtasks", s.createSubtask)
	api.Put("/tasks/:id/parent", s.changeParent)
	api.Get("/tasks/:id/hierarchy", s.getHierarchy)

	// Dependency routes
	api.Get("/tasks/:id/dependencies", s.getDependencies)
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

	// Client management routes (for smart shutdown)
	api.Get("/daemon/clients/count", s.getActiveClientCountHandler)
	api.Post("/daemon/clients/register", s.registerClientHandler)
	api.Post("/daemon/clients/unregister", s.unregisterClientHandler)
	api.Post("/daemon/heartbeat", s.heartbeatHandler)

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
	ctx := s.getContextFromRequest(c)
	gorevler, err := iy.GorevListele(ctx, filters)
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
	ctx := s.getContextFromRequest(c)
	task, err := iy.VeriYonetici().GorevGetir(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("failed to get task with ID %s: %v", id, err))
	}

	// Load subtasks for the task
	subtasks, err := iy.AltGorevleriGetir(ctx, id)
	if err == nil && len(subtasks) > 0 {
		task.Subtasks = subtasks
	}

	// Load dependencies for the task (for VS Code extension task detail panel)
	baglantilari, err := iy.VeriYonetici().BaglantilariGetir(ctx, id)
	if err == nil && len(baglantilari) > 0 {
		bagimliliklar := make([]gorev.Bagimlilik, 0, len(baglantilari))
		for _, b := range baglantilari {
			dep := gorev.Bagimlilik{
				KaynakID:    b.SourceID,
				HedefID:     b.TargetID,
				BaglantiTip: b.ConnectionType,
			}
			// Get target task info for display
			if targetTask, err := iy.VeriYonetici().GorevGetir(ctx, b.TargetID); err == nil {
				dep.HedefBaslik = targetTask.Title
				dep.HedefDurum = targetTask.Status
			}
			bagimliliklar = append(bagimliliklar, dep)
		}
		task.Bagimliliklar = bagimliliklar
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    task,
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

	// Handle task updates - use GorevGuncelle with params map and workspace context
	// Note: API accepts Turkish field names but DB columns are English since v0.17.0
	// Mapping: durum -> status, proje_id -> project_id, oncelik -> priority
	iy := s.getIsYoneticiFromContext(c)
	ctx := s.getContextFromRequest(c)

	// Build update params with proper field name mapping
	params := make(map[string]interface{})

	if durum, ok := req["durum"].(string); ok {
		params["status"] = durum
	}
	if projeID, ok := req["proje_id"].(string); ok {
		params["project_id"] = projeID
	}
	if oncelik, ok := req["oncelik"].(string); ok {
		params["priority"] = oncelik
	}
	if baslik, ok := req["baslik"].(string); ok {
		params["title"] = baslik
	}
	if aciklama, ok := req["aciklama"].(string); ok {
		params["description"] = aciklama
	}

	if len(params) > 0 {
		if err := iy.VeriYonetici().GorevGuncelle(ctx, id, params); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update task %s: %v", id, err))
		}
	}

	// Get updated task
	gorev, err := iy.VeriYonetici().GorevGetir(ctx, id)
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
	ctx := s.getContextFromRequest(c)
	if err := iy.VeriYonetici().GorevSil(ctx, id); err != nil {
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
	ctx := s.getContextFromRequest(c)
	gorev, err := iy.TemplatedenGorevOlustur(ctx, req.TemplateID, req.Degerler)
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
		return s.isYonetici
	}

	iy, ok := isYonetici.(*gorev.IsYonetici)
	if !ok {
		return s.isYonetici
	}

	return iy
}

// getContextFromRequest extracts language from request and creates language-aware context
func (s *APIServer) getContextFromRequest(c *fiber.Ctx) context.Context {
	lang := c.Get("Accept-Language", "tr")
	// Simple language code extraction (take first 2 chars)
	if len(lang) >= 2 {
		lang = lang[:2]
	}
	if lang != "tr" && lang != "en" {
		lang = "tr" // default fallback
	}
	return i18n.WithLanguage(c.UserContext(), lang)
}

// getProjects retrieves all projects
func (s *APIServer) getProjects(c *fiber.Ctx) error {
	iy := s.getIsYoneticiFromContext(c)
	ctx := s.getContextFromRequest(c)

	projeler, err := iy.ProjeListele(ctx)
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
	ctx := s.getContextFromRequest(c)
	proje, err := iy.VeriYonetici().ProjeGetir(ctx, id)
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
	ctx := s.getContextFromRequest(c)

	// Create project using business logic
	proje, err := iy.ProjeOlustur(ctx, req.Name, req.Definition)
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
	ctx := s.getContextFromRequest(c)
	gorevler, err := iy.GorevListele(ctx, filters)
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
	ctx := s.getContextFromRequest(c)
	if err := iy.VeriYonetici().AktifProjeAyarla(ctx, projeID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to activate project with ID %s: %v", projeID, err))
	}

	// Return updated project
	proje, err := iy.VeriYonetici().ProjeGetir(ctx, projeID)
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
	ctx := s.getContextFromRequest(c)
	templateler, err := iy.VeriYonetici().TemplateListele(ctx, kategori)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to list templates with category '%s': %v", kategori, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    templateler,
		"total":   len(templateler),
	})
}

// initializeTemplates creates default templates (TR/EN pairs) in the database
func (s *APIServer) initializeTemplates(c *fiber.Ctx) error {
	iy := s.getIsYoneticiFromContext(c)
	ctx := s.getContextFromRequest(c)

	// Call the business logic to initialize templates
	if err := iy.VeriYonetici().VarsayilanTemplateleriOlustur(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to initialize templates: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Default templates initialized successfully",
	})
}

// getSummary retrieves system-wide summary statistics
func (s *APIServer) getSummary(c *fiber.Ctx) error {
	iy := s.getIsYoneticiFromContext(c)
	ctx := s.getContextFromRequest(c)

	// Get all tasks
	allTasks, err := iy.VeriYonetici().GorevleriGetir(ctx, "", "", "")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get tasks: %v", err))
	}

	// Get all projects
	projects, err := iy.VeriYonetici().ProjeleriGetir(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get projects: %v", err))
	}

	// Get templates
	templates, err := iy.VeriYonetici().TemplateListele(ctx, "")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get templates: %v", err))
	}

	// Get active project
	activeProjectID, _ := iy.VeriYonetici().AktifProjeGetir(ctx)
	var activeProject *gorev.Proje
	if activeProjectID != "" {
		activeProject, _ = iy.VeriYonetici().ProjeGetir(ctx, activeProjectID)
	}

	// Calculate statistics
	statusCounts := map[string]int{
		"pending":     0,
		"in_progress": 0,
		"completed":   0,
	}
	priorityCounts := map[string]int{
		"high":   0,
		"medium": 0,
		"low":    0,
	}

	var overdueTasks []*gorev.Gorev
	var dueTodayTasks []*gorev.Gorev
	var dueThisWeekTasks []*gorev.Gorev
	var highPriorityTasks []*gorev.Gorev
	var blockedTasks []*gorev.Gorev
	var recentTasks []*gorev.Gorev

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekEnd := today.AddDate(0, 0, 7)

	for _, task := range allTasks {
		// Status counts
		switch task.Status {
		case "pending", "beklemede":
			statusCounts["pending"]++
		case "in_progress", "devam_ediyor":
			statusCounts["in_progress"]++
		case "completed", "tamamlandi":
			statusCounts["completed"]++
		}

		// Priority counts
		switch task.Priority {
		case "high", "yuksek":
			priorityCounts["high"]++
			if task.Status != "completed" && task.Status != "tamamlandi" {
				highPriorityTasks = append(highPriorityTasks, task)
			}
		case "medium", "orta":
			priorityCounts["medium"]++
		case "low", "dusuk":
			priorityCounts["low"]++
		}

		// Due date analysis
		if task.DueDate != nil && task.Status != "completed" && task.Status != "tamamlandi" {
			dueDate := *task.DueDate
			if dueDate.Before(today) {
				overdueTasks = append(overdueTasks, task)
			} else if dueDate.Before(today.AddDate(0, 0, 1)) {
				dueTodayTasks = append(dueTodayTasks, task)
			} else if dueDate.Before(weekEnd) {
				dueThisWeekTasks = append(dueThisWeekTasks, task)
			}
		}

		// Blocked tasks (has uncompleted dependencies)
		if task.UncompletedDependencyCount > 0 {
			blockedTasks = append(blockedTasks, task)
		}

		// Recent tasks (last 5 updated)
		if len(recentTasks) < 5 {
			recentTasks = append(recentTasks, task)
		}
	}

	// Limit lists to reasonable size
	if len(highPriorityTasks) > 5 {
		highPriorityTasks = highPriorityTasks[:5]
	}
	if len(overdueTasks) > 5 {
		overdueTasks = overdueTasks[:5]
	}
	if len(blockedTasks) > 5 {
		blockedTasks = blockedTasks[:5]
	}

	// Calculate completion rate
	totalTasks := len(allTasks)
	completedTasks := statusCounts["completed"]
	completionRate := 0.0
	if totalTasks > 0 {
		completionRate = float64(completedTasks) / float64(totalTasks) * 100
	}

	// Build response
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"total_tasks":     totalTasks,
			"total_projects":  len(projects),
			"total_templates": len(templates),
			"active_project":  activeProject,
			"status_counts": fiber.Map{
				"pending":     statusCounts["pending"],
				"in_progress": statusCounts["in_progress"],
				"completed":   statusCounts["completed"],
			},
			"priority_counts": fiber.Map{
				"high":   priorityCounts["high"],
				"medium": priorityCounts["medium"],
				"low":    priorityCounts["low"],
			},
			"due_date_summary": fiber.Map{
				"overdue":       len(overdueTasks),
				"due_today":     len(dueTodayTasks),
				"due_this_week": len(dueThisWeekTasks),
			},
			"completion_rate":     completionRate,
			"high_priority_tasks": highPriorityTasks,
			"overdue_tasks":       overdueTasks,
			"blocked_tasks":       blockedTasks,
			"recent_tasks":        recentTasks,
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

// getSubtasks retrieves all subtasks for a parent task
func (s *APIServer) getSubtasks(c *fiber.Ctx) error {
	parentID := c.Params("id")
	if parentID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Parent task ID is required")
	}

	iy := s.getIsYoneticiFromContext(c)
	ctx := s.getContextFromRequest(c)

	subtasks, err := iy.AltGorevleriGetir(ctx, parentID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get subtasks: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    subtasks,
		"total":   len(subtasks),
	})
}

// createSubtask creates a subtask under a parent task
func (s *APIServer) createSubtask(c *fiber.Ctx) error {
	parentID := c.Params("id")
	if parentID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Parent task ID is required")
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		DueDate     string `json:"due_date"`
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
	ctx := s.getContextFromRequest(c)
	gorev, err := iy.AltGorevOlustur(ctx, parentID, req.Title, req.Description, req.Priority, req.DueDate, etiketler)
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
	ctx := s.getContextFromRequest(c)
	if err := iy.GorevUstDegistir(ctx, taskID, req.NewParentID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to change parent for task %s: %v", taskID, err))
	}

	// Get updated task
	gorev, err := iy.VeriYonetici().GorevGetir(ctx, taskID)
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
	ctx := s.getContextFromRequest(c)
	hierarchy, err := iy.GorevHiyerarsiGetir(ctx, taskID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get hierarchy for task %s: %v", taskID, err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    hierarchy,
	})
}

// getDependencies retrieves all dependencies for a task
func (s *APIServer) getDependencies(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Task ID is required")
	}

	iy := s.getIsYoneticiFromContext(c)
	ctx := s.getContextFromRequest(c)

	// Get all connections for this task
	baglantilari, err := iy.VeriYonetici().BaglantilariGetir(ctx, taskID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get dependencies for task %s: %v", taskID, err))
	}

	// Build dependency list with task details
	type DependencyInfo struct {
		SourceID       string `json:"source_id"`
		TargetID       string `json:"target_id"`
		ConnectionType string `json:"connection_type"`
		SourceTitle    string `json:"source_title,omitempty"`
		SourceStatus   string `json:"source_status,omitempty"`
		TargetTitle    string `json:"target_title,omitempty"`
		TargetStatus   string `json:"target_status,omitempty"`
	}

	dependencies := make([]DependencyInfo, 0)
	for _, b := range baglantilari {
		dep := DependencyInfo{
			SourceID:       b.SourceID,
			TargetID:       b.TargetID,
			ConnectionType: b.ConnectionType,
		}

		// Get source task info
		if sourceTask, err := iy.VeriYonetici().GorevGetir(ctx, b.SourceID); err == nil {
			dep.SourceTitle = sourceTask.Title
			dep.SourceStatus = sourceTask.Status
		}

		// Get target task info
		if targetTask, err := iy.VeriYonetici().GorevGetir(ctx, b.TargetID); err == nil {
			dep.TargetTitle = targetTask.Title
			dep.TargetStatus = targetTask.Status
		}

		dependencies = append(dependencies, dep)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    dependencies,
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
	ctx := s.getContextFromRequest(c)
	if _, err := iy.GorevBagimlilikEkle(ctx, req.KaynakID, hedefID, req.BaglantiTipi); err != nil {
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
	ctx := s.getContextFromRequest(c)
	err := iy.VeriYonetici().BaglantiSil(ctx, kaynakID, hedefID)
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
	ctx := s.getContextFromRequest(c)
	aktifProjeID, err := iy.VeriYonetici().AktifProjeGetir(ctx)
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
	proje, err := iy.VeriYonetici().ProjeGetir(ctx, aktifProjeID)
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
	ctx := s.getContextFromRequest(c)
	if err := iy.VeriYonetici().AktifProjeKaldir(ctx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to remove active project: %v", err))
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Active project removed successfully",
	})
}

// getActiveClientCountHandler returns the number of active clients
func (s *APIServer) getActiveClientCountHandler(c *fiber.Ctx) error {
	count := s.clientTracker.GetActiveClientCount()
	return c.JSON(fiber.Map{
		"success":      true,
		"client_count": count,
	})
}

// registerClientHandler registers a new client connection
func (s *APIServer) registerClientHandler(c *fiber.Ctx) error {
	var req struct {
		ClientID    string `json:"client_id"`
		ClientType  string `json:"client_type"`
		WorkspaceID string `json:"workspace_id"`
		TTL         int    `json:"ttl_seconds"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.ClientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "client_id is required")
	}

	if req.ClientType == "" {
		req.ClientType = "unknown"
	}

	if req.TTL <= 0 {
		req.TTL = 300 // Default 5 minutes
	}

	client := &daemon.ClientInfo{
		ClientID:     req.ClientID,
		ClientType:   req.ClientType,
		WorkspaceID:  req.WorkspaceID,
		ConnectedAt:  time.Now(),
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(time.Duration(req.TTL) * time.Second),
	}

	s.clientTracker.RegisterClient(client)

	log.Printf("[ClientTracker] Client registered: %s (%s)", req.ClientID, req.ClientType)

	return c.JSON(fiber.Map{
		"success":    true,
		"client_id":  req.ClientID,
		"expires_at": client.ExpiresAt.Unix(),
	})
}

// unregisterClientHandler removes a client connection
func (s *APIServer) unregisterClientHandler(c *fiber.Ctx) error {
	var req struct {
		ClientID string `json:"client_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.ClientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "client_id is required")
	}

	s.clientTracker.UnregisterClient(req.ClientID)

	log.Printf("[ClientTracker] Client unregistered: %s", req.ClientID)

	return c.JSON(fiber.Map{
		"success":   true,
		"client_id": req.ClientID,
	})
}

// heartbeatHandler updates client activity to extend TTL
func (s *APIServer) heartbeatHandler(c *fiber.Ctx) error {
	var req struct {
		ClientID string `json:"client_id"`
		TTL      int    `json:"ttl_seconds"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if req.ClientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "client_id is required")
	}

	if req.TTL <= 0 {
		req.TTL = 300
	}

	s.clientTracker.UpdateActivity(req.ClientID, time.Duration(req.TTL)*time.Second)

	return c.JSON(fiber.Map{
		"success":   true,
		"client_id": req.ClientID,
	})
}
