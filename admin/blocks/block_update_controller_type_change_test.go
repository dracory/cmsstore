package admin

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/admin/shared"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/test"
	_ "modernc.org/sqlite"
)

func initBlockHandler(store cmsstore.StoreInterface) func(w http.ResponseWriter, r *http.Request) string {
	ui := UI(shared.UiConfig{
		Layout: shared.Layout,
		Logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		Store:  store,
	})

	return NewBlockUpdateController(ui).Handler
}

func TestBlockTypeChangeForDraftBlocks(t *testing.T) {
	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("InitStore should succeed, got error: %v", err)
	}

	handler := initBlockHandler(store)

	// Create a site first
	site, err := testutils.SeedSite(store, "Test Site")
	if err != nil {
		t.Fatalf("SeedSite should succeed, got error: %v", err)
	}

	// Create a draft block
	draftBlock := cmsstore.NewBlock()
	draftBlock.SetName("Test Draft Block")
	draftBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	draftBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	draftBlock.SetSiteID(site.ID())
	draftBlock.SetContent("<p>Original HTML content</p>")

	// Set some metadata
	err = draftBlock.SetMetas(map[string]string{
		"test_meta": "test_value",
	})
	if err != nil {
		t.Fatalf("Setting metas should succeed, got error: %v", err)
	}

	// Save the block
	err = store.BlockCreate(context.Background(), draftBlock)
	if err != nil {
		t.Fatalf("BlockCreate should succeed, got error: %v", err)
	}

	// Verify the block was saved with draft status
	savedBlock, err := store.BlockFindByID(context.Background(), draftBlock.ID())
	if err != nil {
		t.Fatalf("BlockFindByID should succeed, got error: %v", err)
	}
	if savedBlock.Status() != cmsstore.BLOCK_STATUS_DRAFT {
		t.Fatalf("Expected block status to be %s, got %s", cmsstore.BLOCK_STATUS_DRAFT, savedBlock.Status())
	}

	// Test 1: Verify type field is editable for draft blocks
	t.Run("DraftBlockTypeFieldIsEditable", func(t *testing.T) {
		// Create a request to get the block update form
		getValues := url.Values{
			"block_id": {draftBlock.ID()},
			"view":     {"settings"},
		}

		body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
			GetValues: getValues,
		})

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
		}

		// Verify that the type field is a select dropdown (editable)
		if !strings.Contains(body, "Block Settings") {
			t.Errorf("Expected settings view, but not found in response")
		}
		
		// Check for select element with block_type name
		hasSelect := strings.Contains(body, `<select`) && strings.Contains(body, `name="block_type"`)
		hasReadonlyText := strings.Contains(body, `name="block_type"`) && strings.Contains(body, `type="text"`) && strings.Contains(body, `readonly`)
		
		if !hasSelect {
			t.Errorf("Expected editable block type field (select), but not found in response")
		}
		if hasReadonlyText {
			t.Errorf("Expected block_type field to be editable (select), but found readonly text input")
		}
		if !strings.Contains(body, `Can only be changed while in draft status`) {
			t.Errorf("Expected help text about draft status, but not found in response")
		}
	})

	// Test 2: Successfully change block type for draft block
	t.Run("DraftBlockCanChangeType", func(t *testing.T) {
		// Create a POST request to change the block type
		postValues := url.Values{
			"block_id":     {draftBlock.ID()},
			"view":         {"settings"},
			"block_name":   {"Updated Block Name"},
			"block_type":   {cmsstore.BLOCK_TYPE_MENU},
			"block_status": {cmsstore.BLOCK_STATUS_DRAFT},
		}

		body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
			PostValues: postValues,
		})

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
		}

		if strings.Contains(body, "System error") {
			t.Errorf("Expected successful update, but found system error in response")
		}

		// Verify the block was updated
		updatedBlock, err := store.BlockFindByID(context.Background(), draftBlock.ID())
		if err != nil {
			t.Errorf("Expected to find updated block, got error: %v", err)
		}

		if updatedBlock == nil {
			t.Errorf("Expected to find updated block, but got nil")
		}

		if updatedBlock.Type() != cmsstore.BLOCK_TYPE_MENU {
			t.Errorf("Expected block type to be %s, got %s", cmsstore.BLOCK_TYPE_MENU, updatedBlock.Type())
		}

		if updatedBlock.Name() != "Updated Block Name" {
			t.Errorf("Expected block name to be 'Updated Block Name', got %s", updatedBlock.Name())
		}

		// Verify content and metadata were cleared
		if updatedBlock.Content() != "" {
			t.Errorf("Expected content to be cleared after type change, got: %s", updatedBlock.Content())
		}

		metas, err := updatedBlock.Metas()
		if err != nil {
			t.Errorf("Expected to get metas, got error: %v", err)
		}

		if len(metas) != 0 {
			t.Errorf("Expected metas to be cleared after type change, got: %v", metas)
		}
	})

	// Test 3: Cannot change type for published blocks
	t.Run("PublishedBlockCannotChangeType", func(t *testing.T) {
		// Create a published block
		publishedBlock := cmsstore.NewBlock()
		publishedBlock.SetName("Test Published Block")
		publishedBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
		publishedBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
		publishedBlock.SetContent("<p>Published content</p>")

		err = store.BlockCreate(context.Background(), publishedBlock)
		if err != nil {
			t.Fatalf("BlockCreate should succeed, got error: %v", err)
		}

		// Try to change the type of a published block
		postValues := url.Values{
			"block_id":     {publishedBlock.ID()},
			"view":         {"settings"},
			"block_name":   {"Updated Published Block"},
			"block_type":   {cmsstore.BLOCK_TYPE_MENU},
			"block_status": {cmsstore.BLOCK_STATUS_ACTIVE},
		}

		body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
			PostValues: postValues,
		})

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
		}

		if !strings.Contains(body, "Block type can only be changed while the block is in draft status") {
			t.Errorf("Expected error message about draft status, but not found in response")
		}

		// Verify the block type was NOT changed
		unchangedBlock, err := store.BlockFindByID(context.Background(), publishedBlock.ID())
		if err != nil {
			t.Errorf("Expected to find unchanged block, got error: %v", err)
		}

		if unchangedBlock.Type() != cmsstore.BLOCK_TYPE_HTML {
			t.Errorf("Expected block type to remain %s, got %s", cmsstore.BLOCK_TYPE_HTML, unchangedBlock.Type())
		}
	})

	// Test 4: Verify type field is readonly for published blocks
	t.Run("PublishedBlockTypeFieldIsReadonly", func(t *testing.T) {
		// First publish the block
		draftBlock.SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)
		err = store.BlockUpdate(context.Background(), draftBlock)
		if err != nil {
			t.Fatalf("BlockUpdate should succeed, got error: %v", err)
		}

		// Create a request to get the published block update form
		getValues := url.Values{
			"block_id": {draftBlock.ID()},
			"view":     {"settings"},
		}

		body, response, err := test.CallStringEndpoint(http.MethodGet, handler, test.NewRequestOptions{
			GetValues: getValues,
		})

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
		}

		// Verify that the type field is readonly
		if !strings.Contains(body, `readonly`) {
			t.Errorf("Expected readonly attribute for published block, but not found in response")
		}
		if !strings.Contains(body, `Block type cannot be changed after publication`) {
			t.Errorf("Expected help text about publication, but not found in response")
		}
		if strings.Contains(body, `<select name="block_type"`) {
			t.Errorf("Expected field to be readonly, but found select element")
		}
	})

	// Test 5: Invalid block type validation
	t.Run("InvalidBlockTypeRejected", func(t *testing.T) {
		// Reset block to draft
		draftBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
		err = store.BlockUpdate(context.Background(), draftBlock)
		if err != nil {
			t.Fatalf("BlockUpdate should succeed, got error: %v", err)
		}

		// Try to set an invalid block type
		postValues := url.Values{
			"block_id":     {draftBlock.ID()},
			"view":         {"settings"},
			"block_name":   {"Test Block"},
			"block_type":   {"invalid_type"},
			"block_status": {cmsstore.BLOCK_STATUS_DRAFT},
		}

		body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
			PostValues: postValues,
		})

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
		}

		if !strings.Contains(body, "Invalid block type: invalid_type") {
			t.Errorf("Expected error message about invalid type, but not found in response")
		}
	})
}

