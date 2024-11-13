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
	TranslationCreate(w http.ResponseWriter, r *http.Request)
	TranslationManager(w http.ResponseWriter, r *http.Request)
	TranslationDelete(w http.ResponseWriter, r *http.Request)
	TranslationUpdate(w http.ResponseWriter, r *http.Request)
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

func (ui ui) TranslationCreate(w http.ResponseWriter, r *http.Request) {
	controller := NewTranslationCreateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) TranslationManager(w http.ResponseWriter, r *http.Request) {
	controller := NewTranslationManagerController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) TranslationDelete(w http.ResponseWriter, r *http.Request) {
	controller := NewTranslationDeleteController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}

func (ui ui) TranslationUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewTranslationUpdateController(ui)
	html := controller.Handler(w, r)
	responses.HTMLResponse(w, r, html)
}
