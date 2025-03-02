# [Draft] Enhanced Template System

## Summary
- **Problem**: Current template system lacks advanced features needed for complex layouts and dynamic content
- **Solution**: Enhance template system with inheritance, composition, dynamic sections, and better caching

## Background

The CMS template system currently provides basic functionality:
- Simple template rendering
- Basic variable substitution
- Limited layout reuse
- No template inheritance
- Basic caching
- Limited dynamic content support

## Detailed Design

### 1. Template Inheritance

```go
type Template struct {
    ID          string
    Name        string
    Content     string
    Parent      *Template
    Sections    map[string]*Section
    Variables   map[string]interface{}
    LastModified time.Time
}

type Section struct {
    Name      string
    Content   string
    Override  bool
    Append    bool
    Prepend   bool
}

// Example template definition
const baseTemplate = `
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
`

const pageTemplate = `
{% extends "base.html" %}

{% block title %}Home Page{% endblock %}

{% block content %}
    <h1>Welcome</h1>
    <div class="content">
        {{ .PageContent }}
    </div>
{% endblock %}
`
```

### 2. Dynamic Sections

```go
type DynamicSection interface {
    Render(ctx *Context) (string, error)
    CacheKey() string
    CacheDuration() time.Duration
}

type MenuSection struct {
    MenuID string
    Cache  CacheInterface
}

func (s *MenuSection) Render(ctx *Context) (string, error) {
    cacheKey := s.CacheKey()
    if cached, ok := s.Cache.Get(cacheKey); ok {
        return cached.(string), nil
    }

    menu, err := ctx.Store.GetMenu(s.MenuID)
    if err != nil {
        return "", err
    }

    html := renderMenu(menu)
    s.Cache.Set(cacheKey, html, s.CacheDuration())
    
    return html, nil
}
```

### 3. Template Composition

```go
type TemplateComponent struct {
    Name     string
    Template string
    Data     interface{}
    Cache    bool
}

// Example component usage
const headerComponent = `
<header>
    {% component "menu" data=site.MainMenu %}
    {% component "search" %}
    {% if user %}
        {% component "user-profile" data=user %}
    {% endif %}
</header>
`

// Component registration
func (t *TemplateEngine) RegisterComponent(name string, component *TemplateComponent) {
    t.components[name] = component
}

// Component rendering
func (t *TemplateEngine) RenderComponent(name string, data interface{}) (string, error) {
    component, ok := t.components[name]
    if !ok {
        return "", fmt.Errorf("component not found: %s", name)
    }

    if component.Cache {
        cacheKey := fmt.Sprintf("component:%s:%v", name, data)
        if cached, ok := t.cache.Get(cacheKey); ok {
            return cached.(string), nil
        }
    }

    result, err := t.renderTemplate(component.Template, data)
    if err != nil {
        return "", err
    }

    if component.Cache {
        t.cache.Set(cacheKey, result, defaultCacheDuration)
    }

    return result, nil
}
```

### 4. Template Functions

```go
type TemplateFuncMap map[string]interface{}

// Register built-in functions
func (t *TemplateEngine) registerBuiltinFuncs() {
    t.funcMap = TemplateFuncMap{
        "date": func(format string, date time.Time) string {
            return date.Format(format)
        },
        "truncate": func(s string, length int) string {
            if len(s) <= length {
                return s
            }
            return s[:length] + "..."
        },
        "markdown": func(content string) template.HTML {
            html := markdown.ToHTML([]byte(content), nil, nil)
            return template.HTML(html)
        },
        "asset": func(path string) string {
            return t.assetManager.GetURL(path)
        },
    }
}

// Custom function registration
func (t *TemplateEngine) RegisterFunc(name string, fn interface{}) error {
    if _, exists := t.funcMap[name]; exists {
        return fmt.Errorf("function already registered: %s", name)
    }
    t.funcMap[name] = fn
    return nil
}
```

### 5. Template Context

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

func NewTemplateContext(site *Site, page *Page) *TemplateContext {
    return &TemplateContext{
        Site:       site,
        Page:       page,
        Data:       make(map[string]interface{}),
        Components: make(map[string]*TemplateComponent),
        Logger:     slog.Default(),
    }
}
```

### 6. Asset Management

```go
type AssetManager struct {
    baseURL     string
    filesystem  fs.FS
    cache       CacheInterface
    fingerprint bool
}

func (am *AssetManager) GetURL(path string) string {
    if !am.fingerprint {
        return filepath.Join(am.baseURL, path)
    }

    hash, err := am.getFileHash(path)
    if err != nil {
        return filepath.Join(am.baseURL, path)
    }

    ext := filepath.Ext(path)
    base := strings.TrimSuffix(path, ext)
    return filepath.Join(am.baseURL, fmt.Sprintf("%s.%s%s", base, hash, ext))
}

func (am *AssetManager) getFileHash(path string) (string, error) {
    cacheKey := fmt.Sprintf("asset:hash:%s", path)
    if hash, ok := am.cache.Get(cacheKey); ok {
        return hash.(string), nil
    }

    content, err := fs.ReadFile(am.filesystem, path)
    if err != nil {
        return "", err
    }

    hash := fmt.Sprintf("%x", md5.Sum(content))[:8]
    am.cache.Set(cacheKey, hash, 24*time.Hour)

    return hash, nil
}
```

## Alternatives Considered

1. **Pure Go Templates**
   - Pros: Native Go support, good performance
   - Cons: Limited features, no inheritance
   - Rejected: Need more advanced templating features

2. **Third-party Template Engine**
   - Pros: Feature-rich, maintained externally
   - Cons: Additional dependency, learning curve
   - Rejected: Need tight CMS integration

3. **Custom DSL**
   - Pros: Complete control, domain-specific features
   - Cons: Complex implementation, maintenance burden
   - Rejected: Template syntax is sufficient

## Implementation Plan

1. Phase 1: Core Features (2 weeks)
   - Implement template inheritance
   - Add dynamic sections
   - Create component system

2. Phase 2: Functions & Context (2 weeks)
   - Add template functions
   - Enhance context system
   - Implement asset management

3. Phase 3: Caching & Performance (1 week)
   - Optimize template parsing
   - Implement caching strategy
   - Add performance metrics

4. Phase 4: Migration & Documentation (2 weeks)
   - Update existing templates
   - Create documentation
   - Add examples

## Risks and Mitigations

1. **Performance**
   - Risk: Template parsing overhead
   - Mitigation: Aggressive caching, precompilation

2. **Complexity**
   - Risk: Template system becomes too complex
   - Mitigation: Good documentation, helper functions

3. **Migration**
   - Risk: Breaking existing templates
   - Mitigation: Backward compatibility layer

4. **Security**
   - Risk: Template injection vulnerabilities
   - Mitigation: Context escaping, security review 