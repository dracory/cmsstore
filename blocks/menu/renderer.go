package menu

import (
	"context"
	"fmt"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// renderMenuHTML provides a simplified menu rendering implementation.
//
// IMPORTANT: This is a simplified implementation that renders menus as flat lists.
// It does NOT support:
//   - Hierarchical/nested menu structures
//   - startLevel and maxDepth filtering
//   - Different style rendering (vertical, horizontal, dropdown, breadcrumb)
//
// For production use, you should:
//  1. Import and use the comprehensive MenuRenderer from frontend/blocks/menu package
//  2. Or implement full hierarchical rendering here
//
// This simplified version is provided to avoid circular dependencies between
// the blocks package and the frontend package, while still allowing basic
// menu block functionality.
func renderMenuHTML(ctx context.Context, store cmsstore.StoreInterface, menuItems []cmsstore.MenuItemInterface, style, renderingMode, cssClass, cssID string, startLevel, maxDepth int) (string, error) {
	// Handle Bootstrap 5 dropdown separately as it has a different structure
	if renderingMode == cmsstore.BLOCK_MENU_RENDERING_BOOTSTRAP5 {
		return renderBootstrap5Dropdown(menuItems, cssClass, cssID)
	}

	// Build nav element using hb library for other styles
	nav := hb.Nav()

	// Add CSS classes
	nav.Class(fmt.Sprintf("menu menu-style-%s", style))
	if cssClass != "" {
		nav.Class(cssClass)
	}

	// Add CSS ID if provided
	if cssID != "" {
		nav.ID(cssID)
	}

	// Add menu items
	for _, item := range menuItems {
		nav.AddChild(hb.A().Href(item.URL()).Text(item.Name()))
	}

	return nav.ToHTML(), nil
}

// renderBootstrap5Dropdown renders a Bootstrap 5 dropdown menu
func renderBootstrap5Dropdown(menuItems []cmsstore.MenuItemInterface, cssClass, cssID string) (string, error) {
	// Build Bootstrap 5 dropdown structure
	div := hb.Div()

	// Add CSS classes
	div.Class("dropdown")
	if cssClass != "" {
		div.Class(cssClass)
	}

	// Add CSS ID if provided
	if cssID != "" {
		div.ID(cssID)
	}

	// Create dropdown toggle button
	button := hb.Button()
	button.Class("btn btn-secondary dropdown-toggle")
	button.Attr("type", "button")
	button.Attr("data-bs-toggle", "dropdown")
	button.Attr("aria-expanded", "false")
	button.Text("Dropdown")

	// Create dropdown menu
	dropdownMenu := hb.Div()
	dropdownMenu.Class("dropdown-menu")

	// Add menu items as dropdown items
	for _, item := range menuItems {
		dropdownItem := hb.A()
		dropdownItem.Class("dropdown-item")
		dropdownItem.Href(item.URL())
		dropdownItem.Text(item.Name())
		dropdownMenu.AddChild(dropdownItem)
	}

	// Assemble the dropdown
	div.AddChild(button)
	div.AddChild(dropdownMenu)

	return div.ToHTML(), nil
}
