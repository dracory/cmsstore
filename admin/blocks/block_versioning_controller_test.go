package admin

import (
	"context"
	"io"
	"log/slog"
	"net/http"
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

func initBlockVersioningHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewBlockVersioningController(ui).Handler
}

func Test_BlockVersioningController_BlockIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "block id is required")
}

func Test_BlockVersioningController_BlockIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {"invalid-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "block not found")
}

func Test_BlockVersioningController_ListRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and block
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<p>Original content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	// Create some versioning entries
	versioning1 := cmsstore.NewVersioning()
	versioning1.SetEntityID(block.ID())
	versioning1.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning1.SetContent(`{"content": "<p>Version 1 content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning1)
	require.NoError(t, err)

	versioning2 := cmsstore.NewVersioning()
	versioning2.SetEntityID(block.ID())
	versioning2.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning2.SetContent(`{"content": "<p>Version 2 content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning2)
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Block Revisions")
	assert.Contains(t, body, "Version")
	assert.Contains(t, body, "Created")
	assert.Contains(t, body, "Actions")
}

func Test_BlockVersioningController_PreviewRevision(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and block
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Historical content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id":      {block.ID()},
			"versioning_id": {versioning.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Block Revision")
	assert.Contains(t, body, "Historical content")
}

func Test_BlockVersioningController_RestoreAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and block
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	// Create a versioning entry with different content
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Restored content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":            {block.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {"content"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "revision attributes restored successfully")

	// Verify block content was restored
	restoredBlock, err := store.BlockFindByID(context.Background(), block.ID())
	require.NoError(t, err)
	assert.Equal(t, "<p>Restored content</p>", restoredBlock.Content())
}

func Test_BlockVersioningController_RestoreNoAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and block
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Restored content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":      {block.ID()},
			"versioning_id": {versioning.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "no revision attributes were selected")
}

func Test_BlockVersioningController_RestoreMultipleAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and block
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Original Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	block.SetContent("<p>Original content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	// Create a versioning entry with different values
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Restored content</p>", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":            {block.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {"content", "status"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "revision attributes restored successfully")
}

func Test_BlockVersioningController_VersioningNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and block
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id":      {block.ID()},
			"versioning_id": {"non-existent-versioning-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "versioning not found")
}

func Test_BlockVersioningController_EmptyRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and block (no versioning entries)
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Version")
	assert.Contains(t, body, "Created")
	assert.Contains(t, body, "Actions")
}

func Test_BlockVersioningController_DifferentBlockTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initBlockVersioningHandler(store)

	// Test HTML block versioning
	htmlBlock := cmsstore.NewBlock()
	htmlBlock.SetName("HTML Block")
	htmlBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	htmlBlock.SetSiteID(site.ID())
	htmlBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	htmlBlock.SetContent("<p>HTML content</p>")
	err = store.BlockCreate(context.Background(), htmlBlock)
	require.NoError(t, err)

	// Create versioning for HTML block
	versioning1 := cmsstore.NewVersioning()
	versioning1.SetEntityID(htmlBlock.ID())
	versioning1.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning1.SetContent(`{"content": "<p>HTML version 1</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning1)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {htmlBlock.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Revisions")
	assert.Contains(t, body, "Preview")

	// Test Navbar block versioning
	navbarBlock := cmsstore.NewBlock()
	navbarBlock.SetName("Navbar Block")
	navbarBlock.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	navbarBlock.SetSiteID(site.ID())
	navbarBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), navbarBlock)
	require.NoError(t, err)

	// Create versioning for Navbar block
	versioning2 := cmsstore.NewVersioning()
	versioning2.SetEntityID(navbarBlock.ID())
	versioning2.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning2.SetContent(`{"content": "navbar data"}`)
	err = store.VersioningCreate(context.Background(), versioning2)
	require.NoError(t, err)

	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {navbarBlock.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Revisions")
	assert.Contains(t, body, "Preview")

	// Test preview modal shows content for navbar block
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id":      {navbarBlock.ID()},
			"versioning_id": {versioning2.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "navbar data")
}
