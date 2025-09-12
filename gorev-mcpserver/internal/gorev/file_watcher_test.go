package gorev

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/msenol/gorev/internal/constants"
)

// TestDefaultFileWatcherConfig tests the default configuration
func TestDefaultFileWatcherConfig(t *testing.T) {
	config := DefaultFileWatcherConfig()
	
	// Test default extensions
	expectedExtensions := []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h", ".md", ".txt", ".json", ".yaml", ".yml"}
	if len(config.WatchedExtensions) != len(expectedExtensions) {
		t.Errorf("Expected %d extensions, got %d", len(expectedExtensions), len(config.WatchedExtensions))
	}
	
	// Test some key extensions
	found := make(map[string]bool)
	for _, ext := range config.WatchedExtensions {
		found[ext] = true
	}
	for _, expected := range []string{".go", ".js", ".py", ".md"} {
		if !found[expected] {
			t.Errorf("Expected extension %s not found in default config", expected)
		}
	}
	
	// Test ignore patterns
	expectedIgnores := []string{"node_modules", ".git", ".vscode", "vendor", "build", "dist", "*.tmp", "*.log", "*.swp"}
	if len(config.IgnorePatterns) != len(expectedIgnores) {
		t.Errorf("Expected %d ignore patterns, got %d", len(expectedIgnores), len(config.IgnorePatterns))
	}
	
	// Test other defaults
	if config.DebounceDuration != 500*time.Millisecond {
		t.Errorf("Expected debounce duration 500ms, got %v", config.DebounceDuration)
	}
	
	if !config.AutoUpdateStatus {
		t.Error("Expected AutoUpdateStatus to be true by default")
	}
	
	if config.MaxFileSize != 10*1024*1024 {
		t.Errorf("Expected max file size 10MB, got %d", config.MaxFileSize)
	}
}

// TestNewFileWatcher tests the FileWatcher constructor
func TestNewFileWatcher(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Test that FileWatcher is properly initialized
	if fw.veriYonetici != vy {
		t.Error("VeriYonetici not set correctly")
	}
	
	if fw.watcher == nil {
		t.Error("FSNotify watcher not created")
	}
	
	if fw.watchedPaths == nil {
		t.Error("watchedPaths map not initialized")
	}
	
	if fw.taskPaths == nil {
		t.Error("taskPaths map not initialized")
	}
	
	if fw.ctx == nil {
		t.Error("Context not created")
	}
	
	if fw.cancel == nil {
		t.Error("Cancel function not created")
	}
	
	// Test configuration is stored
	if len(fw.config.WatchedExtensions) == 0 {
		t.Error("Config not stored properly")
	}
}

// TestAddTaskPath tests adding task-path associations
func TestFileWatcher_AddTaskPath(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testFile := filepath.Join(tempDir, "test.go")
	
	// Test adding a file path
	err = fw.AddTaskPath("task-1", testFile)
	if err != nil {
		t.Errorf("Failed to add task path: %v", err)
	}
	
	// Verify the path was added
	fw.mu.RLock()
	if tasks, exists := fw.watchedPaths[testFile]; !exists || len(tasks) != 1 || tasks[0] != "task-1" {
		t.Errorf("Task path not added correctly: %v", fw.watchedPaths)
	}
	
	if paths, exists := fw.taskPaths["task-1"]; !exists || len(paths) != 1 || paths[0] != testFile {
		t.Errorf("Task paths not updated correctly: %v", fw.taskPaths)
	}
	fw.mu.RUnlock()
	
	// Test adding multiple tasks to same path
	err = fw.AddTaskPath("task-2", testFile)
	if err != nil {
		t.Errorf("Failed to add second task to path: %v", err)
	}
	
	fw.mu.RLock()
	tasks := fw.watchedPaths[testFile]
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks for path, got %d", len(tasks))
	}
	fw.mu.RUnlock()
	
	// Test adding directory path
	err = fw.AddTaskPath("task-3", tempDir)
	if err != nil {
		t.Errorf("Failed to add directory path: %v", err)
	}
	
	// Test adding non-existent path (should not error)
	nonExistentPath := filepath.Join(tempDir, "nonexistent.txt")
	err = fw.AddTaskPath("task-4", nonExistentPath)
	if err != nil {
		t.Errorf("Failed to add non-existent path: %v", err)
	}
}

