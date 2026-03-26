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

func initBlockUpdateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewBlockUpdateController(ui).Handler
}

func Test_BlockUpdateController_BlockIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "block id is required")
}

func Test_BlockUpdateController_BlockIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initBlockUpdateHandler(store)

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

func Test_BlockUpdateController_ViewSettings(t *testing.T) {
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
	block.SetContent("<p>Test content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
			"view":     {"settings"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Edit Block")
	assert.Contains(t, body, "block_name")
	assert.Contains(t, body, "block_type")
	assert.Contains(t, body, "block_site_id")
	assert.Contains(t, body, block.Name())
	assert.Contains(t, body, block.Type())
	assert.Contains(t, body, site.ID())
}

func Test_BlockUpdateController_ViewContent(t *testing.T) {
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
	block.SetContent("<p>Test content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
			"view":     {"content"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Edit Block")
	assert.Contains(t, body, "Content")
	assert.Contains(t, body, "block_content")
}

func Test_BlockUpdateController_UpdateSettings(t *testing.T) {
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
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":      {block.ID()},
			"block_name":    {"Updated Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_NAVBAR},
			"block_site_id": {site.ID()},
			"block_status":  {cmsstore.BLOCK_STATUS_ACTIVE},
			"view":          {"settings"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "block updated successfully")

	// Verify block was updated
	updatedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	require.NoError(t, err)
	assert.Equal(t, "Updated Block", updatedBlock.Name())
	assert.Equal(t, cmsstore.BLOCK_TYPE_NAVBAR, updatedBlock.Type())
	assert.Equal(t, cmsstore.BLOCK_STATUS_ACTIVE, updatedBlock.Status())
}

func Test_BlockUpdateController_UpdateContent(t *testing.T) {
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

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":      {block.ID()},
			"block_content": {"<div>Updated content</div>"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")

	// Verify block content was updated
	updatedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	require.NoError(t, err)
	assert.Equal(t, "<div>Updated content</div>", updatedBlock.Content())
}

func Test_BlockUpdateController_Update_ValidationError_MissingName(t *testing.T) {
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

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":   {block.ID()},
			"block_type": {cmsstore.BLOCK_TYPE_HTML},
			"site_id":    {site.ID()},
			"view":       {"settings"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_BlockUpdateController_Update_WithMemo(t *testing.T) {
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

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":   {block.ID()},
			"block_name": {"Updated Block"},
			"block_type": {cmsstore.BLOCK_TYPE_HTML},
			"site_id":    {site.ID()},
			"block_memo": {"Updated memo"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
}

func Test_BlockUpdateController_Update_WithHandle(t *testing.T) {
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

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":     {block.ID()},
			"block_name":   {"Updated Block"},
			"block_type":   {cmsstore.BLOCK_TYPE_HTML},
			"site_id":      {site.ID()},
			"block_handle": {"updated-block-handle"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
}

func Test_BlockUpdateController_BlockTypeChange(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and create a draft block
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	block.SetContent("<p>Original HTML content</p>")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

	handler := initBlockUpdateHandler(store)

	// Change block type from HTML to Navbar
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":     {block.ID()},
			"block_name":   {"Test Block"},
			"block_type":   {cmsstore.BLOCK_TYPE_NAVBAR},
			"block_status": {cmsstore.BLOCK_STATUS_DRAFT},
			"site_id":      {site.ID()},
			"view":         {"settings"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")

	// Verify block type was changed
	updatedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	require.NoError(t, err)
	assert.Equal(t, cmsstore.BLOCK_TYPE_NAVBAR, updatedBlock.Type())
}

func Test_BlockUpdateController_DifferentBlockTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initBlockUpdateHandler(store)

	// Test HTML block update
	htmlBlock := cmsstore.NewBlock()
	htmlBlock.SetName("HTML Block")
	htmlBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	htmlBlock.SetSiteID(site.ID())
	htmlBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), htmlBlock)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":      {htmlBlock.ID()},
			"block_name":    {"Updated HTML Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"site_id":       {site.ID()},
			"block_content": {"<div>Updated HTML</div>"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")

	// Test Navbar block update
	navbarBlock := cmsstore.NewBlock()
	navbarBlock.SetName("Navbar Block")
	navbarBlock.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	navbarBlock.SetSiteID(site.ID())
	navbarBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), navbarBlock)
	require.NoError(t, err)

	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":   {navbarBlock.ID()},
			"block_name": {"Updated Navbar Block"},
			"block_type": {cmsstore.BLOCK_TYPE_NAVBAR},
			"site_id":    {site.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")
}
