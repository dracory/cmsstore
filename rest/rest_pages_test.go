package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gouniverse/cmsstore"
)

func createTestSite(t *testing.T, store cmsstore.StoreInterface) cmsstore.SiteInterface {
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site for Pages - " + t.Name())
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}
	return testSite
}

func createTestPage(t *testing.T, store cmsstore.StoreInterface, siteID string) cmsstore.PageInterface {
	testPage := cmsstore.NewPage()
	testPage.SetTitle("Test Page - " + t.Name())
	testPage.SetContent("Test Content - " + t.Name())
	testPage.SetSiteID(siteID)
	testPage.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err := store.PageCreate(context.Background(), testPage)
	if err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	return testPage
}

func TestPageEndpoints(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	t.Run("List Pages", func(t *testing.T) {
		// Create test data
		testSite := createTestSite(t, store)
		createTestPage(t, store, testSite.ID())

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
		// Create test data
		testSite := createTestSite(t, store)
		testPage := createTestPage(t, store, testSite.ID())

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

		expectedTitle := "Test Page - " + t.Name()
		if title, ok := result["title"].(string); !ok || title != expectedTitle {
			t.Errorf("Expected title to be '%s', got %v", expectedTitle, title)
		}
	})

	t.Run("Create Page", func(t *testing.T) {
		// Create test site first
		testSite := createTestSite(t, store)
		
		pageData := map[string]interface{}{
			"title":   "New Test Page - " + t.Name(),
			"content": "New Test Content - " + t.Name(),
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

		expectedTitle := "New Test Page - " + t.Name()
		if title, ok := result["title"].(string); !ok || title != expectedTitle {
			t.Errorf("Expected title to be '%s', got %v", expectedTitle, title)
		}
	})

	t.Run("Update Page", func(t *testing.T) {
		// Create test data
		testSite := createTestSite(t, store)
		testPage := createTestPage(t, store, testSite.ID())

		updateData := map[string]interface{}{
			"title":   "Updated Test Page - " + t.Name(),
			"content": "Updated Test Content - " + t.Name(),
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

		expectedTitle := "Updated Test Page - " + t.Name()
		if title, ok := result["title"].(string); !ok || title != expectedTitle {
			t.Errorf("Expected title to be '%s', got %v", expectedTitle, title)
		}
	})

	t.Run("Delete Page", func(t *testing.T) {
		// Create test data
		testSite := createTestSite(t, store)
		testPage := createTestPage(t, store, testSite.ID())

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

		// Verify the page was soft deleted by querying with soft-deleted included
		list, err := store.PageList(context.Background(), 
			cmsstore.PageQuery().
				SetID(testPage.ID()).
				SetSoftDeletedIncluded(true).
				SetLimit(1))
		if err != nil {
			t.Fatalf("Failed to find page: %v", err)
		}
		if len(list) == 0 {
			t.Fatalf("Page should still exist after soft delete")
		}
		page := list[0]
		if !page.IsSoftDeleted() {
			t.Errorf("Page should be marked as soft deleted")
		}
	})
}
