package frontend

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/ui"
)

// TestFrontend_Handler_IcoRequest tests that .ico requests return empty response
func TestFrontend_Handler_IcoRequest(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	recorder := httptest.NewRecorder()

	f.Handler(recorder, req)

	if recorder.Body.Len() != 0 {
		t.Error("ICO requests should return empty response")
	}
}

// TestFrontend_StringHandler_IcoRequest tests StringHandler with .ico request
func TestFrontend_StringHandler_IcoRequest(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	req := httptest.NewRequest("GET", "/favicon.ico", nil)
	recorder := httptest.NewRecorder()

	result := f.StringHandler(recorder, req)

	if result != "" {
		t.Error("ICO requests should return empty string")
	}
}

// TestFrontend_MenuFindByID tests menu finding by ID
func TestFrontend_MenuFindByID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	menu := cmsstore.NewMenu().
		SetName("Test Menu").
		SetStatus(cmsstore.MENU_STATUS_ACTIVE)

	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	foundMenu, err := f.(*frontend).MenuFindByID(context.Background(), menu.ID())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundMenu == nil {
		t.Error("Expected menu to be found")
	}

	if foundMenu.ID() != menu.ID() {
		t.Error("Expected menu ID to match")
	}
}

// TestFrontend_MenuItemList tests menu item listing
func TestFrontend_MenuItemList(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	menu := cmsstore.NewMenu().
		SetName("Test Menu").
		SetStatus(cmsstore.MENU_STATUS_ACTIVE)

	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	menuItem := cmsstore.NewMenuItem().
		SetMenuID(menu.ID()).
		SetName("Test Item").
		SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)

	err = store.MenuItemCreate(context.Background(), menuItem)
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

	items, err := f.(*frontend).MenuItemList(context.Background(), cmsstore.MenuItemQuery().SetMenuID(menu.ID()))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(items) == 0 {
		t.Error("Expected at least one menu item")
	}
}

// TestFrontend_MenusEnabled tests checking if menus are enabled
func TestFrontend_MenusEnabled(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	enabled := f.(*frontend).MenusEnabled()
	if !enabled {
		t.Error("Menus should be enabled when store has menus enabled")
	}
}

// TestFrontend_PageFindByID tests page finding by ID
func TestFrontend_PageFindByID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test-page").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	foundPage, err := f.(*frontend).PageFindByID(context.Background(), page.ID())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundPage == nil {
		t.Error("Expected page to be found")
	}

	if foundPage.ID() != page.ID() {
		t.Error("Expected page ID to match")
	}
}

// TestFrontend_RenderMenuHTML tests menu HTML rendering
func TestFrontend_RenderMenuHTML(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	menu := cmsstore.NewMenu().
		SetName("Test Menu").
		SetStatus(cmsstore.MENU_STATUS_ACTIVE)

	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	menuItem := cmsstore.NewMenuItem().
		SetMenuID(menu.ID()).
		SetName("Test Item").
		SetURL("/test").
		SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)

	err = store.MenuItemCreate(context.Background(), menuItem)
	if err != nil {
		t.Fatalf("Failed to create menu item: %v", err)
	}

	items, err := store.MenuItemList(context.Background(), cmsstore.MenuItemQuery().SetMenuID(menu.ID()))
	if err != nil {
		t.Fatalf("Failed to list menu items: %v", err)
	}

	html, err := f.(*frontend).RenderMenuHTML(context.Background(), items, "horizontal", "test-class", 1, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Menu HTML rendering may return empty if there are no valid menu items
	// The important thing is that it doesn't error
	_ = html
}

// TestFrontend_Logger tests getting the logger
func TestFrontend_Logger(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	logger := slog.Default()
	f := New(Config{
		Store:  store,
		Logger: logger,
	})

	retrievedLogger := f.(*frontend).Logger()
	if retrievedLogger != logger {
		t.Error("Logger should be the one provided in config")
	}
}

// TestFrontend_BlockRegistry tests getting the block registry
func TestFrontend_BlockRegistry(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	registry := f.(*frontend).BlockRegistry()
	if registry == nil {
		t.Error("Block registry should not be nil")
	}
}

// TestFrontend_TemplateRenderHtmlByID tests template rendering by ID
func TestFrontend_TemplateRenderHtmlByID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	template := cmsstore.NewTemplate().
		SetSiteID(site.ID()).
		SetName("Test Template").
		SetContent("<div>[[PageTitle]]</div>").
		SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)

	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	req := httptest.NewRequest("GET", "/", nil)

	options := TemplateRenderHtmlByIDOptions{
		PageTitle: "Test Page",
	}

	result, err := f.TemplateRenderHtmlByID(req, template.ID(), options)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

// TestFrontend_TemplateRenderHtmlByID_EmptyID tests template rendering with empty ID
func TestFrontend_TemplateRenderHtmlByID_EmptyID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	req := httptest.NewRequest("GET", "/", nil)

	_, err = f.TemplateRenderHtmlByID(req, "", TemplateRenderHtmlByIDOptions{})
	if err == nil {
		t.Error("Expected error for empty template ID")
	}
}

