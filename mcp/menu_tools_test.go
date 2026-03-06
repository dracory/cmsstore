package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestMenuGet(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	tests := []struct {
		name        string
		menuID      string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "get menu with full ID",
			menuID:      menu.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menu.ID()),
		},
		{
			name:        "get menu with shortened ID",
			menuID:      cmsstore.ShortenID(menu.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menu.ID()),
		},
		{
			name:        "get non-existent menu",
			menuID:      "non_existent_id",
			expectError: true,
			expectedErr: "menu not found",
		},
		{
			name:        "get menu with empty ID",
			menuID:      "",
			expectError: true,
			expectedErr: "missing required parameter: id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the tool
			getPayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "get",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "menu_get",
					"arguments": map[string]any{
						"id": tt.menuID,
					},
				},
			}

			getBody, err := json.Marshal(getPayload)
			require.NoError(t, err)

			getResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(getBody))
			require.NoError(t, err)
			defer getResp.Body.Close()

			getRespBytes, err := io.ReadAll(getResp.Body)
			require.NoError(t, err)

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(getRespBytes, &response)
			require.NoError(t, err)

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				assert.True(t, hasError, "Expected error in response")
				if hasError {
					errorObj := response["error"].(map[string]any)
					assert.Equal(t, tt.expectedErr, errorObj["message"])
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				require.True(t, ok, "Expected response to have result")

				content, ok := result["content"].([]any)
				require.True(t, ok, "Expected response result.content")
				require.Len(t, content, 1, "Expected response result.content to have one item")

				item0, ok := content[0].(map[string]any)
				require.True(t, ok, "Expected response result.content[0] object")

				text, ok := item0["text"].(string)
				require.True(t, ok, "Expected response result.content[0].text")

				var menuData map[string]any
				err = json.Unmarshal([]byte(text), &menuData)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedID, menuData["id"].(string))
				assert.Equal(t, "Main Menu", menuData["name"].(string))
				assert.Equal(t, cmsstore.MENU_STATUS_ACTIVE, menuData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), menuData["site_id"].(string))
			}
		})
	}
}

func TestMenuList(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create menus with different properties
	activeMenu := cmsstore.NewMenu()
	activeMenu.SetName("Main Menu")
	activeMenu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	activeMenu.SetSiteID(site.ID())
	activeMenu.SetHandle("main-menu")
	activeMenu.SetMemo("Active menu")
	err = store.MenuCreate(context.Background(), activeMenu)
	require.NoError(t, err)

	draftMenu := cmsstore.NewMenu()
	draftMenu.SetName("Draft Menu")
	draftMenu.SetStatus(cmsstore.MENU_STATUS_DRAFT)
	draftMenu.SetSiteID(site.ID())
	draftMenu.SetHandle("draft-menu")
	draftMenu.SetMemo("Draft menu")
	err = store.MenuCreate(context.Background(), draftMenu)
	require.NoError(t, err)

	// Test listing all menus
	t.Run("list all menus", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"limit":  10,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return both menus
		assert.Equal(t, 2, len(items), "Expected both menus to be returned")
	})

	// Test filtering by site_id
	t.Run("list menus by site_id", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"site_id": cmsstore.ShortenID(site.ID()),
					"limit":   10,
					"offset":  0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return both menus for the site
		assert.Equal(t, 2, len(items), "Expected both menus for the site")
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			assert.Equal(t, cmsstore.ShortenID(site.ID()), itemMap["site_id"].(string))
		}
	})

	// Test filtering by status
	t.Run("list menus by status", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"status": cmsstore.MENU_STATUS_ACTIVE,
					"limit":  10,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only active menu
		assert.Equal(t, 1, len(items), "Expected only active menu")
		item := items[0].(map[string]interface{})
		assert.Equal(t, cmsstore.MENU_STATUS_ACTIVE, item["status"].(string))
	})

	// Test filtering by handle
	t.Run("list menus by handle", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"handle": "main-menu",
					"limit":  10,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only the menu with matching handle
		assert.Equal(t, 1, len(items), "Expected only menu with matching handle")
		item := items[0].(map[string]interface{})
		assert.Equal(t, "main-menu", item["handle"].(string))
	})

	// Test pagination
	t.Run("list menus with pagination", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"limit":  1,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only 1 menu due to limit
		assert.Equal(t, 1, len(items), "Expected only 1 menu due to limit")
	})
}

