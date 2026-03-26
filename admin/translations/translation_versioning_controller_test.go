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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "translation id is required")
}

func Test_TranslationVersioningController_TranslationIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {"invalid-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "translation not found")
}

func Test_TranslationVersioningController_ListRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetContent(map[string]string{"en": "Original content"})
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	// Create some versioning entries
	versioning1 := cmsstore.NewVersioning()
	versioning1.SetEntityID(translation.ID())
	versioning1.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content1, _ := json.Marshal(map[string]string{"en": "Version 1 content"})
	versioning1.SetContent(string(content1))
	err = store.VersioningCreate(context.Background(), versioning1)
	require.NoError(t, err)

	versioning2 := cmsstore.NewVersioning()
	versioning2.SetEntityID(translation.ID())
	versioning2.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content2, _ := json.Marshal(map[string]string{"en": "Version 2 content"})
	versioning2.SetContent(string(content2))
	err = store.VersioningCreate(context.Background(), versioning2)
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Translation Revisions")
	assert.Contains(t, body, "Preview")
}

func Test_TranslationVersioningController_PreviewRevision(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(translation.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content, _ := json.Marshal(map[string]string{"en": "Historical content"})
	versioning.SetContent(string(content))
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
			"versioning_id":  {versioning.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Translation Revision")
	assert.Contains(t, body, "Close")
}

func Test_TranslationVersioningController_RestoreAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

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
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id":      {translation.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {cmsstore.COLUMN_CONTENT},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "revision attributes restored successfully")

	// Verify translation content was restored
	restoredTranslation, err := store.TranslationFindByID(context.Background(), translation.ID())
	require.NoError(t, err)
	restoredContent, _ := restoredTranslation.Content()
	assert.Equal(t, "Restored content", restoredContent["en"])
}

func Test_TranslationVersioningController_RestoreNoAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(translation.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_TRANSLATION)
	content, _ := json.Marshal(map[string]string{"en": "Restored content"})
	versioning.SetContent(string(content))
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {translation.ID()},
			"versioning_id":  {versioning.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "no revision attributes were selected")
}

func Test_TranslationVersioningController_RestoreMultipleAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Original Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	translation.SetContent(map[string]string{"en": "Original content"})
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

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
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id":      {translation.ID()},
			"versioning_id":       {versioning.ID()},
			"revision_attributes": {cmsstore.COLUMN_NAME, cmsstore.COLUMN_STATUS},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "revision attributes restored successfully")
}

func Test_TranslationVersioningController_VersioningNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
			"versioning_id":  {"non-existent-versioning-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// The versioning controller should show existing revisions even with invalid versioning_id
	assert.Contains(t, strings.ToLower(body), "translation revisions")
}

func Test_TranslationVersioningController_EmptyRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation (no versioning entries)
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetContent(map[string]string{"en": "Current content"})
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	handler := initTranslationVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Translation Revisions")
	assert.Contains(t, body, "Version")
}
