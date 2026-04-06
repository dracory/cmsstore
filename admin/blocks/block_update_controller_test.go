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
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockUpdateHandler(store)

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

func Test_BlockUpdateController_BlockIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockUpdateHandler(store)

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

func Test_BlockUpdateController_ViewSettings(t *testing.T) {
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
	block.SetContent("<p>Test content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
			"view":     {"settings"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Edit Block") {
		t.Errorf("Expected body to contain 'Edit Block'")
	}
	if !strings.Contains(body, "block_name") {
		t.Errorf("Expected body to contain 'block_name'")
	}
	if !strings.Contains(body, "block_type") {
		t.Errorf("Expected body to contain 'block_type'")
	}
	if !strings.Contains(body, "block_site_id") {
		t.Errorf("Expected body to contain 'block_site_id'")
	}
	if !strings.Contains(body, block.Name()) {
		t.Errorf("Expected body to contain block.Name()")
	}
	if !strings.Contains(body, block.Type()) {
		t.Errorf("Expected body to contain block.Type()")
	}
	if !strings.Contains(body, site.ID()) {
		t.Errorf("Expected body to contain site.ID()")
	}
}

func Test_BlockUpdateController_ViewContent(t *testing.T) {
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
	block.SetContent("<p>Test content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"block_id": {block.ID()},
			"view":     {"content"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Edit Block") {
		t.Errorf("Expected body to contain 'Edit Block'")
	}
	if !strings.Contains(body, "Content") {
		t.Errorf("Expected body to contain 'Content'")
	}
	if !strings.Contains(body, "block_content") {
		t.Errorf("Expected body to contain 'block_content'")
	}
}

func Test_BlockUpdateController_UpdateSettings(t *testing.T) {
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
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

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
	if !strings.Contains(bodyLower, "block updated successfully") {
		t.Errorf("Expected body to contain 'block updated successfully'")
	}

	// Verify block was updated
	updatedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	if err != nil {
		t.Fatalf("Failed to find block: %v", err)
	}
	if updatedBlock.Name() != "Updated Block" {
		t.Errorf("Expected name %q, got %q", "Updated Block", updatedBlock.Name())
	}
	if updatedBlock.Type() != cmsstore.BLOCK_TYPE_NAVBAR {
		t.Errorf("Expected type %q, got %q", cmsstore.BLOCK_TYPE_NAVBAR, updatedBlock.Type())
	}
	if updatedBlock.Status() != cmsstore.BLOCK_STATUS_ACTIVE {
		t.Errorf("Expected status %q, got %q", cmsstore.BLOCK_STATUS_ACTIVE, updatedBlock.Status())
	}
}

func Test_BlockUpdateController_UpdateContent(t *testing.T) {
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

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":      {block.ID()},
			"block_content": {"<div>Updated content</div>"},
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

	// Verify block content was updated
	updatedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	if err != nil {
		t.Fatalf("Failed to find block: %v", err)
	}
	if updatedBlock.Content() != "<div>Updated content</div>" {
		t.Errorf("Expected content %q, got %q", "<div>Updated content</div>", updatedBlock.Content())
	}
}

func Test_BlockUpdateController_Update_ValidationError_MissingName(t *testing.T) {
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

	handler := initBlockUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":   {block.ID()},
			"block_type": {cmsstore.BLOCK_TYPE_HTML},
			"site_id":    {site.ID()},
			"view":       {"settings"},
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
}

func Test_BlockUpdateController_Update_WithMemo(t *testing.T) {
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
}

func Test_BlockUpdateController_Update_WithHandle(t *testing.T) {
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

func Test_BlockUpdateController_BlockTypeChange(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and create a draft block
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	block.SetContent("<p>Original HTML content</p>")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

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

	// Verify block type was changed
	updatedBlock, err := store.BlockFindByID(context.Background(), block.ID())
	if err != nil {
		t.Fatalf("Failed to find block: %v", err)
	}
	if updatedBlock.Type() != cmsstore.BLOCK_TYPE_NAVBAR {
		t.Errorf("Expected type %q, got %q", cmsstore.BLOCK_TYPE_NAVBAR, updatedBlock.Type())
	}
}

func Test_BlockUpdateController_DifferentBlockTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockUpdateHandler(store)

	// Test HTML block update
	htmlBlock := cmsstore.NewBlock()
	htmlBlock.SetName("HTML Block")
	htmlBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	htmlBlock.SetSiteID(site.ID())
	htmlBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), htmlBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":      {htmlBlock.ID()},
			"block_name":    {"Updated HTML Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"site_id":       {site.ID()},
			"block_content": {"<div>Updated HTML</div>"},
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

	// Test Navbar block update
	navbarBlock := cmsstore.NewBlock()
	navbarBlock.SetName("Navbar Block")
	navbarBlock.SetType(cmsstore.BLOCK_TYPE_NAVBAR)
	navbarBlock.SetSiteID(site.ID())
	navbarBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	err = store.BlockCreate(context.Background(), navbarBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_id":   {navbarBlock.ID()},
			"block_name": {"Updated Navbar Block"},
			"block_type": {cmsstore.BLOCK_TYPE_NAVBAR},
			"site_id":    {site.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	bodyLower = strings.ToLower(body)
	if !strings.Contains(bodyLower, "success") {
		t.Errorf("Expected body to contain 'success'")
	}
}
