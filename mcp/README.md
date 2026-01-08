# CMS Store MCP Handler

This package provides an MCP (Model Context Protocol) handler for the CMS Store, enabling seamless integration between LLM applications and the CMS data sources and tools. The handler supports MCP protocol and can be attached to any existing HTTP server.

> **Note:** This package only handles MCP protocol requests. For RESTful HTTP API functionality, please use the separate `rest` package.

## Features

- Page management (create, list, read, update, delete)
- Menu management (create, list, read)
- Site listing
- Schema discovery for LLMs (`cms_schema`)
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
	http.HandleFunc("/mcp/cms", mcpHandler.Handler)
	
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

This handler supports both MCP-standard JSON-RPC methods and legacy aliases:

- MCP-standard:
  - `initialize`
  - `notifications/initialized`
  - `tools/list`
  - `tools/call`
- Legacy aliases (supported for compatibility):
  - `list_tools` (alias of `tools/list`)
  - `call_tool` (alias of `tools/call`)

### Request/Response Shape

This MCP handler uses JSON-RPC 2.0 over HTTP `POST`.

#### List tools

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/list",
  "params": {}
}
```

#### Call a tool

MCP-standard:

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "tools/call",
  "params": {
    "name": "page_list",
    "arguments": {
      "limit": 10,
      "offset": 0
    }
  }
}
```

Legacy alias:

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "method": "call_tool",
  "params": {
    "tool_name": "page_list",
    "arguments": {
      "limit": 10,
      "offset": 0
    }
  }
}
```

Tool results are returned in the MCP format:

```json
{
  "jsonrpc": "2.0",
  "id": "1",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"items\":[]}" 
      }
    ]
  }
}
```

### Supported Tools

This MCP handler exposes the following tools via `tools/list`:

- `cms_schema`
- `page_list`
- `page_create`
- `page_get`
- `page_update`
- `page_delete`
- `menu_list`
- `menu_create`
- `menu_get`
- `site_list`

### Schema discovery (`cms_schema`)

LLMs can call `cms_schema` to retrieve a JSON document describing CMS entities and supported tool arguments.

Example call:

```json
{
  "jsonrpc": "2.0",
  "id": "schema",
  "method": "tools/call",
  "params": {
    "name": "cms_schema",
    "arguments": {}
  }
}
```

The response `result.content[0].text` contains a JSON document with `entities` and `tools` keys.

### ID typing: always use strings

CMS IDs can be very large. Some LLM clients may convert large integer-looking strings into JSON numbers (including scientific notation), which is lossy and can break lookups.

- Always send identifiers as strings:

```json
{ "id": "20260108160058473" }
```

- Do not send identifiers as numbers:

```json
{ "id": 2.0260108160058473e+31 }
```

The handler's `tools/list` response includes `inputSchema` definitions that mark identifiers as `type: "string"` to help clients send the correct types.

### Page Operations

##### Create a Page

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "create",
  "method": "tools/call",
  "params": {
    "name": "page_create",
    "arguments": {
      "title": "My New Page",
      "content": "<p>Page content goes here</p>",
      "status": "published",
      "site_id": "site-123"
    }
  }
}
```

**Response:**
```json
{
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"id\":\"page_123\",\"title\":\"My New Page\"}"
      }
    ]
  }
}
```

##### Get a Page

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "get",
  "method": "tools/call",
  "params": {
    "name": "page_get",
    "arguments": {
      "id": "page_123"
    }
  }
}
```

**Response:**
```json
{
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"id\":\"page_123\",\"title\":\"My New Page\"}"
      }
    ]
  }
}
```

##### List Pages

**Request:**

```json
{
  "jsonrpc": "2.0",
  "id": "list",
  "method": "tools/call",
  "params": {
    "name": "page_list",
    "arguments": {
      "limit": 10,
      "offset": 0,
      "site_id": "site-123"
    }
  }
}
```

#### Menu Operations

##### Create a Menu

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "create",
  "method": "tools/call",
  "params": {
    "name": "menu_create",
    "arguments": {
      "name": "Main Menu",
      "site_id": "site-123"
    }
  }
}
```

##### Get a Menu

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": "get",
  "method": "tools/call",
  "params": {
    "name": "menu_get",
    "arguments": {
      "id": "menu_123"
    }
  }
}
```

##### List Menus

```json
{
  "jsonrpc": "2.0",
  "id": "list",
  "method": "tools/call",
  "params": {
    "name": "menu_list",
    "arguments": {
      "limit": 10,
      "offset": 0,
      "site_id": "site-123"
    }
  }
}
```

### Site Operations

##### List Sites

```json
{
  "jsonrpc": "2.0",
  "id": "list",
  "method": "tools/call",
  "params": {
    "name": "site_list",
    "arguments": {
      "limit": 10,
      "offset": 0
    }
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

1. Add the tool definition to `handleToolsList` in `mcp.go`.

2. Add the tool dispatch in `dispatchTool` in `mcp.go`.

3. Implement the tool handler method on `*MCP` (for example `toolMyTool(ctx, args)`).

### Best Practices

- Follow consistent error handling patterns
- Add comprehensive documentation for new tools
- Write tests for all new functionality
- Ensure proper validation of input parameters
- Return well-structured responses

## License

This package is part of the CMS Store and is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).
