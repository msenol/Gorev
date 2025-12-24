package mcp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/msenol/gorev/internal/daemon"
)

// Proxy forwards MCP messages between stdio and HTTP daemon
type Proxy struct {
	daemonURL    string
	workspaceCtx *WorkspaceContext
	client       *http.Client
	debug        bool
	clientID     string
	stopChan     chan struct{}
}

// NewProxy creates a new MCP proxy instance
func NewProxy(daemonURL string, ctx *WorkspaceContext, debug bool) *Proxy {
	return &Proxy{
		daemonURL:    daemonURL,
		workspaceCtx: ctx,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		debug: debug,
	}
}

// Serve starts the MCP proxy, reading from stdin and writing to stdout
func (p *Proxy) Serve() error {
	if p.debug {
		log.Printf("[MCP Proxy] Starting proxy for workspace: %s (%s)", p.workspaceCtx.Name, p.workspaceCtx.ID)
		log.Printf("[MCP Proxy] Daemon URL: %s", p.daemonURL)
	}

	// Generate client ID and register with daemon
	p.clientID = daemon.GenerateClientID()
	p.stopChan = make(chan struct{})

	if err := p.registerClient(); err != nil {
		log.Printf("[MCP Proxy] Warning: Failed to register client: %v", err)
	} else {
		log.Printf("[MCP Proxy] Client registered: %s", p.clientID)
	}

	// Start heartbeat loop
	go p.heartbeatLoop()

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()

		if p.debug {
			log.Printf("[MCP Proxy] <- stdin: %s", line)
		}

		// Parse JSON-RPC request
		var req JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			if p.debug {
				log.Printf("[MCP Proxy] Parse error: %v", err)
			}
			p.writeError(writer, nil, ParseError, "Parse error", err.Error())
			continue
		}

		if p.debug {
			log.Printf("[MCP Proxy] Request: method=%s, id=%v (type=%T)", req.Method, req.ID, req.ID)
		}

		// Check if this is a notification (id is null) - don't send response
		if req.ID == nil {
			if p.debug {
				log.Printf("[MCP Proxy] Notification (no response expected): %s", req.Method)
			}
			continue
		}

		// Forward to daemon HTTP API
		response, err := p.forwardToDaemon(req)
		if err != nil {
			if p.debug {
				log.Printf("[MCP Proxy] Forward error: %v", err)
			}
			p.writeError(writer, req.ID, InternalError, "Internal error", err.Error())
			continue
		}

		// Write response to stdout
		if p.debug {
			log.Printf("[MCP Proxy] -> stdout: %s", response)
		}

		writer.WriteString(response + "\n")
		writer.Flush()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// forwardToDaemon forwards JSON-RPC request to daemon HTTP API
func (p *Proxy) forwardToDaemon(req JSONRPCRequest) (string, error) {
	// Map MCP method to HTTP endpoint
	endpoint := p.mapMethodToEndpoint(req.Method)

	if p.debug {
		log.Printf("[MCP Proxy] Forwarding %s to %s (id=%v)", req.Method, endpoint, req.ID)
	}

	// Create HTTP request
	body, err := json.Marshal(req.Params)
	if err != nil {
		return "", fmt.Errorf("failed to marshal params: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", p.daemonURL, endpoint),
		bytes.NewReader(body),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Inject workspace headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Workspace-Id", p.workspaceCtx.ID)
	httpReq.Header.Set("X-Workspace-Path", p.workspaceCtx.Path)
	httpReq.Header.Set("X-Workspace-Name", p.workspaceCtx.Name)

	// Execute request
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Convert HTTP response to JSON-RPC response
	id := req.ID
	if p.debug {
		log.Printf("[MCP Proxy] Converting response with id=%v (type=%T)", id, id)
	}
	return p.convertHTTPToJSONRPC(resp, id)
}

// mapMethodToEndpoint maps MCP tool names to HTTP API endpoints
func (p *Proxy) mapMethodToEndpoint(method string) string {
	// Map of MCP tool names to REST endpoints
	// For now, use a generic MCP bridge endpoint that handles all tools
	return "/api/v1/mcp/" + method
}

// convertHTTPToJSONRPC converts HTTP response to JSON-RPC response
func (p *Proxy) convertHTTPToJSONRPC(resp *http.Response, id interface{}) (string, error) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response body as generic JSON
	var result interface{}
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return "", fmt.Errorf("failed to parse response: %w", err)
		}
	}

	if p.debug {
		log.Printf("[MCP Proxy] HTTP response: status=%d, id=%v (type=%T), result_type=%T",
			resp.StatusCode, id, id, result)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		// HTTP error, convert to JSON-RPC error
		errorResp := NewErrorResponse(id, ServerError, fmt.Sprintf("HTTP %d", resp.StatusCode), result)
		data, _ := json.Marshal(errorResp)
		return string(data), nil
	}

	// Success response
	response := NewSuccessResponse(id, result)
	data, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(data), nil
}

// writeError writes JSON-RPC error to stdout
func (p *Proxy) writeError(w *bufio.Writer, id interface{}, code int, message string, data interface{}) {
	response := NewErrorResponse(id, code, message, data)

	jsonData, _ := json.Marshal(response)
	w.WriteString(string(jsonData) + "\n")
	w.Flush()
}

// registerClient registers this MCP proxy as an active client
func (p *Proxy) registerClient() error {
	reqBody := map[string]interface{}{
		"client_id":    p.clientID,
		"client_type":  "mcp-proxy",
		"workspace_id": p.workspaceCtx.ID,
		"ttl_seconds":  300, // 5 minutes TTL
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/daemon/clients/register", p.daemonURL),
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status %d", resp.StatusCode)
	}

	return nil
}

// unregisterClient removes this client from tracking
func (p *Proxy) unregisterClient() error {
	reqBody := map[string]interface{}{
		"client_id": p.clientID,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/daemon/clients/unregister", p.daemonURL),
		bytes.NewReader(data),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// heartbeatLoop periodically sends heartbeats to keep client active
func (p *Proxy) heartbeatLoop() {
	ticker := time.NewTicker(60 * time.Second) // Heartbeat every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := p.sendHeartbeat(); err != nil {
				log.Printf("[MCP Proxy] Heartbeat failed: %v", err)
			}
		case <-p.stopChan:
			return
		}
	}
}

// sendHeartbeat sends a heartbeat to extend TTL
func (p *Proxy) sendHeartbeat() error {
	reqBody := map[string]interface{}{
		"client_id":   p.clientID,
		"ttl_seconds": 300,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/daemon/heartbeat", p.daemonURL),
		bytes.NewReader(data),
	)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
