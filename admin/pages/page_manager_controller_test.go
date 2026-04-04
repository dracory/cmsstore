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

func initPageManagerHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewPageManagerController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_PageManagerController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Page Manager") {
		t.Errorf("Expected body to contain 'Page Manager'")
	}
	if !strings.Contains(body, "New Page") {
		t.Errorf("Expected body to contain 'New Page'")
	}
	if !strings.Contains(body, "table table-striped") {
		t.Errorf("Expected body to contain table class")
	}
	if !strings.Contains(body, "Filters") {
		t.Errorf("Expected body to contain 'Filters'")
	}
}

func Test_PageManagerController_WithPages(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a page directly without using SeedPage to have more control
	page := cmsstore.NewPage()
	page.SetSiteID(site.ID())
	page.SetName("Test Page 1")
	page.SetAlias("test-page-1")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Should contain the page content
	if !strings.Contains(body, "Page Manager") {
		t.Errorf("Expected body to contain 'Page Manager'")
	}
	if !strings.Contains(body, "Test Page 1") {
		t.Errorf("Expected body to contain 'Test Page 1'")
	}
	if !strings.Contains(body, "Test Site") {
		t.Errorf("Expected body to contain 'Test Site'")
	}
}

func Test_PageManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"action": {ActionModalPageFilterShow},
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
	if !strings.Contains(body, "name=\"filter_status\"") {
		t.Errorf("Expected body to contain filter_status input")
	}
	if !strings.Contains(body, "name=\"filter_name\"") {
		t.Errorf("Expected body to contain filter_name input")
	}
	if !strings.Contains(body, "name=\"filter_site_id\"") {
		t.Errorf("Expected body to contain filter_site_id input")
	}
	if !strings.Contains(body, "name=\"filter_created_from\"") {
		t.Errorf("Expected body to contain filter_created_from input")
	}
	if !strings.Contains(body, "name=\"filter_created_to\"") {
		t.Errorf("Expected body to contain filter_created_to input")
	}
}

func Test_PageManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a page directly
	page := cmsstore.NewPage()
	page.SetSiteID(site.ID())
	page.SetName("Test Page 1")
	page.SetAlias("test-page-1")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	testCases := []struct {
		name      string
		sortBy    string
		sortOrder string
	}{
		{
			name:      "Sort by name ASC",
			sortBy:    cmsstore.COLUMN_NAME,
			sortOrder: "asc",
		},
		{
			name:      "Sort by name DESC",
			sortBy:    cmsstore.COLUMN_NAME,
			sortOrder: "desc",
		},
		{
			name:      "Sort by created_at ASC",
			sortBy:    cmsstore.COLUMN_CREATED_AT,
			sortOrder: "asc",
		},
		{
			name:      "Sort by created_at DESC",
			sortBy:    cmsstore.COLUMN_CREATED_AT,
			sortOrder: "desc",
		},
		{
			name:      "Sort by status ASC",
			sortBy:    cmsstore.COLUMN_STATUS,
			sortOrder: "asc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
				GetValues: url.Values{
					"by":   {tc.sortBy},
					"sort": {tc.sortOrder},
				},
			})
			if err != nil {
				t.Fatalf("Failed to call endpoint: %v", err)
			}
			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			if !strings.Contains(body, "Page Manager") {
				t.Errorf("Expected body to contain 'Page Manager'")
			}
			if !strings.Contains(body, "table table-striped") {
				t.Errorf("Expected body to contain table class")
			}
		})
	}
}

func Test_PageManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a page directly
	page1 := cmsstore.NewPage()
	page1.SetSiteID(site.ID())
	page1.SetName("Test Page 1")
	page1.SetAlias("test-page-1")
	page1.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page1)
	if err != nil {
		t.Fatalf("Failed to create page1: %v", err)
	}

	// Create a second page directly with alias
	page2 := cmsstore.NewPage()
	page2.SetSiteID(site.ID())
	page2.SetName("Test Page 2")
	page2.SetAlias("test-page-2")
	page2.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page2)
	if err != nil {
		t.Fatalf("Failed to create page2: %v", err)
	}

	testCases := []struct {
		name       string
		filterName string
		expected   string
	}{
		{
			name:       "Filter by page name",
			filterName: page1.Name(),
			expected:   page1.Name(),
		},
		{
			name:       "Filter by status",
			filterName: cmsstore.PAGE_STATUS_ACTIVE,
			expected:   "status: " + cmsstore.PAGE_STATUS_ACTIVE,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
				GetValues: url.Values{
					"filter_name": {tc.filterName},
				},
			})
			if err != nil {
				t.Fatalf("Failed to call endpoint: %v", err)
			}
			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			if !strings.Contains(body, "Page Manager") {
				t.Errorf("Expected body to contain 'Page Manager'")
			}
			if tc.name == "Filter by page name" {
				if !strings.Contains(body, tc.expected) {
					t.Errorf("Expected body to contain '%s'", tc.expected)
				}
			}
		})
	}
}

func Test_PageManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Seed multiple pages for pagination
	for i := 0; i < 25; i++ {
		page := cmsstore.NewPage()
		page.SetSiteID(site.ID())
		page.SetName("Test Page " + string(rune(i+65))) // Use A, B, C... instead of numbers
		page.SetAlias("test-page-" + string(rune(i+65)))
		page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
		err := store.PageCreate(context.Background(), page)
		if err != nil {
			t.Fatalf("Failed to create page %d: %v", i, err)
		}
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page":     {"1"},
			"per_page": {"10"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Page Manager") {
		t.Errorf("Expected body to contain 'Page Manager'")
	}
	if !strings.Contains(body, "pagination") {
		t.Errorf("Expected body to contain 'pagination'")
	}
}

func Test_PageManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	// Test with no pages
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Page Manager") {
		t.Errorf("Expected body to contain 'Page Manager'")
	}
	if !strings.Contains(body, "table table-striped") {
		t.Errorf("Expected body to contain table class")
	}
	if !strings.Contains(body, "Showing pages") {
		t.Errorf("Expected body to contain 'Showing pages'")
	}
}

func Test_PageManagerController_TableActions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a page directly
	page := cmsstore.NewPage()
	page.SetSiteID(site.ID())
	page.SetName("Test Page 1")
	page.SetAlias("test-page-1")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Check for action buttons
	if !strings.Contains(body, "btn btn-primary") {
		t.Errorf("Expected body to contain edit button")
	}
	if !strings.Contains(body, "btn btn-danger") {
		t.Errorf("Expected body to contain delete button")
	}
	if !strings.Contains(body, "page-update") {
		t.Errorf("Expected body to contain update path")
	}
	if !strings.Contains(body, "page-delete") {
		t.Errorf("Expected body to contain delete path")
	}
	if !strings.Contains(body, page.ID()) {
		t.Errorf("Expected body to contain page ID")
	}
}
