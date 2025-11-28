package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/msenol/gorev/internal/daemon"
)

// getOrStartDaemon detects existing daemon or starts a new one based on connection mode
// Connection mode can be: auto (default), local, docker, remote
// This is used by MCP proxy and other clients to ensure daemon is running
func getOrStartDaemon() (string, error) {
	// Check for existing daemon first (fast path)
	lock, err := daemon.GetDaemonInfo()
	if err == nil && daemon.IsDaemonHealthy(lock.DaemonURL) {
		log.Printf("✅ Using existing daemon: %s (PID: %d)", lock.DaemonURL, lock.PID)
		return lock.DaemonURL, nil
	}

	// Get connection mode from environment
	mode := getConnectionMode()

	switch mode {
	case "remote":
		return handleRemoteMode()
	case "docker":
		return handleDockerMode()
	case "local":
		return handleLocalMode()
	case "auto":
		fallthrough
	default:
		return handleAutoMode()
	}
}

// getConnectionMode returns the connection mode from environment variable
func getConnectionMode() string {
	// Check GOREV_MCP_CONNECTION_MODE env var first
	if mode := getEnvDefault("GOREV_MCP_CONNECTION_MODE", ""); mode != "" {
		return mode
	}
	// Default to auto
	return "auto"
}

// handleRemoteMode tries to connect to a remote daemon
func handleRemoteMode() (string, error) {
	host := getEnvDefault("GOREV_API_HOST", "localhost")
	port := getEnvDefault("GOREV_API_PORT", "5082")
	daemonURL := fmt.Sprintf("http://%s:%s", host, port)

	log.Printf("🔍 Trying remote daemon: %s", daemonURL)

	if daemon.IsDaemonHealthy(daemonURL) {
		log.Printf("✅ Remote daemon is healthy: %s", daemonURL)
		return daemonURL, nil
	}

	return "", fmt.Errorf("remote daemon at %s is not accessible", daemonURL)
}

// handleDockerMode starts docker container
func handleDockerMode() (string, error) {
	composeFile := getEnvDefault("GOREV_DOCKER_COMPOSE_FILE", "./docker-compose.yml")

	log.Printf("🐳 Starting Gorev with Docker (compose: %s)...", composeFile)

	cmd := exec.Command("docker-compose", "-f", composeFile, "up", "-d")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker-compose failed: %w\nOutput: %s", err, string(output))
	}

	daemonURL := "http://localhost:5082"
	if err := daemon.WaitForDaemon(daemonURL, 30*time.Second); err != nil {
		return "", fmt.Errorf("docker daemon failed to start: %w", err)
	}

	log.Printf("✅ Docker daemon started: %s", daemonURL)
	return daemonURL, nil
}

// handleLocalMode starts local gorev binary
func handleLocalMode() (string, error) {
	serverPath := getEnvDefault("GOREV_SERVER_PATH", "")
	if serverPath == "" {
		// Try to find in PATH
		serverPath = "gorev"
	}

	log.Printf("📦 Starting local daemon from: %s", serverPath)

	cmd := exec.Command(serverPath, "daemon", "--detach", "--port", "5082")
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to start local daemon: %w. Make sure gorev is installed: npm install -g gorev", err)
	}

	daemonURL := "http://localhost:5082"
	if err := daemon.WaitForDaemon(daemonURL, 15*time.Second); err != nil {
		return "", fmt.Errorf("local daemon failed to start: %w", err)
	}

	log.Printf("✅ Local daemon started: %s", daemonURL)
	return daemonURL, nil
}

// handleAutoMode tries auto-detection and best effort
func handleAutoMode() (string, error) {
	// 1. Try remote if specified
	if host := getEnvDefault("GOREV_API_HOST", ""); host != "" && host != "localhost" {
		return handleRemoteMode()
	}

	// 2. Try docker if docker-compose.yml exists
	if _, err := exec.LookPath("docker-compose"); err == nil {
		if _, err := os.Stat("./docker-compose.yml"); err == nil {
			log.Printf("📄 Found docker-compose.yml, trying Docker mode...")
			if daemonURL, err := handleDockerMode(); err == nil {
				return daemonURL, nil
			}
		}
	}

	// 3. Try local binary
	log.Printf("🔍 Trying local mode...")
	return handleLocalMode()
}

// getEnvDefault returns environment variable or default value
func getEnvDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
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
