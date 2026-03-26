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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func initPageDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewPageDeleteController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_PageDeleteController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageDeleteHandler(store)

	// Seed a site and page
	site, err := testutils.SeedSite(store, testutils.SITE_01)
	require.NoError(t, err)

	page, err := testutils.SeedPage(store, testutils.PAGE_01, site.ID())
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {page.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Delete Page")
	assert.Contains(t, body, "Are you sure you want to delete this page?")
	assert.Contains(t, body, "This action cannot be undone")
	assert.Contains(t, body, "name=\"page_id\"")
	assert.Contains(t, body, "value=\""+page.ID()+"\"")
}

func Test_PageDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageDeleteHandler(store)

	// Seed a site and page
	site, err := testutils.SeedSite(store, testutils.SITE_01)
	require.NoError(t, err)

	page, err := testutils.SeedPage(store, testutils.PAGE_01, site.ID())
	require.NoError(t, err)

	// Verify page exists before deletion
	pages, _ := store.PageList(context.Background(), cmsstore.PageQuery().SetID(page.ID()))
	assert.Len(t, pages, 1)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"page_id": {page.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "page deleted successfully")

	// Verify page is soft deleted
	pages, _ = store.PageList(context.Background(), cmsstore.PageQuery().SetID(page.ID()))
	assert.Len(t, pages, 0) // Should be empty since it's soft deleted
}

func Test_PageDeleteController_Delete_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageDeleteHandler(store)

	testCases := []struct {
		name       string
		postValues url.Values
		errorMsg   string
	}{
		{
			name: "Missing page ID",
			postValues: url.Values{},
			errorMsg:   "page id is required",
		},
		{
			name: "Empty page ID",
			postValues: url.Values{
				"page_id": {""},
			},
			errorMsg: "page id is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
				PostValues: tc.postValues,
			})
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)

			assert.Contains(t, strings.ToLower(body), "error")
			assert.Contains(t, strings.ToLower(body), tc.errorMsg)
		})
	}
}

func Test_PageDeleteController_PageNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageDeleteHandler(store)

	// Test with non-existent page ID
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {"non-existent-page-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "page not found")

	// Test POST with non-existent page ID
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"page_id": {"non-existent-page-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "page not found")
}

func Test_PageDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageDeleteHandler(store)

	// Seed a site and page
	site, err := testutils.SeedSite(store, testutils.SITE_01)
	require.NoError(t, err)

	page, err := testutils.SeedPage(store, testutils.PAGE_01, site.ID())
	require.NoError(t, err)

	// First, GET the delete modal to confirm it shows correctly
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page_id": {page.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, page.ID())

	// Then, POST to delete the page
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"page_id": {page.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")
}
