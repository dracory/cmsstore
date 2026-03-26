# Unified Block Type System

## Overview

The **unified block type system** allows you to define custom block types in **one place**, ensuring frontend rendering and admin UI stay perfectly synchronized.

## Project Structure

Built-in block types are organized in the `blocks/` folder at the cmsstore package level:

```
cmsstore/
├── blocks/
│   ├── html/
│   │   └── html_block_type.go    # HTML block implementation
│   ├── menu/
│   │   └── menu_block_type.go    # Menu block implementation
│   └── README.md                  # Built-in blocks documentation
├── block_type.go                  # BlockType interface & registry
├── block_type_adapters.go         # Backward compatibility adapters
└── block_type_example_test.go    # Complete examples
```

This follows the same organizational pattern as `frontend/blocks/` and `admin/blocks/`, making the codebase consistent and easy to navigate.

### The Problem (Old Way)

Previously, you had to register two separate things:

```go
// Frontend renderer
frontend.BlockRegistry().Register("gallery", galleryRenderer)

// Admin provider (easy to forget!)
adminUI.BlockAdminRegistry().Register("gallery", galleryAdminProvider)
```

**Issues:**
- ❌ Easy to forget one registration
- ❌ Renderer and admin UI can get out of sync
- ❌ Two separate files to maintain
- ❌ No compile-time guarantee they match

### The Solution (New Way)

Define everything in **one struct** and register **once**:

```go
type GalleryBlockType struct {
    store cmsstore.StoreInterface
}

// Single registration - both frontend and admin work!
cmsstore.RegisterBlockType(&GalleryBlockType{store: store})
```

**Benefits:**
- ✅ Single registration point
- ✅ Frontend and admin always in sync
- ✅ One file to maintain
- ✅ Compile-time type safety
- ✅ Impossible to forget admin UI

## Quick Start

### 1. Define Your Block Type

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    
    "github.com/dracory/cmsstore"
    "github.com/dracory/form"
    "github.com/dracory/req"
)

type GalleryBlockType struct {
    store cmsstore.StoreInterface
}

// TypeKey: Unique identifier (stored in database)
func (t *GalleryBlockType) TypeKey() string {
    return "gallery"
}

// TypeLabel: Display name (shown in admin UI)
func (t *GalleryBlockType) TypeLabel() string {
    return "Gallery Block"
}

// Render: Frontend rendering
func (t *GalleryBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    var images []Image
    json.Unmarshal([]byte(block.Content()), &images)
    
    layout := block.Meta("layout")
    html := fmt.Sprintf(`<div class="gallery gallery-%s">`, layout)
    for _, img := range images {
        html += fmt.Sprintf(`<img src="%s" alt="%s">`, img.URL, img.Alt)
    }
    html += `</div>`
    
    return html, nil
}

// GetAdminFields: Admin form fields
func (t *GalleryBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
    return []form.FieldInterface{
        form.NewField(form.FieldOptions{
            Label: "Gallery Images (JSON)",
            Name:  "gallery_images",
            Type:  form.FORM_FIELD_TYPE_TEXTAREA,
            Value: block.Content(),
        }),
        form.NewField(form.FieldOptions{
            Label: "Layout",
            Name:  "gallery_layout",
            Type:  form.FORM_FIELD_TYPE_SELECT,
            Value: block.Meta("layout"),
            Options: []form.FieldOption{
                {Value: "Grid", Key: "grid"},
                {Value: "Masonry", Key: "masonry"},
            },
        }),
    }
}

// SaveAdminFields: Form submission handling
func (t *GalleryBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
    images := req.GetStringTrimmed(r, "gallery_images")
    layout := req.GetStringTrimmed(r, "gallery_layout")
    
    // Validation
    if images == "" {
        return fmt.Errorf("images are required")
    }
    
    // Save
    block.SetContent(images)
    block.SetMeta("layout", layout)
    return nil
}

type Image struct {
    URL string `json:"url"`
    Alt string `json:"alt"`
}
```

### 2. Register Once

```go
func main() {
    store := cmsstore.NewStore(...)
    
    // Single registration - that's it!
    cmsstore.RegisterBlockType(&GalleryBlockType{store: store})
    
    // Now "gallery" blocks work everywhere:
    // ✅ Frontend rendering
    // ✅ Admin UI dropdown
    // ✅ Admin edit forms
    // ✅ Admin save logic
}
```

## Complete Examples

### Example 1: Video Block

```go
type VideoBlockType struct{}

func (t *VideoBlockType) TypeKey() string { return "video" }
func (t *VideoBlockType) TypeLabel() string { return "Video Block" }

func (t *VideoBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    videoURL := block.Meta("video_url")
    autoplay := block.Meta("autoplay") == "true"
    
    return fmt.Sprintf(`
<video src="%s" controls %s>
    Your browser does not support video.
</video>
`, videoURL, map[bool]string{true: "autoplay", false: ""}[autoplay]), nil
}

func (t *VideoBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
    return []form.FieldInterface{
        form.NewField(form.FieldOptions{
            Label: "Video URL",
            Name:  "video_url",
            Type:  form.FORM_FIELD_TYPE_STRING,
            Value: block.Meta("video_url"),
        }),
        form.NewField(form.FieldOptions{
            Label: "Autoplay",
            Name:  "video_autoplay",
            Type:  form.FORM_FIELD_TYPE_CHECKBOX,
            Value: block.Meta("autoplay"),
        }),
    }
}

func (t *VideoBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
    block.SetMeta("video_url", req.GetStringTrimmed(r, "video_url"))
    block.SetMeta("autoplay", req.GetStringTrimmed(r, "video_autoplay"))
    return nil
}
```

### Example 2: Interactive Vue.js Tree

```go
type VueTreeBlockType struct{}

