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

func initSiteDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewSiteDeleteController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_SiteDeleteController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initSiteDeleteHandler(store)

	site, _ := testutils.SeedSite(store, "Site to Delete")

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"site_id": {site.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Delete Site")
	// The modal might only show ID, check based on output from previous turn
	assert.Contains(t, body, site.ID())
}

func Test_SiteDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initSiteDeleteHandler(store)

	site, _ := testutils.SeedSite(store, "Site to Delete")

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"site_id": {site.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "site deleted successfully")

	// Verify in DB (soft deleted)
	// Use SiteList with SoftDeletedIncluded(true) because SiteFindByID filters them out
	list, err := store.SiteList(context.Background(), cmsstore.SiteQuery().SetID(site.ID()).SetSoftDeletedIncluded(true))
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.True(t, list[0].IsSoftDeleted())
}
