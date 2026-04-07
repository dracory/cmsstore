# Proposal: Custom Block Variable Substitution

**Date:** 2026-04-06  
**Status:** Draft  
**Author:** AI Assistant

---

## Problem Statement

Currently, custom blocks in the CMS have no mechanism to expose dynamic data that can be used in the surrounding page or template content. This is a significant limitation for use cases such as:

- **Blog blocks** that need to expose post title, author, date, category for use in page headers or meta tags
- **Product blocks** that need to expose price, name, SKU for use in structured data or page titles
- **Event blocks** that need to expose event date, location, organizer information
- **User profile blocks** that need to expose user data for personalized content
- **Any dynamic content** that should be referenceable elsewhere in the page/template

**Important consideration:** Blocks can be displayed in both **pages** and **templates**. The solution must work for both contexts.

The current `BlockType.Render()` interface only returns an HTML string, providing no mechanism for blocks to expose data for variable substitution.

---

## Goals

1. Enable custom blocks to set arbitrary variables that can be referenced anywhere in page/template content
2. Maintain backward compatibility with existing block implementations
3. Ensure thread-safe, context-aware variable storage
4. Keep the solution simple and idiomatic to the existing codebase
5. Allow complete flexibility in variable naming (no enforced conventions)

---

## Current Rendering Flow Problem

Looking at `frontend/frontend.go` lines 562-578, the current rendering order is:

```go
// Line 571-574: Standard placeholders replaced FIRST
for keyWord, value := range replacementsKeywords {
    content = strings.ReplaceAll(content, "[["+keyWord+"]]", value)
}

// Line 578: Blocks rendered AFTER (blocks cannot affect placeholder replacement)
content, err = frontend.contentRenderBlocks(ctx, content)
```

This means blocks cannot set variables because placeholder replacement happens before blocks are rendered.

---

## Solution: Context-Based Variable Storage

### Overview

Allow blocks to set arbitrary custom variables via context that can be referenced anywhere in page or template content using `[[variable_name]]` syntax. Users have complete freedom in naming variables.

### Key Concept

**Page/Template Content:**
```html
<h1>Blog. [[blog_title]]</h1>
<p>By [[author_name]] on [[publish_date]]</p>
<div>Price: [[product_price]] [[currency]]</div>
```

**Block Sets Variables:**
```go
vars.Set("blog_title", "Hello World")
vars.Set("author_name", "John Doe")
vars.Set("publish_date", "2026-04-07")
```

**Rendered Output:**
```html
<h1>Blog. Hello World</h1>
<p>By John Doe on 2026-04-07</p>
```

---

## Implementation

### Step 1: Add VarsContext to context.go

```go
// In context.go

type contextKey string

const (
    varsContextKey contextKey = "vars"
)

// VarsContext stores custom variables set by blocks during rendering
type VarsContext struct {
    vars map[string]string
    mu   sync.RWMutex
}

// NewVarsContext creates a new variable context
func NewVarsContext() *VarsContext {
    return &VarsContext{
        vars: make(map[string]string),
    }
}

// Set stores a variable that can be referenced as [[key]] in content
func (v *VarsContext) Set(key, value string) {
    v.mu.Lock()
    defer v.mu.Unlock()
    v.vars[key] = value
}

// Get retrieves a variable value
func (v *VarsContext) Get(key string) (string, bool) {
    v.mu.RLock()
    defer v.mu.RUnlock()
    val, ok := v.vars[key]
    return val, ok
}

// All returns all variables (creates a copy for thread safety)
func (v *VarsContext) All() map[string]string {
    v.mu.RLock()
    defer v.mu.RUnlock()
    
    result := make(map[string]string, len(v.vars))
    for k, v := range v.vars {
        result[k] = v
    }
    return result
}

// WithVarsContext adds a VarsContext to the context
func WithVarsContext(ctx context.Context) context.Context {
    return context.WithValue(ctx, varsContextKey, NewVarsContext())
}

// VarsFromContext retrieves the VarsContext from context
func VarsFromContext(ctx context.Context) *VarsContext {
    if vars, ok := ctx.Value(varsContextKey).(*VarsContext); ok {
        return vars
    }
    return nil
}
```

### Step 2: Update frontend.go renderContentToHtml()

