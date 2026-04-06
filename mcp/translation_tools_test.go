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

func TestTranslationGet(t *testing.T) {
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

	// Create a translation
	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetHandle("test-translation")
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetSiteID(site.ID())

	// Set content with multiple languages
	content := map[string]string{
		"en": "Hello World",
		"fr": "Bonjour le monde",
		"es": "Hola Mundo",
	}
	err = translation.SetContent(content)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	tests := []struct {
		name          string
		translationID string
		expectError   bool
		expectedID    string
		expectedErr   string
	}{
		{
			name:          "get translation with full ID",
			translationID: translation.ID(),
			expectError:   false,
			expectedID:    cmsstore.ShortenID(translation.ID()),
		},
		{
			name:          "get translation with shortened ID",
			translationID: cmsstore.ShortenID(translation.ID()),
			expectError:   false,
			expectedID:    cmsstore.ShortenID(translation.ID()),
		},
		{
			name:          "get non-existent translation",
			translationID: "non_existent_id",
			expectError:   true,
			expectedErr:   "translation not found",
		},
		{
			name:          "get translation with empty ID",
			translationID: "",
			expectError:   true,
			expectedErr:   "missing required parameter: id",
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
					"tool_name": "translation_get",
					"arguments": map[string]any{
						"id": tt.translationID,
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

				var translationData map[string]any
				err = json.Unmarshal([]byte(text), &translationData)
				if err != nil {
					t.Fatalf("Failed to unmarshal translation data: %v", err)
				}

				if translationData["id"].(string) != tt.expectedID {
					t.Errorf("Expected id '%s', got '%s'", tt.expectedID, translationData["id"].(string))
				}
				if translationData["name"].(string) != "Test Translation" {
					t.Errorf("Expected name 'Test Translation', got '%s'", translationData["name"].(string))
				}
				if translationData["handle"].(string) != "test-translation" {
					t.Errorf("Expected handle 'test-translation', got '%s'", translationData["handle"].(string))
				}
				if translationData["status"].(string) != cmsstore.TRANSLATION_STATUS_ACTIVE {
					t.Errorf("Expected status '%s', got '%s'", cmsstore.TRANSLATION_STATUS_ACTIVE, translationData["status"].(string))
				}
				if translationData["site_id"].(string) != cmsstore.ShortenID(site.ID()) {
					t.Errorf("Expected site_id '%s', got '%s'", cmsstore.ShortenID(site.ID()), translationData["site_id"].(string))
				}

				// Check content
				contentMap, ok := translationData["content"].(map[string]any)
				if !ok {
					t.Fatalf("Expected content to be a map")
				}
				if contentMap["en"].(string) != "Hello World" {
					t.Errorf("Expected content.en 'Hello World', got '%s'", contentMap["en"].(string))
				}
				if contentMap["fr"].(string) != "Bonjour le monde" {
					t.Errorf("Expected content.fr 'Bonjour le monde', got '%s'", contentMap["fr"].(string))
				}
				if contentMap["es"].(string) != "Hola Mundo" {
					t.Errorf("Expected content.es 'Hola Mundo', got '%s'", contentMap["es"].(string))
				}
			}
		})
	}
}

