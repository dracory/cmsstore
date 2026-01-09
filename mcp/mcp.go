package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dracory/cmsstore"
)

type MCP struct {
	store cmsstore.StoreInterface
}

func NewMCP(store cmsstore.StoreInterface) *MCP {
	return &MCP{store: store}
}

// Handler is an HTTP handler intended to be mounted at a dedicated route.
//
// The protocol is JSON-RPC 2.0 compatible and currently supports:
// - method: "call_tool" with params {"tool_name": string, "arguments": object}
// - method: "list_tools" with params {}
func (m *MCP) Handler(w http.ResponseWriter, r *http.Request) {
	if m == nil || m.store == nil {
		writeJSON(w, http.StatusInternalServerError, jsonRPCErrorResponse(nil, -32603, "store is not initialized"))
		return
	}

	// Check if versioning is enabled
	if !m.store.VersioningEnabled() {
		writeJSON(w, http.StatusForbidden, jsonRPCErrorResponse(nil, -32000, "mcp disabled as versioning is required"))
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, jsonRPCErrorResponse(nil, -32602, "failed to read request body"))
		return
	}

	defer r.Body.Close()

	var req jsonRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(nil, -32700, "parse error"))
		return
	}

	if strings.TrimSpace(req.JSONRPC) == "" {
		req.JSONRPC = "2.0"
	}

	switch req.Method {
	case "initialize":
		m.handleInitialize(w, r.Context(), req.ID, req.Params)
		return
	case "notifications/initialized":
		m.handleInitialized(w, r.Context())
		return
	case "tools/list":
		m.handleToolsList(w, r.Context(), req.ID)
		return
	case "tools/call":
		m.handleToolsCall(w, r.Context(), req.ID, req.Params)
		return
	case "list_tools":
		m.handleToolsList(w, r.Context(), req.ID)
		return
	case "call_tool":
		m.handleToolsCall(w, r.Context(), req.ID, req.Params)
		return
	default:
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(req.ID, -32601, "method not found"))
		return
	}
}

func (m *MCP) handleInitialize(w http.ResponseWriter, _ context.Context, id any, params json.RawMessage) {
	var p struct {
		ProtocolVersion string `json:"protocolVersion"`
		ClientInfo      any    `json:"clientInfo"`
		Capabilities    any    `json:"capabilities"`
	}
	_ = json.Unmarshal(params, &p)

	result := map[string]any{
		"protocolVersion": "2025-06-18",
		"serverInfo": map[string]any{
			"name":    "cmsstore",
			"version": "0.1.0",
		},
		"capabilities": map[string]any{
			"tools": map[string]any{},
		},
		"echo": map[string]any{
			"clientProtocolVersion": p.ProtocolVersion,
			"clientInfo":            p.ClientInfo,
			"clientCapabilities":    p.Capabilities,
		},
	}

	writeJSON(w, http.StatusOK, jsonRPCResultResponse(id, result))
}

func (m *MCP) handleInitialized(w http.ResponseWriter, _ context.Context) {
	// JSON-RPC notifications do not expect a response.
	w.WriteHeader(http.StatusOK)
}