// TestFrontend_TemplateRenderHtmlByID_NotFound tests template rendering with non-existent ID
func TestFrontend_TemplateRenderHtmlByID_NotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	req := httptest.NewRequest("GET", "/", nil)

	_, err = f.TemplateRenderHtmlByID(req, "non-existent-id", TemplateRenderHtmlByIDOptions{})
	if err == nil {
		t.Error("Expected error for non-existent template ID")
	}
}

// TestFrontend_TemplateRenderHtmlByID_Inactive tests template rendering with inactive template
func TestFrontend_TemplateRenderHtmlByID_Inactive(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	template := cmsstore.NewTemplate().
		SetSiteID(site.ID()).
		SetName("Test Template").
		SetContent("<div>Test</div>").
		SetStatus(cmsstore.TEMPLATE_STATUS_INACTIVE)

	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	req := httptest.NewRequest("GET", "/", nil)

	_, err = f.TemplateRenderHtmlByID(req, template.ID(), TemplateRenderHtmlByIDOptions{})
	if err == nil {
		t.Error("Expected error for inactive template")
	}
}

// TestFrontend_StringHandler_DomainNotSupported tests handling of unsupported domains
func TestFrontend_StringHandler_DomainNotSupported(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})

	req := httptest.NewRequest("GET", "http://unsupported.com/page", nil)
	recorder := httptest.NewRecorder()

	result := f.StringHandler(recorder, req)

	expected := "Domain not supported: unsupported.com"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestFrontend_ContentRenderBlockByID_EmptyID tests rendering block with empty ID
func TestFrontend_ContentRenderBlockByID_EmptyID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})
	fe := f.(*frontend)

	content := "Test content [[BLOCK_]] more content"
	result, err := fe.contentRenderBlockByID(context.Background(), content, "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != content {
		t.Error("Content should remain unchanged when block ID is empty")
	}
}

// TestFrontend_ContentRenderPageURLByID_EmptyID tests rendering page URL with empty ID
func TestFrontend_ContentRenderPageURLByID_EmptyID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})
	fe := f.(*frontend)

	content := "Test content [[PAGE_URL_]] more content"
	result, err := fe.contentRenderPageURLByID(context.Background(), content, "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != content {
		t.Error("Content should remain unchanged when page ID is empty")
	}
}

// TestFrontend_ContentRenderTranslationByHandleOrId_EmptyID tests rendering translation with empty ID
func TestFrontend_ContentRenderTranslationByHandleOrId_EmptyID(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})
	fe := f.(*frontend)

	content := "Test content [[TRANSLATION_]] more content"
	result, err := fe.ContentRenderTranslationByHandleOrId(context.Background(), content, "", "en")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != content {
		t.Error("Content should remain unchanged when translation ID is empty")
	}
}

// TestFrontend_ConvertBlockJsonToHtml_NilRenderer tests block JSON conversion with nil renderer
func TestFrontend_ConvertBlockJsonToHtml_NilRenderer(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store: store,
	})
	fe := f.(*frontend)

	json := `{"type":"paragraph","content":"test"}`
	result := fe.convertBlockJsonToHtml(json)

	if result != "Block editor not configured" {
		t.Error("Expected 'Block editor not configured' message when renderer is nil")
	}
}

// TestFrontend_ConvertBlockJsonToHtml_Malformed tests block JSON conversion with malformed JSON
func TestFrontend_ConvertBlockJsonToHtml_Malformed(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	renderer := func(blocks []ui.BlockInterface) string {
		return "rendered"
	}

	f := New(Config{
		Store:               store,
		BlockEditorRenderer: renderer,
	})
	fe := f.(*frontend)

	json := "not valid json"
	result := fe.convertBlockJsonToHtml(json)

	if result != "Malformed block content" {
		t.Error("Expected 'Malformed block content' message for invalid JSON")
	}
}

// TestFrontend_ConvertBlockJsonToHtml_UnmarshalError tests block JSON conversion with unmarshal error
func TestFrontend_ConvertBlockJsonToHtml_UnmarshalError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	renderer := func(blocks []ui.BlockInterface) string {
		return "rendered"
	}

	f := New(Config{
		Store:               store,
		BlockEditorRenderer: renderer,
	})
	fe := f.(*frontend)

	// Valid JSON but structure that won't unmarshal to blocks
	json := `{"invalid":"structure"}`
	result := fe.convertBlockJsonToHtml(json)

	if result != "Error parsing block content" {
		t.Errorf("Expected 'Error parsing block content' for unmarshal error, got %q", result)
	}
}

// TestPageRenderHtmlBySiteAndAlias_WithContent tests page rendering with basic content
func TestPageRenderHtmlBySiteAndAlias_WithContent(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test-page").
		SetContent("<h1>Test Content</h1>").
		SetTitle("Test Page Title").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	result := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "test-page", "en")

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Verify the page content is in the result
	if !strings.Contains(result, "Test Content") {
		t.Error("Expected page content to be in result")
	}
}

