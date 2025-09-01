package frontend

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/dracory/cmsstore"
)

// applyMiddlewares applies middlewares to the page content
func (frontend *frontend) applyMiddlewares(
	w http.ResponseWriter,
	r *http.Request,
	pageContent string,
	pageMiddlewareIdentifiersBefore []string,
	pageMiddlewareIdentifiersAfter []string,
) string {
	middlewares := frontend.store.Middlewares()
	return applyMiddlewares(w, r, middlewares, pageContent, pageMiddlewareIdentifiersBefore, pageMiddlewareIdentifiersAfter)
}

// applyMiddlewares applies the given middlewares to the page content.
//
// Middlewares can be of two types: "before" and "after":
// - "before" middlewares are executed before the page content is processed.
// - "after" middlewares are executed after the page content is processed.
//
// Each middleware can modify the request, response, or page content. If a "before" middleware
// writes a response and does not call the next handler, execution is halted.
//
// Parameters:
// - w: The HTTP response writer.
// - r: The HTTP request.
// - pageContent: The original page content to be processed.
// - middlewares: A list of available middleware instances.
// - middlewareAliases: A list of middleware aliases to be applied in order.
//
// Returns:
// - The processed page content after all applicable middlewares have been applied.
func applyMiddlewares(
	w http.ResponseWriter,
	r *http.Request,
	middlewares []cmsstore.MiddlewareInterface,
	pageContent string,
	pageMiddlewareIdentifiersBefore []string,
	pageMiddlewareIdentifiersAfter []string,
) string {
	var beforeHandlers []func(http.Handler) http.Handler
	var afterHandlers []func(http.Handler) http.Handler

	// Merge the before and after middleware identifiers
	middlewareIdentifiers := append(pageMiddlewareIdentifiersBefore, pageMiddlewareIdentifiersAfter...)

	// Retrieve and categorize the middlewares based on their alias
	for _, identifier := range middlewareIdentifiers {
		var foundMiddleware cmsstore.MiddlewareInterface
		for _, mw := range middlewares {
			if mw.Identifier() == identifier {
				foundMiddleware = mw
				break
			}
		}

		// If middleware is not found, return an error response and stop processing
		if foundMiddleware == nil {
			http.Error(w, "Middleware not found: "+identifier, http.StatusNotFound)
			return ""
		}

		handler := foundMiddleware.Handler()
		if handler == nil {
			http.Error(w, "Middleware handler is nil: "+identifier, http.StatusInternalServerError)
			return ""
		}

		// Categorize middleware based on type
		switch foundMiddleware.Type() {
		case cmsstore.MIDDLEWARE_TYPE_BEFORE:
			beforeHandlers = append(beforeHandlers, handler)
		case cmsstore.MIDDLEWARE_TYPE_AFTER:
			afterHandlers = append(afterHandlers, handler)
		}
	}

	// Define the final handler that serves the original page content
	var finalHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, pageContent)
	})

	// Apply "before" middlewares in reverse order to maintain expected execution flow
	for i := len(beforeHandlers) - 1; i >= 0; i-- {
		finalHandler = beforeHandlers[i](finalHandler)
	}

	// Capture the response before applying "after" middlewares
	recorder := httptest.NewRecorder()
	finalHandler.ServeHTTP(recorder, r)

	// If a "before" middleware wrote a response (e.g., an authentication check failed),
	// execution stops and the response is immediately returned.
	if recorder.Code != http.StatusOK {
		w.WriteHeader(recorder.Code)
		w.Write(recorder.Body.Bytes())
		return ""
	}

	// Get the modified content after "before" middlewares have executed
	modifiedContent := recorder.Body.String()

	// Define a handler for serving the modified content
	finalHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, modifiedContent)
	})

	// Apply "after" middlewares in the given order
	for _, handler := range afterHandlers {
		finalHandler = handler(finalHandler)
	}

	// Capture the final response after applying "after" middlewares
	recorder = httptest.NewRecorder()
	finalHandler.ServeHTTP(recorder, r)

	// Return the final modified page content
	return recorder.Body.String()
}
