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

func initMenuManagerHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuManagerController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuManagerController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Menu Manager")
	assert.Contains(t, body, "New Menu")
	assert.Contains(t, body, "table table-striped")
	assert.Contains(t, body, "Filters")
}

func Test_MenuManagerController_WithMenus(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

	// Seed a site and menus
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create menus directly
	menu1 := cmsstore.NewMenu()
	menu1.SetSiteID(site.ID())
	menu1.SetName("Test Menu 1")
	menu1.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu1)
	require.NoError(t, err)

	menu2 := cmsstore.NewMenu()
	menu2.SetSiteID(site.ID())
	menu2.SetName("Test Menu 2")
	menu2.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu2)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Should contain the menus
	assert.Contains(t, body, menu1.Name())
	assert.Contains(t, body, menu2.Name())
	assert.Contains(t, body, site.Name())
	assert.Contains(t, body, menu1.ID())
	assert.Contains(t, body, menu2.ID())
}

func Test_MenuManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

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
}

func Test_MenuManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a menu directly
	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu 1")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
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

			assert.Contains(t, body, "Menu Manager")
			assert.Contains(t, body, "table table-striped")
		})
	}
}

func Test_MenuManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a menu directly
	menu1 := cmsstore.NewMenu()
	menu1.SetSiteID(site.ID())
	menu1.SetName("Test Menu 1")
	menu1.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu1)
	require.NoError(t, err)

	// Create a second menu directly
	menu2 := cmsstore.NewMenu()
	menu2.SetSiteID(site.ID())
	menu2.SetName("Test Menu 2")
	menu2.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu2)
	require.NoError(t, err)

	testCases := []struct {
		name       string
		filterName string
		expected   string
	}{
		{
			name:       "Filter by menu name",
			filterName: menu1.Name(),
			expected:   menu1.Name(),
		},
		{
			name:       "Filter by status",
			filterName: cmsstore.MENU_STATUS_ACTIVE,
			expected:   "status: " + cmsstore.MENU_STATUS_ACTIVE,
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

			assert.Contains(t, body, "Menu Manager")
			if tc.name == "Filter by menu name" {
				assert.Contains(t, body, tc.expected)
			}
		})
	}
}

func Test_MenuManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Seed multiple menus for pagination
	for i := 0; i < 25; i++ {
		menu := cmsstore.NewMenu()
		menu.SetSiteID(site.ID())
		menu.SetName("Test Menu " + string(rune(i+65))) // Use A, B, C... instead of numbers
		menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
		err := store.MenuCreate(context.Background(), menu)
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

	assert.Contains(t, body, "Menu Manager")
	assert.Contains(t, body, "pagination")
}

func Test_MenuManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

	// Test with no menus
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Menu Manager")
	assert.Contains(t, body, "table table-striped")
	assert.Contains(t, body, "Showing menus")
}

func Test_MenuManagerController_TableActions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a menu directly
	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu 1")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Check for action buttons
	assert.Contains(t, body, "btn btn-primary") // Edit button
	assert.Contains(t, body, "btn btn-danger")  // Delete button
	assert.Contains(t, body, "menu-update")      // Update path
	assert.Contains(t, body, "menu-delete")      // Delete path
	assert.Contains(t, body, menu.ID())
}
