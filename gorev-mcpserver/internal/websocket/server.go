package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

// EventType represents different types of database change events
type EventType string

const (
	EventTaskCreated   EventType = "task_created"
	EventTaskUpdated   EventType = "task_updated"
	EventTaskDeleted   EventType = "task_deleted"
	EventProjectCreated EventType = "project_created"
	EventProjectUpdated EventType = "project_updated"
	EventProjectDeleted EventType = "project_deleted"
	EventTemplateChanged EventType = "template_changed"
	EventWorkspaceSync  EventType = "workspace_sync"
)

// ChangeEvent represents a database change event
type ChangeEvent struct {
	Type        EventType              `json:"type"`
	WorkspaceID string                 `json:"workspace_id"`
	EntityID    string                 `json:"entity_id,omitempty"`
	EntityType  string                 `json:"entity_type,omitempty"`
	Action      string                 `json:"action,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	Timestamp   int64                  `json:"timestamp"`
}

// Client represents a connected WebSocket client
type Client struct {
	ID          string
	WorkspaceID string
	Conn        *websocket.Conn
	Send        chan *ChangeEvent
}

// Hub manages WebSocket connections and broadcasting
type Hub struct {
	// Registered clients by workspace ID
	clients map[string]map[*Client]bool

	// Inbound events from database changes
	broadcast chan *ChangeEvent

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe access
	mu sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan *ChangeEvent, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.WorkspaceID]; !ok {
				h.clients[client.WorkspaceID] = make(map[*Client]bool)
			}
			h.clients[client.WorkspaceID][client] = true
			h.mu.Unlock()
			log.Printf("[WebSocket Hub] ✓ Client registered: %s (workspace: %s)", client.ID, client.WorkspaceID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.WorkspaceID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.clients, client.WorkspaceID)
					}
					log.Printf("[WebSocket Hub] ✗ Client unregistered: %s (workspace: %s)", client.ID, client.WorkspaceID)
				}
			}
			h.mu.Unlock()

		case event := <-h.broadcast:
			h.mu.RLock()
			if clients, ok := h.clients[event.WorkspaceID]; ok {
				for client := range clients {
					select {
					case client.Send <- event:
						// Event sent successfully
					default:
						// Client's send buffer is full, close connection
						close(client.Send)
						delete(clients, client)
						log.Printf("[WebSocket Hub] ⚠ Client send buffer full, disconnecting: %s", client.ID)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastEvent sends an event to all clients in a workspace
func (h *Hub) BroadcastEvent(event *ChangeEvent) {
	h.broadcast <- event
}

// GetClientCount returns the total number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	count := 0
	for _, clients := range h.clients {
		count += len(clients)
	}
	return count
}

// GetWorkspaceClientCount returns the number of clients for a specific workspace
func (h *Hub) GetWorkspaceClientCount(workspaceID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.clients[workspaceID]; ok {
		return len(clients)
	}
	return 0
}

// ReadPump reads messages from the WebSocket connection
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg map[string]interface{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket Client] ⚠ Read error: %v", err)
			}
			break
		}

		// Handle ping/pong or other client messages
		if msgType, ok := msg["type"].(string); ok {
			if msgType == "ping" {
				// Send pong
				pong := map[string]interface{}{
					"type": "pong",
					"timestamp": msg["timestamp"],
				}
				c.Conn.WriteJSON(pong)
			}
		}
	}
}

// WritePump writes messages to the WebSocket connection
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for event := range c.Send {
		data, err := json.Marshal(event)
		if err != nil {
			log.Printf("[WebSocket Client] ⚠ Marshal error: %v", err)
			continue
		}

		if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("[WebSocket Client] ⚠ Write error: %v", err)
			return
		}
	}
}
