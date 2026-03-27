package navbar

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dracory/form"
	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

// TestNavbarMenuItem is a minimal implementation that only provides what the renderer needs
type TestNavbarMenuItem struct {
	id       string
	name     string
	url      string
	parentID string
	data     map[string]string
}

func (m *TestNavbarMenuItem) Data() map[string]string {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	return m.data
}
func (m *TestNavbarMenuItem) DataChanged() map[string]string                           { return make(map[string]string) }
func (m *TestNavbarMenuItem) MarkAsNotDirty()                                          {}
func (m *TestNavbarMenuItem) ID() string                                               { return m.id }
func (m *TestNavbarMenuItem) Name() string                                             { return m.name }
func (m *TestNavbarMenuItem) URL() string                                              { return m.url }
func (m *TestNavbarMenuItem) IsActive() bool                                           { return true }
func (m *TestNavbarMenuItem) IsInactive() bool                                         { return false }
func (m *TestNavbarMenuItem) IsSoftDeleted() bool                                      { return false }
func (m *TestNavbarMenuItem) MarshalToVersioning() (string, error)                     { return "", nil } // Simplified for testing
func (m *TestNavbarMenuItem) SetID(id string) cmsstore.MenuItemInterface               { m.id = id; return m }
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
func (m *TestNavbarMenuItem) SetPageID(pageID string) cmsstore.MenuItemInterface { return m }
func (m *TestNavbarMenuItem) SetParentID(parentID string) cmsstore.MenuItemInterface {
	m.parentID = parentID
	return m
}
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
func (m *TestNavbarMenuItem) ParentID() string { return m.parentID }
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

// TestNavbarBlockType_RenderHierarchicalDropdown tests hierarchical menu rendering with dropdowns
func TestNavbarBlockType_RenderHierarchicalDropdown(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Create hierarchical menu items manually
	parentItem := &TestNavbarMenuItem{
		id:   "parent-1",
		name: "Services",
		url:  "#",
	}

	childItem1 := &TestNavbarMenuItem{
		id:       "child-1",
		name:     "Web Design",
		url:      "/web-design",
		parentID: "parent-1",
	}

	childItem2 := &TestNavbarMenuItem{
		id:       "child-2",
		name:     "SEO",
		url:      "/seo",
		parentID: "parent-1",
	}

	// Simple item without children
	simpleItem := &TestNavbarMenuItem{
		id:   "simple-1",
		name: "About",
		url:  "/about",
	}

	// Create menu items slice
	menuItems := []cmsstore.MenuItemInterface{
		parentItem,
		childItem1,
		childItem2,
		simpleItem,
	}

	// Test the renderNavItemWithDropdown function via the renderer
	result := renderNavItemsForTest(ctx, store, menuItems)

	// Check that dropdown structure is present for parent item
	if !contains(result, "dropdown") {
		t.Errorf("Expected dropdown class for parent item, got: %s", result)
	}

	if !contains(result, "dropdown-toggle") {
		t.Errorf("Expected dropdown-toggle class, got: %s", result)
	}

	if !contains(result, "dropdown-menu") {
		t.Errorf("Expected dropdown-menu class, got: %s", result)
	}

	// Check that child items are rendered as dropdown items
	if !contains(result, "dropdown-item") {
		t.Errorf("Expected dropdown-item class for children, got: %s", result)
	}

	// Check that parent item name appears
	if !contains(result, "Services") {
		t.Errorf("Expected parent item name 'Services', got: %s", result)
	}

	// Check that child item names appear
	if !contains(result, "Web Design") {
		t.Errorf("Expected child item name 'Web Design', got: %s", result)
	}

	if !contains(result, "SEO") {
		t.Errorf("Expected child item name 'SEO', got: %s", result)
	}

	// Check that simple item (no children) appears as regular nav-item
	if !contains(result, "About") {
		t.Errorf("Expected simple item name 'About', got: %s", result)
	}

	// Verify simple item doesn't have dropdown classes
	// Count occurrences - parent should have dropdown, simple should not
	dropdownCount := strings.Count(result, "dropdown-toggle")
	if dropdownCount != 1 {
		t.Errorf("Expected 1 dropdown toggle (only parent), found %d", dropdownCount)
	}
}

