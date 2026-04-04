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

func initSiteDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewSiteDeleteController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_SiteDeleteController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initSiteDeleteHandler(store)

	site, _ := testutils.SeedSite(store, "Site to Delete")

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"site_id": {site.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Delete Site") {
		t.Errorf("Expected body to contain 'Delete Site'")
	}
	// The modal might only show ID, check based on output from previous turn
	if !strings.Contains(body, site.ID()) {
		t.Errorf("Expected body to contain site ID")
	}
}

func Test_SiteDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initSiteDeleteHandler(store)

	site, _ := testutils.SeedSite(store, "Site to Delete")

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"site_id": {site.ID()},
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
	if !strings.Contains(bodyLower, "site deleted successfully") {
		t.Errorf("Expected body to contain success message")
	}

	// Verify in DB (soft deleted)
	// Use SiteList with SoftDeletedIncluded(true) because SiteFindByID filters them out
	list, err := store.SiteList(context.Background(), cmsstore.SiteQuery().SetID(site.ID()).SetSoftDeletedIncluded(true))
	if err != nil {
		t.Fatalf("Failed to list sites: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("Expected 1 site, got %d", len(list))
	}
	if !list[0].IsSoftDeleted() {
		t.Errorf("Expected site to be soft deleted")
	}
}
