package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/spf13/cast"
	_ "modernc.org/sqlite"
)

func TestBlockGet(t *testing.T) {
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

	// Create a block
	block := cmsstore.NewBlock()
	block.SetType("text")
	block.SetContent("Test content")
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetSiteID(site.ID())
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	tests := []struct {
		name        string
		blockID     string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "get block with full ID",
			blockID:     block.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(block.ID()),
		},
		{
			name:        "get block with shortened ID",
			blockID:     cmsstore.ShortenID(block.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(block.ID()),
		},
		{
			name:        "get non-existent block",
			blockID:     "non_existent_id",
			expectError: true,
			expectedErr: "block not found",
		},
		{
			name:        "get block with empty ID",
			blockID:     "",
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
					"tool_name": "block_get",
					"arguments": map[string]any{
						"id": tt.blockID,
					},
				},
			}

			getBody, err := json.Marshal(getPayload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
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

				var blockData map[string]any
				err = json.Unmarshal([]byte(text), &blockData)
				if err != nil {
					t.Fatalf("Failed to unmarshal block data: %v", err)
				}

				if blockData["id"].(string) != tt.expectedID {
					t.Errorf("Expected id '%s', got '%s'", tt.expectedID, blockData["id"].(string))
				}
				if blockData["type"].(string) != "text" {
					t.Errorf("Expected type 'text', got '%s'", blockData["type"].(string))
				}
				if blockData["content"].(string) != "Test content" {
					t.Errorf("Expected content 'Test content', got '%s'", blockData["content"].(string))
				}
				if blockData["status"].(string) != cmsstore.BLOCK_STATUS_ACTIVE {
					t.Errorf("Expected status '%s', got '%s'", cmsstore.BLOCK_STATUS_ACTIVE, blockData["status"].(string))
				}
				if blockData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
				}
				if blockData["page_id"].(string) != "" {
					t.Errorf("Expected empty page_id, got '%s'", blockData["page_id"].(string))
				}
				// New fields assertions
				if blockData["name"].(string) != "" {
					t.Errorf("Expected empty name, got '%s'", blockData["name"].(string))
				}
				if blockData["handle"].(string) != "" {
					t.Errorf("Expected empty handle, got '%s'", blockData["handle"].(string))
				}
				if blockData["template_id"].(string) != "" {
					t.Errorf("Expected empty template_id, got '%s'", blockData["template_id"].(string))
				}
				if blockData["parent_id"].(string) != "" {
					t.Errorf("Expected empty parent_id, got '%s'", blockData["parent_id"].(string))
				}
				if blockData["sequence"].(string) != "0" {
					t.Errorf("Expected sequence '0', got '%s'", blockData["sequence"].(string))
				}
				if blockData["editor"].(string) != "" {
					t.Errorf("Expected empty editor, got '%s'", blockData["editor"].(string))
				}
				if blockData["memo"].(string) != "" {
					t.Errorf("Expected empty memo, got '%s'", blockData["memo"].(string))
				}
				if blockData["created_at"].(string) == "" {
					t.Errorf("Expected non-empty created_at")
				}
				if blockData["updated_at"].(string) == "" {
					t.Errorf("Expected non-empty updated_at")
				}
				if blockData["soft_deleted_at"].(string) != "9999-12-31 23:59:59" {
					t.Errorf("Expected soft_deleted_at '9999-12-31 23:59:59', got '%s'", blockData["soft_deleted_at"].(string))
				}
				if blockData["metas"] == nil {
					t.Errorf("Expected non-nil metas")
				}
			}
		})
	}
}