func TestTranslationList(t *testing.T) {
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

	// Create translations with different properties
	activeTranslation := cmsstore.NewTranslation()
	activeTranslation.SetName("Active Translation")
	activeTranslation.SetHandle("active-translation")
	activeTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	activeTranslation.SetSiteID(site.ID())
	err = store.TranslationCreate(context.Background(), activeTranslation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	draftTranslation := cmsstore.NewTranslation()
	draftTranslation.SetName("Draft Translation")
	draftTranslation.SetHandle("draft-translation")
	draftTranslation.SetStatus(cmsstore.TRANSLATION_STATUS_DRAFT)
	draftTranslation.SetSiteID(site.ID())
	err = store.TranslationCreate(context.Background(), draftTranslation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Test listing all translations
	t.Run("list all translations", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		if err != nil {
			t.Fatalf("Failed to unmarshal translation list: %v", err)
		}

		items, ok := translationList["items"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'items' to be a slice")
		}

		// Should return both translations
		if len(items) != 2 {
			t.Errorf("Expected 2 translations, got %d", len(items))
		}
	})

	// Test filtering by site_id
	t.Run("list translations by site_id", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return both translations for the site
		assert.Equal(t, 2, len(items), "Expected both translations for the site")
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			assert.Equal(t, cmsstore.ShortenID(site.ID()), itemMap["site_id"].(string))
		}
	})

	// Test filtering by status
	t.Run("list translations by status", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
				"arguments": map[string]any{
					"status": cmsstore.TRANSLATION_STATUS_ACTIVE,
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only active translation
		assert.Equal(t, 1, len(items), "Expected only active translation")
		item := items[0].(map[string]interface{})
		assert.Equal(t, cmsstore.TRANSLATION_STATUS_ACTIVE, item["status"].(string))
	})

	// Test filtering by handle
	t.Run("list translations by handle", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
				"arguments": map[string]any{
					"handle": "active-translation",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only the translation with matching handle
		assert.Equal(t, 1, len(items), "Expected only translation with matching handle")
		item := items[0].(map[string]interface{})
		assert.Equal(t, "active-translation", item["handle"].(string))
	})

	// Test pagination
	t.Run("list translations with pagination", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return only 1 translation due to limit
		assert.Equal(t, 1, len(items), "Expected only 1 translation due to limit")
	})
}

func TestTranslationUpsert_Create(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	tests := []struct {
		name            string
		translationName string
		handle          string
		status          string
		siteID          string
		content         map[string]string
		expectError     bool
		expectedErr     string
	}{
		{
			name:            "create translation with all fields",
			translationName: "New Translation",
			handle:          "new-translation",
			status:          cmsstore.TRANSLATION_STATUS_ACTIVE,
			siteID:          cmsstore.ShortenID(site.ID()),
			content:         map[string]string{"en": "Hello", "fr": "Bonjour"},
			expectError:     false,
		},
		{
			name:            "create translation with minimal fields",
			translationName: "Minimal Translation",
			handle:          "",
			status:          cmsstore.TRANSLATION_STATUS_DRAFT,
			siteID:          "",
			content:         nil,
			expectError:     false,
		},
		{
			name:            "create translation with empty name",
			translationName: "",
			handle:          "test-translation",
			status:          cmsstore.TRANSLATION_STATUS_ACTIVE,
			siteID:          cmsstore.ShortenID(site.ID()),
			content:         map[string]string{"en": "Test"},
			expectError:     true,
			expectedErr:     "missing required parameter: name",
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
					"tool_name": "translation_upsert",
					"arguments": map[string]any{
						"name":    tt.translationName,
						"handle":  tt.handle,
						"status":  tt.status,
						"site_id": tt.siteID,
						"content": tt.content,
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

				var translationData map[string]any
				err = json.Unmarshal([]byte(text), &translationData)
				require.NoError(t, err)

				assert.Equal(t, tt.translationName, translationData["name"].(string))
				assert.Equal(t, tt.status, translationData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), translationData["site_id"].(string))

				if tt.handle != "" {
					assert.Equal(t, tt.handle, translationData["handle"].(string))
				}

				if tt.content != nil {
					contentMap, ok := translationData["content"].(map[string]any)
					require.True(t, ok, "Expected content to be a map")
					for lang, text := range tt.content {
						assert.Equal(t, text, contentMap[lang])
					}
				}

				// Check new fields
				assert.NotEmpty(t, translationData["created_at"].(string))
				assert.NotEmpty(t, translationData["updated_at"].(string))
				assert.NotEmpty(t, translationData["soft_deleted_at"].(string))
				assert.NotNil(t, translationData["metas"])
			}
		})
	}
}

