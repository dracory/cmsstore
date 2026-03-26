package admin

import (
	"net/http"
	"sync"

	"github.com/dracory/cmsstore"
)

// BlockAdminFieldProvider defines how custom block types render their admin UI.
//
// Custom block types must implement this interface to provide:
//   - Form fields for the content editing tab
//   - Display label for the block type
//   - Logic to save form data back to the block
//
// Example implementation for a gallery block:
//
//	type GalleryAdminProvider struct {
//	    store cmsstore.StoreInterface
//	}
//
//	func (p *GalleryAdminProvider) GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface {
//	    return []form.FieldInterface{
//	        form.NewField(form.FieldOptions{
//	            Label: "Gallery Images (JSON)",
//	            Name:  "gallery_images",
//	            Type:  form.FORM_FIELD_TYPE_TEXTAREA,
//	            Value: block.Content(),
//	        }),
//	        form.NewField(form.FieldOptions{
//	            Label: "Layout Style",
//	            Name:  "gallery_layout",
//	            Type:  form.FORM_FIELD_TYPE_SELECT,
//	            Value: block.Meta("layout"),
//	            Options: []form.FieldOption{
//	                {Value: "Grid", Key: "grid"},
//	                {Value: "Masonry", Key: "masonry"},
//	            },
//	        }),
//	    }
//	}
//
//	func (p *GalleryAdminProvider) GetTypeLabel() string {
//	    return "Gallery Block"
//	}
//
//	func (p *GalleryAdminProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
//	    block.SetContent(req.GetStringTrimmed(r, "gallery_images"))
//	    block.SetMeta("layout", req.GetStringTrimmed(r, "gallery_layout"))
//	    return nil
//	}
//
// To register a custom admin provider:
//
//	adminUI := admin.UI(config)
//	adminUI.BlockAdminRegistry().Register("gallery", &GalleryAdminProvider{store: store})
//
// See admin/blocks/README.md for detailed examples and best practices.
type BlockAdminFieldProvider interface {
	// GetContentFields returns form fields for the block content editing tab.
	//
	// The fields should allow users to configure all block-specific properties.
	// Use the block parameter to read current values, and the request for context.
	//
	// Parameters:
	//   - block: The block being edited (use for reading current values)
	//   - r: The HTTP request (use for context, user info, etc.)
	//
	// Returns:
	//   - Array of form fields to display in the content tab (should return []form.FieldInterface)
	GetContentFields(block cmsstore.BlockInterface, r *http.Request) interface{}

	// GetTypeLabel returns the human-readable display name for this block type.
	//
	// This label appears in:
	//   - Block type dropdown during creation
	//   - Block type display field in settings
	//
	// Example: "Gallery Block", "Video Block", "Custom Tree Block"
	GetTypeLabel() string

	// SaveContentFields processes form submission and updates the block.
	//
	// This method is called when the user saves the content tab. It should:
	//   1. Read form values from the request
	//   2. Validate the input
	//   3. Update the block's content and metadata
	//   4. Return an error if validation fails
	//
	// Parameters:
	//   - r: The HTTP request containing form data
	//   - block: The block to update (modify in place)
	//
	// Returns:
	//   - error: Validation or processing error, or nil on success
	SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error
}

// BlockAdminFieldProviderRegistry manages all registered block admin field providers.
//
// The registry is thread-safe and maps block types (strings) to their admin providers.
// It's used by the admin controllers to dynamically generate forms for different block types.
//
// Built-in block types:
//   - cmsstore.BLOCK_TYPE_HTML: HTML content editor
//   - cmsstore.BLOCK_TYPE_MENU: Menu configuration form
//
// Custom block types can be registered after admin UI initialization:
//
//	adminUI.BlockAdminRegistry().Register("custom_type", customProvider)
type BlockAdminFieldProviderRegistry struct {
	providers map[string]BlockAdminFieldProvider
	mu        sync.RWMutex
}

// NewBlockAdminFieldProviderRegistry creates a new registry for block admin field providers.
// This is typically called internally during admin UI initialization.
func NewBlockAdminFieldProviderRegistry() *BlockAdminFieldProviderRegistry {
	return &BlockAdminFieldProviderRegistry{
		providers: make(map[string]BlockAdminFieldProvider),
	}
}

// Register registers an admin field provider for a specific block type.
//
// The blockType should match the value stored in the block's Type() field.
// If a provider already exists for the given type, it will be replaced.
//
// This method is thread-safe and can be called after admin UI initialization
// to add custom block types.
//
// Example:
//
//	adminUI.BlockAdminRegistry().Register("video", &VideoAdminProvider{})
//	adminUI.BlockAdminRegistry().Register("carousel", &CarouselAdminProvider{store: store})
func (r *BlockAdminFieldProviderRegistry) Register(blockType string, provider BlockAdminFieldProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[blockType] = provider
}

// GetProvider returns the admin field provider for the given block type.
//
// If no provider is registered for the type, returns nil.
// Callers should check for nil and provide a fallback.
//
// This method is thread-safe.
func (r *BlockAdminFieldProviderRegistry) GetProvider(blockType string) BlockAdminFieldProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.providers[blockType]
}

// GetAllProviders returns a map of all registered block type providers.
//
// This is useful for:
//   - Generating the block type dropdown in create forms
//   - Listing available block types
//
// Returns a copy of the internal map to prevent external modification.
// This method is thread-safe.
func (r *BlockAdminFieldProviderRegistry) GetAllProviders() map[string]BlockAdminFieldProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	copy := make(map[string]BlockAdminFieldProvider, len(r.providers))
	for k, v := range r.providers {
		copy[k] = v
	}
	return copy
}

// GetRegisteredTypes returns a sorted list of all registered block type keys.
//
// This is useful for generating dropdowns and validating block types.
// The list is sorted alphabetically for consistent UI ordering.
func (r *BlockAdminFieldProviderRegistry) GetRegisteredTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.providers))
	for k := range r.providers {
		types = append(types, k)
	}
	return types
}