func TestBlockList(t *testing.T) {
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

	// Create blocks with different properties
	activeBlock := cmsstore.NewBlock()
	activeBlock.SetType("text")
	activeBlock.SetContent("Active content")
	activeBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	activeBlock.SetSiteID(site.ID())
	activeBlock.SetName("Active Block")
	activeBlock.SetHandle("active-block")
	err = store.BlockCreate(context.Background(), activeBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	draftBlock := cmsstore.NewBlock()
	draftBlock.SetType("image")
	draftBlock.SetContent("Draft content")
	draftBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	draftBlock.SetSiteID(site.ID())
	draftBlock.SetName("Draft Block")
	draftBlock.SetHandle("draft-block")
	err = store.BlockCreate(context.Background(), draftBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Test listing all blocks
	t.Run("list all blocks", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "block_list",
				"arguments": map[string]any{
					"limit":  10,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		if err != nil {
			t.Fatalf("Failed to unmarshal block list: %v", err)
		}

		items, ok := blockList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return both blocks
		if len(items) != 2 {
			t.Errorf("Expected both blocks to be returned, got %d", len(items))
		}
	})

	// Test filtering by site_id
	t.Run("list blocks by site_id", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "block_list",
				"arguments": map[string]any{
					"site_id": cmsstore.ShortenID(site.ID()),
					"limit":   10,
					"offset":  0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		if err != nil {
			t.Fatalf("Failed to unmarshal block list: %v", err)
		}

		items, ok := blockList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return both blocks for the site
		if len(items) != 2 {
			t.Errorf("Expected both blocks for the site, got %d", len(items))
		}
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if itemMap["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
				t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), itemMap["site_id"].(string))
			}
		}
	})

	// Test filtering by status
	t.Run("list blocks by status", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "block_list",
				"arguments": map[string]any{
					"status": cmsstore.BLOCK_STATUS_ACTIVE,
					"limit":  10,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		if err != nil {
			t.Fatalf("Failed to unmarshal block list: %v", err)
		}

		items, ok := blockList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only active block
		if len(items) != 1 {
			t.Errorf("Expected only active block, got %d", len(items))
		}
		item := items[0].(map[string]interface{})
		if item["status"].(string) != cmsstore.BLOCK_STATUS_ACTIVE {
			t.Errorf("Expected status '%s', got '%s'", cmsstore.BLOCK_STATUS_ACTIVE, item["status"].(string))
		}
	})

	// Test filtering by handle
	t.Run("list blocks by handle", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "block_list",
				"arguments": map[string]any{
					"handle": "active-block",
					"limit":  10,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		if err != nil {
			t.Fatalf("Failed to unmarshal block list: %v", err)
		}

		items, ok := blockList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only the block with matching handle
		if len(items) != 1 {
			t.Errorf("Expected only block with matching handle, got %d", len(items))
		}
		item := items[0].(map[string]interface{})
		if item["handle"].(string) != "active-block" {
			t.Errorf("Expected handle 'active-block', got '%s'", item["handle"].(string))
		}
	})

	// Test pagination
	t.Run("list blocks with pagination", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "block_list",
				"arguments": map[string]any{
					"limit":  1,
					"offset": 0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		if err != nil {
			t.Fatalf("Failed to unmarshal block list: %v", err)
		}

		items, ok := blockList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only 1 block due to limit
		if len(items) != 1 {
			t.Errorf("Expected only 1 block due to limit, got %d", len(items))
		}
	})
}

