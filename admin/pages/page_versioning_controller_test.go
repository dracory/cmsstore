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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	handler := initVersioningHandler(store)

	// Seed a page
	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)
	require.NoError(t, err)

	// Create another version by updating
	seededPage.SetTitle("Updated Title")
	err = store.PageUpdate(context.Background(), seededPage)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {seededPage.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Check for modal elements
	assert.Contains(t, body, "ModalPageVersioning")
	assert.Contains(t, body, "Page Revisions")
	
	// Should contain two revisions (one from create, one from update)
	assert.Equal(t, 2, strings.Count(body, "Preview"), "Expected 2 preview buttons for 2 revisions")
}

func Test_PageVersioningController_PreviewRevision(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initVersioningHandler(store)

	// Seed a page (creates 1st version)
	seededPage, err := testutils.SeedPage(store, testutils.SITE_01, testutils.PAGE_01)
	require.NoError(t, err)

	// Get versions
	versions, err := store.VersioningList(context.Background(), cmsstore.NewVersioningQuery().
		SetEntityType(cmsstore.VERSIONING_TYPE_PAGE).
		SetEntityID(seededPage.ID()))
	require.NoError(t, err)
	require.Len(t, versions, 1)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id":       {seededPage.ID()},
			"versioning_id": {versions[0].ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Check for attribute table
	assert.Contains(t, body, "Attribute")
	assert.Contains(t, body, "Value")
	assert.Contains(t, body, "Apply")
	assert.Contains(t, body, "title")
	assert.Contains(t, body, "content")
	assert.Contains(t, body, "Restore Selected Attributes")
}

func Test_PageVersioningController_RestoreAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initVersioningHandler(store)

	// 1. Create a page manually with a title (Version 1)
	page := cmsstore.NewPage().
		SetID("page-to-restore").
		SetName("Original Name").
		SetTitle("Original Title").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE).
		SetSiteID(testutils.SITE_01)
	
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)
	v1ID := page.ID()

	// Get Version 1
	versions, err := store.VersioningList(context.Background(), cmsstore.NewVersioningQuery().SetEntityID(v1ID))
	require.NoError(t, err)
	require.NotEmpty(t, versions)
	version1ID := versions[0].ID()

	// 2. Update the page (Version 2: "Updated Title")
	page.SetTitle("Updated Title")
	err = store.PageUpdate(context.Background(), page)
	require.NoError(t, err)

	// Verify current title is "Updated Title"
	currentPage, _ := store.PageFindByID(context.Background(), v1ID)
	assert.Equal(t, "Updated Title", currentPage.Title())

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
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Check for success message (Swal)
	assert.Contains(t, body, "success")
	assert.Contains(t, body, "restored successfully")

	// 4. Verify title is restored in database
	restoredPage, _ := store.PageFindByID(context.Background(), v1ID)
	assert.Equal(t, "Original Title", restoredPage.Title(), "Title should be restored to original value")
}

func Test_PageVersioningController_RestoreNoAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

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

	assert.Contains(t, body, "error")
	assert.Contains(t, body, "No revision attributes were selected")
}
