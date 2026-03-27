package menu

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/spf13/cast"
)

// MenuBlockType provides both frontend rendering and admin UI for menu blocks.
//
// This is a built-in block type that renders navigation menus with various styles.
type MenuBlockType struct {
	store  cmsstore.StoreInterface
	logger interface {
		Error(msg string, args ...interface{})
	}
}

// NewMenuBlockType creates a new menu block type.
func NewMenuBlockType(store cmsstore.StoreInterface, logger interface {
	Error(msg string, args ...interface{})
}) *MenuBlockType {
	return &MenuBlockType{
		store:  store,
		logger: logger,
	}
}

// TypeKey returns the unique identifier for menu blocks.
func (t *MenuBlockType) TypeKey() string {
	return cmsstore.BLOCK_TYPE_MENU
}

// TypeLabel returns the display name for menu blocks.
func (t *MenuBlockType) TypeLabel() string {
	return "Menu Block"
}

// Render renders a menu block by loading menu items and generating HTML.
// Supports runtime attributes: depth, style, class, id for dynamic configuration.
func (t *MenuBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
	if block == nil {
		t.logger.Error("renderMenuBlock: Block is nil")
		return "<!-- Block is nil -->", nil
	}

	if !t.store.MenusEnabled() {
		t.logger.Error("renderMenuBlock: Menus not enabled")
		return "<!-- Menus not enabled -->", nil
	}

	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		t.logger.Error("renderMenuBlock: No menu ID in block meta", "blockID", block.ID())
		return "<!-- No menu selected -->", nil
	}

	menu, err := t.store.MenuFindByID(ctx, menuID)
	if err != nil {
		t.logger.Error("renderMenuBlock: Error finding menu", "menuID", menuID, "error", err)
		return "", err
	}

	if menu == nil {
		t.logger.Error("renderMenuBlock: Menu not found", "menuID", menuID)
		return "<!-- Menu not found -->", nil
	}

	if !menu.IsActive() {
		t.logger.Error("renderMenuBlock: Menu not active", "menuID", menuID, "status", menu.Status())
		return "<!-- Menu not active -->", nil
	}

	menuItems, err := t.store.MenuItemList(ctx, cmsstore.MenuItemQuery().
		SetMenuID(menuID).
		SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_SEQUENCE).
		SetSortOrder("asc"))

	if err != nil {
		t.logger.Error("renderMenuBlock: Error listing menu items", "menuID", menuID, "error", err)
		return "", err
	}

	if len(menuItems) == 0 {
		t.logger.Error("renderMenuBlock: No active menu items found", "menuID", menuID, "menuName", menu.Name())
	}

	// Parse render options
	options := &cmsstore.RenderOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Get style: runtime attribute > block meta > default
	style := options.Attributes["style"]
	if style == "" {
		style = block.Meta(cmsstore.BLOCK_META_MENU_STYLE)
	}
	if style == "" {
		style = cmsstore.BLOCK_MENU_STYLE_VERTICAL
	}

	// Get rendering mode: runtime attribute > block meta > default
	renderingMode := options.Attributes["mode"]
	if renderingMode == "" {
		renderingMode = block.Meta(cmsstore.BLOCK_META_MENU_RENDERING_MODE)
	}
	if renderingMode == "" {
		renderingMode = cmsstore.BLOCK_MENU_RENDERING_PLAIN
	}

	// Get CSS class: runtime attribute > block meta
	cssClass := options.Attributes["class"]
	if cssClass == "" {
		cssClass = block.Meta(cmsstore.BLOCK_META_MENU_CSS_CLASS)
	}

	// Get CSS ID: runtime attribute > block meta
	cssID := options.Attributes["id"]
	if cssID == "" {
		cssID = block.Meta(cmsstore.BLOCK_META_MENU_CSS_ID)
	}

	// Get start level: runtime attribute > block meta > 0
	startLevel := cast.ToInt(options.Attributes["start-level"])
	if startLevel == 0 {
		startLevel = cast.ToInt(block.Meta(cmsstore.BLOCK_META_MENU_START_LEVEL))
	}

	// Get max depth: runtime attribute > block meta > 0 (unlimited)
	maxDepth := cast.ToInt(options.Attributes["depth"])
	if maxDepth == 0 {
		maxDepth = cast.ToInt(block.Meta(cmsstore.BLOCK_META_MENU_MAX_DEPTH))
	}

	// Use the menu renderer from frontend/blocks/menu package
	// This delegates to the existing comprehensive menu rendering logic
	return renderMenuHTML(ctx, t.store, menuItems, renderingMode, style, cssClass, cssID, startLevel, maxDepth)
}

