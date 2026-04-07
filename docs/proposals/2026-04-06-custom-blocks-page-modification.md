# Proposal: Custom Blocks Affecting Page Properties

**Date:** 2026-04-06  
**Status:** Draft  
**Author:** AI Assistant

---

## Problem Statement

Currently, custom blocks in the CMS have no mechanism to modify or influence the page or template they are displayed on. This is a significant limitation for use cases such as:

- **Blog blocks** that need to set the page title based on the blog post being viewed
- **SEO blocks** that need to dynamically set meta descriptions, keywords, or canonical URLs
- **Product detail blocks** that should propagate product information to page-level metadata
- **Dynamic content blocks** that need to modify page-level robots directives

**Important consideration:** Blocks can be displayed in both **pages** and **templates**. The solution must work for both contexts.

The current `BlockType.Render()` interface only returns an HTML string, providing no mechanism for blocks to communicate back to the rendering process.

---

## Goals

1. Enable custom blocks to read and modify page-level properties
2. Maintain backward compatibility with existing block implementations
3. Ensure thread-safe, context-aware access to page data
4. Keep the solution simple and idiomatic to the existing codebase

---

## Proposed Solutions

### Option 1: Context-Based Access (Recommended)

**Overview:** Extend the context system to expose the current renderable (page or template) via context accessors, similar to the existing `RequestFromContext()` pattern.

**Implementation:**

```go
// In context.go

// Context keys for page and template access
const (
    pageContextKey     contextKey = "page"
    templateContextKey contextKey = "template"
)

// PageFromContext retrieves the PageInterface from context when rendering page content
func PageFromContext(ctx context.Context) PageInterface {
    if page, ok := ctx.Value(pageContextKey).(PageInterface); ok {
        return page
    }
    return nil
}

// TemplateFromContext retrieves the TemplateInterface from context when rendering template content
func TemplateFromContext(ctx context.Context) TemplateInterface {
    if template, ok := ctx.Value(templateContextKey).(TemplateInterface); ok {
        return template
    }
    return nil
}

// Renderable interface for common metadata operations on pages and templates
type Renderable interface {
    // Getters
    Title() string
    MetaDescription() string
    MetaKeywords() string
    MetaRobots() string
    CanonicalUrl() string
    
    // Setters
    SetTitle(title string)
    SetMetaDescription(metaDescription string)
    SetMetaKeywords(metaKeywords string)
    SetMetaRobots(metaRobots string)
    SetCanonicalUrl(canonicalUrl string)
}

func RenderableFromContext(ctx context.Context) Renderable {
    if page := PageFromContext(ctx); page != nil {
        return page
    }
    return TemplateFromContext(ctx)
}
```

**Usage in custom blocks:**

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    // Option A: Check for specific type
    if page := cmsstore.PageFromContext(ctx); page != nil {
        page.SetTitle("My Blog Post | " + page.Title())
        page.SetMetaDescription("Blog post excerpt...")
    }
    
    // Option B: Generic renderable (works for both pages and templates)
    if r := cmsstore.RenderableFromContext(ctx); r != nil {
        r.SetTitle("My Blog Post")
        r.SetMetaDescription("Blog post excerpt...")
    }
    
    return renderBlogContent(block), nil
}
```

**Changes Required:**
1. Move `pageContextKey` from `frontend/frontend.go` to `context.go` and add `templateContextKey` constant
2. Add `PageFromContext()`, `TemplateFromContext()`, and `RenderableFromContext()` functions to `context.go`
3. **CRITICAL**: Update `PageRenderHtmlBySiteAndAlias()` to set page in context **before** calling `renderContentToHtml()`
4. Update `TemplateRenderHtmlByID()` to set template in context **before** calling `renderContentToHtml()`
5. Add complete `Renderable` interface with both getters and setters for metadata operations
6. Update both `PageInterface` and `TemplateInterface` to implement `Renderable` interface
7. Update documentation

---

### Option 2: Return Value Extension

**Overview:** Extend the `Render()` return type to include optional page metadata modifications.

**Implementation:**

```go
// New types
type PageModifications struct {
    Title           string
    MetaDescription string
    MetaKeywords    string
    MetaRobots      string
    CanonicalUrl    string
    CustomMetas     map[string]string
}

