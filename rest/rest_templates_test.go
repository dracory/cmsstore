package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
)

type templateTestFixtures struct {
	serverURL    string
	store        cmsstore.StoreInterface
	cleanup      func()
	testSite     cmsstore.SiteInterface
	testTemplate cmsstore.TemplateInterface
}

func setupTemplateTest(t *testing.T) *templateTestFixtures {
	serverURL, store, cleanup := setupTestAPI(t)

	// Create a test site for our tests
	testSite := cmsstore.NewSite()
	testSite.SetName("Test Site for Templates")
	testSite.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err := store.SiteCreate(context.Background(), testSite)
	if err != nil {
		t.Fatalf("Failed to create test site: %v", err)
	}

	// Create a test template for our tests
	testTemplate := cmsstore.NewTemplate()
	testTemplate.SetName("Test Template")
	testTemplate.SetContent("Test Content")
	testTemplate.SetSiteID(testSite.ID())
	testTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	err = store.TemplateCreate(context.Background(), testTemplate)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	return &templateTestFixtures{
		serverURL:    serverURL,
		store:        store,
		cleanup:      cleanup,
		testSite:     testSite,
		testTemplate: testTemplate,
	}
}

func TestListTemplates(t *testing.T) {
	fixtures := setupTemplateTest(t)
	defer fixtures.cleanup()

	resp, err := http.Get(fixtures.serverURL + "/api/templates?site_id=" + fixtures.testSite.ID())
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

	templates, ok := result["templates"].([]interface{})
	if !ok {
		t.Fatalf("Expected templates to be an array, got %T", result["templates"])
	}

	if len(templates) < 1 {
		t.Errorf("Expected at least one template, got %d", len(templates))
	}
}

func TestGetTemplate(t *testing.T) {
	fixtures := setupTemplateTest(t)
	defer fixtures.cleanup()

	resp, err := http.Get(fixtures.serverURL + "/api/templates/" + fixtures.testTemplate.ID())
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

	if id, ok := result["id"].(string); !ok || id != fixtures.testTemplate.ID() {
		t.Errorf("Expected id to be %s, got %v", fixtures.testTemplate.ID(), id)
	}

	if name, ok := result["name"].(string); !ok || name != "Test Template" {
		t.Errorf("Expected name to be 'Test Template', got %v", name)
	}
}

func TestCreateTemplate(t *testing.T) {
	fixtures := setupTemplateTest(t)
	defer fixtures.cleanup()

	templateData := map[string]interface{}{
		"name":    "New Test Template",
		"content": "New Test Content",
		"site_id": fixtures.testSite.ID(),
		"status":  cmsstore.TEMPLATE_STATUS_ACTIVE,
	}

	jsonData, err := json.Marshal(templateData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(fixtures.serverURL+"/api/templates", "application/json", bytes.NewBuffer(jsonData))
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

	if name, ok := result["name"].(string); !ok || name != "New Test Template" {
		t.Errorf("Expected name to be 'New Test Template', got %v", name)
	}
}

func TestUpdateTemplate(t *testing.T) {
	fixtures := setupTemplateTest(t)
	defer fixtures.cleanup()

	updateData := map[string]interface{}{
		"name":    "Updated Test Template",
		"content": "Updated Test Content",
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fixtures.serverURL+"/api/templates/"+fixtures.testTemplate.ID(), bytes.NewBuffer(jsonData))
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

	if name, ok := result["name"].(string); !ok || name != "Updated Test Template" {
		t.Errorf("Expected name to be 'Updated Test Template', got %v", name)
	}
}

func TestDeleteTemplate(t *testing.T) {
	fixtures := setupTemplateTest(t)
	defer fixtures.cleanup()

	req, err := http.NewRequest(http.MethodDelete, fixtures.serverURL+"/api/templates/"+fixtures.testTemplate.ID(), nil)
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

	// Verify the template was soft deleted by querying with soft-deleted included
	list, err := fixtures.store.TemplateList(context.Background(),
		cmsstore.TemplateQuery().
			SetID(fixtures.testTemplate.ID()).
			SetSoftDeletedIncluded(true).
			SetLimit(1))
	if err != nil {
		t.Fatalf("Failed to find template: %v", err)
	}
	if len(list) == 0 {
		t.Fatalf("Template should still exist after soft delete")
	}
	template := list[0]
	if !template.IsSoftDeleted() {
		t.Errorf("Template should be marked as soft deleted")
	}
}
