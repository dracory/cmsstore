package testutils

import (
	"context"
	"database/sql"
	"os"

	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/utils"
	_ "modernc.org/sqlite"
)

const SITE_01 = "SITE_01"
const SITE_02 = "SITE_02"
const PAGE_01 = "PAGE_01"
const PAGE_02 = "PAGE_02"
const TEMPLATE_01 = "TEMPLATE_01"
const TEMPLATE_02 = "TEMPLATE_02"
const BLOCK_01 = "BLOCK_01"
const BLOCK_02 = "BLOCK_02"

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

func InitStore(filepath string) (cmsstore.StoreInterface, error) {
	db := initDB(filepath)

	store, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
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

func SeedSite(store cmsstore.StoreInterface, siteID string) (cmsstore.SiteInterface, error) {
	site := cmsstore.NewSite().SetID(siteID).SetName(siteID).SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err := store.SiteCreate(context.Background(), site)

	return site, err
}

func SeedPage(store cmsstore.StoreInterface, siteID string, pageID string) (cmsstore.PageInterface, error) {
	page := cmsstore.NewPage().
		SetSiteID(siteID).
		SetID(pageID).
		SetName(pageID).
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err := store.PageCreate(context.Background(), page)

	return page, err
}