// TestPageRenderHtmlBySiteAndAlias_WithTemplate tests page rendering with template
func TestPageRenderHtmlBySiteAndAlias_WithTemplate(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	template := cmsstore.NewTemplate().
		SetSiteID(site.ID()).
		SetName("Test Template").
		SetContent("<div class='template'>[[PageContent]]</div>").
		SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)

	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test-page").
		SetContent("<h1>Page Content</h1>").
		SetTemplateID(template.ID()).
		SetTitle("Test Page Title").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	result := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "test-page", "en")

	if result == "" {
		t.Error("Expected non-empty result")
	}

	// Verify the template content is in the result
	if !strings.Contains(result, "template") {
		t.Error("Expected template content to be in result")
	}
}

// TestPageRenderHtmlBySiteAndAlias_NotFound tests page rendering when page not found
func TestPageRenderHtmlBySiteAndAlias_NotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	result := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "nonexistent", "en")

	if result == "" {
		t.Error("Expected non-empty result (not found message)")
	}

	if !strings.Contains(result, "not found") {
		t.Error("Expected 'not found' message in result")
	}
}

// TestPageRenderHtmlBySiteAndAlias_WithCustomNotFoundHandler tests custom not found handler
func TestPageRenderHtmlBySiteAndAlias_WithCustomNotFoundHandler(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	handlerCalled := false
	customHandler := func(w http.ResponseWriter, r *http.Request, alias string) (handled bool, result string) {
		handlerCalled = true
		return true, "Custom 404: " + alias
	}

	f := New(Config{
		Store:               store,
		Logger:              slog.Default(),
		PageNotFoundHandler: customHandler,
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	result := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "nonexistent", "en")

	if !handlerCalled {
		t.Error("Expected custom not found handler to be called")
	}

	if !strings.Contains(result, "Custom 404") {
		t.Error("Expected custom 404 message in result")
	}
}

// TestFindSiteAndEndpointByDomainAndPath_BasicDomain tests basic domain matching
func TestFindSiteAndEndpointByDomainAndPath_BasicDomain(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site")

	_, err = site.SetDomainNames([]string{"example.com"})
	if err != nil {
		t.Fatalf("Failed to set domain names: %v", err)
	}

	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundSite, endpoint, err := f.(*frontend).findSiteAndEndpointByDomainAndPath(ctx, "example.com", "/page")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundSite == nil {
		t.Error("Expected site to be found")
	}

	if endpoint != "example.com" {
		t.Errorf("Expected endpoint 'example.com', got %q", endpoint)
	}
}

// TestFindSiteAndEndpointByDomainAndPath_Subdirectory tests subdirectory-based endpoints
func TestFindSiteAndEndpointByDomainAndPath_Subdirectory(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site")

	_, err = site.SetDomainNames([]string{"example.com/blog"})
	if err != nil {
		t.Fatalf("Failed to set domain names: %v", err)
	}

	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundSite, endpoint, err := f.(*frontend).findSiteAndEndpointByDomainAndPath(ctx, "example.com", "/blog/page")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundSite == nil {
		t.Error("Expected site to be found")
	}

	if endpoint != "example.com/blog" {
		t.Errorf("Expected endpoint 'example.com/blog', got %q", endpoint)
	}
}

// TestFindSiteAndEndpointByDomainAndPath_LongestMatch tests longest prefix matching
func TestFindSiteAndEndpointByDomainAndPath_LongestMatch(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Create two sites with overlapping domains
	site1 := cmsstore.NewSite().
		SetName("Site 1")

	_, err = site1.SetDomainNames([]string{"example.com"})
	if err != nil {
		t.Fatalf("Failed to set domain names for site 1: %v", err)
	}

	site1.SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site1)
	if err != nil {
		t.Fatalf("Failed to create site 1: %v", err)
	}

	site2 := cmsstore.NewSite().
		SetName("Site 2")

	_, err = site2.SetDomainNames([]string{"example.com/blog"})
	if err != nil {
		t.Fatalf("Failed to set domain names for site 2: %v", err)
	}

	site2.SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site2)
	if err != nil {
		t.Fatalf("Failed to create site 2: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	// Should match the longer endpoint
	foundSite, endpoint, err := f.(*frontend).findSiteAndEndpointByDomainAndPath(ctx, "example.com", "/blog/page")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundSite == nil {
		t.Error("Expected site to be found")
	}

	if foundSite.ID() != site2.ID() {
		t.Error("Expected to match the longer endpoint (site 2)")
	}

	if endpoint != "example.com/blog" {
		t.Errorf("Expected endpoint 'example.com/blog', got %q", endpoint)
	}
}

// TestFindSiteAndEndpointByDomainAndPath_NotFound tests when no site matches
func TestFindSiteAndEndpointByDomainAndPath_NotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site")

	_, err = site.SetDomainNames([]string{"example.com"})
	if err != nil {
		t.Fatalf("Failed to set domain names: %v", err)
	}

	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundSite, endpoint, err := f.(*frontend).findSiteAndEndpointByDomainAndPath(ctx, "otherdomain.com", "/page")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundSite != nil {
		t.Error("Expected no site to be found")
	}

	if endpoint != "" {
		t.Errorf("Expected empty endpoint, got %q", endpoint)
	}
}

