package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gouniverse/cmsstore"
)

func TestBlockEndpoints(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	// Create a test site for our tests
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site for Blocks")
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	// Create a test block for our tests
	testBlock := cmsstore.NewBlock()
	testBlock.SetName("Test Block")
	testBlock.SetContent("Test Content")
	testBlock.SetSiteID(testSite.ID())
	testBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(context.Background(), testBlock)
	if err != nil {
		t.Fatalf("Failed to create test block: %v", err)
	}

	t.Run("List Blocks", func(t *testing.T) {
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
	})

	t.Run("Get Block", func(t *testing.T) {
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

		if id, ok := result["id"].(string); !ok || id != testBlock.ID() {
			t.Errorf("Expected id to be %s, got %v", testBlock.ID(), id)
		}

		if name, ok := result["name"].(string); !ok || name != "Test Block" {
			t.Errorf("Expected name to be 'Test Block', got %v", name)
		}
	})

	t.Run("Create Block", func(t *testing.T) {
		blockData := map[string]interface{}{
			"name":    "New Test Block",
			"content": "New Test Content",
			"site_id": testSite.ID(),
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

		if _, ok := result["id"].(string); !ok {
			t.Errorf("Expected id to be a string")
		}

		if name, ok := result["name"].(string); !ok || name != "New Test Block" {
			t.Errorf("Expected name to be 'New Test Block', got %v", name)
		}
	})

	t.Run("Update Block", func(t *testing.T) {
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
	})

	t.Run("Delete Block", func(t *testing.T) {
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
		block, err := store.BlockFindByID(context.Background(), testBlock.ID())
		if err != nil {
			t.Fatalf("Failed to find block: %v", err)
		}
		if block == nil {
			t.Fatalf("Block should still exist after soft delete")
		}
		if !block.IsSoftDeleted() {
			t.Errorf("Block should be marked as soft deleted")
		}
	})
}
