package admin

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"net/http"

	"github.com/gouniverse/blockeditor"
	adminBlocks "github.com/gouniverse/cmsstore/admin/blocks"
	adminPages "github.com/gouniverse/cmsstore/admin/pages"
	"github.com/gouniverse/cmsstore/admin/shared"
	adminSites "github.com/gouniverse/cmsstore/admin/sites"
	adminTemplates "github.com/gouniverse/cmsstore/admin/templates"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func New(options AdminOptions) (*admin, error) {
	if options.Store == nil {
		return nil, errors.New(shared.ERROR_STORE_IS_NIL)
	}

	if options.Logger == nil {
		return nil, errors.New(shared.ERROR_LOGGER_IS_NIL)
	}

	return &admin{
		blockEditorDefinitions: options.BlockEditorDefinitions,
		logger:                 options.Logger,
		store:                  options.Store,
		funcLayout:             options.FuncLayout,
	}, nil
}

type Admin interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type AdminOptions struct {
	BlockEditorDefinitions []blockeditor.BlockDefinition
	FuncLayout             func(title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger *slog.Logger
	Store  cmsstore.StoreInterface
}

var _ Admin = (*admin)(nil)

type admin struct {
	blockEditorDefinitions []blockeditor.BlockDefinition
	funcLayout             func(title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger *slog.Logger
	store  cmsstore.StoreInterface
}

func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	path := utils.Req(r, "path", "home")

	if path == "" {
		path = shared.PathHome
	}

	ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)

	routeFunc := a.getRoute(path)
	routeFunc(w, r.WithContext(ctx))
}

func (a *admin) getRoute(route string) func(w http.ResponseWriter, r *http.Request) {

	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathHome: a.pageHome,
	}

	maps.Copy(routes, a.blockRoutes())
	maps.Copy(routes, a.pageRoutes())
	maps.Copy(routes, a.siteRoutes())
	maps.Copy(routes, a.templateRoutes())

	if val, ok := routes[route]; ok {
		return val
	}

	return routes[shared.PathHome]
}

func (a *admin) pageHome(w http.ResponseWriter, r *http.Request) {
	header := a.adminHeader(shared.Endpoint(r))
	breadcrumbs := a.adminBreadcrumbs([]bs.Breadcrumb{
		{
			URL:  shared.Endpoint(r),
			Name: "Home",
		},
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
			URL:        shared.URL(shared.Endpoint(r), shared.PathSitesSiteManager, nil),
		},
		{
			Count:      cast.ToString(pagesCount),
			Title:      "Total Pages",
			Background: "bg-info",
			Icon:       "bi-journals",
			URL:        shared.URL(shared.Endpoint(r), shared.PathPagesPageManager, nil),
		},
		{
			Count:      cast.ToString(templatesCount),
			Title:      "Total Templates",
			Background: "bg-warning",
			Icon:       "bi-file-earmark-text-fill",
			URL:        shared.URL(shared.Endpoint(r), shared.PathTemplatesTemplateManager, nil),
		},
		{
			Count:      cast.ToString(blocksCount),
			Title:      "Total Blocks",
			Background: "bg-primary",
			Icon:       "bi-grid-3x3-gap-fill",
			URL:        shared.URL(shared.Endpoint(r), shared.PathBlocksBlockManager, nil),
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

	heading := hb.NewHeading1().
		HTML("Content Management Dashboard")

	container := hb.NewDiv().
		ID("page-manager").
		Class("container").
		Child(header).
		Child(heading).
		Child(hb.NewHTML(breadcrumbs)).
		Child(hb.Div().Class("row g-3").Children(cards))

	a.render(w, r, "Home", container.ToHTML(), struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}{})
}

func (a *admin) render(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	webpage := webpageComplete(webpageTitle, webpageHtml, options).ToHTML()

	if a.funcLayout != nil {
		isNotEmpty := a.funcLayout("", "", struct {
			Styles     []string
			StyleURLs  []string
			Scripts    []string
			ScriptURLs []string
		}{}) != ""
		if isNotEmpty {
			webpage = a.funcLayout(webpageTitle, webpageHtml, options)
		}
	}

	responses.HTMLResponse(w, r, webpage)
	return ""
}

func (a *admin) blockUI(r *http.Request) adminBlocks.UiInterface {
	options := adminBlocks.UiConfig{
		// BlockEditorDefinitions: a.blockEditorDefinitions,
		AdminHeader: a.adminHeader(shared.Endpoint(r)),
		Endpoint:    shared.Endpoint(r),
		Layout:      a.render,
		Logger:      a.logger,
		Store:       a.store,
	}
	return adminBlocks.UI(options)
}

func (a *admin) blockRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	blockRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathBlocksBlockCreate: func(w http.ResponseWriter, r *http.Request) {
			a.blockUI(r).BlockCreate(w, r)
		},
		shared.PathBlocksBlockDelete: func(w http.ResponseWriter, r *http.Request) {
			a.blockUI(r).BlockDelete(w, r)
		},
		shared.PathBlocksBlockManager: func(w http.ResponseWriter, r *http.Request) {
			a.blockUI(r).BlockManager(w, r)
		},
		shared.PathBlocksBlockUpdate: func(w http.ResponseWriter, r *http.Request) {
			a.blockUI(r).BlockUpdate(w, r)
		},
	}
	return blockRoutes
}

