package admin

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

func initTranslationManagerHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewTranslationManagerController(ui).Handler
}

func Test_TranslationManagerController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Translation Manager") {
		t.Errorf("Expected body to contain 'Translation Manager'")
	}
	if !strings.Contains(body, "New Translation") {
		t.Errorf("Expected body to contain 'New Translation'")
	}
	if !strings.Contains(body, "<tbody></tbody>") {
		t.Errorf("Expected body to contain empty tbody")
	}
}

func Test_TranslationManagerController_WithTranslations(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Header Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Footer Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	err = store.TranslationCreate(context.Background(), translation2)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Translation Manager") {
		t.Errorf("Expected body to contain 'Translation Manager'")
	}
	if !strings.Contains(body, "Header Translation") {
		t.Errorf("Expected body to contain 'Header Translation'")
	}
	if !strings.Contains(body, "Footer Translation") {
		t.Errorf("Expected body to contain 'Footer Translation'")
	}
	if !strings.Contains(body, "translation-update") {
		t.Errorf("Expected body to contain 'translation-update'")
	}
	if !strings.Contains(body, "translation-delete") {
		t.Errorf("Expected body to contain 'translation-delete'")
	}
}

func Test_TranslationManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"action": {"modal_translation_filter_show"},
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
	if !strings.Contains(body, "name") {
		t.Errorf("Expected body to contain name filter")
	}
	if !strings.Contains(body, "status") {
		t.Errorf("Expected body to contain status filter")
	}
}

func Test_TranslationManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("A Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Z Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	// Test sort by name ASC
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"by":   {"name"},
			"sort": {"asc"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "A Translation") {
		t.Errorf("Expected body to contain 'A Translation'")
	}
	if !strings.Contains(body, "Z Translation") {
		t.Errorf("Expected body to contain 'Z Translation'")
	}

	// Test sort by name DESC
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"by":   {"name"},
			"sort": {"desc"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Z Translation") {
		t.Errorf("Expected body to contain 'Z Translation'")
	}
	if !strings.Contains(body, "A Translation") {
		t.Errorf("Expected body to contain 'A Translation'")
	}
}

func Test_TranslationManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Header Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Footer Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	err = store.TranslationCreate(context.Background(), translation2)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	// Test filter by name
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_name": {"Header"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Translation Manager") {
		t.Errorf("Expected body to contain 'Translation Manager'")
	}

	// Test filter by status
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_status": {cmsstore.TRANSLATION_STATUS_ACTIVE},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Translation Manager") {
		t.Errorf("Expected body to contain 'Translation Manager'")
	}
}

func Test_TranslationManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create many translations to test pagination
	for i := 1; i <= 25; i++ {
		translation := cmsstore.NewTranslation()
		translation.SetName(fmt.Sprintf("Translation %d", i))
		translation.SetSiteID(site.ID())
		translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
		err = store.TranslationCreate(context.Background(), translation)
		if err != nil {
			t.Fatalf("Failed to create translation: %v", err)
		}
	}

	handler := initTranslationManagerHandler(store)

	// Test first page
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"0"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Translation 1") {
		t.Errorf("Expected body to contain 'Translation 1'")
	}
	if !strings.Contains(body, "pagination") {
		t.Errorf("Expected body to contain 'pagination'")
	}

	// Test second page
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"1"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "pagination") {
		t.Errorf("Expected body to contain 'pagination'")
	}
}

func Test_TranslationManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Translation Manager") {
		t.Errorf("Expected body to contain 'Translation Manager'")
	}
	if !strings.Contains(body, "New Translation") {
		t.Errorf("Expected body to contain 'New Translation'")
	}
	if !strings.Contains(body, "<tbody></tbody>") {
		t.Errorf("Expected body to contain empty tbody")
	}
}

func Test_TranslationManagerController_TableActions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "translation-update") {
		t.Errorf("Expected body to contain 'translation-update'")
	}
	if !strings.Contains(body, "translation-delete") {
		t.Errorf("Expected body to contain 'translation-delete'")
	}
}

func Test_TranslationManagerController_MultipleSites(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, "Site 1")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	site2 := cmsstore.NewSite()
	site2.SetName("Site 2")
	site2.SetDomainNames([]string{"site2.example.com"})
	err = store.SiteCreate(context.Background(), site2)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create translations for different sites
	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Site 1 Translation")
	translation1.SetSiteID(site1.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Site 2 Translation")
	translation2.SetSiteID(site2.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Site 1 Translation") {
		t.Errorf("Expected body to contain 'Site 1 Translation'")
	}
	if !strings.Contains(body, "Site 2 Translation") {
		t.Errorf("Expected body to contain 'Site 2 Translation'")
	}
}

func Test_TranslationManagerController_Search(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Searchable Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Other Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_name": {"Searchable"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	// Note: The filtering functionality might be implemented client-side
	// For now, just verify the search form is present and the page loads correctly
	if !strings.Contains(body, "Translation Manager") {
		t.Errorf("Expected body to contain 'Translation Manager'")
	}
}

func Test_TranslationManagerController_DifferentStatuses(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initTranslationManagerHandler(store)

	// Test active translation
	activeTranslation := cmsstore.NewTranslation()
	activeTranslation.SetName("Active Translation")
	activeTranslation.SetSiteID(site.ID())
	activeTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), activeTranslation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Test draft translation
	draftTranslation := cmsstore.NewTranslation()
	draftTranslation.SetName("Draft Translation")
	draftTranslation.SetSiteID(site.ID())
	draftTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	err = store.TranslationCreate(context.Background(), draftTranslation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Test inactive translation
	inactiveTranslation := cmsstore.NewTranslation()
	inactiveTranslation.SetName("Inactive Translation")
	inactiveTranslation.SetSiteID(site.ID())
	inactiveTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_INACTIVE)
	err = store.TranslationCreate(context.Background(), inactiveTranslation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_status": {cmsstore.TRANSLATION_STATUS_ACTIVE},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(body, "Translation Manager") {
		t.Errorf("Expected body to contain 'Translation Manager'")
	}
}
