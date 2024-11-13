package shared

import (
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

func AdminBreadcrumbs(r *http.Request, pageBreadcrumbs []Breadcrumb) hb.TagInterface {
	adminHomeURL := AdminHomeURL(r)
	endpoint := Endpoint(r)
	adminHomeBreadcrumb := lo.If(adminHomeURL != "", Breadcrumb{
		Name: "Home",
		URL:  adminHomeURL,
	}).Else(Breadcrumb{})

	breadcrumbItems := []Breadcrumb{
		adminHomeBreadcrumb,
		{
			Name: "CMS",
			URL:  URL(endpoint, PathHome, nil),
		},
	}

	breadcrumbItems = append(breadcrumbItems, pageBreadcrumbs...)

	breadcrumbs := Breadcrumbs(breadcrumbItems)

	return breadcrumbs
}
