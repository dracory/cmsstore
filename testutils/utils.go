package testutils

import (
	"context"
	"database/sql"
	"os"

	"github.com/dracory/cmsstore"
)

const SITE_01 = "SITE_01"
const SITE_02 = "SITE_02"
const PAGE_01 = "PAGE_01"
const PAGE_02 = "PAGE_02"
const TEMPLATE_01 = "TEMPLATE_01"
const TEMPLATE_02 = "TEMPLATE_02"
const TRANSLATION_01 = "TRANSLATION_01"
const TRANSLATION_02 = "TRANSLATION_02"

func initDB(filepath string) *sql.DB {
	if filepath != ":memory:" && fileExists(filepath) {
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
		DB:                         db,
		BlockTableName:             "block_table",
		PageTableName:              "page_table",
		SiteTableName:              "site_table",
		TemplateTableName:          "template_table",
		MenusEnabled:               true,
		MenuTableName:              "menu_table",
		MenuItemTableName:          "menu_item_table",
		TranslationsEnabled:        true,
		TranslationTableName:       "translation_table",
		TranslationLanguageDefault: "en",
		VersioningEnabled:          true,
		VersioningTableName:        "versioning_table",
		AutomigrateEnabled:         true,
	})

	if err != nil {
		return nil, err
	}

	return store, nil
}

func SeedPage(store cmsstore.StoreInterface, siteID string, pageID string) (cmsstore.PageInterface, error) {
	page := cmsstore.NewPage().
		SetSiteID(siteID).
		SetName(pageID).
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err := store.PageCreate(context.Background(), page)

	return page, err
}

func SeedSite(store cmsstore.StoreInterface, siteID string) (cmsstore.SiteInterface, error) {
	site := cmsstore.NewSite().
		SetName(siteID).
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err := store.SiteCreate(context.Background(), site)

	return site, err
}

func SeedTemplate(store cmsstore.StoreInterface, siteID string, templateID string) (cmsstore.TemplateInterface, error) {
	template := cmsstore.NewTemplate().
		SetSiteID(siteID).
		SetName("Template" + templateID).
		SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)

	err := store.TemplateCreate(context.Background(), template)

	return template, err
}

func SeedTranslation(store cmsstore.StoreInterface, siteID string, translationID string) (cmsstore.TranslationInterface, error) {
	translation := cmsstore.NewTranslation().
		SetSiteID(siteID).
		SetName("Translation" + translationID).
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	err := store.TranslationCreate(context.Background(), translation)

	return translation, err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
