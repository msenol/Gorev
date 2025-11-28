package api

import (
	"time"

	"github.com/msenol/gorev/internal/gorev"
	ws "github.com/msenol/gorev/internal/websocket"
)

// WorkspaceContextInterface is satisfied by *WorkspaceContext
// This allows the middleware package to work with workspace contexts without circular dependencies
type WorkspaceContextInterface interface {
	GetIsYonetici() any // Returns *gorev.IsYonetici but typed as 'any' for middleware compatibility
	ToWorkspaceInfo() *WorkspaceInfo
}

// WorkspaceContext represents a registered workspace with its database connection
type WorkspaceContext struct {
	ID           string              `json:"id"`            // Workspace unique identifier (typically absolute path)
	Name         string              `json:"name"`          // Workspace display name
	Path         string              `json:"path"`          // Absolute filesystem path
	DatabasePath string              `json:"database_path"` // Path to workspace database
	VeriYonetici *gorev.VeriYonetici `json:"-"`             // Database manager (not serialized)
	IsYonetici   *gorev.IsYonetici   `json:"-"`             // Business logic manager (not serialized)
	EventEmitter ws.EventEmitter     `json:"-"`             // Event emitter for real-time updates (not serialized)
	LastAccessed time.Time           `json:"last_accessed"` // Last time this workspace was accessed
	CreatedAt    time.Time           `json:"created_at"`    // When workspace was registered
	TaskCount    int                 `json:"task_count"`    // Cached task count for UI
}

// WorkspaceRegistration represents a workspace registration request
type WorkspaceRegistration struct {
	Path        string `json:"path"`         // Absolute path to workspace folder (required in local mode)
	Name        string `json:"name"`         // Optional display name (defaults to folder name)
	WorkspaceID string `json:"workspace_id"` // Explicit workspace ID (for centralized mode)
}

// WorkspaceInfo is a lightweight representation of a workspace for API responses
type WorkspaceInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	TaskCount    int       `json:"task_count"`
	LastAccessed time.Time `json:"last_accessed"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToWorkspaceInfo converts WorkspaceContext to WorkspaceInfo for API responses
func (wc *WorkspaceContext) ToWorkspaceInfo() *WorkspaceInfo {
	return &WorkspaceInfo{
		ID:           wc.ID,
		Name:         wc.Name,
		Path:         wc.Path,
		TaskCount:    wc.TaskCount,
		LastAccessed: wc.LastAccessed,
		CreatedAt:    wc.CreatedAt,
	}
}

// GetIsYonetici returns the business logic manager for this workspace
// Returns as 'any' for interface compatibility with middleware
func (wc *WorkspaceContext) GetIsYonetici() any {
	return wc.IsYonetici
}
