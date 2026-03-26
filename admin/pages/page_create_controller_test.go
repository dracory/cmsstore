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

func initPageCreateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewPageCreateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_PageCreateController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "New Page")
	assert.Contains(t, body, "name=\"page_name\"")
	assert.Contains(t, body, "name=\"site_id\"")
}

func Test_PageCreateController_Create(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageCreateHandler(store)

	// First seed a site
	site, err := testutils.SeedSite(store, testutils.SITE_01)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"page_name": {"Test New Page"},
			"site_id":   {site.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "page created successfully")

	// Verify in DB
	pages, _ := store.PageList(context.Background(), cmsstore.PageQuery().SetNameLike("Test New Page"))
	assert.Len(t, pages, 1)
	assert.Equal(t, "Test New Page", pages[0].Name())
	assert.Equal(t, site.ID(), pages[0].SiteID())
}

func Test_PageCreateController_Create_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageCreateHandler(store)

	testCases := []struct {
		name       string
		postValues url.Values
		errorMsg   string
	}{
		{
			name: "Missing page name",
			postValues: url.Values{
				"site_id": {"test-site-id"},
			},
			errorMsg: "page name is required",
		},
		{
			name: "Missing site ID",
			postValues: url.Values{
				"page_name": {"Test Page"},
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
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, response.StatusCode)

			assert.Contains(t, strings.ToLower(body), "error")
			assert.Contains(t, strings.ToLower(body), tc.errorMsg)
		})
	}
}

func Test_PageCreateController_Create_WithSiteList(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initPageCreateHandler(store)

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, testutils.SITE_01)
	require.NoError(t, err)

	site2, err := testutils.SeedSite(store, testutils.SITE_02)
	require.NoError(t, err)

	// Test GET request to see site options
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Should contain both sites in the dropdown
	assert.Contains(t, body, site1.Name())
	assert.Contains(t, body, site2.Name())
	assert.Contains(t, body, "value=\""+site1.ID()+"\"")
	assert.Contains(t, body, "value=\""+site2.ID()+"\"")
}
