package breadcrumbs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dracory/cmsstore"
)

// BreadcrumbsBlockType represents a breadcrumbs block for navigation
type BreadcrumbsBlockType struct {
	store cmsstore.StoreInterface
}

// NewBreadcrumbsBlockType creates a new breadcrumbs block type
func NewBreadcrumbsBlockType(store cmsstore.StoreInterface) *BreadcrumbsBlockType {
	return &BreadcrumbsBlockType{
		store: store,
	}
}

// TypeKey returns the unique identifier for this block type
func (t *BreadcrumbsBlockType) TypeKey() string {
	return cmsstore.BLOCK_TYPE_BREADCRUMBS
}

// TypeLabel returns the human-readable display name
func (t *BreadcrumbsBlockType) TypeLabel() string {
	return "Breadcrumbs"
}

// Render renders the breadcrumbs block for frontend display
func (t *BreadcrumbsBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	style := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE)
	if style == "" {
		style = cmsstore.BLOCK_BREADCRUMBS_STYLE_DEFAULT
	}

	renderingMode := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE)
	if renderingMode == "" {
		renderingMode = cmsstore.BLOCK_BREADCRUMBS_RENDERING_BOOTSTRAP5
	}

	cssClass := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS)
	cssID := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_ID)
	separator := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR)
	if separator == "" {
		separator = "/"
	}

	homeText := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_TEXT)
	if homeText == "" {
		homeText = "Home"
	}

	homeURL := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_URL)
	if homeURL == "" {
		homeURL = "/"
	}

	// Generate breadcrumb items based on current page
	breadcrumbs := t.generateBreadcrumbs(ctx, block, homeText, homeURL)

	// Use the breadcrumbs renderer
	return renderBreadcrumbsHTML(breadcrumbs, style, renderingMode, cssClass, cssID, separator)
}

// GetAdminFields returns form fields for editing breadcrumbs block configuration.
func (t *BreadcrumbsBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	fields := map[string]interface{}{
		"breadcrumbs_menu_id":        block.Meta(cmsstore.BLOCK_META_MENU_ID),
		"breadcrumbs_style":          block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE),
		"breadcrumbs_rendering_mode": block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE),
		"breadcrumbs_separator":      block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR),
		"breadcrumbs_home_text":      block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_TEXT),
		"breadcrumbs_home_url":       block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_URL),
		"breadcrumbs_css_class":      block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS),
		"breadcrumbs_css_id":         block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_ID),
	}

	return fields
}

// SaveAdminFields processes form submission and updates the breadcrumbs block.
func (t *BreadcrumbsBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	r.ParseForm()

	menuID := r.FormValue("breadcrumbs_menu_id")
	style := r.FormValue("breadcrumbs_style")
	renderingMode := r.FormValue("breadcrumbs_rendering_mode")
	separator := r.FormValue("breadcrumbs_separator")
	homeText := r.FormValue("breadcrumbs_home_text")
	homeURL := r.FormValue("breadcrumbs_home_url")
	cssClass := r.FormValue("breadcrumbs_css_class")
	cssID := r.FormValue("breadcrumbs_css_id")

	block.SetMeta(cmsstore.BLOCK_META_MENU_ID, menuID)
	block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE, style)
	block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE, renderingMode)
	block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR, separator)
	block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_TEXT, homeText)
	block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_URL, homeURL)
	block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS, cssClass)
	block.SetMeta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_ID, cssID)

	return nil
}

// Validate validates the breadcrumbs block configuration
func (t *BreadcrumbsBlockType) Validate(block cmsstore.BlockInterface) error {
	// Breadcrumbs don't require any specific configuration
	return nil
}

// GetPreview returns a preview of the breadcrumbs block
func (t *BreadcrumbsBlockType) GetPreview(block cmsstore.BlockInterface) string {
	style := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE)
	renderingMode := block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE)

	if style == "" {
		style = "default"
	}
	if renderingMode == "" {
		renderingMode = "bootstrap5"
	}

	return fmt.Sprintf("Breadcrumbs: %s (%s)", style, renderingMode)
}

