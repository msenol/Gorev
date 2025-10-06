package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/msenol/gorev/internal/mcp"
	"github.com/spf13/cobra"
)

func createMCPProxyCommand() *cobra.Command {
	var debugMode bool

	cmd := &cobra.Command{
		Use:   "mcp-proxy",
		Short: "MCP proxy to daemon server (for AI assistants)",
		Long: `Lightweight MCP proxy that forwards stdio MCP messages to daemon HTTP API.

This command is used by MCP clients (VS Code, Claude Code, Cursor, etc.) to connect
to the Gorev daemon. It automatically:
- Detects or starts the daemon
- Auto-detects the current workspace
- Registers workspace with daemon
- Forwards MCP tool calls to daemon HTTP API

The proxy runs in stdio mode, making it compatible with all MCP-compatible editors.`,
		Example: `  # Start MCP proxy (auto-detect workspace)
  gorev mcp-proxy

  # Start MCP proxy with debug logging
  gorev mcp-proxy --debug

  # Usage in MCP config (.kilocode/mcp.json):
  {
    "mcpServers": {
      "gorev": {
        "command": "npx",
        "args": ["-y", "@mehmetsenol/gorev-mcp-server", "mcp-proxy"]
      }
    }
  }`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPProxy(debugMode)
		},
	}

	cmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug logging (logs to stderr)")

	return cmd
}

func runMCPProxy(debug bool) error {
	// 1. Detect or start daemon
	if debug {
		log.SetOutput(os.Stderr) // Send logs to stderr so stdout is clean for MCP protocol
		log.Printf("[MCP Proxy] Starting MCP proxy...")
	}

	daemonURL, err := getOrStartDaemon()
	if err != nil {
		return fmt.Errorf("failed to get daemon: %w", err)
	}

	if debug {
		log.Printf("[MCP Proxy] Daemon URL: %s", daemonURL)
	}

	// 2. Auto-detect workspace
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	workspacePath := findWorkspaceRoot(cwd)
	if debug {
		log.Printf("[MCP Proxy] Workspace path: %s", workspacePath)
	}

	// 3. Register workspace with daemon
	workspaceCtx, err := registerWorkspaceWithDaemon(daemonURL, workspacePath, debug)
	if err != nil {
		return fmt.Errorf("failed to register workspace: %w", err)
	}

	if debug {
		log.Printf("[MCP Proxy] Workspace registered:")
		log.Printf("  ID: %s", workspaceCtx.ID)
		log.Printf("  Name: %s", workspaceCtx.Name)
		log.Printf("  Path: %s", workspaceCtx.Path)
		log.Printf("[MCP Proxy] Ready to forward MCP messages")
	}

	// 4. Start MCP proxy (stdio ↔ HTTP bridge)
	proxy := mcp.NewProxy(daemonURL, workspaceCtx, debug)
	return proxy.Serve() // Blocking - reads from stdin, writes to stdout
}

// findWorkspaceRoot walks up directory tree to find .gorev/ directory
func findWorkspaceRoot(startPath string) string {
	current := startPath

	for {
		gorevDir := filepath.Join(current, ".gorev")
		if _, err := os.Stat(gorevDir); err == nil {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root, fallback to home
			home, err := os.UserHomeDir()
			if err != nil {
				return startPath // Last resort
			}
			return home
		}

		current = parent
	}
}

// registerWorkspaceWithDaemon registers workspace with daemon and returns workspace context
func registerWorkspaceWithDaemon(daemonURL, workspacePath string, debug bool) (*mcp.WorkspaceContext, error) {
	body := map[string]string{
		"path": workspacePath,
		"name": filepath.Base(workspacePath),
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if debug {
		log.Printf("[MCP Proxy] Registering workspace: %s", filepath.Base(workspacePath))
	}

	resp, err := http.Post(
		daemonURL+"/api/v1/workspaces/register",
		"application/json",
		bytes.NewReader(data),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to register workspace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("workspace registration failed with status %d", resp.StatusCode)
	}

	var result struct {
		WorkspaceID string `json:"workspace_id"`
		Workspace   struct {
			Name string `json:"name"`
			Path string `json:"path"`
		} `json:"workspace"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &mcp.WorkspaceContext{
		ID:   result.WorkspaceID,
		Name: result.Workspace.Name,
		Path: result.Workspace.Path,
	}, nil
}
