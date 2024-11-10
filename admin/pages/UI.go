package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/responses"
)

type UiConfig struct {
	BlockEditorDefinitions []blockeditor.BlockDefinition
	Endpoint               string
	Layout                 func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger          *slog.Logger
	Store           cmsstore.StoreInterface
	URL             func(endpoint string, path string, params map[string]string) string
	PathPageCreate  string
	PathPageDelete  string
	PathPageManager string
	PathPageUpdate  string
}

func UI(config UiConfig) UiInterface {
	return ui{
		blockEditorDefinitions: config.BlockEditorDefinitions,
		endpoint:               config.Endpoint,
		layout:                 config.Layout,
		logger:                 config.Logger,
		store:                  config.Store,
		url:                    config.URL,
		pathPageCreate:         config.PathPageCreate,
		pathPageDelete:         config.PathPageDelete,
		pathPageManager:        config.PathPageManager,
		pathPageUpdate:         config.PathPageUpdate,
	}
}

type UiInterface interface {
	BlockEditorDefinitions() []blockeditor.BlockDefinition
	Endpoint() string
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	PathPageCreate() string
	PathPageDelete() string
	PathPageManager() string
	PathPageUpdate() string
	PageCreate(w http.ResponseWriter, r *http.Request)
	PageManager(w http.ResponseWriter, r *http.Request)
	PageDelete(w http.ResponseWriter, r *http.Request)
	PageUpdate(w http.ResponseWriter, r *http.Request)
	Store() cmsstore.StoreInterface
	URL(endpoint string, path string, params map[string]string) string
}

type ui struct {
	blockEditorDefinitions []blockeditor.BlockDefinition
	endpoint               string
	layout                 func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger          *slog.Logger
	store           cmsstore.StoreInterface
	url             func(endpoint string, path string, params map[string]string) string
	pathPageCreate  string
	pathPageDelete  string
	pathPageManager string
	pathPageUpdate  string
}

func (ui ui) BlockEditorDefinitions() []blockeditor.BlockDefinition {
	return ui.blockEditorDefinitions
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

func (ui ui) PathPageCreate() string {
	return ui.pathPageCreate
}

func (ui ui) PathPageDelete() string {
	return ui.pathPageDelete
}

func (ui ui) PathPageManager() string {
	return ui.pathPageManager
}

func (ui ui) PathPageUpdate() string {
	return ui.pathPageUpdate
}

func (ui ui) Store() cmsstore.StoreInterface {
	return ui.store
}

func (ui ui) URL(endpoint string, path string, params map[string]string) string {
	return ui.url(endpoint, path, params)
}

func (ui ui) PageCreate(w http.ResponseWriter, r *http.Request) {
	controller := NewPageCreateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) PageManager(w http.ResponseWriter, r *http.Request) {
	controller := NewPageManagerController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) PageDelete(w http.ResponseWriter, r *http.Request) {
	// controller := NewPageDeleteController(ui)
	// html := controller.Handler(w, r)
	// responses.HTMLResponse(w, r, html)
	responses.HTMLResponse(w, r, "Not implemented")
}

func (ui ui) PageUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewPageUpdateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}
