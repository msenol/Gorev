package api

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// getMigrationsPath returns the absolute path to the migrations directory
func getMigrationsPath() string {
	// Assuming tests run from project root or internal/api directory
	paths := []string{
		"internal/veri/migrations",       // From project root
		"../veri/migrations",             // From internal/api
		"../../internal/veri/migrations", // From test subdirectory
		"/home/msenol/Projects/Gorev/gorev-mcpserver/internal/veri/migrations", // Absolute fallback
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			absPath, _ := filepath.Abs(p)
			return absPath
		}
	}

	return ""
}

// setupTestWorkspaceDir creates a temporary workspace directory for testing
func setupTestWorkspaceDir(t *testing.T) string {
	tmpDir := t.TempDir()
	workspaceDir := filepath.Join(tmpDir, "test-workspace")
	err := os.MkdirAll(workspaceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test workspace dir: %v", err)
	}
	return workspaceDir
}

func TestNewWorkspaceManager(t *testing.T) {
	wm := NewWorkspaceManager()
	if wm == nil {
		t.Fatal("NewWorkspaceManager returned nil")
	}
	if wm.workspaces == nil {
		t.Fatal("WorkspaceManager.workspaces is nil")
	}
	if len(wm.workspaces) != 0 {
		t.Fatalf("Expected empty workspaces map, got %d items", len(wm.workspaces))
	}
}

func TestRegisterWorkspace_Success(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	// Register workspace
	ws, err := wm.RegisterWorkspace(workspaceDir, "Test Workspace")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Verify workspace properties
	if ws.Name != "Test Workspace" {
		t.Errorf("Expected name 'Test Workspace', got '%s'", ws.Name)
	}
	if ws.Path != workspaceDir {
		t.Errorf("Expected path '%s', got '%s'", workspaceDir, ws.Path)
	}
	if ws.ID == "" {
		t.Error("Workspace ID is empty")
	}
	if ws.VeriYonetici == nil {
		t.Error("VeriYonetici is nil")
	}
	if ws.IsYonetici == nil {
		t.Error("IsYonetici is nil")
	}

	// Verify .gorev directory was created
	gorevDir := filepath.Join(workspaceDir, ".gorev")
	if _, err := os.Stat(gorevDir); os.IsNotExist(err) {
		t.Errorf(".gorev directory not created at %s", gorevDir)
	}

	// Verify database file exists
	dbPath := filepath.Join(gorevDir, "gorev.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file not created at %s", dbPath)
	}
}

func TestRegisterWorkspace_AutoName(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	// Register without name
	ws, err := wm.RegisterWorkspace(workspaceDir, "")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Should use directory name
	expectedName := filepath.Base(workspaceDir)
	if ws.Name != expectedName {
		t.Errorf("Expected auto-generated name '%s', got '%s'", expectedName, ws.Name)
	}
}

func TestRegisterWorkspace_Duplicate(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	// Register first time
	ws1, err := wm.RegisterWorkspace(workspaceDir, "First")
	if err != nil {
		t.Fatalf("First RegisterWorkspace failed: %v", err)
	}

	// Register same path again
	ws2, err := wm.RegisterWorkspace(workspaceDir, "Second")
	if err != nil {
		t.Fatalf("Second RegisterWorkspace failed: %v", err)
	}

	// Should return existing workspace (same ID)
	if ws1.ID != ws2.ID {
		t.Errorf("Expected same workspace ID, got %s and %s", ws1.ID, ws2.ID)
	}

	// Name should remain unchanged (from first registration)
	if ws2.Name != "First" {
		t.Errorf("Expected name to remain 'First', got '%s'", ws2.Name)
	}

	// Should only have one workspace in registry
	workspaces := wm.ListWorkspaces()
	if len(workspaces) != 1 {
		t.Errorf("Expected 1 workspace in registry, got %d", len(workspaces))
	}
}

func TestRegisterWorkspace_InvalidPath(t *testing.T) {
	wm := NewWorkspaceManager()
	invalidPath := "/this/path/does/not/exist/at/all"

	ws, err := wm.RegisterWorkspace(invalidPath, "Invalid")
	if err == nil {
		t.Fatal("Expected error for invalid path, got nil")
	}
	if ws != nil {
		t.Errorf("Expected nil workspace for invalid path, got %v", ws)
	}
}