func TestTranslationUpsert_Update(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a translation
	translation := cmsstore.NewTranslation()
	translation.SetName("Original Translation")
	translation.SetHandle("original-translation")
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetSiteID(site.ID())

	// Set original content
	originalContent := map[string]string{
		"en": "Original English",
		"fr": "Original French",
	}
	err = translation.SetContent(originalContent)
	require.NoError(t, err)

	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	tests := []struct {
		name            string
		translationID   string
		translationName string
		handle          string
		status          string
		content         map[string]string
		expectError     bool
		expectedErr     string
	}{
		{
			name:            "update translation with full ID",
			translationID:   translation.ID(),
			translationName: "Updated Translation",
			handle:          "updated-translation",
			status:          cmsstore.TRANSLATION_STATUS_DRAFT,
			content:         map[string]string{"en": "Updated English", "de": "Updated German"},
			expectError:     false,
		},
		{
			name:            "update translation with shortened ID",
			translationID:   cmsstore.ShortenID(translation.ID()),
			translationName: "Updated Translation",
			handle:          "updated-translation",
			status:          cmsstore.TRANSLATION_STATUS_DRAFT,
			content:         map[string]string{"en": "Updated English", "de": "Updated German"},
			expectError:     false,
		},
		{
			name:            "update non-existent translation",
			translationID:   "non_existent_id",
			translationName: "Updated Translation",
			handle:          "updated-translation",
			status:          cmsstore.TRANSLATION_STATUS_DRAFT,
			content:         map[string]string{"en": "Updated English"},
			expectError:     true,
			expectedErr:     "translation not found",
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
					"tool_name": "translation_upsert",
					"arguments": map[string]any{
						"id":      tt.translationID,
						"name":    tt.translationName,
						"handle":  tt.handle,
						"status":  tt.status,
						"content": tt.content,
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

				var translationData map[string]any
				err = json.Unmarshal([]byte(text), &translationData)
				require.NoError(t, err)

				assert.Equal(t, tt.translationName, translationData["name"].(string))
				assert.Equal(t, tt.handle, translationData["handle"].(string))
				assert.Equal(t, tt.status, translationData["status"].(string))
				assert.Equal(t, cmsstore.ShortenID(site.ID()), translationData["site_id"].(string))

				// Check updated content
				if tt.content != nil {
					contentMap, ok := translationData["content"].(map[string]any)
					require.True(t, ok, "Expected content to be a map")
					for lang, text := range tt.content {
						assert.Equal(t, text, contentMap[lang])
					}
				}
			}
		})
	}
}

