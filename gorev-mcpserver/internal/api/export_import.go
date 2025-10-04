package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	mcpHandlers "github.com/msenol/gorev/internal/mcp"
)

// ExportRequest represents the export request payload
type ExportRequest struct {
	OutputPath           string   `json:"output_path"`
	Format               string   `json:"format"`
	IncludeCompleted     bool     `json:"include_completed"`
	IncludeDependencies  bool     `json:"include_dependencies"`
	IncludeTemplates     bool     `json:"include_templates"`
	IncludeAIContext     bool     `json:"include_ai_context"`
	ProjectFilter        []string `json:"project_filter"`
}

// ImportRequest represents the import request payload
type ImportRequest struct {
	FilePath             string            `json:"file_path"`
	ImportMode           string            `json:"import_mode"`
	ConflictResolution   string            `json:"conflict_resolution"`
	DryRun               bool              `json:"dry_run"`
	PreserveIDs          bool              `json:"preserve_ids"`
	ProjectMapping       map[string]string `json:"project_mapping"`
}

// exportData handles data export requests
func (s *APIServer) exportData(c *fiber.Ctx) error {
	var req ExportRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
	}

	// Validate required fields
	if req.OutputPath == "" {
		return fiber.NewError(fiber.StatusBadRequest, "output_path is required")
	}

	// Default format to JSON
	if req.Format == "" {
		req.Format = "json"
	}

	// Build params for MCP handler
	params := map[string]interface{}{
		"output_path":           req.OutputPath,
		"format":                req.Format,
		"include_completed":     req.IncludeCompleted,
		"include_dependencies":  req.IncludeDependencies,
		"include_templates":     req.IncludeTemplates,
		"include_ai_context":    req.IncludeAIContext,
	}

	if len(req.ProjectFilter) > 0 {
		// Convert []string to []interface{} for MCP handler
		projectFilterInterface := make([]interface{}, len(req.ProjectFilter))
		for i, p := range req.ProjectFilter {
			projectFilterInterface[i] = p
		}
		params["project_filter"] = projectFilterInterface
	}

	// Call MCP handler through server's handlers field
	if s.handlers == nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "export service not available")
	}

	// Type assert to MCP Handlers
	handlers, ok := s.handlers.(*mcpHandlers.Handlers)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid handler type")
	}

	result, err := handlers.GorevExport(params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("export failed: %v", err))
	}

	// Extract message from result
	message := "Export completed successfully"
	if result != nil && len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(map[string]interface{}); ok {
			if text, ok := textContent["text"].(string); ok {
				message = text
			}
		}
	}

	if result != nil && result.IsError {
		return fiber.NewError(fiber.StatusInternalServerError, message)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
		"path":    req.OutputPath,
	})
}

// importData handles data import requests
func (s *APIServer) importData(c *fiber.Ctx) error {
	var req ImportRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
	}

	// Validate required fields
	if req.FilePath == "" {
		return fiber.NewError(fiber.StatusBadRequest, "file_path is required")
	}

	// Default values
	if req.ImportMode == "" {
		req.ImportMode = "merge"
	}
	if req.ConflictResolution == "" {
		req.ConflictResolution = "skip"
	}

	// Build params for MCP handler
	params := map[string]interface{}{
		"file_path":           req.FilePath,
		"import_mode":         req.ImportMode,
		"conflict_resolution": req.ConflictResolution,
		"dry_run":             req.DryRun,
		"preserve_ids":        req.PreserveIDs,
	}

	if len(req.ProjectMapping) > 0 {
		params["project_mapping"] = req.ProjectMapping
	}

	// Call MCP handler through server's handlers field
	if s.handlers == nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, "import service not available")
	}

	// Type assert to MCP Handlers
	handlers, ok := s.handlers.(*mcpHandlers.Handlers)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid handler type")
	}

	result, err := handlers.GorevImport(params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("import failed: %v", err))
	}

	// Extract message from result
	message := "Import completed successfully"
	if result != nil && len(result.Content) > 0 {
		if textContent, ok := result.Content[0].(map[string]interface{}); ok {
			if text, ok := textContent["text"].(string); ok {
				message = text
			}
		}
	}

	if result != nil && result.IsError {
		return fiber.NewError(fiber.StatusInternalServerError, message)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": message,
	})
}