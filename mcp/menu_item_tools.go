package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/dracory/cmsstore"
)

func (m *MCP) toolMenuItemGet(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	menuItem, err := m.store.MenuItemFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if menuItem == nil {
		return "", errors.New("menu item not found")
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      menuItem.ID(),
		"name":    menuItem.Name(),
		"url":     menuItem.URL(),
		"target":  menuItem.Target(),
		"status":  menuItem.Status(),
		"menu_id": menuItem.MenuID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolMenuItemList(ctx context.Context, args map[string]any) (string, error) {
	q := cmsstore.MenuItemQuery()

	if v, ok := args["menu_id"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetMenuID(v)
	}
	if v, ok := args["status"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetStatus(v)
	}
	if v, ok := args["name_like"].(string); ok && strings.TrimSpace(v) != "" {
		q.SetNameLike(v)
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

	menuItems, err := m.store.MenuItemList(ctx, q)
	if err != nil {
		return "", err
	}

	items := make([]map[string]any, 0, len(menuItems))
	for _, menuItem := range menuItems {
		if menuItem == nil {
			continue
		}
		items = append(items, map[string]any{
			"id":         menuItem.ID(),
			"name":       menuItem.Name(),
			"handle":     menuItem.Handle(),
			"url":        menuItem.URL(),
			"target":     menuItem.Target(),
			"status":     menuItem.Status(),
			"menu_id":    menuItem.MenuID(),
			"page_id":    menuItem.PageID(),
			"parent_id":  menuItem.ParentID(),
			"sequence":   menuItem.Sequence(),
			"created_at": menuItem.CreatedAt(),
			"updated_at": menuItem.UpdatedAt(),
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

func (m *MCP) toolMenuItemUpsert(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	name, _ := args["name"].(string)
	url, _ := args["url"].(string)
	target, _ := args["target"].(string)
	status, _ := args["status"].(string)
	menuID, _ := args["menu_id"].(string)

	if strings.TrimSpace(name) == "" {
		return "", errors.New("missing required parameter: name")
	}

	var menuItem cmsstore.MenuItemInterface
	var err error

	// If ID is provided, try to find existing menu item
	if strings.TrimSpace(id) != "" {
		menuItem, err = m.store.MenuItemFindByID(ctx, id)
		if err != nil {
			return "", err
		}
		if menuItem == nil {
			return "", errors.New("menu item not found")
		}
	} else {
		// Create new menu item
		menuItem = cmsstore.NewMenuItem()

		// Set menu ID if not provided
		if menuID == "" {
			menus, err := m.store.MenuList(ctx, cmsstore.MenuQuery())
			if err != nil {
				return "", err
			}
			if len(menus) > 0 {
				menuID = menus[0].ID()
			} else {
				// Create a default menu if none exists
				defaultMenu := cmsstore.NewMenu()
				defaultMenu.SetName("Default Menu")
				defaultMenu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
				if err := m.store.MenuCreate(ctx, defaultMenu); err != nil {
					return "", err
				}
				menuID = defaultMenu.ID()
			}
		}
		menuItem.SetMenuID(menuID)
	}

	// Set/update fields
	menuItem.SetName(name)
	if url != "" {
		menuItem.SetURL(url)
	}
	if target != "" {
		menuItem.SetTarget(target)
	}
	if status != "" {
		menuItem.SetStatus(status)
	}
	if v := strings.TrimSpace(argString(args, "menu_id")); v != "" {
		menuItem.SetMenuID(v)
	}
	if v := strings.TrimSpace(argString(args, "page_id")); v != "" {
		menuItem.SetPageID(v)
	}
	if v := strings.TrimSpace(argString(args, "parent_id")); v != "" {
		menuItem.SetParentID(v)
	}
	if v := strings.TrimSpace(argString(args, "handle")); v != "" {
		menuItem.SetHandle(v)
	}
	if v := strings.TrimSpace(argString(args, "memo")); v != "" {
		menuItem.SetMemo(v)
	}
	if v, ok := argInt(args, "sequence"); ok {
		menuItem.SetSequenceInt(int(v))
	}

	// Save menu item
	if strings.TrimSpace(id) != "" {
		// Update existing menu item
		if err := m.store.MenuItemUpdate(ctx, menuItem); err != nil {
			return "", err
		}
	} else {
		// Create new menu item
		if err := m.store.MenuItemCreate(ctx, menuItem); err != nil {
			return "", err
		}
	}

	// Create versioning record if versioning is enabled
	if m.store.VersioningEnabled() {
		if err := m.createMenuItemVersioning(ctx, menuItem); err != nil {
			// Log error but don't fail the operation
			// In a production environment, you might want to handle this differently
		}
	}

	respBytes, err := json.Marshal(map[string]any{
		"id":      menuItem.ID(),
		"name":    menuItem.Name(),
		"url":     menuItem.URL(),
		"target":  menuItem.Target(),
		"status":  menuItem.Status(),
		"menu_id": menuItem.MenuID(),
	})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

func (m *MCP) toolMenuItemDelete(ctx context.Context, args map[string]any) (string, error) {
	id := argString(args, "id")
	if strings.TrimSpace(id) == "" {
		return "", errors.New("missing required parameter: id")
	}

	menuItem, err := m.store.MenuItemFindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if menuItem == nil {
		return "", errors.New("menu item not found")
	}

	if err := m.store.MenuItemSoftDeleteByID(ctx, id); err != nil {
		if err := m.store.MenuItemDelete(ctx, menuItem); err != nil {
			return "", err
		}
	}

	respBytes, err := json.Marshal(map[string]any{"id": id})
	if err != nil {
		return "", err
	}
	return string(respBytes), nil
}

// createMenuItemVersioning creates a versioning record for a menu item if versioning is enabled
func (m *MCP) createMenuItemVersioning(ctx context.Context, menuItem cmsstore.MenuItemInterface) error {
	if !m.store.VersioningEnabled() {
		return nil
	}

	if menuItem == nil {
		return errors.New("menu item is nil")
	}

	// Get last versioning to check if content has changed
	lastVersioningList, err := m.store.VersioningList(ctx, cmsstore.NewVersioningQuery().
		SetEntityType(cmsstore.VERSIONING_TYPE_MENU_ITEM).
		SetEntityID(menuItem.ID()).
		SetOrderBy("created_at").
		SetSortOrder("DESC").
		SetLimit(1))

	if err != nil {
		return err
	}

	// Marshal menu item content for versioning
	menuItemData := map[string]any{
		"id":        menuItem.ID(),
		"name":      menuItem.Name(),
		"url":       menuItem.URL(),
		"target":    menuItem.Target(),
		"handle":    menuItem.Handle(),
		"status":    menuItem.Status(),
		"menu_id":   menuItem.MenuID(),
		"page_id":   menuItem.PageID(),
		"parent_id": menuItem.ParentID(),
		"sequence":  menuItem.Sequence(),
		"memo":      menuItem.Memo(),
	}
	content, err := json.Marshal(menuItemData)
	if err != nil {
		return err
	}

	// Check if last versioning has the same content
	if len(lastVersioningList) > 0 {
		lastVersioning := lastVersioningList[0]
		if lastVersioning.Content() == string(content) {
			return nil // No change needed
		}
	}

	// Create new versioning record
	return m.store.VersioningCreate(ctx, cmsstore.NewVersioning().
		SetEntityID(menuItem.ID()).
		SetEntityType(cmsstore.VERSIONING_TYPE_MENU_ITEM).
		SetContent(string(content)))
}
