package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
)

type UiConfig struct {
	BlockEditorDefinitions []blockeditor.BlockDefinition
	AdminHeader            hb.TagInterface
	AdminHomeURL           string
	Endpoint               string
	Layout                 func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
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
		blockEditorDefinitions: config.BlockEditorDefinitions,
		adminHeader:            config.AdminHeader,
		adminHomeURL:           config.AdminHomeURL,
		endpoint:               config.Endpoint,
		layout:                 config.Layout,
		logger:                 config.Logger,
		store:                  config.Store,
	}
}

type UiInterface interface {
	BlockEditorDefinitions() []blockeditor.BlockDefinition
	AdminHeader() hb.TagInterface
	AdminHomeURL() string
	Endpoint() string
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	Store() cmsstore.StoreInterface
	PageCreate(w http.ResponseWriter, r *http.Request)
	PageManager(w http.ResponseWriter, r *http.Request)
	PageDelete(w http.ResponseWriter, r *http.Request)
	PageUpdate(w http.ResponseWriter, r *http.Request)
}

type ui struct {
	blockEditorDefinitions []blockeditor.BlockDefinition
	adminHeader            hb.TagInterface
	adminHomeURL           string
	endpoint               string
	layout                 func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
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

func (ui ui) AdminHomeURL() string {
	return ui.adminHomeURL
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

func (ui ui) Store() cmsstore.StoreInterface {
	return ui.store
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
	controller := NewPageDeleteController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) PageUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewPageUpdateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}
