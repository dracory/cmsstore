# CMS Store MCP Handler

This package provides an MCP (Model Context Protocol) handler for the CMS Store, enabling seamless integration between LLM applications and the CMS data sources and tools. The handler supports MCP protocol and can be attached to any existing HTTP server.

> **Note:** This package only handles MCP protocol requests. For RESTful HTTP API functionality, please use the separate `rest` package.

## Features

- Page management (create, read, update, delete)
- Menu management (create, read)
- Extensible architecture for adding more handlers
- JSON-RPC 2.0 compatible
- Attachable to any existing HTTP server

## Getting Started

### Prerequisites

- Go 1.24 or later
- A running instance of the CMS Store

### Installation

```bash
go get github.com/dracory/cmsstore/mcp
```

### Basic Usage

```go
package main

import (
	"log"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/mcp"
)

func main() {
	// Initialize your CMS store
	store := cmsstore.NewStore(db) // Your database connection

	// Create the MCP handler
	mcpHandler := mcp.NewMCP(store)
	
	// Register the MCP handler with your existing router/mux
	http.HandleFunc("/mcp/", mcpHandler.Handler)
	
	// Start your server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

## API Reference

### MCP Protocol API

The MCP (Model Context Protocol) API allows LLMs to interact with the CMS Store using a structured protocol designed for AI applications.

#### Page Operations

##### Create a Page

**Request:**
```json
{
  "name": "page_create",
  "params": {
    "title": "My New Page",
    "content": "<p>Page content goes here</p>",
    "status": "published",
    "site_id": "site-123"
  }
}
```

**Response:**
```json
{
  "result": {
    "id": "page_123",
    "title": "My New Page",
    "content": "<p>Page content goes here</p>",
    "status": "published",
    "site_id": "site-123",
    "success": true
  }
}
```

##### Get a Page

**Request:**
```json
{
  "name": "page_get",
  "params": {
    "id": "page_123"
  }
}
```

**Response:**
```json
{
  "result": {
    "id": "page_123",
    "title": "My New Page",
    "content": "<p>Page content goes here</p>",
    "status": "published",
    "site_id": "site-123"
  }
}
```

#### Menu Operations

##### Create a Menu

**Request:**
```json
{
  "name": "menu_create",
  "params": {
    "title": "Main Menu",
    "site_id": "site-123"
  }
}
```

##### Get a Menu

**Request:**
```json
{
  "name": "menu_get",
  "params": {
    "id": "menu_123"
  }
}
```

## Error Handling

### MCP Protocol Errors

Errors are returned in the MCP format:

```json
{
  "error": {
    "message": "Error message here"
  }
}
```

## Extending the Server

The MCP server is designed to be easily extensible for the MCP protocol interface.

### Adding New MCP Tools

1. Define a new tool in the `registerHandlers` method in `server.go`:

```go
newTool := mcp.NewTool("tool_name",
    mcp.WithDescription("Tool description"),
    mcp.WithString("param_name", mcp.Required(), mcp.Description("Parameter description")),
    // Add more parameters as needed
)
s.server.AddTool(newTool, s.handleNewTool)
```

2. Implement the handler function:

```go
func (s *Server) handleNewTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Parse parameters
    // Perform operations
    // Return results
    return mcp.NewToolResultSuccess(map[string]interface{}{
        "success": true,
        "data": result,
    }), nil
}
```

### Best Practices

- Follow consistent error handling patterns
- Add comprehensive documentation for new tools
- Write tests for all new functionality
- Ensure proper validation of input parameters
- Return well-structured responses

## License

This package is part of the CMS Store and is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).
