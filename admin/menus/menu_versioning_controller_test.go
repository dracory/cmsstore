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

func initMenuVersioningHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuVersioningController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuVersioningController_MenuIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "menu id is required")
}

func Test_MenuVersioningController_MenuIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuVersioningHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {"invalid-menu-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "menu not found")
}

func Test_MenuVersioningController_ListRevisions(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuVersioningHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Menu Revisions")
}

func Test_MenuVersioningController_PreviewRevision(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuVersioningHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(menu.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_MENU)
	versioning.SetContent(`{"name": "Test Menu", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id":       {menu.ID()},
			"versioning_id": {versioning.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Menu Revision")
	assert.Contains(t, body, "Test Menu")
}

func Test_MenuVersioningController_RestoreAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuVersioningHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	// Create a versioning entry with different data
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(menu.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_MENU)
	versioning.SetContent(`{"name": "Restored Menu", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id":       {menu.ID()},
			"versioning_id": {versioning.ID()},
		},
		PostValues: url.Values{
			"action":              {"restore"},
			"revision_attributes": {"name"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "restored successfully")

	// Verify the menu was updated
	updatedMenu, err := store.MenuFindByID(context.Background(), menu.ID())
	require.NoError(t, err)
	require.NotNil(t, updatedMenu)
	assert.Equal(t, "Restored Menu", updatedMenu.Name())
}

func Test_MenuVersioningController_RestoreNoAttributes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuVersioningHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	// Create a versioning entry
	versioning := cmsstore.NewVersioning()
	versioning.SetEntityID(menu.ID())
	versioning.SetEntityType(cmsstore.VERSIONING_TYPE_MENU)
	versioning.SetContent(`{"name": "Test Menu", "status": "active"}`)
	err = store.VersioningCreate(context.Background(), versioning)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id":       {menu.ID()},
			"versioning_id": {versioning.ID()},
		},
		PostValues: url.Values{
			"action":              {"restore"},
			"revision_attributes": {"name"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "restored successfully")
}
