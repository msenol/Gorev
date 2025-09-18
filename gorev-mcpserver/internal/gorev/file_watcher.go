package gorev

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/msenol/gorev/internal/constants"
)

// FileWatcher monitors file system changes and automatically updates related tasks
type FileWatcher struct {
	watcher      *fsnotify.Watcher
	veriYonetici VeriYoneticiInterface
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.RWMutex

	// Watched paths and their associated task IDs
	watchedPaths map[string][]string

	// Task-to-paths mapping for easy cleanup
	taskPaths map[string][]string

	// Configuration
	config FileWatcherConfig
}

// FileWatcherConfig holds configuration for file watching
type FileWatcherConfig struct {
	// Extensions to watch (e.g., [".go", ".js", ".py"])
	WatchedExtensions []string

	// Patterns to ignore (e.g., ["node_modules", ".git", "*.tmp"])
	IgnorePatterns []string

	// Debounce duration to avoid multiple events for same file
	DebounceDuration time.Duration

	// Auto-update task status on file changes
	AutoUpdateStatus bool

	// Maximum file size to watch (in bytes)
	MaxFileSize int64
}

// DefaultFileWatcherConfig returns sensible defaults
func DefaultFileWatcherConfig() FileWatcherConfig {
	return FileWatcherConfig{
		WatchedExtensions: []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h", ".md", ".txt", ".json", ".yaml", ".yml"},
		IgnorePatterns:    []string{"node_modules", ".git", ".vscode", "vendor", "build", "dist", "*.tmp", "*.log", "*.swp"},
		DebounceDuration:  500 * time.Millisecond,
		AutoUpdateStatus:  true,
		MaxFileSize:       10 * 1024 * 1024, // 10MB
	}
}

// FileChangeEvent represents a file system change event
type FileChangeEvent struct {
	Path      string    `json:"path"`
	Operation string    `json:"operation"` // "create", "write", "remove", "rename"
	Timestamp time.Time `json:"timestamp"`
	TaskIDs   []string  `json:"task_ids"`
}

// NewFileWatcher creates a new file system watcher
func NewFileWatcher(veriYonetici VeriYoneticiInterface, config FileWatcherConfig) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	fw := &FileWatcher{
		watcher:      watcher,
		veriYonetici: veriYonetici,
		ctx:          ctx,
		cancel:       cancel,
		watchedPaths: make(map[string][]string),
		taskPaths:    make(map[string][]string),
		config:       config,
	}

	// Start the event processing goroutine
	go fw.processEvents()

	return fw, nil
}

// AddTaskPath associates a file path with a task ID for monitoring
func (fw *FileWatcher) AddTaskPath(taskID string, path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	// Clean and validate path
	cleanPath := filepath.Clean(path)

	// Check if path exists and get info
	info, err := filepath.Glob(cleanPath)
	if err != nil {
		return fmt.Errorf("invalid path pattern: %w", err)
	}

	if len(info) == 0 {
		// Path doesn't exist yet, still add it for future monitoring
		log.Printf("Path %s doesn't exist yet, will monitor for creation", cleanPath)
	}

	// Add to watcher if it's a directory or if we need to watch parent directory
	watchPath := cleanPath
	if !fw.isDirectory(cleanPath) {
		// For files, watch the parent directory
		watchPath = filepath.Dir(cleanPath)
	}

	// Add to filesystem watcher
	if err := fw.watcher.Add(watchPath); err != nil {
		// Ignore "already watching" errors
		if !strings.Contains(err.Error(), "already watching") {
			return fmt.Errorf("failed to add path to watcher: %w", err)
		}
	}

	// Update internal mappings
	fw.watchedPaths[cleanPath] = append(fw.watchedPaths[cleanPath], taskID)
	fw.taskPaths[taskID] = append(fw.taskPaths[taskID], cleanPath)

	log.Printf("Added path %s for task %s", cleanPath, taskID)
	return nil
}

