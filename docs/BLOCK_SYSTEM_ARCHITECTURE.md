# Block System Architecture

## Overview

The cmsstore block system follows a **consistent organizational pattern** across all layers, making it intuitive to navigate and extend.

## Directory Structure

```
cmsstore/
├── blocks/                          # Built-in block types (unified)
│   ├── html/
│   │   └── html_block_type.go      # HTML block (complete implementation)
│   ├── menu/
│   │   └── menu_block_type.go      # Menu block (complete implementation)
│   └── README.md
│
├── frontend/
│   ├── blocks/                      # Frontend-specific renderers (legacy)
│   │   ├── html/
│   │   │   └── renderer.go
│   │   ├── menu/
│   │   │   ├── renderer.go
│   │   │   └── menu_renderer.go
│   │   └── README.md
│   ├── block_renderer.go            # BlockRenderer interface & registry
│   └── frontend.go
│
├── admin/blocks/                    # Admin UI controllers & providers
│   ├── admin_field_provider.go      # BlockAdminFieldProvider interface
│   ├── admin_provider_html.go       # HTML admin provider (legacy)
│   ├── admin_provider_menu.go       # Menu admin provider (legacy)
│   ├── block_create_controller.go
│   ├── block_update_controller.go
│   ├── UI.go                        # Initialization & registration
│   └── README.md
│
├── block_type.go                    # Unified BlockType interface & registry
├── block_type_adapters.go           # Backward compatibility
├── block_type_example_test.go       # Complete examples
│
└── docs/
    ├── UNIFIED_BLOCK_TYPES.md       # Main guide for unified system
    ├── BLOCK_EXTENSIBILITY.md       # Overview of all systems
    └── BLOCK_SYSTEM_ARCHITECTURE.md # This file
```

## Three Layers of Organization

### 1. Core Package Level (`blocks/`)

**Purpose**: Built-in block types using the unified `BlockType` interface

**Pattern**: Each block type in its own folder with complete implementation

**Example**:
```
blocks/
├── html/
│   └── html_block_type.go    # TypeKey, TypeLabel, Render, GetAdminFields, SaveAdminFields
└── menu/
    └── menu_block_type.go    # TypeKey, TypeLabel, Render, GetAdminFields, SaveAdminFields
```

**Registration**: Automatic during admin UI initialization
```go
cmsstore.RegisterBlockType(html.NewHTMLBlockType())
cmsstore.RegisterBlockType(menu.NewMenuBlockType(store, logger))
```

### 2. Frontend Layer (`frontend/blocks/`)

**Purpose**: Frontend-specific rendering logic (legacy separate system)

**Pattern**: Each block type has specialized renderer

**Example**:
```
frontend/blocks/
├── html/
│   └── renderer.go           # HTMLRenderer.Render()
└── menu/
    ├── renderer.go           # BlockRenderer.Render()
    └── menu_renderer.go      # MenuRenderer (comprehensive)
```

**Registration**: Manual via `frontend.BlockRegistry()`
```go
frontend.BlockRegistry().Register("html", html.NewHTMLRenderer())
```

### 3. Admin Layer (`admin/blocks/`)

**Purpose**: Admin UI controllers and field providers

**Pattern**: Controllers + field providers per block type

**Example**:
```
admin/blocks/
├── admin_provider_html.go    # HTMLAdminProvider
├── admin_provider_menu.go    # MenuAdminProvider
├── block_create_controller.go
└── block_update_controller.go
```

**Registration**: Manual via `adminUI.BlockAdminRegistry()`
```go
adminUI.BlockAdminRegistry().Register("html", &HTMLAdminProvider{})
```

## Lookup Priority

When rendering or editing blocks, the system checks in this order:

### Frontend Rendering
1. **Global `BlockType` registry** (`cmsstore.GetBlockType()`) ← Unified system
2. **Local `BlockRenderer` registry** (`frontend.BlockRegistry()`) ← Legacy
3. **Fallback**: `NoOpRenderer` (HTML comment)

