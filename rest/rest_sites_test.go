package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gouniverse/cmsstore"
)

func TestSiteEndpoints(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	// Create a test site for our tests
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site")
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	_, err := testSite.SetDomainNames([]string{"example.com", "www.example.com"})
	if err != nil {
		t.Fatalf("Failed to set domain names: %v", err)
	}

	err = store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	t.Run("List Sites", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/sites")
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

		sites, ok := result["sites"].([]interface{})
		if !ok {
			t.Fatalf("Expected sites to be an array, got %T", result["sites"])
		}

		if len(sites) < 1 {
			t.Errorf("Expected at least one site, got %d", len(sites))
		}
	})

	t.Run("Get Site", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/sites/" + testSite.ID())
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

		if id, ok := result["id"].(string); !ok || id != testSite.ID() {
			t.Errorf("Expected id to be %s, got %v", testSite.ID(), id)
		}

		if name, ok := result["name"].(string); !ok || name != "Test Site" {
			t.Errorf("Expected name to be 'Test Site', got %v", name)
		}
	})

	t.Run("Create Site", func(t *testing.T) {
		siteData := map[string]interface{}{
			"name":         "New Test Site",
			"status":       cmsstore.SITE_STATUS_ACTIVE,
			"domain_names": []string{"newsite.com"},
		}

		jsonData, err := json.Marshal(siteData)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		resp, err := http.Post(serverURL+"/api/sites", "application/json", bytes.NewBuffer(jsonData))
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

		if name, ok := result["name"].(string); !ok || name != "New Test Site" {
			t.Errorf("Expected name to be 'New Test Site', got %v", name)
		}
	})

	t.Run("Update Site", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name":         "Updated Test Site",
			"domain_names": []string{"updated.com"},
		}

		jsonData, err := json.Marshal(updateData)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest(http.MethodPut, serverURL+"/api/sites/"+testSite.ID(), bytes.NewBuffer(jsonData))
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

		if name, ok := result["name"].(string); !ok || name != "Updated Test Site" {
			t.Errorf("Expected name to be 'Updated Test Site', got %v", name)
		}
	})

	t.Run("Delete Site", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/sites/"+testSite.ID(), nil)
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

		// Verify the site was soft deleted
		// Create a query that includes soft-deleted items
		query := cmsstore.SiteQuery().SetSoftDeletedIncluded(true).SetID(testSite.ID())
		sites, err := store.SiteList(context.Background(), query)
		if err != nil {
			t.Fatalf("Failed to find site: %v", err)
		}
		if len(sites) == 0 {
			t.Fatalf("Site should still exist after soft delete")
		}
		site := sites[0]
		if !site.IsSoftDeleted() {
			t.Errorf("Site should be marked as soft deleted")
		}
	})
}
