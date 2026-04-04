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

func initMenuDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuDeleteController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuDeleteController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuDeleteHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Delete Menu") {
		t.Errorf("Expected body to contain 'Delete Menu'")
	}
	if !strings.Contains(body, "Are you sure") {
		t.Errorf("Expected body to contain confirmation message")
	}
	if !strings.Contains(body, menu.ID()) {
		t.Errorf("Expected body to contain menu ID")
	}
	if !strings.Contains(body, "This action cannot be undone") {
		t.Errorf("Expected body to contain 'This action cannot be undone'")
	}
	if !strings.Contains(body, "name=\"menu_id\"") {
		t.Errorf("Expected body to contain 'name=\"menu_id\"'")
	}
	if !strings.Contains(body, "value=\""+menu.ID()+"\"") {
		t.Errorf("Expected body to contain menu ID value")
	}
}

func Test_MenuDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuDeleteHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	// Verify menu exists before deletion
	menus, _ := store.MenuList(context.Background(), cmsstore.MenuQuery().SetID(menu.ID()))
	if len(menus) != 1 {
		t.Errorf("Expected 1 menu before deletion, got %d", len(menus))
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"menu_id": {menu.ID()},
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
	if !strings.Contains(bodyLower, "menu deleted successfully") {
		t.Errorf("Expected body to contain success message")
	}

	// Verify menu is soft deleted
	menus, _ = store.MenuList(context.Background(), cmsstore.MenuQuery().SetID(menu.ID()))
	if len(menus) != 0 {
		t.Errorf("Expected 0 menus after deletion, got %d (should be soft deleted)", len(menus))
	}
}

func Test_MenuDeleteController_Delete_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuDeleteHandler(store)

	testCases := []struct {
		name       string
		postValues url.Values
		errorMsg   string
	}{
		{
			name:       "Missing menu ID",
			postValues: url.Values{},
			errorMsg:   "menu id is required",
		},
		{
			name: "Empty menu ID",
			postValues: url.Values{
				"menu_id": {""},
			},
			errorMsg: "menu id is required",
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

func Test_MenuDeleteController_MenuNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuDeleteHandler(store)

	// Test with non-existent menu ID
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {"non-existent-menu-id"},
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
	if !strings.Contains(bodyLower, "menu not found") {
		t.Errorf("Expected body to contain 'menu not found'")
	}

	// Test POST with non-existent menu ID
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"menu_id": {"non-existent-menu-id"},
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
	if !strings.Contains(bodyLower, "menu not found") {
		t.Errorf("Expected body to contain 'menu not found'")
	}
}

func Test_MenuDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuDeleteHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	// First, GET the delete modal to confirm it shows correctly
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, menu.ID()) {
		t.Errorf("Expected body to contain menu ID")
	}

	// Then, POST to delete the menu
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"menu_id": {menu.ID()},
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
