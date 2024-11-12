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
	adminTranslations "github.com/gouniverse/cmsstore/admin/translations"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

type AdminOptions struct {
	BlockEditorDefinitions []blockeditor.BlockDefinition
	FuncLayout             func(title string, body string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger       *slog.Logger
	Store        cmsstore.StoreInterface
	AdminHomeURL string
}

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
		adminHomeURL:           options.AdminHomeURL,
	}, nil
}

type Admin interface {
	Handle(w http.ResponseWriter, r *http.Request)
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
	logger       *slog.Logger
	store        cmsstore.StoreInterface
	adminHomeURL string
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
	maps.Copy(routes, a.translationRoutes())

	if val, ok := routes[route]; ok {
		return val
	}

	return routes[shared.PathHome]
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

func (a *admin) blockRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	blockRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathBlocksBlockCreate: func(w http.ResponseWriter, r *http.Request) {
			adminBlocks.UI(a.uiConfig(r)).BlockCreate(w, r)
		},
		shared.PathBlocksBlockDelete: func(w http.ResponseWriter, r *http.Request) {
			adminBlocks.UI(a.uiConfig(r)).BlockDelete(w, r)
		},
		shared.PathBlocksBlockManager: func(w http.ResponseWriter, r *http.Request) {
			adminBlocks.UI(a.uiConfig(r)).BlockManager(w, r)
		},
		shared.PathBlocksBlockUpdate: func(w http.ResponseWriter, r *http.Request) {
			adminBlocks.UI(a.uiConfig(r)).BlockUpdate(w, r)
		},
	}
	return blockRoutes
}

func (a *admin) pageRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	pageRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathPagesPageCreate: func(w http.ResponseWriter, r *http.Request) {
			adminPages.UI(a.uiConfig(r)).PageCreate(w, r)
		},
		shared.PathPagesPageDelete: func(w http.ResponseWriter, r *http.Request) {
			adminPages.UI(a.uiConfig(r)).PageDelete(w, r)
		},
		shared.PathPagesPageManager: func(w http.ResponseWriter, r *http.Request) {
			adminPages.UI(a.uiConfig(r)).PageManager(w, r)
		},
		shared.PathPagesPageUpdate: func(w http.ResponseWriter, r *http.Request) {
			adminPages.UI(a.uiConfig(r)).PageUpdate(w, r)
		},
	}
	return pageRoutes
}

func (a *admin) siteRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	siteRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathSitesSiteCreate: func(w http.ResponseWriter, r *http.Request) {
			adminSites.UI(a.uiConfig(r)).SiteCreate(w, r)
		},
		shared.PathSitesSiteDelete: func(w http.ResponseWriter, r *http.Request) {
			adminSites.UI(a.uiConfig(r)).SiteDelete(w, r)
		},
		shared.PathSitesSiteUpdate: func(w http.ResponseWriter, r *http.Request) {
			adminSites.UI(a.uiConfig(r)).SiteUpdate(w, r)
		},
		shared.PathSitesSiteManager: func(w http.ResponseWriter, r *http.Request) {
			adminSites.UI(a.uiConfig(r)).SiteManager(w, r)
		},
	}

	return siteRoutes
}

func (a *admin) templateRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	templateRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathTemplatesTemplateCreate: func(w http.ResponseWriter, r *http.Request) {
			adminTemplates.UI(a.uiConfig(r)).TemplateCreate(w, r)
		},
		shared.PathTemplatesTemplateDelete: func(w http.ResponseWriter, r *http.Request) {
			adminTemplates.UI(a.uiConfig(r)).TemplateDelete(w, r)
		},
		shared.PathTemplatesTemplateManager: func(w http.ResponseWriter, r *http.Request) {
			adminTemplates.UI(a.uiConfig(r)).TemplateManager(w, r)
		},
		shared.PathTemplatesTemplateUpdate: func(w http.ResponseWriter, r *http.Request) {
			adminTemplates.UI(a.uiConfig(r)).TemplateUpdate(w, r)
		},
	}
	return templateRoutes
}

func (a *admin) translationRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	translationsRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathTranslationsTranslationCreate: func(w http.ResponseWriter, r *http.Request) {
			adminTranslations.UI(a.uiConfig(r)).TranslationCreate(w, r)
		},
		shared.PathTranslationsTranslationDelete: func(w http.ResponseWriter, r *http.Request) {
			adminTranslations.UI(a.uiConfig(r)).TranslationDelete(w, r)
		},
		shared.PathTranslationsTranslationManager: func(w http.ResponseWriter, r *http.Request) {
			adminTranslations.UI(a.uiConfig(r)).TranslationManager(w, r)
		},
		shared.PathTranslationsTranslationUpdate: func(w http.ResponseWriter, r *http.Request) {
			adminTranslations.UI(a.uiConfig(r)).TranslationUpdate(w, r)
		},
	}
	return translationsRoutes
}

func (a *admin) adminBreadcrumbs(endpoint string, pageBreadcrumbs []shared.Breadcrumb) hb.TagInterface {
	return shared.AdminBreadcrumbs(a.adminHomeURL, endpoint, pageBreadcrumbs)
}

func (a *admin) uiConfig(r *http.Request) shared.UiConfig {
	return shared.UiConfig{
		BlockEditorDefinitions: a.blockEditorDefinitions,
		AdminBreadcrumbs:       a.adminBreadcrumbs,
		Endpoint:               shared.Endpoint(r),
		Layout:                 a.render,
		Logger:                 a.logger,
		Store:                  a.store,
	}
}
