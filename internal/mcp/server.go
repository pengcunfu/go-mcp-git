package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// Server represents an MCP server
type Server struct {
	name         string
	version      string
	capabilities ServerCapabilities
	tools        []Tool
	toolHandlers map[string]ToolHandler
	initialized  bool
}

// ToolHandler is a function that handles tool calls
type ToolHandler func(ctx context.Context, arguments map[string]interface{}) ([]TextContent, error)

// NewServer creates a new MCP server
func NewServer(name, version string) *Server {
	return &Server{
		name:         name,
		version:      version,
		capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
		},
		tools:        make([]Tool, 0),
		toolHandlers: make(map[string]ToolHandler),
		initialized:  false,
	}
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	s.tools = append(s.tools, tool)
	s.toolHandlers[tool.Name] = handler
}

// Serve starts the MCP server using stdio
func (s *Server) Serve(ctx context.Context) error {
	reader := bufio.NewReader(os.Stdin)
	writer := os.Stdout

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read request
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("failed to read request: %w", err)
			}

			// Process request
			response, err := s.handleRequest(ctx, line)
			if err != nil {
				log.Printf("Error handling request: %v", err)
				continue
			}

			// Write response
			if response != nil {
				responseBytes, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error marshaling response: %v", err)
					continue
				}

				if _, err := writer.Write(append(responseBytes, '\n')); err != nil {
					log.Printf("Error writing response: %v", err)
					continue
				}
			}
		}
	}
}

// handleRequest processes a single JSON-RPC request
func (s *Server) handleRequest(ctx context.Context, requestBytes []byte) (*JSONRPCResponse, error) {
	var request JSONRPCRequest
	if err := json.Unmarshal(requestBytes, &request); err != nil {
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			Error: &RPCError{
				Code:    -32700,
				Message: "Parse error",
			},
		}, nil
	}

	switch request.Method {
	case MethodInitialize:
		return s.handleInitialize(ctx, request)
	case MethodListTools:
		return s.handleListTools(ctx, request)
	case MethodCallTool:
		return s.handleCallTool(ctx, request)
	default:
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      request.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}, nil
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(ctx context.Context, request JSONRPCRequest) (*JSONRPCResponse, error) {
	var initReq InitializeRequest
	if err := json.Unmarshal(request.Params, &initReq); err != nil {
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      request.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}, nil
	}

	s.initialized = true

	response := InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo: ServerInfo{
			Name:    s.name,
			Version: s.version,
		},
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      request.ID,
		Result:  response,
	}, nil
}

// handleListTools handles the list_tools request
func (s *Server) handleListTools(ctx context.Context, request JSONRPCRequest) (*JSONRPCResponse, error) {
	if !s.initialized {
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      request.ID,
			Error: &RPCError{
				Code:    -32002,
				Message: "Server not initialized",
			},
		}, nil
	}

	response := ListToolsResponse{
		Tools: s.tools,
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      request.ID,
		Result:  response,
	}, nil
}

// handleCallTool handles the call_tool request
func (s *Server) handleCallTool(ctx context.Context, request JSONRPCRequest) (*JSONRPCResponse, error) {
	if !s.initialized {
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      request.ID,
			Error: &RPCError{
				Code:    -32002,
				Message: "Server not initialized",
			},
		}, nil
	}

	var callReq CallToolRequest
	if err := json.Unmarshal(request.Params, &callReq); err != nil {
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      request.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}, nil
	}

	handler, exists := s.toolHandlers[callReq.Name]
	if !exists {
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      request.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: fmt.Sprintf("Unknown tool: %s", callReq.Name),
			},
		}, nil
	}

	content, err := handler(ctx, callReq.Arguments)
	if err != nil {
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      request.ID,
			Error: &RPCError{
				Code:    -32603,
				Message: fmt.Sprintf("Tool execution error: %v", err),
			},
		}, nil
	}

	response := CallToolResponse{
		Content: content,
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      request.ID,
		Result:  response,
	}, nil
}
