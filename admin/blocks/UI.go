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
	Logger *slog.Logger
	Store  cmsstore.StoreInterface
	// URL    func(endpoint string, path string, params map[string]string) string
	// PathBlockCreate  string
	// PathBlockDelete  string
	// PathBlockManager string
	// PathBlockUpdate  string
}

func UI(config UiConfig) UiInterface {
	return ui{
		endpoint: config.Endpoint,
		layout:   config.Layout,
		logger:   config.Logger,
		store:    config.Store,
		// url:      config.URL,
		// pathBlockCreate:  config.PathBlockCreate,
		// pathBlockDelete:  config.PathBlockDelete,
		// pathBlockManager: config.PathBlockManager,
		// pathBlockUpdate:  config.PathBlockUpdate,
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
	// PathBlockCreate() string
	// PathBlockDelete() string
	// PathBlockManager() string
	// PathBlockUpdate() string
	BlockCreate(w http.ResponseWriter, r *http.Request)
	BlockManager(w http.ResponseWriter, r *http.Request)
	BlockDelete(w http.ResponseWriter, r *http.Request)
	BlockUpdate(w http.ResponseWriter, r *http.Request)
	Store() cmsstore.StoreInterface
	// URL(endpoint string, path string, params map[string]string) string
}

type ui struct {
	endpoint string
	layout   func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger           *slog.Logger
	store            cmsstore.StoreInterface
	url              func(endpoint string, path string, params map[string]string) string
	pathBlockCreate  string
	pathBlockDelete  string
	pathBlockManager string
	pathBlockUpdate  string
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

// func (ui ui) PathBlockCreate() string {
// 	return ui.pathBlockCreate
// }

// func (ui ui) PathBlockDelete() string {
// 	return ui.pathBlockDelete
// }

// func (ui ui) PathBlockManager() string {
// 	return ui.pathBlockManager
// }

// func (ui ui) PathBlockUpdate() string {
// 	return ui.pathBlockUpdate
// }

func (ui ui) Store() cmsstore.StoreInterface {
	return ui.store
}

// func (ui ui) URL(endpoint string, path string, params map[string]string) string {
// 	return ui.url(endpoint, path, params)
// }

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
