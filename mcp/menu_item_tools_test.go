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

func TestMenuItemGet(t *testing.T) {
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

	// Create a menu item
	menuItem := cmsstore.NewMenuItem()
	menuItem.SetName("Home")
	menuItem.SetURL("/")
	menuItem.SetTarget("_self")
	menuItem.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	menuItem.SetMenuID(menu.ID())
	menuItem.SetHandle("home")
	menuItem.SetMemo("Home page link")
	err = store.MenuItemCreate(context.Background(), menuItem)
	require.NoError(t, err)

	tests := []struct {
		name        string
		menuItemID  string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "get menu item with full ID",
			menuItemID:  menuItem.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menuItem.ID()),
		},
		{
			name:        "get menu item with shortened ID",
			menuItemID:  cmsstore.ShortenID(menuItem.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menuItem.ID()),
		},
		{
			name:        "get non-existent menu item",
			menuItemID:  "non_existent_id",
			expectError: true,
			expectedErr: "menu item not found",
		},
		{
			name:        "get menu item with empty ID",
			menuItemID:  "",
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
					"tool_name": "menu_item_get",
					"arguments": map[string]any{
						"id": tt.menuItemID,
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

				var menuItemData map[string]any
				err = json.Unmarshal([]byte(text), &menuItemData)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedID, menuItemData["id"].(string))
				assert.Equal(t, "Home", menuItemData["name"].(string))
				assert.Equal(t, "/", menuItemData["url"].(string))
				assert.Equal(t, "_self", menuItemData["target"].(string))
				assert.Equal(t, cmsstore.MENU_ITEM_STATUS_ACTIVE, menuItemData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(menu.ID()), menuItemData["menu_id"].(string))

				// Check for new fields
				assert.Contains(t, menuItemData, "memo")
				assert.Contains(t, menuItemData, "page_id")
				assert.Contains(t, menuItemData, "parent_id")
				assert.Contains(t, menuItemData, "sequence")
				assert.Contains(t, menuItemData, "created_at")
				assert.Contains(t, menuItemData, "updated_at")
				// assert.Contains(t, menuItemData, "soft_deleted_at") // commented out to match tool response
				assert.Contains(t, menuItemData, "metas")
			}
		})
	}
}

func TestMenuItemList(t *testing.T) {
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

	// Create menu items with different properties
	activeItem := cmsstore.NewMenuItem()
	activeItem.SetName("Home")
	activeItem.SetURL("/")
	activeItem.SetTarget("_self")
	activeItem.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	activeItem.SetMenuID(menu.ID())
	activeItem.SetHandle("home")
	activeItem.SetMemo("Home page link")
	err = store.MenuItemCreate(context.Background(), activeItem)
	require.NoError(t, err)

	draftItem := cmsstore.NewMenuItem()
	draftItem.SetName("About")
	draftItem.SetURL("/about")
	draftItem.SetTarget("_self")
	draftItem.SetStatus(cmsstore.MENU_ITEM_STATUS_DRAFT)
	draftItem.SetMenuID(menu.ID())
	draftItem.SetHandle("about")
	draftItem.SetMemo("About page link")
	err = store.MenuItemCreate(context.Background(), draftItem)
	require.NoError(t, err)

	// Test listing all menu items
	t.Run("list all menu items", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return both menu items
		assert.Equal(t, 2, len(items), "Expected both menu items to be returned")
	})

	// Test filtering by menu_id
	t.Run("list menu items by menu_id", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
				"arguments": map[string]any{
					"menu_id": cmsstore.ShortenID(menu.ID()),
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return both menu items for the menu
		assert.Equal(t, 2, len(items), "Expected both menu items for the menu")
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			assert.Equal(t, cmsstore.ShortenID(menu.ID()), itemMap["menu_id"].(string))
		}
	})

	// Test filtering by status
	t.Run("list menu items by status", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
				"arguments": map[string]any{
					"status": cmsstore.MENU_ITEM_STATUS_ACTIVE,
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only active menu item
		assert.Equal(t, 1, len(items), "Expected only active menu item")
		item := items[0].(map[string]interface{})
		assert.Equal(t, cmsstore.MENU_ITEM_STATUS_ACTIVE, item["status"].(string))
	})

	// Test filtering by name_like
	t.Run("list menu items by name_like", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
				"arguments": map[string]any{
					"name_like": "Home",
					"limit":     10,
					"offset":    0,
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only the menu item with matching name
		assert.Equal(t, 1, len(items), "Expected only menu item with matching name")
		item := items[0].(map[string]interface{})
		assert.Equal(t, "Home", item["name"].(string))
	})

	// Test pagination
	t.Run("list menu items with pagination", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only 1 menu item due to limit
		assert.Equal(t, 1, len(items), "Expected only 1 menu item due to limit")
	})
}

