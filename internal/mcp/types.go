package mcp

import "encoding/json"

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool represents an MCP tool
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// TextContent represents text content in MCP
type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// InitializeRequest represents the initialize request
type InitializeRequest struct {
	ProtocolVersion string            `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo        `json:"clientInfo"`
}

// InitializeResponse represents the initialize response
type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

// ClientCapabilities represents client capabilities
type ClientCapabilities struct {
	Roots *RootsCapability `json:"roots,omitempty"`
}

// RootsCapability represents roots capability
type RootsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ServerCapabilities represents server capabilities
type ServerCapabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

// ToolsCapability represents tools capability
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ClientInfo represents client information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerInfo represents server information
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ListToolsResponse represents the response to list_tools
type ListToolsResponse struct {
	Tools []Tool `json:"tools"`
}

// CallToolRequest represents a tool call request
type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// CallToolResponse represents a tool call response
type CallToolResponse struct {
	Content []TextContent `json:"content"`
}

// ListRootsResponse represents the response to list_roots
type ListRootsResponse struct {
	Roots []Root `json:"roots"`
}

// Root represents a root directory
type Root struct {
	URI  string `json:"uri"`
	Name string `json:"name,omitempty"`
}

// Constants for JSON-RPC
const (
	JSONRPCVersion = "2.0"
)

// MCP method names
const (
	MethodInitialize = "initialize"
	MethodListTools  = "tools/list"
	MethodCallTool   = "tools/call"
	MethodListRoots  = "roots/list"
)
