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
	_ "modernc.org/sqlite"
)

func initPageDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewPageDeleteController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_PageDeleteController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageDeleteHandler(store)

	// Seed a site and page
	site, err := testutils.SeedSite(store, testutils.SITE_01)
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	page, err := testutils.SeedPage(store, testutils.PAGE_01, site.ID())
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {page.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Delete Page") {
		t.Errorf("Expected body to contain 'Delete Page'")
	}
	if !strings.Contains(body, "Are you sure you want to delete this page?") {
		t.Errorf("Expected body to contain confirmation message")
	}
	if !strings.Contains(body, "This action cannot be undone") {
		t.Errorf("Expected body to contain warning message")
	}
	if !strings.Contains(body, "name=\"page_id\"") {
		t.Errorf("Expected body to contain page_id input")
	}
	if !strings.Contains(body, "value=\""+page.ID()+"\"") {
		t.Errorf("Expected body to contain page ID value")
	}
}

func Test_PageDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageDeleteHandler(store)

	// Seed a site and page
	site, err := testutils.SeedSite(store, testutils.SITE_01)
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	page, err := testutils.SeedPage(store, testutils.PAGE_01, site.ID())
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	// Verify page exists before deletion
	pages, _ := store.PageList(context.Background(), cmsstore.PageQuery().SetID(page.ID()))
	if len(pages) != 1 {
		t.Errorf("Expected 1 page before deletion, got %d", len(pages))
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"page_id": {page.ID()},
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
	if !strings.Contains(bodyLower, "page deleted successfully") {
		t.Errorf("Expected body to contain success message")
	}

	// Verify page is soft deleted
	pages, _ = store.PageList(context.Background(), cmsstore.PageQuery().SetID(page.ID()))
	if len(pages) != 0 {
		t.Errorf("Expected 0 pages after deletion, got %d (should be soft deleted)", len(pages))
	}
}

func Test_PageDeleteController_Delete_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageDeleteHandler(store)

	testCases := []struct {
		name       string
		postValues url.Values
		errorMsg   string
	}{
		{
			name:       "Missing page ID",
			postValues: url.Values{},
			errorMsg:   "page id is required",
		},
		{
			name: "Empty page ID",
			postValues: url.Values{
				"page_id": {""},
			},
			errorMsg: "page id is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
				PostValues: tc.postValues,
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
			if !strings.Contains(bodyLower, tc.errorMsg) {
				t.Errorf("Expected body to contain '%s'", tc.errorMsg)
			}
		})
	}
}

func Test_PageDeleteController_PageNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageDeleteHandler(store)

	// Test with non-existent page ID
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {"non-existent-page-id"},
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
	if !strings.Contains(bodyLower, "page not found") {
		t.Errorf("Expected body to contain 'page not found'")
	}

	// Test POST with non-existent page ID
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"page_id": {"non-existent-page-id"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	bodyLower = strings.ToLower(body)
	if !strings.Contains(bodyLower, "error") {
		t.Errorf("Expected body to contain 'error'")
	}
	if !strings.Contains(bodyLower, "page not found") {
		t.Errorf("Expected body to contain 'page not found'")
	}
}

func Test_PageDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageDeleteHandler(store)

	// Seed a site and page
	site, err := testutils.SeedSite(store, testutils.SITE_01)
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	page, err := testutils.SeedPage(store, testutils.PAGE_01, site.ID())
	if err != nil {
		t.Fatalf("Failed to seed page: %v", err)
	}

	// First, GET the delete modal to confirm it shows correctly
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {page.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, page.ID()) {
		t.Errorf("Expected body to contain page ID")
	}

	// Then, POST to delete the page
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"page_id": {page.ID()},
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
