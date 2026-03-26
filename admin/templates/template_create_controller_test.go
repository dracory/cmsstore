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

func initTemplateCreateHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewTemplateCreateController(ui).Handler
}

func Test_TemplateCreateController_Index(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "New Template")
	assert.Contains(t, body, "template_name")
	assert.Contains(t, body, "site_id")
}

func Test_TemplateCreateController_Create(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_name":    {"Test Template"},
			"site_id":          {site.ID()},
			"template_content": {"<html><body>Test content</body></html>"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "template created successfully")
}

func Test_TemplateCreateController_Create_ValidationError_MissingName(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"site_id":          {site.ID()},
			"template_content": {"<html><body>Test content</body></html>"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_TemplateCreateController_Create_ValidationError_MissingSiteID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_name":    {"Test Template"},
			"template_content": {"<html><body>Test content</body></html>"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "error")
}

func Test_TemplateCreateController_Create_WithSiteList(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed multiple sites
	site1, err := testutils.SeedSite(store, "Test Site 1")
	require.NoError(t, err)

	site2 := cmsstore.NewSite()
	site2.SetName("Test Site 2")
	site2.SetDomainNames([]string{"site2.example.com"})
	err = store.SiteCreate(context.Background(), site2)
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, body, "Test Site 1")
	assert.Contains(t, body, "Test Site 2")
	assert.Contains(t, body, site1.ID())
	assert.Contains(t, body, site2.ID())
}

func Test_TemplateCreateController_Create_WithMemo(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_name":    {"Test Template"},
			"site_id":          {site.ID()},
			"template_content": {"<html><body>Test content</body></html>"},
			"template_memo":    {"This is a test memo"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
	assert.Contains(t, strings.ToLower(body), "template created successfully")
}

func Test_TemplateCreateController_Create_WithHandle(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_name":    {"Test Template"},
			"site_id":          {site.ID()},
			"template_content": {"<html><body>Test content</body></html>"},
			"template_handle":  {"test-template-handle"},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
}

func Test_TemplateCreateController_Create_EmptyContent(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_name":    {"Test Template"},
			"site_id":          {site.ID()},
			"template_content": {""},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Empty content might be allowed
	assert.Contains(t, strings.ToLower(body), "success")
}

func Test_TemplateCreateController_Create_WithHTMLContent(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	require.NoError(t, err)

	// Seed a site
	site, err := testutils.SeedSite(store, "Test Site")
	require.NoError(t, err)

	handler := initTemplateCreateHandler(store)

	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Test Template</title>
</head>
<body>
    <h1>Hello World</h1>
    <p>This is a test template</p>
</body>
</html>`

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: map[string][]string{
			"template_name":    {"HTML Template"},
			"site_id":          {site.ID()},
			"template_content": {htmlContent},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	assert.Contains(t, strings.ToLower(body), "success")
}
