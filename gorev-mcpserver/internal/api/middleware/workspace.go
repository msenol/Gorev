package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// WorkspaceGetter is a minimal interface for getting workspaces
// This is satisfied by *api.WorkspaceManager
type WorkspaceGetter interface {
	GetWorkspace(workspaceID string) (any, error)
}

// WorkspaceMiddleware creates a middleware that extracts workspace context from request headers
// and attaches it to the Fiber context for downstream handlers to use
func WorkspaceMiddleware(getter WorkspaceGetter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract workspace identifiers from headers
		workspaceID := c.Get("X-Workspace-Id")
		workspacePath := c.Get("X-Workspace-Path")
		workspaceName := c.Get("X-Workspace-Name")

		// Log ALL requests for debugging (including empty headers)
		log.Printf("[Workspace Middleware] %s %s | Headers: ID=%q, Path=%q, Name=%q",
			c.Method(), c.Path(), workspaceID, workspacePath, workspaceName)

		// Store workspace identifiers in context
		c.Locals("workspace_id", workspaceID)
		c.Locals("workspace_path", workspacePath)
		c.Locals("workspace_name", workspaceName)

		// If workspace ID is provided, try to load the workspace context
		if workspaceID != "" {
			workspace, err := getter.GetWorkspace(workspaceID)
			if err != nil {
				log.Printf("[Workspace Middleware] ❌ Failed to get workspace %s: %v", workspaceID, err)
				// Don't fail the request, just log the error
				// This allows workspace registration requests to proceed
			} else {
				log.Printf("[Workspace Middleware] ✓ Found workspace %s, storing in context", workspaceID)
				// Store workspace context in Fiber locals for handlers to access
				c.Locals("workspace", workspace)

				// Extract IsYonetici using interface with 'any' return type
				// This matches WorkspaceContext.GetIsYonetici() any method
				type isYoneticiGetter interface {
					GetIsYonetici() any
				}

				if wg, ok := workspace.(isYoneticiGetter); ok {
					isYonetici := wg.GetIsYonetici()
					log.Printf("[Workspace Middleware] ✓ Extracted IsYonetici, type: %T", isYonetici)
					c.Locals("is_yonetici", isYonetici)
				} else {
					log.Printf("[Workspace Middleware] ❌ Workspace doesn't implement GetIsYonetici() any, type: %T", workspace)
				}
			}
		}

		return c.Next()
	}
}

// GetWorkspaceID extracts the workspace ID from Fiber context
func GetWorkspaceID(c *fiber.Ctx) string {
	if id, ok := c.Locals("workspace_id").(string); ok {
		return id
	}
	return ""
}

// GetWorkspacePath extracts the workspace path from Fiber context
func GetWorkspacePath(c *fiber.Ctx) string {
	if path, ok := c.Locals("workspace_path").(string); ok {
		return path
	}
	return ""
}

// GetWorkspaceName extracts the workspace name from Fiber context
func GetWorkspaceName(c *fiber.Ctx) string {
	if name, ok := c.Locals("workspace_name").(string); ok {
		return name
	}
	return ""
}

// GetIsYonetici extracts the IsYonetici from Fiber context
// Returns nil if workspace context is not available
func GetIsYonetici(c *fiber.Ctx) any {
	return c.Locals("is_yonetici")
}
