# Built-in Block Types

This directory contains the built-in block types that ship with cmsstore. Each block type is a complete, unified implementation that includes both frontend rendering and admin UI.

## Structure

```
blocks/
├── html/
│   └── html_block_type.go    # HTML block (raw HTML content)
├── menu/
│   └── menu_block_type.go    # Menu block (navigation menus)
└── README.md                  # This file
```

## Built-in Block Types

### HTML Block (`html`)
- **Type Key**: `cmsstore.BLOCK_TYPE_HTML` ("html")
- **Purpose**: Renders raw HTML content
- **Admin UI**: CodeMirror editor with syntax highlighting
- **Use Case**: Custom HTML snippets, embedded content, scripts

### Menu Block (`menu`)
- **Type Key**: `cmsstore.BLOCK_TYPE_MENU` ("menu")
- **Purpose**: Renders navigation menus
- **Admin UI**: Menu selector, style options, depth controls
- **Styles**: Vertical, Horizontal, Dropdown, Breadcrumb
- **Use Case**: Site navigation, footer menus, sidebar menus

## Architecture

Each built-in block type follows the unified `BlockType` interface:

```go
type BlockType interface {
    TypeKey() string
    TypeLabel() string
    Render(ctx context.Context, block BlockInterface) (string, error)
    GetAdminFields(block BlockInterface, r *http.Request) interface{}
    SaveAdminFields(r *http.Request, block BlockInterface) error
}
```

This ensures:
- ✅ Frontend rendering and admin UI stay in sync
- ✅ Single source of truth for each block type
- ✅ Consistent pattern for all block types
- ✅ Easy to understand and maintain

## Registration

Built-in block types are automatically registered during initialization:

```go
// In admin/blocks/UI.go
func initBlockAdminProviders(store cmsstore.StoreInterface, logger *slog.Logger) *BlockAdminFieldProviderRegistry {
    registry := NewBlockAdminFieldProviderRegistry()
    
    // Register built-in block types globally
    cmsstore.RegisterBlockType(html.NewHTMLBlockType())
    cmsstore.RegisterBlockType(menu.NewMenuBlockType(store, logger))
    
    return registry
}
```

## Adding Custom Block Types

Custom block types should **not** be added to this directory. Instead, define them in your own project and register them globally:

```go
// In your project
type GalleryBlockType struct { /* ... */ }

func main() {
    // Register your custom block type
    cmsstore.RegisterBlockType(&GalleryBlockType{})
    
    // Now it works everywhere!
}
```

See `docs/UNIFIED_BLOCK_TYPES.md` for complete examples and documentation.

## Design Principles

1. **Self-contained**: Each block type folder contains everything needed
2. **Unified**: Frontend and admin in one struct
3. **Consistent**: All follow the same `BlockType` interface
4. **Documented**: Clear purpose and usage for each type
5. **Testable**: Easy to unit test in isolation

## Future Built-in Types

Potential candidates for future built-in block types:
- **Image Block**: Responsive images with captions
- **Video Block**: Embedded video players
- **Form Block**: Contact forms, newsletters
- **Code Block**: Syntax-highlighted code snippets
- **Accordion Block**: Collapsible content sections

These would follow the same unified pattern and live in their own folders.
