package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/mark3labs/mcp-go/mcp"
	mcpServer "github.com/mark3labs/mcp-go/server"
)

// MCPInterface defines the interface for the MCP handler
type MCPInterface interface {
	// Handler returns the HTTP handler for the MCP server
	Handler(w http.ResponseWriter, r *http.Request)
}

// mcpHandler represents the MCP handler for CMS operations
type mcpHandler struct {
	store  cmsstore.StoreInterface
	server *mcpServer.MCPServer
}

// NewMCP creates a new MCP handler instance
func NewMCP(store cmsstore.StoreInterface) MCPInterface {
	handler := &mcpHandler{
		store: store,
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

// Handler is the main handler for the MCP server.
// It processes MCP protocol requests and can be attached to any existing HTTP server.
func (m *mcpHandler) Handler(w http.ResponseWriter, r *http.Request) {
	// Process the MCP protocol request
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		// Use the MCP server's handler function directly
		m.server.HandleMCPRequest(w, r)
		return
	}

	// Return error for non-MCP requests
	http.Error(w, `{"success":false,"error":"This endpoint only accepts MCP protocol requests"}`, http.StatusBadRequest)
}

// registerHandlers registers all the MCP handlers
func (m *mcpHandler) registerHandlers() {
	// Register page operations
	pageCreateTool := mcp.NewTool("page_create",
		mcp.WithDescription("Create a new page"),
		mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Page content")),
		mcp.WithString("status", mcp.Description("Page status (draft, published, etc.)"), mcp.Enum("draft", "published")),
	)
	m.server.AddTool(pageCreateTool, m.handlePageCreate)

	// Add more tools for other operations (get, update, delete, etc.)
	pageGetTool := mcp.NewTool("page_get",
		mcp.WithDescription("Get a page by ID"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Page ID")),
	)
	m.server.AddTool(pageGetTool, m.handlePageGet)

	pageUpdateTool := mcp.NewTool("page_update",
		mcp.WithDescription("Update a page"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Page ID")),
		mcp.WithObject("updates", mcp.Required(), mcp.Description("Page updates")),
	)
	m.server.AddTool(pageUpdateTool, m.handlePageUpdate)

	pageDeleteTool := mcp.NewTool("page_delete",
		mcp.WithDescription("Delete a page"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Page ID")),
	)
	m.server.AddTool(pageDeleteTool, m.handlePageDelete)

	// Menu operations
	menuCreateTool := mcp.NewTool("menu_create",
		mcp.WithDescription("Create a new menu"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Menu name")),
		mcp.WithString("status", mcp.Description("Menu status (draft, published, etc.)"), mcp.Enum("draft", "published")),
		mcp.WithArray("items", mcp.Description("Menu items"), mcp.WithObject(mcp.WithString("title", mcp.Required(), mcp.Description("Item title")), mcp.WithString("url", mcp.Required(), mcp.Description("Item URL")))),
	)
	m.server.AddTool(menuCreateTool, m.handleMenuCreate)

	menuGetTool := mcp.NewTool("menu_get",
		mcp.WithDescription("Get a menu"),
		mcp.WithString("id", mcp.Required(), mcp.Description("Menu ID")),
	)
	m.server.AddTool(menuGetTool, m.handleMenuGet)
}

// handlePageCreate handles the page_create tool
func (m *mcpHandler) handlePageCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (m *mcpHandler) handlePageGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (m *mcpHandler) handlePageUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (m *mcpHandler) handlePageDelete(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
func (m *mcpHandler) handleMenuCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		items := []map[string]interface{}{}
		for _, item := range params.Items {
			items = append(items, map[string]interface{}{
				"title": item.Title,
				"url":   item.URL,
			})
		}
		// Store items as metadata since SetItems is not available
		menu.SetMetadata(map[string]interface{}{"items": items})
	}

	// Save the menu
	if err := m.store.MenuCreate(ctx, menu); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create menu: %v", err)), nil
	}

	// Return success response
	result, err := json.Marshal(map[string]interface{}{
		"id":      menu.ID(),
		"name":    menu.Name(),
		"status":  menu.Status(),
		"items":   menu.GetMetadata()["items"],
		"success": true,
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal menu: %v", err)), nil
	}

	return mcp.NewToolResultText(string(result)), nil
}

// handleMenuGet handles menu retrieval requests
func (m *mcpHandler) handleMenuGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	result, err := json.Marshal(map[string]interface{}{
		"id":     menu.ID(),
		"name":   menu.Name(),
		"status": menu.Status(),
		"items":  menu.GetMetadata()["items"],
	})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal menu: %v", err)), nil
	}

	// Return the menu data as a text result
	return mcp.NewToolResultText(string(result)), nil
}
