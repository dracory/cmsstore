package breadcrumbs

import (
	"context"
	"net/http"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/form"
	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

// TestBreadcrumbsBlockType_BasicProperties tests the breadcrumbs block type basic functionality
func TestBreadcrumbsBlockType_BasicProperties(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	breadcrumbsBlock := NewBreadcrumbsBlockType(store)

	// Test basic properties
	if breadcrumbsBlock.TypeKey() != cmsstore.BLOCK_TYPE_BREADCRUMBS {
		t.Errorf("Expected type %s, got %s", cmsstore.BLOCK_TYPE_BREADCRUMBS, breadcrumbsBlock.TypeKey())
	}

	if breadcrumbsBlock.TypeLabel() != "Breadcrumbs" {
		t.Errorf("Expected name 'Breadcrumbs', got '%s'", breadcrumbsBlock.TypeLabel())
	}
}

// TestBreadcrumbsBlockType_RenderBootstrap5 tests Bootstrap 5 rendering functionality
func TestBreadcrumbsBlockType_RenderBootstrap5(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	breadcrumbsBlock := NewBreadcrumbsBlockType(store)

	// Create a mock block with Bootstrap 5 configuration
	block := &TestBreadcrumbsBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_BREADCRUMBS_STYLE:          "default",
			cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE: "bootstrap5",
			cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS:      "custom-breadcrumbs",
			cmsstore.BLOCK_META_BREADCRUMBS_CSS_ID:         "main-breadcrumbs",
			cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR:      "→",
			cmsstore.BLOCK_META_BREADCRUMBS_HOME_TEXT:      "Home",
			cmsstore.BLOCK_META_BREADCRUMBS_HOME_URL:       "https://example.com",
		},
	}

	// Test rendering with Bootstrap 5 (no page needed for basic test)
	result, err := breadcrumbsBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Check that result contains expected Bootstrap 5 classes
	if len(result) == 0 {
		t.Error("Render returned empty result")
		return
	}

	// Basic checks for Bootstrap 5 structure
	expectedClasses := []string{
		"breadcrumb",
		"custom-breadcrumbs",
		"breadcrumb-item",
	}

	for _, class := range expectedClasses {
		if !contains(result, class) {
			t.Errorf("Expected result to contain class '%s', got: %s", class, result)
		}
	}

	// Check for separator
	if !contains(result, "→") {
		t.Errorf("Expected result to contain separator '→', got: %s", result)
	}

	// Check for CSS ID
	if !contains(result, `id="main-breadcrumbs"`) {
		t.Errorf("Expected result to contain CSS ID 'main-breadcrumbs', got: %s", result)
	}
}

// TestBreadcrumbsBlockType_RenderPlain tests plain rendering mode
func TestBreadcrumbsBlockType_RenderPlain(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	breadcrumbsBlock := NewBreadcrumbsBlockType(store)

	// Create a mock block with plain rendering
	block := &TestBreadcrumbsBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_BREADCRUMBS_STYLE:          "centered",
			cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE: "plain",
			cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS:      "plain-breadcrumbs",
			cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR:      ">",
		},
	}

	// Test rendering with plain mode
	result, err := breadcrumbsBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Check that result contains expected plain breadcrumb classes
	expectedClasses := []string{
		"breadcrumbs",
		"breadcrumbs-style-centered",
		"plain-breadcrumbs",
	}

	for _, class := range expectedClasses {
		if !contains(result, class) {
			t.Errorf("Expected result to contain class '%s', got: %s", class, result)
		}
	}
}

