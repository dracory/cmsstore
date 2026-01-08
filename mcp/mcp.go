package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func argString(args map[string]any, key string) string {
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case json.Number:
		return t.String()
	case float64:
		return fmt.Sprintf("%.0f", t)
	case int:
		return fmt.Sprintf("%d", t)
	case int64:
		return fmt.Sprintf("%d", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func argInt(args map[string]any, key string) (int, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case json.Number:
		i64, err := t.Int64()
		if err != nil {
			return 0, false
		}
		return int(i64), true
	case float64:
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	default:
		return 0, false
	}
}

func argBool(args map[string]any, key string) (bool, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return false, false
	}
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		vv := strings.TrimSpace(strings.ToLower(t))
		if vv == "true" || vv == "1" || vv == "yes" {
			return true, true
		}
		if vv == "false" || vv == "0" || vv == "no" {
			return false, true
		}
		return false, false
	default:
		return false, false
	}
}

type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

func (m *MCP) handleInitialize(w http.ResponseWriter, ctx context.Context, id any, params json.RawMessage) {
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

func (m *MCP) handleInitialized(w http.ResponseWriter, ctx context.Context) {
	// JSON-RPC notifications do not expect a response.
	w.WriteHeader(http.StatusOK)
}

func (m *MCP) handleToolsList(w http.ResponseWriter, ctx context.Context, id any) {
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
			"name":        "page_create",
			"description": "Create a CMS page",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"title"},
				"properties": map[string]any{
					"title":       map[string]any{"type": "string"},
					"content":     map[string]any{"type": "string"},
					"status":      map[string]any{"type": "string"},
					"site_id":     map[string]any{"type": "string"},
					"template_id": map[string]any{"type": "string"},
					"alias":       map[string]any{"type": "string"},
					"name":        map[string]any{"type": "string"},
					"handle":      map[string]any{"type": "string"},
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
			"name":        "page_update",
			"description": "Update a CMS page",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"id"},
				"properties": map[string]any{
					"id": map[string]any{"type": "string"},
					"updates": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"title":            map[string]any{"type": "string"},
							"content":          map[string]any{"type": "string"},
							"status":           map[string]any{"type": "string"},
							"alias":            map[string]any{"type": "string"},
							"name":             map[string]any{"type": "string"},
							"handle":           map[string]any{"type": "string"},
							"site_id":          map[string]any{"type": "string"},
							"template_id":      map[string]any{"type": "string"},
							"canonical_url":    map[string]any{"type": "string"},
							"meta_description": map[string]any{"type": "string"},
							"meta_keywords":    map[string]any{"type": "string"},
							"meta_robots":      map[string]any{"type": "string"},
							"memo":             map[string]any{"type": "string"},
						},
					},
				},
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
			"name":        "menu_create",
			"description": "Create a CMS menu",
			"inputSchema": map[string]any{
				"type":     "object",
				"required": []string{"name"},
				"properties": map[string]any{
					"name":    map[string]any{"type": "string"},
					"status":  map[string]any{"type": "string"},
					"site_id": map[string]any{"type": "string"},
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
	case "page_create":
		return m.toolPageCreate(ctx, args)
	case "page_get":
		return m.toolPageGet(ctx, args)
	case "page_update":
		return m.toolPageUpdate(ctx, args)
	case "page_delete":
		return m.toolPageDelete(ctx, args)
	case "menu_list":
		return m.toolMenuList(ctx, args)
	case "menu_create":
		return m.toolMenuCreate(ctx, args)
	case "menu_get":
		return m.toolMenuGet(ctx, args)
	case "site_list":
		return m.toolSiteList(ctx, args)
	default:
		return "", errors.New("tool not found")
	}
}

