package admin

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

func initTemplateDeleteHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewTemplateDeleteController(ui).Handler
}

func Test_TemplateDeleteController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and template
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	template := cmsstore.NewTemplate()
	template.SetName("Test Template")
	template.SetSiteID(site.ID())
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetContent("<html><body>Test content</body></html>")
	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"template_id": {template.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "Delete Template") {
		t.Errorf("Expected body to contain 'Delete Template'")
	}
	if !strings.Contains(body, "Are you sure you want to delete this template?") {
		t.Errorf("Expected body to contain confirmation message")
	}
	if !strings.Contains(body, template.ID()) {
		t.Errorf("Expected body to contain template ID")
	}
}

func Test_TemplateDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site and template
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	template := cmsstore.NewTemplate()
	template.SetName("Test Template")
	template.SetSiteID(site.ID())
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetContent("<html><body>Test content</body></html>")
	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {template.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	bodyLower := strings.ToLower(body)
	if !strings.Contains(bodyLower, "success") {
		t.Errorf("Expected body to contain 'success'")
	}
	if !strings.Contains(bodyLower, "template deleted successfully") {
		t.Errorf("Expected body to contain success message")
	}

	// Verify template is deleted (soft delete)
	deletedTemplate, err := store.TemplateFindByID(context.Background(), template.ID())
	if err == nil && deletedTemplate != nil {
		if deletedTemplate.Status() != cmsstore.TEMPLATE_STATUS_INACTIVE {
			t.Errorf("Expected template status to be INACTIVE, got %s", deletedTemplate.Status())
		}
	}
}

func Test_TemplateDeleteController_Delete_ValidationError_MissingTemplateID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(strings.ToLower(body), "error") {
		t.Errorf("Expected body to contain 'error'")
	}
}

func Test_TemplateDeleteController_Delete_ValidationError_EmptyTemplateID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {""},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(strings.ToLower(body), "error") {
		t.Errorf("Expected body to contain 'error'")
	}
}

func Test_TemplateDeleteController_TemplateNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {"non-existent-id"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	bodyLower := strings.ToLower(body)
	if !strings.Contains(bodyLower, "error") {
		t.Errorf("Expected body to contain 'error'")
	}
	if !strings.Contains(bodyLower, "template not found") {
		t.Errorf("Expected body to contain error message")
	}
}

func Test_TemplateDeleteController_Delete_DifferentTemplateTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	// Test HTML template deletion
	htmlTemplate := cmsstore.NewTemplate()
	htmlTemplate.SetName("HTML Template")
	htmlTemplate.SetSiteID(site.ID())
	htmlTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	htmlTemplate.SetContent("<div>HTML template content</div>")
	err = store.TemplateCreate(context.Background(), htmlTemplate)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {htmlTemplate.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(strings.ToLower(body), "success") {
		t.Errorf("Expected body to contain 'success'")
	}

	// Test simple template deletion
	simpleTemplate := cmsstore.NewTemplate()
	simpleTemplate.SetName("Simple Template")
	simpleTemplate.SetSiteID(site.ID())
	simpleTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	simpleTemplate.SetContent("Simple template content")
	err = store.TemplateCreate(context.Background(), simpleTemplate)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {simpleTemplate.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(strings.ToLower(body), "success") {
		t.Errorf("Expected body to contain 'success'")
	}
}

func Test_TemplateDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	// Create multiple templates
	template1 := cmsstore.NewTemplate()
	template1.SetName("Template 1")
	template1.SetSiteID(site.ID())
	template1.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template1.SetContent("<div>Template 1 content</div>")
	err = store.TemplateCreate(context.Background(), template1)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	template2 := cmsstore.NewTemplate()
	template2.SetName("Template 2")
	template2.SetSiteID(site.ID())
	template2.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template2.SetContent("<div>Template 2 content</div>")
	err = store.TemplateCreate(context.Background(), template2)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Delete first template
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {template1.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(strings.ToLower(body), "success") {
		t.Errorf("Expected body to contain 'success'")
	}

	// Verify first template is deleted but second template remains
	deletedTemplate1, err := store.TemplateFindByID(context.Background(), template1.ID())
	if err == nil && deletedTemplate1 != nil {
		if deletedTemplate1.Status() != cmsstore.TEMPLATE_STATUS_INACTIVE {
			t.Errorf("Expected deleted template status to be INACTIVE, got %s", deletedTemplate1.Status())
		}
	}

	activeTemplate2, err := store.TemplateFindByID(context.Background(), template2.ID())
	if err == nil && activeTemplate2 != nil {
		if activeTemplate2.Status() != cmsstore.TEMPLATE_STATUS_ACTIVE {
			t.Errorf("Expected active template status to be ACTIVE, got %s", activeTemplate2.Status())
		}
	}
}

func Test_TemplateDeleteController_Delete_WithContent(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("Failed to seed site: %v", err)
	}

	// Create a template with content
	template := cmsstore.NewTemplate()
	template.SetName("Template with Content")
	template.SetSiteID(site.ID())
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetContent("<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Test Template</h1></body></html>")
	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {template.ID()},
		},
	})
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	bodyLower := strings.ToLower(body)
	if !strings.Contains(bodyLower, "success") {
		t.Errorf("Expected body to contain 'success'")
	}
	if !strings.Contains(bodyLower, "template deleted successfully") {
		t.Errorf("Expected body to contain success message")
	}
}
