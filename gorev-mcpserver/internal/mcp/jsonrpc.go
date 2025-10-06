package mcp

// JSON-RPC 2.0 protocol types for MCP proxy

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string                 `json:"jsonrpc"` // Must be "2.0"
	ID      interface{}            `json:"id"`      // Request ID (can be string, number, or null)
	Method  string                 `json:"method"`  // MCP tool name (e.g., "gorev_listele")
	Params  map[string]interface{} `json:"params"`  // Tool parameters
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"` // Must be "2.0"
	ID      interface{}   `json:"id"`      // Same ID as request
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC 2.0 error object
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes
const (
	// ParseError - Invalid JSON was received by the server
	ParseError = -32700

	// InvalidRequest - The JSON sent is not a valid Request object
	InvalidRequest = -32600

	// MethodNotFound - The method does not exist / is not available
	MethodNotFound = -32601

	// InvalidParams - Invalid method parameter(s)
	InvalidParams = -32602

	// InternalError - Internal JSON-RPC error
	InternalError = -32603

	// ServerError - Reserved for implementation-defined server-errors
	ServerError = -32000
)

// NewErrorResponse creates a JSON-RPC error response
func NewErrorResponse(id interface{}, code int, message string, data interface{}) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
}

// NewSuccessResponse creates a JSON-RPC success response
func NewSuccessResponse(id interface{}, result interface{}) JSONRPCResponse {
	return JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// WorkspaceContext represents workspace information for proxy
type WorkspaceContext struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}
