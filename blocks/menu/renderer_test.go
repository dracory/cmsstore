package menu

import (
	"context"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
	"github.com/dromara/carbon/v2"
	_ "modernc.org/sqlite"
)

// TestMenuItem is a minimal implementation that only provides what the renderer needs
type TestMenuItem struct {
	name string
	url  string
	data map[string]string
}

func (m *TestMenuItem) Data() map[string]string {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	return m.data
}
func (m *TestMenuItem) DataChanged() map[string]string                                   { return make(map[string]string) }
func (m *TestMenuItem) MarkAsNotDirty()                                                  {}
func (m *TestMenuItem) ID() string                                                       { return "test-id" }
func (m *TestMenuItem) Name() string                                                     { return m.name }
func (m *TestMenuItem) URL() string                                                      { return m.url }
func (m *TestMenuItem) IsActive() bool                                                   { return true }
func (m *TestMenuItem) IsInactive() bool                                                 { return false }
func (m *TestMenuItem) IsSoftDeleted() bool                                              { return false }
func (m *TestMenuItem) MarshalToVersioning() (string, error)                             { return "", nil } // Simplified for testing
func (m *TestMenuItem) SetID(id string) cmsstore.MenuItemInterface                       { return m }
func (m *TestMenuItem) SetName(name string) cmsstore.MenuItemInterface                   { m.name = name; return m }
func (m *TestMenuItem) SetURL(url string) cmsstore.MenuItemInterface                     { m.url = url; return m }
func (m *TestMenuItem) SetCreatedAt(createdAt string) cmsstore.MenuItemInterface         { return m }
func (m *TestMenuItem) SetUpdatedAt(updatedAt string) cmsstore.MenuItemInterface         { return m }
func (m *TestMenuItem) SetSoftDeletedAt(softDeletedAt string) cmsstore.MenuItemInterface { return m }
func (m *TestMenuItem) SetHandle(handle string) cmsstore.MenuItemInterface               { return m }
func (m *TestMenuItem) SetMemo(memo string) cmsstore.MenuItemInterface                   { return m }
func (m *TestMenuItem) SetMenuID(menuID string) cmsstore.MenuItemInterface               { return m }
func (m *TestMenuItem) SetMeta(key, value string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[key] = value
	return nil
}
func (m *TestMenuItem) SetMetas(metas map[string]string) error { m.data = metas; return nil }
func (m *TestMenuItem) UpsertMetas(metas map[string]string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	for k, v := range metas {
		m.data[k] = v
	}
	return nil
}
func (m *TestMenuItem) SetPageID(pageID string) cmsstore.MenuItemInterface     { return m }
func (m *TestMenuItem) SetParentID(parentID string) cmsstore.MenuItemInterface { return m }
func (m *TestMenuItem) SetSequence(sequence string) cmsstore.MenuItemInterface { return m }
func (m *TestMenuItem) SetSequenceInt(sequence int) cmsstore.MenuItemInterface { return m }
func (m *TestMenuItem) SetStatus(status string) cmsstore.MenuItemInterface     { return m }
func (m *TestMenuItem) SetTarget(target string) cmsstore.MenuItemInterface     { return m }
func (m *TestMenuItem) CreatedAt() string                                      { return "2023-01-01 00:00:00" }
func (m *TestMenuItem) CreatedAtCarbon() *carbon.Carbon                        { return carbon.Parse("2023-01-01 00:00:00") }
func (m *TestMenuItem) UpdatedAt() string                                      { return "2023-01-01 00:00:00" }
func (m *TestMenuItem) UpdatedAtCarbon() *carbon.Carbon                        { return carbon.Parse("2023-01-01 00:00:00") }
func (m *TestMenuItem) SoftDeletedAt() string                                  { return "" }
func (m *TestMenuItem) SoftDeletedAtCarbon() *carbon.Carbon                    { return nil }
func (m *TestMenuItem) Handle() string                                         { return "" }
func (m *TestMenuItem) Memo() string                                           { return "" }
func (m *TestMenuItem) MenuID() string                                         { return "test-menu" }
func (m *TestMenuItem) Meta(key string) string {
	if m.data == nil {
		return ""
	}
	return m.data[key]
}
func (m *TestMenuItem) Metas() (map[string]string, error) {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	return m.data, nil
}
func (m *TestMenuItem) PageID() string   { return "" }
func (m *TestMenuItem) ParentID() string { return "" }
func (m *TestMenuItem) Sequence() string { return "1" }
func (m *TestMenuItem) SequenceInt() int { return 1 }
func (m *TestMenuItem) Status() string   { return "active" }
func (m *TestMenuItem) Target() string   { return "_self" }