// TestBreadcrumbsBlockType_Validate tests validation functionality
func TestBreadcrumbsBlockType_Validate(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	breadcrumbsBlock := NewBreadcrumbsBlockType(store)

	// Test validation (breadcrumbs don't require any specific configuration)
	block := &TestBreadcrumbsBlock{
		meta: map[string]string{},
	}

	err = breadcrumbsBlock.Validate(block)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

// TestBreadcrumbsBlockType_GetPreview tests preview functionality
func TestBreadcrumbsBlockType_GetPreview(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	breadcrumbsBlock := NewBreadcrumbsBlockType(store)

	// Test preview with default configuration
	block := &TestBreadcrumbsBlock{
		meta: map[string]string{},
	}

	preview := breadcrumbsBlock.GetPreview(block)
	if preview != "Breadcrumbs: default (bootstrap5)" {
		t.Errorf("Expected preview 'Breadcrumbs: default (bootstrap5)', got '%s'", preview)
	}

	// Test preview with custom configuration
	block.meta[cmsstore.BLOCK_META_BREADCRUMBS_STYLE] = "centered"
	block.meta[cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE] = "plain"

	preview = breadcrumbsBlock.GetPreview(block)
	if preview != "Breadcrumbs: centered (plain)" {
		t.Errorf("Expected preview 'Breadcrumbs: centered (plain)', got '%s'", preview)
	}
}

// TestBreadcrumbsBlockType_AdminFields tests admin fields functionality
func TestBreadcrumbsBlockType_AdminFields(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	breadcrumbsBlock := NewBreadcrumbsBlockType(store)

	// Create a mock block with some meta data
	block := &TestBreadcrumbsBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_BREADCRUMBS_STYLE:          "default",
			cmsstore.BLOCK_META_BREADCRUMBS_RENDERING_MODE: "bootstrap5",
			cmsstore.BLOCK_META_BREADCRUMBS_SEPARATOR:      "→",
			cmsstore.BLOCK_META_BREADCRUMBS_HOME_TEXT:      "Home",
			cmsstore.BLOCK_META_BREADCRUMBS_CSS_CLASS:      "custom-class",
		},
	}

	// Test getting admin fields with valid request
	req, _ := http.NewRequest("GET", "/test", nil)
	fields := breadcrumbsBlock.GetAdminFields(block, req)
	if fields == nil {
		t.Error("Expected admin fields, got nil")
		return
	}

	// Check if it's a slice of form.FieldInterface
	fieldsSlice, ok := fields.([]form.FieldInterface)
	if !ok {
		t.Errorf("Expected fields to be []form.FieldInterface, got %T", fields)
		return
	}

	// We should have form fields now (style, rendering mode, separator, etc.)
	if len(fieldsSlice) == 0 {
		t.Error("Expected at least one form field, got none")
	}
}

// TestBreadcrumbsBlockType_SaveAdminFields tests saving admin fields
func TestBreadcrumbsBlockType_SaveAdminFields(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	breadcrumbsBlock := NewBreadcrumbsBlockType(store)

	// Create a mock block
	block := &TestBreadcrumbsBlock{
		meta: map[string]string{},
	}

	// Create a real http.Request with form data
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Form = map[string][]string{
		"breadcrumbs_style":          {"default"},
		"breadcrumbs_rendering_mode": {"bootstrap5"},
		"breadcrumbs_separator":      {"→"},
		"breadcrumbs_home_text":      {"Home"},
		"breadcrumbs_home_url":       {"https://example.com"},
		"breadcrumbs_css_class":      {"custom-class"},
		"breadcrumbs_css_id":         {"main-breadcrumbs"},
	}

	// Test saving admin fields
	err = breadcrumbsBlock.SaveAdminFields(req, block)
	if err != nil {
		t.Errorf("Expected no error saving admin fields, got: %v", err)
		return
	}

	// Verify that meta data was saved
	if block.meta["breadcrumbs_style"] != "default" {
		t.Errorf("Expected breadcrumbs_style to be 'default', got '%s'", block.meta["breadcrumbs_style"])
	}

	if block.meta["breadcrumbs_separator"] != "→" {
		t.Errorf("Expected breadcrumbs_separator to be '→', got '%s'", block.meta["breadcrumbs_separator"])
	}
}

// TestBreadcrumbsBlock is a mock implementation of BlockInterface for testing
type TestBreadcrumbsBlock struct {
	meta map[string]string
}

func (b *TestBreadcrumbsBlock) Data() map[string]string {
	if b.meta == nil {
		b.meta = make(map[string]string)
	}
	return b.meta
}
func (b *TestBreadcrumbsBlock) DataChanged() map[string]string { return make(map[string]string) }
func (b *TestBreadcrumbsBlock) MarkAsNotDirty()                {}

