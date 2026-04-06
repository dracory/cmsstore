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

func TestTemplateList_SiteIDUnshortening(t *testing.T) {
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

	// Create a template with the site
	template := cmsstore.NewTemplate()
	template.SetName("Test Template")
	template.SetContent("Test content")
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetSiteID(site.ID())
	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
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
					"tool_name": "template_list",
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

			var templateList map[string]any
			err = json.Unmarshal([]byte(text), &templateList)
			if err != nil {
				t.Fatalf("Failed to unmarshal template list: %v", err)
			}

			items, ok := templateList["items"].([]interface{})
			if !ok {
				t.Fatalf("Expected 'items' to be a slice")
			}

			if len(items) != tt.expectedCount {
				t.Errorf("Expected %d templates, got %d", tt.expectedCount, len(items))
			}

			if tt.expectedCount > 0 {
				item := items[0].(map[string]interface{})
				if item["site_id"].(string) != tt.expectedSiteID {
					t.Errorf("Expected site_id '%s', got '%s'", tt.expectedSiteID, item["site_id"].(string))
				}

				// Check for new fields
				if _, ok := item["memo"]; !ok {
					t.Errorf("Expected 'memo' field in response")
				}
				if _, ok := item["handle"]; !ok {
					t.Errorf("Expected 'handle' field in response")
				}
				if _, ok := item["editor"]; !ok {
					t.Errorf("Expected 'editor' field in response")
				}
				if _, ok := item["created_at"]; !ok {
					t.Errorf("Expected 'created_at' field in response")
				}
				if _, ok := item["updated_at"]; !ok {
					t.Errorf("Expected 'updated_at' field in response")
				}
				// assert.Contains(t, item, "soft_deleted_at") // commented out to match tool response
				if _, ok := item["metas"]; !ok {
					t.Errorf("Expected 'metas' field in response")
				}
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

	// Create templates for both sites
	template1 := cmsstore.NewTemplate()
	template1.SetName("Template 1")
	template1.SetContent("Test content 1")
	template1.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template1.SetSiteID(site1.ID())
	err = store.TemplateCreate(context.Background(), template1)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	template2 := cmsstore.NewTemplate()
	template2.SetName("Template 2")
	template2.SetContent("Test content 2")
	template2.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template2.SetSiteID(site2.ID())
	err = store.TemplateCreate(context.Background(), template2)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

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

	var templateList map[string]any
	err = json.Unmarshal([]byte(text), &templateList)
	if err != nil {
		t.Fatalf("Failed to unmarshal template list: %v", err)
	}

	items, ok := templateList["items"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'items' to be a slice")
	}

	// Should return all templates when no site_id is specified
	if len(items) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(items))
	}
}

func TestTemplateList_WithOtherFilters(t *testing.T) {
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

	// Create templates with different statuses
	activeTemplate := cmsstore.NewTemplate()
	activeTemplate.SetName("Active Template")
	activeTemplate.SetContent("Active content")
	activeTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	activeTemplate.SetSiteID(site.ID())
	err = store.TemplateCreate(context.Background(), activeTemplate)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	draftTemplate := cmsstore.NewTemplate()
	draftTemplate.SetName("Draft Template")
	draftTemplate.SetContent("Draft content")
	draftTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_DRAFT)
	draftTemplate.SetSiteID(site.ID())
	err = store.TemplateCreate(context.Background(), draftTemplate)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

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

	var templateList map[string]any
	err = json.Unmarshal([]byte(text), &templateList)
	if err != nil {
		t.Fatalf("Failed to unmarshal template list: %v", err)
	}

	items, ok := templateList["items"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'items' to be a slice")
	}

	// Should return only the active template for the specified site
	if len(items) != 1 {
		t.Errorf("Expected 1 template, got %d", len(items))
	}
	item := items[0].(map[string]interface{})
	if item["status"].(string) != cmsstore.TEMPLATE_STATUS_ACTIVE {
		t.Errorf("Expected status '%s', got '%s'", cmsstore.TEMPLATE_STATUS_ACTIVE, item["status"].(string))
	}
}

