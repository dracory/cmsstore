# [Draft] Enhanced Middleware System

## Status
**[Draft]** - Basic middleware implemented, enhanced features pending

## Summary
- **Problem**: Current middleware implementation lacks flexibility and observability needed for complex request processing
- **Solution**: Enhance middleware system with dynamic configuration, better chaining, monitoring, and recovery capabilities

## Current Implementation (As-Is)

The CMS Store currently uses a standard Go middleware pattern:

```go
// middleware.go
type MiddlewareInterface interface {
    Identifier() string  // Unique ID (e.g., "auth_before")
    Name() string      // Human-friendly label
    Description() string
    Type() string      // "before", "after", "replace"
    Handler() func(next http.Handler) http.Handler
}
```

**Files:**
- `middleware.go` - Core interface and implementation
- `frontend/frontend_middleware.go` - Middleware application logic
- `frontend/frontend.go` - Integration in `PageRenderHtmlBySiteAndAlias()`

**Current Features:**
- ✅ Standard Go http.Handler middleware pattern
- ✅ Identifier-based middleware selection
- ✅ "before" and "after" middleware types
- ✅ Page-level middleware assignment (MiddlewaresBefore/MiddlewaresAfter)
- ✅ Store-level middleware registration
- ✅ Response capture and modification via httptest.ResponseRecorder

**Current Usage:**
```go
// Creating middleware
mw := cmsstore.Middleware().
    SetIdentifier("auth_check").
    SetName("Authentication Check").
    SetType(cmsstore.MIDDLEWARE_TYPE_BEFORE).
    SetHandler(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check auth...
            next.ServeHTTP(w, r)
        })
    })

// Register with store
store.SetMiddlewares([]cmsstore.MiddlewareInterface{mw})

// Assign to page
page.SetMiddlewaresBefore([]string{"auth_check"})
```

## Limitations (Why Enhancement Needed)

Current approach has limitations:
- ❌ No middleware metrics/monitoring
- ❌ No priority/ordering system (manual ordering only)
- ❌ No conditional execution based on request context
- ❌ No panic recovery middleware
- ❌ No per-middleware configuration
- ❌ No rich context sharing (Site, Page, User objects)
- ❌ No Prometheus metrics integration
- ❌ Middleware chain management is manual

## Proposed Enhanced Design (To-Be)

### 1. Enhanced Middleware Interface

```go
type Middleware interface {
    // Core functionality
    Process(ctx *Context, next MiddlewareFunc) error
    
    // Configuration
    Configure(config map[string]interface{}) error
    Priority() int
    
    // Metadata
    Name() string
    Description() string
    Version() string
}

type Context struct {
    Request        *http.Request
    Response       http.ResponseWriter
    Site          *Site
    Page          *Page
    User          *User
    Cache         CacheInterface
    Logger        *slog.Logger
    StartTime     time.Time
    Metrics       *MiddlewareMetrics
    Store         map[string]interface{}
}
```

### 2. Middleware Chain Management

```go
type MiddlewareChain struct {
    middlewares []Middleware
    metrics     *MiddlewareMetrics
    logger      *slog.Logger
}

func (mc *MiddlewareChain) Use(m Middleware) {
    mc.middlewares = append(mc.middlewares, m)
    sort.Slice(mc.middlewares, func(i, j int) bool {
        return mc.middlewares[i].Priority() < mc.middlewares[j].Priority()
    })
}

func (mc *MiddlewareChain) Remove(name string) { ... }
func (mc *MiddlewareChain) Process(ctx *Context) error { ... }
```

### 3. Middleware Metrics (Prometheus)

```go
type MiddlewareMetrics struct {
    executions   *prometheus.CounterVec
    duration     *prometheus.HistogramVec
    errors       *prometheus.CounterVec
    activeCount  *prometheus.GaugeVec
}
```

### 4. Conditional & Recovery Middleware

```go
// Conditional execution
type ConditionalMiddleware struct {
    middleware Middleware
    condition  func(*Context) bool
}

// Panic recovery
type RecoveryMiddleware struct {
    logger *slog.Logger
}
```

## Implementation Status

| Feature | Status | Notes |
|---------|--------|-------|
| Basic middleware interface | ✅ Implemented | `middleware.go` |
| Standard Go handler pattern | ✅ Implemented | `func(next http.Handler) http.Handler` |
| Before/After types | ✅ Implemented | `MIDDLEWARE_TYPE_BEFORE`, `MIDDLEWARE_TYPE_AFTER` |
| Page middleware assignment | ✅ Implemented | `MiddlewaresBefore()`, `MiddlewaresAfter()` |
| Middleware metrics | ❌ Not implemented | No Prometheus integration |
| Priority system | ❌ Not implemented | Manual ordering only |
| Conditional middleware | ❌ Not implemented | No condition functions |
| Recovery middleware | ❌ Not implemented | No panic recovery |
| Rich context | ❌ Not implemented | Standard http.Request only |
| Configuration | ❌ Not implemented | No Configure() method |

## Migration Path

### Option 1: Extend Existing Interface
Add optional methods to current interface with adapter pattern.

### Option 2: New Enhanced Interface
Create separate `EnhancedMiddlewareInterface` while keeping existing one.

```go
type EnhancedMiddlewareInterface interface {
    MiddlewareInterface  // Embed existing
    
    // New methods
    Version() string
    Priority() int
    Configure(config map[string]interface{}) error
    Process(ctx *MiddlewareContext) error
}
```

## Files to Modify (If Implementing)

1. `middleware.go` - Extend interface or create enhanced version
2. `frontend/frontend_middleware.go` - Add chain management and metrics
3. New: `middleware_chain.go` - Chain builder with priority sorting
4. New: `middleware_metrics.go` - Prometheus metrics integration
5. New: `middleware_recovery.go` - Panic recovery middleware
6. New: `middleware_conditional.go` - Conditional execution wrapper

## Risks and Mitigations

1. **Performance Impact**
   - Risk: Overhead from metrics and context wrapping
   - Mitigation: Make enhanced features opt-in

2. **Backward Compatibility**
   - Risk: Breaking existing middleware
   - Mitigation: Keep existing interface, add new as extension

3. **Complexity**
   - Risk: System becomes hard to understand
   - Mitigation: Clear migration guide, examples

4. **Resource Usage**
   - Risk: Memory from rich context objects
   - Mitigation: Lazy loading, object pooling