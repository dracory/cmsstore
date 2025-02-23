package cmsstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/base/database"
	"github.com/gouniverse/versionstore"
	"github.com/samber/lo"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	// Context is the context used if the AutoMigrateEnabled option is true
	// If not set, a background context is used
	Context context.Context

	// DB is the database connection
	DB *sql.DB

	// DbDriverName is the database driver name
	// If not set, an attempt will be made to detect it
	DbDriverName string

	// AutomigrateEnabled enables automigrate
	AutomigrateEnabled bool

	// DebugEnabled enables debug
	DebugEnabled bool

	// BlockTableName is the name of the block database table to be created/used
	BlockTableName string

	// MenusEnabled enables menus
	MenusEnabled bool

	// MenuTableName is the name of the menu database table to be created/used
	MenuTableName string

	// MenuItemTableName is the name of the menu item database table to be created/used
	MenuItemTableName string

	// PageTableName is the name of the page database table to be created/used
	PageTableName string

	// SiteTableName is the name of the site database table to be created/used
	SiteTableName string

	TemplateTableName string

	// TranslationsEnabled enables translations
	TranslationsEnabled bool

	// TranslationTableName is the name of the translation database table to be created/used
	TranslationTableName string

	// TranslationLanguageDefault is the default language, i.e en
	TranslationLanguageDefault string

	// TranslationLanguages is the list of supported languages
	TranslationLanguages map[string]string

	// VersioningEnabled enables versioning
	VersioningEnabled bool

	// VersioningTableName is the name of the versioning database table to be created/used
	VersioningTableName string

	// Shortcodes is a list of shortcodes to be registered
	Shortcodes []ShortcodeInterface

	// Middlewares is a list of middlewares to be registered
	Middlewares []MiddlewareInterface
}

// NewStore creates a new block store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.BlockTableName == "" {
		return nil, errors.New("cms store: BlockTableName is required")
	}

	if opts.PageTableName == "" {
		return nil, errors.New("cms store: PageTableName is required")
	}

	if opts.SiteTableName == "" {
		return nil, errors.New("cms store: SiteTableName is required")
	}

	if opts.TemplateTableName == "" {
		return nil, errors.New("cms store: TemplateTableName is required")
	}

	if opts.MenusEnabled && opts.MenuTableName == "" {
		return nil, errors.New("cms store: MenuTableName is required")
	}

	if opts.MenusEnabled && opts.MenuItemTableName == "" {
		return nil, errors.New("cms store: MenuItemTableName is required")
	}

	if opts.TranslationsEnabled && opts.TranslationTableName == "" {
		return nil, errors.New("cms store: TranslationTableName is required")
	}

	if opts.VersioningEnabled && opts.VersioningTableName == "" {
		return nil, errors.New("cms store: VersioningTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("cms store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = database.DatabaseType(opts.DB)
	}

	if len(opts.Shortcodes) == 0 {
		opts.Shortcodes = []ShortcodeInterface{}
	}

	if len(opts.Middlewares) == 0 {
		opts.Middlewares = []MiddlewareInterface{}
	}

	versionStore, err := initializeVersioningStore(opts)

	if err != nil {
		return nil, err
	}

	store := &store{
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,

		blockTableName:    opts.BlockTableName,
		pageTableName:     opts.PageTableName,
		siteTableName:     opts.SiteTableName,
		templateTableName: opts.TemplateTableName,

		menusEnabled:      opts.MenusEnabled,
		menuTableName:     opts.MenuTableName,
		menuItemTableName: opts.MenuItemTableName,

		translationsEnabled:        opts.TranslationsEnabled,
		translationTableName:       opts.TranslationTableName,
		translationLanguageDefault: opts.TranslationLanguageDefault,
		translationLanguages:       opts.TranslationLanguages,

		versioningEnabled: opts.VersioningEnabled,
		// versioningTableName: opts.VersioningTableName,
		versioningStore: versionStore,

		shortcodes:  opts.Shortcodes,
		middlewares: opts.Middlewares,
	}

	if store.automigrateEnabled {
		context := lo.If(opts.Context != nil, opts.Context).Else(context.Background())

		if err := store.AutoMigrate(context); err != nil {
			return nil, err
		}

		if opts.VersioningEnabled {
			if err := versionStore.AutoMigrate(); err != nil {
				return nil, err
			}
		}
	}

	return store, nil
}

func initializeVersioningStore(opts NewStoreOptions) (versionstore.StoreInterface, error) {
	if !opts.VersioningEnabled {
		return nil, nil
	}

	if opts.VersioningTableName == "" {
		return nil, errors.New("cms store: VersioningTableName is required")
	}

	versionStore, err := versionstore.NewStore(versionstore.NewStoreOptions{
		TableName:          opts.VersioningTableName,
		DB:                 opts.DB,
		AutomigrateEnabled: opts.AutomigrateEnabled,
		DebugEnabled:       opts.DebugEnabled,
		// Logger:             nil,
	})

	if err != nil {
		return nil, err
	}

	if versionStore == nil {
		return nil, errors.New("cms store: version store is nil")
	}

	return versionStore, nil
}
