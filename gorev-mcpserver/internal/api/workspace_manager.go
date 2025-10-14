package api

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	ws "github.com/msenol/gorev/internal/websocket"
)

// WorkspaceManager manages multiple workspace contexts with their database connections
type WorkspaceManager struct {
	workspaces   map[string]*WorkspaceContext // Keyed by workspace ID
	migrationsFS fs.FS                        // Embedded migrations filesystem (optional)
	wsHub        *ws.Hub                      // WebSocket hub for real-time updates
	mu           sync.RWMutex
}

// NewWorkspaceManager creates a new workspace manager
func NewWorkspaceManager() *WorkspaceManager {
	return &WorkspaceManager{
		workspaces: make(map[string]*WorkspaceContext),
	}
}

// SetMigrationsFS sets the embedded migrations filesystem for database initialization
func (wm *WorkspaceManager) SetMigrationsFS(migrationsFS fs.FS) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.migrationsFS = migrationsFS
}

// RegisterWorkspace registers a new workspace or returns existing one
func (wm *WorkspaceManager) RegisterWorkspace(path string, name string) (*WorkspaceContext, error) {
	// Clean and validate path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.invalidWorkspacePath", map[string]interface{}{"Error": err}))
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf(i18n.T("error.workspacePathNotExist", map[string]interface{}{"Path": absPath}))
	}

	// Generate workspace ID from path
	workspaceID := generateWorkspaceID(absPath)

	wm.mu.Lock()
	defer wm.mu.Unlock()

	// Check if workspace already registered
	if ws, exists := wm.workspaces[workspaceID]; exists {
		// Update last accessed time
		ws.LastAccessed = time.Now()
		return ws, nil
	}

	// Set default name if not provided
	if name == "" {
		name = filepath.Base(absPath)
	}

	// Determine database path (prefer workspace-local .gorev folder)
	dbPath := filepath.Join(absPath, ".gorev", "gorev.db")

	// Ensure .gorev directory exists
	gorevDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(gorevDir, 0755); err != nil {
		return nil, fmt.Errorf(i18n.T("error.gorevDirCreateFailed", map[string]interface{}{"Error": err}))
	}

	// Create event emitter for real-time updates BEFORE initializing database
	var eventEmitter ws.EventEmitter
	if wm.wsHub != nil {
		eventEmitter = ws.NewHubEventEmitter(wm.wsHub)
	} else {
		// Fallback to no-op emitter if hub not available
		eventEmitter = ws.NewNoOpEventEmitter()
	}

	// Initialize database manager with event emitter
	var veriYonetici *gorev.VeriYonetici
	// Already holding write lock, no need for read lock
	migrationsFS := wm.migrationsFS

	if migrationsFS != nil {
		// Use embedded migrations if available with event emitter
		veriYonetici, err = gorev.YeniVeriYoneticiWithEmbeddedMigrationsAndEventEmitter(dbPath, migrationsFS, eventEmitter, workspaceID)
	} else {
		// Fallback to filesystem migrations (for tests and backwards compatibility)
		// Search for migrations in standard locations
		migrationsPath := findMigrationsPath()
		if migrationsPath == "" {
			// Use embedded path as last resort (will use embedded FS if available in binary)
			migrationsPath = "embedded://migrations"
		}
		veriYonetici, err = gorev.YeniVeriYoneticiWithEventEmitter(dbPath, migrationsPath, eventEmitter, workspaceID)
	}

	if err != nil {
		return nil, fmt.Errorf(i18n.T("error.workspaceDbInitFailed", map[string]interface{}{"Name": name, "Error": err}))
	}

	// Initialize business logic manager
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Get task count
	taskCount, _ := wm.getTaskCount(isYonetici)

	// Create workspace context (eventEmitter already created above)
	workspace := &WorkspaceContext{
		ID:           workspaceID,
		Name:         name,
		Path:         absPath,
		DatabasePath: dbPath,
		VeriYonetici: veriYonetici,
		IsYonetici:   isYonetici,
		EventEmitter: eventEmitter,
		LastAccessed: time.Now(),
		CreatedAt:    time.Now(),
		TaskCount:    taskCount,
	}

	wm.workspaces[workspaceID] = workspace
	return workspace, nil
}

// GetWorkspace retrieves a workspace by ID
// Returns any to satisfy middleware.WorkspaceGetter interface
func (wm *WorkspaceManager) GetWorkspace(workspaceID string) (any, error) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return nil, fmt.Errorf(i18n.T("error.workspaceNotFound", map[string]interface{}{"ID": workspaceID}))
	}

	// Update last accessed time
	workspace.LastAccessed = time.Now()
	return workspace, nil
}

