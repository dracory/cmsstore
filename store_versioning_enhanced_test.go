package cmsstore

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dracory/versionstore"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// TestVersioningMultipleUpdates tests multiple updates create separate versions
func TestVersioningMultipleUpdates(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Version 1").
		SetContent("Content 1")

	// Create - should create version 1
	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	// Update 1 - should create version 2
	page.SetTitle("Version 2")
	page.SetContent("Content 2")
	err = store.PageUpdate(ctx, page)
	require.NoError(t, err)

	// Update 2 - should create version 3
	page.SetTitle("Version 3")
	page.SetContent("Content 3")
	err = store.PageUpdate(ctx, page)
	require.NoError(t, err)

	// Query versions
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("ASC"))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(versions), 2, "Should have at least 2 versions")

	// Verify version content progression
	if len(versions) >= 2 {
		// First version should contain "Version 1" or "Content 1"
		require.Contains(t, versions[0].Content(), "Version 1")

		// Later versions should contain updated content
		lastVersion := versions[len(versions)-1]
		require.Contains(t, lastVersion.Content(), "Version 3")
	}
}

// TestVersioningNoChangeNoVersion tests that updating without changes doesn't create new version
func TestVersioningNoChangeNoVersion(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Unchanged Title").
		SetContent("Unchanged Content")

	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	// Get initial version count
	versions1, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.NoError(t, err)
	initialCount := len(versions1)

	// Update without changing anything
	err = store.PageUpdate(ctx, page)
	require.NoError(t, err)

	// Version count should remain the same
	versions2, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.NoError(t, err)
	require.Equal(t, initialCount, len(versions2), "No new version should be created for unchanged update")
}

// TestVersioningContentParsing tests that version content can be parsed back to entity
func TestVersioningContentParsing(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Test Title").
		SetContent("Test Content").
		SetAlias("/test-alias")

	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	// Get the version
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.NoError(t, err)
	require.Len(t, versions, 1)

	// Parse version content
	var pageData map[string]interface{}
	err = json.Unmarshal([]byte(versions[0].Content()), &pageData)
	require.NoError(t, err)

	// Verify key fields are present
	require.Equal(t, "Test Title", pageData["title"])
	require.Equal(t, "Test Content", pageData["content"])
	require.Equal(t, "/test-alias", pageData["alias"])
	require.Equal(t, page.ID(), pageData["id"])
}

// TestVersioningDifferentEntityTypes tests versioning for different entity types
func TestVersioningDifferentEntityTypes(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create and version a page
	page := NewPage().SetSiteID("site1").SetTitle("Page")
	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	// Create and version a template
	template := NewTemplate().SetSiteID("site1").SetName("Template")
	err = store.TemplateCreate(ctx, template)
	require.NoError(t, err)

	// Create and version a site
	site := NewSite().SetName("Site")
	err = store.SiteCreate(ctx, site)
	require.NoError(t, err)

	// Verify each has versions
	pageVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.NoError(t, err)
	require.Len(t, pageVersions, 1)

	templateVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_TEMPLATE).
		SetEntityID(template.ID()))
	require.NoError(t, err)
	require.Len(t, templateVersions, 1)

	siteVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_SITE).
		SetEntityID(site.ID()))
	require.NoError(t, err)
	require.Len(t, siteVersions, 1)
}

// TestVersioningQueryFilters tests various version query filters
func TestVersioningQueryFilters(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create multiple pages with versions
	for i := 0; i < 3; i++ {
		page := NewPage().SetSiteID("site1").SetTitle("Page")
		err = store.PageCreate(ctx, page)
		require.NoError(t, err)
	}

	// Test entity type filter
	allPageVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(allPageVersions), 3)

	// Test limit
	limitedVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetLimit(2))
	require.NoError(t, err)
	require.Len(t, limitedVersions, 2)

	// Test offset
	offsetVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetOffset(1).
		SetLimit(2))
	require.NoError(t, err)
	require.Len(t, offsetVersions, 2)
}

// TestVersioningListCount tests counting versions via list length
func TestVersioningListCount(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create a page
	page := NewPage().SetSiteID("site1").SetTitle("Page")
	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	// Update multiple times
	for i := 0; i < 5; i++ {
		page.SetTitle("Updated")
		err = store.PageUpdate(ctx, page)
		require.NoError(t, err)
	}

	// List versions for this page
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(versions), 1)
}

// TestVersioningSoftDelete tests soft deleting versions
func TestVersioningSoftDelete(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create a version directly
	version := NewVersioning().
		SetEntityType("test-type").
		SetEntityID("test-id").
		SetContent("test-content")

	err = store.VersioningCreate(ctx, version)
	require.NoError(t, err)

	// Soft delete it
	err = store.VersioningSoftDeleteByID(ctx, version.ID())
	require.NoError(t, err)

	// Should not find it normally
	found, err := store.VersioningFindByID(ctx, version.ID())
	require.NoError(t, err)
	require.Nil(t, found)

	// Should find it with soft delete included
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetID(version.ID()).
		SetSoftDeletedIncluded(true))
	require.NoError(t, err)
	require.Len(t, versions, 1)
	require.NotEmpty(t, versions[0].SoftDeletedAt())
}

// TestVersioningWithMetadata tests that metadata is captured in versions
func TestVersioningWithMetadata(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()

	// Create page with metadata
	page := NewPage().
		SetSiteID("site1").
		SetTitle("Page with Meta").
		SetEditor("user-123")

	err = page.SetMetas(map[string]string{
		"custom_field_1": "value1",
		"custom_field_2": "value2",
	})
	require.NoError(t, err)

	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	// Get version
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	require.NoError(t, err)
	require.Len(t, versions, 1)

	// Parse content
	var pageData map[string]interface{}
	err = json.Unmarshal([]byte(versions[0].Content()), &pageData)
	require.NoError(t, err)

	// Verify metadata is in version
	require.Contains(t, versions[0].Content(), "custom_field_1")
	require.Contains(t, versions[0].Content(), "value1")
}

// TestVersioningOrderBy tests ordering of versions
func TestVersioningOrderBy(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()

	page := NewPage().SetSiteID("site1").SetTitle("V1")
	err = store.PageCreate(ctx, page)
	require.NoError(t, err)

	// Create multiple versions
	for i := 2; i <= 5; i++ {
		page.SetTitle("V" + string(rune('0'+i)))
		err = store.PageUpdate(ctx, page)
		require.NoError(t, err)
	}

	// Get versions in ascending order
	versionsAsc, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("ASC"))
	require.NoError(t, err)

	// Get versions in descending order
	versionsDesc, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("DESC"))
	require.NoError(t, err)

	// Both should have same count
	require.Equal(t, len(versionsAsc), len(versionsDesc))

	// Verify ordering is reversed - compare created_at timestamps
	if len(versionsAsc) > 0 && len(versionsDesc) > 0 {
		// First of ascending should have earliest timestamp
		// Last of descending should have earliest timestamp
		// Compare the timestamps, not IDs (IDs may differ due to precision)
		require.Equal(t, versionsAsc[0].CreatedAt(), versionsDesc[len(versionsDesc)-1].CreatedAt())

		// Last of ascending should have latest timestamp
		// First of descending should have latest timestamp
		require.Equal(t, versionsAsc[len(versionsAsc)-1].CreatedAt(), versionsDesc[0].CreatedAt())
	}
}