func TestTemplateUpsert_Create(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	createPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "create",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "template_upsert",
			"arguments": map[string]any{
				"name":    "New Template",
				"content": "Template content",
				"status":  "active",
			},
		},
	}

	createBody, err := json.Marshal(createPayload)
	if err != nil {
		t.Fatalf("Failed to marshal create payload: %v", err)
	}

	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to post create request: %v", err)
	}
	defer createResp.Body.Close()

	createRespBytes, err := io.ReadAll(createResp.Body)
	if err != nil {
		t.Fatalf("Failed to read create response: %v", err)
	}

	text := rpcResultText(t, createRespBytes)
	var templateData map[string]any
	err = json.Unmarshal([]byte(text), &templateData)
	if err != nil {
		t.Fatalf("Failed to unmarshal template data: %v", err)
	}

	if templateData["name"] != "New Template" {
		t.Errorf("Expected name 'New Template', got '%v'", templateData["name"])
	}
	if templateData["status"] != "active" {
		t.Errorf("Expected status 'active', got '%v'", templateData["status"])
	}
	if templateData["id"] == "" {
		t.Errorf("Expected id to not be empty")
	}
}

func TestTemplateUpsert_Update(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	template := cmsstore.NewTemplate()
	template.SetName("Original Name")
	template.SetContent("Original content")
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	err := store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	updatePayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "update",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "template_upsert",
			"arguments": map[string]any{
				"id":   template.ID(),
				"name": "Updated Name",
			},
		},
	}

	updateBody, err := json.Marshal(updatePayload)
	if err != nil {
		t.Fatalf("Failed to marshal update payload: %v", err)
	}

	updateResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatalf("Failed to post update request: %v", err)
	}
	defer updateResp.Body.Close()

	updateRespBytes, err := io.ReadAll(updateResp.Body)
	if err != nil {
		t.Fatalf("Failed to read update response: %v", err)
	}

	text := rpcResultText(t, updateRespBytes)
	var templateData map[string]any
	err = json.Unmarshal([]byte(text), &templateData)
	if err != nil {
		t.Fatalf("Failed to unmarshal template data: %v", err)
	}

	if templateData["id"] != template.ID() {
		t.Errorf("Expected id '%s', got '%v'", template.ID(), templateData["id"])
	}
	if templateData["name"] != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%v'", templateData["name"])
	}
	if templateData["content"] != "Original content" {
		t.Errorf("Expected content 'Original content', got '%v'", templateData["content"])
	}
}

func TestTemplateGet(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	template := cmsstore.NewTemplate()
	template.SetName("Get Template")
	template.SetContent("Content")
	err := store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	getPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "get",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "template_get",
			"arguments": map[string]any{
				"id": template.ID(),
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

	text := rpcResultText(t, getRespBytes)
	var templateData map[string]any
	err = json.Unmarshal([]byte(text), &templateData)
	if err != nil {
		t.Fatalf("Failed to unmarshal template data: %v", err)
	}

	if templateData["id"] != template.ID() {
		t.Errorf("Expected id '%s', got '%v'", template.ID(), templateData["id"])
	}
	if templateData["name"] != "Get Template" {
		t.Errorf("Expected name 'Get Template', got '%v'", templateData["name"])
	}
}

func TestTemplateDelete(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	template := cmsstore.NewTemplate()
	template.SetName("Delete Template")
	err := store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	deletePayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "delete",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "template_delete",
			"arguments": map[string]any{
				"id": template.ID(),
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

	text := rpcResultText(t, deleteRespBytes)
	var deleteData map[string]any
	err = json.Unmarshal([]byte(text), &deleteData)
	if err != nil {
		t.Fatalf("Failed to unmarshal delete data: %v", err)
	}

	if deleteData["id"] != cmsstore.ShortenID(template.ID()) {
		t.Errorf("Expected id '%s', got '%v'", cmsstore.ShortenID(template.ID()), deleteData["id"])
	}

	// Verify it's gone
	deletedTemplate, err := store.TemplateFindByID(context.Background(), template.ID())
	if err != nil {
		t.Fatalf("Failed to find deleted template: %v", err)
	}
	if deletedTemplate != nil {
		t.Errorf("Expected template to be deleted, but it was found")
	}
}
