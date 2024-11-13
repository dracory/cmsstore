package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/cmsstore/admin/shared"
	"github.com/gouniverse/responses"
)

func UI(config shared.UiConfig) UiInterface {
	return ui{
		blockEditorDefinitions: config.BlockEditorDefinitions,
		layout:                 config.Layout,
		logger:                 config.Logger,
		store:                  config.Store,
	}
}

type UiInterface interface {
	shared.UiInterface
	BlockEditorDefinitions() []blockeditor.BlockDefinition
	PageCreate(w http.ResponseWriter, r *http.Request)
	PageManager(w http.ResponseWriter, r *http.Request)
	PageDelete(w http.ResponseWriter, r *http.Request)
	PageUpdate(w http.ResponseWriter, r *http.Request)
}

type ui struct {
	blockEditorDefinitions []blockeditor.BlockDefinition
	layout                 func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger *slog.Logger
	store  cmsstore.StoreInterface
}

func (ui ui) BlockEditorDefinitions() []blockeditor.BlockDefinition {
	return ui.blockEditorDefinitions
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
