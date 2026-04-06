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
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
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
	if !strings.Contains(bodyLower, "block id is required") {
		t.Errorf("Expected body to contain 'block id is required'")
	}
}

func Test_BlockVersioningController_BlockIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {"invalid-id"},
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

func Test_BlockVersioningController_ListRevisions(t *testing.T) {
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
	block.SetContent("<p>Original content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create some versioning entries
	versioning1 := cmsstore.NewVersioning()
	versioning1.SetEntityID(block.ID())
	versioning1.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning1.SetContent(`{"content": "<p>Version 1 content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning1)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	versioning2 := cmsstore.NewVersioning()
	versioning2.SetEntityID(block.ID())
	versioning2.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning2.SetContent(`{"content": "<p>Version 2 content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning2)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initBlockVersioningHandler(store)

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

	if !strings.Contains(body, "Block Revisions") {
		t.Errorf("Expected body to contain 'Block Revisions'")
	}
	if !strings.Contains(body, "Version") {
		t.Errorf("Expected body to contain 'Version'")
	}
	if !strings.Contains(body, "Created") {
		t.Errorf("Expected body to contain 'Created'")
	}
	if !strings.Contains(body, "Actions") {
		t.Errorf("Expected body to contain 'Actions'")
	}
}

func Test_BlockVersioningController_PreviewRevision(t *testing.T) {
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
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Historical content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id":      {block.ID()},
			"versioning_id": {versioning.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Block Revision") {
		t.Errorf("Expected body to contain 'Block Revision'")
	}
	if !strings.Contains(body, "Historical content") {
		t.Errorf("Expected body to contain 'Historical content'")
	}
}

func Test_BlockVersioningController_RestoreAttributes(t *testing.T) {
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
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create a versioning entry with different content
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Restored content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":            {block.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {"content"},
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
	if !strings.Contains(bodyLower, "revision attributes restored successfully") {
		t.Errorf("Expected body to contain 'revision attributes restored successfully'")
	}

	// Verify block content was restored
	restoredBlock, err := store.BlockFindByID(context.Background(), block.ID())
	if err != nil {
		t.Fatalf("Failed to find block: %v", err)
	}
	if restoredBlock.Content() != "<p>Restored content</p>" {
		t.Errorf("Expected content %q, got %q", "<p>Restored content</p>", restoredBlock.Content())
	}
}

func Test_BlockVersioningController_RestoreNoAttributes(t *testing.T) {
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
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Restored content</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":      {block.ID()},
			"versioning_id": {versioning.ID()},
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
	if !strings.Contains(bodyLower, "no revision attributes were selected") {
		t.Errorf("Expected body to contain 'no revision attributes were selected'")
	}
}

func Test_BlockVersioningController_RestoreMultipleAttributes(t *testing.T) {
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
	block.SetName("Original Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	block.SetContent("<p>Original content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create a versioning entry with different values
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(block.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning.SetContent(`{"content": "<p>Restored content</p>", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":            {block.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {"content", "status"},
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
	if !strings.Contains(bodyLower, "revision attributes restored successfully") {
		t.Errorf("Expected body to contain 'revision attributes restored successfully'")
	}
}

func Test_BlockVersioningController_VersioningNotFound(t *testing.T) {
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
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id":      {block.ID()},
			"versioning_id": {"non-existent-versioning-id"},
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
	if !strings.Contains(bodyLower, "versioning not found") {
		t.Errorf("Expected body to contain 'versioning not found'")
	}
}

func Test_BlockVersioningController_EmptyRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and block (no versioning entries)
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetContent("<p>Current content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockVersioningHandler(store)

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

	if !strings.Contains(body, "Version") {
		t.Errorf("Expected body to contain 'Version'")
	}
	if !strings.Contains(body, "Created") {
		t.Errorf("Expected body to contain 'Created'")
	}
	if !strings.Contains(body, "Actions") {
		t.Errorf("Expected body to contain 'Actions'")
	}
}

func Test_BlockVersioningController_DifferentBlockTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockVersioningHandler(store)

	// Test HTML block versioning
	htmlBlock := cmsstore.NewBlock()
	htmlBlock.SetName("HTML Block")
	htmlBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	htmlBlock.SetSiteID(site.ID())
	htmlBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	htmlBlock.SetContent("<p>HTML content</p>")
	err = store.BlockCreate(context.Background(), htmlBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create versioning for HTML block
	versioning1 := cmsstore.NewVersioning()
	versioning1.SetEntityID(htmlBlock.ID())
	versioning1.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning1.SetContent(`{"content": "<p>HTML version 1</p>"}`)
	err = store.VersioningCreate(context.Background(), versioning1)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {htmlBlock.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Revisions") {
		t.Errorf("Expected body to contain 'Revisions'")
	}
	if !strings.Contains(body, "Preview") {
		t.Errorf("Expected body to contain 'Preview'")
	}

	// Test Navbar block versioning
	navbarBlock := cmsstore.NewBlock()
	navbarBlock.SetName("Navbar Block")
	navbarBlock.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	navbarBlock.SetSiteID(site.ID())
	navbarBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), navbarBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create versioning for Navbar block
	versioning2 := cmsstore.NewVersioning()
	versioning2.SetEntityID(navbarBlock.ID())
	versioning2.SetEntityType(cmsstore.VERSIONING_TYPE_BLOCK)
	versioning2.SetContent(`{"content": "navbar data"}`)
	err = store.VersioningCreate(context.Background(), versioning2)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {navbarBlock.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Revisions") {
		t.Errorf("Expected body to contain 'Revisions'")
	}
	if !strings.Contains(body, "Preview") {
		t.Errorf("Expected body to contain 'Preview'")
	}

	// Test preview modal shows content for navbar block
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id":      {navbarBlock.ID()},
			"versioning_id": {versioning2.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "navbar data") {
		t.Errorf("Expected body to contain 'navbar data'")
	}
}
