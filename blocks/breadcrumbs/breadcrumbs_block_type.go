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
	breadcrumbs := t.generateBreadcrumbs(ctx, homeText, homeURL)

	// Use the breadcrumbs renderer
	return renderBreadcrumbsHTML(breadcrumbs, style, renderingMode, cssClass, cssID, separator)
}

// GetAdminFields returns form fields for editing breadcrumbs block configuration.
func (t *BreadcrumbsBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	fields := map[string]interface{}{
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

	style := r.FormValue("breadcrumbs_style")
	renderingMode := r.FormValue("breadcrumbs_rendering_mode")
	separator := r.FormValue("breadcrumbs_separator")
	homeText := r.FormValue("breadcrumbs_home_text")
	homeURL := r.FormValue("breadcrumbs_home_url")
	cssClass := r.FormValue("breadcrumbs_css_class")
	cssID := r.FormValue("breadcrumbs_css_id")

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
func (t *BreadcrumbsBlockType) generateBreadcrumbs(ctx context.Context, homeText, homeURL string) []BreadcrumbItem {
	var breadcrumbs []BreadcrumbItem

	// Add home breadcrumb
	breadcrumbs = append(breadcrumbs, BreadcrumbItem{
		Name:   homeText,
		URL:    homeURL,
		Active: false,
	})

	// Add current page breadcrumb (simplified for now)
	breadcrumbs = append(breadcrumbs, BreadcrumbItem{
		Name:   "Current Page",
		URL:    "", // Current page has no URL
		Active: true,
	})

	return breadcrumbs
}

// BreadcrumbItem represents a single breadcrumb item
type BreadcrumbItem struct {
	Name   string
	URL    string
	Active bool
}
