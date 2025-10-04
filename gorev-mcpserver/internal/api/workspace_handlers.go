package api

import (
	"github.com/gofiber/fiber/v2"
)

// registerWorkspaceHandler handles workspace registration requests
func (s *APIServer) registerWorkspaceHandler(c *fiber.Ctx) error {
	var req WorkspaceRegistration
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate request
	if req.Path == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Workspace path is required",
		})
	}

	// Register workspace
	workspace, err := s.workspaceManager.RegisterWorkspace(req.Path, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":      true,
		"workspace_id": workspace.ID,
		"workspace":    workspace.ToWorkspaceInfo(),
	})
}

// listWorkspacesHandler returns all registered workspaces
func (s *APIServer) listWorkspacesHandler(c *fiber.Ctx) error {
	workspaces := s.workspaceManager.ListWorkspaces()

	// Convert to WorkspaceInfo for API response
	infos := make([]*WorkspaceInfo, 0, len(workspaces))
	for _, ws := range workspaces {
		infos = append(infos, ws.ToWorkspaceInfo())
	}

	return c.JSON(fiber.Map{
		"success":    true,
		"workspaces": infos,
		"total":      len(infos),
	})
}

// getWorkspaceHandler returns workspace details by ID
func (s *APIServer) getWorkspaceHandler(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Workspace ID is required",
		})
	}

	workspace, err := s.workspaceManager.GetWorkspaceContext(workspaceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"workspace": workspace.ToWorkspaceInfo(),
	})
}

// unregisterWorkspaceHandler removes a workspace
func (s *APIServer) unregisterWorkspaceHandler(c *fiber.Ctx) error {
	workspaceID := c.Params("id")
	if workspaceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Workspace ID is required",
		})
	}

	if err := s.workspaceManager.UnregisterWorkspace(workspaceID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Workspace unregistered successfully",
	})
}
