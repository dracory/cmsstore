package cmsstore

import (
	"context"
	"net/http"
)

// BlockTypeAdapter wraps separate renderer and admin provider into a unified BlockType.
//
// This adapter allows the new unified BlockType system to work with code that
// still uses the separate BlockRenderer and BlockAdminFieldProvider interfaces.
//
// This is useful for:
//   - Backward compatibility with existing code
//   - Gradual migration to the unified system
//   - Wrapping built-in block types (HTML, Menu)
//
// Example:
//
//	renderer := &HTMLRenderer{}
//	adminProvider := &HTMLAdminProvider{}
//	blockType := cmsstore.NewBlockTypeAdapter("html", "HTML Block", renderer, adminProvider)
//	cmsstore.RegisterBlockType(blockType)
type BlockTypeAdapter struct {
	typeKey       string
	typeLabel     string
	renderer      BlockRenderer
	adminProvider BlockAdminFieldProvider
}

// BlockRenderer interface for frontend rendering (from frontend package).
type BlockRenderer interface {
	Render(ctx context.Context, block BlockInterface) (string, error)
}

// BlockAdminFieldProvider interface for admin UI (from admin/blocks package).
type BlockAdminFieldProvider interface {
	GetContentFields(block BlockInterface, r *http.Request) interface{}
	GetTypeLabel() string
	SaveContentFields(r *http.Request, block BlockInterface) error
}

// NewBlockTypeAdapter creates a new adapter that combines a renderer and admin provider.
//
// Parameters:
//   - typeKey: Unique identifier for the block type (e.g., "html", "menu")
//   - typeLabel: Display name for the block type (e.g., "HTML Block")
//   - renderer: Frontend renderer implementation
//   - adminProvider: Admin UI provider implementation
//
// Returns:
//   - A BlockType that can be registered with RegisterBlockType()
func NewBlockTypeAdapter(typeKey, typeLabel string, renderer BlockRenderer, adminProvider BlockAdminFieldProvider) BlockType {
	return &BlockTypeAdapter{
		typeKey:       typeKey,
		typeLabel:     typeLabel,
		renderer:      renderer,
		adminProvider: adminProvider,
	}
}

// TypeKey returns the unique identifier for this block type.
func (a *BlockTypeAdapter) TypeKey() string {
	return a.typeKey
}

// TypeLabel returns the human-readable display name.
func (a *BlockTypeAdapter) TypeLabel() string {
	return a.typeLabel
}

// Render delegates to the wrapped renderer.
// Attributes are ignored for legacy adapters.
func (a *BlockTypeAdapter) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
	if a.renderer == nil {
		return "<!-- No renderer configured -->", nil
	}
	return a.renderer.Render(ctx, block)
}

// GetAdminFields delegates to the wrapped admin provider.
func (a *BlockTypeAdapter) GetAdminFields(block BlockInterface, r *http.Request) interface{} {
	if a.adminProvider == nil {
		return []interface{}{}
	}
	return a.adminProvider.GetContentFields(block, r)
}

// SaveAdminFields delegates to the wrapped admin provider.
func (a *BlockTypeAdapter) SaveAdminFields(r *http.Request, block BlockInterface) error {
	if a.adminProvider == nil {
		return nil
	}
	return a.adminProvider.SaveContentFields(r, block)
}

// GetCustomVariables returns nil as the adapter has no knowledge of custom variables.
// Implement GetCustomVariables directly on your BlockType for variable metadata.
func (a *BlockTypeAdapter) GetCustomVariables() []BlockCustomVariable {
	return nil
}

// RendererOnlyAdapter creates a BlockType from just a renderer.
//
// This is useful when you only need frontend rendering and don't need admin UI.
// The admin UI will fall back to a basic textarea editor.
//
// Example:
//
//	renderer := &CustomRenderer{}
//	blockType := cmsstore.RendererOnlyAdapter("custom", "Custom Block", renderer)
//	cmsstore.RegisterBlockType(blockType)
func RendererOnlyAdapter(typeKey, typeLabel string, renderer BlockRenderer) BlockType {
	return &BlockTypeAdapter{
		typeKey:   typeKey,
		typeLabel: typeLabel,
		renderer:  renderer,
	}
}

// AdminOnlyAdapter creates a BlockType from just an admin provider.
//
// This is useful when you only need admin UI and the frontend rendering
// is handled elsewhere (e.g., by a separate system).
//
// Example:
//
//	adminProvider := &CustomAdminProvider{}
//	blockType := cmsstore.AdminOnlyAdapter("custom", "Custom Block", adminProvider)
//	cmsstore.RegisterBlockType(blockType)
func AdminOnlyAdapter(typeKey, typeLabel string, adminProvider BlockAdminFieldProvider) BlockType {
	return &BlockTypeAdapter{
		typeKey:       typeKey,
		typeLabel:     typeLabel,
		adminProvider: adminProvider,
	}
}
