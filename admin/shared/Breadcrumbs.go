package shared

import "github.com/dracory/hb"

func Breadcrumbs(breadcrumbs []Breadcrumb) hb.TagInterface {

	ol := hb.OL().Attr("class", "breadcrumb").Style("border:none;")

	for _, breadcrumb := range breadcrumbs {

		link := hb.Hyperlink().
			HTML(breadcrumb.Name).
			Href(breadcrumb.URL)

		li := hb.LI().
			Class("breadcrumb-item").
			Child(link)

		ol.AddChild(li)
	}

	nav := hb.Nav().
		Attr("aria-label", "breadcrumb").
		Child(ol)

	return nav
}
