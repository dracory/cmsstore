package mcp

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/dracory/cmsstore"
)

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
		q.SetLimit(int(v))
	}
	if v, ok := argInt(args, "offset"); ok {
		q.SetOffset(int(v))
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
