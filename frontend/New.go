// Package frontend provides a CMS frontend rendering system.
// It handles page rendering, block rendering, caching, and template processing
// for a content management system.
package frontend

import (
	"log/slog"
	"net/http"

	"github.com/dracory/cmsstore"
	"github.com/dracory/ui"
)

// Config holds the configuration options for creating a new Frontend instance.
// All fields are optional and have sensible defaults if not provided.
type Config struct {
	// BlockEditorRenderer converts block editor JSON content to HTML.
	// Required if using block editor pages.
	BlockEditorRenderer func(blocks []ui.BlockInterface) string
	// Logger is used for structured logging. If nil, logging is disabled.
	Logger *slog.Logger

	// Shortcodes are custom shortcode handlers for content processing.
	Shortcodes []cmsstore.ShortcodeInterface

	// Store provides access to CMS data (pages, blocks, templates, etc.).
	// Required for frontend operation.
	Store cmsstore.StoreInterface

	// CacheEnabled enables in-memory caching of rendered content.
	CacheEnabled bool

	// CacheExpireSeconds sets the TTL for cached items.
	// Defaults to 600 seconds (10 minutes) if not set or <= 0.
	CacheExpireSeconds int

	// PageNotFoundHandler is called when a page is not found.
	// If it returns handled=true, the frontend will use the result and skip the default 404 response.
	PageNotFoundHandler func(w http.ResponseWriter, r *http.Request, alias string) (handled bool, result string)
}

// New creates a new Frontend instance with the provided configuration.
//
// It initializes the frontend with the following features:
//   - Block renderer registry for custom block types
//   - Optional caching system with TTL-based expiration
//   - Cache warming on startup if caching is enabled
//
// Cache Configuration:
//   - If CacheEnabled is true and CacheExpireSeconds is not set or <= 0,
//     it defaults to 10 minutes (600 seconds)
//   - A background goroutine warms up the cache after initialization
//
// Example usage:
//
//	frontend := New(Config{
//	    Store:               store,
//	    Logger:              logger,
//	    BlockEditorRenderer: myRenderer,
//	    CacheEnabled:        true,
//	    CacheExpireSeconds:  300,
//	})
//
// Parameters:
//   - config: the configuration options for the frontend
//
// Returns:
//   - FrontendInterface: a configured frontend instance ready to handle requests
func New(config Config) FrontendInterface {
	if config.CacheEnabled && config.CacheExpireSeconds <= 0 {
		config.CacheExpireSeconds = 10 * 60 // 10 minutes
	}

	f := frontend{
		blockEditorRenderer: config.BlockEditorRenderer,
		logger:              config.Logger,
		shortcodes:          config.Shortcodes,
		store:               config.Store,
		cacheEnabled:        config.CacheEnabled,
		cacheExpireSeconds:  config.CacheExpireSeconds,
		pageNotFoundHandler: config.PageNotFoundHandler,
	}
	f.blockRenderers = initBlockRenderers(&f, config.Store)

	if config.CacheEnabled {
		cache := initCache()

		if cache != nil {
			f.cache = cache

			go f.warmUpCache()
		}

	}

	return &f
}
