package admin

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
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
	SiteCreate(w http.ResponseWriter, r *http.Request)
	SiteManager(w http.ResponseWriter, r *http.Request)
	SiteDelete(w http.ResponseWriter, r *http.Request)
	SiteUpdate(w http.ResponseWriter, r *http.Request)
}

type ui struct {
	endpoint string
	layout   func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger *slog.Logger
	store  cmsstore.StoreInterface
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

func (ui ui) SiteCreate(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteCreateController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (ui ui) SiteManager(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteManagerController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (ui ui) SiteDelete(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteDeleteController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (ui ui) SiteUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteUpdateController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}
