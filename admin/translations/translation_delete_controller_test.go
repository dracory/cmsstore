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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"translation_id": {translation.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Delete Translation")
	assert.Contains(t, body, "Are you sure you want to delete this translation?")
	assert.Contains(t, body, translation.ID())
}

func Test_TranslationDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translation
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetSiteID(site.ID())
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {translation.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "translation deleted successfully")

	// Verify translation is deleted (soft delete)
	deletedTranslation, err := store.TranslationFindByID(context.Background(), translation.ID())
	if err == nil && deletedTranslation != nil {
		assert.Equal(t, cmsstore.TRANSLATION_STATUS_INACTIVE, deletedTranslation.Status())
	}
}

func Test_TranslationDeleteController_Delete_ValidationError_MissingTranslationID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_TranslationDeleteController_Delete_ValidationError_EmptyTranslationID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {""},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_TranslationDeleteController_TranslationNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTranslationDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {"non-existent-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "translation not found")
}

func Test_TranslationDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTranslationDeleteHandler(store)

	// Create multiple translations
	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Translation 1")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	require.NoError(t, err)

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Translation 2")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	require.NoError(t, err)

	// Delete first translation
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"translation_id": {translation1.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")

	// Verify first translation is deleted but second translation remains
	deletedTranslation1, err := store.TranslationFindByID(context.Background(), translation1.ID())
	if err == nil && deletedTranslation1 != nil {
		assert.Equal(t, cmsstore.TRANSLATION_STATUS_INACTIVE, deletedTranslation1.Status())
	}

	activeTranslation2, err := store.TranslationFindByID(context.Background(), translation2.ID())
	if err == nil && activeTranslation2 != nil {
		assert.Equal(t, cmsstore.TRANSLATION_STATUS_ACTIVE, activeTranslation2.Status())
	}
}
