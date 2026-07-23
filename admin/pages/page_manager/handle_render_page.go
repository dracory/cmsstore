package page_manager

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cdn"
	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
)

func handleRenderPage(ui shared.UiInterface, store cmsstore.StoreInterface, w http.ResponseWriter, r *http.Request) string {
	siteList := []cmsstore.SiteInterface{}
	if store != nil {
		sites, err := store.SiteList(r.Context(), cmsstore.SiteQuery().
			SetOrderBy(cmsstore.COLUMN_NAME).
			SetSortOrder(cmsstore.SORT_ORDER_ASC).
			SetOffset(0).
			SetLimit(100))
		if err != nil {
			slog.Error("Failed to load sites for breadcrumbs", "error", err)
		} else {
			siteList = sites
		}
	}

	breadcrumbs := shared.AdminBreadcrumbs(r, []shared.Breadcrumb{
		{
			Name: "Page Manager",
			URL:  shared.URLR(r, shared.PathPagesPageManager, nil),
		},
	}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: siteList,
	})

	actionButtons := hb.Div().
		Class("d-flex gap-2 float-end")

	buttonPageNew := hb.Button().
		Class("btn btn-primary d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-plus-circle me-2")).
		HTML("New Page").
		ID("btn-page-new")

	actionButtons = actionButtons.Child(buttonPageNew)

	heading := hb.Heading1().HTML("Page Manager").Child(actionButtons)

	htmlContent, err := pageFiles.ReadFile("pages.html")
	if err != nil {
		slog.Error("Failed to read pages HTML template", "error", err)
		return hb.Div().HTML("Error loading pages component").ToHTML()
	}

	jsContent, err := pageFiles.ReadFile("pages.js")
	if err != nil {
		slog.Error("Failed to read pages JavaScript file", "error", err)
		return hb.Div().HTML("Error loading posts component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const urlPagesLoad = '` + shared.URLR(r, shared.PathPagesPageManager, map[string]string{"action": actionLoadPages}) + `';
		const urlPageDelete = '` + shared.URLR(r, shared.PathPagesPageManager, map[string]string{"action": actionDeletePage}) + `';
		const urlPageCreate = '` + shared.URLR(r, shared.PathPagesPageManager, map[string]string{"action": actionCreatePage}) + `';
		const urlPageUpdate = '` + shared.URLR(r, shared.PathPagesPageUpdate, map[string]string{"page_id": "PAGE_ID_PLACEHOLDER"}) + `';
		const urlPageDeleteConfirm = '` + shared.URLR(r, shared.PathPagesPageDelete, map[string]string{"page_id": "PAGE_ID_PLACEHOLDER"}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueContainer)

	options := struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
		},
	}

	return ui.Layout(w, r, "Page Manager | CMS", content.ToHTML(), options)
}