// TestRemoveTaskPath tests removing task-path associations
func TestFileWatcher_RemoveTaskPath(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testFile := filepath.Join(tempDir, "test.go")
	
	// Add some tasks
	fw.AddTaskPath("task-1", testFile)
	fw.AddTaskPath("task-2", testFile)
	
	// Remove one task
	err = fw.RemoveTaskPath("task-1", testFile)
	if err != nil {
		t.Errorf("Failed to remove task path: %v", err)
	}
	
	// Verify task-1 was removed but task-2 remains
	fw.mu.RLock()
	tasks := fw.watchedPaths[testFile]
	if len(tasks) != 1 || tasks[0] != "task-2" {
		t.Errorf("Expected only task-2 to remain, got: %v", tasks)
	}
	
	if _, exists := fw.taskPaths["task-1"]; exists {
		t.Error("task-1 should have been removed from taskPaths")
	}
	fw.mu.RUnlock()
	
	// Remove last task (should clean up path completely)
	err = fw.RemoveTaskPath("task-2", testFile)
	if err != nil {
		t.Errorf("Failed to remove last task from path: %v", err)
	}
	
	fw.mu.RLock()
	if _, exists := fw.watchedPaths[testFile]; exists {
		t.Error("Path should have been removed when no tasks remain")
	}
	
	if _, exists := fw.taskPaths["task-2"]; exists {
		t.Error("task-2 should have been removed from taskPaths")
	}
	fw.mu.RUnlock()
}

// TestRemoveTask tests removing all paths for a task
func TestFileWatcher_RemoveTask(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Create temporary files for testing
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testFile1 := filepath.Join(tempDir, "test1.go")
	testFile2 := filepath.Join(tempDir, "test2.go")
	
	// Add multiple paths for one task
	fw.AddTaskPath("task-1", testFile1)
	fw.AddTaskPath("task-1", testFile2)
	fw.AddTaskPath("task-2", testFile1) // Another task using same path
	
	// Remove all paths for task-1
	err = fw.RemoveTask("task-1")
	if err != nil {
		t.Errorf("Failed to remove task: %v", err)
	}
	
	// Verify task-1 is completely removed
	fw.mu.RLock()
	if _, exists := fw.taskPaths["task-1"]; exists {
		t.Error("task-1 should have been completely removed")
	}
	
	// Verify task-2 is still there for testFile1
	if tasks := fw.watchedPaths[testFile1]; len(tasks) != 1 || tasks[0] != "task-2" {
		t.Errorf("task-2 should still be watching testFile1, got: %v", tasks)
	}
	
	// testFile2 should have no watchers now
	if _, exists := fw.watchedPaths[testFile2]; exists {
		t.Error("testFile2 should have no watchers after task-1 removal")
	}
	fw.mu.RUnlock()
}

// TestShouldIgnore tests the ignore patterns functionality
func TestFileWatcher_ShouldIgnore(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	testCases := []struct {
		name           string
		path           string
		expectedIgnore bool
		description    string
	}{
		{
			name:           "Go file should not be ignored",
			path:           "/path/to/file.go",
			expectedIgnore: false,
			description:    "Go files are in watched extensions",
		},
		{
			name:           "Tmp file should be ignored",
			path:           "/path/to/file.tmp",
			expectedIgnore: true,
			description:    "Tmp files match ignore pattern",
		},
		{
			name:           "Node modules should be ignored",
			path:           "/project/node_modules/package/index.js",
			expectedIgnore: true,
			description:    "Node modules directory should be ignored",
		},
		{
			name:           "Git file should be ignored",
			path:           "/project/.git/config",
			expectedIgnore: true,
			description:    "Git directory should be ignored",
		},
		{
			name:           "Log file should be ignored",
			path:           "/path/to/app.log",
			expectedIgnore: true,
			description:    "Log files match ignore pattern",
		},
		{
			name:           "Unknown extension should be ignored",
			path:           "/path/to/file.xyz",
			expectedIgnore: true,
			description:    "Unknown extensions not in watch list should be ignored",
		},
		{
			name:           "Python file should not be ignored",
			path:           "/path/to/script.py",
			expectedIgnore: false,
			description:    "Python files are in watched extensions",
		},
		{
			name:           "VSCode file should be ignored",
			path:           "/project/.vscode/settings.json",
			expectedIgnore: true,
			description:    "VSCode directory should be ignored",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fw.shouldIgnore(tc.path)
			if result != tc.expectedIgnore {
				t.Errorf("Expected shouldIgnore(%s) to be %v, got %v. %s", 
					tc.path, tc.expectedIgnore, result, tc.description)
			}
		})
	}
}

