package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gouniverse/cmsstore"
)

// TestApplyMiddlewares_BeforeAndAfter tests that the before and after middlewares are applied correctly
func TestApplyMiddlewares_BeforeAndAfter(t *testing.T) {
	// Define test middlewares
	beforeMiddleware := &MockMiddleware{
		identifier:     "before_mw",
		name:           "Before Middleware",
		description:    "Modifies the content before it is served",
		middlewareType: "before",
		handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Before|"))
				next.ServeHTTP(w, r)
			})
		},
	}

	afterMiddleware := &MockMiddleware{
		identifier:     "after_mw",
		name:           "After Middleware",
		description:    "Modifies the content after it is served",
		middlewareType: "after",
		handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rec := httptest.NewRecorder()
				next.ServeHTTP(rec, r)
				w.Write([]byte(rec.Body.String() + "|After"))
			})
		},
	}

	middlewares := []cmsstore.MiddlewareInterface{beforeMiddleware, afterMiddleware}
	middlewaresBefore := []string{"before_mw"}
	middlewaresAfter := []string{"after_mw"}

	// Test request
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	// Apply middlewares
	pageContent := "PageContent"
	result := applyMiddlewares(recorder, req, middlewares, pageContent, middlewaresBefore, middlewaresAfter)

	// Expected output
	expectedOutput := "Before|PageContent|After"

	// Assert result
	if result != expectedOutput {
		t.Errorf("Expected output %q but got %q", expectedOutput, result)
	}
}

// TestApplyMiddlewares_MiddlewareNotFound ensures that a missing middleware returns an error
func TestApplyMiddlewares_MiddlewareNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	middlewares := []cmsstore.MiddlewareInterface{}

	// Apply middleware with a non-existent alias
	pageContent := "PageContent"
	result := applyMiddlewares(recorder, req, middlewares, pageContent, []string{"non_existent"}, nil)

	// Assert that result is empty due to error response
	if result != "" {
		t.Errorf("Expected empty result but got %q", result)
	}
	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected status %d but got %d", http.StatusNotFound, recorder.Code)
	}
}

// TestApplyMiddlewares_NilHandler ensures that a nil handler is handled properly
func TestApplyMiddlewares_NilHandler(t *testing.T) {
	nilMiddleware := &MockMiddleware{
		identifier:     "nil_mw",
		name:           "Nil Middleware",
		description:    "Has a nil handler",
		middlewareType: "before",
		handler:        nil,
	}

	middlewares := []cmsstore.MiddlewareInterface{nilMiddleware}
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	pageContent := "PageContent"
	result := applyMiddlewares(recorder, req, middlewares, pageContent, []string{"nil_mw"}, nil)

	// Assert that result is empty due to error response
	if result != "" {
		t.Errorf("Expected empty result but got %q", result)
	}
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d but got %d", http.StatusInternalServerError, recorder.Code)
	}
}

// TestApplyMiddlewares_StopExecution ensures that a middleware can stop execution
func TestApplyMiddlewares_StopExecution(t *testing.T) {
	stopMiddleware := &MockMiddleware{
		identifier:     "stop_mw",
		name:           "Stop Middleware",
		description:    "Stops execution",
		middlewareType: "before",
		handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Blocked", http.StatusForbidden)
			})
		},
	}

	middlewares := []cmsstore.MiddlewareInterface{stopMiddleware}
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	pageContent := "PageContent"
	result := applyMiddlewares(recorder, req, middlewares, pageContent, []string{"stop_mw"}, nil)

	// Assert that result is empty due to middleware stopping execution
	if result != "" {
		t.Errorf("Expected empty result but got %q", result)
	}
	if recorder.Code != http.StatusForbidden {
		t.Errorf("Expected status %d but got %d", http.StatusForbidden, recorder.Code)
	}
}

// TestApplyMiddlewares_EmptyMiddlewareList ensures it works without any middleware
func TestApplyMiddlewares_EmptyMiddlewareList(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	// middlewares := []cmsstore.MiddlewareInterface{}

	pageContent := "PageContent"
	result := applyMiddlewares(recorder, req, nil, pageContent, nil, nil)

	// Assert that content remains unchanged
	if result != pageContent {
		t.Errorf("Expected output %q but got %q", pageContent, result)
	}
}

// TestApplyMiddlewares_ModifyContent ensures middleware can modify page content
func TestApplyMiddlewares_ModifyContent(t *testing.T) {
	modifyMiddleware := &MockMiddleware{
		identifier:     "modify_mw",
		name:           "Modify Middleware",
		description:    "Changes page content",
		middlewareType: "before",
		handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ModifiedContent"))
			})
		},
	}

	middlewares := []cmsstore.MiddlewareInterface{modifyMiddleware}
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	pageContent := "PageContent"
	result := applyMiddlewares(recorder, req, middlewares, pageContent, nil, []string{"modify_mw"})

	// Assert that content is modified
	expectedOutput := "ModifiedContent"
	if result != expectedOutput {
		t.Errorf("Expected output %q but got %q", expectedOutput, result)
	}
}

// TestApplyMiddlewares_OrderMatters ensures that before and after middleware execute in the correct order
func TestApplyMiddlewares_OrderMatters(t *testing.T) {
	beforeMiddleware := &MockMiddleware{
		identifier:     "before_mw",
		name:           "Before Middleware",
		description:    "Adds before tag",
		middlewareType: "before",
		handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("[Before]"))
				next.ServeHTTP(w, r)
			})
		},
	}

	afterMiddleware := &MockMiddleware{
		identifier:     "after_mw",
		name:           "After Middleware",
		description:    "Adds after tag",
		middlewareType: "after",
		handler: func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				rec := httptest.NewRecorder()
				next.ServeHTTP(rec, r)
				w.Write([]byte(rec.Body.String() + "[After]"))
			})
		},
	}

	middlewares := []cmsstore.MiddlewareInterface{beforeMiddleware, afterMiddleware}
	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	pageContent := "PageContent"
	result := applyMiddlewares(recorder, req, middlewares, pageContent, []string{"before_mw"}, []string{"after_mw"})

	// Expected output
	expectedOutput := "[Before]PageContent[After]"

	// Assert order
	if result != expectedOutput {
		t.Errorf("Expected output %q but got %q", expectedOutput, result)
	}
}
