# Block Renderers

This directory contains individual block renderers organized by block type. Each block type has its own folder with its renderer implementation and any associated assets.

## Structure

```
blocks/
├── html/
│   ├── renderer.go          # HTML block renderer implementation
│   └── assets/              # CSS, JS, images for HTML blocks (if needed)
├── menu/
│   ├── renderer.go          # Menu block renderer implementation
│   └── assets/              # CSS, JS, images for menu blocks (if needed)
└── [block-type]/
    ├── renderer.go          # Block renderer implementation
    └── assets/              # Block-specific assets
```

## Adding New Block Types

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
   registry.Register(cmsstore.BLOCK_TYPE_[TYPE], [type].NewBlockRenderer(f))
   ```
5. Add any assets to the `assets/` subfolder

## Interface Requirements

Each renderer should depend on the `FrontendStore` interface rather than the concrete `frontend` struct to maintain loose coupling:

```go
type FrontendStore interface {
    MenuFindByID(ctx context.Context, id string) (cmsstore.MenuInterface, error)
    MenuItemList(ctx context.Context, query cmsstore.MenuItemQueryInterface) ([]cmsstore.MenuItemInterface, error)
    MenusEnabled() bool
    RenderMenuHTML(ctx context.Context, menuItems []cmsstore.MenuItemInterface, style, cssClass string, startLevel, maxDepth int) (string, error)
    Logger() *slog.Logger
}
```
