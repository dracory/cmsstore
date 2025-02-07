package admin

import (
	"context"
	"log/slog"
	"maps"
	"net/http"

	"github.com/gouniverse/blockeditor"
	adminBlocks "github.com/gouniverse/cmsstore/admin/blocks"
	adminMenus "github.com/gouniverse/cmsstore/admin/menus"
	adminPages "github.com/gouniverse/cmsstore/admin/pages"
	"github.com/gouniverse/cmsstore/admin/shared"
	adminSites "github.com/gouniverse/cmsstore/admin/sites"
	adminTemplates "github.com/gouniverse/cmsstore/admin/templates"
	adminTranslations "github.com/gouniverse/cmsstore/admin/translations"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/responses"
	"github.com/gouniverse/utils"
)

// == TYPE ====================================================================

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
	flags        map[string]bool
}

// == INTERFACE IMPLEMENTATION CHECK ==========================================

var _ shared.Admin = (*admin)(nil)

// == INTERFACE IMPLEMENTATION ================================================

func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	path := utils.Req(r, "path", "home")

	if path == "" {
		path = shared.PathHome
	}

	ctx := context.WithValue(r.Context(), shared.KeyEndpoint, r.URL.Path)
	ctx = context.WithValue(ctx, shared.KeyAdminHomeURL, a.adminHomeURL)

	routeFunc := a.getRoute(path)
	routeFunc(w, r.WithContext(ctx))
}

// == PRIVATE METHODS =========================================================

func (a *admin) getRoute(route string) func(w http.ResponseWriter, r *http.Request) {

	routes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathHome: a.pageHome,
	}

	maps.Copy(routes, a.blockRoutes())

	if a.store.MenusEnabled() {
		maps.Copy(routes, a.menuRoutes())
	}

	maps.Copy(routes, a.pageRoutes())
	maps.Copy(routes, a.siteRoutes())
	maps.Copy(routes, a.templateRoutes())

	if a.store.TranslationsEnabled() {
		maps.Copy(routes, a.translationRoutes())
	}

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
		shared.PathBlocksBlockCreate:  adminBlocks.UI(a.uiConfig()).BlockCreate,
		shared.PathBlocksBlockDelete:  adminBlocks.UI(a.uiConfig()).BlockDelete,
		shared.PathBlocksBlockManager: adminBlocks.UI(a.uiConfig()).BlockManager,
		shared.PathBlocksBlockUpdate:  adminBlocks.UI(a.uiConfig()).BlockUpdate,
	}
	return blockRoutes
}

func (a *admin) menuRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	menuRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathMenusMenuCreate:  adminMenus.UI(a.uiConfig()).MenuCreate,
		shared.PathMenusMenuDelete:  adminMenus.UI(a.uiConfig()).MenuDelete,
		shared.PathMenusMenuManager: adminMenus.UI(a.uiConfig()).MenuManager,
		shared.PathMenusMenuUpdate:  adminMenus.UI(a.uiConfig()).MenuUpdate,
	}
	return menuRoutes
}

func (a *admin) pageRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	pageRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathPagesPageCreate:     adminPages.UI(a.uiConfig()).PageCreate,
		shared.PathPagesPageDelete:     adminPages.UI(a.uiConfig()).PageDelete,
		shared.PathPagesPageManager:    adminPages.UI(a.uiConfig()).PageManager,
		shared.PathPagesPageUpdate:     adminPages.UI(a.uiConfig()).PageUpdate,
		shared.PathPagesPageVersioning: adminPages.UI(a.uiConfig()).PageVersioning,
	}
	return pageRoutes
}

func (a *admin) siteRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	siteRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathSitesSiteCreate:  adminSites.UI(a.uiConfig()).SiteCreate,
		shared.PathSitesSiteDelete:  adminSites.UI(a.uiConfig()).SiteDelete,
		shared.PathSitesSiteUpdate:  adminSites.UI(a.uiConfig()).SiteUpdate,
		shared.PathSitesSiteManager: adminSites.UI(a.uiConfig()).SiteManager,
	}

	return siteRoutes
}

func (a *admin) templateRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	templateRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathTemplatesTemplateCreate:  adminTemplates.UI(a.uiConfig()).TemplateCreate,
		shared.PathTemplatesTemplateDelete:  adminTemplates.UI(a.uiConfig()).TemplateDelete,
		shared.PathTemplatesTemplateManager: adminTemplates.UI(a.uiConfig()).TemplateManager,
		shared.PathTemplatesTemplateUpdate:  adminTemplates.UI(a.uiConfig()).TemplateUpdate,
	}
	return templateRoutes
}

func (a *admin) translationRoutes() map[string]func(w http.ResponseWriter, r *http.Request) {
	translationsRoutes := map[string]func(w http.ResponseWriter, r *http.Request){
		shared.PathTranslationsTranslationCreate:  adminTranslations.UI(a.uiConfig()).TranslationCreate,
		shared.PathTranslationsTranslationDelete:  adminTranslations.UI(a.uiConfig()).TranslationDelete,
		shared.PathTranslationsTranslationManager: adminTranslations.UI(a.uiConfig()).TranslationManager,
		shared.PathTranslationsTranslationUpdate:  adminTranslations.UI(a.uiConfig()).TranslationUpdate,
	}
	return translationsRoutes
}

// func (a *admin) adminBreadcrumbs(r *http.Request, pageBreadcrumbs []shared.Breadcrumb) hb.TagInterface {
// 	return shared.AdminBreadcrumbs(r, pageBreadcrumbs)
// }

func (a *admin) uiConfig() shared.UiConfig {
	return shared.UiConfig{
		BlockEditorDefinitions: a.blockEditorDefinitions,
		Layout:                 a.render,
		Logger:                 a.logger,
		Store:                  a.store,
	}
}
