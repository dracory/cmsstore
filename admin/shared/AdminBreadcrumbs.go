package shared

import (
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
)

func AdminBreadcrumbs(r *http.Request, pageBreadcrumbs []Breadcrumb, options struct {
	SiteList []cmsstore.SiteInterface
}) hb.TagInterface {
	adminHomeURL := AdminHomeURL(r)
	siteID := utils.Req(r, "filter_site_id", "")
	site, siteFound := lo.Find(options.SiteList, func(site cmsstore.SiteInterface) bool {
		return site.ID() == siteID
	})
	path := utils.Req(r, "path", "")

	adminHomeBreadcrumb := lo.
		If(adminHomeURL != "", Breadcrumb{
			Name: "Home",
			URL:  adminHomeURL,
		}).
		Else(Breadcrumb{})

	breadcrumbItems := []Breadcrumb{
		adminHomeBreadcrumb,
		{
			Name: "CMS",
			URL:  URLR(r, PathHome, nil),
		},
	}

	breadcrumbItems = append(breadcrumbItems, pageBreadcrumbs...)

	breadcrumbs := Breadcrumbs(breadcrumbItems)

	dropdown := hb.Div().Class("dropdown float-end").
		Child(hb.Button().
			Class("btn btn-secondary dropdown-toggle").
			Type("button").
			Attr("data-bs-toggle", "dropdown").
			Attr("aria-expanded", "false").
			Text("Site: ").
			Text(lo.IfF(siteFound, func() string { return site.Name() }).Else("all sites")))

	dropdownMenu := hb.UL().Class("dropdown-menu")

	for _, site := range options.SiteList {
		link := hb.Hyperlink().
			Text(site.Name()).
			Href(URLR(r, path, map[string]string{
				"filter_site_id": site.ID(),
			}))

		dropdownMenu.Child(hb.LI().
			Class("dropdown-item").
			Child(link))
	}

	dropdown.Child(dropdownMenu)

	return hb.Div().
		Child(breadcrumbs).
		Child(dropdown)
}
