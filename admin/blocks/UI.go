package admin

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
)

func UI(config shared.UiConfig) UiInterface {
	registry := initBlockAdminProviders(config.Store, config.Logger)

	return ui{
		//
		layout:             config.Layout,
		logger:             config.Logger,
		store:              config.Store,
		blockAdminRegistry: registry,
	}
}

// initBlockAdminProviders initializes and registers all built-in block admin providers
func initBlockAdminProviders(store cmsstore.StoreInterface, logger *slog.Logger) *BlockAdminFieldProviderRegistry {
	registry := NewBlockAdminFieldProviderRegistry()

	// Register built-in HTML provider
	registry.Register(cmsstore.BLOCK_TYPE_HTML, NewHTMLAdminProvider())

	// Register built-in Menu provider
	registry.Register(cmsstore.BLOCK_TYPE_MENU, NewMenuAdminProvider(store, logger))

	return registry
}

type UiInterface interface {
	shared.UiInterface
	BlockCreate(w http.ResponseWriter, r *http.Request)
	BlockManager(w http.ResponseWriter, r *http.Request)
	BlockDelete(w http.ResponseWriter, r *http.Request)
	BlockUpdate(w http.ResponseWriter, r *http.Request)
	BlockVersioning(w http.ResponseWriter, r *http.Request)
	BlockAdminRegistry() *BlockAdminFieldProviderRegistry
}

type ui struct {
	// endpoint string
	layout func(w http.ResponseWriter, r *http.Request, webpageTitle, webpageHtml string, options struct {
		Styles     []string
		StyleURLs  []string
		Scripts    []string
		ScriptURLs []string
	}) string
	logger             *slog.Logger
	store              cmsstore.StoreInterface
	blockAdminRegistry *BlockAdminFieldProviderRegistry
}

// func (ui ui) Endpoint() string {
// 	return ui.endpoint
// }

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (ui ui) BlockManager(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockManagerController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (ui ui) BlockDelete(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockDeleteController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (ui ui) BlockUpdate(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockUpdateController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (ui ui) BlockVersioning(w http.ResponseWriter, r *http.Request) {
	controller := NewBlockVersioningController(ui)
	html := controller.Handler(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

// BlockAdminRegistry returns the block admin field provider registry, allowing
// external packages to register custom block type admin UI providers.
//
// Example usage:
//
//	adminUI := admin.UI(config)
//	adminUI.BlockAdminRegistry().Register("gallery", &GalleryAdminProvider{store: store})
//
// Custom admin providers must implement the BlockAdminFieldProvider interface:
//
//	type BlockAdminFieldProvider interface {
//	    GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface
//	    GetTypeLabel() string
//	    SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error
//	}
//
// See admin/blocks/admin_field_provider.go for detailed interface documentation
// and admin/blocks/README.md for complete examples.
func (ui ui) BlockAdminRegistry() *BlockAdminFieldProviderRegistry {
	return ui.blockAdminRegistry
}
