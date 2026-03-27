package navbar

import (
	"context"
	"fmt"
	"strings"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// renderNavbarHTML renders a navbar with different styles and rendering modes
func renderNavbarHTML(ctx context.Context, store cmsstore.StoreInterface, blockID string, menuItems []cmsstore.MenuItemInterface, style, renderingMode, cssClass, cssID, brandText, brandURL, brandImageURL, brandImageWidth, brandImageHeight, brandImageAlt string, fixed, dark bool, customCSS string) (string, error) {
	// Handle Bootstrap 5 rendering
	if renderingMode == cmsstore.BLOCK_NAVBAR_RENDERING_BOOTSTRAP5 {
		return renderBootstrap5Navbar(ctx, store, blockID, menuItems, style, cssClass, cssID, brandText, brandURL, brandImageURL, brandImageWidth, brandImageHeight, brandImageAlt, fixed, dark, customCSS)
	}

	// Handle plain rendering
	return renderPlainNavbar(ctx, store, blockID, menuItems, style, cssClass, cssID, brandText, brandURL, brandImageURL, brandImageWidth, brandImageHeight, brandImageAlt, fixed, dark, customCSS)
}

// renderBootstrap5Navbar renders a Bootstrap 5 navbar
func renderBootstrap5Navbar(ctx context.Context, store cmsstore.StoreInterface, blockID string, menuItems []cmsstore.MenuItemInterface, style, cssClass, cssID, brandText, brandURL, brandImageURL, brandImageWidth, brandImageHeight, brandImageAlt string, fixed, dark bool, customCSS string) (string, error) {
	var result strings.Builder

	// Add custom CSS if provided
	if customCSS != "" {
		result.WriteString("<style>")
		result.WriteString(customCSS)
		result.WriteString("</style>")
	}

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

	// Create brand with image and/or text support
	if brandText != "" || brandImageURL != "" {
		brand := hb.A()
		brand.Class("navbar-brand")
		if brandURL != "" {
			brand.Href(brandURL)
		} else {
			brand.Href("/")
		}

		// Determine what to render: image only, text only, or both
		hasImage := brandImageURL != ""
		hasText := brandText != ""

		if hasImage && hasText {
			// Image and text together (Bootstrap 5 pattern)
			img := hb.Img(brandImageURL)

			// Set default dimensions if not provided
			width := brandImageWidth
			if width == "" {
				width = "30"
			}
			height := brandImageHeight
			if height == "" {
				height = "24"
			}

			img.Attr("width", width)
			img.Attr("height", height)
			img.Class("d-inline-block align-text-top")

			alt := brandImageAlt
			if alt == "" {
				alt = "Logo"
			}
			img.Alt(alt)

			brand.AddChild(img)
			brand.Text(" " + brandText) // Add space between image and text
		} else if hasImage {
			// Image only
			img := hb.Img(brandImageURL)

			// Set default dimensions if not provided
			width := brandImageWidth
			if width == "" {
				width = "30"
			}
			height := brandImageHeight
			if height == "" {
				height = "24"
			}

			img.Attr("width", width)
			img.Attr("height", height)

			alt := brandImageAlt
			if alt == "" {
				alt = "Logo"
			}
			img.Alt(alt)

			brand.AddChild(img)
		} else {
			// Text only
			brand.Text(brandText)
		}

		nav.AddChild(brand)
	}

	// Create navbar toggler for mobile with unique target
	contentID := "navbarContent-" + blockID
	toggler := hb.Button()
	toggler.Class("navbar-toggler")
	toggler.Attr("type", "button")
	toggler.Attr("data-bs-toggle", "collapse")
	toggler.Attr("data-bs-target", "#"+contentID)
	toggler.Attr("aria-controls", contentID)
	toggler.Attr("aria-expanded", "false")
	toggler.Attr("aria-label", "Toggle navigation")

	// Add hamburger icon
	span1 := hb.Span()
	span1.Class("navbar-toggler-icon")
	toggler.AddChild(span1)
	nav.AddChild(toggler)

	// Create collapsible content with unique ID
	collapse := hb.Div()
	collapse.Class("collapse")
	collapse.Class("navbar-collapse")
	collapse.ID(contentID)

	// Create navbar menu with hierarchical support
	navbarMenu := hb.Ul()
	navbarMenu.Class("navbar-nav")
	if style == cmsstore.BLOCK_NAVBAR_STYLE_CENTERED {
		navbarMenu.Class("justify-content-center")
	} else if style == cmsstore.BLOCK_NAVBAR_STYLE_BOTTOM {
		navbarMenu.Class("ms-auto")
	}

	// Build menu item lookup map for parent-child relationships
	menuItemMap := make(map[string]cmsstore.MenuItemInterface)
	for _, item := range menuItems {
		menuItemMap[item.ID()] = item
	}

	// Find top-level items (no parent or empty parent ID)
	var topLevelItems []cmsstore.MenuItemInterface
	for _, item := range menuItems {
		if item.ParentID() == "" {
			topLevelItems = append(topLevelItems, item)
		}
	}

	// Render top-level items with dropdown support
	for _, item := range topLevelItems {
		navItem := renderNavItemWithDropdown(ctx, store, item, menuItemMap, 0)
		navbarMenu.AddChild(navItem)
	}

	collapse.AddChild(navbarMenu)
	nav.AddChild(collapse)

	result.WriteString(nav.ToHTML())
	return result.String(), nil
}

// renderNavItemWithDropdown renders a nav item with dropdown support for children
func renderNavItemWithDropdown(ctx context.Context, store cmsstore.StoreInterface, item cmsstore.MenuItemInterface, menuItemMap map[string]cmsstore.MenuItemInterface, depth int) *hb.Tag {
	url := resolveMenuItemURL(ctx, store, item)

	// Find children of this item
	var children []cmsstore.MenuItemInterface
	for _, mi := range menuItemMap {
		if mi.ParentID() == item.ID() {
			children = append(children, mi)
		}
	}

	hasChildren := len(children) > 0

	if hasChildren {
		// Render as dropdown
		navItem := hb.Li()
		navItem.Class("nav-item")
		navItem.Class("dropdown")

		// Dropdown toggle link
		dropdownToggle := hb.A()
		dropdownToggle.Class("nav-link")
		dropdownToggle.Class("dropdown-toggle")
		dropdownToggle.Href("#")
		dropdownToggle.Attr("role", "button")
		dropdownToggle.Attr("data-bs-toggle", "dropdown")
		dropdownToggle.Attr("aria-expanded", "false")
		dropdownToggle.Text(item.Name())

		// Add target attribute if set
		if item.Target() != "" {
			dropdownToggle.Attr("target", item.Target())
		}

		navItem.AddChild(dropdownToggle)

		// Dropdown menu
		dropdownMenu := hb.Ul()
		dropdownMenu.Class("dropdown-menu")

		// Add child items to dropdown
		for _, child := range children {
			childURL := resolveMenuItemURL(ctx, store, child)
			dropdownItem := hb.Li()

			childLink := hb.A()
			childLink.Class("dropdown-item")
			childLink.Text(child.Name())

			if childURL != "" {
				childLink.Href(childURL)
				if child.Target() != "" {
					childLink.Attr("target", child.Target())
				}
			} else {
				childLink.Href("#")
			}

			dropdownItem.AddChild(childLink)
			dropdownMenu.AddChild(dropdownItem)
		}

		navItem.AddChild(dropdownMenu)
		return navItem
	}

	// Render as simple nav item (no children)
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

	return navItem
}

// renderPlainNavbar renders a plain navbar without Bootstrap classes
func renderPlainNavbar(ctx context.Context, store cmsstore.StoreInterface, blockID string, menuItems []cmsstore.MenuItemInterface, style, cssClass, cssID, brandText, brandURL, brandImageURL, brandImageWidth, brandImageHeight, brandImageAlt string, fixed, dark bool, customCSS string) (string, error) {
	var result strings.Builder

	// Add custom CSS if provided
	if customCSS != "" {
		result.WriteString("<style>")
		result.WriteString(customCSS)
		result.WriteString("</style>")
	}

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

	// Create brand with image and/or text support
	if brandText != "" || brandImageURL != "" {
		brand := hb.A()
		brand.Class("navbar-brand")
		if brandURL != "" {
			brand.Href(brandURL)
		} else {
			brand.Href("/")
		}

		// Determine what to render: image only, text only, or both
		hasImage := brandImageURL != ""
		hasText := brandText != ""

		if hasImage && hasText {
			// Image and text together
			img := hb.Img(brandImageURL)

			// Set default dimensions if not provided
			width := brandImageWidth
			if width == "" {
				width = "30"
			}
			height := brandImageHeight
			if height == "" {
				height = "24"
			}

			img.Attr("width", width)
			img.Attr("height", height)

			alt := brandImageAlt
			if alt == "" {
				alt = "Logo"
			}
			img.Alt(alt)

			brand.AddChild(img)
			brand.Text(" " + brandText) // Add space between image and text
		} else if hasImage {
			// Image only
			img := hb.Img(brandImageURL)

			// Set default dimensions if not provided
			width := brandImageWidth
			if width == "" {
				width = "30"
			}
			height := brandImageHeight
			if height == "" {
				height = "24"
			}

			img.Attr("width", width)
			img.Attr("height", height)

			alt := brandImageAlt
			if alt == "" {
				alt = "Logo"
			}
			img.Alt(alt)

			brand.AddChild(img)
		} else {
			// Text only
			brand.Text(brandText)
		}

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
	result.WriteString(nav.ToHTML())
	return result.String(), nil
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
