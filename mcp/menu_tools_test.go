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

func TestMenuGet(t *testing.T) {
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
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}

			getResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(getBody))
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}
			defer getResp.Body.Close()

			getRespBytes, err := io.ReadAll(getResp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(getRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
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

				var menuData map[string]any
				err = json.Unmarshal([]byte(text), &menuData)
				if err != nil {
					t.Fatalf("Failed to unmarshal menu data: %v", err)
				}

				if menuData["id"].(string) != tt.expectedID {
					t.Errorf("Expected id '%s', got '%s'", tt.expectedID, menuData["id"].(string))
				}
				if menuData["name"].(string) != "Main Menu" {
					t.Errorf("Expected name 'Main Menu', got '%s'", menuData["name"].(string))
				}
				if menuData["status"].(string) != cmsstore.MENU_STATUS_ACTIVE {
					t.Errorf("Expected status '%s', got '%s'", cmsstore.MENU_STATUS_ACTIVE, menuData["status"].(string))
				}
				if menuData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), menuData["site_id"].(string))
				}

				// Check for new fields
				if _, ok := menuData["memo"]; !ok {
					t.Errorf("Expected menuData to contain 'memo'")
				}
				if _, ok := menuData["created_at"]; !ok {
					t.Errorf("Expected menuData to contain 'created_at'")
				}
				if _, ok := menuData["updated_at"]; !ok {
					t.Errorf("Expected menuData to contain 'updated_at'")
				}
				// if _, ok := menuData["soft_deleted_at"]; !ok {
				// 	t.Errorf("Expected menuData to contain 'soft_deleted_at'")
				// }
				if _, ok := menuData["metas"]; !ok {
					t.Errorf("Expected menuData to contain 'metas'")
				}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create menus with different properties
	activeMenu := cmsstore.NewMenu()
	activeMenu.SetName("Main Menu")
	activeMenu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	activeMenu.SetSiteID(site.ID())
	activeMenu.SetHandle("main-menu")
	activeMenu.SetMemo("Active menu")
	err = store.MenuCreate(context.Background(), activeMenu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	draftMenu := cmsstore.NewMenu()
	draftMenu.SetName("Draft Menu")
	draftMenu.SetStatus(cmsstore.MENU_STATUS_DRAFT)
	draftMenu.SetSiteID(site.ID())
	draftMenu.SetHandle("draft-menu")
	draftMenu.SetMemo("Draft menu")
	err = store.MenuCreate(context.Background(), draftMenu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return both menus
		if len(items) != 2 {
			t.Errorf("Expected both menus to be returned, got %d", len(items))
		}
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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return both menus for the site
		if len(items) != 2 {
			t.Errorf("Expected both menus for the site, got %d", len(items))
		}
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if itemMap["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
				t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), itemMap["site_id"].(string))
			}
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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only active menu
		if len(items) != 1 {
			t.Errorf("Expected only active menu, got %d", len(items))
		}
		item := items[0].(map[string]interface{})
		if item["status"].(string) != cmsstore.MENU_STATUS_ACTIVE {
			t.Errorf("Expected status '%s', got '%s'", cmsstore.MENU_STATUS_ACTIVE, item["status"].(string))
		}
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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only the menu with matching handle
		if len(items) != 1 {
			t.Errorf("Expected only menu with matching handle, got %d", len(items))
		}
		item := items[0].(map[string]interface{})
		if item["handle"].(string) != "main-menu" {
			t.Errorf("Expected handle 'main-menu', got '%s'", item["handle"].(string))
		}
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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only 1 menu due to limit
		if len(items) != 1 {
			t.Errorf("Expected only 1 menu due to limit, got %d", len(items))
		}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

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
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
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

				var menuData map[string]any
				err = json.Unmarshal([]byte(text), &menuData)
				if err != nil {
					t.Fatalf("Failed to unmarshal menu data: %v", err)
				}

				if menuData["name"].(string) != tt.menuName {
					t.Errorf("Expected name '%s', got '%s'", tt.menuName, menuData["name"].(string))
				}
				if menuData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, menuData["status"].(string))
				}
				if menuData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), menuData["site_id"].(string))
				}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Original Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	menu.SetHandle("original-menu")
	menu.SetMemo("Original memo")
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

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
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
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

				var menuData map[string]any
				err = json.Unmarshal([]byte(text), &menuData)
				if err != nil {
					t.Fatalf("Failed to unmarshal menu data: %v", err)
				}

				if menuData["name"].(string) != tt.menuName {
					t.Errorf("Expected name '%s', got '%s'", tt.menuName, menuData["name"].(string))
				}
				if menuData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, menuData["status"].(string))
				}
				if menuData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), menuData["site_id"].(string))
				}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Test Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

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
				if err != nil {
					t.Fatalf("Failed to post request: %v", err)
				}

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
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}

			deleteResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(deleteBody))
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}
			defer deleteResp.Body.Close()

			deleteRespBytes, err := io.ReadAll(deleteResp.Body)
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(deleteRespBytes, &response)
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
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

	var menuData map[string]any
	err = json.Unmarshal([]byte(text), &menuData)
	if err != nil {
		t.Fatalf("Failed to unmarshal menu data: %v", err)
	}

	if menuData["name"].(string) != "Default Site Menu" {
		t.Errorf("Expected name 'Default Site Menu', got '%s'", menuData["name"].(string))
	}
	if menuData["status"].(string) != cmsstore.MENU_STATUS_ACTIVE {
		t.Errorf("Expected status '%s', got '%s'", cmsstore.MENU_STATUS_ACTIVE, menuData["status"].(string))
	}
	if menuData["handle"].(string) != "default-menu" {
		t.Errorf("Expected handle 'default-menu', got '%s'", menuData["handle"].(string))
	}
	if menuData["memo"].(string) != "Menu with default site" {
		t.Errorf("Expected memo 'Menu with default site', got '%s'", menuData["memo"].(string))
	}
	// site_id should be set to the default site
	if menuData["site_id"].(string) == "" {
		t.Errorf("Expected site_id to not be empty")
	}
}