type RenderResult struct {
    HTML          string
    PageModifications *PageModifications
}

// Updated BlockType interface
type BlockType interface {
    // ... other methods
    Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (*RenderResult, error)
}
```

**Usage in custom blocks:**

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (*RenderResult, error) {
    content := renderBlogContent(block)
    
    return &RenderResult{
        HTML: content,
        PageModifications: &PageModifications{
            Title:           "My Blog Post",
            MetaDescription: "Blog post excerpt...",
        },
    }, nil
}
```

**Changes Required:**
1. Add new `RenderResult` and `PageModifications` types
2. Modify `BlockType` interface `Render()` signature
3. Update `BlockRendererRegistry.RenderBlock()` to handle return value
4. Update all existing block implementations (HTML, Menu, Navbar, Breadcrumbs)
5. Update `renderContentToHtml()` to collect and merge modifications

**Pros:**
- Explicit modification mechanism
- No side effects during rendering
- Can merge multiple block modifications deterministically

**Cons:**
- **BREAKING CHANGE** - all existing block types must be updated
- More complex API
- Requires coordination between multiple blocks

---

### Option 3: Two-Pass Rendering

**Overview:** Split rendering into two phases: metadata collection, then HTML generation.

**Implementation:**

```go
type MetadataCollector interface {
    CollectMetadata(ctx context.Context, block BlockInterface, page PageInterface) error
}

// In frontend:
func (f *frontend) renderContentToHtml(...) (string, error) {
    // Phase 1: Collect metadata from all blocks
    blocks := extractBlocks(content)
    for _, block := range blocks {
        if collector, ok := blockType.(MetadataCollector); ok {
            collector.CollectMetadata(ctx, block, page)
        }
    }
    
    // Phase 2: Render HTML
    html := renderBlocks(blocks)
    return html, nil
}
```

**Usage in custom blocks:**

```go
type BlogBlockType struct{}

func (b *BlogBlockType) CollectMetadata(ctx context.Context, block BlockInterface, page PageInterface) error {
    page.SetTitle("My Blog Post")
    return nil
}

func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    return renderBlogContent(block), nil
}
```

**Changes Required:**
1. Add `MetadataCollector` interface
2. Refactor `renderContentToHtml()` to two-phase approach
3. Update block extraction to happen before rendering
4. Update documentation

**Pros:**
- Separation of concerns (metadata vs rendering)
- Metadata collected before HTML generation
- Optional interface - existing blocks unaffected

**Cons:**
- More complex rendering pipeline
- Requires block extraction to be done separately
- Potential double work (blocks parsed twice)

---

### Option 4: Event/Hook System

**Overview:** Blocks emit events during rendering that the frontend collects and processes.

**Implementation:**

```go
type PageModifierEvent struct {
    BlockID         string
    Title           string
    MetaDescription string
    MetaKeywords    string
    MetaRobots      string
    CanonicalUrl    string
}

func EmitPageModifier(ctx context.Context, event PageModifierEvent) {
    if hooks, ok := ctx.Value(pageModifierHooksKey).([]PageModifierEvent); ok {
        // Append to collection
    }
}

// In frontend:
func (f *frontend) renderContentToHtml(...) (string, error) {
    ctx = withPageModifierHooks(ctx)
    
    html := renderBlocks(ctx, content)
    
    // Apply collected modifications
    events := getPageModifierHooks(ctx)
    for _, event := range events {
        applyPageModification(page, event)
    }
    
    return html, nil
}
```

**Usage in custom blocks:**

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    // Emit modification event
    cmsstore.EmitPageModifier(ctx, cmsstore.PageModifierEvent{
        BlockID:         block.ID(),
        Title:           "My Blog Post",
        MetaDescription: "Blog post excerpt...",
    })
    
    return renderBlogContent(block), nil
}
```

**Changes Required:**
1. Add event types and context key
2. Add `EmitPageModifier()` function
3. Update `renderContentToHtml()` to support hooks
4. Update documentation

**Pros:**
- Decoupled communication
- Can handle multiple block modifications
- Event history available for debugging

**Cons:**
- More complex infrastructure
- Requires understanding of event system
- Less direct than context access

---

### Option 5: Template Variable Injection

**Overview:** Instead of modifying the page directly, blocks set template variables that can be used in the template's `<head>` section.

**Implementation:**

```go
type TemplateContext interface {
    SetVariable(key string, value string)
    GetVariable(key string) string
}

