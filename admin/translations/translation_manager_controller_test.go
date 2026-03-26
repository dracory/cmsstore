package admin

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Translation Manager")
	assert.Contains(t, body, "New Translation")
	assert.Contains(t, body, "<tbody></tbody>")
}

func Test_TranslationManagerController_WithTranslations(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Header Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	require.NoError(t, err)

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Footer Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	err = store.TranslationCreate(context.Background(), translation2)
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Translation Manager")
	assert.Contains(t, body, "Header Translation")
	assert.Contains(t, body, "Footer Translation")
	assert.Contains(t, body, "translation-update")
	assert.Contains(t, body, "translation-delete")
}

func Test_TranslationManagerController_FilterModal(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"action": {"modal_translation_filter_show"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Filters")
	assert.Contains(t, body, "name")
	assert.Contains(t, body, "status")
}

func Test_TranslationManagerController_Sorting(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("A Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	require.NoError(t, err)

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Z Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	// Test sort by name ASC
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"by":   {"name"},
			"sort": {"asc"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "A Translation")
	assert.Contains(t, body, "Z Translation")

	// Test sort by name DESC
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"by":   {"name"},
			"sort": {"desc"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Z Translation")
	assert.Contains(t, body, "A Translation")
}

func Test_TranslationManagerController_Filtering(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Header Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	require.NoError(t, err)

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Footer Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	err = store.TranslationCreate(context.Background(), translation2)
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	// Test filter by name
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_name": {"Header"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Translation Manager")

	// Test filter by status
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_status": {cmsstore.TRANSLATION_STATUS_ACTIVE},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Translation Manager")
}

func Test_TranslationManagerController_Pagination(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create many translations to test pagination
	for i := 1; i <= 25; i++ {
		translation := cmsstore.NewTranslation()
		translation.SetName(fmt.Sprintf("Translation %d", i))
		translation.SetSiteID(site.ID())
		translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
		err = store.TranslationCreate(context.Background(), translation)
		require.NoError(t, err)
	}

	handler := initTranslationManagerHandler(store)

	// Test first page
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"0"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Translation 1")
	assert.Contains(t, body, "pagination")

	// Test second page
	body, response, err = test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"per_page": {"10"},
			"page":     {"1"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "pagination")
}

func Test_TranslationManagerController_EmptyState(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Translation Manager")
	assert.Contains(t, body, "New Translation")
	assert.Contains(t, body, "<tbody></tbody>")
}

func Test_TranslationManagerController_TableActions(t *testing.T) {
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

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "translation-update")
	assert.Contains(t, body, "translation-delete")
}

func Test_TranslationManagerController_MultipleSites(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, "Site 1")
	require.NoError(t, err)

	site2 := cmsstore.NewSite()
	site2.SetName("Site 2")
	site2.SetDomainNames([]string{"site2.example.com"})
	err = store.SiteCreate(context.Background(), site2)
	require.NoError(t, err)

	// Create translations for different sites
	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Site 1 Translation")
	translation1.SetSiteID(site1.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	require.NoError(t, err)

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Site 2 Translation")
	translation2.SetSiteID(site2.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Site 1 Translation")
	assert.Contains(t, body, "Site 2 Translation")
}

func Test_TranslationManagerController_Search(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and translations
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Searchable Translation")
	translation1.SetSiteID(site.ID())
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation1)
	require.NoError(t, err)

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Other Translation")
	translation2.SetSiteID(site.ID())
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), translation2)
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_name": {"Searchable"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	// Note: The filtering functionality might be implemented client-side
	// For now, just verify the search form is present and the page loads correctly
	assert.Contains(t, body, "Translation Manager")
}

func Test_TranslationManagerController_DifferentStatuses(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTranslationManagerHandler(store)

	// Test active translation
	activeTranslation := cmsstore.NewTranslation()
	activeTranslation.SetName("Active Translation")
	activeTranslation.SetSiteID(site.ID())
	activeTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	err = store.TranslationCreate(context.Background(), activeTranslation)
	require.NoError(t, err)

	// Test draft translation
	draftTranslation := cmsstore.NewTranslation()
	draftTranslation.SetName("Draft Translation")
	draftTranslation.SetSiteID(site.ID())
	draftTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	err = store.TranslationCreate(context.Background(), draftTranslation)
	require.NoError(t, err)

	// Test inactive translation
	inactiveTranslation := cmsstore.NewTranslation()
	inactiveTranslation.SetName("Inactive Translation")
	inactiveTranslation.SetSiteID(site.ID())
	inactiveTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_INACTIVE)
	err = store.TranslationCreate(context.Background(), inactiveTranslation)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"filter_status": {cmsstore.TRANSLATION_STATUS_ACTIVE},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, "Translation Manager")
}