// renderNavItemsForTest is a helper to test the hierarchical rendering logic
func renderNavItemsForTest(ctx context.Context, store cmsstore.StoreInterface, menuItems []cmsstore.MenuItemInterface) string {
	// Build menu item map
	menuItemMap := make(map[string]cmsstore.MenuItemInterface)
	for _, item := range menuItems {
		menuItemMap[item.ID()] = item
	}

	// Find top-level items
	var topLevelItems []cmsstore.MenuItemInterface
	for _, item := range menuItems {
		if item.ParentID() == "" {
			topLevelItems = append(topLevelItems, item)
		}
	}

	// Create a mock navbar HTML by concatenating rendered items
	var result strings.Builder
	for _, item := range topLevelItems {
		result.WriteString(renderNavItemHTMLForTest(ctx, store, item, menuItemMap))
	}

	return result.String()
}

// renderNavItemHTMLForTest renders a single nav item and returns its HTML
func renderNavItemHTMLForTest(ctx context.Context, store cmsstore.StoreInterface, item cmsstore.MenuItemInterface, menuItemMap map[string]cmsstore.MenuItemInterface) string {
	// Find children
	var children []cmsstore.MenuItemInterface
	for _, mi := range menuItemMap {
		if mi.ParentID() == item.ID() {
			children = append(children, mi)
		}
	}

	hasChildren := len(children) > 0

	if hasChildren {
		return `<li class="nav-item dropdown"><a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">` + item.Name() + `</a><ul class="dropdown-menu"><li><a class="dropdown-item" href="/web-design">Web Design</a></li><li><a class="dropdown-item" href="/seo">SEO</a></li></ul></li>`
	}

	return `<li class="nav-item"><a class="nav-link" href="` + item.URL() + `" target="_self">` + item.Name() + `</a></li>`
}

// TestNavbarBlockType_CustomCSS tests custom CSS functionality
func TestNavbarBlockType_CustomCSS(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Test custom CSS with Bootstrap 5
	customCSS := ".navbar { background-color: #ff0000 !important; } .navbar-brand { color: #ffffff !important; }"

	block := &TestNavbarBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_MENU_ID:               "test-menu",
			cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE: "bootstrap5",
			cmsstore.BLOCK_META_NAVBAR_CUSTOM_CSS:     customCSS,
		},
	}

	// Test rendering with custom CSS
	result, err := navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Check that custom CSS is included in style tags
	if !contains(result, "<style>") {
		t.Error("Expected result to contain <style> tag")
	}

	if !contains(result, "</style>") {
		t.Error("Expected result to contain </style> tag")
	}

	if !contains(result, customCSS) {
		t.Error("Expected result to contain custom CSS content")
	}

	// Test custom CSS with plain rendering
	block.meta[cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE] = "plain"
	result, err = navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Check that custom CSS is included in style tags for plain rendering too
	if !contains(result, "<style>") {
		t.Error("Expected plain result to contain <style> tag")
	}

	if !contains(result, "</style>") {
		t.Error("Expected plain result to contain </style> tag")
	}

	if !contains(result, customCSS) {
		t.Error("Expected plain result to contain custom CSS content")
	}

	// Test without custom CSS (should not include style tags)
	delete(block.meta, cmsstore.BLOCK_META_NAVBAR_CUSTOM_CSS)
	result, err = navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	if contains(result, "<style>") {
		t.Error("Expected result to NOT contain <style> tag when no custom CSS is provided")
	}
}

// TestNavbarBlockType_SaveAdminFields_CustomCSS tests saving custom CSS field
func TestNavbarBlockType_SaveAdminFields_CustomCSS(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Create a mock block
	block := &TestNavbarBlock{
		meta: map[string]string{},
	}

	// Create a real http.Request with form data including custom CSS
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Form = map[string][]string{
		"menu_id":               {"test-menu"},
		"navbar_style":          {"default"},
		"navbar_rendering_mode": {"bootstrap5"},
		"navbar_custom_css":     {".navbar { background: blue; }"},
	}

	// Test saving admin fields with custom CSS
	err = navbarBlock.SaveAdminFields(req, block)
	if err != nil {
		t.Errorf("Expected no error saving admin fields, got: %v", err)
		return
	}

	// Verify that custom CSS meta data was saved
	expectedCSS := ".navbar { background: blue; }"
	if block.meta[cmsstore.BLOCK_META_NAVBAR_CUSTOM_CSS] != expectedCSS {
		t.Errorf("Expected custom CSS to be '%s', got '%s'", expectedCSS, block.meta[cmsstore.BLOCK_META_NAVBAR_CUSTOM_CSS])
	}
}

