package cmsstore

import (
	"context"
	"testing"
	"time"

	"github.com/dracory/versionstore"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestStoreVersioningTrack(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		VersioningEnabled:    true,
		VersioningTableName:  "version_table",
		AutomigrateEnabled:   true,
	})
	require.NoError(t, err)

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Initial Title").
		SetContent("Initial Content")

	// Test versioning on create
	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.NoError(t, err)
	require.Len(t, versions, 1)
	require.Contains(t, versions[0].Content(), "Initial Title")

	// Sleep to ensure different timestamp (SQLite 1s precision)
	time.Sleep(1100 * time.Millisecond)

	// Test versioning on update with change
	page.SetTitle("Updated Title")
	err = store.PageUpdate(ctx, page)
	require.NoError(t, err)

	versions, err = store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("DESC"))
	require.NoError(t, err)
	require.Len(t, versions, 2)
	require.Contains(t, versions[0].Content(), "Updated Title")

	// Test no new version if no changes
	err = store.PageUpdate(ctx, page)
	require.NoError(t, err)

	versions, _ = store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.Len(t, versions, 2)
}

func TestStoreVersioningDirectCRUD(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table",
		PageTableName:        "page_table",
		SiteTableName:        "site_table",
		TemplateTableName:    "template_table",
		VersioningEnabled:    true,
		VersioningTableName:  "version_table",
		AutomigrateEnabled:   true,
	})
	require.NoError(t, err)

	ctx := context.Background()
	version := NewVersioning().
		SetEntityType("test-entity").
		SetEntityID("test-id").
		SetContent("test-content")

	// Create
	err = store.VersioningCreate(ctx, version)
	require.NoError(t, err)

	// Find
	found, err := store.VersioningFindByID(ctx, version.ID())
	require.NoError(t, err)
	require.NotNil(t, found)
	require.Equal(t, "test-content", found.Content())

	// Soft Delete
	err = store.VersioningSoftDeleteByID(ctx, version.ID())
	require.NoError(t, err)

	softDeleted, _ := store.VersioningFindByID(ctx, version.ID())
	require.Nil(t, softDeleted)

	// Delete
	err = store.VersioningDeleteByID(ctx, version.ID())
	require.NoError(t, err)
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
	_, err := store.versioningContentFromEntity(nil)
	require.Error(t, err)
	require.Equal(t, "entity is nil", err.Error())

	// Unsupported entity
	_, err = store.versioningContentFromEntity("string-is-not-supported")
	require.Error(t, err)
	require.Equal(t, "entity does not support versioning", err.Error())
}

func TestStoreVersioningEnabledDisabled(t *testing.T) {
	db := initDB(":memory:")

	// Disabled
	storeOff, _ := NewStore(NewStoreOptions{
		DB:                  db,
		BlockTableName:      "block_table_off",
		PageTableName:       "page_table_off",
		SiteTableName:       "site_table_off",
		TemplateTableName:   "template_table_off",
		VersioningEnabled:   false,
		AutomigrateEnabled:  true,
	})
	require.False(t, storeOff.VersioningEnabled())

	// Enabled
	storeOn, _ := NewStore(NewStoreOptions{
		DB:                   db,
		BlockTableName:       "block_table_on",
		PageTableName:        "page_table_on",
		SiteTableName:        "site_table_on",
		TemplateTableName:    "template_table_on",
		VersioningEnabled:    true,
		VersioningTableName:  "version_table_on",
		AutomigrateEnabled:   true,
	})
	require.True(t, storeOn.VersioningEnabled())
}
