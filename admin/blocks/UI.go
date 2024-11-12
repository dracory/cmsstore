package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
)

func UI(config shared.UiConfig) UiInterface {
	return ui{
		// adminHeader:      config.AdminHeader,
		adminBreadcrumbs: config.AdminBreadcrumbs,
		endpoint:         config.Endpoint,
		layout:           config.Layout,
		logger:           config.Logger,
		store:            config.Store,
	}
}

type UiInterface interface {
	shared.UiInterface
	BlockCreate(w http.ResponseWriter, r *http.Request)
	BlockManager(w http.ResponseWriter, r *http.Request)
	BlockDelete(w http.ResponseWriter, r *http.Request)
	BlockUpdate(w http.ResponseWriter, r *http.Request)
}

type ui struct {
	// adminHeader      hb.TagInterface
	adminBreadcrumbs func(endpoint string, breadcrumbs []shared.Breadcrumb) hb.TagInterface
	endpoint         string
	layout           func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger *slog.Logger
	store  cmsstore.StoreInterface
}

// func (ui ui) AdminHeader() hb.TagInterface {
// 	return ui.adminHeader
// }

func (ui ui) AdminBreadcrumbs(endpoint string, breadcrumbs []shared.Breadcrumb) hb.TagInterface {
	return ui.adminBreadcrumbs(endpoint, breadcrumbs)
}

func (ui ui) Endpoint() string {
	return ui.endpoint
}

func (ui ui) Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
	Styles     []string
	StyleURLs  []string
	Scripts    []string
	ScriptURLs []string
}) string {
	return ui.layout(w, r, webpageTitle, webpageHtml, options)
}

func (ui ui) Logger() *slog.Logger {
	return ui.logger
}

func (ui ui) Store() cmsstore.StoreInterface {
	return ui.store
}

func (ui ui) BlockCreate(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockCreateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) BlockManager(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockManagerController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) BlockDelete(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockDeleteController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) BlockUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockUpdateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}
