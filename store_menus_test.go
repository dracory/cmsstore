package cmsstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dracory/sb"
)

func TestStoreMenuCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_create",
		PageTableName:      "page_table_create",
		SiteTableName:      "site_table_create",
		TemplateTableName:  "template_table_create",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_create",
		MenuItemTableName:  "menu_item_table_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	menu := NewMenu().SetSiteID("Site1")

	ctx := context.Background()
	err = store.MenuCreate(ctx, menu)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreMenuFindByHandle(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_find_by_handle",
		PageTableName:      "page_table_find_by_handle",
		SiteTableName:      "site_table_find_by_handle",
		TemplateTableName:  "template_table_find_by_handle",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_find_by_handle",
		MenuItemTableName:  "menu_item_table_find_by_handle",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menu := NewMenu().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE).
		SetHandle("test-handle")

	err = menu.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.MenuCreate(ctx, menu)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menuFound, errFind := store.MenuFindByHandle(ctx, menu.Handle())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if menuFound == nil {
		t.Fatal("Menu MUST NOT be nil")
	}

	if menuFound.ID() != menu.ID() {
		t.Fatal("IDs do not match")
	}

	if menuFound.Status() != menu.Status() {
		t.Fatal("Statuses do not match")
	}

	if menuFound.Meta("education_1") != menu.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if menuFound.Meta("education_2") != menu.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if menuFound.Meta("education_3") != menu.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreMenuFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_find_by_id",
		PageTableName:      "page_table_find_by_id",
		SiteTableName:      "site_table_find_by_id",
		TemplateTableName:  "template_table_find_by_id",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_find_by_id",
		MenuItemTableName:  "menu_item_table_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menu := NewMenu().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE)

	err = menu.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.MenuCreate(ctx, menu)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	menuFound, errFind := store.MenuFindByID(ctx, menu.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if menuFound == nil {
		t.Fatal("Menu MUST NOT be nil")
	}

	if menuFound.ID() != menu.ID() {
		t.Fatal("IDs do not match")
	}

	if menuFound.Status() != menu.Status() {
		t.Fatal("Statuses do not match")
	}

	if menuFound.Meta("education_1") != menu.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if menuFound.Meta("education_2") != menu.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if menuFound.Meta("education_3") != menu.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreMenuSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_soft_delete",
		PageTableName:      "page_table_soft_delete",
		SiteTableName:      "site_table_soft_delete",
		TemplateTableName:  "template_table_soft_delete",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_soft_delete",
		MenuItemTableName:  "menu_item_table_soft_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	menu := NewMenu().
		SetSiteID("Site1")

	ctx := context.Background()
	err = store.MenuCreate(ctx, menu)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MenuSoftDeleteByID(ctx, menu.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if menu.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Menu MUST NOT be soft deleted")
	}

	menuFound, errFind := store.MenuFindByID(ctx, menu.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if menuFound != nil {
		t.Fatal("Menu MUST be nil")
	}

	menuFindWithSoftDeleted, err := store.MenuList(ctx, MenuQuery().
		SetSoftDeletedIncluded(true).
		SetID(menu.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(menuFindWithSoftDeleted) == 0 {
		t.Fatal("Menu MUST be soft deleted")
	}

	if strings.Contains(menuFindWithSoftDeleted[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Menu MUST be soft deleted", menuFindWithSoftDeleted[0].SoftDeletedAt())
	}

	if !menuFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("Menu MUST be soft deleted")
	}
}

func TestStoreMenuDeleteByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_delete",
		PageTableName:      "page_table_delete",
		SiteTableName:      "site_table_delete",
		TemplateTableName:  "template_table_delete",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_delete",
		MenuItemTableName:  "menu_item_table_delete",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	menu := NewMenu().
		SetSiteID("Site1")

	ctx := context.Background()
	err = store.MenuCreate(ctx, menu)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MenuDeleteByID(ctx, menu.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menuFindWithDeleted, err := store.MenuList(ctx, MenuQuery().
		SetSoftDeletedIncluded(true).
		SetID(menu.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(menuFindWithDeleted) != 0 {
		t.Fatal("Menu MUST be deleted, but it is not")
	}
}

func TestStoreMenuCount(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_count",
		PageTableName:      "page_table_count",
		SiteTableName:      "site_table_count",
		TemplateTableName:  "template_table_count",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_count",
		MenuItemTableName:  "menu_item_table_count",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	// Create 3 menus
	for i := 0; i < 3; i++ {
		menu := NewMenu().
			SetSiteID("Site1").
			SetStatus(PAGE_STATUS_ACTIVE)
		err = store.MenuCreate(ctx, menu)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
	}

	count, err := store.MenuCount(ctx, MenuQuery().SetSiteID("Site1"))
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if count != 3 {
		t.Fatalf("Expected count 3, got %d", count)
	}
}

func TestStoreMenuDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_delete_op",
		PageTableName:      "page_table_delete_op",
		SiteTableName:      "site_table_delete_op",
		TemplateTableName:  "template_table_delete_op",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_delete_op",
		MenuItemTableName:  "menu_item_table_delete_op",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()

	menu := NewMenu().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE).
		SetHandle("delete-me")

	err = store.MenuCreate(ctx, menu)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	// Delete by entity
	err = store.MenuDelete(ctx, menu)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	found, err := store.MenuFindByHandle(ctx, "delete-me")
	if err != nil && !strings.Contains(err.Error(), "not found") {
		t.Fatal("unexpected error:", err)
	}

	if found != nil {
		t.Fatal("Menu should have been deleted")
	}
}

func TestStoreMenuErrorPaths(t *testing.T) {
	ctx := context.Background()

	// Test with nil DB
	store := &storeImplementation{db: nil}

	_, err := store.MenuCount(ctx, MenuQuery())
	if err == nil {
		t.Error("Expected error for nil db")
	}

	err = store.MenuCreate(ctx, NewMenu())
	if err == nil {
		t.Error("Expected error for nil db")
	}

	err = store.MenuDelete(ctx, NewMenu())
	if err == nil {
		t.Error("Expected error for nil db")
	}

	err = store.MenuDeleteByID(ctx, "id")
	if err == nil {
		t.Error("Expected error for nil db")
	}

	_, err = store.MenuFindByHandle(ctx, "handle")
	if err == nil {
		t.Error("Expected error for nil db")
	}

	_, err = store.MenuFindByID(ctx, "id")
	if err == nil {
		t.Error("Expected error for nil db")
	}

	_, err = store.MenuList(ctx, MenuQuery())
	if err == nil {
		t.Error("Expected error for nil db")
	}

	err = store.MenuSoftDelete(ctx, NewMenu())
	if err == nil {
		t.Error("Expected error for nil db")
	}

	err = store.MenuSoftDeleteByID(ctx, "id")
	if err == nil {
		t.Error("Expected error for nil db")
	}

	err = store.MenuUpdate(ctx, NewMenu())
	if err == nil {
		t.Error("Expected error for nil db")
	}

	// Test with nil entity
	store.db = initDB(":memory:")
	err = store.MenuCreate(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil menu")
	}

	err = store.MenuDelete(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil menu")
	}

	err = store.MenuSoftDelete(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil menu")
	}

	err = store.MenuUpdate(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil menu")
	}

	// Test with empty ID/handle
	_, err = store.MenuFindByHandle(ctx, "")
	if err == nil {
		t.Error("Expected error for empty handle")
	}

	_, err = store.MenuFindByID(ctx, "")
	if err == nil {
		t.Error("Expected error for empty id")
	}

	err = store.MenuDeleteByID(ctx, "")
	if err == nil {
		t.Error("Expected error for empty id")
	}
}

func TestStoreMenuUpdate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_update",
		PageTableName:      "page_table_update",
		SiteTableName:      "site_table_update",
		TemplateTableName:  "template_table_update",
		MenusEnabled:       true,
		MenuTableName:      "menu_table_update",
		MenuItemTableName:  "menu_item_table_update",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	menu := NewMenu().
		SetSiteID("Site1").
		SetStatus(PAGE_STATUS_ACTIVE)

	ctx := context.Background()
	err = store.MenuCreate(ctx, menu)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menu.SetStatus(PAGE_STATUS_INACTIVE)

	err = menu.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MenuUpdate(ctx, menu)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menuFound, errFind := store.MenuFindByID(ctx, menu.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if menuFound == nil {
		t.Fatal("Menu MUST NOT be nil")
	}

	if menuFound.Status() != PAGE_STATUS_INACTIVE {
		t.Fatal("Status MUST be INACTIVE, found: ", menuFound.Status())
	}

	metas, err := menuFound.Metas()

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
