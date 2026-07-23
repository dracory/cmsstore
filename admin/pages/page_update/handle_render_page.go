package page_update

import (
	"net/http"

	"github.com/dracory/bs"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
)

func handleRenderPage(ui uiInterface, store cmsstore.StoreInterface, page cmsstore.PageInterface, view string, w http.ResponseWriter, r *http.Request) string {
	siteList := []cmsstore.SiteInterface{}
	if store != nil {
		sites, err := store.SiteList(r.Context(), cmsstore.SiteQuery().
			SetOrderBy(cmsstore.COLUMN_NAME).
			SetSortOrder(cmsstore.SORT_ORDER_ASC).
			SetOffset(0).
			SetLimit(100))
		if err == nil {
			siteList = sites
		}
	}

	breadcrumbs := shared.AdminBreadcrumbs(r, []shared.Breadcrumb{
		{
			Name: "Page Manager",
			URL:  shared.URLR(r, shared.PathPagesPageManager, nil),
		},
		{
			Name: "Edit Page",
			URL:  shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": page.ID()}),
		},
	}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: siteList,
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		ID("btn-page-save")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.URLR(r, shared.PathPagesPageManager, nil))

	badgeStatus := hb.Div().
		Class("badge fs-6 ms-3").
		ClassIf(page.Status() == cmsstore.PAGE_STATUS_ACTIVE, "bg-success").
		ClassIf(page.Status() == cmsstore.PAGE_STATUS_INACTIVE, "bg-secondary").
		ClassIf(page.Status() == cmsstore.PAGE_STATUS_DRAFT, "bg-warning").
		Text(page.Status())

	pageTitle := hb.Heading1().
		HTML("Edit Page: ").
		Text(page.Name()).
		Child(hb.Sup().Child(badgeStatus)).
		Child(buttonSave).
		Child(buttonCancel)

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == viewContent, "active").
				Href(shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{
					"page_id": page.ID(),
					"view":    viewContent,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == viewSEO, "active").
				Href(shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{
					"page_id": page.ID(),
					"view":    viewSEO,
				})).
				HTML("SEO"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == viewMiddlewares, "active").
				Href(shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{
					"page_id": page.ID(),
					"view":    viewMiddlewares,
				})).
				HTML("Middlewares"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == viewSettings, "active").
				Href(shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{
					"page_id": page.ID(),
					"view":    viewSettings,
				})).
				HTML("Settings")))

	var body hb.TagInterface

	switch view {
	case viewContent:
		body = renderContentTab(r, page)
	case viewSEO:
		body = renderSEOTab(r, page)
	case viewSettings:
		body = renderSettingsTab(r, page)
	case viewMiddlewares:
		body = renderMiddlewaresTab(r, page)
	default:
		body = renderContentTab(r, page)
	}

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTMLIf(view == viewContent, "Page Contents").
					HTMLIf(view == viewSEO, "Page SEO").
					HTMLIf(view == viewMiddlewares, "Page Middlewares").
					HTMLIf(view == viewSettings, "Page Settings").
					Style("margin-bottom:0;display:inline-block;"))).
		Child(
			hb.Div().
				Class("card-body").
				Child(body))

	content := hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(pageTitle).
		Child(tabs).
		Child(card)

	options := struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{}

	return ui.Layout(w, r, "Edit Page | CMS", content.ToHTML(), options)
}
