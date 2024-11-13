package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
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
	MenuCreate(w http.ResponseWriter, r *http.Request)
	MenuManager(w http.ResponseWriter, r *http.Request)
	MenuDelete(w http.ResponseWriter, r *http.Request)
	MenuUpdate(w http.ResponseWriter, r *http.Request)
}

type ui struct {
	layout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger *slog.Logger
	store  cmsstore.StoreInterface
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

func (ui ui) MenuCreate(w http.ResponseWriter, r *http.Request) {
	controller := NewMenuCreateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) MenuManager(w http.ResponseWriter, r *http.Request) {
	controller := NewMenuManagerController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) MenuDelete(w http.ResponseWriter, r *http.Request) {
	controller := NewMenuDeleteController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) MenuUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewMenuUpdateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}
