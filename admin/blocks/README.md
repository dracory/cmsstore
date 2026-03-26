# Block Admin UI System

This directory contains the admin UI system for managing CMS blocks. The system is fully extensible, allowing external packages to register custom block types with their own admin interfaces.

## Architecture Overview

The admin system uses a **provider pattern** similar to the frontend block rendering system:

- **BlockAdminFieldProvider**: Interface for defining admin UI for block types
- **BlockAdminFieldProviderRegistry**: Thread-safe registry for managing providers
- **Built-in Providers**: HTML and Menu block admin providers

## Adding Custom Block Admin UI (External Packages)

**Projects that import this package can register their own custom block type admin UI** without modifying the cmsstore package.

### Quick Example

```go
package main

import (
    "net/http"
    "github.com/dracory/cmsstore"
    "github.com/dracory/cmsstore/admin/blocks"
    "github.com/dracory/form"
    "github.com/dracory/req"
)

// 1. Define your custom admin provider
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
            Help:  "Enter image URLs as JSON array",
        }),
        form.NewField(form.FieldOptions{
            Label: "Layout Style",
            Name:  "gallery_layout",
            Type:  form.FORM_FIELD_TYPE_SELECT,
            Value: block.Meta("layout"),
            Options: []form.FieldOption{
                {Value: "Grid", Key: "grid"},
                {Value: "Masonry", Key: "masonry"},
                {Value: "Carousel", Key: "carousel"},
            },
        }),
        form.NewField(form.FieldOptions{
            Label: "Images Per Row",
            Name:  "gallery_columns",
            Type:  form.FORM_FIELD_TYPE_NUMBER,
            Value: block.Meta("columns"),
            Help:  "Number of images per row (for grid layout)",
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
    
    // Validation
    if images == "" {
        return &blocks.ValidationError{Message: "Gallery images are required"}
    }
    
    // Save to block
    block.SetContent(images)
    block.SetMeta("layout", layout)
    block.SetMeta("columns", columns)
    
    return nil
}

// 2. Register it after creating the admin UI
func main() {
    store := cmsstore.NewStore(...)
    adminUI := blocks.UI(blocks.UiConfig{
        Store:  store,
        Logger: logger,
        Layout: layoutFunc,
    })
    
    // Register your custom block admin provider
    adminUI.BlockAdminRegistry().Register("gallery", &GalleryAdminProvider{store: store})
    
    // Now "Gallery Block" will appear in the block type dropdown
    // and the admin UI will use your custom fields
}
```

### Complete Example: Video Block with Advanced Features