func TestMenuUpsert_Create(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	tests := []struct {
		name        string
		menuName    string
		status      string
		siteID      string
		handle      string
		memo        string
		expectError bool
		expectedErr string
	}{
		{
			name:        "create menu with all fields",
			menuName:    "New Menu",
			status:      cmsstore.MENU_STATUS_ACTIVE,
			siteID:      cmsstore.ShortenID(site.ID()),
			handle:      "new-menu",
			memo:        "Test memo",
			expectError: false,
		},
		{
			name:        "create menu with minimal fields",
			menuName:    "Minimal Menu",
			status:      cmsstore.MENU_STATUS_DRAFT,
			siteID:      "",
			handle:      "",
			memo:        "",
			expectError: false,
		},
		{
			name:        "create menu with empty name",
			menuName:    "",
			status:      cmsstore.MENU_STATUS_ACTIVE,
			siteID:      cmsstore.ShortenID(site.ID()),
			handle:      "test-menu",
			memo:        "Test memo",
			expectError: true,
			expectedErr: "missing required parameter: name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the tool
			upsertPayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "upsert",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "menu_upsert",
					"arguments": map[string]any{
						"name":    tt.menuName,
						"status":  tt.status,
						"site_id": tt.siteID,
						"handle":  tt.handle,
						"memo":    tt.memo,
					},
				},
			}

			upsertBody, err := json.Marshal(upsertPayload)
			require.NoError(t, err)

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			require.NoError(t, err)
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			require.NoError(t, err)

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
			require.NoError(t, err)

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				assert.True(t, hasError, "Expected error in response")
				if hasError {
					errorObj := response["error"].(map[string]any)
					assert.Equal(t, tt.expectedErr, errorObj["message"])
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				require.True(t, ok, "Expected response to have result")

				content, ok := result["content"].([]any)
				require.True(t, ok, "Expected response result.content")
				require.Len(t, content, 1, "Expected response result.content to have one item")

				item0, ok := content[0].(map[string]any)
				require.True(t, ok, "Expected response result.content[0] object")

				text, ok := item0["text"].(string)
				require.True(t, ok, "Expected response result.content[0].text")

				var menuData map[string]any
				err = json.Unmarshal([]byte(text), &menuData)
				require.NoError(t, err)

				assert.Equal(t, tt.menuName, menuData["name"].(string))
				assert.Equal(t, tt.status, menuData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), menuData["site_id"].(string))
			}
		})
	}
}

func TestMenuUpsert_Update(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Original Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("original-menu")
	menu.SetMemo("Original memo")
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	tests := []struct {
		name        string
		menuID      string
		menuName    string
		status      string
		handle      string
		expectError bool
		expectedErr string
	}{
		{
			name:        "update menu with full ID",
			menuID:      menu.ID(),
			menuName:    "Updated Menu",
			status:      cmsstore.MENU_STATUS_DRAFT,
			handle:      "updated-menu",
			expectError: false,
		},
		{
			name:        "update menu with shortened ID",
			menuID:      cmsstore.ShortenID(menu.ID()),
			menuName:    "Updated Menu",
			status:      cmsstore.MENU_STATUS_DRAFT,
			handle:      "updated-menu",
			expectError: false,
		},
		{
			name:        "update non-existent menu",
			menuID:      "non_existent_id",
			menuName:    "Updated Menu",
			status:      cmsstore.MENU_STATUS_DRAFT,
			handle:      "updated-menu",
			expectError: true,
			expectedErr: "menu not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the tool
			upsertPayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "upsert",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "menu_upsert",
					"arguments": map[string]any{
						"id":     tt.menuID,
						"name":   tt.menuName,
						"status": tt.status,
						"handle": tt.handle,
					},
				},
			}

			upsertBody, err := json.Marshal(upsertPayload)
			require.NoError(t, err)

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			require.NoError(t, err)
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			require.NoError(t, err)

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
			require.NoError(t, err)

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				assert.True(t, hasError, "Expected error in response")
				if hasError {
					errorObj := response["error"].(map[string]any)
					assert.Equal(t, tt.expectedErr, errorObj["message"])
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				require.True(t, ok, "Expected response to have result")

				content, ok := result["content"].([]any)
				require.True(t, ok, "Expected response result.content")
				require.Len(t, content, 1, "Expected response result.content to have one item")

				item0, ok := content[0].(map[string]any)
				require.True(t, ok, "Expected response result.content[0] object")

				text, ok := item0["text"].(string)
				require.True(t, ok, "Expected response result.content[0].text")

				var menuData map[string]any
				err = json.Unmarshal([]byte(text), &menuData)
				require.NoError(t, err)

				assert.Equal(t, tt.menuName, menuData["name"].(string))
				assert.Equal(t, tt.status, menuData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), menuData["site_id"].(string))
			}
		})
	}
}

