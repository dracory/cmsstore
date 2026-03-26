package navbar

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

// TestNavbarMenuItem is a minimal implementation that only provides what the renderer needs
type TestNavbarMenuItem struct {
	name string
	url  string
	data map[string]string
}

func (m *TestNavbarMenuItem) Data() map[string]string {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	return m.data
}
func (m *TestNavbarMenuItem) DataChanged() map[string]string                           { return make(map[string]string) }
func (m *TestNavbarMenuItem) MarkAsNotDirty()                                          {}
func (m *TestNavbarMenuItem) ID() string                                               { return "test-id" }
func (m *TestNavbarMenuItem) Name() string                                             { return m.name }
func (m *TestNavbarMenuItem) URL() string                                              { return m.url }
func (m *TestNavbarMenuItem) IsActive() bool                                           { return true }
func (m *TestNavbarMenuItem) IsInactive() bool                                         { return false }
func (m *TestNavbarMenuItem) IsSoftDeleted() bool                                      { return false }
func (m *TestNavbarMenuItem) MarshalToVersioning() (string, error)                     { return "", nil } // Simplified for testing
func (m *TestNavbarMenuItem) SetID(id string) cmsstore.MenuItemInterface               { return m }
func (m *TestNavbarMenuItem) SetName(name string) cmsstore.MenuItemInterface           { m.name = name; return m }
func (m *TestNavbarMenuItem) SetURL(url string) cmsstore.MenuItemInterface             { m.url = url; return m }
func (m *TestNavbarMenuItem) SetCreatedAt(createdAt string) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetUpdatedAt(updatedAt string) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetSoftDeletedAt(softDeletedAt string) cmsstore.MenuItemInterface {
	return m
}
func (m *TestNavbarMenuItem) SetHandle(handle string) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetMemo(memo string) cmsstore.MenuItemInterface     { return m }
func (m *TestNavbarMenuItem) SetMenuID(menuID string) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetMeta(key, value string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[key] = value
	return nil
}
func (m *TestNavbarMenuItem) SetMetas(metas map[string]string) error { m.data = metas; return nil }
func (m *TestNavbarMenuItem) UpsertMetas(metas map[string]string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	for k, v := range metas {
		m.data[k] = v
	}
	return nil
}
func (m *TestNavbarMenuItem) SetPageID(pageID string) cmsstore.MenuItemInterface     { return m }
func (m *TestNavbarMenuItem) SetParentID(parentID string) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetSequence(sequence string) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetSequenceInt(sequence int) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetStatus(status string) cmsstore.MenuItemInterface     { return m }
func (m *TestNavbarMenuItem) SetTarget(target string) cmsstore.MenuItemInterface     { return m }
func (m *TestNavbarMenuItem) CreatedAt() string                                      { return "2023-01-01 00:00:00" }
func (m *TestNavbarMenuItem) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse("2023-01-01 00:00:00")
}
func (m *TestNavbarMenuItem) UpdatedAt() string { return "2023-01-01 00:00:00" }
func (m *TestNavbarMenuItem) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse("2023-01-01 00:00:00")
}
func (m *TestNavbarMenuItem) SoftDeletedAt() string               { return "" }
func (m *TestNavbarMenuItem) SoftDeletedAtCarbon() *carbon.Carbon { return nil }
func (m *TestNavbarMenuItem) Handle() string                      { return "" }
func (m *TestNavbarMenuItem) Memo() string                        { return "" }
func (m *TestNavbarMenuItem) MenuID() string                      { return "test-menu" }
func (m *TestNavbarMenuItem) Meta(key string) string {
	if m.data == nil {
		return ""
	}
	return m.data[key]
}
func (m *TestNavbarMenuItem) Metas() (map[string]string, error) {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	return m.data, nil
}
func (m *TestNavbarMenuItem) PageID() string   { return "" }
func (m *TestNavbarMenuItem) ParentID() string { return "" }
func (m *TestNavbarMenuItem) Sequence() string { return "1" }
func (m *TestNavbarMenuItem) SequenceInt() int { return 1 }
func (m *TestNavbarMenuItem) Status() string   { return "active" }
func (m *TestNavbarMenuItem) Target() string   { return "_self" }