func (m *MCP) handleToolsList(w http.ResponseWriter, _ context.Context, id any) {
	tools := []map[string]any{
		{
			"name":        "cms_schema",
			"description": "Get a JSON schema-like description of CMS entities and supported MCP tools",
			"inputSchema": map[string]any{"type": "object"},
		},
		{
			"name":        "page_list",
			"description": "List CMS pages",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":                map[string]any{"type": "integer"},
					"offset":               map[string]any{"type": "integer"},
					"site_id":              map[string]any{"type": "string"},
					"status":               map[string]any{"type": "string"},
					"name_like":            map[string]any{"type": "string"},
					"alias_like":           map[string]any{"type": "string"},
					"handle":               map[string]any{"type": "string"},
					"template_id":          map[string]any{"type": "string"},
					"include_soft_deleted": map[string]any{"type": "boolean"},
					"order_by":             map[string]any{"type": "string"},
					"sort_order":           map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "page_upsert",
			"description": "Create or update a CMS page (if ID is provided, updates existing page; otherwise creates new page)",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"title"},
				"properties": map[string]any{
					"id":               map[string]any{"type": "string"},
					"title":            map[string]any{"type": "string"},
					"content":          map[string]any{"type": "string"},
					"status":           map[string]any{"type": "string"},
					"site_id":          map[string]any{"type": "string"},
					"template_id":      map[string]any{"type": "string"},
					"alias":            map[string]any{"type": "string"},
					"name":             map[string]any{"type": "string"},
					"handle":           map[string]any{"type": "string"},
					"canonical_url":    map[string]any{"type": "string"},
					"meta_description": map[string]any{"type": "string"},
					"meta_keywords":    map[string]any{"type": "string"},
					"meta_robots":      map[string]any{"type": "string"},
					"memo":             map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "page_get",
			"description": "Get a CMS page by ID",
			"inputSchema": map[string]any{
				"type":       "object",
				"required":   []string{"id"},
				"properties": map[string]any{"id": map[string]any{"type": "string"}},
			},
		},
		{
			"name":        "page_delete",
			"description": "Delete a CMS page",
			"inputSchema": map[string]any{
				"type":       "object",
				"required":   []string{"id"},
				"properties": map[string]any{"id": map[string]any{"type": "string"}},
			},
		},
		{
			"name":        "menu_upsert",
			"description": "Create or update a CMS menu (if ID is provided, updates existing menu; otherwise creates new menu)",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"name"},
				"properties": map[string]any{
					"id":      map[string]any{"type": "string"},
					"name":    map[string]any{"type": "string"},
					"status":  map[string]any{"type": "string"},
					"site_id": map[string]any{"type": "string"},
					"handle":  map[string]any{"type": "string"},
					"memo":    map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "menu_list",
			"description": "List CMS menus",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":                map[string]any{"type": "integer"},
					"offset":               map[string]any{"type": "integer"},
					"site_id":              map[string]any{"type": "string"},
					"status":               map[string]any{"type": "string"},
					"name_like":            map[string]any{"type": "string"},
					"handle":               map[string]any{"type": "string"},
					"include_soft_deleted": map[string]any{"type": "boolean"},
					"order_by":             map[string]any{"type": "string"},
					"sort_order":           map[string]any{"type": "string"},
				},
			},
		},
		{
			"name":        "menu_get",
			"description": "Get a CMS menu by ID",
			"inputSchema": map[string]any{
				"type":       "object",
				"required":   []string{"id"},
				"properties": map[string]any{"id": map[string]any{"type": "string"}},
			},
		},
		{
			"name":        "site_list",
			"description": "List CMS sites",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":                map[string]any{"type": "integer"},
					"offset":               map[string]any{"type": "integer"},
					"status":               map[string]any{"type": "string"},
					"name_like":            map[string]any{"type": "string"},
					"domain_name":          map[string]any{"type": "string"},
					"handle":               map[string]any{"type": "string"},
					"include_soft_deleted": map[string]any{"type": "boolean"},
					"order_by":             map[string]any{"type": "string"},
					"sort_order":           map[string]any{"type": "string"},
				},
			},
		},
	}

	result := map[string]any{
		"tools": tools,
	}
	writeJSON(w, http.StatusOK, jsonRPCResultResponse(id, result))
}

func (m *MCP) handleToolsCall(w http.ResponseWriter, ctx context.Context, id any, params json.RawMessage) {
	// Support both MCP standard params and legacy ones:
	// - MCP: {"name": "tool", "arguments": {...}}
	// - Legacy: {"tool_name": "tool", "arguments": {...}}
	var p struct {
		Name     string          `json:"name"`
		ToolName string          `json:"tool_name"`
		ArgsRaw  json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(params, &p); err != nil {
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(id, -32602, "invalid params"))
		return
	}

	toolName := strings.TrimSpace(p.Name)
	if toolName == "" {
		toolName = strings.TrimSpace(p.ToolName)
	}
	if toolName == "" {
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(id, -32602, "missing tool name"))
		return
	}

	args := map[string]any{}
	if len(p.ArgsRaw) > 0 {
		dec := json.NewDecoder(strings.NewReader(string(p.ArgsRaw)))
		dec.UseNumber()
		if err := dec.Decode(&args); err != nil {
			writeJSON(w, http.StatusOK, jsonRPCErrorResponse(id, -32602, "invalid arguments"))
			return
		}
	}

	textResult, err := m.dispatchTool(ctx, toolName, args)
	if err != nil {
		writeJSON(w, http.StatusOK, jsonRPCErrorResponse(id, -32000, err.Error()))
		return
	}

	result := map[string]any{
		"content": []map[string]any{
			{
				"type": "text",
				"text": textResult,
			},
		},
	}
	writeJSON(w, http.StatusOK, jsonRPCResultResponse(id, result))
}

func (m *MCP) dispatchTool(ctx context.Context, toolName string, args map[string]any) (string, error) {
	switch toolName {
	case "cms_schema":
		return m.toolCmsSchema(ctx, args)
	case "page_list":
		return m.toolPageList(ctx, args)
	case "page_upsert":
		return m.toolPageUpsert(ctx, args)
	case "page_get":
		return m.toolPageGet(ctx, args)
	case "page_delete":
		return m.toolPageDelete(ctx, args)
	case "menu_list":
		return m.toolMenuList(ctx, args)
	case "menu_upsert":
		return m.toolMenuUpsert(ctx, args)
	case "menu_get":
		return m.toolMenuGet(ctx, args)
	case "site_list":
		return m.toolSiteList(ctx, args)
	default:
		return "", errors.New("tool not found")
	}
}

// 	// Return the menu data as a text result
// 	return mcp.NewToolResultText(string(result)), nil
// }
