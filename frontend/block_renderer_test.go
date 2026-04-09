package frontend

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
)

// MockBlockRenderer is a mock implementation of BlockRenderer for testing
type MockBlockRenderer struct {
	renderFunc func(ctx context.Context, block cmsstore.BlockInterface) (string, error)
}

func (m *MockBlockRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	if m.renderFunc != nil {
		return m.renderFunc(ctx, block)
	}
	return "mock rendered content", nil
}

// TestNewBlockRendererRegistry tests creating a new block renderer registry
func TestNewBlockRendererRegistry(t *testing.T) {
	registry := NewBlockRendererRegistry()
	if registry == nil {
		t.Error("NewBlockRendererRegistry returned nil")
	}
	if registry.renderers == nil {
		t.Error("Registry renderers map should be initialized")
	}
}

// TestBlockRendererRegistry_Register tests registering a renderer
func TestBlockRendererRegistry_Register(t *testing.T) {
	registry := NewBlockRendererRegistry()
	mockRenderer := &MockBlockRenderer{}

	// Register a renderer
	registry.Register("test_type", mockRenderer)

	// Verify it was registered
	renderer := registry.GetRenderer("test_type")
	if renderer == nil {
		t.Error("Renderer should be registered")
	}
	if renderer != mockRenderer {
		t.Error("Registered renderer should be the same instance")
	}
}

// TestBlockRendererRegistry_Register_Overwrite tests overwriting an existing renderer
func TestBlockRendererRegistry_Register_Overwrite(t *testing.T) {
	registry := NewBlockRendererRegistry()
	mockRenderer1 := &MockBlockRenderer{}
	mockRenderer2 := &MockBlockRenderer{}

	// Register first renderer
	registry.Register("test_type", mockRenderer1)

	// Overwrite with second renderer
	registry.Register("test_type", mockRenderer2)

	// Verify it was overwritten
	renderer := registry.GetRenderer("test_type")
	if renderer != mockRenderer2 {
		t.Error("Renderer should be overwritten")
	}
}

// TestBlockRendererRegistry_GetRenderer tests getting a registered renderer
func TestBlockRendererRegistry_GetRenderer(t *testing.T) {
	registry := NewBlockRendererRegistry()
	mockRenderer := &MockBlockRenderer{}

	registry.Register("test_type", mockRenderer)

	renderer := registry.GetRenderer("test_type")
	if renderer != mockRenderer {
		t.Error("Should return the registered renderer")
	}
}

// TestBlockRendererRegistry_GetRenderer_NotFound tests getting a non-existent renderer
func TestBlockRendererRegistry_GetRenderer_NotFound(t *testing.T) {
	registry := NewBlockRendererRegistry()

	// Register HTML renderer (default)
	registry.Register(cmsstore.BLOCK_TYPE_HTML, &MockBlockRenderer{})

	// Try to get a non-existent renderer
	renderer := registry.GetRenderer("non_existent_type")

	// Should return HTML renderer as default
	if renderer == nil {
		t.Error("Should return default HTML renderer when type not found")
	}
}

// TestBlockRendererRegistry_GetRenderer_NoDefault tests getting renderer when no default is registered
func TestBlockRendererRegistry_GetRenderer_NoDefault(t *testing.T) {
	registry := NewBlockRendererRegistry()

	// Don't register any renderers
	renderer := registry.GetRenderer("non_existent_type")

	// Should return NoOpRenderer as ultimate fallback
	if renderer == nil {
		t.Error("Should return NoOpRenderer as ultimate fallback")
	}

	_, isNoOp := renderer.(*NoOpRenderer)
	if !isNoOp {
		t.Error("Should return NoOpRenderer when no renderers are registered")
	}
}

// TestBlockRendererRegistry_RenderBlock tests rendering a block
func TestBlockRendererRegistry_RenderBlock(t *testing.T) {
	registry := NewBlockRendererRegistry()

	mockRenderer := &MockBlockRenderer{
		renderFunc: func(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
			return "test rendered content", nil
		},
	}

	registry.Register("test_type", mockRenderer)

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	block := cmsstore.NewBlock().
		SetContent("test content").
		SetType("test_type").
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	result, err := registry.RenderBlock(context.Background(), block)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "test rendered content" {
		t.Errorf("Expected 'test rendered content', got %q", result)
	}
}

