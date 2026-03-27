package navbar

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
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

// TypeKey returns the unique identifier for this block type
func (t *NavbarBlockType) TypeKey() string {
	return cmsstore.BLOCK_TYPE_NAVBAR
}

// TypeLabel returns the human-readable display name
func (t *NavbarBlockType) TypeLabel() string {
	return "Navbar"
}

// Render renders the navbar block for frontend display
func (t *NavbarBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	menuID := block.Meta(cmsstore.BLOCK_META_MENU_ID)
	if menuID == "" {
		return "", fmt.Errorf("no menu ID specified for navbar block")
	}

	// Get menu
	_, err := t.store.MenuFindByID(ctx, menuID)
	if err != nil {
		return "", fmt.Errorf("failed to get menu %s: %w", menuID, err)
	}

	// Get menu items
	menuItems, err := t.store.MenuItemList(ctx, cmsstore.MenuItemQuery().
		SetMenuID(menuID).
		SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_SEQUENCE).
		SetSortOrder("asc"))
	if err != nil {
		return "", fmt.Errorf("failed to get menu items: %w", err)
	}

	// Get rendering configuration
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
	if brandURL == "" {
		brandURL = "/"
	}

	fixed := block.Meta(cmsstore.BLOCK_META_NAVBAR_FIXED) == "true"
	dark := block.Meta(cmsstore.BLOCK_META_NAVBAR_DARK) == "true"

	// Use the navbar renderer with unique ID based on block ID
	return renderNavbarHTML(ctx, t.store, block.ID(), menuItems, style, renderingMode, cssClass, cssID, brandText, brandURL, fixed, dark)
}

// GetAdminFields returns form fields for editing navbar block configuration.
func (t *NavbarBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	menuList, err := t.store.MenuList(r.Context(), cmsstore.MenuQuery().
		SetStatus(cmsstore.MENU_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder("asc"))

	if err != nil {
		log.Printf("Error loading menu list for navbar admin fields: %v", err)
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
			Help:     "Select the menu to display in this navbar",
			Options:  menuOptions,
		}),
		form.NewField(form.FieldOptions{
			Label: "Navbar Style",
			Name:  "navbar_style",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_STYLE),
			Help:  "Choose the navbar layout style",
			Options: []form.FieldOption{
				{
					Value: "Default",
					Key:   cmsstore.BLOCK_NAVBAR_STYLE_DEFAULT,
				},
				{
					Value: "Centered",
					Key:   cmsstore.BLOCK_NAVBAR_STYLE_CENTERED,
				},
				{
					Value: "Bottom",
					Key:   cmsstore.BLOCK_NAVBAR_STYLE_BOTTOM,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Rendering Mode",
			Name:  "navbar_rendering_mode",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE),
			Help:  "Choose the rendering framework",
			Options: []form.FieldOption{
				{
					Value: "Plain",
					Key:   cmsstore.BLOCK_NAVBAR_RENDERING_PLAIN,
				},
				{
					Value: "Bootstrap 5",
					Key:   cmsstore.BLOCK_NAVBAR_RENDERING_BOOTSTRAP5,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Brand Text",
			Name:  "navbar_brand_text",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT),
			Help:  "Text displayed as the navbar brand/logo",
		}),
		form.NewField(form.FieldOptions{
			Label: "Brand URL",
			Name:  "navbar_brand_url",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_BRAND_URL),
			Help:  "URL for the brand link (default: /)",
		}),
		form.NewField(form.FieldOptions{
			Label: "CSS ID",
			Name:  "navbar_css_id",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_CSS_ID),
			Help:  "Optional CSS ID for the navbar",
		}),
		form.NewField(form.FieldOptions{
			Label: "CSS Class",
			Name:  "navbar_css_class",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_CSS_CLASS),
			Help:  "Optional CSS classes for styling",
		}),
		form.NewField(form.FieldOptions{
			Label: "Fixed Position",
			Name:  "navbar_fixed",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_FIXED),
			Help:  "Fix the navbar to top or bottom of the page",
			Options: []form.FieldOption{
				{
					Value: "No",
					Key:   "",
				},
				{
					Value: "Top",
					Key:   "true",
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Dark Theme",
			Name:  "navbar_dark",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: block.Meta(cmsstore.BLOCK_META_NAVBAR_DARK),
			Help:  "Use dark color scheme",
			Options: []form.FieldOption{
				{
					Value: "No",
					Key:   "",
				},
				{
					Value: "Yes",
					Key:   "true",
				},
			},
		}),
	}

	return fieldsContent
}

// SaveAdminFields processes form submission and updates the navbar block.
func (t *NavbarBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	r.ParseForm()

	menuID := r.FormValue("menu_id")
	style := r.FormValue("navbar_style")
	renderingMode := r.FormValue("navbar_rendering_mode")
	brandText := r.FormValue("navbar_brand_text")
	brandURL := r.FormValue("navbar_brand_url")
	fixed := r.FormValue("navbar_fixed")
	dark := r.FormValue("navbar_dark")
	cssClass := r.FormValue("navbar_css_class")
	cssID := r.FormValue("navbar_css_id")

	block.SetMeta(cmsstore.BLOCK_META_MENU_ID, menuID)
	block.SetMeta(cmsstore.BLOCK_META_NAVBAR_STYLE, style)
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

	// Validate menu exists
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
