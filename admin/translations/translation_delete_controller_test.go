package admin

import (
	"context"
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

func initTranslationDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewTranslationDeleteController(ui).Handler
}

func Test_TranslationDeleteController_Index(t *testing.T) {
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

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Delete Translation") {
		t.Errorf("Expected body to contain 'Delete Translation'")
	}
	if !strings.Contains(body, "Are you sure you want to delete this translation?") {
		t.Errorf("Expected body to contain confirmation message")
	}
	if !strings.Contains(body, translation.ID()) {
		t.Errorf("Expected body to contain translation ID")
	}
}

func Test_TranslationDeleteController_Delete(t *testing.T) {
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

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {translation.ID()},
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
	if !strings.Contains(bodyLower, "translation deleted successfully") {
		t.Errorf("Expected body to contain success message")
	}

	// Verify translation is deleted (soft delete)
	deletedTranslation, err := store.TranslationFindByID(context.Background(), translation.ID())
	if err == nil && deletedTranslation != nil {
		if deletedTranslation.Status() != cmsstore.TRANSLATION_STATUS_INACTIVE {
			t.Errorf("Expected translation status to be INACTIVE, got %s", deletedTranslation.Status())
		}
	}
}

func Test_TranslationDeleteController_Delete_ValidationError_MissingTranslationID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(strings.ToLower(body), "error") {
		t.Errorf("Expected body to contain 'error'")
	}
}

func Test_TranslationDeleteController_Delete_ValidationError_EmptyTranslationID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {""},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(strings.ToLower(body), "error") {
		t.Errorf("Expected body to contain 'error'")
	}
}

func Test_TranslationDeleteController_TranslationNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {"non-existent-id"},
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
	if !strings.Contains(bodyLower, "translation not found") {
		t.Errorf("Expected body to contain error message")
	}
}

func Test_TranslationDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initTranslationDeleteHandler(store)

	// Create multiple translations
	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Translation 1")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Translation 2")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Delete first translation
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {translation1.ID()},
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

	// Verify first translation is deleted but second translation remains
	deletedTranslation1, err := store.TranslationFindByID(context.Background(), translation1.ID())
	if err == nil && deletedTranslation1 != nil {
		if deletedTranslation1.Status() != cmsstore.TRANSLATION_STATUS_INACTIVE {
			t.Errorf("Expected deleted translation status to be INACTIVE, got %s", deletedTranslation1.Status())
		}
	}

	activeTranslation2, err := store.TranslationFindByID(context.Background(), translation2.ID())
	if err == nil && activeTranslation2 != nil {
		if activeTranslation2.Status() != cmsstore.TRANSLATION_STATUS_ACTIVE {
			t.Errorf("Expected active translation status to be ACTIVE, got %s", activeTranslation2.Status())
		}
	}
}
