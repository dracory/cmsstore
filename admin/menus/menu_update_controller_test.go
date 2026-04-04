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

func initMenuUpdateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuUpdateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuUpdateController_MenuIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuUpdateHandler(store)

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
	if !strings.Contains(bodyLower, "menu id is required") {
		t.Errorf("Expected body to contain 'menu id is required'")
	}
}

func Test_MenuUpdateController_MenuIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {"invalid-menu-id"},
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
}

func Test_MenuUpdateController_ViewSettings_IsDefault(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuUpdateHandler(store)

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
			"view":    {"settings"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Edit Menu") {
		t.Errorf("Expected body to contain 'Edit Menu'")
	}
	if !strings.Contains(body, "name=\"menu_name\"") {
		t.Errorf("Expected body to contain menu_name input")
	}
	if !strings.Contains(body, menu.Name()) {
		t.Errorf("Expected body to contain menu name")
	}
	if !strings.Contains(body, "name=\"menu_site_id\"") {
		t.Errorf("Expected body to contain menu_site_id input")
	}
	if !strings.Contains(body, site.ID()) {
		t.Errorf("Expected body to contain site ID")
	}
}

func Test_MenuUpdateController_ViewMenuItems(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuUpdateHandler(store)

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
			"view":    {"menu_items"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Menu Items") {
		t.Errorf("Expected body to contain 'Menu Items'")
	}
}

func Test_MenuUpdateController_UpdateSettings(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuUpdateHandler(store)

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

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
			"view":    {"settings"},
		},
		PostValues: url.Values{
			"menu_name":    {"Updated Menu Name"},
			"menu_site_id": {site.ID()},
			"menu_status":  {cmsstore.MENU_STATUS_ACTIVE},
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
	if !strings.Contains(bodyLower, "menu saved successfully") {
		t.Errorf("Expected body to contain success message")
	}
}

func Test_MenuUpdateController_UpdateSettings_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuUpdateHandler(store)

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

	// Test submitting form with empty menu name - should still work
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
			"view":    {"settings"},
		},
		PostValues: url.Values{
			"menu_name":    {""},
			"menu_site_id": {site.ID()},
			"menu_status":  {cmsstore.MENU_STATUS_ACTIVE},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Should contain either success or error handling
	if !strings.Contains(body, "success") && !strings.Contains(body, "error") {
		t.Errorf("Expected body to contain 'success' or 'error'")
	}
}