func (t *VueTreeBlockType) TypeKey() string { return "vue_tree" }
func (t *VueTreeBlockType) TypeLabel() string { return "Interactive Tree (Vue.js)" }

func (t *VueTreeBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    treeData := block.Content()
    blockID := block.ID()
    
    return fmt.Sprintf(`
<div id="vue-tree-%s"></div>
<script type="module">
import { createApp } from 'https://unpkg.com/vue@3/dist/vue.esm-browser.js'

createApp({
  data() { return { treeData: %s, expanded: {} } },
  methods: {
    toggleNode(id) { this.expanded[id] = !this.expanded[id] }
  }
}).mount('#vue-tree-%s')
</script>
`, blockID, treeData, blockID), nil
}

func (t *VueTreeBlockType) GetAdminFields(block cmsstore.BlockInterface, r *http.Request) interface{} {
    return []form.FieldInterface{
        form.NewField(form.FieldOptions{
            Label: "Tree Data (JSON)",
            Name:  "tree_data",
            Type:  form.FORM_FIELD_TYPE_TEXTAREA,
            Value: block.Content(),
        }),
    }
}

func (t *VueTreeBlockType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
    treeData := req.GetStringTrimmed(r, "tree_data")
    
    // Validate JSON
    var test interface{}
    if err := json.Unmarshal([]byte(treeData), &test); err != nil {
        return fmt.Errorf("invalid JSON: %v", err)
    }
    
    block.SetContent(treeData)
    return nil
}
```

## BlockType Interface

```go
type BlockType interface {
    // TypeKey: Unique identifier (e.g., "gallery", "video")
    TypeKey() string
    
    // TypeLabel: Display name (e.g., "Gallery Block")
    TypeLabel() string
    
    // Render: Frontend rendering logic
    Render(ctx context.Context, block BlockInterface) (string, error)
    
    // GetAdminFields: Admin form fields
    GetAdminFields(block BlockInterface, r *http.Request) interface{}
    
    // SaveAdminFields: Form submission handling
    SaveAdminFields(r *http.Request, block BlockInterface) error
}
```

## Registration

```go
// Global registration (recommended)
cmsstore.RegisterBlockType(blockType)

// Retrieve a registered type
blockType := cmsstore.GetBlockType("gallery")

// Get all registered types
allTypes := cmsstore.GetAllBlockTypes()
```

## Backward Compatibility

The old separate registration still works! The system checks:

1. **Global BlockType registry** (new unified way) ← checked first
2. **Local BlockRenderer registry** (old frontend way)
3. **Local BlockAdminFieldProvider registry** (old admin way)

### Adapters for Gradual Migration

If you have existing separate renderers/providers:

```go
// Wrap existing components into a BlockType
renderer := &ExistingRenderer{}
adminProvider := &ExistingAdminProvider{}

blockType := cmsstore.NewBlockTypeAdapter(
    "custom",
    "Custom Block",
    renderer,
    adminProvider,
)

cmsstore.RegisterBlockType(blockType)
```

## Best Practices

### 1. Keep Everything Together

```go
// ✅ Good: Everything in one file
type GalleryBlockType struct {
    store cmsstore.StoreInterface
}
// ... all methods in same file
```

```go
// ❌ Bad: Split across files
type GalleryRenderer struct{} // in renderer.go
type GalleryAdminProvider struct{} // in admin.go
```

### 2. Use Metadata for Configuration

```go
// ✅ Good: Configuration in metadata
block.SetContent(actualImages)
block.SetMeta("layout", "grid")
block.SetMeta("columns", "3")
```

### 3. Validate in SaveAdminFields

```go
func (t *CustomType) SaveAdminFields(r *http.Request, block cmsstore.BlockInterface) error {
    value := req.GetStringTrimmed(r, "field")
    
    if value == "" {
        return fmt.Errorf("field is required")
    }
    
    // Validate format, etc.
    
    block.SetContent(value)
    return nil
}
```

### 4. Handle Empty States Gracefully

```go
func (t *CustomType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
    if block.Content() == "" {
        return "<!-- Empty block -->", nil
    }
    // ... render logic
}
```

## Migration Guide

### From Separate Registrations

**Before:**
```go
frontend.BlockRegistry().Register("gallery", galleryRenderer)
adminUI.BlockAdminRegistry().Register("gallery", galleryAdminProvider)
```

**After:**
```go
cmsstore.RegisterBlockType(&GalleryBlockType{store: store})
```

### Steps

1. Create a new struct implementing `BlockType`
2. Move rendering logic to `Render()` method
3. Move admin fields to `GetAdminFields()` method
4. Move save logic to `SaveAdminFields()` method
5. Replace two registrations with one `RegisterBlockType()`
6. Delete old renderer and admin provider files

## Testing

```go
func TestGalleryBlockType(t *testing.T) {
    blockType := &GalleryBlockType{}
    block := cmsstore.NewBlock()
    block.SetContent(`[{"url":"img.jpg"}]`)
    block.SetMeta("layout", "grid")
    
    // Test rendering
    html, err := blockType.Render(context.Background(), block)
    assert.NoError(t, err)
    assert.Contains(t, html, "gallery-grid")
    
    // Test admin fields
    fields := blockType.GetAdminFields(block, nil)
    assert.NotEmpty(t, fields)
}
```

## Summary

✅ **One struct** - Everything in one place  
✅ **One registration** - `RegisterBlockType()` once  
✅ **Always in sync** - Frontend and admin can't diverge  
✅ **Type safe** - Compiler catches errors  
✅ **Easy to maintain** - Single source of truth  
✅ **Backward compatible** - Old way still works  

**Recommended:** Always use the unified `BlockType` system for new block types.
