package cmsstore

import (
	"context"
	"strings"
	"testing"

	"github.com/dracory/sb"
)

func TestStoreMenuItemCreate(t *testing.T) {

	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	menuItem := NewMenuItem().
		SetMenuID("Menu1").
		SetParentID("0").
		SetSequence("0")

	ctx := context.Background()
	err = store.MenuItemCreate(ctx, menuItem)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreMenuItemFindByID(t *testing.T) {
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

	menuItem := NewMenuItem().
		SetMenuID("Menu1").
		SetStatus(MENU_ITEM_STATUS_ACTIVE).
		SetParentID("0").
		SetSequence("0")

	err = menuItem.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	ctx := context.Background()
	err = store.MenuItemCreate(ctx, menuItem)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	menuItemFound, errFind := store.MenuItemFindByID(ctx, menuItem.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if menuItemFound == nil {
		t.Fatal("MenuItem MUST NOT be nil")
	}

	if menuItemFound.ID() != menuItem.ID() {
		t.Fatal("IDs do not match")
	}

	if menuItemFound.Status() != menuItem.Status() {
		t.Fatal("Statuses do not match")
	}

	if menuItemFound.Meta("education_1") != menuItem.Meta("education_1") {
		t.Fatal("Metas do not match")
	}

	if menuItemFound.Meta("education_2") != menuItem.Meta("education_2") {
		t.Fatal("Metas do not match")
	}

	if menuItemFound.Meta("education_3") != menuItem.Meta("education_3") {
		t.Fatal("Metas do not match")
	}
}

func TestStoreMenuItemSoftDelete(t *testing.T) {
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

	menuItem := NewMenuItem().
		SetMenuID("Menu1").
		SetParentID("0").
		SetSequence("0")

	ctx := context.Background()
	err = store.MenuItemCreate(ctx, menuItem)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MenuItemSoftDeleteByID(ctx, menuItem.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if menuItem.SoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("MenuItem MUST NOT be soft deleted")
	}

	menuItemFound, errFind := store.MenuItemFindByID(ctx, menuItem.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if menuItemFound != nil {
		t.Fatal("MenuItem MUST be nil")
	}

	menuItemFindWithSoftDeleted, err := store.MenuItemList(ctx, MenuItemQuery().
		SetSoftDeletedIncluded(true).
		SetID(menuItem.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(menuItemFindWithSoftDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(menuItemFindWithSoftDeleted[0].SoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("MenuItem MUST be soft deleted", menuItem.SoftDeletedAt())
	}

	if !menuItemFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("MenuItem MUST be soft deleted")
	}
}

func TestStoreMenuItemDelete(t *testing.T) {
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

	menuItem := NewMenuItem().
		SetMenuID("Menu1").
		SetParentID("0").
		SetSequence("0")

	ctx := context.Background()
	err = store.MenuItemCreate(ctx, menuItem)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MenuItemDeleteByID(ctx, menuItem.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menuItemFindWithDeleted, err := store.MenuItemList(ctx, MenuItemQuery().
		SetSoftDeletedIncluded(true).
		SetID(menuItem.ID()).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(menuItemFindWithDeleted) != 0 {
		t.Fatal("MenuItem MUST be deleted, but it is not")
	}
}

func TestStoreMenuItemUpdate(t *testing.T) {
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

	menuItem := NewMenuItem().
		SetMenuID("Menu1").
		SetStatus(PAGE_STATUS_ACTIVE).
		SetParentID("0").
		SetSequence("0")

	ctx := context.Background()
	err = store.MenuItemCreate(ctx, menuItem)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menuItem.SetStatus(MENU_ITEM_STATUS_INACTIVE)

	err = menuItem.SetMetas(map[string]string{
		"education_1": "Education 1",
		"education_2": "Education 2",
		"education_3": "Education 3",
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.MenuItemUpdate(ctx, menuItem)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	menuItemFound, errFind := store.MenuItemFindByID(ctx, menuItem.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if menuItemFound == nil {
		t.Fatal("MenuItem MUST NOT be nil")
	}

	if menuItemFound.Status() != MENU_ITEM_STATUS_INACTIVE {
		t.Fatal("Status MUST be INACTIVE, found: ", menuItemFound.Status())
	}

	metas, err := menuItemFound.Metas()

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
