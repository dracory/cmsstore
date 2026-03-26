package menu

import (
	"context"
	"log/slog"

	"github.com/dracory/cmsstore"
	"github.com/spf13/cast"
)

// BlockRenderer renders menu blocks by loading menu items and generating HTML
type BlockRenderer struct {
	store FrontendStore
}

// FrontendStore interface for store operations needed by menu renderer
type FrontendStore interface {
	MenuFindByID(ctx context.Context, id string) (cmsstore.MenuInterface, error)
	MenuItemList(ctx context.Context, query cmsstore.MenuItemQueryInterface) ([]cmsstore.MenuItemInterface, error)
	MenusEnabled() bool
	PageFindByID(ctx context.Context, id string) (cmsstore.PageInterface, error)
	Logger() *slog.Logger
}

// NewBlockRenderer creates a new menu block renderer
func NewBlockRenderer(store FrontendStore) *BlockRenderer {
	return &BlockRenderer{
		store: store,
	}
}

// Render renders a menu block by loading the menu and rendering it as HTML
func (r *BlockRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	if block == nil {
		r.store.Logger().Debug("renderMenuBlock: Block is nil")
		return "<!-- Block is nil -->", nil
	}

	if !r.store.MenusEnabled() {
		r.store.Logger().Debug("renderMenuBlock: Menus not enabled")
		return "<!-- Menus not enabled -->", nil
	}

	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		r.store.Logger().Debug("renderMenuBlock: No menu ID in block meta", "blockID", block.ID())
		return "<!-- No menu selected -->", nil
	}

	menu, err := r.store.MenuFindByID(ctx, menuID)
	if err != nil {
		r.store.Logger().Error("renderMenuBlock: Error finding menu", "menuID", menuID, "error", err)
		return "", err
	}

	if menu == nil {
		r.store.Logger().Debug("renderMenuBlock: Menu not found", "menuID", menuID)
		return "<!-- Menu not found -->", nil
	}

	if !menu.IsActive() {
		r.store.Logger().Debug("renderMenuBlock: Menu not active", "menuID", menuID, "status", menu.Status())
		return "<!-- Menu not active -->", nil
	}

	menuItems, err := r.store.MenuItemList(ctx, cmsstore.MenuItemQuery().
		SetMenuID(menuID).
		SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_SEQUENCE).
		SetSortOrder("asc"))

	if err != nil {
		r.store.Logger().Error("renderMenuBlock: Error listing menu items", "menuID", menuID, "error", err)
		return "", err
	}

	if len(menuItems) == 0 {
		r.store.Logger().Debug("renderMenuBlock: No active menu items found", "menuID", menuID, "menuName", menu.Name())
	}

	style := block.Meta(cmsstore.BLOCK_META_MENU_STYLE)
	if style == "" {
		style = cmsstore.BLOCK_MENU_STYLE_VERTICAL
	}

	cssClass := block.Meta(cmsstore.BLOCK_META_MENU_CSS_CLASS)
	startLevelStr := block.Meta(cmsstore.BLOCK_META_MENU_START_LEVEL)
	maxDepthStr := block.Meta(cmsstore.BLOCK_META_MENU_MAX_DEPTH)

	startLevel, err := cast.ToIntE(startLevelStr)
	if err != nil {
		startLevel = 0 // default value
	}

	maxDepth, err := cast.ToIntE(maxDepthStr)
	if err != nil {
		maxDepth = 0 // default value
	}

	menuRenderer := NewMenuRenderer(r.store)
	return menuRenderer.RenderMenuHTML(ctx, menuItems, style, cssClass, startLevel, maxDepth)
}
