package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/cmsstore"
)

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

func (m *MCP) toolMenuUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	name, _ := args["name"].(string)
	status, _ := args["status"].(string)
	siteID := strings.TrimSpace(argString(args, "site_id"))

	if strings.TrimSpace(name) == "" {
		return "", errors.New("missing required parameter: name")
	}

	var menu cmsstore.MenuInterface
	var err error

	// If ID is provided, try to find existing menu
	if strings.TrimSpace(id) != "" {
		menu, err = m.store.MenuFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if menu == nil {
			return "", errors.New("menu not found")
		}
	} else {
		// Create new menu
		menu = cmsstore.NewMenu()

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
		menu.SetSiteID(siteID)
	}

	// Set/update fields
	menu.SetName(name)
	if status != "" {
		menu.SetStatus(status)
	}
	if v := strings.TrimSpace(argString(args, "site_id")); v != "" {
		menu.SetSiteID(v)
	}
	if v := strings.TrimSpace(argString(args, "handle")); v != "" {
		menu.SetHandle(v)
	}
	if v := strings.TrimSpace(argString(args, "memo")); v != "" {
		menu.SetMemo(v)
	}

	// Save menu
	if strings.TrimSpace(id) != "" {
		// Update existing menu
		if err := m.store.MenuUpdate(ctx, menu); err != nil {
			return "", err
		}
	} else {
		// Create new menu
		if err := m.store.MenuCreate(ctx, menu); err != nil {
			return "", err
		}
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      menu.ID(),
		"name":    menu.Name(),
		"status":  menu.Status(),
		"site_id": menu.SiteID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}
