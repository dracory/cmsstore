package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gouniverse/cmsstore"
)

func TestPageEndpoints(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	// Create a test site for our tests
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site for Pages")
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	// Create a test page for our tests
	testPage := cmsstore.NewPage()
	testPage.SetTitle("Test Page")
	testPage.SetContent("Test Content")
	testPage.SetSiteID(testSite.ID())
	testPage.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(context.Background(), testPage)
	if err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	t.Run("List Pages", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/pages")
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

		pages, ok := result["pages"].([]interface{})
		if !ok {
			t.Fatalf("Expected pages to be an array, got %T", result["pages"])
		}

		if len(pages) < 1 {
			t.Errorf("Expected at least one page, got %d", len(pages))
		}
	})

	t.Run("Get Page", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/pages/" + testPage.ID())
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

		if id, ok := result["id"].(string); !ok || id != testPage.ID() {
			t.Errorf("Expected id to be %s, got %v", testPage.ID(), id)
		}

		if title, ok := result["title"].(string); !ok || title != "Test Page" {
			t.Errorf("Expected title to be 'Test Page', got %v", title)
		}
	})

	t.Run("Create Page", func(t *testing.T) {
		pageData := map[string]interface{}{
			"title":   "New Test Page",
			"content": "New Test Content",
			"site_id": testSite.ID(),
			"status":  cmsstore.PAGE_STATUS_ACTIVE,
		}

		jsonData, err := json.Marshal(pageData)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		resp, err := http.Post(serverURL+"/api/pages", "application/json", bytes.NewBuffer(jsonData))
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

		if title, ok := result["title"].(string); !ok || title != "New Test Page" {
			t.Errorf("Expected title to be 'New Test Page', got %v", title)
		}
	})

	t.Run("Update Page", func(t *testing.T) {
		updateData := map[string]interface{}{
			"title":   "Updated Test Page",
			"content": "Updated Test Content",
		}

		jsonData, err := json.Marshal(updateData)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest(http.MethodPut, serverURL+"/api/pages/"+testPage.ID(), bytes.NewBuffer(jsonData))
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

		if title, ok := result["title"].(string); !ok || title != "Updated Test Page" {
			t.Errorf("Expected title to be 'Updated Test Page', got %v", title)
		}
	})

	t.Run("Delete Page", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/pages/"+testPage.ID(), nil)
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

		// Verify the page was soft deleted
		page, err := store.PageFindByID(context.Background(), testPage.ID())
		if err != nil {
			t.Fatalf("Failed to find page: %v", err)
		}
		if page == nil {
			t.Fatalf("Page should still exist after soft delete")
		}
		if !page.IsSoftDeleted() {
			t.Errorf("Page should be marked as soft deleted")
		}
	})
}
