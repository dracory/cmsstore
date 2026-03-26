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

func initMenuDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	return NewMenuDeleteController(UI(shared.UiConfig{
		Store:  store,
		Logger: slog.New(slog.NewTextHandler(os.Stderr, nil)),
		Layout: shared.Layout,
	})).Handler
}

func Test_MenuDeleteController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuDeleteHandler(store)

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

	assert.Contains(t, body, "Delete Menu")
	assert.Contains(t, body, "Are you sure you want to delete this menu?")
	assert.Contains(t, body, "This action cannot be undone")
	assert.Contains(t, body, "name=\"menu_id\"")
	assert.Contains(t, body, "value=\""+menu.ID()+"\"")
}

func Test_MenuDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuDeleteHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	// Verify menu exists before deletion
	menus, _ := store.MenuList(context.Background(), cmsstore.MenuQuery().SetID(menu.ID()))
	assert.Len(t, menus, 1)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"menu_id": {menu.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "menu deleted successfully")

	// Verify menu is soft deleted
	menus, _ = store.MenuList(context.Background(), cmsstore.MenuQuery().SetID(menu.ID()))
	assert.Len(t, menus, 0) // Should be empty since it's soft deleted
}

func Test_MenuDeleteController_Delete_ValidationError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuDeleteHandler(store)

	testCases := []struct {
		name       string
		postValues url.Values
		errorMsg   string
	}{
		{
			name:       "Missing menu ID",
			postValues: url.Values{},
			errorMsg:   "menu id is required",
		},
		{
			name: "Empty menu ID",
			postValues: url.Values{
				"menu_id": {""},
			},
			errorMsg: "menu id is required",
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

func Test_MenuDeleteController_MenuNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuDeleteHandler(store)

	// Test with non-existent menu ID
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {"non-existent-menu-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "menu not found")

	// Test POST with non-existent menu ID
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"menu_id": {"non-existent-menu-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "menu not found")
}

func Test_MenuDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initMenuDeleteHandler(store)

	// Seed a site and menu
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	menu := cmsstore.NewMenu()
	menu.SetSiteID(site.ID())
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	// First, GET the delete modal to confirm it shows correctly
	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: url.Values{
			"menu_id": {menu.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, body, menu.ID())

	// Then, POST to delete the menu
	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: url.Values{
			"menu_id": {menu.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")
}