func (a *admin) pageUI(r *http.Request) adminPages.UiInterface {
	options := adminPages.UiConfig{
		BlockEditorDefinitions: a.blockEditorDefinitions,
		AdminHeader:            a.adminHeader(shared.Endpoint(r)),
		Endpoint:               shared.Endpoint(r),
		Layout:                 a.render,
		Logger:                 a.logger,
		Store:                  a.store,
	}
	return adminPages.UI(options)
}

func (a *admin) pageRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	pageRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathPagesPageCreate: func(w http.ResponseWriter, r *http.Request) {
			a.pageUI(r).PageCreate(w, r)
		},
		shared.PathPagesPageDelete: func(w http.ResponseWriter, r *http.Request) {
			a.pageUI(r).PageDelete(w, r)
		},
		shared.PathPagesPageManager: func(w http.ResponseWriter, r *http.Request) {
			a.pageUI(r).PageManager(w, r)
		},
		shared.PathPagesPageUpdate: func(w http.ResponseWriter, r *http.Request) {
			a.pageUI(r).PageUpdate(w, r)
		},
	}
	return pageRoutes
}

func (a *admin) siteUI(r *http.Request) adminSites.UiInterface {
	options := adminSites.UiConfig{
		AdminHeader: a.adminHeader(shared.Endpoint(r)),
		Endpoint:    shared.Endpoint(r),
		Layout:      a.render,
		Logger:      a.logger,
		Store:       a.store,
	}
	return adminSites.UI(options)
}

func (a *admin) siteRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	siteRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathSitesSiteCreate: func(w http.ResponseWriter, r *http.Request) {
			a.siteUI(r).SiteCreate(w, r)
		},
		shared.PathSitesSiteDelete: func(w http.ResponseWriter, r *http.Request) {
			a.siteUI(r).SiteDelete(w, r)
		},
		shared.PathSitesSiteUpdate: func(w http.ResponseWriter, r *http.Request) {
			a.siteUI(r).SiteUpdate(w, r)
		},
		shared.PathSitesSiteManager: func(w http.ResponseWriter, r *http.Request) {
			a.siteUI(r).SiteManager(w, r)
		},
	}

	return siteRoutes
}

func (a *admin) templateUI(r *http.Request) adminTemplates.UiInterface {
	options := adminTemplates.UiConfig{
		Endpoint:    shared.Endpoint(r),
		AdminHeader: a.adminHeader(shared.Endpoint(r)),
		Layout:      a.render,
		Logger:      a.logger,
		Store:       a.store,
	}
	return adminTemplates.UI(options)
}

func (a *admin) templateRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	templateRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathTemplatesTemplateCreate: func(w http.ResponseWriter, r *http.Request) {
			a.templateUI(r).TemplateCreate(w, r)
		},
		shared.PathTemplatesTemplateDelete: func(w http.ResponseWriter, r *http.Request) {
			a.templateUI(r).TemplateDelete(w, r)
		},
		shared.PathTemplatesTemplateManager: func(w http.ResponseWriter, r *http.Request) {
			a.templateUI(r).TemplateManager(w, r)
		},
		shared.PathTemplatesTemplateUpdate: func(w http.ResponseWriter, r *http.Request) {
			a.templateUI(r).TemplateUpdate(w, r)
		},
	}
	return templateRoutes
}