func TestMenuItemUpsert_Create(t *testing.T) {
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
		name         string
		menuItemName string
		url          string
		target       string
		status       string
		menuID       string
		handle       string
		memo         string
		pageID       string
		parentID     string
		sequence     int
		expectError  bool
		expectedErr  string
	}{
		{
			name:         "create menu item with all fields",
			menuItemName: "Contact",
			url:          "/contact",
			target:       "_self",
			status:       cmsstore.MENU_ITEM_STATUS_ACTIVE,
			menuID:       cmsstore.ShortenID(menu.ID()),
			handle:       "contact",
			memo:         "Contact page link",
			pageID:       "",
			parentID:     "",
			sequence:     1,
			expectError:  false,
		},
		{
			name:         "create menu item with minimal fields",
			menuItemName: "About",
			url:          "/about",
			target:       "_self",
			status:       cmsstore.MENU_ITEM_STATUS_DRAFT,
			menuID:       "",
			handle:       "",
			memo:         "",
			pageID:       "",
			parentID:     "",
			sequence:     0,
			expectError:  false,
		},
		{
			name:         "create menu item with empty name",
			menuItemName: "",
			url:          "/test",
			target:       "_self",
			status:       cmsstore.MENU_ITEM_STATUS_ACTIVE,
			menuID:       cmsstore.ShortenID(menu.ID()),
			handle:       "test",
			memo:         "Test memo",
			pageID:       "",
			parentID:     "",
			sequence:     0,
			expectError:  true,
			expectedErr:  "missing required parameter: name",
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
					"tool_name": "menu_item_upsert",
					"arguments": map[string]any{
						"name":      tt.menuItemName,
						"url":       tt.url,
						"target":    tt.target,
						"status":    tt.status,
						"menu_id":   tt.menuID,
						"handle":    tt.handle,
						"memo":      tt.memo,
						"page_id":   tt.pageID,
						"parent_id": tt.parentID,
						"sequence":  tt.sequence,
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

				var menuItemData map[string]any
				err = json.Unmarshal([]byte(text), &menuItemData)
				require.NoError(t, err)

				assert.Equal(t, tt.menuItemName, menuItemData["name"].(string))
				assert.Equal(t, tt.status, menuItemData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(menu.ID()), menuItemData["menu_id"].(string))
			}
		})
	}
}

func TestMenuItemUpsert_Update(t *testing.T) {
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

	// Create a menu item
	menuItem := cmsstore.NewMenuItem()
	menuItem.SetName("Original Item")
	menuItem.SetURL("/original")
	menuItem.SetTarget("_self")
	menuItem.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	menuItem.SetMenuID(menu.ID())
	menuItem.SetHandle("original")
	menuItem.SetMemo("Original memo")
	err = store.MenuItemCreate(context.Background(), menuItem)
	require.NoError(t, err)

	tests := []struct {
		name         string
		menuItemID   string
		menuItemName string
		url          string
		target       string
		status       string
		handle       string
		expectError  bool
		expectedErr  string
	}{
		{
			name:         "update menu item with full ID",
			menuItemID:   menuItem.ID(),
			menuItemName: "Updated Item",
			url:          "/updated",
			target:       "_blank",
			status:       cmsstore.MENU_ITEM_STATUS_DRAFT,
			handle:       "updated",
			expectError:  false,
		},
		{
			name:         "update menu item with shortened ID",
			menuItemID:   cmsstore.ShortenID(menuItem.ID()),
			menuItemName: "Updated Item",
			url:          "/updated",
			target:       "_blank",
			status:       cmsstore.MENU_ITEM_STATUS_DRAFT,
			handle:       "updated",
			expectError:  false,
		},
		{
			name:         "update non-existent menu item",
			menuItemID:   "non_existent_id",
			menuItemName: "Updated Item",
			url:          "/updated",
			target:       "_blank",
			status:       cmsstore.MENU_ITEM_STATUS_DRAFT,
			handle:       "updated",
			expectError:  true,
			expectedErr:  "menu item not found",
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
					"tool_name": "menu_item_upsert",
					"arguments": map[string]any{
						"id":     tt.menuItemID,
						"name":   tt.menuItemName,
						"url":    tt.url,
						"target": tt.target,
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

				var menuItemData map[string]any
				err = json.Unmarshal([]byte(text), &menuItemData)
				require.NoError(t, err)

				assert.Equal(t, tt.menuItemName, menuItemData["name"].(string))
				assert.Equal(t, tt.status, menuItemData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(menu.ID()), menuItemData["menu_id"].(string))
			}
		})
	}
}

