package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/msenol/gorev/internal/api"
	"github.com/msenol/gorev/internal/daemon"
	"github.com/spf13/cobra"
)

func createDaemonCommand() *cobra.Command {
	var daemonPort string
	var detach bool

	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Run Gorev as background daemon process",
		Long: `Start Gorev daemon server that provides:
- HTTP REST API for workspace management
- WebSocket for real-time updates (future)
- MCP proxy connection endpoint
- Shared database connection pool

The daemon runs as a single process that multiple MCP clients can connect to,
eliminating port conflicts and reducing resource usage.`,
		Example: `  # Start daemon in foreground
  gorev daemon

  # Start daemon in background (detached)
  gorev daemon --detach

  # Start daemon on custom port
  gorev daemon --port 5083

  # Check daemon status
  curl http://localhost:5082/api/health`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if detach {
				return runDetachedDaemon(daemonPort)
			}
			return runDaemon(daemonPort)
		},
	}

	cmd.Flags().StringVar(&daemonPort, "port", "5082", "Daemon HTTP API port")
	cmd.Flags().BoolVar(&detach, "detach", false, "Run as background process (daemon)")

	return cmd
}

func runDaemon(port string) error {
	log.Printf("🚀 Starting Gorev Daemon on port %s...", port)

	// Create lock file
	if err := daemon.CreateLockFile(os.Getpid(), port, version); err != nil {
		return fmt.Errorf("failed to create lock file: %w", err)
	}
	defer func() {
		if err := daemon.RemoveLockFile(); err != nil {
			log.Printf("Warning: Failed to remove lock file: %v", err)
		}
	}()

	// Initialize workspace manager (multi-workspace support)
	workspaceManager := api.NewWorkspaceManager()

	// Set embedded migrations
	migrationsFS, err := getEmbeddedMigrationsFS()
	if err != nil {
		return fmt.Errorf("failed to get embedded migrations: %w", err)
	}
	workspaceManager.SetMigrationsFS(migrationsFS)

	// Create API server (pure multi-workspace, no legacy single workspace)
	apiServer := api.NewAPIServer(port, nil) // nil for legacy isYonetici

	// Serve static files (Web UI)
	if err := api.ServeStaticFiles(apiServer.App(), WebDistFS); err != nil {
		log.Printf("Warning: Failed to serve web UI: %v", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in background
	errChan := make(chan error, 1)
	go func() {
		if err := apiServer.Start(); err != nil {
			errChan <- fmt.Errorf("daemon server failed: %w", err)
		}
	}()

	log.Printf("✅ Gorev Daemon started successfully")
	log.Printf("📱 Web UI: http://localhost:%s", port)
	log.Printf("🔧 API: http://localhost:%s/api/v1", port)
	log.Printf("🔌 WebSocket: ws://localhost:%s/ws (future)", port)
	log.Printf("💾 Lock file: %s", daemon.GetLockFilePath())
	log.Printf("\nPress Ctrl+C to stop daemon")

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		log.Printf("\n🛑 Received signal %v, shutting down gracefully...", sig)
	case err := <-errChan:
		log.Printf("\n❌ Server error: %v", err)
		return err
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("daemon shutdown failed: %w", err)
	}

	log.Printf("✅ Daemon stopped successfully")
	return nil
}

func runDetachedDaemon(port string) error {
	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Fork process and run in background
	cmd := exec.Command(exePath, "daemon", "--port", port)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	// Start detached process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start daemon process: %w", err)
	}

	// Wait a moment for daemon to initialize
	time.Sleep(500 * time.Millisecond)

	// Verify daemon is running
	daemonURL := fmt.Sprintf("http://localhost:%s", port)
	if err := daemon.WaitForDaemon(daemonURL, 10*time.Second); err != nil {
		return fmt.Errorf("daemon failed to start: %w", err)
	}

	log.Printf("✅ Daemon started in background")
	log.Printf("📋 PID: %d", cmd.Process.Pid)
	log.Printf("🔗 URL: %s", daemonURL)
	log.Printf("💾 Lock file: %s", daemon.GetLockFilePath())
	log.Printf("\nUse 'curl %s/api/health' to check status", daemonURL)

	return nil
}

// stopDaemonCommand creates a command to stop running daemon
func createDaemonStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon-stop",
		Short: "Stop running Gorev daemon",
		Long:  "Gracefully stop the running Gorev daemon process",
		RunE: func(cmd *cobra.Command, args []string) error {
			lock, err := daemon.ReadLockFile()
			if err != nil {
				return fmt.Errorf("no daemon running: %w", err)
			}

			log.Printf("🛑 Stopping daemon (PID: %d)...", lock.PID)

			// Send SIGTERM for graceful shutdown
			process, err := os.FindProcess(lock.PID)
			if err != nil {
				return fmt.Errorf("failed to find daemon process: %w", err)
			}

			if err := process.Signal(syscall.SIGTERM); err != nil {
				return fmt.Errorf("failed to send stop signal: %w", err)
			}

			// Wait for process to exit
			for i := 0; i < 30; i++ {
				if !daemon.IsProcessRunning(lock.PID) {
					log.Printf("✅ Daemon stopped successfully")
					return nil
				}
				time.Sleep(1 * time.Second)
			}

			// Force kill if still running
			log.Printf("⚠️  Daemon did not stop gracefully, forcing kill...")
			if err := process.Kill(); err != nil {
				return fmt.Errorf("failed to kill daemon: %w", err)
			}

			log.Printf("✅ Daemon forcefully stopped")
			return nil
		},
	}
}

// statusDaemonCommand creates a command to check daemon status
func createDaemonStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "daemon-status",
		Short: "Check Gorev daemon status",
		Long:  "Display status and information about running Gorev daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			lock, err := daemon.GetDaemonInfo()
			if err != nil {
				fmt.Println("❌ Daemon Status: NOT RUNNING")
				fmt.Printf("   Reason: %v\n", err)
				return nil
			}

			fmt.Println("✅ Daemon Status: RUNNING")
			fmt.Printf("   PID: %d\n", lock.PID)
			fmt.Printf("   URL: %s\n", lock.DaemonURL)
			fmt.Printf("   Port: %s\n", lock.Port)
			fmt.Printf("   Version: %s\n", lock.Version)
			fmt.Printf("   Started: %s\n", lock.StartTime.Format(time.RFC3339))
			fmt.Printf("   Uptime: %s\n", time.Since(lock.StartTime).Round(time.Second))
			fmt.Printf("   Lock File: %s\n", daemon.GetLockFilePath())

			return nil
		},
	}
}
