package shared

import "github.com/dracory/hb"

// PageHeader renders the page heading with optional right-aligned action
// buttons, matching the admin section pattern (breadcrumb + H2 + actions).
func PageHeader(title string, actions ...hb.TagInterface) hb.TagInterface {
	header := hb.Div().
		Class("d-flex justify-content-between align-items-center mb-3").
		Child(hb.H2().Text(title).Class("mb-0"))
	if len(actions) > 0 {
		header.Child(hb.Div().Class("d-flex gap-2 align-items-center").Children(actions))
	}
	return header
}
