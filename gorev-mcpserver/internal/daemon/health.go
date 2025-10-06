package daemon

import (
	"fmt"
	"net/http"
	"time"
)

// IsDaemonHealthy checks if daemon at given URL is responding to health checks
func IsDaemonHealthy(daemonURL string) bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(daemonURL + "/api/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// WaitForDaemon waits for daemon to be ready at given URL with timeout
func WaitForDaemon(daemonURL string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if IsDaemonHealthy(daemonURL) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for daemon to be ready after %v", timeout)
}

// GetDaemonInfo returns daemon information if running, or nil if not found/unhealthy
func GetDaemonInfo() (*LockFile, error) {
	lock, err := ReadLockFile()
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	// Verify process is running
	if !IsProcessRunning(lock.PID) {
		// Stale lock file
		RemoveLockFile()
		return nil, fmt.Errorf("stale lock file found (PID %d not running)", lock.PID)
	}

	// Verify daemon is healthy
	if !IsDaemonHealthy(lock.DaemonURL) {
		// Process running but daemon unhealthy
		RemoveLockFile()
		return nil, fmt.Errorf("daemon unhealthy at %s", lock.DaemonURL)
	}

	return lock, nil
}
