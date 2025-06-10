package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gouniverse/cmsstore"
)

// setupBlockTest creates the necessary test data for block tests
func setupBlockTest(t *testing.T) (string, cmsstore.StoreInterface, cmsstore.SiteInterface, cmsstore.PageInterface, cmsstore.BlockInterface, func()) {
	serverURL, store, cleanup := setupTestAPI(t)

	// Create a test site
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site for Blocks")
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	// Create a test page
	testPage := cmsstore.NewPage()
	testPage.SetName("Test Page for Blocks")
	testPage.SetSiteID(testSite.ID())
	testPage.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), testPage)
	if err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	// Create a test block
	testBlock := cmsstore.NewBlock()
	testBlock.SetName("Test Block")
	testBlock.SetContent("Test Content")
	testBlock.SetSiteID(testSite.ID())
	testBlock.SetPageID(testPage.ID())
	testBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	testBlock.SetTemplateID("")
	testBlock.SetParentID("")
	testBlock.SetSequenceInt(-1)
	err = store.BlockCreate(context.Background(), testBlock)
	if err != nil {
		t.Fatalf("Failed to create test block: %v", err)
	}

	return serverURL, store, testSite, testPage, testBlock, cleanup
}

// TestBlockList tests the GET /api/blocks endpoint (list blocks)
func TestBlockList(t *testing.T) {
	serverURL, _, testSite, _, _, cleanup := setupBlockTest(t)
	defer cleanup()

	resp, err := http.Get(serverURL + "/api/blocks?site_id=" + testSite.ID())
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", result["success"])
	}

	blocks, ok := result["blocks"].([]interface{})
	if !ok {
		t.Fatalf("Expected blocks to be an array, got %T", result["blocks"])
	}

	if len(blocks) < 1 {
		t.Errorf("Expected at least one block, got %d", len(blocks))
	}
}

// TestBlockGet tests the GET /api/blocks/{id} endpoint (get single block)
func TestBlockGet(t *testing.T) {
	serverURL, _, _, _, testBlock, cleanup := setupBlockTest(t)
	defer cleanup()

	resp, err := http.Get(serverURL + "/api/blocks/" + testBlock.ID())
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", result["success"])
	}

	block, ok := result["block"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected block to be an object, got %T", result["block"])
	}

	if id, ok := block["id"].(string); !ok || id != testBlock.ID() {
		t.Errorf("Expected id to be %s, got %v", testBlock.ID(), id)
	}

	if name, ok := block["name"].(string); !ok || name != testBlock.Name() {
		t.Errorf("Expected name to be %s, got %v", testBlock.Name(), name)
	}

	if content, ok := block["content"].(string); !ok || content != testBlock.Content() {
		t.Errorf("Expected content to be %s, got %v", testBlock.Content(), content)
	}
}

// TestBlockCreate tests the POST /api/blocks endpoint (create block)
func TestBlockCreate(t *testing.T) {
	serverURL, _, testSite, testPage, _, cleanup := setupBlockTest(t)
	defer cleanup()

	blockData := map[string]interface{}{
		"name":    "New Test Block",
		"content": "New Test Content",
		"site_id": testSite.ID(),
		"page_id": testPage.ID(),
		"status":  cmsstore.BLOCK_STATUS_ACTIVE,
	}

	jsonData, err := json.Marshal(blockData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(serverURL+"/api/blocks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", result["success"])
	}

	// Print the response for debugging
	jsonResponse, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("Response: %s", jsonResponse)

	if _, ok := result["id"].(string); !ok {
		t.Errorf("Expected id to be a string")
	}

	if name, ok := result["name"].(string); !ok || name != "New Test Block" {
		t.Errorf("Expected name to be 'New Test Block', got %v", name)
	}
}

// TestBlockUpdate tests the PUT /api/blocks/{id} endpoint (update block)
func TestBlockUpdate(t *testing.T) {
	serverURL, _, _, _, testBlock, cleanup := setupBlockTest(t)
	defer cleanup()

	updateData := map[string]interface{}{
		"name":    "Updated Test Block",
		"content": "Updated Test Content",
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, serverURL+"/api/blocks/"+testBlock.ID(), bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
		req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", result["success"])
	}

	if name, ok := result["name"].(string); !ok || name != "Updated Test Block" {
		t.Errorf("Expected name to be 'Updated Test Block', got %v", name)
	}
}

// TestBlockDelete tests the DELETE /api/blocks/{id} endpoint (delete block)
func TestBlockDelete(t *testing.T) {
	serverURL, store, _, _, testBlock, cleanup := setupBlockTest(t)
	defer cleanup()

	req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/blocks/"+testBlock.ID(), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := result["success"].(bool); !ok || !success {
		t.Errorf("Expected success to be true, got %v", result["success"])
	}

	// Verify the block was soft deleted
	// Create a query that includes soft-deleted items
	query := cmsstore.BlockQuery().SetSoftDeleteIncluded(true).SetID(testBlock.ID())
	blocks, err := store.BlockList(context.Background(), query)
	if err != nil {
		t.Fatalf("Failed to find block: %v", err)
	}

	if len(blocks) == 0 {
		t.Fatalf("Block should still exist after soft delete")
	}

	block := blocks[0]
	if !block.IsSoftDeleted() {
		t.Errorf("Block should be marked as soft deleted")
	}
}
