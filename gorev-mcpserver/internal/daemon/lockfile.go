package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// LockFile represents the daemon process lock file structure
type LockFile struct {
	PID       int       `json:"pid"`
	Port      string    `json:"port"`
	StartTime time.Time `json:"start_time"`
	DaemonURL string    `json:"daemon_url"`
	Version   string    `json:"version"`
}

// GetLockFilePath returns the absolute path to the daemon lock file
func GetLockFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to /tmp if home directory not available
		return filepath.Join("/tmp", ".gorev-daemon", ".lock")
	}
	return filepath.Join(home, ".gorev-daemon", ".lock")
}

// CreateLockFile creates a lock file for the daemon process
func CreateLockFile(pid int, port, version string) error {
	lockPath := GetLockFilePath()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(lockPath), 0755); err != nil {
		return fmt.Errorf("failed to create lock directory: %w", err)
	}

	lock := LockFile{
		PID:       pid,
		Port:      port,
		StartTime: time.Now(),
		DaemonURL: fmt.Sprintf("http://localhost:%s", port),
		Version:   version,
	}

	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	if err := os.WriteFile(lockPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}

// ReadLockFile reads and parses the daemon lock file
func ReadLockFile() (*LockFile, error) {
	lockPath := GetLockFilePath()

	data, err := os.ReadFile(lockPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("lock file not found")
		}
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	var lock LockFile
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lock, nil
}

// IsProcessRunning checks if a process with given PID is running
func IsProcessRunning(pid int) bool {
	// Try to send signal 0 (no-op signal) to check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix systems, FindProcess always succeeds, so we need to send a signal
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}

	// Process doesn't exist or we don't have permission
	return false
}

// RemoveLockFile removes the daemon lock file
func RemoveLockFile() error {
	lockPath := GetLockFilePath()
	if err := os.Remove(lockPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already removed
		}
		return fmt.Errorf("failed to remove lock file: %w", err)
	}
	return nil
}

// UpdateLastAccess updates the lock file's start time (can be used for heartbeat)
func UpdateLastAccess() error {
	lock, err := ReadLockFile()
	if err != nil {
		return err
	}

	// Update start time to current time (acts as "last seen" timestamp)
	lock.StartTime = time.Now()

	lockPath := GetLockFilePath()
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	if err := os.WriteFile(lockPath, data, 0644); err != nil {
		return fmt.Errorf("failed to update lock file: %w", err)
	}

	return nil
}