func TestBlockUpsert_Create(t *testing.T) {
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
		blockType   string
		content     string
		status      string
		siteID      string
		blockName   string
		handle      string
		editor      string
		memo        string
		sequence    int
		expectError bool
		expectedErr string
	}{
		{
			name:        "create block with all fields",
			blockType:   "text",
			content:     "New content",
			status:      cmsstore.BLOCK_STATUS_ACTIVE,
			siteID:      cmsstore.ShortenID(site.ID()),
			blockName:   "New Block",
			handle:      "new-block",
			editor:      "html",
			memo:        "Test memo",
			sequence:    1,
			expectError: false,
		},
		{
			name:        "create block with minimal fields",
			blockType:   "image",
			content:     "",
			status:      cmsstore.BLOCK_STATUS_DRAFT,
			siteID:      "",
			blockName:   "",
			handle:      "",
			editor:      "",
			memo:        "",
			sequence:    0,
			expectError: false,
		},
		{
			name:        "create block with empty type",
			blockType:   "",
			content:     "Content",
			status:      cmsstore.BLOCK_STATUS_ACTIVE,
			siteID:      cmsstore.ShortenID(site.ID()),
			blockName:   "Test Block",
			handle:      "test-block",
			editor:      "html",
			memo:        "Test memo",
			sequence:    1,
			expectError: true,
			expectedErr: "missing required parameter: type",
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
					"tool_name": "block_upsert",
					"arguments": map[string]any{
						"type":     tt.blockType,
						"content":  tt.content,
						"status":   tt.status,
						"site_id":  tt.siteID,
						"name":     tt.blockName,
						"handle":   tt.handle,
						"editor":   tt.editor,
						"memo":     tt.memo,
						"sequence": tt.sequence,
					},
				},
			}

			upsertBody, err := json.Marshal(upsertPayload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
			}

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
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

				var blockData map[string]any
				err = json.Unmarshal([]byte(text), &blockData)
				if err != nil {
					t.Fatalf("Failed to unmarshal block data: %v", err)
				}

				if blockData["type"].(string) != tt.blockType {
					t.Errorf("Expected type '%s', got '%s'", tt.blockType, blockData["type"].(string))
				}
				if blockData["content"].(string) != tt.content {
					t.Errorf("Expected content '%s', got '%s'", tt.content, blockData["content"].(string))
				}
				if blockData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, blockData["status"].(string))
				}
				if blockData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
				}
				// New fields assertions
				if blockData["name"].(string) != tt.blockName {
					t.Errorf("Expected name '%s', got '%s'", tt.blockName, blockData["name"].(string))
				}
				if blockData["handle"].(string) != tt.handle {
					t.Errorf("Expected handle '%s', got '%s'", tt.handle, blockData["handle"].(string))
				}
				if blockData["page_id"].(string) != "" {
					t.Errorf("Expected empty page_id, got '%s'", blockData["page_id"].(string))
				}
				if blockData["template_id"].(string) != "" {
					t.Errorf("Expected empty template_id, got '%s'", blockData["template_id"].(string))
				}
				if blockData["parent_id"].(string) != "" {
					t.Errorf("Expected empty parent_id, got '%s'", blockData["parent_id"].(string))
				}
				if blockData["sequence"].(string) != cast.ToString(tt.sequence) {
					t.Errorf("Expected sequence '%s', got '%s'", cast.ToString(tt.sequence), blockData["sequence"].(string))
				}
				if blockData["editor"].(string) != tt.editor {
					t.Errorf("Expected editor '%s', got '%s'", tt.editor, blockData["editor"].(string))
				}
				if blockData["memo"].(string) != tt.memo {
					t.Errorf("Expected memo '%s', got '%s'", tt.memo, blockData["memo"].(string))
				}
				if blockData["created_at"].(string) == "" {
					t.Errorf("Expected non-empty created_at")
				}
				if blockData["updated_at"].(string) == "" {
					t.Errorf("Expected non-empty updated_at")
				}
				if blockData["soft_deleted_at"].(string) != "9999-12-31 23:59:59" {
					t.Errorf("Expected soft_deleted_at '9999-12-31 23:59:59', got '%s'", blockData["soft_deleted_at"].(string))
				}
				if blockData["metas"] == nil {
					t.Errorf("Expected non-nil metas")
				}
			}
		})
	}
}

