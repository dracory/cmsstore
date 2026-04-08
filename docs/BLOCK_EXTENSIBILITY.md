# Block System Extensibility Guide

This document provides a comprehensive overview of the extensible block system in cmsstore, covering both frontend rendering and admin UI customization.

## Overview

The cmsstore block system is fully extensible, allowing external packages to:

1. **Register custom block types** for frontend rendering
2. **Provide admin UI** for managing custom block types
3. **Use any frontend framework** (Vue.js, React, Alpine.js, Mithril, etc.)

## Architecture

### Unified Block Types (Recommended)

- **Interface**: `BlockType`
- **Registry**: Global `BlockTypeRegistry`
- **Location**: `block_type.go`
- **Built-in Types**: `blocks/html/`, `blocks/menu/`
- **Registration**: `cmsstore.RegisterBlockType()`

### Legacy Separate Systems (Backward Compatible)

#### Frontend Rendering

- **Interface**: `BlockRenderer`
- **Registry**: `BlockRendererRegistry`
- **Location**: `frontend/block_renderer.go`
- **Access**: `frontend.BlockRegistry()`

#### Admin UI

- **Interface**: `BlockAdminFieldProvider`
- **Registry**: `BlockAdminFieldProviderRegistry`
- **Location**: `admin/blocks/admin_field_provider.go`
- **Access**: `adminUI.BlockAdminRegistry()`

## Quick Start Example

Here's a complete example of adding a custom "Gallery" block type:

### 1. Frontend Renderer

```go
package main

import (
    "context"
    "fmt"
    "github.com/dracory/cmsstore"
)

// GalleryRenderer renders gallery blocks
type GalleryRenderer struct {
    store cmsstore.StoreInterface
}

func (r *GalleryRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    images := parseImages(block.Content())
    layout := block.Meta("layout")
    columns := block.Meta("columns")
    
    html := fmt.Sprintf(`
<div class="gallery gallery-layout-%s" data-columns="%s">
    %s
</div>
<script>
    // Initialize gallery with your preferred library
    initGallery('.gallery');
</script>
`, layout, columns, renderImages(images))
    
    return html, nil
}
```

### 2. Admin Provider

```go
// GalleryAdminProvider provides admin UI for gallery blocks
type GalleryAdminProvider struct {
    store cmsstore.StoreInterface
}

func (p *GalleryAdminProvider) GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface {
    return []form.FieldInterface{
        form.NewField(form.FieldOptions{
            Label: "Gallery Images (JSON)",
            Name:  "gallery_images",
            Type:  form.FORM_FIELD_TYPE_TEXTAREA,
            Value: block.Content(),
        }),
        form.NewField(form.FieldOptions{
            Label: "Layout Style",
            Name:  "gallery_layout",
            Type:  form.FORM_FIELD_TYPE_SELECT,
            Value: block.Meta("layout"),
            Options: []form.FieldOption{
                {Value: "Grid", Key: "grid"},
                {Value: "Masonry", Key: "masonry"},
            },
        }),
        form.NewField(form.FieldOptions{
            Label: "Columns",
            Name:  "gallery_columns",
            Type:  form.FORM_FIELD_TYPE_NUMBER,
            Value: block.Meta("columns"),
        }),
    }
}

func (p *GalleryAdminProvider) GetTypeLabel() string {
    return "Gallery Block"
}

func (p *GalleryAdminProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
    images := req.GetStringTrimmed(r, "gallery_images")
    layout := req.GetStringTrimmed(r, "gallery_layout")
    columns := req.GetStringTrimmed(r, "gallery_columns")
    
    if images == "" {
        return &admin.ValidationError{Message: "Images are required"}
    }
    
    block.SetContent(images)
    block.SetMeta("layout", layout)
    block.SetMeta("columns", columns)
    return nil
}
```

### 3. Registration

```go
func main() {
    store := cmsstore.NewStore(...)
    
    // Frontend
    frontend := cmsstore.NewFrontend(store, ...)
    frontend.BlockRegistry().Register("gallery", &GalleryRenderer{store: store})
    
    // Admin
    adminUI := admin.UI(admin.UiConfig{...})
    adminUI.BlockAdminRegistry().Register("gallery", &GalleryAdminProvider{store: store})
    
    // Now "gallery" blocks work in both frontend and admin!
}
```

## Interactive Blocks with JavaScript Frameworks

### Vue.js Example

```go
type VueTreeRenderer struct{}

func (r *VueTreeRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    treeData := block.Content()
    blockID := block.ID()
    
    return fmt.Sprintf(`
<div id="vue-tree-%s"></div>
<script type="module">
import { createApp } from 'https://unpkg.com/vue@3/dist/vue.esm-browser.js'

createApp({
  data() {
    return { treeData: %s }
  },
  template: '<div class="tree">{{ treeData }}</div>'
}).mount('#vue-tree-%s')
</script>
`, blockID, treeData, blockID), nil
}
```

