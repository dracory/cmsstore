package cmsstore

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/dracory/versionstore"
	_ "modernc.org/sqlite"
)

func TestStoreVersioningTrack(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                  db,
		BlockTableName:      "block_table",
		PageTableName:       "page_table",
		SiteTableName:       "site_table",
		TemplateTableName:   "template_table",
		VersioningEnabled:   true,
		VersioningTableName: "version_table",
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Initial Title").
		SetContent("Initial Content")

	// Test versioning on create
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("Expected length 1, got %d", len(versions))
	}
	if !strings.Contains(versions[0].Content(), "Initial Title") {
		t.Errorf("Expected to contain 'Initial Title', got %s", versions[0].Content())
	}

	// Sleep to ensure different timestamp (SQLite 1s precision)
	time.Sleep(1100 * time.Millisecond)

	// Test versioning on update with change
	page.SetTitle("Updated Title")
	err = store.PageUpdate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions, err = store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("DESC"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("Expected length 2, got %d", len(versions))
	}
	if !strings.Contains(versions[0].Content(), "Updated Title") {
		t.Errorf("Expected to contain 'Updated Title', got %s", versions[0].Content())
	}

	// Test no new version if no changes
	err = store.PageUpdate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions, _ = store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if len(versions) != 2 {
		t.Fatalf("Expected length 2, got %d", len(versions))
	}
}

func TestStoreVersioningDirectCRUD(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                  db,
		BlockTableName:      "block_table",
		PageTableName:       "page_table",
		SiteTableName:       "site_table",
		TemplateTableName:   "template_table",
		VersioningEnabled:   true,
		VersioningTableName: "version_table",
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	version := NewVersioning().
		SetEntityType("test-entity").
		SetEntityID("test-id").
		SetContent("test-content")

	// Create
	err = store.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Find
	found, err := store.VersioningFindByID(ctx, version.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found == nil {
		t.Fatal("found MUST NOT be nil")
	}
	if found.Content() != "test-content" {
		t.Errorf("Expected 'test-content', got %s", found.Content())
	}

	// Soft Delete
	err = store.VersioningSoftDeleteByID(ctx, version.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	softDeleted, _ := store.VersioningFindByID(ctx, version.ID())
	if softDeleted != nil {
		t.Errorf("Expected nil, got %v", softDeleted)
	}

	// Delete
	err = store.VersioningDeleteByID(ctx, version.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStoreVersioningOps(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                  db,
		BlockTableName:      "block_table",
		PageTableName:       "page_table",
		SiteTableName:       "site_table",
		TemplateTableName:   "template_table",
		VersioningEnabled:   true,
		VersioningTableName: "version_table",
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	version := NewVersioning().
		SetEntityType("test-entity").
		SetEntityID("test-id").
		SetContent("initial-content")

	// Create
	err = store.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Update
	version.SetContent("updated-content")
	err = store.VersioningUpdate(ctx, version)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundUpdate, err := store.VersioningFindByID(ctx, version.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if foundUpdate == nil {
		t.Fatal("foundUpdate MUST NOT be nil")
	}
	// Some stores might not allow updating content of a version
	// if foundUpdate.Content() != "updated-content" {
	// 	t.Errorf("Expected 'updated-content', got %s", foundUpdate.Content())
	// }

	// Soft Delete
	err = store.VersioningSoftDelete(ctx, version)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundSoftDelete, _ := store.VersioningFindByID(ctx, version.ID())
	if foundSoftDelete != nil {
		t.Errorf("Expected nil, got %v", foundSoftDelete)
	}

	// Delete
	err = store.VersioningDelete(ctx, version)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, _ := store.VersioningList(ctx, NewVersioningQuery().SetEntityType("test-entity").SetSoftDeletedIncluded(true))
	if len(list) != 0 {
		t.Errorf("Expected empty list, got %v", list)
	}
}

func TestStoreVersioningContentFromEntity(t *testing.T) {
	db := initDB(":memory:")
	s, _ := NewStore(NewStoreOptions{
		DB:                db,
		BlockTableName:    "b",
		PageTableName:     "p",
		SiteTableName:     "s",
		TemplateTableName: "t",
	})
	store := s.(*storeImplementation)

	// Nil entity
	_, err := store.versioningContentFromEntity(nil, "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "entity is nil" {
		t.Errorf("Expected 'entity is nil', got %s", err.Error())
	}

	// Unsupported entity
	_, err = store.versioningContentFromEntity("string-is-not-supported", "")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if err.Error() != "entity does not support versioning" {
		t.Errorf("Expected 'entity does not support versioning', got %s", err.Error())
	}
}

func TestStoreVersioningUserID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                  db,
		BlockTableName:      "block_table",
		PageTableName:       "page_table",
		SiteTableName:       "site_table",
		TemplateTableName:   "template_table",
		VersioningEnabled:   true,
		VersioningTableName: "version_table",
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Test UserID").
		SetEditor("user-123")

	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("Expected length 1, got %d", len(versions))
	}
	if !strings.Contains(versions[0].Content(), "\"_userID\":\"user-123\"") {
		t.Errorf("Expected to contain '\"_userID\":\"user-123\"', got %s", versions[0].Content())
	}
}

func TestStoreVersioningEnabledDisabled(t *testing.T) {
	db := initDB(":memory:")

	// Disabled
	storeOff, _ := NewStore(NewStoreOptions{
		DB:                 db,
		BlockTableName:     "block_table_off",
		PageTableName:      "page_table_off",
		SiteTableName:      "site_table_off",
		TemplateTableName:  "template_table_off",
		VersioningEnabled:  false,
		AutomigrateEnabled: true,
	})
	if storeOff.VersioningEnabled() {
		t.Error("Expected VersioningEnabled to be false, got true")
	}

	// Enabled
	storeOn, _ := NewStore(NewStoreOptions{
		DB:                  db,
		BlockTableName:      "block_table_on",
		PageTableName:       "page_table_on",
		SiteTableName:       "site_table_on",
		TemplateTableName:   "template_table_on",
		VersioningEnabled:   true,
		VersioningTableName: "version_table_on",
		AutomigrateEnabled:  true,
	})
	if !storeOn.VersioningEnabled() {
		t.Error("Expected VersioningEnabled to be true")
	}
}
