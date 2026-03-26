package frontend

import (
	"context"

	"github.com/dracory/cmsstore"
	"github.com/spf13/cast"
)

// MenuBlockRenderer renders menu blocks by loading menu items and generating HTML
type MenuBlockRenderer struct {
	store *frontend
}

// NewMenuBlockRenderer creates a new menu block renderer
func NewMenuBlockRenderer(store *frontend) *MenuBlockRenderer {
	return &MenuBlockRenderer{
		store: store,
	}
}

// Render renders a menu block by loading the menu and rendering it as HTML
func (r *MenuBlockRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	if !r.store.store.MenusEnabled() {
		r.store.logger.Debug("renderMenuBlock: Menus not enabled")
		return "<!-- Menus not enabled -->", nil
	}

	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		r.store.logger.Debug("renderMenuBlock: No menu ID in block meta", "blockID", block.ID())
		return "<!-- No menu selected -->", nil
	}

	menu, err := r.store.store.MenuFindByID(ctx, menuID)
	if err != nil {
		r.store.logger.Error("renderMenuBlock: Error finding menu", "menuID", menuID, "error", err)
		return "", err
	}

	if menu == nil {
		r.store.logger.Debug("renderMenuBlock: Menu not found", "menuID", menuID)
		return "<!-- Menu not found -->", nil
	}

	if !menu.IsActive() {
		r.store.logger.Debug("renderMenuBlock: Menu not active", "menuID", menuID, "status", menu.Status())
		return "<!-- Menu not active -->", nil
	}

	menuItems, err := r.store.store.MenuItemList(ctx, cmsstore.MenuItemQuery().
		SetMenuID(menuID).
		SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_SEQUENCE).
		SetSortOrder("asc"))

	if err != nil {
		r.store.logger.Error("renderMenuBlock: Error listing menu items", "menuID", menuID, "error", err)
		return "", err
	}

	if len(menuItems) == 0 {
		r.store.logger.Debug("renderMenuBlock: No active menu items found", "menuID", menuID, "menuName", menu.Name())
	}

	style := block.Meta(cmsstore.BLOCK_META_MENU_STYLE)
	if style == "" {
		style = cmsstore.BLOCK_MENU_STYLE_VERTICAL
	}

	cssClass := block.Meta(cmsstore.BLOCK_META_MENU_CSS_CLASS)
	startLevel := cast.ToInt(block.Meta(cmsstore.BLOCK_META_MENU_START_LEVEL))
	maxDepth := cast.ToInt(block.Meta(cmsstore.BLOCK_META_MENU_MAX_DEPTH))

	return r.store.renderMenuHTML(ctx, menuItems, style, cssClass, startLevel, maxDepth)
}