func TestMenuDelete(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	tests := []struct {
		name        string
		menuID      string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "delete menu with full ID",
			menuID:      menu.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menu.ID()),
		},
		{
			name:        "delete menu with shortened ID",
			menuID:      cmsstore.ShortenID(menu.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menu.ID()),
		},
		{
			name:        "delete non-existent menu",
			menuID:      "non_existent_id",
			expectError: true,
			expectedErr: "menu not found",
		},
		{
			name:        "delete menu with empty ID",
			menuID:      "",
			expectError: true,
			expectedErr: "missing required parameter: id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetID := tt.menuID
			if tt.name == "delete menu with full ID" || tt.name == "delete menu with shortened ID" {
				// Create a fresh menu for each positive test case
				m := cmsstore.NewMenu()
				m.SetName("Test Menu")
				m.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
				m.SetSiteID(site.ID())
				err = store.MenuCreate(context.Background(), m)
				require.NoError(t, err)

				if tt.name == "delete menu with full ID" {
					targetID = m.ID()
				} else {
					targetID = cmsstore.ShortenID(m.ID())
				}
				// Update expectedID to match the new menu
				tt.expectedID = cmsstore.ShortenID(m.ID())
			}

			// Call the tool
			deletePayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "delete",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "menu_delete",
					"arguments": map[string]any{
						"id": targetID,
					},
				},
			}

			deleteBody, err := json.Marshal(deletePayload)
			require.NoError(t, err)

			deleteResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(deleteBody))
			require.NoError(t, err)
			defer deleteResp.Body.Close()

			deleteRespBytes, err := io.ReadAll(deleteResp.Body)
			require.NoError(t, err)

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(deleteRespBytes, &response)
			require.NoError(t, err)

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				assert.True(t, hasError, "Expected error in response")
				if hasError {
					errorObj := response["error"].(map[string]any)
					assert.Equal(t, tt.expectedErr, errorObj["message"])
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				require.True(t, ok, "Expected response to have result")

				content, ok := result["content"].([]any)
				require.True(t, ok, "Expected response result.content")
				require.Len(t, content, 1, "Expected response result.content to have one item")

				item0, ok := content[0].(map[string]any)
				require.True(t, ok, "Expected response result.content[0] object")

				text, ok := item0["text"].(string)
				require.True(t, ok, "Expected response result.content[0].text")

				var deleteData map[string]any
				err = json.Unmarshal([]byte(text), &deleteData)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedID, deleteData["id"].(string))
			}
		})
	}
}

