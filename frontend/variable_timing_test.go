package frontend

import (
	"context"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	_ "modernc.org/sqlite"
)

// TestVariableTimingVerification verifies that variables captured before
// page content rendering are still available after (because VarsContext is a pointer)
func TestVariableTimingVerification(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Create site and block
	site := cmsstore.NewSite()
	site.SetName("Test Site")
	site.SetStatus(cmsstore.SITE_STATUS_ACTIVE)
	err = store.SiteCreate(ctx, site)
	if err != nil {
		t.Fatal(err)
	}

	testBlockType := &testVariableSettingBlockType{}
	cmsstore.RegisterCustomBlockType(testBlockType)

	block := cmsstore.NewBlock()
	block.SetName("Test Variable Block")
	block.SetType("test_variable_setter")
	block.SetSiteID(site.ID())
	block.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
	err = store.BlockCreate(ctx, block)
	if err != nil {
		t.Fatal(err)
	}

	// Simulate the exact flow in renderContentToHtml
	ctx = cmsstore.RequestToContext(ctx, &http.Request{})
	ctx = cmsstore.WithVarsContext(ctx)

	// Step 1: Get customVars BEFORE rendering page content (line 581)
	customVarsBefore := cmsstore.VarsFromContext(ctx)
	if customVarsBefore == nil {
		t.Fatal("customVars is nil")
	}

	beforeCount := len(customVarsBefore.All())
	t.Logf("Variables before page content rendering: %d", beforeCount)

	// Step 2: Render page content blocks (line 589)
	// Manually render the block to set variables
	renderedBlock, err := testBlockType.Render(ctx, block)
	if err != nil {
		t.Fatal(err)
	}
	_ = renderedBlock

	// Step 3: Check if customVarsBefore (captured earlier) now has the new variables
	afterCount := len(customVarsBefore.All())
	t.Logf("Variables after page content rendering: %d", afterCount)

	if afterCount <= beforeCount {
		t.Errorf("Variables were NOT added to the captured customVars pointer. Before: %d, After: %d", beforeCount, afterCount)
	}

	// Verify specific variable
	if val, exists := customVarsBefore.All()["blog_title"]; !exists {
		t.Errorf("blog_title not found in customVars captured before rendering")
	} else if val != "My Blog Post Title" {
		t.Errorf("blog_title has wrong value: %s", val)
	} else {
		t.Logf("SUCCESS: Variables added after capture are accessible via the pointer")
	}
}
