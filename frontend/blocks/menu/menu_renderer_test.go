package menu

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/dracory/cmsstore"
	"github.com/dromara/carbon/v2"
)

// TestGetCurrentPath verifies that getCurrentPath extracts the path from request context
func TestGetCurrentPath(t *testing.T) {
	// Test with request in context
	req, _ := http.NewRequest("GET", "/about?q=test", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)

	path := getCurrentPath(ctx)
	if path != "/about" {
		t.Errorf("Expected path '/about', got: %s", path)
	}

	// Test without request in context
	ctx2 := context.Background()
	path2 := getCurrentPath(ctx2)
	if path2 != "" {
		t.Errorf("Expected empty path when no request, got: %s", path2)
	}
}

// TestIsActiveItem verifies the active item detection logic
func TestIsActiveItem(t *testing.T) {
	tests := []struct {
		name        string
		itemURL     string
		currentPath string
		expected    bool
	}{
		{"exact match", "/about", "/about", true},
		{"trailing slash on item", "/about/", "/about", true},
		{"trailing slash on current", "/about", "/about/", true},
		{"different paths", "/about", "/contact", false},
		{"root match", "/", "/", true},
		{"empty item URL", "", "/about", false},
		{"empty current path", "/about", "", false},
		{"both empty", "", "", false},
		{"partial match", "/about", "/about-us", false},
		{"subpath", "/about", "/about/team", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isActiveItem(test.itemURL, test.currentPath)
			if result != test.expected {
				t.Errorf("isActiveItem(%q, %q) = %v, want %v", test.itemURL, test.currentPath, result, test.expected)
			}
		})
	}
}

// TestHasActiveChild verifies the recursive active child detection
func TestHasActiveChild(t *testing.T) {
	// Create a simple tree: Parent -> Child -> Grandchild
	grandchild := &menuTreeNode{
		Item: &mockMenuItem{id: "3", url: "/about/team"},
	}
	child := &menuTreeNode{
		Item:     &mockMenuItem{id: "2", url: "/about"},
		Children: []*menuTreeNode{grandchild},
	}
	parent := &menuTreeNode{
		Item:     &mockMenuItem{id: "1", url: "/"},
		Children: []*menuTreeNode{child},
	}

	// Test: grandchild URL matches
	if !hasActiveChild(parent, "/about/team") {
		t.Error("Expected parent to have active grandchild")
	}

	// Test: child URL matches
	if !hasActiveChild(parent, "/about") {
		t.Error("Expected parent to have active child")
	}

	// Test: parent URL matches
	if !hasActiveChild(parent, "/") {
		t.Error("Expected parent to be active")
	}

	// Test: no match
	if hasActiveChild(parent, "/contact") {
		t.Error("Expected no active items for /contact")
	}
}

// TestRenderMenuItemHTML_ActiveState verifies active CSS classes are added
func TestRenderMenuItemHTML_ActiveState(t *testing.T) {
	renderer := &MenuRenderer{}

	node := &menuTreeNode{
		Item: &mockMenuItem{id: "1", name: "About", url: "/about"},
	}

	// Create request context for /about
	req, _ := http.NewRequest("GET", "/about", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)

	html := renderer.renderMenuItemHTML(ctx, node, false)

	// Check li has active class
	if !strings.Contains(html, `<li class="active">`) {
		t.Errorf("Expected li to have 'active' class. HTML: %s", html)
	}

	// Check a has active class
	if !strings.Contains(html, `<a href="/about" class="active">About</a>`) {
		t.Errorf("Expected a to have 'active' class. HTML: %s", html)
	}
}

// TestRenderMenuItemHTML_InactiveState verifies no active class when not matching
func TestRenderMenuItemHTML_InactiveState(t *testing.T) {
	renderer := &MenuRenderer{}

	node := &menuTreeNode{
		Item: &mockMenuItem{id: "1", name: "About", url: "/about"},
	}

	// Create request context for different path
	req, _ := http.NewRequest("GET", "/contact", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)

	html := renderer.renderMenuItemHTML(ctx, node, false)

	// Check no active class
	if strings.Contains(html, `class="active"`) {
		t.Errorf("Expected no 'active' class when paths don't match. HTML: %s", html)
	}

	// Verify the link still renders
	if !strings.Contains(html, `<a href="/about">About</a>`) {
		t.Errorf("Expected link to render without active class. HTML: %s", html)
	}
}

// TestRenderMenuItemHTML_ActiveParent verifies active-parent class for dropdowns
func TestRenderMenuItemHTML_ActiveParent(t *testing.T) {
	renderer := &MenuRenderer{}

	// Create parent with active child
	child := &menuTreeNode{
		Item: &mockMenuItem{id: "2", name: "Team", url: "/about/team"},
	}
	parent := &menuTreeNode{
		Item:     &mockMenuItem{id: "1", name: "About", url: "/about"},
		Children: []*menuTreeNode{child},
	}

	// Create request context for child path
	req, _ := http.NewRequest("GET", "/about/team", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)

	// Render parent without children (simulates dropdown parent view)
	html := renderer.renderMenuItemHTML(ctx, parent, false)

	// Check parent has active-parent class (since child is active but we're not rendering children)
	if !strings.Contains(html, `class="active-parent"`) {
		t.Errorf("Expected parent to have 'active-parent' class. HTML: %s", html)
	}
}

