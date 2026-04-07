package cmsstore

import (
	"context"
	"net/http"
	"sync"
)

// contextKey is a private type to avoid key collisions with other packages
type contextKey string

const (
	httpRequestContextKey contextKey = "http_request"
	varsContextKey        contextKey = "vars"
)

// RequestFromContext retrieves the *http.Request from the context if it was
// previously added using RequestToContext. This is useful for custom block types
// that need to access request data (e.g., query parameters, headers).
//
// Example usage:
//
//	func (b *myBlockType) Render(ctx context.Context, block cmsstore.BlockInterface) (string, error) {
//		req := cmsstore.RequestFromContext(ctx)
//		if req != nil {
//			queryParam := req.URL.Query().Get("q")
//			// ... use queryParam
//		}
//		// ... render block
//	}
func RequestFromContext(ctx context.Context) *http.Request {
	if req, ok := ctx.Value(httpRequestContextKey).(*http.Request); ok {
		return req
	}
	return nil
}

// RequestToContext adds the *http.Request to the context. This is called
// internally by the frontend when rendering blocks to ensure the request
// is available to custom block types via RequestFromContext.
func RequestToContext(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, httpRequestContextKey, req)
}

// VarsContext stores custom variables set by blocks during rendering.
// Blocks can set arbitrary variables that will be replaced in the final content.
//
// Example usage in a custom block:
//
//	func (b *BlogBlockType) Render(ctx context.Context, block BlockInterface, opts ...RenderOption) (string, error) {
//		if vars := cmsstore.VarsFromContext(ctx); vars != nil {
//			vars.Set("blog_title", "My Blog Post")
//			vars.Set("author_name", "John Doe")
//			vars.Set("publish_date", "2026-04-07")
//		}
//		return html, nil
//	}
//
// Variables can then be referenced in page/template content as [[blog_title]], [[author_name]], etc.
type VarsContext struct {
	vars map[string]string
	mu   sync.RWMutex
}

// NewVarsContext creates a new variable context for storing custom variables.
func NewVarsContext() *VarsContext {
	return &VarsContext{
		vars: make(map[string]string),
	}
}

// Set stores a variable that can be referenced as [[key]] in content.
// If the key already exists, it will be overwritten.
//
// Example:
//
//	vars.Set("blog_title", "Understanding Go Contexts")
//	vars.Set("product_price", "$99.99")
func (v *VarsContext) Set(key, value string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.vars[key] = value
}

// Get retrieves a variable value. Returns the value and true if found,
// or empty string and false if not found.
func (v *VarsContext) Get(key string) (string, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	val, ok := v.vars[key]
	return val, ok
}

// All returns all variables as a map. Creates a copy for thread safety.
func (v *VarsContext) All() map[string]string {
	v.mu.RLock()
	defer v.mu.RUnlock()

	result := make(map[string]string, len(v.vars))
	for k, v := range v.vars {
		result[k] = v
	}
	return result
}

// WithVarsContext adds a VarsContext to the context. This is called internally
// by the frontend during rendering to enable blocks to set custom variables.
func WithVarsContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, varsContextKey, NewVarsContext())
}

// VarsFromContext retrieves the VarsContext from context. Returns nil if
// no VarsContext was added to the context.
//
// Example usage in a custom block:
//
//	if vars := cmsstore.VarsFromContext(ctx); vars != nil {
//		vars.Set("my_variable", "my value")
//	}
func VarsFromContext(ctx context.Context) *VarsContext {
	if vars, ok := ctx.Value(varsContextKey).(*VarsContext); ok {
		return vars
	}
	return nil
}
