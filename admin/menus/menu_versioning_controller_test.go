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

func initMenuVersioningHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuVersioningController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuVersioningController_MenuIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuVersioningHandler(store)

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

func Test_MenuVersioningController_MenuIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuVersioningHandler(store)

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

func Test_MenuVersioningController_ListRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuVersioningHandler(store)

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

	if !strings.Contains(body, "Menu Revisions") {
		t.Errorf("Expected body to contain 'Menu Revisions'")
	}
}

func Test_MenuVersioningController_PreviewRevision(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuVersioningHandler(store)

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

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(menu.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_MENU)
	versioning.SetContent(`{"name": "Test Menu", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id":       {menu.ID()},
			"versioning_id": {versioning.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Menu Revision") {
		t.Errorf("Expected body to contain 'Menu Revision'")
	}
	if !strings.Contains(body, "Test Menu") {
		t.Errorf("Expected body to contain 'Test Menu'")
	}
}

func Test_MenuVersioningController_RestoreAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuVersioningHandler(store)

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

	// Create a versioning entry with different data
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(menu.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_MENU)
	versioning.SetContent(`{"name": "Restored Menu", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id":       {menu.ID()},
			"versioning_id": {versioning.ID()},
		},
		PostValues: url.Values{
			"action":              {"restore"},
			"revision_attributes": {"name"},
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
	if !strings.Contains(bodyLower, "restored successfully") {
		t.Errorf("Expected body to contain success message")
	}

	// Verify the menu was updated
	updatedMenu, err := store.MenuFindByID(context.Background(), menu.ID())
	if err != nil {
		t.Fatalf("Failed to find menu: %v", err)
	}
	if updatedMenu == nil {
		t.Fatalf("Expected menu to not be nil")
	}
	if updatedMenu.Name() != "Restored Menu" {
		t.Errorf("Expected menu name 'Restored Menu', got '%s'", updatedMenu.Name())
	}
}

func Test_MenuVersioningController_RestoreNoAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuVersioningHandler(store)

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

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(menu.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_MENU)
	versioning.SetContent(`{"name": "Test Menu", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id":       {menu.ID()},
			"versioning_id": {versioning.ID()},
		},
		PostValues: url.Values{
			"action":              {"restore"},
			"revision_attributes": {"name"},
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
	if !strings.Contains(bodyLower, "restored successfully") {
		t.Errorf("Expected body to contain success message")
	}
}
