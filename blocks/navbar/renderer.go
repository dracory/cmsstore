package navbar

import (
	"context"
	"fmt"
	"strings"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// renderNavbarHTML renders a navbar with different styles and rendering modes
func renderNavbarHTML(ctx context.Context, store cmsstore.StoreInterface, menuItems []cmsstore.MenuItemInterface, style, renderingMode, cssClass, cssID, brandText, brandURL string, fixed, dark bool) (string, error) {
	// Handle Bootstrap 5 rendering
	if renderingMode == cmsstore.BLOCK_NAVBAR_RENDERING_BOOTSTRAP5 {
		return renderBootstrap5Navbar(ctx, store, menuItems, style, cssClass, cssID, brandText, brandURL, fixed, dark)
	}

	// Handle plain rendering
	return renderPlainNavbar(ctx, store, menuItems, style, cssClass, cssID, brandText, brandURL, fixed, dark)
}

// renderBootstrap5Navbar renders a Bootstrap 5 navbar
func renderBootstrap5Navbar(ctx context.Context, store cmsstore.StoreInterface, menuItems []cmsstore.MenuItemInterface, style, cssClass, cssID, brandText, brandURL string, fixed, dark bool) (string, error) {
	nav := hb.Nav()

	// Base Bootstrap navbar classes
	nav.Class("navbar")
	nav.Class("navbar-expand-lg")

	// Add theme class
	if dark {
		nav.Class("navbar-dark")
		nav.Class("bg-dark")
	} else {
		nav.Class("navbar-light")
		nav.Class("bg-light")
	}

	// Add fixed positioning
	if fixed {
		nav.Class("fixed-top")
	}

	// Add custom CSS classes
	if cssClass != "" {
		nav.Class(cssClass)
	}

	// Add CSS ID
	if cssID != "" {
		nav.ID(cssID)
	}

	// Create brand
	if brandText != "" {
		brand := hb.A()
		brand.Class("navbar-brand")
		if brandURL != "" {
			brand.Href(brandURL)
		} else {
			brand.Href("/")
		}
		brand.Text(brandText)
		nav.AddChild(brand)
	}

	// Create navbar toggler for mobile
	toggler := hb.Button()
	toggler.Class("navbar-toggler")
	toggler.Attr("type", "button")
	toggler.Attr("data-bs-toggle", "collapse")
	toggler.Attr("data-bs-target", "#navbarContent")
	toggler.Attr("aria-controls", "navbarContent")
	toggler.Attr("aria-expanded", "false")
	toggler.Attr("aria-label", "Toggle navigation")

	// Add hamburger icon
	span1 := hb.Span()
	span1.Class("navbar-toggler-icon")
	toggler.AddChild(span1)
	nav.AddChild(toggler)

	// Create collapsible content
	collapse := hb.Div()
	collapse.Class("collapse")
	collapse.Class("navbar-collapse")
	collapse.ID("navbarContent")

	// Create navbar menu
	navbarMenu := hb.Ul()
	navbarMenu.Class("navbar-nav")
	if style == cmsstore.BLOCK_NAVBAR_STYLE_CENTERED {
		navbarMenu.Class("justify-content-center")
	} else if style == cmsstore.BLOCK_NAVBAR_STYLE_BOTTOM {
		navbarMenu.Class("ms-auto")
	}

	// Add menu items (flat list for now)
	for _, item := range menuItems {
		url := resolveMenuItemURL(ctx, store, item)

		navItem := hb.Li()
		navItem.Class("nav-item")

		if url != "" {
			navLink := hb.A()
			navLink.Class("nav-link")
			navLink.Href(url)
			navLink.Text(item.Name())

			// Add target attribute if set
			if item.Target() != "" {
				navLink.Attr("target", item.Target())
			}

			navItem.AddChild(navLink)
		} else {
			// Render as plain text without link
			navItem.Text(item.Name())
		}

		navbarMenu.AddChild(navItem)
	}

	collapse.AddChild(navbarMenu)
	nav.AddChild(collapse)

	return nav.ToHTML(), nil
}

// renderPlainNavbar renders a plain navbar without Bootstrap classes
func renderPlainNavbar(ctx context.Context, store cmsstore.StoreInterface, menuItems []cmsstore.MenuItemInterface, style, cssClass, cssID, brandText, brandURL string, fixed, dark bool) (string, error) {
	nav := hb.Nav()

	// Base classes
	nav.Class("navbar")
	nav.Class(fmt.Sprintf("navbar-style-%s", style))

	// Add fixed positioning
	if fixed {
		nav.Class("navbar-fixed")
	}

	// Add theme class
	if dark {
		nav.Class("navbar-dark")
	}

	// Add custom CSS classes
	if cssClass != "" {
		nav.Class(cssClass)
	}

	// Add CSS ID
	if cssID != "" {
		nav.ID(cssID)
	}

	// Create brand
	if brandText != "" {
		brand := hb.A()
		brand.Class("navbar-brand")
		if brandURL != "" {
			brand.Href(brandURL)
		} else {
			brand.Href("/")
		}
		brand.Text(brandText)
		nav.AddChild(brand)
	}

	// Create menu
	menu := hb.Ul()
	menu.Class("navbar-menu")
	if style == cmsstore.BLOCK_NAVBAR_STYLE_CENTERED {
		menu.Class("navbar-centered")
	} else if style == cmsstore.BLOCK_NAVBAR_STYLE_BOTTOM {
		menu.Class("navbar-bottom")
	}

	// Add menu items (flat list)
	for _, item := range menuItems {
		url := resolveMenuItemURL(ctx, store, item)

		menuItem := hb.Li()
		menuItem.Class("navbar-item")

		if url != "" {
			link := hb.A()
			link.Class("navbar-link")
			link.Href(url)
			link.Text(item.Name())

			// Add target attribute if set
			if item.Target() != "" {
				link.Attr("target", item.Target())
			}

			menuItem.AddChild(link)
		} else {
			// Render as plain text without link
			menuItem.Text(item.Name())
		}

		menu.AddChild(menuItem)
	}

	nav.AddChild(menu)

	return nav.ToHTML(), nil
}

// resolveMenuItemURL resolves the URL for a menu item
// If the item has a direct URL, it returns that.
// Otherwise, if the item has a PageID, it looks up the page and returns its alias.
func resolveMenuItemURL(ctx context.Context, store cmsstore.StoreInterface, item cmsstore.MenuItemInterface) string {
	if item.URL() != "" {
		return item.URL()
	}

	if item.PageID() != "" {
		page, err := store.PageFindByID(ctx, item.PageID())
		if err != nil || page == nil {
			return ""
		}
		return "/" + strings.TrimPrefix(page.Alias(), "/")
	}

	return ""
}
