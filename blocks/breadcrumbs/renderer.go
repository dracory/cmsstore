package breadcrumbs

import (
	"fmt"

	"github.com/dracory/cmsstore"
	"github.com/dracory/hb"
)

// renderBreadcrumbsHTML renders breadcrumbs with different styles and rendering modes
func renderBreadcrumbsHTML(breadcrumbs []BreadcrumbItem, style, renderingMode, cssClass, cssID, separator string) (string, error) {
	// Handle Bootstrap 5 rendering
	if renderingMode == cmsstore.BLOCK_BREADCRUMBS_RENDERING_BOOTSTRAP5 {
		return renderBootstrap5Breadcrumbs(breadcrumbs, style, cssClass, cssID, separator)
	}

	// Handle plain rendering
	return renderPlainBreadcrumbs(breadcrumbs, style, cssClass, cssID, separator)
}

// renderBootstrap5Breadcrumbs renders Bootstrap 5 breadcrumbs
func renderBootstrap5Breadcrumbs(breadcrumbs []BreadcrumbItem, style, cssClass, cssID, separator string) (string, error) {
	nav := hb.Nav()

	// Base Bootstrap breadcrumb classes
	nav.Class("breadcrumb")
	nav.Attr("aria-label", "breadcrumb")

	// Add custom CSS classes
	if cssClass != "" {
		nav.Class(cssClass)
	}

	// Add CSS ID
	if cssID != "" {
		nav.ID(cssID)
	}

	// Add breadcrumb items
	for _, item := range breadcrumbs {
		if item.Active {
			// Active breadcrumb (current page)
			active := hb.Span()
			active.Class("breadcrumb-item active")
			active.Attr("aria-current", "page")
			active.Text(item.Name)
			nav.AddChild(active)
		} else {
			// Regular breadcrumb with link
			li := hb.Li()
			li.Class("breadcrumb-item")

			link := hb.A()
			link.Href(item.URL)
			link.Text(item.Name)

			li.AddChild(link)
			nav.AddChild(li)
		}
	}

	return nav.ToHTML(), nil
}

// renderPlainBreadcrumbs renders plain breadcrumbs without Bootstrap classes
func renderPlainBreadcrumbs(breadcrumbs []BreadcrumbItem, style, cssClass, cssID, separator string) (string, error) {
	nav := hb.Nav()

	// Base classes
	nav.Class("breadcrumbs")
	nav.Class(fmt.Sprintf("breadcrumbs-style-%s", style))

	// Add custom CSS classes
	if cssClass != "" {
		nav.Class(cssClass)
	}

	// Add CSS ID
	if cssID != "" {
		nav.ID(cssID)
	}

	// Add breadcrumb items
	for _, item := range breadcrumbs {
		if len(breadcrumbs) > 1 && item != breadcrumbs[0] {
			// Add separator
			sep := hb.Span()
			sep.Class("breadcrumb-separator")
			sep.Text(separator)
			nav.AddChild(sep)
		}

		if item.Active {
			// Active breadcrumb (current page)
			active := hb.Span()
			active.Class("breadcrumb-item active")
			active.Text(item.Name)
			nav.AddChild(active)
		} else {
			// Regular breadcrumb with link
			link := hb.A()
			link.Class("breadcrumb-link")
			link.Href(item.URL)
			link.Text(item.Name)
			nav.AddChild(link)
		}
	}

	return nav.ToHTML(), nil
}
