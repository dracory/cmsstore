package mcp

import (
	"context"
	"encoding/json"
)

func (m *MCP) toolCmsSchema(_ context.Context, _ map[string]any) (string, error) {
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
