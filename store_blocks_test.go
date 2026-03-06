package cmsstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dracory/sb"
	"github.com/stretchr/testify/require"
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

	ctx := context.Background()

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("").
		SetParentID("").
		SetSequenceInt(0)

	err = store.BlockCreate(ctx, block)

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

	ctx := context.Background()

	block := NewBlock().
		SetSiteID("Site1").
		SetPageID("").
		SetTemplateID("").
		SetParentID("").
		SetSequenceInt(0).
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

	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	blockFound, errFind := store.BlockFindByHandle(ctx, block.Handle())

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
		SetParentID("").
		SetSequenceInt(0).
		SetStatus(PAGE_STATUS_ACTIVE)

	err = block.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	blockFound, errFind := store.BlockFindByID(ctx, block.ID())

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
		SetTemplateID("").
		SetParentID("").
		SetSequenceInt(0)

	ctx := context.Background()
	err = store.BlockCreate(ctx, block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.BlockSoftDeleteByID(ctx, block.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if block.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Block MUST NOT be soft deleted")
	}

	blockFound, errFind := store.BlockFindByID(ctx, block.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if blockFound != nil {
		t.Fatal("Block MUST be nil")
	}

	blockFindWithSoftDeleted, err := store.BlockList(ctx, BlockQuery().
		SetID(block.ID()).
		SetSoftDeleteIncluded(true).
		SetLimit(1))

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

func TestStoreBlockDeleteByID(t *testing.T) {
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
		SetTemplateID("").
		SetParentID("").
		SetSequenceInt(0)

	ctx := context.Background()
	err = store.BlockCreate(ctx, block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.BlockDeleteByID(ctx, block.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	blockFindWithDeleted, err := store.BlockList(ctx, BlockQuery().
		SetID(block.ID()).
		SetLimit(1).
		SetSoftDeleteIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(blockFindWithDeleted) != 0 {
		t.Fatal("Block MUST be deleted, but it is not")
	}
}

func TestStoreBlockCount(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_count",
		PageTableName:      "page_table_count",
		SiteTableName:      "site_table_count",
		TemplateTableName:  "template_table_count",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create 3 blocks
	for i := 0; i < 3; i++ {
		block := NewBlock().
			SetSiteID("Site1").
			SetStatus(PAGE_STATUS_ACTIVE)
		err = store.BlockCreate(ctx, block)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err := store.BlockCount(ctx, BlockQuery().SetSiteID("Site1"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 3 {
		t.Fatalf("Expected count 3, got %d", count)
	}

	count, err = store.BlockCount(ctx, BlockQuery().SetSiteID("NonExistent"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 0 {
		t.Fatalf("Expected count 0, got %d", count)
	}
}

func TestStoreBlockDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_delete_op",
		PageTableName:      "page_table_delete_op",
		SiteTableName:      "site_table_delete_op",
		TemplateTableName:  "template_table_delete_op",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	block := NewBlock().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE).
		SetHandle("delete-me")

	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete by entity
	err = store.BlockDelete(ctx, block)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.BlockFindByHandle(ctx, "delete-me")
	if err != nil && !strings.Contains(err.Error(), "not found") {
		t.Fatal("unexpected error:", err)
	}

	if found != nil {
		t.Fatal("Block should have been deleted")
	}
}

func TestStoreBlockErrorPaths(t *testing.T) {
	ctx := context.Background()
	
	// Test with nil DB
	store := &storeImplementation{db: nil}
	
	_, err := store.BlockCount(ctx, BlockQuery())
	require.Error(t, err)
	require.Contains(t, err.Error(), "db is nil")

	err = store.BlockCreate(ctx, NewBlock())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	err = store.BlockDelete(ctx, NewBlock())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	err = store.BlockDeleteByID(ctx, "id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	_, err = store.BlockFindByHandle(ctx, "handle")
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	_, err = store.BlockFindByID(ctx, "id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	_, err = store.BlockList(ctx, BlockQuery())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	err = store.BlockSoftDelete(ctx, NewBlock())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	err = store.BlockSoftDeleteByID(ctx, "id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	err = store.BlockUpdate(ctx, NewBlock())
	require.Error(t, err)
	require.Contains(t, err.Error(), "database is nil")

	// Test with nil entity
	store.db = initDB(":memory:")
	err = store.BlockCreate(ctx, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "block is nil")

	err = store.BlockDelete(ctx, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "block is nil")

	err = store.BlockSoftDelete(ctx, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "block is nil")

	err = store.BlockUpdate(ctx, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "block is nil")

	// Test with empty ID/handle
	_, err = store.BlockFindByHandle(ctx, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "block handle is empty")

	_, err = store.BlockFindByID(ctx, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "block id is empty")

	err = store.BlockDeleteByID(ctx, "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "block id is empty")
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
		SetParentID("").
		SetSequenceInt(0).
		SetStatus(PAGE_STATUS_ACTIVE)

	ctx := context.Background()

	err = store.BlockCreate(ctx, block)

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

	err = store.BlockUpdate(ctx, block)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	blockFound, errFind := store.BlockFindByID(ctx, block.ID())

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
