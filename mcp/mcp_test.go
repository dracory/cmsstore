package mcp_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/mcp"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

func rpcResultText(t *testing.T, respBytes []byte) string {
	t.Helper()

	var rpcResp map[string]any
	if err := json.Unmarshal(respBytes, &rpcResp); err != nil {
		t.Fatalf("Failed to unmarshal json-rpc response: %v. Body=%s", err, string(respBytes))
	}

	result, ok := rpcResp["result"].(map[string]any)
	if !ok {
		t.Fatalf("Expected response to have result: %s", string(respBytes))
	}

	content, ok := result["content"].([]any)
	if !ok || len(content) == 0 {
		t.Fatalf("Expected response result.content: %s", string(respBytes))
	}

	item0, ok := content[0].(map[string]any)
	if !ok {
		t.Fatalf("Expected response result.content[0] object: %s", string(respBytes))
	}

	text, ok := item0["text"].(string)
	if !ok {
		t.Fatalf("Expected response result.content[0].text: %s", string(respBytes))
	}

	return text
}

func initMCPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	mcpHandler := mcp.NewMCP(store)

	server := httptest.NewServer(http.HandlerFunc(mcpHandler.Handler))
	return server, server.Close
}

func initMCPServerWithStore(t *testing.T) (*httptest.Server, cmsstore.StoreInterface, func()) {
	t.Helper()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize store: %v", err)
	}

	mcpHandler := mcp.NewMCP(store)
	server := httptest.NewServer(http.HandlerFunc(mcpHandler.Handler))
	return server, store, server.Close
}

func Test_MCP_ListTools(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	reqPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "list_tools",
		"params":  map[string]any{},
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	respStr := string(bodyBytes)
	if !strings.Contains(respStr, "list_tools") && !strings.Contains(respStr, "tools") {
		t.Fatalf("Unexpected response: %s", respStr)
	}
	if !strings.Contains(respStr, "page_upsert") {
		t.Fatalf("Expected tools list to contain page_upsert: %s", respStr)
	}
	if !strings.Contains(respStr, "cms_schema") {
		t.Fatalf("Expected tools list to contain cms_schema: %s", respStr)
	}
}

func Test_MCP_PageCreate_And_PageGet(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	createPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "create",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_upsert",
			"arguments": map[string]any{
				"title":   "Test Page",
				"content": "Hello",
				"status":  "published",
			},
		},
	}

	createBody, err := json.Marshal(createPayload)
	if err != nil {
		t.Fatalf("Failed to marshal create request: %v", err)
	}

	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	defer createResp.Body.Close()

	createRespBytes, err := io.ReadAll(createResp.Body)
	if err != nil {
		t.Fatalf("Failed to read create response body: %v", err)
	}

	text := rpcResultText(t, createRespBytes)

	var pageData map[string]any
	if err := json.Unmarshal([]byte(text), &pageData); err != nil {
		t.Fatalf("Failed to unmarshal page data: %v", err)
	}

	pageID, _ := pageData["id"].(string)
	if strings.TrimSpace(pageID) == "" {
		t.Fatalf("Expected page id to be returned: %v", pageData)
	}

	getPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "get",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_get",
			"arguments": map[string]any{
				"id": pageID,
			},
		},
	}

	getBody, err := json.Marshal(getPayload)
	if err != nil {
		t.Fatalf("Failed to marshal get request: %v", err)
	}

	getResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(getBody))
	if err != nil {
		t.Fatalf("Failed to send get request: %v", err)
	}
	defer getResp.Body.Close()

	getRespBytes, err := io.ReadAll(getResp.Body)
	if err != nil {
		t.Fatalf("Failed to read get response body: %v", err)
	}

	getRespStr := string(getRespBytes)
	if !strings.Contains(getRespStr, "Test Page") {
		t.Fatalf("Expected get response to contain title: %s", getRespStr)
	}
}

func Test_MCP_PageList(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	createPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "create",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_upsert",
			"arguments": map[string]any{
				"title":   "List Page",
				"content": "Hello",
				"status":  "published",
			},
		},
	}

	createBody, err := json.Marshal(createPayload)
	if err != nil {
		t.Fatalf("Failed to marshal create request: %v", err)
	}

	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	createRespBytes, err := io.ReadAll(createResp.Body)
	_ = createResp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read create response body: %v", err)
	}
	createText := rpcResultText(t, createRespBytes)
	if !strings.Contains(createText, "List Page") {
		t.Fatalf("Expected page_upsert response to contain created page title. Got: %s", createText)
	}

	listPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "list",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_list",
			"arguments": map[string]any{"limit": 10, "offset": 0},
		},
	}

	listBody, err := json.Marshal(listPayload)
	if err != nil {
		t.Fatalf("Failed to marshal list request: %v", err)
	}

	listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
	if err != nil {
		t.Fatalf("Failed to send list request: %v", err)
	}
	defer listResp.Body.Close()

	listRespBytes, err := io.ReadAll(listResp.Body)
	if err != nil {
		t.Fatalf("Failed to read list response body: %v", err)
	}

	if !strings.Contains(string(listRespBytes), "List Page") {
		t.Fatalf("Expected page_list response to contain created page title: %s", string(listRespBytes))
	}
}

