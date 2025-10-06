package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/msenol/gorev/internal/daemon"
)

// getOrStartDaemon detects existing daemon or starts a new one
// This is used by MCP proxy and other clients to ensure daemon is running
func getOrStartDaemon() (string, error) {
	// Try to get existing daemon info
	lock, err := daemon.GetDaemonInfo()
	if err == nil {
		log.Printf("✅ Using existing daemon: %s (PID: %d)", lock.DaemonURL, lock.PID)
		return lock.DaemonURL, nil
	}

	// No healthy daemon found, start new one
	log.Printf("🚀 No daemon found, starting new daemon...")

	// Start daemon in detached mode
	cmd := exec.Command("gorev", "daemon", "--detach", "--port", "5082")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to start daemon: %w", err)
	}

	// Wait for daemon to be ready
	daemonURL := "http://localhost:5082"
	if err := daemon.WaitForDaemon(daemonURL, 15*time.Second); err != nil {
		return "", fmt.Errorf("daemon failed to start: %w", err)
	}

	log.Printf("✅ Daemon started successfully: %s", daemonURL)
	return daemonURL, nil
}

// ensureDaemonRunning ensures daemon is running and returns its URL
// Similar to getOrStartDaemon but with more verbose logging
func ensureDaemonRunning(port string) (string, error) {
	// Check if daemon is already running
	lock, err := daemon.ReadLockFile()
	if err == nil && daemon.IsProcessRunning(lock.PID) {
		// Verify daemon is healthy
		if daemon.IsDaemonHealthy(lock.DaemonURL) {
			log.Printf("✅ Daemon already running: %s (PID: %d)", lock.DaemonURL, lock.PID)
			return lock.DaemonURL, nil
		}

		// Stale lock file
		log.Printf("⚠️  Stale lock file found, removing...")
		daemon.RemoveLockFile()
	}

	// Start new daemon
	log.Printf("🚀 Starting new daemon on port %s...", port)

	cmd := exec.Command("gorev", "daemon", "--detach", "--port", port)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to start daemon: %w", err)
	}

	daemonURL := fmt.Sprintf("http://localhost:%s", port)

	// Wait for daemon to be ready
	if err := daemon.WaitForDaemon(daemonURL, 15*time.Second); err != nil {
		return "", fmt.Errorf("daemon failed to become healthy: %w", err)
	}

	log.Printf("✅ Daemon started successfully: %s", daemonURL)
	return daemonURL, nil
}

// stopDaemon stops the running daemon gracefully
func stopDaemon() error {
	lock, err := daemon.ReadLockFile()
	if err != nil {
		return fmt.Errorf("no daemon running: %w", err)
	}

	log.Printf("🛑 Stopping daemon (PID: %d)...", lock.PID)

	// Use daemon-stop command
	cmd := exec.Command("gorev", "daemon-stop")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop daemon: %w", err)
	}

	return nil
}

// getDaemonStatus returns daemon status information
func getDaemonStatus() (*daemon.LockFile, error) {
	return daemon.GetDaemonInfo()
}
