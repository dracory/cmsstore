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
	pageSql := store.pageTableCreateSql()
	tableSql := store.siteTableCreateSql()
	templateSql := store.templateTableCreateSql()

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

	for _, sql := range []string{blockSql, pageSql, tableSql, templateSql} {
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