func Test_MCP_CmsSchema(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	schemaPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "schema",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "cms_schema",
			"arguments": map[string]any{},
		},
	}

	schemaBody, err := json.Marshal(schemaPayload)
	if err != nil {
		t.Fatalf("Failed to marshal schema request: %v", err)
	}

	schemaResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(schemaBody))
	if err != nil {
		t.Fatalf("Failed to send schema request: %v", err)
	}
	defer schemaResp.Body.Close()

	schemaRespBytes, err := io.ReadAll(schemaResp.Body)
	if err != nil {
		t.Fatalf("Failed to read schema response body: %v", err)
	}

	schemaText := rpcResultText(t, schemaRespBytes)
	if !strings.Contains(schemaText, "\"entities\"") || !strings.Contains(schemaText, "\"tools\"") {
		t.Fatalf("Expected cms_schema response to contain entities and tools. Got: %s", schemaText)
	}
	if !strings.Contains(schemaText, "\"page\"") {
		t.Fatalf("Expected cms_schema response to contain page entity. Got: %s", schemaText)
	}
}

func Test_MCP_Initialize(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	reqPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "init",
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2025-06-18",
			"clientInfo": map[string]any{
				"name":    "test",
				"version": "0",
			},
			"capabilities": map[string]any{},
		},
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	respStr := string(bodyBytes)
	if !strings.Contains(respStr, "\"protocolVersion\"") {
		t.Fatalf("Expected initialize response to contain protocolVersion: %s", respStr)
	}
	if !strings.Contains(respStr, "\"serverInfo\"") {
		t.Fatalf("Expected initialize response to contain serverInfo: %s", respStr)
	}
}

func Test_MCP_ToolsList_StandardMethod(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	reqPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "tools/list",
		"params":  map[string]any{},
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	respStr := string(bodyBytes)
	if !strings.Contains(respStr, "page_upsert") {
		t.Fatalf("Expected tools list to contain page_upsert: %s", respStr)
	}
	if !strings.Contains(respStr, "cms_schema") {
		t.Fatalf("Expected tools list to contain cms_schema: %s", respStr)
	}
}

func Test_MCP_ToolsCall_StandardMethod(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	createPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "create",
		"method":  "tools/call",
		"params": map[string]any{
			"name": "page_upsert",
			"arguments": map[string]any{
				"title":   "Tools Call Page",
				"content": "Hello",
				"status":  "published",
			},
		},
	}

	createBody, err := json.Marshal(createPayload)
	if err != nil {
		t.Fatalf("Failed to marshal create request: %v", err)
	}

	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	defer createResp.Body.Close()

	createRespBytes, err := io.ReadAll(createResp.Body)
	if err != nil {
		t.Fatalf("Failed to read create response body: %v", err)
	}

	if !strings.Contains(string(createRespBytes), "Tools Call Page") {
		t.Fatalf("Expected tools/call page_upsert response to contain title: %s", string(createRespBytes))
	}
}

func Test_MCP_PageUpdate_WithUpdatesObject(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	createPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "create",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_upsert",
			"arguments": map[string]any{
				"title":   "Before Update",
				"content": "Hello",
				"status":  "published",
			},
		},
	}

	createBody, err := json.Marshal(createPayload)
	if err != nil {
		t.Fatalf("Failed to marshal create request: %v", err)
	}

	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	defer createResp.Body.Close()

	createRespBytes, err := io.ReadAll(createResp.Body)
	if err != nil {
		t.Fatalf("Failed to read create response body: %v", err)
	}

	text := rpcResultText(t, createRespBytes)
	var pageData map[string]any
	if err := json.Unmarshal([]byte(text), &pageData); err != nil {
		t.Fatalf("Failed to unmarshal page data: %v", err)
	}
	pageID, _ := pageData["id"].(string)
	if strings.TrimSpace(pageID) == "" {
		t.Fatalf("Expected page id to be returned: %v", pageData)
	}

	updatePayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "update",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_upsert",
			"arguments": map[string]any{
				"id":    pageID,
				"title": "After Update",
				"alias": "/after-update",
			},
		},
	}

	updateBody, err := json.Marshal(updatePayload)
	if err != nil {
		t.Fatalf("Failed to marshal update request: %v", err)
	}

	updateResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatalf("Failed to send update request: %v", err)
	}
	defer updateResp.Body.Close()

	updateRespBytes, err := io.ReadAll(updateResp.Body)
	if err != nil {
		t.Fatalf("Failed to read update response body: %v", err)
	}

	if !strings.Contains(string(updateRespBytes), "After Update") {
		t.Fatalf("Expected page_upsert response to contain updated title: %s", string(updateRespBytes))
	}
}