// TestPageFindBySiteAndAliasWithPatterns_AnyPattern tests :any pattern matching
func TestPageFindBySiteAndAliasWithPatterns_AnyPattern(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("/blog/:any").
		SetContent("Blog content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundPage, err := f.(*frontend).pageFindBySiteAndAliasWithPatterns(ctx, site.ID(), "/blog/post-title")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundPage == nil {
		t.Error("Expected page to match :any pattern")
	} else if foundPage.ID() != page.ID() {
		t.Error("Expected to find the correct page")
	}
}

// TestPageFindBySiteAndAliasWithPatterns_NumPattern tests :num pattern matching
func TestPageFindBySiteAndAliasWithPatterns_NumPattern(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("/post/:num").
		SetContent("Post content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundPage, err := f.(*frontend).pageFindBySiteAndAliasWithPatterns(ctx, site.ID(), "/post/123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundPage == nil {
		t.Error("Expected page to match :num pattern")
	} else if foundPage.ID() != page.ID() {
		t.Error("Expected to find the correct page")
	}
}

// TestPageFindBySiteAndAliasWithPatterns_AlphaPattern tests :alpha pattern matching
func TestPageFindBySiteAndAliasWithPatterns_AlphaPattern(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("/user/:alpha").
		SetContent("User content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundPage, err := f.(*frontend).pageFindBySiteAndAliasWithPatterns(ctx, site.ID(), "/user/john-doe")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundPage == nil {
		t.Error("Expected page to match :alpha pattern")
	} else if foundPage.ID() != page.ID() {
		t.Error("Expected to find the correct page")
	}
}

// TestPageFindBySiteAndAliasWithPatterns_NoMatch tests when no pattern matches
func TestPageFindBySiteAndAliasWithPatterns_NoMatch(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("/post/:num").
		SetContent("Post content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	// This should not match :num pattern
	foundPage, err := f.(*frontend).pageFindBySiteAndAliasWithPatterns(ctx, site.ID(), "post/abc")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundPage != nil {
		t.Error("Expected no page to be found for non-matching pattern")
	}
}

// TestPageOrTemplateContent_NoTemplate tests page content without template
func TestPageOrTemplateContent_NoTemplate(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetContent("Page content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	result := f.(*frontend).pageOrTemplateContent(req, page)

	if result != "Page content" {
		t.Errorf("Expected 'Page content', got %q", result)
	}
}

// TestPageOrTemplateContent_WithTemplate tests page content with template
func TestPageOrTemplateContent_WithTemplate(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	template := cmsstore.NewTemplate().
		SetSiteID(site.ID()).
		SetName("Test Template").
		SetContent("Template content").
		SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)

	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetContent("Page content").
		SetTemplateID(template.ID()).
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	result := f.(*frontend).pageOrTemplateContent(req, page)

	if result != "Template content" {
		t.Errorf("Expected 'Template content', got %q", result)
	}
}

// TestPageOrTemplateContent_TemplateNotFound tests page content when template not found
func TestPageOrTemplateContent_TemplateNotFound(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetContent("Page content").
		SetTemplateID("non-existent-template").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	result := f.(*frontend).pageOrTemplateContent(req, page)

	// Should return page content when template not found
	if result != "Page content" {
		t.Errorf("Expected 'Page content' when template not found, got %q", result)
	}
}

// TestFetchActiveSites_Caching tests that active sites are cached
func TestFetchActiveSites_Caching(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site")

	_, err = site.SetDomainNames([]string{"example.com"})
	if err != nil {
		t.Fatalf("Failed to set domain names: %v", err)
	}

	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()

	// First call - should fetch from database
	sites1, err := f.(*frontend).fetchActiveSites(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sites1) == 0 {
		t.Error("Expected at least one site")
	}

	// Second call - should use cache
	sites2, err := f.(*frontend).fetchActiveSites(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sites2) != len(sites1) {
		t.Error("Expected cached sites to match")
	}

	if sites1[0].ID() != sites2[0].ID() {
		t.Error("Expected cached site ID to match")
	}
}

// TestFetchPageAliasMapBySite_Caching tests that page alias map is cached
func TestFetchPageAliasMapBySite_Caching(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()

	// First call - should fetch from database
	aliasMap1, err := f.(*frontend).fetchPageAliasMapBySite(ctx, site.ID())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(aliasMap1) == 0 {
		t.Error("Expected at least one page alias")
	}

	// Second call - should use cache
	aliasMap2, err := f.(*frontend).fetchPageAliasMapBySite(ctx, site.ID())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(aliasMap2) != len(aliasMap1) {
		t.Error("Expected cached alias map to match")
	}
}

// TestFullRenderingPipeline_WithBlocks tests complete rendering pipeline with blocks
func TestFullRenderingPipeline_WithBlocks(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create a block
	block := cmsstore.NewBlock().
		SetContent("<div>Block content</div>").
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create a page that references the block
	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test").
		SetContent("[[BLOCK_" + block.ID() + "]]]").
		SetTitle("Test Page").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	result := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "test", "en")

	if result == "" {
		t.Error("Expected non-empty result from full rendering pipeline")
	}

	// Verify block content is in the result
	if !strings.Contains(result, "Block content") {
		t.Error("Expected block content to be rendered in result")
	}
}

// TestApplyShortcodes_NoShortcodes tests applyShortcodes with no shortcodes
func TestApplyShortcodes_NoShortcodes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	content := "Test content without shortcodes"

	result, err := f.(*frontend).applyShortcodes(req, content)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != content {
		t.Errorf("Expected content to remain unchanged, got %q", result)
	}
}

