package admin

import (
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
)

// MenuAdminProvider provides admin UI for Menu block types.
type MenuAdminProvider struct {
	store  cmsstore.StoreInterface
	logger interface{ Error(msg string, args ...interface{}) }
}

// NewMenuAdminProvider creates a new Menu block admin provider.
func NewMenuAdminProvider(store cmsstore.StoreInterface, logger interface{ Error(msg string, args ...interface{}) }) *MenuAdminProvider {
	return &MenuAdminProvider{
		store:  store,
		logger: logger,
	}
}

// GetContentFields returns form fields for Menu block content editing.
func (p *MenuAdminProvider) GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface {
	menusEnabled := p.store.MenusEnabled()

	if !menusEnabled {
		return []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Menu Blocks Not Available",
				Type:  form.FORM_FIELD_TYPE_RAW,
				Value: hb.Div().Class("alert alert-warning").Text("Menu functionality is not enabled in this CMS installation.").ToHTML(),
			}),
		}
	}

	menuList, err := p.store.MenuList(r.Context(), cmsstore.MenuQuery().
		SetStatus(cmsstore.MENU_STATUS_ACTIVE).
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC))

	if err != nil {
		p.logger.Error("Error loading menus", "error", err.Error())
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

// GetTypeLabel returns the display label for Menu blocks.
func (p *MenuAdminProvider) GetTypeLabel() string {
	return "Menu Block"
}

// SaveContentFields processes form data and updates the Menu block.
func (p *MenuAdminProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
	menuID := req.GetStringTrimmed(r, "menu_id")
	menuStyle := req.GetStringTrimmed(r, "menu_style")
	menuCSSClass := req.GetStringTrimmed(r, "menu_css_class")
	menuStartLevel := req.GetStringTrimmed(r, "menu_start_level")
	menuMaxDepth := req.GetStringTrimmed(r, "menu_max_depth")

	if menuID == "" {
		return &ValidationError{Message: "Menu selection is required"}
	}

	block.SetMeta(cmsstore.BLOCK_META_MENU_ID, menuID)
	block.SetMeta(cmsstore.BLOCK_META_MENU_STYLE, menuStyle)
	block.SetMeta(cmsstore.BLOCK_META_MENU_CSS_CLASS, menuCSSClass)
	block.SetMeta(cmsstore.BLOCK_META_MENU_START_LEVEL, menuStartLevel)
	block.SetMeta(cmsstore.BLOCK_META_MENU_MAX_DEPTH, menuMaxDepth)

	return nil
}

// ValidationError represents a validation error from SaveContentFields.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
