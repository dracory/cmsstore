package cmsstore

import (
	"context"
	"net/http"
)

// contextKey is a private type to avoid key collisions with other packages
type contextKey string

const httpRequestContextKey contextKey = "http_request"

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
