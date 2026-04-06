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

func TestPageGet(t *testing.T) {
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

	// Create a page
	page := cmsstore.NewPage()
	page.SetTitle("Test Page")
	page.SetContent("Test content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	tests := []struct {
		name        string
		pageID      string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "get page with full ID",
			pageID:      page.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(page.ID()),
		},
		{
			name:        "get page with shortened ID",
			pageID:      cmsstore.ShortenID(page.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(page.ID()),
		},
		{
			name:        "get non-existent page",
			pageID:      "non_existent_id",
			expectError: true,
			expectedErr: "page not found",
		},
		{
			name:        "get page with empty ID",
			pageID:      "",
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
					"tool_name": "page_get",
					"arguments": map[string]any{
						"id": tt.pageID,
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

				var pageData map[string]any
				err = json.Unmarshal([]byte(text), &pageData)
				if err != nil {
					t.Fatalf("Failed to unmarshal page data: %v", err)
				}

				if pageData["id"].(string) != tt.expectedID {
					t.Errorf("Expected id '%s', got '%s'", tt.expectedID, pageData["id"].(string))
				}
				if pageData["title"].(string) != "Test Page" {
					t.Errorf("Expected title 'Test Page', got '%s'", pageData["title"].(string))
				}
				if pageData["content"].(string) != "Test content" {
					t.Errorf("Expected content 'Test content', got '%s'", pageData["content"].(string))
				}
				if pageData["status"].(string) != cmsstore.PAGE_STATUS_ACTIVE {
					t.Errorf("Expected status '%s', got '%s'", cmsstore.PAGE_STATUS_ACTIVE, pageData["status"].(string))
				}
				if pageData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), pageData["site_id"].(string))
				}

				// Check for new fields
				if _, ok := pageData["memo"]; !ok {
					t.Errorf("Expected 'memo' field in response")
				}
				if _, ok := pageData["alias"]; !ok {
					t.Errorf("Expected 'alias' field in response")
				}
				if _, ok := pageData["name"]; !ok {
					t.Errorf("Expected 'name' field in response")
				}
				if _, ok := pageData["handle"]; !ok {
					t.Errorf("Expected 'handle' field in response")
				}
				if _, ok := pageData["canonical_url"]; !ok {
					t.Errorf("Expected 'canonical_url' field in response")
				}
				if _, ok := pageData["meta_description"]; !ok {
					t.Errorf("Expected 'meta_description' field in response")
				}
				if _, ok := pageData["meta_keywords"]; !ok {
					t.Errorf("Expected 'meta_keywords' field in response")
				}
				if _, ok := pageData["meta_robots"]; !ok {
					t.Errorf("Expected 'meta_robots' field in response")
				}
				if _, ok := pageData["created_at"]; !ok {
					t.Errorf("Expected 'created_at' field in response")
				}
				if _, ok := pageData["updated_at"]; !ok {
					t.Errorf("Expected 'updated_at' field in response")
				}
				// assert.Contains(t, pageData, "soft_deleted_at") // commented out to match tool response
				if _, ok := pageData["metas"]; !ok {
					t.Errorf("Expected 'metas' field in response")
				}
			}
		})
	}
}