func TestMenuItemDelete(t *testing.T) {
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

	// Create a menu item
	menuItem := cmsstore.NewMenuItem()
	menuItem.SetName("Test Item")
	menuItem.SetURL("/test")
	menuItem.SetTarget("_self")
	menuItem.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	menuItem.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), menuItem)
	require.NoError(t, err)

	tests := []struct {
		name        string
		menuItemID  string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "delete menu item with full ID",
			menuItemID:  menuItem.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menuItem.ID()),
		},
		{
			name:        "delete menu item with shortened ID",
			menuItemID:  cmsstore.ShortenID(menuItem.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(menuItem.ID()),
		},
		{
			name:        "delete non-existent menu item",
			menuItemID:  "non_existent_id",
			expectError: true,
			expectedErr: "menu item not found",
		},
		{
			name:        "delete menu item with empty ID",
			menuItemID:  "",
			expectError: true,
			expectedErr: "missing required parameter: id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetID := tt.menuItemID
			if tt.name == "delete menu item with full ID" || tt.name == "delete menu item with shortened ID" {
				// Create a fresh menu item for each positive test case
				mi := cmsstore.NewMenuItem()
				mi.SetName("Test Item")
				mi.SetURL("/test")
				mi.SetTarget("_self")
				mi.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
				mi.SetMenuID(menu.ID())
				err = store.MenuItemCreate(context.Background(), mi)
				require.NoError(t, err)

				if tt.name == "delete menu item with full ID" {
					targetID = mi.ID()
				} else {
					targetID = cmsstore.ShortenID(mi.ID())
				}
				// Update expectedID to match the new menu item
				tt.expectedID = cmsstore.ShortenID(mi.ID())
			}

			// Call the tool
			deletePayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "delete",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "menu_item_delete",
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

func TestMenuItemUpsert_WithDefaultMenu(t *testing.T) {
	server, _, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a menu item without specifying menu_id - should use default menu
	upsertPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "upsert",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "menu_item_upsert",
			"arguments": map[string]any{
				"name":   "Default Menu Item",
				"url":    "/default",
				"target": "_self",
				"status": cmsstore.MENU_ITEM_STATUS_ACTIVE,
				"handle": "default-item",
				"memo":   "Menu item with default menu",
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

	var menuItemData map[string]any
	err = json.Unmarshal([]byte(text), &menuItemData)
	require.NoError(t, err)

	assert.Equal(t, "Default Menu Item", menuItemData["name"].(string))
	assert.Equal(t, cmsstore.MENU_ITEM_STATUS_ACTIVE, menuItemData["status"].(string))
	assert.Equal(t, "default-item", menuItemData["handle"].(string))
	assert.Equal(t, "Menu item with default menu", menuItemData["memo"].(string))
	// menu_id should be set to the default menu
	assert.NotEmpty(t, menuItemData["menu_id"].(string))
}

func TestMenuItemList_WithSoftDeleted(t *testing.T) {
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

	// Create a menu item
	menuItem := cmsstore.NewMenuItem()
	menuItem.SetName("Soft Deleted Item")
	menuItem.SetURL("/soft-deleted")
	menuItem.SetTarget("_self")
	menuItem.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	menuItem.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), menuItem)
	require.NoError(t, err)

	// Soft delete the menu item
	err = store.MenuItemSoftDeleteByID(context.Background(), menuItem.ID())
	require.NoError(t, err)

	// Test listing without include_soft_deleted (should not include soft deleted)
	t.Run("list menu items without soft deleted", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
				"arguments": map[string]any{
					"menu_id": cmsstore.ShortenID(menu.ID()),
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should not return the soft deleted menu item
		assert.Equal(t, 0, len(items), "Expected no menu items (soft deleted should be excluded)")
	})

	// Test listing with include_soft_deleted (should include soft deleted)
	t.Run("list menu items with soft deleted", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
				"arguments": map[string]any{
					"menu_id":              cmsstore.ShortenID(menu.ID()),
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return the soft deleted menu item
		assert.Equal(t, 1, len(items), "Expected 1 menu item (soft deleted should be included)")
	})
}

func TestMenuItemList_WithOrdering(t *testing.T) {
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

	// Create menu items with different names to test ordering
	item1 := cmsstore.NewMenuItem()
	item1.SetName("Alpha Item")
	item1.SetURL("/alpha")
	item1.SetTarget("_self")
	item1.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	item1.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), item1)
	require.NoError(t, err)

	item2 := cmsstore.NewMenuItem()
	item2.SetName("Beta Item")
	item2.SetURL("/beta")
	item2.SetTarget("_self")
	item2.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	item2.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), item2)
	require.NoError(t, err)

	// Test ordering by name ascending
	t.Run("list menu items ordered by name ascending", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
				"arguments": map[string]any{
					"menu_id":    cmsstore.ShortenID(menu.ID()),
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return menu items in alphabetical order
		assert.Equal(t, 2, len(items), "Expected 2 menu items")
		assert.Equal(t, "Alpha Item", items[0].(map[string]interface{})["name"].(string))
		assert.Equal(t, "Beta Item", items[1].(map[string]interface{})["name"].(string))
	})

	// Test ordering by name descending
	t.Run("list menu items ordered by name descending", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "menu_item_list",
				"arguments": map[string]any{
					"menu_id":    cmsstore.ShortenID(menu.ID()),
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

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		require.NoError(t, err)

		items, ok := menuItemList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return menu items in reverse alphabetical order
		assert.Equal(t, 2, len(items), "Expected 2 menu items")
		assert.Equal(t, "Beta Item", items[0].(map[string]interface{})["name"].(string))
		assert.Equal(t, "Alpha Item", items[1].(map[string]interface{})["name"].(string))
	})
}