```go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/dracory/cmsstore"
    "github.com/dracory/cmsstore/admin/blocks"
    "github.com/dracory/form"
    "github.com/dracory/hb"
    "github.com/dracory/req"
)

type VideoAdminProvider struct {
    store  cmsstore.StoreInterface
    logger interface{ Error(msg string, args ...interface{}) }
}

func (p *VideoAdminProvider) GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface {
    return []form.FieldInterface{
        form.NewField(form.FieldOptions{
            Label:    "Video URL",
            Name:     "video_url",
            Type:     form.FORM_FIELD_TYPE_STRING,
            Value:    block.Meta("video_url"),
            Required: true,
            Help:     "YouTube, Vimeo, or direct video URL",
        }),
        form.NewField(form.FieldOptions{
            Label: "Video Provider",
            Name:  "video_provider",
            Type:  form.FORM_FIELD_TYPE_SELECT,
            Value: block.Meta("provider"),
            Options: []form.FieldOption{
                {Value: "YouTube", Key: "youtube"},
                {Value: "Vimeo", Key: "vimeo"},
                {Value: "Direct Link", Key: "direct"},
            },
        }),
        form.NewField(form.FieldOptions{
            Label: "Autoplay",
            Name:  "video_autoplay",
            Type:  form.FORM_FIELD_TYPE_CHECKBOX,
            Value: block.Meta("autoplay"),
            Help:  "Start playing automatically when page loads",
        }),
        form.NewField(form.FieldOptions{
            Label: "Show Controls",
            Name:  "video_controls",
            Type:  form.FORM_FIELD_TYPE_CHECKBOX,
            Value: block.Meta("controls"),
            Help:  "Display video player controls",
        }),
        form.NewField(form.FieldOptions{
            Label: "Aspect Ratio",
            Name:  "video_aspect",
            Type:  form.FORM_FIELD_TYPE_SELECT,
            Value: block.Meta("aspect"),
            Options: []form.FieldOption{
                {Value: "16:9 (Widescreen)", Key: "16:9"},
                {Value: "4:3 (Standard)", Key: "4:3"},
                {Value: "1:1 (Square)", Key: "1:1"},
            },
        }),
        // Add a preview section
        &form.Field{
            Type: form.FORM_FIELD_TYPE_RAW,
            Value: hb.Div().
                Class("alert alert-info mt-3").
                Child(hb.Strong().Text("Preview:")).
                Child(hb.BR()).
                Child(hb.Div().
                    ID("video-preview").
                    HTML(p.generatePreview(block))).
                ToHTML(),
        },
    }
}

func (p *VideoAdminProvider) GetTypeLabel() string {
    return "Video Block"
}

func (p *VideoAdminProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
    videoURL := req.GetStringTrimmed(r, "video_url")
    provider := req.GetStringTrimmed(r, "video_provider")
    autoplay := req.GetStringTrimmed(r, "video_autoplay")
    controls := req.GetStringTrimmed(r, "video_controls")
    aspect := req.GetStringTrimmed(r, "video_aspect")
    
    // Validation
    if videoURL == "" {
        return &blocks.ValidationError{Message: "Video URL is required"}
    }
    
    // Store configuration as JSON in content
    config := map[string]string{
        "url":      videoURL,
        "provider": provider,
        "autoplay": autoplay,
        "controls": controls,
        "aspect":   aspect,
    }
    
    configJSON, err := json.Marshal(config)
    if err != nil {
        p.logger.Error("Failed to marshal video config", "error", err)
        return err
    }
    
    block.SetContent(string(configJSON))
    block.SetMeta("video_url", videoURL)
    block.SetMeta("provider", provider)
    block.SetMeta("autoplay", autoplay)
    block.SetMeta("controls", controls)
    block.SetMeta("aspect", aspect)
    
    return nil
}

func (p *VideoAdminProvider) generatePreview(block cmsstore.BlockInterface) string {
    videoURL := block.Meta("video_url")
    if videoURL == "" {
        return "<em>No video URL set</em>"
    }
    return `<iframe src="` + videoURL + `" width="100%" height="300"></iframe>`
}
```

### Example: Interactive Tree Block with Vue.js

```go
package main

import (
    "encoding/json"
    "net/http"
    "github.com/dracory/cmsstore"
    "github.com/dracory/cmsstore/admin/blocks"
    "github.com/dracory/form"
    "github.com/dracory/hb"
    "github.com/dracory/req"
)

type TreeAdminProvider struct {
    store cmsstore.StoreInterface
}

func (p *TreeAdminProvider) GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface {
    return []form.FieldInterface{
        form.NewField(form.FieldOptions{
            Label: "Tree Data (JSON)",
            Name:  "tree_data",
            Type:  form.FORM_FIELD_TYPE_TEXTAREA,
            Value: block.Content(),
            Help:  "Tree structure in JSON format",
        }),
        form.NewField(form.FieldOptions{
            Label: "Tree Style",
            Name:  "tree_style",
            Type:  form.FORM_FIELD_TYPE_SELECT,
            Value: block.Meta("style"),
            Options: []form.FieldOption{
                {Value: "Collapsible", Key: "collapsible"},
                {Value: "Expandable", Key: "expandable"},
                {Value: "Flat List", Key: "flat"},
            },
        }),
        form.NewField(form.FieldOptions{
            Label: "Use Vue.js Renderer",
            Name:  "tree_use_vue",
            Type:  form.FORM_FIELD_TYPE_CHECKBOX,
            Value: block.Meta("use_vue"),
            Help:  "Enable interactive Vue.js rendering on frontend",
        }),
        form.NewField(form.FieldOptions{
            Label: "Custom CSS Class",
            Name:  "tree_css_class",
            Type:  form.FORM_FIELD_TYPE_STRING,
            Value: block.Meta("css_class"),
        }),
        // JSON editor with syntax highlighting
        &form.Field{
            Type: form.FORM_FIELD_TYPE_RAW,
            Value: hb.Script(`
                setTimeout(function() {
                    var textarea = document.querySelector('textarea[name="tree_data"]');
                    if (textarea && typeof CodeMirror !== 'undefined') {
                        var editor = CodeMirror.fromTextArea(textarea, {
                            mode: "application/json",
                            lineNumbers: true,
                            matchBrackets: true,
                            autoCloseBrackets: true
                        });
                        editor.on('change', function() {
                            textarea.value = editor.getValue();
                        });
                    }
                }, 500);
            `).ToHTML(),
        },
    }
}

func (p *TreeAdminProvider) GetTypeLabel() string {
    return "Interactive Tree Block"
}

func (p *TreeAdminProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
    treeData := req.GetStringTrimmed(r, "tree_data")
    style := req.GetStringTrimmed(r, "tree_style")
    useVue := req.GetStringTrimmed(r, "tree_use_vue")
    cssClass := req.GetStringTrimmed(r, "tree_css_class")
    
    // Validate JSON
    if treeData != "" {
        var test interface{}
        if err := json.Unmarshal([]byte(treeData), &test); err != nil {
            return &blocks.ValidationError{Message: "Invalid JSON format: " + err.Error()}
        }
    }
    
    block.SetContent(treeData)
    block.SetMeta("style", style)
    block.SetMeta("use_vue", useVue)
    block.SetMeta("css_class", cssClass)
    
    return nil
}
```

