# [Draft] Enhanced Shortcode System

## Status
**[Draft]** - Basic shortcode interface implemented, enhanced features pending

## Summary
- **Problem**: Current shortcode implementation is basic and lacks features needed for complex content generation
- **Solution**: Enhance shortcode system with validation, caching, async processing, and better developer tooling

## Current Implementation (As-Is)

The CMS Store currently uses a basic shortcode interface:

```go
// shortcode_interface.go
type ShortcodeInterface interface {
    Alias() string
    Description() string
    Render(r *http.Request, s string, m map[string]string) string
}
```

**Files:**
- `shortcode_interface.go` - Core interface definition
- `frontend/frontend.go` - Shortcode application in `applyShortcodes()`
- `frontend/new.go` - Shortcode registration via `Config.Shortcodes`
- `store.go` - `Shortcodes()` method on store

**Current Features:**
- Basic `Alias()`, `Description()`, `Render()` interface
- Shortcode registration at frontend initialization
- Integration with `shortcode` package (uses `<shortcode>` brackets)
- Access to HTTP request context and parameters map
- Store-level shortcodes + frontend-level shortcodes

**Current Usage:**
```go
// Registering shortcodes
frontend := cmsstore.NewFrontend(cmsstore.Config{
    Shortcodes: []cmsstore.ShortcodeInterface{
        myCustomShortcode,
    },
})

// In content
<myshortcode param1="value1" param2="value2">
```

## Limitations (Why Enhancement Needed)

Current approach has limitations:
- No parameter validation (runtime errors only)
- No structured error handling (returns string, not error)
- No built-in caching support per shortcode
- No async/render streaming support
- No parameter type definitions or defaults
- No versioning information
- No dedicated shortcode context (just raw HTTP request)

## Proposed Enhanced Design (To-Be)

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
    
    // Rendering with error handling
    Render(ctx ShortcodeContext) (string, error)
    
    // Async support
    RenderAsync(ctx ShortcodeContext) (<-chan ShortcodeResult, error)
    
    // Per-shortcode caching
    CacheKey(params map[string]string) string
    CacheDuration() time.Duration
}
```

### 2. Parameter Definition System

```go
type ShortcodeParameter struct {
    Name        string
    Type        string // string, int, float, bool, array
    Required    bool
    Default     interface{}
    Validation  []string // regex, min, max, etc.
    Description string
}
```

### 3. Structured Context

```go
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

### 4. Validation System

```go
type ShortcodeValidator interface {
    Validate(param ShortcodeParameter, value string) error
}

// Built-in validators
// - RegexValidator
// - RangeValidator  
// - RequiredValidator
// - TypeValidator
```

## Implementation Status

| Feature | Status | Notes |
|---------|--------|-------|
| Basic interface (Alias, Description, Render) | Implemented | `shortcode_interface.go` |
| Parameter validation | Not implemented | No validation framework |
| Structured error handling | Not implemented | Returns string, not (string, error) |
| Shortcode context | Not implemented | Raw request + map only |
| Per-shortcode caching | Not implemented | No CacheKey/CacheDuration methods |
| Async processing | Not implemented | No RenderAsync method |
| Version tracking | Not implemented | No Version() method |
| Parameter definitions | Not implemented | No Parameters() method |

## Migration Path

### Option 1: Extend Interface (Breaking Change)
Add new methods to existing interface - requires updating all shortcodes.

### Option 2: New Interface (Backward Compatible)
Create `EnhancedShortcodeInterface` that embeds `ShortcodeInterface`.

```go
type EnhancedShortcodeInterface interface {
    ShortcodeInterface
    
    // New methods
    Version() string
    Parameters() []ShortcodeParameter
    Validate(params map[string]string) error
    RenderEnhanced(ctx ShortcodeContext) (string, error)
}
```

## Files to Modify (If Implementing)

1. `shortcode_interface.go` - Extend or create new enhanced interface
2. `frontend/frontend.go` - Update `applyShortcodes()` to support new features
3. New: `shortcode_validator.go` - Validation framework
4. New: `shortcode_context.go` - Context struct and cache interface
5. New: `shortcode_test_helper.go` - Developer testing tools

## Risks and Mitigations

1. **Backward Compatibility**
   - Risk: Breaking existing shortcodes
   - Mitigation: Use new interface or adapter pattern

2. **Performance**
   - Risk: Validation/caching overhead
   - Mitigation: Make features opt-in per shortcode

3. **Complexity**
   - Risk: System becomes too complex
   - Mitigation: Helper functions, good docs

4. **Resource Usage**
   - Risk: Memory from per-shortcode caching
   - Mitigation: Configurable limits