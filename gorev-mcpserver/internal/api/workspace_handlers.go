package api

import (
	"crypto/sha256"
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/msenol/gorev/internal/config"
)

// generateWorkspaceIDFromPath generates a unique workspace ID from path
func generateWorkspaceIDFromPath(path string) string {
	hash := sha256.Sum256([]byte(path))
	return fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes of hash
}

// registerWorkspaceHandler handles workspace registration requests
func (s *APIServer) registerWorkspaceHandler(c *fiber.Ctx) error {
	var req WorkspaceRegistration
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	var workspace *WorkspaceContext
	var err error

	// Centralized mode: skip path validation, use workspace_id or generate from path
	if config.IsCentralizedMode() {
		workspaceID := req.WorkspaceID
		workspaceName := req.Name

		// If no workspace_id provided, generate from path or name
		if workspaceID == "" {
			if req.Path != "" {
				// Use path hash as workspace_id for consistency
				workspaceID = generateWorkspaceIDFromPath(req.Path)
				if workspaceName == "" {
					// Use last component of path as name
					workspaceName = filepath.Base(req.Path)
				}
			} else if req.Name != "" {
				// Use name as workspace_id
				workspaceID = req.Name
			} else {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"error":   "Workspace path, name, or workspace_id is required",
				})
			}
		}

		// Register by workspace ID (no local path validation in centralized mode)
		workspace, err = s.workspaceManager.RegisterWorkspaceByID(workspaceID, workspaceName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
	} else {
		// Local mode: require and validate path
		if req.Path == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Workspace path is required",
			})
		}

		// Register workspace with path (validates path exists)
		workspace, err = s.workspaceManager.RegisterWorkspace(req.Path, req.Name)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}
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
