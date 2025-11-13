package cmsstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dracory/sb"
)

func TestStorePageCreate(t *testing.T) {
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

	page := NewPage().SetSiteID("Site1")

	ctx := context.Background()
	err = store.PageCreate(ctx, page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStorePageFindByHandle(t *testing.T) {
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

	page := NewPage().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE).
		SetHandle("test-handle")

	err = page.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	pageFound, errFind := store.PageFindByHandle(ctx, page.Handle())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if pageFound == nil {
		t.Fatal("Page MUST NOT be nil")
	}

	if pageFound.ID() != page.ID() {
		t.Fatal("IDs do not match")
	}

	if pageFound.Status() != page.Status() {
		t.Fatal("Statuses do not match")
	}

	if pageFound.Meta("education_1") != page.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if pageFound.Meta("education_2") != page.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if pageFound.Meta("education_3") != page.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStorePageFindByID(t *testing.T) {
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

	page := NewPage().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE)

	err = page.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	pageFound, errFind := store.PageFindByID(ctx, page.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if pageFound == nil {
		t.Fatal("Page MUST NOT be nil")
	}

	if pageFound.ID() != page.ID() {
		t.Fatal("IDs do not match")
	}

	if pageFound.Status() != page.Status() {
		t.Fatal("Statuses do not match")
	}

	if pageFound.Meta("education_1") != page.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if pageFound.Meta("education_2") != page.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if pageFound.Meta("education_3") != page.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStorePageSoftDelete(t *testing.T) {
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

	page := NewPage().
		SetSiteID("Site1")

	ctx := context.Background()
	err = store.PageCreate(ctx, page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.PageSoftDeleteByID(ctx, page.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if page.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Page MUST NOT be soft deleted")
	}

	pageFound, errFind := store.PageFindByID(ctx, page.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if pageFound != nil {
		t.Fatal("Page MUST be nil")
	}

	pageFindWithSoftDeleted, err := store.PageList(ctx, PageQuery().
		SetSoftDeletedIncluded(true).
		SetID(page.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(pageFindWithSoftDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(pageFindWithSoftDeleted[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Page MUST be soft deleted", page.SoftDeletedAt())
	}

	if !pageFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("Page MUST be soft deleted")
	}
}

func TestStorePageDelete(t *testing.T) {
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

	page := NewPage().
		SetSiteID("Site1")

	ctx := context.Background()
	err = store.PageCreate(ctx, page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.PageDeleteByID(ctx, page.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	pageFindWithDeleted, err := store.PageList(ctx, PageQuery().
		SetSoftDeletedIncluded(true).
		SetID(page.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(pageFindWithDeleted) != 0 {
		t.Fatal("Page MUST be deleted, but it is not")
	}
}

func TestStorePageUpdate(t *testing.T) {
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

	page := NewPage().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.PageCreate(ctx, page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	page.SetStatus(PAGE_STATUS_INACTIVE)

	err = page.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.PageUpdate(ctx, page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	pageFound, errFind := store.PageFindByID(ctx, page.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if pageFound == nil {
		t.Fatal("Page MUST NOT be nil")
	}

	if pageFound.Status() != PAGE_STATUS_INACTIVE {
		t.Fatal("Status MUST be INACTIVE, found: ", pageFound.Status())
	}

	metas, err := pageFound.Metas()

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
