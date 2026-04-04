package frontend

import (
	"context"
	"log/slog"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

func TestRenderMenuBlock(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatal(err)
	}

	// Create pages for menu items
	page1 := cmsstore.NewPage()
	page1.SetTitle("Page 1")
	page1.SetAlias("/page-1")
	page1.SetSiteID(site.ID())
	page1.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(ctx, page1)
	if err != nil {
		t.Fatal(err)
	}

	page2 := cmsstore.NewPage()
	page2.SetTitle("Page 2")
	page2.SetAlias("/page-2")
	page2.SetSiteID(site.ID())
	page2.SetStatus(cmsstore.PAGE_STATUS_ACTIVE)
	err = store.PageCreate(ctx, page2)
	if err != nil {
		t.Fatal(err)
	}

	// Create a menu
	menu := cmsstore.NewMenu()
	menu.SetName("Test Menu")
	menu.SetSiteID(site.ID())
	menu.SetStatus(cmsstore.MENU_STATUS_ACTIVE)
	err = store.MenuCreate(ctx, menu)
	if err != nil {
		t.Fatal(err)
	}

	// Create menu items
	menuItem1 := cmsstore.NewMenuItem()
	menuItem1.SetName("Home")
	menuItem1.SetMenuID(menu.ID())
	menuItem1.SetPageID(page1.ID())
	menuItem1.SetSequenceInt(1)
	menuItem1.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	err = store.MenuItemCreate(ctx, menuItem1)
	if err != nil {
		t.Fatal(err)
	}

	menuItem2 := cmsstore.NewMenuItem()
	menuItem2.SetName("About")
	menuItem2.SetMenuID(menu.ID())
	menuItem2.SetPageID(page2.ID())
	menuItem2.SetSequenceInt(2)
	menuItem2.SetStatus(cmsstore.MENU_ITEM_STATUS_ACTIVE)
	err = store.MenuItemCreate(ctx, menuItem2)
	if err != nil {
		t.Fatal(err)
	}

	// Create a menu block
	block := cmsstore.NewBlock()
	block.SetName("Test Menu Block")
	block.SetType(cmsstore.BLOCK_TYPE_MENU)
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	block.SetMeta(cmsstore.BLOCK_META_MENU_ID, menu.ID())
	block.SetMeta(cmsstore.BLOCK_META_MENU_STYLE, cmsstore.BLOCK_MENU_STYLE_VERTICAL)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Create frontend
	logger := slog.New(slog.NewTextHandler(nil, nil))
	fe := New(Config{
		Store:  store,
		Logger: logger,
	})

	// Test rendering
	content, err := fe.(*frontend).fetchBlockContent(ctx, block.ID())
	if err != nil {
		t.Fatal(err)
	}

	if content == "" {
		t.Fatal("Expected menu block content, got empty string")
	}

	// Check that the content contains menu items
	if !contains(content, "Home") {
		t.Errorf("Expected content to contain 'Home', got: %s", content)
	}

	if !contains(content, "About") {
		t.Errorf("Expected content to contain 'About', got: %s", content)
	}

	// Check that it contains links
	if !contains(content, "/page-1") {
		t.Errorf("Expected content to contain '/page-1', got: %s", content)
	}

	if !contains(content, "/page-2") {
		t.Errorf("Expected content to contain '/page-2', got: %s", content)
	}
}

func TestRenderHTMLBlock(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatal(err)
	}

	// Create an HTML block
	block := cmsstore.NewBlock()
	block.SetName("Test HTML Block")
	block.SetType(cmsstore.BLOCK_TYPE_HTML)
	block.SetContent("<div>Test Content</div>")
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Create frontend
	logger := slog.New(slog.NewTextHandler(nil, nil))
	fe := New(Config{
		Store:  store,
		Logger: logger,
	})

	// Test rendering
	content, err := fe.(*frontend).fetchBlockContent(ctx, block.ID())
	if err != nil {
		t.Fatal(err)
	}

	if content != "<div>Test Content</div>" {
		t.Errorf("Expected '<div>Test Content</div>', got: %s", content)
	}
}

