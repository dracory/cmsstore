package cmsstore

import (
	"database/sql"
	"os"

	"github.com/gouniverse/utils"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	if filepath != ":memory:" && utils.FileExists(filepath) {
		err := os.Remove(filepath) // remove database

		if err != nil {
			panic(err)
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}