// GetWorkspaceContext retrieves a workspace by ID and returns the concrete type
// Use this when you need the full WorkspaceContext struct
func (wm *WorkspaceManager) GetWorkspaceContext(workspaceID string) (*WorkspaceContext, error) {
	workspace, err := wm.GetWorkspace(workspaceID)
	if err != nil {
		return nil, err
	}
	if ws, ok := workspace.(*WorkspaceContext); ok {
		return ws, nil
	}
	return nil, fmt.Errorf(i18n.T("error.unexpectedWorkspaceType"))
}

// ListWorkspaces returns all registered workspaces
func (wm *WorkspaceManager) ListWorkspaces() []*WorkspaceContext {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	workspaces := make([]*WorkspaceContext, 0, len(wm.workspaces))
	for _, ws := range wm.workspaces {
		workspaces = append(workspaces, ws)
	}

	return workspaces
}

// UnregisterWorkspace removes a workspace and closes its database connection
func (wm *WorkspaceManager) UnregisterWorkspace(workspaceID string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	workspace, exists := wm.workspaces[workspaceID]
	if !exists {
		return fmt.Errorf(i18n.T("error.workspaceNotFound", map[string]interface{}{"ID": workspaceID}))
	}

	// Close database connection
	if workspace.VeriYonetici != nil {
		if err := workspace.VeriYonetici.Kapat(); err != nil {
			return fmt.Errorf(i18n.T("error.dbCloseFailed", map[string]interface{}{"Error": err}))
		}
	}

	delete(wm.workspaces, workspaceID)
	return nil
}

// UpdateTaskCount updates the cached task count for a workspace
func (wm *WorkspaceManager) UpdateTaskCount(workspaceID string) error {
	workspace, err := wm.GetWorkspaceContext(workspaceID)
	if err != nil {
		return err
	}

	taskCount, err := wm.getTaskCount(workspace.IsYonetici)
	if err != nil {
		return err
	}

	wm.mu.Lock()
	workspace.TaskCount = taskCount
	wm.mu.Unlock()

	return nil
}

// Cleanup removes workspaces that haven't been accessed recently
func (wm *WorkspaceManager) Cleanup(maxAge time.Duration) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	now := time.Now()
	for id, ws := range wm.workspaces {
		if now.Sub(ws.LastAccessed) > maxAge {
			// Close database connection
			if ws.VeriYonetici != nil {
				_ = ws.VeriYonetici.Kapat()
			}
			delete(wm.workspaces, id)
		}
	}

	return nil
}

// CloseAll closes all workspace database connections
func (wm *WorkspaceManager) CloseAll() error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	var errors []error
	for id, ws := range wm.workspaces {
		if ws.VeriYonetici != nil {
			if err := ws.VeriYonetici.Kapat(); err != nil {
				errors = append(errors, fmt.Errorf(i18n.T("error.workspaceCloseFailed", map[string]interface{}{"ID": id, "Error": err})))
			}
		}
	}

	// Clear workspace map
	wm.workspaces = make(map[string]*WorkspaceContext)

	if len(errors) > 0 {
		return fmt.Errorf(i18n.T("error.workspacesCloseFailed", map[string]interface{}{"Errors": errors}))
	}

	return nil
}

// Helper functions

// generateWorkspaceID generates a unique ID for a workspace based on its path
func generateWorkspaceID(path string) string {
	// Use SHA256 hash of absolute path for consistent IDs
	hash := sha256.Sum256([]byte(path))
	return fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes of hash
}

// findMigrationsPath searches for the migrations directory in standard locations
func findMigrationsPath() string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// List of possible migration paths (relative to CWD)
	possiblePaths := []string{
		"internal/veri/migrations",       // From project root
		"../veri/migrations",             // From internal/api
		"../../internal/veri/migrations", // From test subdirectory
		filepath.Join(cwd, "internal/veri/migrations"),
		filepath.Join(filepath.Dir(cwd), "veri/migrations"),
		filepath.Join(filepath.Dir(filepath.Dir(cwd)), "internal/veri/migrations"),
	}

	for _, path := range possiblePaths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		if stat, err := os.Stat(absPath); err == nil && stat.IsDir() {
			return absPath
		}
	}

	return ""
}

// getTaskCount gets the total task count for a workspace
func (wm *WorkspaceManager) getTaskCount(isYonetici *gorev.IsYonetici) (int, error) {
	filters := make(map[string]interface{})
	gorevler, err := isYonetici.GorevListele(filters)
	if err != nil {
		return 0, err
	}
	return len(gorevler), nil
}
