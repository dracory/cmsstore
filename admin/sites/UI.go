package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
)

type UiConfig struct {
	AdminHeader hb.TagInterface
	Endpoint    string
	Layout      func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger *slog.Logger
	Store  cmsstore.StoreInterface
}

func UI(config UiConfig) UiInterface {
	return ui{
		adminHeader: config.AdminHeader,
		endpoint:    config.Endpoint,
		layout:      config.Layout,
		logger:      config.Logger,
		store:       config.Store,
	}
}

type UiInterface interface {
	AdminHeader() hb.TagInterface
	Endpoint() string
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	SiteCreate(w http.ResponseWriter, r *http.Request)
	SiteManager(w http.ResponseWriter, r *http.Request)
	SiteDelete(w http.ResponseWriter, r *http.Request)
	SiteUpdate(w http.ResponseWriter, r *http.Request)
	Store() cmsstore.StoreInterface
}

type ui struct {
	adminHeader hb.TagInterface
	endpoint    string
	layout      func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
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
