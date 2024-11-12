package shared

import "github.com/gouniverse/hb"

func Breadcrumbs(breadcrumbs []Breadcrumb) hb.TagInterface {
	nav := hb.Nav().Attr("aria-label", "breadcrumb")
	ol := hb.OL().Attr("class", "breadcrumb")

	for _, breadcrumb := range breadcrumbs {
		li := hb.LI().Attr("class", "breadcrumb-item")
		link := hb.Hyperlink().HTML(breadcrumb.Name).Attr("href", breadcrumb.URL)

		li.AddChild(link)

		ol.AddChild(li)
	}

	nav.AddChild(ol)

	return nav
}
