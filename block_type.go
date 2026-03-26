package cmsstore

import (
	"context"
	"net/http"
	"sync"
)

// BlockType defines a complete block type with both frontend rendering and admin UI.
//
// This is the recommended way to define custom block types. It ensures that both
// the frontend renderer and admin UI are registered together, preventing the common
// mistake of registering one but forgetting the other.
//
// Example implementation:
//
//	type GalleryBlockType struct {
//	    store StoreInterface
//	}
//
//	func (t *GalleryBlockType) TypeKey() string {
//	    return "gallery"
//	}
//
//	func (t *GalleryBlockType) TypeLabel() string {
//	    return "Gallery Block"
//	}
//
//	func (t *GalleryBlockType) Render(ctx context.Context, block BlockInterface) (string, error) {
//	    // Frontend rendering logic
//	    images := parseImages(block.Content())
//	    return renderGalleryHTML(images, block.Meta("layout")), nil
//	}
//
//	func (t *GalleryBlockType) GetAdminFields(block BlockInterface, r *http.Request) []form.FieldInterface {
//	    // Admin form fields
//	    return []form.FieldInterface{
//	        form.NewField(form.FieldOptions{
//	            Label: "Gallery Images",
//	            Name:  "gallery_images",
//	            Type:  form.FORM_FIELD_TYPE_TEXTAREA,
//	            Value: block.Content(),
//	        }),
//	    }
//	}
//
//	func (t *GalleryBlockType) SaveAdminFields(r *http.Request, block BlockInterface) error {
//	    // Save form data
//	    images := req.GetStringTrimmed(r, "gallery_images")
//	    block.SetContent(images)
//	    return nil
//	}
//
// To register:
//
//	cmsstore.RegisterBlockType(&GalleryBlockType{store: store})
//
// This single registration automatically makes the block type available in both
// frontend rendering and admin UI.
type BlockType interface {
	// TypeKey returns the unique identifier for this block type.
	// This is stored in the database and used for lookups.
	// Example: "gallery", "video", "custom_tree"
	TypeKey() string

	// TypeLabel returns the human-readable display name.
	// This appears in the admin UI block type dropdown.
	// Example: "Gallery Block", "Video Block", "Custom Tree Block"
	TypeLabel() string

	// Render renders the block for frontend display.
	// This is called when the block appears on a page.
	//
	// Parameters:
	//   - ctx: Request context
	//   - block: The block to render
	//
	// Returns:
	//   - HTML string to display on the page
	//   - Error if rendering fails
	Render(ctx context.Context, block BlockInterface) (string, error)

	// GetAdminFields returns form fields for the admin content editing tab.
	//
	// This method is called when displaying the block edit form in the admin panel.
	// Return an array of form fields that allow users to configure the block.
	//
	// Parameters:
	//   - block: The block being edited (use for reading current values)
	//   - r: The HTTP request (use for context, loading related data, etc.)
	//
	// Returns:
	//   - Array of form fields to display
	GetAdminFields(block BlockInterface, r *http.Request) interface{}

	// SaveAdminFields processes form submission and updates the block.
	//
	// This method is called when the user saves the content tab in the admin panel.
	// Read form values, validate them, and update the block accordingly.
	//
	// Parameters:
	//   - r: The HTTP request containing form data
	//   - block: The block to update (modify in place)
	//
	// Returns:
	//   - Error if validation fails, or nil on success
	SaveAdminFields(r *http.Request, block BlockInterface) error
}

// BlockTypeRegistry manages all registered block types.
//
// This is a global registry that stores block type definitions. Both the frontend
// and admin systems use this registry to look up block types.
//
// The registry is thread-safe and can be accessed concurrently.
type BlockTypeRegistry struct {
	types map[string]BlockType
	mu    sync.RWMutex
}

var globalBlockTypeRegistry = &BlockTypeRegistry{
	types: make(map[string]BlockType),
}

// RegisterBlockType registers a block type globally.
//
// This is the recommended way to add custom block types. A single registration
// makes the block type available in both frontend rendering and admin UI.
//
// Example:
//
//	cmsstore.RegisterBlockType(&GalleryBlockType{store: store})
//
// The block type will automatically be available:
//   - In the admin UI block type dropdown
//   - For frontend rendering when blocks of this type are displayed
//
// This method is thread-safe.
func RegisterBlockType(blockType BlockType) {
	globalBlockTypeRegistry.Register(blockType)
}

// GetBlockType retrieves a registered block type by its key.
//
// Returns nil if no block type is registered with the given key.
// This method is thread-safe.
func GetBlockType(typeKey string) BlockType {
	return globalBlockTypeRegistry.Get(typeKey)
}

// GetAllBlockTypes returns all registered block types.
//
// Returns a map of type keys to block types.
// The returned map is a copy to prevent external modification.
// This method is thread-safe.
func GetAllBlockTypes() map[string]BlockType {
	return globalBlockTypeRegistry.GetAll()
}

// Register registers a block type in the registry.
func (r *BlockTypeRegistry) Register(blockType BlockType) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.types[blockType.TypeKey()] = blockType
}

// Get retrieves a block type by its key.
func (r *BlockTypeRegistry) Get(typeKey string) BlockType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.types[typeKey]
}

// GetAll returns all registered block types.
func (r *BlockTypeRegistry) GetAll() map[string]BlockType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy to prevent external modification
	copy := make(map[string]BlockType, len(r.types))
	for k, v := range r.types {
		copy[k] = v
	}
	return copy
}

// GetKeys returns all registered block type keys.
func (r *BlockTypeRegistry) GetKeys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]string, 0, len(r.types))
	for k := range r.types {
		keys = append(keys, k)
	}
	return keys
}