// GetAdminFields returns form fields for editing menu block configuration.
func (t *MenuBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	menusEnabled := t.store.MenusEnabled()

	if !menusEnabled {
		return []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Menu Blocks Not Available",
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: hb.Div().Class("alert alert-warning").Text("Menu functionality is not enabled in this CMS installation.").ToHTML(),
			}),
		}
	}

	menuList, err := t.store.MenuList(r.Context(), cmsstore.MenuQuery().
		SetStatus(cmsstore.MENU_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC))

	if err != nil {
		t.logger.Error("Error loading menus", "error", err.Error())
	}

	menuOptions := []form.FieldOption{
		{
			Value: "- Select Menu -",
			Key:   "",
		},
	}

	for _, menu := range menuList {
		menuOptions = append(menuOptions, form.FieldOption{
			Value: menu.Name(),
			Key:   menu.ID(),
		})
	}

	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label:    "Menu",
			Name:     "menu_id",
			Type:     form.FORM_FIELD_TYPE_SELECT,
			Value:    block.Meta(cmsstore.BLOCK_META_MENU_ID),
			Required: true,
			Help:     "Select the menu to display in this block",
			Options:  menuOptions,
		}),
		form.NewField(form.FieldOptions{
			Label: "Menu Style",
			Name:  "menu_style",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_MENU_STYLE),
			Help:  "Choose how the menu should be displayed",
			Options: []form.FieldOption{
				{
					Value: "Vertical (default)",
					Key:   cmsstore.BLOCK_MENU_STYLE_VERTICAL,
				},
				{
					Value: "Horizontal",
					Key:   cmsstore.BLOCK_MENU_STYLE_HORIZONTAL,
				},
				{
					Value: "Dropdown",
					Key:   cmsstore.BLOCK_MENU_STYLE_DROPDOWN,
				},
				{
					Value: "Breadcrumb",
					Key:   cmsstore.BLOCK_MENU_STYLE_BREADCRUMB,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Rendering Mode",
			Name:  "menu_rendering_mode",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_MENU_RENDERING_MODE),
			Help:  "Choose the rendering framework for the menu",
			Options: []form.FieldOption{
				{
					Value: "Plain (default)",
					Key:   cmsstore.BLOCK_MENU_RENDERING_PLAIN,
				},
				{
					Value: "Bootstrap 5",
					Key:   cmsstore.BLOCK_MENU_RENDERING_BOOTSTRAP5,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "CSS ID",
			Name:  "menu_css_id",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_MENU_CSS_ID),
			Help:  "Optional CSS ID for unique identification",
		}),
		form.NewField(form.FieldOptions{
			Label: "CSS Class",
			Name:  "menu_css_class",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_MENU_CSS_CLASS),
			Help:  "Optional CSS class for custom styling",
		}),
		form.NewField(form.FieldOptions{
			Label: "Start Level",
			Name:  "menu_start_level",
			Type:  form.FORM_FIELD_TYPE_NUMBER,
			Value: block.Meta(cmsstore.BLOCK_META_MENU_START_LEVEL),
			Help:  "Start rendering from this level (0 = root level)",
		}),
		form.NewField(form.FieldOptions{
			Label: "Max Depth",
			Name:  "menu_max_depth",
			Type:  form.FORM_FIELD_TYPE_NUMBER,
			Value: block.Meta(cmsstore.BLOCK_META_MENU_MAX_DEPTH),
			Help:  "Maximum depth to render (0 = unlimited)",
		}),
	}

	return fieldsContent
}

// SaveAdminFields processes form submission and updates the menu block.
func (t *MenuBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	menuID := req.GetStringTrimmed(r, "menu_id")
	menuStyle := req.GetStringTrimmed(r, "menu_style")
	menuRenderingMode := req.GetStringTrimmed(r, "menu_rendering_mode")
	menuCSSClass := req.GetStringTrimmed(r, "menu_css_class")
	menuCSSID := req.GetStringTrimmed(r, "menu_css_id")
	menuStartLevel := req.GetStringTrimmed(r, "menu_start_level")
	menuMaxDepth := req.GetStringTrimmed(r, "menu_max_depth")

	if menuID == "" {
		return fmt.Errorf("menu selection is required")
	}

	block.SetMeta(cmsstore.BLOCK_META_MENU_ID, menuID)
	block.SetMeta(cmsstore.BLOCK_META_MENU_STYLE, menuStyle)
	block.SetMeta(cmsstore.BLOCK_META_MENU_RENDERING_MODE, menuRenderingMode)
	block.SetMeta(cmsstore.BLOCK_META_MENU_CSS_CLASS, menuCSSClass)
	block.SetMeta(cmsstore.BLOCK_META_MENU_CSS_ID, menuCSSID)
	block.SetMeta(cmsstore.BLOCK_META_MENU_START_LEVEL, menuStartLevel)
	block.SetMeta(cmsstore.BLOCK_META_MENU_MAX_DEPTH, menuMaxDepth)

	return nil
}
