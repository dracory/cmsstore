package breadcrumbs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
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
// Supports runtime attributes for dynamic configuration
func (t *BreadcrumbsBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
	// Parse render options
	options := &cmsstore.RenderOptions{}
	for _, opt := range opts {
		opt(options)
	}

	style := options.Attributes["style"]
	if style == "" {
		style = block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE)
	}
	if style == "" {
		style = cmsstore.BLOCK_BREADCRUMBS_STYLE_DEFAULT
	}

	renderingMode := options.Attributes["mode"]
	if renderingMode == "" {
		renderingMode = block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE)
	}
	if renderingMode == "" {
		renderingMode = cmsstore.BLOCK_BREADCRUMBS_RENDERING_BOOTSTRAP5
	}

	cssClass := options.Attributes["class"]
	if cssClass == "" {
		cssClass = block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS)
	}

	cssID := options.Attributes["id"]
	if cssID == "" {
		cssID = block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_ID)
	}

	separator := options.Attributes["separator"]
	if separator == "" {
		separator = block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR)
	}
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
	menuList, err := t.store.MenuList(r.Context(), cmsstore.MenuQuery().
		SetStatus(cmsstore.MENU_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder("asc"))

	if err != nil {
		// Continue with empty menu list
	}

	menuOptions := []form.FieldOption{
		{
			Value: "- No Menu -",
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
			Label:   "Menu (optional)",
			Name:    "breadcrumbs_menu_id",
			Type:    form.FORM_FIELD_TYPE_SELECT,
			Value:   block.Meta(cmsstore.BLOCK_META_MENU_ID),
			Help:    "Select menu to generate breadcrumbs from page hierarchy (optional)",
			Options: menuOptions,
		}),
		form.NewField(form.FieldOptions{
			Label: "Breadcrumbs Style",
			Name:  "breadcrumbs_style",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_STYLE),
			Help:  "Choose the breadcrumbs layout style",
			Options: []form.FieldOption{
				{
					Value: "Default",
					Key:   cmsstore.BLOCK_BREADCRUMBS_STYLE_DEFAULT,
				},
				{
					Value: "Centered",
					Key:   cmsstore.BLOCK_BREADCRUMBS_STYLE_CENTERED,
				},
				{
					Value: "Right",
					Key:   cmsstore.BLOCK_BREADCRUMBS_STYLE_RIGHT,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Rendering Mode",
			Name:  "breadcrumbs_rendering_mode",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE),
			Help:  "Choose the rendering framework",
			Options: []form.FieldOption{
				{
					Value: "Plain",
					Key:   cmsstore.BLOCK_BREADCRUMBS_RENDERING_PLAIN,
				},
				{
					Value: "Bootstrap 5",
					Key:   cmsstore.BLOCK_BREADCRUMBS_RENDERING_BOOTSTRAP5,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Home Text",
			Name:  "breadcrumbs_home_text",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_TEXT),
			Help:  "Text for the home breadcrumb (default: Home)",
		}),
		form.NewField(form.FieldOptions{
			Label: "Home URL",
			Name:  "breadcrumbs_home_url",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_HOME_URL),
			Help:  "URL for the home link (default: /)",
		}),
		form.NewField(form.FieldOptions{
			Label: "Separator",
			Name:  "breadcrumbs_separator",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR),
			Help:  "Separator between breadcrumbs (default: /)",
		}),
		form.NewField(form.FieldOptions{
			Label: "CSS ID",
			Name:  "breadcrumbs_css_id",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_ID),
			Help:  "Optional CSS ID for the breadcrumbs container",
		}),
		form.NewField(form.FieldOptions{
			Label: "CSS Class",
			Name:  "breadcrumbs_css_class",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS),
			Help:  "Optional CSS classes for styling",
		}),
	}

	return fieldsContent
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
