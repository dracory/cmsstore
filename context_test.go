package cmsstore

import (
	"context"
	"net/http"
	"testing"
)

func TestRequestToContext(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/test?q=hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add to context
	ctx := RequestToContext(context.Background(), req)

	// Verify we can retrieve it
	retrieved := RequestFromContext(ctx)
	if retrieved == nil {
		t.Fatal("Expected to retrieve request from context, got nil")
	}

	// Verify it's the same request
	if retrieved.URL.Query().Get("q") != "hello" {
		t.Errorf("Expected query param 'q' to be 'hello', got: %s", retrieved.URL.Query().Get("q"))
	}
}

func TestRequestFromContext_NotFound(t *testing.T) {
	// Try to get request from empty context
	ctx := context.Background()
	req := RequestFromContext(ctx)

	if req != nil {
		t.Error("Expected nil when request not in context, got a request")
	}
}

func TestRequestFromContext_WrongType(t *testing.T) {
	// Add wrong type to context with same key (simulates collision)
	ctx := context.WithValue(context.Background(), httpRequestContextKey, "not a request")
	req := RequestFromContext(ctx)

	if req != nil {
		t.Error("Expected nil when wrong type in context, got a request")
	}
}

func TestRequestFromContext_NilRequest(t *testing.T) {
	// Add nil request to context
	ctx := context.WithValue(context.Background(), httpRequestContextKey, (*http.Request)(nil))
	req := RequestFromContext(ctx)

	if req != nil {
		t.Error("Expected nil when nil request in context, got a request")
	}
}

func TestRequestToContext_PreservesExistingContext(t *testing.T) {
	// Create a context with existing values
	type otherKey string
	ctx := context.WithValue(context.Background(), otherKey("existing"), "value")

	// Add request
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx = RequestToContext(ctx, req)

	// Verify existing value is preserved
	if ctx.Value(otherKey("existing")) != "value" {
		t.Error("Expected existing context values to be preserved")
	}

	// Verify request is accessible
	retrieved := RequestFromContext(ctx)
	if retrieved == nil {
		t.Error("Expected to retrieve request")
	}
}
