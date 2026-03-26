package admin

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Page Manager")
	assert.Contains(t, body, "New Page")
	assert.Contains(t, body, "table table-striped")
	assert.Contains(t, body, "Filters")
}

func Test_PageManagerController_WithPages(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a page directly without using SeedPage to have more control
	page := cmsstore.NewPage()
	page.SetSiteID(site.ID())
	page.SetName("Test Page 1")
	page.SetAlias("test-page-1")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Should contain the page content
	assert.Contains(t, body, "Page Manager")
	assert.Contains(t, body, "Test Page 1")
	assert.Contains(t, body, "Test Site")
}

func Test_PageManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"action": {ActionModalPageFilterShow},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Filters")
	assert.Contains(t, body, "name=\"filter_status\"")
	assert.Contains(t, body, "name=\"filter_name\"")
	assert.Contains(t, body, "name=\"filter_site_id\"")
	assert.Contains(t, body, "name=\"filter_created_from\"")
	assert.Contains(t, body, "name=\"filter_created_to\"")
}

func Test_PageManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a page directly
	page := cmsstore.NewPage()
	page.SetSiteID(site.ID())
	page.SetName("Test Page 1")
	page.SetAlias("test-page-1")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)

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
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)

			assert.Contains(t, body, "Page Manager")
			assert.Contains(t, body, "table table-striped")
		})
	}
}

func Test_PageManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a page directly
	page1 := cmsstore.NewPage()
	page1.SetSiteID(site.ID())
	page1.SetName("Test Page 1")
	page1.SetAlias("test-page-1")
	page1.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page1)
	require.NoError(t, err)

	// Create a second page directly with alias
	page2 := cmsstore.NewPage()
	page2.SetSiteID(site.ID())
	page2.SetName("Test Page 2")
	page2.SetAlias("test-page-2")
	page2.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page2)
	require.NoError(t, err)

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
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)

			assert.Contains(t, body, "Page Manager")
			if tc.name == "Filter by page name" {
				assert.Contains(t, body, tc.expected)
			}
		})
	}
}

func Test_PageManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Seed multiple pages for pagination
	for i := 0; i < 25; i++ {
		page := cmsstore.NewPage()
		page.SetSiteID(site.ID())
		page.SetName("Test Page " + string(rune(i+65))) // Use A, B, C... instead of numbers
		page.SetAlias("test-page-" + string(rune(i+65)))
		page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
		err := store.PageCreate(context.Background(), page)
		require.NoError(t, err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page":     {"1"},
			"per_page": {"10"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Page Manager")
	assert.Contains(t, body, "pagination")
}

func Test_PageManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	// Test with no pages
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Page Manager")
	assert.Contains(t, body, "table table-striped")
	assert.Contains(t, body, "Showing pages")
}

func Test_PageManagerController_TableActions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a page directly
	page := cmsstore.NewPage()
	page.SetSiteID(site.ID())
	page.SetName("Test Page 1")
	page.SetAlias("test-page-1")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Check for action buttons
	assert.Contains(t, body, "btn btn-primary") // Edit button
	assert.Contains(t, body, "btn btn-danger")  // Delete button
	assert.Contains(t, body, "page-update")     // Update path
	assert.Contains(t, body, "page-delete")     // Delete path
	assert.Contains(t, body, page.ID())
}
