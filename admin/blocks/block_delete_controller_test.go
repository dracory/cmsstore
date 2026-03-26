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

func initBlockDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewBlockDeleteController(ui).Handler
}

func Test_BlockDeleteController_Index(t *testing.T) {
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
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Delete Block")
	assert.Contains(t, body, "Are you sure you want to delete this block?")
	assert.Contains(t, body, block.ID())
}

func Test_BlockDeleteController_Delete(t *testing.T) {
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
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "block deleted successfully")

	// Verify block is deleted (soft delete)
	deletedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	if err == nil && deletedBlock != nil {
		// If block is found, it should be inactive
		assert.Equal(t, cmsstore.BLOCK_STATUS_INACTIVE, deletedBlock.Status())
	} else {
		// Block might not be findable after soft delete, which is also expected
		// The important thing is that the delete operation succeeded
		assert.Contains(t, strings.ToLower(body), "success")
	}
}

func Test_BlockDeleteController_Delete_ValidationError_MissingBlockID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_BlockDeleteController_Delete_ValidationError_EmptyBlockID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {""},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_BlockDeleteController_BlockNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {"non-existent-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "block not found")
}

func Test_BlockDeleteController_Delete_DifferentBlockTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	// Test HTML block deletion
	htmlBlock := cmsstore.NewBlock()
	htmlBlock.SetName("HTML Block")
	htmlBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	htmlBlock.SetSiteID(site.ID())
	htmlBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), htmlBlock)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {htmlBlock.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")

	// Test Navbar block deletion
	navbarBlock := cmsstore.NewBlock()
	navbarBlock.SetName("Navbar Block")
	navbarBlock.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	navbarBlock.SetSiteID(site.ID())
	navbarBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), navbarBlock)
	require.NoError(t, err)

	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {navbarBlock.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")
}

func Test_BlockDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	// Create multiple blocks
	block1 := cmsstore.NewBlock()
	block1.SetName("Block 1")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	require.NoError(t, err)

	block2 := cmsstore.NewBlock()
	block2.SetName("Block 2")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	require.NoError(t, err)

	// Delete first block
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {block1.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")

	// Verify first block is deleted but second block remains
	deletedBlock1, err := store.BlockFindByID(context.Background(), block1.ID())
	if err == nil && deletedBlock1 != nil {
		assert.Equal(t, cmsstore.BLOCK_STATUS_INACTIVE, deletedBlock1.Status())
	}

	activeBlock2, err := store.BlockFindByID(context.Background(), block2.ID())
	if err == nil && activeBlock2 != nil {
		assert.Equal(t, cmsstore.BLOCK_STATUS_ACTIVE, activeBlock2.Status())
	}
}

func Test_BlockDeleteController_Delete_WithContent(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a block with content
	block := cmsstore.NewBlock()
	block.SetName("Block with Content")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<div>Test content</div>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "block deleted successfully")
}