func TestTranslationDelete(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a translation
	translation := cmsstore.NewTranslation()
	translation.SetName("Test Translation")
	translation.SetHandle("test-translation")
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetSiteID(site.ID())
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	tests := []struct {
		name          string
		translationID string
		expectError   bool
		expectedID    string
		expectedErr   string
	}{
		{
			name:          "delete translation with full ID",
			translationID: translation.ID(),
			expectError:   false,
			expectedID:    cmsstore.ShortenID(translation.ID()),
		},
		{
			name:          "delete translation with shortened ID",
			translationID: cmsstore.ShortenID(translation.ID()),
			expectError:   false,
			expectedID:    cmsstore.ShortenID(translation.ID()),
		},
		{
			name:          "delete non-existent translation",
			translationID: "non_existent_id",
			expectError:   true,
			expectedErr:   "translation not found",
		},
		{
			name:          "delete translation with empty ID",
			translationID: "",
			expectError:   true,
			expectedErr:   "missing required parameter: id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetID := tt.translationID
			if tt.name == "delete translation with full ID" || tt.name == "delete translation with shortened ID" {
				// Create a fresh translation for each positive test case
				translationObj := cmsstore.NewTranslation()
				translationObj.SetName("Test Translation")
				translationObj.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
				translationObj.SetSiteID(site.ID())
				err = store.TranslationCreate(context.Background(), translationObj)
				require.NoError(t, err)

				if tt.name == "delete translation with full ID" {
					targetID = translationObj.ID()
				} else {
					targetID = cmsstore.ShortenID(translationObj.ID())
				}
				// Update expectedID to match the new translation
				tt.expectedID = cmsstore.ShortenID(translationObj.ID())
			}

			// Call the tool
			deletePayload := map[string]any{
				"jsonrpc": "2.0",
				"id":      "delete",
				"method":  "call_tool",
				"params": map[string]any{
					"tool_name": "translation_delete",
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

func TestTranslationUpsert_WithDefaultSite(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	// Create a translation without specifying site_id - should use default site
	upsertPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "upsert",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "translation_upsert",
			"arguments": map[string]any{
				"name":    "Default Site Translation",
				"status":  cmsstore.TRANSLATION_STATUS_ACTIVE,
				"handle":  "default-translation",
				"content": map[string]string{"en": "Default content", "fr": "Contenu par défaut"},
				"memo":    "Translation with default site",
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

	var translationData map[string]any
	err = json.Unmarshal([]byte(text), &translationData)
	require.NoError(t, err)

	assert.Equal(t, "Default Site Translation", translationData["name"].(string))
	assert.Equal(t, cmsstore.TRANSLATION_STATUS_ACTIVE, translationData["status"].(string))
	assert.Equal(t, "default-translation", translationData["handle"].(string))
	assert.Equal(t, "Translation with default site", translationData["memo"].(string))
	// site_id should be set to the default site
	assert.NotEmpty(t, translationData["site_id"].(string))

	// Check content
	contentMap, ok := translationData["content"].(map[string]any)
	require.True(t, ok, "Expected content to be a map")
	assert.Equal(t, "Default content", contentMap["en"].(string))
	assert.Equal(t, "Contenu par défaut", contentMap["fr"].(string))
}

func TestTranslationList_WithSoftDeleted(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create a translation
	translation := cmsstore.NewTranslation()
	translation.SetName("Soft Deleted Translation")
	translation.SetHandle("soft-deleted-translation")
	translation.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation.SetSiteID(site.ID())
	err = store.TranslationCreate(context.Background(), translation)
	require.NoError(t, err)

	// Soft delete the translation
	err = store.TranslationSoftDeleteByID(context.Background(), translation.ID())
	require.NoError(t, err)

	// Test listing without include_soft_deleted (should not include soft deleted)
	t.Run("list translations without soft deleted", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should not return the soft deleted translation
		assert.Equal(t, 0, len(items), "Expected no translations (soft deleted should be excluded)")
	})

	// Test listing with include_soft_deleted (should include soft deleted)
	t.Run("list translations with soft deleted", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return the soft deleted translation
		assert.Equal(t, 1, len(items), "Expected 1 translation (soft deleted should be included)")
	})
}

func TestTranslationList_WithOrdering(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), site)
	require.NoError(t, err)

	// Create translations with different names to test ordering
	translation1 := cmsstore.NewTranslation()
	translation1.SetName("Alpha Translation")
	translation1.SetHandle("alpha-translation")
	translation1.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation1.SetSiteID(site.ID())
	err = store.TranslationCreate(context.Background(), translation1)
	require.NoError(t, err)

	translation2 := cmsstore.NewTranslation()
	translation2.SetName("Beta Translation")
	translation2.SetHandle("beta-translation")
	translation2.SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)
	translation2.SetSiteID(site.ID())
	err = store.TranslationCreate(context.Background(), translation2)
	require.NoError(t, err)

	// Test ordering by name ascending
	t.Run("list translations ordered by name ascending", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return translations in alphabetical order
		assert.Equal(t, 2, len(items), "Expected 2 translations")
		assert.Equal(t, "Alpha Translation", items[0].(map[string]interface{})["name"].(string))
		assert.Equal(t, "Beta Translation", items[1].(map[string]interface{})["name"].(string))
	})

	// Test ordering by name descending
	t.Run("list translations ordered by name descending", func(t *testing.T) {
		listPayload := map[string]any{
			"jsonrpc": "2.0",
			"id":      "list",
			"method":  "call_tool",
			"params": map[string]any{
				"tool_name": "translation_list",
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

		var translationList map[string]any
		err = json.Unmarshal([]byte(text), &translationList)
		require.NoError(t, err)

		items, ok := translationList["items"].([]interface{})
		require.True(t, ok, "Expected 'items' to be a slice")

		// Should return translations in reverse alphabetical order
		assert.Equal(t, 2, len(items), "Expected 2 translations")
		assert.Equal(t, "Beta Translation", items[0].(map[string]interface{})["name"].(string))
		assert.Equal(t, "Alpha Translation", items[1].(map[string]interface{})["name"].(string))
	})
}
