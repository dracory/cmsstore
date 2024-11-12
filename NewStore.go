package cmsstore

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	DB                         *sql.DB
	DbDriverName               string
	AutomigrateEnabled         bool
	DebugEnabled               bool
	BlockTableName             string
	PageTableName              string
	SiteTableName              string
	TemplateTableName          string
	TranslationTableName       string
	TranslationsEnabled        bool
	TranslationLanguageDefault string
	TranslationLanguages       map[string]string
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

	if opts.TranslationsEnabled && opts.TranslationTableName == "" {
		return nil, errors.New("cms store: TranslationTableName is required")
	}

	if opts.DB == nil {
		return nil, errors.New("cms store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	store := &store{
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,

		blockTableName:       opts.BlockTableName,
		pageTableName:        opts.PageTableName,
		siteTableName:        opts.SiteTableName,
		templateTableName:    opts.TemplateTableName,
		translationTableName: opts.TranslationTableName,

		translationsEnabled:        opts.TranslationsEnabled,
		translationLanguageDefault: opts.TranslationLanguageDefault,
		translationLanguages:       opts.TranslationLanguages,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
