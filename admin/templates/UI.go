package admin

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/hb"
	"github.com/gouniverse/responses"
)

func UI(config shared.UiConfig) UiInterface {
	return ui{

		layout: config.Layout,
		logger: config.Logger,
		store:  config.Store,
	}
}

type UiInterface interface {
	shared.UiInterface
	TemplateCreate(w http.ResponseWriter, r *http.Request)
	TemplateManager(w http.ResponseWriter, r *http.Request)
	TemplateDelete(w http.ResponseWriter, r *http.Request)
	TemplateUpdate(w http.ResponseWriter, r *http.Request)
}

type ui struct {
	endpoint    string
	adminHeader hb.TagInterface

	layout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger *slog.Logger
	store  cmsstore.StoreInterface
}

func (ui ui) AdminHeader() hb.TagInterface {
	return ui.adminHeader
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

func (ui ui) TemplateCreate(w http.ResponseWriter, r *http.Request) {
	controller := NewTemplateCreateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) TemplateManager(w http.ResponseWriter, r *http.Request) {
	controller := NewTemplateManagerController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) TemplateDelete(w http.ResponseWriter, r *http.Request) {
	controller := NewTemplateDeleteController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) TemplateUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewTemplateUpdateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}
