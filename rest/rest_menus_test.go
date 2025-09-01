package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
)

func TestMenuEndpoints(t *testing.T) {
	serverURL, store, cleanup := setupTestAPI(t)
	defer cleanup()

	// Create a test site for our tests
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site for Menus")
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	// Create a test menu for our tests
	testMenu := cmsstore.NewMenu()
	testMenu.SetName("Test Menu")
	testMenu.SetSiteID(testSite.ID())
	testMenu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(context.Background(), testMenu)
	if err != nil {
		t.Fatalf("Failed to create test menu: %v", err)
	}

	t.Run("List Menus", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/menus?site_id=" + testSite.ID())
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

		menus, ok := result["menus"].([]interface{})
		if !ok {
			t.Fatalf("Expected menus to be an array, got %T", result["menus"])
		}

		if len(menus) < 1 {
			t.Errorf("Expected at least one menu, got %d", len(menus))
		}
	})

	t.Run("Get Menu", func(t *testing.T) {
		resp, err := http.Get(serverURL + "/api/menus/" + testMenu.ID())
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

		if id, ok := result["id"].(string); !ok || id != testMenu.ID() {
			t.Errorf("Expected id to be %s, got %v", testMenu.ID(), id)
		}

		if name, ok := result["name"].(string); !ok || name != "Test Menu" {
			t.Errorf("Expected name to be 'Test Menu', got %v", name)
		}
	})

	t.Run("Create Menu", func(t *testing.T) {
		menuData := map[string]interface{}{
			"name":    "New Test Menu",
			"site_id": testSite.ID(),
			"status":  cmsstore.MENU_STATUS_ACTIVE,
		}

		jsonData, err := json.Marshal(menuData)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		resp, err := http.Post(serverURL+"/api/menus", "application/json", bytes.NewBuffer(jsonData))
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

		if name, ok := result["name"].(string); !ok || name != "New Test Menu" {
			t.Errorf("Expected name to be 'New Test Menu', got %v", name)
		}
	})

	t.Run("Update Menu", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name": "Updated Test Menu",
		}

		jsonData, err := json.Marshal(updateData)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v", err)
		}

		req, err := http.NewRequest(http.MethodPut, serverURL+"/api/menus/"+testMenu.ID(), bytes.NewBuffer(jsonData))
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

		if name, ok := result["name"].(string); !ok || name != "Updated Test Menu" {
			t.Errorf("Expected name to be 'Updated Test Menu', got %v", name)
		}
	})

	t.Run("Delete Menu", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodDelete, serverURL+"/api/menus/"+testMenu.ID(), nil)
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

		// Verify the menu was soft deleted by querying with soft-deleted included
		list, err := store.MenuList(context.Background(),
			cmsstore.MenuQuery().
				SetID(testMenu.ID()).
				SetSoftDeletedIncluded(true).
				SetLimit(1))
		if err != nil {
			t.Fatalf("Failed to find menu: %v", err)
		}
		if len(list) == 0 {
			t.Fatalf("Menu should still exist after soft delete")
		}
		menu := list[0]
		if !menu.IsSoftDeleted() {
			t.Errorf("Menu should be marked as soft deleted")
		}
	})
}