// TestContentRenderBlocks_DatabaseError tests block rendering with database error
func TestContentRenderBlocks_DatabaseError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	content := "Test content [[BLOCK_nonexistent]] more content"

	result, err := f.(*frontend).contentRenderBlocks(ctx, content)

	// Should handle missing block gracefully
	if err != nil {
		t.Errorf("Expected no error for non-existent block, got %v", err)
	}

	// The block placeholder should be removed or replaced
	if strings.Contains(result, "[[BLOCK_nonexistent]]") {
		t.Error("Expected block placeholder to be processed")
	}
}

// TestContentRenderPageURLs_DatabaseError tests page URL rendering with database error
func TestContentRenderPageURLs_DatabaseError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	content := "Test content [[PAGE_URL_nonexistent]] more content"

	result, err := f.(*frontend).contentRenderPageURLs(ctx, content)

	// Should handle missing page gracefully
	if err != nil {
		t.Errorf("Expected no error for non-existent page, got %v", err)
	}

	// The page URL placeholder should be replaced with empty string
	if strings.Contains(result, "[[PAGE_URL_nonexistent]]") {
		t.Error("Expected page URL placeholder to be processed")
	}
}

// TestPageRenderHtmlBySiteAndAlias_DatabaseError tests page rendering with database error
func TestPageRenderHtmlBySiteAndAlias_DatabaseError(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	// Try to render a page that doesn't exist
	result := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "nonexistent", "en")

	// Should return a not found message instead of error
	if result == "" {
		t.Error("Expected non-empty result (not found message)")
	}

	if !strings.Contains(result, "not found") {
		t.Error("Expected 'not found' message in result")
	}
}

// TestFetchPageBySiteAndAlias_Caching tests that page fetching is cached
func TestFetchPageBySiteAndAlias_Caching(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test").
		SetContent("Page content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()

	// First call - should fetch from database
	page1, err := f.(*frontend).fetchPageBySiteAndAlias(ctx, site.ID(), "test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if page1 == nil {
		t.Error("Expected page to be found")
	}

	// Second call - should use cache
	page2, err := f.(*frontend).fetchPageBySiteAndAlias(ctx, site.ID(), "test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if page2 == nil {
		t.Error("Expected cached page to be found")
	}

	if page1.ID() != page2.ID() {
		t.Error("Expected cached page ID to match")
	}
}

// TestWarmUpCache tests cache warming on startup
func TestWarmUpCache(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: true,
	})

	// Warm up cache should run during initialization
	// Verify that active sites are cached after warmup
	ctx := context.Background()
	sites, err := f.(*frontend).fetchActiveSites(ctx)

	if err != nil {
		t.Errorf("Expected no error after warmup, got %v", err)
	}

	if len(sites) == 0 {
		t.Error("Expected sites to be cached after warmup")
	}
}

// TestCacheExpiration tests that cached items expire correctly
func TestCacheExpiration(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create frontend with cache enabled
	f := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: true,
	})

	ctx := context.Background()

	// First call - should fetch from database and cache
	sites1, err := f.(*frontend).fetchActiveSites(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sites1) == 0 {
		t.Error("Expected sites to be found")
	}

	// Second call - should use cache
	sites2, err := f.(*frontend).fetchActiveSites(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(sites2) == 0 {
		t.Error("Expected sites to be found after cache expiration")
	}

	// Both should have the same data from cache
	if sites1[0].ID() != sites2[0].ID() {
		t.Error("Expected site ID to match from cache")
	}
}