func (m *MCP) toolCmsSchema(ctx context.Context, args map[string]any) (string, error) {
	entities := map[string]any{
		"page": map[string]any{
			"fields": []map[string]any{
				{"name": "id", "type": "string"},
				{"name": "site_id", "type": "string"},
				{"name": "template_id", "type": "string"},
				{"name": "title", "type": "string"},
				{"name": "name", "type": "string"},
				{"name": "handle", "type": "string"},
				{"name": "alias", "type": "string"},
				{"name": "canonical_url", "type": "string"},
				{"name": "content", "type": "string"},
				{"name": "editor", "type": "string"},
				{"name": "memo", "type": "string"},
				{"name": "meta_description", "type": "string"},
				{"name": "meta_keywords", "type": "string"},
				{"name": "meta_robots", "type": "string"},
				{"name": "middlewares_before", "type": "array", "items": map[string]any{"type": "string"}},
				{"name": "middlewares_after", "type": "array", "items": map[string]any{"type": "string"}},
				{"name": "status", "type": "string"},
				{"name": "created_at", "type": "string"},
				{"name": "updated_at", "type": "string"},
				{"name": "soft_deleted_at", "type": "string"},
			},
		},
		"menu": map[string]any{
			"fields": []map[string]any{
				{"name": "id", "type": "string"},
				{"name": "site_id", "type": "string"},
				{"name": "name", "type": "string"},
				{"name": "handle", "type": "string"},
				{"name": "memo", "type": "string"},
				{"name": "status", "type": "string"},
				{"name": "created_at", "type": "string"},
				{"name": "updated_at", "type": "string"},
				{"name": "soft_deleted_at", "type": "string"},
			},
		},
		"site": map[string]any{
			"fields": []map[string]any{
				{"name": "id", "type": "string"},
				{"name": "name", "type": "string"},
				{"name": "handle", "type": "string"},
				{"name": "domain_names", "type": "array", "items": map[string]any{"type": "string"}},
				{"name": "memo", "type": "string"},
				{"name": "status", "type": "string"},
				{"name": "created_at", "type": "string"},
				{"name": "updated_at", "type": "string"},
				{"name": "soft_deleted_at", "type": "string"},
			},
		},
	}

	tools := map[string]any{
		"page_list": map[string]any{
			"arguments": []map[string]any{
				{"name": "limit", "type": "integer"},
				{"name": "offset", "type": "integer"},
				{"name": "site_id", "type": "string"},
				{"name": "status", "type": "string"},
				{"name": "name_like", "type": "string"},
				{"name": "alias_like", "type": "string"},
				{"name": "handle", "type": "string"},
				{"name": "template_id", "type": "string"},
				{"name": "include_soft_deleted", "type": "boolean"},
				{"name": "order_by", "type": "string"},
				{"name": "sort_order", "type": "string"},
			},
			"returns": map[string]any{
				"items": "array[page]",
			},
		},
		"menu_list": map[string]any{
			"arguments": []map[string]any{
				{"name": "limit", "type": "integer"},
				{"name": "offset", "type": "integer"},
				{"name": "site_id", "type": "string"},
				{"name": "status", "type": "string"},
				{"name": "name_like", "type": "string"},
				{"name": "handle", "type": "string"},
				{"name": "include_soft_deleted", "type": "boolean"},
				{"name": "order_by", "type": "string"},
				{"name": "sort_order", "type": "string"},
			},
			"returns": map[string]any{
				"items": "array[menu]",
			},
		},
		"site_list": map[string]any{
			"arguments": []map[string]any{
				{"name": "limit", "type": "integer"},
				{"name": "offset", "type": "integer"},
				{"name": "status", "type": "string"},
				{"name": "name_like", "type": "string"},
				{"name": "domain_name", "type": "string"},
				{"name": "handle", "type": "string"},
				{"name": "include_soft_deleted", "type": "boolean"},
				{"name": "order_by", "type": "string"},
				{"name": "sort_order", "type": "string"},
			},
			"returns": map[string]any{
				"items": "array[site]",
			},
		},
		"page_create": map[string]any{
			"arguments": []map[string]any{
				{"name": "title", "type": "string", "required": true},
				{"name": "content", "type": "string"},
				{"name": "status", "type": "string"},
				{"name": "site_id", "type": "string"},
			},
			"returns": map[string]any{
				"page": "page",
			},
		},
		"page_get": map[string]any{
			"arguments": []map[string]any{{"name": "id", "type": "string", "required": true}},
			"returns":   map[string]any{"page": "page"},
		},
		"page_update": map[string]any{
			"arguments": []map[string]any{{"name": "id", "type": "string", "required": true}},
			"returns":   map[string]any{"page": "page"},
		},
		"page_delete": map[string]any{
			"arguments": []map[string]any{{"name": "id", "type": "string", "required": true}},
			"returns":   map[string]any{"deleted": "boolean"},
		},
		"menu_create": map[string]any{
			"arguments": []map[string]any{{"name": "name", "type": "string", "required": true}},
			"returns":   map[string]any{"menu": "menu"},
		},
		"menu_get": map[string]any{
			"arguments": []map[string]any{{"name": "id", "type": "string", "required": true}},
			"returns":   map[string]any{"menu": "menu"},
		},
	}

	respBytes, err := json.Marshal(map[string]any{
		"entities": entities,
		"tools":    tools,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolPageList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.PageQuery()

	if v, ok := args["site_id"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSiteID(v)
	}
	if v, ok := args["status"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetStatus(v)
	}
	if v, ok := args["name_like"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetNameLike(v)
	}
	if v, ok := args["alias_like"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetAliasLike(v)
	}
	if v, ok := args["handle"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetHandle(v)
	}
	if v := strings.TrimSpace(argString(args, "template_id")); v != "" {
		q.SetTemplateID(v)
	}
	if v, ok := argBool(args, "include_soft_deleted"); ok {
		q.SetSoftDeletedIncluded(v)
	}
	if v, ok := argInt(args, "limit"); ok {
		q.SetLimit(v)
	}
	if v, ok := argInt(args, "offset"); ok {
		q.SetOffset(v)
	}
	if v, ok := args["order_by"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetOrderBy(v)
	}
	if v, ok := args["sort_order"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSortOrder(v)
	}

	pages, err := m.store.PageList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(pages))
	for _, p := range pages {
		if p == nil {
			continue
		}
		items = append(items, map[string]any{
			"id":          p.ID(),
			"title":       p.Title(),
			"name":        p.Name(),
			"handle":      p.Handle(),
			"alias":       p.Alias(),
			"status":      p.Status(),
			"site_id":     p.SiteID(),
			"template_id": p.TemplateID(),
			"created_at":  p.CreatedAt(),
			"updated_at":  p.UpdatedAt(),
		})
	}

	respBytes, err := json.Marshal(map[string]any{
		"items": items,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolMenuList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.MenuQuery()

	if v, ok := args["site_id"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSiteID(v)
	}
	if v, ok := args["status"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetStatus(v)
	}
	if v, ok := args["name_like"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetNameLike(v)
	}
	if v, ok := args["handle"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetHandle(v)
	}
	if v, ok := argBool(args, "include_soft_deleted"); ok {
		q.SetSoftDeletedIncluded(v)
	}
	if v, ok := argInt(args, "limit"); ok {
		q.SetLimit(v)
	}
	if v, ok := argInt(args, "offset"); ok {
		q.SetOffset(v)
	}
	if v, ok := args["order_by"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetOrderBy(v)
	}
	if v, ok := args["sort_order"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSortOrder(v)
	}

	menus, err := m.store.MenuList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(menus))
	for _, menu := range menus {
		if menu == nil {
			continue
		}
		items = append(items, map[string]any{
			"id":         menu.ID(),
			"name":       menu.Name(),
			"handle":     menu.Handle(),
			"status":     menu.Status(),
			"site_id":    menu.SiteID(),
			"created_at": menu.CreatedAt(),
			"updated_at": menu.UpdatedAt(),
		})
	}

	respBytes, err := json.Marshal(map[string]any{
		"items": items,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolSiteList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.SiteQuery()

	if v, ok := args["status"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetStatus(v)
	}
	if v, ok := args["name_like"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetNameLike(v)
	}
	if v, ok := args["domain_name"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetDomainName(v)
	}
	if v, ok := args["handle"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetHandle(v)
	}
	if v, ok := argBool(args, "include_soft_deleted"); ok {
		q.SetSoftDeletedIncluded(v)
	}
	if v, ok := argInt(args, "limit"); ok {
		q.SetLimit(v)
	}
	if v, ok := argInt(args, "offset"); ok {
		q.SetOffset(v)
	}
	if v, ok := args["order_by"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetOrderBy(v)
	}
	if v, ok := args["sort_order"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetSortOrder(v)
	}

	sites, err := m.store.SiteList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(sites))
	for _, s := range sites {
		if s == nil {
			continue
		}
		domains, _ := s.DomainNames()
		items = append(items, map[string]any{
			"id":          s.ID(),
			"name":        s.Name(),
			"handle":      s.Handle(),
			"status":      s.Status(),
			"domainNames": domains,
			"created_at":  s.CreatedAt(),
			"updated_at":  s.UpdatedAt(),
		})
	}

	respBytes, err := json.Marshal(map[string]any{
		"items": items,
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolPageCreate(ctx context.Context, args map[string]any) (string, error) {
	title, _ := args["title"].(string)
	content, _ := args["content"].(string)
	status, _ := args["status"].(string)
	siteID, _ := args["site_id"].(string)

	if strings.TrimSpace(title) == "" {
		return "", errors.New("missing required parameter: title")
	}

	if siteID == "" {
		sites, err := m.store.SiteList(ctx, cmsstore.SiteQuery())
		if err != nil {
			return "", err
		}
		if len(sites) > 0 {
			siteID = sites[0].ID()
		} else {
			defaultSite := cmsstore.NewSite()
			defaultSite.SetName("Default Site")
			defaultSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
			if err := m.store.SiteCreate(ctx, defaultSite); err != nil {
				return "", err
			}
			siteID = defaultSite.ID()
		}
	}

	page := cmsstore.NewPage()
	page.SetTitle(title)
	page.SetContent(content)
	page.SetSiteID(siteID)
	if status != "" {
		page.SetStatus(status)
	}

	if err := m.store.PageCreate(ctx, page); err != nil {
		return "", err
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
		"site_id": page.SiteID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolPageGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	page, err := m.store.PageFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if page == nil {
		return "", errors.New("page not found")
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
		"site_id": page.SiteID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolPageUpdate(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	updatesAny, _ := args["updates"].(map[string]any)
	if updatesAny == nil {
		updatesAny = map[string]any{}
		for k, v := range args {
			if k == "id" {
				continue
			}
			updatesAny[k] = v
		}
	}
	if len(updatesAny) == 0 {
		return "", errors.New("missing required parameter: updates")
	}

	page, err := m.store.PageFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if page == nil {
		return "", errors.New("page not found")
	}

	if v := strings.TrimSpace(argString(updatesAny, "title")); v != "" {
		page.SetTitle(v)
	}
	if v, ok := updatesAny["content"].(string); ok {
		page.SetContent(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "status")); v != "" {
		page.SetStatus(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "alias")); v != "" {
		page.SetAlias(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "name")); v != "" {
		page.SetName(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "handle")); v != "" {
		page.SetHandle(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "template_id")); v != "" {
		page.SetTemplateID(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "site_id")); v != "" {
		page.SetSiteID(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "canonical_url")); v != "" {
		page.SetCanonicalUrl(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "meta_description")); v != "" {
		page.SetMetaDescription(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "meta_keywords")); v != "" {
		page.SetMetaKeywords(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "meta_robots")); v != "" {
		page.SetMetaRobots(v)
	}
	if v := strings.TrimSpace(argString(updatesAny, "memo")); v != "" {
		page.SetMemo(v)
	}

	if err := m.store.PageUpdate(ctx, page); err != nil {
		return "", err
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      page.ID(),
		"title":   page.Title(),
		"content": page.Content(),
		"status":  page.Status(),
		"site_id": page.SiteID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolPageDelete(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	page, err := m.store.PageFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if page == nil {
		return "", errors.New("page not found")
	}

	if err := m.store.PageSoftDeleteByID(ctx, id); err != nil {
		if err := m.store.PageDelete(ctx, page); err != nil {
			return "", err
		}
	}

	respBytes, err := json.Marshal(map[string]any{"id": id})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolMenuCreate(ctx context.Context, args map[string]any) (string, error) {
	name, _ := args["name"].(string)
	status, _ := args["status"].(string)
	siteID := strings.TrimSpace(argString(args, "site_id"))
	if strings.TrimSpace(name) == "" {
		return "", errors.New("missing required parameter: name")
	}

	if siteID == "" {
		sites, err := m.store.SiteList(ctx, cmsstore.SiteQuery())
		if err != nil {
			return "", err
		}
		if len(sites) > 0 {
			siteID = sites[0].ID()
		} else {
			defaultSite := cmsstore.NewSite()
			defaultSite.SetName("Default Site")
			defaultSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
			if err := m.store.SiteCreate(ctx, defaultSite); err != nil {
				return "", err
			}
			siteID = defaultSite.ID()
		}
	}

	menu := cmsstore.NewMenu()
	menu.SetName(name)
	menu.SetSiteID(siteID)
	if status != "" {
		menu.SetStatus(status)
	}

	if err := m.store.MenuCreate(ctx, menu); err != nil {
		return "", err
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":     menu.ID(),
		"name":   menu.Name(),
		"status": menu.Status(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolMenuGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	menu, err := m.store.MenuFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if menu == nil {
		return "", errors.New("menu not found")
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":     menu.ID(),
		"name":   menu.Name(),
		"status": menu.Status(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func jsonRPCErrorResponse(id any, code int, message string) map[string]any {
	return map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]any{
			"code":    code,
			"message": message,
		},
	}
}

func jsonRPCResultResponse(id any, result any) map[string]any {
	return map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}
}

// 	// Return the menu data as a text result
// 	return mcp.NewToolResultText(string(result)), nil
// }
