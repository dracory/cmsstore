package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

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
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

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
			if err != nil {
				t.Fatalf("Failed to marshal get payload: %v", err)
			}

			getResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(getBody))
			if err != nil {
				t.Fatalf("Failed to post get request: %v", err)
			}
			defer getResp.Body.Close()

			getRespBytes, err := io.ReadAll(getResp.Body)
			if err != nil {
				t.Fatalf("Failed to read get response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(getRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal get response: %v", err)
			}

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				if !hasError {
					t.Errorf("Expected error in response")
				}
				if hasError {
					errorObj := response["error"].(map[string]any)
					if errorObj["message"] != tt.expectedErr {
						t.Errorf("Expected error message '%s', got '%s'", tt.expectedErr, errorObj["message"])
					}
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				if !ok {
					t.Fatalf("Expected response to have result")
				}

				content, ok := result["content"].([]any)
				if !ok {
					t.Fatalf("Expected response result.content")
				}
				if len(content) != 1 {
					t.Fatalf("Expected response result.content to have one item")
				}

				item0, ok := content[0].(map[string]any)
				if !ok {
					t.Fatalf("Expected response result.content[0] object")
				}

				text, ok := item0["text"].(string)
				if !ok {
					t.Fatalf("Expected response result.content[0].text")
				}

				var menuItemData map[string]any
				err = json.Unmarshal([]byte(text), &menuItemData)
				if err != nil {
					t.Fatalf("Failed to unmarshal menu item data: %v", err)
				}

				if menuItemData["id"].(string) != tt.expectedID {
					t.Errorf("Expected id '%s', got '%s'", tt.expectedID, menuItemData["id"].(string))
				}
				if menuItemData["name"].(string) != "Home" {
					t.Errorf("Expected name 'Home', got '%s'", menuItemData["name"].(string))
				}
				if menuItemData["url"].(string) != "/" {
					t.Errorf("Expected url '/', got '%s'", menuItemData["url"].(string))
				}
				if menuItemData["target"].(string) != "_self" {
					t.Errorf("Expected target '_self', got '%s'", menuItemData["target"].(string))
				}
				if menuItemData["status"].(string) != cmsstore.MENU_ITEM_STATUS_ACTIVE {
					t.Errorf("Expected status '%s', got '%s'", cmsstore.MENU_ITEM_STATUS_ACTIVE, menuItemData["status"].(string))
				}
				if menuItemData["menu_id"].(string) != cmsstore.ShortenID(menu.ID()) {
					t.Errorf("Expected menu_id '%s', got '%s'", cmsstore.ShortenID(menu.ID()), menuItemData["menu_id"].(string))
				}

				// Check for new fields
				if _, ok := menuItemData["memo"]; !ok {
					t.Errorf("Expected menuItemData to have 'memo' key")
				}
				if _, ok := menuItemData["page_id"]; !ok {
					t.Errorf("Expected menuItemData to have 'page_id' key")
				}
				if _, ok := menuItemData["parent_id"]; !ok {
					t.Errorf("Expected menuItemData to have 'parent_id' key")
				}
				if _, ok := menuItemData["sequence"]; !ok {
					t.Errorf("Expected menuItemData to have 'sequence' key")
				}
				if _, ok := menuItemData["created_at"]; !ok {
					t.Errorf("Expected menuItemData to have 'created_at' key")
				}
				if _, ok := menuItemData["updated_at"]; !ok {
					t.Errorf("Expected menuItemData to have 'updated_at' key")
				}
				if _, ok := menuItemData["metas"]; !ok {
					t.Errorf("Expected menuItemData to have 'metas' key")
				}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

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
	if err != nil {
		t.Fatalf("Failed to create active menu item: %v", err)
	}

	draftItem := cmsstore.NewMenuItem()
	draftItem.SetName("About")
	draftItem.SetURL("/about")
	draftItem.SetTarget("_self")
	draftItem.SetStatus(cmsstore.MENU_ITEM_STATUS_DRAFT)
	draftItem.SetMenuID(menu.ID())
	draftItem.SetHandle("about")
	draftItem.SetMemo("About page link")
	err = store.MenuItemCreate(context.Background(), draftItem)
	if err != nil {
		t.Fatalf("Failed to create draft menu item: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return both menu items
		if len(items) != 2 {
			t.Errorf("Expected both menu items to be returned, got %d", len(items))
		}
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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return both menu items for the menu
		if len(items) != 2 {
			t.Errorf("Expected both menu items for the menu, got %d", len(items))
		}
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if itemMap["menu_id"].(string) != cmsstore.ShortenID(menu.ID()) {
				t.Errorf("Expected menu_id '%s', got '%s'", cmsstore.ShortenID(menu.ID()), itemMap["menu_id"].(string))
			}
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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only active menu item
		if len(items) != 1 {
			t.Errorf("Expected only active menu item, got %d", len(items))
		}
		item := items[0].(map[string]interface{})
		if item["status"].(string) != cmsstore.MENU_ITEM_STATUS_ACTIVE {
			t.Errorf("Expected status '%s', got '%s'", cmsstore.MENU_ITEM_STATUS_ACTIVE, item["status"].(string))
		}
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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only the menu item with matching name
		if len(items) != 1 {
			t.Errorf("Expected only menu item with matching name, got %d", len(items))
		}
		item := items[0].(map[string]interface{})
		if item["name"].(string) != "Home" {
			t.Errorf("Expected name 'Home', got '%s'", item["name"].(string))
		}
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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only 1 menu item due to limit
		if len(items) != 1 {
			t.Errorf("Expected only 1 menu item due to limit, got %d", len(items))
		}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

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
			if err != nil {
				t.Fatalf("Failed to marshal upsert payload: %v", err)
			}

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			if err != nil {
				t.Fatalf("Failed to post upsert request: %v", err)
			}
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			if err != nil {
				t.Fatalf("Failed to read upsert response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal upsert response: %v", err)
			}

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				if !hasError {
					t.Errorf("Expected error in response")
				}
				if hasError {
					errorObj := response["error"].(map[string]any)
					if errorObj["message"] != tt.expectedErr {
						t.Errorf("Expected error message '%s', got '%s'", tt.expectedErr, errorObj["message"])
					}
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				if !ok {
					t.Fatalf("Expected response to have result")
				}

				content, ok := result["content"].([]any)
				if !ok {
					t.Fatalf("Expected response result.content")
				}
				if len(content) != 1 {
					t.Fatalf("Expected response result.content to have one item")
				}

				item0, ok := content[0].(map[string]any)
				if !ok {
					t.Fatalf("Expected response result.content[0] object")
				}

				text, ok := item0["text"].(string)
				if !ok {
					t.Fatalf("Expected response result.content[0].text")
				}

				var menuItemData map[string]any
				err = json.Unmarshal([]byte(text), &menuItemData)
				if err != nil {
					t.Fatalf("Failed to unmarshal menu item data: %v", err)
				}

				if menuItemData["name"].(string) != tt.menuItemName {
					t.Errorf("Expected name '%s', got '%s'", tt.menuItemName, menuItemData["name"].(string))
				}
				if menuItemData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, menuItemData["status"].(string))
				}
				if menuItemData["menu_id"].(string) != cmsstore.ShortenID(menu.ID()) {
					t.Errorf("Expected menu_id '%s', got '%s'", cmsstore.ShortenID(menu.ID()), menuItemData["menu_id"].(string))
				}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

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
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

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
			if err != nil {
				t.Fatalf("Failed to marshal upsert payload: %v", err)
			}

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			if err != nil {
				t.Fatalf("Failed to post upsert request: %v", err)
			}
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			if err != nil {
				t.Fatalf("Failed to read upsert response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal upsert response: %v", err)
			}

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				if !hasError {
					t.Errorf("Expected error in response")
				}
				if hasError {
					errorObj := response["error"].(map[string]any)
					if errorObj["message"] != tt.expectedErr {
						t.Errorf("Expected error message '%s', got '%s'", tt.expectedErr, errorObj["message"])
					}
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				if !ok {
					t.Fatalf("Expected response to have result")
				}

				content, ok := result["content"].([]any)
				if !ok {
					t.Fatalf("Expected response result.content")
				}
				if len(content) != 1 {
					t.Fatalf("Expected response result.content to have one item")
				}

				item0, ok := content[0].(map[string]any)
				if !ok {
					t.Fatalf("Expected response result.content[0] object")
				}

				text, ok := item0["text"].(string)
				if !ok {
					t.Fatalf("Expected response result.content[0].text")
				}

				var menuItemData map[string]any
				err = json.Unmarshal([]byte(text), &menuItemData)
				if err != nil {
					t.Fatalf("Failed to unmarshal menu item data: %v", err)
				}

				if menuItemData["name"].(string) != tt.menuItemName {
					t.Errorf("Expected name '%s', got '%s'", tt.menuItemName, menuItemData["name"].(string))
				}
				if menuItemData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, menuItemData["status"].(string))
				}
				if menuItemData["menu_id"].(string) != cmsstore.ShortenID(menu.ID()) {
					t.Errorf("Expected menu_id '%s', got '%s'", cmsstore.ShortenID(menu.ID()), menuItemData["menu_id"].(string))
				}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	// Create a menu item
	menuItem := cmsstore.NewMenuItem()
	menuItem.SetName("Test Item")
	menuItem.SetURL("/test")
	menuItem.SetTarget("_self")
	menuItem.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	menuItem.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), menuItem)
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

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
				err := store.MenuItemCreate(context.Background(), mi)
				if err != nil {
					t.Fatalf("Failed to create menu item: %v", err)
				}

				if tt.name == "delete menu item with full ID" {
					targetID = mi.ID()
				} else {
					targetID = cmsstore.ShortenID(mi.ID())
				}
				// Update expectedID to match the new menu item
				tt.expectedID = cmsstore.ShortenID(mi.ID())
			} else {
				targetID = tt.menuItemID
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
			if err != nil {
				t.Fatalf("Failed to marshal delete payload: %v", err)
			}

			deleteResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(deleteBody))
			if err != nil {
				t.Fatalf("Failed to post delete request: %v", err)
			}
			defer deleteResp.Body.Close()

			deleteRespBytes, err := io.ReadAll(deleteResp.Body)
			if err != nil {
				t.Fatalf("Failed to read delete response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(deleteRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal delete response: %v", err)
			}

			if tt.expectError {
				// Check for error
				_, hasError := response["error"]
				if !hasError {
					t.Errorf("Expected error in response")
				}
				if hasError {
					errorObj := response["error"].(map[string]any)
					if errorObj["message"] != tt.expectedErr {
						t.Errorf("Expected error message '%s', got '%s'", tt.expectedErr, errorObj["message"])
					}
				}
			} else {
				// Check for success
				result, ok := response["result"].(map[string]any)
				if !ok {
					t.Fatalf("Expected response to have result")
				}

				content, ok := result["content"].([]any)
				if !ok {
					t.Fatalf("Expected response result.content")
				}
				if len(content) != 1 {
					t.Fatalf("Expected response result.content to have one item")
				}

				item0, ok := content[0].(map[string]any)
				if !ok {
					t.Fatalf("Expected response result.content[0] object")
				}

				text, ok := item0["text"].(string)
				if !ok {
					t.Fatalf("Expected response result.content[0].text")
				}

				var deleteData map[string]any
				err = json.Unmarshal([]byte(text), &deleteData)
				if err != nil {
					t.Fatalf("Failed to unmarshal delete data: %v", err)
				}

				if deleteData["id"].(string) != tt.expectedID {
					t.Errorf("Expected id '%s', got '%s'", tt.expectedID, deleteData["id"].(string))
				}
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
	if err != nil {
		t.Fatalf("Failed to marshal upsert payload: %v", err)
	}

	upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
	if err != nil {
		t.Fatalf("Failed to post upsert request: %v", err)
	}
	defer upsertResp.Body.Close()

	upsertRespBytes, err := io.ReadAll(upsertResp.Body)
	if err != nil {
		t.Fatalf("Failed to read upsert response: %v", err)
	}

	// Parse the result
	var response map[string]any
	err = json.Unmarshal(upsertRespBytes, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal upsert response: %v", err)
	}

	// Check for success
	result, ok := response["result"].(map[string]any)
	if !ok {
		t.Fatalf("Expected response to have result")
	}

	content, ok := result["content"].([]any)
	if !ok {
		t.Fatalf("Expected response result.content")
	}
	if len(content) != 1 {
		t.Fatalf("Expected response result.content to have one item")
	}

	item0, ok := content[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected response result.content[0] object")
	}

	text, ok := item0["text"].(string)
	if !ok {
		t.Fatalf("Expected response result.content[0].text")
	}

	var menuItemData map[string]any
	err = json.Unmarshal([]byte(text), &menuItemData)
	if err != nil {
		t.Fatalf("Failed to unmarshal menu item data: %v", err)
	}

	if menuItemData["name"].(string) != "Default Menu Item" {
		t.Errorf("Expected name 'Default Menu Item', got '%s'", menuItemData["name"].(string))
	}
	if menuItemData["status"].(string) != cmsstore.MENU_ITEM_STATUS_ACTIVE {
		t.Errorf("Expected status '%s', got '%s'", cmsstore.MENU_ITEM_STATUS_ACTIVE, menuItemData["status"].(string))
	}
	if menuItemData["handle"].(string) != "default-item" {
		t.Errorf("Expected handle 'default-item', got '%s'", menuItemData["handle"].(string))
	}
	if menuItemData["memo"].(string) != "Menu item with default menu" {
		t.Errorf("Expected memo 'Menu item with default menu', got '%s'", menuItemData["memo"].(string))
	}
	// menu_id should be set to the default menu
	if menuItemData["menu_id"].(string) == "" {
		t.Errorf("Expected menu_id to not be empty")
	}
}