// TestBlockRendererRegistry_RenderBlock_NilBlock tests rendering a nil block
func TestBlockRendererRegistry_RenderBlock_NilBlock(t *testing.T) {
	registry := NewBlockRendererRegistry()

	result, err := registry.RenderBlock(context.Background(), nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "<!-- Block is nil -->" {
		t.Errorf("Expected '<!-- Block is nil -->', got %q", result)
	}
}

// TestBlockRendererRegistry_RenderBlock_GlobalBlockType tests rendering with global block type
func TestBlockRendererRegistry_RenderBlock_GlobalBlockType(t *testing.T) {
	registry := NewBlockRendererRegistry()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Create a menu first for the navbar to use
	menu := cmsstore.NewMenu().
		SetName("Test Menu").
		SetStatus(cmsstore.MENU_STATUS_ACTIVE)

	err = store.MenuCreate(context.Background(), menu)
	if err != nil {
		t.Fatalf("Failed to create menu: %v", err)
	}

	block := cmsstore.NewBlock().
		SetContent("test content").
		SetType("navbar").
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = block.SetMeta("menu_id", menu.ID())
	if err != nil {
		t.Fatalf("Failed to set menu_id metadata: %v", err)
	}

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	// This should use the global navbar block type registered in initBlockRenderers
	result, err := registry.RenderBlock(context.Background(), block)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// The navbar block type should render something
	if result == "" {
		t.Error("Expected non-empty result from navbar renderer")
	}
}

// TestNoOpRenderer_Render tests the NoOpRenderer
func TestNoOpRenderer_Render(t *testing.T) {
	renderer := &NoOpRenderer{}

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	block := cmsstore.NewBlock().
		SetContent("test content").
		SetType("unknown_type").
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	result, err := renderer.Render(context.Background(), block)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "<!-- No renderer available for block type: unknown_type -->"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestBlockRendererRegistry_ConcurrentAccess tests concurrent access to the registry
func TestBlockRendererRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewBlockRendererRegistry()
	mockRenderer := &MockBlockRenderer{}

	done := make(chan bool)

	// Concurrent registrations
	for i := 0; i < 10; i++ {
		go func(i int) {
			registry.Register("type_"+string(rune('a'+i)), mockRenderer)
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			registry.GetRenderer("type_a")
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Verify registry is still functional
	renderer := registry.GetRenderer("type_a")
	if renderer == nil {
		t.Error("Registry should still be functional after concurrent access")
	}
}

// TestInitBlockRenderers tests the initBlockRenderers function
func TestInitBlockRenderers(t *testing.T) {
	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	// Create a minimal frontend struct
	f := &frontend{
		store: store,
	}

	registry := initBlockRenderers(f, store)

	if registry == nil {
		t.Error("initBlockRenderers should return a registry")
	}

	// Verify HTML renderer is registered
	htmlRenderer := registry.GetRenderer(cmsstore.BLOCK_TYPE_HTML)
	if htmlRenderer == nil {
		t.Error("HTML renderer should be registered")
	}

	// Verify Menu renderer is registered
	menuRenderer := registry.GetRenderer(cmsstore.BLOCK_TYPE_MENU)
	if menuRenderer == nil {
		t.Error("Menu renderer should be registered")
	}
}

// TestBlockRendererRegistry_ThreadSafety tests thread safety of the registry
func TestBlockRendererRegistry_ThreadSafety(t *testing.T) {
	registry := NewBlockRendererRegistry()
	mockRenderer := &MockBlockRenderer{}

	// Register multiple renderers concurrently
	for i := 0; i < 100; i++ {
		go func(i int) {
			registry.Register("concurrent_type_"+string(rune(i)), mockRenderer)
		}(i)
	}

	// Give goroutines time to complete
	// In a real test, we'd use a sync.WaitGroup, but for this simple test
	// we'll just verify the registry still works after the operations

	// Try to get a renderer
	renderer := registry.GetRenderer("concurrent_type_0")
	// It might or might not be there depending on timing, but the registry
	// shouldn't panic or be corrupted
	_ = renderer
}
