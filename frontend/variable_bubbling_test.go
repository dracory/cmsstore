package frontend

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

// setupVariableBubblingTest creates shared test infrastructure
func setupVariableBubblingTest(t *testing.T) (*frontend, *cmsstore.BlockInterface, func()) {
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

	// Create a custom block type that sets variables via VarsFromContext
	testBlockType := &testVariableSettingBlockType{}
	cmsstore.RegisterCustomBlockType(testBlockType)

	// Create a block using the custom type
	block := cmsstore.NewBlock()
	block.SetName("Test Variable Block")
	block.SetType("test_variable_setter")
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Create frontend
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	fe := New(Config{
		Store:  store,
		Logger: logger,
	})

	f := fe.(*frontend)

	cleanup := func() {
		// Cleanup if needed
	}

	return f, &block, cleanup
}

// TestCustomVariableReplacement verifies that variables set by blocks
// are properly replaced in the content after block rendering.
func TestCustomVariableReplacement(t *testing.T) {
	f, block, cleanup := setupVariableBubblingTest(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/test", nil)
	content := "[[BLOCK_" + (*block).ID() + "]] Title: [[blog_title]]"

	html, err := f.renderContentToHtml(req, content, TemplateRenderHtmlByIDOptions{
		PageTitle: "Static Page Title", // This should NOT affect [[blog_title]]
	})
	if err != nil {
		t.Fatal(err)
	}

	// The block sets blog_title to "My Blog Post Title"
	// After rendering, [[blog_title]] should be replaced
	expected := "Title: My Blog Post Title"
	t.Logf("Actual HTML output: %q", html)
	if !contains(html, expected) {
		t.Errorf("Expected content to contain '%s', got: %s", expected, html)
	}
}

// TestVariableBubblingToPageTitle verifies that variables set by blocks IN THE TEMPLATE
// bubble up and are available when standard placeholders like [[PageTitle]] are processed.
// NOTE: This test only covers blocks in template content. See TestVariableBubblingFromPageContentToTemplate
// for the critical test of blocks in page content bubbling up to template placeholders.
func TestVariableBubblingToPageTitle(t *testing.T) {
	f, block, cleanup := setupVariableBubblingTest(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/test", nil)

	// Template content that uses both PageTitle and custom variables
	// In the fixed version, we should be able to use blog_title in PageTitle
	content := "[[BLOCK_" + (*block).ID() + "]]Title: [[blog_title]] | PageTitle was: [[PageTitle]]"

	html, err := f.renderContentToHtml(req, content, TemplateRenderHtmlByIDOptions{
		PageTitle: "[[blog_title]]", // Attempt to use block variable in PageTitle
	})
	if err != nil {
		t.Fatal(err)
	}

	// Currently this will fail because:
	// 1. [[PageTitle]] gets replaced with literal string "[[blog_title]]"
	// 2. Block renders and sets blog_title = "My Blog Post Title"
	// 3. [[blog_title]] replacement happens, but PageTitle was already replaced
	//
	// Expected after fix: PageTitle should be processed AFTER blocks render,
	// so [[blog_title]] inside PageTitle should be replaced with "My Blog Post Title"
	//
	// Currently we expect: "PageTitle was: [[blog_title]]" (literal string replacement)
	// After fix we expect: "PageTitle was: My Blog Post Title"

	// This test documents the CURRENT behavior (which is the bug)
	// After the fix, this assertion should be updated
	if contains(html, "PageTitle was: [[blog_title]]") {
		t.Logf("BUG CONFIRMED: PageTitle was replaced with literal '[[blog_title]]' instead of the block variable value")
		t.Logf("Current output: %s", html)
		// This is the failing test - it demonstrates the bug exists
		t.Errorf("Variable bubbling failed: PageTitle should contain the resolved blog_title value, got: %s", html)
	} else if contains(html, "PageTitle was: My Blog Post Title") {
		t.Logf("FIX CONFIRMED: Variable bubbling works correctly")
	}
}

// TestMultipleVariablesFromBlock verifies that a single block can set
// multiple variables that are all properly replaced.
func TestMultipleVariablesFromBlock(t *testing.T) {
	f, block, cleanup := setupVariableBubblingTest(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/test", nil)
	content := "[[BLOCK_" + (*block).ID() + "]][[blog_title]] - [[blog_summary]] - [[blog_date]]"

	html, err := f.renderContentToHtml(req, content, TemplateRenderHtmlByIDOptions{})
	if err != nil {
		t.Fatal(err)
	}

	expectedVars := []string{"My Blog Post Title", "This is a summary", "2024-01-15"}
	for _, expected := range expectedVars {
		if !contains(html, expected) {
			t.Errorf("Expected content to contain '%s', got: %s", expected, html)
		}
	}
}

// TestVariableBubblingFromPageContentToTemplate is the CRITICAL test that verifies
// variables set by blocks in PAGE CONTENT bubble up to TEMPLATE placeholders.
// This is the real-world scenario where:
// - Template has: <title>[[PageTitle]]</title>
// - Page title contains: "Blog Post [[blog_title]]"
// - Page content has: [[BLOCK_xxx]] (the Blog Post block)
// - The block sets: blog_title = "My Blog Post Title"
// Expected: The <title> should show "Blog Post My Blog Post Title"
func TestVariableBubblingFromPageContentToTemplate(t *testing.T) {
	f, block, cleanup := setupVariableBubblingTest(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/test", nil)

	// This simulates the TEMPLATE content with [[PageTitle]] placeholder
	templateContent := "<html><head><title>[[PageTitle]]</title></head><body>[[PageContent]]</body></html>"

	// This simulates the PAGE CONTENT with the Blog Post block
	pageContent := "[[BLOCK_" + (*block).ID() + "]]<p>Blog content here</p>"

	// This simulates the page title from the database containing a variable reference
	pageTitle := "Blog Post [[blog_title]]"

	html, err := f.renderContentToHtml(req, templateContent, TemplateRenderHtmlByIDOptions{
		PageContent: pageContent, // Block is in page content, not template
		PageTitle:   pageTitle,   // Title contains variable reference
	})
	if err != nil {
		t.Fatal(err)
	}

	// The critical assertion: [[blog_title]] in PageTitle should be resolved
	// to "My Blog Post Title" (set by the block in page content)
	expectedTitle := "<title>Blog Post My Blog Post Title</title>"
	if !contains(html, expectedTitle) {
		t.Errorf("Variable bubbling from page content to template failed.\nExpected title: %s\nActual HTML: %s", expectedTitle, html)
	}

	// Also verify the block rendered in the page content
	if !contains(html, "<!-- Block rendered -->") {
		t.Errorf("Block in page content did not render")
	}

	// Verify the variable is NOT left unreplaced
	if contains(html, "[[blog_title]]") {
		t.Errorf("Variable [[blog_title]] was not replaced in the output: %s", html)
	}
}

// TestVariableBubblingFromPageContentToMultiplePlaceholders verifies that
// variables from page content blocks are available in ALL template placeholders,
// not just PageTitle.
func TestVariableBubblingFromPageContentToMultiplePlaceholders(t *testing.T) {
	f, block, cleanup := setupVariableBubblingTest(t)
	defer cleanup()

	req, _ := http.NewRequest("GET", "/test", nil)

	// Template uses variables in multiple places
	templateContent := `<html>
		<head>
			<title>[[PageTitle]]</title>
			<meta name="description" content="[[PageMetaDescription]]">
		</head>
		<body>[[PageContent]]</body>
	</html>`

	// Page content has the block that sets variables
	pageContent := "[[BLOCK_" + (*block).ID() + "]]"

	html, err := f.renderContentToHtml(req, templateContent, TemplateRenderHtmlByIDOptions{
		PageContent:         pageContent,
		PageTitle:           "[[blog_title]]",            // Variable in title
		PageMetaDescription: "Summary: [[blog_summary]]", // Variable in meta
	})
	if err != nil {
		t.Fatal(err)
	}

	// Both should be resolved
	if !contains(html, "<title>My Blog Post Title</title>") {
		t.Errorf("PageTitle variable not resolved: %s", html)
	}
	if !contains(html, `content="Summary: This is a summary"`) {
		t.Errorf("PageMetaDescription variable not resolved: %s", html)
	}
}

// testVariableSettingBlockType is a test block type that sets variables via VarsFromContext
type testVariableSettingBlockType struct{}

func (t *testVariableSettingBlockType) TypeKey() string {
	return "test_variable_setter"
}

func (t *testVariableSettingBlockType) TypeLabel() string {
	return "Test Variable Setter Block"
}

func (t *testVariableSettingBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, options ...cmsstore.RenderOption) (string, error) {
	// Set custom variables that should be available in the page/template
	if vars := cmsstore.VarsFromContext(ctx); vars != nil {
		vars.Set("blog_title", "My Blog Post Title")
		vars.Set("blog_summary", "This is a summary")
		vars.Set("blog_date", "2024-01-15")
		vars.Set("blog_slug", "my-blog-post")
	}
	return "<!-- Block rendered -->", nil
}

func (t *testVariableSettingBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
	return nil
}

func (t *testVariableSettingBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
	return nil
}

func (t *testVariableSettingBlockType) GetCustomVariables() []cmsstore.BlockCustomVariable {
	return []cmsstore.BlockCustomVariable{
		{Name: "blog_title", Description: "The blog post title"},
		{Name: "blog_summary", Description: "The blog post summary"},
		{Name: "blog_date", Description: "The blog post date"},
		{Name: "blog_slug", Description: "The blog post slug"},
	}
}

func (t *testVariableSettingBlockType) Validate(block cmsstore.BlockInterface) error {
	return nil
}

func (t *testVariableSettingBlockType) GetPreview(block cmsstore.BlockInterface) string {
	return "Test Variable Setter"
}