// TestEventOpToString tests the event operation conversion
func TestFileWatcher_EventOpToString(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	testCases := []struct {
		op       fsnotify.Op
		expected string
	}{
		{fsnotify.Create, "create"},
		{fsnotify.Write, "write"},
		{fsnotify.Remove, "remove"},
		{fsnotify.Rename, "rename"},
		{fsnotify.Chmod, "chmod"},
		{fsnotify.Op(0), "unknown"}, // Unknown operation (0 value)
	}
	
	for _, tc := range testCases {
		result := fw.eventOpToString(tc.op)
		if result != tc.expected {
			t.Errorf("Expected eventOpToString(%v) to be %s, got %s", tc.op, tc.expected, result)
		}
	}
}

// TestGetWatchedPaths tests retrieving watched paths
func TestFileWatcher_GetWatchedPaths(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Initially should be empty
	paths := fw.GetWatchedPaths()
	if len(paths) != 0 {
		t.Errorf("Expected no watched paths initially, got %d", len(paths))
	}
	
	// Create temporary files for testing
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testFile1 := filepath.Join(tempDir, "file1.go")
	testFile2 := filepath.Join(tempDir, "file2.go")
	
	// Create the files
	if err := os.WriteFile(testFile1, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create test file1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create test file2: %v", err)
	}
	
	// Add some paths
	fw.AddTaskPath("task-1", testFile1)
	fw.AddTaskPath("task-2", testFile1)
	fw.AddTaskPath("task-1", testFile2)
	
	paths = fw.GetWatchedPaths()
	if len(paths) != 2 {
		t.Errorf("Expected 2 watched paths, got %d", len(paths))
	}
	
	// Check that file1.go has 2 tasks
	if tasks, exists := paths[testFile1]; !exists || len(tasks) != 2 {
		t.Errorf("Expected file1.go to have 2 tasks, got: %v", tasks)
	}
	
	// Check that file2.go has 1 task  
	if tasks, exists := paths[testFile2]; !exists || len(tasks) != 1 {
		t.Errorf("Expected file2.go to have 1 task, got: %v", tasks)
	}
	
	// Verify the returned map is a copy (mutations shouldn't affect original)
	delete(paths, testFile1)
	paths2 := fw.GetWatchedPaths()
	if len(paths2) != 2 {
		t.Error("GetWatchedPaths should return a copy, original should not be affected")
	}
}

