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

func initMenuManagerHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuManagerController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuManagerController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Menu Manager") {
		t.Errorf("Expected body to contain 'Menu Manager'")
	}
	if !strings.Contains(body, "New Menu") {
		t.Errorf("Expected body to contain 'New Menu'")
	}
	if !strings.Contains(body, "table table-striped") {
		t.Errorf("Expected body to contain table class")
	}
	if !strings.Contains(body, "Filters") {
		t.Errorf("Expected body to contain 'Filters'")
	}
}

func Test_MenuManagerController_WithMenus(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

	// Seed a site and menus
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create menus directly
	menu1 := cmsstore.NewMenu()
	menu1.SetSiteID(site.ID())
	menu1.SetName("Test Menu 1")
	menu1.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu1)
	if err != nil {
		t.Fatalf("Failed to create menu1: %v", err)
	}

	menu2 := cmsstore.NewMenu()
	menu2.SetSiteID(site.ID())
	menu2.SetName("Test Menu 2")
	menu2.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu2)
	if err != nil {
		t.Fatalf("Failed to create menu2: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Should contain the menus
	if !strings.Contains(body, menu1.Name()) {
		t.Errorf("Expected body to contain menu1 name")
	}
	if !strings.Contains(body, menu2.Name()) {
		t.Errorf("Expected body to contain menu2 name")
	}
	if !strings.Contains(body, site.Name()) {
		t.Errorf("Expected body to contain site name")
	}
	if !strings.Contains(body, menu1.ID()) {
		t.Errorf("Expected body to contain menu1 ID")
	}
	if !strings.Contains(body, menu2.ID()) {
		t.Errorf("Expected body to contain menu2 ID")
	}
}

func Test_MenuManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

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
}

func Test_MenuManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a menu directly
	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu 1")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
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

			if !strings.Contains(body, "Menu Manager") {
				t.Errorf("Expected body to contain 'Menu Manager'")
			}
			if !strings.Contains(body, "table table-striped") {
				t.Errorf("Expected body to contain table class")
			}
		})
	}
}

func Test_MenuManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a menu directly
	menu1 := cmsstore.NewMenu()
	menu1.SetSiteID(site.ID())
	menu1.SetName("Test Menu 1")
	menu1.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu1)
	if err != nil {
		t.Fatalf("Failed to create menu1: %v", err)
	}

	// Create a second menu directly
	menu2 := cmsstore.NewMenu()
	menu2.SetSiteID(site.ID())
	menu2.SetName("Test Menu 2")
	menu2.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu2)
	if err != nil {
		t.Fatalf("Failed to create menu2: %v", err)
	}

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
			if err != nil {
				t.Fatalf("Failed to call endpoint: %v", err)
			}
			if response.StatusCode != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
			}

			if !strings.Contains(body, "Menu Manager") {
				t.Errorf("Expected body to contain 'Menu Manager'")
			}
			if tc.name == "Filter by menu name" {
				if !strings.Contains(body, tc.expected) {
					t.Errorf("Expected body to contain '%s'", tc.expected)
				}
			}
		})
	}
}

func Test_MenuManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Seed multiple menus for pagination
	for i := 0; i < 25; i++ {
		menu := cmsstore.NewMenu()
		menu.SetSiteID(site.ID())
		menu.SetName("Test Menu " + string(rune(i+65))) // Use A, B, C... instead of numbers
		menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
		err := store.MenuCreate(context.Background(), menu)
		if err != nil {
			t.Fatalf("Failed to create menu %d: %v", i, err)
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

	if !strings.Contains(body, "Menu Manager") {
		t.Errorf("Expected body to contain 'Menu Manager'")
	}
	if !strings.Contains(body, "pagination") {
		t.Errorf("Expected body to contain 'pagination'")
	}
}

func Test_MenuManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

	// Test with no menus
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Menu Manager") {
		t.Errorf("Expected body to contain 'Menu Manager'")
	}
	if !strings.Contains(body, "table table-striped") {
		t.Errorf("Expected body to contain table class")
	}
	if !strings.Contains(body, "Showing menus") {
		t.Errorf("Expected body to contain 'Showing menus'")
	}
}

func Test_MenuManagerController_TableActions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuManagerHandler(store)

	// Seed a site using the same pattern as site manager tests
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a menu directly
	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu 1")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
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
	if !strings.Contains(body, "menu-update") {
		t.Errorf("Expected body to contain update path")
	}
	if !strings.Contains(body, "menu-delete") {
		t.Errorf("Expected body to contain delete path")
	}
	if !strings.Contains(body, menu.ID()) {
		t.Errorf("Expected body to contain menu ID")
	}
}