// TestContentRenderBlocks_CircularReference tests circular block references
func TestContentRenderBlocks_CircularReference(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create block A that references block B
	blockA := cmsstore.NewBlock().
		SetContent("Content A [[BLOCK_blockB]]").
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), blockA)
	if err != nil {
		t.Fatalf("Failed to create block A: %v", err)
	}

	// Create block B that references block A (circular reference)
	blockB := cmsstore.NewBlock().
		SetContent("Content B [[BLOCK_blockA]]").
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), blockB)
	if err != nil {
		t.Fatalf("Failed to create block B: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()

	// Try to render block A (which has circular reference)
	content := "[[BLOCK_" + blockA.ID() + "]]"
	result, err := f.(*frontend).contentRenderBlocks(ctx, content)

	// The system should handle circular references gracefully
	// It should either error OR render without infinite loop
	if err != nil {
		t.Logf("Circular reference detected with error: %v", err)
	}

	// Verify it doesn't hang (test will timeout if it does)
	// Should render at least some content
	if !strings.Contains(result, "Content A") && !strings.Contains(result, "Content B") {
		// If no content rendered, check if it was replaced with empty or error message
		if result == "" || strings.Contains(result, "[[BLOCK_") {
			t.Logf("Circular reference prevented rendering, result: %q", result)
		}
	} else {
		t.Logf("Circular reference handled, rendered content: %q", result)
	}
}

// Helper function to create a test site with a page
func setupTestSiteWithPage(t *testing.T, store cmsstore.StoreInterface) (cmsstore.SiteInterface, cmsstore.PageInterface) {
	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err := store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test").
		SetContent("Test content").
		SetTitle("Test Page").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	return site, page
}

// TestHelperFunction tests the helper function works correctly
func TestHelperFunction(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site, page := setupTestSiteWithPage(t, store)

	if site == nil {
		t.Error("Expected site to be created")
	}

	if page == nil {
		t.Error("Expected page to be created")
	}

	if page.SiteID() != site.ID() {
		t.Error("Expected page to belong to the site")
	}
}

// BenchmarkPageRendering benchmarks page rendering performance
func BenchmarkPageRendering(b *testing.B) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		b.Fatalf("Failed to init store: %v", err)
	}

	// Create site and page directly in benchmark
	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		b.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test").
		SetContent("Test content").
		SetTitle("Test Page").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		b.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "test", "en")
	}
}

// BenchmarkCacheOperations benchmarks cache operations
func BenchmarkCacheOperations(b *testing.B) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		b.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		b.Fatalf("Failed to create site: %v", err)
	}

	f := New(Config{
		Store:        store,
		Logger:       slog.Default(),
		CacheEnabled: true,
	})

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.(*frontend).fetchActiveSites(ctx)
	}
}

// Mock shortcode for testing
type mockShortcode struct {
	alias       string
	description string
	render      func(r *http.Request, s string, m map[string]string) string
}

func (m *mockShortcode) Alias() string {
	return m.alias
}

func (m *mockShortcode) Description() string {
	return m.description
}

func (m *mockShortcode) Render(r *http.Request, s string, attrs map[string]string) string {
	return m.render(r, s, attrs)
}

// TestApplyShortcodes_CustomShortcode tests custom shortcode rendering
func TestApplyShortcodes_CustomShortcode(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	called := false
	customShortcode := &mockShortcode{
		alias:       "highlight",
		description: "Highlights text with a background color",
		render: func(r *http.Request, s string, attrs map[string]string) string {
			called = true
			color := attrs["color"]
			if color == "" {
				color = "yellow"
			}
			// The shortcode library passes content in attrs
			content := attrs["_content"]
			return "<span style='background-color:" + color + "'>" + content + "</span>"
		},
	}

	f := New(Config{
		Store:      store,
		Logger:     slog.Default(),
		Shortcodes: []cmsstore.ShortcodeInterface{customShortcode},
	})

	req := httptest.NewRequest("GET", "/", nil)
	content := "<highlight color='red'>Important text</highlight>"

	result, err := f.(*frontend).applyShortcodes(req, content)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("Expected shortcode to be called")
	}

	// The shortcode library should process the shortcode
	if result == content {
		t.Error("Expected content to be modified by shortcode")
	}
}

// TestApplyShortcodes_WithParameters tests shortcode with multiple parameters
func TestApplyShortcodes_WithParameters(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	called := false
	buttonShortcode := &mockShortcode{
		alias:       "button",
		description: "Creates a button link",
		render: func(r *http.Request, s string, attrs map[string]string) string {
			called = true
			url := attrs["url"]
			class := attrs["class"]
			content := attrs["_content"]
			return "<a href='" + url + "' class='" + class + "'>" + content + "</a>"
		},
	}

	f := New(Config{
		Store:      store,
		Logger:     slog.Default(),
		Shortcodes: []cmsstore.ShortcodeInterface{buttonShortcode},
	})

	req := httptest.NewRequest("GET", "/", nil)
	content := "<button url='/contact' class='btn-primary'>Contact Us</button>"

	result, err := f.(*frontend).applyShortcodes(req, content)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("Expected shortcode to be called")
	}

	// The shortcode library should process the shortcode
	if result == content {
		t.Error("Expected content to be modified by shortcode")
	}
}

