package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/mark3labs/mcp-go/mcp"
	mcpServer "github.com/mark3labs/mcp-go/server"
)

// MCPInterface defines the interface for MCP handlers
type MCPInterface interface {
	Handler() http.HandlerFunc
}

// MCP represents the MCP handler for CMS operations
type MCP struct {
	store  cmsstore.StoreInterface
	server *mcpServer.MCPServer
	tools  map[string]mcp.Tool
}

// NewMCP creates a new MCP handler instance
func NewMCP(store cmsstore.StoreInterface) *MCP {
	handler := &MCP{
		store: store,
		tools: make(map[string]mcp.Tool),
	}

	// Initialize MCP server
	handler.server = mcpServer.NewMCPServer(
		"CMS Store",
		"1.0.0",
	)

	// Register handlers
	handler.registerHandlers()

	return handler
}

// Handler returns an http.HandlerFunc that can be attached to any router
func (m *MCP) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Process the MCP protocol request
		if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
			// Handle the MCP request
			m.handleMCPRequest(w, r)
			return
		}

		// Return error for non-MCP requests
		http.Error(w, `{"success":false,"error":"This endpoint only accepts MCP protocol requests"}`, http.StatusBadRequest)
	}
}

// handleMCPRequest processes an MCP protocol request
func (m *MCP) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"success":false,"error":"Failed to read request body: %v"}`, err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON-RPC request
	var request struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      string          `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      nil,
			"error": map[string]interface{}{
				"code":    -32700,
				"message": "Parse error",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle the call_tool method
	if request.Method == "call_tool" {
		m.handleCallTool(w, r.Context(), request.ID, request.Params)
		return
	}

	// For other methods, let the MCP server handle it
	response := m.server.HandleMessage(r.Context(), body)

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// registerHandlers registers all the MCP tools and their handlers
func (m *MCP) registerHandlers() {
	// Register page operations
	pageCreateTool := mcp.NewTool("page_create",
		mcp.WithDescription("Create a new page"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Page content")),
		mcp.WithString("status", mcp.Description("Page status (draft, published, etc.)"), mcp.Enum("draft", "published")),
	)
	m.tools["page_create"] = pageCreateTool
	m.server.AddTool(pageCreateTool, m.handlePageCreate)

	// Add more tools for other operations (get, update, delete, etc.)
	pageGetTool := mcp.NewTool("page_get",
		mcp.WithDescription("Get a page by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Page ID")),
	)
	m.tools["page_get"] = pageGetTool
	m.server.AddTool(pageGetTool, m.handlePageGet)

	pageUpdateTool := mcp.NewTool("page_update",
		mcp.WithDescription("Update a page"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Page ID")),
		mcp.WithObject("updates", mcp.Required(), mcp.Description("Page updates")),
	)
	m.tools["page_update"] = pageUpdateTool
	m.server.AddTool(pageUpdateTool, m.handlePageUpdate)

	pageDeleteTool := mcp.NewTool("page_delete",
		mcp.WithDescription("Delete a page"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Page ID")),
	)
	m.tools["page_delete"] = pageDeleteTool
	m.server.AddTool(pageDeleteTool, m.handlePageDelete)

	// Menu operations
	menuCreateTool := mcp.NewTool("menu_create",
		mcp.WithDescription("Create a new menu"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Menu name")),
		mcp.WithString("status", mcp.Description("Menu status (draft, active, inactive)"), 
			mcp.Enum(cmsstore.MENU_STATUS_DRAFT, cmsstore.MENU_STATUS_ACTIVE, cmsstore.MENU_STATUS_INACTIVE)),
		// Note: The MCP library API doesn't support nested objects in the way we were trying to use them
		// So we'll use a simpler approach for menu items
		mcp.WithString("items_json", mcp.Description("Menu items as JSON string")),
	)
	m.tools["menu_create"] = menuCreateTool
	m.server.AddTool(menuCreateTool, m.handleMenuCreate)

	menuGetTool := mcp.NewTool("menu_get",
		mcp.WithDescription("Get a menu"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Menu ID")),
	)
	m.tools["menu_get"] = menuGetTool
	m.server.AddTool(menuGetTool, m.handleMenuGet)

	// We need to handle call_tool method manually
}

// handlePageCreate handles the page_create tool
// handlePageCreate handles the page_create tool
func (m *MCP) handlePageCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse the request parameters
	var params struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  string `json:"status,omitempty"`
	}

	// Convert request.Arguments to JSON and then unmarshal into our struct
	argsJSON, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &params); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	// Create page using the store
	page := cmsstore.NewPage()
	page.SetTitle(params.Title)
	page.SetContent(params.Content)
	if params.Status != "" {
		page.SetStatus(params.Status)
	}

	// Save the page
	if err := m.store.PageCreate(ctx, page); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create page: %v", err)), nil
	}

	// Return success response
	result, err := json.Marshal(map[string]interface{}{
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
		"success": true,
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal page: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

// handlePageGet handles page retrieval requests
func (m *MCP) handlePageGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Get the page ID from request parameters
	pageID, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing or invalid page ID: %v", err)), nil
	}

	// Get the page from the store
	page, err := m.store.PageFindByID(ctx, pageID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get page: %v", err)), nil
	}

	if page == nil {
		return mcp.NewToolResultError("page not found"), nil
	}

	// Convert the page to JSON for the response
	result, err := json.Marshal(map[string]interface{}{
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal page: %v", err)), nil
	}

	// Return the page data as a text result
	return mcp.NewToolResultText(string(result)), nil
}

// handlePageUpdate handles page update requests
func (m *MCP) handlePageUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse the arguments
	var args struct {
		ID      string         `json:"id"`
		Updates map[string]any `json:"updates"`
	}

	// Convert request.Arguments to JSON and then unmarshal into our struct
	argsJSON, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &args); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	// Get the existing page
	page, err := m.store.PageFindByID(ctx, args.ID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to find page: %v", err)), nil
	}

	if page == nil {
		return mcp.NewToolResultError("page not found"), nil
	}

	// Apply updates
	for key, value := range args.Updates {
		switch key {
		case "title":
			if title, ok := value.(string); ok {
				page.SetTitle(title)
			}
		case "content":
			if content, ok := value.(string); ok {
				page.SetContent(content)
			}
		case "status":
			if status, ok := value.(string); ok {
				page.SetStatus(status)
			}
		}
	}

	// Save the updated page
	if err := m.store.PageUpdate(ctx, page); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to update page: %v", err)), nil
	}

	// Return success response
	result, err := json.Marshal(map[string]interface{}{
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
		"success": true,
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal page: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

// handlePageDelete handles page deletion requests
func (m *MCP) handlePageDelete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse the arguments
	var data struct {
		ID string `json:"id"`
	}

	// Convert request.Arguments to JSON and then unmarshal into our struct
	argsJSON, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &data); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	// First try to get the page to ensure it exists
	page, err := m.store.PageFindByID(ctx, data.ID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to find page: %v", err)), nil
	}
	if page == nil {
		return mcp.NewToolResultError("page not found"), nil
	}

	// Try soft delete first, fall back to hard delete if needed
	err = m.store.PageSoftDeleteByID(ctx, data.ID)
	if err != nil {
		// If soft delete fails, try hard delete
		if err := m.store.PageDelete(ctx, page); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to delete page: %v", err)), nil
		}
	}

	// Return success response
	result, err := json.Marshal(map[string]interface{}{
		"id":      data.ID,
		"success": true,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal response: %v", err)), nil
	}
	return mcp.NewToolResultText(string(result)), nil
}

// handleMenuCreate handles menu creation requests
func (m *MCP) handleMenuCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse the request parameters
	var params struct {
		Name   string `json:"name"`
		Status string `json:"status,omitempty"`
		Items  []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"items,omitempty"`
	}

	// Convert request.Arguments to JSON and then unmarshal into our struct
	argsJSON, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &params); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	// Create menu using the store
	menu := cmsstore.NewMenu()
	menu.SetName(params.Name)
	if params.Status != "" {
		menu.SetStatus(params.Status)
	}

	// Add menu items if provided
	if len(params.Items) > 0 {
		// Convert items to JSON string
		itemsJSON, err := json.Marshal(params.Items)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal menu items: %v", err)), nil
		}
		// Store items as meta since SetItems is not available
		err = menu.SetMeta("items", string(itemsJSON))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to set menu items: %v", err)), nil
		}
	}

	// Save the menu
	if err := m.store.MenuCreate(ctx, menu); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create menu: %v", err)), nil
	}

	// Return success response
	// Get items from meta
	itemsJSON := menu.Meta("items")
	
	result, err := json.Marshal(map[string]interface{}{
		"id":      menu.ID(),
		"name":    menu.Name(),
		"status":  menu.Status(),
		"items":   itemsJSON,
		"success": true,
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal menu: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

// handleCallTool processes the JSON-RPC call_tool method
func (m *MCP) handleCallTool(w http.ResponseWriter, ctx context.Context, id string, params json.RawMessage) {
	// Parse the call_tool parameters
	var callParams struct {
		ToolName  string          `json:"tool_name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(params, &callParams); err != nil {
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"error": map[string]interface{}{
				"code":    -32602,
				"message": "Invalid params",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if the tool exists
	_, ok := m.tools[callParams.ToolName]
	if !ok {
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"error": map[string]interface{}{
				"code":    -32601,
				"message": fmt.Sprintf("Tool %s not found", callParams.ToolName),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create a CallToolRequest with the correct structure
	callRequest := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      callParams.ToolName,
			Arguments: json.RawMessage(callParams.Arguments),
		},
	}

	// Call the appropriate handler based on the tool name
	var result *mcp.CallToolResult
	var err error

	switch callParams.ToolName {
	case "page_create":
		result, err = m.handlePageCreate(ctx, callRequest)
	case "page_get":
		result, err = m.handlePageGet(ctx, callRequest)
	case "page_update":
		result, err = m.handlePageUpdate(ctx, callRequest)
	case "page_delete":
		result, err = m.handlePageDelete(ctx, callRequest)
	case "menu_create":
		result, err = m.handleMenuCreate(ctx, callRequest)
	case "menu_get":
		result, err = m.handleMenuGet(ctx, callRequest)
	default:
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"error": map[string]interface{}{
				"code":    -32601,
				"message": fmt.Sprintf("Handler for tool %s not implemented", callParams.ToolName),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if err != nil {
		response := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"error": map[string]interface{}{
				"code":    -32603,
				"message": err.Error(),
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return the successful result
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result.Result,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMenuGet handles menu retrieval requests
func (m *MCP) handleMenuGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse the request parameters
	var params struct {
		ID string `json:"id"`
	}

	// Convert request.Arguments to JSON and then unmarshal into our struct
	argsJSON, err := json.Marshal(request.Params.Arguments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	if err := json.Unmarshal(argsJSON, &params); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse request: %v", err)), nil
	}

	// Get the menu from the store
	menu, err := m.store.MenuFindByID(ctx, params.ID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to find menu: %v", err)), nil
	}

	if menu == nil {
		return mcp.NewToolResultError("menu not found"), nil
	}

	// Convert the menu to JSON for the response
	// Get items from meta
	itemsJSON := menu.Meta("items")
	
	result, err := json.Marshal(map[string]interface{}{
		"id":     menu.ID(),
		"name":   menu.Name(),
		"status": menu.Status(),
		"items":   itemsJSON,
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal menu: %v", err)), nil
	}

	// Return the menu data as a text result
	return mcp.NewToolResultText(string(result)), nil
}