func Test_MCP_PageUpdate_FlatArgs_WithNumericID(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	page := cmsstore.NewPage()
	page.SetID("20260108160058473")
	page.SetTitle("Numeric ID")
	page.SetStatus("published")

	sites, err := store.SiteList(context.Background(), cmsstore.SiteQuery())
	if err != nil {
		t.Fatalf("Failed to list sites: %v", err)
	}
	if len(sites) == 0 {
		defaultSite := cmsstore.NewSite()
		defaultSite.SetName("Default Site")
		defaultSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
		if err := store.SiteCreate(context.Background(), defaultSite); err != nil {
			t.Fatalf("Failed to create site: %v", err)
		}
		page.SetSiteID(defaultSite.ID())
	} else {
		page.SetSiteID(sites[0].ID())
	}

	if err := store.PageCreate(context.Background(), page); err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	updatePayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "update",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "page_update",
			"arguments": map[string]any{
				"id":          json.Number("20260108160058473"),
				"alias":       "/ai-created",
				"name":        "AI Created",
				"template_id": json.Number("20240604061502639"),
			},
		},
	}

	updateBody, err := json.Marshal(updatePayload)
	if err != nil {
		t.Fatalf("Failed to marshal update request: %v", err)
	}

	updateResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(updateBody))
	if err != nil {
		t.Fatalf("Failed to send update request: %v", err)
	}
	defer updateResp.Body.Close()

	updateRespBytes, err := io.ReadAll(updateResp.Body)
	if err != nil {
		t.Fatalf("Failed to read update response body: %v", err)
	}

	if strings.Contains(string(updateRespBytes), "missing required parameter: id") {
		t.Fatalf("Expected numeric id to be accepted, got: %s", string(updateRespBytes))
	}
}

func Test_MCP_MenuList(t *testing.T) {
	server, cleanup := initMCPServer(t)
	defer cleanup()

	createPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "create",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "menu_upsert",
			"arguments": map[string]any{
				"name":   "Main Menu",
				"status": "active",
			},
		},
	}

	createBody, err := json.Marshal(createPayload)
	if err != nil {
		t.Fatalf("Failed to marshal create request: %v", err)
	}
	createResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		t.Fatalf("Failed to send create request: %v", err)
	}
	_ = createResp.Body.Close()

	listPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "list",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "menu_list",
			"arguments": map[string]any{"limit": 10, "offset": 0},
		},
	}

	listBody, err := json.Marshal(listPayload)
	if err != nil {
		t.Fatalf("Failed to marshal list request: %v", err)
	}

	listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
	if err != nil {
		t.Fatalf("Failed to send list request: %v", err)
	}
	defer listResp.Body.Close()

	listRespBytes, err := io.ReadAll(listResp.Body)
	if err != nil {
		t.Fatalf("Failed to read list response body: %v", err)
	}
	if !strings.Contains(string(listRespBytes), "Main Menu") {
		t.Fatalf("Expected menu_list response to contain created menu name: %s", string(listRespBytes))
	}
}

func Test_MCP_SiteList(t *testing.T) {
	server, store, cleanup := initMCPServerWithStore(t)
	defer cleanup()

	site := cmsstore.NewSite()
	site.SetID("20260108160000001")
	site.SetName("Example Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	_, _ = site.SetDomainNames([]string{"example.test"})
	if err := store.SiteCreate(context.Background(), site); err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	listPayload := map[string]any{
		"jsonrpc": "2.0",
		"id":      "list",
		"method":  "call_tool",
		"params": map[string]any{
			"tool_name": "site_list",
			"arguments": map[string]any{"limit": 10, "offset": 0},
		},
	}

	listBody, err := json.Marshal(listPayload)
	if err != nil {
		t.Fatalf("Failed to marshal list request: %v", err)
	}

	listResp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(listBody))
	if err != nil {
		t.Fatalf("Failed to send list request: %v", err)
	}
	defer listResp.Body.Close()

	listRespBytes, err := io.ReadAll(listResp.Body)
	if err != nil {
		t.Fatalf("Failed to read list response body: %v", err)
	}
	if !strings.Contains(string(listRespBytes), "Example Site") {
		t.Fatalf("Expected site_list response to contain created site name: %s", string(listRespBytes))
	}
}
