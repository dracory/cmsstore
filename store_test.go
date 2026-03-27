package cmsstore

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	if filepath != ":memory:" && fileExists(filepath) {
		err := os.Remove(filepath) // remove database

		if err != nil {
			panic(err)
		}
	}

	// For in-memory databases, use cache=shared to allow concurrent access
	// from multiple goroutines using the same connection pool
	dsn := filepath
	if filepath == ":memory:" {
		dsn = "file::memory:?cache=shared"
	}
	dsn += "?parseTime=true"
	if filepath == ":memory:" {
		dsn = "file::memory:?cache=shared&parseTime=true"
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}

	// For in-memory databases, set connection pool to 1 to ensure
	// all goroutines share the same database instance
	if filepath == ":memory:" {
		db.SetMaxOpenConns(1)
	}

	return db
}

func initStore(filepath string) (StoreInterface, error) {
	db := initDB(filepath)

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table",
		PageTableName:      "page_table",
		SiteTableName:      "site_table",
		TemplateTableName:  "template_table",
		MenusEnabled:       true,
		MenuTableName:      "menu_table",
		MenuItemTableName:  "menu_item_table",
		AutomigrateEnabled: true,
	})

	if err != nil {
		return nil, err
	}

	return store, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