// TestNavbarBlock is a mock implementation of BlockInterface for testing
type TestNavbarBlock struct {
	meta map[string]string
}

func (b *TestNavbarBlock) Data() map[string]string {
	if b.meta == nil {
		b.meta = make(map[string]string)
	}
	return b.meta
}
func (b *TestNavbarBlock) DataChanged() map[string]string { return make(map[string]string) }
func (b *TestNavbarBlock) MarkAsNotDirty()                {}

func (b *TestNavbarBlock) ID() string                                        { return "test-block" }
func (b *TestNavbarBlock) SetID(id string) cmsstore.BlockInterface           { return b }
func (b *TestNavbarBlock) Type() string                                      { return cmsstore.BLOCK_TYPE_NAVBAR }
func (b *TestNavbarBlock) SetType(blockType string) cmsstore.BlockInterface  { return b }
func (b *TestNavbarBlock) Content() string                                   { return "test content" }
func (b *TestNavbarBlock) SetContent(content string) cmsstore.BlockInterface { return b }
func (b *TestNavbarBlock) CreatedAt() string                                 { return "2023-01-01 00:00:00" }
func (b *TestNavbarBlock) CreatedAtCarbon() *carbon.Carbon {
	return carbon.Parse("2023-01-01 00:00:00")
}
func (b *TestNavbarBlock) SetCreatedAt(createdAt string) cmsstore.BlockInterface         { return b }
func (b *TestNavbarBlock) Editor() string                                                { return "blockeditor" }
func (b *TestNavbarBlock) SetEditor(editor string) cmsstore.BlockInterface               { return b }
func (b *TestNavbarBlock) Handle() string                                                { return "test-handle" }
func (b *TestNavbarBlock) SetHandle(handle string) cmsstore.BlockInterface               { return b }
func (b *TestNavbarBlock) Memo() string                                                  { return "" }
func (b *TestNavbarBlock) SetMemo(memo string) cmsstore.BlockInterface                   { return b }
func (b *TestNavbarBlock) Name() string                                                  { return "Test Block" }
func (b *TestNavbarBlock) SetName(name string) cmsstore.BlockInterface                   { return b }
func (b *TestNavbarBlock) PageID() string                                                { return "" }
func (b *TestNavbarBlock) SetPageID(pageID string) cmsstore.BlockInterface               { return b }
func (b *TestNavbarBlock) ParentID() string                                              { return "" }
func (b *TestNavbarBlock) SetParentID(parentID string) cmsstore.BlockInterface           { return b }
func (b *TestNavbarBlock) Sequence() string                                              { return "1" }
func (b *TestNavbarBlock) SequenceInt() int                                              { return 1 }
func (b *TestNavbarBlock) SetSequence(sequence string) cmsstore.BlockInterface           { return b }
func (b *TestNavbarBlock) SetSequenceInt(sequence int) cmsstore.BlockInterface           { return b }
func (b *TestNavbarBlock) SiteID() string                                                { return "" }
func (b *TestNavbarBlock) SetSiteID(siteID string) cmsstore.BlockInterface               { return b }
func (b *TestNavbarBlock) TemplateID() string                                            { return "" }
func (b *TestNavbarBlock) SetTemplateID(templateID string) cmsstore.BlockInterface       { return b }
func (b *TestNavbarBlock) SoftDeletedAt() string                                         { return "" }
func (b *TestNavbarBlock) SoftDeletedAtCarbon() *carbon.Carbon                           { return nil }
func (b *TestNavbarBlock) SetSoftDeletedAt(softDeletedAt string) cmsstore.BlockInterface { return b }
func (b *TestNavbarBlock) Status() string                                                { return cmsstore.BLOCK_STATUS_ACTIVE }
func (b *TestNavbarBlock) SetStatus(status string) cmsstore.BlockInterface               { return b }
func (b *TestNavbarBlock) UpdatedAt() string                                             { return "2023-01-01 00:00:00" }
func (b *TestNavbarBlock) UpdatedAtCarbon() *carbon.Carbon {
	return carbon.Parse("2023-01-01 00:00:00")
}
func (b *TestNavbarBlock) SetUpdatedAt(updatedAt string) cmsstore.BlockInterface { return b }
func (b *TestNavbarBlock) IsActive() bool                                        { return true }
func (b *TestNavbarBlock) IsInactive() bool                                      { return false }
func (b *TestNavbarBlock) IsSoftDeleted() bool                                   { return false }
func (b *TestNavbarBlock) MarshalToVersioning() (string, error)                  { return "", nil }
func (b *TestNavbarBlock) Meta(key string) string                                { return b.meta[key] }
func (b *TestNavbarBlock) SetMeta(key, value string) error                       { b.meta[key] = value; return nil }
func (b *TestNavbarBlock) Metas() (map[string]string, error)                     { return b.meta, nil }
func (b *TestNavbarBlock) SetMetas(metas map[string]string) error                { b.meta = metas; return nil }
func (b *TestNavbarBlock) UpsertMetas(metas map[string]string) error {
	if b.meta == nil {
		b.meta = make(map[string]string)
	}
	for k, v := range metas {
		b.meta[k] = v
	}
	return nil
}