// RemoveTaskPath removes monitoring for a specific task-path combination
func (fw *FileWatcher) RemoveTaskPath(taskID string, path string) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	cleanPath := filepath.Clean(path)

	// Remove task from path's task list
	if tasks, exists := fw.watchedPaths[cleanPath]; exists {
		newTasks := make([]string, 0, len(tasks))
		for _, t := range tasks {
			if t != taskID {
				newTasks = append(newTasks, t)
			}
		}

		if len(newTasks) == 0 {
			// No more tasks watching this path, remove from watcher
			delete(fw.watchedPaths, cleanPath)

			watchPath := cleanPath
			if !fw.isDirectory(cleanPath) {
				watchPath = filepath.Dir(cleanPath)
			}

			// Only remove if no other paths use this watch path
			stillNeeded := false
			for p := range fw.watchedPaths {
				if fw.isDirectory(p) && p == watchPath {
					stillNeeded = true
					break
				}
				if !fw.isDirectory(p) && filepath.Dir(p) == watchPath {
					stillNeeded = true
					break
				}
			}

			if !stillNeeded {
				if err := fw.watcher.Remove(watchPath); err != nil {
					log.Printf("Failed to remove watch path %s: %v", watchPath, err)
				}
			}
		} else {
			fw.watchedPaths[cleanPath] = newTasks
		}
	}

	// Remove path from task's path list
	if paths, exists := fw.taskPaths[taskID]; exists {
		newPaths := make([]string, 0, len(paths))
		for _, p := range paths {
			if p != cleanPath {
				newPaths = append(newPaths, p)
			}
		}

		if len(newPaths) == 0 {
			delete(fw.taskPaths, taskID)
		} else {
			fw.taskPaths[taskID] = newPaths
		}
	}

	log.Printf("Removed path %s for task %s", cleanPath, taskID)
	return nil
}

// RemoveTask removes all paths associated with a task
func (fw *FileWatcher) RemoveTask(taskID string) error {
	fw.mu.RLock()
	paths := make([]string, len(fw.taskPaths[taskID]))
	copy(paths, fw.taskPaths[taskID])
	fw.mu.RUnlock()

	for _, path := range paths {
		if err := fw.RemoveTaskPath(taskID, path); err != nil {
			log.Printf("Error removing path %s for task %s: %v", path, taskID, err)
		}
	}

	return nil
}

// processEvents handles file system events from the watcher
func (fw *FileWatcher) processEvents() {
	debounceMap := make(map[string]*time.Timer)
	debounceMu := sync.Mutex{}

	for {
		select {
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Skip if file should be ignored
			if fw.shouldIgnore(event.Name) {
				continue
			}

			// Debounce events for the same file
			debounceMu.Lock()
			if timer, exists := debounceMap[event.Name]; exists {
				timer.Stop()
			}

			debounceMap[event.Name] = time.AfterFunc(fw.config.DebounceDuration, func() {
				fw.handleFileEvent(event)
				debounceMu.Lock()
				delete(debounceMap, event.Name)
				debounceMu.Unlock()
			})
			debounceMu.Unlock()

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("File watcher error: %v", err)

		case <-fw.ctx.Done():
			return
		}
	}
}

// handleFileEvent processes a single file system event
func (fw *FileWatcher) handleFileEvent(event fsnotify.Event) {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	// Find tasks associated with this file
	var affectedTasks []string
	eventPath := filepath.Clean(event.Name)

	// Check exact path matches
	if tasks, exists := fw.watchedPaths[eventPath]; exists {
		affectedTasks = append(affectedTasks, tasks...)
	}

	// Check if event affects files within watched directories
	for watchedPath, tasks := range fw.watchedPaths {
		if fw.isDirectory(watchedPath) {
			if strings.HasPrefix(eventPath, watchedPath+string(filepath.Separator)) {
				affectedTasks = append(affectedTasks, tasks...)
			}
		}
	}

	// Check if event affects watched files by checking if the file is within a watched parent directory
	for watchedPath, tasks := range fw.watchedPaths {
		if !fw.isDirectory(watchedPath) && filepath.Dir(watchedPath) == filepath.Dir(eventPath) {
			if matched, _ := filepath.Match(filepath.Base(watchedPath), filepath.Base(eventPath)); matched {
				affectedTasks = append(affectedTasks, tasks...)
			}
		}
	}

	if len(affectedTasks) == 0 {
		return
	}

	// Remove duplicates
	taskMap := make(map[string]bool)
	uniqueTasks := make([]string, 0, len(affectedTasks))
	for _, taskID := range affectedTasks {
		if !taskMap[taskID] {
			taskMap[taskID] = true
			uniqueTasks = append(uniqueTasks, taskID)
		}
	}

	// Create change event
	changeEvent := FileChangeEvent{
		Path:      eventPath,
		Operation: fw.eventOpToString(event.Op),
		Timestamp: time.Now(),
		TaskIDs:   uniqueTasks,
	}

	// Log the event
	log.Printf("File change detected: %s %s (affects tasks: %v)", changeEvent.Operation, changeEvent.Path, uniqueTasks)

	// Update affected tasks
	for _, taskID := range uniqueTasks {
		if err := fw.updateTaskOnFileChange(taskID, changeEvent); err != nil {
			log.Printf("Error updating task %s for file change: %v", taskID, err)
		}
	}
}