func (a *admin) adminBreadcrumbs(breadcrumbs []bs.Breadcrumb) string {
	return bs.Breadcrumbs(breadcrumbs).
		Style("margin-bottom:10px;").
		ToHTML()
}

func (a *admin) adminHeader(endpoint string) hb.TagInterface {
	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(endpoint + "").
		Class("nav-link")
	linkBlocks := hb.NewHyperlink().
		HTML("Blocks ").
		Href(endpoint + "?path=" + shared.PathBlocksBlockManager).
		Class("nav-link")
	// linkMenus := hb.NewHyperlink().
	// 	HTML("Menus ").
	// 	Href(endpoint + "?path=" + PathMenusMenuManager).
	// 	Class("nav-link")
	linkPages := hb.NewHyperlink().
		HTML("Pages ").
		Href(endpoint + "?path=" + shared.PathPagesPageManager).
		Class("nav-link")
	linkTemplates := hb.NewHyperlink().
		HTML("Templates ").
		Href(endpoint + "?path=" + shared.PathTemplatesTemplateManager).
		Class("nav-link")
	linkSites := hb.NewHyperlink().
		HTML("Sites ").
		Href(endpoint + "?path=" + shared.PathBlocksBlockManager).
		Class("nav-link")
	// linkWidgets := hb.NewHyperlink().
	// 	HTML("Widgets ").
	// 	Href(endpoint + "?path=" + PathWidgetsWidgetManager).
	// 	Class("nav-link")
	// linkSettings := hb.NewHyperlink().
	// 	HTML("Settings").
	// 	Href(endpoint + "?path=" + PathSettingsSettingManager).
	// 	Class("nav-link")
	// linkTranslations := hb.NewHyperlink().
	// 	HTML("Translations").
	// 	Href(endpoint + "?path=" + PathTranslationsTranslationManager).
	// 	Class("nav-link")

	templatesCount, err := a.store.TemplateCount(cmsstore.TemplateQuery())

	if err != nil {
		a.logger.Error(err.Error())
		templatesCount = -1
	}

	blocksCount, err := a.store.BlockCount(cmsstore.BlockQuery())

	if err != nil {
		a.logger.Error(err.Error())
		blocksCount = -1
	}

	pagesCount, err := a.store.PageCount(cmsstore.PageQuery())

	if err != nil {
		a.logger.Error(err.Error())
		pagesCount = -1
	}

	sitesCount, err := a.store.SiteCount(cmsstore.SiteQuery())

	if err != nil {
		a.logger.Error(err.Error())
		sitesCount = -1
	}

	ulNav := hb.NewUL().Class("nav  nav-pills justify-content-center")
	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkHome))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkSites.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(sitesCount)))))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkTemplates.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(templatesCount)))))

	ulNav.Child(hb.
		LI().
		Class("nav-item").
		Child(linkPages.
			Child(hb.NewSpan().
				Class("badge bg-secondary").
				HTML(cast.ToString(pagesCount)))))

	// if cms.menusEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkMenus.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(menusCount, 10)))))
	// }

	ulNav.Child(hb.
		LI().
		Class("nav-item").
		Child(linkBlocks.
			Child(hb.NewSpan().
				Class("badge bg-secondary").
				HTML(cast.ToString(blocksCount)))))

	// if cms.widgetsEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkWidgets.AddChild(hb.NewSpan().Class("badge bg-secondary").HTML(strconv.FormatInt(widgetsCount, 10)))))
	// }

	// if cms.translationsEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkTranslations.Child(hb.NewSpan().Class("badge bg-secondary").HTML(utils.ToString(translationsCount)))))
	// }

	// if cms.settingsEnabled {
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").AddChild(linkSettings))
	// }
	// add Translations

	// for _, entity := range cms.customEntityList {
	// 	linkEntity := hb.NewHyperlink().HTML(entity.TypeLabel).Href(endpoint + "?path=entities/entity-manager&type=" + entity.Type).Class("nav-link")
	// 	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkEntity))
	// }

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav))
}