func TestGetWorkspace_Success(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	// Register workspace
	registered, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Get workspace
	retrieved, err := wm.GetWorkspace(registered.ID)
	if err != nil {
		t.Fatalf("GetWorkspace failed: %v", err)
	}

	// Verify it's the same workspace
	if retrieved == nil {
		t.Fatal("GetWorkspace returned nil")
	}

	// Type assert to WorkspaceContext
	ws, ok := retrieved.(*WorkspaceContext)
	if !ok {
		t.Fatalf("GetWorkspace returned wrong type: %T", retrieved)
	}

	if ws.ID != registered.ID {
		t.Errorf("Expected ID %s, got %s", registered.ID, ws.ID)
	}
}

func TestGetWorkspace_NotFound(t *testing.T) {
	wm := NewWorkspaceManager()

	ws, err := wm.GetWorkspace("nonexistent-id")
	if err == nil {
		t.Fatal("Expected error for non-existent workspace, got nil")
	}
	if ws != nil {
		t.Errorf("Expected nil workspace, got %v", ws)
	}
}

func TestGetWorkspaceContext_Success(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	// Register workspace
	registered, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Get workspace context
	ws, err := wm.GetWorkspaceContext(registered.ID)
	if err != nil {
		t.Fatalf("GetWorkspaceContext failed: %v", err)
	}

	if ws.ID != registered.ID {
		t.Errorf("Expected ID %s, got %s", registered.ID, ws.ID)
	}
}

func TestListWorkspaces_Empty(t *testing.T) {
	wm := NewWorkspaceManager()

	workspaces := wm.ListWorkspaces()
	if workspaces == nil {
		t.Fatal("ListWorkspaces returned nil")
	}
	if len(workspaces) != 0 {
		t.Errorf("Expected empty list, got %d workspaces", len(workspaces))
	}
}

func TestListWorkspaces_Multiple(t *testing.T) {
	wm := NewWorkspaceManager()

	// Create multiple test workspaces
	tmpDir := t.TempDir()
	ws1Dir := filepath.Join(tmpDir, "workspace1")
	ws2Dir := filepath.Join(tmpDir, "workspace2")
	ws3Dir := filepath.Join(tmpDir, "workspace3")

	for _, dir := range []string{ws1Dir, ws2Dir, ws3Dir} {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}
		_, err = wm.RegisterWorkspace(dir, filepath.Base(dir))
		if err != nil {
			t.Fatalf("RegisterWorkspace failed: %v", err)
		}
	}

	// List all workspaces
	workspaces := wm.ListWorkspaces()
	if len(workspaces) != 3 {
		t.Errorf("Expected 3 workspaces, got %d", len(workspaces))
	}
}

func TestUnregisterWorkspace_Success(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	// Register workspace
	ws, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Verify it's registered
	workspaces := wm.ListWorkspaces()
	if len(workspaces) != 1 {
		t.Fatalf("Expected 1 workspace before unregister, got %d", len(workspaces))
	}

	// Unregister
	err = wm.UnregisterWorkspace(ws.ID)
	if err != nil {
		t.Fatalf("UnregisterWorkspace failed: %v", err)
	}

	// Verify it's removed
	workspaces = wm.ListWorkspaces()
	if len(workspaces) != 0 {
		t.Errorf("Expected 0 workspaces after unregister, got %d", len(workspaces))
	}

	// Verify can't get it anymore
	_, err = wm.GetWorkspace(ws.ID)
	if err == nil {
		t.Error("Expected error when getting unregistered workspace")
	}
}

func TestUnregisterWorkspace_NotFound(t *testing.T) {
	wm := NewWorkspaceManager()

	err := wm.UnregisterWorkspace("nonexistent-id")
	if err == nil {
		t.Fatal("Expected error for non-existent workspace, got nil")
	}
}

func TestUpdateTaskCount(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	// Register workspace
	ws, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Initial task count should be 0
	if ws.TaskCount != 0 {
		t.Errorf("Expected initial task count 0, got %d", ws.TaskCount)
	}

	// Update task count
	err = wm.UpdateTaskCount(ws.ID)
	if err != nil {
		t.Fatalf("UpdateTaskCount failed: %v", err)
	}

	// Task count should still be 0 (no tasks created)
	updatedWs, _ := wm.GetWorkspaceContext(ws.ID)
	if updatedWs.TaskCount != 0 {
		t.Errorf("Expected task count 0 after update, got %d", updatedWs.TaskCount)
	}
}

