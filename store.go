package cmsstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gouniverse/versionstore"
)

// == TYPE ====================================================================

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

	versioningEnabled   bool
	versioningTableName string
	versioningStore     versionstore.StoreInterface
}

// == INTERFACE ===============================================================

var _ StoreInterface = (*store)(nil) // verify it extends the interface

// PUBLIC METHODS ============================================================

// AutoMigrate auto migrate
func (store *store) AutoMigrate(context context.Context, opts ...Option) error {
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

	if store.versioningEnabled && store.versioningTableName == "" {
		return errors.New("versioning table name is empty")
	}

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
			_, err := transaction.ExecContext(context, sql)

			if err != nil {
				return err
			}

			continue
		} else {
			_, err := store.db.ExecContext(context, sql)

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

// EnableDebug - enables the debug option
func (st *store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

func (store *store) MenusEnabled() bool {
	return store.menusEnabled
}

func (store *store) TranslationsEnabled() bool {
	return store.translationsEnabled
}

func (store *store) VersioningEnabled() bool {
	return store.versioningEnabled
}

func (store *store) VersioningCreate(version VersioningInterface) error {
	return store.versioningStore.VersionCreate(version)
}

func (store *store) VersioningDelete(version VersioningInterface) error {
	return store.versioningStore.VersionDelete(version)
}

func (store *store) VersioningDeleteByID(id string) error {
	return store.versioningStore.VersionDeleteByID(id)
}

func (store *store) VersioningFindByID(versioningID string) (VersioningInterface, error) {
	return store.versioningStore.VersionFindByID(versioningID)
}

func (store *store) VersioningList(query VersioningQueryInterface) ([]VersioningInterface, error) {
	list, err := store.versioningStore.VersionList(query)

	if err != nil {
		return nil, err
	}

	newlist := make([]VersioningInterface, len(list))

	for i, v := range list {
		newlist[i] = v
	}

	return newlist, nil
}

func (store *store) VersioningSoftDelete(versioning VersioningInterface) error {
	return store.versioningStore.VersionSoftDelete(versioning)
}

func (store *store) VersioningSoftDeleteByID(id string) error {
	return store.versioningStore.VersionSoftDeleteByID(id)
}

func (store *store) VersioningUpdate(version VersioningInterface) error {
	return store.versioningStore.VersionUpdate(version)
}