Modify the rendering pipeline to support custom variable replacement:

```go
// In frontend/frontend.go, update renderContentToHtml() around line 562-578

func (frontend *frontend) renderContentToHtml(
    ctx context.Context,
    content string,
    replacementsKeywords map[string]string,
) (string, error) {
    // Phase 1: Replace standard placeholders (PageTitle, SiteName, etc.)
    for keyWord, value := range replacementsKeywords {
        content = strings.ReplaceAll(content, "[["+keyWord+"]]", value)
    }
    
    // Phase 2: Initialize vars context and render blocks
    // Context flows through: contentRenderBlocks -> contentRenderBlockByID -> 
    // fetchBlockContent -> renderBlockByType -> BlockRendererRegistry.RenderBlock -> 
    // BlockType.Render(ctx, block) where blocks can access VarsFromContext(ctx)
    ctx = cmsstore.WithVarsContext(ctx)
    content, err := frontend.contentRenderBlocks(ctx, content)
    if err != nil {
        return "", err
    }
    
    // Phase 3: Replace custom variables set by blocks
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        for key, value := range vars.All() {
            content = strings.ReplaceAll(content, "[["+key+"]]", value)
        }
    }
    
    return content, nil
}
```

### Context Flow

The context flows through the entire rendering pipeline:

1. `renderContentToHtml(ctx, ...)` - Adds VarsContext **pointer** to ctx
2. `contentRenderBlocks(ctx, ...)` - Passes ctx through
3. `contentRenderBlockByID(ctx, ...)` - Passes ctx through
4. `fetchBlockContent(ctx, ...)` - Passes ctx through
5. `renderBlockByType(ctx, block)` - Passes ctx through
6. `BlockRendererRegistry.RenderBlock(ctx, block)` - Passes ctx through
7. `BlockType.Render(ctx, block)` - Block receives ctx with VarsContext!

This is verified in `frontend/block_renderer.go` line 120-121:

```go
globalBlockType := cmsstore.GetBlockType(blockType)
if globalBlockType != nil {
    return globalBlockType.Render(ctx, block)  // Context passed to block!
}
```

### How Variables Get Back

Even though Go passes context by value, the `VarsContext` stored inside is a **pointer**. This means:

```go
// In renderContentToHtml:
ctx = cmsstore.WithVarsContext(ctx)  // Stores *VarsContext pointer in ctx
varsPtr := cmsstore.VarsFromContext(ctx)  // Get the pointer (e.g., 0x12345)

// Context is passed by value through the call chain, but the pointer remains the same

// In BlockType.Render (deep in the call stack):
vars := cmsstore.VarsFromContext(ctx)  // Gets the SAME pointer (0x12345)
vars.Set("blog_title", "Hello")        // Modifies the map at 0x12345

// Back in renderContentToHtml:
vars := cmsstore.VarsFromContext(ctx)  // Still the SAME pointer (0x12345)
vars.All()                             // Returns map with "blog_title" = "Hello"
```

The `VarsContext` struct contains a map and mutex, and since we store a pointer to it in the context, all functions in the call chain share the same underlying data structure. When blocks call `vars.Set()`, they're modifying the shared map, which is then accessible back in `renderContentToHtml`.

**Visual Representation:**

```
renderContentToHtml:
  ctx contains: varsContextKey -> *VarsContext{vars: map[string]string{}}
                                        ↓ (pointer: 0x12345)
  ↓ pass ctx by value
contentRenderBlocks:
  ctx contains: varsContextKey -> *VarsContext{vars: map[string]string{}}
                                        ↓ (same pointer: 0x12345)
  ↓ pass ctx by value
BlockType.Render:
  ctx contains: varsContextKey -> *VarsContext{vars: map[string]string{}}
                                        ↓ (same pointer: 0x12345)
  vars.Set("key", "value")  // Modifies map at 0x12345
  
  ↑ return (context not returned, but pointer was modified)
renderContentToHtml:
  ctx contains: varsContextKey -> *VarsContext{vars: map["key"]="value"}
                                        ↑ (same pointer: 0x12345, now has data)
```

This is why the `VarsContext` type must be a struct with a pointer receiver, not just a plain map.

### Step 3: Usage in Custom Blocks