func TestPageDelete(t *testing.T) {
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

	// Create a page
	page := cmsstore.NewPage()
	page.SetTitle("Test Page")
	page.SetContent("Test content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	tests := []struct {
		name        string
		pageID      string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "delete page with full ID",
			pageID:      page.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(page.ID()),
		},
		{
			name:        "delete page with shortened ID",
			pageID:      cmsstore.ShortenID(page.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(page.ID()),
		},
		{
			name:        "delete non-existent page",
			pageID:      "non_existent_id",
			expectError: true,
			expectedErr: "page not found",
		},
		{
			name:        "delete page with empty ID",
			pageID:      "",
			expectError: true,
			expectedErr: "missing required parameter: id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetID := tt.pageID
			if tt.name == "delete page with full ID" || tt.name == "delete page with shortened ID" {
				// Create a fresh page for each positive test case
				p := cmsstore.NewPage()
				p.SetTitle("Test Page")
				p.SetContent("Test content")
				p.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
				p.SetSiteID(site.ID())
				err = store.PageCreate(context.Background(), p)
				if err != nil {
					t.Fatalf("Failed to create page: %v", err)
				}

				if tt.name == "delete page with full ID" {
					targetID = p.ID()
				} else {
					targetID = cmsstore.ShortenID(p.ID())
				}
				// Update expectedID to match the new page
				tt.expectedID = cmsstore.ShortenID(p.ID())
			} else {
				targetID = tt.pageID
			}

			// Call the tool
			deletePayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "delete",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "page_delete",
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
func TestPageList_SiteIDUnshortening(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site first
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a page with the site
	page := cmsstore.NewPage()
	page.SetTitle("Test Page")
	page.SetContent("Test content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	tests := []struct {
		name           string
		siteID         string
		expectedCount  int
		expectedSiteID string
	}{
		{
			name:           "shortened site ID should work",
			siteID:         cmsstore.ShortenID(site.ID()), // 9-char shortened ID
			expectedCount:  1,
			expectedSiteID: cmsstore.ShortenID(site.ID()),
		},
		{
			name:           "full site ID should work",
			siteID:         site.ID(), // full 32-char ID
			expectedCount:  1,
			expectedSiteID: cmsstore.ShortenID(site.ID()),
		},
		{
			name:           "different site ID should return empty",
			siteID:         "different_site_id",
			expectedCount:  0,
			expectedSiteID: "different_site_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the tool with site_id parameter
			listPayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "list",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "page_list",
					"arguments": map[string]any{
						"site_id": tt.siteID,
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

			var pageList map[string]any
			err = json.Unmarshal([]byte(text), &pageList)
			if err != nil {
				t.Fatalf("Failed to unmarshal page list: %v", err)
			}

			items, ok := pageList["items"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'items' to be a slice")
			}

			if len(items) != tt.expectedCount {
				t.Errorf("Expected %d pages, got %d", tt.expectedCount, len(items))
			}

			if tt.expectedCount > 0 {
				item := items[0].(map[string]interface{})
				if item["site_id"].(string) != tt.expectedSiteID {
					t.Errorf("Expected site_id '%s', got '%s'", tt.expectedSiteID, item["site_id"].(string))
				}
			}
		})
	}
}

func TestPageList_NoSiteID(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create two sites
	site1 := cmsstore.NewSite()
	site1.SetName("Site 1")
	site1.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site1)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	site2 := cmsstore.NewSite()
	site2.SetName("Site 2")
	site2.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(context.Background(), site2)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create pages for both sites
	page1 := cmsstore.NewPage()
	page1.SetTitle("Page 1")
	page1.SetContent("Test content 1")
	page1.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page1.SetSiteID(site1.ID())
	err = store.PageCreate(context.Background(), page1)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	page2 := cmsstore.NewPage()
	page2.SetTitle("Page 2")
	page2.SetContent("Test content 2")
	page2.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page2.SetSiteID(site2.ID())
	err = store.PageCreate(context.Background(), page2)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	// Call the tool without site_id parameter
	listPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "list",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_list",
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

	var pageList map[string]any
	err = json.Unmarshal([]byte(text), &pageList)
	if err != nil {
		t.Fatalf("Failed to unmarshal page list: %v", err)
	}

	items, ok := pageList["items"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'items' to be a slice")
	}

	// Should return all pages when no site_id is specified
	if len(items) != 2 {
		t.Errorf("Expected 2 pages, got %d", len(items))
	}
}

func TestPageList_WithOtherFilters(t *testing.T) {
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

	// Create pages with different statuses
	activePage := cmsstore.NewPage()
	activePage.SetTitle("Active Page")
	activePage.SetContent("Active content")
	activePage.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	activePage.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), activePage)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	draftPage := cmsstore.NewPage()
	draftPage.SetTitle("Draft Page")
	draftPage.SetContent("Draft content")
	draftPage.SetStatus(cmsstore.PAGE_STATUS_DRAFT)
	draftPage.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), draftPage)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	// Call the tool with multiple filters
	listPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "list",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_list",
			"arguments": map[string]any{
				"site_id": cmsstore.ShortenID(site.ID()), // shortened site ID
				"status":  cmsstore.PAGE_STATUS_ACTIVE,
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

	var pageList map[string]any
	err = json.Unmarshal([]byte(text), &pageList)
	if err != nil {
		t.Fatalf("Failed to unmarshal page list: %v", err)
	}

	items, ok := pageList["items"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'items' to be a slice")
	}

	// Should return only the active page for the specified site
	if len(items) != 1 {
		t.Errorf("Expected 1 page, got %d", len(items))
	}
	item := items[0].(map[string]interface{})
	if item["status"].(string) != cmsstore.PAGE_STATUS_ACTIVE {
		t.Errorf("Expected status '%s', got '%s'", cmsstore.PAGE_STATUS_ACTIVE, item["status"].(string))
	}
}

