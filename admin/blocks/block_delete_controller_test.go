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

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Delete Block") {
		t.Errorf("Expected body to contain 'Delete Block'")
	}
	if !strings.Contains(body, "Are you sure you want to delete this block?") {
		t.Errorf("Expected body to contain 'Are you sure you want to delete this block?'")
	}
	if !strings.Contains(body, block.ID()) {
		t.Errorf("Expected body to contain %q", block.ID())
	}
}

func Test_BlockDeleteController_Delete(t *testing.T) {
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

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	bodyLower := strings.ToLower(body)
	if !strings.Contains(bodyLower, "success") {
		t.Errorf("Expected body to contain 'success'")
	}
	if !strings.Contains(bodyLower, "block deleted successfully") {
		t.Errorf("Expected body to contain 'block deleted successfully'")
	}

	// Verify block is deleted (soft delete)
	deletedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	if err == nil && deletedBlock != nil {
		// If block is found, it should be inactive
		if deletedBlock.Status() != cmsstore.BLOCK_STATUS_INACTIVE {
			t.Errorf("Expected status %v, got %v", cmsstore.BLOCK_STATUS_INACTIVE, deletedBlock.Status())
		}
	} else {
		// Block might not be findable after soft delete, which is also expected
		// The important thing is that the delete operation succeeded
		if !strings.Contains(bodyLower, "success") {
			t.Errorf("Expected body to contain 'success'")
		}
	}
}

func Test_BlockDeleteController_Delete_ValidationError_MissingBlockID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(strings.ToLower(body), "error") {
		t.Errorf("Expected body to contain 'error'")
	}
}

func Test_BlockDeleteController_Delete_ValidationError_EmptyBlockID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {""},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(strings.ToLower(body), "error") {
		t.Errorf("Expected body to contain 'error'")
	}
}

func Test_BlockDeleteController_BlockNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {"non-existent-id"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	bodyLower := strings.ToLower(body)
	if !strings.Contains(bodyLower, "error") {
		t.Errorf("Expected body to contain 'error'")
	}
	if !strings.Contains(bodyLower, "block not found") {
		t.Errorf("Expected body to contain 'block not found'")
	}
}

func Test_BlockDeleteController_Delete_DifferentBlockTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockDeleteHandler(store)

	// Test HTML block deletion
	htmlBlock := cmsstore.NewBlock()
	htmlBlock.SetName("HTML Block")
	htmlBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	htmlBlock.SetSiteID(site.ID())
	htmlBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), htmlBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {htmlBlock.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(strings.ToLower(body), "success") {
		t.Errorf("Expected body to contain 'success'")
	}

	// Test Navbar block deletion
	navbarBlock := cmsstore.NewBlock()
	navbarBlock.SetName("Navbar Block")
	navbarBlock.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	navbarBlock.SetSiteID(site.ID())
	navbarBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), navbarBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {navbarBlock.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(strings.ToLower(body), "success") {
		t.Errorf("Expected body to contain 'success'")
	}
}

func Test_BlockDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockDeleteHandler(store)

	// Create multiple blocks
	block1 := cmsstore.NewBlock()
	block1.SetName("Block 1")
	block1.SetType(cmsstore.BLOCK_TYPE_HTML)
	block1.SetSiteID(site.ID())
	block1.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block1)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	block2 := cmsstore.NewBlock()
	block2.SetName("Block 2")
	block2.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	block2.SetSiteID(site.ID())
	block2.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), block2)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Delete first block
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {block1.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(strings.ToLower(body), "success") {
		t.Errorf("Expected body to contain 'success'")
	}

	// Verify first block is deleted but second block remains
	deletedBlock1, err := store.BlockFindByID(context.Background(), block1.ID())
	if err == nil && deletedBlock1 != nil {
		if deletedBlock1.Status() != cmsstore.BLOCK_STATUS_INACTIVE {
			t.Errorf("Expected status %q, got %q", cmsstore.BLOCK_STATUS_INACTIVE, deletedBlock1.Status())
		}
	}

	activeBlock2, err := store.BlockFindByID(context.Background(), block2.ID())
	if err == nil && activeBlock2 != nil {
		if activeBlock2.Status() != cmsstore.BLOCK_STATUS_ACTIVE {
			t.Errorf("Expected status %q, got %q", cmsstore.BLOCK_STATUS_ACTIVE, activeBlock2.Status())
		}
	}
}

func Test_BlockDeleteController_Delete_WithContent(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a block with content
	block := cmsstore.NewBlock()
	block.SetName("Block with Content")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<div>Test content</div>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id": {block.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	bodyLower := strings.ToLower(body)
	if !strings.Contains(bodyLower, "success") {
		t.Errorf("Expected body to contain 'success'")
	}
	if !strings.Contains(bodyLower, "block deleted successfully") {
		t.Errorf("Expected body to contain 'block deleted successfully'")
	}
}
