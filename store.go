package cmsstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/database"
	"github.com/dracory/versionstore"
)

// == TYPE ====================================================================

// store represents the core structure for managing CMS data.
type storeImplementation struct {
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

	// Custom Entities
	customEntitiesEnabled bool
	customEntityStore     *CustomEntityStore

	// Shortcodes
	shortcodes  []ShortcodeInterface
	middlewares []MiddlewareInterface
}

// == INTERFACE ===============================================================

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

// PUBLIC METHODS ============================================================

// AutoMigrate performs automatic database migrations.
func (store *storeImplementation) AutoMigrate(ctx context.Context, opts ...Option) error {
	if store.db == nil {
		return errors.New("cms store: database is nil")
	}

	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	transaction, hasTransaction := options.params["tx"].(*sql.Tx)
	isDryRun, hasDryRun := options.params["dryRun"].(bool)

	blockSql, err := store.blockTableCreateSql()
	if err != nil {
		return err
	}
	pageSql, err := store.pageTableCreateSql()
	if err != nil {
		return err
	}
	tableSql, err := store.siteTableCreateSql()
	if err != nil {
		return err
	}
	templateSql, err := store.templateTableCreateSql()
	if err != nil {
		return err
	}

	var menuSql, menuItemSql, translationSql string
	if store.menusEnabled {
		menuSql, err = store.menuTableCreateSql()
		if err != nil {
			return err
		}
		menuItemSql, err = store.menuItemTableCreateSql()
		if err != nil {
			return err
		}
	}

	if store.translationsEnabled {
		translationSql, err = store.translationTableCreateSql()
		if err != nil {
			return err
		}
	}

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
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// MenusEnabled checks if menus are enabled.
func (store *storeImplementation) MenusEnabled() bool {
	return store.menusEnabled
}

// TranslationsEnabled checks if translations are enabled.
func (store *storeImplementation) TranslationsEnabled() bool {
	return store.translationsEnabled
}

// VersioningEnabled checks if versioning is enabled.
func (store *storeImplementation) VersioningEnabled() bool {
	return store.versioningEnabled
}

// CustomEntitiesEnabled checks if custom entities are enabled.
func (store *storeImplementation) CustomEntitiesEnabled() bool {
	return store.customEntitiesEnabled
}

// CustomEntityStore returns the custom entity store.
func (store *storeImplementation) CustomEntityStore() *CustomEntityStore {
	return store.customEntityStore
}

// Shortcodes returns the list of shortcodes.
func (store *storeImplementation) Shortcodes() []ShortcodeInterface {
	return store.shortcodes
}

// AddShortcode adds a shortcode to the store.
func (store *storeImplementation) AddShortcode(shortcode ShortcodeInterface) {
	store.shortcodes = append(store.shortcodes, shortcode)
}

// AddShortcodes adds multiple shortcodes to the store.
func (store *storeImplementation) AddShortcodes(shortcodes []ShortcodeInterface) {
	store.shortcodes = append(store.shortcodes, shortcodes...)
}

// SetShortcodes sets the list of shortcodes.
func (store *storeImplementation) SetShortcodes(shortcodes []ShortcodeInterface) {
	store.shortcodes = shortcodes
}

// Middlewares returns the list of middlewares.
func (store *storeImplementation) Middlewares() []MiddlewareInterface {
	return store.middlewares
}

// AddMiddleware adds a middleware to the store.
func (store *storeImplementation) AddMiddleware(middleware MiddlewareInterface) {
	store.middlewares = append(store.middlewares, middleware)
}

// AddMiddlewares adds multiple middlewares to the store.
func (store *storeImplementation) AddMiddlewares(middlewares []MiddlewareInterface) {
	store.middlewares = append(store.middlewares, middlewares...)
}

// SetMiddlewares sets the list of middlewares.
func (store *storeImplementation) SetMiddlewares(middlewares []MiddlewareInterface) {
	store.middlewares = middlewares
}

// toQuerableContext converts a context to a queryable context.
func (store *storeImplementation) toQuerableContext(ctx context.Context) database.QueryableContext {
	if database.IsQueryableContext(ctx) {
		return ctx.(database.QueryableContext)
	}

	return database.Context(ctx, store.db)
}

func (store *storeImplementation) withTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	if !store.VersioningEnabled() || database.IsQueryableContext(ctx) {
		return fn(ctx)
	}

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	txCtx := database.Context(ctx, tx)
	if err := fn(txCtx); err != nil {
		return err
	}

	return tx.Commit()
}
