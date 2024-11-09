package cmsstore

import (
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func TestStoreBlockCreate(t *testing.T) {
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

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("")

	err = store.BlockCreate(block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreBlockFindByHandle(t *testing.T) {
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

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("").
		SetStatus(PAGE_STATUS_ACTIVE).
		SetHandle("test-handle")

	err = block.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.BlockCreate(block)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	blockFound, errFind := store.BlockFindByHandle(block.Handle())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if blockFound == nil {
		t.Fatal("Block MUST NOT be nil")
	}

	if blockFound.ID() != block.ID() {
		t.Fatal("IDs do not match")
	}

	if blockFound.Status() != block.Status() {
		t.Fatal("Statuses do not match")
	}

	if blockFound.Meta("education_1") != block.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if blockFound.Meta("education_2") != block.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if blockFound.Meta("education_3") != block.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreBlockFindByID(t *testing.T) {
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

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("").
		SetStatus(PAGE_STATUS_ACTIVE)

	err = block.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.BlockCreate(block)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	blockFound, errFind := store.BlockFindByID(block.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if blockFound == nil {
		t.Fatal("Block MUST NOT be nil")
	}

	if blockFound.ID() != block.ID() {
		t.Fatal("IDs do not match")
	}

	if blockFound.Status() != block.Status() {
		t.Fatal("Statuses do not match")
	}

	if blockFound.Meta("education_1") != block.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if blockFound.Meta("education_2") != block.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if blockFound.Meta("education_3") != block.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreBlockSoftDelete(t *testing.T) {
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

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("")

	err = store.BlockCreate(block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.BlockSoftDeleteByID(block.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if block.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Block MUST NOT be soft deleted")
	}

	blockFound, errFind := store.BlockFindByID(block.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if blockFound != nil {
		t.Fatal("Block MUST be nil")
	}
	query := NewBlockQuery().SetWithSoftDeleted(true)

	query, err = query.SetID(block.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query, err = query.SetLimit(1)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	blockFindWithSoftDeleted, err := store.BlockList(query)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(blockFindWithSoftDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(blockFindWithSoftDeleted[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Block MUST be soft deleted", block.SoftDeletedAt())
	}

	if !blockFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("Block MUST be soft deleted")
	}
}

func TestStoreBlockDelete(t *testing.T) {
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

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("")

	err = store.BlockCreate(block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.BlockDeleteByID(block.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query := NewBlockQuery().SetWithSoftDeleted(true)

	query, err = query.SetID(block.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	query, err = query.SetLimit(1)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	blockFindWithDeleted, err := store.BlockList(query)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(blockFindWithDeleted) != 0 {
		t.Fatal("Block MUST be deleted, but it is not")
	}
}

func TestStoreBlockUpdate(t *testing.T) {
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

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("").
		SetStatus(PAGE_STATUS_ACTIVE)

	err = store.BlockCreate(block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	block.SetStatus(PAGE_STATUS_INACTIVE)

	err = block.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.BlockUpdate(block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	blockFound, errFind := store.BlockFindByID(block.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if blockFound == nil {
		t.Fatal("Block MUST NOT be nil")
	}

	if blockFound.Status() != PAGE_STATUS_INACTIVE {
		t.Fatal("Status MUST be INACTIVE, found: ", blockFound.Status())
	}

	metas, err := blockFound.Metas()

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