func TestMenuItemList_WithSoftDeleted(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	// Create a menu item
	menuItem := cmsstore.NewMenuItem()
	menuItem.SetName("Soft Deleted Item")
	menuItem.SetURL("/soft-deleted")
	menuItem.SetTarget("_self")
	menuItem.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	menuItem.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), menuItem)
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

	// Soft delete the menu item
	err = store.MenuItemSoftDeleteByID(context.Background(), menuItem.ID())
	if err != nil {
		t.Fatalf("Failed to soft delete menu item: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should not return the soft deleted menu item
		if len(items) != 0 {
			t.Errorf("Expected no menu items (soft deleted should be excluded), got %d", len(items))
		}
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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return the soft deleted menu item
		if len(items) != 1 {
			t.Errorf("Expected 1 menu item (soft deleted should be included), got %d", len(items))
		}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Main Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("main-menu")
	menu.SetMemo("Test menu")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	// Create menu items with different names to test ordering
	item1 := cmsstore.NewMenuItem()
	item1.SetName("Alpha Item")
	item1.SetURL("/alpha")
	item1.SetTarget("_self")
	item1.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	item1.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), item1)
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

	item2 := cmsstore.NewMenuItem()
	item2.SetName("Beta Item")
	item2.SetURL("/beta")
	item2.SetTarget("_self")
	item2.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	item2.SetMenuID(menu.ID())
	err = store.MenuItemCreate(context.Background(), item2)
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return menu items in alphabetical order
		if len(items) != 2 {
			t.Errorf("Expected 2 menu items, got %d", len(items))
		}
		if items[0].(map[string]interface{})["name"].(string) != "Alpha Item" {
			t.Errorf("Expected first item name 'Alpha Item', got '%s'", items[0].(map[string]interface{})["name"].(string))
		}
		if items[1].(map[string]interface{})["name"].(string) != "Beta Item" {
			t.Errorf("Expected second item name 'Beta Item', got '%s'", items[1].(map[string]interface{})["name"].(string))
		}
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
		if err != nil {
			t.Fatalf("Failed to marshal list payload: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post list request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read list response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal list response: %v", err)
		}

		result, ok := response["result"].(map[string]any)
		if !ok {
			t.Fatalf("Expected response to have result")
		}

		content, ok := result["content"].([]any)
		if !ok {
			t.Fatalf("Expected response result.content")
		}
		if len(content) != 1 {
			t.Fatalf("Expected response result.content to have one item")
		}

		item0, ok := content[0].(map[string]any)
		if !ok {
			t.Fatalf("Expected response result.content[0] object")
		}

		text, ok := item0["text"].(string)
		if !ok {
			t.Fatalf("Expected response result.content[0].text")
		}

		var menuItemList map[string]any
		err = json.Unmarshal([]byte(text), &menuItemList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu item list: %v", err)
		}

		items, ok := menuItemList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return menu items in reverse alphabetical order
		if len(items) != 2 {
			t.Errorf("Expected 2 menu items, got %d", len(items))
		}
		if items[0].(map[string]interface{})["name"].(string) != "Beta Item" {
			t.Errorf("Expected first item name 'Beta Item', got '%s'", items[0].(map[string]interface{})["name"].(string))
		}
		if items[1].(map[string]interface{})["name"].(string) != "Alpha Item" {
			t.Errorf("Expected second item name 'Alpha Item', got '%s'", items[1].(map[string]interface{})["name"].(string))
		}
	})
}
