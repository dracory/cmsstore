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

func TestPageGet(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a page
	page := cmsstore.NewPage()
	page.SetTitle("Test Page")
	page.SetContent("Test content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)

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

				var pageData map[string]any
				err = json.Unmarshal([]byte(text), &pageData)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedID, pageData["id"].(string))
				assert.Equal(t, "Test Page", pageData["title"].(string))
				assert.Equal(t, "Test content", pageData["content"].(string))
				assert.Equal(t, cmsstore.PAGE_STATUS_ACTIVE, pageData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), pageData["site_id"].(string))
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
	require.NoError(t, err)

	// Create a page
	page := cmsstore.NewPage()
	page.SetTitle("Test Page")
	page.SetContent("Test content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)

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
				require.NoError(t, err)

				if tt.name == "delete page with full ID" {
					targetID = p.ID()
				} else {
					targetID = cmsstore.ShortenID(p.ID())
				}
				// Update expectedID to match the new page
				tt.expectedID = cmsstore.ShortenID(p.ID())
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

func TestPageList_SiteIDUnshortening(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site first
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a page with the site
	page := cmsstore.NewPage()
	page.SetTitle("Test Page")
	page.SetContent("Test content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)

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

			var pageList map[string]any
			err = json.Unmarshal([]byte(text), &pageList)
			require.NoError(t, err)

			items, ok := pageList["items"].([]interface{})
			require.True(t, ok, "Expected 'items' to be a slice")

			assert.Equal(t, tt.expectedCount, len(items), "Unexpected number of pages")

			if tt.expectedCount > 0 {
				item := items[0].(map[string]interface{})
				assert.Equal(t, tt.expectedSiteID, item["site_id"].(string))
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
	require.NoError(t, err)

	site2 := cmsstore.NewSite()
	site2.SetName("Site 2")
	site2.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(context.Background(), site2)
	require.NoError(t, err)

	// Create pages for both sites
	page1 := cmsstore.NewPage()
	page1.SetTitle("Page 1")
	page1.SetContent("Test content 1")
	page1.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page1.SetSiteID(site1.ID())
	err = store.PageCreate(context.Background(), page1)
	require.NoError(t, err)

	page2 := cmsstore.NewPage()
	page2.SetTitle("Page 2")
	page2.SetContent("Test content 2")
	page2.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page2.SetSiteID(site2.ID())
	err = store.PageCreate(context.Background(), page2)
	require.NoError(t, err)

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

	var pageList map[string]any
	err = json.Unmarshal([]byte(text), &pageList)
	require.NoError(t, err)

	items, ok := pageList["items"].([]interface{})
	require.True(t, ok, "Expected 'items' to be a slice")

	// Should return all pages when no site_id is specified
	assert.Equal(t, 2, len(items), "Expected all pages to be returned")
}

func TestPageList_WithOtherFilters(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create pages with different statuses
	activePage := cmsstore.NewPage()
	activePage.SetTitle("Active Page")
	activePage.SetContent("Active content")
	activePage.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	activePage.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), activePage)
	require.NoError(t, err)

	draftPage := cmsstore.NewPage()
	draftPage.SetTitle("Draft Page")
	draftPage.SetContent("Draft content")
	draftPage.SetStatus(cmsstore.PAGE_STATUS_DRAFT)
	draftPage.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), draftPage)
	require.NoError(t, err)

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

	var pageList map[string]any
	err = json.Unmarshal([]byte(text), &pageList)
	require.NoError(t, err)

	items, ok := pageList["items"].([]interface{})
	require.True(t, ok, "Expected 'items' to be a slice")

	// Should return only the active page for the specified site
	assert.Equal(t, 1, len(items), "Expected only active page")
	item := items[0].(map[string]interface{})
	assert.Equal(t, cmsstore.PAGE_STATUS_ACTIVE, item["status"].(string))
}

func TestPageUpsert_Create(t *testing.T) {
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

				var pageData map[string]any
				err = json.Unmarshal([]byte(text), &pageData)
				require.NoError(t, err)

				assert.Equal(t, tt.title, pageData["title"].(string))
				assert.Equal(t, tt.content, pageData["content"].(string))
				assert.Equal(t, tt.status, pageData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), pageData["site_id"].(string))
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
	require.NoError(t, err)

	// Create a page
	page := cmsstore.NewPage()
	page.SetTitle("Original Page")
	page.SetContent("Original content")
	page.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	page.SetSiteID(site.ID())
	err = store.PageCreate(context.Background(), page)
	require.NoError(t, err)

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

				var pageData map[string]any
				err = json.Unmarshal([]byte(text), &pageData)
				require.NoError(t, err)

				assert.Equal(t, tt.title, pageData["title"].(string))
				assert.Equal(t, tt.content, pageData["content"].(string))
				assert.Equal(t, tt.status, pageData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), pageData["site_id"].(string))
			}
		})
	}
}
