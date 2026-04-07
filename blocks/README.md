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


## Custom Variables

Blocks can set custom variables that are automatically replaced in the final rendered content. This allows blocks to expose dynamic data that can be referenced anywhere in page or template content.

### Setting Variables in Blocks

In your block's `Render` method, use `VarsFromContext` to set variables:

```go
func (b *BlogBlockType) Render(ctx context.Context, block cmsstore.BlockInterface, opts ...cmsstore.RenderOption) (string, error) {
    // Fetch your data
    post := fetchBlogPost(block.Data())
    
    // Set custom variables
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        vars.Set("blog_title", post.Title)
        vars.Set("blog_author", post.Author)
        vars.Set("blog_date", post.Date.Format("2006-01-02"))
    }
    
    return renderHTML(post), nil
}
```

### Using Variables in Content

Reference variables using `[[variable_name]]` syntax in your page or template content:

```html
<!DOCTYPE html>
<html>
<head>
    <title>[[PageTitle]] - [[blog_title]]</title>
    <meta property="og:title" content="[[blog_title]]">
    <meta property="article:author" content="[[blog_author]]">
</head>
<body>
    <h1>[[blog_title]]</h1>
    <p>By [[blog_author]] on [[blog_date]]</p>
    
    [[BLOCK_blog_content]]
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

### Multiple Blocks

Multiple blocks can set different variables:

```go
// Blog block
vars.Set("blog_title", "My Post")

// Product block
vars.Set("product_price", "$99")

// User block
vars.Set("user_name", "Jane")
```

All variables are available in the final content.

### Variable Collisions

If multiple blocks set the same variable name, the last block (in content order) wins. To avoid collisions:

1. Use unique variable names
2. Use namespacing (e.g., `blog:title`, `product:title`)
3. Use prefixes (e.g., `blog_title`, `product_title`)

### Example Use Cases

**Blog Post SEO:**
```go
vars.Set("blog_title", post.Title)
vars.Set("blog_excerpt", post.Excerpt)
vars.Set("blog_image", post.FeaturedImage)
```

**E-commerce Product:**
```go
vars.Set("product_name", product.Name)
vars.Set("product_price", product.Price)
vars.Set("product_sku", product.SKU)
```

**Event Information:**
```go
vars.Set("event_name", event.Name)
vars.Set("event_date", event.Date.Format("2006-01-02"))
vars.Set("event_location", event.Location)
```

See `examples/custom-variables/` for complete examples.
