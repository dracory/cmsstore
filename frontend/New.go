package frontend

import (
	"log/slog"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/ui"
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

	frontend := frontend{
		blockEditorRenderer: config.BlockEditorRenderer,
		logger:              config.Logger,
		shortcodes:          config.Shortcodes,
		store:               config.Store,
		cacheEnabled:        config.CacheEnabled,
		cacheExpireSeconds:  config.CacheExpireSeconds,
	}

	if config.CacheEnabled {
		cache := initCache()

		if cache != nil {
			frontend.cache = cache

			go frontend.warmUpCache()
		}

	}

	return &frontend
}
