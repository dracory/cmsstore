package mcp_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	"github.com/dracory/cmsstore/mcp"
// 	"github.com/dracory/cmsstore/testutils"
// 	_ "github.com/mattn/go-sqlite3"
// )

// func initMCP(t *testing.T) (mcp.MCPInterface, *httptest.Server) {
// 	t.Helper()

// 	t.Log("Initializing test store...")
// 	store, err := testutils.InitStore(":memory:")
// 	if err != nil {
// 		t.Fatalf("Failed to initialize store: %v", err)
// 	}

// 	t.Log("Creating new MCP handler instance...")
// 	mcpHandler := mcp.NewMCP(store)

// 	// Create a test HTTP server with the MCP handler
// 	server := httptest.NewServer(mcpHandler.Handler())
// 	t.Logf("Test server started successfully on %s", server.URL)

// 	return mcpHandler, server
// }

// func Test_MCP_CreatePage(t *testing.T) {
// 	_, server := initMCP(t)
// 	defer server.Close()

// 	tests := []struct {
// 		name           string
// 		request        map[string]interface{}
// 		wantStatus     int
// 		wantInResponse []string
// 	}{
// 		{
// 			name: "Valid page creation",
// 			request: map[string]interface{}{
// 				"jsonrpc": "2.0",
// 				"id":      "1",
// 				"method":  "call_tool",
// 				"params": map[string]interface{}{
// 					"tool_name": "page_create",
// 					"arguments": map[string]interface{}{
// 						"title":   "Test Page",
// 						"content": "This is a test page content",
// 						"status":  "published",
// 					},
// 				},
// 			},
// 			wantStatus:     http.StatusOK,
// 			wantInResponse: []string{"success", "true", "id", "title", "Test Page"},
// 		},
// 		{
// 			name: "Missing required title",
// 			request: map[string]interface{}{
// 				"jsonrpc": "2.0",
// 				"id":      "2",
// 				"method":  "call_tool",
// 				"params": map[string]interface{}{
// 					"tool_name": "page_create",
// 					"arguments": map[string]interface{}{
// 						"content": "This is a test page content",
// 					},
// 				},
// 			},
// 			wantStatus:     http.StatusOK,
// 			wantInResponse: []string{"error", "missing required parameter"},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Convert request to JSON
// 			reqBody, err := json.Marshal(tt.request)
// 			if err != nil {
// 				t.Fatalf("Failed to marshal request: %v", err)
// 			}

// 			// Create HTTP request
// 			req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(reqBody))
// 			if err != nil {
// 				t.Fatalf("Failed to create request: %v", err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			// Send request
// 			client := &http.Client{}
// 			resp, err := client.Do(req)
// 			if err != nil {
// 				t.Fatalf("Failed to send request: %v", err)
// 			}
// 			defer resp.Body.Close()

// 			// Check status code
// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
// 			}

// 			// Read response body
// 			body, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Fatalf("Failed to read response body: %v", err)
// 			}

// 			// Check response contains expected strings
// 			respStr := string(body)
// 			for _, want := range tt.wantInResponse {
// 				if !strings.Contains(respStr, want) {
// 					t.Errorf("Response does not contain %q: %s", want, respStr)
// 				}
// 			}
// 		})
// 	}
// }

// func Test_MCP_GetPage(t *testing.T) {
// 	_, server := initMCP(t)
// 	defer server.Close()

// 	// First create a page to retrieve
// 	createReq := map[string]interface{}{
// 		"jsonrpc": "2.0",
// 		"id":      "create",
// 		"method":  "call_tool",
// 		"params": map[string]interface{}{
// 			"tool_name": "page_create",
// 			"arguments": map[string]interface{}{
// 				"title":   "Test Get Page",
// 				"content": "This is a page to retrieve",
// 				"status":  "published",
// 			},
// 		},
// 	}

// 	// Convert request to JSON
// 	reqBody, err := json.Marshal(createReq)
// 	if err != nil {
// 		t.Fatalf("Failed to marshal request: %v", err)
// 	}

// 	// Create HTTP request
// 	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		t.Fatalf("Failed to create request: %v", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Fatalf("Failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Read response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Fatalf("Failed to read response body: %v", err)
// 	}

// 	// Extract page ID from response
// 	var createResp map[string]interface{}
// 	if err := json.Unmarshal(body, &createResp); err != nil {
// 		t.Fatalf("Failed to unmarshal response: %v", err)
// 	}

// 	result := createResp["result"].(map[string]interface{})
// 	text := result["text"].(string)
// 	var pageData map[string]interface{}
// 	if err := json.Unmarshal([]byte(text), &pageData); err != nil {
// 		t.Fatalf("Failed to unmarshal page data: %v", err)
// 	}

// 	pageID := pageData["id"].(string)

// 	// Now test getting the page
// 	tests := []struct {
// 		name           string
// 		request        map[string]interface{}
// 		wantStatus     int
// 		wantInResponse []string
// 	}{
// 		{
// 			name: "Get existing page",
// 			request: map[string]interface{}{
// 				"jsonrpc": "2.0",
// 				"id":      "get",
// 				"method":  "call_tool",
// 				"params": map[string]interface{}{
// 					"tool_name": "page_get",
// 					"arguments": map[string]interface{}{
// 						"id": pageID,
// 					},
// 				},
// 			},
// 			wantStatus:     http.StatusOK,
// 			wantInResponse: []string{"Test Get Page", "This is a page to retrieve"},
// 		},
// 		{
// 			name: "Get non-existent page",
// 			request: map[string]interface{}{
// 				"jsonrpc": "2.0",
// 				"id":      "get_nonexistent",
// 				"method":  "call_tool",
// 				"params": map[string]interface{}{
// 					"tool_name": "page_get",
// 					"arguments": map[string]interface{}{
// 						"id": "nonexistent-id",
// 					},
// 				},
// 			},
// 			wantStatus:     http.StatusOK,
// 			wantInResponse: []string{"page not found"},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Convert request to JSON
// 			reqBody, err := json.Marshal(tt.request)
// 			if err != nil {
// 				t.Fatalf("Failed to marshal request: %v", err)
// 			}

// 			// Create HTTP request
// 			req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(reqBody))
// 			if err != nil {
// 				t.Fatalf("Failed to create request: %v", err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			// Send request
// 			client := &http.Client{}
// 			resp, err := client.Do(req)
// 			if err != nil {
// 				t.Fatalf("Failed to send request: %v", err)
// 			}
// 			defer resp.Body.Close()

// 			// Check status code
// 			if resp.StatusCode != tt.wantStatus {
// 				t.Errorf("Expected status %d, got %d", tt.wantStatus, resp.StatusCode)
// 			}

// 			// Read response body
// 			body, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Fatalf("Failed to read response body: %v", err)
// 			}

// 			// Check response contains expected strings
// 			respStr := string(body)
// 			for _, want := range tt.wantInResponse {
// 				if !strings.Contains(respStr, want) {
// 					t.Errorf("Response does not contain %q: %s", want, respStr)
// 				}
// 			}
// 		})
// 	}
// }
