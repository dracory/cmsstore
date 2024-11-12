package cmsstore

import (
	"database/sql"
	"errors"
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
}

// == INTERFACE ===============================================================

var _ StoreInterface = (*store)(nil) // verify it extends the interface

// PUBLIC METHODS ============================================================

// AutoMigrate auto migrate
func (store *store) AutoMigrate() error {
	if store.db == nil {
		return errors.New("cms store: database is nil")
	}

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
		_, err := store.db.Exec(sql)

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