// updateTaskOnFileChange updates a task based on file system changes
func (fw *FileWatcher) updateTaskOnFileChange(taskID string, event FileChangeEvent) error {
	// Get current task
	gorev, err := fw.veriYonetici.GorevGetir(taskID)
	if err != nil {
		return fmt.Errorf("failed to get task %s: %w", taskID, err)
	}

	// Create interaction record
	interactionData, _ := json.Marshal(map[string]interface{}{
		"file_change":    event,
		"auto_generated": true,
	})

	// Record the file change interaction
	if err := fw.veriYonetici.AIEtkilemasimKaydet(taskID, "file_change", string(interactionData), "file_watcher"); err != nil {
		log.Printf("Failed to record AI interaction for task %s: error.interactionSaveFailed", taskID, err)
	}

	// Auto-update task status if configured
	if fw.config.AutoUpdateStatus && gorev.Durum == constants.TaskStatusPending {
		// Transition to in-progress when files are modified
		if event.Operation == "write" || event.Operation == "create" {
			gorev.Durum = constants.TaskStatusInProgress
			gorev.GuncellemeTarih = time.Now()

			// Update task status using interface-compatible parameters
			updateParams := map[string]interface{}{
				"durum": constants.TaskStatusInProgress,
			}
			if err := fw.veriYonetici.GorevGuncelle(taskID, updateParams); err != nil {
				return fmt.Errorf("failed to update task status: %w", err)
			}

			log.Printf("Auto-transitioned task %s to in-progress due to file changes", taskID)
		}
	}

	// Update last AI interaction timestamp
	if err := fw.veriYonetici.GorevSonAIEtkilesiminiGuncelle(taskID, time.Now()); err != nil {
		log.Printf("Failed to update last AI interaction for task %s: %v", taskID, err)
	}

	return nil
}

// isDirectory checks if a path represents a directory
func (fw *FileWatcher) isDirectory(path string) bool {
	info, err := filepath.Glob(path)
	if err != nil {
		return false
	}

	if len(info) == 0 {
		// Path doesn't exist, guess based on extension
		return filepath.Ext(path) == ""
	}

	// Check the first match (for glob patterns)
	if len(info) > 0 {
		stat, err := filepath.Glob(info[0])
		return err == nil && len(stat) > 0 && filepath.Ext(stat[0]) == ""
	}

	return false
}

// shouldIgnore checks if a file should be ignored based on configuration
func (fw *FileWatcher) shouldIgnore(path string) bool {
	// Check file extension
	ext := filepath.Ext(path)
	if len(fw.config.WatchedExtensions) > 0 {
		found := false
		for _, watchedExt := range fw.config.WatchedExtensions {
			if ext == watchedExt {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	// Check ignore patterns
	for _, pattern := range fw.config.IgnorePatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}

		// Check if any directory in the path matches the pattern
		parts := strings.Split(filepath.Clean(path), string(filepath.Separator))
		for _, part := range parts {
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
		}
	}

	return false
}

// eventOpToString converts fsnotify.Op to string
func (fw *FileWatcher) eventOpToString(op fsnotify.Op) string {
	switch {
	case op&fsnotify.Create == fsnotify.Create:
		return "create"
	case op&fsnotify.Write == fsnotify.Write:
		return "write"
	case op&fsnotify.Remove == fsnotify.Remove:
		return "remove"
	case op&fsnotify.Rename == fsnotify.Rename:
		return "rename"
	case op&fsnotify.Chmod == fsnotify.Chmod:
		return "chmod"
	default:
		return "unknown"
	}
}

// GetWatchedPaths returns all currently watched paths and their associated tasks
func (fw *FileWatcher) GetWatchedPaths() map[string][]string {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	result := make(map[string][]string)
	for path, tasks := range fw.watchedPaths {
		result[path] = make([]string, len(tasks))
		copy(result[path], tasks)
	}
	return result
}

// GetTaskPaths returns all paths associated with a specific task
func (fw *FileWatcher) GetTaskPaths(taskID string) []string {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	if paths, exists := fw.taskPaths[taskID]; exists {
		result := make([]string, len(paths))
		copy(result, paths)
		return result
	}
	return nil
}

// Stop stops the file watcher and cleans up resources
func (fw *FileWatcher) Stop() error {
	fw.cancel()
	return fw.watcher.Close()
}

// Close is an alias for Stop() for consistent resource cleanup interface
func (fw *FileWatcher) Close() error {
	return fw.Stop()
}

// GetStats returns statistics about the file watcher
func (fw *FileWatcher) GetStats() map[string]interface{} {
	fw.mu.RLock()
	defer fw.mu.RUnlock()

	return map[string]interface{}{
		"watched_paths_count": len(fw.watchedPaths),
		"watched_tasks_count": len(fw.taskPaths),
		"config":              fw.config,
	}
}