func TestMenuUpsert_WithDefaultSite(t *testing.T) {
	server, _, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a menu without specifying site_id - should use default site
	upsertPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "upsert",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "menu_upsert",
			"arguments": map[string]any{
				"name":   "Default Site Menu",
				"status": cmsstore.MENU_STATUS_ACTIVE,
				"handle": "default-menu",
				"memo":   "Menu with default site",
			},
		},
	}

	upsertBody, err := json.Marshal(upsertPayload)
	require.NoError(t, err)

	upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
	require.NoError(t, err)
	defer upsertResp.Body.Close()

	upsertRespBytes, err := io.ReadAll(upsertResp.Body)
	require.NoError(t, err)

	// Parse the result
	var response map[string]any
	err = json.Unmarshal(upsertRespBytes, &response)
	require.NoError(t, err)

	// Check for success
	result, ok := response["result"].(map[string]any)
	require.True(t, ok, "Expected response to have result")

	content, ok := result["content"].([]any)
	require.True(t, ok, "Expected response result.content")
	require.Len(t, content, 1, "Expected response result.content to have one item")

	item0, ok := content[0].(map[string]any)
	require.True(t, ok, "Expected response result.content[0] object")

	text, ok := item0["text"].(string)
	require.True(t, ok, "Expected response result.content[0].text")

	var menuData map[string]any
	err = json.Unmarshal([]byte(text), &menuData)
	require.NoError(t, err)

	assert.Equal(t, "Default Site Menu", menuData["name"].(string))
	assert.Equal(t, cmsstore.MENU_STATUS_ACTIVE, menuData["status"].(string))
	assert.Equal(t, "default-menu", menuData["handle"].(string))
	assert.Equal(t, "Menu with default site", menuData["memo"].(string))
	// site_id should be set to the default site
	assert.NotEmpty(t, menuData["site_id"].(string))
}

func TestMenuList_WithSoftDeleted(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Soft Deleted Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu)
	require.NoError(t, err)

	// Soft delete the menu
	err = store.MenuSoftDeleteByID(context.Background(), menu.ID())
	require.NoError(t, err)

	// Test listing without include_soft_deleted (should not include soft deleted)
	t.Run("list menus without soft deleted", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"site_id": cmsstore.ShortenID(site.ID()),
					"limit":   10,
					"offset":  0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should not return the soft deleted menu
		assert.Equal(t, 0, len(items), "Expected no menus (soft deleted should be excluded)")
	})

	// Test listing with include_soft_deleted (should include soft deleted)
	t.Run("list menus with soft deleted", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"site_id":              cmsstore.ShortenID(site.ID()),
					"include_soft_deleted": true,
					"limit":                10,
					"offset":               0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return the soft deleted menu
		assert.Equal(t, 1, len(items), "Expected 1 menu (soft deleted should be included)")
	})
}

func TestMenuList_WithOrdering(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create menus with different names to test ordering
	menu1 := cmsstore.NewMenu()
	menu1.SetName("Alpha Menu")
	menu1.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu1.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu1)
	require.NoError(t, err)

	menu2 := cmsstore.NewMenu()
	menu2.SetName("Beta Menu")
	menu2.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu2.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu2)
	require.NoError(t, err)

	// Test ordering by name ascending
	t.Run("list menus ordered by name ascending", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"site_id":    cmsstore.ShortenID(site.ID()),
					"order_by":   "name",
					"sort_order": "asc",
					"limit":      10,
					"offset":     0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return menus in alphabetical order
		assert.Equal(t, 2, len(items), "Expected 2 menus")
		assert.Equal(t, "Alpha Menu", items[0].(map[string]interface{})["name"].(string))
		assert.Equal(t, "Beta Menu", items[1].(map[string]interface{})["name"].(string))
	})

	// Test ordering by name descending
	t.Run("list menus ordered by name descending", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_list",
				"arguments": map[string]any{
					"site_id":    cmsstore.ShortenID(site.ID()),
					"order_by":   "name",
					"sort_order": "desc",
					"limit":      10,
					"offset":     0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		require.NoError(t, err)

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		require.NoError(t, err)
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		require.NoError(t, err)

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		require.NoError(t, err)

		result, ok := response["result"].(map[string]any)
		require.True(t, ok, "Expected response to have result")

		content, ok := result["content"].([]any)
		require.True(t, ok, "Expected response result.content")
		require.Len(t, content, 1, "Expected response result.content to have one item")

		item0, ok := content[0].(map[string]any)
		require.True(t, ok, "Expected response result.content[0] object")

		text, ok := item0["text"].(string)
		require.True(t, ok, "Expected response result.content[0].text")

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		require.NoError(t, err)

		items, ok := menuList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return menus in reverse alphabetical order
		assert.Equal(t, 2, len(items), "Expected 2 menus")
		assert.Equal(t, "Beta Menu", items[0].(map[string]interface{})["name"].(string))
		assert.Equal(t, "Alpha Menu", items[1].(map[string]interface{})["name"].(string))
	})
}
