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
func renderMenuHTML(ctx context.Context, store cmsstore.StoreInterface, menuItems []cmsstore.MenuItemInterface, style, cssClass, cssID string, startLevel, maxDepth int) (string, error) {
	// Build nav element using hb library
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
