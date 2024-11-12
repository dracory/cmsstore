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
	BlockCreate(w http.ResponseWriter, r *http.Request)
	BlockManager(w http.ResponseWriter, r *http.Request)
	BlockDelete(w http.ResponseWriter, r *http.Request)
	BlockUpdate(w http.ResponseWriter, r *http.Request)
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
	logger           *slog.Logger
	store            cmsstore.StoreInterface
	url              func(endpoint string, path string, params map[string]string) string
	pathBlockCreate  string
	pathBlockDelete  string
	pathBlockManager string
	pathBlockUpdate  string
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
