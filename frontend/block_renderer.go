package frontend

import (
	"context"
	"sync"

	"github.com/dracory/cmsstore"
	"github.com/dracory/cmsstore/frontend/blocks/html"
	"github.com/dracory/cmsstore/frontend/blocks/menu"
)

// BlockRenderer interface defines how different block types are rendered.
//
// Custom block renderers must implement this interface to integrate with the
// block rendering system. Each renderer is responsible for converting a block's
// content into HTML output.
//
// Example implementation:
//
//	type GalleryRenderer struct {
//	    store FrontendStore
//	}
//
//	func (r *GalleryRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
//	    // Parse block content and metadata
//	    images := parseGalleryImages(block.Content())
//	    layout := block.Meta("layout")
//
//	    // Generate HTML
//	    return renderGalleryHTML(images, layout), nil
//	}
//
// To register a custom renderer:
//
//	frontend := cmsstore.NewFrontend(...)
//	frontend.BlockRegistry().Register("gallery", &GalleryRenderer{store: frontend})
//
// See frontend/blocks/README.md for detailed examples and best practices.
type BlockRenderer interface {
	// Render renders the block content and returns the HTML
	Render(ctx context.Context, block cmsstore.BlockInterface) (string, error)
}

// BlockRendererRegistry manages all registered block renderers.
//
// The registry is thread-safe and can be accessed concurrently. It maps block
// types (strings) to their corresponding BlockRenderer implementations.
//
// Built-in block types:
//   - cmsstore.BLOCK_TYPE_HTML: Renders raw HTML content
//   - cmsstore.BLOCK_TYPE_MENU: Renders navigation menus
//
// Custom block types can be registered after frontend initialization:
//
//	frontend.BlockRegistry().Register("custom_type", customRenderer)
type BlockRendererRegistry struct {
	renderers map[string]BlockRenderer
	mu        sync.RWMutex
}

// NewBlockRendererRegistry creates a new registry for block renderers.
// This is typically called internally during frontend initialization.
func NewBlockRendererRegistry() *BlockRendererRegistry {
	return &BlockRendererRegistry{
		renderers: make(map[string]BlockRenderer),
	}
}

// Register registers a renderer for a specific block type.
//
// The blockType should match the value stored in the block's Type() field.
// If a renderer already exists for the given type, it will be replaced.
//
// This method is thread-safe and can be called after frontend initialization
// to add custom block types.
//
// Example:
//
//	frontend.BlockRegistry().Register("video", &VideoRenderer{})
//	frontend.BlockRegistry().Register("carousel", &CarouselRenderer{store: frontend})
func (r *BlockRendererRegistry) Register(blockType string, renderer BlockRenderer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.renderers[blockType] = renderer
}

// GetRenderer returns the renderer for the given block type
func (r *BlockRendererRegistry) GetRenderer(blockType string) BlockRenderer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if renderer, exists := r.renderers[blockType]; exists && renderer != nil {
		return renderer
	}
	// Return HTML renderer as default if it exists
	if renderer, exists := r.renderers[cmsstore.BLOCK_TYPE_HTML]; exists && renderer != nil {
		return renderer
	}
	// Return no-op renderer as ultimate fallback
	return &NoOpRenderer{}
}

// RenderBlock renders a block using the appropriate renderer
func (r *BlockRendererRegistry) RenderBlock(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	if block == nil {
		return "<!-- Block is nil -->", nil
	}
	renderer := r.GetRenderer(block.Type())
	return renderer.Render(ctx, block)
}

// NoOpRenderer is a fallback renderer that returns empty content
type NoOpRenderer struct{}

// Render implements BlockRenderer interface
func (r *NoOpRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	return "<!-- No renderer available for block type: " + block.Type() + " -->", nil
}

// initBlockRenderers initializes and registers all block renderers
func initBlockRenderers(f *frontend) *BlockRendererRegistry {
	registry := NewBlockRendererRegistry()

	// Register HTML renderer (default)
	registry.Register(cmsstore.BLOCK_TYPE_HTML, html.NewHTMLRenderer())

	// Register Menu renderer
	registry.Register(cmsstore.BLOCK_TYPE_MENU, menu.NewBlockRenderer(f))

	return registry
}