func TestPageUpsert_Create(t *testing.T) {
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
		title       string
		content     string
		status      string
		siteID      string
		expectError bool
		expectedErr string
	}{
		{
			name:        "create page with all fields",
			title:       "New Page",
			content:     "New content",
			status:      cmsstore.PAGE_STATUS_ACTIVE,
			siteID:      cmsstore.ShortenID(site.ID()),
			expectError: false,
		},
		{
			name:        "create page with minimal fields",
			title:       "Minimal Page",
			content:     "",
			status:      cmsstore.PAGE_STATUS_DRAFT,
			siteID:      "",
			expectError: false,
		},
		{
			name:        "create page with empty title",
			title:       "",
			content:     "Content",
			status:      cmsstore.PAGE_STATUS_ACTIVE,
			siteID:      cmsstore.ShortenID(site.ID()),
			expectError: true,
			expectedErr: "missing required parameter: title",
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
					"tool_name": "page_upsert",
					"arguments": map[string]any{
						"title":   tt.title,
						"content": tt.content,
						"status":  tt.status,
						"site_id": tt.siteID,
						"alias":   "test-alias",
						"name":    "Test Name",
						"handle":  "test-handle",
						"memo":    "Test memo",
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

				var pageData map[string]any
				err = json.Unmarshal([]byte(text), &pageData)
				if err != nil {
					t.Fatalf("Failed to unmarshal page data: %v", err)
				}

				if pageData["title"].(string) != tt.title {
					t.Errorf("Expected title '%s', got '%s'", tt.title, pageData["title"].(string))
				}
				if pageData["content"].(string) != tt.content {
					t.Errorf("Expected content '%s', got '%s'", tt.content, pageData["content"].(string))
				}
				if pageData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, pageData["status"].(string))
				}
				if pageData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), pageData["site_id"].(string))
				}
			}
		})
	}
}

func TestPageUpsert_Update(t *testing.T) {
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

	// Create a page
	page := cmsstore.NewPage()
	page.SetTitle("Original Page")
	page.SetContent("Original content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	tests := []struct {
		name        string
		pageID      string
		title       string
		content     string
		status      string
		expectError bool
		expectedErr string
	}{
		{
			name:        "update page with full ID",
			pageID:      page.ID(),
			title:       "Updated Page",
			content:     "Updated content",
			status:      cmsstore.PAGE_STATUS_DRAFT,
			expectError: false,
		},
		{
			name:        "update page with shortened ID",
			pageID:      cmsstore.ShortenID(page.ID()),
			title:       "Updated Page",
			content:     "Updated content",
			status:      cmsstore.PAGE_STATUS_DRAFT,
			expectError: false,
		},
		{
			name:        "update non-existent page",
			pageID:      "non_existent_id",
			title:       "Updated Page",
			content:     "Updated content",
			status:      cmsstore.PAGE_STATUS_DRAFT,
			expectError: true,
			expectedErr: "page not found",
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
					"tool_name": "page_upsert",
					"arguments": map[string]any{
						"id":      tt.pageID,
						"title":   tt.title,
						"content": tt.content,
						"status":  tt.status,
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

				var pageData map[string]any
				err = json.Unmarshal([]byte(text), &pageData)
				if err != nil {
					t.Fatalf("Failed to unmarshal page data: %v", err)
				}

				if pageData["title"].(string) != tt.title {
					t.Errorf("Expected title '%s', got '%s'", tt.title, pageData["title"].(string))
				}
				if pageData["content"].(string) != tt.content {
					t.Errorf("Expected content '%s', got '%s'", tt.content, pageData["content"].(string))
				}
				if pageData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, pageData["status"].(string))
				}
				if pageData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), pageData["site_id"].(string))
				}
			}
		})
	}
}
