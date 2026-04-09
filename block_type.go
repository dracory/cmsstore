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
const (
	// BLOCK_ORIGIN_SYSTEM for built-in block types
	BLOCK_ORIGIN_SYSTEM = "system"

	// BLOCK_ORIGIN_CUSTOM for user-defined block types
	BLOCK_ORIGIN_CUSTOM = "custom"
)

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
	// Options can include runtime attributes via WithAttributes(attrs).
	// Example:
	//   blockType.Render(ctx, block) // Without attributes
	//   blockType.Render(ctx, block, WithAttributes(map[string]string{"depth": "2"})) // With attributes
	//
	// Parameters:
	//   - ctx: Request context
	//   - block: The block to render
	//   - opts: Optional render options (variadic)
	//
	// Returns:
	//   - HTML string to display on the page
	//   - Error if rendering fails
	Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error)

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

	// GetCustomVariables returns metadata about custom variables this block type
	// can set during rendering via VarsFromContext. Returns nil or empty slice if none.
	//
	// Variables are always strings and are referenced in content as [[name]].
	//
	// Example:
	//   func (b *BlogBlockType) GetCustomVariables() []BlockCustomVariable {
	//       return []BlockCustomVariable{
	//           {Name: "blog_title",  Description: "The blog post title"},
	//           {Name: "blog_author", Description: "The post author name"},
	//       }
	//   }
	GetCustomVariables() []BlockCustomVariable
}

// BlockCustomVariable describes a custom variable that a block type can set during rendering.
// Variables are always strings and are referenced in content as [[name]].
type BlockCustomVariable struct {
	Name        string // Variable name, e.g. "blog_title"
	Description string // Human-readable description of what the variable contains
}

// BlockTypeRegistry manages all registered block types.
//
// This is a global registry that stores block type definitions. Both the frontend
// and admin systems use this registry to look up block types.
//
// The registry is thread-safe and can be accessed concurrently.
type BlockTypeRegistry struct {
	types   map[string]BlockType
	origins map[string]string // typeKey -> origin (system/custom)
	mu      sync.RWMutex
}

var globalBlockTypeRegistry = &BlockTypeRegistry{
	types:   make(map[string]BlockType),
	origins: make(map[string]string),
}

// RegisterSystemBlockType registers a system (built-in) block type.
//
// This is used internally by the CMS for built-in types like HTML, Menu, Navbar.
// For user-defined custom block types, use RegisterCustomBlockType instead.
func RegisterSystemBlockType(blockType BlockType) {
	globalBlockTypeRegistry.RegisterWithOrigin(blockType, BLOCK_ORIGIN_SYSTEM)
}

// RegisterCustomBlockType registers a user-defined (custom) block type.
//
// Use this for your own custom block types. The CMS will treat these as
// user-defined and display them with a different badge color (cyan).
func RegisterCustomBlockType(blockType BlockType) {
	globalBlockTypeRegistry.RegisterWithOrigin(blockType, BLOCK_ORIGIN_CUSTOM)
}

// GetBlockTypeOrigin returns the origin for a block type (system/custom).
func GetBlockTypeOrigin(typeKey string) string {
	return globalBlockTypeRegistry.GetOrigin(typeKey)
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

// RegisterWithOrigin registers a block type with a specified origin.
func (r *BlockTypeRegistry) RegisterWithOrigin(blockType BlockType, origin string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.types[blockType.TypeKey()] = blockType
	r.origins[blockType.TypeKey()] = origin
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

// GetOrigin returns the origin for a block type (system/custom).
// Returns BLOCK_ORIGIN_CUSTOM as default if not explicitly set.
func (r *BlockTypeRegistry) GetOrigin(typeKey string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if origin, ok := r.origins[typeKey]; ok {
		return origin
	}
	return BLOCK_ORIGIN_CUSTOM
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

// GetByOrigin returns all block types with the specified origin.
//
// Use BLOCK_ORIGIN_SYSTEM for built-in types, BLOCK_ORIGIN_CUSTOM for user-defined.
func (r *BlockTypeRegistry) GetByOrigin(origin string) map[string]BlockType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make(map[string]BlockType)
	for typeKey, typeOrigin := range r.origins {
		if typeOrigin == origin {
			if bt, ok := r.types[typeKey]; ok {
				result[typeKey] = bt
			}
		}
	}
	return result
}

// GetSystemBlockTypes returns all system (built-in) block types.
func GetSystemBlockTypes() map[string]BlockType {
	return globalBlockTypeRegistry.GetByOrigin(BLOCK_ORIGIN_SYSTEM)
}

// GetCustomBlockTypes returns all custom (user-defined) block types.
func GetCustomBlockTypes() map[string]BlockType {
	return globalBlockTypeRegistry.GetByOrigin(BLOCK_ORIGIN_CUSTOM)
}

// RenderOption configures block rendering behavior.
type RenderOption func(*RenderOptions)

// RenderOptions holds rendering configuration.
type RenderOptions struct {
	// Attributes contains runtime attributes passed via block reference syntax.
	// Example: <block id="menu" depth="2" /> results in {"depth": "2"}
	Attributes map[string]string
}

// WithAttributes passes runtime attributes to the block renderer.
// Used when blocks are referenced with attribute syntax:
//
//	<block id="menu_main" depth="2" style="sidebar" />
//	[[block id='menu_main' depth='2' style='sidebar']]
func WithAttributes(attrs map[string]string) RenderOption {
	return func(opts *RenderOptions) {
		opts.Attributes = attrs
	}
}

// BlockAttributeDefinition describes a single runtime attribute.
type BlockAttributeDefinition struct {
	// Name is the attribute name (e.g., "depth", "style")
	Name string

	// Type is the attribute type: "string", "int", "bool", "enum", "float"
	Type string

	// Required indicates whether the attribute is required
	Required bool

	// Default is the default value if not provided
	Default interface{}

	// Description is a human-readable description
	Description string

	// EnumValues lists valid values for enum type
	EnumValues []string

	// Validation is a validation rule: "range:1,10", "regex:^[a-z]+$", etc.
	Validation string

	// MinValue is the minimum value (for int/float types)
	MinValue *float64

	// MaxValue is the maximum value (for int/float types)
	MaxValue *float64
}