### Alpine.js Example

```go
type AlpineAccordionRenderer struct{}

func (r *AlpineAccordionRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    items := parseItems(block.Content())
    
    html := `<div x-data="{ active: null }">`
    for i, item := range items {
        html += fmt.Sprintf(`
  <div>
    <button @click="active = active === %d ? null : %d">%s</button>
    <div x-show="active === %d" x-transition>%s</div>
  </div>
`, i, i, item.Title, i, item.Content)
    }
    html += `</div>`
    return html, nil
}
```

### Mithril.js Example

```go
type MithrilMenuRenderer struct {
    store cmsstore.StoreInterface
}

func (r *MithrilMenuRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    menuData := fetchMenuData(ctx, r.store, block)
    menuJSON, _ := json.Marshal(menuData)
    
    return fmt.Sprintf(`
<div id="mithril-menu-%s"></div>
<script src="https://unpkg.com/mithril/mithril.js"></script>
<script>
m.mount(document.getElementById('mithril-menu-%s'), {
  view: () => m('nav', %s.map(item => m('a', {href: item.url}, item.name)))
})
</script>
`, block.ID(), block.ID(), menuJSON), nil
}
```

## Built-in Block Types

### HTML Block
- **Type**: `cmsstore.BLOCK_TYPE_HTML`
- **Frontend**: Renders raw HTML content
- **Admin**: CodeMirror editor with syntax highlighting

### Menu Block
- **Type**: `cmsstore.BLOCK_TYPE_MENU`
- **Frontend**: Renders navigation menus with multiple styles
- **Admin**: Menu selector, style options, CSS class, depth controls

## Documentation References

- **Frontend Rendering**: `frontend/blocks/README.md`
- **Admin UI**: `admin/blocks/README.md`
- **Interface Docs**: `frontend/block_renderer.go`, `admin/blocks/admin_field_provider.go`

## Best Practices

### 1. Use Metadata for Configuration
```go
block.SetContent(actualContent)      // Main content
block.SetMeta("layout", "grid")      // Configuration
block.SetMeta("columns", "3")        // Configuration
```

### 2. Validate Input
```go
func (p *CustomProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
    value := req.GetStringTrimmed(r, "field")
    if value == "" {
        return &admin.ValidationError{Message: "Required"}
    }
    return nil
}
```

### 3. Graceful Error Handling
```go
func (r *CustomRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    if block.Content() == "" {
        return "<!-- Empty block -->", nil
    }
    // ... render logic
}
```

### 4. Thread Safety
Both registries are thread-safe. You can register providers at any time after initialization.

### 5. Fallback Behavior
- **Frontend**: Falls back to `NoOpRenderer` (returns HTML comment)
- **Admin**: Falls back to HTML provider (basic textarea)

### 6. Accessing Request Data in Blocks

All block types can access the HTTP request from the context during frontend rendering:

```go
func (b *myBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    req := cmsstore.RequestFromContext(ctx)
    if req != nil {
        searchQuery := req.URL.Query().Get("q")
        // ... use request data
    }
    // ... render block
}
```

**Important**: Always check if `RequestFromContext()` returns `nil` - the request may not be available in admin preview or CLI contexts.

For complete details, see [Request Context in Blocks](./BLOCK_SYSTEM_ARCHITECTURE.md#request-context-in-blocks).

### 7. Setting Custom Variables

Blocks can expose data as custom variables that can be referenced anywhere in page or template content:

```go
func (b *BlogBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
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

Then use in page/template content:

```html
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
```

**Variable naming**: Use any convention - `snake_case`, `camelCase`, `namespaced:var`, `$prefixed`, etc.

For complete details and examples, see [Custom Variables in Blocks](./BLOCK_SYSTEM_ARCHITECTURE.md#custom-variables-in-blocks).

## Testing Blocks

```go
func TestBlockRenderer(t *testing.T) {
    renderer := &GalleryRenderer{}
    block := cmsstore.NewBlock()
    block.SetContent(`[{"url":"image1.jpg"},{"url":"image2.jpg"}]`)
    block.SetMeta("layout", "grid")
    
    html, err := renderer.Render(context.Background(), block)
    
    assert.NoError(t, err)
    assert.Contains(t, html, "gallery-layout-grid")
}
```

## Migration Guide

If you have existing custom blocks hardcoded in the system:

1. Extract rendering logic into a `BlockRenderer` implementation
2. Extract admin form logic into a `BlockAdminFieldProvider` implementation
3. Register both providers after initialization
4. Remove hardcoded switch statements

## Support

For questions or issues:
- Check `frontend/blocks/README.md` for frontend examples
- Check `admin/blocks/README.md` for admin UI examples
- Review built-in providers: `admin/blocks/admin_provider_html.go`, `admin/blocks/admin_provider_menu.go`