// generateBreadcrumbs creates breadcrumb items based on the current page context
func (t *BreadcrumbsBlockType) generateBreadcrumbs(ctx context.Context, block cmsstore.BlockInterface, homeText, homeURL string) []BreadcrumbItem {
	var breadcrumbs []BreadcrumbItem

	// Add home breadcrumb
	breadcrumbs = append(breadcrumbs, BreadcrumbItem{
		Name:   homeText,
		URL:    homeURL,
		Active: false,
	})

	// Get current page from context
	currentPage, found := getCurrentPageFromContext(ctx)
	if !found {
		// If no current page found, return only home breadcrumb
		return breadcrumbs
	}

	// Get menu ID from block configuration
	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		// If no menu configured, return home + current page
		breadcrumbs = append(breadcrumbs, BreadcrumbItem{
			Name:   currentPage.Name(),
			URL:    "", // Current page has no URL
			Active: true,
		})
		return breadcrumbs
	}

	// Build breadcrumb path from menu hierarchy
	menuPath, err := t.buildMenuPath(ctx, menuID, currentPage.ID())
	if err != nil || len(menuPath) == 0 {
		// Fallback to home + current page if menu navigation fails
		breadcrumbs = append(breadcrumbs, BreadcrumbItem{
			Name:   currentPage.Name(),
			URL:    "", // Current page has no URL
			Active: true,
		})
		return breadcrumbs
	}

	// Add menu path items (excluding home which is already added)
	for i, item := range menuPath {
		isActive := i == len(menuPath)-1 // Last item is active
		breadcrumbs = append(breadcrumbs, BreadcrumbItem{
			Name:   item.Name,
			URL:    item.URL,
			Active: isActive,
		})
	}

	return breadcrumbs
}

// getCurrentPageFromContext extracts the current page from the context
func getCurrentPageFromContext(ctx context.Context) (cmsstore.PageInterface, bool) {
	type contextKey string
	const pageContextKey contextKey = "page"

	if page, ok := ctx.Value(pageContextKey).(cmsstore.PageInterface); ok {
		return page, true
	}
	return nil, false
}

// MenuPathItem represents an item in the menu path
type MenuPathItem struct {
	Name string
	URL  string
}

// buildMenuPath builds the breadcrumb path from menu hierarchy
func (t *BreadcrumbsBlockType) buildMenuPath(ctx context.Context, menuID, currentPageID string) ([]MenuPathItem, error) {
	// Get all menu items for the specified menu
	menuItems, err := t.store.MenuItemList(ctx, cmsstore.MenuItemQuery().
		SetMenuID(menuID).
		SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE))
	if err != nil {
		return nil, err
	}

	// Create a map of menu items by ID for quick lookup
	itemMap := make(map[string]cmsstore.MenuItemInterface)
	for _, item := range menuItems {
		itemMap[item.ID()] = item
	}

	// Find the current page in the menu
	var currentItem cmsstore.MenuItemInterface
	for _, item := range menuItems {
		if item.PageID() == currentPageID {
			currentItem = item
			break
		}
	}

	if currentItem == nil {
		return nil, fmt.Errorf("current page not found in menu")
	}

	// Build the path by walking up the hierarchy
	var path []MenuPathItem
	current := currentItem

	// Walk up to root
	for current != nil {
		// Get page details for this menu item
		page, err := t.store.PageFindByID(ctx, current.PageID())
		if err != nil {
			return nil, err
		}

		// Determine URL for this item
		var url string
		if current.URL() != "" {
			url = current.URL()
		} else {
			// Use page alias if no custom URL
			url = "/" + page.Alias()
		}

		// Add to beginning of path (reverse order)
		path = append([]MenuPathItem{{
			Name: page.Name(),
			URL:  url,
		}}, path...)

		// Move to parent
		parentID := current.ParentID()
		if parentID == "" {
			break
		}
		current = itemMap[parentID]
	}

	return path, nil
}

// BreadcrumbItem represents a single breadcrumb item
type BreadcrumbItem struct {
	Name   string
	URL    string
	Active bool
}