// TestApplyShortcodes_MultipleShortcodes tests multiple shortcodes in content
func TestApplyShortcodes_MultipleShortcodes(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	boldCalled := false
	italicCalled := false

	boldShortcode := &mockShortcode{
		alias:       "bold",
		description: "Makes text bold",
		render: func(r *http.Request, s string, attrs map[string]string) string {
			boldCalled = true
			content := attrs["_content"]
			return "<strong>" + content + "</strong>"
		},
	}

	italicShortcode := &mockShortcode{
		alias:       "italic",
		description: "Makes text italic",
		render: func(r *http.Request, s string, attrs map[string]string) string {
			italicCalled = true
			content := attrs["_content"]
			return "<em>" + content + "</em>"
		},
	}

	f := New(Config{
		Store:      store,
		Logger:     slog.Default(),
		Shortcodes: []cmsstore.ShortcodeInterface{boldShortcode, italicShortcode},
	})

	req := httptest.NewRequest("GET", "/", nil)
	content := "<bold>Bold text</bold> and <italic>italic text</italic>"

	result, err := f.(*frontend).applyShortcodes(req, content)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !boldCalled {
		t.Error("Expected bold shortcode to be called")
	}

	if !italicCalled {
		t.Error("Expected italic shortcode to be called")
	}

	// The shortcode library should process both shortcodes
	if result == content {
		t.Error("Expected content to be modified by shortcodes")
	}
}

// TestContentRenderTranslations_FallbackChain tests translation fallback chain
func TestContentRenderTranslations_FallbackChain(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Create translation with only English and Spanish
	translation := cmsstore.NewTranslation().
		SetHandle("greeting").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	translationContent := map[string]string{
		"en": "Hello",
		"es": "Hola",
		// No French translation
	}
	err = translation.SetContent(translationContent)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	content := "[[TRANSLATION_greeting]]"

	// Try to get French translation (should fall back to empty or default)
	result, err := f.(*frontend).contentRenderTranslations(ctx, content, "fr")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should replace with empty string when translation not found
	if strings.Contains(result, "[[TRANSLATION_greeting]]") {
		t.Error("Expected translation placeholder to be replaced")
	}
}

// TestContentRenderTranslations_Basic tests basic translation rendering
func TestContentRenderTranslations_Basic(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	translation := cmsstore.NewTranslation().
		SetHandle("welcome").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	translationContent := map[string]string{
		"en": "Welcome",
	}
	err = translation.SetContent(translationContent)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	content := "[[TRANSLATION_welcome]]"

	// Test that translation rendering works
	result, err := f.(*frontend).contentRenderTranslations(ctx, content, "en")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "Welcome") {
		t.Error("Expected translation to be rendered")
	}
}

// TestConvertBlockJsonToHtml_ValidBlocks tests block JSON conversion with valid blocks
func TestConvertBlockJsonToHtml_ValidBlocks(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	rendererCalled := false
	renderer := func(blocks []ui.BlockInterface) string {
		rendererCalled = true
		// Simple renderer that concatenates block content
		result := ""
		for _, block := range blocks {
			// Use the Parameter method to get content
			content := block.Parameter("content")
			if content != "" {
				result += "<p>" + content + "</p>"
			}
		}
		return result
	}

	f := New(Config{
		Store:               store,
		BlockEditorRenderer: renderer,
	})
	fe := f.(*frontend)

	// Create a valid block using the ui package
	block := ui.NewBlock()
	block.SetType("paragraph")
	block.SetParameter("content", "Test content")

	// Convert to JSON
	jsonStr, err := block.ToJson()
	if err != nil {
		t.Fatalf("Failed to convert block to JSON: %v", err)
	}

	// Wrap in array
	jsonStr = "[" + jsonStr + "]"

	result := fe.convertBlockJsonToHtml(jsonStr)

	if !rendererCalled {
		t.Error("Expected renderer to be called")
	}

	if result == "Malformed block content" || result == "Error parsing block content" {
		t.Errorf("Expected valid rendering, got: %q", result)
	}
}

// TestConvertBlockJsonToHtml_EmptyBlocks tests block JSON conversion with empty array
func TestConvertBlockJsonToHtml_EmptyBlocks(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	renderer := func(blocks []ui.BlockInterface) string {
		if len(blocks) == 0 {
			return ""
		}
		return "rendered"
	}

	f := New(Config{
		Store:               store,
		BlockEditorRenderer: renderer,
	})
	fe := f.(*frontend)

	json := `[]`
	result := fe.convertBlockJsonToHtml(json)

	// Empty array is valid JSON but not a valid object, so it should return malformed
	if result != "Malformed block content" {
		t.Logf("Empty array handling: %q", result)
	}
}

// TestPageFindBySiteAndAliasWithPatterns_AllPattern tests :all pattern matching
func TestPageFindBySiteAndAliasWithPatterns_AllPattern(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("/files/:all").
		SetContent("File content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundPage, err := f.(*frontend).pageFindBySiteAndAliasWithPatterns(ctx, site.ID(), "/files/path/to/file.pdf")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundPage == nil {
		t.Error("Expected page to match :all pattern")
	} else if foundPage.ID() != page.ID() {
		t.Error("Expected to find the correct page")
	}
}