## BlockAdminFieldProvider Interface

```go
type BlockAdminFieldProvider interface {
    // GetContentFields returns form fields for the content editing tab
    GetContentFields(block cmsstore.BlockInterface, r *http.Request) []form.FieldInterface
    
    // GetTypeLabel returns the display name for this block type
    GetTypeLabel() string
    
    // SaveContentFields processes form data and updates the block
    SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error
}
```

## Best Practices

### 1. **Use Metadata for Configuration**
Store block-specific settings in metadata rather than mixing with content:
```go
block.SetContent(actualContent)           // Main content
block.SetMeta("layout", "grid")          // Configuration
block.SetMeta("columns", "3")            // Configuration
```

### 2. **Validate Input**
Always validate form data before saving:
```go
func (p *CustomProvider) SaveContentFields(r *http.Request, block cmsstore.BlockInterface) error {
    value := req.GetStringTrimmed(r, "field_name")
    if value == "" {
        return &blocks.ValidationError{Message: "Field is required"}
    }
    // Save...
    return nil
}
```

### 3. **Provide Helpful Help Text**
Guide users with clear help text:
```go
form.NewField(form.FieldOptions{
    Label: "Image URL",
    Name:  "image_url",
    Help:  "Enter a full URL starting with https://",
})
```

### 4. **Use Appropriate Field Types**
Choose the right form field type for better UX:
- `FORM_FIELD_TYPE_STRING` - Short text
- `FORM_FIELD_TYPE_TEXTAREA` - Long text, JSON, HTML
- `FORM_FIELD_TYPE_SELECT` - Dropdown options
- `FORM_FIELD_TYPE_CHECKBOX` - Boolean flags
- `FORM_FIELD_TYPE_NUMBER` - Numeric values
- `FORM_FIELD_TYPE_RAW` - Custom HTML (for previews, scripts)

### 5. **Add Visual Enhancements**
Use RAW fields for previews, instructions, or custom UI:
```go
&form.Field{
    Type: form.FORM_FIELD_TYPE_RAW,
    Value: hb.Div().
        Class("alert alert-info").
        HTML("<strong>Tip:</strong> Use JSON format for best results").
        ToHTML(),
}
```

### 6. **Handle Store Dependencies**
If your provider needs store access, pass it in the constructor:
```go
type CustomProvider struct {
    store cmsstore.StoreInterface
}

func NewCustomProvider(store cmsstore.StoreInterface) *CustomProvider {
    return &CustomProvider{store: store}
}
```

## Registration

Register your custom admin provider after initializing the admin UI:

```go
adminUI := blocks.UI(config)
adminUI.BlockAdminRegistry().Register("custom_type", customProvider)
```

## What Happens After Registration

1. **Block Creation**: Your block type appears in the "Block type" dropdown
2. **Block Editing**: Content tab uses your `GetContentFields()` method
3. **Block Saving**: Form submission calls your `SaveContentFields()` method
4. **Type Display**: Settings tab shows your `GetTypeLabel()` value

## Thread Safety

The registry is thread-safe and can be accessed concurrently. You can register providers at any time after admin UI initialization.

## Fallback Behavior

If a block type has no registered admin provider:
- Falls back to the HTML provider (basic textarea)
- Shows a warning in the help text
- Still allows basic content editing

This ensures the admin UI never breaks, even with unregistered block types.
