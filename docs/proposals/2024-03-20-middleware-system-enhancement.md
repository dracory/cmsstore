# [Draft] Enhanced Middleware System

## Summary
- **Problem**: Current middleware implementation lacks flexibility and observability needed for complex request processing
- **Solution**: Enhance middleware system with dynamic configuration, better chaining, monitoring, and recovery capabilities

## Background

The CMS uses middleware for request processing, but the current implementation has limitations:
- Static middleware configuration
- Limited error recovery
- No middleware-specific metrics
- Basic middleware chaining
- No conditional middleware execution
- Limited context sharing

## Detailed Design

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

type MiddlewareFunc func(ctx *Context) error

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

func (mc *MiddlewareChain) Remove(name string) {
    for i, m := range mc.middlewares {
        if m.Name() == name {
            mc.middlewares = append(mc.middlewares[:i], mc.middlewares[i+1:]...)
            return
        }
    }
}

func (mc *MiddlewareChain) Process(ctx *Context) error {
    if len(mc.middlewares) == 0 {
        return nil
    }
    
    return mc.processMiddleware(ctx, 0)
}

func (mc *MiddlewareChain) processMiddleware(ctx *Context, index int) error {
    if index >= len(mc.middlewares) {
        return nil
    }
    
    m := mc.middlewares[index]
    
    // Create next function
    next := func(ctx *Context) error {
        return mc.processMiddleware(ctx, index+1)
    }
    
    // Process with metrics
    start := time.Now()
    err := m.Process(ctx, next)
    duration := time.Since(start)
    
    mc.metrics.RecordMiddleware(m.Name(), duration, err)
    
    return err
}
```

### 3. Middleware Metrics

```go
type MiddlewareMetrics struct {
    executions   *prometheus.CounterVec
    duration     *prometheus.HistogramVec
    errors       *prometheus.CounterVec
    activeCount  *prometheus.GaugeVec
}

func NewMiddlewareMetrics() *MiddlewareMetrics {
    return &MiddlewareMetrics{
        executions: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "cms_middleware_executions_total",
                Help: "Total number of middleware executions",
            },
            []string{"middleware"},
        ),
        duration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "cms_middleware_duration_seconds",
                Help: "Duration of middleware execution",
            },
            []string{"middleware"},
        ),
        errors: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "cms_middleware_errors_total",
                Help: "Total number of middleware errors",
            },
            []string{"middleware", "error_type"},
        ),
        activeCount: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "cms_middleware_active_requests",
                Help: "Number of active requests in middleware",
            },
            []string{"middleware"},
        ),
    }
}
```

### 4. Conditional Middleware

```go
type ConditionFunc func(*Context) bool

type ConditionalMiddleware struct {
    middleware Middleware
    condition  ConditionFunc
}

func (cm *ConditionalMiddleware) Process(ctx *Context, next MiddlewareFunc) error {
    if cm.condition(ctx) {
        return cm.middleware.Process(ctx, next)
    }
    return next(ctx)
}

// Example usage
chain.Use(&ConditionalMiddleware{
    middleware: NewAuthMiddleware(),
    condition: func(ctx *Context) bool {
        return !strings.HasPrefix(ctx.Request.URL.Path, "/public/")
    },
})
```

### 5. Recovery Middleware

```go
type RecoveryMiddleware struct {
    logger *slog.Logger
}

func (m *RecoveryMiddleware) Process(ctx *Context, next MiddlewareFunc) error {
    defer func() {
        if r := recover(); r != nil {
            stack := debug.Stack()
            m.logger.Error("Middleware panic recovered",
                "error", r,
                "stack", string(stack),
                "url", ctx.Request.URL.String(),
            )
            
            if ctx.Response.Header().Get("Content-Type") == "" {
                ctx.Response.Header().Set("Content-Type", "text/html")
            }
            ctx.Response.WriteHeader(http.StatusInternalServerError)
            
            // Render error page
            errorPage := NewErrorPage(http.StatusInternalServerError)
            errorPage.Render(ctx.Response)
        }
    }()
    
    return next(ctx)
}
```

### 6. Example Implementation

```go
// Cache middleware example
type CacheMiddleware struct {
    cache  CacheInterface
    config CacheConfig
}

type CacheConfig struct {
    Enabled     bool
    Duration    time.Duration
    KeyPrefix   string
    IgnorePaths []string
}

func (m *CacheMiddleware) Process(ctx *Context, next MiddlewareFunc) error {
    if !m.config.Enabled {
        return next(ctx)
    }
    
    // Check if path should be cached
    for _, path := range m.config.IgnorePaths {
        if strings.HasPrefix(ctx.Request.URL.Path, path) {
            return next(ctx)
        }
    }
    
    // Generate cache key
    key := fmt.Sprintf("%s:%s:%s",
        m.config.KeyPrefix,
        ctx.Site.ID,
        ctx.Request.URL.Path,
    )
    
    // Check cache
    if cached, ok := m.cache.Get(key); ok {
        ctx.Response.Write(cached.([]byte))
        return nil
    }
    
    // Create response recorder
    recorder := httptest.NewRecorder()
    ctx.Response = recorder
    
    // Process request
    if err := next(ctx); err != nil {
        return err
    }
    
    // Cache response
    response := recorder.Result()
    body, _ := io.ReadAll(response.Body)
    m.cache.Set(key, body, m.config.Duration)
    
    // Write to original response
    for k, v := range response.Header {
        ctx.Response.Header()[k] = v
    }
    ctx.Response.WriteHeader(response.StatusCode)
    ctx.Response.Write(body)
    
    return nil
}
```

## Alternatives Considered

1. **Function-based Middleware**
   - Pros: Simpler implementation
   - Cons: Limited functionality, no configuration
   - Rejected: Need more features and flexibility

2. **Event-based System**
   - Pros: More decoupled
   - Cons: Complex flow control, harder to debug
   - Rejected: Direct middleware chain is more predictable

3. **Aspect-oriented Approach**
   - Pros: Clean separation of concerns
   - Cons: Complex implementation, runtime overhead
   - Rejected: Traditional middleware pattern is sufficient

## Implementation Plan

1. Phase 1: Core Enhancement (2 weeks)
   - Implement new interfaces
   - Add middleware chain management
   - Create basic metrics

2. Phase 2: Features (2 weeks)
   - Add conditional middleware
   - Implement recovery system
   - Create configuration system

3. Phase 3: Monitoring (1 week)
   - Add detailed metrics
   - Create dashboards
   - Add logging

4. Phase 4: Migration (2 weeks)
   - Update existing middleware
   - Add tests
   - Update documentation

## Risks and Mitigations

1. **Performance Impact**
   - Risk: Overhead from metrics and chaining
   - Mitigation: Benchmark-driven optimization

2. **Complexity**
   - Risk: System becomes hard to understand
   - Mitigation: Clear documentation, examples

3. **Migration**
   - Risk: Breaking existing middleware
   - Mitigation: Compatibility layer, gradual rollout

4. **Resource Usage**
   - Risk: Memory leaks from context
   - Mitigation: Context cleanup, monitoring 