// mockMenuItem is a test implementation of MenuItemInterface
type mockMenuItem struct {
	id        string
	name      string
	url       string
	parentID  string
	target    string
	menuID    string
	pageID    string
	status    string
	sequence  int
	createdAt string
	updatedAt string
	handle    string
	memo      string
	meta      map[string]string
}

// Core methods used by renderer
func (m *mockMenuItem) ID() string       { return m.id }
func (m *mockMenuItem) Name() string     { return m.name }
func (m *mockMenuItem) URL() string      { return m.url }
func (m *mockMenuItem) ParentID() string { return m.parentID }
func (m *mockMenuItem) Target() string   { return m.target }

// Setters
func (m *mockMenuItem) SetID(id string) cmsstore.MenuItemInterface     { m.id = id; return m }
func (m *mockMenuItem) SetName(name string) cmsstore.MenuItemInterface { m.name = name; return m }
func (m *mockMenuItem) SetURL(url string) cmsstore.MenuItemInterface   { m.url = url; return m }
func (m *mockMenuItem) SetParentID(parentID string) cmsstore.MenuItemInterface {
	m.parentID = parentID
	return m
}
func (m *mockMenuItem) SetTarget(target string) cmsstore.MenuItemInterface {
	m.target = target
	return m
}
func (m *mockMenuItem) SetMenuID(menuID string) cmsstore.MenuItemInterface {
	m.menuID = menuID
	return m
}
func (m *mockMenuItem) SetPageID(pageID string) cmsstore.MenuItemInterface {
	m.pageID = pageID
	return m
}
func (m *mockMenuItem) SetStatus(status string) cmsstore.MenuItemInterface {
	m.status = status
	return m
}
func (m *mockMenuItem) SetSequenceInt(sequence int) cmsstore.MenuItemInterface {
	m.sequence = sequence
	return m
}
func (m *mockMenuItem) SetSequence(sequence string) cmsstore.MenuItemInterface { return m }
func (m *mockMenuItem) SetHandle(handle string) cmsstore.MenuItemInterface {
	m.handle = handle
	return m
}
func (m *mockMenuItem) SetMemo(memo string) cmsstore.MenuItemInterface { m.memo = memo; return m }
func (m *mockMenuItem) SetCreatedAt(createdAt string) cmsstore.MenuItemInterface {
	m.createdAt = createdAt
	return m
}
func (m *mockMenuItem) SetUpdatedAt(updatedAt string) cmsstore.MenuItemInterface {
	m.updatedAt = updatedAt
	return m
}
func (m *mockMenuItem) SetSoftDeletedAt(at string) cmsstore.MenuItemInterface { return m }

// Getters (stubs)
func (m *mockMenuItem) MenuID() string                      { return m.menuID }
func (m *mockMenuItem) PageID() string                      { return m.pageID }
func (m *mockMenuItem) Status() string                      { return m.status }
func (m *mockMenuItem) Sequence() string                    { return "" }
func (m *mockMenuItem) SequenceInt() int                    { return m.sequence }
func (m *mockMenuItem) Handle() string                      { return m.handle }
func (m *mockMenuItem) Memo() string                        { return m.memo }
func (m *mockMenuItem) CreatedAt() string                   { return m.createdAt }
func (m *mockMenuItem) UpdatedAt() string                   { return m.updatedAt }
func (m *mockMenuItem) SoftDeletedAt() string               { return "" }
func (m *mockMenuItem) CreatedAtCarbon() *carbon.Carbon     { return nil }
func (m *mockMenuItem) UpdatedAtCarbon() *carbon.Carbon     { return nil }
func (m *mockMenuItem) SoftDeletedAtCarbon() *carbon.Carbon { return nil }
func (m *mockMenuItem) IsActive() bool                      { return m.status == "active" }
func (m *mockMenuItem) IsInactive() bool                    { return false }
func (m *mockMenuItem) IsSoftDeleted() bool                 { return false }

// Data methods (stubs)
func (m *mockMenuItem) Data() map[string]string              { return nil }
func (m *mockMenuItem) DataChanged() map[string]string       { return nil }
func (m *mockMenuItem) MarkAsNotDirty()                      {}
func (m *mockMenuItem) MarshalToVersioning() (string, error) { return "", nil }

// Meta methods
func (m *mockMenuItem) Meta(key string) string {
	if m.meta != nil {
		return m.meta[key]
	}
	return ""
}
func (m *mockMenuItem) SetMeta(key, value string) error {
	if m.meta == nil {
		m.meta = make(map[string]string)
	}
	m.meta[key] = value
	return nil
}
func (m *mockMenuItem) Metas() (map[string]string, error)      { return m.meta, nil }
func (m *mockMenuItem) SetMetas(metas map[string]string) error { m.meta = metas; return nil }
func (m *mockMenuItem) UpsertMetas(metas map[string]string) error {
	if m.meta == nil {
		m.meta = make(map[string]string)
	}
	for k, v := range metas {
		m.meta[k] = v
	}
	return nil
}

// Soft delete setters (stubs)