func TestBlockTypeDispatcher(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(nil, nil))
	fe := New(Config{
		Store:  store,
		Logger: logger,
	})

	f := fe.(*frontend)

	// Test HTML block
	htmlBlock := cmsstore.NewBlock()
	htmlBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	htmlBlock.SetContent("<p>HTML Content</p>")

	content, err := f.renderBlockByType(ctx, htmlBlock)
	if err != nil {
		t.Fatal(err)
	}

	if content != "<p>HTML Content</p>" {
		t.Errorf("Expected '<p>HTML Content</p>', got: %s", content)
	}

	// Test empty type defaults to HTML
	emptyBlock := cmsstore.NewBlock()
	emptyBlock.SetType("")
	emptyBlock.SetContent("<p>Default Content</p>")

	content, err = f.renderBlockByType(ctx, emptyBlock)
	if err != nil {
		t.Fatal(err)
	}

	if content != "<p>Default Content</p>" {
		t.Errorf("Expected '<p>Default Content</p>', got: %s", content)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestRequestContextPassedToBlocks verifies that the http request is passed
// in the context when rendering blocks
func TestRequestContextPassedToBlocks(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create a site
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatal(err)
	}

	// Create a custom block type that checks for request in context
	var requestReceived *http.Request
	testBlockType := &testBlockTypeWithRequestCheck{
		onRender: func(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
			requestReceived = cmsstore.RequestFromContext(ctx)
			return "test content", nil
		},
	}
	cmsstore.RegisterCustomBlockType(testBlockType)

	// Create a block using the custom type
	block := cmsstore.NewBlock()
	block.SetName("Test Block")
	block.SetType("test_request_check")
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Create frontend
	logger := slog.New(slog.NewTextHandler(nil, nil))
	fe := New(Config{
		Store:  store,
		Logger: logger,
	})

	// Create a request with query parameters
	req, err := http.NewRequest("GET", "/test?q=searchterm", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Render a page that contains the block
	f := fe.(*frontend)
	content := "[[BLOCK_" + block.ID() + "]]"
	html, err := f.renderContentToHtml(req, content, TemplateRenderHtmlByIDOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// Verify the block received the request
	if requestReceived == nil {
		t.Fatal("Expected block to receive http_request in context, but got nil")
	}

	// Verify the query parameter was preserved
	if requestReceived.URL.Query().Get("q") != "searchterm" {
		t.Errorf("Expected query parameter 'q' to be 'searchterm', got: %s", requestReceived.URL.Query().Get("q"))
	}

	// Verify the content was rendered
	if html != "test content" {
		t.Errorf("Expected rendered content to be 'test content', got: %s", html)
	}
}

// testBlockTypeWithRequestCheck is a test block type that captures the request from context
type testBlockTypeWithRequestCheck struct {
	onRender func(ctx context.Context, block cmsstore.BlockInterface) (string, error)
}

func (t *testBlockTypeWithRequestCheck) TypeKey() string {
	return "test_request_check"
}

func (t *testBlockTypeWithRequestCheck) TypeLabel() string {
	return "Test Request Check Block"
}

func (t *testBlockTypeWithRequestCheck) Render(ctx context.Context, block cmsstore.BlockInterface, options ...cmsstore.RenderOption) (string, error) {
	return t.onRender(ctx, block)
}

func (t *testBlockTypeWithRequestCheck) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	return nil
}

func (t *testBlockTypeWithRequestCheck) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	return nil
}

func (t *testBlockTypeWithRequestCheck) Validate(block cmsstore.BlockInterface) error {
	return nil
}

func (t *testBlockTypeWithRequestCheck) GetPreview(block cmsstore.BlockInterface) string {
	return "Test Request Check"
}
