package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
)

type UiConfig struct {
	Endpoint    string
	AdminHeader hb.TagInterface
	Layout      func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger              *slog.Logger
	Store               cmsstore.StoreInterface
	URL                 func(endpoint string, path string, params map[string]string) string
	PathTemplateCreate  string
	PathTemplateDelete  string
	PathTemplateManager string
	PathTemplateUpdate  string
}

func UI(config UiConfig) UiInterface {
	return ui{
		endpoint:            config.Endpoint,
		adminHeader:         config.AdminHeader,
		layout:              config.Layout,
		logger:              config.Logger,
		store:               config.Store,
		url:                 config.URL,
		pathTemplateCreate:  config.PathTemplateCreate,
		pathTemplateDelete:  config.PathTemplateDelete,
		pathTemplateManager: config.PathTemplateManager,
		pathTemplateUpdate:  config.PathTemplateUpdate,
	}
}

type UiInterface interface {
	Endpoint() string
	AdminHeader() hb.TagInterface
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	PathTemplateCreate() string
	PathTemplateDelete() string
	PathTemplateManager() string
	PathTemplateUpdate() string
	TemplateCreate(w http.ResponseWriter, r *http.Request)
	TemplateManager(w http.ResponseWriter, r *http.Request)
	TemplateDelete(w http.ResponseWriter, r *http.Request)
	TemplateUpdate(w http.ResponseWriter, r *http.Request)
	Store() cmsstore.StoreInterface
	URL(endpoint string, path string, params map[string]string) string
}

type ui struct {
	endpoint    string
	adminHeader hb.TagInterface
	layout      func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger              *slog.Logger
	store               cmsstore.StoreInterface
	url                 func(endpoint string, path string, params map[string]string) string
	pathTemplateCreate  string
	pathTemplateDelete  string
	pathTemplateManager string
	pathTemplateUpdate  string
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

func (ui ui) PathTemplateCreate() string {
	return ui.pathTemplateCreate
}

func (ui ui) PathTemplateDelete() string {
	return ui.pathTemplateDelete
}

func (ui ui) PathTemplateManager() string {
	return ui.pathTemplateManager
}

func (ui ui) PathTemplateUpdate() string {
	return ui.pathTemplateUpdate
}

func (ui ui) Store() cmsstore.StoreInterface {
	return ui.store
}

func (ui ui) URL(endpoint string, path string, params map[string]string) string {
	return ui.url(endpoint, path, params)
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
