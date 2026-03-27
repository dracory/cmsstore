package frontend

import (
	"log/slog"

	"github.com/dracory/cmsstore"
	"github.com/dracory/ui"
)

type Config struct {
	BlockEditorRenderer func(blocks []ui.BlockInterface) string
	Logger              *slog.Logger
	Shortcodes          []cmsstore.ShortcodeInterface
	Store               cmsstore.StoreInterface
	CacheEnabled        bool
	CacheExpireSeconds  int
}

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
