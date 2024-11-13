package admin

import (
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func (a *admin) pageHome(w http.ResponseWriter, r *http.Request) {
	siteList, err := a.store.SiteList(cmsstore.SiteQuery().
		SetOrderBy(cmsstore.COLUMN_NAME).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		siteList = []cmsstore.SiteInterface{}
	}

	adminHeader := shared.AdminHeader(a.store, a.logger, r)
	adminBreadcrumbs := shared.AdminBreadcrumbs(r, []shared.Breadcrumb{}, struct{ SiteList []cmsstore.SiteInterface }{
		SiteList: siteList,
	})

	pagesCount, errPagesCount := a.store.PageCount(cmsstore.PageQuery())

	if errPagesCount != nil {
		pagesCount = 0
	}

	sitesCount, errSitesCount := a.store.SiteCount(cmsstore.SiteQuery())

	if errSitesCount != nil {
		sitesCount = 0
	}

	templatesCount, errTemplatesCount := a.store.TemplateCount(cmsstore.TemplateQuery())

	if errTemplatesCount != nil {
		templatesCount = 0
	}

	blocksCount, errBlocksCount := a.store.BlockCount(cmsstore.BlockQuery())

	if errBlocksCount != nil {
		blocksCount = 0
	}

	tiles := []struct {
		Count      string
		Title      string
		Background string
		Icon       string
		URL        string
	}{

		{
			Count:      cast.ToString(sitesCount),
			Title:      "Total Sites",
			Background: "bg-success",
			Icon:       "bi-globe",
			URL:        shared.URLR(r, shared.PathSitesSiteManager, nil),
		},
		{
			Count:      cast.ToString(pagesCount),
			Title:      "Total Pages",
			Background: "bg-info",
			Icon:       "bi-journals",
			URL:        shared.URLR(r, shared.PathPagesPageManager, nil),
		},
		{
			Count:      cast.ToString(templatesCount),
			Title:      "Total Templates",
			Background: "bg-warning",
			Icon:       "bi-file-earmark-text-fill",
			URL:        shared.URLR(r, shared.PathTemplatesTemplateManager, nil),
		},
		{
			Count:      cast.ToString(blocksCount),
			Title:      "Total Blocks",
			Background: "bg-primary",
			Icon:       "bi-grid-3x3-gap-fill",
			URL:        shared.URLR(r, shared.PathBlocksBlockManager, nil),
		},
	}

	cards := lo.Map(tiles, func(tile struct {
		Count      string
		Title      string
		Background string
		Icon       string
		URL        string
	}, index int) hb.TagInterface {
		card := hb.Div().
			Class("card").
			Class("bg-transparent border round-10 shadow-lg h-100").
			// OnMouseOver(`this.style.setProperty('background-color', 'beige', 'important');this.style.setProperty('scale', 1.1);this.style.setProperty('border', '4px solid moccasin', 'important');`).
			// OnMouseOut(`this.style.setProperty('background-color', 'transparent', 'important');this.style.setProperty('scale', 1);this.style.setProperty('border', '4px solid transparent', 'important');`).
			Child(hb.Div().
				Class("card-body").
				Class(tile.Background).
				Style("--bs-bg-opacity:0.3;").
				Child(hb.Div().Class("row").
					Child(hb.Div().Class("col-8").
						Child(hb.Div().
							Style("margin-top:-4px;margin-right:8px;font-size:32px;").
							Text(tile.Count)).
						Child(hb.NewDiv().
							Style("margin-top:-4px;margin-right:8px;font-size:16px;").
							Text(tile.Title)),
					).
					Child(hb.Div().Class("col-4").
						Child(hb.I().
							Class("bi float-end").
							Class(tile.Icon).
							Style(`color:silver;opacity:0.6;`).
							Style("margin-top:-4px;margin-right:8px;font-size:48px;")),
					),
				)).
			Child(hb.Div().
				Class("card-footer text-center").
				Class(tile.Background).
				Style("--bs-bg-opacity:0.5;").
				Child(hb.A().
					Class("text-white").
					Href(tile.URL).
					Text("More info").
					Child(hb.I().Class("bi bi-arrow-right-circle-fill ms-3").Style("margin-top:-4px;margin-right:8px;font-size:16px;")),
				))
		return hb.Div().Class("col-xs-12 col-sm-6 col-md-3").Child(card)
	})

	pageTitle := hb.NewHeading1().
		HTML("Content Management Dashboard")

	container := hb.NewDiv().
		ID("page-manager").
		Class("container").
		Child(adminBreadcrumbs).
		Child(hb.HR()).
		Child(adminHeader).
		Child(hb.HR()).
		Child(pageTitle).
		Child(hb.Div().Class("row g-3").Children(cards))

	a.render(w, r, "Home", container.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{})
}