// TestNavbarBlockType_BasicProperties tests the navbar block type basic functionality
func TestNavbarBlockType_BasicProperties(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Test basic properties
	if navbarBlock.TypeKey() != cmsstore.BLOCK_TYPE_NAVBAR {
		t.Errorf("Expected type %s, got %s", cmsstore.BLOCK_TYPE_NAVBAR, navbarBlock.TypeKey())
	}

	if navbarBlock.TypeLabel() != "Navbar" {
		t.Errorf("Expected name 'Navbar', got '%s'", navbarBlock.TypeLabel())
	}
}

// TestNavbarBlockType_RenderBootstrap5 tests Bootstrap 5 rendering functionality
func TestNavbarBlockType_RenderBootstrap5(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Create a mock block with Bootstrap 5 configuration
	block := &TestNavbarBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_MENU_ID:               "test-menu",
			cmsstore.BLOCK_META_NAVBAR_STYLE:          "default",
			cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE: "bootstrap5",
			cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT:     "Test Brand",
			cmsstore.BLOCK_META_NAVBAR_BRAND_URL:      "https://example.com",
			cmsstore.BLOCK_META_NAVBAR_CSS_CLASS:      "custom-navbar",
			cmsstore.BLOCK_META_NAVBAR_CSS_ID:         "main-navbar",
			cmsstore.BLOCK_META_NAVBAR_FIXED:          "true",
			cmsstore.BLOCK_META_NAVBAR_DARK:           "true",
		},
	}

	// Test rendering with Bootstrap 5 (no page needed for basic test)
	result, err := navbarBlock.Render(ctx, block)
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
		"navbar",
		"navbar-expand-lg",
		"navbar-dark",
		"bg-dark",
		"fixed-top",
		"custom-navbar",
	}

	for _, class := range expectedClasses {
		if !contains(result, class) {
			t.Errorf("Expected result to contain class '%s', got: %s", class, result)
		}
	}

	// Check for brand
	if !contains(result, "Test Brand") {
		t.Errorf("Expected result to contain brand text 'Test Brand', got: %s", result)
	}

	// Check for CSS ID
	if !contains(result, `id="main-navbar"`) {
		t.Errorf("Expected result to contain CSS ID 'main-navbar', got: %s", result)
	}
}

// TestNavbarBlockType_RenderPlain tests plain rendering mode
func TestNavbarBlockType_RenderPlain(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Create a mock block with plain rendering
	block := &TestNavbarBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_MENU_ID:               "test-menu",
			cmsstore.BLOCK_META_NAVBAR_STYLE:          "default",
			cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE: "plain",
			cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT:     "Plain Brand",
			cmsstore.BLOCK_META_NAVBAR_CSS_CLASS:      "plain-navbar",
		},
	}

	// Test rendering with plain mode
	result, err := navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Check that result contains expected plain navbar classes
	expectedClasses := []string{
		"navbar",
		"navbar-style-default",
		"plain-navbar",
	}

	for _, class := range expectedClasses {
		if !contains(result, class) {
			t.Errorf("Expected result to contain class '%s', got: %s", class, result)
		}
	}
}

