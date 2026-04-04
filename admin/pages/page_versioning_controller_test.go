package admin

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

func initVersioningHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewPageVersioningController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_PageVersioningController_ListRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initVersioningHandler(store)

	// Seed a page
	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	// Create another version by updating
	seededPage.SetTitle("Updated Title")
	err = store.PageUpdate(context.Background(), seededPage)
	if err != nil {
		t.Fatalf("Failed to update page: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check for modal elements
	if !strings.Contains(body, "ModalPageVersioning") {
		t.Errorf("Expected body to contain 'ModalPageVersioning'")
	}
	if !strings.Contains(body, "Page Revisions") {
		t.Errorf("Expected body to contain 'Page Revisions'")
	}

	// Should contain two revisions (one from create, one from update)
	previewCount := strings.Count(body, "Preview")
	if previewCount != 2 {
		t.Errorf("Expected 2 preview buttons for 2 revisions, got %d", previewCount)
	}
}

func Test_PageVersioningController_PreviewRevision(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initVersioningHandler(store)

	// Seed a page (creates 1st version)
	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	// Get versions
	versions, err := store.VersioningList(context.Background(), cmsstore.NewVersioningQuery().
		SetEntityType(cmsstore.VERSIONING_TYPE_PAGE).
		SetEntityID(seededPage.ID()))
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("Expected 1 version, got %d", len(versions))
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id":       {seededPage.ID()},
			"versioning_id": {versions[0].ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check for attribute table
	if !strings.Contains(body, "Attribute") {
		t.Errorf("Expected body to contain 'Attribute'")
	}
	if !strings.Contains(body, "Value") {
		t.Errorf("Expected body to contain 'Value'")
	}
	if !strings.Contains(body, "Apply") {
		t.Errorf("Expected body to contain 'Apply'")
	}
	if !strings.Contains(body, "title") {
		t.Errorf("Expected body to contain 'title'")
	}
	if !strings.Contains(body, "content") {
		t.Errorf("Expected body to contain 'content'")
	}
	if !strings.Contains(body, "Restore Selected Attributes") {
		t.Errorf("Expected body to contain 'Restore Selected Attributes'")
	}
}

func Test_PageVersioningController_RestoreAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initVersioningHandler(store)

	// 1. Create a page manually with a title (Version 1)
	page := cmsstore.NewPage().
		SetID("page-to-restore").
		SetName("Original Name").
		SetTitle("Original Title").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE).
		SetSiteID(testutils.SITE_01)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}
	v1ID := page.ID()

	// Get Version 1
	versions, err := store.VersioningList(context.Background(), cmsstore.NewVersioningQuery().SetEntityID(v1ID))
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}
	if len(versions) == 0 {
		t.Fatalf("Expected at least 1 version, got 0")
	}
	version1ID := versions[0].ID()

	// 2. Update the page (Version 2: "Updated Title")
	page.SetTitle("Updated Title")
	err = store.PageUpdate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to update page: %v", err)
	}

	// Verify current title is "Updated Title"
	currentPage, _ := store.PageFindByID(context.Background(), v1ID)
	if currentPage.Title() != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%s'", currentPage.Title())
	}

	// 3. Restore title from Version 1 via POST
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id":       {v1ID},
			"versioning_id": {version1ID},
		},
		PostValues: url.Values{
			"revision_attributes": {cmsstore.COLUMN_TITLE},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check for success message (Swal)
	if !strings.Contains(body, "success") {
		t.Errorf("Expected body to contain 'success'")
	}
	if !strings.Contains(body, "restored successfully") {
		t.Errorf("Expected body to contain 'restored successfully'")
	}

	// 4. Verify title is restored in database
	restoredPage, _ := store.PageFindByID(context.Background(), v1ID)
	if restoredPage.Title() != "Original Title" {
		t.Errorf("Title should be restored to 'Original Title', got '%s'", restoredPage.Title())
	}
}

func Test_PageVersioningController_RestoreNoAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initVersioningHandler(store)

	seededPage, _ := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)
	versions, _ := store.VersioningList(context.Background(), cmsstore.NewVersioningQuery().SetEntityID(seededPage.ID()))

	body, _, _ := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id":       {seededPage.ID()},
			"versioning_id": {versions[0].ID()},
		},
		// No PostValues for revision_attributes
	})

	if !strings.Contains(body, "error") {
		t.Errorf("Expected body to contain 'error'")
	}
	if !strings.Contains(body, "No revision attributes were selected") {
		t.Errorf("Expected body to contain 'No revision attributes were selected'")
	}
}
