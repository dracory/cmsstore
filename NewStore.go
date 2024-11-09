package cmsstore

import (
	"database/sql"
	"errors"

	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new block store
type NewStoreOptions struct {
	BlockTableName     string
	PageTableName      string
	SiteTableName      string
	TemplateTableName  string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
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

	if opts.DB == nil {
		return nil, errors.New("cms store: DB is required")
	}

	if opts.DbDriverName == "" {
		opts.DbDriverName = sb.DatabaseDriverName(opts.DB)
	}

	store := &store{
		blockTableName:     opts.BlockTableName,
		pageTableName:      opts.PageTableName,
		siteTableName:      opts.SiteTableName,
		templateTableName:  opts.TemplateTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.automigrateEnabled {
		err := store.AutoMigrate()

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}
