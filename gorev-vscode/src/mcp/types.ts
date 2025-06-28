// MCP Protocol Types

export interface MCPRequest {
  jsonrpc: '2.0';
  id: number;
  method: string;
  params?: any;
}

export interface MCPResponse {
  jsonrpc: '2.0';
  id: number;
  result?: any;
  error?: MCPErrorResponse;
}

export interface MCPErrorResponse {
  code: number;
  message: string;
  data?: any;
}

export interface MCPNotification {
  jsonrpc: '2.0';
  method: string;
  params?: any;
}

export interface MCPInitializeParams {
  protocolVersion: string;
  capabilities: MCPCapabilities;
  clientInfo?: {
    name: string;
    version: string;
  };
}

export interface MCPCapabilities {
  tools?: {
    listChanged?: boolean;
  };
  logging?: {
    levels?: string[];
  };
}

export interface MCPTool {
  name: string;
  description?: string;
  inputSchema?: any;
}

export interface MCPToolCallParams {
  name: string;
  arguments?: any;
}

export interface MCPToolResult {
  content: Array<{
    type: 'text';
    text: string;
  }>;
  isError?: boolean;
}