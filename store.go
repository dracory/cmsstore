package cmsstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/database"
	"github.com/gouniverse/versionstore"
)

// == TYPE ====================================================================

// store represents the core structure for managing CMS data.
type store struct {
	blockTableName     string
	pageTableName      string
	siteTableName      string
	templateTableName  string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool

	// Menus
	menusEnabled      bool
	menuTableName     string
	menuItemTableName string

	// Translations
	translationsEnabled        bool
	translationTableName       string
	translationLanguages       map[string]string
	translationLanguageDefault string

	versioningEnabled bool
	//versioningTableName string
	versioningStore versionstore.StoreInterface

	// Shortcodes
	shortcodes  []ShortcodeInterface
	middlewares []MiddlewareInterface
}

// == INTERFACE ===============================================================

var _ StoreInterface = (*store)(nil) // verify it extends the interface

// PUBLIC METHODS ============================================================

// AutoMigrate performs automatic database migrations.
func (store *store) AutoMigrate(ctx context.Context, opts ...Option) error {
	if store.db == nil {
		return errors.New("cms store: database is nil")
	}

	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	transaction, hasTransaction := options.params["tx"].(*sql.Tx)
	isDryRun, hasDryRun := options.params["dryRun"].(bool)

	blockSql := store.blockTableCreateSql()
	menuSql := store.menuTableCreateSql()
	menuItemSql := store.menuItemTableCreateSql()
	pageSql := store.pageTableCreateSql()
	tableSql := store.siteTableCreateSql()
	templateSql := store.templateTableCreateSql()
	translationSql := store.translationTableCreateSql()

	if blockSql == "" {
		return errors.New("block table create sql is empty")
	}

	if pageSql == "" {
		return errors.New("page table create sql is empty")
	}

	if tableSql == "" {
		return errors.New("site table create sql is empty")
	}

	if templateSql == "" {
		return errors.New("template table create sql is empty")
	}

	if store.menusEnabled && store.menuTableName == "" {
		return errors.New("menu table name is empty")
	}

	if store.menusEnabled && store.menuItemTableName == "" {
		return errors.New("menu item table name is empty")
	}

	if store.translationsEnabled && translationSql == "" {
		return errors.New("translation table create sql is empty")
	}

	// if store.versioningEnabled && store.versioningTableName == "" {
	// 	return errors.New("versioning table name is empty")
	// }

	sqlList := []string{
		blockSql,
		pageSql,
		tableSql,
		templateSql,
	}

	if store.menusEnabled {
		sqlList = append(sqlList, menuSql)
		sqlList = append(sqlList, menuItemSql)
	}

	if store.translationsEnabled {
		sqlList = append(sqlList, translationSql)
	}

	for _, sql := range sqlList {
		if hasDryRun && isDryRun {
			continue
		}

		if hasTransaction {
			_, err := transaction.ExecContext(ctx, sql)

			if err != nil {
				return err
			}

			continue
		} else {
			_, err := store.db.ExecContext(ctx, sql)

			if err != nil {
				return err
			}
		}
	}

	if store.versioningEnabled {
		err := store.versioningStore.AutoMigrate()

		if err != nil {
			return err
		}
	}

	return nil
}

// EnableDebug enables or disables debug mode.
func (st *store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// MenusEnabled checks if menus are enabled.
func (store *store) MenusEnabled() bool {
	return store.menusEnabled
}

// TranslationsEnabled checks if translations are enabled.
func (store *store) TranslationsEnabled() bool {
	return store.translationsEnabled
}

// VersioningEnabled checks if versioning is enabled.
func (store *store) VersioningEnabled() bool {
	return store.versioningEnabled
}

// Shortcodes returns the list of shortcodes.
func (store *store) Shortcodes() []ShortcodeInterface {
	return store.shortcodes
}

// AddShortcode adds a shortcode to the store.
func (store *store) AddShortcode(shortcode ShortcodeInterface) {
	store.shortcodes = append(store.shortcodes, shortcode)
}

// AddShortcodes adds multiple shortcodes to the store.
func (store *store) AddShortcodes(shortcodes []ShortcodeInterface) {
	store.shortcodes = append(store.shortcodes, shortcodes...)
}

// SetShortcodes sets the list of shortcodes.
func (store *store) SetShortcodes(shortcodes []ShortcodeInterface) {
	store.shortcodes = shortcodes
}

// Middlewares returns the list of middlewares.
func (store *store) Middlewares() []MiddlewareInterface {
	return store.middlewares
}

// AddMiddleware adds a middleware to the store.
func (store *store) AddMiddleware(middleware MiddlewareInterface) {
	store.middlewares = append(store.middlewares, middleware)
}

// AddMiddlewares adds multiple middlewares to the store.
func (store *store) AddMiddlewares(middlewares []MiddlewareInterface) {
	store.middlewares = append(store.middlewares, middlewares...)
}

// SetMiddlewares sets the list of middlewares.
func (store *store) SetMiddlewares(middlewares []MiddlewareInterface) {
	store.middlewares = middlewares
}

// toQuerableContext converts a context to a queryable context.
func (store *store) toQuerableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}

	return database.Context(ctx, store.db)
}