func TestMenuList_WithSoftDeleted(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatal(err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Soft Deleted Menu")
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	// Soft delete the menu
	err = store.MenuSoftDeleteByID(context.Background(), menu.ID())
	if err != nil {
		t.Fatalf("Failed to soft delete menu: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should not return the soft deleted menu
		if len(items) != 0 {
			t.Errorf("Expected no menus (soft deleted should be excluded), got %d", len(items))
		}
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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return the soft deleted menu
		if len(items) != 1 {
			t.Errorf("Expected 1 menu (soft deleted should be included), got %d", len(items))
		}
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
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create menus with different names to test ordering
	menu1 := cmsstore.NewMenu()
	menu1.SetName("Alpha Menu")
	menu1.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu1.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu1)
	if err != nil {
		t.Fatalf("Failed to create menu1: %v", err)
	}

	menu2 := cmsstore.NewMenu()
	menu2.SetName("Beta Menu")
	menu2.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	menu2.SetSiteID(site.ID())
	err = store.MenuCreate(context.Background(), menu2)
	if err != nil {
		t.Fatalf("Failed to create menu2: %v", err)
	}

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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return menus in alphabetical order
		if len(items) != 2 {
			t.Errorf("Expected 2 menus, got %d", len(items))
		}
		if items[0].(map[string]interface{})["name"].(string) != "Alpha Menu" {
			t.Errorf("Expected first menu name 'Alpha Menu', got '%s'", items[0].(map[string]interface{})["name"].(string))
		}
		if items[1].(map[string]interface{})["name"].(string) != "Beta Menu" {
			t.Errorf("Expected second menu name 'Beta Menu', got '%s'", items[1].(map[string]interface{})["name"].(string))
		}
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
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}
		defer listResp.Body.Close()

		listRespBytes, err := io.ReadAll(listResp.Body)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
		}

		// Parse the result
		var response map[string]any
		err = json.Unmarshal(listRespBytes, &response)
		if err != nil {
			t.Fatalf("Failed to post request: %v", err)
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

		var menuList map[string]any
		err = json.Unmarshal([]byte(text), &menuList)
		if err != nil {
			t.Fatalf("Failed to unmarshal menu list: %v", err)
		}

		items, ok := menuList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return menus in reverse alphabetical order
		if len(items) != 2 {
			t.Errorf("Expected 2 menus, got %d", len(items))
		}
		if items[0].(map[string]interface{})["name"].(string) != "Beta Menu" {
			t.Errorf("Expected first menu name 'Beta Menu', got '%s'", items[0].(map[string]interface{})["name"].(string))
		}
		if items[1].(map[string]interface{})["name"].(string) != "Alpha Menu" {
			t.Errorf("Expected second menu name 'Alpha Menu', got '%s'", items[1].(map[string]interface{})["name"].(string))
		}
	})
}
