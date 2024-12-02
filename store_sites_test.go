package cmsstore

import (
	"context"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func TestStoreSiteCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_create",
		PageTableName:      "page_table_create",
		SiteTableName:      "site_table_create",
		TemplateTableName:  "template_table_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	site := NewSite()

	ctx := context.Background()
	err = store.SiteCreate(ctx, site)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreSiteFindByHandle(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_find_by_handle",
		PageTableName:      "page_table_find_by_handle",
		SiteTableName:      "site_table_find_by_handle",
		TemplateTableName:  "template_table_find_by_handle",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	site := NewSite().
		SetStatus(PAGE_STATUS_ACTIVE).
		SetHandle("test-handle")

	err = site.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	siteFound, errFind := store.SiteFindByHandle(ctx, site.Handle())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if siteFound == nil {
		t.Fatal("Site MUST NOT be nil")
	}

	if siteFound.ID() != site.ID() {
		t.Fatal("IDs do not match")
	}

	if siteFound.Status() != site.Status() {
		t.Fatal("Statuses do not match")
	}

	if siteFound.Meta("education_1") != site.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if siteFound.Meta("education_2") != site.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if siteFound.Meta("education_3") != site.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreSiteFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_find_by_id",
		PageTableName:      "page_table_find_by_id",
		SiteTableName:      "site_table_find_by_id",
		TemplateTableName:  "template_table_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	site := NewSite().
		SetStatus(PAGE_STATUS_ACTIVE)

	err = site.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	siteFound, errFind := store.SiteFindByID(ctx, site.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if siteFound == nil {
		t.Fatal("Site MUST NOT be nil")
	}

	if siteFound.ID() != site.ID() {
		t.Fatal("IDs do not match")
	}

	if siteFound.Status() != site.Status() {
		t.Fatal("Statuses do not match")
	}

	if siteFound.Meta("education_1") != site.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if siteFound.Meta("education_2") != site.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if siteFound.Meta("education_3") != site.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreSiteSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_soft_delete",
		PageTableName:      "page_table_soft_delete",
		SiteTableName:      "site_table_soft_delete",
		TemplateTableName:  "template_table_soft_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	site := NewSite()

	ctx := context.Background()
	err = store.SiteCreate(ctx, site)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SiteSoftDeleteByID(ctx, site.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if site.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Site MUST NOT be soft deleted")
	}

	siteFound, errFind := store.SiteFindByID(ctx, site.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if siteFound != nil {
		t.Fatal("Site MUST be nil")
	}

	siteFindWithSoftDeleted, err := store.SiteList(ctx, SiteQuery().
		SetSoftDeletedIncluded(true).
		SetID(site.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(siteFindWithSoftDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(siteFindWithSoftDeleted[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Site MUST be soft deleted", site.SoftDeletedAt())
	}

	if !siteFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("Site MUST be soft deleted")
	}
}

func TestStoreSiteDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_delete",
		PageTableName:      "page_table_delete",
		SiteTableName:      "site_table_delete",
		TemplateTableName:  "template_table_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	site := NewSite()

	ctx := context.Background()
	err = store.SiteCreate(ctx, site)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SiteDeleteByID(ctx, site.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	siteFindWithDeleted, err := store.SiteList(ctx, SiteQuery().
		SetSoftDeletedIncluded(true).
		SetID(site.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(siteFindWithDeleted) != 0 {
		t.Fatal("Site MUST be deleted, but it is not")
	}
}

func TestStoreSiteUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_update",
		PageTableName:      "page_table_update",
		SiteTableName:      "site_table_update",
		TemplateTableName:  "template_table_update",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	site := NewSite().
		SetStatus(PAGE_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.SiteCreate(ctx, site)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	site.SetStatus(PAGE_STATUS_INACTIVE)

	err = site.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SiteUpdate(ctx, site)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	siteFound, errFind := store.SiteFindByID(ctx, site.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if siteFound == nil {
		t.Fatal("Site MUST NOT be nil")
	}

	if siteFound.Status() != PAGE_STATUS_INACTIVE {
		t.Fatal("Status MUST be INACTIVE, found: ", siteFound.Status())
	}

	metas, err := siteFound.Metas()

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(metas) < 3 {
		t.Fatal("Metas MUST be 3, found: ", len(metas))
	}

	if metas["education_1"] != "Education 1" {
		t.Fatal("Metas do not match")
	}

	if metas["education_2"] != "Education 2" {
		t.Fatal("Metas do not match")
	}

	if metas["education_3"] != "Education 3" {
		t.Fatal("Metas do not match")
	}
}
