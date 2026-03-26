# Block Renderers

This directory contains individual block renderers organized by block type. Each block type has its own folder with its renderer implementation and any associated assets.

## Structure

```
blocks/
├── html/
│   ├── renderer.go          # HTMLRenderer implementation
│   └── assets/              # CSS, JS, images for HTML blocks (if needed)
├── menu/
│   ├── renderer.go          # Menu BlockRenderer implementation
│   ├── menu_renderer.go     # MenuRenderer (comprehensive menu rendering)
│   └── assets/              # CSS, JS, images for menu blocks (if needed)
└── [block-type]/
    ├── renderer.go          # Block renderer implementation
    └── assets/              # Block-specific assets
```

## Adding Custom Block Types (External Packages)

**Projects that import this package can register their own custom block types** without modifying the cmsstore package. This is the recommended approach for extending functionality.

### Quick Example

```go
package main

import (
    "context"
    "github.com/dracory/cmsstore"
)

// 1. Define your custom renderer
type GalleryRenderer struct {
    store cmsstore.StoreInterface
}

func (r *GalleryRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    // Your custom rendering logic
    images := parseImages(block.Content())
    layout := block.Meta("layout")
    return renderGalleryHTML(images, layout), nil
}

// 2. Register it after creating the frontend
func main() {
    store := cmsstore.NewStore(...)
    frontend := cmsstore.NewFrontend(store, ...)
    
    // Register your custom block type
    frontend.BlockRegistry().Register("gallery", &GalleryRenderer{store: store})
    
    // Now blocks with Type() == "gallery" will use your renderer
}
```

### Complete Example with Multiple Custom Block Types

```go
package main

import (
    "context"
    "fmt"
    "github.com/dracory/cmsstore"
)

// Video block renderer
type VideoRenderer struct{}

func (r *VideoRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    videoURL := block.Meta("video_url")
    autoplay := block.Meta("autoplay") == "true"
    
    html := fmt.Sprintf(`
        <video src="%s" controls %s>
            Your browser does not support the video tag.
        </video>
    `, videoURL, map[bool]string{true: "autoplay", false: ""}[autoplay])
    
    return html, nil
}

// Carousel block renderer
type CarouselRenderer struct {
    store cmsstore.StoreInterface
}

func (r *CarouselRenderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    items := parseCarouselItems(block.Content())
    interval := block.Meta("interval")
    
    return buildCarouselHTML(items, interval), nil
}

func main() {
    store := cmsstore.NewStore(...)
    frontend := cmsstore.NewFrontend(store, ...)
    
    // Register multiple custom block types
    frontend.BlockRegistry().Register("video", &VideoRenderer{})
    frontend.BlockRegistry().Register("carousel", &CarouselRenderer{store: store})
    
    // Start your application
    http.HandleFunc("/", frontend.Handler)
    http.ListenAndServe(":8080", nil)
}
```

### Best Practices for Custom Renderers

1. **Implement the BlockRenderer interface**
   ```go
   type BlockRenderer interface {
       Render(ctx context.Context, block cmsstore.BlockInterface) (string, error)
   }
   ```

2. **Use block metadata for configuration**
   ```go
   layout := block.Meta("layout")
   cssClass := block.Meta("css_class")
   ```

3. **Handle errors gracefully**
   ```go
   if block.Content() == "" {
       return "<!-- Empty block -->", nil
   }
   ```

4. **Access store if needed**
   ```go
   type CustomRenderer struct {
       store cmsstore.StoreInterface
   }
   ```

5. **Return HTML comments for debugging**
   ```go
   return "<!-- Custom block rendered successfully -->", nil
   ```

## Adding Built-in Block Types (Internal)

If you're contributing to the cmsstore package itself:

1. Create a new folder: `blocks/[block-type]/`
2. Create `renderer.go` with your block renderer implementation
3. Implement the `BlockRenderer` interface:
   ```go
   type BlockRenderer interface {
       Render(ctx context.Context, block cmsstore.BlockInterface) (string, error)
   }
   ```
4. Register the renderer in `frontend/block_renderer.go`:
   ```go
   registry.Register(cmsstore.BLOCK_TYPE_[TYPE], [type].New[Type]Renderer(f))
   ```
5. Add any assets to the `assets/` subfolder

## Interface Requirements

Each renderer should depend on the `FrontendStore` interface rather than the concrete `frontend` struct to maintain loose coupling:

```go
type FrontendStore interface {
    MenuFindByID(ctx context.Context, id string) (cmsstore.MenuInterface, error)
    MenuItemList(ctx context.Context, query cmsstore.MenuItemQueryInterface) ([]cmsstore.MenuItemInterface, error)
    MenusEnabled() bool
    PageFindByID(ctx context.Context, id string) (cmsstore.PageInterface, error)
    Logger() *slog.Logger
}
```

## Architecture Pattern

All renderers follow a consistent pattern:

1. **Renderer Struct**: Each renderer has its own struct (e.g., `HTMLRenderer`, `MenuRenderer`)
2. **Constructor Function**: `New[Type]Renderer()` creates renderer instances
3. **Interface Implementation**: All implement the `BlockRenderer` interface
4. **Delegation**: Main `frontend.go` delegates to specialized renderers for consistency

**Example Pattern:**
```go
// In frontend/blocks/[type]/renderer.go
type [Type]Renderer struct { ... }
func New[Type]Renderer(store FrontendStore) *[Type]Renderer { ... }
func (r *[Type]Renderer) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) { ... }

// In frontend/block_renderer.go
registry.Register(cmsstore.BLOCK_TYPE_[TYPE], [type].New[Type]Renderer(f))
```
