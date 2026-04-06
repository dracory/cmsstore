package admin

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
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
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Block Manager") {
		t.Errorf("Expected body to contain 'Block Manager'")
	}
	if !strings.Contains(body, "No blocks found") {
		t.Errorf("Expected body to contain 'No blocks found'")
	}
	if !strings.Contains(body, "New Block") {
		t.Errorf("Expected body to contain 'New Block'")
	}
}

func Test_BlockManagerController_WithBlocks(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	block1 := cmsstore.NewBlock()
	block1.SetName("Header Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	block2 := cmsstore.NewBlock()
	block2.SetName("Footer Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), block2)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Block Manager") {
		t.Errorf("Expected body to contain 'Block Manager'")
	}
	if !strings.Contains(body, "Header Block") {
		t.Errorf("Expected body to contain 'Header Block'")
	}
	if !strings.Contains(body, "Footer Block") {
		t.Errorf("Expected body to contain 'Footer Block'")
	}
	if !strings.Contains(body, "block-update") {
		t.Errorf("Expected body to contain 'block-update'")
	}
	if !strings.Contains(body, "block-delete") {
		t.Errorf("Expected body to contain 'block-delete'")
	}
}

func Test_BlockManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"action": {"modal_block_filter_show"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Filters") {
		t.Errorf("Expected body to contain 'Filters'")
	}
	if !strings.Contains(body, "name=\"name\"") {
		t.Errorf("Expected body to contain 'name=\"name\"'")
	}
	if !strings.Contains(body, "name=\"type\"") {
		t.Errorf("Expected body to contain 'name=\"type\"'")
	}
	if !strings.Contains(body, "name=\"status\"") {
		t.Errorf("Expected body to contain 'name=\"status\"'")
	}
}

func Test_BlockManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	block1 := cmsstore.NewBlock()
	block1.SetName("A Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	block2 := cmsstore.NewBlock()
	block2.SetName("Z Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockManagerHandler(store)

	// Test sort by name ASC
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"sort_by":    {"name"},
			"sort_order": {"asc"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "A Block") {
		t.Errorf("Expected body to contain 'A Block'")
	}
	if !strings.Contains(body, "Z Block") {
		t.Errorf("Expected body to contain 'Z Block'")
	}

	// Test sort by name DESC
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"sort_by":    {"name"},
			"sort_order": {"desc"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Z Block") {
		t.Errorf("Expected body to contain 'Z Block'")
	}
	if !strings.Contains(body, "A Block") {
		t.Errorf("Expected body to contain 'A Block'")
	}
}

func Test_BlockManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	block1 := cmsstore.NewBlock()
	block1.SetName("Header Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	block2 := cmsstore.NewBlock()
	block2.SetName("Footer Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), block2)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockManagerHandler(store)

	// Test filter by name
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"name": {"Header"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Header Block") {
		t.Errorf("Expected body to contain 'Header Block'")
	}
	if strings.Contains(body, "Footer Block") {
		t.Errorf("Expected body not to contain 'Footer Block'")
	}

	// Test filter by status
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"status": {cmsstore.BLOCK_STATUS_ACTIVE},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Header Block") {
		t.Errorf("Expected body to contain 'Header Block'")
	}
	if strings.Contains(body, "Footer Block") {
		t.Errorf("Expected body not to contain 'Footer Block'")
	}

	// Test filter by type
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"type": {cmsstore.BLOCK_TYPE_HTML},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Header Block") {
		t.Errorf("Expected body to contain 'Header Block'")
	}
	if strings.Contains(body, "Footer Block") {
		t.Errorf("Expected body not to contain 'Footer Block'")
	}
}

func Test_BlockManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create many blocks to test pagination
	for i := 1; i <= 25; i++ {
		block := cmsstore.NewBlock()
		block.SetName(fmt.Sprintf("Block %d", i))
		block.SetType(cmsstore.BLOCK_TYPE_HTML)
		block.SetSiteID(site.ID())
		block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
		err = store.BlockCreate(context.Background(), block)
		if err != nil {
			t.Fatalf("Failed to create block: %v", err)
		}
	}

	handler := initBlockManagerHandler(store)

	// Test first page
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"1"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Block 1") {
		t.Errorf("Expected body to contain 'Block 1'")
	}
	if !strings.Contains(body, "pagination") {
		t.Errorf("Expected body to contain 'pagination'")
	}

	// Test second page
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"2"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Block 11") {
		t.Errorf("Expected body to contain 'Block 11'")
	}
}

func Test_BlockManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "No blocks found") {
		t.Errorf("Expected body to contain 'No blocks found'")
	}
	if !strings.Contains(body, "New Block") {
		t.Errorf("Expected body to contain 'New Block'")
	}
	if !strings.Contains(body, "block-create") {
		t.Errorf("Expected body to contain 'block-create'")
	}
}

func Test_BlockManagerController_TableActions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and block
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "block-update") {
		t.Errorf("Expected body to contain 'block-update'")
	}
	if !strings.Contains(body, "block-delete") {
		t.Errorf("Expected body to contain 'block-delete'")
	}
	if !strings.Contains(body, "block-versioning") {
		t.Errorf("Expected body to contain 'block-versioning'")
	}
}

func Test_BlockManagerController_MultipleSites(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, "Site 1")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	site2 := cmsstore.NewSite()
	site2.SetName("Site 2")
	site2.SetDomainNames([]string{"site2.example.com"})
	err = store.SiteCreate(context.Background(), site2)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create blocks for different sites
	block1 := cmsstore.NewBlock()
	block1.SetName("Site 1 Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site1.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	block2 := cmsstore.NewBlock()
	block2.SetName("Site 2 Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site2.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Site 1 Block") {
		t.Errorf("Expected body to contain 'Site 1 Block'")
	}
	if !strings.Contains(body, "Site 2 Block") {
		t.Errorf("Expected body to contain 'Site 2 Block'")
	}
}

func Test_BlockManagerController_Search(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and blocks
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	block1 := cmsstore.NewBlock()
	block1.SetName("Searchable Block")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	block2 := cmsstore.NewBlock()
	block2.SetName("Other Block")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockManagerHandler(store)

	// Test search functionality
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"name": {"Searchable"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Searchable Block") {
		t.Errorf("Expected body to contain 'Searchable Block'")
	}
	if strings.Contains(body, "Other Block") {
		t.Errorf("Expected body not to contain 'Other Block'")
	}
}
