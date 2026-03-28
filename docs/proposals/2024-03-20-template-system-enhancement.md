# [Draft] Enhanced Template System

## Status
**[Draft]** - Basic template system implemented, advanced features pending

## Summary
- **Problem**: Current template system provides basic content wrapping but lacks advanced features for complex layouts
- **Solution**: Enhance template system with inheritance, composition, dynamic sections, and template functions

## Current Implementation (As-Is)

The CMS Store currently has a simple template system:

**Template Entity:**
```go
// template_implementation.go
type TemplateInterface interface {
    ID() string
    Name() string
    Content() string
    Handle() string
    Status() string
    // ... other fields
}
```

**Files:**
- `template_implementation.go` - Template entity with dataobject pattern
- `template_query.go` - Database queries for templates
- `frontend/frontend.go` - Template rendering in `pageOrTemplateContent()`

**Current Features:**
- Template entity with CRUD operations
- Page-to-template assignment via `page.TemplateID()`
- Template content serves as wrapper for page content
- Simple placeholder replacement (`[[PageTitle]]`, `[[PageContent]]`, etc.)
- Template status (active/inactive)
- Versioning support

**Current Template Usage Flow:**
```go
// 1. Page requests render
// 2. If page has TemplateID, load template content
// 3. Template content becomes the base
// 4. Page content available via [[PageContent]] placeholder
// 5. Other placeholders: [[PageTitle]], [[PageMetaDescription]], etc.

// Example template content:
// <html>
//   <head><title>[[PageTitle]]</title></head>
//   <body>[[PageContent]]</body>
// </html>
```

**Current Limitations:**
- No template inheritance (no extends/block system)
- No dynamic sections or section overriding
- No template composition/components
- No template functions (date, truncate, markdown, etc.)
- No rich template context (only basic placeholders)
- No asset management or fingerprinting
- No conditional logic in templates
- No template caching beyond basic TTL

## Proposed Enhanced Design (To-Be)

### 1. Template Inheritance (Django/Jinja-style)

```go
type Template struct {
    ID          string
    Name        string
    Content     string
    Parent      *Template        // Parent template for inheritance
    Sections    map[string]*Section
}

type Section struct {
    Name     string
    Content  string
    Override bool
    Append   bool
    Prepend  bool
}
```

**Template Syntax:**
```html
<!-- base.html -->
<!DOCTYPE html>
<html>
<head>
    {% block head %}
    <title>{% block title %}{% endblock %} - My Site</title>
    {% endblock %}
</head>
<body>
    {% block header %}{% endblock %}
    {% block content %}{% endblock %}
    {% block footer %}{% endblock %}
</body>
</html>

<!-- page.html extends base -->
{% extends "base.html" %}
{% block title %}Home Page{% endblock %}
{% block content %}
    <h1>Welcome</h1>
    <div class="content">{{ .PageContent }}</div>
{% endblock %}
```

### 2. Template Composition

```go
type TemplateComponent struct {
    Name     string
    Template string
    Data     interface{}
    Cache    bool
}
```

**Component Usage:**
```html
<header>
    {% component "menu" data=site.MainMenu %}
    {% component "search" %}
    {% if user %}
        {% component "user-profile" data=user %}
    {% endif %}
</header>
```

### 3. Template Functions

```go
type TemplateFuncMap map[string]interface{}

// Built-in functions:
// - date(format, time) - Format dates
// - truncate(string, length) - Truncate text
// - markdown(content) - Convert markdown to HTML
// - asset(path) - Get asset URL with fingerprinting
// - url(pageID) - Generate page URL
```

### 4. Rich Template Context

```go
type TemplateContext struct {
    Site       *Site
    Page       *Page
    User       *User
    Request    *http.Request
    Data       map[string]interface{}
    Components map[string]*TemplateComponent
    Cache      CacheInterface
    Logger     *slog.Logger
}
```

### 5. Asset Management

```go
type AssetManager struct {
    baseURL     string
    filesystem  fs.FS
    cache       CacheInterface
    fingerprint bool  // Add content hash to filenames
}

// Usage in template:
// <link rel="stylesheet" href="{{ asset('css/main.css') }}">
// → /css/main.a3f7b2c.css
```

## Implementation Status

| Feature | Status | Notes |
|---------|--------|-------|
| Template entity (CRUD) | Implemented | `template_implementation.go` |
| Page-template assignment | Implemented | `page.TemplateID()` |
| Basic placeholder replacement | Implemented | `[[PageTitle]]`, `[[PageContent]]` |
| Template status | Implemented | Active/Inactive |
| Template versioning | Implemented | Via versioning system |
| Template inheritance | Not implemented | No extends/block syntax |
| Dynamic sections | Not implemented | No section override |
| Template components | Not implemented | No component system |
| Template functions | Not implemented | No date/truncate/etc |
| Rich context | Not implemented | Only basic placeholders |
| Asset management | Not implemented | No fingerprinting |
| Conditional logic | Not implemented | No if/else in templates |

## Migration Strategy

### Option 1: New Template Engine
Create separate `TemplateEngine` with enhanced features, keep existing for backward compatibility.

### Option 2: Extend Existing
Add features incrementally without breaking existing templates.

**Example Migration Path:**
1. Phase 1: Add template functions (backward compatible)
2. Phase 2: Add component system (new opt-in syntax)
3. Phase 3: Add inheritance (new opt-in syntax)
4. Phase 4: Deprecate simple placeholders over time

## Files to Modify (If Implementing)

1. New: `template_engine.go` - Template parsing and rendering
2. New: `template_funcs.go` - Built-in template functions
3. New: `template_inheritance.go` - Block/extends system
4. New: `template_components.go` - Component system
5. New: `asset_manager.go` - Asset fingerprinting
6. `frontend/frontend.go` - Integrate new engine alongside existing

## Risks and Mitigations

1. **Performance**
   - Risk: Template parsing overhead
   - Mitigation: Pre-compile templates, aggressive caching

2. **Complexity**
   - Risk: Template system becomes too complex
   - Mitigation: Start simple, add features incrementally

3. **Migration**
   - Risk: Breaking existing templates
   - Mitigation: Backward compatibility layer, gradual rollout

4. **Security**
   - Risk: Template injection vulnerabilities
   - Mitigation: Context escaping, security review