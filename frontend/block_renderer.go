package frontend

import (
	"context"

	"github.com/dracory/cmsstore"
)

// BlockRenderer interface defines how different block types are rendered
type BlockRenderer interface {
	// Render renders the block content and returns the HTML
	Render(ctx context.Context, block cmsstore.BlockInterface) (string, error)
}

// BlockRendererRegistry manages all registered block renderers
type BlockRendererRegistry struct {
	renderers map[string]BlockRenderer
}

// NewBlockRendererRegistry creates a new registry for block renderers
func NewBlockRendererRegistry() *BlockRendererRegistry {
	return &BlockRendererRegistry{
		renderers: make(map[string]BlockRenderer),
	}
}

// Register registers a renderer for a specific block type
func (r *BlockRendererRegistry) Register(blockType string, renderer BlockRenderer) {
	r.renderers[blockType] = renderer
}

// GetRenderer returns the renderer for the given block type
func (r *BlockRendererRegistry) GetRenderer(blockType string) BlockRenderer {
	if renderer, exists := r.renderers[blockType]; exists {
		return renderer
	}
	// Return HTML renderer as default
	return r.renderers[cmsstore.BLOCK_TYPE_HTML]
}

// RenderBlock renders a block using the appropriate renderer
func (r *BlockRendererRegistry) RenderBlock(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
	renderer := r.GetRenderer(block.Type())
	return renderer.Render(ctx, block)
}

// initBlockRenderers initializes and registers all block renderers
func initBlockRenderers(f *frontend) *BlockRendererRegistry {
	registry := NewBlockRendererRegistry()

	// Register HTML renderer (default)
	registry.Register(cmsstore.BLOCK_TYPE_HTML, NewHTMLBlockRenderer())

	// Register Menu renderer
	registry.Register(cmsstore.BLOCK_TYPE_MENU, NewMenuBlockRenderer(f))

	return registry
}
