package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/cmsstore"
)

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

func (m *MCP) toolPageUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	title, _ := args["title"].(string)
	content, _ := args["content"].(string)
	status, _ := args["status"].(string)
	siteID, _ := args["site_id"].(string)

	if strings.TrimSpace(title) == "" {
		return "", errors.New("missing required parameter: title")
	}

	var page cmsstore.PageInterface
	var err error

	// If ID is provided, try to find existing page
	if strings.TrimSpace(id) != "" {
		page, err = m.store.PageFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if page == nil {
			return "", errors.New("page not found")
		}
	} else {
		// Create new page
		page = cmsstore.NewPage()

		// Set site ID
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
		page.SetSiteID(siteID)
	}

	// Set/update fields
	page.SetTitle(title)
	if content != "" {
		page.SetContent(content)
	}
	if status != "" {
		page.SetStatus(status)
	}
	if v := strings.TrimSpace(argString(args, "alias")); v != "" {
		page.SetAlias(v)
	}
	if v := strings.TrimSpace(argString(args, "name")); v != "" {
		page.SetName(v)
	}
	if v := strings.TrimSpace(argString(args, "handle")); v != "" {
		page.SetHandle(v)
	}
	if v := strings.TrimSpace(argString(args, "template_id")); v != "" {
		page.SetTemplateID(v)
	}
	if v := strings.TrimSpace(argString(args, "site_id")); v != "" {
		page.SetSiteID(v)
	}
	if v := strings.TrimSpace(argString(args, "canonical_url")); v != "" {
		page.SetCanonicalUrl(v)
	}
	if v := strings.TrimSpace(argString(args, "meta_description")); v != "" {
		page.SetMetaDescription(v)
	}
	if v := strings.TrimSpace(argString(args, "meta_keywords")); v != "" {
		page.SetMetaKeywords(v)
	}
	if v := strings.TrimSpace(argString(args, "meta_robots")); v != "" {
		page.SetMetaRobots(v)
	}
	if v := strings.TrimSpace(argString(args, "memo")); v != "" {
		page.SetMemo(v)
	}

	// Save page
	if strings.TrimSpace(id) != "" {
		// Update existing page
		if err := m.store.PageUpdate(ctx, page); err != nil {
			return "", err
		}
	} else {
		// Create new page
		if err := m.store.PageCreate(ctx, page); err != nil {
			return "", err
		}
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
