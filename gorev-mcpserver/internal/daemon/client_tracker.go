package daemon

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// ClientInfo holds information about an active client connection
type ClientInfo struct {
	ClientID     string    `json:"client_id"`
	ClientType   string    `json:"client_type"` // "vscode", "mcp-proxy", "web-ui"
	WorkspaceID  string    `json:"workspace_id"`
	ConnectedAt  time.Time `json:"connected_at"`
	LastActivity time.Time `json:"last_activity"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// ClientTracker manages active client connections with thread-safe operations
type ClientTracker struct {
	mu       sync.RWMutex
	clients  map[string]*ClientInfo
	stopChan chan struct{}
}

// NewClientTracker creates a new client tracker with automatic cleanup
func NewClientTracker() *ClientTracker {
	ct := &ClientTracker{
		clients:  make(map[string]*ClientInfo),
		stopChan: make(chan struct{}),
	}
	// Start cleanup goroutine
	go ct.cleanupLoop()
	return ct
}

// RegisterClient adds a new client connection
func (ct *ClientTracker) RegisterClient(client *ClientInfo) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.clients[client.ClientID] = client
}

// UnregisterClient removes a client connection
func (ct *ClientTracker) UnregisterClient(clientID string) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	delete(ct.clients, clientID)
}

// GetActiveClientCount returns the number of active clients (non-expired)
func (ct *ClientTracker) GetActiveClientCount() int {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	count := 0
	now := time.Now()
	for _, client := range ct.clients {
		if now.Before(client.ExpiresAt) {
			count++
		}
	}
	return count
}

// GetClientCountByType returns the number of clients of a specific type
func (ct *ClientTracker) GetClientCountByType(clientType string) int {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	count := 0
	now := time.Now()
	for _, client := range ct.clients {
		if client.ClientType == clientType && now.Before(client.ExpiresAt) {
			count++
		}
	}
	return count
}

// GetClients returns a snapshot of all active clients
func (ct *ClientTracker) GetClients() []*ClientInfo {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	now := time.Now()
	clients := make([]*ClientInfo, 0)
	for _, client := range ct.clients {
		if now.Before(client.ExpiresAt) {
			clients = append(clients, client)
		}
	}
	return clients
}

// UpdateActivity updates the last activity time for a client
func (ct *ClientTracker) UpdateActivity(clientID string, duration time.Duration) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if client, exists := ct.clients[clientID]; exists {
		client.LastActivity = time.Now()
		client.ExpiresAt = client.LastActivity.Add(duration)
	}
}

// cleanupLoop periodically removes expired clients
func (ct *ClientTracker) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ct.cleanupExpired()
		case <-ct.stopChan:
			return
		}
	}
}

// cleanupExpired removes clients that have exceeded their TTL
func (ct *ClientTracker) cleanupExpired() {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	now := time.Now()
	for id, client := range ct.clients {
		if now.After(client.ExpiresAt) {
			delete(ct.clients, id)
		}
	}
}

// Shutdown stops the cleanup loop
func (ct *ClientTracker) Shutdown() {
	close(ct.stopChan)
}

// generateClientID generates a unique client ID
func GenerateClientID() string {
	return fmt.Sprintf("client-%d-%d", os.Getpid(), time.Now().UnixNano())
}