func TestCleanup(t *testing.T) {
	wm := NewWorkspaceManager()

	// Create workspace
	workspaceDir := setupTestWorkspaceDir(t)
	ws, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Set last accessed to old time
	wm.mu.Lock()
	ws.LastAccessed = time.Now().Add(-2 * time.Hour)
	wm.mu.Unlock()

	// Cleanup workspaces older than 1 hour
	err = wm.Cleanup(1 * time.Hour)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Workspace should be removed
	workspaces := wm.ListWorkspaces()
	if len(workspaces) != 0 {
		t.Errorf("Expected 0 workspaces after cleanup, got %d", len(workspaces))
	}
}

func TestCleanup_RecentWorkspace(t *testing.T) {
	wm := NewWorkspaceManager()

	// Create workspace
	workspaceDir := setupTestWorkspaceDir(t)
	ws, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Last accessed is recent (now)
	// Cleanup workspaces older than 1 hour - should not remove this one
	err = wm.Cleanup(1 * time.Hour)
	if err != nil {
		t.Fatalf("Cleanup failed: %v", err)
	}

	// Workspace should still exist
	workspaces := wm.ListWorkspaces()
	if len(workspaces) != 1 {
		t.Errorf("Expected 1 workspace after cleanup, got %d", len(workspaces))
	}
	if workspaces[0].ID != ws.ID {
		t.Error("Wrong workspace remained after cleanup")
	}
}

func TestCloseAll(t *testing.T) {
	wm := NewWorkspaceManager()

	// Create multiple workspaces
	tmpDir := t.TempDir()
	for i := 1; i <= 3; i++ {
		dir := filepath.Join(tmpDir, "workspace"+string(rune('0'+i)))
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		_, err = wm.RegisterWorkspace(dir, "Test"+string(rune('0'+i)))
		if err != nil {
			t.Fatalf("RegisterWorkspace failed: %v", err)
		}
	}

	// Verify workspaces exist
	workspaces := wm.ListWorkspaces()
	if len(workspaces) != 3 {
		t.Fatalf("Expected 3 workspaces, got %d", len(workspaces))
	}

	// Close all
	err := wm.CloseAll()
	if err != nil {
		t.Fatalf("CloseAll failed: %v", err)
	}

	// All workspaces should be removed
	workspaces = wm.ListWorkspaces()
	if len(workspaces) != 0 {
		t.Errorf("Expected 0 workspaces after CloseAll, got %d", len(workspaces))
	}
}

func TestConcurrentAccess(t *testing.T) {
	wm := NewWorkspaceManager()

	// Create test workspace
	workspaceDir := setupTestWorkspaceDir(t)
	ws, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Concurrent reads
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := wm.GetWorkspace(ws.ID)
			if err != nil {
				t.Errorf("Concurrent GetWorkspace failed: %v", err)
			}
		}()
	}

	// Concurrent list operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = wm.ListWorkspaces()
		}()
	}

	wg.Wait()
}

func TestWorkspaceIDGeneration(t *testing.T) {
	wm := NewWorkspaceManager()

	// Create two different workspaces
	tmpDir := t.TempDir()
	ws1Dir := filepath.Join(tmpDir, "workspace1")
	ws2Dir := filepath.Join(tmpDir, "workspace2")

	for _, dir := range []string{ws1Dir, ws2Dir} {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
	}

	ws1, _ := wm.RegisterWorkspace(ws1Dir, "WS1")
	ws2, _ := wm.RegisterWorkspace(ws2Dir, "WS2")

	// IDs should be different
	if ws1.ID == ws2.ID {
		t.Error("Different workspaces should have different IDs")
	}

	// Same path should generate same ID
	ws1Again, _ := wm.RegisterWorkspace(ws1Dir, "WS1-Again")
	if ws1.ID != ws1Again.ID {
		t.Error("Same workspace path should generate same ID")
	}
}

func TestWorkspaceContextInterface(t *testing.T) {
	wm := NewWorkspaceManager()
	workspaceDir := setupTestWorkspaceDir(t)

	ws, err := wm.RegisterWorkspace(workspaceDir, "Test")
	if err != nil {
		t.Fatalf("RegisterWorkspace failed: %v", err)
	}

	// Test GetIsYonetici
	isYonetici := ws.GetIsYonetici()
	if isYonetici == nil {
		t.Error("GetIsYonetici returned nil")
	}

	// Test ToWorkspaceInfo
	info := ws.ToWorkspaceInfo()
	if info == nil {
		t.Fatal("ToWorkspaceInfo returned nil")
	}
	if info.ID != ws.ID {
		t.Errorf("WorkspaceInfo ID mismatch: expected %s, got %s", ws.ID, info.ID)
	}
	if info.Name != ws.Name {
		t.Errorf("WorkspaceInfo Name mismatch: expected %s, got %s", ws.Name, info.Name)
	}
}
