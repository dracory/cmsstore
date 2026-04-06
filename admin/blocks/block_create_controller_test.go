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

func initBlockCreateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewBlockCreateController(ui).Handler
}

func Test_BlockCreateController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "New Block") {
		t.Errorf("Expected body to contain 'New Block'")
	}
	if !strings.Contains(body, "block_name") {
		t.Errorf("Expected body to contain block_name")
	}
	if !strings.Contains(body, "block_type") {
		t.Errorf("Expected body to contain block_type")
	}
	if !strings.Contains(body, "site_id") {
		t.Errorf("Expected body to contain site_id")
	}
}

func Test_BlockCreateController_Create(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_name":    {"Test Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"site_id":       {site.ID()},
			"block_content": {"<p>Test content</p>"},
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
	if !strings.Contains(bodyLower, "block created successfully") {
		t.Errorf("Expected body to contain 'block created successfully'")
	}
}

func Test_BlockCreateController_Create_ValidationError_MissingName(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"site_id":       {site.ID()},
			"block_content": {"<p>Test content</p>"},
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

func Test_BlockCreateController_Create_MissingTypeUsesDefault(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_name":    {"Test Block"},
			"site_id":       {site.ID()},
			"block_content": {"<p>Test content</p>"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Block type has a default value, so missing type succeeds
	if !strings.Contains(strings.ToLower(body), "success") {
		t.Errorf("Expected body to contain 'success'")
	}
}

func Test_BlockCreateController_Create_ValidationError_MissingSiteID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_name":    {"Test Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"block_content": {"<p>Test content</p>"},
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

func Test_BlockCreateController_Create_WithSiteList(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, "Test Site 1")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	site2 := cmsstore.NewSite()
	site2.SetName("Test Site 2")
	site2.SetDomainNames([]string{"site2.example.com"})
	err = store.SiteCreate(context.Background(), site2)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Test Site 1") {
		t.Errorf("Expected body to contain 'Test Site 1'")
	}
	if !strings.Contains(body, "Test Site 2") {
		t.Errorf("Expected body to contain 'Test Site 2'")
	}
	if !strings.Contains(body, site1.ID()) {
		t.Errorf("Expected body to contain %q", site1.ID())
	}
	if !strings.Contains(body, site2.ID()) {
		t.Errorf("Expected body to contain %q", site2.ID())
	}
}

func Test_BlockCreateController_Create_WithDifferentBlockTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockCreateHandler(store)

	// Test HTML block creation
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_name":    {"HTML Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"site_id":       {site.ID()},
			"block_content": {"<div>HTML content</div>"},
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

	// Test Navbar block creation
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_name": {"Navbar Block"},
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
	// Navbar blocks may require additional fields, so we'll just check for no error
	if strings.Contains(strings.ToLower(body), "error") {
		t.Errorf("Expected body to NOT contain 'error'")
	}
}

func Test_BlockCreateController_Create_WithMemo(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_name":    {"Test Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"site_id":       {site.ID()},
			"block_content": {"<p>Test content</p>"},
			"block_memo":    {"This is a test memo"},
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
	if !strings.Contains(bodyLower, "block created successfully") {
		t.Errorf("Expected body to contain 'block created successfully'")
	}
}

func Test_BlockCreateController_Create_WithHandle(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initBlockCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"block_name":    {"Test Block"},
			"block_type":    {cmsstore.BLOCK_TYPE_HTML},
			"site_id":       {site.ID()},
			"block_content": {"<p>Test content</p>"},
			"block_handle":  {"test-block-handle"},
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