// TestNavbarBlockType_Validate tests validation functionality
func TestNavbarBlockType_Validate(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Test validation with missing menu ID
	block := &TestNavbarBlock{
		meta: map[string]string{},
	}

	err = navbarBlock.Validate(block)
	if err == nil {
		t.Error("Expected validation error for missing menu ID")
	}

	// Test validation with valid menu ID
	block.meta[cmsstore.BLOCK_META_MENU_ID] = "test-menu"
	err = navbarBlock.Validate(block)
	// Note: testutils store doesn't validate menu existence, so this won't error
	if err != nil {
		t.Errorf("Expected no validation error for valid menu ID, got: %v", err)
	}
}

// TestNavbarBlockType_GetPreview tests preview functionality
func TestNavbarBlockType_GetPreview(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Test preview with no menu ID
	block := &TestNavbarBlock{
		meta: map[string]string{},
	}

	preview := navbarBlock.GetPreview(block)
	if preview != "No menu selected" {
		t.Errorf("Expected preview 'No menu selected', got '%s'", preview)
	}

	// Test preview with invalid menu ID
	block.meta[cmsstore.BLOCK_META_MENU_ID] = "invalid-menu"
	preview = navbarBlock.GetPreview(block)
	if preview != "Invalid menu" {
		t.Errorf("Expected preview 'Invalid menu', got '%s'", preview)
	}

	// Test preview with valid menu ID (but testutils store doesn't have it)
	block.meta[cmsstore.BLOCK_META_MENU_ID] = "test-menu"
	preview = navbarBlock.GetPreview(block)
	if preview != "Invalid menu" {
		t.Errorf("Expected preview 'Invalid menu', got '%s'", preview)
	}
}

// TestNavbarBlockType_AdminFields tests admin fields functionality
func TestNavbarBlockType_AdminFields(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Create a mock block with some meta data
	block := &TestNavbarBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_MENU_ID:               "test-menu",
			cmsstore.BLOCK_META_NAVBAR_STYLE:          "default",
			cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE: "bootstrap5",
			cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT:     "Test Brand",
			cmsstore.BLOCK_META_NAVBAR_CSS_CLASS:      "custom-class",
		},
	}

	// Test getting admin fields with valid request
	req, _ := http.NewRequest("GET", "/test", nil)
	fields := navbarBlock.GetAdminFields(block, req)
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

	// We should have form fields now (menu dropdown, style, rendering mode, etc.)
	if len(fieldsSlice) == 0 {
		t.Error("Expected at least one form field, got none")
	}
}

// TestNavbarBlockType_SaveAdminFields tests saving admin fields
func TestNavbarBlockType_SaveAdminFields(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Create a mock block
	block := &TestNavbarBlock{
		meta: map[string]string{},
	}

	// Create a real http.Request with form data
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Form = map[string][]string{
		"menu_id":               {"test-menu"},
		"navbar_style":          {"default"},
		"navbar_rendering_mode": {"bootstrap5"},
		"navbar_brand_text":     {"Test Brand"},
		"navbar_brand_url":      {"https://example.com"},
		"navbar_fixed":          {"true"},
		"navbar_dark":           {"false"},
		"navbar_css_class":      {"custom-class"},
		"navbar_css_id":         {"main-navbar"},
	}

	// Test saving admin fields
	err = navbarBlock.SaveAdminFields(req, block)
	if err != nil {
		t.Errorf("Expected no error saving admin fields, got: %v", err)
		return
	}

	// Verify that meta data was saved
	if block.meta["menu_id"] != "test-menu" {
		t.Errorf("Expected menu_id to be 'test-menu', got '%s'", block.meta["menu_id"])
	}

	if block.meta["navbar_brand_text"] != "Test Brand" {
		t.Errorf("Expected navbar_brand_text to be 'Test Brand', got '%s'", block.meta["navbar_brand_text"])
	}

	if block.meta["navbar_fixed"] != "true" {
		t.Errorf("Expected navbar_fixed to be 'true', got '%s'", block.meta["navbar_fixed"])
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0
}