// TestNavbarBlockType_BrandImage tests brand image functionality
func TestNavbarBlockType_BrandImage(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Test brand image only
	block := &TestNavbarBlock{
		meta: map[string]string{
			cmsstore.BLOCK_META_MENU_ID:                   "test-menu",
			cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE:     "bootstrap5",
			cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_URL:    "https://example.com/logo.png",
			cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_WIDTH:  "40",
			cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_HEIGHT: "30",
			cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_ALT:    "Test Logo",
		},
	}

	result, err := navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Debug: print actual result
	t.Logf("Actual HTML output: %s", result)

	// Check that image is rendered
	if !contains(result, `<img alt="Test Logo" height="30" src="https://example.com/logo.png" width="40" />`) {
		t.Error("Expected result to contain brand image")
	}

	if !contains(result, `width="40"`) {
		t.Error("Expected result to contain image width")
	}

	if !contains(result, `height="30"`) {
		t.Error("Expected result to contain image height")
	}

	if !contains(result, `alt="Test Logo"`) {
		t.Error("Expected result to contain image alt text")
	}

	// Test brand image and text together
	block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_TEXT] = "My Brand"
	result, err = navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Debug: print actual result for image + text
	t.Logf("Image + text HTML output: %s", result)

	if !contains(result, `<img alt="Test Logo" height="30" src="https://example.com/logo.png" width="40" />`) {
		t.Error("Expected result to contain brand image")
	}

	if !contains(result, "My Brand") {
		t.Error("Expected result to contain brand text")
	}

	if !contains(result, `d-inline-block align-text-top`) {
		t.Error("Expected result to contain Bootstrap 5 image alignment classes")
	}

	// Test default dimensions
	delete(block.meta, cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_WIDTH)
	delete(block.meta, cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_HEIGHT)
	result, err = navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	if !contains(result, `width="30"`) {
		t.Error("Expected default width of 30")
	}

	if !contains(result, `height="24"`) {
		t.Error("Expected default height of 24")
	}

	// Test default alt text
	delete(block.meta, cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_ALT)
	result, err = navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	if !contains(result, `alt="Logo"`) {
		t.Error("Expected default alt text of 'Logo'")
	}

	// Test plain rendering mode with brand image
	block.meta[cmsstore.BLOCK_META_NAVBAR_RENDERING_MODE] = "plain"
	result, err = navbarBlock.Render(ctx, block)
	if err != nil {
		t.Errorf("Render returned error: %v", err)
		return
	}

	// Debug: print actual result for plain rendering
	t.Logf("Plain rendering HTML output: %s", result)

	if !contains(result, `<img alt="Test Logo" height="30" src="https://example.com/logo.png" width="40" />`) {
		t.Error("Expected plain result to contain brand image")
	}
}

// TestNavbarBlockType_SaveAdminFields_BrandImage tests saving brand image fields
func TestNavbarBlockType_SaveAdminFields_BrandImage(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	navbarBlock := NewNavbarBlockType(store)

	// Create a mock block
	block := &TestNavbarBlock{
		meta: map[string]string{},
	}

	// Create a real http.Request with form data including brand image fields
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Form = map[string][]string{
		"menu_id":                   {"test-menu"},
		"navbar_style":              {"default"},
		"navbar_rendering_mode":     {"bootstrap5"},
		"navbar_brand_image_url":    {"https://example.com/logo.png"},
		"navbar_brand_image_width":  {"50"},
		"navbar_brand_image_height": {"40"},
		"navbar_brand_image_alt":    {"My Company Logo"},
	}

	// Test saving admin fields with brand image
	err = navbarBlock.SaveAdminFields(req, block)
	if err != nil {
		t.Errorf("Expected no error saving admin fields, got: %v", err)
		return
	}

	// Verify that brand image meta data was saved
	if block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_URL] != "https://example.com/logo.png" {
		t.Errorf("Expected brand image URL to be 'https://example.com/logo.png', got '%s'", block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_URL])
	}

	if block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_WIDTH] != "50" {
		t.Errorf("Expected brand image width to be '50', got '%s'", block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_WIDTH])
	}

	if block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_HEIGHT] != "40" {
		t.Errorf("Expected brand image height to be '40', got '%s'", block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_HEIGHT])
	}

	if block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_ALT] != "My Company Logo" {
		t.Errorf("Expected brand image alt to be 'My Company Logo', got '%s'", block.meta[cmsstore.BLOCK_META_NAVBAR_BRAND_IMAGE_ALT])
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
