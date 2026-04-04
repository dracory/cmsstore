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

func initMenuCreateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuCreateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuCreateController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "New Menu") {
		t.Errorf("Expected body to contain 'New Menu'")
	}
	if !strings.Contains(body, "name=\"menu_name\"") {
		t.Errorf("Expected body to contain menu_name input")
	}
	if !strings.Contains(body, "name=\"site_id\"") {
		t.Errorf("Expected body to contain site_id input")
	}
}

func Test_MenuCreateController_Create(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuCreateHandler(store)

	// First seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"menu_name": {"Test New Menu"},
			"site_id":   {site.ID()},
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
	if !strings.Contains(bodyLower, "menu created successfully") {
		t.Errorf("Expected body to contain 'menu created successfully'")
	}

	// Verify in DB
	menus, _ := store.MenuList(context.Background(), cmsstore.MenuQuery().SetNameLike("Test New Menu"))
	if len(menus) != 1 {
		t.Errorf("Expected 1 menu, got %d", len(menus))
	}
	if menus[0].Name() != "Test New Menu" {
		t.Errorf("Expected menu name 'Test New Menu', got '%s'", menus[0].Name())
	}
	if menus[0].SiteID() != site.ID() {
		t.Errorf("Expected site ID '%s', got '%s'", site.ID(), menus[0].SiteID())
	}
}

func Test_MenuCreateController_Create_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuCreateHandler(store)

	testCases := []struct {
		name       string
		postValues url.Values
		errorMsg   string
	}{
		{
			name: "Missing menu name",
			postValues: url.Values{
				"site_id": {"test-site-id"},
			},
			errorMsg: "menu name is required",
		},
		{
			name: "Missing site ID",
			postValues: url.Values{
				"menu_name": {"Test Menu"},
			},
			errorMsg: "site id is required",
		},
		{
			name:       "Empty form",
			postValues: url.Values{},
			errorMsg:   "site id is required",
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

func Test_MenuCreateController_Create_WithSiteList(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initMenuCreateHandler(store)

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, "Test Site 1")
	if err != nil {
		t.Fatalf("Failed to seed site1: %v", err)
	}

	site2, err := testutils.SeedSite(store, "Test Site 2")
	if err != nil {
		t.Fatalf("Failed to seed site2: %v", err)
	}

	// Test GET request to see site options
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Should contain both sites in the dropdown
	if !strings.Contains(body, site1.Name()) {
		t.Errorf("Expected body to contain site1 name '%s'", site1.Name())
	}
	if !strings.Contains(body, site2.Name()) {
		t.Errorf("Expected body to contain site2 name '%s'", site2.Name())
	}
	if !strings.Contains(body, "value=\""+site1.ID()+"\"") {
		t.Errorf("Expected body to contain site1 ID")
	}
	if !strings.Contains(body, "value=\""+site2.ID()+"\"") {
		t.Errorf("Expected body to contain site2 ID")
	}
}