func (b *TestBreadcrumbsBlock) ID() string                                        { return "test-block" }
func (b *TestBreadcrumbsBlock) SetID(id string) cmsstore.BlockInterface           { return b }
func (b *TestBreadcrumbsBlock) Type() string                                      { return cmsstore.BLOCK_TYPE_BREADCRUMBS }
func (b *TestBreadcrumbsBlock) SetType(blockType string) cmsstore.BlockInterface  { return b }
func (b *TestBreadcrumbsBlock) Content() string                                   { return "test content" }
func (b *TestBreadcrumbsBlock) SetContent(content string) cmsstore.BlockInterface { return b }
func (b *TestBreadcrumbsBlock) CreatedAt() string                                 { return "2023-01-01 00:00:00" }
func (b *TestBreadcrumbsBlock) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse("2023-01-01 00:00:00")
}
func (b *TestBreadcrumbsBlock) SetCreatedAt(createdAt string) cmsstore.BlockInterface   { return b }
func (b *TestBreadcrumbsBlock) Editor() string                                          { return "blockeditor" }
func (b *TestBreadcrumbsBlock) SetEditor(editor string) cmsstore.BlockInterface         { return b }
func (b *TestBreadcrumbsBlock) Handle() string                                          { return "test-handle" }
func (b *TestBreadcrumbsBlock) SetHandle(handle string) cmsstore.BlockInterface         { return b }
func (b *TestBreadcrumbsBlock) Memo() string                                            { return "" }
func (b *TestBreadcrumbsBlock) SetMemo(memo string) cmsstore.BlockInterface             { return b }
func (b *TestBreadcrumbsBlock) Name() string                                            { return "Test Block" }
func (b *TestBreadcrumbsBlock) SetName(name string) cmsstore.BlockInterface             { return b }
func (b *TestBreadcrumbsBlock) PageID() string                                          { return "" }
func (b *TestBreadcrumbsBlock) SetPageID(pageID string) cmsstore.BlockInterface         { return b }
func (b *TestBreadcrumbsBlock) ParentID() string                                        { return "" }
func (b *TestBreadcrumbsBlock) SetParentID(parentID string) cmsstore.BlockInterface     { return b }
func (b *TestBreadcrumbsBlock) Sequence() string                                        { return "1" }
func (b *TestBreadcrumbsBlock) SequenceInt() int                                        { return 1 }
func (b *TestBreadcrumbsBlock) SetSequence(sequence string) cmsstore.BlockInterface     { return b }
func (b *TestBreadcrumbsBlock) SetSequenceInt(sequence int) cmsstore.BlockInterface     { return b }
func (b *TestBreadcrumbsBlock) SiteID() string                                          { return "" }
func (b *TestBreadcrumbsBlock) SetSiteID(siteID string) cmsstore.BlockInterface         { return b }
func (b *TestBreadcrumbsBlock) TemplateID() string                                      { return "" }
func (b *TestBreadcrumbsBlock) SetTemplateID(templateID string) cmsstore.BlockInterface { return b }
func (b *TestBreadcrumbsBlock) SoftDeletedAt() string                                   { return "" }
func (b *TestBreadcrumbsBlock) SoftDeletedAtCarbon() *carbon.Carbon                     { return nil }
func (b *TestBreadcrumbsBlock) SetSoftDeletedAt(softDeletedAt string) cmsstore.BlockInterface {
	return b
}
func (b *TestBreadcrumbsBlock) Status() string                                  { return cmsstore.BLOCK_STATUS_ACTIVE }
func (b *TestBreadcrumbsBlock) SetStatus(status string) cmsstore.BlockInterface { return b }
func (b *TestBreadcrumbsBlock) UpdatedAt() string                               { return "2023-01-01 00:00:00" }
func (b *TestBreadcrumbsBlock) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse("2023-01-01 00:00:00")
}
func (b *TestBreadcrumbsBlock) SetUpdatedAt(updatedAt string) cmsstore.BlockInterface { return b }
func (b *TestBreadcrumbsBlock) IsActive() bool                                        { return true }
func (b *TestBreadcrumbsBlock) IsInactive() bool                                      { return false }
func (b *TestBreadcrumbsBlock) IsSoftDeleted() bool                                   { return false }
func (b *TestBreadcrumbsBlock) MarshalToVersioning() (string, error)                  { return "", nil }
func (b *TestBreadcrumbsBlock) Meta(key string) string                                { return b.meta[key] }
func (b *TestBreadcrumbsBlock) SetMeta(key, value string) error                       { b.meta[key] = value; return nil }
func (b *TestBreadcrumbsBlock) Metas() (map[string]string, error)                     { return b.meta, nil }
func (b *TestBreadcrumbsBlock) SetMetas(metas map[string]string) error                { b.meta = metas; return nil }
func (b *TestBreadcrumbsBlock) UpsertMetas(metas map[string]string) error {
	if b.meta == nil {
		b.meta = make(map[string]string)
	}
	for k, v := range metas {
		b.meta[k] = v
	}
	return nil
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0
}
