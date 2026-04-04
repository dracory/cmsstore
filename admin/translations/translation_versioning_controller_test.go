package admin

import (
	"context"
	"encoding/json"
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

func initTranslationVersioningHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewTranslationVersioningController(ui).Handler
}

func Test_TranslationVersioningController_TranslationIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

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
	if !strings.Contains(bodyLower, "translation id is required") {
		t.Errorf("Expected body to contain error message")
	}
}

func Test_TranslationVersioningController_TranslationIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {"invalid-id"},
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

func Test_TranslationVersioningController_ListRevisions(t *testing.T) {
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
	translation.SetContent(map[string]string{"en": "Original content"})
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create some versioning entries
	versioning1 := cmsstore.NewVersioning()
	versioning1.SetEntityID(translation.ID())
	versioning1.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content1, _ := json.Marshal(map[string]string{"en": "Version 1 content"})
	versioning1.SetContent(string(content1))
	err = store.VersioningCreate(context.Background(), versioning1)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	versioning2 := cmsstore.NewVersioning()
	versioning2.SetEntityID(translation.ID())
	versioning2.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content2, _ := json.Marshal(map[string]string{"en": "Version 2 content"})
	versioning2.SetContent(string(content2))
	err = store.VersioningCreate(context.Background(), versioning2)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

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

	if !strings.Contains(body, "Translation Revisions") {
		t.Errorf("Expected body to contain 'Translation Revisions'")
	}
	if !strings.Contains(body, "Preview") {
		t.Errorf("Expected body to contain 'Preview'")
	}
}

func Test_TranslationVersioningController_PreviewRevision(t *testing.T) {
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
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(translation.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content, _ := json.Marshal(map[string]string{"en": "Historical content"})
	versioning.SetContent(string(content))
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
			"versioning_id":  {versioning.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Translation Revision") {
		t.Errorf("Expected body to contain 'Translation Revision'")
	}
	if !strings.Contains(body, "Close") {
		t.Errorf("Expected body to contain 'Close'")
	}
}

func Test_TranslationVersioningController_RestoreAttributes(t *testing.T) {
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
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create a versioning entry with different content
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(translation.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	// Content must be stored as a JSON string that can be unmarshalled into map[string]string
	contentValue, _ := json.Marshal(map[string]string{"en": "Restored content"})
	content, _ := json.Marshal(map[string]any{
		cmsstore.COLUMN_CONTENT: string(contentValue),
	})
	versioning.SetContent(string(content))
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id":      {translation.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {cmsstore.COLUMN_CONTENT},
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
	if !strings.Contains(bodyLower, "revision attributes restored successfully") {
		t.Errorf("Expected body to contain success message")
	}

	// Verify translation content was restored
	restoredTranslation, err := store.TranslationFindByID(context.Background(), translation.ID())
	if err != nil {
		t.Fatalf("Failed to find translation: %v", err)
	}
	restoredContent, _ := restoredTranslation.Content()
	if restoredContent["en"] != "Restored content" {
		t.Errorf("Expected restored content 'Restored content', got '%s'", restoredContent["en"])
	}
}

func Test_TranslationVersioningController_RestoreNoAttributes(t *testing.T) {
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
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(translation.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content, _ := json.Marshal(map[string]string{"en": "Restored content"})
	versioning.SetContent(string(content))
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {translation.ID()},
			"versioning_id":  {versioning.ID()},
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
	if !strings.Contains(bodyLower, "no revision attributes were selected") {
		t.Errorf("Expected body to contain error message")
	}
}

func Test_TranslationVersioningController_RestoreMultipleAttributes(t *testing.T) {
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
	translation.SetName("Original Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	translation.SetContent(map[string]string{"en": "Original content"})
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create a versioning entry with different values
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(translation.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content, _ := json.Marshal(map[string]string{
		cmsstore.COLUMN_NAME:   "Restored Translation",
		cmsstore.COLUMN_STATUS: cmsstore.TRANSLATION_STATUS_ACTIVE,
	})
	versioning.SetContent(string(content))
	err = store.VersioningCreate(context.Background(), versioning)
	if err != nil {
		t.Fatalf("Failed to create versioning: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id":      {translation.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {cmsstore.COLUMN_NAME, cmsstore.COLUMN_STATUS},
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
	if !strings.Contains(bodyLower, "revision attributes restored successfully") {
		t.Errorf("Expected body to contain success message")
	}
}

func Test_TranslationVersioningController_VersioningNotFound(t *testing.T) {
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
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
			"versioning_id":  {"non-existent-versioning-id"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	// The versioning controller should show existing revisions even with invalid versioning_id
	if !strings.Contains(strings.ToLower(body), "translation revisions") {
		t.Errorf("Expected body to contain 'translation revisions'")
	}
}

func Test_TranslationVersioningController_EmptyRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and translation (no versioning entries)
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	handler := initTranslationVersioningHandler(store)

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

	if !strings.Contains(body, "Translation Revisions") {
		t.Errorf("Expected body to contain 'Translation Revisions'")
	}
	if !strings.Contains(body, "Version") {
		t.Errorf("Expected body to contain 'Version'")
	}
}
