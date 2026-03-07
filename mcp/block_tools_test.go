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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	// Create a block
	block := cmsstore.NewBlock()
	block.SetType("text")
	block.SetContent("Test content")
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetSiteID(site.ID())
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

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

				var blockData map[string]any
				err = json.Unmarshal([]byte(text), &blockData)
				require.NoError(t, err)

				assert.Equal(t, tt.expectedID, blockData["id"].(string))
				assert.Equal(t, "text", blockData["type"].(string))
				assert.Equal(t, "Test content", blockData["content"].(string))
				assert.Equal(t, cmsstore.BLOCK_STATUS_ACTIVE, blockData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
				assert.Equal(t, "", blockData["page_id"].(string))
				// New fields assertions
				assert.Equal(t, "", blockData["name"].(string))
				assert.Equal(t, "", blockData["handle"].(string))
				assert.Equal(t, "", blockData["template_id"].(string))
				assert.Equal(t, "", blockData["parent_id"].(string))
				assert.Equal(t, "0", blockData["sequence"].(string))
				assert.Equal(t, "", blockData["editor"].(string))
				assert.Equal(t, "", blockData["memo"].(string))
				assert.NotEmpty(t, blockData["created_at"].(string))
				assert.NotEmpty(t, blockData["updated_at"].(string))
				assert.Equal(t, "9999-12-31 23:59:59", blockData["soft_deleted_at"].(string))
				assert.NotNil(t, blockData["metas"])
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
	require.NoError(t, err)

	// Create blocks with different properties
	activeBlock := cmsstore.NewBlock()
	activeBlock.SetType("text")
	activeBlock.SetContent("Active content")
	activeBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	activeBlock.SetSiteID(site.ID())
	activeBlock.SetName("Active Block")
	activeBlock.SetHandle("active-block")
	err = store.BlockCreate(context.Background(), activeBlock)
	require.NoError(t, err)

	draftBlock := cmsstore.NewBlock()
	draftBlock.SetType("image")
	draftBlock.SetContent("Draft content")
	draftBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	draftBlock.SetSiteID(site.ID())
	draftBlock.SetName("Draft Block")
	draftBlock.SetHandle("draft-block")
	err = store.BlockCreate(context.Background(), draftBlock)
	require.NoError(t, err)

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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		require.NoError(t, err)

		items, ok := blockList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return both blocks
		assert.Equal(t, 2, len(items), "Expected both blocks to be returned")
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		require.NoError(t, err)

		items, ok := blockList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return both blocks for the site
		assert.Equal(t, 2, len(items), "Expected both blocks for the site")
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			assert.Equal(t, cmsstore.ShortenID(site.ID()), itemMap["site_id"].(string))
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		require.NoError(t, err)

		items, ok := blockList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only active block
		assert.Equal(t, 1, len(items), "Expected only active block")
		item := items[0].(map[string]interface{})
		assert.Equal(t, cmsstore.BLOCK_STATUS_ACTIVE, item["status"].(string))
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		require.NoError(t, err)

		items, ok := blockList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only the block with matching handle
		assert.Equal(t, 1, len(items), "Expected only block with matching handle")
		item := items[0].(map[string]interface{})
		assert.Equal(t, "active-block", item["handle"].(string))
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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		require.NoError(t, err)

		items, ok := blockList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only 1 block due to limit
		assert.Equal(t, 1, len(items), "Expected only 1 block due to limit")
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
	require.NoError(t, err)

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

				var blockData map[string]any
				err = json.Unmarshal([]byte(text), &blockData)
				require.NoError(t, err)

				assert.Equal(t, tt.blockType, blockData["type"].(string))
				assert.Equal(t, tt.content, blockData["content"].(string))
				assert.Equal(t, tt.status, blockData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
				// New fields assertions
				assert.Equal(t, tt.blockName, blockData["name"].(string))
				assert.Equal(t, tt.handle, blockData["handle"].(string))
				assert.Equal(t, "", blockData["page_id"].(string))
				assert.Equal(t, "", blockData["template_id"].(string))
				assert.Equal(t, "", blockData["parent_id"].(string))
				assert.Equal(t, cast.ToString(tt.sequence), blockData["sequence"].(string))
				assert.Equal(t, tt.editor, blockData["editor"].(string))
				assert.Equal(t, tt.memo, blockData["memo"].(string))
				assert.NotEmpty(t, blockData["created_at"].(string))
				assert.NotEmpty(t, blockData["updated_at"].(string))
				assert.Equal(t, "9999-12-31 23:59:59", blockData["soft_deleted_at"].(string))
				assert.NotNil(t, blockData["metas"])
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
	require.NoError(t, err)

	// Create a block
	block := cmsstore.NewBlock()
	block.SetType("text")
	block.SetContent("Original content")
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetSiteID(site.ID())
	block.SetName("Original Block")
	block.SetHandle("original-block")
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

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

				var blockData map[string]any
				err = json.Unmarshal([]byte(text), &blockData)
				require.NoError(t, err)

				assert.Equal(t, tt.blockType, blockData["type"].(string))
				assert.Equal(t, tt.content, blockData["content"].(string))
				assert.Equal(t, tt.status, blockData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
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
	require.NoError(t, err)

	// Create a block
	block := cmsstore.NewBlock()
	block.SetType("text")
	block.SetContent("Test content")
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetSiteID(site.ID())
	err = store.BlockCreate(context.Background(), block)
	require.NoError(t, err)

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
				require.NoError(t, err)

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

func TestBlockUpsert_WithPageID(t *testing.T) {
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

	var blockData map[string]any
	err = json.Unmarshal([]byte(text), &blockData)
	require.NoError(t, err)

	assert.Equal(t, "text", blockData["type"].(string))
	assert.Equal(t, "Page block content", blockData["content"].(string))
	assert.Equal(t, cmsstore.BLOCK_STATUS_ACTIVE, blockData["status"].(string))
	assert.Equal(t, cmsstore.ShortenID(site.ID()), blockData["site_id"].(string))
	assert.Equal(t, cmsstore.ShortenID(page.ID()), blockData["page_id"].(string))
}

func TestBlockList_WithPageID(t *testing.T) {
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

	// Create blocks for different pages
	pageBlock := cmsstore.NewBlock()
	pageBlock.SetType("text")
	pageBlock.SetContent("Page block content")
	pageBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	pageBlock.SetSiteID(site.ID())
	pageBlock.SetPageID(page.ID())
	pageBlock.SetName("Page Block")
	err = store.BlockCreate(context.Background(), pageBlock)
	require.NoError(t, err)

	siteBlock := cmsstore.NewBlock()
	siteBlock.SetType("image")
	siteBlock.SetContent("Site block content")
	siteBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	siteBlock.SetSiteID(site.ID())
	siteBlock.SetName("Site Block")
	err = store.BlockCreate(context.Background(), siteBlock)
	require.NoError(t, err)

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

		var blockList map[string]any
		err = json.Unmarshal([]byte(text), &blockList)
		require.NoError(t, err)

		items, ok := blockList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only the block for the specific page
		assert.Equal(t, 1, len(items), "Expected only block for the page")
		item := items[0].(map[string]interface{})
		assert.Equal(t, cmsstore.ShortenID(page.ID()), item["page_id"].(string))
		assert.Equal(t, cmsstore.ShortenID(site.ID()), item["site_id"].(string))
	})
}
