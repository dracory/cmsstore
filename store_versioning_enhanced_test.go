package cmsstore

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dracory/versionstore"
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Version 1").
		SetContent("Content 1")

	// Create - should create version 1
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Update 1 - should create version 2
	page.SetTitle("Version 2")
	page.SetContent("Content 2")
	err = store.PageUpdate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Update 2 - should create version 3
	page.SetTitle("Version 3")
	page.SetContent("Content 3")
	err = store.PageUpdate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Query versions
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("ASC"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) < 2 {
		t.Fatalf("Should have at least 2 versions, got %d", len(versions))
	}

	// Verify version content progression
	if len(versions) >= 2 {
		// First version should contain "Version 1" or "Content 1"
		if !strings.Contains(versions[0].Content(), "Version 1") {
			t.Errorf("Expected first version to contain 'Version 1', got %s", versions[0].Content())
		}

		// Later versions should contain updated content
		lastVersion := versions[len(versions)-1]
		if !strings.Contains(lastVersion.Content(), "Version 3") {
			t.Errorf("Expected last version to contain 'Version 3', got %s", lastVersion.Content())
		}
	}
}

// TestVersioningNoChangeNoVersion tests that updating without changes doesn't create new version
func TestVersioningNoChangeNoVersion(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                  db,
		BlockTableName:      "block_table_no_change",
		PageTableName:       "page_table_no_change",
		SiteTableName:       "site_table_no_change",
		TemplateTableName:   "template_table_no_change",
		VersioningEnabled:   true,
		VersioningTableName: "version_table_no_change",
		AutomigrateEnabled:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Unchanged Title").
		SetContent("Unchanged Content")

	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get initial version count
	versions1, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	initialCount := len(versions1)

	// Update without changing anything
	err = store.PageUpdate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Version count should remain the same
	versions2, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if initialCount != len(versions2) {
		t.Errorf("Expected version count to remain the same, got %d instead of %d", len(versions2), initialCount)
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()
	page := NewPage().
		SetSiteID("test-site").
		SetTitle("Test Title").
		SetContent("Test Content").
		SetAlias("/test-alias")

	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get the version
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(versions))
	}

	// Parse version content
	var pageData map[string]interface{}
	err = json.Unmarshal([]byte(versions[0].Content()), &pageData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify key fields are present
	if pageData["title"] != "Test Title" {
		t.Errorf("Expected title to be 'Test Title', got %v", pageData["title"])
	}
	if pageData["content"] != "Test Content" {
		t.Errorf("Expected content to be 'Test Content', got %v", pageData["content"])
	}
	if pageData["alias"] != "/test-alias" {
		t.Errorf("Expected alias to be '/test-alias', got %v", pageData["alias"])
	}
	if pageData["id"] != page.ID() {
		t.Errorf("Expected id to be %v, got %v", page.ID(), pageData["id"])
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create and version a page
	page := NewPage().SetSiteID("site1").SetTitle("Page")
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create and version a template
	template := NewTemplate().SetSiteID("site1").SetName("Template")
	err = store.TemplateCreate(ctx, template)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create and version a site
	site := NewSite().SetName("Site")
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify each has versions
	pageVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pageVersions) != 1 {
		t.Errorf("Expected 1 page version, got %d", len(pageVersions))
	}

	templateVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_TEMPLATE).
		SetEntityID(template.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(templateVersions) != 1 {
		t.Errorf("Expected 1 template version, got %d", len(templateVersions))
	}

	siteVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_SITE).
		SetEntityID(site.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(siteVersions) != 1 {
		t.Errorf("Expected 1 site version, got %d", len(siteVersions))
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create multiple pages with versions
	for i := 0; i < 3; i++ {
		page := NewPage().SetSiteID("site1").SetTitle("Page")
		err = store.PageCreate(ctx, page)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Test entity type filter
	allPageVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(allPageVersions) < 3 {
		t.Errorf("Expected at least 3 page versions, got %d", len(allPageVersions))
	}

	// Test limit
	limitedVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetLimit(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(limitedVersions) != 2 {
		t.Errorf("Expected 2 limited versions, got %d", len(limitedVersions))
	}

	// Test offset
	offsetVersions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetOffset(1).
		SetLimit(2))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(offsetVersions) != 2 {
		t.Errorf("Expected 2 offset versions, got %d", len(offsetVersions))
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create a page
	page := NewPage().SetSiteID("site1").SetTitle("Page")
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Update multiple times
	for i := 0; i < 5; i++ {
		page.SetTitle("Updated")
		err = store.PageUpdate(ctx, page)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// List versions for this page
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) < 1 {
		t.Errorf("Expected at least 1 version, got %d", len(versions))
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	// Create a version directly
	version := NewVersioning().
		SetEntityType("test-type").
		SetEntityID("test-id").
		SetContent("test-content")

	err = store.VersioningCreate(ctx, version)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Soft delete it
	err = store.VersioningSoftDeleteByID(ctx, version.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should not find it normally
	found, err := store.VersioningFindByID(ctx, version.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found != nil {
		t.Errorf("Expected version to be nil, got %v", found)
	}

	// Should find it with soft delete included
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetID(version.ID()).
		SetSoftDeletedIncluded(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(versions))
	}
	if versions[0].SoftDeletedAt() == "" {
		t.Errorf("Expected soft deleted at to be not empty, got empty")
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get version
	versions, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(versions))
	}

	// Parse content
	var pageData map[string]interface{}
	err = json.Unmarshal([]byte(versions[0].Content()), &pageData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify metadata is in version
	if !strings.Contains(versions[0].Content(), "custom_field_1") {
		t.Errorf("Expected version to contain 'custom_field_1', got %s", versions[0].Content())
	}
	if !strings.Contains(versions[0].Content(), "value1") {
		t.Errorf("Expected version to contain 'value1', got %s", versions[0].Content())
	}
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
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx := context.Background()

	page := NewPage().SetSiteID("site1").SetTitle("V1")
	err = store.PageCreate(ctx, page)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Create multiple versions
	for i := 2; i <= 5; i++ {
		page.SetTitle("V" + string(rune('0'+i)))
		err = store.PageUpdate(ctx, page)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// Get versions in ascending order
	versionsAsc, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("ASC"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Get versions in descending order
	versionsDesc, err := store.VersioningList(ctx, NewVersioningQuery().
		SetEntityType(VERSIONING_TYPE_PAGE).
		SetEntityID(page.ID()).
		SetOrderBy(versionstore.COLUMN_CREATED_AT).
		SetSortOrder("DESC"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both should have same count
	if len(versionsAsc) != len(versionsDesc) {
		t.Errorf("Expected same count for ascending and descending versions, got %d and %d", len(versionsAsc), len(versionsDesc))
	}

	// Verify ordering is reversed - compare created_at timestamps
	if len(versionsAsc) > 0 && len(versionsDesc) > 0 {
		// First of ascending should have earliest timestamp
		// Last of descending should have earliest timestamp
		// Compare the timestamps - Since CreatedAt returns strings, compare them directly
		if versionsAsc[0].CreatedAt() > versionsDesc[len(versionsDesc)-1].CreatedAt() {
			t.Errorf("Expected ascending first version to have earlier timestamp, got %v and %v", versionsAsc[0].CreatedAt(), versionsDesc[len(versionsDesc)-1].CreatedAt())
		}

		// Last of ascending should have latest timestamp
		// First of descending should have latest timestamp
		if versionsAsc[len(versionsAsc)-1].CreatedAt() < versionsDesc[0].CreatedAt() {
			t.Errorf("Expected ascending last version to have later timestamp, got %v and %v", versionsAsc[len(versionsAsc)-1].CreatedAt(), versionsDesc[0].CreatedAt())
		}
	}
}
