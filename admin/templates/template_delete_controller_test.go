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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)

	// Seed a site and template
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	template := cmsstore.NewTemplate()
	template.SetName("Test Template")
	template.SetSiteID(site.ID())
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetContent("<html><body>Test content</body></html>")
	err = store.TemplateCreate(context.Background(), template)
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
		GetValues: map[string][]string{
			"template_id": {template.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Delete Template")
	assert.Contains(t, body, "Are you sure you want to delete this template?")
	assert.Contains(t, body, template.ID())
}

func Test_TemplateDeleteController_Delete(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site and template
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	template := cmsstore.NewTemplate()
	template.SetName("Test Template")
	template.SetSiteID(site.ID())
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetContent("<html><body>Test content</body></html>")
	err = store.TemplateCreate(context.Background(), template)
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {template.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "template deleted successfully")

	// Verify template is deleted (soft delete)
	deletedTemplate, err := store.TemplateFindByID(context.Background(), template.ID())
	if err == nil && deletedTemplate != nil {
		assert.Equal(t, cmsstore.TEMPLATE_STATUS_INACTIVE, deletedTemplate.Status())
	}
}

func Test_TemplateDeleteController_Delete_ValidationError_MissingTemplateID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_TemplateDeleteController_Delete_ValidationError_EmptyTemplateID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {""},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_TemplateDeleteController_TemplateNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {"non-existent-id"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
	assert.Contains(t, strings.ToLower(body), "template not found")
}

func Test_TemplateDeleteController_Delete_DifferentTemplateTypes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	// Test HTML template deletion
	htmlTemplate := cmsstore.NewTemplate()
	htmlTemplate.SetName("HTML Template")
	htmlTemplate.SetSiteID(site.ID())
	htmlTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	htmlTemplate.SetContent("<div>HTML template content</div>")
	err = store.TemplateCreate(context.Background(), htmlTemplate)
	require.NoError(t, err)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {htmlTemplate.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")

	// Test simple template deletion
	simpleTemplate := cmsstore.NewTemplate()
	simpleTemplate.SetName("Simple Template")
	simpleTemplate.SetSiteID(site.ID())
	simpleTemplate.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	simpleTemplate.SetContent("Simple template content")
	err = store.TemplateCreate(context.Background(), simpleTemplate)
	require.NoError(t, err)

	body, response, err = test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {simpleTemplate.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")
}

func Test_TemplateDeleteController_Integration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	// Create multiple templates
	template1 := cmsstore.NewTemplate()
	template1.SetName("Template 1")
	template1.SetSiteID(site.ID())
	template1.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template1.SetContent("<div>Template 1 content</div>")
	err = store.TemplateCreate(context.Background(), template1)
	require.NoError(t, err)

	template2 := cmsstore.NewTemplate()
	template2.SetName("Template 2")
	template2.SetSiteID(site.ID())
	template2.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template2.SetContent("<div>Template 2 content</div>")
	err = store.TemplateCreate(context.Background(), template2)
	require.NoError(t, err)

	// Delete first template
	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {template1.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, strings.ToLower(body), "success")

	// Verify first template is deleted but second template remains
	deletedTemplate1, err := store.TemplateFindByID(context.Background(), template1.ID())
	if err == nil && deletedTemplate1 != nil {
		assert.Equal(t, cmsstore.TEMPLATE_STATUS_INACTIVE, deletedTemplate1.Status())
	}

	activeTemplate2, err := store.TemplateFindByID(context.Background(), template2.ID())
	if err == nil && activeTemplate2 != nil {
		assert.Equal(t, cmsstore.TEMPLATE_STATUS_ACTIVE, activeTemplate2.Status())
	}
}

func Test_TemplateDeleteController_Delete_WithContent(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	// Create a template with content
	template := cmsstore.NewTemplate()
	template.SetName("Template with Content")
	template.SetSiteID(site.ID())
	template.SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)
	template.SetContent("<!DOCTYPE html><html><head><title>Test</title></head><body><h1>Test Template</h1></body></html>")
	err = store.TemplateCreate(context.Background(), template)
	require.NoError(t, err)

	handler := initTemplateDeleteHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_id": {template.ID()},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "template deleted successfully")
}