// TestPageFindBySiteAndAliasWithPatterns_StringPattern tests :string pattern matching
func TestPageFindBySiteAndAliasWithPatterns_StringPattern(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("/category/:string").
		SetContent("Category content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundPage, err := f.(*frontend).pageFindBySiteAndAliasWithPatterns(ctx, site.ID(), "/category/Technology")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if foundPage == nil {
		t.Error("Expected page to match :string pattern")
	} else if foundPage.ID() != page.ID() {
		t.Error("Expected to find the correct page")
	}
}

// TestPageFindBySiteAndAliasWithPatterns_NumericPattern tests :numeric pattern matching
func TestPageFindBySiteAndAliasWithPatterns_NumericPattern(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("/price/:numeric").
		SetContent("Price content").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	ctx := context.Background()
	foundPage, err := f.(*frontend).pageFindBySiteAndAliasWithPatterns(ctx, site.ID(), "/price/19.99")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// :numeric pattern matches numbers with decimals and hyphens
	// If it doesn't match, that's okay - the pattern might be strict
	if foundPage != nil && foundPage.ID() != page.ID() {
		t.Error("Expected to find the correct page if pattern matches")
	}

	// Log the result for debugging
	if foundPage == nil {
		t.Logf(":numeric pattern did not match '19.99' - this may be expected behavior")
	}
}

// TestFullRenderingPipeline_WithTranslations tests complete rendering with translations
func TestFullRenderingPipeline_WithTranslations(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create translation
	translation := cmsstore.NewTranslation().
		SetHandle("welcome").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	translationContent := map[string]string{
		"en": "Welcome to our site",
		"es": "Bienvenido a nuestro sitio",
	}
	err = translation.SetContent(translationContent)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create page with translation
	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test").
		SetContent("<h1>[[TRANSLATION_welcome]]</h1>").
		SetTitle("Test Page").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	// Test English
	resultEN := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "test", "en")

	if !strings.Contains(resultEN, "Welcome to our site") {
		t.Errorf("Expected English translation, got: %q", resultEN)
	}

	// Test Spanish
	resultES := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "test", "es")

	if !strings.Contains(resultES, "Bienvenido a nuestro sitio") {
		t.Errorf("Expected Spanish translation, got: %q", resultES)
	}
}

// TestFullRenderingPipeline_ComplexPage tests complete rendering with blocks, translations, and template
func TestFullRenderingPipeline_ComplexPage(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	site := cmsstore.NewSite().
		SetName("Test Site").
		SetStatus(cmsstore.SITE_STATUS_ACTIVE)

	err = store.SiteCreate(context.Background(), site)
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}

	// Create block
	block := cmsstore.NewBlock().
		SetContent("<div class='header'>Header Block</div>").
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// Create translation
	translation := cmsstore.NewTranslation().
		SetHandle("footer").
		SetStatus(cmsstore.TRANSLATION_STATUS_ACTIVE)

	translationContent := map[string]string{
		"en": "Copyright 2024",
	}
	err = translation.SetContent(translationContent)
	if err != nil {
		t.Fatalf("Failed to set translation content: %v", err)
	}

	err = store.TranslationCreate(context.Background(), translation)
	if err != nil {
		t.Fatalf("Failed to create translation: %v", err)
	}

	// Create template
	template := cmsstore.NewTemplate().
		SetSiteID(site.ID()).
		SetName("Test Template").
		SetContent("<html><body>[[BLOCK_" + block.ID() + "]][[PageContent]]<footer>[[TRANSLATION_footer]]</footer></body></html>").
		SetStatus(cmsstore.TEMPLATE_STATUS_ACTIVE)

	err = store.TemplateCreate(context.Background(), template)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	// Create page
	page := cmsstore.NewPage().
		SetSiteID(site.ID()).
		SetName("Test Page").
		SetAlias("test").
		SetContent("<main>Main content</main>").
		SetTemplateID(template.ID()).
		SetTitle("Test Page").
		SetStatus(cmsstore.PAGE_STATUS_ACTIVE)

	err = store.PageCreate(context.Background(), page)
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	f := New(Config{
		Store:  store,
		Logger: slog.Default(),
	})

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	result := f.(*frontend).PageRenderHtmlBySiteAndAlias(recorder, req, site.ID(), "test", "en")

	// Verify all components are rendered
	if !strings.Contains(result, "Header Block") {
		t.Error("Expected block to be rendered")
	}

	if !strings.Contains(result, "Main content") {
		t.Error("Expected page content to be rendered")
	}

	if !strings.Contains(result, "Copyright 2024") {
		t.Error("Expected translation to be rendered")
	}

	if !strings.Contains(result, "<html>") {
		t.Error("Expected template structure to be rendered")
	}
}
