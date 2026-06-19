package cmsstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dracory/database"
	"github.com/dracory/neat"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
)

// == TYPE ====================================================================

// store represents the core structure for managing CMS data.
type storeImplementation struct {
	blockTableName     string
	pageTableName      string
	siteTableName      string
	templateTableName  string
	db                 *sql.DB
	neatDB             *neat.Database
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
	versioningStore *versioningStore

	// Custom Entities
	customEntitiesEnabled bool
	customEntityStore     *CustomEntityStore

	// Shortcodes
	shortcodes  []ShortcodeInterface
	middlewares []MiddlewareInterface

	// Pending versioning operations to execute after transaction commit
	pendingVersioningOps []pendingVersioningOp
}

type pendingVersioningOp struct {
	entityType string
	entityID   string
	entity     any
}

// == INTERFACE ===============================================================

var _ StoreInterface = (*storeImplementation)(nil) // verify it extends the interface

// PUBLIC METHODS ============================================================

// AutoMigrate performs automatic database migrations.
// Deprecated: Use MigrateUp instead.
func (store *storeImplementation) AutoMigrate(ctx context.Context, opts ...Option) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	return store.MigrateUp(ctx)
}

// MigrateUp creates the cms store tables
func (store *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	// Create page table
	if !store.neatDB.Schema().HasTable(store.pageTableName) {
		err := store.neatDB.Schema().Create(store.pageTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_SITE_ID, 40)
			table.String(COLUMN_STATUS, 40)
			table.String(COLUMN_ALIAS, 255)
			table.String(COLUMN_NAME, 255)
			table.String(COLUMN_TITLE, 255)
			table.Text(COLUMN_CONTENT)
			table.String(COLUMN_EDITOR, 40)
			table.String(COLUMN_TEMPLATE_ID, 40)
			table.String(COLUMN_CANONICAL_URL, 255)
			table.String(COLUMN_META_KEYWORDS, 255)
			table.String(COLUMN_META_DESCRIPTION, 255)
			table.String(COLUMN_META_ROBOTS, 255)
			table.String(COLUMN_HANDLE, 40)
			table.Text(COLUMN_MIDDLEWARES_AFTER)
			table.Text(COLUMN_MIDDLEWARES_BEFORE)
			table.Text(COLUMN_METAS)
			table.Text(COLUMN_MEMO)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			return err
		}
	}

	// Create block table
	if !store.neatDB.Schema().HasTable(store.blockTableName) {
		err := store.neatDB.Schema().Create(store.blockTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_SITE_ID, 40)
			table.String(COLUMN_PAGE_ID, 40)
			table.String(COLUMN_TEMPLATE_ID, 40)
			table.String(COLUMN_STATUS, 40)
			table.String(COLUMN_PARENT_ID, 40)
			table.Integer(COLUMN_SEQUENCE)
			table.String(COLUMN_TYPE, 100)
			table.String(COLUMN_NAME, 255)
			table.Text(COLUMN_CONTENT)
			table.String(COLUMN_EDITOR, 40)
			table.String(COLUMN_HANDLE, 40)
			table.Text(COLUMN_METAS)
			table.Text(COLUMN_MEMO)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			return err
		}
	}

	// Create site table
	if !store.neatDB.Schema().HasTable(store.siteTableName) {
		err := store.neatDB.Schema().Create(store.siteTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_STATUS, 40)
			table.String(COLUMN_NAME, 255)
			table.Text(COLUMN_DOMAIN_NAMES)
			table.String(COLUMN_HANDLE, 40)
			table.Text(COLUMN_METAS)
			table.Text(COLUMN_MEMO)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			return err
		}
	}

	// Create template table
	if !store.neatDB.Schema().HasTable(store.templateTableName) {
		err := store.neatDB.Schema().Create(store.templateTableName, func(table contractsschema.Blueprint) {
			table.String(COLUMN_ID, 40)
			table.Primary(COLUMN_ID)
			table.String(COLUMN_SITE_ID, 40)
			table.String(COLUMN_STATUS, 40)
			table.String(COLUMN_NAME, 255)
			table.Text(COLUMN_CONTENT)
			table.String(COLUMN_EDITOR, 40)
			table.String(COLUMN_HANDLE, 40)
			table.Text(COLUMN_METAS)
			table.Text(COLUMN_MEMO)
			table.DateTime(COLUMN_CREATED_AT)
			table.DateTime(COLUMN_UPDATED_AT)
			table.DateTime(COLUMN_SOFT_DELETED_AT)
		})
		if err != nil {
			return err
		}
	}

	// Create menu table if enabled
	if store.menusEnabled {
		if !store.neatDB.Schema().HasTable(store.menuTableName) {
			err := store.neatDB.Schema().Create(store.menuTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_ID, 40)
				table.Primary(COLUMN_ID)
				table.String(COLUMN_SITE_ID, 40)
				table.String(COLUMN_STATUS, 40)
				table.String(COLUMN_NAME, 255)
				table.String(COLUMN_HANDLE, 40)
				table.Text(COLUMN_METAS)
				table.Text(COLUMN_MEMO)
				table.DateTime(COLUMN_CREATED_AT)
				table.DateTime(COLUMN_UPDATED_AT)
				table.DateTime(COLUMN_SOFT_DELETED_AT)
			})
			if err != nil {
				return err
			}
		}

		// Create menu item table
		if !store.neatDB.Schema().HasTable(store.menuItemTableName) {
			err := store.neatDB.Schema().Create(store.menuItemTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_ID, 40)
				table.Primary(COLUMN_ID)
				table.String(COLUMN_MENU_ID, 40)
				table.String(COLUMN_STATUS, 40)
				table.String(COLUMN_NAME, 255)
				table.String(COLUMN_PARENT_ID, 40)
				table.Integer(COLUMN_SEQUENCE)
				table.String(COLUMN_PAGE_ID, 40)
				table.String(COLUMN_URL, 255)
				table.String(COLUMN_TARGET, 40)
				table.String(COLUMN_HANDLE, 40)
				table.Text(COLUMN_METAS)
				table.Text(COLUMN_MEMO)
				table.DateTime(COLUMN_CREATED_AT)
				table.DateTime(COLUMN_UPDATED_AT)
				table.DateTime(COLUMN_SOFT_DELETED_AT)
			})
			if err != nil {
				return err
			}
		}
	}

	// Create translation table if enabled
	if store.translationsEnabled {
		if !store.neatDB.Schema().HasTable(store.translationTableName) {
			err := store.neatDB.Schema().Create(store.translationTableName, func(table contractsschema.Blueprint) {
				table.String(COLUMN_ID, 40)
				table.Primary(COLUMN_ID)
				table.String(COLUMN_SITE_ID, 40)
				table.String(COLUMN_STATUS, 40)
				table.String(COLUMN_NAME, 255)
				table.String(COLUMN_HANDLE, 40)
				table.Text(COLUMN_CONTENT)
				table.Text(COLUMN_METAS)
				table.Text(COLUMN_MEMO)
				table.DateTime(COLUMN_CREATED_AT)
				table.DateTime(COLUMN_UPDATED_AT)
				table.DateTime(COLUMN_SOFT_DELETED_AT)
			})
			if err != nil {
				return err
			}
		}
	}

	if store.versioningEnabled {
		err := store.versioningStore.MigrateUp(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// MigrateDown drops the cms store tables
func (store *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if store.neatDB == nil {
		return errors.New("cms store: database is nil")
	}

	// Drop in reverse order of creation
	if store.translationsEnabled {
		if store.neatDB.Schema().HasTable(store.translationTableName) {
			err := store.neatDB.Schema().Drop(store.translationTableName)
			if err != nil {
				return err
			}
		}
	}

	if store.menusEnabled {
		if store.neatDB.Schema().HasTable(store.menuItemTableName) {
			err := store.neatDB.Schema().Drop(store.menuItemTableName)
			if err != nil {
				return err
			}
		}
		if store.neatDB.Schema().HasTable(store.menuTableName) {
			err := store.neatDB.Schema().Drop(store.menuTableName)
			if err != nil {
				return err
			}
		}
	}

	if store.neatDB.Schema().HasTable(store.templateTableName) {
		err := store.neatDB.Schema().Drop(store.templateTableName)
		if err != nil {
			return err
		}
	}

	if store.neatDB.Schema().HasTable(store.siteTableName) {
		err := store.neatDB.Schema().Drop(store.siteTableName)
		if err != nil {
			return err
		}
	}

	if store.neatDB.Schema().HasTable(store.pageTableName) {
		err := store.neatDB.Schema().Drop(store.pageTableName)
		if err != nil {
			return err
		}
	}

	if store.neatDB.Schema().HasTable(store.blockTableName) {
		err := store.neatDB.Schema().Drop(store.blockTableName)
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

	// Clear pending operations before starting new transaction
	store.pendingVersioningOps = nil

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	txCtx := database.Context(ctx, tx)
	if err := fn(txCtx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Execute pending versioning operations after successful commit
	if len(store.pendingVersioningOps) > 0 {
		for _, op := range store.pendingVersioningOps {
			content, err := store.versioningContentFromEntity(op.entity, "")
			if err != nil {
				// Log error but don't fail the transaction
				continue
			}
			store.versioningCreateIfChanged(ctx, op.entityType, op.entityID, content)
		}
		store.pendingVersioningOps = nil
	}

	return nil
}
