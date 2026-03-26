package admin

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func initBlockManagerHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewBlockManagerController(ui).Handler
}

func Test_BlockManagerController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Block Manager")
	assert.Contains(t, body, "No blocks found")
	assert.Contains(t, body, "New Block")
}

func Test_BlockManagerController_WithBlocks(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block1 := cmsstore.NewBlock()
	block1.SetName("Header Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	require.NoError(t, err)

	block2 := cmsstore.NewBlock()
	block2.SetName("Footer Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), block2)
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Block Manager")
	assert.Contains(t, body, "Header Block")
	assert.Contains(t, body, "Footer Block")
	assert.Contains(t, body, "block-update")
	assert.Contains(t, body, "block-delete")
}

func Test_BlockManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"action": {"modal_block_filter_show"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Filters")
	assert.Contains(t, body, "name=\"name\"")
	assert.Contains(t, body, "name=\"type\"")
	assert.Contains(t, body, "name=\"status\"")
}

func Test_BlockManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block1 := cmsstore.NewBlock()
	block1.SetName("A Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	require.NoError(t, err)

	block2 := cmsstore.NewBlock()
	block2.SetName("Z Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	// Test sort by name ASC
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"sort_by":    {"name"},
			"sort_order": {"asc"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "A Block")
	assert.Contains(t, body, "Z Block")

	// Test sort by name DESC
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"sort_by":    {"name"},
			"sort_order": {"desc"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Z Block")
	assert.Contains(t, body, "A Block")
}

func Test_BlockManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block1 := cmsstore.NewBlock()
	block1.SetName("Header Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	require.NoError(t, err)

	block2 := cmsstore.NewBlock()
	block2.SetName("Footer Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), block2)
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	// Test filter by name
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"name": {"Header"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Header Block")
	assert.NotContains(t, body, "Footer Block")

	// Test filter by status
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"status": {cmsstore.BLOCK_STATUS_ACTIVE},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Header Block")
	assert.NotContains(t, body, "Footer Block")

	// Test filter by type
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"type": {cmsstore.BLOCK_TYPE_HTML},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Header Block")
	assert.NotContains(t, body, "Footer Block")
}

func Test_BlockManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create many blocks to test pagination
	for i := 1; i <= 25; i++ {
		block := cmsstore.NewBlock()
		block.SetName(fmt.Sprintf("Block %d", i))
		block.SetType(cmsstore.BLOCK_TYPE_HTML)
		block.SetSiteID(site.ID())
		block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(context.Background(), block)
		require.NoError(t, err)
	}

	handler := initBlockManagerHandler(store)

	// Test first page
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"1"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Block 1")
	assert.Contains(t, body, "pagination")

	// Test second page
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"2"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Block 11")
}

func Test_BlockManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "No blocks found")
	assert.Contains(t, body, "New Block")
	assert.Contains(t, body, "block-create")
}

func Test_BlockManagerController_TableActions(t *testing.T) {
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

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "block-update")
	assert.Contains(t, body, "block-delete")
	assert.Contains(t, body, "block-versioning")
}

func Test_BlockManagerController_MultipleSites(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, "Site 1")
	require.NoError(t, err)

	site2 := cmsstore.NewSite()
	site2.SetName("Site 2")
	site2.SetDomainNames([]string{"site2.example.com"})
	err = store.SiteCreate(context.Background(), site2)
	require.NoError(t, err)

	// Create blocks for different sites
	block1 := cmsstore.NewBlock()
	block1.SetName("Site 1 Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site1.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	require.NoError(t, err)

	block2 := cmsstore.NewBlock()
	block2.SetName("Site 2 Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site2.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Site 1 Block")
	assert.Contains(t, body, "Site 2 Block")
}

func Test_BlockManagerController_Search(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block1 := cmsstore.NewBlock()
	block1.SetName("Searchable Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	require.NoError(t, err)

	block2 := cmsstore.NewBlock()
	block2.SetName("Other Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	require.NoError(t, err)

	handler := initBlockManagerHandler(store)

	// Test search functionality
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"name": {"Searchable"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Searchable Block")
	assert.NotContains(t, body, "Other Block")
}