func TestBlockUpsert_Update(t *testing.T) {
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

	// Create a block
	block := cmsstore.NewBlock()
	block.SetType("text")
	block.SetContent("Original content")
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetSiteID(site.ID())
	block.SetName("Original Block")
	block.SetHandle("original-block")
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	tests := []struct {
		name        string
		blockID     string
		blockType   string
		content     string
		status      string
		blockName   string
		handle      string
		expectError bool
		expectedErr string
	}{
		{
			name:        "update block with full ID",
			blockID:     block.ID(),
			blockType:   "image",
			content:     "Updated content",
			status:      cmsstore.BLOCK_STATUS_DRAFT,
			blockName:   "Updated Block",
			handle:      "updated-block",
			expectError: false,
		},
		{
			name:        "update block with shortened ID",
			blockID:     cmsstore.ShortenID(block.ID()),
			blockType:   "image",
			content:     "Updated content",
			status:      cmsstore.BLOCK_STATUS_DRAFT,
			blockName:   "Updated Block",
			handle:      "updated-block",
			expectError: false,
		},
		{
			name:        "update non-existent block",
			blockID:     "non_existent_id",
			blockType:   "image",
			content:     "Updated content",
			status:      cmsstore.BLOCK_STATUS_DRAFT,
			blockName:   "Updated Block",
			handle:      "updated-block",
			expectError: true,
			expectedErr: "block not found",
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
					"tool_name": "block_upsert",
					"arguments": map[string]any{
						"id":      tt.blockID,
						"type":    tt.blockType,
						"content": tt.content,
						"status":  tt.status,
						"name":    tt.blockName,
						"handle":  tt.handle,
					},
				},
			}

			upsertBody, err := json.Marshal(upsertPayload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
			}

			upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}
			defer upsertResp.Body.Close()

			upsertRespBytes, err := io.ReadAll(upsertResp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(upsertRespBytes, &response)
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

				var blockData map[string]any
				err = json.Unmarshal([]byte(text), &blockData)
				if err != nil {
					t.Fatalf("Failed to unmarshal block data: %v", err)
				}

				if blockData["type"].(string) != tt.blockType {
					t.Errorf("Expected type '%s', got '%s'", tt.blockType, blockData["type"].(string))
				}
				if blockData["content"].(string) != tt.content {
					t.Errorf("Expected content '%s', got '%s'", tt.content, blockData["content"].(string))
				}
				if blockData["status"].(string) != tt.status {
					t.Errorf("Expected status '%s', got '%s'", tt.status, blockData["status"].(string))
				}
				if blockData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
				}
			}
		})
	}
}

