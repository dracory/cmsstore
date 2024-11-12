package shared

import (
	"github.com/gouniverse/hb"
	"github.com/samber/lo"
)

func AdminBreadcrumbs(adminHomeURL string, endpoint string, pageBreadcrumbs []Breadcrumb) hb.TagInterface {
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
