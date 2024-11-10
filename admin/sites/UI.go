package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/responses"
)

type UiConfig struct {
	Endpoint string
	Layout   func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger          *slog.Logger
	Store           cmsstore.StoreInterface
	URL             func(endpoint string, path string, params map[string]string) string
	PathSiteCreate  string
	PathSiteDelete  string
	PathSiteManager string
	PathSiteUpdate  string
}

func UI(config UiConfig) UiInterface {
	return ui{
		endpoint:        config.Endpoint,
		layout:          config.Layout,
		logger:          config.Logger,
		store:           config.Store,
		url:             config.URL,
		pathSiteCreate:  config.PathSiteCreate,
		pathSiteDelete:  config.PathSiteDelete,
		pathSiteManager: config.PathSiteManager,
		pathSiteUpdate:  config.PathSiteUpdate,
	}
}

type UiInterface interface {
	Endpoint() string
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	PathSiteCreate() string
	PathSiteDelete() string
	PathSiteManager() string
	PathSiteUpdate() string
	SiteCreate(w http.ResponseWriter, r *http.Request)
	SiteManager(w http.ResponseWriter, r *http.Request)
	SiteDelete(w http.ResponseWriter, r *http.Request)
	SiteUpdate(w http.ResponseWriter, r *http.Request)
	Store() cmsstore.StoreInterface
	URL(endpoint string, path string, params map[string]string) string
}

type ui struct {
	endpoint string
	layout   func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger          *slog.Logger
	store           cmsstore.StoreInterface
	url             func(endpoint string, path string, params map[string]string) string
	pathSiteCreate  string
	pathSiteDelete  string
	pathSiteManager string
	pathSiteUpdate  string
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

func (ui ui) PathSiteCreate() string {
	return ui.pathSiteCreate
}

func (ui ui) PathSiteDelete() string {
	return ui.pathSiteDelete
}

func (ui ui) PathSiteManager() string {
	return ui.pathSiteManager
}

func (ui ui) PathSiteUpdate() string {
	return ui.pathSiteUpdate
}

func (ui ui) Store() cmsstore.StoreInterface {
	return ui.store
}

func (ui ui) URL(endpoint string, path string, params map[string]string) string {
	return ui.url(endpoint, path, params)
}

func (ui ui) SiteCreate(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteCreateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) SiteManager(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteManagerController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) SiteDelete(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteDeleteController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) SiteUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewSiteUpdateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}