func TestBlockDelete(t *testing.T) {
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

	// Create a block
	block := cmsstore.NewBlock()
	block.SetType("text")
	block.SetContent("Test content")
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetSiteID(site.ID())
	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	tests := []struct {
		name        string
		blockID     string
		expectError bool
		expectedID  string
		expectedErr string
	}{
		{
			name:        "delete block with full ID",
			blockID:     block.ID(),
			expectError: false,
			expectedID:  cmsstore.ShortenID(block.ID()),
		},
		{
			name:        "delete block with shortened ID",
			blockID:     cmsstore.ShortenID(block.ID()),
			expectError: false,
			expectedID:  cmsstore.ShortenID(block.ID()),
		},
		{
			name:        "delete non-existent block",
			blockID:     "non_existent_id",
			expectError: true,
			expectedErr: "block not found",
		},
		{
			name:        "delete block with empty ID",
			blockID:     "",
			expectError: true,
			expectedErr: "missing required parameter: id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetID := tt.blockID
			if tt.name == "delete block with full ID" || tt.name == "delete block with shortened ID" {
				// Create a fresh block for each positive test case
				b := cmsstore.NewBlock()
				b.SetType("text")
				b.SetContent("Test content")
				b.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
				b.SetSiteID(site.ID())
				err = store.BlockCreate(context.Background(), b)
				if err != nil {
					t.Fatalf("Failed to create block: %v", err)
				}

				if tt.name == "delete block with full ID" {
					targetID = b.ID()
				} else {
					targetID = cmsstore.ShortenID(b.ID())
				}
				// Update expectedID to match the new block
				tt.expectedID = cmsstore.ShortenID(b.ID())
			}

			// Call the tool
			deletePayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "delete",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "block_delete",
					"arguments": map[string]any{
						"id": targetID,
					},
				},
			}

			deleteBody, err := json.Marshal(deletePayload)
			if err != nil {
				t.Fatalf("Failed to marshal payload: %v", err)
			}

			deleteResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(deleteBody))
			if err != nil {
				t.Fatalf("Failed to post request: %v", err)
			}
			defer deleteResp.Body.Close()

			deleteRespBytes, err := io.ReadAll(deleteResp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			// Parse the result
			var response map[string]any
			err = json.Unmarshal(deleteRespBytes, &response)
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

func TestBlockUpsert_WithPageID(t *testing.T) {
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

	// Create a block with page_id
	upsertPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "upsert",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "block_upsert",
			"arguments": map[string]any{
				"type":    "text",
				"content": "Page block content",
				"status":  cmsstore.BLOCK_STATUS_ACTIVE,
				"site_id": cmsstore.ShortenID(site.ID()),
				"page_id": cmsstore.ShortenID(page.ID()),
				"name":    "Page Block",
				"handle":  "page-block",
			},
		},
	}

	upsertBody, err := json.Marshal(upsertPayload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	upsertResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(upsertBody))
	if err != nil {
		t.Fatalf("Failed to post request: %v", err)
	}
	defer upsertResp.Body.Close()

	upsertRespBytes, err := io.ReadAll(upsertResp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	// Parse the result
	var response map[string]any
	err = json.Unmarshal(upsertRespBytes, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
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

	var blockData map[string]any
	err = json.Unmarshal([]byte(text), &blockData)
	if err != nil {
		t.Fatalf("Failed to unmarshal block data: %v", err)
	}

	if blockData["type"].(string) != "text" {
		t.Errorf("Expected type 'text', got '%s'", blockData["type"].(string))
	}
	if blockData["content"].(string) != "Page block content" {
		t.Errorf("Expected content 'Page block content', got '%s'", blockData["content"].(string))
	}
	if blockData["status"].(string) != cmsstore.BLOCK_STATUS_ACTIVE {
		t.Errorf("Expected status '%s', got '%s'", cmsstore.BLOCK_STATUS_ACTIVE, blockData["status"].(string))
	}
	if blockData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
		t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
	}
	if blockData["page_id"].(string) != cmsstore.ShortenID(page.ID()) {
		t.Errorf("Expected page_id '%s', got '%s'", cmsstore.ShortenID(page.ID()), blockData["page_id"].(string))
	}
}

func TestBlockList_WithPageID(t *testing.T) {
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

	// Create blocks for different pages
	pageBlock := cmsstore.NewBlock()
	pageBlock.SetType("text")
	pageBlock.SetContent("Page block content")
	pageBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	pageBlock.SetSiteID(site.ID())
	pageBlock.SetPageID(page.ID())
	pageBlock.SetName("Page Block")
	err = store.BlockCreate(context.Background(), pageBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	siteBlock := cmsstore.NewBlock()
	siteBlock.SetType("image")
	siteBlock.SetContent("Site block content")
	siteBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	siteBlock.SetSiteID(site.ID())
	siteBlock.SetName("Site Block")
	err = store.BlockCreate(context.Background(), siteBlock)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Test filtering by page_id
	t.Run("list blocks by page_id", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "block_list",
				"arguments": map[string]any{
					"page_id": cmsstore.ShortenID(page.ID()),
					"limit":   10,
					"offset":  0,
				},
			},
		}

		listBody, err := json.Marshal(listPayload)
		if err != nil {
			t.Fatalf("Failed to marshal payload: %v", err)
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		if err != nil {
			t.Fatalf("Failed to unmarshal block list: %v", err)
		}

		items, ok := blockList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return only the block for the specific page
		if len(items) != 1 {
			t.Errorf("Expected only block for the page, got %d", len(items))
		}
		item := items[0].(map[string]interface{})
		if item["page_id"].(string) != cmsstore.ShortenID(page.ID()) {
			t.Errorf("Expected page_id '%s', got '%s'", cmsstore.ShortenID(page.ID()), item["page_id"].(string))
		}
		if item["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
			t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), item["site_id"].(string))
		}
	})
}
