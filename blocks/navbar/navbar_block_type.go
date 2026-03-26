package navbar

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dracory/cmsstore"
)

// NavbarBlockType represents a navbar block for navigation
type NavbarBlockType struct {
	store cmsstore.StoreInterface
}

// NewNavbarBlockType creates a new navbar block type
func NewNavbarBlockType(store cmsstore.StoreInterface) *NavbarBlockType {
	return &NavbarBlockType{
		store: store,
	}
}

// Type returns the block type identifier
func (t *NavbarBlockType) Type() string {
	return cmsstore.BLOCK_TYPE_NAVBAR
}

// Name returns the display name for this block type
func (t *NavbarBlockType) Name() string {
	return "Navbar"
}

// Description returns a description of this block type
func (t *NavbarBlockType) Description() string {
	return "Navigation bar with Bootstrap 5 dropdown support and plain rendering options"
}

// Icon returns the icon for this block type
func (t *NavbarBlockType) Icon() string {
	return "🧭"
}

// Category returns the category for this block type
func (t *NavbarBlockType) Category() string {
	return "Navigation"
}

// Render renders the navbar block
func (t *NavbarBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, page cmsstore.PageInterface, site cmsstore.SiteInterface) (string, error) {
	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		return "", fmt.Errorf("no menu ID specified for navbar block")
	}

	// Get menu
	_, err := t.store.MenuFindByID(ctx, menuID)
	if err != nil {
		return "", fmt.Errorf("failed to get menu %s: %w", menuID, err)
	}

	// Get active menu items
	query := cmsstore.MenuItemQuery()
	query.SetMenuID(menuID)
	query.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)

	menuItems, err := t.store.MenuItemList(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to get menu items for menu %s: %w", menuID, err)
	}

	style := block.Meta(cmsstore.BLOCK_META_NAVBAR_STYLE)
	if style == "" {
		style = cmsstore.BLOCK_NAVBAR_STYLE_DEFAULT
	}

	renderingMode := block.Meta(cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE)
	if renderingMode == "" {
		renderingMode = cmsstore.BLOCK_NAVBAR_RENDERING_BOOTSTRAP5
	}

	cssClass := block.Meta(cmsstore.BLOCK_META_NAVBAR_CSS_CLASS)
	cssID := block.Meta(cmsstore.BLOCK_META_NAVBAR_CSS_ID)
	brandText := block.Meta(cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT)
	brandURL := block.Meta(cmsstore.BLOCK_META_NAVBAR_BRAND_URL)
	fixed := block.Meta(cmsstore.BLOCK_META_NAVBAR_FIXED) == "true"
	dark := block.Meta(cmsstore.BLOCK_META_NAVBAR_DARK) == "true"

	return renderNavbarHTML(ctx, t.store, menuItems, style, renderingMode, cssClass, cssID, brandText, brandURL, fixed, dark)
}

// GetAdminFields returns form fields for editing navbar block configuration.
func (t *NavbarBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	fields := map[string]interface{}{
		"menu_id":               block.Meta(cmsstore.BLOCK_META_MENU_ID),
		"navbar_style":          block.Meta(cmsstore.BLOCK_META_NAVBAR_STYLE),
		"navbar_rendering_mode": block.Meta(cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE),
		"navbar_brand_text":     block.Meta(cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT),
		"navbar_brand_url":      block.Meta(cmsstore.BLOCK_META_NAVBAR_BRAND_URL),
		"navbar_fixed":          block.Meta(cmsstore.BLOCK_META_NAVBAR_FIXED),
		"navbar_dark":           block.Meta(cmsstore.BLOCK_META_NAVBAR_DARK),
		"navbar_css_class":      block.Meta(cmsstore.BLOCK_META_NAVBAR_CSS_CLASS),
		"navbar_css_id":         block.Meta(cmsstore.BLOCK_META_NAVBAR_CSS_ID),
	}

	return fields
}

// SaveAdminFields processes form submission and updates the navbar block.
func (t *NavbarBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	r.ParseForm()

	menuID := r.FormValue("menu_id")
	navbarStyle := r.FormValue("navbar_style")
	renderingMode := r.FormValue("navbar_rendering_mode")
	brandText := r.FormValue("navbar_brand_text")
	brandURL := r.FormValue("navbar_brand_url")
	fixed := r.FormValue("navbar_fixed")
	dark := r.FormValue("navbar_dark")
	cssClass := r.FormValue("navbar_css_class")
	cssID := r.FormValue("navbar_css_id")

	if menuID == "" {
		return fmt.Errorf("menu selection is required")
	}

	block.SetMeta(cmsstore.BLOCK_META_MENU_ID, menuID)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_STYLE, navbarStyle)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE, renderingMode)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT, brandText)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_BRAND_URL, brandURL)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_FIXED, fixed)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_DARK, dark)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_CSS_CLASS, cssClass)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_CSS_ID, cssID)

	return nil
}

// Validate validates the navbar block configuration
func (t *NavbarBlockType) Validate(block cmsstore.BlockInterface) error {
	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		return fmt.Errorf("menu ID is required for navbar block")
	}

	// Validate that the menu exists
	_, err := t.store.MenuFindByID(context.Background(), menuID)
	if err != nil {
		return fmt.Errorf("invalid menu ID: %s", menuID)
	}

	return nil
}

// GetPreview returns a preview of the navbar block
func (t *NavbarBlockType) GetPreview(block cmsstore.BlockInterface) string {
	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		return "No menu selected"
	}

	menu, err := t.store.MenuFindByID(context.Background(), menuID)
	if err != nil || menu == nil {
		return "Invalid menu"
	}

	return fmt.Sprintf("Navbar: %s", menu.Name())
}