During frontend rendering, the `*http.Request` is automatically injected into the context. Blocks can access it via `cmsstore.RequestFromContext(ctx)` to read query parameters, headers, and other request data. See [Request Context in Blocks](#request-context-in-blocks) for details.

### Admin UI
1. **Global `BlockType` registry** (`cmsstore.GetBlockType()`) ← Unified system
2. **Local `BlockAdminFieldProvider` registry** (`adminUI.BlockAdminRegistry()`) ← Legacy
3. **Fallback**: Basic textarea editor

## Request Context in Blocks

All block types (both built-in and custom) can access the HTTP request from the context during frontend rendering. This enables blocks to read query parameters, headers, cookies, and other request-specific data.

### Usage Example

```go
func (b *myBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    // Get the request from context
    req := cmsstore.RequestFromContext(ctx)
    if req == nil {
        // Request not available (e.g., admin preview, CLI rendering)
        return "<!-- No request available -->", nil
    }
    
    // Access query parameters
    searchQuery := req.URL.Query().Get("q")
    pageNum := req.URL.Query().Get("page")
    
    // Access headers
    userAgent := req.Header.Get("User-Agent")
    
    // Access cookies
    cookie, err := req.Cookie("session_id")
    
    // Render with request data
    return fmt.Sprintf("Search: %s (Page: %s)", searchQuery, pageNum), nil
}
```

### Important Notes

- **Always check for nil**: `RequestFromContext()` returns `nil` when the request is not available (admin preview, CLI rendering, background jobs)
- **Automatic injection**: The frontend automatically injects the request before calling `Render()`
- **Read-only**: Blocks should treat the request as read-only; modifications won't affect the actual HTTP response
- **Thread safety**: The request object is not cloned; blocks should not store references to it for async operations

### Use Cases

- **Search blocks**: Read `?q=` query parameter to display search results
- **Filtering blocks**: Read filter parameters from URL
- **Pagination**: Read `?page=` parameter
- **Geolocation**: Use IP from request
- **A/B Testing**: Read cookies or headers

---

## Custom Variables in Blocks

Blocks can set custom variables that are automatically replaced in the final rendered content. This allows blocks to expose dynamic data that can be referenced anywhere in page or template content.

### How It Works

During rendering, blocks can set variables via `VarsFromContext()`. After all blocks are rendered, these variables are replaced in the content using `[[variable_name]]` syntax.

### Setting Variables in Blocks

```go
func (b *BlogBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
    // Fetch your data
    post := fetchBlogPost(block.Data())
    
    // Set custom variables
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        vars.Set("blog_title", post.Title)
        vars.Set("blog_author", post.Author)
        vars.Set("blog_date", post.Date.Format("2006-01-02"))
        vars.Set("blog_excerpt", post.Excerpt)
    }
    
    // Return the block's HTML
    return renderBlogHTML(post), nil
}
```

### Using Variables in Content

Reference variables in your page or template content:

```html
<!DOCTYPE html>
<html>
<head>
    <title>[[PageTitle]] - [[blog_title]]</title>
    <meta name="description" content="[[blog_excerpt]]">
    <meta property="og:title" content="[[blog_title]]">
    <meta property="article:author" content="[[blog_author]]">
    <meta property="article:published_time" content="[[blog_date]]">
</head>
<body>
    <header>
        <h1>[[blog_title]]</h1>
        <p>By [[blog_author]] on [[blog_date]]</p>
    </header>
    
    <main>
        [[BLOCK_blog_content]]
    </main>
</body>
</html>
```

### Variable Naming

You have complete freedom in naming variables. Common patterns:

- **snake_case**: `blog_title`, `user_name`, `product_price`
- **camelCase**: `blogTitle`, `userName`, `productPrice`
- **PascalCase**: `BlogTitle`, `UserName`, `ProductPrice`
- **namespaced**: `blog:title`, `product:price`, `event:date`
- **prefixed**: `$blogTitle`, `@userName`, `#productPrice`

### Important Notes

- **Always check for nil**: `VarsFromContext()` returns `nil` if variables are not available
- **Thread-safe**: Variable storage is protected by mutex
- **Collision handling**: If multiple blocks set the same variable, the last block wins (based on content order)
- **Scope**: Variables are only available during the current render cycle

### Use Cases

**Blog Post SEO:**
```go
vars.Set("post_title", post.Title)
vars.Set("post_excerpt", post.Excerpt)
vars.Set("post_author", post.Author)
vars.Set("post_image", post.FeaturedImage)
```

**E-commerce Product:**
```go
vars.Set("product_name", product.Name)
vars.Set("product_price", product.Price)
vars.Set("product_sku", product.SKU)
vars.Set("product_availability", product.Availability)
```

**Event Information:**
```go
vars.Set("event_name", event.Name)
vars.Set("event_date", event.Date.Format("2006-01-02"))
vars.Set("event_location", event.Location)
vars.Set("event_organizer", event.Organizer)
```

**User Profile:**
```go
vars.Set("user_display_name", user.Name)
vars.Set("user_role", user.Role)
vars.Set("user_join_date", user.CreatedAt.Format("2006-01-02"))
```

### Complete Example

See `examples/custom-variables/` for a complete working example with blog and product blocks.

---

### Current State (Built-in Types)
- ✅ HTML and Menu blocks now use unified `BlockType` in `blocks/` folder
- ✅ Legacy providers in `admin/blocks/admin_provider_*.go` kept for reference
- ✅ Frontend renderers in `frontend/blocks/` still used by menu block
- ✅ All systems work together seamlessly

### For Custom Blocks

**Recommended (New Projects)**:
```go
// Define once in your project
type GalleryBlockType struct { /* ... */ }

// Register once
cmsstore.RegisterBlockType(&GalleryBlockType{})

// Works everywhere!
```

**Legacy (Existing Projects)**:
```go
// Still supported - register separately
frontend.BlockRegistry().Register("gallery", galleryRenderer)
adminUI.BlockAdminRegistry().Register("gallery", galleryAdminProvider)
```

## Design Principles

1. **Consistent Pattern**: Same folder structure at all levels
2. **Single Source of Truth**: Unified types in `blocks/` folder
3. **Backward Compatible**: Legacy systems still work
4. **Progressive Enhancement**: New code uses unified system
5. **Clear Separation**: Core → Frontend → Admin

## Benefits of This Architecture

✅ **Easy to Navigate**: Same pattern everywhere  
✅ **Hard to Break**: Frontend and admin always in sync (unified)  
✅ **Flexible**: Support both old and new registration methods  
✅ **Scalable**: Add new block types by creating new folders  
✅ **Documented**: README in each blocks/ folder  

## Adding a New Built-in Block Type

If you're contributing to cmsstore:

1. Create `blocks/newtype/newtype_block_type.go`
2. Implement `BlockType` interface
3. Register in `admin/blocks/UI.go`:
   ```go
   cmsstore.RegisterBlockType(newtype.NewTypeBlockType())
   ```
4. Update `blocks/README.md`
5. Add examples to `block_type_example_test.go`

## Summary

The block system now follows a **consistent three-layer architecture**:

- **`blocks/`** - Unified built-in types (recommended)
- **`frontend/blocks/`** - Frontend renderers (legacy)
- **`admin/blocks/`** - Admin UI (legacy)

All layers work together, with the unified system taking priority while maintaining full backward compatibility.
