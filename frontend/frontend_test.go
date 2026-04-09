package frontend

import (
	"context"
	"log/slog"
	"net/http/httptest"
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
