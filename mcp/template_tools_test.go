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

func TestTemplateList_SiteIDUnshortening(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site first
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a template with the site
	template := cmsstore.NewTemplate()
	template.SetName("Test Template")
	template.SetContent("Test content")
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetSiteID(site.ID())
	err = store.TemplateCreate(context.Background(), template)
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
					"tool_name": "template_list",
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

			var templateList map[string]any
			err = json.Unmarshal([]byte(text), &templateList)
			require.NoError(t, err)

			items, ok := templateList["items"].([]interface{})
			require.True(t, ok, "Expected 'items' to be a slice")

			assert.Equal(t, tt.expectedCount, len(items), "Unexpected number of templates")

			if tt.expectedCount > 0 {
				item := items[0].(map[string]interface{})
				assert.Equal(t, tt.expectedSiteID, item["site_id"].(string))
			}
		})
	}
}

func TestTemplateList_NoSiteID(t *testing.T) {
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

	// Create templates for both sites
	template1 := cmsstore.NewTemplate()
	template1.SetName("Template 1")
	template1.SetContent("Test content 1")
	template1.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template1.SetSiteID(site1.ID())
	err = store.TemplateCreate(context.Background(), template1)
	require.NoError(t, err)

	template2 := cmsstore.NewTemplate()
	template2.SetName("Template 2")
	template2.SetContent("Test content 2")
	template2.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template2.SetSiteID(site2.ID())
	err = store.TemplateCreate(context.Background(), template2)
	require.NoError(t, err)

	// Call the tool without site_id parameter
	listPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "list",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "template_list",
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

	var templateList map[string]any
	err = json.Unmarshal([]byte(text), &templateList)
	require.NoError(t, err)

	items, ok := templateList["items"].([]interface{})
	require.True(t, ok, "Expected 'items' to be a slice")

	// Should return all templates when no site_id is specified
	assert.Equal(t, 2, len(items), "Expected all templates to be returned")
}

func TestTemplateList_WithOtherFilters(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create templates with different statuses
	activeTemplate := cmsstore.NewTemplate()
	activeTemplate.SetName("Active Template")
	activeTemplate.SetContent("Active content")
	activeTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	activeTemplate.SetSiteID(site.ID())
	err = store.TemplateCreate(context.Background(), activeTemplate)
	require.NoError(t, err)

	draftTemplate := cmsstore.NewTemplate()
	draftTemplate.SetName("Draft Template")
	draftTemplate.SetContent("Draft content")
	draftTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_DRAFT)
	draftTemplate.SetSiteID(site.ID())
	err = store.TemplateCreate(context.Background(), draftTemplate)
	require.NoError(t, err)

	// Call the tool with multiple filters
	listPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "list",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "template_list",
			"arguments": map[string]any{
				"site_id": cmsstore.ShortenID(site.ID()), // shortened site ID
				"status":  cmsstore.TEMPLATE_STATUS_ACTIVE,
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

	var templateList map[string]any
	err = json.Unmarshal([]byte(text), &templateList)
	require.NoError(t, err)

	items, ok := templateList["items"].([]interface{})
	require.True(t, ok, "Expected 'items' to be a slice")

	// Should return only the active template for the specified site
	assert.Equal(t, 1, len(items), "Expected only active template")
	item := items[0].(map[string]interface{})
	assert.Equal(t, cmsstore.TEMPLATE_STATUS_ACTIVE, item["status"].(string))
}
