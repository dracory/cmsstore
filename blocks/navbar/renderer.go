package navbar

import (
	"context"
	"fmt"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// renderNavbarHTML renders a navbar with different styles and rendering modes
func renderNavbarHTML(ctx context.Context, store cmsstore.StoreInterface, menuItems []cmsstore.MenuItemInterface, style, renderingMode, cssClass, cssID, brandText, brandURL string, fixed, dark bool) (string, error) {
	// Handle Bootstrap 5 rendering
	if renderingMode == cmsstore.BLOCK_NAVBAR_RENDERING_BOOTSTRAP5 {
		return renderBootstrap5Navbar(menuItems, style, cssClass, cssID, brandText, brandURL, fixed, dark)
	}

	// Handle plain rendering
	return renderPlainNavbar(menuItems, style, cssClass, cssID, brandText, brandURL, fixed, dark)
}

// renderBootstrap5Navbar renders a Bootstrap 5 navbar
func renderBootstrap5Navbar(menuItems []cmsstore.MenuItemInterface, style, cssClass, cssID, brandText, brandURL string, fixed, dark bool) (string, error) {
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

	// Add menu items
	for _, item := range menuItems {
		navItem := hb.Li()
		navItem.Class("nav-item")

		navLink := hb.A()
		navLink.Class("nav-link")
		navLink.Href(item.URL())
		navLink.Text(item.Name())
		
		navItem.AddChild(navLink)
		navbarMenu.AddChild(navItem)
	}

	collapse.AddChild(navbarMenu)
	nav.AddChild(collapse)

	return nav.ToHTML(), nil
}

// renderPlainNavbar renders a plain navbar without Bootstrap classes
func renderPlainNavbar(menuItems []cmsstore.MenuItemInterface, style, cssClass, cssID, brandText, brandURL string, fixed, dark bool) (string, error) {
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

	// Add menu items
	for _, item := range menuItems {
		 menuItem := hb.Li()
	 menuItem.Class("navbar-item")

	 link := hb.A()
	 link.Class("navbar-link")
	 link.Href(item.URL())
	 link.Text(item.Name())
	 
	 menuItem.AddChild(link)
	 menu.AddChild(menuItem)
	}

	nav.AddChild(menu)

	return nav.ToHTML(), nil
}
