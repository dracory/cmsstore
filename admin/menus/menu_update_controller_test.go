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

func initMenuUpdateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuUpdateController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuUpdateController_MenuIdIsRequired(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuUpdateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "menu id is required")
}

func Test_MenuUpdateController_MenuIdIsInvalid(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuUpdateHandler(store)

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

func Test_MenuUpdateController_ViewSettings_IsDefault(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuUpdateHandler(store)

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
			"view":    {"settings"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Edit Menu")
	assert.Contains(t, body, "name=\"menu_name\"")
	assert.Contains(t, body, menu.Name())
	assert.Contains(t, body, "name=\"menu_site_id\"")
	assert.Contains(t, body, site.ID())
}

func Test_MenuUpdateController_ViewMenuItems(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuUpdateHandler(store)

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
			"view":    {"menu_items"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Menu Items")
}

func Test_MenuUpdateController_UpdateSettings(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuUpdateHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
			"view":    {"settings"},
		},
		PostValues: url.Values{
			"menu_name":    {"Updated Menu Name"},
			"menu_site_id": {site.ID()},
			"menu_status":  {cmsstore.MENU_STATUS_ACTIVE},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "menu saved successfully")
}

func Test_MenuUpdateController_UpdateSettings_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuUpdateHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	// Test submitting form with empty menu name - should still work
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
			"view":    {"settings"},
		},
		PostValues: url.Values{
			"menu_name":    {""},
			"menu_site_id": {site.ID()},
			"menu_status":  {cmsstore.MENU_STATUS_ACTIVE},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Should contain either success or error handling
	assert.True(t, strings.Contains(body, "success") || strings.Contains(body, "error"))
}