func TemplateContextFromContext(ctx context.Context) TemplateContext {
    // ... similar to RequestFromContext
}
```

**Usage in templates:**

```html
<head>
    <title>[[PageTitle]]</title>
    <meta name="description" content="[[PageMetaDescription]]">
    <!-- Can also use custom variables set by blocks -->
    <meta property="og:title" content="[[OGTitle]]">
</head>
<body>
    [[BLOCK_blog_content]]
</body>
```

**Usage in blocks:**

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    if tc := cmsstore.TemplateContextFromContext(ctx); tc != nil {
        tc.SetVariable("OGTitle", "My Blog Post")
        tc.SetVariable("OGDescription", "Blog post excerpt...")
    }
    return renderBlogContent(block), nil
}
```

**Changes Required:**
1. Add `TemplateContext` interface
2. Add context accessor function
3. Initialize context before rendering
4. Update templates to use new placeholders

**Pros:**
- Non-invasive to Page interface
- Flexible variable system
- Can be used for any template variable, not just SEO
- Could complement Option 1 for additional template variables

**Cons:**
- Requires template updates
- Indirect (not modifying actual page properties)
- Variables only available in current render context
- Doesn't solve the core problem of modifying persistent page metadata

---

### Option 6: Custom Template Variables from Blocks

**Overview:** Allow blocks to set custom template variables (e.g., `[[blog:title]]`, `[[blog:author]]`) that can be used anywhere in the page or template content. This solves the need for blocks to expose data that the surrounding content can reference.

**Current Rendering Order Problem:**

Looking at `frontend/frontend.go` line 562-578, the current flow is:

```go
// Line 571-574: Standard placeholders replaced FIRST
for keyWord, value := range replacementsKeywords {
    content = strings.ReplaceAll(content, "[["+keyWord+"]]", value)
}

// Line 578: Blocks rendered AFTER (can't affect prior replacement)
content, err = frontend.contentRenderBlocks(ctx, content)
```

This means blocks cannot affect `[[PageTitle]]` because the replacement already happened.

**Solution A: Two-Pass Placeholder Replacement**

```go
// In renderContentToHtml():

// Phase 1: Replace standard placeholders
for keyWord, value := range replacementsKeywords {
    content = strings.ReplaceAll(content, "[["+keyWord+"]]", value)
}

// Phase 2: Render blocks (blocks can set custom vars via context)
ctx = WithVarsContext(ctx)  // Initialize variable storage in context
content, err = frontend.contentRenderBlocks(ctx, content)

// Phase 3: Replace custom variables set by blocks
vars := GetVarsFromContext(ctx)
for key, value := range vars {
    content = strings.ReplaceAll(content, "[["+key+"]]", value)
}
```

**Block Usage:**

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    // Set custom variables via context
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        vars.Set("blog:title", "My Blog Post")
        vars.Set("blog:author", "John Doe")
        vars.Set("blog:date", "2024-01-15")
    }
    
    return renderBlogContent(block), nil
}
```

**Template/Page Usage:**

```html
<head>
    <!-- Standard placeholders -->
    <title>[[PageTitle]]</title>
    
    <!-- Custom variables from blog block -->
    <meta property="og:title" content="[[blog:title]]">
    <meta property="article:author" content="[[blog:author]]">
    <meta property="article:published_time" content="[[blog:date]]">
</head>
<body>
    <h1>[[blog:title]]</h1>
    <p>By [[blog:author]] on [[blog:date]]</p>
    
    <!-- The block that sets these variables -->
    [[BLOCK_blog_content]]
</body>
```

**Solution B: Variable Collection Interface**

```go
// Optional interface for blocks that want to expose variables
type VariableCollector interface {
    CollectVars(ctx context.Context, block BlockInterface) map[string]string
}

