package admin

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

func initSiteManagerHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewSiteManagerController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_SiteManagerController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initSiteManagerHandler(store)

	// Seed a site
	_, err = testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Site Manager") {
		t.Errorf("Expected body to contain 'Site Manager'")
	}
	if !strings.Contains(body, "Test Site") {
		t.Errorf("Expected body to contain 'Test Site'")
	}
	if !strings.Contains(body, "New Site") {
		t.Errorf("Expected body to contain 'New Site'")
	}
}

func Test_SiteManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initSiteManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
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
	if !strings.Contains(body, "filter_status") {
		t.Errorf("Expected body to contain filter_status")
	}
	if !strings.Contains(body, "filter_name") {
		t.Errorf("Expected body to contain filter_name")
	}
}
