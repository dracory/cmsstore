# [Draft] Enhanced Shortcode System

## Summary
- **Problem**: Current shortcode implementation is basic and lacks features needed for complex content generation
- **Solution**: Enhance shortcode system with validation, caching, async processing, and better developer tooling

## Background

The CMS currently supports shortcodes through a basic interface:
```go
type ShortcodeInterface interface {
    Alias() string
    Description() string
    Render(r *http.Request, s string, m map[string]string) string
}
```

While functional, this approach has limitations:
- No parameter validation
- Limited error handling
- No caching support
- Synchronous processing only
- Basic parameter parsing
- Limited development tools

## Detailed Design

### 1. Enhanced Shortcode Interface

```go
type ShortcodeInterface interface {
    // Basic Information
    Alias() string
    Description() string
    Version() string
    
    // Parameter Definition & Validation
    Parameters() []ShortcodeParameter
    Validate(params map[string]string) error
    
    // Rendering
    Render(ctx ShortcodeContext) (string, error)
    RenderAsync(ctx ShortcodeContext) (<-chan ShortcodeResult, error)
    
    // Caching
    CacheKey(params map[string]string) string
    CacheDuration() time.Duration
}

type ShortcodeParameter struct {
    Name        string
    Type        string // string, int, float, bool, array
    Required    bool
    Default     interface{}
    Validation  []string // regex, min, max, etc.
    Description string
}

type ShortcodeContext struct {
    Request     *http.Request
    Content     string
    Parameters  map[string]string
    Cache       ShortcodeCacheInterface
    Logger      *slog.Logger
}

type ShortcodeResult struct {
    Content string
    Error   error
    Cache   bool
}
```

### 2. Parameter Validation System

```go
type ShortcodeValidator interface {
    Validate(param ShortcodeParameter, value string) error
}

// Example validators
type RegexValidator struct {
    Pattern string
}

type RangeValidator struct {
    Min float64
    Max float64
}

// Example usage
func NewProductListShortcode() ShortcodeInterface {
    return &ProductListShortcode{
        parameters: []ShortcodeParameter{
            {
                Name:     "category",
                Type:     "string",
                Required: true,
                Validation: []string{
                    "regex:^[a-zA-Z0-9-]+$",
                },
            },
            {
                Name:     "limit",
                Type:     "int",
                Default:  10,
                Validation: []string{
                    "range:1,100",
                },
            },
        },
    }
}
```

### 3. Caching System

```go
type ShortcodeCacheInterface interface {
    Get(key string) (string, bool)
    Set(key string, value string, duration time.Duration)
    Delete(key string)
}

// Example implementation
func (s *ProductListShortcode) CacheKey(params map[string]string) string {
    return fmt.Sprintf("shortcode:productlist:%s:%s",
        params["category"],
        params["limit"],
    )
}

func (s *ProductListShortcode) CacheDuration() time.Duration {
    return 5 * time.Minute
}
```

### 4. Async Processing

```go
func (s *ProductListShortcode) RenderAsync(ctx ShortcodeContext) (<-chan ShortcodeResult, error) {
    resultChan := make(chan ShortcodeResult)
    
    go func() {
        defer close(resultChan)
        
        // Fetch products asynchronously
        products, err := s.fetchProducts(ctx.Parameters["category"])
        if err != nil {
            resultChan <- ShortcodeResult{Error: err}
            return
        }
        
        // Render HTML
        html := s.renderProducts(products)
        
        resultChan <- ShortcodeResult{
            Content: html,
            Cache:   true,
        }
    }()
    
    return resultChan, nil
}
```

### 5. Developer Tools

```go
// CLI tool for shortcode development
type ShortcodeTester struct {
    shortcode ShortcodeInterface
    logger    *slog.Logger
}

func (t *ShortcodeTester) Test(params map[string]string) {
    // Validate parameters
    if err := t.shortcode.Validate(params); err != nil {
        t.logger.Error("Validation failed", "error", err)
        return
    }
    
    // Test rendering
    ctx := ShortcodeContext{
        Parameters: params,
        Cache:     NewMemoryCache(),
        Logger:    t.logger,
    }
    
    result, err := t.shortcode.Render(ctx)
    if err != nil {
        t.logger.Error("Rendering failed", "error", err)
        return
    }
    
    t.logger.Info("Render successful",
        "result", result,
        "cache_key", t.shortcode.CacheKey(params),
    )
}
```

### 6. Example Implementation

```go
type ProductListShortcode struct {
    parameters []ShortcodeParameter
    db        *sql.DB
    cache     ShortcodeCacheInterface
}

func (s *ProductListShortcode) Render(ctx ShortcodeContext) (string, error) {
    // Parameter validation
    if err := s.Validate(ctx.Parameters); err != nil {
        return "", err
    }
    
    // Check cache
    if cached, ok := ctx.Cache.Get(s.CacheKey(ctx.Parameters)); ok {
        return cached, nil
    }
    
    // Fetch and render products
    products, err := s.fetchProducts(ctx.Parameters["category"])
    if err != nil {
        return "", err
    }
    
    html := s.renderProducts(products)
    
    // Cache result
    ctx.Cache.Set(s.CacheKey(ctx.Parameters), html, s.CacheDuration())
    
    return html, nil
}
```

## Alternatives Considered

1. **Template-based System**
   - Pros: Familiar syntax, built-in functions
   - Cons: Less flexibility, harder to extend
   - Rejected: Need more programmatic control

2. **Plugin System**
   - Pros: Complete isolation, security
   - Cons: Complex deployment, performance overhead
   - Rejected: Overkill for current needs

3. **DSL-based Approach**
   - Pros: More expressive
   - Cons: Learning curve, parsing complexity
   - Rejected: Simple tag-based system is sufficient

## Implementation Plan

1. Phase 1: Core Enhancement (2 weeks)
   - Implement new interfaces
   - Add parameter validation
   - Create basic caching

2. Phase 2: Async & Tools (2 weeks)
   - Add async processing
   - Create developer tools
   - Write documentation

3. Phase 3: Migration (2 weeks)
   - Update existing shortcodes
   - Add tests
   - Create examples

4. Phase 4: Optimization (1 week)
   - Performance testing
   - Cache tuning
   - Documentation updates

## Risks and Mitigations

1. **Backward Compatibility**
   - Risk: Breaking existing shortcodes
   - Mitigation: Adapter layer, gradual migration

2. **Performance**
   - Risk: Overhead from validation/caching
   - Mitigation: Benchmark-driven optimization

3. **Complexity**
   - Risk: System becomes too complex
   - Mitigation: Good documentation, helper functions

4. **Resource Usage**
   - Risk: Memory usage from caching
   - Mitigation: Configurable limits, monitoring 