// In frontend - pre-collect before rendering:
func (f *frontend) renderContentToHtml(...) (string, error) {
    // Step 1: Extract and pre-render blocks to collect variables
    blockIDs := extractBlockIDs(content)
    vars := make(map[string]string)
    
    for _, blockID := range blockIDs {
        block := f.store.BlockFindByID(ctx, blockID)
        blockType := cmsstore.GetBlockType(block.Type())
        
        if collector, ok := blockType.(VariableCollector); ok {
            blockVars := collector.CollectVars(ctx, block)
            for k, v := range blockVars {
                vars[k] = v
            }
        }
    }
    
    // Step 2: Merge custom vars into replacements
    for key, value := range vars {
        replacementsKeywords[key] = value
    }
    
    // Step 3: Replace all placeholders at once
    for keyWord, value := range replacementsKeywords {
        content = strings.ReplaceAll(content, "[["+keyWord+"]]", value)
    }
    
    // Step 4: Render blocks for HTML output
    content = f.contentRenderBlocks(ctx, content)
    
    return content, nil
}
```

**Changes Required:**
1. Add `VarsContext` and `VarsFromContext()` to manage custom variables
2. Modify `renderContentToHtml()` to support two-pass replacement OR variable pre-collection
3. Document naming conventions (e.g., `blocktype:key` format to avoid collisions)
4. Update block development guide

**Pros:**
- Blocks can expose arbitrary data to the surrounding content
- Enables rich integration between blocks and templates
- No breaking changes to existing `BlockType` interface
- Can complement Option 1 (blocks can both modify page metadata AND set custom vars)
- Solves real-world use case: blog post metadata in page `<head>`

**Cons:**
- Requires changes to rendering pipeline
- Variable namespacing needed (prefix with block type to avoid collisions)
- Two-pass rendering adds slight overhead
- Custom variables only available during current render cycle

**Naming Convention Recommendation:**

To avoid collisions between multiple blocks, use namespaced keys:

```go
// Format: blocktype:key
vars.Set("blog:title", "My Post")       // Blog block
echo "product:price", "$99.99") // Product block
vars.Set("event:date", "2024-03-15")    // Event block
```

---

## Comparison Matrix

| Criteria | Option 1 | Option 2 | Option 3 | Option 4 | Option 5 | Option 6 |
|----------|----------|----------|----------|----------|----------|----------|
| **Breaking Changes** | No | Yes | No | No | No | No |
| **Works with Pages** | Yes | Yes | Yes | Yes | Yes | Yes |
| **Works with Templates** | Yes | Yes | Yes | Yes | Yes | Yes |
| **Complexity** | Low | Medium | Medium | High | Medium | Medium |
| **API Surface** | Minimal | Medium | Medium | Large | Medium | Medium |
| **Backward Compat** | Full | Requires migration | Full | Full | Full | Full |
| **Multi-block Coordination** | Last wins | Merge strategy | First wins | Merge strategy | N/A | Last wins |
| **Implementation Effort** | Low | High | Medium | High | Medium | Medium |
| **Pattern Consistency** | High (matches RequestFromContext) | Low | Medium | Medium | Medium | Medium |
| **Custom Variables** | No | No | No | No | Yes | Yes |
| **Modifies Page Metadata** | Yes | Yes | Yes | Yes | No | No |

---

## Recommendation

**Primary Recommendation: Option 1 (Context-Based Page Access)**

This option is recommended because:

1. **Minimal changes** - Only requires moving context key and exporting a function
2. **Pattern consistency** - Follows the existing `RequestFromContext()` pattern already in use
3. **No breaking changes** - Existing block implementations continue to work
4. **Low complexity** - Simple, idiomatic Go solution
5. **Immediate availability** - Can be implemented quickly

**Trade-offs and Considerations:**

1. **"Last block wins" semantics** for multiple blocks modifying the same property:
   - Acceptable because typically only one content block modifies SEO metadata
   - Block rendering order is deterministic (follows content order)
   - Should be clearly documented

2. **Template vs Page context priority:**
   - When rendering a page with a template, blocks will access the page context
   - When rendering a template directly, blocks will access the template context
   - This is the correct behavior as pages override template defaults

3. **Thread safety:**
   - Page/template modifications happen during single-threaded rendering
   - No additional synchronization needed

---

## Implementation Plan (Option 1)

### Phase 1: Core Changes

1. **Update `context.go`:**
   - Move `pageContextKey` from `frontend/frontend.go` to `context.go` (line 85)
   - Add `templateContextKey` constant
   - Add `PageFromContext()`, `TemplateFromContext()` functions
   - Add complete `Renderable` interface with getters and setters
   - Add `RenderableFromContext()` function

2. **Update `frontend.go`:**
   - Remove local `pageContextKey` definition (line 85)
   - **CRITICAL FIX**: Move page context setting in `PageRenderHtmlBySiteAndAlias()` from line 471 to **before** line 455 (before `renderContentToHtml()` call)
   - Update `TemplateRenderHtmlByID()` to set template in context **before** calling `renderContentToHtml()`

3. **Update interfaces:**
   - Ensure `PageInterface` implements `Renderable` (already has all required methods)
   - **ISSUE**: `TemplateInterface` is missing SEO-related methods and needs to be extended:
     - Add `Title()` and `SetTitle()` methods
     - Add `MetaDescription()` and `SetMetaDescription()` methods  
     - Add `MetaKeywords()` and `SetMetaKeywords()` methods
     - Add `MetaRobots()` and `SetMetaRobots()` methods
     - Add `CanonicalUrl()` and `SetCanonicalUrl()` methods

### Phase 2: Documentation

1. Update block development documentation
2. Add examples for page modification in custom blocks
3. Document "last block wins" behavior

### Phase 3: Example Implementation

1. Create example blog block that modifies page title
2. Add tests demonstrating the feature

---

## Important Considerations

### Breaking Changes for TemplateInterface

**Warning**: Extending `TemplateInterface` with SEO methods is a **breaking change** that affects:
- All existing template implementations
- Database schema (may need migration for new fields)
- Template creation/update logic

**Mitigation strategies:**
1. **Gradual rollout**: Implement with default empty values for new fields
2. **Interface segregation**: Create `SEORenderable` interface that only pages implement initially
3. **Template metadata**: Store SEO data in template `Meta()` fields instead of dedicated columns

### Alternative: Page-Only Implementation

For faster implementation with zero breaking changes:

```go
// Modified RenderableFromContext that only works with pages
func RenderableFromContext(ctx context.Context) PageInterface {
    return PageFromContext(ctx) // Returns nil for template-only rendering
}
```

This approach:
- ✅ Solves the primary use case (page SEO modification)
- ✅ Zero breaking changes
- ✅ Can be extended to templates later
- ❌ Blocks in template-only rendering cannot modify metadata

---

## Recommended Combination: Option 1 + Option 6

For a complete solution that covers all use cases:

1. **Option 1 (Context-Based Access)** for modifying standard page metadata like `[[PageTitle]]`
2. **Option 6 (Custom Variables)** for block-specific data like `[[blog:title]]`

**Example combined usage:**

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    // Option 1: Modify page metadata
    if page := cmsstore.PageFromContext(ctx); page != nil {
        page.SetTitle("My Blog Post")
        page.SetMetaDescription("A post about...")
    }
    
    // Option 6: Set custom variables for template
    if vars := cmsstore.VarsFromContext(ctx); vars != nil {
        vars.Set("blog:title", "My Blog Post")
        vars.Set("blog:author", "John Doe")
    }
    
    return renderBlogContent(block), nil
}
```

This gives maximum flexibility:
- Standard metadata (title, description) flow through `[[PageTitle]]`, `[[PageMetaDescription]]`
- Block-specific data flows through `[[blog:title]]`, `[[product:price]]`, etc.
- Both can be used in templates for OpenGraph, Schema.org, or custom display

---

## Alternative Consideration

If the "last block wins" limitation becomes problematic, **Option 3 (Two-Pass Rendering)** could be implemented as an enhancement later. This would:

- Add the `MetadataCollector` interface
- Allow explicit metadata collection before HTML rendering
- Maintain backward compatibility with blocks that don't implement the interface

This hybrid approach would provide the best of both worlds: simple API for basic cases, explicit control for complex scenarios.

---

## Conclusion

The context-based approach (Option 1) provides the most pragmatic solution with minimal changes and maximum backward compatibility. It enables the desired functionality while maintaining the simplicity of the current architecture.