Blocks can now set arbitrary variables:

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    // Fetch blog post data
    post := fetchBlogPost(block.Data())
    
    // Set custom variables - user chooses any naming convention
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        vars.Set("blog_title", post.Title)
        vars.Set("author_name", post.Author)
        vars.Set("publish_date", post.Date.Format("2006-01-02"))
        vars.Set("category", post.Category)
        
        // Or use different naming styles
        vars.Set("BlogTitle", post.Title)        // PascalCase
        vars.Set("blog:title", post.Title)       // Namespaced
        vars.Set("$blogTitle", post.Title)       // Prefixed
        vars.Set("user_name", post.Author)       // Snake case
    }
    
    return renderBlogHTML(post), nil
}
```

### Step 4: Usage in Pages/Templates

Variables can be referenced anywhere in content:

```html
<!DOCTYPE html>
<html>
<head>
    <title>[[PageTitle]] - [[blog_title]]</title>
    <meta name="description" content="[[PageMetaDescription]]">
    <meta property="og:title" content="[[blog_title]]">
    <meta property="article:author" content="[[author_name]]">
    <meta property="article:published_time" content="[[publish_date]]">
</head>
<body>
    <header>
        <h1>[[blog_title]]</h1>
        <p>By [[author_name]] on [[publish_date]] in [[category]]</p>
    </header>
    
    <main>
        [[BLOCK_blog_content]]
    </main>
    
    <footer>
        <p>Author: [[user_name]]</p>
    </footer>
</body>
</html>
```

---

## Variable Naming

Users have complete freedom in naming variables. Common patterns:

- `snake_case`: `blog_title`, `user_name`, `product_price`
- `camelCase`: `blogTitle`, `userName`, `productPrice`
- `PascalCase`: `BlogTitle`, `UserName`, `ProductPrice`
- `namespaced`: `blog:title`, `product:price`, `event:date`
- `prefixed`: `$blogTitle`, `@userName`, `#productPrice`

The system doesn't enforce any convention - users choose what works for their project.

---

## Multiple Blocks

Multiple blocks can set different variables without conflicts:

```go
// Blog block sets blog-related variables
vars.Set("blog_title", "My Post")
vars.Set("blog_author", "John")

// Product block sets product-related variables
vars.Set("product_name", "Widget")
vars.Set("product_price", "$99")

// User block sets user-related variables
vars.Set("user_name", "Jane")
vars.Set("user_role", "Admin")
```

All variables are available for substitution in the final content.

---

## Variable Collision Handling

If multiple blocks set the same variable name, the last block wins (based on rendering order in content). This is simple and predictable. Users should choose unique names or use namespacing to avoid collisions.

---

## Changes Required

1. **context.go**
   - Add `VarsContext` type with thread-safe variable storage
   - Add `WithVarsContext()` and `VarsFromContext()` functions
   - Add `varsContextKey` constant

2. **frontend/frontend.go**
   - Update `renderContentToHtml()` to initialize vars context before block rendering
   - Add third phase to replace custom variables after blocks are rendered
   - Ensure context is passed through to `contentRenderBlocks()`

3. **Documentation**
   - Add block development guide section on setting custom variables
   - Provide examples of different naming conventions
   - Document variable collision behavior
   - Add examples for common use cases (blog, products, events, users)

---

## Benefits

1. **Complete flexibility** - Users choose any variable naming convention
2. **No breaking changes** - Existing blocks continue to work unchanged
3. **Simple API** - Just `vars.Set(key, value)` in block render
4. **Thread-safe** - Mutex-protected variable storage
5. **Composable** - Multiple blocks can set different variables
6. **Works everywhere** - Variables can be used in pages, templates, and block content
7. **Predictable** - Clear rendering order and collision handling

---

## Use Cases

### Blog Post Metadata

```go
vars.Set("post_title", "Understanding Go Contexts")
vars.Set("post_author", "Jane Developer")
vars.Set("post_date", "2026-04-07")
vars.Set("post_reading_time", "5 min read")
```

### E-commerce Product

```go
vars.Set("product_name", "Wireless Headphones")
vars.Set("product_price", "$149.99")
vars.Set("product_sku", "WH-1000XM4")
vars.Set("product_availability", "In Stock")
```

