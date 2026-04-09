# Block Custom Variables Metadata

**Date**: 2026-04-09  
**Status**: Proposed

## Summary

Add a method to the `BlockType` interface to expose metadata about custom variables that a block sets during rendering. This enables the admin interface to display available variables and their descriptions to content editors.

## Problem Statement

Blocks can set custom variables using `VarsFromContext(ctx)` during rendering:

```go
func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
    if vars := VarsFromContext(ctx); vars != nil {
        vars.Set("blog_title", post.Title)
        vars.Set("blog_author", post.Author)
        vars.Set("blog_date", post.Date.Format("2006-01-02"))
    }
    return html, nil
}
```

These variables are referenced in page/template content as `[[blog_title]]`, `[[blog_author]]`, etc. However:

1. Content editors have no way to discover what variables a block exposes
2. Variable names and purposes are not documented anywhere in the admin UI
3. Editors must guess names or consult external documentation

## Proposed Solution

Add `GetCustomVariables()` to the `BlockType` interface and a simple `BlockCustomVariable` struct with just name and description:

```go
type BlockType interface {
    TypeKey() string
    TypeLabel() string
    Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error)
    GetAdminFields(block BlockInterface, r *http.Request) interface{}
    SaveAdminFields(r *http.Request, block BlockInterface) error

    // GetCustomVariables returns metadata about custom variables this block type
    // can set during rendering. Returns nil or empty slice if none.
    GetCustomVariables() []BlockCustomVariable
}

// BlockCustomVariable describes a custom variable that a block type can set.
// Variables are always strings and are referenced in content as [[name]].
type BlockCustomVariable struct {
    Name        string // Variable name, e.g. "blog_title"
    Description string // Human-readable description of what the variable contains
}
```

## Implementation Details

### 1. Struct and Interface in `block_type.go`

Add `BlockCustomVariable` struct and `GetCustomVariables()` to the `BlockType` interface.

### 2. Built-in Block Types

Built-in blocks (HTML, Menu, Navbar, Breadcrumbs) don't set custom variables:

```go
func (t *HTMLBlockType) GetCustomVariables() []BlockCustomVariable {
    return nil
}
```

### 3. Custom Block Example

```go
func (b *BlogBlockType) GetCustomVariables() []BlockCustomVariable {
    return []BlockCustomVariable{
        {Name: "blog_title",   Description: "The blog post title"},
        {Name: "blog_author",  Description: "The post author name"},
        {Name: "blog_date",    Description: "Publication date in YYYY-MM-DD format"},
        {Name: "blog_excerpt", Description: "Short summary of the post"},
    }
}
```

### 4. Admin UI Integration

In the block edit form, when a block type exposes custom variables, show a reference table:

```html
<div class="custom-variables-section">
    <h3>Available Variables</h3>
    <p>This block sets the following variables for use in page/template content:</p>
    <table>
        <thead>
            <tr><th>Variable</th><th>Description</th></tr>
        </thead>
        <tbody>
            <tr><td><code>[[blog_title]]</code></td><td>The blog post title</td></tr>
            <tr><td><code>[[blog_author]]</code></td><td>The post author name</td></tr>
        </tbody>
    </table>
</div>
```

## Backward Compatibility

Returning `nil` is valid, so existing block types are unaffected once they add the method stub. No breaking changes to rendering or storage.

## Migration Path

1. Add `BlockCustomVariable` struct and `GetCustomVariables()` to `BlockType` interface in `block_type.go`
2. Update all built-in block types to return `nil`
3. Update documentation and examples
4. Implement admin UI display in `admin/blocks/block_update_controller.go`