func TestBlockTypeChangeContentCleanup(t *testing.T) {
	// Initialize test store
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("InitStore should succeed, got error: %v", err)
	}

	handler := initBlockHandler(store)

	// Create a draft block with content and metadata
	draftBlock := cmsstore.NewBlock()
	draftBlock.SetName("Test Block")
	draftBlock.SetType(cmsstore.BLOCK_TYPE_HTML)
	draftBlock.SetStatus(cmsstore.BLOCK_STATUS_DRAFT)
	draftBlock.SetContent("<div>Original content</div>")

	// Add metadata
	err = draftBlock.SetMetas(map[string]string{
		"css_class": "original-class",
		"custom_id": "original-id",
	})
	if err != nil {
		t.Fatalf("Setting metas should succeed, got error: %v", err)
	}

	err = store.BlockCreate(context.Background(), draftBlock)
	if err != nil {
		t.Fatalf("BlockCreate should succeed, got error: %v", err)
	}

	// Change the block type
	postValues := url.Values{
		"block_id":     {draftBlock.ID()},
		"view":         {"settings"},
		"block_name":   {"Test Block"},
		"block_type":   {cmsstore.BLOCK_TYPE_MENU},
		"block_status": {cmsstore.BLOCK_STATUS_DRAFT},
	}

	body, response, err := test.CallStringEndpoint(http.MethodPost, handler, test.NewRequestOptions{
		PostValues: postValues,
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got: %d", http.StatusOK, response.StatusCode)
	}

	if strings.Contains(body, "System error") {
		t.Errorf("Expected successful update, but found system error in response")
	}

	// Verify content and metadata were cleared
	updatedBlock, err := store.BlockFindByID(context.Background(), draftBlock.ID())
	if err != nil {
		t.Errorf("Expected to find updated block, got error: %v", err)
	}

	if updatedBlock.Content() != "" {
		t.Errorf("Expected content to be cleared after type change, got: %s", updatedBlock.Content())
	}

	metas, err := updatedBlock.Metas()
	if err != nil {
		t.Errorf("Expected to get metas, got error: %v", err)
	}

	if len(metas) != 0 {
		t.Errorf("Expected metas to be cleared after type change, got: %v", metas)
	}

	if updatedBlock.Type() != cmsstore.BLOCK_TYPE_MENU {
		t.Errorf("Expected block type to be %s, got %s", cmsstore.BLOCK_TYPE_MENU, updatedBlock.Type())
	}
}