### Event Information

```go
vars.Set("event_name", "Go Conference 2026")
vars.Set("event_date", "2026-09-15")
vars.Set("event_location", "San Francisco, CA")
vars.Set("event_capacity", "500 attendees")
```

### User Profile

```go
vars.Set("user_display_name", "John Doe")
vars.Set("user_role", "Administrator")
vars.Set("user_join_date", "2024-01-15")
vars.Set("user_post_count", "42")
```

---

## Testing Strategy

### Unit Tests

```go
func TestVarsContext(t *testing.T) {
    ctx := cmsstore.WithVarsContext(context.Background())
    vars := cmsstore.VarsFromContext(ctx)
    
    // Test Set and Get
    vars.Set("test_key", "test_value")
    val, ok := vars.Get("test_key")
    assert.True(t, ok)
    assert.Equal(t, "test_value", val)
    
    // Test All
    vars.Set("key1", "value1")
    vars.Set("key2", "value2")
    all := vars.All()
    assert.Equal(t, 3, len(all))
}

func TestVariableReplacement(t *testing.T) {
    // Test that custom variables are replaced after blocks render
    content := "<h1>[[custom_title]]</h1>[[BLOCK_test]]"
    
    // Mock block that sets custom_title
    // ... render content
    
    assert.Contains(t, result, "<h1>My Custom Title</h1>")
}

func TestMultipleBlocksVariables(t *testing.T) {
    // Test that multiple blocks can set different variables
    // ... test scenario
}
```

### Integration Tests

1. Create test page with custom variable placeholders
2. Create test block that sets those variables
3. Render page and verify variables are replaced
4. Test variable collision (last block wins)
5. Test with both pages and templates

---

## Migration Path

This feature requires no migration:

1. Existing blocks work unchanged (they simply don't set variables)
2. Existing pages/templates work unchanged (unknown variables remain as `[[var]]`)
3. New blocks can opt-in by using `VarsFromContext()`
4. New pages/templates can reference custom variables as needed

---

## Performance Considerations

1. **Variable storage** - Small map in context, minimal memory overhead
2. **Replacement** - Additional string replacement pass after block rendering
3. **Thread safety** - Mutex overhead only when blocks actually set variables
4. **Optimization** - Skip variable replacement if no variables were set

Potential optimization:

```go
// Only do replacement if variables were actually set
if vars := cmsstore.VarsFromContext(ctx); vars != nil && len(vars.All()) > 0 {
    for key, value := range vars.All() {
        content = strings.ReplaceAll(content, "[["+key+"]]", value)
    }
}
```

---

## Documentation Updates

### Block Development Guide

Add new section: "Setting Custom Variables"

```markdown
## Setting Custom Variables

Blocks can expose data as custom variables that can be referenced anywhere in the page or template content.

### Basic Usage

```go
func (b *MyBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        vars.Set("my_variable", "my value")
    }
    return html, nil
}
```

### Naming Conventions

You can use any naming convention that suits your project:

- Snake case: `blog_title`, `user_name`
- Camel case: `blogTitle`, `userName`
- Namespaced: `blog:title`, `user:name`
- Prefixed: `$blogTitle`, `@userName`

### Using Variables in Content

Reference variables using `[[variable_name]]` syntax:

```html
<h1>[[blog_title]]</h1>
<p>By [[author_name]] on [[publish_date]]</p>
```

### Variable Collisions

If multiple blocks set the same variable name, the last block (in content order) wins. Use unique names or namespacing to avoid conflicts.
```

---

## Future Enhancements

Potential future additions (not in initial implementation):

1. **Variable scoping** - Block-local vs global variables
2. **Variable types** - Support for non-string values (numbers, booleans)
3. **Variable functions** - `[[uppercase:blog_title]]`, `[[format_date:publish_date]]`
4. **Variable defaults** - `[[blog_title|Default Title]]`
5. **Conditional variables** - `[[if:user_logged_in]]...[[endif]]`

These can be added later without breaking changes to the core variable system.

---

## Conclusion

The custom variable substitution approach provides maximum flexibility with minimal complexity. Users can set any variable names they want, and blocks can expose arbitrary data to the surrounding content. This solves the original problem while maintaining backward compatibility and keeping the implementation simple.