// TestGetTaskPaths tests retrieving paths for a specific task
func TestFileWatcher_GetTaskPaths(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Test non-existent task
	paths := fw.GetTaskPaths("non-existent")
	if paths != nil {
		t.Errorf("Expected nil for non-existent task, got: %v", paths)
	}
	
	// Create temporary files for testing
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testFile1 := filepath.Join(tempDir, "file1.go")
	testFile2 := filepath.Join(tempDir, "file2.go")
	testFile3 := filepath.Join(tempDir, "file3.go")
	
	// Create the files
	for _, file := range []string{testFile1, testFile2, testFile3} {
		if err := os.WriteFile(file, []byte("package main"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}
	
	// Add paths for a task
	fw.AddTaskPath("task-1", testFile1)
	fw.AddTaskPath("task-1", testFile2)
	fw.AddTaskPath("task-2", testFile3)
	
	// Test task-1 paths
	paths = fw.GetTaskPaths("task-1")
	if len(paths) != 2 {
		t.Errorf("Expected 2 paths for task-1, got %d", len(paths))
	}
	
	expectedPaths := map[string]bool{
		testFile1: true,
		testFile2: true,
	}
	
	for _, path := range paths {
		if !expectedPaths[path] {
			t.Errorf("Unexpected path for task-1: %s", path)
		}
	}
	
	// Test task-2 paths
	paths = fw.GetTaskPaths("task-2")
	if len(paths) != 1 || paths[0] != testFile3 {
		t.Errorf("Expected 1 path '%s' for task-2, got: %v", testFile3, paths)
	}
	
	// Verify the returned slice is a copy
	originalPath := paths[0]
	paths[0] = "modified"
	paths2 := fw.GetTaskPaths("task-2")
	if len(paths2) == 0 {
		t.Error("GetTaskPaths should return a copy, got empty slice")
	} else if paths2[0] != originalPath {
		t.Error("GetTaskPaths should return a copy, original should not be affected")
	}
}

// TestGetStats tests the statistics functionality
func TestFileWatcher_GetStats(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Test initial stats
	stats := fw.GetStats()
	
	if watchedCount, ok := stats["watched_paths_count"].(int); !ok || watchedCount != 0 {
		t.Errorf("Expected watched_paths_count to be 0, got: %v", stats["watched_paths_count"])
	}
	
	if tasksCount, ok := stats["watched_tasks_count"].(int); !ok || tasksCount != 0 {
		t.Errorf("Expected watched_tasks_count to be 0, got: %v", stats["watched_tasks_count"])
	}
	
	if configVal, ok := stats["config"].(FileWatcherConfig); !ok {
		t.Error("Expected config to be present in stats")
	} else {
		if len(configVal.WatchedExtensions) == 0 {
			t.Error("Config in stats should contain watched extensions")
		}
	}
	
	// Create temporary files for testing
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	testFile1 := filepath.Join(tempDir, "file1.go")
	testFile2 := filepath.Join(tempDir, "file2.go")
	
	// Create the files
	for _, file := range []string{testFile1, testFile2} {
		if err := os.WriteFile(file, []byte("package main"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}
	
	// Add some paths and check updated stats
	fw.AddTaskPath("task-1", testFile1)
	fw.AddTaskPath("task-1", testFile2)
	fw.AddTaskPath("task-2", testFile1)
	
	stats = fw.GetStats()
	
	if watchedCount := stats["watched_paths_count"].(int); watchedCount != 2 {
		t.Errorf("Expected watched_paths_count to be 2, got: %v", watchedCount)
	}
	
	if tasksCount := stats["watched_tasks_count"].(int); tasksCount != 2 {
		t.Errorf("Expected watched_tasks_count to be 2, got: %v", tasksCount)
	}
}

// TestUpdateTaskOnFileChange tests task updates based on file changes
func TestFileWatcher_UpdateTaskOnFileChange(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Create test task
	testTask := &Gorev{
		ID:              "test-task",
		Baslik:          "Test Task",
		Durum:           constants.TaskStatusPending,
		Oncelik:         constants.PriorityMedium,
		OlusturmaTarih:  time.Now().Add(-1 * time.Hour),
		GuncellemeTarih: time.Now().Add(-1 * time.Hour),
	}
	vy.gorevler["test-task"] = testTask
	
	// Create file change event
	changeEvent := FileChangeEvent{
		Path:      "/path/to/file.go",
		Operation: "write",
		Timestamp: time.Now(),
		TaskIDs:   []string{"test-task"},
	}
	
	// Test task update on file change
	err = fw.updateTaskOnFileChange("test-task", changeEvent)
	if err != nil {
		t.Errorf("updateTaskOnFileChange failed: %v", err)
	}
	
	// Verify task status was updated to in-progress
	updatedTask := vy.gorevler["test-task"]
	if updatedTask.Durum != constants.TaskStatusInProgress {
		t.Errorf("Expected task status to be in-progress, got: %s", updatedTask.Durum)
	}
	
	// Test with task that doesn't exist
	err = fw.updateTaskOnFileChange("non-existent", changeEvent)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
	
	// Test with task already in progress (should not change status again)
	inProgressTask := &Gorev{
		ID:              "progress-task", 
		Baslik:          "In Progress Task",
		Durum:           constants.TaskStatusInProgress,
		Oncelik:         constants.PriorityMedium,
		OlusturmaTarih:  time.Now().Add(-2 * time.Hour),
		GuncellemeTarih: time.Now().Add(-1 * time.Hour),
	}
	vy.gorevler["progress-task"] = inProgressTask
	
	err = fw.updateTaskOnFileChange("progress-task", changeEvent)
	if err != nil {
		t.Errorf("updateTaskOnFileChange failed for in-progress task: %v", err)
	}
	
	// Should remain in progress
	if vy.gorevler["progress-task"].Durum != constants.TaskStatusInProgress {
		t.Error("In-progress task status should remain unchanged")
	}
}

// TestStopAndClose tests cleanup functionality
func TestFileWatcher_StopAndClose(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	
	// Add some paths
	fw.AddTaskPath("task-1", "/path/to/file.go")
	
	// Test Stop
	err = fw.Stop()
	if err != nil {
		t.Errorf("Stop() failed: %v", err)
	}
	
	// Test Close (should be alias for Stop)
	fw2, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create second FileWatcher: %v", err)
	}
	
	err = fw2.Close()
	if err != nil {
		t.Errorf("Close() failed: %v", err)
	}
}

// TestIsDirectory tests directory detection
func TestFileWatcher_IsDirectory(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Create a temporary directory for testing real paths
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test file
	testFile := filepath.Join(tempDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	testCases := []struct {
		name        string
		path        string
		expectedDir bool
		description string
	}{
		{
			name:        "Path with extension should not be directory",
			path:        "/path/to/file.go",
			expectedDir: false,
			description: "Files with extensions are not directories (non-existent paths)",
		},
		{
			name:        "Path without extension should be directory",
			path:        "/path/to/directory",
			expectedDir: true,
			description: "Paths without extensions are assumed to be directories (non-existent paths)",
		},
		{
			name:        "Root path should be directory",
			path:        "/",
			expectedDir: true,
			description: "Root path is a directory",
		},
		{
			name:        "Existing directory should be directory",
			path:        tempDir,
			expectedDir: true,
			description: "Existing directories are correctly identified",
		},
		{
			name:        "Existing file should not be directory",
			path:        testFile,
			expectedDir: false,
			description: "Existing files with extensions are not directories",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fw.isDirectory(tc.path)
			if result != tc.expectedDir {
				t.Errorf("Expected isDirectory(%s) to be %v, got %v. %s", 
					tc.path, tc.expectedDir, result, tc.description)
			}
		})
	}
}

// TestFileChangeEvent tests FileChangeEvent struct
func TestFileChangeEvent(t *testing.T) {
	event := FileChangeEvent{
		Path:      "/path/to/file.go",
		Operation: "write",
		Timestamp: time.Now(),
		TaskIDs:   []string{"task-1", "task-2"},
	}
	
	if event.Path != "/path/to/file.go" {
		t.Errorf("Expected path '/path/to/file.go', got %s", event.Path)
	}
	
	if event.Operation != "write" {
		t.Errorf("Expected operation 'write', got %s", event.Operation)
	}
	
	if len(event.TaskIDs) != 2 {
		t.Errorf("Expected 2 task IDs, got %d", len(event.TaskIDs))
	}
	
	if event.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero")
	}
}

// TestConcurrentOperations tests thread safety
func TestFileWatcher_ConcurrentOperations(t *testing.T) {
	vy := NewMockVeriYonetici()
	config := DefaultFileWatcherConfig()
	
	fw, err := NewFileWatcher(vy, config)
	if err != nil {
		t.Fatalf("Failed to create FileWatcher: %v", err)
	}
	defer fw.Close()
	
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "filewatch_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Test concurrent adds and removes
	var wg sync.WaitGroup
	numGoroutines := 10
	
	// Create test files
	testFiles := make([]string, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		testFiles[i] = filepath.Join(tempDir, fmt.Sprintf("file-%d.go", i))
		if err := os.WriteFile(testFiles[i], []byte("package main"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", testFiles[i], err)
		}
	}
	
	// Concurrent AddTaskPath operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			taskID := fmt.Sprintf("task-%d", id)
			path := testFiles[id]
			err := fw.AddTaskPath(taskID, path)
			if err != nil {
				t.Errorf("Concurrent AddTaskPath failed: %v", err)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify all tasks were added
	stats := fw.GetStats()
	if watchedPaths := stats["watched_paths_count"].(int); watchedPaths != numGoroutines {
		t.Errorf("Expected %d watched paths, got %d", numGoroutines, watchedPaths)
	}
	
	// Concurrent GetWatchedPaths operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			paths := fw.GetWatchedPaths()
			if len(paths) != numGoroutines {
				t.Errorf("Concurrent GetWatchedPaths returned wrong count: %d", len(paths))
			}
		}()
	}
	
	wg.Wait()
	
	// Concurrent RemoveTask operations
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			taskID := fmt.Sprintf("task-%d", id)
			err := fw.RemoveTask(taskID)
			if err != nil {
				t.Errorf("Concurrent RemoveTask failed: %v", err)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify some tasks were removed
	stats = fw.GetStats()
	remainingTasks := stats["watched_tasks_count"].(int)
	if remainingTasks != numGoroutines/2 {
		t.Errorf("Expected %d remaining tasks, got %d", numGoroutines/2, remainingTasks)
	}
}
