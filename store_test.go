package cmsstore

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
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

func TestStorePageCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		PageTableName:      "page_table_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	page := NewPage()

	err = store.PageCreate(page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStorePageFindByHandle(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		PageTableName:      "page_table_find_by_handle",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	page := NewPage().
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

	err = store.PageCreate(page)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	pageFound, errFind := store.PageFindByHandle(page.Handle())

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
		PageTableName:      "page_table_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	page := NewPage().
		SetStatus(PAGE_STATUS_ACTIVE)

	err = page.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.PageCreate(page)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	pageFound, errFind := store.PageFindByID(page.ID())

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
		PageTableName:      "page_table_soft_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	page := NewPage()

	err = store.PageCreate(page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.PageSoftDeleteByID(page.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if page.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Page MUST NOT be soft deleted")
	}

	pageFound, errFind := store.PageFindByID(page.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if pageFound != nil {
		t.Fatal("Page MUST be nil")
	}
	query := NewPageQuery().SetWithSoftDeleted(true)

	query, err = query.SetID(page.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query, err = query.SetLimit(1)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	pageFindWithDeleted, err := store.PageList(query)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(pageFindWithDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(pageFindWithDeleted[0].SoftDeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Page MUST be soft deleted", page.SoftDeletedAt())
	}
}

func TestStorePageDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		PageTableName:      "page_table_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	page := NewPage()

	err = store.PageCreate(page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.PageDeleteByID(page.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query := NewPageQuery().SetWithSoftDeleted(true)

	query, err = query.SetID(page.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query, err = query.SetLimit(1)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	pageFindWithDeleted, err := store.PageList(query)

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
		PageTableName:      "page_table_update",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	page := NewPage().
		SetStatus(PAGE_STATUS_ACTIVE)

	err = store.PageCreate(page)

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

	err = store.PageUpdate(page)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	pageFound, errFind := store.PageFindByID(page.ID())

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
