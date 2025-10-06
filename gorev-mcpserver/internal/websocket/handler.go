package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// HandleWebSocket handles WebSocket connection requests
func HandleWebSocket(hub *Hub) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Extract workspace ID from query params or headers
		workspaceID := c.Query("workspace_id")
		if workspaceID == "" {
			workspaceID = c.Headers("X-Workspace-Id")
		}

		if workspaceID == "" {
			log.Printf("[WebSocket Handler] ❌ Missing workspace_id, closing connection")
			c.Close()
			return
		}

		// Create new client
		client := &Client{
			ID:          uuid.New().String(),
			WorkspaceID: workspaceID,
			Conn:        c,
			Send:        make(chan *ChangeEvent, 256),
		}

		// Register client
		hub.register <- client

		// Send welcome message
		welcomeEvent := &ChangeEvent{
			Type:        EventWorkspaceSync,
			WorkspaceID: workspaceID,
			Action:      "connected",
			Data: map[string]interface{}{
				"client_id": client.ID,
				"message":   "WebSocket connection established",
			},
		}
		client.Send <- welcomeEvent

		// Start read and write pumps
		go client.WritePump()
		client.ReadPump(hub)
	})
}

// WebSocketUpgradeMiddleware checks if the request can be upgraded to WebSocket
func WebSocketUpgradeMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
