package html

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/testutils"
)

// TestNewHTMLRenderer tests the NewHTMLRenderer function
func TestNewHTMLRenderer(t *testing.T) {
	renderer := NewHTMLRenderer()
	if renderer == nil {
		t.Error("NewHTMLRenderer returned nil")
	}
}

// TestHTMLRenderer_Render_NilBlock tests rendering with a nil block
func TestHTMLRenderer_Render_NilBlock(t *testing.T) {
	renderer := NewHTMLRenderer()
	result, err := renderer.Render(context.Background(), nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "<!-- Block is nil -->"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestHTMLRenderer_Render_EmptyContent tests rendering with empty content
func TestHTMLRenderer_Render_EmptyContent(t *testing.T) {
	renderer := NewHTMLRenderer()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	block := cmsstore.NewBlock().
		SetContent("").
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	result, err := renderer.Render(context.Background(), block)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "<!-- Empty block content -->"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestHTMLRenderer_Render_ValidContent tests rendering with valid HTML content
func TestHTMLRenderer_Render_ValidContent(t *testing.T) {
	renderer := NewHTMLRenderer()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	content := "<div>Hello World</div>"
	block := cmsstore.NewBlock().
		SetContent(content).
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	result, err := renderer.Render(context.Background(), block)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

// TestHTMLRenderer_Render_ComplexHTML tests rendering with complex HTML
func TestHTMLRenderer_Render_ComplexHTML(t *testing.T) {
	renderer := NewHTMLRenderer()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	content := `<div class="container">
		<h1>Title</h1>
		<p>Paragraph with <strong>bold</strong> text</p>
		<ul>
			<li>Item 1</li>
			<li>Item 2</li>
		</ul>
	</div>`
	block := cmsstore.NewBlock().
		SetContent(content).
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	result, err := renderer.Render(context.Background(), block)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

// TestHTMLRenderer_Render_WithScriptTags tests rendering HTML with script tags
func TestHTMLRenderer_Render_WithScriptTags(t *testing.T) {
	renderer := NewHTMLRenderer()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	content := `<script>console.log("test");</script><div>Content</div>`
	block := cmsstore.NewBlock().
		SetContent(content).
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	result, err := renderer.Render(context.Background(), block)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

// TestHTMLRenderer_Render_WhitespaceContent tests rendering with whitespace-only content
func TestHTMLRenderer_Render_WhitespaceContent(t *testing.T) {
	renderer := NewHTMLRenderer()

	store, err := testutils.InitStore(":memory:")
	if err != nil {
		t.Fatalf("Failed to init store: %v", err)
	}

	content := "   \n\t  "
	block := cmsstore.NewBlock().
		SetContent(content).
		SetType(cmsstore.BLOCK_TYPE_HTML).
		SetStatus(cmsstore.BLOCK_STATUS_ACTIVE)

	err = store.BlockCreate(context.Background(), block)
	if err != nil {
		t.Fatalf("Failed to create block: %v", err)
	}

	result, err := renderer.Render(context.Background(), block)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Whitespace is not empty string, so it should be returned as-is
	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

// TestHTMLRenderer_Render_ImplementsInterface tests that HTMLRenderer implements cmsstore.BlockRenderer
func TestHTMLRenderer_Render_ImplementsInterface(t *testing.T) {
	renderer := NewHTMLRenderer()

	// This test ensures that HTMLRenderer can be used wherever BlockRenderer is expected
	var _ interface {
		Render(ctx context.Context, block cmsstore.BlockInterface) (string, error)
	} = renderer
}