// TestRenderMenuHTMLBasic tests the core rendering functionality
func TestRenderMenuHTMLBasic(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Create test menu items with minimal implementation
	menuItems := []cmsstore.MenuItemInterface{
		&TestMenuItem{name: "Home", url: "/"},
		&TestMenuItem{name: "About", url: "/about"},
		&TestMenuItem{name: "Contact", url: "/contact"},
	}

	tests := []struct {
		name          string
		style         string
		renderingMode string
		cssClass      string
		cssID         string
		startLevel    int
		maxDepth      int
		expectedHTML  string
	}{
		{
			name:          "Basic vertical menu",
			style:         "vertical",
			renderingMode: "plain",
			cssClass:      "nav",
			cssID:         "main-nav",
			startLevel:    0,
			maxDepth:      0,
			expectedHTML: `<nav class="menu menu-style-vertical nav" id="main-nav">` +
				`<a href="/">Home</a>` +
				`<a href="/about">About</a>` +
				`<a href="/contact">Contact</a>` +
				`</nav>`,
		},
		{
			name:          "Horizontal menu without CSS class and ID",
			style:         "horizontal",
			renderingMode: "plain",
			cssClass:      "",
			cssID:         "",
			startLevel:    0,
			maxDepth:      0,
			expectedHTML: `<nav class="menu menu-style-horizontal">` +
				`<a href="/">Home</a>` +
				`<a href="/about">About</a>` +
				`<a href="/contact">Contact</a>` +
				`</nav>`,
		},
		{
			name:          "Dropdown menu with only CSS class",
			style:         "dropdown",
			renderingMode: "plain",
			cssClass:      "dropdown-menu",
			cssID:         "",
			startLevel:    0,
			maxDepth:      0,
			expectedHTML: `<nav class="menu menu-style-dropdown dropdown-menu">` +
				`<a href="/">Home</a>` +
				`<a href="/about">About</a>` +
				`<a href="/contact">Contact</a>` +
				`</nav>`,
		},
		{
			name:          "Bootstrap 5 dropdown",
			style:         "dropdown",
			renderingMode: "bootstrap5",
			cssClass:      "my-dropdown",
			cssID:         "main-dropdown",
			startLevel:    0,
			maxDepth:      0,
			expectedHTML: `<div class="dropdown my-dropdown" id="main-dropdown">` +
				`<button aria-expanded="false" class="btn btn-secondary dropdown-toggle" data-bs-toggle="dropdown" type="button">Dropdown</button>` +
				`<div class="dropdown-menu">` +
				`<a class="dropdown-item" href="/">Home</a>` +
				`<a class="dropdown-item" href="/about">About</a>` +
				`<a class="dropdown-item" href="/contact">Contact</a>` +
				`</div>` +
				`</div>`,
		},
		{
			name:         "Empty menu items",
			style:        "vertical",
			cssClass:     "",
			cssID:        "",
			startLevel:   0,
			maxDepth:     0,
			expectedHTML: `<nav class="menu menu-style-vertical"></nav>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testMenuItems []cmsstore.MenuItemInterface
			if tt.name != "Empty menu items" {
				testMenuItems = menuItems
			}

			result, err := renderMenuHTML(ctx, store, testMenuItems, tt.style, tt.renderingMode, tt.cssClass, tt.cssID, tt.startLevel, tt.maxDepth)

			if err != nil {
				t.Errorf("renderMenuHTML() returned error: %v", err)
				return
			}

			if result != tt.expectedHTML {
				t.Errorf("renderMenuHTML() = %q, want %q", result, tt.expectedHTML)
			}
		})
	}
}

func TestRenderMenuHTMLWithSpecialCharacters(t *testing.T) {
	ctx := context.Background()
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Test with special characters in menu item names and URLs
	menuItems := []cmsstore.MenuItemInterface{
		&TestMenuItem{name: "Home & About", url: "/home?param=value&other=test"},
		&TestMenuItem{name: "Products > Services", url: "/products/services"},
	}

	result, err := renderMenuHTML(ctx, store, menuItems, "vertical", "plain", "test-class", "test-id", 0, 0)

	if err != nil {
		t.Errorf("renderMenuHTML() returned error: %v", err)
		return
	}

	// The hb library should handle HTML escaping
	expected := `<nav class="menu menu-style-vertical test-class" id="test-id">` +
		`<a href="/home?param=value&amp;other=test">Home &amp; About</a>` +
		`<a href="/products/services">Products &gt; Services</a>` +
		`</nav>`

	if result != expected {
		t.Errorf("renderMenuHTML() with special characters = %q, want %q", result, expected)
	}
}
