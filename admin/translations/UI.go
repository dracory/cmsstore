package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/responses"
)

type UiConfig struct {
	Endpoint     string
	AdminHeader  hb.TagInterface
	AdminHomeURL string
	Layout       func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
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
		endpoint:     config.Endpoint,
		adminHeader:  config.AdminHeader,
		adminHomeURL: config.AdminHomeURL,
		layout:       config.Layout,
		logger:       config.Logger,
		store:        config.Store,
	}
}

type UiInterface interface {
	Endpoint() string
	AdminHeader() hb.TagInterface
	AdminHomeURL() string
	Layout(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	Logger() *slog.Logger
	TranslationCreate(w http.ResponseWriter, r *http.Request)
	TranslationManager(w http.ResponseWriter, r *http.Request)
	TranslationDelete(w http.ResponseWriter, r *http.Request)
	TranslationUpdate(w http.ResponseWriter, r *http.Request)
	Store() cmsstore.StoreInterface
}

type ui struct {
	endpoint     string
	adminHeader  hb.TagInterface
	adminHomeURL string
	layout       func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger                 *slog.Logger
	store                  cmsstore.StoreInterface
	url                    func(endpoint string, path string, params map[string]string) string
	pathTranslationCreate  string
	pathTranslationDelete  string
	pathTranslationManager string
	pathTranslationUpdate  string
}

func (ui ui) AdminHeader() hb.TagInterface {
	return ui.adminHeader
}

func (ui ui) AdminHomeURL() string {
	return ui.adminHomeURL
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
